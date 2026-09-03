package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/MarinJursic/production-readiness-checklist/scanner/adapter"
	"github.com/MarinJursic/production-readiness-checklist/scanner/adapterfixture"
	"github.com/MarinJursic/production-readiness-checklist/scanner/automaticcoverage"
	"github.com/MarinJursic/production-readiness-checklist/scanner/benchmark"
	"github.com/MarinJursic/production-readiness-checklist/scanner/catalog"
	projectconfig "github.com/MarinJursic/production-readiness-checklist/scanner/config"
	"github.com/MarinJursic/production-readiness-checklist/scanner/controlprogram"
	"github.com/MarinJursic/production-readiness-checklist/scanner/controlprogramcatalog"
	"github.com/MarinJursic/production-readiness-checklist/scanner/controlreview"
	"github.com/MarinJursic/production-readiness-checklist/scanner/controlruntime"
	"github.com/MarinJursic/production-readiness-checklist/scanner/doctor"
	"github.com/MarinJursic/production-readiness-checklist/scanner/engine"
	"github.com/MarinJursic/production-readiness-checklist/scanner/evidence"
	"github.com/MarinJursic/production-readiness-checklist/scanner/evidencebundle"
	"github.com/MarinJursic/production-readiness-checklist/scanner/evidencerequirements"
	"github.com/MarinJursic/production-readiness-checklist/scanner/evidenceset"
	"github.com/MarinJursic/production-readiness-checklist/scanner/exception"
	"github.com/MarinJursic/production-readiness-checklist/scanner/fullscan"
	"github.com/MarinJursic/production-readiness-checklist/scanner/invalidation"
	"github.com/MarinJursic/production-readiness-checklist/scanner/inventory"
	"github.com/MarinJursic/production-readiness-checklist/scanner/mcpserver"
	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
	"github.com/MarinJursic/production-readiness-checklist/scanner/pack"
	"github.com/MarinJursic/production-readiness-checklist/scanner/provider"
	"github.com/MarinJursic/production-readiness-checklist/scanner/remediation"
	"github.com/MarinJursic/production-readiness-checklist/scanner/report"
	"github.com/MarinJursic/production-readiness-checklist/scanner/repositoryevidence"
	"github.com/MarinJursic/production-readiness-checklist/scanner/state"
	"github.com/MarinJursic/production-readiness-checklist/scanner/trust"
	"github.com/MarinJursic/production-readiness-checklist/scanner/verifier"
)

var (
	version            = "0.1.0-dev"
	revision           = "unknown"
	builtAt            = "unknown"
	userCacheDirectory = os.UserCacheDir
	commandContext     = context.Background()
)

type versionInformation struct {
	SchemaVersion string `json:"schema_version"`
	Version       string `json:"version"`
	Revision      string `json:"revision"`
	BuiltAt       string `json:"built_at"`
	GoVersion     string `json:"go_version"`
}

const (
	exitSuccess            = 0
	exitGateFailed         = 1
	exitIncomplete         = 2
	exitConfiguration      = 3
	exitExecution          = 4
	exitPolicyDenied       = 5
	exitInternal           = 6
	exitCancelled          = 7
	exitCandidateRejected  = 8
	maxScanAdapters        = 16
	defaultReportRetention = 5
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
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	commandContext = ctx
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	return runWithInput(args, os.Stdin, stdout, stderr)
}

func runWithInput(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	args = friendlyTopLevelArgs(args)
	if args[0] == "--version" || args[0] == "-v" {
		args = append([]string{"version"}, args[1:]...)
	}
	var err error
	outcome := exitSuccess
	errorFallback := exitInternal
	switch args[0] {
	case "version":
		errorFallback = exitConfiguration
		err = runVersion(args[1:], stdout, stderr)
	case "setup":
		outcome, err = runSetup(args[1:], stdout, stderr)
	case "update":
		errorFallback = exitExecution
		err = runUpdate(args[1:], stdout, stderr)
	case "report":
		errorFallback = exitConfiguration
		err = runReport(args[1:], stdout, stderr)
	case "cache":
		errorFallback = exitConfiguration
		err = runCache(args[1:], stdout, stderr)
	case "completion":
		errorFallback = exitConfiguration
		err = runCompletion(args[1:], stdout, stderr)
	case "inventory":
		errorFallback = exitConfiguration
		err = runInventory(args[1:], stdout, stderr)
	case "config":
		errorFallback = exitConfiguration
		err = runConfig(args[1:], stdout, stderr)
	case "catalog":
		errorFallback = exitConfiguration
		err = runCatalog(args[1:], stdout, stderr)
	case "coverage":
		errorFallback = exitConfiguration
		err = runCoverage(args[1:], stdout, stderr)
	case "evidence":
		errorFallback = exitConfiguration
		err = runEvidence(args[1:], stdout, stderr)
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
	case "quick":
		outcome, err = runScanAlias("quick", args[1:], stdout, stderr)
	case "verify":
		outcome, err = runScanAlias("verify", args[1:], stdout, stderr)
	case "ci":
		outcome, err = runScanAlias("ci", args[1:], stdout, stderr)
	case "full":
		outcome, err = runScanAlias("full", args[1:], stdout, stderr)
	case "fix":
		outcome, err = runFix(args[1:], stdout, stderr)
	case "doctor":
		outcome, err = runDoctor(args[1:], stdout, stderr)
	case "login", "logout", "auth":
		errorFallback = exitConfiguration
		err = runAuthentication(args[0], args[1:], stdin, stdout, stderr)
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
		if len(args) > 1 && args[1] == "advanced" {
			advancedUsage(stdout)
			return 0
		}
		if len(args) > 1 {
			fmt.Fprintf(stderr, "unknown help topic %q\n", args[1])
			return exitConfiguration
		}
		usage(stdout)
		return 0
	default:
		fmt.Fprintf(stderr, "unknown command %q\n", args[0])
		usage(stderr)
		return exitConfiguration
	}
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return exitSuccess
		}
		code := errorExitCode(err, errorFallback)
		fmt.Fprintf(stderr, "error [PRC-EXIT-%d]: %v\n", code, err)
		return code
	}
	return outcome
}

// friendlyTopLevelArgs keeps the common path intentionally small: `prc`
// scans the current directory, and `prc /path/to/project` treats any value that
// is not a known command or root flag as the scan target. The inventory layer
// then performs all path resolution and hardened filesystem access.
func friendlyTopLevelArgs(args []string) []string {
	if len(args) == 0 {
		return []string{"scan"}
	}
	if args[0] == "--help" || args[0] == "-h" || args[0] == "--version" || args[0] == "-v" ||
		strings.HasPrefix(args[0], "-") || isTopLevelCommand(args[0]) {
		return args
	}
	return append([]string{"scan"}, args...)
}

func isTopLevelCommand(value string) bool {
	switch value {
	case "version", "setup", "update", "report", "cache", "completion", "inventory", "config", "catalog", "coverage", "evidence", "benchmark", "pack", "plan", "scan", "quick", "verify", "ci", "full",
		"fix", "doctor", "login", "logout", "auth", "history", "diff", "remediate", "remediate-proposal",
		"explain", "adapter", "exception", "provider", "mcp", "help":
		return true
	default:
		return false
	}
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
	printProductBanner(output, newTerminalStyle("auto", output))
	fmt.Fprintln(output, "Usage: prc [project path]")
	fmt.Fprintln(output, "       prc <command> [options]")
	fmt.Fprintln(output)
	fmt.Fprintln(output, "Start here:")
	fmt.Fprintln(output, "  prc                      Scan the current project and write a detailed report")
	fmt.Fprintln(output, "  prc /path/to/project     Scan another project")
	fmt.Fprintln(output, "  prc quick                Run a small local risk screen")
	fmt.Fprintln(output, "  prc verify               Add the pinned offline secret check")
	fmt.Fprintln(output, "  prc setup                Check this project and show the safest next command")
	fmt.Fprintln(output, "  prc report               Open the newest scan report")
	fmt.Fprintln(output, "  prc login codex          Sign in for an optional Codex review")
	fmt.Fprintln(output, "  prc full codex           Deep AI advice for all nondeterministic controls")
	fmt.Fprintln(output, "  prc update               Check for a newer scanner version")
	fmt.Fprintln(output, "  prc help advanced        Show policy, evidence, CI, and integration commands")
	fmt.Fprintln(output)
	fmt.Fprintln(output, "Common commands: setup, quick, scan, full, report, login, logout, auth, update, version")
}

func advancedUsage(output io.Writer) {
	fmt.Fprintln(output, "Advanced Production Readiness Checklist commands")
	fmt.Fprintln(output)
	fmt.Fprintln(output, "Inspection and automation:")
	fmt.Fprintln(output, "  prc ci                   Write SARIF to stdout for a CI job")
	fmt.Fprintln(output, "  prc doctor               Check local scanning capabilities")
	fmt.Fprintln(output, "  prc coverage             Show rule, predicate, collector, and import coverage")
	fmt.Fprintln(output, "  prc inventory            Inspect the safe content-hashed project inventory")
	fmt.Fprintln(output, "  prc plan                 Show the selected local execution plan")
	fmt.Fprintln(output, "  prc history              Inspect immutable prior runs")
	fmt.Fprintln(output, "  prc diff                 Find prior results invalidated by current changes")
	fmt.Fprintln(output, "  prc cache                Inspect or clean scanner-owned cache files")
	fmt.Fprintln(output, "  prc completion           Generate shell completion")
	fmt.Fprintln(output)
	fmt.Fprintln(output, "Evidence and integrations:")
	fmt.Fprintln(output, "  prc evidence             Inspect or verify exact evidence contracts")
	fmt.Fprintln(output, "  prc catalog              Validate or bundle the rule catalog")
	fmt.Fprintln(output, "  prc config               Validate a declared project configuration")
	fmt.Fprintln(output, "  prc adapter              Inspect explicitly authorized local adapters")
	fmt.Fprintln(output, "  prc exception            Verify signed risk exceptions")
	fmt.Fprintln(output, "  prc mcp serve            Start the read-only agent integration server")
	fmt.Fprintln(output)
	fmt.Fprintln(output, "Separate change workflow:")
	fmt.Fprintln(output, "  prc fix                  Build and verify an isolated change candidate")
	fmt.Fprintln(output, "  prc remediate            Run the lower-level bounded remediation contract")
	fmt.Fprintln(output)
	fmt.Fprintln(output, "Run `prc <command> --help` for exact options. Scanning never invokes the change workflow.")
}

