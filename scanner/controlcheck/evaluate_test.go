package controlcheck

import (
	"context"
	"errors"
	"reflect"
	"regexp"
	"sort"
	"testing"
	"time"
)

type fakeProvider struct {
	descriptor ProviderDescriptor
	now        time.Time
	mutate     func(*Evidence)
	err        error
	calls      int
}

func (provider *fakeProvider) Descriptor() ProviderDescriptor { return provider.descriptor }

func (provider *fakeProvider) Collect(_ context.Context, request EvidenceRequest) (Evidence, error) {
	provider.calls++
	if provider.err != nil {
		return Evidence{}, provider.err
	}
	payload := []byte("bounded evidence for " + request.ClauseID)
	evidence := Evidence{
		EvidenceID:            "evidence-" + request.ClauseID,
		ProviderID:            provider.descriptor.ID,
		ProviderDigest:        provider.descriptor.Digest,
		ControlID:             request.ControlID,
		ControlRevision:       request.ControlRevision,
		ControlSemanticSHA256: request.ControlSemanticSHA256,
		ClauseID:              request.ClauseID,
		ClauseSHA256:          request.ClauseSHA256,
		SubjectID:             request.SubjectID,
		InventoryDigest:       request.InventoryDigest,
		Target:                request.Target,
		ObservedAt:            provider.now.Add(-time.Minute),
		PayloadSHA256:         digestBytes(payload),
		Payload:               payload,
		ObservedSubjects:      append([]string(nil), request.ExpectedSubjects...),
	}
	if provider.mutate != nil {
		provider.mutate(&evidence)
	}
	return evidence, nil
}

type fakeVerifier struct {
	descriptor Descriptor
	mutate     func(*Verification)
	err        error
	calls      int
}

func (verifier *fakeVerifier) Descriptor() Descriptor { return cloneDescriptor(verifier.descriptor) }

func (verifier *fakeVerifier) Verify(_ context.Context, request VerificationRequest) (Verification, error) {
	verifier.calls++
	if verifier.err != nil {
		return Verification{}, verifier.err
	}
	verification := Verification{
		ImplementationID:      request.Descriptor.ImplementationID,
		ImplementationDigest:  request.Descriptor.ImplementationDigest,
		ControlID:             request.Binding.ControlID,
		ControlRevision:       request.Binding.ControlRevision,
		ControlSemanticSHA256: request.Binding.ControlSemanticSHA256,
		ClauseID:              request.Clause.ID,
		ClauseSHA256:          request.Clause.SHA256,
		SubjectID:             request.Binding.SubjectID,
		InventoryDigest:       request.Binding.InventoryDigest,
		Target:                request.Clause.Target,
		EvidenceSHA256:        request.Evidence.PayloadSHA256,
		Applicability:         ApplicabilityApplicable,
		Decisions:             []Outcome{OutcomePass},
	}
	if verifier.mutate != nil {
		verifier.mutate(&verification)
	}
	return verification, nil
}

type providerResolver struct {
	providers map[string]EvidenceProvider
	calls     []string
}

func (resolver *providerResolver) ProviderFor(_ context.Context, clause Clause) (EvidenceProvider, bool) {
	resolver.calls = append(resolver.calls, clause.ID)
	provider, ok := resolver.providers[clause.ID]
	return provider, ok
}

type verifierResolver struct {
	verifiers map[string]Verifier
	calls     []string
}

func (resolver *verifierResolver) VerifierFor(_ context.Context, clause Clause) (Verifier, bool) {
	resolver.calls = append(resolver.calls, clause.ID)
	verifier, ok := resolver.verifiers[clause.ID]
	return verifier, ok
}

type fixture struct {
	now              time.Time
	binding          Binding
	provider         *fakeProvider
	verifier         *fakeVerifier
	providerResolver *providerResolver
	verifierResolver *verifierResolver
}

