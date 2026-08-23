package engine

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/MarinJursic/production-readiness-checklist/scanner/adapter"
	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
)

type adapterBinding struct {
	AdapterID       string
	ManifestSHA256  string
	ObservationKind string
}

func assertionAdapterBindings(assertion model.Assertion) ([]adapterBinding, error) {
	raw, exists := assertion.Parameters["adapter_bindings"]
	if !exists {
		return []adapterBinding{}, nil
	}
	items, ok := raw.([]any)
	if !ok {
		return nil, fmt.Errorf("assertion %s adapter_bindings must be a list", assertion.ID)
	}
	bindings := make([]adapterBinding, 0, len(items))
	seen := map[string]bool{}
	for _, item := range items {
		mapping, ok := item.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("assertion %s adapter binding must be an object", assertion.ID)
		}
		binding := adapterBinding{}
		binding.AdapterID, _ = mapping["adapter_id"].(string)
		binding.ManifestSHA256, _ = mapping["manifest_sha256"].(string)
		binding.ObservationKind, _ = mapping["observation_kind"].(string)
		decodedDigest, digestErr := hex.DecodeString(binding.ManifestSHA256)
		if strings.TrimSpace(binding.AdapterID) == "" || digestErr != nil || len(decodedDigest) != sha256.Size ||
			binding.ManifestSHA256 != strings.ToLower(binding.ManifestSHA256) || strings.TrimSpace(binding.ObservationKind) == "" {
			return nil, fmt.Errorf("assertion %s has an incomplete adapter binding", assertion.ID)
		}
		key := binding.AdapterID + "\x00" + binding.ManifestSHA256 + "\x00" + binding.ObservationKind
		if seen[key] {
			return nil, fmt.Errorf("assertion %s has a duplicate adapter binding", assertion.ID)
		}
		seen[key] = true
		bindings = append(bindings, binding)
	}
	sort.Slice(bindings, func(i, j int) bool {
		if bindings[i].AdapterID != bindings[j].AdapterID {
			return bindings[i].AdapterID < bindings[j].AdapterID
		}
		if bindings[i].ManifestSHA256 != bindings[j].ManifestSHA256 {
			return bindings[i].ManifestSHA256 < bindings[j].ManifestSHA256
		}
		return bindings[i].ObservationKind < bindings[j].ObservationKind
	})
	return bindings, nil
}

func (e *Engine) AuthorizesAdapter(profileID string, inventory model.Inventory, adapterID, manifestSHA256 string) (bool, error) {
	return e.AuthorizesAdapterMode(profileID, inventory, ExecutionModeVerifyLocal, adapterID, manifestSHA256)
}

func (e *Engine) AuthorizesAdapterMode(
	profileID string,
	inventory model.Inventory,
	mode string,
	adapterID string,
	manifestSHA256 string,
) (bool, error) {
	plan, err := e.PlanMode(profileID, inventory, mode)
	if err != nil {
		return false, err
	}
	return e.authorizesAdapterInPlan(plan, adapterID, manifestSHA256)
}

func (e *Engine) authorizesAdapterInPlan(plan model.Plan, adapterID, manifestSHA256 string) (bool, error) {
	for _, planned := range plan.Adapters {
		if planned.AdapterID == adapterID && planned.ManifestSHA256 == manifestSHA256 {
			return planned.Status == "authorized", nil
		}
	}
	return false, nil
}

func validateAdapterExecutions(inventory model.Inventory, executions []model.AdapterExecution) ([]model.AdapterExecution, error) {
	validated := append([]model.AdapterExecution{}, executions...)
	sort.Slice(validated, func(i, j int) bool { return validated[i].ExecutionID < validated[j].ExecutionID })
	seen := map[string]bool{}
	for _, execution := range validated {
		if err := adapter.ValidateExecution(execution); err != nil {
			return nil, fmt.Errorf("validate adapter execution: %w", err)
		}
		if seen[execution.ExecutionID] {
			return nil, fmt.Errorf("duplicate adapter execution %s", execution.ExecutionID)
		}
		seen[execution.ExecutionID] = true
		if execution.Subject.TargetName != inventory.TargetName ||
			execution.Subject.TargetCommit != inventory.GitCommit ||
			execution.Subject.InventoryDigest != inventory.Digest {
			return nil, fmt.Errorf("adapter execution %s does not match the scanned inventory", execution.ExecutionID)
		}
	}
	return validated, nil
}

