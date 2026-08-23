package adapter

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type OCIPlan struct {
	RuntimePath   string   `json:"runtime_path"`
	RuntimeSHA256 string   `json:"runtime_sha256"`
	Arguments     []string `json:"arguments"`
	Cleanup       []string `json:"cleanup_arguments"`
	Name          string   `json:"container_name"`
	seal          [32]byte
}

func BuildOCIPlan(runtime, workspace, runID string, manifest Manifest) (OCIPlan, error) {
	if err := manifest.Validate(); err != nil {
		return OCIPlan{}, err
	}
	runtimePath, err := exec.LookPath(runtime)
	if err != nil {
		return OCIPlan{}, fmt.Errorf("find OCI runtime %q: %w", runtime, err)
	}
	runtimePath, err = filepath.Abs(runtimePath)
	if err != nil {
		return OCIPlan{}, fmt.Errorf("resolve OCI runtime: %w", err)
	}
	runtimeName := strings.TrimSuffix(strings.ToLower(filepath.Base(runtimePath)), ".exe")
	if runtimeName != "docker" && runtimeName != "podman" {
		return OCIPlan{}, fmt.Errorf("unsupported OCI runtime %q", runtimePath)
	}
	if resolved, resolveErr := filepath.EvalSymlinks(runtimePath); resolveErr == nil {
		runtimePath = resolved
	}
	runtimeDigest, err := hashExecutable(runtimePath)
	if err != nil {
		return OCIPlan{}, err
	}
	workspace, err = filepath.Abs(workspace)
	if err != nil {
		return OCIPlan{}, fmt.Errorf("resolve adapter workspace: %w", err)
	}
	info, err := os.Stat(workspace)
	if err != nil || !info.IsDir() {
		return OCIPlan{}, fmt.Errorf("adapter workspace is not an accessible directory")
	}
	if resolved, resolveErr := filepath.EvalSymlinks(workspace); resolveErr == nil {
		workspace = resolved
	}
	if strings.Contains(workspace, ",") {
		return OCIPlan{}, fmt.Errorf("adapter workspace path cannot contain a comma")
	}
	if len(runID) < 12 {
		return OCIPlan{}, fmt.Errorf("run ID must contain at least 12 characters")
	}
	if !hexDigestPattern.MatchString(runID) {
		return OCIPlan{}, fmt.Errorf("run ID must be a lowercase SHA-256 digest")
	}
	name := "prc-adapter-" + sanitizeContainerName(runID[:24])
	mount := "type=bind,src=" + workspace + ",dst=/workspace,readonly"
	pidsLimit := manifest.Resources.PIDs
	if !manifest.Capabilities.ChildProcesses {
		pidsLimit = 1
	}
	arguments := []string{
		"run", "--rm", "--interactive", "--name", name, "--pull=never",
		"--network=none", "--read-only", "--cap-drop=ALL",
		"--security-opt=no-new-privileges", "--pids-limit=" + strconv.Itoa(pidsLimit),
		"--memory=" + strconv.Itoa(manifest.Resources.MemoryMB) + "m",
		"--cpus=" + strconv.FormatFloat(manifest.Resources.CPUs, 'f', -1, 64),
		"--user=65532:65532",
	}
	if manifest.Capabilities.WriteScratch {
		arguments = append(arguments, "--tmpfs=/tmp:rw,noexec,nosuid,size="+strconv.Itoa(manifest.Resources.TmpfsMB)+"m")
	}
	arguments = append(arguments, "--mount="+mount, "--workdir=/workspace", manifest.Image)
	arguments = append(arguments, manifest.Command...)
	plan := OCIPlan{
		RuntimePath: runtimePath, RuntimeSHA256: runtimeDigest, Arguments: arguments,
		Cleanup: []string{"rm", "--force", name}, Name: name,
	}
	plan.seal, err = sealOCIPlan(plan, manifest)
	if err != nil {
		return OCIPlan{}, err
	}
	return plan, nil
}

type RunOutput struct {
	Transcript        Transcript `json:"transcript"`
	DiagnosticsSHA256 string     `json:"diagnostics_sha256"`
	DiagnosticsBytes  int        `json:"diagnostics_bytes"`
	DurationMS        int64      `json:"duration_ms"`
}

