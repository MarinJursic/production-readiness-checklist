package catalog

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDigestBindsGoverningDefinitions(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	loaded, err := Load(root)
	if err != nil {
		t.Fatal(err)
	}
	first, err := loaded.Digest()
	if err != nil {
		t.Fatal(err)
	}
	second, err := loaded.Digest()
	if err != nil {
		t.Fatal(err)
	}
	if first != second || len(first) != 64 {
		t.Fatalf("catalog digest is not deterministic: %q != %q", first, second)
	}
	for id, objective := range loaded.Objectives {
		objective.Statement += " changed"
		loaded.Objectives[id] = objective
		break
	}
	changed, err := loaded.Digest()
	if err != nil {
		t.Fatal(err)
	}
	if changed == first {
		t.Fatal("objective change did not change catalog identity")
	}
}

func TestBundleIsDeterministicCompleteAndPathIndependent(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	loaded, err := Load(root)
	if err != nil {
		t.Fatal(err)
	}
	first, err := loaded.Bundle()
	if err != nil {
		t.Fatal(err)
	}
	second, err := loaded.Bundle()
	if err != nil {
		t.Fatal(err)
	}
	firstJSON, err := json.Marshal(first)
	if err != nil {
		t.Fatal(err)
	}
	secondJSON, err := json.Marshal(second)
	if err != nil {
		t.Fatal(err)
	}
	if string(firstJSON) != string(secondJSON) {
		t.Fatal("catalog bundle output is not deterministic")
	}
	if first.SchemaVersion != BundleSchema || first.Manifest.SchemaVersion != ManifestSchema ||
		first.Manifest.CatalogVersion != loaded.Version || len(first.Manifest.CatalogDigest) != 64 ||
		first.Manifest.ObjectiveCount != len(first.Objectives) ||
		first.Manifest.AssertionCount != len(first.Assertions) ||
		first.Manifest.ProfileCount != len(first.Profiles) {
		t.Fatalf("incomplete catalog bundle: %+v", first.Manifest)
	}
	for index := 1; index < len(first.Objectives); index++ {
		if first.Objectives[index-1].ID >= first.Objectives[index].ID {
			t.Fatal("bundle objectives are not ordered by ID")
		}
	}
	if strings.Contains(string(firstJSON), root) {
		t.Fatal("catalog bundle leaked its filesystem root")
	}
}

