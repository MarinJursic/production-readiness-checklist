package engine

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	workspaceinventory "github.com/MarinJursic/production-readiness-checklist/scanner/inventory"
	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
	"gopkg.in/yaml.v3"
)

const maximumNativeParseBytes = 4 * 1024 * 1024

const (
	maximumSensitiveMaterialFiles     = 10_000
	maximumSensitiveMaterialBytes     = 256 * 1024 * 1024
	maximumSensitiveMaterialFileBytes = 8 * 1024 * 1024
	maximumSensitiveLocations         = 100
	maximumKubernetesProblems         = 100
)

var (
	errNativeInputLimit = errors.New("native parser input limit exceeded")
	gitRevisionPattern  = regexp.MustCompile(`^[0-9a-f]{40,64}$`)
	fromDigestPattern   = regexp.MustCompile(`@sha256:[0-9a-fA-F]{64}$`)
	conflictStart       = regexp.MustCompile(`(?m)^<<<<<<<(?: |\r?$)`)
	conflictMiddle      = regexp.MustCompile(`(?m)^=======[[:space:]]*$`)
	conflictEnd         = regexp.MustCompile(`(?m)^>>>>>>>(?: |\r?$)`)
)

var sensitiveMaterialExtensions = map[string]bool{
	".asc": true, ".bak": true, ".cfg": true, ".conf": true, ".config": true,
	".css": true, ".env": true, ".gql": true, ".graphql": true, ".gradle": true,
	".hcl": true, ".html": true, ".ini": true, ".json": true, ".jsonl": true,
	".key": true, ".kts": true, ".md": true, ".mk": true, ".pem": true,
	".properties": true, ".proto": true, ".rst": true, ".scss": true, ".sql": true,
	".svelte": true, ".tf": true, ".tfvars": true, ".toml": true, ".txt": true,
	".vue": true, ".xml": true, ".yaml": true, ".yml": true,
}

var sensitiveMaterialNames = map[string]bool{
	"dockerfile": true, "gemfile": true, "jenkinsfile": true, "makefile": true,
	"procfile": true, "rakefile": true, "vagrantfile": true,
	"id_dsa": true, "id_ecdsa": true, "id_ed25519": true, "id_mldsa": true, "id_rsa": true,
}

func inventoryEvidence(
	inventory model.Inventory,
	kind, producer, source, summary string,
	observedAt time.Time,
) model.Evidence {
	evidence := model.Evidence{
		SchemaVersion: model.EvidenceSchema, Kind: kind, Authority: "repository",
		Producer: producer, TargetDigest: inventory.Digest, Source: source,
		ObservedAt: observedAt, Summary: summary,
	}
	identity, _ := json.Marshal(evidence)
	digest := sha256.Sum256(identity)
	evidence.ID = hex.EncodeToString(digest[:])
	return evidence
}

func expectedFile(inventory model.Inventory, relative string) (model.FileRecord, bool) {
	for _, file := range inventory.Files {
		if file.Path == filepath.ToSlash(relative) {
			return file, true
		}
	}
	return model.FileRecord{}, false
}

func readVerifiedEvidence(
	inventory model.Inventory,
	relative, kind, producer, summary string,
	observedAt time.Time,
) ([]byte, model.Evidence, error) {
	return readVerifiedEvidenceWithLimit(inventory, relative, kind, producer, summary, observedAt, maximumNativeParseBytes)
}

func readVerifiedEvidenceWithLimit(
	inventory model.Inventory,
	relative, kind, producer, summary string,
	observedAt time.Time,
	maximumBytes int64,
) ([]byte, model.Evidence, error) {
	if maximumBytes < 1 || maximumBytes > maximumSensitiveMaterialFileBytes {
		return nil, model.Evidence{}, fmt.Errorf("native evidence read limit is invalid")
	}
	record, ok := expectedFile(inventory, relative)
	if !ok {
		return nil, model.Evidence{}, fmt.Errorf("file is absent from inventory: %s", relative)
	}
	if record.Size > maximumBytes {
		return nil, model.Evidence{}, fmt.Errorf("%w: %s has %d bytes", errNativeInputLimit, relative, record.Size)
	}
	path, err := safePath(inventory.Root, relative)
	if err != nil {
		return nil, model.Evidence{}, err
	}
	pathInfo, err := os.Lstat(path)
	if err != nil {
		return nil, model.Evidence{}, fmt.Errorf("inspect %s: %w", relative, err)
	}
	if !pathInfo.Mode().IsRegular() {
		return nil, model.Evidence{}, fmt.Errorf("file changed to a non-regular entry: %s", relative)
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, model.Evidence{}, fmt.Errorf("open %s: %w", relative, err)
	}
	defer file.Close()
	openedInfo, err := file.Stat()
	if err != nil {
		return nil, model.Evidence{}, fmt.Errorf("inspect open file %s: %w", relative, err)
	}
	if !openedInfo.Mode().IsRegular() || !os.SameFile(pathInfo, openedInfo) {
		return nil, model.Evidence{}, fmt.Errorf("file changed while opening: %s", relative)
	}
	data, err := io.ReadAll(io.LimitReader(file, maximumBytes+1))
	if err != nil {
		return nil, model.Evidence{}, fmt.Errorf("read %s: %w", relative, err)
	}
	if int64(len(data)) > maximumBytes {
		return nil, model.Evidence{}, fmt.Errorf("%w: %s", errNativeInputLimit, relative)
	}
	digest := sha256.Sum256(data)
	contentDigest := hex.EncodeToString(digest[:])
	if int64(len(data)) != record.Size || contentDigest != record.SHA256 || uint32(openedInfo.Mode().Perm()) != record.Mode {
		return nil, model.Evidence{}, fmt.Errorf("target changed after inventory: %s", relative)
	}
	evidence := model.Evidence{
		SchemaVersion: model.EvidenceSchema, Kind: kind, Authority: "repository",
		Producer: producer, TargetDigest: inventory.Digest, Source: filepath.ToSlash(relative),
		ContentSHA256: contentDigest, Size: int64(len(data)), ObservedAt: observedAt, Summary: summary,
	}
	identity, err := json.Marshal(evidence)
	if err != nil {
		return nil, model.Evidence{}, err
	}
	evidenceDigest := sha256.Sum256(identity)
	evidence.ID = hex.EncodeToString(evidenceDigest[:])
	return data, evidence, nil
}

func nativeReadFailure(result model.AssertionResult, err error) model.AssertionResult {
	if errors.Is(err, errNativeInputLimit) {
		result.Execution, result.Assessment = "blocked", "unknown"
	} else {
		result.Execution, result.Assessment = "error", "unknown"
	}
	result.Summary = err.Error()
	return result
}

