package controlreview

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/MarinJursic/production-readiness-checklist/scanner/aiplan"
	"github.com/MarinJursic/production-readiness-checklist/scanner/fullscan"
	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
)

const (
	maximumContextFileBytes = 128 * 1024
	maximumContextTotal     = 4 * 1024 * 1024
	maximumPathContextBytes = 1024 * 1024
	maximumCachedOutput     = 1024 * 1024
)

var reviewedDeterministicRoutes = map[string]bool{
	"local_static": true, "artifact_verification": true, "bounded_execution": true,
	"external_readonly_query": true, "structured_record_validation": true,
	"deterministic_composite": true,
}

var reviewedNondeterministicRoutes = map[string]bool{
	"contextual_judgment": true, "accountable_human_decision": true,
	"specialist_or_legal_judgment": true, "empirical_protocol_undefined": true,
	"contract_incomplete": true, "mixed": true, "unbounded_claim": true,
}

// Apply runs an explicitly authorized advisory review for every selected,
// active, reviewed nondeterministic control. Repository content is copied,
// screened, and embedded as untrusted context; provider tools never receive
// the source workspace path.
func Apply(ctx context.Context, run model.RunResult, options Options) (model.RunResult, Summary, error) {
	started := time.Now()
	if run.ControlCatalog == nil || len(run.ControlResults) == 0 {
		return model.RunResult{}, Summary{}, fmt.Errorf("AI review requires a complete control-catalog scan")
	}
	if !options.AllowRemoteSourceProcessing {
		return model.RunResult{}, Summary{}, fmt.Errorf("AI review requires explicit --allow-remote-source-processing acknowledgement")
	}
	if err := normalizeOptions(&options); err != nil {
		return model.RunResult{}, Summary{}, err
	}
	// Reject legacy, unknown, and explicitly deterministic selections before
	// copying any source context or creating resumable state.
	if _, _, err := selectReviewControls(run, options.ControlIDs); err != nil {
		return model.RunResult{}, Summary{}, err
	}
	snapshot, err := createSnapshot(run.Inventory)
	if err != nil {
		return model.RunResult{}, Summary{}, fmt.Errorf("prepare safe remote-review snapshot: %w", err)
	}
	defer os.RemoveAll(snapshot.Directory)
	stateDirectory, err := prepareStateDirectory(run, options)
	if err != nil {
		return model.RunResult{}, Summary{}, err
	}
	options.StateDirectory = stateDirectory
	tasks, focused, err := buildTasks(run, snapshot, options)
	if err != nil {
		return model.RunResult{}, Summary{}, err
	}
	runner := options.Runner
	if runner == nil {
		runner, err = newCLIRunner(options)
		if err != nil {
			return model.RunResult{}, Summary{}, err
		}
	}
	outputs := make([]Output, len(tasks))
	reused := 0
	totalControls := 0
	reusedControls := 0
	pending := make([]int, 0, len(tasks))
	for index, task := range tasks {
		totalControls += len(task.Controls)
		cached, ok, cacheErr := loadCachedOutput(stateDirectory, task)
		if cacheErr != nil {
			return model.RunResult{}, Summary{}, cacheErr
		}
		if ok {
			outputs[index] = cached
			reused++
			reusedControls += len(task.Controls)
		} else {
			pending = append(pending, index)
		}
	}
	baseProgress := Progress{
		Phase: "prepared", Provider: options.Provider, Model: options.Model,
		ReviewDepth:    options.ReviewDepth,
		StateDirectory: stateDirectory, Workers: options.Workers,
		TotalBatches: len(tasks), CompletedBatches: reused, ReusedBatches: reused,
		TotalControls: totalControls, CompletedControls: reusedControls,
		MaxCostUSD: options.MaxCostUSD, Elapsed: time.Since(started),
	}
	if options.Progress != nil {
		options.Progress(baseProgress)
	}
	statistics, batchErr := runPendingBatches(
		ctx, runner, tasks, pending, outputs, stateDirectory, options.Workers,
		options.Progress, baseProgress, started,
	)
	reviews := map[string]Review{}
	for _, output := range outputs {
		if output.SchemaVersion == "" {
			continue
		}
		for _, review := range output.Reviews {
			if _, exists := reviews[review.ControlID]; exists {
				return model.RunResult{}, Summary{}, fmt.Errorf("provider repeated advisory review for %s", review.ControlID)
			}
			reviews[review.ControlID] = review
		}
	}
	advisoryFailures := 0
	for index := range run.ControlResults {
		review, ok := reviews[run.ControlResults[index].ControlID]
		if !ok {
			continue
		}
		if review.AssessmentCandidate == "advisory_fail_candidate" {
			advisoryFailures++
		}
		run.ControlResults[index].AIReview = &model.AIControlReview{
			Provider: options.Provider, Model: options.Model,
			ReviewDepth:            options.ReviewDepth,
			AssessmentCandidate:    review.AssessmentCandidate,
			ApplicabilityCandidate: review.ApplicabilityCandidate,
			Confidence:             review.Confidence, Priority: review.Priority,
			RootCause: review.RootCause, RootCauseKey: review.RootCauseKey,
			Effort: review.Effort, BlastRadius: review.BlastRadius,
			Reason: review.Reason, Challenge: review.Challenge,
			RiskIfIgnored: review.RiskIfIgnored, Advice: review.Advice,
			RemediationSteps:     append([]string{}, review.RemediationSteps...),
			VerificationSteps:    append([]string{}, review.VerificationSteps...),
			EvidenceNeeded:       append([]string{}, review.EvidenceNeeded...),
			Evidence:             append([]model.FindingLocation{}, review.Evidence...),
			Limitations:          append([]string{}, review.Limitations...),
			CitationVerification: citationVerification(review.Evidence),
			ClaimVerification:    "advisory_unverified",
			TaskID:               taskForControl(tasks, review.ControlID),
		}
	}
	run.ControlCatalog.AIReviewProvider = options.Provider
	run.ControlCatalog.AIReviewModel = options.Model
	run.ControlCatalog.AIReviewDepth = options.ReviewDepth
	run.ControlCatalog.AIReviewedCount = len(reviews)
	run.ControlCatalog.AIAdvisoryFailCount = advisoryFailures
	if batchErr != nil {
		run.ControlCatalog.AIReviewState = "partial"
	} else if focused {
		run.ControlCatalog.AIReviewState = "focused"
	} else {
		run.ControlCatalog.AIReviewState = "complete"
	}
	run.AIImprovementPlan, err = aiplan.Build(run, run.RunID)
	if err != nil {
		return model.RunResult{}, Summary{}, err
	}
	run, err = fullscan.Reidentify(run)
	if err != nil {
		return model.RunResult{}, Summary{}, err
	}
	summary := Summary{
		Provider: options.Provider, Model: options.Model, ReviewDepth: options.ReviewDepth,
		ReviewedControls: len(reviews),
		AdvisoryFailures: advisoryFailures, ReusedBatches: reused,
		CompletedBatches: reused + statistics.CompletedBatches, StateDirectory: stateDirectory, Focused: focused,
		TokenUsage: statistics.TokenUsage, TokenUsageBatches: statistics.TokenUsageBatches,
		EstimatedCostUSD: statistics.EstimatedCostUSD, EstimatedCostBatches: statistics.EstimatedCostBatches,
		MaxCostUSD: options.MaxCostUSD, Duration: time.Since(started),
	}
	if batchErr != nil {
		return run, summary, batchErr
	}
	return run, summary, nil
}

