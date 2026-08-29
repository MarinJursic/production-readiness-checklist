// Package controlbinding loads the reviewed deterministic-control bindings.
//
// Loading is deliberately separate from executing a check. This package only
// authenticates the catalog shape and its binding to controlcheck's closed
// implementation registry. Evidence providers and verifier registrations are
// runtime capabilities; their absence must remain Blocked at evaluation time.
package controlbinding

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
	"regexp"
	"slices"
	"strings"
	"unicode/utf8"

	"github.com/MarinJursic/production-readiness-checklist/scanner/controlcheck"
)

const (
	catalogRelativePath = "catalog/control-check-bindings.json"
	schemaVersion       = "prc.control-check-bindings/v0.1"
	generatorID         = "prc.build-control-check-bindings@0.1"

	expectedImplementationCount = 11
	expectedBindingCount        = 686
	expectedClauseCount         = 765
	maximumClauseCount          = 50
	maximumStatementCharacters  = 2000

	maximumPlainBytes      int64 = 4 * 1024 * 1024
	maximumCompressedBytes int64 = 2 * 1024 * 1024
	maximumExpandedBytes   int64 = maximumPlainBytes
)

var (
	digestPattern         = regexp.MustCompile(`^[0-9a-f]{64}$`)
	versionPattern        = regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+$`)
	controlIDPattern      = regexp.MustCompile(`^(?:PRC-[0-9]{2}-[0-9]{3}|USEQ-[A-F0-9]{8})$`)
	implementationPattern = regexp.MustCompile(`^prc\.check\.[a-z0-9-]+@[0-9]+\.[0-9]+$`)
)

var validRoutes = map[string]struct{}{
	"local_static": {}, "artifact_verification": {}, "bounded_execution": {},
	"external_readonly_query": {}, "structured_record_validation": {},
	"deterministic_composite": {},
}

var validAuthorities = map[string]struct{}{
	"repository": {}, "artifact": {}, "executed": {}, "environment": {},
	"external_registry": {}, "structured_record": {},
}

var passRequirements = []string{
	"applicability_resolved_applicable",
	"exact_implementation_registered",
	"provider_contract_registered",
	"required_evidence_authority_available",
	"complete_bounded_evidence",
	"full_clause_proven",
}

var failRequirements = []string{
	"applicability_resolved_applicable",
	"exact_implementation_registered",
	"provider_contract_registered",
	"required_evidence_authority_available",
	"complete_bounded_evidence",
	"bounded_violation_observed",
}

var blockedConditions = []string{
	"applicability_missing_or_ambiguous",
	"implementation_unregistered",
	"provider_contract_missing",
	"required_evidence_authority_unavailable",
	"unsupported_target",
	"complete_inventory_unavailable",
	"evidence_missing_stale_partial_or_conflicting",
	"rate_limited_or_ambiguous_external_result",
}

// Catalog is an immutable, validated snapshot. Every method that returns a
// definition returns a defensive deep copy.
type Catalog struct {
	digest                     string
	registryVersion            string
	registrySHA256             string
	methodologySHA256          string
	classificationCorpusSHA256 string
	implementations            []ImplementationDefinition
	definitions                []Definition
	byControl                  map[string]int
	clauseCount                int
}

// ImplementationDefinition is one implementation contract named by the
// catalog. Runtime code must still be registered separately.
type ImplementationDefinition struct {
	CheckerFamily                string
	ImplementationID             string
	SupportedEvidenceAuthorities []string
	CapabilityClass              string
	ExternalProviderRegistration string
	ImplementationContractSHA256 string
	RegistrationState            string
	OnUnregistered               string
	ExternalProviderClaimed      bool
}

// Definition is the reviewed, revision-bound definition of one deterministic
// control. It is not an executable controlcheck.Binding until a later caller
// supplies a bounded subject inventory and runtime capabilities.
type Definition struct {
	ControlID             string
	Revision              int
	SemanticSHA256        string
	FinalRowSHA256        string
	Route                 string
	Aggregation           string
	ApplicabilityContract ApplicabilityContract
	Clauses               []ClauseDefinition
}

// ApplicabilityContract describes the fail-closed applicability result.
type ApplicabilityContract struct {
	ApplicableResult      string
	NotApplicableRequires string
	UnresolvedResult      string
}

// ClauseDefinition is one atomic, implementation-bound deterministic clause.
type ClauseDefinition struct {
	ClauseID                     string
	Ordinal                      int
	Statement                    string
	CheckerFamily                string
	EvidenceAuthority            string
	ImplementationID             string
	ImplementationContractSHA256 string
	ImplementationRegistration   string
	ProviderContract             string
	ExternalProviderClaimed      bool
	ResultContract               ResultContract
}

// ResultContract records the exact Pass, Fail, and Blocked preconditions.
type ResultContract struct {
	Pass    ResultBranch
	Fail    ResultBranch
	Blocked BlockedBranch
}

// ResultBranch is an all-required Pass or Fail branch.
type ResultBranch struct {
	Result      string
	RequiresAll []string
}

// BlockedBranch is an any-condition Blocked branch.
type BlockedBranch struct {
	Result  string
	WhenAny []string
}

type rawDocument struct {
	SchemaVersion              string              `json:"schema_version"`
	GeneratorID                string              `json:"generator_id"`
	RegistryVersion            string              `json:"registry_version"`
	RegistrySHA256             string              `json:"registry_sha256"`
	MethodologySHA256          string              `json:"methodology_sha256"`
	ClassificationCorpusSHA256 string              `json:"classification_corpus_sha256"`
	ImplementationRegistry     []rawImplementation `json:"implementation_registry"`
	BindingCount               int                 `json:"binding_count"`
	Bindings                   []rawBinding        `json:"bindings"`
}

type rawImplementation struct {
	CheckerFamily                string   `json:"checker_family"`
	ImplementationID             string   `json:"implementation_id"`
	SupportedEvidenceAuthorities []string `json:"supported_evidence_authorities"`
	CapabilityClass              string   `json:"capability_class"`
	ExternalProviderRegistration string   `json:"external_provider_registration"`
	ImplementationContractSHA256 string   `json:"implementation_contract_sha256"`
	RegistrationState            string   `json:"registration_state"`
	OnUnregistered               string   `json:"on_unregistered"`
	ExternalProviderClaimed      *bool    `json:"external_provider_claimed"`
}

type rawBinding struct {
	ControlID             string                   `json:"control_id"`
	Revision              int                      `json:"revision"`
	SemanticSHA256        string                   `json:"semantic_sha256"`
	FinalRowSHA256        string                   `json:"final_row_sha256"`
	Route                 string                   `json:"route"`
	Aggregation           string                   `json:"aggregation"`
	ApplicabilityContract rawApplicabilityContract `json:"applicability_contract"`
	Clauses               []rawClause              `json:"clauses"`
}

type rawApplicabilityContract struct {
	ApplicableResult      string `json:"applicable_result"`
	NotApplicableRequires string `json:"not_applicable_requires"`
	UnresolvedResult      string `json:"unresolved_result"`
}

type rawClause struct {
	ClauseID                     string            `json:"clause_id"`
	Ordinal                      int               `json:"ordinal"`
	Statement                    string            `json:"statement"`
	CheckerFamily                string            `json:"checker_family"`
	EvidenceAuthority            string            `json:"evidence_authority"`
	ImplementationID             string            `json:"implementation_id"`
	ImplementationContractSHA256 string            `json:"implementation_contract_sha256"`
	ImplementationRegistration   string            `json:"implementation_registration"`
	ProviderContract             string            `json:"provider_contract"`
	ExternalProviderClaimed      *bool             `json:"external_provider_claimed"`
	ResultContract               rawResultContract `json:"result_contract"`
}

type rawResultContract struct {
	Pass    rawResultBranch  `json:"pass"`
	Fail    rawResultBranch  `json:"fail"`
	Blocked rawBlockedBranch `json:"blocked"`
}

type rawResultBranch struct {
	Result      string   `json:"result"`
	RequiresAll []string `json:"requires_all"`
}

type rawBlockedBranch struct {
	Result  string   `json:"result"`
	WhenAny []string `json:"when_any"`
}

// Load reads catalog/control-check-bindings.json or its compact .json.gz form
// below root. Exactly one form must exist.
func Load(root string) (*Catalog, error) {
	absolute, err := filepath.Abs(root)
	if err != nil {
		return nil, fmt.Errorf("resolve binding catalog root: %w", err)
	}
	data, err := readCatalogDocument(filepath.Join(absolute, filepath.FromSlash(catalogRelativePath)))
	if err != nil {
		return nil, err
	}
	if !utf8.Valid(data) {
		return nil, fmt.Errorf("control check bindings are not valid UTF-8")
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	var document rawDocument
	if err := decoder.Decode(&document); err != nil {
		return nil, fmt.Errorf("parse control check bindings: %w", err)
	}
	var trailing any
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		return nil, fmt.Errorf("control check bindings contain trailing JSON")
	}
	catalog, err := validate(document)
	if err != nil {
		return nil, err
	}
	digest := sha256.Sum256(data)
	catalog.digest = hex.EncodeToString(digest[:])
	return catalog, nil
}

// Digest returns the SHA-256 of the expanded JSON bytes. Plain and compressed
// copies of identical catalog content therefore have the same identity.
func (catalog *Catalog) Digest() string { return catalog.digest }

// RegistryVersion returns the catalog's bound control-registry version.
func (catalog *Catalog) RegistryVersion() string { return catalog.registryVersion }

// RegistrySHA256 returns the catalog's bound control-registry digest.
func (catalog *Catalog) RegistrySHA256() string { return catalog.registrySHA256 }

// MethodologySHA256 returns the reviewed classification-methodology digest.
func (catalog *Catalog) MethodologySHA256() string { return catalog.methodologySHA256 }

// ClassificationCorpusSHA256 returns the final-review corpus digest.
func (catalog *Catalog) ClassificationCorpusSHA256() string {
	return catalog.classificationCorpusSHA256
}

// BindingCount returns the exact number of reviewed deterministic controls.
func (catalog *Catalog) BindingCount() int { return len(catalog.definitions) }

// ClauseCount returns the exact number of atomic deterministic clauses.
func (catalog *Catalog) ClauseCount() int { return catalog.clauseCount }

// Implementations returns defensive copies of the catalog contracts.
func (catalog *Catalog) Implementations() []ImplementationDefinition {
	result := make([]ImplementationDefinition, len(catalog.implementations))
	for index, implementation := range catalog.implementations {
		result[index] = cloneImplementation(implementation)
	}
	return result
}

// Definitions returns every binding in stable control-ID order.
func (catalog *Catalog) Definitions() []Definition {
	result := make([]Definition, len(catalog.definitions))
	for index, definition := range catalog.definitions {
		result[index] = cloneDefinition(definition)
	}
	return result
}

// Definition returns one binding as a defensive copy.
func (catalog *Catalog) Definition(controlID string) (Definition, bool) {
	index, ok := catalog.byControl[controlID]
	if !ok {
		return Definition{}, false
	}
	return cloneDefinition(catalog.definitions[index]), true
}

func validate(document rawDocument) (*Catalog, error) {
	if document.SchemaVersion != schemaVersion {
		return nil, fmt.Errorf("unsupported control check binding schema %q", document.SchemaVersion)
	}
	if document.GeneratorID != generatorID {
		return nil, fmt.Errorf("unsupported control check binding generator %q", document.GeneratorID)
	}
	if !versionPattern.MatchString(document.RegistryVersion) ||
		!isDigest(document.RegistrySHA256) || !isDigest(document.MethodologySHA256) ||
		!isDigest(document.ClassificationCorpusSHA256) {
		return nil, fmt.Errorf("control check binding envelope has invalid version or digest data")
	}
	if document.BindingCount != expectedBindingCount || len(document.Bindings) != expectedBindingCount {
		return nil, fmt.Errorf("control check binding count must be exactly %d", expectedBindingCount)
	}
	implementations, implementationByFamily, err := validateImplementations(document.ImplementationRegistry)
	if err != nil {
		return nil, err
	}
	definitions := make([]Definition, 0, len(document.Bindings))
	byControl := make(map[string]int, len(document.Bindings))
	clauseCount := 0
	previousControlID := ""
	for _, binding := range document.Bindings {
		if _, duplicate := byControl[binding.ControlID]; duplicate {
			return nil, fmt.Errorf("duplicate control binding %q", binding.ControlID)
		}
		if previousControlID != "" && binding.ControlID <= previousControlID {
			return nil, fmt.Errorf("control bindings are not in strict control-ID order at %q", binding.ControlID)
		}
		definition, err := validateBinding(binding, implementationByFamily)
		if err != nil {
			return nil, err
		}
		byControl[definition.ControlID] = len(definitions)
		definitions = append(definitions, definition)
		clauseCount += len(definition.Clauses)
		previousControlID = definition.ControlID
	}
	if clauseCount != expectedClauseCount {
		return nil, fmt.Errorf("control check clause count must be exactly %d", expectedClauseCount)
	}
	return &Catalog{
		registryVersion:            document.RegistryVersion,
		registrySHA256:             document.RegistrySHA256,
		methodologySHA256:          document.MethodologySHA256,
		classificationCorpusSHA256: document.ClassificationCorpusSHA256,
		implementations:            implementations,
		definitions:                definitions,
		byControl:                  byControl,
		clauseCount:                clauseCount,
	}, nil
}

func validateImplementations(raw []rawImplementation) ([]ImplementationDefinition, map[string]ImplementationDefinition, error) {
	runtimeDescriptors := controlcheck.Implementations()
	if len(raw) != expectedImplementationCount || len(runtimeDescriptors) != expectedImplementationCount {
		return nil, nil, fmt.Errorf("implementation registry must contain exactly %d reviewed contracts", expectedImplementationCount)
	}
	byRuntimeFamily := make(map[string]controlcheck.Descriptor, len(runtimeDescriptors))
	for _, descriptor := range runtimeDescriptors {
		family := string(descriptor.Family)
		if _, duplicate := byRuntimeFamily[family]; duplicate {
			return nil, nil, fmt.Errorf("runtime implementation registry has duplicate family %q", family)
		}
		byRuntimeFamily[family] = descriptor
	}
	result := make([]ImplementationDefinition, 0, len(raw))
	byFamily := make(map[string]ImplementationDefinition, len(raw))
	previousFamily := ""
	for _, candidate := range raw {
		if _, duplicate := byFamily[candidate.CheckerFamily]; duplicate {
			return nil, nil, fmt.Errorf("duplicate catalog implementation family %q", candidate.CheckerFamily)
		}
		if previousFamily != "" && candidate.CheckerFamily <= previousFamily {
			return nil, nil, fmt.Errorf("catalog implementations are not in strict family order at %q", candidate.CheckerFamily)
		}
		if !implementationPattern.MatchString(candidate.ImplementationID) ||
			!isDigest(candidate.ImplementationContractSHA256) ||
			candidate.CapabilityClass == "" ||
			candidate.ExternalProviderRegistration != "required_when_authority_is_environment_or_external_registry" ||
			candidate.RegistrationState != "runtime_registration_required" ||
			candidate.OnUnregistered != "blocked" ||
			candidate.ExternalProviderClaimed == nil || *candidate.ExternalProviderClaimed {
			return nil, nil, fmt.Errorf("implementation %q has an invalid fail-closed contract", candidate.CheckerFamily)
		}
		if len(candidate.SupportedEvidenceAuthorities) == 0 || hasDuplicateStrings(candidate.SupportedEvidenceAuthorities) {
			return nil, nil, fmt.Errorf("implementation %q has invalid evidence authorities", candidate.CheckerFamily)
		}
		for _, authority := range candidate.SupportedEvidenceAuthorities {
			if _, ok := validAuthorities[authority]; !ok {
				return nil, nil, fmt.Errorf("implementation %q has unsupported evidence authority %q", candidate.CheckerFamily, authority)
			}
		}
		computed, err := implementationContractDigest(candidate)
		if err != nil || computed != candidate.ImplementationContractSHA256 {
			return nil, nil, fmt.Errorf("implementation %q contract digest does not match its content", candidate.CheckerFamily)
		}
		runtimeDescriptor, ok := byRuntimeFamily[candidate.CheckerFamily]
		if !ok || runtimeDescriptor.ImplementationID != candidate.ImplementationID ||
			runtimeDescriptor.ImplementationDigest != candidate.ImplementationContractSHA256 ||
			!sameDescriptorAuthorities(runtimeDescriptor, candidate.SupportedEvidenceAuthorities) {
			return nil, nil, fmt.Errorf("implementation %q does not match the closed runtime descriptor", candidate.CheckerFamily)
		}
		definition := ImplementationDefinition{
			CheckerFamily: candidate.CheckerFamily, ImplementationID: candidate.ImplementationID,
			SupportedEvidenceAuthorities: append([]string(nil), candidate.SupportedEvidenceAuthorities...),
			CapabilityClass:              candidate.CapabilityClass,
			ExternalProviderRegistration: candidate.ExternalProviderRegistration,
			ImplementationContractSHA256: candidate.ImplementationContractSHA256,
			RegistrationState:            candidate.RegistrationState, OnUnregistered: candidate.OnUnregistered,
			ExternalProviderClaimed: false,
		}
		result = append(result, definition)
		byFamily[definition.CheckerFamily] = definition
		previousFamily = candidate.CheckerFamily
	}
	if len(byFamily) != len(byRuntimeFamily) {
		return nil, nil, fmt.Errorf("catalog implementation registry does not exactly cover the runtime registry")
	}
	return result, byFamily, nil
}

func validateBinding(binding rawBinding, implementations map[string]ImplementationDefinition) (Definition, error) {
	if !controlIDPattern.MatchString(binding.ControlID) || binding.Revision < 1 ||
		!isDigest(binding.SemanticSHA256) || !isDigest(binding.FinalRowSHA256) {
		return Definition{}, fmt.Errorf("control binding %q has invalid identity data", binding.ControlID)
	}
	if _, ok := validRoutes[binding.Route]; !ok || binding.Aggregation != "all_clauses_pass" {
		return Definition{}, fmt.Errorf("control binding %s has invalid routing or aggregation", binding.ControlID)
	}
	applicability := binding.ApplicabilityContract
	if applicability.ApplicableResult != "applicable" ||
		applicability.NotApplicableRequires != "authoritative_bounded_absence_proof" ||
		applicability.UnresolvedResult != "blocked" {
		return Definition{}, fmt.Errorf("control binding %s has an invalid applicability contract", binding.ControlID)
	}
	if len(binding.Clauses) < 1 || len(binding.Clauses) > maximumClauseCount {
		return Definition{}, fmt.Errorf("control binding %s has an invalid clause count", binding.ControlID)
	}
	clauses := make([]ClauseDefinition, 0, len(binding.Clauses))
	seenClauses := make(map[string]struct{}, len(binding.Clauses))
	for index, clause := range binding.Clauses {
		if _, duplicate := seenClauses[clause.ClauseID]; duplicate {
			return Definition{}, fmt.Errorf("control binding %s has duplicate clause %s", binding.ControlID, clause.ClauseID)
		}
		definition, err := validateClause(binding.ControlID, index+1, clause, implementations)
		if err != nil {
			return Definition{}, err
		}
		seenClauses[definition.ClauseID] = struct{}{}
		clauses = append(clauses, definition)
	}
	return Definition{
		ControlID: binding.ControlID, Revision: binding.Revision,
		SemanticSHA256: binding.SemanticSHA256, FinalRowSHA256: binding.FinalRowSHA256,
		Route: binding.Route, Aggregation: binding.Aggregation,
		ApplicabilityContract: ApplicabilityContract(applicability),
		Clauses: clauses,
	}, nil
}

func validateClause(controlID string, expectedOrdinal int, clause rawClause,
	implementations map[string]ImplementationDefinition) (ClauseDefinition, error) {
	if clause.Ordinal != expectedOrdinal || !isDigest(clause.ClauseID) ||
		strings.TrimSpace(clause.Statement) == "" || utf8.RuneCountInString(clause.Statement) > maximumStatementCharacters {
		return ClauseDefinition{}, fmt.Errorf("control binding %s has an invalid clause at ordinal %d", controlID, expectedOrdinal)
	}
	implementation, ok := implementations[clause.CheckerFamily]
	if !ok || clause.ImplementationID != implementation.ImplementationID ||
		clause.ImplementationContractSHA256 != implementation.ImplementationContractSHA256 ||
		!slices.Contains(implementation.SupportedEvidenceAuthorities, clause.EvidenceAuthority) {
		return ClauseDefinition{}, fmt.Errorf("control binding %s clause %d has an invalid implementation binding", controlID, expectedOrdinal)
	}
	if clause.ImplementationRegistration != "required_before_execution" ||
		clause.ProviderContract != "required_before_execution" ||
		clause.ExternalProviderClaimed == nil || *clause.ExternalProviderClaimed {
		return ClauseDefinition{}, fmt.Errorf("control binding %s clause %d does not fail closed", controlID, expectedOrdinal)
	}
	computed, err := clauseDigest(clause)
	if err != nil || computed != clause.ClauseID {
		return ClauseDefinition{}, fmt.Errorf("control binding %s clause %d digest does not match its content", controlID, expectedOrdinal)
	}
	if !validResultContract(clause.ResultContract) {
		return ClauseDefinition{}, fmt.Errorf("control binding %s clause %d has an invalid result contract", controlID, expectedOrdinal)
	}
	return ClauseDefinition{
		ClauseID: clause.ClauseID, Ordinal: clause.Ordinal, Statement: clause.Statement,
		CheckerFamily: clause.CheckerFamily, EvidenceAuthority: clause.EvidenceAuthority,
		ImplementationID:             clause.ImplementationID,
		ImplementationContractSHA256: clause.ImplementationContractSHA256,
		ImplementationRegistration:   clause.ImplementationRegistration,
		ProviderContract:             clause.ProviderContract, ExternalProviderClaimed: false,
		ResultContract: ResultContract{
			Pass:    ResultBranch{Result: clause.ResultContract.Pass.Result, RequiresAll: append([]string(nil), clause.ResultContract.Pass.RequiresAll...)},
			Fail:    ResultBranch{Result: clause.ResultContract.Fail.Result, RequiresAll: append([]string(nil), clause.ResultContract.Fail.RequiresAll...)},
			Blocked: BlockedBranch{Result: clause.ResultContract.Blocked.Result, WhenAny: append([]string(nil), clause.ResultContract.Blocked.WhenAny...)},
		},
	}, nil
}

func validResultContract(contract rawResultContract) bool {
	return contract.Pass.Result == "pass" && slices.Equal(contract.Pass.RequiresAll, passRequirements) &&
		contract.Fail.Result == "fail" && slices.Equal(contract.Fail.RequiresAll, failRequirements) &&
		contract.Blocked.Result == "blocked" && slices.Equal(contract.Blocked.WhenAny, blockedConditions)
}

func implementationContractDigest(candidate rawImplementation) (string, error) {
	contract := struct {
		CapabilityClass              string   `json:"capability_class"`
		CheckerFamily                string   `json:"checker_family"`
		ExternalProviderRegistration string   `json:"external_provider_registration"`
		ImplementationID             string   `json:"implementation_id"`
		SupportedEvidenceAuthorities []string `json:"supported_evidence_authorities"`
	}{
		CapabilityClass: candidate.CapabilityClass, CheckerFamily: candidate.CheckerFamily,
		ExternalProviderRegistration: candidate.ExternalProviderRegistration,
		ImplementationID:             candidate.ImplementationID,
		SupportedEvidenceAuthorities: candidate.SupportedEvidenceAuthorities,
	}
	return digestCanonical(contract)
}

func clauseDigest(clause rawClause) (string, error) {
	contract := struct {
		CheckerFamily     string `json:"checker_family"`
		EvidenceAuthority string `json:"evidence_authority"`
		Statement         string `json:"statement"`
	}{
		CheckerFamily:     clause.CheckerFamily,
		EvidenceAuthority: clause.EvidenceAuthority,
		Statement:         clause.Statement,
	}
	return digestCanonical(contract)
}

func digestCanonical(value any) (string, error) {
	var output bytes.Buffer
	encoder := json.NewEncoder(&output)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(value); err != nil {
		return "", err
	}
	data := bytes.TrimSuffix(output.Bytes(), []byte{'\n'})
	digest := sha256.Sum256(data)
	return hex.EncodeToString(digest[:]), nil
}

func sameDescriptorAuthorities(descriptor controlcheck.Descriptor, authorities []string) bool {
	descriptorAuthorities := make([]string, 0, len(descriptor.Capabilities))
	seen := map[string]struct{}{}
	for _, capability := range descriptor.Capabilities {
		authority := string(capability.Authority)
		if _, duplicate := seen[authority]; duplicate {
			continue
		}
		seen[authority] = struct{}{}
		descriptorAuthorities = append(descriptorAuthorities, authority)
	}
	slices.Sort(descriptorAuthorities)
	candidate := append([]string(nil), authorities...)
	slices.Sort(candidate)
	return slices.Equal(descriptorAuthorities, candidate)
}

func hasDuplicateStrings(values []string) bool {
	seen := make(map[string]struct{}, len(values))
	for _, value := range values {
		if _, duplicate := seen[value]; duplicate {
			return true
		}
		seen[value] = struct{}{}
	}
	return false
}

func isDigest(value string) bool { return digestPattern.MatchString(value) }

func cloneImplementation(value ImplementationDefinition) ImplementationDefinition {
	value.SupportedEvidenceAuthorities = append([]string(nil), value.SupportedEvidenceAuthorities...)
	return value
}

func cloneDefinition(value Definition) Definition {
	clauses := value.Clauses
	value.Clauses = make([]ClauseDefinition, len(clauses))
	for index, clause := range clauses {
		clause.ResultContract.Pass.RequiresAll = append([]string(nil), clause.ResultContract.Pass.RequiresAll...)
		clause.ResultContract.Fail.RequiresAll = append([]string(nil), clause.ResultContract.Fail.RequiresAll...)
		clause.ResultContract.Blocked.WhenAny = append([]string(nil), clause.ResultContract.Blocked.WhenAny...)
		value.Clauses[index] = clause
	}
	return value
}

func readCatalogDocument(path string) ([]byte, error) {
	plainExists, plainErr := pathExists(path)
	compressedPath := path + ".gz"
	compressedExists, compressedErr := pathExists(compressedPath)
	if plainErr != nil {
		return nil, fmt.Errorf("inspect control check bindings %s: %w", path, plainErr)
	}
	if compressedErr != nil {
		return nil, fmt.Errorf("inspect compressed control check bindings %s: %w", compressedPath, compressedErr)
	}
	if plainExists && compressedExists {
		return nil, fmt.Errorf("control check bindings have ambiguous plain and compressed files")
	}
	if !plainExists && !compressedExists {
		return nil, fmt.Errorf("%w: %s or %s", os.ErrNotExist, path, compressedPath)
	}
	if plainExists {
		return readRegularBounded(path, maximumPlainBytes, "control check bindings")
	}
	encoded, err := readRegularBounded(compressedPath, maximumCompressedBytes, "compressed control check bindings")
	if err != nil {
		return nil, err
	}
	compressed := bytes.NewReader(encoded)
	reader, err := gzip.NewReader(compressed)
	if err != nil {
		return nil, fmt.Errorf("open compressed control check bindings %s: %w", compressedPath, err)
	}
	// The npm artifact is one deterministic gzip member. Disabling multistream
	// also makes appended members or bytes observable instead of silently
	// treating them as part of the trusted catalog representation.
	reader.Multistream(false)
	data, readErr := io.ReadAll(io.LimitReader(reader, maximumExpandedBytes+1))
	closeErr := reader.Close()
	if readErr != nil {
		return nil, fmt.Errorf("read compressed control check bindings %s: %w", compressedPath, readErr)
	}
	if closeErr != nil {
		return nil, fmt.Errorf("close compressed control check bindings %s: %w", compressedPath, closeErr)
	}
	if int64(len(data)) > maximumExpandedBytes {
		return nil, fmt.Errorf("expanded control check bindings exceed the %d-byte limit", maximumExpandedBytes)
	}
	if compressed.Len() != 0 {
		return nil, fmt.Errorf("compressed control check bindings contain a trailing gzip member or data")
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
