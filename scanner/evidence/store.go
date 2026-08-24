package evidence

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/MarinJursic/production-readiness-checklist/scanner/finding"
	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
)

var digestPattern = regexp.MustCompile(`^[0-9a-f]{64}$`)

func evidenceID(item model.Evidence) (string, error) {
	item.ID = ""
	payload, err := json.Marshal(item)
	if err != nil {
		return "", err
	}
	digest := sha256.Sum256(payload)
	return fmt.Sprintf("%x", digest), nil
}

func runID(run model.RunResult) (string, error) {
	run.RunID = ""
	payload, err := json.Marshal(run)
	if err != nil {
		return "", err
	}
	digest := sha256.Sum256(payload)
	return fmt.Sprintf("%x", digest), nil
}

func WriteRun(root string, run model.RunResult) error {
	return WriteRunWithArtifacts(root, run, nil)
}

// WriteRunWithArtifacts stores reviewed native-adapter payloads before the
// immutable records that reference them. The caller supplies payloads by
// sha256: descriptor; every supplied payload must be declared by an adapter
// execution in this run and match its declared size and digest.
func WriteRunWithArtifacts(root string, run model.RunResult, artifacts map[string][]byte) error {
	if root == "" {
		return nil
	}
	if err := validateArtifactPayloads(run, artifacts); err != nil {
		return err
	}
	if !digestPattern.MatchString(run.RunID) {
		return fmt.Errorf("invalid run ID %q", run.RunID)
	}
	expectedRunID, err := runID(run)
	if err != nil {
		return fmt.Errorf("calculate run ID: %w", err)
	}
	if run.RunID != expectedRunID {
		return fmt.Errorf("run ID does not match record content")
	}
	if run.SchemaVersion == model.RunSchema || run.SchemaVersion == "prc.run/v0.11" || run.SchemaVersion == "prc.run/v0.10" || run.SchemaVersion == "prc.run/v0.9" || run.SchemaVersion == "prc.run/v0.8" || run.SchemaVersion == "prc.run/v0.7" || run.SchemaVersion == "prc.run/v0.6" {
		if run.Findings == nil {
			return fmt.Errorf("current run findings must encode as an array")
		}
		for _, item := range run.Findings {
			if err := finding.Validate(item); err != nil {
				return fmt.Errorf("invalid finding %s: %w", item.ID, err)
			}
		}
	}
	for _, result := range run.Results {
		for _, item := range result.EvidenceObserved {
			if !digestPattern.MatchString(item.ID) {
				return fmt.Errorf("invalid evidence ID %q", item.ID)
			}
			expectedEvidenceID, err := evidenceID(item)
			if err != nil {
				return fmt.Errorf("calculate evidence ID: %w", err)
			}
			if item.ID != expectedEvidenceID {
				return fmt.Errorf("evidence ID does not match record content")
			}
			path := filepath.Join(root, "evidence", item.ID[:2], item.ID+".json")
			if err := writeJSONImmutable(path, item); err != nil {
				return fmt.Errorf("write evidence %s: %w", item.ID, err)
			}
		}
	}
	for descriptor, payload := range artifacts {
		digest := descriptor[len("sha256:"):]
		path := filepath.Join(root, "artifacts", "sha256", digest[:2], digest)
		if err := writeBytesImmutable(path, payload); err != nil {
			return fmt.Errorf("write artifact %s: %w", descriptor, err)
		}
	}
	path := filepath.Join(root, "runs", run.RunID+".json")
	if err := writeJSONImmutable(path, run); err != nil {
		return fmt.Errorf("write run %s: %w", run.RunID, err)
	}
	return nil
}

func validateArtifactPayloads(run model.RunResult, payloads map[string][]byte) error {
	declared := map[string]map[int64]bool{}
	for _, execution := range run.AdapterExecutions {
		for _, artifact := range execution.Transcript.Artifacts {
			digest := strings.TrimPrefix(artifact.Digest, "sha256:")
			if !strings.HasPrefix(artifact.Digest, "sha256:") || !digestPattern.MatchString(digest) {
				return fmt.Errorf("adapter execution declares invalid artifact digest %q", artifact.Digest)
			}
			if artifact.Size < 0 {
				return fmt.Errorf("adapter execution declares invalid artifact size")
			}
			if declared[artifact.Digest] == nil {
				declared[artifact.Digest] = map[int64]bool{}
			}
			declared[artifact.Digest][artifact.Size] = true
		}
	}
	for descriptor, payload := range payloads {
		digest := strings.TrimPrefix(descriptor, "sha256:")
		if !strings.HasPrefix(descriptor, "sha256:") || !digestPattern.MatchString(digest) {
			return fmt.Errorf("invalid artifact payload descriptor %q", descriptor)
		}
		if !declared[descriptor][int64(len(payload))] {
			return fmt.Errorf("artifact payload %s is undeclared or has the wrong size", descriptor)
		}
		actual := sha256.Sum256(payload)
		if fmt.Sprintf("%x", actual) != digest {
			return fmt.Errorf("artifact payload %s does not match its digest", descriptor)
		}
	}
	return nil
}

func writeJSONImmutable(path string, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return writeBytesImmutable(path, data)
}

func writeBytesImmutable(path string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	temporary, err := os.CreateTemp(filepath.Dir(path), ".prc-*.tmp")
	if err != nil {
		return err
	}
	temporaryName := temporary.Name()
	defer os.Remove(temporaryName)
	if err := temporary.Chmod(0o600); err != nil {
		temporary.Close()
		return err
	}
	if _, err := temporary.Write(data); err != nil {
		temporary.Close()
		return err
	}
	if err := temporary.Sync(); err != nil {
		temporary.Close()
		return err
	}
	if err := temporary.Close(); err != nil {
		return err
	}
	if err := os.Link(temporaryName, path); err == nil {
		return nil
	} else if !os.IsExist(err) {
		return err
	}
	existing, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if !bytes.Equal(existing, data) {
		return fmt.Errorf("immutable record already exists with different content")
	}
	return nil
}
