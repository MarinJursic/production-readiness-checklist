// Package trust verifies offline publisher signatures for immutable scanner
// artifacts. It deliberately contains no key discovery, network fetching, or
// private-key handling.
package trust

import (
	"bytes"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	StoreSchema        = "prc.trust-store/v0.1"
	SignatureSchema    = "prc.signature/v0.1"
	VerificationSchema = "prc.signature-verification/v0.1"
	AlgorithmEd25519   = "ed25519"
)

const maximumTrustFileBytes = 1024 * 1024

var (
	identifierPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9.-]{0,127}$`)
	digestPattern     = regexp.MustCompile(`^[0-9a-f]{64}$`)
	artifactIDPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9./@_-]{0,255}$`)
)

type Key struct {
	ID        string    `json:"id" yaml:"id"`
	Algorithm string    `json:"algorithm" yaml:"algorithm"`
	PublicKey string    `json:"public_key" yaml:"public_key"`
	Scopes    []string  `json:"scopes" yaml:"scopes"`
	Status    string    `json:"status" yaml:"status"`
	NotBefore time.Time `json:"not_before" yaml:"not_before"`
	NotAfter  time.Time `json:"not_after" yaml:"not_after"`
	Reason    string    `json:"reason,omitempty" yaml:"reason,omitempty"`
}

type Store struct {
	SchemaVersion string `json:"schema_version" yaml:"schema_version"`
	ID            string `json:"id" yaml:"id"`
	Revision      int    `json:"revision" yaml:"revision"`
	Keys          []Key  `json:"keys" yaml:"keys"`
}

type LoadedStore struct {
	Store  Store
	Digest string
}

type Signature struct {
	SchemaVersion string    `json:"schema_version" yaml:"schema_version"`
	ArtifactKind  string    `json:"artifact_kind" yaml:"artifact_kind"`
	ArtifactID    string    `json:"artifact_id" yaml:"artifact_id"`
	SHA256        string    `json:"sha256" yaml:"sha256"`
	KeyID         string    `json:"key_id" yaml:"key_id"`
	Algorithm     string    `json:"algorithm" yaml:"algorithm"`
	IssuedAt      time.Time `json:"issued_at" yaml:"issued_at"`
	Value         string    `json:"signature" yaml:"signature"`
}

type Verification struct {
	SchemaVersion    string    `json:"schema_version"`
	ArtifactKind     string    `json:"artifact_kind"`
	ArtifactID       string    `json:"artifact_id"`
	SHA256           string    `json:"sha256"`
	KeyID            string    `json:"key_id"`
	Algorithm        string    `json:"algorithm"`
	IssuedAt         time.Time `json:"issued_at"`
	VerifiedAt       time.Time `json:"verified_at"`
	TrustStoreID     string    `json:"trust_store_id"`
	TrustStoreDigest string    `json:"trust_store_digest"`
	SignatureDigest  string    `json:"signature_digest"`
	Verified         bool      `json:"verified"`
}

type signingPayload struct {
	Domain        string    `json:"domain"`
	SchemaVersion string    `json:"schema_version"`
	ArtifactKind  string    `json:"artifact_kind"`
	ArtifactID    string    `json:"artifact_id"`
	SHA256        string    `json:"sha256"`
	KeyID         string    `json:"key_id"`
	Algorithm     string    `json:"algorithm"`
	IssuedAt      time.Time `json:"issued_at"`
}

func LoadStore(path string) (LoadedStore, error) {
	data, err := readBoundedFile(path, "trust store")
	if err != nil {
		return LoadedStore{}, err
	}
	var store Store
	if err := decodeStrict(data, &store, "trust store"); err != nil {
		return LoadedStore{}, err
	}
	if err := normalizeStore(&store); err != nil {
		return LoadedStore{}, err
	}
	payload, err := json.Marshal(store)
	if err != nil {
		return LoadedStore{}, fmt.Errorf("encode trust store identity: %w", err)
	}
	digest := sha256.Sum256(payload)
	return LoadedStore{Store: store, Digest: hex.EncodeToString(digest[:])}, nil
}

func LoadSignature(path string) (Signature, error) {
	data, err := readBoundedFile(path, "signature envelope")
	if err != nil {
		return Signature{}, err
	}
	var signature Signature
	if err := decodeStrict(data, &signature, "signature envelope"); err != nil {
		return Signature{}, err
	}
	if err := signature.Validate(); err != nil {
		return Signature{}, err
	}
	return signature, nil
}

func (signature Signature) Validate() error {
	if signature.SchemaVersion != SignatureSchema || !artifactKind(signature.ArtifactKind) ||
		!artifactIDPattern.MatchString(signature.ArtifactID) || !digestPattern.MatchString(signature.SHA256) ||
		!identifierPattern.MatchString(signature.KeyID) || signature.Algorithm != AlgorithmEd25519 ||
		!utcTime(signature.IssuedAt) {
		return fmt.Errorf("signature envelope has an unsupported or invalid identity")
	}
	value, err := base64.StdEncoding.Strict().DecodeString(signature.Value)
	if err != nil || len(value) != ed25519.SignatureSize {
		return fmt.Errorf("signature envelope requires a canonical Ed25519 signature")
	}
	return nil
}

