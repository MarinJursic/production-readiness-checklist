// Package benchmark evaluates scanner behavior against small, labeled,
// reproducible fixture repositories. It measures deterministic assertion
// outcomes; it does not turn fixture counts into a readiness score.
package benchmark

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/MarinJursic/production-readiness-checklist/scanner/catalog"
	"github.com/MarinJursic/production-readiness-checklist/scanner/engine"
	"github.com/MarinJursic/production-readiness-checklist/scanner/inventory"
	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
	"gopkg.in/yaml.v3"
)

const (
	SuiteSchema  = "prc.benchmark-suite/v0.1"
	ReportSchema = "prc.benchmark-report/v0.1"
)

const maximumSuiteBytes = 1024 * 1024

var (
	suiteIDPattern = regexp.MustCompile(`^prc\.benchmark\.[a-z0-9.-]+@[0-9]+\.[0-9]+$`)
	caseIDPattern  = regexp.MustCompile(`^[a-z0-9][a-z0-9.-]{0,127}$`)
	assertionID    = regexp.MustCompile(`^PRC-A-[A-Z0-9]+-[0-9]{3}$`)
)

type Expectation struct {
	AssertionID string `json:"assertion_id" yaml:"assertion_id"`
	Assessment  string `json:"assessment" yaml:"assessment"`
	Execution   string `json:"execution" yaml:"execution"`
}

type Case struct {
	ID           string        `json:"id" yaml:"id"`
	Target       string        `json:"target" yaml:"target"`
	Expectations []Expectation `json:"expectations" yaml:"expectations"`
}

type QualityBudget struct {
	MinimumPrecision         float64 `json:"minimum_precision" yaml:"minimum_precision"`
	MinimumRecall            float64 `json:"minimum_recall" yaml:"minimum_recall"`
	MaximumFalsePositiveRate float64 `json:"maximum_false_positive_rate" yaml:"maximum_false_positive_rate"`
	MaximumMismatches        int     `json:"maximum_mismatches" yaml:"maximum_mismatches"`
	RequireDeterminism       bool    `json:"require_determinism" yaml:"require_determinism"`
}

func (budget *QualityBudget) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return fmt.Errorf("quality_budget must be a mapping")
	}
	seen := map[string]bool{}
	for index := 0; index < len(node.Content); index += 2 {
		key, value := node.Content[index].Value, node.Content[index+1]
		if seen[key] {
			return fmt.Errorf("quality_budget contains duplicate field %q", key)
		}
		seen[key] = true
		var destination any
		switch key {
		case "minimum_precision":
			destination = &budget.MinimumPrecision
		case "minimum_recall":
			destination = &budget.MinimumRecall
		case "maximum_false_positive_rate":
			destination = &budget.MaximumFalsePositiveRate
		case "maximum_mismatches":
			destination = &budget.MaximumMismatches
		case "require_determinism":
			destination = &budget.RequireDeterminism
		default:
			return fmt.Errorf("quality_budget contains unknown field %q", key)
		}
		if err := value.Decode(destination); err != nil {
			return fmt.Errorf("decode quality_budget field %s: %w", key, err)
		}
	}
	for _, required := range []string{
		"minimum_precision", "minimum_recall", "maximum_false_positive_rate", "maximum_mismatches", "require_determinism",
	} {
		if !seen[required] {
			return fmt.Errorf("quality_budget requires field %q", required)
		}
	}
	return nil
}

type Suite struct {
	SchemaVersion string        `json:"schema_version" yaml:"schema_version"`
	ID            string        `json:"id" yaml:"id"`
	Title         string        `json:"title" yaml:"title"`
	ProfileID     string        `json:"profile_id" yaml:"profile_id"`
	Cases         []Case        `json:"cases" yaml:"cases"`
	QualityBudget QualityBudget `json:"quality_budget" yaml:"quality_budget"`
}

type LoadedSuite struct {
	Suite  Suite
	Digest string
	base   string
}

type ExpectationResult struct {
	AssertionID        string `json:"assertion_id"`
	ExpectedAssessment string `json:"expected_assessment"`
	ExpectedExecution  string `json:"expected_execution"`
	ActualAssessment   string `json:"actual_assessment"`
	ActualExecution    string `json:"actual_execution"`
	Summary            string `json:"summary"`
	Matched            bool   `json:"matched"`
}

