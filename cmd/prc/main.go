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
	"sort"
	"strings"

	"github.com/MarinJursic/production-readiness-checklist/scanner/adapter"
	"github.com/MarinJursic/production-readiness-checklist/scanner/catalog"
	projectconfig "github.com/MarinJursic/production-readiness-checklist/scanner/config"
	"github.com/MarinJursic/production-readiness-checklist/scanner/engine"
	"github.com/MarinJursic/production-readiness-checklist/scanner/evidence"
	"github.com/MarinJursic/production-readiness-checklist/scanner/inventory"
	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
	"github.com/MarinJursic/production-readiness-checklist/scanner/provider"
	"github.com/MarinJursic/production-readiness-checklist/scanner/remediation"
	"github.com/MarinJursic/production-readiness-checklist/scanner/report"
)

const version = "0.1.0-dev"

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		usage(stderr)
		return 2
	}
	var err error
	var successful = true
	var resultExit bool
	switch args[0] {
	case "version":
		fmt.Fprintf(stdout, "prc %s\n", version)
		return 0
	case "inventory":
		err = runInventory(args[1:], stdout, stderr)
	case "config":
		err = runConfig(args[1:], stdout, stderr)
	case "plan":
		err = runPlan(args[1:], stdout, stderr)
	case "scan":
		resultExit = true
		successful, err = runScan(args[1:], stdout, stderr)
	case "remediate":
		resultExit = true
		successful, err = runRemediate(args[1:], stdout, stderr)
	case "remediate-proposal":
		resultExit = true
		successful, err = runRemediateProposal(args[1:], stdout, stderr)
	case "explain":
		err = runExplain(args[1:], stdout, stderr)
	case "adapter":
		err = runAdapter(args[1:], stdout, stderr)
	case "provider":
		err = runProvider(args[1:], stdout, stderr)
	case "help", "-h", "--help":
		usage(stdout)
		return 0
	default:
		fmt.Fprintf(stderr, "unknown command %q\n", args[0])
		usage(stderr)
		return 2
	}
	if err != nil {
		fmt.Fprintln(stderr, "error:", err)
		return 2
	}
	if resultExit && !successful {
		return 1
	}
	return 0
}

