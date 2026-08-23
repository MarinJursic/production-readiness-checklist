package adapter

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
)

const (
	CheckovProtocolVersion     = "prc-adapter-checkov-json-v1"
	CheckovOutputSchemaVersion = "checkov.report/v3.3.8+prc-v1"
	CheckovToolVersion         = "3.3.8"
	CheckovObservationKind     = "iac-policy"
	CheckovImage               = "docker.io/bridgecrew/checkov@sha256:c64ffb6d6fc8087c896341a2c697770a04a1cf558db04fa7b8129d8ca6bce336"
	CheckovConfigSnapshotPath  = ".prc/checkov-config.json"
	checkovArtifactID          = "checkov-iac-policy-report"
	checkovArtifactMediaType   = "application/vnd.prc.checkov.iac-policy-report+json;version=1"
	checkovReportURL           = "Add an api key '--bc-api-key <api-key>' to see more detailed insights via https://bridgecrew.cloud"
	maxCheckovInputFiles       = 10_000
)

var (
	checkovCheckIDPattern = regexp.MustCompile(`^(?:CKV2?|BC)_[A-Z0-9][A-Z0-9_]{0,31}_[A-Za-z0-9][A-Za-z0-9._-]{0,127}$`)
	checkovFrameworks     = []string{"terraform", "terraform_json", "kubernetes", "dockerfile"}
)

func checkovCommand() []string {
	// Checkov's pre-parser validates an explicitly trusted config only when the
	// flag and path are separate arguments. The equals form falls through to its
	// default .checkov.yml discovery and would let target-owned settings hide
	// results.
	return []string{"--config-file", "/workspace/" + CheckovConfigSnapshotPath}
}

func validateCheckovManifest(manifest Manifest) error {
	if manifest.Protocol != CheckovProtocolVersion || manifest.OutputSchema != CheckovOutputSchemaVersion {
		return fmt.Errorf("checkov adapter requires its exact protocol and output schema")
	}
	if manifest.Image != CheckovImage || !slices.Equal(manifest.Command, checkovCommand()) {
		return fmt.Errorf("checkov adapter requires the reviewed immutable image and scanner-owned command")
	}
	if manifest.Tool.Name != "checkov" || manifest.Tool.Version != CheckovToolVersion ||
		manifest.Tool.Upstream != "https://github.com/bridgecrewio/checkov" ||
		len(manifest.Tool.Formats) != 1 || manifest.Tool.Formats[0].Name != "checkov-json" ||
		!slices.Equal(manifest.Tool.Formats[0].Versions, []string{CheckovToolVersion}) {
		return fmt.Errorf("checkov adapter tool identity does not match the reviewed normalizer")
	}
	if !slices.Equal(manifest.ObservationKinds, []string{CheckovObservationKind}) {
		return fmt.Errorf("checkov adapter must declare only %s observations", CheckovObservationKind)
	}
	if !manifest.Capabilities.WriteScratch || !manifest.Capabilities.ChildProcesses {
		return fmt.Errorf("checkov adapter requires a bounded scratch and OS-task allowance")
	}
	if len(manifest.DataMounts) != 0 || manifest.Resources.TimeoutSeconds < 60 ||
		manifest.Resources.MemoryMB < 1024 || manifest.Resources.PIDs < 32 || manifest.Resources.PIDs > 256 ||
		manifest.Resources.TmpfsMB < 128 || manifest.Resources.MaxStdout < 16*1024*1024 {
		return fmt.Errorf("checkov adapter resource limits cannot support or cannot safely bound the reviewed command")
	}
	return nil
}

type checkovConfigDocument struct {
	Files                   []string `json:"file"`
	Frameworks              []string `json:"framework"`
	Output                  []string `json:"output"`
	Quiet                   bool     `json:"quiet"`
	Compact                 bool     `json:"compact"`
	SoftFail                bool     `json:"soft-fail"`
	SkipDownload            bool     `json:"skip-download"`
	SkipResultsUpload       bool     `json:"skip-results-upload"`
	DownloadExternalModules bool     `json:"download-external-modules"`
	EvaluateVariables       bool     `json:"evaluate-variables"`
}

