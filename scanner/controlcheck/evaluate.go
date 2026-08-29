package controlcheck

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"reflect"
	"sort"
	"strings"
)

// verifiedValue can only be constructed after this package has checked every
// binding. Pass, fail, and not-applicable results are emitted only from one of
// these values; raw provider or verifier output is never aggregated directly.
type verifiedValue struct {
	clauseID     string
	clauseSHA256 string
	outcome      Outcome
	reason       ReasonCode
}

type clauseEvaluation struct {
	verified *verifiedValue
	result   ClauseResult
}

// ClauseSHA256 returns the digest used to bind an exact clause statement.
func ClauseSHA256(statement string) string {
	sum := sha256.Sum256([]byte(statement))
	return hex.EncodeToString(sum[:])
}

// Evaluate applies the immutable implementation registry, obtains evidence,
// validates it, invokes the registered verifier, and aggregates every clause.
// It never mutates the supplied binding.
func Evaluate(ctx context.Context, request Request) Result {
	binding, ok := normalizedBinding(request.Binding)
	if !ok || request.Now.IsZero() {
		return aggregateInvalid(request.Binding)
	}

	evaluations := make([]clauseEvaluation, 0, len(binding.Clauses))
	for _, clause := range binding.Clauses {
		evaluations = append(evaluations, evaluateClause(ctx, request, binding, clause))
	}
	return aggregate(binding, evaluations)
}

func evaluateClause(ctx context.Context, request Request, binding Binding, clause Clause) clauseEvaluation {
	descriptor, exists := LookupImplementation(clause.ImplementationID)
	if !exists {
		return unverified(clause, OutcomeBlocked, ReasonImplementationUnavailable)
	}
	if descriptor.Family != clause.Family || descriptor.ImplementationDigest != clause.ImplementationDigest {
		return unverified(clause, OutcomeBlocked, ReasonImplementationDigestMismatch)
	}
	if !supportsTarget(descriptor, clause.Target) {
		return unverified(clause, OutcomeUnknown, ReasonUnsupportedTarget)
	}
	if !supports(descriptor, clause.Target, clause.RequiredAuthority) {
		return unverified(clause, OutcomeBlocked, ReasonWrongAuthority)
	}

	if isNilInterface(request.Providers) {
		return unverified(clause, OutcomeBlocked, ReasonProviderMissing)
	}
	provider, exists := request.Providers.ProviderFor(ctx, cloneClause(clause))
	if !exists || isNilInterface(provider) {
		return unverified(clause, OutcomeBlocked, ReasonProviderMissing)
	}
	providerDescriptor := provider.Descriptor()
	if providerDescriptor.ID == "" || !isSHA256(providerDescriptor.Digest) {
		return unverified(clause, OutcomeBlocked, ReasonProviderUnavailable)
	}
	if providerDescriptor.Authority != clause.RequiredAuthority {
		return unverified(clause, OutcomeBlocked, ReasonWrongAuthority)
	}

	if isNilInterface(request.Verifiers) {
		return unverified(clause, OutcomeBlocked, ReasonVerifierMissing)
	}
	verifier, exists := request.Verifiers.VerifierFor(ctx, cloneClause(clause))
	if !exists || isNilInterface(verifier) {
		return unverified(clause, OutcomeBlocked, ReasonVerifierMissing)
	}
	if !descriptorEqual(verifier.Descriptor(), descriptor) {
		return unverified(clause, OutcomeBlocked, ReasonImplementationDigestMismatch)
	}

	evidenceRequest := EvidenceRequest{
		ControlID:             binding.ControlID,
		ControlRevision:       binding.ControlRevision,
		ControlSemanticSHA256: binding.ControlSemanticSHA256,
		ClauseID:              clause.ID,
		ClauseSHA256:          clause.SHA256,
		SubjectID:             binding.SubjectID,
		ExpectedSubjects:      append([]string(nil), binding.Subjects...),
		InventoryDigest:       binding.InventoryDigest,
		Target:                clause.Target,
	}
	evidence, err := provider.Collect(ctx, evidenceRequest)
	if err != nil {
		return unverified(clause, OutcomeBlocked, ReasonProviderUnavailable)
	}
	if evidence.EvidenceID == "" || evidence.Payload == nil {
		return unverified(clause, OutcomeBlocked, ReasonEvidenceMissing)
	}
	if !evidenceMatches(evidence, evidenceRequest, providerDescriptor) {
		return unverified(clause, OutcomeBlocked, ReasonEvidenceBindingMismatch)
	}
	computedEvidenceDigest := digestBytes(evidence.Payload)
	if !isSHA256(evidence.PayloadSHA256) || evidence.PayloadSHA256 != computedEvidenceDigest {
		return unverified(clause, OutcomeBlocked, ReasonEvidenceDigestMismatch)
	}
	if evidence.ObservedAt.IsZero() {
		return unverified(clause, OutcomeBlocked, ReasonEvidenceMissing)
	}
	if evidence.ObservedAt.After(request.Now) {
		return unverified(clause, OutcomeBlocked, ReasonFutureEvidence)
	}
	if request.Now.Sub(evidence.ObservedAt) > clause.MaximumEvidenceAge {
		return unverified(clause, OutcomeBlocked, ReasonStaleEvidence)
	}
	if clause.RequireCompleteEvidence && !sameSubjectSet(binding.Subjects, evidence.ObservedSubjects) {
		return unverified(clause, OutcomeBlocked, ReasonIncompleteEvidence)
	}
	if len(evidence.ContradictionDigests) != 0 {
		return unverified(clause, OutcomeBlocked, ReasonContradictoryEvidence)
	}

	verification, err := verifier.Verify(ctx, VerificationRequest{
		Descriptor: cloneDescriptor(descriptor),
		Binding:    cloneBinding(binding),
		Clause:     cloneClause(clause),
		Evidence:   cloneEvidence(evidence),
	})
	if err != nil {
		return unverified(clause, OutcomeBlocked, ReasonVerifierUnavailable)
	}
	if !verificationMatches(verification, binding, clause, descriptor, computedEvidenceDigest) {
		return unverified(clause, OutcomeBlocked, ReasonVerificationBindingMismatch)
	}
	if len(verification.ContradictionDigests) != 0 || contradictoryDecisions(verification.Decisions) {
		return unverified(clause, OutcomeBlocked, ReasonContradictoryEvidence)
	}

	return sealVerification(clause, verification)
}

