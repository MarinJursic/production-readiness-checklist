package provider

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	workspaceinventory "github.com/MarinJursic/production-readiness-checklist/scanner/inventory"
)

var errLimit = errors.New("provider output limit exceeded")

type limitedBuffer struct {
	data  []byte
	limit int
	err   error
}

func (buffer *limitedBuffer) Write(input []byte) (int, error) {
	if buffer.err != nil {
		return 0, buffer.err
	}
	remaining := buffer.limit - len(buffer.data)
	if len(input) > remaining {
		written := max(remaining, 0)
		buffer.data = append(buffer.data, input[:written]...)
		buffer.err = errLimit
		return written, buffer.err
	}
	buffer.data = append(buffer.data, input...)
	return len(input), nil
}

func Run(ctx context.Context, plan Plan, task Task) (Execution, error) {
	if err := task.Validate(); err != nil {
		return Execution{}, err
	}
	invocationStarted := time.Now().UTC()
	failPreflight := func(cause error) (Execution, error) {
		return Execution{}, newRunFailure(
			plan, task, invocationStarted, time.Now().UTC(), "preflight", "preflight_failed",
			"Provider invocation failed scanner preflight before launch.", cause, "", nil, false, "", nil, false,
		)
	}
	workspace, err := existingDirectory(plan.Workspace, "agent workspace")
	if err != nil || workspace != plan.Workspace {
		return failPreflight(fmt.Errorf("agent workspace changed after execution planning"))
	}
	outputDirectory, err := existingDirectory(plan.OutputDirectory, "agent output directory")
	if err != nil || outputDirectory != plan.OutputDirectory {
		return failPreflight(fmt.Errorf("agent output directory changed after execution planning"))
	}
	if err := validatePrivateOutputDirectory(outputDirectory); err != nil {
		return failPreflight(err)
	}
	executionDirectory, err := existingDirectory(plan.ExecutionDirectory, "agent execution directory")
	if err != nil || executionDirectory != plan.ExecutionDirectory || executionDirectory != outputDirectory {
		return failPreflight(fmt.Errorf("agent execution directory changed after execution planning"))
	}
	item, err := workspaceinventory.Build(plan.Workspace)
	if err != nil || item.Digest != task.WorkspaceInventoryDigest {
		return failPreflight(fmt.Errorf("agent workspace changed after execution planning"))
	}
	expectedSeal, err := sealPlan(plan, task)
	if err != nil || expectedSeal != plan.seal {
		return failPreflight(fmt.Errorf("provider execution plan or task changed after capability evaluation"))
	}
	executableDigest, err := hashFile(plan.ExecutablePath, 1024*1024*1024)
	if err != nil || executableDigest != plan.ExecutableSHA256 {
		return failPreflight(fmt.Errorf("provider executable changed after execution planning"))
	}
	_, schemaDigest, _, err := verifiedFile(plan.OutputSchemaPath, "agent output schema", 1024*1024)
	if err != nil || schemaDigest != plan.OutputSchemaSHA256 {
		return failPreflight(fmt.Errorf("agent output schema changed after execution planning"))
	}
	stdoutPath := filepath.Join(plan.OutputDirectory, "provider-stdout.log")
	stderrPath := filepath.Join(plan.OutputDirectory, "provider-stderr.log")
	paths := []string{stdoutPath, stderrPath}
	if plan.ResultPath != "" {
		paths = append(paths, plan.ResultPath)
	}
	for _, path := range paths {
		if _, err := os.Lstat(path); !os.IsNotExist(err) {
			return failPreflight(fmt.Errorf("provider output path already exists"))
		}
	}
	stdout := &limitedBuffer{data: make([]byte, 0, min(task.MaxOutputBytes, 64*1024)), limit: task.MaxOutputBytes}
	stderr := &limitedBuffer{data: make([]byte, 0, min(64*1024, task.MaxOutputBytes)), limit: min(64*1024, task.MaxOutputBytes)}
	runContext, cancel := context.WithTimeout(ctx, time.Duration(task.TimeoutSeconds)*time.Second)
	defer cancel()
	command := exec.CommandContext(runContext, plan.ExecutablePath, plan.Arguments...)
	command.Dir = plan.ExecutionDirectory
	command.Env = filteredEnvironment(plan.EnvironmentVariables)
	command.Stdin = strings.NewReader(plan.prompt)
	command.Stdout = stdout
	command.Stderr = stderr
	started := time.Now().UTC()
	runErr := command.Run()
	completed := time.Now().UTC()
	if writeErr := writeExclusive(stdoutPath, stdout.data); writeErr != nil {
		return Execution{}, newRunFailure(
			plan, task, started, time.Now().UTC(), "transcript", "transcript_failed",
			"Provider transcript persistence failed; no candidate was accepted.", writeErr,
			stdoutPath, stdout.data, false, stderrPath, stderr.data, false,
		)
	}
	if writeErr := writeExclusive(stderrPath, stderr.data); writeErr != nil {
		return Execution{}, newRunFailure(
			plan, task, started, time.Now().UTC(), "transcript", "transcript_failed",
			"Provider transcript persistence failed; no candidate was accepted.", writeErr,
			stdoutPath, stdout.data, true, stderrPath, stderr.data, false,
		)
	}
	if runContext.Err() != nil {
		if ctx.Err() != nil {
			cause := fmt.Errorf("provider execution cancelled: %w", ctx.Err())
			return Execution{}, newRunFailure(
				plan, task, started, completed, "execution", "cancelled",
				"Provider invocation was cancelled; no candidate was accepted.", cause,
				stdoutPath, stdout.data, true, stderrPath, stderr.data, true,
			)
		}
		cause := fmt.Errorf("provider timed out after %d seconds: %w", task.TimeoutSeconds, runContext.Err())
		return Execution{}, newRunFailure(
			plan, task, started, completed, "execution", "timeout",
			"Provider invocation exceeded its scanner-owned timeout; no candidate was accepted.", cause,
			stdoutPath, stdout.data, true, stderrPath, stderr.data, true,
		)
	}
	if errors.Is(stdout.err, errLimit) || errors.Is(stderr.err, errLimit) {
		cause := fmt.Errorf("provider output exceeded its configured byte limit")
		return Execution{}, newRunFailure(
			plan, task, started, completed, "execution", "output_limit",
			"Provider output exceeded its scanner-owned byte limit; no candidate was accepted.", cause,
			stdoutPath, stdout.data, true, stderrPath, stderr.data, true,
		)
	}
	if runErr != nil {
		cause := fmt.Errorf("provider process failed; diagnostics digest %s: %w", digestBytes(stderr.data), runErr)
		return Execution{}, newRunFailure(
			plan, task, started, completed, "execution", "process_failed",
			"Provider process exited unsuccessfully; no candidate was accepted.", cause,
			stdoutPath, stdout.data, true, stderrPath, stderr.data, true,
		)
	}
	item, err = workspaceinventory.Build(plan.Workspace)
	if err != nil || item.Digest != task.WorkspaceInventoryDigest {
		cause := fmt.Errorf("agent workspace changed during provider execution")
		return Execution{}, newRunFailure(
			plan, task, started, completed, "postflight", "workspace_changed",
			"Source workspace integrity changed during provider invocation; no candidate was accepted.", cause,
			stdoutPath, stdout.data, true, stderrPath, stderr.data, true,
		)
	}
	var outputData []byte
	if plan.Provider == "codex" {
		info, err := os.Stat(plan.ResultPath)
		if err != nil || !info.Mode().IsRegular() || info.Size() > int64(task.MaxOutputBytes) {
			cause := fmt.Errorf("codex result is missing or exceeds the output limit")
			return Execution{}, newRunFailure(
				plan, task, started, completed, "postflight", "result_missing",
				"Provider result was unavailable within its scanner-owned bound; no candidate was accepted.", cause,
				stdoutPath, stdout.data, true, stderrPath, stderr.data, true,
			)
		}
		outputData, err = os.ReadFile(plan.ResultPath)
		if err != nil {
			cause := fmt.Errorf("read codex result: %w", err)
			return Execution{}, newRunFailure(
				plan, task, started, completed, "postflight", "result_missing",
				"Provider result was unavailable within its scanner-owned bound; no candidate was accepted.", cause,
				stdoutPath, stdout.data, true, stderrPath, stderr.data, true,
			)
		}
	} else {
		outputData = stdout.data
	}
	output, err := ParseOutput(plan.Provider, outputData, task)
	if err != nil {
		return Execution{}, newRunFailure(
			plan, task, started, completed, "protocol", "protocol_rejected",
			"Provider result failed the sealed output protocol; no candidate was accepted.", err,
			stdoutPath, stdout.data, true, stderrPath, stderr.data, true,
		)
	}
	execution := Execution{
		SchemaVersion: ExecutionSchema, Provider: plan.Provider, TaskID: task.TaskID,
		ExecutableSHA256: plan.ExecutableSHA256, OutputSchemaSHA256: plan.OutputSchemaSHA256,
		StartedAt: started, CompletedAt: completed, DurationMS: completed.Sub(started).Milliseconds(),
		StdoutPath: stdoutPath, StdoutSHA256: digestBytes(stdout.data), StdoutBytes: len(stdout.data),
		StderrPath: stderrPath, StderrSHA256: digestBytes(stderr.data), StderrBytes: len(stderr.data),
		Output: output,
	}
	execution.ExecutionID, err = executionID(execution)
	if err != nil {
		return Execution{}, err
	}
	return execution, nil
}

func filteredEnvironment(names []string) []string {
	allowed := map[string]bool{}
	for _, name := range names {
		allowed[name] = true
	}
	result := make([]string, 0, len(names))
	for _, item := range os.Environ() {
		name, _, found := strings.Cut(item, "=")
		if found && allowed[name] {
			result = append(result, item)
		}
	}
	return result
}

func writeExclusive(path string, data []byte) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		return fmt.Errorf("create provider transcript: %w", err)
	}
	if _, err := bytes.NewReader(data).WriteTo(file); err != nil {
		file.Close()
		return fmt.Errorf("write provider transcript: %w", err)
	}
	if err := file.Sync(); err != nil {
		file.Close()
		return fmt.Errorf("sync provider transcript: %w", err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("close provider transcript: %w", err)
	}
	return nil
}

func digestBytes(data []byte) string {
	digest := sha256.Sum256(data)
	return hex.EncodeToString(digest[:])
}

func executionID(execution Execution) (string, error) {
	execution.ExecutionID = ""
	payload, err := json.Marshal(execution)
	if err != nil {
		return "", fmt.Errorf("encode provider execution identity: %w", err)
	}
	return digestBytes(payload), nil
}
