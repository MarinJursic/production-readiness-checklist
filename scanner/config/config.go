package config

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
	"slices"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	Schema             = "prc.config/v0.1"
	ValidationSchema   = "prc.config-validation/v0.1"
	maximumConfigBytes = 1024 * 1024
)

var (
	identifierPattern = regexp.MustCompile(`^[a-z0-9](?:[a-z0-9-]{0,62})$`)
	profilePattern    = regexp.MustCompile(`^prc/[a-z0-9-]+$`)
	factKeyPattern    = regexp.MustCompile(`^[a-z][a-z0-9-]{0,62}$`)
	digestPattern     = regexp.MustCompile(`^sha256:[0-9a-f]{64}$`)
	sourceRefPattern  = regexp.MustCompile(`^[0-9a-f]{40,64}$`)
)

type Document struct {
	SchemaVersion string          `yaml:"schema_version" json:"schema_version"`
	Project       Project         `yaml:"project" json:"project"`
	Assessment    Assessment      `yaml:"assessment" json:"assessment"`
	Components    Components      `yaml:"components" json:"components"`
	Features      map[string]bool `yaml:"features" json:"features"`
	Data          Data            `yaml:"data" json:"data"`
	Execution     Execution       `yaml:"execution" json:"execution"`
	Remediation   Remediation     `yaml:"remediation" json:"remediation"`
}

type Project struct {
	ID          string `yaml:"id" json:"id"`
	Name        string `yaml:"name" json:"name"`
	RiskProfile string `yaml:"risk_profile" json:"risk_profile"`
}

type Assessment struct {
	Profile            string   `yaml:"profile" json:"profile"`
	SourceRef          string   `yaml:"source_ref" json:"source_ref"`
	ArtifactDigests    []string `yaml:"artifact_digests" json:"artifact_digests"`
	TargetEnvironments []string `yaml:"target_environments" json:"target_environments"`
}

type Component struct {
	Path string `yaml:"path" json:"path"`
	Type string `yaml:"type" json:"type"`
}

type Exclusion struct {
	Path      string `yaml:"path" json:"path"`
	Rationale string `yaml:"rationale" json:"rationale"`
}

type Components struct {
	Include []Component `yaml:"include" json:"include"`
	Exclude []Exclusion `yaml:"exclude" json:"exclude"`
}

type Data struct {
	Classifications      []string `yaml:"classifications" json:"classifications"`
	Regulated            []string `yaml:"regulated" json:"regulated"`
	ProhibitedInEvidence []string `yaml:"prohibited_in_evidence" json:"prohibited_in_evidence"`
}

type Execution struct {
	Network             string     `yaml:"network" json:"network"`
	AllowCommands       [][]string `yaml:"allow_commands" json:"allow_commands"`
	ProductionConnected bool       `yaml:"production_connected" json:"production_connected"`
	MaxParallel         int        `yaml:"max_parallel" json:"max_parallel"`
	MaxDurationSeconds  int        `yaml:"max_duration_seconds" json:"max_duration_seconds"`
}

type Remediation struct {
	Enabled         bool     `yaml:"enabled" json:"enabled"`
	MaxAttempts     int      `yaml:"max_attempts" json:"max_attempts"`
	MaxFiles        int      `yaml:"max_files" json:"max_files"`
	MaxChangedLines int      `yaml:"max_changed_lines" json:"max_changed_lines"`
	ProtectedPaths  []string `yaml:"protected_paths" json:"protected_paths"`
}

type Validation struct {
	SchemaVersion string   `json:"schema_version"`
	Digest        string   `json:"digest"`
	SourceSHA256  string   `json:"source_sha256"`
	Configuration Document `json:"configuration"`
}

