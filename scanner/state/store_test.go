package state

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/MarinJursic/production-readiness-checklist/scanner/catalog"
	"github.com/MarinJursic/production-readiness-checklist/scanner/engine"
	"github.com/MarinJursic/production-readiness-checklist/scanner/evidence"
	"github.com/MarinJursic/production-readiness-checklist/scanner/inventory"
	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
)

func testRun(t *testing.T) model.RunResult {
	t.Helper()
	target := t.TempDir()
	if err := os.WriteFile(filepath.Join(target, "app.go"), []byte("package app\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	c, err := catalog.Load(root)
	if err != nil {
		t.Fatal(err)
	}
	item, err := inventory.Build(target)
	if err != nil {
		t.Fatal(err)
	}
	run, err := engine.New(c).Scan("prc/core-repository", item)
	if err != nil {
		t.Fatal(err)
	}
	return run
}

func writeAndOpen(t *testing.T, run model.RunResult) (*Store, string) {
	t.Helper()
	root := t.TempDir()
	if err := os.Chmod(root, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := evidence.WriteRun(root, run); err != nil {
		t.Fatal(err)
	}
	store, err := Open(context.Background(), root)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { store.Close() })
	return store, root
}

func expectedEvidenceCount(run model.RunResult) int {
	ids := map[string]bool{}
	for _, result := range run.Results {
		for _, item := range result.EvidenceObserved {
			ids[item.ID] = true
		}
	}
	return len(ids)
}

func TestStoreIndexesQueriesAndReloadsCanonicalRun(t *testing.T) {
	ctx := context.Background()
	run := testRun(t)
	store, _ := writeAndOpen(t, run)
	if err := store.IndexRun(ctx, run); err != nil {
		t.Fatal(err)
	}
	counts, err := store.Counts(ctx)
	if err != nil {
		t.Fatal(err)
	}
	wantCounts := Counts{
		Runs: 1, Results: len(run.Results), Findings: len(run.Findings), Evidence: expectedEvidenceCount(run),
		InventoryFiles: len(run.Inventory.Files), InventoryFacts: len(run.Inventory.Facts), AuditEvents: 1,
	}
	if !reflect.DeepEqual(counts, wantCounts) {
		t.Fatalf("counts = %+v, want %+v", counts, wantCounts)
	}
	if len(run.Findings) > 0 {
		var fingerprint string
		if err := store.db.QueryRowContext(ctx,
			"SELECT fingerprint FROM findings WHERE run_id = ? AND finding_id = ?",
			run.RunID, run.Findings[0].ID).Scan(&fingerprint); err != nil || fingerprint != run.Findings[0].Fingerprint {
			t.Fatalf("indexed finding fingerprint=%q err=%v", fingerprint, err)
		}
	}
	history, err := store.ListRuns(ctx, Query{Limit: 10, TargetName: run.Inventory.TargetName, ProfileID: run.Plan.ProfileID})
	if err != nil {
		t.Fatal(err)
	}
	if len(history) != 1 || history[0].RunID != run.RunID || history[0].TerminalState != run.TerminalState {
		t.Fatalf("history = %+v", history)
	}
	missing, err := store.ListRuns(ctx, Query{Limit: 10, TargetName: "different"})
	if err != nil || len(missing) != 0 {
		t.Fatalf("filtered history = %+v, %v", missing, err)
	}
	loaded, err := store.LoadRun(ctx, run.RunID)
	if err != nil || loaded.RunID != run.RunID || loaded.Inventory.Root != "" {
		t.Fatalf("loaded run = %s root=%q err=%v", loaded.RunID, loaded.Inventory.Root, err)
	}
	if err := store.IntegrityCheck(ctx); err != nil {
		t.Fatal(err)
	}
}

func TestStoreReadsLegacyV03RunAfterRuleBindingUpgrade(t *testing.T) {
	run := testRun(t)
	run.SchemaVersion = "prc.run/v0.3"
	run.Plan.SchemaVersion = "prc.plan/v0.3"
	run.Plan.EngineVersion = ""
	run.Plan.ProfileDigest = ""
	run.Plan.CatalogDigest = ""
	for index := range run.Plan.Assertions {
		run.Plan.Assertions[index].AssertionRevision = 0
		run.Plan.Assertions[index].DefinitionDigest = ""
	}
	run.Plan.Digest = ""
	payload, err := json.Marshal(run.Plan)
	if err != nil {
		t.Fatal(err)
	}
	digest := sha256.Sum256(payload)
	run.Plan.Digest = hex.EncodeToString(digest[:])
	run.RunID = runIdentity(run)
	store, _ := writeAndOpen(t, run)
	if err := store.IndexRun(context.Background(), run); err != nil {
		t.Fatal(err)
	}
	loaded, err := store.LoadRun(context.Background(), run.RunID)
	if err != nil {
		t.Fatal(err)
	}
	if loaded.SchemaVersion != "prc.run/v0.3" || loaded.Plan.SchemaVersion != "prc.plan/v0.3" || loaded.Plan.EngineVersion != "" {
		t.Fatalf("legacy run was not preserved: %+v", loaded.Plan)
	}
}

func TestStoreReadsLegacyV04RunAfterCatalogBindingUpgrade(t *testing.T) {
	run := testRun(t)
	run.SchemaVersion = "prc.run/v0.4"
	run.Plan.SchemaVersion = "prc.plan/v0.4"
	run.Plan.CatalogDigest = ""
	run.Plan.Digest = ""
	payload, err := json.Marshal(run.Plan)
	if err != nil {
		t.Fatal(err)
	}
	digest := sha256.Sum256(payload)
	run.Plan.Digest = hex.EncodeToString(digest[:])
	run.RunID = runIdentity(run)
	store, _ := writeAndOpen(t, run)
	if err := store.IndexRun(context.Background(), run); err != nil {
		t.Fatal(err)
	}
	loaded, err := store.LoadRun(context.Background(), run.RunID)
	if err != nil {
		t.Fatal(err)
	}
	if loaded.SchemaVersion != "prc.run/v0.4" || loaded.Plan.SchemaVersion != "prc.plan/v0.4" || loaded.Plan.CatalogDigest != "" {
		t.Fatalf("legacy run was not preserved: %+v", loaded.Plan)
	}
}

func TestStoreReadsLegacyV05RunAfterFindingUpgrade(t *testing.T) {
	run := testRun(t)
	run.SchemaVersion = "prc.run/v0.5"
	run.RunID = runIdentity(run)
	store, _ := writeAndOpen(t, run)
	if err := store.IndexRun(context.Background(), run); err != nil {
		t.Fatal(err)
	}
	loaded, err := store.LoadRun(context.Background(), run.RunID)
	if err != nil {
		t.Fatal(err)
	}
	if loaded.SchemaVersion != "prc.run/v0.5" || loaded.Findings != nil {
		t.Fatalf("legacy run was not preserved: %+v", loaded)
	}
}

func TestValidateRunRejectsTamperedNestedIdentities(t *testing.T) {
	t.Run("plan", func(t *testing.T) {
		run := testRun(t)
		run.Plan.ProfileVersion = "tampered"
		run.RunID = runIdentity(run)
		if err := validateRun(run); err == nil || !strings.Contains(err.Error(), "plan digest") {
			t.Fatalf("unexpected error: %v", err)
		}
	})
	t.Run("inventory", func(t *testing.T) {
		run := testRun(t)
		run.Inventory.Files[0].SHA256 = strings.Repeat("0", 64)
		run.RunID = runIdentity(run)
		if err := validateRun(run); err == nil || !strings.Contains(err.Error(), "inventory digest") {
			t.Fatalf("unexpected error: %v", err)
		}
	})
	t.Run("finding-content", func(t *testing.T) {
		run := testRun(t)
		if len(run.Findings) == 0 {
			t.Fatal("fixture produced no findings")
		}
		run.Findings[0].Summary = "tampered"
		run.RunID = runIdentity(run)
		if err := validateRun(run); err == nil || !strings.Contains(err.Error(), "finding") {
			t.Fatalf("unexpected error: %v", err)
		}
	})
	t.Run("missing-finding", func(t *testing.T) {
		run := testRun(t)
		if len(run.Findings) == 0 {
			t.Fatal("fixture produced no findings")
		}
		run.Findings = run.Findings[1:]
		run.RunID = runIdentity(run)
		if err := validateRun(run); err == nil || !strings.Contains(err.Error(), "exactly one finding") {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestStoreReindexIsIdempotentAndRepairsDerivedRows(t *testing.T) {
	ctx := context.Background()
	run := testRun(t)
	store, _ := writeAndOpen(t, run)
	if err := store.IndexRun(ctx, run); err != nil {
		t.Fatal(err)
	}
	if _, err := store.db.ExecContext(ctx, "UPDATE results SET summary = 'tampered' WHERE run_id = ?", run.RunID); err != nil {
		t.Fatal(err)
	}
	if err := store.IndexRun(ctx, run); err != nil {
		t.Fatal(err)
	}
	counts, err := store.Counts(ctx)
	if err != nil || counts.Runs != 1 || counts.Results != len(run.Results) || counts.AuditEvents != 1 {
		t.Fatalf("reindex counts = %+v, %v", counts, err)
	}
	var summary string
	if err := store.db.QueryRowContext(ctx,
		"SELECT summary FROM results WHERE run_id = ? AND assertion_id = ?", run.RunID, run.Results[0].AssertionID).Scan(&summary); err != nil {
		t.Fatal(err)
	}
	if summary != run.Results[0].Summary {
		t.Fatalf("derived row was not repaired: %q", summary)
	}
}

func TestStoreMigratesV1ToFindingIndex(t *testing.T) {
	ctx := context.Background()
	root := t.TempDir()
	if err := os.Chmod(root, 0o700); err != nil {
		t.Fatal(err)
	}
	store, err := Open(ctx, root)
	if err != nil {
		t.Fatal(err)
	}
	for _, table := range []string{"finding_evidence", "finding_locations", "findings"} {
		if _, err := store.db.ExecContext(ctx, "DROP TABLE "+table); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := store.db.ExecContext(ctx, "PRAGMA user_version = 1"); err != nil {
		t.Fatal(err)
	}
	if err := store.Close(); err != nil {
		t.Fatal(err)
	}
	store, err = Open(ctx, root)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	var version int
	if err := store.db.QueryRowContext(ctx, "PRAGMA user_version").Scan(&version); err != nil || version != SchemaVersion {
		t.Fatalf("state schema version=%d err=%v", version, err)
	}
	if err := store.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM findings").Scan(new(int)); err != nil {
		t.Fatalf("finding table unavailable after migration: %v", err)
	}
}

func TestStoreRefusesToIndexWithoutCanonicalRecords(t *testing.T) {
	run := testRun(t)
	root := t.TempDir()
	if err := os.Chmod(root, 0o700); err != nil {
		t.Fatal(err)
	}
	store, err := Open(context.Background(), root)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	if err := store.IndexRun(context.Background(), run); err == nil || !strings.Contains(err.Error(), "indexed record") {
		t.Fatalf("unexpected error: %v", err)
	}
	counts, err := store.Counts(context.Background())
	if err != nil || counts.Runs != 0 {
		t.Fatalf("failed index left state: %+v, %v", counts, err)
	}
}

func TestStoreIndexTransactionRollsBackWhenEvidenceIsMissing(t *testing.T) {
	run := testRun(t)
	store, root := writeAndOpen(t, run)
	var evidenceID string
	for _, result := range run.Results {
		if len(result.EvidenceObserved) > 0 {
			evidenceID = result.EvidenceObserved[0].ID
			break
		}
	}
	if evidenceID == "" {
		t.Fatal("test run produced no evidence")
	}
	if err := os.Remove(filepath.Join(root, "evidence", evidenceID[:2], evidenceID+".json")); err != nil {
		t.Fatal(err)
	}
	if err := store.IndexRun(context.Background(), run); err == nil || !strings.Contains(err.Error(), "evidence") {
		t.Fatalf("unexpected error: %v", err)
	}
	counts, err := store.Counts(context.Background())
	if err != nil || counts.Runs != 0 || counts.Results != 0 {
		t.Fatalf("failed transaction left state: %+v, %v", counts, err)
	}
}

func TestStoreRejectsPublicStateDirectory(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Windows does not expose POSIX directory modes")
	}
	root := t.TempDir()
	if err := os.Chmod(root, 0o755); err != nil {
		t.Fatal(err)
	}
	if _, err := Open(context.Background(), root); err == nil || !strings.Contains(err.Error(), "0700") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLoadRunRejectsRecordSymlinkEscape(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink creation requires additional Windows privileges")
	}
	run := testRun(t)
	store, root := writeAndOpen(t, run)
	if err := store.IndexRun(context.Background(), run); err != nil {
		t.Fatal(err)
	}
	outside := t.TempDir()
	moved := filepath.Join(outside, "runs")
	if err := os.Rename(filepath.Join(root, "runs"), moved); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(moved, filepath.Join(root, "runs")); err != nil {
		t.Fatal(err)
	}
	if _, err := store.LoadRun(context.Background(), run.RunID); err == nil || !strings.Contains(err.Error(), "escapes") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStoreRejectsFutureSchemaVersion(t *testing.T) {
	ctx := context.Background()
	root := t.TempDir()
	if err := os.Chmod(root, 0o700); err != nil {
		t.Fatal(err)
	}
	store, err := Open(ctx, root)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.db.ExecContext(ctx, "PRAGMA user_version = 99"); err != nil {
		t.Fatal(err)
	}
	if err := store.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err := Open(ctx, root); err == nil || !strings.Contains(err.Error(), "newer") {
		t.Fatalf("unexpected error: %v", err)
	}
}
