// Package controlprogram evaluates bounded declarative predicates over
// normalized authoritative facts. It does not execute code or obtain evidence.
package controlprogram

import (
	"encoding/json"
	"time"
)

const (
	ProgramSchemaVersion  = "prc.control-check-program/v0.1"
	EvidenceSchemaVersion = "prc.control-check-evidence/v0.1"

	MaxDocumentBytes   = 1 << 20
	MaxExpressionDepth = 16
	MaxExpressionNodes = 256
	MaxJSONDepth       = 32
	MaxCompositeArgs   = 64
	MaxFacts           = 512
	MaxSubjects        = 1024
	MaxListValues      = 1024
	MaxStringBytes     = 4096
	MaxFactKeyBytes    = 128
	MaxEvidenceAgeSecs = 31_536_000
)

// Authority is the already-authenticated evidence boundary.
type Authority string

const (
	AuthorityRepository       Authority = "repository"
	AuthorityArtifact         Authority = "artifact"
	AuthorityExecuted         Authority = "executed"
	AuthorityEnvironment      Authority = "environment"
	AuthorityExternalRegistry Authority = "external_registry"
	AuthorityStructuredRecord Authority = "structured_record"
)

// Applicability is evidence-bound. Unknown never becomes Not Applicable.
type Applicability string

const (
	ApplicabilityApplicable    Applicability = "applicable"
	ApplicabilityNotApplicable Applicability = "not_applicable"
	ApplicabilityUnknown       Applicability = "unknown"
)

// Outcome is the sealed predicate result.
type Outcome string

const (
	OutcomePass          Outcome = "pass"
	OutcomeFail          Outcome = "fail"
	OutcomeBlocked       Outcome = "blocked"
	OutcomeNotApplicable Outcome = "not_applicable"
)

// ReasonCode gives stable fail-closed provenance.
type ReasonCode string

const (
	ReasonPassed                    ReasonCode = "passed"
	ReasonPredicateFalse            ReasonCode = "predicate_false"
	ReasonNotApplicable             ReasonCode = "not_applicable"
	ReasonInvalidProgram            ReasonCode = "invalid_program"
	ReasonInvalidEvidence           ReasonCode = "invalid_evidence"
	ReasonEvidenceBindingMismatch   ReasonCode = "evidence_binding_mismatch"
	ReasonWrongAuthority            ReasonCode = "wrong_authority"
	ReasonEvidenceStale             ReasonCode = "evidence_stale"
	ReasonEvidenceFromFuture        ReasonCode = "evidence_from_future"
	ReasonEvidenceIncomplete        ReasonCode = "evidence_incomplete"
	ReasonEvidenceConflicting       ReasonCode = "evidence_conflicting"
	ReasonApplicabilityUnknown      ReasonCode = "applicability_unknown"
	ReasonApplicabilityDisallowed   ReasonCode = "not_applicable_disallowed"
	ReasonApplicabilityProofMissing ReasonCode = "applicability_proof_missing"
	ReasonFactMissing               ReasonCode = "fact_missing"
	ReasonFactTypeMismatch          ReasonCode = "fact_type_mismatch"
)

// FactType prevents a generic string from being used as an identity, digest,
// schema, state, or timestamp without an explicit reviewed type.
type FactType string

const (
	FactIdentity      FactType = "identity"
	FactSchema        FactType = "schema"
	FactDigest        FactType = "digest"
	FactState         FactType = "state"
	FactString        FactType = "string"
	FactBoolean       FactType = "boolean"
	FactNumber        FactType = "number"
	FactTime          FactType = "time"
	FactStringSet     FactType = "string_set"
	FactIdentityMap   FactType = "identity_map"
	FactSchemaMap     FactType = "schema_map"
	FactDigestMap     FactType = "digest_map"
	FactStateMap      FactType = "state_map"
	FactStringMap     FactType = "string_map"
	FactBooleanMap    FactType = "boolean_map"
	FactNumberMap     FactType = "number_map"
	FactTimeMap       FactType = "time_map"
	FactDirectedGraph FactType = "directed_graph"
)

// Operation is a closed predicate language. There are no functions, regular
// expressions, interpolation, traversal, or extension hooks.
type Operation string

