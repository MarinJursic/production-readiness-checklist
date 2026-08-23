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
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
)

type OCIPlan struct {
	RuntimePath     string           `json:"runtime_path"`
	RuntimeSHA256   string           `json:"runtime_sha256"`
	Arguments       []string         `json:"arguments"`
	Cleanup         []string         `json:"cleanup_arguments"`
	Name            string           `json:"container_name"`
	Workspace       string           `json:"workspace"`
	WorkspaceSHA256 string           `json:"workspace_sha256,omitempty"`
	DataMounts      []BoundDataMount `json:"data_mounts"`
	seal            [32]byte
}

func BuildOCIPlan(runtimeName, workspace, runID string, manifest Manifest) (OCIPlan, error) {
	return BuildOCIPlanWithData(runtimeName, workspace, runID, manifest, nil)
}

// BuildOCIPlanWithData binds every manifest-declared external data directory
// by content before adding its read-only OCI mount.
func BuildOCIPlanWithData(runtimeName, workspace, runID string, manifest Manifest, dataSources map[string]string) (OCIPlan, error) {
	if err := manifest.ValidateForCurrentEngine(); err != nil {
		return OCIPlan{}, err
	}
	dataMounts, err := bindDataMounts(manifest, dataSources)
	if err != nil {
		return OCIPlan{}, err
	}
	if runtime.GOOS == "windows" {
		return OCIPlan{}, fmt.Errorf("the current OCI adapter runner does not support Windows hosts")
	}
	if os.Geteuid() == 0 {
		return OCIPlan{}, fmt.Errorf("OCI adapter runner refuses to launch containers as host root")
	}
	runtimePath, err := exec.LookPath(runtimeName)
	if err != nil {
		return OCIPlan{}, fmt.Errorf("find OCI runtime %q: %w", runtimeName, err)
	}
	runtimePath, err = filepath.Abs(runtimePath)
	if err != nil {
		return OCIPlan{}, fmt.Errorf("resolve OCI runtime: %w", err)
	}
	detectedRuntime := strings.TrimSuffix(strings.ToLower(filepath.Base(runtimePath)), ".exe")
	if detectedRuntime != "docker" && detectedRuntime != "podman" {
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
	userID := strconv.Itoa(os.Getuid())
	groupID := strconv.Itoa(os.Getgid())
	arguments := []string{
		"run", "--rm", "--interactive", "--name", name, "--pull=never",
		"--network=none", "--read-only", "--cap-drop=ALL",
		"--security-opt=no-new-privileges=true", "--pids-limit=" + strconv.Itoa(pidsLimit),
		"--memory=" + strconv.Itoa(manifest.Resources.MemoryMB) + "m",
		"--memory-swap=" + strconv.Itoa(manifest.Resources.MemoryMB) + "m",
		"--cpus=" + strconv.FormatFloat(manifest.Resources.CPUs, 'f', -1, 64),
		"--user=" + userID + ":" + groupID, "--ulimit=nofile=1024:1024",
	}
	if manifest.Capabilities.WriteScratch {
		arguments = append(arguments, "--tmpfs=/tmp:rw,noexec,nosuid,nodev,mode=1777,size="+strconv.Itoa(manifest.Resources.TmpfsMB)+"m")
	}
	arguments = append(arguments, "--mount="+mount)
	for _, dataMount := range dataMounts {
		arguments = append(arguments, "--mount=type=bind,src="+dataMount.Source+",dst="+dataMount.Destination+",readonly")
	}
	arguments = append(arguments, "--workdir=/workspace", manifest.Image)
	arguments = append(arguments, manifest.Command...)
	plan := OCIPlan{
		RuntimePath: runtimePath, RuntimeSHA256: runtimeDigest, Arguments: arguments,
		Cleanup: []string{"rm", "--force", name}, Name: name, Workspace: workspace,
		DataMounts: dataMounts,
	}
	plan.seal, err = sealOCIPlan(plan, manifest)
	if err != nil {
		return OCIPlan{}, err
	}
	return plan, nil
}

// BuildSnapshotOCIPlan binds a previously verified snapshot digest into the
// sealed plan. RunOCI rechecks the digest immediately before and after the
// container so target drift cannot be mistaken for evidence about the sealed
// inventory.
func BuildSnapshotOCIPlan(runtime string, snapshot *Snapshot, runID string, manifest Manifest) (OCIPlan, error) {
	return BuildSnapshotOCIPlanWithData(runtime, snapshot, runID, manifest, nil)
}

// BuildSnapshotOCIPlanWithData binds both the sealed project snapshot and all
// declared read-only data dependencies into one tamper-evident plan.
func BuildSnapshotOCIPlanWithData(runtime string, snapshot *Snapshot, runID string, manifest Manifest, dataSources map[string]string) (OCIPlan, error) {
	if snapshot == nil || snapshot.Path == "" || !hexDigestPattern.MatchString(snapshot.Digest) {
		return OCIPlan{}, fmt.Errorf("OCI adapter execution requires a valid prepared snapshot")
	}
	actual, err := snapshotDigest(snapshot.Path)
	if err != nil {
		return OCIPlan{}, err
	}
	if actual != snapshot.Digest {
		return OCIPlan{}, fmt.Errorf("adapter snapshot changed before execution planning")
	}
	plan, err := BuildOCIPlanWithData(runtime, snapshot.Path, runID, manifest, dataSources)
	if err != nil {
		return OCIPlan{}, err
	}
	plan.WorkspaceSHA256 = snapshot.Digest
	plan.seal, err = sealOCIPlan(plan, manifest)
	if err != nil {
		return OCIPlan{}, err
	}
	return plan, nil
}

type RunOutput struct {
	Transcript        Transcript `json:"transcript"`
	StartedAt         time.Time  `json:"started_at"`
	CompletedAt       time.Time  `json:"completed_at"`
	DiagnosticsSHA256 string     `json:"diagnostics_sha256"`
	DiagnosticsBytes  int        `json:"diagnostics_bytes"`
	DurationMS        int64      `json:"duration_ms"`
	// ArtifactPayloads carries scanner-owned artifact bytes out of the native
	// normalizer. Payloads are keyed by their sha256: descriptor and are never
	// serialized into an execution or run record.
	ArtifactPayloads map[string][]byte        `json:"-"`
	DataInputs       []model.AdapterDataInput `json:"data_inputs"`
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
	if plan.WorkspaceSHA256 != "" {
		workspaceDigest, digestErr := snapshotDigest(plan.Workspace)
		if digestErr != nil || workspaceDigest != plan.WorkspaceSHA256 {
			return RunOutput{}, fmt.Errorf("adapter snapshot changed after execution planning")
		}
	}
	if err := verifyBoundDataMounts(plan.DataMounts); err != nil {
		return RunOutput{}, err
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
	started := time.Now().UTC()
	err = command.Run()
	completed := time.Now().UTC()
	if plan.WorkspaceSHA256 != "" {
		workspaceDigest, digestErr := snapshotDigest(plan.Workspace)
		if digestErr != nil || workspaceDigest != plan.WorkspaceSHA256 {
			cleanupOCI(plan)
			return outputMetadata(stderr.String(), started, completed), fmt.Errorf("adapter snapshot changed during execution")
		}
	}
	if dataErr := verifyBoundDataMounts(plan.DataMounts); dataErr != nil {
		cleanupOCI(plan)
		return outputMetadata(stderr.String(), started, completed), dataErr
	}
	if ctx.Err() != nil {
		cleanupOCI(plan)
		return outputMetadata(stderr.String(), started, completed), fmt.Errorf("adapter execution deadline reached: %w", ctx.Err())
	}
	if errors.Is(stdout.Err(), errOutputLimit) || errors.Is(stderr.Err(), errOutputLimit) {
		cleanupOCI(plan)
		return outputMetadata(stderr.String(), started, completed), fmt.Errorf("adapter output exceeded its configured byte limit")
	}
	if err != nil {
		return outputMetadata(stderr.String(), started, completed), fmt.Errorf("OCI adapter process failed: %w", err)
	}
	transcript, artifactPayloads, err := ParseManifestOutputWithArtifacts(manifest, bytes.NewReader(stdout.Bytes()))
	if err != nil {
		return outputMetadata(stderr.String(), started, completed), err
	}
	if err := ValidateTranscriptContract(manifest, transcript); err != nil {
		return outputMetadata(stderr.String(), started, completed), err
	}
	output := outputMetadata(stderr.String(), started, completed)
	output.Transcript = transcript
	output.ArtifactPayloads = artifactPayloads
	output.DataInputs = make([]model.AdapterDataInput, 0, len(plan.DataMounts))
	for _, mount := range plan.DataMounts {
		output.DataInputs = append(output.DataInputs, model.AdapterDataInput{
			Name: mount.Name, Destination: mount.Destination, SHA256: mount.SHA256,
			Files: mount.Files, Bytes: mount.Bytes,
		})
	}
	return output, nil
}

func sealOCIPlan(plan OCIPlan, manifest Manifest) ([32]byte, error) {
	payload, err := json.Marshal(struct {
		RuntimePath     string           `json:"runtime_path"`
		RuntimeSHA256   string           `json:"runtime_sha256"`
		Arguments       []string         `json:"arguments"`
		Cleanup         []string         `json:"cleanup_arguments"`
		Name            string           `json:"container_name"`
		Workspace       string           `json:"workspace"`
		WorkspaceSHA256 string           `json:"workspace_sha256,omitempty"`
		DataMounts      []BoundDataMount `json:"data_mounts"`
		Manifest        Manifest         `json:"manifest"`
	}{
		RuntimePath: plan.RuntimePath, RuntimeSHA256: plan.RuntimeSHA256,
		Arguments: plan.Arguments, Cleanup: plan.Cleanup,
		Name: plan.Name, Workspace: plan.Workspace, WorkspaceSHA256: plan.WorkspaceSHA256,
		DataMounts: plan.DataMounts,
		Manifest:   manifest,
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

func outputMetadata(diagnostics string, started, completed time.Time) RunOutput {
	digest := sha256.Sum256([]byte(diagnostics))
	return RunOutput{
		StartedAt: started, CompletedAt: completed,
		DiagnosticsSHA256: fmt.Sprintf("%x", digest), DiagnosticsBytes: len(diagnostics),
		DurationMS: completed.Sub(started).Milliseconds(),
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
func (buffer *boundedBuffer) Bytes() []byte  { return bytes.Clone(buffer.data) }
func (buffer *boundedBuffer) Err() error     { return buffer.err }
