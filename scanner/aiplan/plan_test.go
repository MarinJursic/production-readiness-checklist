package aiplan

import (
	"strings"
	"testing"

	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
)

func TestBuildGroupsOnlyExactDomainKeyAndCause(t *testing.T) {
	digest := strings.Repeat("a", 64)
	review := func(taskID, cause, key, priority string) *model.AIControlReview {
		return &model.AIControlReview{
			Provider: "codex", ReviewDepth: "deep", AssessmentCandidate: "needs_evidence",
			Priority: priority, RootCause: cause, RootCauseKey: key, Effort: "medium",
			BlastRadius: "component", TaskID: taskID,
		}
	}
	run := model.RunResult{
		RunID: digest,
		ControlCatalog: &model.ControlCatalogSummary{
			AIReviewProvider: "codex", AIReviewDepth: "deep", AIReviewState: "complete",
		},
		ControlResults: []model.ControlResult{
			{ControlID: "PRC-01-001", Source: model.Source{Path: "docs/engineering/08-security-and-cryptography.md"}, AIReview: review(digest, "Ownership is unclear.", "ownership-unclear", "medium")},
			{ControlID: "PRC-01-002", Source: model.Source{Path: "docs/engineering/08-security-and-cryptography.md"}, AIReview: review(digest, "  Ownership   is unclear. ", "ownership-unclear", "high")},
			{ControlID: "PRC-01-003", Source: model.Source{Path: "docs/checklists/05-application-security.md"}, AIReview: review(digest, "Ownership is unclear.", "ownership-unclear", "critical")},
			{ControlID: "PRC-01-004", Source: model.Source{Path: "docs/engineering/08-security-and-cryptography.md"}, AIReview: review(digest, "The runtime proof is missing.", "ownership-unclear", "low")},
		},
	}
	sourceRunID, err := sourceIdentity(run)
	if err != nil {
		t.Fatal(err)
	}
	plan, err := Build(run, sourceRunID)
	if err != nil {
		t.Fatal(err)
	}
	if plan.SchemaVersion != SchemaVersion || plan.Authority != "advisory_only" ||
		plan.ReviewedControlCount != 4 || plan.ItemCount != 3 {
		t.Fatalf("unexpected plan: %+v", plan)
	}
	grouped := plan.Items[1]
	if grouped.Priority != "high" || grouped.ControlCount != 2 ||
		len(grouped.ControlIDs) != 2 || grouped.ItemID == "" {
		t.Fatalf("exact cause was not grouped conservatively: %+v", grouped)
	}
	second, err := Build(run, sourceRunID)
	if err != nil || second.Items[1].ItemID != grouped.ItemID {
		t.Fatalf("plan identity was not stable: %+v err=%v", second, err)
	}
}

func TestBuildRejectsUnboundSourceRun(t *testing.T) {
	run := model.RunResult{ControlCatalog: &model.ControlCatalogSummary{AIReviewProvider: "codex"}}
	if _, err := Build(run, "not-a-digest"); err == nil {
		t.Fatal("unbound improvement plan source was accepted")
	}
}
