package engine

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/MarinJursic/production-readiness-checklist/scanner/catalog"
	workspaceinventory "github.com/MarinJursic/production-readiness-checklist/scanner/inventory"
	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
	"gopkg.in/yaml.v3"
)

const applicabilityEvaluator = "prc-expr/v0.1"

var (
	actionUsePattern = regexp.MustCompile(`(?m)^\s*(?:-\s*)?uses:\s*["']?([^"'\s#]+)`)
	commitPattern    = regexp.MustCompile(`^[0-9a-fA-F]{40}$`)
)

type Engine struct {
	Catalog *catalog.Catalog
	Now     func() time.Time
}

func New(c *catalog.Catalog) *Engine {
	return &Engine{Catalog: c, Now: func() time.Time { return time.Now().UTC() }}
}

func (e *Engine) Plan(profileID string, inventory model.Inventory) (model.Plan, error) {
	profile, err := e.Catalog.Profile(profileID)
	if err != nil {
		return model.Plan{}, err
	}
	plan := model.Plan{
		SchemaVersion: model.PlanSchema, TargetName: inventory.TargetName,
		TargetCommit: inventory.GitCommit, InventoryDigest: inventory.Digest,
		ProfileID: profile.ID, ProfileVersion: profile.Version,
	}
	for _, assertionID := range profile.AssertionIDs {
		assertion := e.Catalog.Assertions[assertionID]
		applicability, _ := evaluateApplicability(assertion.Applicability, inventory)
		plan.Assertions = append(plan.Assertions, model.PlannedAssertion{
			AssertionID: assertion.ID, Implementation: assertion.ImplementationID,
			Applicability: applicability, ApplicabilityBy: applicabilityEvaluator,
		})
	}
	payload, err := json.Marshal(plan)
	if err != nil {
		return model.Plan{}, fmt.Errorf("encode plan: %w", err)
	}
	digest := sha256.Sum256(payload)
	plan.Digest = hex.EncodeToString(digest[:])
	return plan, nil
}

func (e *Engine) Scan(profileID string, inventory model.Inventory) (model.RunResult, error) {
	return e.ScanWithAdapterEvidence(profileID, inventory, nil)
}

func (e *Engine) ScanWithAdapterEvidence(
	profileID string,
	inventory model.Inventory,
	executions []model.AdapterExecution,
) (model.RunResult, error) {
	profile, err := e.Catalog.Profile(profileID)
	if err != nil {
		return model.RunResult{}, err
	}
	plan, err := e.Plan(profileID, inventory)
	if err != nil {
		return model.RunResult{}, err
	}
	started := e.Now()
	validatedExecutions, err := validateAdapterExecutions(inventory, executions)
	if err != nil {
		return model.RunResult{}, err
	}
	for _, execution := range validatedExecutions {
		authorized, err := e.authorizesAdapterInPlan(plan, execution.AdapterID, execution.ManifestSHA256)
		if err != nil {
			return model.RunResult{}, err
		}
		if !authorized {
			return model.RunResult{}, fmt.Errorf("adapter execution %s is not authorized by an applicable assertion", execution.ExecutionID)
		}
	}
	run := model.RunResult{
		SchemaVersion: model.RunSchema, StartedAt: started, Plan: plan, Inventory: inventory,
		AdapterExecutions: validatedExecutions,
	}
	for _, planned := range plan.Assertions {
		assertion := e.Catalog.Assertions[planned.AssertionID]
		result := model.AssertionResult{
			AssertionID: assertion.ID, ControlIDs: assertion.ControlIDs,
			Applicability: planned.Applicability, Severity: assertion.Severity,
			Gate: assertion.Gate, EvidenceRequired: assertion.EvidenceRequired,
			EvidenceObserved: make([]model.Evidence, 0),
			RemediationClass: assertion.RemediationClass,
		}
		if planned.Applicability == "not_applicable" {
			result.Execution = "not_run"
			result.Assessment = "not_applicable"
			result.Summary = "Assertion is not applicable to the detected inventory."
		} else if planned.Applicability == "undetermined" {
			result.Execution = "not_run"
			result.Assessment = "unknown"
			result.Summary = "Applicability could not be determined by the supported expression evaluator."
		} else if assertion.ImplementationID == "prc.native.analysis-evidence@0.1" {
			result = evaluateAnalysisEvidence(assertion, inventory, validatedExecutions, result)
		} else {
			result = e.evaluate(assertion, inventory, result, started)
		}
		run.Results = append(run.Results, result)
	}
	run.CompletedAt = e.Now()
	run.TerminalState = terminalState(profile, run.Results)
	runPayload, err := json.Marshal(run)
	if err != nil {
		return model.RunResult{}, fmt.Errorf("encode run identity: %w", err)
	}
	runDigest := sha256.Sum256(runPayload)
	run.RunID = hex.EncodeToString(runDigest[:])
	return run, nil
}

