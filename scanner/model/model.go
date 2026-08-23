package model

import (
	"encoding/json"
	"time"
)

const (
	EngineVersion          = "prc.engine/v0.1"
	InventorySchema        = "prc.inventory/v0.3"
	PlanSchema             = "prc.plan/v0.6"
	RunSchema              = "prc.run/v0.8"
	EvidenceSchema         = "prc.evidence/v0.1"
	FindingSchema          = "prc.finding/v0.1"
	AdapterExecutionSchema = "prc.adapter-execution/v0.2"
)

type Source struct {
	Path string `yaml:"path" json:"path"`
	Line int    `yaml:"line" json:"line"`
}

type Objective struct {
	ID              string   `yaml:"id" json:"id"`
	Revision        int      `yaml:"revision" json:"revision"`
	Title           string   `yaml:"title" json:"title"`
	Statement       string   `yaml:"statement" json:"statement"`
	Source          Source   `yaml:"source" json:"source"`
	Domains         []string `yaml:"domains" json:"domains"`
	AutomationClass string   `yaml:"automation_class" json:"automation_class"`
	AssertionIDs    []string `yaml:"assertion_ids" json:"assertion_ids"`
}

type EvidenceRequirement struct {
	Kind             string `yaml:"kind" json:"kind"`
	MinimumAuthority string `yaml:"minimum_authority" json:"minimum_authority"`
	Description      string `yaml:"description" json:"description"`
}

type Assertion struct {
	ID               string                `yaml:"id" json:"id"`
	Revision         int                   `yaml:"revision" json:"revision"`
	ControlIDs       []string              `yaml:"control_ids" json:"control_ids"`
	Title            string                `yaml:"title" json:"title"`
	Statement        string                `yaml:"statement" json:"statement"`
	Applicability    string                `yaml:"applicability" json:"applicability"`
	EvidenceRequired []EvidenceRequirement `yaml:"evidence_required" json:"evidence_required"`
	ImplementationID string                `yaml:"implementation_id" json:"implementation_id"`
	Parameters       map[string]any        `yaml:"parameters" json:"parameters,omitempty"`
	Severity         string                `yaml:"severity" json:"severity"`
	Gate             string                `yaml:"gate" json:"gate"`
	RemediationClass string                `yaml:"remediation_class" json:"remediation_class"`
}

type TerminalPolicy struct {
	BlockOn              []string `yaml:"block_on" json:"block_on"`
	AllowManualRemaining bool     `yaml:"allow_manual_remaining" json:"allow_manual_remaining"`
}

type Profile struct {
	SchemaVersion  string         `yaml:"schema_version" json:"schema_version"`
	ID             string         `yaml:"id" json:"id"`
	Version        string         `yaml:"version" json:"version"`
	Title          string         `yaml:"title" json:"title"`
	Description    string         `yaml:"description" json:"description"`
	AssertionIDs   []string       `yaml:"assertion_ids" json:"assertion_ids"`
	TerminalPolicy TerminalPolicy `yaml:"terminal_policy" json:"terminal_policy"`
}

type CIInventory struct {
	GitHubActions bool     `json:"github_actions"`
	WorkflowFiles []string `json:"workflow_files"`
}

type InfrastructureInventory struct {
	TerraformFiles  []string `json:"terraform_files"`
	KubernetesFiles []string `json:"kubernetes_files"`
}

type InventoryComponent struct {
	ID        string `json:"id"`
	Kind      string `json:"kind"`
	Path      string `json:"path"`
	Ecosystem string `json:"ecosystem,omitempty"`
}

type InventoryRelation struct {
	From string `json:"from"`
	To   string `json:"to"`
	Kind string `json:"kind"`
}

type InventoryFact struct {
	Key             string   `json:"key"`
	Value           string   `json:"value"`
	Source          string   `json:"source"`
	Detector        string   `json:"detector"`
	DetectorVersion string   `json:"detector_version"`
	Confidence      float64  `json:"confidence"`
	ScopePath       string   `json:"scope_path"`
	Limitations     []string `json:"limitations"`
}

