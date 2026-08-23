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
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/MarinJursic/production-readiness-checklist/scanner/adapter"
	"github.com/MarinJursic/production-readiness-checklist/scanner/adapterfixture"
	"github.com/MarinJursic/production-readiness-checklist/scanner/benchmark"
	"github.com/MarinJursic/production-readiness-checklist/scanner/catalog"
	projectconfig "github.com/MarinJursic/production-readiness-checklist/scanner/config"
	"github.com/MarinJursic/production-readiness-checklist/scanner/doctor"
	"github.com/MarinJursic/production-readiness-checklist/scanner/engine"
	"github.com/MarinJursic/production-readiness-checklist/scanner/evidence"
	"github.com/MarinJursic/production-readiness-checklist/scanner/exception"
	"github.com/MarinJursic/production-readiness-checklist/scanner/invalidation"
	"github.com/MarinJursic/production-readiness-checklist/scanner/inventory"
	"github.com/MarinJursic/production-readiness-checklist/scanner/mcpserver"
	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
	"github.com/MarinJursic/production-readiness-checklist/scanner/pack"
	"github.com/MarinJursic/production-readiness-checklist/scanner/provider"
	"github.com/MarinJursic/production-readiness-checklist/scanner/remediation"
	"github.com/MarinJursic/production-readiness-checklist/scanner/report"
	"github.com/MarinJursic/production-readiness-checklist/scanner/state"
	"github.com/MarinJursic/production-readiness-checklist/scanner/trust"
	"github.com/MarinJursic/production-readiness-checklist/scanner/verifier"
)

var (
	version  = "0.1.0-dev"
	revision = "unknown"
	builtAt  = "unknown"
)

type versionInformation struct {
	SchemaVersion string `json:"schema_version"`
	Version       string `json:"version"`
	Revision      string `json:"revision"`
	BuiltAt       string `json:"built_at"`
	GoVersion     string `json:"go_version"`
}

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
	maxScanAdapters       = 16
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
	return runWithInput(args, os.Stdin, stdout, stderr)
}

