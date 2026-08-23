// Package adapterfixture validates recorded adapter protocol transcripts against
// an exact adapter manifest. It is a release-quality gate, not an adapter
// execution mechanism.
package adapterfixture

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
	"strings"

	"github.com/MarinJursic/production-readiness-checklist/scanner/adapter"
	"gopkg.in/yaml.v3"
)

const (
	SuiteSchema  = "prc.adapter-fixture-suite/v0.1"
	ReportSchema = "prc.adapter-fixture-report/v0.1"

	maximumSuiteBytes = 1024 * 1024
	maximumCases      = 100
)

var (
	suiteIDPattern = regexp.MustCompile(`^prc\.adapter-fixture\.[a-z0-9.-]+@[0-9]+\.[0-9]+$`)
	caseIDPattern  = regexp.MustCompile(`^[a-z0-9][a-z0-9.-]{0,127}$`)
	hexPattern     = regexp.MustCompile(`^[0-9a-f]{64}$`)
)

type ManifestReference struct {
	Path   string `json:"path" yaml:"path"`
	SHA256 string `json:"sha256" yaml:"sha256"`
}

type LimitOverride struct {
	MaxLineBytes *int `json:"max_line_bytes,omitempty" yaml:"max_line_bytes,omitempty"`
	MaxMessages  *int `json:"max_messages,omitempty" yaml:"max_messages,omitempty"`
	MaxStdout    *int `json:"max_stdout_bytes,omitempty" yaml:"max_stdout_bytes,omitempty"`
}

type ObservationExpectation struct {
	ID      string `json:"id" yaml:"id"`
	Kind    string `json:"kind" yaml:"kind"`
	Outcome string `json:"outcome" yaml:"outcome"`
}

type Expectation struct {
	Disposition   string                   `json:"disposition" yaml:"disposition"`
	SummaryStatus string                   `json:"summary_status,omitempty" yaml:"summary_status,omitempty"`
	ErrorCode     string                   `json:"error_code,omitempty" yaml:"error_code,omitempty"`
	Observations  []ObservationExpectation `json:"observations,omitempty" yaml:"observations,omitempty"`
}

type Case struct {
	ID       string         `json:"id" yaml:"id"`
	Output   string         `json:"output" yaml:"output"`
	Limits   *LimitOverride `json:"limits,omitempty" yaml:"limits,omitempty"`
	Expected Expectation    `json:"expected" yaml:"expected"`
}

type Suite struct {
	SchemaVersion string            `json:"schema_version" yaml:"schema_version"`
	ID            string            `json:"id" yaml:"id"`
	Title         string            `json:"title" yaml:"title"`
	Manifest      ManifestReference `json:"manifest" yaml:"manifest"`
	Cases         []Case            `json:"cases" yaml:"cases"`
}

type LoadedSuite struct {
	Suite    Suite
	Manifest adapter.Manifest
	Digest   string
	base     string
}

type Actual struct {
	Disposition   string                   `json:"disposition"`
	SummaryStatus string                   `json:"summary_status,omitempty"`
	ErrorCode     string                   `json:"error_code,omitempty"`
	Error         string                   `json:"error,omitempty"`
	Observations  []ObservationExpectation `json:"observations"`
}

type CaseResult struct {
	ID            string      `json:"id"`
	Output        string      `json:"output"`
	OutputSHA256  string      `json:"output_sha256"`
	Expected      Expectation `json:"expected"`
	Actual        Actual      `json:"actual"`
	Deterministic bool        `json:"deterministic"`
	Passed        bool        `json:"passed"`
}

type Summary struct {
	Cases              int `json:"cases"`
	Matched            int `json:"matched"`
	Mismatched         int `json:"mismatched"`
	DeterministicCases int `json:"deterministic_cases"`
}

