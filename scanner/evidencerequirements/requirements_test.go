package evidencerequirements

import (
	"path/filepath"
	"testing"

	"github.com/MarinJursic/production-readiness-checklist/scanner/controlprogram"
)

func TestBuildExportsEveryExactClauseAndHonestCollectorCoverage(t *testing.T) {
	report, err := Build(filepath.Join("..", ".."), Filter{CollectorStatus: "all"})
	if err != nil {
		t.Fatal(err)
	}
	if report.SchemaVersion != SchemaVersion || report.ExactControlCount != 686 || report.ExactClauseCount != 765 ||
		report.SelectedControlCount != 686 || report.SelectedClauseCount != 765 ||
		report.BuiltInCollectorCount != 1 || report.MissingCollectorCount != 764 ||
		report.SignedImportRouteCount != 765 || len(report.Requirements) != 765 || len(report.Authorities) != 6 {
		t.Fatalf("unexpected report summary: controls=%d clauses=%d selected_controls=%d selected_clauses=%d built_in=%d missing=%d imports=%d authorities=%d",
			report.ExactControlCount, report.ExactClauseCount, report.SelectedControlCount, report.SelectedClauseCount,
			report.BuiltInCollectorCount, report.MissingCollectorCount, report.SignedImportRouteCount, len(report.Authorities))
	}
	total := AuthoritySummary{}
	for _, authority := range report.Authorities {
		total.SelectedClauseCount += authority.SelectedClauseCount
		total.BuiltInCollectorCount += authority.BuiltInCollectorCount
		total.MissingCollectorCount += authority.MissingCollectorCount
		total.SignedImportRouteCount += authority.SignedImportRouteCount
	}
	if total.SelectedClauseCount != report.SelectedClauseCount ||
		total.BuiltInCollectorCount != report.BuiltInCollectorCount ||
		total.MissingCollectorCount != report.MissingCollectorCount ||
		total.SignedImportRouteCount != report.SignedImportRouteCount {
		t.Fatalf("authority totals do not match report: %+v", total)
	}
	for _, requirement := range report.Requirements {
		if requirement.MissingEvidenceResult != "blocked" || !requirement.SignedImportSupported ||
			len(requirement.Facts) == 0 || requirement.TemplateID == "" || requirement.ImplementationContractSHA256 == "" {
			t.Fatalf("incomplete requirement: %+v", requirement)
		}
	}
}

func TestBuildFiltersMissingRepositoryCollectors(t *testing.T) {
	report, err := Build(filepath.Join("..", ".."), Filter{
		Authority: controlprogram.AuthorityRepository, CollectorStatus: "missing",
	})
	if err != nil {
		t.Fatal(err)
	}
	if report.SelectedClauseCount != 25 || report.SelectedControlCount != 25 || report.BuiltInCollectorCount != 0 ||
		report.MissingCollectorCount != 25 || report.SignedImportRouteCount != 25 || len(report.Authorities) != 1 ||
		report.Authorities[0].Name != "repository" {
		t.Fatalf("unexpected repository filter report: %+v", report)
	}
	for _, requirement := range report.Requirements {
		if requirement.Authority != controlprogram.AuthorityRepository || requirement.CollectorStatus != "missing" {
			t.Fatalf("filter leaked requirement: %+v", requirement)
		}
	}
}

func TestBuildFindsTheBuiltInDocumentedCommandsCollector(t *testing.T) {
	report, err := Build(filepath.Join("..", ".."), Filter{ControlID: "PRC-36-004", CollectorStatus: "built_in"})
	if err != nil {
		t.Fatal(err)
	}
	if report.SelectedClauseCount != 1 || report.BuiltInCollectorCount != 1 || len(report.Requirements) != 1 {
		t.Fatalf("unexpected built-in report: %+v", report)
	}
	requirement := report.Requirements[0]
	if requirement.ControlID != "PRC-36-004" || requirement.CollectorStatus != "built_in" ||
		requirement.CollectorID == "" || len(requirement.Facts) == 0 || !requirement.RequiresInventoryInput {
		t.Fatalf("unexpected documented-commands requirement: %+v", requirement)
	}
	for _, parameter := range requirement.Parameters {
		if parameter.Origin != "scanner_inventory" {
			t.Fatalf("unexpected parameter origin: %+v", parameter)
		}
	}
}

func TestBuildRejectsInvalidOrEmptySelections(t *testing.T) {
	root := filepath.Join("..", "..")
	for name, filter := range map[string]Filter{
		"authority": {Authority: "website"},
		"control":   {ControlID: "PRC-NOT-A-CONTROL"},
		"status":    {CollectorStatus: "available"},
		"empty":     {ControlID: "PRC-36-004", CollectorStatus: "missing"},
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := Build(root, filter); err == nil {
				t.Fatalf("filter %+v was accepted", filter)
			}
		})
	}
}
