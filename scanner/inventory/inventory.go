package inventory

import (
	"bufio"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	projectconfig "github.com/MarinJursic/production-readiness-checklist/scanner/config"
	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
)

const (
	maxEntries             = 200_000
	maxInventoryFileBytes  = 1024 * 1024 * 1024
	maxInventoryTotalBytes = 8 * 1024 * 1024 * 1024
	maxStructuredFileBytes = 2 * 1024 * 1024
	maxGitReferenceBytes   = 4 * 1024
	maxPackedRefsBytes     = 8 * 1024 * 1024
	maxGitStatusBytes      = 32 * 1024 * 1024
	gitInspectionTimeout   = 10 * time.Second
)

const detectorVersion = "0.4"

var (
	kubernetesAPIVersion = regexp.MustCompile(`(?m)^apiVersion:[[:space:]]*[^[:space:]#]+`)
	kubernetesKind       = regexp.MustCompile(`(?m)^kind:[[:space:]]*[^[:space:]#]+`)
	openAPIYAMLMarker    = regexp.MustCompile(`(?m)^openapi[ \t]*:[ \t]*["']?3\.`)
	gitCommitPattern     = regexp.MustCompile(`^[0-9a-f]{40}(?:[0-9a-f]{24})?$`)
)

var excludedDirectories = map[string]string{
	".angular":          "Generated Angular cache content is not project source.",
	".aws-sam":          "Generated AWS SAM build content is not project source.",
	".build":            "Generated Swift package build content is not project source.",
	".dart_tool":        "Generated Dart and Flutter tool state is not project source.",
	".expo":             "Generated Expo tool state is not project source.",
	".git":              "Git metadata is inspected separately with bounded readers.",
	".gradle":           "Generated Gradle cache content is not project source.",
	".hypothesis":       "Generated Hypothesis test cache content is not project source.",
	".mypy_cache":       "Generated mypy cache content is not project source.",
	".next":             "Generated Next.js build and cache content is not project source.",
	".nox":              "Generated nox environment content is not project source.",
	".nuxt":             "Generated Nuxt build content is not project source.",
	".nyc_output":       "Generated test coverage process data is not project source.",
	".output":           "Generated framework build output is not project source.",
	".parcel-cache":     "Generated Parcel cache content is not project source.",
	".pnpm-store":       "Installed pnpm package-store content is third-party generated data.",
	".prc":              "Scanner-owned output is excluded to avoid scanning prior reports as project input.",
	".pytest_cache":     "Generated pytest cache content is not project source.",
	".ruff_cache":       "Generated Ruff cache content is not project source.",
	".serverless":       "Generated Serverless Framework package content is not project source.",
	".svelte-kit":       "Generated SvelteKit build content is not project source.",
	".terraform":        "Downloaded Terraform providers and module cache content are not project source.",
	".terragrunt-cache": "Downloaded Terragrunt working and module cache content is not project source.",
	".tox":              "Generated tox environment content is not project source.",
	".turbo":            "Generated Turborepo cache content is not project source.",
	".venv":             "A local Python virtual environment is third-party generated content.",
	".vite":             "Generated Vite cache content is not project source.",
	"DerivedData":       "Generated Apple build and indexing content is not project source.",
	"Pods":              "Installed CocoaPods dependencies are third-party generated content.",
	"__pycache__":       "Generated Python bytecode cache content is not project source.",
	"bower_components":  "Installed Bower dependencies are third-party generated content.",
	"node_modules":      "Installed Node dependencies are third-party generated content; lockfiles and manifests remain inventoried.",
	"venv":              "A local Python virtual environment is third-party generated content.",
}

func directoryExclusion(path string, entry fs.DirEntry) (string, bool) {
	if !entry.IsDir() {
		return "", false
	}
	if reason, excluded := excludedDirectories[entry.Name()]; excluded {
		return reason, true
	}
	parent := filepath.Dir(path)
	switch entry.Name() {
	case "target":
		if regularMarker(filepath.Join(path, "CACHEDIR.TAG")) ||
			regularMarker(filepath.Join(parent, "Cargo.toml")) || regularMarker(filepath.Join(parent, "pom.xml")) {
			return "A manifest or cache marker identifies this target directory as generated build output.", true
		}
	case "build":
		if regularMarker(filepath.Join(parent, "CMakeLists.txt")) ||
			regularMarker(filepath.Join(parent, "build.gradle")) || regularMarker(filepath.Join(parent, "build.gradle.kts")) ||
			regularMarker(filepath.Join(path, "CMakeCache.txt")) {
			return "A build-system marker identifies this build directory as generated output.", true
		}
	case "coverage":
		if regularMarker(filepath.Join(path, "coverage-final.json")) || regularMarker(filepath.Join(path, "lcov.info")) ||
			directoryMarker(filepath.Join(path, "lcov-report")) {
			return "Coverage-tool markers identify this directory as generated test output.", true
		}
	}
	return "", false
}

func regularMarker(path string) bool {
	info, err := os.Lstat(path)
	return err == nil && info.Mode().IsRegular()
}