func evaluateApplicability(expression string, inventory model.Inventory) (string, string) {
	switch strings.TrimSpace(expression) {
	case "true":
		return "applicable", ""
	case "false":
		return "not_applicable", ""
	case "inventory.package_ecosystems.size() > 0":
		if len(inventory.PackageEcosystems) > 0 {
			return "applicable", ""
		}
		return "not_applicable", ""
	case "inventory.ci.github_actions == true":
		if inventory.CI.GitHubActions {
			return "applicable", ""
		}
		return "not_applicable", ""
	case "inventory.source_files > 0":
		if inventory.SourceFiles > 0 {
			return "applicable", ""
		}
		return "not_applicable", ""
	case "inventory.container_files.size() > 0":
		if len(inventory.ContainerFiles) > 0 {
			return "applicable", ""
		}
		return "not_applicable", ""
	case "inventory.infrastructure.terraform_files.size() > 0":
		if len(inventory.Infrastructure.TerraformFiles) > 0 {
			return "applicable", ""
		}
		return "not_applicable", ""
	case "inventory.infrastructure.kubernetes_files.size() > 0":
		if len(inventory.Infrastructure.KubernetesFiles) > 0 {
			return "applicable", ""
		}
		return "not_applicable", ""
	default:
		return "undetermined", fmt.Sprintf("unsupported applicability expression %q", expression)
	}
}

