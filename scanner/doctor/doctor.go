// Package doctor inspects local scanner prerequisites without executing target
// code, containers, or agent providers.
package doctor

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/MarinJursic/production-readiness-checklist/scanner/catalog"
	"github.com/MarinJursic/production-readiness-checklist/scanner/inventory"
)

const Schema = "prc.doctor/v0.1"

type Options struct {
	Target          string
	CatalogRoot     string
	StateDirectory  string
	CandidateParent string
	OCIRuntime      string
	Providers       []string
	Now             func() time.Time
}

type Check struct {
	ID       string   `json:"id"`
	Status   string   `json:"status"`
	Required bool     `json:"required"`
	Summary  string   `json:"summary"`
	Details  []string `json:"details"`
}

type Summary struct {
	Passed   int `json:"passed"`
	Warnings int `json:"warnings"`
	Failed   int `json:"failed"`
}

type Report struct {
	SchemaVersion string    `json:"schema_version"`
	GeneratedAt   time.Time `json:"generated_at"`
	Platform      string    `json:"platform"`
	Architecture  string    `json:"architecture"`
	Target        string    `json:"target"`
	CatalogRoot   string    `json:"catalog_root"`
	Ready         bool      `json:"ready"`
	Summary       Summary   `json:"summary"`
	Checks        []Check   `json:"checks"`
}

// Run performs read-only validation except for private, self-cleaning filesystem
// probes in explicitly supplied state and candidate directories.
func Run(options Options) Report {
	if options.Now == nil {
		options.Now = func() time.Time { return time.Now().UTC() }
	}
	report := Report{
		SchemaVersion: Schema, GeneratedAt: options.Now().UTC(), Platform: runtime.GOOS,
		Architecture: runtime.GOARCH, Checks: []Check{},
	}

	target, targetCheck := inspectTarget(options.Target)
	report.Target = target
	report.Checks = append(report.Checks, targetCheck)
	catalogRoot, catalogCheck := inspectCatalog(options.CatalogRoot)
	report.CatalogRoot = catalogRoot
	report.Checks = append(report.Checks, catalogCheck)
	report.Checks = append(report.Checks, inspectStateDirectory(options.StateDirectory))
	report.Checks = append(report.Checks, inspectCandidateParent(options.CandidateParent, target, targetCheck.Status == "pass"))
	report.Checks = append(report.Checks, inspectExecutable("oci-runtime", options.OCIRuntime, []string{"docker", "podman"}))

	providers := append([]string(nil), options.Providers...)
	sort.Strings(providers)
	for _, provider := range providers {
		providerName := strings.TrimSuffix(strings.ToLower(filepath.Base(provider)), ".exe")
		if providerName != "codex" && providerName != "claude" {
			report.Checks = append(report.Checks, failed("provider.invalid", "The requested provider is not supported.", fmt.Errorf("provider=%s", providerName)))
			continue
		}
		report.Checks = append(report.Checks, inspectExecutable("provider."+providerName, provider, []string{providerName}))
	}
	for _, check := range report.Checks {
		switch check.Status {
		case "pass":
			report.Summary.Passed++
		case "warn":
			report.Summary.Warnings++
		case "fail":
			report.Summary.Failed++
		}
	}
	report.Ready = report.Summary.Failed == 0
	return report
}

func inspectTarget(value string) (string, Check) {
	if strings.TrimSpace(value) == "" {
		value = "."
	}
	absolute, err := filepath.Abs(value)
	if err != nil {
		return value, failed("target.inventory", "The target path could not be resolved.", err)
	}
	item, err := inventory.Build(absolute)
	if err != nil {
		return absolute, failed("target.inventory", "The target cannot be inventoried safely.", err)
	}
	return item.Root, passed("target.inventory", true, "The target can be inventoried without executing project code.",
		fmt.Sprintf("inventory_digest=%s", item.Digest), fmt.Sprintf("files=%d", item.FileCount))
}

func inspectCatalog(value string) (string, Check) {
	if strings.TrimSpace(value) == "" {
		value = "."
	}
	absolute, err := filepath.Abs(value)
	if err != nil {
		return value, failed("catalog.load", "The catalog root could not be resolved.", err)
	}
	c, err := catalog.Load(absolute)
	if err != nil {
		return absolute, failed("catalog.load", "The catalog is unavailable or invalid.", err)
	}
	return c.Root, passed("catalog.load", true, "The catalog and its references are valid.",
		fmt.Sprintf("objectives=%d", len(c.Objectives)), fmt.Sprintf("assertions=%d", len(c.Assertions)),
		fmt.Sprintf("profiles=%d", len(c.Profiles)))
}

func inspectStateDirectory(value string) Check {
	if strings.TrimSpace(value) == "" {
		return warning("state-store", "No state directory was requested; persistent evidence-store capabilities were not tested.")
	}
	root, err := existingDirectory(value)
	if err != nil {
		return failed("state-store", "The requested state directory is not usable.", err)
	}
	if runtime.GOOS != "windows" {
		info, statErr := os.Stat(root)
		if statErr != nil || info.Mode().Perm()&0o077 != 0 {
			return failed("state-store", "The state directory is accessible by group or other users.", fmt.Errorf("use mode 0700 or stricter"))
		}
	}
	probe, err := os.MkdirTemp(root, ".prc-doctor-")
	if err != nil {
		return failed("state-store", "The scanner cannot create private state records.", err)
	}
	defer os.RemoveAll(probe)
	if err := os.Chmod(probe, 0o700); err != nil {
		return failed("state-store", "The scanner cannot protect a state directory.", err)
	}
	source := filepath.Join(probe, "source")
	linked := filepath.Join(probe, "record")
	if err := writeSyncedProbe(source, []byte("prc-doctor\n"), 0o600); err != nil {
		return failed("state-store", "The scanner cannot write a private state record.", err)
	}
	if err := os.Link(source, linked); err != nil {
		return failed("state-store", "The filesystem does not support immutable evidence publication by hard link.", err)
	}
	return passed("state-store", true, "Private content-addressed evidence publication is supported.", "path="+root)
}

