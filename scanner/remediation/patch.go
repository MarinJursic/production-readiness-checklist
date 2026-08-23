package remediation

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
	"github.com/MarinJursic/production-readiness-checklist/scanner/provider"
)

const maximumPatchedFileBytes = 16 * 1024 * 1024

var (
	hunkHeader = regexp.MustCompile(`^@@ -(\d+)(?:,(\d+))? \+(\d+)(?:,(\d+))? @@(?: .*)?$`)
	indexLine  = regexp.MustCompile(`^index [0-9a-f]{7,64}\.\.[0-9a-f]{7,64}( 100644)?$`)
)

type patchLine struct {
	kind byte
	text string
}

type patchHunk struct {
	oldStart int
	oldCount int
	newStart int
	newCount int
	lines    []patchLine
}

type patchFile struct {
	path    string
	newFile bool
	hunks   []patchHunk
}

type textLine struct {
	text    string
	newline bool
}

func applyProviderPatch(candidateRoot string, baseline model.Inventory, task provider.Task, output provider.Output, maxFiles, maxChangedLines int) ([]Change, error) {
	files, err := parseProviderPatch(output.Patch)
	if err != nil {
		return nil, err
	}
	if len(files) != len(output.ChangedFiles) || len(files) > maxFiles {
		return nil, fmt.Errorf("proposal changes %d files, above or inconsistent with the %d-file contract", len(files), maxFiles)
	}
	baselineFiles := make(map[string]model.FileRecord, len(baseline.Files))
	for _, record := range baseline.Files {
		baselineFiles[record.Path] = record
	}
	allowed := make(map[string]bool, len(task.AllowedPaths))
	for _, path := range task.AllowedPaths {
		allowed[path] = true
	}
	protectedPaths := proposalProtectedPaths(task.ProtectedPaths)
	type pendingWrite struct {
		path    string
		content []byte
		mode    os.FileMode
		change  Change
	}
	pending := make([]pendingWrite, 0, len(files))
	changedLines := 0
	for index, filePatch := range files {
		if index >= len(output.ChangedFiles) || filePatch.path != output.ChangedFiles[index] ||
			!allowed[filePatch.path] || protected(filePatch.path, protectedPaths) {
			return nil, fmt.Errorf("proposal patch path is outside the R2 fix contract: %s", filePatch.path)
		}
		record, exists := baselineFiles[filePatch.path]
		if filePatch.newFile == exists {
			return nil, fmt.Errorf("proposal file status does not match the baseline: %s", filePatch.path)
		}
		var source []byte
		mode := os.FileMode(0o644)
		change := Change{Path: filePatch.path, Kind: "added", AfterMode: uint32(mode.Perm())}
		if exists {
			path, err := safeJoin(candidateRoot, filePatch.path)
			if err != nil {
				return nil, err
			}
			source, err = readCandidateText(path, record)
			if err != nil {
				return nil, err
			}
			mode = os.FileMode(record.Mode)
			change.Kind = "modified"
			change.BeforeSHA256 = record.SHA256
			change.BeforeMode = record.Mode
			change.AfterMode = record.Mode
		}
		content, added, removed, err := applyFileHunks(source, filePatch)
		if err != nil {
			return nil, fmt.Errorf("apply proposal patch to %s: %w", filePatch.path, err)
		}
		if exists && string(content) == string(source) {
			return nil, fmt.Errorf("proposal patch makes no change to %s", filePatch.path)
		}
		changedLines += added + removed
		if changedLines > maxChangedLines {
			return nil, fmt.Errorf("proposal changes more than the %d-line contract", maxChangedLines)
		}
		digest := sha256.Sum256(content)
		change.AfterSHA256 = hex.EncodeToString(digest[:])
		change.AddedLines = added
		change.RemovedLines = removed
		pending = append(pending, pendingWrite{path: filePatch.path, content: content, mode: mode, change: change})
	}
	for _, write := range pending {
		if err := writePatchedFile(candidateRoot, write.path, write.content, write.mode, write.change.Kind == "added"); err != nil {
			return nil, err
		}
	}
	changes := make([]Change, 0, len(pending))
	for _, write := range pending {
		changes = append(changes, write.change)
	}
	sort.Slice(changes, func(left, right int) bool { return changes[left].Path < changes[right].Path })
	return changes, nil
}

