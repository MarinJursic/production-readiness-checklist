package evidence

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

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
	if root == "" {
		return nil
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
	if run.SchemaVersion == model.RunSchema || run.SchemaVersion == "prc.run/v0.6" {
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
	path := filepath.Join(root, "runs", run.RunID+".json")
	if err := writeJSONImmutable(path, run); err != nil {
		return fmt.Errorf("write run %s: %w", run.RunID, err)
	}
	return nil
}

func writeJSONImmutable(path string, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
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