func (e *Engine) evaluate(
	assertion model.Assertion,
	inventory model.Inventory,
	result model.AssertionResult,
	observedAt time.Time,
) model.AssertionResult {
	result.Execution = "completed"
	switch assertion.ImplementationID {
	case "prc.native.file-present@0.1":
		return evaluateFilePresent(assertion, inventory, result, observedAt)
	case "prc.native.dependency-lock@0.1":
		return evaluateDependencyLocks(assertion, inventory, result, observedAt)
	case "prc.native.github-action-pin@0.1":
		return evaluateActionPins(assertion, inventory, result, observedAt)
	case "prc.native.ci-present@0.1":
		return evaluateCIPresent(assertion, inventory, result, observedAt)
	case "prc.native.test-suite@0.1":
		return evaluateTestSuite(assertion, inventory, result, observedAt)
	case "prc.native.github-workflow-permissions@0.1":
		return evaluateWorkflowPermissions(assertion, inventory, result, observedAt)
	case "prc.native.final-newline@0.1":
		return evaluateFinalNewline(assertion, inventory, result, observedAt)
	case "prc.native.git-revision@0.1":
		return evaluateGitRevision(assertion, inventory, result, observedAt)
	case "prc.native.github-workflow-valid@0.1":
		return evaluateWorkflowValidity(assertion, inventory, result, observedAt)
	case "prc.native.github-workflow-jobs@0.1":
		return evaluateWorkflowJobs(assertion, inventory, result, observedAt)
	case "prc.native.github-workflow-timeouts@0.1":
		return evaluateWorkflowTimeouts(assertion, inventory, result, observedAt)
	case "prc.native.github-no-pull-request-target@0.1":
		return evaluateNoPullRequestTarget(assertion, inventory, result, observedAt)
	case "prc.native.merge-conflict-markers@0.1":
		return evaluateMergeConflictMarkers(assertion, inventory, result, observedAt)
	case "prc.native.restrictive-file-modes@0.1":
		return evaluateRestrictiveFileModes(assertion, inventory, result, observedAt)
	case "prc.native.inventory-files-nonempty@0.1":
		return evaluateInventoryFilesNonempty(assertion, inventory, result, observedAt)
	case "prc.native.runtime-version@0.1":
		return evaluateRuntimeVersions(assertion, inventory, result, observedAt)
	case "prc.native.container-base-pin@0.1":
		return evaluateContainerBasePins(assertion, inventory, result, observedAt)
	case "prc.native.container-nonroot@0.1":
		return evaluateContainerNonRoot(assertion, inventory, result, observedAt)
	case "prc.native.terraform-lock@0.1":
		return evaluateTerraformLocks(assertion, inventory, result, observedAt)
	case "prc.native.kubernetes-nonroot@0.1":
		return evaluateKubernetesNonRoot(assertion, inventory, result, observedAt)
	case "prc.native.kubernetes-resources@0.1":
		return evaluateKubernetesResources(assertion, inventory, result, observedAt)
	case "prc.native.manual-evidence@0.1":
		result.Assessment = "manual_review"
		result.Summary = "This assertion requires scoped evidence from an accountable reviewer."
		return result
	case "prc.native.analysis-evidence@0.1":
		result.Execution = "error"
		result.Assessment = "unknown"
		result.Summary = "Analysis evidence must be evaluated through the bound adapter-evidence path."
		return result
	default:
		result.Execution = "error"
		result.Assessment = "unknown"
		result.Summary = "Unsupported implementation: " + assertion.ImplementationID
		return result
	}
}

func stringSlice(value any) ([]string, error) {
	items, ok := value.([]any)
	if !ok {
		if strings, ok := value.([]string); ok {
			return strings, nil
		}
		return nil, fmt.Errorf("expected a string list")
	}
	result := make([]string, 0, len(items))
	for _, item := range items {
		text, ok := item.(string)
		if !ok {
			return nil, fmt.Errorf("expected a string list")
		}
		result = append(result, text)
	}
	return result, nil
}

func safePath(root, relative string) (string, error) {
	if relative == "" || filepath.IsAbs(relative) {
		return "", fmt.Errorf("unsafe repository path %q", relative)
	}
	clean := filepath.Clean(filepath.FromSlash(relative))
	if clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("repository path escapes target: %q", relative)
	}
	resolved := filepath.Join(root, clean)
	rootWithSeparator := filepath.Clean(root) + string(filepath.Separator)
	if resolved != filepath.Clean(root) && !strings.HasPrefix(resolved, rootWithSeparator) {
		return "", fmt.Errorf("repository path escapes target: %q", relative)
	}
	realRoot, rootErr := filepath.EvalSymlinks(root)
	realResolved, resolvedErr := filepath.EvalSymlinks(resolved)
	if rootErr == nil && resolvedErr == nil {
		realRootWithSeparator := filepath.Clean(realRoot) + string(filepath.Separator)
		if realResolved != filepath.Clean(realRoot) && !strings.HasPrefix(realResolved, realRootWithSeparator) {
			return "", fmt.Errorf("repository symlink escapes target: %q", relative)
		}
	}
	return resolved, nil
}

func fileEvidence(
	inventory model.Inventory,
	relative, kind, producer, summary string,
	observedAt time.Time,
) (model.Evidence, error) {
	evidence, _, _, err := fileEvidenceAndLastByte(inventory, relative, kind, producer, summary, observedAt)
	return evidence, err
}

