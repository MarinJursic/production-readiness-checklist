package catalog

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
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
	"gopkg.in/yaml.v3"
)

const maximumCatalogFileBytes = 4 * 1024 * 1024

var (
	catalogVersionPattern = regexp.MustCompile(`^[0-9]+\.[0-9]+\.[0-9]+$`)
	objectiveIDPattern    = regexp.MustCompile(`^(?:PRC-[0-9]{2}-[0-9]{3}|USEQ-[A-F0-9]{8})$`)
	assertionIDPattern    = regexp.MustCompile(`^PRC-A-[A-Z0-9]+-[0-9]{3}$`)
	implementationPattern = regexp.MustCompile(`^prc\.[a-z0-9.-]+@[0-9]+\.[0-9]+$`)
	profileIDPattern      = regexp.MustCompile(`^prc/[a-z0-9-]+$`)
	profileVersionPattern = regexp.MustCompile(`^[0-9]+\.[0-9]+$`)
	namePattern           = regexp.MustCompile(`^[a-z0-9-]+$`)
)

type Catalog struct {
	Root       string
	Version    string
	Objectives map[string]model.Objective
	Assertions map[string]model.Assertion
	Profiles   map[string]model.Profile
}

type objectiveDocument struct {
	SchemaVersion  string            `yaml:"schema_version"`
	CatalogVersion string            `yaml:"catalog_version"`
	Objectives     []model.Objective `yaml:"objectives"`
}

type assertionDocument struct {
	SchemaVersion  string            `yaml:"schema_version"`
	CatalogVersion string            `yaml:"catalog_version"`
	Assertions     []model.Assertion `yaml:"assertions"`
}

func Load(root string) (*Catalog, error) {
	abs, err := filepath.Abs(root)
	if err != nil {
		return nil, fmt.Errorf("resolve catalog root: %w", err)
	}
	c := &Catalog{
		Root:       abs,
		Objectives: map[string]model.Objective{},
		Assertions: map[string]model.Assertion{},
		Profiles:   map[string]model.Profile{},
	}
	if err := loadObjectives(c); err != nil {
		return nil, err
	}
	if err := loadAssertions(c); err != nil {
		return nil, err
	}
	if err := loadProfiles(c); err != nil {
		return nil, err
	}
	if err := c.validateReferences(); err != nil {
		return nil, err
	}
	return c, nil
}

func yamlFiles(root, directory string) ([]string, error) {
	paths, err := filepath.Glob(filepath.Join(root, "catalog", directory, "*.yaml"))
	if err != nil {
		return nil, fmt.Errorf("list %s catalog: %w", directory, err)
	}
	sort.Strings(paths)
	if len(paths) == 0 {
		return nil, fmt.Errorf("no %s catalog files under %s", directory, root)
	}
	return paths, nil
}

func decodeYAML(path string, destination any) error {
	data, err := readBoundedRegularFile(path, "catalog YAML")
	if err != nil {
		return err
	}
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)
	if err := decoder.Decode(destination); err != nil {
		return fmt.Errorf("parse %s: %w", path, err)
	}
	var trailing any
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		if err == nil {
			return fmt.Errorf("parse %s: multiple YAML documents are not allowed", path)
		}
		return fmt.Errorf("parse trailing %s content: %w", path, err)
	}
	return nil
}

func loadObjectives(c *Catalog) error {
	paths, err := yamlFiles(c.Root, "objectives")
	if err != nil {
		return err
	}
	for _, path := range paths {
		var document objectiveDocument
		if err := decodeYAML(path, &document); err != nil {
			return err
		}
		if document.SchemaVersion != "prc.objectives/v0.1" {
			return fmt.Errorf("unsupported objective schema %q in %s", document.SchemaVersion, path)
		}
		if err := bindCatalogVersion(c, document.CatalogVersion, path); err != nil {
			return err
		}
		if len(document.Objectives) == 0 {
			return fmt.Errorf("objective catalog %s contains no objectives", path)
		}
		for _, objective := range document.Objectives {
			if err := validateObjective(c.Root, objective); err != nil {
				return fmt.Errorf("invalid objective in %s: %w", path, err)
			}
			if _, exists := c.Objectives[objective.ID]; exists {
				return fmt.Errorf("duplicate objective ID %s", objective.ID)
			}
			c.Objectives[objective.ID] = objective
		}
	}
	return nil
}

