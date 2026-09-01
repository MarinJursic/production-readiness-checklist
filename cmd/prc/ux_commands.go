package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/MarinJursic/production-readiness-checklist/scanner/doctor"
)

const maximumUpdateResponse = 1024 * 1024

var (
	npmPackageRegistryURL = "https://registry.npmjs.org/@marinjursic%2fprc"
	updateHTTPClient      = &http.Client{
		Timeout: 8 * time.Second,
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	launchReport = defaultLaunchReport
)

func runSetup(args []string, stdout, stderr io.Writer) (int, error) {
	set := flag.NewFlagSet("setup", flag.ContinueOnError)
	set.SetOutput(stderr)
	set.Usage = func() {
		fmt.Fprintln(stderr, "Usage: prc setup [project path]")
		fmt.Fprintln(stderr, "Checks the project and bundled catalog without running project code or AI.")
	}
	if err := set.Parse(args); err != nil {
		return exitInternal, exitError(exitConfiguration, err)
	}
	if set.NArg() > 1 {
		return exitInternal, exitError(exitConfiguration, errors.New("setup accepts at most one project path"))
	}
	target := "."
	if set.NArg() == 1 {
		target = set.Arg(0)
	}
	result := doctor.Run(doctor.Options{Target: target, CatalogRoot: defaultCatalogRoot()})
	fmt.Fprintf(stdout, "Production Readiness Checklist %s\n\n", terminalText(version))
	fmt.Fprintf(stdout, "Project: %s\n", result.Target)
	fmt.Fprintf(stdout, "Bundled rules: %s\n", result.CatalogRoot)
	for _, check := range result.Checks {
		if check.Required {
			mark := "✓"
			if check.Status != "pass" {
				mark = "×"
			}
			fmt.Fprintf(stdout, "%s %s — %s\n", mark, check.ID, check.Summary)
		}
	}
	fmt.Fprintln(stdout)
	for _, providerName := range []string{"codex", "claude"} {
		if _, err := exec.LookPath(providerName); err == nil {
			fmt.Fprintf(stdout, "Optional AI: %s is installed; check login with `prc auth %s`.\n",
				authenticationProviderTitle(providerName), providerName)
		} else {
			fmt.Fprintf(stdout, "Optional AI: %s is not installed; local scanning still works.\n",
				authenticationProviderTitle(providerName))
		}
	}
	fmt.Fprintln(stdout)
	if !result.Ready {
		fmt.Fprintln(stdout, "Setup is incomplete. Fix the failed item above, then run `prc setup` again.")
		return exitIncomplete, nil
	}
	fmt.Fprintln(stdout, "Ready. Run `prc` for the normal local scan. No AI or fixes will start.")
	return exitSuccess, nil
}

type npmRegistryDocument struct {
	DistTags map[string]string `json:"dist-tags"`
}

func runUpdate(args []string, stdout, stderr io.Writer) error {
	set := flag.NewFlagSet("update", flag.ContinueOnError)
	set.SetOutput(stderr)
	set.Usage = func() {
		fmt.Fprintln(stderr, "Usage: prc update")
		fmt.Fprintln(stderr, "Checks the official npm registry once and prints an update command; never installs automatically.")
	}
	if err := set.Parse(args); err != nil {
		return exitError(exitConfiguration, err)
	}
	if set.NArg() != 0 {
		return exitError(exitConfiguration, errors.New("update accepts no positional arguments"))
	}
	request, err := http.NewRequestWithContext(commandContext, http.MethodGet, npmPackageRegistryURL, nil)
	if err != nil {
		return exitError(exitInternal, err)
	}
	request.Header.Set("Accept", "application/json")
	request.Header.Set("User-Agent", "prc/"+version+" explicit-update-check")
	response, err := updateHTTPClient.Do(request)
	if err != nil {
		return exitError(exitExecution, fmt.Errorf("check the npm registry: %w", err))
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return exitError(exitExecution, fmt.Errorf("npm registry returned HTTP %d", response.StatusCode))
	}
	reader := io.LimitReader(response.Body, maximumUpdateResponse+1)
	data, err := io.ReadAll(reader)
	if err != nil {
		return exitError(exitExecution, fmt.Errorf("read npm registry response: %w", err))
	}
	if len(data) > maximumUpdateResponse {
		return exitError(exitExecution, errors.New("npm registry response exceeded the 1 MiB safety limit"))
	}
	var document npmRegistryDocument
	if err := json.Unmarshal(data, &document); err != nil {
		return exitError(exitExecution, fmt.Errorf("decode npm registry response: %w", err))
	}
	latest := document.DistTags["latest"]
	if !validReleaseVersion(latest) {
		return exitError(exitExecution, errors.New("npm registry did not return a valid latest release version"))
	}
	fmt.Fprintf(stdout, "Installed: %s\nLatest:    %s\n", version, latest)
	if version == latest {
		fmt.Fprintln(stdout, "You are up to date.")
		return nil
	}
	if !validReleaseVersion(version) {
		fmt.Fprintln(stdout, "This is a development build, so it is not compared with the published release.")
		return nil
	}
	if compareReleaseVersions(version, latest) > 0 {
		fmt.Fprintln(stdout, "This installed build is newer than the current npm latest tag.")
		return nil
	}
	fmt.Fprintln(stdout, "A newer version is available.")
	fmt.Fprintln(stdout, "Update with: npm install -g @marinjursic/prc@latest")
	return nil
}

func compareReleaseVersions(left, right string) int {
	leftParts, rightParts := strings.Split(left, "."), strings.Split(right, ".")
	for index := 0; index < 3; index++ {
		leftValue, _ := strconv.Atoi(leftParts[index])
		rightValue, _ := strconv.Atoi(rightParts[index])
		if leftValue < rightValue {
			return -1
		}
		if leftValue > rightValue {
			return 1
		}
	}
	return 0
}

func validReleaseVersion(value string) bool {
	parts := strings.Split(value, ".")
	if len(parts) != 3 {
		return false
	}
	for _, part := range parts {
		if part == "" || len(part) > 9 || len(part) > 1 && part[0] == '0' {
			return false
		}
		for _, character := range part {
			if character < '0' || character > '9' {
				return false
			}
		}
	}
	return true
}

func runReport(args []string, stdout, stderr io.Writer) error {
	action := "open"
	if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		action, args = args[0], args[1:]
	}
	set := flag.NewFlagSet("report "+action, flag.ContinueOnError)
	set.SetOutput(stderr)
	set.Usage = func() {
		fmt.Fprintln(stderr, "Usage: prc report [open|path|list] [--limit N]")
		fmt.Fprintln(stderr, "Opens, locates, or lists private scanner-generated HTML reports.")
	}
	limit := set.Int("limit", 10, "maximum reports to list")
	if err := set.Parse(args); err != nil {
		return exitError(exitConfiguration, err)
	}
	if set.NArg() != 0 {
		return exitError(exitConfiguration, fmt.Errorf("unexpected report arguments: %s", strings.Join(set.Args(), " ")))
	}
	reports, err := cachedReports()
	if err != nil {
		return exitError(exitConfiguration, err)
	}
	if action == "list" {
		if *limit < 1 || *limit > 1000 {
			return exitError(exitConfiguration, errors.New("--limit must be between 1 and 1000"))
		}
		if len(reports) == 0 {
			fmt.Fprintln(stdout, "No cached scan reports were found. Run `prc` first.")
			return nil
		}
		for _, item := range reports[:min(*limit, len(reports))] {
			fmt.Fprintf(stdout, "%s  %s  %s\n", item.modified.Format(time.RFC3339), formatByteCount(item.size), item.path)
		}
		return nil
	}
	if flagWasSet(set, "limit") {
		return exitError(exitConfiguration, errors.New("--limit applies only to `prc report list`"))
	}
	if action != "open" && action != "path" {
		return exitError(exitConfiguration, fmt.Errorf("unknown report action %q; use open, path, or list", action))
	}
	if len(reports) == 0 {
		return exitError(exitConfiguration, errors.New("no cached scan report was found; run `prc` first"))
	}
	latest := reports[0].path
	if action == "path" {
		fmt.Fprintln(stdout, latest)
		return nil
	}
	if err := launchReport(latest); err != nil {
		return exitError(exitExecution, fmt.Errorf("open the latest report: %w; open it manually at %s", err, latest))
	}
	fmt.Fprintf(stdout, "Opened: %s\n", latest)
	return nil
}

