package controlreview

import (
	"context"
	"time"

	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
)

const (
	TaskSchema   = "prc.control-review-task/v0.2"
	OutputSchema = "prc.control-review-output/v0.1"
)

type AssertionContext struct {
	AssertionID string                  `json:"assertion_id"`
	Assessment  string                  `json:"assessment"`
	Summary     string                  `json:"summary"`
	Locations   []model.FindingLocation `json:"locations"`
}

type TaskControl struct {
	ControlID              string             `json:"control_id"`
	Statement              string             `json:"statement"`
	ChecklistSource        model.Source       `json:"checklist_source"`
	ContractSHA256         string             `json:"contract_sha256"`
	ContractStatus         string             `json:"contract_status"`
	CanonicalControlID     string             `json:"canonical_control_id"`
	EvaluationClass        string             `json:"evaluation_class"`
	AutomationClass        string             `json:"automation_class"`
	ApplicabilityClass     string             `json:"applicability_class"`
	Atomicity              string             `json:"atomicity"`
	CompleteInventory      bool               `json:"complete_inventory_required"`
	NegativeCondition      bool               `json:"negative_condition"`
	ProjectThresholds      bool               `json:"project_thresholds_required"`
	EvidenceAuthorities    []string           `json:"evidence_authorities"`
	NotApplicableProof     string             `json:"not_applicable_proof"`
	CurrentDisposition     string             `json:"current_disposition"`
	CurrentCoverage        string             `json:"current_coverage"`
	CurrentAssertionChecks []AssertionContext `json:"current_assertion_checks"`
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
	RequireOneSubagentPerRule bool          `json:"require_one_subagent_per_rule"`
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
	Reason                 string                  `json:"reason"`
	Advice                 string                  `json:"advice"`
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
	BatchSize                   int
	Workers                     int
	Timeout                     time.Duration
	MaxCostUSD                  float64
	ControlIDs                  []string
	Runner                      BatchRunner
	Progress                    func(Progress)
}

type Summary struct {
	Provider             string
	Model                string
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