func newFixture(t *testing.T) *fixture {
	t.Helper()
	now := time.Date(2026, time.August, 29, 12, 0, 0, 0, time.UTC)
	descriptor, ok := LookupImplementation("prc.check.inventory-fact@0.1")
	if !ok {
		t.Fatal("inventory implementation missing")
	}
	statement := "Every declared workspace subject has a readable manifest."
	clause := Clause{
		ID: digestBytes([]byte("clause-one")), Statement: statement, SHA256: ClauseSHA256(statement),
		Family: FamilyInventoryFact, Target: TargetRepository, RequiredAuthority: AuthorityRepository,
		MaximumEvidenceAge: time.Hour, RequireCompleteEvidence: true, AllowNotApplicable: true,
		ImplementationID: descriptor.ImplementationID, ImplementationDigest: descriptor.ImplementationDigest,
	}
	binding := Binding{
		ControlID: "PRC-TEST-001", ControlRevision: 3,
		ControlSemanticSHA256: digestBytes([]byte("semantic-v3")),
		SubjectID:             "workspace", Subjects: []string{"service-b", "service-a"},
		InventoryDigest: digestBytes([]byte("inventory")), Aggregation: AggregationAllClausesPass,
		Clauses: []Clause{clause},
	}
	provider := &fakeProvider{
		descriptor: ProviderDescriptor{ID: "repository-snapshot", Digest: digestBytes([]byte("provider-v1")), Authority: AuthorityRepository},
		now:        now,
	}
	verifier := &fakeVerifier{descriptor: descriptor}
	return &fixture{
		now: now, binding: binding, provider: provider, verifier: verifier,
		providerResolver: &providerResolver{providers: map[string]EvidenceProvider{clause.ID: provider}},
		verifierResolver: &verifierResolver{verifiers: map[string]Verifier{clause.ID: verifier}},
	}
}

func (fixture *fixture) request() Request {
	return Request{
		Binding: fixture.binding, Now: fixture.now,
		Providers: fixture.providerResolver, Verifiers: fixture.verifierResolver,
	}
}

