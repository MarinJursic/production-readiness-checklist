package evidenceset

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/MarinJursic/production-readiness-checklist/scanner/controlprogram"
	"github.com/MarinJursic/production-readiness-checklist/scanner/controlprogramcatalog"
	"github.com/MarinJursic/production-readiness-checklist/scanner/controlruntime"
	"github.com/MarinJursic/production-readiness-checklist/scanner/evidencebundle"
	workspaceinventory "github.com/MarinJursic/production-readiness-checklist/scanner/inventory"
	"github.com/MarinJursic/production-readiness-checklist/scanner/repositoryevidence"
	"github.com/MarinJursic/production-readiness-checklist/scanner/trust"
)

func TestManifestRequiresUniqueOrderedAuthoritiesAndSiblingFiles(t *testing.T) {
	valid := Manifest{
		SchemaVersion: SchemaVersion, TrustStoreFile: "trust-store.json",
		Bundles: []BundleReference{
			{Authority: controlprogram.AuthorityArtifact, BundleFile: "artifact.json", PolicySignatureFile: "artifact-policy.json", EvidenceSignatureFile: "artifact-evidence.json"},
			{Authority: controlprogram.AuthorityRepository, BundleFile: "repository.json", PolicySignatureFile: "repository-policy.json", EvidenceSignatureFile: "repository-evidence.json"},
		},
	}
	if err := validate(valid); err != nil {
		t.Fatal(err)
	}

	duplicateAuthority := valid
	duplicateAuthority.Bundles = append([]BundleReference(nil), valid.Bundles...)
	duplicateAuthority.Bundles[1].Authority = controlprogram.AuthorityArtifact
	if err := validate(duplicateAuthority); err == nil {
		t.Fatal("duplicate authority was accepted")
	}

	escape := valid
	escape.Bundles = append([]BundleReference(nil), valid.Bundles...)
	escape.Bundles[0].BundleFile = "../artifact.json"
	if err := validate(escape); err == nil {
		t.Fatal("path escape was accepted")
	}

	reusedFile := valid
	reusedFile.Bundles = append([]BundleReference(nil), valid.Bundles...)
	reusedFile.Bundles[1].EvidenceSignatureFile = valid.Bundles[0].BundleFile
	if err := validate(reusedFile); err == nil {
		t.Fatal("one file reused in two trust roles was accepted")
	}
}

