package engine

import (
	"fmt"
	"sort"
	"strings"

	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
)

const (
	ExecutionModeInspect     = "inspect"
	ExecutionModeVerifyLocal = "verify-local"
)

type implementationSpec struct {
	kind         string
	capabilities model.CapabilitySet
}

var implementationRegistry = func() map[string]implementationSpec {
	readOnly := capabilitySet(true, false, "deny")
	none := capabilitySet(false, false, "deny")
	registry := map[string]implementationSpec{}
	for _, identifier := range []string{
		"prc.native.file-present@0.1",
		"prc.native.dependency-lock@0.1",
		"prc.native.github-action-pin@0.1",
		"prc.native.ci-present@0.1",
		"prc.native.test-suite@0.1",
		"prc.native.github-workflow-permissions@0.1",
		"prc.native.final-newline@0.1",
		"prc.native.git-revision@0.1",
		"prc.native.github-workflow-valid@0.1",
		"prc.native.github-workflow-jobs@0.1",
		"prc.native.github-workflow-timeouts@0.1",
		"prc.native.github-no-pull-request-target@0.1",
		"prc.native.merge-conflict-markers@0.1",
		"prc.native.restrictive-file-modes@0.1",
		"prc.native.inventory-files-nonempty@0.1",
		"prc.native.runtime-version@0.1",
		"prc.native.container-base-pin@0.1",
		"prc.native.container-nonroot@0.1",
		"prc.native.terraform-lock@0.1",
		"prc.native.kubernetes-nonroot@0.2",
		"prc.native.kubernetes-resources@0.2",
		"prc.native.kubernetes-host-access@0.1",
		"prc.native.kubernetes-privilege-escalation@0.1",
		"prc.native.kubernetes-capabilities@0.1",
		"prc.native.kubernetes-seccomp@0.1",
		"prc.native.private-key-armor@0.1",
		"prc.native.go-http-timeout@0.1",
		"prc.native.go-http-server-timeout@0.1",
		"prc.native.openapi-root@0.1",
		"prc.native.openapi-operation-responses@0.1",
		"prc.native.openapi-operation-ids@0.1",
	} {
		registry[identifier] = implementationSpec{kind: "built-in", capabilities: readOnly}
	}
	registry["prc.native.manual-evidence@0.1"] = implementationSpec{kind: "manual", capabilities: none}
	registry["prc.native.analysis-evidence@0.1"] = implementationSpec{kind: "adapter-evidence", capabilities: none}
	return registry
}()

func capabilitySet(readWorkspace, writeScratch bool, process string) model.CapabilitySet {
	return model.CapabilitySet{
		ReadWorkspace: readWorkspace, WriteScratch: writeScratch, Process: process,
		Network: "deny", NetworkHosts: []string{}, SecretHandles: []string{},
	}
}

func capabilityPolicy(mode string) (model.CapabilitySet, error) {
	switch mode {
	case ExecutionModeInspect:
		return capabilitySet(true, false, "deny"), nil
	case ExecutionModeVerifyLocal:
		return capabilitySet(true, true, "oci"), nil
	default:
		return model.CapabilitySet{}, fmt.Errorf("unsupported execution mode %q", mode)
	}
}

func capabilityAllowed(policy, required model.CapabilitySet) (bool, string) {
	if required.ReadWorkspace && !policy.ReadWorkspace {
		return false, "read_workspace is denied by the selected execution mode"
	}
	if required.WriteScratch && !policy.WriteScratch {
		return false, "write_scratch is denied by the selected execution mode"
	}
	if required.Process != "deny" && policy.Process != required.Process {
		return false, "process capability " + required.Process + " is denied by the selected execution mode"
	}
	if required.Network != "deny" || len(required.NetworkHosts) != 0 {
		return false, "network capability is outside the current fail-closed policy"
	}
	if len(required.SecretHandles) != 0 {
		return false, "secret handles are outside the current fail-closed policy"
	}
	return true, ""
}