type Report struct {
	SchemaVersion   string       `json:"schema_version"`
	SuiteID         string       `json:"suite_id"`
	SuiteDigest     string       `json:"suite_digest"`
	CorpusDigest    string       `json:"corpus_digest"`
	AdapterID       string       `json:"adapter_id"`
	ManifestSHA256  string       `json:"manifest_sha256"`
	Summary         Summary      `json:"summary"`
	Cases           []CaseResult `json:"cases"`
	QualityFailures []string     `json:"quality_failures"`
	Passed          bool         `json:"passed"`
}

func Load(path string) (LoadedSuite, error) {
	data, base, err := readSuite(path)
	if err != nil {
		return LoadedSuite{}, err
	}
	if err := validateYAML(data); err != nil {
		return LoadedSuite{}, fmt.Errorf("decode adapter fixture suite: %w", err)
	}
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)
	var suite Suite
	if err := decoder.Decode(&suite); err != nil {
		return LoadedSuite{}, fmt.Errorf("decode adapter fixture suite: %w", err)
	}
	var trailing any
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		if err == nil {
			return LoadedSuite{}, fmt.Errorf("adapter fixture suite contains more than one YAML document")
		}
		return LoadedSuite{}, fmt.Errorf("decode trailing adapter fixture suite content: %w", err)
	}
	if err := validateSuite(suite); err != nil {
		return LoadedSuite{}, err
	}
	manifestPath, err := resolveRegular(base, suite.Manifest.Path, maximumSuiteBytes)
	if err != nil {
		return LoadedSuite{}, fmt.Errorf("adapter fixture manifest: %w", err)
	}
	manifest, err := adapter.LoadManifest(manifestPath)
	if err != nil {
		return LoadedSuite{}, fmt.Errorf("load adapter fixture manifest: %w", err)
	}
	manifestDigest, err := adapter.ManifestDigest(manifest)
	if err != nil {
		return LoadedSuite{}, err
	}
	if manifestDigest != suite.Manifest.SHA256 {
		return LoadedSuite{}, fmt.Errorf("adapter fixture manifest digest mismatch: expected %s, got %s", suite.Manifest.SHA256, manifestDigest)
	}
	for _, item := range suite.Cases {
		_, err := caseLimits(manifest.Resources.Limits, item.Limits)
		if err != nil {
			return LoadedSuite{}, fmt.Errorf("adapter fixture case %s limits: %w", item.ID, err)
		}
		if _, err := resolveRegular(base, item.Output, int64(manifest.Resources.MaxStdout)+1); err != nil {
			return LoadedSuite{}, fmt.Errorf("adapter fixture case %s output: %w", item.ID, err)
		}
	}
	payload, err := json.Marshal(suite)
	if err != nil {
		return LoadedSuite{}, fmt.Errorf("encode adapter fixture suite identity: %w", err)
	}
	digest := sha256.Sum256(payload)
	return LoadedSuite{Suite: suite, Manifest: manifest, Digest: hex.EncodeToString(digest[:]), base: base}, nil
}

