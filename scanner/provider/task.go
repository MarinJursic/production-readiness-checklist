package provider

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strings"
	"unicode/utf8"

	workspaceinventory "github.com/MarinJursic/production-readiness-checklist/scanner/inventory"
	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
)

var (
	hexDigest = regexp.MustCompile(`^[0-9a-f]{64}$`)
	assertion = regexp.MustCompile(`^PRC-A-[A-Z0-9]+-[0-9]{3}$`)
	control   = regexp.MustCompile(`^[A-Z]+-[0-9A-F]{8}$`)
)

func LoadTask(path string) (Task, error) {
	task, err := readTask(path)
	if err != nil {
		return Task{}, err
	}
	if err := task.Validate(); err != nil {
		return Task{}, err
	}
	return task, nil
}

// SealTask reads a draft task, binds its workspace inventory, replaces its
// task_id with its canonical content digest, and returns it without changing the
// draft file.
func SealTask(path, workspace string) (Task, error) {
	item, err := workspaceinventory.Build(workspace)
	if err != nil {
		return Task{}, err
	}
	return SealTaskWithInventory(path, workspace, item, nil)
}

// SealTaskWithInventory binds a previously constructed current inventory and
// mandatory protected paths into a task. This lets configured callers use the
// same declared-scope identity that remediation will independently verify.
func SealTaskWithInventory(path, workspace string, item model.Inventory, requiredProtectedPaths []string) (Task, error) {
	task, err := readTask(path)
	if err != nil {
		return Task{}, err
	}
	return SealTaskValueWithInventory(task, workspace, item, requiredProtectedPaths)
}

// SealTaskValueWithInventory seals a scanner-constructed draft without first
// persisting it. It applies the same inventory, input, and protected-path
// checks as the file-based task workflow.
func SealTaskValueWithInventory(task Task, workspace string, item model.Inventory, requiredProtectedPaths []string) (Task, error) {
	task.TaskID = ""
	task.Inputs = []InputFile{}
	if err := validatePaths(task.RelevantPaths, "relevant"); err != nil {
		return Task{}, err
	}
	if !slices.Equal(task.RelevantPaths, sortedCopy(task.RelevantPaths)) {
		return Task{}, fmt.Errorf("agent task relevant paths must be sorted")
	}
	if err := validatePaths(task.ProtectedPaths, "protected"); err != nil {
		return Task{}, err
	}
	if !slices.Equal(task.ProtectedPaths, sortedCopy(task.ProtectedPaths)) {
		return Task{}, fmt.Errorf("agent task protected paths must be sorted")
	}
	task.ProtectedPaths = mergedSortedUnique(task.ProtectedPaths, requiredProtectedPaths)
	workspaceRoot, err := filepath.Abs(workspace)
	if err != nil {
		return Task{}, fmt.Errorf("resolve agent workspace: %w", err)
	}
	if resolved, resolveErr := filepath.EvalSymlinks(workspaceRoot); resolveErr == nil {
		workspaceRoot = resolved
	}
	if item.SchemaVersion != model.InventorySchema || item.Root != workspaceRoot || item.Digest == "" {
		return Task{}, fmt.Errorf("agent task requires the current rooted workspace inventory")
	}
	task.WorkspaceInventoryDigest = item.Digest
	records := make(map[string]string, len(item.Files))
	for _, record := range item.Files {
		records[record.Path] = record.SHA256
	}
	for _, relevantPath := range task.RelevantPaths {
		expectedDigest, ok := records[relevantPath]
		if !ok {
			return Task{}, fmt.Errorf("relevant path %s is not an inventoried regular file", relevantPath)
		}
		content, err := readInputFile(workspace, relevantPath, expectedDigest)
		if err != nil {
			return Task{}, err
		}
		task.Inputs = append(task.Inputs, InputFile{Path: relevantPath, SHA256: expectedDigest, Content: content})
	}
	task.TaskID, err = TaskID(task)
	if err != nil {
		return Task{}, err
	}
	if err := task.Validate(); err != nil {
		return Task{}, err
	}
	return task, nil
}

func mergedSortedUnique(groups ...[]string) []string {
	seen := map[string]bool{}
	for _, group := range groups {
		for _, value := range group {
			seen[value] = true
		}
	}
	values := make([]string, 0, len(seen))
	for value := range seen {
		values = append(values, value)
	}
	sort.Strings(values)
	return values
}