func checkovConfig(item model.Inventory) ([]byte, []string, error) {
	selected := append([]string{}, item.ContainerFiles...)
	selected = append(selected, item.Infrastructure.TerraformFiles...)
	selected = append(selected, item.Infrastructure.KubernetesFiles...)
	sort.Strings(selected)
	selected = slices.Compact(selected)
	if len(selected) == 0 {
		return nil, nil, fmt.Errorf("checkov adapter requires at least one inventoried Terraform, Kubernetes, or container file")
	}
	if len(selected) > maxCheckovInputFiles {
		return nil, nil, fmt.Errorf("checkov adapter input exceeds %d inventoried files", maxCheckovInputFiles)
	}
	inventoried := make(map[string]bool, len(item.Files))
	for _, record := range item.Files {
		inventoried[record.Path] = true
	}
	absolute := make([]string, 0, len(selected))
	for _, path := range selected {
		if err := validateRelativePath(path); err != nil || !inventoried[path] || path == CheckovConfigSnapshotPath {
			return nil, nil, fmt.Errorf("checkov adapter inventory contains an invalid policy input path %q", path)
		}
		absolute = append(absolute, "/workspace/"+path)
	}
	document := checkovConfigDocument{
		Files: absolute, Frameworks: append([]string{}, checkovFrameworks...), Output: []string{"json"},
		Quiet: false, Compact: true, SoftFail: true, SkipDownload: true, SkipResultsUpload: true,
		DownloadExternalModules: false, EvaluateVariables: false,
	}
	payload, err := json.Marshal(document)
	if err != nil {
		return nil, nil, fmt.Errorf("encode scanner-owned Checkov configuration: %w", err)
	}
	return append(payload, '\n'), selected, nil
}

func readCheckovExpectedPaths(workspace string) ([]string, error) {
	path := filepath.Join(workspace, filepath.FromSlash(CheckovConfigSnapshotPath))
	info, err := os.Lstat(path)
	if err != nil || !info.Mode().IsRegular() || info.Size() < 2 || info.Size() > 1024*1024 {
		return nil, fmt.Errorf("sealed Checkov configuration is missing or invalid")
	}
	payload, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read sealed Checkov configuration: %w", err)
	}
	if err := rejectDuplicateKeys(payload); err != nil {
		return nil, fmt.Errorf("sealed Checkov configuration: %w", err)
	}
	decoder := json.NewDecoder(bytes.NewReader(payload))
	decoder.DisallowUnknownFields()
	var document checkovConfigDocument
	if err := decoder.Decode(&document); err != nil {
		return nil, fmt.Errorf("decode sealed Checkov configuration: %w", err)
	}
	var trailing any
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		return nil, fmt.Errorf("sealed Checkov configuration contains trailing content")
	}
	if !slices.Equal(document.Frameworks, checkovFrameworks) || !slices.Equal(document.Output, []string{"json"}) ||
		document.Quiet || !document.Compact || !document.SoftFail || !document.SkipDownload ||
		!document.SkipResultsUpload || document.DownloadExternalModules || document.EvaluateVariables ||
		len(document.Files) == 0 || len(document.Files) > maxCheckovInputFiles {
		return nil, fmt.Errorf("sealed Checkov configuration does not match the reviewed offline policy")
	}
	expected := make([]string, 0, len(document.Files))
	seen := map[string]bool{}
	for _, absolute := range document.Files {
		if !strings.HasPrefix(absolute, "/workspace/") {
			return nil, fmt.Errorf("sealed Checkov configuration contains an unbound file path")
		}
		relative := strings.TrimPrefix(absolute, "/workspace/")
		if err := validateRelativePath(relative); err != nil || seen[relative] {
			return nil, fmt.Errorf("sealed Checkov configuration contains an invalid or duplicate file path")
		}
		seen[relative] = true
		expected = append(expected, relative)
	}
	if !sort.StringsAreSorted(expected) {
		return nil, fmt.Errorf("sealed Checkov configuration file paths are not canonical")
	}
	return expected, nil
}

