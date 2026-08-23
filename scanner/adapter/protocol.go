package adapter

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"regexp"
	"strings"
)

const ProtocolVersion = "prc-adapter-jsonl-v1"

var (
	hexDigestPattern = regexp.MustCompile(`^[0-9a-f]{64}$`)
	commitPattern    = regexp.MustCompile(`^[0-9a-f]{40,64}$`)
	sha256Pattern    = regexp.MustCompile(`^sha256:[0-9a-f]{64}$`)
)

type Limits struct {
	MaxLineBytes int `json:"max_line_bytes" yaml:"max_line_bytes"`
	MaxMessages  int `json:"max_messages" yaml:"max_messages"`
	MaxStdin     int `json:"max_stdin_bytes" yaml:"max_stdin_bytes"`
	MaxStdout    int `json:"max_stdout_bytes" yaml:"max_stdout_bytes"`
	MaxStderr    int `json:"max_stderr_bytes" yaml:"max_stderr_bytes"`
}

func DefaultLimits() Limits {
	return Limits{
		MaxLineBytes: 256 * 1024, MaxMessages: 10_000, MaxStdin: 1024 * 1024,
		MaxStdout: 4 * 1024 * 1024, MaxStderr: 64 * 1024,
	}
}

type Subject struct {
	TargetName      string `json:"target_name"`
	TargetCommit    string `json:"target_commit,omitempty"`
	InventoryDigest string `json:"inventory_digest"`
}

type Hello struct {
	Type     string `json:"type"`
	Protocol string `json:"protocol"`
	RunID    string `json:"run_id"`
}

type Input struct {
	Type    string         `json:"type"`
	Subject Subject        `json:"subject"`
	Facts   map[string]any `json:"facts"`
	Config  map[string]any `json:"config"`
}

type Execute struct {
	Type string `json:"type"`
}

type Log struct {
	Type    string `json:"type"`
	Level   string `json:"level"`
	Message string `json:"message"`
}

type Location struct {
	Path   string `json:"path"`
	Line   int    `json:"line,omitempty"`
	Column int    `json:"column,omitempty"`
}

type Observation struct {
	ID        string         `json:"id"`
	Kind      string         `json:"kind"`
	Outcome   string         `json:"outcome"`
	Summary   string         `json:"summary"`
	Locations []Location     `json:"locations"`
	Data      map[string]any `json:"data,omitempty"`
}

type ObservationMessage struct {
	Type        string      `json:"type"`
	Observation Observation `json:"observation"`
}

type Artifact struct {
	ID        string `json:"id"`
	MediaType string `json:"media_type"`
	Digest    string `json:"digest"`
	Size      int64  `json:"size"`
	Path      string `json:"path,omitempty"`
}

type ArtifactMessage struct {
	Type     string   `json:"type"`
	Artifact Artifact `json:"artifact"`
}

type Summary struct {
	Type   string         `json:"type"`
	Status string         `json:"status"`
	Counts map[string]int `json:"counts"`
	Reason string         `json:"reason,omitempty"`
}

type Transcript struct {
	Logs         []Log         `json:"logs"`
	Observations []Observation `json:"observations"`
	Artifacts    []Artifact    `json:"artifacts"`
	Summary      Summary       `json:"summary"`
}

func InputJSONL(runID string, subject Subject, facts, config map[string]any) ([]byte, error) {
	if !hexDigestPattern.MatchString(runID) {
		return nil, fmt.Errorf("adapter run ID must be a lowercase SHA-256 digest")
	}
	if strings.TrimSpace(subject.TargetName) == "" || !hexDigestPattern.MatchString(subject.InventoryDigest) {
		return nil, fmt.Errorf("adapter subject requires a target name and lowercase SHA-256 inventory digest")
	}
	if subject.TargetCommit != "" && !commitPattern.MatchString(subject.TargetCommit) {
		return nil, fmt.Errorf("adapter target commit must be 40 to 64 lowercase hexadecimal characters")
	}
	if facts == nil {
		facts = map[string]any{}
	}
	if config == nil {
		config = map[string]any{}
	}
	messages := []any{
		Hello{Type: "hello", Protocol: ProtocolVersion, RunID: runID},
		Input{Type: "input", Subject: subject, Facts: facts, Config: config},
		Execute{Type: "execute"},
	}
	var output bytes.Buffer
	encoder := json.NewEncoder(&output)
	encoder.SetEscapeHTML(false)
	for _, message := range messages {
		if err := encoder.Encode(message); err != nil {
			return nil, fmt.Errorf("encode adapter input: %w", err)
		}
	}
	return output.Bytes(), nil
}

