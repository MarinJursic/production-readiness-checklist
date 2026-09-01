package report

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/MarinJursic/production-readiness-checklist/scanner/controlprogram"
	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
	"github.com/MarinJursic/production-readiness-checklist/scanner/trust"
)

func reportRun() model.RunResult {
	digest := strings.Repeat("a", 64)
	started := time.Date(2026, 8, 23, 10, 0, 0, 0, time.UTC)
	exactEvidence := controlprogram.Evidence{
		SchemaVersion: controlprogram.EvidenceSchemaVersion, EvidenceID: "report-fixture",
		ProgramSHA256: digest, ControlID: "USEQ-11111111", ControlRevision: 1,
		ControlSemanticSHA256: digest, ClauseID: digest, ClauseSHA256: digest,
		ImplementationContractSHA256: digest, SubjectID: "fixture", ObservedSubjects: []string{"fixture"},
		InventorySHA256: digest, Authority: controlprogram.AuthorityRepository, ObservedAt: started,
		Complete: true, Applicability: controlprogram.ApplicabilityApplicable,
		ApplicabilityProofContractSHA256: digest,
		Facts:                            map[string]controlprogram.Fact{"fixture.flag": {Type: controlprogram.FactBoolean, Complete: true, Boolean: boolPointer(true)}},
	}
	exactEvidenceSHA256 := controlprogram.EvidenceSHA256(exactEvidence)
	return model.RunResult{
		SchemaVersion: model.RunSchema, RunID: digest, StartedAt: started, CompletedAt: started.Add(time.Second),
		Plan: model.Plan{
			ProfileID: "prc/core-repository", ProfileVersion: "0.3",
			ConfigurationDigest: strings.Repeat("b", 64), ProjectID: "example-product",
			TargetEnvironments: []string{"staging"}, ArtifactDigests: []string{"sha256:" + digest},
		},
		Inventory: model.Inventory{TargetName: "<unsafe & target>", Digest: digest}, TerminalState: "no_go",
		AdapterExecutions: []model.AdapterExecution{{
			AdapterID: "<unsafe-adapter>", ManifestSHA256: digest, ExecutionID: digest,
			Resolution: model.AdapterResolution{Source: "explicit-local", PublisherID: "test", Trust: "local-explicit"},
			Transcript: model.AdapterTranscript{Summary: model.AdapterSummary{Status: "completed"}},
		}},
		Findings: []model.Finding{{
			SchemaVersion: model.FindingSchema, ID: digest, Fingerprint: strings.Repeat("f", 64),
			AssertionID: "PRC-A-CORE-001", ControlIDs: []string{"USEQ-FDCA6C71"},
			Title: "README present", Summary: "Missing README | required.", Severity: "high", Gate: "required",
			RemediationClass: "R2", Subject: model.FindingSubject{Kind: "project", ID: "example-product", InventoryDigest: digest},
			Locations: []model.FindingLocation{{Path: "README.md", Line: 1, Column: 2}}, EvidenceIDs: []string{"evidence-001"},
		}},
		Results: []model.AssertionResult{
			{AssertionID: "PRC-A-CORE-001", Applicability: "applicable", Assessment: "fail", Execution: "completed", Severity: "high", Gate: "required", Summary: "Missing README | required.", RemediationClass: "R2", ControlIDs: []string{"USEQ-FDCA6C71"}, Locations: []model.FindingLocation{{Path: "README.md", Line: 1, Column: 2}}, EvidenceRequired: []model.EvidenceRequirement{{Kind: "file", MinimumAuthority: "repository", Description: "A root README must exist."}}, EvidenceObserved: []model.Evidence{{ID: "evidence-001", Kind: "file", Source: "README.md", Summary: "README was not present."}}},
			{AssertionID: "PRC-A-CORE-012", Assessment: "manual_review", Execution: "completed", Severity: "high", Gate: "required", Summary: "Reviewer required.", RemediationClass: "R0", EvidenceObserved: []model.Evidence{}},
			{AssertionID: "PRC-A-CORE-013", Assessment: "unknown", Execution: "blocked", Severity: "high", Gate: "required", Summary: "Adapter unavailable.", RemediationClass: "R2", EvidenceObserved: []model.Evidence{}},
		},
		ControlCatalog: &model.ControlCatalogSummary{
			ControlCount: 2, ActiveControlCount: 2, AIReviewProvider: "codex",
			AIReviewState: "focused", AIReviewedCount: 1, AIAdvisoryFailCount: 1,
			ReviewedDeterministicCount: 1, ReviewedNondeterministicCount: 1,
			DeterministicBindingCount: 1, ClassificationCorpusSHA256: digest,
			ControlCheckBindingsSHA256: digest, DeterministicProgramTemplateCount: 1,
			DeterministicProgramExecutedCount: 1, DeterministicProgramPassCount: 1,
		},
		ControlResults: []model.ControlResult{
			{ControlID: "USEQ-11111111", Disposition: "verified_pass", Statement: "A broad rule has exact evidence.",
				Classification: "deterministic", ClassificationRoute: "local_static",
				ClassificationDecisionBasis: "strength_audit_confirmed", ClassificationRowSHA256: digest,
				DeterministicBindingID: "USEQ-11111111@1", DeterministicBindingSHA256: digest,
				DeterministicClauseResults: []model.DeterministicClauseResult{{
					TemplateID: "USEQ-11111111@1#1", CollectorID: "repository.fixture.v1", ClauseID: "c1",
					ClauseOrdinal: 1, ImplementationID: "repository.fixture.v1", ImplementationContractSHA256: digest,
					RequiredAuthority: "repository", ProviderID: "fixture", ProgramSHA256: digest,
					EvidenceSHA256: exactEvidenceSHA256, Status: "passed", Outcome: "pass", EvaluatedAt: started,
				}}},
			{ControlID: "USEQ-22222222", Disposition: "needs_review", Statement: "A cited rule remains advisory.",
				Classification: "nondeterministic", ClassificationRoute: "contextual_judgment",
				ClassificationDecisionBasis: "primary_nondeterministic", ClassificationRowSHA256: digest,
				AIReview: &model.AIControlReview{
					Provider: "codex", AssessmentCandidate: "advisory_fail_candidate", ApplicabilityCandidate: "applicable",
					ReviewDepth: "deep", Confidence: "medium", Priority: "high",
					RootCause: "The ownership boundary is not established.", RootCauseKey: "ownership-boundary-undefined",
					Effort: "medium", BlastRadius: "component",
					Reason: "The cited line looks risky.", Challenge: "The line might not describe live behavior.",
					RiskIfIgnored: "The behavior could fail in production.", Advice: "Review the real behavior.",
					RemediationSteps:  []string{"Correct the behavior in an isolated change."},
					VerificationSteps: []string{"Run the independent behavior check."},
					EvidenceNeeded:    []string{"Current runtime evidence."},
					Evidence:          []model.FindingLocation{{Path: "README.md", Line: 1}}, Limitations: []string{"Only repository text was visible."},
					CitationVerification: "snapshot_location_validated", ClaimVerification: "advisory_unverified", TaskID: digest,
				}},
		},
		DeterministicEvidence: []controlprogram.Evidence{exactEvidence},
		AuthoritativeEvidence: []model.AuthoritativeEvidenceVerification{{
			SchemaVersion: "prc.authoritative-evidence-verification/v0.1", BundleID: "fixture-bundle",
			BundleSHA256: digest, PolicySHA256: strings.Repeat("b", 64),
			CatalogSHA256: digest, InventorySHA256: digest,
			Authority: "repository", EntryCount: 1,
			PolicySignature: trust.Verification{
				KeyID: "policy-key", VerifiedAt: started.Add(time.Second), TrustStoreID: "fixture-store",
				TrustStoreDigest: digest,
			},
			EvidenceSignature: trust.Verification{
				KeyID: "evidence-key", VerifiedAt: started.Add(time.Second), TrustStoreID: "fixture-store",
				TrustStoreDigest: digest,
			},
		}},
		AIImprovementPlan: &model.AIImprovementPlan{
			SchemaVersion: "prc.ai-improvement-plan/v0.1", Authority: "advisory_only", SourceRunID: digest,
			ReviewProvider: "codex", ReviewDepth: "deep", ReviewState: "focused",
			ReviewedControlCount: 1, ItemCount: 1,
			Items: []model.AIImprovementPlanItem{{
				ItemID: digest, Domain: "engineering/governance-and-foundations",
				RootCauseKey: "ownership-boundary-undefined", RootCause: "The ownership boundary is not established.",
				Priority: "high", Effort: "medium", BlastRadius: "component",
				AssessmentCandidates: []string{"advisory_fail_candidate"}, ControlCount: 1,
				ControlIDs: []string{"USEQ-22222222"}, TaskIDs: []string{digest},
			}},
		},
	}
}

