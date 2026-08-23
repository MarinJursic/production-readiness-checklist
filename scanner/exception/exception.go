// Package exception validates signed, scoped, expiring human risk decisions.
// An accepted exception never rewrites a scanner assertion or terminal state.
package exception

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
	"github.com/MarinJursic/production-readiness-checklist/scanner/trust"
	"gopkg.in/yaml.v3"
)

const (
	Schema       = "prc.risk-exception/v0.1"
	ReportSchema = "prc.risk-exception-verification/v0.1"
)

const maximumRecordBytes = 1024 * 1024

var (
	exceptionIDPattern = regexp.MustCompile(`^PRC-EXC-[A-Z0-9][A-Z0-9-]{2,63}$`)
	actorIDPattern     = regexp.MustCompile(`^[a-z0-9][a-z0-9._-]{0,127}$`)
	profileIDPattern   = regexp.MustCompile(`^prc/[a-z0-9-]+$`)
	assertionIDPattern = regexp.MustCompile(`^PRC-A-[A-Z0-9]+-[0-9]{3}$`)
	controlIDPattern   = regexp.MustCompile(`^(?:USEQ-[A-F0-9]{8}|PRC-[0-9]{2}-[0-9]{3})$`)
	digestPattern      = regexp.MustCompile(`^[0-9a-f]{64}$`)
	revisionPattern    = regexp.MustCompile(`^(?:[0-9a-f]{40,64})?$`)
)

type Actor struct {
	ID        string `json:"id" yaml:"id"`
	Name      string `json:"name" yaml:"name"`
	Authority string `json:"authority" yaml:"authority"`
}

type RunBinding struct {
	RunID              string   `json:"run_id" yaml:"run_id"`
	InventoryDigest    string   `json:"inventory_digest" yaml:"inventory_digest"`
	ProfileID          string   `json:"profile_id" yaml:"profile_id"`
	ProfileVersion     string   `json:"profile_version" yaml:"profile_version"`
	TargetName         string   `json:"target_name" yaml:"target_name"`
	TargetCommit       string   `json:"target_commit" yaml:"target_commit"`
	ProjectID          string   `json:"project_id" yaml:"project_id"`
	ArtifactDigests    []string `json:"artifact_digests" yaml:"artifact_digests"`
	TargetEnvironments []string `json:"target_environments" yaml:"target_environments"`
}

type FindingBinding struct {
	FindingID          string   `json:"finding_id" yaml:"finding_id"`
	FindingFingerprint string   `json:"finding_fingerprint" yaml:"finding_fingerprint"`
	AssertionID        string   `json:"assertion_id" yaml:"assertion_id"`
	ControlIDs         []string `json:"control_ids" yaml:"control_ids"`
}

type RiskAnalysis struct {
	Title                string `json:"title" yaml:"title"`
	Rationale            string `json:"rationale" yaml:"rationale"`
	Likelihood           string `json:"likelihood" yaml:"likelihood"`
	Impact               string `json:"impact" yaml:"impact"`
	WorstCredibleOutcome string `json:"worst_credible_outcome" yaml:"worst_credible_outcome"`
}

type CompensatingControl struct {
	Description        string   `json:"description" yaml:"description"`
	EvidenceReferences []string `json:"evidence_references" yaml:"evidence_references"`
}

type Monitoring struct {
	Owner    string `json:"owner" yaml:"owner"`
	Signal   string `json:"signal" yaml:"signal"`
	Response string `json:"response" yaml:"response"`
}

type Remediation struct {
	Owner string    `json:"owner" yaml:"owner"`
	Plan  string    `json:"plan" yaml:"plan"`
	DueAt time.Time `json:"due_at" yaml:"due_at"`
}

type Record struct {
	SchemaVersion        string                `json:"schema_version" yaml:"schema_version"`
	ID                   string                `json:"id" yaml:"id"`
	Status               string                `json:"status" yaml:"status"`
	Run                  RunBinding            `json:"run" yaml:"run"`
	Finding              FindingBinding        `json:"finding" yaml:"finding"`
	RequestedBy          Actor                 `json:"requested_by" yaml:"requested_by"`
	RiskOwner            Actor                 `json:"risk_owner" yaml:"risk_owner"`
	Reviewers            []Actor               `json:"reviewers" yaml:"reviewers"`
	Risk                 RiskAnalysis          `json:"risk" yaml:"risk"`
	CompensatingControls []CompensatingControl `json:"compensating_controls" yaml:"compensating_controls"`
	Monitoring           Monitoring            `json:"monitoring" yaml:"monitoring"`
	Remediation          Remediation           `json:"remediation" yaml:"remediation"`
	ApprovedAt           time.Time             `json:"approved_at" yaml:"approved_at"`
	ExpiresAt            time.Time             `json:"expires_at" yaml:"expires_at"`
}

