package remediation

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/MarinJursic/production-readiness-checklist/scanner/catalog"
	"github.com/MarinJursic/production-readiness-checklist/scanner/engine"
	"github.com/MarinJursic/production-readiness-checklist/scanner/inventory"
	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
	"github.com/MarinJursic/production-readiness-checklist/scanner/provider"
	"github.com/MarinJursic/production-readiness-checklist/scanner/verifier"
)

func passingVerifier(t *testing.T) *verifier.Options {
	t.Helper()
	runtimePath := filepath.Join(t.TempDir(), "docker")
	if err := os.WriteFile(runtimePath, []byte("#!/bin/sh\nexit 0\n"), 0o700); err != nil {
		t.Fatal(err)
	}
	options := verifier.Defaults(runtimePath,
		"registry.example/prc/test-verifier@sha256:"+strings.Repeat("a", 64), "")
	return &options
}

func failingVerifier(t *testing.T) *verifier.Options {
	t.Helper()
	runtimePath := filepath.Join(t.TempDir(), "docker")
	script := "#!/bin/sh\nif [ \"$1\" = image ]; then exit 0; fi\nprintf test-failed\nexit 1\n"
	if err := os.WriteFile(runtimePath, []byte(script), 0o700); err != nil {
		t.Fatal(err)
	}
	options := verifier.Defaults(runtimePath,
		"registry.example/prc/test-verifier@sha256:"+strings.Repeat("a", 64), "")
	return &options
}

func unavailableVerifier(t *testing.T) *verifier.Options {
	t.Helper()
	runtimePath := filepath.Join(t.TempDir(), "docker")
	if err := os.WriteFile(runtimePath, []byte("#!/bin/sh\nexit 125\n"), 0o700); err != nil {
		t.Fatal(err)
	}
	options := verifier.Defaults(runtimePath,
		"registry.example/prc/test-verifier@sha256:"+strings.Repeat("a", 64), "")
	return &options
}

func proposalFinding(t *testing.T, item model.Inventory) model.Finding {
	t.Helper()
	c, err := catalog.Load(testCatalogRoot(t))
	if err != nil {
		t.Fatal(err)
	}
	run, err := engine.New(c).Scan("prc/core-repository", item)
	if err != nil {
		t.Fatal(err)
	}
	finding, ok := findingFor(run, agentTestSuiteAssertion)
	if !ok {
		t.Fatal("missing test-suite finding")
	}
	return finding
}

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
	item, err := inventory.Build(target)
	if err != nil {
		t.Fatal(err)
	}
	finding := proposalFinding(t, item)
	draft := provider.Task{
		SchemaVersion: provider.TaskSchema, Mode: "suggest",
		WorkspaceInventoryDigest: strings.Repeat("b", 64), AssertionID: "PRC-A-CORE-010",
		FindingID: finding.ID, FindingFingerprint: finding.Fingerprint,
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
		MaxFiles: 2, MaxChangedLines: 10, Verifier: passingVerifier(t),
	})
	if err != nil {
		t.Fatal(err)
	}
	if !candidate.Accepted || candidate.Contract.RemediationClass != "R2" || candidate.Contract.ProposalTaskID != task.TaskID ||
		candidate.Contract.ProposalFindingID != task.FindingID || candidate.Contract.FindingFingerprint != task.FindingFingerprint {
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
	if candidate.Verification == nil || candidate.Verification.Outcome != "pass" ||
		candidate.Verification.CandidateInventoryDigest != candidate.CandidateInventoryDigest {
		t.Fatalf("candidate lacks bound independent verification: %+v", candidate.Verification)
	}
}

func TestRunProposalRequiresAndEnforcesIndependentVerification(t *testing.T) {
	target := proposalTarget(t)
	task := sealedProposalTask(t, target, []string{"app_test.py"}, defaultProposalProtectedPaths())
	patch := "diff --git a/app_test.py b/app_test.py\n" +
		"new file mode 100644\n--- /dev/null\n+++ b/app_test.py\n" +
		"@@ -0,0 +1,4 @@\n+from app import ready\n+\n+def test_ready():\n+    assert ready() is True\n"
	missingPath := filepath.Join(t.TempDir(), "missing-verifier")
	_, err := RunProposal(ProposalOptions{
		CatalogRoot: testCatalogRoot(t), Target: target, CandidateDir: missingPath,
		ProfileID: "prc/core-repository", Provider: "codex", Task: task,
		Output: proposalOutput(task, "app_test.py", patch), MaxFiles: 2, MaxChangedLines: 10,
	})
	if err == nil || !IsPolicyDenied(err) || !strings.Contains(err.Error(), "sandbox verifier") {
		t.Fatalf("missing verifier error = %v", err)
	}
	if _, statErr := os.Stat(missingPath); !os.IsNotExist(statErr) {
		t.Fatal("missing verifier created a candidate")
	}

	candidate, err := RunProposal(ProposalOptions{
		CatalogRoot: testCatalogRoot(t), Target: target, CandidateDir: filepath.Join(t.TempDir(), "failed-tests"),
		ProfileID: "prc/core-repository", Provider: "codex", Task: task,
		Output: proposalOutput(task, "app_test.py", patch), MaxFiles: 2, MaxChangedLines: 10,
		Verifier: failingVerifier(t),
	})
	if err != nil {
		t.Fatal(err)
	}
	if candidate.Accepted || candidate.Verification == nil || candidate.Verification.Outcome != "fail" ||
		!strings.Contains(strings.Join(candidate.Reasons, " "), "tests failed") {
		t.Fatalf("failing tests were not preserved as a rejected candidate: %+v", candidate)
	}
}