func evaluateGitRevision(
	assertion model.Assertion,
	inventory model.Inventory,
	result model.AssertionResult,
	observedAt time.Time,
) model.AssertionResult {
	if !gitRevisionPattern.MatchString(inventory.GitCommit) {
		result.Assessment = "fail"
		result.Summary = "The target did not resolve to an immutable Git commit identity."
		return result
	}
	result.Assessment = "pass"
	result.Summary = "The inventory identifies source revision " + inventory.GitCommit + "."
	result.EvidenceObserved = []model.Evidence{inventoryEvidence(
		inventory, "repository-identity", assertion.ImplementationID, ".git/HEAD",
		"The inventory resolved the repository HEAD to an immutable commit.", observedAt,
	)}
	return result
}

func decodeWorkflow(data []byte) (*yaml.Node, error) {
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	var document yaml.Node
	if err := decoder.Decode(&document); err != nil {
		return nil, err
	}
	if len(document.Content) == 0 || document.Content[0].Kind != yaml.MappingNode {
		return nil, errors.New("workflow document root is not a mapping")
	}
	var extra yaml.Node
	if err := decoder.Decode(&extra); err != io.EOF {
		if err == nil {
			return nil, errors.New("workflow contains multiple YAML documents")
		}
		return nil, err
	}
	return &document, nil
}

func readWorkflows(
	assertion model.Assertion,
	inventory model.Inventory,
	observedAt time.Time,
) ([]string, []*yaml.Node, []model.Evidence, error) {
	paths := workflowPaths(inventory)
	documents := make([]*yaml.Node, 0, len(paths))
	evidence := make([]model.Evidence, 0, len(paths))
	for _, relative := range paths {
		data, item, err := readVerifiedEvidence(
			inventory, relative, "workflow-parse", assertion.ImplementationID,
			"Parsed a GitHub Actions workflow from the inventoried content.", observedAt,
		)
		if err != nil {
			return nil, nil, nil, err
		}
		document, err := decodeWorkflow(data)
		if err != nil {
			return paths, documents, evidence, fmt.Errorf("cannot parse %s: %w", relative, err)
		}
		documents = append(documents, document)
		evidence = append(evidence, item)
	}
	return paths, documents, evidence, nil
}

func evaluateWorkflowValidity(assertion model.Assertion, inventory model.Inventory, result model.AssertionResult, observedAt time.Time) model.AssertionResult {
	paths, _, evidence, err := readWorkflows(assertion, inventory, observedAt)
	if err != nil {
		return nativeReadFailure(result, err)
	}
	result.Assessment = "pass"
	result.Summary = fmt.Sprintf("All %d detected workflow files parse as one YAML mapping document.", len(paths))
	result.EvidenceObserved = evidence
	return result
}

func evaluateWorkflowJobs(assertion model.Assertion, inventory model.Inventory, result model.AssertionResult, observedAt time.Time) model.AssertionResult {
	paths, documents, evidence, err := readWorkflows(assertion, inventory, observedAt)
	if err != nil {
		return nativeReadFailure(result, err)
	}
	missing := []string{}
	for index, document := range documents {
		jobs := mappingValue(document, "jobs")
		if jobs == nil || jobs.Kind != yaml.MappingNode || len(jobs.Content) == 0 {
			missing = append(missing, paths[index])
		}
	}
	result.EvidenceObserved = evidence
	if len(missing) > 0 {
		result.Assessment = "fail"
		result.Summary = "Workflows without a nonempty jobs mapping: " + strings.Join(missing, ", ") + "."
		return result
	}
	result.Assessment = "pass"
	result.Summary = "Every detected workflow defines at least one job."
	return result
}

func intParameter(assertion model.Assertion, key string, fallback int) int {
	switch value := assertion.Parameters[key].(type) {
	case int:
		return value
	case int64:
		return int(value)
	case float64:
		return int(value)
	default:
		return fallback
	}
}

func scalarPositiveInteger(node *yaml.Node) (int, bool) {
	if node == nil || node.Kind != yaml.ScalarNode || strings.Contains(node.Value, "${{") {
		return 0, false
	}
	value, err := strconv.Atoi(node.Value)
	return value, err == nil && value > 0
}

func evaluateWorkflowTimeouts(assertion model.Assertion, inventory model.Inventory, result model.AssertionResult, observedAt time.Time) model.AssertionResult {
	paths, documents, evidence, err := readWorkflows(assertion, inventory, observedAt)
	if err != nil {
		return nativeReadFailure(result, err)
	}
	maximum := intParameter(assertion, "maximum_minutes", 360)
	violations := []string{}
	for index, document := range documents {
		jobs := mappingValue(document, "jobs")
		if jobs == nil || jobs.Kind != yaml.MappingNode {
			continue
		}
		for item := 0; item+1 < len(jobs.Content); item += 2 {
			name, job := jobs.Content[item].Value, jobs.Content[item+1]
			if mappingValue(job, "runs-on") == nil {
				continue
			}
			value, ok := scalarPositiveInteger(mappingValue(job, "timeout-minutes"))
			if !ok || value > maximum {
				violations = append(violations, paths[index]+":"+name)
			}
		}
	}
	sort.Strings(violations)
	result.EvidenceObserved = evidence
	if len(violations) > 0 {
		result.Assessment = "fail"
		result.Summary = fmt.Sprintf("Jobs without a static timeout from 1 through %d minutes: %s.", maximum, strings.Join(violations, ", "))
		return result
	}
	result.Assessment = "pass"
	result.Summary = fmt.Sprintf("Every directly executed workflow job has a static timeout no greater than %d minutes.", maximum)
	return result
}

func triggerDeclared(node *yaml.Node, trigger string) bool {
	if node == nil {
		return false
	}
	switch node.Kind {
	case yaml.ScalarNode:
		return node.Value == trigger
	case yaml.SequenceNode:
		for _, child := range node.Content {
			if child.Kind == yaml.ScalarNode && child.Value == trigger {
				return true
			}
		}
	case yaml.MappingNode:
		for index := 0; index+1 < len(node.Content); index += 2 {
			if node.Content[index].Value == trigger {
				return true
			}
		}
	}
	return false
}

func evaluateNoPullRequestTarget(assertion model.Assertion, inventory model.Inventory, result model.AssertionResult, observedAt time.Time) model.AssertionResult {
	paths, documents, evidence, err := readWorkflows(assertion, inventory, observedAt)
	if err != nil {
		return nativeReadFailure(result, err)
	}
	violations := []string{}
	for index, document := range documents {
		if triggerDeclared(mappingValue(document, "on"), "pull_request_target") {
			violations = append(violations, paths[index])
		}
	}
	result.EvidenceObserved = evidence
	if len(violations) > 0 {
		result.Assessment = "manual_review"
		result.Summary = "Workflows using pull_request_target require security review of untrusted-code data flow, permissions, and secrets: " + strings.Join(violations, ", ") + ". Trigger presence alone is not treated as proof of a vulnerability."
		return result
	}
	result.Assessment = "pass"
	result.Summary = "No detected workflow uses pull_request_target, so the privileged-trigger review condition is absent."
	return result
}

