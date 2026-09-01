package automaticcoverage_test

import (
	"path/filepath"
	"testing"

	"github.com/MarinJursic/production-readiness-checklist/scanner/automaticcoverage"
)

func TestReportKeepsCollectionSeparateFromImportCompatibility(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	reportValue, err := automaticcoverage.Build(root)
	if err != nil {
		t.Fatal(err)
	}
	if reportValue.SchemaVersion != automaticcoverage.SchemaVersion || reportValue.ControlCount != 10_042 ||
		reportValue.ReviewedRoutingControlCount != 10_042 || reportValue.DeterministicControlCount != 686 ||
		reportValue.NondeterministicControlCount != 9_356 || reportValue.ExactClauseCount != 765 ||
		reportValue.AdvisoryAIReviewControlCount != 9_356 ||
		reportValue.ExactPredicateClauseCount != 765 || reportValue.BuiltInCollectorClauseCount != 3 ||
		reportValue.SignedImportSupportedClauseCount != 765 || len(reportValue.Authorities) != 6 {
		t.Fatalf("automatic coverage = %+v", reportValue)
	}
	total, builtIn, signed := 0, 0, 0
	previous := ""
	for _, authority := range reportValue.Authorities {
		if previous != "" && previous >= authority.Name {
			t.Fatalf("authorities are not ordered: %+v", reportValue.Authorities)
		}
		previous = authority.Name
		total += authority.ExactClauseCount
		builtIn += authority.BuiltInCollectorClauseCount
		signed += authority.SignedImportSupportedClauseCount
	}
	if total != reportValue.ExactClauseCount || builtIn != reportValue.BuiltInCollectorClauseCount ||
		signed != reportValue.SignedImportSupportedClauseCount {
		t.Fatalf("authority totals do not match report: %+v", reportValue)
	}
}
