package controlprogram

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"maps"
	"math/big"
	"slices"
	"sort"
	"time"
)

// ProgramSHA256 returns the stable digest evidence must bind.
func ProgramSHA256(program Program) string { return digestJSON(program) }

// EvidenceSHA256 returns the stable digest included in Result.
func EvidenceSHA256(evidence Evidence) string { return digestJSON(evidence) }

func digestJSON(value any) string {
	data, err := json.Marshal(value)
	if err != nil {
		return ""
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func digestString(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}

type truth struct {
	known  bool
	value  bool
	reason ReasonCode
}

// Evaluate is pure and side-effect free. Now is explicit so the same sealed
// inputs always produce the same output.
func Evaluate(program Program, evidence Evidence, now time.Time) Result {
	result := Result{
		ControlID: program.ControlID, ControlRevision: program.ControlRevision,
		ControlSemanticSHA256: program.ControlSemanticSHA256,
		ClauseID:              program.ClauseID, ClauseSHA256: program.ClauseSHA256,
		ImplementationContractSHA256: program.ImplementationContractSHA256,
		ProgramSHA256:                ProgramSHA256(program), EvidenceSHA256: EvidenceSHA256(evidence),
		EvaluatedAt: now.UTC(), Outcome: OutcomeBlocked,
	}
	if err := ValidateProgram(program); err != nil || now.IsZero() {
		result.ReasonCode = ReasonInvalidProgram
		return result
	}
	if err := ValidateEvidence(evidence); err != nil {
		result.ReasonCode = ReasonInvalidEvidence
		return result
	}
	if evidence.ProgramSHA256 != result.ProgramSHA256 || evidence.ControlID != program.ControlID ||
		evidence.ControlRevision != program.ControlRevision ||
		evidence.ControlSemanticSHA256 != program.ControlSemanticSHA256 ||
		evidence.ClauseID != program.ClauseID || evidence.ClauseSHA256 != program.ClauseSHA256 ||
		evidence.ImplementationContractSHA256 != program.ImplementationContractSHA256 ||
		evidence.SubjectID != program.SubjectID || evidence.InventorySHA256 != program.InventorySHA256 ||
		evidence.ApplicabilityProofContractSHA256 != program.ApplicabilityProofContractSHA256 {
		result.ReasonCode = ReasonEvidenceBindingMismatch
		return result
	}
	if evidence.Authority != program.RequiredAuthority {
		result.ReasonCode = ReasonWrongAuthority
		return result
	}
	if len(evidence.ContradictionDigests) != 0 {
		result.ReasonCode = ReasonEvidenceConflicting
		return result
	}
	if !evidence.Complete || !slices.Equal(evidence.ObservedSubjects, program.Subjects) {
		result.ReasonCode = ReasonEvidenceIncomplete
		return result
	}
	factKeys := make([]string, 0, len(evidence.Facts))
	for key := range evidence.Facts {
		factKeys = append(factKeys, key)
	}
	sort.Strings(factKeys)
	for _, key := range factKeys {
		fact := evidence.Facts[key]
		if len(fact.ConflictDigests) != 0 {
			result.ReasonCode = ReasonEvidenceConflicting
			return result
		}
		if !fact.Complete {
			result.ReasonCode = ReasonEvidenceIncomplete
			return result
		}
	}
	if evidence.ObservedAt.After(now) {
		result.ReasonCode = ReasonEvidenceFromFuture
		return result
	}
	if now.Sub(evidence.ObservedAt) > time.Duration(program.MaximumEvidenceAgeSeconds)*time.Second {
		result.ReasonCode = ReasonEvidenceStale
		return result
	}
	switch evidence.Applicability {
	case ApplicabilityUnknown:
		result.ReasonCode = ReasonApplicabilityUnknown
		return result
	case ApplicabilityNotApplicable:
		if !program.AllowNotApplicable {
			result.ReasonCode = ReasonApplicabilityDisallowed
			return result
		}
		if evidence.ApplicabilityProof == "" ||
			evidence.ApplicabilityProofSHA256 != digestString(evidence.ApplicabilityProof) {
			result.ReasonCode = ReasonApplicabilityProofMissing
			return result
		}
		result.Outcome, result.ReasonCode = OutcomeNotApplicable, ReasonNotApplicable
		return result
	}

	evaluated := evaluateExpression(program.Predicate, program.Parameters, evidence.Facts)
	if !evaluated.known {
		result.ReasonCode = evaluated.reason
		return result
	}
	if evaluated.value {
		result.Outcome, result.ReasonCode = OutcomePass, ReasonPassed
	} else {
		// Every evidence, scope, authority, freshness, conflict, completeness,
		// and type gate above has passed. Only this path may emit Fail.
		result.Outcome, result.ReasonCode = OutcomeFail, ReasonPredicateFalse
	}
	return result
}

func evaluateExpression(expression Expression, parameters map[string]Parameter, facts map[string]Fact) truth {
	switch expression.Op {
	case OpAll, OpAny:
		values := make([]truth, len(expression.Args))
		for i, child := range expression.Args {
			values[i] = evaluateExpression(child, parameters, facts)
		}
		for _, value := range values {
			if !value.known {
				return value
			}
		}
		if expression.Op == OpAll {
			for _, value := range values {
				if !value.value {
					return truth{known: true}
				}
			}
			return truth{known: true, value: true}
		}
		for _, value := range values {
			if value.value {
				return truth{known: true, value: true}
			}
		}
		return truth{known: true}
	case OpNot:
		value := evaluateExpression(*expression.Arg, parameters, facts)
		if value.known {
			value.value = !value.value
		}
		return value
	}

	fact, ok := facts[expression.Fact]
	if !ok {
		return truth{reason: ReasonFactMissing}
	}
	if expression.Op == OpMapKeysEqualSetFact {
		other, found := facts[expression.OtherFact]
		if !found {
			return truth{reason: ReasonFactMissing}
		}
		keys, isMap := factMapKeys(fact)
		if !isMap || other.Type != FactStringSet {
			return truth{reason: ReasonFactTypeMismatch}
		}
		return truth{known: true, value: slices.Equal(keys, other.Strings)}
	}
	if expression.Op == OpMapKeysEqualSet {
		keys, isMap := factMapKeys(fact)
		if !isMap {
			return truth{reason: ReasonFactTypeMismatch}
		}
		return truth{known: true, value: slices.Equal(keys, expression.Strings)}
	}
	if expression.Op == OpTimeMapDeltaLessEqParameter {
		other, otherFound := facts[expression.OtherFact]
		parameter, parameterFound := parameters[expression.Parameter]
		if !otherFound {
			return truth{reason: ReasonFactMissing}
		}
		if !parameterFound {
			return truth{reason: ReasonInvalidProgram}
		}
		if fact.Type != FactTimeMap || other.Type != FactTimeMap || parameter.Type != FactNumber {
			return truth{reason: ReasonFactTypeMismatch}
		}
		return truth{known: true, value: timeMapDeltasLessEq(fact.Timestamps, other.Timestamps, parameter.Number)}
	}
	if expression.Op == OpTimeMapDeltaEqualNumberMapFact {
		ends, endsFound := facts[expression.OtherFact]
		durations, durationsFound := facts[expression.ThirdFact]
		if !endsFound || !durationsFound {
			return truth{reason: ReasonFactMissing}
		}
		if fact.Type != FactTimeMap || ends.Type != FactTimeMap || durations.Type != FactNumberMap {
			return truth{reason: ReasonFactTypeMismatch}
		}
		return truth{known: true, value: timeMapDeltasEqualNumberMap(fact.Timestamps, ends.Timestamps, durations.Numbers)}
	}
	if expression.Op == OpMapKeyPartitionEqualSetParameter {
		other, otherFound := facts[expression.OtherFact]
		parameter, parameterFound := parameters[expression.Parameter]
		if !otherFound {
			return truth{reason: ReasonFactMissing}
		}
		if !parameterFound {
			return truth{reason: ReasonInvalidProgram}
		}
		leftKeys, leftIsMap := factMapKeys(fact)
		rightKeys, rightIsMap := factMapKeys(other)
		if !leftIsMap || !rightIsMap || parameter.Type != FactStringSet {
			return truth{reason: ReasonFactTypeMismatch}
		}
		return truth{known: true, value: mapKeyPartitionEqualSet(leftKeys, rightKeys, parameter.Strings)}
	}
	factPairTypes := map[Operation][2]FactType{
		OpIdentityEqualFact:             {FactIdentity, FactIdentity},
		OpSchemaEqualFact:               {FactSchema, FactSchema},
		OpDigestEqualFact:               {FactDigest, FactDigest},
		OpStringEqualFact:               {FactString, FactString},
		OpStateInSetFact:                {FactState, FactStringSet},
		OpBooleanEqualFact:              {FactBoolean, FactBoolean},
		OpNumberEqualFact:               {FactNumber, FactNumber},
		OpNumberLessFact:                {FactNumber, FactNumber},
		OpNumberLessEqFact:              {FactNumber, FactNumber},
		OpNumberGreaterFact:             {FactNumber, FactNumber},
		OpNumberGreaterEqFact:           {FactNumber, FactNumber},
		OpTimeBeforeFact:                {FactTime, FactTime},
		OpTimeBeforeEqFact:              {FactTime, FactTime},
		OpTimeAfterFact:                 {FactTime, FactTime},
		OpTimeAfterEqFact:               {FactTime, FactTime},
		OpSetEqualFact:                  {FactStringSet, FactStringSet},
		OpSetContainsAllFact:            {FactStringSet, FactStringSet},
		OpSetDisjointFact:               {FactStringSet, FactStringSet},
		OpIdentityMapEqualFact:          {FactIdentityMap, FactIdentityMap},
		OpIdentityMapValuesInFact:       {FactIdentityMap, FactIdentityMap},
		OpIdentityMapValuesNotInFact:    {FactIdentityMap, FactIdentityMap},
		OpIdentityMapValuesNotEqualFact: {FactIdentityMap, FactIdentityMap},
		OpSchemaMapEqualFact:            {FactSchemaMap, FactSchemaMap},
		OpDigestMapEqualFact:            {FactDigestMap, FactDigestMap},
		OpStateMapEqualFact:             {FactStateMap, FactStateMap},
		OpStringMapEqualFact:            {FactStringMap, FactStringMap},
		OpBooleanMapEqualFact:           {FactBooleanMap, FactBooleanMap},
		OpBooleanMapAnyTrueFact:         {FactBooleanMap, FactBooleanMap},
		OpBooleanMapImpliesFact:         {FactBooleanMap, FactBooleanMap},
		OpStringMapAnyNonemptyFact:      {FactStringMap, FactStringMap},
		OpNumberMapEqualFact:            {FactNumberMap, FactNumberMap},
		OpNumberMapLessFact:             {FactNumberMap, FactNumberMap},
		OpNumberMapLessEqFact:           {FactNumberMap, FactNumberMap},
		OpNumberMapGreaterFact:          {FactNumberMap, FactNumberMap},
		OpNumberMapGreaterEqFact:        {FactNumberMap, FactNumberMap},
		OpTimeMapEqualFact:              {FactTimeMap, FactTimeMap},
		OpTimeMapBeforeFact:             {FactTimeMap, FactTimeMap},
		OpTimeMapBeforeEqFact:           {FactTimeMap, FactTimeMap},
		OpTimeMapAfterFact:              {FactTimeMap, FactTimeMap},
		OpTimeMapAfterEqFact:            {FactTimeMap, FactTimeMap},
	}
	if expectedPair, isPair := factPairTypes[expression.Op]; isPair {
		other, found := facts[expression.OtherFact]
		if !found {
			return truth{reason: ReasonFactMissing}
		}
		if fact.Type != expectedPair[0] || other.Type != expectedPair[1] {
			return truth{reason: ReasonFactTypeMismatch}
		}
		switch expression.Op {
		case OpIdentityEqualFact, OpSchemaEqualFact, OpDigestEqualFact, OpStringEqualFact:
			return truth{known: true, value: *fact.String == *other.String}
		case OpStateInSetFact:
			_, found := slices.BinarySearch(other.Strings, *fact.String)
			return truth{known: true, value: found}
		case OpBooleanEqualFact:
			return truth{known: true, value: *fact.Boolean == *other.Boolean}
		case OpNumberEqualFact, OpNumberLessFact, OpNumberLessEqFact, OpNumberGreaterFact, OpNumberGreaterEqFact:
			left, _ := new(big.Rat).SetString(fact.Number.String())
			right, _ := new(big.Rat).SetString(other.Number.String())
			comparison := left.Cmp(right)
			matches := map[Operation]bool{
				OpNumberEqualFact: comparison == 0, OpNumberLessFact: comparison < 0,
				OpNumberLessEqFact: comparison <= 0, OpNumberGreaterFact: comparison > 0,
				OpNumberGreaterEqFact: comparison >= 0,
			}[expression.Op]
			return truth{known: true, value: matches}
		case OpTimeBeforeFact, OpTimeBeforeEqFact, OpTimeAfterFact, OpTimeAfterEqFact:
			left, _ := time.Parse(time.RFC3339Nano, *fact.Timestamp)
			right, _ := time.Parse(time.RFC3339Nano, *other.Timestamp)
			comparison := left.Compare(right)
			matches := map[Operation]bool{
				OpTimeBeforeFact: comparison < 0, OpTimeBeforeEqFact: comparison <= 0,
				OpTimeAfterFact: comparison > 0, OpTimeAfterEqFact: comparison >= 0,
			}[expression.Op]
			return truth{known: true, value: matches}
		case OpSetEqualFact:
			return truth{known: true, value: slices.Equal(fact.Strings, other.Strings)}
		case OpSetContainsAllFact:
			return truth{known: true, value: containsAll(fact.Strings, other.Strings)}
		case OpSetDisjointFact:
			return truth{known: true, value: disjoint(fact.Strings, other.Strings)}
		case OpIdentityMapEqualFact, OpSchemaMapEqualFact, OpDigestMapEqualFact,
			OpStateMapEqualFact, OpStringMapEqualFact:
			return truth{known: true, value: maps.Equal(fact.Values, other.Values)}
		case OpIdentityMapValuesInFact:
			return truth{known: true, value: identityMapValuesIn(fact.Values, other.Values)}
		case OpIdentityMapValuesNotInFact:
			return truth{known: true, value: identityMapValuesNotIn(fact.Values, other.Values)}
		case OpIdentityMapValuesNotEqualFact:
			if len(fact.Values) != len(other.Values) {
				return truth{known: true}
			}
			for key, left := range fact.Values {
				right, found := other.Values[key]
				if !found || left == right {
					return truth{known: true}
				}
			}
			return truth{known: true, value: true}
		case OpBooleanMapEqualFact:
			return truth{known: true, value: maps.Equal(fact.Booleans, other.Booleans)}
		case OpBooleanMapAnyTrueFact:
			if len(fact.Booleans) != len(other.Booleans) {
				return truth{known: true}
			}
			for key, left := range fact.Booleans {
				right, found := other.Booleans[key]
				if !found || (!left && !right) {
					return truth{known: true}
				}
			}
			return truth{known: true, value: true}
		case OpBooleanMapImpliesFact:
			if len(fact.Booleans) != len(other.Booleans) {
				return truth{known: true}
			}
			for key, antecedent := range fact.Booleans {
				consequent, found := other.Booleans[key]
				if !found || (antecedent && !consequent) {
					return truth{known: true}
				}
			}
			return truth{known: true, value: true}
		case OpStringMapAnyNonemptyFact:
			if len(fact.Values) != len(other.Values) {
				return truth{known: true}
			}
			for key, left := range fact.Values {
				right, found := other.Values[key]
				if !found || (left == "" && right == "") {
					return truth{known: true}
				}
			}
			return truth{known: true, value: true}
		case OpNumberMapEqualFact, OpNumberMapLessFact, OpNumberMapLessEqFact,
			OpNumberMapGreaterFact, OpNumberMapGreaterEqFact:
			return truth{known: true, value: compareNumberMaps(fact.Numbers, other.Numbers, expression.Op)}
		case OpTimeMapEqualFact, OpTimeMapBeforeFact, OpTimeMapBeforeEqFact, OpTimeMapAfterFact, OpTimeMapAfterEqFact:
			return truth{known: true, value: compareTimeMaps(fact.Timestamps, other.Timestamps, expression.Op)}
		}
	}
	parameterTypes := map[Operation][2]FactType{
		OpIdentityEqualParameter:                   {FactIdentity, FactIdentity},
		OpSchemaEqualParameter:                     {FactSchema, FactSchema},
		OpDigestEqualParameter:                     {FactDigest, FactDigest},
		OpStringEqualParameter:                     {FactString, FactString},
		OpStateInParameter:                         {FactState, FactStringSet},
		OpBooleanEqualParameter:                    {FactBoolean, FactBoolean},
		OpNumberEqualParameter:                     {FactNumber, FactNumber},
		OpNumberLessParameter:                      {FactNumber, FactNumber},
		OpNumberLessEqParameter:                    {FactNumber, FactNumber},
		OpNumberGreaterParameter:                   {FactNumber, FactNumber},
		OpNumberGreaterEqParameter:                 {FactNumber, FactNumber},
		OpTimeBeforeParameter:                      {FactTime, FactTime},
		OpTimeBeforeEqParameter:                    {FactTime, FactTime},
		OpTimeAfterParameter:                       {FactTime, FactTime},
		OpTimeAfterEqParameter:                     {FactTime, FactTime},
		OpSetEqualParameter:                        {FactStringSet, FactStringSet},
		OpSetContainsAllParameter:                  {FactStringSet, FactStringSet},
		OpSetDisjointParameter:                     {FactStringSet, FactStringSet},
		OpIdentityMapEqualParameter:                {FactIdentityMap, FactIdentityMap},
		OpIdentityMapValuesDifferForPairsParameter: {FactIdentityMap, FactDirectedGraph},
		OpSchemaMapEqualParameter:                  {FactSchemaMap, FactSchemaMap},
		OpDigestMapEqualParameter:                  {FactDigestMap, FactDigestMap},
		OpStateMapEqualParameter:                   {FactStateMap, FactStateMap},
		OpStringMapEqualParameter:                  {FactStringMap, FactStringMap},
		OpBooleanMapEqualParameter:                 {FactBooleanMap, FactBooleanMap},
		OpIdentityMapValuesInParameter:             {FactIdentityMap, FactStringSet},
		OpIdentityMapValuesBijectSetParameter:      {FactIdentityMap, FactStringSet},
		OpStateMapValuesInParameter:                {FactStateMap, FactStringSet},
		OpStringMapValuesInParameter:               {FactStringMap, FactStringSet},
		OpBooleanMapAllEqualParameter:              {FactBooleanMap, FactBoolean},
		OpNumberMapEqualParameter:                  {FactNumberMap, FactNumberMap},
		OpNumberMapLessParameter:                   {FactNumberMap, FactNumberMap},
		OpNumberMapLessEqParameter:                 {FactNumberMap, FactNumberMap},
		OpNumberMapGreaterParameter:                {FactNumberMap, FactNumberMap},
		OpNumberMapGreaterEqParameter:              {FactNumberMap, FactNumberMap},
		OpTimeMapEqualParameter:                    {FactTimeMap, FactTimeMap},
		OpTimeMapBeforeParameter:                   {FactTimeMap, FactTimeMap},
		OpTimeMapBeforeEqParameter:                 {FactTimeMap, FactTimeMap},
		OpTimeMapAfterParameter:                    {FactTimeMap, FactTimeMap},
		OpTimeMapAfterEqParameter:                  {FactTimeMap, FactTimeMap},
	}
	if expectedTypes, usesParameter := parameterTypes[expression.Op]; usesParameter {
		parameter, found := parameters[expression.Parameter]
		if !found {
			return truth{reason: ReasonInvalidProgram}
		}
		if fact.Type != expectedTypes[0] || parameter.Type != expectedTypes[1] {
			return truth{reason: ReasonFactTypeMismatch}
		}
		switch expression.Op {
		case OpIdentityEqualParameter, OpSchemaEqualParameter, OpDigestEqualParameter, OpStringEqualParameter:
			return truth{known: true, value: *fact.String == *parameter.String}
		case OpStateInParameter:
			_, found := slices.BinarySearch(parameter.Strings, *fact.String)
			return truth{known: true, value: found}
		case OpBooleanEqualParameter:
			return truth{known: true, value: *fact.Boolean == *parameter.Boolean}
		case OpNumberEqualParameter, OpNumberLessParameter, OpNumberLessEqParameter, OpNumberGreaterParameter, OpNumberGreaterEqParameter:
			left, _ := new(big.Rat).SetString(fact.Number.String())
			right, _ := new(big.Rat).SetString(parameter.Number.String())
			comparison := left.Cmp(right)
			matches := map[Operation]bool{
				OpNumberEqualParameter: comparison == 0, OpNumberLessParameter: comparison < 0,
				OpNumberLessEqParameter: comparison <= 0, OpNumberGreaterParameter: comparison > 0,
				OpNumberGreaterEqParameter: comparison >= 0,
			}[expression.Op]
			return truth{known: true, value: matches}
		case OpTimeBeforeParameter, OpTimeBeforeEqParameter, OpTimeAfterParameter, OpTimeAfterEqParameter:
			left, _ := time.Parse(time.RFC3339Nano, *fact.Timestamp)
			right, _ := time.Parse(time.RFC3339Nano, *parameter.Timestamp)
			comparison := left.Compare(right)
			matches := map[Operation]bool{
				OpTimeBeforeParameter: comparison < 0, OpTimeBeforeEqParameter: comparison <= 0,
				OpTimeAfterParameter: comparison > 0, OpTimeAfterEqParameter: comparison >= 0,
			}[expression.Op]
			return truth{known: true, value: matches}
		case OpSetEqualParameter:
			return truth{known: true, value: slices.Equal(fact.Strings, parameter.Strings)}
		case OpSetContainsAllParameter:
			return truth{known: true, value: containsAll(fact.Strings, parameter.Strings)}
		case OpSetDisjointParameter:
			return truth{known: true, value: disjoint(fact.Strings, parameter.Strings)}
		case OpIdentityMapEqualParameter, OpSchemaMapEqualParameter, OpDigestMapEqualParameter,
			OpStateMapEqualParameter, OpStringMapEqualParameter:
			return truth{known: true, value: maps.Equal(fact.Values, parameter.Values)}
		case OpIdentityMapValuesDifferForPairsParameter:
			return truth{known: true, value: identityMapValuesDifferForPairs(fact.Values, parameter.Edges)}
		case OpBooleanMapEqualParameter:
			return truth{known: true, value: maps.Equal(fact.Booleans, parameter.Booleans)}
		case OpIdentityMapValuesInParameter, OpStateMapValuesInParameter, OpStringMapValuesInParameter:
			for _, value := range fact.Values {
				if _, found := slices.BinarySearch(parameter.Strings, value); !found {
					return truth{known: true}
				}
			}
			return truth{known: true, value: true}
		case OpIdentityMapValuesBijectSetParameter:
			return truth{known: true, value: identityMapValuesBijectSet(fact.Values, parameter.Strings)}
		case OpBooleanMapAllEqualParameter:
			for _, value := range fact.Booleans {
				if value != *parameter.Boolean {
					return truth{known: true}
				}
			}
			return truth{known: true, value: true}
		case OpNumberMapEqualParameter, OpNumberMapLessParameter, OpNumberMapLessEqParameter,
			OpNumberMapGreaterParameter, OpNumberMapGreaterEqParameter:
			return truth{known: true, value: compareNumberMaps(fact.Numbers, parameter.Numbers, expression.Op)}
		case OpTimeMapEqualParameter, OpTimeMapBeforeParameter, OpTimeMapBeforeEqParameter, OpTimeMapAfterParameter, OpTimeMapAfterEqParameter:
			return truth{known: true, value: compareTimeMaps(fact.Timestamps, parameter.Timestamps, expression.Op)}
		}
	}
	if expression.Op == OpMapKeysEqualSetParameter {
		parameter, found := parameters[expression.Parameter]
		if !found {
			return truth{reason: ReasonInvalidProgram}
		}
		keys, isMap := factMapKeys(fact)
		if !isMap || parameter.Type != FactStringSet {
			return truth{reason: ReasonFactTypeMismatch}
		}
		return truth{known: true, value: slices.Equal(keys, parameter.Strings)}
	}
	expectedType := map[Operation]FactType{
		OpIdentityEqual: FactIdentity, OpSchemaEqual: FactSchema, OpDigestEqual: FactDigest,
		OpStringEqual: FactString, OpStateIn: FactState, OpBooleanEqual: FactBoolean,
		OpNumberEqual: FactNumber, OpNumberLess: FactNumber, OpNumberLessEq: FactNumber,
		OpNumberGreater: FactNumber, OpNumberGreaterEq: FactNumber,
		OpTimeBefore: FactTime, OpTimeBeforeEq: FactTime, OpTimeAfter: FactTime, OpTimeAfterEq: FactTime,
		OpSetEqual: FactStringSet, OpSetContainsAll: FactStringSet, OpSetDisjoint: FactStringSet,
		OpIdentityMapValuesIn:     FactIdentityMap,
		OpStateMapValuesIn:        FactStateMap,
		OpStringMapValuesIn:       FactStringMap,
		OpBooleanMapAllEqual:      FactBooleanMap,
		OpStringMapAllNonempty:    FactStringMap,
		OpIdentityMapValuesUnique: FactIdentityMap,
		OpDirectedGraphAcyclic:    FactDirectedGraph,
	}[expression.Op]
	if fact.Type != expectedType {
		return truth{reason: ReasonFactTypeMismatch}
	}
	switch expression.Op {
	case OpIdentityEqual, OpSchemaEqual, OpDigestEqual, OpStringEqual:
		return truth{known: true, value: *fact.String == *expression.String}
	case OpStateIn:
		_, found := slices.BinarySearch(expression.Strings, *fact.String)
		return truth{known: true, value: found}
	case OpBooleanEqual:
		return truth{known: true, value: *fact.Boolean == *expression.Boolean}
	case OpNumberEqual, OpNumberLess, OpNumberLessEq, OpNumberGreater, OpNumberGreaterEq:
		left, _ := new(big.Rat).SetString(fact.Number.String())
		right, _ := new(big.Rat).SetString(expression.Number.String())
		comparison := left.Cmp(right)
		matches := map[Operation]bool{
			OpNumberEqual: comparison == 0, OpNumberLess: comparison < 0,
			OpNumberLessEq: comparison <= 0, OpNumberGreater: comparison > 0,
			OpNumberGreaterEq: comparison >= 0,
		}[expression.Op]
		return truth{known: true, value: matches}
	case OpTimeBefore, OpTimeBeforeEq, OpTimeAfter, OpTimeAfterEq:
		left, _ := time.Parse(time.RFC3339Nano, *fact.Timestamp)
		right, _ := time.Parse(time.RFC3339Nano, *expression.Timestamp)
		comparison := left.Compare(right)
		matches := map[Operation]bool{
			OpTimeBefore: comparison < 0, OpTimeBeforeEq: comparison <= 0,
			OpTimeAfter: comparison > 0, OpTimeAfterEq: comparison >= 0,
		}[expression.Op]
		return truth{known: true, value: matches}
	case OpSetEqual:
		return truth{known: true, value: slices.Equal(fact.Strings, expression.Strings)}
	case OpSetContainsAll:
		return truth{known: true, value: containsAll(fact.Strings, expression.Strings)}
	case OpSetDisjoint:
		return truth{known: true, value: disjoint(fact.Strings, expression.Strings)}
	case OpIdentityMapValuesIn, OpStateMapValuesIn, OpStringMapValuesIn:
		for _, value := range fact.Values {
			if _, found := slices.BinarySearch(expression.Strings, value); !found {
				return truth{known: true}
			}
		}
		return truth{known: true, value: true}
	case OpBooleanMapAllEqual:
		for _, value := range fact.Booleans {
			if value != *expression.Boolean {
				return truth{known: true}
			}
		}
		return truth{known: true, value: true}
	case OpDirectedGraphAcyclic:
		return truth{known: true, value: directedGraphAcyclic(fact.Edges)}
	case OpStringMapAllNonempty:
		for _, value := range fact.Values {
			if value == "" {
				return truth{known: true}
			}
		}
		return truth{known: true, value: true}
	case OpIdentityMapValuesUnique:
		seen := make(map[string]struct{}, len(fact.Values))
		for _, value := range fact.Values {
			if _, duplicate := seen[value]; duplicate {
				return truth{known: true}
			}
			seen[value] = struct{}{}
		}
		return truth{known: true, value: true}
	default:
		return truth{reason: ReasonInvalidProgram}
	}
}

func timeMapDeltasLessEq(starts, ends map[string]string, maximumSeconds json.Number) bool {
	if len(starts) != len(ends) {
		return false
	}
	maximum, ok := new(big.Rat).SetString(maximumSeconds.String())
	if !ok || maximum.Sign() < 0 {
		return false
	}
	maximum.Mul(maximum, big.NewRat(int64(time.Second), 1))
	for key, startValue := range starts {
		endValue, found := ends[key]
		if !found {
			return false
		}
		start, _ := time.Parse(time.RFC3339Nano, startValue)
		end, _ := time.Parse(time.RFC3339Nano, endValue)
		startNanos := new(big.Int).Mul(big.NewInt(start.Unix()), big.NewInt(int64(time.Second)))
		startNanos.Add(startNanos, big.NewInt(int64(start.Nanosecond())))
		endNanos := new(big.Int).Mul(big.NewInt(end.Unix()), big.NewInt(int64(time.Second)))
		endNanos.Add(endNanos, big.NewInt(int64(end.Nanosecond())))
		delta := new(big.Int).Sub(endNanos, startNanos)
		if delta.Sign() < 0 || new(big.Rat).SetInt(delta).Cmp(maximum) > 0 {
			return false
		}
	}
	return true
}

func timeMapDeltasEqualNumberMap(starts, ends map[string]string, durations map[string]json.Number) bool {
	if len(starts) == 0 || len(starts) != len(ends) || len(starts) != len(durations) {
		return false
	}
	for key, startValue := range starts {
		endValue, endFound := ends[key]
		durationValue, durationFound := durations[key]
		if !endFound || !durationFound {
			return false
		}
		start, _ := time.Parse(time.RFC3339Nano, startValue)
		end, _ := time.Parse(time.RFC3339Nano, endValue)
		if end.Before(start) {
			return false
		}
		observed, ok := new(big.Rat).SetString(durationValue.String())
		if !ok {
			return false
		}
		wholeSeconds := new(big.Int).Sub(big.NewInt(end.Unix()), big.NewInt(start.Unix()))
		totalNanoseconds := new(big.Int).Mul(wholeSeconds, big.NewInt(int64(time.Second)))
		totalNanoseconds.Add(totalNanoseconds, big.NewInt(int64(end.Nanosecond()-start.Nanosecond())))
		computedSeconds := new(big.Rat).SetFrac(totalNanoseconds, big.NewInt(int64(time.Second)))
		if observed.Cmp(computedSeconds) != 0 {
			return false
		}
	}
	return true
}

func directedGraphAcyclic(edges []DirectedEdge) bool {
	adjacency := make(map[string][]string)
	nodes := make(map[string]struct{})
	for _, edge := range edges {
		adjacency[edge.From] = append(adjacency[edge.From], edge.To)
		nodes[edge.From] = struct{}{}
		nodes[edge.To] = struct{}{}
	}
	const (
		unseen = iota
		visiting
		visited
	)
	states := make(map[string]int, len(nodes))
	var visit func(string) bool
	visit = func(node string) bool {
		switch states[node] {
		case visiting:
			return false
		case visited:
			return true
		}
		states[node] = visiting
		for _, next := range adjacency[node] {
			if !visit(next) {
				return false
			}
		}
		states[node] = visited
		return true
	}
	for node := range nodes {
		if states[node] == unseen && !visit(node) {
			return false
		}
	}
	return true
}

func factMapKeys(fact Fact) ([]string, bool) {
	var keys []string
	switch fact.Type {
	case FactIdentityMap, FactSchemaMap, FactDigestMap, FactStateMap, FactStringMap:
		keys = make([]string, 0, len(fact.Values))
		for key := range fact.Values {
			keys = append(keys, key)
		}
	case FactBooleanMap:
		keys = make([]string, 0, len(fact.Booleans))
		for key := range fact.Booleans {
			keys = append(keys, key)
		}
	case FactNumberMap:
		keys = make([]string, 0, len(fact.Numbers))
		for key := range fact.Numbers {
			keys = append(keys, key)
		}
	case FactTimeMap:
		keys = make([]string, 0, len(fact.Timestamps))
		for key := range fact.Timestamps {
			keys = append(keys, key)
		}
	default:
		return nil, false
	}
	sort.Strings(keys)
	return keys, true
}

func compareNumberMaps(observed, expected map[string]json.Number, operation Operation) bool {
	// Map predicates are exact-domain comparisons. Silently ignoring an extra
	// observed key can hide an out-of-policy subject while still returning
	// Pass. Callers that need subset semantics must model that scope explicitly
	// with a reviewed subject/set predicate.
	if len(observed) != len(expected) {
		return false
	}
	for key, rightValue := range expected {
		leftValue, found := observed[key]
		if !found {
			return false
		}
		left, _ := new(big.Rat).SetString(leftValue.String())
		right, _ := new(big.Rat).SetString(rightValue.String())
		comparison := left.Cmp(right)
		if !map[Operation]bool{
			OpNumberMapEqualParameter:     comparison == 0,
			OpNumberMapLessParameter:      comparison < 0,
			OpNumberMapLessEqParameter:    comparison <= 0,
			OpNumberMapGreaterParameter:   comparison > 0,
			OpNumberMapGreaterEqParameter: comparison >= 0,
			OpNumberMapEqualFact:          comparison == 0,
			OpNumberMapLessFact:           comparison < 0,
			OpNumberMapLessEqFact:         comparison <= 0,
			OpNumberMapGreaterFact:        comparison > 0,
			OpNumberMapGreaterEqFact:      comparison >= 0,
		}[operation] {
			return false
		}
	}
	return true
}

func compareTimeMaps(observed, expected map[string]string, operation Operation) bool {
	if len(observed) != len(expected) {
		return false
	}
	for key, rightValue := range expected {
		leftValue, found := observed[key]
		if !found {
			return false
		}
		left, _ := time.Parse(time.RFC3339Nano, leftValue)
		right, _ := time.Parse(time.RFC3339Nano, rightValue)
		comparison := left.Compare(right)
		if !map[Operation]bool{
			OpTimeMapEqualParameter:    comparison == 0,
			OpTimeMapBeforeParameter:   comparison < 0,
			OpTimeMapBeforeEqParameter: comparison <= 0,
			OpTimeMapAfterParameter:    comparison > 0,
			OpTimeMapAfterEqParameter:  comparison >= 0,
			OpTimeMapEqualFact:         comparison == 0,
			OpTimeMapBeforeFact:        comparison < 0,
			OpTimeMapBeforeEqFact:      comparison <= 0,
			OpTimeMapAfterFact:         comparison > 0,
			OpTimeMapAfterEqFact:       comparison >= 0,
		}[operation] {
			return false
		}
	}
	return true
}

func identityMapValuesIn(references, inventory map[string]string) bool {
	allowed := make(map[string]struct{}, len(inventory))
	for _, identity := range inventory {
		if _, duplicate := allowed[identity]; duplicate {
			return false
		}
		allowed[identity] = struct{}{}
	}
	for _, identity := range references {
		if _, found := allowed[identity]; !found {
			return false
		}
	}
	return true
}

func identityMapValuesNotIn(references, inventory map[string]string) bool {
	denied := make(map[string]struct{}, len(inventory))
	for _, identity := range inventory {
		if _, duplicate := denied[identity]; duplicate {
			return false
		}
		denied[identity] = struct{}{}
	}
	for _, identity := range references {
		if _, found := denied[identity]; found {
			return false
		}
	}
	return true
}

func identityMapValuesDifferForPairs(identities map[string]string, pairs []DirectedEdge) bool {
	if len(identities) == 0 || len(pairs) == 0 {
		return false
	}
	for _, pair := range pairs {
		if pair.From >= pair.To {
			return false
		}
		left, leftFound := identities[pair.From]
		right, rightFound := identities[pair.To]
		if !leftFound || !rightFound || left == right {
			return false
		}
	}
	return true
}

func identityMapValuesBijectSet(values map[string]string, expected []string) bool {
	if len(values) == 0 || len(values) != len(expected) {
		return false
	}
	observed := make([]string, 0, len(values))
	seen := make(map[string]struct{}, len(values))
	for _, value := range values {
		if _, duplicate := seen[value]; duplicate {
			return false
		}
		seen[value] = struct{}{}
		observed = append(observed, value)
	}
	slices.Sort(observed)
	return slices.Equal(observed, expected)
}

func mapKeyPartitionEqualSet(left, right, expected []string) bool {
	combined := make([]string, 0, len(left)+len(right))
	i, j := 0, 0
	for i < len(left) && j < len(right) {
		if left[i] == right[j] {
			return false
		}
		if left[i] < right[j] {
			combined = append(combined, left[i])
			i++
		} else {
			combined = append(combined, right[j])
			j++
		}
	}
	combined = append(combined, left[i:]...)
	combined = append(combined, right[j:]...)
	return slices.Equal(combined, expected)
}

func containsAll(actual, required []string) bool {
	for _, value := range required {
		if _, found := slices.BinarySearch(actual, value); !found {
			return false
		}
	}
	return true
}

func disjoint(left, right []string) bool {
	for _, value := range right {
		if _, found := slices.BinarySearch(left, value); found {
			return false
		}
	}
	return true
}
