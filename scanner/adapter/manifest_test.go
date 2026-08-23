package adapter

import (
	"context"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
)

func validManifest() Manifest {
	limits := DefaultLimits()
	return Manifest{
		SchemaVersion: ManifestSchema, ID: "prc.adapter.fixture@0.1",
		Title: "Fixture adapter", Description: "Produces bounded fixture observations for protocol tests.",
		Publisher: Publisher{ID: "prc-project", Name: "Production Readiness Checklist"},
		Owner:     "scanner-maintainers", Maintenance: "active",
		Protocol: ProtocolVersion, OutputSchema: OutputSchemaVersion,
		ObservationKinds: []string{"fixture-result"},
		Compatibility:    Compatibility{EngineAPIs: []string{model.EngineVersion}},
		Tool: Tool{
			Name: "prc-fixture", Version: "1.0.0", Upstream: "https://example.com/prc-fixture",
			Formats: []ToolFormat{{Name: "fixture-json", Versions: []string{"1.0"}}},
		},
		Limitations: []string{"This adapter is a protocol fixture and does not inspect production projects."},
		Runner:      "oci", Image: "registry.example/prc/fixture@sha256:" + strings.Repeat("a", 64),
		Command: []string{"/adapter", "scan"},
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
		"mutable image":  func(value *Manifest) { value.Image = "registry.example/prc/fixture:latest" },
		"network":        func(value *Manifest) { value.Capabilities.Network = "allow" },
		"secret":         func(value *Manifest) { value.Capabilities.SecretHandles = []string{"TOKEN"} },
		"process":        func(value *Manifest) { value.Capabilities.ChildProcesses = true },
		"native runner":  func(value *Manifest) { value.Runner = "process" },
		"short timeout":  func(value *Manifest) { value.Resources.TimeoutSeconds = 0 },
		"protocol":       func(value *Manifest) { value.Protocol = "prc-adapter-jsonl-v2" },
		"latest tool":    func(value *Manifest) { value.Tool.Version = "latest" },
		"no limitations": func(value *Manifest) { value.Limitations = nil },
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
	data := []byte("schema_version: prc.adapter-manifest/v0.4\nid: prc.adapter.fixture@0.1\nunexpected: true\n")
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

func TestManifestCompatibilityAndObservationContractFailClosed(t *testing.T) {
	manifest := validManifest()
	manifest.Compatibility.EngineAPIs = []string{"prc.engine/v9.9"}
	if err := manifest.Validate(); err != nil {
		t.Fatalf("structurally valid future compatibility was rejected: %v", err)
	}
	if err := manifest.ValidateForCurrentEngine(); err == nil || !strings.Contains(err.Error(), "does not support") {
		t.Fatalf("unexpected compatibility error: %v", err)
	}

	manifest = validManifest()
	output := validRunOutput("not_found")
	output.Transcript.Observations[0].Kind = "undeclared-result"
	if _, err := BindExecution(strings.Repeat("f", 64), Subject{
		TargetName: "fixture", InventoryDigest: strings.Repeat("e", 64),
	}, manifest, output); err == nil || !strings.Contains(err.Error(), "undeclared observation kind") {
		t.Fatalf("unexpected observation contract error: %v", err)
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
		"--security-opt=no-new-privileges=true", "--pids-limit=1", "--memory=512m",
		"--memory-swap=512m", "--cpus=1",
		"--user=" + strconv.Itoa(os.Getuid()) + ":" + strconv.Itoa(os.Getgid()),
		"--ulimit=nofile=1024:1024", "--tmpfs=/tmp:rw,noexec,nosuid,nodev,mode=1777,size=64m",
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
