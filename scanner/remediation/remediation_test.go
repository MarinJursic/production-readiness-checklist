package remediation

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/MarinJursic/production-readiness-checklist/scanner/inventory"
)

const testPinnedCheckout = "actions/checkout@3d3c42e5aac5ba805825da76410c181273ba90b1"

func testCatalogRoot(t *testing.T) string {
	t.Helper()
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	return root
}

func writeTestFile(t *testing.T, root, relative, content string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(relative))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func remediationTarget(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	files := map[string]string{
		"README.md":                "# Example\n",
		"LICENSE":                  "MIT\n",
		"CONTRIBUTING.md":          "# Contributing\n",
		"SECURITY.md":              "# Security\n",
		"CODE_OF_CONDUCT.md":       "# Conduct\n",
		".github/CODEOWNERS":       "* @owner\n",
		"requirements.txt":         "example==1.0\n",
		"requirements.lock.txt":    "example==1.0\n",
		"app.py":                   "def ready(): return True",
		"tests/test_app.py":        "def test_ready(): assert True\n",
		".github/workflows/ci.yml": "name: CI\non: [push]\npermissions:\n  contents: read\njobs:\n  test:\n    runs-on: ubuntu-latest\n    steps:\n      - uses: " + testPinnedCheckout + "\n",
	}
	for relative, content := range files {
		writeTestFile(t, root, relative, content)
	}
	return root
}

func runAcceptedCandidate(t *testing.T) (string, Candidate) {
	t.Helper()
	target := remediationTarget(t)
	candidatePath := filepath.Join(t.TempDir(), "candidate")
	candidate, err := Run(Options{
		CatalogRoot: testCatalogRoot(t), Target: target, CandidateDir: candidatePath,
		ProfileID: "prc/core-repository", AssertionID: finalNewlineAssertion,
		MaxFiles: 20, MaxChangedLines: 20,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !candidate.Accepted {
		t.Fatalf("candidate rejected: %v", candidate.Reasons)
	}
	return target, candidate
}

func TestRunCreatesAcceptedIsolatedFinalNewlineCandidate(t *testing.T) {
	target, candidate := runAcceptedCandidate(t)
	original, err := os.ReadFile(filepath.Join(target, "app.py"))
	if err != nil {
		t.Fatal(err)
	}
	fixed, err := os.ReadFile(filepath.Join(candidate.CandidatePath, "app.py"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.HasSuffix(string(original), "\n") {
		t.Fatal("original workspace was modified")
	}
	if string(fixed) != string(original)+"\n" {
		t.Fatalf("unexpected fix %q", fixed)
	}
	if len(candidate.Changes) != 1 || candidate.Changes[0].Path != "app.py" || candidate.Changes[0].AddedLines != 1 {
		t.Fatalf("unexpected changes: %+v", candidate.Changes)
	}
	if candidate.Changes[0].BeforeMode != candidate.Changes[0].AfterMode {
		t.Fatal("fix changed the file mode")
	}
	if candidate.BeforeAssessment != "fail" || candidate.AfterAssessment != "pass" || candidate.Reasons == nil {
		t.Fatalf("unexpected assessments: %+v", candidate)
	}
	contractIdentity := candidate.Contract
	contractIdentity.TaskID = ""
	wantTaskID, err := contentID(contractIdentity)
	if err != nil || wantTaskID != candidate.Contract.TaskID {
		t.Fatalf("task identity mismatch: %v %s", err, wantTaskID)
	}
	wantCandidateID, err := candidateContentID(candidate)
	if err != nil || wantCandidateID != candidate.CandidateID {
		t.Fatalf("candidate identity mismatch: %v %s", err, wantCandidateID)
	}
}

func TestAuditRejectsProtectedAndExcludedPathChanges(t *testing.T) {
	target, candidate := runAcceptedCandidate(t)
	writeTestFile(t, candidate.CandidatePath, ".github/workflows/ci.yml", "tampered\n")
	writeTestFile(t, candidate.CandidatePath, ".prc/hidden", "unexpected\n")
	baseline, err := inventory.Build(target)
	if err != nil {
		t.Fatal(err)
	}
	after, err := inventory.Build(candidate.CandidatePath)
	if err != nil {
		t.Fatal(err)
	}
	_, reasons := auditCandidate(baseline, after, candidate.Contract)
	joined := strings.Join(reasons, " ")
	if !strings.Contains(joined, "protected path .github/workflows/ci.yml") ||
		!strings.Contains(joined, "unexpected path .prc/hidden") {
		t.Fatalf("audit did not catch attacks: %v", reasons)
	}
}

func TestAuditRejectsModeDeletionAndUnexpectedFile(t *testing.T) {
	target, candidate := runAcceptedCandidate(t)
	if err := os.Chmod(filepath.Join(candidate.CandidatePath, "app.py"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(filepath.Join(candidate.CandidatePath, "LICENSE")); err != nil {
		t.Fatal(err)
	}
	writeTestFile(t, candidate.CandidatePath, "extra.txt", "unexpected\n")
	baseline, err := inventory.Build(target)
	if err != nil {
		t.Fatal(err)
	}
	after, err := inventory.Build(candidate.CandidatePath)
	if err != nil {
		t.Fatal(err)
	}
	_, reasons := auditCandidate(baseline, after, candidate.Contract)
	joined := strings.Join(reasons, " ")
	for _, expected := range []string{"changed file mode for app.py", "deleted LICENSE", "added unexpected path extra.txt"} {
		if !strings.Contains(joined, expected) {
			t.Errorf("missing %q in %v", expected, reasons)
		}
	}
}

func TestCandidateMustBeOutsideTarget(t *testing.T) {
	target := remediationTarget(t)
	_, err := Run(Options{
		CatalogRoot: testCatalogRoot(t), Target: target,
		CandidateDir: filepath.Join(target, "candidate"), ProfileID: "prc/core-repository",
		AssertionID: finalNewlineAssertion, MaxFiles: 20, MaxChangedLines: 20,
	})
	if err == nil || !strings.Contains(err.Error(), "outside the target tree") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSafeJoinRejectsTraversalAndBackslashes(t *testing.T) {
	for _, path := range []string{"../escape", "/absolute", "a\\b", "a/../b", ""} {
		if _, err := safeJoin(t.TempDir(), path); err == nil {
			t.Errorf("accepted unsafe path %q", path)
		}
	}
}
