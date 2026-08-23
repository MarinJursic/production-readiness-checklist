package invalidation

import (
	"strings"
	"testing"
	"time"

	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
)

func TestAnalyzeInvalidatesOnlyRelevantRuleInputs(t *testing.T) {
	baseInventory := testInventory("base", []model.FileRecord{
		{Path: "LICENSE", Size: 3, SHA256: strings.Repeat("a", 64), Mode: 0o644},
		{Path: "README.md", Size: 3, SHA256: strings.Repeat("b", 64), Mode: 0o644},
	})
	currentInventory := testInventory("current", []model.FileRecord{
		{Path: "LICENSE", Size: 3, SHA256: strings.Repeat("a", 64), Mode: 0o644},
		{Path: "README.md", Size: 4, SHA256: strings.Repeat("c", 64), Mode: 0o644},
	})
	readme := testAssertion("README", []string{"README.md"})
	license := testAssertion("LICENSE", []string{"LICENSE"})
	basePlan := testPlan(baseInventory.Digest, "base-plan", readme, license)
	currentPlan := testPlan(currentInventory.Digest, "current-plan", readme, license)
	base := model.RunResult{RunID: strings.Repeat("f", 64), Inventory: baseInventory, Plan: basePlan, Results: []model.AssertionResult{
		{AssertionID: readme.ID, Assessment: "pass"}, {AssertionID: license.ID, Assessment: "pass"},
	}}
	report, err := Analyze(base, currentInventory, currentPlan, map[string]model.Assertion{
		readme.ID: readme, license.ID: license,
	}, time.Unix(1, 0))
	if err != nil {
		t.Fatal(err)
	}
	readmeImpact := impactByID(t, report, readme.ID)
	if readmeImpact.Conclusion != "invalidated" || !hasReason(readmeImpact, "relevant_repository_input_changed") ||
		len(readmeImpact.RelevantChangedPaths) != 1 || readmeImpact.RelevantChangedPaths[0] != "README.md" {
		t.Fatalf("README impact = %+v", readmeImpact)
	}
	licenseImpact := impactByID(t, report, license.ID)
	if licenseImpact.Conclusion != "unchanged_inputs" || licenseImpact.ReuseAllowed ||
		!hasReason(licenseImpact, "evidence_target_requires_rebinding") {
		t.Fatalf("license impact = %+v", licenseImpact)
	}
}

func TestAnalyzeAllowsReuseOnlyForExactBoundInput(t *testing.T) {
	inventory := testInventory("same", []model.FileRecord{{Path: "README.md", Size: 3, SHA256: strings.Repeat("a", 64), Mode: 0o644}})
	assertion := testAssertion("README", []string{"README.md"})
	plan := testPlan(inventory.Digest, "same-plan", assertion)
	base := model.RunResult{RunID: strings.Repeat("f", 64), Inventory: inventory, Plan: plan,
		Results: []model.AssertionResult{{AssertionID: assertion.ID, Assessment: "pass"}}}
	report, err := Analyze(base, inventory, plan, map[string]model.Assertion{assertion.ID: assertion}, time.Unix(1, 0))
	if err != nil {
		t.Fatal(err)
	}
	impact := impactByID(t, report, assertion.ID)
	if impact.Conclusion != "unchanged_inputs" || !impact.ReuseAllowed || impact.FreshEvaluationRequired || report.Summary.Reusable != 1 {
		t.Fatalf("impact = %+v report=%+v", impact, report.Summary)
	}
}

func TestAnalyzeInvalidatesLegacyUnboundPlanAndFreshEvidence(t *testing.T) {
	inventory := testInventory("same", nil)
	manual := model.Assertion{ID: "PRC-A-MANUAL", Revision: 1, ImplementationID: "prc.native.manual-evidence@0.1", Parameters: map[string]any{}}
	currentPlan := testPlan(inventory.Digest, "current", manual)
	legacyPlan := currentPlan
	legacyPlan.EngineVersion, legacyPlan.ProfileDigest, legacyPlan.CatalogDigest = "", "", ""
	legacyPlan.Assertions[0].AssertionRevision, legacyPlan.Assertions[0].DefinitionDigest = 0, ""
	base := model.RunResult{RunID: strings.Repeat("f", 64), Inventory: inventory, Plan: legacyPlan,
		Results: []model.AssertionResult{{AssertionID: manual.ID, Assessment: "manual_review"}}}
	report, err := Analyze(base, inventory, currentPlan, map[string]model.Assertion{manual.ID: manual}, time.Unix(1, 0))
	if err != nil {
		t.Fatal(err)
	}
	impact := impactByID(t, report, manual.ID)
	if !hasReason(impact, "base_plan_lacks_rule_binding") || !hasReason(impact, "freshness_not_proven") || impact.ReuseAllowed {
		t.Fatalf("impact = %+v", impact)
	}
}

func TestAnalyzeInvalidatesGitRevisionWithoutFileChanges(t *testing.T) {
	baseInventory := testInventory("a", nil)
	baseInventory.GitCommit = strings.Repeat("1", 40)
	currentInventory := testInventory("b", nil)
	currentInventory.GitCommit = strings.Repeat("2", 40)
	assertion := model.Assertion{ID: "PRC-A-GIT", Revision: 1, ImplementationID: "prc.native.git-revision@0.1", Parameters: map[string]any{}}
	basePlan := testPlan(baseInventory.Digest, "a", assertion)
	currentPlan := testPlan(currentInventory.Digest, "b", assertion)
	base := model.RunResult{RunID: strings.Repeat("f", 64), Inventory: baseInventory, Plan: basePlan,
		Results: []model.AssertionResult{{AssertionID: assertion.ID, Assessment: "pass"}}}
	report, err := Analyze(base, currentInventory, currentPlan, map[string]model.Assertion{assertion.ID: assertion}, time.Unix(1, 0))
	if err != nil {
		t.Fatal(err)
	}
	if impact := impactByID(t, report, assertion.ID); !hasReason(impact, "git_revision_changed") || impact.Conclusion != "invalidated" {
		t.Fatalf("impact = %+v", impact)
	}
}