func cachedReports() ([]cachedReport, error) {
	root, err := userCacheDirectory()
	if err != nil {
		return nil, fmt.Errorf("locate scanner cache: %w", err)
	}
	directory := filepath.Join(root, "prc", "reports")
	entries, err := os.ReadDir(directory)
	if os.IsNotExist(err) {
		return []cachedReport{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read report cache: %w", err)
	}
	reports := make([]cachedReport, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || entry.Type()&os.ModeSymlink != 0 || !scannerDefaultReportName(entry.Name()) {
			continue
		}
		info, err := entry.Info()
		if err != nil || !info.Mode().IsRegular() {
			continue
		}
		reports = append(reports, cachedReport{
			path: filepath.Join(directory, entry.Name()), name: entry.Name(), modified: info.ModTime(), size: info.Size(),
		})
	}
	sort.SliceStable(reports, func(left, right int) bool {
		if !reports[left].modified.Equal(reports[right].modified) {
			return reports[left].modified.After(reports[right].modified)
		}
		return reports[left].name < reports[right].name
	})
	return reports, nil
}

func defaultLaunchReport(path string) error {
	var command *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		command = exec.Command("open", path)
	case "windows":
		command = exec.Command("rundll32", "url.dll,FileProtocolHandler", path)
	default:
		command = exec.Command("xdg-open", path)
	}
	if err := command.Start(); err != nil {
		return err
	}
	return command.Process.Release()
}

