package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/MarinJursic/production-readiness-checklist/scanner/adapter"
	"github.com/MarinJursic/production-readiness-checklist/scanner/catalog"
	projectconfig "github.com/MarinJursic/production-readiness-checklist/scanner/config"
	"github.com/MarinJursic/production-readiness-checklist/scanner/doctor"
	"github.com/MarinJursic/production-readiness-checklist/scanner/engine"
	"github.com/MarinJursic/production-readiness-checklist/scanner/evidence"
	"github.com/MarinJursic/production-readiness-checklist/scanner/invalidation"
	"github.com/MarinJursic/production-readiness-checklist/scanner/inventory"
	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
	"github.com/MarinJursic/production-readiness-checklist/scanner/provider"
	"github.com/MarinJursic/production-readiness-checklist/scanner/remediation"
	"github.com/MarinJursic/production-readiness-checklist/scanner/report"
	"github.com/MarinJursic/production-readiness-checklist/scanner/state"
)

const version = "0.1.0-dev"

const (
	exitSuccess           = 0
	exitGateFailed        = 1
	exitIncomplete        = 2
	exitConfiguration     = 3
	exitExecution         = 4
	exitPolicyDenied      = 5
	exitInternal          = 6
	exitCancelled         = 7
	exitCandidateRejected = 8
)

type classifiedError struct {
	code int
	err  error
}

func (err classifiedError) Error() string { return err.err.Error() }
func (err classifiedError) Unwrap() error { return err.err }

func exitError(code int, err error) error {
	if err == nil {
		return nil
	}
	var classified classifiedError
	if errors.As(err, &classified) || errors.Is(err, context.Canceled) {
		return err
	}
	return classifiedError{code: code, err: err}
}

func errorExitCode(err error, fallback int) int {
	if errors.Is(err, context.Canceled) {
		return exitCancelled
	}
	var classified classifiedError
	if errors.As(err, &classified) {
		return classified.code
	}
	return fallback
}

func remediationCommandError(err error) error {
	if remediation.IsPolicyDenied(err) {
		return exitError(exitPolicyDenied, err)
	}
	if remediation.IsProviderExecution(err) {
		return exitError(exitExecution, err)
	}
	return exitError(exitConfiguration, err)
}

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		usage(stderr)
		return exitConfiguration
	}
	var err error
	outcome := exitSuccess
	errorFallback := exitInternal
	switch args[0] {
	case "version":
		fmt.Fprintf(stdout, "prc %s\n", version)
		return 0
	case "inventory":
		errorFallback = exitConfiguration
		err = runInventory(args[1:], stdout, stderr)
	case "config":
		errorFallback = exitConfiguration
		err = runConfig(args[1:], stdout, stderr)
	case "catalog":
		errorFallback = exitConfiguration
		err = runCatalog(args[1:], stdout, stderr)
	case "plan":
		errorFallback = exitConfiguration
		err = runPlan(args[1:], stdout, stderr)
	case "scan":
		outcome, err = runScan(args[1:], stdout, stderr)
	case "fix":
		outcome, err = runFix(args[1:], stdout, stderr)
	case "doctor":
		outcome, err = runDoctor(args[1:], stdout, stderr)
	case "history":
		errorFallback = exitConfiguration
		err = runHistory(args[1:], stdout, stderr)
	case "diff":
		errorFallback = exitConfiguration
		err = runDiff(args[1:], stdout, stderr)
	case "remediate":
		outcome, err = runRemediate(args[1:], stdout, stderr)
	case "remediate-proposal":
		outcome, err = runRemediateProposal(args[1:], stdout, stderr)
	case "explain":
		errorFallback = exitConfiguration
		err = runExplain(args[1:], stdout, stderr)
	case "adapter":
		errorFallback = exitExecution
		err = runAdapter(args[1:], stdout, stderr)
	case "provider":
		errorFallback = exitExecution
		err = runProvider(args[1:], stdout, stderr)
	case "help", "-h", "--help":
		usage(stdout)
		return 0
	default:
		fmt.Fprintf(stderr, "unknown command %q\n", args[0])
		usage(stderr)
		return exitConfiguration
	}
	if err != nil {
		code := errorExitCode(err, errorFallback)
		fmt.Fprintf(stderr, "error [PRC-EXIT-%d]: %v\n", code, err)
		return code
	}
	return outcome
}

func usage(output io.Writer) {
	fmt.Fprintln(output, "Production Readiness Scanner")
	fmt.Fprintln(output, "usage: prc <catalog|config|inventory|plan|scan|diff|fix|remediate|remediate-proposal|doctor|history|explain|adapter|provider|version> [options]")
}

func runCatalog(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 || (args[0] != "validate" && args[0] != "bundle") {
		return errors.New("catalog requires validate or bundle")
	}
	command := args[0]
	set := flag.NewFlagSet("catalog "+command, flag.ContinueOnError)
	set.SetOutput(stderr)
	root := set.String("catalog-root", ".", "repository containing the PRC catalog")
	format := "human"
	if command == "validate" {
		set.StringVar(&format, "format", "human", "human or json")
	}
	if err := set.Parse(args[1:]); err != nil {
		return err
	}
	if set.NArg() != 0 {
		return fmt.Errorf("unexpected catalog arguments: %s", strings.Join(set.Args(), " "))
	}
	if command == "validate" && format != "human" && format != "json" {
		return fmt.Errorf("unsupported format %q", format)
	}
	loaded, err := catalog.Load(*root)
	if err != nil {
		return err
	}
	if command == "bundle" {
		bundle, err := loaded.Bundle()
		if err != nil {
			return err
		}
		return encodeJSON(stdout, bundle)
	}
	manifest, err := loaded.Manifest()
	if err != nil {
		return err
	}
	if format == "json" {
		return encodeJSON(stdout, manifest)
	}
	fmt.Fprintf(stdout, "Catalog: %s\n", manifest.CatalogDigest)
	fmt.Fprintf(stdout, "Version: %s\n", manifest.CatalogVersion)
	fmt.Fprintf(stdout, "Definitions: %d objectives, %d assertions, %d profiles\n",
		manifest.ObjectiveCount, manifest.AssertionCount, manifest.ProfileCount)
	return nil
}

