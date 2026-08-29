// Package controlruntime binds reviewed deterministic program templates to
// registered evidence providers and evaluates the resulting sealed evidence.
// Providers produce typed facts; only controlprogram.Evaluate owns Pass/Fail.
package controlruntime

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"time"

	"github.com/MarinJursic/production-readiness-checklist/scanner/controlprogram"
	"github.com/MarinJursic/production-readiness-checklist/scanner/controlprogramcatalog"
)

// ExecutionStatus describes how far one exact clause evaluation progressed.
type ExecutionStatus string

const (
	StatusPassed            ExecutionStatus = "passed"
	StatusFailed            ExecutionStatus = "failed"
	StatusNotApplicable     ExecutionStatus = "not_applicable"
	StatusBlockedBinding    ExecutionStatus = "blocked_binding"
	StatusBlockedProvider   ExecutionStatus = "blocked_provider_unregistered"
	StatusBlockedAuthority  ExecutionStatus = "blocked_provider_authority"
	StatusBlockedCollection ExecutionStatus = "blocked_collection"
	StatusBlockedEvidence   ExecutionStatus = "blocked_evidence"
	StatusBlockedCanceled   ExecutionStatus = "blocked_canceled"
)

// Request is the immutable input a provider receives. The Program already
// contains scanner-sealed policy, inventory, applicability, and freshness.
// A provider cannot supply or replace those parameters.
type Request struct {
	Template       controlprogramcatalog.Template
	Program        controlprogram.Program
	EvaluationTime time.Time
}

// Provider obtains typed raw evidence for one exact collector contract.
// Implementations must not make a compliance decision.
type Provider interface {
	ID() string
	Authority() controlprogram.Authority
	Collect(context.Context, Request) (controlprogram.Evidence, error)
}

// Registry is an immutable exact collector-ID to provider mapping.
type Registry struct {
	providers map[string]Provider
}

// NewRegistry rejects nil, duplicate, empty, or authority-less providers.
func NewRegistry(providers ...Provider) (*Registry, error) {
	registered := make(map[string]Provider, len(providers))
	for _, provider := range providers {
		if isNilProvider(provider) {
			return nil, fmt.Errorf("register control evidence provider: provider is nil")
		}
		id := provider.ID()
		if id == "" {
			return nil, fmt.Errorf("register control evidence provider: provider ID is empty")
		}
		if _, exists := registered[id]; exists {
			return nil, fmt.Errorf("register control evidence provider %s: duplicate ID", id)
		}
		if !validAuthority(provider.Authority()) {
			return nil, fmt.Errorf("register control evidence provider %s: unsupported authority", id)
		}
		registered[id] = provider
	}
	return &Registry{providers: registered}, nil
}

