// Package fullscan binds the complete source control registry to one scanner
// run. It deliberately keeps broad control dispositions separate from the
// narrower executable assertion results.
package fullscan

import (
	"bufio"
	"bytes"
	"compress/gzip"
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
	"github.com/MarinJursic/production-readiness-checklist/scanner/controlbinding"
	"github.com/MarinJursic/production-readiness-checklist/scanner/controlprogram"
	"github.com/MarinJursic/production-readiness-checklist/scanner/controlprogramcatalog"
	"github.com/MarinJursic/production-readiness-checklist/scanner/controlruntime"
	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
)

const (
	registrySchema          = "prc.control-id-registry/v0.1"
	contractSchema          = "prc.control-contracts/v0.2"
	contractGenerator       = "prc.control-contracts@0.2"
	bindingSchema           = "prc.control-check-bindings/v0.1"
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
	SchemaVersion                     string            `json:"schema_version"`
	GeneratorID                       string            `json:"generator_id"`
	RegistrySHA256                    string            `json:"registry_sha256"`
	RegistryVersion                   string            `json:"registry_version"`
	ClassificationMethodologySHA256   string            `json:"classification_methodology_sha256"`
	ClassificationSummarySHA256       string            `json:"classification_summary_sha256"`
	ClassificationCorpusSHA256        string            `json:"classification_corpus_sha256"`
	ControlCheckBindingsSchemaVersion string            `json:"control_check_bindings_schema_version"`
	ControlCheckBindingsSHA256        string            `json:"control_check_bindings_sha256"`
	BindingCount                      int               `json:"binding_count"`
	ContractCount                     int               `json:"contract_count"`
	Contracts                         []controlContract `json:"contracts"`
}

type controlContract struct {
	ControlID                   string   `json:"control_id"`
	Revision                    int      `json:"revision"`
	ContractStatus              string   `json:"contract_status"`
	ReviewerStatus              string   `json:"reviewer_status"`
	Classification              string   `json:"classification"`
	ClassificationRoute         string   `json:"classification_route"`
	ClassificationDecisionBasis string   `json:"classification_decision_basis"`
	ClassificationRowSHA256     string   `json:"classification_row_sha256"`
	DeterministicBindingID      *string  `json:"deterministic_binding_id"`
	DeterministicBindingSHA256  *string  `json:"deterministic_binding_sha256"`
	CanonicalControlID          string   `json:"canonical_control_id"`
	EvaluationClass             string   `json:"evaluation_class"`
	AutomationClass             string   `json:"automation_class"`
	ApplicabilityClass          string   `json:"applicability_class"`
	Atomicity                   string   `json:"atomicity"`
	CompleteInventoryRequired   bool     `json:"complete_inventory_required"`
	NegativeCondition           bool     `json:"negative_condition"`
	ProjectThresholdsRequired   bool     `json:"project_thresholds_required"`
	EvidenceAuthorities         []string `json:"evidence_authorities"`
	NotApplicableProof          string   `json:"not_applicable_proof"`
	ContractSHA256              string   `json:"contract_sha256"`
}

type bindingDocument struct {
	SchemaVersion              string            `json:"schema_version"`
	GeneratorID                string            `json:"generator_id"`
	RegistryVersion            string            `json:"registry_version"`
	RegistrySHA256             string            `json:"registry_sha256"`
	MethodologySHA256          string            `json:"methodology_sha256"`
	ClassificationCorpusSHA256 string            `json:"classification_corpus_sha256"`
	ImplementationRegistry     []json.RawMessage `json:"implementation_registry"`
	BindingCount               int               `json:"binding_count"`
	Bindings                   []json.RawMessage `json:"bindings"`
}

type bindingIdentity struct {
	ControlID      string `json:"control_id"`
	Revision       int    `json:"revision"`
	SemanticSHA256 string `json:"semantic_sha256"`
	FinalRowSHA256 string `json:"final_row_sha256"`
	Route          string `json:"route"`
}

// Attach loads and verifies the complete registry, adds one result for every
// registered control, and recomputes the content-addressed run identity.
func Attach(root string, scannerCatalog *catalog.Catalog, run model.RunResult) (model.RunResult, error) {
	return AttachProgramExecutions(root, scannerCatalog, run, nil)
}