func runEvidence(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 || args[0] == "help" || args[0] == "--help" || args[0] == "-h" {
		fmt.Fprintln(stdout, "Usage:")
		fmt.Fprintln(stdout, "  prc evidence requirements [--authority NAME] [--control ID] [--collector-status all|built_in|missing] [--format human|json]")
		fmt.Fprintln(stdout, "  prc evidence verify-set --set FILE [PROJECT] [--config FILE] [--format human|json]")
		return nil
	}
	switch args[0] {
	case "requirements":
		return runEvidenceRequirements(args[1:], stdout, stderr)
	case "verify-set":
		return runEvidenceVerifySet(args[1:], stdout, stderr)
	default:
		return fmt.Errorf("unknown evidence command %q", args[0])
	}
}

func runEvidenceRequirements(args []string, stdout, stderr io.Writer) error {
	set := flag.NewFlagSet("evidence requirements", flag.ContinueOnError)
	set.SetOutput(stderr)
	root := set.String("catalog-root", defaultCatalogRoot(), "repository containing the PRC catalog")
	authority := set.String("authority", "", "artifact, environment, executed, external_registry, repository, or structured_record")
	controlID := set.String("control", "", "one exact control ID")
	collectorStatus := set.String("collector-status", "all", "all, built_in, or missing")
	format := set.String("format", "human", "human or json")
	if err := set.Parse(args); err != nil {
		return err
	}
	if set.NArg() != 0 {
		return fmt.Errorf("unexpected evidence requirements arguments: %s", strings.Join(set.Args(), " "))
	}
	if *format != "human" && *format != "json" {
		return fmt.Errorf("unsupported format %q", *format)
	}
	reportValue, err := evidencerequirements.Build(*root, evidencerequirements.Filter{
		Authority: controlprogram.Authority(*authority), ControlID: *controlID, CollectorStatus: *collectorStatus,
	})
	if err != nil {
		return err
	}
	if *format == "json" {
		return encodeJSON(stdout, reportValue)
	}
	fmt.Fprintln(stdout, "Evidence producer requirements")
	fmt.Fprintf(stdout, "  Selected clauses    %d/%d across %d controls\n", reportValue.SelectedClauseCount, reportValue.ExactClauseCount, reportValue.SelectedControlCount)
	fmt.Fprintf(stdout, "  Built-in collectors %d\n", reportValue.BuiltInCollectorCount)
	fmt.Fprintf(stdout, "  Missing collectors  %d\n", reportValue.MissingCollectorCount)
	fmt.Fprintf(stdout, "  Signed import route %d\n", reportValue.SignedImportRouteCount)
	fmt.Fprintln(stdout)
	fmt.Fprintln(stdout, "Evidence authority        Selected  Built in  Missing  Signed import")
	for _, item := range reportValue.Authorities {
		fmt.Fprintf(stdout, "  %-23s %8d  %8d  %7d  %13d\n", item.Name, item.SelectedClauseCount,
			item.BuiltInCollectorCount, item.MissingCollectorCount, item.SignedImportRouteCount)
	}
	if *controlID == "" {
		if *authority != "" {
			fmt.Fprintln(stdout)
			fmt.Fprintln(stdout, "Selected clauses")
			for _, requirement := range reportValue.Requirements {
				fmt.Fprintf(stdout, "  %s clause %d · %s · %s\n", requirement.ControlID,
					requirement.ClauseOrdinal, requirement.CollectorStatus, requirement.CollectorID)
			}
		}
		fmt.Fprintln(stdout)
		fmt.Fprintln(stdout, "Add --control ID to print one control's full clause contracts.")
		fmt.Fprintln(stdout, "Use --format json for every selected machine-readable contract.")
		return nil
	}
	for _, requirement := range reportValue.Requirements {
		fmt.Fprintln(stdout)
		fmt.Fprintf(stdout, "%s clause %d · %s · %s\n", requirement.ControlID, requirement.ClauseOrdinal, requirement.Authority, requirement.CollectorStatus)
		fmt.Fprintf(stdout, "  %s\n", requirement.ClauseStatement)
		fmt.Fprintf(stdout, "  Collector: %s\n", requirement.CollectorID)
		fmt.Fprintln(stdout, "  Required facts:")
		for _, fact := range requirement.Facts {
			fmt.Fprintf(stdout, "    - %s (%s; complete=%t)\n", fact.ID, fact.Type, fact.CompleteRequired)
			fmt.Fprintf(stdout, "      %s\n", fact.SourceRequirement)
		}
		if len(requirement.Parameters) > 0 {
			fmt.Fprintln(stdout, "  Inputs sealed before collection:")
			for _, parameter := range requirement.Parameters {
				fmt.Fprintf(stdout, "    - %s (%s; %s)\n", parameter.ID, parameter.Type, parameter.Origin)
				fmt.Fprintf(stdout, "      %s\n", parameter.SourceRequirement)
			}
		}
		fmt.Fprintf(stdout, "  Completeness: %s\n", requirement.CompletenessContract)
		fmt.Fprintf(stdout, "  Freshness: %s\n", requirement.FreshnessContract)
	}
	return nil
}

func runEvidenceVerifySet(args []string, stdout, stderr io.Writer) error {
	set := flag.NewFlagSet("evidence verify-set", flag.ContinueOnError)
	set.SetOutput(stderr)
	root := set.String("catalog-root", defaultCatalogRoot(), "repository containing the PRC catalog")
	target := set.String("target", ".", "project directory to inventory")
	configPath := set.String("config", "", "optional project configuration")
	manifestPath := set.String("set", "", "signed multi-authority evidence-set manifest")
	format := set.String("format", "human", "human or json")
	if err := set.Parse(reorderInterspersedFlags(set, args)); err != nil {
		return err
	}
	if set.NArg() > 1 {
		return fmt.Errorf("evidence verify-set accepts at most one project path")
	}
	if set.NArg() == 1 {
		if flagWasSet(set, "target") {
			return fmt.Errorf("use either a project path or --target, not both")
		}
		*target = set.Arg(0)
	}
	if strings.TrimSpace(*manifestPath) == "" {
		return fmt.Errorf("evidence verify-set requires --set")
	}
	if *format != "human" && *format != "json" {
		return fmt.Errorf("unsupported format %q", *format)
	}
	item, _, err := configuredInventory(*target, *configPath)
	if err != nil {
		return err
	}
	programCatalog, err := controlprogramcatalog.Load(*root)
	if err != nil {
		return err
	}
	verifiedAt := time.Now().UTC()
	executions, verifications, err := evidenceset.VerifyAndEvaluate(programCatalog, item, *manifestPath, verifiedAt)
	if err != nil {
		return exitError(exitPolicyDenied, fmt.Errorf("verify signed evidence set: %w", err))
	}
	reportValue, err := evidenceset.SummarizeVerification(programCatalog, item, executions, verifications, verifiedAt)
	if err != nil {
		return exitError(exitInternal, err)
	}
	if *format == "json" {
		return encodeJSON(stdout, reportValue)
	}
	fmt.Fprintln(stdout, "Evidence set verified")
	fmt.Fprintf(stdout, "  Signatures  verified for %d bundle(s)\n", reportValue.BundleCount)
	fmt.Fprintf(stdout, "  Entries     %d\n", reportValue.EntryCount)
	fmt.Fprintf(stdout, "  Outcomes    %d passed · %d failed · %d not applicable · %d blocked\n",
		reportValue.Outcomes.Passed, reportValue.Outcomes.Failed,
		reportValue.Outcomes.NotApplicable, reportValue.Outcomes.Blocked)
	for _, item := range reportValue.Authorities {
		fmt.Fprintf(stdout, "  %-17s %d entries · policy %s · evidence %s\n",
			item.Authority, item.EntryCount, item.PolicyKeyID, item.EvidenceKeyID)
	}
	fmt.Fprintln(stdout)
	fmt.Fprintln(stdout, "This confirms the signatures and internal catalog, inventory, scope, and predicate bindings.")
	fmt.Fprintln(stdout, "It does not prove producer claims or readiness; failed and blocked entries remain unchanged.")
	return nil
}

func runCoverage(args []string, stdout, stderr io.Writer) error {
	set := flag.NewFlagSet("coverage", flag.ContinueOnError)
	set.SetOutput(stderr)
	root := set.String("catalog-root", defaultCatalogRoot(), "repository containing the PRC catalog")
	format := set.String("format", "human", "human or json")
	if err := set.Parse(args); err != nil {
		return err
	}
	if set.NArg() != 0 {
		return fmt.Errorf("unexpected coverage arguments: %s", strings.Join(set.Args(), " "))
	}
	if *format != "human" && *format != "json" {
		return fmt.Errorf("unsupported format %q", *format)
	}
	reportValue, err := automaticcoverage.Build(*root)
	if err != nil {
		return err
	}
	if *format == "json" {
		return encodeJSON(stdout, reportValue)
	}
	fmt.Fprintln(stdout, "Automatic coverage")
	fmt.Fprintf(stdout, "  Reviewed routing     %d/%d controls (%.1f%%)\n",
		reportValue.ReviewedRoutingControlCount, reportValue.ControlCount,
		percentage(reportValue.ReviewedRoutingControlCount, reportValue.ControlCount))
	fmt.Fprintf(stdout, "  Exact predicates     %d/%d clauses (%.1f%%)\n",
		reportValue.ExactPredicateClauseCount, reportValue.ExactClauseCount,
		percentage(reportValue.ExactPredicateClauseCount, reportValue.ExactClauseCount))
	fmt.Fprintf(stdout, "  Advisory AI route    %d/%d controls (%.1f%%)\n",
		reportValue.AdvisoryAIReviewControlCount, reportValue.NondeterministicControlCount,
		percentage(reportValue.AdvisoryAIReviewControlCount, reportValue.NondeterministicControlCount))
	fmt.Fprintf(stdout, "  Built-in collectors  %d/%d clauses (%.1f%%)\n",
		reportValue.BuiltInCollectorClauseCount, reportValue.ExactClauseCount,
		percentage(reportValue.BuiltInCollectorClauseCount, reportValue.ExactClauseCount))
	fmt.Fprintf(stdout, "  Signed import route  %d/%d clauses (%.1f%%)\n",
		reportValue.SignedImportSupportedClauseCount, reportValue.ExactClauseCount,
		percentage(reportValue.SignedImportSupportedClauseCount, reportValue.ExactClauseCount))
	fmt.Fprintln(stdout)
	fmt.Fprintln(stdout, "Evidence authority        Exact  Built in  Signed import")
	for _, authority := range reportValue.Authorities {
		fmt.Fprintf(stdout, "  %-23s %5d  %8d  %13d\n", authority.Name, authority.ExactClauseCount,
			authority.BuiltInCollectorClauseCount, authority.SignedImportSupportedClauseCount)
	}
	fmt.Fprintln(stdout)
	fmt.Fprintln(stdout, "Built-in means the scanner can collect the evidence itself.")
	fmt.Fprintln(stdout, "Signed import means a trusted external producer can supply facts; it is not a built-in observation.")
	fmt.Fprintln(stdout, "The AI route gives advice for subjective rules; it never creates a verified Pass or Fail.")
	fmt.Fprintln(stdout, "Inspect producer gaps with `prc evidence requirements --collector-status missing`.")
	return nil
}