func sealVerification(clause Clause, verification Verification) clauseEvaluation {
	switch verification.Applicability {
	case ApplicabilityUnknown:
		if len(verification.Decisions) != 0 {
			return unverified(clause, OutcomeBlocked, ReasonVerificationIncomplete)
		}
		return unverified(clause, OutcomeUnknown, ReasonApplicabilityUnknown)
	case ApplicabilityNotApplicable:
		if !clause.AllowNotApplicable {
			return unverified(clause, OutcomeBlocked, ReasonNotApplicableDisallowed)
		}
		if len(verification.Decisions) != 0 {
			return unverified(clause, OutcomeBlocked, ReasonVerificationIncomplete)
		}
		if verification.ApplicabilityProof == nil || !isSHA256(verification.ApplicabilityProofSHA256) ||
			digestBytes(verification.ApplicabilityProof) != verification.ApplicabilityProofSHA256 {
			return unverified(clause, OutcomeBlocked, ReasonApplicabilityProofMissing)
		}
		return verified(clause, OutcomeNotApplicable, ReasonNotApplicable)
	case ApplicabilityApplicable:
		if len(verification.Decisions) != 1 {
			return unverified(clause, OutcomeBlocked, ReasonVerificationIncomplete)
		}
		switch verification.Decisions[0] {
		case OutcomePass:
			return verified(clause, OutcomePass, ReasonPassed)
		case OutcomeFail:
			return verified(clause, OutcomeFail, ReasonFailed)
		default:
			return unverified(clause, OutcomeBlocked, ReasonVerificationIncomplete)
		}
	default:
		return unverified(clause, OutcomeUnknown, ReasonApplicabilityUnknown)
	}
}

func verified(clause Clause, outcome Outcome, reason ReasonCode) clauseEvaluation {
	value := &verifiedValue{clauseID: clause.ID, clauseSHA256: clause.SHA256, outcome: outcome, reason: reason}
	return clauseEvaluation{verified: value, result: ClauseResult{
		ClauseID: value.clauseID, ClauseSHA256: value.clauseSHA256, Outcome: value.outcome, ReasonCode: value.reason,
	}}
}

func unverified(clause Clause, outcome Outcome, reason ReasonCode) clauseEvaluation {
	return clauseEvaluation{result: ClauseResult{
		ClauseID: clause.ID, ClauseSHA256: clause.SHA256, Outcome: outcome, ReasonCode: reason,
	}}
}

