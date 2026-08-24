package controlreview

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/MarinJursic/production-readiness-checklist/scanner/provider"
)

func sealTask(task Task) (Task, error) {
	task.TaskID = ""
	id, err := taskID(task)
	if err != nil {
		return Task{}, err
	}
	task.TaskID = id
	if err := validateTask(task); err != nil {
		return Task{}, err
	}
	return task, nil
}

func taskID(task Task) (string, error) {
	task.TaskID = ""
	payload, err := json.Marshal(task)
	if err != nil {
		return "", fmt.Errorf("encode control-review task identity: %w", err)
	}
	digest := sha256.Sum256(payload)
	return hex.EncodeToString(digest[:]), nil
}

func validateTask(task Task) error {
	if task.SchemaVersion != TaskSchema || task.Provider != "codex" && task.Provider != "claude" ||
		!task.RequireOneSubagentPerRule || len(task.Controls) < 1 || len(task.Controls) > 8 ||
		!lowerHexDigest(task.InventoryDigest) || !lowerHexDigest(task.RegistrySHA256) {
		return fmt.Errorf("invalid control-review task envelope")
	}
	expected, err := taskID(task)
	if err != nil || task.TaskID != expected {
		return fmt.Errorf("control-review task ID does not match canonical content")
	}
	seen := map[string]bool{}
	validEvaluation := map[string]bool{"repository": true, "environment": true, "human_external": true, "mixed": true, "unclassified": true}
	validAutomation := map[string]bool{
		"deterministic_candidate": true, "ai_advisory_candidate": true, "environment_evidence_required": true,
		"human_or_external_required": true, "mixed_evidence_required": true,
	}
	validAuthority := map[string]bool{"declared": true, "repository": true, "artifact": true, "environment": true, "human": true}
	for index, control := range task.Controls {
		if control.ControlID == "" || seen[control.ControlID] || strings.TrimSpace(control.Statement) == "" ||
			len(control.Statement) > 16*1024 || !lowerHexDigest(control.ContractSHA256) ||
			(control.ContractStatus != "generated_unreviewed" && control.ContractStatus != "reviewed") ||
			control.CanonicalControlID == "" || control.CanonicalControlID > control.ControlID ||
			!validEvaluation[control.EvaluationClass] || !validAutomation[control.AutomationClass] ||
			(control.ApplicabilityClass != "conditional" && control.ApplicabilityClass != "scope_required") ||
			(control.Atomicity != "apparently_atomic" && control.Atomicity != "compound_review_required") ||
			len(control.EvidenceAuthorities) == 0 || strings.TrimSpace(control.NotApplicableProof) == "" ||
			len(control.NotApplicableProof) > 1000 || len(control.CurrentAssertionChecks) > 256 {
			return fmt.Errorf("invalid control-review task control %q", control.ControlID)
		}
		if index > 0 && task.Controls[index-1].ControlID >= control.ControlID {
			return fmt.Errorf("control-review task controls must be strictly ordered")
		}
		seen[control.ControlID] = true
		seenAuthorities := map[string]bool{}
		for _, authority := range control.EvidenceAuthorities {
			if !validAuthority[authority] || seenAuthorities[authority] {
				return fmt.Errorf("control-review task %s has an invalid evidence authority", control.ControlID)
			}
			seenAuthorities[authority] = true
		}
	}
	seenPaths := map[string]bool{}
	pathBytes := 0
	for index, path := range task.RepositoryPaths {
		if err := safeReviewPath(path); err != nil || seenPaths[path] || index > 0 && task.RepositoryPaths[index-1] >= path {
			return fmt.Errorf("invalid, duplicate, or unordered repository path context")
		}
		seenPaths[path] = true
		pathBytes += len(path) + 1
	}
	if pathBytes > maximumPathContextBytes {
		return fmt.Errorf("repository path context exceeds its byte limit")
	}
	contextBytes := 0
	seenContext := map[string]bool{}
	for _, input := range task.ContextFiles {
		if err := safeReviewPath(input.Path); err != nil || seenContext[input.Path] || input.StartLine < 1 ||
			input.EndLine < input.StartLine || !seenPaths[input.Path] || input.Content == "" ||
			len(input.Content) > maximumContextFileBytes || !utf8.ValidString(input.Content) ||
			input.EndLine-input.StartLine+1 != visibleLineCount(input.Content) {
			return fmt.Errorf("invalid or duplicate review context file %q", input.Path)
		}
		digest := sha256.Sum256([]byte(input.Content))
		if input.SHA256 != hex.EncodeToString(digest[:]) {
			return fmt.Errorf("review context file %s has an invalid content digest", input.Path)
		}
		seenContext[input.Path] = true
		contextBytes += len(input.Content)
	}
	if contextBytes > maximumContextTotal || len(task.SnapshotLimitations) > 64 {
		return fmt.Errorf("review context exceeds its total limit")
	}
	return nil
}

func visibleLineCount(content string) int {
	lines := strings.Count(content, "\n")
	if !strings.HasSuffix(content, "\n") {
		lines++
	}
	return lines
}

