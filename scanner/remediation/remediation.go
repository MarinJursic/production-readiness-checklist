package remediation

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/MarinJursic/production-readiness-checklist/scanner/catalog"
	"github.com/MarinJursic/production-readiness-checklist/scanner/engine"
	"github.com/MarinJursic/production-readiness-checklist/scanner/inventory"
	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
)

const (
	finalNewlineAssertion     = "PRC-A-CORE-014"
	finalNewlineFixer         = "prc.fixer.final-newline@0.1"
	restrictiveModesAssertion = "PRC-A-CORE-022"
	restrictiveModesFixer     = "prc.fixer.restrictive-file-modes@0.1"
	maximumFixFiles           = 1000
	maximumFixLines           = 1000
)

// Run creates and validates one isolated deterministic remediation candidate.
// It never mutates the target workspace and never merges, deploys, or releases.
func Run(options Options) (Candidate, error) {
	if options.ProfileID == "" {
		options.ProfileID = "prc/core-repository"
	}
	if options.AssertionID == "" {
		options.AssertionID = finalNewlineAssertion
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
	c, err := catalog.Load(options.CatalogRoot)
	if err != nil {
		return Candidate{}, err
	}
	assertion, exists := c.Assertions[options.AssertionID]
	if !exists {
		return Candidate{}, fmt.Errorf("unknown assertion %q", options.AssertionID)
	}
	if assertion.RemediationClass != "R1" {
		return Candidate{}, fmt.Errorf("assertion %s is not eligible for the registered R1 fixer", assertion.ID)
	}

	baseline, err := inventory.Build(options.Target)
	if err != nil {
		return Candidate{}, err
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
	var paths []string
	var fixerID, goal string
	var acceptance []string
	switch {
	case assertion.ID == finalNewlineAssertion && assertion.ImplementationID == "prc.native.final-newline@0.1":
		paths, err = finalNewlineViolations(baseline)
		fixerID = finalNewlineFixer
		goal = "Append one line-feed byte to each recognized source file whose final byte is not a line-feed; an empty file has no final byte."
		acceptance = []string{
			"Every allowed path differs only by one appended line-feed byte.",
			"No path outside the allowlist changes and no protected path changes.",
			"The target assertion passes and every baseline passing assertion remains passing.",
			"The original workspace remains byte-for-byte and mode-for-mode unchanged.",
		}
	case assertion.ID == restrictiveModesAssertion && assertion.ImplementationID == "prc.native.restrictive-file-modes@0.1":
		paths = restrictiveModeViolations(baseline)
		fixerID = restrictiveModesFixer
		goal = "Clear group-write and other-write permission bits on each inventoried regular file that has them, preserving file bytes and all other permission bits."
		acceptance = []string{
			"Every allowed path has only its group-write and other-write permission bits cleared.",
			"Every allowed path remains byte-for-byte identical to its baseline content.",
			"No path outside the allowlist changes and no protected path changes.",
			"The target assertion passes and every baseline passing assertion remains passing.",
			"The original workspace remains byte-for-byte and mode-for-mode unchanged.",
		}
	default:
		return Candidate{}, fmt.Errorf("no deterministic fixer is registered for %s", options.AssertionID)
	}
	if err != nil {
		return Candidate{}, err
	}
	if len(paths) == 0 {
		return Candidate{}, fmt.Errorf("assertion %s reported no deterministic violations", assertion.ID)
	}
	if len(paths) > options.MaxFiles {
		return Candidate{}, fmt.Errorf("fix requires %d files, above the configured budget", len(paths))
	}
	if fixerID == finalNewlineFixer && len(paths) > options.MaxChangedLines {
		return Candidate{}, fmt.Errorf("fix requires %d changed lines, above the configured budget", len(paths))
	}
	for _, path := range paths {
		if protected(path, defaultProtectedPaths) {
			return Candidate{}, fmt.Errorf("fix would touch protected path %s", path)
		}
	}

	contract := FixContract{
		SchemaVersion: FixContractSchema, BaselineRunID: beforeRun.RunID,
		BaselineInventoryDigest: baseline.Digest, AssertionID: assertion.ID,
		ControlIDs: append([]string(nil), assertion.ControlIDs...),
		Goal:       goal,
		FixerID:    fixerID, RemediationClass: "R1",
		AllowedPaths:   append([]string(nil), paths...),
		ProtectedPaths: append([]string(nil), defaultProtectedPaths...),
		Network:        "deny", MaxChangedLines: options.MaxChangedLines, MaxFiles: options.MaxFiles,
		MaxAttempts: 1,
		Acceptance:  acceptance,
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
	switch fixerID {
	case finalNewlineFixer:
		err = applyFinalNewline(candidateRoot, paths)
	case restrictiveModesFixer:
		err = applyRestrictiveModes(candidateRoot, baseline, paths)
	}
	if err != nil {
		return Candidate{}, err
	}
	candidateInventory, err := inventory.Build(candidateRoot)
	if err != nil {
		return Candidate{}, err
	}
	changes, reasons := auditCandidate(baseline, candidateInventory, contract)
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
		return Candidate{}, fmt.Errorf("source workspace changed during deterministic remediation")
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

func resultFor(run model.RunResult, assertionID string) (model.AssertionResult, bool) {
	for _, result := range run.Results {
		if result.AssertionID == assertionID {
			return result, true
		}
	}
	return model.AssertionResult{}, false
}

func contentID(value any) (string, error) {
	payload, err := json.Marshal(value)
	if err != nil {
		return "", fmt.Errorf("encode content identity: %w", err)
	}
	digest := sha256.Sum256(payload)
	return hex.EncodeToString(digest[:]), nil
}

func candidateContentID(candidate Candidate) (string, error) {
	candidate.CandidateID = ""
	candidate.CandidatePath = ""
	return contentID(candidate)
}
