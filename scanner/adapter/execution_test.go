package adapter

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
)

func validRunOutput(outcome string) RunOutput {
	started := time.Date(2026, 8, 23, 12, 0, 0, 0, time.UTC)
	return RunOutput{
		Transcript: Transcript{
			Logs: []Log{}, Artifacts: []Artifact{},
			Observations: []Observation{{
				ID: "OBS-1", Kind: "fixture-result", Outcome: outcome,
				Summary: "Fixture observation.", Locations: []Location{},
			}},
			Summary: Summary{Type: "summary", Status: "completed", Counts: map[string]int{"observations": 1}},
		},
		StartedAt: started, CompletedAt: started.Add(time.Second), DurationMS: 1000,
		DiagnosticsSHA256: strings.Repeat("0", 64), DiagnosticsBytes: 0,
	}
}

func TestBindExecutionProducesContentAddressedValidatedRecord(t *testing.T) {
	manifest := validManifest()
	execution, err := BindExecution(strings.Repeat("a", 64), Subject{
		TargetName: "example", TargetCommit: strings.Repeat("b", 40), InventoryDigest: strings.Repeat("c", 64),
	}, manifest, validRunOutput("not_found"))
	if err != nil {
		t.Fatal(err)
	}
	if execution.SchemaVersion != "prc.adapter-execution/v0.2" || len(execution.ExecutionID) != 64 ||
		execution.Resolution.Source != ResolutionSourceExplicitLocal ||
		execution.Resolution.PublisherID != manifest.Publisher.ID {
		t.Fatalf("unexpected execution identity: %+v", execution)
	}
	if err := ValidateExecution(execution); err != nil {
		t.Fatal(err)
	}
	execution.Transcript.Observations[0].Outcome = "found"
	if err := ValidateExecution(execution); err == nil {
		t.Fatal("mutated bound execution was accepted")
	}
}

func TestExecutionIdentityBindsRegistryResolutionAndPreservesLegacyV01(t *testing.T) {
	manifest := validManifest()
	resolution := model.AdapterResolution{
		Source: ResolutionSourceRegistry, PublisherID: manifest.Publisher.ID, Trust: "verified-community",
		RegistryID: "prc.adapter-registry.test@0.1", RegistryRevision: 7,
		RegistryDigest: strings.Repeat("d", 64),
	}
	execution, err := BindExecutionWithResolution(strings.Repeat("a", 64), Subject{
		TargetName: "example", InventoryDigest: strings.Repeat("c", 64),
	}, manifest, resolution, validRunOutput("not_found"))
	if err != nil {
		t.Fatal(err)
	}
	execution.Resolution.RegistryRevision++
	if err := ValidateExecution(execution); err == nil {
		t.Fatal("mutated registry provenance was accepted")
	}

	legacy, err := BindExecution(strings.Repeat("a", 64), Subject{
		TargetName: "example", InventoryDigest: strings.Repeat("c", 64),
	}, manifest, validRunOutput("not_found"))
	if err != nil {
		t.Fatal(err)
	}
	legacy.SchemaVersion = "prc.adapter-execution/v0.1"
	legacy.Resolution = model.AdapterResolution{}
	legacy.ExecutionID, err = executionID(legacy)
	if err != nil {
		t.Fatal(err)
	}
	if err := ValidateExecution(legacy); err != nil {
		t.Fatalf("legacy execution no longer validates: %v", err)
	}
	payload, err := json.Marshal(legacy)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(payload), "resolution") {
		t.Fatalf("legacy execution drifted: %s", payload)
	}
}

func TestBindExecutionRejectsInvalidTranscriptAndSubject(t *testing.T) {
	output := validRunOutput("not_found")
	output.Transcript.Observations[0].Locations = nil
	if _, err := BindExecution(strings.Repeat("a", 64), Subject{
		TargetName: "example", InventoryDigest: strings.Repeat("c", 64),
	}, validManifest(), output); err == nil {
		t.Fatal("invalid transcript was bound")
	}
	if _, err := BindExecution(strings.Repeat("a", 64), Subject{
		TargetName: "example", InventoryDigest: "wrong",
	}, validManifest(), validRunOutput("not_found")); err == nil {
		t.Fatal("invalid subject was bound")
	}
}