func adapterObservationEvidence(
	inventory model.Inventory,
	execution model.AdapterExecution,
	observation model.AdapterObservation,
) model.Evidence {
	payload, _ := json.Marshal(observation)
	evidence := model.Evidence{
		SchemaVersion: model.EvidenceSchema, Kind: "analysis-result", Authority: "executed",
		Producer:      execution.AdapterID,
		TargetDigest:  inventory.Digest,
		Source:        "adapter:" + execution.AdapterID + "/" + execution.AdapterRunID + "/" + observation.ID,
		ContentSHA256: execution.ExecutionID, Size: int64(len(payload)), ObservedAt: execution.CompletedAt,
		Summary: observation.Summary,
	}
	identity, _ := json.Marshal(evidence)
	digest := sha256.Sum256(identity)
	evidence.ID = hex.EncodeToString(digest[:])
	return evidence
}

func evaluateAnalysisEvidence(
	assertion model.Assertion,
	inventory model.Inventory,
	executions []model.AdapterExecution,
	result model.AssertionResult,
) model.AssertionResult {
	bindings, err := assertionAdapterBindings(assertion)
	if err != nil {
		result.Execution, result.Assessment, result.Summary = "error", "unknown", err.Error()
		return result
	}
	if len(bindings) == 0 {
		result.Execution, result.Assessment = "blocked", "unknown"
		result.Summary = "No immutable adapter bindings are configured for this assertion."
		return result
	}
	result.Execution = "completed"
	unresolved := []string{}
	conflicts := []string{}
	failures := []string{}
	for _, binding := range bindings {
		matchedExecution := false
		outcomes := map[string]bool{}
		for _, execution := range executions {
			if execution.AdapterID != binding.AdapterID || execution.ManifestSHA256 != binding.ManifestSHA256 {
				continue
			}
			matchedExecution = true
			if execution.Transcript.Summary.Status != "completed" {
				result.Execution = "blocked"
				unresolved = append(unresolved, binding.ObservationKind+":"+execution.Transcript.Summary.Status)
				continue
			}
			for _, observation := range execution.Transcript.Observations {
				if observation.Kind != binding.ObservationKind {
					continue
				}
				outcomes[observation.Outcome] = true
				result.EvidenceObserved = append(result.EvidenceObserved,
					adapterObservationEvidence(inventory, execution, observation))
			}
		}
		name := binding.AdapterID + "/" + binding.ObservationKind
		if !matchedExecution || len(outcomes) == 0 {
			unresolved = append(unresolved, name+":missing")
			continue
		}
		if outcomes["found"] {
			failures = append(failures, name)
			continue
		}
		if outcomes["not_found"] && len(outcomes) == 1 {
			continue
		}
		if outcomes["not_found"] {
			conflicts = append(conflicts, name)
			continue
		}
		unresolved = append(unresolved, name+":unsupported")
	}
	if len(failures) > 0 {
		result.Assessment = "fail"
		result.Summary = "Authorized executed analysis found violations for: " + strings.Join(failures, ", ") + "."
		return result
	}
	if len(conflicts) > 0 {
		result.Assessment = "conflicting"
		result.Summary = "Authorized analysis observations conflict for: " + strings.Join(conflicts, ", ") + "."
		return result
	}
	if len(unresolved) > 0 {
		result.Assessment = "unknown"
		if result.Execution != "blocked" {
			result.Execution = "completed"
		}
		result.Summary = "Required authorized analysis evidence is unresolved: " + strings.Join(unresolved, ", ") + "."
		return result
	}
	result.Assessment = "pass"
	result.Summary = "Every immutable adapter binding produced completed not-found observations for the scanned inventory."
	return result
}