func TestEvaluateConformance(t *testing.T) {
	tests := []struct {
		name        string
		change      func(*fixture)
		outcome     Outcome
		reason      ReasonCode
		verifyCalls int
	}{
		{name: "pass", outcome: OutcomePass, reason: ReasonPassed, verifyCalls: 1},
		{name: "fail", change: func(f *fixture) {
			f.verifier.mutate = func(v *Verification) { v.Decisions = []Outcome{OutcomeFail} }
		}, outcome: OutcomeFail, reason: ReasonFailed, verifyCalls: 1},
		{name: "not applicable with bounded proof", change: func(f *fixture) {
			f.verifier.mutate = func(v *Verification) {
				v.Applicability, v.Decisions = ApplicabilityNotApplicable, nil
				v.ApplicabilityProof = []byte("bounded absence proof")
				v.ApplicabilityProofSHA256 = digestBytes(v.ApplicabilityProof)
			}
		}, outcome: OutcomeNotApplicable, reason: ReasonNotApplicable, verifyCalls: 1},
		{name: "not applicable disallowed", change: func(f *fixture) {
			f.binding.Clauses[0].AllowNotApplicable = false
			f.verifier.mutate = func(v *Verification) {
				v.Applicability, v.Decisions = ApplicabilityNotApplicable, nil
				v.ApplicabilityProof = []byte("proof")
				v.ApplicabilityProofSHA256 = digestBytes(v.ApplicabilityProof)
			}
		}, outcome: OutcomeBlocked, reason: ReasonNotApplicableDisallowed, verifyCalls: 1},
		{name: "not applicable proof digest mismatch", change: func(f *fixture) {
			f.verifier.mutate = func(v *Verification) {
				v.Applicability, v.Decisions = ApplicabilityNotApplicable, nil
				v.ApplicabilityProof = []byte("proof")
				v.ApplicabilityProofSHA256 = digestBytes([]byte("different"))
			}
		}, outcome: OutcomeBlocked, reason: ReasonApplicabilityProofMissing, verifyCalls: 1},
		{name: "applicability unknown", change: func(f *fixture) {
			f.verifier.mutate = func(v *Verification) { v.Applicability, v.Decisions = ApplicabilityUnknown, nil }
		}, outcome: OutcomeUnknown, reason: ReasonApplicabilityUnknown, verifyCalls: 1},
		{name: "missing provider resolver", change: func(f *fixture) {
			f.providerResolver = nil
		}, outcome: OutcomeBlocked, reason: ReasonProviderMissing},
		{name: "missing provider registration", change: func(f *fixture) {
			f.providerResolver.providers = nil
		}, outcome: OutcomeBlocked, reason: ReasonProviderMissing},
		{name: "provider unavailable", change: func(f *fixture) {
			f.provider.err = errors.New("snapshot unavailable")
		}, outcome: OutcomeBlocked, reason: ReasonProviderUnavailable},
		{name: "missing evidence", change: func(f *fixture) {
			f.provider.mutate = func(e *Evidence) { e.EvidenceID, e.Payload = "", nil }
		}, outcome: OutcomeBlocked, reason: ReasonEvidenceMissing},
		{name: "wrong provider authority", change: func(f *fixture) {
			f.provider.descriptor.Authority = AuthorityEnvironment
		}, outcome: OutcomeBlocked, reason: ReasonWrongAuthority},
		{name: "stale evidence", change: func(f *fixture) {
			f.provider.mutate = func(e *Evidence) { e.ObservedAt = f.now.Add(-2 * time.Hour) }
		}, outcome: OutcomeBlocked, reason: ReasonStaleEvidence},
		{name: "future evidence", change: func(f *fixture) {
			f.provider.mutate = func(e *Evidence) { e.ObservedAt = f.now.Add(time.Second) }
		}, outcome: OutcomeBlocked, reason: ReasonFutureEvidence},
		{name: "incomplete inventory", change: func(f *fixture) {
			f.provider.mutate = func(e *Evidence) { e.ObservedSubjects = []string{"service-a"} }
		}, outcome: OutcomeBlocked, reason: ReasonIncompleteEvidence},
		{name: "extra inventory subject", change: func(f *fixture) {
			f.provider.mutate = func(e *Evidence) { e.ObservedSubjects = append(e.ObservedSubjects, "undeclared") }
		}, outcome: OutcomeBlocked, reason: ReasonIncompleteEvidence},
		{name: "payload digest mismatch", change: func(f *fixture) {
			f.provider.mutate = func(e *Evidence) { e.PayloadSHA256 = digestBytes([]byte("other")) }
		}, outcome: OutcomeBlocked, reason: ReasonEvidenceDigestMismatch},
		{name: "evidence control binding mismatch", change: func(f *fixture) {
			f.provider.mutate = func(e *Evidence) { e.ControlRevision++ }
		}, outcome: OutcomeBlocked, reason: ReasonEvidenceBindingMismatch},
		{name: "evidence contradiction", change: func(f *fixture) {
			f.provider.mutate = func(e *Evidence) { e.ContradictionDigests = []string{digestBytes([]byte("conflict"))} }
		}, outcome: OutcomeBlocked, reason: ReasonContradictoryEvidence},
		{name: "missing verifier resolver", change: func(f *fixture) {
			f.verifierResolver = nil
		}, outcome: OutcomeBlocked, reason: ReasonVerifierMissing},
		{name: "missing verifier registration", change: func(f *fixture) {
			f.verifierResolver.verifiers = nil
		}, outcome: OutcomeBlocked, reason: ReasonVerifierMissing},
		{name: "verifier unavailable", change: func(f *fixture) {
			f.verifier.err = errors.New("implementation unavailable")
		}, outcome: OutcomeBlocked, reason: ReasonVerifierUnavailable, verifyCalls: 1},
		{name: "verifier descriptor digest mismatch", change: func(f *fixture) {
			f.verifier.descriptor.ImplementationDigest = digestBytes([]byte("replacement"))
		}, outcome: OutcomeBlocked, reason: ReasonImplementationDigestMismatch},
		{name: "verification evidence digest mismatch", change: func(f *fixture) {
			f.verifier.mutate = func(v *Verification) { v.EvidenceSHA256 = digestBytes([]byte("other evidence")) }
		}, outcome: OutcomeBlocked, reason: ReasonVerificationBindingMismatch, verifyCalls: 1},
		{name: "verification contradiction", change: func(f *fixture) {
			f.verifier.mutate = func(v *Verification) { v.Decisions = []Outcome{OutcomePass, OutcomeFail} }
		}, outcome: OutcomeBlocked, reason: ReasonContradictoryEvidence, verifyCalls: 1},
		{name: "unsupported target", change: func(f *fixture) {
			f.binding.Clauses[0].Target = TargetArtifact
		}, outcome: OutcomeUnknown, reason: ReasonUnsupportedTarget},
		{name: "wrong bound authority", change: func(f *fixture) {
			f.binding.Clauses[0].RequiredAuthority = AuthorityEnvironment
		}, outcome: OutcomeBlocked, reason: ReasonWrongAuthority},
		{name: "implementation unavailable", change: func(f *fixture) {
			f.binding.Clauses[0].ImplementationID = "prc.check.unregistered@0.1"
		}, outcome: OutcomeBlocked, reason: ReasonImplementationUnavailable},
		{name: "bound implementation digest mismatch", change: func(f *fixture) {
			f.binding.Clauses[0].ImplementationDigest = digestBytes([]byte("wrong contract"))
		}, outcome: OutcomeBlocked, reason: ReasonImplementationDigestMismatch},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fixture := newFixture(t)
			if test.change != nil {
				test.change(fixture)
			}
			result := Evaluate(context.Background(), fixture.request())
			if result.Outcome != test.outcome || result.ReasonCode != test.reason {
				t.Fatalf("got %s/%s, want %s/%s: %#v", result.Outcome, result.ReasonCode, test.outcome, test.reason, result)
			}
			if fixture.verifier.calls != test.verifyCalls {
				t.Fatalf("verifier calls = %d, want %d", fixture.verifier.calls, test.verifyCalls)
			}
		})
	}
}