func runConfig(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 || args[0] != "validate" {
		return errors.New("config requires validate")
	}
	set := flag.NewFlagSet("config validate", flag.ContinueOnError)
	set.SetOutput(stderr)
	path := set.String("file", "production-readiness.yaml", "project configuration YAML file")
	format := set.String("format", "human", "human or json")
	if err := set.Parse(args[1:]); err != nil {
		return err
	}
	if *format != "human" && *format != "json" {
		return fmt.Errorf("unsupported format %q", *format)
	}
	validation, err := projectconfig.Load(*path)
	if err != nil {
		return err
	}
	if *format == "json" {
		return encodeJSON(stdout, validation)
	}
	document := validation.Configuration
	fmt.Fprintf(stdout, "Configuration: %s\n", validation.Digest)
	fmt.Fprintf(stdout, "Project: %s (%s, %s risk)\n", document.Project.Name, document.Project.ID, document.Project.RiskProfile)
	fmt.Fprintf(stdout, "Profile: %s\n", document.Assessment.Profile)
	fmt.Fprintf(stdout, "Execution: network=%s commands=%d production-connected=%t\n",
		document.Execution.Network, len(document.Execution.AllowCommands), document.Execution.ProductionConnected)
	fmt.Fprintf(stdout, "Remediation: enabled=%t attempts=%d files=%d lines=%d\n",
		document.Remediation.Enabled, document.Remediation.MaxAttempts,
		document.Remediation.MaxFiles, document.Remediation.MaxChangedLines)
	return nil
}

func runProvider(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		return errors.New("provider requires capabilities, seal-task, plan, validate-output, or run")
	}
	set := flag.NewFlagSet("provider "+args[0], flag.ContinueOnError)
	set.SetOutput(stderr)
	providerName := set.String("provider", "", "codex or claude")
	if args[0] == "seal-task" {
		path := set.String("file", "", "draft agent task JSON file")
		workspace := set.String("workspace", ".", "workspace to bind into the task")
		configPath := set.String("config", "", "optional validated project configuration")
		if err := set.Parse(args[1:]); err != nil {
			return err
		}
		if *path == "" {
			return errors.New("--file is required")
		}
		item, validation, err := configuredInventory(*workspace, *configPath)
		if err != nil {
			return err
		}
		var requiredProtectedPaths []string
		if validation != nil {
			configuration := &remediation.ProjectConfiguration{Validation: *validation, SourcePath: *configPath}
			requiredProtectedPaths, err = remediation.RequiredProtectedPaths(*workspace, configuration)
			if err != nil {
				return err
			}
		}
		task, err := provider.SealTaskWithInventory(*path, *workspace, item, requiredProtectedPaths)
		if err != nil {
			return err
		}
		return encodeJSON(stdout, task)
	}
	if args[0] == "capabilities" {
		if err := set.Parse(args[1:]); err != nil {
			return err
		}
		capabilities, err := provider.ProviderCapabilities(*providerName)
		if err != nil {
			return err
		}
		return encodeJSON(stdout, capabilities)
	}
	taskPath := set.String("task", "", "agent task JSON file")
	if args[0] == "validate-output" {
		path := set.String("file", "", "provider output file")
		if err := set.Parse(args[1:]); err != nil {
			return err
		}
		if *taskPath == "" || *path == "" {
			return errors.New("--task and --file are required")
		}
		task, err := provider.LoadTask(*taskPath)
		if err != nil {
			return err
		}
		data, err := readBoundedRegularFile(*path, int64(task.MaxOutputBytes))
		if err != nil {
			return err
		}
		output, err := provider.ParseOutput(*providerName, data, task)
		if err != nil {
			return err
		}
		return encodeJSON(stdout, output)
	}
	if args[0] != "plan" && args[0] != "run" {
		return fmt.Errorf("unknown provider command %q", args[0])
	}
	workspace := set.String("workspace", ".", "read-only target workspace")
	outputDirectory := set.String("output-dir", "", "existing disjoint output directory")
	schemaPath := set.String("output-schema", "schemas/agent-output.schema.json", "agent output JSON Schema")
	executable := set.String("executable", "", "provider CLI executable")
	if err := set.Parse(args[1:]); err != nil {
		return err
	}
	if *providerName == "" || *taskPath == "" || *outputDirectory == "" {
		return errors.New("--provider, --task, and --output-dir are required")
	}
	if *executable == "" {
		*executable = *providerName
	}
	task, err := provider.LoadTask(*taskPath)
	if err != nil {
		return err
	}
	launchPlan, err := provider.BuildPlan(*providerName, *executable, *workspace, *outputDirectory, *schemaPath, task)
	if err != nil {
		return err
	}
	if args[0] == "plan" {
		return encodeJSON(stdout, launchPlan)
	}
	execution, err := provider.Run(context.Background(), launchPlan, task)
	if err != nil {
		return err
	}
	return encodeJSON(stdout, execution)
}

func runRemediate(args []string, stdout, stderr io.Writer) (int, error) {
	set, target, catalogRoot, profile := parseCommon("remediate", args, stderr)
	configPath := set.String("config", "", "optional validated project configuration")
	assertionID := set.String("assertion", "PRC-A-CORE-014", "registered R1 assertion ID to remediate")
	candidateDirectory := set.String("candidate-dir", "", "new isolated candidate directory (required)")
	maxFiles := set.Int("max-files", 20, "maximum files the fix may change")
	maxChangedLines := set.Int("max-changed-lines", 20, "maximum changed lines")
	format := set.String("format", "human", "human or json")
	if err := set.Parse(args); err != nil {
		return exitInternal, exitError(exitConfiguration, err)
	}
	if *format != "human" && *format != "json" {
		return exitInternal, exitError(exitConfiguration, fmt.Errorf("unsupported format %q", *format))
	}
	configuration, err := loadRemediationConfiguration(*configPath)
	if err != nil {
		return exitInternal, remediationCommandError(err)
	}
	maxAttempts := 1
	if configuration != nil {
		document := configuration.Validation.Configuration
		if !flagWasSet(set, "profile") {
			*profile = document.Assessment.Profile
		}
		if !flagWasSet(set, "max-files") {
			*maxFiles = 0
		}
		if !flagWasSet(set, "max-changed-lines") {
			*maxChangedLines = 0
		}
		maxAttempts = document.Remediation.MaxAttempts
	}
	candidate, err := remediation.Run(remediation.Options{
		CatalogRoot: *catalogRoot, Target: *target, CandidateDir: *candidateDirectory,
		ProfileID: *profile, AssertionID: *assertionID,
		MaxFiles: *maxFiles, MaxChangedLines: *maxChangedLines,
		Attempt: 1, MaxAttempts: maxAttempts, Configuration: configuration,
	})
	if err != nil {
		return exitInternal, remediationCommandError(err)
	}
	if *format == "json" {
		if err := encodeJSON(stdout, candidate); err != nil {
			return exitInternal, exitError(exitInternal, err)
		}
	} else {
		printCandidate(stdout, candidate)
	}
	if !candidate.Accepted {
		return exitCandidateRejected, nil
	}
	return exitSuccess, nil
}

