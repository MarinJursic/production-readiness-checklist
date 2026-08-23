// Package mcpserver exposes the scanner's read-only inspect surface over the
// Model Context Protocol. It deliberately excludes adapters, providers,
// remediation, state mutation, and caller-selected filesystem paths.
package mcpserver

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"time"

	"github.com/MarinJursic/production-readiness-checklist/scanner/catalog"
	projectconfig "github.com/MarinJursic/production-readiness-checklist/scanner/config"
	"github.com/MarinJursic/production-readiness-checklist/scanner/engine"
	"github.com/MarinJursic/production-readiness-checklist/scanner/inventory"
	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
)

const (
	PlanResultSchema    = "prc.mcp-plan-result/v0.1"
	ScanResultSchema    = "prc.mcp-scan-result/v0.1"
	ExplainResultSchema = "prc.mcp-explain-result/v0.1"
	defaultProfileID    = "prc/core-repository"
)

// Options are resolved once at server startup. Tool callers cannot replace
// these trust-boundary values with arguments from a model.
type Options struct {
	CatalogRoot string
	Target      string
	ConfigPath  string
	ProfileID   string
	Now         func() time.Time
}

// Service owns the immutable catalog and path-locked target configuration for
// one MCP process. The target is inventoried again for every plan or scan so a
// client can assess edits performed outside this read-only server.
type Service struct {
	catalog    *catalog.Catalog
	target     string
	configPath string
	profileID  string
	now        func() time.Time
}

type PlanResult struct {
	SchemaVersion string     `json:"schema_version"`
	Plan          model.Plan `json:"plan"`
}

type InventorySummary struct {
	TargetName        string   `json:"target_name"`
	GitCommit         string   `json:"git_commit,omitempty"`
	Digest            string   `json:"digest"`
	FileCount         int      `json:"file_count"`
	SourceFiles       int      `json:"source_files"`
	PackageEcosystems []string `json:"package_ecosystems"`
}

type ResultSummary struct {
	Total          int `json:"total"`
	Pass           int `json:"pass"`
	Fail           int `json:"fail"`
	Unknown        int `json:"unknown"`
	Manual         int `json:"manual"`
	NotApplicable  int `json:"not_applicable"`
	Blocked        int `json:"blocked"`
	ExecutionError int `json:"execution_error"`
	Findings       int `json:"findings"`
}

// ScanResult omits the full file inventory and raw adapter transcripts from
// the MCP response. It retains the content-addressed inventory identity,
// evidence-linked results, and actionable findings needed by an agent.
type ScanResult struct {
	SchemaVersion  string                  `json:"schema_version"`
	RunID          string                  `json:"run_id"`
	StartedAt      time.Time               `json:"started_at"`
	CompletedAt    time.Time               `json:"completed_at"`
	TerminalState  string                  `json:"terminal_state"`
	ProfileID      string                  `json:"profile_id"`
	ProfileVersion string                  `json:"profile_version"`
	PlanDigest     string                  `json:"plan_digest"`
	Inventory      InventorySummary        `json:"inventory"`
	Summary        ResultSummary           `json:"summary"`
	Results        []model.AssertionResult `json:"results"`
	Findings       []model.Finding         `json:"findings"`
}

type ExplainResult struct {
	SchemaVersion string            `json:"schema_version"`
	Assertion     model.Assertion   `json:"assertion"`
	Objectives    []model.Objective `json:"objectives"`
}

func NewService(options Options) (*Service, error) {
	catalogRoot, err := canonicalDirectory(options.CatalogRoot, "catalog root")
	if err != nil {
		return nil, err
	}
	target, err := canonicalDirectory(options.Target, "target")
	if err != nil {
		return nil, err
	}
	loadedCatalog, err := catalog.Load(catalogRoot)
	if err != nil {
		return nil, fmt.Errorf("load locked catalog: %w", err)
	}

	configPath := ""
	profileID := options.ProfileID
	if options.ConfigPath != "" {
		configPath, err = canonicalFile(options.ConfigPath, "project configuration")
		if err != nil {
			return nil, err
		}
		validation, loadErr := projectconfig.Load(configPath)
		if loadErr != nil {
			return nil, fmt.Errorf("load locked project configuration: %w", loadErr)
		}
		if profileID == "" {
			profileID = validation.Configuration.Assessment.Profile
		} else if validation.Configuration.Assessment.Profile != profileID {
			return nil, fmt.Errorf("configured profile %s does not match locked profile %s",
				validation.Configuration.Assessment.Profile, profileID)
		}
	}
	if profileID == "" {
		profileID = defaultProfileID
	}
	if _, err := loadedCatalog.Profile(profileID); err != nil {
		return nil, err
	}
	now := options.Now
	if now == nil {
		now = func() time.Time { return time.Now().UTC() }
	}
	return &Service{
		catalog: loadedCatalog, target: target, configPath: configPath,
		profileID: profileID, now: now,
	}, nil
}