type FileRecord struct {
	Path   string `json:"path"`
	Size   int64  `json:"size"`
	SHA256 string `json:"sha256"`
	Mode   uint32 `json:"mode"`
}

type DeclaredComponent struct {
	Path string `json:"path"`
	Type string `json:"type"`
}

type DeclaredExclusion struct {
	Path      string `json:"path"`
	Rationale string `json:"rationale"`
}

type DeclaredScope struct {
	ConfigurationDigest string              `json:"configuration_digest"`
	ProjectID           string              `json:"project_id"`
	ProjectName         string              `json:"project_name"`
	RiskProfile         string              `json:"risk_profile"`
	ProfileID           string              `json:"profile_id"`
	SourceRef           string              `json:"source_ref"`
	ArtifactDigests     []string            `json:"artifact_digests"`
	TargetEnvironments  []string            `json:"target_environments"`
	Components          []DeclaredComponent `json:"components"`
	Exclusions          []DeclaredExclusion `json:"exclusions"`
	Features            map[string]bool     `json:"features"`
	DataClassifications []string            `json:"data_classifications"`
	RegulatedData       []string            `json:"regulated_data"`
	ProhibitedEvidence  []string            `json:"prohibited_evidence"`
}

type Inventory struct {
	SchemaVersion     string                  `json:"schema_version"`
	TargetName        string                  `json:"target_name"`
	GitCommit         string                  `json:"git_commit,omitempty"`
	Digest            string                  `json:"digest"`
	FileCount         int                     `json:"file_count"`
	SourceFiles       int                     `json:"source_files"`
	PackageEcosystems []string                `json:"package_ecosystems"`
	Manifests         []string                `json:"manifests"`
	LockFiles         []string                `json:"lock_files"`
	ContainerFiles    []string                `json:"container_files"`
	Symlinks          []string                `json:"symlinks"`
	CI                CIInventory             `json:"ci"`
	Infrastructure    InfrastructureInventory `json:"infrastructure"`
	Components        []InventoryComponent    `json:"components"`
	Relations         []InventoryRelation     `json:"relations"`
	Facts             []InventoryFact         `json:"facts"`
	Files             []FileRecord            `json:"files"`
	DeclaredScope     *DeclaredScope          `json:"declared_scope,omitempty"`
	Root              string                  `json:"-"`
}

type PlannedAssertion struct {
	AssertionID         string `json:"assertion_id"`
	AssertionRevision   int    `json:"assertion_revision,omitempty"`
	DefinitionDigest    string `json:"definition_digest,omitempty"`
	Implementation      string `json:"implementation_id"`
	Applicability       string `json:"applicability"`
	ApplicabilityBy     string `json:"applicability_evaluator"`
	ApplicabilityReason string `json:"applicability_reason"`
}

// CapabilitySet is the closed execution envelope used by a plan policy or
// required by one DAG node. Empty host and secret arrays are significant:
// undeclared external access is denied rather than inherited.
type CapabilitySet struct {
	ReadWorkspace bool     `json:"read_workspace"`
	WriteScratch  bool     `json:"write_scratch"`
	Process       string   `json:"process"`
	Network       string   `json:"network"`
	NetworkHosts  []string `json:"network_hosts"`
	SecretHandles []string `json:"secret_handles"`
}

type PlannedImplementation struct {
	ID           string        `json:"id"`
	Kind         string        `json:"kind"`
	AssertionIDs []string      `json:"assertion_ids"`
	Capabilities CapabilitySet `json:"capabilities"`
	Status       string        `json:"status"`
	Reason       string        `json:"reason,omitempty"`
}

type PlannedAdapter struct {
	AdapterID        string        `json:"adapter_id"`
	ManifestSHA256   string        `json:"manifest_sha256"`
	ObservationKinds []string      `json:"observation_kinds"`
	Capabilities     CapabilitySet `json:"capabilities"`
	Status           string        `json:"status"`
	Reason           string        `json:"reason,omitempty"`
}

