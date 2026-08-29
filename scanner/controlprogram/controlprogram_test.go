package controlprogram

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strings"
	"testing"
	"time"
)

func stringPointer(value string) *string { return &value }
func boolPointer(value bool) *bool       { return &value }

func fixture() (Program, Evidence, time.Time) {
	now := time.Date(2026, time.August, 29, 12, 0, 0, 0, time.UTC)
	digest := digestString("release artifact")
	deadline := "2026-08-30T00:00:00Z"
	program := Program{
		SchemaVersion: ProgramSchemaVersion, ControlID: "PRC-03-003", ControlRevision: 2,
		ControlSemanticSHA256: digestString("semantic-v2"), ClauseID: digestString("clause-object"),
		ClauseSHA256:                 digestString("clause statement"),
		ImplementationContractSHA256: digestString("implementation-contract"),
		SubjectID:                    "release-2026-08-29", Subjects: []string{"artifact", "manifest"},
		InventorySHA256: digestString("artifact\x00manifest"), RequiredAuthority: AuthorityArtifact,
		AllowNotApplicable: true, ApplicabilityProofContractSHA256: digestString("na-proof-contract-v1"),
		MaximumEvidenceAgeSeconds: 3600,
		Parameters:                map[string]Parameter{},
		Predicate: Expression{Op: OpAll, Args: []Expression{
			{Op: OpIdentityEqual, Fact: "artifact.identity", String: stringPointer("release/app@sha256:abc")},
			{Op: OpSchemaEqual, Fact: "manifest.schema", String: stringPointer("spdx-2.3")},
			{Op: OpDigestEqual, Fact: "artifact.digest", String: &digest},
			{Op: OpStateIn, Fact: "signature.state", Strings: []string{"trusted", "verified"}},
			{Op: OpBooleanEqual, Fact: "signature.present", Boolean: boolPointer(true)},
			{Op: OpNumberGreaterEq, Fact: "policy.score", Number: json.Number("90.0")},
			{Op: OpTimeBeforeEq, Fact: "signature.expires_at", Timestamp: &deadline},
			{Op: OpSetContainsAll, Fact: "manifest.subjects", Strings: []string{"artifact", "manifest"}},
			{Op: OpSetDisjoint, Fact: "manifest.flags", Strings: []string{"revoked", "unsigned"}},
			{Op: OpNot, Arg: &Expression{Op: OpStringEqual, Fact: "release.channel", String: stringPointer("development")}},
		}},
	}
	evidence := Evidence{
		SchemaVersion: EvidenceSchemaVersion, EvidenceID: "artifact-proof-1",
		ControlID: program.ControlID, ControlRevision: program.ControlRevision,
		ControlSemanticSHA256: program.ControlSemanticSHA256, ClauseID: program.ClauseID,
		ClauseSHA256:                 program.ClauseSHA256,
		ImplementationContractSHA256: program.ImplementationContractSHA256,
		SubjectID:                    program.SubjectID, ObservedSubjects: append([]string(nil), program.Subjects...),
		InventorySHA256: program.InventorySHA256, Authority: program.RequiredAuthority,
		ObservedAt: now.Add(-time.Minute), Complete: true, Applicability: ApplicabilityApplicable,
		ApplicabilityProofContractSHA256: program.ApplicabilityProofContractSHA256,
		Facts: map[string]Fact{
			"artifact.identity":    {Type: FactIdentity, Complete: true, String: stringPointer("release/app@sha256:abc")},
			"manifest.schema":      {Type: FactSchema, Complete: true, String: stringPointer("spdx-2.3")},
			"artifact.digest":      {Type: FactDigest, Complete: true, String: &digest},
			"signature.state":      {Type: FactState, Complete: true, String: stringPointer("verified")},
			"signature.present":    {Type: FactBoolean, Complete: true, Boolean: boolPointer(true)},
			"policy.score":         {Type: FactNumber, Complete: true, Number: json.Number("90")},
			"signature.expires_at": {Type: FactTime, Complete: true, Timestamp: stringPointer("2026-08-29T23:00:00Z")},
			"manifest.subjects":    {Type: FactStringSet, Complete: true, Strings: []string{"artifact", "manifest", "signature"}},
			"manifest.flags":       {Type: FactStringSet, Complete: true, Strings: []string{"signed"}},
			"release.channel":      {Type: FactString, Complete: true, String: stringPointer("production")},
		},
	}
	evidence.ProgramSHA256 = ProgramSHA256(program)
	return program, evidence, now
}

func TestEvaluatePassFailAndNoPartialFail(t *testing.T) {
	program, evidence, now := fixture()
	if result := Evaluate(program, evidence, now); result.Outcome != OutcomePass || result.ReasonCode != ReasonPassed {
		t.Fatalf("complete true predicate did not pass: %#v", result)
	}

	failing := evidence
	failing.Facts = cloneFacts(evidence.Facts)
	failing.Facts["policy.score"] = Fact{Type: FactNumber, Complete: true, Number: json.Number("89.999")}
	if result := Evaluate(program, failing, now); result.Outcome != OutcomeFail || result.ReasonCode != ReasonPredicateFalse {
		t.Fatalf("complete authoritative negative did not fail: %#v", result)
	}

	partial := failing
	partial.Facts = cloneFacts(failing.Facts)
	delete(partial.Facts, "manifest.schema")
	if result := Evaluate(program, partial, now); result.Outcome != OutcomeBlocked || result.ReasonCode != ReasonFactMissing {
		t.Fatalf("false plus missing evidence must block, not fail: %#v", result)
	}
}

func TestAllAnyNotCompositionRequiresEveryReferencedFact(t *testing.T) {
	program, evidence, now := fixture()
	program.Predicate = Expression{Op: OpAll, Args: []Expression{
		{Op: OpAny, Args: []Expression{
			{Op: OpStringEqual, Fact: "release.channel", String: stringPointer("development")},
			{Op: OpStateIn, Fact: "signature.state", Strings: []string{"trusted", "verified"}},
		}},
		{Op: OpNot, Arg: &Expression{Op: OpBooleanEqual, Fact: "signature.present", Boolean: boolPointer(false)}},
	}}
	evidence.ProgramSHA256 = ProgramSHA256(program)
	if result := Evaluate(program, evidence, now); result.Outcome != OutcomePass {
		t.Fatalf("valid all/any/not expression did not pass: %#v", result)
	}
	delete(evidence.Facts, "release.channel")
	if result := Evaluate(program, evidence, now); result.Outcome != OutcomeBlocked || result.ReasonCode != ReasonFactMissing {
		t.Fatalf("true any sibling must not hide a missing fact: %#v", result)
	}
}