func runFix(args []string, stdout, stderr io.Writer) (int, error) {
	set, target, catalogRoot, profile := parseCommon("fix", args, stderr)
	configPath := set.String("config", "", "optional validated project configuration")
	candidateRoot := set.String("candidate-root", "", "new root for isolated attempt workspaces (required)")
	maxAttempts := set.Int("max-attempts", 3, "maximum remediation attempts")
	maxFiles := set.Int("max-files", 20, "maximum changed files across all attempts")
	maxChangedLines := set.Int("max-changed-lines", 200, "maximum changed lines across all attempts")
	providerName := set.String("provider", "", "optional suggest-only provider: codex or claude")
	providerExecutable := set.String("provider-executable", "", "provider CLI executable; defaults to provider name")
	agentOutputSchema := set.String("agent-output-schema", "", "agent output JSON schema; defaults to catalog schema")
	allowRemoteSource := set.Bool("allow-remote-source-processing", false, "explicitly allow sealed task inputs to be processed by the provider")
	agentTimeout := set.Int("agent-timeout-seconds", 300, "provider timeout per attempt")
	agentMaxOutput := set.Int("agent-max-output-bytes", 256*1024, "provider output byte limit per attempt")
	agentMaxCost := set.Float64("agent-max-cost-usd", 0, "provider-enforced cost limit per attempt; unsupported by Codex")
	format := set.String("format", "human", "human or json")
	if err := set.Parse(args); err != nil {
		return exitInternal, exitError(exitConfiguration, err)
	}
	if *format != "human" && *format != "json" {
		return exitInternal, exitError(exitConfiguration, fmt.Errorf("unsupported format %q", *format))
	}
	if *providerName == "" {
		for _, name := range []string{"provider-executable", "agent-output-schema", "allow-remote-source-processing", "agent-timeout-seconds", "agent-max-output-bytes", "agent-max-cost-usd"} {
			if flagWasSet(set, name) {
				return exitInternal, exitError(exitConfiguration, fmt.Errorf("--%s requires --provider", name))
			}
		}
	}
	configuration, err := loadRemediationConfiguration(*configPath)
	if err != nil {
		return exitInternal, remediationCommandError(err)
	}
	if configuration != nil {
		document := configuration.Validation.Configuration
		if !flagWasSet(set, "profile") {
			*profile = document.Assessment.Profile
		}
		if !flagWasSet(set, "max-attempts") {
			*maxAttempts = document.Remediation.MaxAttempts
		}
		if !flagWasSet(set, "max-files") {
			*maxFiles = 0
		}
		if !flagWasSet(set, "max-changed-lines") {
			*maxChangedLines = 0
		}
	}
	var agent *remediation.AgentOptions
	if *providerName != "" {
		if *providerExecutable == "" {
			*providerExecutable = *providerName
		}
		if *agentOutputSchema == "" {
			*agentOutputSchema = filepath.Join(*catalogRoot, "schemas", "agent-output.schema.json")
		}
		agent = &remediation.AgentOptions{
			Provider: *providerName, Executable: *providerExecutable, OutputSchemaPath: *agentOutputSchema,
			AllowRemoteSourceProcessing: *allowRemoteSource, TimeoutSeconds: *agentTimeout,
			MaxOutputBytes: *agentMaxOutput, MaxCostUSD: *agentMaxCost,
		}
	}
	result, err := remediation.RunLoop(remediation.LoopOptions{
		CatalogRoot: *catalogRoot, Target: *target, CandidateRoot: *candidateRoot,
		ProfileID: *profile, MaxAttempts: *maxAttempts, MaxFiles: *maxFiles,
		MaxChangedLines: *maxChangedLines, Configuration: configuration, Agent: agent,
	})
	if err != nil {
		return exitInternal, remediationCommandError(err)
	}
	if *format == "json" {
		if err := encodeJSON(stdout, result); err != nil {
			return exitInternal, exitError(exitInternal, err)
		}
	} else {
		printRemediationRun(stdout, result)
	}
	return remediationExitCode(result), nil
}

func remediationExitCode(result remediation.RemediationRun) int {
	switch result.TerminalState {
	case "profile_satisfied":
		return exitSuccess
	case "candidate_rejected":
		return exitCandidateRejected
	case "stopped_by_policy_or_budget":
		return exitPolicyDenied
	case "provider_stopped":
		return exitIncomplete
	}
	return scanTerminalExitCode(result.GateState)
}

func printRemediationRun(output io.Writer, result remediation.RemediationRun) {
	fmt.Fprintf(output, "Remediation run: %s\n", result.RunID)
	fmt.Fprintf(output, "Terminal state: %s\n", result.TerminalState)
	fmt.Fprintf(output, "Assessment gate: %s\n", result.GateState)
	fmt.Fprintf(output, "Result workspace: %s\n", result.ResultWorkspace)
	fmt.Fprintf(output, "Original unchanged: %t\n", result.OriginalUnchanged)
	fmt.Fprintf(output, "Budget: %d/%d attempts, %d/%d files, %d/%d lines\n",
		result.Usage.Attempts, result.MaxAttempts, result.Usage.ChangedFiles, result.MaxFiles,
		result.Usage.ChangedLines, result.MaxChangedLines)
	for _, candidate := range result.Candidates {
		status := "rejected"
		if candidate.Accepted {
			status = "accepted"
		}
		fmt.Fprintf(output, "- %s candidate %d: %s (%s)\n", status, candidate.Contract.Attempt,
			candidate.Contract.AssertionID, candidate.CandidatePath)
	}
	for _, execution := range result.ProviderExecutions {
		fmt.Fprintf(output, "- provider %s execution %s: %s\n", execution.Provider, execution.ExecutionID, execution.Output.Status)
	}
	for _, item := range result.Remaining {
		fmt.Fprintf(output, "- remaining %s [%s]: %s\n", item.AssertionID, item.ReasonCode, item.Reason)
	}
	for _, reason := range result.StopReasons {
		fmt.Fprintf(output, "- stop: %s\n", reason)
	}
}

type repeatedStringFlag []string

func (values *repeatedStringFlag) String() string { return strings.Join(*values, ",") }

func (values *repeatedStringFlag) Set(value string) error {
	*values = append(*values, value)
	return nil
}

func runDoctor(args []string, stdout, stderr io.Writer) (int, error) {
	set := flag.NewFlagSet("doctor", flag.ContinueOnError)
	set.SetOutput(stderr)
	target := set.String("target", ".", "repository to inspect")
	catalogRoot := set.String("catalog-root", ".", "repository containing the PRC catalog")
	stateDirectory := set.String("state-dir", "", "existing private state directory to probe")
	candidateParent := set.String("candidate-parent", "", "existing external directory to probe for isolated candidates")
	ociRuntime := set.String("oci-runtime", "", "docker or podman executable to inspect without running")
	format := set.String("format", "human", "human or json")
	providers := repeatedStringFlag{}
	set.Var(&providers, "provider", "codex or claude executable to inspect without running; repeatable")
	if err := set.Parse(args); err != nil {
		return exitInternal, exitError(exitConfiguration, err)
	}
	if *format != "human" && *format != "json" {
		return exitInternal, exitError(exitConfiguration, fmt.Errorf("unsupported format %q", *format))
	}
	report := doctor.Run(doctor.Options{
		Target: *target, CatalogRoot: *catalogRoot, StateDirectory: *stateDirectory,
		CandidateParent: *candidateParent, OCIRuntime: *ociRuntime, Providers: providers,
	})
	if *format == "json" {
		if err := encodeJSON(stdout, report); err != nil {
			return exitInternal, exitError(exitInternal, err)
		}
	} else {
		printDoctor(stdout, report)
	}
	if !report.Ready {
		return exitIncomplete, nil
	}
	return exitSuccess, nil
}