func boolPointer(value bool) *bool { return &value }

func TestMarkdownReportIsScopedAndEscapesTableCells(t *testing.T) {
	var output bytes.Buffer
	if err := Write("markdown", &output, reportRun()); err != nil {
		t.Fatal(err)
	}
	text := output.String()
	for _, expected := range []string{"# Production readiness assessment", "Missing README \\| required.", "## What the result means", "Verified problems", "Narrow checks passed", "Reviewed deterministic controls", "Exact programs attempted", "Retained replayable exact evidence documents", "AI review improvement plan", "The ownership boundary is not established.", "Signed authoritative evidence", "policy-key", "evidence-key", "Classification corpus digest", "USEQ-11111111@1", "strength_audit_confirmed", "citation=snapshot_location_validated; claim=advisory_unverified", "## Adapter executions", "local-explicit", "## Findings", "example-product", "staging", "not an unqualified production-readiness"} {
		if !strings.Contains(text, expected) {
			t.Errorf("missing %q in report", expected)
		}
	}
}

func TestSARIFContainsOnlyFailedFindings(t *testing.T) {
	var output bytes.Buffer
	if err := Write("sarif", &output, reportRun()); err != nil {
		t.Fatal(err)
	}
	var decoded map[string]any
	if err := json.Unmarshal(output.Bytes(), &decoded); err != nil {
		t.Fatal(err)
	}
	runs := decoded["runs"].([]any)
	properties := runs[0].(map[string]any)["properties"].(map[string]any)
	if properties["configuration_digest"] != strings.Repeat("b", 64) || properties["project_id"] != "example-product" {
		t.Fatalf("configured scope missing from SARIF: %+v", properties)
	}
	results := runs[0].(map[string]any)["results"].([]any)
	if len(results) != 1 || results[0].(map[string]any)["ruleId"] != "PRC-A-CORE-001" {
		t.Fatalf("unexpected SARIF results: %+v", results)
	}
	findingProperties := results[0].(map[string]any)["properties"].(map[string]any)
	if findingProperties["finding_id"] != strings.Repeat("a", 64) || findingProperties["fingerprint"] != strings.Repeat("f", 64) {
		t.Fatalf("canonical finding identity missing from SARIF: %+v", findingProperties)
	}
}