type plannedAdapterBuilder struct {
	adapter model.PlannedAdapter
	nodeID  string
}

func buildExecutionGraph(plan *model.Plan, assertions map[string]model.Assertion) error {
	policy, err := capabilityPolicy(plan.ExecutionMode)
	if err != nil {
		return err
	}
	plan.CapabilityPolicy = policy
	plan.Implementations = []model.PlannedImplementation{}
	plan.Adapters = []model.PlannedAdapter{}
	plan.Nodes = []model.PlanNode{}

	implementationAssertions := map[string][]string{}
	adapterBuilders := map[string]*plannedAdapterBuilder{}
	assertionBindings := map[string][]adapterBinding{}
	for _, planned := range plan.Assertions {
		implementationAssertions[planned.Implementation] = append(
			implementationAssertions[planned.Implementation], planned.AssertionID,
		)
		if planned.Applicability != "applicable" || planned.Implementation != "prc.native.analysis-evidence@0.1" {
			continue
		}
		bindings, bindingErr := assertionAdapterBindings(assertions[planned.AssertionID])
		if bindingErr != nil {
			return bindingErr
		}
		assertionBindings[planned.AssertionID] = bindings
		for _, binding := range bindings {
			key := binding.AdapterID + "\x00" + binding.ManifestSHA256
			builder, exists := adapterBuilders[key]
			if !exists {
				required := capabilitySet(true, true, "oci")
				available, reason := capabilityAllowed(policy, required)
				status := "authorized"
				if !available {
					status = "blocked"
				}
				builder = &plannedAdapterBuilder{
					adapter: model.PlannedAdapter{
						AdapterID: binding.AdapterID, ManifestSHA256: binding.ManifestSHA256,
						ObservationKinds: []string{}, Capabilities: required, Status: status, Reason: reason,
					},
					nodeID: "adapter:" + binding.ManifestSHA256,
				}
				adapterBuilders[key] = builder
			}
			builder.adapter.ObservationKinds = append(builder.adapter.ObservationKinds, binding.ObservationKind)
		}
	}

	implementationIDs := make([]string, 0, len(implementationAssertions))
	for identifier := range implementationAssertions {
		implementationIDs = append(implementationIDs, identifier)
	}
	sort.Strings(implementationIDs)
	for _, identifier := range implementationIDs {
		assertionIDs := uniqueSorted(implementationAssertions[identifier])
		spec, exists := implementationRegistry[identifier]
		implementation := model.PlannedImplementation{
			ID: identifier, Kind: "unavailable", AssertionIDs: assertionIDs,
			Capabilities: capabilitySet(false, false, "deny"), Status: "blocked",
			Reason: "No scanner implementation is registered for this exact implementation ID.",
		}
		if exists {
			implementation.Kind = spec.kind
			implementation.Capabilities = spec.capabilities
			implementation.Status = "available"
			implementation.Reason = ""
		}
		plan.Implementations = append(plan.Implementations, implementation)
	}

	adapterKeys := make([]string, 0, len(adapterBuilders))
	for key := range adapterBuilders {
		adapterKeys = append(adapterKeys, key)
	}
	sort.Strings(adapterKeys)
	for _, key := range adapterKeys {
		builder := adapterBuilders[key]
		builder.adapter.ObservationKinds = uniqueSorted(builder.adapter.ObservationKinds)
		plan.Adapters = append(plan.Adapters, builder.adapter)
	}

	plan.Nodes = append(plan.Nodes, model.PlanNode{
		ID: "inventory:" + plan.InventoryDigest, Kind: "inventory", DependsOn: []string{},
		Capabilities: capabilitySet(true, false, "deny"), Status: "ready",
	})
	for _, key := range adapterKeys {
		builder := adapterBuilders[key]
		status := builder.adapter.Status
		if status == "authorized" {
			status = "ready"
		}
		plan.Nodes = append(plan.Nodes, model.PlanNode{
			ID: builder.nodeID, Kind: "adapter", DependsOn: []string{"inventory:" + plan.InventoryDigest},
			AdapterID: builder.adapter.AdapterID, ManifestSHA256: builder.adapter.ManifestSHA256,
			Capabilities: builder.adapter.Capabilities, Status: status, Reason: builder.adapter.Reason,
		})
	}

	assertionNodeIDs := make([]string, 0, len(plan.Assertions))
	implementationsByID := map[string]model.PlannedImplementation{}
	for _, implementation := range plan.Implementations {
		implementationsByID[implementation.ID] = implementation
	}
	for _, planned := range plan.Assertions {
		assertion := assertions[planned.AssertionID]
		implementation := implementationsByID[planned.Implementation]
		node := model.PlanNode{
			ID: "assertion:" + planned.AssertionID, Kind: "assertion",
			DependsOn: []string{"inventory:" + plan.InventoryDigest}, AssertionID: planned.AssertionID,
			ImplementationID: planned.Implementation, Capabilities: implementation.Capabilities, Status: "ready",
		}
		switch {
		case planned.Applicability == "not_applicable":
			node.Status, node.Reason = "skipped", "The applicability expression evaluated to false."
		case planned.Applicability == "undetermined":
			node.Status, node.Reason = "blocked", planned.ApplicabilityReason
		case implementation.Status == "blocked":
			node.Status, node.Reason = "blocked", implementation.Reason
		case assertion.ImplementationID == "prc.native.analysis-evidence@0.1":
			bindings := assertionBindings[planned.AssertionID]
			if len(bindings) == 0 {
				node.Status, node.Reason = "blocked", "No immutable adapter bindings are configured for this assertion."
				break
			}
			node.DependsOn = []string{"inventory:" + plan.InventoryDigest}
			for _, binding := range bindings {
				builder := adapterBuilders[binding.AdapterID+"\x00"+binding.ManifestSHA256]
				node.DependsOn = append(node.DependsOn, builder.nodeID)
				if builder.adapter.Status == "blocked" {
					node.Status, node.Reason = "blocked", builder.adapter.Reason
				}
			}
			node.DependsOn = uniqueSorted(node.DependsOn)
		}
		plan.Nodes = append(plan.Nodes, node)
		assertionNodeIDs = append(assertionNodeIDs, node.ID)
	}
	gate := model.PlanNode{
		ID: "gate:" + strings.TrimPrefix(plan.ProfileID, "prc/"), Kind: "gate",
		DependsOn: assertionNodeIDs, Capabilities: capabilitySet(false, false, "deny"), Status: "ready",
	}
	for _, node := range plan.Nodes {
		if node.Kind == "assertion" && node.Status == "blocked" {
			gate.Status, gate.Reason = "blocked", "One or more planned assertions cannot execute under the resolved plan."
			break
		}
	}
	plan.Nodes = append(plan.Nodes, gate)
	return ValidateExecutionPlan(*plan)
}

