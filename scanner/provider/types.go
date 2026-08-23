package provider

import "time"

const (
	TaskSchema      = "prc.agent-task/v0.2"
	OutputSchema    = "prc.agent-output/v0.1"
	ExecutionSchema = "prc.agent-execution/v0.1"
	FailureSchema   = "prc.agent-failure/v0.1"
)

type Task struct {
	SchemaVersion               string      `json:"schema_version"`
	TaskID                      string      `json:"task_id"`
	Mode                        string      `json:"mode"`
	WorkspaceInventoryDigest    string      `json:"workspace_inventory_digest"`
	FindingID                   string      `json:"finding_id"`
	FindingFingerprint          string      `json:"finding_fingerprint"`
	AssertionID                 string      `json:"assertion_id"`
	ControlIDs                  []string    `json:"control_ids"`
	Goal                        string      `json:"goal"`
	ReadScope                   string      `json:"read_scope"`
	RelevantPaths               []string    `json:"relevant_paths"`
	Inputs                      []InputFile `json:"inputs"`
	AllowedPaths                []string    `json:"allowed_paths"`
	ProtectedPaths              []string    `json:"protected_paths"`
	AllowedCommands             [][]string  `json:"allowed_commands"`
	Network                     string      `json:"network"`
	Secrets                     string      `json:"secrets"`
	AllowRemoteSourceProcessing bool        `json:"allow_remote_source_processing"`
	TimeoutSeconds              int         `json:"timeout_seconds"`
	MaxOutputBytes              int         `json:"max_output_bytes"`
	MaxCostUSD                  float64     `json:"max_cost_usd"`
}

type InputFile struct {
	Path    string `json:"path"`
	SHA256  string `json:"sha256"`
	Content string `json:"content"`
}

type CommandResult struct {
	Argv   []string `json:"argv"`
	Result string   `json:"result"`
}

type Output struct {
	SchemaVersion              string          `json:"schema_version"`
	TaskID                     string          `json:"task_id"`
	Status                     string          `json:"status"`
	RootCause                  string          `json:"root_cause"`
	ChangedFiles               []string        `json:"changed_files"`
	Patch                      string          `json:"patch"`
	CommandsRequestedOrRun     []CommandResult `json:"commands_requested_or_run"`
	Limitations                []string        `json:"limitations"`
	RequestedCapabilityChanges []string        `json:"requested_capability_changes"`
}

type Capabilities struct {
	Provider                  string `json:"provider"`
	Mode                      string `json:"mode"`
	NonInteractive            bool   `json:"non_interactive"`
	StructuredOutput          bool   `json:"structured_output"`
	ReadScope                 string `json:"read_scope"`
	WorkspaceMutation         bool   `json:"workspace_mutation"`
	Shell                     bool   `json:"shell"`
	NetworkTools              bool   `json:"network_tools"`
	MCP                       bool   `json:"mcp"`
	SessionPersistence        bool   `json:"session_persistence"`
	ScannerTimeout            bool   `json:"scanner_timeout"`
	ProviderCostLimit         bool   `json:"provider_cost_limit"`
	ExactReadPathEnforcement  bool   `json:"exact_read_path_enforcement"`
	ExactWritePathEnforcement bool   `json:"exact_write_path_enforcement"`
	RemoteSourceProcessing    bool   `json:"remote_source_processing"`
}

type Plan struct {
	Provider             string       `json:"provider"`
	ExecutablePath       string       `json:"executable_path"`
	ExecutableSHA256     string       `json:"executable_sha256"`
	Workspace            string       `json:"workspace"`
	ExecutionDirectory   string       `json:"execution_directory"`
	OutputDirectory      string       `json:"output_directory"`
	ResultPath           string       `json:"result_path,omitempty"`
	OutputSchemaPath     string       `json:"output_schema_path"`
	OutputSchemaSHA256   string       `json:"output_schema_sha256"`
	Arguments            []string     `json:"arguments"`
	EnvironmentVariables []string     `json:"environment_variables"`
	Capabilities         Capabilities `json:"capabilities"`
	TaskID               string       `json:"task_id"`
	PromptSHA256         string       `json:"prompt_sha256"`
	prompt               string
	seal                 [32]byte
}

type Execution struct {
	SchemaVersion      string    `json:"schema_version"`
	ExecutionID        string    `json:"execution_id"`
	Provider           string    `json:"provider"`
	TaskID             string    `json:"task_id"`
	ExecutableSHA256   string    `json:"executable_sha256"`
	OutputSchemaSHA256 string    `json:"output_schema_sha256"`
	StartedAt          time.Time `json:"started_at"`
	CompletedAt        time.Time `json:"completed_at"`
	DurationMS         int64     `json:"duration_ms"`
	StdoutPath         string    `json:"stdout_path"`
	StdoutSHA256       string    `json:"stdout_sha256"`
	StdoutBytes        int       `json:"stdout_bytes"`
	StderrPath         string    `json:"stderr_path"`
	StderrSHA256       string    `json:"stderr_sha256"`
	StderrBytes        int       `json:"stderr_bytes"`
	Output             Output    `json:"output"`
}

// Failure records a provider invocation that did not produce a valid Output.
// Reason is scanner-authored and deliberately excludes raw provider content;
// bounded transcript digests preserve the underlying evidence separately.
type Failure struct {
	SchemaVersion       string    `json:"schema_version"`
	FailureID           string    `json:"failure_id"`
	Provider            string    `json:"provider"`
	TaskID              string    `json:"task_id"`
	ExecutableSHA256    string    `json:"executable_sha256"`
	OutputSchemaSHA256  string    `json:"output_schema_sha256"`
	StartedAt           time.Time `json:"started_at"`
	CompletedAt         time.Time `json:"completed_at"`
	DurationMS          int64     `json:"duration_ms"`
	Stage               string    `json:"stage"`
	ReasonCode          string    `json:"reason_code"`
	Reason              string    `json:"reason"`
	TranscriptsComplete bool      `json:"transcripts_complete"`
	StdoutPath          string    `json:"stdout_path,omitempty"`
	StdoutSHA256        string    `json:"stdout_sha256,omitempty"`
	StdoutBytes         int       `json:"stdout_bytes"`
	StderrPath          string    `json:"stderr_path,omitempty"`
	StderrSHA256        string    `json:"stderr_sha256,omitempty"`
	StderrBytes         int       `json:"stderr_bytes"`
}
