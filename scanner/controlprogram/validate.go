package controlprogram

import (
	"encoding/json"
	"fmt"
	"math/big"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

var (
	digestPattern  = regexp.MustCompile(`^[0-9a-f]{64}$`)
	controlPattern = regexp.MustCompile(`^(?:PRC-[0-9]{2}-[0-9]{3}|USEQ-[A-F0-9]{8})$`)
	factKeyPattern = regexp.MustCompile(`^[A-Za-z0-9_.-]+$`)
	numberPattern  = regexp.MustCompile(`^-?(?:0|[1-9][0-9]{0,63})(?:\.[0-9]{1,64})?(?:[eE]([+-]?[0-9]{1,3}))?$`)
)

var authorities = map[Authority]bool{
	AuthorityRepository: true, AuthorityArtifact: true, AuthorityExecuted: true,
	AuthorityEnvironment: true, AuthorityExternalRegistry: true, AuthorityStructuredRecord: true,
}

func validDigest(value string) bool { return digestPattern.MatchString(value) }

func validBoundedString(value string, allowEmpty bool) bool {
	return utf8.ValidString(value) && len(value) <= MaxStringBytes && (allowEmpty || value != "")
}

func sortedUnique(values []string, maximum int) bool {
	if len(values) > maximum {
		return false
	}
	for i, value := range values {
		if !validBoundedString(value, false) || (i > 0 && values[i-1] >= value) {
			return false
		}
	}
	return true
}

func validDigestList(values []string) bool {
	if !sortedUnique(values, MaxListValues) {
		return false
	}
	for _, value := range values {
		if !validDigest(value) {
			return false
		}
	}
	return true
}

func validFactKey(key string) bool {
	return len(key) <= MaxFactKeyBytes && factKeyPattern.MatchString(key)
}

// ValidateProgram applies every language and resource bound before evaluation.
func ValidateProgram(program Program) error {
	if program.SchemaVersion != ProgramSchemaVersion || !controlPattern.MatchString(program.ControlID) {
		return fmt.Errorf("invalid program schema or control ID")
	}
	if program.ControlRevision < 1 || !validDigest(program.ControlSemanticSHA256) ||
		!validDigest(program.ClauseID) || !validDigest(program.ClauseSHA256) ||
		!validDigest(program.ImplementationContractSHA256) || !validDigest(program.InventorySHA256) ||
		!validDigest(program.ApplicabilityProofContractSHA256) {
		return fmt.Errorf("invalid program identity digest")
	}
	if !validBoundedString(program.SubjectID, false) || !sortedUnique(program.Subjects, MaxSubjects) ||
		len(program.Subjects) == 0 || !authorities[program.RequiredAuthority] {
		return fmt.Errorf("invalid program subject inventory or authority")
	}
	if program.MaximumEvidenceAgeSeconds < 1 || program.MaximumEvidenceAgeSeconds > MaxEvidenceAgeSecs {
		return fmt.Errorf("invalid maximum evidence age")
	}
	if program.Parameters == nil || len(program.Parameters) > MaxFacts {
		return fmt.Errorf("invalid sealed parameter map")
	}
	for key, parameter := range program.Parameters {
		if !validFactKey(key) {
			return fmt.Errorf("invalid sealed parameter key")
		}
		if err := validateParameter(parameter); err != nil {
			return fmt.Errorf("parameter %s: %w", key, err)
		}
	}
	nodes := 0
	if err := validateExpression(program.Predicate, program.Parameters, 1, &nodes); err != nil {
		return err
	}
	return nil
}

func validateExpression(expression Expression, parameters map[string]Parameter, depth int, nodes *int) error {
	(*nodes)++
	if depth > MaxExpressionDepth || *nodes > MaxExpressionNodes {
		return fmt.Errorf("predicate exceeds expression bounds")
	}
	noLeafValues := func() bool {
		return expression.Fact == "" && expression.OtherFact == "" && expression.ThirdFact == "" && expression.Parameter == "" && expression.String == nil && expression.Boolean == nil &&
			expression.Number == "" && expression.Strings == nil && expression.Timestamp == nil
	}
	noCompositeValues := func() bool { return expression.Args == nil && expression.Arg == nil }
	singleFact := func() bool {
		return validFactKey(expression.Fact) && expression.OtherFact == "" && expression.ThirdFact == "" && expression.Parameter == "" && noCompositeValues()
	}
	twoFacts := func() bool {
		return validFactKey(expression.Fact) && validFactKey(expression.OtherFact) &&
			expression.ThirdFact == "" && expression.Parameter == "" && noCompositeValues() && expression.String == nil && expression.Boolean == nil &&
			expression.Number == "" && expression.Strings == nil && expression.Timestamp == nil
	}
	threeFacts := func() bool {
		return validFactKey(expression.Fact) && validFactKey(expression.OtherFact) && validFactKey(expression.ThirdFact) &&
			expression.Parameter == "" && noCompositeValues() && expression.String == nil && expression.Boolean == nil &&
			expression.Number == "" && expression.Strings == nil && expression.Timestamp == nil
	}
	factAndParameter := func(expectedFactType, expectedParameterType FactType) bool {
		parameter, ok := parameters[expression.Parameter]
		return validFactKey(expression.Fact) && expression.OtherFact == "" && expression.ThirdFact == "" && validFactKey(expression.Parameter) && ok &&
			parameter.Type == expectedParameterType && expectedFactType != "" && noCompositeValues() &&
			expression.String == nil && expression.Boolean == nil && expression.Number == "" &&
			expression.Strings == nil && expression.Timestamp == nil
	}
	factAndMapKeyParameter := func() bool {
		parameter, ok := parameters[expression.Parameter]
		return validFactKey(expression.Fact) && expression.OtherFact == "" && expression.ThirdFact == "" && validFactKey(expression.Parameter) && ok &&
			parameter.Type == FactStringSet && noCompositeValues() && expression.String == nil && expression.Boolean == nil &&
			expression.Number == "" && expression.Strings == nil && expression.Timestamp == nil
	}
	twoFactsAndParameter := func(expectedParameterType FactType) bool {
		parameter, ok := parameters[expression.Parameter]
		return validFactKey(expression.Fact) && validFactKey(expression.OtherFact) && expression.ThirdFact == "" && validFactKey(expression.Parameter) && ok &&
			parameter.Type == expectedParameterType && noCompositeValues() && expression.String == nil && expression.Boolean == nil &&
			expression.Number == "" && expression.Strings == nil && expression.Timestamp == nil
	}
	switch expression.Op {
	case OpAll, OpAny:
		if len(expression.Args) == 0 || len(expression.Args) > MaxCompositeArgs || expression.Arg != nil || !noLeafValues() {
			return fmt.Errorf("%s has invalid operands", expression.Op)
		}
		for _, child := range expression.Args {
			if err := validateExpression(child, parameters, depth+1, nodes); err != nil {
				return err
			}
		}
	case OpNot:
		if expression.Arg == nil || expression.Args != nil || !noLeafValues() {
			return fmt.Errorf("not has invalid operands")
		}
		return validateExpression(*expression.Arg, parameters, depth+1, nodes)
	case OpIdentityEqual, OpSchemaEqual, OpStringEqual:
		if !singleFact() || expression.String == nil || !validBoundedString(*expression.String, expression.Op == OpStringEqual) ||
			expression.Boolean != nil || expression.Number != "" || expression.Strings != nil || expression.Timestamp != nil {
			return fmt.Errorf("%s has invalid operands", expression.Op)
		}
	case OpDigestEqual:
		if !singleFact() || expression.String == nil || !validDigest(*expression.String) || expression.Boolean != nil ||
			expression.Number != "" || expression.Strings != nil || expression.Timestamp != nil {
			return fmt.Errorf("digest_eq has invalid operands")
		}
	case OpStateIn:
		if !singleFact() || expression.Strings == nil || !sortedUnique(expression.Strings, MaxListValues) || len(expression.Strings) == 0 ||
			expression.String != nil || expression.Boolean != nil || expression.Number != "" || expression.Timestamp != nil {
			return fmt.Errorf("state_in has invalid operands")
		}
	case OpBooleanEqual:
		if !singleFact() || expression.Boolean == nil || expression.String != nil || expression.Number != "" ||
			expression.Strings != nil || expression.Timestamp != nil {
			return fmt.Errorf("boolean_eq has invalid operands")
		}
	case OpBooleanMapAllEqual:
		if !singleFact() || expression.Boolean == nil || expression.String != nil || expression.Number != "" ||
			expression.Strings != nil || expression.Timestamp != nil {
			return fmt.Errorf("boolean_map_all_eq has invalid operands")
		}
	case OpDirectedGraphAcyclic:
		if !singleFact() || expression.String != nil || expression.Boolean != nil || expression.Number != "" ||
			expression.Strings != nil || expression.Timestamp != nil {
			return fmt.Errorf("directed_graph_acyclic has invalid operands")
		}
	case OpStringMapAllNonempty:
		if !singleFact() || expression.String != nil || expression.Boolean != nil || expression.Number != "" ||
			expression.Strings != nil || expression.Timestamp != nil {
			return fmt.Errorf("string_map_all_nonempty has invalid operands")
		}
	case OpIdentityMapValuesUnique:
		if !singleFact() || expression.String != nil || expression.Boolean != nil || expression.Number != "" ||
			expression.Strings != nil || expression.Timestamp != nil {
			return fmt.Errorf("identity_map_values_unique has invalid operands")
		}
	case OpNumberEqual, OpNumberLess, OpNumberLessEq, OpNumberGreater, OpNumberGreaterEq:
		if !singleFact() || !validNumber(string(expression.Number)) || expression.String != nil || expression.Boolean != nil ||
			expression.Strings != nil || expression.Timestamp != nil {
			return fmt.Errorf("%s has invalid operands", expression.Op)
		}
	case OpTimeBefore, OpTimeBeforeEq, OpTimeAfter, OpTimeAfterEq:
		if !singleFact() || expression.Timestamp == nil || !validTimestamp(*expression.Timestamp) ||
			expression.String != nil || expression.Boolean != nil || expression.Number != "" || expression.Strings != nil {
			return fmt.Errorf("%s has invalid operands", expression.Op)
		}
	case OpSetEqual, OpSetContainsAll, OpSetDisjoint:
		if !singleFact() || expression.Strings == nil || !sortedUnique(expression.Strings, MaxListValues) || expression.String != nil ||
			expression.Boolean != nil || expression.Number != "" || expression.Timestamp != nil {
			return fmt.Errorf("%s has invalid operands", expression.Op)
		}
	case OpMapKeysEqualSet, OpIdentityMapValuesIn, OpStateMapValuesIn, OpStringMapValuesIn:
		if !singleFact() || expression.Strings == nil || len(expression.Strings) == 0 || !sortedUnique(expression.Strings, MaxListValues) ||
			expression.String != nil || expression.Boolean != nil || expression.Number != "" || expression.Timestamp != nil {
			return fmt.Errorf("%s has invalid operands", expression.Op)
		}
	case OpIdentityEqualFact, OpSchemaEqualFact, OpDigestEqualFact, OpStringEqualFact,
		OpStateInSetFact, OpBooleanEqualFact, OpNumberEqualFact, OpNumberLessFact,
		OpNumberLessEqFact, OpNumberGreaterFact, OpNumberGreaterEqFact,
		OpTimeBeforeFact, OpTimeBeforeEqFact, OpTimeAfterFact, OpTimeAfterEqFact,
		OpSetEqualFact, OpSetContainsAllFact, OpSetDisjointFact,
		OpIdentityMapEqualFact, OpIdentityMapValuesInFact, OpIdentityMapValuesNotInFact, OpIdentityMapValuesNotEqualFact, OpSchemaMapEqualFact, OpDigestMapEqualFact,
		OpStateMapEqualFact, OpStringMapEqualFact, OpBooleanMapEqualFact, OpMapKeysEqualSetFact,
		OpBooleanMapAnyTrueFact, OpBooleanMapImpliesFact, OpStringMapAnyNonemptyFact,
		OpNumberMapEqualFact, OpNumberMapLessFact, OpNumberMapLessEqFact,
		OpNumberMapGreaterFact, OpNumberMapGreaterEqFact,
		OpTimeMapEqualFact, OpTimeMapBeforeFact, OpTimeMapBeforeEqFact, OpTimeMapAfterFact, OpTimeMapAfterEqFact:
		if !twoFacts() {
			return fmt.Errorf("%s has invalid operands", expression.Op)
		}
	case OpIdentityEqualParameter:
		if !factAndParameter(FactIdentity, FactIdentity) {
			return fmt.Errorf("%s has invalid operands", expression.Op)
		}
	case OpSchemaEqualParameter:
		if !factAndParameter(FactSchema, FactSchema) {
			return fmt.Errorf("%s has invalid operands", expression.Op)
		}
	case OpDigestEqualParameter:
		if !factAndParameter(FactDigest, FactDigest) {
			return fmt.Errorf("%s has invalid operands", expression.Op)
		}
	case OpStringEqualParameter:
		if !factAndParameter(FactString, FactString) {
			return fmt.Errorf("%s has invalid operands", expression.Op)
		}
	case OpStateInParameter:
		if !factAndParameter(FactState, FactStringSet) {
			return fmt.Errorf("%s has invalid operands", expression.Op)
		}
	case OpBooleanEqualParameter:
		if !factAndParameter(FactBoolean, FactBoolean) {
			return fmt.Errorf("%s has invalid operands", expression.Op)
		}
	case OpNumberEqualParameter, OpNumberLessParameter, OpNumberLessEqParameter,
		OpNumberGreaterParameter, OpNumberGreaterEqParameter:
		if !factAndParameter(FactNumber, FactNumber) {
			return fmt.Errorf("%s has invalid operands", expression.Op)
		}
	case OpTimeBeforeParameter, OpTimeBeforeEqParameter, OpTimeAfterParameter, OpTimeAfterEqParameter:
		if !factAndParameter(FactTime, FactTime) {
			return fmt.Errorf("%s has invalid operands", expression.Op)
		}
	case OpSetEqualParameter, OpSetContainsAllParameter, OpSetDisjointParameter:
		if !factAndParameter(FactStringSet, FactStringSet) {
			return fmt.Errorf("%s has invalid operands", expression.Op)
		}
	case OpIdentityMapEqualParameter:
		if !factAndParameter(FactIdentityMap, FactIdentityMap) {
			return fmt.Errorf("%s has invalid operands", expression.Op)
		}
	case OpIdentityMapValuesDifferForPairsParameter:
		parameter := parameters[expression.Parameter]
		if !factAndParameter(FactIdentityMap, FactDirectedGraph) || !validIdentityPairSet(parameter.Edges) {
			return fmt.Errorf("%s has invalid operands", expression.Op)
		}
	case OpSchemaMapEqualParameter:
		if !factAndParameter(FactSchemaMap, FactSchemaMap) {
			return fmt.Errorf("%s has invalid operands", expression.Op)
		}
	case OpDigestMapEqualParameter:
		if !factAndParameter(FactDigestMap, FactDigestMap) {
			return fmt.Errorf("%s has invalid operands", expression.Op)
		}
	case OpStateMapEqualParameter:
		if !factAndParameter(FactStateMap, FactStateMap) {
			return fmt.Errorf("%s has invalid operands", expression.Op)
		}
	case OpStringMapEqualParameter:
		if !factAndParameter(FactStringMap, FactStringMap) {
			return fmt.Errorf("%s has invalid operands", expression.Op)
		}
	case OpBooleanMapEqualParameter:
		if !factAndParameter(FactBooleanMap, FactBooleanMap) {
			return fmt.Errorf("%s has invalid operands", expression.Op)
		}
	case OpMapKeysEqualSetParameter:
		if !factAndMapKeyParameter() {
			return fmt.Errorf("%s has invalid operands", expression.Op)
		}
	case OpIdentityMapValuesInParameter:
		if !factAndParameter(FactIdentityMap, FactStringSet) {
			return fmt.Errorf("%s has invalid operands", expression.Op)
		}
	case OpIdentityMapValuesBijectSetParameter:
		parameter := parameters[expression.Parameter]
		if !factAndParameter(FactIdentityMap, FactStringSet) || len(parameter.Strings) == 0 {
			return fmt.Errorf("%s has invalid operands", expression.Op)
		}
	case OpStateMapValuesInParameter:
		if !factAndParameter(FactStateMap, FactStringSet) {
			return fmt.Errorf("%s has invalid operands", expression.Op)
		}
	case OpStringMapValuesInParameter:
		if !factAndParameter(FactStringMap, FactStringSet) {
			return fmt.Errorf("%s has invalid operands", expression.Op)
		}
	case OpBooleanMapAllEqualParameter:
		if !factAndParameter(FactBooleanMap, FactBoolean) {
			return fmt.Errorf("%s has invalid operands", expression.Op)
		}
	case OpTimeMapDeltaLessEqParameter:
		parameter := parameters[expression.Parameter]
		maximum, ok := new(big.Rat).SetString(parameter.Number.String())
		if !twoFactsAndParameter(FactNumber) || !ok || maximum.Sign() < 0 {
			return fmt.Errorf("%s has invalid operands", expression.Op)
		}
	case OpTimeMapDeltaEqualNumberMapFact:
		if !threeFacts() {
			return fmt.Errorf("%s has invalid operands", expression.Op)
		}
	case OpMapKeyPartitionEqualSetParameter:
		parameter := parameters[expression.Parameter]
		if !twoFactsAndParameter(FactStringSet) || parameter.Strings == nil {
			return fmt.Errorf("%s has invalid operands", expression.Op)
		}
	case OpNumberMapEqualParameter, OpNumberMapLessParameter, OpNumberMapLessEqParameter,
		OpNumberMapGreaterParameter, OpNumberMapGreaterEqParameter:
		if !factAndParameter(FactNumberMap, FactNumberMap) {
			return fmt.Errorf("%s has invalid operands", expression.Op)
		}
	case OpTimeMapEqualParameter, OpTimeMapBeforeParameter, OpTimeMapBeforeEqParameter, OpTimeMapAfterParameter, OpTimeMapAfterEqParameter:
		if !factAndParameter(FactTimeMap, FactTimeMap) {
			return fmt.Errorf("%s has invalid operands", expression.Op)
		}
	default:
		return fmt.Errorf("unsupported predicate operation %q", expression.Op)
	}
	return nil
}

func validateParameter(parameter Parameter) error {
	fact := Fact{
		Type: parameter.Type, Complete: true, String: parameter.String, Boolean: parameter.Boolean,
		Number: parameter.Number, Strings: parameter.Strings, Timestamp: parameter.Timestamp,
		Values: parameter.Values, Booleans: parameter.Booleans,
		Numbers: parameter.Numbers, Timestamps: parameter.Timestamps, Edges: parameter.Edges,
	}
	return validateFact(fact)
}

func validNumber(value string) bool {
	if value == "" || len(value) > 140 || strings.ContainsAny(value, "/") {
		return false
	}
	match := numberPattern.FindStringSubmatch(value)
	if match == nil {
		return false
	}
	if match[1] != "" {
		exponent, err := strconv.Atoi(match[1])
		if err != nil || exponent < -1000 || exponent > 1000 {
			return false
		}
	}
	_, ok := new(big.Rat).SetString(value)
	return ok
}

func validTimestamp(value string) bool {
	if !validBoundedString(value, false) {
		return false
	}
	_, err := time.Parse(time.RFC3339Nano, value)
	return err == nil
}

// ValidateEvidence validates shape and normalized values. Incomplete or
// conflicting facts remain structurally valid so Evaluate can return Blocked.
func ValidateEvidence(evidence Evidence) error {
	if evidence.SchemaVersion != EvidenceSchemaVersion || !validBoundedString(evidence.EvidenceID, false) ||
		!validDigest(evidence.ProgramSHA256) || !controlPattern.MatchString(evidence.ControlID) ||
		evidence.ControlRevision < 1 || !validDigest(evidence.ControlSemanticSHA256) ||
		!validDigest(evidence.ClauseID) || !validDigest(evidence.ClauseSHA256) ||
		!validDigest(evidence.ImplementationContractSHA256) || !validBoundedString(evidence.SubjectID, false) ||
		evidence.ObservedSubjects == nil || !sortedUnique(evidence.ObservedSubjects, MaxSubjects) || !validDigest(evidence.InventorySHA256) ||
		!authorities[evidence.Authority] || evidence.ObservedAt.IsZero() ||
		!validDigest(evidence.ApplicabilityProofContractSHA256) || evidence.Facts == nil || len(evidence.Facts) > MaxFacts {
		return fmt.Errorf("invalid evidence envelope")
	}
	if !validDigestList(evidence.ContradictionDigests) {
		return fmt.Errorf("invalid evidence contradiction digests")
	}
	switch evidence.Applicability {
	case ApplicabilityApplicable, ApplicabilityUnknown:
		if evidence.ApplicabilityProof != "" || evidence.ApplicabilityProofSHA256 != "" {
			return fmt.Errorf("unexpected applicability proof")
		}
	case ApplicabilityNotApplicable:
		if !validBoundedString(evidence.ApplicabilityProof, false) || !validDigest(evidence.ApplicabilityProofSHA256) {
			return fmt.Errorf("invalid applicability proof")
		}
	default:
		return fmt.Errorf("invalid applicability")
	}
	keys := make([]string, 0, len(evidence.Facts))
	for key := range evidence.Facts {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		if !validFactKey(key) {
			return fmt.Errorf("invalid fact key")
		}
		if err := validateFact(evidence.Facts[key]); err != nil {
			return fmt.Errorf("fact %s: %w", key, err)
		}
	}
	return nil
}

func validateFact(fact Fact) error {
	if !validDigestList(fact.ConflictDigests) {
		return fmt.Errorf("invalid conflict digests")
	}
	noString := fact.String == nil
	noBoolean := fact.Boolean == nil
	noNumber := fact.Number == ""
	noStrings := fact.Strings == nil
	noTime := fact.Timestamp == nil
	noValues := fact.Values == nil
	noBooleans := fact.Booleans == nil
	noNumbers := fact.Numbers == nil
	noTimestamps := fact.Timestamps == nil
	if fact.Type != FactDirectedGraph && fact.Edges != nil {
		return fmt.Errorf("unexpected directed graph edges")
	}
	switch fact.Type {
	case FactIdentity, FactSchema, FactState:
		if noString || !validBoundedString(*fact.String, false) || !noBoolean || !noNumber || !noStrings || !noTime || !noValues || !noBooleans || !noNumbers || !noTimestamps {
			return fmt.Errorf("invalid typed string value")
		}
	case FactString:
		if noString || !validBoundedString(*fact.String, true) || !noBoolean || !noNumber || !noStrings || !noTime || !noValues || !noBooleans || !noNumbers || !noTimestamps {
			return fmt.Errorf("invalid string value")
		}
	case FactDigest:
		if noString || !validDigest(*fact.String) || !noBoolean || !noNumber || !noStrings || !noTime || !noValues || !noBooleans || !noNumbers || !noTimestamps {
			return fmt.Errorf("invalid digest value")
		}
	case FactBoolean:
		if !noString || noBoolean || !noNumber || !noStrings || !noTime || !noValues || !noBooleans || !noNumbers || !noTimestamps {
			return fmt.Errorf("invalid boolean value")
		}
	case FactNumber:
		if !noString || !noBoolean || !validNumber(string(fact.Number)) || !noStrings || !noTime || !noValues || !noBooleans || !noNumbers || !noTimestamps {
			return fmt.Errorf("invalid number value")
		}
	case FactTime:
		if !noString || !noBoolean || !noNumber || !noStrings || fact.Timestamp == nil || !validTimestamp(*fact.Timestamp) || !noValues || !noBooleans || !noNumbers || !noTimestamps {
			return fmt.Errorf("invalid time value")
		}
	case FactStringSet:
		if !noString || !noBoolean || !noNumber || fact.Strings == nil || !sortedUnique(fact.Strings, MaxListValues) || !noTime || !noValues || !noBooleans || !noNumbers || !noTimestamps {
			return fmt.Errorf("invalid string set")
		}
	case FactIdentityMap, FactSchemaMap, FactStateMap, FactStringMap:
		allowEmpty := fact.Type == FactStringMap
		if !noString || !noBoolean || !noNumber || !noStrings || !noTime || !validStringMap(fact.Values, allowEmpty, false) || !noBooleans || !noNumbers || !noTimestamps {
			return fmt.Errorf("invalid %s", fact.Type)
		}
	case FactDigestMap:
		if !noString || !noBoolean || !noNumber || !noStrings || !noTime || !validStringMap(fact.Values, false, true) || !noBooleans || !noNumbers || !noTimestamps {
			return fmt.Errorf("invalid digest map")
		}
	case FactBooleanMap:
		if !noString || !noBoolean || !noNumber || !noStrings || !noTime || !noValues || fact.Booleans == nil || len(fact.Booleans) > MaxListValues || !noNumbers || !noTimestamps {
			return fmt.Errorf("invalid boolean map")
		}
		for key := range fact.Booleans {
			if !validFactKey(key) {
				return fmt.Errorf("invalid boolean map")
			}
		}
	case FactNumberMap:
		if !noString || !noBoolean || !noNumber || !noStrings || !noTime || !noValues || !noBooleans || !validNumberMap(fact.Numbers) || !noTimestamps {
			return fmt.Errorf("invalid number map")
		}
	case FactTimeMap:
		if !noString || !noBoolean || !noNumber || !noStrings || !noTime || !noValues || !noBooleans || !noNumbers || !validTimeMap(fact.Timestamps) {
			return fmt.Errorf("invalid time map")
		}
	case FactDirectedGraph:
		if !noString || !noBoolean || !noNumber || !noStrings || !noTime || !noValues || !noBooleans || !noNumbers || !noTimestamps ||
			fact.Edges == nil || len(fact.Edges) > MaxListValues || !validDirectedEdges(fact.Edges) {
			return fmt.Errorf("invalid directed graph")
		}
	default:
		return fmt.Errorf("unsupported fact type")
	}
	return nil
}

func validDirectedEdges(edges []DirectedEdge) bool {
	for index, edge := range edges {
		if !validBoundedString(edge.From, false) || !validBoundedString(edge.To, false) {
			return false
		}
		if index > 0 {
			previous := edges[index-1]
			if previous.From > edge.From || (previous.From == edge.From && previous.To >= edge.To) {
				return false
			}
		}
	}
	return true
}

func validIdentityPairSet(pairs []DirectedEdge) bool {
	if len(pairs) == 0 || !validDirectedEdges(pairs) {
		return false
	}
	for _, pair := range pairs {
		if pair.From >= pair.To {
			return false
		}
	}
	return true
}

func validStringMap(values map[string]string, allowEmpty bool, digestValues bool) bool {
	if values == nil || len(values) > MaxListValues {
		return false
	}
	for key, value := range values {
		if !validFactKey(key) {
			return false
		}
		if digestValues {
			if !validDigest(value) {
				return false
			}
		} else if !validBoundedString(value, allowEmpty) {
			return false
		}
	}
	return true
}

func validNumberMap(values map[string]json.Number) bool {
	if values == nil || len(values) > MaxListValues {
		return false
	}
	for key, value := range values {
		if !validFactKey(key) || !validNumber(value.String()) {
			return false
		}
	}
	return true
}

func validTimeMap(values map[string]string) bool {
	if values == nil || len(values) > MaxListValues {
		return false
	}
	for key, value := range values {
		if !validFactKey(key) || !validTimestamp(value) {
			return false
		}
	}
	return true
}