func loadAssertions(c *Catalog) error {
	paths, err := yamlFiles(c.Root, "assertions")
	if err != nil {
		return err
	}
	for _, path := range paths {
		var document assertionDocument
		if err := decodeYAML(path, &document); err != nil {
			return err
		}
		if document.SchemaVersion != "prc.assertions/v0.1" {
			return fmt.Errorf("unsupported assertion schema %q in %s", document.SchemaVersion, path)
		}
		if err := bindCatalogVersion(c, document.CatalogVersion, path); err != nil {
			return err
		}
		if len(document.Assertions) == 0 {
			return fmt.Errorf("assertion catalog %s contains no assertions", path)
		}
		for _, assertion := range document.Assertions {
			if err := validateAssertion(assertion); err != nil {
				return fmt.Errorf("invalid assertion in %s: %w", path, err)
			}
			if _, exists := c.Assertions[assertion.ID]; exists {
				return fmt.Errorf("duplicate assertion ID %s", assertion.ID)
			}
			c.Assertions[assertion.ID] = assertion
		}
	}
	return nil
}

func loadProfiles(c *Catalog) error {
	paths, err := yamlFiles(c.Root, "profiles")
	if err != nil {
		return err
	}
	for _, path := range paths {
		var profile model.Profile
		if err := decodeYAML(path, &profile); err != nil {
			return err
		}
		if err := validateProfile(profile); err != nil {
			return fmt.Errorf("invalid profile in %s: %w", path, err)
		}
		if _, exists := c.Profiles[profile.ID]; exists {
			return fmt.Errorf("duplicate profile ID %s", profile.ID)
		}
		c.Profiles[profile.ID] = profile
	}
	return nil
}

func (c *Catalog) validateReferences() error {
	for _, objectiveID := range sortedKeys(c.Objectives) {
		objective := c.Objectives[objectiveID]
		if err := uniqueReferences(objective.AssertionIDs, "objective "+objective.ID+" assertion"); err != nil {
			return err
		}
		for _, assertionID := range objective.AssertionIDs {
			linked, ok := c.Assertions[assertionID]
			if !ok {
				return fmt.Errorf("objective %s references missing assertion %s", objective.ID, assertionID)
			}
			if !slices.Contains(linked.ControlIDs, objective.ID) {
				return fmt.Errorf("assertion %s does not map back to objective %s", assertionID, objective.ID)
			}
		}
	}
	for _, assertionID := range sortedKeys(c.Assertions) {
		assertion := c.Assertions[assertionID]
		if err := uniqueReferences(assertion.ControlIDs, "assertion "+assertion.ID+" control"); err != nil {
			return err
		}
		for _, controlID := range assertion.ControlIDs {
			objective, ok := c.Objectives[controlID]
			if !ok {
				return fmt.Errorf("assertion %s references unavailable objective %s", assertion.ID, controlID)
			}
			if !slices.Contains(objective.AssertionIDs, assertion.ID) {
				return fmt.Errorf("objective %s does not map back to assertion %s", controlID, assertion.ID)
			}
		}
	}
	for _, profileID := range sortedKeys(c.Profiles) {
		profile := c.Profiles[profileID]
		if len(profile.AssertionIDs) == 0 {
			return fmt.Errorf("profile %s contains no assertions", profile.ID)
		}
		if err := uniqueReferences(profile.AssertionIDs, "profile "+profile.ID+" assertion"); err != nil {
			return err
		}
		if err := uniqueReferences(profile.TerminalPolicy.BlockOn, "profile "+profile.ID+" terminal severity"); err != nil {
			return err
		}
		for _, severity := range profile.TerminalPolicy.BlockOn {
			if severity != "critical" && severity != "high" && severity != "medium" && severity != "low" && severity != "info" {
				return fmt.Errorf("profile %s has invalid terminal severity %q", profile.ID, severity)
			}
		}
		for _, assertionID := range profile.AssertionIDs {
			if _, ok := c.Assertions[assertionID]; !ok {
				return fmt.Errorf("profile %s references missing assertion %s", profile.ID, assertionID)
			}
		}
	}
	return nil
}

func bindCatalogVersion(c *Catalog, version, path string) error {
	if !catalogVersionPattern.MatchString(version) {
		return fmt.Errorf("catalog version %q in %s is not semantic version X.Y.Z", version, path)
	}
	if c.Version != "" && c.Version != version {
		return fmt.Errorf("catalog version %q in %s does not match %q", version, path, c.Version)
	}
	c.Version = version
	return nil
}

func uniqueReferences(values []string, label string) error {
	seen := map[string]bool{}
	for _, value := range values {
		if value == "" {
			return fmt.Errorf("%s reference is empty", label)
		}
		if seen[value] {
			return fmt.Errorf("duplicate %s reference %s", label, value)
		}
		seen[value] = true
	}
	return nil
}

