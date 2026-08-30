// Package automaticcoverage reports separate, non-inflated measures for the
// reviewed rule corpus, exact predicates, built-in evidence collection, and
// signed external evidence support.
package automaticcoverage

import (
	"fmt"
	"sort"

	"github.com/MarinJursic/production-readiness-checklist/scanner/controlprogramcatalog"
	"github.com/MarinJursic/production-readiness-checklist/scanner/evidencebundle"
	"github.com/MarinJursic/production-readiness-checklist/scanner/fullscan"
	"github.com/MarinJursic/production-readiness-checklist/scanner/providercapability"
)

const SchemaVersion = "prc.automatic-coverage/v0.1"

type Authority struct {
	Name                             string `json:"name"`
	ExactClauseCount                 int    `json:"exact_clause_count"`
	BuiltInCollectorClauseCount      int    `json:"built_in_collector_clause_count"`
	SignedImportSupportedClauseCount int    `json:"signed_import_supported_clause_count"`
}

// Report deliberately does not combine the measures into one percentage.
// Exact predicates, import compatibility, and collected observations answer
// different questions and must not be added together.
type Report struct {
	SchemaVersion                    string      `json:"schema_version"`
	ControlCount                     int         `json:"control_count"`
	ReviewedRoutingControlCount      int         `json:"reviewed_routing_control_count"`
	DeterministicControlCount        int         `json:"deterministic_control_count"`
	NondeterministicControlCount     int         `json:"nondeterministic_control_count"`
	AdvisoryAIReviewControlCount     int         `json:"advisory_ai_review_control_count"`
	ExactClauseCount                 int         `json:"exact_clause_count"`
	ExactPredicateClauseCount        int         `json:"exact_predicate_clause_count"`
	BuiltInCollectorClauseCount      int         `json:"built_in_collector_clause_count"`
	SignedImportSupportedClauseCount int         `json:"signed_import_supported_clause_count"`
	Authorities                      []Authority `json:"authorities"`
}

func Build(root string) (Report, error) {
	registryCoverage, err := fullscan.LoadRegistryCoverage(root)
	if err != nil {
		return Report{}, err
	}
	programs, err := controlprogramcatalog.Load(root)
	if err != nil {
		return Report{}, err
	}
	capabilities, err := providercapability.Load()
	if err != nil {
		return Report{}, err
	}
	templates := programs.Templates()
	byCollector := make(map[string]controlprogramcatalog.Template, len(templates))
	deterministicControls := map[string]bool{}
	authorities := map[string]*Authority{}
	for _, template := range templates {
		if _, exists := byCollector[template.CollectorContract.CollectorID]; exists {
			return Report{}, fmt.Errorf("exact program catalog repeats collector %s", template.CollectorContract.CollectorID)
		}
		byCollector[template.CollectorContract.CollectorID] = template
		deterministicControls[template.ControlID] = true
		name := string(template.RequiredAuthority)
		item := authorities[name]
		if item == nil {
			item = &Authority{Name: name}
			authorities[name] = item
		}
		item.ExactClauseCount++
		if evidencebundle.SupportsAuthority(template.RequiredAuthority) {
			item.SignedImportSupportedClauseCount++
		}
	}
	for _, capability := range capabilities {
		template, ok := byCollector[capability.CollectorID]
		if !ok || template.ControlID != capability.ControlID || template.ClauseOrdinal != capability.ClauseOrdinal ||
			template.RequiredAuthority != capability.Authority {
			return Report{}, fmt.Errorf("built-in collector capability %s does not match an exact clause", capability.CollectorID)
		}
		authorities[string(capability.Authority)].BuiltInCollectorClauseCount++
	}
	items := make([]Authority, 0, len(authorities))
	signedImportCount := 0
	for _, item := range authorities {
		items = append(items, *item)
		signedImportCount += item.SignedImportSupportedClauseCount
	}
	sort.Slice(items, func(left, right int) bool { return items[left].Name < items[right].Name })
	controlCount := registryCoverage.ControlCount
	if registryCoverage.DeterministicCount != len(deterministicControls) ||
		registryCoverage.NondeterministicCount != controlCount-len(deterministicControls) {
		return Report{}, fmt.Errorf("exact program controls do not match the reviewed classification split")
	}
	return Report{
		SchemaVersion: SchemaVersion,
		ControlCount:  controlCount, ReviewedRoutingControlCount: controlCount,
		DeterministicControlCount:    len(deterministicControls),
		NondeterministicControlCount: registryCoverage.NondeterministicCount,
		AdvisoryAIReviewControlCount: registryCoverage.NondeterministicCount,
		ExactClauseCount:             len(templates), ExactPredicateClauseCount: len(templates),
		BuiltInCollectorClauseCount:      len(capabilities),
		SignedImportSupportedClauseCount: signedImportCount,
		Authorities:                      items,
	}, nil
}