func percentage(numerator, denominator int) float64 {
	if denominator <= 0 {
		return 0
	}
	return float64(numerator) * 100 / float64(denominator)
}

func runScanAlias(name string, args []string, stdout, stderr io.Writer) (int, error) {
	if len(args) == 1 && (args[0] == "--help" || args[0] == "-h") {
		switch name {
		case "quick":
			fmt.Fprintln(stdout, "Usage: prc quick [project path] [scan options]")
			fmt.Fprintln(stdout, "Runs 18 high-signal local checks, includes all 10,042 controls in the report, and never starts AI.")
			fmt.Fprintln(stdout, "Use prc scan --help for the shared report and output options.")
		case "full":
			fmt.Fprintln(stdout, "Usage: prc full <codex|claude> [project path] [scan options]")
			fmt.Fprintln(stdout, "Runs the 40-check core scan and deep advisory AI review of all 9,356 nondeterministic controls.")
			fmt.Fprintln(stdout, "Requests one primary subagent per rule and a skeptical review per batch; provider internals cannot be independently attested.")
			fmt.Fprintln(stdout, "The provider name explicitly allows bounded, screened remote source processing.")
			fmt.Fprintln(stdout, "Add --plan to inspect source and exact batch limits without starting the provider.")
			fmt.Fprintln(stdout, "Add --details to print every local check and completed AI rule review; the HTML report always contains all of them.")
			fmt.Fprintln(stdout, "Use prc scan --help for advanced review, report, and output options.")
		case "verify":
			fmt.Fprintln(stdout, "Usage: prc verify [project path] [scan options]")
			fmt.Fprintln(stdout, "Runs the core scan plus the bundled, digest-pinned Gitleaks adapter in a locked-down local container.")
			fmt.Fprintln(stdout, "The image must already exist locally: PRC passes --pull=never, disables networking, and mounts a sealed read-only snapshot.")
			fmt.Fprintln(stdout, "Use prc doctor first to inspect local capabilities. Use prc scan for custom adapters or AI review.")
		case "ci":
			fmt.Fprintln(stdout, "Usage: prc ci [project path] [scan options] > prc-results.sarif")
			fmt.Fprintln(stdout, "Runs the normal local checks, writes SARIF only, and never starts AI or creates an HTML report.")
		}
		return exitSuccess, nil
	}
	for _, argument := range args {
		if argument == "--profile" || strings.HasPrefix(argument, "--profile=") {
			return exitInternal, exitError(exitConfiguration, fmt.Errorf("prc %s selects its own profile; use prc scan --profile for a custom profile", name))
		}
	}
	switch name {
	case "quick":
		for _, argument := range args {
			if argument == "--ai" || strings.HasPrefix(argument, "--ai=") ||
				argument == "--review-provider" || strings.HasPrefix(argument, "--review-provider=") {
				return exitInternal, exitError(exitConfiguration, errors.New("prc quick is local only; use prc full codex or prc full claude for AI review"))
			}
		}
		return runScan(append([]string{"--profile", "prc/quick"}, args...), stdout, stderr)
	case "verify":
		owned := []string{"catalog-root", "mode", "adapter-manifest", "adapter-registry", "adapter-id", "adapter-data"}
		for _, name := range owned {
			if containsFlag(args, name) {
				return exitInternal, exitError(exitConfiguration, fmt.Errorf("prc verify owns --%s; use prc scan for custom adapter execution", name))
			}
		}
		for _, argument := range args {
			if argument == "--ai" || strings.HasPrefix(argument, "--ai=") ||
				argument == "--review-provider" || strings.HasPrefix(argument, "--review-provider=") ||
				argument == "--allow-remote-source-processing" || strings.HasPrefix(argument, "--allow-remote-source-processing=") {
				return exitInternal, exitError(exitConfiguration, errors.New("prc verify is local only; use prc full codex, prc full claude, or prc scan for AI review"))
			}
		}
		return runScan(verifyScanArguments(args), stdout, stderr)
	case "ci":
		for _, argument := range args {
			if argument == "--ai" || strings.HasPrefix(argument, "--ai=") ||
				argument == "--review-provider" || strings.HasPrefix(argument, "--review-provider=") {
				return exitInternal, exitError(exitConfiguration, errors.New("prc ci is local only; use prc scan for an advisory AI review"))
			}
			if argument == "--format" || strings.HasPrefix(argument, "--format=") ||
				argument == "--report" || strings.HasPrefix(argument, "--report=") {
				return exitInternal, exitError(exitConfiguration, errors.New("prc ci owns SARIF output; use prc scan for another format or an HTML report"))
			}
		}
		return runScan(append([]string{"--format", "sarif", "--no-report"}, args...), stdout, stderr)
	case "full":
		if len(args) == 0 || args[0] != "codex" && args[0] != "claude" {
			return exitInternal, exitError(exitConfiguration, errors.New("usage: prc full codex [project path] or prc full claude [project path]"))
		}
		for _, argument := range args[1:] {
			if argument == "--ai" || strings.HasPrefix(argument, "--ai=") ||
				argument == "--review-provider" || strings.HasPrefix(argument, "--review-provider=") {
				return exitInternal, exitError(exitConfiguration, errors.New("prc full already selects the AI provider"))
			}
		}
		forwarded := []string{"--ai", args[0]}
		remaining := append([]string(nil), args[1:]...)
		for index, argument := range remaining {
			if argument == "--plan" {
				remaining[index] = "--review-plan"
			}
		}
		if !containsFlag(remaining, "review-depth") {
			forwarded = append(forwarded, "--review-depth", "deep")
		}
		if !containsFlag(remaining, "review-workers") {
			forwarded = append(forwarded, "--review-workers", "4")
		}
		if args[0] == "codex" && !containsFlag(remaining, "review-effort") {
			forwarded = append(forwarded, "--review-effort", "xhigh")
		}
		return runScan(append(forwarded, remaining...), stdout, stderr)
	default:
		panic("unknown scan alias")
	}
}

func verifyScanArguments(args []string) []string {
	manifest := filepath.Join(defaultCatalogRoot(), "adapters", "gitleaks-v8.30.0.yaml")
	return append([]string{
		"--profile", "prc/core-repository",
		"--mode", engine.ExecutionModeVerifyLocal,
		"--adapter-manifest", manifest,
	}, args...)
}