func citationVerification(evidence []model.FindingLocation) string {
	if len(evidence) == 0 {
		return "not_cited"
	}
	return "snapshot_location_validated"
}

func normalizeOptions(options *Options) error {
	if options.Provider != "codex" && options.Provider != "claude" {
		return fmt.Errorf("unsupported AI review provider %q", options.Provider)
	}
	if options.Executable == "" {
		options.Executable = options.Provider
	}
	if options.ReasoningEffort == "" {
		options.ReasoningEffort = "high"
	}
	if options.ReasoningEffort != "high" && options.ReasoningEffort != "xhigh" {
		return fmt.Errorf("AI review reasoning effort must be high or xhigh")
	}
	if options.Provider == "claude" && options.ReasoningEffort != "high" {
		return fmt.Errorf("claude AI review supports high reasoning effort; xhigh is Codex-only")
	}
	if options.ReviewDepth == "" {
		options.ReviewDepth = "standard"
	}
	if options.ReviewDepth != "standard" && options.ReviewDepth != "deep" {
		return fmt.Errorf("AI review depth must be standard or deep")
	}
	if options.BatchSize == 0 {
		options.BatchSize = 8
	}
	if options.BatchSize < 1 || options.BatchSize > 8 {
		return fmt.Errorf("AI review batch size must be between 1 and 8")
	}
	if options.Workers == 0 {
		options.Workers = 1
	}
	if options.Workers < 1 || options.Workers > 4 {
		return fmt.Errorf("AI review workers must be between 1 and 4")
	}
	if options.Timeout == 0 {
		options.Timeout = 30 * time.Minute
	}
	if options.Timeout < 30*time.Second || options.Timeout > 2*time.Hour {
		return fmt.Errorf("AI review timeout must be between 30 seconds and 2 hours per batch")
	}
	if options.MaxCostUSD < 0 || options.Provider == "codex" && options.MaxCostUSD > 0 {
		return fmt.Errorf("only Claude can enforce a nonzero per-batch AI review cost limit")
	}
	return nil
}