func TestRunProposalOCIIntegrationAcceptsPassingJavaScriptCandidate(t *testing.T) {
	image := os.Getenv("PRC_TEST_VERIFIER_IMAGE")
	if image == "" {
		t.Skip("set PRC_TEST_VERIFIER_IMAGE to an already-present digest-pinned Node image")
	}
	target := remediationTarget(t)
	if err := os.Remove(filepath.Join(target, "tests", "test_app.py")); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(filepath.Join(target, "app.py")); err != nil {
		t.Fatal(err)
	}
	writeTestFile(t, target, "ready.js", "exports.ready = () => true;\n")
	item, err := inventory.Build(target)
	if err != nil {
		t.Fatal(err)
	}
	c, err := catalog.Load(testCatalogRoot(t))
	if err != nil {
		t.Fatal(err)
	}
	run, err := engine.New(c).Scan("prc/core-repository", item)
	if err != nil {
		t.Fatal(err)
	}
	finding, ok := findingFor(run, agentTestSuiteAssertion)
	if !ok {
		t.Fatal("missing test-suite finding")
	}
	protectedPaths, err := RequiredProtectedPaths(target, nil)
	if err != nil {
		t.Fatal(err)
	}
	task, supported, err := planAgentTask(item, c.Assertions[agentTestSuiteAssertion], finding, protectedPaths,
		AgentOptions{Provider: "codex", AllowRemoteSourceProcessing: true})
	if err != nil || !supported || task.RelevantPaths[0] != "ready.js" || len(task.AllowedPaths) != 1 {
		t.Fatalf("task=%+v supported=%t err=%v", task, supported, err)
	}
	testPath := task.AllowedPaths[0]
	patch := "diff --git a/" + testPath + " b/" + testPath + "\n" +
		"new file mode 100644\n--- /dev/null\n+++ b/" + testPath + "\n" +
		"@@ -0,0 +1,4 @@\n+const test = require('node:test');\n+const assert = require('node:assert/strict');\n+const { ready } = require('../ready');\n+test('ready', () => assert.equal(ready(), true));\n"
	verification := verifier.Defaults("docker", image, "")
	candidate, err := RunProposal(ProposalOptions{
		CatalogRoot: testCatalogRoot(t), Target: target, CandidateDir: filepath.Join(t.TempDir(), "candidate"),
		ProfileID: "prc/core-repository", Provider: "codex", Task: task,
		Output: proposalOutput(task, testPath, patch), MaxFiles: 2, MaxChangedLines: 10,
		Verifier: &verification,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !candidate.Accepted || candidate.Verification == nil || candidate.Verification.Kind != "javascript" ||
		candidate.Verification.Outcome != "pass" {
		t.Fatalf("candidate = %+v", candidate)
	}
}

func TestRunProposalBindsConfiguredPolicyAndInventory(t *testing.T) {
	target := proposalTarget(t)
	configuration := configuredRemediation(t, target, nil)
	item, err := inventory.Build(target)
	if err != nil {
		t.Fatal(err)
	}
	item, err = inventory.BindConfiguration(item, configuration.Validation, configuration.SourcePath)
	if err != nil {
		t.Fatal(err)
	}
	finding := proposalFinding(t, item)
	draft := provider.Task{
		SchemaVersion: provider.TaskSchema, Mode: "suggest", WorkspaceInventoryDigest: strings.Repeat("b", 64),
		FindingID: finding.ID, FindingFingerprint: finding.Fingerprint,
		AssertionID: "PRC-A-CORE-010", ControlIDs: []string{"USEQ-12655775"},
		Goal: "Add one focused, discoverable regression test.", ReadScope: "task-inputs",
		RelevantPaths: []string{"app.py"}, Inputs: []provider.InputFile{}, AllowedPaths: []string{"app_test.py"},
		ProtectedPaths: []string{".git/"}, AllowedCommands: [][]string{}, Network: "deny", Secrets: "none",
		AllowRemoteSourceProcessing: true, TimeoutSeconds: 60, MaxOutputBytes: 256 * 1024,
	}
	data, err := json.Marshal(draft)
	if err != nil {
		t.Fatal(err)
	}
	draftPath := filepath.Join(t.TempDir(), "draft.json")
	if err := os.WriteFile(draftPath, data, 0o600); err != nil {
		t.Fatal(err)
	}
	protectedPaths, err := RequiredProtectedPaths(target, configuration)
	if err != nil {
		t.Fatal(err)
	}
	task, err := provider.SealTaskWithInventory(draftPath, target, item, protectedPaths)
	if err != nil {
		t.Fatal(err)
	}
	patch := "diff --git a/app_test.py b/app_test.py\n" +
		"new file mode 100644\n--- /dev/null\n+++ b/app_test.py\n" +
		"@@ -0,0 +1,4 @@\n+from app import ready\n+\n+def test_ready():\n+    assert ready() is True\n"
	candidate, err := RunProposal(ProposalOptions{
		CatalogRoot: testCatalogRoot(t), Target: target, CandidateDir: filepath.Join(t.TempDir(), "candidate"),
		ProfileID: "prc/core-repository", Provider: "codex", Task: task,
		Output: proposalOutput(task, "app_test.py", patch), MaxFiles: 2, MaxChangedLines: 10,
		Attempt: 1, MaxAttempts: 3, Configuration: configuration, Verifier: passingVerifier(t),
	})
	if err != nil {
		t.Fatal(err)
	}
	if !candidate.Accepted || candidate.Contract.ConfigurationDigest != configuration.Validation.Digest ||
		candidate.Contract.ProjectID != "example-product" || candidate.Contract.MaxAttempts != 3 {
		t.Fatalf("configured proposal candidate lost policy identity: %+v", candidate)
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
		MaxFiles: 2, MaxChangedLines: 10, Verifier: passingVerifier(t),
	})
	if err == nil || !strings.Contains(err.Error(), "hunk counts") {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, statErr := os.Stat(filepath.Join(candidatePath, "app_test.py")); !os.IsNotExist(statErr) {
		t.Fatal("malformed proposal wrote its target file")
	}
}

func TestRunProposalRejectsTaskForDifferentFindingFingerprint(t *testing.T) {
	target := proposalTarget(t)
	task := sealedProposalTask(t, target, []string{"app_test.py"}, defaultProposalProtectedPaths())
	task.FindingFingerprint = strings.Repeat("e", 64)
	identifier, err := provider.TaskID(task)
	if err != nil {
		t.Fatal(err)
	}
	task.TaskID = identifier
	patch := "diff --git a/app_test.py b/app_test.py\n" +
		"new file mode 100644\n--- /dev/null\n+++ b/app_test.py\n" +
		"@@ -0,0 +1,4 @@\n+from app import ready\n+\n+def test_ready():\n+    assert ready() is True\n"
	output := proposalOutput(task, "app_test.py", patch)
	_, err = RunProposal(ProposalOptions{
		CatalogRoot: testCatalogRoot(t), Target: target, CandidateDir: filepath.Join(t.TempDir(), "candidate"),
		Provider: "codex", Task: task, Output: output, MaxFiles: 2, MaxChangedLines: 10, Verifier: passingVerifier(t),
	})
	if err == nil || !strings.Contains(err.Error(), "exact current finding") || !IsPolicyDenied(err) {
		t.Fatalf("unexpected error: %v", err)
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
		MaxFiles: 2, MaxChangedLines: 10, Verifier: passingVerifier(t),
	})
	if err == nil || !strings.Contains(err.Error(), "outside the R2 fix contract") || !IsPolicyDenied(err) {
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
		MaxFiles: 2, MaxChangedLines: 10, Verifier: passingVerifier(t),
	})
	if err == nil || !strings.Contains(err.Error(), "constant assertion") || !IsPolicyDenied(err) {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, statErr := os.Stat(candidatePath); !os.IsNotExist(statErr) {
		t.Fatal("anti-gaming rejection created a candidate directory")
	}
}

func TestRunProposalRejectsTestWithoutCollectableDeclarationOrBehaviorCheck(t *testing.T) {
	tests := []struct {
		name, patch, reason string
	}{
		{
			name: "comment only",
			patch: "diff --git a/app_test.py b/app_test.py\n--- /dev/null\n+++ b/app_test.py\n" +
				"@@ -0,0 +1 @@\n+# A test will be added later.\n",
			reason: "without a recognized test declaration",
		},
		{
			name: "invocation only",
			patch: "diff --git a/app_test.py b/app_test.py\n--- /dev/null\n+++ b/app_test.py\n" +
				"@@ -0,0 +1,4 @@\n+from app import ready\n+\n+def test_ready():\n+    ready()\n",
			reason: "without a recognized behavioral assertion",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			target := proposalTarget(t)
			task := sealedProposalTask(t, target, []string{"app_test.py"}, defaultProposalProtectedPaths())
			candidatePath := filepath.Join(t.TempDir(), "candidate")
			_, err := RunProposal(ProposalOptions{
				CatalogRoot: testCatalogRoot(t), Target: target, CandidateDir: candidatePath,
				Provider: "codex", Task: task, Output: proposalOutput(task, "app_test.py", test.patch),
				MaxFiles: 2, MaxChangedLines: 10, Verifier: passingVerifier(t),
			})
			if err == nil || !strings.Contains(err.Error(), test.reason) || !IsPolicyDenied(err) {
				t.Fatalf("unexpected error: %v", err)
			}
			if _, statErr := os.Stat(candidatePath); !os.IsNotExist(statErr) {
				t.Fatal("anti-gaming rejection created a candidate directory")
			}
		})
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

func TestTestPayloadAuditRejectsExecutableCapabilities(t *testing.T) {
	tests := []struct {
		path     string
		content  string
		category string
	}{
		{"test_app.py", "import subprocess\n\ndef test_ready():\n    assert subprocess.run(['true']).returncode == 0\n", "process execution"},
		{"test_app.py", "import requests\n\ndef test_ready():\n    assert requests.get('https://example.test').ok\n", "network access"},
		{"app.test.js", "const test = require('node:test');\nconst fs = require('node:fs');\ntest('ready', () => fs.writeFileSync('/tmp/x', 'x'));\n", "filesystem access"},
		{"app_test.go", "package app\nimport (\"os\"; \"testing\")\nfunc TestReady(t *testing.T) { if os.Getenv(\"TOKEN\") == \"\" { t.Fail() } }\n", "environment access"},
		{"app_test.py", "def test_ready():\n    assert eval('1 + 1') == 2\n", "dynamic or encoded"},
	}
	for _, test := range tests {
		t.Run(test.path+"-"+test.category, func(t *testing.T) {
			reasons := strings.Join(auditTestPayload(test.path, []byte(test.content)), " ")
			if !strings.Contains(reasons, test.category) {
				t.Fatalf("payload category %q was missed: %s", test.category, reasons)
			}
		})
	}
}

func TestTestPayloadAuditAllowsFocusedLocalAssertions(t *testing.T) {
	for path, content := range map[string]string{
		"app_test.py": "from app import ready\n\ndef test_ready():\n    assert ready() is True\n",
		"app.test.js": "const test = require('node:test');\nconst assert = require('node:assert/strict');\ntest('ready', () => assert.equal(2 + 2, 4));\n",
		"app_test.go": "package app\nimport \"testing\"\nfunc TestReady(t *testing.T) { if !ready() { t.Fail() } }\n",
	} {
		if reasons := auditTestPayload(path, []byte(content)); len(reasons) != 0 {
			t.Fatalf("focused %s test was rejected: %v", path, reasons)
		}
	}
}

func TestRunProposalRejectsExecutableTestPayloadBeforeCandidateOrVerifier(t *testing.T) {
	target := proposalTarget(t)
	task := sealedProposalTask(t, target, []string{"app_test.py"}, defaultProposalProtectedPaths())
	patch := "diff --git a/app_test.py b/app_test.py\n" +
		"new file mode 100644\n--- /dev/null\n+++ b/app_test.py\n" +
		"@@ -0,0 +1,5 @@\n+import subprocess\n+from app import ready\n+\n+def test_ready():\n+    assert subprocess.run(['true']).returncode == 0 and ready() is True\n"
	candidatePath := filepath.Join(t.TempDir(), "candidate")
	_, err := RunProposal(ProposalOptions{
		CatalogRoot: testCatalogRoot(t), Target: target, CandidateDir: candidatePath,
		Provider: "codex", Task: task, Output: proposalOutput(task, "app_test.py", patch),
		MaxFiles: 2, MaxChangedLines: 10,
	})
	if err == nil || !IsPolicyDenied(err) || !strings.Contains(err.Error(), "prohibited process execution") {
		t.Fatalf("unexpected executable-payload result: %v", err)
	}
	if _, statErr := os.Stat(candidatePath); !os.IsNotExist(statErr) {
		t.Fatal("executable payload created a candidate workspace")
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
	changes, err := applyProviderPatch(candidateRoot, baseline, task, output, nil, 1, 2)
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
	if _, err := applyProviderPatch(otherCandidate, baseline, task, output, nil, 1, 2); err == nil || !strings.Contains(err.Error(), "does not match") {
		t.Fatalf("unexpected error: %v", err)
	}
}
