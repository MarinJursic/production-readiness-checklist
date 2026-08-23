package remediation

import (
	"errors"
	"time"

	projectconfig "github.com/MarinJursic/production-readiness-checklist/scanner/config"
	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
	"github.com/MarinJursic/production-readiness-checklist/scanner/provider"
)

// PolicyDeniedError marks a requested remediation that is syntactically valid
// but forbidden by the scanner or project capability policy.
type PolicyDeniedError struct{ Err error }

func (err PolicyDeniedError) Error() string { return err.Err.Error() }
func (err PolicyDeniedError) Unwrap() error { return err.Err }

func policyDenied(err error) error { return PolicyDeniedError{Err: err} }

// IsPolicyDenied reports whether an error represents a denied operation rather
// than invalid configuration or an execution failure.
func IsPolicyDenied(err error) bool {
	var denied PolicyDeniedError
	return errors.As(err, &denied)
}

const (
	FixContractSchema = "prc.fix-contract/v0.2"
	CandidateSchema   = "prc.remediation-candidate/v0.2"
	RunSchema         = "prc.remediation-run/v0.1"
)

type FixContract struct {
	SchemaVersion           string   `json:"schema_version"`
	TaskID                  string   `json:"task_id"`
	BaselineRunID           string   `json:"baseline_run_id"`
	BaselineInventoryDigest string   `json:"baseline_inventory_digest"`
	ConfigurationDigest     string   `json:"configuration_digest,omitempty"`
	ProjectID               string   `json:"project_id,omitempty"`
	AssertionID             string   `json:"assertion_id"`
	ControlIDs              []string `json:"control_ids"`
	Goal                    string   `json:"goal"`
	FixerID                 string   `json:"fixer_id"`
	RemediationClass        string   `json:"remediation_class"`
	Provider                string   `json:"provider,omitempty"`
	ProposalTaskID          string   `json:"proposal_task_id,omitempty"`
	ProposalSHA256          string   `json:"proposal_sha256,omitempty"`
	AllowedPaths            []string `json:"allowed_paths"`
	ProtectedPaths          []string `json:"protected_paths"`
	Network                 string   `json:"network"`
	MaxChangedLines         int      `json:"max_changed_lines"`
	MaxFiles                int      `json:"max_files"`
	Attempt                 int      `json:"attempt"`
	MaxAttempts             int      `json:"max_attempts"`
	Acceptance              []string `json:"acceptance"`
}

// ProjectConfiguration binds one already validated project policy and its exact
// source file into every inventory used during remediation.
type ProjectConfiguration struct {
	Validation projectconfig.Validation
	SourcePath string
}

type Change struct {
	Path         string `json:"path"`
	Kind         string `json:"kind"`
	BeforeSHA256 string `json:"before_sha256"`
	AfterSHA256  string `json:"after_sha256"`
	BeforeMode   uint32 `json:"before_mode"`
	AfterMode    uint32 `json:"after_mode"`
	AddedLines   int    `json:"added_lines"`
	RemovedLines int    `json:"removed_lines"`
}

type Candidate struct {
	SchemaVersion            string      `json:"schema_version"`
	CandidateID              string      `json:"candidate_id"`
	CandidatePath            string      `json:"candidate_path"`
	Contract                 FixContract `json:"contract"`
	CandidateInventoryDigest string      `json:"candidate_inventory_digest"`
	CandidateRunID           string      `json:"candidate_run_id"`
	Changes                  []Change    `json:"changes"`
	BeforeAssessment         string      `json:"before_assessment"`
	AfterAssessment          string      `json:"after_assessment"`
	Accepted                 bool        `json:"accepted"`
	Reasons                  []string    `json:"reasons"`
}

type Options struct {
	CatalogRoot     string
	Target          string
	CandidateDir    string
	ProfileID       string
	TargetName      string
	AssertionID     string
	MaxFiles        int
	MaxChangedLines int
	Attempt         int
	MaxAttempts     int
	Configuration   *ProjectConfiguration
}

type ProposalOptions struct {
	CatalogRoot     string
	Target          string
	CandidateDir    string
	ProfileID       string
	TargetName      string
	Provider        string
	Task            provider.Task
	Output          provider.Output
	MaxFiles        int
	MaxChangedLines int
	Attempt         int
	MaxAttempts     int
	Configuration   *ProjectConfiguration
}

type LoopOptions struct {
	CatalogRoot     string
	Target          string
	CandidateRoot   string
	ProfileID       string
	MaxFiles        int
	MaxChangedLines int
	MaxAttempts     int
	Configuration   *ProjectConfiguration
	Now             func() time.Time
}

type BudgetUsage struct {
	Attempts     int `json:"attempts"`
	ChangedFiles int `json:"changed_files"`
	ChangedLines int `json:"changed_lines"`
}

type RemainingWork struct {
	AssertionID      string   `json:"assertion_id"`
	ControlIDs       []string `json:"control_ids"`
	Title            string   `json:"title"`
	Assessment       string   `json:"assessment"`
	Execution        string   `json:"execution"`
	Severity         string   `json:"severity"`
	Gate             string   `json:"gate"`
	RemediationClass string   `json:"remediation_class"`
	Summary          string   `json:"summary"`
	ReasonCode       string   `json:"reason_code"`
	Reason           string   `json:"reason"`
}

type RemediationRun struct {
	SchemaVersion         string          `json:"schema_version"`
	RunID                 string          `json:"run_id"`
	StartedAt             time.Time       `json:"started_at"`
	CompletedAt           time.Time       `json:"completed_at"`
	ProfileID             string          `json:"profile_id"`
	ConfigurationDigest   string          `json:"configuration_digest,omitempty"`
	ProjectID             string          `json:"project_id,omitempty"`
	SourceInventoryDigest string          `json:"source_inventory_digest"`
	CandidateRoot         string          `json:"candidate_root"`
	ResultWorkspace       string          `json:"result_workspace"`
	FinalInventoryDigest  string          `json:"final_inventory_digest"`
	OriginalUnchanged     bool            `json:"original_unchanged"`
	MaxAttempts           int             `json:"max_attempts"`
	MaxFiles              int             `json:"max_files"`
	MaxChangedLines       int             `json:"max_changed_lines"`
	Usage                 BudgetUsage     `json:"usage"`
	Candidates            []Candidate     `json:"candidates"`
	FinalRun              model.RunResult `json:"final_run"`
	GateState             string          `json:"gate_state"`
	TerminalState         string          `json:"terminal_state"`
	Remaining             []RemainingWork `json:"remaining"`
	StopReasons           []string        `json:"stop_reasons"`
}