func prepareStateDirectory(run model.RunResult, options Options) (string, error) {
	path := options.StateDirectory
	if path == "" {
		cache, err := os.UserCacheDir()
		if err != nil {
			return "", fmt.Errorf("locate user cache for resumable AI review: %w", err)
		}
		path = filepath.Join(cache, "prc", "control-reviews", run.Inventory.Digest, run.ControlCatalog.RegistrySHA256, options.Provider)
	}
	absolute, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("resolve AI review state directory: %w", err)
	}
	if pathInside(run.Inventory.Root, absolute) {
		return "", fmt.Errorf("AI review state directory must be outside the scanned project")
	}
	if err := os.MkdirAll(absolute, 0o700); err != nil {
		return "", fmt.Errorf("create AI review state directory: %w", err)
	}
	info, err := os.Lstat(absolute)
	if err != nil || !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
		return "", fmt.Errorf("AI review state path is not a regular directory")
	}
	if err := os.Chmod(absolute, 0o700); err != nil {
		return "", fmt.Errorf("protect AI review state directory: %w", err)
	}
	return absolute, nil
}

func pathInside(root, path string) bool {
	relative, err := filepath.Rel(filepath.Clean(root), filepath.Clean(path))
	return err == nil && relative != ".." && !strings.HasPrefix(relative, ".."+string(filepath.Separator))
}

