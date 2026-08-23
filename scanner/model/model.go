package model

import "time"

const (
	InventorySchema        = "prc.inventory/v0.3"
	PlanSchema             = "prc.plan/v0.3"
	RunSchema              = "prc.run/v0.3"
	EvidenceSchema         = "prc.evidence/v0.1"
	AdapterExecutionSchema = "prc.adapter-execution/v0.1"
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
	Implementation      string `json:"implementation_id"`
	Applicability       string `json:"applicability"`
	ApplicabilityBy     string `json:"applicability_evaluator"`
	ApplicabilityReason string `json:"applicability_reason"`
}

type Plan struct {
	SchemaVersion       string             `json:"schema_version"`
	Digest              string             `json:"digest"`
	TargetName          string             `json:"target_name"`
	TargetCommit        string             `json:"target_commit,omitempty"`
	InventoryDigest     string             `json:"inventory_digest"`
	ProfileID           string             `json:"profile_id"`
	ProfileVersion      string             `json:"profile_version"`
	ConfigurationDigest string             `json:"configuration_digest,omitempty"`
	ProjectID           string             `json:"project_id,omitempty"`
	ArtifactDigests     []string           `json:"artifact_digests"`
	TargetEnvironments  []string           `json:"target_environments"`
	Assertions          []PlannedAssertion `json:"assertions"`
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

type RunResult struct {
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

type AdapterExecution struct {
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