type checkovReport struct {
	CheckType string         `json:"check_type"`
	Results   checkovResults `json:"results"`
	Summary   checkovSummary `json:"summary"`
	URL       string         `json:"url,omitempty"`
}

type checkovResults struct {
	PassedChecks  []checkovRecord `json:"passed_checks"`
	FailedChecks  []checkovRecord `json:"failed_checks"`
	SkippedChecks []checkovRecord `json:"skipped_checks"`
	ParsingErrors []string        `json:"parsing_errors"`
}

type checkovSummary struct {
	Passed         int    `json:"passed"`
	Failed         int    `json:"failed"`
	Skipped        int    `json:"skipped"`
	ParsingErrors  int    `json:"parsing_errors"`
	ResourceCount  int    `json:"resource_count"`
	CheckovVersion string `json:"checkov_version"`
}

type checkovRecord struct {
	CheckID                   string             `json:"check_id"`
	BCCheckID                 json.RawMessage    `json:"bc_check_id"`
	CheckName                 string             `json:"check_name"`
	CheckResult               checkovCheckResult `json:"check_result"`
	CodeBlock                 json.RawMessage    `json:"code_block"`
	FilePath                  string             `json:"file_path"`
	FileAbsPath               string             `json:"file_abs_path"`
	RepoFilePath              string             `json:"repo_file_path"`
	FileLineRange             []int              `json:"file_line_range"`
	Resource                  string             `json:"resource"`
	Evaluations               json.RawMessage    `json:"evaluations"`
	CheckClass                string             `json:"check_class"`
	FixedDefinition           json.RawMessage    `json:"fixed_definition"`
	EntityTags                json.RawMessage    `json:"entity_tags"`
	CallerFilePath            json.RawMessage    `json:"caller_file_path"`
	CallerFileLineRange       json.RawMessage    `json:"caller_file_line_range"`
	ResourceAddress           json.RawMessage    `json:"resource_address"`
	Severity                  json.RawMessage    `json:"severity"`
	BCCategory                json.RawMessage    `json:"bc_category"`
	Benchmarks                json.RawMessage    `json:"benchmarks"`
	Description               json.RawMessage    `json:"description"`
	ShortDescription          json.RawMessage    `json:"short_description"`
	VulnerabilityDetails      json.RawMessage    `json:"vulnerability_details"`
	ConnectedNode             json.RawMessage    `json:"connected_node"`
	Guideline                 json.RawMessage    `json:"guideline"`
	Details                   []json.RawMessage  `json:"details"`
	CheckLen                  json.RawMessage    `json:"check_len"`
	DefinitionContextFilePath json.RawMessage    `json:"definition_context_file_path"`
	Breadcrumbs               json.RawMessage    `json:"breadcrumbs,omitempty"`
}

type checkovCheckResult struct {
	Result               string            `json:"result"`
	EvaluatedKeys        []json.RawMessage `json:"evaluated_keys,omitempty"`
	ResultsConfiguration json.RawMessage   `json:"results_configuration,omitempty"`
	SuppressComment      string            `json:"suppress_comment,omitempty"`
	Entity               json.RawMessage   `json:"entity,omitempty"`
}

type checkovNormalizedTool struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type checkovNormalizedFramework struct {
	Name      string `json:"name"`
	Passed    int    `json:"passed"`
	Failed    int    `json:"failed"`
	Resources int    `json:"resources"`
}