func TestEveryLeafOperation(t *testing.T) {
	digest := digestString("release artifact")
	tests := []Expression{
		{Op: OpIdentityEqual, Fact: "artifact.identity", String: stringPointer("release/app@sha256:abc")},
		{Op: OpSchemaEqual, Fact: "manifest.schema", String: stringPointer("spdx-2.3")},
		{Op: OpDigestEqual, Fact: "artifact.digest", String: &digest},
		{Op: OpStringEqual, Fact: "release.channel", String: stringPointer("production")},
		{Op: OpStateIn, Fact: "signature.state", Strings: []string{"trusted", "verified"}},
		{Op: OpBooleanEqual, Fact: "signature.present", Boolean: boolPointer(true)},
		{Op: OpNumberEqual, Fact: "policy.score", Number: json.Number("90.0")},
		{Op: OpNumberLess, Fact: "policy.score", Number: json.Number("91")},
		{Op: OpNumberLessEq, Fact: "policy.score", Number: json.Number("90")},
		{Op: OpNumberGreater, Fact: "policy.score", Number: json.Number("89")},
		{Op: OpNumberGreaterEq, Fact: "policy.score", Number: json.Number("90")},
		{Op: OpTimeBefore, Fact: "signature.expires_at", Timestamp: stringPointer("2026-08-30T00:00:00Z")},
		{Op: OpTimeBeforeEq, Fact: "signature.expires_at", Timestamp: stringPointer("2026-08-29T23:00:00Z")},
		{Op: OpTimeAfter, Fact: "signature.expires_at", Timestamp: stringPointer("2026-08-29T22:00:00Z")},
		{Op: OpTimeAfterEq, Fact: "signature.expires_at", Timestamp: stringPointer("2026-08-29T23:00:00Z")},
		{Op: OpSetEqual, Fact: "manifest.subjects", Strings: []string{"artifact", "manifest", "signature"}},
		{Op: OpSetContainsAll, Fact: "manifest.subjects", Strings: []string{"artifact", "manifest"}},
		{Op: OpSetDisjoint, Fact: "manifest.flags", Strings: []string{"revoked", "unsigned"}},
	}
	for _, predicate := range tests {
		t.Run(string(predicate.Op), func(t *testing.T) {
			program, evidence, now := fixture()
			program.Predicate = predicate
			evidence.ProgramSHA256 = ProgramSHA256(program)
			if result := Evaluate(program, evidence, now); result.Outcome != OutcomePass {
				t.Fatalf("%s did not pass exact matching evidence: %#v", predicate.Op, result)
			}
		})
	}
}

func TestEveryFactToFactOperation(t *testing.T) {
	digest := digestString("release artifact")
	tests := []struct {
		predicate Expression
		facts     map[string]Fact
	}{
		{Expression{Op: OpIdentityEqualFact, Fact: "left", OtherFact: "right"}, map[string]Fact{
			"left":  {Type: FactIdentity, Complete: true, String: stringPointer("release/app@sha256:abc")},
			"right": {Type: FactIdentity, Complete: true, String: stringPointer("release/app@sha256:abc")},
		}},
		{Expression{Op: OpSchemaEqualFact, Fact: "left", OtherFact: "right"}, map[string]Fact{
			"left":  {Type: FactSchema, Complete: true, String: stringPointer("spdx-2.3")},
			"right": {Type: FactSchema, Complete: true, String: stringPointer("spdx-2.3")},
		}},
		{Expression{Op: OpDigestEqualFact, Fact: "left", OtherFact: "right"}, map[string]Fact{
			"left":  {Type: FactDigest, Complete: true, String: &digest},
			"right": {Type: FactDigest, Complete: true, String: &digest},
		}},
		{Expression{Op: OpStringEqualFact, Fact: "left", OtherFact: "right"}, map[string]Fact{
			"left":  {Type: FactString, Complete: true, String: stringPointer("production")},
			"right": {Type: FactString, Complete: true, String: stringPointer("production")},
		}},
		{Expression{Op: OpStateInSetFact, Fact: "left", OtherFact: "right"}, map[string]Fact{
			"left":  {Type: FactState, Complete: true, String: stringPointer("verified")},
			"right": {Type: FactStringSet, Complete: true, Strings: []string{"trusted", "verified"}},
		}},
		{Expression{Op: OpBooleanEqualFact, Fact: "left", OtherFact: "right"}, map[string]Fact{
			"left":  {Type: FactBoolean, Complete: true, Boolean: boolPointer(true)},
			"right": {Type: FactBoolean, Complete: true, Boolean: boolPointer(true)},
		}},
		{Expression{Op: OpNumberEqualFact, Fact: "left", OtherFact: "right"}, numberPair("90", "90.0")},
		{Expression{Op: OpNumberLessFact, Fact: "left", OtherFact: "right"}, numberPair("89", "90")},
		{Expression{Op: OpNumberLessEqFact, Fact: "left", OtherFact: "right"}, numberPair("90", "90")},
		{Expression{Op: OpNumberGreaterFact, Fact: "left", OtherFact: "right"}, numberPair("91", "90")},
		{Expression{Op: OpNumberGreaterEqFact, Fact: "left", OtherFact: "right"}, numberPair("90", "90")},
		{Expression{Op: OpTimeBeforeFact, Fact: "left", OtherFact: "right"}, timePair("2026-08-29T22:00:00Z", "2026-08-29T23:00:00Z")},
		{Expression{Op: OpTimeBeforeEqFact, Fact: "left", OtherFact: "right"}, timePair("2026-08-29T23:00:00Z", "2026-08-29T23:00:00Z")},
		{Expression{Op: OpTimeAfterFact, Fact: "left", OtherFact: "right"}, timePair("2026-08-30T00:00:00Z", "2026-08-29T23:00:00Z")},
		{Expression{Op: OpTimeAfterEqFact, Fact: "left", OtherFact: "right"}, timePair("2026-08-29T23:00:00Z", "2026-08-29T23:00:00Z")},
		{Expression{Op: OpSetEqualFact, Fact: "left", OtherFact: "right"}, setPair([]string{"a", "b"}, []string{"a", "b"})},
		{Expression{Op: OpSetContainsAllFact, Fact: "left", OtherFact: "right"}, setPair([]string{"a", "b", "c"}, []string{"a", "b"})},
		{Expression{Op: OpSetDisjointFact, Fact: "left", OtherFact: "right"}, setPair([]string{"a", "b"}, []string{"c", "d"})},
	}
	for _, test := range tests {
		t.Run(string(test.predicate.Op), func(t *testing.T) {
			program, evidence, now := fixture()
			program.Predicate = test.predicate
			evidence.Facts = test.facts
			evidence.ProgramSHA256 = ProgramSHA256(program)
			if result := Evaluate(program, evidence, now); result.Outcome != OutcomePass {
				t.Fatalf("%s did not pass matching facts: %#v", test.predicate.Op, result)
			}

			delete(evidence.Facts, "right")
			if result := Evaluate(program, evidence, now); result.Outcome != OutcomeBlocked || result.ReasonCode != ReasonFactMissing {
				t.Fatalf("%s missing right fact did not block: %#v", test.predicate.Op, result)
			}
		})
	}
}

func TestFactToFactTypeMismatchBlocksInsteadOfFailing(t *testing.T) {
	program, evidence, now := fixture()
	program.Predicate = Expression{Op: OpDigestEqualFact, Fact: "artifact.digest", OtherFact: "other"}
	evidence.Facts = map[string]Fact{
		"artifact.digest": {Type: FactDigest, Complete: true, String: stringPointer(digestString("artifact"))},
		"other":           {Type: FactString, Complete: true, String: stringPointer(digestString("artifact"))},
	}
	evidence.ProgramSHA256 = ProgramSHA256(program)
	if result := Evaluate(program, evidence, now); result.Outcome != OutcomeBlocked || result.ReasonCode != ReasonFactTypeMismatch {
		t.Fatalf("fact-to-fact type mismatch must block: %#v", result)
	}
}