func ParseOutput(input io.Reader, limits Limits) (Transcript, error) {
	if err := validateLimits(limits); err != nil {
		return Transcript{}, err
	}
	output, err := io.ReadAll(io.LimitReader(input, int64(limits.MaxStdout)+1))
	if err != nil {
		return Transcript{}, fmt.Errorf("read adapter output: %w", err)
	}
	if len(output) > limits.MaxStdout {
		return Transcript{}, fmt.Errorf("adapter output exceeds %d bytes", limits.MaxStdout)
	}
	scanner := bufio.NewScanner(bytes.NewReader(output))
	scanner.Buffer(make([]byte, min(limits.MaxLineBytes, 64*1024)), limits.MaxLineBytes)
	transcript := Transcript{
		Logs: make([]Log, 0), Observations: make([]Observation, 0), Artifacts: make([]Artifact, 0),
	}
	messageCount := 0
	summarySeen := false
	observationIDs := map[string]bool{}
	artifactIDs := map[string]bool{}
	for scanner.Scan() {
		messageCount++
		if messageCount > limits.MaxMessages {
			return Transcript{}, fmt.Errorf("adapter output exceeds %d messages", limits.MaxMessages)
		}
		line := scanner.Bytes()
		if len(bytes.TrimSpace(line)) == 0 {
			return Transcript{}, fmt.Errorf("adapter output line %d is blank", messageCount)
		}
		if err := rejectDuplicateKeys(line); err != nil {
			return Transcript{}, protocolError(messageCount, err)
		}
		if summarySeen {
			return Transcript{}, fmt.Errorf("adapter emitted a message after its summary")
		}
		var envelope struct {
			Type string `json:"type"`
		}
		if err := json.Unmarshal(line, &envelope); err != nil {
			return Transcript{}, fmt.Errorf("adapter output line %d is not valid JSON: %w", messageCount, err)
		}
		switch envelope.Type {
		case "log":
			var message Log
			if err := decodeStrict(line, &message); err != nil {
				return Transcript{}, protocolError(messageCount, err)
			}
			if err := validateLog(message); err != nil {
				return Transcript{}, protocolError(messageCount, err)
			}
			transcript.Logs = append(transcript.Logs, message)
		case "observation":
			var message ObservationMessage
			if err := decodeStrict(line, &message); err != nil {
				return Transcript{}, protocolError(messageCount, err)
			}
			if err := validateObservation(message.Observation); err != nil {
				return Transcript{}, protocolError(messageCount, err)
			}
			if observationIDs[message.Observation.ID] {
				return Transcript{}, protocolError(messageCount, fmt.Errorf("duplicate observation ID %q", message.Observation.ID))
			}
			observationIDs[message.Observation.ID] = true
			transcript.Observations = append(transcript.Observations, message.Observation)
		case "artifact":
			var message ArtifactMessage
			if err := decodeStrict(line, &message); err != nil {
				return Transcript{}, protocolError(messageCount, err)
			}
			if err := validateArtifact(message.Artifact); err != nil {
				return Transcript{}, protocolError(messageCount, err)
			}
			if artifactIDs[message.Artifact.ID] {
				return Transcript{}, protocolError(messageCount, fmt.Errorf("duplicate artifact ID %q", message.Artifact.ID))
			}
			artifactIDs[message.Artifact.ID] = true
			transcript.Artifacts = append(transcript.Artifacts, message.Artifact)
		case "summary":
			if err := decodeStrict(line, &transcript.Summary); err != nil {
				return Transcript{}, protocolError(messageCount, err)
			}
			if err := validateSummary(transcript.Summary); err != nil {
				return Transcript{}, protocolError(messageCount, err)
			}
			summarySeen = true
		default:
			return Transcript{}, fmt.Errorf("adapter output line %d has unsupported message type %q", messageCount, envelope.Type)
		}
	}
	if err := scanner.Err(); err != nil {
		if errors.Is(err, bufio.ErrTooLong) || strings.Contains(err.Error(), "token too long") {
			return Transcript{}, fmt.Errorf("adapter output line exceeds %d bytes", limits.MaxLineBytes)
		}
		return Transcript{}, fmt.Errorf("read adapter output: %w", err)
	}
	if !summarySeen {
		return Transcript{}, fmt.Errorf("adapter output ended without a summary")
	}
	actualCounts := map[string]int{
		"logs": len(transcript.Logs), "observations": len(transcript.Observations), "artifacts": len(transcript.Artifacts),
	}
	for name, actual := range actualCounts {
		if declared, ok := transcript.Summary.Counts[name]; ok && declared != actual {
			return Transcript{}, fmt.Errorf("adapter summary count %q is %d, observed %d", name, declared, actual)
		}
	}
	return transcript, nil
}

func validateLimits(limits Limits) error {
	if limits.MaxLineBytes < 1 || limits.MaxMessages < 1 || limits.MaxStdin < 1 || limits.MaxStdout < 1 || limits.MaxStderr < 1 {
		return fmt.Errorf("adapter protocol limits must all be positive")
	}
	return nil
}

