package exception

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/MarinJursic/production-readiness-checklist/scanner/catalog"
	"github.com/MarinJursic/production-readiness-checklist/scanner/engine"
	"github.com/MarinJursic/production-readiness-checklist/scanner/inventory"
	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
	"github.com/MarinJursic/production-readiness-checklist/scanner/trust"
	"gopkg.in/yaml.v3"
)

func fixtureRun(t *testing.T) model.RunResult {
	t.Helper()
	target := t.TempDir()
	if err := os.WriteFile(filepath.Join(target, "fixture.txt"), []byte("fixture\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	item, err := inventory.Build(target)
	if err != nil {
		t.Fatal(err)
	}
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	catalogValue, err := catalog.Load(root)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 8, 23, 12, 0, 0, 0, time.UTC)
	runner := engine.New(catalogValue)
	runner.Now = func() time.Time { return now }
	run, err := runner.Scan("prc/core-repository", item)
	if err != nil {
		t.Fatal(err)
	}
	if len(run.Findings) == 0 {
		t.Fatal("fixture run has no finding")
	}
	return run
}

func fixtureRecord(run model.RunResult) Record {
	finding := run.Findings[0]
	return Record{
		SchemaVersion: Schema, ID: "PRC-EXC-FIXTURE-001", Status: "approved",
		Run: RunBinding{
			RunID: run.RunID, InventoryDigest: run.Inventory.Digest,
			ProfileID: run.Plan.ProfileID, ProfileVersion: run.Plan.ProfileVersion,
			TargetName: run.Inventory.TargetName, TargetCommit: run.Inventory.GitCommit, ProjectID: run.Plan.ProjectID,
			ArtifactDigests:    append([]string(nil), run.Plan.ArtifactDigests...),
			TargetEnvironments: append([]string(nil), run.Plan.TargetEnvironments...),
		},
		Finding: FindingBinding{
			FindingID: finding.ID, FindingFingerprint: finding.Fingerprint,
			AssertionID: finding.AssertionID, ControlIDs: append([]string(nil), finding.ControlIDs...),
		},
		RequestedBy: Actor{ID: "requester", Name: "Requesting engineer", Authority: "engineering"},
		RiskOwner:   Actor{ID: "risk-owner", Name: "Accountable owner", Authority: "executive"},
		Reviewers:   []Actor{{ID: "security-reviewer", Name: "Security reviewer", Authority: "security"}},
		Risk: RiskAnalysis{
			Title: "Time-bounded fixture exception", Rationale: "The fixture exercises signed exception verification.",
			Likelihood: "unlikely", Impact: "high", WorstCredibleOutcome: "The missing control contributes to a release defect.",
		},
		CompensatingControls: []CompensatingControl{{
			Description:        "A separately verified temporary control limits exposure.",
			EvidenceReferences: []string{strings.Repeat("b", 64)},
		}},
		Monitoring: Monitoring{Owner: "operations", Signal: "Alert on the affected behavior.", Response: "Disable the affected feature."},
		Remediation: Remediation{
			Owner: "engineering", Plan: "Implement and independently verify the missing control.",
			DueAt: run.CompletedAt.Add(48 * time.Hour),
		},
		ApprovedAt: run.CompletedAt.Add(time.Hour), ExpiresAt: run.CompletedAt.Add(7 * 24 * time.Hour),
	}
}

func signedException(t *testing.T, record Record) (Loaded, trust.LoadedStore, trust.Signature) {
	t.Helper()
	directory := t.TempDir()
	recordPath := filepath.Join(directory, "exception.yaml")
	recordPayload, err := yaml.Marshal(record)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(recordPath, recordPayload, 0o600); err != nil {
		t.Fatal(err)
	}
	loaded, err := Load(recordPath)
	if err != nil {
		t.Fatal(err)
	}
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	storeDocument := trust.Store{
		SchemaVersion: trust.StoreSchema, ID: "exception-keys", Revision: 1,
		Keys: []trust.Key{{
			ID: "risk-owners", Algorithm: trust.AlgorithmEd25519,
			PublicKey: base64.StdEncoding.EncodeToString(publicKey), Scopes: []string{"risk-exception"}, Status: "active",
			NotBefore: record.ApprovedAt.Add(-time.Hour), NotAfter: record.ExpiresAt.Add(time.Hour),
		}},
	}
	storePath := filepath.Join(directory, "store.json")
	storePayload, err := json.Marshal(storeDocument)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(storePath, storePayload, 0o600); err != nil {
		t.Fatal(err)
	}
	store, err := trust.LoadStore(storePath)
	if err != nil {
		t.Fatal(err)
	}
	signature := trust.Signature{
		SchemaVersion: trust.SignatureSchema, ArtifactKind: "risk-exception", ArtifactID: loaded.Record.ID,
		SHA256: loaded.Digest, KeyID: "risk-owners", Algorithm: trust.AlgorithmEd25519,
		IssuedAt: record.ApprovedAt,
	}
	payload, err := trust.SigningPayload(signature)
	if err != nil {
		t.Fatal(err)
	}
	signature.Value = base64.StdEncoding.EncodeToString(ed25519.Sign(privateKey, payload))
	return loaded, store, signature
}

func TestVerifyAcceptsCurrentSignedExceptionWithoutChangingGate(t *testing.T) {
	run := fixtureRun(t)
	record := fixtureRecord(run)
	loaded, store, signature := signedException(t, record)
	verification, err := Verify(loaded, run, store, signature, record.ApprovedAt.Add(time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	if verification.Disposition != "accepted_risk_exception" || !verification.Signature.Verified ||
		verification.ExceptionDigest != loaded.Digest || !strings.Contains(verification.GateEffect, "unchanged") ||
		run.TerminalState == "profile_satisfied" {
		t.Fatalf("verification = %+v run state = %s", verification, run.TerminalState)
	}
}

func TestVerifyRejectsExpiredMismatchedAndMutatedExceptions(t *testing.T) {
	run := fixtureRun(t)
	record := fixtureRecord(run)
	loaded, store, signature := signedException(t, record)

	if _, err := Verify(loaded, run, store, signature, record.ExpiresAt); err == nil || !strings.Contains(err.Error(), "current") {
		t.Fatalf("expiry error = %v", err)
	}

	mismatchedRun := run
	mismatchedRun.RunID = strings.Repeat("c", 64)
	if _, err := Verify(loaded, mismatchedRun, store, signature, record.ApprovedAt.Add(time.Hour)); err == nil ||
		!strings.Contains(err.Error(), "run scope") {
		t.Fatalf("run mismatch error = %v", err)
	}

	loaded.Record.Risk.Rationale = "A changed rationale invalidates the canonical exception."
	if _, err := Verify(loaded, run, store, signature, record.ApprovedAt.Add(time.Hour)); err == nil ||
		!strings.Contains(err.Error(), "digest") {
		t.Fatalf("mutation error = %v", err)
	}
}

func TestRecordRejectsSelfApprovalAndUnboundedExpiry(t *testing.T) {
	run := fixtureRun(t)
	record := fixtureRecord(run)
	record.RiskOwner = record.RequestedBy
	if err := normalizeAndValidate(&record); err == nil || !strings.Contains(err.Error(), "independent") {
		t.Fatalf("self-approval error = %v", err)
	}
	record = fixtureRecord(run)
	record.ExpiresAt = record.ApprovedAt.Add(367 * 24 * time.Hour)
	record.Remediation.DueAt = record.ExpiresAt
	if err := normalizeAndValidate(&record); err == nil || !strings.Contains(err.Error(), "dates") {
		t.Fatalf("unbounded expiry error = %v", err)
	}
}
