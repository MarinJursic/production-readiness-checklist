package trust

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func signedFixture(t *testing.T) (LoadedStore, Signature, time.Time) {
	t.Helper()
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	issuedAt := time.Date(2026, 8, 23, 12, 0, 0, 0, time.UTC)
	storeDocument := Store{
		SchemaVersion: StoreSchema, ID: "prc-release-keys", Revision: 1,
		Keys: []Key{{
			ID: "release-2026", Algorithm: AlgorithmEd25519,
			PublicKey: base64.StdEncoding.EncodeToString(publicKey), Scopes: []string{"pack", "adapter-registry"},
			Status: "active", NotBefore: issuedAt.Add(-time.Hour), NotAfter: issuedAt.Add(24 * time.Hour),
		}},
	}
	directory := t.TempDir()
	storePath := filepath.Join(directory, "trust-store.json")
	payload, err := json.Marshal(storeDocument)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(storePath, payload, 0o600); err != nil {
		t.Fatal(err)
	}
	store, err := LoadStore(storePath)
	if err != nil {
		t.Fatal(err)
	}
	signature := Signature{
		SchemaVersion: SignatureSchema, ArtifactKind: "pack", ArtifactID: "prc.pack.example@1.0",
		SHA256: strings.Repeat("a", 64), KeyID: "release-2026", Algorithm: AlgorithmEd25519, IssuedAt: issuedAt,
	}
	signingPayload, err := SigningPayload(signature)
	if err != nil {
		t.Fatal(err)
	}
	signature.Value = base64.StdEncoding.EncodeToString(ed25519.Sign(privateKey, signingPayload))
	return store, signature, issuedAt.Add(time.Hour)
}

func TestVerifyAcceptsScopedCurrentSignature(t *testing.T) {
	store, signature, verifiedAt := signedFixture(t)
	verification, err := Verify(store, signature, signature.ArtifactKind, signature.ArtifactID, signature.SHA256, verifiedAt)
	if err != nil {
		t.Fatal(err)
	}
	if !verification.Verified || verification.SchemaVersion != VerificationSchema ||
		verification.TrustStoreDigest != store.Digest || len(verification.SignatureDigest) != 64 {
		t.Fatalf("verification = %+v", verification)
	}
}

func TestVerifyRejectsTamperingRevocationScopeAndTime(t *testing.T) {
	tests := map[string]func(*LoadedStore, *Signature, *time.Time){
		"subject tampering": func(_ *LoadedStore, signature *Signature, _ *time.Time) {
			signature.SHA256 = strings.Repeat("b", 64)
		},
		"signature tampering": func(_ *LoadedStore, signature *Signature, _ *time.Time) {
			value, _ := base64.StdEncoding.DecodeString(signature.Value)
			value[0] ^= 0xff
			signature.Value = base64.StdEncoding.EncodeToString(value)
		},
		"revocation": func(store *LoadedStore, _ *Signature, _ *time.Time) {
			store.Store.Keys[0].Status = "revoked"
			store.Store.Keys[0].Reason = "fixture revocation"
			refreshStoreDigest(t, store)
		},
		"wrong scope": func(store *LoadedStore, _ *Signature, _ *time.Time) {
			store.Store.Keys[0].Scopes = []string{"adapter-registry"}
			refreshStoreDigest(t, store)
		},
		"expired": func(store *LoadedStore, _ *Signature, verifiedAt *time.Time) {
			*verifiedAt = store.Store.Keys[0].NotAfter.Add(time.Second)
		},
		"non-UTC verification": func(_ *LoadedStore, _ *Signature, verifiedAt *time.Time) {
			*verifiedAt = time.Date(2026, 8, 23, 14, 0, 0, 0, time.FixedZone("fixture", 3600))
		},
	}
	for name, mutate := range tests {
		t.Run(name, func(t *testing.T) {
			store, signature, verifiedAt := signedFixture(t)
			mutate(&store, &signature, &verifiedAt)
			if _, err := Verify(store, signature, "pack", "prc.pack.example@1.0", strings.Repeat("a", 64), verifiedAt); err == nil {
				t.Fatal("invalid signature state was accepted")
			}
		})
	}
}

func TestTrustFilesRejectYAMLAliases(t *testing.T) {
	data := []byte("schema_version: &schema prc.trust-store/v0.1\nid: *schema\nrevision: 1\nkeys: []\n")
	var store Store
	if err := decodeStrict(data, &store, "trust store"); err == nil || !strings.Contains(err.Error(), "aliases") {
		t.Fatalf("alias error = %v", err)
	}
}

func TestLoadedTrustStoreCannotBeMutatedAfterValidation(t *testing.T) {
	store, signature, verifiedAt := signedFixture(t)
	store.Store.Revision++
	if _, err := Verify(store, signature, "pack", signature.ArtifactID, signature.SHA256, verifiedAt); err == nil ||
		!strings.Contains(err.Error(), "digest") {
		t.Fatalf("mutated trust store error = %v", err)
	}
}

func refreshStoreDigest(t *testing.T, store *LoadedStore) {
	t.Helper()
	payload, err := json.Marshal(store.Store)
	if err != nil {
		t.Fatal(err)
	}
	directory := t.TempDir()
	path := filepath.Join(directory, "store.json")
	if err := os.WriteFile(path, payload, 0o600); err != nil {
		t.Fatal(err)
	}
	loaded, err := LoadStore(path)
	if err != nil {
		t.Fatal(err)
	}
	*store = loaded
}
