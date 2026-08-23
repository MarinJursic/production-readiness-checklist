package adapter

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/MarinJursic/production-readiness-checklist/scanner/inventory"
)

func TestLiveGitleaksOCIUsesPinnedRulesAndIgnoresTargetSuppressions(t *testing.T) {
	configuredImage := os.Getenv("PRC_TEST_GITLEAKS_IMAGE")
	if configuredImage == "" {
		t.Skip("set PRC_TEST_GITLEAKS_IMAGE to the reviewed image after pulling it by digest")
	}
	if configuredImage != GitleaksImage {
		t.Fatalf("PRC_TEST_GITLEAKS_IMAGE does not match the reviewed image: %s", configuredImage)
	}
	manifest, err := LoadManifest(filepath.Join("..", "..", "adapters", "gitleaks-v8.30.0.yaml"))
	if err != nil {
		t.Fatal(err)
	}

	clean := t.TempDir()
	writeSnapshotFixture(t, clean, "README.md", "No credentials are stored here.\n")
	cleanTranscript := runLiveGitleaks(t, clean, manifest)
	if len(cleanTranscript.Observations) != 1 || cleanTranscript.Observations[0].Outcome != "not_found" {
		t.Fatalf("clean Gitleaks transcript = %+v", cleanTranscript)
	}

	leaky := t.TempDir()
	synthetic := "ghp_" + "aB3dE5fG7hJ9kL1mN3pQ5rS7tU9vW1xY3zA5"
	writeSnapshotFixture(t, leaky, "leak.env", "token = \""+synthetic+"\" # gitleaks:allow\n")
	writeSnapshotFixture(t, leaky, ".gitleaks.toml", "title = \"target suppression\"\n")
	writeSnapshotFixture(t, leaky, ".gitleaksignore", "/workspace/leak.env:github-pat:1\n")
	leakyTranscript := runLiveGitleaks(t, leaky, manifest)
	if len(leakyTranscript.Observations) != 1 || leakyTranscript.Observations[0].Outcome != "found" ||
		leakyTranscript.Observations[0].Locations[0].Path != "leak.env" {
		t.Fatalf("leaky Gitleaks transcript = %+v", leakyTranscript)
	}
	encoded, err := json.Marshal(leakyTranscript)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(encoded), synthetic) || strings.Contains(string(encoded), "REDACTED") {
		t.Fatal("bound transcript retained matched secret material or raw report text")
	}
}

func runLiveGitleaks(t *testing.T, root string, manifest Manifest) Transcript {
	t.Helper()
	item, err := inventory.Build(root)
	if err != nil {
		t.Fatal(err)
	}
	snapshot, err := PrepareSnapshotForManifest(item, manifest)
	if err != nil {
		t.Fatal(err)
	}
	defer snapshot.Close()
	runID := strings.Repeat("d", 64)
	subject := Subject{TargetName: item.TargetName, TargetCommit: item.GitCommit, InventoryDigest: item.Digest}
	input, err := ExecutionInput(manifest, runID, subject, map[string]any{}, map[string]any{})
	if err != nil {
		t.Fatal(err)
	}
	plan, err := BuildSnapshotOCIPlan("docker", snapshot, runID, manifest)
	if err != nil {
		t.Fatal(err)
	}
	output, err := RunOCI(context.Background(), plan, manifest, input)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := BindExecution(runID, subject, manifest, output); err != nil {
		t.Fatal(err)
	}
	return output.Transcript
}
