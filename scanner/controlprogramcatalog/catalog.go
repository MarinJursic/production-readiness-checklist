// Package controlprogramcatalog authenticates the generated deterministic
// clause programs. It never obtains evidence and never executes target code.
package controlprogramcatalog

import (
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
	"reflect"
	"regexp"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/MarinJursic/production-readiness-checklist/scanner/controlbinding"
	"github.com/MarinJursic/production-readiness-checklist/scanner/controlprogram"
	"github.com/MarinJursic/production-readiness-checklist/scanner/providercapability"
)

const (
	catalogRelativePath       = "catalog/control-check-programs.json"
	programSchemaRelativePath = "schemas/control-check-program.schema.json"
	schemaVersion             = "prc.control-check-program-catalog/v0.4"
	generatorID               = "prc.build-control-check-programs@0.4"
	bindingSchemaVersion      = "prc.control-check-bindings/v0.1"

	expectedControlCount  = 686
	expectedTemplateCount = 765
	maximumClauseCount    = 50
	maximumStatementRunes = 2000
	maximumJSONDepth      = 32

	maximumPlainBytes      int64 = 16 * 1024 * 1024
	maximumCompressedBytes int64 = 4 * 1024 * 1024
	maximumExpandedBytes   int64 = maximumPlainBytes
	maximumSchemaBytes     int64 = 256 * 1024
)