func fileEvidenceAndLastByte(
	inventory model.Inventory,
	relative, kind, producer, summary string,
	observedAt time.Time,
) (model.Evidence, byte, bool, error) {
	path, err := safePath(inventory.Root, relative)
	if err != nil {
		return model.Evidence{}, 0, false, err
	}
	file, err := os.Open(path)
	if err != nil {
		return model.Evidence{}, 0, false, err
	}
	defer file.Close()
	hasher := sha256.New()
	buffer := make([]byte, 32*1024)
	var size int64
	var lastByte byte
	hasContent := false
	for {
		count, readErr := file.Read(buffer)
		if count > 0 {
			hasContent = true
			lastByte = buffer[count-1]
			size += int64(count)
			if _, err := hasher.Write(buffer[:count]); err != nil {
				return model.Evidence{}, 0, false, fmt.Errorf("hash %s: %w", relative, err)
			}
		}
		if readErr != nil {
			if readErr == io.EOF {
				break
			}
			return model.Evidence{}, 0, false, fmt.Errorf("read %s: %w", relative, readErr)
		}
	}
	contentDigest := hex.EncodeToString(hasher.Sum(nil))
	expectedDigest := ""
	for _, record := range inventory.Files {
		if record.Path == filepath.ToSlash(relative) {
			expectedDigest = record.SHA256
			break
		}
	}
	if expectedDigest == "" || expectedDigest != contentDigest {
		return model.Evidence{}, 0, false, fmt.Errorf("target changed after inventory: %s", relative)
	}
	evidence := model.Evidence{
		SchemaVersion: model.EvidenceSchema, Kind: kind, Authority: "repository",
		Producer: producer, TargetDigest: inventory.Digest, Source: filepath.ToSlash(relative),
		ContentSHA256: contentDigest, Size: size, ObservedAt: observedAt, Summary: summary,
	}
	identity, err := json.Marshal(evidence)
	if err != nil {
		return model.Evidence{}, 0, false, err
	}
	digest := sha256.Sum256(identity)
	evidence.ID = hex.EncodeToString(digest[:])
	return evidence, lastByte, hasContent, nil
}

func evaluateFilePresent(
	assertion model.Assertion,
	inventory model.Inventory,
	result model.AssertionResult,
	observedAt time.Time,
) model.AssertionResult {
	paths, err := stringSlice(assertion.Parameters["paths"])
	if err != nil {
		result.Execution, result.Assessment, result.Summary = "error", "unknown", err.Error()
		return result
	}
	minimumBytes := int64(1)
	if value, ok := assertion.Parameters["minimum_bytes"].(int); ok {
		minimumBytes = int64(value)
	}
	for _, relative := range paths {
		path, err := safePath(inventory.Root, relative)
		if err != nil {
			result.Execution, result.Assessment, result.Summary = "error", "unknown", err.Error()
			return result
		}
		info, err := os.Stat(path)
		if err != nil || !info.Mode().IsRegular() || info.Size() < minimumBytes {
			continue
		}
		evidence, err := fileEvidence(inventory, relative, "repository-file", assertion.ImplementationID,
			"Required repository file is present and nonempty.", observedAt)
		if err != nil {
			result.Execution, result.Assessment, result.Summary = "error", "unknown", err.Error()
			return result
		}
		result.Assessment = "pass"
		result.Summary = "Observed " + filepath.ToSlash(relative) + "."
		result.EvidenceObserved = []model.Evidence{evidence}
		return result
	}
	result.Assessment = "fail"
	result.Summary = "None of the required repository files were present and nonempty."
	return result
}

