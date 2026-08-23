package remediation

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	projectconfig "github.com/MarinJursic/production-readiness-checklist/scanner/config"
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

func configuredRemediation(t *testing.T, target string, mutate func(string) string) *ProjectConfiguration {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(testCatalogRoot(t), "fixtures", "config", "production-readiness.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	content := string(data)
	if mutate != nil {
		content = mutate(content)
	}
	path := filepath.Join(target, "production-readiness.yaml")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	validation, err := projectconfig.Load(path)
	if err != nil {
		t.Fatal(err)
	}
	return &ProjectConfiguration{Validation: validation, SourcePath: path}
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

func TestRunEnforcesAndRecordsConfiguredRemediationPolicy(t *testing.T) {
	target := remediationTarget(t)
	configuration := configuredRemediation(t, target, nil)
	candidate, err := Run(Options{
		CatalogRoot: testCatalogRoot(t), Target: target, CandidateDir: filepath.Join(t.TempDir(), "candidate"),
		ProfileID: "prc/core-repository", AssertionID: finalNewlineAssertion,
		MaxFiles: 20, MaxChangedLines: 200, Attempt: 1, MaxAttempts: 3, Configuration: configuration,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !candidate.Accepted || candidate.SchemaVersion != CandidateSchema ||
		candidate.Contract.ConfigurationDigest != configuration.Validation.Digest ||
		candidate.Contract.ProjectID != "example-product" || candidate.Contract.Attempt != 1 ||
		candidate.Contract.MaxAttempts != 3 || !slices.Contains(candidate.Contract.ProtectedPaths, "production-readiness.yaml") {
		t.Fatalf("configured candidate lost policy identity: %+v", candidate)
	}
	source, err := os.ReadFile(configuration.SourcePath)
	if err != nil {
		t.Fatal(err)
	}
	copied, err := os.ReadFile(filepath.Join(candidate.CandidatePath, "production-readiness.yaml"))
	if err != nil || string(copied) != string(source) {
		t.Fatalf("candidate configuration changed: %v", err)
	}
}

func TestRunRejectsConfiguredBudgetExpansionAndDisabledPolicy(t *testing.T) {
	target := remediationTarget(t)
	configuration := configuredRemediation(t, target, nil)
	_, err := Run(Options{
		CatalogRoot: testCatalogRoot(t), Target: target, CandidateDir: filepath.Join(t.TempDir(), "candidate"),
		ProfileID: "prc/core-repository", AssertionID: finalNewlineAssertion,
		MaxFiles: 21, MaxChangedLines: 200, MaxAttempts: 3, Configuration: configuration,
	})
	if err == nil || !strings.Contains(err.Error(), "budget exceeds") {
		t.Fatalf("unexpected budget error: %v", err)
	}

	disabledTarget := remediationTarget(t)
	disabled := configuredRemediation(t, disabledTarget, func(content string) string {
		return strings.Replace(content, "  enabled: true", "  enabled: false", 1)
	})
	_, err = Run(Options{
		CatalogRoot: testCatalogRoot(t), Target: disabledTarget, CandidateDir: filepath.Join(t.TempDir(), "candidate"),
		ProfileID: "prc/core-repository", AssertionID: finalNewlineAssertion,
		MaxFiles: 20, MaxChangedLines: 200, MaxAttempts: 3, Configuration: disabled,
	})
	if err == nil || !strings.Contains(err.Error(), "disabled") {
		t.Fatalf("unexpected disabled-policy error: %v", err)
	}
}

func TestRunAddsConfiguredProtectedPaths(t *testing.T) {
	target := remediationTarget(t)
	configuration := configuredRemediation(t, target, func(content string) string {
		return strings.Replace(content,
			"protected_paths: [.git/, .github/workflows/, .prc/, catalog/",
			"protected_paths: [.git/, .github/workflows/, .prc/, app.py, catalog/", 1)
	})
	_, err := Run(Options{
		CatalogRoot: testCatalogRoot(t), Target: target, CandidateDir: filepath.Join(t.TempDir(), "candidate"),
		ProfileID: "prc/core-repository", AssertionID: finalNewlineAssertion,
		MaxFiles: 20, MaxChangedLines: 200, MaxAttempts: 3, Configuration: configuration,
	})
	if err == nil || !strings.Contains(err.Error(), "protected path app.py") {
		t.Fatalf("unexpected protected-path error: %v", err)
	}
}

func TestRunCreatesBytePreservingRestrictiveModeCandidate(t *testing.T) {
	target := remediationTarget(t)
	originalPath := filepath.Join(target, "app.py")
	if err := os.Chmod(originalPath, 0o666); err != nil {
		t.Fatal(err)
	}
	original, err := os.ReadFile(originalPath)
	if err != nil {
		t.Fatal(err)
	}
	candidate, err := Run(Options{
		CatalogRoot: testCatalogRoot(t), Target: target,
		CandidateDir: filepath.Join(t.TempDir(), "candidate"), ProfileID: "prc/core-repository",
		AssertionID: restrictiveModesAssertion, MaxFiles: 20, MaxChangedLines: 1,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !candidate.Accepted || candidate.Contract.FixerID != restrictiveModesFixer {
		t.Fatalf("candidate rejected: %+v", candidate)
	}
	originalInfo, err := os.Stat(originalPath)
	if err != nil {
		t.Fatal(err)
	}
	candidatePath := filepath.Join(candidate.CandidatePath, "app.py")
	candidateInfo, err := os.Stat(candidatePath)
	if err != nil {
		t.Fatal(err)
	}
	fixed, err := os.ReadFile(candidatePath)
	if err != nil {
		t.Fatal(err)
	}
	if originalInfo.Mode().Perm() != 0o666 {
		t.Fatalf("original mode changed to %#o", originalInfo.Mode().Perm())
	}
	if candidateInfo.Mode().Perm() != 0o644 {
		t.Fatalf("candidate mode = %#o, want 0644", candidateInfo.Mode().Perm())
	}
	if string(fixed) != string(original) {
		t.Fatal("mode fixer changed file bytes")
	}
	if len(candidate.Changes) != 1 || candidate.Changes[0].Path != "app.py" ||
		candidate.Changes[0].BeforeMode != 0o666 || candidate.Changes[0].AfterMode != 0o644 ||
		candidate.Changes[0].AddedLines != 0 || candidate.Changes[0].RemovedLines != 0 {
		t.Fatalf("unexpected mode changes: %+v", candidate.Changes)
	}
}

func TestRestrictiveModeFixRejectsProtectedPath(t *testing.T) {
	target := remediationTarget(t)
	workflow := filepath.Join(target, ".github", "workflows", "ci.yml")
	if err := os.Chmod(workflow, 0o666); err != nil {
		t.Fatal(err)
	}
	_, err := Run(Options{
		CatalogRoot: testCatalogRoot(t), Target: target,
		CandidateDir: filepath.Join(t.TempDir(), "candidate"), ProfileID: "prc/core-repository",
		AssertionID: restrictiveModesAssertion, MaxFiles: 20, MaxChangedLines: 1,
	})
	if err == nil || !strings.Contains(err.Error(), "protected path .github/workflows/ci.yml") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRestrictiveModeAuditRejectsContentChange(t *testing.T) {
	target := remediationTarget(t)
	if err := os.Chmod(filepath.Join(target, "app.py"), 0o666); err != nil {
		t.Fatal(err)
	}
	candidate, err := Run(Options{
		CatalogRoot: testCatalogRoot(t), Target: target,
		CandidateDir: filepath.Join(t.TempDir(), "candidate"), ProfileID: "prc/core-repository",
		AssertionID: restrictiveModesAssertion, MaxFiles: 20, MaxChangedLines: 1,
	})
	if err != nil {
		t.Fatal(err)
	}
	writeTestFile(t, candidate.CandidatePath, "app.py", "tampered\n")
	baseline, err := inventory.Build(target)
	if err != nil {
		t.Fatal(err)
	}
	after, err := inventory.Build(candidate.CandidatePath)
	if err != nil {
		t.Fatal(err)
	}
	_, reasons := auditCandidate(baseline, after, candidate.Contract)
	if joined := strings.Join(reasons, " "); !strings.Contains(joined, "changed file content while restricting mode") {
		t.Fatalf("audit did not catch content change: %v", reasons)
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