func TestSealedParametersAreProgramBoundAndProviderCannotChooseExpectedValue(t *testing.T) {
	program, evidence, now := fixture()
	program.Parameters = map[string]Parameter{
		"minimum_score": {Type: FactNumber, Number: json.Number("90")},
	}
	program.Predicate = Expression{
		Op: OpNumberGreaterEqParameter, Fact: "policy.score", Parameter: "minimum_score",
	}
	evidence.ProgramSHA256 = ProgramSHA256(program)
	if result := Evaluate(program, evidence, now); result.Outcome != OutcomePass {
		t.Fatalf("measurement meeting sealed threshold did not pass: %#v", result)
	}

	changed := program
	changed.Parameters = map[string]Parameter{
		"minimum_score": {Type: FactNumber, Number: json.Number("95")},
	}
	// Evidence remains bound to the old program. A provider cannot silently
	// lower or replace the expected value because that changes ProgramSHA256.
	if result := Evaluate(changed, evidence, now); result.Outcome != OutcomeBlocked || result.ReasonCode != ReasonEvidenceBindingMismatch {
		t.Fatalf("policy-parameter change without evidence rebinding must block: %#v", result)
	}

	evidence.ProgramSHA256 = ProgramSHA256(changed)
	if result := Evaluate(changed, evidence, now); result.Outcome != OutcomeFail || result.ReasonCode != ReasonPredicateFalse {
		t.Fatalf("measurement below sealed threshold must fail: %#v", result)
	}
}

func TestMissingOrWrongTypedProgramParameterIsInvalid(t *testing.T) {
	program, evidence, now := fixture()
	program.Predicate = Expression{Op: OpDigestEqualParameter, Fact: "artifact.digest", Parameter: "expected_digest"}
	if result := Evaluate(program, evidence, now); result.Outcome != OutcomeBlocked || result.ReasonCode != ReasonInvalidProgram {
		t.Fatalf("missing sealed parameter must invalidate program: %#v", result)
	}
	program.Parameters["expected_digest"] = Parameter{Type: FactString, String: stringPointer(digestString("artifact"))}
	if result := Evaluate(program, evidence, now); result.Outcome != OutcomeBlocked || result.ReasonCode != ReasonInvalidProgram {
		t.Fatalf("wrong typed sealed parameter must invalidate program: %#v", result)
	}
}

func TestSealedNumericAndTimeMapsCompareEveryRequiredKey(t *testing.T) {
	program, evidence, now := fixture()
	program.Parameters = map[string]Parameter{
		"maximums":  {Type: FactNumberMap, Numbers: map[string]json.Number{"latency.p95": "250", "queue.age": "30"}},
		"deadlines": {Type: FactTimeMap, Timestamps: map[string]string{"certificate": "2026-09-01T00:00:00Z"}},
	}
	program.Predicate = Expression{Op: OpAll, Args: []Expression{
		{Op: OpNumberMapLessEqParameter, Fact: "measurements", Parameter: "maximums"},
		{Op: OpTimeMapBeforeParameter, Fact: "events", Parameter: "deadlines"},
	}}
	evidence.Facts = map[string]Fact{
		"measurements": {Type: FactNumberMap, Complete: true, Numbers: map[string]json.Number{"latency.p95": "240", "queue.age": "20"}},
		"events":       {Type: FactTimeMap, Complete: true, Timestamps: map[string]string{"certificate": "2026-08-31T00:00:00Z"}},
	}
	evidence.ProgramSHA256 = ProgramSHA256(program)
	if result := Evaluate(program, evidence, now); result.Outcome != OutcomePass {
		t.Fatalf("complete maps meeting sealed policy did not pass: %#v", result)
	}

	missing := evidence
	missing.Facts = cloneFacts(evidence.Facts)
	missing.Facts["measurements"] = Fact{Type: FactNumberMap, Complete: true, Numbers: map[string]json.Number{"latency.p95": "240"}}
	if result := Evaluate(program, missing, now); result.Outcome != OutcomeFail {
		t.Fatalf("complete observed map missing a required metric must fail: %#v", result)
	}

	extra := evidence
	extra.Facts = cloneFacts(evidence.Facts)
	extra.Facts["measurements"] = Fact{Type: FactNumberMap, Complete: true, Numbers: map[string]json.Number{
		"latency.p95": "240", "queue.age": "20", "unreviewed.metric": "0",
	}}
	if result := Evaluate(program, extra, now); result.Outcome != OutcomeFail {
		t.Fatalf("complete observed map with an unsealed metric must fail: %#v", result)
	}

	incomplete := missing
	incomplete.Facts = cloneFacts(missing.Facts)
	fact := incomplete.Facts["measurements"]
	fact.Complete = false
	incomplete.Facts["measurements"] = fact
	if result := Evaluate(program, incomplete, now); result.Outcome != OutcomeBlocked || result.ReasonCode != ReasonEvidenceIncomplete {
		t.Fatalf("incomplete observed map must block, not fail: %#v", result)
	}
}

func TestTimeMapEqualityUsesInstantsAndExactDomains(t *testing.T) {
	program, evidence, now := fixture()
	program.Parameters = map[string]Parameter{
		"expected": {Type: FactTimeMap, Timestamps: map[string]string{
			"release": "2026-08-29T12:00:00Z",
		}},
	}
	program.Predicate = Expression{Op: OpAll, Args: []Expression{
		{Op: OpTimeMapEqualParameter, Fact: "observed", Parameter: "expected"},
		{Op: OpTimeMapEqualFact, Fact: "observed", OtherFact: "canonical"},
	}}
	evidence.Facts = map[string]Fact{
		"observed":  {Type: FactTimeMap, Complete: true, Timestamps: map[string]string{"release": "2026-08-29T14:00:00+02:00"}},
		"canonical": {Type: FactTimeMap, Complete: true, Timestamps: map[string]string{"release": "2026-08-29T12:00:00Z"}},
	}
	evidence.ProgramSHA256 = ProgramSHA256(program)
	if result := Evaluate(program, evidence, now); result.Outcome != OutcomePass {
		t.Fatalf("equal instants with different offsets did not pass: %#v", result)
	}

	extra := evidence
	extra.Facts = cloneFacts(evidence.Facts)
	extra.Facts["observed"] = Fact{Type: FactTimeMap, Complete: true, Timestamps: map[string]string{
		"release": "2026-08-29T12:00:00Z", "unsealed": "2026-08-29T12:00:00Z",
	}}
	if result := Evaluate(program, extra, now); result.Outcome != OutcomeFail {
		t.Fatalf("time equality ignored an extra map key: %#v", result)
	}
}