func TestManifestReaderRejectsSymlinksAndUnknownFields(t *testing.T) {
	directory := t.TempDir()
	validPath := filepath.Join(directory, "set.json")
	if err := os.WriteFile(validPath, []byte(`{
  "schema_version":"prc.authoritative-evidence-set/v0.1",
  "trust_store_file":"trust-store.json",
  "bundles":[{"authority":"repository","bundle_file":"bundle.json","policy_signature_file":"policy.json","evidence_signature_file":"evidence.json"}]
}`), 0o600); err != nil {
		t.Fatal(err)
	}
	manifest, loadedDirectory, err := load(validPath)
	if err != nil || manifest.SchemaVersion != SchemaVersion || loadedDirectory != directory {
		t.Fatalf("loaded manifest = %+v directory=%q err=%v", manifest, loadedDirectory, err)
	}

	unknownPath := filepath.Join(directory, "unknown.json")
	if err := os.WriteFile(unknownPath, []byte(`{"schema_version":"prc.authoritative-evidence-set/v0.1","trust_store_file":"trust-store.json","bundles":[],"unexpected":true}`), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, _, err := load(unknownPath); err == nil {
		t.Fatal("unknown manifest field was accepted")
	}

	linkPath := filepath.Join(directory, "link.json")
	if err := os.Symlink(validPath, linkPath); err != nil {
		t.Fatal(err)
	}
	if _, _, err := load(linkPath); err == nil {
		t.Fatal("symlink manifest was accepted")
	}
}

func TestEvidenceSetFilesMustBeDistinctRegularAndUnchanged(t *testing.T) {
	directory := t.TempDir()
	manifest := Manifest{
		SchemaVersion: SchemaVersion, TrustStoreFile: "trust-store.json",
		Bundles: []BundleReference{{
			Authority: controlprogram.AuthorityRepository, BundleFile: "bundle.json",
			PolicySignatureFile: "policy.json", EvidenceSignatureFile: "evidence.json",
		}},
	}
	for _, name := range []string{"trust-store.json", "bundle.json", "policy.json", "evidence.json"} {
		if err := os.WriteFile(filepath.Join(directory, name), []byte(name), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	snapshots, err := snapshotFiles(directory, manifest)
	if err != nil || len(snapshots) != 4 {
		t.Fatalf("snapshots=%+v err=%v", snapshots, err)
	}
	if err := os.WriteFile(filepath.Join(directory, "bundle.json"), []byte("changed bundle"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := verifyUnchanged(snapshots); err == nil {
		t.Fatal("changed referenced file was accepted")
	}

	if err := os.Remove(filepath.Join(directory, "evidence.json")); err != nil {
		t.Fatal(err)
	}
	if err := os.Link(filepath.Join(directory, "policy.json"), filepath.Join(directory, "evidence.json")); err != nil {
		t.Fatal(err)
	}
	if _, err := snapshotFiles(directory, manifest); err == nil {
		t.Fatal("one inode reused for policy and evidence signatures was accepted")
	}
}

func TestEvidenceSetVerifiesMultipleSignedAuthoritiesTogether(t *testing.T) {
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	programs, err := controlprogramcatalog.Load(root)
	if err != nil {
		t.Fatal(err)
	}
	var template controlprogramcatalog.Template
	for _, candidate := range programs.Templates() {
		if candidate.CollectorContract.CollectorID == repositoryevidence.DocumentedCommandsCollectorID {
			template = candidate
			break
		}
	}
	workspace := t.TempDir()
	writeTestFile(t, filepath.Join(workspace, "package.json"), []byte(`{"scripts":{"build":"node build.mjs","test":"node --test"}}`))
	writeTestFile(t, filepath.Join(workspace, "README.md"), []byte("```sh\nnpm run build\nnpm test\n```\n"))
	item, err := workspaceinventory.Build(workspace)
	if err != nil {
		t.Fatal(err)
	}
	binding, ok := repositoryevidence.Binding(item, template)
	if !ok {
		t.Fatal("documented command binding is missing")
	}
	program, err := template.Program(binding)
	if err != nil {
		t.Fatal(err)
	}
	collector, err := repositoryevidence.NewDocumentedCommandsProvider(item)
	if err != nil {
		t.Fatal(err)
	}
	registry, err := controlruntime.NewRegistry(collector)
	if err != nil {
		t.Fatal(err)
	}
	observedAt := time.Date(2026, 8, 30, 10, 0, 0, 0, time.UTC)
	execution := controlruntime.Evaluate(context.Background(), template, binding, registry, observedAt)
	evidence, ok := execution.SealedEvidence()
	if !ok || execution.Status != controlruntime.StatusPassed {
		t.Fatalf("fixture execution = %+v", execution)
	}

	directory := t.TempDir()
	bundle := evidencebundle.Bundle{
		SchemaVersion: evidencebundle.SchemaVersion, ID: "set.repository@1",
		CatalogSHA256: programs.Digest(), InventorySHA256: item.Digest,
		Authority: controlprogram.AuthorityRepository,
		Entries: []evidencebundle.Entry{{
			TemplateID: template.TemplateID, ProviderID: "set.repository.collector@1",
			Program: program, Evidence: evidence,
		}},
	}
	bundleData, err := json.Marshal(bundle)
	if err != nil {
		t.Fatal(err)
	}
	writeTestFile(t, filepath.Join(directory, "repository.json"), bundleData)
	bundleDigest, err := evidencebundle.BundleSHA256(bundle)
	if err != nil {
		t.Fatal(err)
	}
	policyDigest, err := evidencebundle.PolicySHA256(bundle)
	if err != nil {
		t.Fatal(err)
	}

	var artifactTemplate controlprogramcatalog.Template
	for _, candidate := range programs.Templates() {
		if candidate.RequiredAuthority == controlprogram.AuthorityArtifact {
			artifactTemplate = candidate
			break
		}
	}
	if artifactTemplate.TemplateID == "" {
		t.Fatal("artifact template is missing")
	}
	digestValue := strings.Repeat("b", 64)
	digestAlgorithm := "sha256"
	proofContract := strings.Repeat("a", 64)
	artifactProgram, err := artifactTemplate.Program(controlprogramcatalog.RuntimeBinding{
		SubjectID: "project", Subjects: []string{"project"}, InventorySHA256: item.Digest,
		ApplicabilityProofContractSHA256: proofContract, MaximumEvidenceAgeSeconds: 3600,
		AuthenticatedPolicyParameters: map[string]controlprogram.Parameter{
			"prc_03_003_c1.approved_digest_algorithms": {
				Type: controlprogram.FactStringSet, Strings: []string{digestAlgorithm},
			},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	artifactEvidence := controlprogram.Evidence{
		SchemaVersion: controlprogram.EvidenceSchemaVersion, EvidenceID: "set-artifact-evidence",
		ProgramSHA256: controlprogram.ProgramSHA256(artifactProgram),
		ControlID:     artifactProgram.ControlID, ControlRevision: artifactProgram.ControlRevision,
		ControlSemanticSHA256: artifactProgram.ControlSemanticSHA256,
		ClauseID:              artifactProgram.ClauseID, ClauseSHA256: artifactProgram.ClauseSHA256,
		ImplementationContractSHA256: artifactProgram.ImplementationContractSHA256,
		SubjectID:                    artifactProgram.SubjectID, ObservedSubjects: artifactProgram.Subjects,
		InventorySHA256: item.Digest, Authority: controlprogram.AuthorityArtifact,
		ObservedAt: observedAt, Complete: true, Applicability: controlprogram.ApplicabilityApplicable,
		ApplicabilityProofContractSHA256: proofContract,
		Facts: map[string]controlprogram.Fact{
			"prc_03_003_c1.release_record_digest": {
				Type: controlprogram.FactDigest, Complete: true, String: &digestValue,
			},
			"prc_03_003_c1.exact_artifact_bytes_digest": {
				Type: controlprogram.FactDigest, Complete: true, String: &digestValue,
			},
			"prc_03_003_c1.recorded_digest_algorithm": {
				Type: controlprogram.FactState, Complete: true, String: &digestAlgorithm,
			},
		},
	}
	artifactBundle := evidencebundle.Bundle{
		SchemaVersion: evidencebundle.SchemaVersion, ID: "set.artifact@1",
		CatalogSHA256: programs.Digest(), InventorySHA256: item.Digest,
		Authority: controlprogram.AuthorityArtifact,
		Entries: []evidencebundle.Entry{{
			TemplateID: artifactTemplate.TemplateID, ProviderID: "set.artifact.collector@1",
			Program: artifactProgram, Evidence: artifactEvidence,
		}},
	}
	artifactBundleData, err := json.Marshal(artifactBundle)
	if err != nil {
		t.Fatal(err)
	}
	writeTestFile(t, filepath.Join(directory, "artifact.json"), artifactBundleData)
	artifactBundleDigest, err := evidencebundle.BundleSHA256(artifactBundle)
	if err != nil {
		t.Fatal(err)
	}
	artifactPolicyDigest, err := evidencebundle.PolicySHA256(artifactBundle)
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
	artifactEvidencePublic, artifactEvidencePrivate, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	store := trust.Store{
		SchemaVersion: trust.StoreSchema, ID: "set-test-keys", Revision: 1,
		Keys: []trust.Key{
			{ID: "artifact-evidence-key", Algorithm: trust.AlgorithmEd25519, PublicKey: base64.StdEncoding.EncodeToString(artifactEvidencePublic), Scopes: []string{"control-evidence-artifact"}, Status: "active", NotBefore: observedAt.Add(-time.Hour), NotAfter: observedAt.Add(time.Hour)},
			{ID: "evidence-key", Algorithm: trust.AlgorithmEd25519, PublicKey: base64.StdEncoding.EncodeToString(evidencePublic), Scopes: []string{"control-evidence-repository"}, Status: "active", NotBefore: observedAt.Add(-time.Hour), NotAfter: observedAt.Add(time.Hour)},
			{ID: "policy-key", Algorithm: trust.AlgorithmEd25519, PublicKey: base64.StdEncoding.EncodeToString(policyPublic), Scopes: []string{"control-policy-bundle"}, Status: "active", NotBefore: observedAt.Add(-time.Hour), NotAfter: observedAt.Add(time.Hour)},
		},
	}
	storeData, err := json.Marshal(store)
	if err != nil {
		t.Fatal(err)
	}
	writeTestFile(t, filepath.Join(directory, "trust-store.json"), storeData)
	writeTestSignature(t, filepath.Join(directory, "artifact-policy.json"), "control-policy-bundle", artifactBundle.ID, artifactPolicyDigest, "policy-key", observedAt.Add(-time.Minute), policyPrivate)
	writeTestSignature(t, filepath.Join(directory, "artifact-evidence.json"), "control-evidence-artifact", artifactBundle.ID, artifactBundleDigest, "artifact-evidence-key", observedAt.Add(time.Minute), artifactEvidencePrivate)
	writeTestSignature(t, filepath.Join(directory, "repository-policy.json"), "control-policy-bundle", bundle.ID, policyDigest, "policy-key", observedAt.Add(-time.Minute), policyPrivate)
	writeTestSignature(t, filepath.Join(directory, "repository-evidence.json"), "control-evidence-repository", bundle.ID, bundleDigest, "evidence-key", observedAt.Add(time.Minute), evidencePrivate)
	manifest := Manifest{
		SchemaVersion: SchemaVersion, TrustStoreFile: "trust-store.json",
		Bundles: []BundleReference{
			{
				Authority: controlprogram.AuthorityArtifact, BundleFile: "artifact.json",
				PolicySignatureFile: "artifact-policy.json", EvidenceSignatureFile: "artifact-evidence.json",
			},
			{
				Authority: controlprogram.AuthorityRepository, BundleFile: "repository.json",
				PolicySignatureFile: "repository-policy.json", EvidenceSignatureFile: "repository-evidence.json",
			},
		},
	}
	manifestData, err := json.Marshal(manifest)
	if err != nil {
		t.Fatal(err)
	}
	manifestPath := filepath.Join(directory, "evidence-set.json")
	writeTestFile(t, manifestPath, manifestData)

	executions, verifications, err := VerifyAndEvaluate(programs, item, manifestPath, observedAt.Add(2*time.Minute))
	if err != nil {
		t.Fatal(err)
	}
	if len(executions) != 2 || executions[0].Status != controlruntime.StatusPassed ||
		executions[1].Status != controlruntime.StatusPassed || len(verifications) != 2 ||
		verifications[0].Authority != "artifact" || verifications[0].EntryCount != 1 ||
		verifications[1].Authority != "repository" || verifications[1].EntryCount != 1 {
		t.Fatalf("executions=%+v verifications=%+v", executions, verifications)
	}

	repositorySignaturePath := filepath.Join(directory, "repository-evidence.json")
	repositorySignature, err := os.ReadFile(repositorySignaturePath)
	if err != nil {
		t.Fatal(err)
	}
	repositorySignature[len(repositorySignature)-2] ^= 1
	writeTestFile(t, repositorySignaturePath, repositorySignature)
	executions, verifications, err = VerifyAndEvaluate(programs, item, manifestPath, observedAt.Add(2*time.Minute))
	if err == nil || executions != nil || verifications != nil {
		t.Fatalf("tampered second authority returned partial results: executions=%+v verifications=%+v err=%v", executions, verifications, err)
	}
}

func writeTestSignature(t *testing.T, path, kind, artifactID, digest, keyID string, issuedAt time.Time, privateKey ed25519.PrivateKey) {
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
	writeTestFile(t, path, data)
}

func writeTestFile(t *testing.T, path string, data []byte) {
	t.Helper()
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
}
