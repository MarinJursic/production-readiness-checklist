package adapter

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/MarinJursic/production-readiness-checklist/scanner/inventory"
)

func TestLiveGrypeOCIUsesOfflineBoundDatabaseAndIgnoresTargetPolicy(t *testing.T) {
	image := os.Getenv("PRC_TEST_GRYPE_IMAGE")
	database := os.Getenv("PRC_TEST_GRYPE_DB")
	if image == "" || database == "" {
		t.Skip("set PRC_TEST_GRYPE_IMAGE and PRC_TEST_GRYPE_DB to exercise the reviewed image and offline database")
	}
	manifest, err := LoadManifest(filepath.Join("..", "..", "adapters", "grype-v0.116.1.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	if image != manifest.Image {
		t.Fatalf("test image %q does not match manifest %q", image, manifest.Image)
	}
	root := t.TempDir()
	writeSnapshotFixture(t, root, "package-lock.json", `{
  "name": "vulnerable-fixture",
  "version": "1.0.0",
  "lockfileVersion": 3,
  "requires": true,
  "packages": {
    "": {"name": "vulnerable-fixture", "version": "1.0.0", "dependencies": {"lodash": "4.17.15"}},
    "node_modules/lodash": {"version": "4.17.15"}
  },
  "dependencies": {"lodash": {"version": "4.17.15"}}
}
`)
	writeSnapshotFixture(t, root, ".grype.yaml", "ignore:\n  - vulnerability: CVE-2021-23337\n")

	first := runLiveGrype(t, root, database, manifest)
	second := runLiveGrype(t, root, database, manifest)
	if len(first.Transcript.Artifacts) != 1 || len(second.Transcript.Artifacts) != 1 ||
		first.Transcript.Artifacts[0].Digest != second.Transcript.Artifacts[0].Digest ||
		!bytes.Equal(first.ArtifactPayloads[first.Transcript.Artifacts[0].Digest], second.ArtifactPayloads[second.Transcript.Artifacts[0].Digest]) {
		t.Fatalf("Grype output was not deterministic: first=%+v second=%+v", first.Transcript.Artifacts, second.Transcript.Artifacts)
	}
	foundCVE := false
	for _, observation := range first.Transcript.Observations {
		aliases, _ := observation.Data["aliases"].([]string)
		for _, alias := range aliases {
			if alias == "CVE-2021-23337" {
				foundCVE = true
			}
		}
	}
	if !foundCVE {
		t.Fatalf("hostile target policy suppressed the expected CVE: %+v", first.Transcript.Observations)
	}
	if len(first.DataInputs) != 1 || first.DataInputs[0].Name != GrypeDataMountName || first.DataInputs[0].Bytes < 1024*1024 {
		t.Fatalf("offline database identity was not bound: %+v", first.DataInputs)
	}
}

func runLiveGrype(t *testing.T, root, database string, manifest Manifest) RunOutput {
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
	config, err := os.ReadFile(filepath.Join(snapshot.Path, filepath.FromSlash(GrypeConfigSnapshotPath)))
	if err != nil || !bytes.Equal(config, grypeConfig) {
		t.Fatalf("scanner-owned Grype configuration was not injected: %v", err)
	}
	runID := item.Digest
	plan, err := BuildSnapshotOCIPlanWithData("docker", snapshot, runID, manifest, map[string]string{
		GrypeDataMountName: database,
	})
	if err != nil {
		t.Fatal(err)
	}
	input, err := ExecutionInput(manifest, runID, Subject{
		TargetName: item.TargetName, InventoryDigest: item.Digest,
	}, map[string]any{}, map[string]any{})
	if err != nil {
		t.Fatal(err)
	}
	output, err := RunOCI(context.Background(), plan, manifest, input)
	if err != nil {
		t.Fatal(err)
	}
	return output
}