type CaseResult struct {
	ID              string              `json:"id"`
	Target          string              `json:"target"`
	InventoryDigest string              `json:"inventory_digest"`
	RunID           string              `json:"run_id"`
	Deterministic   bool                `json:"deterministic"`
	Passed          bool                `json:"passed"`
	Expectations    []ExpectationResult `json:"expectations"`
}

type OutcomeCounts struct {
	Pass          int `json:"pass"`
	Fail          int `json:"fail"`
	NotApplicable int `json:"not_applicable"`
	Unknown       int `json:"unknown"`
	ManualReview  int `json:"manual_review"`
	Stale         int `json:"stale"`
	Conflicting   int `json:"conflicting"`
}

type ClassificationMetrics struct {
	TruePositive      int     `json:"true_positive"`
	FalsePositive     int     `json:"false_positive"`
	FalseNegative     int     `json:"false_negative"`
	TrueNegative      int     `json:"true_negative"`
	Precision         float64 `json:"precision"`
	Recall            float64 `json:"recall"`
	FalsePositiveRate float64 `json:"false_positive_rate"`
}

type Summary struct {
	Cases              int           `json:"cases"`
	Expectations       int           `json:"expectations"`
	Matched            int           `json:"matched"`
	Mismatched         int           `json:"mismatched"`
	DeterministicCases int           `json:"deterministic_cases"`
	ExpectedOutcomes   OutcomeCounts `json:"expected_outcomes"`
}

type Report struct {
	SchemaVersion   string                `json:"schema_version"`
	SuiteID         string                `json:"suite_id"`
	SuiteDigest     string                `json:"suite_digest"`
	CorpusDigest    string                `json:"corpus_digest"`
	CatalogDigest   string                `json:"catalog_digest"`
	ProfileID       string                `json:"profile_id"`
	EvaluatedAt     time.Time             `json:"evaluated_at"`
	QualityBudget   QualityBudget         `json:"quality_budget"`
	Summary         Summary               `json:"summary"`
	Metrics         ClassificationMetrics `json:"metrics"`
	Cases           []CaseResult          `json:"cases"`
	QualityFailures []string              `json:"quality_failures"`
	Passed          bool                  `json:"passed"`
}

func Load(path string, catalogValue *catalog.Catalog) (LoadedSuite, error) {
	data, base, err := readSuite(path)
	if err != nil {
		return LoadedSuite{}, err
	}
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)
	var suite Suite
	if err := decoder.Decode(&suite); err != nil {
		return LoadedSuite{}, fmt.Errorf("decode benchmark suite: %w", err)
	}
	var trailing any
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		if err == nil {
			return LoadedSuite{}, fmt.Errorf("benchmark suite contains more than one YAML document")
		}
		return LoadedSuite{}, fmt.Errorf("decode trailing benchmark suite content: %w", err)
	}
	if err := validateAndNormalize(&suite, catalogValue); err != nil {
		return LoadedSuite{}, err
	}
	for _, item := range suite.Cases {
		if _, err := resolveTarget(base, item.Target); err != nil {
			return LoadedSuite{}, fmt.Errorf("benchmark case %s target: %w", item.ID, err)
		}
	}
	payload, err := json.Marshal(suite)
	if err != nil {
		return LoadedSuite{}, fmt.Errorf("encode benchmark suite identity: %w", err)
	}
	digest := sha256.Sum256(payload)
	return LoadedSuite{Suite: suite, Digest: hex.EncodeToString(digest[:]), base: base}, nil
}

