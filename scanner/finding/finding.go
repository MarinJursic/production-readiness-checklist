package finding

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strings"

	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
)

var digestPattern = regexp.MustCompile(`^[0-9a-f]{64}$`)

var (
	assertionIDPattern = regexp.MustCompile(`^PRC-A-[A-Z0-9]+-[0-9]{3}$`)
	controlIDPattern   = regexp.MustCompile(`^(?:PRC-[0-9]{2}-[0-9]{3}|USEQ-[A-F0-9]{8})$`)
)

type Input struct {
	AssertionID      string
	ControlIDs       []string
	Title            string
	Summary          string
	Severity         string
	Gate             string
	RemediationClass string
	Subject          model.FindingSubject
	Locations        []model.FindingLocation
	EvidenceIDs      []string
}

func New(input Input) (model.Finding, error) {
	item := model.Finding{
		SchemaVersion: model.FindingSchema,
		AssertionID:   input.AssertionID, ControlIDs: sortedUnique(input.ControlIDs),
		Title: input.Title, Summary: input.Summary, Severity: input.Severity,
		Gate: input.Gate, RemediationClass: input.RemediationClass, Subject: input.Subject,
		Locations: canonicalLocations(input.Locations), EvidenceIDs: sortedUnique(input.EvidenceIDs),
	}
	if err := validateFields(item); err != nil {
		return model.Finding{}, err
	}
	item.Fingerprint = fingerprint(item)
	item.ID = identity(item)
	return item, nil
}

func Validate(item model.Finding) error {
	if err := validateFields(item); err != nil {
		return err
	}
	if !slices.Equal(item.ControlIDs, sortedUnique(item.ControlIDs)) ||
		!slices.Equal(item.EvidenceIDs, sortedUnique(item.EvidenceIDs)) ||
		!slices.Equal(item.Locations, canonicalLocations(item.Locations)) {
		return fmt.Errorf("finding arrays are not sorted and unique")
	}
	if !digestPattern.MatchString(item.Fingerprint) || item.Fingerprint != fingerprint(item) {
		return fmt.Errorf("finding fingerprint does not match its stable subject")
	}
	if !digestPattern.MatchString(item.ID) || item.ID != identity(item) {
		return fmt.Errorf("finding ID does not match its content")
	}
	return nil
}

func validateFields(item model.Finding) error {
	if item.SchemaVersion != model.FindingSchema || !assertionIDPattern.MatchString(item.AssertionID) ||
		strings.TrimSpace(item.Title) == "" || strings.TrimSpace(item.Summary) == "" ||
		strings.TrimSpace(item.Subject.Kind) == "" || strings.TrimSpace(item.Subject.ID) == "" ||
		!digestPattern.MatchString(item.Subject.InventoryDigest) {
		return fmt.Errorf("finding has incomplete identity or text")
	}
	if len(item.ControlIDs) == 0 {
		return fmt.Errorf("finding must map to at least one control")
	}
	if !slices.Contains([]string{"repository", "project", "component", "artifact", "environment", "endpoint"}, item.Subject.Kind) {
		return fmt.Errorf("finding has an invalid subject kind")
	}
	if !slices.Contains([]string{"info", "low", "medium", "high", "critical"}, item.Severity) ||
		!slices.Contains([]string{"advisory", "required", "no-go"}, item.Gate) ||
		!slices.Contains([]string{"R0", "R1", "R2", "R3", "R4", "R5", "R6"}, item.RemediationClass) {
		return fmt.Errorf("finding has an invalid severity, gate, or remediation class")
	}
	for _, controlID := range item.ControlIDs {
		if !controlIDPattern.MatchString(controlID) {
			return fmt.Errorf("finding contains an invalid control ID")
		}
	}
	for _, evidenceID := range item.EvidenceIDs {
		if !digestPattern.MatchString(evidenceID) {
			return fmt.Errorf("finding contains an invalid evidence ID")
		}
	}
	for _, location := range item.Locations {
		clean := filepath.ToSlash(filepath.Clean(filepath.FromSlash(location.Path)))
		if strings.TrimSpace(location.Path) == "" || filepath.IsAbs(location.Path) || strings.Contains(location.Path, "\\") ||
			clean != location.Path || clean == "." || strings.HasPrefix(clean, "../") ||
			location.Line < 0 || location.Column < 0 || (location.Column > 0 && location.Line == 0) {
			return fmt.Errorf("finding contains an invalid location")
		}
	}
	return nil
}

func fingerprint(item model.Finding) string {
	payload, _ := json.Marshal(struct {
		SchemaVersion string                  `json:"schema_version"`
		AssertionID   string                  `json:"assertion_id"`
		SubjectKind   string                  `json:"subject_kind"`
		SubjectID     string                  `json:"subject_id"`
		Locations     []model.FindingLocation `json:"locations"`
	}{
		SchemaVersion: "prc.finding-fingerprint/v0.1", AssertionID: item.AssertionID,
		SubjectKind: item.Subject.Kind, SubjectID: item.Subject.ID, Locations: item.Locations,
	})
	digest := sha256.Sum256(payload)
	return hex.EncodeToString(digest[:])
}

func identity(item model.Finding) string {
	item.ID = ""
	payload, _ := json.Marshal(item)
	digest := sha256.Sum256(payload)
	return hex.EncodeToString(digest[:])
}

func sortedUnique(values []string) []string {
	seen := map[string]bool{}
	result := make([]string, 0, len(values))
	for _, value := range values {
		if !seen[value] {
			seen[value] = true
			result = append(result, value)
		}
	}
	sort.Strings(result)
	return result
}

func canonicalLocations(values []model.FindingLocation) []model.FindingLocation {
	seen := map[model.FindingLocation]bool{}
	result := make([]model.FindingLocation, 0, len(values))
	for _, value := range values {
		if !seen[value] {
			seen[value] = true
			result = append(result, value)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Path != result[j].Path {
			return result[i].Path < result[j].Path
		}
		if result[i].Line != result[j].Line {
			return result[i].Line < result[j].Line
		}
		return result[i].Column < result[j].Column
	})
	return result
}
