package aiplan

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
)

const SchemaVersion = "prc.ai-improvement-plan/v0.1"

type improvementGroup struct {
	item                 model.AIImprovementPlanItem
	assessmentCandidates map[string]bool
	taskIDs              map[string]bool
}

// Build performs only stable, exact grouping. It never asks a
// model to merge semantically similar text and never changes a control result.
// The source reviews remain the complete advice; this compact index only shows
// which reviews share an exact domain, cause key, and normalized cause.
func Build(run model.RunResult, sourceRunID string) (*model.AIImprovementPlan, error) {
	if run.ControlCatalog == nil || run.ControlCatalog.AIReviewProvider == "" {
		return nil, nil
	}
	if !lowerHexDigest(sourceRunID) {
		return nil, fmt.Errorf("AI improvement plan requires a content-addressed source run")
	}
	expectedSourceRunID, err := sourceIdentity(run)
	if err != nil {
		return nil, err
	}
	if sourceRunID != expectedSourceRunID {
		return nil, fmt.Errorf("AI improvement plan source does not match the pre-review run")
	}
	groups := map[string]*improvementGroup{}
	reviewedControls := 0
	for _, control := range run.ControlResults {
		review := control.AIReview
		if review == nil {
			continue
		}
		reviewedControls++
		domain := reviewDomain(control.Source.Path)
		normalizedCause := strings.ToLower(strings.Join(strings.Fields(review.RootCause), " "))
		groupKey := domain + "\x00" + review.RootCauseKey + "\x00" + normalizedCause
		group := groups[groupKey]
		if group == nil {
			group = &improvementGroup{
				item: model.AIImprovementPlanItem{
					Domain: domain, RootCauseKey: review.RootCauseKey,
					RootCause: strings.TrimSpace(review.RootCause), Priority: review.Priority,
					Effort: review.Effort, BlastRadius: review.BlastRadius,
					ControlIDs: []string{}, TaskIDs: []string{}, AssessmentCandidates: []string{},
				},
				assessmentCandidates: map[string]bool{}, taskIDs: map[string]bool{},
			}
			groups[groupKey] = group
		}
		if strings.TrimSpace(review.RootCause) < group.item.RootCause {
			group.item.RootCause = strings.TrimSpace(review.RootCause)
		}
		if priorityRank(review.Priority) < priorityRank(group.item.Priority) {
			group.item.Priority = review.Priority
		}
		group.item.Effort = conservativeEffort(group.item.Effort, review.Effort)
		group.item.BlastRadius = conservativeBlastRadius(group.item.BlastRadius, review.BlastRadius)
		group.item.ControlIDs = append(group.item.ControlIDs, control.ControlID)
		group.assessmentCandidates[review.AssessmentCandidate] = true
		group.taskIDs[review.TaskID] = true
	}

	items := make([]model.AIImprovementPlanItem, 0, len(groups))
	for _, group := range groups {
		sort.Strings(group.item.ControlIDs)
		group.item.ControlCount = len(group.item.ControlIDs)
		group.item.AssessmentCandidates = sortedKeys(group.assessmentCandidates)
		group.item.TaskIDs = sortedKeys(group.taskIDs)
		itemID, err := improvementItemID(group.item)
		if err != nil {
			return nil, err
		}
		group.item.ItemID = itemID
		items = append(items, group.item)
	}
	sort.Slice(items, func(left, right int) bool {
		if leftRank, rightRank := priorityRank(items[left].Priority), priorityRank(items[right].Priority); leftRank != rightRank {
			return leftRank < rightRank
		}
		if items[left].ControlCount != items[right].ControlCount {
			return items[left].ControlCount > items[right].ControlCount
		}
		if items[left].Domain != items[right].Domain {
			return items[left].Domain < items[right].Domain
		}
		if items[left].RootCauseKey != items[right].RootCauseKey {
			return items[left].RootCauseKey < items[right].RootCauseKey
		}
		return items[left].RootCause < items[right].RootCause
	})
	return &model.AIImprovementPlan{
		SchemaVersion: SchemaVersion, Authority: "advisory_only", SourceRunID: sourceRunID,
		ReviewProvider: run.ControlCatalog.AIReviewProvider, ReviewModel: run.ControlCatalog.AIReviewModel,
		ReviewDepth: run.ControlCatalog.AIReviewDepth, ReviewState: run.ControlCatalog.AIReviewState,
		ReviewedControlCount: reviewedControls, ItemCount: len(items), Items: items,
	}, nil
}

