package engine

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/MarinJursic/production-readiness-checklist/scanner/catalog"
	projectconfig "github.com/MarinJursic/production-readiness-checklist/scanner/config"
	"github.com/MarinJursic/production-readiness-checklist/scanner/finding"
	"github.com/MarinJursic/production-readiness-checklist/scanner/inventory"
	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
)

const pinnedCheckout = "actions/checkout@3d3c42e5aac5ba805825da76410c181273ba90b1"

func repositoryRoot(t *testing.T) string {
	t.Helper()
	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	return root
}

func writeFixture(t *testing.T, root, relative, content string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(relative))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func healthyRepository(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	files := map[string]string{
		".git/HEAD":                "ref: refs/heads/main\n",
		".git/refs/heads/main":     strings.Repeat("a", 40) + "\n",
		"README.md":                "# Example\n",
		"LICENSE":                  "MIT\n",
		"CONTRIBUTING.md":          "# Contributing\n",
		"SECURITY.md":              "# Security\n",
		"CODE_OF_CONDUCT.md":       "# Conduct\n",
		".github/CODEOWNERS":       "* @owner\n",
		".github/dependabot.yml":   "version: 2\nupdates: []\n",
		".python-version":          "3.12\n",
		"requirements.txt":         "example==1.0\n",
		"requirements.lock.txt":    "example==1.0\n",
		"app.py":                   "def ready(): return True\n",
		"tests/test_app.py":        "def test_ready(): assert True\n",
		".github/workflows/ci.yml": "name: CI\non: [push]\npermissions:\n  contents: read\njobs:\n  test:\n    runs-on: ubuntu-latest\n    timeout-minutes: 15\n    steps:\n      - uses: " + pinnedCheckout + "\n",
	}
	for relative, content := range files {
		writeFixture(t, root, relative, content)
	}
	return root
}

func scanner(t *testing.T) *Engine {
	t.Helper()
	c, err := catalog.Load(repositoryRoot(t))
	if err != nil {
		t.Fatal(err)
	}
	result := New(c)
	times := []time.Time{
		time.Date(2026, 8, 23, 10, 0, 0, 0, time.UTC),
		time.Date(2026, 8, 23, 10, 0, 1, 0, time.UTC),
	}
	index := 0
	result.Now = func() time.Time {
		value := times[index%len(times)]
		index++
		return value
	}
	return result
}

func findResult(t *testing.T, run model.RunResult, assertionID string) model.AssertionResult {
	t.Helper()
	for _, result := range run.Results {
		if result.AssertionID == assertionID {
			return result
		}
	}
	t.Fatalf("missing result %s", assertionID)
	return model.AssertionResult{}
}

