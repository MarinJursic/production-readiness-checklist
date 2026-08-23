package adapter

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"

	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
)

func validSyftManifest() Manifest {
	limits := DefaultLimits()
	limits.MaxStdin = 1024
	limits.MaxStdout = 64 * 1024 * 1024
	return Manifest{
		SchemaVersion: ManifestSchema, ID: "prc.adapter.syft@1.51",
		Title: "Syft repository CycloneDX SBOM", Description: "Produces a normalized CycloneDX SBOM for engine tests.",
		Publisher: Publisher{ID: "prc-project", Name: "Production Readiness Checklist"},
		Owner:     "scanner-maintainers", Maintenance: "active",
		Protocol: SyftProtocolVersion, OutputSchema: SyftOutputSchemaVersion,
		ObservationKinds: []string{SyftObservationKind},
		Compatibility:    Compatibility{EngineAPIs: []string{model.EngineVersion}},
		Tool: Tool{
			Name: "syft", Version: SyftToolVersion, Upstream: "https://github.com/anchore/syft",
			Formats: []ToolFormat{{Name: "cyclonedx-json", Versions: []string{"1.7"}}},
		},
		Limitations: []string{"Test fixture exercises only the reviewed repository protocol."},
		Runner:      "oci", Image: SyftImage, Command: syftCommand(),
		Capabilities: Capabilities{
			ReadWorkspace: true, WriteScratch: false, Network: "deny",
			NetworkHosts: []string{}, SecretHandles: []string{}, ChildProcesses: true,
		},
		Resources: Resources{
			TimeoutSeconds: 180, MemoryMB: 512, CPUs: 1, PIDs: 64, TmpfsMB: 1, Limits: limits,
		},
	}
}

func validSyftDocument(serial, timestamp string) map[string]any {
	return map[string]any{
		"$schema": "http://cyclonedx.org/schema/bom-1.7.schema.json", "bomFormat": "CycloneDX",
		"specVersion": "1.7", "serialNumber": serial, "version": 1,
		"metadata": map[string]any{
			"timestamp": timestamp,
			"tools": map[string]any{"components": []any{map[string]any{
				"type": "application", "author": "anchore", "name": "syft", "version": "1.51.0",
			}}},
			"component": map[string]any{"bom-ref": "source-123", "type": "file", "name": "/workspace"},
		},
		"components": []any{
			map[string]any{
				"bom-ref": "pkg:golang/example.com/module@v1.2.3?package-id=123", "type": "library",
				"name": "example.com/module", "version": "v1.2.3", "purl": "pkg:golang/example.com/module@v1.2.3",
				"properties": []any{
					map[string]any{"name": "syft:package:type", "value": "go-module"},
					map[string]any{"name": "syft:location:0:path", "value": "/go.mod"},
				},
			},
			map[string]any{"bom-ref": "file-123", "type": "file", "name": "/workspace/go.mod"},
		},
		"dependencies": []any{
			map[string]any{"ref": "source-123", "dependsOn": []any{"pkg:golang/example.com/module@v1.2.3?package-id=123"}},
		},
	}
}

func TestSyftManifestPinsReviewedExecutionContract(t *testing.T) {
	manifest := validSyftManifest()
	if err := manifest.ValidateForCurrentEngine(); err != nil {
		t.Fatal(err)
	}
	tests := map[string]func(*Manifest){
		"image":   func(value *Manifest) { value.Image = "ghcr.io/anchore/syft@sha256:" + strings.Repeat("f", 64) },
		"command": func(value *Manifest) { value.Command = append(value.Command, "--verbose") },
		"tool":    func(value *Manifest) { value.Tool.Version = "1.50.0" },
		"kind":    func(value *Manifest) { value.ObservationKinds = []string{"package"} },
		"scratch": func(value *Manifest) { value.Capabilities.WriteScratch = true },
		"tasks":   func(value *Manifest) { value.Capabilities.ChildProcesses = false },
	}
	for name, mutate := range tests {
		t.Run(name, func(t *testing.T) {
			value := validSyftManifest()
			mutate(&value)
			if err := value.Validate(); err == nil {
				t.Fatal("tampered Syft manifest was accepted")
			}
		})
	}
}

