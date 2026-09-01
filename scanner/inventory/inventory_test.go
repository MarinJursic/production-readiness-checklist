package inventory

import (
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	projectconfig "github.com/MarinJursic/production-readiness-checklist/scanner/config"
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
	writeFile(t, root, "api/openapi.yaml", "openapi: 3.2.0\ninfo: {title: Example, version: 1.0.0}\npaths: {}\n")

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
	if first.SchemaVersion != "prc.inventory/v0.3" || len(first.Components) != 7 || len(first.Relations) != 6 {
		t.Fatalf("unexpected graph: schema=%s components=%+v relations=%+v", first.SchemaVersion, first.Components, first.Relations)
	}
	if !slices.Contains(first.ContainerFiles, "Dockerfile") ||
		!slices.Contains(first.Infrastructure.TerraformFiles, "infra/main.tf") ||
		!slices.Contains(first.Infrastructure.KubernetesFiles, "deploy/app.yaml") {
		t.Fatalf("missing deployment inventory: containers=%v infrastructure=%+v", first.ContainerFiles, first.Infrastructure)
	}
	if len(first.Facts) != 8 {
		t.Fatalf("expected eight provenance facts, got %+v", first.Facts)
	}
	openAPIDetected := false
	for _, component := range first.Components {
		if component.Kind == "api-description" && component.Ecosystem == "openapi" && component.Path == "api/openapi.yaml" {
			openAPIDetected = true
		}
	}
	if !openAPIDetected {
		t.Fatalf("missing OpenAPI inventory component: %+v", first.Components)
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

func TestBuildDiscoversOpenAPIByRootMarkerAndKnownName(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "contracts/service.yml", "openapi: 3.1.2\ninfo: {title: Service, version: 1}\npaths: {}\n")
	writeFile(t, root, "docs/openapi.json", "{not-json\n")
	item, err := Build(root)
	if err != nil {
		t.Fatal(err)
	}
	paths := []string{}
	for _, component := range item.Components {
		if component.Kind == "api-description" && component.Ecosystem == "openapi" {
			paths = append(paths, component.Path)
		}
	}
	if !slices.Equal(paths, []string{"contracts/service.yml", "docs/openapi.json"}) {
		t.Fatalf("OpenAPI components = %v", paths)
	}
	facts := 0
	for _, fact := range item.Facts {
		if fact.Key == "api.openapi" && fact.Detector == "prc.inventory.openapi-document" && len(fact.Limitations) == 1 {
			facts++
		}
	}
	if facts != 2 {
		t.Fatalf("OpenAPI provenance facts = %+v", item.Facts)
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

func TestBuildInventoriesAmbiguousSiteAndVendorDirectoriesAndReportsExclusions(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "site/app.py", "print('site')\n")
	writeFile(t, root, "vendor/first_party.go", "package vendor\n")
	writeFile(t, root, "node_modules/dependency.js", "module.exports = {}\n")
	item, err := Build(root)
	if err != nil {
		t.Fatal(err)
	}
	paths := []string{}
	for _, record := range item.Files {
		paths = append(paths, record.Path)
	}
	if !slices.Equal(paths, []string{"site/app.py", "vendor/first_party.go"}) {
		t.Fatalf("ambiguous first-party directories were skipped or dependencies were included: %v", paths)
	}
	foundExclusion := false
	for _, fact := range item.Facts {
		if fact.Key == "repository.exclusion" && fact.Source == "node_modules" {
			foundExclusion = true
		}
	}
	if !foundExclusion {
		t.Fatalf("inventory exclusion was not reported: %+v", item.Facts)
	}
}

func TestBuildSkipsOnlyTheDefaultMkDocsSiteOutput(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "mkdocs.yml", "site_name: Example\n")
	writeFile(t, root, "docs/index.md", "# Example\n")
	writeFile(t, root, "site/index.html", "<html>generated</html>")
	item, err := Build(root)
	if err != nil {
		t.Fatal(err)
	}
	paths := []string{}
	found := false
	for _, record := range item.Files {
		paths = append(paths, record.Path)
	}
	for _, fact := range item.Facts {
		if fact.Key == "repository.exclusion" && fact.Source == "site" && strings.Contains(fact.Value, "MkDocs") {
			found = true
		}
	}
	if slices.Contains(paths, "site/index.html") || !found {
		t.Fatalf("default MkDocs output was not excluded transparently: paths=%v facts=%+v", paths, item.Facts)
	}

	explicit := t.TempDir()
	writeFile(t, explicit, "mkdocs.yml", "site_name: Example\nsite_dir: public-docs\n")
	writeFile(t, explicit, "site/source.html", "<html>authored</html>\n")
	second, err := Build(explicit)
	if err != nil {
		t.Fatal(err)
	}
	if len(second.Files) != 2 || second.Files[1].Path != "site/source.html" {
		t.Fatalf("site was excluded despite an explicit non-default MkDocs output: %+v", second.Files)
	}

	for name, declaration := range map[string]string{
		"quoted":   "'site_dir': public-docs\n",
		"indented": "  site_dir: public-docs\n",
	} {
		t.Run(name, func(t *testing.T) {
			root := t.TempDir()
			writeFile(t, root, "mkdocs.yml", "site_name: Example\n"+declaration)
			writeFile(t, root, "site/source.html", "generated but conservatively retained\n")
			item, buildErr := Build(root)
			if buildErr != nil {
				t.Fatal(buildErr)
			}
			found := false
			for _, file := range item.Files {
				found = found || file.Path == "site/source.html"
			}
			if !found {
				t.Fatalf("site was excluded despite an explicit %s site_dir key", name)
			}
		})
	}
}

func TestBuildSkipsOnlyClearOrMarkedGeneratedDirectories(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "main.go", "package main\n")
	writeFile(t, root, ".next/cache/bundle.bin", "generated\n")
	writeFile(t, root, ".terraform/providers/provider.bin", "generated\n")
	writeFile(t, root, "rust/Cargo.toml", "[package]\nname = \"sample\"\nversion = \"0.1.0\"\n")
	writeFile(t, root, "rust/target/debug/sample", "generated\n")
	writeFile(t, root, "cmake/CMakeLists.txt", "cmake_minimum_required(VERSION 3.20)\n")
	writeFile(t, root, "cmake/build/app", "generated\n")
	writeFile(t, root, "generated/coverage/lcov.info", "TN:\n")
	writeFile(t, root, "generated/coverage/lcov-report/index.html", "generated\n")
	writeFile(t, root, "ambiguous/target/source.go", "package target\n")
	writeFile(t, root, "ambiguous/build/source.go", "package build\n")
	writeFile(t, root, "ambiguous/coverage/source.go", "package coverage\n")

	item, err := Build(root)
	if err != nil {
		t.Fatal(err)
	}
	paths := make([]string, 0, len(item.Files))
	for _, record := range item.Files {
		paths = append(paths, record.Path)
	}
	want := []string{
		"ambiguous/build/source.go", "ambiguous/coverage/source.go", "ambiguous/target/source.go",
		"cmake/CMakeLists.txt", "main.go", "rust/Cargo.toml",
	}
	if !slices.Equal(paths, want) {
		t.Fatalf("generated and ambiguous directory handling = %v, want %v", paths, want)
	}
	wantExclusions := map[string]bool{
		".next": true, ".terraform": true, "cmake/build": true,
		"generated/coverage": true, "rust/target": true,
	}
	for _, fact := range item.Facts {
		if fact.Key == "repository.exclusion" {
			delete(wantExclusions, fact.Source)
		}
	}
	if len(wantExclusions) != 0 {
		t.Fatalf("missing transparent exclusion facts: %v", wantExclusions)
	}
}

