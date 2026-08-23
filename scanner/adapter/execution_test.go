package adapter

import (
	"strings"
	"testing"
	"time"
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
	if execution.SchemaVersion != "prc.adapter-execution/v0.1" || len(execution.ExecutionID) != 64 {
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