func printDoctor(output io.Writer, report doctor.Report) {
	fmt.Fprintf(output, "Scanner environment ready: %t\n", report.Ready)
	fmt.Fprintf(output, "Platform: %s/%s\n", report.Platform, report.Architecture)
	fmt.Fprintf(output, "Target: %s\n", report.Target)
	fmt.Fprintf(output, "Catalog: %s\n", report.CatalogRoot)
	fmt.Fprintf(output, "Checks: %d passed, %d warnings, %d failed\n", report.Summary.Passed, report.Summary.Warnings, report.Summary.Failed)
	for _, check := range report.Checks {
		requirement := "optional"
		if check.Required {
			requirement = "required"
		}
		fmt.Fprintf(output, "[%s] %s (%s): %s\n", strings.ToUpper(check.Status), check.ID, requirement, check.Summary)
		for _, detail := range check.Details {
			fmt.Fprintf(output, "  - %s\n", detail)
		}
	}
}

func runHistory(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 || (args[0] != "list" && args[0] != "show" && args[0] != "check") {
		return errors.New("history requires list, show, or check")
	}
	set := flag.NewFlagSet("history "+args[0], flag.ContinueOnError)
	set.SetOutput(stderr)
	stateDirectory := set.String("state-dir", "", "private scanner state directory (required)")
	format := set.String("format", "human", "human or json")
	limit := set.Int("limit", 20, "maximum runs to list")
	targetName := set.String("target-name", "", "exact target name filter")
	profileID := set.String("profile", "", "exact profile ID filter")
	terminalState := set.String("terminal-state", "", "exact terminal state filter")
	if err := set.Parse(args[1:]); err != nil {
		return err
	}
	if strings.TrimSpace(*stateDirectory) == "" {
		return errors.New("state directory is required")
	}
	if *format != "human" && *format != "json" {
		return fmt.Errorf("unsupported format %q", *format)
	}
	store, err := state.Open(context.Background(), *stateDirectory)
	if err != nil {
		return err
	}
	defer store.Close()
	if args[0] == "check" {
		if set.NArg() != 0 {
			return errors.New("history check accepts no positional arguments")
		}
		if err := store.IntegrityCheck(context.Background()); err != nil {
			return err
		}
		counts, err := store.Counts(context.Background())
		if err != nil {
			return err
		}
		report := state.CheckReport{
			SchemaVersion: state.CheckSchema, CheckedAt: time.Now().UTC(),
			StatePath: store.Path(), Integrity: "ok", Counts: counts,
		}
		if *format == "json" {
			return encodeJSON(stdout, report)
		}
		fmt.Fprintf(stdout, "State: %s\n", report.StatePath)
		fmt.Fprintln(stdout, "Integrity: ok")
		fmt.Fprintf(stdout, "Records: %d runs, %d results, %d findings, %d evidence, %d files, %d facts, %d audit events\n",
			counts.Runs, counts.Results, counts.Findings, counts.Evidence,
			counts.InventoryFiles, counts.InventoryFacts, counts.AuditEvents)
		return nil
	}
	if args[0] == "show" {
		if set.NArg() != 1 {
			return errors.New("history show requires one run ID")
		}
		run, err := store.LoadRun(context.Background(), set.Arg(0))
		if err != nil {
			return err
		}
		if *format == "json" {
			return encodeJSON(stdout, run)
		}
		printRun(stdout, run)
		return nil
	}
	if set.NArg() != 0 {
		return errors.New("history list accepts no positional arguments")
	}
	runs, err := store.ListRuns(context.Background(), state.Query{
		Limit: *limit, TargetName: *targetName, ProfileID: *profileID, TerminalState: *terminalState,
	})
	if err != nil {
		return err
	}
	report := state.HistoryReport{
		SchemaVersion: state.HistorySchema, GeneratedAt: time.Now().UTC(),
		StatePath: store.Path(), Runs: runs,
	}
	if *format == "json" {
		return encodeJSON(stdout, report)
	}
	fmt.Fprintf(stdout, "State: %s\n", report.StatePath)
	fmt.Fprintf(stdout, "Runs: %d\n", len(report.Runs))
	for _, item := range report.Runs {
		fmt.Fprintf(stdout, "- %s %s %s %s pass=%d fail=%d blocked=%d\n",
			item.CompletedAt.Format(time.RFC3339), item.RunID, item.TargetName, item.TerminalState,
			item.PassCount, item.FailCount, item.BlockedCount)
	}
	return nil
}

func runDiff(args []string, stdout, stderr io.Writer) error {
	set := flag.NewFlagSet("diff", flag.ContinueOnError)
	set.SetOutput(stderr)
	stateDirectory := set.String("state-dir", "", "private scanner state directory (required)")
	baseRunID := set.String("base-run", "", "canonical prior run ID (required)")
	target := set.String("target", ".", "current repository to inspect")
	catalogRoot := set.String("catalog-root", ".", "repository containing the PRC catalog")
	configPath := set.String("config", "", "optional validated project configuration")
	profile := set.String("profile", "", "profile to compare; defaults to configuration or base run")
	format := set.String("format", "human", "human or json")
	if err := set.Parse(args); err != nil {
		return err
	}
	if set.NArg() != 0 {
		return errors.New("diff accepts no positional arguments")
	}
	if strings.TrimSpace(*stateDirectory) == "" || strings.TrimSpace(*baseRunID) == "" {
		return errors.New("--state-dir and --base-run are required")
	}
	if *format != "human" && *format != "json" {
		return fmt.Errorf("unsupported format %q", *format)
	}
	store, err := state.Open(context.Background(), *stateDirectory)
	if err != nil {
		return err
	}
	defer store.Close()
	base, err := store.LoadRun(context.Background(), *baseRunID)
	if err != nil {
		return err
	}
	item, validation, err := configuredInventory(*target, *configPath)
	if err != nil {
		return err
	}
	if *profile == "" {
		if validation != nil {
			*profile = validation.Configuration.Assessment.Profile
		} else {
			*profile = base.Plan.ProfileID
		}
	}
	scanner, err := loadEngine(*catalogRoot)
	if err != nil {
		return err
	}
	plan, err := scanner.Plan(*profile, item)
	if err != nil {
		return err
	}
	report, err := invalidation.Analyze(base, item, plan, scanner.Catalog.Assertions, time.Now().UTC())
	if err != nil {
		return err
	}
	if *format == "json" {
		return encodeJSON(stdout, report)
	}
	printInvalidation(stdout, report)
	return nil
}