type checkovNormalizedSummary struct {
	Passed     int                          `json:"passed"`
	Failed     int                          `json:"failed"`
	Resources  int                          `json:"resources"`
	Frameworks []checkovNormalizedFramework `json:"frameworks"`
}

type checkovNormalizedFinding struct {
	CheckID   string `json:"check_id"`
	CheckName string `json:"check_name"`
	Framework string `json:"framework"`
	Path      string `json:"path"`
	LineStart int    `json:"line_start"`
	LineEnd   int    `json:"line_end"`
	Resource  string `json:"resource"`
}

type checkovNormalizedReport struct {
	Schema           string                     `json:"schema"`
	Tool             checkovNormalizedTool      `json:"tool"`
	Summary          checkovNormalizedSummary   `json:"summary"`
	ExpectedPaths    []string                   `json:"expected_paths"`
	AnalyzedPaths    []string                   `json:"analyzed_paths"`
	UnsupportedPaths []string                   `json:"unsupported_paths"`
	Findings         []checkovNormalizedFinding `json:"findings"`
}

func parseCheckovOutput(data []byte, maxRecords int, expectedPaths []string) (Transcript, map[string][]byte, error) {
	if len(bytes.TrimSpace(data)) == 0 {
		return Transcript{}, nil, fmt.Errorf("checkov output is empty")
	}
	if !utf8.Valid(data) {
		return Transcript{}, nil, fmt.Errorf("checkov output is not valid UTF-8")
	}
	if err := validateCheckovExpectedPaths(expectedPaths); err != nil {
		return Transcript{}, nil, err
	}
	if err := rejectDuplicateKeys(data); err != nil {
		return Transcript{}, nil, fmt.Errorf("checkov output: %w", err)
	}
	rawReports, err := decodeCheckovReports(data)
	if err != nil {
		return Transcript{}, nil, err
	}
	expected := make(map[string]bool, len(expectedPaths))
	for _, path := range expectedPaths {
		expected[path] = true
	}
	analyzed := map[string]bool{}
	findingsByKey := map[string]checkovNormalizedFinding{}
	frameworks := make([]checkovNormalizedFramework, 0, len(rawReports))
	seenFrameworks := map[string]bool{}
	totalRecords := 0
	for index, raw := range rawReports {
		report, summaryOnly, err := decodeCheckovReport(raw)
		if err != nil {
			return Transcript{}, nil, fmt.Errorf("checkov report %d: %w", index, err)
		}
		if summaryOnly {
			if len(rawReports) != 1 {
				return Transcript{}, nil, fmt.Errorf("checkov summary-only output cannot be combined with framework reports")
			}
			continue
		}
		if seenFrameworks[report.CheckType] {
			return Transcript{}, nil, fmt.Errorf("checkov output contains duplicate %s reports", report.CheckType)
		}
		seenFrameworks[report.CheckType] = true
		if err := validateCheckovReportEnvelope(report); err != nil {
			return Transcript{}, nil, fmt.Errorf("checkov %s report: %w", report.CheckType, err)
		}
		totalRecords += report.Summary.Passed + report.Summary.Failed + report.Summary.Skipped + report.Summary.ParsingErrors
		if totalRecords > maxRecords {
			return Transcript{}, nil, fmt.Errorf("checkov output exceeds %d policy records", maxRecords)
		}
		if report.Summary.ParsingErrors != 0 || report.Summary.Skipped != 0 {
			return Transcript{}, nil, fmt.Errorf("checkov report contains parsing errors or suppressed policy checks")
		}
		for _, record := range report.Results.PassedChecks {
			path, err := validateCheckovRecord(record, "PASSED", report.CheckType, expected)
			if err != nil {
				return Transcript{}, nil, err
			}
			analyzed[path] = true
		}
		for _, record := range report.Results.FailedChecks {
			path, err := validateCheckovRecord(record, "FAILED", report.CheckType, expected)
			if err != nil {
				return Transcript{}, nil, err
			}
			analyzed[path] = true
			finding := checkovNormalizedFinding{
				CheckID: record.CheckID, CheckName: record.CheckName, Framework: report.CheckType,
				Path: path, LineStart: record.FileLineRange[0], LineEnd: record.FileLineRange[1], Resource: record.Resource,
			}
			keyBytes, _ := json.Marshal(finding)
			keyDigest := sha256.Sum256(keyBytes)
			findingsByKey[hex.EncodeToString(keyDigest[:])] = finding
		}
		frameworks = append(frameworks, checkovNormalizedFramework{
			Name: report.CheckType, Passed: report.Summary.Passed,
			Failed: report.Summary.Failed, Resources: report.Summary.ResourceCount,
		})
	}
	sort.Slice(frameworks, func(i, j int) bool { return frameworks[i].Name < frameworks[j].Name })
	analyzedPaths := sortedCheckovKeys(analyzed)
	unsupportedPaths := make([]string, 0)
	for _, path := range expectedPaths {
		if !analyzed[path] {
			unsupportedPaths = append(unsupportedPaths, path)
		}
	}
	findings := make([]checkovNormalizedFinding, 0, len(findingsByKey))
	for _, finding := range findingsByKey {
		findings = append(findings, finding)
	}
	sort.Slice(findings, func(i, j int) bool {
		left, right := findings[i], findings[j]
		if left.Path != right.Path {
			return left.Path < right.Path
		}
		if left.LineStart != right.LineStart {
			return left.LineStart < right.LineStart
		}
		if left.CheckID != right.CheckID {
			return left.CheckID < right.CheckID
		}
		if left.Resource != right.Resource {
			return left.Resource < right.Resource
		}
		return left.Framework < right.Framework
	})
	summary := checkovNormalizedSummary{Frameworks: frameworks}
	for _, framework := range frameworks {
		summary.Passed += framework.Passed
		summary.Failed += framework.Failed
		summary.Resources += framework.Resources
	}
	report := checkovNormalizedReport{
		Schema: "prc.checkov-iac-policy-report/v1", Tool: checkovNormalizedTool{Name: "checkov", Version: CheckovToolVersion},
		Summary: summary, ExpectedPaths: append([]string{}, expectedPaths...), AnalyzedPaths: analyzedPaths,
		UnsupportedPaths: unsupportedPaths, Findings: findings,
	}
	canonical, err := json.Marshal(report)
	if err != nil {
		return Transcript{}, nil, fmt.Errorf("encode normalized Checkov report: %w", err)
	}
	canonical = append(canonical, '\n')
	digest := sha256.Sum256(canonical)
	descriptor := "sha256:" + hex.EncodeToString(digest[:])
	observations := make([]Observation, 0, max(1, len(findings)+len(unsupportedPaths)))
	for _, finding := range findings {
		findingBytes, _ := json.Marshal(finding)
		findingDigest := sha256.Sum256(findingBytes)
		observations = append(observations, Observation{
			ID: hex.EncodeToString(findingDigest[:]), Kind: CheckovObservationKind, Outcome: "found",
			Summary:   fmt.Sprintf("%s: %s", finding.CheckID, finding.CheckName),
			Locations: []Location{{Path: finding.Path, Line: finding.LineStart}},
			Data: map[string]any{
				"artifact_digest": descriptor, "check_id": finding.CheckID, "framework": finding.Framework,
				"line_end": finding.LineEnd, "resource": finding.Resource, "tool_version": CheckovToolVersion,
			},
		})
	}
	for _, path := range unsupportedPaths {
		pathDigest := sha256.Sum256([]byte("checkov-unsupported\x00" + path))
		observations = append(observations, Observation{
			ID: hex.EncodeToString(pathDigest[:]), Kind: CheckovObservationKind, Outcome: "unsupported",
			Summary:   "Checkov produced no evaluated policy record for an inventoried IaC or container file.",
			Locations: []Location{{Path: path}},
			Data:      map[string]any{"artifact_digest": descriptor, "reason": "no_evaluated_records", "tool_version": CheckovToolVersion},
		})
	}
	if len(observations) == 0 {
		locations := make([]Location, 0, len(analyzedPaths))
		for _, path := range analyzedPaths {
			locations = append(locations, Location{Path: path})
		}
		observations = append(observations, Observation{
			ID: hex.EncodeToString(digest[:]), Kind: CheckovObservationKind, Outcome: "not_found",
			Summary: "Checkov found no policy violations in the inventoried IaC and container files.", Locations: locations,
			Data: map[string]any{"analyzed_file_count": len(analyzedPaths), "artifact_digest": descriptor, "tool_version": CheckovToolVersion},
		})
	}
	transcript := Transcript{
		Logs: []Log{}, Observations: observations,
		Artifacts: []Artifact{{ID: checkovArtifactID, MediaType: checkovArtifactMediaType, Digest: descriptor, Size: int64(len(canonical)), Path: "iac-policy-report.json"}},
		Summary:   Summary{Type: "summary", Status: "completed", Counts: map[string]int{"logs": 0, "observations": len(observations), "artifacts": 1}},
	}
	return transcript, map[string][]byte{descriptor: canonical}, nil
}

