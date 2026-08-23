package engine

import (
	"fmt"
	"path/filepath"
	"slices"

	"github.com/MarinJursic/production-readiness-checklist/scanner/catalog"
	"github.com/MarinJursic/production-readiness-checklist/scanner/finding"
	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
)

func buildFindings(catalog *catalog.Catalog, run model.RunResult) ([]model.Finding, error) {
	findings := make([]model.Finding, 0)
	for _, result := range run.Results {
		if result.Assessment != "fail" {
			continue
		}
		assertion, ok := catalog.Assertions[result.AssertionID]
		if !ok {
			return nil, fmt.Errorf("cannot create finding for unknown assertion %s", result.AssertionID)
		}
		evidenceIDs := make([]string, 0, len(result.EvidenceObserved))
		for _, evidence := range result.EvidenceObserved {
			evidenceIDs = append(evidenceIDs, evidence.ID)
		}
		item, err := finding.New(finding.Input{
			AssertionID: result.AssertionID, ControlIDs: result.ControlIDs,
			Title: assertion.Title, Summary: result.Summary, Severity: result.Severity,
			Gate: result.Gate, RemediationClass: result.RemediationClass,
			Subject: model.FindingSubject{
				Kind: subjectKind(run.Inventory), ID: subjectID(run.Inventory), InventoryDigest: run.Inventory.Digest,
			},
			Locations:   findingLocations(assertion, result, run.Inventory, run.AdapterExecutions),
			EvidenceIDs: evidenceIDs,
		})
		if err != nil {
			return nil, fmt.Errorf("create finding for %s: %w", result.AssertionID, err)
		}
		findings = append(findings, item)
	}
	return findings, nil
}

func subjectKind(inventory model.Inventory) string {
	if inventory.DeclaredScope != nil && inventory.DeclaredScope.ProjectID != "" {
		return "project"
	}
	return "repository"
}

func subjectID(inventory model.Inventory) string {
	if inventory.DeclaredScope != nil && inventory.DeclaredScope.ProjectID != "" {
		return inventory.DeclaredScope.ProjectID
	}
	return inventory.TargetName
}

func findingLocations(
	assertion model.Assertion,
	result model.AssertionResult,
	inventory model.Inventory,
	executions []model.AdapterExecution,
) []model.FindingLocation {
	inventoryPaths := make(map[string]bool, len(inventory.Files))
	for _, file := range inventory.Files {
		inventoryPaths[file.Path] = true
	}
	locations := make([]model.FindingLocation, 0)
	detailedPaths := map[string]bool{}
	for _, location := range result.Locations {
		path := filepath.ToSlash(filepath.Clean(filepath.FromSlash(location.Path)))
		if path != location.Path || !inventoryPaths[path] || location.Line < 0 || location.Column < 0 ||
			(location.Column > 0 && location.Line == 0) {
			continue
		}
		locations = append(locations, location)
		detailedPaths[path] = true
	}
	for _, evidence := range result.EvidenceObserved {
		if inventoryPaths[evidence.Source] && !detailedPaths[evidence.Source] {
			locations = append(locations, model.FindingLocation{Path: evidence.Source})
		}
	}
	bindings, err := assertionAdapterBindings(assertion)
	if err != nil {
		return locations
	}
	for _, binding := range bindings {
		for _, execution := range executions {
			if execution.AdapterID != binding.AdapterID || execution.ManifestSHA256 != binding.ManifestSHA256 {
				continue
			}
			for _, observation := range execution.Transcript.Observations {
				if observation.Kind != binding.ObservationKind || !slices.Contains(binding.FailOutcomes, observation.Outcome) {
					continue
				}
				for _, location := range observation.Locations {
					path := filepath.ToSlash(filepath.Clean(filepath.FromSlash(location.Path)))
					if path != location.Path || !inventoryPaths[path] {
						continue
					}
					locations = append(locations, model.FindingLocation{
						Path: path, Line: location.Line, Column: location.Column,
					})
				}
			}
		}
	}
	return locations
}
