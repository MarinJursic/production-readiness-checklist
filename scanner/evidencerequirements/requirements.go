// Package evidencerequirements exposes the complete, authenticated contract a
// data-only evidence producer must satisfy. It never collects facts and never
// turns provider claims into scanner verdicts.
package evidencerequirements

import (
	"fmt"
	"regexp"
	"sort"

	"github.com/MarinJursic/production-readiness-checklist/scanner/controlprogram"
	"github.com/MarinJursic/production-readiness-checklist/scanner/controlprogramcatalog"
	"github.com/MarinJursic/production-readiness-checklist/scanner/evidencebundle"
	"github.com/MarinJursic/production-readiness-checklist/scanner/providercapability"
)

const SchemaVersion = "prc.evidence-requirements/v0.1"

var controlIDPattern = regexp.MustCompile(`^(?:PRC-[0-9]{2}-[0-9]{3}|USEQ-[A-F0-9]{8})$`)

// Filter narrows the exported producer contract without changing its meaning.
// Empty authority and control ID select all values. CollectorStatus is all,
// built_in, or missing.
type Filter struct {
	Authority       controlprogram.Authority
	ControlID       string
	CollectorStatus string
}

type Selection struct {
	Authority       string `json:"authority"`
	ControlID       string `json:"control_id"`
	CollectorStatus string `json:"collector_status"`
}

type AuthoritySummary struct {
	Name                   string `json:"name"`
	SelectedClauseCount    int    `json:"selected_clause_count"`
	BuiltInCollectorCount  int    `json:"built_in_collector_count"`
	MissingCollectorCount  int    `json:"missing_collector_count"`
	SignedImportRouteCount int    `json:"signed_import_route_count"`
}

type FactRequirement struct {
	ID                string                  `json:"id"`
	Type              controlprogram.FactType `json:"type"`
	RawValueSemantics string                  `json:"raw_value_semantics"`
	SourceRequirement string                  `json:"source_requirement"`
	CompleteRequired  bool                    `json:"complete_required"`
}

type ParameterRequirement struct {
	ID                string                                `json:"id"`
	Type              controlprogram.FactType               `json:"type"`
	Origin            controlprogramcatalog.ParameterOrigin `json:"origin"`
	SourceRequirement string                                `json:"source_requirement"`
}

type Requirement struct {
	TemplateID                       string                   `json:"template_id"`
	ControlID                        string                   `json:"control_id"`
	ControlRevision                  int                      `json:"control_revision"`
	ClauseID                         string                   `json:"clause_id"`
	ClauseOrdinal                    int                      `json:"clause_ordinal"`
	ClauseStatement                  string                   `json:"clause_statement"`
	CheckerFamily                    string                   `json:"checker_family"`
	Authority                        controlprogram.Authority `json:"authority"`
	ImplementationID                 string                   `json:"implementation_id"`
	ImplementationContractSHA256     string                   `json:"implementation_contract_sha256"`
	CollectorID                      string                   `json:"collector_id"`
	CollectorStatus                  string                   `json:"collector_status"`
	SignedImportSupported            bool                     `json:"signed_import_supported"`
	RequiredSources                  []string                 `json:"required_sources"`
	InventoryContract                string                   `json:"inventory_contract"`
	NormalizationContract            string                   `json:"normalization_contract"`
	CompletenessContract             string                   `json:"completeness_contract"`
	FreshnessContract                string                   `json:"freshness_contract"`
	Facts                            []FactRequirement        `json:"facts"`
	Parameters                       []ParameterRequirement   `json:"parameters"`
	RequiresInventoryInput           bool                     `json:"requires_inventory_input"`
	RequiresAuthenticatedPolicyInput bool                     `json:"requires_authenticated_policy_input"`
	RequiresAuthenticatedContext     bool                     `json:"requires_authenticated_context_input"`
	MissingEvidenceResult            string                   `json:"missing_evidence_result"`
}

type Report struct {
	SchemaVersion          string             `json:"schema_version"`
	CatalogSHA256          string             `json:"catalog_sha256"`
	RegistrySHA256         string             `json:"registry_sha256"`
	BindingCatalogSHA256   string             `json:"binding_catalog_sha256"`
	ExactControlCount      int                `json:"exact_control_count"`
	ExactClauseCount       int                `json:"exact_clause_count"`
	SelectedControlCount   int                `json:"selected_control_count"`
	SelectedClauseCount    int                `json:"selected_clause_count"`
	BuiltInCollectorCount  int                `json:"built_in_collector_count"`
	MissingCollectorCount  int                `json:"missing_collector_count"`
	SignedImportRouteCount int                `json:"signed_import_route_count"`
	Selection              Selection          `json:"selection"`
	Authorities            []AuthoritySummary `json:"authorities"`
	Requirements           []Requirement      `json:"requirements"`
}