// ValidateExecutionPlan enforces the semantic DAG and capability contract that
// JSON Schema alone cannot express. It is also used when durable run records
// are reopened, so rehashed but internally inconsistent plans fail closed.
func ValidateExecutionPlan(plan model.Plan) error {
	if len(plan.Nodes) != len(plan.Assertions)+len(plan.Adapters)+2 {
		return fmt.Errorf("execution plan must contain one inventory node, one node per adapter and assertion, and one gate node")
	}
	policy, err := capabilityPolicy(plan.ExecutionMode)
	if err != nil {
		return err
	}
	if !sameCapabilities(policy, plan.CapabilityPolicy) {
		return fmt.Errorf("execution plan capability policy does not match mode %s", plan.ExecutionMode)
	}
	implementations, err := validateImplementationRegistry(plan)
	if err != nil {
		return err
	}

	seen := map[string]bool{}
	seenNodes := map[string]model.PlanNode{}
	assertionNodes := map[string]bool{}
	assertions := map[string]model.PlannedAssertion{}
	for _, assertion := range plan.Assertions {
		assertions[assertion.AssertionID] = assertion
	}
	adapters := map[string]model.PlannedAdapter{}
	for _, adapter := range plan.Adapters {
		key := adapter.AdapterID + "\x00" + adapter.ManifestSHA256
		if _, exists := adapters[key]; exists {
			return fmt.Errorf("execution plan contains a duplicate adapter binding")
		}
		required := capabilitySet(true, true, "oci")
		allowed, reason := capabilityAllowed(policy, required)
		expectedStatus := "authorized"
		if !allowed {
			expectedStatus = "blocked"
		}
		if strings.TrimSpace(adapter.AdapterID) == "" || strings.TrimSpace(adapter.ManifestSHA256) == "" ||
			len(adapter.ObservationKinds) == 0 || !uniqueNonempty(adapter.ObservationKinds) ||
			!sameCapabilities(adapter.Capabilities, required) || adapter.Status != expectedStatus ||
			(adapter.Status == "blocked" && adapter.Reason != reason) ||
			(adapter.Status == "authorized" && adapter.Reason != "") {
			return fmt.Errorf("execution plan adapter %s does not match mode policy", adapter.AdapterID)
		}
		adapters[key] = adapter
	}
	adapterNodes := map[string]bool{}
	adapterUses := map[string]bool{}
	inventoryNodes, gateNodes := 0, 0
	assertionNodeIDs := make([]string, 0, len(plan.Assertions))
	for index, node := range plan.Nodes {
		if strings.TrimSpace(node.ID) == "" || seen[node.ID] {
			return fmt.Errorf("execution plan contains an empty or duplicate node ID %q", node.ID)
		}
		dependencySeen := map[string]bool{}
		for _, dependency := range node.DependsOn {
			if dependencySeen[dependency] {
				return fmt.Errorf("execution plan node %s contains duplicate dependency %s", node.ID, dependency)
			}
			if !seen[dependency] {
				return fmt.Errorf("execution plan node %s depends on unavailable or later node %s", node.ID, dependency)
			}
			dependencySeen[dependency] = true
		}
		seen[node.ID] = true
		seenNodes[node.ID] = node
		if node.Status != "ready" && node.Status != "skipped" && node.Status != "blocked" {
			return fmt.Errorf("execution plan node %s has unsupported status %q", node.ID, node.Status)
		}
		if node.Status == "ready" {
			if allowed, reason := capabilityAllowed(policy, node.Capabilities); !allowed {
				return fmt.Errorf("execution plan node %s exceeds its capability policy: %s", node.ID, reason)
			}
		}

		switch node.Kind {
		case "inventory":
			inventoryNodes++
			if index != 0 || node.ID != "inventory:"+plan.InventoryDigest || len(node.DependsOn) != 0 ||
				node.Status != "ready" || node.Reason != "" ||
				node.AssertionID != "" || node.ImplementationID != "" || node.AdapterID != "" || node.ManifestSHA256 != "" ||
				!sameCapabilities(node.Capabilities, capabilitySet(true, false, "deny")) {
				return fmt.Errorf("execution plan inventory node must be the first ready root node")
			}
		case "adapter":
			key := node.AdapterID + "\x00" + node.ManifestSHA256
			adapter, exists := adapters[key]
			if !exists || adapterNodes[key] {
				return fmt.Errorf("execution plan contains an unbound or duplicate adapter node %s", node.ID)
			}
			expectedStatus := adapter.Status
			if expectedStatus == "authorized" {
				expectedStatus = "ready"
			}
			if node.ID != "adapter:"+node.ManifestSHA256 ||
				len(node.DependsOn) != 1 || node.DependsOn[0] != "inventory:"+plan.InventoryDigest ||
				node.AssertionID != "" || node.ImplementationID != "" ||
				node.Status != expectedStatus || node.Reason != adapter.Reason ||
				!sameCapabilities(node.Capabilities, adapter.Capabilities) {
				return fmt.Errorf("execution plan adapter node %s does not match its immutable binding", node.ID)
			}
			adapterNodes[key] = true
		case "assertion":
			assertion, exists := assertions[node.AssertionID]
			if !exists || assertionNodes[node.AssertionID] {
				return fmt.Errorf("execution plan contains an empty or duplicate assertion node")
			}
			implementation := implementations[node.ImplementationID]
			if node.ID != "assertion:"+node.AssertionID || node.ImplementationID != assertion.Implementation ||
				node.AdapterID != "" || node.ManifestSHA256 != "" ||
				!sameCapabilities(node.Capabilities, implementation.Capabilities) ||
				!dependencySeen["inventory:"+plan.InventoryDigest] {
				return fmt.Errorf("execution plan assertion node %s does not match its planned assertion", node.ID)
			}
			for dependency := range dependencySeen {
				if dependency == "inventory:"+plan.InventoryDigest {
					continue
				}
				if seenNodes[dependency].Kind != "adapter" || implementation.Kind != "adapter-evidence" || assertion.Applicability != "applicable" {
					return fmt.Errorf("execution plan assertion node %s has a non-adapter or unnecessary dependency", node.ID)
				}
				adapterUses[dependency] = true
			}
			expectedStatus, expectedReason := "ready", ""
			switch assertion.Applicability {
			case "not_applicable":
				expectedStatus = "skipped"
				expectedReason = "The applicability expression evaluated to false."
			case "undetermined":
				expectedStatus = "blocked"
				expectedReason = assertion.ApplicabilityReason
			case "applicable":
				if implementation.Status == "blocked" {
					expectedStatus = "blocked"
					expectedReason = implementation.Reason
				} else if implementation.Kind == "adapter-evidence" && len(node.DependsOn) == 1 {
					expectedStatus = "blocked"
					expectedReason = "No immutable adapter bindings are configured for this assertion."
				}
				for dependency := range dependencySeen {
					if dependencyNode, found := findDependencyNode(plan.Nodes[:index], dependency); found && dependencyNode.Status == "blocked" {
						expectedStatus = "blocked"
						expectedReason = dependencyNode.Reason
					}
				}
			default:
				return fmt.Errorf("execution plan assertion %s has unsupported applicability %q", assertion.AssertionID, assertion.Applicability)
			}
			if node.Status != expectedStatus || node.Reason != expectedReason {
				return fmt.Errorf("execution plan assertion node %s has an inconsistent execution state", node.ID)
			}
			assertionNodes[node.AssertionID] = true
			assertionNodeIDs = append(assertionNodeIDs, node.ID)
		case "gate":
			gateNodes++
			if index != len(plan.Nodes)-1 || node.ID != "gate:"+strings.TrimPrefix(plan.ProfileID, "prc/") ||
				node.AssertionID != "" || node.ImplementationID != "" || node.AdapterID != "" || node.ManifestSHA256 != "" ||
				!sameStringSet(node.DependsOn, assertionNodeIDs) {
				return fmt.Errorf("execution plan gate node must be the final node and depend on every assertion")
			}
			expectedStatus := "ready"
			for _, prior := range plan.Nodes[:index] {
				if prior.Kind == "assertion" && prior.Status == "blocked" {
					expectedStatus = "blocked"
					break
				}
			}
			if node.Status != expectedStatus || (node.Status == "blocked" && node.Reason == "") {
				return fmt.Errorf("execution plan gate node has an inconsistent execution state")
			}
			expectedReason := ""
			if expectedStatus == "blocked" {
				expectedReason = "One or more planned assertions cannot execute under the resolved plan."
			}
			if node.Reason != expectedReason || !sameCapabilities(node.Capabilities, capabilitySet(false, false, "deny")) {
				return fmt.Errorf("execution plan gate node does not match its capability or reason contract")
			}
		default:
			return fmt.Errorf("execution plan node %s has unsupported kind %q", node.ID, node.Kind)
		}
	}
	if inventoryNodes != 1 || gateNodes != 1 || len(assertionNodes) != len(plan.Assertions) || len(adapterNodes) != len(plan.Adapters) {
		return fmt.Errorf("execution plan topology is incomplete")
	}
	for _, node := range plan.Nodes {
		if node.Kind == "adapter" && !adapterUses[node.ID] {
			return fmt.Errorf("execution plan adapter node %s is not required by an assertion", node.ID)
		}
	}
	return nil
}