func Evaluate(path string) (Report, error) {
	loaded, err := Load(path)
	if err != nil {
		return Report{}, err
	}
	report := Report{
		SchemaVersion: ReportSchema, SuiteID: loaded.Suite.ID, SuiteDigest: loaded.Digest,
		AdapterID: loaded.Manifest.ID, ManifestSHA256: loaded.Suite.Manifest.SHA256,
		Cases: []CaseResult{}, QualityFailures: []string{},
	}
	corpus := sha256.New()
	_, _ = io.WriteString(corpus, loaded.Digest)
	_, _ = io.WriteString(corpus, loaded.Suite.Manifest.SHA256)
	for _, item := range loaded.Suite.Cases {
		limits, _ := caseLimits(loaded.Manifest.Resources.Limits, item.Limits)
		outputPath, _ := resolveRegular(loaded.base, item.Output, int64(loaded.Manifest.Resources.MaxStdout)+1)
		data, err := os.ReadFile(outputPath)
		if err != nil {
			return Report{}, fmt.Errorf("read adapter fixture case %s: %w", item.ID, err)
		}
		outputDigest := sha256.Sum256(data)
		outputSHA := hex.EncodeToString(outputDigest[:])
		_, _ = io.WriteString(corpus, item.ID)
		_, _ = io.WriteString(corpus, outputSHA)
		first := evaluateOnce(loaded.Manifest, data, limits)
		second := evaluateOnce(loaded.Manifest, data, limits)
		firstJSON, _ := json.Marshal(first)
		secondJSON, _ := json.Marshal(second)
		deterministic := bytes.Equal(firstJSON, secondJSON)
		matched := matches(item.Expected, first)
		result := CaseResult{
			ID: item.ID, Output: item.Output, OutputSHA256: outputSHA, Expected: item.Expected,
			Actual: first, Deterministic: deterministic, Passed: matched && deterministic,
		}
		report.Cases = append(report.Cases, result)
		report.Summary.Cases++
		if deterministic {
			report.Summary.DeterministicCases++
		}
		if result.Passed {
			report.Summary.Matched++
		} else {
			report.Summary.Mismatched++
			report.QualityFailures = append(report.QualityFailures, fmt.Sprintf("case %s did not match its deterministic expectation", item.ID))
		}
	}
	report.CorpusDigest = hex.EncodeToString(corpus.Sum(nil))
	report.Passed = report.Summary.Cases > 0 && report.Summary.Mismatched == 0 && report.Summary.DeterministicCases == report.Summary.Cases
	return report, nil
}

func evaluateOnce(manifest adapter.Manifest, data []byte, limits adapter.Limits) Actual {
	transcript, err := adapter.ParseOutput(bytes.NewReader(data), limits)
	if err != nil {
		return Actual{Disposition: "protocol_rejected", ErrorCode: protocolErrorCode(err), Error: err.Error(), Observations: []ObservationExpectation{}}
	}
	if err := adapter.ValidateTranscriptContract(manifest, transcript); err != nil {
		code := "contract_violation"
		if strings.Contains(err.Error(), "undeclared observation kind") {
			code = "undeclared_observation_kind"
		}
		return Actual{Disposition: "contract_rejected", ErrorCode: code, Error: err.Error(), Observations: observations(transcript)}
	}
	return Actual{Disposition: "accepted", SummaryStatus: transcript.Summary.Status, Observations: observations(transcript)}
}

func observations(transcript adapter.Transcript) []ObservationExpectation {
	result := make([]ObservationExpectation, 0, len(transcript.Observations))
	for _, item := range transcript.Observations {
		result = append(result, ObservationExpectation{ID: item.ID, Kind: item.Kind, Outcome: item.Outcome})
	}
	return result
}

func matches(expected Expectation, actual Actual) bool {
	if expected.Disposition != actual.Disposition || expected.SummaryStatus != actual.SummaryStatus || expected.ErrorCode != actual.ErrorCode {
		return false
	}
	if len(expected.Observations) != len(actual.Observations) {
		return false
	}
	for index := range expected.Observations {
		if expected.Observations[index] != actual.Observations[index] {
			return false
		}
	}
	return true
}

func protocolErrorCode(err error) string {
	message := err.Error()
	switch {
	case strings.Contains(message, "exceeds"):
		return "resource_limit"
	case strings.Contains(message, "not valid JSON"), strings.Contains(message, "unexpected end of JSON"),
		strings.Contains(message, "is blank"), strings.Contains(message, "unsupported message type"):
		return "malformed_output"
	case strings.Contains(message, "adapter protocol error"):
		return "schema_violation"
	case strings.Contains(message, "after its summary"), strings.Contains(message, "without a summary"):
		return "sequence_violation"
	default:
		return "protocol_violation"
	}
}