func evaluateDependencyLocks(
	assertion model.Assertion,
	inventory model.Inventory,
	result model.AssertionResult,
	observedAt time.Time,
) model.AssertionResult {
	locks := map[string][]string{"node": {}, "python": {}, "go": {}, "rust": {}, "java": {}}
	for _, relative := range inventory.LockFiles {
		name := filepath.Base(relative)
		switch {
		case name == "package-lock.json", name == "npm-shrinkwrap.json", name == "pnpm-lock.yaml", name == "yarn.lock":
			locks["node"] = append(locks["node"], relative)
		case name == "uv.lock", name == "Pipfile.lock", name == "poetry.lock", strings.HasSuffix(name, ".lock.txt"):
			locks["python"] = append(locks["python"], relative)
		case name == "go.sum":
			locks["go"] = append(locks["go"], relative)
		case name == "Cargo.lock":
			locks["rust"] = append(locks["rust"], relative)
		case name == "gradle.lockfile":
			locks["java"] = append(locks["java"], relative)
		}
	}
	missing := []string{}
	for _, ecosystem := range inventory.PackageEcosystems {
		if len(locks[ecosystem]) == 0 {
			missing = append(missing, ecosystem)
		}
	}
	if len(missing) > 0 {
		result.Assessment = "fail"
		result.Summary = "No supported lock or checksum file was found for: " + strings.Join(missing, ", ") + "."
		return result
	}
	for _, ecosystem := range inventory.PackageEcosystems {
		for _, relative := range locks[ecosystem] {
			evidence, err := fileEvidence(inventory, relative, "repository-file", assertion.ImplementationID,
				"Dependency resolution file for "+ecosystem+".", observedAt)
			if err != nil {
				result.Execution, result.Assessment, result.Summary = "error", "unknown", err.Error()
				return result
			}
			result.EvidenceObserved = append(result.EvidenceObserved, evidence)
		}
	}
	result.Assessment = "pass"
	result.Summary = "Every detected package ecosystem has a supported lock or checksum file."
	return result
}

func workflowPaths(inventory model.Inventory) []string {
	paths := []string{}
	for _, file := range inventory.Files {
		relative := file.Path
		if strings.HasPrefix(relative, ".github/workflows/") &&
			(strings.HasSuffix(relative, ".yml") || strings.HasSuffix(relative, ".yaml")) {
			paths = append(paths, relative)
		}
	}
	return paths
}

func evaluateActionPins(
	assertion model.Assertion,
	inventory model.Inventory,
	result model.AssertionResult,
	observedAt time.Time,
) model.AssertionResult {
	unpinned := []string{}
	for _, relative := range workflowPaths(inventory) {
		path, err := safePath(inventory.Root, relative)
		if err != nil {
			result.Execution, result.Assessment, result.Summary = "error", "unknown", err.Error()
			return result
		}
		data, err := os.ReadFile(path)
		if err != nil {
			result.Execution, result.Assessment, result.Summary = "error", "unknown", err.Error()
			return result
		}
		for _, match := range actionUsePattern.FindAllStringSubmatch(string(data), -1) {
			reference := match[1]
			if strings.HasPrefix(reference, "./") || strings.HasPrefix(reference, "docker://") {
				continue
			}
			at := strings.LastIndex(reference, "@")
			if at < 0 || !commitPattern.MatchString(reference[at+1:]) {
				unpinned = append(unpinned, relative+":"+reference)
			}
		}
		evidence, err := fileEvidence(inventory, relative, "workflow-parse", assertion.ImplementationID,
			"Parsed workflow action references.", observedAt)
		if err != nil {
			result.Execution, result.Assessment, result.Summary = "error", "unknown", err.Error()
			return result
		}
		result.EvidenceObserved = append(result.EvidenceObserved, evidence)
	}
	sort.Strings(unpinned)
	if len(unpinned) > 0 {
		result.Assessment = "fail"
		result.Summary = "Mutable or invalid action references: " + strings.Join(unpinned, ", ") + "."
		return result
	}
	result.Assessment = "pass"
	result.Summary = "Every external GitHub Action reference is a full commit SHA."
	return result
}