func TestPassWhileBrokenIsBlockedBeforeVerification(t *testing.T) {
	fixture := newFixture(t)
	fixture.provider.mutate = func(evidence *Evidence) {
		evidence.ObservedSubjects = []string{"service-a"}
	}
	fixture.verifier.mutate = func(verification *Verification) {
		verification.Decisions = []Outcome{OutcomePass}
	}

	result := Evaluate(context.Background(), fixture.request())
	if result.Outcome != OutcomeBlocked || result.ReasonCode != ReasonIncompleteEvidence {
		t.Fatalf("partial inventory must not pass: %#v", result)
	}
	if fixture.verifier.calls != 0 {
		t.Fatalf("verifier was called %d times with incomplete evidence", fixture.verifier.calls)
	}
}

func TestMixedFamilyClausesAreResolvedIndependently(t *testing.T) {
	fixture := newFixture(t)
	secondDescriptor, ok := LookupImplementation("prc.check.structured-record@0.1")
	if !ok {
		t.Fatal("structured-record implementation missing")
	}
	statement := "The release record binds its approval to the exact artifact digest."
	second := Clause{
		ID: digestBytes([]byte("clause-two")), Statement: statement, SHA256: ClauseSHA256(statement),
		Family: FamilyStructuredRecord, Target: TargetStructuredRecord, RequiredAuthority: AuthorityStructuredRecord,
		MaximumEvidenceAge: time.Hour, RequireCompleteEvidence: true,
		ImplementationID: secondDescriptor.ImplementationID, ImplementationDigest: secondDescriptor.ImplementationDigest,
	}
	fixture.binding.Clauses = append(fixture.binding.Clauses, second)
	secondProvider := &fakeProvider{
		descriptor: ProviderDescriptor{ID: "release-record", Digest: digestBytes([]byte("record-provider")), Authority: AuthorityStructuredRecord},
		now:        fixture.now,
	}
	secondVerifier := &fakeVerifier{descriptor: secondDescriptor}
	fixture.providerResolver.providers[second.ID] = secondProvider
	fixture.verifierResolver.verifiers[second.ID] = secondVerifier

	result := Evaluate(context.Background(), fixture.request())
	if result.Outcome != OutcomePass || len(result.Clauses) != 2 {
		t.Fatalf("mixed control did not pass both clauses: %#v", result)
	}
	if fixture.provider.calls != 1 || secondProvider.calls != 1 || fixture.verifier.calls != 1 || secondVerifier.calls != 1 {
		t.Fatalf("each family must execute once: providers %d/%d verifiers %d/%d",
			fixture.provider.calls, secondProvider.calls, fixture.verifier.calls, secondVerifier.calls)
	}
}