func buildTasks(run model.RunResult, snapshot reviewSnapshot, options Options) ([]Task, bool, error) {
	if options.ReviewDepth == "" {
		options.ReviewDepth = "standard"
	}
	controls, focusedRequested, err := selectReviewControls(run, options.ControlIDs)
	if err != nil {
		return nil, false, err
	}
	assertions := map[string]model.AssertionResult{}
	for _, result := range run.Results {
		assertions[result.AssertionID] = result
	}
	tasks := make([]Task, 0, (len(controls)+options.BatchSize-1)/options.BatchSize)
	for start := 0; start < len(controls); start += options.BatchSize {
		end := min(start+options.BatchSize, len(controls))
		task := Task{
			SchemaVersion: TaskSchema, InventoryDigest: run.Inventory.Digest,
			RegistrySHA256: run.ControlCatalog.RegistrySHA256, Provider: options.Provider,
			ReviewDepth: options.ReviewDepth, RequireOneSubagentPerRule: true,
			RequireBatchSkeptic: options.ReviewDepth == "deep", Controls: []TaskControl{},
		}
		for _, control := range controls[start:end] {
			item := TaskControl{
				ControlID: control.ControlID, Statement: control.Statement,
				ChecklistSource: control.Source, CurrentDisposition: control.Disposition,
				ContractSHA256: control.ContractSHA256, ContractStatus: control.ContractStatus,
				Classification: control.Classification, ClassificationRoute: control.ClassificationRoute,
				ClassificationDecisionBasis: control.ClassificationDecisionBasis,
				ClassificationRowSHA256:     control.ClassificationRowSHA256,
				CanonicalControlID:          control.CanonicalControlID, EvaluationClass: control.EvaluationClass,
				AutomationClass: control.AutomationClass, ApplicabilityClass: control.ApplicabilityClass,
				Atomicity: control.Atomicity, CompleteInventory: control.CompleteInventoryRequired,
				NegativeCondition: control.NegativeCondition, ProjectThresholds: control.ProjectThresholdsRequired,
				EvidenceAuthorities: append([]string{}, control.EvidenceAuthorities...),
				NotApplicableProof:  control.NotApplicableProof,
				CurrentCoverage:     control.Coverage, CurrentAssertionChecks: []AssertionContext{},
			}
			for _, assertionID := range control.ExecutedAssertionIDs {
				result := assertions[assertionID]
				item.CurrentAssertionChecks = append(item.CurrentAssertionChecks, AssertionContext{
					AssertionID: assertionID, Assessment: result.Assessment,
					Summary: result.Summary, Locations: append([]model.FindingLocation{}, result.Locations...),
				})
			}
			task.Controls = append(task.Controls, item)
		}
		task.RepositoryPaths, task.SnapshotLimitations = boundedPaths(snapshot.Paths, snapshot.Limitations)
		task.ContextFiles, task.SnapshotLimitations = selectContext(task.Controls, snapshot, task.SnapshotLimitations)
		sealed, err := sealTask(task)
		if err != nil {
			return nil, false, err
		}
		tasks = append(tasks, sealed)
	}
	return tasks, focusedRequested, nil
}

func selectReviewControls(run model.RunResult, requested []string) ([]model.ControlResult, bool, error) {
	selected := map[string]bool{}
	for _, id := range requested {
		if selected[id] {
			return nil, false, fmt.Errorf("duplicate --review-control %s", id)
		}
		selected[id] = true
	}
	focusedRequested := len(selected) > 0
	controls := make([]model.ControlResult, 0, len(run.ControlResults))
	seen := make(map[string]bool, len(run.ControlResults))
	for _, control := range run.ControlResults {
		if control.ControlID == "" || seen[control.ControlID] {
			return nil, false, fmt.Errorf("AI review control results contain a missing or duplicate control ID %q", control.ControlID)
		}
		seen[control.ControlID] = true
		classification, err := reviewedClassification(control)
		if err != nil {
			return nil, false, err
		}
		requestedControl := selected[control.ControlID]
		if requestedControl && classification == "deterministic" {
			return nil, false, fmt.Errorf("--review-control %s is reviewed deterministic and cannot receive an AI advisory verdict", control.ControlID)
		}
		if control.Disposition == "retired" || focusedRequested && !requestedControl {
			continue
		}
		if requestedControl {
			delete(selected, control.ControlID)
		}
		if classification == "nondeterministic" {
			controls = append(controls, control)
		}
	}
	if len(selected) > 0 {
		unknown := make([]string, 0, len(selected))
		for id := range selected {
			unknown = append(unknown, id)
		}
		sort.Strings(unknown)
		return nil, false, fmt.Errorf("unknown or retired --review-control values: %s", strings.Join(unknown, ", "))
	}
	return controls, focusedRequested, nil
}

