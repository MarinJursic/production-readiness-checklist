package adapter

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
)

func manifestWithDataMount() Manifest {
	manifest := validManifest()
	manifest.ID = "prc.adapter.grype@0.116"
	manifest.Title = "Grype fixture"
	manifest.Description = "Reviewed Grype data mount fixture."
	manifest.Protocol = GrypeProtocolVersion
	manifest.OutputSchema = GrypeOutputSchemaVersion
	manifest.ObservationKinds = []string{GrypeObservationKind}
	manifest.Tool = Tool{
		Name: "grype", Version: GrypeToolVersion, Upstream: "https://github.com/anchore/grype",
		Formats: []ToolFormat{{Name: "grype-json", Versions: []string{GrypeToolVersion}}},
	}
	manifest.Image = GrypeImage
	manifest.Command = grypeCommand()
	manifest.Capabilities.WriteScratch = true
	manifest.Capabilities.ChildProcesses = true
	manifest.Resources.TimeoutSeconds = 300
	manifest.Resources.MemoryMB = 2048
	manifest.Resources.PIDs = 128
	manifest.Resources.TmpfsMB = 256
	manifest.Resources.MaxStdout = 128 * 1024 * 1024
	manifest.DataMounts = []DataMount{{
		Name: GrypeDataMountName, Destination: grypeDataMountPath,
		MaxFiles: 16, MaxBytes: 3 * 1024 * 1024 * 1024,
	}}
	return manifest
}

func TestOCIDataMountIsBoundReadOnlyAndRecordedWithoutHostPath(t *testing.T) {
	runtimePath := filepath.Join(t.TempDir(), "docker")
	if err := os.WriteFile(runtimePath, []byte("fixture runtime"), 0o700); err != nil {
		t.Fatal(err)
	}
	workspace := t.TempDir()
	data := t.TempDir()
	if err := os.WriteFile(filepath.Join(data, "database.bin"), []byte("reviewed database"), 0o600); err != nil {
		t.Fatal(err)
	}
	plan, err := BuildOCIPlanWithData(runtimePath, workspace, strings.Repeat("a", 64), manifestWithDataMount(), map[string]string{
		GrypeDataMountName: data,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.DataMounts) != 1 || plan.DataMounts[0].Files != 1 || plan.DataMounts[0].Bytes != 17 ||
		len(plan.DataMounts[0].SHA256) != 64 {
		t.Fatalf("unexpected bound mount: %+v", plan.DataMounts)
	}
	expected := "--mount=type=bind,src=" + plan.DataMounts[0].Source + ",dst=" + grypeDataMountPath + ",readonly"
	if !slices.Contains(plan.Arguments, expected) {
		t.Fatalf("missing read-only mount %q in %v", expected, plan.Arguments)
	}

	output := validRunOutput("not_found")
	output.Transcript.Observations[0].Kind = GrypeObservationKind
	output.DataInputs = []model.AdapterDataInput{{
		Name: GrypeDataMountName, Destination: grypeDataMountPath,
		SHA256: plan.DataMounts[0].SHA256, Files: 1, Bytes: 17,
	}}
	execution, err := BindExecution(strings.Repeat("b", 64), Subject{
		TargetName: "fixture", InventoryDigest: strings.Repeat("c", 64),
	}, manifestWithDataMount(), output)
	if err != nil {
		t.Fatal(err)
	}
	payload, err := json.Marshal(execution)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(payload), data) || !strings.Contains(string(payload), plan.DataMounts[0].SHA256) {
		t.Fatalf("durable execution leaks path or omits identity: %s", payload)
	}
}

func TestDataMountFailsClosedOnMissingExtraSymlinkOrDrift(t *testing.T) {
	manifest := manifestWithDataMount()
	runtimePath := filepath.Join(t.TempDir(), "docker")
	if err := os.WriteFile(runtimePath, []byte("#!/bin/sh\nexit 0\n"), 0o700); err != nil {
		t.Fatal(err)
	}
	workspace := t.TempDir()
	if _, err := BuildOCIPlanWithData(runtimePath, workspace, strings.Repeat("a", 64), manifest, nil); err == nil {
		t.Fatal("missing required data mount was accepted")
	}
	data := t.TempDir()
	if err := os.WriteFile(filepath.Join(data, "database.bin"), []byte("one"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := BuildOCIPlanWithData(runtimePath, workspace, strings.Repeat("a", 64), manifest, map[string]string{
		GrypeDataMountName: data, "extra": data,
	}); err == nil {
		t.Fatal("undeclared data mount was accepted")
	}
	linkData := t.TempDir()
	if err := os.Symlink(filepath.Join(data, "database.bin"), filepath.Join(linkData, "link")); err != nil {
		t.Fatal(err)
	}
	if _, err := BuildOCIPlanWithData(runtimePath, workspace, strings.Repeat("a", 64), manifest, map[string]string{
		GrypeDataMountName: linkData,
	}); err == nil {
		t.Fatal("symlinked data content was accepted")
	}
	plan, err := BuildOCIPlanWithData(runtimePath, workspace, strings.Repeat("a", 64), manifest, map[string]string{
		GrypeDataMountName: data,
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(data, "database.bin"), []byte("two"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := RunOCI(context.Background(), plan, manifest, nil); err == nil || !strings.Contains(err.Error(), "changed") {
		t.Fatalf("unexpected data drift result: %v", err)
	}
}
