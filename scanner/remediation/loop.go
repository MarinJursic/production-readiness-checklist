package remediation

import (
	"context"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/MarinJursic/production-readiness-checklist/scanner/catalog"
	"github.com/MarinJursic/production-readiness-checklist/scanner/engine"
	"github.com/MarinJursic/production-readiness-checklist/scanner/inventory"
	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
	"github.com/MarinJursic/production-readiness-checklist/scanner/provider"
	"github.com/MarinJursic/production-readiness-checklist/scanner/verifier"
)

// RunLoop closes registered deterministic R1 findings and, when explicitly
// configured, scanner-planned R2 suggestions. Every accepted attempt becomes
// the immutable source for a fresh sibling candidate; the original workspace
// is never modified.
func RunLoop(options LoopOptions) (RemediationRun, error) {
	if options.ProfileID == "" {
		options.ProfileID = "prc/core-repository"
	}
	if options.Now == nil {
		options.Now = func() time.Time { return time.Now().UTC() }
	}
	if options.Context == nil {
		options.Context = context.Background()
	}
	started := options.Now().UTC()
	policy, err := resolvePolicy(options.Target, options.ProfileID, options.MaxFiles, options.MaxChangedLines,
		1, options.MaxAttempts, options.Configuration)
	if err != nil {
		return RemediationRun{}, err
	}
	maxDurationSeconds, err := resolveLoopDuration(options.MaxDurationSeconds, options.Configuration)
	if err != nil {
		return RemediationRun{}, err
	}
	runContext, cancelRun := context.WithTimeout(options.Context, time.Duration(maxDurationSeconds)*time.Second)
	defer cancelRun()
	c, err := catalog.Load(options.CatalogRoot)
	if err != nil {
		return RemediationRun{}, err
	}
	if options.Agent != nil {
		if err := validateAgentOptions(*options.Agent); err != nil {
			return RemediationRun{}, err
		}
	}
	scanner := engine.New(c)
	original, err := inventory.Build(options.Target)
	if err != nil {
		return RemediationRun{}, err
	}
	original, err = policy.bind(original, "")
	if err != nil {
		return RemediationRun{}, err
	}
	candidateRoot, err := resolveCandidateContainer(original.Root, options.CandidateRoot)
	if err != nil {
		return RemediationRun{}, err
	}
	candidateContainerPrepared := false
	activeTarget := original.Root
	activeConfiguration := rebasedConfiguration(options.Configuration, policy.configRelative, activeTarget)
	activeInventory := original
	activeRun, err := scanner.Scan(options.ProfileID, activeInventory)
	if err != nil {
		return RemediationRun{}, err
	}
	candidates := []Candidate{}
	attemptRecords := []AttemptRecord{}
	providerExecutions := []provider.Execution{}
	providerFailures := []provider.Failure{}
	usage := BudgetUsage{}
	stopReasons := []string{}
	stoppedAssertion, stoppedCode := "", ""

	for {
		assertionID := nextDeterministicFailure(activeRun, c)
		mode := "deterministic"
		var agentTask provider.Task
		if assertionID == "" && options.Agent != nil {
			assertionID, agentTask, err = nextAgentFailure(activeRun, activeInventory, c, policy.protectedPaths, *options.Agent)
			if err != nil {
				return RemediationRun{}, err
			}
			mode = "agent"
		}
		if assertionID == "" {
			break
		}
		if remediationDeadlineReached(runContext) {
			stoppedAssertion, stoppedCode = assertionID, "time_budget_exhausted"
			stopReasons = append(stopReasons, "The configured remediation duration budget is exhausted.")
			break
		}
		if usage.Attempts >= policy.maxAttempts {
			stoppedAssertion, stoppedCode = assertionID, "budget_exhausted"
			stopReasons = append(stopReasons, "The configured remediation attempt budget is exhausted.")
			break
		}
		remainingFiles := policy.maxFiles - usage.ChangedFiles
		remainingLines := policy.maxChangedLines - usage.ChangedLines
		requiredFiles, requiredLines := 1, 1
		if mode == "deterministic" {
			requiredFiles, requiredLines, err = deterministicCost(activeInventory, assertionID)
			if err != nil {
				return RemediationRun{}, err
			}
		}
		if requiredFiles > remainingFiles || requiredLines > remainingLines {
			stoppedAssertion, stoppedCode = assertionID, "budget_exhausted"
			stopReasons = append(stopReasons,
				fmt.Sprintf("Assertion %s requires %d files and %d changed lines, exceeding the remaining budget of %d files and %d lines.",
					assertionID, requiredFiles, requiredLines, remainingFiles, remainingLines))
			break
		}
		if mode == "agent" {
			if options.Verifier == nil {
				return RemediationRun{}, policyDenied(fmt.Errorf("agent remediation requires an independent sandbox verifier"))
			}
			kind, kindErr := verifier.InferKind(agentTask.RelevantPaths[0])
			if kindErr != nil {
				return RemediationRun{}, policyDenied(kindErr)
			}
			verificationOptions := *options.Verifier
			verificationOptions.Kind = kind
			if err := verifier.ValidateOptions(verificationOptions); err != nil {
				return RemediationRun{}, policyDenied(err)
			}
			if err := verifier.Preflight(runContext, verificationOptions); err != nil {
				if runContext.Err() == context.DeadlineExceeded {
					stoppedAssertion, stoppedCode = assertionID, "time_budget_exhausted"
					stopReasons = append(stopReasons, "The configured remediation duration budget expired during verifier preflight.")
					break
				}
				return RemediationRun{}, policyDenied(err)
			}
		}
		if !candidateContainerPrepared {
			candidateRoot, err = prepareCandidateContainer(original.Root, candidateRoot)
			if err != nil {
				return RemediationRun{}, err
			}
			candidateContainerPrepared = true
		}
		attempt := usage.Attempts + 1
		finding, ok := findingFor(activeRun, assertionID)
		if !ok {
			return RemediationRun{}, fmt.Errorf("failed assertion %s has no canonical finding", assertionID)
		}
		attemptRecord := AttemptRecord{
			Attempt: attempt, Mode: mode, AssertionID: assertionID,
			FindingID: finding.ID, FindingFingerprint: finding.Fingerprint,
			StartedAt: options.Now().UTC(), BeforeInventoryDigest: activeInventory.Digest,
		}
		destination := filepath.Join(candidateRoot, fmt.Sprintf("attempt-%03d", attempt))
		lineAllowance := remainingLines
		if lineAllowance < 1 {
			lineAllowance = 1
		}
		var candidate Candidate
		if mode == "deterministic" {
			candidate, err = Run(Options{
				CatalogRoot: options.CatalogRoot, Target: activeTarget, CandidateDir: destination,
				ProfileID: options.ProfileID, TargetName: original.TargetName, AssertionID: assertionID,
				MaxFiles: remainingFiles, MaxChangedLines: lineAllowance,
				Attempt: attempt, MaxAttempts: policy.maxAttempts, Configuration: activeConfiguration,
			})
			if err != nil {
				return RemediationRun{}, err
			}
			attemptRecord.TaskID = candidate.Contract.TaskID
			attemptRecord.FindingID = candidate.Contract.FindingID
			attemptRecord.FindingFingerprint = candidate.Contract.FindingFingerprint
		} else {
			attemptRecord.TaskID = agentTask.TaskID
			outputDirectory := filepath.Join(candidateRoot, fmt.Sprintf("attempt-%03d-provider", attempt))
			if err := os.Mkdir(outputDirectory, 0o700); err != nil {
				return RemediationRun{}, fmt.Errorf("create provider output directory: %w", err)
			}
			if err := writeSealedAgentTask(outputDirectory, agentTask); err != nil {
				return RemediationRun{}, err
			}
			launchPlan, err := provider.BuildPlan(options.Agent.Provider, options.Agent.Executable, activeTarget,
				outputDirectory, options.Agent.OutputSchemaPath, agentTask)
			if err != nil {
				return RemediationRun{}, providerExecution(err)
			}
			execution, executionErr := provider.Run(runContext, launchPlan, agentTask)
			if err := verifySealedAgentTask(outputDirectory, agentTask); err != nil {
				return RemediationRun{}, providerExecution(err)
			}
			if executionErr != nil {
				failure, ok := provider.FailureFromError(executionErr)
				if !ok {
					return RemediationRun{}, providerExecution(executionErr)
				}
				providerFailures = append(providerFailures, failure)
				attemptRecord.ProviderFailureID = failure.FailureID
				attemptRecord.CompletedAt = options.Now().UTC()
				attemptRecord.Outcome = "provider_failed"
				attemptRecord.ReasonCode = "provider_" + failure.ReasonCode
				attemptRecord.Reason = failure.Reason
				attemptRecords = append(attemptRecords, attemptRecord)
				usage.Attempts++
				if runContext.Err() == context.DeadlineExceeded {
					stoppedAssertion, stoppedCode = assertionID, "time_budget_exhausted"
					stopReasons = append(stopReasons, failure.Reason,
						"The configured remediation duration budget expired during provider execution.")
				} else {
					stoppedAssertion, stoppedCode = assertionID, "provider_failed"
					stopReasons = append(stopReasons, failure.Reason)
				}
				break
			}
			providerExecutions = append(providerExecutions, execution)
			attemptRecord.ProviderExecutionID = execution.ExecutionID
			usage.Attempts++
			if execution.Output.Status != "candidate" {
				attemptRecord.CompletedAt = options.Now().UTC()
				attemptRecord.Outcome = "provider_stopped"
				attemptRecord.ReasonCode = "provider_" + execution.Output.Status
				attemptRecord.Reason = "Provider returned " + execution.Output.Status + " without a candidate patch."
				attemptRecords = append(attemptRecords, attemptRecord)
				stoppedAssertion, stoppedCode = assertionID, "provider_stopped"
				stopReasons = append(stopReasons, "The provider returned "+execution.Output.Status+"; the loop stopped without applying a patch.")
				break
			}
			candidate, err = RunProposal(ProposalOptions{
				CatalogRoot: options.CatalogRoot, Target: activeTarget, CandidateDir: destination,
				ProfileID: options.ProfileID, TargetName: original.TargetName,
				Provider: options.Agent.Provider, Task: agentTask, Output: execution.Output,
				MaxFiles: remainingFiles, MaxChangedLines: lineAllowance,
				Attempt: attempt, MaxAttempts: policy.maxAttempts, Configuration: activeConfiguration,
				Verifier: options.Verifier, Context: runContext,
			})
			if err != nil {
				if IsPolicyDenied(err) {
					attemptRecord.CompletedAt = options.Now().UTC()
					attemptRecord.Outcome = "rejected"
					attemptRecord.ReasonCode = proposalPolicyReasonCode(err)
					attemptRecord.Reason = err.Error()
					attemptRecords = append(attemptRecords, attemptRecord)
					if runContext.Err() == context.DeadlineExceeded {
						stoppedAssertion, stoppedCode = assertionID, "time_budget_exhausted"
						stopReasons = append(stopReasons,
							"The configured remediation duration budget expired while evaluating the provider proposal.")
					} else {
						stoppedAssertion, stoppedCode = assertionID, "candidate_rejected"
						stopReasons = append(stopReasons, "Provider proposal rejected by scanner policy: "+err.Error())
					}
					break
				}
				return RemediationRun{}, err
			}
		}
		deadlineRejected := false
		candidate, deadlineRejected, err = enforceCandidateDeadline(runContext, candidate)
		if err != nil {
			return RemediationRun{}, err
		}
		candidates = append(candidates, candidate)
		attemptRecord.CompletedAt = options.Now().UTC()
		attemptRecord.CandidateID = candidate.CandidateID
		attemptRecord.AfterInventoryDigest = candidate.CandidateInventoryDigest
		if candidate.Verification != nil {
			attemptRecord.VerificationExecutionID = candidate.Verification.ExecutionID
		}
		if mode == "deterministic" {
			usage.Attempts++
		}
		for _, change := range candidate.Changes {
			usage.ChangedFiles++
			usage.ChangedLines += change.AddedLines + change.RemovedLines
		}
		if !candidate.Accepted {
			attemptRecord.Outcome = "rejected"
			attemptRecord.ReasonCode = "verification_rejected"
			if deadlineRejected {
				attemptRecord.ReasonCode = "budget_rejected"
			}
			attemptRecord.Reason = candidateRejectionReason(candidate)
			attemptRecords = append(attemptRecords, attemptRecord)
			if runContext.Err() == context.DeadlineExceeded {
				stoppedAssertion, stoppedCode = assertionID, "time_budget_exhausted"
				stopReasons = append(stopReasons,
					"The configured remediation duration budget expired during candidate verification.")
			} else {
				stoppedAssertion, stoppedCode = assertionID, "candidate_rejected"
				stopReasons = append(stopReasons,
					"The independently verified candidate was rejected; the loop did not retry or advance its source workspace.")
			}
			break
		}
		attemptRecord.Outcome = "accepted"
		attemptRecord.ReasonCode = "accepted"
		attemptRecord.Reason = "Candidate passed independent scanner verification."
		attemptRecords = append(attemptRecords, attemptRecord)
		activeTarget = candidate.CandidatePath
		activeConfiguration = rebasedConfiguration(options.Configuration, policy.configRelative, activeTarget)
		activeInventory, err = inventory.Build(activeTarget)
		if err != nil {
			return RemediationRun{}, err
		}
		activeInventory.TargetName = original.TargetName
		if activeConfiguration != nil {
			activeInventory, err = inventory.BindConfiguration(activeInventory, activeConfiguration.Validation, activeConfiguration.SourcePath)
			if err != nil {
				return RemediationRun{}, err
			}
		}
		activeRun, err = scanner.Scan(options.ProfileID, activeInventory)
		if err != nil {
			return RemediationRun{}, err
		}
	}

	currentOriginal, err := inventory.Build(options.Target)
	if err == nil {
		currentOriginal, err = policy.bind(currentOriginal, "")
	}
	if err != nil || currentOriginal.Root != original.Root || currentOriginal.Digest != original.Digest {
		return RemediationRun{}, fmt.Errorf("original workspace changed during remediation loop")
	}
	remaining := classifyRemaining(activeRun, c, stoppedAssertion, stoppedCode, options.Agent != nil)
	terminalState := "machine_work_complete"
	if stoppedCode == "budget_exhausted" || stoppedCode == "time_budget_exhausted" {
		terminalState = "stopped_by_policy_or_budget"
	} else if stoppedCode == "candidate_rejected" {
		terminalState = "candidate_rejected"
	} else if stoppedCode == "provider_stopped" {
		terminalState = "provider_stopped"
	} else if stoppedCode == "provider_failed" {
		terminalState = "provider_failed"
	} else if activeRun.TerminalState == "profile_satisfied" {
		terminalState = "profile_satisfied"
	}
	result := RemediationRun{
		SchemaVersion: RunSchema, StartedAt: started, CompletedAt: options.Now().UTC(),
		ProfileID: options.ProfileID, ConfigurationDigest: policy.configurationID, ProjectID: policy.projectID,
		SourceInventoryDigest: original.Digest, CandidateRoot: candidateRoot, ResultWorkspace: activeTarget,
		FinalInventoryDigest: activeInventory.Digest, OriginalUnchanged: true,
		MaxAttempts: policy.maxAttempts, MaxFiles: policy.maxFiles, MaxChangedLines: policy.maxChangedLines,
		MaxDurationSeconds: maxDurationSeconds,
		Usage:              usage, Attempts: attemptRecords, Candidates: candidates, ProviderExecutions: providerExecutions,
		ProviderFailures: providerFailures,
		FinalRun:         activeRun, GateState: activeRun.TerminalState,
		TerminalState: terminalState, Remaining: remaining, StopReasons: uniqueSorted(stopReasons),
	}
	if err := validateAttemptAudit(result); err != nil {
		return RemediationRun{}, err
	}
	result.RunID, err = remediationRunID(result)
	if err != nil {
		return RemediationRun{}, err
	}
	return result, nil
}