func RunOCI(
	ctx context.Context,
	plan OCIPlan,
	manifest Manifest,
	input []byte,
) (RunOutput, error) {
	expectedSeal, err := sealOCIPlan(plan, manifest)
	if err != nil {
		return RunOutput{}, err
	}
	if plan.seal != expectedSeal {
		return RunOutput{}, fmt.Errorf("OCI execution plan or adapter manifest changed after capability evaluation")
	}
	runtimeDigest, err := hashExecutable(plan.RuntimePath)
	if err != nil {
		return RunOutput{}, err
	}
	if runtimeDigest != plan.RuntimeSHA256 {
		return RunOutput{}, fmt.Errorf("OCI runtime changed after execution planning")
	}
	if len(input) > manifest.Resources.MaxStdin {
		return RunOutput{}, fmt.Errorf("adapter input exceeds %d bytes", manifest.Resources.MaxStdin)
	}
	ctx, cancel := context.WithTimeout(ctx, manifest.Timeout())
	defer cancel()
	stdout := newBoundedBuffer(manifest.Resources.MaxStdout)
	stderr := newBoundedBuffer(manifest.Resources.MaxStderr)
	command := exec.CommandContext(ctx, plan.RuntimePath, plan.Arguments...)
	command.Env = []string{"PATH=/usr/bin:/bin", "LANG=C.UTF-8", "LC_ALL=C.UTF-8"}
	command.Stdin = bytes.NewReader(input)
	command.Stdout = stdout
	command.Stderr = stderr
	started := time.Now()
	err = command.Run()
	duration := time.Since(started)
	if ctx.Err() != nil {
		cleanupOCI(plan)
		return outputMetadata(stderr.String(), duration), fmt.Errorf("adapter timed out after %s: %w", manifest.Timeout(), ctx.Err())
	}
	if errors.Is(stdout.Err(), errOutputLimit) || errors.Is(stderr.Err(), errOutputLimit) {
		cleanupOCI(plan)
		return outputMetadata(stderr.String(), duration), fmt.Errorf("adapter output exceeded its configured byte limit")
	}
	if err != nil {
		return outputMetadata(stderr.String(), duration), fmt.Errorf("OCI adapter process failed: %w", err)
	}
	transcript, err := ParseOutput(strings.NewReader(stdout.String()), manifest.Resources.Limits)
	if err != nil {
		return outputMetadata(stderr.String(), duration), err
	}
	output := outputMetadata(stderr.String(), duration)
	output.Transcript = transcript
	return output, nil
}

func sealOCIPlan(plan OCIPlan, manifest Manifest) ([32]byte, error) {
	payload, err := json.Marshal(struct {
		RuntimePath   string   `json:"runtime_path"`
		RuntimeSHA256 string   `json:"runtime_sha256"`
		Arguments     []string `json:"arguments"`
		Cleanup       []string `json:"cleanup_arguments"`
		Name          string   `json:"container_name"`
		Manifest      Manifest `json:"manifest"`
	}{
		RuntimePath: plan.RuntimePath, RuntimeSHA256: plan.RuntimeSHA256,
		Arguments: plan.Arguments, Cleanup: plan.Cleanup,
		Name: plan.Name, Manifest: manifest,
	})
	if err != nil {
		return [32]byte{}, fmt.Errorf("seal validated OCI plan: %w", err)
	}
	return sha256.Sum256(payload), nil
}

func hashExecutable(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("open OCI runtime: %w", err)
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil || !info.Mode().IsRegular() {
		return "", fmt.Errorf("OCI runtime is not a regular file")
	}
	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", fmt.Errorf("hash OCI runtime: %w", err)
	}
	return fmt.Sprintf("%x", hasher.Sum(nil)), nil
}

func outputMetadata(diagnostics string, duration time.Duration) RunOutput {
	digest := sha256.Sum256([]byte(diagnostics))
	return RunOutput{
		DiagnosticsSHA256: fmt.Sprintf("%x", digest), DiagnosticsBytes: len(diagnostics), DurationMS: duration.Milliseconds(),
	}
}

func cleanupOCI(plan OCIPlan) {
	digest, err := hashExecutable(plan.RuntimePath)
	if err != nil || digest != plan.RuntimeSHA256 {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	command := exec.CommandContext(ctx, plan.RuntimePath, plan.Cleanup...)
	command.Env = []string{"PATH=/usr/bin:/bin", "LANG=C.UTF-8", "LC_ALL=C.UTF-8"}
	command.Stdout = nil
	command.Stderr = nil
	_ = command.Run()
}

func sanitizeContainerName(value string) string {
	var result strings.Builder
	for _, character := range strings.ToLower(value) {
		if (character >= 'a' && character <= 'z') || (character >= '0' && character <= '9') || character == '-' || character == '_' {
			result.WriteRune(character)
		} else {
			result.WriteByte('-')
		}
	}
	return strings.Trim(result.String(), "-_")
}

var errOutputLimit = errors.New("output limit exceeded")

type boundedBuffer struct {
	data  []byte
	limit int
	err   error
}

func newBoundedBuffer(limit int) *boundedBuffer {
	return &boundedBuffer{data: make([]byte, 0, min(limit, 64*1024)), limit: limit}
}

func (buffer *boundedBuffer) Write(input []byte) (int, error) {
	if buffer.err != nil {
		return 0, buffer.err
	}
	remaining := buffer.limit - len(buffer.data)
	if len(input) > remaining {
		written := max(remaining, 0)
		buffer.data = append(buffer.data, input[:written]...)
		buffer.err = errOutputLimit
		return written, buffer.err
	}
	buffer.data = append(buffer.data, input...)
	return len(input), nil
}

func (buffer *boundedBuffer) String() string { return string(buffer.data) }
func (buffer *boundedBuffer) Err() error     { return buffer.err }
