package main

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/MarinJursic/production-readiness-checklist/scanner/adapter"
	"github.com/MarinJursic/production-readiness-checklist/scanner/adapterfixture"
	"github.com/MarinJursic/production-readiness-checklist/scanner/benchmark"
	"github.com/MarinJursic/production-readiness-checklist/scanner/catalog"
	"github.com/MarinJursic/production-readiness-checklist/scanner/engine"
	"github.com/MarinJursic/production-readiness-checklist/scanner/exception"
	"github.com/MarinJursic/production-readiness-checklist/scanner/invalidation"
	"github.com/MarinJursic/production-readiness-checklist/scanner/inventory"
	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
	"github.com/MarinJursic/production-readiness-checklist/scanner/pack"
	"github.com/MarinJursic/production-readiness-checklist/scanner/provider"
	"github.com/MarinJursic/production-readiness-checklist/scanner/remediation"
	"github.com/MarinJursic/production-readiness-checklist/scanner/trust"
)

func commandFinding(t *testing.T, target, assertionID string) model.Finding {
	t.Helper()
	item, err := inventory.Build(target)
	if err != nil {
		t.Fatal(err)
	}
	c, err := catalog.Load(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	run, err := engine.New(c).Scan("prc/core-repository", item)
	if err != nil {
		t.Fatal(err)
	}
	for _, finding := range run.Findings {
		if finding.AssertionID == assertionID {
			return finding
		}
	}
	t.Fatalf("missing finding for %s", assertionID)
	return model.Finding{}
}

func TestVersionCommand(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if code := run([]string{"version"}, &stdout, &stderr); code != 0 {
		t.Fatalf("exit=%d stderr=%s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "prc 0.1.0-dev (revision unknown, built unknown, go1.") {
		t.Fatalf("unexpected output %q", stdout.String())
	}
	stdout.Reset()
	if code := run([]string{"version", "--format", "json"}, &stdout, &stderr); code != 0 {
		t.Fatalf("json exit=%d stderr=%s", code, stderr.String())
	}
	var information versionInformation
	if err := json.Unmarshal(stdout.Bytes(), &information); err != nil {
		t.Fatal(err)
	}
	if information.SchemaVersion != "prc.version/v0.1" || information.Version != "0.1.0-dev" ||
		information.Revision != "unknown" || information.BuiltAt != "unknown" ||
		!strings.HasPrefix(information.GoVersion, "go1.") {
		t.Fatalf("version information = %+v", information)
	}
	if code := run([]string{"version", "unexpected"}, &stdout, &stderr); code != exitConfiguration {
		t.Fatalf("unexpected argument exit=%d", code)
	}
}

func TestMCPServeCommandCompletesHandshakeAndListsReadOnlyTools(t *testing.T) {
	repository := filepath.Join("..", "..")
	input := strings.NewReader(strings.Join([]string{
		`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-11-25","capabilities":{},"clientInfo":{"name":"cli-test","version":"1.0.0"}}}`,
		`{"jsonrpc":"2.0","method":"notifications/initialized"}`,
		`{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}`,
	}, "\n") + "\n")
	var stdout, stderr bytes.Buffer
	code := runWithInput([]string{
		"mcp", "serve", "--catalog-root", repository, "--target", repository,
		"--profile", "prc/core-repository",
	}, input, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit=%d stderr=%s", code, stderr.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("MCP server wrote unexpected stderr: %s", stderr.String())
	}
	lines := strings.Split(strings.TrimSpace(stdout.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("MCP responses = %d: %s", len(lines), stdout.String())
	}
	if !strings.Contains(lines[0], `"protocolVersion":"2025-11-25"`) ||
		!strings.Contains(lines[1], `"name":"prc_scan"`) ||
		!strings.Contains(lines[1], `"readOnlyHint":true`) ||
		strings.Contains(stdout.String(), "Production Readiness Scanner\n") {
		t.Fatalf("unexpected MCP output: %s", stdout.String())
	}
}

func TestMCPServeRejectsUnboundArguments(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := runWithInput([]string{"mcp", "serve", "unexpected"}, strings.NewReader(""), &stdout, &stderr)
	if code != exitConfiguration || !strings.Contains(stderr.String(), "unexpected mcp serve arguments") || stdout.Len() != 0 {
		t.Fatalf("exit=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
}

func TestConfigValidateCommand(t *testing.T) {
	path := filepath.Join("..", "..", "fixtures", "config", "production-readiness.yaml")
	var stdout, stderr bytes.Buffer
	if code := run([]string{"config", "validate", "--file", path, "--format", "json"}, &stdout, &stderr); code != 0 {
		t.Fatalf("exit=%d stderr=%s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), `"schema_version": "prc.config-validation/v0.1"`) ||
		!strings.Contains(stdout.String(), `"network": "deny"`) {
		t.Fatalf("unexpected output %s", stdout.String())
	}
}

func TestCatalogValidateAndBundleCommands(t *testing.T) {
	root := filepath.Join("..", "..")
	var manifestOutput, stderr bytes.Buffer
	if code := run([]string{"catalog", "validate", "--catalog-root", root, "--format", "json"}, &manifestOutput, &stderr); code != 0 {
		t.Fatalf("validate exit=%d stderr=%s", code, stderr.String())
	}
	var manifest struct {
		SchemaVersion  string `json:"schema_version"`
		CatalogDigest  string `json:"catalog_digest"`
		ObjectiveCount int    `json:"objective_count"`
		AssertionCount int    `json:"assertion_count"`
		ProfileCount   int    `json:"profile_count"`
	}
	if err := json.Unmarshal(manifestOutput.Bytes(), &manifest); err != nil {
		t.Fatal(err)
	}
	if manifest.SchemaVersion != "prc.catalog-manifest/v0.1" || len(manifest.CatalogDigest) != 64 ||
		manifest.ObjectiveCount == 0 || manifest.AssertionCount == 0 || manifest.ProfileCount == 0 {
		t.Fatalf("invalid catalog manifest: %+v", manifest)
	}
	var bundleOutput bytes.Buffer
	stderr.Reset()
	if code := run([]string{"catalog", "bundle", "--catalog-root", root}, &bundleOutput, &stderr); code != 0 {
		t.Fatalf("bundle exit=%d stderr=%s", code, stderr.String())
	}
	var bundle struct {
		SchemaVersion string `json:"schema_version"`
		Manifest      struct {
			CatalogDigest string `json:"catalog_digest"`
		} `json:"manifest"`
		Objectives []json.RawMessage `json:"objectives"`
		Assertions []json.RawMessage `json:"assertions"`
		Profiles   []json.RawMessage `json:"profiles"`
	}
	if err := json.Unmarshal(bundleOutput.Bytes(), &bundle); err != nil {
		t.Fatal(err)
	}
	if bundle.SchemaVersion != "prc.catalog-bundle/v0.1" ||
		bundle.Manifest.CatalogDigest != manifest.CatalogDigest ||
		len(bundle.Objectives) != manifest.ObjectiveCount ||
		len(bundle.Assertions) != manifest.AssertionCount || len(bundle.Profiles) != manifest.ProfileCount {
		t.Fatalf("bundle does not match manifest: %+v", bundle)
	}
}

func TestBenchmarkRunCommand(t *testing.T) {
	repository := filepath.Join("..", "..")
	var stdout, stderr bytes.Buffer
	code := run([]string{
		"benchmark", "run", "--catalog-root", repository,
		"--suite", filepath.Join(repository, "fixtures", "benchmarks", "core-native", "suite.yaml"),
		"--evaluated-at", "2026-08-23T12:00:00Z", "--format", "json",
	}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit=%d stderr=%s", code, stderr.String())
	}
	var report benchmark.Report
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
		t.Fatal(err)
	}
	if report.SchemaVersion != benchmark.ReportSchema || !report.Passed || report.Summary.Expectations != 10 {
		t.Fatalf("benchmark report = %+v", report)
	}

	directory := t.TempDir()
	target := filepath.Join(directory, "target")
	if err := os.Mkdir(target, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(target, "README.md"), []byte("present\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	suite := `schema_version: prc.benchmark-suite/v0.1
id: prc.benchmark.cli-mismatch@0.1
title: CLI mismatch
profile_id: prc/core-repository
cases:
  - id: mismatch
    target: target
    expectations:
      - assertion_id: PRC-A-CORE-001
        assessment: fail
        execution: completed
quality_budget:
  minimum_precision: 1
  minimum_recall: 1
  maximum_false_positive_rate: 0
  maximum_mismatches: 0
  require_determinism: true
`
	suitePath := filepath.Join(directory, "suite.yaml")
	if err := os.WriteFile(suitePath, []byte(suite), 0o600); err != nil {
		t.Fatal(err)
	}
	stdout.Reset()
	stderr.Reset()
	code = run([]string{
		"benchmark", "run", "--catalog-root", repository, "--suite", suitePath,
		"--evaluated-at", "2026-08-23T12:00:00Z", "--format", "json",
	}, &stdout, &stderr)
	if code != exitGateFailed {
		t.Fatalf("mismatch exit=%d stderr=%s", code, stderr.String())
	}
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil || report.Passed {
		t.Fatalf("mismatch report=%+v err=%v", report, err)
	}
}

func TestPackValidateCommand(t *testing.T) {
	repository := filepath.Join("..", "..")
	var stdout, stderr bytes.Buffer
	code := run([]string{
		"pack", "validate", "--catalog-root", repository,
		"--file", filepath.Join(repository, "packs", "core-foundation.yaml"), "--format", "json",
	}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit=%d stderr=%s", code, stderr.String())
	}
	var report pack.Report
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
		t.Fatal(err)
	}
	if report.SchemaVersion != pack.ReportSchema || len(report.Digest) != 64 || len(report.Manifest.Assertions) != 3 {
		t.Fatalf("pack report = %+v", report)
	}
}

func TestPackAndRegistryPublisherVerificationCommands(t *testing.T) {
	repository := filepath.Join("..", "..")
	catalogValue, err := catalog.Load(repository)
	if err != nil {
		t.Fatal(err)
	}
	loadedPack, err := pack.Load(repository, filepath.Join(repository, "packs", "core-foundation.yaml"), catalogValue)
	if err != nil {
		t.Fatal(err)
	}
	registryPath := filepath.Join(repository, "fixtures", "adapters", "fixture-registry.yaml")
	registry, err := adapter.LoadRegistry(registryPath)
	if err != nil {
		t.Fatal(err)
	}
	issuedAt := time.Date(2026, 8, 23, 12, 0, 0, 0, time.UTC)
	storePath, packSignature := writePublisherFixture(
		t, "pack", loadedPack.Manifest.ID, loadedPack.Digest, issuedAt,
	)
	_, registrySignature := writePublisherFixtureWithStore(
		t, storePath, "adapter-registry", registry.ID, registry.Digest, issuedAt,
	)

	tests := []struct {
		name string
		args []string
	}{
		{
			name: "pack",
			args: []string{
				"pack", "verify", "--catalog-root", repository,
				"--file", filepath.Join(repository, "packs", "core-foundation.yaml"),
				"--trust-store", storePath, "--signature", packSignature,
				"--verified-at", "2026-08-23T13:00:00Z", "--format", "json",
			},
		},
		{
			name: "registry",
			args: []string{
				"adapter", "registry-verify", "--file", registryPath,
				"--trust-store", storePath, "--signature", registrySignature,
				"--verified-at", "2026-08-23T13:00:00Z", "--format", "json",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			if code := run(test.args, &stdout, &stderr); code != 0 {
				t.Fatalf("exit=%d stderr=%s", code, stderr.String())
			}
			var verification trust.Verification
			if err := json.Unmarshal(stdout.Bytes(), &verification); err != nil {
				t.Fatal(err)
			}
			if !verification.Verified || verification.SchemaVersion != trust.VerificationSchema ||
				len(verification.TrustStoreDigest) != 64 {
				t.Fatalf("verification = %+v", verification)
			}
		})
	}

	var stdout, stderr bytes.Buffer
	code := run([]string{
		"pack", "verify", "--catalog-root", repository,
		"--file", filepath.Join(repository, "packs", "core-foundation.yaml"),
		"--trust-store", storePath, "--signature", registrySignature,
		"--verified-at", "2026-08-23T13:00:00Z", "--format", "json",
	}, &stdout, &stderr)
	if code != exitPolicyDenied || !strings.Contains(stderr.String(), "subject") {
		t.Fatalf("mismatched signature exit=%d stderr=%s", code, stderr.String())
	}
}

func TestSignedRiskExceptionVerificationCommand(t *testing.T) {
	repository := filepath.Join("..", "..")
	target, stateDirectory := t.TempDir(), t.TempDir()
	if err := os.Chmod(stateDirectory, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(target, "fixture.txt"), []byte("fixture\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	code := run([]string{
		"scan", "--target", target, "--catalog-root", repository,
		"--state-dir", stateDirectory, "--exit-policy", "never", "--format", "json",
	}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("scan exit=%d stderr=%s", code, stderr.String())
	}
	var scanned model.RunResult
	if err := json.Unmarshal(stdout.Bytes(), &scanned); err != nil {
		t.Fatal(err)
	}
	if len(scanned.Findings) == 0 {
		t.Fatal("fixture scan has no finding")
	}
	finding := scanned.Findings[0]
	approvedAt := scanned.CompletedAt.Add(time.Hour)
	record := exception.Record{
		SchemaVersion: exception.Schema, ID: "PRC-EXC-CLI-001", Status: "approved",
		Run: exception.RunBinding{
			RunID: scanned.RunID, InventoryDigest: scanned.Inventory.Digest,
			ProfileID: scanned.Plan.ProfileID, ProfileVersion: scanned.Plan.ProfileVersion,
			TargetName: scanned.Inventory.TargetName, TargetCommit: scanned.Inventory.GitCommit,
			ProjectID: scanned.Plan.ProjectID, ArtifactDigests: scanned.Plan.ArtifactDigests,
			TargetEnvironments: scanned.Plan.TargetEnvironments,
		},
		Finding: exception.FindingBinding{
			FindingID: finding.ID, FindingFingerprint: finding.Fingerprint,
			AssertionID: finding.AssertionID, ControlIDs: finding.ControlIDs,
		},
		RequestedBy: exception.Actor{ID: "requester", Name: "Requesting engineer", Authority: "engineering"},
		RiskOwner:   exception.Actor{ID: "risk-owner", Name: "Risk owner", Authority: "executive"},
		Reviewers:   []exception.Actor{{ID: "security-reviewer", Name: "Security reviewer", Authority: "security"}},
		Risk: exception.RiskAnalysis{
			Title: "CLI fixture exception", Rationale: "Exercise the signed exception command.",
			Likelihood: "unlikely", Impact: "high", WorstCredibleOutcome: "A release defect remains until remediation.",
		},
		CompensatingControls: []exception.CompensatingControl{{
			Description:        "A temporary reviewed control limits exposure.",
			EvidenceReferences: []string{strings.Repeat("b", 64)},
		}},
		Monitoring: exception.Monitoring{Owner: "operations", Signal: "Alert on the affected behavior.", Response: "Disable the affected feature."},
		Remediation: exception.Remediation{
			Owner: "engineering", Plan: "Implement and verify the missing control.", DueAt: approvedAt.Add(24 * time.Hour),
		},
		ApprovedAt: approvedAt, ExpiresAt: approvedAt.Add(7 * 24 * time.Hour),
	}
	directory := t.TempDir()
	recordPath := filepath.Join(directory, "exception.json")
	writeJSONFixture(t, recordPath, record)
	loaded, err := exception.Load(recordPath)
	if err != nil {
		t.Fatal(err)
	}
	storePath, signaturePath := writePublisherFixture(
		t, "risk-exception", loaded.Record.ID, loaded.Digest, approvedAt,
	)
	stdout.Reset()
	stderr.Reset()
	code = run([]string{
		"exception", "verify", "--file", recordPath, "--state-dir", stateDirectory,
		"--trust-store", storePath, "--signature", signaturePath,
		"--verified-at", approvedAt.Add(time.Hour).Format(time.RFC3339Nano), "--format", "json",
	}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exception verify exit=%d stderr=%s", code, stderr.String())
	}
	var verification exception.Verification
	if err := json.Unmarshal(stdout.Bytes(), &verification); err != nil {
		t.Fatal(err)
	}
	if verification.Disposition != "accepted_risk_exception" || !verification.Signature.Verified ||
		!strings.Contains(verification.GateEffect, "unchanged") {
		t.Fatalf("exception verification = %+v", verification)
	}
}

func writePublisherFixture(t *testing.T, kind, artifactID, digest string, issuedAt time.Time) (string, string) {
	t.Helper()
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	directory := t.TempDir()
	storePath := filepath.Join(directory, "trust-store.json")
	store := trust.Store{
		SchemaVersion: trust.StoreSchema, ID: "test-release-keys", Revision: 1,
		Keys: []trust.Key{{
			ID: "test-release", Algorithm: trust.AlgorithmEd25519,
			PublicKey: base64.StdEncoding.EncodeToString(publicKey), Scopes: []string{"adapter-registry", "pack", "risk-exception"},
			Status: "active", NotBefore: issuedAt.Add(-time.Hour), NotAfter: issuedAt.Add(24 * time.Hour),
		}},
	}
	writeJSONFixture(t, storePath, store)
	privateKeyPath := filepath.Join(directory, "private-key")
	if err := os.WriteFile(privateKeyPath, privateKey, 0o600); err != nil {
		t.Fatal(err)
	}
	return writePublisherFixtureWithStore(t, storePath, kind, artifactID, digest, issuedAt)
}

func writePublisherFixtureWithStore(t *testing.T, storePath, kind, artifactID, digest string, issuedAt time.Time) (string, string) {
	t.Helper()
	privateKeyPath := filepath.Join(filepath.Dir(storePath), "private-key")
	privateKey, err := os.ReadFile(privateKeyPath)
	if err != nil {
		t.Fatal(err)
	}
	signature := trust.Signature{
		SchemaVersion: trust.SignatureSchema, ArtifactKind: kind, ArtifactID: artifactID,
		SHA256: digest, KeyID: "test-release", Algorithm: trust.AlgorithmEd25519, IssuedAt: issuedAt,
	}
	payload, err := trust.SigningPayload(signature)
	if err != nil {
		t.Fatal(err)
	}
	signature.Value = base64.StdEncoding.EncodeToString(ed25519.Sign(ed25519.PrivateKey(privateKey), payload))
	path := filepath.Join(filepath.Dir(storePath), strings.ReplaceAll(kind, "-", "_")+"-signature.json")
	writeJSONFixture(t, path, signature)
	return storePath, path
}

func writeJSONFixture(t *testing.T, path string, value any) {
	t.Helper()
	payload, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, payload, 0o600); err != nil {
		t.Fatal(err)
	}
}

func TestConfigValidateRejectsCapabilityExpansion(t *testing.T) {
	path := filepath.Join(t.TempDir(), "production-readiness.yaml")
	data, err := os.ReadFile(filepath.Join("..", "..", "fixtures", "config", "production-readiness.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	data = bytes.Replace(data, []byte("network: deny"), []byte("network: allow"), 1)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	if code := run([]string{"config", "validate", "--file", path}, &stdout, &stderr); code != exitConfiguration {
		t.Fatalf("exit=%d stderr=%s", code, stderr.String())
	}
}

func TestPlanConsumesValidatedConfiguration(t *testing.T) {
	target := t.TempDir()
	if err := os.WriteFile(filepath.Join(target, "app.py"), []byte("print('ready')\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	configPath := filepath.Join("..", "..", "fixtures", "config", "production-readiness.yaml")
	var stdout, stderr bytes.Buffer
	code := run([]string{
		"plan", "--target", target, "--catalog-root", filepath.Join("..", ".."),
		"--config", configPath, "--format", "json",
	}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit=%d stderr=%s", code, stderr.String())
	}
	var plan model.Plan
	if err := json.Unmarshal(stdout.Bytes(), &plan); err != nil {
		t.Fatal(err)
	}
	if plan.SchemaVersion != "prc.plan/v0.6" || plan.EngineVersion != "prc.engine/v0.1" ||
		len(plan.ProfileDigest) != 64 || len(plan.CatalogDigest) != 64 || len(plan.ConfigurationDigest) != 64 ||
		plan.ProjectID != "example-product" || !slices.Equal(plan.TargetEnvironments, []string{"staging"}) ||
		plan.ExecutionMode != "inspect" || len(plan.Nodes) != len(plan.Assertions)+len(plan.Adapters)+2 {
		t.Fatalf("configured plan = %+v", plan)
	}
}

func TestAdapterValidateOutputCommand(t *testing.T) {
	var stdout, stderr bytes.Buffer
	path := filepath.Join("..", "..", "fixtures", "adapters", "valid-output.jsonl")
	if code := run([]string{"adapter", "validate-output", "--file", path}, &stdout, &stderr); code != 0 {
		t.Fatalf("exit=%d stderr=%s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), `"status": "completed"`) {
		t.Fatalf("unexpected output %q", stdout.String())
	}
}

func TestAdapterRegistryValidateCommand(t *testing.T) {
	var stdout, stderr bytes.Buffer
	path := filepath.Join("..", "..", "fixtures", "adapters", "fixture-registry.yaml")
	if code := run([]string{"adapter", "registry-validate", "--file", path, "--format", "json"}, &stdout, &stderr); code != 0 {
		t.Fatalf("exit=%d stderr=%s", code, stderr.String())
	}
	var report adapter.RegistryReport
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
		t.Fatal(err)
	}
	if report.SchemaVersion != adapter.RegistryReportSchema || len(report.Digest) != 64 ||
		len(report.Registry.Entries) != 1 || report.Registry.Entries[0].Status != "active" {
		t.Fatalf("registry output = %+v", report)
	}
}

func TestAdapterValidateOutputRejectsAuthorityAttack(t *testing.T) {
	var stdout, stderr bytes.Buffer
	path := filepath.Join("..", "..", "fixtures", "adapters", "malicious-authority-output.jsonl")
	if code := run([]string{"adapter", "validate-output", "--file", path}, &stdout, &stderr); code != exitExecution {
		t.Fatalf("exit=%d stderr=%s", code, stderr.String())
	}
}

func TestAdapterValidateOutputRejectsUndeclaredObservationKind(t *testing.T) {
	path := filepath.Join(t.TempDir(), "output.jsonl")
	data := []byte("{\"type\":\"observation\",\"observation\":{\"id\":\"OBS-1\",\"kind\":\"undeclared\",\"outcome\":\"not_found\",\"summary\":\"fixture\",\"locations\":[]}}\n{\"type\":\"summary\",\"status\":\"completed\",\"counts\":{\"observations\":1}}\n")
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	if code := run([]string{
		"adapter", "validate-output", "--manifest", filepath.Join("..", "..", "fixtures", "adapters", "fixture-adapter.yaml"),
		"--file", path,
	}, &stdout, &stderr); code != exitExecution || !strings.Contains(stderr.String(), "undeclared observation kind") {
		t.Fatalf("exit=%d stderr=%s", code, stderr.String())
	}
}

func TestAdapterFixtureValidateCommand(t *testing.T) {
	var stdout, stderr bytes.Buffer
	path := filepath.Join("..", "..", "fixtures", "adapters", "fixture-suite.yaml")
	if code := run([]string{"adapter", "fixture-validate", "--suite", path, "--format", "json"}, &stdout, &stderr); code != 0 {
		t.Fatalf("exit=%d stderr=%s", code, stderr.String())
	}
	var report adapterfixture.Report
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
		t.Fatal(err)
	}
	if !report.Passed || report.Summary.Cases != 7 || report.Summary.Mismatched != 0 {
		t.Fatalf("unexpected fixture report: %+v", report)
	}
}

func TestProviderCapabilitiesAreReadOnly(t *testing.T) {
	for _, name := range []string{"codex", "claude"} {
		t.Run(name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			if code := run([]string{"provider", "capabilities", "--provider", name}, &stdout, &stderr); code != 0 {
				t.Fatalf("exit=%d stderr=%s", code, stderr.String())
			}
			for _, forbidden := range []string{`"workspace_mutation": true`, `"network_tools": true`, `"shell": true`} {
				if strings.Contains(stdout.String(), forbidden) {
					t.Fatalf("unsafe capability in %s", stdout.String())
				}
			}
		})
	}
}

func TestProviderRunEmitsFailureRecordBeforeExecutionExit(t *testing.T) {
	workspace := t.TempDir()
	if err := os.WriteFile(filepath.Join(workspace, "app.go"), []byte("package app\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	item, err := inventory.Build(workspace)
	if err != nil {
		t.Fatal(err)
	}
	draft := provider.Task{
		SchemaVersion: provider.TaskSchema, Mode: "suggest",
		FindingID: strings.Repeat("a", 64), FindingFingerprint: strings.Repeat("b", 64),
		AssertionID: "PRC-A-CORE-010", ControlIDs: []string{"USEQ-12655775"},
		Goal: "Add one focused test.", ReadScope: "task-inputs", RelevantPaths: []string{"app.go"},
		Inputs: []provider.InputFile{}, AllowedPaths: []string{"app_test.go"},
		ProtectedPaths: []string{".git/", "catalog/", "schemas/"}, AllowedCommands: [][]string{},
		Network: "deny", Secrets: "none", AllowRemoteSourceProcessing: true,
		TimeoutSeconds: 30, MaxOutputBytes: 64 * 1024,
	}
	task, err := provider.SealTaskValueWithInventory(draft, workspace, item, nil)
	if err != nil {
		t.Fatal(err)
	}
	taskData, _ := json.Marshal(task)
	taskPath := filepath.Join(t.TempDir(), "task.json")
	if err := os.WriteFile(taskPath, taskData, 0o600); err != nil {
		t.Fatal(err)
	}
	executable := filepath.Join(t.TempDir(), "codex")
	if err := os.WriteFile(executable, []byte("#!/bin/sh\nprintf diagnostic >&2\nexit 9\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	outputDirectory := t.TempDir()
	if err := os.Chmod(outputDirectory, 0o700); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	code := run([]string{
		"provider", "run", "--provider", "codex", "--task", taskPath,
		"--workspace", workspace, "--output-dir", outputDirectory,
		"--output-schema", filepath.Join("..", "..", "schemas", "agent-output.schema.json"),
		"--executable", executable,
	}, &stdout, &stderr)
	if code != exitExecution || !strings.Contains(stderr.String(), "PRC-EXIT-4") {
		t.Fatalf("exit=%d stderr=%s stdout=%s", code, stderr.String(), stdout.String())
	}
	var failure provider.Failure
	if err := json.Unmarshal(stdout.Bytes(), &failure); err != nil || failure.Validate() != nil ||
		failure.ReasonCode != "process_failed" || !failure.TranscriptsComplete {
		t.Fatalf("provider failure=%+v decode=%v", failure, err)
	}
}

func TestScanReportFormats(t *testing.T) {
	target := t.TempDir()
	expected := map[string]string{
		"markdown": "# Production readiness assessment",
		"html":     "<!doctype html>",
		"sarif":    `"version": "2.1.0"`,
		"junit":    "<testsuite ",
	}
	for format, marker := range expected {
		t.Run(format, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code := run([]string{
				"scan", "--target", target, "--catalog-root", filepath.Join("..", ".."),
				"--format", format, "--exit-policy", "never",
			}, &stdout, &stderr)
			if code != 0 || !strings.Contains(stdout.String(), marker) {
				t.Fatalf("exit=%d stderr=%s output=%s", code, stderr.String(), stdout.String())
			}
		})
	}
}

func TestSimpleScanDiscoversCatalogCreatesPrivateReportAndNeverFixesTarget(t *testing.T) {
	repository, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	workingDirectory, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(repository); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(workingDirectory) })

	cacheRoot := t.TempDir()
	previousCacheDirectory := userCacheDirectory
	userCacheDirectory = func() (string, error) { return cacheRoot, nil }
	t.Cleanup(func() { userCacheDirectory = previousCacheDirectory })

	target := t.TempDir()
	sourcePath := filepath.Join(target, "app.py")
	before := []byte("print('ready')")
	if err := os.WriteFile(sourcePath, before, 0o666); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(sourcePath, 0o666); err != nil {
		t.Fatal(err)
	}
	beforeInfo, err := os.Stat(sourcePath)
	if err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	code := run([]string{"scan", target, "--exit-policy", "never"}, &stdout, &stderr)
	if code != 0 || stderr.Len() != 0 {
		t.Fatalf("exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	if !strings.Contains(stdout.String(), "Scan mode: report only; no fixes were applied.") ||
		!strings.Contains(stdout.String(), "Detailed report: ") {
		t.Fatalf("simple scan output = %s", stdout.String())
	}
	reportDirectory := filepath.Join(cacheRoot, "prc", "reports")
	entries, err := os.ReadDir(reportDirectory)
	if err != nil || len(entries) != 1 {
		t.Fatalf("reports=%v err=%v", entries, err)
	}
	reportPath := filepath.Join(reportDirectory, entries[0].Name())
	reportInfo, err := os.Stat(reportPath)
	if err != nil || reportInfo.Mode().Perm() != 0o600 {
		t.Fatalf("report mode=%v err=%v", reportInfo.Mode().Perm(), err)
	}
	reportBody, err := os.ReadFile(reportPath)
	if err != nil || !bytes.Contains(reportBody, []byte("Detailed findings")) ||
		!bytes.Contains(reportBody, []byte("All assertion results")) {
		t.Fatalf("detailed report is incomplete: %v", err)
	}
	after, err := os.ReadFile(sourcePath)
	if err != nil || !bytes.Equal(after, before) {
		t.Fatalf("scan changed target bytes: %q (%v)", after, err)
	}
	afterInfo, err := os.Stat(sourcePath)
	if err != nil || afterInfo.Mode().Perm() != beforeInfo.Mode().Perm() {
		t.Fatalf("scan changed target mode: before=%v after=%v err=%v", beforeInfo.Mode().Perm(), afterInfo.Mode().Perm(), err)
	}
}

func TestScanCustomReportNeverOverwrites(t *testing.T) {
	target := t.TempDir()
	reportPath := filepath.Join(t.TempDir(), "readiness.html")
	args := []string{
		"scan", "--catalog-root", filepath.Join("..", ".."), "--report", reportPath,
		"--exit-policy", "never", target,
	}
	var stdout, stderr bytes.Buffer
	if code := run(args, &stdout, &stderr); code != 0 {
		t.Fatalf("exit=%d stderr=%s", code, stderr.String())
	}
	original, err := os.ReadFile(reportPath)
	if err != nil || len(original) == 0 {
		t.Fatalf("custom report was not created: %v", err)
	}
	stdout.Reset()
	stderr.Reset()
	if code := run(args, &stdout, &stderr); code != exitInternal || !strings.Contains(stderr.String(), "without overwriting") {
		t.Fatalf("overwrite exit=%d stderr=%s", code, stderr.String())
	}
	after, err := os.ReadFile(reportPath)
	if err != nil || !bytes.Equal(after, original) {
		t.Fatalf("existing report changed: %v", err)
	}
}

func TestIaCProfileAloneAuthorizesCheckedInCheckovAdapter(t *testing.T) {
	repository := filepath.Join("..", "..")
	target := t.TempDir()
	if err := os.WriteFile(filepath.Join(target, "main.tf"), []byte("resource \"aws_s3_bucket\" \"logs\" {}\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	item, err := inventory.Build(target)
	if err != nil {
		t.Fatal(err)
	}
	manifest, err := adapter.LoadManifest(filepath.Join(repository, "adapters", "checkov-v3.3.8.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	digest, err := adapter.ManifestDigest(manifest)
	if err != nil {
		t.Fatal(err)
	}
	scanner, err := loadEngine(repository)
	if err != nil {
		t.Fatal(err)
	}
	authorized, err := scanner.AuthorizesAdapter("prc/iac", item, manifest.ID, digest)
	if err != nil || !authorized {
		t.Fatalf("IaC profile Checkov authorization=%t err=%v", authorized, err)
	}
	authorized, err = scanner.AuthorizesAdapter("prc/core-repository", item, manifest.ID, digest)
	if err != nil || authorized {
		t.Fatalf("core profile unexpectedly authorized Checkov=%t err=%v", authorized, err)
	}
}

func copyCatalogFile(t *testing.T, sourceRoot, targetRoot, relative string) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(sourceRoot, filepath.FromSlash(relative)))
	if err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(targetRoot, filepath.FromSlash(relative))
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(target, data, 0o644); err != nil {
		t.Fatal(err)
	}
	return data
}

func fakeDockerRuntime(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "docker")
	script := `#!/bin/sh
printf '%s\n' \
  '{"type":"log","level":"info","message":"fixture executed"}' \
  '{"type":"observation","observation":{"id":"OBS-FIXTURE-001","kind":"fixture-result","outcome":"not_found","summary":"Fixture pattern absent.","locations":[]}}' \
  '{"type":"summary","status":"completed","counts":{"logs":1,"observations":1,"artifacts":0}}'
`
	if err := os.WriteFile(path, []byte(script), 0o700); err != nil {
		t.Fatal(err)
	}
	return path
}

func fixtureRegistry(t *testing.T, manifestPath, manifestDigest, trust, status, reason string) string {
	t.Helper()
	directory := t.TempDir()
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(directory, "adapter.yaml"), data, 0o600); err != nil {
		t.Fatal(err)
	}
	registry := "schema_version: prc.adapter-registry/v0.1\n" +
		"id: prc.adapter-registry.test@0.1\nrevision: 1\nentries:\n" +
		"  - adapter_id: prc.adapter.fixture@0.1\n" +
		"    manifest_path: adapter.yaml\n" +
		"    manifest_sha256: " + manifestDigest + "\n" +
		"    publisher_id: prc-project\n" +
		"    trust: " + trust + "\n" +
		"    status: " + status + "\n"
	if reason != "" {
		registry += "    reason: " + reason + "\n"
	}
	path := filepath.Join(directory, "registry.yaml")
	if err := os.WriteFile(path, []byte(registry), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func replaceCatalogAdapterBindingBlock(t *testing.T, assertions []byte, binding string) []byte {
	t.Helper()
	bindingStart := bytes.Index(assertions, []byte("      adapter_bindings:\n"))
	if bindingStart < 0 {
		t.Fatal("test catalog has no adapter binding block")
	}
	bindingEndOffset := bytes.Index(assertions[bindingStart:], []byte("    severity:"))
	if bindingEndOffset < 0 {
		t.Fatal("test catalog adapter binding block has no severity boundary")
	}
	bindingEnd := bindingStart + bindingEndOffset
	updated := make([]byte, 0, len(assertions)-bindingEndOffset+len(binding))
	updated = append(updated, assertions[:bindingStart]...)
	updated = append(updated, binding...)
	updated = append(updated, assertions[bindingEnd:]...)
	return updated
}

func TestScanConsumesOnlyProfileAuthorizedLiveAdapterEvidence(t *testing.T) {
	repository := filepath.Join("..", "..")
	manifestPath := filepath.Join(repository, "fixtures", "adapters", "fixture-adapter.yaml")
	manifest, err := adapter.LoadManifest(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	manifestDigest, err := adapter.ManifestDigest(manifest)
	if err != nil {
		t.Fatal(err)
	}
	catalogRoot := t.TempDir()
	copyCatalogFile(t, repository, catalogRoot, "catalog/objectives/core-repository.yaml")
	for _, source := range []string{
		"docs/checklists/08-maintenance-vendors-compliance.md",
		"docs/engineering/01-governance-and-foundations.md",
		"docs/engineering/05-code-quality-and-implementation.md",
		"docs/engineering/06-application-services-and-apis.md",
		"docs/engineering/08-security-and-cryptography.md",
		"docs/engineering/10-verification-and-testing.md",
		"docs/engineering/11-developer-experience-platform-and-delivery.md",
		"docs/engineering/16-specialized-domains-and-release-assurance.md",
	} {
		copyCatalogFile(t, repository, catalogRoot, source)
	}
	assertions := copyCatalogFile(t, repository, catalogRoot, "catalog/assertions/core-repository.yaml")
	copyCatalogFile(t, repository, catalogRoot, "catalog/profiles/core-repository.yaml")
	binding := "      adapter_bindings:\n" +
		"        - adapter_id: " + manifest.ID + "\n" +
		"          manifest_sha256: " + manifestDigest + "\n" +
		"          observation_kind: fixture-result\n"
	assertions = replaceCatalogAdapterBindingBlock(t, assertions, binding)
	if !bytes.Contains(assertions, []byte(manifestDigest)) {
		t.Fatal("test catalog binding replacement failed")
	}
	if err := os.WriteFile(filepath.Join(catalogRoot, "catalog", "assertions", "core-repository.yaml"), assertions, 0o644); err != nil {
		t.Fatal(err)
	}
	target := t.TempDir()
	if err := os.WriteFile(filepath.Join(target, "app.py"), []byte("print('ready')\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	code := run([]string{
		"scan", "--target", target, "--catalog-root", catalogRoot,
		"--mode", "verify-local",
		"--adapter-manifest", manifestPath, "--adapter-runtime", fakeDockerRuntime(t),
		"--format", "json", "--exit-policy", "never",
	}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit=%d stderr=%s", code, stderr.String())
	}
	var result model.RunResult
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if len(result.AdapterExecutions) != 1 || result.AdapterExecutions[0].AdapterID != manifest.ID {
		t.Fatalf("bound executions = %+v", result.AdapterExecutions)
	}
	if result.AdapterExecutions[0].Resolution.Source != adapter.ResolutionSourceExplicitLocal ||
		result.AdapterExecutions[0].Resolution.PublisherID != manifest.Publisher.ID {
		t.Fatalf("local resolution = %+v", result.AdapterExecutions[0].Resolution)
	}
	found := false
	for _, assertion := range result.Results {
		if assertion.AssertionID == "PRC-A-CORE-013" {
			found = true
			if assertion.Assessment != "pass" || len(assertion.EvidenceObserved) != 1 {
				t.Fatalf("adapter-backed assertion = %+v", assertion)
			}
		}
	}
	if !found {
		t.Fatal("missing adapter-backed assertion")
	}

	registryPath := fixtureRegistry(t, manifestPath, manifestDigest, "first-party-sandboxed", "active", "")
	stdout.Reset()
	stderr.Reset()
	code = run([]string{
		"scan", "--target", target, "--catalog-root", catalogRoot,
		"--mode", "verify-local", "--adapter-registry", registryPath,
		"--adapter-id", manifest.ID, "--adapter-runtime", fakeDockerRuntime(t),
		"--format", "json", "--exit-policy", "never",
	}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("registry scan exit=%d stderr=%s", code, stderr.String())
	}
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if len(result.AdapterExecutions) != 1 || result.AdapterExecutions[0].ManifestSHA256 != manifestDigest ||
		result.AdapterExecutions[0].Resolution.Source != adapter.ResolutionSourceRegistry ||
		result.AdapterExecutions[0].Resolution.RegistryID != "prc.adapter-registry.test@0.1" ||
		len(result.AdapterExecutions[0].Resolution.RegistryDigest) != 64 {
		t.Fatalf("registry-resolved executions = %+v", result.AdapterExecutions)
	}
}

func TestScanExecutesEveryAuthorizedRepeatedManifest(t *testing.T) {
	repository := filepath.Join("..", "..")
	firstPath := filepath.Join(repository, "fixtures", "adapters", "fixture-adapter.yaml")
	first, err := adapter.LoadManifest(firstPath)
	if err != nil {
		t.Fatal(err)
	}
	firstDigest, err := adapter.ManifestDigest(first)
	if err != nil {
		t.Fatal(err)
	}
	secondData, err := os.ReadFile(firstPath)
	if err != nil {
		t.Fatal(err)
	}
	secondData = bytes.Replace(
		secondData,
		[]byte("id: prc.adapter.fixture@0.1"),
		[]byte("id: prc.adapter.fixture-two@0.1"),
		1,
	)
	secondPath := filepath.Join(t.TempDir(), "second-adapter.yaml")
	if err := os.WriteFile(secondPath, secondData, 0o600); err != nil {
		t.Fatal(err)
	}
	second, err := adapter.LoadManifest(secondPath)
	if err != nil {
		t.Fatal(err)
	}
	secondDigest, err := adapter.ManifestDigest(second)
	if err != nil {
		t.Fatal(err)
	}

	catalogRoot := t.TempDir()
	copyCatalogFile(t, repository, catalogRoot, "catalog/objectives/core-repository.yaml")
	for _, source := range []string{
		"docs/checklists/08-maintenance-vendors-compliance.md",
		"docs/engineering/01-governance-and-foundations.md",
		"docs/engineering/05-code-quality-and-implementation.md",
		"docs/engineering/06-application-services-and-apis.md",
		"docs/engineering/08-security-and-cryptography.md",
		"docs/engineering/10-verification-and-testing.md",
		"docs/engineering/11-developer-experience-platform-and-delivery.md",
		"docs/engineering/16-specialized-domains-and-release-assurance.md",
	} {
		copyCatalogFile(t, repository, catalogRoot, source)
	}
	assertions := copyCatalogFile(t, repository, catalogRoot, "catalog/assertions/core-repository.yaml")
	copyCatalogFile(t, repository, catalogRoot, "catalog/profiles/core-repository.yaml")
	binding := "      adapter_bindings:\n" +
		"        - adapter_id: " + first.ID + "\n" +
		"          manifest_sha256: " + firstDigest + "\n" +
		"          observation_kind: fixture-result\n" +
		"        - adapter_id: " + second.ID + "\n" +
		"          manifest_sha256: " + secondDigest + "\n" +
		"          observation_kind: fixture-result\n"
	assertions = replaceCatalogAdapterBindingBlock(t, assertions, binding)
	if err := os.WriteFile(
		filepath.Join(catalogRoot, "catalog", "assertions", "core-repository.yaml"),
		assertions,
		0o644,
	); err != nil {
		t.Fatal(err)
	}

	target := t.TempDir()
	if err := os.WriteFile(filepath.Join(target, "app.py"), []byte("print('ready')\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	code := run([]string{
		"scan", "--target", target, "--catalog-root", catalogRoot,
		"--mode", "verify-local", "--adapter-manifest", secondPath,
		"--adapter-manifest", firstPath, "--adapter-runtime", fakeDockerRuntime(t),
		"--format", "json", "--exit-policy", "never",
	}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit=%d stderr=%s", code, stderr.String())
	}
	var result model.RunResult
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if len(result.AdapterExecutions) != 2 {
		t.Fatalf("adapter executions = %+v", result.AdapterExecutions)
	}
	executedIDs := map[string]bool{}
	for _, execution := range result.AdapterExecutions {
		executedIDs[execution.AdapterID] = true
		if execution.Resolution.Source != adapter.ResolutionSourceExplicitLocal {
			t.Fatalf("unexpected adapter resolution = %+v", execution.Resolution)
		}
	}
	if !executedIDs[first.ID] || !executedIDs[second.ID] {
		t.Fatalf("executed adapter IDs = %+v", executedIDs)
	}
	foundAssertion := false
	for _, assertion := range result.Results {
		if assertion.AssertionID == "PRC-A-CORE-013" {
			foundAssertion = true
			if assertion.Assessment != "pass" || len(assertion.EvidenceObserved) != 2 {
				t.Fatalf("multi-adapter assertion = %+v", assertion)
			}
		}
	}
	if !foundAssertion {
		t.Fatal("missing multi-adapter assertion result")
	}
}

func TestAdapterDataFlagParsingIsQualifiedAndDuplicateSafe(t *testing.T) {
	parsed, err := parseScanAdapterData(repeatedStringFlag{
		"prc.adapter.grype@0.116/grype-db=/tmp/db=cache",
	})
	if err != nil || parsed["prc.adapter.grype@0.116"]["grype-db"] != "/tmp/db=cache" {
		t.Fatalf("parsed=%+v err=%v", parsed, err)
	}
	for _, values := range []repeatedStringFlag{
		{"missing-assignment"},
		{"name=/tmp/db"},
		{"prc.adapter.grype@0.116/grype-db=/a", "prc.adapter.grype@0.116/grype-db=/b"},
	} {
		if _, err := parseScanAdapterData(values); err == nil {
			t.Fatalf("invalid adapter data flags accepted: %v", values)
		}
	}
}

func TestScanRejectsDuplicateAdapterRequests(t *testing.T) {
	manifestPath := filepath.Join("..", "..", "fixtures", "adapters", "fixture-adapter.yaml")
	var stdout, stderr bytes.Buffer
	code := run([]string{
		"scan", "--target", t.TempDir(), "--catalog-root", filepath.Join("..", ".."),
		"--mode", "verify-local", "--adapter-manifest", manifestPath,
		"--adapter-manifest", manifestPath, "--exit-policy", "never",
	}, &stdout, &stderr)
	if code != exitConfiguration || !strings.Contains(stderr.String(), "cannot repeat") {
		t.Fatalf("exit=%d stderr=%s", code, stderr.String())
	}
}

func TestScanRejectsUnauthorizedAdapterBeforeRuntimeExecution(t *testing.T) {
	repository := filepath.Join("..", "..")
	target := t.TempDir()
	if err := os.WriteFile(filepath.Join(target, "app.py"), []byte("print('ready')\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	marker := filepath.Join(t.TempDir(), "executed")
	runtimeDirectory := t.TempDir()
	runtime := filepath.Join(runtimeDirectory, "docker")
	if err := os.WriteFile(runtime, []byte("#!/bin/sh\ntouch \""+marker+"\"\n"), 0o700); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	code := run([]string{
		"scan", "--target", target, "--catalog-root", repository,
		"--mode", "verify-local",
		"--adapter-manifest", filepath.Join(repository, "fixtures", "adapters", "fixture-adapter.yaml"),
		"--adapter-runtime", runtime, "--format", "json", "--exit-policy", "never",
	}, &stdout, &stderr)
	if code != exitPolicyDenied || !strings.Contains(stderr.String(), "not authorized") {
		t.Fatalf("exit=%d stderr=%s", code, stderr.String())
	}
	if _, err := os.Stat(marker); !os.IsNotExist(err) {
		t.Fatal("unauthorized adapter runtime was executed")
	}
}

func TestScanRequiresExplicitCapabilityModeBeforeAdapterExecution(t *testing.T) {
	repository := filepath.Join("..", "..")
	target := t.TempDir()
	if err := os.WriteFile(filepath.Join(target, "app.py"), []byte("print('ready')\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	marker := filepath.Join(t.TempDir(), "executed")
	runtime := filepath.Join(t.TempDir(), "docker")
	if err := os.WriteFile(runtime, []byte("#!/bin/sh\ntouch \""+marker+"\"\n"), 0o700); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	code := run([]string{
		"scan", "--target", target, "--catalog-root", repository,
		"--adapter-manifest", filepath.Join(repository, "fixtures", "adapters", "fixture-adapter.yaml"),
		"--adapter-runtime", runtime, "--exit-policy", "never",
	}, &stdout, &stderr)
	if code != exitPolicyDenied || !strings.Contains(stderr.String(), "explicit --mode verify-local") {
		t.Fatalf("exit=%d stderr=%s", code, stderr.String())
	}
	if _, err := os.Stat(marker); !os.IsNotExist(err) {
		t.Fatal("adapter runtime executed without an explicit capability mode")
	}
	stdout.Reset()
	stderr.Reset()
	code = run([]string{
		"scan", "--target", target, "--catalog-root", repository,
		"--mode", "inspect",
		"--adapter-manifest", filepath.Join(repository, "fixtures", "adapters", "fixture-adapter.yaml"),
		"--adapter-runtime", runtime, "--exit-policy", "never",
	}, &stdout, &stderr)
	if code != exitPolicyDenied || !strings.Contains(stderr.String(), "explicit --mode verify-local") {
		t.Fatalf("inspect exit=%d stderr=%s", code, stderr.String())
	}
	if _, err := os.Stat(marker); !os.IsNotExist(err) {
		t.Fatal("adapter runtime executed in inspect mode")
	}
}

func TestScanRejectsUnsupportedExecutionModeAsConfiguration(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{
		"scan", "--target", t.TempDir(), "--catalog-root", filepath.Join("..", ".."),
		"--mode", "production-write", "--exit-policy", "never",
	}, &stdout, &stderr)
	if code != exitConfiguration || !strings.Contains(stderr.String(), "unsupported execution mode") {
		t.Fatalf("exit=%d stderr=%s", code, stderr.String())
	}
}

func TestUnknownCommandIsUsageError(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if code := run([]string{"unknown"}, &stdout, &stderr); code != exitConfiguration {
		t.Fatalf("exit=%d stderr=%s", code, stderr.String())
	}
}

func TestRemediateCommandCreatesAcceptedCandidate(t *testing.T) {
	target := t.TempDir()
	if err := os.WriteFile(filepath.Join(target, "app.py"), []byte("print('ready')"), 0o644); err != nil {
		t.Fatal(err)
	}
	candidate := filepath.Join(t.TempDir(), "candidate")
	var stdout, stderr bytes.Buffer
	code := run([]string{
		"remediate", "--target", target, "--catalog-root", filepath.Join("..", ".."),
		"--candidate-dir", candidate, "--format", "json",
	}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit=%d stderr=%s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), `"accepted": true`) {
		t.Fatalf("unexpected output %q", stdout.String())
	}
	original, err := os.ReadFile(filepath.Join(target, "app.py"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.HasSuffix(string(original), "\n") {
		t.Fatal("command modified original target")
	}
}

func TestFixCommandRunsBoundedDeterministicLoop(t *testing.T) {
	target := t.TempDir()
	targetPath := filepath.Join(target, "app.py")
	if err := os.WriteFile(targetPath, []byte("print('ready')"), 0o666); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(targetPath, 0o666); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	code := run([]string{
		"fix", "--target", target, "--catalog-root", filepath.Join("..", ".."),
		"--candidate-root", filepath.Join(t.TempDir(), "candidates"), "--format", "json",
	}, &stdout, &stderr)
	if code != 1 {
		t.Fatalf("exit=%d stderr=%s", code, stderr.String())
	}
	var result remediation.RemediationRun
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if result.SchemaVersion != "prc.remediation-run/v0.9" || result.MaxDurationSeconds != 1800 ||
		len(result.Candidates) != 2 || len(result.Attempts) != 2 ||
		result.ProviderExecutions == nil || len(result.ProviderExecutions) != 0 ||
		result.TerminalState != "machine_work_complete" || !result.OriginalUnchanged {
		t.Fatalf("unexpected remediation run: %+v", result)
	}
	original, err := os.ReadFile(targetPath)
	if err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(targetPath)
	if err != nil {
		t.Fatal(err)
	}
	if strings.HasSuffix(string(original), "\n") || info.Mode().Perm() != 0o666 {
		t.Fatal("fix command modified the original workspace")
	}
}

func TestFixCommandRunsBoundedSuggestOnlyProvider(t *testing.T) {
	target := t.TempDir()
	if err := os.WriteFile(filepath.Join(target, "app.py"), []byte("def ready(): return True\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	providerPath := fakeFixCodex(t)
	candidateRoot := filepath.Join(t.TempDir(), "candidates")
	var stdout, stderr bytes.Buffer
	code := run([]string{
		"fix", "--target", target, "--catalog-root", filepath.Join("..", ".."),
		"--candidate-root", candidateRoot, "--max-attempts", "1", "--format", "json",
		"--provider", "codex", "--provider-executable", providerPath,
		"--verifier-runtime", fakeVerifierRuntime(t), "--verifier-image", testVerifierImage,
		"--allow-remote-source-processing",
	}, &stdout, &stderr)
	if code != exitIncomplete {
		t.Fatalf("exit=%d stderr=%s", code, stderr.String())
	}
	var result remediation.RemediationRun
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	if len(result.ProviderExecutions) != 1 || len(result.Candidates) != 1 ||
		!result.Candidates[0].Accepted || result.Candidates[0].Contract.Provider != "codex" {
		t.Fatalf("unexpected provider remediation: %+v", result)
	}
	if _, err := os.Stat(filepath.Join(target, "app_test.py")); !os.IsNotExist(err) {
		t.Fatal("fix provider modified the original workspace")
	}
	if _, err := os.Stat(filepath.Join(result.ResultWorkspace, "app_test.py")); err != nil {
		t.Fatalf("provider candidate test is missing: %v", err)
	}
}

func TestFixCommandRequiresRemoteSourceAcknowledgement(t *testing.T) {
	target := t.TempDir()
	if err := os.WriteFile(filepath.Join(target, "app.py"), []byte("def ready(): return True\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	candidateRoot := filepath.Join(t.TempDir(), "candidates")
	var stdout, stderr bytes.Buffer
	code := run([]string{
		"fix", "--target", target, "--catalog-root", filepath.Join("..", ".."),
		"--candidate-root", candidateRoot, "--provider", "codex",
		"--provider-executable", fakeFixCodex(t),
		"--verifier-runtime", fakeVerifierRuntime(t), "--verifier-image", testVerifierImage,
	}, &stdout, &stderr)
	if code != exitPolicyDenied || !strings.Contains(stderr.String(), "PRC-EXIT-5") {
		t.Fatalf("exit=%d stderr=%s", code, stderr.String())
	}
	if _, err := os.Stat(candidateRoot); !os.IsNotExist(err) {
		t.Fatal("source-processing denial created a candidate root")
	}
}

func TestFixProviderLaunchFailureUsesExecutionExitCode(t *testing.T) {
	target := t.TempDir()
	if err := os.WriteFile(filepath.Join(target, "app.py"), []byte("def ready(): return True\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	code := run([]string{
		"fix", "--target", target, "--catalog-root", filepath.Join("..", ".."),
		"--candidate-root", filepath.Join(t.TempDir(), "candidates"),
		"--provider", "codex", "--provider-executable", filepath.Join(t.TempDir(), "codex"),
		"--verifier-runtime", fakeVerifierRuntime(t), "--verifier-image", testVerifierImage,
		"--allow-remote-source-processing",
	}, &stdout, &stderr)
	if code != exitExecution || !strings.Contains(stderr.String(), "PRC-EXIT-4") {
		t.Fatalf("exit=%d stderr=%s", code, stderr.String())
	}
	if _, err := os.Stat(filepath.Join(target, "app_test.py")); !os.IsNotExist(err) {
		t.Fatal("provider launch failure modified the original workspace")
	}
}

func fakeFixCodex(t *testing.T) string {
	t.Helper()
	body := `prompt=$(/bin/cat)
task_id=$(printf '%s' "$prompt" | /usr/bin/sed -n 's/.*"task_id": "\([0-9a-f]*\)".*/\1/p' | /usr/bin/head -n 1)
result=""
while [ "$#" -gt 0 ]; do
  if [ "$1" = "--output-last-message" ]; then shift; result="$1"; fi
  shift
done
printf '%s' '{"schema_version":"prc.agent-output/v0.1","task_id":"'"$task_id"'","status":"candidate","root_cause":"No recognized automated test exists.","changed_files":["app_test.py"],"patch":"diff --git a/app_test.py b/app_test.py\nnew file mode 100644\n--- /dev/null\n+++ b/app_test.py\n@@ -0,0 +1,4 @@\n+from app import ready\n+\n+def test_ready():\n+    assert ready() is True\n","commands_requested_or_run":[],"limitations":[],"requested_capability_changes":[]}' > "$result"
printf '%s\n' '{"type":"turn.completed"}'
`
	path := filepath.Join(t.TempDir(), "codex")
	if err := os.WriteFile(path, []byte("#!/bin/sh\nset -eu\n"+body), 0o755); err != nil {
		t.Fatal(err)
	}
	return path
}

const testVerifierImage = "registry.example/prc/test-verifier@sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

func fakeVerifierRuntime(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "docker")
	if err := os.WriteFile(path, []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestDoctorCommandReportsVerifiedEnvironment(t *testing.T) {
	target := t.TempDir()
	state := t.TempDir()
	if err := os.Chmod(state, 0o700); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	code := run([]string{
		"doctor", "--target", target, "--catalog-root", filepath.Join("..", ".."),
		"--state-dir", state, "--candidate-parent", t.TempDir(), "--format", "json",
	}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit=%d stderr=%s", code, stderr.String())
	}
	var report struct {
		SchemaVersion string `json:"schema_version"`
		Ready         bool   `json:"ready"`
		Summary       struct {
			Failed int `json:"failed"`
		} `json:"summary"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &report); err != nil {
		t.Fatal(err)
	}
	if report.SchemaVersion != "prc.doctor/v0.1" || !report.Ready || report.Summary.Failed != 0 {
		t.Fatalf("unexpected doctor report: %+v", report)
	}
}

func TestDoctorCommandFailsForUnsafeCandidateParent(t *testing.T) {
	target := t.TempDir()
	var stdout, stderr bytes.Buffer
	code := run([]string{
		"doctor", "--target", target, "--catalog-root", filepath.Join("..", ".."),
		"--candidate-parent", target, "--format", "json",
	}, &stdout, &stderr)
	if code != exitIncomplete || !strings.Contains(stdout.String(), `"ready": false`) {
		t.Fatalf("exit=%d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
}

func TestStableExitCodeContract(t *testing.T) {
	wantTerminal := map[string]int{
		"profile_satisfied":     exitSuccess,
		"no_go":                 exitGateFailed,
		"assessment_incomplete": exitIncomplete,
		"environment_blocked":   exitIncomplete,
		"machine_work_complete_manual_evidence_remaining": exitIncomplete,
		"policy_stopped":   exitPolicyDenied,
		"budget_exhausted": exitPolicyDenied,
		"unrecognized":     exitInternal,
	}
	for terminal, want := range wantTerminal {
		if got := scanTerminalExitCode(terminal); got != want {
			t.Errorf("terminal %s exit=%d want=%d", terminal, got, want)
		}
	}
	if got := errorExitCode(context.Canceled, exitInternal); got != exitCancelled {
		t.Fatalf("cancelled exit=%d", got)
	}
	if got := remediationExitCode(remediation.RemediationRun{TerminalState: "candidate_rejected"}); got != exitCandidateRejected {
		t.Fatalf("candidate rejection exit=%d", got)
	}
	if got := remediationExitCode(remediation.RemediationRun{TerminalState: "stopped_by_policy_or_budget"}); got != exitPolicyDenied {
		t.Fatalf("policy stop exit=%d", got)
	}
	if got := remediationExitCode(remediation.RemediationRun{TerminalState: "provider_stopped"}); got != exitIncomplete {
		t.Fatalf("provider stop exit=%d", got)
	}
	if got := remediationExitCode(remediation.RemediationRun{TerminalState: "provider_failed"}); got != exitExecution {
		t.Fatalf("provider failure exit=%d", got)
	}
	if got := remediationExitCode(remediation.RemediationRun{
		TerminalState: "provider_failed", Attempts: []remediation.AttemptRecord{{ReasonCode: "provider_cancelled"}},
	}); got != exitCancelled {
		t.Fatalf("provider cancellation exit=%d", got)
	}
	if got := errorExitCode(errors.New("unclassified"), exitInternal); got != exitInternal {
		t.Fatalf("internal fallback exit=%d", got)
	}
}

type failingWriter struct{}

func (failingWriter) Write([]byte) (int, error) { return 0, errors.New("fixture write failed") }

func TestJSONOutputFailureUsesInternalExitCode(t *testing.T) {
	path := filepath.Join("..", "..", "fixtures", "config", "production-readiness.yaml")
	var stderr bytes.Buffer
	code := run([]string{"config", "validate", "--file", path, "--format", "json"}, failingWriter{}, &stderr)
	if code != exitInternal || !strings.Contains(stderr.String(), "PRC-EXIT-6") {
		t.Fatalf("exit=%d stderr=%s", code, stderr.String())
	}
}

func TestRemediationPolicyDenialUsesPolicyExitCode(t *testing.T) {
	target := t.TempDir()
	if err := os.WriteFile(filepath.Join(target, "app.py"), []byte("print('ready')"), 0o644); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	code := run([]string{
		"remediate", "--target", target, "--catalog-root", filepath.Join("..", ".."),
		"--candidate-dir", filepath.Join(t.TempDir(), "candidate"), "--assertion", "PRC-A-CORE-001",
	}, &stdout, &stderr)
	if code != exitPolicyDenied || !strings.Contains(stderr.String(), "PRC-EXIT-5") {
		t.Fatalf("exit=%d stderr=%s", code, stderr.String())
	}
}

func TestScanIndexesDurableHistoryAndHistoryCommandsReadCanonicalRun(t *testing.T) {
	target := t.TempDir()
	if err := os.WriteFile(filepath.Join(target, "app.py"), []byte("print('ready')\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	stateDirectory := t.TempDir()
	if err := os.Chmod(stateDirectory, 0o700); err != nil {
		t.Fatal(err)
	}
	var scanOutput, stderr bytes.Buffer
	code := run([]string{
		"scan", "--target", target, "--catalog-root", filepath.Join("..", ".."),
		"--state-dir", stateDirectory, "--format", "json", "--exit-policy", "never",
	}, &scanOutput, &stderr)
	if code != 0 {
		t.Fatalf("scan exit=%d stderr=%s", code, stderr.String())
	}
	var scanned model.RunResult
	if err := json.Unmarshal(scanOutput.Bytes(), &scanned); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(stateDirectory, "state.sqlite")); err != nil {
		t.Fatalf("state index was not created: %v", err)
	}
	var historyOutput bytes.Buffer
	stderr.Reset()
	code = run([]string{
		"history", "list", "--state-dir", stateDirectory, "--target-name", scanned.Inventory.TargetName,
		"--format", "json",
	}, &historyOutput, &stderr)
	if code != 0 {
		t.Fatalf("history exit=%d stderr=%s", code, stderr.String())
	}
	var history struct {
		SchemaVersion string `json:"schema_version"`
		Runs          []struct {
			RunID string `json:"run_id"`
		} `json:"runs"`
	}
	if err := json.Unmarshal(historyOutput.Bytes(), &history); err != nil {
		t.Fatal(err)
	}
	if history.SchemaVersion != "prc.history/v0.1" || len(history.Runs) != 1 || history.Runs[0].RunID != scanned.RunID {
		t.Fatalf("unexpected history: %+v", history)
	}
	var showOutput bytes.Buffer
	stderr.Reset()
	code = run([]string{
		"history", "show", "--state-dir", stateDirectory, "--format", "json", scanned.RunID,
	}, &showOutput, &stderr)
	if code != 0 {
		t.Fatalf("history show exit=%d stderr=%s", code, stderr.String())
	}
	var shown model.RunResult
	if err := json.Unmarshal(showOutput.Bytes(), &shown); err != nil || shown.RunID != scanned.RunID {
		t.Fatalf("shown run=%s err=%v", shown.RunID, err)
	}
	var checkOutput bytes.Buffer
	stderr.Reset()
	code = run([]string{
		"history", "check", "--state-dir", stateDirectory, "--format", "json",
	}, &checkOutput, &stderr)
	if code != 0 || !strings.Contains(checkOutput.String(), `"integrity": "ok"`) {
		t.Fatalf("history check exit=%d stdout=%s stderr=%s", code, checkOutput.String(), stderr.String())
	}
	if err := os.WriteFile(filepath.Join(target, "app.py"), []byte("print('changed')\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	var diffOutput bytes.Buffer
	stderr.Reset()
	code = run([]string{
		"diff", "--state-dir", stateDirectory, "--base-run", scanned.RunID,
		"--target", target, "--catalog-root", filepath.Join("..", ".."), "--format", "json",
	}, &diffOutput, &stderr)
	if code != 0 {
		t.Fatalf("diff exit=%d stderr=%s", code, stderr.String())
	}
	var diff invalidation.Report
	if err := json.Unmarshal(diffOutput.Bytes(), &diff); err != nil {
		t.Fatal(err)
	}
	if diff.SchemaVersion != invalidation.SchemaVersion || len(diff.ChangedFiles) != 1 || diff.ChangedFiles[0].Path != "app.py" {
		t.Fatalf("unexpected invalidation report: %+v", diff)
	}
	var sourceInvalidated, readmeUnchanged bool
	for _, impact := range diff.Assertions {
		if impact.AssertionID == "PRC-A-CORE-014" && impact.Conclusion == "invalidated" {
			sourceInvalidated = true
		}
		if impact.AssertionID == "PRC-A-CORE-001" && impact.Conclusion == "unchanged_inputs" && !impact.ReuseAllowed {
			readmeUnchanged = true
		}
	}
	if !sourceInvalidated || !readmeUnchanged {
		t.Fatalf("unexpected assertion impact: %+v", diff.Assertions)
	}
}

func TestRemediateCommandConsumesConfiguredPolicy(t *testing.T) {
	target := t.TempDir()
	if err := os.WriteFile(filepath.Join(target, "app.py"), []byte("print('ready')"), 0o644); err != nil {
		t.Fatal(err)
	}
	configData, err := os.ReadFile(filepath.Join("..", "..", "fixtures", "config", "production-readiness.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	configPath := filepath.Join(target, "production-readiness.yaml")
	if err := os.WriteFile(configPath, configData, 0o600); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	code := run([]string{
		"remediate", "--target", target, "--catalog-root", filepath.Join("..", ".."),
		"--config", configPath, "--candidate-dir", filepath.Join(t.TempDir(), "candidate"), "--format", "json",
	}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit=%d stderr=%s", code, stderr.String())
	}
	var candidate struct {
		SchemaVersion string `json:"schema_version"`
		Contract      struct {
			ConfigurationDigest string `json:"configuration_digest"`
			ProjectID           string `json:"project_id"`
			MaxAttempts         int    `json:"max_attempts"`
		} `json:"contract"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &candidate); err != nil {
		t.Fatal(err)
	}
	if candidate.SchemaVersion != "prc.remediation-candidate/v0.4" ||
		len(candidate.Contract.ConfigurationDigest) != 64 || candidate.Contract.ProjectID != "example-product" ||
		candidate.Contract.MaxAttempts != 3 {
		t.Fatalf("configured candidate = %+v", candidate)
	}
}

func TestRemediateCommandRestrictsCandidateModeOnly(t *testing.T) {
	target := t.TempDir()
	targetPath := filepath.Join(target, "app.py")
	if err := os.WriteFile(targetPath, []byte("print('ready')\n"), 0o666); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(targetPath, 0o666); err != nil {
		t.Fatal(err)
	}
	candidate := filepath.Join(t.TempDir(), "candidate")
	var stdout, stderr bytes.Buffer
	code := run([]string{
		"remediate", "--target", target, "--catalog-root", filepath.Join("..", ".."),
		"--assertion", "PRC-A-CORE-022", "--candidate-dir", candidate, "--format", "json",
	}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit=%d stderr=%s", code, stderr.String())
	}
	originalInfo, err := os.Stat(targetPath)
	if err != nil {
		t.Fatal(err)
	}
	candidateInfo, err := os.Stat(filepath.Join(candidate, "app.py"))
	if err != nil {
		t.Fatal(err)
	}
	if originalInfo.Mode().Perm() != 0o666 || candidateInfo.Mode().Perm() != 0o644 {
		t.Fatalf("original=%#o candidate=%#o", originalInfo.Mode().Perm(), candidateInfo.Mode().Perm())
	}
}

func TestRemediateProposalCommandCreatesAcceptedIsolatedCandidate(t *testing.T) {
	target := t.TempDir()
	if err := os.WriteFile(filepath.Join(target, "app.py"), []byte("def ready(): return True\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	finding := commandFinding(t, target, "PRC-A-CORE-010")
	draft := provider.Task{
		SchemaVersion: provider.TaskSchema, Mode: "suggest", WorkspaceInventoryDigest: strings.Repeat("b", 64),
		FindingID: finding.ID, FindingFingerprint: finding.Fingerprint,
		AssertionID: "PRC-A-CORE-010", ControlIDs: []string{"USEQ-12655775"},
		Goal: "Add one focused test.", ReadScope: "task-inputs", RelevantPaths: []string{"app.py"},
		Inputs: []provider.InputFile{}, AllowedPaths: []string{"app_test.py"},
		ProtectedPaths:  []string{".git/", ".github/workflows/", ".prc/", "catalog/", "production-readiness.yaml", "schemas/"},
		AllowedCommands: [][]string{}, Network: "deny", Secrets: "none",
		AllowRemoteSourceProcessing: true, TimeoutSeconds: 60, MaxOutputBytes: 256 * 1024,
	}
	draftData, _ := json.Marshal(draft)
	draftPath := filepath.Join(t.TempDir(), "draft.json")
	if err := os.WriteFile(draftPath, draftData, 0o600); err != nil {
		t.Fatal(err)
	}
	task, err := provider.SealTask(draftPath, target)
	if err != nil {
		t.Fatal(err)
	}
	taskData, _ := json.Marshal(task)
	taskPath := filepath.Join(t.TempDir(), "task.json")
	if err := os.WriteFile(taskPath, taskData, 0o600); err != nil {
		t.Fatal(err)
	}
	proposal := provider.Output{
		SchemaVersion: provider.OutputSchema, TaskID: task.TaskID, Status: "candidate",
		RootCause: "No test is present.", ChangedFiles: []string{"app_test.py"},
		Patch:                  "diff --git a/app_test.py b/app_test.py\n--- /dev/null\n+++ b/app_test.py\n@@ -0,0 +1,4 @@\n+from app import ready\n+\n+def test_ready():\n+    assert ready() is True\n",
		CommandsRequestedOrRun: []provider.CommandResult{}, Limitations: []string{}, RequestedCapabilityChanges: []string{},
	}
	proposalData, _ := json.Marshal(proposal)
	proposalPath := filepath.Join(t.TempDir(), "proposal.json")
	if err := os.WriteFile(proposalPath, proposalData, 0o600); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	candidatePath := filepath.Join(t.TempDir(), "candidate")
	code := run([]string{
		"remediate-proposal", "--target", target, "--catalog-root", filepath.Join("..", ".."),
		"--provider", "codex", "--task", taskPath, "--output", proposalPath,
		"--verifier-runtime", fakeVerifierRuntime(t), "--verifier-image", testVerifierImage,
		"--candidate-dir", candidatePath, "--format", "json",
	}, &stdout, &stderr)
	if code != 0 || !strings.Contains(stdout.String(), `"accepted": true`) {
		t.Fatalf("exit=%d stderr=%s stdout=%s", code, stderr.String(), stdout.String())
	}
	if _, err := os.Stat(filepath.Join(target, "app_test.py")); !os.IsNotExist(err) {
		t.Fatal("source workspace was modified")
	}
}