func Load(path string) (Validation, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return Validation{}, fmt.Errorf("inspect project configuration: %w", err)
	}
	if !info.Mode().IsRegular() || info.Size() > maximumConfigBytes {
		return Validation{}, fmt.Errorf("project configuration must be a regular file no larger than 1 MiB")
	}
	file, err := os.Open(path)
	if err != nil {
		return Validation{}, fmt.Errorf("open project configuration: %w", err)
	}
	defer file.Close()
	opened, err := file.Stat()
	if err != nil || !opened.Mode().IsRegular() || !os.SameFile(info, opened) {
		return Validation{}, fmt.Errorf("project configuration changed while opening")
	}
	data, err := io.ReadAll(io.LimitReader(file, maximumConfigBytes+1))
	if err != nil || len(data) > maximumConfigBytes {
		return Validation{}, fmt.Errorf("read project configuration: file exceeds 1 MiB")
	}
	var document Document
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)
	if err := decoder.Decode(&document); err != nil {
		return Validation{}, fmt.Errorf("decode project configuration: %w", err)
	}
	var trailing any
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		return Validation{}, fmt.Errorf("project configuration contains multiple YAML documents or trailing data")
	}
	normalize(&document)
	if err := document.Validate(); err != nil {
		return Validation{}, err
	}
	payload, err := json.Marshal(document)
	if err != nil {
		return Validation{}, fmt.Errorf("encode canonical project configuration: %w", err)
	}
	digest := sha256.Sum256(payload)
	sourceDigest := sha256.Sum256(data)
	return Validation{
		SchemaVersion: ValidationSchema,
		Digest:        hex.EncodeToString(digest[:]),
		SourceSHA256:  hex.EncodeToString(sourceDigest[:]),
		Configuration: document,
	}, nil
}

func normalize(document *Document) {
	if document.Features == nil {
		document.Features = map[string]bool{}
	}
	if document.Assessment.ArtifactDigests == nil {
		document.Assessment.ArtifactDigests = []string{}
	}
	if document.Assessment.TargetEnvironments == nil {
		document.Assessment.TargetEnvironments = []string{}
	}
	if document.Components.Include == nil {
		document.Components.Include = []Component{}
	}
	if document.Components.Exclude == nil {
		document.Components.Exclude = []Exclusion{}
	}
	if document.Data.Classifications == nil {
		document.Data.Classifications = []string{}
	}
	if document.Data.Regulated == nil {
		document.Data.Regulated = []string{}
	}
	if document.Data.ProhibitedInEvidence == nil {
		document.Data.ProhibitedInEvidence = []string{}
	}
	if document.Execution.AllowCommands == nil {
		document.Execution.AllowCommands = [][]string{}
	}
	if document.Remediation.ProtectedPaths == nil {
		document.Remediation.ProtectedPaths = []string{}
	}
}

func (document Document) Validate() error {
	if document.SchemaVersion != Schema {
		return fmt.Errorf("unsupported project configuration schema %q", document.SchemaVersion)
	}
	if !identifierPattern.MatchString(document.Project.ID) || strings.TrimSpace(document.Project.Name) == "" || len(document.Project.Name) > 200 {
		return fmt.Errorf("project requires a valid id and name")
	}
	if !oneOf(document.Project.RiskProfile, "low", "medium", "high", "critical") {
		return fmt.Errorf("project risk_profile must be low, medium, high, or critical")
	}
	if !profilePattern.MatchString(document.Assessment.Profile) {
		return fmt.Errorf("assessment requires a valid profile ID")
	}
	if document.Assessment.SourceRef != "" && !sourceRefPattern.MatchString(document.Assessment.SourceRef) {
		return fmt.Errorf("assessment source_ref must be empty or an exact 40-64 character lowercase hexadecimal revision")
	}
	if err := validateSortedUnique(document.Assessment.ArtifactDigests, "artifact digests", func(value string) bool {
		return digestPattern.MatchString(value)
	}); err != nil {
		return err
	}
	if err := validateNames(document.Assessment.TargetEnvironments, "target environments"); err != nil {
		return err
	}
	if err := validateComponents(document.Components); err != nil {
		return err
	}
	if len(document.Features) > 1000 {
		return fmt.Errorf("features exceed 1000 entries")
	}
	for key := range document.Features {
		if !factKeyPattern.MatchString(key) {
			return fmt.Errorf("invalid feature key %q", key)
		}
	}
	if err := validateNames(document.Data.Classifications, "data classifications"); err != nil {
		return err
	}
	if err := validateNames(document.Data.Regulated, "regulated data tags"); err != nil {
		return err
	}
	if err := validateNames(document.Data.ProhibitedInEvidence, "prohibited evidence classes"); err != nil {
		return err
	}
	if document.Execution.Network != "deny" || document.Execution.ProductionConnected || len(document.Execution.AllowCommands) != 0 {
		return fmt.Errorf("configuration v0.1 requires network deny, production_connected false, and no allowed commands")
	}
	if document.Execution.MaxParallel < 1 || document.Execution.MaxParallel > 32 ||
		document.Execution.MaxDurationSeconds < 1 || document.Execution.MaxDurationSeconds > 86400 {
		return fmt.Errorf("execution limits must use 1-32 workers and 1-86400 seconds")
	}
	if document.Remediation.MaxAttempts < 1 || document.Remediation.MaxAttempts > 10 ||
		document.Remediation.MaxFiles < 1 || document.Remediation.MaxFiles > 1000 ||
		document.Remediation.MaxChangedLines < 1 || document.Remediation.MaxChangedLines > 10000 {
		return fmt.Errorf("remediation budgets are outside supported limits")
	}
	if len(document.Remediation.ProtectedPaths) == 0 {
		return fmt.Errorf("remediation requires at least one protected path")
	}
	if err := validatePaths(document.Remediation.ProtectedPaths, "protected paths"); err != nil {
		return err
	}
	return nil
}