// AttachProgramExecutions authenticates and attaches exact deterministic
// clause executions before aggregating the complete control registry. Missing
// executions remain Blocked and cannot be upgraded by narrower assertions.
func AttachProgramExecutions(root string, scannerCatalog *catalog.Catalog, run model.RunResult, executions []controlruntime.Execution) (model.RunResult, error) {
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
	programs, err := controlprogramcatalog.Load(root)
	if err != nil {
		return model.RunResult{}, fmt.Errorf("load exact deterministic programs: %w", err)
	}
	if programs.BindingCatalogSHA256() != contracts.ControlCheckBindingsSHA256 ||
		programs.RegistrySHA256() != registryDigest || programs.MethodologySHA256() != contracts.ClassificationMethodologySHA256 ||
		programs.ClassificationCorpusSHA256() != contracts.ClassificationCorpusSHA256 || programs.ControlCount() != contracts.BindingCount {
		return model.RunResult{}, fmt.Errorf("exact deterministic programs are stale against reviewed contracts")
	}
	programResults, deterministicEvidence, err := indexProgramExecutions(programs, executions)
	if err != nil {
		return model.RunResult{}, err
	}
	run.DeterministicEvidence = deterministicEvidence
	profileState := run.TerminalState
	run.ControlCatalog = &model.ControlCatalogSummary{
		SchemaVersion: controlCatalogSchema, RegistryVersion: registry.RegistryVersion,
		RegistrySHA256: registryDigest, SourceSHA256: registry.SourceSHA256,
		ContractSchemaVersion: contracts.SchemaVersion, ContractGeneratorID: contracts.GeneratorID,
		ContractSHA256:                     contractDigest,
		ClassificationMethodologySHA256:    contracts.ClassificationMethodologySHA256,
		ClassificationSummarySHA256:        contracts.ClassificationSummarySHA256,
		ClassificationCorpusSHA256:         contracts.ClassificationCorpusSHA256,
		ControlCheckBindingsSchemaVersion:  contracts.ControlCheckBindingsSchemaVersion,
		ControlCheckBindingsSHA256:         contracts.ControlCheckBindingsSHA256,
		ControlCheckProgramsSchemaVersion:  "prc.control-check-program-catalog/v0.4",
		ControlCheckProgramsSHA256:         programs.Digest(),
		ControlCheckProgramsCatalogSHA256:  programs.CatalogSHA256(),
		ControlCheckDefinitionSchemaSHA256: programs.DefinitionSchemaSHA256(),
		ControlCheckDefinitionCorpusSHA256: programs.DefinitionCorpusSHA256(),
		ControlCount:                       len(registry.Entries), ContractCount: len(contracts.Contracts),
		DeterministicBindingCount:         contracts.BindingCount,
		DeterministicProgramTemplateCount: programs.TemplateCount(),
		DeterministicProgramBlockedCount:  contracts.BindingCount,
		ProfileTerminalState:              profileState,
	}
	run.ControlResults = make([]model.ControlResult, 0, len(registry.Entries))
	active := 0
	allAssertions := assertionsByControl(scannerCatalog)
	executed := executedByControl(run.Results)
	programCountByControl := map[string]int{}
	for _, template := range programs.Templates() {
		programCountByControl[template.ControlID]++
	}
	for index, control := range registry.Entries {
		if control.Status == "active" {
			active++
		}
		contract := contracts.Contracts[index]
		switch contract.Classification {
		case "deterministic":
			run.ControlCatalog.ReviewedDeterministicCount++
		case "nondeterministic":
			run.ControlCatalog.ReviewedNondeterministicCount++
		}
		if contract.ReviewerStatus == "agent_reviewed" {
			run.ControlCatalog.AgentReviewedContractCount++
		} else {
			run.ControlCatalog.GeneratedContractCount++
		}
		run.ControlResults = append(run.ControlResults, controlResult(control, contract, allAssertions[control.ID], executed[control.ID], programCountByControl[control.ID], programResults[control.ID]))
	}
	run.ControlCatalog.ActiveControlCount = active
	run.ControlCatalog.DeterministicProgramBlockedCount = 0
	for _, result := range run.ControlResults {
		if result.Classification == "deterministic" && result.Disposition == "blocked" {
			run.ControlCatalog.DeterministicProgramBlockedCount++
		}
		for _, clause := range result.DeterministicClauseResults {
			switch clause.Status {
			case string(controlruntime.StatusPassed):
				run.ControlCatalog.DeterministicProgramExecutedCount++
				run.ControlCatalog.DeterministicProgramPassCount++
			case string(controlruntime.StatusFailed):
				run.ControlCatalog.DeterministicProgramExecutedCount++
				run.ControlCatalog.DeterministicProgramFailCount++
			case string(controlruntime.StatusNotApplicable):
				run.ControlCatalog.DeterministicProgramExecutedCount++
				run.ControlCatalog.DeterministicProgramNACount++
			case string(controlruntime.StatusBlockedEvidence):
				run.ControlCatalog.DeterministicProgramExecutedCount++
			}
		}
	}
	if hasConfirmedControlFailure(run.ControlResults) {
		run.TerminalState = "no_go"
	} else if run.TerminalState == "profile_satisfied" && hasIncompleteControls(run.ControlResults) {
		run.TerminalState = "assessment_incomplete"
	}
	return Reidentify(run)
}