func Evaluate(catalogValue *catalog.Catalog, suitePath string, evaluatedAt time.Time) (Report, error) {
	if evaluatedAt.IsZero() {
		return Report{}, fmt.Errorf("benchmark evaluation time is required")
	}
	loaded, err := Load(suitePath, catalogValue)
	if err != nil {
		return Report{}, err
	}
	catalogDigest, err := catalogValue.Digest()
	if err != nil {
		return Report{}, err
	}
	report := Report{
		SchemaVersion: ReportSchema, SuiteID: loaded.Suite.ID, SuiteDigest: loaded.Digest,
		CatalogDigest: catalogDigest, ProfileID: loaded.Suite.ProfileID, EvaluatedAt: evaluatedAt.UTC(),
		QualityBudget: loaded.Suite.QualityBudget, Cases: []CaseResult{}, QualityFailures: []string{},
	}
	for _, benchmarkCase := range loaded.Suite.Cases {
		target, err := resolveTarget(loaded.base, benchmarkCase.Target)
		if err != nil {
			return Report{}, fmt.Errorf("benchmark case %s target: %w", benchmarkCase.ID, err)
		}
		item, err := inventory.Build(target)
		if err != nil {
			return Report{}, fmt.Errorf("benchmark case %s inventory: %w", benchmarkCase.ID, err)
		}
		runner := engine.New(catalogValue)
		runner.Now = func() time.Time { return evaluatedAt.UTC() }
		first, err := runner.Scan(loaded.Suite.ProfileID, item)
		if err != nil {
			return Report{}, fmt.Errorf("benchmark case %s first scan: %w", benchmarkCase.ID, err)
		}
		second, err := runner.Scan(loaded.Suite.ProfileID, item)
		if err != nil {
			return Report{}, fmt.Errorf("benchmark case %s repeated scan: %w", benchmarkCase.ID, err)
		}
		caseResult := evaluateCase(benchmarkCase, first, second)
		report.Cases = append(report.Cases, caseResult)
		accumulate(&report, benchmarkCase, caseResult)
	}
	report.CorpusDigest, err = corpusDigest(report.SuiteDigest, report.Cases)
	if err != nil {
		return Report{}, err
	}
	finalize(&report)
	return report, nil
}

func corpusDigest(suiteDigest string, cases []CaseResult) (string, error) {
	type identity struct {
		ID              string `json:"id"`
		InventoryDigest string `json:"inventory_digest"`
	}
	entries := make([]identity, 0, len(cases))
	for _, item := range cases {
		entries = append(entries, identity{ID: item.ID, InventoryDigest: item.InventoryDigest})
	}
	payload, err := json.Marshal(struct {
		SuiteDigest string     `json:"suite_digest"`
		Cases       []identity `json:"cases"`
	}{SuiteDigest: suiteDigest, Cases: entries})
	if err != nil {
		return "", fmt.Errorf("encode benchmark corpus identity: %w", err)
	}
	digest := sha256.Sum256(payload)
	return hex.EncodeToString(digest[:]), nil
}

func evaluateCase(benchmarkCase Case, first, second model.RunResult) CaseResult {
	actual := make(map[string]model.AssertionResult, len(first.Results))
	for _, result := range first.Results {
		actual[result.AssertionID] = result
	}
	caseResult := CaseResult{
		ID: benchmarkCase.ID, Target: benchmarkCase.Target, InventoryDigest: first.Inventory.Digest,
		RunID: first.RunID, Deterministic: first.RunID == second.RunID,
		Passed: true, Expectations: []ExpectationResult{},
	}
	for _, expected := range benchmarkCase.Expectations {
		observed := actual[expected.AssertionID]
		matched := observed.Assessment == expected.Assessment && observed.Execution == expected.Execution
		caseResult.Expectations = append(caseResult.Expectations, ExpectationResult{
			AssertionID: expected.AssertionID, ExpectedAssessment: expected.Assessment,
			ExpectedExecution: expected.Execution, ActualAssessment: observed.Assessment,
			ActualExecution: observed.Execution, Summary: observed.Summary, Matched: matched,
		})
		caseResult.Passed = caseResult.Passed && matched
	}
	caseResult.Passed = caseResult.Passed && caseResult.Deterministic
	return caseResult
}

func accumulate(report *Report, benchmarkCase Case, result CaseResult) {
	report.Summary.Cases++
	if result.Deterministic {
		report.Summary.DeterministicCases++
	}
	for index, expected := range benchmarkCase.Expectations {
		actual := result.Expectations[index]
		report.Summary.Expectations++
		if actual.Matched {
			report.Summary.Matched++
		} else {
			report.Summary.Mismatched++
		}
		switch expected.Assessment {
		case "pass":
			report.Summary.ExpectedOutcomes.Pass++
		case "fail":
			report.Summary.ExpectedOutcomes.Fail++
		case "not_applicable":
			report.Summary.ExpectedOutcomes.NotApplicable++
		case "unknown":
			report.Summary.ExpectedOutcomes.Unknown++
		case "manual_review":
			report.Summary.ExpectedOutcomes.ManualReview++
		case "stale":
			report.Summary.ExpectedOutcomes.Stale++
		case "conflicting":
			report.Summary.ExpectedOutcomes.Conflicting++
		}
		expectedFailure := expected.Assessment == "fail"
		actualFailure := actual.ActualAssessment == "fail"
		switch {
		case expectedFailure && actualFailure:
			report.Metrics.TruePositive++
		case !expectedFailure && actualFailure:
			report.Metrics.FalsePositive++
		case expectedFailure && !actualFailure:
			report.Metrics.FalseNegative++
		default:
			report.Metrics.TrueNegative++
		}
	}
}

