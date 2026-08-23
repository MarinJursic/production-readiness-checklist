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
	manifest.DataMounts = []DataMount{{
		Name: "fixture-db", Destination: "/prc-inputs/fixture-db",
		MaxFiles: 4, MaxBytes: 1024,
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
		"fixture-db": data,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.DataMounts) != 1 || plan.DataMounts[0].Files != 1 || plan.DataMounts[0].Bytes != 17 ||
		len(plan.DataMounts[0].SHA256) != 64 {
		t.Fatalf("unexpected bound mount: %+v", plan.DataMounts)
	}
	expected := "--mount=type=bind,src=" + plan.DataMounts[0].Source + ",dst=/prc-inputs/fixture-db,readonly"
	if !slices.Contains(plan.Arguments, expected) {
		t.Fatalf("missing read-only mount %q in %v", expected, plan.Arguments)
	}

	output := validRunOutput("not_found")
	output.DataInputs = []model.AdapterDataInput{{
		Name: "fixture-db", Destination: "/prc-inputs/fixture-db",
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
		"fixture-db": data, "extra": data,
	}); err == nil {
		t.Fatal("undeclared data mount was accepted")
	}
	linkData := t.TempDir()
	if err := os.Symlink(filepath.Join(data, "database.bin"), filepath.Join(linkData, "link")); err != nil {
		t.Fatal(err)
	}
	if _, err := BuildOCIPlanWithData(runtimePath, workspace, strings.Repeat("a", 64), manifest, map[string]string{
		"fixture-db": linkData,
	}); err == nil {
		t.Fatal("symlinked data content was accepted")
	}
	plan, err := BuildOCIPlanWithData(runtimePath, workspace, strings.Repeat("a", 64), manifest, map[string]string{
		"fixture-db": data,
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
