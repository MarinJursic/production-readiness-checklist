// Package invalidation determines which prior assertion conclusions remain
// input-equivalent after a target changes. It never rewrites or silently
// rebinds evidence; callers must perform a fresh evaluation when reuse is not
// explicitly allowed.
package invalidation

import (
	"fmt"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"time"

	workspaceinventory "github.com/MarinJursic/production-readiness-checklist/scanner/inventory"
	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
)

const (
	SchemaVersion = "prc.invalidation/v0.1"
	RulesVersion  = "prc.invalidation-rules/v0.1"
)

type FileChange struct {
	Path           string `json:"path"`
	Change         string `json:"change"`
	BeforeSHA256   string `json:"before_sha256,omitempty"`
	AfterSHA256    string `json:"after_sha256,omitempty"`
	BeforeSize     int64  `json:"before_size,omitempty"`
	AfterSize      int64  `json:"after_size,omitempty"`
	BeforeMode     uint32 `json:"before_mode,omitempty"`
	AfterMode      uint32 `json:"after_mode,omitempty"`
	ContentChanged bool   `json:"content_changed"`
	ModeChanged    bool   `json:"mode_changed"`
}

type Reason struct {
	Code   string `json:"code"`
	Detail string `json:"detail"`
}

type AssertionImpact struct {
	AssertionID             string   `json:"assertion_id"`
	PreviousAssessment      string   `json:"previous_assessment,omitempty"`
	PreviousApplicability   string   `json:"previous_applicability,omitempty"`
	CurrentApplicability    string   `json:"current_applicability,omitempty"`
	Conclusion              string   `json:"conclusion"`
	ReuseAllowed            bool     `json:"reuse_allowed"`
	FreshEvaluationRequired bool     `json:"fresh_evaluation_required"`
	RelevantChangedPaths    []string `json:"relevant_changed_paths"`
	Reasons                 []Reason `json:"reasons"`
}

type Summary struct {
	Invalidated     int `json:"invalidated"`
	UnchangedInputs int `json:"unchanged_inputs"`
	Reusable        int `json:"reusable"`
	New             int `json:"new"`
	Removed         int `json:"removed"`
}

type Report struct {
	SchemaVersion          string            `json:"schema_version"`
	RulesVersion           string            `json:"rules_version"`
	GeneratedAt            time.Time         `json:"generated_at"`
	TargetName             string            `json:"target_name"`
	BaseRunID              string            `json:"base_run_id"`
	BaseInventoryDigest    string            `json:"base_inventory_digest"`
	CurrentInventoryDigest string            `json:"current_inventory_digest"`
	BasePlanDigest         string            `json:"base_plan_digest"`
	CurrentPlanDigest      string            `json:"current_plan_digest"`
	ChangedDimensions      []string          `json:"changed_dimensions"`
	ChangedFiles           []FileChange      `json:"changed_files"`
	Assertions             []AssertionImpact `json:"assertions"`
	Summary                Summary           `json:"summary"`
}

type dependency struct {
	paths       map[string]bool
	content     bool
	mode        bool
	presence    bool
	allFiles    bool
	alwaysFresh bool
}