// sourceIdentity reconstructs the exact pre-review run. This prevents a plan
// from claiming an unrelated source while keeping the original scanner result
// independently content-addressed.
func sourceIdentity(run model.RunResult) (string, error) {
	copyRun := run
	copyRun.RunID = ""
	copyRun.AIImprovementPlan = nil
	if run.ControlCatalog != nil {
		catalog := *run.ControlCatalog
		catalog.AIReviewProvider = ""
		catalog.AIReviewModel = ""
		catalog.AIReviewDepth = ""
		catalog.AIReviewState = ""
		catalog.AIReviewedCount = 0
		catalog.AIAdvisoryFailCount = 0
		copyRun.ControlCatalog = &catalog
	}
	copyRun.ControlResults = append([]model.ControlResult(nil), run.ControlResults...)
	for index := range copyRun.ControlResults {
		copyRun.ControlResults[index].AIReview = nil
	}
	payload, err := json.Marshal(copyRun)
	if err != nil {
		return "", fmt.Errorf("encode pre-review run identity: %w", err)
	}
	digest := sha256.Sum256(payload)
	return hex.EncodeToString(digest[:]), nil
}

func lowerHexDigest(value string) bool {
	if len(value) != 64 {
		return false
	}
	for _, character := range value {
		if (character < '0' || character > '9') && (character < 'a' || character > 'f') {
			return false
		}
	}
	return true
}

func improvementItemID(item model.AIImprovementPlanItem) (string, error) {
	item.ItemID = ""
	payload, err := json.Marshal(item)
	if err != nil {
		return "", fmt.Errorf("encode AI improvement-plan item identity: %w", err)
	}
	digest := sha256.Sum256(payload)
	return hex.EncodeToString(digest[:]), nil
}

func reviewDomain(sourcePath string) string {
	normalized := filepath.ToSlash(sourcePath)
	scope := "catalog"
	if strings.Contains(normalized, "/engineering/") {
		scope = "engineering"
	} else if strings.Contains(normalized, "/checklists/") {
		scope = "checklist"
	}
	base := strings.TrimSuffix(filepath.Base(normalized), filepath.Ext(normalized))
	if separator := strings.IndexByte(base, '-'); separator > 0 {
		allDigits := true
		for _, character := range base[:separator] {
			if character < '0' || character > '9' {
				allDigits = false
				break
			}
		}
		if allDigits {
			base = base[separator+1:]
		}
	}
	var domain strings.Builder
	lastSeparator := false
	for _, character := range strings.ToLower(base) {
		switch {
		case character >= 'a' && character <= 'z', character >= '0' && character <= '9':
			domain.WriteRune(character)
			lastSeparator = false
		case domain.Len() > 0 && !lastSeparator:
			domain.WriteByte('-')
			lastSeparator = true
		}
	}
	slug := strings.Trim(domain.String(), "-")
	if slug == "" {
		slug = "uncategorized"
	}
	return scope + "/" + slug
}

func priorityRank(value string) int {
	switch value {
	case "critical":
		return 0
	case "high":
		return 1
	case "medium":
		return 2
	case "low":
		return 3
	default:
		return 4
	}
}

func conservativeEffort(left, right string) string {
	rank := map[string]int{"small": 0, "medium": 1, "large": 2, "unknown": 3}
	if rank[right] > rank[left] {
		return right
	}
	return left
}

func conservativeBlastRadius(left, right string) string {
	rank := map[string]int{"local": 0, "component": 1, "system": 2, "organization": 3, "unknown": 4}
	if rank[right] > rank[left] {
		return right
	}
	return left
}

func sortedKeys(values map[string]bool) []string {
	result := make([]string, 0, len(values))
	for value := range values {
		result = append(result, value)
	}
	sort.Strings(result)
	return result
}
