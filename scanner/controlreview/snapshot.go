package controlreview

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
	"github.com/MarinJursic/production-readiness-checklist/scanner/provider"
)

const (
	maximumSnapshotFileBytes  = 2 * 1024 * 1024
	maximumSnapshotTotalBytes = 128 * 1024 * 1024
	maximumReviewIgnoreBytes  = 64 * 1024
	maximumReviewIgnoreFiles  = 100
	reviewIgnoreName          = ".prcreviewignore"
)

var reviewExtensions = map[string]bool{
	".c": true, ".cc": true, ".conf": true, ".cpp": true, ".cs": true, ".css": true,
	".go": true, ".gradle": true, ".graphql": true, ".h": true, ".hcl": true,
	".html": true, ".ini": true, ".java": true, ".js": true, ".json": true,
	".jsx": true, ".kt": true, ".kts": true, ".md": true, ".php": true,
	".properties": true, ".proto": true, ".py": true, ".rb": true, ".rs": true,
	".rst": true, ".scala": true, ".sh": true, ".sql": true, ".swift": true,
	".tf": true, ".toml": true, ".ts": true, ".tsx": true, ".txt": true,
	".xml": true, ".yaml": true, ".yml": true,
}

var reviewNames = map[string]bool{
	"dockerfile": true, "makefile": true, "procfile": true,
}

type reviewSnapshot struct {
	Directory   string
	LineCounts  map[string]int
	Paths       []string
	Contents    map[string][]byte
	Digests     map[string]string
	Limitations []string
	Bytes       int64
	Omitted     int
}

type reviewExclusion struct {
	Path   string
	Reason string
}