func Verify(
	store LoadedStore,
	signature Signature,
	artifactKindValue string,
	artifactID string,
	digest string,
	verifiedAt time.Time,
) (Verification, error) {
	if err := normalizeAndVerifyLoadedStore(store); err != nil {
		return Verification{}, err
	}
	if err := signature.Validate(); err != nil {
		return Verification{}, err
	}
	if !utcTime(verifiedAt) {
		return Verification{}, fmt.Errorf("signature verification time must be a nonzero UTC timestamp")
	}
	if signature.ArtifactKind != artifactKindValue || signature.ArtifactID != artifactID || signature.SHA256 != digest {
		return Verification{}, fmt.Errorf("signature subject does not match the requested artifact identity")
	}
	var selected *Key
	for index := range store.Store.Keys {
		if store.Store.Keys[index].ID == signature.KeyID {
			selected = &store.Store.Keys[index]
			break
		}
	}
	if selected == nil {
		return Verification{}, fmt.Errorf("signature key %s is not present in trust store %s", signature.KeyID, store.Store.ID)
	}
	if selected.Status == "revoked" {
		return Verification{}, fmt.Errorf("signature key %s is revoked: %s", selected.ID, selected.Reason)
	}
	if !contains(selected.Scopes, signature.ArtifactKind) {
		return Verification{}, fmt.Errorf("signature key %s is not authorized for %s artifacts", selected.ID, signature.ArtifactKind)
	}
	if signature.IssuedAt.Before(selected.NotBefore) || signature.IssuedAt.After(selected.NotAfter) ||
		verifiedAt.Before(selected.NotBefore) || verifiedAt.After(selected.NotAfter) || signature.IssuedAt.After(verifiedAt) {
		return Verification{}, fmt.Errorf("signature or verification time is outside key %s validity", selected.ID)
	}
	publicKey, err := base64.StdEncoding.Strict().DecodeString(selected.PublicKey)
	if err != nil || len(publicKey) != ed25519.PublicKeySize {
		return Verification{}, fmt.Errorf("trust key %s has an invalid public key", selected.ID)
	}
	value, _ := base64.StdEncoding.Strict().DecodeString(signature.Value)
	payload, err := SigningPayload(signature)
	if err != nil {
		return Verification{}, err
	}
	if !ed25519.Verify(ed25519.PublicKey(publicKey), payload, value) {
		return Verification{}, fmt.Errorf("signature verification failed for %s", signature.ArtifactID)
	}
	signaturePayload, err := json.Marshal(signature)
	if err != nil {
		return Verification{}, fmt.Errorf("encode signature identity: %w", err)
	}
	signatureDigest := sha256.Sum256(signaturePayload)
	return Verification{
		SchemaVersion: VerificationSchema, ArtifactKind: signature.ArtifactKind, ArtifactID: signature.ArtifactID,
		SHA256: signature.SHA256, KeyID: signature.KeyID, Algorithm: signature.Algorithm,
		IssuedAt: signature.IssuedAt.UTC(), VerifiedAt: verifiedAt.UTC(), TrustStoreID: store.Store.ID,
		TrustStoreDigest: store.Digest, SignatureDigest: hex.EncodeToString(signatureDigest[:]), Verified: true,
	}, nil
}

func SigningPayload(signature Signature) ([]byte, error) {
	if err := signature.ValidateWithoutValue(); err != nil {
		return nil, err
	}
	return json.Marshal(signingPayload{
		Domain: "prc.publisher-signature/v1", SchemaVersion: signature.SchemaVersion,
		ArtifactKind: signature.ArtifactKind, ArtifactID: signature.ArtifactID, SHA256: signature.SHA256,
		KeyID: signature.KeyID, Algorithm: signature.Algorithm, IssuedAt: signature.IssuedAt.UTC(),
	})
}

func (signature Signature) ValidateWithoutValue() error {
	value := signature.Value
	signature.Value = base64.StdEncoding.EncodeToString(make([]byte, ed25519.SignatureSize))
	err := signature.Validate()
	signature.Value = value
	return err
}