type PlanNode struct {
	ID               string        `json:"id"`
	Kind             string        `json:"kind"`
	DependsOn        []string      `json:"depends_on"`
	AssertionID      string        `json:"assertion_id,omitempty"`
	ImplementationID string        `json:"implementation_id,omitempty"`
	AdapterID        string        `json:"adapter_id,omitempty"`
	ManifestSHA256   string        `json:"manifest_sha256,omitempty"`
	Capabilities     CapabilitySet `json:"capabilities"`
	Status           string        `json:"status"`
	Reason           string        `json:"reason,omitempty"`
}

type Plan struct {
	SchemaVersion       string                  `json:"schema_version"`
	Digest              string                  `json:"digest"`
	EngineVersion       string                  `json:"engine_version,omitempty"`
	TargetName          string                  `json:"target_name"`
	TargetCommit        string                  `json:"target_commit,omitempty"`
	InventoryDigest     string                  `json:"inventory_digest"`
	ProfileID           string                  `json:"profile_id"`
	ProfileVersion      string                  `json:"profile_version"`
	ProfileDigest       string                  `json:"profile_digest,omitempty"`
	CatalogDigest       string                  `json:"catalog_digest,omitempty"`
	ConfigurationDigest string                  `json:"configuration_digest,omitempty"`
	ProjectID           string                  `json:"project_id,omitempty"`
	ArtifactDigests     []string                `json:"artifact_digests"`
	TargetEnvironments  []string                `json:"target_environments"`
	ExecutionMode       string                  `json:"execution_mode"`
	CapabilityPolicy    CapabilitySet           `json:"capability_policy"`
	Implementations     []PlannedImplementation `json:"implementations"`
	Adapters            []PlannedAdapter        `json:"adapters"`
	Nodes               []PlanNode              `json:"nodes"`
	Assertions          []PlannedAssertion      `json:"assertions"`
}

// MarshalJSON preserves the byte contract of archived v0.5 through v0.1 plan
// records after the execution DAG was added in v0.6.
func (plan Plan) MarshalJSON() ([]byte, error) {
	type current Plan
	if plan.SchemaVersion == PlanSchema {
		return json.Marshal(current(plan))
	}
	type legacy struct {
		SchemaVersion       string             `json:"schema_version"`
		Digest              string             `json:"digest"`
		EngineVersion       string             `json:"engine_version,omitempty"`
		TargetName          string             `json:"target_name"`
		TargetCommit        string             `json:"target_commit,omitempty"`
		InventoryDigest     string             `json:"inventory_digest"`
		ProfileID           string             `json:"profile_id"`
		ProfileVersion      string             `json:"profile_version"`
		ProfileDigest       string             `json:"profile_digest,omitempty"`
		CatalogDigest       string             `json:"catalog_digest,omitempty"`
		ConfigurationDigest string             `json:"configuration_digest,omitempty"`
		ProjectID           string             `json:"project_id,omitempty"`
		ArtifactDigests     []string           `json:"artifact_digests"`
		TargetEnvironments  []string           `json:"target_environments"`
		Assertions          []PlannedAssertion `json:"assertions"`
	}
	return json.Marshal(legacy{
		SchemaVersion: plan.SchemaVersion, Digest: plan.Digest, EngineVersion: plan.EngineVersion,
		TargetName: plan.TargetName, TargetCommit: plan.TargetCommit, InventoryDigest: plan.InventoryDigest,
		ProfileID: plan.ProfileID, ProfileVersion: plan.ProfileVersion, ProfileDigest: plan.ProfileDigest,
		CatalogDigest: plan.CatalogDigest, ConfigurationDigest: plan.ConfigurationDigest,
		ProjectID: plan.ProjectID, ArtifactDigests: plan.ArtifactDigests,
		TargetEnvironments: plan.TargetEnvironments, Assertions: plan.Assertions,
	})
}

type Evidence struct {
	SchemaVersion string    `json:"schema_version"`
	ID            string    `json:"id"`
	Kind          string    `json:"kind"`
	Authority     string    `json:"authority"`
	Producer      string    `json:"producer"`
	TargetDigest  string    `json:"target_digest"`
	Source        string    `json:"source"`
	ContentSHA256 string    `json:"content_sha256,omitempty"`
	Size          int64     `json:"size,omitempty"`
	ObservedAt    time.Time `json:"observed_at"`
	Summary       string    `json:"summary"`
}