func TestNumericAndTimeMapsCompareExactDomainsFactToFact(t *testing.T) {
	program, evidence, now := fixture()
	program.Predicate = Expression{Op: OpAll, Args: []Expression{
		{Op: OpNumberMapLessEqFact, Fact: "observed_numbers", OtherFact: "maximum_numbers"},
		{Op: OpTimeMapBeforeFact, Fact: "review_times", OtherFact: "landing_times"},
	}}
	evidence.Facts = map[string]Fact{
		"observed_numbers": {Type: FactNumberMap, Complete: true, Numbers: map[string]json.Number{"a": "1", "b": "2"}},
		"maximum_numbers":  {Type: FactNumberMap, Complete: true, Numbers: map[string]json.Number{"a": "1", "b": "3"}},
		"review_times": {Type: FactTimeMap, Complete: true, Timestamps: map[string]string{
			"change-a": "2026-08-29T10:00:00Z", "change-b": "2026-08-29T10:30:00Z",
		}},
		"landing_times": {Type: FactTimeMap, Complete: true, Timestamps: map[string]string{
			"change-a": "2026-08-29T11:00:00Z", "change-b": "2026-08-29T11:30:00Z",
		}},
	}
	evidence.ProgramSHA256 = ProgramSHA256(program)
	if result := Evaluate(program, evidence, now); result.Outcome != OutcomePass {
		t.Fatalf("exact fact-to-fact map comparisons did not pass: %#v", result)
	}

	extraKey := evidence
	extraKey.Facts = cloneFacts(evidence.Facts)
	extraKey.Facts["review_times"] = Fact{Type: FactTimeMap, Complete: true, Timestamps: map[string]string{
		"change-a": "2026-08-29T10:00:00Z", "change-b": "2026-08-29T10:30:00Z", "unmatched": "2026-08-29T09:00:00Z",
	}}
	if result := Evaluate(program, extraKey, now); result.Outcome != OutcomeFail {
		t.Fatalf("map comparison ignored an unmatched relation key: %#v", result)
	}
}