func TestInventoryTotalLimitErrorNamesTheTargetAndNextFile(t *testing.T) {
	err := inventoryTotalLimitError(
		filepath.FromSlash("/workspace/project"), filepath.FromSlash("/workspace/project/assets/archive.bin"),
		7*1024*1024*1024, 2*1024*1024*1024,
	)
	for _, expected := range []string{
		"assets/archive.bin", "7.0 GiB", "2.0 GiB", "8.0 GiB", "/workspace/project", ".prcignore",
	} {
		if !strings.Contains(filepath.ToSlash(err.Error()), expected) {
			t.Fatalf("actionable limit error missing %q: %v", expected, err)
		}
	}
}

func TestBuildAllowsOnlyReviewedNonSourceDirectoryExclusions(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "app.go", "package app\n")
	writeFile(t, root, "recordings/demo.mp4", "large generated recording\n")
	writeFile(t, root, "recordings/measurements.csv", "value\n1\n")
	writeFile(t, root, ".prcignore", "recordings | Local generated demo recordings are not project source.\n")

	item, err := Build(root)
	if err != nil {
		t.Fatal(err)
	}
	paths := make([]string, 0, len(item.Files))
	for _, record := range item.Files {
		paths = append(paths, record.Path)
	}
	if !slices.Equal(paths, []string{".prcignore", "app.go"}) {
		t.Fatalf("excluded non-source files were inventoried: %v", paths)
	}
	found := false
	for _, fact := range item.Facts {
		if fact.Key == "repository.user_exclusion" && fact.Source == "recordings" &&
			fact.Value == "Local generated demo recordings are not project source." && len(fact.Limitations) == 2 {
			found = true
		}
	}
	if !found {
		t.Fatalf("reviewed exclusion is not visible in inventory facts: %+v", item.Facts)
	}

	before := item.Digest
	writeFile(t, root, ".prcignore", "recordings | Local generated media files are not project source.\n")
	item, err = Build(root)
	if err != nil {
		t.Fatal(err)
	}
	if item.Digest == before {
		t.Fatal("inventory identity did not include the exclusion declaration")
	}
}

