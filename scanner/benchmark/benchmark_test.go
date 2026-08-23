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
}