func parseProviderPatch(patch string) ([]patchFile, error) {
	if patch == "" || !strings.HasSuffix(patch, "\n") || strings.Contains(patch, "\r") {
		return nil, fmt.Errorf("proposal patch must be LF-terminated unified diff text")
	}
	lines := strings.Split(strings.TrimSuffix(patch, "\n"), "\n")
	files := make([]patchFile, 0)
	for index := 0; index < len(lines); {
		fields := strings.Fields(lines[index])
		if len(fields) != 4 || fields[0] != "diff" || fields[1] != "--git" ||
			!strings.HasPrefix(fields[2], "a/") || !strings.HasPrefix(fields[3], "b/") {
			return nil, fmt.Errorf("proposal patch has an unsupported file header")
		}
		path := strings.TrimPrefix(fields[2], "a/")
		if strings.TrimPrefix(fields[3], "b/") != path {
			return nil, fmt.Errorf("proposal patch renames or redirects %s", path)
		}
		index++
		newMode := false
		if index < len(lines) && strings.HasPrefix(lines[index], "new file mode ") {
			if lines[index] != "new file mode 100644" {
				return nil, fmt.Errorf("proposal patch requests an unsupported new-file mode")
			}
			newMode = true
			index++
		}
		if index < len(lines) && strings.HasPrefix(lines[index], "index ") {
			if !indexLine.MatchString(lines[index]) {
				return nil, fmt.Errorf("proposal patch has malformed index metadata")
			}
			index++
		}
		if index+1 >= len(lines) {
			return nil, fmt.Errorf("proposal patch section for %s has no path headers", path)
		}
		oldHeader, newHeader := lines[index], lines[index+1]
		newFile := oldHeader == "--- /dev/null"
		if (!newFile && oldHeader != "--- a/"+path) || newHeader != "+++ b/"+path || (newMode && !newFile) {
			return nil, fmt.Errorf("proposal patch has invalid path headers for %s", path)
		}
		index += 2
		filePatch := patchFile{path: path, newFile: newFile}
		for index < len(lines) && !strings.HasPrefix(lines[index], "diff --git ") {
			matches := hunkHeader.FindStringSubmatch(lines[index])
			if matches == nil {
				return nil, fmt.Errorf("proposal patch has unsupported metadata or malformed hunk for %s", path)
			}
			hunk, err := parsedHunkHeader(matches)
			if err != nil {
				return nil, err
			}
			index++
			oldLines, newLines := 0, 0
			for index < len(lines) && !strings.HasPrefix(lines[index], "@@ ") && !strings.HasPrefix(lines[index], "diff --git ") {
				line := lines[index]
				if line == "\\ No newline at end of file" {
					return nil, fmt.Errorf("proposal patch uses unsupported non-newline markers")
				}
				if line == "" || !strings.Contains(" +-", line[:1]) {
					return nil, fmt.Errorf("proposal patch has malformed hunk content for %s", path)
				}
				entry := patchLine{kind: line[0], text: line[1:]}
				hunk.lines = append(hunk.lines, entry)
				if entry.kind != '+' {
					oldLines++
				}
				if entry.kind != '-' {
					newLines++
				}
				index++
			}
			if oldLines != hunk.oldCount || newLines != hunk.newCount || len(hunk.lines) == 0 {
				return nil, fmt.Errorf("proposal patch hunk counts do not match for %s", path)
			}
			filePatch.hunks = append(filePatch.hunks, hunk)
		}
		if len(filePatch.hunks) == 0 {
			return nil, fmt.Errorf("proposal patch has no hunks for %s", path)
		}
		files = append(files, filePatch)
	}
	return files, nil
}

func parsedHunkHeader(matches []string) (patchHunk, error) {
	values := [4]int{}
	for index, value := range matches[1:] {
		if value == "" {
			if index == 1 || index == 3 {
				values[index] = 1
			}
			continue
		}
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed < 0 {
			return patchHunk{}, fmt.Errorf("proposal patch has invalid hunk coordinates")
		}
		values[index] = parsed
	}
	return patchHunk{oldStart: values[0], oldCount: values[1], newStart: values[2], newCount: values[3]}, nil
}

