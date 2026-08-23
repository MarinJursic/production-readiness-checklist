package remediation

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	workspaceinventory "github.com/MarinJursic/production-readiness-checklist/scanner/inventory"
	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
)

var defaultProtectedPaths = []string{
	".git/",
	".github/workflows/",
	".prc/",
	"catalog/",
	"production-readiness.yaml",
	"schemas/",
}

func finalNewlineViolations(item model.Inventory) ([]string, error) {
	violations := make([]string, 0)
	for _, record := range item.Files {
		if !workspaceinventory.IsSourcePath(record.Path) {
			continue
		}
		lastByte, hasContent, err := verifiedLastByte(item.Root, record)
		if err != nil {
			return nil, err
		}
		if !hasContent || lastByte != '\n' {
			violations = append(violations, record.Path)
		}
	}
	sort.Strings(violations)
	return violations, nil
}

func verifiedLastByte(root string, record model.FileRecord) (byte, bool, error) {
	path, err := safeJoin(root, record.Path)
	if err != nil {
		return 0, false, err
	}
	pathInfo, err := os.Lstat(path)
	if err != nil || !pathInfo.Mode().IsRegular() || uint32(pathInfo.Mode().Perm()) != record.Mode {
		return 0, false, fmt.Errorf("target changed after inventory: %s", record.Path)
	}
	file, err := os.Open(path)
	if err != nil {
		return 0, false, fmt.Errorf("open %s: %w", record.Path, err)
	}
	defer file.Close()
	openedInfo, err := file.Stat()
	if err != nil || !os.SameFile(pathInfo, openedInfo) || uint32(openedInfo.Mode().Perm()) != record.Mode {
		return 0, false, fmt.Errorf("target changed while opening: %s", record.Path)
	}
	hasher := sha256.New()
	buffer := make([]byte, 32*1024)
	var size int64
	var lastByte byte
	hasContent := false
	for {
		count, readErr := file.Read(buffer)
		if count > 0 {
			hasContent = true
			lastByte = buffer[count-1]
			size += int64(count)
			_, _ = hasher.Write(buffer[:count])
		}
		if readErr != nil {
			if readErr == io.EOF {
				break
			}
			return 0, false, fmt.Errorf("read %s: %w", record.Path, readErr)
		}
	}
	if size != record.Size || hex.EncodeToString(hasher.Sum(nil)) != record.SHA256 {
		return 0, false, fmt.Errorf("target changed after inventory: %s", record.Path)
	}
	return lastByte, hasContent, nil
}

func applyFinalNewline(candidateRoot string, paths []string) error {
	for _, relative := range paths {
		path, err := safeJoin(candidateRoot, relative)
		if err != nil {
			return err
		}
		pathInfo, err := os.Lstat(path)
		if err != nil || !pathInfo.Mode().IsRegular() {
			return fmt.Errorf("candidate path is not a regular file: %s", relative)
		}
		file, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND, 0)
		if err != nil {
			return fmt.Errorf("open candidate file %s: %w", relative, err)
		}
		openedInfo, statErr := file.Stat()
		if statErr != nil || !os.SameFile(pathInfo, openedInfo) {
			file.Close()
			return fmt.Errorf("candidate file changed while opening: %s", relative)
		}
		_, writeErr := file.Write([]byte{'\n'})
		syncErr := file.Sync()
		closeErr := file.Close()
		if writeErr != nil || syncErr != nil || closeErr != nil {
			return fmt.Errorf("append final newline to %s", relative)
		}
	}
	return nil
}

func auditCandidate(
	baseline, candidate model.Inventory,
	contract FixContract,
) ([]Change, []string) {
	allowed := make(map[string]bool, len(contract.AllowedPaths))
	for _, path := range contract.AllowedPaths {
		allowed[path] = true
	}
	before := make(map[string]model.FileRecord, len(baseline.Files))
	after := make(map[string]model.FileRecord, len(candidate.Files))
	for _, record := range baseline.Files {
		before[record.Path] = record
	}
	for _, record := range candidate.Files {
		after[record.Path] = record
	}
	changes := make([]Change, 0)
	reasons := make([]string, 0)
	changedAllowed := make(map[string]bool, len(contract.AllowedPaths))
	reasons = append(reasons, auditRawWorkspace(baseline, candidate.Root, contract.ProtectedPaths)...)
	for path, beforeRecord := range before {
		afterRecord, present := after[path]
		if !present {
			changes = append(changes, Change{
				Path: path, Kind: "deleted", BeforeSHA256: beforeRecord.SHA256, BeforeMode: beforeRecord.Mode,
			})
			reasons = append(reasons, "Candidate deleted "+path+".")
			continue
		}
		if beforeRecord.SHA256 == afterRecord.SHA256 && beforeRecord.Mode == afterRecord.Mode {
			continue
		}
		change := Change{
			Path: path, Kind: "modified", BeforeSHA256: beforeRecord.SHA256, AfterSHA256: afterRecord.SHA256,
			BeforeMode: beforeRecord.Mode, AfterMode: afterRecord.Mode,
		}
		if !allowed[path] {
			reasons = append(reasons, "Candidate changed a path outside the fix contract: "+path+".")
		} else {
			changedAllowed[path] = true
		}
		if protected(path, contract.ProtectedPaths) {
			reasons = append(reasons, "Candidate changed protected path "+path+".")
		}
		if beforeRecord.Mode != afterRecord.Mode {
			reasons = append(reasons, "Candidate changed file mode for "+path+".")
		}
		expectedDigest, err := digestWithAppendedNewline(baseline.Root, beforeRecord)
		if err != nil {
			reasons = append(reasons, err.Error()+".")
		} else if afterRecord.Size != beforeRecord.Size+1 || afterRecord.SHA256 != expectedDigest {
			reasons = append(reasons, "Candidate made a non-deterministic change to "+path+".")
		} else {
			change.AddedLines = 1
		}
		changes = append(changes, change)
	}
	for path, record := range after {
		if _, present := before[path]; present {
			continue
		}
		changes = append(changes, Change{
			Path: path, Kind: "added", AfterSHA256: record.SHA256, AfterMode: record.Mode,
		})
		reasons = append(reasons, "Candidate added unexpected path "+path+".")
	}
	sort.Slice(changes, func(left, right int) bool { return changes[left].Path < changes[right].Path })
	if len(changes) > contract.MaxFiles {
		reasons = append(reasons, fmt.Sprintf("Candidate changed %d files, above the %d-file limit.", len(changes), contract.MaxFiles))
	}
	changedLines := 0
	for _, change := range changes {
		changedLines += change.AddedLines + change.RemovedLines
	}
	if changedLines > contract.MaxChangedLines {
		reasons = append(reasons, fmt.Sprintf("Candidate changed %d lines, above the %d-line limit.", changedLines, contract.MaxChangedLines))
	}
	for _, path := range contract.AllowedPaths {
		if !changedAllowed[path] {
			reasons = append(reasons, "Candidate did not apply the contracted deterministic change to "+path+".")
		}
	}
	return changes, uniqueSorted(reasons)
}

