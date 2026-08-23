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