func TestEvaluationIsStableAndDoesNotMutateBinding(t *testing.T) {
	fixture := newFixture(t)
	original := cloneBinding(fixture.binding)
	first := Evaluate(context.Background(), fixture.request())
	if !reflect.DeepEqual(fixture.binding, original) {
		t.Fatal("evaluation mutated the supplied binding")
	}

	reversed := newFixture(t)
	reversed.binding.Subjects[0], reversed.binding.Subjects[1] = reversed.binding.Subjects[1], reversed.binding.Subjects[0]
	second := Evaluate(context.Background(), reversed.request())
	if !reflect.DeepEqual(first, second) {
		t.Fatalf("equivalent subject orders produced different results:\nfirst: %#v\nsecond: %#v", first, second)
	}
}

func TestAllClausesAggregation(t *testing.T) {
	clause := func(id string, outcome Outcome, reason ReasonCode, sealed bool) clauseEvaluation {
		result := ClauseResult{ClauseID: digestBytes([]byte(id)), ClauseSHA256: digestBytes([]byte("hash-" + id)), Outcome: outcome, ReasonCode: reason}
		evaluation := clauseEvaluation{result: result}
		if sealed {
			evaluation.verified = &verifiedValue{clauseID: result.ClauseID, clauseSHA256: result.ClauseSHA256, outcome: outcome, reason: reason}
		}
		return evaluation
	}
	binding := Binding{ControlID: "aggregate", ControlRevision: 1, ControlSemanticSHA256: digestBytes([]byte("aggregate"))}
	tests := []struct {
		name    string
		values  []clauseEvaluation
		outcome Outcome
	}{
		{name: "pass and NA is pass", values: []clauseEvaluation{
			clause("a", OutcomePass, ReasonPassed, true), clause("b", OutcomeNotApplicable, ReasonNotApplicable, true),
		}, outcome: OutcomePass},
		{name: "fail beats blocked", values: []clauseEvaluation{
			clause("a", OutcomeBlocked, ReasonStaleEvidence, false), clause("b", OutcomeFail, ReasonFailed, true),
		}, outcome: OutcomeFail},
		{name: "blocked beats unknown", values: []clauseEvaluation{
			clause("a", OutcomeUnknown, ReasonApplicabilityUnknown, false), clause("b", OutcomeBlocked, ReasonProviderMissing, false),
		}, outcome: OutcomeBlocked},
		{name: "all NA is NA", values: []clauseEvaluation{
			clause("a", OutcomeNotApplicable, ReasonNotApplicable, true), clause("b", OutcomeNotApplicable, ReasonNotApplicable, true),
		}, outcome: OutcomeNotApplicable},
		{name: "unsealed pass is blocked", values: []clauseEvaluation{
			clause("a", OutcomePass, ReasonPassed, false),
		}, outcome: OutcomeBlocked},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := aggregate(binding, test.values)
			if result.Outcome != test.outcome {
				t.Fatalf("got %s, want %s", result.Outcome, test.outcome)
			}
		})
	}
}

func TestInvalidBindingFailsClosedWithoutCapabilities(t *testing.T) {
	tests := []struct {
		name   string
		change func(*fixture)
	}{
		{name: "semantic hash", change: func(f *fixture) { f.binding.ControlSemanticSHA256 = "not-a-digest" }},
		{name: "clause hash", change: func(f *fixture) { f.binding.Clauses[0].SHA256 = digestBytes([]byte("wrong clause")) }},
		{name: "inventory hash", change: func(f *fixture) { f.binding.InventoryDigest = "not-a-digest" }},
		{name: "completeness cannot be disabled", change: func(f *fixture) { f.binding.Clauses[0].RequireCompleteEvidence = false }},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fixture := newFixture(t)
			test.change(fixture)
			result := Evaluate(context.Background(), fixture.request())
			if result.Outcome != OutcomeBlocked || result.ReasonCode != ReasonInvalidBinding {
				t.Fatalf("invalid binding did not fail closed: %#v", result)
			}
			if len(fixture.providerResolver.calls) != 0 || len(fixture.verifierResolver.calls) != 0 {
				t.Fatal("invalid binding reached a capability resolver")
			}
		})
	}
}

