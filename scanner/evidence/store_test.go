package evidence

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
)

func signRun(t *testing.T, run model.RunResult) model.RunResult {
	t.Helper()
	for resultIndex := range run.Results {
		for evidenceIndex := range run.Results[resultIndex].EvidenceObserved {
			item := &run.Results[resultIndex].EvidenceObserved[evidenceIndex]
			id, err := evidenceID(*item)
			if err != nil {
				t.Fatal(err)
			}
			item.ID = id
		}
	}
	id, err := runID(run)
	if err != nil {
		t.Fatal(err)
	}
	run.RunID = id
	return run
}

func TestWriteRunStoresEvidenceAndRunAtomically(t *testing.T) {
	root := t.TempDir()
	run := signRun(t, model.RunResult{
		SchemaVersion: model.RunSchema,
		Results: []model.AssertionResult{{
			AssertionID:      "PRC-A-CORE-001",
			EvidenceObserved: []model.Evidence{{SchemaVersion: model.EvidenceSchema}},
		}},
	})
	evidenceID := run.Results[0].EvidenceObserved[0].ID
	runID := run.RunID
	if err := WriteRun(root, run); err != nil {
		t.Fatal(err)
	}
	paths := []string{
		filepath.Join(root, "evidence", evidenceID[:2], evidenceID+".json"),
		filepath.Join(root, "runs", runID+".json"),
	}
	for _, path := range paths {
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("expected %s: %v", path, err)
		}
	}
	matches, err := filepath.Glob(filepath.Join(root, "**", ".prc-*.tmp"))
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) != 0 {
		t.Fatalf("temporary files remain: %v", matches)
	}
}

func TestWriteRunRejectsMalformedEvidenceID(t *testing.T) {
	run := model.RunResult{
		Results: []model.AssertionResult{{EvidenceObserved: []model.Evidence{{ID: "x"}}}},
	}
	id, err := runID(run)
	if err != nil {
		t.Fatal(err)
	}
	run.RunID = id
	if err := WriteRun(t.TempDir(), run); err == nil {
		t.Fatal("expected malformed evidence ID error")
	}
}

func TestWriteRunRejectsMalformedRunID(t *testing.T) {
	if err := WriteRun(t.TempDir(), model.RunResult{RunID: "../escape"}); err == nil {
		t.Fatal("expected malformed run ID error")
	}
}

func TestImmutableWriterDoesNotOverwriteDifferentContent(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "record.json")
	if err := writeJSONImmutable(path, map[string]string{"value": "first"}); err != nil {
		t.Fatal(err)
	}
	if err := writeJSONImmutable(path, map[string]string{"value": "different"}); err == nil {
		t.Fatal("expected immutable record collision")
	}
}