func TestTypedMapsCompareExactDomains(t *testing.T) {
	program, evidence, now := fixture()
	digestA := digestString("artifact-a")
	program.Parameters = map[string]Parameter{
		"expected_ids":     {Type: FactIdentityMap, Values: map[string]string{"release": "release-42"}},
		"expected_schemas": {Type: FactSchemaMap, Values: map[string]string{"manifest": "spdx-2.3"}},
		"expected_digests": {Type: FactDigestMap, Values: map[string]string{"artifact": digestA}},
		"expected_states":  {Type: FactStateMap, Values: map[string]string{"release": "approved"}},
		"expected_labels":  {Type: FactStringMap, Values: map[string]string{"channel": "production"}},
		"expected_flags":   {Type: FactBooleanMap, Booleans: map[string]bool{"blocked": false}},
	}
	program.Predicate = Expression{Op: OpAll, Args: []Expression{
		{Op: OpIdentityMapEqualParameter, Fact: "ids", Parameter: "expected_ids"},
		{Op: OpSchemaMapEqualParameter, Fact: "schemas", Parameter: "expected_schemas"},
		{Op: OpDigestMapEqualParameter, Fact: "digests", Parameter: "expected_digests"},
		{Op: OpStateMapEqualParameter, Fact: "states", Parameter: "expected_states"},
		{Op: OpStringMapEqualParameter, Fact: "labels", Parameter: "expected_labels"},
		{Op: OpBooleanMapEqualParameter, Fact: "flags", Parameter: "expected_flags"},
	}}
	evidence.Facts = map[string]Fact{
		"ids":     {Type: FactIdentityMap, Complete: true, Values: map[string]string{"release": "release-42"}},
		"schemas": {Type: FactSchemaMap, Complete: true, Values: map[string]string{"manifest": "spdx-2.3"}},
		"digests": {Type: FactDigestMap, Complete: true, Values: map[string]string{"artifact": digestA}},
		"states":  {Type: FactStateMap, Complete: true, Values: map[string]string{"release": "approved"}},
		"labels":  {Type: FactStringMap, Complete: true, Values: map[string]string{"channel": "production"}},
		"flags":   {Type: FactBooleanMap, Complete: true, Booleans: map[string]bool{"blocked": false}},
	}
	evidence.ProgramSHA256 = ProgramSHA256(program)
	if result := Evaluate(program, evidence, now); result.Outcome != OutcomePass {
		t.Fatalf("exact typed maps did not pass: %#v", result)
	}

	extra := evidence
	extra.Facts = cloneFacts(evidence.Facts)
	extra.Facts["states"] = Fact{Type: FactStateMap, Complete: true, Values: map[string]string{
		"release": "approved", "unreviewed": "approved",
	}}
	if result := Evaluate(program, extra, now); result.Outcome != OutcomeFail {
		t.Fatalf("typed map comparison ignored an extra key: %#v", result)
	}

	program.Parameters = map[string]Parameter{
		"state_keys":     {Type: FactStringSet, Strings: []string{"release"}},
		"allowed_ids":    {Type: FactStringSet, Strings: []string{"release-42", "release-43"}},
		"allowed_states": {Type: FactStringSet, Strings: []string{"approved", "released"}},
		"allowed_labels": {Type: FactStringSet, Strings: []string{"production", "staging"}},
		"required_flag":  {Type: FactBoolean, Boolean: boolPointer(false)},
	}
	program.Predicate = Expression{Op: OpAll, Args: []Expression{
		{Op: OpMapKeysEqualSetParameter, Fact: "states", Parameter: "state_keys"},
		{Op: OpIdentityMapValuesInParameter, Fact: "ids", Parameter: "allowed_ids"},
		{Op: OpStateMapValuesInParameter, Fact: "states", Parameter: "allowed_states"},
		{Op: OpStringMapValuesInParameter, Fact: "labels", Parameter: "allowed_labels"},
		{Op: OpBooleanMapAllEqualParameter, Fact: "flags", Parameter: "required_flag"},
		{Op: OpBooleanMapAllEqual, Fact: "flags", Boolean: boolPointer(false)},
		{Op: OpStringMapAllNonempty, Fact: "labels"},
	}}
	evidence.ProgramSHA256 = ProgramSHA256(program)
	if result := Evaluate(program, evidence, now); result.Outcome != OutcomePass {
		t.Fatalf("typed map key/value predicates did not pass: %#v", result)
	}
	unauthorizedIdentity := evidence
	unauthorizedIdentity.Facts = cloneFacts(evidence.Facts)
	unauthorizedIdentity.Facts["ids"] = Fact{Type: FactIdentityMap, Complete: true, Values: map[string]string{"release": "release-99"}}
	if result := Evaluate(program, unauthorizedIdentity, now); result.Outcome != OutcomeFail {
		t.Fatalf("identity map accepted a value outside the sealed identity set: %#v", result)
	}
	evidence.Facts["state_keys"] = Fact{Type: FactStringSet, Complete: true, Strings: []string{"release"}}
	program.Parameters = map[string]Parameter{}
	program.Predicate = Expression{Op: OpMapKeysEqualSetFact, Fact: "states", OtherFact: "state_keys"}
	evidence.ProgramSHA256 = ProgramSHA256(program)
	if result := Evaluate(program, evidence, now); result.Outcome != OutcomePass {
		t.Fatalf("fact-to-fact map key comparison did not pass: %#v", result)
	}

	program.Parameters = map[string]Parameter{}
	program.Predicate = Expression{Op: OpAll, Args: []Expression{
		{Op: OpIdentityMapEqualFact, Fact: "ids", OtherFact: "ids_copy"},
		{Op: OpSchemaMapEqualFact, Fact: "schemas", OtherFact: "schemas_copy"},
		{Op: OpDigestMapEqualFact, Fact: "digests", OtherFact: "digests_copy"},
		{Op: OpStateMapEqualFact, Fact: "states", OtherFact: "states_copy"},
		{Op: OpStringMapEqualFact, Fact: "labels", OtherFact: "labels_copy"},
		{Op: OpBooleanMapEqualFact, Fact: "flags", OtherFact: "flags_copy"},
	}}
	for key, fact := range cloneFacts(evidence.Facts) {
		evidence.Facts[key+"_copy"] = fact
	}
	evidence.ProgramSHA256 = ProgramSHA256(program)
	if result := Evaluate(program, evidence, now); result.Outcome != OutcomePass {
		t.Fatalf("fact-to-fact typed maps did not pass: %#v", result)
	}

	program.Predicate = Expression{Op: OpIdentityMapValuesNotEqualFact, Fact: "ids", OtherFact: "approver_ids"}
	evidence.Facts["approver_ids"] = Fact{Type: FactIdentityMap, Complete: true, Values: map[string]string{"release": "approver-7"}}
	evidence.ProgramSHA256 = ProgramSHA256(program)
	if result := Evaluate(program, evidence, now); result.Outcome != OutcomePass {
		t.Fatalf("per-key identity separation did not pass: %#v", result)
	}
	evidence.Facts["approver_ids"] = Fact{Type: FactIdentityMap, Complete: true, Values: map[string]string{"release": "release-42"}}
	if result := Evaluate(program, evidence, now); result.Outcome != OutcomeFail {
		t.Fatalf("per-key identity separation accepted the same identity: %#v", result)
	}

	program.Predicate = Expression{Op: OpIdentityMapValuesUnique, Fact: "ids"}
	evidence.ProgramSHA256 = ProgramSHA256(program)
	if result := Evaluate(program, evidence, now); result.Outcome != OutcomePass {
		t.Fatalf("unique identity map did not pass: %#v", result)
	}
	evidence.Facts["ids"] = Fact{Type: FactIdentityMap, Complete: true, Values: map[string]string{"first": "same-id", "second": "same-id"}}
	if result := Evaluate(program, evidence, now); result.Outcome != OutcomeFail {
		t.Fatalf("duplicate identity-map values were accepted: %#v", result)
	}

	program.Predicate = Expression{Op: OpIdentityMapValuesInFact, Fact: "references", OtherFact: "inventory"}
	evidence.Facts["references"] = Fact{Type: FactIdentityMap, Complete: true, Values: map[string]string{"record-a": "item-1", "record-b": "item-1"}}
	evidence.Facts["inventory"] = Fact{Type: FactIdentityMap, Complete: true, Values: map[string]string{"first": "item-1", "second": "item-2"}}
	evidence.ProgramSHA256 = ProgramSHA256(program)
	if result := Evaluate(program, evidence, now); result.Outcome != OutcomePass {
		t.Fatalf("valid identity references did not resolve: %#v", result)
	}
	evidence.Facts["references"] = Fact{Type: FactIdentityMap, Complete: true, Values: map[string]string{"record-a": "missing-item"}}
	if result := Evaluate(program, evidence, now); result.Outcome != OutcomeFail {
		t.Fatalf("unknown identity reference was accepted: %#v", result)
	}
	evidence.Facts["references"] = Fact{Type: FactIdentityMap, Complete: true, Values: map[string]string{"record-a": "item-1"}}
	evidence.Facts["inventory"] = Fact{Type: FactIdentityMap, Complete: true, Values: map[string]string{"first": "item-1", "duplicate": "item-1"}}
	if result := Evaluate(program, evidence, now); result.Outcome != OutcomeFail {
		t.Fatalf("ambiguous duplicate inventory identity was accepted: %#v", result)
	}

	program.Parameters = map[string]Parameter{
		"incompatible_roles": {Type: FactDirectedGraph, Edges: []DirectedEdge{{From: "approver", To: "author"}}},
	}
	program.Predicate = Expression{Op: OpIdentityMapValuesDifferForPairsParameter, Fact: "role_holders", Parameter: "incompatible_roles"}
	evidence.Facts = map[string]Fact{
		"role_holders": {Type: FactIdentityMap, Complete: true, Values: map[string]string{
			"approver": "person-2", "author": "person-1", "reader": "person-1",
		}},
	}
	evidence.ProgramSHA256 = ProgramSHA256(program)
	if result := Evaluate(program, evidence, now); result.Outcome != OutcomePass {
		t.Fatalf("compatible role reuse was rejected: %#v", result)
	}
	evidence.Facts["role_holders"] = Fact{Type: FactIdentityMap, Complete: true, Values: map[string]string{
		"approver": "person-1", "author": "person-1", "reader": "person-3",
	}}
	if result := Evaluate(program, evidence, now); result.Outcome != OutcomeFail {
		t.Fatalf("incompatible roles with the same identity were accepted: %#v", result)
	}
	program.Parameters["incompatible_roles"] = Parameter{
		Type: FactDirectedGraph, Edges: []DirectedEdge{{From: "author", To: "approver"}},
	}
	if err := ValidateProgram(program); err == nil {
		t.Fatal("noncanonical identity pair set was accepted")
	}

	program.Parameters = map[string]Parameter{
		"required_approvals": {Type: FactStringSet, Strings: []string{"approval-a", "approval-b"}},
	}
	program.Predicate = Expression{Op: OpIdentityMapValuesBijectSetParameter, Fact: "approval_ids", Parameter: "required_approvals"}
	evidence.Facts = map[string]Fact{
		"approval_ids": {Type: FactIdentityMap, Complete: true, Values: map[string]string{"security": "approval-a", "service": "approval-b"}},
	}
	evidence.ProgramSHA256 = ProgramSHA256(program)
	if result := Evaluate(program, evidence, now); result.Outcome != OutcomePass {
		t.Fatalf("exact approval bijection did not pass: %#v", result)
	}
	evidence.Facts["approval_ids"] = Fact{Type: FactIdentityMap, Complete: true, Values: map[string]string{"security": "approval-a", "service": "approval-a"}}
	if result := Evaluate(program, evidence, now); result.Outcome != OutcomeFail {
		t.Fatalf("duplicate approval identity was accepted: %#v", result)
	}
}