type Loaded struct {
	Record Record
	Digest string
}

type Verification struct {
	SchemaVersion   string             `json:"schema_version"`
	Exception       Record             `json:"exception"`
	ExceptionDigest string             `json:"exception_digest"`
	VerifiedAt      time.Time          `json:"verified_at"`
	Signature       trust.Verification `json:"signature"`
	Disposition     string             `json:"disposition"`
	GateEffect      string             `json:"gate_effect"`
}

func Load(path string) (Loaded, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return Loaded{}, fmt.Errorf("inspect risk exception: %w", err)
	}
	if info.Mode()&os.ModeSymlink != 0 || !info.Mode().IsRegular() || info.Size() > maximumRecordBytes {
		return Loaded{}, fmt.Errorf("risk exception must be a non-symlink regular file no larger than 1 MiB")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return Loaded{}, fmt.Errorf("read risk exception: %w", err)
	}
	if err := rejectAmbiguousYAML(data); err != nil {
		return Loaded{}, err
	}
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)
	var record Record
	if err := decoder.Decode(&record); err != nil {
		return Loaded{}, fmt.Errorf("decode risk exception: %w", err)
	}
	var trailing any
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		if err == nil {
			return Loaded{}, fmt.Errorf("risk exception contains more than one YAML document")
		}
		return Loaded{}, fmt.Errorf("decode trailing risk exception content: %w", err)
	}
	if err := normalizeAndValidate(&record); err != nil {
		return Loaded{}, err
	}
	digest, err := recordDigest(record)
	if err != nil {
		return Loaded{}, err
	}
	return Loaded{Record: record, Digest: digest}, nil
}

func Verify(
	loaded Loaded,
	run model.RunResult,
	store trust.LoadedStore,
	signature trust.Signature,
	verifiedAt time.Time,
) (Verification, error) {
	record := loaded.Record
	if err := normalizeAndValidate(&record); err != nil {
		return Verification{}, err
	}
	digest, err := recordDigest(record)
	if err != nil {
		return Verification{}, err
	}
	if digest != loaded.Digest {
		return Verification{}, fmt.Errorf("risk exception digest does not match its canonical content")
	}
	if !utcTime(verifiedAt) || verifiedAt.Before(record.ApprovedAt) || !verifiedAt.Before(record.ExpiresAt) {
		return Verification{}, fmt.Errorf("risk exception is not current at the UTC verification time")
	}
	if signature.IssuedAt.Before(record.ApprovedAt) || signature.IssuedAt.After(verifiedAt) {
		return Verification{}, fmt.Errorf("risk exception signature time is inconsistent with approval or verification")
	}
	if err := validateRunBinding(record, run); err != nil {
		return Verification{}, err
	}
	signatureVerification, err := trust.Verify(
		store, signature, "risk-exception", record.ID, loaded.Digest, verifiedAt,
	)
	if err != nil {
		return Verification{}, err
	}
	return Verification{
		SchemaVersion: ReportSchema, Exception: record, ExceptionDigest: loaded.Digest,
		VerifiedAt: verifiedAt.UTC(), Signature: signatureVerification,
		Disposition: "accepted_risk_exception",
		GateEffect:  "finding remains failed and the scanner terminal state is unchanged",
	}, nil
}