func printInvalidation(output io.Writer, report invalidation.Report) {
	fmt.Fprintf(output, "Invalidation analysis: %s -> %s\n", report.BaseInventoryDigest, report.CurrentInventoryDigest)
	fmt.Fprintf(output, "Base run: %s\n", report.BaseRunID)
	fmt.Fprintf(output, "Changed files: %d\n", len(report.ChangedFiles))
	fmt.Fprintf(output, "Assertions: %d invalidated, %d input-equivalent with rebinding required, %d reusable, %d new, %d removed\n",
		report.Summary.Invalidated, report.Summary.UnchangedInputs, report.Summary.Reusable, report.Summary.New, report.Summary.Removed)
	for _, impact := range report.Assertions {
		fmt.Fprintf(output, "- %s: %s", impact.AssertionID, impact.Conclusion)
		if impact.ReuseAllowed {
			fmt.Fprint(output, " (reuse allowed)")
		}
		fmt.Fprintln(output)
		for _, reason := range impact.Reasons {
			fmt.Fprintf(output, "  - %s: %s\n", reason.Code, reason.Detail)
		}
	}
}

func runRemediateProposal(args []string, stdout, stderr io.Writer) (int, error) {
	set, target, catalogRoot, profile := parseCommon("remediate-proposal", args, stderr)
	configPath := set.String("config", "", "optional validated project configuration")
	providerName := set.String("provider", "", "codex or claude")
	taskPath := set.String("task", "", "sealed agent task JSON file")
	outputPath := set.String("output", "", "provider output JSON file")
	candidateDirectory := set.String("candidate-dir", "", "new isolated candidate directory (required)")
	maxFiles := set.Int("max-files", 20, "maximum files the proposal may change")
	maxChangedLines := set.Int("max-changed-lines", 200, "maximum changed lines")
	format := set.String("format", "human", "human or json")
	if err := set.Parse(args); err != nil {
		return exitInternal, exitError(exitConfiguration, err)
	}
	if *providerName == "" || *taskPath == "" || *outputPath == "" || *candidateDirectory == "" {
		return exitInternal, exitError(exitConfiguration, errors.New("--provider, --task, --output, and --candidate-dir are required"))
	}
	if *format != "human" && *format != "json" {
		return exitInternal, exitError(exitConfiguration, fmt.Errorf("unsupported format %q", *format))
	}
	configuration, err := loadRemediationConfiguration(*configPath)
	if err != nil {
		return exitInternal, exitError(exitConfiguration, err)
	}
	maxAttempts := 1
	if configuration != nil {
		document := configuration.Validation.Configuration
		if !flagWasSet(set, "profile") {
			*profile = document.Assessment.Profile
		}
		if !flagWasSet(set, "max-files") {
			*maxFiles = 0
		}
		if !flagWasSet(set, "max-changed-lines") {
			*maxChangedLines = 0
		}
		maxAttempts = document.Remediation.MaxAttempts
	}
	task, err := provider.LoadTask(*taskPath)
	if err != nil {
		return exitInternal, exitError(exitConfiguration, err)
	}
	data, err := readBoundedRegularFile(*outputPath, int64(task.MaxOutputBytes))
	if err != nil {
		return exitInternal, exitError(exitConfiguration, err)
	}
	proposal, err := provider.ParseOutput(*providerName, data, task)
	if err != nil {
		return exitInternal, exitError(exitExecution, err)
	}
	candidate, err := remediation.RunProposal(remediation.ProposalOptions{
		CatalogRoot: *catalogRoot, Target: *target, CandidateDir: *candidateDirectory,
		ProfileID: *profile, Provider: *providerName, Task: task, Output: proposal,
		MaxFiles: *maxFiles, MaxChangedLines: *maxChangedLines,
		Attempt: 1, MaxAttempts: maxAttempts, Configuration: configuration,
	})
	if err != nil {
		return exitInternal, remediationCommandError(err)
	}
	if *format == "json" {
		if err := encodeJSON(stdout, candidate); err != nil {
			return exitInternal, exitError(exitInternal, err)
		}
	} else {
		printCandidate(stdout, candidate)
	}
	if !candidate.Accepted {
		return exitCandidateRejected, nil
	}
	return exitSuccess, nil
}

func loadRemediationConfiguration(path string) (*remediation.ProjectConfiguration, error) {
	if path == "" {
		return nil, nil
	}
	validation, err := projectconfig.Load(path)
	if err != nil {
		return nil, err
	}
	return &remediation.ProjectConfiguration{Validation: validation, SourcePath: path}, nil
}

func readBoundedRegularFile(path string, limit int64) ([]byte, error) {
	info, err := os.Stat(path)
	if err != nil || !info.Mode().IsRegular() || info.Size() > limit {
		return nil, fmt.Errorf("input must be a regular file no larger than %d bytes", limit)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if int64(len(data)) > limit {
		return nil, fmt.Errorf("input exceeds %d bytes", limit)
	}
	return data, nil
}

func printCandidate(output io.Writer, candidate remediation.Candidate) {
	fmt.Fprintf(output, "Candidate: %s\n", candidate.CandidateID)
	fmt.Fprintf(output, "Task: %s\n", candidate.Contract.TaskID)
	fmt.Fprintf(output, "Workspace: %s\n", candidate.CandidatePath)
	fmt.Fprintf(output, "Assertion: %s (%s -> %s)\n", candidate.Contract.AssertionID,
		candidate.BeforeAssessment, candidate.AfterAssessment)
	fmt.Fprintf(output, "Accepted: %t\n", candidate.Accepted)
	for _, change := range candidate.Changes {
		fmt.Fprintf(output, "- %s: %s (+%d/-%d lines)\n", change.Path, change.Kind, change.AddedLines, change.RemovedLines)
	}
	for _, reason := range candidate.Reasons {
		fmt.Fprintf(output, "- rejection: %s\n", reason)
	}
}

func runAdapter(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		return errors.New("adapter requires validate-output, registry-validate, plan-oci, or run-oci")
	}
	switch args[0] {
	case "validate-output":
		set := flag.NewFlagSet("adapter validate-output", flag.ContinueOnError)
		set.SetOutput(stderr)
		path := set.String("file", "-", "adapter JSONL output file, or - for stdin")
		manifestPath := set.String("manifest", "", "optional adapter manifest supplying limits")
		if err := set.Parse(args[1:]); err != nil {
			return err
		}
		limits := adapter.DefaultLimits()
		var manifest *adapter.Manifest
		if *manifestPath != "" {
			loaded, err := adapter.LoadManifest(*manifestPath)
			if err != nil {
				return err
			}
			manifest = &loaded
			limits = loaded.Resources.Limits
		}
		input := io.Reader(os.Stdin)
		if *path != "-" {
			file, err := os.Open(*path)
			if err != nil {
				return err
			}
			defer file.Close()
			input = file
		}
		transcript, err := adapter.ParseOutput(input, limits)
		if err != nil {
			return err
		}
		if manifest != nil {
			if err := adapter.ValidateTranscriptContract(*manifest, transcript); err != nil {
				return err
			}
		}
		return encodeJSON(stdout, transcript)
	case "registry-validate":
		set := flag.NewFlagSet("adapter registry-validate", flag.ContinueOnError)
		set.SetOutput(stderr)
		path := set.String("file", "", "adapter registry lockfile")
		format := set.String("format", "human", "human or json")
		if err := set.Parse(args[1:]); err != nil {
			return exitError(exitConfiguration, err)
		}
		if *path == "" {
			return exitError(exitConfiguration, errors.New("--file is required"))
		}
		if set.NArg() != 0 {
			return exitError(exitConfiguration, fmt.Errorf("unexpected registry arguments: %s", strings.Join(set.Args(), " ")))
		}
		registry, err := adapter.LoadRegistry(*path)
		if err != nil {
			return exitError(exitConfiguration, err)
		}
		if *format == "json" {
			return encodeJSON(stdout, registry.Report())
		}
		if *format != "human" {
			return exitError(exitConfiguration, fmt.Errorf("unsupported format %q", *format))
		}
		fmt.Fprintf(stdout, "Registry: %s revision %d\n", registry.ID, registry.Revision)
		fmt.Fprintf(stdout, "Digest: %s\n", registry.Digest)
		for _, entry := range registry.Entries {
			fmt.Fprintf(stdout, "- %s: %s, %s (%s)\n", entry.AdapterID, entry.Status, entry.Trust, entry.ManifestSHA256)
		}
		return nil
	case "plan-oci", "run-oci":
		return runOCIAdapter(args[0], args[1:], stdout, stderr)
	default:
		return fmt.Errorf("unknown adapter command %q", args[0])
	}
}