func validateImplementationRegistry(plan model.Plan) (map[string]model.PlannedImplementation, error) {
	uses := map[string][]string{}
	seenAssertions := map[string]bool{}
	for _, assertion := range plan.Assertions {
		if assertion.AssertionID == "" || seenAssertions[assertion.AssertionID] {
			return nil, fmt.Errorf("execution plan contains an empty or duplicate planned assertion")
		}
		seenAssertions[assertion.AssertionID] = true
		uses[assertion.Implementation] = append(uses[assertion.Implementation], assertion.AssertionID)
	}
	implementations := map[string]model.PlannedImplementation{}
	for _, implementation := range plan.Implementations {
		if implementation.ID == "" || implementations[implementation.ID].ID != "" {
			return nil, fmt.Errorf("execution plan contains an empty or duplicate implementation registry entry")
		}
		expectedUses, expected := uses[implementation.ID]
		if !expected || !sameStringSet(implementation.AssertionIDs, expectedUses) {
			return nil, fmt.Errorf("execution plan implementation %s does not map exactly to its assertions", implementation.ID)
		}
		spec, available := implementationRegistry[implementation.ID]
		if available {
			if implementation.Kind != spec.kind || implementation.Status != "available" ||
				!sameCapabilities(implementation.Capabilities, spec.capabilities) || implementation.Reason != "" {
				return nil, fmt.Errorf("execution plan implementation %s does not match the built-in registry", implementation.ID)
			}
		} else if implementation.Kind != "unavailable" || implementation.Status != "blocked" ||
			!sameCapabilities(implementation.Capabilities, capabilitySet(false, false, "deny")) || implementation.Reason == "" {
			return nil, fmt.Errorf("execution plan implementation %s must fail closed as unavailable", implementation.ID)
		}
		implementations[implementation.ID] = implementation
	}
	if len(implementations) != len(uses) {
		return nil, fmt.Errorf("execution plan implementation registry is incomplete")
	}
	return implementations, nil
}

