package inventory

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func writeFile(t *testing.T, root, relative, content string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(relative))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestBuildHashesContentAndDetectsEcosystems(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "app.py", "print('one')\n")
	writeFile(t, root, "requirements.txt", "example==1.0\n")
	writeFile(t, root, "requirements.lock.txt", "example==1.0\n")
	writeFile(t, root, ".github/workflows/ci.yml", "name: CI\n")
	writeFile(t, root, "Dockerfile", "FROM example@sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\n")
	writeFile(t, root, "infra/main.tf", "terraform {}\n")
	writeFile(t, root, "deploy/app.yaml", "apiVersion: apps/v1\nkind: Deployment\n")

	first, err := Build(root)
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Contains(first.PackageEcosystems, "python") {
		t.Fatalf("expected Python ecosystem, got %v", first.PackageEcosystems)
	}
	if !first.CI.GitHubActions {
		t.Fatal("expected GitHub Actions detection")
	}
	if first.SourceFiles != 1 {
		t.Fatalf("expected one source file, got %d", first.SourceFiles)
	}
	if first.SchemaVersion != "prc.inventory/v0.2" || len(first.Components) != 6 || len(first.Relations) != 5 {
		t.Fatalf("unexpected graph: schema=%s components=%+v relations=%+v", first.SchemaVersion, first.Components, first.Relations)
	}
	if !slices.Contains(first.ContainerFiles, "Dockerfile") ||
		!slices.Contains(first.Infrastructure.TerraformFiles, "infra/main.tf") ||
		!slices.Contains(first.Infrastructure.KubernetesFiles, "deploy/app.yaml") {
		t.Fatalf("missing deployment inventory: containers=%v infrastructure=%+v", first.ContainerFiles, first.Infrastructure)
	}
	if len(first.Facts) != 6 {
		t.Fatalf("expected six provenance facts, got %+v", first.Facts)
	}

	writeFile(t, root, "app.py", "print('two')\n")
	second, err := Build(root)
	if err != nil {
		t.Fatal(err)
	}
	if first.Digest == second.Digest {
		t.Fatal("inventory digest did not change with file content")
	}
}

func TestBuildSkipsCachesAndSymlinks(t *testing.T) {
	root := t.TempDir()
	outside := filepath.Join(t.TempDir(), "outside.py")
	if err := os.WriteFile(outside, []byte("secret\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	writeFile(t, root, "__pycache__/test_fake.pyc", "compiled")
	if err := os.Symlink(outside, filepath.Join(root, "linked.py")); err != nil {
		t.Fatal(err)
	}
	writeFile(t, root, "main.go", "package main\n")

	item, err := Build(root)
	if err != nil {
		t.Fatal(err)
	}
	if item.FileCount != 1 || item.Files[0].Path != "main.go" {
		t.Fatalf("unexpected inventory files: %+v", item.Files)
	}
	if !slices.Equal(item.Symlinks, []string{"linked.py"}) {
		t.Fatalf("symlinks = %v", item.Symlinks)
	}
}

func TestBuildGraphAndFactsAreDeterministic(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "package.json", "{}\n")
	writeFile(t, root, "Dockerfile.worker", "FROM example:1\n")
	first, err := Build(root)
	if err != nil {
		t.Fatal(err)
	}
	second, err := Build(root)
	if err != nil {
		t.Fatal(err)
	}
	if first.Digest != second.Digest {
		t.Fatalf("deterministic inventory digests differ: %s != %s", first.Digest, second.Digest)
	}
	if len(first.Facts) == 0 || first.Facts[0].DetectorVersion == "" || first.Facts[0].Confidence <= 0 {
		t.Fatalf("facts lack provenance: %+v", first.Facts)
	}
}

func TestBuildIncludesPermissionModeInIdentity(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "script.sh", "#!/bin/sh\n")
	before, err := Build(root)
	if err != nil {
		t.Fatal(err)
	}
	if before.Files[0].Mode != 0o644 {
		t.Fatalf("mode = %#o", before.Files[0].Mode)
	}
	if err := os.Chmod(filepath.Join(root, "script.sh"), 0o755); err != nil {
		t.Fatal(err)
	}
	after, err := Build(root)
	if err != nil {
		t.Fatal(err)
	}
	if after.Files[0].Mode != 0o755 || before.Digest == after.Digest {
		t.Fatalf("mode=%#o digest changed=%t", after.Files[0].Mode, before.Digest != after.Digest)
	}
}
