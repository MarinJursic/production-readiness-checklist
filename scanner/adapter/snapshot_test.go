package adapter

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/MarinJursic/production-readiness-checklist/scanner/inventory"
)

func TestPrepareSnapshotMaterializesOnlySealedRegularFiles(t *testing.T) {
	root := t.TempDir()
	writeSnapshotFixture(t, root, "src/main.go", "package main\n")
	writeSnapshotFixture(t, root, ".git/private", "excluded\n")
	writeSnapshotFixture(t, root, "site/generated.html", "excluded\n")
	if err := os.Symlink("src/main.go", filepath.Join(root, "source-link")); err != nil {
		t.Fatal(err)
	}
	item, err := inventory.Build(root)
	if err != nil {
		t.Fatal(err)
	}
	snapshot, err := PrepareSnapshot(item)
	if err != nil {
		t.Fatal(err)
	}
	defer snapshot.Close()
	if snapshot.Files != 1 || snapshot.Bytes != int64(len("package main\n")) ||
		!hexDigestPattern.MatchString(snapshot.Digest) {
		t.Fatalf("snapshot identity = %+v", snapshot)
	}
	data, err := os.ReadFile(filepath.Join(snapshot.Path, "src", "main.go"))
	if err != nil || string(data) != "package main\n" {
		t.Fatalf("snapshot content = %q, %v", data, err)
	}
	for _, excluded := range []string{".git/private", "site/generated.html", "source-link"} {
		if _, err := os.Lstat(filepath.Join(snapshot.Path, filepath.FromSlash(excluded))); !os.IsNotExist(err) {
			t.Fatalf("excluded path %s was materialized: %v", excluded, err)
		}
	}
	info, err := os.Stat(filepath.Join(snapshot.Path, "src", "main.go"))
	if err != nil || info.Mode().Perm() != 0o400 {
		t.Fatalf("snapshot file mode = %v, %v", info.Mode().Perm(), err)
	}
}

func TestPrepareSnapshotRejectsWorkspaceDrift(t *testing.T) {
	root := t.TempDir()
	writeSnapshotFixture(t, root, "README.md", "before\n")
	item, err := inventory.Build(root)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "README.md"), []byte("after\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := PrepareSnapshot(item); err == nil || !strings.Contains(err.Error(), "changed after inventory") {
		t.Fatalf("workspace drift error = %v", err)
	}
}

func TestGitleaksSnapshotRelocatesButPreservesIgnoreFileBytes(t *testing.T) {
	root := t.TempDir()
	writeSnapshotFixture(t, root, gitleaksIgnoreSourcePath, "untrusted-fingerprint\n")
	writeSnapshotFixture(t, root, "README.md", "sealed\n")
	item, err := inventory.Build(root)
	if err != nil {
		t.Fatal(err)
	}
	snapshot, err := PrepareSnapshotForManifest(item, validGitleaksManifest())
	if err != nil {
		t.Fatal(err)
	}
	defer snapshot.Close()
	if _, err := os.Stat(filepath.Join(snapshot.Path, gitleaksIgnoreSourcePath)); !os.IsNotExist(err) {
		t.Fatalf("target-controlled ignore file remains active: %v", err)
	}
	data, err := os.ReadFile(filepath.Join(snapshot.Path, filepath.FromSlash(gitleaksIgnoreSnapshotPath)))
	if err != nil || string(data) != "untrusted-fingerprint\n" {
		t.Fatalf("relocated ignore content = %q, %v", data, err)
	}
}