func findDependencyNode(nodes []model.PlanNode, identifier string) (model.PlanNode, bool) {
	for _, node := range nodes {
		if node.ID == identifier {
			return node, true
		}
	}
	return model.PlanNode{}, false
}

func uniqueNonempty(values []string) bool {
	seen := map[string]bool{}
	for _, value := range values {
		if strings.TrimSpace(value) == "" || seen[value] {
			return false
		}
		seen[value] = true
	}
	return true
}

func sameCapabilities(left, right model.CapabilitySet) bool {
	return left.ReadWorkspace == right.ReadWorkspace && left.WriteScratch == right.WriteScratch &&
		left.Process == right.Process && left.Network == right.Network &&
		(left.NetworkHosts == nil) == (right.NetworkHosts == nil) &&
		(left.SecretHandles == nil) == (right.SecretHandles == nil) &&
		sameStringSet(left.NetworkHosts, right.NetworkHosts) && sameStringSet(left.SecretHandles, right.SecretHandles)
}

func sameStringSet(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	leftSet := map[string]bool{}
	for _, value := range left {
		if leftSet[value] {
			return false
		}
		leftSet[value] = true
	}
	rightSet := map[string]bool{}
	for _, value := range right {
		if rightSet[value] || !leftSet[value] {
			return false
		}
		rightSet[value] = true
	}
	return true
}

func uniqueSorted(values []string) []string {
	seen := map[string]bool{}
	for _, value := range values {
		seen[value] = true
	}
	result := make([]string, 0, len(seen))
	for value := range seen {
		result = append(result, value)
	}
	sort.Strings(result)
	return result
}
