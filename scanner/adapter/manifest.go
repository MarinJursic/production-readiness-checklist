package adapter

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
	"gopkg.in/yaml.v3"
)

const (
	ManifestSchema      = "prc.adapter-manifest/v0.4"
	OutputSchemaVersion = "prc.adapter-message/v0.1"
)

var (
	adapterIDPattern       = regexp.MustCompile(`^prc\.adapter\.[a-z0-9.-]+@[0-9]+\.[0-9]+$`)
	imagePattern           = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:/-]*@sha256:[0-9a-f]{64}$`)
	publisherIDPattern     = regexp.MustCompile(`^[a-z0-9][a-z0-9.-]{0,127}$`)
	observationKindPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9.-]{0,127}$`)
	versionTokenPattern    = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._+-]{0,127}$`)
	engineAPIPattern       = regexp.MustCompile(`^prc\.engine/v[0-9]+\.[0-9]+$`)
)

type Publisher struct {
	ID   string `json:"id" yaml:"id"`
	Name string `json:"name" yaml:"name"`
}

type ToolFormat struct {
	Name     string   `json:"name" yaml:"name"`
	Versions []string `json:"versions" yaml:"versions"`
}

type Tool struct {
	Name     string       `json:"name" yaml:"name"`
	Version  string       `json:"version" yaml:"version"`
	Upstream string       `json:"upstream" yaml:"upstream"`
	Formats  []ToolFormat `json:"formats" yaml:"formats"`
}

type Compatibility struct {
	EngineAPIs []string `json:"engine_apis" yaml:"engine_apis"`
}

type Capabilities struct {
	ReadWorkspace  bool     `json:"read_workspace" yaml:"read_workspace"`
	WriteScratch   bool     `json:"write_scratch" yaml:"write_scratch"`
	Network        string   `json:"network" yaml:"network"`
	NetworkHosts   []string `json:"network_hosts" yaml:"network_hosts"`
	SecretHandles  []string `json:"secret_handles" yaml:"secret_handles"`
	ChildProcesses bool     `json:"child_processes" yaml:"child_processes"`
}

type Resources struct {
	TimeoutSeconds int     `json:"timeout_seconds" yaml:"timeout_seconds"`
	MemoryMB       int     `json:"memory_mb" yaml:"memory_mb"`
	CPUs           float64 `json:"cpus" yaml:"cpus"`
	PIDs           int     `json:"pids" yaml:"pids"`
	TmpfsMB        int     `json:"tmpfs_mb" yaml:"tmpfs_mb"`
	Limits         `json:",inline" yaml:",inline"`
}

type Manifest struct {
	SchemaVersion    string        `json:"schema_version" yaml:"schema_version"`
	ID               string        `json:"id" yaml:"id"`
	Title            string        `json:"title" yaml:"title"`
	Description      string        `json:"description" yaml:"description"`
	Publisher        Publisher     `json:"publisher" yaml:"publisher"`
	Owner            string        `json:"owner" yaml:"owner"`
	Maintenance      string        `json:"maintenance" yaml:"maintenance"`
	Protocol         string        `json:"protocol" yaml:"protocol"`
	OutputSchema     string        `json:"output_schema" yaml:"output_schema"`
	ObservationKinds []string      `json:"observation_kinds" yaml:"observation_kinds"`
	Compatibility    Compatibility `json:"compatibility" yaml:"compatibility"`
	Tool             Tool          `json:"tool" yaml:"tool"`
	Limitations      []string      `json:"limitations" yaml:"limitations"`
	Runner           string        `json:"runner" yaml:"runner"`
	Image            string        `json:"image" yaml:"image"`
	Command          []string      `json:"command" yaml:"command"`
	Capabilities     Capabilities  `json:"capabilities" yaml:"capabilities"`
	Resources        Resources     `json:"resources" yaml:"resources"`
}

func LoadManifest(path string) (Manifest, error) {
	info, err := os.Stat(path)
	if err != nil {
		return Manifest{}, fmt.Errorf("inspect adapter manifest: %w", err)
	}
	if !info.Mode().IsRegular() || info.Size() > 1024*1024 {
		return Manifest{}, fmt.Errorf("adapter manifest must be a regular file no larger than 1 MiB")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return Manifest{}, fmt.Errorf("read adapter manifest: %w", err)
	}
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)
	var manifest Manifest
	if err := decoder.Decode(&manifest); err != nil {
		return Manifest{}, fmt.Errorf("decode adapter manifest: %w", err)
	}
	var trailing any
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		if err == nil {
			return Manifest{}, fmt.Errorf("adapter manifest contains more than one YAML document")
		}
		return Manifest{}, fmt.Errorf("decode trailing adapter manifest content: %w", err)
	}
	if err := manifest.Validate(); err != nil {
		return Manifest{}, err
	}
	return manifest, nil
}

func (manifest Manifest) Validate() error {
	if manifest.SchemaVersion != ManifestSchema {
		return fmt.Errorf("unsupported adapter manifest schema %q", manifest.SchemaVersion)
	}
	if !adapterIDPattern.MatchString(manifest.ID) {
		return fmt.Errorf("invalid adapter ID %q", manifest.ID)
	}
	if strings.TrimSpace(manifest.Title) == "" || strings.TrimSpace(manifest.Description) == "" {
		return fmt.Errorf("adapter title and description are required")
	}
	if !publisherIDPattern.MatchString(manifest.Publisher.ID) || strings.TrimSpace(manifest.Publisher.Name) == "" {
		return fmt.Errorf("adapter publisher identity is invalid")
	}
	if strings.TrimSpace(manifest.Owner) == "" {
		return fmt.Errorf("adapter owner is required")
	}
	if manifest.Maintenance != "active" && manifest.Maintenance != "deprecated" {
		return fmt.Errorf("unsupported adapter maintenance state %q", manifest.Maintenance)
	}
	switch manifest.Protocol {
	case ProtocolVersion:
		if manifest.OutputSchema != OutputSchemaVersion {
			return fmt.Errorf("unsupported JSONL adapter output schema %q", manifest.OutputSchema)
		}
	case GitleaksProtocolVersion:
		if manifest.OutputSchema != GitleaksOutputSchemaVersion {
			return fmt.Errorf("unsupported Gitleaks adapter output schema %q", manifest.OutputSchema)
		}
	case SyftProtocolVersion:
		if manifest.OutputSchema != SyftOutputSchemaVersion {
			return fmt.Errorf("unsupported Syft adapter output schema %q", manifest.OutputSchema)
		}
	default:
		return fmt.Errorf("unsupported adapter protocol %q", manifest.Protocol)
	}
	if !uniqueTokens(manifest.ObservationKinds, observationKindPattern) {
		return fmt.Errorf("adapter observation kinds must be a nonempty unique token array")
	}
	if !uniqueTokens(manifest.Compatibility.EngineAPIs, engineAPIPattern) {
		return fmt.Errorf("adapter engine compatibility must be a nonempty unique API array")
	}
	if err := validateTool(manifest.Tool); err != nil {
		return err
	}
	if !uniqueNonemptyText(manifest.Limitations) {
		return fmt.Errorf("adapter limitations must be a nonempty unique text array")
	}
	if manifest.Runner != "oci" {
		return fmt.Errorf("external adapters require the OCI runner")
	}
	if !imagePattern.MatchString(manifest.Image) {
		return fmt.Errorf("adapter image must be a fully qualified immutable SHA-256 reference")
	}
	registry := strings.SplitN(manifest.Image, "/", 2)[0]
	if !strings.Contains(manifest.Image, "/") ||
		(!strings.Contains(registry, ".") && !strings.Contains(registry, ":") && registry != "localhost") {
		return fmt.Errorf("adapter image must include an explicit registry host")
	}
	if len(manifest.Command) == 0 || len(manifest.Command) > 128 || strings.TrimSpace(manifest.Command[0]) == "" {
		return fmt.Errorf("adapter command must be a nonempty argument array")
	}
	commandBytes := 0
	for _, argument := range manifest.Command {
		if strings.ContainsRune(argument, '\x00') {
			return fmt.Errorf("adapter command contains a NUL byte")
		}
		commandBytes += len(argument)
	}
	if commandBytes > 256*1024 {
		return fmt.Errorf("adapter command exceeds 256 KiB")
	}
	if !manifest.Capabilities.ReadWorkspace {
		return fmt.Errorf("the current OCI adapter protocol requires read_workspace")
	}
	if manifest.Capabilities.Network != "deny" {
		return fmt.Errorf("network mode %q is not supported; the current runner fails closed to deny", manifest.Capabilities.Network)
	}
	if len(manifest.Capabilities.NetworkHosts) != 0 {
		return fmt.Errorf("network hosts must be empty when network is denied")
	}
	if len(manifest.Capabilities.SecretHandles) != 0 {
		return fmt.Errorf("secret handles are not supported by the current runner")
	}
	if manifest.Protocol == ProtocolVersion && manifest.Capabilities.ChildProcesses {
		return fmt.Errorf("child-process capability is not supported by the current runner")
	}
	if err := validateResources(manifest.Resources); err != nil {
		return err
	}
	if manifest.Protocol == GitleaksProtocolVersion {
		if err := validateGitleaksManifest(manifest); err != nil {
			return err
		}
	} else if manifest.Protocol == SyftProtocolVersion {
		if err := validateSyftManifest(manifest); err != nil {
			return err
		}
	}
	return nil
}

func validateTool(tool Tool) error {
	if strings.TrimSpace(tool.Name) == "" || !versionTokenPattern.MatchString(tool.Version) || strings.EqualFold(tool.Version, "latest") {
		return fmt.Errorf("adapter tool requires a name and immutable version token")
	}
	upstream, err := url.Parse(tool.Upstream)
	if err != nil || upstream.Scheme != "https" || upstream.Host == "" || upstream.User != nil {
		return fmt.Errorf("adapter tool upstream must be an absolute HTTPS URL without credentials")
	}
	if len(tool.Formats) == 0 {
		return fmt.Errorf("adapter tool must declare at least one supported format")
	}
	seen := map[string]bool{}
	for _, format := range tool.Formats {
		if !observationKindPattern.MatchString(format.Name) || seen[format.Name] ||
			!uniqueTokens(format.Versions, versionTokenPattern) {
			return fmt.Errorf("adapter tool formats must have unique names and nonempty unique versions")
		}
		seen[format.Name] = true
	}
	return nil
}

func uniqueTokens(values []string, pattern *regexp.Regexp) bool {
	if len(values) == 0 {
		return false
	}
	seen := map[string]bool{}
	for _, value := range values {
		if !pattern.MatchString(value) || seen[value] {
			return false
		}
		seen[value] = true
	}
	return true
}

func uniqueNonemptyText(values []string) bool {
	if len(values) == 0 {
		return false
	}
	seen := map[string]bool{}
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" || seen[value] {
			return false
		}
		seen[value] = true
	}
	return true
}

func (manifest Manifest) SupportsEngine(engineAPI string) bool {
	for _, supported := range manifest.Compatibility.EngineAPIs {
		if supported == engineAPI {
			return true
		}
	}
	return false
}

func (manifest Manifest) ValidateForCurrentEngine() error {
	if err := manifest.Validate(); err != nil {
		return err
	}
	if !manifest.SupportsEngine(model.EngineVersion) {
		return fmt.Errorf("adapter %s does not support engine API %s", manifest.ID, model.EngineVersion)
	}
	return nil
}

func validateResources(resources Resources) error {
	if resources.TimeoutSeconds < 1 || resources.TimeoutSeconds > 3600 {
		return fmt.Errorf("adapter timeout must be between 1 and 3600 seconds")
	}
	if resources.MemoryMB < 32 || resources.MemoryMB > 32768 {
		return fmt.Errorf("adapter memory must be between 32 and 32768 MiB")
	}
	if math.IsNaN(resources.CPUs) || math.IsInf(resources.CPUs, 0) || resources.CPUs <= 0 || resources.CPUs > 32 {
		return fmt.Errorf("adapter CPUs must be greater than zero and at most 32")
	}
	if resources.PIDs < 1 || resources.PIDs > 4096 {
		return fmt.Errorf("adapter PID limit must be between 1 and 4096")
	}
	if resources.TmpfsMB < 1 || resources.TmpfsMB > 4096 {
		return fmt.Errorf("adapter tmpfs must be between 1 and 4096 MiB")
	}
	if err := validateLimits(resources.Limits); err != nil {
		return err
	}
	if resources.MaxLineBytes > 4*1024*1024 || resources.MaxMessages > 1_000_000 ||
		resources.MaxStdin > 16*1024*1024 || resources.MaxStdout > 256*1024*1024 ||
		resources.MaxStderr > 16*1024*1024 {
		return fmt.Errorf("adapter protocol limits exceed scanner safety ceilings")
	}
	if resources.MaxLineBytes > resources.MaxStdout {
		return fmt.Errorf("adapter line limit cannot exceed stdout limit")
	}
	return nil
}

func (manifest Manifest) Timeout() time.Duration {
	return time.Duration(manifest.Resources.TimeoutSeconds) * time.Second
}