func safeReviewPath(path string) error {
	if path == "" || filepath.IsAbs(path) || strings.Contains(path, "\\") ||
		filepath.ToSlash(filepath.Clean(filepath.FromSlash(path))) != path || strings.HasPrefix(path, "../") {
		return fmt.Errorf("unsafe review path %q", path)
	}
	return nil
}

func parseOutput(providerName string, data []byte, task Task) (Output, error) {
	if providerName == "claude" {
		var envelope map[string]json.RawMessage
		if err := provider.DecodeStrictJSON(data, &envelope); err != nil {
			return Output{}, fmt.Errorf("decode Claude review envelope: %w", err)
		}
		var isError bool
		if raw := envelope["is_error"]; len(raw) > 0 && json.Unmarshal(raw, &isError) != nil {
			return Output{}, fmt.Errorf("claude review envelope has invalid is_error")
		}
		structured := envelope["structured_output"]
		if isError || len(structured) == 0 || string(structured) == "null" {
			return Output{}, fmt.Errorf("claude did not return a successful structured review")
		}
		data = structured
	} else if providerName != "codex" {
		return Output{}, fmt.Errorf("unsupported control-review provider %q", providerName)
	}
	var output Output
	if err := decodeInnerOutput(data, &output); err != nil {
		return Output{}, fmt.Errorf("decode control-review output: %w", err)
	}
	if err := validateOutput(output, task); err != nil {
		return Output{}, err
	}
	return output, nil
}

func decodeInnerOutput(data []byte, output *Output) error {
	return provider.DecodeStrictJSON(data, output)
}

func validateOutput(output Output, task Task) error {
	if output.SchemaVersion != OutputSchema {
		return fmt.Errorf("control-review output uses unsupported schema %q", output.SchemaVersion)
	}
	if output.TaskID != task.TaskID {
		return fmt.Errorf("control-review output task ID does not match its sealed task")
	}
	if len(output.Reviews) != len(task.Controls) {
		return fmt.Errorf("control-review output has %d reviews for %d controls", len(output.Reviews), len(task.Controls))
	}
	for index, review := range output.Reviews {
		if review.ControlID != task.Controls[index].ControlID {
			return fmt.Errorf("control-review output is missing, duplicate, or out of order for %s", task.Controls[index].ControlID)
		}
		if err := validateReview(review, task.ContextFiles); err != nil {
			return fmt.Errorf("invalid advisory review for %s: %w", review.ControlID, err)
		}
	}
	return nil
}

func validateReview(review Review, contextFiles []ContextFile) error {
	validAssessment := map[string]bool{
		"advisory_pass_candidate": true, "advisory_fail_candidate": true,
		"needs_evidence": true, "not_applicable_candidate": true,
	}
	validApplicability := map[string]bool{"applicable": true, "not_applicable": true, "undetermined": true}
	validConfidence := map[string]bool{"low": true, "medium": true, "high": true}
	if !validAssessment[review.AssessmentCandidate] || !validApplicability[review.ApplicabilityCandidate] ||
		!validConfidence[review.Confidence] || strings.TrimSpace(review.Reason) == "" ||
		len(review.Reason) > 16*1024 || len(review.Advice) > 16*1024 || len(review.Evidence) > 256 ||
		len(review.Limitations) < 1 || len(review.Limitations) > 256 {
		return fmt.Errorf("invalid candidate values or field limits")
	}
	if review.Confidence == "high" && len(review.Evidence) == 0 {
		return fmt.Errorf("high confidence requires cited repository evidence")
	}
	seenEvidence := map[string]bool{}
	visible := map[string]ContextFile{}
	for _, contextFile := range contextFiles {
		visible[contextFile.Path] = contextFile
	}
	for _, evidence := range review.Evidence {
		clean := filepath.ToSlash(filepath.Clean(filepath.FromSlash(evidence.Path)))
		excerpt, exists := visible[evidence.Path]
		key := fmt.Sprintf("%s:%d:%d", evidence.Path, evidence.Line, evidence.Column)
		if evidence.Path == "" || filepath.IsAbs(evidence.Path) || clean != evidence.Path ||
			!exists || evidence.Line < excerpt.StartLine || evidence.Line > excerpt.EndLine ||
			evidence.Column < 0 || seenEvidence[key] {
			return fmt.Errorf("unsafe, unseen, or duplicate evidence location %q", key)
		}
		seenEvidence[key] = true
	}
	limitations := append([]string{}, review.Limitations...)
	sort.Strings(limitations)
	for index, limitation := range limitations {
		if strings.TrimSpace(limitation) == "" || len(limitation) > 16*1024 || index > 0 && limitations[index-1] == limitation {
			return fmt.Errorf("invalid or duplicate limitation")
		}
	}
	return nil
}

func lowerHexDigest(value string) bool {
	if len(value) != 64 {
		return false
	}
	for _, character := range value {
		if (character < '0' || character > '9') && (character < 'a' || character > 'f') {
			return false
		}
	}
	return true
}