func runOCIAdapter(commandName string, args []string, stdout, stderr io.Writer) error {
	set := flag.NewFlagSet("adapter "+commandName, flag.ContinueOnError)
	set.SetOutput(stderr)
	manifestPath := set.String("manifest", "", "adapter manifest path")
	target := set.String("target", ".", "read-only target workspace")
	runtime := set.String("runtime", "docker", "docker or podman executable")
	if err := set.Parse(args); err != nil {
		return err
	}
	if *manifestPath == "" {
		return errors.New("--manifest is required")
	}
	manifest, err := adapter.LoadManifest(*manifestPath)
	if err != nil {
		return err
	}
	item, err := inventory.Build(*target)
	if err != nil {
		return err
	}
	runID := item.Digest
	if commandName == "run-oci" {
		runID, err = randomRunID()
		if err != nil {
			return err
		}
	}
	ociPlan, err := adapter.BuildOCIPlan(*runtime, item.Root, runID, manifest)
	if err != nil {
		return err
	}
	if commandName == "plan-oci" {
		return encodeJSON(stdout, ociPlan)
	}
	input, err := adapter.InputJSONL(runID, adapter.Subject{
		TargetName: item.TargetName, TargetCommit: item.GitCommit, InventoryDigest: item.Digest,
	}, inventoryFacts(item), map[string]any{})
	if err != nil {
		return err
	}
	result, err := adapter.RunOCI(context.Background(), ociPlan, manifest, input)
	if err != nil {
		return err
	}
	execution, err := adapter.BindExecution(runID, adapter.Subject{
		TargetName: item.TargetName, TargetCommit: item.GitCommit, InventoryDigest: item.Digest,
	}, manifest, result)
	if err != nil {
		return err
	}
	return encodeJSON(stdout, execution)
}

func randomRunID() (string, error) {
	value := make([]byte, 32)
	if _, err := rand.Read(value); err != nil {
		return "", fmt.Errorf("create adapter run ID: %w", err)
	}
	return hex.EncodeToString(value), nil
}

func parseCommon(name string, args []string, stderr io.Writer) (*flag.FlagSet, *string, *string, *string) {
	set := flag.NewFlagSet(name, flag.ContinueOnError)
	set.SetOutput(stderr)
	target := set.String("target", ".", "repository to inspect")
	catalogRoot := set.String("catalog-root", ".", "repository containing the PRC catalog")
	profile := set.String("profile", "prc/core-repository", "profile ID")
	return set, target, catalogRoot, profile
}

func encodeJSON(output io.Writer, value any) error {
	encoder := json.NewEncoder(output)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(value); err != nil {
		return exitError(exitInternal, fmt.Errorf("write JSON output: %w", err))
	}
	return nil
}

func runInventory(args []string, stdout, stderr io.Writer) error {
	set := flag.NewFlagSet("inventory", flag.ContinueOnError)
	set.SetOutput(stderr)
	target := set.String("target", ".", "repository to inspect")
	configPath := set.String("config", "", "optional validated project configuration")
	format := set.String("format", "human", "human or json")
	if err := set.Parse(args); err != nil {
		return err
	}
	item, _, err := configuredInventory(*target, *configPath)
	if err != nil {
		return err
	}
	if *format == "json" {
		return encodeJSON(stdout, item)
	}
	if *format != "human" {
		return fmt.Errorf("unsupported format %q", *format)
	}
	fmt.Fprintf(stdout, "Target: %s\n", item.TargetName)
	fmt.Fprintf(stdout, "Inventory digest: %s\n", item.Digest)
	fmt.Fprintf(stdout, "Files: %d (%d recognized source files)\n", item.FileCount, item.SourceFiles)
	fmt.Fprintf(stdout, "Package ecosystems: %s\n", displayList(item.PackageEcosystems))
	fmt.Fprintf(stdout, "Manifests: %s\n", displayList(item.Manifests))
	fmt.Fprintf(stdout, "Lock files: %s\n", displayList(item.LockFiles))
	fmt.Fprintf(stdout, "Container files: %s\n", displayList(item.ContainerFiles))
	fmt.Fprintf(stdout, "Terraform files: %s\n", displayList(item.Infrastructure.TerraformFiles))
	fmt.Fprintf(stdout, "Kubernetes files: %s\n", displayList(item.Infrastructure.KubernetesFiles))
	fmt.Fprintf(stdout, "Symlinks: %s\n", displayList(item.Symlinks))
	fmt.Fprintf(stdout, "GitHub Actions: %t\n", item.CI.GitHubActions)
	fmt.Fprintf(stdout, "Inventory graph: %d components, %d relations, %d sourced facts\n", len(item.Components), len(item.Relations), len(item.Facts))
	if item.DeclaredScope != nil {
		fmt.Fprintf(stdout, "Declared scope: %s (%s) via %s\n", item.DeclaredScope.ProjectID,
			item.DeclaredScope.RiskProfile, item.DeclaredScope.ConfigurationDigest)
	}
	return nil
}