func Analyze(
	base model.RunResult,
	currentInventory model.Inventory,
	currentPlan model.Plan,
	assertions map[string]model.Assertion,
	now time.Time,
) (Report, error) {
	if base.RunID == "" || base.Inventory.Digest == "" || currentInventory.Digest == "" {
		return Report{}, fmt.Errorf("base run and current inventory must have canonical identities")
	}
	if base.Plan.InventoryDigest != base.Inventory.Digest || currentPlan.InventoryDigest != currentInventory.Digest {
		return Report{}, fmt.Errorf("plan and inventory identities do not match")
	}
	if (base.Plan.ProjectID != "" || currentPlan.ProjectID != "") && base.Plan.ProjectID != currentPlan.ProjectID {
		return Report{}, fmt.Errorf("base project %q does not match current project %q", base.Plan.ProjectID, currentPlan.ProjectID)
	}
	if base.Plan.ProjectID == "" && currentPlan.ProjectID == "" && base.Inventory.TargetName != currentInventory.TargetName {
		return Report{}, fmt.Errorf("base target %q does not match current target %q", base.Inventory.TargetName, currentInventory.TargetName)
	}

	changes := compareFiles(base.Inventory.Files, currentInventory.Files)
	dimensions := changedDimensions(base.Inventory, currentInventory, base.Plan, currentPlan)
	report := Report{
		SchemaVersion: SchemaVersion, RulesVersion: RulesVersion, GeneratedAt: now.UTC(),
		TargetName: currentInventory.TargetName, BaseRunID: base.RunID,
		BaseInventoryDigest: base.Inventory.Digest, CurrentInventoryDigest: currentInventory.Digest,
		BasePlanDigest: base.Plan.Digest, CurrentPlanDigest: currentPlan.Digest,
		ChangedDimensions: dimensions, ChangedFiles: changes, Assertions: []AssertionImpact{},
	}

	basePlans := map[string]model.PlannedAssertion{}
	for _, planned := range base.Plan.Assertions {
		basePlans[planned.AssertionID] = planned
	}
	currentPlans := map[string]model.PlannedAssertion{}
	for _, planned := range currentPlan.Assertions {
		currentPlans[planned.AssertionID] = planned
	}
	baseResults := map[string]model.AssertionResult{}
	for _, result := range base.Results {
		baseResults[result.AssertionID] = result
	}

	for _, current := range currentPlan.Assertions {
		impact := AssertionImpact{
			AssertionID: current.AssertionID, CurrentApplicability: current.Applicability,
			Conclusion: "invalidated", ReuseAllowed: false, FreshEvaluationRequired: true,
			RelevantChangedPaths: []string{}, Reasons: []Reason{},
		}
		previous, existed := basePlans[current.AssertionID]
		if !existed {
			impact.Conclusion = "new"
			impact.Reasons = append(impact.Reasons, reason("assertion_added", "The current profile newly includes this assertion."))
			report.Summary.New++
			report.Assertions = append(report.Assertions, impact)
			continue
		}
		impact.PreviousApplicability = previous.Applicability
		priorResult, hasResult := baseResults[current.AssertionID]
		if hasResult {
			impact.PreviousAssessment = priorResult.Assessment
		} else {
			impact.Reasons = append(impact.Reasons, reason("prior_result_missing", "The base run has no result for this planned assertion."))
		}
		if base.Plan.EngineVersion == "" || base.Plan.ProfileDigest == "" || previous.DefinitionDigest == "" || previous.AssertionRevision == 0 {
			impact.Reasons = append(impact.Reasons, reason(
				"base_plan_lacks_rule_binding",
				"The base plan predates cryptographic engine, profile, or assertion-definition binding.",
			))
		}
		if base.Plan.EngineVersion != currentPlan.EngineVersion {
			impact.Reasons = append(impact.Reasons, reason("engine_version_changed", "The scanner engine contract changed."))
		}
		if base.Inventory.SchemaVersion != currentInventory.SchemaVersion {
			impact.Reasons = append(impact.Reasons, reason("inventory_schema_changed", "The inventory detection contract changed."))
		}
		if len(changes) == 0 && (!reflect.DeepEqual(base.Inventory.Facts, currentInventory.Facts) ||
			!reflect.DeepEqual(base.Inventory.Components, currentInventory.Components) ||
			!reflect.DeepEqual(base.Inventory.Relations, currentInventory.Relations)) {
			impact.Reasons = append(impact.Reasons, reason("inventory_detector_output_changed", "Sourced facts or graph structure changed without a repository file change."))
		}
		if base.Plan.ProfileID != currentPlan.ProfileID || base.Plan.ProfileVersion != currentPlan.ProfileVersion || base.Plan.ProfileDigest != currentPlan.ProfileDigest {
			impact.Reasons = append(impact.Reasons, reason("profile_definition_changed", "The selected profile identity or definition changed."))
		}
		if previous.Implementation != current.Implementation || previous.AssertionRevision != current.AssertionRevision || previous.DefinitionDigest != current.DefinitionDigest {
			impact.Reasons = append(impact.Reasons, reason("assertion_definition_changed", "The assertion revision, parameters, or implementation binding changed."))
		}
		if previous.ApplicabilityBy != current.ApplicabilityBy || previous.Applicability != current.Applicability {
			impact.Reasons = append(impact.Reasons, reason("applicability_changed", "The applicability evaluator or outcome changed."))
		}
		if base.Plan.ConfigurationDigest != currentPlan.ConfigurationDigest ||
			base.Plan.ProjectID != currentPlan.ProjectID ||
			!reflect.DeepEqual(base.Plan.ArtifactDigests, currentPlan.ArtifactDigests) ||
			!reflect.DeepEqual(base.Plan.TargetEnvironments, currentPlan.TargetEnvironments) {
			impact.Reasons = append(impact.Reasons, reason("declared_scope_changed", "Configuration, artifacts, project identity, or target environments changed."))
		}
		assertion, known := assertions[current.AssertionID]
		if !known {
			impact.Reasons = append(impact.Reasons, reason("assertion_definition_missing", "The current assertion definition is unavailable."))
		} else {
			dep := dependencies(assertion, base.Inventory, currentInventory)
			impact.RelevantChangedPaths = relevantChanges(changes, dep)
			if len(impact.RelevantChangedPaths) > 0 {
				impact.Reasons = append(impact.Reasons, reason("relevant_repository_input_changed", "One or more repository inputs used by this implementation changed."))
			}
			if dep.alwaysFresh {
				impact.Reasons = append(impact.Reasons, reason("freshness_not_proven", "The assertion depends on evidence whose validity window cannot be proven from repository identity alone."))
			}
			if assertion.ImplementationID == "prc.native.git-revision@0.1" && base.Inventory.GitCommit != currentInventory.GitCommit {
				impact.Reasons = append(impact.Reasons, reason("git_revision_changed", "The immutable repository revision identity changed."))
			}
		}
		if hasResult && hasNonRepositoryEvidence(priorResult) {
			impact.Reasons = append(impact.Reasons, reason("non_repository_evidence_requires_revalidation", "Prior executed, artifact, environment, or human evidence is not implicitly carried forward."))
		}
		impact.Reasons = deduplicateReasons(impact.Reasons)
		if len(impact.Reasons) == 0 {
			impact.Conclusion = "unchanged_inputs"
			if base.Inventory.Digest == currentInventory.Digest && base.Plan.Digest == currentPlan.Digest {
				impact.ReuseAllowed = true
				impact.FreshEvaluationRequired = false
				report.Summary.Reusable++
			} else {
				impact.Reasons = append(impact.Reasons, reason(
					"evidence_target_requires_rebinding",
					"The implementation inputs are unchanged, but prior evidence is bound to a different complete inventory digest.",
				))
				report.Summary.UnchangedInputs++
			}
		} else {
			report.Summary.Invalidated++
		}
		report.Assertions = append(report.Assertions, impact)
	}

	for _, previous := range base.Plan.Assertions {
		if _, exists := currentPlans[previous.AssertionID]; exists {
			continue
		}
		impact := AssertionImpact{
			AssertionID: previous.AssertionID, PreviousApplicability: previous.Applicability,
			Conclusion: "removed", ReuseAllowed: false, FreshEvaluationRequired: false,
			RelevantChangedPaths: []string{}, Reasons: []Reason{reason("assertion_removed", "The current profile no longer includes this assertion.")},
		}
		if result, ok := baseResults[previous.AssertionID]; ok {
			impact.PreviousAssessment = result.Assessment
		}
		report.Summary.Removed++
		report.Assertions = append(report.Assertions, impact)
	}
	sort.Slice(report.Assertions, func(i, j int) bool { return report.Assertions[i].AssertionID < report.Assertions[j].AssertionID })
	return report, nil
}

