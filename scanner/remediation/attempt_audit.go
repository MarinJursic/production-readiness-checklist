package remediation

import (
	"encoding/hex"
	"fmt"
)

func validateAttemptAudit(run RemediationRun) error {
	if len(run.Attempts) != run.Usage.Attempts {
		return fmt.Errorf("remediation attempt audit count does not match budget usage")
	}
	candidates := make(map[string]Candidate, len(run.Candidates))
	changedFiles, changedLines := 0, 0
	for _, candidate := range run.Candidates {
		expectedID, err := candidateContentID(candidate)
		if err != nil || expectedID != candidate.CandidateID || !digestIdentifier(candidate.CandidateID) ||
			candidates[candidate.CandidateID].CandidateID != "" {
			return fmt.Errorf("remediation candidate audit identity is invalid or duplicated")
		}
		candidates[candidate.CandidateID] = candidate
		for _, change := range candidate.Changes {
			changedFiles++
			changedLines += change.AddedLines + change.RemovedLines
		}
	}
	if changedFiles != run.Usage.ChangedFiles || changedLines != run.Usage.ChangedLines ||
		run.Usage.Attempts > run.MaxAttempts || changedFiles > run.MaxFiles || changedLines > run.MaxChangedLines {
		return fmt.Errorf("remediation attempt audit does not match cumulative change budgets")
	}
	executions := make(map[string]providerExecutionLink, len(run.ProviderExecutions))
	for _, execution := range run.ProviderExecutions {
		if !digestIdentifier(execution.ExecutionID) || executions[execution.ExecutionID].id != "" {
			return fmt.Errorf("remediation provider execution identity is invalid or duplicated")
		}
		executions[execution.ExecutionID] = providerExecutionLink{id: execution.ExecutionID, taskID: execution.TaskID}
	}
	failures := make(map[string]providerExecutionLink, len(run.ProviderFailures))
	for _, failure := range run.ProviderFailures {
		if failure.Validate() != nil || failures[failure.FailureID].id != "" {
			return fmt.Errorf("remediation provider failure identity is invalid or duplicated")
		}
		failures[failure.FailureID] = providerExecutionLink{
			id: failure.FailureID, taskID: failure.TaskID, reasonCode: failure.ReasonCode, reason: failure.Reason,
		}
	}
	seenCandidates := map[string]bool{}
	seenExecutions := map[string]bool{}
	seenFailures := map[string]bool{}
	expectedBefore := run.SourceInventoryDigest
	for index, record := range run.Attempts {
		if record.Attempt != index+1 || !digestIdentifier(record.FindingID) ||
			!digestIdentifier(record.FindingFingerprint) || !digestIdentifier(record.TaskID) ||
			!digestIdentifier(record.BeforeInventoryDigest) || record.BeforeInventoryDigest != expectedBefore ||
			record.StartedAt.IsZero() || record.CompletedAt.IsZero() || record.CompletedAt.Before(record.StartedAt) ||
			record.Reason == "" {
			return fmt.Errorf("remediation attempt %d has invalid identity, ordering, time, or reason", index+1)
		}
		switch record.Mode {
		case "deterministic":
			if record.ProviderExecutionID != "" || record.ProviderFailureID != "" {
				return fmt.Errorf("deterministic remediation attempt %d links a provider invocation", record.Attempt)
			}
		case "agent":
			execution, executionOK := executions[record.ProviderExecutionID]
			failure, failureOK := failures[record.ProviderFailureID]
			if executionOK == failureOK {
				return fmt.Errorf("agent remediation attempt %d must link exactly one provider result", record.Attempt)
			}
			if executionOK {
				if seenExecutions[record.ProviderExecutionID] || execution.taskID != record.TaskID {
					return fmt.Errorf("agent remediation attempt %d has invalid provider execution linkage", record.Attempt)
				}
				seenExecutions[record.ProviderExecutionID] = true
			} else {
				if seenFailures[record.ProviderFailureID] || failure.taskID != record.TaskID {
					return fmt.Errorf("agent remediation attempt %d has invalid provider failure linkage", record.Attempt)
				}
				seenFailures[record.ProviderFailureID] = true
			}
		default:
			return fmt.Errorf("remediation attempt %d has unsupported mode", record.Attempt)
		}
		switch record.Outcome {
		case "accepted":
			if record.ReasonCode != "accepted" ||
				record.Mode == "agent" && (record.ProviderExecutionID == "" || record.ProviderFailureID != "") {
				return fmt.Errorf("accepted remediation attempt %d has an invalid reason code", record.Attempt)
			}
			candidate, err := linkedAttemptCandidate(record, candidates, seenCandidates)
			if err != nil || !candidate.Accepted || candidateAttemptTaskID(record, candidate) != record.TaskID {
				return fmt.Errorf("accepted remediation attempt %d has invalid candidate linkage: %v", record.Attempt, err)
			}
			expectedBefore = record.AfterInventoryDigest
		case "rejected":
			if record.Mode == "agent" && (record.ProviderExecutionID == "" || record.ProviderFailureID != "") {
				return fmt.Errorf("rejected remediation attempt %d has invalid provider linkage", record.Attempt)
			}
			if record.CandidateID == "" {
				if record.AfterInventoryDigest != "" || !proposalRejectionCode(record.ReasonCode) {
					return fmt.Errorf("pre-candidate rejection %d has invalid audit fields", record.Attempt)
				}
			} else {
				candidate, err := linkedAttemptCandidate(record, candidates, seenCandidates)
				if err != nil || candidate.Accepted || candidateAttemptTaskID(record, candidate) != record.TaskID ||
					record.ReasonCode != "verification_rejected" {
					return fmt.Errorf("rejected remediation attempt %d has invalid candidate linkage: %v", record.Attempt, err)
				}
			}
			if index != len(run.Attempts)-1 {
				return fmt.Errorf("rejected remediation attempt %d is not terminal", record.Attempt)
			}
		case "provider_stopped":
			if record.Mode != "agent" || record.ProviderExecutionID == "" || record.ProviderFailureID != "" ||
				record.CandidateID != "" || record.AfterInventoryDigest != "" ||
				(record.ReasonCode != "provider_unable" && record.ReasonCode != "provider_needs_escalation") ||
				index != len(run.Attempts)-1 {
				return fmt.Errorf("stopped provider attempt %d has invalid audit fields", record.Attempt)
			}
		case "provider_failed":
			failure := failures[record.ProviderFailureID]
			if record.Mode != "agent" || record.ProviderExecutionID != "" || record.ProviderFailureID == "" ||
				record.CandidateID != "" || record.AfterInventoryDigest != "" ||
				!providerFailureReasonCode(record.ReasonCode) || record.ReasonCode != "provider_"+failure.reasonCode ||
				record.Reason != failure.reason || index != len(run.Attempts)-1 {
				return fmt.Errorf("failed provider attempt %d has invalid audit fields", record.Attempt)
			}
		default:
			return fmt.Errorf("remediation attempt %d has unsupported outcome", record.Attempt)
		}
	}
	if expectedBefore != run.FinalInventoryDigest || len(seenCandidates) != len(candidates) ||
		len(seenExecutions) != len(executions) || len(seenFailures) != len(failures) {
		return fmt.Errorf("remediation attempt audit does not cover the final inventory, candidates, and provider invocations")
	}
	return nil
}