func TestHealthyRepositoryProducesVerifiedAndExplicitUnresolvedResults(t *testing.T) {
	root := healthyRepository(t)
	item, err := inventory.Build(root)
	if err != nil {
		t.Fatal(err)
	}
	run, err := scanner(t).Scan("prc/core-repository", item)
	if err != nil {
		t.Fatal(err)
	}
	for index := 1; index <= 11; index++ {
		assertionID := fmt.Sprintf("PRC-A-CORE-%03d", index)
		if result := findResult(t, run, assertionID); result.Assessment != "pass" {
			t.Errorf("%s = %s: %s", assertionID, result.Assessment, result.Summary)
		}
	}
	if got := findResult(t, run, "PRC-A-CORE-014").Assessment; got != "pass" {
		t.Fatalf("final-newline assertion = %s", got)
	}
	for index := 15; index <= 25; index++ {
		assertionID := fmt.Sprintf("PRC-A-CORE-%03d", index)
		result := findResult(t, run, assertionID)
		if result.Assessment != "pass" {
			t.Errorf("%s = %s: %s", assertionID, result.Assessment, result.Summary)
		}
		if len(result.EvidenceObserved) == 0 {
			t.Errorf("%s passed without observed evidence", assertionID)
		}
	}
	for index := 26; index <= 30; index++ {
		assertionID := fmt.Sprintf("PRC-A-CORE-%03d", index)
		if result := findResult(t, run, assertionID); result.Assessment != "not_applicable" {
			t.Errorf("%s = %s: %s", assertionID, result.Assessment, result.Summary)
		}
	}
	if got := findResult(t, run, "PRC-A-CORE-012").Assessment; got != "manual_review" {
		t.Fatalf("manual assertion = %s", got)
	}
	if got := findResult(t, run, "PRC-A-CORE-013").Execution; got != "blocked" {
		t.Fatalf("analysis assertion execution = %s", got)
	}
	if run.TerminalState != "environment_blocked" {
		t.Fatalf("terminal state = %s", run.TerminalState)
	}
	for _, result := range run.Results {
		if result.EvidenceObserved == nil {
			t.Fatalf("%s encoded evidence_observed as null", result.AssertionID)
		}
	}
	runIdentity := run
	runIdentity.RunID = ""
	runPayload, err := json.Marshal(runIdentity)
	if err != nil {
		t.Fatal(err)
	}
	runDigest := sha256.Sum256(runPayload)
	if run.RunID != hex.EncodeToString(runDigest[:]) {
		t.Fatal("run ID does not address the canonical run envelope")
	}
	firstEvidence := findResult(t, run, "PRC-A-CORE-001").EvidenceObserved[0]
	evidenceIdentity := firstEvidence
	evidenceIdentity.ID = ""
	evidencePayload, err := json.Marshal(evidenceIdentity)
	if err != nil {
		t.Fatal(err)
	}
	evidenceDigest := sha256.Sum256(evidencePayload)
	if firstEvidence.ID != hex.EncodeToString(evidenceDigest[:]) {
		t.Fatal("evidence ID does not address the canonical evidence envelope")
	}
}

func TestMutableActionReferenceTriggersNoGo(t *testing.T) {
	root := healthyRepository(t)
	workflow := filepath.Join(root, ".github", "workflows", "ci.yml")
	data, err := os.ReadFile(workflow)
	if err != nil {
		t.Fatal(err)
	}
	writeFixture(t, root, ".github/workflows/ci.yml", strings.ReplaceAll(string(data), pinnedCheckout, "actions/checkout@v7"))
	item, err := inventory.Build(root)
	if err != nil {
		t.Fatal(err)
	}
	run, err := scanner(t).Scan("prc/core-repository", item)
	if err != nil {
		t.Fatal(err)
	}
	result := findResult(t, run, "PRC-A-CORE-008")
	if result.Assessment != "fail" || run.TerminalState != "no_go" {
		t.Fatalf("action result=%s terminal=%s summary=%s", result.Assessment, run.TerminalState, result.Summary)
	}
	if len(run.Findings) != 1 {
		t.Fatalf("findings = %+v", run.Findings)
	}
	findingItem := run.Findings[0]
	if findingItem.AssertionID != result.AssertionID || findingItem.Subject.InventoryDigest != run.Inventory.Digest ||
		len(findingItem.ID) != 64 || len(findingItem.Fingerprint) != 64 ||
		!slices.Contains(findingItem.Locations, model.FindingLocation{Path: ".github/workflows/ci.yml"}) {
		t.Fatalf("canonical finding = %+v", findingItem)
	}
	if err := finding.Validate(findingItem); err != nil {
		t.Fatal(err)
	}
}

