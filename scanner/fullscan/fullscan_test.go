package fullscan

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"

	"github.com/MarinJursic/production-readiness-checklist/scanner/catalog"
	"github.com/MarinJursic/production-readiness-checklist/scanner/engine"
	workspaceinventory "github.com/MarinJursic/production-readiness-checklist/scanner/inventory"
	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
)

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
		complete.ControlCatalog.GeneratedContractCount != 10_042 || complete.ControlCatalog.ExpertReviewedContractCount != 0 ||
		complete.ControlCatalog.ContractSchemaVersion != contractSchema || complete.ControlCatalog.ContractSHA256 == "" ||
		len(complete.ControlResults) != 10_042 {
		t.Fatalf("incomplete catalog attachment: %+v / %d", complete.ControlCatalog, len(complete.ControlResults))
	}
	seen := map[string]bool{}
	for index, result := range complete.ControlResults {
		if seen[result.ControlID] || index > 0 && complete.ControlResults[index-1].ControlID >= result.ControlID {
			t.Fatalf("duplicate or unordered control result %s", result.ControlID)
		}
		seen[result.ControlID] = true
		if result.ContractSHA256 == "" || result.ContractStatus != "generated_unreviewed" ||
			result.CanonicalControlID == "" || result.AutomationClass == "" || result.EvaluationClass == "" ||
			len(result.EvidenceAuthorities) == 0 {
			t.Fatalf("control %s is missing its bound contract: %+v", result.ControlID, result)
		}
		if result.Disposition == "partially_verified" && result.Coverage != "partial_assertions" {
			t.Fatalf("control %s overstates its coverage: %+v", result.ControlID, result)
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

func TestControlResultNeverTurnsPartialEvidenceIntoAFullPass(t *testing.T) {
	control := model.Control{ID: "USEQ-AAAAAAAA", Status: "active", Revision: 1, Statement: "A broad outcome is true."}
	contract := controlContract{
		ControlID: control.ID, Revision: 1, ContractStatus: "generated_unreviewed", CanonicalControlID: control.ID,
		EvaluationClass: "repository", AutomationClass: "ai_advisory_candidate", ApplicabilityClass: "scope_required",
		Atomicity: "apparently_atomic", EvidenceAuthorities: []string{"repository"}, ContractSHA256: strings.Repeat("a", 64),
	}
	passed := controlResult(control, contract, []string{"PRC-A-TEST-001"}, []model.AssertionResult{{
		AssertionID: "PRC-A-TEST-001", Assessment: "pass",
	}})
	if passed.Disposition != "partially_verified" || passed.Authority != "deterministic_partial" {
		t.Fatalf("partial assertion was overstated: %+v", passed)
	}
	failed := controlResult(control, contract, []string{"PRC-A-TEST-001"}, []model.AssertionResult{{
		AssertionID: "PRC-A-TEST-001", Assessment: "fail",
	}})
	if failed.Disposition != "confirmed_failure" {
		t.Fatalf("observed failure was hidden: %+v", failed)
	}
}
