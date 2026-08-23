package finding

import (
	"strings"
	"testing"

	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
)

func validInput() Input {
	return Input{
		AssertionID: "PRC-A-CORE-001", ControlIDs: []string{"USEQ-BBBBBBBB", "USEQ-AAAAAAAA"},
		Title: "Repository has a README", Summary: "README.md is missing.",
		Severity: "high", Gate: "required", RemediationClass: "R2",
		Subject: model.FindingSubject{
			Kind: "project", ID: "example", InventoryDigest: strings.Repeat("a", 64),
		},
		Locations:   []model.FindingLocation{{Path: "z.go", Line: 2}, {Path: "a.go"}, {Path: "a.go"}},
		EvidenceIDs: []string{strings.Repeat("c", 64), strings.Repeat("b", 64)},
	}
}

func TestFindingIdentityIsStableAndContentAddressed(t *testing.T) {
	first, err := New(validInput())
	if err != nil {
		t.Fatal(err)
	}
	input := validInput()
	input.ControlIDs[0], input.ControlIDs[1] = input.ControlIDs[1], input.ControlIDs[0]
	input.EvidenceIDs[0], input.EvidenceIDs[1] = input.EvidenceIDs[1], input.EvidenceIDs[0]
	second, err := New(input)
	if err != nil {
		t.Fatal(err)
	}
	if first.ID != second.ID || first.Fingerprint != second.Fingerprint {
		t.Fatalf("canonical finding identity changed: %+v != %+v", first, second)
	}
	if err := Validate(first); err != nil {
		t.Fatal(err)
	}
	changed := first
	changed.Summary = "Different failure."
	if changed.Fingerprint != first.Fingerprint {
		t.Fatal("test mutated fingerprint unexpectedly")
	}
	if err := Validate(changed); err == nil || !strings.Contains(err.Error(), "ID") {
		t.Fatalf("content mutation was not detected: %v", err)
	}
}

func TestFindingFingerprintSurvivesEvidenceAndSummaryChanges(t *testing.T) {
	first, err := New(validInput())
	if err != nil {
		t.Fatal(err)
	}
	input := validInput()
	input.Summary = "A more detailed explanation."
	input.EvidenceIDs = []string{strings.Repeat("d", 64)}
	second, err := New(input)
	if err != nil {
		t.Fatal(err)
	}
	if first.Fingerprint != second.Fingerprint || first.ID == second.ID {
		t.Fatalf("fingerprint/content identity semantics are wrong: %+v %+v", first, second)
	}
}
