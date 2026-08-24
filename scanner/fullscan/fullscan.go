// Package fullscan binds the complete source control registry to one scanner
// run. It deliberately keeps broad control dispositions separate from the
// narrower executable assertion results.
package fullscan

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/MarinJursic/production-readiness-checklist/scanner/catalog"
	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
)

const (
	registrySchema          = "prc.control-id-registry/v0.1"
	contractSchema          = "prc.control-contracts/v0.1"
	controlCatalogSchema    = "prc.control-catalog-summary/v0.1"
	maximumRegistryBytes    = 16 * 1024 * 1024
	maximumContractBytes    = 16 * 1024 * 1024
	maximumSourceFileBytes  = 16 * 1024 * 1024
	maximumSourceTotalBytes = 128 * 1024 * 1024
	maximumControls         = 100_000
)

var (
	controlIDPattern = regexp.MustCompile(`^(?:PRC-[0-9]{2}-[0-9]{3}|USEQ-[A-F0-9]{8})$`)
	versionPattern   = regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+$`)
	digestPattern    = regexp.MustCompile(`^[0-9a-f]{64}$`)
)

// ErrRegistryUnavailable lets deliberately minimal custom test catalogs keep
// using profile-only scans. Released scanner bundles always include the full
// registry, and a present but invalid registry still fails closed.
var ErrRegistryUnavailable = errors.New("complete control registry is unavailable")

type registryDocument struct {
	SchemaVersion   string          `json:"schema_version"`
	RegistryVersion string          `json:"registry_version"`
	SourceSHA256    string          `json:"source_sha256"`
	EntryCount      int             `json:"entry_count"`
	Entries         []model.Control `json:"entries"`
}

type contractDocument struct {
	SchemaVersion   string            `json:"schema_version"`
	RegistrySHA256  string            `json:"registry_sha256"`
	RegistryVersion string            `json:"registry_version"`
	ContractCount   int               `json:"contract_count"`
	Contracts       []controlContract `json:"contracts"`
}

type controlContract struct {
	ControlID                 string   `json:"control_id"`
	Revision                  int      `json:"revision"`
	ContractStatus            string   `json:"contract_status"`
	CanonicalControlID        string   `json:"canonical_control_id"`
	EvaluationClass           string   `json:"evaluation_class"`
	AutomationClass           string   `json:"automation_class"`
	ApplicabilityClass        string   `json:"applicability_class"`
	Atomicity                 string   `json:"atomicity"`
	CompleteInventoryRequired bool     `json:"complete_inventory_required"`
	NegativeCondition         bool     `json:"negative_condition"`
	ProjectThresholdsRequired bool     `json:"project_thresholds_required"`
	EvidenceAuthorities       []string `json:"evidence_authorities"`
	NotApplicableProof        string   `json:"not_applicable_proof"`
	ReviewerStatus            string   `json:"reviewer_status"`
	ContractSHA256            string   `json:"contract_sha256"`
}

// Attach loads and verifies the complete registry, adds one result for every
// registered control, and recomputes the content-addressed run identity.
func Attach(root string, scannerCatalog *catalog.Catalog, run model.RunResult) (model.RunResult, error) {
	registry, registryDigest, err := loadRegistry(root)
	if err != nil {
		return model.RunResult{}, err
	}
	if err := validateCatalogReferences(scannerCatalog, registry); err != nil {
		return model.RunResult{}, err
	}
	contracts, contractDigest, err := loadContracts(root, registry, registryDigest)
	if err != nil {
		return model.RunResult{}, err
	}
	profileState := run.TerminalState
	run.ControlCatalog = &model.ControlCatalogSummary{
		SchemaVersion: controlCatalogSchema, RegistryVersion: registry.RegistryVersion,
		RegistrySHA256: registryDigest, SourceSHA256: registry.SourceSHA256,
		ContractSchemaVersion: contracts.SchemaVersion, ContractSHA256: contractDigest,
		ControlCount: len(registry.Entries), ContractCount: len(contracts.Contracts),
		ProfileTerminalState: profileState,
	}
	run.ControlResults = make([]model.ControlResult, 0, len(registry.Entries))
	active := 0
	allAssertions := assertionsByControl(scannerCatalog)
	executed := executedByControl(run.Results)
	for index, control := range registry.Entries {
		if control.Status == "active" {
			active++
		}
		contract := contracts.Contracts[index]
		if contract.ReviewerStatus == "expert_reviewed" {
			run.ControlCatalog.ExpertReviewedContractCount++
		} else {
			run.ControlCatalog.GeneratedContractCount++
		}
		run.ControlResults = append(run.ControlResults, controlResult(control, contract, allAssertions[control.ID], executed[control.ID]))
	}
	run.ControlCatalog.ActiveControlCount = active
	if run.TerminalState == "profile_satisfied" && hasIncompleteControls(run.ControlResults) {
		run.TerminalState = "assessment_incomplete"
	}
	return Reidentify(run)
}

func loadContracts(root string, registry registryDocument, registryDigest string) (contractDocument, string, error) {
	absolute, err := filepath.Abs(root)
	if err != nil {
		return contractDocument{}, "", fmt.Errorf("resolve catalog root: %w", err)
	}
	path := filepath.Join(absolute, "catalog", "control-contracts.json")
	data, err := readRegularBounded(path, maximumContractBytes, "control contracts")
	if err != nil {
		return contractDocument{}, "", err
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	var document contractDocument
	if err := decoder.Decode(&document); err != nil {
		return contractDocument{}, "", fmt.Errorf("parse control contracts: %w", err)
	}
	var trailing any
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		return contractDocument{}, "", fmt.Errorf("control contracts contain trailing JSON")
	}
	if err := validateContracts(document, registry, registryDigest); err != nil {
		return contractDocument{}, "", err
	}
	digest := sha256.Sum256(data)
	return document, hex.EncodeToString(digest[:]), nil
}

func validateContracts(document contractDocument, registry registryDocument, registryDigest string) error {
	if document.SchemaVersion != contractSchema || document.RegistrySHA256 != registryDigest ||
		document.RegistryVersion != registry.RegistryVersion || document.ContractCount != len(document.Contracts) ||
		document.ContractCount != len(registry.Entries) || document.ContractCount < 1 || document.ContractCount > maximumControls {
		return fmt.Errorf("control contracts do not exactly bind the current registry")
	}
	registered := make(map[string]model.Control, len(registry.Entries))
	for _, control := range registry.Entries {
		registered[control.ID] = control
	}
	validEvaluation := map[string]bool{"repository": true, "environment": true, "human_external": true, "mixed": true, "unclassified": true}
	validAutomation := map[string]bool{
		"deterministic_candidate": true, "ai_advisory_candidate": true, "environment_evidence_required": true,
		"human_or_external_required": true, "mixed_evidence_required": true,
	}
	validAuthority := map[string]bool{"declared": true, "repository": true, "artifact": true, "environment": true, "human": true}
	for index, contract := range document.Contracts {
		control := registry.Entries[index]
		if contract.ControlID != control.ID || contract.Revision != control.Revision ||
			!digestPattern.MatchString(contract.ContractSHA256) || !registeredCanonical(contract, control, registered) ||
			!validEvaluation[contract.EvaluationClass] || !validAutomation[contract.AutomationClass] ||
			(contract.ApplicabilityClass != "conditional" && contract.ApplicabilityClass != "scope_required") ||
			(contract.Atomicity != "apparently_atomic" && contract.Atomicity != "compound_review_required") ||
			strings.TrimSpace(contract.NotApplicableProof) == "" || len(contract.NotApplicableProof) > 1000 ||
			len(contract.EvidenceAuthorities) == 0 {
			return fmt.Errorf("control %s has an invalid machine-readable contract", control.ID)
		}
		if control.Status == "retired" {
			if contract.ContractStatus != "retired" {
				return fmt.Errorf("retired control %s has an active contract", control.ID)
			}
		} else if contract.ContractStatus != "generated_unreviewed" && contract.ContractStatus != "reviewed" {
			return fmt.Errorf("active control %s has an unsupported contract status", control.ID)
		}
		if contract.ContractStatus == "reviewed" && contract.ReviewerStatus != "expert_reviewed" ||
			contract.ContractStatus != "reviewed" && contract.ReviewerStatus != "generated_unreviewed" {
			return fmt.Errorf("control %s contract review status is inconsistent", control.ID)
		}
		seenAuthorities := map[string]bool{}
		for _, authority := range contract.EvidenceAuthorities {
			if !validAuthority[authority] || seenAuthorities[authority] {
				return fmt.Errorf("control %s has an invalid evidence authority", control.ID)
			}
			seenAuthorities[authority] = true
		}
		if contractDigest(contract) != contract.ContractSHA256 {
			return fmt.Errorf("control %s contract digest does not match its content", control.ID)
		}
	}
	return nil
}

func registeredCanonical(contract controlContract, control model.Control, registered map[string]model.Control) bool {
	canonical, ok := registered[contract.CanonicalControlID]
	return ok && contract.CanonicalControlID <= control.ID &&
		normalizedStatement(canonical.Statement) == normalizedStatement(control.Statement)
}

func normalizedStatement(value string) string {
	return strings.TrimSuffix(strings.Join(strings.Fields(strings.ToLower(value)), " "), ".")
}

func contractDigest(contract controlContract) string {
	unsigned := map[string]any{
		"control_id": contract.ControlID, "revision": contract.Revision,
		"contract_status": contract.ContractStatus, "canonical_control_id": contract.CanonicalControlID,
		"evaluation_class": contract.EvaluationClass, "automation_class": contract.AutomationClass,
		"applicability_class": contract.ApplicabilityClass, "atomicity": contract.Atomicity,
		"complete_inventory_required": contract.CompleteInventoryRequired,
		"negative_condition":          contract.NegativeCondition,
		"project_thresholds_required": contract.ProjectThresholdsRequired,
		"evidence_authorities":        contract.EvidenceAuthorities, "not_applicable_proof": contract.NotApplicableProof,
		"reviewer_status": contract.ReviewerStatus,
	}
	payload, _ := json.Marshal(unsigned)
	value := sha256.Sum256(payload)
	return hex.EncodeToString(value[:])
}

// Reidentify returns a run whose ID binds every current field. It is used again
// after optional advisory AI reviews are merged into the complete control set.
func Reidentify(run model.RunResult) (model.RunResult, error) {
	run.RunID = ""
	payload, err := json.Marshal(run)
	if err != nil {
		return model.RunResult{}, fmt.Errorf("encode complete run identity: %w", err)
	}
	digest := sha256.Sum256(payload)
	run.RunID = hex.EncodeToString(digest[:])
	return run, nil
}

func loadRegistry(root string) (registryDocument, string, error) {
	absolute, err := filepath.Abs(root)
	if err != nil {
		return registryDocument{}, "", fmt.Errorf("resolve catalog root: %w", err)
	}
	path := filepath.Join(absolute, "catalog", "control-id-registry.json")
	if _, err := os.Lstat(path); os.IsNotExist(err) {
		return registryDocument{}, "", fmt.Errorf("%w: %s", ErrRegistryUnavailable, path)
	}
	data, err := readRegularBounded(path, maximumRegistryBytes, "control registry")
	if err != nil {
		return registryDocument{}, "", err
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	var document registryDocument
	if err := decoder.Decode(&document); err != nil {
		return registryDocument{}, "", fmt.Errorf("parse control registry: %w", err)
	}
	var trailing any
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		return registryDocument{}, "", fmt.Errorf("control registry contains trailing JSON")
	}
	if err := validateRegistry(absolute, document); err != nil {
		return registryDocument{}, "", err
	}
	digest := sha256.Sum256(data)
	return document, hex.EncodeToString(digest[:]), nil
}

func validateRegistry(root string, document registryDocument) error {
	if document.SchemaVersion != registrySchema {
		return fmt.Errorf("unsupported control registry schema %q", document.SchemaVersion)
	}
	if !versionPattern.MatchString(document.RegistryVersion) {
		return fmt.Errorf("invalid control registry version %q", document.RegistryVersion)
	}
	if !digestPattern.MatchString(document.SourceSHA256) {
		return fmt.Errorf("invalid control registry source digest")
	}
	if document.EntryCount != len(document.Entries) || document.EntryCount == 0 || document.EntryCount > maximumControls {
		return fmt.Errorf("control registry entry count %d does not match its bounded entries", document.EntryCount)
	}
	seen := make(map[string]bool, len(document.Entries))
	for index, entry := range document.Entries {
		if !controlIDPattern.MatchString(entry.ID) || seen[entry.ID] {
			return fmt.Errorf("control registry contains invalid or duplicate ID %q", entry.ID)
		}
		seen[entry.ID] = true
		if index > 0 && document.Entries[index-1].ID >= entry.ID {
			return fmt.Errorf("control registry entries must be ordered by ID")
		}
		if entry.Status != "active" && entry.Status != "retired" {
			return fmt.Errorf("control %s has unsupported status %q", entry.ID, entry.Status)
		}
		if entry.Revision < 1 || strings.TrimSpace(entry.Statement) == "" || len(entry.Statement) > 16*1024 ||
			!digestPattern.MatchString(entry.SemanticSHA256) || entry.Source.Line < 1 {
			return fmt.Errorf("control %s has invalid revision, statement, digest, or source line", entry.ID)
		}
		if err := safeSourcePath(entry.Source.Path); err != nil {
			return fmt.Errorf("control %s source: %w", entry.ID, err)
		}
	}
	sourceDigest, lines, err := sourceSnapshot(root)
	if err != nil {
		return err
	}
	if sourceDigest != document.SourceSHA256 {
		return fmt.Errorf("control registry source digest does not match the complete source documents")
	}
	for _, entry := range document.Entries {
		if entry.Status != "active" {
			continue
		}
		pathLines, ok := lines[entry.Source.Path]
		if !ok || entry.Source.Line > len(pathLines) {
			return fmt.Errorf("control %s points outside its source document", entry.ID)
		}
		expected := "- [ ] **" + entry.ID + "** — " + entry.Statement
		if pathLines[entry.Source.Line-1] != expected {
			return fmt.Errorf("control %s does not exactly match its source line", entry.ID)
		}
	}
	return nil
}

func sourceSnapshot(root string) (string, map[string][]string, error) {
	patterns := []string{
		filepath.Join(root, "docs", "checklists", "*.md"),
		filepath.Join(root, "docs", "engineering", "[0-9][0-9]-*.md"),
	}
	paths := make([]string, 0)
	for _, pattern := range patterns {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return "", nil, fmt.Errorf("list control source documents: %w", err)
		}
		paths = append(paths, matches...)
	}
	sort.Strings(paths)
	if len(paths) == 0 {
		return "", nil, fmt.Errorf("no control source documents found")
	}
	hasher := sha256.New()
	lines := make(map[string][]string, len(paths))
	total := 0
	for _, path := range paths {
		data, err := readRegularBounded(path, maximumSourceFileBytes, "control source")
		if err != nil {
			return "", nil, err
		}
		total += len(data)
		if total > maximumSourceTotalBytes {
			return "", nil, fmt.Errorf("control source documents exceed their total byte limit")
		}
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return "", nil, fmt.Errorf("resolve control source path: %w", err)
		}
		relative = filepath.ToSlash(relative)
		hasher.Write([]byte(relative))
		hasher.Write([]byte{0})
		hasher.Write(data)
		hasher.Write([]byte{0})
		lines[relative] = splitLines(data)
	}
	return hex.EncodeToString(hasher.Sum(nil)), lines, nil
}

func splitLines(data []byte) []string {
	result := make([]string, 0)
	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Buffer(make([]byte, 64*1024), maximumSourceFileBytes)
	for scanner.Scan() {
		result = append(result, scanner.Text())
	}
	return result
}

func safeSourcePath(path string) error {
	if path == "" || filepath.IsAbs(path) || filepath.ToSlash(filepath.Clean(filepath.FromSlash(path))) != path ||
		strings.HasPrefix(path, "../") || !strings.HasSuffix(path, ".md") ||
		(!strings.HasPrefix(path, "docs/checklists/") && !strings.HasPrefix(path, "docs/engineering/")) {
		return fmt.Errorf("unsafe or unsupported path %q", path)
	}
	return nil
}

func readRegularBounded(path string, limit int64, label string) ([]byte, error) {
	before, err := os.Lstat(path)
	if err != nil || !before.Mode().IsRegular() || before.Mode()&os.ModeSymlink != 0 || before.Size() > limit {
		return nil, fmt.Errorf("%s %s must be a regular file no larger than %d bytes", label, path, limit)
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open %s %s: %w", label, path, err)
	}
	defer file.Close()
	opened, err := file.Stat()
	if err != nil || !opened.Mode().IsRegular() || !os.SameFile(before, opened) {
		return nil, fmt.Errorf("%s %s changed while opening", label, path)
	}
	data, err := io.ReadAll(io.LimitReader(file, limit+1))
	if err != nil {
		return nil, fmt.Errorf("read %s %s: %w", label, path, err)
	}
	if int64(len(data)) > limit || int64(len(data)) != opened.Size() {
		return nil, fmt.Errorf("%s %s changed or exceeded its byte limit", label, path)
	}
	after, err := file.Stat()
	if err != nil || !os.SameFile(opened, after) || after.Size() != opened.Size() {
		return nil, fmt.Errorf("%s %s changed while reading", label, path)
	}
	return data, nil
}

func validateCatalogReferences(scannerCatalog *catalog.Catalog, document registryDocument) error {
	registered := make(map[string]model.Control, len(document.Entries))
	for _, control := range document.Entries {
		registered[control.ID] = control
	}
	for id, objective := range scannerCatalog.Objectives {
		control, ok := registered[id]
		if !ok || control.Status != "active" || control.Revision != objective.Revision ||
			control.Statement != objective.Statement || control.Source != objective.Source {
			return fmt.Errorf("executable objective %s does not exactly match the active control registry", id)
		}
	}
	for _, assertion := range scannerCatalog.Assertions {
		for _, id := range assertion.ControlIDs {
			if control, ok := registered[id]; !ok || control.Status != "active" {
				return fmt.Errorf("assertion %s references missing or retired control %s", assertion.ID, id)
			}
		}
	}
	return nil
}

func assertionsByControl(scannerCatalog *catalog.Catalog) map[string][]string {
	result := map[string][]string{}
	for _, assertion := range scannerCatalog.Assertions {
		for _, controlID := range assertion.ControlIDs {
			result[controlID] = append(result[controlID], assertion.ID)
		}
	}
	for id := range result {
		sort.Strings(result[id])
	}
	return result
}

func executedByControl(results []model.AssertionResult) map[string][]model.AssertionResult {
	linked := map[string][]model.AssertionResult{}
	for _, result := range results {
		for _, controlID := range result.ControlIDs {
			linked[controlID] = append(linked[controlID], result)
		}
	}
	for id := range linked {
		sort.Slice(linked[id], func(left, right int) bool {
			return linked[id][left].AssertionID < linked[id][right].AssertionID
		})
	}
	return linked
}

func controlResult(control model.Control, contract controlContract, assertionIDs []string, executed []model.AssertionResult) model.ControlResult {
	result := model.ControlResult{
		ControlID: control.ID, Revision: control.Revision, Statement: control.Statement, Source: control.Source,
		ContractSHA256: contract.ContractSHA256, ContractStatus: contract.ContractStatus,
		CanonicalControlID: contract.CanonicalControlID, EvaluationClass: contract.EvaluationClass,
		AutomationClass: contract.AutomationClass, ApplicabilityClass: contract.ApplicabilityClass,
		Atomicity: contract.Atomicity, CompleteInventoryRequired: contract.CompleteInventoryRequired,
		NegativeCondition: contract.NegativeCondition, ProjectThresholdsRequired: contract.ProjectThresholdsRequired,
		EvidenceAuthorities: append([]string{}, contract.EvidenceAuthorities...),
		NotApplicableProof:  contract.NotApplicableProof,
		Disposition:         "needs_review", Coverage: "unmapped", Authority: "none",
		AssertionIDs: append([]string{}, assertionIDs...), ExecutedAssertionIDs: []string{},
		Summary: "No executable assertion currently covers this broad control; scoped evidence and review are still required.",
	}
	if control.Status == "retired" {
		result.Disposition, result.Coverage, result.Summary = "retired", "retired", "This historical control is retired and is not part of the active assessment."
		return result
	}
	if len(assertionIDs) > 0 {
		result.Coverage = "not_in_selected_profile"
		result.Summary = "Executable assertions exist, but none ran in the selected profile; the complete control still needs review."
	}
	if len(executed) == 0 {
		return result
	}
	result.Coverage = "partial_assertions"
	result.Authority = "deterministic_partial"
	counts := map[string]int{}
	for _, assertion := range executed {
		result.ExecutedAssertionIDs = append(result.ExecutedAssertionIDs, assertion.AssertionID)
		counts[assertion.Assessment]++
	}
	switch {
	case counts["fail"] > 0:
		result.Disposition = "confirmed_failure"
		result.Summary = "At least one deterministic assertion failed. This proves a problem, while the rest of the broad control may still need review."
	case counts["unknown"] > 0 || counts["stale"] > 0 || counts["conflicting"] > 0:
		result.Disposition = "blocked"
		result.Summary = "A linked assertion lacks trustworthy evidence or could not run; the broad control remains blocked."
	case counts["manual_review"] > 0:
		result.Disposition = "needs_review"
		result.Summary = "A linked assertion requires accountable human evidence; the broad control remains under review."
	case counts["pass"] == len(executed):
		result.Disposition = "partially_verified"
		result.Summary = "Every linked assertion in this profile passed, but those assertions cover only part of the broad control."
	default:
		result.Disposition = "needs_review"
		result.Summary = "Linked assertions did not establish the complete broad control or a justified Not Applicable result."
	}
	return result
}

func hasIncompleteControls(results []model.ControlResult) bool {
	for _, result := range results {
		if result.Disposition != "retired" && result.Disposition != "confirmed_failure" {
			return true
		}
	}
	return false
}