func TestFindingFingerprintIsStableAcrossUnrelatedInventoryChanges(t *testing.T) {
	root := healthyRepository(t)
	workflow := filepath.Join(root, ".github", "workflows", "ci.yml")
	data, err := os.ReadFile(workflow)
	if err != nil {
		t.Fatal(err)
	}
	writeFixture(t, root, ".github/workflows/ci.yml", strings.ReplaceAll(string(data), pinnedCheckout, "actions/checkout@v7"))
	firstInventory, err := inventory.Build(root)
	if err != nil {
		t.Fatal(err)
	}
	first, err := scanner(t).Scan("prc/core-repository", firstInventory)
	if err != nil {
		t.Fatal(err)
	}
	writeFixture(t, root, "unrelated.txt", "unrelated\n")
	secondInventory, err := inventory.Build(root)
	if err != nil {
		t.Fatal(err)
	}
	second, err := scanner(t).Scan("prc/core-repository", secondInventory)
	if err != nil {
		t.Fatal(err)
	}
	if len(first.Findings) != 1 || len(second.Findings) != 1 ||
		first.Findings[0].Fingerprint != second.Findings[0].Fingerprint ||
		first.Findings[0].ID == second.Findings[0].ID {
		t.Fatalf("finding correlation changed incorrectly: %+v %+v", first.Findings, second.Findings)
	}
}

func TestMissingFinalNewlineProducesR1Finding(t *testing.T) {
	root := healthyRepository(t)
	writeFixture(t, root, "app.py", "def ready(): return True")
	item, err := inventory.Build(root)
	if err != nil {
		t.Fatal(err)
	}
	run, err := scanner(t).Scan("prc/core-repository", item)
	if err != nil {
		t.Fatal(err)
	}
	result := findResult(t, run, "PRC-A-CORE-014")
	if result.Assessment != "fail" || result.RemediationClass != "R1" || !strings.Contains(result.Summary, "app.py") {
		t.Fatalf("final-newline result=%+v", result)
	}
}

func TestTargetMutationAfterInventoryCannotProducePass(t *testing.T) {
	root := healthyRepository(t)
	item, err := inventory.Build(root)
	if err != nil {
		t.Fatal(err)
	}
	writeFixture(t, root, "README.md", "# Changed after inventory\n")
	run, err := scanner(t).Scan("prc/core-repository", item)
	if err != nil {
		t.Fatal(err)
	}
	result := findResult(t, run, "PRC-A-CORE-001")
	if result.Execution != "error" || result.Assessment == "pass" {
		t.Fatalf("mutation result=%+v", result)
	}
}

func TestPlanDigestIsDeterministic(t *testing.T) {
	item, err := inventory.Build(healthyRepository(t))
	if err != nil {
		t.Fatal(err)
	}
	scanner := scanner(t)
	first, err := scanner.Plan("prc/core-repository", item)
	if err != nil {
		t.Fatal(err)
	}
	second, err := scanner.Plan("prc/core-repository", item)
	if err != nil {
		t.Fatal(err)
	}
	if first.Digest != second.Digest {
		t.Fatalf("plan digests differ: %s != %s", first.Digest, second.Digest)
	}
	for _, planned := range first.Assertions {
		if planned.ApplicabilityReason == "" || planned.AssertionRevision < 1 || len(planned.DefinitionDigest) != 64 {
			t.Fatalf("%s has incomplete rule binding: %+v", planned.AssertionID, planned)
		}
		if planned.Applicability == "undetermined" {
			t.Fatalf("checked-in applicability for %s is undetermined: %s", planned.AssertionID, planned.ApplicabilityReason)
		}
	}
	if first.EngineVersion != "prc.engine/v0.1" || len(first.ProfileDigest) != 64 || len(first.CatalogDigest) != 64 {
		t.Fatalf("plan has incomplete engine/catalog/profile binding: %+v", first)
	}
	if first.ExecutionMode != ExecutionModeInspect || first.CapabilityPolicy.Process != "deny" ||
		first.CapabilityPolicy.WriteScratch || len(first.Nodes) != len(first.Assertions)+len(first.Adapters)+2 {
		t.Fatalf("plan has an invalid inspect-mode DAG: %+v", first)
	}
	if first.Nodes[0].Kind != "inventory" || first.Nodes[len(first.Nodes)-1].Kind != "gate" {
		t.Fatalf("plan DAG boundaries are invalid: %+v", first.Nodes)
	}
	for _, assertion := range first.Assertions {
		if _, registered := implementationRegistry[assertion.Implementation]; !registered {
			t.Errorf("catalog implementation is not registered: %s", assertion.Implementation)
		}
	}
	filePresentUses := 0
	for _, implementation := range first.Implementations {
		if implementation.ID == "prc.native.file-present@0.2" {
			filePresentUses = len(implementation.AssertionIDs)
		}
	}
	if filePresentUses != 7 {
		t.Fatalf("shared implementation registry entry was not deduplicated: %d", filePresentUses)
	}
}