func configuredInventory(target, configPath string) (model.Inventory, *projectconfig.Validation, error) {
	var validation *projectconfig.Validation
	if configPath != "" {
		loaded, err := projectconfig.Load(configPath)
		if err != nil {
			return model.Inventory{}, nil, err
		}
		validation = &loaded
	}
	item, err := inventory.Build(target)
	if err != nil {
		return model.Inventory{}, nil, err
	}
	if validation != nil {
		item, err = inventory.BindConfiguration(item, *validation, configPath)
		if err != nil {
			return model.Inventory{}, nil, err
		}
	}
	return item, validation, nil
}

func loadEngine(catalogRoot string) (*engine.Engine, error) {
	c, err := catalog.Load(catalogRoot)
	if err != nil {
		return nil, err
	}
	return engine.New(c), nil
}

func runPlan(args []string, stdout, stderr io.Writer) error {
	set, target, catalogRoot, profile := parseCommon("plan", args, stderr)
	configPath := set.String("config", "", "optional validated project configuration")
	mode := set.String("mode", engine.ExecutionModeInspect, "execution mode: inspect or verify-local")
	format := set.String("format", "human", "human or json")
	if err := set.Parse(args); err != nil {
		return err
	}
	item, validation, err := configuredInventory(*target, *configPath)
	if err != nil {
		return err
	}
	if validation != nil && !flagWasSet(set, "profile") {
		*profile = validation.Configuration.Assessment.Profile
	}
	scanner, err := loadEngine(*catalogRoot)
	if err != nil {
		return err
	}
	plan, err := scanner.PlanMode(*profile, item, *mode)
	if err != nil {
		return err
	}
	if *format == "json" {
		return encodeJSON(stdout, plan)
	}
	if *format != "human" {
		return fmt.Errorf("unsupported format %q", *format)
	}
	fmt.Fprintf(stdout, "Plan: %s@%s\n", plan.ProfileID, plan.ProfileVersion)
	fmt.Fprintf(stdout, "Execution mode: %s\n", plan.ExecutionMode)
	fmt.Fprintf(stdout, "Target: %s (%s)\n", plan.TargetName, plan.InventoryDigest)
	fmt.Fprintf(stdout, "Plan digest: %s\n", plan.Digest)
	if plan.ConfigurationDigest != "" {
		fmt.Fprintf(stdout, "Configuration: %s (%s)\n", plan.ConfigurationDigest, plan.ProjectID)
		fmt.Fprintf(stdout, "Environments: %s\n", displayList(plan.TargetEnvironments))
		fmt.Fprintf(stdout, "Artifacts: %s\n", displayList(plan.ArtifactDigests))
	}
	for _, assertion := range plan.Assertions {
		fmt.Fprintf(stdout, "- %s: %s via %s\n", assertion.AssertionID, assertion.Applicability, assertion.Implementation)
	}
	return nil
}

func runScan(args []string, stdout, stderr io.Writer) (int, error) {
	set, target, catalogRoot, profile := parseCommon("scan", args, stderr)
	configPath := set.String("config", "", "optional validated project configuration")
	mode := set.String("mode", engine.ExecutionModeInspect, "execution mode: inspect or verify-local")
	format := set.String("format", "human", "human, json, markdown, html, sarif, or junit")
	stateDirectory := set.String("state-dir", "", "optional directory for content-addressed evidence and run records")
	exitPolicy := set.String("exit-policy", "profile", "profile, no-go, or never")
	adapterManifest := set.String("adapter-manifest", "", "optional immutable OCI adapter manifest authorized by the selected profile")
	adapterRegistry := set.String("adapter-registry", "", "optional adapter registry lockfile")
	adapterID := set.String("adapter-id", "", "adapter ID to resolve from --adapter-registry")
	adapterRuntime := set.String("adapter-runtime", "docker", "docker or podman executable for the authorized adapter")
	if err := set.Parse(args); err != nil {
		return exitInternal, exitError(exitConfiguration, err)
	}
	if *format != "human" && *format != "json" && *format != "markdown" &&
		*format != "html" && *format != "sarif" && *format != "junit" {
		return exitInternal, exitError(exitConfiguration, fmt.Errorf("unsupported format %q", *format))
	}
	if *exitPolicy != "profile" && *exitPolicy != "no-go" && *exitPolicy != "never" {
		return exitInternal, exitError(exitConfiguration, fmt.Errorf("unsupported exit policy %q", *exitPolicy))
	}
	if *mode != engine.ExecutionModeInspect && *mode != engine.ExecutionModeVerifyLocal {
		return exitInternal, exitError(exitConfiguration, fmt.Errorf("unsupported execution mode %q", *mode))
	}
	if *adapterManifest != "" && *adapterRegistry != "" {
		return exitInternal, exitError(exitConfiguration, fmt.Errorf("--adapter-manifest and --adapter-registry are mutually exclusive"))
	}
	if (*adapterRegistry == "") != (*adapterID == "") {
		return exitInternal, exitError(exitConfiguration, fmt.Errorf("--adapter-registry and --adapter-id must be supplied together"))
	}
	item, validation, err := configuredInventory(*target, *configPath)
	if err != nil {
		return exitInternal, exitError(exitConfiguration, err)
	}
	if validation != nil && !flagWasSet(set, "profile") {
		*profile = validation.Configuration.Assessment.Profile
	}
	scanner, err := loadEngine(*catalogRoot)
	if err != nil {
		return exitInternal, exitError(exitConfiguration, err)
	}
	executions := []model.AdapterExecution{}
	if *adapterManifest != "" || *adapterRegistry != "" {
		if !flagWasSet(set, "mode") || *mode != engine.ExecutionModeVerifyLocal {
			return exitInternal, exitError(exitPolicyDenied, fmt.Errorf("adapter execution requires an explicit --mode verify-local capability grant"))
		}
		var manifest adapter.Manifest
		var registryResolution *model.AdapterResolution
		if *adapterManifest != "" {
			manifest, err = adapter.LoadManifest(*adapterManifest)
			if err != nil {
				return exitInternal, exitError(exitConfiguration, err)
			}
		} else {
			registry, registryErr := adapter.LoadRegistry(*adapterRegistry)
			if registryErr != nil {
				return exitInternal, exitError(exitConfiguration, registryErr)
			}
			resolved, resolveErr := registry.Resolve(*adapterID, "", nil, adapter.DefaultRegistryPolicy())
			if resolveErr != nil {
				return exitInternal, exitError(exitPolicyDenied, resolveErr)
			}
			manifest = resolved.Manifest
			registryResolution = &resolved.Resolution
		}
		manifestDigest, err := adapter.ManifestDigest(manifest)
		if err != nil {
			return exitInternal, exitError(exitConfiguration, err)
		}
		authorized, err := scanner.AuthorizesAdapterMode(*profile, item, *mode, manifest.ID, manifestDigest)
		if err != nil {
			return exitInternal, exitError(exitConfiguration, err)
		}
		if !authorized {
			return exitInternal, exitError(exitPolicyDenied, fmt.Errorf("adapter %s with manifest digest %s is not authorized by an applicable assertion in %s", manifest.ID, manifestDigest, *profile))
		}
		adapterRunID, err := randomRunID()
		if err != nil {
			return exitInternal, exitError(exitInternal, err)
		}
		subject := adapter.Subject{
			TargetName: item.TargetName, TargetCommit: item.GitCommit, InventoryDigest: item.Digest,
		}
		input, err := adapter.InputJSONL(adapterRunID, subject, inventoryFacts(item), map[string]any{})
		if err != nil {
			return exitInternal, exitError(exitExecution, err)
		}
		ociPlan, err := adapter.BuildOCIPlan(*adapterRuntime, item.Root, adapterRunID, manifest)
		if err != nil {
			return exitInternal, exitError(exitExecution, err)
		}
		adapterContext := context.Background()
		cancel := func() {}
		if validation != nil {
			adapterContext, cancel = context.WithTimeout(
				adapterContext,
				time.Duration(validation.Configuration.Execution.MaxDurationSeconds)*time.Second,
			)
		}
		output, err := adapter.RunOCI(adapterContext, ociPlan, manifest, input)
		cancel()
		if err != nil {
			return exitInternal, exitError(exitExecution, err)
		}
		var execution model.AdapterExecution
		if registryResolution == nil {
			execution, err = adapter.BindExecution(adapterRunID, subject, manifest, output)
		} else {
			execution, err = adapter.BindExecutionWithResolution(
				adapterRunID, subject, manifest, *registryResolution, output,
			)
		}
		if err != nil {
			return exitInternal, exitError(exitExecution, err)
		}
		executions = append(executions, execution)
	}
	run, err := scanner.ScanMode(*profile, item, executions, *mode)
	if err != nil {
		return exitInternal, exitError(exitExecution, err)
	}
	var stateStore *state.Store
	if *stateDirectory != "" {
		stateStore, err = state.Open(context.Background(), *stateDirectory)
		if err != nil {
			return exitInternal, exitError(exitConfiguration, err)
		}
		defer stateStore.Close()
	}
	if err := evidence.WriteRun(*stateDirectory, run); err != nil {
		return exitInternal, exitError(exitInternal, err)
	}
	if stateStore != nil {
		if err := stateStore.IndexRun(context.Background(), run); err != nil {
			return exitInternal, exitError(exitInternal, err)
		}
	}
	if *format == "json" {
		if err := encodeJSON(stdout, run); err != nil {
			return exitInternal, exitError(exitInternal, err)
		}
	} else if *format == "human" {
		printRun(stdout, run)
	} else if err := report.Write(*format, stdout, run); err != nil {
		return exitInternal, exitError(exitInternal, err)
	}
	switch *exitPolicy {
	case "profile", "no-go":
		return scanTerminalExitCode(run.TerminalState), nil
	case "never":
		return exitSuccess, nil
	}
	panic("validated exit policy was not handled")
}