func TestAnalyzeInvalidatesCatalogIdentityChange(t *testing.T) {
	inventory := testInventory("same", nil)
	assertion := testAssertion("README", []string{"README.md"})
	basePlan := testPlan(inventory.Digest, "base", assertion)
	currentPlan := testPlan(inventory.Digest, "current", assertion)
	currentPlan.CatalogDigest = strings.Repeat("e", 64)
	base := model.RunResult{RunID: strings.Repeat("f", 64), Inventory: inventory, Plan: basePlan,
		Results: []model.AssertionResult{{AssertionID: assertion.ID, Assessment: "pass"}}}
	report, err := Analyze(base, inventory, currentPlan, map[string]model.Assertion{assertion.ID: assertion}, time.Unix(1, 0))
	if err != nil {
		t.Fatal(err)
	}
	if impact := impactByID(t, report, assertion.ID); !hasReason(impact, "catalog_definition_changed") || impact.ReuseAllowed {
		t.Fatalf("impact = %+v", impact)
	}
}

func TestGoHTTPDependenciesTrackOnlyNonTestGoSource(t *testing.T) {
	files := []model.FileRecord{
		{Path: "server.go", Size: 1},
		{Path: "server_test.go", Size: 1},
		{Path: "README.md", Size: 1},
	}
	before := testInventory("a", files)
	after := testInventory("b", files)
	for _, implementation := range []string{
		"prc.native.go-http-timeout@0.1",
		"prc.native.go-http-server-timeout@0.1",
	} {
		assertion := model.Assertion{ImplementationID: implementation}
		dependency := dependencies(assertion, before, after)
		if !dependency.paths["server.go"] || dependency.paths["server_test.go"] || dependency.paths["README.md"] || dependency.allFiles {
			t.Fatalf("%s dependencies = %+v", implementation, dependency)
		}
	}
}

func TestAnalyzeRejectsConfiguredToUnconfiguredComparison(t *testing.T) {
	inventory := testInventory("a", nil)
	assertion := testAssertion("README", []string{"README.md"})
	basePlan := testPlan(inventory.Digest, "a", assertion)
	basePlan.ProjectID = "project"
	currentPlan := testPlan(inventory.Digest, "a", assertion)
	base := model.RunResult{RunID: strings.Repeat("f", 64), Inventory: inventory, Plan: basePlan}
	if _, err := Analyze(base, inventory, currentPlan, map[string]model.Assertion{assertion.ID: assertion}, time.Unix(1, 0)); err == nil ||
		!strings.Contains(err.Error(), "does not match") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func testInventory(digest string, files []model.FileRecord) model.Inventory {
	return model.Inventory{SchemaVersion: model.InventorySchema, TargetName: "repository", Digest: strings.Repeat(digest[:1], 64), Files: files,
		PackageEcosystems: []string{}, Manifests: []string{}, LockFiles: []string{}, ContainerFiles: []string{}, Symlinks: []string{},
		CI: model.CIInventory{WorkflowFiles: []string{}}, Infrastructure: model.InfrastructureInventory{TerraformFiles: []string{}, KubernetesFiles: []string{}},
		Components: []model.InventoryComponent{}, Relations: []model.InventoryRelation{}, Facts: []model.InventoryFact{}}
}

func testAssertion(name string, paths []string) model.Assertion {
	return model.Assertion{ID: "PRC-A-" + name, Revision: 1, ImplementationID: "prc.native.file-present@0.1",
		Parameters: map[string]any{"paths": paths}}
}

func testPlan(inventoryDigest, digest string, assertions ...model.Assertion) model.Plan {
	plan := model.Plan{SchemaVersion: model.PlanSchema, Digest: strings.Repeat(digest[:1], 64), EngineVersion: "prc.engine/v0.1",
		TargetName: "repository", InventoryDigest: inventoryDigest, ProfileID: "prc/test", ProfileVersion: "1", ProfileDigest: strings.Repeat("d", 64),
		CatalogDigest:   strings.Repeat("c", 64),
		ArtifactDigests: []string{}, TargetEnvironments: []string{}, Assertions: []model.PlannedAssertion{}}
	for _, assertion := range assertions {
		plan.Assertions = append(plan.Assertions, model.PlannedAssertion{AssertionID: assertion.ID, AssertionRevision: assertion.Revision,
			DefinitionDigest: strings.Repeat(strings.ToLower(assertion.ID[len(assertion.ID)-1:]), 64), Implementation: assertion.ImplementationID,
			Applicability: "applicable", ApplicabilityBy: "cel-go/v0.30.0+prc-inventory/v0.3", ApplicabilityReason: "test"})
	}
	return plan
}

func impactByID(t *testing.T, report Report, id string) AssertionImpact {
	t.Helper()
	for _, impact := range report.Assertions {
		if impact.AssertionID == id {
			return impact
		}
	}
	t.Fatalf("impact %s not found", id)
	return AssertionImpact{}
}

func hasReason(impact AssertionImpact, code string) bool {
	for _, item := range impact.Reasons {
		if item.Code == code {
			return true
		}
	}
	return false
}