func evidenceMatches(evidence Evidence, request EvidenceRequest, provider ProviderDescriptor) bool {
	return evidence.ProviderID == provider.ID && evidence.ProviderDigest == provider.Digest &&
		evidence.ControlID == request.ControlID && evidence.ControlRevision == request.ControlRevision &&
		evidence.ControlSemanticSHA256 == request.ControlSemanticSHA256 && evidence.ClauseID == request.ClauseID &&
		evidence.ClauseSHA256 == request.ClauseSHA256 && evidence.SubjectID == request.SubjectID &&
		evidence.InventoryDigest == request.InventoryDigest && evidence.Target == request.Target
}

func verificationMatches(verification Verification, binding Binding, clause Clause, descriptor Descriptor, evidenceDigest string) bool {
	return verification.ImplementationID == descriptor.ImplementationID &&
		verification.ImplementationDigest == descriptor.ImplementationDigest &&
		verification.ControlID == binding.ControlID && verification.ControlRevision == binding.ControlRevision &&
		verification.ControlSemanticSHA256 == binding.ControlSemanticSHA256 && verification.ClauseID == clause.ID &&
		verification.ClauseSHA256 == clause.SHA256 && verification.SubjectID == binding.SubjectID &&
		verification.InventoryDigest == binding.InventoryDigest && verification.Target == clause.Target &&
		verification.EvidenceSHA256 == evidenceDigest
}

func contradictoryDecisions(decisions []Outcome) bool {
	hasPass, hasFail := false, false
	for _, decision := range decisions {
		hasPass = hasPass || decision == OutcomePass
		hasFail = hasFail || decision == OutcomeFail
	}
	return hasPass && hasFail
}

func aggregate(binding Binding, evaluations []clauseEvaluation) Result {
	result := Result{
		ControlID: binding.ControlID, ControlRevision: binding.ControlRevision,
		ControlSemanticSHA256: binding.ControlSemanticSHA256,
		Clauses:               make([]ClauseResult, len(evaluations)),
	}
	for i, evaluation := range evaluations {
		result.Clauses[i] = sealedResult(evaluation)
	}
	sort.Slice(result.Clauses, func(i, j int) bool {
		if result.Clauses[i].ClauseID == result.Clauses[j].ClauseID {
			return result.Clauses[i].ClauseSHA256 < result.Clauses[j].ClauseSHA256
		}
		return result.Clauses[i].ClauseID < result.Clauses[j].ClauseID
	})

	for _, outcome := range []Outcome{OutcomeFail, OutcomeBlocked, OutcomeUnknown, OutcomePass, OutcomeNotApplicable} {
		for _, clause := range result.Clauses {
			if clause.Outcome == outcome {
				result.Outcome = outcome
				switch outcome {
				case OutcomeFail:
					result.ReasonCode = ReasonFailed
				case OutcomePass:
					result.ReasonCode = ReasonPassed
				case OutcomeNotApplicable:
					result.ReasonCode = ReasonNotApplicable
				default:
					result.ReasonCode = clause.ReasonCode
				}
				return result
			}
		}
	}
	result.Outcome, result.ReasonCode = OutcomeBlocked, ReasonInvalidBinding
	return result
}

func aggregateInvalid(binding Binding) Result {
	clauses := append([]Clause(nil), binding.Clauses...)
	sort.Slice(clauses, func(i, j int) bool { return clauses[i].ID < clauses[j].ID })
	result := Result{
		ControlID: binding.ControlID, ControlRevision: binding.ControlRevision,
		ControlSemanticSHA256: binding.ControlSemanticSHA256,
		Outcome:               OutcomeBlocked, ReasonCode: ReasonInvalidBinding,
		Clauses: make([]ClauseResult, len(clauses)),
	}
	for i, clause := range clauses {
		result.Clauses[i] = ClauseResult{ClauseID: clause.ID, ClauseSHA256: clause.SHA256, Outcome: OutcomeBlocked, ReasonCode: ReasonInvalidBinding}
	}
	return result
}

func normalizedBinding(input Binding) (Binding, bool) {
	binding := cloneBinding(input)
	if strings.TrimSpace(binding.ControlID) == "" || binding.ControlRevision <= 0 ||
		!isSHA256(binding.ControlSemanticSHA256) || strings.TrimSpace(binding.SubjectID) == "" ||
		!isSHA256(binding.InventoryDigest) || binding.Aggregation != AggregationAllClausesPass ||
		len(binding.Subjects) == 0 || len(binding.Clauses) == 0 {
		return Binding{}, false
	}
	if !validUniqueStrings(binding.Subjects) {
		return Binding{}, false
	}
	sort.Strings(binding.Subjects)
	seenClauses := make(map[string]struct{}, len(binding.Clauses))
	for _, clause := range binding.Clauses {
		if !validClause(clause) {
			return Binding{}, false
		}
		if _, duplicate := seenClauses[clause.ID]; duplicate {
			return Binding{}, false
		}
		seenClauses[clause.ID] = struct{}{}
	}
	sort.Slice(binding.Clauses, func(i, j int) bool { return binding.Clauses[i].ID < binding.Clauses[j].ID })
	return binding, true
}

