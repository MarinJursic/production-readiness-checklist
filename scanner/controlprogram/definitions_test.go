package controlprogram

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"
)

const expectedReviewedDefinitionCount = 765

type definitionEnvelope struct {
	SchemaVersion              string              `json:"schema_version"`
	Scope                      string              `json:"scope"`
	SourceBindingCatalogSHA256 string              `json:"source_binding_catalog_sha256"`
	DefinitionCount            int                 `json:"definition_count"`
	Definitions                []programDefinition `json:"definitions"`
}

type rawFactContract struct {
	FactKey          string    `json:"fact_id"`
	FactType         FactType  `json:"fact_type"`
	Authority        Authority `json:"authority"`
	RawSemantics     string    `json:"raw_value_semantics"`
	Source           string    `json:"source_requirement"`
	CompleteRequired bool      `json:"complete_required"`
}

type sealedParameterContract struct {
	ParameterKey  string   `json:"parameter_id"`
	ParameterType FactType `json:"parameter_type"`
	ValueOrigin   string   `json:"value_origin"`
	Source        string   `json:"source_requirement"`
}

type programFixture struct {
	Description     string               `json:"description"`
	Parameters      map[string]Parameter `json:"parameters"`
	Facts           map[string]Fact      `json:"facts"`
	ExpectedOutcome Outcome              `json:"expected_outcome"`
}

type programFixtures struct {
	Pass           programFixture `json:"pass"`
	Fail           programFixture `json:"fail"`
	Blocked        programFixture `json:"blocked"`
	Counterexample programFixture `json:"counterexample"`
}

type programDefinition struct {
	ControlID                  string                    `json:"control_id"`
	ControlRevision            int                       `json:"control_revision"`
	ControlSemanticSHA256      string                    `json:"control_semantic_sha256"`
	ClauseOrdinal              int                       `json:"clause_ordinal"`
	ClauseID                   string                    `json:"clause_id"`
	ClauseStatement            string                    `json:"clause_statement"`
	CheckerFamily              string                    `json:"checker_family"`
	RequiredAuthority          Authority                 `json:"required_authority"`
	CorrectedCheckerFamily     string                    `json:"corrected_checker_family"`
	CorrectedRequiredAuthority Authority                 `json:"corrected_required_authority"`
	ClassificationStatus       string                    `json:"classification_status"`
	ClassificationErrorReason  *string                   `json:"classification_error_reason"`
	RawFactContracts           []rawFactContract         `json:"raw_fact_contracts"`
	SealedParameterContracts   []sealedParameterContract `json:"sealed_parameter_contracts"`
	Predicate                  Expression                `json:"predicate"`
	RequiredRuntimeOps         []Operation               `json:"required_runtime_ops"`
	CollectorContract          json.RawMessage           `json:"collector_contract"`
	Fixtures                   programFixtures           `json:"fixtures"`
	ReviewReason               string                    `json:"review_reason"`
	CounterexampleAnalysis     string                    `json:"counterexample_analysis"`
}

func TestEveryReviewedDeterministicDefinitionExecutesAdversarialFixtures(t *testing.T) {
	root := filepath.Clean(filepath.Join("..", ".."))
	bindingBytes, err := os.ReadFile(filepath.Join(root, "catalog", "control-check-bindings.json"))
	if err != nil {
		t.Fatal(err)
	}
	bindingDigest := sha256.Sum256(bindingBytes)
	expectedBindingDigest := hex.EncodeToString(bindingDigest[:])
	paths, err := filepath.Glob(filepath.Join(root, "research", "control-classification", "program-specs", "*.json"))
	if err != nil || len(paths) == 0 {
		t.Fatalf("find deterministic program definitions: paths=%d err=%v", len(paths), err)
	}
	sort.Strings(paths)
	seen := make(map[string]struct{}, expectedReviewedDefinitionCount)
	definitionCount := 0
	for _, path := range paths {
		info, err := os.Stat(path)
		if err != nil {
			t.Fatal(err)
		}
		if info.Size() > 3_000_000 {
			t.Fatalf("definition part %s exceeds the 3 MB review-file limit", path)
		}
		envelope := decodeDefinitionEnvelope(t, path)
		if envelope.SchemaVersion != "prc.control-check-program-definitions/v0.1" || envelope.SourceBindingCatalogSHA256 != expectedBindingDigest {
			t.Fatalf("definition part %s has stale schema or binding digest", path)
		}
		if envelope.Scope != "structured_record" && envelope.Scope != "non_structured" {
			t.Fatalf("definition part %s has invalid scope %q", path, envelope.Scope)
		}
		if envelope.DefinitionCount != len(envelope.Definitions) {
			t.Fatalf("definition part %s has stale definition_count", path)
		}
		for _, definition := range envelope.Definitions {
			definitionCount++
			identity := fmt.Sprintf("%s/%d/%s", definition.ControlID, definition.ClauseOrdinal, definition.ClauseID)
			if _, duplicate := seen[identity]; duplicate {
				t.Fatalf("duplicate deterministic definition %s", identity)
			}
			seen[identity] = struct{}{}
			if definition.ClassificationStatus != "exact_predicate" || definition.ClassificationErrorReason != nil {
				t.Fatalf("definition %s is not an exact predicate", identity)
			}
			encoded, err := json.Marshal(definition)
			if err != nil {
				t.Fatal(err)
			}
			lower := strings.ToLower(string(encoded))
			for _, forbidden := range []string{
				`"/semantic/`, "schema binding for every promise", "tokenized statement words",
				"record-paths", "required-paths", "record_paths", "required_paths",
				"clause_satisfied", "provider_verdict", "is bounded by direct raw facts",
				"closed operations", "jointly decide every named part of", "can look correct while",
				"is missing or mismatched for the exact promise",
				"the program reads", "the passing relation is specific", "false-pass scenario for",
				"evidence that may still appear healthy", "decisive broken observation",
				"is bounded to the sealed", "its pinned schema exposes these concrete members rather than statement tokens",
				"fixture rejects a complete-looking record",
			} {
				if strings.Contains(lower, forbidden) {
					t.Fatalf("definition %s contains forbidden generic delegation %q", identity, forbidden)
				}
			}
			if definition.CorrectedRequiredAuthority == "" || len(definition.RawFactContracts) == 0 || len(definition.CollectorContract) == 0 {
				t.Fatalf("definition %s is missing its exact evidence contract", identity)
			}
			runDefinitionFixtures(t, definition)
		}
	}
	if definitionCount != expectedReviewedDefinitionCount {
		t.Fatalf("executed %d definitions; want %d", definitionCount, expectedReviewedDefinitionCount)
	}
}

