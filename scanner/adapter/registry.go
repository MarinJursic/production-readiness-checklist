package adapter

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
	"gopkg.in/yaml.v3"
)

const (
	RegistrySchema       = "prc.adapter-registry/v0.1"
	RegistryReportSchema = "prc.adapter-registry-report/v0.1"
)

var registryIDPattern = regexp.MustCompile(`^prc\.adapter-registry\.[a-z0-9.-]+@[0-9]+\.[0-9]+$`)

type RegistryEntry struct {
	AdapterID      string `json:"adapter_id" yaml:"adapter_id"`
	ManifestPath   string `json:"manifest_path" yaml:"manifest_path"`
	ManifestSHA256 string `json:"manifest_sha256" yaml:"manifest_sha256"`
	PublisherID    string `json:"publisher_id" yaml:"publisher_id"`
	Trust          string `json:"trust" yaml:"trust"`
	Status         string `json:"status" yaml:"status"`
	Reason         string `json:"reason,omitempty" yaml:"reason,omitempty"`
}

type RegistryDocument struct {
	SchemaVersion string          `json:"schema_version" yaml:"schema_version"`
	ID            string          `json:"id" yaml:"id"`
	Revision      int             `json:"revision" yaml:"revision"`
	Entries       []RegistryEntry `json:"entries" yaml:"entries"`
}

type Registry struct {
	SchemaVersion string
	ID            string
	Revision      int
	Digest        string
	Entries       []RegistryEntry
	manifests     map[string]Manifest
}

type RegistryReport struct {
	SchemaVersion string           `json:"schema_version"`
	Registry      RegistryDocument `json:"registry"`
	Digest        string           `json:"digest"`
}

type RegistryPolicy struct {
	AllowedTrust    []string
	AllowDeprecated bool
}

type ResolvedAdapter struct {
	Entry      RegistryEntry           `json:"entry"`
	Manifest   Manifest                `json:"manifest"`
	Resolution model.AdapterResolution `json:"resolution"`
}

func DefaultRegistryPolicy() RegistryPolicy {
	return RegistryPolicy{AllowedTrust: []string{"first-party-sandboxed", "verified-community"}}
}

func (registry Registry) Report() RegistryReport {
	return RegistryReport{
		SchemaVersion: RegistryReportSchema,
		Registry: RegistryDocument{
			SchemaVersion: registry.SchemaVersion, ID: registry.ID,
			Revision: registry.Revision, Entries: append([]RegistryEntry(nil), registry.Entries...),
		},
		Digest: registry.Digest,
	}
}

