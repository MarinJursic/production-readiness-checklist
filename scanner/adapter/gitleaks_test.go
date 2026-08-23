package adapter

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"

	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
)

func validGitleaksManifest() Manifest {
	limits := DefaultLimits()
	limits.MaxStdout = 16 * 1024 * 1024
	return Manifest{
		SchemaVersion: ManifestSchema, ID: "prc.adapter.gitleaks@8.30",
		Title:       "Gitleaks current-tree secret detection",
		Description: "Normalizes a redacted Gitleaks report into bounded secret-detection observations.",
		Publisher:   Publisher{ID: "prc-project", Name: "Production Readiness Checklist"},
		Owner:       "scanner-maintainers", Maintenance: "active",
		Protocol: GitleaksProtocolVersion, OutputSchema: GitleaksOutputSchemaVersion,
		ObservationKinds: []string{GitleaksObservationKind},
		Compatibility:    Compatibility{EngineAPIs: []string{model.EngineVersion}},
		Tool: Tool{
			Name: "gitleaks", Version: GitleaksToolVersion, Upstream: "https://github.com/gitleaks/gitleaks",
			Formats: []ToolFormat{{Name: "gitleaks-json", Versions: []string{GitleaksToolVersion}}},
		},
		Limitations: []string{"Scans the current inventoried file tree, not Git history."},
		Runner:      "oci", Image: GitleaksImage, Command: gitleaksCommand(),
		Capabilities: Capabilities{
			ReadWorkspace: true, WriteScratch: false, Network: "deny",
			NetworkHosts: []string{}, SecretHandles: []string{}, ChildProcesses: true,
		},
		Resources: Resources{
			TimeoutSeconds: 120, MemoryMB: 256, CPUs: 1, PIDs: 64, TmpfsMB: 1, Limits: limits,
		},
	}
}

func validGitleaksFinding() gitleaksFinding {
	return gitleaksFinding{
		RuleID: "github-pat", Description: "Potential GitHub personal access token.",
		StartLine: 7, EndLine: 7, StartColumn: 10, EndColumn: 49,
		Match: "REDACTED", Secret: "REDACTED", File: "/workspace/config/example.env",
		Entropy: 4.8, Tags: []string{}, Fingerprint: "/workspace/config/example.env:github-pat:7",
	}
}

func TestGitleaksManifestPinsReviewedExecutionContract(t *testing.T) {
	manifest := validGitleaksManifest()
	if err := manifest.ValidateForCurrentEngine(); err != nil {
		t.Fatal(err)
	}
	tests := map[string]func(*Manifest){
		"image":   func(value *Manifest) { value.Image = "ghcr.io/gitleaks/gitleaks@sha256:" + strings.Repeat("f", 64) },
		"command": func(value *Manifest) { value.Command = append(value.Command, "--verbose") },
		"tool":    func(value *Manifest) { value.Tool.Version = "8.30.1" },
		"kind":    func(value *Manifest) { value.ObservationKinds = []string{"static-analysis"} },
		"scratch": func(value *Manifest) { value.Capabilities.WriteScratch = true },
		"tasks":   func(value *Manifest) { value.Capabilities.ChildProcesses = false },
	}
	for name, mutate := range tests {
		t.Run(name, func(t *testing.T) {
			value := validGitleaksManifest()
			mutate(&value)
			if err := value.Validate(); err == nil {
				t.Fatal("tampered Gitleaks manifest was accepted")
			}
		})
	}
}

