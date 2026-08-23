package provider

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

func ParseOutput(provider string, data []byte, task Task) (Output, error) {
	if err := task.Validate(); err != nil {
		return Output{}, err
	}
	if len(data) == 0 || len(data) > task.MaxOutputBytes {
		return Output{}, fmt.Errorf("agent output must be between 1 and %d bytes", task.MaxOutputBytes)
	}
	if provider == "claude" {
		unwrapped, err := unwrapClaude(data)
		if err != nil {
			return Output{}, err
		}
		data = unwrapped
	} else if provider != "codex" {
		return Output{}, fmt.Errorf("unsupported agent provider %q", provider)
	}
	var output Output
	if err := strictJSON(data, &output); err != nil {
		return Output{}, fmt.Errorf("decode agent output: %w", err)
	}
	if err := validateOutput(provider, output, task); err != nil {
		return Output{}, err
	}
	return output, nil
}

// ValidateOutput validates an already-decoded proposal against its sealed task.
// Callers must treat the proposal as untrusted until this returns nil.
func ValidateOutput(provider string, output Output, task Task) error {
	if err := task.Validate(); err != nil {
		return err
	}
	return validateOutput(provider, output, task)
}

// OutputID returns the canonical SHA-256 identity of one decoded proposal.
func OutputID(output Output) (string, error) {
	payload, err := json.Marshal(output)
	if err != nil {
		return "", fmt.Errorf("encode agent output identity: %w", err)
	}
	digest := sha256.Sum256(payload)
	return hex.EncodeToString(digest[:]), nil
}

func unwrapClaude(data []byte) ([]byte, error) {
	if err := rejectDuplicateKeys(data); err != nil {
		return nil, fmt.Errorf("decode claude envelope: %w", err)
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	var envelope map[string]json.RawMessage
	if err := decoder.Decode(&envelope); err != nil {
		return nil, fmt.Errorf("decode claude envelope: %w", err)
	}
	var trailing any
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		return nil, fmt.Errorf("claude envelope contains trailing JSON")
	}
	if raw, ok := envelope["is_error"]; ok {
		var isError bool
		if err := json.Unmarshal(raw, &isError); err != nil {
			return nil, fmt.Errorf("claude envelope has invalid is_error")
		}
		if isError {
			return nil, fmt.Errorf("claude reported an execution error")
		}
	}
	structured, ok := envelope["structured_output"]
	if !ok || bytes.Equal(bytes.TrimSpace(structured), []byte("null")) {
		return nil, fmt.Errorf("claude envelope has no structured_output")
	}
	return structured, nil
}

func validateOutput(provider string, output Output, task Task) error {
	if output.SchemaVersion != OutputSchema || output.TaskID != task.TaskID {
		return fmt.Errorf("agent output schema or task identity does not match")
	}
	if output.ChangedFiles == nil || output.CommandsRequestedOrRun == nil ||
		output.Limitations == nil || output.RequestedCapabilityChanges == nil {
		return fmt.Errorf("agent output list fields must be arrays")
	}
	if strings.TrimSpace(output.RootCause) == "" {
		return fmt.Errorf("agent output requires a root-cause explanation")
	}
	if len(output.RequestedCapabilityChanges) != 0 {
		return fmt.Errorf("agent requested a capability change")
	}
	if len(output.CommandsRequestedOrRun) != 0 {
		return fmt.Errorf("suggest provider reported command execution")
	}
	if provider == "claude" || provider == "codex" {
		// Provider launch plans expose no shell, edit, web, or MCP tools.
	} else {
		return fmt.Errorf("unsupported agent provider %q", provider)
	}
	allowed := map[string]bool{}
	for _, path := range task.AllowedPaths {
		allowed[path] = true
	}
	seen := map[string]bool{}
	for index, path := range output.ChangedFiles {
		if err := validatePathPattern(path); err != nil {
			return fmt.Errorf("agent changed file: %w", err)
		}
		if !allowed[path] || matchesProtected(path, task.ProtectedPaths) {
			return fmt.Errorf("agent proposed path outside the fix contract: %s", path)
		}
		if seen[path] {
			return fmt.Errorf("agent output repeats changed file %s", path)
		}
		seen[path] = true
		if index > 0 && output.ChangedFiles[index-1] > path {
			return fmt.Errorf("agent changed files must be sorted")
		}
	}
	if len(output.Patch) > task.MaxOutputBytes {
		return fmt.Errorf("agent patch exceeds the output budget")
	}
	switch output.Status {
	case "candidate":
		if len(output.ChangedFiles) == 0 || strings.TrimSpace(output.Patch) == "" {
			return fmt.Errorf("candidate output requires changed files and a patch")
		}
		patchFiles, err := validateUnifiedDiff(output.Patch, task)
		if err != nil {
			return err
		}
		if len(patchFiles) != len(output.ChangedFiles) {
			return fmt.Errorf("agent patch paths do not match changed_files")
		}
		for index := range patchFiles {
			if patchFiles[index] != output.ChangedFiles[index] {
				return fmt.Errorf("agent patch paths do not match changed_files")
			}
		}
	case "unable", "needs_escalation":
		if len(output.ChangedFiles) != 0 || output.Patch != "" {
			return fmt.Errorf("non-candidate output cannot contain a patch")
		}
	default:
		return fmt.Errorf("unsupported agent output status %q", output.Status)
	}
	return nil
}

