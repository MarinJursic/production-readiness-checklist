package inventory

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	projectconfig "github.com/MarinJursic/production-readiness-checklist/scanner/config"
	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
)

const maxEntries = 200_000

const detectorVersion = "0.3"

var (
	kubernetesAPIVersion = regexp.MustCompile(`(?m)^apiVersion:[[:space:]]*[^[:space:]#]+`)
	kubernetesKind       = regexp.MustCompile(`(?m)^kind:[[:space:]]*[^[:space:]#]+`)
)

var excludedDirectories = map[string]bool{
	".git": true, ".mypy_cache": true, ".prc": true, ".pytest_cache": true,
	".venv": true, "__pycache__": true, "node_modules": true, "site": true,
	"vendor": true,
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
	err = filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return fmt.Errorf("walk %s: %w", path, walkErr)
		}
		if path == root {
			return nil
		}
		if entry.IsDir() && excludedDirectories[entry.Name()] {
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
		if len(result.Files) >= maxEntries {
			return fmt.Errorf("target exceeds %d files", maxEntries)
		}
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		relative = filepath.ToSlash(relative)
		fileRecord, err := hashFile(path, relative)
		if err != nil {
			return err
		}
		result.Files = append(result.Files, fileRecord)
		name := entry.Name()
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
		if isYAML(name) && entryInfo.Size() <= 2*1024*1024 {
			data, readErr := os.ReadFile(path)
			if readErr != nil {
				return fmt.Errorf("read inventory candidate %s: %w", relative, readErr)
			}
			if kubernetesAPIVersion.Match(data) && kubernetesKind.Match(data) {
				result.Infrastructure.KubernetesFiles = append(result.Infrastructure.KubernetesFiles, relative)
				addComponent("infrastructure", relative, "kubernetes")
				addFact("infrastructure.kubernetes", "true", relative, "prc.inventory.kubernetes-yaml", 0.9,
					"Top-level Kubernetes fields do not prove that the resource is valid or deployed.")
			}
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
	result.GitCommit = gitCommit(root)

	if err := seal(&result); err != nil {
		return model.Inventory{}, err
	}
	return result, nil
}

func seal(result *model.Inventory) error {
	payload, err := json.Marshal(digestInput{
		GitCommit: result.GitCommit, Files: result.Files, PackageEcosystems: result.PackageEcosystems,
		Manifests: result.Manifests, LockFiles: result.LockFiles, ContainerFiles: result.ContainerFiles,
		Symlinks: result.Symlinks, CI: result.CI, Infrastructure: result.Infrastructure,
		Components: result.Components, Relations: result.Relations, Facts: result.Facts,
		DeclaredScope: result.DeclaredScope,
	})
	if err != nil {
		return fmt.Errorf("encode inventory digest: %w", err)
	}
	digest := sha256.Sum256(payload)
	result.Digest = hex.EncodeToString(digest[:])
	return nil
}

// BindConfiguration adds reviewed declarations as explicitly sourced facts and
// reseals the inventory. Declarations influence applicability but do not prove
// deployed behavior, and exclusions do not remove files from inventory.
func BindConfiguration(item model.Inventory, validation projectconfig.Validation, sourcePath string) (model.Inventory, error) {
	if item.SchemaVersion != model.InventorySchema || item.Root == "" {
		return model.Inventory{}, fmt.Errorf("configuration requires a current rooted inventory")
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

func hashFile(path, relative string) (model.FileRecord, error) {
	pathInfo, err := os.Lstat(path)
	if err != nil {
		return model.FileRecord{}, fmt.Errorf("inspect %s: %w", relative, err)
	}
	if !pathInfo.Mode().IsRegular() {
		return model.FileRecord{}, fmt.Errorf("file changed to a non-regular entry: %s", relative)
	}
	file, err := os.Open(path)
	if err != nil {
		return model.FileRecord{}, fmt.Errorf("open %s: %w", relative, err)
	}
	defer file.Close()
	openedInfo, err := file.Stat()
	if err != nil {
		return model.FileRecord{}, fmt.Errorf("inspect open file %s: %w", relative, err)
	}
	if !openedInfo.Mode().IsRegular() || !os.SameFile(pathInfo, openedInfo) {
		return model.FileRecord{}, fmt.Errorf("file changed while opening: %s", relative)
	}
	hasher := sha256.New()
	size, err := io.Copy(hasher, file)
	if err != nil {
		return model.FileRecord{}, fmt.Errorf("hash %s: %w", relative, err)
	}
	return model.FileRecord{
		Path: relative, Size: size, SHA256: hex.EncodeToString(hasher.Sum(nil)), Mode: uint32(openedInfo.Mode().Perm()),
	}, nil
}

func gitCommit(root string) string {
	gitDirectory := filepath.Join(root, ".git")
	info, err := os.Stat(gitDirectory)
	if err != nil {
		return ""
	}
	if !info.IsDir() {
		data, err := os.ReadFile(gitDirectory)
		if err != nil {
			return ""
		}
		line := strings.TrimSpace(string(data))
		if !strings.HasPrefix(line, "gitdir: ") {
			return ""
		}
		gitDirectory = strings.TrimPrefix(line, "gitdir: ")
		if !filepath.IsAbs(gitDirectory) {
			gitDirectory = filepath.Join(root, gitDirectory)
		}
	}
	head, err := os.ReadFile(filepath.Join(gitDirectory, "HEAD"))
	if err != nil {
		return ""
	}
	value := strings.TrimSpace(string(head))
	if !strings.HasPrefix(value, "ref: ") {
		return value
	}
	ref := strings.TrimPrefix(value, "ref: ")
	if data, err := os.ReadFile(filepath.Join(gitDirectory, filepath.FromSlash(ref))); err == nil {
		return strings.TrimSpace(string(data))
	}
	file, err := os.Open(filepath.Join(gitDirectory, "packed-refs"))
	if err != nil {
		return ""
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) == 2 && fields[1] == ref {
			return fields[0]
		}
	}
	return ""
}