func indexProgramExecutions(programs *controlprogramcatalog.Catalog, executions []controlruntime.Execution) (map[string][]model.DeterministicClauseResult, []controlprogram.Evidence, error) {
	byTemplate := make(map[string]controlprogramcatalog.Template, programs.TemplateCount())
	for _, template := range programs.Templates() {
		byTemplate[template.TemplateID] = template
	}
	indexed := make(map[string][]model.DeterministicClauseResult)
	deterministicEvidence := make([]controlprogram.Evidence, 0, len(executions))
	seen := make(map[string]bool, len(executions))
	for _, execution := range executions {
		template, ok := byTemplate[execution.TemplateID]
		if !ok || seen[execution.TemplateID] {
			return nil, nil, fmt.Errorf("deterministic program execution references an unknown or duplicate template")
		}
		seen[execution.TemplateID] = true
		if !execution.Authenticated() || execution.ControlID != template.ControlID || execution.ControlRevision != template.ControlRevision ||
			execution.ClauseID != template.ClauseID || execution.ClauseOrdinal != template.ClauseOrdinal ||
			execution.CollectorID != template.CollectorContract.CollectorID || execution.ImplementationID != template.ImplementationID ||
			execution.ImplementationContractSHA256 != template.ImplementationContractSHA256 ||
			execution.RequiredAuthority != template.RequiredAuthority || execution.EvaluatedAt.IsZero() ||
			!validProgramExecution(execution) {
			return nil, nil, fmt.Errorf("deterministic program execution does not match authenticated template %s", template.TemplateID)
		}
		if evidence, exists := execution.SealedEvidence(); exists {
			if controlprogram.EvidenceSHA256(evidence) != execution.EvidenceSHA256 {
				return nil, nil, fmt.Errorf("deterministic program execution evidence does not match template %s", template.TemplateID)
			}
			deterministicEvidence = append(deterministicEvidence, evidence)
		} else if execution.Status == controlruntime.StatusPassed || execution.Status == controlruntime.StatusFailed ||
			execution.Status == controlruntime.StatusNotApplicable || execution.Status == controlruntime.StatusBlockedEvidence {
			return nil, nil, fmt.Errorf("deterministic program execution omits replayable evidence for template %s", template.TemplateID)
		}
		indexed[execution.ControlID] = append(indexed[execution.ControlID], model.DeterministicClauseResult{
			TemplateID: execution.TemplateID, CollectorID: execution.CollectorID, ClauseID: execution.ClauseID,
			ClauseOrdinal: execution.ClauseOrdinal, ImplementationID: execution.ImplementationID,
			ImplementationContractSHA256: execution.ImplementationContractSHA256,
			RequiredAuthority:            string(execution.RequiredAuthority), ProviderID: execution.ProviderID,
			ProgramSHA256: execution.ProgramSHA256, EvidenceSHA256: execution.EvidenceSHA256,
			Status: string(execution.Status), Outcome: string(execution.Outcome), ReasonCode: string(execution.ReasonCode),
			EvaluatedAt: execution.EvaluatedAt.UTC(),
		})
	}
	for controlID := range indexed {
		sort.Slice(indexed[controlID], func(left, right int) bool {
			return indexed[controlID][left].ClauseOrdinal < indexed[controlID][right].ClauseOrdinal
		})
	}
	return indexed, deterministicEvidence, nil
}