func remediationDeadlineReached(ctx context.Context) bool {
	return ctx.Err() == context.DeadlineExceeded
}

func enforceCandidateDeadline(ctx context.Context, candidate Candidate) (Candidate, bool, error) {
	if !candidate.Accepted || !remediationDeadlineReached(ctx) {
		return candidate, false, nil
	}
	candidate.Accepted = false
	candidate.Reasons = []string{"Candidate completed after the scanner-owned remediation duration budget expired."}
	identifier, err := candidateContentID(candidate)
	if err != nil {
		return Candidate{}, false, err
	}
	candidate.CandidateID = identifier
	return candidate, true, nil
}

func proposalPolicyReasonCode(err error) string {
	message := strings.ToLower(err.Error())
	switch {
	case strings.Contains(message, "anti-gaming audit"):
		return "anti_gaming_rejected"
	case strings.Contains(message, "outside the r2 fix contract") || strings.Contains(message, "protected path"):
		return "scope_rejected"
	case strings.Contains(message, "budget") || strings.Contains(message, "line contract") || strings.Contains(message, "file contract"):
		return "budget_rejected"
	default:
		return "policy_rejected"
	}
}

func candidateRejectionReason(candidate Candidate) string {
	if len(candidate.Reasons) == 0 {
		return "Candidate failed independent scanner verification."
	}
	return strings.Join(candidate.Reasons, " ")
}

