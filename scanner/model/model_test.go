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
		SchemaVersion: RunSchema,
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