func TestJUnitDistinguishesFailureErrorAndSkipped(t *testing.T) {
	var output bytes.Buffer
	if err := Write("junit", &output, reportRun()); err != nil {
		t.Fatal(err)
	}
	var suite struct {
		Tests      int `xml:"tests,attr"`
		Failures   int `xml:"failures,attr"`
		Errors     int `xml:"errors,attr"`
		Skipped    int `xml:"skipped,attr"`
		Properties []struct {
			Name  string `xml:"name,attr"`
			Value string `xml:"value,attr"`
		} `xml:"properties>property"`
	}
	if err := xml.Unmarshal(output.Bytes(), &suite); err != nil {
		t.Fatal(err)
	}
	if suite.Tests != 3 || suite.Failures != 1 || suite.Errors != 1 || suite.Skipped != 1 {
		t.Fatalf("unexpected JUnit counts: %+v", suite)
	}
	foundConfiguration := false
	for _, property := range suite.Properties {
		if property.Name == "configuration_digest" && property.Value == strings.Repeat("b", 64) {
			foundConfiguration = true
		}
	}
	if !foundConfiguration {
		t.Fatalf("configured scope missing from JUnit: %+v", suite.Properties)
	}
}

func TestHTMLReportEscapesUntrustedText(t *testing.T) {
	var output bytes.Buffer
	if err := Write("html", &output, reportRun()); err != nil {
		t.Fatal(err)
	}
	text := output.String()
	if strings.Contains(text, "<unsafe & target>") || !strings.Contains(text, "&lt;unsafe &amp; target&gt;") ||
		strings.Contains(text, "<unsafe-adapter>") || !strings.Contains(text, "&lt;unsafe-adapter&gt;") {
		t.Fatalf("HTML output was not escaped: %s", text)
	}
	if !strings.Contains(text, "example-product") || !strings.Contains(text, "staging") {
		t.Fatal("configured scope missing from HTML")
	}
	for _, expected := range []string{"report-brand", "Scan report", "Local checks only", "hero-metrics", "About this score", "score-gauge", "pathLength=\"100\"", "NOT READY", "simple pass rate", "wider controls still need evidence or review", "What this scan actually checked", "does not certify the whole project", "Scores by category", "Not scored in this scan", "AI review improvement plan", "The ownership boundary is not established.", "advisory", "control-category", "Verified pass", "Reviewed classification", "strength_audit_confirmed", "Deterministic binding", "USEQ-11111111@1", "Exact deterministic programs", "repository.fixture.v1", "Exact deterministic execution", "replayable evidence documents retained", "Signed authoritative evidence", "fixture-bundle", "policy-key", "evidence-key", "Replayable exact evidence", "report-fixture", "fixture.flag", "Classification corpus digest", "AI verification state", "citation_verification", "advisory_unverified", "What to fix first", "Local check details", "README.md:1:2", "USEQ-FDCA6C71", "A root README must exist.", "README was not present.", "authority:", "observed: not recorded", "Evidence time:", "evidence-001", "Remediation class", ">R2<", "isolated agent-authored candidate", "Report-only scan", "Show more controls", "Technical evidence and IDs", "control-catalog-data", "pageSize = 25"} {
		if !strings.Contains(text, expected) {
			t.Errorf("detailed HTML report missing %q", expected)
		}
	}
	overallAt := strings.Index(text, "Local check pass rate")
	categoriesAt := strings.Index(text, "Scores by category")
	findingsAt := strings.Index(text, "What to fix first")
	detailsAt := strings.Index(text, "Local check details")
	if overallAt < 0 || categoriesAt <= overallAt || findingsAt <= categoriesAt || detailsAt <= findingsAt {
		t.Fatalf("report hierarchy should be overall score, category scores, findings, then details: %d %d %d %d", overallAt, categoriesAt, findingsAt, detailsAt)
	}
	if strings.Contains(text, `class="assertion-row result-fail" open`) {
		t.Fatal("failed assertions should start collapsed so large evidence lists do not break report navigation")
	}
	if strings.Contains(text, `class="finding-card severity-high" open`) || !strings.Contains(text, `class="finding-card severity-high"`) {
		t.Fatal("verified failure cards should start as compact, severity-colored disclosures")
	}
	for _, collapsed := range []string{`class="section-disclosure" id="local-checks"`, `class="section-disclosure" id="control-explorer"`, `class="section-disclosure" id="technical"`} {
		if !strings.Contains(text, collapsed) {
			t.Errorf("advanced report section is not an explicit collapsed disclosure: %s", collapsed)
		}
	}
}

