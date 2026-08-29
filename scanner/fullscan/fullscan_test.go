package fullscan

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/MarinJursic/production-readiness-checklist/scanner/catalog"
	"github.com/MarinJursic/production-readiness-checklist/scanner/controlprogram"
	"github.com/MarinJursic/production-readiness-checklist/scanner/controlprogramcatalog"
	"github.com/MarinJursic/production-readiness-checklist/scanner/controlruntime"
	"github.com/MarinJursic/production-readiness-checklist/scanner/engine"
	workspaceinventory "github.com/MarinJursic/production-readiness-checklist/scanner/inventory"
	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
)

type exactArtifactFixtureProvider struct {
	id       string
	disagree bool
}

func (provider exactArtifactFixtureProvider) ID() string { return provider.id }
func (provider exactArtifactFixtureProvider) Authority() controlprogram.Authority {
	return controlprogram.AuthorityArtifact
}
func (provider exactArtifactFixtureProvider) Collect(_ context.Context, request controlruntime.Request) (controlprogram.Evidence, error) {
	facts := map[string]controlprogram.Fact{}
	for index, contract := range request.Template.RawFactContracts {
		value := strings.Repeat("a", 64)
		if provider.disagree && index == len(request.Template.RawFactContracts)-1 {
			value = strings.Repeat("b", 64)
		}
		facts[contract.FactKey] = controlprogram.Fact{Type: controlprogram.FactDigest, Complete: true, String: &value}
	}
	return controlprogram.Evidence{
		SchemaVersion: controlprogram.EvidenceSchemaVersion, EvidenceID: "exact-artifact-fixture",
		ProgramSHA256: controlprogram.ProgramSHA256(request.Program), ControlID: request.Program.ControlID,
		ControlRevision: request.Program.ControlRevision, ControlSemanticSHA256: request.Program.ControlSemanticSHA256,
		ClauseID: request.Program.ClauseID, ClauseSHA256: request.Program.ClauseSHA256,
		ImplementationContractSHA256: request.Program.ImplementationContractSHA256,
		SubjectID:                    request.Program.SubjectID, ObservedSubjects: append([]string(nil), request.Program.Subjects...),
		InventorySHA256: request.Program.InventorySHA256, Authority: request.Program.RequiredAuthority,
		ObservedAt: time.Date(2026, 8, 29, 12, 0, 0, 0, time.UTC), Complete: true,
		Applicability:                    controlprogram.ApplicabilityApplicable,
		ApplicabilityProofContractSHA256: request.Program.ApplicabilityProofContractSHA256,
		Facts:                            facts,
	}, nil
}

func gzipBytes(t *testing.T, data []byte) []byte {
	t.Helper()
	var output bytes.Buffer
	writer, err := gzip.NewWriterLevel(&output, gzip.BestCompression)
	if err != nil {
		t.Fatal(err)
	}
	writer.Header.ModTime = time.Unix(0, 0)
	if _, err := writer.Write(data); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	return output.Bytes()
}

func TestReadCatalogDocumentAcceptsOneBoundedPlainOrCompressedFile(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "catalog.json")
	payload := []byte("{\"safe\":true}\n")
	if err := os.WriteFile(path, payload, 0o600); err != nil {
		t.Fatal(err)
	}
	plain, err := readCatalogDocument(path, 1024, "test catalog")
	if err != nil || !bytes.Equal(plain, payload) {
		t.Fatalf("plain data=%q err=%v", plain, err)
	}
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path+".gz", gzipBytes(t, payload), 0o600); err != nil {
		t.Fatal(err)
	}
	compressed, err := readCatalogDocument(path, 1024, "test catalog")
	if err != nil || !bytes.Equal(compressed, payload) {
		t.Fatalf("compressed data=%q err=%v", compressed, err)
	}
	if err := os.WriteFile(path, payload, 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := readCatalogDocument(path, 1024, "test catalog"); err == nil || !strings.Contains(err.Error(), "ambiguous") {
		t.Fatalf("plain and compressed error = %v", err)
	}
}