func directoryMarker(path string) bool {
	info, err := os.Lstat(path)
	return err == nil && info.IsDir() && info.Mode()&os.ModeSymlink == 0
}

func inventoryFileLimitError(root, path string, size int64) error {
	return fmt.Errorf(
		"inventory stopped at %q: file size %s exceeds the %s inventory limit per file; target is %q. Run prc scan from the exact project root and remove generated or unrelated large files from that target",
		inventoryRelativePath(root, path), formatInventoryBytes(size), formatInventoryBytes(maxInventoryFileBytes), root,
	)
}

func inventoryTotalLimitError(root, path string, counted, next int64) error {
	return fmt.Errorf(
		"inventory stopped at %q after counting %s: adding the next %s file would exceed the %s total inventory limit; target is %q. Run prc scan from the exact project root, for example prc scan /path/to/project, and remove generated or unrelated large data from that target",
		inventoryRelativePath(root, path), formatInventoryBytes(counted), formatInventoryBytes(next),
		formatInventoryBytes(maxInventoryTotalBytes), root,
	)
}

func inventoryRelativePath(root, path string) string {
	relative, err := filepath.Rel(root, path)
	if err != nil || relative == "." || relative == "" || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return filepath.ToSlash(path)
	}
	return filepath.ToSlash(relative)
}

func formatInventoryBytes(value int64) string {
	const (
		kib = int64(1024)
		mib = 1024 * kib
		gib = 1024 * mib
	)
	switch {
	case value >= gib:
		return fmt.Sprintf("%.1f GiB", float64(value)/float64(gib))
	case value >= mib:
		return fmt.Sprintf("%.1f MiB", float64(value)/float64(mib))
	case value >= kib:
		return fmt.Sprintf("%.1f KiB", float64(value)/float64(kib))
	default:
		return fmt.Sprintf("%d bytes", value)
	}
}

var sourceExtensions = map[string]bool{
	".c": true, ".cc": true, ".cpp": true, ".cs": true, ".go": true,
	".java": true, ".js": true, ".jsx": true, ".kt": true, ".php": true,
	".py": true, ".rb": true, ".rs": true, ".scala": true, ".sh": true,
	".swift": true, ".ts": true, ".tsx": true,
}

func IsSourcePath(path string) bool {
	return sourceExtensions[strings.ToLower(filepath.Ext(path))]
}

func isOpenAPIFileName(name string) bool {
	switch strings.ToLower(name) {
	case "openapi.json", "openapi.yaml", "openapi.yml":
		return true
	default:
		return false
	}
}

func hasTopLevelOpenAPIMarker(data []byte) bool {
	return openAPIYAMLMarker.Match(data)
}

var manifests = map[string]string{
	"package.json":     "node",
	"pyproject.toml":   "python",
	"requirements.txt": "python",
	"setup.py":         "python",
	"go.mod":           "go",
	"Cargo.toml":       "rust",
	"pom.xml":          "java",
	"build.gradle":     "java",
	"build.gradle.kts": "java",
}

var lockFiles = map[string]bool{
	"package-lock.json": true, "npm-shrinkwrap.json": true, "pnpm-lock.yaml": true,
	"yarn.lock": true, "uv.lock": true, "Pipfile.lock": true, "poetry.lock": true,
	"go.sum": true, "Cargo.lock": true, "gradle.lockfile": true,
}

type digestInput struct {
	GitCommit         string                        `json:"git_commit,omitempty"`
	Files             []model.FileRecord            `json:"files"`
	PackageEcosystems []string                      `json:"package_ecosystems"`
	Manifests         []string                      `json:"manifests"`
	LockFiles         []string                      `json:"lock_files"`
	ContainerFiles    []string                      `json:"container_files"`
	Symlinks          []string                      `json:"symlinks"`
	CI                model.CIInventory             `json:"ci"`
	Infrastructure    model.InfrastructureInventory `json:"infrastructure"`
	Components        []model.InventoryComponent    `json:"components"`
	Relations         []model.InventoryRelation     `json:"relations"`
	Facts             []model.InventoryFact         `json:"facts"`
	DeclaredScope     *model.DeclaredScope          `json:"declared_scope,omitempty"`
}