func readInputFile(workspace, relativePath, expectedDigest string) (string, error) {
	path := filepath.Join(workspace, filepath.FromSlash(relativePath))
	before, err := os.Lstat(path)
	if err != nil || !before.Mode().IsRegular() || before.Size() > 256*1024 {
		return "", fmt.Errorf("agent input %s must be a regular text file no larger than 256 KiB", relativePath)
	}
	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("open agent input %s: %w", relativePath, err)
	}
	defer file.Close()
	after, err := file.Stat()
	if err != nil || !after.Mode().IsRegular() || !os.SameFile(before, after) {
		return "", fmt.Errorf("agent input changed while opening: %s", relativePath)
	}
	data, err := io.ReadAll(io.LimitReader(file, 256*1024+1))
	if err != nil || len(data) > 256*1024 || !utf8.Valid(data) || bytes.IndexByte(data, 0) >= 0 {
		return "", fmt.Errorf("agent input %s must be valid UTF-8 text no larger than 256 KiB", relativePath)
	}
	digest := sha256.Sum256(data)
	if hex.EncodeToString(digest[:]) != expectedDigest {
		return "", fmt.Errorf("agent input changed after inventory: %s", relativePath)
	}
	return string(data), nil
}

func readTask(path string) (Task, error) {
	info, err := os.Stat(path)
	if err != nil {
		return Task{}, fmt.Errorf("inspect agent task: %w", err)
	}
	if !info.Mode().IsRegular() || info.Size() > 1024*1024 {
		return Task{}, fmt.Errorf("agent task must be a regular file no larger than 1 MiB")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return Task{}, fmt.Errorf("read agent task: %w", err)
	}
	var task Task
	if err := strictJSON(data, &task); err != nil {
		return Task{}, fmt.Errorf("decode agent task: %w", err)
	}
	return task, nil
}

func (task Task) Validate() error {
	if task.SchemaVersion != TaskSchema {
		return fmt.Errorf("unsupported agent task schema %q", task.SchemaVersion)
	}
	if task.Mode != "suggest" {
		return fmt.Errorf("agent mode %q is unsupported; current providers fail closed to suggest", task.Mode)
	}
	if !hexDigest.MatchString(task.WorkspaceInventoryDigest) {
		return fmt.Errorf("agent task requires a workspace inventory digest")
	}
	if !assertion.MatchString(task.AssertionID) || len(task.ControlIDs) == 0 || strings.TrimSpace(task.Goal) == "" || len(task.Goal) > 16384 {
		return fmt.Errorf("agent task requires a valid assertion, controls, and goal")
	}
	if task.ReadScope != "task-inputs" {
		return fmt.Errorf("agent read scope must be task-inputs")
	}
	if task.RelevantPaths == nil || task.Inputs == nil || task.AllowedPaths == nil || task.ProtectedPaths == nil || task.AllowedCommands == nil {
		return fmt.Errorf("agent task path and command fields must be arrays")
	}
	if len(task.AllowedPaths) == 0 || len(task.ProtectedPaths) == 0 {
		return fmt.Errorf("agent task requires allowed and protected paths")
	}
	if len(task.AllowedCommands) != 0 {
		return fmt.Errorf("suggest providers do not permit command execution")
	}
	if task.Network != "deny" || task.Secrets != "none" {
		return fmt.Errorf("suggest providers require denied tool network and no task secrets")
	}
	if !task.AllowRemoteSourceProcessing {
		return fmt.Errorf("remote source processing requires explicit acknowledgement")
	}
	if task.TimeoutSeconds < 1 || task.TimeoutSeconds > 3600 {
		return fmt.Errorf("agent timeout must be between 1 and 3600 seconds")
	}
	if task.MaxOutputBytes < 1024 || task.MaxOutputBytes > 4*1024*1024 {
		return fmt.Errorf("agent output limit must be between 1 KiB and 4 MiB")
	}
	if math.IsNaN(task.MaxCostUSD) || math.IsInf(task.MaxCostUSD, 0) || task.MaxCostUSD < 0 || task.MaxCostUSD > 1000 {
		return fmt.Errorf("agent cost limit must be between 0 and 1000 USD")
	}
	if err := validateUniqueStrings(task.ControlIDs, control, "control ID"); err != nil {
		return err
	}
	if !slices.Equal(task.ControlIDs, sortedCopy(task.ControlIDs)) {
		return fmt.Errorf("agent task control IDs must be sorted")
	}
	pathGroups := []struct {
		name  string
		paths []string
	}{
		{"relevant", task.RelevantPaths}, {"allowed", task.AllowedPaths}, {"protected", task.ProtectedPaths},
	}
	for _, group := range pathGroups {
		if err := validatePaths(group.paths, group.name); err != nil {
			return err
		}
		if !slices.Equal(group.paths, sortedCopy(group.paths)) {
			return fmt.Errorf("agent task %s paths must be sorted", group.name)
		}
	}
	if len(task.Inputs) != len(task.RelevantPaths) {
		return fmt.Errorf("agent inputs must exactly match relevant paths")
	}
	totalInputBytes := 0
	for index, input := range task.Inputs {
		if input.Path != task.RelevantPaths[index] || !hexDigest.MatchString(input.SHA256) ||
			!utf8.ValidString(input.Content) || strings.IndexByte(input.Content, 0) >= 0 || len(input.Content) > 256*1024 {
			return fmt.Errorf("agent input does not match relevant path %s", task.RelevantPaths[index])
		}
		digest := sha256.Sum256([]byte(input.Content))
		if hex.EncodeToString(digest[:]) != input.SHA256 {
			return fmt.Errorf("agent input digest does not match %s", input.Path)
		}
		totalInputBytes += len(input.Content)
	}
	if totalInputBytes > 768*1024 {
		return fmt.Errorf("agent inputs exceed 768 KiB")
	}
	encodedTask, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("encode agent task: %w", err)
	}
	if len(encodedTask) > 1024*1024 {
		return fmt.Errorf("encoded agent task exceeds 1 MiB")
	}
	for _, path := range task.AllowedPaths {
		if matchesProtected(path, task.ProtectedPaths) {
			return fmt.Errorf("allowed path %s is protected", path)
		}
	}
	want, err := TaskID(task)
	if err != nil {
		return err
	}
	if !hexDigest.MatchString(task.TaskID) || task.TaskID != want {
		return fmt.Errorf("agent task ID does not match its canonical content")
	}
	return nil
}