func TestPerKeyRelationsTimeDeltaAndAcyclicGraph(t *testing.T) {
	program, evidence, now := fixture()
	program.Parameters = map[string]Parameter{
		"maximum_seconds": {Type: FactNumber, Number: json.Number("3600")},
	}
	program.Predicate = Expression{Op: OpAll, Args: []Expression{
		{Op: OpBooleanMapAnyTrueFact, Fact: "primary", OtherFact: "secondary"},
		{Op: OpBooleanMapImpliesFact, Fact: "required", OtherFact: "completed"},
		{Op: OpStringMapAnyNonemptyFact, Fact: "person", OtherFact: "organization"},
		{Op: OpTimeMapDeltaLessEqParameter, Fact: "started", OtherFact: "finished", Parameter: "maximum_seconds"},
		{Op: OpTimeMapDeltaEqualNumberMapFact, Fact: "started", OtherFact: "finished", ThirdFact: "duration_seconds"},
		{Op: OpDirectedGraphAcyclic, Fact: "dependencies"},
	}}
	evidence.Facts = map[string]Fact{
		"primary":      {Type: FactBooleanMap, Complete: true, Booleans: map[string]bool{"a": true, "b": false}},
		"secondary":    {Type: FactBooleanMap, Complete: true, Booleans: map[string]bool{"a": false, "b": true}},
		"required":     {Type: FactBooleanMap, Complete: true, Booleans: map[string]bool{"a": true, "b": false}},
		"completed":    {Type: FactBooleanMap, Complete: true, Booleans: map[string]bool{"a": true, "b": false}},
		"person":       {Type: FactStringMap, Complete: true, Values: map[string]string{"a": "alice", "b": ""}},
		"organization": {Type: FactStringMap, Complete: true, Values: map[string]string{"a": "", "b": "operations"}},
		"started": {Type: FactTimeMap, Complete: true, Timestamps: map[string]string{
			"a": "2026-08-29T10:00:00Z", "b": "2026-08-29T11:00:00Z",
		}},
		"finished": {Type: FactTimeMap, Complete: true, Timestamps: map[string]string{
			"a": "2026-08-29T10:30:00Z", "b": "2026-08-29T12:00:00Z",
		}},
		"duration_seconds": {Type: FactNumberMap, Complete: true, Numbers: map[string]json.Number{"a": "1800", "b": "3600"}},
		"dependencies": {Type: FactDirectedGraph, Complete: true, Edges: []DirectedEdge{
			{From: "api", To: "database"}, {From: "web", To: "api"},
		}},
	}
	evidence.ProgramSHA256 = ProgramSHA256(program)
	if result := Evaluate(program, evidence, now); result.Outcome != OutcomePass {
		t.Fatalf("bounded per-key relations did not pass: %#v", result)
	}

	cycle := evidence
	cycle.Facts = cloneFacts(evidence.Facts)
	cycle.Facts["dependencies"] = Fact{Type: FactDirectedGraph, Complete: true, Edges: []DirectedEdge{
		{From: "api", To: "database"}, {From: "database", To: "web"}, {From: "web", To: "api"},
	}}
	if result := Evaluate(program, cycle, now); result.Outcome != OutcomeFail {
		t.Fatalf("directed cycle did not fail: %#v", result)
	}

	overtime := evidence
	overtime.Facts = cloneFacts(evidence.Facts)
	overtime.Facts["finished"] = Fact{Type: FactTimeMap, Complete: true, Timestamps: map[string]string{
		"a": "2026-08-29T10:30:00Z", "b": "2026-08-29T12:00:01Z",
	}}
	if result := Evaluate(program, overtime, now); result.Outcome != OutcomeFail {
		t.Fatalf("over-limit time delta did not fail: %#v", result)
	}

	misreported := evidence
	misreported.Facts = cloneFacts(evidence.Facts)
	misreported.Facts["duration_seconds"] = Fact{Type: FactNumberMap, Complete: true, Numbers: map[string]json.Number{"a": "1799", "b": "3600"}}
	if result := Evaluate(program, misreported, now); result.Outcome != OutcomeFail {
		t.Fatalf("misreported stage duration did not fail: %#v", result)
	}
}

func TestMapKeyPartitionRequiresOneExactRoutePerSubject(t *testing.T) {
	program, evidence, now := fixture()
	program.Parameters = map[string]Parameter{
		"subjects": {Type: FactStringSet, Strings: []string{"change-a", "change-b"}},
	}
	program.Predicate = Expression{
		Op: OpMapKeyPartitionEqualSetParameter, Fact: "direct", OtherFact: "delegated", Parameter: "subjects",
	}
	evidence.Facts = map[string]Fact{
		"direct":    {Type: FactIdentityMap, Complete: true, Values: map[string]string{"change-a": "owner-1"}},
		"delegated": {Type: FactStringMap, Complete: true, Values: map[string]string{"change-b": "owner-1>delegate-2"}},
	}
	evidence.ProgramSHA256 = ProgramSHA256(program)
	if result := Evaluate(program, evidence, now); result.Outcome != OutcomePass {
		t.Fatalf("exact route partition did not pass: %#v", result)
	}

	overlap := evidence
	overlap.Facts = cloneFacts(evidence.Facts)
	overlap.Facts["delegated"] = Fact{Type: FactStringMap, Complete: true, Values: map[string]string{
		"change-a": "owner-1>delegate-2", "change-b": "owner-1>delegate-2",
	}}
	if result := Evaluate(program, overlap, now); result.Outcome != OutcomeFail {
		t.Fatalf("subject assigned to two routes was accepted: %#v", result)
	}

	missing := evidence
	missing.Facts = cloneFacts(evidence.Facts)
	missing.Facts["delegated"] = Fact{Type: FactStringMap, Complete: true, Values: map[string]string{}}
	if result := Evaluate(program, missing, now); result.Outcome != OutcomeFail {
		t.Fatalf("partition missing a sealed subject was accepted: %#v", result)
	}
}

func TestLiteralMapVocabulariesCannotBeChangedByPolicy(t *testing.T) {
	program, evidence, now := fixture()
	program.Parameters = map[string]Parameter{}
	program.Predicate = Expression{Op: OpAll, Args: []Expression{
		{Op: OpMapKeysEqualSet, Fact: "states", Strings: []string{"build", "deploy"}},
		{Op: OpStateMapValuesIn, Fact: "states", Strings: []string{"blocked", "passed"}},
		{Op: OpStringMapValuesIn, Fact: "classes", Strings: []string{"command", "query"}},
		{Op: OpIdentityMapValuesIn, Fact: "roles", Strings: []string{"approver", "author"}},
	}}
	evidence.Facts = map[string]Fact{
		"states":  {Type: FactStateMap, Complete: true, Values: map[string]string{"build": "passed", "deploy": "blocked"}},
		"classes": {Type: FactStringMap, Complete: true, Values: map[string]string{"a": "command", "b": "query"}},
		"roles":   {Type: FactIdentityMap, Complete: true, Values: map[string]string{"first": "approver", "second": "author"}},
	}
	evidence.ProgramSHA256 = ProgramSHA256(program)
	if result := Evaluate(program, evidence, now); result.Outcome != OutcomePass {
		t.Fatalf("fixed map vocabularies did not pass: %#v", result)
	}
	evidence.Facts["states"] = Fact{Type: FactStateMap, Complete: true, Values: map[string]string{"build": "passed", "deploy": "provider-chosen"}}
	if result := Evaluate(program, evidence, now); result.Outcome != OutcomeFail {
		t.Fatalf("provider-selected absolute state was accepted: %#v", result)
	}
}