func prepareCandidateContainer(targetRoot, destination string) (string, error) {
	destination, err := resolveCandidateContainer(targetRoot, destination)
	if err != nil {
		return "", err
	}
	if err := os.Mkdir(destination, 0o700); err != nil {
		return "", fmt.Errorf("create candidate root: %w", err)
	}
	return destination, nil
}

func resolveCandidateContainer(targetRoot, destination string) (string, error) {
	if strings.TrimSpace(destination) == "" {
		return "", fmt.Errorf("candidate root is required")
	}
	destination, err := filepath.Abs(destination)
	if err != nil {
		return "", fmt.Errorf("resolve candidate root: %w", err)
	}
	parent := filepath.Dir(destination)
	parentInfo, err := os.Stat(parent)
	if err != nil || !parentInfo.IsDir() {
		return "", fmt.Errorf("candidate root parent must be an existing directory")
	}
	if resolved, resolveErr := filepath.EvalSymlinks(parent); resolveErr == nil {
		destination = filepath.Join(resolved, filepath.Base(destination))
	}
	if pathWithin(targetRoot, destination) || pathWithin(destination, targetRoot) {
		return "", fmt.Errorf("candidate root must be outside the target tree")
	}
	if _, err := os.Lstat(destination); !os.IsNotExist(err) {
		return "", fmt.Errorf("candidate root already exists")
	}
	return destination, nil
}

