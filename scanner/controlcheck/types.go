// Package controlcheck evaluates reviewed deterministic control bindings.
//
// It is intentionally independent from the catalog decoder and scan engine. A
// caller must supply a fully bound request, an evidence provider, and a
// verifier. Raw evidence and verifier reports are never trusted directly.
package controlcheck

import (
	"context"
	"time"
)

// Family identifies one of the reviewed deterministic checker families.
type Family string

const (
	FamilyInventoryFact       Family = "inventory_fact"
	FamilyStructuredDocument  Family = "structured_document"
	FamilyPackageMetadata     Family = "package_metadata"
	FamilyCIPolicy            Family = "ci_policy"
	FamilyContainerIaC        Family = "container_iac"
	FamilySourceAST           Family = "source_ast"
	FamilyArtifactIntegrity   Family = "artifact_integrity"
	FamilyAnalysisAdapter     Family = "analysis_adapter"
	FamilyExecutionEvidence   Family = "execution_evidence"
	FamilyEnvironmentEvidence Family = "environment_evidence"
	FamilyStructuredRecord    Family = "structured_record"
)

// Authority is the authoritative evidence class required by a binding.
type Authority string

const (
	AuthorityRepository       Authority = "repository"
	AuthorityArtifact         Authority = "artifact"
	AuthorityExecuted         Authority = "executed"
	AuthorityEnvironment      Authority = "environment"
	AuthorityExternalRegistry Authority = "external_registry"
	AuthorityStructuredRecord Authority = "structured_record"
)

// Target identifies the bounded subject type an implementation can inspect.
type Target string

const (
	TargetRepository       Target = "repository"
	TargetArtifact         Target = "artifact"
	TargetExecution        Target = "execution"
	TargetEnvironment      Target = "environment"
	TargetExternalRegistry Target = "external_registry"
	TargetStructuredRecord Target = "structured_record"
)

// Applicability is a verifier conclusion, not a trust signal.
type Applicability string

const (
	ApplicabilityApplicable    Applicability = "applicable"
	ApplicabilityNotApplicable Applicability = "not_applicable"
	ApplicabilityUnknown       Applicability = "unknown"
)

// Outcome is the scanner-owned result after all binding checks.
type Outcome string

const (
	OutcomePass          Outcome = "pass"
	OutcomeFail          Outcome = "fail"
	OutcomeBlocked       Outcome = "blocked"
	OutcomeUnknown       Outcome = "unknown"
	OutcomeNotApplicable Outcome = "not_applicable"
)

// Aggregation defines how atomic clauses form the control result. The first
// runtime supports only the reviewed all-clauses contract.
type Aggregation string

const (
	AggregationAllClausesPass Aggregation = "all_clauses_pass"
)

// ReasonCode is stable machine-readable result provenance.
type ReasonCode string

const (
	ReasonPassed                       ReasonCode = "passed"
	ReasonFailed                       ReasonCode = "failed"
	ReasonNotApplicable                ReasonCode = "not_applicable"
	ReasonInvalidBinding               ReasonCode = "invalid_binding"
	ReasonImplementationUnavailable    ReasonCode = "implementation_unavailable"
	ReasonImplementationDigestMismatch ReasonCode = "implementation_digest_mismatch"
	ReasonUnsupportedTarget            ReasonCode = "unsupported_target"
	ReasonProviderMissing              ReasonCode = "provider_missing"
	ReasonProviderUnavailable          ReasonCode = "provider_unavailable"
	ReasonVerifierMissing              ReasonCode = "verifier_missing"
	ReasonVerifierUnavailable          ReasonCode = "verifier_unavailable"
	ReasonEvidenceMissing              ReasonCode = "evidence_missing"
	ReasonEvidenceBindingMismatch      ReasonCode = "evidence_binding_mismatch"
	ReasonEvidenceDigestMismatch       ReasonCode = "evidence_digest_mismatch"
	ReasonWrongAuthority               ReasonCode = "wrong_authority"
	ReasonStaleEvidence                ReasonCode = "stale_evidence"
	ReasonFutureEvidence               ReasonCode = "future_evidence"
	ReasonIncompleteEvidence           ReasonCode = "incomplete_evidence"
	ReasonContradictoryEvidence        ReasonCode = "contradictory_evidence"
	ReasonApplicabilityUnknown         ReasonCode = "applicability_unknown"
	ReasonNotApplicableDisallowed      ReasonCode = "not_applicable_disallowed"
	ReasonApplicabilityProofMissing    ReasonCode = "applicability_proof_missing"
	ReasonVerificationBindingMismatch  ReasonCode = "verification_binding_mismatch"
	ReasonVerificationIncomplete       ReasonCode = "verification_incomplete"
)

