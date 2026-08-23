package remediation

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/MarinJursic/production-readiness-checklist/scanner/catalog"
	"github.com/MarinJursic/production-readiness-checklist/scanner/engine"
	"github.com/MarinJursic/production-readiness-checklist/scanner/inventory"
	"github.com/MarinJursic/production-readiness-checklist/scanner/provider"
)

func fakeAgentExecutable(t *testing.T, name, status string) string {
	t.Helper()
	patch := `diff --git a/app_test.py b/app_test.py\nnew file mode 100644\n--- /dev/null\n+++ b/app_test.py\n@@ -0,0 +1,4 @@\n+from app import ready\n+\n+def test_ready():\n+    assert ready() is True\n`
	return fakeAgentExecutableWithPatch(t, name, status, patch)
}

func fakeAgentExecutableWithPatch(t *testing.T, name, status, patch string) string {
	t.Helper()
	body := `prompt=$(/bin/cat)
task_id=$(printf '%s' "$prompt" | /usr/bin/sed -n 's/.*"task_id": "\([0-9a-f]*\)".*/\1/p' | /usr/bin/head -n 1)
result=""
while [ "$#" -gt 0 ]; do
  if [ "$1" = "--output-last-message" ]; then shift; result="$1"; fi
  shift
done
printf '%s' '{"schema_version":"prc.agent-output/v0.1","task_id":"'"$task_id"'","status":"` + status + `","root_cause":"No recognized automated test exists.","changed_files":` + map[string]string{
		"candidate": `["app_test.py"]`, "unable": `[]`, "needs_escalation": `[]`,
	}[status] + `,"patch":"` + map[string]string{
		"candidate": patch, "unable": ``, "needs_escalation": ``,
	}[status] + `","commands_requested_or_run":[],"limitations":[],"requested_capability_changes":[]}' > "$result"
printf '%s\n' '{"type":"turn.completed"}'
`
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, []byte("#!/bin/sh\nset -eu\n"+body), 0o755); err != nil {
		t.Fatal(err)
	}
	return path
}

