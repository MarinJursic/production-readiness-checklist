// Package verifier runs scanner-owned test commands in an immutable,
// network-denied OCI sandbox and binds the outcome to a candidate inventory.
package verifier

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/MarinJursic/production-readiness-checklist/scanner/inventory"
)

const ExecutionSchema = "prc.verification-execution/v0.1"

var (
	digestPattern = regexp.MustCompile(`^[0-9a-f]{64}$`)
	imagePattern  = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:/-]*@sha256:[0-9a-f]{64}$`)
	errLimit      = errors.New("verification output limit exceeded")
)

// Options are operator-owned verifier capabilities. They are deliberately not
// loaded from the target repository or supplied by a remediation provider.
type Options struct {
	Runtime        string
	Image          string
	Kind           string
	TimeoutSeconds int
	MemoryMB       int
	CPUs           float64
	PIDs           int
	TmpfsMB        int
	MaxStdoutBytes int
	MaxStderrBytes int
}

// Plan is a sealed exact-argv OCI launch plan. The seal is intentionally not
// serializable: callers cannot deserialize a plan and bypass capability checks.
type Plan struct {
	PlanID             string         `json:"plan_id"`
	RuntimePath        string         `json:"runtime_path"`
	RuntimeSHA256      string         `json:"runtime_sha256"`
	Image              string         `json:"image"`
	Kind               string         `json:"kind"`
	Command            []string       `json:"command"`
	Arguments          []string       `json:"arguments"`
	Cleanup            []string       `json:"cleanup_arguments"`
	ContainerName      string         `json:"container_name"`
	Workspace          string         `json:"workspace"`
	Inventory          string         `json:"candidate_inventory_digest"`
	WorkspaceInventory string         `json:"workspace_inventory_digest"`
	Timeout            int            `json:"timeout_seconds"`
	MaxStdout          int            `json:"max_stdout_bytes"`
	MaxStderr          int            `json:"max_stderr_bytes"`
	Policy             PolicyEvidence `json:"policy"`
	seal               [32]byte
}

// StreamEvidence preserves bounded diagnostic provenance without embedding
// potentially sensitive project output in the durable execution record.
type StreamEvidence struct {
	SHA256 string `json:"sha256"`
	Bytes  int    `json:"bytes"`
}

// PolicyEvidence records the scanner-owned sandbox contract without exposing
// host paths or raw container arguments.
type PolicyEvidence struct {
	Network             string  `json:"network"`
	Workspace           string  `json:"workspace"`
	RootFilesystem      string  `json:"root_filesystem"`
	Pull                string  `json:"pull"`
	Capabilities        string  `json:"capabilities"`
	PrivilegeEscalation bool    `json:"privilege_escalation"`
	User                string  `json:"user"`
	TimeoutSeconds      int     `json:"timeout_seconds"`
	MemoryMB            int     `json:"memory_mb"`
	CPUs                float64 `json:"cpus"`
	PIDs                int     `json:"pids"`
	TmpfsMB             int     `json:"tmpfs_mb"`
	MaxStdoutBytes      int     `json:"max_stdout_bytes"`
	MaxStderrBytes      int     `json:"max_stderr_bytes"`
}

// Execution is a content-addressed record for every launched verification,
// including failing tests, timeouts, output limits, and runtime failures.
type Execution struct {
	SchemaVersion            string         `json:"schema_version"`
	ExecutionID              string         `json:"execution_id"`
	PlanID                   string         `json:"plan_id"`
	RuntimeSHA256            string         `json:"runtime_sha256"`
	Kind                     string         `json:"kind"`
	Image                    string         `json:"image"`
	Command                  []string       `json:"command"`
	Policy                   PolicyEvidence `json:"policy"`
	CandidateInventoryDigest string         `json:"candidate_inventory_digest"`
	StartedAt                time.Time      `json:"started_at"`
	CompletedAt              time.Time      `json:"completed_at"`
	DurationMS               int64          `json:"duration_ms"`
	Outcome                  string         `json:"outcome"`
	ReasonCode               string         `json:"reason_code"`
	ExitCode                 int            `json:"exit_code"`
	Stdout                   StreamEvidence `json:"stdout"`
	Stderr                   StreamEvidence `json:"stderr"`
	WorkspaceUnchanged       bool           `json:"workspace_unchanged"`
}

// Defaults returns conservative scanner-owned resource ceilings.
func Defaults(runtimeName, image, kind string) Options {
	return Options{
		Runtime: runtimeName, Image: image, Kind: kind, TimeoutSeconds: 300,
		MemoryMB: 1024, CPUs: 1, PIDs: 128, TmpfsMB: 512,
		MaxStdoutBytes: 1024 * 1024, MaxStderrBytes: 1024 * 1024,
	}
}

// InferKind chooses a supported scanner-owned verifier from the sealed source
// input path. It does not inspect project scripts or agent prose.
func InferKind(sourcePath string) (string, error) {
	switch strings.ToLower(filepath.Ext(sourcePath)) {
	case ".go":
		return "go", nil
	case ".py":
		return "python", nil
	case ".js":
		return "javascript", nil
	default:
		return "", fmt.Errorf("no scanner-owned verifier for source path %q", sourcePath)
	}
}

func commandForKind(kind string) ([]string, error) {
	switch kind {
	case "go":
		return []string{"go", "test", "./..."}, nil
	case "python":
		return []string{"python", "-m", "pytest", "-q"}, nil
	case "javascript":
		return []string{"node", "--test"}, nil
	default:
		return nil, fmt.Errorf("unsupported verifier kind %q", kind)
	}
}

// BuildPlan validates every capability and seals an exact OCI invocation.
func BuildPlan(workspace, candidateInventoryDigest string, options Options) (Plan, error) {
	if err := ValidateOptions(options); err != nil {
		return Plan{}, err
	}
	if !digestPattern.MatchString(candidateInventoryDigest) {
		return Plan{}, fmt.Errorf("verifier requires a candidate inventory digest")
	}
	if runtime.GOOS == "windows" {
		return Plan{}, fmt.Errorf("the current verifier does not support Windows hosts")
	}
	if os.Geteuid() == 0 {
		return Plan{}, fmt.Errorf("verifier refuses to launch containers as host root")
	}
	runtimePath, err := exec.LookPath(options.Runtime)
	if err != nil {
		return Plan{}, fmt.Errorf("find OCI verifier runtime %q: %w", options.Runtime, err)
	}
	runtimePath, err = filepath.Abs(runtimePath)
	if err != nil {
		return Plan{}, fmt.Errorf("resolve OCI verifier runtime: %w", err)
	}
	if resolved, resolveErr := filepath.EvalSymlinks(runtimePath); resolveErr == nil {
		runtimePath = resolved
	}
	runtimeName := strings.TrimSuffix(strings.ToLower(filepath.Base(runtimePath)), ".exe")
	if runtimeName != "docker" && runtimeName != "podman" {
		return Plan{}, fmt.Errorf("unsupported OCI verifier runtime %q", runtimePath)
	}
	runtimeDigest, err := hashRegularFile(runtimePath)
	if err != nil {
		return Plan{}, fmt.Errorf("hash OCI verifier runtime: %w", err)
	}
	workspace, err = filepath.Abs(workspace)
	if err != nil {
		return Plan{}, fmt.Errorf("resolve verification workspace: %w", err)
	}
	info, err := os.Stat(workspace)
	if err != nil || !info.IsDir() {
		return Plan{}, fmt.Errorf("verification workspace is not an accessible directory")
	}
	if resolved, resolveErr := filepath.EvalSymlinks(workspace); resolveErr == nil {
		workspace = resolved
	}
	if strings.Contains(workspace, ",") {
		return Plan{}, fmt.Errorf("verification workspace path cannot contain a comma")
	}
	workspaceInventory, err := inventory.Build(workspace)
	if err != nil {
		return Plan{}, fmt.Errorf("inventory verification workspace: %w", err)
	}
	command, err := commandForKind(options.Kind)
	if err != nil {
		return Plan{}, err
	}
	policy := policyEvidence(options)
	planID, err := verificationPlanID(runtimeDigest, options.Image, options.Kind, command, candidateInventoryDigest, policy)
	if err != nil {
		return Plan{}, err
	}
	name := "prc-verify-" + planID[:24]
	mount := "type=bind,src=" + workspace + ",dst=/workspace,readonly"
	scratchMB := options.TmpfsMB / 2
	executableScratchMB := options.TmpfsMB - scratchMB
	userID := strconv.Itoa(os.Getuid())
	groupID := strconv.Itoa(os.Getgid())
	arguments := []string{
		"run", "--rm", "--interactive", "--name", name, "--pull=never",
		"--network=none", "--read-only", "--cap-drop=ALL",
		"--security-opt=no-new-privileges=true", "--pids-limit=" + strconv.Itoa(options.PIDs),
		"--memory=" + strconv.Itoa(options.MemoryMB) + "m", "--memory-swap=" + strconv.Itoa(options.MemoryMB) + "m",
		"--cpus=" + strconv.FormatFloat(options.CPUs, 'f', -1, 64),
		"--user=" + userID + ":" + groupID,
		"--ulimit=nofile=1024:1024", "--tmpfs=/tmp:rw,noexec,nosuid,nodev,mode=1777,size=" + strconv.Itoa(scratchMB) + "m",
		"--tmpfs=/prc-exec:rw,exec,nosuid,nodev,uid=" + userID + ",gid=" + groupID + ",mode=0700,size=" + strconv.Itoa(executableScratchMB) + "m",
		"--env=HOME=/tmp/home", "--env=TMPDIR=/tmp", "--env=GOTMPDIR=/prc-exec/go-tmp",
		"--env=GOCACHE=/prc-exec/go-build",
		"--env=GOMODCACHE=/prc-exec/go-mod",
		"--env=GOPROXY=off", "--env=GOSUMDB=off",
		"--env=PYTHONDONTWRITEBYTECODE=1", "--env=PYTHONPYCACHEPREFIX=/tmp/pycache",
		"--mount=" + mount, "--workdir=/workspace", options.Image,
	}
	arguments = append(arguments, command...)
	plan := Plan{
		PlanID: planID, RuntimePath: runtimePath, RuntimeSHA256: runtimeDigest,
		Image: options.Image, Kind: options.Kind, Command: command, Arguments: arguments,
		Cleanup: []string{"rm", "--force", name}, ContainerName: name,
		Workspace: workspace, Inventory: candidateInventoryDigest,
		WorkspaceInventory: workspaceInventory.Digest,
		Timeout:            options.TimeoutSeconds, MaxStdout: options.MaxStdoutBytes, MaxStderr: options.MaxStderrBytes,
		Policy: policy,
	}
	plan.seal, err = sealPlan(plan)
	if err != nil {
		return Plan{}, err
	}
	return plan, nil
}

// Preflight proves that the exact runtime can reach an already-present image
// before a remote provider is invoked. It never pulls, builds, or runs an image.
func Preflight(ctx context.Context, options Options) error {
	if err := ValidateOptions(options); err != nil {
		return err
	}
	if ctx == nil {
		ctx = context.Background()
	}
	if runtime.GOOS == "windows" {
		return fmt.Errorf("the current verifier does not support Windows hosts")
	}
	if os.Geteuid() == 0 {
		return fmt.Errorf("verifier refuses to launch containers as host root")
	}
	runtimePath, err := exec.LookPath(options.Runtime)
	if err != nil {
		return fmt.Errorf("find OCI verifier runtime %q: %w", options.Runtime, err)
	}
	runtimePath, err = filepath.Abs(runtimePath)
	if err != nil {
		return fmt.Errorf("resolve OCI verifier runtime: %w", err)
	}
	if resolved, resolveErr := filepath.EvalSymlinks(runtimePath); resolveErr == nil {
		runtimePath = resolved
	}
	runtimeName := strings.TrimSuffix(strings.ToLower(filepath.Base(runtimePath)), ".exe")
	if runtimeName != "docker" && runtimeName != "podman" {
		return fmt.Errorf("unsupported OCI verifier runtime %q", runtimePath)
	}
	beforeDigest, err := hashRegularFile(runtimePath)
	if err != nil {
		return fmt.Errorf("hash OCI verifier runtime: %w", err)
	}
	preflightContext, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	stdout := newBoundedBuffer(1024 * 1024)
	stderr := newBoundedBuffer(1024 * 1024)
	command := exec.CommandContext(preflightContext, runtimePath, "image", "inspect", options.Image)
	command.Env = []string{"PATH=/usr/bin:/bin", "LANG=C.UTF-8", "LC_ALL=C.UTF-8"}
	command.Stdout = stdout
	command.Stderr = stderr
	runErr := command.Run()
	if preflightContext.Err() != nil {
		return fmt.Errorf("OCI verifier image preflight timed out or was cancelled")
	}
	if errors.Is(stdout.Err(), errLimit) || errors.Is(stderr.Err(), errLimit) {
		return fmt.Errorf("OCI verifier image preflight exceeded its output limit")
	}
	if runErr != nil {
		return fmt.Errorf("OCI verifier image is unavailable to the configured runtime")
	}
	afterDigest, err := hashRegularFile(runtimePath)
	if err != nil || afterDigest != beforeDigest {
		return fmt.Errorf("OCI verifier runtime changed during image preflight")
	}
	return nil
}

// ValidateOptions checks operator-owned verifier configuration without
// touching a target workspace or launching a process.
func ValidateOptions(options Options) error {
	if strings.TrimSpace(options.Runtime) == "" {
		return fmt.Errorf("verifier runtime is required")
	}
	if !imagePattern.MatchString(options.Image) || !strings.Contains(strings.SplitN(options.Image, "/", 2)[0], ".") {
		return fmt.Errorf("verifier image must include a registry and immutable SHA-256 digest")
	}
	if _, err := commandForKind(options.Kind); err != nil {
		return err
	}
	if options.TimeoutSeconds < 1 || options.TimeoutSeconds > 3600 || options.MemoryMB < 64 || options.MemoryMB > 32768 ||
		math.IsNaN(options.CPUs) || math.IsInf(options.CPUs, 0) || options.CPUs <= 0 || options.CPUs > 32 ||
		options.PIDs < 2 || options.PIDs > 4096 || options.TmpfsMB < 32 || options.TmpfsMB > 4096 ||
		options.MaxStdoutBytes < 1024 || options.MaxStdoutBytes > 16*1024*1024 ||
		options.MaxStderrBytes < 1024 || options.MaxStderrBytes > 16*1024*1024 {
		return fmt.Errorf("verifier resource limits are outside scanner safety bounds")
	}
	return nil
}

// Run executes a sealed plan and always returns a durable record once a
// process launch is attempted. A non-pass outcome is not a Go error.
func Run(ctx context.Context, plan Plan) (Execution, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	expectedSeal, err := sealPlan(plan)
	if err != nil || expectedSeal != plan.seal {
		return Execution{}, fmt.Errorf("verification plan changed after capability evaluation")
	}
	digest, err := hashRegularFile(plan.RuntimePath)
	if err != nil || digest != plan.RuntimeSHA256 {
		return Execution{}, fmt.Errorf("OCI verifier runtime changed after execution planning")
	}
	before, err := inventory.Build(plan.Workspace)
	if err != nil || before.Digest != plan.WorkspaceInventory {
		return Execution{}, fmt.Errorf("verification workspace does not match the sealed candidate inventory")
	}
	boundedContext, cancel := context.WithTimeout(ctx, time.Duration(plan.Timeout)*time.Second)
	defer cancel()
	stdout := newBoundedBuffer(plan.MaxStdout)
	stderr := newBoundedBuffer(plan.MaxStderr)
	command := exec.CommandContext(boundedContext, plan.RuntimePath, plan.Arguments...)
	command.Env = []string{"PATH=/usr/bin:/bin", "LANG=C.UTF-8", "LC_ALL=C.UTF-8"}
	command.Stdin = bytes.NewReader(nil)
	command.Stdout = stdout
	command.Stderr = stderr
	started := time.Now().UTC()
	runErr := command.Run()
	completed := time.Now().UTC()
	exitCode := 0
	if command.ProcessState != nil {
		exitCode = command.ProcessState.ExitCode()
	} else if runErr != nil {
		exitCode = -1
	}
	outcome, reason := "pass", "tests_passed"
	if boundedContext.Err() != nil {
		cleanup(plan)
		if errors.Is(context.Cause(boundedContext), context.Canceled) || errors.Is(ctx.Err(), context.Canceled) {
			outcome, reason = "cancelled", "cancelled"
		} else {
			outcome, reason = "timeout", "timeout"
		}
	} else if errors.Is(stdout.Err(), errLimit) || errors.Is(stderr.Err(), errLimit) {
		cleanup(plan)
		outcome, reason = "output_limit", "output_limit"
	} else if runErr != nil {
		if exitCode == 125 || exitCode == 126 || exitCode == 127 || exitCode < 0 {
			outcome, reason = "infrastructure_error", "runtime_failed"
		} else {
			outcome, reason = "fail", "tests_failed"
		}
	}
	after, inventoryErr := inventory.Build(plan.Workspace)
	unchanged := inventoryErr == nil && after.Digest == plan.WorkspaceInventory
	if !unchanged {
		outcome, reason = "workspace_changed", "workspace_changed"
	}
	execution := Execution{
		SchemaVersion: ExecutionSchema, PlanID: plan.PlanID, RuntimeSHA256: plan.RuntimeSHA256,
		Kind: plan.Kind, Image: plan.Image, Command: append([]string(nil), plan.Command...), Policy: plan.Policy,
		CandidateInventoryDigest: plan.Inventory,
		StartedAt:                started, CompletedAt: completed, DurationMS: completed.Sub(started).Milliseconds(),
		Outcome: outcome, ReasonCode: reason, ExitCode: exitCode,
		Stdout: streamEvidence(stdout.Bytes()), Stderr: streamEvidence(stderr.Bytes()), WorkspaceUnchanged: unchanged,
	}
	execution.ExecutionID, err = executionID(execution)
	if err != nil {
		return Execution{}, err
	}
	if err := execution.Validate(); err != nil {
		return Execution{}, err
	}
	return execution, nil
}

func (execution Execution) Validate() error {
	if execution.SchemaVersion != ExecutionSchema || !digestPattern.MatchString(execution.ExecutionID) ||
		!digestPattern.MatchString(execution.PlanID) || !digestPattern.MatchString(execution.RuntimeSHA256) ||
		!digestPattern.MatchString(execution.CandidateInventoryDigest) ||
		!imagePattern.MatchString(execution.Image) || execution.StartedAt.IsZero() ||
		execution.CompletedAt.Before(execution.StartedAt) || execution.DurationMS < 0 ||
		execution.CompletedAt.Sub(execution.StartedAt).Milliseconds() != execution.DurationMS ||
		!digestPattern.MatchString(execution.Stdout.SHA256) || !digestPattern.MatchString(execution.Stderr.SHA256) ||
		execution.Stdout.Bytes < 0 || execution.Stderr.Bytes < 0 {
		return fmt.Errorf("verification execution has invalid identity, provenance, or timing")
	}
	expectedCommand, err := commandForKind(execution.Kind)
	if err != nil || !equalStrings(expectedCommand, execution.Command) {
		return fmt.Errorf("verification execution command is not scanner-owned")
	}
	policyOptions := Options{
		Runtime: "docker", Image: execution.Image, Kind: execution.Kind,
		TimeoutSeconds: execution.Policy.TimeoutSeconds, MemoryMB: execution.Policy.MemoryMB,
		CPUs: execution.Policy.CPUs, PIDs: execution.Policy.PIDs, TmpfsMB: execution.Policy.TmpfsMB,
		MaxStdoutBytes: execution.Policy.MaxStdoutBytes, MaxStderrBytes: execution.Policy.MaxStderrBytes,
	}
	if execution.Policy.Network != "deny" || execution.Policy.Workspace != "read_only" ||
		execution.Policy.RootFilesystem != "read_only" || execution.Policy.Pull != "never" ||
		execution.Policy.Capabilities != "none" || execution.Policy.PrivilegeEscalation ||
		execution.Policy.User != "non_root" || ValidateOptions(policyOptions) != nil {
		return fmt.Errorf("verification execution sandbox policy is invalid")
	}
	expectedPlanID, err := verificationPlanID(execution.RuntimeSHA256, execution.Image, execution.Kind, execution.Command,
		execution.CandidateInventoryDigest, execution.Policy)
	if err != nil || expectedPlanID != execution.PlanID {
		return fmt.Errorf("verification execution does not match its sealed plan identity")
	}
	valid := map[string]string{
		"pass": "tests_passed", "fail": "tests_failed", "timeout": "timeout",
		"cancelled": "cancelled", "output_limit": "output_limit",
		"infrastructure_error": "runtime_failed", "workspace_changed": "workspace_changed",
	}
	failedTestExit := execution.ExitCode > 0 && execution.ExitCode != 125 && execution.ExitCode != 126 && execution.ExitCode != 127
	runtimeFailureExit := execution.ExitCode < 0 || execution.ExitCode == 125 || execution.ExitCode == 126 || execution.ExitCode == 127
	if valid[execution.Outcome] != execution.ReasonCode ||
		(execution.Outcome != "workspace_changed" && !execution.WorkspaceUnchanged) ||
		(execution.Outcome == "workspace_changed" && execution.WorkspaceUnchanged) ||
		(execution.Outcome == "pass" && execution.ExitCode != 0) ||
		(execution.Outcome == "fail" && !failedTestExit) ||
		(execution.Outcome == "infrastructure_error" && !runtimeFailureExit) {
		return fmt.Errorf("verification execution has inconsistent outcome")
	}
	expectedID, err := executionID(execution)
	if err != nil || expectedID != execution.ExecutionID {
		return fmt.Errorf("verification execution ID does not match its content")
	}
	return nil
}

func policyEvidence(options Options) PolicyEvidence {
	return PolicyEvidence{
		Network: "deny", Workspace: "read_only", RootFilesystem: "read_only", Pull: "never",
		Capabilities: "none", PrivilegeEscalation: false, User: "non_root",
		TimeoutSeconds: options.TimeoutSeconds, MemoryMB: options.MemoryMB, CPUs: options.CPUs,
		PIDs: options.PIDs, TmpfsMB: options.TmpfsMB,
		MaxStdoutBytes: options.MaxStdoutBytes, MaxStderrBytes: options.MaxStderrBytes,
	}
}

func sealPlan(plan Plan) ([32]byte, error) {
	copy := plan
	copy.seal = [32]byte{}
	payload, err := json.Marshal(copy)
	if err != nil {
		return [32]byte{}, fmt.Errorf("seal verification plan: %w", err)
	}
	return sha256.Sum256(payload), nil
}

func executionID(execution Execution) (string, error) {
	execution.ExecutionID = ""
	payload, err := json.Marshal(execution)
	if err != nil {
		return "", fmt.Errorf("encode verification execution identity: %w", err)
	}
	digest := sha256.Sum256(payload)
	return hex.EncodeToString(digest[:]), nil
}

func verificationPlanID(runtimeSHA256, image, kind string, command []string, inventory string, policy PolicyEvidence) (string, error) {
	identity := struct {
		RuntimeSHA256 string         `json:"runtime_sha256"`
		Image         string         `json:"image"`
		Kind          string         `json:"kind"`
		Command       []string       `json:"command"`
		Inventory     string         `json:"candidate_inventory_digest"`
		Policy        PolicyEvidence `json:"policy"`
	}{runtimeSHA256, image, kind, command, inventory, policy}
	payload, err := json.Marshal(identity)
	if err != nil {
		return "", fmt.Errorf("encode verification plan identity: %w", err)
	}
	digest := sha256.Sum256(payload)
	return hex.EncodeToString(digest[:]), nil
}

func hashRegularFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil || !info.Mode().IsRegular() {
		return "", fmt.Errorf("path is not a regular file")
	}
	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func streamEvidence(data []byte) StreamEvidence {
	digest := sha256.Sum256(data)
	return StreamEvidence{SHA256: hex.EncodeToString(digest[:]), Bytes: len(data)}
}

func cleanup(plan Plan) {
	digest, err := hashRegularFile(plan.RuntimePath)
	if err != nil || digest != plan.RuntimeSHA256 {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	command := exec.CommandContext(ctx, plan.RuntimePath, plan.Cleanup...)
	command.Env = []string{"PATH=/usr/bin:/bin", "LANG=C.UTF-8", "LC_ALL=C.UTF-8"}
	_ = command.Run()
}

type boundedBuffer struct {
	buffer bytes.Buffer
	limit  int
	err    error
}

func newBoundedBuffer(limit int) *boundedBuffer { return &boundedBuffer{limit: limit} }

func (buffer *boundedBuffer) Write(data []byte) (int, error) {
	remaining := buffer.limit - buffer.buffer.Len()
	if remaining <= 0 {
		buffer.err = errLimit
		return 0, errLimit
	}
	if len(data) > remaining {
		_, _ = buffer.buffer.Write(data[:remaining])
		buffer.err = errLimit
		return remaining, errLimit
	}
	return buffer.buffer.Write(data)
}

func (buffer *boundedBuffer) Bytes() []byte { return buffer.buffer.Bytes() }
func (buffer *boundedBuffer) Err() error    { return buffer.err }

func equalStrings(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	for index := range left {
		if left[index] != right[index] {
			return false
		}
	}
	return true
}
