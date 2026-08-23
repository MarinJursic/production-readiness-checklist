package engine

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
	"gopkg.in/yaml.v3"
)

const (
	maximumOpenAPIFiles    = 256
	maximumOpenAPIBytes    = 64 * 1024 * 1024
	maximumOpenAPIProblems = 100
)

var (
	openAPIVersionPattern = regexp.MustCompile(`^3\.([0-9]+)\.[0-9]+(?:-[0-9A-Za-z.-]+)?(?:\+[0-9A-Za-z.-]+)?$`)
	errUnsupportedOpenAPI = errors.New("unsupported OpenAPI feature version")
)

type openAPIProblem struct {
	Path    string
	Line    int
	Column  int
	Message string
}

func openAPIDescriptionPaths(inventory model.Inventory) ([]string, error) {
	paths := make([]string, 0)
	seen := map[string]bool{}
	var totalBytes int64
	for _, component := range inventory.Components {
		if component.Kind != "api-description" || component.Ecosystem != "openapi" {
			continue
		}
		if seen[component.Path] {
			continue
		}
		seen[component.Path] = true
		record, ok := expectedFile(inventory, component.Path)
		if !ok {
			return nil, fmt.Errorf("OpenAPI inventory component %s has no file record", component.Path)
		}
		paths = append(paths, component.Path)
		if len(paths) > maximumOpenAPIFiles {
			return nil, fmt.Errorf("%w: OpenAPI analysis found more than %d documents", errNativeInputLimit, maximumOpenAPIFiles)
		}
		if record.Size < 0 || record.Size > maximumOpenAPIBytes-totalBytes {
			return nil, fmt.Errorf("%w: OpenAPI analysis exceeds %d total bytes", errNativeInputLimit, maximumOpenAPIBytes)
		}
		totalBytes += record.Size
	}
	sort.Strings(paths)
	return paths, nil
}

func decodeOpenAPIDocument(data []byte) (*yaml.Node, error) {
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	var document yaml.Node
	if err := decoder.Decode(&document); err != nil {
		return nil, err
	}
	if len(document.Content) != 1 {
		return nil, fmt.Errorf("document must contain one root value")
	}
	var trailing yaml.Node
	if err := decoder.Decode(&trailing); err != io.EOF {
		if err == nil {
			return nil, fmt.Errorf("document contains more than one YAML document")
		}
		return nil, err
	}
	return document.Content[0], nil
}

func appendOpenAPIProblem(problems *[]openAPIProblem, path string, node *yaml.Node, message string) {
	line, column := 1, 1
	if node != nil {
		if node.Line > 0 {
			line = node.Line
		}
		if node.Column > 0 {
			column = node.Column
		}
	}
	*problems = append(*problems, openAPIProblem{Path: path, Line: line, Column: column, Message: message})
}

func validateOpenAPIYAMLShape(path string, node *yaml.Node, problems *[]openAPIProblem) {
	type pendingNode struct {
		node  *yaml.Node
		depth int
	}
	stack := []pendingNode{{node: node}}
	visited := 0
	for len(stack) > 0 {
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		visited++
		if visited > 100_000 || current.depth > 128 {
			appendOpenAPIProblem(problems, path, current.node, "YAML structure exceeds the bounded node or nesting limit")
			return
		}
		if current.node.Kind == yaml.AliasNode {
			continue
		}
		if current.node.Kind == yaml.MappingNode {
			seen := map[string]bool{}
			for index := 0; index+1 < len(current.node.Content); index += 2 {
				key, value := current.node.Content[index], current.node.Content[index+1]
				if key.Kind != yaml.ScalarNode || key.Tag != "!!str" {
					appendOpenAPIProblem(problems, path, key, "mapping keys must be scalar strings")
				} else if seen[key.Value] {
					appendOpenAPIProblem(problems, path, key, "duplicate mapping key "+key.Value)
				} else {
					seen[key.Value] = true
				}
				stack = append(stack, pendingNode{node: value, depth: current.depth + 1})
			}
			continue
		}
		for _, child := range current.node.Content {
			stack = append(stack, pendingNode{node: child, depth: current.depth + 1})
		}
	}
}

func resolveOpenAPINode(node *yaml.Node) *yaml.Node {
	seen := map[*yaml.Node]bool{}
	for steps := 0; node != nil && node.Kind == yaml.AliasNode && steps < 32; steps++ {
		if node.Alias == nil || seen[node] {
			return nil
		}
		seen[node] = true
		node = node.Alias
	}
	if node != nil && node.Kind == yaml.AliasNode {
		return nil
	}
	return node
}

func requireOpenAPIString(path, field string, node *yaml.Node, problems *[]openAPIProblem) {
	node = resolveOpenAPINode(node)
	if node == nil || node.Kind != yaml.ScalarNode || node.Tag != "!!str" || strings.TrimSpace(node.Value) == "" {
		appendOpenAPIProblem(problems, path, node, field+" must be a nonempty string")
	}
}