func validateCheckovExpectedPaths(paths []string) error {
	if len(paths) == 0 || len(paths) > maxCheckovInputFiles || !sort.StringsAreSorted(paths) {
		return fmt.Errorf("checkov output parsing requires a canonical nonempty expected file set")
	}
	for index, path := range paths {
		if err := validateRelativePath(path); err != nil || (index > 0 && paths[index-1] == path) {
			return fmt.Errorf("checkov output parsing received an invalid or duplicate expected path")
		}
	}
	return nil
}

func decodeCheckovReports(data []byte) ([]json.RawMessage, error) {
	decoder := json.NewDecoder(bytes.NewReader(data))
	var document json.RawMessage
	if err := decoder.Decode(&document); err != nil {
		return nil, fmt.Errorf("decode Checkov output: %w", err)
	}
	var trailing any
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		return nil, fmt.Errorf("Checkov output contains trailing content")
	}
	trimmed := bytes.TrimSpace(document)
	if len(trimmed) == 0 {
		return nil, fmt.Errorf("Checkov output is empty")
	}
	if trimmed[0] == '{' {
		return []json.RawMessage{document}, nil
	}
	if trimmed[0] != '[' {
		return nil, fmt.Errorf("Checkov output must be one report object or report array")
	}
	var reports []json.RawMessage
	if err := json.Unmarshal(document, &reports); err != nil || reports == nil || len(reports) > len(checkovFrameworks) {
		return nil, fmt.Errorf("Checkov output report array is invalid")
	}
	return reports, nil
}