func compareFiles(before, after []model.FileRecord) []FileChange {
	left, right := map[string]model.FileRecord{}, map[string]model.FileRecord{}
	for _, file := range before {
		left[file.Path] = file
	}
	for _, file := range after {
		right[file.Path] = file
	}
	paths := map[string]bool{}
	for path := range left {
		paths[path] = true
	}
	for path := range right {
		paths[path] = true
	}
	ordered := make([]string, 0, len(paths))
	for path := range paths {
		ordered = append(ordered, path)
	}
	sort.Strings(ordered)
	changes := []FileChange{}
	for _, path := range ordered {
		old, oldOK := left[path]
		current, currentOK := right[path]
		if oldOK && currentOK && old.SHA256 == current.SHA256 && old.Size == current.Size && old.Mode == current.Mode {
			continue
		}
		change := FileChange{Path: path}
		switch {
		case !oldOK:
			change.Change, change.ContentChanged = "added", true
		case !currentOK:
			change.Change, change.ContentChanged = "deleted", true
		case old.SHA256 != current.SHA256 || old.Size != current.Size:
			change.Change, change.ContentChanged = "modified", true
		default:
			change.Change = "mode_changed"
		}
		if oldOK {
			change.BeforeSHA256, change.BeforeSize, change.BeforeMode = old.SHA256, old.Size, old.Mode
		}
		if currentOK {
			change.AfterSHA256, change.AfterSize, change.AfterMode = current.SHA256, current.Size, current.Mode
		}
		change.ModeChanged = oldOK && currentOK && old.Mode != current.Mode
		changes = append(changes, change)
	}
	return changes
}

