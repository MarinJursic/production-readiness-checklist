package repositoryevidence

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/MarinJursic/production-readiness-checklist/scanner/controlprogramcatalog"
	"github.com/MarinJursic/production-readiness-checklist/scanner/controlruntime"
	workspaceinventory "github.com/MarinJursic/production-readiness-checklist/scanner/inventory"
)

func documentedTopicsTemplate(t *testing.T, collectorID string) controlprogramcatalog.Template {
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
		if template.CollectorContract.CollectorID == collectorID {
			return template
		}
	}
	t.Fatalf("documented-topic template %s is missing", collectorID)
	return controlprogramcatalog.Template{}
}

func evaluateTopicFixture(t *testing.T, collectorID, markdown string) controlruntime.Execution {
	t.Helper()
	root := t.TempDir()
	writeFixture(t, root, "docs/operations.md", markdown)
	item, err := workspaceinventory.Build(root)
	if err != nil {
		t.Fatal(err)
	}
	providers, err := NewDocumentedTopicsProviders(item)
	if err != nil {
		t.Fatal(err)
	}
	var selected *DocumentedTopicsProvider
	for _, provider := range providers {
		if provider.ID() == collectorID {
			selected = provider
		}
	}
	if selected == nil {
		t.Fatalf("provider %s is missing", collectorID)
	}
	registry, err := controlruntime.NewRegistry(selected)
	if err != nil {
		t.Fatal(err)
	}
	template := documentedTopicsTemplate(t, collectorID)
	binding, ok := Binding(item, template)
	if !ok {
		t.Fatal("documented-topic binding was not created")
	}
	return controlruntime.Evaluate(
		context.Background(), template, binding, registry,
		time.Date(2026, 9, 1, 12, 0, 0, 0, time.UTC),
	)
}

func TestDocumentedArchitectureTopicsPassOnlyExactHeadingsWithContent(t *testing.T) {
	markdown := `# Architecture
The service boundaries and responsibility split are recorded here.

## Data flow
Requests enter through the API and validated events leave through the queue.

# Dependencies
Direct runtime and service dependencies have named owners and versions.

# Deployment
The release pipeline deploys the sealed artifact through staged environments.

# Recovery
Restore, rollback, and validation steps are recorded for the service owner.
`
	result := evaluateTopicFixture(t, ArchitectureTopicsCollectorID, markdown)
	if result.Status != controlruntime.StatusPassed || result.EvidenceSHA256 == "" {
		t.Fatalf("documented architecture topics = %+v", result)
	}
}

func TestDocumentedConventionTopicsPassWithoutPrescribingFilesOrProjectLayout(t *testing.T) {
	markdown := `# Coding
Changes follow the repository style checks and keep ownership boundaries clear.

# Review
Another maintainer reviews risk, behavior, and evidence before merge.

# Testing
The documented test command covers changed behavior before a release.

# Release
The maintainer verifies the built artifact and recorded checks before publishing.
`
	result := evaluateTopicFixture(t, ConventionsTopicsCollectorID, markdown)
	if result.Status != controlruntime.StatusPassed {
		t.Fatalf("documented conventions = %+v", result)
	}
}

func TestDocumentedTopicsBlockPlaceholdersAliasesAndFencedHeadings(t *testing.T) {
	for name, markdown := range map[string]string{
		"placeholders": `# Architecture
TBD
# Data flow
Coming soon
# Dependencies
To be documented
# Deployment
Not yet documented
# Recovery
To be determined
`,
		"aliases": `# System design
This describes architecture without the explicit required heading contract.
# Information movement
This describes data flow without the exact required topic identity.
`,
		"fenced": "```markdown\n# Architecture\nThis heading is example data, not project documentation.\n```\n",
		"comment-only": `# Architecture
<!-- TODO: write the actual architecture section before release. -->
# Data flow
<!-- This is not visible project documentation. -->
`,
		"fenced-content": "# Architecture\n```text\nThis code example is not enough to document the project architecture.\n```\n",
	} {
		t.Run(name, func(t *testing.T) {
			result := evaluateTopicFixture(t, ArchitectureTopicsCollectorID, markdown)
			if result.Status != controlruntime.StatusBlockedEvidence || result.Outcome != "blocked" {
				t.Fatalf("unsupported documentation created a verdict: %+v", result)
			}
		})
	}
}