func Build(target string) (model.Inventory, error) {
	root, err := filepath.Abs(target)
	if err != nil {
		return model.Inventory{}, fmt.Errorf("resolve target: %w", err)
	}
	info, err := os.Stat(root)
	if err != nil {
		return model.Inventory{}, fmt.Errorf("inspect target: %w", err)
	}
	if !info.IsDir() {
		return model.Inventory{}, fmt.Errorf("target is not a directory: %s", root)
	}
	if resolved, err := filepath.EvalSymlinks(root); err == nil {
		root = resolved
	}

	result := model.Inventory{
		SchemaVersion: model.InventorySchema, Root: root, TargetName: filepath.Base(root),
		Files: []model.FileRecord{}, PackageEcosystems: []string{}, Manifests: []string{}, LockFiles: []string{},
		ContainerFiles: []string{}, Symlinks: []string{},
		CI:             model.CIInventory{WorkflowFiles: []string{}},
		Infrastructure: model.InfrastructureInventory{TerraformFiles: []string{}, KubernetesFiles: []string{}},
		Components:     []model.InventoryComponent{{ID: "repository:.", Kind: "repository", Path: "."}},
		Relations:      []model.InventoryRelation{},
		Facts: []model.InventoryFact{{
			Key: "repository.detected", Value: "true", Source: ".", Detector: "prc.inventory.repository",
			DetectorVersion: detectorVersion, Confidence: 1, ScopePath: ".", Limitations: []string{},
		}},
	}
	ecosystems := map[string]bool{}
	componentIDs := map[string]bool{"repository:.": true}
	addComponent := func(kind, relative, ecosystem string) {
		id := kind + ":" + relative
		if componentIDs[id] {
			return
		}
		componentIDs[id] = true
		result.Components = append(result.Components, model.InventoryComponent{
			ID: id, Kind: kind, Path: relative, Ecosystem: ecosystem,
		})
		result.Relations = append(result.Relations, model.InventoryRelation{
			From: "repository:.", To: id, Kind: "contains",
		})
	}
	addFact := func(key, value, source, detector string, confidence float64, limitations ...string) {
		result.Facts = append(result.Facts, model.InventoryFact{
			Key: key, Value: value, Source: source, Detector: detector,
			DetectorVersion: detectorVersion, Confidence: confidence, ScopePath: filepath.ToSlash(filepath.Dir(source)),
			Limitations: append([]string{}, limitations...),
		})
	}
	visitedEntries := 0
	var totalBytes int64
	err = filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return fmt.Errorf("walk %s: %w", path, walkErr)
		}
		if path == root {
			return nil
		}
		visitedEntries++
		if visitedEntries > maxEntries {
			return fmt.Errorf("target exceeds %d filesystem entries", maxEntries)
		}
		if reason, excluded := directoryExclusion(path, entry); excluded {
			relative, relativeErr := filepath.Rel(root, path)
			if relativeErr != nil {
				return relativeErr
			}
			addFact("repository.exclusion", reason, filepath.ToSlash(relative), "prc.inventory.exclusion", 1,
				"Entries below this directory were not counted or hashed.")
			return filepath.SkipDir
		}
		if entry.IsDir() {
			return nil
		}
		if entry.Type()&os.ModeSymlink != 0 {
			relative, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}
			relative = filepath.ToSlash(relative)
			result.Symlinks = append(result.Symlinks, relative)
			addFact("repository.symlink", "true", relative, "prc.inventory.symlink", 1,
				"The link target is deliberately not followed or inventoried.")
			return nil
		}
		entryInfo, err := entry.Info()
		if err != nil {
			return fmt.Errorf("inspect %s: %w", path, err)
		}
		if !entryInfo.Mode().IsRegular() {
			return nil
		}
		if entryInfo.Size() > maxInventoryFileBytes {
			return inventoryFileLimitError(root, path, entryInfo.Size())
		}
		if entryInfo.Size() > maxInventoryTotalBytes-totalBytes {
			return inventoryTotalLimitError(root, path, totalBytes, entryInfo.Size())
		}
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		relative = filepath.ToSlash(relative)
		name := entry.Name()
		structuredCandidate := isYAML(name) || isOpenAPIFileName(name)
		fileRecord, structuredData, err := inspectFile(
			path, relative, structuredCandidate, maxStructuredFileBytes,
		)
		if err != nil {
			return err
		}
		if fileRecord.Size > maxInventoryTotalBytes-totalBytes {
			return inventoryTotalLimitError(root, path, totalBytes, fileRecord.Size)
		}
		result.Files = append(result.Files, fileRecord)
		totalBytes += fileRecord.Size
		ecosystem, manifest := manifests[name]
		if strings.HasPrefix(name, "requirements") && strings.HasSuffix(name, ".txt") &&
			!strings.HasSuffix(name, ".lock.txt") {
			ecosystem, manifest = "python", true
		}
		if manifest {
			result.Manifests = append(result.Manifests, relative)
			ecosystems[ecosystem] = true
			addComponent("package-manifest", relative, ecosystem)
			addFact("package.ecosystem", ecosystem, relative, "prc.inventory.package-manifest", 1,
				"Manifest presence does not prove that the component is built or deployed.")
		}
		if lockFiles[name] || strings.HasSuffix(name, ".lock.txt") {
			result.LockFiles = append(result.LockFiles, relative)
		}
		if IsSourcePath(name) {
			result.SourceFiles++
		}
		if strings.HasPrefix(relative, ".github/workflows/") &&
			(strings.HasSuffix(name, ".yml") || strings.HasSuffix(name, ".yaml")) {
			result.CI.GitHubActions = true
			result.CI.WorkflowFiles = append(result.CI.WorkflowFiles, relative)
			addComponent("ci-workflow", relative, "github-actions")
			addFact("ci.github_actions", "true", relative, "prc.inventory.github-actions", 1,
				"Workflow presence does not prove that a run completed for the assessed commit.")
		}
		if isContainerFile(name) {
			result.ContainerFiles = append(result.ContainerFiles, relative)
			addComponent("container-build", relative, "dockerfile")
			addFact("container.build_definition", "true", relative, "prc.inventory.container-file", 1,
				"A build definition does not prove that an image is built or deployed.")
		}
		if strings.HasSuffix(strings.ToLower(name), ".tf") {
			result.Infrastructure.TerraformFiles = append(result.Infrastructure.TerraformFiles, relative)
			addComponent("infrastructure", relative, "terraform")
			addFact("infrastructure.terraform", "true", relative, "prc.inventory.terraform-file", 0.98,
				"A Terraform file may be a module or example that is not applied.")
		}
		if isYAML(name) && kubernetesAPIVersion.Match(structuredData) && kubernetesKind.Match(structuredData) {
			result.Infrastructure.KubernetesFiles = append(result.Infrastructure.KubernetesFiles, relative)
			addComponent("infrastructure", relative, "kubernetes")
			addFact("infrastructure.kubernetes", "true", relative, "prc.inventory.kubernetes-yaml", 0.9,
				"Top-level Kubernetes fields do not prove that the resource is valid or deployed.")
		}
		if isOpenAPIFileName(name) || hasTopLevelOpenAPIMarker(structuredData) {
			addComponent("api-description", relative, "openapi")
			addFact("api.openapi", "true", relative, "prc.inventory.openapi-document", 0.95,
				"An OpenAPI document does not prove that the described API is implemented, reachable, or deployed.")
		}
		return nil
	})
	if err != nil {
		return model.Inventory{}, err
	}
	sort.Slice(result.Files, func(i, j int) bool { return result.Files[i].Path < result.Files[j].Path })
	sort.Strings(result.Manifests)
	sort.Strings(result.LockFiles)
	sort.Strings(result.ContainerFiles)
	sort.Strings(result.Symlinks)
	sort.Strings(result.CI.WorkflowFiles)
	sort.Strings(result.Infrastructure.TerraformFiles)
	sort.Strings(result.Infrastructure.KubernetesFiles)
	sort.Slice(result.Components, func(i, j int) bool { return result.Components[i].ID < result.Components[j].ID })
	sort.Slice(result.Relations, func(i, j int) bool {
		if result.Relations[i].From != result.Relations[j].From {
			return result.Relations[i].From < result.Relations[j].From
		}
		if result.Relations[i].To != result.Relations[j].To {
			return result.Relations[i].To < result.Relations[j].To
		}
		return result.Relations[i].Kind < result.Relations[j].Kind
	})
	sort.Slice(result.Facts, func(i, j int) bool {
		if result.Facts[i].Key != result.Facts[j].Key {
			return result.Facts[i].Key < result.Facts[j].Key
		}
		if result.Facts[i].Source != result.Facts[j].Source {
			return result.Facts[i].Source < result.Facts[j].Source
		}
		return result.Facts[i].Value < result.Facts[j].Value
	})
	for ecosystem := range ecosystems {
		result.PackageEcosystems = append(result.PackageEcosystems, ecosystem)
	}
	sort.Strings(result.PackageEcosystems)
	result.FileCount = len(result.Files)
	result.GitCommit, err = gitRepositoryState(root)
	if err != nil {
		addFact("repository.git_worktree_state", "unverified", ".", "prc.inventory.git-state", 0,
			"The Git working-tree state could not be established safely: "+err.Error())
	} else if result.GitCommit != "" {
		state := gitWorktreeState(root)
		addFact("repository.git_worktree_state", state.Value, ".", "prc.inventory.git-state", state.Confidence, state.Limitations...)
	}
	addFact("repository.inventory_bytes", fmt.Sprintf("%d", totalBytes), ".", "prc.inventory.bytes", 1)
	sort.Slice(result.Facts, func(i, j int) bool {
		if result.Facts[i].Key != result.Facts[j].Key {
			return result.Facts[i].Key < result.Facts[j].Key
		}
		if result.Facts[i].Source != result.Facts[j].Source {
			return result.Facts[i].Source < result.Facts[j].Source
		}
		return result.Facts[i].Value < result.Facts[j].Value
	})

	if err := seal(&result); err != nil {
		return model.Inventory{}, err
	}
	return result, nil
}