func providerFailureReasonCode(code string) bool {
	switch code {
	case "provider_preflight_failed", "provider_transcript_failed", "provider_cancelled", "provider_timeout",
		"provider_output_limit", "provider_process_failed", "provider_workspace_changed",
		"provider_result_missing", "provider_protocol_rejected":
		return true
	default:
		return false
	}
}

type providerExecutionLink struct {
	id         string
	taskID     string
	reasonCode string
	reason     string
}

func linkedAttemptCandidate(record AttemptRecord, candidates map[string]Candidate, seen map[string]bool) (Candidate, error) {
	candidate, ok := candidates[record.CandidateID]
	if !ok {
		return Candidate{}, fmt.Errorf("candidate ID is absent")
	}
	if seen[record.CandidateID] {
		return Candidate{}, fmt.Errorf("candidate ID is duplicated")
	}
	if !digestIdentifier(record.AfterInventoryDigest) || candidate.CandidateInventoryDigest != record.AfterInventoryDigest {
		return Candidate{}, fmt.Errorf("candidate inventory digest does not match")
	}
	if candidate.Contract.Attempt != record.Attempt || candidate.Contract.AssertionID != record.AssertionID {
		return Candidate{}, fmt.Errorf("candidate attempt or assertion does not match")
	}
	contractFindingID := candidate.Contract.FindingID
	if record.Mode == "agent" {
		contractFindingID = candidate.Contract.ProposalFindingID
	}
	if contractFindingID != record.FindingID || candidate.Contract.FindingFingerprint != record.FindingFingerprint {
		return Candidate{}, fmt.Errorf("candidate finding identity does not match")
	}
	seen[record.CandidateID] = true
	return candidate, nil
}

func candidateAttemptTaskID(record AttemptRecord, candidate Candidate) string {
	if record.Mode == "agent" {
		return candidate.Contract.ProposalTaskID
	}
	return candidate.Contract.TaskID
}

func proposalRejectionCode(code string) bool {
	switch code {
	case "anti_gaming_rejected", "scope_rejected", "budget_rejected", "policy_rejected":
		return true
	default:
		return false
	}
}

func digestIdentifier(value string) bool {
	if len(value) != 64 {
		return false
	}
	_, err := hex.DecodeString(value)
	return err == nil
}