func evaluateCIPresent(
	assertion model.Assertion,
	inventory model.Inventory,
	result model.AssertionResult,
	observedAt time.Time,
) model.AssertionResult {
	paths := workflowPaths(inventory)
	if len(paths) == 0 {
		result.Assessment = "fail"
		result.Summary = "No supported continuous-integration configuration was found."
		return result
	}
	for _, relative := range paths {
		evidence, err := fileEvidence(inventory, relative, "repository-file", assertion.ImplementationID,
			"Continuous-integration workflow is present.", observedAt)
		if err != nil {
			result.Execution, result.Assessment, result.Summary = "error", "unknown", err.Error()
			return result
		}
		result.EvidenceObserved = append(result.EvidenceObserved, evidence)
	}
	result.Assessment = "pass"
	result.Summary = "Observed supported continuous-integration configuration."
	return result
}

func evaluateTestSuite(
	assertion model.Assertion,
	inventory model.Inventory,
	result model.AssertionResult,
	observedAt time.Time,
) model.AssertionResult {
	candidates := []string{}
	for _, file := range inventory.Files {
		relative := file.Path
		parts := strings.Split(relative, "/")
		name := filepath.Base(relative)
		isSource := workspaceinventory.IsSourcePath(name)
		if isSource && ((len(parts) > 1 && (parts[0] == "tests" || parts[0] == "test" || parts[0] == "__tests__")) ||
			strings.HasSuffix(name, "_test.go") || strings.HasPrefix(name, "test_") || strings.HasSuffix(name, "_test.py")) {
			candidates = append(candidates, relative)
		}
	}
	if len(candidates) == 0 {
		result.Assessment = "fail"
		result.Summary = "No recognized test path or test file was found."
		return result
	}
	sort.Strings(candidates)
	evidence, err := fileEvidence(inventory, candidates[0], "repository-structure", assertion.ImplementationID,
		"Recognized test file is present.", observedAt)
	if err != nil {
		result.Execution, result.Assessment, result.Summary = "error", "unknown", err.Error()
		return result
	}
	result.EvidenceObserved = []model.Evidence{evidence}
	result.Assessment = "pass"
	result.Summary = fmt.Sprintf("Observed recognized tests; first deterministic evidence path is %s.", candidates[0])
	return result
}

func evaluateFinalNewline(
	assertion model.Assertion,
	inventory model.Inventory,
	result model.AssertionResult,
	observedAt time.Time,
) model.AssertionResult {
	const maximumEvidenceFiles = 2048
	paths := make([]string, 0)
	for _, file := range inventory.Files {
		if workspaceinventory.IsSourcePath(file.Path) {
			paths = append(paths, file.Path)
		}
	}
	if len(paths) > maximumEvidenceFiles {
		result.Execution = "blocked"
		result.Assessment = "unknown"
		result.Summary = fmt.Sprintf("Source-format evidence requires %d files, above the %d-file limit.", len(paths), maximumEvidenceFiles)
		return result
	}
	violations := make([]string, 0)
	for _, relative := range paths {
		evidence, lastByte, hasContent, err := fileEvidenceAndLastByte(
			inventory, relative, "source-format", assertion.ImplementationID,
			"Verified the final byte of a recognized source file.", observedAt,
		)
		if err != nil {
			result.Execution, result.Assessment, result.Summary = "error", "unknown", err.Error()
			return result
		}
		result.EvidenceObserved = append(result.EvidenceObserved, evidence)
		if !hasContent || lastByte != '\n' {
			violations = append(violations, relative)
		}
	}
	if len(violations) == 0 {
		result.Assessment = "pass"
		result.Summary = fmt.Sprintf("All %d recognized source files end with a line-feed byte.", len(paths))
		return result
	}
	result.Assessment = "fail"
	visible := violations
	if len(visible) > 10 {
		visible = visible[:10]
	}
	result.Summary = "Source files without a final line-feed byte: " + strings.Join(visible, ", ")
	if len(violations) > len(visible) {
		result.Summary += fmt.Sprintf(" (and %d more)", len(violations)-len(visible))
	}
	result.Summary += "."
	return result
}