func containsFlag(args []string, name string) bool {
	prefix := "--" + name
	for _, argument := range args {
		if argument == prefix || strings.HasPrefix(argument, prefix+"=") {
			return true
		}
	}
	return false
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
	store, err := state.Open(commandContext, *stateDirectory)
	if err != nil {
		return exitError(exitConfiguration, err)
	}
	defer store.Close()
	run, err := store.LoadRun(commandContext, loaded.Record.Run.RunID)
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
	catalogRoot := set.String("catalog-root", defaultCatalogRoot(), "directory containing the bundled PRC catalog")
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
	execution, err := provider.Run(commandContext, launchPlan, task)
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
	catalogRoot := set.String("catalog-root", defaultCatalogRoot(), "directory containing the bundled PRC catalog")
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
	store, err := state.Open(commandContext, *stateDirectory)
	if err != nil {
		return err
	}
	defer store.Close()
	if args[0] == "check" {
		if set.NArg() != 0 {
			return errors.New("history check accepts no positional arguments")
		}
		if err := store.IntegrityCheck(commandContext); err != nil {
			return err
		}
		counts, err := store.Counts(commandContext)
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
		run, err := store.LoadRun(commandContext, set.Arg(0))
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
	runs, err := store.ListRuns(commandContext, state.Query{
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
	store, err := state.Open(commandContext, *stateDirectory)
	if err != nil {
		return err
	}
	defer store.Close()
	base, err := store.LoadRun(commandContext, *baseRunID)
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
	dataFlags := repeatedStringFlag{}
	set.Var(&dataFlags, "data", "manifest data mount as NAME=PATH; repeatable")
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
	dataSources, err := parseNamedPaths(dataFlags, "--data")
	if err != nil {
		return err
	}
	if commandName == "run-oci" {
		runID, err = randomRunID()
		if err != nil {
			return err
		}
	}
	if commandName == "plan-oci" {
		ociPlan, err := adapter.BuildOCIPlanWithData(*runtime, item.Root, runID, manifest, dataSources)
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
	ociPlan, err := adapter.BuildSnapshotOCIPlanWithData(*runtime, snapshot, runID, manifest, dataSources)
	if err != nil {
		return err
	}
	input, err := adapter.ExecutionInput(manifest, runID, adapter.Subject{
		TargetName: item.TargetName, TargetCommit: item.GitCommit, InventoryDigest: item.Digest,
	}, inventoryFacts(item), map[string]any{})
	if err != nil {
		return err
	}
	result, err := adapter.RunOCI(commandContext, ociPlan, manifest, input)
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
	set.Usage = func() {
		if name == "scan" {
			fmt.Fprintln(stderr, "Usage: prc scan [project path] [options]")
			fmt.Fprintln(stderr, "Scans without modifying the project and writes a detailed report by default.")
			fmt.Fprintln(stderr, "Examples: prc scan   |   prc scan --ai codex   |   prc scan ../project")
		} else {
			fmt.Fprintf(stderr, "Usage: prc %s [options]\n", name)
		}
		fmt.Fprintln(stderr)
		set.PrintDefaults()
	}
	target := set.String("target", ".", "repository to inspect")
	catalogRoot := set.String("catalog-root", defaultCatalogRoot(), "directory containing the PRC catalog")
	profile := set.String("profile", "prc/core-repository", "profile ID")
	return set, target, catalogRoot, profile
}

func defaultCatalogRoot() string {
	executable, err := os.Executable()
	if err != nil {
		return "."
	}
	return resolveDefaultCatalogRoot(executable, ".")
}

func resolveDefaultCatalogRoot(executable, fallback string) string {
	if executable != "" {
		root := filepath.Dir(executable)
		// Installed scanner resources must win over look-alike catalog files in
		// the project being scanned. Development binaries have no adjacent
		// bundle and therefore continue to use the source checkout fallback.
		if hasBundledCatalog(root) {
			return root
		}
	}
	return fallback
}

func hasBundledCatalog(root string) bool {
	for _, relative := range []string{
		"catalog/profiles/core-repository.yaml",
		"catalog/assertions/core-repository.yaml",
		"catalog/objectives/core-repository.yaml",
	} {
		info, err := os.Lstat(filepath.Join(root, filepath.FromSlash(relative)))
		if err != nil || !info.Mode().IsRegular() {
			return false
		}
	}
	return true
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
	if item.GitCommit != "" {
		fmt.Fprintf(stdout, "Git HEAD: %s (%s worktree)\n", item.GitCommit, inventoryFact(item, "repository.git_worktree_state"))
	}
	fmt.Fprintf(stdout, "Files: %d (%d recognized source files)\n", item.FileCount, item.SourceFiles)
	fmt.Fprintf(stdout, "Hashed bytes: %s\n", inventoryFact(item, "repository.inventory_bytes"))
	fmt.Fprintf(stdout, "Automatic exclusions reported: %d\n", inventoryFactCount(item, "repository.exclusion"))
	fmt.Fprintf(stdout, "Reviewed .prcignore exclusions: %d\n", inventoryFactCount(item, "repository.user_exclusion"))
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

func parseNamedPaths(values repeatedStringFlag, label string) (map[string]string, error) {
	result := map[string]string{}
	for _, value := range values {
		parts := strings.SplitN(value, "=", 2)
		if len(parts) != 2 || strings.TrimSpace(parts[0]) == "" || strings.TrimSpace(parts[1]) == "" {
			return nil, fmt.Errorf("%s requires NAME=PATH, got %q", label, value)
		}
		if _, exists := result[parts[0]]; exists {
			return nil, fmt.Errorf("%s repeats name %q", label, parts[0])
		}
		result[parts[0]] = parts[1]
	}
	return result, nil
}

func parseScanAdapterData(values repeatedStringFlag) (map[string]map[string]string, error) {
	result := map[string]map[string]string{}
	for _, value := range values {
		parts := strings.SplitN(value, "=", 2)
		if len(parts) != 2 || strings.TrimSpace(parts[1]) == "" {
			return nil, fmt.Errorf("--adapter-data requires ADAPTER_ID/NAME=PATH, got %q", value)
		}
		separator := strings.LastIndex(parts[0], "/")
		if separator < 1 || separator == len(parts[0])-1 {
			return nil, fmt.Errorf("--adapter-data requires ADAPTER_ID/NAME=PATH, got %q", value)
		}
		adapterID, name := parts[0][:separator], parts[0][separator+1:]
		if result[adapterID] == nil {
			result[adapterID] = map[string]string{}
		}
		if _, exists := result[adapterID][name]; exists {
			return nil, fmt.Errorf("--adapter-data repeats %s/%s", adapterID, name)
		}
		result[adapterID][name] = parts[1]
	}
	return result, nil
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
	dataSources map[string]string,
) (model.AdapterExecution, map[string][]byte, error) {
	if err := ctx.Err(); err != nil {
		return model.AdapterExecution{}, nil, fmt.Errorf("adapter execution budget is exhausted: %w", err)
	}
	runID, err := randomRunID()
	if err != nil {
		return model.AdapterExecution{}, nil, err
	}
	subject := adapter.Subject{
		TargetName: item.TargetName, TargetCommit: item.GitCommit, InventoryDigest: item.Digest,
	}
	input, err := adapter.ExecutionInput(
		resolved.Manifest, runID, subject, inventoryFacts(item), map[string]any{},
	)
	if err != nil {
		return model.AdapterExecution{}, nil, err
	}
	snapshot, err := adapter.PrepareSnapshotForManifest(item, resolved.Manifest)
	if err != nil {
		return model.AdapterExecution{}, nil, err
	}
	defer snapshot.Close()
	plan, err := adapter.BuildSnapshotOCIPlanWithData(runtime, snapshot, runID, resolved.Manifest, dataSources)
	if err != nil {
		return model.AdapterExecution{}, nil, err
	}
	output, err := adapter.RunOCI(ctx, plan, resolved.Manifest, input)
	if err != nil {
		return model.AdapterExecution{}, nil, err
	}
	var execution model.AdapterExecution
	if resolved.Resolution == nil {
		execution, err = adapter.BindExecution(runID, subject, resolved.Manifest, output)
	} else {
		execution, err = adapter.BindExecutionWithResolution(
			runID, subject, resolved.Manifest, *resolved.Resolution, output,
		)
	}
	if err != nil {
		return model.AdapterExecution{}, nil, err
	}
	return execution, output.ArtifactPayloads, nil
}

func runScan(args []string, stdout, stderr io.Writer) (int, error) {
	set, target, catalogRoot, profile := parseCommon("scan", args, stderr)
	configPath := set.String("config", "", "optional validated project configuration")
	mode := set.String("mode", engine.ExecutionModeInspect, "execution mode: inspect or verify-local")
	format := set.String("format", "human", "human, json, markdown, html, sarif, or junit")
	colorMode := set.String("color", "auto", "human output color: auto, always, or never")
	details := set.Bool("details", false, "show every local check and completed AI rule review in human output")
	reportPath := set.String("report", "", "write a detailed HTML report to this new file")
	noReport := set.Bool("no-report", false, "do not create the default HTML report")
	stateDirectory := set.String("state-dir", "", "optional directory for content-addressed evidence and run records")
	exitPolicy := set.String("exit-policy", "profile", "profile, no-go, or never")
	reviewProvider := set.String("review-provider", "none", "advisory review: none, codex, or claude")
	aiProvider := set.String("ai", "", "simple AI review: codex or claude; also acknowledges screened remote source processing")
	reviewExecutable := set.String("review-executable", "", "optional exact Codex or Claude executable")
	reviewModel := set.String("review-model", "", "optional model override for advisory review")
	reviewEffort := set.String("review-effort", "high", "advisory reasoning effort: high, or xhigh for Codex")
	reviewDepth := set.String("review-depth", "standard", "advisory review depth: standard or deep")
	reviewStateDirectory := set.String("review-state-dir", "", "private resumable AI review state outside the target")
	allowRemoteReview := set.Bool("allow-remote-source-processing", false, "allow screened source excerpts to be sent to the selected remote AI provider")
	reviewBatchSize := set.Int("review-batch-size", 8, "controls per provider call, from 1 to 8; each still gets one subagent")
	reviewWorkers := set.Int("review-workers", 1, "parallel provider calls, from 1 to 4")
	reviewTimeout := set.Duration("review-timeout", 30*time.Minute, "timeout for each resumable provider batch")
	reviewMaxCost := set.Float64("review-max-cost-usd", 0, "optional Claude cost limit in USD for each batch")
	reviewMaxBatches := set.Int("review-max-batches", 1500, "hard limit on provider batches for the complete review")
	reviewMaxDuration := set.Duration("review-max-duration", 24*time.Hour, "total resumable AI review time limit")
	reviewPlanOnly := set.Bool("review-plan", false, "screen source and print the exact AI work plan without starting a provider")
	reviewDetails := set.Bool("review-details", false, "print every completed advisory rule review in the terminal")
	reviewControls := repeatedStringFlag{}
	set.Var(&reviewControls, "review-control", "review only this control ID for debugging; repeatable")
	adapterManifests := repeatedStringFlag{}
	set.Var(&adapterManifests, "adapter-manifest", "immutable OCI adapter manifest authorized by the selected profile; repeatable")
	adapterRegistry := set.String("adapter-registry", "", "optional adapter registry lockfile")
	adapterIDs := repeatedStringFlag{}
	set.Var(&adapterIDs, "adapter-id", "adapter ID to resolve from --adapter-registry; repeatable")
	adapterRuntime := set.String("adapter-runtime", "docker", "docker or podman executable for the authorized adapter")
	adapterData := repeatedStringFlag{}
	set.Var(&adapterData, "adapter-data", "read-only adapter data as ADAPTER_ID/NAME=PATH; repeatable")
	evidenceBundlePath := set.String("evidence-bundle", "", "signed authoritative deterministic evidence bundle")
	evidenceTrustStorePath := set.String("evidence-trust-store", "", "trust store for the evidence bundle signatures")
	evidencePolicySignaturePath := set.String("evidence-policy-signature", "", "policy signature for the exact evidence bundle bytes")
	evidenceSignaturePath := set.String("evidence-signature", "", "authority-scoped evidence signature for the exact bundle bytes")
	evidenceSetPath := set.String("evidence-set", "", "one manifest for signed evidence from multiple authorities")
	if err := set.Parse(reorderInterspersedFlags(set, args)); err != nil {
		return exitInternal, exitError(exitConfiguration, err)
	}
	if set.NArg() > 1 {
		return exitInternal, exitError(exitConfiguration, fmt.Errorf("scan accepts at most one project path"))
	}
	if set.NArg() == 1 {
		if flagWasSet(set, "target") {
			return exitInternal, exitError(exitConfiguration, fmt.Errorf("use either a project path or --target, not both"))
		}
		*target = set.Arg(0)
	}
	if *noReport && *reportPath != "" {
		return exitInternal, exitError(exitConfiguration, fmt.Errorf("--report and --no-report are mutually exclusive"))
	}
	if *format != "human" && *format != "json" && *format != "markdown" &&
		*format != "html" && *format != "sarif" && *format != "junit" {
		return exitInternal, exitError(exitConfiguration, fmt.Errorf("unsupported format %q", *format))
	}
	if *reviewPlanOnly && *format != "human" && *format != "json" {
		return exitInternal, exitError(exitConfiguration, fmt.Errorf("--review-plan supports human or json output"))
	}
	if *colorMode != "auto" && *colorMode != "always" && *colorMode != "never" {
		return exitInternal, exitError(exitConfiguration, fmt.Errorf("unsupported color mode %q", *colorMode))
	}
	if *format != "human" && flagWasSet(set, "color") {
		return exitInternal, exitError(exitConfiguration, fmt.Errorf("--color applies only to human output"))
	}
	if *format != "human" && flagWasSet(set, "details") {
		return exitInternal, exitError(exitConfiguration, fmt.Errorf("--details applies only to human output"))
	}
	if *exitPolicy != "profile" && *exitPolicy != "no-go" && *exitPolicy != "never" {
		return exitInternal, exitError(exitConfiguration, fmt.Errorf("unsupported exit policy %q", *exitPolicy))
	}
	if *mode != engine.ExecutionModeInspect && *mode != engine.ExecutionModeVerifyLocal {
		return exitInternal, exitError(exitConfiguration, fmt.Errorf("unsupported execution mode %q", *mode))
	}
	if *aiProvider != "" {
		if *aiProvider != "codex" && *aiProvider != "claude" {
			return exitInternal, exitError(exitConfiguration, fmt.Errorf("unsupported AI provider %q; use codex or claude", *aiProvider))
		}
		if flagWasSet(set, "review-provider") && *reviewProvider != *aiProvider {
			return exitInternal, exitError(exitConfiguration, fmt.Errorf("--ai and --review-provider select different providers"))
		}
		if flagWasSet(set, "allow-remote-source-processing") && !*allowRemoteReview {
			return exitInternal, exitError(exitConfiguration, fmt.Errorf("--ai cannot be combined with --allow-remote-source-processing=false"))
		}
		*reviewProvider = *aiProvider
		*allowRemoteReview = true
	}
	if *reviewProvider != "none" && *reviewProvider != "codex" && *reviewProvider != "claude" {
		return exitInternal, exitError(exitConfiguration, fmt.Errorf("unsupported review provider %q", *reviewProvider))
	}
	reviewFlags := []string{"review-executable", "review-model", "review-effort", "review-depth", "review-state-dir", "allow-remote-source-processing", "review-batch-size", "review-workers", "review-timeout", "review-max-cost-usd", "review-max-batches", "review-max-duration", "review-plan", "review-details", "review-control"}
	if *reviewProvider == "none" {
		for _, name := range reviewFlags {
			if flagWasSet(set, name) {
				return exitInternal, exitError(exitConfiguration, fmt.Errorf("--%s requires --review-provider codex or claude", name))
			}
		}
	}
	if err := validateUniqueFlagValues(reviewControls, "--review-control"); err != nil {
		return exitInternal, exitError(exitConfiguration, err)
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
	if requestedAdapters == 0 && len(adapterData) > 0 {
		return exitInternal, exitError(exitConfiguration, fmt.Errorf("--adapter-data requires a requested adapter"))
	}
	evidenceFlags := map[string]string{
		"--evidence-bundle": *evidenceBundlePath, "--evidence-trust-store": *evidenceTrustStorePath,
		"--evidence-policy-signature": *evidencePolicySignaturePath, "--evidence-signature": *evidenceSignaturePath,
	}
	evidenceFlagCount := 0
	for _, value := range evidenceFlags {
		if strings.TrimSpace(value) != "" {
			evidenceFlagCount++
		}
	}
	if evidenceFlagCount != 0 && evidenceFlagCount != len(evidenceFlags) {
		missing := make([]string, 0, len(evidenceFlags)-evidenceFlagCount)
		for name, value := range evidenceFlags {
			if strings.TrimSpace(value) == "" {
				missing = append(missing, name)
			}
		}
		sort.Strings(missing)
		return exitInternal, exitError(exitConfiguration, fmt.Errorf("signed evidence import requires all four evidence flags; missing %s", strings.Join(missing, ", ")))
	}
	if strings.TrimSpace(*evidenceSetPath) != "" && evidenceFlagCount != 0 {
		return exitInternal, exitError(exitConfiguration, fmt.Errorf("--evidence-set cannot be combined with the four single-bundle evidence flags"))
	}
	if err := validateUniqueFlagValues(adapterManifests, "--adapter-manifest"); err != nil {
		return exitInternal, exitError(exitConfiguration, err)
	}
	if err := validateUniqueFlagValues(adapterIDs, "--adapter-id"); err != nil {
		return exitInternal, exitError(exitConfiguration, err)
	}
	dataByAdapter, err := parseScanAdapterData(adapterData)
	if err != nil {
		return exitInternal, exitError(exitConfiguration, err)
	}
	var progressStyle terminalStyle
	if *format == "human" {
		progressStyle = newTerminalStyle(*colorMode, stdout)
		printProductBanner(stdout, progressStyle)
		displayTarget := *target
		if absolute, absoluteErr := filepath.Abs(*target); absoluteErr == nil {
			displayTarget = absolute
		}
		fmt.Fprintf(stdout, "  Project  %s\n\n", progressStyle.paint(ansiCyan, terminalText(displayTarget)))
		fmt.Fprintln(stdout, progressStyle.paint(ansiBlue, "[1/4] Building a safe, content-hashed inventory..."))
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
	if *format == "human" {
		fmt.Fprintln(stdout, progressStyle.paint(ansiBlue, "[2/4] Running the selected deterministic checks..."))
	}
	executions := []model.AdapterExecution{}
	artifactPayloads := map[string][]byte{}
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
		resolvedIDs := map[string]bool{}
		for _, resolved := range resolvedAdapters {
			resolvedIDs[resolved.Manifest.ID] = true
		}
		for adapterID := range dataByAdapter {
			if !resolvedIDs[adapterID] {
				return exitInternal, exitError(exitConfiguration, fmt.Errorf("--adapter-data references unrequested adapter %s", adapterID))
			}
		}
		adapterContext := commandContext
		cancelAdapters := func() {}
		if validation != nil {
			adapterContext, cancelAdapters = context.WithTimeout(
				adapterContext,
				time.Duration(validation.Configuration.Execution.MaxDurationSeconds)*time.Second,
			)
		}
		defer cancelAdapters()
		for _, resolved := range resolvedAdapters {
			execution, payloads, executeErr := executeScanAdapter(
				adapterContext, item, resolved, *adapterRuntime, dataByAdapter[resolved.Manifest.ID],
			)
			if executeErr != nil {
				return exitInternal, exitError(exitExecution, executeErr)
			}
			executions = append(executions, execution)
			for descriptor, payload := range payloads {
				if existing, exists := artifactPayloads[descriptor]; exists && !bytes.Equal(existing, payload) {
					return exitInternal, exitError(exitExecution, fmt.Errorf("adapters produced conflicting payloads for %s", descriptor))
				}
				artifactPayloads[descriptor] = payload
			}
		}
	}
	run, err := scanner.ScanMode(*profile, item, executions, *mode)
	if err != nil {
		return exitInternal, exitError(exitExecution, err)
	}
	profileRun := run
	programCatalog, programErr := controlprogramcatalog.Load(*catalogRoot)
	var completeRun model.RunResult
	var attachErr error
	if programErr == nil {
		exactExecutions, executionErr := repositoryevidence.EvaluateSupported(
			commandContext, programCatalog, item, profileRun.CompletedAt,
		)
		if executionErr != nil {
			return exitInternal, exitError(exitExecution, fmt.Errorf("collect exact repository evidence: %w", executionErr))
		}
		if evidenceFlagCount == len(evidenceFlags) || strings.TrimSpace(*evidenceSetPath) != "" {
			var imported []controlruntime.Execution
			var verifications []evidencebundle.Verification
			var importErr error
			if strings.TrimSpace(*evidenceSetPath) != "" {
				imported, verifications, importErr = evidenceset.VerifyAndEvaluate(
					programCatalog, item, *evidenceSetPath, profileRun.CompletedAt,
				)
			} else {
				var verification evidencebundle.Verification
				imported, verification, importErr = evidencebundle.VerifyAndEvaluate(
					programCatalog, item, *evidenceBundlePath, *evidenceTrustStorePath,
					*evidencePolicySignaturePath, *evidenceSignaturePath, profileRun.CompletedAt,
				)
				if importErr == nil {
					verifications = []evidencebundle.Verification{verification}
				}
			}
			if importErr != nil {
				return exitInternal, exitError(exitPolicyDenied, fmt.Errorf("import signed authoritative evidence: %w", importErr))
			}
			seenTemplates := make(map[string]bool, len(exactExecutions)+len(imported))
			for _, execution := range exactExecutions {
				seenTemplates[execution.TemplateID] = true
			}
			for _, execution := range imported {
				if seenTemplates[execution.TemplateID] {
					return exitInternal, exitError(exitConfiguration, fmt.Errorf("signed evidence duplicates an already evaluated template %s", execution.TemplateID))
				}
				seenTemplates[execution.TemplateID] = true
			}
			exactExecutions = append(exactExecutions, imported...)
			for _, verification := range verifications {
				profileRun.AuthoritativeEvidence = append(profileRun.AuthoritativeEvidence, model.AuthoritativeEvidenceVerification{
					SchemaVersion: verification.SchemaVersion, BundleID: verification.BundleID,
					BundleSHA256: verification.BundleSHA256, PolicySHA256: verification.PolicySHA256,
					CatalogSHA256:   verification.CatalogSHA256,
					InventorySHA256: verification.InventorySHA256, Authority: verification.Authority,
					EntryCount:        verification.EntryCount,
					Entries:           append([]model.AuthoritativeEvidenceEntry(nil), verification.Entries...),
					PolicySignature:   verification.PolicySignature,
					EvidenceSignature: verification.EvidenceSignature,
				})
			}
		}
		completeRun, attachErr = fullscan.AttachProgramExecutions(*catalogRoot, scanner.Catalog, profileRun, exactExecutions)
	} else {
		if evidenceFlagCount != 0 || strings.TrimSpace(*evidenceSetPath) != "" {
			return exitInternal, exitError(exitConfiguration, fmt.Errorf("signed evidence import requires the complete released control catalog: %w", programErr))
		}
		// Deliberately minimal custom catalogs used for focused local scans do
		// not contain the complete registry or exact-program corpus. Attach
		// decides whether that absence is the supported minimal case or a
		// broken released bundle.
		completeRun, attachErr = fullscan.Attach(*catalogRoot, scanner.Catalog, profileRun)
	}
	if attachErr == nil {
		run = completeRun
	} else if !errors.Is(attachErr, fullscan.ErrRegistryUnavailable) {
		return exitInternal, exitError(exitConfiguration, attachErr)
	}
	if *format == "human" {
		fmt.Fprintln(stdout, progressStyle.paint(ansiBlue, "[3/4] Classifying complete-catalog coverage without guessing..."))
	}
	var reviewSummary controlreview.Summary
	var reviewFailure error
	if *reviewProvider != "none" {
		reviewOptions := controlreview.Options{
			Provider: *reviewProvider, Executable: *reviewExecutable, Model: *reviewModel,
			ReasoningEffort: *reviewEffort, ReviewDepth: *reviewDepth, StateDirectory: *reviewStateDirectory,
			SchemaPath:                  filepath.Join(*catalogRoot, "schemas", "control-review-output.schema.json"),
			AllowRemoteSourceProcessing: *allowRemoteReview, BatchSize: *reviewBatchSize,
			Workers: *reviewWorkers, Timeout: *reviewTimeout, MaxCostUSD: *reviewMaxCost,
			MaxBatches: *reviewMaxBatches, MaxDuration: *reviewMaxDuration,
			ControlIDs: append([]string{}, reviewControls...),
		}
		if *reviewPlanOnly {
			preview, previewErr := controlreview.BuildPreview(commandContext, run, reviewOptions)
			if previewErr != nil {
				return exitInternal, exitError(exitExecution, previewErr)
			}
			if *format == "json" {
				if err := encodeJSON(stdout, preview); err != nil {
					return exitInternal, err
				}
			} else {
				fmt.Fprintln(stdout, "AI review plan — no provider was started")
				fmt.Fprintf(stdout, "  Provider          %s%s\n", preview.Provider, optionalModel(preview.Model))
				fmt.Fprintf(stdout, "  Work              %d controls in %d batches of up to %d\n", preview.Controls, preview.Batches, preview.BatchSize)
				fmt.Fprintf(stdout, "  Parallel workers  %d\n", preview.Workers)
				fmt.Fprintf(stdout, "  Source snapshot   %d screened text files · %s\n", preview.SourceFiles, formatByteCount(preview.SourceBytes))
				fmt.Fprintf(stdout, "  Selected context  %d excerpts · %s across all batches · max %d files / %s per batch\n",
					preview.ContextFiles, formatByteCount(preview.ContextBytes), preview.MaxContextFiles, formatByteCount(int64(preview.MaxContextBytes)))
				if preview.ContextLimited > 0 {
					fmt.Fprintf(stdout, "  Context limits    %d/%d batches reached an excerpt cap; each sealed task names the limitation\n",
						preview.ContextLimited, preview.Batches)
				}
				fmt.Fprintf(stdout, "  Omitted files     %d\n", preview.OmittedFiles)
				fmt.Fprintf(stdout, "  Safety limits     %d batches · %s total · %s per batch\n", preview.MaximumBatches, preview.MaximumDuration, preview.TimeoutPerBatch)
				if *reviewMaxCost > 0 {
					fmt.Fprintf(stdout, "  Claude limit      $%.4f per batch; this is not a whole-run cost cap\n", *reviewMaxCost)
				}
				for _, limitation := range preview.Limitations {
					fmt.Fprintf(stdout, "  Limit             %s\n", terminalText(limitation))
				}
				fmt.Fprintln(stdout, "\nAI output would be advice only. Run the same command without --plan to start it.")
			}
			return exitSuccess, nil
		}
		var reviewProgress func(controlreview.Progress)
		if *format == "human" {
			reviewProgress = aiReviewProgressPrinter(stdout, progressStyle)
		}
		reviewOptions.Progress = reviewProgress
		reviewed, summary, reviewErr := controlreview.Apply(commandContext, run, reviewOptions)
		if reviewed.ControlCatalog != nil {
			run, reviewSummary = reviewed, summary
		}
		if reviewErr != nil {
			reviewFailure = reviewErr
			if *format == "human" {
				message := "      AI review stopped before a valid batch completed. The local scan report will still be written."
				if summary.CompletedBatches > 0 {
					message = "      AI review stopped early. Completed batches will be saved in a partial report and reused on resume."
				}
				fmt.Fprintln(stdout, progressStyle.paint(ansiYellow, message))
			}
		}
	}
	var stateStore *state.Store
	if *stateDirectory != "" {
		stateStore, err = state.Open(commandContext, *stateDirectory)
		if err != nil {
			return exitInternal, exitError(exitConfiguration, err)
		}
		defer stateStore.Close()
	}
	if err := evidence.WriteRunWithArtifacts(*stateDirectory, run, artifactPayloads); err != nil {
		return exitInternal, exitError(exitInternal, err)
	}
	if stateStore != nil {
		if err := stateStore.IndexRun(commandContext, run); err != nil {
			return exitInternal, exitError(exitInternal, err)
		}
	}
	writtenReport := ""
	if *format == "human" {
		finalStage := "[4/4] Preparing the final result..."
		if *reportPath != "" || !*noReport {
			finalStage = "[4/4] Writing the report..."
		}
		fmt.Fprintln(stdout, progressStyle.paint(ansiBlue, finalStage))
		fmt.Fprintln(stdout)
	}
	if *reportPath != "" || *format == "human" && !*noReport {
		writtenReport, err = writeScanReport(run, *reportPath)
		if err != nil {
			return exitInternal, exitError(exitInternal, err)
		}
	}
	if *format == "json" {
		if err := encodeJSON(stdout, run); err != nil {
			return exitInternal, exitError(exitInternal, err)
		}
	} else if *format == "human" {
		style := newTerminalStyle(*colorMode, stdout)
		printScanSummary(stdout, run, reviewSummary, style, writtenReport, *details, *details || *reviewDetails)
	} else if err := report.Write(*format, stdout, run); err != nil {
		return exitInternal, exitError(exitInternal, err)
	}
	if reviewFailure != nil {
		return exitInternal, exitError(exitExecution, reviewFailure)
	}
	switch *exitPolicy {
	case "profile":
		return scanTerminalExitCode(profileTerminalState(run)), nil
	case "no-go":
		return scanNoGoExitCode(profileTerminalState(run)), nil
	case "never":
		return exitSuccess, nil
	}
	panic("validated exit policy was not handled")
}

func aiReviewProgressPrinter(output io.Writer, style terminalStyle) func(controlreview.Progress) {
	lastCompleted := -1
	lastPercent := -1
	lastElapsedBucket := -1
	liveLine := false
	return func(progress controlreview.Progress) {
		if progress.Phase == "prepared" {
			newBatches := max(progress.TotalBatches-progress.ReusedBatches, 0)
			parallelJobs := min(progress.Workers, newBatches)
			fmt.Fprintln(output, style.paint(ansiBlue, "      ╭─ AI REVIEW"))
			fmt.Fprintf(output, "      │ %s · %s thinking · %s review\n",
				style.paint(ansiCyan, authenticationProviderTitle(progress.Provider)),
				terminalText(progress.ReasoningEffort), terminalText(progress.ReviewDepth))
			jobPlan := "no new " + authenticationProviderTitle(progress.Provider) + " jobs"
			if parallelJobs > 0 {
				jobPlan = "up to " + countedNoun(parallelJobs, authenticationProviderTitle(progress.Provider)+" job", authenticationProviderTitle(progress.Provider)+" jobs")
			}
			fmt.Fprintf(output, "      │ %s · %s · %s\n",
				countedNoun(progress.TotalControls, "rule", "rules"), countedNoun(progress.TotalBatches, "batch", "batches"), jobPlan)
			fmt.Fprintf(output, "      │ %s will be reused\n", countedNoun(progress.ReusedBatches, "cached batch", "cached batches"))
			if newBatches == 0 {
				fmt.Fprintln(output, "      ╰─ All AI work is cached; the provider will not start.")
			} else {
				fmt.Fprintf(output, "      │ %s across this run\n", countedNoun(progress.TotalAgentSlots, "requested rule-review assignment", "requested rule-review assignments"))
				if progress.MaxCostUSD > 0 {
					fmt.Fprintf(output, "      │ Enforced limit: $%.4f per new Claude batch\n", progress.MaxCostUSD)
				}
				fmt.Fprintln(output, "      │ Inner-agent counts are requested work; exact state is not exposed.")
				fmt.Fprintln(output, "      ╰─ Ctrl+C stops safely. Finished batches are kept for the next run.")
			}
			if progress.CompletedBatches == 0 {
				return
			}
		}
		if progress.TotalBatches == 0 {
			return
		}
		percent := progress.CompletedBatches * 100 / progress.TotalBatches
		if !style.interactive {
			elapsedBucket := int(progress.Elapsed / (15 * time.Second))
			completedChanged := progress.CompletedBatches != lastCompleted
			if progress.Phase == "running" && elapsedBucket == lastElapsedBucket {
				return
			}
			if progress.Phase == "batch_started" && lastCompleted >= 0 {
				return
			}
			if progress.Phase == "batch_completed" && progress.TotalBatches > 20 &&
				progress.CompletedBatches != progress.TotalBatches && percent == lastPercent && elapsedBucket == lastElapsedBucket {
				return
			}
			if !completedChanged && progress.Phase == "prepared" {
				return
			}
			lastElapsedBucket = elapsedBucket
		}
		lastCompleted, lastPercent = progress.CompletedBatches, percent
		spinner := []string{"◐", "◓", "◑", "◒"}[int(progress.Elapsed/time.Second)%4]
		if progress.CompletedBatches == progress.TotalBatches {
			spinner = "✓"
		} else if progress.Phase == "batch_failed" {
			spinner = "!"
		}
		line := fmt.Sprintf("      %s %s %3d%%  %d/%d batches  %d/%d rules",
			spinner, localCheckBar(percent, 10), percent, progress.CompletedBatches, progress.TotalBatches,
			progress.CompletedControls, progress.TotalControls)
		if progress.CompletedBatches < progress.TotalBatches {
			line += fmt.Sprintf("  %d jobs  ≤%d requested", progress.ActiveBatches, progress.ActiveAgentSlots)
		}
		line += "  " + formatProgressDuration(progress.Elapsed)
		if eta, ok := reviewETA(progress); ok {
			line += "  ETA " + formatProgressDuration(eta)
		}
		if progress.TokenUsageBatches > 0 {
			line += fmt.Sprintf("  tok %s in/%s out",
				formatCompactCount(progress.TokenUsage.InputTokens),
				formatCompactCount(progress.TokenUsage.OutputTokens))
		}
		if progress.EstimatedCostBatches > 0 {
			line += fmt.Sprintf("  est. $%.4f", progress.EstimatedCostUSD)
		}
		colored := style.paint(ansiBlue, line)
		if style.interactive {
			fmt.Fprintf(output, "\r\x1b[2K%s", colored)
			liveLine = true
			if progress.CompletedBatches == progress.TotalBatches || progress.Phase == "batch_failed" {
				fmt.Fprintln(output)
				liveLine = false
			}
			return
		}
		if liveLine {
			fmt.Fprintln(output)
			liveLine = false
		}
		fmt.Fprintln(output, colored)
	}
}

func reviewETA(progress controlreview.Progress) (time.Duration, bool) {
	completed := progress.CompletedBatches - progress.ReusedBatches
	remaining := progress.TotalBatches - progress.CompletedBatches
	if completed < 1 || remaining < 1 || progress.Elapsed <= 0 {
		return 0, false
	}
	return time.Duration(float64(progress.Elapsed) * float64(remaining) / float64(completed)), true
}

func formatProgressDuration(value time.Duration) string {
	if value < 0 {
		value = 0
	}
	value = value.Round(time.Second)
	hours := int(value / time.Hour)
	minutes := int((value % time.Hour) / time.Minute)
	seconds := int((value % time.Minute) / time.Second)
	if hours > 0 {
		return fmt.Sprintf("%dh %02dm", hours, minutes)
	}
	if minutes > 0 {
		return fmt.Sprintf("%dm %02ds", minutes, seconds)
	}
	return fmt.Sprintf("%ds", seconds)
}

func formatCompactCount(value int64) string {
	switch {
	case value >= 1_000_000_000:
		return fmt.Sprintf("%.1fB", float64(value)/1_000_000_000)
	case value >= 1_000_000:
		return fmt.Sprintf("%.1fM", float64(value)/1_000_000)
	case value >= 1_000:
		return fmt.Sprintf("%.1fk", float64(value)/1_000)
	default:
		return strconv.FormatInt(value, 10)
	}
}

func countedNoun(count int, singular, plural string) string {
	label := plural
	if count == 1 {
		label = singular
	}
	return fmt.Sprintf("%d %s", count, label)
}

func printAIReviewResults(output io.Writer, run model.RunResult, summary controlreview.Summary, style terminalStyle, showAll bool) {
	reviewed := make([]model.ControlResult, 0)
	counts := map[string]int{}
	for _, result := range run.ControlResults {
		if result.AIReview == nil {
			continue
		}
		reviewed = append(reviewed, result)
		counts[result.AIReview.AssessmentCandidate]++
	}
	if len(reviewed) == 0 {
		return
	}
	sort.SliceStable(reviewed, func(left, right int) bool {
		leftReview, rightReview := reviewed[left].AIReview, reviewed[right].AIReview
		if leftRank, rightRank := scanSeverityRank(leftReview.Priority), scanSeverityRank(rightReview.Priority); leftRank != rightRank {
			return leftRank < rightRank
		}
		if leftRank, rightRank := aiCandidateRank(leftReview.AssessmentCandidate), aiCandidateRank(rightReview.AssessmentCandidate); leftRank != rightRank {
			return leftRank < rightRank
		}
		return reviewed[left].ControlID < reviewed[right].ControlID
	})

	printTerminalSection(output, style, "AI REVIEW")
	if summary.Provider != "" {
		fmt.Fprintf(output, "  %s · %s review · %s · %s completed · %s cached · %s\n",
			authenticationProviderTitle(summary.Provider), terminalText(summary.ReviewDepth),
			countedNoun(summary.ReviewedControls, "rule", "rules"),
			countedNoun(summary.CompletedBatches, "batch", "batches"),
			countedNoun(summary.ReusedBatches, "batch", "batches"), formatProgressDuration(summary.Duration))
	}
	fmt.Fprintf(output, "  %s %d work · %s %d proof · %s %d okay · %s %d may not apply\n",
		style.paint(ansiRed, "✗"), counts["advisory_fail_candidate"],
		style.paint(ansiYellow, "?"), counts["needs_evidence"],
		style.paint(ansiGreen, "✓"), counts["advisory_pass_candidate"],
		style.paint(ansiBlue, "–"), counts["not_applicable_candidate"])
	if summary.TokenUsageBatches > 0 {
		fmt.Fprintf(output, "  Tokens  %s in (%s cached) · %s out · %s thinking\n",
			formatCompactCount(summary.TokenUsage.InputTokens), formatCompactCount(summary.TokenUsage.CachedInputTokens),
			formatCompactCount(summary.TokenUsage.OutputTokens), formatCompactCount(summary.TokenUsage.ReasoningOutputTokens))
	} else if summary.Provider == "codex" && summary.CompletedBatches > summary.ReusedBatches {
		fmt.Fprintln(output, "  Tokens  not reported by the completed Codex events")
	}
	if summary.EstimatedCostBatches > 0 {
		fmt.Fprintf(output, "  Cost    $%.4f provider estimate · not a bill\n", summary.EstimatedCostUSD)
	}
	if summary.MaxCostUSD > 0 {
		fmt.Fprintf(output, "  Limit   $%.4f for each new Claude batch\n", summary.MaxCostUSD)
	}
	fmt.Fprintln(output, "  AI advice only — it does not change verified pass or fail results.")

	visible := len(reviewed)
	if !showAll && visible > 8 {
		visible = 8
	}
	fmt.Fprintln(output)
	for _, result := range reviewed[:visible] {
		review := result.AIReview
		priority := strings.ToUpper(review.Priority)
		if priority == "" || priority == "NONE" {
			priority = "INFO"
		}
		fmt.Fprintf(output, "  %s  %-8s %s  %s\n",
			aiCandidateLabel(review.AssessmentCandidate, style),
			style.paint(terminalSeverityColor(review.Priority), priority),
			terminalText(result.ControlID), terminalBrief(review.Advice, 96))
	}
	if hidden := len(reviewed) - visible; hidden > 0 {
		fmt.Fprintf(output, "  + %d more — open the report or run again with --details.\n", hidden)
	}
	if showAll && summary.StateDirectory != "" {
		fmt.Fprintf(output, "\n  Resume data  %s\n", terminalText(summary.StateDirectory))
	}
}

func printTerminalSection(output io.Writer, style terminalStyle, title string) {
	line := "── " + title + " " + strings.Repeat("─", max(2, 54-len(title)))
	fmt.Fprintln(output)
	fmt.Fprintln(output, style.paint(ansiBlue, line))
}

func aiCandidateLabel(candidate string, style terminalStyle) string {
	switch candidate {
	case "advisory_fail_candidate":
		return style.paint(ansiRed, "✗ WORK  ")
	case "needs_evidence":
		return style.paint(ansiYellow, "? PROOF ")
	case "advisory_pass_candidate":
		return style.paint(ansiGreen, "✓ OKAY  ")
	case "not_applicable_candidate":
		return style.paint(ansiBlue, "– N/A?  ")
	default:
		return style.paint(ansiRed, "× ERROR ")
	}
}

func aiCandidateRank(candidate string) int {
	switch candidate {
	case "advisory_fail_candidate":
		return 0
	case "needs_evidence":
		return 1
	case "not_applicable_candidate":
		return 2
	case "advisory_pass_candidate":
		return 3
	default:
		return 4
	}
}

func reorderInterspersedFlags(set *flag.FlagSet, args []string) []string {
	flags := make([]string, 0, len(args))
	positionals := make([]string, 0, 1)
	usedDelimiter := false
	for index := 0; index < len(args); index++ {
		argument := args[index]
		if argument == "--" {
			positionals = append(positionals, args[index+1:]...)
			usedDelimiter = true
			break
		}
		if len(argument) > 1 && argument[0] == '-' {
			flags = append(flags, argument)
			name := strings.TrimLeft(argument, "-")
			if separator := strings.IndexByte(name, '='); separator >= 0 {
				continue
			}
			defined := set.Lookup(name)
			boolean := false
			if defined != nil {
				if value, ok := defined.Value.(interface{ IsBoolFlag() bool }); ok {
					boolean = value.IsBoolFlag()
				}
			}
			if defined != nil && !boolean && index+1 < len(args) {
				index++
				flags = append(flags, args[index])
			}
			continue
		}
		positionals = append(positionals, argument)
	}
	if usedDelimiter {
		flags = append(flags, "--")
	}
	return append(flags, positionals...)
}

func printScanSummary(
	output io.Writer, run model.RunResult, reviewSummary controlreview.Summary,
	style terminalStyle, writtenReport string, showLocalDetails, showReviewDetails bool,
) {
	local := report.SummarizeLocalChecks(run)
	tone := terminalToneColor(local.Tone)
	fmt.Fprintln(output, style.paint(tone, "  ╭─ SCAN RESULT"))
	fmt.Fprintf(output, "  │ %s · %d%%\n", style.paint(tone, local.Label), local.Percentage)
	fmt.Fprintf(output, "  │ %s %d/%d applicable checks passed\n",
		style.paint(tone, localCheckBar(local.Percentage, 16)), local.Passed, local.Applicable)
	fmt.Fprintf(output, "  │ %d passed · %d failed · %d unresolved · %d manual · %d not needed\n",
		local.Passed, local.Failed, local.Unresolved, local.Manual, local.NotApplicable)
	fmt.Fprintf(output, "  ╰─ %s\n", local.Explanation)
	attention := sortedScanAttention(run.Results)
	printTerminalSection(output, style, "LOCAL CHECKS")
	if showLocalDetails {
		fmt.Fprintf(output, "  Every local check · %d total\n\n", len(run.Results))
		for _, result := range run.Results {
			fmt.Fprintf(output, "  %s  %s  %s\n", assessmentLabel(result.Assessment, result.Execution, style),
				terminalText(result.AssertionID), terminalBrief(result.Summary, 108))
		}
	} else if len(attention) == 0 {
		fmt.Fprintf(output, "  %s All applicable local checks passed.\n", style.paint(ansiGreen, "✓"))
	} else {
		visible := len(attention)
		if visible > 6 {
			visible = 6
		}
		fmt.Fprintf(output, "  %s · showing the highest priority\n\n", countedNoun(len(attention), "check needs attention", "checks need attention"))
		for _, result := range attention[:visible] {
			severity := strings.ToUpper(result.Severity)
			if severity == "" {
				severity = "UNSET"
			}
			fmt.Fprintf(output, "  %s  %-8s %s  %s\n",
				assessmentLabel(result.Assessment, result.Execution, style),
				style.paint(terminalSeverityColor(result.Severity), severity),
				terminalText(result.AssertionID), terminalBrief(result.Summary, 96))
		}
		if remaining := len(attention) - visible; remaining > 0 {
			fmt.Fprintf(output, "  + %d more — open the report or run again with --details.\n", remaining)
		}
	}
	for _, result := range run.Results {
		if result.AssertionID == "PRC-A-CORE-013" && result.Execution == "blocked" {
			fmt.Fprintln(output, "  Next  Run `prc verify` for the pinned secret check; its image must already be local.")
			break
		}
	}

	printAIReviewResults(output, run, reviewSummary, style, showReviewDetails)

	printTerminalSection(output, style, "COVERAGE")
	fmt.Fprintf(output, "  Local checks  %d total · %d applicable · %d not needed\n",
		len(run.Results), local.Applicable, local.NotApplicable)
	if run.ControlCatalog != nil {
		controlCounts := map[string]int{}
		for _, result := range run.ControlResults {
			controlCounts[result.Disposition]++
		}
		fmt.Fprintf(output, "  Full catalog  %d/%d included · %d need evidence or review\n",
			len(run.ControlResults), run.ControlCatalog.ControlCount, controlCounts["needs_review"]+controlCounts["blocked"])
		if run.ControlCatalog.AIReviewedCount > 0 {
			fmt.Fprintf(output, "  AI review     %s · never counted as verified passes\n",
				countedNoun(run.ControlCatalog.AIReviewedCount, "advisory result", "advisory results"))
		}
		if run.AIImprovementPlan != nil {
			fmt.Fprintf(output, "  Action plan   %s · estimates remain advisory\n",
				countedNoun(run.AIImprovementPlan.ItemCount, "exact cause group", "exact cause groups"))
		}
		if len(run.AuthoritativeEvidence) > 0 {
			entries := 0
			for _, verification := range run.AuthoritativeEvidence {
				entries += verification.EntryCount
			}
			fmt.Fprintf(output, "  Evidence      %d exact entries from %d signed bundle(s)\n", entries, len(run.AuthoritativeEvidence))
		}
	}
	if excluded := inventoryFactCount(run.Inventory, "repository.user_exclusion"); excluded > 0 {
		label := "directories"
		if excluded == 1 {
			label = "directory"
		}
		fmt.Fprintf(output, "  Scope warning     %d reviewed .prcignore %s omitted and bound into this inventory\n", excluded, label)
	}

	printTerminalSection(output, style, "REPORT")
	if writtenReport == "" {
		fmt.Fprintln(output, "  Detailed report disabled")
	} else {
		reportLink, clickable := style.fileLink(writtenReport)
		fmt.Fprintf(output, "  %s\n", reportLink)
		if clickable {
			fmt.Fprintln(output, "  Click to open every check, finding, category score, and next step.")
		} else {
			fmt.Fprintln(output, "  Open it for every check, finding, category score, and next step.")
		}
	}
	fmt.Fprintln(output, "  Read-only scan · no fixes applied · no project scripts run")
	fmt.Fprintln(output)
}

func sortedScanAttention(results []model.AssertionResult) []model.AssertionResult {
	attention := make([]model.AssertionResult, 0, len(results))
	for _, result := range results {
		if result.Assessment == "pass" || result.Assessment == "not_applicable" {
			continue
		}
		attention = append(attention, result)
	}
	sort.SliceStable(attention, func(left, right int) bool {
		if leftRank, rightRank := scanSeverityRank(attention[left].Severity), scanSeverityRank(attention[right].Severity); leftRank != rightRank {
			return leftRank < rightRank
		}
		if leftRank, rightRank := scanAttentionRank(attention[left]), scanAttentionRank(attention[right]); leftRank != rightRank {
			return leftRank < rightRank
		}
		return attention[left].AssertionID < attention[right].AssertionID
	})
	return attention
}

func scanSeverityRank(severity string) int {
	switch severity {
	case "critical":
		return 0
	case "high":
		return 1
	case "medium":
		return 2
	case "low":
		return 3
	case "info":
		return 4
	default:
		return 5
	}
}

func scanAttentionRank(result model.AssertionResult) int {
	if result.Assessment == "fail" {
		return 0
	}
	if result.Execution == "blocked" || result.Execution == "error" || result.Assessment == "unknown" {
		return 1
	}
	if result.Assessment == "manual_review" {
		return 2
	}
	return 3
}

func writeScanReport(run model.RunResult, requestedPath string) (string, error) {
	path := requestedPath
	defaultReport := requestedPath == ""
	directory := ""
	if path == "" {
		cacheRoot, err := userCacheDirectory()
		if err != nil {
			return "", fmt.Errorf("locate the user cache for the scan report: %w", err)
		}
		directory = filepath.Join(cacheRoot, "prc", "reports")
		if pathWithin(run.Inventory.Root, directory) {
			return "", fmt.Errorf("default report directory is inside the scanned project; use --report with a path outside the project or --no-report")
		}
		if err := os.MkdirAll(directory, 0o700); err != nil {
			return "", fmt.Errorf("create private scan report directory: %w", err)
		}
		info, err := os.Lstat(directory)
		if err != nil || !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
			return "", fmt.Errorf("scan report directory is not a regular directory")
		}
		if err := os.Chmod(directory, 0o700); err != nil {
			return "", fmt.Errorf("protect scan report directory: %w", err)
		}
		runID := run.RunID
		if len(runID) > 16 {
			runID = runID[:16]
		}
		path = filepath.Join(directory, sanitizeReportName(run.Inventory.TargetName)+"-"+runID+".html")
	} else {
		absolute, err := filepath.Abs(path)
		if err != nil {
			return "", fmt.Errorf("resolve scan report path: %w", err)
		}
		path = absolute
		if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
			return "", fmt.Errorf("create scan report parent: %w", err)
		}
	}
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		return "", fmt.Errorf("create scan report without overwriting an existing file: %w", err)
	}
	if err := report.Write("html", file, run); err != nil {
		_ = file.Close()
		_ = os.Remove(path)
		return "", fmt.Errorf("write detailed scan report: %w", err)
	}
	if err := file.Close(); err != nil {
		_ = os.Remove(path)
		return "", fmt.Errorf("finish detailed scan report: %w", err)
	}
	absolute, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("resolve completed scan report path: %w", err)
	}
	if defaultReport {
		_, _, _ = pruneDefaultReports(directory, defaultReportRetention, absolute)
	}
	return absolute, nil
}

type cachedReport struct {
	path     string
	name     string
	modified time.Time
	size     int64
}

func pruneDefaultReports(directory string, keep int, preserve string) (int, int64, error) {
	if keep < 1 {
		return 0, 0, fmt.Errorf("default report retention must keep at least one report")
	}
	entries, err := os.ReadDir(directory)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, 0, nil
		}
		return 0, 0, err
	}
	reports := make([]cachedReport, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || entry.Type()&os.ModeSymlink != 0 || !scannerDefaultReportName(entry.Name()) {
			continue
		}
		info, infoErr := entry.Info()
		if infoErr != nil || !info.Mode().IsRegular() {
			continue
		}
		reports = append(reports, cachedReport{
			path: filepath.Join(directory, entry.Name()), name: entry.Name(),
			modified: info.ModTime(), size: info.Size(),
		})
	}
	preserved := filepath.Clean(preserve)
	sort.SliceStable(reports, func(left, right int) bool {
		leftPreserved := filepath.Clean(reports[left].path) == preserved
		rightPreserved := filepath.Clean(reports[right].path) == preserved
		if leftPreserved != rightPreserved {
			return leftPreserved
		}
		if !reports[left].modified.Equal(reports[right].modified) {
			return reports[left].modified.After(reports[right].modified)
		}
		return reports[left].name < reports[right].name
	})
	removed := 0
	var freed int64
	for _, item := range reports[min(keep, len(reports)):] {
		if err := os.Remove(item.path); err != nil {
			return removed, freed, err
		}
		removed++
		freed += item.size
	}
	return removed, freed, nil
}