func TestBuildRefusesUserExclusionContainingReviewableFiles(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "archive/app.go", "package archive\n")
	writeFile(t, root, ".prcignore", "archive | This directory is believed to contain generated output.\n")
	if _, err := Build(root); err == nil || !strings.Contains(err.Error(), "archive/app.go") || !strings.Contains(err.Error(), "refuse") {
		t.Fatalf("source-bearing exclusion was accepted: %v", err)
	}
}

func TestBuildRefusesUserExclusionContainingLessCommonSourceOrConfiguration(t *testing.T) {
	for _, relative := range []string{"archive/component.vue", "archive/settings.json", "archive/service.proto"} {
		t.Run(filepath.Base(relative), func(t *testing.T) {
			root := t.TempDir()
			writeFile(t, root, relative, "reviewable project input\n")
			writeFile(t, root, ".prcignore", "archive | This directory is believed to contain generated output.\n")
			if _, err := Build(root); err == nil || !strings.Contains(err.Error(), relative) {
				t.Fatalf("reviewable exclusion %s was accepted: %v", relative, err)
			}
		})
	}
}

func TestBuildRejectsUnsafeOrHiddenUserExclusions(t *testing.T) {
	t.Run("path traversal", func(t *testing.T) {
		root := t.TempDir()
		writeFile(t, root, ".prcignore", "../outside | This path must never escape the selected project.\n")
		if _, err := Build(root); err == nil || !strings.Contains(err.Error(), "unsafe path") {
			t.Fatalf("unsafe path was accepted: %v", err)
		}
	})

	t.Run("symlinked ignore file", func(t *testing.T) {
		root := t.TempDir()
		outside := filepath.Join(t.TempDir(), "rules")
		if err := os.WriteFile(outside, []byte("media | Generated media is not project source.\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink(outside, filepath.Join(root, ".prcignore")); err != nil {
			t.Fatal(err)
		}
		if _, err := Build(root); err == nil || !strings.Contains(err.Error(), "regular file") {
			t.Fatalf("symlinked ignore file was accepted: %v", err)
		}
	})

	t.Run("hidden by automatic exclusion", func(t *testing.T) {
		root := t.TempDir()
		writeFile(t, root, "node_modules/cache/blob.bin", "generated\n")
		writeFile(t, root, ".prcignore", "node_modules/cache | Redundant generated dependency cache directory.\n")
		if _, err := Build(root); err == nil || !strings.Contains(err.Error(), "automatically excluded") {
			t.Fatalf("hidden exclusion was accepted: %v", err)
		}
	})

	t.Run("symlink inside excluded directory", func(t *testing.T) {
		root := t.TempDir()
		if err := os.Mkdir(filepath.Join(root, "recordings"), 0o700); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink(filepath.Join(t.TempDir(), "outside"), filepath.Join(root, "recordings", "latest")); err != nil {
			t.Fatal(err)
		}
		writeFile(t, root, ".prcignore", "recordings | Local generated recordings are not project source.\n")
		if _, err := Build(root); err == nil || !strings.Contains(err.Error(), "symbolic link") {
			t.Fatalf("an exclusion containing a symlink was accepted: %v", err)
		}
	})
}

func TestBuildRecordsCleanDirtyAndExternalGitWorktreeState(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("Git is not available")
	}
	repository := t.TempDir()
	runGit(t, repository, "init", "-q")
	runGit(t, repository, "config", "user.name", "PRC Test")
	runGit(t, repository, "config", "user.email", "prc-test@example.invalid")
	writeFile(t, repository, "app.go", "package app\n")
	runGit(t, repository, "add", "app.go")
	runGit(t, repository, "commit", "-q", "-m", "fixture")
	clean, err := Build(repository)
	if err != nil {
		t.Fatal(err)
	}
	if clean.GitCommit == "" || inventoryFactValue(clean.Facts, "repository.git_worktree_state") != "clean" {
		t.Fatalf("clean repository identity was not established: commit=%q facts=%+v", clean.GitCommit, clean.Facts)
	}
	candidateRoot := t.TempDir()
	writeFile(t, candidateRoot, "app.go", "package candidate\n")
	candidate, err := Build(candidateRoot)
	if err != nil {
		t.Fatal(err)
	}
	candidate, err = BindDerivedSource(candidate, clean)
	if err != nil {
		t.Fatal(err)
	}
	if candidate.GitCommit != clean.GitCommit || inventoryFactValue(candidate.Facts, "repository.git_worktree_state") != "dirty" {
		t.Fatalf("derived candidate lost source provenance: commit=%q facts=%+v", candidate.GitCommit, candidate.Facts)
	}
	if err := VerifyIdentity(candidate); err != nil {
		t.Fatalf("derived candidate identity is invalid: %v", err)
	}
	writeFile(t, repository, "app.go", "package changed\n")
	dirty, err := Build(repository)
	if err != nil {
		t.Fatal(err)
	}
	if dirty.GitCommit != clean.GitCommit || inventoryFactValue(dirty.Facts, "repository.git_worktree_state") != "dirty" {
		t.Fatalf("dirty repository was mislabeled: commit=%q facts=%+v", dirty.GitCommit, dirty.Facts)
	}
	runGit(t, repository, "restore", "app.go")
	worktree := filepath.Join(t.TempDir(), "external-worktree")
	runGit(t, repository, "worktree", "add", "--detach", worktree)
	worktreeInventory, err := Build(worktree)
	if err != nil {
		t.Fatal(err)
	}
	if worktreeInventory.GitCommit != clean.GitCommit || inventoryFactValue(worktreeInventory.Facts, "repository.git_worktree_state") != "clean" {
		t.Fatalf("standard external Git worktree was not supported: commit=%q facts=%+v", worktreeInventory.GitCommit, worktreeInventory.Facts)
	}
}

