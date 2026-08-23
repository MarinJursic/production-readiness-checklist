package adapter

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
)

const (
	ResolutionSourceExplicitLocal = "explicit-local"
	ResolutionSourceRegistry      = "registry"
	ResolutionTrustLocalExplicit  = "local-explicit"
)

func ManifestDigest(manifest Manifest) (string, error) {
	if err := manifest.Validate(); err != nil {
		return "", err
	}
	payload, err := json.Marshal(manifest)
	if err != nil {
		return "", fmt.Errorf("encode adapter manifest identity: %w", err)
	}
	digest := sha256.Sum256(payload)
	return hex.EncodeToString(digest[:]), nil
}

func executionID(execution model.AdapterExecution) (string, error) {
	execution.ExecutionID = ""
	payload, err := json.Marshal(execution)
	if err != nil {
		return "", fmt.Errorf("encode adapter execution identity: %w", err)
	}
	digest := sha256.Sum256(payload)
	return hex.EncodeToString(digest[:]), nil
}

func BindExecution(runID string, subject Subject, manifest Manifest, output RunOutput) (model.AdapterExecution, error) {
	return BindExecutionWithResolution(runID, subject, manifest, model.AdapterResolution{
		Source: ResolutionSourceExplicitLocal, PublisherID: manifest.Publisher.ID, Trust: ResolutionTrustLocalExplicit,
	}, output)
}

// BindExecutionWithResolution makes the authorization provenance part of the
// content-addressed execution. Callers resolving an adapter from a Registry
// must use the resolution returned by Registry.Resolve.
func BindExecutionWithResolution(
	runID string,
	subject Subject,
	manifest Manifest,
	resolution model.AdapterResolution,
	output RunOutput,
) (model.AdapterExecution, error) {
	if !hexDigestPattern.MatchString(runID) {
		return model.AdapterExecution{}, fmt.Errorf("adapter run ID must be a lowercase SHA-256 digest")
	}
	if strings.TrimSpace(subject.TargetName) == "" || !hexDigestPattern.MatchString(subject.InventoryDigest) {
		return model.AdapterExecution{}, fmt.Errorf("adapter subject requires a target name and inventory digest")
	}
	if subject.TargetCommit != "" && !commitPattern.MatchString(subject.TargetCommit) {
		return model.AdapterExecution{}, fmt.Errorf("adapter subject commit is invalid")
	}
	manifestDigest, err := ManifestDigest(manifest)
	if err != nil {
		return model.AdapterExecution{}, err
	}
	if output.StartedAt.IsZero() || output.CompletedAt.Before(output.StartedAt) || output.DurationMS < 0 {
		return model.AdapterExecution{}, fmt.Errorf("adapter execution timestamps are invalid")
	}
	if err := ValidateTranscriptContract(manifest, output.Transcript); err != nil {
		return model.AdapterExecution{}, err
	}
	if err := validateResolutionForPublisher(resolution, manifest.Publisher.ID); err != nil {
		return model.AdapterExecution{}, err
	}
	dataInputs := make([]model.AdapterDataInput, len(output.DataInputs))
	copy(dataInputs, output.DataInputs)
	execution := model.AdapterExecution{
		SchemaVersion: model.AdapterExecutionSchema,
		AdapterRunID:  runID, AdapterID: manifest.ID, ManifestSHA256: manifestDigest, Image: manifest.Image,
		Resolution: resolution, DataInputs: dataInputs,
		Subject: model.AdapterSubject{
			TargetName: subject.TargetName, TargetCommit: subject.TargetCommit, InventoryDigest: subject.InventoryDigest,
		},
		StartedAt: output.StartedAt, CompletedAt: output.CompletedAt, DurationMS: output.DurationMS,
		DiagnosticsSHA256: output.DiagnosticsSHA256, DiagnosticsBytes: output.DiagnosticsBytes,
		Transcript: transcriptToModel(output.Transcript),
	}
	execution.ExecutionID, err = executionID(execution)
	if err != nil {
		return model.AdapterExecution{}, err
	}
	if err := ValidateExecution(execution); err != nil {
		return model.AdapterExecution{}, err
	}
	return execution, nil
}

// ValidateTranscriptContract binds structurally valid protocol output to the
// exact current-engine manifest contract. An adapter cannot introduce a new
// observation class merely by emitting it.
func ValidateTranscriptContract(manifest Manifest, transcript Transcript) error {
	if err := manifest.ValidateForCurrentEngine(); err != nil {
		return err
	}
	if err := validateTranscript(transcript); err != nil {
		return err
	}
	allowed := make(map[string]bool, len(manifest.ObservationKinds))
	for _, kind := range manifest.ObservationKinds {
		allowed[kind] = true
	}
	for _, observation := range transcript.Observations {
		if !allowed[observation.Kind] {
			return fmt.Errorf("adapter emitted undeclared observation kind %q", observation.Kind)
		}
	}
	return nil
}

