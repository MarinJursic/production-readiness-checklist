package model

import (
	"encoding/json"
	"time"

	"github.com/MarinJursic/production-readiness-checklist/scanner/controlprogram"
)

const (
	EngineVersion          = "prc.engine/v0.1"
	InventorySchema        = "prc.inventory/v0.3"
	PlanSchema             = "prc.plan/v0.6"
	RunSchema              = "prc.run/v0.12"
	EvidenceSchema         = "prc.evidence/v0.1"
	FindingSchema          = "prc.finding/v0.1"
	AdapterExecutionSchema = "prc.adapter-execution/v0.3"
)

type Source struct {
	Path string `yaml:"path" json:"path"`
	Line int    `yaml:"line" json:"line"`
}

// Control is one stable source control from the complete checklist registry.
// Controls are intentionally broader than executable assertions: a control can
// need several independent repository, environment, and human checks before it
// can be called satisfied.
type Control struct {
	ID             string `json:"id"`
	Status         string `json:"status"`
	Revision       int    `json:"revision"`
	Statement      string `json:"statement"`
	SemanticSHA256 string `json:"semantic_sha256"`
	Source         Source `json:"source"`
}

type ControlCatalogSummary struct {
	SchemaVersion                      string `json:"schema_version"`
	RegistryVersion                    string `json:"registry_version"`
	RegistrySHA256                     string `json:"registry_sha256"`
	SourceSHA256                       string `json:"source_sha256"`
	ContractSchemaVersion              string `json:"contract_schema_version"`
	ContractGeneratorID                string `json:"contract_generator_id,omitempty"`
	ContractSHA256                     string `json:"contract_sha256"`
	ClassificationMethodologySHA256    string `json:"classification_methodology_sha256,omitempty"`
	ClassificationSummarySHA256        string `json:"classification_summary_sha256,omitempty"`
	ClassificationCorpusSHA256         string `json:"classification_corpus_sha256,omitempty"`
	ControlCheckBindingsSchemaVersion  string `json:"control_check_bindings_schema_version,omitempty"`
	ControlCheckBindingsSHA256         string `json:"control_check_bindings_sha256,omitempty"`
	ControlCheckProgramsSchemaVersion  string `json:"control_check_programs_schema_version,omitempty"`
	ControlCheckProgramsSHA256         string `json:"control_check_programs_sha256,omitempty"`
	ControlCheckProgramsCatalogSHA256  string `json:"control_check_programs_catalog_sha256,omitempty"`
	ControlCheckDefinitionSchemaSHA256 string `json:"control_check_definition_schema_sha256,omitempty"`
	ControlCheckDefinitionCorpusSHA256 string `json:"control_check_definition_corpus_sha256,omitempty"`
	ControlCount                       int    `json:"control_count"`
	ActiveControlCount                 int    `json:"active_control_count"`
	ContractCount                      int    `json:"contract_count"`
	GeneratedContractCount             int    `json:"generated_contract_count"`
	AgentReviewedContractCount         int    `json:"agent_reviewed_contract_count"`
	ReviewedDeterministicCount         int    `json:"reviewed_deterministic_count,omitempty"`
	ReviewedNondeterministicCount      int    `json:"reviewed_nondeterministic_count,omitempty"`
	DeterministicBindingCount          int    `json:"deterministic_binding_count,omitempty"`
	DeterministicProgramTemplateCount  int    `json:"deterministic_program_template_count,omitempty"`
	DeterministicProgramBlockedCount   int    `json:"deterministic_program_blocked_count,omitempty"`
	DeterministicProgramExecutedCount  int    `json:"deterministic_program_executed_count,omitempty"`
	DeterministicProgramPassCount      int    `json:"deterministic_program_pass_count,omitempty"`
	DeterministicProgramFailCount      int    `json:"deterministic_program_fail_count,omitempty"`
	DeterministicProgramNACount        int    `json:"deterministic_program_not_applicable_count,omitempty"`
	ProfileTerminalState               string `json:"profile_terminal_state"`
	AIReviewProvider                   string `json:"ai_review_provider,omitempty"`
	AIReviewModel                      string `json:"ai_review_model,omitempty"`
	AIReviewDepth                      string `json:"ai_review_depth,omitempty"`
	AIReviewState                      string `json:"ai_review_state,omitempty"`
	AIReviewedCount                    int    `json:"ai_reviewed_count,omitempty"`
	AIAdvisoryFailCount                int    `json:"ai_advisory_fail_count,omitempty"`
}