func validProgramExecution(execution controlruntime.Execution) bool {
	programDigestValid := execution.ProgramSHA256 == "" || digestPattern.MatchString(execution.ProgramSHA256)
	evidenceDigestValid := execution.EvidenceSHA256 == "" || digestPattern.MatchString(execution.EvidenceSHA256)
	if !programDigestValid || !evidenceDigestValid {
		return false
	}
	switch execution.Status {
	case controlruntime.StatusPassed:
		return execution.ProviderID != "" && execution.ProgramSHA256 != "" && execution.EvidenceSHA256 != "" &&
			execution.Outcome == "pass" && execution.ReasonCode == "passed"
	case controlruntime.StatusFailed:
		return execution.ProviderID != "" && execution.ProgramSHA256 != "" && execution.EvidenceSHA256 != "" &&
			execution.Outcome == "fail" && execution.ReasonCode == "predicate_false"
	case controlruntime.StatusNotApplicable:
		return execution.ProviderID != "" && execution.ProgramSHA256 != "" && execution.EvidenceSHA256 != "" &&
			execution.Outcome == "not_applicable" && execution.ReasonCode == "not_applicable"
	case controlruntime.StatusBlockedEvidence:
		return execution.ProviderID != "" && execution.ProgramSHA256 != "" && execution.EvidenceSHA256 != "" &&
			execution.Outcome == "blocked" && execution.ReasonCode != ""
	case controlruntime.StatusBlockedProvider:
		return execution.ProviderID == "" && execution.ProgramSHA256 != "" && execution.EvidenceSHA256 == "" && execution.Outcome == "" && execution.ReasonCode == ""
	case controlruntime.StatusBlockedBinding:
		return execution.ProviderID == "" && execution.ProgramSHA256 == "" && execution.EvidenceSHA256 == "" && execution.Outcome == "" && execution.ReasonCode == ""
	case controlruntime.StatusBlockedAuthority, controlruntime.StatusBlockedCollection:
		return execution.ProviderID != "" && execution.ProgramSHA256 != "" && execution.EvidenceSHA256 == "" && execution.Outcome == "" && execution.ReasonCode == ""
	case controlruntime.StatusBlockedCanceled:
		return execution.EvidenceSHA256 == "" && execution.Outcome == "" && execution.ReasonCode == ""
	default:
		return false
	}
}

func loadContracts(root string, registry registryDocument, registryDigest string) (contractDocument, string, error) {
	absolute, err := filepath.Abs(root)
	if err != nil {
		return contractDocument{}, "", fmt.Errorf("resolve catalog root: %w", err)
	}
	path := filepath.Join(absolute, "catalog", "control-contracts.json")
	data, err := readCatalogDocument(path, maximumContractBytes, "control contracts")
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
	if err := validateBindingArtifact(absolute, document, registry); err != nil {
		return contractDocument{}, "", err
	}
	digest := sha256.Sum256(data)
	return document, hex.EncodeToString(digest[:]), nil
}

