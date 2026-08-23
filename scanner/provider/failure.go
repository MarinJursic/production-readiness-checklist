package provider

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

var failureReasons = map[string]string{
	"preflight_failed":  "Provider invocation failed scanner preflight before launch.",
	"transcript_failed": "Provider transcript persistence failed; no candidate was accepted.",
	"cancelled":         "Provider invocation was cancelled; no candidate was accepted.",
	"timeout":           "Provider invocation exceeded its scanner-owned timeout; no candidate was accepted.",
	"output_limit":      "Provider output exceeded its scanner-owned byte limit; no candidate was accepted.",
	"process_failed":    "Provider process exited unsuccessfully; no candidate was accepted.",
	"workspace_changed": "Source workspace integrity changed during provider invocation; no candidate was accepted.",
	"result_missing":    "Provider result was unavailable within its scanner-owned bound; no candidate was accepted.",
	"protocol_rejected": "Provider result failed the sealed output protocol; no candidate was accepted.",
}

// RunFailureError preserves a content-addressed safe failure record while
// retaining the original cause for errors.Is/errors.As and CLI classification.
type RunFailureError struct {
	Failure Failure
	Err     error
}

func (failure RunFailureError) Error() string { return failure.Err.Error() }
func (failure RunFailureError) Unwrap() error { return failure.Err }

// FailureFromError extracts a scanner-authored provider failure record.
func FailureFromError(err error) (Failure, bool) {
	var failure RunFailureError
	if !errors.As(err, &failure) {
		return Failure{}, false
	}
	return failure.Failure, true
}

func (failure Failure) Validate() error {
	if failure.SchemaVersion != FailureSchema || failure.FailureID == "" ||
		failure.Provider != "codex" && failure.Provider != "claude" ||
		!hexDigest.MatchString(failure.TaskID) || !hexDigest.MatchString(failure.ExecutableSHA256) ||
		!hexDigest.MatchString(failure.OutputSchemaSHA256) {
		return fmt.Errorf("provider failure identity or binding is invalid")
	}
	if failure.StartedAt.IsZero() || failure.CompletedAt.IsZero() ||
		failure.CompletedAt.Before(failure.StartedAt) || failure.DurationMS < 0 ||
		failure.Reason == "" || len(failure.Reason) > 4096 {
		return fmt.Errorf("provider failure time or reason is invalid")
	}
	valid := map[string]map[string]bool{
		"preflight":  {"preflight_failed": true},
		"transcript": {"transcript_failed": true},
		"execution": {
			"cancelled": true, "timeout": true, "output_limit": true, "process_failed": true,
		},
		"postflight": {"workspace_changed": true, "result_missing": true},
		"protocol":   {"protocol_rejected": true},
	}
	if !valid[failure.Stage][failure.ReasonCode] {
		return fmt.Errorf("provider failure stage and reason code are invalid")
	}
	if failure.Reason != failureReasons[failure.ReasonCode] {
		return fmt.Errorf("provider failure reason is not the canonical safe text")
	}
	if failure.StdoutBytes < 0 || failure.StderrBytes < 0 ||
		!validTranscript(failure.StdoutPath, failure.StdoutSHA256, failure.StdoutBytes) ||
		!validTranscript(failure.StderrPath, failure.StderrSHA256, failure.StderrBytes) ||
		failure.TranscriptsComplete != (failure.StdoutPath != "" && failure.StderrPath != "") {
		return fmt.Errorf("provider failure transcript metadata is invalid")
	}
	want, err := failureID(failure)
	if err != nil || failure.FailureID != want || !hexDigest.MatchString(failure.FailureID) {
		return fmt.Errorf("provider failure ID does not match its canonical content")
	}
	return nil
}

func validTranscript(path, digest string, bytes int) bool {
	return path == "" && digest == "" && bytes == 0 || path != "" && hexDigest.MatchString(digest)
}

func failureID(failure Failure) (string, error) {
	failure.FailureID = ""
	payload, err := json.Marshal(failure)
	if err != nil {
		return "", fmt.Errorf("encode provider failure identity: %w", err)
	}
	return digestBytes(payload), nil
}

func newRunFailure(
	plan Plan,
	task Task,
	started, completed time.Time,
	stage, code, reason string,
	cause error,
	stdoutPath string,
	stdout []byte,
	stdoutWritten bool,
	stderrPath string,
	stderr []byte,
	stderrWritten bool,
) error {
	failure := Failure{
		SchemaVersion: FailureSchema, Provider: plan.Provider, TaskID: task.TaskID,
		ExecutableSHA256: plan.ExecutableSHA256, OutputSchemaSHA256: plan.OutputSchemaSHA256,
		StartedAt: started, CompletedAt: completed, DurationMS: completed.Sub(started).Milliseconds(),
		Stage: stage, ReasonCode: code, Reason: reason,
		TranscriptsComplete: stdoutWritten && stderrWritten,
		StdoutBytes:         len(stdout), StderrBytes: len(stderr),
	}
	if stdoutWritten {
		failure.StdoutPath, failure.StdoutSHA256 = stdoutPath, digestBytes(stdout)
	}
	if stderrWritten {
		failure.StderrPath, failure.StderrSHA256 = stderrPath, digestBytes(stderr)
	}
	identifier, err := failureID(failure)
	if err != nil {
		return fmt.Errorf("%w; record provider failure: %v", cause, err)
	}
	failure.FailureID = identifier
	if err := failure.Validate(); err != nil {
		return fmt.Errorf("%w; validate provider failure: %v", cause, err)
	}
	return RunFailureError{Failure: failure, Err: cause}
}
