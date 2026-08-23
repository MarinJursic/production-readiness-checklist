package engine

import (
	"strings"
	"testing"
	"time"

	"github.com/MarinJursic/production-readiness-checklist/scanner/adapter"
	"github.com/MarinJursic/production-readiness-checklist/scanner/inventory"
	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
)

func analysisManifest() adapter.Manifest {
	limits := adapter.DefaultLimits()
	return adapter.Manifest{
		SchemaVersion: adapter.ManifestSchema, ID: "prc.adapter.analysis-fixture@0.1",
		Title: "Analysis evidence fixture", Description: "Produces static-analysis observations for engine tests.",
		Publisher: adapter.Publisher{ID: "prc-project", Name: "Production Readiness Checklist"},
		Owner:     "scanner-maintainers", Maintenance: "active",
		Protocol: adapter.ProtocolVersion, OutputSchema: adapter.OutputSchemaVersion,
		ObservationKinds: []string{"static-analysis"},
		Compatibility:    adapter.Compatibility{EngineAPIs: []string{model.EngineVersion}},
		Tool: adapter.Tool{
			Name: "analysis-fixture", Version: "1.0.0", Upstream: "https://example.com/analysis-fixture",
			Formats: []adapter.ToolFormat{{Name: "fixture-json", Versions: []string{"1.0"}}},
		},
		Limitations: []string{"Test-only fixture; it does not execute a real analyzer."}, Runner: "oci",
		Image:   "registry.example/prc/analysis-fixture@sha256:" + strings.Repeat("a", 64),
		Command: []string{"/adapter", "scan"},
		Capabilities: adapter.Capabilities{
			ReadWorkspace: true, WriteScratch: true, Network: "deny",
			NetworkHosts: []string{}, SecretHandles: []string{}, ChildProcesses: false,
		},
		Resources: adapter.Resources{
			TimeoutSeconds: 60, MemoryMB: 512, CPUs: 1, PIDs: 64, TmpfsMB: 64, Limits: limits,
		},
	}
}

func boundAnalysisExecution(t *testing.T, item model.Inventory, outcome, status string) model.AdapterExecution {
	t.Helper()
	manifest := analysisManifest()
	started := time.Date(2026, 8, 23, 12, 0, 0, 0, time.UTC)
	reason := ""
	if status != "completed" {
		reason = "fixture did not complete"
	}
	output := adapter.RunOutput{
		Transcript: adapter.Transcript{
			Logs: []adapter.Log{}, Artifacts: []adapter.Artifact{},
			Observations: []adapter.Observation{{
				ID: "OBS-ANALYSIS-1", Kind: "static-analysis", Outcome: outcome,
				Summary: "Bound static-analysis observation.", Locations: []adapter.Location{},
			}},
			Summary: adapter.Summary{Type: "summary", Status: status, Reason: reason, Counts: map[string]int{"observations": 1}},
		},
		StartedAt: started, CompletedAt: started.Add(time.Second), DurationMS: 1000,
		DiagnosticsSHA256: strings.Repeat("0", 64), DiagnosticsBytes: 0,
	}
	execution, err := adapter.BindExecution(strings.Repeat("d", 64), adapter.Subject{
		TargetName: item.TargetName, TargetCommit: item.GitCommit, InventoryDigest: item.Digest,
	}, manifest, output)
	if err != nil {
		t.Fatal(err)
	}
	return execution
}

func scannerWithAnalysisBinding(t *testing.T, item model.Inventory) *Engine {
	t.Helper()
	result := scanner(t)
	manifestDigest, err := adapter.ManifestDigest(analysisManifest())
	if err != nil {
		t.Fatal(err)
	}
	assertion := result.Catalog.Assertions["PRC-A-CORE-013"]
	assertion.Parameters = map[string]any{
		"adapter_bindings": []any{map[string]any{
			"adapter_id": analysisManifest().ID, "manifest_sha256": manifestDigest, "observation_kind": "static-analysis",
		}},
	}
	result.Catalog.Assertions[assertion.ID] = assertion
	authorized, err := result.AuthorizesAdapter("prc/core-repository", item, analysisManifest().ID, manifestDigest)
	if err != nil || !authorized {
		t.Fatalf("configured adapter authorization = %t, %v", authorized, err)
	}
	return result
}