func decodeCheckovReport(raw json.RawMessage) (checkovReport, bool, error) {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(raw, &fields); err != nil || fields == nil {
		return checkovReport{}, false, fmt.Errorf("report must be an object")
	}
	if len(fields) == 1 && fields["summary"] != nil {
		var summary checkovSummary
		if err := decodeCheckovStrict(fields["summary"], &summary); err != nil {
			return checkovReport{}, false, fmt.Errorf("decode summary-only report: %w", err)
		}
		if summary != (checkovSummary{CheckovVersion: CheckovToolVersion}) {
			return checkovReport{}, false, fmt.Errorf("summary-only report must describe zero evaluated resources")
		}
		return checkovReport{Summary: summary}, true, nil
	}
	for _, required := range []string{"check_type", "results", "summary"} {
		if fields[required] == nil {
			return checkovReport{}, false, fmt.Errorf("report omits required property %q", required)
		}
	}
	var report checkovReport
	if err := decodeCheckovStrict(raw, &report); err != nil {
		return checkovReport{}, false, fmt.Errorf("decode framework report: %w", err)
	}
	return report, false, nil
}

func decodeCheckovStrict(data []byte, destination any) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(destination); err != nil {
		return err
	}
	var trailing any
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		return fmt.Errorf("unexpected trailing JSON content")
	}
	return nil
}