func TestWhitespaceOnlyOrInvalidFoundationTextCannotPass(t *testing.T) {
	for _, test := range []struct {
		name    string
		content []byte
	}{
		{"ASCII whitespace", []byte(" \t\r\n")},
		{"Unicode whitespace and BOM", []byte("\ufeff\u00a0\u2003\n")},
		{"invalid UTF-8", []byte{0xff, 0xfe}},
		{"NUL control", []byte{'o', 'k', 0}},
	} {
		t.Run(test.name, func(t *testing.T) {
			root := healthyRepository(t)
			if err := os.WriteFile(filepath.Join(root, "README.md"), test.content, 0o644); err != nil {
				t.Fatal(err)
			}
			item, err := inventory.Build(root)
			if err != nil {
				t.Fatal(err)
			}
			run, err := scanner(t).Scan("prc/core-repository", item)
			if err != nil {
				t.Fatal(err)
			}
			result := findResult(t, run, "PRC-A-CORE-001")
			if result.Execution != "completed" || result.Assessment != "fail" ||
				!strings.Contains(result.Summary, "non-whitespace UTF-8 text") {
				t.Fatalf("invalid foundation text passed: %+v", result)
			}
		})
	}
}

func TestReadableFoundationTextIsLanguageNeutral(t *testing.T) {
	root := healthyRepository(t)
	if err := os.WriteFile(filepath.Join(root, "README.md"), []byte("\ufeffΈτοιμο για παραγωγή\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	item, err := inventory.Build(root)
	if err != nil {
		t.Fatal(err)
	}
	run, err := scanner(t).Scan("prc/core-repository", item)
	if err != nil {
		t.Fatal(err)
	}
	result := findResult(t, run, "PRC-A-CORE-001")
	if result.Execution != "completed" || result.Assessment != "pass" {
		t.Fatalf("readable non-English UTF-8 text did not pass: %+v", result)
	}
}

func TestPlanRejectsUnsupportedExecutionMode(t *testing.T) {
	item, err := inventory.Build(healthyRepository(t))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := scanner(t).PlanMode("prc/core-repository", item, "production-write"); err == nil ||
		!strings.Contains(err.Error(), "unsupported execution mode") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExecutionGraphRejectsForwardDependencies(t *testing.T) {
	item, err := inventory.Build(healthyRepository(t))
	if err != nil {
		t.Fatal(err)
	}
	plan, err := scanner(t).Plan("prc/core-repository", item)
	if err != nil {
		t.Fatal(err)
	}
	plan.Nodes[1].DependsOn = []string{plan.Nodes[len(plan.Nodes)-1].ID}
	if err := ValidateExecutionPlan(plan); err == nil || !strings.Contains(err.Error(), "later node") {
		t.Fatalf("unexpected graph validation error: %v", err)
	}
}

func TestExecutionGraphRejectsIncompleteTopologyAndPolicyDrift(t *testing.T) {
	item, err := inventory.Build(healthyRepository(t))
	if err != nil {
		t.Fatal(err)
	}
	plan, err := scanner(t).Plan("prc/core-repository", item)
	if err != nil {
		t.Fatal(err)
	}
	withoutGate := plan
	withoutGate.Nodes = append([]model.PlanNode(nil), plan.Nodes[:len(plan.Nodes)-1]...)
	if err := ValidateExecutionPlan(withoutGate); err == nil || !strings.Contains(err.Error(), "one gate node") {
		t.Fatalf("unexpected incomplete topology error: %v", err)
	}
	driftedPolicy := plan
	driftedPolicy.CapabilityPolicy.Network = "allow"
	if err := ValidateExecutionPlan(driftedPolicy); err == nil || !strings.Contains(err.Error(), "does not match mode") {
		t.Fatalf("unexpected policy drift error: %v", err)
	}
	duplicateImplementation := plan
	duplicateImplementation.Implementations = append(
		append([]model.PlannedImplementation(nil), plan.Implementations...),
		plan.Implementations[0],
	)
	if err := ValidateExecutionPlan(duplicateImplementation); err == nil || !strings.Contains(err.Error(), "duplicate implementation") {
		t.Fatalf("unexpected registry duplication error: %v", err)
	}
}

func TestUnregisteredImplementationIsBlockedWithoutDispatch(t *testing.T) {
	item, err := inventory.Build(healthyRepository(t))
	if err != nil {
		t.Fatal(err)
	}
	engine := scanner(t)
	assertion := engine.Catalog.Assertions["PRC-A-CORE-001"]
	assertion.ImplementationID = "prc.native.unregistered@0.1"
	engine.Catalog.Assertions[assertion.ID] = assertion
	run, err := engine.Scan("prc/core-repository", item)
	if err != nil {
		t.Fatal(err)
	}
	result := findResult(t, run, assertion.ID)
	if result.Execution != "blocked" || result.Assessment != "unknown" ||
		!strings.Contains(result.Summary, "No scanner implementation is registered") {
		t.Fatalf("unregistered implementation result = %+v", result)
	}
}

func TestPlanDigestBindsCompleteRuleAndProfileDefinitions(t *testing.T) {
	item, err := inventory.Build(healthyRepository(t))
	if err != nil {
		t.Fatal(err)
	}
	engine := scanner(t)
	base, err := engine.Plan("prc/core-repository", item)
	if err != nil {
		t.Fatal(err)
	}
	assertion := engine.Catalog.Assertions["PRC-A-CORE-001"]
	assertion.Parameters = map[string]any{"paths": []string{"README.md"}, "minimum_bytes": 2}
	engine.Catalog.Assertions[assertion.ID] = assertion
	changedRule, err := engine.Plan("prc/core-repository", item)
	if err != nil {
		t.Fatal(err)
	}
	if changedRule.Digest == base.Digest || changedRule.Assertions[0].DefinitionDigest == base.Assertions[0].DefinitionDigest {
		t.Fatal("assertion parameter change did not change the bound plan identity")
	}
	profile := engine.Catalog.Profiles["prc/core-repository"]
	profile.Description += " Changed."
	engine.Catalog.Profiles[profile.ID] = profile
	changedProfile, err := engine.Plan("prc/core-repository", item)
	if err != nil {
		t.Fatal(err)
	}
	if changedProfile.Digest == changedRule.Digest || changedProfile.ProfileDigest == changedRule.ProfileDigest {
		t.Fatal("profile definition change did not change the bound plan identity")
	}
	objective := engine.Catalog.Objectives[assertion.ControlIDs[0]]
	objective.Statement += " Changed."
	engine.Catalog.Objectives[objective.ID] = objective
	changedCatalog, err := engine.Plan("prc/core-repository", item)
	if err != nil {
		t.Fatal(err)
	}
	if changedCatalog.Digest == changedProfile.Digest || changedCatalog.CatalogDigest == changedProfile.CatalogDigest {
		t.Fatal("objective change did not change the bound catalog and plan identities")
	}
}

func TestCELApplicabilityUsesBoundedInventoryProjection(t *testing.T) {
	item, err := inventory.Build(healthyRepository(t))
	if err != nil {
		t.Fatal(err)
	}
	engine := scanner(t)
	tests := []struct {
		name       string
		expression string
		want       string
	}{
		{
			name:       "composes source and ecosystem facts",
			expression: `inventory.source_files > 0 && "python" in inventory.package_ecosystems`,
			want:       "applicable",
		},
		{
			name:       "component macros are supported",
			expression: `inventory.components.exists(c, c.kind == "package-manifest")`,
			want:       "applicable",
		},
		{
			name:       "negative expressions are explicit",
			expression: `inventory.infrastructure.kubernetes_files.size() > 0`,
			want:       "not_applicable",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, reason := engine.evaluateApplicability(test.expression, item)
			if got != test.want {
				t.Fatalf("applicability = %s, want %s: %s", got, test.want, reason)
			}
			if reason == "" {
				t.Fatal("applicability result has no reason")
			}
		})
	}
}

func TestCELApplicabilityFailsClosed(t *testing.T) {
	item := model.Inventory{
		SourceFiles: 1,
		Components:  []model.InventoryComponent{{ID: "repository:.", Kind: "repository", Path: "."}},
	}
	engine := scanner(t)
	tests := []struct {
		name       string
		expression string
	}{
		{name: "missing field", expression: `inventory.not_a_real_field == true`},
		{name: "non boolean result", expression: `inventory.source_files`},
		{name: "empty", expression: "   "},
		{name: "syntax error", expression: `inventory.source_files >`},
		{name: "expression size limit", expression: strings.Repeat("true || ", 600) + "false"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, reason := engine.evaluateApplicability(test.expression, item)
			if got != "undetermined" {
				t.Fatalf("applicability = %s, want undetermined", got)
			}
			if reason == "" {
				t.Fatal("fail-closed result has no reason")
			}
		})
	}
}

func TestCELApplicabilityEnforcesRuntimeCostLimit(t *testing.T) {
	components := make([]model.InventoryComponent, 20_000)
	for index := range components {
		components[index] = model.InventoryComponent{
			ID: fmt.Sprintf("component:%05d", index), Kind: "source", Path: fmt.Sprintf("src/%05d", index),
		}
	}
	item := model.Inventory{Components: components}
	got, reason := scanner(t).evaluateApplicability(
		`inventory.components.all(c, c.kind == "source")`, item,
	)
	if got != "undetermined" {
		t.Fatalf("cost-exhausting applicability = %s, want undetermined", got)
	}
	if !strings.Contains(strings.ToLower(reason), "cost") {
		t.Fatalf("cost-exhaustion reason = %q", reason)
	}
}

func TestCELApplicabilityProgramCacheIsConcurrent(t *testing.T) {
	engine := scanner(t)
	item := model.Inventory{SourceFiles: 1}
	const workers = 32
	var wait sync.WaitGroup
	errors := make(chan string, workers)
	for range workers {
		wait.Add(1)
		go func() {
			defer wait.Done()
			got, reason := engine.evaluateApplicability(`inventory.source_files > 0`, item)
			if got != "applicable" {
				errors <- fmt.Sprintf("result=%s reason=%s", got, reason)
			}
		}()
	}
	wait.Wait()
	close(errors)
	for failure := range errors {
		t.Error(failure)
	}
}

func TestConfiguredPlanBindsScopeAndExposesDeclarationsToCEL(t *testing.T) {
	item, err := inventory.Build(healthyRepository(t))
	if err != nil {
		t.Fatal(err)
	}
	configPath := filepath.Join(repositoryRoot(t), "fixtures", "config", "production-readiness.yaml")
	validation, err := projectconfig.Load(configPath)
	if err != nil {
		t.Fatal(err)
	}
	item, err = inventory.BindConfiguration(item, validation, configPath)
	if err != nil {
		t.Fatal(err)
	}
	engine := scanner(t)
	plan, err := engine.Plan("prc/core-repository", item)
	if err != nil {
		t.Fatal(err)
	}
	if plan.ConfigurationDigest != validation.Digest || plan.ProjectID != "example-product" ||
		!slices.Equal(plan.TargetEnvironments, []string{"staging"}) || plan.ArtifactDigests == nil {
		t.Fatalf("configured plan lost declared identity: %+v", plan)
	}
	got, reason := engine.evaluateApplicability(`inventory.declared.features.authentication == true`, item)
	if got != "applicable" {
		t.Fatalf("declared feature applicability = %s: %s", got, reason)
	}
	item.DeclaredScope.ProfileID = "prc/other-profile"
	if _, err := engine.Plan("prc/core-repository", item); err == nil {
		t.Fatalf("unexpected profile mismatch error: %v", err)
	}
}

func TestSymlinkEscapeCannotSatisfyFileAssertion(t *testing.T) {
	root := healthyRepository(t)
	if err := os.Remove(filepath.Join(root, "README.md")); err != nil {
		t.Fatal(err)
	}
	outside := filepath.Join(t.TempDir(), "README.md")
	if err := os.WriteFile(outside, []byte("outside\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(root, "README.md")); err != nil {
		t.Fatal(err)
	}
	item, err := inventory.Build(root)
	if err != nil {
		t.Fatal(err)
	}
	run, err := scanner(t).Scan("prc/core-repository", item)
	if err != nil {
		t.Fatal(err)
	}
	result := findResult(t, run, "PRC-A-CORE-001")
	if result.Assessment == "pass" {
		t.Fatalf("symlink escape produced pass: %+v", result)
	}
}

func TestTerminalStateHonorsGateSemantics(t *testing.T) {
	profile := model.Profile{TerminalPolicy: model.TerminalPolicy{
		BlockOn: []string{"critical", "high"}, AllowManualRemaining: true,
	}}
	tests := []struct {
		name    string
		results []model.AssertionResult
		want    string
	}{
		{
			name: "advisory failures and unavailable evidence remain visible without blocking",
			results: []model.AssertionResult{
				{Applicability: "applicable", Execution: "completed", Assessment: "fail", Severity: "high", Gate: "advisory"},
				{Applicability: "applicable", Execution: "blocked", Assessment: "unknown", Severity: "high", Gate: "advisory"},
			},
			want: "profile_satisfied",
		},
		{
			name: "required blocking severity fails closed",
			results: []model.AssertionResult{
				{Applicability: "applicable", Execution: "completed", Assessment: "fail", Severity: "high", Gate: "required"},
			},
			want: "no_go",
		},
		{
			name: "no-go gate blocks at any severity",
			results: []model.AssertionResult{
				{Applicability: "applicable", Execution: "completed", Assessment: "fail", Severity: "low", Gate: "no-go"},
			},
			want: "no_go",
		},
		{
			name: "required nonblocking failure is incomplete",
			results: []model.AssertionResult{
				{Applicability: "applicable", Execution: "completed", Assessment: "fail", Severity: "medium", Gate: "required"},
			},
			want: "assessment_incomplete",
		},
		{
			name: "required execution failure is environment blocked",
			results: []model.AssertionResult{
				{Applicability: "applicable", Execution: "error", Assessment: "unknown", Severity: "medium", Gate: "required"},
			},
			want: "environment_blocked",
		},
		{
			name: "undetermined required applicability is incomplete",
			results: []model.AssertionResult{
				{Applicability: "undetermined", Execution: "not_run", Assessment: "unknown", Severity: "medium", Gate: "required"},
			},
			want: "assessment_incomplete",
		},
		{
			name: "allowed required manual evidence has an explicit state",
			results: []model.AssertionResult{
				{Applicability: "applicable", Execution: "completed", Assessment: "manual_review", Severity: "medium", Gate: "required"},
			},
			want: "machine_work_complete_manual_evidence_remaining",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := terminalState(profile, test.results); got != test.want {
				t.Fatalf("terminal state = %s, want %s", got, test.want)
			}
		})
	}

	profile.TerminalPolicy.AllowManualRemaining = false
	manual := []model.AssertionResult{{
		Applicability: "applicable", Execution: "completed", Assessment: "manual_review",
		Severity: "medium", Gate: "required",
	}}
	if got := terminalState(profile, manual); got != "assessment_incomplete" {
		t.Fatalf("disallowed manual evidence terminal state = %s", got)
	}
}