func ValidateExecution(execution model.AdapterExecution) error {
	if execution.SchemaVersion != model.AdapterExecutionSchema && execution.SchemaVersion != "prc.adapter-execution/v0.2" &&
		execution.SchemaVersion != "prc.adapter-execution/v0.1" {
		return fmt.Errorf("unsupported adapter execution schema %q", execution.SchemaVersion)
	}
	if !hexDigestPattern.MatchString(execution.ExecutionID) || !hexDigestPattern.MatchString(execution.AdapterRunID) ||
		!hexDigestPattern.MatchString(execution.ManifestSHA256) || !hexDigestPattern.MatchString(execution.DiagnosticsSHA256) {
		return fmt.Errorf("adapter execution contains an invalid digest")
	}
	if !adapterIDPattern.MatchString(execution.AdapterID) || !imagePattern.MatchString(execution.Image) {
		return fmt.Errorf("adapter execution identity is invalid")
	}
	if execution.SchemaVersion == model.AdapterExecutionSchema || execution.SchemaVersion == "prc.adapter-execution/v0.2" {
		if err := validateResolution(execution.Resolution); err != nil {
			return fmt.Errorf("adapter execution resolution: %w", err)
		}
	} else if execution.Resolution != (model.AdapterResolution{}) {
		return fmt.Errorf("legacy adapter execution cannot contain resolution provenance")
	}
	if execution.SchemaVersion == model.AdapterExecutionSchema {
		if err := validateExecutionDataInputs(execution.DataInputs); err != nil {
			return err
		}
	} else if execution.DataInputs != nil {
		return fmt.Errorf("legacy adapter execution cannot contain data inputs")
	}
	if strings.TrimSpace(execution.Subject.TargetName) == "" || !hexDigestPattern.MatchString(execution.Subject.InventoryDigest) {
		return fmt.Errorf("adapter execution subject is invalid")
	}
	if execution.Subject.TargetCommit != "" && !commitPattern.MatchString(execution.Subject.TargetCommit) {
		return fmt.Errorf("adapter execution subject commit is invalid")
	}
	if execution.StartedAt.IsZero() || execution.CompletedAt.Before(execution.StartedAt) ||
		execution.DurationMS < 0 || execution.DiagnosticsBytes < 0 {
		return fmt.Errorf("adapter execution metadata is invalid")
	}
	if execution.CompletedAt.Sub(execution.StartedAt).Milliseconds() != execution.DurationMS {
		return fmt.Errorf("adapter execution duration does not match its timestamps")
	}
	if err := validateModelTranscript(execution.Transcript); err != nil {
		return err
	}
	expected, err := executionID(execution)
	if err != nil {
		return err
	}
	if expected != execution.ExecutionID {
		return fmt.Errorf("adapter execution ID does not match its content")
	}
	return nil
}

func validateExecutionDataInputs(inputs []model.AdapterDataInput) error {
	if inputs == nil {
		return fmt.Errorf("adapter execution data inputs cannot be null")
	}
	seen := map[string]bool{}
	for _, input := range inputs {
		if !observationKindPattern.MatchString(input.Name) || seen[input.Name] ||
			input.Destination != "/prc-inputs/"+input.Name || !hexDigestPattern.MatchString(input.SHA256) ||
			input.Files < 1 || input.Bytes < 1 {
			return fmt.Errorf("adapter execution contains an invalid or duplicate data input")
		}
		seen[input.Name] = true
	}
	return nil
}

func validateResolutionForPublisher(resolution model.AdapterResolution, manifestPublisherID string) error {
	if err := validateResolution(resolution); err != nil {
		return err
	}
	if resolution.PublisherID != manifestPublisherID {
		return fmt.Errorf("resolution publisher does not match the manifest publisher")
	}
	return nil
}

func validateResolution(resolution model.AdapterResolution) error {
	if !publisherIDPattern.MatchString(resolution.PublisherID) {
		return fmt.Errorf("resolution publisher is invalid")
	}
	switch resolution.Source {
	case ResolutionSourceExplicitLocal:
		if resolution.Trust != ResolutionTrustLocalExplicit || resolution.RegistryID != "" ||
			resolution.RegistryRevision != 0 || resolution.RegistryDigest != "" {
			return fmt.Errorf("explicit-local resolution has inconsistent trust or registry fields")
		}
	case ResolutionSourceRegistry:
		if !registryIDPattern.MatchString(resolution.RegistryID) || resolution.RegistryRevision < 1 ||
			!hexDigestPattern.MatchString(resolution.RegistryDigest) ||
			(resolution.Trust != "first-party-sandboxed" && resolution.Trust != "verified-community" &&
				resolution.Trust != "unverified-community" && resolution.Trust != "local") {
			return fmt.Errorf("registry resolution has invalid registry identity or trust")
		}
	default:
		return fmt.Errorf("unsupported resolution source %q", resolution.Source)
	}
	return nil
}