func TestRawReportsCannotClaimTrustWithBooleanFields(t *testing.T) {
	for _, value := range []any{Evidence{}, Verification{}} {
		typeOf := reflect.TypeOf(value)
		for i := 0; i < typeOf.NumField(); i++ {
			field := typeOf.Field(i)
			if field.IsExported() && field.Type.Kind() == reflect.Bool {
				t.Fatalf("%s.%s is an exported trust-like boolean", typeOf.Name(), field.Name)
			}
		}
	}
}

func TestImplementationRegistryIsClosedStableAndDefensive(t *testing.T) {
	implementations := Implementations()
	if len(implementations) != 11 {
		t.Fatalf("registry has %d implementations, want 11", len(implementations))
	}
	versionedID := regexp.MustCompile(`^prc\.check\.[a-z0-9-]+@0\.1$`)
	seenFamilies := map[Family]bool{}
	seenIDs := map[string]bool{}
	expected := map[string]string{
		"prc.check.analysis-adapter@0.1":     "58874f628ccf411d678417c387f45cde2457575f8c9d6d3f504627a9b96cf8d3",
		"prc.check.artifact-integrity@0.1":   "b80869c071be1f64a4d131cc3a4738f5a45d6d65b63529c8031932d79d1ac9b5",
		"prc.check.ci-policy@0.1":            "3ba1e384d5ab180ccb2d806c13f664f0e87576d889cc821ebd96004037ea0899",
		"prc.check.container-iac@0.1":        "d32efaf4030b652d15a08b45f2a21fd2410dbcccc9e15a53608325facfaf4338",
		"prc.check.environment-evidence@0.1": "a6e267998e29bb26034c8d7d75ee887df6f806118e4f89cb3d6123a51dd137fc",
		"prc.check.execution-evidence@0.1":   "a75fe01ec7fb60727324ee533ec183f1b55cd8a66b39f12b6401959caa37dfca",
		"prc.check.inventory-fact@0.1":       "e109ca581c021b1c2439e3e2dba40e470938cace7f79d03cd1f2a28327d5dc56",
		"prc.check.package-metadata@0.1":     "d7a26047868ab897e9ae0e35a758c0ab787971d05309ad5e4a7012426b4695ff",
		"prc.check.source-ast@0.1":           "6f6b32e965f38c924b7fd9600128237c6ff4c9b1d0668da56e3553b85982ebff",
		"prc.check.structured-document@0.1":  "cd05c2c2c6cd7981722127c94b9e93e967fa72df62a99398fb0aacd77ad2d354",
		"prc.check.structured-record@0.1":    "44d9d08b14d2aa7ae1b57c511ce9c9fa16cf232fdc127cd198580be90fbd6039",
	}
	for _, descriptor := range implementations {
		if seenFamilies[descriptor.Family] || seenIDs[descriptor.ImplementationID] {
			t.Fatalf("duplicate registry descriptor: %#v", descriptor)
		}
		seenFamilies[descriptor.Family], seenIDs[descriptor.ImplementationID] = true, true
		if !versionedID.MatchString(descriptor.ImplementationID) || !isSHA256(descriptor.ImplementationDigest) || len(descriptor.Capabilities) == 0 {
			t.Fatalf("invalid descriptor: %#v", descriptor)
		}
		if expected[descriptor.ImplementationID] != descriptor.ImplementationDigest {
			t.Fatalf("unexpected implementation contract: %#v", descriptor)
		}
	}
	if !sort.SliceIsSorted(implementations, func(i, j int) bool { return implementations[i].Family < implementations[j].Family }) {
		t.Fatal("implementation registry is not in stable family order")
	}

	original := cloneDescriptor(implementations[0])
	implementations[0].Capabilities[0].Target = TargetRepository
	lookedUp, ok := LookupImplementation(original.ImplementationID)
	if !ok || !descriptorEqual(lookedUp, original) {
		t.Fatal("caller mutation changed the immutable registry")
	}
}
