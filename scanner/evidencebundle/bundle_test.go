package evidencebundle_test

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/MarinJursic/production-readiness-checklist/scanner/catalog"
	"github.com/MarinJursic/production-readiness-checklist/scanner/controlprogram"
	"github.com/MarinJursic/production-readiness-checklist/scanner/controlprogramcatalog"
	"github.com/MarinJursic/production-readiness-checklist/scanner/controlruntime"
	"github.com/MarinJursic/production-readiness-checklist/scanner/engine"
	"github.com/MarinJursic/production-readiness-checklist/scanner/evidence"
	"github.com/MarinJursic/production-readiness-checklist/scanner/evidencebundle"
	"github.com/MarinJursic/production-readiness-checklist/scanner/fullscan"
	workspaceinventory "github.com/MarinJursic/production-readiness-checklist/scanner/inventory"
	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
	"github.com/MarinJursic/production-readiness-checklist/scanner/repositoryevidence"
	"github.com/MarinJursic/production-readiness-checklist/scanner/state"
	"github.com/MarinJursic/production-readiness-checklist/scanner/trust"
)

type signedBundleFixture struct {
	catalog                                         *controlprogramcatalog.Catalog
	item                                            model.Inventory
	bundlePath, storePath, policyPath, evidencePath string
	verifiedAt                                      time.Time
}