func TestHTMLFindingsAndAttentionChecksSortBySeverity(t *testing.T) {
	run := reportRun()
	base := run.Findings[0]
	low := base
	low.Title, low.Severity, low.AssertionID = "Low issue", "low", "PRC-A-LOW"
	critical := base
	critical.Title, critical.Severity, critical.AssertionID = "Critical issue", "critical", "PRC-A-CRITICAL"
	medium := base
	medium.Title, medium.Severity, medium.AssertionID = "Medium issue", "medium", "PRC-A-MEDIUM"
	run.Findings = []model.Finding{low, medium, critical}
	run.Results = []model.AssertionResult{
		{AssertionID: "PRC-A-LOW", Assessment: "fail", Severity: "low"},
		{AssertionID: "PRC-A-CRITICAL", Assessment: "fail", Severity: "critical"},
		{AssertionID: "PRC-A-MEDIUM", Assessment: "manual_review", Severity: "medium"},
	}

	var output bytes.Buffer
	if err := Write("html", &output, run); err != nil {
		t.Fatal(err)
	}
	text := output.String()
	criticalAt := strings.Index(text, "Critical issue")
	mediumAt := strings.Index(text, "Medium issue")
	lowAt := strings.Index(text, "Low issue")
	if criticalAt < 0 || mediumAt <= criticalAt || lowAt <= mediumAt {
		t.Fatalf("findings should be ordered critical, medium, low: %d %d %d", criticalAt, mediumAt, lowAt)
	}
	localDetailsAt := strings.Index(text, `id="local-checks"`)
	if localDetailsAt < 0 {
		t.Fatal("local check details section is missing")
	}
	localDetails := text[localDetailsAt:]
	criticalCheckAt := strings.Index(localDetails, "PRC-A-CRITICAL")
	mediumCheckAt := strings.Index(localDetails, "PRC-A-MEDIUM")
	lowCheckAt := strings.Index(localDetails, "PRC-A-LOW")
	if criticalCheckAt < 0 || mediumCheckAt <= criticalCheckAt || lowCheckAt <= mediumCheckAt {
		t.Fatalf("attention checks should be ordered critical, medium, low: %d %d %d", criticalCheckAt, mediumCheckAt, lowCheckAt)
	}
	for _, expected := range []string{"severity-critical", "severity-medium", "severity-low", "Technical details"} {
		if !strings.Contains(text, expected) {
			t.Errorf("severity presentation missing %q", expected)
		}
	}
}