func createSnapshot(inventory model.Inventory) (reviewSnapshot, error) {
	directory, err := os.MkdirTemp("", "prc-review-snapshot-")
	if err != nil {
		return reviewSnapshot{}, fmt.Errorf("create private review snapshot: %w", err)
	}
	cleanup := func(cause error) (reviewSnapshot, error) {
		_ = os.RemoveAll(directory)
		return reviewSnapshot{}, cause
	}
	if err := os.Chmod(directory, 0o700); err != nil {
		return cleanup(fmt.Errorf("protect review snapshot: %w", err))
	}
	result := reviewSnapshot{
		Directory: directory, LineCounts: map[string]int{}, Paths: []string{},
		Contents: map[string][]byte{}, Digests: map[string]string{}, Limitations: []string{},
	}
	exclusions, err := loadReviewExclusions(inventory)
	if err != nil {
		return cleanup(err)
	}
	omitted := map[string]int{}
	total := int64(0)
	for _, record := range inventory.Files {
		candidate, reason := remoteReviewCandidate(record.Path)
		if !candidate {
			omitted[reason]++
			continue
		}
		if record.Size > maximumSnapshotFileBytes {
			omitted["larger than the 2 MiB per-file remote-review limit"]++
			continue
		}
		if exclusion, excluded := exclusions[record.Path]; excluded {
			// Reopen and hash the exact inventoried file immediately before
			// omitting it. The bytes stay local, but a changed path cannot use
			// the ignore file to widen the remote-review snapshot unnoticed.
			if _, err := readInventoryFile(inventory.Root, record); err != nil {
				return cleanup(err)
			}
			result.Omitted++
			result.Limitations = append(result.Limitations, fmt.Sprintf(
				"Remote AI review intentionally omitted %q: %s", exclusion.Path, exclusion.Reason,
			))
			continue
		}
		if total+record.Size > maximumSnapshotTotalBytes {
			return cleanup(fmt.Errorf("safe remote-review snapshot exceeds %d bytes", maximumSnapshotTotalBytes))
		}
		data, err := readInventoryFile(inventory.Root, record)
		if err != nil {
			return cleanup(err)
		}
		if bytes.IndexByte(data, 0) >= 0 || !utf8.Valid(data) {
			omitted["binary or not valid UTF-8 text"]++
			continue
		}
		if err := provider.ScreenRemoteContent(record.Path, data); err != nil {
			return cleanup(err)
		}
		destination := filepath.Join(directory, filepath.FromSlash(record.Path))
		if err := os.MkdirAll(filepath.Dir(destination), 0o700); err != nil {
			return cleanup(fmt.Errorf("create review snapshot directory: %w", err))
		}
		file, err := os.OpenFile(destination, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
		if err != nil {
			return cleanup(fmt.Errorf("create review snapshot file: %w", err))
		}
		if _, err := file.Write(data); err != nil {
			_ = file.Close()
			return cleanup(fmt.Errorf("write review snapshot file: %w", err))
		}
		if err := file.Close(); err != nil {
			return cleanup(fmt.Errorf("finish review snapshot file: %w", err))
		}
		result.LineCounts[record.Path] = bytes.Count(data, []byte{'\n'}) + 1
		result.Paths = append(result.Paths, record.Path)
		result.Contents[record.Path] = data
		result.Digests[record.Path] = record.SHA256
		total += record.Size
	}
	for reason, count := range omitted {
		result.Omitted += count
		result.Limitations = append(result.Limitations, fmt.Sprintf(
			"The scanner omitted %d inventoried file(s) from remote review because they were %s.", count, reason,
		))
	}
	result.Bytes = total
	result.Limitations = uniqueSorted(result.Limitations)
	return result, nil
}

// loadReviewExclusions reads an optional, inventoried root file that limits
// only what may be sent to a remote AI reviewer. It does not remove files from
// the local inventory, local checks, adapters, or authoritative evidence.
// Entries are exact regular-file paths; globs and directories are rejected so
// every omission remains visible and reviewable.
func loadReviewExclusions(inventory model.Inventory) (map[string]reviewExclusion, error) {
	records := make(map[string]model.FileRecord, len(inventory.Files))
	var config *model.FileRecord
	for index := range inventory.Files {
		record := inventory.Files[index]
		records[record.Path] = record
		if record.Path == reviewIgnoreName {
			candidate := record
			config = &candidate
		}
	}
	if config == nil {
		return map[string]reviewExclusion{}, nil
	}
	if config.Size > maximumReviewIgnoreBytes {
		return nil, fmt.Errorf("%s exceeds the %d-byte limit", reviewIgnoreName, maximumReviewIgnoreBytes)
	}
	data, err := readInventoryFile(inventory.Root, *config)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", reviewIgnoreName, err)
	}
	result := map[string]reviewExclusion{}
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for line := 1; scanner.Scan(); line++ {
		value := strings.TrimSpace(scanner.Text())
		if value == "" || strings.HasPrefix(value, "#") {
			continue
		}
		parts := strings.SplitN(value, "|", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("%s line %d must use `relative/file | reviewed reason`", reviewIgnoreName, line)
		}
		relative := strings.TrimSpace(parts[0])
		reason := strings.Join(strings.Fields(parts[1]), " ")
		if relative == "" || relative == reviewIgnoreName || filepath.IsAbs(relative) ||
			strings.Contains(relative, `\`) || filepath.ToSlash(filepath.Clean(filepath.FromSlash(relative))) != relative ||
			relative == "." || relative == ".." || strings.HasPrefix(relative, "../") {
			return nil, fmt.Errorf("%s line %d contains an unsafe file path %q", reviewIgnoreName, line, relative)
		}
		if len(reason) < 10 || len(reason) > 300 {
			return nil, fmt.Errorf("%s line %d requires a 10-300 character reviewed reason", reviewIgnoreName, line)
		}
		if _, exists := result[relative]; exists {
			return nil, fmt.Errorf("%s repeats %q", reviewIgnoreName, relative)
		}
		if len(result) >= maximumReviewIgnoreFiles {
			return nil, fmt.Errorf("%s exceeds the %d-entry limit", reviewIgnoreName, maximumReviewIgnoreFiles)
		}
		record, exists := records[relative]
		if !exists {
			return nil, fmt.Errorf("%s path %q is not an inventoried regular file", reviewIgnoreName, relative)
		}
		candidate, candidateReason := remoteReviewCandidate(record.Path)
		if !candidate {
			return nil, fmt.Errorf("%s path %q is already omitted because it is %s", reviewIgnoreName, relative, candidateReason)
		}
		if record.Size > maximumSnapshotFileBytes {
			return nil, fmt.Errorf("%s path %q is already omitted because it exceeds the 2 MiB per-file remote-review limit", reviewIgnoreName, relative)
		}
		result[relative] = reviewExclusion{Path: relative, Reason: reason}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read %s: %w", reviewIgnoreName, err)
	}
	return result, nil
}

func remoteReviewCandidate(path string) (bool, string) {
	lower := strings.ToLower(filepath.ToSlash(path))
	parts := strings.Split(lower, "/")
	for _, part := range parts {
		if part == ".claude" || part == ".codex" || part == ".agents" || part == ".git" {
			return false, "provider-instruction or scanner-private content"
		}
	}
	base := parts[len(parts)-1]
	if base == "agents.md" || base == "claude.md" || base == "copilot-instructions.md" {
		return false, "provider-instruction content"
	}
	if base == ".env" || strings.HasPrefix(base, ".env.") || base == ".npmrc" || base == ".pypirc" ||
		base == ".netrc" || base == "auth.json" || base == "credentials" ||
		base == "id_rsa" || base == "id_ed25519" || strings.Contains(base, "secret") ||
		strings.Contains(base, "credential") {
		return false, "named like sensitive configuration or credentials"
	}
	switch strings.ToLower(filepath.Ext(base)) {
	case ".key", ".pem", ".p12", ".pfx", ".jks", ".keystore", ".kdbx":
		return false, "a private-key, credential-store, or certificate container"
	}
	extension := strings.ToLower(filepath.Ext(base))
	if reviewExtensions[extension] || reviewNames[base] {
		return true, ""
	}
	return false, "not a supported remote-review text type"
}

func readInventoryFile(root string, record model.FileRecord) ([]byte, error) {
	if record.Path == "" || filepath.IsAbs(record.Path) ||
		filepath.ToSlash(filepath.Clean(filepath.FromSlash(record.Path))) != record.Path {
		return nil, fmt.Errorf("inventory contains unsafe review path %q", record.Path)
	}
	path := filepath.Join(root, filepath.FromSlash(record.Path))
	before, err := os.Lstat(path)
	if err != nil || !before.Mode().IsRegular() || before.Mode()&os.ModeSymlink != 0 || before.Size() != record.Size {
		return nil, fmt.Errorf("review input %s changed after inventory", record.Path)
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open review input %s: %w", record.Path, err)
	}
	defer file.Close()
	opened, err := file.Stat()
	if err != nil || !os.SameFile(before, opened) {
		return nil, fmt.Errorf("review input %s changed while opening", record.Path)
	}
	data, err := io.ReadAll(io.LimitReader(file, maximumSnapshotFileBytes+1))
	if err != nil || int64(len(data)) != record.Size {
		return nil, fmt.Errorf("review input %s changed or exceeded its byte limit", record.Path)
	}
	digest := sha256.Sum256(data)
	if hex.EncodeToString(digest[:]) != record.SHA256 {
		return nil, fmt.Errorf("review input %s no longer matches its inventory digest", record.Path)
	}
	return data, nil
}
