package inventory

import (
	"os"
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
	if len(first.Facts) != 7 {
		t.Fatalf("expected seven provenance facts, got %+v", first.Facts)
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