func decodeStrict(data []byte, target any) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return err
	}
	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		if err == nil {
			return fmt.Errorf("multiple JSON values in one protocol line")
		}
		return err
	}
	return nil
}

func rejectDuplicateKeys(data []byte) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	var walk func() error
	walk = func() error {
		token, err := decoder.Token()
		if err != nil {
			return err
		}
		delimiter, ok := token.(json.Delim)
		if !ok {
			return nil
		}
		switch delimiter {
		case '{':
			seen := map[string]bool{}
			for decoder.More() {
				keyToken, err := decoder.Token()
				if err != nil {
					return err
				}
				key, ok := keyToken.(string)
				if !ok {
					return fmt.Errorf("object key is not a string")
				}
				if seen[key] {
					return fmt.Errorf("duplicate JSON key %q", key)
				}
				seen[key] = true
				if err := walk(); err != nil {
					return err
				}
			}
			_, err = decoder.Token()
			return err
		case '[':
			for decoder.More() {
				if err := walk(); err != nil {
					return err
				}
			}
			_, err = decoder.Token()
			return err
		default:
			return fmt.Errorf("unexpected JSON delimiter %q", delimiter)
		}
	}
	if err := walk(); err != nil {
		return err
	}
	if token, err := decoder.Token(); !errors.Is(err, io.EOF) {
		if err == nil {
			return fmt.Errorf("unexpected trailing JSON value %v", token)
		}
		return err
	}
	return nil
}

func protocolError(line int, err error) error {
	return fmt.Errorf("adapter protocol error on line %d: %w", line, err)
}

func validateLog(message Log) error {
	if message.Type != "log" {
		return fmt.Errorf("log message has invalid type")
	}
	if message.Level != "debug" && message.Level != "info" && message.Level != "warning" && message.Level != "error" {
		return fmt.Errorf("log level %q is invalid", message.Level)
	}
	if strings.TrimSpace(message.Message) == "" {
		return fmt.Errorf("log message is empty")
	}
	return nil
}

func validateObservation(observation Observation) error {
	if strings.TrimSpace(observation.ID) == "" || strings.TrimSpace(observation.Kind) == "" || strings.TrimSpace(observation.Summary) == "" {
		return fmt.Errorf("observation ID, kind, and summary are required")
	}
	if observation.Outcome != "found" && observation.Outcome != "not_found" && observation.Outcome != "value" && observation.Outcome != "unsupported" {
		return fmt.Errorf("observation outcome %q is invalid", observation.Outcome)
	}
	if observation.Locations == nil {
		return fmt.Errorf("observation locations must be an array")
	}
	for _, location := range observation.Locations {
		if err := validateRelativePath(location.Path); err != nil {
			return fmt.Errorf("observation location: %w", err)
		}
		if location.Line < 0 || location.Column < 0 {
			return fmt.Errorf("observation location coordinates cannot be negative")
		}
	}
	return nil
}

func validateArtifact(artifact Artifact) error {
	if strings.TrimSpace(artifact.ID) == "" || strings.TrimSpace(artifact.MediaType) == "" {
		return fmt.Errorf("artifact ID and media type are required")
	}
	if !sha256Pattern.MatchString(artifact.Digest) {
		return fmt.Errorf("artifact digest must be a lowercase SHA-256 descriptor")
	}
	if artifact.Size < 0 {
		return fmt.Errorf("artifact size cannot be negative")
	}
	if artifact.Path != "" {
		if err := validateRelativePath(artifact.Path); err != nil {
			return fmt.Errorf("artifact path: %w", err)
		}
	}
	return nil
}

func validateSummary(summary Summary) error {
	if summary.Type != "summary" {
		return fmt.Errorf("summary message has invalid type")
	}
	switch summary.Status {
	case "completed", "partial", "unsupported", "configuration_error", "execution_error", "timeout", "parse_error":
	default:
		return fmt.Errorf("summary status %q is invalid", summary.Status)
	}
	if summary.Status != "completed" && strings.TrimSpace(summary.Reason) == "" {
		return fmt.Errorf("non-completed summary requires a reason")
	}
	if summary.Counts == nil {
		return fmt.Errorf("summary counts must be an object")
	}
	for name, count := range summary.Counts {
		if strings.TrimSpace(name) == "" || count < 0 {
			return fmt.Errorf("summary counts must use nonempty names and nonnegative values")
		}
	}
	return nil
}

func validateRelativePath(path string) error {
	if path == "" || strings.Contains(path, "\\") || filepath.IsAbs(path) {
		return fmt.Errorf("path %q is not a nonempty relative path", path)
	}
	clean := filepath.Clean(filepath.FromSlash(path))
	if clean == "." || clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		return fmt.Errorf("path %q escapes its root", path)
	}
	if filepath.ToSlash(clean) != path {
		return fmt.Errorf("path %q is not normalized", path)
	}
	return nil
}