func validateSuite(suite Suite) error {
	if suite.SchemaVersion != SuiteSchema {
		return fmt.Errorf("unsupported adapter fixture suite schema %q", suite.SchemaVersion)
	}
	if !suiteIDPattern.MatchString(suite.ID) || strings.TrimSpace(suite.Title) == "" {
		return fmt.Errorf("adapter fixture suite requires a valid ID and title")
	}
	if err := validateRelativePath(suite.Manifest.Path); err != nil || !strings.HasSuffix(suite.Manifest.Path, ".yaml") {
		return fmt.Errorf("adapter fixture manifest path must be a normalized relative YAML path")
	}
	if !hexPattern.MatchString(suite.Manifest.SHA256) {
		return fmt.Errorf("adapter fixture manifest requires a lowercase SHA-256 digest")
	}
	if len(suite.Cases) == 0 || len(suite.Cases) > maximumCases {
		return fmt.Errorf("adapter fixture suite must contain between 1 and %d cases", maximumCases)
	}
	seen := map[string]bool{}
	for _, item := range suite.Cases {
		if !caseIDPattern.MatchString(item.ID) || seen[item.ID] {
			return fmt.Errorf("adapter fixture case IDs must be valid and unique")
		}
		seen[item.ID] = true
		if err := validateRelativePath(item.Output); err != nil || !strings.HasSuffix(item.Output, ".jsonl") {
			return fmt.Errorf("adapter fixture case %s output must be a normalized relative JSONL path", item.ID)
		}
		if err := validateExpectation(item.Expected); err != nil {
			return fmt.Errorf("adapter fixture case %s: %w", item.ID, err)
		}
	}
	return nil
}

func validateExpectation(expected Expectation) error {
	switch expected.Disposition {
	case "accepted":
		if expected.SummaryStatus == "" || expected.ErrorCode != "" {
			return fmt.Errorf("accepted expectation requires summary_status and forbids error_code")
		}
		if !validSummaryStatus(expected.SummaryStatus) {
			return fmt.Errorf("accepted expectation has unsupported summary_status %q", expected.SummaryStatus)
		}
	case "protocol_rejected":
		if expected.SummaryStatus != "" || expected.ErrorCode == "" || len(expected.Observations) != 0 {
			return fmt.Errorf("rejected expectation requires error_code and forbids summary_status and observations")
		}
		if !validProtocolErrorCode(expected.ErrorCode) {
			return fmt.Errorf("protocol-rejected expectation has unsupported error_code %q", expected.ErrorCode)
		}
	case "contract_rejected":
		if expected.SummaryStatus != "" || expected.ErrorCode == "" || len(expected.Observations) == 0 {
			return fmt.Errorf("contract-rejected expectation requires error_code and at least one exact observation")
		}
		if expected.ErrorCode != "undeclared_observation_kind" && expected.ErrorCode != "contract_violation" {
			return fmt.Errorf("contract-rejected expectation has unsupported error_code %q", expected.ErrorCode)
		}
	default:
		return fmt.Errorf("unsupported expected disposition %q", expected.Disposition)
	}
	seen := map[string]bool{}
	for _, item := range expected.Observations {
		if strings.TrimSpace(item.ID) == "" || strings.TrimSpace(item.Kind) == "" || seen[item.ID] {
			return fmt.Errorf("expected observations require unique nonempty IDs and kinds")
		}
		seen[item.ID] = true
		switch item.Outcome {
		case "found", "not_found", "value", "unsupported":
		default:
			return fmt.Errorf("expected observation %s has invalid outcome %q", item.ID, item.Outcome)
		}
	}
	return nil
}

func validSummaryStatus(value string) bool {
	switch value {
	case "completed", "partial", "unsupported", "configuration_error", "execution_error", "timeout", "parse_error":
		return true
	default:
		return false
	}
}

func validProtocolErrorCode(value string) bool {
	switch value {
	case "resource_limit", "schema_violation", "malformed_output", "sequence_violation", "protocol_violation":
		return true
	default:
		return false
	}
}