func normalizeAndValidate(record *Record) error {
	if record.SchemaVersion != Schema || !exceptionIDPattern.MatchString(record.ID) || record.Status != "approved" {
		return fmt.Errorf("risk exception requires a supported schema, valid ID, and approved status")
	}
	if !digestPattern.MatchString(record.Run.RunID) || !digestPattern.MatchString(record.Run.InventoryDigest) ||
		!profileIDPattern.MatchString(record.Run.ProfileID) || !cleanText(record.Run.ProfileVersion, 64) ||
		!cleanText(record.Run.TargetName, 255) || !revisionPattern.MatchString(record.Run.TargetCommit) ||
		(record.Run.ProjectID != "" && !actorIDPattern.MatchString(record.Run.ProjectID)) {
		return fmt.Errorf("risk exception run binding is invalid")
	}
	if !digestPattern.MatchString(record.Finding.FindingID) || !digestPattern.MatchString(record.Finding.FindingFingerprint) ||
		!assertionIDPattern.MatchString(record.Finding.AssertionID) || len(record.Finding.ControlIDs) == 0 {
		return fmt.Errorf("risk exception finding binding is invalid")
	}
	for _, controlID := range record.Finding.ControlIDs {
		if !controlIDPattern.MatchString(controlID) {
			return fmt.Errorf("risk exception has an invalid control ID")
		}
	}
	if err := normalizeUnique(&record.Finding.ControlIDs, 256, controlIDPattern); err != nil {
		return fmt.Errorf("risk exception control IDs: %w", err)
	}
	if err := normalizeDigests(&record.Run.ArtifactDigests); err != nil {
		return fmt.Errorf("risk exception artifact digests: %w", err)
	}
	if err := normalizeNames(&record.Run.TargetEnvironments); err != nil {
		return fmt.Errorf("risk exception target environments: %w", err)
	}
	if err := validateActor(record.RequestedBy); err != nil {
		return fmt.Errorf("risk exception requester: %w", err)
	}
	if err := validateActor(record.RiskOwner); err != nil {
		return fmt.Errorf("risk exception owner: %w", err)
	}
	if record.RequestedBy.ID == record.RiskOwner.ID || len(record.Reviewers) == 0 || len(record.Reviewers) > 32 {
		return fmt.Errorf("risk exception requires an independent risk owner and reviewers")
	}
	seenActors := map[string]bool{record.RequestedBy.ID: true, record.RiskOwner.ID: true}
	for _, reviewer := range record.Reviewers {
		if err := validateActor(reviewer); err != nil || seenActors[reviewer.ID] {
			return fmt.Errorf("risk exception reviewers must be valid, unique, and independent")
		}
		seenActors[reviewer.ID] = true
	}
	sort.Slice(record.Reviewers, func(left, right int) bool { return record.Reviewers[left].ID < record.Reviewers[right].ID })
	if !cleanText(record.Risk.Title, 200) || !cleanText(record.Risk.Rationale, 4000) ||
		!member([]string{"rare", "unlikely", "possible", "likely", "almost-certain"}, record.Risk.Likelihood) ||
		!member([]string{"low", "medium", "high", "critical"}, record.Risk.Impact) ||
		!cleanText(record.Risk.WorstCredibleOutcome, 4000) {
		return fmt.Errorf("risk exception analysis is incomplete or invalid")
	}
	if len(record.CompensatingControls) == 0 || len(record.CompensatingControls) > 32 {
		return fmt.Errorf("risk exception requires bounded compensating controls")
	}
	for index := range record.CompensatingControls {
		control := &record.CompensatingControls[index]
		if !cleanText(control.Description, 2000) || len(control.EvidenceReferences) == 0 {
			return fmt.Errorf("risk exception compensating control requires description and evidence")
		}
		if err := normalizeUnique(&control.EvidenceReferences, 64, digestPattern); err != nil {
			return fmt.Errorf("risk exception compensating evidence: %w", err)
		}
	}
	if !cleanText(record.Monitoring.Owner, 200) || !cleanText(record.Monitoring.Signal, 2000) ||
		!cleanText(record.Monitoring.Response, 2000) || !cleanText(record.Remediation.Owner, 200) ||
		!cleanText(record.Remediation.Plan, 4000) {
		return fmt.Errorf("risk exception monitoring or remediation is incomplete")
	}
	if !utcTime(record.ApprovedAt) || !utcTime(record.ExpiresAt) || !utcTime(record.Remediation.DueAt) ||
		!record.ApprovedAt.Before(record.ExpiresAt) || record.ExpiresAt.Sub(record.ApprovedAt) > 366*24*time.Hour ||
		record.Remediation.DueAt.Before(record.ApprovedAt) || record.Remediation.DueAt.After(record.ExpiresAt) {
		return fmt.Errorf("risk exception approval, expiry, or remediation dates are invalid")
	}
	record.ApprovedAt, record.ExpiresAt, record.Remediation.DueAt =
		record.ApprovedAt.UTC(), record.ExpiresAt.UTC(), record.Remediation.DueAt.UTC()
	return nil
}

