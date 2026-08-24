package report

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"strings"
	"testing"
	"time"

	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
)

func reportRun() model.RunResult {
	digest := strings.Repeat("a", 64)
	started := time.Date(2026, 8, 23, 10, 0, 0, 0, time.UTC)
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
		},
		ControlResults: []model.ControlResult{
			{ControlID: "USEQ-11111111", Disposition: "needs_review", Statement: "A broad rule needs more evidence."},
			{ControlID: "USEQ-22222222", Disposition: "needs_review", Statement: "A cited rule remains advisory.", AIReview: &model.AIControlReview{
				Provider: "codex", AssessmentCandidate: "advisory_fail_candidate", ApplicabilityCandidate: "applicable",
				Confidence: "medium", Reason: "The cited line looks risky.", Advice: "Review the real behavior.",
				Evidence: []model.FindingLocation{{Path: "README.md", Line: 1}}, Limitations: []string{"Only repository text was visible."},
				CitationVerification: "snapshot_location_validated", ClaimVerification: "advisory_unverified", TaskID: digest,
			}},
		},
	}
}

func TestMarkdownReportIsScopedAndEscapesTableCells(t *testing.T) {
	var output bytes.Buffer
	if err := Write("markdown", &output, reportRun()); err != nil {
		t.Fatal(err)
	}
	text := output.String()
	for _, expected := range []string{"# Production readiness assessment", "Missing README \\| required.", "## What the result means", "Verified problems", "Narrow checks passed", "citation=snapshot_location_validated; claim=advisory_unverified", "## Adapter executions", "local-explicit", "## Findings", "example-product", "staging", "not an unqualified production-readiness"} {
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
	for _, expected := range []string{"What the result means", "Verified problems", "Controls needing review or evidence", "AI verification state", "citation=snapshot_location_validated; claim=advisory_unverified", "does not prove that the line supports the AI claim", "Detailed findings", "All assertion results", "README.md:1:2", "USEQ-FDCA6C71", "A root README must exist.", "README was not present.", "authority:", "observed: not recorded", "Evidence time:", "evidence-001", "remediation class R2", "isolated agent-authored candidate", "Report-only scan"} {
		if !strings.Contains(text, expected) {
			t.Errorf("detailed HTML report missing %q", expected)
		}
	}
}