func TestBriefTextCollapsesWhitespaceAndPreservesUnicode(t *testing.T) {
	if got := briefText("  a\n  clearer   result  ", 100); got != "a clearer result" {
		t.Fatalf("unexpected short result: %q", got)
	}
	if got := briefText("readiness ✓ evidence", 11); got != "readiness ✓…" {
		t.Fatalf("unicode truncation was not rune-safe: %q", got)
	}
}

func TestHTMLReportKeepsLongDiagnosticsBehindNestedDetails(t *testing.T) {
	run := reportRun()
	locations := make([]model.FindingLocation, 0, 7)
	for index := 1; index <= 7; index++ {
		locations = append(locations, model.FindingLocation{Path: fmt.Sprintf("src/file-%d.go", index), Line: index})
	}
	run.Findings[0].Locations = locations
	run.Results[0].Locations = locations
	run.Results[0].Summary = strings.Repeat("long diagnostic output ", 12)

	var output bytes.Buffer
	if err := Write("html", &output, run); err != nil {
		t.Fatal(err)
	}
	text := output.String()
	if strings.Count(text, "Show all 7 locations") != 2 {
		t.Fatalf("long finding and assertion locations should each require a second disclosure")
	}
	if !strings.Contains(text, `class="raw-details"><summary>Technical evidence and IDs</summary>`) {
		t.Fatal("raw assertion evidence is not separated from the plain result")
	}
	if strings.Contains(text, `class="raw-details" open`) || strings.Contains(text, `class="finding-details" open`) {
		t.Fatal("technical evidence should start collapsed")
	}
	if !strings.Contains(text, "long diagnostic output long diagnostic output") || !strings.Contains(text, "…") {
		t.Fatal("the complete diagnostic and its shortened row label must both remain available")
	}
}

