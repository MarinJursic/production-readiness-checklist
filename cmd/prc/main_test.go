package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/MarinJursic/production-readiness-checklist/scanner/provider"
)

func TestVersionCommand(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if code := run([]string{"version"}, &stdout, &stderr); code != 0 {
		t.Fatalf("exit=%d stderr=%s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "prc 0.1.0-dev") {
		t.Fatalf("unexpected output %q", stdout.String())
	}
}

func TestAdapterValidateOutputCommand(t *testing.T) {
	var stdout, stderr bytes.Buffer
	path := filepath.Join("..", "..", "fixtures", "adapters", "valid-output.jsonl")
	if code := run([]string{"adapter", "validate-output", "--file", path}, &stdout, &stderr); code != 0 {
		t.Fatalf("exit=%d stderr=%s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), `"status": "completed"`) {
		t.Fatalf("unexpected output %q", stdout.String())
	}
}

func TestAdapterValidateOutputRejectsAuthorityAttack(t *testing.T) {
	var stdout, stderr bytes.Buffer
	path := filepath.Join("..", "..", "fixtures", "adapters", "malicious-authority-output.jsonl")
	if code := run([]string{"adapter", "validate-output", "--file", path}, &stdout, &stderr); code != 2 {
		t.Fatalf("exit=%d stderr=%s", code, stderr.String())
	}
}

func TestProviderCapabilitiesAreReadOnly(t *testing.T) {
	for _, name := range []string{"codex", "claude"} {
		t.Run(name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			if code := run([]string{"provider", "capabilities", "--provider", name}, &stdout, &stderr); code != 0 {
				t.Fatalf("exit=%d stderr=%s", code, stderr.String())
			}
			for _, forbidden := range []string{`"workspace_mutation": true`, `"network_tools": true`, `"shell": true`} {
				if strings.Contains(stdout.String(), forbidden) {
					t.Fatalf("unsafe capability in %s", stdout.String())
				}
			}
		})
	}
}

func TestUnknownCommandIsUsageError(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if code := run([]string{"unknown"}, &stdout, &stderr); code != 2 {
		t.Fatalf("exit=%d stderr=%s", code, stderr.String())
	}
}

func TestRemediateCommandCreatesAcceptedCandidate(t *testing.T) {
	target := t.TempDir()
	if err := os.WriteFile(filepath.Join(target, "app.py"), []byte("print('ready')"), 0o644); err != nil {
		t.Fatal(err)
	}
	candidate := filepath.Join(t.TempDir(), "candidate")
	var stdout, stderr bytes.Buffer
	code := run([]string{
		"remediate", "--target", target, "--catalog-root", filepath.Join("..", ".."),
		"--candidate-dir", candidate, "--format", "json",
	}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit=%d stderr=%s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), `"accepted": true`) {
		t.Fatalf("unexpected output %q", stdout.String())
	}
	original, err := os.ReadFile(filepath.Join(target, "app.py"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.HasSuffix(string(original), "\n") {
		t.Fatal("command modified original target")
	}
}

func TestRemediateProposalCommandCreatesAcceptedIsolatedCandidate(t *testing.T) {
	target := t.TempDir()
	if err := os.WriteFile(filepath.Join(target, "app.py"), []byte("def ready(): return True\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	draft := provider.Task{
		SchemaVersion: provider.TaskSchema, Mode: "suggest", WorkspaceInventoryDigest: strings.Repeat("b", 64),
		AssertionID: "PRC-A-CORE-010", ControlIDs: []string{"USEQ-12655775"},
		Goal: "Add one focused test.", ReadScope: "task-inputs", RelevantPaths: []string{"app.py"},
		Inputs: []provider.InputFile{}, AllowedPaths: []string{"app_test.py"},
		ProtectedPaths:  []string{".git/", ".github/workflows/", ".prc/", "catalog/", "production-readiness.yaml", "schemas/"},
		AllowedCommands: [][]string{}, Network: "deny", Secrets: "none",
		AllowRemoteSourceProcessing: true, TimeoutSeconds: 60, MaxOutputBytes: 256 * 1024,
	}
	draftData, _ := json.Marshal(draft)
	draftPath := filepath.Join(t.TempDir(), "draft.json")
	if err := os.WriteFile(draftPath, draftData, 0o600); err != nil {
		t.Fatal(err)
	}
	task, err := provider.SealTask(draftPath, target)
	if err != nil {
		t.Fatal(err)
	}
	taskData, _ := json.Marshal(task)
	taskPath := filepath.Join(t.TempDir(), "task.json")
	if err := os.WriteFile(taskPath, taskData, 0o600); err != nil {
		t.Fatal(err)
	}
	proposal := provider.Output{
		SchemaVersion: provider.OutputSchema, TaskID: task.TaskID, Status: "candidate",
		RootCause: "No test is present.", ChangedFiles: []string{"app_test.py"},
		Patch:                  "diff --git a/app_test.py b/app_test.py\n--- /dev/null\n+++ b/app_test.py\n@@ -0,0 +1,2 @@\n+def test_ready():\n+    assert True\n",
		CommandsRequestedOrRun: []provider.CommandResult{}, Limitations: []string{}, RequestedCapabilityChanges: []string{},
	}
	proposalData, _ := json.Marshal(proposal)
	proposalPath := filepath.Join(t.TempDir(), "proposal.json")
	if err := os.WriteFile(proposalPath, proposalData, 0o600); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	candidatePath := filepath.Join(t.TempDir(), "candidate")
	code := run([]string{
		"remediate-proposal", "--target", target, "--catalog-root", filepath.Join("..", ".."),
		"--provider", "codex", "--task", taskPath, "--output", proposalPath,
		"--candidate-dir", candidatePath, "--format", "json",
	}, &stdout, &stderr)
	if code != 0 || !strings.Contains(stdout.String(), `"accepted": true`) {
		t.Fatalf("exit=%d stderr=%s stdout=%s", code, stderr.String(), stdout.String())
	}
	if _, err := os.Stat(filepath.Join(target, "app_test.py")); !os.IsNotExist(err) {
		t.Fatal("source workspace was modified")
	}
}