const (
	OpAll                                      Operation = "all"
	OpAny                                      Operation = "any"
	OpNot                                      Operation = "not"
	OpIdentityEqual                            Operation = "identity_eq"
	OpSchemaEqual                              Operation = "schema_eq"
	OpDigestEqual                              Operation = "digest_eq"
	OpStringEqual                              Operation = "string_eq"
	OpStateIn                                  Operation = "state_in"
	OpBooleanEqual                             Operation = "boolean_eq"
	OpNumberEqual                              Operation = "number_eq"
	OpNumberLess                               Operation = "number_lt"
	OpNumberLessEq                             Operation = "number_lte"
	OpNumberGreater                            Operation = "number_gt"
	OpNumberGreaterEq                          Operation = "number_gte"
	OpTimeBefore                               Operation = "time_before"
	OpTimeBeforeEq                             Operation = "time_lte"
	OpTimeAfter                                Operation = "time_after"
	OpTimeAfterEq                              Operation = "time_gte"
	OpSetEqual                                 Operation = "set_eq"
	OpSetContainsAll                           Operation = "set_contains_all"
	OpSetDisjoint                              Operation = "set_disjoint"
	OpMapKeysEqualSet                          Operation = "map_keys_eq_set"
	OpIdentityMapValuesIn                      Operation = "identity_map_values_in"
	OpStateMapValuesIn                         Operation = "state_map_values_in"
	OpStringMapValuesIn                        Operation = "string_map_values_in"
	OpBooleanMapAllEqual                       Operation = "boolean_map_all_eq"
	OpStringMapAllNonempty                     Operation = "string_map_all_nonempty"
	OpIdentityMapValuesUnique                  Operation = "identity_map_values_unique"
	OpDirectedGraphAcyclic                     Operation = "directed_graph_acyclic"
	OpIdentityEqualFact                        Operation = "identity_eq_fact"
	OpSchemaEqualFact                          Operation = "schema_eq_fact"
	OpDigestEqualFact                          Operation = "digest_eq_fact"
	OpStringEqualFact                          Operation = "string_eq_fact"
	OpStateInSetFact                           Operation = "state_in_set_fact"
	OpBooleanEqualFact                         Operation = "boolean_eq_fact"
	OpNumberEqualFact                          Operation = "number_eq_fact"
	OpNumberLessFact                           Operation = "number_lt_fact"
	OpNumberLessEqFact                         Operation = "number_lte_fact"
	OpNumberGreaterFact                        Operation = "number_gt_fact"
	OpNumberGreaterEqFact                      Operation = "number_gte_fact"
	OpTimeBeforeFact                           Operation = "time_before_fact"
	OpTimeBeforeEqFact                         Operation = "time_lte_fact"
	OpTimeAfterFact                            Operation = "time_after_fact"
	OpTimeAfterEqFact                          Operation = "time_gte_fact"
	OpSetEqualFact                             Operation = "set_eq_fact"
	OpSetContainsAllFact                       Operation = "set_contains_all_fact"
	OpSetDisjointFact                          Operation = "set_disjoint_fact"
	OpIdentityMapEqualFact                     Operation = "identity_map_eq_fact"
	OpIdentityMapValuesInFact                  Operation = "identity_map_values_in_fact"
	OpIdentityMapValuesNotInFact               Operation = "identity_map_values_not_in_fact"
	OpIdentityMapValuesNotEqualFact            Operation = "identity_map_values_not_equal_fact"
	OpSchemaMapEqualFact                       Operation = "schema_map_eq_fact"
	OpDigestMapEqualFact                       Operation = "digest_map_eq_fact"
	OpStateMapEqualFact                        Operation = "state_map_eq_fact"
	OpStringMapEqualFact                       Operation = "string_map_eq_fact"
	OpBooleanMapEqualFact                      Operation = "boolean_map_eq_fact"
	OpMapKeysEqualSetFact                      Operation = "map_keys_eq_set_fact"
	OpBooleanMapAnyTrueFact                    Operation = "boolean_map_any_true_fact"
	OpBooleanMapImpliesFact                    Operation = "boolean_map_implies_fact"
	OpStringMapAnyNonemptyFact                 Operation = "string_map_any_nonempty_fact"
	OpNumberMapEqualFact                       Operation = "number_map_eq_fact"
	OpNumberMapLessFact                        Operation = "number_map_lt_fact"
	OpNumberMapLessEqFact                      Operation = "number_map_lte_fact"
	OpNumberMapGreaterFact                     Operation = "number_map_gt_fact"
	OpNumberMapGreaterEqFact                   Operation = "number_map_gte_fact"
	OpTimeMapEqualFact                         Operation = "time_map_eq_fact"
	OpTimeMapBeforeFact                        Operation = "time_map_before_fact"
	OpTimeMapBeforeEqFact                      Operation = "time_map_lte_fact"
	OpTimeMapAfterFact                         Operation = "time_map_after_fact"
	OpTimeMapAfterEqFact                       Operation = "time_map_gte_fact"
	OpIdentityEqualParameter                   Operation = "identity_eq_parameter"
	OpSchemaEqualParameter                     Operation = "schema_eq_parameter"
	OpDigestEqualParameter                     Operation = "digest_eq_parameter"
	OpStringEqualParameter                     Operation = "string_eq_parameter"
	OpStateInParameter                         Operation = "state_in_parameter"
	OpBooleanEqualParameter                    Operation = "boolean_eq_parameter"
	OpNumberEqualParameter                     Operation = "number_eq_parameter"
	OpNumberLessParameter                      Operation = "number_lt_parameter"
	OpNumberLessEqParameter                    Operation = "number_lte_parameter"
	OpNumberGreaterParameter                   Operation = "number_gt_parameter"
	OpNumberGreaterEqParameter                 Operation = "number_gte_parameter"
	OpTimeBeforeParameter                      Operation = "time_before_parameter"
	OpTimeBeforeEqParameter                    Operation = "time_lte_parameter"
	OpTimeAfterParameter                       Operation = "time_after_parameter"
	OpTimeAfterEqParameter                     Operation = "time_gte_parameter"
	OpSetEqualParameter                        Operation = "set_eq_parameter"
	OpSetContainsAllParameter                  Operation = "set_contains_all_parameter"
	OpSetDisjointParameter                     Operation = "set_disjoint_parameter"
	OpIdentityMapEqualParameter                Operation = "identity_map_eq_parameter"
	OpIdentityMapValuesDifferForPairsParameter Operation = "identity_map_values_differ_for_pairs_parameter"
	OpSchemaMapEqualParameter                  Operation = "schema_map_eq_parameter"
	OpDigestMapEqualParameter                  Operation = "digest_map_eq_parameter"
	OpStateMapEqualParameter                   Operation = "state_map_eq_parameter"
	OpStringMapEqualParameter                  Operation = "string_map_eq_parameter"
	OpBooleanMapEqualParameter                 Operation = "boolean_map_eq_parameter"
	OpMapKeysEqualSetParameter                 Operation = "map_keys_eq_set_parameter"
	OpIdentityMapValuesInParameter             Operation = "identity_map_values_in_parameter"
	OpIdentityMapValuesBijectSetParameter      Operation = "identity_map_values_biject_set_parameter"
	OpStateMapValuesInParameter                Operation = "state_map_values_in_parameter"
	OpStringMapValuesInParameter               Operation = "string_map_values_in_parameter"
	OpBooleanMapAllEqualParameter              Operation = "boolean_map_all_eq_parameter"
	OpTimeMapDeltaLessEqParameter              Operation = "time_map_delta_lte_parameter"
	OpTimeMapDeltaEqualNumberMapFact           Operation = "time_map_delta_eq_number_map_fact"
	OpMapKeyPartitionEqualSetParameter         Operation = "map_key_partition_eq_set_parameter"
	OpNumberMapEqualParameter                  Operation = "number_map_eq_parameter"
	OpNumberMapLessParameter                   Operation = "number_map_lt_parameter"
	OpNumberMapLessEqParameter                 Operation = "number_map_lte_parameter"
	OpNumberMapGreaterParameter                Operation = "number_map_gt_parameter"
	OpNumberMapGreaterEqParameter              Operation = "number_map_gte_parameter"
	OpTimeMapEqualParameter                    Operation = "time_map_eq_parameter"
	OpTimeMapBeforeParameter                   Operation = "time_map_before_parameter"
	OpTimeMapBeforeEqParameter                 Operation = "time_map_lte_parameter"
	OpTimeMapAfterParameter                    Operation = "time_map_after_parameter"
	OpTimeMapAfterEqParameter                  Operation = "time_map_gte_parameter"
)