func runWithInput(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		usage(stderr)
		return exitConfiguration
	}
	var err error
	outcome := exitSuccess
	errorFallback := exitInternal
	switch args[0] {
	case "version":
		errorFallback = exitConfiguration
		err = runVersion(args[1:], stdout, stderr)
	case "inventory":
		errorFallback = exitConfiguration
		err = runInventory(args[1:], stdout, stderr)
	case "config":
		errorFallback = exitConfiguration
		err = runConfig(args[1:], stdout, stderr)
	case "catalog":
		errorFallback = exitConfiguration
		err = runCatalog(args[1:], stdout, stderr)
	case "benchmark":
		outcome, err = runBenchmark(args[1:], stdout, stderr)
	case "pack":
		errorFallback = exitConfiguration
		err = runPack(args[1:], stdout, stderr)
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
	case "exception":
		errorFallback = exitConfiguration
		err = runException(args[1:], stdout, stderr)
	case "provider":
		errorFallback = exitExecution
		err = runProvider(args[1:], stdout, stderr)
	case "mcp":
		errorFallback = exitConfiguration
		err = runMCP(args[1:], stdin, stdout, stderr)
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

func runVersion(args []string, stdout, stderr io.Writer) error {
	set := flag.NewFlagSet("version", flag.ContinueOnError)
	set.SetOutput(stderr)
	format := set.String("format", "text", "output format: text or json")
	if err := set.Parse(args); err != nil {
		return exitError(exitConfiguration, err)
	}
	if set.NArg() != 0 {
		return exitError(exitConfiguration, fmt.Errorf("unexpected version arguments: %s", strings.Join(set.Args(), " ")))
	}
	information := versionInformation{
		SchemaVersion: "prc.version/v0.1", Version: version, Revision: revision,
		BuiltAt: builtAt, GoVersion: runtime.Version(),
	}
	switch *format {
	case "text":
		fmt.Fprintf(stdout, "prc %s (revision %s, built %s, %s)\n",
			information.Version, information.Revision, information.BuiltAt, information.GoVersion)
	case "json":
		encoder := json.NewEncoder(stdout)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(information); err != nil {
			return exitError(exitInternal, fmt.Errorf("encode version information: %w", err))
		}
	default:
		return exitError(exitConfiguration, fmt.Errorf("unsupported version format %q", *format))
	}
	return nil
}

func usage(output io.Writer) {
	fmt.Fprintln(output, "Production Readiness Scanner")
	fmt.Fprintln(output, "usage: prc <catalog|pack|benchmark|config|inventory|plan|scan|diff|fix|remediate|remediate-proposal|doctor|history|explain|adapter|exception|provider|mcp|version> [options]")
}

func runMCP(args []string, stdin io.Reader, stdout, stderr io.Writer) error {
	if len(args) == 0 || args[0] != "serve" {
		return exitError(exitConfiguration, errors.New("mcp requires serve"))
	}
	set := flag.NewFlagSet("mcp serve", flag.ContinueOnError)
	set.SetOutput(stderr)
	target := set.String("target", ".", "path-locked repository to inspect")
	catalogRoot := set.String("catalog-root", ".", "path-locked repository containing the PRC catalog")
	configPath := set.String("config", "", "optional path-locked validated project configuration")
	profile := set.String("profile", "", "locked profile; defaults to configuration or prc/core-repository")
	if err := set.Parse(args[1:]); err != nil {
		return exitError(exitConfiguration, err)
	}
	if set.NArg() != 0 {
		return exitError(exitConfiguration, fmt.Errorf("unexpected mcp serve arguments: %s", strings.Join(set.Args(), " ")))
	}
	service, err := mcpserver.NewService(mcpserver.Options{
		CatalogRoot: *catalogRoot, Target: *target, ConfigPath: *configPath, ProfileID: *profile,
	})
	if err != nil {
		return exitError(exitConfiguration, err)
	}
	server, err := mcpserver.NewServer(service, version)
	if err != nil {
		return exitError(exitInternal, err)
	}
	if err := server.Serve(stdin, stdout); err != nil {
		return exitError(exitExecution, err)
	}
	return nil
}

func runException(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 || args[0] != "verify" {
		return exitError(exitConfiguration, errors.New("exception requires verify"))
	}
	set := flag.NewFlagSet("exception verify", flag.ContinueOnError)
	set.SetOutput(stderr)
	path := set.String("file", "", "risk exception YAML file")
	stateDirectory := set.String("state-dir", "", "private scanner state containing the bound run")
	trustStorePath := set.String("trust-store", "", "risk-owner trust store")
	signaturePath := set.String("signature", "", "detached risk exception signature")
	verifiedAtText := set.String("verified-at", "", "required UTC RFC3339 verification time")
	format := set.String("format", "human", "human or json")
	if err := set.Parse(args[1:]); err != nil {
		return exitError(exitConfiguration, err)
	}
	if set.NArg() != 0 || *path == "" || *stateDirectory == "" || *trustStorePath == "" ||
		*signaturePath == "" || *verifiedAtText == "" {
		return exitError(exitConfiguration, errors.New("exception verify requires --file, --state-dir, --trust-store, --signature, and --verified-at"))
	}
	if *format != "human" && *format != "json" {
		return exitError(exitConfiguration, fmt.Errorf("unsupported format %q", *format))
	}
	loaded, err := exception.Load(*path)
	if err != nil {
		return exitError(exitConfiguration, err)
	}
	store, err := state.Open(context.Background(), *stateDirectory)
	if err != nil {
		return exitError(exitConfiguration, err)
	}
	defer store.Close()
	run, err := store.LoadRun(context.Background(), loaded.Record.Run.RunID)
	if err != nil {
		return exitError(exitConfiguration, err)
	}
	trustStore, err := trust.LoadStore(*trustStorePath)
	if err != nil {
		return exitError(exitConfiguration, err)
	}
	signature, err := trust.LoadSignature(*signaturePath)
	if err != nil {
		return exitError(exitConfiguration, err)
	}
	verifiedAt, err := time.Parse(time.RFC3339Nano, *verifiedAtText)
	if err != nil {
		return exitError(exitConfiguration, fmt.Errorf("parse --verified-at: %w", err))
	}
	verification, err := exception.Verify(loaded, run, trustStore, signature, verifiedAt)
	if err != nil {
		return exitError(exitPolicyDenied, err)
	}
	if *format == "json" {
		return encodeJSON(stdout, verification)
	}
	fmt.Fprintf(stdout, "Verified risk exception: %s (%s)\n", verification.Exception.ID, verification.ExceptionDigest)
	fmt.Fprintf(stdout, "Finding: %s (%s)\n", verification.Exception.Finding.FindingID, verification.Exception.Finding.AssertionID)
	fmt.Fprintf(stdout, "Risk owner: %s (%s)\n", verification.Exception.RiskOwner.Name, verification.Exception.RiskOwner.ID)
	fmt.Fprintf(stdout, "Expires at: %s\n", verification.Exception.ExpiresAt.Format(time.RFC3339))
	fmt.Fprintf(stdout, "Gate effect: %s\n", verification.GateEffect)
	return nil
}

func runPack(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 || (args[0] != "validate" && args[0] != "verify") {
		return exitError(exitConfiguration, errors.New("pack requires validate or verify"))
	}
	action := args[0]
	set := flag.NewFlagSet("pack "+action, flag.ContinueOnError)
	set.SetOutput(stderr)
	path := set.String("file", "", "pack manifest YAML file")
	root := set.String("catalog-root", ".", "repository containing the catalog, packs, and fixtures")
	format := set.String("format", "human", "human or json")
	trustStorePath := set.String("trust-store", "", "publisher trust store for signature verification")
	signaturePath := set.String("signature", "", "detached publisher signature envelope")
	verifiedAtText := set.String("verified-at", "", "required UTC RFC3339 verification time")
	if err := set.Parse(args[1:]); err != nil {
		return exitError(exitConfiguration, err)
	}
	if *path == "" {
		return exitError(exitConfiguration, errors.New("--file is required"))
	}
	if set.NArg() != 0 {
		return exitError(exitConfiguration, fmt.Errorf("unexpected pack arguments: %s", strings.Join(set.Args(), " ")))
	}
	if *format != "human" && *format != "json" {
		return exitError(exitConfiguration, fmt.Errorf("unsupported format %q", *format))
	}
	catalogValue, err := catalog.Load(*root)
	if err != nil {
		return exitError(exitConfiguration, err)
	}
	loaded, err := pack.Load(*root, *path, catalogValue)
	if err != nil {
		return exitError(exitConfiguration, err)
	}
	if action == "verify" {
		if *trustStorePath == "" || *signaturePath == "" || *verifiedAtText == "" {
			return exitError(exitConfiguration, errors.New("pack verify requires --trust-store, --signature, and --verified-at"))
		}
		verifiedAt, err := time.Parse(time.RFC3339Nano, *verifiedAtText)
		if err != nil {
			return exitError(exitConfiguration, fmt.Errorf("parse --verified-at: %w", err))
		}
		store, err := trust.LoadStore(*trustStorePath)
		if err != nil {
			return exitError(exitConfiguration, err)
		}
		signature, err := trust.LoadSignature(*signaturePath)
		if err != nil {
			return exitError(exitConfiguration, err)
		}
		verification, err := trust.Verify(store, signature, "pack", loaded.Manifest.ID, loaded.Digest, verifiedAt)
		if err != nil {
			return exitError(exitPolicyDenied, err)
		}
		if *format == "json" {
			return encodeJSON(stdout, verification)
		}
		fmt.Fprintf(stdout, "Verified pack: %s (%s)\n", verification.ArtifactID, verification.SHA256)
		fmt.Fprintf(stdout, "Publisher key: %s\n", verification.KeyID)
		fmt.Fprintf(stdout, "Trust store: %s (%s)\n", verification.TrustStoreID, verification.TrustStoreDigest)
		fmt.Fprintf(stdout, "Verified at: %s\n", verification.VerifiedAt.Format(time.RFC3339))
		return nil
	}
	if *trustStorePath != "" || *signaturePath != "" || *verifiedAtText != "" {
		return exitError(exitConfiguration, errors.New("signature flags require pack verify"))
	}
	if *format == "json" {
		return encodeJSON(stdout, loaded.Report())
	}
	fmt.Fprintf(stdout, "Pack: %s\n", loaded.Manifest.ID)
	fmt.Fprintf(stdout, "Digest: %s\n", loaded.Digest)
	fmt.Fprintf(stdout, "Catalog: %s\n", loaded.CatalogDigest)
	fmt.Fprintf(stdout, "Benchmark: %s (%s)\n", loaded.Manifest.Benchmark.SuiteID, loaded.SuiteDigest)
	fmt.Fprintf(stdout, "Benchmark corpus: %s\n", loaded.BenchmarkCorpusDigest)
	fmt.Fprintf(stdout, "Validated assertions: %d\n", len(loaded.Manifest.Assertions))
	return nil
}

func runBenchmark(args []string, stdout, stderr io.Writer) (int, error) {
	if len(args) == 0 || args[0] != "run" {
		return exitInternal, exitError(exitConfiguration, errors.New("benchmark requires run"))
	}
	set := flag.NewFlagSet("benchmark run", flag.ContinueOnError)
	set.SetOutput(stderr)
	suitePath := set.String("suite", "", "benchmark suite YAML file")
	catalogRoot := set.String("catalog-root", ".", "repository containing the PRC catalog")
	format := set.String("format", "human", "human or json")
	evaluatedAtText := set.String("evaluated-at", "", "optional RFC3339 evaluation time for reproducible output")
	if err := set.Parse(args[1:]); err != nil {
		return exitInternal, exitError(exitConfiguration, err)
	}
	if set.NArg() != 0 {
		return exitInternal, exitError(exitConfiguration, fmt.Errorf("unexpected benchmark arguments: %s", strings.Join(set.Args(), " ")))
	}
	if *suitePath == "" {
		return exitInternal, exitError(exitConfiguration, errors.New("--suite is required"))
	}
	if *format != "human" && *format != "json" {
		return exitInternal, exitError(exitConfiguration, fmt.Errorf("unsupported format %q", *format))
	}
	evaluatedAt := time.Now().UTC()
	if *evaluatedAtText != "" {
		parsed, err := time.Parse(time.RFC3339Nano, *evaluatedAtText)
		if err != nil {
			return exitInternal, exitError(exitConfiguration, fmt.Errorf("parse --evaluated-at: %w", err))
		}
		evaluatedAt = parsed.UTC()
	}
	catalogValue, err := catalog.Load(*catalogRoot)
	if err != nil {
		return exitInternal, exitError(exitConfiguration, err)
	}
	reportValue, err := benchmark.Evaluate(catalogValue, *suitePath, evaluatedAt)
	if err != nil {
		return exitInternal, exitError(exitConfiguration, err)
	}
	if *format == "json" {
		if err := encodeJSON(stdout, reportValue); err != nil {
			return exitInternal, err
		}
	} else {
		fmt.Fprintf(stdout, "Benchmark: %s\n", reportValue.SuiteID)
		fmt.Fprintf(stdout, "Suite digest: %s\n", reportValue.SuiteDigest)
		fmt.Fprintf(stdout, "Corpus digest: %s\n", reportValue.CorpusDigest)
		fmt.Fprintf(stdout, "Cases: %d/%d deterministic\n", reportValue.Summary.DeterministicCases, reportValue.Summary.Cases)
		fmt.Fprintf(stdout, "Expectations: %d matched, %d mismatched\n", reportValue.Summary.Matched, reportValue.Summary.Mismatched)
		fmt.Fprintf(stdout, "Precision: %.4f  Recall: %.4f  False-positive rate: %.4f\n",
			reportValue.Metrics.Precision, reportValue.Metrics.Recall, reportValue.Metrics.FalsePositiveRate)
		if reportValue.Passed {
			fmt.Fprintln(stdout, "Quality budget: passed")
		} else {
			fmt.Fprintln(stdout, "Quality budget: failed")
			for _, failure := range reportValue.QualityFailures {
				fmt.Fprintln(stdout, "- "+failure)
			}
		}
	}
	if !reportValue.Passed {
		return exitGateFailed, nil
	}
	return exitSuccess, nil
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
		if failure, ok := provider.FailureFromError(err); ok {
			if encodeErr := encodeJSON(stdout, failure); encodeErr != nil {
				return encodeErr
			}
		}
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
	maxDuration := set.Int("max-duration-seconds", 1800, "maximum wall-clock duration for the full remediation loop")
	providerName := set.String("provider", "", "optional suggest-only provider: codex or claude")
	providerExecutable := set.String("provider-executable", "", "provider CLI executable; defaults to provider name")
	agentOutputSchema := set.String("agent-output-schema", "", "agent output JSON schema; defaults to catalog schema")
	allowRemoteSource := set.Bool("allow-remote-source-processing", false, "explicitly allow sealed task inputs to be processed by the provider")
	agentTimeout := set.Int("agent-timeout-seconds", 300, "provider timeout per attempt")
	agentMaxOutput := set.Int("agent-max-output-bytes", 256*1024, "provider output byte limit per attempt")
	agentMaxCost := set.Float64("agent-max-cost-usd", 0, "provider-enforced cost limit per attempt; unsupported by Codex")
	verifierRuntime := set.String("verifier-runtime", "docker", "docker or podman executable for independent candidate tests")
	verifierImage := set.String("verifier-image", "", "immutable digest-pinned verifier image (required with --provider)")
	verifierTimeout := set.Int("verifier-timeout-seconds", 300, "independent test timeout per candidate")
	verifierMemory := set.Int("verifier-memory-mb", 1024, "independent test memory limit in MiB")
	verifierCPUs := set.Float64("verifier-cpus", 1, "independent test CPU limit")
	verifierPIDs := set.Int("verifier-pids", 128, "independent test process limit")
	verifierTmpfs := set.Int("verifier-tmpfs-mb", 512, "independent test scratch limit in MiB")
	verifierMaxOutput := set.Int("verifier-max-output-bytes", 1024*1024, "independent test stdout and stderr limit")
	format := set.String("format", "human", "human or json")
	if err := set.Parse(args); err != nil {
		return exitInternal, exitError(exitConfiguration, err)
	}
	if *format != "human" && *format != "json" {
		return exitInternal, exitError(exitConfiguration, fmt.Errorf("unsupported format %q", *format))
	}
	if *providerName == "" {
		for _, name := range []string{"provider-executable", "agent-output-schema", "allow-remote-source-processing", "agent-timeout-seconds", "agent-max-output-bytes", "agent-max-cost-usd", "verifier-runtime", "verifier-image", "verifier-timeout-seconds", "verifier-memory-mb", "verifier-cpus", "verifier-pids", "verifier-tmpfs-mb", "verifier-max-output-bytes"} {
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
		if !flagWasSet(set, "max-duration-seconds") {
			*maxDuration = document.Execution.MaxDurationSeconds
		}
	}
	var agent *remediation.AgentOptions
	var verification *verifier.Options
	if *providerName != "" {
		if strings.TrimSpace(*verifierImage) == "" {
			return exitInternal, exitError(exitConfiguration, errors.New("--verifier-image is required with --provider"))
		}
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
		configured := verifier.Defaults(*verifierRuntime, *verifierImage, "go")
		configured.TimeoutSeconds = *verifierTimeout
		configured.MemoryMB = *verifierMemory
		configured.CPUs = *verifierCPUs
		configured.PIDs = *verifierPIDs
		configured.TmpfsMB = *verifierTmpfs
		configured.MaxStdoutBytes = *verifierMaxOutput
		configured.MaxStderrBytes = *verifierMaxOutput
		if err := verifier.ValidateOptions(configured); err != nil {
			return exitInternal, exitError(exitConfiguration, err)
		}
		configured.Kind = ""
		verification = &configured
	}
	result, err := remediation.RunLoop(remediation.LoopOptions{
		CatalogRoot: *catalogRoot, Target: *target, CandidateRoot: *candidateRoot,
		ProfileID: *profile, MaxAttempts: *maxAttempts, MaxFiles: *maxFiles,
		MaxChangedLines: *maxChangedLines, MaxDurationSeconds: *maxDuration,
		Configuration: configuration, Agent: agent, Verifier: verification,
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
	case "provider_failed":
		if len(result.Attempts) > 0 && result.Attempts[len(result.Attempts)-1].ReasonCode == "provider_cancelled" {
			return exitCancelled
		}
		return exitExecution
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
	fmt.Fprintf(output, "Duration ceiling: %d seconds\n", result.MaxDurationSeconds)
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
	for _, failure := range result.ProviderFailures {
		fmt.Fprintf(output, "- provider %s failure %s: %s\n", failure.Provider, failure.FailureID, failure.ReasonCode)
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
	value = strings.TrimSpace(value)
	if value == "" {
		return fmt.Errorf("repeated flag value cannot be empty")
	}
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
	verifierRuntime := set.String("verifier-runtime", "docker", "docker or podman executable for independent candidate tests")
	verifierImage := set.String("verifier-image", "", "immutable digest-pinned verifier image (required)")
	verifierTimeout := set.Int("verifier-timeout-seconds", 300, "independent test timeout")
	verifierMemory := set.Int("verifier-memory-mb", 1024, "independent test memory limit in MiB")
	format := set.String("format", "human", "human or json")
	if err := set.Parse(args); err != nil {
		return exitInternal, exitError(exitConfiguration, err)
	}
	if *providerName == "" || *taskPath == "" || *outputPath == "" || *candidateDirectory == "" || *verifierImage == "" {
		return exitInternal, exitError(exitConfiguration, errors.New("--provider, --task, --output, --candidate-dir, and --verifier-image are required"))
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
	verification := verifier.Defaults(*verifierRuntime, *verifierImage, "go")
	verification.TimeoutSeconds = *verifierTimeout
	verification.MemoryMB = *verifierMemory
	if err := verifier.ValidateOptions(verification); err != nil {
		return exitInternal, exitError(exitConfiguration, err)
	}
	verification.Kind = ""
	candidate, err := remediation.RunProposal(remediation.ProposalOptions{
		CatalogRoot: *catalogRoot, Target: *target, CandidateDir: *candidateDirectory,
		ProfileID: *profile, Provider: *providerName, Task: task, Output: proposal,
		MaxFiles: *maxFiles, MaxChangedLines: *maxChangedLines,
		Attempt: 1, MaxAttempts: maxAttempts, Configuration: configuration, Verifier: &verification,
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
	// The path is an explicit operator-selected CLI input, not a target-derived
	// relative path. Its type and size are validated before it is read.
	// #nosec G703 -- arbitrary explicit input paths are the command contract.
	info, err := os.Stat(path)
	if err != nil || !info.Mode().IsRegular() || info.Size() > limit {
		return nil, fmt.Errorf("input must be a regular file no larger than %d bytes", limit)
	}
	// #nosec G304 G703 -- see the explicit-input boundary above.
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
	if candidate.Verification != nil {
		fmt.Fprintf(output, "- verification %s: %s (%s)\n", candidate.Verification.ExecutionID,
			candidate.Verification.Outcome, strings.Join(candidate.Verification.Command, " "))
	}
	for _, reason := range candidate.Reasons {
		fmt.Fprintf(output, "- rejection: %s\n", reason)
	}
}

func runAdapter(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		return errors.New("adapter requires validate-output, fixture-validate, registry-validate, registry-verify, plan-oci, or run-oci")
	}
	switch args[0] {
	case "fixture-validate":
		set := flag.NewFlagSet("adapter fixture-validate", flag.ContinueOnError)
		set.SetOutput(stderr)
		suitePath := set.String("suite", "", "adapter fixture suite YAML file")
		format := set.String("format", "human", "human or json")
		if err := set.Parse(args[1:]); err != nil {
			return exitError(exitConfiguration, err)
		}
		if set.NArg() != 0 || *suitePath == "" {
			return exitError(exitConfiguration, errors.New("adapter fixture-validate requires --suite"))
		}
		if *format != "human" && *format != "json" {
			return exitError(exitConfiguration, fmt.Errorf("unsupported format %q", *format))
		}
		reportValue, err := adapterfixture.Evaluate(*suitePath)
		if err != nil {
			return exitError(exitConfiguration, err)
		}
		if *format == "json" {
			if err := encodeJSON(stdout, reportValue); err != nil {
				return err
			}
		} else {
			fmt.Fprintf(stdout, "Adapter fixture suite: %s\n", reportValue.SuiteID)
			fmt.Fprintf(stdout, "Suite digest: %s\n", reportValue.SuiteDigest)
			fmt.Fprintf(stdout, "Corpus digest: %s\n", reportValue.CorpusDigest)
			fmt.Fprintf(stdout, "Adapter: %s (%s)\n", reportValue.AdapterID, reportValue.ManifestSHA256)
			fmt.Fprintf(stdout, "Cases: %d matched, %d mismatched, %d/%d deterministic\n",
				reportValue.Summary.Matched, reportValue.Summary.Mismatched,
				reportValue.Summary.DeterministicCases, reportValue.Summary.Cases)
			if reportValue.Passed {
				fmt.Fprintln(stdout, "Fixture quality gate: passed")
			} else {
				fmt.Fprintln(stdout, "Fixture quality gate: failed")
				for _, failure := range reportValue.QualityFailures {
					fmt.Fprintln(stdout, "- "+failure)
				}
			}
		}
		if !reportValue.Passed {
			return exitError(exitGateFailed, errors.New("adapter fixture quality gate failed"))
		}
		return nil
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
		var transcript adapter.Transcript
		var err error
		if manifest == nil {
			transcript, err = adapter.ParseOutput(input, limits)
		} else {
			transcript, err = adapter.ParseManifestOutput(*manifest, input)
		}
		if err != nil {
			return err
		}
		if manifest != nil {
			if err := adapter.ValidateTranscriptContract(*manifest, transcript); err != nil {
				return err
			}
		}
		return encodeJSON(stdout, transcript)
	case "registry-validate", "registry-verify":
		action := args[0]
		set := flag.NewFlagSet("adapter "+action, flag.ContinueOnError)
		set.SetOutput(stderr)
		path := set.String("file", "", "adapter registry lockfile")
		format := set.String("format", "human", "human or json")
		trustStorePath := set.String("trust-store", "", "publisher trust store for signature verification")
		signaturePath := set.String("signature", "", "detached publisher signature envelope")
		verifiedAtText := set.String("verified-at", "", "required UTC RFC3339 verification time")
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
		if action == "registry-verify" {
			if *trustStorePath == "" || *signaturePath == "" || *verifiedAtText == "" {
				return exitError(exitConfiguration, errors.New("adapter registry-verify requires --trust-store, --signature, and --verified-at"))
			}
			verifiedAt, err := time.Parse(time.RFC3339Nano, *verifiedAtText)
			if err != nil {
				return exitError(exitConfiguration, fmt.Errorf("parse --verified-at: %w", err))
			}
			store, err := trust.LoadStore(*trustStorePath)
			if err != nil {
				return exitError(exitConfiguration, err)
			}
			signature, err := trust.LoadSignature(*signaturePath)
			if err != nil {
				return exitError(exitConfiguration, err)
			}
			verification, err := trust.Verify(store, signature, "adapter-registry", registry.ID, registry.Digest, verifiedAt)
			if err != nil {
				return exitError(exitPolicyDenied, err)
			}
			if *format == "json" {
				return encodeJSON(stdout, verification)
			}
			if *format != "human" {
				return exitError(exitConfiguration, fmt.Errorf("unsupported format %q", *format))
			}
			fmt.Fprintf(stdout, "Verified adapter registry: %s (%s)\n", verification.ArtifactID, verification.SHA256)
			fmt.Fprintf(stdout, "Publisher key: %s\n", verification.KeyID)
			fmt.Fprintf(stdout, "Trust store: %s (%s)\n", verification.TrustStoreID, verification.TrustStoreDigest)
			fmt.Fprintf(stdout, "Verified at: %s\n", verification.VerifiedAt.Format(time.RFC3339))
			return nil
		}
		if *trustStorePath != "" || *signaturePath != "" || *verifiedAtText != "" {
			return exitError(exitConfiguration, errors.New("signature flags require adapter registry-verify"))
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
	if commandName == "plan-oci" {
		ociPlan, err := adapter.BuildOCIPlan(*runtime, item.Root, runID, manifest)
		if err != nil {
			return err
		}
		return encodeJSON(stdout, ociPlan)
	}
	snapshot, err := adapter.PrepareSnapshotForManifest(item, manifest)
	if err != nil {
		return err
	}
	defer snapshot.Close()
	ociPlan, err := adapter.BuildSnapshotOCIPlan(*runtime, snapshot, runID, manifest)
	if err != nil {
		return err
	}
	input, err := adapter.ExecutionInput(manifest, runID, adapter.Subject{
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

type resolvedScanAdapter struct {
	Manifest       adapter.Manifest
	ManifestDigest string
	Resolution     *model.AdapterResolution
}

func validateUniqueFlagValues(values repeatedStringFlag, label string) error {
	seen := map[string]bool{}
	for _, value := range values {
		if seen[value] {
			return fmt.Errorf("%s cannot repeat %q", label, value)
		}
		seen[value] = true
	}
	return nil
}

func resolveScanAdapters(
	scanner *engine.Engine,
	profile string,
	item model.Inventory,
	mode string,
	manifestPaths repeatedStringFlag,
	registryPath string,
	adapterIDs repeatedStringFlag,
) ([]resolvedScanAdapter, error) {
	resolved := make([]resolvedScanAdapter, 0, max(len(manifestPaths), len(adapterIDs)))
	if registryPath == "" {
		for _, manifestPath := range manifestPaths {
			manifest, err := adapter.LoadManifest(manifestPath)
			if err != nil {
				return nil, exitError(exitConfiguration, err)
			}
			resolved = append(resolved, resolvedScanAdapter{Manifest: manifest})
		}
	} else {
		registry, err := adapter.LoadRegistry(registryPath)
		if err != nil {
			return nil, exitError(exitConfiguration, err)
		}
		for _, adapterID := range adapterIDs {
			entry, err := registry.Resolve(adapterID, "", nil, adapter.DefaultRegistryPolicy())
			if err != nil {
				return nil, exitError(exitPolicyDenied, err)
			}
			resolution := entry.Resolution
			resolved = append(resolved, resolvedScanAdapter{
				Manifest: entry.Manifest, Resolution: &resolution,
			})
		}
	}

	seen := map[string]bool{}
	for index := range resolved {
		digest, err := adapter.ManifestDigest(resolved[index].Manifest)
		if err != nil {
			return nil, exitError(exitConfiguration, err)
		}
		resolved[index].ManifestDigest = digest
		key := resolved[index].Manifest.ID + "\x00" + digest
		if seen[key] {
			return nil, exitError(exitConfiguration, fmt.Errorf(
				"adapter %s with manifest digest %s was requested more than once",
				resolved[index].Manifest.ID, digest,
			))
		}
		seen[key] = true
		authorized, err := scanner.AuthorizesAdapterMode(
			profile, item, mode, resolved[index].Manifest.ID, digest,
		)
		if err != nil {
			return nil, exitError(exitConfiguration, err)
		}
		if !authorized {
			return nil, exitError(exitPolicyDenied, fmt.Errorf(
				"adapter %s with manifest digest %s is not authorized by an applicable assertion in %s",
				resolved[index].Manifest.ID, digest, profile,
			))
		}
	}
	sort.Slice(resolved, func(left, right int) bool {
		if resolved[left].Manifest.ID != resolved[right].Manifest.ID {
			return resolved[left].Manifest.ID < resolved[right].Manifest.ID
		}
		return resolved[left].ManifestDigest < resolved[right].ManifestDigest
	})
	return resolved, nil
}

func executeScanAdapter(
	ctx context.Context,
	item model.Inventory,
	resolved resolvedScanAdapter,
	runtime string,
) (model.AdapterExecution, error) {
	if err := ctx.Err(); err != nil {
		return model.AdapterExecution{}, fmt.Errorf("adapter execution budget is exhausted: %w", err)
	}
	runID, err := randomRunID()
	if err != nil {
		return model.AdapterExecution{}, err
	}
	subject := adapter.Subject{
		TargetName: item.TargetName, TargetCommit: item.GitCommit, InventoryDigest: item.Digest,
	}
	input, err := adapter.ExecutionInput(
		resolved.Manifest, runID, subject, inventoryFacts(item), map[string]any{},
	)
	if err != nil {
		return model.AdapterExecution{}, err
	}
	snapshot, err := adapter.PrepareSnapshotForManifest(item, resolved.Manifest)
	if err != nil {
		return model.AdapterExecution{}, err
	}
	defer snapshot.Close()
	plan, err := adapter.BuildSnapshotOCIPlan(runtime, snapshot, runID, resolved.Manifest)
	if err != nil {
		return model.AdapterExecution{}, err
	}
	output, err := adapter.RunOCI(ctx, plan, resolved.Manifest, input)
	if err != nil {
		return model.AdapterExecution{}, err
	}
	if resolved.Resolution == nil {
		return adapter.BindExecution(runID, subject, resolved.Manifest, output)
	}
	return adapter.BindExecutionWithResolution(
		runID, subject, resolved.Manifest, *resolved.Resolution, output,
	)
}

func runScan(args []string, stdout, stderr io.Writer) (int, error) {
	set, target, catalogRoot, profile := parseCommon("scan", args, stderr)
	configPath := set.String("config", "", "optional validated project configuration")
	mode := set.String("mode", engine.ExecutionModeInspect, "execution mode: inspect or verify-local")
	format := set.String("format", "human", "human, json, markdown, html, sarif, or junit")
	stateDirectory := set.String("state-dir", "", "optional directory for content-addressed evidence and run records")
	exitPolicy := set.String("exit-policy", "profile", "profile, no-go, or never")
	adapterManifests := repeatedStringFlag{}
	set.Var(&adapterManifests, "adapter-manifest", "immutable OCI adapter manifest authorized by the selected profile; repeatable")
	adapterRegistry := set.String("adapter-registry", "", "optional adapter registry lockfile")
	adapterIDs := repeatedStringFlag{}
	set.Var(&adapterIDs, "adapter-id", "adapter ID to resolve from --adapter-registry; repeatable")
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
	if len(adapterManifests) > 0 && *adapterRegistry != "" {
		return exitInternal, exitError(exitConfiguration, fmt.Errorf("--adapter-manifest and --adapter-registry are mutually exclusive"))
	}
	if (*adapterRegistry == "") != (len(adapterIDs) == 0) {
		return exitInternal, exitError(exitConfiguration, fmt.Errorf("--adapter-registry and at least one --adapter-id must be supplied together"))
	}
	requestedAdapters := len(adapterManifests)
	if *adapterRegistry != "" {
		requestedAdapters = len(adapterIDs)
	}
	if requestedAdapters > maxScanAdapters {
		return exitInternal, exitError(exitConfiguration, fmt.Errorf("scan accepts at most %d adapters", maxScanAdapters))
	}
	if err := validateUniqueFlagValues(adapterManifests, "--adapter-manifest"); err != nil {
		return exitInternal, exitError(exitConfiguration, err)
	}
	if err := validateUniqueFlagValues(adapterIDs, "--adapter-id"); err != nil {
		return exitInternal, exitError(exitConfiguration, err)
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
	if requestedAdapters > 0 {
		if !flagWasSet(set, "mode") || *mode != engine.ExecutionModeVerifyLocal {
			return exitInternal, exitError(exitPolicyDenied, fmt.Errorf("adapter execution requires an explicit --mode verify-local capability grant"))
		}
		resolvedAdapters, resolveErr := resolveScanAdapters(
			scanner, *profile, item, *mode, adapterManifests, *adapterRegistry, adapterIDs,
		)
		if resolveErr != nil {
			return exitInternal, resolveErr
		}
		adapterContext := context.Background()
		cancelAdapters := func() {}
		if validation != nil {
			adapterContext, cancelAdapters = context.WithTimeout(
				adapterContext,
				time.Duration(validation.Configuration.Execution.MaxDurationSeconds)*time.Second,
			)
		}
		defer cancelAdapters()
		for _, resolved := range resolvedAdapters {
			execution, executeErr := executeScanAdapter(adapterContext, item, resolved, *adapterRuntime)
			if executeErr != nil {
				return exitInternal, exitError(exitExecution, executeErr)
			}
			executions = append(executions, execution)
		}
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