func isNilProvider(provider Provider) bool {
	if provider == nil {
		return true
	}
	value := reflect.ValueOf(provider)
	return value.Kind() == reflect.Ptr && value.IsNil()
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

// Provider returns the exact registered collector provider, if present.
func (registry *Registry) Provider(collectorID string) (Provider, bool) {
	if registry == nil {
		return nil, false
	}
	provider, ok := registry.providers[collectorID]
	return provider, ok
}

// IDs returns registered collector IDs in stable order.
func (registry *Registry) IDs() []string {
	if registry == nil {
		return []string{}
	}
	ids := make([]string, 0, len(registry.providers))
	for id := range registry.providers {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}

// BindingResolver supplies scope and policy inputs that were established
// independently of the evidence provider.
type BindingResolver interface {
	Resolve(controlprogramcatalog.Template) (controlprogramcatalog.RuntimeBinding, bool)
}

// ResolveFunc adapts a function into a BindingResolver.
type ResolveFunc func(controlprogramcatalog.Template) (controlprogramcatalog.RuntimeBinding, bool)

func (resolve ResolveFunc) Resolve(template controlprogramcatalog.Template) (controlprogramcatalog.RuntimeBinding, bool) {
	return resolve(template)
}

// Execution is the bounded, evidence-linked result for one exact clause.
type Execution struct {
	TemplateID                   string                    `json:"template_id"`
	CollectorID                  string                    `json:"collector_id"`
	ControlID                    string                    `json:"control_id"`
	ControlRevision              int                       `json:"control_revision"`
	ClauseID                     string                    `json:"clause_id"`
	ClauseOrdinal                int                       `json:"clause_ordinal"`
	ImplementationID             string                    `json:"implementation_id"`
	ImplementationContractSHA256 string                    `json:"implementation_contract_sha256"`
	RequiredAuthority            controlprogram.Authority  `json:"required_authority"`
	ProviderID                   string                    `json:"provider_id,omitempty"`
	ProgramSHA256                string                    `json:"program_sha256,omitempty"`
	EvidenceSHA256               string                    `json:"evidence_sha256,omitempty"`
	Status                       ExecutionStatus           `json:"status"`
	Outcome                      controlprogram.Outcome    `json:"outcome,omitempty"`
	ReasonCode                   controlprogram.ReasonCode `json:"reason_code,omitempty"`
	EvaluatedAt                  time.Time                 `json:"evaluated_at"`
	validated                    bool
	evidence                     *controlprogram.Evidence
}

func baseExecution(template controlprogramcatalog.Template, now time.Time) Execution {
	return Execution{
		TemplateID: template.TemplateID, CollectorID: template.CollectorContract.CollectorID,
		ControlID: template.ControlID, ControlRevision: template.ControlRevision,
		ClauseID: template.ClauseID, ClauseOrdinal: template.ClauseOrdinal,
		ImplementationID:             template.ImplementationID,
		ImplementationContractSHA256: template.ImplementationContractSHA256,
		RequiredAuthority:            template.RequiredAuthority, EvaluatedAt: now.UTC(),
		validated: true,
	}
}

// Authenticated reports whether this value was constructed by this runtime.
// Serialized provider output cannot manufacture an attachable execution.
func (execution Execution) Authenticated() bool { return execution.validated }

// SealedEvidence returns a defensive copy of structurally valid evidence
// collected by this runtime. Authoritative outcomes can therefore be replayed
// from the run record instead of trusting only a result digest.
func (execution Execution) SealedEvidence() (controlprogram.Evidence, bool) {
	if !execution.validated || execution.evidence == nil {
		return controlprogram.Evidence{}, false
	}
	data, err := json.Marshal(execution.evidence)
	if err != nil {
		return controlprogram.Evidence{}, false
	}
	copy, err := controlprogram.DecodeEvidence(data)
	if err != nil {
		return controlprogram.Evidence{}, false
	}
	return copy, true
}

// Evaluate executes one exact template. All inability to establish a complete,
// current, correctly bound fact set remains Blocked.
func Evaluate(ctx context.Context, template controlprogramcatalog.Template, binding controlprogramcatalog.RuntimeBinding, registry *Registry, now time.Time) Execution {
	execution := baseExecution(template, now)
	if now.IsZero() {
		execution.Status = StatusBlockedBinding
		return execution
	}
	program, err := template.Program(binding)
	if err != nil {
		execution.Status = StatusBlockedBinding
		return execution
	}
	execution.ProgramSHA256 = controlprogram.ProgramSHA256(program)

	provider, ok := registry.Provider(template.CollectorContract.CollectorID)
	if !ok {
		execution.Status = StatusBlockedProvider
		return execution
	}
	execution.ProviderID = provider.ID()
	if provider.Authority() != template.RequiredAuthority {
		execution.Status = StatusBlockedAuthority
		return execution
	}
	if err := ctx.Err(); err != nil {
		execution.Status = StatusBlockedCanceled
		return execution
	}
	evidence, err := provider.Collect(ctx, Request{Template: template, Program: program, EvaluationTime: now.UTC()})
	if err != nil {
		if ctx.Err() != nil {
			execution.Status = StatusBlockedCanceled
		} else {
			execution.Status = StatusBlockedCollection
		}
		return execution
	}
	if controlprogram.ValidateEvidence(evidence) == nil {
		copy := evidence
		execution.evidence = &copy
	}
	result := controlprogram.Evaluate(program, evidence, now)
	execution.EvidenceSHA256 = result.EvidenceSHA256
	execution.Outcome = result.Outcome
	execution.ReasonCode = result.ReasonCode
	switch result.Outcome {
	case controlprogram.OutcomePass:
		execution.Status = StatusPassed
	case controlprogram.OutcomeFail:
		execution.Status = StatusFailed
	case controlprogram.OutcomeNotApplicable:
		execution.Status = StatusNotApplicable
	default:
		execution.Status = StatusBlockedEvidence
	}
	return execution
}

// NewApplicableEvidence creates the scanner-owned envelope around provider
// facts. Identity, authority, scope, inventory, and observation time always
// come from the already-bound request rather than provider-selected values.
func NewApplicableEvidence(request Request, evidenceID string, facts map[string]controlprogram.Fact, complete bool) (controlprogram.Evidence, error) {
	if request.EvaluationTime.IsZero() {
		return controlprogram.Evidence{}, fmt.Errorf("construct applicable evidence: evaluation time is missing")
	}
	evidence := controlprogram.Evidence{
		SchemaVersion: controlprogram.EvidenceSchemaVersion, EvidenceID: evidenceID,
		ProgramSHA256: controlprogram.ProgramSHA256(request.Program), ControlID: request.Program.ControlID,
		ControlRevision: request.Program.ControlRevision, ControlSemanticSHA256: request.Program.ControlSemanticSHA256,
		ClauseID: request.Program.ClauseID, ClauseSHA256: request.Program.ClauseSHA256,
		ImplementationContractSHA256: request.Program.ImplementationContractSHA256,
		SubjectID:                    request.Program.SubjectID, ObservedSubjects: append([]string(nil), request.Program.Subjects...),
		InventorySHA256: request.Program.InventorySHA256, Authority: request.Program.RequiredAuthority,
		ObservedAt: request.EvaluationTime.UTC(), Complete: complete,
		Applicability:                    controlprogram.ApplicabilityApplicable,
		ApplicabilityProofContractSHA256: request.Program.ApplicabilityProofContractSHA256,
		Facts:                            facts,
	}
	if err := controlprogram.ValidateEvidence(evidence); err != nil {
		return controlprogram.Evidence{}, fmt.Errorf("construct applicable evidence: %w", err)
	}
	return evidence, nil
}

// EvaluateCatalog evaluates every exact template in stable catalog order. A
// missing binding is recorded per template rather than aborting unrelated work.
func EvaluateCatalog(ctx context.Context, catalog *controlprogramcatalog.Catalog, resolver BindingResolver, registry *Registry, now time.Time) []Execution {
	if catalog == nil {
		return []Execution{}
	}
	templates := catalog.Templates()
	results := make([]Execution, 0, len(templates))
	for _, template := range templates {
		if ctx.Err() != nil {
			execution := baseExecution(template, now)
			execution.Status = StatusBlockedCanceled
			results = append(results, execution)
			continue
		}
		if resolver == nil {
			execution := baseExecution(template, now)
			execution.Status = StatusBlockedBinding
			results = append(results, execution)
			continue
		}
		binding, ok := resolver.Resolve(template)
		if !ok {
			execution := baseExecution(template, now)
			execution.Status = StatusBlockedBinding
			results = append(results, execution)
			continue
		}
		results = append(results, Evaluate(ctx, template, binding, registry, now))
	}
	return results
}
