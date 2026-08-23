package adapter

import (
	"context"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func validManifest() Manifest {
	limits := DefaultLimits()
	return Manifest{
		SchemaVersion: ManifestSchema,
		ID:            "prc.adapter.fixture@0.1",
		Title:         "Fixture adapter",
		Trust:         "first-party-sandboxed",
		Runner:        "oci",
		Image:         "registry.example/prc/fixture@sha256:" + strings.Repeat("a", 64),
		Command:       []string{"/adapter", "scan"},
		Capabilities: Capabilities{
			ReadWorkspace: true, WriteScratch: true, Network: "deny",
			NetworkHosts: []string{}, SecretHandles: []string{},
		},
		Resources: Resources{
			TimeoutSeconds: 60, MemoryMB: 512, CPUs: 1, PIDs: 64, TmpfsMB: 64, Limits: limits,
		},
	}
}

func TestManifestFailsClosedOnMutableOrPrivilegedConfiguration(t *testing.T) {
	tests := map[string]func(*Manifest){
		"mutable image": func(value *Manifest) { value.Image = "registry.example/prc/fixture:latest" },
		"network":       func(value *Manifest) { value.Capabilities.Network = "allow" },
		"secret":        func(value *Manifest) { value.Capabilities.SecretHandles = []string{"TOKEN"} },
		"process":       func(value *Manifest) { value.Capabilities.ChildProcesses = true },
		"native runner": func(value *Manifest) { value.Runner = "process" },
		"short timeout": func(value *Manifest) { value.Resources.TimeoutSeconds = 0 },
	}
	for name, mutate := range tests {
		t.Run(name, func(t *testing.T) {
			manifest := validManifest()
			mutate(&manifest)
			if err := manifest.Validate(); err == nil {
				t.Fatal("expected manifest rejection")
			}
		})
	}
}

func TestLoadManifestRejectsUnknownFields(t *testing.T) {
	path := filepath.Join(t.TempDir(), "adapter.yaml")
	data := []byte("schema_version: prc.adapter-manifest/v0.1\nid: prc.adapter.fixture@0.1\nunexpected: true\n")
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadManifest(path); err == nil {
		t.Fatal("expected unknown field rejection")
	}
}

func TestCheckedInManifestLoads(t *testing.T) {
	manifest, err := LoadManifest(fixturePath("fixture-adapter.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	if manifest.ID != "prc.adapter.fixture@0.1" {
		t.Fatalf("unexpected fixture manifest: %+v", manifest)
	}
}

func TestOCIPlanContainsRequiredIsolationFlags(t *testing.T) {
	runtimeDirectory := t.TempDir()
	runtime := filepath.Join(runtimeDirectory, "docker")
	if err := os.WriteFile(runtime, []byte("#!/bin/sh\nexit 0\n"), 0o700); err != nil {
		t.Fatal(err)
	}
	workspace := t.TempDir()
	plan, err := BuildOCIPlan(runtime, workspace, strings.Repeat("b", 64), validManifest())
	if err != nil {
		t.Fatal(err)
	}
	required := []string{
		"--pull=never", "--network=none", "--read-only", "--cap-drop=ALL",
		"--security-opt=no-new-privileges", "--pids-limit=1", "--memory=512m",
		"--cpus=1", "--user=65532:65532", "--tmpfs=/tmp:rw,noexec,nosuid,size=64m",
	}
	for _, argument := range required {
		if !slices.Contains(plan.Arguments, argument) {
			t.Errorf("missing isolation argument %q in %v", argument, plan.Arguments)
		}
	}
	for _, argument := range plan.Arguments {
		if strings.HasPrefix(argument, "--mount=") && !strings.Contains(argument, "readonly") {
			t.Fatalf("workspace mount is not read-only: %s", argument)
		}
	}
}

func TestOCIPlanCannotBeExpandedAfterCapabilityEvaluation(t *testing.T) {
	runtime := filepath.Join(t.TempDir(), "docker")
	if err := os.WriteFile(runtime, []byte("fixture"), 0o700); err != nil {
		t.Fatal(err)
	}
	manifest := validManifest()
	plan, err := BuildOCIPlan(runtime, t.TempDir(), strings.Repeat("c", 64), manifest)
	if err != nil {
		t.Fatal(err)
	}
	plan.Arguments = append(plan.Arguments, "--privileged")
	if _, err := RunOCI(context.Background(), plan, manifest, nil); err == nil {
		t.Fatal("mutated OCI plan was accepted")
	}
}

func TestOCIInputLimitIsEnforcedBeforeExecution(t *testing.T) {
	runtime := filepath.Join(t.TempDir(), "docker")
	if err := os.WriteFile(runtime, []byte("fixture"), 0o700); err != nil {
		t.Fatal(err)
	}
	manifest := validManifest()
	plan, err := BuildOCIPlan(runtime, t.TempDir(), strings.Repeat("d", 64), manifest)
	if err != nil {
		t.Fatal(err)
	}
	input := make([]byte, manifest.Resources.MaxStdin+1)
	if _, err := RunOCI(context.Background(), plan, manifest, input); err == nil {
		t.Fatal("oversized OCI input was accepted")
	}
}

func TestOCIRuntimeCannotChangeAfterPlanning(t *testing.T) {
	runtime := filepath.Join(t.TempDir(), "docker")
	if err := os.WriteFile(runtime, []byte("first"), 0o700); err != nil {
		t.Fatal(err)
	}
	manifest := validManifest()
	plan, err := BuildOCIPlan(runtime, t.TempDir(), strings.Repeat("e", 64), manifest)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(runtime, []byte("second"), 0o700); err != nil {
		t.Fatal(err)
	}
	if _, err := RunOCI(context.Background(), plan, manifest, nil); err == nil {
		t.Fatal("changed OCI runtime was accepted")
	}
}

func TestBoundedBufferPreservesLimitAndReturnsError(t *testing.T) {
	buffer := newBoundedBuffer(4)
	written, err := buffer.Write([]byte("abcdef"))
	if err == nil || written != 4 || buffer.String() != "abcd" {
		t.Fatalf("written=%d err=%v content=%q", written, err, buffer.String())
	}
}
