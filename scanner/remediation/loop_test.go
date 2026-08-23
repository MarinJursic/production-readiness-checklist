package remediation

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func loopClock() func() time.Time {
	current := time.Date(2026, 8, 23, 12, 0, 0, 0, time.UTC)
	return func() time.Time {
		value := current
		current = current.Add(time.Second)
		return value
	}
}

func TestRunLoopClosesAllRegisteredR1FindingsInIsolatedCandidates(t *testing.T) {
	target := remediationTarget(t)
	path := filepath.Join(target, "app.py")
	if err := os.Chmod(path, 0o666); err != nil {
		t.Fatal(err)
	}
	result, err := RunLoop(LoopOptions{
		CatalogRoot: testCatalogRoot(t), Target: target, CandidateRoot: filepath.Join(t.TempDir(), "candidates"),
		ProfileID: "prc/core-repository", MaxAttempts: 3, MaxFiles: 20, MaxChangedLines: 20, Now: loopClock(),
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Candidates) != 2 || result.Usage.Attempts != 2 || result.Usage.ChangedFiles != 2 ||
		result.Usage.ChangedLines != 1 || len(result.Attempts) != 2 || !result.OriginalUnchanged ||
		result.TerminalState != "machine_work_complete" {
		t.Fatalf("unexpected loop result: %+v", result)
	}
	if result.Attempts[0].Mode != "deterministic" || result.Attempts[0].Outcome != "accepted" ||
		result.Attempts[0].CandidateID != result.Candidates[0].CandidateID ||
		result.Attempts[1].BeforeInventoryDigest != result.Candidates[0].CandidateInventoryDigest ||
		result.Attempts[1].CandidateID != result.Candidates[1].CandidateID {
		t.Fatalf("attempt audit linkage is incomplete: %+v", result.Attempts)
	}
	if result.FinalRun.Inventory.TargetName != filepath.Base(target) {
		t.Fatalf("candidate changed logical target name to %q", result.FinalRun.Inventory.TargetName)
	}
	if result.Candidates[0].Contract.AssertionID != finalNewlineAssertion ||
		result.Candidates[1].Contract.AssertionID != restrictiveModesAssertion ||
		result.Candidates[0].Contract.Attempt != 1 || result.Candidates[1].Contract.Attempt != 2 {
		t.Fatalf("unexpected candidate sequence: %+v", result.Candidates)
	}
	original, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	originalInfo, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	finalPath := filepath.Join(result.ResultWorkspace, "app.py")
	final, err := os.ReadFile(finalPath)
	if err != nil {
		t.Fatal(err)
	}
	finalInfo, err := os.Stat(finalPath)
	if err != nil {
		t.Fatal(err)
	}
	if strings.HasSuffix(string(original), "\n") || originalInfo.Mode().Perm() != 0o666 ||
		!strings.HasSuffix(string(final), "\n") || finalInfo.Mode().Perm() != 0o644 {
		t.Fatalf("original/final state mismatch: original=%q/%#o final=%q/%#o",
			original, originalInfo.Mode().Perm(), final, finalInfo.Mode().Perm())
	}
	for _, item := range result.Remaining {
		if item.AssertionID == finalNewlineAssertion || item.AssertionID == restrictiveModesAssertion {
			t.Fatalf("closed R1 finding remains: %+v", item)
		}
	}
	wantID, err := remediationRunID(result)
	if err != nil || wantID != result.RunID {
		t.Fatalf("run identity mismatch: %v %s", err, wantID)
	}
	tampered := result
	tampered.Attempts = append([]AttemptRecord(nil), result.Attempts...)
	tampered.Attempts[1].BeforeInventoryDigest = result.SourceInventoryDigest
	if err := validateAttemptAudit(tampered); err == nil || !strings.Contains(err.Error(), "identity, ordering") {
		t.Fatalf("attempt-chain tampering was accepted: %v", err)
	}
}

func TestRunLoopStopsPredictablyAtAttemptBudget(t *testing.T) {
	target := remediationTarget(t)
	if err := os.Chmod(filepath.Join(target, "app.py"), 0o666); err != nil {
		t.Fatal(err)
	}
	result, err := RunLoop(LoopOptions{
		CatalogRoot: testCatalogRoot(t), Target: target, CandidateRoot: filepath.Join(t.TempDir(), "candidates"),
		ProfileID: "prc/core-repository", MaxAttempts: 1, MaxFiles: 20, MaxChangedLines: 20, Now: loopClock(),
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.TerminalState != "stopped_by_policy_or_budget" || len(result.Candidates) != 1 || len(result.StopReasons) != 1 {
		t.Fatalf("unexpected budget stop: %+v", result)
	}
	found := false
	for _, item := range result.Remaining {
		if item.AssertionID == restrictiveModesAssertion && item.ReasonCode == "budget_exhausted" {
			found = true
		}
	}
	if !found {
		t.Fatalf("remaining work does not explain budget stop: %+v", result.Remaining)
	}
}

func TestRunLoopRebasesInternalConfigurationAcrossCandidates(t *testing.T) {
	target := remediationTarget(t)
	if err := os.Chmod(filepath.Join(target, "app.py"), 0o666); err != nil {
		t.Fatal(err)
	}
	configuration := configuredRemediation(t, target, nil)
	result, err := RunLoop(LoopOptions{
		CatalogRoot: testCatalogRoot(t), Target: target, CandidateRoot: filepath.Join(t.TempDir(), "candidates"),
		ProfileID: "prc/core-repository", MaxAttempts: 3, MaxFiles: 20, MaxChangedLines: 200,
		Configuration: configuration, Now: loopClock(),
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Candidates) != 2 || result.ConfigurationDigest != configuration.Validation.Digest ||
		result.ProjectID != "example-product" {
		t.Fatalf("configured loop identity mismatch: %+v", result)
	}
	for _, candidate := range result.Candidates {
		if candidate.Contract.ConfigurationDigest != configuration.Validation.Digest || candidate.Contract.MaxAttempts != 3 {
			t.Fatalf("candidate lost configured policy: %+v", candidate.Contract)
		}
	}
	source, err := os.ReadFile(configuration.SourcePath)
	if err != nil {
		t.Fatal(err)
	}
	final, err := os.ReadFile(filepath.Join(result.ResultWorkspace, "production-readiness.yaml"))
	if err != nil || string(final) != string(source) {
		t.Fatalf("configuration was not preserved across candidates: %v", err)
	}
}

func TestRunLoopReportsMachineWorkCompleteWithoutCandidates(t *testing.T) {
	target := remediationTarget(t)
	if err := os.WriteFile(filepath.Join(target, "app.py"), []byte("def ready(): return True\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	result, err := RunLoop(LoopOptions{
		CatalogRoot: testCatalogRoot(t), Target: target, CandidateRoot: filepath.Join(t.TempDir(), "candidates"),
		ProfileID: "prc/core-repository", MaxAttempts: 3, MaxFiles: 20, MaxChangedLines: 20, Now: loopClock(),
	})
	if err != nil {
		t.Fatal(err)
	}
	resolvedTarget, err := filepath.EvalSymlinks(target)
	if err != nil {
		t.Fatal(err)
	}
	if result.TerminalState != "machine_work_complete" || len(result.Candidates) != 0 ||
		result.Attempts == nil || len(result.Attempts) != 0 || result.Candidates == nil || result.Remaining == nil ||
		result.StopReasons == nil || result.ResultWorkspace != resolvedTarget {
		t.Fatalf("unexpected no-op loop state: terminal=%s candidates=%d workspace=%s",
			result.TerminalState, len(result.Candidates), result.ResultWorkspace)
	}
}

func TestRunLoopRejectsCandidateRootInsideTarget(t *testing.T) {
	target := remediationTarget(t)
	_, err := RunLoop(LoopOptions{
		CatalogRoot: testCatalogRoot(t), Target: target, CandidateRoot: filepath.Join(target, "candidates"),
		ProfileID: "prc/core-repository", MaxAttempts: 3, MaxFiles: 20, MaxChangedLines: 20,
	})
	if err == nil || !strings.Contains(err.Error(), "outside the target tree") {
		t.Fatalf("unexpected error: %v", err)
	}
}
