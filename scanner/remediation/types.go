package remediation

import (
	projectconfig "github.com/MarinJursic/production-readiness-checklist/scanner/config"
	"github.com/MarinJursic/production-readiness-checklist/scanner/provider"
)

const (
	FixContractSchema = "prc.fix-contract/v0.2"
	CandidateSchema   = "prc.remediation-candidate/v0.2"
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
	Provider        string
	Task            provider.Task
	Output          provider.Output
	MaxFiles        int
	MaxChangedLines int
	Attempt         int
	MaxAttempts     int
	Configuration   *ProjectConfiguration
}
