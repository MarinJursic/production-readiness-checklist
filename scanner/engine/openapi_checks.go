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
	openAPIVersionPattern      = regexp.MustCompile(`^3\.([0-9]+)\.[0-9]+(?:-[0-9A-Za-z.-]+)?(?:\+[0-9A-Za-z.-]+)?$`)
	openAPIResponseCodePattern = regexp.MustCompile(`^(?:default|[1-5](?:[0-9]{2}|XX))$`)
	errUnsupportedOpenAPI      = errors.New("unsupported OpenAPI feature version")
)

var openAPIOperationFields = map[string]bool{
	"get": true, "put": true, "post": true, "delete": true, "options": true,
	"head": true, "patch": true, "trace": true,
}

type openAPIProblem struct {
	Path    string
	Line    int
	Column  int
	Message string
}

type openAPIOperation struct {
	Label string
	Node  *yaml.Node
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

func openAPIFeatureMinor(root *yaml.Node) int {
	root = resolveOpenAPINode(root)
	version := resolveOpenAPINode(mappingValue(root, "openapi"))
	if version == nil || version.Kind != yaml.ScalarNode || version.Tag != "!!str" {
		return -1
	}
	match := openAPIVersionPattern.FindStringSubmatch(strings.TrimSpace(version.Value))
	if match == nil {
		return -1
	}
	switch match[1] {
	case "0":
		return 0
	case "1":
		return 1
	case "2":
		return 2
	default:
		return -1
	}
}

func appendOpenAPIOperations(
	documentPath, sectionName string,
	section *yaml.Node,
	featureMinor int,
	operations *[]openAPIOperation,
	problems *[]openAPIProblem,
) {
	section = resolveOpenAPINode(section)
	if section == nil {
		return
	}
	if section.Kind != yaml.MappingNode {
		appendOpenAPIProblem(problems, documentPath, section, sectionName+" must be an object")
		return
	}
	for item := 0; item+1 < len(section.Content); item += 2 {
		pathKey, pathItem := section.Content[item], resolveOpenAPINode(section.Content[item+1])
		if pathItem == nil || pathItem.Kind != yaml.MappingNode {
			appendOpenAPIProblem(problems, documentPath, section.Content[item+1], sectionName+" path item "+pathKey.Value+" must be an object")
			continue
		}
		for field := 0; field+1 < len(pathItem.Content); field += 2 {
			methodKey, operationNode := pathItem.Content[field], resolveOpenAPINode(pathItem.Content[field+1])
			method := strings.ToLower(methodKey.Value)
			if openAPIOperationFields[method] || (featureMinor == 2 && method == "query") {
				label := sectionName + " " + pathKey.Value + " " + strings.ToUpper(method)
				if methodKey.Value != method {
					appendOpenAPIProblem(problems, documentPath, methodKey, label+" fixed operation field must be lowercase")
				}
				if operationNode == nil || operationNode.Kind != yaml.MappingNode {
					appendOpenAPIProblem(problems, documentPath, pathItem.Content[field+1], label+" operation must be an object")
					continue
				}
				*operations = append(*operations, openAPIOperation{Label: label, Node: operationNode})
				continue
			}
			if featureMinor != 2 || methodKey.Value != "additionalOperations" {
				continue
			}
			if operationNode == nil || operationNode.Kind != yaml.MappingNode {
				appendOpenAPIProblem(problems, documentPath, pathItem.Content[field+1], sectionName+" "+pathKey.Value+" additionalOperations must be an object")
				continue
			}
			for additional := 0; additional+1 < len(operationNode.Content); additional += 2 {
				additionalKey := operationNode.Content[additional]
				additionalOperation := resolveOpenAPINode(operationNode.Content[additional+1])
				label := sectionName + " " + pathKey.Value + " " + additionalKey.Value
				fixedMethod := strings.ToLower(additionalKey.Value)
				if openAPIOperationFields[fixedMethod] || fixedMethod == "query" {
					appendOpenAPIProblem(problems, documentPath, additionalKey, label+" duplicates a fixed operation field")
				}
				if additionalOperation == nil || additionalOperation.Kind != yaml.MappingNode {
					appendOpenAPIProblem(problems, documentPath, operationNode.Content[additional+1], label+" operation must be an object")
					continue
				}
				*operations = append(*operations, openAPIOperation{Label: label, Node: additionalOperation})
			}
		}
	}
}

func directOpenAPIOperations(path string, root *yaml.Node) ([]openAPIOperation, []openAPIProblem) {
	operations := make([]openAPIOperation, 0)
	problems := make([]openAPIProblem, 0)
	root = resolveOpenAPINode(root)
	if root == nil || root.Kind != yaml.MappingNode {
		return operations, problems
	}
	featureMinor := openAPIFeatureMinor(root)
	appendOpenAPIOperations(path, "paths", mappingValue(root, "paths"), featureMinor, &operations, &problems)
	if featureMinor >= 1 {
		appendOpenAPIOperations(path, "webhooks", mappingValue(root, "webhooks"), featureMinor, &operations, &problems)
	}
	return operations, problems
}

func openAPIOperationResponseProblems(path string, operations []openAPIOperation) []openAPIProblem {
	problems := make([]openAPIProblem, 0)
	for _, operation := range operations {
		responsesNode := mappingValue(operation.Node, "responses")
		responses := resolveOpenAPINode(responsesNode)
		if responses == nil || responses.Kind != yaml.MappingNode {
			appendOpenAPIProblem(&problems, path, operation.Node, operation.Label+" responses must be an object")
			continue
		}
		responseCount := 0
		for item := 0; item+1 < len(responses.Content); item += 2 {
			responseKey, responseNode := responses.Content[item], resolveOpenAPINode(responses.Content[item+1])
			if strings.HasPrefix(strings.ToLower(responseKey.Value), "x-") {
				continue
			}
			if responseKey.Kind != yaml.ScalarNode || responseKey.Tag != "!!str" || !openAPIResponseCodePattern.MatchString(responseKey.Value) {
				appendOpenAPIProblem(&problems, path, responseKey, operation.Label+" has invalid response key "+responseKey.Value)
				continue
			}
			responseCount++
			if responseNode == nil || responseNode.Kind != yaml.MappingNode {
				appendOpenAPIProblem(&problems, path, responses.Content[item+1], operation.Label+" response "+responseKey.Value+" must be an object")
				continue
			}
			if reference := mappingValue(responseNode, "$ref"); reference != nil {
				reference = resolveOpenAPINode(reference)
				if reference == nil || reference.Kind != yaml.ScalarNode || reference.Tag != "!!str" || strings.TrimSpace(reference.Value) == "" {
					appendOpenAPIProblem(&problems, path, responseNode, operation.Label+" response "+responseKey.Value+" $ref must be a nonempty string")
				}
				continue
			}
			description := resolveOpenAPINode(mappingValue(responseNode, "description"))
			if description == nil || description.Kind != yaml.ScalarNode || description.Tag != "!!str" || strings.TrimSpace(description.Value) == "" {
				appendOpenAPIProblem(&problems, path, responseNode, operation.Label+" response "+responseKey.Value+" description must be a nonempty string")
			}
		}
		if responseCount == 0 {
			appendOpenAPIProblem(&problems, path, responses, operation.Label+" responses must contain at least one response code or default")
		}
	}
	return problems
}

func openAPIOperationIDProblems(path string, operations []openAPIOperation) ([]openAPIProblem, int) {
	problems := make([]openAPIProblem, 0)
	seen := map[string]string{}
	declared := 0
	for _, operation := range operations {
		identifier := mappingValue(operation.Node, "operationId")
		if identifier == nil {
			continue
		}
		declared++
		identifier = resolveOpenAPINode(identifier)
		if identifier == nil || identifier.Kind != yaml.ScalarNode || identifier.Tag != "!!str" || strings.TrimSpace(identifier.Value) == "" {
			appendOpenAPIProblem(&problems, path, identifier, operation.Label+" operationId must be a nonempty string")
			continue
		}
		if previous, exists := seen[identifier.Value]; exists {
			appendOpenAPIProblem(&problems, path, identifier, operation.Label+" operationId "+identifier.Value+" duplicates "+previous)
			continue
		}
		seen[identifier.Value] = operation.Label
	}
	return problems, declared
}

func applyOpenAPIProblems(
	result model.AssertionResult,
	problems []openAPIProblem,
	successSummary, failureKind string,
) model.AssertionResult {
	if len(problems) == 0 {
		result.Assessment = "pass"
		result.Summary = successSummary
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
	result.Summary = fmt.Sprintf("Detected %d %s problem(s); reporting up to %d: %s.",
		len(problems), failureKind, maximumOpenAPIProblems, strings.Join(messages, ", "))
	return result
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
	return applyOpenAPIProblems(result, problems,
		fmt.Sprintf("All %d detected OpenAPI document(s) have supported 3.0, 3.1, or 3.2 root structure and required metadata.", len(paths)),
		"OpenAPI root-structure")
}

func evaluateOpenAPIOperationResponses(
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
	operationCount := 0
	for _, path := range paths {
		data, evidence, readErr := readVerifiedEvidence(
			inventory, path, "openapi-parse", assertion.ImplementationID,
			"Parsed directly declared OpenAPI operations and response objects.", observedAt,
		)
		if readErr != nil {
			return nativeReadFailure(result, readErr)
		}
		result.EvidenceObserved = append(result.EvidenceObserved, evidence)
		root, parseErr := decodeOpenAPIDocument(data)
		if parseErr != nil {
			return nativeReadFailure(result, fmt.Errorf("cannot parse %s as one OpenAPI YAML or JSON document: %w", path, parseErr))
		}
		if _, validationErr := openAPIRootProblems(path, root); validationErr != nil {
			return nativeReadFailure(result, validationErr)
		}
		operations, traversalProblems := directOpenAPIOperations(path, root)
		operationCount += len(operations)
		problems = append(problems, traversalProblems...)
		problems = append(problems, openAPIOperationResponseProblems(path, operations)...)
	}
	result.EvidenceObserved = append([]model.Evidence{inventoryEvidence(
		inventory, "openapi-parse", assertion.ImplementationID, ".",
		fmt.Sprintf("Bounded OpenAPI operation-response analysis inspected %d directly declared operation(s) in %d document(s).", operationCount, len(paths)), observedAt,
	)}, result.EvidenceObserved...)
	return applyOpenAPIProblems(result, problems,
		fmt.Sprintf("All %d directly declared OpenAPI operation(s) provide a nonempty Responses Object with structurally valid inline responses or references.", operationCount),
		"OpenAPI operation-response")
}

func evaluateOpenAPIOperationIDs(
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
	operationCount, declaredCount := 0, 0
	for _, path := range paths {
		data, evidence, readErr := readVerifiedEvidence(
			inventory, path, "openapi-parse", assertion.ImplementationID,
			"Parsed directly declared OpenAPI operationId values.", observedAt,
		)
		if readErr != nil {
			return nativeReadFailure(result, readErr)
		}
		result.EvidenceObserved = append(result.EvidenceObserved, evidence)
		root, parseErr := decodeOpenAPIDocument(data)
		if parseErr != nil {
			return nativeReadFailure(result, fmt.Errorf("cannot parse %s as one OpenAPI YAML or JSON document: %w", path, parseErr))
		}
		if _, validationErr := openAPIRootProblems(path, root); validationErr != nil {
			return nativeReadFailure(result, validationErr)
		}
		operations, traversalProblems := directOpenAPIOperations(path, root)
		operationCount += len(operations)
		problems = append(problems, traversalProblems...)
		identifierProblems, count := openAPIOperationIDProblems(path, operations)
		declaredCount += count
		problems = append(problems, identifierProblems...)
	}
	result.EvidenceObserved = append([]model.Evidence{inventoryEvidence(
		inventory, "openapi-parse", assertion.ImplementationID, ".",
		fmt.Sprintf("Bounded OpenAPI operation identity analysis inspected %d directly declared operation(s) and %d declared operationId value(s) in %d document(s).", operationCount, declaredCount, len(paths)), observedAt,
	)}, result.EvidenceObserved...)
	return applyOpenAPIProblems(result, problems,
		fmt.Sprintf("All %d declared operationId value(s) across %d directly declared OpenAPI operation(s) are nonempty and unique within each document.", declaredCount, operationCount),
		"OpenAPI operationId")
}