func LoadRegistry(path string) (Registry, error) {
	data, base, err := readRegistryFile(path)
	if err != nil {
		return Registry{}, err
	}
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)
	var document RegistryDocument
	if err := decoder.Decode(&document); err != nil {
		return Registry{}, fmt.Errorf("decode adapter registry: %w", err)
	}
	var trailing any
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		if err == nil {
			return Registry{}, fmt.Errorf("adapter registry contains more than one YAML document")
		}
		return Registry{}, fmt.Errorf("decode trailing adapter registry content: %w", err)
	}
	if err := validateRegistryDocument(document); err != nil {
		return Registry{}, err
	}

	entries := append([]RegistryEntry(nil), document.Entries...)
	sort.Slice(entries, func(left, right int) bool { return entries[left].AdapterID < entries[right].AdapterID })
	registry := Registry{
		SchemaVersion: document.SchemaVersion, ID: document.ID, Revision: document.Revision,
		Entries: entries, manifests: map[string]Manifest{},
	}
	seenIDs, seenDigests := map[string]bool{}, map[string]bool{}
	for _, entry := range entries {
		if seenIDs[entry.AdapterID] || seenDigests[entry.ManifestSHA256] {
			return Registry{}, fmt.Errorf("adapter registry contains a duplicate adapter ID or manifest digest")
		}
		seenIDs[entry.AdapterID], seenDigests[entry.ManifestSHA256] = true, true
		if entry.Status == "revoked" {
			continue
		}
		manifestPath, err := resolveRegistryManifest(base, entry.ManifestPath)
		if err != nil {
			return Registry{}, fmt.Errorf("adapter %s manifest path: %w", entry.AdapterID, err)
		}
		manifest, err := LoadManifest(manifestPath)
		if err != nil {
			return Registry{}, fmt.Errorf("load registered adapter %s: %w", entry.AdapterID, err)
		}
		digest, err := ManifestDigest(manifest)
		if err != nil {
			return Registry{}, err
		}
		if manifest.ID != entry.AdapterID || digest != entry.ManifestSHA256 || manifest.Publisher.ID != entry.PublisherID {
			return Registry{}, fmt.Errorf("registered adapter %s does not match its ID, digest, or publisher pin", entry.AdapterID)
		}
		if (entry.Status == "active" && manifest.Maintenance != "active") ||
			(entry.Status == "deprecated" && manifest.Maintenance != "deprecated") {
			return Registry{}, fmt.Errorf("registered adapter %s lifecycle does not match its manifest", entry.AdapterID)
		}
		if !manifest.SupportsEngine(model.EngineVersion) {
			return Registry{}, fmt.Errorf("registered adapter %s does not support engine API %s", entry.AdapterID, model.EngineVersion)
		}
		registry.manifests[entry.AdapterID] = manifest
	}
	identity := RegistryDocument{
		SchemaVersion: registry.SchemaVersion, ID: registry.ID, Revision: registry.Revision, Entries: registry.Entries,
	}
	payload, err := json.Marshal(identity)
	if err != nil {
		return Registry{}, fmt.Errorf("encode adapter registry identity: %w", err)
	}
	digest := sha256.Sum256(payload)
	registry.Digest = hex.EncodeToString(digest[:])
	return registry, nil
}

func (registry Registry) Resolve(
	adapterID string,
	expectedDigest string,
	requiredObservationKinds []string,
	policy RegistryPolicy,
) (ResolvedAdapter, error) {
	var entry *RegistryEntry
	for index := range registry.Entries {
		if registry.Entries[index].AdapterID == adapterID {
			entry = &registry.Entries[index]
			break
		}
	}
	if entry == nil {
		return ResolvedAdapter{}, fmt.Errorf("adapter %s is not present in registry %s", adapterID, registry.ID)
	}
	if entry.Status == "revoked" {
		return ResolvedAdapter{}, fmt.Errorf("adapter %s is revoked: %s", adapterID, entry.Reason)
	}
	if entry.Status == "deprecated" && !policy.AllowDeprecated {
		return ResolvedAdapter{}, fmt.Errorf("adapter %s is deprecated and denied by policy: %s", adapterID, entry.Reason)
	}
	if expectedDigest != "" && expectedDigest != entry.ManifestSHA256 {
		return ResolvedAdapter{}, fmt.Errorf("adapter %s registry digest does not match the required manifest digest", adapterID)
	}
	if !stringMember(policy.AllowedTrust, entry.Trust) {
		return ResolvedAdapter{}, fmt.Errorf("adapter %s trust level %s is denied by registry policy", adapterID, entry.Trust)
	}
	manifest, ok := registry.manifests[adapterID]
	if !ok {
		return ResolvedAdapter{}, fmt.Errorf("adapter %s has no resolvable active manifest", adapterID)
	}
	for _, kind := range requiredObservationKinds {
		if !stringMember(manifest.ObservationKinds, kind) {
			return ResolvedAdapter{}, fmt.Errorf("adapter %s does not declare required observation kind %s", adapterID, kind)
		}
	}
	return ResolvedAdapter{
		Entry: *entry, Manifest: manifest,
		Resolution: model.AdapterResolution{
			Source: ResolutionSourceRegistry, PublisherID: entry.PublisherID, Trust: entry.Trust,
			RegistryID: registry.ID, RegistryRevision: registry.Revision, RegistryDigest: registry.Digest,
		},
	}, nil
}