func mappingValue(node *yaml.Node, key string) *yaml.Node {
	if node == nil {
		return nil
	}
	if node.Kind == yaml.DocumentNode && len(node.Content) > 0 {
		node = node.Content[0]
	}
	if node.Kind != yaml.MappingNode {
		return nil
	}
	for index := 0; index+1 < len(node.Content); index += 2 {
		if node.Content[index].Value == key {
			return node.Content[index+1]
		}
	}
	return nil
}

func permissionsExplicit(document *yaml.Node) bool {
	if permissions := mappingValue(document, "permissions"); permissions != nil {
		return !(permissions.Kind == yaml.ScalarNode && permissions.Value == "write-all")
	}
	jobs := mappingValue(document, "jobs")
	if jobs == nil || jobs.Kind != yaml.MappingNode || len(jobs.Content) == 0 {
		return false
	}
	for index := 1; index < len(jobs.Content); index += 2 {
		permissions := mappingValue(jobs.Content[index], "permissions")
		if permissions == nil || (permissions.Kind == yaml.ScalarNode && permissions.Value == "write-all") {
			return false
		}
	}
	return true
}

func evaluateWorkflowPermissions(
	assertion model.Assertion,
	inventory model.Inventory,
	result model.AssertionResult,
	observedAt time.Time,
) model.AssertionResult {
	missing := []string{}
	for _, relative := range workflowPaths(inventory) {
		path, err := safePath(inventory.Root, relative)
		if err != nil {
			result.Execution, result.Assessment, result.Summary = "error", "unknown", err.Error()
			return result
		}
		data, err := os.ReadFile(path)
		if err != nil {
			result.Execution, result.Assessment, result.Summary = "error", "unknown", err.Error()
			return result
		}
		var document yaml.Node
		if err := yaml.Unmarshal(data, &document); err != nil {
			result.Execution, result.Assessment = "error", "unknown"
			result.Summary = fmt.Sprintf("Cannot parse %s: %v", relative, err)
			return result
		}
		if !permissionsExplicit(&document) {
			missing = append(missing, relative)
		}
		evidence, err := fileEvidence(inventory, relative, "workflow-parse", assertion.ImplementationID,
			"Parsed workflow permission declarations.", observedAt)
		if err != nil {
			result.Execution, result.Assessment, result.Summary = "error", "unknown", err.Error()
			return result
		}
		result.EvidenceObserved = append(result.EvidenceObserved, evidence)
	}
	if len(missing) > 0 {
		result.Assessment = "fail"
		result.Summary = "Workflows without explicit least-privilege permissions: " + strings.Join(missing, ", ") + "."
		return result
	}
	result.Assessment = "pass"
	result.Summary = "Every GitHub Actions workflow declares explicit permissions."
	return result
}

func terminalState(profile model.Profile, results []model.AssertionResult) string {
	blockingSeverity := map[string]bool{}
	for _, severity := range profile.TerminalPolicy.BlockOn {
		blockingSeverity[severity] = true
	}
	hasFailure, hasManual, hasIncomplete, environmentBlocked := false, false, false, false
	for _, result := range results {
		if result.Applicability == "not_applicable" || result.Gate == "advisory" {
			continue
		}
		if result.Assessment == "fail" {
			hasFailure = true
			if result.Gate == "no-go" || blockingSeverity[result.Severity] {
				return "no_go"
			}
		}
		if result.Assessment == "manual_review" {
			hasManual = true
		}
		if result.Assessment == "unknown" || result.Assessment == "stale" ||
			result.Assessment == "conflicting" || result.Applicability == "undetermined" {
			hasIncomplete = true
		}
		if result.Execution == "error" || result.Execution == "blocked" {
			environmentBlocked = true
		}
	}
	if environmentBlocked {
		return "environment_blocked"
	}
	if hasFailure || hasIncomplete {
		return "assessment_incomplete"
	}
	if hasManual {
		if !profile.TerminalPolicy.AllowManualRemaining {
			return "assessment_incomplete"
		}
		return "machine_work_complete_manual_evidence_remaining"
	}
	return "profile_satisfied"
}