func seal(result *model.Inventory) error {
	digest, err := identity(*result)
	if err != nil {
		return err
	}
	result.Digest = digest
	return nil
}

func identity(result model.Inventory) (string, error) {
	payload, err := json.Marshal(digestInput{
		GitCommit: result.GitCommit, Files: result.Files, PackageEcosystems: result.PackageEcosystems,
		Manifests: result.Manifests, LockFiles: result.LockFiles, ContainerFiles: result.ContainerFiles,
		Symlinks: result.Symlinks, CI: result.CI, Infrastructure: result.Infrastructure,
		Components: result.Components, Relations: result.Relations, Facts: result.Facts,
		DeclaredScope: result.DeclaredScope,
	})
	if err != nil {
		return "", fmt.Errorf("encode inventory digest: %w", err)
	}
	digest := sha256.Sum256(payload)
	return hex.EncodeToString(digest[:]), nil
}

// VerifyIdentity recomputes the deterministic inventory digest without
// reading the workspace. It is safe to use for canonical stored records.
func VerifyIdentity(result model.Inventory) error {
	expected, err := identity(result)
	if err != nil {
		return err
	}
	if result.Digest != expected {
		return fmt.Errorf("inventory digest does not match record content")
	}
	return nil
}

// BindDerivedSource records the source revision from which a scanner-owned
// candidate copy was made. The candidate is deliberately marked dirty: this
// preserves provenance for checks that ask whether a source revision is known
// without pretending the uncommitted candidate itself is a clean checkout.
func BindDerivedSource(item, source model.Inventory) (model.Inventory, error) {
	if err := VerifyIdentity(item); err != nil {
		return model.Inventory{}, fmt.Errorf("verify derived candidate inventory: %w", err)
	}
	if err := VerifyIdentity(source); err != nil {
		return model.Inventory{}, fmt.Errorf("verify source inventory: %w", err)
	}
	if source.GitCommit == "" {
		return item, nil
	}
	if item.GitCommit != "" && item.GitCommit != source.GitCommit {
		return model.Inventory{}, fmt.Errorf("candidate belongs to a different Git revision")
	}
	item.GitCommit = source.GitCommit
	facts := item.Facts[:0]
	for _, fact := range item.Facts {
		if fact.Key != "repository.git_worktree_state" {
			facts = append(facts, fact)
		}
	}
	item.Facts = append(facts, model.InventoryFact{
		Key: "repository.git_worktree_state", Value: "dirty", Source: "scanner-derived-candidate",
		Detector: "prc.remediation.candidate-origin", DetectorVersion: "0.1", Confidence: 1, ScopePath: ".",
		Limitations: []string{
			"The candidate is an uncommitted scanner-owned copy derived from the recorded Git revision; it is not a clean checkout of that revision.",
		},
	})
	sort.Slice(item.Facts, func(left, right int) bool {
		if item.Facts[left].Key != item.Facts[right].Key {
			return item.Facts[left].Key < item.Facts[right].Key
		}
		if item.Facts[left].Source != item.Facts[right].Source {
			return item.Facts[left].Source < item.Facts[right].Source
		}
		return item.Facts[left].Value < item.Facts[right].Value
	})
	if err := seal(&item); err != nil {
		return model.Inventory{}, err
	}
	return item, nil
}