func finalize(report *Report) {
	report.Metrics.Precision = positiveRate(
		report.Metrics.TruePositive, report.Metrics.TruePositive+report.Metrics.FalsePositive, 1,
	)
	report.Metrics.Recall = positiveRate(
		report.Metrics.TruePositive, report.Metrics.TruePositive+report.Metrics.FalseNegative, 1,
	)
	report.Metrics.FalsePositiveRate = positiveRate(
		report.Metrics.FalsePositive, report.Metrics.FalsePositive+report.Metrics.TrueNegative, 0,
	)
	budget := report.QualityBudget
	if report.Summary.Mismatched > budget.MaximumMismatches {
		report.QualityFailures = append(report.QualityFailures, fmt.Sprintf(
			"mismatched expectations %d exceed budget %d", report.Summary.Mismatched, budget.MaximumMismatches,
		))
	}
	if budget.RequireDeterminism && report.Summary.DeterministicCases != report.Summary.Cases {
		report.QualityFailures = append(report.QualityFailures, "one or more benchmark cases were nondeterministic")
	}
	if report.Metrics.Precision < budget.MinimumPrecision {
		report.QualityFailures = append(report.QualityFailures, fmt.Sprintf(
			"precision %.6f is below budget %.6f", report.Metrics.Precision, budget.MinimumPrecision,
		))
	}
	if report.Metrics.Recall < budget.MinimumRecall {
		report.QualityFailures = append(report.QualityFailures, fmt.Sprintf(
			"recall %.6f is below budget %.6f", report.Metrics.Recall, budget.MinimumRecall,
		))
	}
	if report.Metrics.FalsePositiveRate > budget.MaximumFalsePositiveRate {
		report.QualityFailures = append(report.QualityFailures, fmt.Sprintf(
			"false-positive rate %.6f exceeds budget %.6f", report.Metrics.FalsePositiveRate, budget.MaximumFalsePositiveRate,
		))
	}
	report.Passed = len(report.QualityFailures) == 0
}

func positiveRate(numerator, denominator int, empty float64) float64 {
	if denominator == 0 {
		return empty
	}
	return float64(numerator) / float64(denominator)
}

func readSuite(path string) ([]byte, string, error) {
	absolute, err := filepath.Abs(path)
	if err != nil {
		return nil, "", fmt.Errorf("resolve benchmark suite path: %w", err)
	}
	info, err := os.Lstat(absolute)
	if err != nil {
		return nil, "", fmt.Errorf("inspect benchmark suite: %w", err)
	}
	if info.Mode()&os.ModeSymlink != 0 || !info.Mode().IsRegular() || info.Size() > maximumSuiteBytes {
		return nil, "", fmt.Errorf("benchmark suite must be a non-symlink regular file no larger than 1 MiB")
	}
	data, err := os.ReadFile(absolute)
	if err != nil {
		return nil, "", fmt.Errorf("read benchmark suite: %w", err)
	}
	base, err := filepath.EvalSymlinks(filepath.Dir(absolute))
	if err != nil {
		return nil, "", fmt.Errorf("resolve benchmark suite directory: %w", err)
	}
	return data, base, nil
}