func validateCheckovReportEnvelope(report checkovReport) error {
	if !slices.Contains(checkovFrameworks, report.CheckType) || report.Summary.CheckovVersion != CheckovToolVersion ||
		report.Summary.Passed < 0 || report.Summary.Failed < 0 || report.Summary.Skipped < 0 ||
		report.Summary.ParsingErrors < 0 || report.Summary.ResourceCount < 0 {
		return fmt.Errorf("report summary or tool identity is invalid")
	}
	if report.URL != "" && report.URL != checkovReportURL {
		return fmt.Errorf("report contains an unexpected service URL")
	}
	if report.Results.PassedChecks == nil || report.Results.FailedChecks == nil ||
		report.Results.SkippedChecks == nil || report.Results.ParsingErrors == nil {
		return fmt.Errorf("report result collections must be arrays")
	}
	if report.Summary.Passed != len(report.Results.PassedChecks) || report.Summary.Failed != len(report.Results.FailedChecks) ||
		report.Summary.Skipped != len(report.Results.SkippedChecks) || report.Summary.ParsingErrors != len(report.Results.ParsingErrors) {
		return fmt.Errorf("report summary counts do not match its result collections")
	}
	for _, path := range report.Results.ParsingErrors {
		if !validBoundedText(path, 16*1024) {
			return fmt.Errorf("report contains an invalid parsing-error path")
		}
	}
	if report.Summary.Passed+report.Summary.Failed > 0 && report.Summary.ResourceCount == 0 {
		return fmt.Errorf("report evaluates policies without any parsed resources")
	}
	return nil
}

func validateCheckovRecord(record checkovRecord, outcome, framework string, expected map[string]bool) (string, error) {
	if !checkovCheckIDPattern.MatchString(record.CheckID) || !validBoundedText(record.CheckName, 16*1024) ||
		!validBoundedText(record.CheckClass, 16*1024) || !strings.HasPrefix(record.CheckClass, "checkov.") ||
		!validBoundedText(record.Resource, 16*1024) || record.CheckResult.Result != outcome ||
		record.CheckResult.SuppressComment != "" || len(record.CheckResult.EvaluatedKeys) > 256 ||
		len(record.CheckResult.ResultsConfiguration) > 64*1024 || len(record.CheckResult.Entity) > 1024*1024 ||
		!checkovObjectOrNull(record.CheckResult.Entity) {
		return "", fmt.Errorf("checkov %s record %q has an invalid policy identity or result", framework, record.CheckID)
	}
	for _, key := range record.CheckResult.EvaluatedKeys {
		var text string
		if err := json.Unmarshal(key, &text); err != nil || !validBoundedText(text, 16*1024) {
			return "", fmt.Errorf("checkov record %q has an invalid evaluated key", record.CheckID)
		}
	}
	if !checkovNull(record.BCCheckID) || !checkovNull(record.CodeBlock) || !checkovNull(record.FixedDefinition) ||
		!checkovNull(record.EntityTags) || !checkovNull(record.ResourceAddress) || !checkovNull(record.Severity) ||
		!checkovNull(record.BCCategory) || !checkovEmptyObjectOrNull(record.Benchmarks) || !checkovNull(record.Description) ||
		!checkovNull(record.ShortDescription) || !checkovNull(record.VulnerabilityDetails) || !checkovNull(record.ConnectedNode) ||
		!checkovNull(record.Guideline) || !checkovNull(record.CheckLen) || len(record.Details) != 0 ||
		len(record.Evaluations) > 64*1024 || !checkovObjectOrNull(record.Evaluations) {
		return "", fmt.Errorf("checkov record %q contains unexpected online, code, or policy metadata", record.CheckID)
	}
	if len(record.FileLineRange) != 2 || record.FileLineRange[0] < 1 || record.FileLineRange[1] < record.FileLineRange[0] {
		return "", fmt.Errorf("checkov record %q has an invalid source range", record.CheckID)
	}
	path, err := normalizeCheckovPath(record.FileAbsPath)
	if err != nil || !expected[path] || record.RepoFilePath != "/workspace/"+path {
		return "", fmt.Errorf("checkov record %q refers to an unbound file", record.CheckID)
	}
	filePath, err := normalizeCheckovPath(record.FilePath)
	if err != nil || filePath != path && filePath != filepath.Base(path) || !checkovFrameworkMatchesPath(framework, path) {
		return "", fmt.Errorf("checkov record %q has inconsistent file or framework identity", record.CheckID)
	}
	if err := validateOptionalCheckovPath(record.CallerFilePath, expected); err != nil ||
		err == nil && !checkovOptionalRange(record.CallerFileLineRange) ||
		validateOptionalCheckovPath(record.DefinitionContextFilePath, expected) != nil {
		return "", fmt.Errorf("checkov record %q has invalid source context", record.CheckID)
	}
	return path, nil
}