func validClause(clause Clause) bool {
	return isSHA256(clause.ID) && strings.TrimSpace(clause.Statement) != "" &&
		isSHA256(clause.SHA256) && clause.SHA256 == ClauseSHA256(clause.Statement) &&
		knownFamily(clause.Family) && knownTarget(clause.Target) && knownAuthority(clause.RequiredAuthority) &&
		clause.MaximumEvidenceAge > 0 && clause.RequireCompleteEvidence && strings.TrimSpace(clause.ImplementationID) != "" &&
		isSHA256(clause.ImplementationDigest)
}

func sealedResult(evaluation clauseEvaluation) ClauseResult {
	switch evaluation.result.Outcome {
	case OutcomePass, OutcomeFail, OutcomeNotApplicable:
		value := evaluation.verified
		if value == nil || value.clauseID != evaluation.result.ClauseID || value.clauseSHA256 != evaluation.result.ClauseSHA256 ||
			value.outcome != evaluation.result.Outcome || value.reason != evaluation.result.ReasonCode {
			return ClauseResult{
				ClauseID: evaluation.result.ClauseID, ClauseSHA256: evaluation.result.ClauseSHA256,
				Outcome: OutcomeBlocked, ReasonCode: ReasonVerificationIncomplete,
			}
		}
	}
	return evaluation.result
}

func validUniqueStrings(values []string) bool {
	seen := make(map[string]struct{}, len(values))
	for _, value := range values {
		if value == "" || strings.TrimSpace(value) != value {
			return false
		}
		if _, duplicate := seen[value]; duplicate {
			return false
		}
		seen[value] = struct{}{}
	}
	return true
}

func sameSubjectSet(expected, observed []string) bool {
	if !validUniqueStrings(observed) || len(expected) != len(observed) {
		return false
	}
	copyObserved := append([]string(nil), observed...)
	sort.Strings(copyObserved)
	for i := range expected {
		if expected[i] != copyObserved[i] {
			return false
		}
	}
	return true
}

func cloneBinding(binding Binding) Binding {
	binding.Subjects = append([]string(nil), binding.Subjects...)
	binding.Clauses = append([]Clause(nil), binding.Clauses...)
	return binding
}

func cloneClause(clause Clause) Clause { return clause }

func cloneEvidence(evidence Evidence) Evidence {
	evidence.Payload = append([]byte(nil), evidence.Payload...)
	evidence.ObservedSubjects = append([]string(nil), evidence.ObservedSubjects...)
	evidence.ContradictionDigests = append([]string(nil), evidence.ContradictionDigests...)
	return evidence
}

func supportsTarget(descriptor Descriptor, target Target) bool {
	for _, capability := range descriptor.Capabilities {
		if capability.Target == target {
			return true
		}
	}
	return false
}

func knownFamily(value Family) bool {
	for _, descriptor := range implementationRegistry {
		if descriptor.Family == value {
			return true
		}
	}
	return false
}

func knownTarget(value Target) bool {
	switch value {
	case TargetRepository, TargetArtifact, TargetExecution, TargetEnvironment, TargetExternalRegistry, TargetStructuredRecord:
		return true
	default:
		return false
	}
}

func knownAuthority(value Authority) bool {
	switch value {
	case AuthorityRepository, AuthorityArtifact, AuthorityExecuted, AuthorityEnvironment, AuthorityExternalRegistry, AuthorityStructuredRecord:
		return true
	default:
		return false
	}
}

func digestBytes(value []byte) string {
	sum := sha256.Sum256(value)
	return hex.EncodeToString(sum[:])
}

func isSHA256(value string) bool {
	if len(value) != sha256.Size*2 || strings.ToLower(value) != value {
		return false
	}
	decoded, err := hex.DecodeString(value)
	return err == nil && len(decoded) == sha256.Size
}

func isNilInterface(value any) bool {
	if value == nil {
		return true
	}
	reflected := reflect.ValueOf(value)
	switch reflected.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return reflected.IsNil()
	default:
		return false
	}
}