func validateObjective(root string, objective model.Objective) error {
	if !objectiveIDPattern.MatchString(objective.ID) || objective.Revision < 1 {
		return fmt.Errorf("objective ID or revision is invalid")
	}
	if err := validateText(objective.Title, 512, "objective title"); err != nil {
		return err
	}
	if err := validateText(objective.Statement, 16*1024, "objective statement"); err != nil {
		return err
	}
	if objective.Source.Line < 1 || filepath.IsAbs(objective.Source.Path) || strings.Contains(objective.Source.Path, "\\") ||
		!strings.HasPrefix(objective.Source.Path, "docs/") || !strings.HasSuffix(objective.Source.Path, ".md") ||
		filepath.ToSlash(filepath.Clean(filepath.FromSlash(objective.Source.Path))) != objective.Source.Path {
		return fmt.Errorf("objective source is not a safe Markdown path and positive line")
	}
	if len(objective.Domains) == 0 {
		return fmt.Errorf("objective requires at least one domain")
	}
	if err := validatePatternValues(objective.Domains, namePattern, "objective domain"); err != nil {
		return err
	}
	switch objective.AutomationClass {
	case "automated", "partial", "manual", "external", "contextual", "unmapped":
	default:
		return fmt.Errorf("objective automation class %q is invalid", objective.AutomationClass)
	}
	if err := validatePatternValues(objective.AssertionIDs, assertionIDPattern, "objective assertion ID"); err != nil {
		return err
	}
	data, err := readBoundedRegularFile(filepath.Join(root, filepath.FromSlash(objective.Source.Path)), "objective source")
	if err != nil {
		return err
	}
	lines := strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n")
	if objective.Source.Line > len(lines) {
		return fmt.Errorf("objective source line %d is outside %s", objective.Source.Line, objective.Source.Path)
	}
	want := "- [ ] **" + objective.ID + "** — " + objective.Statement
	if lines[objective.Source.Line-1] != want {
		return fmt.Errorf("objective source %s:%d does not exactly match its ID and statement", objective.Source.Path, objective.Source.Line)
	}
	return nil
}

func validateAssertion(assertion model.Assertion) error {
	if !assertionIDPattern.MatchString(assertion.ID) || assertion.Revision < 1 {
		return fmt.Errorf("assertion ID or revision is invalid")
	}
	if len(assertion.ControlIDs) == 0 {
		return fmt.Errorf("assertion requires at least one control")
	}
	if err := validatePatternValues(assertion.ControlIDs, objectiveIDPattern, "assertion control ID"); err != nil {
		return err
	}
	if err := validateText(assertion.Title, 512, "assertion title"); err != nil {
		return err
	}
	if err := validateText(assertion.Statement, 16*1024, "assertion statement"); err != nil {
		return err
	}
	if err := validateText(assertion.Applicability, 4096, "assertion applicability"); err != nil {
		return err
	}
	if len(assertion.EvidenceRequired) == 0 {
		return fmt.Errorf("assertion requires at least one evidence requirement")
	}
	seenEvidence := map[string]bool{}
	for _, requirement := range assertion.EvidenceRequired {
		if !namePattern.MatchString(requirement.Kind) {
			return fmt.Errorf("evidence kind %q is invalid", requirement.Kind)
		}
		switch requirement.MinimumAuthority {
		case "declared", "repository", "executed", "artifact", "environment", "human":
		default:
			return fmt.Errorf("evidence authority %q is invalid", requirement.MinimumAuthority)
		}
		if err := validateText(requirement.Description, 16*1024, "evidence description"); err != nil {
			return err
		}
		key := requirement.Kind + "\x00" + requirement.MinimumAuthority + "\x00" + requirement.Description
		if seenEvidence[key] {
			return fmt.Errorf("assertion repeats an evidence requirement")
		}
		seenEvidence[key] = true
	}
	if !implementationPattern.MatchString(assertion.ImplementationID) {
		return fmt.Errorf("implementation ID %q is invalid", assertion.ImplementationID)
	}
	parameters, err := json.Marshal(assertion.Parameters)
	if err != nil || len(parameters) > 64*1024 {
		return fmt.Errorf("assertion parameters must be canonical JSON no larger than 64 KiB")
	}
	if !oneOf(assertion.Severity, "info", "low", "medium", "high", "critical") {
		return fmt.Errorf("assertion severity %q is invalid", assertion.Severity)
	}
	if !oneOf(assertion.Gate, "advisory", "required", "no-go") {
		return fmt.Errorf("assertion gate %q is invalid", assertion.Gate)
	}
	if !oneOf(assertion.RemediationClass, "R0", "R1", "R2", "R3", "R4", "R5", "R6") {
		return fmt.Errorf("assertion remediation class %q is invalid", assertion.RemediationClass)
	}
	return nil
}

