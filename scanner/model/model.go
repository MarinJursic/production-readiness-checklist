package model

import "time"

const (
	InventorySchema = "prc.inventory/v0.1"
	PlanSchema      = "prc.plan/v0.1"
	RunSchema       = "prc.run/v0.1"
	EvidenceSchema  = "prc.evidence/v0.1"
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
	GitHubActions bool `json:"github_actions"`
}

type FileRecord struct {
	Path   string `json:"path"`
	Size   int64  `json:"size"`
	SHA256 string `json:"sha256"`
}

type Inventory struct {
	SchemaVersion     string       `json:"schema_version"`
	TargetName        string       `json:"target_name"`
	GitCommit         string       `json:"git_commit,omitempty"`
	Digest            string       `json:"digest"`
	FileCount         int          `json:"file_count"`
	SourceFiles       int          `json:"source_files"`
	PackageEcosystems []string     `json:"package_ecosystems"`
	Manifests         []string     `json:"manifests"`
	LockFiles         []string     `json:"lock_files"`
	CI                CIInventory  `json:"ci"`
	Files             []FileRecord `json:"files"`
	Root              string       `json:"-"`
}

type PlannedAssertion struct {
	AssertionID     string `json:"assertion_id"`
	Implementation  string `json:"implementation_id"`
	Applicability   string `json:"applicability"`
	ApplicabilityBy string `json:"applicability_evaluator"`
}

type Plan struct {
	SchemaVersion   string             `json:"schema_version"`
	Digest          string             `json:"digest"`
	TargetName      string             `json:"target_name"`
	TargetCommit    string             `json:"target_commit,omitempty"`
	InventoryDigest string             `json:"inventory_digest"`
	ProfileID       string             `json:"profile_id"`
	ProfileVersion  string             `json:"profile_version"`
	Assertions      []PlannedAssertion `json:"assertions"`
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
	SchemaVersion string            `json:"schema_version"`
	RunID         string            `json:"run_id"`
	StartedAt     time.Time         `json:"started_at"`
	CompletedAt   time.Time         `json:"completed_at"`
	Plan          Plan              `json:"plan"`
	Inventory     Inventory         `json:"inventory"`
	Results       []AssertionResult `json:"results"`
	TerminalState string            `json:"terminal_state"`
}