func (service *Service) Plan() (PlanResult, error) {
	item, err := service.currentInventory()
	if err != nil {
		return PlanResult{}, err
	}
	scanner := engine.New(service.catalog)
	plan, err := scanner.PlanMode(service.profileID, item, engine.ExecutionModeInspect)
	if err != nil {
		return PlanResult{}, err
	}
	return PlanResult{SchemaVersion: PlanResultSchema, Plan: plan}, nil
}

func (service *Service) Scan() (ScanResult, error) {
	item, err := service.currentInventory()
	if err != nil {
		return ScanResult{}, err
	}
	scanner := engine.New(service.catalog)
	scanner.Now = func() time.Time { return service.now().UTC() }
	run, err := scanner.ScanMode(service.profileID, item, nil, engine.ExecutionModeInspect)
	if err != nil {
		return ScanResult{}, err
	}
	return projectRun(run), nil
}

func (service *Service) Explain(assertionID string) (ExplainResult, error) {
	assertion, ok := service.catalog.Assertions[assertionID]
	if !ok {
		return ExplainResult{}, fmt.Errorf("unknown assertion %q", assertionID)
	}
	objectives := make([]model.Objective, 0)
	for _, objective := range service.catalog.Objectives {
		if slices.Contains(objective.AssertionIDs, assertionID) {
			objectives = append(objectives, objective)
		}
	}
	slices.SortFunc(objectives, func(left, right model.Objective) int {
		if left.ID < right.ID {
			return -1
		}
		if left.ID > right.ID {
			return 1
		}
		return 0
	})
	return ExplainResult{
		SchemaVersion: ExplainResultSchema, Assertion: assertion, Objectives: objectives,
	}, nil
}

func (service *Service) currentInventory() (model.Inventory, error) {
	var validation *projectconfig.Validation
	if service.configPath != "" {
		loaded, err := projectconfig.Load(service.configPath)
		if err != nil {
			return model.Inventory{}, fmt.Errorf("reload locked project configuration: %w", err)
		}
		if loaded.Configuration.Assessment.Profile != service.profileID {
			return model.Inventory{}, fmt.Errorf("project configuration changed its profile from locked value %s", service.profileID)
		}
		validation = &loaded
	}
	item, err := inventory.Build(service.target)
	if err != nil {
		return model.Inventory{}, err
	}
	if validation != nil {
		item, err = inventory.BindConfiguration(item, *validation, service.configPath)
		if err != nil {
			return model.Inventory{}, err
		}
	}
	return item, nil
}

func projectRun(run model.RunResult) ScanResult {
	summary := ResultSummary{Total: len(run.Results), Findings: len(run.Findings)}
	for _, result := range run.Results {
		switch result.Assessment {
		case "pass":
			summary.Pass++
		case "fail":
			summary.Fail++
		case "unknown":
			summary.Unknown++
		case "manual":
			summary.Manual++
		case "not_applicable":
			summary.NotApplicable++
		}
		switch result.Execution {
		case "blocked":
			summary.Blocked++
		case "error":
			summary.ExecutionError++
		}
	}
	return ScanResult{
		SchemaVersion: ScanResultSchema, RunID: run.RunID,
		StartedAt: run.StartedAt, CompletedAt: run.CompletedAt,
		TerminalState: run.TerminalState, ProfileID: run.Plan.ProfileID,
		ProfileVersion: run.Plan.ProfileVersion, PlanDigest: run.Plan.Digest,
		Inventory: InventorySummary{
			TargetName: run.Inventory.TargetName, GitCommit: run.Inventory.GitCommit,
			Digest: run.Inventory.Digest, FileCount: run.Inventory.FileCount,
			SourceFiles:       run.Inventory.SourceFiles,
			PackageEcosystems: append([]string(nil), run.Inventory.PackageEcosystems...),
		},
		Summary: summary, Results: run.Results, Findings: run.Findings,
	}
}

func canonicalDirectory(path, label string) (string, error) {
	if path == "" {
		path = "."
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("resolve %s: %w", label, err)
	}
	resolved, err := filepath.EvalSymlinks(abs)
	if err != nil {
		return "", fmt.Errorf("resolve %s symlinks: %w", label, err)
	}
	info, err := os.Stat(resolved)
	if err != nil {
		return "", fmt.Errorf("inspect %s: %w", label, err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("%s must be a directory", label)
	}
	return filepath.Clean(resolved), nil
}

func canonicalFile(path, label string) (string, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("resolve %s: %w", label, err)
	}
	resolved, err := filepath.EvalSymlinks(abs)
	if err != nil {
		return "", fmt.Errorf("resolve %s symlinks: %w", label, err)
	}
	info, err := os.Stat(resolved)
	if err != nil {
		return "", fmt.Errorf("inspect %s: %w", label, err)
	}
	if !info.Mode().IsRegular() {
		return "", fmt.Errorf("%s must be a regular file", label)
	}
	return filepath.Clean(resolved), nil
}