func scanTerminalExitCode(terminalState string) int {
	switch terminalState {
	case "profile_satisfied":
		return exitSuccess
	case "no_go":
		return exitGateFailed
	case "assessment_incomplete", "environment_blocked", "machine_work_complete_manual_evidence_remaining":
		return exitIncomplete
	case "policy_stopped", "budget_exhausted":
		return exitPolicyDenied
	default:
		return exitInternal
	}
}

func inventoryFacts(item model.Inventory) map[string]any {
	return map[string]any{
		"source_files": item.SourceFiles, "package_ecosystems": item.PackageEcosystems,
		"manifests": item.Manifests, "lock_files": item.LockFiles, "container_files": item.ContainerFiles,
		"symlinks": item.Symlinks, "ci": item.CI, "infrastructure": item.Infrastructure,
		"components": item.Components, "relations": item.Relations, "facts": item.Facts,
		"declared_scope": item.DeclaredScope,
	}
}

func flagWasSet(set *flag.FlagSet, name string) bool {
	found := false
	set.Visit(func(item *flag.Flag) {
		if item.Name == name {
			found = true
		}
	})
	return found
}

func printRun(output io.Writer, run model.RunResult) {
	fmt.Fprintf(output, "Run: %s\n", run.RunID)
	fmt.Fprintf(output, "Profile: %s@%s\n", run.Plan.ProfileID, run.Plan.ProfileVersion)
	fmt.Fprintf(output, "Target: %s (%s)\n", run.Inventory.TargetName, run.Inventory.Digest)
	if run.Plan.ConfigurationDigest != "" {
		fmt.Fprintf(output, "Configuration: %s (%s)\n", run.Plan.ConfigurationDigest, run.Plan.ProjectID)
		fmt.Fprintf(output, "Environments: %s\n", displayList(run.Plan.TargetEnvironments))
		fmt.Fprintf(output, "Artifacts: %s\n", displayList(run.Plan.ArtifactDigests))
	}
	fmt.Fprintf(output, "Adapter executions: %d\n", len(run.AdapterExecutions))
	fmt.Fprintf(output, "Terminal state: %s\n\n", run.TerminalState)
	for _, result := range run.Results {
		fmt.Fprintf(output, "[%s] %s (%s/%s): %s\n", strings.ToUpper(result.Assessment),
			result.AssertionID, result.Severity, result.Gate, result.Summary)
	}
}

func runExplain(args []string, stdout, stderr io.Writer) error {
	set := flag.NewFlagSet("explain", flag.ContinueOnError)
	set.SetOutput(stderr)
	catalogRoot := set.String("catalog-root", ".", "repository containing the PRC catalog")
	if err := set.Parse(args); err != nil {
		return err
	}
	if set.NArg() != 1 {
		return errors.New("explain requires one assertion ID")
	}
	c, err := catalog.Load(*catalogRoot)
	if err != nil {
		return err
	}
	assertionID := set.Arg(0)
	assertion, ok := c.Assertions[assertionID]
	if !ok {
		return fmt.Errorf("unknown assertion %q", assertionID)
	}
	fmt.Fprintf(stdout, "%s — %s\n\n%s\n", assertion.ID, assertion.Title, assertion.Statement)
	fmt.Fprintf(stdout, "Implementation: %s\n", assertion.ImplementationID)
	fmt.Fprintf(stdout, "Applicability: %s\n", assertion.Applicability)
	fmt.Fprintf(stdout, "Severity/gate: %s/%s\n", assertion.Severity, assertion.Gate)
	fmt.Fprintf(stdout, "Remediation: %s\n", assertion.RemediationClass)
	fmt.Fprintf(stdout, "Controls: %s\n", strings.Join(assertion.ControlIDs, ", "))
	return nil
}

func displayList(values []string) string {
	if len(values) == 0 {
		return "none detected"
	}
	copyOfValues := append([]string(nil), values...)
	sort.Strings(copyOfValues)
	return strings.Join(copyOfValues, ", ")
}