func rebasedConfiguration(configuration *ProjectConfiguration, relative, target string) *ProjectConfiguration {
	if configuration == nil {
		return nil
	}
	copy := *configuration
	if relative != "" {
		copy.SourcePath = filepath.Join(target, filepath.FromSlash(relative))
	}
	return &copy
}

func nextDeterministicFailure(run model.RunResult, c *catalog.Catalog) string {
	for _, result := range run.Results {
		assertion := c.Assertions[result.AssertionID]
		if result.Assessment == "fail" && assertion.RemediationClass == "R1" && deterministicFixerRegistered(assertion) {
			return assertion.ID
		}
	}
	return ""
}

func nextAgentFailure(
	run model.RunResult,
	item model.Inventory,
	c *catalog.Catalog,
	protectedPaths []string,
	options AgentOptions,
) (string, provider.Task, error) {
	for _, result := range run.Results {
		if result.Assessment != "fail" {
			continue
		}
		assertion := c.Assertions[result.AssertionID]
		finding, ok := findingFor(run, result.AssertionID)
		if !ok {
			return "", provider.Task{}, fmt.Errorf("failed assertion %s has no canonical finding", result.AssertionID)
		}
		task, supported, err := planAgentTask(item, assertion, finding, protectedPaths, options)
		if err != nil {
			return "", provider.Task{}, err
		}
		if supported {
			return assertion.ID, task, nil
		}
	}
	return "", provider.Task{}, nil
}

