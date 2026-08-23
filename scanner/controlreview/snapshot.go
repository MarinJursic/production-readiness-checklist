package controlreview

import (
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
	Directory  string
	LineCounts map[string]int
	Paths      []string
	Contents   map[string][]byte
	Digests    map[string]string
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
		Contents: map[string][]byte{}, Digests: map[string]string{},
	}
	total := int64(0)
	for _, record := range inventory.Files {
		if !remoteReviewCandidate(record.Path) || record.Size > maximumSnapshotFileBytes {
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
	return result, nil
}

func remoteReviewCandidate(path string) bool {
	lower := strings.ToLower(filepath.ToSlash(path))
	parts := strings.Split(lower, "/")
	for _, part := range parts {
		if part == ".claude" || part == ".codex" || part == ".agents" || part == ".git" {
			return false
		}
	}
	base := parts[len(parts)-1]
	if base == "agents.md" || base == "claude.md" || base == "copilot-instructions.md" ||
		base == ".env" || strings.HasPrefix(base, ".env.") || base == ".npmrc" || base == ".pypirc" ||
		base == ".netrc" || base == "auth.json" || base == "credentials" ||
		base == "id_rsa" || base == "id_ed25519" || strings.Contains(base, "secret") ||
		strings.Contains(base, "credential") {
		return false
	}
	switch strings.ToLower(filepath.Ext(base)) {
	case ".key", ".pem", ".p12", ".pfx", ".jks", ".keystore", ".kdbx":
		return false
	}
	extension := strings.ToLower(filepath.Ext(base))
	return reviewExtensions[extension] || reviewNames[base]
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
