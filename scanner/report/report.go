package report

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html/template"
	"io"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/MarinJursic/production-readiness-checklist/scanner/controlprogram"
	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
)

func Write(format string, output io.Writer, run model.RunResult) error {
	switch format {
	case "markdown":
		return writeMarkdown(output, run)
	case "sarif":
		return writeSARIF(output, run)
	case "junit":
		return writeJUnit(output, run)
	case "html":
		return writeHTML(output, run)
	default:
		return fmt.Errorf("unsupported report format %q", format)
	}
}

func assessmentCounts(results []model.AssertionResult) [][2]string {
	counts := map[string]int{}
	for _, result := range results {
		counts[result.Assessment]++
	}
	keys := make([]string, 0, len(counts))
	for key := range counts {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	values := make([][2]string, 0, len(keys))
	for _, key := range keys {
		values = append(values, [2]string{key, strconv.Itoa(counts[key])})
	}
	return values
}

func controlDispositionCounts(results []model.ControlResult) [][2]string {
	counts := map[string]int{}
	for _, result := range results {
		counts[result.Disposition]++
	}
	keys := make([]string, 0, len(counts))
	for key := range counts {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	values := make([][2]string, 0, len(keys))
	for _, key := range keys {
		values = append(values, [2]string{key, strconv.Itoa(counts[key])})
	}
	return values
}

// resultMeaning keeps unlike kinds of evidence separate. A deterministic
// assertion pass is narrower than a complete-control pass, while an AI review
// remains advice even when its citation points to a real snapshot line.
type resultMeaning struct {
	VerifiedProblems      int
	NarrowChecksPassed    int
	LocalChecksUnresolved int
	ManualDecisions       int
	ControlsNeedingReview int
	AIAdvice              int
}

// LocalCheckSummary is deliberately narrower than a production-readiness
// score. It reports the pass rate for applicable assertions in the selected
// local profile, while Label and Explanation continue to honor gate policy.
// A single required failure can therefore produce an 88% pass rate and a
// truthful NOT READY result at the same time.
type LocalCheckSummary struct {
	Passed        int
	Applicable    int
	NotApplicable int
	Failed        int
	Blocking      int
	Unresolved    int
	Manual        int
	ToReview      int
	Percentage    int
	Label         string
	Tone          string
	Explanation   string
	ProfileState  string
}

// SummarizeLocalChecks returns a transparent, unweighted pass rate for the
// selected deterministic profile. It must never be presented as proof that
// the complete control catalog or the project as a whole is ready.
func SummarizeLocalChecks(run model.RunResult) LocalCheckSummary {
	summary := LocalCheckSummary{ProfileState: profileTerminalState(run)}
	for _, result := range run.Results {
		if result.Assessment == "not_applicable" || result.Applicability == "not_applicable" {
			summary.NotApplicable++
			continue
		}
		summary.Applicable++
		switch result.Assessment {
		case "pass":
			summary.Passed++
		case "fail":
			summary.Failed++
			if strings.EqualFold(result.Gate, "required") {
				summary.Blocking++
			}
		case "manual_review":
			summary.Manual++
		case "unknown", "stale", "conflicting":
			summary.Unresolved++
		}
	}
	summary.ToReview = summary.Unresolved + summary.Manual
	if summary.Applicable > 0 {
		summary.Percentage = (summary.Passed*100 + summary.Applicable/2) / summary.Applicable
	}
	switch summary.ProfileState {
	case "profile_satisfied":
		if summary.Applicable > 0 && summary.Passed == summary.Applicable {
			summary.Label, summary.Tone = "GREAT", "great"
			summary.Explanation = "All applicable checks in the selected local profile passed."
		} else {
			summary.Label, summary.Tone = "GOOD", "good"
			summary.Explanation = "The selected local profile has no blocking result."
		}
	case "machine_work_complete_manual_evidence_remaining":
		summary.Label, summary.Tone = "REVIEW NEEDED", "review"
		summary.Explanation = "Automated checks are done, but a person still needs to review evidence."
	case "assessment_incomplete":
		summary.Label, summary.Tone = "NEEDS WORK", "review"
		summary.Explanation = "Some required checks failed or still need trustworthy evidence."
	case "environment_blocked":
		summary.Label, summary.Tone = "SCAN BLOCKED", "bad"
		summary.Explanation = "The scanner could not finish one or more required checks."
	case "no_go":
		summary.Label, summary.Tone = "NOT READY", "bad"
		summary.Explanation = "One or more release-blocking checks failed."
	default:
		summary.Label, summary.Tone = "REVIEW NEEDED", "review"
		summary.Explanation = "Review the detailed results before making a release decision."
	}
	return summary
}

type controlCategorySummary struct {
	Key                string
	Name               string
	Total              int
	ConfirmedFailures  int
	VerifiedPasses     int
	NotApplicable      int
	Blocked            int
	NeedsReview        int
	PartiallyVerified  int
	Retired            int
	LocalPassed        int
	LocalApplicable    int
	LocalNotApplicable int
	LocalFailed        int
	LocalUnresolved    int
	LocalManual        int
	LocalPercentage    int
	HasLocalScore      bool
	LocalLabel         string
	LocalTone          string
}

func summarizeControlCategories(run model.RunResult) []controlCategorySummary {
	byKey := map[string]*controlCategorySummary{}
	categoryByControlID := map[string]string{}
	for _, result := range run.ControlResults {
		key, name := controlCategory(result.Source.Path)
		category := byKey[key]
		if category == nil {
			category = &controlCategorySummary{Key: key, Name: name}
			byKey[key] = category
		}
		categoryByControlID[result.ControlID] = key
		category.Total++
		switch result.Disposition {
		case "confirmed_failure":
			category.ConfirmedFailures++
		case "verified_pass":
			category.VerifiedPasses++
		case "not_applicable":
			category.NotApplicable++
		case "blocked":
			category.Blocked++
		case "needs_review":
			category.NeedsReview++
		case "partially_verified":
			category.PartiallyVerified++
		case "retired":
			category.Retired++
		}
	}
	for _, result := range run.Results {
		linkedCategories := map[string]bool{}
		for _, controlID := range result.ControlIDs {
			if key := categoryByControlID[controlID]; key != "" {
				linkedCategories[key] = true
			}
		}
		for key := range linkedCategories {
			category := byKey[key]
			if result.Assessment == "not_applicable" || result.Applicability == "not_applicable" {
				category.LocalNotApplicable++
				continue
			}
			category.LocalApplicable++
			switch result.Assessment {
			case "pass":
				category.LocalPassed++
			case "fail":
				category.LocalFailed++
			case "manual_review":
				category.LocalManual++
			case "unknown", "stale", "conflicting":
				category.LocalUnresolved++
			}
		}
	}
	categories := make([]controlCategorySummary, 0, len(byKey))
	for _, category := range byKey {
		if category.LocalApplicable > 0 {
			category.HasLocalScore = true
			category.LocalPercentage = (category.LocalPassed*100 + category.LocalApplicable/2) / category.LocalApplicable
			switch {
			case category.LocalFailed > 0:
				category.LocalLabel, category.LocalTone = "NEEDS WORK", "bad"
			case category.LocalUnresolved+category.LocalManual > 0:
				category.LocalLabel, category.LocalTone = "REVIEW", "review"
			case category.LocalPassed == category.LocalApplicable && category.LocalApplicable >= 3:
				category.LocalLabel, category.LocalTone = "GREAT", "great"
			case category.LocalPassed == category.LocalApplicable:
				category.LocalLabel, category.LocalTone = "GOOD", "good"
			default:
				category.LocalLabel, category.LocalTone = "REVIEW", "review"
			}
		} else {
			category.LocalLabel, category.LocalTone = "NOT SCORED", "unscored"
		}
		categories = append(categories, *category)
	}
	sort.Slice(categories, func(left, right int) bool {
		leftRank := categoryScoreRank(categories[left])
		rightRank := categoryScoreRank(categories[right])
		if leftRank != rightRank {
			return leftRank < rightRank
		}
		return categories[left].Name < categories[right].Name
	})
	return categories
}

func categoryScoreRank(category controlCategorySummary) int {
	switch category.LocalTone {
	case "bad":
		return 0
	case "review":
		return 1
	case "good":
		return 2
	case "great":
		return 3
	default:
		return 4
	}
}

func controlCategory(sourcePath string) (string, string) {
	normalized := strings.ReplaceAll(sourcePath, "\\", "/")
	base := strings.TrimSuffix(path.Base(normalized), path.Ext(normalized))
	if separator := strings.IndexByte(base, '-'); separator > 0 {
		numeric := true
		for _, character := range base[:separator] {
			if character < '0' || character > '9' {
				numeric = false
				break
			}
		}
		if numeric {
			base = base[separator+1:]
		}
	}
	if base == "" || base == "." {
		base = "uncategorized"
	}
	scope := "catalog"
	if strings.Contains(normalized, "/engineering/") {
		scope = "engineering"
	} else if strings.Contains(normalized, "/checklists/") {
		scope = "checklist"
	}
	key := scope + "-" + safeSlug(base)
	return key, categoryDisplayName(base)
}

func categoryDisplayName(base string) string {
	names := map[string]string{
		"ai-ml-and-ai-assisted-development":          "AI, ML & agent-assisted development",
		"application-services-and-apis":              "Application services & APIs",
		"data-and-information-lifecycle":             "Data & information lifecycle",
		"data-privacy-performance":                   "Data, privacy & performance",
		"developer-experience-platform-and-delivery": "Developer experience, platform & delivery",
		"environments-quality-experience":            "Environments, quality & experience",
		"evidence-and-decision":                      "Evidence & decision",
		"governance-and-foundations":                 "Governance & foundations",
		"maintenance-vendors-compliance":             "Maintenance, vendors & compliance",
		"operations-sre-and-support":                 "Operations, SRE & support",
		"product-and-requirements":                   "Product & requirements",
		"product-risk-architecture":                  "Product risk & architecture",
		"security-and-cryptography":                  "Security & cryptography",
		"source-build-supply-chain":                  "Source, build & supply chain",
		"specialized-domains-and-release-assurance":  "Specialized domains & release assurance",
		"trust-safety-and-ecosystems":                "Trust, safety & ecosystems",
		"user-experience-web-and-content":            "User experience, web & content",
	}
	if name, ok := names[base]; ok {
		return name
	}
	return humanizeWords(base)
}

func safeSlug(value string) string {
	var result strings.Builder
	previousDash := false
	for _, character := range strings.ToLower(value) {
		if character >= 'a' && character <= 'z' || character >= '0' && character <= '9' {
			result.WriteRune(character)
			previousDash = false
		} else if !previousDash && result.Len() > 0 {
			result.WriteByte('-')
			previousDash = true
		}
	}
	return strings.Trim(result.String(), "-")
}

func humanizeWords(value string) string {
	words := strings.Fields(strings.NewReplacer("-", " ", "_", " ").Replace(value))
	acronyms := map[string]string{"ai": "AI", "api": "API", "apis": "APIs", "ml": "ML", "sre": "SRE", "ux": "UX"}
	for index, word := range words {
		lower := strings.ToLower(word)
		if acronym, ok := acronyms[lower]; ok {
			words[index] = acronym
			continue
		}
		if index == 0 && lower != "and" {
			words[index] = strings.ToUpper(lower[:1]) + lower[1:]
		} else {
			words[index] = lower
		}
	}
	if len(words) == 0 {
		return "Uncategorized"
	}
	return strings.Join(words, " ")
}

func statusClass(value string) string {
	switch value {
	case "pass", "verified_pass":
		return "pass"
	case "fail", "confirmed_failure":
		return "fail"
	case "unknown", "stale", "conflicting", "needs_review":
		return "review"
	case "manual_review":
		return "manual"
	case "blocked":
		return "blocked"
	case "not_applicable":
		return "na"
	case "partially_verified":
		return "partial"
	case "retired":
		return "retired"
	default:
		return "good"
	}
}

func statusSymbol(value string) string {
	switch value {
	case "great", "good", "pass", "verified_pass":
		return "✓"
	case "bad", "fail", "confirmed_failure":
		return "×"
	case "not_applicable", "retired":
		return "–"
	case "partially_verified":
		return "◐"
	default:
		return "!"
	}
}

func prettyStatus(value string) string {
	switch value {
	case "confirmed_failure":
		return "Confirmed failure"
	case "verified_pass":
		return "Verified pass"
	case "needs_review":
		return "Needs evidence"
	case "partially_verified":
		return "Partial evidence"
	case "manual_review":
		return "Manual review"
	case "not_applicable":
		return "Not applicable"
	default:
		return humanizeWords(value)
	}
}

func severityRank(value string) int {
	switch strings.ToLower(value) {
	case "critical":
		return 0
	case "high":
		return 1
	case "medium":
		return 2
	case "low":
		return 3
	case "info":
		return 4
	default:
		return 5
	}
}

func gateRank(value string) int {
	if strings.EqualFold(value, "required") {
		return 0
	}
	return 1
}

func sortedFindings(findings []model.Finding) []model.Finding {
	ordered := append([]model.Finding(nil), findings...)
	sort.SliceStable(ordered, func(left, right int) bool {
		if leftRank, rightRank := severityRank(ordered[left].Severity), severityRank(ordered[right].Severity); leftRank != rightRank {
			return leftRank < rightRank
		}
		if leftRank, rightRank := gateRank(ordered[left].Gate), gateRank(ordered[right].Gate); leftRank != rightRank {
			return leftRank < rightRank
		}
		if ordered[left].Title != ordered[right].Title {
			return ordered[left].Title < ordered[right].Title
		}
		return ordered[left].AssertionID < ordered[right].AssertionID
	})
	return ordered
}

func sortedAttentionResults(results []model.AssertionResult) []model.AssertionResult {
	attention := make([]model.AssertionResult, 0, len(results))
	for _, result := range results {
		if result.Assessment != "pass" && result.Assessment != "not_applicable" {
			attention = append(attention, result)
		}
	}
	sort.SliceStable(attention, func(left, right int) bool {
		if leftRank, rightRank := severityRank(attention[left].Severity), severityRank(attention[right].Severity); leftRank != rightRank {
			return leftRank < rightRank
		}
		if leftRank, rightRank := assessmentAttentionRank(attention[left]), assessmentAttentionRank(attention[right]); leftRank != rightRank {
			return leftRank < rightRank
		}
		return attention[left].AssertionID < attention[right].AssertionID
	})
	return attention
}

func sortedAIRecommendations(results []model.ControlResult) ([]model.ControlResult, int) {
	priorities := make([]model.ControlResult, 0)
	for _, result := range results {
		if result.AIReview != nil && result.AIReview.Priority != "none" {
			priorities = append(priorities, result)
		}
	}
	sort.SliceStable(priorities, func(left, right int) bool {
		leftReview, rightReview := priorities[left].AIReview, priorities[right].AIReview
		if leftRank, rightRank := severityRank(leftReview.Priority), severityRank(rightReview.Priority); leftRank != rightRank {
			return leftRank < rightRank
		}
		if leftReview.AssessmentCandidate != rightReview.AssessmentCandidate {
			return leftReview.AssessmentCandidate == "advisory_fail_candidate"
		}
		return priorities[left].ControlID < priorities[right].ControlID
	})
	total := len(priorities)
	if len(priorities) > 50 {
		priorities = priorities[:50]
	}
	return priorities, total
}

type aiPlanView struct {
	Item model.AIImprovementPlanItem
	Lead model.ControlResult
}

func visibleAIPlan(run model.RunResult) ([]aiPlanView, int) {
	if run.AIImprovementPlan == nil {
		return nil, 0
	}
	byControl := make(map[string]model.ControlResult, len(run.ControlResults))
	for _, control := range run.ControlResults {
		byControl[control.ControlID] = control
	}
	visible := min(len(run.AIImprovementPlan.Items), 24)
	items := make([]aiPlanView, 0, visible)
	for _, item := range run.AIImprovementPlan.Items[:visible] {
		lead := model.ControlResult{}
		if len(item.ControlIDs) > 0 {
			lead = byControl[item.ControlIDs[0]]
		}
		items = append(items, aiPlanView{Item: item, Lead: lead})
	}
	return items, len(run.AIImprovementPlan.Items)
}

func assessmentAttentionRank(result model.AssertionResult) int {
	if result.Assessment == "fail" {
		return 0
	}
	if result.Execution == "blocked" || result.Assessment == "unknown" {
		return 1
	}
	if result.Assessment == "manual_review" {
		return 2
	}
	return 3
}

func summarizeMeaning(run model.RunResult) resultMeaning {
	summary := resultMeaning{VerifiedProblems: len(run.Findings)}
	for _, result := range run.Results {
		switch result.Assessment {
		case "pass":
			summary.NarrowChecksPassed++
		case "manual_review":
			summary.ManualDecisions++
		case "unknown", "stale", "conflicting":
			summary.LocalChecksUnresolved++
		}
	}
	for _, result := range run.ControlResults {
		if result.Disposition == "needs_review" || result.Disposition == "blocked" {
			summary.ControlsNeedingReview++
		}
		if result.AIReview != nil {
			summary.AIAdvice++
		}
	}
	return summary
}

func advisoryVerificationText(review *model.AIControlReview) string {
	if review == nil {
		return "—"
	}
	citation := review.CitationVerification
	if citation == "" { // A readable label for archived v0.12 run records.
		citation = "legacy_unrecorded"
	}
	claim := review.ClaimVerification
	if claim == "" {
		claim = "legacy_unrecorded"
	}
	return "citation=" + citation + "; claim=" + claim
}

func markdownCell(value string) string {
	value = strings.ReplaceAll(value, "\\", "\\\\")
	value = strings.ReplaceAll(value, "|", "\\|")
	return strings.Join(strings.Fields(value), " ")
}

func briefText(value string, limit int) string {
	value = strings.Join(strings.Fields(value), " ")
	characters := []rune(value)
	if limit <= 0 || len(characters) <= limit {
		return value
	}
	return strings.TrimSpace(string(characters[:limit])) + "…"
}

func advisoryLocationText(locations []model.FindingLocation) string {
	values := make([]string, 0, len(locations))
	for _, location := range locations {
		value := fmt.Sprintf("%s:%d", location.Path, location.Line)
		if location.Column > 0 {
			value += fmt.Sprintf(":%d", location.Column)
		}
		values = append(values, value)
	}
	return strings.Join(values, ", ")
}

func writeMarkdown(output io.Writer, run model.RunResult) error {
	meaning := summarizeMeaning(run)
	profileState := run.TerminalState
	if run.ControlCatalog != nil && run.ControlCatalog.ProfileTerminalState != "" {
		profileState = run.ControlCatalog.ProfileTerminalState
	}
	if _, err := fmt.Fprintf(output, "# Production readiness assessment\n\n"+
		"- Run: `%s`\n- Profile: `%s@%s`\n- Target: `%s`\n- Inventory: `%s`\n- Local profile result: **%s**\n- Full catalog coverage result: **%s**\n\n",
		run.RunID, run.Plan.ProfileID, run.Plan.ProfileVersion, markdownCell(run.Inventory.TargetName),
		run.Inventory.Digest, profileState, run.TerminalState); err != nil {
		return err
	}
	if run.Plan.ConfigurationDigest != "" {
		if _, err := fmt.Fprintf(output,
			"- Configuration: `%s`\n- Project: `%s`\n- Environments: %s\n- Artifacts: %s\n\n",
			run.Plan.ConfigurationDigest, run.Plan.ProjectID,
			markdownCell(strings.Join(run.Plan.TargetEnvironments, ", ")),
			markdownCell(strings.Join(run.Plan.ArtifactDigests, ", "))); err != nil {
			return err
		}
	}
	if _, err := fmt.Fprintf(output,
		"## What the result means\n\n| Kind | Count | Meaning |\n| --- | ---: | --- |\n"+
			"| Verified problems | %d | Narrow failures backed by scanner evidence. |\n"+
			"| Narrow checks passed | %d | Exact local assertions that passed; not whole-control certification. |\n"+
			"| Local checks unresolved | %d | Local assertions that still lack usable evidence or conflict. |\n"+
			"| Manual decisions | %d | Assertions that require a person or higher-authority evidence. |\n"+
			"| Controls still needing review or evidence | %d | Broad controls not proved by the selected local profile. |\n"+
			"| Advisory AI reviews | %d | Suggestions only; never verified Pass or final Not Applicable. |\n\n",
		meaning.VerifiedProblems, meaning.NarrowChecksPassed, meaning.LocalChecksUnresolved,
		meaning.ManualDecisions, meaning.ControlsNeedingReview, meaning.AIAdvice); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(output, "## Result counts\n\n| Assessment | Count |\n| --- | ---: |"); err != nil {
		return err
	}
	for _, item := range assessmentCounts(run.Results) {
		if _, err := fmt.Fprintf(output, "| %s | %s |\n", item[0], item[1]); err != nil {
			return err
		}
	}
	if run.ControlCatalog != nil {
		if _, err := fmt.Fprintf(output,
			"\n## Complete control catalog\n\n- Registry: `%s`\n- Registry digest: `%s`\n- Source digest: `%s`\n- Registered controls: **%d** (%d active)\n- Narrow profile terminal state before complete-catalog expansion: **%s**\n\n",
			run.ControlCatalog.RegistryVersion, run.ControlCatalog.RegistrySHA256,
			run.ControlCatalog.SourceSHA256, run.ControlCatalog.ControlCount,
			run.ControlCatalog.ActiveControlCount, run.ControlCatalog.ProfileTerminalState); err != nil {
			return err
		}
		if run.ControlCatalog.ReviewedDeterministicCount+run.ControlCatalog.ReviewedNondeterministicCount > 0 {
			if _, err := fmt.Fprintf(output,
				"- Reviewed deterministic controls: **%d**\n- Reviewed nondeterministic controls: **%d**\n- Deterministic bindings: **%d**\n- Exact program templates: **%d**\n- Exact programs attempted: **%d** (**%d** passed, **%d** failed, **%d** proven Not Applicable)\n- Retained replayable exact evidence documents: **%d**\n- Deterministic controls still blocked: **%d**\n- Classification corpus digest: `%s`\n- Deterministic binding artifact digest: `%s`\n\n",
				run.ControlCatalog.ReviewedDeterministicCount,
				run.ControlCatalog.ReviewedNondeterministicCount,
				run.ControlCatalog.DeterministicBindingCount,
				run.ControlCatalog.DeterministicProgramTemplateCount,
				run.ControlCatalog.DeterministicProgramExecutedCount,
				run.ControlCatalog.DeterministicProgramPassCount,
				run.ControlCatalog.DeterministicProgramFailCount,
				run.ControlCatalog.DeterministicProgramNACount,
				len(run.DeterministicEvidence),
				run.ControlCatalog.DeterministicProgramBlockedCount,
				run.ControlCatalog.ClassificationCorpusSHA256,
				run.ControlCatalog.ControlCheckBindingsSHA256); err != nil {
				return err
			}
		}
		if run.ControlCatalog.AIReviewProvider != "" {
			if _, err := fmt.Fprintf(output,
				"- Advisory AI review: **%s** (model `%s`)\n- AI review state: **%s**\n- AI-reviewed controls: **%d**\n- Advisory failure candidates: **%d**\n- AI advice cannot create a verified Pass or final Not Applicable result.\n\n",
				run.ControlCatalog.AIReviewProvider, run.ControlCatalog.AIReviewModel,
				run.ControlCatalog.AIReviewState, run.ControlCatalog.AIReviewedCount,
				run.ControlCatalog.AIAdvisoryFailCount); err != nil {
				return err
			}
		}
		if run.AIImprovementPlan != nil {
			if _, err := fmt.Fprintf(output,
				"### AI review improvement plan\n\nThe scanner grouped %d reviewed controls into %d exact cause groups. These remain advisory.\n\n| Priority | Cause | Domain | Effort | Reach | Controls |\n| --- | --- | --- | --- | --- | --- |\n",
				run.AIImprovementPlan.ReviewedControlCount, run.AIImprovementPlan.ItemCount); err != nil {
				return err
			}
			for _, item := range run.AIImprovementPlan.Items[:min(len(run.AIImprovementPlan.Items), 50)] {
				visibleControls := item.ControlIDs[:min(len(item.ControlIDs), 8)]
				controls := strings.Join(visibleControls, ", ")
				if len(item.ControlIDs) > len(visibleControls) {
					controls += fmt.Sprintf(" (+%d more)", len(item.ControlIDs)-len(visibleControls))
				}
				if _, err := fmt.Fprintf(output, "| %s | %s | %s | %s | %s | %s |\n",
					markdownCell(item.Priority), markdownCell(item.RootCause), markdownCell(item.Domain),
					markdownCell(item.Effort), markdownCell(item.BlastRadius), markdownCell(controls)); err != nil {
					return err
				}
			}
			if _, err := fmt.Fprintln(output); err != nil {
				return err
			}
		}
		if len(run.AuthoritativeEvidence) > 0 {
			if _, err := fmt.Fprintln(output, "### Signed authoritative evidence"); err != nil {
				return err
			}
			for _, verification := range run.AuthoritativeEvidence {
				if _, err := fmt.Fprintf(output,
					"\n- Bundle `%s`: **%d** %s entries; policy key `%s`; evidence key `%s`; policy digest `%s`; bundle digest `%s`.\n",
					verification.BundleID, verification.EntryCount, verification.Authority,
					verification.PolicySignature.KeyID, verification.EvidenceSignature.KeyID,
					verification.PolicySHA256, verification.BundleSHA256); err != nil {
					return err
				}
			}
			if _, err := fmt.Fprintln(output); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprintln(output, "| Disposition | Count |\n| --- | ---: |"); err != nil {
			return err
		}
		for _, item := range controlDispositionCounts(run.ControlResults) {
			if _, err := fmt.Fprintf(output, "| %s | %s |\n", item[0], item[1]); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprintln(output, "\n### Every registered control\n\n| Disposition | Control | Classification | Route / decision | Binding | Coverage | Authority | Statement | Source | Assertions | AI candidate | AI reason and advice | AI citation locations | AI verification | AI limitations |\n| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |"); err != nil {
			return err
		}
		for _, result := range run.ControlResults {
			aiCandidate, aiReason, aiEvidence, aiVerification, aiLimitations := "—", "—", "—", "—", "—"
			if result.AIReview != nil {
				aiCandidate = result.AIReview.AssessmentCandidate + " / " + result.AIReview.ApplicabilityCandidate + " / " + result.AIReview.Confidence
				aiReason = result.AIReview.Reason + " Risk if ignored: " + result.AIReview.RiskIfIgnored +
					" Suggested work: " + strings.Join(result.AIReview.RemediationSteps, "; ") +
					" Verify with: " + strings.Join(result.AIReview.VerificationSteps, "; ")
				aiEvidence = advisoryLocationText(result.AIReview.Evidence)
				aiVerification = advisoryVerificationText(result.AIReview)
				aiLimitations = strings.Join(result.AIReview.Limitations, "; ")
			}
			binding := "—"
			if result.DeterministicBindingID != "" {
				binding = result.DeterministicBindingID + " / " + result.DeterministicBindingSHA256
			}
			if _, err := fmt.Fprintf(output, "| %s | `%s` | %s | %s / %s | %s | %s | %s | %s | `%s:%d` | %s | %s | %s | %s | %s | %s |\n",
				result.Disposition, result.ControlID, result.Classification,
				result.ClassificationRoute, result.ClassificationDecisionBasis, markdownCell(binding),
				result.Coverage, result.Authority,
				markdownCell(result.Statement), result.Source.Path, result.Source.Line,
				markdownCell(strings.Join(result.ExecutedAssertionIDs, ", ")), markdownCell(aiCandidate),
				markdownCell(aiReason), markdownCell(aiEvidence), markdownCell(aiVerification), markdownCell(aiLimitations)); err != nil {
				return err
			}
		}
	}
	if len(run.AdapterExecutions) > 0 {
		if _, err := fmt.Fprintln(output, "\n## Adapter executions\n\n| Adapter | Manifest | Authorization | Trust | Registry | Status | Execution |\n| --- | --- | --- | --- | --- | --- | --- |"); err != nil {
			return err
		}
		for _, execution := range run.AdapterExecutions {
			registry := execution.Resolution.RegistryID
			if registry == "" {
				registry = "—"
			}
			if _, err := fmt.Fprintf(output, "| `%s` | `%s` | %s | %s | %s | %s | `%s` |\n",
				execution.AdapterID, execution.ManifestSHA256, execution.Resolution.Source,
				execution.Resolution.Trust, registry, execution.Transcript.Summary.Status, execution.ExecutionID); err != nil {
				return err
			}
		}
	}
	if len(run.Findings) > 0 {
		if _, err := fmt.Fprintln(output, "\n## Findings\n\n| Severity | Finding | Assertion | Gate | Summary | Locations | Evidence |\n| --- | --- | --- | --- | --- | ---: | ---: |"); err != nil {
			return err
		}
		for _, finding := range sortedFindings(run.Findings) {
			if _, err := fmt.Fprintf(output, "| %s | `%s` | `%s` | %s | %s | %d | %d |\n",
				finding.Severity, finding.ID, finding.AssertionID, finding.Gate,
				markdownCell(finding.Summary), len(finding.Locations), len(finding.EvidenceIDs)); err != nil {
				return err
			}
		}
	}
	if _, err := fmt.Fprintln(output, "\n## Assertions\n\n| Assessment | Assertion | Severity | Gate | Summary | Evidence |\n| --- | --- | --- | --- | --- | ---: |"); err != nil {
		return err
	}
	for _, result := range run.Results {
		if _, err := fmt.Fprintf(output, "| %s | `%s` | %s | %s | %s | %d |\n",
			result.Assessment, result.AssertionID, result.Severity, result.Gate,
			markdownCell(result.Summary), len(result.EvidenceObserved)); err != nil {
			return err
		}
	}
	_, err := fmt.Fprintln(output, "\nThis report is scoped to the named profile, target inventory, and evidence set. It is not an unqualified production-readiness or compliance claim.")
	return err
}

type sarifLog struct {
	Version string     `json:"version"`
	Schema  string     `json:"$schema"`
	Runs    []sarifRun `json:"runs"`
}

type sarifRun struct {
	Tool       sarifTool         `json:"tool"`
	Results    []sarifResult     `json:"results"`
	Properties map[string]string `json:"properties"`
}

type sarifTool struct {
	Driver sarifDriver `json:"driver"`
}

type sarifDriver struct {
	Name           string      `json:"name"`
	InformationURI string      `json:"informationUri"`
	Rules          []sarifRule `json:"rules"`
}

type sarifRule struct {
	ID               string         `json:"id"`
	ShortDescription sarifMessage   `json:"shortDescription"`
	Properties       map[string]any `json:"properties"`
}

type sarifResult struct {
	RuleID     string          `json:"ruleId"`
	Level      string          `json:"level"`
	Message    sarifMessage    `json:"message"`
	Locations  []sarifLocation `json:"locations,omitempty"`
	Properties map[string]any  `json:"properties"`
}

type sarifMessage struct {
	Text string `json:"text"`
}

type sarifLocation struct {
	PhysicalLocation sarifPhysicalLocation `json:"physicalLocation"`
}

type sarifPhysicalLocation struct {
	ArtifactLocation sarifArtifactLocation `json:"artifactLocation"`
	Region           *sarifRegion          `json:"region,omitempty"`
}

type sarifArtifactLocation struct {
	URI string `json:"uri"`
}

type sarifRegion struct {
	StartLine   int `json:"startLine,omitempty"`
	StartColumn int `json:"startColumn,omitempty"`
}

func sarifLevel(severity string) string {
	switch severity {
	case "critical", "high":
		return "error"
	case "medium":
		return "warning"
	default:
		return "note"
	}
}

func writeSARIF(output io.Writer, run model.RunResult) error {
	rules := make([]sarifRule, 0)
	results := make([]sarifResult, 0)
	seenRules := map[string]bool{}
	for _, finding := range sortedFindings(run.Findings) {
		if !seenRules[finding.AssertionID] {
			seenRules[finding.AssertionID] = true
			rules = append(rules, sarifRule{
				ID: finding.AssertionID, ShortDescription: sarifMessage{Text: finding.Title},
				Properties: map[string]any{"severity": finding.Severity, "gate": finding.Gate, "control_ids": finding.ControlIDs},
			})
		}
		locations := make([]sarifLocation, 0, len(finding.Locations))
		for _, location := range finding.Locations {
			physical := sarifPhysicalLocation{ArtifactLocation: sarifArtifactLocation{URI: location.Path}}
			if location.Line > 0 {
				physical.Region = &sarifRegion{StartLine: location.Line, StartColumn: location.Column}
			}
			locations = append(locations, sarifLocation{PhysicalLocation: physical})
		}
		results = append(results, sarifResult{
			RuleID: finding.AssertionID, Level: sarifLevel(finding.Severity),
			Message: sarifMessage{Text: finding.Summary}, Locations: locations,
			Properties: map[string]any{
				"finding_id": finding.ID, "fingerprint": finding.Fingerprint,
				"severity": finding.Severity, "gate": finding.Gate,
				"remediation_class": finding.RemediationClass, "evidence_ids": finding.EvidenceIDs,
			},
		})
	}
	log := sarifLog{
		Version: "2.1.0", Schema: "https://json.schemastore.org/sarif-2.1.0.json",
		Runs: []sarifRun{{
			Tool:    sarifTool{Driver: sarifDriver{Name: "Production Readiness Checklist", InformationURI: "https://marinjursic.github.io/production-readiness-checklist/", Rules: rules}},
			Results: results,
			Properties: map[string]string{
				"run_id": run.RunID, "profile": run.Plan.ProfileID + "@" + run.Plan.ProfileVersion,
				"inventory_digest": run.Inventory.Digest, "terminal_state": run.TerminalState,
				"adapter_execution_count": strconv.Itoa(len(run.AdapterExecutions)),
			},
		}},
	}
	if run.Plan.ConfigurationDigest != "" {
		log.Runs[0].Properties["configuration_digest"] = run.Plan.ConfigurationDigest
		log.Runs[0].Properties["project_id"] = run.Plan.ProjectID
		log.Runs[0].Properties["target_environments"] = strings.Join(run.Plan.TargetEnvironments, ",")
		log.Runs[0].Properties["artifact_digests"] = strings.Join(run.Plan.ArtifactDigests, ",")
	}
	encoder := json.NewEncoder(output)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	return encoder.Encode(log)
}

type junitSuite struct {
	XMLName    xml.Name        `xml:"testsuite"`
	Name       string          `xml:"name,attr"`
	Tests      int             `xml:"tests,attr"`
	Failures   int             `xml:"failures,attr"`
	Errors     int             `xml:"errors,attr"`
	Skipped    int             `xml:"skipped,attr"`
	Time       string          `xml:"time,attr"`
	Properties []junitProperty `xml:"properties>property"`
	Cases      []junitCase     `xml:"testcase"`
}

type junitProperty struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

type junitCase struct {
	Name      string        `xml:"name,attr"`
	Classname string        `xml:"classname,attr"`
	Failure   *junitMessage `xml:"failure,omitempty"`
	Error     *junitMessage `xml:"error,omitempty"`
	Skipped   *junitMessage `xml:"skipped,omitempty"`
}

type junitMessage struct {
	Message string `xml:"message,attr"`
	Type    string `xml:"type,attr,omitempty"`
	Text    string `xml:",chardata"`
}

func writeJUnit(output io.Writer, run model.RunResult) error {
	suite := junitSuite{
		Name: run.Plan.ProfileID + "@" + run.Plan.ProfileVersion, Tests: len(run.Results),
		Time: strconv.FormatFloat(run.CompletedAt.Sub(run.StartedAt).Seconds(), 'f', 3, 64),
		Properties: []junitProperty{
			{Name: "run_id", Value: run.RunID}, {Name: "inventory_digest", Value: run.Inventory.Digest},
			{Name: "terminal_state", Value: run.TerminalState},
			{Name: "adapter_execution_count", Value: strconv.Itoa(len(run.AdapterExecutions))},
		},
	}
	if run.Plan.ConfigurationDigest != "" {
		suite.Properties = append(suite.Properties,
			junitProperty{Name: "configuration_digest", Value: run.Plan.ConfigurationDigest},
			junitProperty{Name: "project_id", Value: run.Plan.ProjectID},
			junitProperty{Name: "target_environments", Value: strings.Join(run.Plan.TargetEnvironments, ",")},
			junitProperty{Name: "artifact_digests", Value: strings.Join(run.Plan.ArtifactDigests, ",")},
		)
	}
	for _, result := range run.Results {
		item := junitCase{Name: result.AssertionID, Classname: run.Plan.ProfileID}
		detail := &junitMessage{Message: result.Summary, Type: result.Assessment, Text: result.Summary}
		switch {
		case result.Assessment == "fail":
			item.Failure = detail
			suite.Failures++
		case result.Execution == "error" || result.Execution == "blocked" || result.Assessment == "unknown" || result.Assessment == "stale" || result.Assessment == "conflicting":
			item.Error = detail
			suite.Errors++
		case result.Assessment == "manual_review" || result.Assessment == "not_applicable":
			item.Skipped = detail
			suite.Skipped++
		}
		suite.Cases = append(suite.Cases, item)
	}
	if _, err := io.WriteString(output, xml.Header); err != nil {
		return err
	}
	encoder := xml.NewEncoder(output)
	encoder.Indent("", "  ")
	if err := encoder.Encode(suite); err != nil {
		return err
	}
	_, err := io.WriteString(output, "\n")
	return err
}

const htmlReport = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>{{.Run.Inventory.TargetName}} — production readiness report</title>
  <style>
    :root {
      color-scheme: light;
      --canvas: #fff; --surface: #fff; --surface-soft: #f6f8fa;
      --ink: #1f2328; --muted: #59636e; --line: #d1d9e0;
      --blue: #0969da; --blue-bg: #ddf4ff; --green: #1a7f37; --green-bg: #dafbe1;
      --yellow: #9a6700; --yellow-bg: #fff8c5; --red: #cf222e; --red-bg: #ffebe9;
      --gray: #59636e; --gray-bg: #f6f8fa; --radius: .375rem;
    }
    * { box-sizing: border-box; }
    html { scroll-behavior: auto; }
    body { background: var(--canvas); color: var(--ink); font: 14px/1.55 -apple-system, BlinkMacSystemFont, "Segoe UI", Helvetica, Arial, sans-serif; margin: 0; }
    a { color: var(--blue); } a:focus-visible, button:focus-visible, input:focus-visible, select:focus-visible, summary:focus-visible { outline: 2px solid var(--blue); outline-offset: 2px; }
    code { background: #f5f7fa; border-radius: 2px; color: #0938c2; font: .88em ui-monospace, "Roboto Mono", SFMono-Regular, Menlo, Consolas, monospace; overflow-wrap: anywhere; padding: .08rem .25rem; }
    h1, h2, h3 { line-height: 1.25; } h2 { font-size: clamp(1.3rem, 2vw, 1.6rem); font-weight: 600; letter-spacing: -.015em; margin: 0 0 .35rem; }
    .skip-link { background: white; left: 1rem; padding: .7rem 1rem; position: fixed; top: -5rem; z-index: 20; }
    .skip-link:focus { top: 1rem; }
    .hero { background: #fff; border-top: 3px solid var(--blue); color: var(--ink); padding: 0 1.25rem; }
    .hero-inner, main, .nav-inner { margin: 0 auto; max-width: 68rem; }
    .report-brand { align-items: center; border-bottom: 1px solid var(--line); display: flex; font-size: .78rem; justify-content: space-between; margin-bottom: 1.5rem; padding: .8rem 0; }
    .report-brand strong { font-weight: 650; }
    .report-brand span { color: var(--muted); }
    .hero-heading { margin: 0 auto; max-width: 52rem; text-align: center; }
    .hero h1 { font-size: clamp(1.8rem, 4vw, 2.45rem); font-weight: 600; letter-spacing: -.03em; margin: 0; overflow-wrap: anywhere; }
    .rating { align-items: center; display: inline-flex; font-size: .72rem; font-weight: 700; gap: .35rem; letter-spacing: .06em; padding-top: .45rem; text-transform: uppercase; white-space: nowrap; }
    .rating-great { color: var(--green); } .rating-good { color: #0d47a1; } .rating-review { color: var(--yellow); } .rating-bad { color: var(--red); }
    .hero-grid { margin: 1rem auto 0; max-width: 46rem; text-align: center; }
    .score-panel { align-items: center; display: flex; flex-direction: column; }
    .score-gauge { --gauge-color: #0c6; aspect-ratio: 1; flex: 0 0 auto; position: relative; width: 8rem; }
    .score-gauge svg { display: block; height: 100%; overflow: visible; width: 100%; }
    .score-gauge-track, .score-gauge-value { fill: none; stroke-width: 8; }
    .score-gauge-track { stroke: #f1f3f4; }
    .score-gauge-value { stroke: var(--gauge-color); stroke-linecap: round; transform: rotate(-90deg); transform-origin: 50% 50%; }
    .tone-good .hero-score-gauge { --gauge-color: #0969da; } .tone-review .hero-score-gauge { --gauge-color: #bf8700; } .tone-bad .hero-score-gauge { --gauge-color: #cf222e; }
    .score-gauge-center { align-items: baseline; display: flex; inset: 0; justify-content: center; position: absolute; top: 3.45rem; }
    .score-gauge-number { font-size: 2.35rem; font-weight: 500; letter-spacing: -.05em; line-height: 1; }
    .score-gauge-scale { color: var(--muted); font-size: .65rem; font-weight: 600; margin-left: .15rem; }
    .score-kicker { color: var(--muted); display: block; font-size: .7rem; font-weight: 600; letter-spacing: .08em; text-transform: uppercase; }
    .hero-score-gauge { width: 10rem; }
    .overall-rating { font-size: .8rem; justify-content: center; margin: .55rem 0 0; padding: 0; }
    .score-status { font-size: 1rem; font-weight: 560; margin: .3rem auto 0; max-width: 34rem; }
    .hero-metrics { border-bottom: 1px solid var(--line); border-top: 1px solid var(--line); display: grid; grid-template-columns: repeat(4, 1fr); margin: 1.15rem auto 0; max-width: 40rem; }
    .hero-metric { padding: .75rem .5rem .7rem; }
    .hero-metric + .hero-metric { border-left: 1px solid var(--line); }
    .hero-metric strong { display: block; font-size: 1.2rem; font-weight: 600; line-height: 1.1; }
    .hero-metric span { color: var(--muted); display: block; font-size: .7rem; margin-top: .2rem; text-transform: uppercase; }
    .metric-pass strong { color: var(--green); } .metric-fail strong { color: var(--red); } .metric-review strong { color: var(--yellow); }
    .hero-help { margin: .6rem auto .8rem; max-width: 40rem; text-align: left; }
    .hero-help > summary { text-align: center; }
    .report-nav { background: rgba(255,255,255,.98); border-bottom: 1px solid var(--line); position: sticky; top: 0; z-index: 10; }
    .nav-inner { display: flex; gap: .15rem; overflow-x: auto; padding: 0 1rem; scrollbar-width: none; }
    .nav-inner::-webkit-scrollbar { display: none; }
    .nav-inner a { border-bottom: 2px solid transparent; color: var(--muted); font-size: .78rem; font-weight: 600; padding: .65rem .7rem .55rem; text-decoration: none; white-space: nowrap; }
    .nav-inner a:hover { border-bottom-color: var(--blue); color: var(--blue); }
    main { padding: 2.25rem 1rem 5rem; }
    .section { scroll-margin-top: 4rem; margin: 0 0 3.5rem; }
    .section-intro { color: var(--muted); margin: 0 0 1.25rem; max-width: 54rem; }
    .notice { background: #f8faff; border: 1px solid #d2e3fc; border-left: 3px solid var(--blue); border-radius: 3px; padding: .75rem .9rem; }
    .finding-card { --severity-color: #59636e; --severity-bg: var(--gray-bg); background: var(--surface); border: 0; border-bottom: 1px solid var(--line); border-top: 1px solid var(--line); border-radius: 0; content-visibility: visible; contain-intrinsic-size: none; margin: 0 0 -.0625rem; overflow: hidden; }
    .finding-card[open] { background: #fff; }
    .finding-card > summary { align-items: center; display: grid; gap: .8rem; grid-template-columns: auto minmax(0, 1fr) auto; list-style: none; padding: .9rem 1rem; }
    .finding-card > summary::-webkit-details-marker { display: none; }
    .finding-card > summary::after { color: var(--muted); content: "Open"; font-size: .7rem; font-weight: 600; line-height: 1; }
    .finding-card[open] > summary::after { content: "Close"; }
    .finding-card > .finding-body { border-top: 1px solid var(--line); padding: .9rem 1rem 1rem; }
    .finding-title { display: block; font-size: .95rem; font-weight: 650; line-height: 1.35; }
    .finding-preview { color: var(--muted); display: block; font-size: .78rem; line-height: 1.45; margin-top: .15rem; }
    .severity-pill { background: var(--severity-bg); color: var(--severity-color); }
    .severity-critical { --severity-color: #a40e26; --severity-bg: #ffebe9; }
    .severity-high { --severity-color: #cf222e; --severity-bg: var(--red-bg); }
    .severity-medium { --severity-color: #9a6700; --severity-bg: var(--yellow-bg); }
    .severity-low { --severity-color: #0969da; --severity-bg: var(--blue-bg); }
    .severity-info, .severity-unknown { --severity-color: var(--gray); --severity-bg: var(--gray-bg); }
    .severity-label { color: var(--severity-color); font-size: .66rem; font-weight: 700; letter-spacing: .04em; margin-left: auto; text-transform: uppercase; white-space: nowrap; }
    .severity-key { align-items: center; color: var(--muted); display: flex; flex-wrap: wrap; font-size: .75rem; gap: .45rem .8rem; margin: -.45rem 0 1rem; }
    .severity-key .key-dot { background: var(--severity-color); }
    .pill { border: 1px solid transparent; border-radius: 3px; display: inline-flex; font-size: .66rem; font-weight: 650; gap: .25rem; letter-spacing: .03em; padding: .18rem .35rem; text-transform: uppercase; white-space: nowrap; }
    .pill-pass { background: var(--green-bg); color: var(--green); } .pill-partial, .pill-good { background: var(--blue-bg); color: #1e40af; }
    .pill-fail, .pill-bad { background: var(--red-bg); color: var(--red); } .pill-review, .pill-blocked, .pill-manual { background: var(--yellow-bg); color: var(--yellow); }
    .pill-na, .pill-retired { background: var(--gray-bg); color: var(--gray); }
    .next-action { background: var(--yellow-bg); border: 1px solid #d4a72c66; margin-top: .7rem; padding: .65rem .75rem; }
    .primary-detail { margin: .3rem 0 .65rem; max-width: 58rem; }
    .meta { color: var(--muted); } .compact { margin: .35rem 0; } ul { padding-left: 1.35rem; }
    details { background: var(--surface); border: 0; border-bottom: 1px solid var(--line); content-visibility: auto; contain-intrinsic-size: auto 7rem; margin: 0; }
    details[open] { background: #fff; }
    summary { cursor: pointer; font-weight: 520; padding: .75rem .8rem; }
    details > .detail-body { border-top: 1px solid var(--line); padding: .9rem 1rem 1rem; }
    .summary-help { background: transparent; border: 0; border-top: 1px solid var(--line); margin-top: .8rem; }
    .summary-help > summary { color: var(--blue); font-size: .78rem; padding: .65rem 0; }
    .summary-help .detail-body { font-size: .82rem; }
    .finding-details { background: transparent; border: 0; margin-top: .7rem; }
    .finding-details > summary, .raw-details > summary { color: var(--blue); font-size: .78rem; font-weight: 650; padding: .45rem 0; }
    .finding-details > .detail-body, .raw-details > .detail-body { background: var(--surface-soft); border: 1px solid var(--line); border-radius: var(--radius); padding: .8rem .9rem; }
    .raw-details { background: transparent; border: 0; margin-top: .7rem; }
    .raw-list { background: transparent; border: 0; margin-top: .35rem; }
    .raw-list > summary { color: var(--blue); font-size: .75rem; padding: .35rem 0; }
    .section-disclosure { border: 0; border-bottom: 1px solid var(--line); border-top: 1px solid var(--line); border-radius: 0; margin-bottom: -.0625rem; scroll-margin-top: 4rem; }
    .section-disclosure[open] { margin-bottom: 2.5rem; }
    .section-disclosure > summary { align-items: center; display: flex; gap: 1rem; justify-content: space-between; list-style: none; padding: 1rem 1.1rem; }
    .section-disclosure > summary::-webkit-details-marker { display: none; }
    .section-disclosure > summary::after { color: var(--blue); content: "Open"; flex: 0 0 auto; font-size: .72rem; font-weight: 650; }
    .section-disclosure[open] > summary::after { content: "Close"; }
    .disclosure-title { display: block; font-size: 1rem; font-weight: 650; }
    .disclosure-note { color: var(--muted); display: block; font-size: .76rem; font-weight: 400; margin-top: .15rem; }
    .section-disclosure > .section-disclosure-body { border-top: 1px solid var(--line); padding: 1.15rem; }
    .section-disclosure-body h2 { font-size: 1.15rem; margin-top: .2rem; }
    .result-group { border: 1px solid var(--line); border-radius: var(--radius); margin-top: .8rem; }
    .result-group > summary { background: var(--surface-soft); display: flex; justify-content: space-between; }
    .result-count { color: var(--muted); font-size: .75rem; font-weight: 500; }
    .assertion-row > summary, .control-row > summary { align-items: flex-start; display: flex; gap: .55rem; }
    .row-summary { color: var(--ink); min-width: 0; overflow-wrap: anywhere; }
    .row-summary code { white-space: nowrap; }
    .advice { background: var(--blue-bg); border-left: 3px solid var(--blue); margin: .8rem 0; padding: .7rem .8rem; }
    .more-controls { display: block; margin: 1rem auto 0; min-width: 12rem; }
    .result-pass { border-left: 3px solid #0c6; } .result-fail { border-left: 3px solid #f33; }
    .result-review, .result-blocked, .result-manual { border-left: 3px solid #fa3; } .result-na { border-left: 3px solid #9e9e9e; } .result-partial { border-left: 3px solid #06f; }
    dl { display: grid; gap: .25rem 1rem; grid-template-columns: minmax(10rem, .3fr) minmax(0, 1fr); margin: .5rem 0; }
    dt { font-weight: 600; } dd { margin: 0 0 .55rem; min-width: 0; }
    .category-key { align-items: center; color: var(--muted); display: flex; flex-wrap: wrap; font-size: .74rem; gap: .45rem .9rem; margin: 0; padding: .35rem 0 .8rem; }
    .category-key strong { color: var(--ink); }
    .key-item { align-items: center; display: inline-flex; gap: .35rem; }
    .key-dot { background: var(--gray); border-radius: 999px; height: .5rem; width: .5rem; } .key-dot.great { background: #1a7f37; } .key-dot.good { background: #0969da; } .key-dot.review { background: #bf8700; } .key-dot.bad { background: #cf222e; }
    .category-grid { border-bottom: 1px solid var(--line); display: grid; grid-template-columns: repeat(4, minmax(0, 1fr)); }
    .category-card { --category-color: #59636e; background: #fff; border: 0; border-top: 1px solid var(--line); color: inherit; display: flex; min-width: 0; padding: 1.25rem .75rem 1.1rem; text-align: center; text-decoration: none; transition: background-color .12s; }
    .category-card:hover { background: var(--surface-soft); }
    .category-great { --category-color: #1a7f37; } .category-good { --category-color: #0969da; } .category-review { --category-color: #bf8700; } .category-bad { --category-color: #cf222e; }
    .category-overview { align-items: center; display: flex; flex: 1; flex-direction: column; min-width: 0; }
    .category-gauge { aspect-ratio: 1; display: block; position: relative; width: 6.75rem; }
    .category-gauge svg { display: block; height: 100%; width: 100%; }
    .category-gauge .score-gauge-track, .category-gauge .score-gauge-value { stroke-width: 8; }
    .category-gauge .score-gauge-track { stroke: #f1f3f4; }
    .category-gauge .score-gauge-value { stroke: var(--category-color); }
    .category-score { align-items: baseline; display: flex; inset: 0; justify-content: center; position: absolute; top: 2.4rem; }
    .category-score strong { font-size: 1.7rem; font-weight: 500; letter-spacing: -.05em; line-height: 1; } .category-score small { color: var(--muted); font-size: .55rem; font-weight: 600; margin-left: .1rem; }
    .category-title { display: block; font-size: .88rem; font-weight: 600; line-height: 1.3; margin-top: .65rem; }
    .category-score-label { color: var(--category-color); display: block; font-size: .65rem; font-weight: 650; letter-spacing: .05em; margin-top: .35rem; text-transform: uppercase; }
    .category-unscored .category-score-label { color: var(--muted); }
    .category-checks { color: var(--muted); display: block; font-size: .72rem; margin-top: .45rem; }
    .unscored-group { border: 0; border-bottom: 1px solid var(--line); }
    .unscored-group summary { align-items: center; background: var(--surface-soft); display: flex; justify-content: space-between; padding: .7rem .8rem; }
    .unscored-group summary span { color: var(--muted); font-size: .75rem; font-weight: 400; }
    .unscored-grid { display: grid; grid-template-columns: repeat(3, minmax(0, 1fr)); }
    .unscored-link { align-items: center; border-top: 1px solid var(--line); color: inherit; display: grid; gap: .65rem; grid-template-columns: 2rem minmax(0, 1fr); padding: .7rem .8rem; text-decoration: none; }
    .unscored-link:hover { background: #fafafa; }
    .unscored-mark { align-items: center; border: 3px solid #e0e0e0; border-radius: 50%; color: var(--muted); display: flex; font-size: 1rem; height: 2rem; justify-content: center; width: 2rem; }
    .unscored-name { display: block; font-size: .78rem; font-weight: 600; line-height: 1.35; } .unscored-meta { color: var(--muted); display: block; font-size: .68rem; }
    .filter-panel { align-items: end; background: #fafafa; border: 1px solid var(--line); display: grid; gap: .75rem; grid-template-columns: minmax(15rem, 2fr) minmax(12rem, 1fr) minmax(12rem, 1fr) auto; padding: .8rem; position: sticky; top: 2.85rem; z-index: 5; }
    .filter-panel label { color: var(--muted); display: grid; font-size: .72rem; font-weight: 600; gap: .25rem; }
    input, select, button { font: inherit; } input, select { background: white; border: 1px solid #bdbdbd; border-radius: 3px; color: var(--ink); min-height: 2.35rem; padding: .45rem .55rem; width: 100%; }
    button { background: white; border: 1px solid #9e9e9e; border-radius: 3px; color: var(--ink); cursor: pointer; font-weight: 600; min-height: 2.35rem; padding: .45rem .7rem; }
    button:hover { background: #f1f3f4; }
    .filter-count { color: var(--muted); font-size: .78rem; font-weight: 600; margin: .65rem .1rem; }
    .control-row summary code, .assertion-row summary code { margin: 0 .2rem; }
    .technical { background: #fafafa; } .technical table { border-collapse: collapse; width: 100%; }
    th, td { border-bottom: 1px solid var(--line); padding: .65rem; text-align: left; vertical-align: top; }
    th { background: #eef1f4; color: var(--muted); font-weight: 600; } caption { font-size: 1rem; font-weight: 600; margin: 1rem 0; text-align: left; }
    .footer-note { border-top: 1px solid var(--line); color: var(--muted); margin-top: 2rem; padding-top: 1.2rem; }
    .back-top { display: none; }
    [hidden] { display: none !important; }
    @media (max-width: 900px) { .category-grid { grid-template-columns: repeat(3, 1fr); } .unscored-grid { grid-template-columns: repeat(2, 1fr); } .filter-panel { grid-template-columns: 1fr 1fr; position: static; } }
    @media (max-width: 600px) { body { font-size: 13px; } .hero { padding: 0 1rem; } .report-brand { margin-bottom: 1.1rem; } .hero-score-gauge { width: 8.5rem; } .score-gauge-center { top: 2.85rem; } .score-gauge-number { font-size: 2rem; } .hero-metric { padding-left: .2rem; padding-right: .2rem; } .hero-metric strong { font-size: 1.05rem; } .hero-metric span { font-size: .62rem; } .filter-panel { grid-template-columns: 1fr; } .category-grid { grid-template-columns: repeat(2, 1fr); } .category-card { padding-left: .45rem; padding-right: .45rem; } .category-gauge { width: 5.75rem; } .category-score { top: 2rem; } .finding-card > summary { align-items: start; grid-template-columns: minmax(0, 1fr) auto; } .finding-card .severity-pill { grid-column: 1 / -1; grid-row: 1; width: max-content; } .finding-card .finding-summary-copy { grid-column: 1; grid-row: 2; } .finding-card > summary::after { grid-column: 2; grid-row: 2; } .unscored-grid { grid-template-columns: 1fr; } .section-disclosure > summary { align-items: flex-start; } .assertion-row > summary, .control-row > summary { flex-wrap: wrap; } dl { grid-template-columns: 1fr; } dt { margin-top: .45rem; } }
    @media print { .report-nav, .filter-panel, .back-top, button { display: none !important; } .hero { border-top: 0; padding: 1rem 0; } main { padding: 1rem 0; } details { break-inside: avoid; } }
  </style>
</head>
<body class="tone-{{.Local.Tone}}" id="top">
  <a class="skip-link" href="#main-content">Skip to report</a>
  <header class="hero">
    <div class="hero-inner">
      <div class="report-brand"><strong>Production Readiness Checklist</strong><span>Scan report</span></div>
      <div class="hero-heading">
        <h1>{{.Run.Inventory.TargetName}}</h1>
      </div>
      <div class="hero-grid">
        <section class="score-panel" aria-labelledby="local-score-title">
          <div class="score-gauge hero-score-gauge" role="img" aria-label="Overall local score: {{.Local.Percentage}} out of 100">
            <svg viewBox="0 0 120 120" aria-hidden="true"><circle class="score-gauge-track" cx="60" cy="60" r="50"></circle><circle class="score-gauge-value" cx="60" cy="60" r="50" pathLength="100" stroke-dasharray="{{.Local.Percentage}} 100"></circle></svg>
            <span class="score-gauge-center"><strong class="score-gauge-number">{{.Local.Percentage}}</strong><span class="score-gauge-scale">/100</span></span>
          </div>
          <div><span class="score-kicker">Local score</span><span class="rating overall-rating rating-{{.Local.Tone}}" id="local-score-title"><span aria-hidden="true">{{statusSymbol .Local.Tone}}</span> {{.Local.Label}}</span><p class="score-status">{{if .Local.Blocking}}{{.Local.Blocking}} release-blocking {{if eq .Local.Blocking 1}}check{{else}}checks{{end}} failed.{{else}}{{.Local.Explanation}}{{end}}</p></div>
        </section>
      </div>
      <div class="hero-metrics" aria-label="Local check counts"><div class="hero-metric metric-pass"><strong>{{.Local.Passed}}</strong><span>Passed</span></div><div class="hero-metric metric-fail"><strong>{{.Local.Failed}}</strong><span>Failed</span></div><div class="hero-metric metric-review"><strong>{{.Local.ToReview}}</strong><span>Review</span></div><div class="hero-metric"><strong>{{.Local.NotApplicable}}</strong><span>Not needed</span></div></div>
      <details class="summary-help hero-help"><summary>About this score</summary><div class="detail-body">
        <p><strong>{{.Local.Passed}} of {{.Local.Applicable}}</strong> applicable checks passed using <code>{{.Run.Plan.ProfileID}}@{{.Run.Plan.ProfileVersion}}</code>. {{.Local.NotApplicable}} checks did not apply.</p>
        <p>The score is a simple local-check pass rate, not a readiness grade. One serious failure can outweigh many passing checks.</p>
        <p><strong>Report-only scan:</strong> no files were fixed and no project scripts were run.{{if .Run.ControlCatalog}} The wider catalog still has {{.Meaning.ControlsNeedingReview}} controls needing evidence or review.{{end}}</p>
        <p>AI suggestions, when present, are advice rather than verified passes.</p>
      </div></details>
    </div>
  </header>
  <nav class="report-nav" aria-label="Report sections"><div class="nav-inner">
    <a href="#top">Overview</a>{{if .Run.ControlCatalog}}<a href="#categories">Categories</a>{{end}}{{if .AIRecommendationCount}}<a href="#ai-priorities">AI priorities</a>{{end}}<a href="#findings">Fix first</a><a href="#local-checks">Details</a>
  </div></nav>
  <main id="main-content">
    {{if .Run.ControlCatalog}}<section class="section" id="categories">
      <h2>Scores by category</h2>
      <p class="section-intro">Each score uses only the local checks that ran for that category. Select a category to inspect its linked controls.</p>
      <div class="category-key"><strong>Score colors</strong><span class="key-item"><span class="key-dot great"></span>Great</span><span class="key-item"><span class="key-dot good"></span>Good</span><span class="key-item"><span class="key-dot review"></span>Review</span><span class="key-item"><span class="key-dot bad"></span>Needs work</span><span class="key-item"><span class="key-dot"></span>Not scored</span></div>
      <div class="category-grid">{{range .ScoredCategories}}<a class="category-card category-{{.LocalTone}}" href="#control-explorer" data-category-link="{{.Key}}">
        <span class="category-overview">
          <span class="category-gauge" role="img" aria-label="{{.Name}}: {{.LocalPercentage}} out of 100 from linked local checks">
            <svg viewBox="0 0 120 120" aria-hidden="true"><circle class="score-gauge-track" cx="60" cy="60" r="50"></circle><circle class="score-gauge-value" cx="60" cy="60" r="50" pathLength="100" stroke-dasharray="{{.LocalPercentage}} 100"></circle></svg>
            <span class="category-score"><strong>{{.LocalPercentage}}</strong><small>/100</small></span>
          </span>
          <span class="category-title">{{.Name}}</span><span class="category-score-label">{{.LocalLabel}}</span>
          <span class="category-checks">{{.LocalPassed}} / {{.LocalApplicable}} checks passed{{if .LocalNotApplicable}} · {{.LocalNotApplicable}} did not apply{{end}}</span>
        </span>
      </a>{{end}}</div>
      {{if .UnscoredCategories}}<details class="unscored-group"><summary>Not scored in this scan <span>{{len .UnscoredCategories}} categories without an applicable linked local check</span></summary><div class="unscored-grid">{{range .UnscoredCategories}}<a class="unscored-link" href="#control-explorer" data-category-link="{{.Key}}"><span class="unscored-mark" aria-hidden="true">—</span><span><span class="unscored-name">{{.Name}}</span><span class="unscored-meta">{{.Total}} controls · {{.NeedsReview}} need evidence{{if .LocalNotApplicable}} · {{.LocalNotApplicable}} did not apply{{end}}</span></span></a>{{end}}</div></details>{{end}}
    </section>{{end}}

    {{if .Run.AIImprovementPlan}}<section class="section" id="ai-plan">
      <h2>AI review improvement plan</h2>
      <p class="section-intro">The scanner grouped reviews only when their category, cause key, and normalized cause matched exactly. This removes repeated work without turning AI advice into a verified finding.</p>
      <p class="notice"><strong>{{.Run.AIImprovementPlan.ReviewedControlCount}} controls became {{.AIPlanItemCount}} cause groups.</strong> The first {{len .AIPlanItems}} groups are shown. Effort and reach are AI estimates; use the listed checks before accepting a change.</p>
      {{range .AIPlanItems}}<details class="finding-card severity-{{severityClass .Item.Priority}}">
        <summary><span class="pill severity-pill severity-{{severityClass .Item.Priority}}">{{prettyStatus .Item.Priority}}</span><span class="finding-summary-copy"><span class="finding-title">{{.Item.RootCause}}</span><span class="finding-preview">{{.Item.ControlCount}} linked control{{if gt .Item.ControlCount 1}}s{{end}} · {{.Item.Domain}}</span></span></summary>
        <div class="finding-body">
          <p class="primary-detail"><strong>Estimated work:</strong> {{prettyStatus .Item.Effort}} effort · {{prettyStatus .Item.BlastRadius}} reach.</p>
          <p class="primary-detail"><strong>Linked controls:</strong> {{if gt (len .Item.ControlIDs) 8}}<code>{{join (slice .Item.ControlIDs 0 8)}}</code> and {{sub (len .Item.ControlIDs) 8}} more{{else}}<code>{{join .Item.ControlIDs}}</code>{{end}}</p>
          {{if .Lead.AIReview}}<p class="primary-detail"><strong>Start with the advice for {{.Lead.ControlID}}:</strong> {{.Lead.AIReview.Advice}}</p>
          <p class="primary-detail"><strong>Suggested work:</strong></p><ol>{{range .Lead.AIReview.RemediationSteps}}<li>{{.}}</li>{{end}}</ol>
          <p class="primary-detail"><strong>How to verify:</strong></p><ol>{{range .Lead.AIReview.VerificationSteps}}<li>{{.}}</li>{{end}}</ol>{{end}}
          <details class="finding-details"><summary>Why this group is advisory</summary><div class="detail-body"><p>The cause key came from AI. The scanner used exact matching only and did not decide that differently worded causes are the same. Open the linked controls for their full evidence, objections, and limits.</p><p><code>{{.Item.ItemID}}</code></p></div></details>
        </div>
      </details>{{end}}
    </section>
    {{else if .AIRecommendationCount}}<section class="section" id="ai-priorities">
      <h2>AI review priorities</h2>
      <p class="section-intro">These are the highest-priority suggestions from the optional AI review, sorted by its risk estimate. They are a work plan, not verified findings. Confirm the evidence and run the listed checks before accepting a change.</p>
      <p class="notice"><strong>{{.AIRecommendationCount}} controls received a non-empty AI priority.</strong> The first {{len .AIRecommendations}} are shown here; every review remains searchable in the complete control catalog.</p>
      {{range .AIRecommendations}}<details class="finding-card severity-{{severityClass .AIReview.Priority}}">
        <summary><span class="pill severity-pill severity-{{severityClass .AIReview.Priority}}">{{prettyStatus .AIReview.Priority}}</span><span class="finding-summary-copy"><span class="finding-title">{{.ControlID}} — {{brief .Statement}}</span><span class="finding-preview">{{brief .AIReview.Advice}}</span></span></summary>
        <div class="finding-body">
          <p class="primary-detail"><strong>Why it was flagged:</strong> {{.AIReview.Reason}}</p>
          <p class="primary-detail"><strong>Risk if ignored:</strong> {{.AIReview.RiskIfIgnored}}</p>
          <p class="primary-detail"><strong>Suggested work:</strong></p><ol>{{range .AIReview.RemediationSteps}}<li>{{.}}</li>{{end}}</ol>
          <p class="primary-detail"><strong>How to verify:</strong></p><ol>{{range .AIReview.VerificationSteps}}<li>{{.}}</li>{{end}}</ol>
          <details class="finding-details"><summary>Evidence, challenge, and limits</summary><div class="detail-body"><dl>
            <dt>Evidence still needed</dt><dd><ul>{{range .AIReview.EvidenceNeeded}}<li>{{.}}</li>{{end}}</ul></dd>
            <dt>Skeptical challenge</dt><dd>{{.AIReview.Challenge}}</dd>
            <dt>Candidate result</dt><dd>{{prettyStatus .AIReview.AssessmentCandidate}} / {{prettyStatus .AIReview.ApplicabilityCandidate}} / confidence {{prettyStatus .AIReview.Confidence}}</dd>
            <dt>Cited lines</dt><dd>{{if .AIReview.Evidence}}<ul>{{range .AIReview.Evidence}}<li><code>{{location .}}</code></li>{{end}}</ul>{{else}}No repository line was cited.{{end}}</dd>
            <dt>Limits</dt><dd><ul>{{range .AIReview.Limitations}}<li>{{.}}</li>{{end}}</ul></dd>
          </dl></div></details>
        </div>
      </details>{{end}}
    </section>{{end}}

    <section class="section" id="findings">
      <h2>What to fix first</h2>
      <p class="section-intro">Verified problems are sorted from most to least serious. Each row shows only the decision-making summary; open it for the full explanation and next action.</p>
      {{if .Findings}}<div class="severity-key" aria-label="Severity colors"><strong>Severity</strong><span class="key-item severity-critical"><span class="key-dot"></span>Critical</span><span class="key-item severity-high"><span class="key-dot"></span>High</span><span class="key-item severity-medium"><span class="key-dot"></span>Medium</span><span class="key-item severity-low"><span class="key-dot"></span>Low</span><span class="key-item severity-info"><span class="key-dot"></span>Info</span></div>
      {{range .Findings}}<details class="finding-card severity-{{severityClass .Severity}}">
        <summary><span class="pill severity-pill severity-{{severityClass .Severity}}"><span aria-hidden="true">×</span> {{prettyStatus .Severity}} · {{prettyStatus .Gate}}</span><span class="finding-summary-copy"><span class="finding-title">{{.Title}}</span><span class="finding-preview">{{brief .Summary}}</span></span></summary>
        <div class="finding-body">
        <p class="primary-detail"><strong>What failed:</strong> {{.Summary}}</p>
        <p class="next-action"><strong>Next action:</strong> {{remediation .RemediationClass}}</p>
        <details class="finding-details"><summary>Technical details</summary><div class="detail-body"><dl>
          <dt>Full scanner message</dt><dd>{{.Summary}}</dd>
          <dt>Severity / gate</dt><dd>{{prettyStatus .Severity}} / {{prettyStatus .Gate}}</dd>
          <dt>Remediation class</dt><dd>{{.RemediationClass}}</dd>
          <dt>Assertion</dt><dd><code>{{.AssertionID}}</code></dd>
          <dt>Control IDs</dt><dd>{{if .ControlIDs}}<code>{{join .ControlIDs}}</code>{{else}}—{{end}}</dd>
          <dt>Locations</dt><dd>{{if .Locations}}{{if gt (len .Locations) 5}}<ul>{{range slice .Locations 0 5}}<li><code>{{location .}}</code></li>{{end}}</ul><details class="raw-list"><summary>Show all {{len .Locations}} locations</summary><ul>{{range .Locations}}<li><code>{{location .}}</code></li>{{end}}</ul></details>{{else}}<ul>{{range .Locations}}<li><code>{{location .}}</code></li>{{end}}</ul>{{end}}{{else}}No source location was emitted.{{end}}</dd>
          <dt>Evidence IDs</dt><dd>{{if .EvidenceIDs}}<ul>{{range .EvidenceIDs}}<li><code>{{.}}</code></li>{{end}}</ul>{{else}}No evidence identifier was attached.{{end}}</dd>
          <dt>Finding ID</dt><dd><code>{{.ID}}</code></dd><dt>Stable fingerprint</dt><dd><code>{{.Fingerprint}}</code></dd>
        </dl></div></details>
        </div>
      </details>{{end}}{{else}}<p class="notice"><strong>No verified failures were found.</strong> Still review incomplete, manual, and complete-catalog results below.</p>{{end}}
    </section>

    <details class="section-disclosure" id="local-checks">
      <summary><span><span class="disclosure-title">Local check details</span><span class="disclosure-note">Open only when you need the exact check results or raw evidence.</span></span><span class="result-count">{{len .Run.Results}} checks</span></summary>
      <div class="section-disclosure-body">
        <h2>Checks needing attention</h2>
        <p class="section-intro">Failures, blocked checks, and manual decisions appear first. Open one result for its plain explanation; technical evidence is one more optional step inside.</p>
        {{if and (eq .Local.Failed 0) (eq .Local.Unresolved 0) (eq .Local.Manual 0)}}<p class="notice"><strong>No local check needs attention.</strong></p>{{end}}
        {{range .AttentionResults}}{{template "assertion-row" .}}{{end}}
        {{if .Local.Passed}}<details class="result-group"><summary><span>Passed checks</span><span class="result-count">{{.Local.Passed}}</span></summary><div>{{range .Run.Results}}{{if eq .Assessment "pass"}}{{template "assertion-row" .}}{{end}}{{end}}</div></details>{{end}}
        {{if .Local.NotApplicable}}<details class="result-group"><summary><span>Not applicable</span><span class="result-count">{{.Local.NotApplicable}}</span></summary><div>{{range .Run.Results}}{{if eq .Assessment "not_applicable"}}{{template "assertion-row" .}}{{end}}{{end}}</div></details>{{end}}
        <p class="notice"><strong>Evidence time:</strong> observation time records when evidence was collected. It does not guarantee that evidence is still current.</p>
      </div>
    </details>

    {{if .Run.ControlCatalog}}<details class="section-disclosure" id="control-explorer">
      <summary><span><span class="disclosure-title">Search the complete control catalog</span><span class="disclosure-note">Use this when you need one specific rule, category, or evidence state.</span></span><span class="result-count">{{.Run.ControlCatalog.ControlCount}} controls</span></summary>
      <div class="section-disclosure-body">
      <h2>Search all controls</h2>
      <p class="section-intro">Search by plain words or ID, then narrow the list by category or evidence state. Only the first 100 matches are shown at once.</p>
      <div class="filter-panel" aria-label="Control filters">
        <label>Search<input id="control-search" type="search" placeholder="Try: backups, authentication, PRC-02-003" autocomplete="off"></label>
        <label>Category<select id="control-category"><option value="">All categories</option>{{range .Categories}}<option value="{{.Key}}">{{.Name}}</option>{{end}}</select></label>
        <label>Evidence state<select id="control-disposition"><option value="">All states</option><option value="confirmed_failure">Confirmed failure</option><option value="verified_pass">Verified pass</option><option value="not_applicable">Not applicable</option><option value="blocked">Blocked</option><option value="needs_review">Needs evidence</option><option value="partially_verified">Partial evidence</option><option value="retired">Retired</option></select></label>
        <button id="clear-control-filters" type="button">Clear</button>
      </div>
      <p class="filter-count" id="control-filter-count" role="status" aria-live="polite"></p>
      <section id="control-results">{{range .Run.ControlResults}}<details class="control-row result-{{statusClass .Disposition}}" data-disposition="{{.Disposition}}" data-classification="{{.Classification}}" data-category="{{categoryKey .Source.Path}}">
        <summary><span class="pill pill-{{statusClass .Disposition}}"><span aria-hidden="true">{{controlSymbol .Disposition}}</span> {{prettyStatus .Disposition}}</span><span class="row-summary"><code>{{.ControlID}}</code> — {{brief .Statement}}</span></summary>
        <div class="detail-body">
          <p class="primary-detail"><strong>What this control asks:</strong> {{.Statement}}</p>
          <p class="primary-detail"><strong>Current result:</strong> {{.Summary}}</p>
          {{if .AIReview}}<div class="advice"><strong>Advisory AI suggestion · {{prettyStatus .AIReview.Priority}} priority</strong><p>{{.AIReview.Reason}}</p><p><strong>Suggested work:</strong> {{.AIReview.Advice}}</p><ol>{{range .AIReview.RemediationSteps}}<li>{{.}}</li>{{end}}</ol><details class="raw-list"><summary>How to verify and what evidence is missing</summary><div><p><strong>Verification:</strong></p><ol>{{range .AIReview.VerificationSteps}}<li>{{.}}</li>{{end}}</ol><p><strong>Evidence needed:</strong></p><ul>{{range .AIReview.EvidenceNeeded}}<li>{{.}}</li>{{end}}</ul><p><strong>Skeptical challenge:</strong> {{.AIReview.Challenge}}</p></div></details><span class="meta">This is advice only and cannot create a verified pass.</span></div>{{end}}
          <details class="raw-details"><summary>Technical evidence and IDs</summary><div class="detail-body"><dl>
          <dt>Category</dt><dd>{{categoryName .Source.Path}}</dd>
          <dt>Evidence state</dt><dd>{{prettyStatus .Disposition}} — {{.Summary}}</dd>
          <dt>Reviewed classification</dt><dd>{{prettyStatus .Classification}} — {{prettyStatus .ClassificationRoute}} / {{prettyStatus .ClassificationDecisionBasis}}</dd>
          <dt>Classification row digest</dt><dd>{{if .ClassificationRowSHA256}}<code>{{.ClassificationRowSHA256}}</code>{{else}}Legacy record.{{end}}</dd>
          <dt>Deterministic binding</dt><dd>{{if .DeterministicBindingID}}<code>{{.DeterministicBindingID}}</code> / <code>{{.DeterministicBindingSHA256}}</code>{{else}}No deterministic binding; review remains contextual or accountable.{{end}}</dd>
          <dt>Coverage / authority</dt><dd>{{prettyStatus .Coverage}} / {{prettyStatus .Authority}}</dd>
          <dt>Source</dt><dd><code>{{.Source.Path}}:{{.Source.Line}}</code></dd>
          <dt>All known assertions</dt><dd>{{if .AssertionIDs}}<code>{{join .AssertionIDs}}</code>{{else}}None yet.{{end}}</dd>
          <dt>Assertions run in this profile</dt><dd>{{if .ExecutedAssertionIDs}}<code>{{join .ExecutedAssertionIDs}}</code>{{else}}None.{{end}}</dd>
          {{if .DeterministicClauseResults}}<dt>Exact deterministic programs</dt><dd><details class="raw-list"><summary>Show {{len .DeterministicClauseResults}} clause result(s)</summary><ul>{{range .DeterministicClauseResults}}<li><code>{{.ImplementationID}}</code> · {{prettyStatus .Status}} · authority <code>{{.RequiredAuthority}}</code>{{if .ReasonCode}} · reason <code>{{.ReasonCode}}</code>{{end}}{{if .EvidenceSHA256}} · evidence <code>{{.EvidenceSHA256}}</code>{{end}}</li>{{end}}</ul></details></dd>{{end}}
          {{if .AIReview}}<dt>Advisory AI review</dt><dd><strong>{{.AIReview.Provider}}{{if .AIReview.Model}} / {{.AIReview.Model}}{{end}}</strong>: {{.AIReview.ReviewDepth}} depth / {{.AIReview.AssessmentCandidate}} / {{.AIReview.ApplicabilityCandidate}} / confidence {{.AIReview.Confidence}} / priority {{.AIReview.Priority}}.</dd>
          <dt>AI risk if ignored</dt><dd>{{.AIReview.RiskIfIgnored}}</dd>
          <dt>AI citation locations</dt><dd>{{if .AIReview.Evidence}}<ul>{{range .AIReview.Evidence}}<li><code>{{location .}}</code></li>{{end}}</ul>{{else}}No repository line was cited.{{end}}</dd>
          <dt>AI verification state</dt><dd><code>{{verification .AIReview}}</code>. A valid snapshot location proves only that the cited line existed in the screened input. It does not prove that the line supports the AI claim; claim text remains advisory and unverified.</dd>
          <dt>AI review limits</dt><dd><ul>{{range .AIReview.Limitations}}<li>{{.}}</li>{{end}}</ul></dd>{{end}}
          </dl></div></details>
        </div>
      </details>{{end}}</section>
      <button class="more-controls" id="show-more-controls" type="button" hidden>Show more controls</button>
      </div>
    </details>{{end}}

    <details class="section-disclosure" id="technical">
      <summary><span><span class="disclosure-title">Scan metadata</span><span class="disclosure-note">Run IDs, hashes, catalog identity, and adapter records.</span></span><span class="result-count">Technical</span></summary>
      <div class="section-disclosure-body">
      <h2>Scan metadata</h2>
      <p class="section-intro">This information makes the report repeatable and auditable. Most readers do not need it.</p>
      <details class="technical"><summary>Run, source, and catalog identity</summary><div class="detail-body"><dl>
        <dt>Run</dt><dd><code>{{.Run.RunID}}</code></dd><dt>Profile</dt><dd><code>{{.Run.Plan.ProfileID}}@{{.Run.Plan.ProfileVersion}}</code></dd>
        <dt>Target</dt><dd>{{.Run.Inventory.TargetName}}</dd><dt>Inventory</dt><dd><code>{{.Run.Inventory.Digest}}</code></dd>
        {{if .Run.Plan.ConfigurationDigest}}<dt>Configuration</dt><dd><code>{{.Run.Plan.ConfigurationDigest}}</code></dd><dt>Project</dt><dd><code>{{.Run.Plan.ProjectID}}</code></dd><dt>Target environments</dt><dd>{{join .Run.Plan.TargetEnvironments}}</dd><dt>Artifact digests</dt><dd>{{join .Run.Plan.ArtifactDigests}}</dd>{{end}}
        <dt>Local profile result</dt><dd>{{.ProfileState}}</dd><dt>Full catalog coverage result</dt><dd>{{.Run.TerminalState}}</dd>
        {{if .Run.Inventory.GitCommit}}<dt>Git HEAD</dt><dd><code>{{.Run.Inventory.GitCommit}}</code> ({{.GitState}} worktree)</dd>{{end}}
        <dt>Hashed inventory bytes</dt><dd>{{.InventoryBytes}}</dd><dt>Reported automatic exclusions</dt><dd>{{.ExclusionCount}}</dd>
        {{if .Run.ControlCatalog}}<dt>Registry version / digest</dt><dd><code>{{.Run.ControlCatalog.RegistryVersion}}</code> / <code>{{.Run.ControlCatalog.RegistrySHA256}}</code></dd><dt>Catalog source digest</dt><dd><code>{{.Run.ControlCatalog.SourceSHA256}}</code></dd><dt>Controls</dt><dd>{{.Run.ControlCatalog.ControlCount}} registered / {{.Run.ControlCatalog.ActiveControlCount}} active</dd>{{if .Run.ControlCatalog.ClassificationCorpusSHA256}}<dt>Reviewed classifications</dt><dd>{{.Run.ControlCatalog.ReviewedDeterministicCount}} deterministic / {{.Run.ControlCatalog.ReviewedNondeterministicCount}} nondeterministic</dd><dt>Classification corpus digest</dt><dd><code>{{.Run.ControlCatalog.ClassificationCorpusSHA256}}</code></dd><dt>Deterministic bindings</dt><dd>{{.Run.ControlCatalog.DeterministicBindingCount}} / <code>{{.Run.ControlCatalog.ControlCheckBindingsSHA256}}</code></dd><dt>Exact deterministic execution</dt><dd>{{.Run.ControlCatalog.DeterministicProgramExecutedCount}} of {{.Run.ControlCatalog.DeterministicProgramTemplateCount}} programs attempted · {{.Run.ControlCatalog.DeterministicProgramPassCount}} passed · {{.Run.ControlCatalog.DeterministicProgramFailCount}} failed · {{.Run.ControlCatalog.DeterministicProgramNACount}} proven Not Applicable · {{.Run.ControlCatalog.DeterministicProgramBlockedCount}} controls still blocked · {{len .Run.DeterministicEvidence}} replayable evidence documents retained</dd>{{end}}{{end}}
      </dl></div></details>
      {{if .Run.AdapterExecutions}}<details class="technical"><summary>Adapter executions</summary><div class="detail-body"><table>
        <caption>Adapter executions</caption><thead><tr><th scope="col">Adapter</th><th scope="col">Manifest</th><th scope="col">Authorization</th><th scope="col">Trust</th><th scope="col">Registry</th><th scope="col">Status</th><th scope="col">Execution</th></tr></thead>
        <tbody>{{range .Run.AdapterExecutions}}<tr><td><code>{{.AdapterID}}</code></td><td><code>{{.ManifestSHA256}}</code></td><td>{{.Resolution.Source}}</td><td>{{.Resolution.Trust}}</td><td><code>{{.Resolution.RegistryID}}</code></td><td>{{.Transcript.Summary.Status}}</td><td><code>{{.ExecutionID}}</code></td></tr>{{end}}</tbody>
      </table></div></details>{{end}}
      {{if .Run.AuthoritativeEvidence}}<details class="technical"><summary>Signed authoritative evidence</summary><div class="detail-body">
        {{range .Run.AuthoritativeEvidence}}<details class="raw-list"><summary><code>{{.BundleID}}</code> · {{.EntryCount}} {{.Authority}} entries</summary><div><dl>
		  <dt>Bundle digest</dt><dd><code>{{.BundleSHA256}}</code></dd>
		  <dt>Pre-collection policy digest</dt><dd><code>{{.PolicySHA256}}</code></dd>
          <dt>Catalog / inventory</dt><dd><code>{{.CatalogSHA256}}</code> / <code>{{.InventorySHA256}}</code></dd>
          <dt>Independent policy signature</dt><dd>key <code>{{.PolicySignature.KeyID}}</code> · verified {{.PolicySignature.VerifiedAt}}</dd>
          <dt>Authority evidence signature</dt><dd>key <code>{{.EvidenceSignature.KeyID}}</code> · verified {{.EvidenceSignature.VerifiedAt}}</dd>
          <dt>Trust store</dt><dd><code>{{.PolicySignature.TrustStoreID}}</code> / <code>{{.PolicySignature.TrustStoreDigest}}</code></dd>
        </dl></div></details>{{end}}
      </div></details>{{end}}
      {{if .Run.DeterministicEvidence}}<details class="technical"><summary>Replayable exact evidence</summary><div class="detail-body">
        {{range .Run.DeterministicEvidence}}<details class="raw-list"><summary><code>{{.ControlID}}</code> · clause <code>{{.ClauseID}}</code> · {{.Applicability}}</summary><div><dl>
          <dt>Evidence digest</dt><dd><code>{{exactEvidenceSHA .}}</code></dd>
          <dt>Evidence / program IDs</dt><dd><code>{{.EvidenceID}}</code> / <code>{{.ProgramSHA256}}</code></dd>
          <dt>Authority / observed</dt><dd>{{.Authority}} / {{.ObservedAt}}</dd>
          <dt>Inventory / subject</dt><dd><code>{{.InventorySHA256}}</code> / <code>{{.SubjectID}}</code></dd>
          <dt>Normalized facts</dt><dd><ul>{{range $key, $fact := .Facts}}<li><code>{{$key}}</code>: <code>{{factJSON $fact}}</code></li>{{end}}</ul></dd>
        </dl></div></details>{{end}}
      </div></details>{{end}}
      <p class="footer-note">This report is scoped to the named profile, target inventory, and evidence set. It is not an unqualified production-readiness or compliance claim.</p>
      </div>
    </details>
  </main>
  <script>
  (() => {
    for (const link of document.querySelectorAll('a[href^="#"]')) {
      link.addEventListener('click', event => {
        const target = document.querySelector(link.getAttribute('href'));
        if (!target) return;
        event.preventDefault();
        if (target instanceof HTMLDetailsElement) target.open = true;
        history.pushState(null, '', link.getAttribute('href'));
        const settle = () => window.scrollTo(0, target.getBoundingClientRect().top + window.scrollY - 48);
        settle();
        requestAnimationFrame(() => { settle(); setTimeout(settle, 80); });
      });
    }
    const root = document.getElementById('control-results');
    if (!root) return;
    const items = Array.from(root.querySelectorAll(':scope > details'));
    const search = document.getElementById('control-search');
    const category = document.getElementById('control-category');
    const disposition = document.getElementById('control-disposition');
    const clear = document.getElementById('clear-control-filters');
    const count = document.getElementById('control-filter-count');
    const more = document.getElementById('show-more-controls');
    const explorer = document.getElementById('control-explorer');
    const pageSize = 100;
    let limit = pageSize;
    let timer;
    const apply = () => {
      const query = search.value.trim().toLocaleLowerCase();
      const matches = [];
      for (const item of items) {
        const matchesFilter = (!query || item.textContent.toLocaleLowerCase().includes(query)) &&
          (!category.value || item.dataset.category === category.value) &&
          (!disposition.value || item.dataset.disposition === disposition.value);
        if (matchesFilter) matches.push(item);
      }
      for (const item of items) item.hidden = true;
      const visible = matches.slice(0, limit);
      for (const item of visible) item.hidden = false;
      count.textContent = 'Showing ' + visible.length.toLocaleString() + ' of ' + matches.length.toLocaleString() + ' matching controls';
      const remaining = matches.length - visible.length;
      more.hidden = remaining <= 0;
      if (remaining > 0) more.textContent = 'Show ' + Math.min(pageSize, remaining).toLocaleString() + ' more';
    };
    const resetAndApply = () => { limit = pageSize; apply(); };
    const schedule = () => { clearTimeout(timer); timer = setTimeout(resetAndApply, 120); };
    search.addEventListener('input', schedule);
    category.addEventListener('change', resetAndApply);
    disposition.addEventListener('change', resetAndApply);
    clear.addEventListener('click', () => { search.value = ''; category.value = ''; disposition.value = ''; resetAndApply(); search.focus(); });
    more.addEventListener('click', () => { limit += pageSize; apply(); });
    for (const link of document.querySelectorAll('[data-category-link]')) {
      link.addEventListener('click', () => { explorer.open = true; category.value = link.dataset.categoryLink; search.value = ''; disposition.value = ''; resetAndApply(); });
    }
    apply();
  })();
  </script>
</body>
</html>
{{define "assertion-row"}}<details class="assertion-row result-{{statusClass .Assessment}}">
  <summary><span class="pill pill-{{statusClass .Assessment}}"><span aria-hidden="true">{{assessmentSymbol .Assessment}}</span> {{prettyStatus .Assessment}}</span><span class="row-summary"><code>{{.AssertionID}}</code> — {{brief .Summary}}</span><span class="severity-label severity-{{severityClass .Severity}}">{{prettyStatus .Severity}}</span></summary>
  <div class="detail-body">
    <p class="primary-detail"><strong>Result:</strong> {{.Summary}}</p>
    <p class="primary-detail"><strong>Suggested next step:</strong> {{remediation .RemediationClass}}</p>
    <details class="raw-details"><summary>Technical evidence and IDs</summary><div class="detail-body"><dl>
      <dt>Applicability / execution</dt><dd>{{prettyStatus .Applicability}} / {{prettyStatus .Execution}}</dd>
      <dt>Severity / gate</dt><dd>{{prettyStatus .Severity}} / {{prettyStatus .Gate}}</dd>
      <dt>Control IDs</dt><dd>{{if .ControlIDs}}<code>{{join .ControlIDs}}</code>{{else}}—{{end}}</dd>
      <dt>Remediation class</dt><dd>{{.RemediationClass}}</dd>
      <dt>Locations</dt><dd>{{if .Locations}}{{if gt (len .Locations) 5}}<ul>{{range slice .Locations 0 5}}<li><code>{{location .}}</code></li>{{end}}</ul><details class="raw-list"><summary>Show all {{len .Locations}} locations</summary><ul>{{range .Locations}}<li><code>{{location .}}</code></li>{{end}}</ul></details>{{else}}<ul>{{range .Locations}}<li><code>{{location .}}</code></li>{{end}}</ul>{{end}}{{else}}No source location was emitted.{{end}}</dd>
      <dt>Required evidence</dt><dd>{{if .EvidenceRequired}}<ul>{{range .EvidenceRequired}}<li><strong>{{.Kind}}</strong> (minimum authority: {{.MinimumAuthority}}): {{.Description}}</li>{{end}}</ul>{{else}}No evidence requirement was declared.{{end}}</dd>
      <dt>Observed evidence</dt><dd>{{if .EvidenceObserved}}<ul>{{range .EvidenceObserved}}<li><strong>{{.Kind}}</strong> from <code>{{.Source}}</code>: {{.Summary}} <span class="meta">(authority: {{.Authority}}; observed: {{observedAt .}}; ID: <code>{{.ID}}</code>)</span></li>{{end}}</ul>{{else}}No evidence was observed.{{end}}</dd>
    </dl></div></details>
  </div>
</details>{{end}}
`

func writeHTML(output io.Writer, run model.RunResult) error {
	categories := summarizeControlCategories(run)
	aiRecommendations, aiRecommendationCount := sortedAIRecommendations(run.ControlResults)
	aiPlanItems, aiPlanItemCount := visibleAIPlan(run)
	scoredCategories := make([]controlCategorySummary, 0, len(categories))
	unscoredCategories := make([]controlCategorySummary, 0, len(categories))
	for _, category := range categories {
		if category.HasLocalScore {
			scoredCategories = append(scoredCategories, category)
		} else {
			unscoredCategories = append(unscoredCategories, category)
		}
	}
	view := struct {
		Run                   model.RunResult
		Findings              []model.Finding
		AttentionResults      []model.AssertionResult
		AIRecommendations     []model.ControlResult
		AIRecommendationCount int
		AIPlanItems           []aiPlanView
		AIPlanItemCount       int
		Meaning               resultMeaning
		Local                 LocalCheckSummary
		Categories            []controlCategorySummary
		ScoredCategories      []controlCategorySummary
		UnscoredCategories    []controlCategorySummary
		ProfileState          string
		GitState              string
		InventoryBytes        string
		ExclusionCount        int
	}{
		Run: run, Findings: sortedFindings(run.Findings), AttentionResults: sortedAttentionResults(run.Results),
		AIRecommendations: aiRecommendations, AIRecommendationCount: aiRecommendationCount,
		AIPlanItems: aiPlanItems, AIPlanItemCount: aiPlanItemCount,
		Meaning: summarizeMeaning(run), Local: SummarizeLocalChecks(run),
		Categories: categories, ScoredCategories: scoredCategories, UnscoredCategories: unscoredCategories,
		ProfileState:   profileTerminalState(run),
		GitState:       inventoryFact(run.Inventory, "repository.git_worktree_state"),
		InventoryBytes: inventoryFact(run.Inventory, "repository.inventory_bytes"),
		ExclusionCount: inventoryFactCount(run.Inventory, "repository.exclusion"),
	}
	return template.Must(template.New("report").Funcs(template.FuncMap{
		"join":         func(values []string) string { return strings.Join(values, ", ") },
		"sub":          func(left, right int) int { return left - right },
		"brief":        func(value string) string { return briefText(value, 150) },
		"verification": advisoryVerificationText,
		"prettyStatus": prettyStatus,
		"statusClass":  statusClass,
		"exactEvidenceSHA": func(evidence controlprogram.Evidence) string {
			return controlprogram.EvidenceSHA256(evidence)
		},
		"factJSON": func(fact controlprogram.Fact) string {
			data, _ := json.Marshal(fact)
			return string(data)
		},
		"severityClass": func(value string) string {
			switch strings.ToLower(value) {
			case "critical", "high", "medium", "low", "info":
				return strings.ToLower(value)
			default:
				return "unknown"
			}
		},
		"statusSymbol":     statusSymbol,
		"assessmentSymbol": statusSymbol,
		"controlSymbol":    statusSymbol,
		"categoryKey": func(sourcePath string) string {
			key, _ := controlCategory(sourcePath)
			return key
		},
		"categoryName": func(sourcePath string) string {
			_, name := controlCategory(sourcePath)
			return name
		},
		"observedAt": func(value model.Evidence) string {
			if value.ObservedAt.IsZero() {
				return "not recorded in this legacy/test record"
			}
			return value.ObservedAt.UTC().Format(time.RFC3339)
		},
		"location": func(value model.FindingLocation) string {
			result := value.Path
			if value.Line > 0 {
				result += ":" + strconv.Itoa(value.Line)
				if value.Column > 0 {
					result += ":" + strconv.Itoa(value.Column)
				}
			}
			return result
		},
		"remediation": func(class string) string {
			switch class {
			case "R0":
				return "Review and resolve this manually; the scanner does not author a change for this class."
			case "R1":
				return "A deterministic, behavior-preserving candidate may be available through the separate prc fix workflow."
			case "R2":
				return "Use an isolated agent-authored candidate only with independent deterministic verification."
			case "R3":
				return "Treat this as a dependency or build-behavior change requiring explicit opt-in and stronger verification."
			case "R4":
				return "Infrastructure or deployment-definition changes require human authorization and environment-specific validation."
			case "R5":
				return "External staging-system changes require a separate connector policy and approval."
			case "R6":
				return "Production, destructive, legal, financial, and risk decisions are prohibited in the general remediation loop."
			default:
				return "Review the finding and its evidence before selecting a remediation path."
			}
		},
	}).Parse(htmlReport)).Execute(output, view)
}

func profileTerminalState(run model.RunResult) string {
	if run.ControlCatalog != nil && run.ControlCatalog.ProfileTerminalState != "" {
		return run.ControlCatalog.ProfileTerminalState
	}
	return run.TerminalState
}

func inventoryFact(item model.Inventory, key string) string {
	for _, fact := range item.Facts {
		if fact.Key == key {
			return fact.Value
		}
	}
	return "not available"
}

func inventoryFactCount(item model.Inventory, key string) int {
	count := 0
	for _, fact := range item.Facts {
		if fact.Key == key {
			count++
		}
	}
	return count
}
