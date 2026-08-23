package doctor

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

func catalogRoot(t *testing.T) string {
	t.Helper()
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	return root
}

func checkByID(t *testing.T, report Report, id string) Check {
	t.Helper()
	for _, check := range report.Checks {
		if check.ID == id {
			return check
		}
	}
	t.Fatalf("missing check %s", id)
	return Check{}
}

func TestRunValidatesRequestedCapabilitiesWithoutExecutingTools(t *testing.T) {
	target := t.TempDir()
	if err := os.WriteFile(filepath.Join(target, "app.go"), []byte("package app\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	state := t.TempDir()
	if err := os.Chmod(state, 0o700); err != nil {
		t.Fatal(err)
	}
	candidateParent := t.TempDir()
	bin := t.TempDir()
	docker := filepath.Join(bin, "docker")
	codex := filepath.Join(bin, "codex")
	for _, path := range []string{docker, codex} {
		if err := os.WriteFile(path, []byte("#!/bin/sh\nexit 99\n"), 0o700); err != nil {
			t.Fatal(err)
		}
	}
	now := time.Date(2026, 8, 23, 12, 0, 0, 0, time.UTC)
	report := Run(Options{
		Target: target, CatalogRoot: catalogRoot(t), StateDirectory: state,
		CandidateParent: candidateParent, OCIRuntime: docker, Providers: []string{codex},
		Now: func() time.Time { return now },
	})
	if !report.Ready || report.SchemaVersion != Schema || report.GeneratedAt != now || report.Summary.Failed != 0 {
		t.Fatalf("unexpected doctor report: %+v", report)
	}
	for _, id := range []string{"target.inventory", "catalog.load", "state-store", "candidate-workspace", "oci-runtime", "provider.codex"} {
		if check := checkByID(t, report, id); check.Status != "pass" || !check.Required {
			t.Fatalf("check %s = %+v", id, check)
		}
	}
	if entries, err := os.ReadDir(state); err != nil || len(entries) != 0 {
		t.Fatalf("state probe was not cleaned up: %v %v", entries, err)
	}
	if entries, err := os.ReadDir(candidateParent); err != nil || len(entries) != 0 {
		t.Fatalf("candidate probe was not cleaned up: %v %v", entries, err)
	}
}

func TestRunReportsOptionalCapabilitiesAsWarnings(t *testing.T) {
	report := Run(Options{Target: t.TempDir(), CatalogRoot: catalogRoot(t)})
	if !report.Ready || report.Summary.Warnings != 3 || report.Summary.Failed != 0 {
		t.Fatalf("unexpected optional-capability report: %+v", report)
	}
	for _, id := range []string{"state-store", "candidate-workspace", "oci-runtime"} {
		if check := checkByID(t, report, id); check.Status != "warn" || check.Required {
			t.Fatalf("check %s = %+v", id, check)
		}
	}
}

func TestRunAllowsCandidateParentThatContainsTargetWhenProbeIsSibling(t *testing.T) {
	target := t.TempDir()
	report := Run(Options{
		Target: target, CatalogRoot: catalogRoot(t), CandidateParent: filepath.Dir(target),
	})
	if check := checkByID(t, report, "candidate-workspace"); check.Status != "pass" {
		t.Fatalf("safe sibling candidate layout was rejected: %+v", check)
	}
}

func TestRunFailsClosedForOverlappingCandidateAndPublicState(t *testing.T) {
	target := t.TempDir()
	state := t.TempDir()
	if runtime.GOOS != "windows" {
		if err := os.Chmod(state, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	report := Run(Options{
		Target: target, CatalogRoot: catalogRoot(t), StateDirectory: state,
		CandidateParent: filepath.Join(target),
	})
	if runtime.GOOS == "windows" {
		if report.Summary.Failed != 1 {
			t.Fatalf("unexpected Windows failures: %+v", report)
		}
	} else if report.Summary.Failed != 2 {
		t.Fatalf("unexpected failures: %+v", report)
	}
	if report.Ready || checkByID(t, report, "candidate-workspace").Status != "fail" {
		t.Fatalf("unsafe environment was accepted: %+v", report)
	}
}

func TestRunRejectsUnexpectedExecutableIdentity(t *testing.T) {
	bin := t.TempDir()
	path := filepath.Join(bin, "not-docker")
	if err := os.WriteFile(path, []byte("#!/bin/sh\n"), 0o700); err != nil {
		t.Fatal(err)
	}
	report := Run(Options{Target: t.TempDir(), CatalogRoot: catalogRoot(t), OCIRuntime: path})
	check := checkByID(t, report, "oci-runtime")
	if report.Ready || check.Status != "fail" {
		t.Fatalf("unexpected executable was accepted: %+v", report)
	}
}
