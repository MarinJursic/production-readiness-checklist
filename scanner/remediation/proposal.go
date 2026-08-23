package remediation

import (
	"fmt"
	"sort"

	"github.com/MarinJursic/production-readiness-checklist/scanner/catalog"
	"github.com/MarinJursic/production-readiness-checklist/scanner/engine"
	"github.com/MarinJursic/production-readiness-checklist/scanner/inventory"
	"github.com/MarinJursic/production-readiness-checklist/scanner/provider"
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
	if options.MaxFiles < 1 || options.MaxFiles > maximumFixFiles {
		return Candidate{}, fmt.Errorf("max files must be between 1 and %d", maximumFixFiles)
	}
	if options.MaxChangedLines < 1 || options.MaxChangedLines > maximumFixLines {
		return Candidate{}, fmt.Errorf("max changed lines must be between 1 and %d", maximumFixLines)
	}
	if err := provider.ValidateOutput(options.Provider, options.Output, options.Task); err != nil {
		return Candidate{}, err
	}
	if options.Output.Status != "candidate" {
		return Candidate{}, fmt.Errorf("provider output status %s has no candidate patch", options.Output.Status)
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
		return Candidate{}, fmt.Errorf("assertion %s is not eligible for isolated R2 proposal application", assertion.ID)
	}
	if !sameStrings(assertion.ControlIDs, options.Task.ControlIDs) {
		return Candidate{}, fmt.Errorf("agent task controls do not match assertion %s", assertion.ID)
	}

	baseline, err := inventory.Build(options.Target)
	if err != nil {
		return Candidate{}, err
	}
	if baseline.Digest != options.Task.WorkspaceInventoryDigest {
		return Candidate{}, fmt.Errorf("source workspace does not match the sealed agent task")
	}
	scanner := engine.New(c)
	beforeRun, err := scanner.Scan(options.ProfileID, baseline)
	if err != nil {
		return Candidate{}, err
	}
	before, ok := resultFor(beforeRun, assertion.ID)
	if !ok || before.Assessment != "fail" {
		return Candidate{}, fmt.Errorf("assertion %s is not a failing finding in the baseline", assertion.ID)
	}

	proposalID, err := provider.OutputID(options.Output)
	if err != nil {
		return Candidate{}, err
	}
	protectedPaths := proposalProtectedPaths(options.Task.ProtectedPaths)
	contract := FixContract{
		SchemaVersion: FixContractSchema, BaselineRunID: beforeRun.RunID,
		BaselineInventoryDigest: baseline.Digest, AssertionID: assertion.ID,
		ControlIDs: append([]string(nil), assertion.ControlIDs...), Goal: options.Task.Goal,
		FixerID: proposalFixer, RemediationClass: "R2", Provider: options.Provider,
		ProposalTaskID: options.Task.TaskID, ProposalSHA256: proposalID,
		AllowedPaths: append([]string(nil), options.Output.ChangedFiles...), ProtectedPaths: protectedPaths,
		Network: "deny", MaxChangedLines: options.MaxChangedLines, MaxFiles: options.MaxFiles,
		MaxAttempts: 1,
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
	changes, err := applyProviderPatch(candidateRoot, baseline, options.Task, options.Output, options.MaxFiles, options.MaxChangedLines)
	if err != nil {
		return Candidate{}, err
	}
	candidateInventory, err := inventory.Build(candidateRoot)
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
	for _, result := range beforeRun.Results {
		if result.Assessment != "pass" {
			continue
		}
		candidateResult, found := resultFor(afterRun, result.AssertionID)
		if !found || candidateResult.Assessment != "pass" {
			reasons = append(reasons, "Baseline passing assertion regressed: "+result.AssertionID+".")
		}
	}
	currentBaseline, err := inventory.Build(options.Target)
	if err != nil || currentBaseline.Digest != baseline.Digest {
		return Candidate{}, fmt.Errorf("source workspace changed during proposal remediation")
	}
	reasons = uniqueSorted(reasons)
	candidate := Candidate{
		SchemaVersion: CandidateSchema, CandidatePath: candidateRoot, Contract: contract,
		CandidateInventoryDigest: candidateInventory.Digest, CandidateRunID: afterRun.RunID,
		Changes: changes, BeforeAssessment: before.Assessment, AfterAssessment: after.Assessment,
		Accepted: len(reasons) == 0, Reasons: reasons,
	}
	candidate.CandidateID, err = candidateContentID(candidate)
	if err != nil {
		return Candidate{}, err
	}
	return candidate, nil
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
