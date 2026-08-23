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

func TestLiveCheckovOCIUsesPinnedOfflinePolicyAndRejectsSuppressions(t *testing.T) {
	configuredImage := os.Getenv("PRC_TEST_CHECKOV_IMAGE")
	if configuredImage == "" {
		t.Skip("set PRC_TEST_CHECKOV_IMAGE to the reviewed image after pulling it by digest")
	}
	if configuredImage != CheckovImage {
		t.Fatalf("PRC_TEST_CHECKOV_IMAGE does not match the reviewed image: %s", configuredImage)
	}
	manifest, err := LoadManifest(filepath.Join("..", "..", "adapters", "checkov-v3.3.8.yaml"))
	if err != nil {
		t.Fatal(err)
	}

	risky := t.TempDir()
	writeSnapshotFixture(t, risky, "infra/main.tf", `resource "aws_s3_bucket" "public" {
  bucket = "example-public-bucket"
}

resource "aws_s3_bucket_public_access_block" "public" {
  bucket                  = aws_s3_bucket.public.id
  block_public_acls       = false
  block_public_policy     = false
  ignore_public_acls      = false
  restrict_public_buckets = false
}
`)
	writeSnapshotFixture(t, risky, ".checkov.yml", "framework: [secrets]\nskip-check: [CKV_AWS_54]\nquiet: true\n")
	first := runLiveCheckov(t, risky, manifest, strings.Repeat("e", 64))
	second := runLiveCheckov(t, risky, manifest, strings.Repeat("f", 64))
	if len(first.Transcript.Artifacts) != 1 || first.Transcript.Artifacts[0].Digest != second.Transcript.Artifacts[0].Digest ||
		!bytes.Equal(first.ArtifactPayloads[first.Transcript.Artifacts[0].Digest], second.ArtifactPayloads[second.Transcript.Artifacts[0].Digest]) {
		t.Fatal("normalized live Checkov artifact is not deterministic for the same sealed inventory")
	}
	foundExpectedPolicy := false
	for _, observation := range first.Transcript.Observations {
		if observation.Outcome == "found" && observation.Data["check_id"] == "CKV_AWS_54" &&
			observation.Locations[0].Path == "infra/main.tf" {
			foundExpectedPolicy = true
		}
	}
	if !foundExpectedPolicy {
		t.Fatalf("target Checkov configuration changed the reviewed policy result: %+v", first.Transcript.Observations)
	}

	suppressed := t.TempDir()
	writeSnapshotFixture(t, suppressed, "deploy/app.yaml", `apiVersion: apps/v1
kind: Deployment
metadata:
  name: suppressed-risk
  annotations:
    checkov.io/skip1: CKV_K8S_21=target-owned suppression
spec:
  selector:
    matchLabels: {app: suppressed-risk}
  template:
    metadata:
      labels: {app: suppressed-risk}
    spec:
      containers:
        - name: app
          image: example.invalid/app:latest
`)
	if _, err := runLiveCheckovResult(suppressed, manifest, strings.Repeat("a", 64)); err == nil ||
		!strings.Contains(err.Error(), "suppressed policy checks") {
		t.Fatalf("inline Checkov suppression did not fail closed: %v", err)
	}
}

func runLiveCheckov(t *testing.T, root string, manifest Manifest, runID string) RunOutput {
	t.Helper()
	output, err := runLiveCheckovResult(root, manifest, runID)
	if err != nil {
		t.Fatal(err)
	}
	return output
}

func runLiveCheckovResult(root string, manifest Manifest, runID string) (RunOutput, error) {
	item, err := inventory.Build(root)
	if err != nil {
		return RunOutput{}, err
	}
	snapshot, err := PrepareSnapshotForManifest(item, manifest)
	if err != nil {
		return RunOutput{}, err
	}
	defer snapshot.Close()
	subject := Subject{TargetName: item.TargetName, TargetCommit: item.GitCommit, InventoryDigest: item.Digest}
	input, err := ExecutionInput(manifest, runID, subject, nil, nil)
	if err != nil {
		return RunOutput{}, err
	}
	plan, err := BuildSnapshotOCIPlan("docker", snapshot, runID, manifest)
	if err != nil {
		return RunOutput{}, err
	}
	output, err := RunOCI(context.Background(), plan, manifest, input)
	if err != nil {
		return RunOutput{}, err
	}
	if _, err := BindExecution(runID, subject, manifest, output); err != nil {
		return RunOutput{}, err
	}
	return output, nil
}