// Expression is a typed syntax tree. Validation requires exactly the fields
// needed by Op and rejects unused values.
type Expression struct {
	Op        Operation    `json:"op"`
	Args      []Expression `json:"args,omitempty"`
	Arg       *Expression  `json:"arg,omitempty"`
	Fact      string       `json:"fact,omitempty"`
	OtherFact string       `json:"other_fact,omitempty"`
	ThirdFact string       `json:"third_fact,omitempty"`
	Parameter string       `json:"parameter,omitempty"`
	String    *string      `json:"string,omitempty"`
	Boolean   *bool        `json:"boolean,omitempty"`
	Number    json.Number  `json:"number,omitempty"`
	Strings   []string     `json:"strings,omitempty"`
	Timestamp *string      `json:"timestamp,omitempty"`
}

// DirectedEdge is a canonical edge in a bounded directed graph fact.
type DirectedEdge struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// Program binds one reviewed deterministic clause to a closed predicate.
type Program struct {
	SchemaVersion                    string               `json:"schema_version"`
	ControlID                        string               `json:"control_id"`
	ControlRevision                  int                  `json:"control_revision"`
	ControlSemanticSHA256            string               `json:"control_semantic_sha256"`
	ClauseID                         string               `json:"clause_id"`
	ClauseSHA256                     string               `json:"clause_sha256"`
	ImplementationContractSHA256     string               `json:"implementation_contract_sha256"`
	SubjectID                        string               `json:"subject_id"`
	Subjects                         []string             `json:"subjects"`
	InventorySHA256                  string               `json:"inventory_sha256"`
	RequiredAuthority                Authority            `json:"required_authority"`
	AllowNotApplicable               bool                 `json:"allow_not_applicable"`
	ApplicabilityProofContractSHA256 string               `json:"applicability_proof_contract_sha256"`
	MaximumEvidenceAgeSeconds        int64                `json:"maximum_evidence_age_seconds"`
	Parameters                       map[string]Parameter `json:"parameters"`
	Predicate                        Expression           `json:"predicate"`
}