// DeterministicClauseResult records one exact program execution without
// copying raw evidence into every control row. The digests bind the result to
// the separately retained program and evidence documents.
type DeterministicClauseResult struct {
	TemplateID                   string    `json:"template_id"`
	CollectorID                  string    `json:"collector_id"`
	ClauseID                     string    `json:"clause_id"`
	ClauseOrdinal                int       `json:"clause_ordinal"`
	ImplementationID             string    `json:"implementation_id"`
	ImplementationContractSHA256 string    `json:"implementation_contract_sha256"`
	RequiredAuthority            string    `json:"required_authority"`
	ProviderID                   string    `json:"provider_id,omitempty"`
	ProgramSHA256                string    `json:"program_sha256,omitempty"`
	EvidenceSHA256               string    `json:"evidence_sha256,omitempty"`
	Status                       string    `json:"status"`
	Outcome                      string    `json:"outcome,omitempty"`
	ReasonCode                   string    `json:"reason_code,omitempty"`
	EvaluatedAt                  time.Time `json:"evaluated_at"`
}

// AIControlReview is advisory evidence produced by an explicitly selected AI
// reviewer. It never has the authority to create a verified pass or a final
// Not Applicable decision.
type AIControlReview struct {
	Provider               string            `json:"provider"`
	Model                  string            `json:"model,omitempty"`
	ReviewDepth            string            `json:"review_depth"`
	AssessmentCandidate    string            `json:"assessment_candidate"`
	ApplicabilityCandidate string            `json:"applicability_candidate"`
	Confidence             string            `json:"confidence"`
	Priority               string            `json:"priority"`
	Reason                 string            `json:"reason"`
	Challenge              string            `json:"challenge"`
	RiskIfIgnored          string            `json:"risk_if_ignored"`
	Advice                 string            `json:"advice"`
	RemediationSteps       []string          `json:"remediation_steps"`
	VerificationSteps      []string          `json:"verification_steps"`
	EvidenceNeeded         []string          `json:"evidence_needed"`
	Evidence               []FindingLocation `json:"evidence"`
	Limitations            []string          `json:"limitations"`
	CitationVerification   string            `json:"citation_verification,omitempty"`
	ClaimVerification      string            `json:"claim_verification,omitempty"`
	TaskID                 string            `json:"task_id"`
}

// ControlResult makes complete-corpus coverage visible without overstating the
// narrower deterministic assertions. "partially_verified" means every linked
// assertion in this run passed, not that the complete broad control passed.
type ControlResult struct {
	ControlID                         string                      `json:"control_id"`
	Revision                          int                         `json:"revision"`
	Statement                         string                      `json:"statement"`
	Source                            Source                      `json:"source"`
	ContractSHA256                    string                      `json:"contract_sha256"`
	ContractStatus                    string                      `json:"contract_status"`
	Classification                    string                      `json:"classification,omitempty"`
	ClassificationRoute               string                      `json:"classification_route,omitempty"`
	ClassificationDecisionBasis       string                      `json:"classification_decision_basis,omitempty"`
	ClassificationRowSHA256           string                      `json:"classification_row_sha256,omitempty"`
	DeterministicBindingID            string                      `json:"deterministic_binding_id,omitempty"`
	DeterministicBindingSHA256        string                      `json:"deterministic_binding_sha256,omitempty"`
	DeterministicProgramTemplateCount int                         `json:"deterministic_program_template_count,omitempty"`
	DeterministicProgramStatus        string                      `json:"deterministic_program_status,omitempty"`
	DeterministicClauseResults        []DeterministicClauseResult `json:"deterministic_clause_results,omitempty"`
	CanonicalControlID                string                      `json:"canonical_control_id"`
	EvaluationClass                   string                      `json:"evaluation_class"`
	AutomationClass                   string                      `json:"automation_class"`
	ApplicabilityClass                string                      `json:"applicability_class"`
	Atomicity                         string                      `json:"atomicity"`
	CompleteInventoryRequired         bool                        `json:"complete_inventory_required"`
	NegativeCondition                 bool                        `json:"negative_condition"`
	ProjectThresholdsRequired         bool                        `json:"project_thresholds_required"`
	EvidenceAuthorities               []string                    `json:"evidence_authorities"`
	NotApplicableProof                string                      `json:"not_applicable_proof"`
	Disposition                       string                      `json:"disposition"`
	Coverage                          string                      `json:"coverage"`
	Authority                         string                      `json:"authority"`
	AssertionIDs                      []string                    `json:"assertion_ids"`
	ExecutedAssertionIDs              []string                    `json:"executed_assertion_ids"`
	Summary                           string                      `json:"summary"`
	AIReview                          *AIControlReview            `json:"ai_review,omitempty"`
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
	Locations        []FindingLocation     `json:"locations,omitempty"`
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
	SchemaVersion         string                    `json:"schema_version"`
	RunID                 string                    `json:"run_id"`
	StartedAt             time.Time                 `json:"started_at"`
	CompletedAt           time.Time                 `json:"completed_at"`
	Plan                  Plan                      `json:"plan"`
	Inventory             Inventory                 `json:"inventory"`
	AdapterExecutions     []AdapterExecution        `json:"adapter_executions"`
	Results               []AssertionResult         `json:"results"`
	Findings              []Finding                 `json:"findings"`
	ControlCatalog        *ControlCatalogSummary    `json:"control_catalog,omitempty"`
	ControlResults        []ControlResult           `json:"control_results,omitempty"`
	DeterministicEvidence []controlprogram.Evidence `json:"deterministic_evidence,omitempty"`
	TerminalState         string                    `json:"terminal_state"`
}

