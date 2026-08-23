package adapter

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"regexp"
	"slices"
	"sort"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	SyftProtocolVersion     = "prc-adapter-syft-cyclonedx-json-v1"
	SyftOutputSchemaVersion = "cyclonedx.bom/v1.7+syft-1.51.0"
	SyftToolVersion         = "1.51.0"
	SyftObservationKind     = "sbom-generation"
	SyftImage               = "ghcr.io/anchore/syft@sha256:d2dc3ec86cb2b4e7ddb226ba0305c4523b7c0694c45d9f576b42b4c2f5ce7aa8"
	SyftConfigSnapshotPath  = ".prc/syft-config.yaml"
	syftArtifactID          = "syft-cyclonedx-json"
	syftArtifactMediaType   = "application/vnd.cyclonedx+json;version=1.7"
	syftSchemaURI           = "http://cyclonedx.org/schema/bom-1.7.schema.json"
)

var (
	syftConfig                    = []byte("{}\n")
	syftSerialPattern             = regexp.MustCompile(`^urn:uuid:[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)
	syftPackageURLPattern         = regexp.MustCompile(`^pkg:[^\s\x00]+$`)
	syftComponentTypePattern      = regexp.MustCompile(`^[a-z][a-z0-9-]{0,63}$`)
	syftAllowedTopLevelProperties = map[string]bool{
		"$schema": true, "bomFormat": true, "specVersion": true, "serialNumber": true,
		"version": true, "metadata": true, "components": true, "dependencies": true,
	}
)

func syftCommand() []string {
	return []string{
		"scan", "--quiet", "--config=/workspace/" + SyftConfigSnapshotPath,
		"--parallelism=1", "--output=cyclonedx-json", "dir:/workspace",
	}
}

func validateSyftManifest(manifest Manifest) error {
	if manifest.Protocol != SyftProtocolVersion || manifest.OutputSchema != SyftOutputSchemaVersion {
		return fmt.Errorf("syft adapter requires its exact protocol and output schema")
	}
	if manifest.Image != SyftImage || !slices.Equal(manifest.Command, syftCommand()) {
		return fmt.Errorf("syft adapter requires the reviewed immutable image and scanner-owned command")
	}
	if manifest.Tool.Name != "syft" || manifest.Tool.Version != SyftToolVersion ||
		manifest.Tool.Upstream != "https://github.com/anchore/syft" ||
		len(manifest.Tool.Formats) != 1 || manifest.Tool.Formats[0].Name != "cyclonedx-json" ||
		!slices.Equal(manifest.Tool.Formats[0].Versions, []string{"1.7"}) {
		return fmt.Errorf("syft adapter tool identity does not match the reviewed normalizer")
	}
	if !slices.Equal(manifest.ObservationKinds, []string{SyftObservationKind}) {
		return fmt.Errorf("syft adapter must declare only %s observations", SyftObservationKind)
	}
	if manifest.Capabilities.WriteScratch || !manifest.Capabilities.ChildProcesses {
		return fmt.Errorf("syft adapter requires no scratch and a bounded OS-task allowance")
	}
	if manifest.Resources.PIDs < 16 || manifest.Resources.PIDs > 256 ||
		manifest.Resources.TimeoutSeconds < 30 || manifest.Resources.MaxStdout < 1024*1024 {
		return fmt.Errorf("syft adapter resource limits cannot support or cannot safely bound the reviewed command")
	}
	return nil
}

func parseSyftOutput(data []byte, maxComponents int) (Transcript, map[string][]byte, error) {
	if len(bytes.TrimSpace(data)) == 0 {
		return Transcript{}, nil, fmt.Errorf("syft output is empty")
	}
	if !utf8.Valid(data) {
		return Transcript{}, nil, fmt.Errorf("syft output is not valid UTF-8")
	}
	if err := rejectDuplicateKeys(data); err != nil {
		return Transcript{}, nil, fmt.Errorf("syft output: %w", err)
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	var document map[string]any
	if err := decoder.Decode(&document); err != nil {
		return Transcript{}, nil, fmt.Errorf("decode syft CycloneDX output: %w", err)
	}
	var trailing any
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		if err == nil {
			return Transcript{}, nil, fmt.Errorf("syft output contains more than one JSON value")
		}
		return Transcript{}, nil, fmt.Errorf("decode trailing syft output: %w", err)
	}
	if document == nil {
		return Transcript{}, nil, fmt.Errorf("syft output must be a CycloneDX object")
	}
	for name := range document {
		if !syftAllowedTopLevelProperties[name] {
			return Transcript{}, nil, fmt.Errorf("syft output contains unexpected top-level property %q", name)
		}
	}
	if stringValue(document["$schema"]) != syftSchemaURI ||
		stringValue(document["bomFormat"]) != "CycloneDX" ||
		stringValue(document["specVersion"]) != "1.7" || stringValue(document["serialNumber"]) == "" ||
		document["version"] != json.Number("1") {
		return Transcript{}, nil, fmt.Errorf("syft output does not use the exact CycloneDX 1.7 document identity")
	}
	if !syftSerialPattern.MatchString(stringValue(document["serialNumber"])) {
		return Transcript{}, nil, fmt.Errorf("syft output has an invalid CycloneDX serial number")
	}
	metadata, ok := document["metadata"].(map[string]any)
	if !ok || metadata == nil {
		return Transcript{}, nil, fmt.Errorf("syft output requires CycloneDX metadata")
	}
	if err := validateSyftMetadata(metadata); err != nil {
		return Transcript{}, nil, err
	}
	components, ok := document["components"].([]any)
	if !ok || components == nil {
		return Transcript{}, nil, fmt.Errorf("syft output components must be an array")
	}
	if len(components) > maxComponents {
		return Transcript{}, nil, fmt.Errorf("syft output exceeds %d components", maxComponents)
	}
	packageCount, fileCount, err := validateAndNormalizeSyftComponents(components)
	if err != nil {
		return Transcript{}, nil, err
	}
	dependencyCount, edgeCount, err := validateAndNormalizeSyftDependencies(document, maxComponents)
	if err != nil {
		return Transcript{}, nil, err
	}

	// Syft emits a fresh UUID and timestamp for identical inputs. They are
	// optional in CycloneDX 1.7, so the reviewed normalizer removes them and the
	// metadata component's non-semantic bom-ref before canonical serialization.
	delete(document, "serialNumber")
	delete(metadata, "timestamp")
	if component, ok := metadata["component"].(map[string]any); ok {
		delete(component, "bom-ref")
	}
	canonical, err := json.Marshal(document)
	if err != nil {
		return Transcript{}, nil, fmt.Errorf("encode normalized CycloneDX artifact: %w", err)
	}
	canonical = append(canonical, '\n')
	digest := sha256.Sum256(canonical)
	descriptor := "sha256:" + hex.EncodeToString(digest[:])
	observationID := hex.EncodeToString(digest[:])
	transcript := Transcript{
		Logs: []Log{},
		Observations: []Observation{{
			ID: observationID, Kind: SyftObservationKind, Outcome: "value",
			Summary:   "Syft generated a normalized CycloneDX 1.7 SBOM for the sealed repository inventory.",
			Locations: []Location{},
			Data: map[string]any{
				"artifact_digest": descriptor, "component_count": len(components),
				"dependency_count": dependencyCount, "dependency_edge_count": edgeCount,
				"file_component_count": fileCount, "format": "CycloneDX JSON",
				"package_component_count": packageCount, "spec_version": "1.7", "tool_version": SyftToolVersion,
			},
		}},
		Artifacts: []Artifact{{
			ID: syftArtifactID, MediaType: syftArtifactMediaType,
			Digest: descriptor, Size: int64(len(canonical)), Path: "sbom.cdx.json",
		}},
		Summary: Summary{
			Type: "summary", Status: "completed",
			Counts: map[string]int{"logs": 0, "observations": 1, "artifacts": 1},
		},
	}
	return transcript, map[string][]byte{descriptor: canonical}, nil
}

func validateSyftMetadata(metadata map[string]any) error {
	for name := range metadata {
		if name != "timestamp" && name != "tools" && name != "component" {
			return fmt.Errorf("syft metadata contains unexpected property %q", name)
		}
	}
	timestamp := stringValue(metadata["timestamp"])
	parsed, err := time.Parse(time.RFC3339Nano, timestamp)
	if err != nil || parsed.Location() != time.UTC {
		return fmt.Errorf("syft metadata timestamp must be RFC 3339 UTC")
	}
	tools, ok := metadata["tools"].(map[string]any)
	if !ok || len(tools) != 1 {
		return fmt.Errorf("syft metadata must identify exactly one tool component collection")
	}
	items, ok := tools["components"].([]any)
	if !ok || len(items) != 1 {
		return fmt.Errorf("syft metadata must identify exactly one tool")
	}
	tool, ok := items[0].(map[string]any)
	if !ok || len(tool) != 4 || stringValue(tool["type"]) != "application" ||
		stringValue(tool["author"]) != "anchore" || stringValue(tool["name"]) != "syft" ||
		stringValue(tool["version"]) != SyftToolVersion {
		return fmt.Errorf("syft metadata tool identity is invalid")
	}
	component, ok := metadata["component"].(map[string]any)
	if !ok || len(component) != 3 || stringValue(component["type"]) != "file" ||
		stringValue(component["name"]) != "/workspace" || !validBoundedText(stringValue(component["bom-ref"]), 4096) {
		return fmt.Errorf("syft metadata source component is invalid")
	}
	return nil
}

func validateAndNormalizeSyftComponents(components []any) (int, int, error) {
	packageCount := 0
	fileCount := 0
	seenReferences := map[string]bool{}
	for index, raw := range components {
		component, ok := raw.(map[string]any)
		if !ok {
			return 0, 0, fmt.Errorf("syft component %d must be an object", index)
		}
		reference := stringValue(component["bom-ref"])
		componentType := stringValue(component["type"])
		name := stringValue(component["name"])
		if !validBoundedText(reference, 16*1024) || seenReferences[reference] ||
			!syftComponentTypePattern.MatchString(componentType) || !validBoundedText(name, 16*1024) {
			return 0, 0, fmt.Errorf("syft component %d has an invalid or duplicate identity", index)
		}
		seenReferences[reference] = true
		if version, exists := component["version"]; exists && !validBoundedText(stringValue(version), 4096) {
			return 0, 0, fmt.Errorf("syft component %d has an invalid version", index)
		}
		if purl, exists := component["purl"]; exists {
			if value := stringValue(purl); len(value) > 4096 || !syftPackageURLPattern.MatchString(value) {
				return 0, 0, fmt.Errorf("syft component %d has an invalid package URL", index)
			}
			packageCount++
		}
		if componentType == "file" {
			fileCount++
		}
		if properties, exists := component["properties"]; exists {
			items, ok := properties.([]any)
			if !ok || len(items) > 1024 {
				return 0, 0, fmt.Errorf("syft component %d properties are invalid", index)
			}
			for _, item := range items {
				property, ok := item.(map[string]any)
				if !ok || len(property) != 2 || !validBoundedText(stringValue(property["name"]), 4096) ||
					!validBoundedText(stringValue(property["value"]), 64*1024) {
					return 0, 0, fmt.Errorf("syft component %d contains an invalid property", index)
				}
			}
			sort.Slice(items, func(left, right int) bool { return canonicalLess(items[left], items[right]) })
		}
	}
	sort.Slice(components, func(left, right int) bool { return canonicalLess(components[left], components[right]) })
	return packageCount, fileCount, nil
}

func validateAndNormalizeSyftDependencies(document map[string]any, maxComponents int) (int, int, error) {
	raw, exists := document["dependencies"]
	if !exists {
		return 0, 0, nil
	}
	dependencies, ok := raw.([]any)
	if !ok || dependencies == nil || len(dependencies) > maxComponents {
		return 0, 0, fmt.Errorf("syft dependencies must be a bounded array")
	}
	edges := 0
	seen := map[string]bool{}
	for index, rawDependency := range dependencies {
		dependency, ok := rawDependency.(map[string]any)
		if !ok || len(dependency) != 2 {
			return 0, 0, fmt.Errorf("syft dependency %d must contain ref and dependsOn", index)
		}
		reference := stringValue(dependency["ref"])
		if !validBoundedText(reference, 16*1024) || seen[reference] {
			return 0, 0, fmt.Errorf("syft dependency %d has an invalid or duplicate ref", index)
		}
		seen[reference] = true
		dependsOn, ok := dependency["dependsOn"].([]any)
		if !ok || len(dependsOn) > maxComponents {
			return 0, 0, fmt.Errorf("syft dependency %d has an invalid dependsOn array", index)
		}
		for _, target := range dependsOn {
			if !validBoundedText(stringValue(target), 16*1024) {
				return 0, 0, fmt.Errorf("syft dependency %d has an invalid target", index)
			}
		}
		sort.Slice(dependsOn, func(left, right int) bool { return stringValue(dependsOn[left]) < stringValue(dependsOn[right]) })
		edges += len(dependsOn)
	}
	sort.Slice(dependencies, func(left, right int) bool { return canonicalLess(dependencies[left], dependencies[right]) })
	return len(dependencies), edges, nil
}

func stringValue(value any) string {
	result, _ := value.(string)
	return result
}

func validBoundedText(value string, maximum int) bool {
	return value != "" && len(value) <= maximum && utf8.ValidString(value) && !strings.ContainsRune(value, '\x00')
}

func canonicalLess(left, right any) bool {
	leftJSON, _ := json.Marshal(left)
	rightJSON, _ := json.Marshal(right)
	return bytes.Compare(leftJSON, rightJSON) < 0
}
