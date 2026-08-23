package pack

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/MarinJursic/production-readiness-checklist/scanner/catalog"
	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
)

func repositoryRoot(t *testing.T) string {
	t.Helper()
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	return root
}

func TestCoreFoundationPackBindsCatalogAndBenchmarkCoverage(t *testing.T) {
	root := repositoryRoot(t)
	catalogValue, err := catalog.Load(root)
	if err != nil {
		t.Fatal(err)
	}
	loaded, err := Load(root, filepath.Join(root, "packs", "core-foundation.yaml"), catalogValue)
	if err != nil {
		t.Fatal(err)
	}
	if len(loaded.Digest) != 64 || loaded.SuiteDigest != loaded.Manifest.Benchmark.SuiteSHA256 ||
		len(loaded.BenchmarkCorpusDigest) != 64 ||
		len(loaded.Manifest.Assertions) != 3 || len(loaded.CatalogDigest) != 64 {
		t.Fatalf("loaded pack = %+v", loaded)
	}
}

func TestPackRejectsUnpinnedBenchmarkAndCatalogDrift(t *testing.T) {
	root := repositoryRoot(t)
	catalogValue, err := catalog.Load(root)
	if err != nil {
		t.Fatal(err)
	}
	directory := t.TempDir()
	path := filepath.Join(directory, "pack.yaml")
	if err := os.WriteFile(path, []byte("outside repository root\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(root, path, catalogValue); err == nil || !strings.Contains(err.Error(), "escapes") {
		t.Fatalf("manifest outside root error = %v", err)
	}

	path = filepath.Join(root, "packs", "core-foundation.yaml")
	wrongCatalog := *catalogValue
	wrongCatalog.Assertions = make(map[string]model.Assertion, len(catalogValue.Assertions))
	for identifier, assertion := range catalogValue.Assertions {
		wrongCatalog.Assertions[identifier] = assertion
	}
	assertion := wrongCatalog.Assertions["PRC-A-CORE-001"]
	assertion.ImplementationID = "prc.native.changed@9.9"
	wrongCatalog.Assertions[assertion.ID] = assertion
	if _, err := Load(root, path, &wrongCatalog); err == nil || !strings.Contains(err.Error(), "implementation") {
		t.Fatalf("catalog drift error = %v", err)
	}

	wrongBehavior := *catalogValue
	wrongBehavior.Assertions = make(map[string]model.Assertion, len(catalogValue.Assertions))
	for identifier, assertion := range catalogValue.Assertions {
		wrongBehavior.Assertions[identifier] = assertion
	}
	assertion = wrongBehavior.Assertions["PRC-A-CORE-001"]
	assertion.Parameters = map[string]any{"paths": []string{"MISSING.md"}, "minimum_bytes": 1}
	wrongBehavior.Assertions[assertion.ID] = assertion
	if _, err := Load(root, path, &wrongBehavior); err == nil || !strings.Contains(err.Error(), "quality budget failed") {
		t.Fatalf("benchmark drift error = %v", err)
	}
}