func applyFileHunks(source []byte, filePatch patchFile) ([]byte, int, int, error) {
	sourceLines := splitTextLines(string(source))
	result := make([]textLine, 0, len(sourceLines))
	cursor, added, removed := 0, 0, 0
	for _, hunk := range filePatch.hunks {
		position := hunk.oldStart - 1
		if hunk.oldStart == 0 && hunk.oldCount == 0 {
			position = 0
		}
		if position < cursor || position > len(sourceLines) {
			return nil, 0, 0, fmt.Errorf("hunks are overlapping or outside the source")
		}
		result = append(result, sourceLines[cursor:position]...)
		if len(result) != hunk.newStart-1 && !(hunk.newStart == 0 && hunk.newCount == 0) {
			return nil, 0, 0, fmt.Errorf("new-file hunk coordinates are inconsistent")
		}
		cursor = position
		for _, line := range hunk.lines {
			switch line.kind {
			case ' ':
				if cursor >= len(sourceLines) || sourceLines[cursor].text != line.text || !sourceLines[cursor].newline {
					return nil, 0, 0, fmt.Errorf("context does not match the source")
				}
				result = append(result, sourceLines[cursor])
				cursor++
			case '-':
				if cursor >= len(sourceLines) || sourceLines[cursor].text != line.text || !sourceLines[cursor].newline {
					return nil, 0, 0, fmt.Errorf("removed line does not match the source")
				}
				cursor++
				removed++
			case '+':
				result = append(result, textLine{text: line.text, newline: true})
				added++
			default:
				return nil, 0, 0, fmt.Errorf("unsupported patch line")
			}
		}
	}
	result = append(result, sourceLines[cursor:]...)
	var builder strings.Builder
	for _, line := range result {
		builder.WriteString(line.text)
		if line.newline {
			builder.WriteByte('\n')
		}
	}
	return []byte(builder.String()), added, removed, nil
}

func splitTextLines(content string) []textLine {
	if content == "" {
		return nil
	}
	lines := make([]textLine, 0, strings.Count(content, "\n")+1)
	for len(content) > 0 {
		index := strings.IndexByte(content, '\n')
		if index < 0 {
			lines = append(lines, textLine{text: content})
			break
		}
		lines = append(lines, textLine{text: content[:index], newline: true})
		content = content[index+1:]
	}
	return lines
}

func readCandidateText(path string, record model.FileRecord) ([]byte, error) {
	info, err := os.Lstat(path)
	if err != nil || !info.Mode().IsRegular() || info.Size() > maximumPatchedFileBytes || uint32(info.Mode().Perm()) != record.Mode {
		return nil, fmt.Errorf("candidate text path changed or exceeds %d bytes: %s", maximumPatchedFileBytes, record.Path)
	}
	data, err := os.ReadFile(path)
	if err != nil || !strings.HasSuffix(string(data), "\n") {
		return nil, fmt.Errorf("candidate text path must be readable LF-terminated text: %s", record.Path)
	}
	digest := sha256.Sum256(data)
	if int64(len(data)) != record.Size || hex.EncodeToString(digest[:]) != record.SHA256 {
		return nil, fmt.Errorf("candidate text path changed after inventory: %s", record.Path)
	}
	return data, nil
}

func writePatchedFile(root, relative string, content []byte, mode os.FileMode, added bool) error {
	path, err := safeJoin(root, relative)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("create proposal parent for %s: %w", relative, err)
	}
	if added {
		file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, mode)
		if err != nil {
			return fmt.Errorf("create proposal file %s: %w", relative, err)
		}
		if _, err := file.Write(content); err != nil {
			file.Close()
			return fmt.Errorf("write proposal file %s: %w", relative, err)
		}
		if err := file.Sync(); err != nil {
			file.Close()
			return fmt.Errorf("sync proposal file %s: %w", relative, err)
		}
		return file.Close()
	}
	info, err := os.Lstat(path)
	if err != nil || !info.Mode().IsRegular() {
		return fmt.Errorf("proposal target is not a regular file: %s", relative)
	}
	temporary, err := os.CreateTemp(filepath.Dir(path), ".prc-proposal-*")
	if err != nil {
		return fmt.Errorf("create proposal temporary file: %w", err)
	}
	temporaryPath := temporary.Name()
	complete := false
	defer func() {
		if !complete {
			_ = os.Remove(temporaryPath)
		}
	}()
	if err := temporary.Chmod(mode); err != nil {
		temporary.Close()
		return fmt.Errorf("set proposal file mode: %w", err)
	}
	if _, err := temporary.Write(content); err != nil {
		temporary.Close()
		return fmt.Errorf("write proposal file %s: %w", relative, err)
	}
	if err := temporary.Sync(); err != nil {
		temporary.Close()
		return fmt.Errorf("sync proposal file %s: %w", relative, err)
	}
	if err := temporary.Close(); err != nil {
		return fmt.Errorf("close proposal file %s: %w", relative, err)
	}
	if err := os.Rename(temporaryPath, path); err != nil {
		return fmt.Errorf("replace proposal file %s: %w", relative, err)
	}
	complete = true
	return nil
}
