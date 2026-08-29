package controlruntime

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/MarinJursic/production-readiness-checklist/scanner/controlprogram"
	"github.com/MarinJursic/production-readiness-checklist/scanner/controlprogramcatalog"
)

type fixtureProvider struct {
	id        string
	authority controlprogram.Authority
	value     bool
	err       error
}

func (provider fixtureProvider) ID() string                          { return provider.id }
func (provider fixtureProvider) Authority() controlprogram.Authority { return provider.authority }
func (provider fixtureProvider) Collect(_ context.Context, request Request) (controlprogram.Evidence, error) {
	if provider.err != nil {
		return controlprogram.Evidence{}, provider.err
	}
	return evidenceFor(request, provider.value)
}

func fixtureTemplate() controlprogramcatalog.Template {
	return controlprogramcatalog.Template{
		TemplateID: strings.Repeat("1", 64), ControlID: "PRC-01-001", ControlRevision: 1,
		ControlSemanticSHA256: strings.Repeat("2", 64), ClauseID: strings.Repeat("3", 64),
		ClauseOrdinal: 1, ClauseStatementSHA256: strings.Repeat("4", 64),
		ImplementationID: "prc.check.fixture@0.1", ImplementationContractSHA256: strings.Repeat("5", 64),
		RequiredAuthority: controlprogram.AuthorityRepository,
		CollectorContract: controlprogramcatalog.CollectorContract{CollectorID: "prc.collect.fixture@0.1"},
		Predicate:         controlprogram.Expression{Op: controlprogram.OpBooleanEqual, Fact: "fixture.result", Boolean: boolPointer(true)},
	}
}

func fixtureBinding() controlprogramcatalog.RuntimeBinding {
	return controlprogramcatalog.RuntimeBinding{
		SubjectID: "fixture", Subjects: []string{"fixture"}, InventorySHA256: strings.Repeat("6", 64),
		ApplicabilityProofContractSHA256: strings.Repeat("7", 64), MaximumEvidenceAgeSeconds: 3600,
	}
}

func evidenceFor(request Request, value bool) (controlprogram.Evidence, error) {
	return NewApplicableEvidence(request, "fixture-evidence", map[string]controlprogram.Fact{
		"fixture.result": {Type: controlprogram.FactBoolean, Complete: true, Boolean: boolPointer(value)},
	}, true)
}

func boolPointer(value bool) *bool { return &value }

func TestRegistryRejectsDuplicateAndWrongAuthority(t *testing.T) {
	provider := fixtureProvider{id: "prc.collect.fixture@0.1", authority: controlprogram.AuthorityRepository}
	if _, err := NewRegistry(provider, provider); err == nil {
		t.Fatal("duplicate provider was accepted")
	}
	if _, err := NewRegistry(fixtureProvider{id: "bad", authority: "declared"}); err == nil {
		t.Fatal("unsupported authority was accepted")
	}
}

func TestEvaluateOwnsPassAndFail(t *testing.T) {
	now := time.Date(2026, 8, 29, 12, 5, 0, 0, time.UTC)
	for _, item := range []struct {
		name    string
		value   bool
		status  ExecutionStatus
		outcome controlprogram.Outcome
	}{
		{name: "pass", value: true, status: StatusPassed, outcome: controlprogram.OutcomePass},
		{name: "fail", value: false, status: StatusFailed, outcome: controlprogram.OutcomeFail},
	} {
		t.Run(item.name, func(t *testing.T) {
			registry, err := NewRegistry(fixtureProvider{id: "prc.collect.fixture@0.1", authority: controlprogram.AuthorityRepository, value: item.value})
			if err != nil {
				t.Fatal(err)
			}
			result := Evaluate(context.Background(), fixtureTemplate(), fixtureBinding(), registry, now)
			if result.Status != item.status || result.Outcome != item.outcome || result.ProgramSHA256 == "" || result.EvidenceSHA256 == "" {
				t.Fatalf("unexpected execution: %+v", result)
			}
			sealed, ok := result.SealedEvidence()
			if !ok || controlprogram.EvidenceSHA256(sealed) != result.EvidenceSHA256 {
				t.Fatal("runtime did not retain replayable validated evidence")
			}
			sealed.Facts["fixture.result"] = controlprogram.Fact{}
			second, ok := result.SealedEvidence()
			if !ok || controlprogram.EvidenceSHA256(second) != result.EvidenceSHA256 {
				t.Fatal("sealed evidence was not returned defensively")
			}
		})
	}
}

func TestEvaluateFailsClosedAtEveryBoundary(t *testing.T) {
	now := time.Date(2026, 8, 29, 12, 5, 0, 0, time.UTC)
	template := fixtureTemplate()
	if result := Evaluate(context.Background(), template, controlprogramcatalog.RuntimeBinding{}, nil, now); result.Status != StatusBlockedBinding {
		t.Fatalf("invalid binding status = %s", result.Status)
	}
	if result := Evaluate(context.Background(), template, fixtureBinding(), nil, now); result.Status != StatusBlockedProvider {
		t.Fatalf("missing provider status = %s", result.Status)
	}
	wrong, _ := NewRegistry(fixtureProvider{id: "prc.collect.fixture@0.1", authority: controlprogram.AuthorityArtifact})
	if result := Evaluate(context.Background(), template, fixtureBinding(), wrong, now); result.Status != StatusBlockedAuthority {
		t.Fatalf("wrong authority status = %s", result.Status)
	}
	failed, _ := NewRegistry(fixtureProvider{id: "prc.collect.fixture@0.1", authority: controlprogram.AuthorityRepository, err: errors.New("secret provider detail")})
	if result := Evaluate(context.Background(), template, fixtureBinding(), failed, now); result.Status != StatusBlockedCollection {
		t.Fatalf("provider failure status = %s", result.Status)
	}
	canceled, cancel := context.WithCancel(context.Background())
	cancel()
	valid, _ := NewRegistry(fixtureProvider{id: "prc.collect.fixture@0.1", authority: controlprogram.AuthorityRepository})
	if result := Evaluate(canceled, template, fixtureBinding(), valid, now); result.Status != StatusBlockedCanceled {
		t.Fatalf("canceled status = %s", result.Status)
	}
}

func TestEvaluateCatalogPreservesAllTemplates(t *testing.T) {
	// A nil catalog is the only empty result; a loaded catalog is covered by
	// integration tests in controlprogramcatalog/fullscan.
	if results := EvaluateCatalog(context.Background(), nil, nil, nil, time.Now()); len(results) != 0 {
		t.Fatalf("nil catalog produced %d results", len(results))
	}
}
