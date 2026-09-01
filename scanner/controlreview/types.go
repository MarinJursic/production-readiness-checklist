package controlreview

import (
	"context"
	"time"

	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
)

const (
	TaskSchema   = "prc.control-review-task/v0.5"
	OutputSchema = "prc.control-review-output/v0.3"
)

type AssertionContext struct {
	AssertionID string                  `json:"assertion_id"`
	Assessment  string                  `json:"assessment"`
	Summary     string                  `json:"summary"`
	Locations   []model.FindingLocation `json:"locations"`
}

type TaskControl struct {
	ControlID                   string             `json:"control_id"`
	Statement                   string             `json:"statement"`
	ChecklistSource             model.Source       `json:"checklist_source"`
	ContractSHA256              string             `json:"contract_sha256"`
	ContractStatus              string             `json:"contract_status"`
	Classification              string             `json:"classification"`
	ClassificationRoute         string             `json:"classification_route"`
	ClassificationDecisionBasis string             `json:"classification_decision_basis"`
	ClassificationRowSHA256     string             `json:"classification_row_sha256"`
	CanonicalControlID          string             `json:"canonical_control_id"`
	EvaluationClass             string             `json:"evaluation_class"`
	AutomationClass             string             `json:"automation_class"`
	ApplicabilityClass          string             `json:"applicability_class"`
	Atomicity                   string             `json:"atomicity"`
	CompleteInventory           bool               `json:"complete_inventory_required"`
	NegativeCondition           bool               `json:"negative_condition"`
	ProjectThresholds           bool               `json:"project_thresholds_required"`
	EvidenceAuthorities         []string           `json:"evidence_authorities"`
	NotApplicableProof          string             `json:"not_applicable_proof"`
	CurrentDisposition          string             `json:"current_disposition"`
	CurrentCoverage             string             `json:"current_coverage"`
	CurrentAssertionChecks      []AssertionContext `json:"current_assertion_checks"`
}

type ContextFile struct {
	Path      string `json:"path"`
	SHA256    string `json:"sha256"`
	StartLine int    `json:"start_line"`
	EndLine   int    `json:"end_line"`
	Content   string `json:"content"`
}

type Task struct {
	SchemaVersion             string        `json:"schema_version"`
	TaskID                    string        `json:"task_id"`
	InventoryDigest           string        `json:"inventory_digest"`
	RegistrySHA256            string        `json:"registry_sha256"`
	Provider                  string        `json:"provider"`
	ReviewDepth               string        `json:"review_depth"`
	RequireOneSubagentPerRule bool          `json:"require_one_subagent_per_rule"`
	RequireBatchSkeptic       bool          `json:"require_batch_skeptic"`
	Controls                  []TaskControl `json:"controls"`
	RepositoryPaths           []string      `json:"repository_paths"`
	ContextFiles              []ContextFile `json:"context_files"`
	SnapshotLimitations       []string      `json:"snapshot_limitations"`
}

type Review struct {
	ControlID              string                  `json:"control_id"`
	AssessmentCandidate    string                  `json:"assessment_candidate"`
	ApplicabilityCandidate string                  `json:"applicability_candidate"`
	Confidence             string                  `json:"confidence"`
	Priority               string                  `json:"priority"`
	RootCause              string                  `json:"root_cause"`
	RootCauseKey           string                  `json:"root_cause_key"`
	Effort                 string                  `json:"effort"`
	BlastRadius            string                  `json:"blast_radius"`
	Reason                 string                  `json:"reason"`
	Challenge              string                  `json:"challenge"`
	RiskIfIgnored          string                  `json:"risk_if_ignored"`
	Advice                 string                  `json:"advice"`
	RemediationSteps       []string                `json:"remediation_steps"`
	VerificationSteps      []string                `json:"verification_steps"`
	EvidenceNeeded         []string                `json:"evidence_needed"`
	Evidence               []model.FindingLocation `json:"evidence"`
	Limitations            []string                `json:"limitations"`
}

type Output struct {
	SchemaVersion string   `json:"schema_version"`
	TaskID        string   `json:"task_id"`
	Reviews       []Review `json:"reviews"`
}

type Execution struct {
	Provider           string
	Model              string
	ExecutableSHA256   string
	StdoutSHA256       string
	StderrSHA256       string
	Duration           time.Duration
	TokenUsage         TokenUsage
	TokenUsageKnown    bool
	EstimatedCostUSD   float64
	EstimatedCostKnown bool
}

// TokenUsage is the provider-reported usage for completed, non-cached review
// batches. It is accounting information, not evidence about a control.
type TokenUsage struct {
	InputTokens           int64
	CachedInputTokens     int64
	OutputTokens          int64
	ReasoningOutputTokens int64
}

// Progress reports scanner-owned batch completion. Callers must treat the
// counts as local orchestration state rather than proof of provider subagents.
type Progress struct {
	Phase                string
	Provider             string
	Model                string
	ReviewDepth          string
	StateDirectory       string
	Workers              int
	TotalBatches         int
	CompletedBatches     int
	ReusedBatches        int
	TotalControls        int
	CompletedControls    int
	TokenUsage           TokenUsage
	TokenUsageBatches    int
	EstimatedCostUSD     float64
	EstimatedCostBatches int
	MaxCostUSD           float64
	Elapsed              time.Duration
}

type BatchRunner interface {
	Run(context.Context, Task) (Output, Execution, error)
}

type Options struct {
	Provider                    string
	Executable                  string
	Model                       string
	ReasoningEffort             string
	StateDirectory              string
	SchemaPath                  string
	AllowRemoteSourceProcessing bool
	ReviewDepth                 string
	BatchSize                   int
	Workers                     int
	Timeout                     time.Duration
	MaxCostUSD                  float64
	MaxBatches                  int
	MaxDuration                 time.Duration
	ControlIDs                  []string
	Runner                      BatchRunner
	Progress                    func(Progress)
}

type Summary struct {
	Provider             string
	Model                string
	ReviewDepth          string
	ReviewedControls     int
	AdvisoryFailures     int
	ReusedBatches        int
	CompletedBatches     int
	StateDirectory       string
	Focused              bool
	TokenUsage           TokenUsage
	TokenUsageBatches    int
	EstimatedCostUSD     float64
	EstimatedCostBatches int
	MaxCostUSD           float64
	Duration             time.Duration
}

// Preview is a no-provider-call preflight for an advisory review. It reports
// the bounded source snapshot and exact amount of provider work before the user
// allows a long or costly run to begin.
type Preview struct {
	SchemaVersion   string   `json:"schema_version"`
	Provider        string   `json:"provider"`
	Model           string   `json:"model,omitempty"`
	ReviewDepth     string   `json:"review_depth"`
	Workers         int      `json:"workers"`
	Controls        int      `json:"controls"`
	Batches         int      `json:"batches"`
	BatchSize       int      `json:"batch_size"`
	TimeoutPerBatch string   `json:"timeout_per_batch"`
	MaximumBatches  int      `json:"maximum_batches"`
	MaximumDuration string   `json:"maximum_duration"`
	SourceFiles     int      `json:"source_files"`
	SourceBytes     int64    `json:"source_bytes"`
	OmittedFiles    int      `json:"omitted_files"`
	ContextFiles    int      `json:"context_files_across_batches"`
	ContextBytes    int64    `json:"context_bytes_across_batches"`
	ContextLimited  int      `json:"context_limited_batches"`
	MaxContextFiles int      `json:"maximum_context_files_per_batch"`
	MaxContextBytes int      `json:"maximum_context_bytes_per_batch"`
	Limitations     []string `json:"limitations"`
}