func TestBuildRejectsOversizedSparseFileBeforeHashing(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "oversized.bin")
	file, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
	if err := os.Truncate(path, maxInventoryFileBytes+1); err != nil {
		t.Skipf("filesystem cannot create a sparse limit fixture: %v", err)
	}
	if _, err := Build(root); err == nil || !strings.Contains(err.Error(), "inventory limit") {
		t.Fatalf("oversized file was accepted: %v", err)
	}
}

func runGit(t *testing.T, root string, arguments ...string) {
	t.Helper()
	command := exec.Command("git", append([]string{"-C", root}, arguments...)...)
	command.Env = append(os.Environ(), "GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL="+os.DevNull, "GIT_TERMINAL_PROMPT=0")
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("git %v failed: %v: %s", arguments, err, output)
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

func TestGitCommitReadsOnlyBoundedRepositoryMetadata(t *testing.T) {
	commit := strings.Repeat("a", 40)
	packedCommit := strings.Repeat("b", 64)
	tests := map[string]struct {
		prepare func(*testing.T, string)
		want    string
	}{
		"detached head": {
			prepare: func(t *testing.T, root string) {
				writeFile(t, root, ".git/HEAD", commit+"\n")
			},
			want: commit,
		},
		"loose branch": {
			prepare: func(t *testing.T, root string) {
				writeFile(t, root, ".git/HEAD", "ref: refs/heads/main\n")
				writeFile(t, root, ".git/refs/heads/main", commit+"\n")
			},
			want: commit,
		},
		"packed branch": {
			prepare: func(t *testing.T, root string) {
				writeFile(t, root, ".git/HEAD", "ref: refs/heads/main\n")
				writeFile(t, root, ".git/packed-refs", "# pack-refs with: peeled\n"+packedCommit+" refs/heads/main\n")
			},
			want: packedCommit,
		},
		"invalid detached value": {
			prepare: func(t *testing.T, root string) {
				writeFile(t, root, ".git/HEAD", "not-a-commit\n")
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			root := t.TempDir()
			test.prepare(t, root)
			if got := gitCommit(root); got != test.want {
				t.Fatalf("gitCommit() = %q, want %q", got, test.want)
			}
		})
	}
}

func TestGitCommitRejectsMetadataEscapes(t *testing.T) {
	commit := strings.Repeat("c", 40)
	t.Run("head ref traversal", func(t *testing.T) {
		parent := t.TempDir()
		root := filepath.Join(parent, "target")
		if err := os.Mkdir(root, 0o755); err != nil {
			t.Fatal(err)
		}
		writeFile(t, root, ".git/HEAD", "ref: ../../secret-ref\n")
		writeFile(t, parent, "secret-ref", commit+"\n")
		if got := gitCommit(root); got != "" {
			t.Fatalf("escaped ref disclosed %q", got)
		}
	})
	t.Run("external gitdir", func(t *testing.T) {
		parent := t.TempDir()
		root := filepath.Join(parent, "target")
		if err := os.Mkdir(root, 0o755); err != nil {
			t.Fatal(err)
		}
		writeFile(t, root, ".git", "gitdir: ../external-git\n")
		writeFile(t, parent, "external-git/HEAD", commit+"\n")
		if got := gitCommit(root); got != "" {
			t.Fatalf("external gitdir disclosed %q", got)
		}
	})
	t.Run("gitdir symlink", func(t *testing.T) {
		root := t.TempDir()
		external := t.TempDir()
		writeFile(t, external, "HEAD", commit+"\n")
		if err := os.Symlink(external, filepath.Join(root, ".git")); err != nil {
			t.Fatal(err)
		}
		if got := gitCommit(root); got != "" {
			t.Fatalf("symlinked gitdir disclosed %q", got)
		}
	})
}

func TestInspectFileBoundsStructuredCapture(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "document.yaml", "openapi: 3.2.0\n")
	path := filepath.Join(root, "document.yaml")
	record, data, err := inspectFile(path, "document.yaml", true, 4)
	if err != nil {
		t.Fatal(err)
	}
	if record.Size == 0 || data != nil {
		t.Fatalf("record=%+v captured=%q", record, data)
	}
	record, data, err = inspectFile(path, "document.yaml", true, maxStructuredFileBytes)
	if err != nil || string(data) != "openapi: 3.2.0\n" {
		t.Fatalf("record=%+v captured=%q err=%v", record, data, err)
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

func TestBindConfigurationAddsSourcedDeclaredScope(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "src/app.py", "print('ready')\n")
	baseline, err := Build(root)
	if err != nil {
		t.Fatal(err)
	}
	configPath, err := filepath.Abs(filepath.Join("..", "..", "fixtures", "config", "production-readiness.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	validation, err := projectconfig.Load(configPath)
	if err != nil {
		t.Fatal(err)
	}
	bound, err := BindConfiguration(baseline, validation, configPath)
	if err != nil {
		t.Fatal(err)
	}
	if bound.Digest == baseline.Digest || bound.DeclaredScope == nil ||
		bound.DeclaredScope.ConfigurationDigest != validation.Digest ||
		bound.DeclaredScope.ProjectID != "example-product" || !bound.DeclaredScope.Features["authentication"] ||
		bound.DeclaredScope.ArtifactDigests == nil {
		t.Fatalf("unexpected declared scope: %+v", bound.DeclaredScope)
	}
	foundFeature, foundComponent := false, false
	for _, fact := range bound.Facts {
		if fact.Key == "feature.authentication" && fact.Value == "true" &&
			fact.Detector == "prc.config.declaration" && len(fact.Limitations) > 0 {
			foundFeature = true
		}
	}
	for _, component := range bound.Components {
		if component.ID == "declared-service:src" {
			foundComponent = true
		}
	}
	if !foundFeature || !foundComponent {
		t.Fatalf("missing declaration provenance: facts=%+v components=%+v", bound.Facts, bound.Components)
	}
	second, err := BindConfiguration(baseline, validation, configPath)
	if err != nil || second.Digest != bound.Digest {
		t.Fatalf("configuration binding is not deterministic: %v %s != %s", err, second.Digest, bound.Digest)
	}
}

func TestBindConfigurationRejectsChangedInTargetSource(t *testing.T) {
	root := t.TempDir()
	source, err := os.ReadFile(filepath.Join("..", "..", "fixtures", "config", "production-readiness.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	configPath := filepath.Join(root, "production-readiness.yaml")
	writeFile(t, root, "production-readiness.yaml", string(source))
	validation, err := projectconfig.Load(configPath)
	if err != nil {
		t.Fatal(err)
	}
	writeFile(t, root, "production-readiness.yaml", strings.Replace(string(source), "Example Product", "Changed Product", 1))
	item, err := Build(root)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := BindConfiguration(item, validation, configPath); err == nil || !strings.Contains(err.Error(), "changed") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestBindConfigurationRejectsUnverifiableDeclaredSourceRef(t *testing.T) {
	root := t.TempDir()
	source, err := os.ReadFile(filepath.Join("..", "..", "fixtures", "config", "production-readiness.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	configured := strings.Replace(string(source), `  source_ref: ""`, `  source_ref: "`+strings.Repeat("a", 40)+`"`, 1)
	configPath := filepath.Join(root, "production-readiness.yaml")
	writeFile(t, root, "production-readiness.yaml", configured)
	validation, err := projectconfig.Load(configPath)
	if err != nil {
		t.Fatal(err)
	}
	item, err := Build(root)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := BindConfiguration(item, validation, configPath); err == nil || !strings.Contains(err.Error(), "no Git revision") {
		t.Fatalf("unexpected error: %v", err)
	}
}
