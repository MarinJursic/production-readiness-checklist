package evidence

import (
	"crypto/sha256"
	"fmt"
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
		SchemaVersion: model.RunSchema, Findings: []model.Finding{},
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

func TestWriteRunStoresDeclaredArtifactPayloadByDigest(t *testing.T) {
	root := t.TempDir()
	payload := []byte("{\"bomFormat\":\"CycloneDX\"}\n")
	digest := fmt.Sprintf("%x", sha256.Sum256(payload))
	run := signRun(t, model.RunResult{
		SchemaVersion: model.RunSchema, Findings: []model.Finding{},
		AdapterExecutions: []model.AdapterExecution{{
			Transcript: model.AdapterTranscript{Artifacts: []model.AdapterArtifact{{
				ID: "sbom", Digest: "sha256:" + digest, Size: int64(len(payload)),
			}}},
		}},
	})
	if err := WriteRunWithArtifacts(root, run, map[string][]byte{"sha256:" + digest: payload}); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(root, "artifacts", "sha256", digest[:2], digest)
	stored, err := os.ReadFile(path)
	if err != nil || string(stored) != string(payload) {
		t.Fatalf("stored artifact = %q, %v", stored, err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Fatalf("artifact mode = %v, %v", info.Mode().Perm(), err)
	}
}

func TestWriteRunRejectsUndeclaredOrMismatchedArtifactPayload(t *testing.T) {
	payload := []byte("sbom")
	digest := fmt.Sprintf("%x", sha256.Sum256(payload))
	run := signRun(t, model.RunResult{SchemaVersion: model.RunSchema, Findings: []model.Finding{}})
	if err := WriteRunWithArtifacts(t.TempDir(), run, map[string][]byte{"sha256:" + digest: payload}); err == nil {
		t.Fatal("expected undeclared artifact rejection")
	}
	run = signRun(t, model.RunResult{
		SchemaVersion: model.RunSchema, Findings: []model.Finding{},
		AdapterExecutions: []model.AdapterExecution{{Transcript: model.AdapterTranscript{
			Artifacts: []model.AdapterArtifact{{Digest: "sha256:" + digest, Size: int64(len(payload))}},
		}}},
	})
	if err := WriteRunWithArtifacts(t.TempDir(), run, map[string][]byte{"sha256:" + digest: []byte("tampered")}); err == nil {
		t.Fatal("expected artifact payload mismatch rejection")
	}
}

func TestWriteRunRejectsImmutableArtifactCollision(t *testing.T) {
	root := t.TempDir()
	payload := []byte("sbom")
	digest := fmt.Sprintf("%x", sha256.Sum256(payload))
	path := filepath.Join(root, "artifacts", "sha256", digest[:2], digest)
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("different"), 0o600); err != nil {
		t.Fatal(err)
	}
	run := signRun(t, model.RunResult{
		SchemaVersion: model.RunSchema, Findings: []model.Finding{},
		AdapterExecutions: []model.AdapterExecution{{Transcript: model.AdapterTranscript{
			Artifacts: []model.AdapterArtifact{{Digest: "sha256:" + digest, Size: int64(len(payload))}},
		}}},
	})
	if err := WriteRunWithArtifacts(root, run, map[string][]byte{"sha256:" + digest: payload}); err == nil {
		t.Fatal("expected immutable artifact collision")
	}
}

func TestWriteRunWithArtifactsIsNoOpWithoutStateRoot(t *testing.T) {
	if err := WriteRunWithArtifacts("", model.RunResult{}, map[string][]byte{"invalid": []byte("ignored")}); err != nil {
		t.Fatal(err)
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