func auditRawWorkspace(baseline model.Inventory, candidateRoot string, protectedPaths []string) []string {
	allowedFiles := make(map[string]bool, len(baseline.Files))
	allowedDirectories := map[string]bool{".": true}
	for _, record := range baseline.Files {
		allowedFiles[record.Path] = true
		parent := filepath.ToSlash(filepath.Dir(filepath.FromSlash(record.Path)))
		for parent != "." {
			allowedDirectories[parent] = true
			parent = filepath.ToSlash(filepath.Dir(filepath.FromSlash(parent)))
		}
	}
	seen := make(map[string]bool, len(baseline.Files))
	reasons := make([]string, 0)
	err := filepath.WalkDir(candidateRoot, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		relative, err := filepath.Rel(candidateRoot, path)
		if err != nil {
			return err
		}
		relative = filepath.ToSlash(relative)
		if relative == "." {
			return nil
		}
		if entry.Type()&os.ModeSymlink != 0 {
			reasons = append(reasons, "Candidate contains symlink "+relative+".")
			if protected(relative, protectedPaths) {
				reasons = append(reasons, "Candidate changed protected path "+relative+".")
			}
			if entry.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if entry.IsDir() {
			if !allowedDirectories[relative] {
				reasons = append(reasons, "Candidate added unexpected directory "+relative+".")
				if protected(relative, protectedPaths) {
					reasons = append(reasons, "Candidate changed protected path "+relative+".")
				}
			}
			return nil
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			reasons = append(reasons, "Candidate contains non-regular path "+relative+".")
			if protected(relative, protectedPaths) {
				reasons = append(reasons, "Candidate changed protected path "+relative+".")
			}
			return nil
		}
		seen[relative] = true
		if !allowedFiles[relative] {
			reasons = append(reasons, "Candidate added unexpected path "+relative+".")
			if protected(relative, protectedPaths) {
				reasons = append(reasons, "Candidate changed protected path "+relative+".")
			}
		}
		return nil
	})
	if err != nil {
		reasons = append(reasons, "Candidate workspace could not be audited: "+err.Error()+".")
	}
	for path := range allowedFiles {
		if !seen[path] {
			reasons = append(reasons, "Candidate deleted "+path+".")
		}
	}
	return uniqueSorted(reasons)
}

func digestWithAppendedNewline(root string, record model.FileRecord) (string, error) {
	path, err := safeJoin(root, record.Path)
	if err != nil {
		return "", err
	}
	pathInfo, err := os.Lstat(path)
	if err != nil || !pathInfo.Mode().IsRegular() || uint32(pathInfo.Mode().Perm()) != record.Mode {
		return "", fmt.Errorf("baseline changed while auditing %s", record.Path)
	}
	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("open baseline %s", record.Path)
	}
	defer file.Close()
	openedInfo, err := file.Stat()
	if err != nil || !os.SameFile(pathInfo, openedInfo) || uint32(openedInfo.Mode().Perm()) != record.Mode {
		return "", fmt.Errorf("baseline changed while auditing %s", record.Path)
	}
	hasher := sha256.New()
	size, err := io.Copy(hasher, file)
	if err != nil || size != record.Size || hex.EncodeToString(hasher.Sum(nil)) != record.SHA256 {
		return "", fmt.Errorf("baseline changed while auditing %s", record.Path)
	}
	hasher = sha256.New()
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", fmt.Errorf("seek baseline %s", record.Path)
	}
	if _, err := io.Copy(hasher, file); err != nil {
		return "", fmt.Errorf("hash baseline %s", record.Path)
	}
	_, _ = hasher.Write([]byte{'\n'})
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func protected(path string, patterns []string) bool {
	for _, pattern := range patterns {
		if strings.HasSuffix(pattern, "/") {
			if path == strings.TrimSuffix(pattern, "/") || strings.HasPrefix(path, pattern) {
				return true
			}
		}
		if path == pattern {
			return true
		}
	}
	return false
}

func uniqueSorted(values []string) []string {
	set := map[string]bool{}
	for _, value := range values {
		set[value] = true
	}
	result := make([]string, 0, len(set))
	for value := range set {
		result = append(result, value)
	}
	sort.Strings(result)
	return result
}