func validateProfile(profile model.Profile) error {
	if profile.SchemaVersion != "prc.profile/v0.1" || !profileIDPattern.MatchString(profile.ID) ||
		!profileVersionPattern.MatchString(profile.Version) {
		return fmt.Errorf("profile schema, ID, or version is invalid")
	}
	if err := validateText(profile.Title, 512, "profile title"); err != nil {
		return err
	}
	if err := validateText(profile.Description, 16*1024, "profile description"); err != nil {
		return err
	}
	if len(profile.AssertionIDs) == 0 {
		return fmt.Errorf("profile requires at least one assertion")
	}
	if err := validatePatternValues(profile.AssertionIDs, assertionIDPattern, "profile assertion ID"); err != nil {
		return err
	}
	if err := validatePatternValues(profile.TerminalPolicy.BlockOn, namePattern, "profile terminal severity"); err != nil {
		return err
	}
	for _, severity := range profile.TerminalPolicy.BlockOn {
		if !oneOf(severity, "info", "low", "medium", "high", "critical") {
			return fmt.Errorf("profile terminal severity %q is invalid", severity)
		}
	}
	return nil
}

func validateText(value string, maximum int, label string) error {
	if strings.TrimSpace(value) == "" || len(value) > maximum || !utf8.ValidString(value) || strings.ContainsRune(value, '\x00') {
		return fmt.Errorf("%s must be nonempty UTF-8 text no longer than %d bytes", label, maximum)
	}
	return nil
}

func validatePatternValues(values []string, pattern *regexp.Regexp, label string) error {
	seen := map[string]bool{}
	for _, value := range values {
		if !pattern.MatchString(value) {
			return fmt.Errorf("invalid %s %q", label, value)
		}
		if seen[value] {
			return fmt.Errorf("duplicate %s %q", label, value)
		}
		seen[value] = true
	}
	return nil
}

func oneOf(value string, allowed ...string) bool {
	return slices.Contains(allowed, value)
}

func sortedKeys[T any](values map[string]T) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func readBoundedRegularFile(path, label string) ([]byte, error) {
	before, err := os.Lstat(path)
	if err != nil {
		return nil, fmt.Errorf("inspect %s %s: %w", label, path, err)
	}
	if !before.Mode().IsRegular() || before.Size() > maximumCatalogFileBytes {
		return nil, fmt.Errorf("%s %s must be a regular file no larger than %d bytes", label, path, maximumCatalogFileBytes)
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open %s %s: %w", label, path, err)
	}
	defer file.Close()
	opened, err := file.Stat()
	if err != nil || !opened.Mode().IsRegular() || !os.SameFile(before, opened) {
		return nil, fmt.Errorf("%s %s changed while opening", label, path)
	}
	data, err := io.ReadAll(io.LimitReader(file, maximumCatalogFileBytes+1))
	if err != nil {
		return nil, fmt.Errorf("read %s %s: %w", label, path, err)
	}
	if len(data) > maximumCatalogFileBytes || int64(len(data)) != opened.Size() {
		return nil, fmt.Errorf("%s %s changed or exceeded its byte limit", label, path)
	}
	after, err := file.Stat()
	if err != nil || !after.Mode().IsRegular() || !os.SameFile(opened, after) || after.Size() != opened.Size() {
		return nil, fmt.Errorf("%s %s changed while reading", label, path)
	}
	return data, nil
}

// Digest returns a deterministic identity for every parsed governing catalog
// definition. The filesystem root is intentionally excluded so identical
// catalogs have the same identity in different workspaces.
func (c *Catalog) Digest() (string, error) {
	payload, err := json.Marshal(struct {
		SchemaVersion string                     `json:"schema_version"`
		Version       string                     `json:"catalog_version"`
		Objectives    map[string]model.Objective `json:"objectives"`
		Assertions    map[string]model.Assertion `json:"assertions"`
		Profiles      map[string]model.Profile   `json:"profiles"`
	}{
		SchemaVersion: "prc.catalog-identity/v0.1", Version: c.Version,
		Objectives: c.Objectives, Assertions: c.Assertions, Profiles: c.Profiles,
	})
	if err != nil {
		return "", fmt.Errorf("encode catalog identity: %w", err)
	}
	digest := sha256.Sum256(payload)
	return hex.EncodeToString(digest[:]), nil
}

func (c *Catalog) Profile(id string) (model.Profile, error) {
	profile, ok := c.Profiles[id]
	if !ok {
		return model.Profile{}, fmt.Errorf("unknown profile %q", id)
	}
	return profile, nil
}