// MarshalJSON preserves the byte contract of archived run records. v0.6 and
// later runs encode the required findings array; v0.5 through v0.3 do not.
func (run RunResult) MarshalJSON() ([]byte, error) {
	type current RunResult
	if run.SchemaVersion == RunSchema {
		return json.Marshal(current(run))
	}
	if run.SchemaVersion == "prc.run/v0.11" {
		type catalogV011 struct {
			SchemaVersion        string `json:"schema_version"`
			RegistryVersion      string `json:"registry_version"`
			RegistrySHA256       string `json:"registry_sha256"`
			SourceSHA256         string `json:"source_sha256"`
			ControlCount         int    `json:"control_count"`
			ActiveControlCount   int    `json:"active_control_count"`
			ProfileTerminalState string `json:"profile_terminal_state"`
			AIReviewProvider     string `json:"ai_review_provider,omitempty"`
			AIReviewModel        string `json:"ai_review_model,omitempty"`
			AIReviewState        string `json:"ai_review_state,omitempty"`
			AIReviewedCount      int    `json:"ai_reviewed_count,omitempty"`
			AIAdvisoryFailCount  int    `json:"ai_advisory_fail_count,omitempty"`
		}
		type controlV011 struct {
			ControlID            string           `json:"control_id"`
			Revision             int              `json:"revision"`
			Statement            string           `json:"statement"`
			Source               Source           `json:"source"`
			Disposition          string           `json:"disposition"`
			Coverage             string           `json:"coverage"`
			Authority            string           `json:"authority"`
			AssertionIDs         []string         `json:"assertion_ids"`
			ExecutedAssertionIDs []string         `json:"executed_assertion_ids"`
			Summary              string           `json:"summary"`
			AIReview             *AIControlReview `json:"ai_review,omitempty"`
		}
		type runV011 struct {
			SchemaVersion     string             `json:"schema_version"`
			RunID             string             `json:"run_id"`
			StartedAt         time.Time          `json:"started_at"`
			CompletedAt       time.Time          `json:"completed_at"`
			Plan              Plan               `json:"plan"`
			Inventory         Inventory          `json:"inventory"`
			AdapterExecutions []AdapterExecution `json:"adapter_executions"`
			Results           []AssertionResult  `json:"results"`
			Findings          []Finding          `json:"findings"`
			ControlCatalog    *catalogV011       `json:"control_catalog,omitempty"`
			ControlResults    []controlV011      `json:"control_results,omitempty"`
			TerminalState     string             `json:"terminal_state"`
		}
		legacy := runV011{
			SchemaVersion: run.SchemaVersion, RunID: run.RunID, StartedAt: run.StartedAt,
			CompletedAt: run.CompletedAt, Plan: run.Plan, Inventory: run.Inventory,
			AdapterExecutions: run.AdapterExecutions, Results: run.Results, Findings: run.Findings,
			TerminalState: run.TerminalState,
		}
		if run.ControlCatalog != nil {
			legacy.ControlCatalog = &catalogV011{
				SchemaVersion: run.ControlCatalog.SchemaVersion, RegistryVersion: run.ControlCatalog.RegistryVersion,
				RegistrySHA256: run.ControlCatalog.RegistrySHA256, SourceSHA256: run.ControlCatalog.SourceSHA256,
				ControlCount: run.ControlCatalog.ControlCount, ActiveControlCount: run.ControlCatalog.ActiveControlCount,
				ProfileTerminalState: run.ControlCatalog.ProfileTerminalState,
				AIReviewProvider:     run.ControlCatalog.AIReviewProvider, AIReviewModel: run.ControlCatalog.AIReviewModel,
				AIReviewState: run.ControlCatalog.AIReviewState, AIReviewedCount: run.ControlCatalog.AIReviewedCount,
				AIAdvisoryFailCount: run.ControlCatalog.AIAdvisoryFailCount,
			}
		}
		for _, control := range run.ControlResults {
			legacy.ControlResults = append(legacy.ControlResults, controlV011{
				ControlID: control.ControlID, Revision: control.Revision, Statement: control.Statement, Source: control.Source,
				Disposition: control.Disposition, Coverage: control.Coverage, Authority: control.Authority,
				AssertionIDs: control.AssertionIDs, ExecutedAssertionIDs: control.ExecutedAssertionIDs,
				Summary: control.Summary, AIReview: control.AIReview,
			})
		}
		return json.Marshal(legacy)
	}
	if run.SchemaVersion == "prc.run/v0.10" || run.SchemaVersion == "prc.run/v0.9" || run.SchemaVersion == "prc.run/v0.8" || run.SchemaVersion == "prc.run/v0.7" || run.SchemaVersion == "prc.run/v0.6" {
		type withFindings struct {
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
		return json.Marshal(withFindings{
			SchemaVersion: run.SchemaVersion, RunID: run.RunID, StartedAt: run.StartedAt,
			CompletedAt: run.CompletedAt, Plan: run.Plan, Inventory: run.Inventory,
			AdapterExecutions: run.AdapterExecutions, Results: run.Results,
			Findings: run.Findings, TerminalState: run.TerminalState,
		})
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

// AdapterDataInput records the content identity of an external read-only data
// dependency without persisting its host filesystem path.
type AdapterDataInput struct {
	Name        string `json:"name"`
	Destination string `json:"destination"`
	SHA256      string `json:"sha256"`
	Files       int    `json:"files"`
	Bytes       int64  `json:"bytes"`
}

type AdapterExecution struct {
	SchemaVersion     string             `json:"schema_version"`
	ExecutionID       string             `json:"execution_id"`
	AdapterRunID      string             `json:"adapter_run_id"`
	AdapterID         string             `json:"adapter_id"`
	ManifestSHA256    string             `json:"manifest_sha256"`
	Image             string             `json:"image"`
	Resolution        AdapterResolution  `json:"resolution"`
	DataInputs        []AdapterDataInput `json:"data_inputs"`
	Subject           AdapterSubject     `json:"subject"`
	StartedAt         time.Time          `json:"started_at"`
	CompletedAt       time.Time          `json:"completed_at"`
	DurationMS        int64              `json:"duration_ms"`
	DiagnosticsSHA256 string             `json:"diagnostics_sha256"`
	DiagnosticsBytes  int                `json:"diagnostics_bytes"`
	Transcript        AdapterTranscript  `json:"transcript"`
}

// MarshalJSON preserves archived v0.1 and v0.2 byte contracts so their
// execution and enclosing run identities remain reproducible.
func (execution AdapterExecution) MarshalJSON() ([]byte, error) {
	type current AdapterExecution
	if execution.SchemaVersion == AdapterExecutionSchema {
		return json.Marshal(current(execution))
	}
	type v02 struct {
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
	if execution.SchemaVersion == "prc.adapter-execution/v0.2" {
		return json.Marshal(v02{
			SchemaVersion: execution.SchemaVersion, ExecutionID: execution.ExecutionID,
			AdapterRunID: execution.AdapterRunID, AdapterID: execution.AdapterID,
			ManifestSHA256: execution.ManifestSHA256, Image: execution.Image,
			Resolution: execution.Resolution, Subject: execution.Subject,
			StartedAt: execution.StartedAt, CompletedAt: execution.CompletedAt, DurationMS: execution.DurationMS,
			DiagnosticsSHA256: execution.DiagnosticsSHA256, DiagnosticsBytes: execution.DiagnosticsBytes,
			Transcript: execution.Transcript,
		})
	}
	type v01 struct {
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
	return json.Marshal(v01{
		SchemaVersion: execution.SchemaVersion, ExecutionID: execution.ExecutionID,
		AdapterRunID: execution.AdapterRunID, AdapterID: execution.AdapterID,
		ManifestSHA256: execution.ManifestSHA256, Image: execution.Image, Subject: execution.Subject,
		StartedAt: execution.StartedAt, CompletedAt: execution.CompletedAt, DurationMS: execution.DurationMS,
		DiagnosticsSHA256: execution.DiagnosticsSHA256, DiagnosticsBytes: execution.DiagnosticsBytes,
		Transcript: execution.Transcript,
	})
}