func scannerDefaultReportName(name string) bool {
	if filepath.Ext(name) != ".html" {
		return false
	}
	base := strings.TrimSuffix(name, ".html")
	separator := strings.LastIndexByte(base, '-')
	if separator < 1 || len(base)-separator-1 != 16 {
		return false
	}
	for _, character := range base[separator+1:] {
		if character < '0' || character > '9' {
			if character < 'a' || character > 'f' {
				return false
			}
		}
	}
	return true
}

func pathWithin(root, path string) bool {
	if root == "" {
		return false
	}
	relative, err := filepath.Rel(filepath.Clean(root), filepath.Clean(path))
	return err == nil && relative != ".." && !strings.HasPrefix(relative, ".."+string(filepath.Separator))
}

func sanitizeReportName(value string) string {
	var result strings.Builder
	for _, character := range value {
		if character >= 'a' && character <= 'z' || character >= 'A' && character <= 'Z' ||
			character >= '0' && character <= '9' || character == '-' || character == '_' {
			result.WriteRune(character)
		} else if result.Len() > 0 && result.String()[result.Len()-1] != '-' {
			result.WriteByte('-')
		}
		if result.Len() >= 48 {
			break
		}
	}
	name := strings.Trim(result.String(), "-")
	if name == "" {
		return "project"
	}
	return name
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

func profileTerminalState(run model.RunResult) string {
	if run.ControlCatalog != nil && run.ControlCatalog.ProfileTerminalState != "" {
		return run.ControlCatalog.ProfileTerminalState
	}
	return run.TerminalState
}

func scanNoGoExitCode(terminalState string) int {
	switch terminalState {
	case "no_go":
		return exitGateFailed
	case "policy_stopped", "budget_exhausted":
		return exitPolicyDenied
	default:
		return exitSuccess
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

func inventoryFact(item model.Inventory, key string) string {
	for _, fact := range item.Facts {
		if fact.Key == key {
			return fact.Value
		}
	}
	return "not available"
}

func inventoryFactCount(item model.Inventory, key string) int {
	count := 0
	for _, fact := range item.Facts {
		if fact.Key == key {
			count++
		}
	}
	return count
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