func validateContracts(document contractDocument, registry registryDocument, registryDigest string) error {
	if document.SchemaVersion != contractSchema || document.GeneratorID != contractGenerator ||
		document.RegistrySHA256 != registryDigest ||
		document.RegistryVersion != registry.RegistryVersion || document.ContractCount != len(document.Contracts) ||
		document.ContractCount != len(registry.Entries) || document.ContractCount < 1 || document.ContractCount > maximumControls ||
		!digestPattern.MatchString(document.ClassificationMethodologySHA256) ||
		!digestPattern.MatchString(document.ClassificationSummarySHA256) ||
		!digestPattern.MatchString(document.ClassificationCorpusSHA256) ||
		document.ControlCheckBindingsSchemaVersion != bindingSchema ||
		!digestPattern.MatchString(document.ControlCheckBindingsSHA256) ||
		document.BindingCount < 1 || document.BindingCount > document.ContractCount {
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
	validAuthority := map[string]bool{
		"declared": true, "repository": true, "artifact": true, "executed": true,
		"environment": true, "external_registry": true, "structured_record": true, "human": true,
	}
	deterministic := 0
	seenBindingIDs := map[string]bool{}
	for index, contract := range document.Contracts {
		control := registry.Entries[index]
		if contract.ControlID != control.ID || contract.Revision != control.Revision ||
			!digestPattern.MatchString(contract.ContractSHA256) || !registeredCanonical(contract, control, registered) ||
			contract.ContractStatus != "reviewed" || contract.ReviewerStatus != "agent_reviewed" ||
			!validClassificationContract(contract) || !digestPattern.MatchString(contract.ClassificationRowSHA256) ||
			!validEvaluation[contract.EvaluationClass] || !validAutomation[contract.AutomationClass] ||
			contract.ApplicabilityClass != "scope_required" ||
			(contract.Atomicity != "apparently_atomic" && contract.Atomicity != "compound_review_required") ||
			strings.TrimSpace(contract.NotApplicableProof) == "" || len(contract.NotApplicableProof) > 1000 ||
			len(contract.EvidenceAuthorities) == 0 {
			return fmt.Errorf("control %s has an invalid machine-readable contract", control.ID)
		}
		if control.Status != "active" {
			return fmt.Errorf("reviewed contract %s does not map to an active control", control.ID)
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
		if contract.Classification == "deterministic" {
			deterministic++
			bindingID := valueOrEmpty(contract.DeterministicBindingID)
			if seenBindingIDs[bindingID] {
				return fmt.Errorf("control %s repeats deterministic binding %s", control.ID, bindingID)
			}
			seenBindingIDs[bindingID] = true
		}
	}
	if deterministic != document.BindingCount || len(seenBindingIDs) != document.BindingCount {
		return fmt.Errorf("reviewed deterministic controls do not match the binding count")
	}
	return nil
}

func validClassificationContract(contract controlContract) bool {
	deterministicRoutes := map[string]bool{
		"local_static": true, "artifact_verification": true, "bounded_execution": true,
		"external_readonly_query": true, "structured_record_validation": true, "deterministic_composite": true,
	}
	nondeterministicRoutes := map[string]bool{
		"contextual_judgment": true, "accountable_human_decision": true,
		"specialist_or_legal_judgment": true, "empirical_protocol_undefined": true,
		"contract_incomplete": true, "mixed": true, "unbounded_claim": true,
	}
	switch contract.Classification {
	case "deterministic":
		return deterministicRoutes[contract.ClassificationRoute] &&
			contract.ClassificationDecisionBasis == "strength_audit_confirmed" &&
			valueOrEmpty(contract.DeterministicBindingID) == fmt.Sprintf("%s@%d", contract.ControlID, contract.Revision) &&
			digestPattern.MatchString(valueOrEmpty(contract.DeterministicBindingSHA256))
	case "nondeterministic":
		return nondeterministicRoutes[contract.ClassificationRoute] &&
			(contract.ClassificationDecisionBasis == "primary_nondeterministic" ||
				contract.ClassificationDecisionBasis == "skeptically_rejected" ||
				contract.ClassificationDecisionBasis == "strength_audit_reclassified") &&
			contract.DeterministicBindingID == nil && contract.DeterministicBindingSHA256 == nil
	default:
		return false
	}
}

func valueOrEmpty(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func validateBindingArtifact(root string, contracts contractDocument, registry registryDocument) error {
	loaded, err := controlbinding.Load(root)
	if err != nil {
		return fmt.Errorf("load strict control check bindings: %w", err)
	}
	if loaded.Digest() != contracts.ControlCheckBindingsSHA256 ||
		loaded.RegistryVersion() != registry.RegistryVersion || loaded.RegistrySHA256() != contracts.RegistrySHA256 ||
		loaded.MethodologySHA256() != contracts.ClassificationMethodologySHA256 ||
		loaded.ClassificationCorpusSHA256() != contracts.ClassificationCorpusSHA256 ||
		loaded.BindingCount() != contracts.BindingCount || len(loaded.Implementations()) != 11 || loaded.ClauseCount() != 765 {
		return fmt.Errorf("strict control check bindings are stale against the reviewed contract envelope")
	}
	path := filepath.Join(root, "catalog", "control-check-bindings.json")
	data, err := readCatalogDocument(path, maximumContractBytes, "control check bindings")
	if err != nil {
		return err
	}
	digest := sha256.Sum256(data)
	if hex.EncodeToString(digest[:]) != contracts.ControlCheckBindingsSHA256 {
		return fmt.Errorf("control check bindings digest does not match the reviewed contract envelope")
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	var document bindingDocument
	if err := decoder.Decode(&document); err != nil {
		return fmt.Errorf("parse control check bindings: %w", err)
	}
	var trailing any
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		return fmt.Errorf("control check bindings contain trailing JSON")
	}
	if document.SchemaVersion != bindingSchema || document.GeneratorID != "prc.build-control-check-bindings@0.1" ||
		document.RegistryVersion != registry.RegistryVersion || document.RegistrySHA256 != contracts.RegistrySHA256 ||
		document.MethodologySHA256 != contracts.ClassificationMethodologySHA256 ||
		document.ClassificationCorpusSHA256 != contracts.ClassificationCorpusSHA256 ||
		document.BindingCount != contracts.BindingCount || document.BindingCount != len(document.Bindings) ||
		len(document.ImplementationRegistry) != 11 {
		return fmt.Errorf("control check bindings are stale against the reviewed contract envelope")
	}

	registered := make(map[string]model.Control, len(registry.Entries))
	for _, control := range registry.Entries {
		registered[control.ID] = control
	}
	contractByID := make(map[string]controlContract, len(contracts.Contracts))
	for _, contract := range contracts.Contracts {
		contractByID[contract.ControlID] = contract
	}
	seen := map[string]bool{}
	previous := ""
	expectedFields := map[string]bool{
		"control_id": true, "revision": true, "semantic_sha256": true, "final_row_sha256": true,
		"route": true, "aggregation": true, "applicability_contract": true, "clauses": true,
	}
	for _, raw := range document.Bindings {
		var fields map[string]json.RawMessage
		if err := json.Unmarshal(raw, &fields); err != nil || len(fields) != len(expectedFields) {
			return fmt.Errorf("control check binding contains invalid fields")
		}
		for field := range fields {
			if !expectedFields[field] {
				return fmt.Errorf("control check binding contains unsupported field %q", field)
			}
		}
		var identity bindingIdentity
		if err := json.Unmarshal(raw, &identity); err != nil {
			return fmt.Errorf("parse control check binding identity: %w", err)
		}
		control, registeredControl := registered[identity.ControlID]
		contract, contracted := contractByID[identity.ControlID]
		bindingDigest, err := canonicalJSONDigest(raw)
		if err != nil {
			return fmt.Errorf("digest control check binding %s: %w", identity.ControlID, err)
		}
		if !registeredControl || !contracted || contract.Classification != "deterministic" ||
			seen[identity.ControlID] || previous >= identity.ControlID || identity.Revision != control.Revision ||
			identity.SemanticSHA256 != control.SemanticSHA256 ||
			identity.FinalRowSHA256 != contract.ClassificationRowSHA256 ||
			identity.Route != contract.ClassificationRoute ||
			valueOrEmpty(contract.DeterministicBindingID) != fmt.Sprintf("%s@%d", identity.ControlID, identity.Revision) ||
			valueOrEmpty(contract.DeterministicBindingSHA256) != bindingDigest {
			return fmt.Errorf("control check binding %s is stale or mismatched", identity.ControlID)
		}
		seen[identity.ControlID], previous = true, identity.ControlID
	}
	if len(seen) != contracts.BindingCount {
		return fmt.Errorf("control check binding identities do not match the reviewed deterministic controls")
	}
	return nil
}

func canonicalJSONDigest(raw json.RawMessage) (string, error) {
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.UseNumber()
	var value any
	if err := decoder.Decode(&value); err != nil {
		return "", err
	}
	var trailing any
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		return "", fmt.Errorf("value contains trailing JSON")
	}
	var canonical bytes.Buffer
	encoder := json.NewEncoder(&canonical)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(value); err != nil {
		return "", err
	}
	encoded := bytes.TrimSuffix(canonical.Bytes(), []byte{'\n'})
	digest := sha256.Sum256(encoded)
	return hex.EncodeToString(digest[:]), nil
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
		"contract_status": contract.ContractStatus, "reviewer_status": contract.ReviewerStatus,
		"classification": contract.Classification, "classification_route": contract.ClassificationRoute,
		"classification_decision_basis": contract.ClassificationDecisionBasis,
		"classification_row_sha256":     contract.ClassificationRowSHA256,
		"deterministic_binding_id":      contract.DeterministicBindingID,
		"deterministic_binding_sha256":  contract.DeterministicBindingSHA256,
		"canonical_control_id":          contract.CanonicalControlID,
		"evaluation_class":              contract.EvaluationClass, "automation_class": contract.AutomationClass,
		"applicability_class": contract.ApplicabilityClass, "atomicity": contract.Atomicity,
		"complete_inventory_required": contract.CompleteInventoryRequired,
		"negative_condition":          contract.NegativeCondition,
		"project_thresholds_required": contract.ProjectThresholdsRequired,
		"evidence_authorities":        contract.EvidenceAuthorities, "not_applicable_proof": contract.NotApplicableProof,
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
	data, err := readCatalogDocument(path, maximumRegistryBytes, "control registry")
	if errors.Is(err, os.ErrNotExist) {
		return registryDocument{}, "", fmt.Errorf("%w: %s or %s.gz", ErrRegistryUnavailable, path, path)
	}
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

// readCatalogDocument accepts the ordinary source JSON or the deterministic
// gzip form used only by the compact npm runtime package. Exactly one form may
// exist. Both the compressed input and the expanded JSON are bounded, and gzip
// checks its checksum before the document is trusted by the existing schema
// and content-digest validation.
func readCatalogDocument(path string, limit int64, label string) ([]byte, error) {
	_, plainErr := os.Lstat(path)
	compressedPath := path + ".gz"
	_, compressedErr := os.Lstat(compressedPath)
	plainExists := plainErr == nil
	compressedExists := compressedErr == nil
	if plainExists && compressedExists {
		return nil, fmt.Errorf("%s has ambiguous plain and compressed files", label)
	}
	if plainExists {
		return readRegularBounded(path, limit, label)
	}
	if compressedExists {
		encoded, err := readRegularBounded(compressedPath, limit, label+" compressed data")
		if err != nil {
			return nil, err
		}
		reader, err := gzip.NewReader(bytes.NewReader(encoded))
		if err != nil {
			return nil, fmt.Errorf("open compressed %s %s: %w", label, compressedPath, err)
		}
		data, readErr := io.ReadAll(io.LimitReader(reader, limit+1))
		closeErr := reader.Close()
		if readErr != nil {
			return nil, fmt.Errorf("read compressed %s %s: %w", label, compressedPath, readErr)
		}
		if closeErr != nil {
			return nil, fmt.Errorf("close compressed %s %s: %w", label, compressedPath, closeErr)
		}
		if int64(len(data)) > limit {
			return nil, fmt.Errorf("expanded %s %s exceeds its %d-byte limit", label, compressedPath, limit)
		}
		return data, nil
	}
	if plainErr != nil && !os.IsNotExist(plainErr) {
		return nil, fmt.Errorf("inspect %s %s: %w", label, path, plainErr)
	}
	if compressedErr != nil && !os.IsNotExist(compressedErr) {
		return nil, fmt.Errorf("inspect compressed %s %s: %w", label, compressedPath, compressedErr)
	}
	return nil, fmt.Errorf("%w: %s or %s", os.ErrNotExist, path, compressedPath)
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

func controlResult(control model.Control, contract controlContract, assertionIDs []string, executed []model.AssertionResult, programTemplateCount int, programResults []model.DeterministicClauseResult) model.ControlResult {
	result := model.ControlResult{
		ControlID: control.ID, Revision: control.Revision, Statement: control.Statement, Source: control.Source,
		ContractSHA256: contract.ContractSHA256, ContractStatus: contract.ContractStatus,
		Classification: contract.Classification, ClassificationRoute: contract.ClassificationRoute,
		ClassificationDecisionBasis:       contract.ClassificationDecisionBasis,
		ClassificationRowSHA256:           contract.ClassificationRowSHA256,
		DeterministicBindingID:            valueOrEmpty(contract.DeterministicBindingID),
		DeterministicBindingSHA256:        valueOrEmpty(contract.DeterministicBindingSHA256),
		DeterministicProgramTemplateCount: programTemplateCount,
		DeterministicClauseResults:        append([]model.DeterministicClauseResult(nil), programResults...),
		CanonicalControlID:                contract.CanonicalControlID, EvaluationClass: contract.EvaluationClass,
		AutomationClass: contract.AutomationClass, ApplicabilityClass: contract.ApplicabilityClass,
		Atomicity: contract.Atomicity, CompleteInventoryRequired: contract.CompleteInventoryRequired,
		NegativeCondition: contract.NegativeCondition, ProjectThresholdsRequired: contract.ProjectThresholdsRequired,
		EvidenceAuthorities: append([]string{}, contract.EvidenceAuthorities...),
		NotApplicableProof:  contract.NotApplicableProof,
		Disposition:         "needs_review", Coverage: "nondeterministic_advisory", Authority: "none",
		AssertionIDs: append([]string{}, assertionIDs...), ExecutedAssertionIDs: []string{},
		Summary: "This reviewed control requires contextual or accountable review. Local checks and AI can provide advice, but cannot produce a verified pass.",
	}
	if control.Status == "retired" {
		result.Disposition, result.Coverage, result.Summary = "retired", "retired", "This historical control is retired and is not part of the active assessment."
		return result
	}
	counts := map[string]int{}
	for _, assertion := range executed {
		result.ExecutedAssertionIDs = append(result.ExecutedAssertionIDs, assertion.AssertionID)
		counts[assertion.Assessment]++
	}
	if contract.Classification == "deterministic" {
		if programTemplateCount == 0 {
			result.Summary = "This deterministic control has no exact program template, so the catalog is incomplete."
			return result
		}
		result.Disposition = "blocked"
		result.Coverage = "deterministic_program_provider_unregistered"
		result.DeterministicProgramStatus = "blocked_provider_unregistered"
		result.Summary = fmt.Sprintf("Exact deterministic programs exist for all %d bound clause(s), but an authoritative collector is not registered, so this control is Blocked.", programTemplateCount)
		if len(executed) > 0 {
			result.Authority = "deterministic_partial"
		}
		programCounts := map[string]int{}
		for _, clause := range programResults {
			programCounts[clause.Status]++
		}
		completeProgramSet := len(programResults) == programTemplateCount
		evaluatedClauses := programCounts[string(controlruntime.StatusPassed)] + programCounts[string(controlruntime.StatusFailed)] +
			programCounts[string(controlruntime.StatusNotApplicable)] + programCounts[string(controlruntime.StatusBlockedEvidence)]
		if len(programResults) > 0 {
			result.Authority = "deterministic_exact"
			result.Coverage = "deterministic_program_partial"
			result.DeterministicProgramStatus = "blocked_program_incomplete"
			result.Summary = fmt.Sprintf("Exact deterministic evidence was supplied for %d of %d clause(s), but the complete control remains Blocked.", evaluatedClauses, programTemplateCount)
		}
		switch {
		case programCounts[string(controlruntime.StatusFailed)] > 0:
			result.Disposition = "confirmed_failure"
			result.DeterministicProgramStatus = "executed_fail"
			if completeProgramSet {
				result.Coverage = "deterministic_program_complete"
			}
			result.Summary = "An exact deterministic clause failed with sealed authoritative evidence. The control is a confirmed failure."
			return result
		case completeProgramSet && programCounts[string(controlruntime.StatusPassed)]+programCounts[string(controlruntime.StatusNotApplicable)] == programTemplateCount && programCounts[string(controlruntime.StatusPassed)] > 0:
			result.Disposition = "verified_pass"
			result.Coverage = "deterministic_program_complete"
			result.DeterministicProgramStatus = "executed_pass"
			result.Summary = fmt.Sprintf("All %d exact deterministic clause program(s) passed or were proven Not Applicable with sealed authoritative evidence.", programTemplateCount)
			return result
		case completeProgramSet && programCounts[string(controlruntime.StatusNotApplicable)] == programTemplateCount:
			result.Disposition = "not_applicable"
			result.Coverage = "deterministic_program_complete"
			result.DeterministicProgramStatus = "executed_not_applicable"
			result.Summary = "Every exact deterministic clause was proven Not Applicable for the sealed assessment scope."
			return result
		case completeProgramSet:
			result.Coverage = "deterministic_program_complete"
			result.DeterministicProgramStatus = "blocked_evidence"
			result.Summary = "Every exact deterministic clause was attempted, but at least one lacked complete, current, non-conflicting authoritative evidence."
		}
		if len(programResults) == 0 {
			switch {
			case counts["fail"] > 0:
				result.Summary = "A narrow assertion found a problem, but the complete reviewed deterministic binding was not executed; the control remains blocked."
			case counts["unknown"] > 0 || counts["stale"] > 0 || counts["conflicting"] > 0:
				result.Summary = "Narrow assertion evidence is unresolved and the complete reviewed deterministic binding was not executed; the control remains blocked."
			case len(executed) > 0 && counts["pass"] == len(executed):
				result.Summary = "Every selected narrow assertion passed, but the complete reviewed deterministic binding was not executed; the control remains blocked."
			}
		}
		return result
	}
	if len(executed) > 0 {
		result.Summary = "Local assertion evidence is available, but this control is reviewed as nondeterministic and still needs contextual or accountable review. AI output remains advisory."
	}
	return result
}

func hasIncompleteControls(results []model.ControlResult) bool {
	for _, result := range results {
		if result.Disposition != "retired" && result.Disposition != "confirmed_failure" && result.Disposition != "verified_pass" && result.Disposition != "not_applicable" {
			return true
		}
	}
	return false
}

func hasConfirmedControlFailure(results []model.ControlResult) bool {
	for _, result := range results {
		if result.Disposition == "confirmed_failure" {
			return true
		}
	}
	return false
}
