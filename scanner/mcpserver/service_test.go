package mcpserver

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

func repositoryRoot(t *testing.T) string {
	t.Helper()
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	return root
}

func writeTargetFile(t *testing.T, directory, name, content string, mode os.FileMode) {
	t.Helper()
	path := filepath.Join(directory, name)
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), mode); err != nil {
		t.Fatal(err)
	}
}

func TestServicePlansScansAndExplainsWithoutMutation(t *testing.T) {
	target := t.TempDir()
	writeTargetFile(t, target, "README.md", "# Example\n", 0o640)
	writeTargetFile(t, target, "LICENSE", "MIT\n", 0o600)
	observedAt := time.Date(2026, 8, 23, 15, 30, 0, 0, time.UTC)
	service, err := NewService(Options{
		CatalogRoot: repositoryRoot(t), Target: target,
		Now: func() time.Time { return observedAt },
	})
	if err != nil {
		t.Fatal(err)
	}

	beforeData, err := os.ReadFile(filepath.Join(target, "README.md"))
	if err != nil {
		t.Fatal(err)
	}
	beforeInfo, err := os.Stat(filepath.Join(target, "README.md"))
	if err != nil {
		t.Fatal(err)
	}
	plan, err := service.Plan()
	if err != nil {
		t.Fatal(err)
	}
	if plan.SchemaVersion != PlanResultSchema || plan.Plan.ExecutionMode != "inspect" ||
		plan.Plan.CapabilityPolicy.Process != "deny" || plan.Plan.CapabilityPolicy.Network != "deny" ||
		plan.Plan.CapabilityPolicy.WriteScratch || len(plan.Plan.Adapters) != 0 {
		t.Fatalf("unsafe MCP plan = %+v", plan.Plan)
	}
	first, err := service.Scan()
	if err != nil {
		t.Fatal(err)
	}
	second, err := service.Scan()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(first, second) {
		t.Fatalf("fixed-time scan is not deterministic\nfirst=%+v\nsecond=%+v", first, second)
	}
	if first.SchemaVersion != ScanResultSchema || first.RunID == "" || first.Inventory.FileCount != 2 ||
		first.Summary.Total != len(first.Results) || first.Summary.Findings != len(first.Findings) ||
		first.StartedAt != observedAt || first.CompletedAt != observedAt {
		t.Fatalf("scan projection = %+v", first)
	}
	afterData, err := os.ReadFile(filepath.Join(target, "README.md"))
	if err != nil {
		t.Fatal(err)
	}
	afterInfo, err := os.Stat(filepath.Join(target, "README.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(afterData) != string(beforeData) || afterInfo.Mode() != beforeInfo.Mode() {
		t.Fatalf("read-only scan mutated target: mode %v -> %v", beforeInfo.Mode(), afterInfo.Mode())
	}

	explanation, err := service.Explain("PRC-A-CORE-001")
	if err != nil {
		t.Fatal(err)
	}
	if explanation.SchemaVersion != ExplainResultSchema || explanation.Assertion.ID != "PRC-A-CORE-001" ||
		len(explanation.Objectives) == 0 {
		t.Fatalf("explanation = %+v", explanation)
	}
	for index := 1; index < len(explanation.Objectives); index++ {
		if explanation.Objectives[index-1].ID >= explanation.Objectives[index].ID {
			t.Fatalf("objectives are not sorted: %+v", explanation.Objectives)
		}
	}
}

func TestServiceReinventoriesExternalEdits(t *testing.T) {
	target := t.TempDir()
	writeTargetFile(t, target, "README.md", "first\n", 0o600)
	service, err := NewService(Options{
		CatalogRoot: repositoryRoot(t), Target: target,
		Now: func() time.Time { return time.Date(2026, 8, 23, 16, 0, 0, 0, time.UTC) },
	})
	if err != nil {
		t.Fatal(err)
	}
	first, err := service.Scan()
	if err != nil {
		t.Fatal(err)
	}
	writeTargetFile(t, target, "README.md", "second\n", 0o600)
	second, err := service.Scan()
	if err != nil {
		t.Fatal(err)
	}
	if first.Inventory.Digest == second.Inventory.Digest || first.RunID == second.RunID || first.PlanDigest == second.PlanDigest {
		t.Fatalf("external edit did not invalidate identities: first=%+v second=%+v", first, second)
	}
}

func TestServiceLocksCanonicalTargetAgainstSymlinkRetargeting(t *testing.T) {
	parent := t.TempDir()
	firstTarget := filepath.Join(parent, "first")
	secondTarget := filepath.Join(parent, "second")
	if err := os.MkdirAll(firstTarget, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(secondTarget, 0o700); err != nil {
		t.Fatal(err)
	}
	writeTargetFile(t, firstTarget, "README.md", "first\n", 0o600)
	writeTargetFile(t, secondTarget, "README.md", "second\n", 0o600)
	writeTargetFile(t, secondTarget, "extra.txt", "must not be scanned\n", 0o600)
	link := filepath.Join(parent, "target")
	if err := os.Symlink(firstTarget, link); err != nil {
		t.Fatal(err)
	}
	service, err := NewService(Options{CatalogRoot: repositoryRoot(t), Target: link})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(link); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(secondTarget, link); err != nil {
		t.Fatal(err)
	}
	plan, err := service.Plan()
	if err != nil {
		t.Fatal(err)
	}
	if plan.Plan.TargetName != filepath.Base(firstTarget) {
		t.Fatalf("service followed retargeted symlink: %s", plan.Plan.TargetName)
	}
}

func TestServiceLocksConfiguredProfile(t *testing.T) {
	target := t.TempDir()
	writeTargetFile(t, target, "README.md", "configured\n", 0o600)
	fixture, err := os.ReadFile(filepath.Join(repositoryRoot(t), "fixtures", "config", "production-readiness.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	configPath := filepath.Join(t.TempDir(), "production-readiness.yaml")
	if err := os.WriteFile(configPath, fixture, 0o600); err != nil {
		t.Fatal(err)
	}
	service, err := NewService(Options{
		CatalogRoot: repositoryRoot(t), Target: target, ConfigPath: configPath,
	})
	if err != nil {
		t.Fatal(err)
	}
	changed := strings.Replace(string(fixture), "prc/core-repository", "prc/other-profile", 1)
	if changed == string(fixture) {
		t.Fatal("configuration fixture did not contain expected profile")
	}
	if err := os.WriteFile(configPath, []byte(changed), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := service.Plan(); err == nil || !strings.Contains(err.Error(), "changed its profile") {
		t.Fatalf("profile mutation error = %v", err)
	}
}

func TestServiceRejectsMismatchedAndUnknownProfiles(t *testing.T) {
	target := t.TempDir()
	writeTargetFile(t, target, "README.md", "profile\n", 0o600)
	configPath := filepath.Join(repositoryRoot(t), "fixtures", "config", "production-readiness.yaml")
	if _, err := NewService(Options{
		CatalogRoot: repositoryRoot(t), Target: target, ConfigPath: configPath, ProfileID: "prc/other-profile",
	}); err == nil || !strings.Contains(err.Error(), "does not match locked profile") {
		t.Fatalf("mismatched profile error = %v", err)
	}
	if _, err := NewService(Options{
		CatalogRoot: repositoryRoot(t), Target: target, ProfileID: "prc/unknown-profile",
	}); err == nil || !strings.Contains(err.Error(), "unknown profile") {
		t.Fatalf("unknown profile error = %v", err)
	}
}
