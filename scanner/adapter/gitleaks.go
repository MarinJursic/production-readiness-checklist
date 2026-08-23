package adapter

import (
	"bytes"
	"crypto/sha256"
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"
)

const (
	GitleaksProtocolVersion     = "prc-adapter-gitleaks-json-v1"
	GitleaksOutputSchemaVersion = "gitleaks.report/v8.30.0"
	GitleaksToolVersion         = "8.30.0"
	GitleaksObservationKind     = "secret-detection"
	GitleaksImage               = "ghcr.io/gitleaks/gitleaks@sha256:691af3c7c5a48b16f187ce3446d5f194838f91238f27270ed36eef6359a574d9"
	GitleaksConfigSHA256        = "e163e53b9e7e8a8511e77271e2b323ed057759542a6d988258afe3a1fa329caf"
	gitleaksReportMediaType     = "application/vnd.gitleaks.report+json;version=8.30.0;redacted=100"
	gitleaksIgnoreSourcePath    = ".gitleaksignore"
	gitleaksIgnoreSnapshotPath  = ".prc/gitleaksignore-source"
)

var gitleaksRuleIDPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]{0,127}$`)

//go:embed gitleaks-v8.30.0.toml
var gitleaksConfig []byte

func gitleaksCommand() []string {
	return []string{
		"dir", "--no-banner", "--no-color", "--log-level=fatal",
		"--config=/dev/stdin", "--gitleaks-ignore-path=/dev/null", "--ignore-gitleaks-allow",
		"--max-archive-depth=0", "--max-decode-depth=0", "--max-target-megabytes=10",
		"--timeout=110", "--redact=100", "--report-format=json", "--report-path=-",
		"--exit-code=0", "/workspace",
	}
}

func validateGitleaksManifest(manifest Manifest) error {
	if manifest.Protocol != GitleaksProtocolVersion || manifest.OutputSchema != GitleaksOutputSchemaVersion {
		return fmt.Errorf("gitleaks adapter requires its exact protocol and output schema")
	}
	if manifest.Image != GitleaksImage || !slices.Equal(manifest.Command, gitleaksCommand()) {
		return fmt.Errorf("gitleaks adapter requires the reviewed immutable image and scanner-owned command")
	}
	if manifest.Tool.Name != "gitleaks" || manifest.Tool.Version != GitleaksToolVersion ||
		manifest.Tool.Upstream != "https://github.com/gitleaks/gitleaks" ||
		len(manifest.Tool.Formats) != 1 || manifest.Tool.Formats[0].Name != "gitleaks-json" ||
		!slices.Equal(manifest.Tool.Formats[0].Versions, []string{GitleaksToolVersion}) {
		return fmt.Errorf("gitleaks adapter tool identity does not match the reviewed normalizer")
	}
	if !slices.Equal(manifest.ObservationKinds, []string{GitleaksObservationKind}) {
		return fmt.Errorf("gitleaks adapter must declare only %s observations", GitleaksObservationKind)
	}
	if manifest.Capabilities.WriteScratch || !manifest.Capabilities.ChildProcesses {
		return fmt.Errorf("gitleaks adapter requires no scratch and a bounded OS-task allowance")
	}
	if manifest.Resources.PIDs < 16 || manifest.Resources.PIDs > 256 ||
		manifest.Resources.TimeoutSeconds < 111 {
		return fmt.Errorf("gitleaks adapter resource limits cannot support or cannot safely bound the reviewed command")
	}
	return nil
}

// ExecutionInput returns only scanner-owned input for the manifest protocol.
// JSONL adapters receive their sealed request; Gitleaks receives the embedded,
// digest-checked upstream ruleset on stdin so target configuration cannot
// weaken the scan.
func ExecutionInput(
	manifest Manifest,
	runID string,
	subject Subject,
	facts, config map[string]any,
) ([]byte, error) {
	if err := manifest.ValidateForCurrentEngine(); err != nil {
		return nil, err
	}
	switch manifest.Protocol {
	case ProtocolVersion:
		return InputJSONL(runID, subject, facts, config)
	case GitleaksProtocolVersion:
		if err := validateInputIdentity(runID, subject); err != nil {
			return nil, err
		}
		digest := sha256.Sum256(gitleaksConfig)
		if hex.EncodeToString(digest[:]) != GitleaksConfigSHA256 {
			return nil, fmt.Errorf("embedded gitleaks configuration digest is invalid")
		}
		return append([]byte(nil), gitleaksConfig...), nil
	default:
		return nil, fmt.Errorf("unsupported adapter input protocol %q", manifest.Protocol)
	}
}

// ParseManifestOutput converts one exact external format into the authority-
// free PRC transcript. Generic JSONL remains the extension protocol; the only
// native format is the reviewed Gitleaks 8.30.0 JSON contract.
func ParseManifestOutput(manifest Manifest, input io.Reader) (Transcript, error) {
	if err := manifest.ValidateForCurrentEngine(); err != nil {
		return Transcript{}, err
	}
	switch manifest.Protocol {
	case ProtocolVersion:
		return ParseOutput(input, manifest.Resources.Limits)
	case GitleaksProtocolVersion:
		data, err := io.ReadAll(io.LimitReader(input, int64(manifest.Resources.MaxStdout)+1))
		if err != nil {
			return Transcript{}, fmt.Errorf("read gitleaks output: %w", err)
		}
		if len(data) > manifest.Resources.MaxStdout {
			return Transcript{}, fmt.Errorf("gitleaks output exceeds %d bytes", manifest.Resources.MaxStdout)
		}
		return parseGitleaksOutput(data, manifest.Resources.MaxMessages)
	default:
		return Transcript{}, fmt.Errorf("unsupported adapter output protocol %q", manifest.Protocol)
	}
}

type gitleaksFinding struct {
	RuleID      string   `json:"RuleID"`
	Description string   `json:"Description"`
	StartLine   int      `json:"StartLine"`
	EndLine     int      `json:"EndLine"`
	StartColumn int      `json:"StartColumn"`
	EndColumn   int      `json:"EndColumn"`
	Match       string   `json:"Match"`
	Secret      string   `json:"Secret"`
	File        string   `json:"File"`
	SymlinkFile string   `json:"SymlinkFile"`
	Commit      string   `json:"Commit"`
	Entropy     float64  `json:"Entropy"`
	Author      string   `json:"Author"`
	Email       string   `json:"Email"`
	Date        string   `json:"Date"`
	Message     string   `json:"Message"`
	Tags        []string `json:"Tags"`
	Fingerprint string   `json:"Fingerprint"`
}

func parseGitleaksOutput(data []byte, maxFindings int) (Transcript, error) {
	if len(bytes.TrimSpace(data)) == 0 {
		return Transcript{}, fmt.Errorf("gitleaks output is empty")
	}
	if err := rejectDuplicateKeys(data); err != nil {
		return Transcript{}, fmt.Errorf("gitleaks output: %w", err)
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	var findings []gitleaksFinding
	if err := decoder.Decode(&findings); err != nil {
		return Transcript{}, fmt.Errorf("decode gitleaks report: %w", err)
	}
	var trailing any
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		if err == nil {
			return Transcript{}, fmt.Errorf("gitleaks output contains more than one JSON value")
		}
		return Transcript{}, fmt.Errorf("decode trailing gitleaks output: %w", err)
	}
	if findings == nil {
		return Transcript{}, fmt.Errorf("gitleaks report must be a JSON array")
	}
	if len(findings) > maxFindings {
		return Transcript{}, fmt.Errorf("gitleaks report exceeds %d findings", maxFindings)
	}

	observations := make([]Observation, 0, max(1, len(findings)))
	seenIDs := map[string]bool{}
	for _, finding := range findings {
		observation, err := normalizeGitleaksFinding(finding)
		if err != nil {
			return Transcript{}, err
		}
		if seenIDs[observation.ID] {
			return Transcript{}, fmt.Errorf("gitleaks report contains duplicate normalized findings")
		}
		seenIDs[observation.ID] = true
		observations = append(observations, observation)
	}
	if len(observations) == 0 {
		observations = append(observations, Observation{
			ID: "gitleaks-no-findings", Kind: GitleaksObservationKind, Outcome: "not_found",
			Summary:   "Gitleaks found no potential secrets in the bounded current-tree scan.",
			Locations: []Location{}, Data: map[string]any{"tool_version": GitleaksToolVersion},
		})
	}
	sort.Slice(observations, func(left, right int) bool {
		leftLocation, rightLocation := "", ""
		if len(observations[left].Locations) > 0 {
			leftLocation = observations[left].Locations[0].Path
		}
		if len(observations[right].Locations) > 0 {
			rightLocation = observations[right].Locations[0].Path
		}
		if leftLocation != rightLocation {
			return leftLocation < rightLocation
		}
		return observations[left].ID < observations[right].ID
	})
	reportDigest := sha256.Sum256(data)
	artifacts := []Artifact{{
		ID: "gitleaks-redacted-report", MediaType: gitleaksReportMediaType,
		Digest: "sha256:" + hex.EncodeToString(reportDigest[:]), Size: int64(len(data)),
	}}
	return Transcript{
		Logs: []Log{}, Observations: observations, Artifacts: artifacts,
		Summary: Summary{Type: "summary", Status: "completed", Counts: map[string]int{
			"logs": 0, "observations": len(observations), "artifacts": len(artifacts),
		}},
	}, nil
}

func normalizeGitleaksFinding(finding gitleaksFinding) (Observation, error) {
	if !gitleaksRuleIDPattern.MatchString(finding.RuleID) || strings.TrimSpace(finding.Description) == "" ||
		len(finding.Description) > 4096 {
		return Observation{}, fmt.Errorf("gitleaks finding has an invalid rule identity or description")
	}
	if finding.StartLine < 1 || finding.EndLine < finding.StartLine || finding.StartColumn < 1 ||
		finding.EndColumn < finding.StartColumn {
		return Observation{}, fmt.Errorf("gitleaks finding has invalid source coordinates")
	}
	if finding.Secret != "REDACTED" || !strings.Contains(finding.Match, "REDACTED") || len(finding.Match) > 8192 {
		return Observation{}, fmt.Errorf("gitleaks finding is not fully redacted")
	}
	if finding.SymlinkFile != "" || finding.Commit != "" || finding.Author != "" || finding.Email != "" ||
		finding.Date != "" || finding.Message != "" {
		return Observation{}, fmt.Errorf("gitleaks current-tree report contains unexpected symlink or history metadata")
	}
	if math.IsNaN(finding.Entropy) || math.IsInf(finding.Entropy, 0) || finding.Entropy < 0 {
		return Observation{}, fmt.Errorf("gitleaks finding has invalid entropy")
	}
	if finding.Tags == nil || len(finding.Tags) > 128 {
		return Observation{}, fmt.Errorf("gitleaks finding tags must be a bounded array")
	}
	for _, tag := range finding.Tags {
		if strings.TrimSpace(tag) == "" || len(tag) > 256 {
			return Observation{}, fmt.Errorf("gitleaks finding contains an invalid tag")
		}
	}
	path, err := normalizeGitleaksPath(finding.File)
	if err != nil {
		return Observation{}, err
	}
	expectedFingerprint := finding.File + ":" + finding.RuleID + ":" + strconv.Itoa(finding.StartLine)
	if finding.Fingerprint != expectedFingerprint {
		return Observation{}, fmt.Errorf("gitleaks finding fingerprint does not match its location")
	}
	identity := sha256.Sum256([]byte(strings.Join([]string{
		path, finding.RuleID, strconv.Itoa(finding.StartLine), strconv.Itoa(finding.StartColumn),
		strconv.Itoa(finding.EndLine), strconv.Itoa(finding.EndColumn),
	}, "\x00")))
	return Observation{
		ID: "gitleaks-" + hex.EncodeToString(identity[:]), Kind: GitleaksObservationKind, Outcome: "found",
		Summary:   "Gitleaks detected a potential secret matching rule " + finding.RuleID + ".",
		Locations: []Location{{Path: path, Line: finding.StartLine, Column: finding.StartColumn}},
		Data: map[string]any{
			"rule_id": finding.RuleID, "end_line": finding.EndLine, "end_column": finding.EndColumn,
			"tool_version": GitleaksToolVersion,
		},
	}, nil
}

func normalizeGitleaksPath(path string) (string, error) {
	const prefix = "/workspace/"
	if !strings.HasPrefix(path, prefix) || len(path) <= len(prefix) || len(path) > 4096 {
		return "", fmt.Errorf("gitleaks finding path is outside the assessment snapshot")
	}
	relative := strings.TrimPrefix(path, prefix)
	if relative == gitleaksIgnoreSnapshotPath {
		relative = gitleaksIgnoreSourcePath
	}
	if err := validateRelativePath(relative); err != nil {
		return "", fmt.Errorf("gitleaks finding path: %w", err)
	}
	return relative, nil
}