func validateTranscript(transcript Transcript) error {
	for _, item := range transcript.Logs {
		if err := validateLog(item); err != nil {
			return fmt.Errorf("adapter log: %w", err)
		}
	}
	observationIDs := map[string]bool{}
	for _, item := range transcript.Observations {
		if err := validateObservation(item); err != nil {
			return fmt.Errorf("adapter observation: %w", err)
		}
		if observationIDs[item.ID] {
			return fmt.Errorf("adapter execution has duplicate observation ID %q", item.ID)
		}
		observationIDs[item.ID] = true
	}
	artifactIDs := map[string]bool{}
	for _, item := range transcript.Artifacts {
		if err := validateArtifact(item); err != nil {
			return fmt.Errorf("adapter artifact: %w", err)
		}
		if artifactIDs[item.ID] {
			return fmt.Errorf("adapter execution has duplicate artifact ID %q", item.ID)
		}
		artifactIDs[item.ID] = true
	}
	if err := validateSummary(transcript.Summary); err != nil {
		return fmt.Errorf("adapter summary: %w", err)
	}
	actual := map[string]int{"logs": len(transcript.Logs), "observations": len(transcript.Observations), "artifacts": len(transcript.Artifacts)}
	for name, count := range actual {
		if declared, ok := transcript.Summary.Counts[name]; ok && declared != count {
			return fmt.Errorf("adapter summary count %q is %d, observed %d", name, declared, count)
		}
	}
	return nil
}

func validateModelTranscript(transcript model.AdapterTranscript) error {
	if transcript.Logs == nil || transcript.Observations == nil || transcript.Artifacts == nil {
		return fmt.Errorf("bound adapter transcript arrays cannot be null")
	}
	observationIDs := map[string]bool{}
	for _, item := range transcript.Logs {
		if err := validateLog(Log{Type: "log", Level: item.Level, Message: item.Message}); err != nil {
			return fmt.Errorf("bound adapter log: %w", err)
		}
	}
	for _, item := range transcript.Observations {
		if item.Locations == nil {
			return fmt.Errorf("bound adapter observation locations cannot be null")
		}
		locations := make([]Location, 0, len(item.Locations))
		for _, location := range item.Locations {
			locations = append(locations, Location{Path: location.Path, Line: location.Line, Column: location.Column})
		}
		observation := Observation{
			ID: item.ID, Kind: item.Kind, Outcome: item.Outcome, Summary: item.Summary,
			Locations: locations, Data: item.Data,
		}
		if err := validateObservation(observation); err != nil {
			return fmt.Errorf("bound adapter observation: %w", err)
		}
		if observationIDs[item.ID] {
			return fmt.Errorf("bound adapter execution has duplicate observation ID %q", item.ID)
		}
		observationIDs[item.ID] = true
	}
	artifactIDs := map[string]bool{}
	for _, item := range transcript.Artifacts {
		artifact := Artifact{ID: item.ID, MediaType: item.MediaType, Digest: item.Digest, Size: item.Size, Path: item.Path}
		if err := validateArtifact(artifact); err != nil {
			return fmt.Errorf("bound adapter artifact: %w", err)
		}
		if artifactIDs[item.ID] {
			return fmt.Errorf("bound adapter execution has duplicate artifact ID %q", item.ID)
		}
		artifactIDs[item.ID] = true
	}
	summary := Summary{Type: "summary", Status: transcript.Summary.Status, Counts: transcript.Summary.Counts, Reason: transcript.Summary.Reason}
	if err := validateSummary(summary); err != nil {
		return fmt.Errorf("bound adapter summary: %w", err)
	}
	actual := map[string]int{"logs": len(transcript.Logs), "observations": len(transcript.Observations), "artifacts": len(transcript.Artifacts)}
	for name, count := range actual {
		if declared, ok := transcript.Summary.Counts[name]; ok && declared != count {
			return fmt.Errorf("bound adapter summary count %q is %d, observed %d", name, declared, count)
		}
	}
	return nil
}

func transcriptToModel(transcript Transcript) model.AdapterTranscript {
	result := model.AdapterTranscript{
		Logs: []model.AdapterLog{}, Observations: []model.AdapterObservation{}, Artifacts: []model.AdapterArtifact{},
		Summary: model.AdapterSummary{Status: transcript.Summary.Status, Counts: transcript.Summary.Counts, Reason: transcript.Summary.Reason},
	}
	for _, item := range transcript.Logs {
		result.Logs = append(result.Logs, model.AdapterLog{Level: item.Level, Message: item.Message})
	}
	for _, item := range transcript.Observations {
		locations := make([]model.AdapterLocation, 0, len(item.Locations))
		for _, location := range item.Locations {
			locations = append(locations, model.AdapterLocation{Path: location.Path, Line: location.Line, Column: location.Column})
		}
		result.Observations = append(result.Observations, model.AdapterObservation{
			ID: item.ID, Kind: item.Kind, Outcome: item.Outcome, Summary: item.Summary,
			Locations: locations, Data: item.Data,
		})
	}
	for _, item := range transcript.Artifacts {
		result.Artifacts = append(result.Artifacts, model.AdapterArtifact{
			ID: item.ID, MediaType: item.MediaType, Digest: item.Digest, Size: item.Size, Path: item.Path,
		})
	}
	if result.Summary.Counts == nil {
		result.Summary.Counts = map[string]int{}
	}
	return result
}