func TestReadCatalogDocumentBoundsExpandedData(t *testing.T) {
	path := filepath.Join(t.TempDir(), "catalog.json")
	if err := os.WriteFile(path+".gz", gzipBytes(t, bytes.Repeat([]byte("x"), 2048)), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := readCatalogDocument(path, 1024, "test catalog"); err == nil || !strings.Contains(err.Error(), "exceeds") {
		t.Fatalf("expanded limit error = %v", err)
	}
}

func repositoryRoot(t *testing.T) string {
	t.Helper()
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	return root
}

func TestAttachIncludesEveryRegisteredControlExactlyOnce(t *testing.T) {
	root := repositoryRoot(t)
	loaded, err := catalog.Load(root)
	if err != nil {
		t.Fatal(err)
	}
	item, err := workspaceinventory.Build(root)
	if err != nil {
		t.Fatal(err)
	}
	run, err := engine.New(loaded).Scan("prc/core-repository", item)
	if err != nil {
		t.Fatal(err)
	}
	complete, err := Attach(root, loaded, run)
	if err != nil {
		t.Fatal(err)
	}
	if complete.ControlCatalog == nil || complete.ControlCatalog.ControlCount != 10_042 ||
		complete.ControlCatalog.ActiveControlCount != 10_042 || complete.ControlCatalog.ContractCount != 10_042 ||
		complete.ControlCatalog.GeneratedContractCount != 0 || complete.ControlCatalog.AgentReviewedContractCount != 10_042 ||
		complete.ControlCatalog.ReviewedDeterministicCount != 686 ||
		complete.ControlCatalog.ReviewedNondeterministicCount != 9_356 ||
		complete.ControlCatalog.DeterministicBindingCount != 686 ||
		complete.ControlCatalog.ContractSchemaVersion != contractSchema || complete.ControlCatalog.ContractSHA256 == "" ||
		complete.ControlCatalog.ContractGeneratorID != contractGenerator ||
		complete.ControlCatalog.ControlCheckBindingsSchemaVersion != bindingSchema ||
		complete.ControlCatalog.ControlCheckBindingsSHA256 == "" ||
		complete.ControlCatalog.ClassificationMethodologySHA256 == "" ||
		complete.ControlCatalog.ClassificationSummarySHA256 == "" ||
		complete.ControlCatalog.ClassificationCorpusSHA256 == "" ||
		len(complete.ControlResults) != 10_042 {
		t.Fatalf("incomplete catalog attachment: %+v / %d", complete.ControlCatalog, len(complete.ControlResults))
	}
	seen := map[string]bool{}
	for index, result := range complete.ControlResults {
		if seen[result.ControlID] || index > 0 && complete.ControlResults[index-1].ControlID >= result.ControlID {
			t.Fatalf("duplicate or unordered control result %s", result.ControlID)
		}
		seen[result.ControlID] = true
		if result.ContractSHA256 == "" || result.ContractStatus != "reviewed" ||
			result.CanonicalControlID == "" || result.AutomationClass == "" || result.EvaluationClass == "" ||
			result.Classification == "" || result.ClassificationRoute == "" ||
			result.ClassificationDecisionBasis == "" || result.ClassificationRowSHA256 == "" ||
			len(result.EvidenceAuthorities) == 0 {
			t.Fatalf("control %s is missing its bound contract: %+v", result.ControlID, result)
		}
		if result.Classification == "deterministic" {
			if result.DeterministicBindingID == "" || result.DeterministicBindingSHA256 == "" ||
				result.DeterministicProgramTemplateCount < 1 || result.DeterministicProgramStatus != "blocked_provider_unregistered" ||
				result.Disposition != "blocked" || result.Coverage != "deterministic_program_provider_unregistered" {
				t.Fatalf("deterministic control %s overstates unexecuted binding coverage: %+v", result.ControlID, result)
			}
		} else if result.DeterministicBindingID != "" || result.DeterministicBindingSHA256 != "" ||
			result.Disposition != "needs_review" || result.Coverage != "nondeterministic_advisory" {
			t.Fatalf("nondeterministic control %s has an unsafe result: %+v", result.ControlID, result)
		}
	}
	identity := complete.RunID
	complete.RunID = ""
	payload, err := json.Marshal(complete)
	if err != nil {
		t.Fatal(err)
	}
	digest := sha256.Sum256(payload)
	if identity != hex.EncodeToString(digest[:]) {
		t.Fatal("complete run ID does not bind the control results")
	}
}

func TestAttachProgramExecutionsPromotesOnlyExactEvaluatorResults(t *testing.T) {
	root := repositoryRoot(t)
	loaded, err := catalog.Load(root)
	if err != nil {
		t.Fatal(err)
	}
	programs, err := controlprogramcatalog.Load(root)
	if err != nil {
		t.Fatal(err)
	}
	counts := map[string]int{}
	for _, template := range programs.Templates() {
		counts[template.ControlID]++
	}
	var template controlprogramcatalog.Template
	for _, candidate := range programs.Templates() {
		if candidate.RequiredAuthority == controlprogram.AuthorityArtifact && len(candidate.SealedParameterContracts) == 0 &&
			counts[candidate.ControlID] == 1 {
			template = candidate
			break
		}
	}
	if template.TemplateID == "" {
		t.Fatal("fixture catalog has no single-clause parameter-free artifact program")
	}
	item, err := workspaceinventory.Build(root)
	if err != nil {
		t.Fatal(err)
	}
	base, err := engine.New(loaded).Scan("prc/core-repository", item)
	if err != nil {
		t.Fatal(err)
	}
	binding := controlprogramcatalog.RuntimeBinding{
		SubjectID: "release-artifact", Subjects: []string{"release-artifact"}, InventorySHA256: item.Digest,
		ApplicabilityProofContractSHA256: strings.Repeat("c", 64), MaximumEvidenceAgeSeconds: 3600,
	}
	now := time.Date(2026, 8, 29, 12, 5, 0, 0, time.UTC)
	for _, scenario := range []struct {
		name        string
		disagree    bool
		disposition string
		terminal    string
	}{
		{name: "pass", disposition: "verified_pass", terminal: base.TerminalState},
		{name: "fail", disagree: true, disposition: "confirmed_failure", terminal: "no_go"},
	} {
		t.Run(scenario.name, func(t *testing.T) {
			registry, registryErr := controlruntime.NewRegistry(exactArtifactFixtureProvider{id: template.CollectorContract.CollectorID, disagree: scenario.disagree})
			if registryErr != nil {
				t.Fatal(registryErr)
			}
			execution := controlruntime.Evaluate(context.Background(), template, binding, registry, now)
			complete, attachErr := AttachProgramExecutions(root, loaded, base, []controlruntime.Execution{execution})
			if attachErr != nil {
				t.Fatal(attachErr)
			}
			var result model.ControlResult
			for _, candidate := range complete.ControlResults {
				if candidate.ControlID == template.ControlID {
					result = candidate
					break
				}
			}
			if result.Disposition != scenario.disposition || result.Coverage != "deterministic_program_complete" ||
				result.Authority != "deterministic_exact" || len(result.DeterministicClauseResults) != 1 || complete.TerminalState != scenario.terminal {
				t.Fatalf("exact result was not aggregated safely: %+v / %s", result, complete.TerminalState)
			}
			if complete.ControlCatalog.DeterministicProgramExecutedCount != 1 {
				t.Fatalf("executed count = %d", complete.ControlCatalog.DeterministicProgramExecutedCount)
			}
			if len(complete.DeterministicEvidence) != 1 ||
				controlprogram.EvidenceSHA256(complete.DeterministicEvidence[0]) != result.DeterministicClauseResults[0].EvidenceSHA256 {
				t.Fatal("exact execution did not retain replayable evidence")
			}
		})
	}
	forged := controlruntime.Execution{TemplateID: template.TemplateID}
	if _, err := AttachProgramExecutions(root, loaded, base, []controlruntime.Execution{forged}); err == nil {
		t.Fatal("externally constructed exact result was accepted")
	}
}

func TestControlResultNeverTurnsPartialEvidenceIntoAFullPass(t *testing.T) {
	control := model.Control{ID: "USEQ-AAAAAAAA", Status: "active", Revision: 1, Statement: "A broad outcome is true."}
	bindingID, bindingDigest := control.ID+"@1", strings.Repeat("b", 64)
	contract := controlContract{
		ControlID: control.ID, Revision: 1, ContractStatus: "reviewed", ReviewerStatus: "agent_reviewed",
		Classification: "deterministic", ClassificationRoute: "local_static", ClassificationDecisionBasis: "strength_audit_confirmed",
		ClassificationRowSHA256: strings.Repeat("c", 64), DeterministicBindingID: &bindingID,
		DeterministicBindingSHA256: &bindingDigest, CanonicalControlID: control.ID,
		EvaluationClass: "repository", AutomationClass: "deterministic_candidate", ApplicabilityClass: "scope_required",
		Atomicity: "apparently_atomic", EvidenceAuthorities: []string{"repository"}, ContractSHA256: strings.Repeat("a", 64),
	}
	passed := controlResult(control, contract, []string{"PRC-A-TEST-001"}, []model.AssertionResult{{
		AssertionID: "PRC-A-TEST-001", Assessment: "pass",
	}}, 1, nil)
	if passed.Disposition != "blocked" || passed.Coverage != "deterministic_program_provider_unregistered" || passed.Authority != "deterministic_partial" {
		t.Fatalf("partial assertion was overstated: %+v", passed)
	}
	failed := controlResult(control, contract, []string{"PRC-A-TEST-001"}, []model.AssertionResult{{
		AssertionID: "PRC-A-TEST-001", Assessment: "fail",
	}}, 1, nil)
	if failed.Disposition != "blocked" || failed.Coverage != "deterministic_program_provider_unregistered" {
		t.Fatalf("observed failure was hidden: %+v", failed)
	}
}

func TestReviewedContractAndBindingMismatchesFailClosed(t *testing.T) {
	root := repositoryRoot(t)
	registry, registryDigest, err := loadRegistry(root)
	if err != nil {
		t.Fatal(err)
	}
	contracts, _, err := loadContracts(root, registry, registryDigest)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("envelope digest", func(t *testing.T) {
		changed := contracts
		changed.ClassificationCorpusSHA256 = strings.Repeat("0", 64)
		if err := validateBindingArtifact(root, changed, registry); err == nil {
			t.Fatal("stale classification corpus digest was accepted")
		}
	})
	t.Run("classification row", func(t *testing.T) {
		changed := contracts
		changed.Contracts = append([]controlContract(nil), contracts.Contracts...)
		for index := range changed.Contracts {
			if changed.Contracts[index].Classification == "deterministic" {
				changed.Contracts[index].ClassificationRowSHA256 = strings.Repeat("0", 64)
				changed.Contracts[index].ContractSHA256 = contractDigest(changed.Contracts[index])
				break
			}
		}
		if err := validateBindingArtifact(root, changed, registry); err == nil {
			t.Fatal("classification-row mismatch against binding was accepted")
		}
	})
	t.Run("missing deterministic identity", func(t *testing.T) {
		changed := contracts
		changed.Contracts = append([]controlContract(nil), contracts.Contracts...)
		for index := range changed.Contracts {
			if changed.Contracts[index].Classification == "deterministic" {
				changed.Contracts[index].DeterministicBindingID = nil
				changed.Contracts[index].ContractSHA256 = contractDigest(changed.Contracts[index])
				break
			}
		}
		if err := validateContracts(changed, registry, registryDigest); err == nil {
			t.Fatal("deterministic contract without binding identity was accepted")
		}
	})
}