func TestBoundAdapterEvidenceCanPassOrFailButAdapterCannotDeclareAssessment(t *testing.T) {
	item, err := inventory.Build(healthyRepository(t))
	if err != nil {
		t.Fatal(err)
	}
	engine := scannerWithAnalysisBinding(t, item)
	passRun, err := engine.ScanWithAdapterEvidence("prc/core-repository", item,
		[]model.AdapterExecution{boundAnalysisExecution(t, item, "not_found", "completed")})
	if err != nil {
		t.Fatal(err)
	}
	pass := findResult(t, passRun, "PRC-A-CORE-013")
	if pass.Assessment != "pass" || pass.Execution != "completed" || len(pass.EvidenceObserved) != 1 ||
		pass.EvidenceObserved[0].Authority != "executed" {
		t.Fatalf("adapter-backed pass = %+v", pass)
	}

	failRun, err := engine.ScanWithAdapterEvidence("prc/core-repository", item,
		[]model.AdapterExecution{boundAnalysisExecution(t, item, "found", "completed")})
	if err != nil {
		t.Fatal(err)
	}
	failure := findResult(t, failRun, "PRC-A-CORE-013")
	if failure.Assessment != "fail" || failRun.TerminalState != "no_go" {
		t.Fatalf("adapter-backed failure = %+v terminal=%s", failure, failRun.TerminalState)
	}
}

func TestAdapterDAGIsBlockedInInspectModeAndReadyInVerifyLocal(t *testing.T) {
	item, err := inventory.Build(healthyRepository(t))
	if err != nil {
		t.Fatal(err)
	}
	engine := scannerWithAnalysisBinding(t, item)
	inspect, err := engine.PlanMode("prc/core-repository", item, ExecutionModeInspect)
	if err != nil {
		t.Fatal(err)
	}
	if len(inspect.Adapters) != 1 || inspect.Adapters[0].Status != "blocked" {
		t.Fatalf("inspect plan adapter = %+v", inspect.Adapters)
	}
	verify, err := engine.PlanMode("prc/core-repository", item, ExecutionModeVerifyLocal)
	if err != nil {
		t.Fatal(err)
	}
	if len(verify.Adapters) != 1 || verify.Adapters[0].Status != "authorized" ||
		verify.CapabilityPolicy.Process != "oci" || !verify.CapabilityPolicy.WriteScratch {
		t.Fatalf("verify-local plan adapter = %+v policy=%+v", verify.Adapters, verify.CapabilityPolicy)
	}
	adapterSeen := false
	for _, node := range verify.Nodes {
		if node.Kind == "adapter" {
			adapterSeen = true
		}
		if node.AssertionID == "PRC-A-CORE-013" && (!adapterSeen || node.Status != "ready" || len(node.DependsOn) != 2) {
			t.Fatalf("analysis assertion was not dependency-ordered after its adapter: %+v", node)
		}
	}
}

func TestIncompleteOrConflictingAdapterEvidenceNeverPasses(t *testing.T) {
	item, err := inventory.Build(healthyRepository(t))
	if err != nil {
		t.Fatal(err)
	}
	engine := scannerWithAnalysisBinding(t, item)
	partial, err := engine.ScanWithAdapterEvidence("prc/core-repository", item,
		[]model.AdapterExecution{boundAnalysisExecution(t, item, "not_found", "partial")})
	if err != nil {
		t.Fatal(err)
	}
	if result := findResult(t, partial, "PRC-A-CORE-013"); result.Assessment != "unknown" || result.Execution != "blocked" {
		t.Fatalf("partial adapter evidence = %+v", result)
	}

	unsupported, err := engine.ScanWithAdapterEvidence("prc/core-repository", item,
		[]model.AdapterExecution{boundAnalysisExecution(t, item, "unsupported", "completed")})
	if err != nil {
		t.Fatal(err)
	}
	if result := findResult(t, unsupported, "PRC-A-CORE-013"); result.Assessment != "unknown" {
		t.Fatalf("unsupported adapter evidence = %+v", result)
	}
}

func TestAdapterEvidenceIsBoundToExactInventory(t *testing.T) {
	item, err := inventory.Build(healthyRepository(t))
	if err != nil {
		t.Fatal(err)
	}
	execution := boundAnalysisExecution(t, item, "not_found", "completed")
	engine := scannerWithAnalysisBinding(t, item)
	other := item
	other.Digest = strings.Repeat("e", 64)
	if _, err := engine.ScanWithAdapterEvidence(
		"prc/core-repository", other, []model.AdapterExecution{execution}); err == nil {
		t.Fatal("adapter execution was reused for a different inventory")
	}
}

func TestEngineRejectsValidButUnconfiguredAdapterExecution(t *testing.T) {
	item, err := inventory.Build(healthyRepository(t))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := scanner(t).ScanWithAdapterEvidence("prc/core-repository", item,
		[]model.AdapterExecution{boundAnalysisExecution(t, item, "not_found", "completed")}); err == nil {
		t.Fatal("valid but profile-unconfigured adapter execution was accepted")
	}
}
