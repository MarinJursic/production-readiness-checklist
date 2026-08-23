package remediation

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
)

func auditProviderCandidate(baseline, candidate model.Inventory, expected []Change, protectedPaths []string) []string {
	before := make(map[string]model.FileRecord, len(baseline.Files))
	after := make(map[string]model.FileRecord, len(candidate.Files))
	wanted := make(map[string]Change, len(expected))
	for _, record := range baseline.Files {
		before[record.Path] = record
	}
	for _, record := range candidate.Files {
		after[record.Path] = record
	}
	for _, change := range expected {
		wanted[change.Path] = change
	}
	reasons := auditRawProviderWorkspace(baseline, candidate.Root, expected, protectedPaths)
	observed := map[string]bool{}
	for path, beforeRecord := range before {
		afterRecord, exists := after[path]
		if !exists {
			reasons = append(reasons, "Candidate deleted "+path+".")
			continue
		}
		change, expectedChange := wanted[path]
		if !expectedChange {
			if beforeRecord.SHA256 != afterRecord.SHA256 || beforeRecord.Mode != afterRecord.Mode {
				reasons = append(reasons, "Candidate changed a path outside the proposal: "+path+".")
			}
			continue
		}
		observed[path] = true
		if change.Kind != "modified" || change.BeforeSHA256 != beforeRecord.SHA256 ||
			change.AfterSHA256 != afterRecord.SHA256 || change.BeforeMode != beforeRecord.Mode ||
			change.AfterMode != afterRecord.Mode {
			reasons = append(reasons, "Candidate bytes or mode do not match the parsed proposal for "+path+".")
		}
		if protected(path, protectedPaths) {
			reasons = append(reasons, "Candidate changed protected path "+path+".")
		}
	}
	for path, afterRecord := range after {
		if _, exists := before[path]; exists {
			continue
		}
		change, expectedChange := wanted[path]
		if !expectedChange {
			reasons = append(reasons, "Candidate added unexpected path "+path+".")
			continue
		}
		observed[path] = true
		if change.Kind != "added" || change.BeforeSHA256 != "" || change.BeforeMode != 0 ||
			change.AfterSHA256 != afterRecord.SHA256 || change.AfterMode != afterRecord.Mode {
			reasons = append(reasons, "Candidate addition does not match the parsed proposal for "+path+".")
		}
		if protected(path, protectedPaths) {
			reasons = append(reasons, "Candidate changed protected path "+path+".")
		}
	}
	for _, change := range expected {
		if !observed[change.Path] {
			reasons = append(reasons, "Candidate did not contain the parsed proposal change for "+change.Path+".")
		}
	}
	return uniqueSorted(reasons)
}

func auditRawProviderWorkspace(baseline model.Inventory, candidateRoot string, expected []Change, protectedPaths []string) []string {
	allowedFiles := make(map[string]bool, len(baseline.Files)+len(expected))
	allowedDirectories := map[string]bool{".": true}
	for _, record := range baseline.Files {
		allowedFiles[record.Path] = true
	}
	for _, change := range expected {
		if change.Kind == "added" {
			allowedFiles[change.Path] = true
		}
	}
	for path := range allowedFiles {
		parent := filepath.ToSlash(filepath.Dir(filepath.FromSlash(path)))
		for parent != "." {
			allowedDirectories[parent] = true
			parent = filepath.ToSlash(filepath.Dir(filepath.FromSlash(parent)))
		}
	}
	seen := map[string]bool{}
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
			return nil
		}
		if entry.IsDir() {
			if !allowedDirectories[relative] {
				reasons = append(reasons, "Candidate added unexpected directory "+relative+".")
			}
			return nil
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			reasons = append(reasons, "Candidate contains non-regular path "+relative+".")
			return nil
		}
		seen[relative] = true
		if !allowedFiles[relative] {
			reasons = append(reasons, "Candidate added unexpected path "+relative+".")
		}
		if protected(relative, protectedPaths) {
			for _, change := range expected {
				if change.Path == relative {
					reasons = append(reasons, "Candidate changed protected path "+relative+".")
				}
			}
		}
		return nil
	})
	if err != nil {
		reasons = append(reasons, fmt.Sprintf("Candidate workspace could not be audited: %v.", err))
	}
	for path := range allowedFiles {
		if !seen[path] {
			reasons = append(reasons, "Candidate deleted "+path+".")
		}
	}
	sort.Strings(reasons)
	return uniqueSorted(reasons)
}

func proposalProtectedPaths(taskPaths []string) []string {
	paths := append([]string(nil), defaultProtectedPaths...)
	paths = append(paths, taskPaths...)
	return uniqueSorted(paths)
}