type AssertionResult struct {
	AssertionID      string                `json:"assertion_id"`
	ControlIDs       []string              `json:"control_ids"`
	Applicability    string                `json:"applicability"`
	Execution        string                `json:"execution"`
	Assessment       string                `json:"assessment"`
	Severity         string                `json:"severity"`
	Gate             string                `json:"gate"`
	Summary          string                `json:"summary"`
	EvidenceRequired []EvidenceRequirement `json:"evidence_required"`
	EvidenceObserved []Evidence            `json:"evidence_observed"`
	RemediationClass string                `json:"remediation_class"`
}

type FindingSubject struct {
	Kind            string `json:"kind"`
	ID              string `json:"id"`
	InventoryDigest string `json:"inventory_digest"`
}

type FindingLocation struct {
	Path   string `json:"path"`
	Line   int    `json:"line,omitempty"`
	Column int    `json:"column,omitempty"`
}

// Finding is an actionable, evidence-linked violation. Blocked, manual,
// unknown, and execution-error results remain assertion results and are not
// mislabeled as observed findings.
type Finding struct {
	SchemaVersion    string            `json:"schema_version"`
	ID               string            `json:"id"`
	Fingerprint      string            `json:"fingerprint"`
	AssertionID      string            `json:"assertion_id"`
	ControlIDs       []string          `json:"control_ids"`
	Title            string            `json:"title"`
	Summary          string            `json:"summary"`
	Severity         string            `json:"severity"`
	Gate             string            `json:"gate"`
	RemediationClass string            `json:"remediation_class"`
	Subject          FindingSubject    `json:"subject"`
	Locations        []FindingLocation `json:"locations"`
	EvidenceIDs      []string          `json:"evidence_ids"`
}

type RunResult struct {
	SchemaVersion     string             `json:"schema_version"`
	RunID             string             `json:"run_id"`
	StartedAt         time.Time          `json:"started_at"`
	CompletedAt       time.Time          `json:"completed_at"`
	Plan              Plan               `json:"plan"`
	Inventory         Inventory          `json:"inventory"`
	AdapterExecutions []AdapterExecution `json:"adapter_executions"`
	Results           []AssertionResult  `json:"results"`
	Findings          []Finding          `json:"findings"`
	TerminalState     string             `json:"terminal_state"`
}

// MarshalJSON preserves the byte contract of archived run records. v0.6 and
// later runs encode the required findings array; v0.5 through v0.3 do not.
func (run RunResult) MarshalJSON() ([]byte, error) {
	type current RunResult
	if run.SchemaVersion == RunSchema || run.SchemaVersion == "prc.run/v0.7" || run.SchemaVersion == "prc.run/v0.6" {
		return json.Marshal(current(run))
	}
	type legacy struct {
		SchemaVersion     string             `json:"schema_version"`
		RunID             string             `json:"run_id"`
		StartedAt         time.Time          `json:"started_at"`
		CompletedAt       time.Time          `json:"completed_at"`
		Plan              Plan               `json:"plan"`
		Inventory         Inventory          `json:"inventory"`
		AdapterExecutions []AdapterExecution `json:"adapter_executions"`
		Results           []AssertionResult  `json:"results"`
		TerminalState     string             `json:"terminal_state"`
	}
	return json.Marshal(legacy{
		SchemaVersion: run.SchemaVersion, RunID: run.RunID, StartedAt: run.StartedAt,
		CompletedAt: run.CompletedAt, Plan: run.Plan, Inventory: run.Inventory,
		AdapterExecutions: run.AdapterExecutions, Results: run.Results, TerminalState: run.TerminalState,
	})
}

type AdapterSubject struct {
	TargetName      string `json:"target_name"`
	TargetCommit    string `json:"target_commit,omitempty"`
	InventoryDigest string `json:"inventory_digest"`
}

type AdapterLog struct {
	Level   string `json:"level"`
	Message string `json:"message"`
}

type AdapterLocation struct {
	Path   string `json:"path"`
	Line   int    `json:"line,omitempty"`
	Column int    `json:"column,omitempty"`
}