func decodeDefinitionEnvelope(t *testing.T, path string) definitionEnvelope {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	decoder.UseNumber()
	var envelope definitionEnvelope
	if err := decoder.Decode(&envelope); err != nil {
		t.Fatalf("decode %s: %v", path, err)
	}
	var trailing any
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		t.Fatalf("definition part %s has trailing JSON", path)
	}
	return envelope
}

func runDefinitionFixtures(t *testing.T, definition programDefinition) {
	t.Helper()
	fixtures := []struct {
		name    string
		fixture programFixture
	}{
		{name: "pass", fixture: definition.Fixtures.Pass},
		{name: "fail", fixture: definition.Fixtures.Fail},
		{name: "blocked", fixture: definition.Fixtures.Blocked},
		{name: "counterexample", fixture: definition.Fixtures.Counterexample},
	}
	for _, item := range fixtures {
		item := item
		t.Run(definition.ControlID+"_"+fmt.Sprint(definition.ClauseOrdinal)+"_"+item.name, func(t *testing.T) {
			now := time.Date(2026, 8, 29, 12, 0, 0, 0, time.UTC)
			program := Program{
				SchemaVersion: ProgramSchemaVersion, ControlID: definition.ControlID,
				ControlRevision: definition.ControlRevision, ControlSemanticSHA256: definition.ControlSemanticSHA256,
				ClauseID: definition.ClauseID, ClauseSHA256: digestString(definition.ClauseStatement),
				ImplementationContractSHA256: digestString("fixture implementation contract"),
				SubjectID:                    "fixture-subject", Subjects: []string{"fixture-subject"}, InventorySHA256: digestString("fixture inventory"),
				RequiredAuthority: definition.CorrectedRequiredAuthority, AllowNotApplicable: false,
				ApplicabilityProofContractSHA256: digestString("fixture applicability contract"),
				MaximumEvidenceAgeSeconds:        3600, Parameters: item.fixture.Parameters, Predicate: definition.Predicate,
			}
			if err := ValidateProgram(program); err != nil {
				t.Fatalf("fixture program is invalid: %v", err)
			}
			evidence := Evidence{
				SchemaVersion: EvidenceSchemaVersion, EvidenceID: "fixture-evidence", ProgramSHA256: ProgramSHA256(program),
				ControlID: program.ControlID, ControlRevision: program.ControlRevision,
				ControlSemanticSHA256: program.ControlSemanticSHA256, ClauseID: program.ClauseID,
				ClauseSHA256: program.ClauseSHA256, ImplementationContractSHA256: program.ImplementationContractSHA256,
				SubjectID: program.SubjectID, ObservedSubjects: append([]string(nil), program.Subjects...),
				InventorySHA256: program.InventorySHA256, Authority: program.RequiredAuthority,
				ObservedAt: now.Add(-time.Minute), Complete: true, Applicability: ApplicabilityApplicable,
				ApplicabilityProofContractSHA256: program.ApplicabilityProofContractSHA256, Facts: item.fixture.Facts,
			}
			result := Evaluate(program, evidence, now)
			if result.Outcome != item.fixture.ExpectedOutcome {
				t.Fatalf("got %s (%s), want %s: %s", result.Outcome, result.ReasonCode, item.fixture.ExpectedOutcome, item.fixture.Description)
			}
		})
	}
}
