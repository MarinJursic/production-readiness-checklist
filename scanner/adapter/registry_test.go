package adapter

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func writeRegistryFixture(t *testing.T, manifest *Manifest, entry RegistryEntry) string {
	t.Helper()
	directory := t.TempDir()
	if manifest != nil {
		data, err := yaml.Marshal(manifest)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(directory, "adapter.yaml"), data, 0o600); err != nil {
			t.Fatal(err)
		}
		digest, err := ManifestDigest(*manifest)
		if err != nil {
			t.Fatal(err)
		}
		if entry.AdapterID == "" {
			entry.AdapterID = manifest.ID
		}
		if entry.ManifestSHA256 == "" {
			entry.ManifestSHA256 = digest
		}
		if entry.PublisherID == "" {
			entry.PublisherID = manifest.Publisher.ID
		}
	}
	if entry.ManifestPath == "" {
		entry.ManifestPath = "adapter.yaml"
	}
	if entry.Trust == "" {
		entry.Trust = "first-party-sandboxed"
	}
	if entry.Status == "" {
		entry.Status = "active"
	}
	document := RegistryDocument{
		SchemaVersion: RegistrySchema, ID: "prc.adapter-registry.test@0.1", Revision: 1,
		Entries: []RegistryEntry{entry},
	}
	data, err := yaml.Marshal(document)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(directory, "registry.yaml")
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestCheckedInRegistryPinsAndResolvesFixture(t *testing.T) {
	registry, err := LoadRegistry(fixturePath("fixture-registry.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	if len(registry.Digest) != 64 || len(registry.Entries) != 1 {
		t.Fatalf("registry identity = %+v", registry)
	}
	resolved, err := registry.Resolve(
		"prc.adapter.fixture@0.1", registry.Entries[0].ManifestSHA256,
		[]string{"fixture-result"}, DefaultRegistryPolicy(),
	)
	if err != nil {
		t.Fatal(err)
	}
	if resolved.Manifest.Publisher.ID != "prc-project" || resolved.Entry.Trust != "first-party-sandboxed" {
		t.Fatalf("resolved adapter = %+v", resolved)
	}
	if resolved.Resolution.Source != ResolutionSourceRegistry || resolved.Resolution.RegistryID != registry.ID ||
		resolved.Resolution.RegistryRevision != registry.Revision || resolved.Resolution.RegistryDigest != registry.Digest {
		t.Fatalf("registry resolution provenance = %+v", resolved.Resolution)
	}
}

func TestRegistryResolutionEnforcesTrustLifecycleDigestAndKinds(t *testing.T) {
	manifest := validManifest()
	path := writeRegistryFixture(t, &manifest, RegistryEntry{Trust: "unverified-community"})
	registry, err := LoadRegistry(path)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := registry.Resolve(manifest.ID, "", nil, DefaultRegistryPolicy()); err == nil || !strings.Contains(err.Error(), "trust level") {
		t.Fatalf("unexpected trust error: %v", err)
	}
	policy := RegistryPolicy{AllowedTrust: []string{"unverified-community"}}
	if _, err := registry.Resolve(manifest.ID, strings.Repeat("0", 64), nil, policy); err == nil || !strings.Contains(err.Error(), "digest") {
		t.Fatalf("unexpected digest error: %v", err)
	}
	if _, err := registry.Resolve(manifest.ID, "", []string{"undeclared"}, policy); err == nil || !strings.Contains(err.Error(), "observation kind") {
		t.Fatalf("unexpected observation error: %v", err)
	}

	manifest.Maintenance = "deprecated"
	path = writeRegistryFixture(t, &manifest, RegistryEntry{Status: "deprecated", Reason: "Superseded by a maintained adapter."})
	registry, err = LoadRegistry(path)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := registry.Resolve(manifest.ID, "", nil, DefaultRegistryPolicy()); err == nil || !strings.Contains(err.Error(), "deprecated") {
		t.Fatalf("unexpected deprecation error: %v", err)
	}
}

func TestRegistryRevocationDoesNotRequireCompromisedManifest(t *testing.T) {
	entry := RegistryEntry{
		AdapterID: "prc.adapter.compromised@1.0", ManifestPath: "removed.yaml",
		ManifestSHA256: strings.Repeat("a", 64), PublisherID: "compromised-publisher",
		Trust: "verified-community", Status: "revoked", Reason: "Compromised signing and release pipeline.",
	}
	registry, err := LoadRegistry(writeRegistryFixture(t, nil, entry))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := registry.Resolve(entry.AdapterID, entry.ManifestSHA256, nil, DefaultRegistryPolicy()); err == nil || !strings.Contains(err.Error(), "revoked") {
		t.Fatalf("unexpected revocation error: %v", err)
	}
}

func TestRegistryRejectsManifestTamperingPublisherMismatchAndTraversal(t *testing.T) {
	manifest := validManifest()
	path := writeRegistryFixture(t, &manifest, RegistryEntry{})
	manifestPath := filepath.Join(filepath.Dir(path), "adapter.yaml")
	if err := os.WriteFile(manifestPath, []byte("tampered\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadRegistry(path); err == nil {
		t.Fatal("tampered manifest was accepted")
	}

	path = writeRegistryFixture(t, &manifest, RegistryEntry{PublisherID: "other-publisher"})
	if _, err := LoadRegistry(path); err == nil || !strings.Contains(err.Error(), "publisher pin") {
		t.Fatalf("unexpected publisher error: %v", err)
	}

	path = writeRegistryFixture(t, nil, RegistryEntry{
		AdapterID: "prc.adapter.fixture@0.1", ManifestPath: "../adapter.yaml",
		ManifestSHA256: strings.Repeat("a", 64), PublisherID: "prc-project",
		Trust: "first-party-sandboxed", Status: "revoked", Reason: "Fixture revocation.",
	})
	if _, err := LoadRegistry(path); err == nil || !strings.Contains(err.Error(), "normalized relative") {
		t.Fatalf("unexpected traversal error: %v", err)
	}
}

func TestRegistryRejectsSymlinkedManifest(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink creation is not generally available")
	}
	manifest := validManifest()
	path := writeRegistryFixture(t, &manifest, RegistryEntry{})
	manifestPath := filepath.Join(filepath.Dir(path), "adapter.yaml")
	realPath := filepath.Join(filepath.Dir(path), "real.yaml")
	if err := os.Rename(manifestPath, realPath); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(realPath, manifestPath); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadRegistry(path); err == nil || !strings.Contains(err.Error(), "symlink") {
		t.Fatalf("unexpected symlink error: %v", err)
	}
}
