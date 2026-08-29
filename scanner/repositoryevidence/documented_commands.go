// Package repositoryevidence implements narrow, fail-closed collectors over
// the exact repository bytes captured by an inventory. It never executes code.
package repositoryevidence

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/MarinJursic/production-readiness-checklist/scanner/controlprogram"
	"github.com/MarinJursic/production-readiness-checklist/scanner/controlprogramcatalog"
	"github.com/MarinJursic/production-readiness-checklist/scanner/controlruntime"
	workspaceinventory "github.com/MarinJursic/production-readiness-checklist/scanner/inventory"
	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
	"github.com/MarinJursic/production-readiness-checklist/scanner/providercapability"
)

const (
	DocumentedCommandsCollectorID = "prc.collect.prc.36.004.c1@0.1"
	documentedCommandsControlID   = "PRC-36-004"
	documentedCommandsFactKey     = "prc_36_004_c1.documented_command_text"
	documentedCommandsParameter   = "prc_36_004_c1.required_documented_command_text_keys"
	maximumManifestBytes          = 1 << 20
	maximumDocumentationFileBytes = 1 << 20
	maximumDocumentationBytes     = 8 << 20
	maximumDocumentationFiles     = 512
)

var applicabilityContractSHA256 = digestText("repository-project-applicability/v0.1")

// DocumentedCommandsProvider proves only the supported positive case: a root
// Node package declares build and test scripts and both public invocations are
// present in bounded inventoried Markdown code. Missing or unsupported input
// remains incomplete evidence, never a guessed failure.
type DocumentedCommandsProvider struct {
	inventory model.Inventory
}

func NewDocumentedCommandsProvider(item model.Inventory) (*DocumentedCommandsProvider, error) {
	if item.Root == "" || item.Digest == "" {
		return nil, fmt.Errorf("documented-command collector requires a sealed inventory")
	}
	return &DocumentedCommandsProvider{inventory: item}, nil
}

func (provider *DocumentedCommandsProvider) ID() string { return DocumentedCommandsCollectorID }

func (provider *DocumentedCommandsProvider) Authority() controlprogram.Authority {
	return controlprogram.AuthorityRepository
}

func (provider *DocumentedCommandsProvider) Collect(ctx context.Context, request controlruntime.Request) (controlprogram.Evidence, error) {
	if provider == nil || request.Template.ControlID != documentedCommandsControlID ||
		request.Template.CollectorContract.CollectorID != DocumentedCommandsCollectorID {
		return controlprogram.Evidence{}, fmt.Errorf("documented-command collector received an unsupported template")
	}
	if err := ctx.Err(); err != nil {
		return controlprogram.Evidence{}, err
	}

	// Keep the statement-derived subject keys even when evidence is incomplete.
	// Besides making the missing fields explicit, this preserves the map through
	// canonical JSON (empty maps are omitted by the shared evidence encoding).
	values := map[string]string{"build-command": "", "test-command": ""}
	complete := false
	manifest, manifestErr := workspaceinventory.ReadVerifiedFile(provider.inventory, "package.json", maximumManifestBytes)
	if manifestErr == nil {
		scripts, err := packageScripts(manifest)
		if err != nil {
			return controlprogram.Evidence{}, fmt.Errorf("parse inventoried package.json: %w", err)
		}
		if usableScript(scripts["build"]) && usableTestScript(scripts["test"]) {
			documented, supported, err := provider.findDocumentedCommands(ctx)
			if err != nil {
				return controlprogram.Evidence{}, err
			}
			if supported {
				for key, value := range documented {
					values[key] = value
				}
				complete = values["build-command"] != "" && values["test-command"] != ""
			}
		}
	}

	facts := map[string]controlprogram.Fact{
		documentedCommandsFactKey: {
			Type: controlprogram.FactStringMap, Complete: complete, Values: values,
		},
	}
	return controlruntime.NewApplicableEvidence(
		request,
		"repository-documented-commands-"+controlprogram.ProgramSHA256(request.Program)[:16],
		facts,
		complete,
	)
}