func openAPIRootProblems(path string, root *yaml.Node) ([]openAPIProblem, error) {
	problems := make([]openAPIProblem, 0)
	root = resolveOpenAPINode(root)
	if root == nil {
		appendOpenAPIProblem(&problems, path, root, "OpenAPI document root alias cannot be resolved safely")
		return problems, nil
	}
	if root.Kind != yaml.MappingNode {
		appendOpenAPIProblem(&problems, path, root, "OpenAPI document root must be an object")
		return problems, nil
	}
	validateOpenAPIYAMLShape(path, root, &problems)
	versionNode := mappingValue(root, "openapi")
	requireOpenAPIString(path, "openapi", versionNode, &problems)
	versionNode = resolveOpenAPINode(versionNode)
	minor := -1
	if versionNode != nil && versionNode.Kind == yaml.ScalarNode && versionNode.Tag == "!!str" {
		match := openAPIVersionPattern.FindStringSubmatch(strings.TrimSpace(versionNode.Value))
		if match == nil {
			appendOpenAPIProblem(&problems, path, versionNode, "openapi must be a semantic 3.x specification version")
		} else {
			switch match[1] {
			case "0":
				minor = 0
			case "1":
				minor = 1
			case "2":
				minor = 2
			default:
				return problems, fmt.Errorf("%w %q in %s", errUnsupportedOpenAPI, versionNode.Value, path)
			}
		}
	}
	info := mappingValue(root, "info")
	info = resolveOpenAPINode(info)
	if info == nil || info.Kind != yaml.MappingNode {
		appendOpenAPIProblem(&problems, path, info, "info must be an object")
	} else {
		requireOpenAPIString(path, "info.title", mappingValue(info, "title"), &problems)
		requireOpenAPIString(path, "info.version", mappingValue(info, "version"), &problems)
	}
	paths := resolveOpenAPINode(mappingValue(root, "paths"))
	components := resolveOpenAPINode(mappingValue(root, "components"))
	webhooks := resolveOpenAPINode(mappingValue(root, "webhooks"))
	if minor == 0 {
		if paths == nil || paths.Kind != yaml.MappingNode {
			appendOpenAPIProblem(&problems, path, paths, "OpenAPI 3.0 paths must be an object")
		}
	} else if minor == 1 || minor == 2 {
		if paths == nil && components == nil && webhooks == nil {
			appendOpenAPIProblem(&problems, path, root, "OpenAPI 3.1 and 3.2 require at least one of paths, components, or webhooks")
		}
		for _, field := range []struct {
			name  string
			value *yaml.Node
		}{{"paths", paths}, {"components", components}, {"webhooks", webhooks}} {
			if field.value != nil && field.value.Kind != yaml.MappingNode {
				appendOpenAPIProblem(&problems, path, field.value, field.name+" must be an object")
			}
		}
	}
	return problems, nil
}

func evaluateOpenAPIRootStructure(
	assertion model.Assertion,
	inventory model.Inventory,
	result model.AssertionResult,
	observedAt time.Time,
) model.AssertionResult {
	paths, err := openAPIDescriptionPaths(inventory)
	if err != nil {
		return nativeReadFailure(result, err)
	}
	problems := make([]openAPIProblem, 0)
	for _, path := range paths {
		data, evidence, readErr := readVerifiedEvidence(
			inventory, path, "openapi-parse", assertion.ImplementationID,
			"Parsed the inventoried OpenAPI document's bounded root structure.", observedAt,
		)
		if readErr != nil {
			return nativeReadFailure(result, readErr)
		}
		result.EvidenceObserved = append(result.EvidenceObserved, evidence)
		root, parseErr := decodeOpenAPIDocument(data)
		if parseErr != nil {
			return nativeReadFailure(result, fmt.Errorf("cannot parse %s as one OpenAPI YAML or JSON document: %w", path, parseErr))
		}
		documentProblems, validationErr := openAPIRootProblems(path, root)
		if validationErr != nil {
			return nativeReadFailure(result, validationErr)
		}
		problems = append(problems, documentProblems...)
	}
	result.EvidenceObserved = append([]model.Evidence{inventoryEvidence(
		inventory, "openapi-parse", assertion.ImplementationID, ".",
		fmt.Sprintf("Bounded OpenAPI root analysis inspected %d inventoried document(s).", len(paths)), observedAt,
	)}, result.EvidenceObserved...)
	if len(problems) == 0 {
		result.Assessment = "pass"
		result.Summary = fmt.Sprintf("All %d detected OpenAPI document(s) have supported 3.0, 3.1, or 3.2 root structure and required metadata.", len(paths))
		return result
	}
	reported := problems
	if len(reported) > maximumOpenAPIProblems {
		reported = reported[:maximumOpenAPIProblems]
	}
	messages := make([]string, 0, len(reported))
	locationSeen := map[string]bool{}
	for _, problem := range reported {
		messages = append(messages, fmt.Sprintf("%s:%d:%d (%s)", problem.Path, problem.Line, problem.Column, problem.Message))
		locationKey := fmt.Sprintf("%s:%d:%d", problem.Path, problem.Line, problem.Column)
		if !locationSeen[locationKey] {
			locationSeen[locationKey] = true
			result.Locations = append(result.Locations, model.FindingLocation{
				Path: problem.Path, Line: problem.Line, Column: problem.Column,
			})
		}
	}
	result.Assessment = "fail"
	result.Summary = fmt.Sprintf("Detected %d OpenAPI root-structure problem(s); reporting up to %d: %s.",
		len(problems), maximumOpenAPIProblems, strings.Join(messages, ", "))
	return result
}
