package adapter

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/MarinJursic/production-readiness-checklist/scanner/inventory"
)

func TestLiveSyftOCIProducesDeterministicCycloneDXDespiteTargetConfig(t *testing.T) {
	configuredImage := os.Getenv("PRC_TEST_SYFT_IMAGE")
	if configuredImage == "" {
		t.Skip("set PRC_TEST_SYFT_IMAGE to the reviewed image after pulling it by digest")
	}
	if configuredImage != SyftImage {
		t.Fatalf("PRC_TEST_SYFT_IMAGE does not match the reviewed image: %s", configuredImage)
	}
	manifest, err := LoadManifest(filepath.Join("..", "..", "adapters", "syft-v1.51.0.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	root := t.TempDir()
	writeSnapshotFixture(t, root, "go.mod", "module example.com/fixture\n\ngo 1.27\n\nrequire github.com/google/uuid v1.6.0\n")
	writeSnapshotFixture(t, root, ".syft.yaml", "output: table\nparallelism: 99\n")

	first := runLiveSyft(t, root, manifest, strings.Repeat("c", 64))
	second := runLiveSyft(t, root, manifest, strings.Repeat("d", 64))
	firstArtifact := first.Transcript.Artifacts[0]
	secondArtifact := second.Transcript.Artifacts[0]
	if firstArtifact.Digest != secondArtifact.Digest ||
		!bytes.Equal(first.ArtifactPayloads[firstArtifact.Digest], second.ArtifactPayloads[secondArtifact.Digest]) {
		t.Fatal("normalized live Syft artifact is not deterministic for the same sealed inventory")
	}
	observation := first.Transcript.Observations[0]
	if observation.Outcome != "value" || observation.Data["package_component_count"].(int) < 1 ||
		!bytes.Contains(first.ArtifactPayloads[firstArtifact.Digest], []byte("pkg:golang/github.com/google/uuid@v1.6.0")) {
		t.Fatalf("live Syft transcript = %+v", first.Transcript)
	}
}

func runLiveSyft(t *testing.T, root string, manifest Manifest, runID string) RunOutput {
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
	subject := Subject{TargetName: item.TargetName, TargetCommit: item.GitCommit, InventoryDigest: item.Digest}
	input, err := ExecutionInput(manifest, runID, subject, nil, nil)
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
	return output
}