func TestLoadRejectsDuplicateProfileAssertion(t *testing.T) {
	root := t.TempDir()
	writeSourceFixture(t, root, "- [ ] **USEQ-AAAAAAAA** — Test objective.\n")
	writeCatalogFixture(t, root, "objectives/core.yaml", `schema_version: prc.objectives/v0.1
catalog_version: 1.0.0
objectives:
  - id: USEQ-AAAAAAAA
    revision: 1
    title: Test
    statement: Test objective.
    source: {path: docs/source.md, line: 1}
    domains: [repository]
    automation_class: automated
    assertion_ids: [PRC-A-TEST-001]
`)
	writeCatalogFixture(t, root, "assertions/core.yaml", `schema_version: prc.assertions/v0.1
catalog_version: 1.0.0
assertions:
  - id: PRC-A-TEST-001
    revision: 1
    control_ids: [USEQ-AAAAAAAA]
    title: Test
    statement: Test assertion.
    applicability: "true"
    evidence_required:
      - kind: repository-file
        minimum_authority: repository
        description: Test evidence.
    implementation_id: prc.native.test@0.1
    severity: high
    gate: required
    remediation_class: R0
`)
	writeCatalogFixture(t, root, "profiles/test.yaml", `schema_version: prc.profile/v0.1
id: prc/test
version: "1.0"
title: Test
description: Test profile.
assertion_ids: [PRC-A-TEST-001, PRC-A-TEST-001]
terminal_policy:
  block_on: [high]
  allow_manual_remaining: true
`)
	_, err := Load(root)
	if err == nil || !strings.Contains(err.Error(), "duplicate profile assertion ID") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLoadRejectsMixedCatalogVersions(t *testing.T) {
	root := t.TempDir()
	writeSourceFixture(t, root, "- [ ] **USEQ-AAAAAAAA** — Test objective.\n")
	writeCatalogFixture(t, root, "objectives/core.yaml", validObjectives("1.0.0", ""))
	writeCatalogFixture(t, root, "assertions/core.yaml", validAssertions("2.0.0", "high", ""))
	writeCatalogFixture(t, root, "profiles/test.yaml", validProfile("PRC-A-TEST-001"))
	_, err := Load(root)
	if err == nil || !strings.Contains(err.Error(), "does not match") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLoadRejectsUnknownYAMLFields(t *testing.T) {
	root := t.TempDir()
	writeSourceFixture(t, root, "- [ ] **USEQ-AAAAAAAA** — Test objective.\n")
	writeCatalogFixture(t, root, "objectives/core.yaml", validObjectives("1.0.0", "    unexpected: true\n"))
	writeCatalogFixture(t, root, "assertions/core.yaml", validAssertions("1.0.0", "high", ""))
	writeCatalogFixture(t, root, "profiles/test.yaml", validProfile("PRC-A-TEST-001"))
	_, err := Load(root)
	if err == nil || !strings.Contains(err.Error(), "field unexpected not found") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLoadRejectsInvalidAssertionSemantics(t *testing.T) {
	root := t.TempDir()
	writeSourceFixture(t, root, "- [ ] **USEQ-AAAAAAAA** — Test objective.\n")
	writeCatalogFixture(t, root, "objectives/core.yaml", validObjectives("1.0.0", ""))
	writeCatalogFixture(t, root, "assertions/core.yaml", validAssertions("1.0.0", "urgent", ""))
	writeCatalogFixture(t, root, "profiles/test.yaml", validProfile("PRC-A-TEST-001"))
	_, err := Load(root)
	if err == nil || !strings.Contains(err.Error(), "severity") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLoadRejectsStaleObjectiveSource(t *testing.T) {
	root := t.TempDir()
	writeSourceFixture(t, root, "- [ ] **USEQ-AAAAAAAA** — Different objective.\n")
	writeCatalogFixture(t, root, "objectives/core.yaml", validObjectives("1.0.0", ""))
	writeCatalogFixture(t, root, "assertions/core.yaml", validAssertions("1.0.0", "high", ""))
	writeCatalogFixture(t, root, "profiles/test.yaml", validProfile("PRC-A-TEST-001"))
	_, err := Load(root)
	if err == nil || !strings.Contains(err.Error(), "does not exactly match") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLoadRejectsBrokenReverseMapping(t *testing.T) {
	root := t.TempDir()
	writeSourceFixture(t, root, "- [ ] **USEQ-AAAAAAAA** — Test objective.\n")
	writeCatalogFixture(t, root, "objectives/core.yaml", validObjectives("1.0.0", ""))
	writeCatalogFixture(t, root, "assertions/core.yaml", strings.Replace(
		validAssertions("1.0.0", "high", ""),
		"control_ids: [USEQ-AAAAAAAA]",
		"control_ids: [USEQ-BBBBBBBB]",
		1,
	))
	writeCatalogFixture(t, root, "profiles/test.yaml", validProfile("PRC-A-TEST-001"))
	_, err := Load(root)
	if err == nil || !strings.Contains(err.Error(), "does not map back") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLoadRejectsSymlinkCatalogFile(t *testing.T) {
	root := t.TempDir()
	writeSourceFixture(t, root, "- [ ] **USEQ-AAAAAAAA** — Test objective.\n")
	outside := filepath.Join(t.TempDir(), "objectives.yaml")
	if err := os.WriteFile(outside, []byte(validObjectives("1.0.0", "")), 0o600); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(root, "catalog", "objectives", "core.yaml")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, path); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	writeCatalogFixture(t, root, "assertions/core.yaml", validAssertions("1.0.0", "high", ""))
	writeCatalogFixture(t, root, "profiles/test.yaml", validProfile("PRC-A-TEST-001"))
	_, err := Load(root)
	if err == nil || !strings.Contains(err.Error(), "must be a regular file") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func validObjectives(version, extra string) string {
	return `schema_version: prc.objectives/v0.1
catalog_version: ` + version + `
objectives:
  - id: USEQ-AAAAAAAA
    revision: 1
    title: Test
    statement: Test objective.
    source: {path: docs/source.md, line: 1}
    domains: [repository]
    automation_class: automated
    assertion_ids: [PRC-A-TEST-001]
` + extra
}

func validAssertions(version, severity, extra string) string {
	return `schema_version: prc.assertions/v0.1
catalog_version: ` + version + `
assertions:
  - id: PRC-A-TEST-001
    revision: 1
    control_ids: [USEQ-AAAAAAAA]
    title: Test
    statement: Test assertion.
    applicability: "true"
    evidence_required:
      - kind: repository-file
        minimum_authority: repository
        description: Test evidence.
    implementation_id: prc.native.test@0.1
    severity: ` + severity + `
    gate: required
    remediation_class: R0
` + extra
}

func validProfile(assertionID string) string {
	return `schema_version: prc.profile/v0.1
id: prc/test
version: "1.0"
title: Test
description: Test profile.
assertion_ids: [` + assertionID + `]
terminal_policy:
  block_on: [high]
  allow_manual_remaining: true
`
}

func writeSourceFixture(t *testing.T, root, content string) {
	t.Helper()
	path := filepath.Join(root, "docs", "source.md")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}

func writeCatalogFixture(t *testing.T, root, relative, content string) {
	t.Helper()
	path := filepath.Join(root, "catalog", filepath.FromSlash(relative))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}
