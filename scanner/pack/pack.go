// Package pack validates distributable groups of catalog assertions against
// their exact implementations and benchmark evidence.
package pack

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
	"sort"
	"strings"
	"time"

	"github.com/MarinJursic/production-readiness-checklist/scanner/benchmark"
	"github.com/MarinJursic/production-readiness-checklist/scanner/catalog"
	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
	"gopkg.in/yaml.v3"
)

const (
	Schema       = "prc.pack/v0.1"
	ReportSchema = "prc.pack-report/v0.1"
)

const maximumPackBytes = 1024 * 1024

var (
	packIDPattern      = regexp.MustCompile(`^prc\.pack\.[a-z0-9.-]+@[0-9]+\.[0-9]+$`)
	publisherIDPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9.-]{0,127}$`)
	engineAPIPattern   = regexp.MustCompile(`^prc\.engine/v[0-9]+\.[0-9]+$`)
)

type AssertionBinding struct {
	AssertionID       string   `json:"assertion_id" yaml:"assertion_id"`
	ImplementationID  string   `json:"implementation_id" yaml:"implementation_id"`
	ValidatedOutcomes []string `json:"validated_outcomes" yaml:"validated_outcomes"`
}

type BenchmarkBinding struct {
	SuiteID     string `json:"suite_id" yaml:"suite_id"`
	SuitePath   string `json:"suite_path" yaml:"suite_path"`
	SuiteSHA256 string `json:"suite_sha256" yaml:"suite_sha256"`
}

type Manifest struct {
	SchemaVersion string             `json:"schema_version" yaml:"schema_version"`
	ID            string             `json:"id" yaml:"id"`
	Title         string             `json:"title" yaml:"title"`
	Description   string             `json:"description" yaml:"description"`
	PublisherID   string             `json:"publisher_id" yaml:"publisher_id"`
	Maintenance   string             `json:"maintenance" yaml:"maintenance"`
	EngineAPIs    []string           `json:"engine_apis" yaml:"engine_apis"`
	ProfileIDs    []string           `json:"profile_ids" yaml:"profile_ids"`
	Assertions    []AssertionBinding `json:"assertions" yaml:"assertions"`
	Benchmark     BenchmarkBinding   `json:"benchmark" yaml:"benchmark"`
	Limitations   []string           `json:"limitations" yaml:"limitations"`
}

type Loaded struct {
	Manifest              Manifest
	Digest                string
	SuiteDigest           string
	BenchmarkCorpusDigest string
	CatalogDigest         string
}

type Report struct {
	SchemaVersion         string   `json:"schema_version"`
	Manifest              Manifest `json:"manifest"`
	Digest                string   `json:"digest"`
	SuiteDigest           string   `json:"suite_digest"`
	BenchmarkCorpusDigest string   `json:"benchmark_corpus_digest"`
	CatalogDigest         string   `json:"catalog_digest"`
}

func (loaded Loaded) Report() Report {
	return Report{
		SchemaVersion: ReportSchema, Manifest: loaded.Manifest, Digest: loaded.Digest,
		SuiteDigest: loaded.SuiteDigest, BenchmarkCorpusDigest: loaded.BenchmarkCorpusDigest,
		CatalogDigest: loaded.CatalogDigest,
	}
}

func Load(root, manifestPath string, catalogValue *catalog.Catalog) (Loaded, error) {
	resolvedRoot, data, err := readManifest(root, manifestPath)
	if err != nil {
		return Loaded{}, err
	}
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)
	var manifest Manifest
	if err := decoder.Decode(&manifest); err != nil {
		return Loaded{}, fmt.Errorf("decode pack manifest: %w", err)
	}
	var trailing any
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		if err == nil {
			return Loaded{}, fmt.Errorf("pack manifest contains more than one YAML document")
		}
		return Loaded{}, fmt.Errorf("decode trailing pack manifest content: %w", err)
	}
	if err := validateAndNormalize(&manifest, catalogValue); err != nil {
		return Loaded{}, err
	}
	suitePath, err := resolveRootFile(resolvedRoot, manifest.Benchmark.SuitePath)
	if err != nil {
		return Loaded{}, fmt.Errorf("resolve pack benchmark suite: %w", err)
	}
	loadedSuite, err := benchmark.Load(suitePath, catalogValue)
	if err != nil {
		return Loaded{}, fmt.Errorf("load pack benchmark suite: %w", err)
	}
	if loadedSuite.Suite.ID != manifest.Benchmark.SuiteID || loadedSuite.Digest != manifest.Benchmark.SuiteSHA256 {
		return Loaded{}, fmt.Errorf("pack benchmark suite does not match its pinned ID or digest")
	}
	if !contains(manifest.ProfileIDs, loadedSuite.Suite.ProfileID) {
		return Loaded{}, fmt.Errorf("pack benchmark profile %s is not declared by the pack", loadedSuite.Suite.ProfileID)
	}
	if err := validateOutcomeCoverage(manifest.Assertions, loadedSuite.Suite); err != nil {
		return Loaded{}, err
	}
	benchmarkReport, err := benchmark.Evaluate(
		catalogValue, suitePath, time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
	)
	if err != nil {
		return Loaded{}, fmt.Errorf("evaluate pack benchmark: %w", err)
	}
	if !benchmarkReport.Passed {
		return Loaded{}, fmt.Errorf("pack benchmark quality budget failed: %s", strings.Join(benchmarkReport.QualityFailures, "; "))
	}
	payload, err := json.Marshal(manifest)
	if err != nil {
		return Loaded{}, fmt.Errorf("encode pack identity: %w", err)
	}
	digest := sha256.Sum256(payload)
	catalogDigest, err := catalogValue.Digest()
	if err != nil {
		return Loaded{}, err
	}
	return Loaded{
		Manifest: manifest, Digest: hex.EncodeToString(digest[:]), SuiteDigest: loadedSuite.Digest,
		BenchmarkCorpusDigest: benchmarkReport.CorpusDigest,
		CatalogDigest:         catalogDigest,
	}, nil
}

func validateAndNormalize(manifest *Manifest, catalogValue *catalog.Catalog) error {
	if manifest.SchemaVersion != Schema || !packIDPattern.MatchString(manifest.ID) ||
		strings.TrimSpace(manifest.Title) == "" || strings.TrimSpace(manifest.Description) == "" ||
		!publisherIDPattern.MatchString(manifest.PublisherID) {
		return fmt.Errorf("pack requires a supported schema, valid identity, title, description, and publisher")
	}
	if manifest.Maintenance != "active" && manifest.Maintenance != "deprecated" {
		return fmt.Errorf("pack maintenance must be active or deprecated")
	}
	if len(manifest.EngineAPIs) == 0 || len(manifest.ProfileIDs) == 0 || len(manifest.Assertions) == 0 || len(manifest.Limitations) == 0 {
		return fmt.Errorf("pack requires engine APIs, profiles, assertions, and honest limitations")
	}
	if len(manifest.EngineAPIs) > 16 || len(manifest.ProfileIDs) > 64 || len(manifest.Assertions) > 4096 || len(manifest.Limitations) > 128 {
		return fmt.Errorf("pack exceeds a collection limit")
	}
	seenEngineAPIs, supportsCurrent := map[string]bool{}, false
	for _, identifier := range manifest.EngineAPIs {
		if !engineAPIPattern.MatchString(identifier) || seenEngineAPIs[identifier] {
			return fmt.Errorf("pack has an invalid or duplicate engine API %q", identifier)
		}
		seenEngineAPIs[identifier] = true
		supportsCurrent = supportsCurrent || identifier == model.EngineVersion
	}
	if !supportsCurrent {
		return fmt.Errorf("pack does not support current engine API %s", model.EngineVersion)
	}
	sort.Strings(manifest.EngineAPIs)
	profiles := map[string]model.Profile{}
	for _, identifier := range manifest.ProfileIDs {
		if _, duplicate := profiles[identifier]; duplicate {
			return fmt.Errorf("pack has duplicate profile %s", identifier)
		}
		profile, err := catalogValue.Profile(identifier)
		if err != nil {
			return fmt.Errorf("pack profile: %w", err)
		}
		profiles[identifier] = profile
	}
	sort.Strings(manifest.ProfileIDs)
	seenAssertions := map[string]bool{}
	for index := range manifest.Assertions {
		binding := &manifest.Assertions[index]
		assertion, ok := catalogValue.Assertions[binding.AssertionID]
		if !ok || seenAssertions[binding.AssertionID] {
			return fmt.Errorf("pack has an unknown or duplicate assertion %s", binding.AssertionID)
		}
		seenAssertions[binding.AssertionID] = true
		if assertion.ImplementationID != binding.ImplementationID {
			return fmt.Errorf("pack assertion %s implementation does not match the catalog", binding.AssertionID)
		}
		profiled := false
		for _, profile := range profiles {
			profiled = profiled || contains(profile.AssertionIDs, binding.AssertionID)
		}
		if !profiled {
			return fmt.Errorf("pack assertion %s is absent from its declared profiles", binding.AssertionID)
		}
		if len(binding.ValidatedOutcomes) == 0 || len(binding.ValidatedOutcomes) > 16 {
			return fmt.Errorf("pack assertion %s requires validated outcomes", binding.AssertionID)
		}
		seenOutcomes := map[string]bool{}
		for _, outcome := range binding.ValidatedOutcomes {
			if !validOutcome(outcome) || seenOutcomes[outcome] {
				return fmt.Errorf("pack assertion %s has invalid or duplicate outcome %q", binding.AssertionID, outcome)
			}
			seenOutcomes[outcome] = true
		}
		sort.Strings(binding.ValidatedOutcomes)
	}
	sort.Slice(manifest.Assertions, func(left, right int) bool {
		return manifest.Assertions[left].AssertionID < manifest.Assertions[right].AssertionID
	})
	if !suiteIDPattern.MatchString(manifest.Benchmark.SuiteID) || !digest(manifest.Benchmark.SuiteSHA256) {
		return fmt.Errorf("pack benchmark binding has an invalid suite ID or digest")
	}
	if err := validateRelativePath(manifest.Benchmark.SuitePath); err != nil {
		return err
	}
	for index := range manifest.Limitations {
		manifest.Limitations[index] = strings.TrimSpace(manifest.Limitations[index])
		if manifest.Limitations[index] == "" {
			return fmt.Errorf("pack limitations cannot be empty")
		}
	}
	sort.Strings(manifest.Limitations)
	return nil
}

func validateOutcomeCoverage(bindings []AssertionBinding, suite benchmark.Suite) error {
	observed := map[string]map[string]bool{}
	for _, item := range suite.Cases {
		for _, expected := range item.Expectations {
			if observed[expected.AssertionID] == nil {
				observed[expected.AssertionID] = map[string]bool{}
			}
			observed[expected.AssertionID][outcomeName(expected)] = true
		}
	}
	for _, binding := range bindings {
		declared := map[string]bool{}
		for _, outcome := range binding.ValidatedOutcomes {
			declared[outcome] = true
		}
		actual := observed[binding.AssertionID]
		if len(actual) != len(declared) {
			return fmt.Errorf("pack assertion %s validated outcomes do not match benchmark coverage", binding.AssertionID)
		}
		for outcome := range declared {
			if !actual[outcome] {
				return fmt.Errorf("pack assertion %s lacks benchmark outcome %s", binding.AssertionID, outcome)
			}
		}
	}
	return nil
}

func outcomeName(expected benchmark.Expectation) string {
	if expected.Assessment == "unknown" && (expected.Execution == "blocked" || expected.Execution == "error") {
		return expected.Execution
	}
	return expected.Assessment
}

func validOutcome(value string) bool {
	return value == "pass" || value == "fail" || value == "not_applicable" || value == "manual_review" ||
		value == "blocked" || value == "error" || value == "unknown" || value == "stale" || value == "conflicting"
}

func readManifest(root, path string) (string, []byte, error) {
	resolvedRoot, err := filepath.Abs(root)
	if err != nil {
		return "", nil, fmt.Errorf("resolve pack root: %w", err)
	}
	resolvedRoot, err = filepath.EvalSymlinks(resolvedRoot)
	if err != nil {
		return "", nil, fmt.Errorf("resolve pack root: %w", err)
	}
	absolute, err := filepath.Abs(path)
	if err != nil {
		return "", nil, fmt.Errorf("resolve pack manifest: %w", err)
	}
	if err := requireInside(resolvedRoot, absolute); err != nil {
		return "", nil, err
	}
	relative, _ := filepath.Rel(resolvedRoot, absolute)
	current := resolvedRoot
	for _, component := range strings.Split(filepath.ToSlash(relative), "/") {
		current = filepath.Join(current, component)
		componentInfo, componentErr := os.Lstat(current)
		if componentErr != nil {
			return "", nil, fmt.Errorf("inspect pack manifest path: %w", componentErr)
		}
		if componentInfo.Mode()&os.ModeSymlink != 0 {
			return "", nil, fmt.Errorf("pack manifest path cannot contain symlinks")
		}
	}
	info, err := os.Lstat(absolute)
	if err != nil {
		return "", nil, fmt.Errorf("inspect pack manifest: %w", err)
	}
	if info.Mode()&os.ModeSymlink != 0 || !info.Mode().IsRegular() || info.Size() > maximumPackBytes {
		return "", nil, fmt.Errorf("pack manifest must be a non-symlink regular file no larger than 1 MiB")
	}
	data, err := os.ReadFile(absolute)
	if err != nil {
		return "", nil, fmt.Errorf("read pack manifest: %w", err)
	}
	return resolvedRoot, data, nil
}

func resolveRootFile(root, relative string) (string, error) {
	if err := validateRelativePath(relative); err != nil {
		return "", err
	}
	path := root
	for _, component := range strings.Split(relative, "/") {
		path = filepath.Join(path, component)
		info, err := os.Lstat(path)
		if err != nil {
			return "", err
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return "", fmt.Errorf("pack path cannot contain symlinks")
		}
	}
	info, err := os.Stat(path)
	if err != nil || !info.Mode().IsRegular() {
		return "", fmt.Errorf("pack path must resolve to a regular file")
	}
	if err := requireInside(root, path); err != nil {
		return "", err
	}
	return path, nil
}

func validateRelativePath(path string) error {
	if path == "" || filepath.IsAbs(path) || strings.Contains(path, "\\") || filepath.ToSlash(filepath.Clean(path)) != path ||
		path == "." || path == ".." || strings.HasPrefix(path, "../") {
		return fmt.Errorf("pack path %q must be a normalized repository-relative slash path", path)
	}
	return nil
}

func requireInside(root, path string) error {
	relative, err := filepath.Rel(root, path)
	if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return fmt.Errorf("pack path escapes the repository root")
	}
	return nil
}

func contains(values []string, wanted string) bool {
	for _, value := range values {
		if value == wanted {
			return true
		}
	}
	return false
}

func digest(value string) bool {
	if len(value) != 64 {
		return false
	}
	for _, character := range value {
		if (character < '0' || character > '9') && (character < 'a' || character > 'f') {
			return false
		}
	}
	return true
}

var suiteIDPattern = regexp.MustCompile(`^prc\.benchmark\.[a-z0-9.-]+@[0-9]+\.[0-9]+$`)