func reviewedClassification(control model.ControlResult) (string, error) {
	if control.ContractStatus != "reviewed" {
		return "", fmt.Errorf("AI review requires a reviewed control contract for %s; contract status %q is not accepted", control.ControlID, control.ContractStatus)
	}
	if !lowerHexDigest(control.ClassificationRowSHA256) {
		return "", fmt.Errorf("AI review requires reviewed classification data for %s; classification row digest is missing or invalid", control.ControlID)
	}
	switch control.Classification {
	case "deterministic":
		if !reviewedDeterministicRoutes[control.ClassificationRoute] || control.ClassificationDecisionBasis != "strength_audit_confirmed" {
			return "", fmt.Errorf("AI review requires reviewed classification data for %s; deterministic route or decision basis is invalid", control.ControlID)
		}
	case "nondeterministic":
		if !reviewedNondeterministicRoutes[control.ClassificationRoute] ||
			(control.ClassificationDecisionBasis != "primary_nondeterministic" &&
				control.ClassificationDecisionBasis != "skeptically_rejected" &&
				control.ClassificationDecisionBasis != "strength_audit_reclassified") {
			return "", fmt.Errorf("AI review requires reviewed classification data for %s; nondeterministic route or decision basis is invalid", control.ControlID)
		}
	default:
		return "", fmt.Errorf("AI review requires reviewed classification data for %s; classification %q is missing or unknown", control.ControlID, control.Classification)
	}
	return control.Classification, nil
}

func boundedPaths(paths []string, limitations []string) ([]string, []string) {
	result := make([]string, 0, len(paths))
	used := 0
	for _, path := range paths {
		if used+len(path)+1 > maximumPathContextBytes {
			return result, uniqueSorted(append(limitations, "Repository path inventory was truncated at its scanner-owned byte limit."))
		}
		result = append(result, path)
		used += len(path) + 1
	}
	return result, uniqueSorted(limitations)
}

type scoredPath struct {
	path       string
	score      int
	anchorLine int
}

func selectContext(controls []TaskControl, snapshot reviewSnapshot, limitations []string) ([]ContextFile, []string) {
	tokens := map[string]bool{}
	for _, control := range controls {
		for _, token := range strings.FieldsFunc(strings.ToLower(control.Statement), func(value rune) bool { return !unicode.IsLetter(value) && !unicode.IsDigit(value) }) {
			if len(token) >= 4 {
				tokens[token] = true
			}
		}
	}
	scores := make([]scoredPath, 0, len(snapshot.Paths))
	for _, path := range snapshot.Paths {
		lowerPath := strings.ToLower(path)
		lowerContent := strings.ToLower(string(snapshot.Contents[path]))
		score := 0
		anchorLine := 1
		firstMatch := -1
		for token := range tokens {
			if strings.Contains(lowerPath, token) {
				score += 20
			}
			if offset := strings.Index(lowerContent, token); offset >= 0 {
				score += 2
				if firstMatch < 0 || offset < firstMatch {
					firstMatch = offset
				}
			}
		}
		base := strings.ToLower(filepath.Base(path))
		if strings.HasPrefix(base, "readme") || base == "security.md" || base == "contributing.md" ||
			base == "package.json" || base == "go.mod" || base == "pyproject.toml" || base == "cargo.toml" {
			score += 10
		}
		if score > 0 {
			if firstMatch >= 0 {
				anchorLine = bytes.Count(snapshot.Contents[path][:firstMatch], []byte{'\n'}) + 1
			}
			scores = append(scores, scoredPath{path: path, score: score, anchorLine: anchorLine})
		}
	}
	sort.Slice(scores, func(left, right int) bool {
		if scores[left].score != scores[right].score {
			return scores[left].score > scores[right].score
		}
		return scores[left].path < scores[right].path
	})
	files := make([]ContextFile, 0)
	used := 0
	for _, scored := range scores {
		if len(snapshot.Contents[scored.path]) == 0 {
			limitations = append(limitations, "At least one relevant empty text file could not provide a citable excerpt.")
			continue
		}
		data, startLine, endLine, truncated := contextExcerpt(snapshot.Contents[scored.path], scored.anchorLine)
		if used+len(data) > maximumContextTotal {
			limitations = append(limitations, "Relevant repository excerpts were capped at the scanner-owned total byte limit.")
			break
		}
		digest := sha256.Sum256(data)
		files = append(files, ContextFile{
			Path: scored.path, SHA256: hex.EncodeToString(digest[:]),
			StartLine: startLine, EndLine: endLine, Content: string(data),
		})
		used += len(data)
		if truncated {
			limitations = append(limitations, "At least one relevant file excerpt was truncated at 128 KiB.")
		}
		if len(files) == 64 {
			limitations = append(limitations, "Relevant repository excerpts were capped at 64 files.")
			break
		}
	}
	limitations = uniqueSorted(limitations)
	return files, limitations
}