func validateAgentOptions(options AgentOptions) error {
	if options.Provider == "" || options.Executable == "" || options.OutputSchemaPath == "" {
		return fmt.Errorf("agent provider, executable, and output schema are required")
	}
	capabilities, err := provider.ProviderCapabilities(options.Provider)
	if err != nil {
		return err
	}
	if options.MaxCostUSD > 0 && !capabilities.ProviderCostLimit {
		return fmt.Errorf("provider %s cannot enforce the requested cost limit", options.Provider)
	}
	if !options.AllowRemoteSourceProcessing {
		return policyDenied(fmt.Errorf("agent remediation requires --allow-remote-source-processing"))
	}
	if options.TimeoutSeconds < 0 || options.TimeoutSeconds > 3600 {
		return fmt.Errorf("agent timeout must be between 1 and 3600 seconds")
	}
	if options.MaxOutputBytes < 0 || options.MaxOutputBytes > 4*1024*1024 {
		return fmt.Errorf("agent output limit must be between 1 KiB and 4 MiB")
	}
	if math.IsNaN(options.MaxCostUSD) || math.IsInf(options.MaxCostUSD, 0) || options.MaxCostUSD < 0 || options.MaxCostUSD > 1000 {
		return fmt.Errorf("agent cost limit must be between 0 and 1000 USD")
	}
	return nil
}