func normalizeCheckovPath(value string) (string, error) {
	var relative string
	if strings.HasPrefix(value, "/workspace/") {
		relative = strings.TrimPrefix(value, "/workspace/")
	} else if strings.HasPrefix(value, "/") {
		relative = strings.TrimPrefix(value, "/")
	} else {
		return "", fmt.Errorf("path is not rooted in the sealed workspace")
	}
	if err := validateRelativePath(relative); err != nil {
		return "", err
	}
	return relative, nil
}

func validateOptionalCheckovPath(raw json.RawMessage, expected map[string]bool) error {
	if checkovNull(raw) {
		return nil
	}
	var value string
	if err := json.Unmarshal(raw, &value); err != nil {
		return err
	}
	path, err := normalizeCheckovPath(value)
	if err != nil || !expected[path] {
		return fmt.Errorf("path is outside the expected file set")
	}
	return nil
}

func checkovOptionalRange(raw json.RawMessage) bool {
	if checkovNull(raw) {
		return true
	}
	var value []int
	return json.Unmarshal(raw, &value) == nil && len(value) == 2 && value[0] >= 1 && value[1] >= value[0]
}

func checkovFrameworkMatchesPath(framework, path string) bool {
	base := strings.ToLower(filepath.Base(path))
	switch framework {
	case "terraform", "terraform_json":
		return strings.HasSuffix(strings.ToLower(path), ".tf") || strings.HasSuffix(strings.ToLower(path), ".tf.json")
	case "kubernetes":
		return strings.HasSuffix(base, ".yaml") || strings.HasSuffix(base, ".yml") || strings.HasSuffix(base, ".json")
	case "dockerfile":
		return base == "dockerfile" || strings.HasPrefix(base, "dockerfile.")
	default:
		return false
	}
}

func checkovNull(raw json.RawMessage) bool {
	return len(raw) == 0 || bytes.Equal(bytes.TrimSpace(raw), []byte("null"))
}
func checkovEmptyObjectOrNull(raw json.RawMessage) bool {
	return checkovNull(raw) || bytes.Equal(bytes.TrimSpace(raw), []byte("{}"))
}
func checkovObjectOrNull(raw json.RawMessage) bool {
	trimmed := bytes.TrimSpace(raw)
	return checkovNull(raw) || len(trimmed) >= 2 && trimmed[0] == '{' && trimmed[len(trimmed)-1] == '}'
}
func sortedCheckovKeys(values map[string]bool) []string {
	result := make([]string, 0, len(values))
	for value := range values {
		result = append(result, value)
	}
	sort.Strings(result)
	return result
}