func contextExcerpt(data []byte, anchorLine int) ([]byte, int, int, bool) {
	lines := bytes.SplitAfter(data, []byte{'\n'})
	if len(lines) > 1 && len(lines[len(lines)-1]) == 0 {
		lines = lines[:len(lines)-1]
	}
	if len(lines) == 0 {
		return []byte{}, 1, 1, false
	}
	anchor := max(1, min(anchorLine, len(lines))) - 1
	if len(data) <= maximumContextFileBytes {
		return append([]byte{}, data...), 1, len(lines), false
	}
	if len(lines[anchor]) > maximumContextFileBytes {
		return utf8Prefix(lines[anchor], maximumContextFileBytes), anchor + 1, anchor + 1, true
	}
	start, end := anchor, anchor+1
	used := len(lines[anchor])
	left, right := anchor-1, anchor+1
	takeLeft := true
	for left >= 0 || right < len(lines) {
		if takeLeft && left >= 0 || right >= len(lines) {
			if used+len(lines[left]) <= maximumContextFileBytes {
				used += len(lines[left])
				start = left
				left--
			} else {
				left = -1
			}
		} else if right < len(lines) {
			if used+len(lines[right]) <= maximumContextFileBytes {
				used += len(lines[right])
				end = right + 1
				right++
			} else {
				right = len(lines)
			}
		}
		takeLeft = !takeLeft
	}
	return bytes.Join(lines[start:end], nil), start + 1, end, true
}

func utf8Prefix(data []byte, limit int) []byte {
	if len(data) <= limit {
		return append([]byte{}, data...)
	}
	end := limit
	for end > 0 && !utf8.Valid(data[:end]) {
		end--
	}
	return append([]byte{}, data[:end]...)
}

func uniqueSorted(values []string) []string {
	seen := map[string]bool{}
	result := make([]string, 0, len(values))
	for _, value := range values {
		if value != "" && !seen[value] {
			seen[value] = true
			result = append(result, value)
		}
	}
	sort.Strings(result)
	return result
}

type batchRunStatistics struct {
	CompletedBatches     int
	CompletedControls    int
	TokenUsage           TokenUsage
	TokenUsageBatches    int
	EstimatedCostUSD     float64
	EstimatedCostBatches int
}

