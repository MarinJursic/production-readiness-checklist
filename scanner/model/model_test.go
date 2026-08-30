package model

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestLegacyControlJSONDoesNotGainReviewedClassificationFields(t *testing.T) {
	run := RunResult{
		SchemaVersion: "prc.run/v0.11",
		ControlCatalog: &ControlCatalogSummary{
			ContractGeneratorID:        "prc.control-contracts@0.2",
			ReviewedDeterministicCount: 1,
		},
		ControlResults: []ControlResult{{
			ControlID: "PRC-03-003", Classification: "deterministic",
			DeterministicBindingID: "PRC-03-003@1",
		}},
	}
	encoded, err := json.Marshal(run)
	if err != nil {
		t.Fatal(err)
	}
	for _, forbidden := range []string{"contract_generator_id", "reviewed_deterministic_count", "classification", "deterministic_binding_id"} {
		if strings.Contains(string(encoded), `"`+forbidden+`"`) {
			t.Fatalf("legacy v0.11 JSON gained %s: %s", forbidden, encoded)
		}
	}
}

func TestArchivedV012ZeroValueClassificationFieldsRemainOmitted(t *testing.T) {
	run := RunResult{
		SchemaVersion: "prc.run/v0.12",
		ControlCatalog: &ControlCatalogSummary{
			SchemaVersion:         "prc.control-catalog-summary/v0.1",
			ContractSchemaVersion: "prc.control-contracts/v0.1",
		},
		ControlResults: []ControlResult{{ControlID: "USEQ-AAAAAAAA"}},
	}
	encoded, err := json.Marshal(run)
	if err != nil {
		t.Fatal(err)
	}
	for _, forbidden := range []string{"contract_generator_id", "classification_methodology_sha256", "reviewed_deterministic_count", "classification_route", "classification_row_sha256", "deterministic_binding_sha256"} {
		if strings.Contains(string(encoded), `"`+forbidden+`"`) {
			t.Fatalf("archived v0.12 JSON gained %s: %s", forbidden, encoded)
		}
	}
}

func TestArchivedV012AIReviewDoesNotGainV013PlanningFields(t *testing.T) {
	run := RunResult{
		SchemaVersion: "prc.run/v0.12",
		ControlResults: []ControlResult{{
			ControlID: "PRC-01-001",
			AIReview: &AIControlReview{
				Provider:               "codex",
				ReviewDepth:            "standard",
				AssessmentCandidate:    "needs_review",
				ApplicabilityCandidate: "applicable",
				Confidence:             "medium",
				Priority:               "medium",
				Reason:                 "Legacy review",
				TaskID:                 strings.Repeat("a", 64),
			},
		}},
	}
	encoded, err := json.Marshal(run)
	if err != nil {
		t.Fatal(err)
	}
	for _, forbidden := range []string{
		"root_cause", "root_cause_key", "effort", "blast_radius",
		"ai_improvement_plan", "authoritative_evidence_bundles",
	} {
		if strings.Contains(string(encoded), `"`+forbidden+`"`) {
			t.Fatalf("archived v0.12 JSON gained %s: %s", forbidden, encoded)
		}
	}
}