func usage(output io.Writer) {
	fmt.Fprintln(output, "Production Readiness Scanner")
	fmt.Fprintln(output, "usage: prc <config|inventory|plan|scan|remediate|remediate-proposal|explain|adapter|provider|version> [options]")
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
		if err := set.Parse(args[1:]); err != nil {
			return err
		}
		if *path == "" {
			return errors.New("--file is required")
		}
		task, err := provider.SealTask(*path, *workspace)
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

func runRemediate(args []string, stdout, stderr io.Writer) (bool, error) {
	set, target, catalogRoot, profile := parseCommon("remediate", args, stderr)
	assertionID := set.String("assertion", "PRC-A-CORE-014", "registered R1 assertion ID to remediate")
	candidateDirectory := set.String("candidate-dir", "", "new isolated candidate directory (required)")
	maxFiles := set.Int("max-files", 20, "maximum files the fix may change")
	maxChangedLines := set.Int("max-changed-lines", 20, "maximum changed lines")
	format := set.String("format", "human", "human or json")
	if err := set.Parse(args); err != nil {
		return false, err
	}
	if *format != "human" && *format != "json" {
		return false, fmt.Errorf("unsupported format %q", *format)
	}
	candidate, err := remediation.Run(remediation.Options{
		CatalogRoot: *catalogRoot, Target: *target, CandidateDir: *candidateDirectory,
		ProfileID: *profile, AssertionID: *assertionID,
		MaxFiles: *maxFiles, MaxChangedLines: *maxChangedLines,
	})
	if err != nil {
		return false, err
	}
	if *format == "json" {
		if err := encodeJSON(stdout, candidate); err != nil {
			return false, err
		}
	} else {
		printCandidate(stdout, candidate)
	}
	return candidate.Accepted, nil
}

func runRemediateProposal(args []string, stdout, stderr io.Writer) (bool, error) {
	set, target, catalogRoot, profile := parseCommon("remediate-proposal", args, stderr)
	providerName := set.String("provider", "", "codex or claude")
	taskPath := set.String("task", "", "sealed agent task JSON file")
	outputPath := set.String("output", "", "provider output JSON file")
	candidateDirectory := set.String("candidate-dir", "", "new isolated candidate directory (required)")
	maxFiles := set.Int("max-files", 20, "maximum files the proposal may change")
	maxChangedLines := set.Int("max-changed-lines", 200, "maximum changed lines")
	format := set.String("format", "human", "human or json")
	if err := set.Parse(args); err != nil {
		return false, err
	}
	if *providerName == "" || *taskPath == "" || *outputPath == "" || *candidateDirectory == "" {
		return false, errors.New("--provider, --task, --output, and --candidate-dir are required")
	}
	if *format != "human" && *format != "json" {
		return false, fmt.Errorf("unsupported format %q", *format)
	}
	task, err := provider.LoadTask(*taskPath)
	if err != nil {
		return false, err
	}
	data, err := readBoundedRegularFile(*outputPath, int64(task.MaxOutputBytes))
	if err != nil {
		return false, err
	}
	proposal, err := provider.ParseOutput(*providerName, data, task)
	if err != nil {
		return false, err
	}
	candidate, err := remediation.RunProposal(remediation.ProposalOptions{
		CatalogRoot: *catalogRoot, Target: *target, CandidateDir: *candidateDirectory,
		ProfileID: *profile, Provider: *providerName, Task: task, Output: proposal,
		MaxFiles: *maxFiles, MaxChangedLines: *maxChangedLines,
	})
	if err != nil {
		return false, err
	}
	if *format == "json" {
		if err := encodeJSON(stdout, candidate); err != nil {
			return false, err
		}
	} else {
		printCandidate(stdout, candidate)
	}
	return candidate.Accepted, nil
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
		return errors.New("adapter requires validate-output, plan-oci, or run-oci")
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
		if *manifestPath != "" {
			manifest, err := adapter.LoadManifest(*manifestPath)
			if err != nil {
				return err
			}
			limits = manifest.Resources.Limits
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
		return encodeJSON(stdout, transcript)
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
	return encoder.Encode(value)
}

func runInventory(args []string, stdout, stderr io.Writer) error {
	set := flag.NewFlagSet("inventory", flag.ContinueOnError)
	set.SetOutput(stderr)
	target := set.String("target", ".", "repository to inspect")
	format := set.String("format", "human", "human or json")
	if err := set.Parse(args); err != nil {
		return err
	}
	item, err := inventory.Build(*target)
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
	return nil
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
	format := set.String("format", "human", "human or json")
	if err := set.Parse(args); err != nil {
		return err
	}
	item, err := inventory.Build(*target)
	if err != nil {
		return err
	}
	scanner, err := loadEngine(*catalogRoot)
	if err != nil {
		return err
	}
	plan, err := scanner.Plan(*profile, item)
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
	fmt.Fprintf(stdout, "Target: %s (%s)\n", plan.TargetName, plan.InventoryDigest)
	fmt.Fprintf(stdout, "Plan digest: %s\n", plan.Digest)
	for _, assertion := range plan.Assertions {
		fmt.Fprintf(stdout, "- %s: %s via %s\n", assertion.AssertionID, assertion.Applicability, assertion.Implementation)
	}
	return nil
}

func runScan(args []string, stdout, stderr io.Writer) (bool, error) {
	set, target, catalogRoot, profile := parseCommon("scan", args, stderr)
	format := set.String("format", "human", "human, json, markdown, html, sarif, or junit")
	stateDirectory := set.String("state-dir", "", "optional directory for content-addressed evidence and run records")
	exitPolicy := set.String("exit-policy", "profile", "profile, no-go, or never")
	adapterManifest := set.String("adapter-manifest", "", "optional immutable OCI adapter manifest authorized by the selected profile")
	adapterRuntime := set.String("adapter-runtime", "docker", "docker or podman executable for the authorized adapter")
	if err := set.Parse(args); err != nil {
		return false, err
	}
	if *format != "human" && *format != "json" && *format != "markdown" &&
		*format != "html" && *format != "sarif" && *format != "junit" {
		return false, fmt.Errorf("unsupported format %q", *format)
	}
	if *exitPolicy != "profile" && *exitPolicy != "no-go" && *exitPolicy != "never" {
		return false, fmt.Errorf("unsupported exit policy %q", *exitPolicy)
	}
	item, err := inventory.Build(*target)
	if err != nil {
		return false, err
	}
	scanner, err := loadEngine(*catalogRoot)
	if err != nil {
		return false, err
	}
	executions := []model.AdapterExecution{}
	if *adapterManifest != "" {
		manifest, err := adapter.LoadManifest(*adapterManifest)
		if err != nil {
			return false, err
		}
		manifestDigest, err := adapter.ManifestDigest(manifest)
		if err != nil {
			return false, err
		}
		authorized, err := scanner.AuthorizesAdapter(*profile, item, manifest.ID, manifestDigest)
		if err != nil {
			return false, err
		}
		if !authorized {
			return false, fmt.Errorf("adapter %s with manifest digest %s is not authorized by an applicable assertion in %s", manifest.ID, manifestDigest, *profile)
		}
		adapterRunID, err := randomRunID()
		if err != nil {
			return false, err
		}
		subject := adapter.Subject{
			TargetName: item.TargetName, TargetCommit: item.GitCommit, InventoryDigest: item.Digest,
		}
		input, err := adapter.InputJSONL(adapterRunID, subject, inventoryFacts(item), map[string]any{})
		if err != nil {
			return false, err
		}
		ociPlan, err := adapter.BuildOCIPlan(*adapterRuntime, item.Root, adapterRunID, manifest)
		if err != nil {
			return false, err
		}
		output, err := adapter.RunOCI(context.Background(), ociPlan, manifest, input)
		if err != nil {
			return false, err
		}
		execution, err := adapter.BindExecution(adapterRunID, subject, manifest, output)
		if err != nil {
			return false, err
		}
		executions = append(executions, execution)
	}
	run, err := scanner.ScanWithAdapterEvidence(*profile, item, executions)
	if err != nil {
		return false, err
	}
	if err := evidence.WriteRun(*stateDirectory, run); err != nil {
		return false, err
	}
	if *format == "json" {
		if err := encodeJSON(stdout, run); err != nil {
			return false, err
		}
	} else if *format == "human" {
		printRun(stdout, run)
	} else if err := report.Write(*format, stdout, run); err != nil {
		return false, err
	}
	switch *exitPolicy {
	case "profile":
		return run.TerminalState == "profile_satisfied", nil
	case "no-go":
		return run.TerminalState != "no_go", nil
	case "never":
		return true, nil
	}
	panic("validated exit policy was not handled")
}

func inventoryFacts(item model.Inventory) map[string]any {
	return map[string]any{
		"source_files": item.SourceFiles, "package_ecosystems": item.PackageEcosystems,
		"manifests": item.Manifests, "lock_files": item.LockFiles, "container_files": item.ContainerFiles,
		"symlinks": item.Symlinks, "ci": item.CI, "infrastructure": item.Infrastructure,
		"components": item.Components, "relations": item.Relations, "facts": item.Facts,
	}
}

func printRun(output io.Writer, run model.RunResult) {
	fmt.Fprintf(output, "Run: %s\n", run.RunID)
	fmt.Fprintf(output, "Profile: %s@%s\n", run.Plan.ProfileID, run.Plan.ProfileVersion)
	fmt.Fprintf(output, "Target: %s (%s)\n", run.Inventory.TargetName, run.Inventory.Digest)
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