func mergeConflictScope(inventory model.Inventory) []string {
	selected := map[string]bool{}
	add := func(paths []string) {
		for _, path := range paths {
			selected[path] = true
		}
	}
	for _, file := range inventory.Files {
		if workspaceinventory.IsSourcePath(file.Path) {
			selected[file.Path] = true
		}
	}
	add(inventory.Manifests)
	add(inventory.LockFiles)
	add(inventory.CI.WorkflowFiles)
	add(inventory.ContainerFiles)
	add(inventory.Infrastructure.TerraformFiles)
	add(inventory.Infrastructure.KubernetesFiles)
	paths := make([]string, 0, len(selected))
	for path := range selected {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	return paths
}

func evaluateMergeConflictMarkers(assertion model.Assertion, inventory model.Inventory, result model.AssertionResult, observedAt time.Time) model.AssertionResult {
	paths := mergeConflictScope(inventory)
	if len(paths) > 4096 {
		result.Execution, result.Assessment = "blocked", "unknown"
		result.Summary = fmt.Sprintf("Conflict-marker inspection requires %d files, above the 4096-file limit.", len(paths))
		return result
	}
	violations := []string{}
	for _, relative := range paths {
		data, _, err := readVerifiedEvidence(inventory, relative, "source-integrity", assertion.ImplementationID,
			"Inspected inventoried content for a complete merge-conflict marker set.", observedAt)
		if err != nil {
			return nativeReadFailure(result, err)
		}
		if conflictStart.Match(data) && conflictMiddle.Match(data) && conflictEnd.Match(data) {
			violations = append(violations, relative)
		}
	}
	result.EvidenceObserved = []model.Evidence{inventoryEvidence(
		inventory, "source-integrity", assertion.ImplementationID, "inventory:scoped-content",
		fmt.Sprintf("Inspected %d content-addressed source, build, workflow, container, and infrastructure files.", len(paths)), observedAt,
	)}
	if len(violations) > 0 {
		result.Assessment = "fail"
		result.Summary = "Files containing complete merge-conflict marker sets: " + strings.Join(violations, ", ") + "."
		return result
	}
	result.Assessment = "pass"
	result.Summary = fmt.Sprintf("No complete merge-conflict marker set was found in %d scoped files.", len(paths))
	return result
}

func evaluateRestrictiveFileModes(assertion model.Assertion, inventory model.Inventory, result model.AssertionResult, observedAt time.Time) model.AssertionResult {
	violations := []string{}
	for _, file := range inventory.Files {
		if file.Mode&0o022 != 0 {
			violations = append(violations, fmt.Sprintf("%s(%#o)", file.Path, file.Mode))
		}
	}
	result.EvidenceObserved = []model.Evidence{inventoryEvidence(
		inventory, "repository-metadata", assertion.ImplementationID, "inventory:file-modes",
		fmt.Sprintf("Evaluated permission modes for %d inventoried regular files.", len(inventory.Files)), observedAt,
	)}
	if len(violations) > 0 {
		result.Assessment = "fail"
		result.Summary = "Group- or other-writable files: " + strings.Join(violations, ", ") + "."
		return result
	}
	result.Assessment = "pass"
	result.Summary = "No inventoried regular file is writable by group or other users."
	return result
}

func evaluateInventoryFilesNonempty(assertion model.Assertion, inventory model.Inventory, result model.AssertionResult, observedAt time.Time) model.AssertionResult {
	field, _ := assertion.Parameters["inventory_field"].(string)
	var paths []string
	switch field {
	case "manifests":
		paths = inventory.Manifests
	case "lock_files":
		paths = inventory.LockFiles
	default:
		result.Execution, result.Assessment = "error", "unknown"
		result.Summary = "Unsupported inventory file field: " + field
		return result
	}
	if len(paths) == 0 {
		result.Assessment = "fail"
		result.Summary = "No files were detected for inventory field " + field + "."
		return result
	}
	empty := []string{}
	for _, relative := range paths {
		_, evidence, err := readVerifiedEvidence(inventory, relative, "repository-file", assertion.ImplementationID,
			"Verified an inventoried dependency file is nonempty.", observedAt)
		if err != nil {
			return nativeReadFailure(result, err)
		}
		result.EvidenceObserved = append(result.EvidenceObserved, evidence)
		if evidence.Size == 0 {
			empty = append(empty, relative)
		}
	}
	if len(empty) > 0 {
		result.Assessment = "fail"
		result.Summary = "Empty inventoried dependency files: " + strings.Join(empty, ", ") + "."
		return result
	}
	result.Assessment = "pass"
	result.Summary = fmt.Sprintf("All %d inventoried %s files are nonempty.", len(paths), field)
	return result
}

func versionDeclaration(ecosystem, relative string, data []byte) bool {
	name := filepath.Base(relative)
	text := string(data)
	switch ecosystem {
	case "go":
		return name == "go.mod" && regexp.MustCompile(`(?m)^go[[:space:]]+[0-9]+\.[0-9]+`).MatchString(text)
	case "python":
		if name == ".python-version" || name == "runtime.txt" {
			return len(strings.TrimSpace(text)) > 0
		}
		return (name == "pyproject.toml" && regexp.MustCompile(`(?m)^requires-python[[:space:]]*=`).MatchString(text)) ||
			(strings.Contains(relative, ".github/workflows/") && regexp.MustCompile(`(?m)^[[:space:]]*python-version:[[:space:]]*["']?[^[:space:]"']+`).MatchString(text)) ||
			(name == ".tool-versions" && regexp.MustCompile(`(?m)^python[[:space:]]+\S+`).MatchString(text))
	case "node":
		if name == ".nvmrc" || name == ".node-version" {
			return len(strings.TrimSpace(text)) > 0
		}
		if name == "package.json" {
			var manifest struct {
				Engines map[string]any `json:"engines"`
				Volta   map[string]any `json:"volta"`
			}
			return json.Unmarshal(data, &manifest) == nil && (manifest.Engines["node"] != nil || manifest.Volta["node"] != nil)
		}
		return (strings.Contains(relative, ".github/workflows/") && regexp.MustCompile(`(?m)^[[:space:]]*node-version:[[:space:]]*["']?[^[:space:]"']+`).MatchString(text)) ||
			(name == ".tool-versions" && regexp.MustCompile(`(?m)^nodejs[[:space:]]+\S+`).MatchString(text))
	case "rust":
		return ((name == "rust-toolchain" || name == "rust-toolchain.toml") && len(strings.TrimSpace(text)) > 0) ||
			(name == "Cargo.toml" && regexp.MustCompile(`(?m)^rust-version[[:space:]]*=`).MatchString(text))
	case "java":
		if name == ".java-version" || name == ".sdkmanrc" {
			return len(strings.TrimSpace(text)) > 0
		}
		return regexp.MustCompile(`(?m)(java-version:|maven\.compiler\.(release|source|target)|languageVersion|^java[[:space:]]+\S+)`).MatchString(text)
	}
	return false
}

func evaluateRuntimeVersions(assertion model.Assertion, inventory model.Inventory, result model.AssertionResult, observedAt time.Time) model.AssertionResult {
	missing := []string{}
	for _, ecosystem := range inventory.PackageEcosystems {
		found := false
		for _, file := range inventory.Files {
			if file.Size == 0 || file.Size > maximumNativeParseBytes {
				continue
			}
			name := filepath.Base(file.Path)
			candidate := name == "go.mod" || name == "pyproject.toml" || name == "package.json" || name == "Cargo.toml" ||
				name == "pom.xml" || strings.HasPrefix(name, "build.gradle") || name == ".python-version" || name == "runtime.txt" ||
				name == ".nvmrc" || name == ".node-version" || name == ".tool-versions" || name == "rust-toolchain" ||
				name == "rust-toolchain.toml" || name == ".java-version" || name == ".sdkmanrc" || strings.Contains(file.Path, ".github/workflows/")
			if !candidate {
				continue
			}
			data, evidence, err := readVerifiedEvidence(inventory, file.Path, "runtime-declaration", assertion.ImplementationID,
				"Inspected a supported runtime version declaration.", observedAt)
			if err != nil {
				return nativeReadFailure(result, err)
			}
			if versionDeclaration(ecosystem, file.Path, data) {
				result.EvidenceObserved = append(result.EvidenceObserved, evidence)
				found = true
				break
			}
		}
		if !found {
			missing = append(missing, ecosystem)
		}
	}
	if len(missing) > 0 {
		result.Assessment = "fail"
		result.Summary = "No supported runtime version declaration was found for: " + strings.Join(missing, ", ") + "."
		return result
	}
	result.Assessment = "pass"
	result.Summary = "Every detected package ecosystem has a supported runtime version declaration."
	return result
}

func dockerLogicalLines(data []byte) []string {
	physical := strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n")
	logical := []string{}
	current := ""
	for _, line := range physical {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || (strings.HasPrefix(trimmed, "#") && current == "") {
			continue
		}
		current += " " + strings.TrimSuffix(trimmed, "\\")
		if strings.HasSuffix(trimmed, "\\") {
			continue
		}
		logical = append(logical, strings.TrimSpace(current))
		current = ""
	}
	if strings.TrimSpace(current) != "" {
		logical = append(logical, strings.TrimSpace(current))
	}
	return logical
}

func dockerFrom(line string) (string, bool) {
	fields := strings.Fields(line)
	if len(fields) < 2 || !strings.EqualFold(fields[0], "FROM") {
		return "", false
	}
	index := 1
	for index < len(fields) && strings.HasPrefix(fields[index], "--") {
		index++
	}
	if index >= len(fields) {
		return "", false
	}
	return fields[index], true
}

func evaluateContainerBasePins(assertion model.Assertion, inventory model.Inventory, result model.AssertionResult, observedAt time.Time) model.AssertionResult {
	violations := []string{}
	for _, relative := range inventory.ContainerFiles {
		data, evidence, err := readVerifiedEvidence(inventory, relative, "containerfile-parse", assertion.ImplementationID,
			"Parsed container build-stage base image declarations.", observedAt)
		if err != nil {
			return nativeReadFailure(result, err)
		}
		result.EvidenceObserved = append(result.EvidenceObserved, evidence)
		stages := 0
		for _, line := range dockerLogicalLines(data) {
			base, ok := dockerFrom(line)
			if !ok {
				continue
			}
			stages++
			if base != "scratch" && !fromDigestPattern.MatchString(base) {
				violations = append(violations, relative+":"+base)
			}
		}
		if stages == 0 {
			violations = append(violations, relative+":missing-FROM")
		}
	}
	if len(violations) > 0 {
		result.Assessment = "fail"
		result.Summary = "Container stages without an immutable base image identity: " + strings.Join(violations, ", ") + "."
		return result
	}
	result.Assessment = "pass"
	result.Summary = "Every detected container build stage uses scratch or a sha256-pinned base image."
	return result
}

func evaluateContainerNonRoot(assertion model.Assertion, inventory model.Inventory, result model.AssertionResult, observedAt time.Time) model.AssertionResult {
	violations := []string{}
	for _, relative := range inventory.ContainerFiles {
		data, evidence, err := readVerifiedEvidence(inventory, relative, "containerfile-parse", assertion.ImplementationID,
			"Parsed final-stage container user declarations.", observedAt)
		if err != nil {
			return nativeReadFailure(result, err)
		}
		result.EvidenceObserved = append(result.EvidenceObserved, evidence)
		stages, user := 0, ""
		for _, line := range dockerLogicalLines(data) {
			fields := strings.Fields(line)
			if _, ok := dockerFrom(line); ok {
				stages++
				user = ""
				continue
			}
			if len(fields) >= 2 && strings.EqualFold(fields[0], "USER") {
				user = fields[1]
			}
		}
		lower := strings.ToLower(user)
		if stages == 0 || user == "" || lower == "root" || lower == "0" || strings.HasPrefix(lower, "0:") || strings.Contains(user, "$") {
			violations = append(violations, relative)
		}
	}
	if len(violations) > 0 {
		result.Assessment = "fail"
		result.Summary = "Container files without a static non-root USER in the final stage: " + strings.Join(violations, ", ") + "."
		return result
	}
	result.Assessment = "pass"
	result.Summary = "Every detected container file declares a static non-root USER in its final stage."
	return result
}

func evaluateTerraformLocks(assertion model.Assertion, inventory model.Inventory, result model.AssertionResult, observedAt time.Time) model.AssertionResult {
	directories := map[string]bool{}
	for _, relative := range inventory.Infrastructure.TerraformFiles {
		directories[filepath.ToSlash(filepath.Dir(relative))] = true
	}
	missing := []string{}
	for directory := range directories {
		lock := filepath.ToSlash(filepath.Join(directory, ".terraform.lock.hcl"))
		if directory == "." {
			lock = ".terraform.lock.hcl"
		}
		record, ok := expectedFile(inventory, lock)
		if !ok || record.Size == 0 {
			missing = append(missing, directory)
			continue
		}
		_, evidence, err := readVerifiedEvidence(inventory, lock, "repository-file", assertion.ImplementationID,
			"Verified a nonempty Terraform provider dependency lock file.", observedAt)
		if err != nil {
			return nativeReadFailure(result, err)
		}
		result.EvidenceObserved = append(result.EvidenceObserved, evidence)
	}
	sort.Strings(missing)
	if len(missing) > 0 {
		result.Assessment = "fail"
		result.Summary = "Terraform directories without a nonempty .terraform.lock.hcl: " + strings.Join(missing, ", ") + "."
		return result
	}
	result.Assessment = "pass"
	result.Summary = "Every detected Terraform configuration directory has a nonempty provider dependency lock file."
	return result
}

func mappingPath(node *yaml.Node, keys ...string) *yaml.Node {
	current := node
	for _, key := range keys {
		current = mappingValue(current, key)
		if current == nil {
			return nil
		}
	}
	return current
}

func nodeScalar(node *yaml.Node) string {
	if node == nil || node.Kind != yaml.ScalarNode {
		return ""
	}
	return node.Value
}

func kubernetesPodSpec(document *yaml.Node) (*yaml.Node, string) {
	kind := nodeScalar(mappingValue(document, "kind"))
	switch kind {
	case "Pod":
		return mappingPath(document, "spec"), kind
	case "Deployment", "StatefulSet", "DaemonSet", "ReplicaSet", "Job":
		return mappingPath(document, "spec", "template", "spec"), kind
	case "CronJob":
		return mappingPath(document, "spec", "jobTemplate", "spec", "template", "spec"), kind
	default:
		return nil, kind
	}
}

func kubernetesDocuments(data []byte) ([]*yaml.Node, error) {
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	documents := []*yaml.Node{}
	for {
		var document yaml.Node
		err := decoder.Decode(&document)
		if err == io.EOF {
			return documents, nil
		}
		if err != nil {
			return nil, err
		}
		if len(document.Content) == 0 {
			continue
		}
		if err := validateKubernetesYAMLShape(&document); err != nil {
			return nil, err
		}
		documents = append(documents, &document)
	}
}

func validateKubernetesYAMLShape(root *yaml.Node) error {
	type pendingNode struct {
		node  *yaml.Node
		depth int
	}
	stack := []pendingNode{{node: root}}
	visited := 0
	for len(stack) > 0 {
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		visited++
		if visited > 100_000 || current.depth > 128 {
			return fmt.Errorf("YAML structure exceeds the bounded node or nesting limit")
		}
		if current.node.Kind == yaml.AliasNode {
			return fmt.Errorf("YAML aliases are outside the bounded Kubernetes manifest parser")
		}
		if current.node.Kind == yaml.MappingNode {
			seen := map[string]bool{}
			for index := 0; index+1 < len(current.node.Content); index += 2 {
				key, value := current.node.Content[index], current.node.Content[index+1]
				if key.Kind != yaml.ScalarNode || key.Tag != "!!str" {
					return fmt.Errorf("mapping key at %d:%d must be a scalar string", key.Line, key.Column)
				}
				if seen[key.Value] {
					return fmt.Errorf("duplicate mapping key %s at %d:%d", key.Value, key.Line, key.Column)
				}
				seen[key.Value] = true
				stack = append(stack, pendingNode{node: value, depth: current.depth + 1})
			}
			continue
		}
		for _, child := range current.node.Content {
			stack = append(stack, pendingNode{node: child, depth: current.depth + 1})
		}
	}
	return nil
}

func boolTrue(node *yaml.Node) bool {
	value, valid := yamlBoolean(node)
	return valid && value
}

func yamlBoolean(node *yaml.Node) (bool, bool) {
	if node == nil || node.Kind != yaml.ScalarNode || node.Tag != "!!bool" {
		return false, false
	}
	switch node.Value {
	case "true":
		return true, true
	case "false":
		return false, true
	default:
		return false, false
	}
}

func absentOrFalse(node *yaml.Node) bool {
	if node == nil || (node.Kind == yaml.ScalarNode && node.Tag == "!!null") {
		return true
	}
	value, valid := yamlBoolean(node)
	return valid && !value
}

func userZero(node *yaml.Node) bool {
	return node != nil && node.Kind == yaml.ScalarNode && node.Value == "0"
}

func containerSequences(podSpec *yaml.Node) []*yaml.Node {
	sequences := []*yaml.Node{}
	for _, key := range []string{"initContainers", "containers"} {
		if sequence := mappingValue(podSpec, key); sequence != nil && sequence.Kind == yaml.SequenceNode {
			sequences = append(sequences, sequence)
		}
	}
	return sequences
}

func securityContainerSequences(podSpec *yaml.Node) []*yaml.Node {
	sequences := containerSequences(podSpec)
	if sequence := mappingValue(podSpec, "ephemeralContainers"); sequence != nil && sequence.Kind == yaml.SequenceNode {
		sequences = append(sequences, sequence)
	}
	return sequences
}

func evaluateKubernetesNonRoot(assertion model.Assertion, inventory model.Inventory, result model.AssertionResult, observedAt time.Time) model.AssertionResult {
	violations := []string{}
	workloads := 0
	for _, relative := range inventory.Infrastructure.KubernetesFiles {
		data, evidence, err := readVerifiedEvidence(inventory, relative, "kubernetes-parse", assertion.ImplementationID,
			"Parsed Kubernetes workload pod and container security contexts.", observedAt)
		if err != nil {
			return nativeReadFailure(result, err)
		}
		documents, err := kubernetesDocuments(data)
		if err != nil {
			return nativeReadFailure(result, fmt.Errorf("cannot parse %s: %w", relative, err))
		}
		result.EvidenceObserved = append(result.EvidenceObserved, evidence)
		for index, document := range documents {
			podSpec, kind := kubernetesPodSpec(document)
			if podSpec == nil {
				continue
			}
			workloads++
			podSecurity := mappingValue(podSpec, "securityContext")
			podNonRoot := boolTrue(mappingValue(podSecurity, "runAsNonRoot"))
			podUserZero := userZero(mappingValue(podSecurity, "runAsUser"))
			sequences := securityContainerSequences(podSpec)
			containerCount := 0
			valid := !podUserZero
			for _, sequence := range sequences {
				for _, container := range sequence.Content {
					containerCount++
					security := mappingValue(container, "securityContext")
					if userZero(mappingValue(security, "runAsUser")) || (!podNonRoot && !boolTrue(mappingValue(security, "runAsNonRoot"))) {
						valid = false
					}
				}
			}
			if containerCount == 0 || !valid {
				violations = append(violations, fmt.Sprintf("%s#document-%d(%s)", relative, index+1, kind))
			}
		}
	}
	if len(violations) > 0 {
		result.Assessment = "fail"
		result.Summary = "Kubernetes workloads without an affirmative non-root policy: " + strings.Join(violations, ", ") + "."
		return result
	}
	result.Assessment = "pass"
	result.Summary = fmt.Sprintf("All %d detected Kubernetes workloads affirmatively require non-root containers.", workloads)
	return result
}

type kubernetesProblem struct {
	Path    string
	Line    int
	Column  int
	Message string
}

func appendKubernetesProblem(problems *[]kubernetesProblem, path string, node *yaml.Node, message string) {
	line, column := 1, 1
	if node != nil {
		if node.Line > 0 {
			line = node.Line
		}
		if node.Column > 0 {
			column = node.Column
		}
	}
	*problems = append(*problems, kubernetesProblem{Path: path, Line: line, Column: column, Message: message})
}

func applyKubernetesProblems(result model.AssertionResult, problems []kubernetesProblem, success, kind string) model.AssertionResult {
	if len(problems) == 0 {
		result.Assessment = "pass"
		result.Summary = success
		return result
	}
	reported := problems
	if len(reported) > maximumKubernetesProblems {
		reported = reported[:maximumKubernetesProblems]
	}
	messages := make([]string, 0, len(reported))
	locations := map[string]bool{}
	for _, problem := range reported {
		messages = append(messages, fmt.Sprintf("%s:%d:%d (%s)", problem.Path, problem.Line, problem.Column, problem.Message))
		key := fmt.Sprintf("%s:%d:%d", problem.Path, problem.Line, problem.Column)
		if !locations[key] {
			locations[key] = true
			result.Locations = append(result.Locations, model.FindingLocation{
				Path: problem.Path, Line: problem.Line, Column: problem.Column,
			})
		}
	}
	result.Assessment = "fail"
	result.Summary = fmt.Sprintf("Detected %d Kubernetes %s problem(s); reporting up to %d: %s.",
		len(problems), kind, maximumKubernetesProblems, strings.Join(messages, ", "))
	return result
}

func kubernetesWorkloadLabel(path string, documentIndex int, kind string) string {
	return fmt.Sprintf("%s#document-%d(%s)", path, documentIndex+1, kind)
}

func nodeOr(node, fallback *yaml.Node) *yaml.Node {
	if node != nil {
		return node
	}
	return fallback
}

func allSecurityContainers(podSpec *yaml.Node) []*yaml.Node {
	containers := make([]*yaml.Node, 0)
	for _, sequence := range securityContainerSequences(podSpec) {
		containers = append(containers, sequence.Content...)
	}
	return containers
}

func windowsWorkload(podSpec *yaml.Node) bool {
	return nodeScalar(mappingPath(podSpec, "os", "name")) == "windows"
}

func evaluateKubernetesHostAccess(assertion model.Assertion, inventory model.Inventory, result model.AssertionResult, observedAt time.Time) model.AssertionResult {
	problems := make([]kubernetesProblem, 0)
	workloads := 0
	for _, relative := range inventory.Infrastructure.KubernetesFiles {
		data, evidence, err := readVerifiedEvidence(inventory, relative, "kubernetes-parse", assertion.ImplementationID,
			"Parsed Kubernetes workload host namespace, volume, and privileged-container fields.", observedAt)
		if err != nil {
			return nativeReadFailure(result, err)
		}
		documents, err := kubernetesDocuments(data)
		if err != nil {
			return nativeReadFailure(result, fmt.Errorf("cannot parse %s: %w", relative, err))
		}
		result.EvidenceObserved = append(result.EvidenceObserved, evidence)
		for documentIndex, document := range documents {
			podSpec, kind := kubernetesPodSpec(document)
			if podSpec == nil {
				continue
			}
			workloads++
			label := kubernetesWorkloadLabel(relative, documentIndex, kind)
			for _, field := range []string{"hostNetwork", "hostPID", "hostIPC"} {
				node := mappingValue(podSpec, field)
				if !absentOrFalse(node) {
					appendKubernetesProblem(&problems, relative, node, label+" must leave "+field+" unset or false")
				}
			}
			podHostProcess := mappingPath(podSpec, "securityContext", "windowsOptions", "hostProcess")
			if !absentOrFalse(podHostProcess) {
				appendKubernetesProblem(&problems, relative, podHostProcess, label+" must leave pod hostProcess unset or false")
			}
			volumes := mappingValue(podSpec, "volumes")
			if volumes != nil && volumes.Kind == yaml.ScalarNode && volumes.Tag == "!!null" {
				volumes = nil
			}
			if volumes != nil && volumes.Kind != yaml.SequenceNode {
				appendKubernetesProblem(&problems, relative, volumes, label+" volumes must be a list")
			} else if volumes != nil {
				for _, volume := range volumes.Content {
					if hostPath := mappingValue(volume, "hostPath"); hostPath != nil {
						appendKubernetesProblem(&problems, relative, hostPath, label+" must not declare a hostPath volume")
					}
				}
			}
			containers := allSecurityContainers(podSpec)
			if len(containers) == 0 {
				appendKubernetesProblem(&problems, relative, podSpec, label+" has no recognized containers")
			}
			for _, container := range containers {
				name := nodeScalar(mappingValue(container, "name"))
				security := mappingValue(container, "securityContext")
				privileged := mappingValue(security, "privileged")
				if !absentOrFalse(privileged) {
					appendKubernetesProblem(&problems, relative, privileged, label+" container "+name+" must leave privileged unset or false")
				}
				hostProcess := mappingPath(security, "windowsOptions", "hostProcess")
				if !absentOrFalse(hostProcess) {
					appendKubernetesProblem(&problems, relative, hostProcess, label+" container "+name+" must leave hostProcess unset or false")
				}
			}
		}
	}
	return applyKubernetesProblems(result, problems,
		fmt.Sprintf("All %d detected Kubernetes workloads avoid host namespaces, hostPath volumes, HostProcess, and privileged containers.", workloads),
		"privileged-host-access")
}

func evaluateKubernetesPrivilegeEscalation(assertion model.Assertion, inventory model.Inventory, result model.AssertionResult, observedAt time.Time) model.AssertionResult {
	problems := make([]kubernetesProblem, 0)
	workloads, containers := 0, 0
	for _, relative := range inventory.Infrastructure.KubernetesFiles {
		data, evidence, err := readVerifiedEvidence(inventory, relative, "kubernetes-parse", assertion.ImplementationID,
			"Parsed Kubernetes Linux container privilege-escalation policy.", observedAt)
		if err != nil {
			return nativeReadFailure(result, err)
		}
		documents, err := kubernetesDocuments(data)
		if err != nil {
			return nativeReadFailure(result, fmt.Errorf("cannot parse %s: %w", relative, err))
		}
		result.EvidenceObserved = append(result.EvidenceObserved, evidence)
		for documentIndex, document := range documents {
			podSpec, kind := kubernetesPodSpec(document)
			if podSpec == nil || windowsWorkload(podSpec) {
				continue
			}
			workloads++
			label := kubernetesWorkloadLabel(relative, documentIndex, kind)
			workloadContainers := allSecurityContainers(podSpec)
			if len(workloadContainers) == 0 {
				appendKubernetesProblem(&problems, relative, podSpec, label+" has no recognized containers")
			}
			for _, container := range workloadContainers {
				containers++
				name := nodeScalar(mappingValue(container, "name"))
				node := mappingPath(container, "securityContext", "allowPrivilegeEscalation")
				value, valid := yamlBoolean(node)
				if !valid || value {
					appendKubernetesProblem(&problems, relative, nodeOr(node, container), label+" container "+name+" must set allowPrivilegeEscalation to false")
				}
			}
		}
	}
	return applyKubernetesProblems(result, problems,
		fmt.Sprintf("All %d containers in %d detected Linux Kubernetes workloads disable privilege escalation.", containers, workloads),
		"privilege-escalation")
}

func stringSequence(node *yaml.Node) ([]string, bool) {
	if node == nil || node.Kind != yaml.SequenceNode {
		return nil, false
	}
	values := make([]string, 0, len(node.Content))
	for _, item := range node.Content {
		if item.Kind != yaml.ScalarNode || item.Tag != "!!str" {
			return nil, false
		}
		values = append(values, item.Value)
	}
	return values, true
}

func evaluateKubernetesCapabilities(assertion model.Assertion, inventory model.Inventory, result model.AssertionResult, observedAt time.Time) model.AssertionResult {
	problems := make([]kubernetesProblem, 0)
	workloads, containers := 0, 0
	for _, relative := range inventory.Infrastructure.KubernetesFiles {
		data, evidence, err := readVerifiedEvidence(inventory, relative, "kubernetes-parse", assertion.ImplementationID,
			"Parsed Kubernetes Linux container capability policy.", observedAt)
		if err != nil {
			return nativeReadFailure(result, err)
		}
		documents, err := kubernetesDocuments(data)
		if err != nil {
			return nativeReadFailure(result, fmt.Errorf("cannot parse %s: %w", relative, err))
		}
		result.EvidenceObserved = append(result.EvidenceObserved, evidence)
		for documentIndex, document := range documents {
			podSpec, kind := kubernetesPodSpec(document)
			if podSpec == nil || windowsWorkload(podSpec) {
				continue
			}
			workloads++
			label := kubernetesWorkloadLabel(relative, documentIndex, kind)
			workloadContainers := allSecurityContainers(podSpec)
			if len(workloadContainers) == 0 {
				appendKubernetesProblem(&problems, relative, podSpec, label+" has no recognized containers")
			}
			for _, container := range workloadContainers {
				containers++
				name := nodeScalar(mappingValue(container, "name"))
				capabilities := mappingPath(container, "securityContext", "capabilities")
				dropNode := mappingValue(capabilities, "drop")
				dropped, valid := stringSequence(dropNode)
				dropsAll := false
				for _, capability := range dropped {
					if capability == "ALL" {
						dropsAll = true
					}
				}
				if !valid || !dropsAll {
					appendKubernetesProblem(&problems, relative, nodeOr(dropNode, container), label+" container "+name+" must drop the ALL capability")
				}
				addNode := mappingValue(capabilities, "add")
				if addNode != nil && !(addNode.Kind == yaml.ScalarNode && addNode.Tag == "!!null") {
					added, addValid := stringSequence(addNode)
					if !addValid {
						appendKubernetesProblem(&problems, relative, addNode, label+" container "+name+" capabilities.add must be a string list")
					} else {
						for _, capability := range added {
							if capability != "NET_BIND_SERVICE" {
								appendKubernetesProblem(&problems, relative, addNode, label+" container "+name+" adds disallowed capability "+capability)
							}
						}
					}
				}
			}
		}
	}
	return applyKubernetesProblems(result, problems,
		fmt.Sprintf("All %d containers in %d detected Linux Kubernetes workloads drop ALL capabilities and add back at most NET_BIND_SERVICE.", containers, workloads),
		"capability-policy")
}

func allowedSeccompType(node *yaml.Node) bool {
	return node != nil && node.Kind == yaml.ScalarNode && node.Tag == "!!str" &&
		(node.Value == "RuntimeDefault" || node.Value == "Localhost")
}

func evaluateKubernetesSeccomp(assertion model.Assertion, inventory model.Inventory, result model.AssertionResult, observedAt time.Time) model.AssertionResult {
	problems := make([]kubernetesProblem, 0)
	workloads, containers := 0, 0
	for _, relative := range inventory.Infrastructure.KubernetesFiles {
		data, evidence, err := readVerifiedEvidence(inventory, relative, "kubernetes-parse", assertion.ImplementationID,
			"Parsed Kubernetes Linux pod and container seccomp profiles.", observedAt)
		if err != nil {
			return nativeReadFailure(result, err)
		}
		documents, err := kubernetesDocuments(data)
		if err != nil {
			return nativeReadFailure(result, fmt.Errorf("cannot parse %s: %w", relative, err))
		}
		result.EvidenceObserved = append(result.EvidenceObserved, evidence)
		for documentIndex, document := range documents {
			podSpec, kind := kubernetesPodSpec(document)
			if podSpec == nil || windowsWorkload(podSpec) {
				continue
			}
			workloads++
			label := kubernetesWorkloadLabel(relative, documentIndex, kind)
			podType := mappingPath(podSpec, "securityContext", "seccompProfile", "type")
			podAllowed := allowedSeccompType(podType)
			workloadContainers := allSecurityContainers(podSpec)
			if len(workloadContainers) == 0 {
				appendKubernetesProblem(&problems, relative, podSpec, label+" has no recognized containers")
			}
			for _, container := range workloadContainers {
				containers++
				name := nodeScalar(mappingValue(container, "name"))
				containerType := mappingPath(container, "securityContext", "seccompProfile", "type")
				if containerType != nil && !allowedSeccompType(containerType) {
					appendKubernetesProblem(&problems, relative, containerType, label+" container "+name+" selects a disallowed seccomp profile")
				} else if containerType == nil && !podAllowed {
					appendKubernetesProblem(&problems, relative, nodeOr(podType, container), label+" container "+name+" has no RuntimeDefault or Localhost seccomp profile")
				}
			}
		}
	}
	return applyKubernetesProblems(result, problems,
		fmt.Sprintf("All %d containers in %d detected Linux Kubernetes workloads inherit or select RuntimeDefault or Localhost seccomp.", containers, workloads),
		"seccomp-policy")
}

func nonemptyMapping(node *yaml.Node) bool {
	return node != nil && node.Kind == yaml.MappingNode && len(node.Content) > 0
}

func evaluateKubernetesResources(assertion model.Assertion, inventory model.Inventory, result model.AssertionResult, observedAt time.Time) model.AssertionResult {
	violations := []string{}
	workloads := 0
	for _, relative := range inventory.Infrastructure.KubernetesFiles {
		data, evidence, err := readVerifiedEvidence(inventory, relative, "kubernetes-parse", assertion.ImplementationID,
			"Parsed Kubernetes container resource requests and limits.", observedAt)
		if err != nil {
			return nativeReadFailure(result, err)
		}
		documents, err := kubernetesDocuments(data)
		if err != nil {
			return nativeReadFailure(result, fmt.Errorf("cannot parse %s: %w", relative, err))
		}
		result.EvidenceObserved = append(result.EvidenceObserved, evidence)
		for index, document := range documents {
			podSpec, kind := kubernetesPodSpec(document)
			if podSpec == nil {
				continue
			}
			workloads++
			containerCount, valid := 0, true
			for _, sequence := range containerSequences(podSpec) {
				for _, container := range sequence.Content {
					containerCount++
					resources := mappingValue(container, "resources")
					if !nonemptyMapping(mappingValue(resources, "requests")) || !nonemptyMapping(mappingValue(resources, "limits")) {
						valid = false
					}
				}
			}
			if containerCount == 0 || !valid {
				violations = append(violations, fmt.Sprintf("%s#document-%d(%s)", relative, index+1, kind))
			}
		}
	}
	if len(violations) > 0 {
		result.Assessment = "fail"
		result.Summary = "Kubernetes workloads with missing container requests or limits: " + strings.Join(violations, ", ") + "."
		return result
	}
	result.Assessment = "pass"
	result.Summary = fmt.Sprintf("All containers in %d detected Kubernetes workloads declare resource requests and limits.", workloads)
	return result
}

func privateKeyArmorLabels() []string {
	return []string{
		"PRIVATE KEY", "ENCRYPTED PRIVATE KEY", "RSA PRIVATE KEY", "EC PRIVATE KEY",
		"DSA PRIVATE KEY", "OPENSSH PRIVATE KEY", "PGP PRIVATE KEY BLOCK",
	}
}

func isSensitiveMaterialPath(path string) bool {
	base := strings.ToLower(filepath.Base(path))
	if workspaceinventory.IsSourcePath(path) || sensitiveMaterialExtensions[strings.ToLower(filepath.Ext(base))] ||
		sensitiveMaterialNames[base] || strings.HasPrefix(base, ".env") ||
		(strings.HasPrefix(base, "ssh_host_") && strings.HasSuffix(base, "_key")) {
		return true
	}
	// Repository dotfiles are commonly text configuration files. Inventory
	// already excludes Git metadata and dependency/cache directories.
	return strings.HasPrefix(base, ".") && base != ".ds_store"
}

func armorPayloadLine(line []byte) (int, bool) {
	line = bytes.TrimSpace(line)
	if len(line) == 0 || bytes.Contains(line, []byte(":")) || line[0] == '=' {
		return 0, true
	}
	for _, character := range line {
		if !(character >= 'A' && character <= 'Z') && !(character >= 'a' && character <= 'z') &&
			!(character >= '0' && character <= '9') && character != '+' && character != '/' && character != '=' {
			return 0, false
		}
	}
	return len(line), true
}

func containsPrivateKeyArmor(data []byte, labels []string) bool {
	lines := bytes.Split(data, []byte("\n"))
	for _, label := range labels {
		begin := []byte("-----BEGIN " + label + "-----")
		end := []byte("-----END " + label + "-----")
		for index, line := range lines {
			if !bytes.Equal(bytes.TrimSpace(line), begin) {
				continue
			}
			payloadBytes := 0
			valid := true
			for cursor := index + 1; cursor < len(lines); cursor++ {
				if bytes.Equal(bytes.TrimSpace(lines[cursor]), end) {
					if valid && payloadBytes >= 32 {
						return true
					}
					break
				}
				count, lineValid := armorPayloadLine(lines[cursor])
				if !lineValid {
					valid = false
				}
				payloadBytes += count
			}
		}
	}
	return false
}

func evaluatePrivateKeyArmor(
	assertion model.Assertion,
	inventory model.Inventory,
	result model.AssertionResult,
	observedAt time.Time,
) model.AssertionResult {
	files := make([]model.FileRecord, 0)
	var totalBytes int64
	for _, file := range inventory.Files {
		if !isSensitiveMaterialPath(file.Path) {
			continue
		}
		files = append(files, file)
		if file.Size < 0 || totalBytes > maximumSensitiveMaterialBytes-file.Size {
			result.Execution, result.Assessment = "blocked", "unknown"
			result.Summary = "Private-key armor inspection could not establish a bounded candidate-file size."
			return result
		}
		totalBytes += file.Size
	}
	if len(files) > maximumSensitiveMaterialFiles {
		result.Execution, result.Assessment = "blocked", "unknown"
		result.Summary = fmt.Sprintf(
			"Private-key armor inspection requires %d candidate files, above the %d-file limit.",
			len(files), maximumSensitiveMaterialFiles,
		)
		return result
	}
	if totalBytes > maximumSensitiveMaterialBytes {
		result.Execution, result.Assessment = "blocked", "unknown"
		result.Summary = fmt.Sprintf(
			"Private-key armor inspection requires %d bytes, above the %d-byte limit.",
			totalBytes, maximumSensitiveMaterialBytes,
		)
		return result
	}
	for _, file := range files {
		if file.Size > maximumSensitiveMaterialFileBytes {
			result.Execution, result.Assessment = "blocked", "unknown"
			result.Summary = fmt.Sprintf(
				"Private-key armor inspection cannot completely inspect %s because its %d bytes exceed the %d-byte per-file limit.",
				file.Path, file.Size, maximumSensitiveMaterialFileBytes,
			)
			return result
		}
	}

	labels := privateKeyArmorLabels()
	violations := make([]string, 0)
	violationEvidence := make([]model.Evidence, 0)
	totalViolations := 0
	for _, file := range files {
		data, evidence, err := readVerifiedEvidenceWithLimit(
			inventory, file.Path, "sensitive-material-scan", assertion.ImplementationID,
			"Inspected content for recognized private-key armor blocks without retaining matched material.", observedAt,
			maximumSensitiveMaterialFileBytes,
		)
		if err != nil {
			return nativeReadFailure(result, err)
		}
		if !containsPrivateKeyArmor(data, labels) {
			continue
		}
		totalViolations++
		if len(violations) < maximumSensitiveLocations {
			violations = append(violations, file.Path)
			violationEvidence = append(violationEvidence, evidence)
		}
	}

	scopeEvidence := inventoryEvidence(
		inventory, "sensitive-material-scan", assertion.ImplementationID, "inventory:sensitive-material-scope",
		fmt.Sprintf("Inspected %d content-addressed candidate files totaling %d bytes; binary and unrecognized file types and Git history are outside this native check.", len(files), totalBytes),
		observedAt,
	)
	result.EvidenceObserved = append([]model.Evidence{scopeEvidence}, violationEvidence...)
	if totalViolations > 0 {
		result.Assessment = "fail"
		result.Summary = fmt.Sprintf("Detected recognized private-key armor in %d inventoried file(s); reporting up to %d paths: %s.", totalViolations, maximumSensitiveLocations, strings.Join(violations, ", "))
		return result
	}
	result.Assessment = "pass"
	result.Summary = fmt.Sprintf("No recognized private-key armor was found in %d bounded candidate files.", len(files))
	return result
}