func (provider *DocumentedCommandsProvider) findDocumentedCommands(ctx context.Context) (map[string]string, bool, error) {
	paths := make([]string, 0)
	total := int64(0)
	for _, file := range provider.inventory.Files {
		extension := strings.ToLower(filepath.Ext(file.Path))
		if extension != ".md" && extension != ".markdown" && extension != ".mdx" {
			continue
		}
		if file.Size > maximumDocumentationFileBytes || len(paths) == maximumDocumentationFiles || total+file.Size > maximumDocumentationBytes {
			return nil, false, nil
		}
		paths = append(paths, file.Path)
		total += file.Size
	}
	if len(paths) == 0 {
		return nil, false, nil
	}
	sort.Strings(paths)
	result := map[string]string{}
	for _, relative := range paths {
		if err := ctx.Err(); err != nil {
			return nil, false, err
		}
		data, err := workspaceinventory.ReadVerifiedFile(provider.inventory, relative, maximumDocumentationFileBytes)
		if err != nil {
			return nil, false, fmt.Errorf("read inventoried documentation %s: %w", relative, err)
		}
		for _, command := range markdownCodeCommands(string(data)) {
			if _, exists := result["build-command"]; !exists && matchesInvocation(command, buildInvocations) {
				result["build-command"] = command
			}
			if _, exists := result["test-command"]; !exists && matchesInvocation(command, testInvocations) {
				result["test-command"] = command
			}
		}
	}
	return result, true, nil
}

var buildInvocations = []string{
	"bun run build", "npm run build", "pnpm build", "pnpm run build", "yarn build", "yarn run build",
}

var testInvocations = []string{
	"bun run test", "bun test", "npm run test", "npm test", "pnpm run test", "pnpm test", "yarn run test", "yarn test",
}

// Binding returns the scanner-owned scope and the exact statement-derived
// parameter set for the one collector implemented by this package.
func Binding(item model.Inventory, template controlprogramcatalog.Template) (controlprogramcatalog.RuntimeBinding, bool) {
	if item.Digest == "" || template.ControlID != documentedCommandsControlID || template.ClauseOrdinal != 1 ||
		template.CollectorContract.CollectorID != DocumentedCommandsCollectorID {
		return controlprogramcatalog.RuntimeBinding{}, false
	}
	subject := "repository@sha256:" + item.Digest
	return controlprogramcatalog.RuntimeBinding{
		SubjectID: subject, Subjects: []string{subject}, InventorySHA256: item.Digest,
		AllowNotApplicable: false, ApplicabilityProofContractSHA256: applicabilityContractSHA256,
		MaximumEvidenceAgeSeconds: 300,
		ScannerInventoryParameters: map[string]controlprogram.Parameter{
			documentedCommandsParameter: {Type: controlprogram.FactStringSet, Strings: []string{"build-command", "test-command"}},
		},
	}, true
}

// EvaluateSupported evaluates every exact repository template for which this
// build has a registered collector. It does not manufacture placeholder
// executions for the other templates; fullscan keeps those controls Blocked.
func EvaluateSupported(ctx context.Context, catalog *controlprogramcatalog.Catalog, item model.Inventory, now time.Time) ([]controlruntime.Execution, error) {
	if catalog == nil || now.IsZero() {
		return nil, fmt.Errorf("evaluate repository evidence requires a catalog and evaluation time")
	}
	provider, err := NewDocumentedCommandsProvider(item)
	if err != nil {
		return nil, err
	}
	registry, err := controlruntime.NewRegistry(provider)
	if err != nil {
		return nil, err
	}
	capabilities, err := providercapability.Load()
	if err != nil {
		return nil, err
	}
	declared := make(map[string]providercapability.Capability, len(capabilities))
	for _, capability := range capabilities {
		declared[capability.CollectorID] = capability
		provider, ok := registry.Provider(capability.CollectorID)
		if !ok || provider.Authority() != capability.Authority {
			return nil, fmt.Errorf("shipped provider %s does not match its capability manifest", capability.CollectorID)
		}
	}
	if len(registry.IDs()) != len(capabilities) {
		return nil, fmt.Errorf("runtime provider registry does not match its capability manifest")
	}
	executions := []controlruntime.Execution{}
	matched := map[string]bool{}
	for _, template := range catalog.Templates() {
		capability, supported := declared[template.CollectorContract.CollectorID]
		if !supported {
			continue
		}
		if capability.ControlID != template.ControlID || capability.ClauseOrdinal != template.ClauseOrdinal ||
			capability.Authority != template.RequiredAuthority || !template.RuntimeRequirements.ProviderClaimed ||
			template.CollectorContract.ProviderStatus != "registered" {
			return nil, fmt.Errorf("catalog template %s does not match its shipped provider", template.TemplateID)
		}
		binding, ok := Binding(item, template)
		if !ok {
			return nil, fmt.Errorf("shipped provider %s has no scanner-owned binding", capability.CollectorID)
		}
		executions = append(executions, controlruntime.Evaluate(ctx, template, binding, registry, now))
		matched[capability.CollectorID] = true
	}
	if len(matched) != len(capabilities) {
		return nil, fmt.Errorf("provider capability manifest references a missing catalog template")
	}
	return executions, nil
}