func TestCheckedInSyftManifestHasCatalogPinnedDigest(t *testing.T) {
	manifest, err := LoadManifest(filepath.Join("..", "..", "adapters", "syft-v1.51.0.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	digest, err := ManifestDigest(manifest)
	if err != nil {
		t.Fatal(err)
	}
	if digest != "f2c75e95c95a2151fb588615d03842a1ca2f5ad2f4e6fd6b98e08cbfcfba6836" {
		t.Fatalf("Syft manifest digest = %s", digest)
	}
}

func TestSyftCycloneDXNormalizesNondeterminismAndPreservesArtifact(t *testing.T) {
	firstDocument := validSyftDocument("urn:uuid:11111111-1111-4111-8111-111111111111", "2026-08-23T12:00:00Z")
	secondDocument := validSyftDocument("urn:uuid:22222222-2222-4222-8222-222222222222", "2026-08-23T12:01:00Z")
	firstData, _ := json.Marshal(firstDocument)
	secondData, _ := json.Marshal(secondDocument)
	first, firstPayloads, err := parseSyftOutput(firstData, 100)
	if err != nil {
		t.Fatal(err)
	}
	second, secondPayloads, err := parseSyftOutput(secondData, 100)
	if err != nil {
		t.Fatal(err)
	}
	if first.Observations[0].Outcome != "value" || first.Observations[0].Kind != SyftObservationKind ||
		first.Observations[0].Data["component_count"] != 2 || first.Observations[0].Data["package_component_count"] != 1 ||
		len(first.Artifacts) != 1 || len(firstPayloads) != 1 {
		t.Fatalf("normalized Syft transcript = %+v", first)
	}
	if first.Artifacts[0].Digest != second.Artifacts[0].Digest ||
		!bytes.Equal(firstPayloads[first.Artifacts[0].Digest], secondPayloads[second.Artifacts[0].Digest]) {
		t.Fatal("Syft UUID or timestamp changed the normalized artifact")
	}
	payload := firstPayloads[first.Artifacts[0].Digest]
	if bytes.Contains(payload, []byte("serialNumber")) || bytes.Contains(payload, []byte("timestamp")) ||
		bytes.Contains(payload, []byte("source-123\",\"type\":\"file\",\"name\":\"/workspace")) {
		t.Fatalf("normalized payload retained nondeterministic metadata: %s", payload)
	}
	var normalized map[string]any
	if err := json.Unmarshal(payload, &normalized); err != nil || normalized["bomFormat"] != "CycloneDX" {
		t.Fatalf("normalized CycloneDX artifact = %+v, %v", normalized, err)
	}
}

func TestSyftCycloneDXRejectsUnreviewedOrMalformedOutput(t *testing.T) {
	tests := map[string]func(map[string]any){
		"unknown top-level property": func(value map[string]any) { value["assessment"] = "pass" },
		"wrong tool": func(value map[string]any) {
			value["metadata"].(map[string]any)["tools"].(map[string]any)["components"].([]any)[0].(map[string]any)["name"] = "other"
		},
		"invalid purl": func(value map[string]any) {
			value["components"].([]any)[0].(map[string]any)["purl"] = "pkg:bad value"
		},
		"duplicate component": func(value map[string]any) {
			items := value["components"].([]any)
			value["components"] = append(items, items[0])
		},
	}
	for name, mutate := range tests {
		t.Run(name, func(t *testing.T) {
			document := validSyftDocument("urn:uuid:11111111-1111-4111-8111-111111111111", "2026-08-23T12:00:00Z")
			mutate(document)
			data, _ := json.Marshal(document)
			if _, _, err := parseSyftOutput(data, 100); err == nil {
				t.Fatal("invalid Syft output was accepted")
			}
		})
	}
	duplicateKeys := []byte(`{"$schema":"x","$schema":"y"}`)
	if _, _, err := parseSyftOutput(duplicateKeys, 100); err == nil {
		t.Fatal("duplicate JSON keys were accepted")
	}
}

func TestSyftExecutionInputIsEmptyAndIdentityBound(t *testing.T) {
	input, err := ExecutionInput(validSyftManifest(), strings.Repeat("a", 64), Subject{
		TargetName: "fixture", InventoryDigest: strings.Repeat("b", 64),
	}, map[string]any{"untrusted": true}, map[string]any{"config": "ignored"})
	if err != nil || len(input) != 0 {
		t.Fatalf("Syft execution input = %q, %v", input, err)
	}
	if _, err := ExecutionInput(validSyftManifest(), "invalid", Subject{}, nil, nil); err == nil {
		t.Fatal("invalid Syft execution identity was accepted")
	}
}