func TestDocumentedTopicsRejectChangedInventoriedBytes(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "README.md")
	writeFixture(t, root, "README.md", "# Coding\nA complete coding convention is recorded for this project.\n")
	item, err := workspaceinventory.Build(root)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(strings.Repeat("changed", 20)), 0o600); err != nil {
		t.Fatal(err)
	}
	providers, _ := NewDocumentedTopicsProviders(item)
	registry, _ := controlruntime.NewRegistry(providers[1])
	template := documentedTopicsTemplate(t, ConventionsTopicsCollectorID)
	binding, _ := Binding(item, template)
	result := controlruntime.Evaluate(context.Background(), template, binding, registry, time.Now().UTC())
	if result.Status != controlruntime.StatusBlockedCollection {
		t.Fatalf("changed documentation status = %s", result.Status)
	}
}

func TestDocumentedTopicsBlockInvalidUTF8Documentation(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "README.md"), []byte{'#', ' ', 'C', 'o', 'd', 'i', 'n', 'g', '\n', 0xff, '\n'}, 0o600); err != nil {
		t.Fatal(err)
	}
	item, err := workspaceinventory.Build(root)
	if err != nil {
		t.Fatal(err)
	}
	providers, _ := NewDocumentedTopicsProviders(item)
	registry, _ := controlruntime.NewRegistry(providers[1])
	template := documentedTopicsTemplate(t, ConventionsTopicsCollectorID)
	binding, _ := Binding(item, template)
	result := controlruntime.Evaluate(context.Background(), template, binding, registry, time.Now().UTC())
	if result.Status != controlruntime.StatusBlockedEvidence {
		t.Fatalf("invalid UTF-8 documentation status = %s", result.Status)
	}
}

func TestDocumentedTopicsParserIgnoresHeadingLikeTextInsideCodeFences(t *testing.T) {
	data := []byte("# Architecture\nReal architecture content names boundaries and responsibilities.\n\n```md\n# Recovery\nFake recovery content inside a code example.\n```\n")
	sections := exactMarkdownTopicSections(data, architectureTopics)
	if sections["architecture"] == "" || sections["recovery"] != "" {
		t.Fatalf("fenced heading handling = %+v", sections)
	}
}

func TestDocumentedTopicsParserRequiresAValidMatchingFenceClosure(t *testing.T) {
	data := []byte("````markdown\n```\n# Architecture\nFake content after a shorter fence.\n``` language\n# Recovery\nFake content after a closing fence with an info string.\n~~~~\n# Dependencies\nFake content after a different fence marker.\n````\n# Deployment\nReal deployment documentation appears outside the code block.\n")
	sections := exactMarkdownTopicSections(data, architectureTopics)
	if sections["architecture"] != "" || sections["recovery"] != "" || sections["dependencies"] != "" || sections["deployment"] == "" {
		t.Fatalf("strict fenced heading handling = %+v", sections)
	}
}

func TestDocumentedTopicsCanProveEvidenceAfterAnOversizedUnrelatedDocument(t *testing.T) {
	root := t.TempDir()
	writeFixture(t, root, "A-archive.md", strings.Repeat("x", maximumDocumentationFileBytes+1))
	writeFixture(t, root, "README.md", `# Coding
Changes follow the repository style checks and keep ownership boundaries clear.
# Review
Another maintainer reviews risk, behavior, and evidence before merge.
# Testing
The documented test command covers changed behavior before a release.
# Release
The maintainer verifies the built artifact and recorded checks before publishing.
`)
	item, err := workspaceinventory.Build(root)
	if err != nil {
		t.Fatal(err)
	}
	providers, _ := NewDocumentedTopicsProviders(item)
	registry, _ := controlruntime.NewRegistry(providers[1])
	template := documentedTopicsTemplate(t, ConventionsTopicsCollectorID)
	binding, _ := Binding(item, template)
	result := controlruntime.Evaluate(context.Background(), template, binding, registry, time.Now().UTC())
	if result.Status != controlruntime.StatusPassed {
		t.Fatalf("valid topic evidence was blocked by an unrelated oversized document: %+v", result)
	}
}
