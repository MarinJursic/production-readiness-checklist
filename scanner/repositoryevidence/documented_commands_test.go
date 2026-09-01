package repositoryevidence

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/MarinJursic/production-readiness-checklist/scanner/controlprogramcatalog"
	"github.com/MarinJursic/production-readiness-checklist/scanner/controlruntime"
	workspaceinventory "github.com/MarinJursic/production-readiness-checklist/scanner/inventory"
)

func writeFixture(t *testing.T, root, relative, content string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(relative))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}

func documentedCommandsTemplate(t *testing.T) controlprogramcatalog.Template {
	t.Helper()
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	catalog, err := controlprogramcatalog.Load(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, template := range catalog.Templates() {
		if template.CollectorContract.CollectorID == DocumentedCommandsCollectorID {
			return template
		}
	}
	t.Fatal("documented-command template is missing")
	return controlprogramcatalog.Template{}
}

func evaluateFixture(t *testing.T, readme string) controlruntime.Execution {
	t.Helper()
	root := t.TempDir()
	writeFixture(t, root, "package.json", `{"scripts":{"build":"node build.mjs","test":"node --test"}}`)
	writeFixture(t, root, "README.md", readme)
	item, err := workspaceinventory.Build(root)
	if err != nil {
		t.Fatal(err)
	}
	provider, err := NewDocumentedCommandsProvider(item)
	if err != nil {
		t.Fatal(err)
	}
	registry, err := controlruntime.NewRegistry(provider)
	if err != nil {
		t.Fatal(err)
	}
	template := documentedCommandsTemplate(t)
	binding, ok := Binding(item, template)
	if !ok {
		t.Fatal("documented-command binding was not created")
	}
	now := time.Date(2026, 8, 29, 12, 0, 0, 0, time.UTC)
	return controlruntime.Evaluate(context.Background(), template, binding, registry, now)
}

func TestDocumentedCommandsProviderPassesOnlyCodeDocumentation(t *testing.T) {
	result := evaluateFixture(t, "# Setup\n\n```bash\nnpm run build\nnpm test\n```\n")
	if result.Status != controlruntime.StatusPassed || result.EvidenceSHA256 == "" {
		t.Fatalf("documented command result = %+v", result)
	}
}

func TestDocumentedCommandsProviderBlocksMissingOrProseOnlyDocumentation(t *testing.T) {
	for _, readme := range []string{
		"# Setup\n\nThe npm run build and npm test commands are useful.\n",
		"# Setup\n\nRun `npm test`.\n",
	} {
		result := evaluateFixture(t, readme)
		if result.Status != controlruntime.StatusBlockedEvidence || result.Outcome != "blocked" {
			t.Fatalf("incomplete documentation result = %+v", result)
		}
	}
}

func TestPackageScriptsRejectsDuplicateKeysAndNonStringScripts(t *testing.T) {
	for _, document := range []string{
		`{"scripts":{"build":"a","build":"b","test":"c"}}`,
		`{"scripts":{"build":["a"],"test":"c"}}`,
		`{"scripts":{},"scripts":{}}`,
	} {
		if _, err := packageScripts([]byte(document)); err == nil {
			t.Fatalf("unsafe package JSON was accepted: %s", document)
		}
	}
}

func TestDocumentedCommandsProviderRejectsChangedInventoriedBytes(t *testing.T) {
	root := t.TempDir()
	writeFixture(t, root, "package.json", `{"scripts":{"build":"node build.mjs","test":"node --test"}}`)
	writeFixture(t, root, "README.md", "```sh\nnpm run build\nnpm test\n```\n")
	item, err := workspaceinventory.Build(root)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "README.md"), []byte(strings.Repeat("x", 10)), 0o600); err != nil {
		t.Fatal(err)
	}
	provider, _ := NewDocumentedCommandsProvider(item)
	registry, _ := controlruntime.NewRegistry(provider)
	template := documentedCommandsTemplate(t)
	binding, _ := Binding(item, template)
	result := controlruntime.Evaluate(context.Background(), template, binding, registry, time.Now().UTC())
	if result.Status != controlruntime.StatusBlockedCollection {
		t.Fatalf("changed documentation status = %s", result.Status)
	}
}

func TestDocumentedCommandsProviderIsDeterministicAndNeverExecutesScripts(t *testing.T) {
	root := t.TempDir()
	marker := filepath.Join(root, "must-not-exist")
	writeFixture(t, root, "package.json", `{"scripts":{"build":"touch must-not-exist","test":"touch must-not-exist"}}`)
	writeFixture(t, root, "README.md", "```sh\nnpm run build\nnpm test\n```\n")
	item, err := workspaceinventory.Build(root)
	if err != nil {
		t.Fatal(err)
	}
	provider, _ := NewDocumentedCommandsProvider(item)
	registry, _ := controlruntime.NewRegistry(provider)
	template := documentedCommandsTemplate(t)
	binding, _ := Binding(item, template)
	now := time.Date(2026, 8, 29, 12, 0, 0, 0, time.UTC)
	first := controlruntime.Evaluate(context.Background(), template, binding, registry, now)
	second := controlruntime.Evaluate(context.Background(), template, binding, registry, now)
	firstJSON, _ := json.Marshal(first)
	secondJSON, _ := json.Marshal(second)
	if string(firstJSON) != string(secondJSON) || first.EvidenceSHA256 != second.EvidenceSHA256 {
		t.Fatalf("repeated exact collection differs:\n%s\n%s", firstJSON, secondJSON)
	}
	if _, err := os.Lstat(marker); !os.IsNotExist(err) {
		t.Fatal("repository script was executed")
	}
}

func TestDocumentedCommandsProviderBlocksOversizedDocumentation(t *testing.T) {
	root := t.TempDir()
	writeFixture(t, root, "package.json", `{"scripts":{"build":"node build.mjs","test":"node --test"}}`)
	writeFixture(t, root, "README.md", strings.Repeat("x", maximumDocumentationFileBytes+1))
	item, err := workspaceinventory.Build(root)
	if err != nil {
		t.Fatal(err)
	}
	provider, _ := NewDocumentedCommandsProvider(item)
	registry, _ := controlruntime.NewRegistry(provider)
	template := documentedCommandsTemplate(t)
	binding, _ := Binding(item, template)
	result := controlruntime.Evaluate(context.Background(), template, binding, registry, time.Now().UTC())
	if result.Status != controlruntime.StatusBlockedEvidence {
		t.Fatalf("oversized documentation status = %s", result.Status)
	}
}

func TestDocumentedCommandsProviderCanProveEarlyEvidenceInALargeDocumentationRepository(t *testing.T) {
	root := t.TempDir()
	writeFixture(t, root, "package.json", `{"scripts":{"build":"node build.mjs","test":"node --test"}}`)
	writeFixture(t, root, "README.md", "```sh\nnpm run build\nnpm test\n```\n")
	for index := 0; index < maximumDocumentationFiles+20; index++ {
		writeFixture(t, root, fmt.Sprintf("docs/archive/%04d.md", index), "Archived unrelated note.\n")
	}
	item, err := workspaceinventory.Build(root)
	if err != nil {
		t.Fatal(err)
	}
	provider, _ := NewDocumentedCommandsProvider(item)
	registry, _ := controlruntime.NewRegistry(provider)
	template := documentedCommandsTemplate(t)
	binding, _ := Binding(item, template)
	result := controlruntime.Evaluate(context.Background(), template, binding, registry, time.Now().UTC())
	if result.Status != controlruntime.StatusPassed {
		t.Fatalf("early positive evidence was blocked by unrelated documentation: %+v", result)
	}
}