var (
	digestPattern         = regexp.MustCompile(`^[0-9a-f]{64}$`)
	versionPattern        = regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+$`)
	controlIDPattern      = regexp.MustCompile(`^(?:PRC-[0-9]{2}-[0-9]{3}|USEQ-[A-F0-9]{8})$`)
	implementationPattern = regexp.MustCompile(`^prc\.check\.[a-z0-9-]+@[0-9]+\.[0-9]+$`)
	collectorPattern      = regexp.MustCompile(`^prc\.collect\.[a-z0-9.-]+@[0-9]+\.[0-9]+$`)
	factKeyPattern        = regexp.MustCompile(`^[A-Za-z0-9_.-]{3,128}$`)
)

var validFamilies = map[string]bool{
	"inventory_fact": true, "structured_document": true, "package_metadata": true,
	"ci_policy": true, "container_iac": true, "source_ast": true,
	"artifact_integrity": true, "analysis_adapter": true, "execution_evidence": true,
	"environment_evidence": true, "structured_record": true,
}

var validAuthorities = map[string]controlprogram.Authority{
	"repository": controlprogram.AuthorityRepository, "artifact": controlprogram.AuthorityArtifact,
	"executed": controlprogram.AuthorityExecuted, "environment": controlprogram.AuthorityEnvironment,
	"external_registry": controlprogram.AuthorityExternalRegistry,
	"structured_record": controlprogram.AuthorityStructuredRecord,
}

var validFactTypes = map[controlprogram.FactType]bool{
	controlprogram.FactIdentity: true, controlprogram.FactSchema: true,
	controlprogram.FactDigest: true, controlprogram.FactState: true,
	controlprogram.FactString: true, controlprogram.FactBoolean: true,
	controlprogram.FactNumber: true, controlprogram.FactTime: true,
	controlprogram.FactStringSet: true, controlprogram.FactIdentityMap: true,
	controlprogram.FactSchemaMap: true, controlprogram.FactDigestMap: true,
	controlprogram.FactStateMap: true, controlprogram.FactStringMap: true,
	controlprogram.FactBooleanMap: true, controlprogram.FactNumberMap: true,
	controlprogram.FactTimeMap: true, controlprogram.FactDirectedGraph: true,
}

var opaqueVerdictNames = map[string]bool{
	"clause_satisfied": true, "is_compliant": true, "compliant": true,
	"control_passed": true, "provider_verdict": true, "verified": true,
	"valid": true, "passed": true, "result": true,
}

var forbiddenGenericKeys = map[string]bool{
	"record-paths": true, "required-paths": true, "record_paths": true, "required_paths": true,
}

// FactContract describes one raw normalized value required by a predicate.
type FactContract struct {
	FactKey           string
	FactType          controlprogram.FactType
	Authority         controlprogram.Authority
	RawValueSemantics string
	SourceRequirement string
	CompleteRequired  bool
}

// ParameterContract describes one expected value sealed before evidence is requested.
type ParameterContract struct {
	ParameterKey      string
	ParameterType     controlprogram.FactType
	ValueOrigin       ParameterOrigin
	SourceRequirement string
}

// ParameterOrigin keeps inventories, project policy, and assessed runtime
// context in separate trust lanes. Absolute invariants belong in predicates,
// not in a caller-selectable parameter lane.
type ParameterOrigin string

const (
	ParameterOriginScannerInventory ParameterOrigin = "scanner_inventory"
	ParameterOriginPolicy           ParameterOrigin = "independently_authenticated_policy"
	ParameterOriginContext          ParameterOrigin = "independently_authenticated_context"
)

// CollectorContract documents the unregistered authority-specific collector.
type CollectorContract struct {
	CollectorID           string
	RequiredSources       []string
	InventoryContract     string
	NormalizationContract string
	CompletenessContract  string
	FreshnessContract     string
	ProviderStatus        string
}

// RuntimeRequirements records fail-closed runtime inputs and capabilities.
type RuntimeRequirements struct {
	SubjectID, Subjects, InventorySHA256, MaximumEvidenceAgeSeconds        string
	AllowNotApplicable, ApplicabilityProofContractSHA256                   string
	SealedParameters, SealedParametersBoundBy                              string
	EvidenceProviderMaySupplyParameters                                    bool
	ProviderRegistration, DomainEvidenceCollector, MissingCapabilityResult string
	ProviderClaimed                                                        bool
}

// RuntimeBinding supplies scanner-sealed scope, policy parameters, and freshness.
type RuntimeBinding struct {
	SubjectID                        string
	Subjects                         []string
	InventorySHA256                  string
	AllowNotApplicable               bool
	ApplicabilityProofContractSHA256 string
	MaximumEvidenceAgeSeconds        int64
	ScannerInventoryParameters       map[string]controlprogram.Parameter
	AuthenticatedPolicyParameters    map[string]controlprogram.Parameter
	AuthenticatedContextParameters   map[string]controlprogram.Parameter
}

// Template is one immutable clause-bound predicate template.
type Template struct {
	TemplateID, ProgramSchemaVersion, ProgramSchemaSHA256 string
	ControlID                                             string
	ControlRevision                                       int
	ControlSemanticSHA256                                 string
	ClauseOrdinal                                         int
	ClauseID, ClauseStatement, ClauseStatementSHA256      string
	CheckerFamily                                         string
	RequiredAuthority                                     controlprogram.Authority
	ImplementationID, ImplementationContractSHA256        string
	PredicateShape, ReviewReason, CounterexampleAnalysis  string
	EndToEndRunnable                                      bool
	RawFactContracts                                      []FactContract
	SealedParameterContracts                              []ParameterContract
	Predicate                                             controlprogram.Expression
	RequiredRuntimeOps                                    []controlprogram.Operation
	CollectorContract                                     CollectorContract
	RuntimeRequirements                                   RuntimeRequirements
	TemplateSHA256                                        string
}

// Catalog is an immutable validated snapshot of all deterministic templates.
type Catalog struct {
	digest, catalogSHA256, programSchemaSHA256, bindingCatalogSHA256               string
	definitionSchemaSHA256, definitionCorpusSHA256                                 string
	registryVersion, registrySHA256, methodologySHA256, classificationCorpusSHA256 string
	templates                                                                      []Template
	byClause                                                                       map[string]int
}

type rawDocument struct {
	SchemaVersion                  string            `json:"schema_version"`
	GeneratorID                    string            `json:"generator_id"`
	ProgramSchemaVersion           string            `json:"program_schema_version"`
	ProgramSchemaSHA256            string            `json:"program_schema_sha256"`
	BindingSchemaVersion           string            `json:"binding_schema_version"`
	BindingCatalogSHA256           string            `json:"binding_catalog_sha256"`
	DefinitionSchemaSHA256         string            `json:"definition_schema_sha256"`
	DefinitionCorpusSHA256         string            `json:"definition_corpus_sha256"`
	RegistryVersion                string            `json:"registry_version"`
	RegistrySHA256                 string            `json:"registry_sha256"`
	MethodologySHA256              string            `json:"methodology_sha256"`
	ClassificationCorpusSHA256     string            `json:"classification_corpus_sha256"`
	ControlCount                   int               `json:"control_count"`
	TemplateCount                  int               `json:"template_count"`
	PredicateDefinedCount          int               `json:"predicate_defined_count"`
	ImplementationMissingCount     int               `json:"implementation_missing_count"`
	ProviderCapabilityMissingCount int               `json:"provider_capability_missing_count"`
	EndToEndRunnableTemplateCount  int               `json:"end_to_end_runnable_template_count"`
	EndToEndRunnableControlCount   int               `json:"end_to_end_runnable_control_count"`
	BlockedControlCount            int               `json:"blocked_control_count"`
	ClassificationErrorCount       int               `json:"classification_error_count"`
	PredicateShapeCounts           map[string]int    `json:"predicate_shape_counts"`
	Templates                      []json.RawMessage `json:"templates"`
	CatalogSHA256                  string            `json:"catalog_sha256"`
}

type rawFactContract struct {
	FactKey           string                  `json:"fact_id"`
	FactType          controlprogram.FactType `json:"fact_type"`
	Authority         string                  `json:"authority"`
	RawValueSemantics string                  `json:"raw_value_semantics"`
	SourceRequirement string                  `json:"source_requirement"`
	CompleteRequired  *bool                   `json:"complete_required"`
}
type rawParameterContract struct {
	ParameterKey      string                  `json:"parameter_id"`
	ParameterType     controlprogram.FactType `json:"parameter_type"`
	ValueOrigin       string                  `json:"value_origin"`
	SourceRequirement string                  `json:"source_requirement"`
}
type rawCollectorContract struct {
	CollectorID           string   `json:"collector_id"`
	RequiredSources       []string `json:"required_sources"`
	InventoryContract     string   `json:"inventory_contract"`
	NormalizationContract string   `json:"normalization_contract"`
	CompletenessContract  string   `json:"completeness_contract"`
	FreshnessContract     string   `json:"freshness_contract"`
	ProviderStatus        string   `json:"provider_status"`
}
type rawRuntimeRequirements struct {
	SubjectID                           string `json:"subject_id"`
	Subjects                            string `json:"subjects"`
	InventorySHA256                     string `json:"inventory_sha256"`
	MaximumEvidenceAgeSeconds           string `json:"maximum_evidence_age_seconds"`
	AllowNotApplicable                  string `json:"allow_not_applicable"`
	ApplicabilityProofContractSHA256    string `json:"applicability_proof_contract_sha256"`
	SealedParameters                    string `json:"sealed_parameters"`
	SealedParametersBoundBy             string `json:"sealed_parameters_bound_by"`
	EvidenceProviderMaySupplyParameters *bool  `json:"evidence_provider_may_supply_parameters"`
	ProviderRegistration                string `json:"provider_registration"`
	ProviderClaimed                     *bool  `json:"provider_claimed"`
	DomainEvidenceCollector             string `json:"domain_evidence_collector"`
	MissingCapabilityResult             string `json:"missing_capability_result"`
}
type rawTemplate struct {
	TemplateID                   string                     `json:"template_id"`
	ProgramSchemaVersion         string                     `json:"program_schema_version"`
	ProgramSchemaSHA256          string                     `json:"program_schema_sha256"`
	ControlID                    string                     `json:"control_id"`
	ControlRevision              int                        `json:"control_revision"`
	ControlSemanticSHA256        string                     `json:"control_semantic_sha256"`
	ClauseOrdinal                int                        `json:"clause_ordinal"`
	ClauseID                     string                     `json:"clause_id"`
	ClauseStatement              string                     `json:"clause_statement"`
	ClauseStatementSHA256        string                     `json:"clause_statement_sha256"`
	CheckerFamily                string                     `json:"checker_family"`
	RequiredAuthority            string                     `json:"required_authority"`
	ImplementationID             string                     `json:"implementation_id"`
	ImplementationContractSHA256 string                     `json:"implementation_contract_sha256"`
	ReviewStatus                 string                     `json:"review_status"`
	PredicateDefined             *bool                      `json:"predicate_defined"`
	EndToEndRunnable             *bool                      `json:"end_to_end_runnable"`
	PredicateShape               string                     `json:"predicate_shape"`
	ReviewReason                 string                     `json:"review_reason"`
	CounterexampleAnalysis       string                     `json:"counterexample_analysis"`
	RawFactContracts             []rawFactContract          `json:"raw_fact_contracts"`
	SealedParameterContracts     []rawParameterContract     `json:"sealed_parameter_contracts"`
	Predicate                    json.RawMessage            `json:"predicate"`
	RequiredRuntimeOps           []controlprogram.Operation `json:"required_runtime_ops"`
	CollectorContract            rawCollectorContract       `json:"collector_contract"`
	ProviderCapabilityStatus     string                     `json:"provider_capability_status"`
	RuntimeRequirements          rawRuntimeRequirements     `json:"runtime_requirements"`
	TemplateSHA256               string                     `json:"template_sha256"`
}

// Load reads a plain or gzip catalog and authenticates it against bindings and schema.
func Load(root string) (*Catalog, error) {
	absolute, err := filepath.Abs(root)
	if err != nil {
		return nil, fmt.Errorf("resolve program catalog root: %w", err)
	}
	data, err := readCatalogDocument(filepath.Join(absolute, filepath.FromSlash(catalogRelativePath)))
	if err != nil {
		return nil, err
	}
	if !utf8.Valid(data) {
		return nil, fmt.Errorf("control check programs are not valid UTF-8")
	}
	schemaData, err := readRegularBounded(filepath.Join(absolute, filepath.FromSlash(programSchemaRelativePath)), maximumSchemaBytes, "control check program schema")
	if err != nil {
		return nil, err
	}
	bindings, err := controlbinding.Load(absolute)
	if err != nil {
		return nil, fmt.Errorf("load referenced control check bindings: %w", err)
	}
	catalog, err := decodeAndValidate(data, sha256Hex(schemaData), bindings)
	if err != nil {
		return nil, err
	}
	catalog.digest = sha256Hex(data)
	return catalog, nil
}

func (catalog *Catalog) Digest() string                 { return catalog.digest }
func (catalog *Catalog) CatalogSHA256() string          { return catalog.catalogSHA256 }
func (catalog *Catalog) ProgramSchemaSHA256() string    { return catalog.programSchemaSHA256 }
func (catalog *Catalog) BindingCatalogSHA256() string   { return catalog.bindingCatalogSHA256 }
func (catalog *Catalog) DefinitionSchemaSHA256() string { return catalog.definitionSchemaSHA256 }
func (catalog *Catalog) DefinitionCorpusSHA256() string { return catalog.definitionCorpusSHA256 }
func (catalog *Catalog) RegistryVersion() string        { return catalog.registryVersion }
func (catalog *Catalog) RegistrySHA256() string         { return catalog.registrySHA256 }
func (catalog *Catalog) MethodologySHA256() string      { return catalog.methodologySHA256 }
func (catalog *Catalog) ClassificationCorpusSHA256() string {
	return catalog.classificationCorpusSHA256
}
func (catalog *Catalog) ControlCount() int  { return expectedControlCount }
func (catalog *Catalog) TemplateCount() int { return len(catalog.templates) }

func (catalog *Catalog) Templates() []Template {
	result := make([]Template, len(catalog.templates))
	for index, template := range catalog.templates {
		result[index] = cloneTemplate(template)
	}
	return result
}

func (catalog *Catalog) Lookup(controlID, clauseID string) (Template, bool) {
	index, ok := catalog.byClause[lookupKey(controlID, clauseID)]
	if !ok {
		return Template{}, false
	}
	return cloneTemplate(catalog.templates[index]), true
}

// Program binds scanner-sealed runtime values and revalidates the exact program.
func (template Template) Program(runtime RuntimeBinding) (controlprogram.Program, error) {
	parameters, ok := parametersByContractOrigin(template.SealedParameterContracts, runtime)
	if !ok {
		return controlprogram.Program{}, fmt.Errorf("bind program template %s: sealed parameters do not match contracts", template.TemplateID)
	}
	program := controlprogram.Program{
		SchemaVersion: controlprogram.ProgramSchemaVersion, ControlID: template.ControlID,
		ControlRevision: template.ControlRevision, ControlSemanticSHA256: template.ControlSemanticSHA256,
		ClauseID: template.ClauseID, ClauseSHA256: template.ClauseStatementSHA256,
		ImplementationContractSHA256: template.ImplementationContractSHA256,
		SubjectID:                    runtime.SubjectID, Subjects: append([]string(nil), runtime.Subjects...),
		InventorySHA256: runtime.InventorySHA256, RequiredAuthority: template.RequiredAuthority,
		AllowNotApplicable:               runtime.AllowNotApplicable,
		ApplicabilityProofContractSHA256: runtime.ApplicabilityProofContractSHA256,
		MaximumEvidenceAgeSeconds:        runtime.MaximumEvidenceAgeSeconds,
		Parameters:                       parameters, Predicate: cloneExpression(template.Predicate),
	}
	if err := controlprogram.ValidateProgram(program); err != nil {
		return controlprogram.Program{}, fmt.Errorf("bind program template %s: %w", template.TemplateID, err)
	}
	return program, nil
}

// ValidateMaterializedProgram verifies that a program supplied through an
// authenticated policy channel is still exactly the reviewed template. The
// policy channel may choose scope, freshness, applicability, and parameter
// values, but it cannot replace the predicate, change a parameter type, move
// the clause to another authority, or alter any reviewed identity.
func (template Template) ValidateMaterializedProgram(program controlprogram.Program) error {
	if err := controlprogram.ValidateProgram(program); err != nil {
		return fmt.Errorf("validate materialized program %s: %w", template.TemplateID, err)
	}
	if program.SchemaVersion != template.ProgramSchemaVersion ||
		program.ControlID != template.ControlID || program.ControlRevision != template.ControlRevision ||
		program.ControlSemanticSHA256 != template.ControlSemanticSHA256 ||
		program.ClauseID != template.ClauseID || program.ClauseSHA256 != template.ClauseStatementSHA256 ||
		program.ImplementationContractSHA256 != template.ImplementationContractSHA256 ||
		program.RequiredAuthority != template.RequiredAuthority ||
		!reflect.DeepEqual(program.Predicate, template.Predicate) {
		return fmt.Errorf("materialized program does not match reviewed template %s", template.TemplateID)
	}
	if len(program.Parameters) != len(template.SealedParameterContracts) {
		return fmt.Errorf("materialized program parameters do not match reviewed template %s", template.TemplateID)
	}
	for _, contract := range template.SealedParameterContracts {
		parameter, ok := program.Parameters[contract.ParameterKey]
		if !ok || parameter.Type != contract.ParameterType {
			return fmt.Errorf("materialized program parameter %s does not match reviewed template %s", contract.ParameterKey, template.TemplateID)
		}
	}
	return nil
}

func decodeAndValidate(data []byte, schemaDigest string, bindings *controlbinding.Catalog) (*Catalog, error) {
	if err := rejectDuplicateKeys(data); err != nil {
		return nil, fmt.Errorf("parse control check programs: %w", err)
	}
	var document rawDocument
	if err := decodeStrict(data, &document); err != nil {
		return nil, fmt.Errorf("parse control check programs: %w", err)
	}
	if document.SchemaVersion != schemaVersion || document.GeneratorID != generatorID || document.ProgramSchemaVersion != controlprogram.ProgramSchemaVersion || document.BindingSchemaVersion != bindingSchemaVersion {
		return nil, fmt.Errorf("unsupported control check program catalog envelope")
	}
	if !isDigest(schemaDigest) || document.ProgramSchemaSHA256 != schemaDigest || !validDocumentDigests(document) {
		return nil, fmt.Errorf("control check program catalog has invalid or stale envelope digests")
	}
	if document.BindingCatalogSHA256 != bindings.Digest() || document.RegistryVersion != bindings.RegistryVersion() ||
		document.RegistrySHA256 != bindings.RegistrySHA256() || document.MethodologySHA256 != bindings.MethodologySHA256() ||
		document.ClassificationCorpusSHA256 != bindings.ClassificationCorpusSHA256() ||
		bindings.BindingCount() != expectedControlCount || bindings.ClauseCount() != expectedTemplateCount {
		return nil, fmt.Errorf("control check program catalog does not match its reviewed binding catalog")
	}
	capabilities, err := providercapability.Load()
	if err != nil {
		return nil, fmt.Errorf("load shipped provider capabilities: %w", err)
	}
	capabilityByCollector := make(map[string]providercapability.Capability, len(capabilities))
	for _, capability := range capabilities {
		capabilityByCollector[capability.CollectorID] = capability
	}
	expectedRunnableCount := len(capabilities)
	if document.ControlCount != expectedControlCount || document.TemplateCount != expectedTemplateCount || len(document.Templates) != expectedTemplateCount ||
		document.PredicateDefinedCount != expectedTemplateCount || document.ImplementationMissingCount != 0 ||
		document.ProviderCapabilityMissingCount != expectedTemplateCount-expectedRunnableCount || document.EndToEndRunnableTemplateCount != expectedRunnableCount ||
		document.EndToEndRunnableControlCount != expectedRunnableCount || document.BlockedControlCount != expectedControlCount-expectedRunnableCount || document.ClassificationErrorCount != 0 {
		return nil, fmt.Errorf("control check program coverage or capability counts are invalid")
	}
	computedCatalogDigest, err := unsignedDigest(data, "catalog_sha256")
	if err != nil || computedCatalogDigest != document.CatalogSHA256 {
		return nil, fmt.Errorf("control check program catalog digest does not match its content")
	}

	templates := make([]Template, 0, len(document.Templates))
	byClause := make(map[string]int, len(document.Templates))
	seenIDs := map[string]bool{}
	seenControls := map[string]bool{}
	shapeCounts := map[string]int{}
	matchedCapabilities := map[string]bool{}
	previousControl := ""
	previousOrdinal := 0
	previousClause := ""
	for _, raw := range document.Templates {
		var candidate rawTemplate
		if err := decodeStrict(raw, &candidate); err != nil {
			return nil, fmt.Errorf("parse control check program template: %w", err)
		}
		binding, ok := bindings.Definition(candidate.ControlID)
		if !ok {
			return nil, fmt.Errorf("program template references unknown binding %q", candidate.ControlID)
		}
		template, err := validateTemplate(raw, candidate, document.ProgramSchemaSHA256, binding)
		if err != nil {
			return nil, err
		}
		capability, shipped := capabilityByCollector[template.CollectorContract.CollectorID]
		if template.EndToEndRunnable != shipped || shipped &&
			(capability.ControlID != template.ControlID || capability.ClauseOrdinal != template.ClauseOrdinal || capability.Authority != template.RequiredAuthority) {
			return nil, fmt.Errorf("program template provider claim does not match the shipped capability manifest")
		}
		if shipped {
			matchedCapabilities[capability.CollectorID] = true
		}
		if previousControl != "" && (template.ControlID < previousControl || (template.ControlID == previousControl && (template.ClauseOrdinal < previousOrdinal || (template.ClauseOrdinal == previousOrdinal && template.ClauseID <= previousClause)))) {
			return nil, fmt.Errorf("program templates are not in strict control/ordinal/clause order")
		}
		if (template.ControlID != previousControl && template.ClauseOrdinal != 1) || (template.ControlID == previousControl && template.ClauseOrdinal != previousOrdinal+1) {
			return nil, fmt.Errorf("program template ordinals are not contiguous for %s", template.ControlID)
		}
		key := lookupKey(template.ControlID, template.ClauseID)
		if _, duplicate := byClause[key]; duplicate || seenIDs[template.TemplateID] {
			return nil, fmt.Errorf("duplicate program template identity")
		}
		byClause[key] = len(templates)
		seenIDs[template.TemplateID] = true
		seenControls[template.ControlID] = true
		shapeCounts[template.PredicateShape]++
		templates = append(templates, template)
		previousControl, previousOrdinal, previousClause = template.ControlID, template.ClauseOrdinal, template.ClauseID
	}
	if len(seenControls) != expectedControlCount || len(matchedCapabilities) != len(capabilities) || !mapsEqualInt(shapeCounts, document.PredicateShapeCounts) {
		return nil, fmt.Errorf("program template control or predicate-shape coverage is incomplete")
	}
	for _, binding := range bindings.Definitions() {
		for _, clause := range binding.Clauses {
			if _, ok := byClause[lookupKey(binding.ControlID, clause.ClauseID)]; !ok {
				return nil, fmt.Errorf("binding clause has no program template")
			}
		}
	}
	return &Catalog{
		catalogSHA256: document.CatalogSHA256, programSchemaSHA256: document.ProgramSchemaSHA256,
		bindingCatalogSHA256: document.BindingCatalogSHA256, definitionSchemaSHA256: document.DefinitionSchemaSHA256,
		definitionCorpusSHA256: document.DefinitionCorpusSHA256, registryVersion: document.RegistryVersion,
		registrySHA256: document.RegistrySHA256, methodologySHA256: document.MethodologySHA256,
		classificationCorpusSHA256: document.ClassificationCorpusSHA256, templates: templates, byClause: byClause,
	}, nil
}

func validateTemplate(raw []byte, candidate rawTemplate, schemaDigest string, binding controlbinding.Definition) (Template, error) {
	if !validTemplateIdentity(candidate, schemaDigest) {
		return Template{}, fmt.Errorf("program template %s/%s has invalid identity data", candidate.ControlID, candidate.ClauseID)
	}
	authority, ok := validAuthorities[candidate.RequiredAuthority]
	if !ok || !validFamilies[candidate.CheckerFamily] {
		return Template{}, fmt.Errorf("program template %s/%s has unsupported family or authority", candidate.ControlID, candidate.ClauseID)
	}
	if binding.Revision != candidate.ControlRevision || binding.SemanticSHA256 != candidate.ControlSemanticSHA256 || candidate.ClauseOrdinal > len(binding.Clauses) {
		return Template{}, fmt.Errorf("program template has a stale control binding")
	}
	clause := binding.Clauses[candidate.ClauseOrdinal-1]
	if clause.Ordinal != candidate.ClauseOrdinal || clause.ClauseID != candidate.ClauseID || clause.Statement != candidate.ClauseStatement || clause.CheckerFamily != candidate.CheckerFamily || clause.EvidenceAuthority != candidate.RequiredAuthority || clause.ImplementationID != candidate.ImplementationID || clause.ImplementationContractSHA256 != candidate.ImplementationContractSHA256 {
		return Template{}, fmt.Errorf("program template does not exactly match its reviewed clause")
	}
	if candidate.ClauseStatementSHA256 != sha256Hex([]byte(candidate.ClauseStatement)) {
		return Template{}, fmt.Errorf("program template clause statement digest mismatch")
	}
	expectedID, _ := canonicalDigest(map[string]any{"control_id": candidate.ControlID, "control_revision": candidate.ControlRevision, "clause_id": candidate.ClauseID})
	if candidate.TemplateID != expectedID {
		return Template{}, fmt.Errorf("program template identity mismatch")
	}
	computed, err := unsignedDigest(raw, "template_sha256")
	if err != nil || computed != candidate.TemplateSHA256 {
		return Template{}, fmt.Errorf("program template digest mismatch")
	}
	facts, err := validateFactContracts(candidate.RawFactContracts, authority)
	if err != nil {
		return Template{}, err
	}
	parameters, err := validateParameterContracts(candidate.SealedParameterContracts)
	if err != nil {
		return Template{}, err
	}
	predicate, err := compilePredicate(candidate, authority, parameters)
	if err != nil {
		return Template{}, fmt.Errorf("program template predicate is invalid: %w", err)
	}
	if candidate.PredicateShape != string(predicate.Op) || !operationsMatch(candidate.RequiredRuntimeOps, predicate) || !referencesMatch(facts, parameters, predicate) {
		return Template{}, fmt.Errorf("program template predicate contracts are stale")
	}
	collector, err := validateCollector(candidate.CollectorContract)
	if err != nil {
		return Template{}, err
	}
	runtime, err := validateRuntime(candidate.RuntimeRequirements)
	if err != nil {
		return Template{}, err
	}
	registered := candidate.EndToEndRunnable != nil && *candidate.EndToEndRunnable
	if (collector.ProviderStatus == "registered") != registered || runtime.ProviderClaimed != registered {
		return Template{}, fmt.Errorf("program template provider declarations are inconsistent")
	}
	return Template{
		TemplateID: candidate.TemplateID, ProgramSchemaVersion: candidate.ProgramSchemaVersion, ProgramSchemaSHA256: candidate.ProgramSchemaSHA256,
		ControlID: candidate.ControlID, ControlRevision: candidate.ControlRevision, ControlSemanticSHA256: candidate.ControlSemanticSHA256,
		ClauseOrdinal: candidate.ClauseOrdinal, ClauseID: candidate.ClauseID, ClauseStatement: candidate.ClauseStatement, ClauseStatementSHA256: candidate.ClauseStatementSHA256,
		CheckerFamily: candidate.CheckerFamily, RequiredAuthority: authority, ImplementationID: candidate.ImplementationID, ImplementationContractSHA256: candidate.ImplementationContractSHA256,
		PredicateShape: candidate.PredicateShape, ReviewReason: candidate.ReviewReason, CounterexampleAnalysis: candidate.CounterexampleAnalysis,
		EndToEndRunnable: registered,
		RawFactContracts: facts, SealedParameterContracts: parameters, Predicate: predicate, RequiredRuntimeOps: append([]controlprogram.Operation(nil), candidate.RequiredRuntimeOps...),
		CollectorContract: collector, RuntimeRequirements: runtime, TemplateSHA256: candidate.TemplateSHA256,
	}, nil
}

func validTemplateIdentity(value rawTemplate, schemaDigest string) bool {
	return isDigest(value.TemplateID) && value.ProgramSchemaVersion == controlprogram.ProgramSchemaVersion && value.ProgramSchemaSHA256 == schemaDigest &&
		controlIDPattern.MatchString(value.ControlID) && value.ControlRevision > 0 && isDigest(value.ControlSemanticSHA256) &&
		value.ClauseOrdinal > 0 && value.ClauseOrdinal <= maximumClauseCount && isDigest(value.ClauseID) && strings.TrimSpace(value.ClauseStatement) != "" &&
		utf8.RuneCountInString(value.ClauseStatement) <= maximumStatementRunes && isDigest(value.ClauseStatementSHA256) &&
		implementationPattern.MatchString(value.ImplementationID) && isDigest(value.ImplementationContractSHA256) && isDigest(value.TemplateSHA256) &&
		value.PredicateDefined != nil && *value.PredicateDefined && providerStatusValid(value) &&
		strings.TrimSpace(value.ReviewReason) != "" && strings.TrimSpace(value.CounterexampleAnalysis) != ""
}

func providerStatusValid(value rawTemplate) bool {
	if value.EndToEndRunnable == nil {
		return false
	}
	if *value.EndToEndRunnable {
		return value.ReviewStatus == "predicate_defined_provider_registered" && value.ProviderCapabilityStatus == "registered"
	}
	return value.ReviewStatus == "predicate_defined_provider_unregistered" && value.ProviderCapabilityStatus == "unregistered"
}

func validateFactContracts(raw []rawFactContract, authority controlprogram.Authority) ([]FactContract, error) {
	if len(raw) == 0 || len(raw) > controlprogram.MaxFacts {
		return nil, fmt.Errorf("raw fact contracts are empty or oversized")
	}
	result := make([]FactContract, len(raw))
	seen := map[string]bool{}
	for index, value := range raw {
		leaf := strings.ToLower(value.FactKey[strings.LastIndex(value.FactKey, ".")+1:])
		mapped, ok := validAuthorities[value.Authority]
		genericText := strings.ToLower(value.RawValueSemantics + " " + value.SourceRequirement)
		if !factKeyPattern.MatchString(value.FactKey) || seen[value.FactKey] || !validFactTypes[value.FactType] || mapped != authority || !ok ||
			value.CompleteRequired == nil || !*value.CompleteRequired || opaqueVerdictNames[leaf] || forbiddenGenericKeys[leaf] || strings.Contains(genericText, "schema binding for every promise") ||
			strings.TrimSpace(value.RawValueSemantics) == "" || strings.TrimSpace(value.SourceRequirement) == "" {
			return nil, fmt.Errorf("invalid raw fact contract %q", value.FactKey)
		}
		seen[value.FactKey] = true
		result[index] = FactContract{value.FactKey, value.FactType, mapped, value.RawValueSemantics, value.SourceRequirement, true}
	}
	return result, nil
}

func validateParameterContracts(raw []rawParameterContract) ([]ParameterContract, error) {
	if len(raw) > controlprogram.MaxFacts {
		return nil, fmt.Errorf("sealed parameter contracts are oversized")
	}
	result := make([]ParameterContract, len(raw))
	seen := map[string]bool{}
	for index, value := range raw {
		leaf := strings.ToLower(value.ParameterKey[strings.LastIndex(value.ParameterKey, ".")+1:])
		origin := ParameterOrigin(value.ValueOrigin)
		if !factKeyPattern.MatchString(value.ParameterKey) || seen[value.ParameterKey] || forbiddenGenericKeys[leaf] || !validFactTypes[value.ParameterType] ||
			(origin != ParameterOriginScannerInventory && origin != ParameterOriginPolicy && origin != ParameterOriginContext) ||
			strings.Contains(strings.ToLower(value.SourceRequirement), "schema binding for every promise") || strings.TrimSpace(value.SourceRequirement) == "" {
			return nil, fmt.Errorf("invalid sealed parameter contract %q", value.ParameterKey)
		}
		seen[value.ParameterKey] = true
		result[index] = ParameterContract{value.ParameterKey, value.ParameterType, origin, value.SourceRequirement}
	}
	return result, nil
}

func validateCollector(raw rawCollectorContract) (CollectorContract, error) {
	if !collectorPattern.MatchString(raw.CollectorID) || len(raw.RequiredSources) == 0 || len(raw.RequiredSources) > 20 ||
		(raw.ProviderStatus != "registered" && raw.ProviderStatus != "unregistered") ||
		strings.TrimSpace(raw.InventoryContract) == "" || strings.TrimSpace(raw.NormalizationContract) == "" || strings.TrimSpace(raw.CompletenessContract) == "" || strings.TrimSpace(raw.FreshnessContract) == "" {
		return CollectorContract{}, fmt.Errorf("invalid collector contract")
	}
	seen := map[string]bool{}
	for _, source := range raw.RequiredSources {
		if strings.TrimSpace(source) == "" || seen[source] {
			return CollectorContract{}, fmt.Errorf("invalid collector source contract")
		}
		seen[source] = true
	}
	return CollectorContract{raw.CollectorID, append([]string(nil), raw.RequiredSources...), raw.InventoryContract, raw.NormalizationContract, raw.CompletenessContract, raw.FreshnessContract, raw.ProviderStatus}, nil
}

func validateRuntime(raw rawRuntimeRequirements) (RuntimeRequirements, error) {
	if raw.SubjectID != "inject_at_runtime" || raw.Subjects != "inject_complete_bounded_inventory_at_runtime" || raw.InventorySHA256 != "inject_at_runtime" ||
		raw.MaximumEvidenceAgeSeconds != "inject_approved_freshness_at_runtime" || raw.AllowNotApplicable != "inject_reviewed_applicability_at_runtime" ||
		raw.ApplicabilityProofContractSHA256 != "inject_at_runtime" || raw.SealedParameters != "compile_from_declared_trusted_origin_before_requesting_evidence" ||
		raw.SealedParametersBoundBy != "program_sha256" || raw.EvidenceProviderMaySupplyParameters == nil || *raw.EvidenceProviderMaySupplyParameters ||
		raw.ProviderRegistration != "required_before_evidence" || raw.ProviderClaimed == nil ||
		((*raw.ProviderClaimed && raw.DomainEvidenceCollector != "shipped_and_registered") ||
			(!*raw.ProviderClaimed && raw.DomainEvidenceCollector != "not_shipped_or_registered")) || raw.MissingCapabilityResult != "blocked" {
		return RuntimeRequirements{}, fmt.Errorf("runtime requirements are not fail closed")
	}
	return RuntimeRequirements{raw.SubjectID, raw.Subjects, raw.InventorySHA256, raw.MaximumEvidenceAgeSeconds, raw.AllowNotApplicable, raw.ApplicabilityProofContractSHA256, raw.SealedParameters, raw.SealedParametersBoundBy, false, raw.ProviderRegistration, raw.DomainEvidenceCollector, raw.MissingCapabilityResult, *raw.ProviderClaimed}, nil
}

func compilePredicate(candidate rawTemplate, authority controlprogram.Authority, contracts []ParameterContract) (controlprogram.Expression, error) {
	var predicate controlprogram.Expression
	if err := decodeStrict(candidate.Predicate, &predicate); err != nil {
		return controlprogram.Expression{}, err
	}
	parameters := map[string]controlprogram.Parameter{}
	for _, contract := range contracts {
		value, err := placeholderParameter(contract.ParameterType)
		if err != nil {
			return controlprogram.Expression{}, err
		}
		parameters[contract.ParameterKey] = value
	}
	program := controlprogram.Program{SchemaVersion: controlprogram.ProgramSchemaVersion, ControlID: candidate.ControlID, ControlRevision: candidate.ControlRevision,
		ControlSemanticSHA256: candidate.ControlSemanticSHA256, ClauseID: candidate.ClauseID, ClauseSHA256: candidate.ClauseStatementSHA256,
		ImplementationContractSHA256: candidate.ImplementationContractSHA256, SubjectID: "placeholder", Subjects: []string{"placeholder"}, InventorySHA256: strings.Repeat("0", 64),
		RequiredAuthority: authority, ApplicabilityProofContractSHA256: strings.Repeat("0", 64), MaximumEvidenceAgeSeconds: 1, Parameters: parameters, Predicate: predicate}
	if err := controlprogram.ValidateProgram(program); err != nil {
		return controlprogram.Expression{}, err
	}
	return cloneExpression(predicate), nil
}

func placeholderParameter(kind controlprogram.FactType) (controlprogram.Parameter, error) {
	text := "value"
	truth := false
	timestamp := "2026-01-01T00:00:00Z"
	digest := strings.Repeat("0", 64)
	switch kind {
	case controlprogram.FactIdentity, controlprogram.FactSchema, controlprogram.FactState:
		return controlprogram.Parameter{Type: kind, String: &text}, nil
	case controlprogram.FactDigest:
		return controlprogram.Parameter{Type: kind, String: &digest}, nil
	case controlprogram.FactString:
		return controlprogram.Parameter{Type: kind, String: &text}, nil
	case controlprogram.FactBoolean:
		return controlprogram.Parameter{Type: kind, Boolean: &truth}, nil
	case controlprogram.FactNumber:
		return controlprogram.Parameter{Type: kind, Number: json.Number("0")}, nil
	case controlprogram.FactTime:
		return controlprogram.Parameter{Type: kind, Timestamp: &timestamp}, nil
	case controlprogram.FactStringSet:
		return controlprogram.Parameter{Type: kind, Strings: []string{"value"}}, nil
	case controlprogram.FactIdentityMap, controlprogram.FactSchemaMap, controlprogram.FactStateMap, controlprogram.FactStringMap:
		return controlprogram.Parameter{Type: kind, Values: map[string]string{"key": "value"}}, nil
	case controlprogram.FactDigestMap:
		return controlprogram.Parameter{Type: kind, Values: map[string]string{"key": digest}}, nil
	case controlprogram.FactBooleanMap:
		return controlprogram.Parameter{Type: kind, Booleans: map[string]bool{"key": false}}, nil
	case controlprogram.FactNumberMap:
		return controlprogram.Parameter{Type: kind, Numbers: map[string]json.Number{"key": "0"}}, nil
	case controlprogram.FactTimeMap:
		return controlprogram.Parameter{Type: kind, Timestamps: map[string]string{"key": timestamp}}, nil
	case controlprogram.FactDirectedGraph:
		return controlprogram.Parameter{Type: kind, Edges: []controlprogram.DirectedEdge{{From: "a", To: "b"}}}, nil
	default:
		return controlprogram.Parameter{}, fmt.Errorf("unsupported sealed parameter type %q", kind)
	}
}

func parametersByContractOrigin(contracts []ParameterContract, runtime RuntimeBinding) (map[string]controlprogram.Parameter, bool) {
	totalRuntime := len(runtime.ScannerInventoryParameters) + len(runtime.AuthenticatedPolicyParameters) + len(runtime.AuthenticatedContextParameters)
	if totalRuntime != len(contracts) {
		return nil, false
	}
	result := make(map[string]controlprogram.Parameter, len(contracts))
	for _, contract := range contracts {
		var source map[string]controlprogram.Parameter
		switch contract.ValueOrigin {
		case ParameterOriginScannerInventory:
			source = runtime.ScannerInventoryParameters
		case ParameterOriginPolicy:
			source = runtime.AuthenticatedPolicyParameters
		case ParameterOriginContext:
			source = runtime.AuthenticatedContextParameters
		default:
			return nil, false
		}
		value, ok := source[contract.ParameterKey]
		if !ok || value.Type != contract.ParameterType {
			return nil, false
		}
		if _, duplicate := result[contract.ParameterKey]; duplicate {
			return nil, false
		}
		result[contract.ParameterKey] = value
	}
	return cloneParameters(result), true
}

func operationsMatch(required []controlprogram.Operation, predicate controlprogram.Expression) bool {
	actual := map[controlprogram.Operation]bool{}
	collectOperations(predicate, actual)
	if len(required) != len(actual) {
		return false
	}
	for index, operation := range required {
		if index > 0 && required[index-1] >= operation || !actual[operation] {
			return false
		}
	}
	return true
}

func collectOperations(expression controlprogram.Expression, result map[controlprogram.Operation]bool) {
	result[expression.Op] = true
	if expression.Arg != nil {
		collectOperations(*expression.Arg, result)
	}
	for _, child := range expression.Args {
		collectOperations(child, result)
	}
}

func referencesMatch(facts []FactContract, parameters []ParameterContract, predicate controlprogram.Expression) bool {
	wantFacts := map[string]bool{}
	for _, value := range facts {
		wantFacts[value.FactKey] = true
	}
	wantParameters := map[string]bool{}
	for _, value := range parameters {
		wantParameters[value.ParameterKey] = true
	}
	gotFacts, gotParameters := map[string]bool{}, map[string]bool{}
	collectReferences(predicate, gotFacts, gotParameters)
	return mapsEqualBool(wantFacts, gotFacts) && mapsEqualBool(wantParameters, gotParameters)
}

func collectReferences(expression controlprogram.Expression, facts, parameters map[string]bool) {
	if expression.Fact != "" {
		facts[expression.Fact] = true
	}
	if expression.OtherFact != "" {
		facts[expression.OtherFact] = true
	}
	if expression.ThirdFact != "" {
		facts[expression.ThirdFact] = true
	}
	if expression.Parameter != "" {
		parameters[expression.Parameter] = true
	}
	if expression.Arg != nil {
		collectReferences(*expression.Arg, facts, parameters)
	}
	for _, child := range expression.Args {
		collectReferences(child, facts, parameters)
	}
}

func cloneTemplate(value Template) Template {
	value.Predicate = cloneExpression(value.Predicate)
	value.RequiredRuntimeOps = append([]controlprogram.Operation(nil), value.RequiredRuntimeOps...)
	value.RawFactContracts = append([]FactContract(nil), value.RawFactContracts...)
	value.SealedParameterContracts = append([]ParameterContract(nil), value.SealedParameterContracts...)
	value.CollectorContract.RequiredSources = append([]string(nil), value.CollectorContract.RequiredSources...)
	return value
}
func cloneExpression(value controlprogram.Expression) controlprogram.Expression {
	value.Args = append([]controlprogram.Expression(nil), value.Args...)
	for index := range value.Args {
		value.Args[index] = cloneExpression(value.Args[index])
	}
	if value.Arg != nil {
		child := cloneExpression(*value.Arg)
		value.Arg = &child
	}
	value.Strings = append([]string(nil), value.Strings...)
	return value
}
func cloneParameters(values map[string]controlprogram.Parameter) map[string]controlprogram.Parameter {
	result := make(map[string]controlprogram.Parameter, len(values))
	for key, value := range values {
		value.Strings = append([]string(nil), value.Strings...)
		value.Values = cloneStringMap(value.Values)
		value.Booleans = cloneBoolMap(value.Booleans)
		value.Numbers = cloneNumberMap(value.Numbers)
		value.Timestamps = cloneStringMap(value.Timestamps)
		value.Edges = append([]controlprogram.DirectedEdge(nil), value.Edges...)
		result[key] = value
	}
	return result
}
func cloneStringMap[V ~string](values map[string]V) map[string]V {
	if values == nil {
		return nil
	}
	result := make(map[string]V, len(values))
	for key, value := range values {
		result[key] = value
	}
	return result
}
func cloneBoolMap(values map[string]bool) map[string]bool {
	if values == nil {
		return nil
	}
	result := make(map[string]bool, len(values))
	for key, value := range values {
		result[key] = value
	}
	return result
}
func cloneNumberMap(values map[string]json.Number) map[string]json.Number {
	if values == nil {
		return nil
	}
	result := make(map[string]json.Number, len(values))
	for key, value := range values {
		result[key] = value
	}
	return result
}

func validDocumentDigests(document rawDocument) bool {
	return isDigest(document.ProgramSchemaSHA256) && isDigest(document.BindingCatalogSHA256) && isDigest(document.DefinitionSchemaSHA256) && isDigest(document.DefinitionCorpusSHA256) && isDigest(document.RegistrySHA256) && isDigest(document.MethodologySHA256) && isDigest(document.ClassificationCorpusSHA256) && isDigest(document.CatalogSHA256) && versionPattern.MatchString(document.RegistryVersion)
}
func mapsEqualBool(left, right map[string]bool) bool {
	if len(left) != len(right) {
		return false
	}
	for key, value := range left {
		if right[key] != value {
			return false
		}
	}
	return true
}
func mapsEqualInt(left, right map[string]int) bool {
	if len(left) != len(right) {
		return false
	}
	for key, value := range left {
		if right[key] != value {
			return false
		}
	}
	return true
}
func lookupKey(controlID, clauseID string) string { return controlID + "\x00" + clauseID }
func isDigest(value string) bool                  { return digestPattern.MatchString(value) }
func sha256Hex(data []byte) string                { value := sha256.Sum256(data); return hex.EncodeToString(value[:]) }

func decodeStrict(data []byte, destination any) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	decoder.UseNumber()
	if err := decoder.Decode(destination); err != nil {
		return err
	}
	var trailing any
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		return fmt.Errorf("document contains trailing JSON")
	}
	return nil
}
func unsignedDigest(data []byte, digestKey string) (string, error) {
	var object map[string]any
	if err := decodeStrict(data, &object); err != nil {
		return "", err
	}
	if _, ok := object[digestKey]; !ok {
		return "", fmt.Errorf("missing %s", digestKey)
	}
	delete(object, digestKey)
	return canonicalDigest(object)
}
func canonicalDigest(value any) (string, error) {
	var output bytes.Buffer
	encoder := json.NewEncoder(&output)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(value); err != nil {
		return "", err
	}
	digest := sha256.Sum256(bytes.TrimSuffix(output.Bytes(), []byte{'\n'}))
	return hex.EncodeToString(digest[:]), nil
}

func rejectDuplicateKeys(data []byte) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	if err := walkJSON(decoder, 1); err != nil {
		return err
	}
	if _, err := decoder.Token(); !errors.Is(err, io.EOF) {
		return fmt.Errorf("document contains trailing JSON")
	}
	return nil
}
func walkJSON(decoder *json.Decoder, depth int) error {
	if depth > maximumJSONDepth {
		return fmt.Errorf("document exceeds JSON nesting depth")
	}
	token, err := decoder.Token()
	if err != nil {
		return err
	}
	delimiter, composite := token.(json.Delim)
	if !composite {
		return nil
	}
	switch delimiter {
	case '{':
		seen := map[string]bool{}
		for decoder.More() {
			keyToken, err := decoder.Token()
			if err != nil {
				return err
			}
			key, ok := keyToken.(string)
			if !ok {
				return fmt.Errorf("JSON key is not a string")
			}
			if seen[key] {
				return fmt.Errorf("duplicate JSON object key %q", key)
			}
			seen[key] = true
			if err := walkJSON(decoder, depth+1); err != nil {
				return err
			}
		}
		end, err := decoder.Token()
		if err != nil || end != json.Delim('}') {
			return fmt.Errorf("invalid JSON object terminator")
		}
	case '[':
		for decoder.More() {
			if err := walkJSON(decoder, depth+1); err != nil {
				return err
			}
		}
		end, err := decoder.Token()
		if err != nil || end != json.Delim(']') {
			return fmt.Errorf("invalid JSON array terminator")
		}
	default:
		return fmt.Errorf("unexpected JSON delimiter")
	}
	return nil
}

func readCatalogDocument(path string) ([]byte, error) {
	plain, plainErr := pathExists(path)
	compressedPath := path + ".gz"
	compressed, compressedErr := pathExists(compressedPath)
	if plainErr != nil || compressedErr != nil {
		return nil, fmt.Errorf("inspect control check programs: %v %v", plainErr, compressedErr)
	}
	if plain && compressed {
		return nil, fmt.Errorf("control check programs have ambiguous plain and compressed files")
	}
	if !plain && !compressed {
		return nil, fmt.Errorf("%w: %s or %s", os.ErrNotExist, path, compressedPath)
	}
	if plain {
		return readRegularBounded(path, maximumPlainBytes, "control check programs")
	}
	encoded, err := readRegularBounded(compressedPath, maximumCompressedBytes, "compressed control check programs")
	if err != nil {
		return nil, err
	}
	reader, err := gzip.NewReader(bytes.NewReader(encoded))
	if err != nil {
		return nil, fmt.Errorf("open compressed control check programs: %w", err)
	}
	reader.Multistream(false)
	data, readErr := io.ReadAll(io.LimitReader(reader, maximumExpandedBytes+1))
	closeErr := reader.Close()
	if readErr != nil || closeErr != nil {
		return nil, fmt.Errorf("read compressed control check programs: %v %v", readErr, closeErr)
	}
	if int64(len(data)) > maximumExpandedBytes {
		return nil, fmt.Errorf("expanded control check programs exceed the %d-byte limit", maximumExpandedBytes)
	}
	compressedReader := bytes.NewReader(encoded)
	check, _ := gzip.NewReader(compressedReader)
	check.Multistream(false)
	_, _ = io.Copy(io.Discard, check)
	_ = check.Close()
	if compressedReader.Len() != 0 {
		return nil, fmt.Errorf("compressed control check programs contain trailing gzip data")
	}
	return data, nil
}
func pathExists(path string) (bool, error) {
	_, err := os.Lstat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
func readRegularBounded(path string, limit int64, label string) ([]byte, error) {
	before, err := os.Lstat(path)
	if err != nil || !before.Mode().IsRegular() || before.Mode()&os.ModeSymlink != 0 || before.Size() > limit {
		return nil, fmt.Errorf("%s %s must be a regular file no larger than %d bytes", label, path, limit)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	after, err := os.Lstat(path)
	if err != nil || !os.SameFile(before, after) || int64(len(data)) != before.Size() {
		return nil, fmt.Errorf("%s changed while being read", label)
	}
	return data, nil
}

var _ = sort.Strings
