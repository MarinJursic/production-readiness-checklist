package catalog

import (
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

func TestLoadRejectsDuplicateProfileAssertion(t *testing.T) {
	root := t.TempDir()
	writeCatalogFixture(t, root, "objectives/core.yaml", `schema_version: prc.objectives/v0.1
catalog_version: 1.0.0
objectives:
  - id: TEST-AAAAAAAA
    revision: 1
    title: Test
    statement: Test objective.
    source: {path: README.md, line: 1}
    domains: [repository]
    automation_class: full
    assertion_ids: [PRC-A-TEST-001]
`)
	writeCatalogFixture(t, root, "assertions/core.yaml", `schema_version: prc.assertions/v0.1
catalog_version: 1.0.0
assertions:
  - id: PRC-A-TEST-001
    revision: 1
    control_ids: [TEST-AAAAAAAA]
    title: Test
    statement: Test assertion.
    applicability: "true"
    evidence_required: []
    implementation_id: prc.native.test@0.1
    severity: high
    gate: required
    remediation_class: R0
`)
	writeCatalogFixture(t, root, "profiles/test.yaml", `schema_version: prc.profile/v0.1
id: prc/test
version: "1"
title: Test
description: Test profile.
assertion_ids: [PRC-A-TEST-001, PRC-A-TEST-001]
terminal_policy:
  block_on: [high]
  allow_manual_remaining: true
`)
	_, err := Load(root)
	if err == nil || !strings.Contains(err.Error(), "duplicate profile prc/test assertion reference") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLoadRejectsMixedCatalogVersions(t *testing.T) {
	root := t.TempDir()
	writeCatalogFixture(t, root, "objectives/core.yaml", "schema_version: prc.objectives/v0.1\ncatalog_version: 1.0.0\nobjectives: []\n")
	writeCatalogFixture(t, root, "assertions/core.yaml", "schema_version: prc.assertions/v0.1\ncatalog_version: 2.0.0\nassertions: []\n")
	writeCatalogFixture(t, root, "profiles/test.yaml", "schema_version: prc.profile/v0.1\nid: prc/test\nversion: \"1\"\ntitle: Test\ndescription: Test.\nassertion_ids: [PRC-A-TEST-001]\nterminal_policy: {block_on: [high], allow_manual_remaining: true}\n")
	_, err := Load(root)
	if err == nil || !strings.Contains(err.Error(), "does not match") {
		t.Fatalf("unexpected error: %v", err)
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
