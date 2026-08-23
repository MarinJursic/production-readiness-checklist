package remediation

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
)

const maximumCandidateBytes int64 = 2 * 1024 * 1024 * 1024

func prepareCandidate(baseline model.Inventory, destination string) (string, error) {
	root := filepath.Clean(baseline.Root)
	destination, err := filepath.Abs(destination)
	if err != nil {
		return "", fmt.Errorf("resolve candidate directory: %w", err)
	}
	parent := filepath.Dir(destination)
	parentInfo, err := os.Stat(parent)
	if err != nil || !parentInfo.IsDir() {
		return "", fmt.Errorf("candidate parent must be an existing directory")
	}
	if resolved, resolveErr := filepath.EvalSymlinks(parent); resolveErr == nil {
		destination = filepath.Join(resolved, filepath.Base(destination))
	}
	if pathWithin(root, destination) || pathWithin(destination, root) {
		return "", fmt.Errorf("candidate directory must be outside the target tree")
	}
	if _, err := os.Lstat(destination); !os.IsNotExist(err) {
		return "", fmt.Errorf("candidate directory already exists")
	}
	var totalBytes int64
	for _, record := range baseline.Files {
		totalBytes += record.Size
		if totalBytes > maximumCandidateBytes {
			return "", fmt.Errorf("candidate exceeds the %d-byte copy limit", maximumCandidateBytes)
		}
	}
	if err := os.Mkdir(destination, 0o700); err != nil {
		return "", fmt.Errorf("create candidate directory: %w", err)
	}
	complete := false
	defer func() {
		if !complete {
			_ = os.RemoveAll(destination)
		}
	}()
	for _, record := range baseline.Files {
		if err := copyInventoryFile(root, destination, record); err != nil {
			return "", err
		}
	}
	complete = true
	return destination, nil
}

func copyInventoryFile(root, destination string, record model.FileRecord) error {
	source, err := safeJoin(root, record.Path)
	if err != nil {
		return err
	}
	pathInfo, err := os.Lstat(source)
	if err != nil || !pathInfo.Mode().IsRegular() {
		return fmt.Errorf("inventory source changed to a non-regular file: %s", record.Path)
	}
	input, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("open candidate source %s: %w", record.Path, err)
	}
	defer input.Close()
	openedInfo, err := input.Stat()
	if err != nil || !os.SameFile(pathInfo, openedInfo) || uint32(openedInfo.Mode().Perm()) != record.Mode {
		return fmt.Errorf("inventory source changed while opening: %s", record.Path)
	}
	outputPath, err := safeJoin(destination, record.Path)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(outputPath), 0o700); err != nil {
		return fmt.Errorf("create candidate parent: %w", err)
	}
	output, err := os.OpenFile(outputPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, os.FileMode(record.Mode))
	if err != nil {
		return fmt.Errorf("create candidate file %s: %w", record.Path, err)
	}
	hasher := sha256.New()
	size, copyErr := io.Copy(io.MultiWriter(output, hasher), input)
	chmodErr := output.Chmod(os.FileMode(record.Mode))
	syncErr := output.Sync()
	closeErr := output.Close()
	if copyErr != nil || chmodErr != nil || syncErr != nil || closeErr != nil {
		return fmt.Errorf("copy candidate file %s", record.Path)
	}
	if size != record.Size || hex.EncodeToString(hasher.Sum(nil)) != record.SHA256 {
		return fmt.Errorf("inventory source changed while copying: %s", record.Path)
	}
	return nil
}

func safeJoin(root, relative string) (string, error) {
	if relative == "" || filepath.IsAbs(relative) || strings.Contains(relative, "\\") {
		return "", fmt.Errorf("unsafe candidate path %q", relative)
	}
	clean := filepath.Clean(filepath.FromSlash(relative))
	if clean == "." || clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) ||
		filepath.ToSlash(clean) != relative {
		return "", fmt.Errorf("unsafe candidate path %q", relative)
	}
	joined := filepath.Join(root, clean)
	if !pathWithin(root, joined) {
		return "", fmt.Errorf("candidate path escapes root: %q", relative)
	}
	return joined, nil
}

func pathWithin(root, path string) bool {
	relative, err := filepath.Rel(filepath.Clean(root), filepath.Clean(path))
	return err == nil && relative != ".." && !strings.HasPrefix(relative, ".."+string(filepath.Separator))
}
