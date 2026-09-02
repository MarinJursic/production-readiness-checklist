package main

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/MarinJursic/production-readiness-checklist/scanner/controlreview"
	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
)

func TestTerminalTextNeutralizesControlAndBidirectionalCharacters(t *testing.T) {
	input := "safe\x1b]8;;https://attacker.invalid\aFAKE PASS\x1b]8;;\a\r\nnext\b\u202ehidden"
	output := terminalText(input)
	if strings.ContainsAny(output, "\x1b\a\r\n\b") || strings.ContainsRune(output, '\u202e') ||
		!strings.Contains(output, `\u001B`) || !strings.Contains(output, `\u202E`) {
		t.Fatalf("unsafe terminal text: %q", output)
	}
}

func TestTerminalBriefIsSafeAndRuneAware(t *testing.T) {
	if got := terminalBrief("  a\n clearer   result  ", 100); got != "a clearer result" {
		t.Fatalf("unexpected compact result: %q", got)
	}
	if got := terminalBrief("readiness ✓ evidence", 11); got != "readiness ✓…" {
		t.Fatalf("unicode truncation was not rune-safe: %q", got)
	}
	if got := terminalBrief("safe\x1b[31m fake", 100); strings.ContainsRune(got, '\x1b') || !strings.Contains(got, `\u001B`) {
		t.Fatalf("terminal control was not neutralized: %q", got)
	}
}

func TestTerminalColorPolicy(t *testing.T) {
	t.Setenv("NO_COLOR", "")
	t.Setenv("TERM", "xterm-256color")
	var output bytes.Buffer
	if newTerminalStyle("auto", &output).color || newTerminalStyle("never", &output).color {
		t.Fatal("redirected or never output unexpectedly enabled color")
	}
	if !newTerminalStyle("always", &output).color {
		t.Fatal("explicit color did not enable color")
	}
	t.Setenv("NO_COLOR", "1")
	if newTerminalStyle("always", &output).color {
		t.Fatal("NO_COLOR did not disable explicit terminal color")
	}
}

func TestTerminalFileLinkIsClickableAndEscapesTheTarget(t *testing.T) {
	path := filepath.Join(t.TempDir(), "readiness report #1.html")
	plain, linked := (terminalStyle{}).fileLink(path)
	if linked || plain != path || strings.ContainsRune(plain, '\x1b') {
		t.Fatalf("plain report path was changed: linked=%t value=%q", linked, plain)
	}

	clickable, linked := (terminalStyle{hyperlink: true}).fileLink(path)
	if !linked || !strings.Contains(clickable, "\x1b]8;;file://") ||
		!strings.Contains(clickable, "%20") || !strings.Contains(clickable, "%231.html") ||
		!strings.Contains(clickable, path) || !strings.HasSuffix(clickable, "\x1b]8;;\x1b\\") {
		t.Fatalf("report path is not a safe OSC 8 file link: %q", clickable)
	}
}

func TestTerminalFileLinkNeutralizesControlCharacters(t *testing.T) {
	path := filepath.Join(t.TempDir(), "safe\x1b]8;;https://attacker.invalid\a.html")
	clickable, linked := (terminalStyle{hyperlink: true}).fileLink(path)
	if !linked || strings.Count(clickable, "\x1b") != 4 || strings.ContainsRune(clickable, '\a') ||
		!strings.Contains(clickable, `\u001B`) || !strings.Contains(clickable, `\u0007`) ||
		!strings.Contains(clickable, "%1B") || !strings.Contains(clickable, "%07") {
		t.Fatalf("unsafe report link: %q", clickable)
	}
}

func TestAssessmentLabelsKeepPlainMeaningWithAndWithoutColor(t *testing.T) {
	cases := []struct{ assessment, execution, word string }{
		{"pass", "completed", "PASS"}, {"fail", "completed", "FAIL"},
		{"unknown", "blocked", "BLOCKED"}, {"manual_review", "completed", "MANUAL"},
		{"not_applicable", "completed", "N/A"}, {"unknown", "error", "ERROR"},
	}
	for _, item := range cases {
		plain := assessmentLabel(item.assessment, item.execution, terminalStyle{})
		colored := assessmentLabel(item.assessment, item.execution, terminalStyle{color: true})
		if !strings.Contains(plain, item.word) || strings.Contains(plain, "\x1b[") ||
			!strings.Contains(colored, item.word) || !strings.Contains(colored, "\x1b[") {
			t.Fatalf("assessment %s/%s plain=%q colored=%q", item.assessment, item.execution, plain, colored)
		}
	}
}

