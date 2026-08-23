package remediation

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/MarinJursic/production-readiness-checklist/scanner/inventory"
	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
)

const maximumFixAttempts = 10

type activePolicy struct {
	configuration   *ProjectConfiguration
	configurationID string
	projectID       string
	protectedPaths  []string
	configRelative  string
	maxFiles        int
	maxChangedLines int
	attempt         int
	maxAttempts     int
}

func resolvePolicy(target, profileID string, maxFiles, maxChangedLines, attempt, maxAttempts int, configuration *ProjectConfiguration) (activePolicy, error) {
	policy := activePolicy{
		configuration: configuration, maxFiles: maxFiles, maxChangedLines: maxChangedLines,
		attempt: attempt, maxAttempts: maxAttempts,
	}
	if policy.attempt == 0 {
		policy.attempt = 1
	}
	if policy.maxAttempts == 0 {
		policy.maxAttempts = 1
	}
	if configuration != nil {
		if err := configuration.Validation.Validate(); err != nil {
			return activePolicy{}, err
		}
		if strings.TrimSpace(configuration.SourcePath) == "" {
			return activePolicy{}, fmt.Errorf("configured remediation requires a configuration source path")
		}
		document := configuration.Validation.Configuration
		if !document.Remediation.Enabled {
			return activePolicy{}, fmt.Errorf("remediation is disabled by project configuration")
		}
		if profileID != document.Assessment.Profile {
			return activePolicy{}, fmt.Errorf("configured profile %s does not match selected profile %s", document.Assessment.Profile, profileID)
		}
		if policy.maxFiles == 0 {
			policy.maxFiles = min(document.Remediation.MaxFiles, maximumFixFiles)
		}
		if policy.maxChangedLines == 0 {
			policy.maxChangedLines = min(document.Remediation.MaxChangedLines, maximumFixLines)
		}
		if maxAttempts == 0 {
			policy.maxAttempts = document.Remediation.MaxAttempts
		}
		if policy.maxFiles > document.Remediation.MaxFiles ||
			policy.maxChangedLines > document.Remediation.MaxChangedLines ||
			policy.maxAttempts > document.Remediation.MaxAttempts {
			return activePolicy{}, fmt.Errorf("command remediation budget exceeds project configuration")
		}
		policy.configurationID = configuration.Validation.Digest
		policy.projectID = document.Project.ID
	}
	if policy.maxFiles < 1 || policy.maxFiles > maximumFixFiles {
		return activePolicy{}, fmt.Errorf("max files must be between 1 and %d", maximumFixFiles)
	}
	if policy.maxChangedLines < 1 || policy.maxChangedLines > maximumFixLines {
		return activePolicy{}, fmt.Errorf("max changed lines must be between 1 and %d", maximumFixLines)
	}
	if policy.maxAttempts < 1 || policy.maxAttempts > maximumFixAttempts || policy.attempt < 1 || policy.attempt > policy.maxAttempts {
		return activePolicy{}, fmt.Errorf("attempt must be within the 1-%d configured attempt budget", maximumFixAttempts)
	}
	protectedPaths, relative, err := requiredProtectedPaths(target, configuration)
	if err != nil {
		return activePolicy{}, err
	}
	policy.protectedPaths = protectedPaths
	policy.configRelative = relative
	return policy, nil
}

// RequiredProtectedPaths returns the scanner defaults, configured policy paths,
// and the in-target configuration source itself as one sorted immutable guard.
func RequiredProtectedPaths(target string, configuration *ProjectConfiguration) ([]string, error) {
	paths, _, err := requiredProtectedPaths(target, configuration)
	return paths, err
}

func requiredProtectedPaths(target string, configuration *ProjectConfiguration) ([]string, string, error) {
	paths := append([]string{}, defaultProtectedPaths...)
	if configuration == nil {
		return uniqueSorted(paths), "", nil
	}
	if err := configuration.Validation.Validate(); err != nil {
		return nil, "", err
	}
	paths = append(paths, configuration.Validation.Configuration.Remediation.ProtectedPaths...)
	root, err := filepath.Abs(target)
	if err != nil {
		return nil, "", fmt.Errorf("resolve remediation target: %w", err)
	}
	if resolved, resolveErr := filepath.EvalSymlinks(root); resolveErr == nil {
		root = resolved
	}
	source, err := filepath.Abs(configuration.SourcePath)
	if err != nil {
		return nil, "", fmt.Errorf("resolve configuration source: %w", err)
	}
	if resolved, resolveErr := filepath.EvalSymlinks(source); resolveErr == nil {
		source = resolved
	}
	relative, err := filepath.Rel(root, source)
	if err != nil {
		return nil, "", fmt.Errorf("relate configuration source to target: %w", err)
	}
	relative = filepath.ToSlash(relative)
	if relative != ".." && !strings.HasPrefix(relative, "../") {
		paths = append(paths, relative)
	} else {
		relative = ""
	}
	sort.Strings(paths)
	return uniqueSorted(paths), relative, nil
}

func (policy activePolicy) bind(item model.Inventory, candidateRoot string) (model.Inventory, error) {
	if policy.configuration == nil {
		return item, nil
	}
	sourcePath := policy.configuration.SourcePath
	if policy.configRelative != "" && candidateRoot != "" {
		sourcePath = filepath.Join(candidateRoot, filepath.FromSlash(policy.configRelative))
	}
	return inventory.BindConfiguration(item, policy.configuration.Validation, sourcePath)
}
