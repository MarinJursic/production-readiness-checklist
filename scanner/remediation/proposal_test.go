package remediation

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/MarinJursic/production-readiness-checklist/scanner/inventory"
	"github.com/MarinJursic/production-readiness-checklist/scanner/provider"
)

func proposalTarget(t *testing.T) string {
	t.Helper()
	target := remediationTarget(t)
	if err := os.Remove(filepath.Join(target, "tests", "test_app.py")); err != nil {
		t.Fatal(err)
	}
	return target
}

func sealedProposalTask(t *testing.T, target string, allowed, protectedPaths []string) provider.Task {
	t.Helper()
	draft := provider.Task{
		SchemaVersion: provider.TaskSchema, Mode: "suggest",
		WorkspaceInventoryDigest: strings.Repeat("b", 64), AssertionID: "PRC-A-CORE-010",
		ControlIDs: []string{"USEQ-12655775"}, Goal: "Add one focused, discoverable regression test.",
		ReadScope: "task-inputs", RelevantPaths: []string{"app.py"}, Inputs: []provider.InputFile{},
		AllowedPaths: allowed, ProtectedPaths: protectedPaths, AllowedCommands: [][]string{},
		Network: "deny", Secrets: "none", AllowRemoteSourceProcessing: true,
		TimeoutSeconds: 60, MaxOutputBytes: 256 * 1024,
	}
	data, err := json.Marshal(draft)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "draft.json")
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
	task, err := provider.SealTask(path, target)
	if err != nil {
		t.Fatal(err)
	}
	return task
}

func proposalOutput(task provider.Task, path, patch string) provider.Output {
	return provider.Output{
		SchemaVersion: provider.OutputSchema, TaskID: task.TaskID, Status: "candidate",
		RootCause: "No recognized automated test exists.", ChangedFiles: []string{path}, Patch: patch,
		CommandsRequestedOrRun: []provider.CommandResult{}, Limitations: []string{},
		RequestedCapabilityChanges: []string{},
	}
}

func defaultProposalProtectedPaths() []string {
	return []string{".git/", ".github/workflows/", ".prc/", "catalog/", "production-readiness.yaml", "schemas/"}
}