func TestProductBannerWorksWithAndWithoutColor(t *testing.T) {
	for name, style := range map[string]terminalStyle{
		"plain": {}, "colored": {color: true},
	} {
		t.Run(name, func(t *testing.T) {
			var output bytes.Buffer
			printProductBanner(&output, style)
			text := output.String()
			if !strings.Contains(text, "PRODUCTION READINESS CHECKLIST") || !strings.Contains(text, "✓") ||
				!strings.Contains(text, "Know what's ready and what still needs work.") {
				t.Fatalf("missing product banner content: %q", text)
			}
			if style.color != strings.Contains(text, "\x1b[") {
				t.Fatalf("unexpected color state: %q", text)
			}
		})
	}
}

func TestAuthenticationPanelsExplainTheBrowserFlowAndNextCommands(t *testing.T) {
	var output bytes.Buffer
	printAuthenticationStart(&output, terminalStyle{}, "codex")
	printAuthenticationSuccess(&output, terminalStyle{}, "codex")
	text := output.String()
	for _, expected := range []string{
		"CONNECT CODEX", "official browser sign-in", "Finish in your browser", "CONNECTED",
		"Codex is ready", "prc full codex --plan", "prc full codex",
	} {
		if !strings.Contains(text, expected) {
			t.Fatalf("authentication panels omitted %q: %s", expected, text)
		}
	}
}

func TestAIReviewResultsUseClearAdvisoryLabelsAndCanBeLimited(t *testing.T) {
	results := make([]model.ControlResult, 0, 18)
	for index := 0; index < 18; index++ {
		candidate := "advisory_pass_candidate"
		priority := "low"
		if index == 0 {
			candidate, priority = "advisory_fail_candidate", "critical"
		} else if index == 1 {
			candidate, priority = "needs_evidence", "high"
		}
		results = append(results, model.ControlResult{
			ControlID: fmt.Sprintf("PRC-TEST-%03d", index+1),
			AIReview: &model.AIControlReview{
				AssessmentCandidate: candidate, Priority: priority,
				Advice: "Review the project-specific evidence and correct the gap.",
			},
		})
	}
	var output bytes.Buffer
	summary := controlreview.Summary{
		Provider: "codex", ReviewDepth: "deep", ReviewedControls: 18,
		CompletedBatches: 3, ReusedBatches: 1, Duration: time.Minute,
		TokenUsageBatches: 2,
		TokenUsage: controlreview.TokenUsage{
			InputTokens: 12_000, CachedInputTokens: 8_000, OutputTokens: 600, ReasoningOutputTokens: 300,
		},
	}
	printAIReviewResults(&output, model.RunResult{ControlResults: results}, summary, terminalStyle{}, false)
	text := output.String()
	for _, expected := range []string{
		"AI REVIEW", "Codex · deep review", "1 work", "1 proof", "16 okay", "AI advice only", "Tokens",
		"✗ WORK", "? PROOF", "PRC-TEST-001", "+ 10 more", "--details",
	} {
		if !strings.Contains(text, expected) {
			t.Fatalf("AI rule review omitted %q: %s", expected, text)
		}
	}
	if strings.Contains(text, "PRC-TEST-018") {
		t.Fatalf("default terminal output was not limited: %s", text)
	}

	output.Reset()
	printAIReviewResults(&output, model.RunResult{ControlResults: results}, summary, terminalStyle{}, true)
	if !strings.Contains(output.String(), "PRC-TEST-018") {
		t.Fatalf("detailed terminal output omitted the final review: %s", output.String())
	}
}

func TestLocalCheckBarIsBoundedAndReadable(t *testing.T) {
	for _, item := range []struct {
		percentage int
		want       string
	}{
		{-10, "░░░░░░░░░░"},
		{50, "█████░░░░░"},
		{88, "█████████░"},
		{120, "██████████"},
	} {
		if got := localCheckBar(item.percentage, 10); got != item.want {
			t.Fatalf("localCheckBar(%d, 10) = %q, want %q", item.percentage, got, item.want)
		}
	}
	if got := localCheckBar(50, 0); got != "" {
		t.Fatalf("zero-width bar = %q", got)
	}
}