type cacheSummary struct {
	Files int
	Bytes int64
}

func runCache(args []string, stdout, stderr io.Writer) error {
	if len(args) > 0 && (args[0] == "help" || args[0] == "--help" || args[0] == "-h") {
		fmt.Fprintln(stdout, "Usage:")
		fmt.Fprintln(stdout, "  prc cache status")
		fmt.Fprintln(stdout, "  prc cache clean --reports|--ai|--all [--older-than DURATION]")
		fmt.Fprintln(stdout, "Only scanner-owned cache paths are inspected or removed; a data class is always required for deletion.")
		return nil
	}
	if len(args) == 0 || args[0] == "status" {
		if len(args) > 1 {
			return exitError(exitConfiguration, errors.New("cache status accepts no options"))
		}
		return printCacheStatus(stdout)
	}
	if args[0] != "clean" {
		return exitError(exitConfiguration, errors.New("cache requires status or clean"))
	}
	set := flag.NewFlagSet("cache clean", flag.ContinueOnError)
	set.SetOutput(stderr)
	reports := set.Bool("reports", false, "remove scanner-generated HTML reports")
	ai := set.Bool("ai", false, "remove resumable AI review cache entries")
	all := set.Bool("all", false, "remove reports and AI review cache entries")
	olderThan := set.Duration("older-than", 0, "only remove files older than this duration")
	if err := set.Parse(args[1:]); err != nil {
		return exitError(exitConfiguration, err)
	}
	if set.NArg() != 0 {
		return exitError(exitConfiguration, fmt.Errorf("unexpected cache clean arguments: %s", strings.Join(set.Args(), " ")))
	}
	if *olderThan < 0 {
		return exitError(exitConfiguration, errors.New("--older-than cannot be negative"))
	}
	if !*reports && !*ai && !*all {
		return exitError(exitConfiguration, errors.New("cache clean requires --reports, --ai, or --all"))
	}
	root, err := userCacheDirectory()
	if err != nil {
		return exitError(exitInternal, err)
	}
	threshold := time.Time{}
	if *olderThan > 0 {
		threshold = time.Now().Add(-*olderThan)
	}
	total := cacheSummary{}
	if *reports || *all {
		removed, err := cleanCacheTree(filepath.Join(root, "prc", "reports"), threshold, true)
		if err != nil {
			return exitError(exitExecution, err)
		}
		total.Files += removed.Files
		total.Bytes += removed.Bytes
	}
	if *ai || *all {
		removed, err := cleanCacheTree(filepath.Join(root, "prc", "control-reviews"), threshold, false)
		if err != nil {
			return exitError(exitExecution, err)
		}
		total.Files += removed.Files
		total.Bytes += removed.Bytes
	}
	fmt.Fprintf(stdout, "Removed %d scanner cache files and freed %s.\n", total.Files, formatByteCount(total.Bytes))
	return nil
}

func printCacheStatus(output io.Writer) error {
	root, err := userCacheDirectory()
	if err != nil {
		return err
	}
	items := []struct {
		label string
		path  string
	}{
		{"Reports", filepath.Join(root, "prc", "reports")},
		{"AI resume data", filepath.Join(root, "prc", "control-reviews")},
	}
	total := cacheSummary{}
	for _, item := range items {
		summary, err := inspectCacheTree(item.path)
		if err != nil {
			return err
		}
		fmt.Fprintf(output, "%-15s %7s in %d file(s)  %s\n", item.label+":", formatByteCount(summary.Bytes), summary.Files, item.path)
		total.Files += summary.Files
		total.Bytes += summary.Bytes
	}
	fmt.Fprintf(output, "%-15s %7s in %d file(s)\n", "Total:", formatByteCount(total.Bytes), total.Files)
	fmt.Fprintln(output, "Nothing was deleted. Use `prc cache clean --reports` or `prc cache clean --ai` explicitly.")
	return nil
}