func TestHTMLReportKeepsLargeCatalogAsDataInsteadOfThousandsOfDOMNodes(t *testing.T) {
	run := reportRun()
	run.Findings = nil
	run.Results = nil
	run.AIImprovementPlan = nil
	run.DeterministicEvidence = nil
	run.AuthoritativeEvidence = nil
	run.ControlResults = make([]model.ControlResult, 0, 1000)
	for index := 0; index < 1000; index++ {
		run.ControlResults = append(run.ControlResults, model.ControlResult{
			ControlID:   fmt.Sprintf("USEQ-%08d", index),
			Statement:   fmt.Sprintf("Review production concern %d using evidence suitable for this project.", index),
			Disposition: "needs_review", Classification: "nondeterministic",
			ClassificationRoute: "contextual_judgment", Summary: "More evidence is needed.",
			Source: model.Source{Path: "docs/engineering/04-architecture-and-design.md", Line: index + 1},
		})
	}
	run.ControlCatalog.ControlCount = len(run.ControlResults)
	run.ControlCatalog.ActiveControlCount = len(run.ControlResults)

	var output bytes.Buffer
	if err := Write("html", &output, run); err != nil {
		t.Fatal(err)
	}
	text := output.String()
	if count := strings.Count(text, "<details"); count >= 100 {
		t.Fatalf("large catalog created %d disclosure nodes before the browser search was used", count)
	}
	if count := strings.Count(text, `class="control-row`); count != 0 {
		t.Fatalf("catalog controls were pre-rendered into %d control rows", count)
	}
	if !strings.Contains(text, `"id":"USEQ-00000000"`) || !strings.Contains(text, `"id":"USEQ-00000999"`) {
		t.Fatal("the complete catalog was not retained in the inert JSON data")
	}
	if output.Len() > 2*1024*1024 {
		t.Fatalf("1000-control report is unexpectedly large: %d bytes", output.Len())
	}
}

func TestLocalCheckSummaryKeepsPassRateSeparateFromGateResult(t *testing.T) {
	run := reportRun()
	run.Results = nil
	for index := 0; index < 14; index++ {
		run.Results = append(run.Results, model.AssertionResult{Assessment: "pass", Applicability: "applicable"})
	}
	for index := 0; index < 2; index++ {
		gate := "optional"
		if index == 0 {
			gate = "required"
		}
		run.Results = append(run.Results, model.AssertionResult{Assessment: "fail", Applicability: "applicable", Gate: gate})
		run.Results = append(run.Results, model.AssertionResult{Assessment: "not_applicable", Applicability: "not_applicable"})
	}
	run.ControlCatalog.ProfileTerminalState = "no_go"

	summary := SummarizeLocalChecks(run)
	if summary.Passed != 14 || summary.Applicable != 16 || summary.NotApplicable != 2 || summary.Blocking != 1 || summary.ToReview != 0 ||
		summary.Percentage != 88 || summary.Label != "NOT READY" || summary.Tone != "bad" {
		t.Fatalf("unexpected local summary: %+v", summary)
	}
	if !strings.Contains(summary.Explanation, "release-blocking") {
		t.Fatalf("gate explanation is unclear: %+v", summary)
	}
}