func TestFinalScanSummaryIsActionableAndOrdered(t *testing.T) {
	run := model.RunResult{
		TerminalState: "no_go",
		Plan:          model.Plan{ProfileID: "prc/core-repository", ProfileVersion: "1.0"},
		Inventory:     model.Inventory{TargetName: "sample-app", Digest: "digest"},
		Results: []model.AssertionResult{
			{AssertionID: "PRC-LOW", Assessment: "fail", Execution: "completed", Severity: "low", Summary: "A lower-priority failure."},
			{AssertionID: "PRC-CRITICAL", Assessment: "fail", Execution: "completed", Severity: "critical", Summary: "A release-blocking failure."},
			{AssertionID: "PRC-MANUAL", Assessment: "manual_review", Execution: "completed", Severity: "medium", Summary: "A reviewer must decide."},
			{AssertionID: "PRC-PASS", Assessment: "pass", Execution: "completed", Severity: "high", Summary: "Passed."},
		},
		ControlCatalog: &model.ControlCatalogSummary{ControlCount: 2, ProfileTerminalState: "no_go"},
		ControlResults: []model.ControlResult{{Disposition: "needs_review"}, {Disposition: "confirmed_failure"}},
	}
	var output bytes.Buffer
	printScanSummary(&output, run, controlreview.Summary{}, terminalStyle{}, "/reports/sample.html", false)
	text := output.String()
	for _, expected := range []string{"SCAN RESULT", "LOCAL CHECKS", "COVERAGE", "REPORT", "/reports/sample.html", "Open it for every check", "no project scripts run"} {
		if !strings.Contains(text, expected) {
			t.Errorf("final summary missing %q: %s", expected, text)
		}
	}
	attentionText := text[strings.Index(text, "LOCAL CHECKS"):]
	criticalAt, manualAt, lowAt := strings.Index(attentionText, "PRC-CRITICAL"), strings.Index(attentionText, "PRC-MANUAL"), strings.Index(attentionText, "PRC-LOW")
	if criticalAt < 0 || manualAt <= criticalAt || lowAt <= manualAt {
		t.Fatalf("attention items are not ordered by severity: %d %d %d\n%s", criticalAt, manualAt, lowAt, text)
	}
	if strings.Contains(text, "Result details") || strings.Contains(text, "Control dispositions") {
		t.Fatalf("dense implementation details leaked into the final handoff: %s", text)
	}
	if strings.Contains(text, "PRC-PASS") {
		t.Fatalf("clean default output showed passing detail: %s", text)
	}

	output.Reset()
	printScanSummary(&output, run, controlreview.Summary{}, terminalStyle{}, "/reports/sample.html", true)
	if !strings.Contains(output.String(), "Every local check") || !strings.Contains(output.String(), "PRC-PASS") {
		t.Fatalf("detailed output omitted a passing local check: %s", output.String())
	}
}

func TestFinalScanSummaryMakesInteractiveReportPathClickable(t *testing.T) {
	run := model.RunResult{
		Plan:      model.Plan{ProfileID: "prc/quick", ProfileVersion: "1.0"},
		Inventory: model.Inventory{TargetName: "sample-app", Digest: "digest"},
	}
	reportPath := filepath.Join(t.TempDir(), "readiness report.html")
	var output bytes.Buffer
	printScanSummary(&output, run, controlreview.Summary{}, terminalStyle{hyperlink: true}, reportPath, false)
	text := output.String()
	if !strings.Contains(text, "\x1b]8;;file://") || !strings.Contains(text, "%20") ||
		!strings.Contains(text, reportPath) || !strings.Contains(text, "Click to open") {
		t.Fatalf("interactive summary did not include a clickable report path: %q", text)
	}
}

func TestFinalScanSummaryPointsBlockedCoreSecretCheckToSafeShortcut(t *testing.T) {
	run := model.RunResult{
		Plan:      model.Plan{ProfileID: "prc/core-repository", ProfileVersion: "1.0"},
		Inventory: model.Inventory{TargetName: "sample-app", Digest: "digest"},
		Results: []model.AssertionResult{{
			AssertionID: "PRC-A-CORE-013", Assessment: "unknown", Execution: "blocked",
			Severity: "high", Summary: "The local adapter was not authorized.",
		}},
	}
	var output bytes.Buffer
	printScanSummary(&output, run, controlreview.Summary{}, terminalStyle{}, "", false)
	if !strings.Contains(output.String(), "Run `prc verify`") || !strings.Contains(output.String(), "must already be local") {
		t.Fatalf("blocked secret check omitted the safe next step: %s", output.String())
	}
}