func validateRunBinding(record Record, run model.RunResult) error {
	if run.RunID != record.Run.RunID || run.Inventory.Digest != record.Run.InventoryDigest ||
		run.Plan.ProfileID != record.Run.ProfileID || run.Plan.ProfileVersion != record.Run.ProfileVersion ||
		run.Inventory.TargetName != record.Run.TargetName || run.Inventory.GitCommit != record.Run.TargetCommit ||
		run.Plan.ProjectID != record.Run.ProjectID || !slices.Equal(run.Plan.ArtifactDigests, record.Run.ArtifactDigests) ||
		!slices.Equal(run.Plan.TargetEnvironments, record.Run.TargetEnvironments) {
		return fmt.Errorf("risk exception does not match the immutable run scope")
	}
	if run.CompletedAt.After(record.ApprovedAt) {
		return fmt.Errorf("risk exception was approved before the bound scan completed")
	}
	for _, finding := range run.Findings {
		if finding.ID != record.Finding.FindingID {
			continue
		}
		controls := append([]string(nil), finding.ControlIDs...)
		sort.Strings(controls)
		if finding.Fingerprint != record.Finding.FindingFingerprint || finding.AssertionID != record.Finding.AssertionID ||
			!slices.Equal(controls, record.Finding.ControlIDs) {
			return fmt.Errorf("risk exception finding binding does not match the run")
		}
		for _, result := range run.Results {
			if result.AssertionID == finding.AssertionID && result.Assessment == "fail" {
				return nil
			}
		}
		return fmt.Errorf("risk exception cannot disposition a non-failing assertion")
	}
	return fmt.Errorf("risk exception finding is absent from the bound run")
}

func recordDigest(record Record) (string, error) {
	payload, err := json.Marshal(record)
	if err != nil {
		return "", fmt.Errorf("encode risk exception identity: %w", err)
	}
	digest := sha256.Sum256(payload)
	return hex.EncodeToString(digest[:]), nil
}

func validateActor(actor Actor) error {
	if !actorIDPattern.MatchString(actor.ID) || !cleanText(actor.Name, 200) ||
		!member([]string{"engineering", "security", "privacy", "legal", "operations", "product", "executive"}, actor.Authority) {
		return fmt.Errorf("actor identity, name, or authority is invalid")
	}
	return nil
}

func normalizeDigests(values *[]string) error {
	if len(*values) > 256 {
		return fmt.Errorf("too many values")
	}
	for _, value := range *values {
		if !strings.HasPrefix(value, "sha256:") || !digestPattern.MatchString(strings.TrimPrefix(value, "sha256:")) {
			return fmt.Errorf("invalid digest")
		}
	}
	sort.Strings(*values)
	if adjacentDuplicate(*values) {
		return fmt.Errorf("duplicate value")
	}
	return nil
}

func normalizeNames(values *[]string) error {
	if len(*values) > 256 {
		return fmt.Errorf("too many values")
	}
	for _, value := range *values {
		if !actorIDPattern.MatchString(value) {
			return fmt.Errorf("invalid name")
		}
	}
	sort.Strings(*values)
	if adjacentDuplicate(*values) {
		return fmt.Errorf("duplicate value")
	}
	return nil
}

func normalizeUnique(values *[]string, maximum int, pattern *regexp.Regexp) error {
	if len(*values) == 0 || len(*values) > maximum {
		return fmt.Errorf("requires between 1 and %d values", maximum)
	}
	for _, value := range *values {
		if !pattern.MatchString(value) {
			return fmt.Errorf("invalid value")
		}
	}
	sort.Strings(*values)
	if adjacentDuplicate(*values) {
		return fmt.Errorf("duplicate value")
	}
	return nil
}

func adjacentDuplicate(values []string) bool {
	for index := 1; index < len(values); index++ {
		if values[index] == values[index-1] {
			return true
		}
	}
	return false
}

func cleanText(value string, maximum int) bool {
	return value != "" && value == strings.TrimSpace(value) && len(value) <= maximum && !strings.ContainsRune(value, '\x00')
}

func member(values []string, wanted string) bool {
	return slices.Contains(values, wanted)
}

func utcTime(value time.Time) bool {
	if value.IsZero() {
		return false
	}
	_, offset := value.Zone()
	return offset == 0
}

func rejectAmbiguousYAML(data []byte) error {
	var syntax yaml.Node
	if err := yaml.Unmarshal(data, &syntax); err != nil {
		return fmt.Errorf("decode risk exception syntax: %w", err)
	}
	if len(syntax.Content) != 1 {
		return fmt.Errorf("risk exception requires one YAML document")
	}
	return walkYAML(syntax.Content[0])
}

func walkYAML(node *yaml.Node) error {
	if node.Kind == yaml.AliasNode || node.Tag == "!!null" || node.Tag == "!!merge" {
		return fmt.Errorf("risk exception cannot contain aliases, merge keys, or null values")
	}
	if node.Kind == yaml.MappingNode {
		seen := map[string]bool{}
		for index := 0; index < len(node.Content); index += 2 {
			key := node.Content[index]
			if key.Kind != yaml.ScalarNode || seen[key.Value] {
				return fmt.Errorf("risk exception contains a non-scalar or duplicate mapping key")
			}
			seen[key.Value] = true
		}
	}
	for _, child := range node.Content {
		if err := walkYAML(child); err != nil {
			return err
		}
	}
	return nil
}