func inspectCacheTree(root string) (cacheSummary, error) {
	result := cacheSummary{}
	info, err := os.Lstat(root)
	if os.IsNotExist(err) {
		return result, nil
	}
	if err != nil {
		return result, fmt.Errorf("inspect cache root %s: %w", root, err)
	}
	if !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
		return result, fmt.Errorf("scanner cache root is not a regular directory: %s", root)
	}
	err = filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("scanner cache contains a symbolic link: %s", path)
		}
		if entry.IsDir() {
			return nil
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return fmt.Errorf("scanner cache contains a non-regular file: %s", path)
		}
		result.Files++
		result.Bytes += info.Size()
		return nil
	})
	return result, err
}

func cleanCacheTree(root string, olderThan time.Time, reportsOnly bool) (cacheSummary, error) {
	result := cacheSummary{}
	if _, err := inspectCacheTree(root); err != nil {
		return result, err
	}
	entries := []string{}
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if os.IsNotExist(walkErr) {
			return nil
		}
		if walkErr != nil {
			return walkErr
		}
		entries = append(entries, path)
		return nil
	})
	if os.IsNotExist(err) {
		return result, nil
	}
	if err != nil {
		return result, err
	}
	sort.Slice(entries, func(left, right int) bool { return len(entries[left]) > len(entries[right]) })
	for _, path := range entries {
		if path == root {
			continue
		}
		info, err := os.Lstat(path)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return result, err
		}
		if info.IsDir() {
			_ = os.Remove(path)
			continue
		}
		if reportsOnly && !scannerDefaultReportName(filepath.Base(path)) {
			continue
		}
		if !olderThan.IsZero() && !info.ModTime().Before(olderThan) {
			continue
		}
		if err := os.Remove(path); err != nil {
			return result, fmt.Errorf("remove scanner cache file %s: %w", path, err)
		}
		result.Files++
		result.Bytes += info.Size()
	}
	return result, nil
}

func formatByteCount(value int64) string {
	const unit = int64(1024)
	if value < unit {
		return strconv.FormatInt(value, 10) + " B"
	}
	divisor, exponent := unit, 0
	for next := value / unit; next >= unit && exponent < 3; next /= unit {
		divisor *= unit
		exponent++
	}
	return fmt.Sprintf("%.1f %ciB", float64(value)/float64(divisor), "KMGT"[exponent])
}

func optionalModel(value string) string {
	if strings.TrimSpace(value) == "" {
		return ""
	}
	return " (" + terminalText(value) + ")"
}

func runCompletion(args []string, stdout, stderr io.Writer) error {
	set := flag.NewFlagSet("completion", flag.ContinueOnError)
	set.SetOutput(stderr)
	set.Usage = func() {
		fmt.Fprintln(stderr, "Usage: prc completion <zsh|bash|fish|powershell>")
	}
	if err := set.Parse(args); err != nil {
		return exitError(exitConfiguration, err)
	}
	if set.NArg() != 1 {
		return exitError(exitConfiguration, errors.New("usage: prc completion <zsh|bash|fish|powershell>"))
	}
	commands := "setup quick scan verify ci full report login logout auth update version doctor coverage inventory plan diff history cache completion help"
	switch set.Arg(0) {
	case "zsh":
		fmt.Fprintf(stdout, "#compdef prc\n_arguments '1:command:(%s)' '*::arg:->args'\n", commands)
	case "bash":
		fmt.Fprintf(stdout, "_prc_complete() { COMPREPLY=( $(compgen -W '%s' -- \"${COMP_WORDS[1]}\") ); }\ncomplete -F _prc_complete prc\n", commands)
	case "fish":
		for _, command := range strings.Fields(commands) {
			fmt.Fprintf(stdout, "complete -c prc -f -n '__fish_use_subcommand' -a %s\n", command)
		}
	case "powershell":
		fmt.Fprintf(stdout, "Register-ArgumentCompleter -Native -CommandName prc -ScriptBlock { param($wordToComplete) '%s'.Split(' ') | Where-Object { $_ -like \"$wordToComplete*\" } }\n", commands)
	default:
		return exitError(exitConfiguration, fmt.Errorf("unsupported shell %q; use zsh, bash, fish, or powershell", set.Arg(0)))
	}
	return nil
}