func agentLoopTarget(t *testing.T) string {
	t.Helper()
	target := proposalTarget(t)
	if err := os.WriteFile(filepath.Join(target, "app.py"), []byte("def ready(): return True\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	return target
}

func agentLoopOptions(t *testing.T, target, status string) LoopOptions {
	t.Helper()
	root := testCatalogRoot(t)
	return LoopOptions{
		CatalogRoot: root, Target: target, CandidateRoot: filepath.Join(t.TempDir(), "candidates"),
		ProfileID: "prc/core-repository", MaxAttempts: 2, MaxFiles: 4, MaxChangedLines: 20,
		Agent: &AgentOptions{
			Provider: "codex", Executable: fakeAgentExecutable(t, "codex", status),
			OutputSchemaPath:            filepath.Join(root, "schemas", "agent-output.schema.json"),
			AllowRemoteSourceProcessing: true, TimeoutSeconds: 30, MaxOutputBytes: 64 * 1024,
		},
		Verifier: passingVerifier(t),
		Now:      loopClock(),
	}
}

func TestPlanAgentTaskBindsOneSourceAndNewTestAllowlist(t *testing.T) {
	target := agentLoopTarget(t)
	item, err := inventory.Build(target)
	if err != nil {
		t.Fatal(err)
	}
	c, err := catalog.Load(testCatalogRoot(t))
	if err != nil {
		t.Fatal(err)
	}
	paths, err := RequiredProtectedPaths(target, nil)
	if err != nil {
		t.Fatal(err)
	}
	run, err := engine.New(c).Scan("prc/core-repository", item)
	if err != nil {
		t.Fatal(err)
	}
	baselineFinding, ok := findingFor(run, agentTestSuiteAssertion)
	if !ok {
		t.Fatal("missing test-suite finding")
	}
	task, supported, err := planAgentTask(item, c.Assertions[agentTestSuiteAssertion], baselineFinding, paths, AgentOptions{
		Provider: "codex", AllowRemoteSourceProcessing: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !supported || task.AssertionID != agentTestSuiteAssertion || task.RelevantPaths[0] != "app.py" ||
		len(task.AllowedPaths) != 1 || !slices.Contains(task.AllowedPaths, "app_test.py") || len(task.Inputs) != 1 ||
		task.Inputs[0].Content != "def ready(): return True\n" || task.TaskID == "" ||
		task.FindingID != baselineFinding.ID || task.FindingFingerprint != baselineFinding.Fingerprint {
		t.Fatalf("unexpected task: %+v", task)
	}
}

func TestRunLoopExecutesReadOnlyProviderAndAcceptsBoundedR2Candidate(t *testing.T) {
	target := agentLoopTarget(t)
	result, err := RunLoop(agentLoopOptions(t, target, "candidate"))
	if err != nil {
		t.Fatal(err)
	}
	if result.SchemaVersion != RunSchema || len(result.ProviderExecutions) != 1 || len(result.Candidates) != 1 ||
		result.Usage.Attempts != 1 || len(result.Attempts) != 1 || !result.Candidates[0].Accepted ||
		result.Candidates[0].Contract.AssertionID != agentTestSuiteAssertion ||
		result.Candidates[0].Contract.Provider != "codex" {
		t.Fatalf("unexpected agent loop result: %+v", result)
	}
	execution := result.ProviderExecutions[0]
	if result.Attempts[0].Mode != "agent" || result.Attempts[0].Outcome != "accepted" ||
		result.Attempts[0].ProviderExecutionID != execution.ExecutionID ||
		result.Attempts[0].CandidateID != result.Candidates[0].CandidateID ||
		result.Candidates[0].Verification == nil ||
		result.Attempts[0].VerificationExecutionID != result.Candidates[0].Verification.ExecutionID {
		t.Fatalf("agent attempt audit linkage is incomplete: %+v", result.Attempts[0])
	}
	tampered := result
	tampered.Attempts = append([]AttemptRecord(nil), result.Attempts...)
	tampered.Attempts[0].VerificationExecutionID = strings.Repeat("f", 64)
	if validateAttemptAudit(tampered) == nil {
		t.Fatal("attempt audit accepted a forged verification linkage")
	}
	if execution.TaskID != result.Candidates[0].Contract.ProposalTaskID || execution.Output.Status != "candidate" {
		t.Fatalf("provider execution is not bound to candidate: %+v", execution)
	}
	if result.Candidates[0].Contract.ProposalFindingID == "" ||
		result.Candidates[0].Contract.ProposalFindingID != agentTaskFindingID(t, execution) {
		t.Fatalf("provider candidate lost triggering finding identity: %+v", result.Candidates[0].Contract)
	}
	if _, err := os.Stat(filepath.Join(filepath.Dir(execution.StdoutPath), "agent-task.json")); err != nil {
		t.Fatalf("sealed task record is missing: %v", err)
	}
	if _, err := os.Stat(filepath.Join(target, "app_test.py")); !os.IsNotExist(err) {
		t.Fatal("agent loop modified the original workspace")
	}
	content, err := os.ReadFile(filepath.Join(result.ResultWorkspace, "app_test.py"))
	if err != nil || !strings.Contains(string(content), "assert ready() is True") {
		t.Fatalf("candidate test = %q, %v", content, err)
	}
	for _, remaining := range result.Remaining {
		if remaining.AssertionID == agentTestSuiteAssertion {
			t.Fatalf("closed test-suite finding remains: %+v", remaining)
		}
	}
}

func TestRunLoopDurationCancelsInFlightProvider(t *testing.T) {
	target := agentLoopTarget(t)
	options := agentLoopOptions(t, target, "candidate")
	providerPath := filepath.Join(t.TempDir(), "codex")
	if err := os.WriteFile(providerPath, []byte("#!/bin/sh\nexec /bin/sleep 5\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	options.Agent.Executable = providerPath
	options.MaxDurationSeconds = 1
	fixed := time.Date(2026, 8, 23, 12, 0, 0, 0, time.UTC)
	options.Now = func() time.Time { return fixed }
	started := time.Now()
	result, err := RunLoop(options)
	if err != nil {
		t.Fatal(err)
	}
	if elapsed := time.Since(started); elapsed > 3*time.Second {
		t.Fatalf("provider exceeded total run deadline: %s", elapsed)
	}
	if result.TerminalState != "stopped_by_policy_or_budget" || result.MaxDurationSeconds != 1 ||
		len(result.ProviderFailures) != 1 || result.ProviderFailures[0].ReasonCode != "cancelled" ||
		len(result.Attempts) != 1 || result.Attempts[0].Outcome != "provider_failed" {
		t.Fatalf("unexpected in-flight duration result: %+v", result)
	}
	found := false
	for _, item := range result.Remaining {
		if item.AssertionID == agentTestSuiteAssertion && item.ReasonCode == "time_budget_exhausted" {
			found = true
		}
	}
	if !found {
		t.Fatalf("remaining work does not explain provider deadline: %+v", result.Remaining)
	}
}

func agentTaskFindingID(t *testing.T, execution provider.Execution) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(filepath.Dir(execution.StdoutPath), "agent-task.json"))
	if err != nil {
		t.Fatal(err)
	}
	var task provider.Task
	if err := json.Unmarshal(data, &task); err != nil {
		t.Fatal(err)
	}
	return task.FindingID
}

func TestRunLoopStopsWithoutPatchWhenProviderIsUnable(t *testing.T) {
	target := agentLoopTarget(t)
	result, err := RunLoop(agentLoopOptions(t, target, "unable"))
	if err != nil {
		t.Fatal(err)
	}
	if result.TerminalState != "provider_stopped" || len(result.ProviderExecutions) != 1 ||
		len(result.Candidates) != 0 || result.Usage.Attempts != 1 || len(result.Attempts) != 1 ||
		len(result.StopReasons) != 1 {
		t.Fatalf("unexpected provider stop: %+v", result)
	}
	if result.Attempts[0].Outcome != "provider_stopped" || result.Attempts[0].ReasonCode != "provider_unable" ||
		result.Attempts[0].CandidateID != "" {
		t.Fatalf("provider stop attempt was not audited: %+v", result.Attempts[0])
	}
	found := false
	for _, remaining := range result.Remaining {
		if remaining.AssertionID == agentTestSuiteAssertion && remaining.ReasonCode == "provider_stopped" {
			found = true
		}
	}
	if !found {
		t.Fatalf("provider stop is not explained: %+v", result.Remaining)
	}
}

func TestRunLoopRecordsProviderProcessFailureWithoutRetry(t *testing.T) {
	target := agentLoopTarget(t)
	options := agentLoopOptions(t, target, "candidate")
	executable := filepath.Join(t.TempDir(), "codex")
	if err := os.WriteFile(executable, []byte("#!/bin/sh\nset -eu\nprintf diagnostic >&2\nexit 9\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	options.Agent.Executable = executable
	result, err := RunLoop(options)
	if err != nil {
		t.Fatal(err)
	}
	if result.TerminalState != "provider_failed" || len(result.ProviderFailures) != 1 ||
		len(result.ProviderExecutions) != 0 || len(result.Candidates) != 0 ||
		result.Usage.Attempts != 1 || len(result.Attempts) != 1 || len(result.StopReasons) != 1 {
		t.Fatalf("unexpected provider failure result: %+v", result)
	}
	failure, attempt := result.ProviderFailures[0], result.Attempts[0]
	if failure.ReasonCode != "process_failed" || !failure.TranscriptsComplete || failure.StderrBytes == 0 ||
		attempt.Outcome != "provider_failed" || attempt.ReasonCode != "provider_process_failed" ||
		attempt.ProviderFailureID != failure.FailureID || attempt.ProviderExecutionID != "" {
		t.Fatalf("provider failure audit linkage is incomplete: failure=%+v attempt=%+v", failure, attempt)
	}
	foundFailureReason := false
	for _, remaining := range result.Remaining {
		if remaining.AssertionID == agentTestSuiteAssertion && remaining.ReasonCode == "provider_failed" &&
			strings.Contains(remaining.Reason, "failed safely") {
			foundFailureReason = true
		}
	}
	if !foundFailureReason {
		t.Fatalf("provider failure is not explained safely: %+v", result.Remaining)
	}
	tampered := result
	tampered.Attempts = append([]AttemptRecord(nil), result.Attempts...)
	tampered.Attempts[0].Reason = "different failure reason"
	if validateAttemptAudit(tampered) == nil {
		t.Fatal("attempt audit accepted a failure reason that did not match its provider record")
	}
	if _, err := os.Stat(filepath.Join(target, "app_test.py")); !os.IsNotExist(err) {
		t.Fatal("failed provider modified the original workspace")
	}
}

func TestRunLoopRecordsAntiGamingProposalAsRejected(t *testing.T) {
	target := agentLoopTarget(t)
	options := agentLoopOptions(t, target, "candidate")
	patch := `diff --git a/app_test.py b/app_test.py\nnew file mode 100644\n--- /dev/null\n+++ b/app_test.py\n@@ -0,0 +1,2 @@\n+def test_ready():\n+    assert True\n`
	options.Agent.Executable = fakeAgentExecutableWithPatch(t, "codex", "candidate", patch)
	result, err := RunLoop(options)
	if err != nil {
		t.Fatal(err)
	}
	if result.TerminalState != "candidate_rejected" || len(result.ProviderExecutions) != 1 ||
		len(result.Candidates) != 0 || result.Usage.Attempts != 1 || len(result.Attempts) != 1 ||
		len(result.StopReasons) != 1 {
		t.Fatalf("unexpected anti-gaming stop: %+v", result)
	}
	if result.Attempts[0].Outcome != "rejected" || result.Attempts[0].ReasonCode != "anti_gaming_rejected" ||
		!strings.Contains(result.Attempts[0].Reason, "constant assertion") ||
		!strings.Contains(result.StopReasons[0], "constant assertion") {
		t.Fatalf("anti-gaming rejection lacks exact audit evidence: %+v", result)
	}
	if _, err := os.Stat(filepath.Join(result.CandidateRoot, "attempt-001")); !os.IsNotExist(err) {
		t.Fatal("anti-gaming rejection created a candidate workspace")
	}
}

func TestRunLoopRequiresRemoteSourceAcknowledgementBeforeCreatingRoot(t *testing.T) {
	target := agentLoopTarget(t)
	options := agentLoopOptions(t, target, "candidate")
	options.Agent.AllowRemoteSourceProcessing = false
	root := options.CandidateRoot
	_, err := RunLoop(options)
	if err == nil || !IsPolicyDenied(err) || !strings.Contains(err.Error(), "allow-remote-source-processing") {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, statErr := os.Stat(root); !os.IsNotExist(statErr) {
		t.Fatal("policy denial created a candidate root")
	}
}

func TestRunLoopRequiresVerifierBeforeCreatingRootOrCallingProvider(t *testing.T) {
	target := agentLoopTarget(t)
	options := agentLoopOptions(t, target, "candidate")
	options.Verifier = nil
	root := options.CandidateRoot
	_, err := RunLoop(options)
	if err == nil || !IsPolicyDenied(err) || !strings.Contains(err.Error(), "sandbox verifier") {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, statErr := os.Stat(root); !os.IsNotExist(statErr) {
		t.Fatal("missing verifier created a candidate root")
	}
}

func TestRunLoopPreflightsVerifierBeforeCreatingRootOrCallingProvider(t *testing.T) {
	target := agentLoopTarget(t)
	options := agentLoopOptions(t, target, "candidate")
	options.Verifier = unavailableVerifier(t)
	options.Agent.Executable = filepath.Join(t.TempDir(), "provider-must-not-run")
	root := options.CandidateRoot
	_, err := RunLoop(options)
	if err == nil || !IsPolicyDenied(err) || !strings.Contains(err.Error(), "image is unavailable") {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, statErr := os.Stat(root); !os.IsNotExist(statErr) {
		t.Fatal("failed verifier preflight created a candidate root")
	}
}

func TestRunLoopRejectsSecretLikeSourceBeforeProviderExecution(t *testing.T) {
	target := agentLoopTarget(t)
	content := "def ready():\n    token = \"" + "ghp_" + strings.Repeat("a", 36) + "\"\n    return True\n"
	if err := os.WriteFile(filepath.Join(target, "app.py"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	options := agentLoopOptions(t, target, "candidate")
	root := options.CandidateRoot
	_, err := RunLoop(options)
	if err == nil || !IsPolicyDenied(err) || !strings.Contains(err.Error(), "app.py (github-token)") ||
		strings.Contains(err.Error(), content) {
		t.Fatalf("unexpected remote-input denial: %v", err)
	}
	if _, statErr := os.Stat(root); !os.IsNotExist(statErr) {
		t.Fatal("remote-input denial created a candidate root")
	}
}

func TestRunLoopDetectsProviderTamperingWithSealedTaskRecord(t *testing.T) {
	target := agentLoopTarget(t)
	options := agentLoopOptions(t, target, "candidate")
	data, err := os.ReadFile(options.Agent.Executable)
	if err != nil {
		t.Fatal(err)
	}
	data = []byte(strings.Replace(string(data), "result=\"\"", "printf tampered > agent-task.json\nresult=\"\"", 1))
	if err := os.WriteFile(options.Agent.Executable, data, 0o755); err != nil {
		t.Fatal(err)
	}
	_, err = RunLoop(options)
	if err == nil || !IsProviderExecution(err) || !strings.Contains(err.Error(), "task record changed") {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, statErr := os.Stat(filepath.Join(target, "app_test.py")); !os.IsNotExist(statErr) {
		t.Fatal("tampering provider changed the original workspace")
	}
}