func validateAndNormalize(suite *Suite, catalogValue *catalog.Catalog) error {
	if suite.SchemaVersion != SuiteSchema || !suiteIDPattern.MatchString(suite.ID) || strings.TrimSpace(suite.Title) == "" {
		return fmt.Errorf("benchmark suite requires a supported schema, valid ID, and title")
	}
	profile, err := catalogValue.Profile(suite.ProfileID)
	if err != nil {
		return fmt.Errorf("benchmark suite profile: %w", err)
	}
	if len(suite.Cases) == 0 || len(suite.Cases) > 100 {
		return fmt.Errorf("benchmark suite must contain between 1 and 100 cases")
	}
	budget := suite.QualityBudget
	if budget.MinimumPrecision < 0 || budget.MinimumPrecision > 1 || budget.MinimumRecall < 0 || budget.MinimumRecall > 1 ||
		budget.MaximumFalsePositiveRate < 0 || budget.MaximumFalsePositiveRate > 1 ||
		budget.MaximumMismatches < 0 || budget.MaximumMismatches > 10000 {
		return fmt.Errorf("benchmark suite quality budget is invalid")
	}
	profileAssertions := make(map[string]bool, len(profile.AssertionIDs))
	for _, identifier := range profile.AssertionIDs {
		profileAssertions[identifier] = true
	}
	seenCases := map[string]bool{}
	for caseIndex := range suite.Cases {
		item := &suite.Cases[caseIndex]
		if !caseIDPattern.MatchString(item.ID) || seenCases[item.ID] {
			return fmt.Errorf("benchmark suite contains an invalid or duplicate case ID %q", item.ID)
		}
		seenCases[item.ID] = true
		if err := validateRelativePath(item.Target); err != nil {
			return fmt.Errorf("benchmark case %s: %w", item.ID, err)
		}
		if len(item.Expectations) == 0 || len(item.Expectations) > 1000 {
			return fmt.Errorf("benchmark case %s must contain between 1 and 1000 expectations", item.ID)
		}
		seenAssertions := map[string]bool{}
		for expectationIndex := range item.Expectations {
			expected := item.Expectations[expectationIndex]
			if !assertionID.MatchString(expected.AssertionID) || !profileAssertions[expected.AssertionID] ||
				seenAssertions[expected.AssertionID] {
				return fmt.Errorf("benchmark case %s has an invalid, unprofiled, or duplicate assertion %q", item.ID, expected.AssertionID)
			}
			seenAssertions[expected.AssertionID] = true
			if !validAssessment(expected.Assessment) || !validExecution(expected.Execution) ||
				!validExpectedOutcome(expected.Assessment, expected.Execution) {
				return fmt.Errorf("benchmark case %s assertion %s has an invalid expected outcome", item.ID, expected.AssertionID)
			}
		}
		sort.Slice(item.Expectations, func(left, right int) bool {
			return item.Expectations[left].AssertionID < item.Expectations[right].AssertionID
		})
	}
	sort.Slice(suite.Cases, func(left, right int) bool { return suite.Cases[left].ID < suite.Cases[right].ID })
	return nil
}

func validAssessment(value string) bool {
	return value == "pass" || value == "fail" || value == "unknown" || value == "manual_review" ||
		value == "not_applicable" || value == "stale" || value == "conflicting"
}

func validExecution(value string) bool {
	return value == "not_run" || value == "completed" || value == "blocked" || value == "error"
}

func validExpectedOutcome(assessment, execution string) bool {
	switch assessment {
	case "pass", "fail", "manual_review":
		return execution == "completed"
	case "not_applicable":
		return execution == "not_run"
	case "unknown", "stale", "conflicting":
		return true
	default:
		return false
	}
}

func validateRelativePath(path string) error {
	if path == "" || filepath.IsAbs(path) || strings.Contains(path, "\\") || filepath.ToSlash(filepath.Clean(path)) != path ||
		path == "." || path == ".." || strings.HasPrefix(path, "../") {
		return fmt.Errorf("target %q must be a normalized relative slash path", path)
	}
	return nil
}

func resolveTarget(base, relative string) (string, error) {
	if err := validateRelativePath(relative); err != nil {
		return "", err
	}
	path := base
	for _, component := range strings.Split(relative, "/") {
		path = filepath.Join(path, component)
		info, err := os.Lstat(path)
		if err != nil {
			return "", err
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return "", fmt.Errorf("benchmark target path cannot contain symlinks")
		}
	}
	info, err := os.Stat(path)
	if err != nil || !info.IsDir() {
		return "", fmt.Errorf("benchmark target must be a directory")
	}
	entries := 0
	if err := filepath.WalkDir(path, func(_ string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		entries++
		if entries > 10000 {
			return fmt.Errorf("benchmark target exceeds 10000 entries")
		}
		if entry.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("benchmark target cannot contain symlinks")
		}
		return nil
	}); err != nil {
		return "", err
	}
	resolved, err := filepath.EvalSymlinks(path)
	if err != nil {
		return "", err
	}
	relativeToBase, err := filepath.Rel(base, resolved)
	if err != nil || relativeToBase == ".." || strings.HasPrefix(relativeToBase, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("benchmark target escapes the suite directory")
	}
	return resolved, nil
}