type AdapterObservation struct {
	ID        string            `json:"id"`
	Kind      string            `json:"kind"`
	Outcome   string            `json:"outcome"`
	Summary   string            `json:"summary"`
	Locations []AdapterLocation `json:"locations"`
	Data      map[string]any    `json:"data,omitempty"`
}

type AdapterArtifact struct {
	ID        string `json:"id"`
	MediaType string `json:"media_type"`
	Digest    string `json:"digest"`
	Size      int64  `json:"size"`
	Path      string `json:"path,omitempty"`
}

type AdapterSummary struct {
	Status string         `json:"status"`
	Counts map[string]int `json:"counts"`
	Reason string         `json:"reason,omitempty"`
}

type AdapterTranscript struct {
	Logs         []AdapterLog         `json:"logs"`
	Observations []AdapterObservation `json:"observations"`
	Artifacts    []AdapterArtifact    `json:"artifacts"`
	Summary      AdapterSummary       `json:"summary"`
}

// AdapterResolution records who authorized the exact manifest. Registry
// fields are present only when the manifest was resolved from a registry trust
// root; explicit local execution records the manifest publisher and the local
// operator grant without implying registry verification.
type AdapterResolution struct {
	Source           string `json:"source"`
	PublisherID      string `json:"publisher_id"`
	Trust            string `json:"trust"`
	RegistryID       string `json:"registry_id,omitempty"`
	RegistryRevision int    `json:"registry_revision,omitempty"`
	RegistryDigest   string `json:"registry_digest,omitempty"`
}

type AdapterExecution struct {
	SchemaVersion     string            `json:"schema_version"`
	ExecutionID       string            `json:"execution_id"`
	AdapterRunID      string            `json:"adapter_run_id"`
	AdapterID         string            `json:"adapter_id"`
	ManifestSHA256    string            `json:"manifest_sha256"`
	Image             string            `json:"image"`
	Resolution        AdapterResolution `json:"resolution"`
	Subject           AdapterSubject    `json:"subject"`
	StartedAt         time.Time         `json:"started_at"`
	CompletedAt       time.Time         `json:"completed_at"`
	DurationMS        int64             `json:"duration_ms"`
	DiagnosticsSHA256 string            `json:"diagnostics_sha256"`
	DiagnosticsBytes  int               `json:"diagnostics_bytes"`
	Transcript        AdapterTranscript `json:"transcript"`
}

// MarshalJSON omits resolution from immutable v0.1 records so their execution
// and enclosing run identities remain byte-for-byte reproducible.
func (execution AdapterExecution) MarshalJSON() ([]byte, error) {
	type current AdapterExecution
	if execution.SchemaVersion != "prc.adapter-execution/v0.1" {
		return json.Marshal(current(execution))
	}
	type legacy struct {
		SchemaVersion     string            `json:"schema_version"`
		ExecutionID       string            `json:"execution_id"`
		AdapterRunID      string            `json:"adapter_run_id"`
		AdapterID         string            `json:"adapter_id"`
		ManifestSHA256    string            `json:"manifest_sha256"`
		Image             string            `json:"image"`
		Subject           AdapterSubject    `json:"subject"`
		StartedAt         time.Time         `json:"started_at"`
		CompletedAt       time.Time         `json:"completed_at"`
		DurationMS        int64             `json:"duration_ms"`
		DiagnosticsSHA256 string            `json:"diagnostics_sha256"`
		DiagnosticsBytes  int               `json:"diagnostics_bytes"`
		Transcript        AdapterTranscript `json:"transcript"`
	}
	return json.Marshal(legacy{
		SchemaVersion: execution.SchemaVersion, ExecutionID: execution.ExecutionID,
		AdapterRunID: execution.AdapterRunID, AdapterID: execution.AdapterID,
		ManifestSHA256: execution.ManifestSHA256, Image: execution.Image, Subject: execution.Subject,
		StartedAt: execution.StartedAt, CompletedAt: execution.CompletedAt, DurationMS: execution.DurationMS,
		DiagnosticsSHA256: execution.DiagnosticsSHA256, DiagnosticsBytes: execution.DiagnosticsBytes,
		Transcript: execution.Transcript,
	})
}