func deterministicFixerRegistered(assertion model.Assertion) bool {
	return (assertion.ID == finalNewlineAssertion && assertion.ImplementationID == "prc.native.final-newline@0.1") ||
		(assertion.ID == restrictiveModesAssertion && assertion.ImplementationID == "prc.native.restrictive-file-modes@0.1")
}

func deterministicCost(item model.Inventory, assertionID string) (int, int, error) {
	switch assertionID {
	case finalNewlineAssertion:
		paths, err := finalNewlineViolations(item)
		if err != nil {
			return 0, 0, err
		}
		if len(paths) == 0 {
			return 0, 0, fmt.Errorf("failing final-newline assertion has no deterministic violations")
		}
		return len(paths), len(paths), nil
	case restrictiveModesAssertion:
		paths := restrictiveModeViolations(item)
		if len(paths) == 0 {
			return 0, 0, fmt.Errorf("failing restrictive-mode assertion has no deterministic violations")
		}
		return len(paths), 0, nil
	default:
		return 0, 0, fmt.Errorf("no deterministic fixer is registered for %s", assertionID)
	}
}

func classifyRemaining(run model.RunResult, c *catalog.Catalog, stoppedAssertion, stoppedCode string, agentEnabled bool) []RemainingWork {
	remaining := make([]RemainingWork, 0)
	for _, result := range run.Results {
		if result.Assessment == "pass" || result.Assessment == "not_applicable" {
			continue
		}
		assertion := c.Assertions[result.AssertionID]
		code, reason := remainingReason(assertion, result, stoppedAssertion, stoppedCode, agentEnabled)
		work := RemainingWork{
			AssertionID: assertion.ID, ControlIDs: append([]string{}, assertion.ControlIDs...), Title: assertion.Title,
			Assessment: result.Assessment, Execution: result.Execution, Severity: result.Severity, Gate: result.Gate,
			RemediationClass: result.RemediationClass, Summary: result.Summary, ReasonCode: code, Reason: reason,
		}
		if finding, ok := findingFor(run, result.AssertionID); ok {
			work.FindingID, work.FindingFingerprint = finding.ID, finding.Fingerprint
		}
		remaining = append(remaining, work)
	}
	return remaining
}

func remainingReason(assertion model.Assertion, result model.AssertionResult, stoppedAssertion, stoppedCode string, agentEnabled bool) (string, string) {
	if assertion.ID == stoppedAssertion && stoppedCode != "" {
		if stoppedCode == "budget_exhausted" {
			return stoppedCode, "The eligible deterministic fix exceeded the remaining configured attempt or change budget."
		}
		if stoppedCode == "time_budget_exhausted" {
			return stoppedCode, "The eligible fix was not started or accepted because the configured remediation duration budget was exhausted."
		}
		if stoppedCode == "provider_stopped" {
			return stoppedCode, "The read-only provider could not produce a candidate and the loop stopped without retrying."
		}
		if stoppedCode == "provider_failed" {
			return stoppedCode, "The provider invocation failed safely and the loop stopped without retrying."
		}
		return stoppedCode, "The independently verified candidate was rejected and was not used as the next baseline."
	}
	if result.Assessment == "manual_review" {
		return "human_evidence_required", "This assertion requires current scoped evidence from an accountable human reviewer."
	}
	if result.Assessment != "fail" {
		return "evidence_or_execution_required", "The result cannot be remediated until admissible evidence and successful deterministic execution establish a failing condition."
	}
	switch assertion.RemediationClass {
	case "R0":
		return "detect_explain_only", "This assertion is detect/explain only and has no authorized automated change."
	case "R1":
		return "no_registered_fixer", "No approved deterministic fixer is registered for this R1 assertion."
	case "R2":
		if agentEnabled {
			return "no_safe_agent_task", "No registered scanner-owned R2 task planner can safely bound and verify this finding."
		}
		return "agent_proposal_required", "A bounded agent proposal and independent deterministic verification are required; this loop does not invoke a provider."
	case "R3":
		return "human_review_required", "This high-impact engineering change requires human review before merge."
	case "R4":
		return "specialized_approval_required", "This data or release-sensitive change requires a specialized approval workflow."
	case "R5":
		return "separate_operational_workflow_required", "This operational or production action is outside the general remediation loop."
	case "R6":
		return "human_decision_required", "This finding requires an accountable human decision."
	default:
		return "unsupported_remediation_class", "The assertion has no supported remediation policy."
	}
}

func remediationRunID(run RemediationRun) (string, error) {
	run.RunID = ""
	return contentID(run)
}