// Parameter is an expected value sealed into the Program before evidence is
// requested. This prevents an evidence provider from choosing both sides of a
// policy, threshold, identity, digest, time, or set comparison.
type Parameter struct {
	Type       FactType               `json:"type"`
	String     *string                `json:"string,omitempty"`
	Boolean    *bool                  `json:"boolean,omitempty"`
	Number     json.Number            `json:"number,omitempty"`
	Strings    []string               `json:"strings,omitempty"`
	Timestamp  *string                `json:"timestamp,omitempty"`
	Values     map[string]string      `json:"values,omitempty"`
	Booleans   map[string]bool        `json:"booleans,omitempty"`
	Numbers    map[string]json.Number `json:"numbers,omitempty"`
	Timestamps map[string]string      `json:"timestamps,omitempty"`
	Edges      []DirectedEdge         `json:"edges,omitempty"`
}

// Fact is one normalized value. Complete is evidence metadata, not the result
// of the predicate. Conflicts always block evaluation.
type Fact struct {
	Type            FactType               `json:"type"`
	Complete        bool                   `json:"complete"`
	String          *string                `json:"string,omitempty"`
	Boolean         *bool                  `json:"boolean,omitempty"`
	Number          json.Number            `json:"number,omitempty"`
	Strings         []string               `json:"strings,omitempty"`
	Timestamp       *string                `json:"timestamp,omitempty"`
	Values          map[string]string      `json:"values,omitempty"`
	Booleans        map[string]bool        `json:"booleans,omitempty"`
	Numbers         map[string]json.Number `json:"numbers,omitempty"`
	Timestamps      map[string]string      `json:"timestamps,omitempty"`
	Edges           []DirectedEdge         `json:"edges,omitempty"`
	ConflictDigests []string               `json:"conflict_digests,omitempty"`
}

// Evidence contains normalized facts from an authority established outside
// this package. The evaluator only verifies the immutable bindings and facts.
type Evidence struct {
	SchemaVersion                    string          `json:"schema_version"`
	EvidenceID                       string          `json:"evidence_id"`
	ProgramSHA256                    string          `json:"program_sha256"`
	ControlID                        string          `json:"control_id"`
	ControlRevision                  int             `json:"control_revision"`
	ControlSemanticSHA256            string          `json:"control_semantic_sha256"`
	ClauseID                         string          `json:"clause_id"`
	ClauseSHA256                     string          `json:"clause_sha256"`
	ImplementationContractSHA256     string          `json:"implementation_contract_sha256"`
	SubjectID                        string          `json:"subject_id"`
	ObservedSubjects                 []string        `json:"observed_subjects"`
	InventorySHA256                  string          `json:"inventory_sha256"`
	Authority                        Authority       `json:"authority"`
	ObservedAt                       time.Time       `json:"observed_at"`
	Complete                         bool            `json:"complete"`
	ContradictionDigests             []string        `json:"contradiction_digests,omitempty"`
	Applicability                    Applicability   `json:"applicability"`
	ApplicabilityProofContractSHA256 string          `json:"applicability_proof_contract_sha256"`
	ApplicabilityProof               string          `json:"applicability_proof,omitempty"`
	ApplicabilityProofSHA256         string          `json:"applicability_proof_sha256,omitempty"`
	Facts                            map[string]Fact `json:"facts"`
}

// Result binds a stable outcome to the exact program and evidence documents.
type Result struct {
	ControlID                    string     `json:"control_id"`
	ControlRevision              int        `json:"control_revision"`
	ControlSemanticSHA256        string     `json:"control_semantic_sha256"`
	ClauseID                     string     `json:"clause_id"`
	ClauseSHA256                 string     `json:"clause_sha256"`
	ImplementationContractSHA256 string     `json:"implementation_contract_sha256"`
	ProgramSHA256                string     `json:"program_sha256"`
	EvidenceSHA256               string     `json:"evidence_sha256"`
	EvaluatedAt                  time.Time  `json:"evaluated_at"`
	Outcome                      Outcome    `json:"outcome"`
	ReasonCode                   ReasonCode `json:"reason_code"`
}