// Clause is one reviewed atomic deterministic statement and its complete
// execution binding. Execution metadata is clause-local because one control
// may combine clauses from different checker families and authorities.
type Clause struct {
	ID                      string        `json:"id"`
	Statement               string        `json:"statement"`
	SHA256                  string        `json:"sha256"`
	Family                  Family        `json:"family"`
	Target                  Target        `json:"target"`
	RequiredAuthority       Authority     `json:"required_authority"`
	MaximumEvidenceAge      time.Duration `json:"maximum_evidence_age"`
	RequireCompleteEvidence bool          `json:"require_complete_evidence"`
	AllowNotApplicable      bool          `json:"allow_not_applicable"`
	ImplementationID        string        `json:"implementation_id"`
	ImplementationDigest    string        `json:"implementation_digest"`
}

// Binding is the trusted policy supplied by the catalog integration. It binds
// control identity and inventory scope. Each clause binds its own checker.
type Binding struct {
	ControlID             string      `json:"control_id"`
	ControlRevision       int         `json:"control_revision"`
	ControlSemanticSHA256 string      `json:"control_semantic_sha256"`
	SubjectID             string      `json:"subject_id"`
	Subjects              []string    `json:"subjects"`
	InventoryDigest       string      `json:"inventory_digest"`
	Aggregation           Aggregation `json:"aggregation"`
	Clauses               []Clause    `json:"clauses"`
}

// Capability is one exact target/authority pair supported by an
// implementation. Keeping the pair intact prevents independently supported
// values from being combined into an unsupported execution mode.
type Capability struct {
	Target    Target    `json:"target"`
	Authority Authority `json:"authority"`
}

// Descriptor is an immutable registry description of one family runtime.
type Descriptor struct {
	Family               Family       `json:"family"`
	ImplementationID     string       `json:"implementation_id"`
	ImplementationDigest string       `json:"implementation_digest"`
	Capabilities         []Capability `json:"capabilities"`
}

// ProviderDescriptor identifies a scanner-configured evidence source. The
// authority belongs to the configured provider boundary, not raw evidence.
type ProviderDescriptor struct {
	ID        string    `json:"id"`
	Digest    string    `json:"digest"`
	Authority Authority `json:"authority"`
}

// EvidenceRequest is the exact binding given to an evidence provider.
type EvidenceRequest struct {
	ControlID             string   `json:"control_id"`
	ControlRevision       int      `json:"control_revision"`
	ControlSemanticSHA256 string   `json:"control_semantic_sha256"`
	ClauseID              string   `json:"clause_id"`
	ClauseSHA256          string   `json:"clause_sha256"`
	SubjectID             string   `json:"subject_id"`
	ExpectedSubjects      []string `json:"expected_subjects"`
	InventoryDigest       string   `json:"inventory_digest"`
	Target                Target   `json:"target"`
}

// Evidence is untrusted provider output. It deliberately has no Trusted,
// Valid, or Complete boolean. Completeness is derived from ObservedSubjects.
type Evidence struct {
	EvidenceID            string    `json:"evidence_id"`
	ProviderID            string    `json:"provider_id"`
	ProviderDigest        string    `json:"provider_digest"`
	ControlID             string    `json:"control_id"`
	ControlRevision       int       `json:"control_revision"`
	ControlSemanticSHA256 string    `json:"control_semantic_sha256"`
	ClauseID              string    `json:"clause_id"`
	ClauseSHA256          string    `json:"clause_sha256"`
	SubjectID             string    `json:"subject_id"`
	InventoryDigest       string    `json:"inventory_digest"`
	Target                Target    `json:"target"`
	ObservedAt            time.Time `json:"observed_at"`
	PayloadSHA256         string    `json:"payload_sha256"`
	Payload               []byte    `json:"payload,omitempty"`
	ObservedSubjects      []string  `json:"observed_subjects"`
	ContradictionDigests  []string  `json:"contradiction_digests,omitempty"`
}

