package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestVersionCommand(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if code := run([]string{"version"}, &stdout, &stderr); code != 0 {
		t.Fatalf("exit=%d stderr=%s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "prc 0.1.0-dev") {
		t.Fatalf("unexpected output %q", stdout.String())
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

func TestAdapterValidateOutputRejectsAuthorityAttack(t *testing.T) {
	var stdout, stderr bytes.Buffer
	path := filepath.Join("..", "..", "fixtures", "adapters", "malicious-authority-output.jsonl")
	if code := run([]string{"adapter", "validate-output", "--file", path}, &stdout, &stderr); code != 2 {
		t.Fatalf("exit=%d stderr=%s", code, stderr.String())
	}
}

func TestUnknownCommandIsUsageError(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if code := run([]string{"unknown"}, &stdout, &stderr); code != 2 {
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