func TestLocalCheckSummaryCallsACompletePassingProfileGreat(t *testing.T) {
	run := reportRun()
	run.ControlCatalog = nil
	run.TerminalState = "profile_satisfied"
	run.Results = []model.AssertionResult{
		{Assessment: "pass", Applicability: "applicable"},
		{Assessment: "pass", Applicability: "applicable"},
		{Assessment: "not_applicable", Applicability: "not_applicable"},
	}
	summary := SummarizeLocalChecks(run)
	if summary.Percentage != 100 || summary.Label != "GREAT" || summary.Tone != "great" {
		t.Fatalf("unexpected passing summary: %+v", summary)
	}
}

func TestLocalCheckSummarySeparatesCompletedScanFromBlockedReadinessEvidence(t *testing.T) {
	run := reportRun()
	run.ControlCatalog = nil
	run.TerminalState = "environment_blocked"
	run.Results = []model.AssertionResult{{
		Assessment: "unknown", Applicability: "applicable", Execution: "blocked", Gate: "required",
	}}
	summary := SummarizeLocalChecks(run)
	if summary.Label != "READINESS BLOCKED" || summary.Tone != "bad" ||
		!strings.Contains(summary.Explanation, "local scan completed") ||
		!strings.Contains(summary.Explanation, "permitted evidence source") {
		t.Fatalf("completed scan and blocked readiness were conflated: %+v", summary)
	}
}

func TestCategoryScoresUseOnlyDistinctLinkedLocalChecks(t *testing.T) {
	run := reportRun()
	run.ControlResults = []model.ControlResult{
		{ControlID: "USEQ-SECURITY-A", Source: model.Source{Path: "docs/engineering/08-security-and-cryptography.md"}},
		{ControlID: "USEQ-SECURITY-B", Source: model.Source{Path: "docs/engineering/08-security-and-cryptography.md"}},
		{ControlID: "USEQ-DOCS", Source: model.Source{Path: "docs/engineering/13-documentation-and-knowledge.md"}},
		{ControlID: "USEQ-ARCH", Source: model.Source{Path: "docs/engineering/04-architecture-and-design.md"}},
	}
	run.Results = []model.AssertionResult{
		{Assessment: "pass", Applicability: "applicable", ControlIDs: []string{"USEQ-SECURITY-A", "USEQ-SECURITY-B"}},
		{Assessment: "fail", Applicability: "applicable", ControlIDs: []string{"USEQ-SECURITY-A"}},
		{Assessment: "pass", Applicability: "applicable", ControlIDs: []string{"USEQ-DOCS"}},
		{Assessment: "not_applicable", Applicability: "not_applicable", ControlIDs: []string{"USEQ-DOCS"}},
	}

	categories := summarizeControlCategories(run)
	byKey := map[string]controlCategorySummary{}
	for _, category := range categories {
		byKey[category.Key] = category
	}
	security := byKey["engineering-security-and-cryptography"]
	if security.LocalApplicable != 2 || security.LocalPassed != 1 || security.LocalFailed != 1 ||
		security.LocalPercentage != 50 || security.LocalLabel != "NEEDS WORK" || security.LocalTone != "bad" {
		t.Fatalf("unexpected security category: %+v", security)
	}
	documentation := byKey["engineering-documentation-and-knowledge"]
	if documentation.LocalApplicable != 1 || documentation.LocalPassed != 1 || documentation.LocalNotApplicable != 1 ||
		documentation.LocalPercentage != 100 || documentation.LocalLabel != "LIMITED EVIDENCE" || documentation.LocalTone != "review" {
		t.Fatalf("unexpected documentation category: %+v", documentation)
	}
	architecture := byKey["engineering-architecture-and-design"]
	if architecture.HasLocalScore || architecture.LocalLabel != "NOT SCORED" || architecture.LocalTone != "unscored" {
		t.Fatalf("unexecuted category was presented as scored: %+v", architecture)
	}
	if len(categories) != 3 || categories[0].Key != security.Key || categories[1].Key != documentation.Key || categories[2].Key != architecture.Key {
		t.Fatalf("categories were not ordered by useful score state: %+v", categories)
	}
}