func caseLimits(base adapter.Limits, override *LimitOverride) (adapter.Limits, error) {
	result := base
	if override == nil {
		return result, nil
	}
	if override.MaxLineBytes == nil && override.MaxMessages == nil && override.MaxStdout == nil {
		return adapter.Limits{}, fmt.Errorf("limit override must contain at least one field")
	}
	values := []struct {
		name        string
		value       *int
		ceiling     int
		destination *int
	}{
		{"max_line_bytes", override.MaxLineBytes, base.MaxLineBytes, &result.MaxLineBytes},
		{"max_messages", override.MaxMessages, base.MaxMessages, &result.MaxMessages},
		{"max_stdout_bytes", override.MaxStdout, base.MaxStdout, &result.MaxStdout},
	}
	for _, item := range values {
		if item.value == nil {
			continue
		}
		if *item.value < 1 || *item.value > item.ceiling {
			return adapter.Limits{}, fmt.Errorf("%s must be positive and cannot exceed manifest value %d", item.name, item.ceiling)
		}
		*item.destination = *item.value
	}
	if result.MaxLineBytes > result.MaxStdout {
		return adapter.Limits{}, fmt.Errorf("max_line_bytes cannot exceed max_stdout_bytes")
	}
	return result, nil
}

func readSuite(path string) ([]byte, string, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return nil, "", fmt.Errorf("inspect adapter fixture suite: %w", err)
	}
	if !info.Mode().IsRegular() || info.Size() > maximumSuiteBytes {
		return nil, "", fmt.Errorf("adapter fixture suite must be a non-symlink regular file no larger than 1 MiB")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, "", fmt.Errorf("read adapter fixture suite: %w", err)
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, "", fmt.Errorf("resolve adapter fixture suite: %w", err)
	}
	base, err := filepath.EvalSymlinks(filepath.Dir(abs))
	if err != nil {
		return nil, "", fmt.Errorf("resolve adapter fixture directory: %w", err)
	}
	return data, base, nil
}

func resolveRegular(base, relative string, maximum int64) (string, error) {
	if err := validateRelativePath(relative); err != nil {
		return "", err
	}
	path := filepath.Join(base, filepath.FromSlash(relative))
	current := base
	for _, part := range strings.Split(filepath.FromSlash(relative), string(filepath.Separator)) {
		current = filepath.Join(current, part)
		info, err := os.Lstat(current)
		if err != nil {
			return "", err
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return "", fmt.Errorf("path %q cannot contain symlinks", relative)
		}
	}
	info, err := os.Stat(path)
	if err != nil {
		return "", err
	}
	if !info.Mode().IsRegular() || info.Size() > maximum {
		return "", fmt.Errorf("path %q must be a regular file no larger than %d bytes", relative, maximum)
	}
	return path, nil
}

func validateRelativePath(path string) error {
	if path == "" || strings.Contains(path, "\\") || filepath.IsAbs(path) {
		return fmt.Errorf("path %q must be nonempty and relative", path)
	}
	clean := filepath.Clean(filepath.FromSlash(path))
	if clean == "." || clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) || filepath.ToSlash(clean) != path {
		return fmt.Errorf("path %q must be normalized and remain inside the suite directory", path)
	}
	return nil
}

func validateYAML(data []byte) error {
	var document yaml.Node
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	if err := decoder.Decode(&document); err != nil {
		return err
	}
	var walk func(*yaml.Node) error
	walk = func(node *yaml.Node) error {
		if node.Alias != nil || node.Kind == yaml.AliasNode {
			return fmt.Errorf("YAML aliases are not allowed")
		}
		if node.Kind == yaml.MappingNode {
			seen := map[string]bool{}
			for index := 0; index < len(node.Content); index += 2 {
				key := node.Content[index]
				if key.Kind != yaml.ScalarNode || seen[key.Value] {
					return fmt.Errorf("YAML mapping keys must be unique scalar values")
				}
				seen[key.Value] = true
			}
		}
		for _, child := range node.Content {
			if err := walk(child); err != nil {
				return err
			}
		}
		return nil
	}
	return walk(&document)
}
