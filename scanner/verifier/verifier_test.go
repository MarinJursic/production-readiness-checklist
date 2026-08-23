package verifier

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/MarinJursic/production-readiness-checklist/scanner/inventory"
)

const testImage = "registry.example/prc/go-verifier@sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

func fakeRuntime(t *testing.T, script string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "docker")
	if err := os.WriteFile(path, []byte("#!/bin/sh\n"+script+"\n"), 0o700); err != nil {
		t.Fatal(err)
	}
	return path
}

func workspace(t *testing.T) (string, string) {
	t.Helper()
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("module example.test/app\n\ngo 1.27.0\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	item, err := inventoryForTest(root)
	if err != nil {
		t.Fatal(err)
	}
	return root, item
}

func inventoryForTest(root string) (string, error) {
	item, err := inventory.Build(root)
	if err != nil {
		return "", err
	}
	return item.Digest, nil
}

func TestBuildPlanOwnsCommandAndIsolation(t *testing.T) {
	root, digest := workspace(t)
	options := Defaults(fakeRuntime(t, "exit 0"), testImage, "go")
	plan, err := BuildPlan(root, digest, options)
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(plan.Command, []string{"go", "test", "./..."}) {
		t.Fatalf("command = %v", plan.Command)
	}
	required := []string{
		"--pull=never", "--network=none", "--read-only", "--cap-drop=ALL",
		"--security-opt=no-new-privileges=true", "--memory=1024m", "--memory-swap=1024m",
		"--cpus=1", "--pids-limit=128", "--ulimit=nofile=1024:1024",
		"--tmpfs=/tmp:rw,noexec,nosuid,nodev,mode=1777,size=256m",
		"--tmpfs=/prc-exec:rw,exec,nosuid,nodev,uid=" + strconv.Itoa(os.Getuid()) +
			",gid=" + strconv.Itoa(os.Getgid()) + ",mode=0700,size=256m",
		"--env=TMPDIR=/tmp", "--env=GOTMPDIR=/prc-exec",
		"--env=GOCACHE=/prc-exec/go-build", "--env=GOMODCACHE=/prc-exec/go-mod",
	}
	for _, value := range required {
		if !slices.Contains(plan.Arguments, value) {
			t.Errorf("missing %q in %v", value, plan.Arguments)
		}
	}
	for _, value := range plan.Arguments {
		if strings.HasPrefix(value, "--mount=") && !strings.Contains(value, "readonly") {
			t.Fatalf("writable workspace mount: %s", value)
		}
	}
}

func TestRunRecordsPassAndFailedTests(t *testing.T) {
	for name, script := range map[string]string{
		"pass": "printf passed; printf diagnostic >&2; exit 0",
		"fail": "printf failed; exit 1",
	} {
		t.Run(name, func(t *testing.T) {
			root, digest := workspace(t)
			plan, err := BuildPlan(root, digest, Defaults(fakeRuntime(t, script), testImage, "go"))
			if err != nil {
				t.Fatal(err)
			}
			execution, err := Run(context.Background(), plan)
			if err != nil {
				t.Fatal(err)
			}
			if execution.Outcome != name || execution.WorkspaceUnchanged != true || execution.ExecutionID == "" ||
				execution.Stdout.Bytes == 0 {
				t.Fatalf("execution = %+v", execution)
			}
		})
	}
}

func TestRunRecordsRuntimeFailureAndWorkspaceMutation(t *testing.T) {
	root, digest := workspace(t)
	plan, err := BuildPlan(root, digest, Defaults(fakeRuntime(t, "exit 125"), testImage, "go"))
	if err != nil {
		t.Fatal(err)
	}
	execution, err := Run(context.Background(), plan)
	if err != nil || execution.Outcome != "infrastructure_error" || execution.ReasonCode != "runtime_failed" {
		t.Fatalf("execution=%+v err=%v", execution, err)
	}

	mutator := fakeRuntime(t, "printf changed > "+shellQuote(filepath.Join(root, "go.mod"))+"; exit 0")
	plan, err = BuildPlan(root, digest, Defaults(mutator, testImage, "go"))
	if err != nil {
		t.Fatal(err)
	}
	execution, err = Run(context.Background(), plan)
	if err != nil || execution.Outcome != "workspace_changed" || execution.WorkspaceUnchanged {
		t.Fatalf("execution=%+v err=%v", execution, err)
	}
}

func TestPlanMutationAndRuntimeReplacementFailClosed(t *testing.T) {
	root, digest := workspace(t)
	runtimePath := fakeRuntime(t, "exit 0")
	plan, err := BuildPlan(root, digest, Defaults(runtimePath, testImage, "go"))
	if err != nil {
		t.Fatal(err)
	}
	plan.Arguments = append(plan.Arguments, "--privileged")
	if _, err := Run(context.Background(), plan); err == nil {
		t.Fatal("mutated plan was accepted")
	}
	plan, err = BuildPlan(root, digest, Defaults(runtimePath, testImage, "go"))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(runtimePath, []byte("changed"), 0o700); err != nil {
		t.Fatal(err)
	}
	if _, err := Run(context.Background(), plan); err == nil {
		t.Fatal("replaced runtime was accepted")
	}
}

func TestExecutionEvidenceTamperingFailsClosed(t *testing.T) {
	root, digest := workspace(t)
	plan, err := BuildPlan(root, digest, Defaults(fakeRuntime(t, "exit 0"), testImage, "go"))
	if err != nil {
		t.Fatal(err)
	}
	execution, err := Run(context.Background(), plan)
	if err != nil {
		t.Fatal(err)
	}

	tamperedPolicy := execution
	tamperedPolicy.Policy.Network = "host"
	tamperedPolicy.ExecutionID, _ = executionID(tamperedPolicy)
	if err := tamperedPolicy.Validate(); err == nil {
		t.Fatal("tampered sandbox policy was accepted")
	}
	tamperedRuntime := execution
	tamperedRuntime.RuntimeSHA256 = strings.Repeat("f", 64)
	tamperedRuntime.ExecutionID, _ = executionID(tamperedRuntime)
	if err := tamperedRuntime.Validate(); err == nil {
		t.Fatal("rewritten runtime identity was accepted without a matching plan")
	}
	tamperedOutcome := execution
	tamperedOutcome.Outcome = "fail"
	tamperedOutcome.ReasonCode = "tests_failed"
	tamperedOutcome.ExecutionID, _ = executionID(tamperedOutcome)
	if err := tamperedOutcome.Validate(); err == nil {
		t.Fatal("failed-test outcome with a zero exit status was accepted")
	}
}

func TestInferKindAndMutableImageFailClosed(t *testing.T) {
	if kind, err := InferKind("service.go"); err != nil || kind != "go" {
		t.Fatalf("kind=%q err=%v", kind, err)
	}
	if kind, err := InferKind("service.js"); err != nil || kind != "javascript" {
		t.Fatalf("kind=%q err=%v", kind, err)
	}
	if _, err := InferKind("service.ts"); err == nil {
		t.Fatal("unsupported source kind was accepted")
	}
	root, digest := workspace(t)
	options := Defaults(fakeRuntime(t, "exit 0"), "registry.example/prc/verifier:latest", "go")
	if _, err := BuildPlan(root, digest, options); err == nil {
		t.Fatal("mutable verifier image was accepted")
	}
}

func TestOCIIntegrationRunsPassingAndFailingJavaScript(t *testing.T) {
	image := os.Getenv("PRC_TEST_VERIFIER_IMAGE")
	if image == "" {
		t.Skip("set PRC_TEST_VERIFIER_IMAGE to an already-present digest-pinned Node image")
	}
	if _, err := exec.LookPath("docker"); err != nil {
		t.Fatal("PRC_TEST_VERIFIER_IMAGE requires a docker executable")
	}
	root := t.TempDir()
	write := func(path, content string) {
		t.Helper()
		if err := os.WriteFile(filepath.Join(root, path), []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	write("add.js", "exports.add = (left, right) => left + right;\n")
	write("add.test.js", "const test = require('node:test');\nconst assert = require('node:assert/strict');\nconst { add } = require('./add');\ntest('adds', () => assert.equal(add(2, 3), 5));\n")
	write("scratch.test.js", "const test = require('node:test');\nconst fs = require('node:fs');\ntest('scanner scratch is writable', () => fs.writeFileSync('/prc-exec/probe', 'ok'));\n")
	run := func() Execution {
		t.Helper()
		digest, err := inventoryForTest(root)
		if err != nil {
			t.Fatal(err)
		}
		plan, err := BuildPlan(root, digest, Defaults("docker", image, "javascript"))
		if err != nil {
			t.Fatal(err)
		}
		execution, err := Run(context.Background(), plan)
		if err != nil {
			t.Fatal(err)
		}
		return execution
	}
	if execution := run(); execution.Outcome != "pass" {
		t.Fatalf("passing suite outcome = %+v", execution)
	}
	write("add.test.js", "const test = require('node:test');\nconst assert = require('node:assert/strict');\nconst { add } = require('./add');\ntest('adds', () => assert.equal(add(2, 3), 6));\n")
	if execution := run(); execution.Outcome != "fail" || execution.ExitCode == 0 {
		t.Fatalf("failing suite outcome = %+v", execution)
	}
}

func shellQuote(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "'\\''") + "'"
}