// EvidenceProvider obtains untrusted evidence. Registration of a provider is
// an operator capability and is outside target-repository control.
type EvidenceProvider interface {
	Descriptor() ProviderDescriptor
	Collect(context.Context, EvidenceRequest) (Evidence, error)
}

// VerificationRequest is passed to a family implementation after the package
// validates provider, binding, freshness, scope, and completeness.
type VerificationRequest struct {
	Descriptor Descriptor `json:"descriptor"`
	Binding    Binding    `json:"binding"`
	Clause     Clause     `json:"clause"`
	Evidence   Evidence   `json:"evidence"`
}

// Verification is an untrusted verifier report. Decisions is a list so the
// package can detect contradictory pass/fail claims instead of accepting the
// first one. It contains no trust boolean.
type Verification struct {
	ImplementationID         string        `json:"implementation_id"`
	ImplementationDigest     string        `json:"implementation_digest"`
	ControlID                string        `json:"control_id"`
	ControlRevision          int           `json:"control_revision"`
	ControlSemanticSHA256    string        `json:"control_semantic_sha256"`
	ClauseID                 string        `json:"clause_id"`
	ClauseSHA256             string        `json:"clause_sha256"`
	SubjectID                string        `json:"subject_id"`
	InventoryDigest          string        `json:"inventory_digest"`
	Target                   Target        `json:"target"`
	EvidenceSHA256           string        `json:"evidence_sha256"`
	Applicability            Applicability `json:"applicability"`
	ApplicabilityProofSHA256 string        `json:"applicability_proof_sha256,omitempty"`
	ApplicabilityProof       []byte        `json:"applicability_proof,omitempty"`
	Decisions                []Outcome     `json:"decisions,omitempty"`
	ContradictionDigests     []string      `json:"contradiction_digests,omitempty"`
}

// Verifier is the bounded implementation boundary. Its output remains
// untrusted until this package validates and seals it internally.
type Verifier interface {
	Descriptor() Descriptor
	Verify(context.Context, VerificationRequest) (Verification, error)
}

// ProviderResolver returns the provider configured for one exact clause. A
// resolver cannot attest trust; the returned provider and all evidence are
// still validated by this package.
type ProviderResolver interface {
	ProviderFor(context.Context, Clause) (EvidenceProvider, bool)
}

// VerifierResolver returns the registered implementation for one exact
// clause. The returned descriptor must match the immutable package registry.
type VerifierResolver interface {
	VerifierFor(context.Context, Clause) (Verifier, bool)
}

// Request contains all runtime capabilities. Now is mandatory so repeated
// evaluation of the same sealed request is deterministic.
type Request struct {
	Binding   Binding
	Now       time.Time
	Providers ProviderResolver
	Verifiers VerifierResolver
}

// ClauseResult is one stable, clause-bound result.
type ClauseResult struct {
	ClauseID     string     `json:"clause_id"`
	ClauseSHA256 string     `json:"clause_sha256"`
	Outcome      Outcome    `json:"outcome"`
	ReasonCode   ReasonCode `json:"reason_code"`
}

// Result is the all-clauses aggregate. Clauses are always sorted by ID.
type Result struct {
	ControlID             string         `json:"control_id"`
	ControlRevision       int            `json:"control_revision"`
	ControlSemanticSHA256 string         `json:"control_semantic_sha256"`
	Outcome               Outcome        `json:"outcome"`
	ReasonCode            ReasonCode     `json:"reason_code"`
	Clauses               []ClauseResult `json:"clauses"`
}