func runPendingBatches(
	ctx context.Context, runner BatchRunner, tasks []Task, pending []int, outputs []Output,
	stateDirectory string, workers int, progress func(Progress), base Progress, started time.Time,
) (batchRunStatistics, error) {
	workContext, cancel := context.WithCancel(ctx)
	defer cancel()
	jobs := make(chan int)
	var wait sync.WaitGroup
	var lock sync.Mutex
	var firstErr error
	statistics := batchRunStatistics{}
	worker := func() {
		defer wait.Done()
		for index := range jobs {
			if workContext.Err() != nil {
				return
			}
			output, execution, err := runner.Run(workContext, tasks[index])
			if err == nil {
				output = attachTaskLimitations(output, tasks[index])
				err = validateOutput(output, tasks[index])
			}
			if err == nil {
				err = storeCachedOutput(stateDirectory, output)
			}
			lock.Lock()
			if err != nil && firstErr == nil {
				firstErr = fmt.Errorf("AI review batch %s failed: %w", tasks[index].TaskID, err)
				cancel()
			} else if err == nil {
				outputs[index] = output
				statistics.CompletedBatches++
				statistics.CompletedControls += len(tasks[index].Controls)
				if execution.TokenUsageKnown {
					if usageErr := addTokenUsage(&statistics.TokenUsage, execution.TokenUsage); usageErr != nil {
						firstErr = usageErr
						cancel()
					} else {
						statistics.TokenUsageBatches++
					}
				}
				if firstErr == nil && execution.EstimatedCostKnown {
					if costErr := addEstimatedCost(&statistics.EstimatedCostUSD, execution.EstimatedCostUSD); costErr != nil {
						firstErr = costErr
						cancel()
					} else {
						statistics.EstimatedCostBatches++
					}
				}
				if firstErr == nil && progress != nil {
					current := base
					current.Phase = "batch_completed"
					current.CompletedBatches += statistics.CompletedBatches
					current.CompletedControls += statistics.CompletedControls
					current.TokenUsage = statistics.TokenUsage
					current.TokenUsageBatches = statistics.TokenUsageBatches
					current.EstimatedCostUSD = statistics.EstimatedCostUSD
					current.EstimatedCostBatches = statistics.EstimatedCostBatches
					current.Elapsed = time.Since(started)
					progress(current)
				}
			}
			lock.Unlock()
		}
	}
	for range min(workers, len(pending)) {
		wait.Add(1)
		go worker()
	}
enqueue:
	for _, index := range pending {
		select {
		case jobs <- index:
		case <-workContext.Done():
			break enqueue
		}
		if workContext.Err() != nil {
			break
		}
	}
	close(jobs)
	wait.Wait()
	if firstErr != nil {
		return statistics, firstErr
	}
	if statistics.CompletedBatches == len(pending) {
		return statistics, nil
	}
	return statistics, ctx.Err()
}

func loadCachedOutput(directory string, task Task) (Output, bool, error) {
	path := filepath.Join(directory, task.TaskID+".json")
	info, err := os.Lstat(path)
	if os.IsNotExist(err) {
		return Output{}, false, nil
	}
	if err != nil || !info.Mode().IsRegular() || info.Mode()&os.ModeSymlink != 0 || info.Size() > maximumCachedOutput {
		return Output{}, false, fmt.Errorf("cached AI review %s is not a bounded regular file", task.TaskID)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return Output{}, false, fmt.Errorf("read cached AI review: %w", err)
	}
	var output Output
	if err := decodeInnerOutput(data, &output); err != nil || validateOutput(output, task) != nil {
		return Output{}, false, fmt.Errorf("cached AI review %s failed its sealed protocol", task.TaskID)
	}
	return output, true, nil
}

func attachTaskLimitations(output Output, task Task) Output {
	for index := range output.Reviews {
		output.Reviews[index].Limitations = uniqueSorted(append(
			append([]string{}, output.Reviews[index].Limitations...), task.SnapshotLimitations...,
		))
	}
	return output
}

func storeCachedOutput(directory string, output Output) error {
	payload, err := json.Marshal(output)
	if err != nil {
		return fmt.Errorf("encode cached AI review: %w", err)
	}
	path := filepath.Join(directory, output.TaskID+".json")
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return nil
		}
		return fmt.Errorf("create cached AI review: %w", err)
	}
	if _, err := file.Write(payload); err != nil {
		_ = file.Close()
		_ = os.Remove(path)
		return fmt.Errorf("write cached AI review: %w", err)
	}
	if err := file.Close(); err != nil {
		_ = os.Remove(path)
		return fmt.Errorf("finish cached AI review: %w", err)
	}
	return nil
}

func taskForControl(tasks []Task, controlID string) string {
	for _, task := range tasks {
		for _, control := range task.Controls {
			if control.ControlID == controlID {
				return task.TaskID
			}
		}
	}
	return ""
}