func TestSyftSnapshotInjectsScannerOwnedConfiguration(t *testing.T) {
	root := t.TempDir()
	writeSnapshotFixture(t, root, ".syft.yaml", "output: table\n")
	writeSnapshotFixture(t, root, ".prc/syft-config.yaml", "target-owned: true\n")
	writeSnapshotFixture(t, root, "go.mod", "module example.com/fixture\n\ngo 1.27\n")
	item, err := inventory.Build(root)
	if err != nil {
		t.Fatal(err)
	}
	snapshot, err := PrepareSnapshotForManifest(item, validSyftManifest())
	if err != nil {
		t.Fatal(err)
	}
	defer snapshot.Close()
	data, err := os.ReadFile(filepath.Join(snapshot.Path, filepath.FromSlash(SyftConfigSnapshotPath)))
	if err != nil || string(data) != string(syftConfig) {
		t.Fatalf("scanner-owned Syft configuration = %q, %v", data, err)
	}
	if snapshot.Files != 2 || snapshot.Bytes <= 0 {
		t.Fatalf("snapshot inventory accounting included protocol input: %+v", snapshot)
	}
	if untrusted, err := os.ReadFile(filepath.Join(snapshot.Path, ".syft.yaml")); err != nil || string(untrusted) != "output: table\n" {
		t.Fatalf("target file was not preserved as scan input: %q, %v", untrusted, err)
	}
}

func TestGrypeSnapshotInjectsScannerOwnedConfiguration(t *testing.T) {
	root := t.TempDir()
	writeSnapshotFixture(t, root, ".grype.yaml", "ignore: everything\n")
	writeSnapshotFixture(t, root, ".prc/grype-config.yaml", "target-owned: true\n")
	writeSnapshotFixture(t, root, "package-lock.json", "{}\n")
	item, err := inventory.Build(root)
	if err != nil {
		t.Fatal(err)
	}
	snapshot, err := PrepareSnapshotForManifest(item, manifestWithDataMount())
	if err != nil {
		t.Fatal(err)
	}
	defer snapshot.Close()
	data, err := os.ReadFile(filepath.Join(snapshot.Path, filepath.FromSlash(GrypeConfigSnapshotPath)))
	if err != nil || string(data) != string(grypeConfig) {
		t.Fatalf("scanner-owned Grype configuration = %q, %v", data, err)
	}
	if snapshot.Files != 2 || snapshot.Bytes <= 0 {
		t.Fatalf("snapshot inventory accounting included protocol input: %+v", snapshot)
	}
	if untrusted, err := os.ReadFile(filepath.Join(snapshot.Path, ".grype.yaml")); err != nil || string(untrusted) != "ignore: everything\n" {
		t.Fatalf("target file was not preserved as scan input: %q, %v", untrusted, err)
	}
}

func TestRunOCIRejectsSnapshotDriftBeforeRuntime(t *testing.T) {
	root := t.TempDir()
	writeSnapshotFixture(t, root, "README.md", "sealed\n")
	item, err := inventory.Build(root)
	if err != nil {
		t.Fatal(err)
	}
	snapshot, err := PrepareSnapshot(item)
	if err != nil {
		t.Fatal(err)
	}
	defer snapshot.Close()
	runtimePath := filepath.Join(t.TempDir(), "docker")
	if err := os.WriteFile(runtimePath, []byte("#!/bin/sh\nexit 99\n"), 0o700); err != nil {
		t.Fatal(err)
	}
	plan, err := BuildSnapshotOCIPlan(runtimePath, snapshot, strings.Repeat("a", 64), validManifest())
	if err != nil {
		t.Fatal(err)
	}
	file := filepath.Join(snapshot.Path, "README.md")
	if err := os.Chmod(file, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(file, []byte("tampered\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := RunOCI(context.Background(), plan, validManifest(), nil); err == nil ||
		!strings.Contains(err.Error(), "snapshot changed") {
		t.Fatalf("snapshot drift execution error = %v", err)
	}
}

func TestSnapshotCloseIsScopedAndIdempotent(t *testing.T) {
	root := t.TempDir()
	writeSnapshotFixture(t, root, "README.md", "sealed\n")
	item, err := inventory.Build(root)
	if err != nil {
		t.Fatal(err)
	}
	snapshot, err := PrepareSnapshot(item)
	if err != nil {
		t.Fatal(err)
	}
	path := snapshot.Path
	if err := snapshot.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("snapshot path remains after close: %v", err)
	}
	if err := snapshot.Close(); err != nil {
		t.Fatal(err)
	}
}

func writeSnapshotFixture(t *testing.T, root, relative, content string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(relative))
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}