// BindConfiguration adds reviewed declarations as explicitly sourced facts and
// reseals the inventory. Declarations influence applicability but do not prove
// deployed behavior, and exclusions do not remove files from inventory.
func BindConfiguration(item model.Inventory, validation projectconfig.Validation, sourcePath string) (model.Inventory, error) {
	if item.SchemaVersion != model.InventorySchema || item.Root == "" {
		return model.Inventory{}, fmt.Errorf("configuration requires a current rooted inventory")
	}
	if err := validation.Validate(); err != nil {
		return model.Inventory{}, err
	}
	if err := validation.VerifySource(sourcePath); err != nil {
		return model.Inventory{}, err
	}
	if item.DeclaredScope != nil {
		return model.Inventory{}, fmt.Errorf("inventory already has a bound configuration")
	}
	source := "configuration:external"
	absoluteSource, err := filepath.Abs(sourcePath)
	if err != nil {
		return model.Inventory{}, fmt.Errorf("resolve configuration source: %w", err)
	}
	if resolvedSource, resolveErr := filepath.EvalSymlinks(absoluteSource); resolveErr == nil {
		absoluteSource = resolvedSource
	}
	if relative, relativeErr := filepath.Rel(item.Root, absoluteSource); relativeErr == nil {
		relative = filepath.ToSlash(relative)
		if relative != ".." && !strings.HasPrefix(relative, "../") {
			matched := false
			for _, record := range item.Files {
				if record.Path == relative {
					if record.SHA256 != validation.SourceSHA256 {
						return model.Inventory{}, fmt.Errorf("configuration changed between validation and inventory")
					}
					matched = true
					break
				}
			}
			if !matched {
				return model.Inventory{}, fmt.Errorf("configuration inside target is not an inventoried regular file")
			}
			source = relative
		}
	}
	document := validation.Configuration
	if document.Assessment.SourceRef != "" {
		if item.GitCommit == "" {
			return model.Inventory{}, fmt.Errorf("declared source_ref cannot be verified because the target has no Git revision")
		}
		if document.Assessment.SourceRef != item.GitCommit {
			return model.Inventory{}, fmt.Errorf("declared source_ref does not match the inventoried Git revision")
		}
		if inventoryFactValue(item.Facts, "repository.git_worktree_state") != "clean" {
			return model.Inventory{}, fmt.Errorf("declared source_ref cannot identify the scanned bytes because the Git working tree is not verified clean")
		}
	}
	scope := &model.DeclaredScope{
		ConfigurationDigest: validation.Digest,
		ProjectID:           document.Project.ID, ProjectName: document.Project.Name,
		RiskProfile: document.Project.RiskProfile, ProfileID: document.Assessment.Profile,
		SourceRef:           document.Assessment.SourceRef,
		ArtifactDigests:     append([]string{}, document.Assessment.ArtifactDigests...),
		TargetEnvironments:  append([]string{}, document.Assessment.TargetEnvironments...),
		Features:            map[string]bool{},
		DataClassifications: append([]string{}, document.Data.Classifications...),
		RegulatedData:       append([]string{}, document.Data.Regulated...),
		ProhibitedEvidence:  append([]string{}, document.Data.ProhibitedInEvidence...),
		Components:          []model.DeclaredComponent{}, Exclusions: []model.DeclaredExclusion{},
	}
	for key, value := range document.Features {
		scope.Features[key] = value
	}
	for _, component := range document.Components.Include {
		scope.Components = append(scope.Components, model.DeclaredComponent{Path: component.Path, Type: component.Type})
		id := "declared-" + component.Type + ":" + component.Path
		item.Components = append(item.Components, model.InventoryComponent{ID: id, Kind: "declared-" + component.Type, Path: component.Path})
		item.Relations = append(item.Relations, model.InventoryRelation{From: "repository:.", To: id, Kind: "declares"})
	}
	for _, exclusion := range document.Components.Exclude {
		scope.Exclusions = append(scope.Exclusions, model.DeclaredExclusion{Path: exclusion.Path, Rationale: exclusion.Rationale})
	}
	item.DeclaredScope = scope
	limitations := []string{"A project declaration is reviewed context, not proof of deployed state or behavior."}
	addDeclaredFact := func(key, value, scopePath string) {
		item.Facts = append(item.Facts, model.InventoryFact{
			Key: key, Value: value, Source: source, Detector: "prc.config.declaration",
			DetectorVersion: "0.1", Confidence: 1, ScopePath: scopePath,
			Limitations: append([]string(nil), limitations...),
		})
	}
	addDeclaredFact("project.id", document.Project.ID, ".")
	addDeclaredFact("project.risk_profile", document.Project.RiskProfile, ".")
	addDeclaredFact("assessment.profile", document.Assessment.Profile, ".")
	if document.Assessment.SourceRef != "" {
		addDeclaredFact("assessment.source_ref", document.Assessment.SourceRef, ".")
	}
	for _, value := range document.Assessment.ArtifactDigests {
		addDeclaredFact("assessment.artifact_digest", value, ".")
	}
	for _, value := range document.Assessment.TargetEnvironments {
		addDeclaredFact("assessment.target_environment", value, ".")
	}
	featureKeys := make([]string, 0, len(document.Features))
	for key := range document.Features {
		featureKeys = append(featureKeys, key)
	}
	sort.Strings(featureKeys)
	for _, key := range featureKeys {
		addDeclaredFact("feature."+key, fmt.Sprintf("%t", document.Features[key]), ".")
	}
	for _, component := range document.Components.Include {
		addDeclaredFact("component.include", component.Type, component.Path)
	}
	for _, exclusion := range document.Components.Exclude {
		addDeclaredFact("component.exclude", exclusion.Rationale, exclusion.Path)
	}
	for _, value := range document.Data.Classifications {
		addDeclaredFact("data.classification", value, ".")
	}
	for _, value := range document.Data.Regulated {
		addDeclaredFact("data.regulated", value, ".")
	}
	for _, value := range document.Data.ProhibitedInEvidence {
		addDeclaredFact("evidence.prohibited_class", value, ".")
	}
	sort.Slice(item.Components, func(i, j int) bool { return item.Components[i].ID < item.Components[j].ID })
	sort.Slice(item.Relations, func(i, j int) bool {
		if item.Relations[i].From != item.Relations[j].From {
			return item.Relations[i].From < item.Relations[j].From
		}
		if item.Relations[i].To != item.Relations[j].To {
			return item.Relations[i].To < item.Relations[j].To
		}
		return item.Relations[i].Kind < item.Relations[j].Kind
	})
	sort.Slice(item.Facts, func(i, j int) bool {
		if item.Facts[i].Key != item.Facts[j].Key {
			return item.Facts[i].Key < item.Facts[j].Key
		}
		if item.Facts[i].Source != item.Facts[j].Source {
			return item.Facts[i].Source < item.Facts[j].Source
		}
		return item.Facts[i].Value < item.Facts[j].Value
	})
	if err := seal(&item); err != nil {
		return model.Inventory{}, err
	}
	return item, nil
}