func TestRunProposalAppliesAndAcceptsIsolatedR2Candidate(t *testing.T) {
	target := proposalTarget(t)
	task := sealedProposalTask(t, target, []string{"app_test.py"}, defaultProposalProtectedPaths())
	patch := "diff --git a/app_test.py b/app_test.py\n" +
		"new file mode 100644\n--- /dev/null\n+++ b/app_test.py\n" +
		"@@ -0,0 +1,4 @@\n+from app import ready\n+\n+def test_ready():\n+    assert ready() is True\n"
	output := proposalOutput(task, "app_test.py", patch)
	candidate, err := RunProposal(ProposalOptions{
		CatalogRoot: testCatalogRoot(t), Target: target, CandidateDir: filepath.Join(t.TempDir(), "candidate"),
		ProfileID: "prc/core-repository", Provider: "codex", Task: task, Output: output,
		MaxFiles: 2, MaxChangedLines: 10,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !candidate.Accepted || candidate.Contract.RemediationClass != "R2" || candidate.Contract.ProposalTaskID != task.TaskID {
		t.Fatalf("unexpected candidate: %+v", candidate)
	}
	if _, err := os.Stat(filepath.Join(target, "app_test.py")); !os.IsNotExist(err) {
		t.Fatal("source workspace was modified")
	}
	content, err := os.ReadFile(filepath.Join(candidate.CandidatePath, "app_test.py"))
	if err != nil || string(content) != "from app import ready\n\ndef test_ready():\n    assert ready() is True\n" {
		t.Fatalf("unexpected candidate content %q: %v", content, err)
	}
	if len(candidate.Changes) != 1 || candidate.Changes[0].Kind != "added" || candidate.Changes[0].AddedLines != 4 {
		t.Fatalf("unexpected changes: %+v", candidate.Changes)
	}
}

func TestRunProposalRejectsMalformedHunkBeforeWriting(t *testing.T) {
	target := proposalTarget(t)
	task := sealedProposalTask(t, target, []string{"app_test.py"}, defaultProposalProtectedPaths())
	patch := "diff --git a/app_test.py b/app_test.py\n--- /dev/null\n+++ b/app_test.py\n" +
		"@@ -0,0 +1,2 @@\n+def test_ready():\n"
	candidatePath := filepath.Join(t.TempDir(), "candidate")
	_, err := RunProposal(ProposalOptions{
		CatalogRoot: testCatalogRoot(t), Target: target, CandidateDir: candidatePath,
		Provider: "codex", Task: task, Output: proposalOutput(task, "app_test.py", patch),
		MaxFiles: 2, MaxChangedLines: 10,
	})
	if err == nil || !strings.Contains(err.Error(), "hunk counts") {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, statErr := os.Stat(filepath.Join(candidatePath, "app_test.py")); !os.IsNotExist(statErr) {
		t.Fatal("malformed proposal wrote its target file")
	}
}

func TestRunProposalAddsDefaultProtectedPathsToWeakTask(t *testing.T) {
	target := proposalTarget(t)
	task := sealedProposalTask(t, target, []string{"catalog/attack.md"}, []string{".git/"})
	patch := "diff --git a/catalog/attack.md b/catalog/attack.md\n--- /dev/null\n+++ b/catalog/attack.md\n" +
		"@@ -0,0 +1 @@\n+attack\n"
	_, err := RunProposal(ProposalOptions{
		CatalogRoot: testCatalogRoot(t), Target: target, CandidateDir: filepath.Join(t.TempDir(), "candidate"),
		Provider: "codex", Task: task, Output: proposalOutput(task, "catalog/attack.md", patch),
		MaxFiles: 2, MaxChangedLines: 10,
	})
	if err == nil || !strings.Contains(err.Error(), "outside the R2 fix contract") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunProposalRejectsVacuousTestBeforeCreatingCandidate(t *testing.T) {
	target := proposalTarget(t)
	task := sealedProposalTask(t, target, []string{"app_test.py"}, defaultProposalProtectedPaths())
	patch := "diff --git a/app_test.py b/app_test.py\n--- /dev/null\n+++ b/app_test.py\n" +
		"@@ -0,0 +1,2 @@\n+def test_ready():\n+    assert True\n"
	candidatePath := filepath.Join(t.TempDir(), "candidate")
	_, err := RunProposal(ProposalOptions{
		CatalogRoot: testCatalogRoot(t), Target: target, CandidateDir: candidatePath,
		Provider: "codex", Task: task, Output: proposalOutput(task, "app_test.py", patch),
		MaxFiles: 2, MaxChangedLines: 10,
	})
	if err == nil || !strings.Contains(err.Error(), "constant assertion") {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, statErr := os.Stat(candidatePath); !os.IsNotExist(statErr) {
		t.Fatal("anti-gaming rejection created a candidate directory")
	}
}

func TestProposalAntiGamingRejectsExistingTestChangesAndSuppressions(t *testing.T) {
	target := remediationTarget(t)
	baseline, err := inventory.Build(target)
	if err != nil {
		t.Fatal(err)
	}
	output := provider.Output{
		ChangedFiles: []string{"tests/test_app.py"},
		Patch: "diff --git a/tests/test_app.py b/tests/test_app.py\n" +
			"--- a/tests/test_app.py\n+++ b/tests/test_app.py\n" +
			"@@ -1 +1,2 @@\n def test_ready(): assert True\n+# noqa\n",
	}
	reasons := strings.Join(auditProposalAntiGaming(baseline, output), " ")
	if !strings.Contains(reasons, "modifies existing test") || !strings.Contains(reasons, "suppression or skip") {
		t.Fatalf("anti-gaming audit missed proposal: %s", reasons)
	}
}

func TestProposalAntiGamingAllowsFocusedBehavioralTest(t *testing.T) {
	target := proposalTarget(t)
	baseline, err := inventory.Build(target)
	if err != nil {
		t.Fatal(err)
	}
	output := provider.Output{
		ChangedFiles: []string{"app_test.py"},
		Patch: "diff --git a/app_test.py b/app_test.py\n--- /dev/null\n+++ b/app_test.py\n" +
			"@@ -0,0 +1,4 @@\n+from app import ready\n+\n+def test_ready():\n+    assert ready() is True\n",
	}
	if reasons := auditProposalAntiGaming(baseline, output); len(reasons) != 0 {
		t.Fatalf("focused behavioral test was rejected: %v", reasons)
	}
}

func TestApplyProviderPatchModifiesOnlyMatchingText(t *testing.T) {
	target := t.TempDir()
	writeTestFile(t, target, "app.py", "def ready():\n    return False\n")
	baseline, err := inventory.Build(target)
	if err != nil {
		t.Fatal(err)
	}
	candidateRoot, err := prepareCandidate(baseline, filepath.Join(t.TempDir(), "candidate"))
	if err != nil {
		t.Fatal(err)
	}
	task := provider.Task{AllowedPaths: []string{"app.py"}, ProtectedPaths: []string{".git/"}}
	output := provider.Output{
		ChangedFiles: []string{"app.py"},
		Patch:        "diff --git a/app.py b/app.py\n--- a/app.py\n+++ b/app.py\n@@ -1,2 +1,2 @@\n def ready():\n-    return False\n+    return True\n",
	}
	changes, err := applyProviderPatch(candidateRoot, baseline, task, output, 1, 2)
	if err != nil {
		t.Fatal(err)
	}
	content, err := os.ReadFile(filepath.Join(candidateRoot, "app.py"))
	if err != nil || string(content) != "def ready():\n    return True\n" {
		t.Fatalf("unexpected content %q: %v", content, err)
	}
	if len(changes) != 1 || changes[0].AddedLines != 1 || changes[0].RemovedLines != 1 {
		t.Fatalf("unexpected changes: %+v", changes)
	}
	output.Patch = strings.Replace(output.Patch, "return False", "return missing", 1)
	otherCandidate, err := prepareCandidate(baseline, filepath.Join(t.TempDir(), "candidate"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := applyProviderPatch(otherCandidate, baseline, task, output, 1, 2); err == nil || !strings.Contains(err.Error(), "does not match") {
		t.Fatalf("unexpected error: %v", err)
	}
}
