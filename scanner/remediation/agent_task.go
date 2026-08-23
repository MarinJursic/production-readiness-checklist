package remediation

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	workspaceinventory "github.com/MarinJursic/production-readiness-checklist/scanner/inventory"
	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
	"github.com/MarinJursic/production-readiness-checklist/scanner/provider"
	"github.com/MarinJursic/production-readiness-checklist/scanner/testdiscovery"
)

const agentTestSuiteAssertion = "PRC-A-CORE-010"

func planAgentTask(
	item model.Inventory,
	assertion model.Assertion,
	finding model.Finding,
	protectedPaths []string,
	options AgentOptions,
) (provider.Task, bool, error) {
	if assertion.ID != agentTestSuiteAssertion || assertion.ImplementationID != "prc.native.test-suite@0.2" ||
		assertion.RemediationClass != "R2" {
		return provider.Task{}, false, nil
	}
	if finding.AssertionID != assertion.ID || finding.Subject.InventoryDigest != item.Digest ||
		finding.ID == "" || finding.Fingerprint == "" {
		return provider.Task{}, false, fmt.Errorf("agent task requires the exact current finding for %s", assertion.ID)
	}
	if !options.AllowRemoteSourceProcessing {
		return provider.Task{}, false, policyDenied(fmt.Errorf("agent remediation requires explicit remote source processing acknowledgement"))
	}
	timeout := options.TimeoutSeconds
	if timeout == 0 {
		timeout = 300
	}
	maxOutput := options.MaxOutputBytes
	if maxOutput == 0 {
		maxOutput = 256 * 1024
	}

	records := append([]model.FileRecord(nil), item.Files...)
	sort.Slice(records, func(left, right int) bool { return records[left].Path < records[right].Path })
	existing := make(map[string]bool, len(records))
	for _, record := range records {
		existing[record.Path] = true
	}
	for _, record := range records {
		if record.Size > 256*1024 || !workspaceinventory.IsSourcePath(record.Path) ||
			testPath(record.Path) || protected(record.Path, protectedPaths) {
			continue
		}
		allowed := testCandidates(record.Path)
		filtered := allowed[:0]
		for _, path := range allowed {
			if !existing[path] && !protected(path, protectedPaths) && testdiscovery.CandidatePath(path) {
				filtered = append(filtered, path)
			}
		}
		allowed = uniqueSorted(filtered)
		if len(allowed) == 0 {
			continue
		}
		// One task may add exactly one file. Multiple alternative paths would
		// let a provider create an unnecessary batch while still staying inside
		// the syntactic allowlist.
		allowed = allowed[:1]
		controls := append([]string(nil), assertion.ControlIDs...)
		sort.Strings(controls)
		draft := provider.Task{
			SchemaVersion: provider.TaskSchema, Mode: "suggest",
			FindingID: finding.ID, FindingFingerprint: finding.Fingerprint,
			AssertionID: assertion.ID, ControlIDs: controls,
			Goal: "Add one focused, non-vacuous, discoverable test for behavior visible in " + record.Path +
				". Add exactly one new test file from allowed_paths. Do not modify production code.",
			ReadScope: "task-inputs", RelevantPaths: []string{record.Path}, Inputs: []provider.InputFile{},
			AllowedPaths: allowed, ProtectedPaths: append([]string(nil), protectedPaths...),
			AllowedCommands: [][]string{}, Network: "deny", Secrets: "none",
			AllowRemoteSourceProcessing: true, TimeoutSeconds: timeout,
			MaxOutputBytes: maxOutput, MaxCostUSD: options.MaxCostUSD,
		}
		task, err := provider.SealTaskValueWithInventory(draft, item.Root, item, protectedPaths)
		if err != nil {
			if errors.Is(err, provider.ErrSensitiveInput) {
				return provider.Task{}, false, policyDenied(err)
			}
			return provider.Task{}, false, err
		}
		return task, true, nil
	}
	return provider.Task{}, false, nil
}

func testCandidates(source string) []string {
	directory := filepath.ToSlash(filepath.Dir(filepath.FromSlash(source)))
	if directory == "." {
		directory = ""
	}
	base := filepath.Base(source)
	extension := strings.ToLower(filepath.Ext(base))
	stem := strings.TrimSuffix(base, filepath.Ext(base))
	join := func(parts ...string) string { return filepath.ToSlash(filepath.Join(parts...)) }
	switch extension {
	case ".go":
		return []string{join(directory, stem+"_test.go")}
	case ".py":
		return []string{join(directory, stem+"_test.py"), join(directory, "test_"+stem+".py"), join("tests", "test_"+stem+".py")}
	case ".js", ".jsx", ".ts", ".tsx":
		return []string{join("tests", stem+".test"+extension), join("tests", stem+".spec"+extension)}
	default:
		return nil
	}
}

func writeSealedAgentTask(directory string, task provider.Task) error {
	data, err := encodedSealedAgentTask(task)
	if err != nil {
		return err
	}
	path := filepath.Join(directory, "agent-task.json")
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		return fmt.Errorf("create sealed agent task record: %w", err)
	}
	if _, err := file.Write(data); err != nil {
		file.Close()
		return fmt.Errorf("write sealed agent task record: %w", err)
	}
	if err := file.Sync(); err != nil {
		file.Close()
		return fmt.Errorf("sync sealed agent task record: %w", err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("close sealed agent task record: %w", err)
	}
	return nil
}

func verifySealedAgentTask(directory string, task provider.Task) error {
	expected, err := encodedSealedAgentTask(task)
	if err != nil {
		return err
	}
	path := filepath.Join(directory, "agent-task.json")
	info, err := os.Lstat(path)
	if err != nil || !info.Mode().IsRegular() || info.Size() != int64(len(expected)) ||
		(runtime.GOOS != "windows" && info.Mode().Perm() != 0o600) {
		return fmt.Errorf("sealed agent task record changed during provider execution")
	}
	actual, err := os.ReadFile(path)
	if err != nil || !bytes.Equal(actual, expected) {
		return fmt.Errorf("sealed agent task record changed during provider execution")
	}
	return nil
}

func encodedSealedAgentTask(task provider.Task) ([]byte, error) {
	data, err := json.MarshalIndent(task, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("encode sealed agent task: %w", err)
	}
	return append(data, '\n'), nil
}