func isContainerFile(name string) bool {
	lower := strings.ToLower(name)
	return lower == "dockerfile" || strings.HasPrefix(lower, "dockerfile.") || strings.HasSuffix(lower, ".dockerfile")
}

func isYAML(name string) bool {
	lower := strings.ToLower(name)
	return strings.HasSuffix(lower, ".yaml") || strings.HasSuffix(lower, ".yml")
}

func inspectFile(path, relative string, capture bool, captureLimit int64) (model.FileRecord, []byte, error) {
	pathInfo, err := os.Lstat(path)
	if err != nil {
		return model.FileRecord{}, nil, fmt.Errorf("inspect %s: %w", relative, err)
	}
	if !pathInfo.Mode().IsRegular() {
		return model.FileRecord{}, nil, fmt.Errorf("file changed to a non-regular entry: %s", relative)
	}
	file, err := os.Open(path)
	if err != nil {
		return model.FileRecord{}, nil, fmt.Errorf("open %s: %w", relative, err)
	}
	defer file.Close()
	openedInfo, err := file.Stat()
	if err != nil {
		return model.FileRecord{}, nil, fmt.Errorf("inspect open file %s: %w", relative, err)
	}
	if !openedInfo.Mode().IsRegular() || !os.SameFile(pathInfo, openedInfo) {
		return model.FileRecord{}, nil, fmt.Errorf("file changed while opening: %s", relative)
	}
	hasher := sha256.New()
	var captured bytes.Buffer
	capturing := capture
	buffer := make([]byte, 32*1024)
	var size int64
	for {
		count, readErr := file.Read(buffer)
		if count > 0 {
			chunk := buffer[:count]
			_, _ = hasher.Write(chunk)
			size += int64(count)
			if size > maxInventoryFileBytes {
				return model.FileRecord{}, nil, fmt.Errorf("file %s exceeded the %d-byte inventory limit while reading", relative, maxInventoryFileBytes)
			}
			if capturing {
				if size <= captureLimit {
					_, _ = captured.Write(chunk)
				} else {
					capturing = false
					captured.Reset()
				}
			}
		}
		if readErr != nil {
			if readErr == io.EOF {
				break
			}
			return model.FileRecord{}, nil, fmt.Errorf("hash %s: %w", relative, readErr)
		}
	}
	finalOpenInfo, err := file.Stat()
	if err != nil {
		return model.FileRecord{}, nil, fmt.Errorf("reinspect open file %s: %w", relative, err)
	}
	finalPathInfo, err := os.Lstat(path)
	if err != nil || !finalPathInfo.Mode().IsRegular() || !os.SameFile(openedInfo, finalOpenInfo) ||
		!os.SameFile(openedInfo, finalPathInfo) || finalOpenInfo.Size() != size ||
		finalOpenInfo.Mode().Perm() != openedInfo.Mode().Perm() ||
		!finalOpenInfo.ModTime().Equal(openedInfo.ModTime()) {
		return model.FileRecord{}, nil, fmt.Errorf("file changed while hashing: %s", relative)
	}
	record := model.FileRecord{
		Path: relative, Size: size, SHA256: hex.EncodeToString(hasher.Sum(nil)), Mode: uint32(openedInfo.Mode().Perm()),
	}
	if !capturing {
		return record, nil, nil
	}
	return record, append([]byte(nil), captured.Bytes()...), nil
}