func changedDimensions(before, after model.Inventory, basePlan, currentPlan model.Plan) []string {
	values := []struct {
		name    string
		changed bool
	}{
		{"inventory_schema", before.SchemaVersion != after.SchemaVersion},
		{"git_commit", before.GitCommit != after.GitCommit},
		{"files", !reflect.DeepEqual(before.Files, after.Files)},
		{"source_files", before.SourceFiles != after.SourceFiles},
		{"package_ecosystems", !reflect.DeepEqual(before.PackageEcosystems, after.PackageEcosystems)},
		{"manifests", !reflect.DeepEqual(before.Manifests, after.Manifests)},
		{"lock_files", !reflect.DeepEqual(before.LockFiles, after.LockFiles)},
		{"container_files", !reflect.DeepEqual(before.ContainerFiles, after.ContainerFiles)},
		{"symlinks", !reflect.DeepEqual(before.Symlinks, after.Symlinks)},
		{"ci", !reflect.DeepEqual(before.CI, after.CI)},
		{"infrastructure", !reflect.DeepEqual(before.Infrastructure, after.Infrastructure)},
		{"components", !reflect.DeepEqual(before.Components, after.Components)},
		{"relations", !reflect.DeepEqual(before.Relations, after.Relations)},
		{"facts", !reflect.DeepEqual(before.Facts, after.Facts)},
		{"declared_scope", !reflect.DeepEqual(before.DeclaredScope, after.DeclaredScope)},
		{"engine", basePlan.EngineVersion != currentPlan.EngineVersion},
		{"profile", basePlan.ProfileID != currentPlan.ProfileID || basePlan.ProfileVersion != currentPlan.ProfileVersion || basePlan.ProfileDigest != currentPlan.ProfileDigest},
	}
	result := []string{}
	for _, value := range values {
		if value.changed {
			result = append(result, value.name)
		}
	}
	return result
}

