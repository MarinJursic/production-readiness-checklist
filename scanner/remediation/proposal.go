package remediation

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/MarinJursic/production-readiness-checklist/scanner/catalog"
	"github.com/MarinJursic/production-readiness-checklist/scanner/engine"
	"github.com/MarinJursic/production-readiness-checklist/scanner/inventory"
	"github.com/MarinJursic/production-readiness-checklist/scanner/provider"
	"github.com/MarinJursic/production-readiness-checklist/scanner/verifier"
)

const proposalFixer = "prc.provider-proposal@0.1"

// RunProposal applies one validated provider proposal in a fresh candidate and
// independently evaluates its deterministic scanner postconditions. It never
// changes the source workspace or runs provider-authored commands.
func RunProposal(options ProposalOptions) (Candidate, error) {
	if options.ProfileID == "" {
		options.ProfileID = "prc/core-repository"
	}
	if options.CandidateDir == "" {
		return Candidate{}, fmt.Errorf("candidate directory is required")
	}
	if options.Context == nil {
		options.Context = context.Background()
	}
	policy, err := resolvePolicy(options.Target, options.ProfileID, options.MaxFiles, options.MaxChangedLines,
		options.Attempt, options.MaxAttempts, options.Configuration)
	if err != nil {
		return Candidate{}, err
	}
	if err := provider.ValidateOutput(options.Provider, options.Output, options.Task); err != nil {
		return Candidate{}, policyDenied(err)
	}
	if options.Output.Status != "candidate" {
		return Candidate{}, policyDenied(fmt.Errorf("provider output status %s has no candidate patch", options.Output.Status))
	}

	c, err := catalog.Load(options.CatalogRoot)
	if err != nil {
		return Candidate{}, err
	}
	assertion, exists := c.Assertions[options.Task.AssertionID]
	if !exists {
		return Candidate{}, fmt.Errorf("unknown assertion %q", options.Task.AssertionID)
	}
	if assertion.RemediationClass != "R2" {
		return Candidate{}, policyDenied(fmt.Errorf("assertion %s is not eligible for isolated R2 proposal application", assertion.ID))
	}
	if !sameStrings(assertion.ControlIDs, options.Task.ControlIDs) {
		return Candidate{}, policyDenied(fmt.Errorf("agent task controls do not match assertion %s", assertion.ID))
	}

	baseline, err := inventory.Build(options.Target)
	if err != nil {
		return Candidate{}, err
	}
	if options.TargetName != "" {
		baseline.TargetName = options.TargetName
	}
	baseline, err = policy.bind(baseline, "")
	if err != nil {
		return Candidate{}, err
	}
	if baseline.Digest != options.Task.WorkspaceInventoryDigest {
		return Candidate{}, policyDenied(fmt.Errorf("source workspace does not match the sealed agent task"))
	}
	scanner := engine.New(c)
	beforeRun, err := scanner.Scan(options.ProfileID, baseline)
	if err != nil {
		return Candidate{}, err
	}
	before, ok := resultFor(beforeRun, assertion.ID)
	if !ok || before.Assessment != "fail" {
		return Candidate{}, policyDenied(fmt.Errorf("assertion %s is not a failing finding in the baseline", assertion.ID))
	}
	baselineFinding, ok := findingFor(beforeRun, assertion.ID)
	if !ok {
		return Candidate{}, fmt.Errorf("failed assertion %s has no canonical baseline finding", assertion.ID)
	}
	if options.Task.FindingID != baselineFinding.ID || options.Task.FindingFingerprint != baselineFinding.Fingerprint {
		return Candidate{}, policyDenied(fmt.Errorf("agent task does not match the exact current finding for assertion %s", assertion.ID))
	}
	if reasons := auditProposalAntiGaming(baseline, options.Output); len(reasons) > 0 {
		return Candidate{}, policyDenied(fmt.Errorf("provider proposal failed anti-gaming audit: %s", strings.Join(reasons, " ")))
	}
	if options.Verifier == nil {
		return Candidate{}, policyDenied(fmt.Errorf("R2 remediation requires an independent sandbox verifier"))
	}
	verificationKind, err := verifier.InferKind(options.Task.RelevantPaths[0])
	if err != nil {
		return Candidate{}, policyDenied(err)
	}
	verificationOptions := *options.Verifier
	verificationOptions.Kind = verificationKind
	if err := verifier.ValidateOptions(verificationOptions); err != nil {
		return Candidate{}, policyDenied(err)
	}
	if err := verifier.Preflight(options.Context, verificationOptions); err != nil {
		return Candidate{}, policyDenied(err)
	}

	proposalID, err := provider.OutputID(options.Output)
	if err != nil {
		return Candidate{}, err
	}
	protectedPaths := proposalProtectedPaths(options.Task.ProtectedPaths, policy.protectedPaths)
	contract := FixContract{
		SchemaVersion: FixContractSchema, BaselineRunID: beforeRun.RunID,
		BaselineInventoryDigest: baseline.Digest,
		FindingID:               baselineFinding.ID,
		FindingFingerprint:      baselineFinding.Fingerprint,
		ProposalFindingID:       options.Task.FindingID,
		AssertionID:             assertion.ID,
		ConfigurationDigest:     policy.configurationID, ProjectID: policy.projectID,
		ControlIDs: append([]string(nil), assertion.ControlIDs...), Goal: options.Task.Goal,
		FixerID: proposalFixer, RemediationClass: "R2", Provider: options.Provider,
		ProposalTaskID: options.Task.TaskID, ProposalSHA256: proposalID,
		AllowedPaths: append([]string(nil), options.Output.ChangedFiles...), ProtectedPaths: protectedPaths,
		Network: "deny", MaxChangedLines: policy.maxChangedLines, MaxFiles: policy.maxFiles,
		Attempt: policy.attempt, MaxAttempts: policy.maxAttempts,
		Acceptance: []string{
			"The scanner-owned parser applies exactly the validated proposal in a new isolated candidate.",
			"No unproposed or protected path changes, no file deletion, and no mode change occur.",
			"The target assertion passes and every baseline passing assertion remains passing.",
			"The original workspace inventory remains byte-for-byte unchanged.",
		},
	}
	sort.Strings(contract.ControlIDs)
	contract.TaskID, err = contentID(contract)
	if err != nil {
		return Candidate{}, err
	}

	candidateRoot, err := prepareCandidate(baseline, options.CandidateDir)
	if err != nil {
		return Candidate{}, err
	}
	changes, err := applyProviderPatch(candidateRoot, baseline, options.Task, options.Output,
		policy.protectedPaths, policy.maxFiles, policy.maxChangedLines)
	if err != nil {
		return Candidate{}, policyDenied(err)
	}
	candidateInventory, err := inventory.Build(candidateRoot)
	if err != nil {
		return Candidate{}, err
	}
	candidateInventory, err = inventory.BindDerivedSource(candidateInventory, baseline)
	if err != nil {
		return Candidate{}, err
	}
	if options.TargetName != "" {
		candidateInventory.TargetName = options.TargetName
	}
	candidateInventory, err = policy.bind(candidateInventory, candidateRoot)
	if err != nil {
		return Candidate{}, err
	}
	reasons := auditProviderCandidate(baseline, candidateInventory, changes, protectedPaths)
	afterRun, err := scanner.Scan(options.ProfileID, candidateInventory)
	if err != nil {
		return Candidate{}, err
	}
	after, ok := resultFor(afterRun, assertion.ID)
	if !ok || after.Assessment != "pass" {
		reasons = append(reasons, "Target assertion did not pass in the candidate.")
	}
	if findingFingerprintPresent(afterRun, baselineFinding.Fingerprint) {
		reasons = append(reasons, "The exact baseline finding fingerprint remains in the candidate.")
	}
	for _, result := range beforeRun.Results {
		if result.Assessment != "pass" {
			continue
		}
		candidateResult, found := resultFor(afterRun, result.AssertionID)
		if !found || candidateResult.Assessment != "pass" {
			reasons = append(reasons, "Baseline passing assertion regressed: "+result.AssertionID+".")
		}
	}
	var verification *verifier.Execution
	if len(reasons) == 0 {
		plan, planErr := verifier.BuildPlan(candidateRoot, candidateInventory.Digest, verificationOptions)
		if planErr != nil {
			reasons = append(reasons, "Independent sandbox verification preflight failed closed: "+planErr.Error()+".")
		} else {
			execution, executionErr := verifier.Run(options.Context, plan)
			if executionErr != nil {
				reasons = append(reasons, "Independent sandbox verification failed closed: "+executionErr.Error()+".")
			} else {
				verification = &execution
				if execution.Outcome == "workspace_changed" {
					return Candidate{}, fmt.Errorf("verification containment failed: candidate workspace changed")
				}
				if execution.Outcome != "pass" {
					reasons = append(reasons, verificationRejectionReason(execution))
				}
			}
		}
	}
	currentBaseline, err := inventory.Build(options.Target)
	if err == nil && options.TargetName != "" {
		currentBaseline.TargetName = options.TargetName
	}
	if err == nil {
		currentBaseline, err = policy.bind(currentBaseline, "")
	}
	if err != nil || currentBaseline.Root != baseline.Root || currentBaseline.Digest != baseline.Digest {
		return Candidate{}, fmt.Errorf("source workspace changed during proposal remediation")
	}
	reasons = uniqueSorted(reasons)
	candidate := Candidate{
		SchemaVersion: CandidateSchema, CandidatePath: candidateRoot, Contract: contract,
		CandidateInventoryDigest: candidateInventory.Digest, CandidateRunID: afterRun.RunID,
		Changes: changes, BeforeAssessment: before.Assessment, AfterAssessment: after.Assessment,
		Accepted: len(reasons) == 0, Reasons: reasons, Verification: verification,
	}
	candidate.CandidateID, err = candidateContentID(candidate)
	if err != nil {
		return Candidate{}, err
	}
	return candidate, nil
}

func verificationRejectionReason(execution verifier.Execution) string {
	switch execution.ReasonCode {
	case "tests_failed":
		return "Independent scanner-owned tests failed in the sandbox."
	case "timeout":
		return "Independent sandbox verification exceeded its scanner-owned timeout."
	case "cancelled":
		return "Independent sandbox verification was cancelled."
	case "output_limit":
		return "Independent sandbox verification exceeded its scanner-owned output limit."
	case "runtime_failed":
		return "Independent sandbox verification infrastructure failed closed."
	default:
		return "Independent sandbox verification did not pass."
	}
}

func sameStrings(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	for index := range left {
		if left[index] != right[index] {
			return false
		}
	}
	return true
}
