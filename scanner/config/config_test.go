package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func fixturePath(t *testing.T) string {
	t.Helper()
	path, err := filepath.Abs(filepath.Join("..", "..", "fixtures", "config", "production-readiness.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	return path
}

func writeConfig(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "production-readiness.yaml")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func validContent(t *testing.T) string {
	t.Helper()
	data, err := os.ReadFile(fixturePath(t))
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func TestLoadValidatesAndCanonicalizesConfiguration(t *testing.T) {
	first, err := Load(fixturePath(t))
	if err != nil {
		t.Fatal(err)
	}
	if first.SchemaVersion != ValidationSchema || len(first.Digest) != 64 || first.Configuration.Project.ID != "example-product" {
		t.Fatalf("unexpected validation: %+v", first)
	}
	reorderedFeatures := strings.Replace(
		validContent(t),
		"  authentication: true\n  payments: false",
		"  payments: false\n  authentication: true",
		1,
	)
	second, err := Load(writeConfig(t, reorderedFeatures))
	if err != nil {
		t.Fatal(err)
	}
	if first.Digest != second.Digest {
		t.Fatalf("canonical digest depends on mapping order: %s != %s", first.Digest, second.Digest)
	}
}

func TestValidationRejectsMutatedCanonicalIdentity(t *testing.T) {
	validation, err := Load(fixturePath(t))
	if err != nil {
		t.Fatal(err)
	}
	validation.Configuration.Project.Name = "Changed after validation"
	if err := validation.Validate(); err == nil || !strings.Contains(err.Error(), "digest does not match") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidationRejectsChangedSourceBytes(t *testing.T) {
	path := writeConfig(t, validContent(t))
	validation, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("# changed\n"+validContent(t)), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := validation.VerifySource(path); err == nil || !strings.Contains(err.Error(), "changed") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLoadRejectsUnknownDuplicateAndTrailingDocuments(t *testing.T) {
	tests := []struct {
		name    string
		content string
	}{
		{name: "unknown field", content: validContent(t) + "unexpected: true\n"},
		{name: "duplicate field", content: strings.Replace(validContent(t), "  id: example-product", "  id: example-product\n  id: duplicate", 1)},
		{name: "trailing document", content: validContent(t) + "---\nextra: true\n"},
		{name: "missing required field", content: strings.Replace(validContent(t), "features:\n  authentication: true\n  payments: false\n", "", 1)},
		{name: "null required value", content: strings.Replace(validContent(t), `  source_ref: ""`, "  source_ref: null", 1)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if _, err := Load(writeConfig(t, test.content)); err == nil {
				t.Fatal("invalid configuration was accepted")
			}
		})
	}
}

func TestLoadRejectsCapabilityExpansionAndUnsafePaths(t *testing.T) {
	tests := []string{
		strings.Replace(validContent(t), "  network: deny", "  network: allow", 1),
		strings.Replace(validContent(t), "  production_connected: false", "  production_connected: true", 1),
		strings.Replace(validContent(t), "  allow_commands: []", "  allow_commands: [[sh]]", 1),
		strings.Replace(validContent(t), "  protected_paths: [.git/", "  protected_paths: [../escape, .git/", 1),
		strings.Replace(validContent(t), `  source_ref: ""`, `  source_ref: "${GIT_SHA}"`, 1),
	}
	for index, content := range tests {
		t.Run(string(rune('a'+index)), func(t *testing.T) {
			if _, err := Load(writeConfig(t, content)); err == nil {
				t.Fatal("unsafe configuration was accepted")
			}
		})
	}
}

func TestLoadRejectsUnsortedSecurityRelevantLists(t *testing.T) {
	content := strings.Replace(
		validContent(t),
		"  classifications: [personal, public]",
		"  classifications: [public, personal]",
		1,
	)
	if _, err := Load(writeConfig(t, content)); err == nil || !strings.Contains(err.Error(), "must be sorted") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLoadRejectsSymlink(t *testing.T) {
	link := filepath.Join(t.TempDir(), "production-readiness.yaml")
	if err := os.Symlink(fixturePath(t), link); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(link); err == nil || !strings.Contains(err.Error(), "regular file") {
		t.Fatalf("unexpected error: %v", err)
	}
}