func dependencies(assertion model.Assertion, before, after model.Inventory) dependency {
	dep := dependency{paths: map[string]bool{}, content: true, presence: true}
	add := func(paths ...[]string) {
		for _, values := range paths {
			for _, path := range values {
				dep.paths[path] = true
			}
		}
	}
	allSource := func(inventory model.Inventory) []string {
		paths := []string{}
		for _, file := range inventory.Files {
			if workspaceinventory.IsSourcePath(file.Path) {
				paths = append(paths, file.Path)
			}
		}
		return paths
	}
	workflows := append(append([]string{}, before.CI.WorkflowFiles...), after.CI.WorkflowFiles...)
	switch assertion.ImplementationID {
	case "prc.native.file-present@0.1":
		if values, ok := assertion.Parameters["paths"].([]any); ok {
			for _, value := range values {
				if path, ok := value.(string); ok {
					dep.paths[path] = true
				}
			}
		} else if values, ok := assertion.Parameters["paths"].([]string); ok {
			add(values)
		}
	case "prc.native.dependency-lock@0.1":
		add(before.Manifests, after.Manifests, before.LockFiles, after.LockFiles)
	case "prc.native.github-action-pin@0.1", "prc.native.github-workflow-permissions@0.1",
		"prc.native.github-workflow-valid@0.1", "prc.native.github-workflow-jobs@0.1",
		"prc.native.github-workflow-timeouts@0.1", "prc.native.github-no-pull-request-target@0.1":
		add(workflows)
	case "prc.native.ci-present@0.1":
		add(workflows)
		dep.content = false
	case "prc.native.test-suite@0.1":
		for _, inventory := range []model.Inventory{before, after} {
			for _, file := range inventory.Files {
				if isTestPath(file.Path) {
					dep.paths[file.Path] = true
				}
			}
		}
	case "prc.native.final-newline@0.1":
		add(allSource(before), allSource(after))
	case "prc.native.git-revision@0.1":
		dep.content, dep.presence = false, false
	case "prc.native.merge-conflict-markers@0.1":
		add(allSource(before), allSource(after), before.Manifests, after.Manifests, before.LockFiles, after.LockFiles,
			workflows, before.ContainerFiles, after.ContainerFiles, before.Infrastructure.TerraformFiles,
			after.Infrastructure.TerraformFiles, before.Infrastructure.KubernetesFiles, after.Infrastructure.KubernetesFiles)
	case "prc.native.restrictive-file-modes@0.1":
		dep.allFiles, dep.content, dep.mode = true, false, true
	case "prc.native.inventory-files-nonempty@0.1":
		field, _ := assertion.Parameters["inventory_field"].(string)
		if field == "manifests" {
			add(before.Manifests, after.Manifests)
		} else {
			add(before.LockFiles, after.LockFiles)
		}
	case "prc.native.runtime-version@0.1":
		for _, inventory := range []model.Inventory{before, after} {
			for _, file := range inventory.Files {
				if isRuntimeCandidate(file.Path) {
					dep.paths[file.Path] = true
				}
			}
		}
	case "prc.native.container-base-pin@0.1", "prc.native.container-nonroot@0.1":
		add(before.ContainerFiles, after.ContainerFiles)
	case "prc.native.terraform-lock@0.1":
		add(before.Infrastructure.TerraformFiles, after.Infrastructure.TerraformFiles)
		for _, path := range append(append([]string{}, before.Infrastructure.TerraformFiles...), after.Infrastructure.TerraformFiles...) {
			directory := filepath.ToSlash(filepath.Dir(path))
			lock := filepath.ToSlash(filepath.Join(directory, ".terraform.lock.hcl"))
			if directory == "." {
				lock = ".terraform.lock.hcl"
			}
			dep.paths[lock] = true
		}
	case "prc.native.kubernetes-nonroot@0.1", "prc.native.kubernetes-resources@0.1":
		add(before.Infrastructure.KubernetesFiles, after.Infrastructure.KubernetesFiles)
	case "prc.native.manual-evidence@0.1":
		dep.content, dep.presence, dep.alwaysFresh = false, false, true
	case "prc.native.analysis-evidence@0.1":
		dep.allFiles, dep.alwaysFresh = true, true
	default:
		dep.allFiles = true
	}
	return dep
}

func relevantChanges(changes []FileChange, dep dependency) []string {
	result := []string{}
	for _, change := range changes {
		if !dep.allFiles && !dep.paths[change.Path] {
			continue
		}
		presenceChanged := change.Change == "added" || change.Change == "deleted"
		if (presenceChanged && dep.presence) || (change.ContentChanged && dep.content) || (change.ModeChanged && dep.mode) {
			result = append(result, change.Path)
		}
	}
	return result
}

func isTestPath(path string) bool {
	parts := strings.Split(path, "/")
	name := filepath.Base(path)
	return workspaceinventory.IsSourcePath(name) && ((len(parts) > 1 && (parts[0] == "tests" || parts[0] == "test" || parts[0] == "__tests__")) ||
		strings.HasSuffix(name, "_test.go") || strings.HasPrefix(name, "test_") || strings.HasSuffix(name, "_test.py"))
}

func isRuntimeCandidate(path string) bool {
	name := filepath.Base(path)
	return name == "go.mod" || name == "pyproject.toml" || name == "package.json" || name == "Cargo.toml" ||
		name == "pom.xml" || strings.HasPrefix(name, "build.gradle") || name == ".python-version" || name == "runtime.txt" ||
		name == ".nvmrc" || name == ".node-version" || name == ".tool-versions" || name == "rust-toolchain" ||
		name == "rust-toolchain.toml" || name == ".java-version" || name == ".sdkmanrc" || strings.Contains(path, ".github/workflows/")
}

func hasNonRepositoryEvidence(result model.AssertionResult) bool {
	for _, evidence := range result.EvidenceObserved {
		if evidence.Authority != "repository" {
			return true
		}
	}
	return false
}

func reason(code, detail string) Reason { return Reason{Code: code, Detail: detail} }

func deduplicateReasons(reasons []Reason) []Reason {
	seen, result := map[string]bool{}, []Reason{}
	for _, item := range reasons {
		if !seen[item.Code] {
			seen[item.Code] = true
			result = append(result, item)
		}
	}
	return result
}