func TaskID(task Task) (string, error) {
	task.TaskID = ""
	payload, err := json.Marshal(task)
	if err != nil {
		return "", fmt.Errorf("encode agent task identity: %w", err)
	}
	digest := sha256.Sum256(payload)
	return hex.EncodeToString(digest[:]), nil
}

func validateUniqueStrings(values []string, pattern *regexp.Regexp, name string) error {
	seen := map[string]bool{}
	for _, value := range values {
		if !pattern.MatchString(value) {
			return fmt.Errorf("invalid %s %q", name, value)
		}
		if seen[value] {
			return fmt.Errorf("duplicate %s %q", name, value)
		}
		seen[value] = true
	}
	return nil
}

func validatePaths(paths []string, name string) error {
	if len(paths) > 10_000 {
		return fmt.Errorf("%s paths exceed 10000 entries", name)
	}
	seen := map[string]bool{}
	for _, path := range paths {
		if err := validatePathPattern(path); err != nil {
			return fmt.Errorf("%s path: %w", name, err)
		}
		if seen[path] {
			return fmt.Errorf("duplicate %s path %q", name, path)
		}
		seen[path] = true
	}
	return nil
}

func validatePathPattern(path string) error {
	if path == "" || filepath.IsAbs(path) || strings.Contains(path, "\\") || strings.ContainsRune(path, '\x00') {
		return fmt.Errorf("unsafe path %q", path)
	}
	value := strings.TrimSuffix(path, "/")
	clean := filepath.ToSlash(filepath.Clean(filepath.FromSlash(value)))
	if clean == "." || clean == ".." || strings.HasPrefix(clean, "../") || clean != value {
		return fmt.Errorf("unsafe path %q", path)
	}
	return nil
}

func matchesProtected(path string, patterns []string) bool {
	for _, pattern := range patterns {
		if strings.HasSuffix(pattern, "/") {
			if path == strings.TrimSuffix(pattern, "/") || strings.HasPrefix(path, pattern) {
				return true
			}
		} else if path == pattern {
			return true
		}
	}
	return false
}

func strictJSON(data []byte, target any) error {
	if err := rejectDuplicateKeys(data); err != nil {
		return err
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return err
	}
	var trailing any
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		if err == nil {
			return fmt.Errorf("multiple JSON values")
		}
		return err
	}
	return nil
}

func rejectDuplicateKeys(data []byte) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	var walk func() error
	walk = func() error {
		token, err := decoder.Token()
		if err != nil {
			return err
		}
		delimiter, ok := token.(json.Delim)
		if !ok {
			return nil
		}
		switch delimiter {
		case '{':
			seen := map[string]bool{}
			for decoder.More() {
				keyToken, err := decoder.Token()
				if err != nil {
					return err
				}
				key, ok := keyToken.(string)
				if !ok {
					return fmt.Errorf("object key is not a string")
				}
				if seen[key] {
					return fmt.Errorf("duplicate JSON key %q", key)
				}
				seen[key] = true
				if err := walk(); err != nil {
					return err
				}
			}
			_, err = decoder.Token()
			return err
		case '[':
			for decoder.More() {
				if err := walk(); err != nil {
					return err
				}
			}
			_, err = decoder.Token()
			return err
		default:
			return fmt.Errorf("unexpected JSON delimiter %q", delimiter)
		}
	}
	if err := walk(); err != nil {
		return err
	}
	if token, err := decoder.Token(); !errors.Is(err, io.EOF) {
		if err == nil {
			return fmt.Errorf("unexpected trailing JSON value %v", token)
		}
		return err
	}
	return nil
}

func sortedCopy(values []string) []string {
	result := append([]string(nil), values...)
	sort.Strings(result)
	return result
}