func TestFindingClassificationUsesExactPartitionAndBaselineMembership(t *testing.T) {
	program, evidence, now := fixture()
	program.Parameters = map[string]Parameter{
		"finding_ids": {Type: FactStringSet, Strings: []string{"finding-a", "finding-b"}},
	}
	program.Predicate = Expression{Op: OpAll, Args: []Expression{
		{Op: OpMapKeyPartitionEqualSetParameter, Fact: "introduced", OtherFact: "inherited", Parameter: "finding_ids"},
		{Op: OpIdentityMapValuesNotInFact, Fact: "introduced", OtherFact: "baseline"},
		{Op: OpIdentityMapValuesInFact, Fact: "inherited", OtherFact: "baseline"},
	}}
	evidence.Facts = map[string]Fact{
		"introduced": {Type: FactIdentityMap, Complete: true, Values: map[string]string{"finding-a": "finding-a"}},
		"inherited":  {Type: FactIdentityMap, Complete: true, Values: map[string]string{"finding-b": "finding-b"}},
		"baseline":   {Type: FactIdentityMap, Complete: true, Values: map[string]string{"old-b": "finding-b"}},
	}
	evidence.ProgramSHA256 = ProgramSHA256(program)
	if result := Evaluate(program, evidence, now); result.Outcome != OutcomePass {
		t.Fatalf("exact finding classification did not pass: %#v", result)
	}
	evidence.Facts["introduced"] = Fact{Type: FactIdentityMap, Complete: true, Values: map[string]string{"finding-a": "finding-b"}}
	if result := Evaluate(program, evidence, now); result.Outcome != OutcomeFail {
		t.Fatalf("baseline finding was accepted as introduced: %#v", result)
	}
}

func numberPair(left, right string) map[string]Fact {
	return map[string]Fact{
		"left":  {Type: FactNumber, Complete: true, Number: json.Number(left)},
		"right": {Type: FactNumber, Complete: true, Number: json.Number(right)},
	}
}

func timePair(left, right string) map[string]Fact {
	return map[string]Fact{
		"left":  {Type: FactTime, Complete: true, Timestamp: stringPointer(left)},
		"right": {Type: FactTime, Complete: true, Timestamp: stringPointer(right)},
	}
}

func setPair(left, right []string) map[string]Fact {
	return map[string]Fact{
		"left":  {Type: FactStringSet, Complete: true, Strings: left},
		"right": {Type: FactStringSet, Complete: true, Strings: right},
	}
}

func TestBlockedEvidenceConditions(t *testing.T) {
	conflict := digestString("conflict")
	tests := []struct {
		name   string
		change func(*Program, *Evidence, time.Time)
		reason ReasonCode
	}{
		{name: "type mismatch", change: func(_ *Program, evidence *Evidence, _ time.Time) {
			evidence.Facts["artifact.identity"] = Fact{Type: FactString, Complete: true, String: stringPointer("release/app@sha256:abc")}
		}, reason: ReasonFactTypeMismatch},
		{name: "incomplete envelope", change: func(_ *Program, evidence *Evidence, _ time.Time) { evidence.Complete = false }, reason: ReasonEvidenceIncomplete},
		{name: "incomplete fact", change: func(_ *Program, evidence *Evidence, _ time.Time) {
			fact := evidence.Facts["artifact.identity"]
			fact.Complete = false
			evidence.Facts["artifact.identity"] = fact
		}, reason: ReasonEvidenceIncomplete},
		{name: "conflicting envelope", change: func(_ *Program, evidence *Evidence, _ time.Time) { evidence.ContradictionDigests = []string{conflict} }, reason: ReasonEvidenceConflicting},
		{name: "conflicting fact", change: func(_ *Program, evidence *Evidence, _ time.Time) {
			fact := evidence.Facts["artifact.identity"]
			fact.ConflictDigests = []string{conflict}
			evidence.Facts["artifact.identity"] = fact
		}, reason: ReasonEvidenceConflicting},
		{name: "wrong authority", change: func(_ *Program, evidence *Evidence, _ time.Time) { evidence.Authority = AuthorityRepository }, reason: ReasonWrongAuthority},
		{name: "incomplete subject inventory", change: func(_ *Program, evidence *Evidence, _ time.Time) { evidence.ObservedSubjects = []string{"artifact"} }, reason: ReasonEvidenceIncomplete},
		{name: "stale", change: func(_ *Program, evidence *Evidence, now time.Time) { evidence.ObservedAt = now.Add(-2 * time.Hour) }, reason: ReasonEvidenceStale},
		{name: "future", change: func(_ *Program, evidence *Evidence, now time.Time) { evidence.ObservedAt = now.Add(time.Second) }, reason: ReasonEvidenceFromFuture},
		{name: "binding mismatch", change: func(_ *Program, evidence *Evidence, _ time.Time) { evidence.ControlRevision++ }, reason: ReasonEvidenceBindingMismatch},
		{name: "applicability unknown", change: func(_ *Program, evidence *Evidence, _ time.Time) { evidence.Applicability = ApplicabilityUnknown }, reason: ReasonApplicabilityUnknown},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			program, evidence, now := fixture()
			evidence.Facts = cloneFacts(evidence.Facts)
			test.change(&program, &evidence, now)
			result := Evaluate(program, evidence, now)
			if result.Outcome != OutcomeBlocked || result.ReasonCode != test.reason {
				t.Fatalf("got %s/%s, want blocked/%s: %#v", result.Outcome, result.ReasonCode, test.reason, result)
			}
		})
	}
}

func TestNotApplicableRequiresBoundProof(t *testing.T) {
	program, evidence, now := fixture()
	evidence.Applicability = ApplicabilityNotApplicable
	evidence.ApplicabilityProof = "authoritative inventory proves the module trigger is absent"
	evidence.ApplicabilityProofSHA256 = digestString(evidence.ApplicabilityProof)
	if result := Evaluate(program, evidence, now); result.Outcome != OutcomeNotApplicable || result.ReasonCode != ReasonNotApplicable {
		t.Fatalf("bounded not-applicable proof was rejected: %#v", result)
	}

	evidence.ApplicabilityProofSHA256 = digestString("different proof")
	if result := Evaluate(program, evidence, now); result.Outcome != OutcomeBlocked || result.ReasonCode != ReasonApplicabilityProofMissing {
		t.Fatalf("mismatched proof must block: %#v", result)
	}

	evidence.ApplicabilityProofSHA256 = digestString(evidence.ApplicabilityProof)
	program.AllowNotApplicable = false
	evidence.ProgramSHA256 = ProgramSHA256(program)
	if result := Evaluate(program, evidence, now); result.Outcome != OutcomeBlocked || result.ReasonCode != ReasonApplicabilityDisallowed {
		t.Fatalf("disallowed not-applicable must block: %#v", result)
	}
}

