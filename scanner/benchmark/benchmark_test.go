package benchmark

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/MarinJursic/production-readiness-checklist/scanner/catalog"
)

func repositoryRoot(t *testing.T) string {
	t.Helper()
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	return root
}

func TestCoreNativeBenchmarkMeetsQualityBudget(t *testing.T) {
	root := repositoryRoot(t)
	catalogValue, err := catalog.Load(root)
	if err != nil {
		t.Fatal(err)
	}
	report, err := Evaluate(catalogValue, filepath.Join(root, "fixtures", "benchmarks", "core-native", "suite.yaml"),
		time.Date(2026, 8, 23, 12, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}
	if !report.Passed || len(report.CorpusDigest) != 64 || report.Summary.Cases != 6 || report.Summary.Expectations != 10 ||
		report.Summary.Mismatched != 0 || report.Summary.DeterministicCases != report.Summary.Cases ||
		report.Metrics.Precision != 1 || report.Metrics.Recall != 1 || report.Metrics.FalsePositiveRate != 0 {
		t.Fatalf("benchmark report = %+v", report)
	}
	if report.Summary.ExpectedOutcomes.Fail == 0 || report.Summary.ExpectedOutcomes.NotApplicable == 0 ||
		report.Summary.ExpectedOutcomes.Unknown == 0 || report.Summary.ExpectedOutcomes.ManualReview == 0 {
		t.Fatalf("benchmark outcome coverage = %+v", report.Summary.ExpectedOutcomes)
	}
}

func TestComprehensiveCoreBenchmarkCoversEveryCatalogAssertion(t *testing.T) {
	root := repositoryRoot(t)
	catalogValue, err := catalog.Load(root)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(root, "fixtures", "benchmarks", "core-native", "suite-comprehensive.yaml")
	loaded, err := Load(path, catalogValue)
	if err != nil {
		t.Fatal(err)
	}
	covered := map[string]bool{}
	for _, benchmarkCase := range loaded.Suite.Cases {
		for _, expectation := range benchmarkCase.Expectations {
			covered[expectation.AssertionID] = true
		}
	}
	profile := catalogValue.Profiles["prc/core-repository"]
	if len(covered) != len(profile.AssertionIDs) {
		t.Fatalf("covered assertions = %d, profile assertions = %d", len(covered), len(profile.AssertionIDs))
	}
	for _, assertionID := range profile.AssertionIDs {
		if !covered[assertionID] {
			t.Errorf("profile assertion is not covered: %s", assertionID)
		}
	}
	report, err := Evaluate(catalogValue, path, time.Date(2026, 8, 23, 12, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}
	if !report.Passed || report.Summary.Cases != 21 || report.Summary.Expectations != 103 ||
		report.Summary.Mismatched != 0 || report.Summary.DeterministicCases != report.Summary.Cases ||
		report.Metrics.Precision != 1 || report.Metrics.Recall != 1 || report.Metrics.FalsePositiveRate != 0 {
		t.Fatalf("comprehensive benchmark report = %+v", report)
	}
}

func TestBenchmarkMismatchFailsBudgetAndTraversalIsRejected(t *testing.T) {
	root := repositoryRoot(t)
	catalogValue, err := catalog.Load(root)
	if err != nil {
		t.Fatal(err)
	}
	directory := t.TempDir()
	target := filepath.Join(directory, "target")
	if err := os.Mkdir(target, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(target, "README.md"), []byte("present\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	suite := `schema_version: prc.benchmark-suite/v0.1
id: prc.benchmark.mismatch@0.1
title: Deliberate mismatch
profile_id: prc/core-repository
cases:
  - id: mismatch
    target: target
    expectations:
      - assertion_id: PRC-A-CORE-001
        assessment: fail
        execution: completed
quality_budget:
  minimum_precision: 1
  minimum_recall: 1
  maximum_false_positive_rate: 0
  maximum_mismatches: 0
  require_determinism: true
`
	path := filepath.Join(directory, "suite.yaml")
	if err := os.WriteFile(path, []byte(suite), 0o600); err != nil {
		t.Fatal(err)
	}
	report, err := Evaluate(catalogValue, path, time.Date(2026, 8, 23, 12, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}
	if report.Passed || report.Summary.Mismatched != 1 || report.Metrics.Recall != 0 || len(report.QualityFailures) == 0 {
		t.Fatalf("mismatch report = %+v", report)
	}

	traversal := strings.Replace(suite, "target: target", "target: ../target", 1)
	if err := os.WriteFile(path, []byte(traversal), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(path, catalogValue); err == nil || !strings.Contains(err.Error(), "normalized relative") {
		t.Fatalf("unexpected traversal error: %v", err)
	}

	missingBudgetField := strings.Replace(suite, "  require_determinism: true\n", "", 1)
	if err := os.WriteFile(path, []byte(missingBudgetField), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(path, catalogValue); err == nil || !strings.Contains(err.Error(), "requires field") {
		t.Fatalf("unexpected incomplete budget error: %v", err)
	}

	legacySetup := strings.Replace(suite, "    target: target\n", "    target: target\n    setup:\n      - operation: truncate\n        path: README.md\n", 1)
	if err := os.WriteFile(path, []byte(legacySetup), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(path, catalogValue); err == nil || !strings.Contains(err.Error(), "setup requires schema") {
		t.Fatalf("unexpected legacy setup error: %v", err)
	}
}

func TestFixtureSetupIsBoundedAndDoesNotMutateSource(t *testing.T) {
	source := t.TempDir()
	path := filepath.Join(source, "app.py")
	if err := os.WriteFile(path, []byte("print('ready')\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	mode := uint32(0o666)
	materialized, cleanup, err := materializeTarget(source, []SetupOperation{
		{Operation: "remove_final_newline", Path: "app.py"},
		{Operation: "chmod", Path: "app.py", Mode: &mode},
		{Operation: "git_head", Value: strings.Repeat("a", 40)},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()
	original, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	changed, err := os.ReadFile(filepath.Join(materialized, "app.py"))
	if err != nil {
		t.Fatal(err)
	}
	head, err := os.ReadFile(filepath.Join(materialized, ".git", "HEAD"))
	if err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(filepath.Join(materialized, "app.py"))
	if err != nil {
		t.Fatal(err)
	}
	if string(original) != "print('ready')\n" || string(changed) != "print('ready')" ||
		string(head) != strings.Repeat("a", 40)+"\n" || info.Mode().Perm() != 0o666 {
		t.Fatalf("materialized fixture original=%q changed=%q head=%q mode=%#o", original, changed, head, info.Mode().Perm())
	}

	invalidMode := uint32(0o1000)
	if err := validateSetupOperation(SetupOperation{Operation: "chmod", Path: "app.py", Mode: &invalidMode}); err == nil {
		t.Fatal("out-of-range fixture mode was accepted")
	}
	if err := validateSetupOperation(SetupOperation{Operation: "truncate", Path: "../outside"}); err == nil {
		t.Fatal("fixture path traversal was accepted")
	}
	if _, _, err := materializeTarget(source, []SetupOperation{{
		Operation: "replace_text", Path: "app.py", Find: "missing", Replace: "value",
	}}); err == nil || !strings.Contains(err.Error(), "exactly once") {
		t.Fatalf("missing replacement token error = %v", err)
	}
}
