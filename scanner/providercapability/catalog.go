// Package providercapability owns the single checked list of deterministic
// evidence collectors shipped in this scanner build.
package providercapability

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"sort"

	"github.com/MarinJursic/production-readiness-checklist/scanner/controlprogram"
)

const SchemaVersion = "prc.provider-capabilities/v0.1"

var (
	collectorPattern = regexp.MustCompile(`^prc\.collect\.[a-z0-9.-]+@[0-9]+\.[0-9]+$`)
	controlPattern   = regexp.MustCompile(`^(?:PRC-[0-9]{2}-[0-9]{3}|USEQ-[A-F0-9]{8})$`)
)

//go:embed providers.json
var documentBytes []byte

// Capability binds one shipped collector to the only reviewed clause and
// evidence authority it is allowed to serve.
type Capability struct {
	CollectorID   string                   `json:"collector_id"`
	ControlID     string                   `json:"control_id"`
	ClauseOrdinal int                      `json:"clause_ordinal"`
	Authority     controlprogram.Authority `json:"authority"`
}

type document struct {
	SchemaVersion string       `json:"schema_version"`
	Capabilities  []Capability `json:"capabilities"`
}

// Load parses and validates the embedded build manifest. It rejects duplicate,
// unordered, unknown, and incomplete declarations.
func Load() ([]Capability, error) {
	decoder := json.NewDecoder(bytes.NewReader(documentBytes))
	decoder.DisallowUnknownFields()
	var value document
	if err := decoder.Decode(&value); err != nil {
		return nil, fmt.Errorf("decode provider capability manifest: %w", err)
	}
	var trailing any
	if err := decoder.Decode(&trailing); err != io.EOF {
		return nil, fmt.Errorf("decode provider capability manifest: trailing JSON")
	}
	if value.SchemaVersion != SchemaVersion || len(value.Capabilities) > 765 {
		return nil, fmt.Errorf("provider capability manifest envelope is invalid")
	}
	seenCollectors := map[string]bool{}
	seenClauses := map[string]bool{}
	previous := ""
	for _, capability := range value.Capabilities {
		key := fmt.Sprintf("%s/%02d/%s", capability.ControlID, capability.ClauseOrdinal, capability.CollectorID)
		clauseKey := fmt.Sprintf("%s/%02d", capability.ControlID, capability.ClauseOrdinal)
		if !collectorPattern.MatchString(capability.CollectorID) || !controlPattern.MatchString(capability.ControlID) ||
			capability.ClauseOrdinal < 1 || capability.ClauseOrdinal > 50 || !validAuthority(capability.Authority) ||
			seenCollectors[capability.CollectorID] || seenClauses[clauseKey] || (previous != "" && key <= previous) {
			return nil, fmt.Errorf("provider capability manifest entry %q is invalid", key)
		}
		seenCollectors[capability.CollectorID] = true
		seenClauses[clauseKey] = true
		previous = key
	}
	result := append([]Capability(nil), value.Capabilities...)
	return result, nil
}

// IDs returns the checked collector identities in stable order.
func IDs() ([]string, error) {
	capabilities, err := Load()
	if err != nil {
		return nil, err
	}
	ids := make([]string, len(capabilities))
	for index, capability := range capabilities {
		ids[index] = capability.CollectorID
	}
	sort.Strings(ids)
	return ids, nil
}

func validAuthority(authority controlprogram.Authority) bool {
	switch authority {
	case controlprogram.AuthorityRepository, controlprogram.AuthorityArtifact,
		controlprogram.AuthorityExecuted, controlprogram.AuthorityEnvironment,
		controlprogram.AuthorityExternalRegistry, controlprogram.AuthorityStructuredRecord:
		return true
	default:
		return false
	}
}