func readRegistryFile(path string) ([]byte, string, error) {
	absolute, err := filepath.Abs(path)
	if err != nil {
		return nil, "", fmt.Errorf("resolve adapter registry path: %w", err)
	}
	info, err := os.Lstat(absolute)
	if err != nil {
		return nil, "", fmt.Errorf("inspect adapter registry: %w", err)
	}
	if info.Mode()&os.ModeSymlink != 0 || !info.Mode().IsRegular() || info.Size() > 1024*1024 {
		return nil, "", fmt.Errorf("adapter registry must be a non-symlink regular file no larger than 1 MiB")
	}
	data, err := os.ReadFile(absolute)
	if err != nil {
		return nil, "", fmt.Errorf("read adapter registry: %w", err)
	}
	base, err := filepath.EvalSymlinks(filepath.Dir(absolute))
	if err != nil {
		return nil, "", fmt.Errorf("resolve adapter registry directory: %w", err)
	}
	return data, base, nil
}

func validateRegistryDocument(document RegistryDocument) error {
	if document.SchemaVersion != RegistrySchema {
		return fmt.Errorf("unsupported adapter registry schema %q", document.SchemaVersion)
	}
	if !registryIDPattern.MatchString(document.ID) || document.Revision < 1 {
		return fmt.Errorf("adapter registry requires a valid ID and positive revision")
	}
	if len(document.Entries) == 0 || len(document.Entries) > 4096 {
		return fmt.Errorf("adapter registry must contain between 1 and 4096 entries")
	}
	for _, entry := range document.Entries {
		if !adapterIDPattern.MatchString(entry.AdapterID) || !hexDigestPattern.MatchString(entry.ManifestSHA256) ||
			!publisherIDPattern.MatchString(entry.PublisherID) {
			return fmt.Errorf("adapter registry entry has an invalid ID, digest, or publisher")
		}
		if err := validateRegistryRelativePath(entry.ManifestPath); err != nil {
			return err
		}
		if entry.Trust != "first-party-sandboxed" && entry.Trust != "verified-community" &&
			entry.Trust != "unverified-community" && entry.Trust != "local" {
			return fmt.Errorf("adapter registry entry %s has unsupported trust %q", entry.AdapterID, entry.Trust)
		}
		if entry.Status != "active" && entry.Status != "deprecated" && entry.Status != "revoked" {
			return fmt.Errorf("adapter registry entry %s has unsupported status %q", entry.AdapterID, entry.Status)
		}
		if (entry.Status == "active" && entry.Reason != "") ||
			(entry.Status != "active" && strings.TrimSpace(entry.Reason) == "") {
			return fmt.Errorf("adapter registry entry %s has an inconsistent lifecycle reason", entry.AdapterID)
		}
	}
	return nil
}

func validateRegistryRelativePath(path string) error {
	if path == "" || strings.Contains(path, "\\") || filepath.IsAbs(path) || filepath.ToSlash(filepath.Clean(path)) != path ||
		path == "." || path == ".." || strings.HasPrefix(path, "../") ||
		(!strings.HasSuffix(path, ".yaml") && !strings.HasSuffix(path, ".yml")) {
		return fmt.Errorf("adapter registry manifest path %q must be a normalized relative slash path", path)
	}
	return nil
}

func resolveRegistryManifest(base, relative string) (string, error) {
	if err := validateRegistryRelativePath(relative); err != nil {
		return "", err
	}
	path := base
	for _, component := range strings.Split(relative, "/") {
		path = filepath.Join(path, component)
		info, err := os.Lstat(path)
		if err != nil {
			return "", err
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return "", fmt.Errorf("registered manifest path cannot contain symlinks")
		}
	}
	info, err := os.Lstat(path)
	if err != nil {
		return "", err
	}
	if info.Mode()&os.ModeSymlink != 0 || !info.Mode().IsRegular() {
		return "", fmt.Errorf("registered manifest must be a non-symlink regular file")
	}
	resolvedBase, err := filepath.EvalSymlinks(base)
	if err != nil {
		return "", err
	}
	resolvedPath, err := filepath.EvalSymlinks(path)
	if err != nil {
		return "", err
	}
	relativeToBase, err := filepath.Rel(resolvedBase, resolvedPath)
	if err != nil || relativeToBase == ".." || strings.HasPrefix(relativeToBase, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("registered manifest escapes the registry directory")
	}
	return resolvedPath, nil
}

func stringMember(values []string, wanted string) bool {
	for _, value := range values {
		if value == wanted {
			return true
		}
	}
	return false
}
