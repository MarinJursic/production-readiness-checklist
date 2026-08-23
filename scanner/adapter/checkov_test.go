package adapter

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/MarinJursic/production-readiness-checklist/scanner/inventory"
)

func validCheckovRecord(path, result string) checkovRecord {
	null := json.RawMessage(`null`)
	return checkovRecord{
		CheckID: "CKV_AWS_54", BCCheckID: null,
		CheckName: "Ensure S3 bucket has block public policy enabled",
		CheckResult: checkovCheckResult{
			Result: result, EvaluatedKeys: []json.RawMessage{json.RawMessage(`"block_public_policy"`)},
			Entity: json.RawMessage(`null`),
		},
		CodeBlock: null, FilePath: "/" + path, FileAbsPath: "/workspace/" + path,
		RepoFilePath: "/workspace/" + path, FileLineRange: []int{10, 20},
		Resource: "aws_s3_bucket_public_access_block.public", Evaluations: json.RawMessage(`null`),
		CheckClass:      "checkov.terraform.checks.resource.aws.S3BlockPublicPolicy",
		FixedDefinition: null, EntityTags: null, CallerFilePath: null, CallerFileLineRange: null,
		ResourceAddress: null, Severity: null, BCCategory: null, Benchmarks: json.RawMessage(`{}`),
		Description: null, ShortDescription: null, VulnerabilityDetails: null, ConnectedNode: null,
		Guideline: null, Details: []json.RawMessage{}, CheckLen: null,
		DefinitionContextFilePath: json.RawMessage(`"/workspace/` + path + `"`),
	}
}

func validCheckovReport(path string, failed bool) checkovReport {
	result := checkovResults{
		PassedChecks: []checkovRecord{}, FailedChecks: []checkovRecord{},
		SkippedChecks: []checkovRecord{}, ParsingErrors: []string{},
	}
	if failed {
		result.FailedChecks = append(result.FailedChecks, validCheckovRecord(path, "FAILED"))
	} else {
		result.PassedChecks = append(result.PassedChecks, validCheckovRecord(path, "PASSED"))
	}
	return checkovReport{
		CheckType: "terraform", Results: result,
		Summary: checkovSummary{
			Passed: len(result.PassedChecks), Failed: len(result.FailedChecks),
			ResourceCount: 1, CheckovVersion: CheckovToolVersion,
		},
		URL: checkovReportURL,
	}
}

func marshalCheckov(t *testing.T, value any) []byte {
	t.Helper()
	payload, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	return payload
}