// Build authenticates the exact-program catalog and shipped collector
// manifest before exporting any requirement. The report is stable: it contains
// no clock time and follows catalog order.
func Build(root string, filter Filter) (Report, error) {
	if err := validateFilter(filter); err != nil {
		return Report{}, err
	}
	catalog, err := controlprogramcatalog.Load(root)
	if err != nil {
		return Report{}, err
	}
	capabilities, err := providercapability.Load()
	if err != nil {
		return Report{}, err
	}
	builtIn := make(map[string]providercapability.Capability, len(capabilities))
	for _, capability := range capabilities {
		builtIn[capability.CollectorID] = capability
	}

	report := Report{
		SchemaVersion: SchemaVersion, CatalogSHA256: catalog.Digest(), RegistrySHA256: catalog.RegistrySHA256(),
		BindingCatalogSHA256: catalog.BindingCatalogSHA256(), ExactControlCount: catalog.ControlCount(),
		ExactClauseCount: catalog.TemplateCount(),
		Selection:        Selection{Authority: string(filter.Authority), ControlID: filter.ControlID, CollectorStatus: filter.CollectorStatus},
		Authorities:      []AuthoritySummary{}, Requirements: []Requirement{},
	}
	authorities := map[string]*AuthoritySummary{}
	controls := map[string]bool{}
	for _, template := range catalog.Templates() {
		status := "missing"
		if _, ok := builtIn[template.CollectorContract.CollectorID]; ok {
			status = "built_in"
		}
		if !selected(template, status, filter) {
			continue
		}
		requirement := requirementFor(template, status)
		report.Requirements = append(report.Requirements, requirement)
		controls[template.ControlID] = true
		summary := authorities[string(template.RequiredAuthority)]
		if summary == nil {
			summary = &AuthoritySummary{Name: string(template.RequiredAuthority)}
			authorities[summary.Name] = summary
		}
		summary.SelectedClauseCount++
		if status == "built_in" {
			report.BuiltInCollectorCount++
			summary.BuiltInCollectorCount++
		} else {
			report.MissingCollectorCount++
			summary.MissingCollectorCount++
		}
		if requirement.SignedImportSupported {
			report.SignedImportRouteCount++
			summary.SignedImportRouteCount++
		}
	}
	for _, summary := range authorities {
		report.Authorities = append(report.Authorities, *summary)
	}
	sort.Slice(report.Authorities, func(left, right int) bool {
		return report.Authorities[left].Name < report.Authorities[right].Name
	})
	report.SelectedControlCount = len(controls)
	report.SelectedClauseCount = len(report.Requirements)
	if report.SelectedClauseCount == 0 {
		return Report{}, fmt.Errorf("evidence requirement filters selected no exact clauses")
	}
	return report, nil
}

func validateFilter(filter Filter) error {
	if filter.Authority != "" && !evidencebundle.SupportsAuthority(filter.Authority) {
		return fmt.Errorf("unsupported evidence authority %q", filter.Authority)
	}
	if filter.ControlID != "" && !controlIDPattern.MatchString(filter.ControlID) {
		return fmt.Errorf("invalid control ID %q", filter.ControlID)
	}
	if filter.CollectorStatus != "all" && filter.CollectorStatus != "built_in" && filter.CollectorStatus != "missing" {
		return fmt.Errorf("collector status must be all, built_in, or missing")
	}
	return nil
}

func selected(template controlprogramcatalog.Template, status string, filter Filter) bool {
	return (filter.Authority == "" || template.RequiredAuthority == filter.Authority) &&
		(filter.ControlID == "" || template.ControlID == filter.ControlID) &&
		(filter.CollectorStatus == "all" || filter.CollectorStatus == status)
}

func requirementFor(template controlprogramcatalog.Template, status string) Requirement {
	requirement := Requirement{
		TemplateID: template.TemplateID, ControlID: template.ControlID, ControlRevision: template.ControlRevision,
		ClauseID: template.ClauseID, ClauseOrdinal: template.ClauseOrdinal, ClauseStatement: template.ClauseStatement,
		CheckerFamily: template.CheckerFamily, Authority: template.RequiredAuthority,
		ImplementationID: template.ImplementationID, ImplementationContractSHA256: template.ImplementationContractSHA256,
		CollectorID: template.CollectorContract.CollectorID, CollectorStatus: status,
		SignedImportSupported: evidencebundle.SupportsAuthority(template.RequiredAuthority),
		RequiredSources:       append([]string(nil), template.CollectorContract.RequiredSources...),
		InventoryContract:     template.CollectorContract.InventoryContract,
		NormalizationContract: template.CollectorContract.NormalizationContract,
		CompletenessContract:  template.CollectorContract.CompletenessContract,
		FreshnessContract:     template.CollectorContract.FreshnessContract,
		Facts:                 []FactRequirement{}, Parameters: []ParameterRequirement{}, MissingEvidenceResult: "blocked",
	}
	for _, fact := range template.RawFactContracts {
		requirement.Facts = append(requirement.Facts, FactRequirement{
			ID: fact.FactKey, Type: fact.FactType, RawValueSemantics: fact.RawValueSemantics,
			SourceRequirement: fact.SourceRequirement, CompleteRequired: fact.CompleteRequired,
		})
	}
	for _, parameter := range template.SealedParameterContracts {
		requirement.Parameters = append(requirement.Parameters, ParameterRequirement{
			ID: parameter.ParameterKey, Type: parameter.ParameterType, Origin: parameter.ValueOrigin,
			SourceRequirement: parameter.SourceRequirement,
		})
		switch parameter.ValueOrigin {
		case controlprogramcatalog.ParameterOriginScannerInventory:
			requirement.RequiresInventoryInput = true
		case controlprogramcatalog.ParameterOriginPolicy:
			requirement.RequiresAuthenticatedPolicyInput = true
		case controlprogramcatalog.ParameterOriginContext:
			requirement.RequiresAuthenticatedContext = true
		}
	}
	return requirement
}