func validateComponents(components Components) error {
	if len(components.Include) > 1000 || len(components.Exclude) > 1000 {
		return fmt.Errorf("component declarations exceed 1000 entries")
	}
	seenPaths := map[string]bool{}
	for _, component := range components.Include {
		if err := validatePath(component.Path); err != nil {
			return fmt.Errorf("included component: %w", err)
		}
		if !oneOf(component.Type, "repository", "web-frontend", "http-api", "service", "worker", "job", "library", "cli", "infrastructure", "container") {
			return fmt.Errorf("included component %s has unsupported type %q", component.Path, component.Type)
		}
		if seenPaths[component.Path] {
			return fmt.Errorf("duplicate component path %q", component.Path)
		}
		seenPaths[component.Path] = true
	}
	for _, exclusion := range components.Exclude {
		if err := validatePath(exclusion.Path); err != nil {
			return fmt.Errorf("excluded component: %w", err)
		}
		if strings.TrimSpace(exclusion.Rationale) == "" || len(exclusion.Rationale) > 1000 {
			return fmt.Errorf("excluded component %s requires a bounded rationale", exclusion.Path)
		}
		if seenPaths[exclusion.Path] {
			return fmt.Errorf("duplicate component path %q", exclusion.Path)
		}
		seenPaths[exclusion.Path] = true
	}
	if !slices.IsSortedFunc(components.Include, func(left, right Component) int { return strings.Compare(left.Path, right.Path) }) ||
		!slices.IsSortedFunc(components.Exclude, func(left, right Exclusion) int { return strings.Compare(left.Path, right.Path) }) {
		return fmt.Errorf("component include and exclude entries must be sorted by path")
	}
	return nil
}

func validateNames(values []string, name string) error {
	return validateSortedUnique(values, name, func(value string) bool { return factKeyPattern.MatchString(value) })
}

func validatePaths(values []string, name string) error {
	return validateSortedUnique(values, name, func(value string) bool { return validatePath(value) == nil })
}

func validateSortedUnique(values []string, name string, valid func(string) bool) error {
	if len(values) > 10_000 {
		return fmt.Errorf("%s exceed 10000 entries", name)
	}
	if !slices.IsSorted(values) {
		return fmt.Errorf("%s must be sorted", name)
	}
	for index, value := range values {
		if !valid(value) {
			return fmt.Errorf("invalid %s value %q", name, value)
		}
		if index > 0 && values[index-1] == value {
			return fmt.Errorf("duplicate %s value %q", name, value)
		}
	}
	return nil
}

func validatePath(path string) error {
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

func oneOf(value string, allowed ...string) bool {
	return slices.Contains(allowed, value)
}