func TestEvaluationIsStableAndDoesNotMutateInputs(t *testing.T) {
	program, evidence, now := fixture()
	evidence.Facts["aaa.incomplete"] = Fact{Type: FactBoolean, Complete: false, Boolean: boolPointer(true)}
	evidence.Facts["zzz.conflict"] = Fact{Type: FactBoolean, Complete: true, Boolean: boolPointer(true), ConflictDigests: []string{digestString("conflict")}}
	originalProgram, originalEvidence := program, cloneEvidence(evidence)
	first := Evaluate(program, evidence, now)
	for range 100 {
		if next := Evaluate(program, evidence, now); !reflect.DeepEqual(first, next) {
			t.Fatalf("repeated evaluation changed output:\nfirst=%#v\nnext=%#v", first, next)
		}
	}
	if !reflect.DeepEqual(program, originalProgram) || !reflect.DeepEqual(evidence, originalEvidence) {
		t.Fatal("evaluation mutated its inputs")
	}
}

func TestStrictDecodeAndResourceBounds(t *testing.T) {
	program, evidence, _ := fixture()
	programJSON, _ := json.Marshal(program)
	evidenceJSON, _ := json.Marshal(evidence)
	if _, err := DecodeProgram(programJSON); err != nil {
		t.Fatalf("valid program did not decode: %v", err)
	}
	if _, err := DecodeEvidence(evidenceJSON); err != nil {
		t.Fatalf("valid evidence did not decode: %v", err)
	}

	duplicate := bytes.Replace(programJSON, []byte(`{"schema_version":`), []byte(`{"schema_version":"duplicate","schema_version":`), 1)
	if _, err := DecodeProgram(duplicate); err == nil || !strings.Contains(err.Error(), "duplicate") {
		t.Fatalf("duplicate key was not rejected: %v", err)
	}
	unknown := bytes.Replace(programJSON, []byte(`{"schema_version":`), []byte(`{"unknown":true,"schema_version":`), 1)
	if _, err := DecodeProgram(unknown); err == nil || !strings.Contains(err.Error(), "unknown") {
		t.Fatalf("unknown field was not rejected: %v", err)
	}
	if _, err := DecodeProgram(append(programJSON, []byte(` {}`)...)); err == nil {
		t.Fatal("trailing JSON was not rejected")
	}
	if _, err := DecodeProgram(bytes.Repeat([]byte{' '}, MaxDocumentBytes+1)); err == nil {
		t.Fatal("oversized document was not rejected")
	}
	deepJSON := []byte(strings.Repeat("[", MaxJSONDepth+1) + "0" + strings.Repeat("]", MaxJSONDepth+1))
	if _, err := DecodeProgram(deepJSON); err == nil || !strings.Contains(err.Error(), "nesting depth") {
		t.Fatalf("over-depth JSON was not rejected before decoding: %v", err)
	}

	deep := Expression{Op: OpBooleanEqual, Fact: "signature.present", Boolean: boolPointer(true)}
	for range MaxExpressionDepth {
		child := deep
		deep = Expression{Op: OpNot, Arg: &child}
	}
	program.Predicate = deep
	if err := ValidateProgram(program); err == nil {
		t.Fatal("over-depth predicate was accepted")
	}

	leaf := Expression{Op: OpBooleanEqual, Fact: "signature.present", Boolean: boolPointer(true)}
	children := make([]Expression, MaxCompositeArgs)
	for i := range children {
		children[i] = Expression{Op: OpAll, Args: []Expression{leaf, leaf, leaf, leaf}}
	}
	program.Predicate = Expression{Op: OpAll, Args: children}
	if err := ValidateProgram(program); err == nil {
		t.Fatal("over-node predicate was accepted")
	}

	program, _, _ = fixture()
	program.Predicate = Expression{Op: OpNumberEqual, Fact: "policy.score", Number: json.Number("1e999999")}
	if err := ValidateProgram(program); err == nil {
		t.Fatal("unbounded numeric exponent was accepted")
	}

	program, _, _ = fixture()
	program.Parameters = map[string]Parameter{"maximum": {Type: FactNumber, Number: json.Number("-1")}}
	program.Predicate = Expression{Op: OpTimeMapDeltaLessEqParameter, Fact: "starts", OtherFact: "ends", Parameter: "maximum"}
	if err := ValidateProgram(program); err == nil {
		t.Fatal("negative maximum time delta was accepted")
	}

	program, _, _ = fixture()
	program.SubjectID = strings.Repeat("x", MaxStringBytes+1)
	if err := ValidateProgram(program); err == nil {
		t.Fatal("oversized string was accepted")
	}

	program, evidence, _ = fixture()
	for i := 0; i <= MaxFacts; i++ {
		evidence.Facts["extra."+string(rune(0x1000+i))] = Fact{Type: FactBoolean, Complete: true, Boolean: boolPointer(true)}
	}
	if err := ValidateEvidence(evidence); err == nil {
		t.Fatal("oversized fact map was accepted")
	}

	_, evidence, _ = fixture()
	evidence.Facts = map[string]Fact{"graph": {Type: FactDirectedGraph, Complete: true, Edges: []DirectedEdge{
		{From: "z", To: "a"}, {From: "a", To: "z"},
	}}}
	if err := ValidateEvidence(evidence); err == nil {
		t.Fatal("non-canonical directed edges were accepted")
	}
	evidence.Facts = map[string]Fact{"flag": {Type: FactBoolean, Complete: true, Boolean: boolPointer(true), Edges: []DirectedEdge{}}}
	if err := ValidateEvidence(evidence); err == nil {
		t.Fatal("directed edges on a non-graph fact were accepted")
	}

	program, _, _ = fixture()
	program.Predicate.Op = Operation("shell")
	if result := Evaluate(program, evidence, time.Date(2026, time.August, 29, 12, 0, 0, 0, time.UTC)); result.Outcome != OutcomeBlocked || result.ReasonCode != ReasonInvalidProgram {
		t.Fatalf("unknown operation was not blocked: %#v", result)
	}
}

func cloneFacts(source map[string]Fact) map[string]Fact {
	result := make(map[string]Fact, len(source))
	for key, fact := range source {
		fact.Strings = append([]string(nil), fact.Strings...)
		if fact.Values != nil {
			fact.Values = make(map[string]string, len(fact.Values))
			for mapKey, value := range source[key].Values {
				fact.Values[mapKey] = value
			}
		}
		if fact.Booleans != nil {
			fact.Booleans = make(map[string]bool, len(fact.Booleans))
			for mapKey, value := range source[key].Booleans {
				fact.Booleans[mapKey] = value
			}
		}
		if fact.Numbers != nil {
			fact.Numbers = make(map[string]json.Number, len(fact.Numbers))
			for mapKey, value := range source[key].Numbers {
				fact.Numbers[mapKey] = value
			}
		}
		if fact.Timestamps != nil {
			fact.Timestamps = make(map[string]string, len(fact.Timestamps))
			for mapKey, value := range source[key].Timestamps {
				fact.Timestamps[mapKey] = value
			}
		}
		fact.Edges = append([]DirectedEdge(nil), fact.Edges...)
		fact.ConflictDigests = append([]string(nil), fact.ConflictDigests...)
		result[key] = fact
	}
	return result
}

func cloneEvidence(evidence Evidence) Evidence {
	evidence.ObservedSubjects = append([]string(nil), evidence.ObservedSubjects...)
	evidence.ContradictionDigests = append([]string(nil), evidence.ContradictionDigests...)
	evidence.Facts = cloneFacts(evidence.Facts)
	return evidence
}