type gitStateResult struct {
	Value       string
	Confidence  float64
	Limitations []string
}

func gitRepositoryState(root string) (string, error) {
	data, err := runBoundedGit(root, "rev-parse", "--verify", "HEAD")
	if err == nil {
		commit := strings.TrimSpace(string(data))
		if gitCommitPattern.MatchString(commit) {
			return commit, nil
		}
		return "", fmt.Errorf("git returned an invalid HEAD object name")
	}
	if commit := gitCommit(root); commit != "" {
		return commit, fmt.Errorf("git HEAD was read with bounded metadata access, but worktree-aware Git inspection failed")
	}
	return "", nil
}

func gitWorktreeState(root string) gitStateResult {
	data, err := runBoundedGit(root, "-c", "core.fsmonitor=false", "-c", "core.untrackedCache=false",
		"-c", "core.hooksPath="+os.DevNull, "status", "--porcelain=v1", "-z", "--untracked-files=all", "--ignore-submodules=none", "--", ".")
	if err != nil {
		return gitStateResult{Value: "unverified", Confidence: 0, Limitations: []string{
			"A clean Git state could not be proven with the bounded, configuration-isolated status check.",
		}}
	}
	if len(data) == 0 {
		return gitStateResult{Value: "clean", Confidence: 1, Limitations: []string{
			"Git cleanliness binds tracked and untracked worktree state to HEAD; the scanner inventory digest separately binds the bytes that were read.",
		}}
	}
	return gitStateResult{Value: "dirty", Confidence: 1, Limitations: []string{
		"The scanned bytes do not exactly identify the recorded Git HEAD because tracked or untracked worktree changes are present.",
	}}
}

