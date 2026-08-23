package adapterfixture

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func fixtureSuitePath(t *testing.T) string {
	t.Helper()
	path, err := filepath.Abs(filepath.Join("..", "..", "fixtures", "adapters", "fixture-suite.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	return path
}

func TestCheckedInAdapterFixtureSuitePasses(t *testing.T) {
	report, err := Evaluate(fixtureSuitePath(t))
	if err != nil {
		t.Fatal(err)
	}
	if !report.Passed || report.SchemaVersion != ReportSchema || report.Summary.Cases != 7 ||
		report.Summary.Matched != 7 || report.Summary.DeterministicCases != 7 ||
		len(report.SuiteDigest) != 64 || len(report.CorpusDigest) != 64 || len(report.ManifestSHA256) != 64 {
		t.Fatalf("unexpected adapter fixture report: %+v", report)
	}
	for _, item := range report.Cases {
		if !item.Passed || !item.Deterministic || len(item.OutputSHA256) != 64 {
			t.Errorf("unexpected case result: %+v", item)
		}
	}
}

func TestAdapterFixtureMismatchProducesFailedReport(t *testing.T) {
	directory, suite := copyCheckedInFixtures(t)
	suite = strings.Replace(suite, "summary_status: timeout", "summary_status: partial", 1)
	path := filepath.Join(directory, "fixture-suite.yaml")
	if err := os.WriteFile(path, []byte(suite), 0o600); err != nil {
		t.Fatal(err)
	}
	report, err := Evaluate(path)
	if err != nil {
		t.Fatal(err)
	}
	if report.Passed || report.Summary.Mismatched != 1 || len(report.QualityFailures) != 1 {
		t.Fatalf("unexpected mismatch report: %+v", report)
	}
}

func TestAdapterFixtureRejectsUntrustedInputs(t *testing.T) {
	tests := map[string]struct {
		mutate func(string) string
		want   string
	}{
		"manifest drift": {
			mutate: func(value string) string {
				return strings.Replace(value, "4e63176c5a0a32c36d47ef705b076d1d90c7d75905bbb5fd85b21c175ed3b7b2", strings.Repeat("a", 64), 1)
			},
			want: "manifest digest mismatch",
		},
		"path traversal": {
			mutate: func(value string) string {
				return strings.Replace(value, "output: valid-output.jsonl", "output: ../valid-output.jsonl", 1)
			},
			want: "normalized relative JSONL path",
		},
		"limit expansion": {
			mutate: func(value string) string { return strings.Replace(value, "max_messages: 2", "max_messages: 10001", 1) },
			want:   "cannot exceed manifest value",
		},
		"unknown field": {
			mutate: func(value string) string {
				return strings.Replace(value, "title: Adapter protocol", "unknown: true\ntitle: Adapter protocol", 1)
			},
			want: "field unknown not found",
		},
		"duplicate field": {
			mutate: func(value string) string {
				return strings.Replace(value, "title: Adapter protocol", "title: duplicate\ntitle: Adapter protocol", 1)
			},
			want: "mapping keys must be unique",
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			directory, suite := copyCheckedInFixtures(t)
			path := filepath.Join(directory, "fixture-suite.yaml")
			if err := os.WriteFile(path, []byte(test.mutate(suite)), 0o600); err != nil {
				t.Fatal(err)
			}
			if _, err := Load(path); err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestAdapterFixtureRejectsSymlinkedTranscript(t *testing.T) {
	directory, suite := copyCheckedInFixtures(t)
	path := filepath.Join(directory, "valid-output.jsonl")
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join("outside", "valid-output.jsonl"), path); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(directory, "fixture-suite.yaml"), []byte(suite), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(filepath.Join(directory, "fixture-suite.yaml")); err == nil || !strings.Contains(err.Error(), "cannot contain symlinks") {
		t.Fatalf("unexpected symlink error: %v", err)
	}
}

func copyCheckedInFixtures(t *testing.T) (string, string) {
	t.Helper()
	source := filepath.Dir(fixtureSuitePath(t))
	destination := t.TempDir()
	entries, err := os.ReadDir(source)
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		data, err := os.ReadFile(filepath.Join(source, entry.Name()))
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(destination, entry.Name()), data, 0o600); err != nil {
			t.Fatal(err)
		}
	}
	suite, err := os.ReadFile(filepath.Join(destination, "fixture-suite.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	return destination, string(suite)
}