func TestVerifiedBundleEvaluatesExactProgramAndRetainsBothSignatures(t *testing.T) {
	fixture := newSignedBundleFixture(t)
	executions, verification, err := evidencebundle.VerifyAndEvaluate(
		fixture.catalog, fixture.item, fixture.bundlePath, fixture.storePath,
		fixture.policyPath, fixture.evidencePath, fixture.verifiedAt,
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(executions) != 1 || executions[0].Status != controlruntime.StatusPassed || !executions[0].Authenticated() {
		t.Fatalf("authenticated executions = %+v", executions)
	}
	if verification.SchemaVersion != evidencebundle.VerificationSchema || verification.EntryCount != 1 ||
		len(verification.Entries) != 1 || verification.Entries[0].TemplateID == "" ||
		!verification.PolicySignature.Verified || !verification.EvidenceSignature.Verified ||
		verification.PolicySHA256 == verification.BundleSHA256 ||
		verification.PolicySignature.SHA256 != verification.PolicySHA256 ||
		verification.EvidenceSignature.SHA256 != verification.BundleSHA256 ||
		verification.PolicySignature.KeyID == verification.EvidenceSignature.KeyID {
		t.Fatalf("bundle verification = %+v", verification)
	}
}

func TestVerifiedBundleRejectsTamperingAndSharedTrustKey(t *testing.T) {
	t.Run("tampered bundle", func(t *testing.T) {
		fixture := newSignedBundleFixture(t)
		data, err := os.ReadFile(fixture.bundlePath)
		if err != nil {
			t.Fatal(err)
		}
		data[len(data)-2] ^= 1
		if err := os.WriteFile(fixture.bundlePath, data, 0o600); err != nil {
			t.Fatal(err)
		}
		if _, _, err := evidencebundle.VerifyAndEvaluate(
			fixture.catalog, fixture.item, fixture.bundlePath, fixture.storePath,
			fixture.policyPath, fixture.evidencePath, fixture.verifiedAt,
		); err == nil {
			t.Fatal("tampered evidence bundle was accepted")
		}
	})

	t.Run("same key", func(t *testing.T) {
		fixture := newSignedBundleFixture(t)
		policyData, err := os.ReadFile(fixture.policyPath)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(fixture.evidencePath, policyData, 0o600); err != nil {
			t.Fatal(err)
		}
		if _, _, err := evidencebundle.VerifyAndEvaluate(
			fixture.catalog, fixture.item, fixture.bundlePath, fixture.storePath,
			fixture.policyPath, fixture.evidencePath, fixture.verifiedAt,
		); err == nil {
			t.Fatal("one key was accepted in both trust roles")
		}
	})
}

func TestVerifiedBundleIdentityIgnoresJSONPresentationWhitespace(t *testing.T) {
	fixture := newSignedBundleFixture(t)
	data, err := os.ReadFile(fixture.bundlePath)
	if err != nil {
		t.Fatal(err)
	}
	var pretty bytes.Buffer
	if err := json.Indent(&pretty, data, "", "  "); err != nil {
		t.Fatal(err)
	}
	pretty.WriteByte('\n')
	if err := os.WriteFile(fixture.bundlePath, pretty.Bytes(), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, _, err := evidencebundle.VerifyAndEvaluate(
		fixture.catalog, fixture.item, fixture.bundlePath, fixture.storePath,
		fixture.policyPath, fixture.evidencePath, fixture.verifiedAt,
	); err != nil {
		t.Fatalf("presentation-only JSON change altered the typed bundle identity: %v", err)
	}
}

func TestVerifiedBundleSurvivesFullRunHistoryReplay(t *testing.T) {
	fixture := newSignedBundleFixture(t)
	catalogRoot, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	assertionCatalog, err := catalog.Load(catalogRoot)
	if err != nil {
		t.Fatal(err)
	}
	run, err := engine.New(assertionCatalog).Scan("prc/core-repository", fixture.item)
	if err != nil {
		t.Fatal(err)
	}
	executions, verification, err := evidencebundle.VerifyAndEvaluate(
		fixture.catalog, fixture.item, fixture.bundlePath, fixture.storePath,
		fixture.policyPath, fixture.evidencePath, run.CompletedAt,
	)
	if err != nil {
		t.Fatal(err)
	}
	run.AuthoritativeEvidence = append(run.AuthoritativeEvidence, model.AuthoritativeEvidenceVerification{
		SchemaVersion: verification.SchemaVersion, BundleID: verification.BundleID,
		BundleSHA256: verification.BundleSHA256, PolicySHA256: verification.PolicySHA256,
		CatalogSHA256: verification.CatalogSHA256, InventorySHA256: verification.InventorySHA256,
		Authority: verification.Authority, EntryCount: verification.EntryCount,
		Entries:         append([]model.AuthoritativeEvidenceEntry(nil), verification.Entries...),
		PolicySignature: verification.PolicySignature, EvidenceSignature: verification.EvidenceSignature,
	})
	run, err = fullscan.AttachProgramExecutions(catalogRoot, assertionCatalog, run, executions)
	if err != nil {
		t.Fatal(err)
	}
	historyRoot := t.TempDir()
	if err := os.Chmod(historyRoot, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := evidence.WriteRun(historyRoot, run); err != nil {
		t.Fatal(err)
	}
	store, err := state.Open(context.Background(), historyRoot)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	if err := store.IndexRun(context.Background(), run); err != nil {
		t.Fatal(err)
	}
	loaded, err := store.LoadRun(context.Background(), run.RunID)
	if err != nil || len(loaded.AuthoritativeEvidence) != 1 ||
		len(loaded.AuthoritativeEvidence[0].Entries) != 1 {
		t.Fatalf("replayed authoritative history = %+v err=%v", loaded.AuthoritativeEvidence, err)
	}
}

func newSignedBundleFixture(t *testing.T) signedBundleFixture {
	t.Helper()
	catalogRoot, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	catalog, err := controlprogramcatalog.Load(catalogRoot)
	if err != nil {
		t.Fatal(err)
	}
	var template controlprogramcatalog.Template
	for _, candidate := range catalog.Templates() {
		if candidate.CollectorContract.CollectorID == repositoryevidence.DocumentedCommandsCollectorID {
			template = candidate
			break
		}
	}
	if template.TemplateID == "" {
		t.Fatal("documented commands template is missing")
	}

	workspace := t.TempDir()
	writeFile(t, filepath.Join(workspace, "package.json"), `{"scripts":{"build":"node build.mjs","test":"node --test"}}`)
	writeFile(t, filepath.Join(workspace, "README.md"), "```sh\nnpm run build\nnpm test\n```\n")
	item, err := workspaceinventory.Build(workspace)
	if err != nil {
		t.Fatal(err)
	}
	binding, ok := repositoryevidence.Binding(item, template)
	if !ok {
		t.Fatal("documented commands binding was not created")
	}
	program, err := template.Program(binding)
	if err != nil {
		t.Fatal(err)
	}
	providerValue, err := repositoryevidence.NewDocumentedCommandsProvider(item)
	if err != nil {
		t.Fatal(err)
	}
	registry, err := controlruntime.NewRegistry(providerValue)
	if err != nil {
		t.Fatal(err)
	}
	observedAt := time.Date(2026, 8, 30, 10, 0, 0, 0, time.UTC)
	execution := controlruntime.Evaluate(context.Background(), template, binding, registry, observedAt)
	evidence, ok := execution.SealedEvidence()
	if !ok || execution.Status != controlruntime.StatusPassed {
		t.Fatalf("fixture evidence execution = %+v", execution)
	}

	directory := t.TempDir()
	bundle := evidencebundle.Bundle{
		SchemaVersion: evidencebundle.SchemaVersion, ID: "fixture.repository@1",
		CatalogSHA256: catalog.Digest(), InventorySHA256: item.Digest,
		Authority: controlprogram.AuthorityRepository,
		Entries: []evidencebundle.Entry{{
			TemplateID: template.TemplateID, ProviderID: "fixture.repository.collector@1",
			Program: program, Evidence: evidence,
		}},
	}
	bundleData, err := json.Marshal(bundle)
	if err != nil {
		t.Fatal(err)
	}
	bundlePath := filepath.Join(directory, "bundle.json")
	writeBytes(t, bundlePath, bundleData)
	bundleDigest, err := evidencebundle.BundleSHA256(bundle)
	if err != nil {
		t.Fatal(err)
	}
	policyDigest, err := evidencebundle.PolicySHA256(bundle)
	if err != nil {
		t.Fatal(err)
	}

	policyPublic, policyPrivate, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	evidencePublic, evidencePrivate, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	policyIssuedAt := observedAt.Add(-time.Minute)
	evidenceIssuedAt := observedAt.Add(time.Minute)
	verifiedAt := observedAt.Add(2 * time.Minute)
	store := trust.Store{
		SchemaVersion: trust.StoreSchema, ID: "fixture-evidence-keys", Revision: 1,
		Keys: []trust.Key{
			{ID: "evidence-2026", Algorithm: trust.AlgorithmEd25519, PublicKey: base64.StdEncoding.EncodeToString(evidencePublic), Scopes: []string{"control-evidence-repository"}, Status: "active", NotBefore: observedAt.Add(-time.Hour), NotAfter: observedAt.Add(time.Hour)},
			{ID: "policy-2026", Algorithm: trust.AlgorithmEd25519, PublicKey: base64.StdEncoding.EncodeToString(policyPublic), Scopes: []string{"control-policy-bundle"}, Status: "active", NotBefore: observedAt.Add(-time.Hour), NotAfter: observedAt.Add(time.Hour)},
		},
	}
	storeData, err := json.Marshal(store)
	if err != nil {
		t.Fatal(err)
	}
	storePath := filepath.Join(directory, "trust-store.json")
	writeBytes(t, storePath, storeData)
	policyPath := filepath.Join(directory, "policy.signature.json")
	evidencePath := filepath.Join(directory, "evidence.signature.json")
	writeSignature(t, policyPath, "control-policy-bundle", bundle.ID, policyDigest, "policy-2026", policyIssuedAt, policyPrivate)
	writeSignature(t, evidencePath, "control-evidence-repository", bundle.ID, bundleDigest, "evidence-2026", evidenceIssuedAt, evidencePrivate)
	return signedBundleFixture{catalog, item, bundlePath, storePath, policyPath, evidencePath, verifiedAt}
}

func writeSignature(t *testing.T, path, kind, artifactID, digest, keyID string, issuedAt time.Time, privateKey ed25519.PrivateKey) {
	t.Helper()
	signature := trust.Signature{
		SchemaVersion: trust.SignatureSchema, ArtifactKind: kind, ArtifactID: artifactID,
		SHA256: digest, KeyID: keyID, Algorithm: trust.AlgorithmEd25519, IssuedAt: issuedAt,
	}
	payload, err := trust.SigningPayload(signature)
	if err != nil {
		t.Fatal(err)
	}
	signature.Value = base64.StdEncoding.EncodeToString(ed25519.Sign(privateKey, payload))
	data, err := json.Marshal(signature)
	if err != nil {
		t.Fatal(err)
	}
	writeBytes(t, path, data)
}

func writeFile(t *testing.T, path, content string) { writeBytes(t, path, []byte(content)) }

func writeBytes(t *testing.T, path string, data []byte) {
	t.Helper()
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
}