func TestCheckedInGitleaksManifestHasCatalogPinnedDigest(t *testing.T) {
	manifest, err := LoadManifest(filepath.Join("..", "..", "adapters", "gitleaks-v8.30.0.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	digest, err := ManifestDigest(manifest)
	if err != nil {
		t.Fatal(err)
	}
	if digest != "ee223fd11f9d32b1799bf2b25c50cea294fd1f7b59686f9d438b351b5c4faac8" {
		t.Fatalf("Gitleaks manifest digest = %s", digest)
	}
}

func TestGitleaksExecutionInputIsPinnedRuleset(t *testing.T) {
	input, err := ExecutionInput(validGitleaksManifest(), strings.Repeat("a", 64), Subject{
		TargetName: "fixture", InventoryDigest: strings.Repeat("b", 64),
	}, map[string]any{"untrusted": true}, map[string]any{"config": "ignored"})
	if err != nil {
		t.Fatal(err)
	}
	expected, err := gitleaksRuleset()
	if err != nil {
		t.Fatal(err)
	}
	digest := sha256.Sum256(input)
	if hex.EncodeToString(digest[:]) != GitleaksConfigSHA256 || !bytes.Equal(input, expected) {
		t.Fatal("Gitleaks execution input is not the pinned embedded ruleset")
	}
	input[0] ^= 0xff
	if bytes.Equal(input, expected) {
		t.Fatal("caller can mutate embedded Gitleaks configuration")
	}
}

func TestGitleaksReportNormalizesAndBindsRedactedRawArtifact(t *testing.T) {
	finding := validGitleaksFinding()
	report, err := json.Marshal([]gitleaksFinding{finding})
	if err != nil {
		t.Fatal(err)
	}
	transcript, err := parseGitleaksOutput(report, 100)
	if err != nil {
		t.Fatal(err)
	}
	if len(transcript.Observations) != 1 || transcript.Observations[0].Outcome != "found" ||
		transcript.Observations[0].Kind != GitleaksObservationKind ||
		transcript.Observations[0].Locations[0].Path != "config/example.env" {
		t.Fatalf("normalized Gitleaks transcript = %+v", transcript)
	}
	encoded, err := json.Marshal(transcript)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Contains(encoded, []byte("REDACTED")) || bytes.Contains(encoded, []byte(finding.Description)) {
		t.Fatalf("normalized transcript retained report content: %s", encoded)
	}
	rawDigest := sha256.Sum256(report)
	if len(transcript.Artifacts) != 1 || transcript.Artifacts[0].Digest != "sha256:"+hex.EncodeToString(rawDigest[:]) ||
		transcript.Artifacts[0].Size != int64(len(report)) {
		t.Fatalf("raw report descriptor = %+v", transcript.Artifacts)
	}
}

func TestEmptyGitleaksReportProducesExplicitNotFoundObservation(t *testing.T) {
	transcript, err := ParseManifestOutput(validGitleaksManifest(), strings.NewReader("[]\n"))
	if err != nil {
		t.Fatal(err)
	}
	if len(transcript.Observations) != 1 || transcript.Observations[0].Outcome != "not_found" ||
		len(transcript.Observations[0].Locations) != 0 || transcript.Summary.Status != "completed" {
		t.Fatalf("empty Gitleaks transcript = %+v", transcript)
	}
}

func TestGitleaksIgnoreSnapshotLocationMapsBackToSubjectPath(t *testing.T) {
	path, err := normalizeGitleaksPath("/workspace/" + gitleaksIgnoreSnapshotPath)
	if err != nil {
		t.Fatal(err)
	}
	if path != gitleaksIgnoreSourcePath {
		t.Fatalf("normalized ignore path = %q", path)
	}
}

func TestGitleaksReportFailsClosedOnUnsafeOrDriftedOutput(t *testing.T) {
	tests := map[string]func(*gitleaksFinding){
		"unredacted": func(value *gitleaksFinding) { value.Secret = "not-redacted" },
		"escape":     func(value *gitleaksFinding) { value.File = "/etc/passwd" },
		"symlink":    func(value *gitleaksFinding) { value.SymlinkFile = "/workspace/link" },
		"history":    func(value *gitleaksFinding) { value.Commit = strings.Repeat("a", 40) },
		"fingerprint": func(value *gitleaksFinding) {
			value.Fingerprint = "mismatched"
		},
	}
	for name, mutate := range tests {
		t.Run(name, func(t *testing.T) {
			finding := validGitleaksFinding()
			mutate(&finding)
			report, err := json.Marshal([]gitleaksFinding{finding})
			if err != nil {
				t.Fatal(err)
			}
			if _, err := parseGitleaksOutput(report, 100); err == nil {
				t.Fatal("unsafe Gitleaks report was accepted")
			}
		})
	}
	for name, report := range map[string]string{
		"null":      "null",
		"unknown":   `[{"RuleID":"x","Unexpected":true}]`,
		"duplicate": `[{"RuleID":"x","RuleID":"y"}]`,
		"trailing":  "[] {}",
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := parseGitleaksOutput([]byte(report), 100); err == nil {
				t.Fatal("drifted Gitleaks report was accepted")
			}
		})
	}
}