func TestCheckedInCheckovManifestLoads(t *testing.T) {
	manifest, err := LoadManifest(filepath.Join("..", "..", "adapters", "checkov-v3.3.8.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	if manifest.ID != "prc.adapter.checkov@3.3" || manifest.Image != CheckovImage {
		t.Fatalf("unexpected Checkov manifest: %+v", manifest)
	}
	digest, err := ManifestDigest(manifest)
	if err != nil {
		t.Fatal(err)
	}
	if digest != "643069c38a3507b5ac36093f16f8e1932ca6c7fa3cd972a8576c8913ae821be7" {
		t.Fatalf("Checkov manifest digest = %s", digest)
	}
}

func TestCheckovManifestPinsReviewedExecutionContract(t *testing.T) {
	path := filepath.Join("..", "..", "adapters", "checkov-v3.3.8.yaml")
	tests := map[string]func(*Manifest){
		"image":   func(value *Manifest) { value.Image = "docker.io/bridgecrew/checkov@sha256:" + strings.Repeat("f", 64) },
		"command": func(value *Manifest) { value.Command = append(value.Command, "--directory=/workspace") },
		"tool":    func(value *Manifest) { value.Tool.Version = "3.3.7" },
		"kind":    func(value *Manifest) { value.ObservationKinds = []string{"configuration"} },
		"network": func(value *Manifest) { value.Capabilities.Network = "allow" },
		"scratch": func(value *Manifest) { value.Capabilities.WriteScratch = false },
		"tasks":   func(value *Manifest) { value.Capabilities.ChildProcesses = false },
	}
	for name, mutate := range tests {
		t.Run(name, func(t *testing.T) {
			manifest, err := LoadManifest(path)
			if err != nil {
				t.Fatal(err)
			}
			mutate(&manifest)
			if err := manifest.Validate(); err == nil {
				t.Fatal("tampered Checkov manifest was accepted")
			}
		})
	}
}

func TestCheckovSnapshotInjectsExactScannerOwnedFileList(t *testing.T) {
	root := t.TempDir()
	writeSnapshotFixture(t, root, "infra/main.tf", "resource \"aws_s3_bucket\" \"logs\" {}\n")
	writeSnapshotFixture(t, root, "deploy/app.yaml", "apiVersion: v1\nkind: Pod\nmetadata: {name: app}\n")
	writeSnapshotFixture(t, root, "Dockerfile", "FROM scratch\n")
	writeSnapshotFixture(t, root, ".checkov.yml", "framework: [secrets]\nskip-check: [CKV_AWS_54]\n")
	item, err := inventory.Build(root)
	if err != nil {
		t.Fatal(err)
	}
	manifest, err := LoadManifest(filepath.Join("..", "..", "adapters", "checkov-v3.3.8.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	snapshot, err := PrepareSnapshotForManifest(item, manifest)
	if err != nil {
		t.Fatal(err)
	}
	defer snapshot.Close()
	expected, err := readCheckovExpectedPaths(snapshot.Path)
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(expected, []string{"Dockerfile", "deploy/app.yaml", "infra/main.tf"}) {
		t.Fatalf("scanner-owned Checkov files = %v", expected)
	}
	payload, err := os.ReadFile(filepath.Join(snapshot.Path, filepath.FromSlash(CheckovConfigSnapshotPath)))
	if err != nil || bytes.Contains(payload, []byte("skip-check")) || !bytes.Contains(payload, []byte(`"skip-download":true`)) {
		t.Fatalf("unexpected scanner-owned Checkov configuration: %s (%v)", payload, err)
	}
	if _, err := os.Stat(filepath.Join(snapshot.Path, ".checkov.yml")); err != nil {
		t.Fatalf("target Checkov configuration was not retained as inert inventory data: %v", err)
	}
}

func TestCheckovOCIPlanUsesIsolatedWorkingDirectory(t *testing.T) {
	root := t.TempDir()
	writeSnapshotFixture(t, root, "infra/main.tf", "resource \"aws_s3_bucket\" \"logs\" {}\n")
	item, err := inventory.Build(root)
	if err != nil {
		t.Fatal(err)
	}
	manifest, err := LoadManifest(filepath.Join("..", "..", "adapters", "checkov-v3.3.8.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	snapshot, err := PrepareSnapshotForManifest(item, manifest)
	if err != nil {
		t.Fatal(err)
	}
	defer snapshot.Close()
	runtimePath := filepath.Join(t.TempDir(), "docker")
	if err := os.WriteFile(runtimePath, []byte("#!/bin/sh\nexit 0\n"), 0o700); err != nil {
		t.Fatal(err)
	}
	plan, err := BuildSnapshotOCIPlan(runtimePath, snapshot, strings.Repeat("c", 64), manifest)
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Contains(plan.Arguments, "--workdir=/tmp") || slices.Contains(plan.Arguments, "--workdir=/workspace") {
		t.Fatalf("Checkov OCI working directory is not isolated: %v", plan.Arguments)
	}
}

func TestCheckovSnapshotRejectsTargetsWithoutApplicableFiles(t *testing.T) {
	root := t.TempDir()
	writeSnapshotFixture(t, root, "README.md", "documentation only\n")
	item, err := inventory.Build(root)
	if err != nil {
		t.Fatal(err)
	}
	manifest, err := LoadManifest(filepath.Join("..", "..", "adapters", "checkov-v3.3.8.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := PrepareSnapshotForManifest(item, manifest); err == nil {
		t.Fatal("Checkov snapshot accepted an inventory without applicable files")
	}
}

func TestCheckovParserNormalizesFindingsDeterministically(t *testing.T) {
	report := validCheckovReport("infra/main.tf", true)
	report.Results.FailedChecks = append(report.Results.FailedChecks, report.Results.FailedChecks[0])
	report.Summary.Failed++
	first, firstPayloads, err := parseCheckovOutput(marshalCheckov(t, report), 100, []string{"infra/main.tf"})
	if err != nil {
		t.Fatal(err)
	}
	second, secondPayloads, err := parseCheckovOutput(marshalCheckov(t, report), 100, []string{"infra/main.tf"})
	if err != nil {
		t.Fatal(err)
	}
	if len(first.Observations) != 1 || first.Observations[0].Outcome != "found" ||
		first.Observations[0].Locations[0] != (Location{Path: "infra/main.tf", Line: 10}) ||
		first.Observations[0].Data["check_id"] != "CKV_AWS_54" {
		t.Fatalf("unexpected Checkov findings: %+v", first.Observations)
	}
	if first.Artifacts[0].Digest != second.Artifacts[0].Digest ||
		!bytes.Equal(firstPayloads[first.Artifacts[0].Digest], secondPayloads[second.Artifacts[0].Digest]) {
		t.Fatal("normalized Checkov output is not deterministic")
	}
	var normalized checkovNormalizedReport
	if err := json.Unmarshal(firstPayloads[first.Artifacts[0].Digest], &normalized); err != nil ||
		len(normalized.Findings) != 1 || normalized.Summary.Failed != 2 {
		t.Fatalf("normalized Checkov report = %+v (%v)", normalized, err)
	}
}

func TestCheckovParserAcceptsButNeverPersistsBoundedGraphEntity(t *testing.T) {
	report := validCheckovReport("infra/main.tf", true)
	report.Results.FailedChecks[0].CheckResult.Entity = json.RawMessage(`{"password":"target-secret","nested":{"enabled":true}}`)
	transcript, payloads, err := parseCheckovOutput(marshalCheckov(t, report), 100, []string{"infra/main.tf"})
	if err != nil {
		t.Fatal(err)
	}
	payload := payloads[transcript.Artifacts[0].Digest]
	if bytes.Contains(payload, []byte("target-secret")) || bytes.Contains(payload, []byte("password")) {
		t.Fatalf("normalized Checkov output retained graph entity data: %s", payload)
	}

	for name, entity := range map[string]json.RawMessage{
		"non-object": json.RawMessage(`"target-secret"`),
		"oversized":  json.RawMessage(`{"value":"` + strings.Repeat("x", 1024*1024) + `"}`),
	} {
		t.Run(name, func(t *testing.T) {
			unsafeReport := validCheckovReport("infra/main.tf", true)
			unsafeReport.Results.FailedChecks[0].CheckResult.Entity = entity
			if _, _, err := parseCheckovOutput(marshalCheckov(t, unsafeReport), 100, []string{"infra/main.tf"}); err == nil {
				t.Fatal("unsafe Checkov graph entity was accepted")
			}
		})
	}
}

func TestCheckovParserDistinguishesCleanAndUnsupportedFiles(t *testing.T) {
	clean, _, err := parseCheckovOutput(marshalCheckov(t, validCheckovReport("infra/main.tf", false)), 100, []string{"infra/main.tf"})
	if err != nil {
		t.Fatal(err)
	}
	if len(clean.Observations) != 1 || clean.Observations[0].Outcome != "not_found" {
		t.Fatalf("clean Checkov transcript = %+v", clean)
	}
	partial, payloads, err := parseCheckovOutput(
		marshalCheckov(t, validCheckovReport("infra/main.tf", false)), 100,
		[]string{"Dockerfile", "infra/main.tf"},
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(partial.Observations) != 1 || partial.Observations[0].Outcome != "unsupported" ||
		partial.Observations[0].Locations[0].Path != "Dockerfile" {
		t.Fatalf("partial Checkov transcript = %+v", partial)
	}
	var normalized checkovNormalizedReport
	if err := json.Unmarshal(payloads[partial.Artifacts[0].Digest], &normalized); err != nil ||
		!slices.Equal(normalized.UnsupportedPaths, []string{"Dockerfile"}) {
		t.Fatalf("partial normalized report = %+v (%v)", normalized, err)
	}
}

func TestCheckovParserFailsClosedOnSuppressionsParsingErrorsAndUnsafeOutput(t *testing.T) {
	tests := map[string]func(*checkovReport){
		"suppressed policy": func(value *checkovReport) {
			skipped := validCheckovRecord("infra/main.tf", "SKIPPED")
			skipped.CheckResult.SuppressComment = "target-owned suppression"
			value.Results.SkippedChecks = append(value.Results.SkippedChecks, skipped)
			value.Summary.Skipped++
		},
		"parsing error": func(value *checkovReport) {
			value.Results.ParsingErrors = append(value.Results.ParsingErrors, "/workspace/infra/main.tf")
			value.Summary.ParsingErrors++
		},
		"summary mismatch": func(value *checkovReport) { value.Summary.Passed++ },
		"wrong version":    func(value *checkovReport) { value.Summary.CheckovVersion = "3.3.7" },
		"escaping path":    func(value *checkovReport) { value.Results.FailedChecks[0].FileAbsPath = "/workspace/../etc/passwd" },
		"code disclosure": func(value *checkovReport) {
			value.Results.FailedChecks[0].CodeBlock = json.RawMessage(`[[1,"secret"]]`)
		},
		"online metadata": func(value *checkovReport) { value.Results.FailedChecks[0].Severity = json.RawMessage(`"HIGH"`) },
		"unbounded records": func(value *checkovReport) {
			value.Results.FailedChecks = append(value.Results.FailedChecks, value.Results.FailedChecks[0])
			value.Summary.Failed++
		},
	}
	for name, mutate := range tests {
		t.Run(name, func(t *testing.T) {
			report := validCheckovReport("infra/main.tf", true)
			mutate(&report)
			limit := 100
			if name == "unbounded records" {
				limit = 1
			}
			if _, _, err := parseCheckovOutput(marshalCheckov(t, report), limit, []string{"infra/main.tf"}); err == nil {
				t.Fatal("unsafe Checkov output was accepted")
			}
		})
	}
	duplicate := []byte(`{"check_type":"terraform","check_type":"dockerfile"}`)
	if _, _, err := parseCheckovOutput(duplicate, 100, []string{"infra/main.tf"}); err == nil || !strings.Contains(err.Error(), "duplicate") {
		t.Fatalf("duplicate Checkov keys were not rejected: %v", err)
	}
}

func TestCheckovPublicParserRequiresSealedExecutionContext(t *testing.T) {
	manifest, err := LoadManifest(filepath.Join("..", "..", "adapters", "checkov-v3.3.8.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ParseManifestOutput(manifest, bytes.NewReader(marshalCheckov(t, validCheckovReport("infra/main.tf", false)))); err == nil ||
		!strings.Contains(err.Error(), "sealed execution snapshot") {
		t.Fatalf("context-free Checkov parser error = %v", err)
	}
}
