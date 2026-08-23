package controlreview

import (
	"context"
	"time"

	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
)

const (
	TaskSchema   = "prc.control-review-task/v0.1"
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
	CurrentDisposition     string             `json:"current_disposition"`
	CurrentCoverage        string             `json:"current_coverage"`
	CurrentAssertionChecks []AssertionContext `json:"current_assertion_checks"`
}

type ContextFile struct {
	Path      string `json:"path"`
	SHA256    string `json:"sha256"`
	StartLine int    `json:"start_line"`
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
	Provider         string
	Model            string
	ExecutableSHA256 string
	StdoutSHA256     string
	StderrSHA256     string
	Duration         time.Duration
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
}

type Summary struct {
	Provider         string
	Model            string
	ReviewedControls int
	AdvisoryFailures int
	ReusedBatches    int
	CompletedBatches int
	StateDirectory   string
	Focused          bool
}