func packageScripts(data []byte) (map[string]string, error) {
	if err := rejectDuplicateJSONKeys(data); err != nil {
		return nil, err
	}
	var document map[string]json.RawMessage
	decoder := json.NewDecoder(bytes.NewReader(data))
	if err := decoder.Decode(&document); err != nil {
		return nil, err
	}
	var trailing any
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		return nil, fmt.Errorf("package.json contains trailing JSON")
	}
	raw, exists := document["scripts"]
	if !exists {
		return map[string]string{}, nil
	}
	var scripts map[string]json.RawMessage
	if err := json.Unmarshal(raw, &scripts); err != nil {
		return nil, fmt.Errorf("scripts is not an object")
	}
	result := make(map[string]string, len(scripts))
	for key, encoded := range scripts {
		var value string
		if err := json.Unmarshal(encoded, &value); err != nil {
			return nil, fmt.Errorf("script %q is not a string", key)
		}
		result[key] = strings.TrimSpace(value)
	}
	return result, nil
}

func rejectDuplicateJSONKeys(data []byte) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	var walk func(int) error
	walk = func(depth int) error {
		if depth > controlprogram.MaxJSONDepth {
			return fmt.Errorf("JSON exceeds nesting limit")
		}
		token, err := decoder.Token()
		if err != nil {
			return err
		}
		delimiter, compound := token.(json.Delim)
		if !compound {
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
				if !ok || seen[key] {
					return fmt.Errorf("duplicate or invalid JSON object key")
				}
				seen[key] = true
				if err := walk(depth + 1); err != nil {
					return err
				}
			}
		case '[':
			for decoder.More() {
				if err := walk(depth + 1); err != nil {
					return err
				}
			}
		default:
			return fmt.Errorf("unexpected JSON delimiter")
		}
		_, err = decoder.Token()
		return err
	}
	if err := walk(1); err != nil {
		return err
	}
	if _, err := decoder.Token(); !errors.Is(err, io.EOF) {
		return fmt.Errorf("package.json contains trailing JSON")
	}
	return nil
}

func usableScript(value string) bool { return strings.TrimSpace(value) != "" }

func usableTestScript(value string) bool {
	lower := strings.ToLower(strings.TrimSpace(value))
	return lower != "" && !strings.Contains(lower, "no test specified") && lower != "exit 0" && lower != "true"
}

func markdownCodeCommands(text string) []string {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	commands := []string{}
	inFence := false
	fence := ""
	for _, line := range strings.Split(text, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "```") || strings.HasPrefix(trimmed, "~~~") {
			marker := trimmed[:3]
			if !inFence {
				inFence, fence = true, marker
			} else if marker == fence {
				inFence, fence = false, ""
			}
			continue
		}
		if inFence || strings.HasPrefix(line, "    ") || strings.HasPrefix(line, "\t") {
			if command := normalizeCommand(trimmed); command != "" {
				commands = append(commands, command)
			}
			continue
		}
		for remainder := line; ; {
			start := strings.IndexByte(remainder, '`')
			if start < 0 {
				break
			}
			remainder = remainder[start+1:]
			end := strings.IndexByte(remainder, '`')
			if end < 0 {
				break
			}
			if command := normalizeCommand(remainder[:end]); command != "" {
				commands = append(commands, command)
			}
			remainder = remainder[end+1:]
		}
	}
	return commands
}

func normalizeCommand(value string) string {
	value = strings.TrimSpace(value)
	if strings.HasPrefix(value, "$ ") || strings.HasPrefix(value, "> ") {
		value = strings.TrimSpace(value[2:])
	}
	return strings.Join(strings.Fields(value), " ")
}

func matchesInvocation(command string, invocations []string) bool {
	for _, invocation := range invocations {
		if command == invocation || strings.HasPrefix(command, invocation+" ") {
			return true
		}
	}
	return false
}

func digestText(value string) string {
	digest := sha256.Sum256([]byte(value))
	return hex.EncodeToString(digest[:])
}
