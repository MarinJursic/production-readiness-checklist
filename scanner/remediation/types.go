package remediation

import (
	"context"
	"errors"
	"time"

	projectconfig "github.com/MarinJursic/production-readiness-checklist/scanner/config"
	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
	"github.com/MarinJursic/production-readiness-checklist/scanner/provider"
	"github.com/MarinJursic/production-readiness-checklist/scanner/verifier"
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

// ProviderExecutionError marks a provider launch, timeout, or protocol failure
// after the scanner has accepted the provider configuration.
type ProviderExecutionError struct{ Err error }

func (err ProviderExecutionError) Error() string { return err.Err.Error() }
func (err ProviderExecutionError) Unwrap() error { return err.Err }

func providerExecution(err error) error { return ProviderExecutionError{Err: err} }

func IsProviderExecution(err error) bool {
	var execution ProviderExecutionError
	return errors.As(err, &execution)
}

const (
	FixContractSchema = "prc.fix-contract/v0.3"
	CandidateSchema   = "prc.remediation-candidate/v0.4"
	RunSchema         = "prc.remediation-run/v0.8"
)

type FixContract struct {
	SchemaVersion           string   `json:"schema_version"`
	TaskID                  string   `json:"task_id"`
	BaselineRunID           string   `json:"baseline_run_id"`
	BaselineInventoryDigest string   `json:"baseline_inventory_digest"`
	FindingID               string   `json:"finding_id"`
	FindingFingerprint      string   `json:"finding_fingerprint"`
	ProposalFindingID       string   `json:"proposal_finding_id,omitempty"`
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
	SchemaVersion            string              `json:"schema_version"`
	CandidateID              string              `json:"candidate_id"`
	CandidatePath            string              `json:"candidate_path"`
	Contract                 FixContract         `json:"contract"`
	CandidateInventoryDigest string              `json:"candidate_inventory_digest"`
	CandidateRunID           string              `json:"candidate_run_id"`
	Changes                  []Change            `json:"changes"`
	BeforeAssessment         string              `json:"before_assessment"`
	AfterAssessment          string              `json:"after_assessment"`
	Accepted                 bool                `json:"accepted"`
	Reasons                  []string            `json:"reasons"`
	Verification             *verifier.Execution `json:"verification,omitempty"`
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
	Verifier        *verifier.Options
	Context         context.Context
}

type LoopOptions struct {
	CatalogRoot        string
	Target             string
	CandidateRoot      string
	ProfileID          string
	MaxFiles           int
	MaxChangedLines    int
	MaxAttempts        int
	MaxDurationSeconds int
	Configuration      *ProjectConfiguration
	Agent              *AgentOptions
	Verifier           *verifier.Options
	Context            context.Context
	Now                func() time.Time
}

// AgentOptions enable the fail-closed suggest-only provider path in the
// bounded remediation loop. Source processing must be acknowledged explicitly.
type AgentOptions struct {
	Provider                    string
	Executable                  string
	OutputSchemaPath            string
	AllowRemoteSourceProcessing bool
	TimeoutSeconds              int
	MaxOutputBytes              int
	MaxCostUSD                  float64
}

type BudgetUsage struct {
	Attempts     int `json:"attempts"`
	ChangedFiles int `json:"changed_files"`
	ChangedLines int `json:"changed_lines"`
}

type RemainingWork struct {
	FindingID          string   `json:"finding_id,omitempty"`
	FindingFingerprint string   `json:"finding_fingerprint,omitempty"`
	AssertionID        string   `json:"assertion_id"`
	ControlIDs         []string `json:"control_ids"`
	Title              string   `json:"title"`
	Assessment         string   `json:"assessment"`
	Execution          string   `json:"execution"`
	Severity           string   `json:"severity"`
	Gate               string   `json:"gate"`
	RemediationClass   string   `json:"remediation_class"`
	Summary            string   `json:"summary"`
	ReasonCode         string   `json:"reason_code"`
	Reason             string   `json:"reason"`
}

// AttemptRecord is the scanner-owned audit linkage for one unit of remediation
// work. It records rejected pre-candidate proposals as well as materialized
// candidates, so an agent attempt cannot disappear from the run history.
type AttemptRecord struct {
	Attempt                 int       `json:"attempt"`
	Mode                    string    `json:"mode"`
	AssertionID             string    `json:"assertion_id"`
	FindingID               string    `json:"finding_id"`
	FindingFingerprint      string    `json:"finding_fingerprint"`
	TaskID                  string    `json:"task_id"`
	StartedAt               time.Time `json:"started_at"`
	CompletedAt             time.Time `json:"completed_at"`
	BeforeInventoryDigest   string    `json:"before_inventory_digest"`
	AfterInventoryDigest    string    `json:"after_inventory_digest,omitempty"`
	ProviderExecutionID     string    `json:"provider_execution_id,omitempty"`
	ProviderFailureID       string    `json:"provider_failure_id,omitempty"`
	VerificationExecutionID string    `json:"verification_execution_id,omitempty"`
	CandidateID             string    `json:"candidate_id,omitempty"`
	Outcome                 string    `json:"outcome"`
	ReasonCode              string    `json:"reason_code"`
	Reason                  string    `json:"reason"`
}

type RemediationRun struct {
	SchemaVersion         string               `json:"schema_version"`
	RunID                 string               `json:"run_id"`
	StartedAt             time.Time            `json:"started_at"`
	CompletedAt           time.Time            `json:"completed_at"`
	ProfileID             string               `json:"profile_id"`
	ConfigurationDigest   string               `json:"configuration_digest,omitempty"`
	ProjectID             string               `json:"project_id,omitempty"`
	SourceInventoryDigest string               `json:"source_inventory_digest"`
	CandidateRoot         string               `json:"candidate_root"`
	ResultWorkspace       string               `json:"result_workspace"`
	FinalInventoryDigest  string               `json:"final_inventory_digest"`
	OriginalUnchanged     bool                 `json:"original_unchanged"`
	MaxAttempts           int                  `json:"max_attempts"`
	MaxFiles              int                  `json:"max_files"`
	MaxChangedLines       int                  `json:"max_changed_lines"`
	MaxDurationSeconds    int                  `json:"max_duration_seconds"`
	Usage                 BudgetUsage          `json:"usage"`
	Attempts              []AttemptRecord      `json:"attempts"`
	Candidates            []Candidate          `json:"candidates"`
	ProviderExecutions    []provider.Execution `json:"provider_executions"`
	ProviderFailures      []provider.Failure   `json:"provider_failures"`
	FinalRun              model.RunResult      `json:"final_run"`
	GateState             string               `json:"gate_state"`
	TerminalState         string               `json:"terminal_state"`
	Remaining             []RemainingWork      `json:"remaining"`
	StopReasons           []string             `json:"stop_reasons"`
}