func runBoundedGit(root string, arguments ...string) ([]byte, error) {
	executable, err := exec.LookPath("git")
	if err != nil {
		return nil, fmt.Errorf("find Git executable: %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), gitInspectionTimeout)
	defer cancel()
	args := append([]string{"-C", root}, arguments...)
	command := exec.CommandContext(ctx, executable, args...)
	command.Env = gitEnvironment()
	stdout := &limitedWriter{limit: maxGitStatusBytes}
	stderr := &limitedWriter{limit: 64 * 1024}
	command.Stdout = stdout
	command.Stderr = stderr
	if err := command.Run(); err != nil {
		if ctx.Err() != nil {
			return nil, fmt.Errorf("git inspection timed out: %w", ctx.Err())
		}
		return nil, fmt.Errorf("git inspection failed with diagnostics digest %x: %w", sha256.Sum256(stderr.data), err)
	}
	if stdout.err != nil || stderr.err != nil {
		return nil, fmt.Errorf("git inspection output exceeded its scanner-owned limit")
	}
	return append([]byte{}, stdout.data...), nil
}

type limitedWriter struct {
	data  []byte
	limit int
	err   error
}

func (writer *limitedWriter) Write(data []byte) (int, error) {
	if writer.err != nil {
		return 0, writer.err
	}
	remaining := writer.limit - len(writer.data)
	if len(data) > remaining {
		if remaining > 0 {
			writer.data = append(writer.data, data[:remaining]...)
		}
		writer.err = fmt.Errorf("bounded output limit exceeded")
		return max(remaining, 0), writer.err
	}
	writer.data = append(writer.data, data...)
	return len(data), nil
}

func gitEnvironment() []string {
	result := []string{
		"GIT_CONFIG_NOSYSTEM=1", "GIT_OPTIONAL_LOCKS=0", "GIT_TERMINAL_PROMPT=0",
		"GIT_CONFIG_GLOBAL=" + os.DevNull,
	}
	for _, name := range []string{"PATH", "SystemRoot", "WINDIR", "TMPDIR", "TEMP", "TMP"} {
		if value := os.Getenv(name); value != "" {
			result = append(result, name+"="+value)
		}
	}
	return result
}

func inventoryFactValue(facts []model.InventoryFact, key string) string {
	for _, fact := range facts {
		if fact.Key == key {
			return fact.Value
		}
	}
	return ""
}

func gitCommit(root string) string {
	targetRoot, err := os.OpenRoot(root)
	if err != nil {
		return ""
	}
	defer targetRoot.Close()
	info, err := targetRoot.Lstat(".git")
	if err != nil {
		return ""
	}
	var gitRoot *os.Root
	if info.IsDir() {
		gitRoot, err = targetRoot.OpenRoot(".git")
	} else if info.Mode().IsRegular() {
		data, readErr := readRootRegularFile(targetRoot, ".git", maxGitReferenceBytes)
		if readErr != nil {
			return ""
		}
		line := strings.TrimSpace(string(data))
		if !strings.HasPrefix(line, "gitdir: ") {
			return ""
		}
		gitDirectory := strings.TrimSpace(strings.TrimPrefix(line, "gitdir: "))
		if gitDirectory == "" || filepath.IsAbs(gitDirectory) || strings.Contains(gitDirectory, `\`) {
			return ""
		}
		gitRoot, err = targetRoot.OpenRoot(filepath.FromSlash(gitDirectory))
	} else {
		return ""
	}
	if err != nil {
		return ""
	}
	defer gitRoot.Close()
	head, err := readRootRegularFile(gitRoot, "HEAD", maxGitReferenceBytes)
	if err != nil {
		return ""
	}
	value := strings.TrimSpace(string(head))
	if !strings.HasPrefix(value, "ref: ") {
		if gitCommitPattern.MatchString(value) {
			return value
		}
		return ""
	}
	ref := strings.TrimSpace(strings.TrimPrefix(value, "ref: "))
	if !validGitRef(ref) {
		return ""
	}
	if data, readErr := readRootRegularFile(gitRoot, filepath.FromSlash(ref), maxGitReferenceBytes); readErr == nil {
		commit := strings.TrimSpace(string(data))
		if gitCommitPattern.MatchString(commit) {
			return commit
		}
	}
	packed, err := readRootRegularFile(gitRoot, "packed-refs", maxPackedRefsBytes)
	if err != nil {
		return ""
	}
	scanner := bufio.NewScanner(bytes.NewReader(packed))
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) == 2 && fields[1] == ref && gitCommitPattern.MatchString(fields[0]) {
			return fields[0]
		}
	}
	return ""
}

func validGitRef(ref string) bool {
	if len(ref) < len("refs/x") || len(ref) > 1024 || !strings.HasPrefix(ref, "refs/") ||
		strings.ContainsAny(ref, " ~^:?*[\\") || strings.Contains(ref, "..") ||
		strings.Contains(ref, "//") || strings.Contains(ref, "@{") ||
		strings.HasSuffix(ref, "/") || strings.HasSuffix(ref, ".") || strings.HasSuffix(ref, ".lock") {
		return false
	}
	for _, component := range strings.Split(ref, "/") {
		if component == "" || strings.HasPrefix(component, ".") || strings.HasSuffix(component, ".") {
			return false
		}
	}
	return true
}

func readRootRegularFile(root *os.Root, name string, limit int64) ([]byte, error) {
	info, err := root.Lstat(name)
	if err != nil || !info.Mode().IsRegular() || info.Size() > limit {
		return nil, fmt.Errorf("root-scoped file is missing, unsafe, or exceeds %d bytes", limit)
	}
	file, err := root.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	openedInfo, err := file.Stat()
	if err != nil || !openedInfo.Mode().IsRegular() || !os.SameFile(info, openedInfo) || openedInfo.Size() > limit {
		return nil, fmt.Errorf("root-scoped file changed while opening")
	}
	data, err := io.ReadAll(io.LimitReader(file, limit+1))
	if err != nil {
		return nil, fmt.Errorf("read bounded root-scoped file: %w", err)
	}
	if int64(len(data)) > limit {
		return nil, fmt.Errorf("root-scoped file exceeded %d bytes while reading", limit)
	}
	return data, nil
}