func validateUnifiedDiff(patch string, task Task) ([]string, error) {
	allowed := map[string]bool{}
	for _, path := range task.AllowedPaths {
		allowed[path] = true
	}
	paths := make([]string, 0)
	seen := map[string]bool{}
	current := ""
	hasOld, hasNew, hasHunk := false, false, false
	finish := func() error {
		if current != "" && (!hasOld || !hasNew || !hasHunk) {
			return fmt.Errorf("agent patch section for %s is incomplete", current)
		}
		return nil
	}
	scanner := bufio.NewScanner(strings.NewReader(patch))
	scanner.Buffer(make([]byte, 64*1024), task.MaxOutputBytes)
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case strings.HasPrefix(line, "diff --git "):
			if err := finish(); err != nil {
				return nil, err
			}
			fields := strings.Fields(line)
			if len(fields) != 4 || !strings.HasPrefix(fields[2], "a/") || !strings.HasPrefix(fields[3], "b/") {
				return nil, fmt.Errorf("agent patch has an unsupported diff header")
			}
			left, right := strings.TrimPrefix(fields[2], "a/"), strings.TrimPrefix(fields[3], "b/")
			if left != right || !allowed[left] || matchesProtected(left, task.ProtectedPaths) {
				return nil, fmt.Errorf("agent patch contains path outside the fix contract: %s", left)
			}
			if err := validatePathPattern(left); err != nil {
				return nil, err
			}
			if seen[left] {
				return nil, fmt.Errorf("agent patch repeats path %s", left)
			}
			seen[left] = true
			paths = append(paths, left)
			current, hasOld, hasNew, hasHunk = left, false, false, false
		case strings.HasPrefix(line, "--- ") && !hasHunk:
			if current == "" || (line != "--- /dev/null" && line != "--- a/"+current) {
				return nil, fmt.Errorf("agent patch has an invalid old path header")
			}
			hasOld = true
		case strings.HasPrefix(line, "+++ ") && !hasHunk:
			if current == "" || line != "+++ b/"+current {
				return nil, fmt.Errorf("agent patch deletes or redirects a file")
			}
			hasNew = true
		case strings.HasPrefix(line, "@@ "):
			if current == "" || !hasOld || !hasNew {
				return nil, fmt.Errorf("agent patch has a hunk outside a valid file section")
			}
			hasHunk = true
		case strings.HasPrefix(line, "GIT binary patch"), strings.HasPrefix(line, "Binary files "),
			strings.HasPrefix(line, "deleted file mode"), strings.HasPrefix(line, "old mode"),
			strings.HasPrefix(line, "new mode"), strings.HasPrefix(line, "rename from"),
			strings.HasPrefix(line, "rename to"), strings.HasPrefix(line, "copy from"),
			strings.HasPrefix(line, "copy to"):
			return nil, fmt.Errorf("agent patch contains unsupported binary, deletion, rename, copy, or mode metadata")
		default:
			if current == "" && strings.TrimSpace(line) != "" {
				return nil, fmt.Errorf("agent patch contains content before its first diff header")
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read agent patch: %w", err)
	}
	if err := finish(); err != nil {
		return nil, err
	}
	if len(paths) == 0 {
		return nil, fmt.Errorf("agent candidate contains no unified diff")
	}
	for index := 1; index < len(paths); index++ {
		if paths[index-1] > paths[index] {
			return nil, fmt.Errorf("agent patch paths must be sorted")
		}
	}
	return paths, nil
}