func inspectCandidateParent(value, target string, targetValid bool) Check {
	if strings.TrimSpace(value) == "" {
		return warning("candidate-workspace", "No candidate parent was requested; isolated remediation filesystem capabilities were not tested.")
	}
	parent, err := existingDirectory(value)
	if err != nil {
		return failed("candidate-workspace", "The requested candidate parent is not usable.", err)
	}
	if !targetValid {
		return failed("candidate-workspace", "Candidate isolation cannot be established until the target is valid.", nil)
	}
	if within(target, parent) {
		return failed("candidate-workspace", "The candidate parent must be outside the target tree.", nil)
	}
	probe, err := os.MkdirTemp(parent, ".prc-doctor-")
	if err != nil {
		return failed("candidate-workspace", "The scanner cannot create an isolated candidate directory.", err)
	}
	defer os.RemoveAll(probe)
	if within(target, probe) || within(probe, target) {
		return failed("candidate-workspace", "The probed candidate directory is not disjoint from the target tree.", nil)
	}
	if err := os.Chmod(probe, 0o700); err != nil {
		return failed("candidate-workspace", "The scanner cannot protect an isolated candidate directory.", err)
	}
	path := filepath.Join(probe, "mode-probe")
	if err := os.WriteFile(path, []byte("prc-doctor\n"), 0o666); err != nil {
		return failed("candidate-workspace", "The scanner cannot create candidate files.", err)
	}
	if err := os.Chmod(path, 0o644); err != nil {
		return failed("candidate-workspace", "The scanner cannot enforce candidate file modes.", err)
	}
	if runtime.GOOS != "windows" {
		info, statErr := os.Stat(path)
		if statErr != nil || info.Mode().Perm() != 0o644 {
			return failed("candidate-workspace", "The filesystem does not preserve required candidate file modes.", statErr)
		}
	}
	return passed("candidate-workspace", true, "Isolated candidate creation and mode enforcement are supported.", "path="+parent)
}

func inspectExecutable(id, value string, allowedNames []string) Check {
	if strings.TrimSpace(value) == "" {
		return warning(id, "This optional executable capability was not requested.")
	}
	path, err := exec.LookPath(value)
	if err != nil {
		return failed(id, "The requested executable was not found.", err)
	}
	path, err = filepath.Abs(path)
	if err != nil {
		return failed(id, "The executable path could not be resolved.", err)
	}
	if resolved, resolveErr := filepath.EvalSymlinks(path); resolveErr == nil {
		path = resolved
	}
	name := strings.TrimSuffix(strings.ToLower(filepath.Base(path)), ".exe")
	allowed := false
	for _, candidate := range allowedNames {
		if name == candidate {
			allowed = true
			break
		}
	}
	if !allowed {
		return failed(id, "The executable name is not allowed for this capability.", fmt.Errorf("resolved to %s", name))
	}
	info, err := os.Stat(path)
	if err != nil || !info.Mode().IsRegular() {
		return failed(id, "The executable is not a readable regular file.", err)
	}
	if runtime.GOOS != "windows" && info.Mode().Perm()&0o111 == 0 {
		return failed(id, "The resolved file is not executable.", nil)
	}
	digest, err := fileDigest(path)
	if err != nil {
		return failed(id, "The executable could not be hashed.", err)
	}
	return passed(id, true, "The requested executable is available and content-addressable.",
		"path="+path, "sha256="+digest)
}

func fileDigest(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	hasher := sha256.New()
	if _, err := io.Copy(hasher, io.LimitReader(file, 1024*1024*1024+1)); err != nil {
		return "", err
	}
	info, err := file.Stat()
	if err != nil {
		return "", err
	}
	if info.Size() > 1024*1024*1024 {
		return "", fmt.Errorf("executable exceeds the 1 GiB inspection limit")
	}
	return fmt.Sprintf("%x", hasher.Sum(nil)), nil
}

func writeSyncedProbe(path string, data []byte, mode os.FileMode) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, mode)
	if err != nil {
		return err
	}
	if _, err := file.Write(data); err != nil {
		file.Close()
		return err
	}
	if err := file.Sync(); err != nil {
		file.Close()
		return err
	}
	return file.Close()
}

func existingDirectory(path string) (string, error) {
	absolute, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	info, err := os.Stat(absolute)
	if err != nil || !info.IsDir() {
		return "", fmt.Errorf("path is not an accessible directory")
	}
	if resolved, resolveErr := filepath.EvalSymlinks(absolute); resolveErr == nil {
		absolute = resolved
	}
	return absolute, nil
}

func within(root, path string) bool {
	relative, err := filepath.Rel(filepath.Clean(root), filepath.Clean(path))
	return err == nil && relative != ".." && !strings.HasPrefix(relative, ".."+string(filepath.Separator))
}

func passed(id string, required bool, summary string, details ...string) Check {
	return Check{ID: id, Status: "pass", Required: required, Summary: summary, Details: append([]string{}, details...)}
}

func warning(id, summary string) Check {
	return Check{ID: id, Status: "warn", Required: false, Summary: summary, Details: []string{}}
}

func failed(id, summary string, err error) Check {
	details := []string{}
	if err != nil {
		details = append(details, err.Error())
	}
	return Check{ID: id, Status: "fail", Required: true, Summary: summary, Details: details}
}