func normalizeStore(store *Store) error {
	if store.SchemaVersion != StoreSchema || !identifierPattern.MatchString(store.ID) || store.Revision < 1 ||
		len(store.Keys) == 0 || len(store.Keys) > 256 {
		return fmt.Errorf("trust store requires a supported schema, valid identity, revision, and bounded keys")
	}
	seenIDs, seenPublicKeys := map[string]bool{}, map[string]bool{}
	for index := range store.Keys {
		key := &store.Keys[index]
		if !identifierPattern.MatchString(key.ID) || key.Algorithm != AlgorithmEd25519 ||
			(key.Status != "active" && key.Status != "revoked") || !utcTime(key.NotBefore) || !utcTime(key.NotAfter) ||
			!key.NotBefore.Before(key.NotAfter) {
			return fmt.Errorf("trust store key %q has invalid identity, algorithm, status, or validity", key.ID)
		}
		publicKey, err := base64.StdEncoding.Strict().DecodeString(key.PublicKey)
		if err != nil || len(publicKey) != ed25519.PublicKeySize {
			return fmt.Errorf("trust store key %s has an invalid Ed25519 public key", key.ID)
		}
		if seenIDs[key.ID] || seenPublicKeys[key.PublicKey] {
			return fmt.Errorf("trust store contains a duplicate key ID or public key")
		}
		seenIDs[key.ID], seenPublicKeys[key.PublicKey] = true, true
		if len(key.Scopes) == 0 || len(key.Scopes) > 16 {
			return fmt.Errorf("trust store key %s requires bounded scopes", key.ID)
		}
		sort.Strings(key.Scopes)
		for scopeIndex, scope := range key.Scopes {
			if !artifactKind(scope) || (scopeIndex > 0 && scope == key.Scopes[scopeIndex-1]) {
				return fmt.Errorf("trust store key %s has invalid or duplicate scopes", key.ID)
			}
		}
		if (key.Status == "active" && key.Reason != "") ||
			(key.Status == "revoked" && strings.TrimSpace(key.Reason) == "") {
			return fmt.Errorf("trust store key %s has an inconsistent revocation reason", key.ID)
		}
		key.NotBefore = key.NotBefore.UTC()
		key.NotAfter = key.NotAfter.UTC()
	}
	sort.Slice(store.Keys, func(left, right int) bool { return store.Keys[left].ID < store.Keys[right].ID })
	return nil
}

func normalizeAndVerifyLoadedStore(store LoadedStore) error {
	copyOfStore := store.Store
	copyOfStore.Keys = append([]Key(nil), store.Store.Keys...)
	for index := range copyOfStore.Keys {
		copyOfStore.Keys[index].Scopes = append([]string(nil), copyOfStore.Keys[index].Scopes...)
	}
	if err := normalizeStore(&copyOfStore); err != nil {
		return err
	}
	payload, err := json.Marshal(copyOfStore)
	if err != nil {
		return err
	}
	digest := sha256.Sum256(payload)
	if store.Digest != hex.EncodeToString(digest[:]) {
		return fmt.Errorf("trust store digest does not match its canonical content")
	}
	return nil
}

func artifactKind(value string) bool {
	return value == "pack" || value == "adapter-registry" || value == "catalog-bundle" || value == "risk-exception"
}

func utcTime(value time.Time) bool {
	if value.IsZero() {
		return false
	}
	_, offset := value.Zone()
	return offset == 0
}

func contains(values []string, wanted string) bool {
	for _, value := range values {
		if value == wanted {
			return true
		}
	}
	return false
}

func readBoundedFile(path, label string) ([]byte, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return nil, fmt.Errorf("inspect %s: %w", label, err)
	}
	if info.Mode()&os.ModeSymlink != 0 || !info.Mode().IsRegular() || info.Size() > maximumTrustFileBytes {
		return nil, fmt.Errorf("%s must be a non-symlink regular file no larger than 1 MiB", label)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", label, err)
	}
	return data, nil
}

func decodeStrict(data []byte, destination any, label string) error {
	var syntax yaml.Node
	if err := yaml.Unmarshal(data, &syntax); err != nil {
		return fmt.Errorf("decode %s syntax: %w", label, err)
	}
	if len(syntax.Content) != 1 {
		return fmt.Errorf("%s requires one YAML document", label)
	}
	if err := rejectAmbiguousYAML(syntax.Content[0], label); err != nil {
		return err
	}
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)
	if err := decoder.Decode(destination); err != nil {
		return fmt.Errorf("decode %s: %w", label, err)
	}
	var trailing any
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		if err == nil {
			return fmt.Errorf("%s contains more than one YAML document", label)
		}
		return fmt.Errorf("decode trailing %s content: %w", label, err)
	}
	return nil
}

func rejectAmbiguousYAML(node *yaml.Node, label string) error {
	if node.Kind == yaml.AliasNode || node.Tag == "!!null" || node.Tag == "!!merge" {
		return fmt.Errorf("%s cannot contain aliases, merge keys, or null values", label)
	}
	if node.Kind == yaml.MappingNode {
		seen := map[string]bool{}
		for index := 0; index < len(node.Content); index += 2 {
			key := node.Content[index]
			if key.Kind != yaml.ScalarNode || seen[key.Value] {
				return fmt.Errorf("%s contains a non-scalar or duplicate mapping key", label)
			}
			seen[key.Value] = true
		}
	}
	for _, child := range node.Content {
		if err := rejectAmbiguousYAML(child, label); err != nil {
			return err
		}
	}
	return nil
}
