#!/usr/bin/env python3
from __future__ import annotations

import copy
import hashlib
import json
import re
from collections import Counter, defaultdict
from pathlib import Path

ROOT = Path(__file__).resolve().parents[3]
BINDINGS = ROOT / "catalog/control-check-bindings.json"
OUT = Path(__file__).with_name("non-structured.json")

GENERIC_PREFIXES = {
    "semantic_wrapper_environment": "Read-only effective configuration, state, and event evidence for every subject in the complete assessed inventory directly demonstrates the operating behavior required by the original control;",
    "semantic_wrapper_execution": "Authenticated raw execution evidence covers every case in the pre-sealed scenario and subject inventory, is bound to the exact tested revisions and inputs, and records the objective observations required by the original control:",
}

RECLASSIFY = {
    "USEQ-27B484EB": "The undefined instruction to minimize failures has no objective optimum or threshold, so two correct implementations can disagree from identical evidence.",
    "USEQ-3847EF4D": "Resolution quality, user effort, and satisfaction lack a fixed population, instrument, sampling protocol, and acceptance interpretation.",
    "USEQ-A96AE686": "User satisfaction and the combined multi-metric evaluation lack a fixed empirical protocol, denominator, and decision rule.",
    "USEQ-ADE1D409": "User satisfaction lacks a fixed population, instrument, sampling protocol, and objective outcome rule.",
}


def prefix(control_id, ordinal):
    return control_id.lower().replace("-", "_") + f"_c{ordinal}"


def fact(key, typ, semantics, authority, source):
    return {"fact_id": key, "fact_type": typ, "authority": authority, "raw_value_semantics": semantics, "source_requirement": source, "complete_required": True}


def parameter(key, typ, semantics, origin=None):
    """Declare which trust lane supplies an expected value.

    Scope/domain inventories are discovered by the scanner before evidence is
    collected.  Project choices (thresholds, approved identities, expected
    values, and policy digests) must come from authenticated policy.  Values
    published by an independent authority use authenticated context.  Keeping
    these lanes separate prevents an evidence provider from selecting the value
    that will make its own observation pass.
    """
    lowered = f"{key} {semantics}".lower()
    if origin is None:
        if any(token in lowered for token in (
            "issuing body", "issuing-body", "external registry", "current baseline",
            "published baseline", "registry authority",
        )):
            origin = "independently_authenticated_context"
        elif any(token in lowered for token in (
            "approved", "expected", "threshold", "maximum", "minimum", "policy",
            "required analyzer", "required config", "accepted", "allowed", "denied",
        )):
            origin = "independently_authenticated_policy"
        else:
            origin = "scanner_inventory"
    return {"parameter_id": key, "parameter_type": typ, "value_origin": origin, "source_requirement": semantics}


def fvalue(typ, value, complete=True):
    if typ == "string_set":
        value = sorted(canonical_identity(item) for item in value)
    elif typ.endswith("_map"):
        value = {canonical_identity(key): item for key, item in value.items()}
    elif typ == "directed_graph":
        value = canonical_pairs(value)
    result = {"type": typ, "complete": complete}
    field = {"identity":"string","schema":"string","digest":"string","state":"string","string":"string","boolean":"boolean","number":"number","time":"timestamp","string_set":"strings","identity_map":"values","schema_map":"values","digest_map":"values","state_map":"values","string_map":"values","boolean_map":"booleans","number_map":"numbers","time_map":"timestamps","directed_graph":"edges"}[typ]
    result[field] = value
    return result


def pvalue(typ, value):
    if typ == "string_set":
        value = sorted(canonical_identity(item) for item in value)
    elif typ.endswith("_map"):
        value = {canonical_identity(key): item for key, item in value.items()}
    elif typ == "directed_graph":
        value = canonical_pairs(value)
    result = {"type": typ}
    field = {"identity":"string","schema":"string","digest":"string","state":"string","string":"string","boolean":"boolean","number":"number","time":"timestamp","string_set":"strings","identity_map":"values","schema_map":"values","digest_map":"values","state_map":"values","string_map":"values","boolean_map":"booleans","number_map":"numbers","time_map":"timestamps","directed_graph":"edges"}[typ]
    result[field] = value
    return result


def canonical_identity(value):
    """Encode example subject identities into the evaluator's canonical key alphabet."""
    return re.sub(r"[^A-Za-z0-9_.-]+", "--", value)


def canonical_pairs(edges):
    pairs = {
        tuple(sorted((canonical_identity(edge["from"]), canonical_identity(edge["to"]))))
        for edge in edges
    }
    return [{"from": left, "to": right} for left, right in sorted(pairs) if left != right]


def exact_profile(facts, params, predicate, pass_facts, fail_facts, param_values, pass_why, fail_why):
    blocked = copy.deepcopy(pass_facts)
    first = next(iter(blocked))
    blocked[first]["complete"] = False
    return {
        "facts": facts, "params": params, "predicate": predicate,
        "fixtures": {
            "pass": {"description": pass_why, "parameters": param_values, "facts": pass_facts, "expected_outcome": "pass"},
            "fail": {"description": fail_why, "parameters": param_values, "facts": fail_facts, "expected_outcome": "fail"},
            "blocked": {"description": "One required raw fact is incomplete; the evaluator must not infer Pass or Fail.", "parameters": param_values, "facts": blocked, "expected_outcome": "blocked"},
            "counterexample": {"description": "A tempting partial signal is present, but the exact predicate still rejects the broken required value.", "parameters": param_values, "facts": copy.deepcopy(fail_facts), "expected_outcome": "fail"},
        },
    }


def digest_pair(control_id, ordinal, authority, left_name, right_name, statement, algorithm=False):
    pre=prefix(control_id,ordinal); left=f"{pre}.{left_name}"; right=f"{pre}.{right_name}"
    facts=[fact(left,"digest",f"Raw SHA-256 digest for {left_name.replace('_',' ')} bound to the assessed subject.",authority,"Exact bytes or authenticated immutable record named by the clause."),fact(right,"digest",f"Independently derived SHA-256 digest for {right_name.replace('_',' ')} bound to the same subject.",authority,"Independent bytes or immutable record on the other side of the clause equality.")]
    pred={"op":"digest_eq_fact","fact":left,"other_fact":right}; params=[]; pvals={}
    pf={left:fvalue("digest","a"*64),right:fvalue("digest","a"*64)}; ff={left:fvalue("digest","a"*64),right:fvalue("digest","b"*64)}
    if algorithm:
        alg=f"{pre}.recorded_digest_algorithm"; allowed=f"{pre}.approved_digest_algorithms"
        facts.append(fact(alg,"state","Raw digest algorithm identifier recorded for the subject.",authority,"Authenticated manifest or record containing the algorithm identifier."))
        params.append(parameter(allowed,"string_set","Approved digest algorithm identifiers sealed before evidence collection.")); pvals[allowed]=pvalue("string_set",["sha256"])
        pf[alg]=fvalue("state","sha256");ff[alg]=fvalue("state","sha256")
        pred={"op":"all","args":[{"op":"state_in_parameter","fact":alg,"parameter":allowed},pred]}
    return exact_profile(facts,params,pred,pf,ff,pvals,f"Both independently obtained digests match for the full clause: {statement}","The complete authoritative digests differ, so the clause is false.")


def set_parameter_profile(control_id, ordinal, authority, observed_suffix, required_values, semantics, mode="contains"):
    pre=prefix(control_id,ordinal); observed=f"{pre}.{observed_suffix}"; expected=f"{pre}.required_{observed_suffix}"
    facts=[fact(observed,"string_set",semantics,authority,"Lossless parser, read-only configuration, raw event stream, or scanner inventory explicitly named by the clause.")]
    params=[parameter(expected,"string_set","Complete required identity/value set sealed independently before observed evidence is requested.")]
    op={"contains":"set_contains_all_parameter","eq":"set_eq_parameter","disjoint":"set_disjoint_parameter"}[mode]
    good=sorted(required_values); bad=good[:-1] if mode!="disjoint" else good
    if mode=="disjoint": good=[]; bad=[good[0] if good else required_values[0]]
    pf={observed:fvalue("string_set",good)}; ff={observed:fvalue("string_set",bad)}; pv={expected:pvalue("string_set",sorted(required_values))}
    return exact_profile(facts,params,{"op":op,"fact":observed,"parameter":expected},pf,ff,pv,"The raw observed set satisfies the independently sealed complete set contract.","A required value is absent or a prohibited value is present in complete raw evidence.")


def property_coverage(control_id, ordinal, authority, property_suffixes, semantics):
    pre=prefix(control_id,ordinal); inventory=f"{pre}.required_subject_ids"
    facts=[];params=[parameter(inventory,"string_set","Complete stable in-scope subject identities sealed independently before collection.")]
    args=[]; pf={};ff={};pv={inventory:pvalue("string_set",["subject-A","subject-B"])}
    for index,suffix in enumerate(property_suffixes):
        key=f"{pre}.{suffix}"
        facts.append(fact(key,"string_set",semantics[suffix],authority,"Raw source fields selecting subject identities with the exact named property; no compliance verdict or score."))
        args.append({"op":"set_eq_parameter","fact":key,"parameter":inventory})
        pf[key]=fvalue("string_set",["subject-A","subject-B"]);ff[key]=fvalue("string_set",["subject-A"] if index==0 else ["subject-A","subject-B"])
    return exact_profile(facts,params,{"op":"all","args":args},pf,ff,pv,"Every raw-property subject set equals the independently sealed complete subject inventory.","At least one complete raw-property set omits an independently inventoried subject.")


def two_sets(control_id,ordinal,authority,left_suffix,right_suffix,semantics,op="set_eq_fact"):
    pre=prefix(control_id,ordinal); left=f"{pre}.{left_suffix}";right=f"{pre}.{right_suffix}"
    facts=[fact(left,"string_set",semantics[0],authority,"First independently enumerated raw source."),fact(right,"string_set",semantics[1],authority,"Second independently enumerated raw source.")]
    left_expected=f"{pre}.required_{left_suffix}";right_expected=f"{pre}.required_{right_suffix}"
    if op=="set_disjoint_fact":
        left_values=["subject-A","subject-B"];right_values=["subject-C"];params=[parameter(left_expected,"string_set","Complete independently sealed first-side subject identities."),parameter(right_expected,"string_set","Complete independently sealed second-side subject identities.")];pred={"op":"all","args":[{"op":"set_eq_parameter","fact":left,"parameter":left_expected},{"op":"set_eq_parameter","fact":right,"parameter":right_expected},{"op":op,"fact":left,"other_fact":right}]};pv={left_expected:pvalue("string_set",left_values),right_expected:pvalue("string_set",right_values)};pf={left:fvalue("string_set",left_values),right:fvalue("string_set",right_values)};ff={left:fvalue("string_set",left_values),right:fvalue("string_set",["subject-B","subject-C"])}
    else:
        values=["subject-A","subject-B"];params=[parameter(left_expected,"string_set","Complete independently sealed subject identities that both authoritative sources must enumerate.")];pred={"op":"all","args":[{"op":"set_eq_parameter","fact":left,"parameter":left_expected},{"op":"set_eq_parameter","fact":right,"parameter":left_expected},{"op":op,"fact":left,"other_fact":right}]};pv={left_expected:pvalue("string_set",values)};pf={left:fvalue("string_set",values),right:fvalue("string_set",values)};ff={left:fvalue("string_set",values),right:fvalue("string_set",["subject-A"])}
    return exact_profile(facts,params,pred,pf,ff,pv,"Both raw sets match their independently sealed complete inventories and satisfy the exact relation.","A raw set omits or misassigns an independently inventoried subject.")


def map_equal_profile(control_id, ordinal, authority, map_type, left_suffix, right_suffix, left_semantics, right_semantics):
    pre=prefix(control_id,ordinal); left=f"{pre}.{left_suffix}"; right=f"{pre}.{right_suffix}"; subjects=f"{pre}.required_subject_ids"
    facts=[fact(left,map_type,left_semantics,authority,"Complete authenticated first-side raw records for every sealed subject."),fact(right,map_type,right_semantics,authority,"Complete independently obtained second-side raw records for every sealed subject.")]
    params=[parameter(subjects,"string_set","Complete stable subject identities sealed from the authoritative scanner inventory before collection.")]
    op_prefix={"identity_map":"identity_map","schema_map":"schema_map","digest_map":"digest_map","state_map":"state_map","string_map":"string_map","boolean_map":"boolean_map","number_map":"number_map","time_map":"time_map"}[map_type]
    args=[{"op":"map_keys_eq_set_parameter","fact":left,"parameter":subjects},{"op":"map_keys_eq_set_parameter","fact":right,"parameter":subjects},{"op":f"{op_prefix}_eq_fact","fact":left,"other_fact":right}]
    if map_type=="digest_map": good={"subject-A":"a"*64,"subject-B":"b"*64};bad={"subject-A":"c"*64,"subject-B":"b"*64}
    elif map_type=="boolean_map": good={"subject-A":True,"subject-B":False};bad={"subject-A":False,"subject-B":False}
    elif map_type=="number_map": good={"subject-A":1,"subject-B":2};bad={"subject-A":3,"subject-B":2}
    elif map_type=="time_map": good={"subject-A":"2026-08-29T10:00:00Z","subject-B":"2026-08-29T11:00:00Z"};bad={"subject-A":"2026-08-29T12:00:00Z","subject-B":"2026-08-29T11:00:00Z"}
    else: good={"subject-A":"value-A","subject-B":"value-B"};bad={"subject-A":"wrong-A","subject-B":"value-B"}
    pf={left:fvalue(map_type,good),right:fvalue(map_type,good)};ff={left:fvalue(map_type,good),right:fvalue(map_type,bad)};pv={subjects:pvalue("string_set",["subject-A","subject-B"])}
    return exact_profile(facts,params,{"op":"all","args":args},pf,ff,pv,"Both direct raw maps cover the complete sealed subject inventory and match for every subject.","One subject has a different raw bound value or is omitted, so the complete map relation fails.")


def map_parameter_profile(control_id, ordinal, authority, map_type, observed_suffix, semantics):
    pre=prefix(control_id,ordinal); observed=f"{pre}.{observed_suffix}"; expected=f"{pre}.required_{observed_suffix}"
    facts=[fact(observed,map_type,semantics,authority,"Complete authenticated raw values for every key in the independently sealed subject inventory.")]
    params=[parameter(expected,map_type,"Exact expected raw values keyed by every sealed subject, authenticated independently before observation.")]
    op_prefix={"identity_map":"identity_map","schema_map":"schema_map","digest_map":"digest_map","state_map":"state_map","string_map":"string_map","boolean_map":"boolean_map"}[map_type]
    if map_type=="digest_map": good={"subject-A":"a"*64,"subject-B":"b"*64};bad={"subject-A":"c"*64,"subject-B":"b"*64}
    elif map_type=="boolean_map": good={"subject-A":True,"subject-B":True};bad={"subject-A":False,"subject-B":True}
    else: good={"subject-A":"value-A","subject-B":"value-B"};bad={"subject-A":"wrong-A","subject-B":"value-B"}
    return exact_profile(facts,params,{"op":f"{op_prefix}_eq_parameter","fact":observed,"parameter":expected},{observed:fvalue(map_type,good)},{observed:fvalue(map_type,bad)},{expected:pvalue(map_type,good)},"The complete raw keyed values exactly equal the independently sealed keyed values.","At least one raw keyed value differs from the independently sealed expected value.")


def scalar_parameter_profile(control_id, ordinal, authority, scalar_type, observed_suffix, semantics, example_value=None):
    pre=prefix(control_id,ordinal);observed=f"{pre}.{observed_suffix}";expected=f"{pre}.required_{observed_suffix}"
    facts=[fact(observed,scalar_type,semantics,authority,"Direct authenticated scalar value from the exact authoritative object named by the clause.")]
    params=[parameter(expected,scalar_type,f"Exact independently approved {observed_suffix.replace('_',' ')}." )]
    if example_value is None:
        example_value={"identity":"identity-A","schema":"schema-A","digest":"a"*64,"state":"state-A","string":"value-A","boolean":True,"number":1,"time":"2026-08-29T10:00:00Z"}[scalar_type]
    if scalar_type=="number": wrong=float(example_value)+1
    elif scalar_type=="boolean": wrong=not example_value
    else: wrong={"identity":"identity-B","schema":"schema-B","digest":"b"*64,"state":"state-B","string":"value-B","time":"2026-08-30T10:00:00Z"}[scalar_type]
    return exact_profile(facts,params,{"op":f"{scalar_type}_eq_parameter","fact":observed,"parameter":expected},{observed:fvalue(scalar_type,example_value)},{observed:fvalue(scalar_type,wrong)},{expected:pvalue(scalar_type,example_value)},f"The direct raw {observed_suffix.replace('_',' ')} equals its independently approved exact value.",f"The direct raw {observed_suffix.replace('_',' ')} differs from its independently approved exact value.")


def boolean_inventory_profile(control_id, ordinal, authority, suffix, expected_value, semantics):
    pre=prefix(control_id,ordinal); observed=f"{pre}.{suffix}"; subjects=f"{pre}.required_subject_ids"; expected=f"{pre}.required_boolean_value"
    facts=[fact(observed,"boolean_map",semantics,authority,"Complete direct raw boolean fields keyed by every sealed subject; never a provider compliance conclusion.")]
    params=[parameter(subjects,"string_set","Complete stable subject identities sealed from the authoritative inventory."),parameter(expected,"boolean",f"Required raw boolean value {str(expected_value).lower()} sealed independently." )]
    args=[{"op":"map_keys_eq_set_parameter","fact":observed,"parameter":subjects},{"op":"boolean_map_all_eq_parameter","fact":observed,"parameter":expected}]
    good={"subject-A":expected_value,"subject-B":expected_value};bad={"subject-A":not expected_value,"subject-B":expected_value}
    return exact_profile(facts,params,{"op":"all","args":args},{observed:fvalue("boolean_map",good)},{observed:fvalue("boolean_map",bad)},{subjects:pvalue("string_set",["subject-A","subject-B"]),expected:pvalue("boolean",expected_value)},"Every sealed subject has the required direct raw boolean field value.","At least one sealed subject has the opposite direct raw field value.")


def boolean_fields_inventory_profile(control_id, ordinal, authority, fields, expected_value, subject_description):
    pre=prefix(control_id,ordinal);subjects=f"{pre}.required_subject_ids"
    facts=[];args=[];pf={};ff={};ids=["subject-A","subject-B"]
    for index,(suffix,semantics) in enumerate(fields):
        observed=f"{pre}.{suffix}"
        facts.append(fact(observed,"boolean_map",semantics,authority,"Complete direct raw boolean field keyed by every independently sealed subject; never a provider conclusion."))
        args += [{"op":"map_keys_eq_set_parameter","fact":observed,"parameter":subjects},{"op":"boolean_map_all_eq","fact":observed,"boolean":expected_value}]
        good={x:expected_value for x in ids};pf[observed]=fvalue("boolean_map",good)
        ff[observed]=fvalue("boolean_map",({"subject-A":not expected_value,"subject-B":expected_value} if index==0 else good))
    params=[parameter(subjects,"string_set",f"Complete {subject_description} sealed independently before collection.")]
    return exact_profile(facts,params,{"op":"all","args":args},pf,ff,{subjects:pvalue("string_set",ids)},f"Every direct raw Boolean field equals the clause's fixed value {str(expected_value).lower()} for every sealed subject.","One direct raw Boolean field has the opposite value for an independently sealed subject.")


def state_map_allowed_profile(control_id, ordinal, authority, suffix, allowed_values, semantics):
    pre=prefix(control_id,ordinal);observed=f"{pre}.{suffix}";subjects=f"{pre}.required_subject_ids";allowed=f"{pre}.approved_state_values"
    facts=[fact(observed,"state_map",semantics,authority,"Complete direct raw state values keyed by every independently sealed subject.")]
    params=[parameter(subjects,"string_set","Complete stable subject identities sealed independently before collection."),parameter(allowed,"string_set","Approved raw state identifiers sealed independently before collection.")]
    args=[{"op":"map_keys_eq_set_parameter","fact":observed,"parameter":subjects},{"op":"state_map_values_in_parameter","fact":observed,"parameter":allowed}]
    good={"subject-A":allowed_values[0],"subject-B":allowed_values[-1]};bad={"subject-A":"unapproved-state","subject-B":allowed_values[-1]}
    return exact_profile(facts,params,{"op":"all","args":args},{observed:fvalue("state_map",good)},{observed:fvalue("state_map",bad)},{subjects:pvalue("string_set",["subject-A","subject-B"]),allowed:pvalue("string_set",sorted(allowed_values))},"Every sealed subject has a direct raw state in the independently approved set.","One sealed subject has a direct raw state outside the approved set.")


def nonempty_string_map_profile(control_id, ordinal, authority, suffix, required_keys, semantics, inventory_semantics):
    pre=prefix(control_id,ordinal);observed=f"{pre}.{suffix}";required=f"{pre}.required_{suffix}_keys"
    facts=[fact(observed,"string_map",semantics,authority,"Lossless parsed direct text values keyed by every independently sealed required identity.")]
    params=[parameter(required,"string_set",inventory_semantics)]
    good={key:f"nonempty-{key}" for key in required_keys};bad=dict(good);bad[required_keys[-1]]=""
    pred={"op":"all","args":[{"op":"map_keys_eq_set_parameter","fact":observed,"parameter":required},{"op":"string_map_all_nonempty","fact":observed}]}
    result=exact_profile(facts,params,pred,{observed:fvalue("string_map",good)},{observed:fvalue("string_map",bad)},{required:pvalue("string_set",required_keys)},"Every independently required key is present with a direct nonempty parsed value.",f"Required key {required_keys[-1]} exists but its direct parsed content is empty.")
    result["fixtures"]["counterexample"]=copy.deepcopy(result["fixtures"]["fail"]);result["fixtures"]["counterexample"]["description"]=f"A heading or key named {required_keys[-1]} exists, but its content is empty; name presence alone cannot pass."
    return result


def expected_field_maps_profile(control_id, ordinal, authority, subject_suffix, fields, subject_semantics, subject_examples=None):
    """Compare clause-specific raw maps to independently sealed exact case maps."""
    pre=prefix(control_id,ordinal);subjects=f"{pre}.required_{subject_suffix}"
    ids=subject_examples or ["case-A","case-B"]
    facts=[];params=[parameter(subjects,"string_set",subject_semantics)];args=[];pf={};ff={};pv={subjects:pvalue("string_set",ids)}
    op_prefix={"identity_map":"identity_map","schema_map":"schema_map","digest_map":"digest_map","state_map":"state_map","string_map":"string_map","boolean_map":"boolean_map","number_map":"number_map","time_map":"time_map"}
    for index,(suffix,map_type,semantics) in enumerate(fields):
        observed=f"{pre}.{suffix}";expected=f"{pre}.expected_{suffix}"
        facts.append(fact(observed,map_type,semantics,authority,"Complete direct raw values keyed by the independently sealed exact case or subject identities."))
        params.append(parameter(expected,map_type,f"Exact independently approved expected {suffix.replace('_',' ')} keyed by every sealed case or subject."))
        args += [{"op":"map_keys_eq_set_parameter","fact":observed,"parameter":subjects},{"op":f"{op_prefix[map_type]}_eq_parameter","fact":observed,"parameter":expected}]
        if map_type=="digest_map": values={key:chr(97+(i%3))*64 for i,key in enumerate(ids)};wrong="f"*64
        elif map_type=="boolean_map": values={key:True for key in ids};wrong=False
        elif map_type=="number_map": values={key:i+1 for i,key in enumerate(ids)};wrong=99
        elif map_type=="time_map": values={key:f"2026-08-29T{10+i:02d}:00:00Z" for i,key in enumerate(ids)};wrong="2026-08-30T00:00:00Z"
        else: values={key:f"expected-{suffix}-{i}" for i,key in enumerate(ids)};wrong=f"wrong-{suffix}"
        pf[observed]=fvalue(map_type,values);pv[expected]=pvalue(map_type,values)
        if index==0:
            broken=dict(values);broken[ids[0]]=wrong;ff[observed]=fvalue(map_type,broken)
        else: ff[observed]=fvalue(map_type,values)
    result=exact_profile(facts,params,{"op":"all","args":args},pf,ff,pv,"Every clause-specific direct raw field covers the exact sealed domain and equals its independently approved case value.","One named direct raw field differs from the independently approved case oracle while every other field remains correct.")
    result["fixtures"]["counterexample"]=copy.deepcopy(result["fixtures"]["fail"]);result["fixtures"]["counterexample"]["description"]="All generic execution states could look successful, but one clause-specific raw oracle field is wrong; the exact field comparison rejects that false pass."
    return result


def keyed_map_coverage_profile(control_id, ordinal, authority, map_type, suffix, semantics):
    pre=prefix(control_id,ordinal); observed=f"{pre}.{suffix}"; required=f"{pre}.required_subject_ids"
    facts=[fact(observed,map_type,semantics,authority,"Complete direct raw keyed records from the exact authoritative source named by the clause.")]
    params=[parameter(required,"string_set","Complete stable subject or case identities sealed independently before collection.")]
    if map_type=="digest_map": good={"subject-A":"a"*64,"subject-B":"b"*64}
    elif map_type=="boolean_map": good={"subject-A":True,"subject-B":False}
    elif map_type=="number_map": good={"subject-A":1,"subject-B":2}
    elif map_type=="time_map": good={"subject-A":"2026-08-29T10:00:00Z","subject-B":"2026-08-29T11:00:00Z"}
    else: good={"subject-A":"value-A","subject-B":"value-B"}
    bad={"subject-A":good["subject-A"]}
    return exact_profile(facts,params,{"op":"map_keys_eq_set_parameter","fact":observed,"parameter":required},{observed:fvalue(map_type,good)},{observed:fvalue(map_type,bad)},{required:pvalue("string_set",["subject-A","subject-B"] )},"Direct raw keyed observations cover every independently sealed subject or case.","A sealed subject or case has no direct raw keyed observation.")


def measurement_coverage_profile(control_id, ordinal, dimensions, key_examples, scope_description):
    """Prove coverage from an independently sealed subject-by-measurement domain."""
    pre=prefix(control_id,ordinal)
    values=f"{pre}.raw_measurement_values_by_subject_measurement_key"
    units=f"{pre}.raw_measurement_unit_identities_by_subject_measurement_key"
    windows=f"{pre}.raw_measurement_window_identities_by_subject_measurement_key"
    required=f"{pre}.required_subject_measurement_keys"
    expected_units=f"{pre}.approved_unit_identities_by_subject_measurement_key"
    expected_windows=f"{pre}.approved_window_identities_by_subject_measurement_key"
    dimension_text=", ".join(dimensions)
    facts=[
        fact(values,"number_map",f"Direct raw numeric measurements for every exact {scope_description} cross-product key. Required dimensions: {dimension_text}.","executed","Authenticated raw execution or observation output keyed only by independently sealed subject and measurement identities."),
        fact(units,"identity_map",f"Direct unit identity attached to every raw measurement across: {dimension_text}.","executed","Authenticated unit metadata from the same raw measurement records."),
        fact(windows,"identity_map",f"Direct observation-window identity attached to every raw measurement across: {dimension_text}.","executed","Authenticated window metadata from the same raw measurement records."),
    ]
    params=[
        parameter(required,"string_set",f"Complete exact {scope_description} cross-product sealed independently before execution; every applicable subject is combined with every required dimension: {dimension_text}."),
        parameter(expected_units,"identity_map",f"Independently approved exact unit identity for every sealed key across: {dimension_text}."),
        parameter(expected_windows,"identity_map",f"Independently approved exact observation-window identity for every sealed key across: {dimension_text}."),
    ]
    args=[{"op":"map_keys_eq_set_parameter","fact":raw,"parameter":required} for raw in [values,units,windows]]+[
        {"op":"identity_map_eq_parameter","fact":units,"parameter":expected_units},
        {"op":"identity_map_eq_parameter","fact":windows,"parameter":expected_windows},
    ]
    good_values={key:index+1 for index,key in enumerate(key_examples)}
    good_units={key:"approved-unit" for key in key_examples}
    good_windows={key:"approved-window" for key in key_examples}
    pf={values:fvalue("number_map",good_values),units:fvalue("identity_map",good_units),windows:fvalue("identity_map",good_windows)}
    missing=key_examples[-1]
    ff=copy.deepcopy(pf);ff[values]=fvalue("number_map",{key:value for key,value in good_values.items() if key != missing})
    pv={required:pvalue("string_set",key_examples),expected_units:pvalue("identity_map",good_units),expected_windows:pvalue("identity_map",good_windows)}
    result=exact_profile(facts,params,{"op":"all","args":args},pf,ff,pv,"Every named dimension is measured for the complete independently sealed subject cross-product, with the exact approved unit and window bindings.",f"The raw value for required subject-measurement key {missing} is absent even though the remaining measurements and metadata look valid.")
    result["fixtures"]["counterexample"]=copy.deepcopy(result["fixtures"]["fail"])
    result["fixtures"]["counterexample"]["description"]=f"A provider reports every other named measure but silently omits {missing}; exact sealed cross-product coverage rejects the false pass."
    return result


def previous_artifact_registry_profile():
    cid="PRC-08-021";ordinal=2;pre=prefix(cid,ordinal)
    recorded=f"{pre}.recorded_known_good_artifact_digests";retrieved=f"{pre}.retrieved_registry_artifact_byte_digests";registries=f"{pre}.source_registry_identities"
    required=f"{pre}.required_previous_artifact_ids";approved=f"{pre}.approved_registry_identities"
    facts=[
        fact(recorded,"digest_map","Direct previous-known-good digest from the authenticated release record, keyed by exact artifact identity.","artifact","Authenticated release record fields for every independently sealed previous-artifact identity."),
        fact(retrieved,"digest_map","Locally recomputed digest of exact bytes currently retrieved from the registry, keyed by the same artifact identity.","artifact","Exact retrieved bytes hashed locally by the pinned scanner digest implementation."),
        fact(registries,"identity_map","Direct registry identity from which each exact artifact byte stream was retrieved.","artifact","Authenticated registry endpoint identity bound to each retrieval."),
    ]
    params=[parameter(required,"string_set","Complete previous-known-good artifact identities sealed from the release record before retrieval."),parameter(approved,"string_set","Independently approved source registry identities sealed by release policy before retrieval.")]
    args=[{"op":"map_keys_eq_set_parameter","fact":x,"parameter":required} for x in [recorded,retrieved,registries]]+[{"op":"digest_map_eq_fact","fact":recorded,"other_fact":retrieved},{"op":"identity_map_values_in_parameter","fact":registries,"parameter":approved}]
    pf={recorded:fvalue("digest_map",{"artifact-A":"a"*64}),retrieved:fvalue("digest_map",{"artifact-A":"a"*64}),registries:fvalue("identity_map",{"artifact-A":"registry-A"})}
    ff=copy.deepcopy(pf);ff[registries]=fvalue("identity_map",{"artifact-A":"unapproved-registry"})
    pv={required:pvalue("string_set",["artifact-A"]),approved:pvalue("string_set",["registry-A","registry-B"])}
    result=exact_profile(facts,params,{"op":"all","args":args},pf,ff,pv,"The recorded artifact bytes are retrievable now, their recomputed digest matches, and the retrieval source is independently approved.","Matching artifact bytes came from a registry outside the independently approved registry set.")
    result["fixtures"]["counterexample"]=copy.deepcopy(result["fixtures"]["fail"]);result["fixtures"]["counterexample"]["description"]="The digest matches, but an unapproved registry supplied the bytes; digest equality alone cannot prove the approved-registry promise."
    return result


def runtime_policy_enforcement_profile():
    cid="PRC-28-032";ordinal=1;pre=prefix(cid,ordinal)
    policies=f"{pre}.effective_runtime_policy_identities_by_workload";states=f"{pre}.runtime_policy_enforcement_states_by_workload";required=f"{pre}.required_workload_ids";approved_policies=f"{pre}.approved_runtime_policy_identities";approved_states=f"{pre}.approved_enforced_state_ids"
    facts=[fact(policies,"identity_map","Direct effective runtime-policy identity attached to each workload by the control plane.","environment","Read-only effective workload-to-policy binding from the control plane."),fact(states,"state_map","Direct runtime-policy enforcement state for each workload, separate from policy attachment.","environment","Read-only effective state from the runtime policy controller.")]
    params=[parameter(required,"string_set","Complete workload identities sealed independently from the deployment inventory."),parameter(approved_policies,"string_set","Independently approved runtime-policy identities."),parameter(approved_states,"string_set","Direct raw controller states that independently mean enforcement is active, not merely attached or configured.")]
    args=[{"op":"map_keys_eq_set_parameter","fact":x,"parameter":required} for x in [policies,states]]+[{"op":"identity_map_values_in_parameter","fact":policies,"parameter":approved_policies},{"op":"state_map_values_in_parameter","fact":states,"parameter":approved_states}]
    ids=["workload-A","workload-B"];policy_values={x:"policy-A" for x in ids};state_values={x:"enforced" for x in ids}
    pf={policies:fvalue("identity_map",policy_values),states:fvalue("state_map",state_values)};ff=copy.deepcopy(pf);ff[states]=fvalue("state_map",{"workload-A":"attached-not-enforced","workload-B":"enforced"});pv={required:pvalue("string_set",ids),approved_policies:pvalue("string_set",["policy-A"]),approved_states:pvalue("string_set",["enforced"])}
    result=exact_profile(facts,params,{"op":"all","args":args},pf,ff,pv,"Every inventoried workload has an approved effective policy identity and a separate raw enforced state.","One workload has the approved policy attached but its direct enforcement state is not active.")
    result["fixtures"]["counterexample"]=copy.deepcopy(result["fixtures"]["fail"]);result["fixtures"]["counterexample"]["description"]="A policy is attached to every workload, but one controller state says attached-not-enforced; attachment alone cannot pass."
    return result


def pipeline_monitor_profile():
    cid="PRC-29-030";ordinal=1;pre=prefix(cid,ordinal);required=f"{pre}.required_pipeline_ids";allowed=f"{pre}.approved_active_monitor_state_ids"
    backlog_ids=f"{pre}.backlog_monitor_identities_by_pipeline";quota_ids=f"{pre}.quota_monitor_identities_by_pipeline";backlog_states=f"{pre}.backlog_monitor_states_by_pipeline";quota_states=f"{pre}.quota_monitor_states_by_pipeline";expected_backlog=f"{pre}.expected_backlog_monitor_identities_by_pipeline";expected_quota=f"{pre}.expected_quota_monitor_identities_by_pipeline"
    facts=[fact(backlog_ids,"identity_map","Direct backlog-depth monitor identity bound to each pipeline.","environment","Read-only effective telemetry monitor bindings."),fact(quota_ids,"identity_map","Direct remaining-quota or quota-exhaustion monitor identity bound to each pipeline.","environment","Read-only effective telemetry monitor bindings."),fact(backlog_states,"state_map","Direct effective state of each bound backlog-depth monitor.","environment","Read-only monitor runtime or effective configuration state."),fact(quota_states,"state_map","Direct effective state of each bound quota monitor.","environment","Read-only monitor runtime or effective configuration state.")]
    params=[parameter(required,"string_set","Complete telemetry-pipeline identities sealed independently before monitor collection."),parameter(expected_backlog,"identity_map","Independently approved backlog-monitor identity for every sealed pipeline."),parameter(expected_quota,"identity_map","Independently approved quota-monitor identity for every sealed pipeline."),parameter(allowed,"string_set","Independently approved direct raw states meaning the bound monitor is active.")]
    args=[{"op":"map_keys_eq_set_parameter","fact":x,"parameter":required} for x in [backlog_ids,quota_ids,backlog_states,quota_states]]+[{"op":"identity_map_eq_parameter","fact":backlog_ids,"parameter":expected_backlog},{"op":"identity_map_eq_parameter","fact":quota_ids,"parameter":expected_quota},{"op":"state_map_values_in_parameter","fact":backlog_states,"parameter":allowed},{"op":"state_map_values_in_parameter","fact":quota_states,"parameter":allowed}]
    ids=["pipeline-A","pipeline-B"];bi={"pipeline-A":"backlog-A","pipeline-B":"backlog-B"};qi={"pipeline-A":"quota-A","pipeline-B":"quota-B"};active={x:"active" for x in ids}
    pf={backlog_ids:fvalue("identity_map",bi),quota_ids:fvalue("identity_map",qi),backlog_states:fvalue("state_map",active),quota_states:fvalue("state_map",active)};ff=copy.deepcopy(pf);ff[quota_states]=fvalue("state_map",{"pipeline-A":"inactive","pipeline-B":"active"});pv={required:pvalue("string_set",ids),expected_backlog:pvalue("identity_map",bi),expected_quota:pvalue("identity_map",qi),allowed:pvalue("string_set",["active"])}
    result=exact_profile(facts,params,{"op":"all","args":args},pf,ff,pv,"Every sealed pipeline has the expected backlog and quota monitors and both direct monitor states are active.","One expected quota monitor is bound but its direct raw state is inactive.")
    result["fixtures"]["counterexample"]=copy.deepcopy(result["fixtures"]["fail"]);result["fixtures"]["counterexample"]["description"]="Every pipeline has monitor records, but one quota monitor is inactive; record presence cannot be treated as active monitoring."
    return result


def resource_label_profile():
    cid="PRC-28-005";ordinal=1;pre=prefix(cid,ordinal);owners=f"{pre}.owner_label_values_by_resource";criticality=f"{pre}.criticality_label_values_by_resource";required=f"{pre}.required_resource_ids";allowed=f"{pre}.approved_criticality_label_values"
    facts=[fact(owners,"string_map","Direct effective owner-label text value for every resource.","environment","Read-only effective resource label map."),fact(criticality,"state_map","Direct effective criticality-label value for every resource.","environment","Read-only effective resource label map.")]
    params=[parameter(required,"string_set","Complete resource identities sealed independently from deployed infrastructure inventory."),parameter(allowed,"string_set","Independently approved nonempty criticality label values.")]
    args=[{"op":"map_keys_eq_set_parameter","fact":owners,"parameter":required},{"op":"map_keys_eq_set_parameter","fact":criticality,"parameter":required},{"op":"string_map_all_nonempty","fact":owners},{"op":"state_map_values_in_parameter","fact":criticality,"parameter":allowed}]
    ids=["resource-A","resource-B"];owner_values={"resource-A":"team-A","resource-B":"team-B"};criticality_values={"resource-A":"critical","resource-B":"standard"}
    pf={owners:fvalue("string_map",owner_values),criticality:fvalue("state_map",criticality_values)};ff=copy.deepcopy(pf);ff[owners]=fvalue("string_map",{"resource-A":"","resource-B":"team-B"});pv={required:pvalue("string_set",ids),allowed:pvalue("string_set",["critical","standard"])}
    result=exact_profile(facts,params,{"op":"all","args":args},pf,ff,pv,"Every inventoried resource has a nonempty direct owner label and an independently approved criticality label.","One resource is present in both label maps but its owner label is empty.")
    result["fixtures"]["counterexample"]=copy.deepcopy(result["fixtures"]["fail"]);result["fixtures"]["counterexample"]["description"]="The resource is counted as owner-labeled by a selector, but its direct owner-label value is empty; direct raw value checking rejects this false pass."
    return result


def release_signature_profile():
    cid="USEQ-E8367B76";ordinal=1;pre=prefix(cid,ordinal)
    verifier=f"{pre}.trust_policy_verifier_identity";policy=f"{pre}.trust_policy_digest";inputs=f"{pre}.verifier_input_artifact_digests";recomputed=f"{pre}.recomputed_release_artifact_digests";states=f"{pre}.signature_or_attestation_verification_states";verified_at=f"{pre}.verification_success_times";deployed_at=f"{pre}.deployment_or_distribution_times"
    required=f"{pre}.required_release_artifact_ids";expected_verifier=f"{pre}.required_trust_policy_verifier_identity";expected_policy=f"{pre}.required_trust_policy_digest";allowed=f"{pre}.approved_successful_verification_state_ids"
    facts=[fact(verifier,"identity","Direct deterministic signature or attestation verifier identity and version.","environment","Authenticated verifier invocation record."),fact(policy,"digest","Direct digest of the selected trust-policy bytes used by the verifier.","environment","Exact verifier policy input bytes."),fact(inputs,"digest_map","Direct artifact digest input to the verifier for every release artifact.","environment","Authenticated verifier input manifest."),fact(recomputed,"digest_map","Locally recomputed digest of the exact release artifact bytes deployed or distributed.","environment","Exact release artifact bytes hashed locally."),fact(states,"state_map","Direct verifier terminal output state for every release artifact.","environment","Authenticated direct output of the pinned verifier."),fact(verified_at,"time_map","Direct successful verification timestamp for every release artifact.","environment","Authenticated verifier completion events."),fact(deployed_at,"time_map","Direct deployment or distribution timestamp for every release artifact.","environment","Authenticated deployment or distribution events.")]
    params=[parameter(required,"string_set","Complete release-artifact identities sealed independently from the effective gate inventory."),parameter(expected_verifier,"identity","Independently approved exact deterministic verifier identity and version."),parameter(expected_policy,"digest","Independently selected exact trust-policy digest."),parameter(allowed,"string_set","Direct verifier output states independently approved to mean successful signature or attestation validation.")]
    args=[{"op":"identity_eq_parameter","fact":verifier,"parameter":expected_verifier},{"op":"digest_eq_parameter","fact":policy,"parameter":expected_policy}]+[{"op":"map_keys_eq_set_parameter","fact":x,"parameter":required} for x in [inputs,recomputed,states,verified_at,deployed_at]]+[{"op":"digest_map_eq_fact","fact":inputs,"other_fact":recomputed},{"op":"state_map_values_in_parameter","fact":states,"parameter":allowed},{"op":"time_map_before_fact","fact":verified_at,"other_fact":deployed_at}]
    ids=["artifact-A","artifact-B"];digests={"artifact-A":"a"*64,"artifact-B":"b"*64};state_values={x:"verified" for x in ids};verify_times={"artifact-A":"2026-08-29T10:00:00Z","artifact-B":"2026-08-29T11:00:00Z"};deploy_times={"artifact-A":"2026-08-29T10:05:00Z","artifact-B":"2026-08-29T11:05:00Z"}
    pf={verifier:fvalue("identity","verifier-A@1"),policy:fvalue("digest","c"*64),inputs:fvalue("digest_map",digests),recomputed:fvalue("digest_map",digests),states:fvalue("state_map",state_values),verified_at:fvalue("time_map",verify_times),deployed_at:fvalue("time_map",deploy_times)};ff=copy.deepcopy(pf);ff[recomputed]=fvalue("digest_map",{"artifact-A":"d"*64,"artifact-B":"b"*64});pv={required:pvalue("string_set",ids),expected_verifier:pvalue("identity","verifier-A@1"),expected_policy:pvalue("digest","c"*64),allowed:pvalue("string_set",["verified"])}
    result=exact_profile(facts,params,{"op":"all","args":args},pf,ff,pv,"Every release artifact is verified by the selected pinned trust policy over its exact recomputed bytes before deployment or distribution.","The verifier reports success, but one verifier-input digest differs from the locally recomputed digest of the bytes actually deployed.")
    result["fixtures"]["counterexample"]=copy.deepcopy(result["fixtures"]["fail"]);result["fixtures"]["counterexample"]["description"]="A successful signature state and timely verification exist, but they cover different bytes for artifact-A; verifier success without exact byte binding cannot pass."
    return result


def exception_gate_profile():
    cid="USEQ-CF1B0C98";ordinal=1;pre=prefix(cid,ordinal);records=f"{pre}.required_exception_record_ids";cases=f"{pre}.required_invalid_exception_case_ids"
    field_specs=[("affected_data_identities","identity_map"),("threat_horizon_identities","identity_map"),("compensating_control_identities","identity_map"),("migration_date_times","time_map"),("owner_identities","identity_map")]
    negative_ids=["missing-affected-data","missing-threat-horizon","missing-compensating-control","missing-migration-date","missing-owner","missing-expiry","expiry-equals-issue","expiry-before-issue"]
    facts=[];params=[parameter(records,"string_set","Complete effective exception-record identities sealed independently from the gate inventory."),parameter(cases,"string_set","The exact eight independently versioned negative cases: one omission for each named field plus expiry equal to issue and expiry before issue.")];args=[];pf={};ff={};pv={records:pvalue("string_set",["exception-A","exception-B"]),cases:pvalue("string_set",negative_ids)}
    for index,(suffix,typ) in enumerate(field_specs):
        key=f"{pre}.{suffix}";facts.append(fact(key,typ,f"Direct submitted nonempty {suffix.replace('_',' ')} for every exception record.","environment","Effective release-gate exception record bytes."));args.append({"op":"map_keys_eq_set_parameter","fact":key,"parameter":records})
        if typ=="time_map": values={"exception-A":"2026-09-01T00:00:00Z","exception-B":"2026-09-02T00:00:00Z"}
        else: values={"exception-A":f"{suffix}-A","exception-B":f"{suffix}-B"}
        pf[key]=fvalue(typ,values);ff[key]=fvalue(typ,({"exception-A":values["exception-A"]} if index==0 else values))
    issue=f"{pre}.exception_issue_times";expiry=f"{pre}.exception_expiry_times";case_inputs=f"{pre}.invalid_case_input_digests";case_states=f"{pre}.invalid_case_gate_states";expected_inputs=f"{pre}.expected_invalid_case_input_digests";expected_states=f"{pre}.expected_invalid_case_gate_states"
    facts += [fact(issue,"time_map","Direct authenticated issue timestamp for every exception record.","environment","Effective exception record timestamps."),fact(expiry,"time_map","Direct authenticated expiry timestamp for every exception record.","environment","Effective exception record timestamps."),fact(case_inputs,"digest_map","Direct digest of each missing-field or expired negative gate fixture.","environment","Authenticated bounded gate invocation input manifest."),fact(case_states,"state_map","Direct terminal gate output for each missing-field or expired fixture.","environment","Authenticated bounded gate output record.")]
    params += [parameter(expected_inputs,"digest_map","Exact independently approved input digest for every sealed negative exception fixture."),parameter(expected_states,"state_map","Exact independently approved rejected state for every sealed negative exception fixture.")]
    args += [{"op":"map_keys_eq_set_parameter","fact":issue,"parameter":records},{"op":"map_keys_eq_set_parameter","fact":expiry,"parameter":records},{"op":"time_map_before_fact","fact":issue,"other_fact":expiry},{"op":"map_keys_eq_set_parameter","fact":case_inputs,"parameter":cases},{"op":"map_keys_eq_set_parameter","fact":case_states,"parameter":cases},{"op":"digest_map_eq_parameter","fact":case_inputs,"parameter":expected_inputs},{"op":"state_map_eq_parameter","fact":case_states,"parameter":expected_states}]
    issue_values={"exception-A":"2026-08-29T10:00:00Z","exception-B":"2026-08-29T11:00:00Z"};expiry_values={"exception-A":"2026-09-29T10:00:00Z","exception-B":"2026-09-29T11:00:00Z"};input_values={case:format(index+1,"x")*64 for index,case in enumerate(negative_ids)};state_values={case:"rejected" for case in negative_ids}
    for target,value in [(issue,issue_values),(expiry,expiry_values),(case_inputs,input_values),(case_states,state_values)]: pf[target]=fvalue("time_map" if target in [issue,expiry] else ("digest_map" if target==case_inputs else "state_map"),value);ff[target]=copy.deepcopy(pf[target])
    pv[expected_inputs]=pvalue("digest_map",input_values);pv[expected_states]=pvalue("state_map",state_values)
    result=exact_profile(facts,params,{"op":"all","args":args},pf,ff,pv,"Every exception record contains each named field with expiry after issue, and every sealed missing-field or expired fixture is rejected.","One exception record omits affected-data while valid timestamps and all negative gate cases still look correct.")
    result["fixtures"]["counterexample"]=copy.deepcopy(result["fixtures"]["fail"]);result["fixtures"]["counterexample"]["description"]="The gate rejects the sampled invalid fixtures, but exception-B is absent from the affected-data map; sampled rejection cannot prove complete record-field coverage."
    return result


def sbom_exact_profile():
    cid="PRC-09-005";ordinal=1;pre=prefix(cid,ordinal);required=f"{pre}.required_component_ids"
    pairs=[
        ("sbom_recorded_component_versions","independently_resolved_component_versions","identity_map","SBOM recorded version identity","version identity resolved independently from the packaged component metadata"),
        ("sbom_recorded_package_identifiers","independently_resolved_package_identifiers","identity_map","SBOM recorded package identifier","package identifier resolved independently from the packaged component metadata"),
        ("sbom_recorded_component_sources","independently_resolved_component_sources","identity_map","SBOM recorded source identity","source identity resolved independently from the lock, provenance, or packaged component metadata"),
        ("sbom_recorded_component_hashes","locally_recomputed_component_byte_digests","digest_map","SBOM recorded component digest","locally recomputed digest of the exact packaged component bytes"),
    ]
    facts=[];args=[];pf={};ff={};ids=["component-A","component-B"]
    for index,(left_suffix,right_suffix,typ,left_meaning,right_meaning) in enumerate(pairs):
        left=f"{pre}.{left_suffix}";right=f"{pre}.{right_suffix}"
        facts += [fact(left,typ,f"Direct {left_meaning} for every component.","artifact","Lossless parsed exact SBOM artifact bytes."),fact(right,typ,f"Direct {right_meaning} for every component.","artifact","Independently parsed package, lock, provenance, or component bytes from the same exact release artifact.")]
        args += [{"op":"map_keys_eq_set_parameter","fact":left,"parameter":required},{"op":"map_keys_eq_set_parameter","fact":right,"parameter":required},{"op":f"{'digest_map' if typ=='digest_map' else 'identity_map'}_eq_fact","fact":left,"other_fact":right}]
        if typ=="digest_map": values={"component-A":"a"*64,"component-B":"b"*64};wrong={"component-A":"c"*64,"component-B":"b"*64}
        else: values={"component-A":f"{left_suffix}-A","component-B":f"{left_suffix}-B"};wrong={"component-A":"wrong-value","component-B":values["component-B"]}
        pf[left]=fvalue(typ,values);pf[right]=fvalue(typ,values);ff[left]=fvalue(typ,values);ff[right]=fvalue(typ,(wrong if index==0 else values))
    params=[parameter(required,"string_set","Complete component identities sealed independently from the exact packaged-artifact inventory before SBOM comparison.")]
    result=exact_profile(facts,params,{"op":"all","args":args},pf,ff,{required:pvalue("string_set",ids)},"Every SBOM field agrees component by component with independently derived package metadata and locally hashed bytes.","The SBOM has every required field, but component-A's recorded version disagrees with the independently resolved packaged version.")
    result["fixtures"]["counterexample"]=copy.deepcopy(result["fixtures"]["fail"]);result["fixtures"]["counterexample"]["description"]="A complete-looking SBOM lists version, package, source, and hash for both components, but component-A's version is wrong relative to the artifact; field presence alone cannot pass."
    return result


def required_load_capacity_profile():
    cid="PRC-33-023";ordinal=1;pre=prefix(cid,ordinal);observed=f"{pre}.tested_capacity_by_subject_load_dimension";units=f"{pre}.tested_capacity_unit_identities";windows=f"{pre}.tested_capacity_window_identities";required=f"{pre}.required_subject_load_dimension_keys";minimum=f"{pre}.approved_required_load_by_key";expected_units=f"{pre}.approved_unit_identities_by_key";expected_windows=f"{pre}.approved_window_identities_by_key"
    facts=[fact(observed,"number_map","Direct tested disaster-recovery capacity for every exact subject x required-load-dimension key.","executed","Authenticated capacity-test raw output."),fact(units,"identity_map","Direct unit identity for each tested load dimension.","executed","Authenticated measurement metadata."),fact(windows,"identity_map","Direct test-window identity for each tested load dimension.","executed","Authenticated execution-window metadata.")]
    params=[parameter(required,"string_set","Complete subject x project-required load-dimension keys sealed independently from the recovery load profile; dimensions are project-defined, never hardcoded to request rate."),parameter(minimum,"number_map","Independently approved minimum required load for every exact sealed key in matching canonical units."),parameter(expected_units,"identity_map","Approved unit identity for every sealed subject-load key."),parameter(expected_windows,"identity_map","Approved test-window identity for every sealed subject-load key.")]
    keys=["recovery-site-A|request_rate","recovery-site-A|concurrent_sessions","recovery-site-A|throughput_bytes_per_second"];limits={key:10 for key in keys};values={key:12 for key in keys};unit_values={key:"approved-unit" for key in keys};window_values={key:"approved-window" for key in keys}
    args=[{"op":"map_keys_eq_set_parameter","fact":x,"parameter":required} for x in [observed,units,windows]]+[{"op":"number_map_gte_parameter","fact":observed,"parameter":minimum},{"op":"identity_map_eq_parameter","fact":units,"parameter":expected_units},{"op":"identity_map_eq_parameter","fact":windows,"parameter":expected_windows}]
    pf={observed:fvalue("number_map",values),units:fvalue("identity_map",unit_values),windows:fvalue("identity_map",window_values)};ff=copy.deepcopy(pf);broken=dict(values);broken[keys[-1]]=9;ff[observed]=fvalue("number_map",broken);pv={required:pvalue("string_set",keys),minimum:pvalue("number_map",limits),expected_units:pvalue("identity_map",unit_values),expected_windows:pvalue("identity_map",window_values)}
    result=exact_profile(facts,params,{"op":"all","args":args},pf,ff,pv,"Tested disaster-recovery capacity meets every independently declared load dimension in its approved unit and window.","Request rate and session capacity pass, but required byte throughput is below its sealed required load.")
    result["fixtures"]["counterexample"]=copy.deepcopy(result["fixtures"]["fail"]);result["fixtures"]["counterexample"]["description"]="A scanner that checks requests per second alone would pass while byte throughput is insufficient; exact sealed load-dimension coverage rejects that result."
    return result


def deployed_reviewed_tested_profile():
    cid="USEQ-4C2E2D56";ordinal=1;pre=prefix(cid,ordinal)
    names=["deployed_artifact_digest","reviewed_artifact_digest","tested_artifact_digest","deployed_configuration_digest","reviewed_configuration_digest","tested_configuration_digest"]
    facts=[fact(f"{pre}.{name}","digest",f"Direct {name.replace('_',' ')} from its exact authoritative bytes or immutable record.","artifact","Exact deployed bytes, reviewed-change record, or authenticated complete test-result binding named by the field.") for name in names]
    da,ra,ta,dc,rc,tc=[f"{pre}.{name}" for name in names]
    args=[{"op":"digest_eq_fact","fact":da,"other_fact":ra},{"op":"digest_eq_fact","fact":da,"other_fact":ta},{"op":"digest_eq_fact","fact":dc,"other_fact":rc},{"op":"digest_eq_fact","fact":dc,"other_fact":tc}]
    pf={key:fvalue("digest",("a"*64 if "artifact" in key else "b"*64)) for key in [da,ra,ta,dc,rc,tc]};ff=copy.deepcopy(pf);ff[ta]=fvalue("digest","c"*64)
    result=exact_profile(facts,[],{"op":"all","args":args},pf,ff,{},"One deployed artifact digest is reused in both reviewed and tested comparisons, and one deployed configuration digest is reused in both comparisons.","The reviewed digest matches the deployed artifact, but the complete test results bind a different artifact digest.")
    result["fixtures"]["counterexample"]=copy.deepcopy(result["fixtures"]["fail"]);result["fixtures"]["counterexample"]["description"]="A split-brain pair of provider-supplied deployed-digest copies could make separate comparisons pass; this predicate references the same deployed artifact fact in both comparisons and rejects the mismatched tested digest."
    return result


def merged_review_verification_profile():
    cid="USEQ-CA466120";ordinal=1;pre=prefix(cid,ordinal);required=f"{pre}.required_post_merge_verification_ids"
    fields=[("locally_recomputed_merged_artifact_digests","digest_map"),("reviewed_artifact_digests_by_verification","digest_map"),("verification_bound_artifact_digests","digest_map"),("locally_recomputed_effective_configuration_digests","digest_map"),("reviewed_configuration_digests_by_verification","digest_map"),("verification_bound_configuration_digests","digest_map")]
    facts=[fact(f"{pre}.{suffix}",typ,f"Direct {suffix.replace('_',' ')} for every post-merge verification result.","artifact","Exact merged bytes, immutable reviewed-change binding, or authenticated verification-result record keyed by the sealed verification identity.") for suffix,typ in fields]
    keys=[f"{pre}.{suffix}" for suffix,_ in fields];aa,ar,av,ca,cr,cv=keys
    args=[{"op":"map_keys_eq_set_parameter","fact":x,"parameter":required} for x in keys]+[{"op":"digest_map_eq_fact","fact":aa,"other_fact":ar},{"op":"digest_map_eq_fact","fact":aa,"other_fact":av},{"op":"digest_map_eq_fact","fact":ca,"other_fact":cr},{"op":"digest_map_eq_fact","fact":ca,"other_fact":cv}]
    ids=["verification-A","verification-B"];artifact_values={x:"a"*64 for x in ids};config_values={x:"b"*64 for x in ids};pf={aa:fvalue("digest_map",artifact_values),ar:fvalue("digest_map",artifact_values),av:fvalue("digest_map",artifact_values),ca:fvalue("digest_map",config_values),cr:fvalue("digest_map",config_values),cv:fvalue("digest_map",config_values)};ff=copy.deepcopy(pf);ff[av]=fvalue("digest_map",{"verification-A":"c"*64,"verification-B":"a"*64})
    result=exact_profile(facts,[parameter(required,"string_set","Complete post-merge verification-result identities sealed independently from the merge workflow inventory.")],{"op":"all","args":args},pf,ff,{required:pvalue("string_set",ids)},"For each post-merge verification, locally recomputed merged bytes, reviewed bindings, and verification-result bindings use the same artifact and configuration digest maps.","One verification result binds a different artifact digest although the reviewed and locally recomputed merged digests agree.")
    result["fixtures"]["counterexample"]=copy.deepcopy(result["fixtures"]["fail"]);result["fixtures"]["counterexample"]["description"]="Separate copied merged-digest fields could agree with different counterparts; comparing all reviewed and verification maps to the same locally recomputed merged map rejects this split binding."
    return result


def branch_policy_profile(control_id, ordinal, kind):
    pre=prefix(control_id,ordinal)
    if kind=="minimum_reviews":
        observed=f"{pre}.effective_minimum_approving_review_counts_by_branch";expected=f"{pre}.approved_minimum_approving_review_counts_by_branch"
        facts=[fact(observed,"number_map","Direct numeric minimum-approving-review field from each effective protected-branch rule.","environment","Read-only effective branch-rule snapshot.")];params=[parameter(expected,"number_map","Independently approved minimum review count keyed by every sealed protected branch; exact keys define the branch inventory.")];good={"branch-main":1,"branch-release":2};bad={"branch-main":0,"branch-release":2};pred={"op":"number_map_gte_parameter","fact":observed,"parameter":expected};pv={expected:pvalue("number_map",good)}
    elif kind=="bypass_denied":
        observed=f"{pre}.effective_review_bypass_permission_states_by_branch_principal";required=f"{pre}.required_branch_principal_relation_ids";allowed=f"{pre}.approved_denied_permission_state_ids"
        facts=[fact(observed,"state_map","Direct effective review-bypass permission state for every protected-branch and assessed-principal relation.","environment","Read-only effective branch permission matrix.")];params=[parameter(required,"string_set","Complete protected-branch x assessed-principal relation identities sealed independently."),parameter(allowed,"string_set","Raw permission state identifiers independently approved to mean bypass is denied.")];good={"branch-main|principal-A":"denied","branch-release|principal-B":"denied"};bad={"branch-main|principal-A":"allowed","branch-release|principal-B":"denied"};pred={"op":"all","args":[{"op":"map_keys_eq_set_parameter","fact":observed,"parameter":required},{"op":"state_map_values_in_parameter","fact":observed,"parameter":allowed}]};pv={required:pvalue("string_set",list(good)),allowed:pvalue("string_set",["denied"])}
    elif kind=="mandatory_checks":
        observed=f"{pre}.effective_required_status_context_identities_by_branch_check";expected=f"{pre}.approved_required_status_context_identities_by_branch_check"
        facts=[fact(observed,"identity_map","Direct required-status-context identity from the effective branch rule for each sealed branch and policy-required check relation.","environment","Read-only effective required-status-check bindings from branch rules.")];params=[parameter(expected,"identity_map","Exact independently approved status-context identity keyed by the complete branch x required-check relation inventory.")];good={"branch-main|check-A":"context-A","branch-release|check-B":"context-B"};bad={"branch-main|check-A":"wrong-context","branch-release|check-B":"context-B"};pred={"op":"identity_map_eq_parameter","fact":observed,"parameter":expected};pv={expected:pvalue("identity_map",good)}
    elif kind=="rewrite_denied":
        observed=f"{pre}.effective_rewrite_permission_states_by_namespace_principal";required=f"{pre}.required_namespace_principal_relation_ids";allowed=f"{pre}.approved_denied_permission_state_ids"
        facts=[fact(observed,"state_map","Direct effective rewrite permission state for each protected branch or release-tag namespace and unauthorized-principal relation.","environment","Read-only effective namespace permission matrix.")];params=[parameter(required,"string_set","Complete protected namespace x unauthorized-principal relation identities sealed independently."),parameter(allowed,"string_set","Raw permission state identifiers independently approved to mean rewrite is denied.")];good={"branch-main|principal-X":"denied","tag-release|principal-X":"denied"};bad={"branch-main|principal-X":"allowed","tag-release|principal-X":"denied"};pred={"op":"all","args":[{"op":"map_keys_eq_set_parameter","fact":observed,"parameter":required},{"op":"state_map_values_in_parameter","fact":observed,"parameter":allowed}]};pv={required:pvalue("string_set",list(good)),allowed:pvalue("string_set",["denied"])}
    else: raise ValueError(kind)
    pf={observed:fvalue("number_map" if kind=="minimum_reviews" else ("identity_map" if kind=="mandatory_checks" else "state_map"),good)};ff={observed:fvalue("number_map" if kind=="minimum_reviews" else ("identity_map" if kind=="mandatory_checks" else "state_map"),bad)}
    return exact_profile(facts,params,pred,pf,ff,pv,f"The effective branch-policy raw field satisfies the independently sealed {kind.replace('_',' ')} relation for every exact branch, check, namespace, or principal key.",f"One effective branch-policy field has a permissive or wrong value for a sealed relation even though other branches remain protected.")


def protected_change_policy_profile():
    cid="USEQ-0B41367E";ordinal=1;pre=prefix(cid,ordinal);reviews=f"{pre}.effective_minimum_approving_review_counts_by_branch";checks=f"{pre}.effective_required_check_set_digests_by_branch";required=f"{pre}.required_protected_branch_ids";minimum=f"{pre}.approved_minimum_review_counts_by_branch";expected_checks=f"{pre}.approved_required_check_set_digests_by_branch"
    facts=[fact(reviews,"number_map","Direct numeric minimum-approving-review field for every protected branch.","environment","Read-only effective protected-branch rules."),fact(checks,"digest_map","Canonical digest of the direct effective required-check identity set for every protected branch.","environment","Read-only effective required-status-check lists canonicalized by the pinned scanner.")]
    params=[parameter(required,"string_set","Complete protected-branch identities sealed independently."),parameter(minimum,"number_map","Approved minimum review count for every sealed branch."),parameter(expected_checks,"digest_map","Approved nonempty required-check set digest for every sealed branch.")]
    ids=["branch-main","branch-release"];minimum_values={x:1 for x in ids};check_values={"branch-main":"a"*64,"branch-release":"b"*64};review_values={x:2 for x in ids}
    args=[{"op":"map_keys_eq_set_parameter","fact":reviews,"parameter":required},{"op":"map_keys_eq_set_parameter","fact":checks,"parameter":required},{"op":"number_map_gte_parameter","fact":reviews,"parameter":minimum},{"op":"digest_map_eq_parameter","fact":checks,"parameter":expected_checks}]
    pf={reviews:fvalue("number_map",review_values),checks:fvalue("digest_map",check_values)};ff=copy.deepcopy(pf);ff[reviews]=fvalue("number_map",{"branch-main":0,"branch-release":2});pv={required:pvalue("string_set",ids),minimum:pvalue("number_map",minimum_values),expected_checks:pvalue("digest_map",check_values)}
    return exact_profile(facts,params,{"op":"all","args":args},pf,ff,pv,"Every protected branch has a sufficient direct approval-count requirement and the exact approved required-check set.","The required-check set is correct, but branch-main's effective minimum approval count is zero.")


def separate_release_authorities_profile():
    cid="USEQ-517E9DBD";ordinal=1;pre=prefix(cid,ordinal);creators=f"{pre}.artifact_creator_identities_by_release";authorizers=f"{pre}.deployment_authorizer_identities_by_release";required=f"{pre}.required_release_ids"
    facts=[fact(creators,"identity_map","Direct authenticated artifact creator identity for every release.","environment","Authenticated artifact creation records."),fact(authorizers,"identity_map","Direct authenticated deployment authorization identity for every release.","environment","Authenticated deployment authorization records collected independently from artifact creation.")]
    params=[parameter(required,"string_set","Complete release identities sealed independently from the release and deployment inventory.")]
    ids=["release-A","release-B"];creator_values={"release-A":"builder-A","release-B":"builder-B"};authorizer_values={"release-A":"approver-A","release-B":"approver-B"}
    args=[{"op":"map_keys_eq_set_parameter","fact":creators,"parameter":required},{"op":"map_keys_eq_set_parameter","fact":authorizers,"parameter":required},{"op":"identity_map_values_not_equal_fact","fact":creators,"other_fact":authorizers}]
    pf={creators:fvalue("identity_map",creator_values),authorizers:fvalue("identity_map",authorizer_values)};ff=copy.deepcopy(pf);ff[authorizers]=fvalue("identity_map",{"release-A":"builder-A","release-B":"approver-B"})
    return exact_profile(facts,params,{"op":"all","args":args},pf,ff,{required:pvalue("string_set",ids)},"For each release, the authenticated artifact creator and deployment authorizer identities are both present and distinct.","release-A was created and authorized by builder-A, so the two duties were not separated even though release-B is separate.")


def policy_negative_cases_profile(control_id, failure_suffix, other_input_suffix, case_ids, failure_semantics, other_semantics):
    fields=[("policy_case_input_digests","digest_map","Direct digest of each independently sealed policy truth-table fixture."),(failure_suffix,"boolean_map",failure_semantics),(other_input_suffix,"digest_map",other_semantics),("aggregate_failure_states","boolean_map","Direct aggregate failed-decision output from the pinned policy evaluator for every sealed fixture.")]
    return combine_profiles(
        pinned_tool_profile(control_id,1,"repository","policy_evaluator_identity","policy_evaluator_configuration_digest","Direct pinned deterministic policy-expression evaluator identity and version.","Direct digest of the exact parsed aggregation policy and evaluator rules."),
        expected_field_maps_profile(control_id,1,"repository","policy_case_ids",fields,"Complete independently generated negative truth-table cases sealed before policy evaluation.",case_ids),
    )


def privileged_approval_profile():
    cid="PRC-28-030";ordinal=1;pre=prefix(cid,ordinal);states=f"{pre}.effective_privileged_or_host_access_states";approvals=f"{pre}.approval_record_digests_by_privileged_subject";configs=f"{pre}.approved_configuration_digests_by_privileged_subject";all_subjects=f"{pre}.required_container_subject_ids";privileged=f"{pre}.required_privileged_subject_ids";expected_states=f"{pre}.approved_effective_access_states";expected_approvals=f"{pre}.approved_exception_record_digests";expected_configs=f"{pre}.approved_privileged_configuration_digests"
    facts=[fact(states,"state_map","Direct effective privileged-container or host-access state for every container subject.","environment","Read-only effective container runtime configuration."),fact(approvals,"digest_map","Direct digest of the independently authenticated explicit approval record for every privileged subject.","environment","Protected approval-record store, collected separately from runtime configuration."),fact(configs,"digest_map","Locally recomputed digest of the exact privileged runtime configuration named by each approval.","environment","Exact effective runtime configuration bytes bound by the approval record.")]
    params=[parameter(all_subjects,"string_set","Complete container subject identities sealed from deployment inventory."),parameter(privileged,"string_set","Exact privileged subject identities sealed independently from reviewed runtime policy."),parameter(expected_states,"state_map","Approved direct effective access state for every container subject."),parameter(expected_approvals,"digest_map","Exact independently approved exception-record digest for every privileged subject."),parameter(expected_configs,"digest_map","Exact approved privileged configuration digest for every privileged subject.")]
    all_ids=["container-A","container-B"];priv_ids=["container-A"];state_values={"container-A":"explicitly-approved-privileged","container-B":"unprivileged"};approval_values={"container-A":"a"*64};config_values={"container-A":"b"*64}
    args=[{"op":"map_keys_eq_set_parameter","fact":states,"parameter":all_subjects},{"op":"state_map_eq_parameter","fact":states,"parameter":expected_states},{"op":"map_keys_eq_set_parameter","fact":approvals,"parameter":privileged},{"op":"map_keys_eq_set_parameter","fact":configs,"parameter":privileged},{"op":"digest_map_eq_parameter","fact":approvals,"parameter":expected_approvals},{"op":"digest_map_eq_parameter","fact":configs,"parameter":expected_configs}]
    pf={states:fvalue("state_map",state_values),approvals:fvalue("digest_map",approval_values),configs:fvalue("digest_map",config_values)};ff=copy.deepcopy(pf);ff[configs]=fvalue("digest_map",{"container-A":"c"*64});pv={all_subjects:pvalue("string_set",all_ids),privileged:pvalue("string_set",priv_ids),expected_states:pvalue("state_map",state_values),expected_approvals:pvalue("digest_map",approval_values),expected_configs:pvalue("digest_map",config_values)}
    return exact_profile(facts,params,{"op":"all","args":args},pf,ff,pv,"Runtime state, independent approval record, and exact approved configuration digest agree for every privileged subject; all other subjects are unprivileged.","The approval record exists, but it names a different privileged configuration digest than the one effectively deployed.")


def interface_binding_profile():
    cid="USEQ-B086FE62";ordinal=1;pre=prefix(cid,ordinal);required=f"{pre}.required_interface_field_ids"
    maps=[
        ("schema_identities","identity_map",{"field-text":"schema-A","field-duration":"schema-A"}),
        ("media_type_identities","identity_map",{"field-text":"application/json","field-duration":"application/json"}),
        ("encoding_identities","identity_map",{"field-text":"utf-8","field-duration":"utf-8"}),
        ("language_tag_or_explicit_na_values","string_map",{"field-text":"en-US","field-duration":"NA"}),
        ("unit_or_explicit_na_values","string_map",{"field-text":"NA","field-duration":"seconds"}),
        ("stable_identifier_values","identity_map",{"field-text":"field-id-text","field-duration":"field-id-duration"}),
        ("time_semantics_identities","identity_map",{"field-text":"not-time","field-duration":"elapsed-duration"}),
    ]
    facts=[];params=[parameter(required,"string_set","Complete interface and data field identities sealed independently from the repository-derived interface inventory.")];args=[];pf={};ff={};pv={required:pvalue("string_set",["field-text","field-duration"])}
    for index,(suffix,typ,values) in enumerate(maps):
        observed=f"{pre}.{suffix}";expected=f"{pre}.approved_{suffix}"
        facts.append(fact(observed,typ,f"Direct explicit versioned {suffix.replace('_',' ')} for every interface field.","repository","Lossless parsed versioned interface contract bytes."));params.append(parameter(expected,typ,f"Exact independently approved {suffix.replace('_',' ')}, including explicit NA only for fields independently classified nonapplicable."));args += [{"op":"map_keys_eq_set_parameter","fact":observed,"parameter":required},{"op":f"{'string_map' if typ=='string_map' else 'identity_map'}_eq_parameter","fact":observed,"parameter":expected}]
        pf[observed]=fvalue(typ,values);pv[expected]=pvalue(typ,values);broken=dict(values);broken["field-duration"]="NA" if suffix=="unit_or_explicit_na_values" else "wrong-value";ff[observed]=fvalue(typ,(broken if suffix=="unit_or_explicit_na_values" else values))
    result=exact_profile(facts,params,{"op":"all","args":args},pf,ff,pv,"Every field's exact versioned bindings match independently reviewed maps; language and unit NA values are sealed per field rather than selected by the parser.","field-duration is independently classified as unit-applicable with unit seconds, but the observed contract says NA.")
    result["fixtures"]["counterexample"]=copy.deepcopy(result["fixtures"]["fail"]);result["fixtures"]["counterexample"]["description"]="A provider could mark every missing language or unit as NA and make presence checks pass; the sealed per-field applicability map rejects NA for field-duration's required unit."
    return result


def mixed_version_policy_profile():
    cid="USEQ-B9F93CB2";ordinal=1;pre=prefix(cid,ordinal);allowed=f"{pre}.effective_allowed_version_combination_configuration_identities";expected_allowed=f"{pre}.approved_allowed_version_combination_ids";cases=f"{pre}.required_unlisted_combination_case_ids";all_combinations=f"{pre}.required_finite_supported_combination_ids";inputs=f"{pre}.unlisted_combination_input_digests";states=f"{pre}.unlisted_combination_serving_states";expected_inputs=f"{pre}.expected_unlisted_combination_input_digests";expected_states=f"{pre}.expected_unlisted_combination_serving_states"
    facts=[fact(allowed,"identity_map","Direct effective configuration record identity keyed by every allowed serving-path and version combination.","environment","Read-only effective serving-path rollout configuration."),fact(inputs,"digest_map","Direct digest of each independently sealed unlisted mixed-version probe input.","environment","Authenticated bounded serving-path probe inputs."),fact(states,"state_map","Direct serving outcome for each unlisted mixed-version probe.","environment","Authenticated bounded serving-path probe outputs.")]
    params=[parameter(expected_allowed,"string_set","Complete allowed serving-path and version-combination identities from independently approved rollout policy."),parameter(cases,"string_set","Complete prohibited serving-path and version-combination identities computed independently."),parameter(all_combinations,"string_set","Complete finite supported serving-path x version cross-product; allowed configuration keys and prohibited probe keys must be an exact disjoint partition."),parameter(expected_inputs,"digest_map","Exact probe input digest for every prohibited combination case."),parameter(expected_states,"state_map","Exact independently approved rejected state for every prohibited combination case.")]
    case_ids=["path-A|v1-v3","path-A|v2-v1"];input_values={case_ids[0]:"a"*64,case_ids[1]:"b"*64};state_values={x:"rejected" for x in case_ids}
    allowed_values={"path-A|v1-v1":"config-entry-A","path-A|v1-v2":"config-entry-B"}
    args=[{"op":"map_keys_eq_set_parameter","fact":allowed,"parameter":expected_allowed},{"op":"map_keys_eq_set_parameter","fact":inputs,"parameter":cases},{"op":"map_keys_eq_set_parameter","fact":states,"parameter":cases},{"op":"map_key_partition_eq_set_parameter","fact":allowed,"other_fact":inputs,"parameter":all_combinations},{"op":"digest_map_eq_parameter","fact":inputs,"parameter":expected_inputs},{"op":"state_map_eq_parameter","fact":states,"parameter":expected_states}]
    pf={allowed:fvalue("identity_map",allowed_values),inputs:fvalue("digest_map",input_values),states:fvalue("state_map",state_values)};ff=copy.deepcopy(pf);ff[allowed]=fvalue("identity_map",{"path-A|v1-v1":"config-entry-A"});pv={expected_allowed:pvalue("string_set",list(allowed_values)),cases:pvalue("string_set",case_ids),all_combinations:pvalue("string_set",list(allowed_values)+case_ids),expected_inputs:pvalue("digest_map",input_values),expected_states:pvalue("state_map",state_values)}
    result=exact_profile(facts,params,{"op":"all","args":args},pf,ff,pv,"Effective configuration enumerates the complete approved set, and every independently computed prohibited combination has an exact rejected probe result.","All unlisted negative probes reject, but effective configuration silently omits approved combination path-A|v1-v2.")
    result["fixtures"]["counterexample"]=copy.deepcopy(result["fixtures"]["fail"]);result["fixtures"]["counterexample"]["description"]="A negative-only suite could pass while an allowed combination is absent from effective config; exact allowed-set equality catches that omission."
    return result


def exact_test_trace_profile():
    cid="USEQ-BD7F03F5";ordinal=1;pre=prefix(cid,ordinal);required=f"{pre}.required_test_result_ids";dimensions=["product_version","configuration_version","feature_flag_state","schema_version","data_state","dependency_set","environment"]
    facts=[];args=[];pf={};ff={};ids=["result-A","result-B"]
    for index,dimension in enumerate(dimensions):
        traced=f"{pre}.traced_{dimension}_identities";actual=f"{pre}.independently_derived_execution_{dimension}_identities"
        facts += [fact(traced,"identity_map",f"Direct {dimension.replace('_',' ')} identity recorded by each test result trace.","executed","Authenticated test result trace record."),fact(actual,"identity_map",f"Direct {dimension.replace('_',' ')} identity independently derived from execution setup and runtime manifests.","executed","Authenticated execution setup, environment, dependency, or runtime input manifest independent of the trace field.")]
        args += [{"op":"map_keys_eq_set_parameter","fact":traced,"parameter":required},{"op":"map_keys_eq_set_parameter","fact":actual,"parameter":required},{"op":"identity_map_eq_fact","fact":traced,"other_fact":actual}]
        values={"result-A":f"{dimension}-A","result-B":f"{dimension}-B"};broken={"result-A":f"wrong-{dimension}","result-B":f"{dimension}-B"};pf[traced]=fvalue("identity_map",values);pf[actual]=fvalue("identity_map",values);ff[traced]=fvalue("identity_map",values);ff[actual]=fvalue("identity_map",(broken if index==0 else values))
    result=exact_profile(facts,[parameter(required,"string_set","Complete authenticated test-result identities sealed independently from the required test inventory.")],{"op":"all","args":args},pf,ff,{required:pvalue("string_set",ids)},"Each trace field equals the corresponding identity independently derived from the actual execution setup for every test result.","The trace says product-version-A, but the execution input manifest shows result-A actually evaluated a different product version.")
    result["fixtures"]["counterexample"]=copy.deepcopy(result["fixtures"]["fail"]);result["fixtures"]["counterexample"]["description"]="Nonempty trace labels can all look complete while naming the intended rather than evaluated product; comparison with independently derived execution identities rejects that trace."
    return result


def fixed_finding_retest_profile():
    cid="PRC-31-021";ordinal=1;pre=prefix(cid,ordinal);required=f"{pre}.required_fixed_finding_ids";runs=f"{pre}.retest_execution_identities_by_finding";states=f"{pre}.retest_terminal_states_by_finding";inputs=f"{pre}.retest_input_digests_by_finding";expected_inputs=f"{pre}.expected_retest_input_digests_by_finding";completed=f"{pre}.approved_completed_retest_state_ids"
    facts=[fact(runs,"identity_map","Direct authenticated retest execution identity for every fixed finding.","executed","Authenticated test invocation records."),fact(states,"state_map","Direct terminal state of the retest for every fixed finding; pass and fail are both completed retests.","executed","Authenticated test completion records."),fact(inputs,"digest_map","Direct digest of the finding-specific retest input and assessed fixed revision.","executed","Authenticated retest input manifests.")]
    params=[parameter(required,"string_set","Complete finding identities independently sealed from authenticated finding state transitions to fixed."),parameter(expected_inputs,"digest_map","Exact approved retest input and fixed-revision digest for every fixed finding."),parameter(completed,"string_set","Direct terminal states independently approved to mean the retest actually completed, regardless of pass or fail.")]
    ids=["finding-A","finding-B"];run_values={"finding-A":"retest-A","finding-B":"retest-B"};state_values={"finding-A":"completed-pass","finding-B":"completed-fail"};input_values={"finding-A":"a"*64,"finding-B":"b"*64}
    args=[{"op":"map_keys_eq_set_parameter","fact":x,"parameter":required} for x in [runs,states,inputs]]+[{"op":"digest_map_eq_parameter","fact":inputs,"parameter":expected_inputs},{"op":"state_map_values_in_parameter","fact":states,"parameter":completed}]
    pf={runs:fvalue("identity_map",run_values),states:fvalue("state_map",state_values),inputs:fvalue("digest_map",input_values)};ff=copy.deepcopy(pf);ff[states]=fvalue("state_map",{"finding-A":"not-run","finding-B":"completed-fail"});pv={required:pvalue("string_set",ids),expected_inputs:pvalue("digest_map",input_values),completed:pvalue("string_set",["completed-pass","completed-fail"])}
    return exact_profile(facts,params,{"op":"all","args":args},pf,ff,pv,"Every independently inventoried fixed finding has a completed retest bound to the exact fixed revision and approved retest input.","finding-A is marked fixed but its retest state is not-run, even though finding-B has a completed retest.")


def state_predicate_wait_profile():
    cid="USEQ-7E74EFCD";ordinal=1;pre=prefix(cid,ordinal);required=f"{pre}.required_async_wait_site_ids";mechanisms=f"{pre}.observed_wait_mechanism_states";predicates=f"{pre}.declared_state_predicate_digests";expected_predicates=f"{pre}.approved_state_predicate_digests";allowed=f"{pre}.approved_state_predicate_wait_mode_ids";analyzer=f"{pre}.trace_analyzer_identity";config=f"{pre}.trace_analyzer_configuration_digest";expected_analyzer=f"{pre}.required_trace_analyzer_identity";expected_config=f"{pre}.required_trace_analyzer_configuration_digest"
    facts=[fact(mechanisms,"state_map","Direct observed wait mechanism for each asynchronous, queue, or eventual-consistency test wait site.","executed","Authenticated execution trace wait events."),fact(predicates,"digest_map","Direct digest of the declared state predicate evaluated at each wait site.","executed","Authenticated execution trace predicate bindings."),fact(analyzer,"identity","Direct pinned trace-analyzer identity and version.","executed","Authenticated trace analyzer invocation."),fact(config,"digest","Direct digest of trace rules that distinguish state predicates from fixed delays.","executed","Exact trace analyzer rules.")]
    params=[parameter(required,"string_set","Complete asynchronous, queue, and eventual-consistency wait-site identities sealed from the versioned test case inventory."),parameter(expected_predicates,"digest_map","Exact independently approved declared state-predicate digest for every wait site."),parameter(allowed,"string_set","Direct wait-mode identifiers approved to mean predicate-driven waiting rather than fixed delay."),parameter(expected_analyzer,"identity","Approved trace analyzer identity and version."),parameter(expected_config,"digest","Approved trace analyzer rule digest.")]
    ids=["async-wait-A","queue-wait-A","eventual-wait-A"];mechanism_values={x:"state-predicate" for x in ids};predicate_values={key:format(index+1,"x")*64 for index,key in enumerate(ids)}
    args=[{"op":"map_keys_eq_set_parameter","fact":mechanisms,"parameter":required},{"op":"map_keys_eq_set_parameter","fact":predicates,"parameter":required},{"op":"state_map_values_in_parameter","fact":mechanisms,"parameter":allowed},{"op":"digest_map_eq_parameter","fact":predicates,"parameter":expected_predicates},{"op":"identity_eq_parameter","fact":analyzer,"parameter":expected_analyzer},{"op":"digest_eq_parameter","fact":config,"parameter":expected_config}]
    pf={mechanisms:fvalue("state_map",mechanism_values),predicates:fvalue("digest_map",predicate_values),analyzer:fvalue("identity","trace-analyzer@1"),config:fvalue("digest","d"*64)};ff=copy.deepcopy(pf);ff[mechanisms]=fvalue("state_map",{"async-wait-A":"fixed-delay","queue-wait-A":"state-predicate","eventual-wait-A":"state-predicate"});pv={required:pvalue("string_set",ids),expected_predicates:pvalue("digest_map",predicate_values),allowed:pvalue("string_set",["state-predicate"]),expected_analyzer:pvalue("identity","trace-analyzer@1"),expected_config:pvalue("digest","d"*64)}
    return exact_profile(facts,params,{"op":"all","args":args},pf,ff,pv,"The pinned trace analyzer observes a declared state predicate at every sealed async wait site and binds each predicate to its approved digest.","The queue and eventual-consistency waits use predicates, but async-wait-A uses a fixed delay.")


def versioned_repository_items_profile():
    cid="USEQ-69B5C63D";ordinal=1;pre=prefix(cid,ordinal);items=f"{pre}.required_schema_transformation_migration_item_ids";references=f"{pre}.required_reference_ids";versions=f"{pre}.immutable_version_identities_by_item";digests=f"{pre}.content_digests_by_item";targets=f"{pre}.resolved_target_version_identities_by_reference";immutable=f"{pre}.target_immutability_states_by_reference";expected_versions=f"{pre}.approved_item_version_identities";expected_digests=f"{pre}.approved_item_content_digests";expected_targets=f"{pre}.approved_reference_target_version_identities";approved_version_ids=f"{pre}.approved_inventoried_version_identities";allowed_immutable=f"{pre}.approved_immutable_target_state_ids"
    facts=[fact(versions,"identity_map","Direct immutable version identity for every repository schema, transformation, and migration-logic item.","repository","Lossless repository inventory and version metadata."),fact(digests,"digest_map","Locally recomputed content digest for every inventoried item.","repository","Exact item bytes."),fact(targets,"identity_map","Direct resolved target version identity for every repository reference.","repository","Lossless parsed references resolved by the pinned repository parser."),fact(immutable,"state_map","Direct target storage/reference state for every resolved reference.","repository","Exact reference and target metadata proving immutable rather than mutable resolution.")]
    params=[parameter(items,"string_set","Complete schema, transformation, and migration-logic item identities sealed from the assessed repository inventory."),parameter(references,"string_set","Complete reference identities sealed from the assessed repository parser inventory."),parameter(expected_versions,"identity_map","Exact independently approved immutable version identity for every item."),parameter(expected_digests,"digest_map","Exact independently recomputed content digest for every item."),parameter(expected_targets,"identity_map","Exact approved single target version for every reference."),parameter(approved_version_ids,"string_set","All and only immutable version identities from the sealed item inventory."),parameter(allowed_immutable,"string_set","Raw target states independently approved as immutable." )]
    item_ids=["schema-A","transform-A","migration-A"];ref_ids=["ref-schema","ref-transform","ref-migration"];version_values={"schema-A":"schema-v1","transform-A":"transform-v2","migration-A":"migration-v3"};digest_values={key:format(index+1,"x")*64 for index,key in enumerate(item_ids)};target_values={"ref-schema":"schema-v1","ref-transform":"transform-v2","ref-migration":"migration-v3"};immutable_values={x:"immutable" for x in ref_ids}
    args=[{"op":"map_keys_eq_set_parameter","fact":versions,"parameter":items},{"op":"map_keys_eq_set_parameter","fact":digests,"parameter":items},{"op":"identity_map_eq_parameter","fact":versions,"parameter":expected_versions},{"op":"digest_map_eq_parameter","fact":digests,"parameter":expected_digests},{"op":"map_keys_eq_set_parameter","fact":targets,"parameter":references},{"op":"map_keys_eq_set_parameter","fact":immutable,"parameter":references},{"op":"identity_map_eq_parameter","fact":targets,"parameter":expected_targets},{"op":"identity_map_values_in_parameter","fact":targets,"parameter":approved_version_ids},{"op":"identity_map_values_in_fact","fact":targets,"other_fact":versions},{"op":"state_map_values_in_parameter","fact":immutable,"parameter":allowed_immutable}]
    pf={versions:fvalue("identity_map",version_values),digests:fvalue("digest_map",digest_values),targets:fvalue("identity_map",target_values),immutable:fvalue("state_map",immutable_values)};ff=copy.deepcopy(pf);ff[targets]=fvalue("identity_map",{"ref-schema":"mutable-latest","ref-transform":"transform-v2","ref-migration":"migration-v3"});pv={items:pvalue("string_set",item_ids),references:pvalue("string_set",ref_ids),expected_versions:pvalue("identity_map",version_values),expected_digests:pvalue("digest_map",digest_values),expected_targets:pvalue("identity_map",target_values),approved_version_ids:pvalue("string_set",list(version_values.values())),allowed_immutable:pvalue("string_set",["immutable"])}
    return exact_profile(facts,params,{"op":"all","args":args},pf,ff,pv,"Every inventoried item has its sealed immutable version and digest, and every reference resolves to one of those exact immutable versions.","All items are versioned and hashed, but ref-schema resolves to mutable-latest instead of the inventoried schema-v1 identity.")


def cryptographic_requirements_profile():
    cid="USEQ-9E10A498";ordinal=1;pre=prefix(cid,ordinal);required=f"{pre}.required_requirement_category_ids";pairs=f"{pre}.required_distinct_category_pairs";categories=["encryption","message-authentication","digital-signature","key-agreement","password-hashing","key-derivation","random-generation"]
    field_specs=[("requirement_record_identities","identity_map"),("effective_version_identities","identity_map"),("requirement_record_content_digests","digest_map"),("approved_purpose_identities","identity_map"),("approved_algorithm_set_identities","identity_map"),("approved_parameter_set_identities","identity_map"),("approved_key_or_input_type_set_identities","identity_map"),("approved_output_contract_identities","identity_map"),("approved_validation_rule_set_identities","identity_map")]
    facts=[];params=[parameter(required,"string_set","Exactly the seven cryptographic requirement category identities named by the clause."),parameter(pairs,"directed_graph","Every unordered pair of distinct cryptographic categories, sealed independently; record identities must differ for each pair.")];args=[];pf={};ff={};pv={required:pvalue("string_set",categories)}
    for index,(suffix,typ) in enumerate(field_specs):
        observed=f"{pre}.{suffix}";expected=f"{pre}.expected_{suffix}"
        facts.append(fact(observed,typ,f"Direct current versioned {suffix.replace('_',' ')} keyed by cryptographic category.","repository","Lossless parsed current cryptographic requirements artifact and locally hashed record bytes."));params.append(parameter(expected,typ,f"Exact independently reviewed current {suffix.replace('_',' ')} for each of the seven categories."));args += [{"op":"map_keys_eq_set_parameter","fact":observed,"parameter":required},{"op":f"{'digest_map' if typ=='digest_map' else 'identity_map'}_eq_parameter","fact":observed,"parameter":expected}]
        if typ=="digest_map": values={category:format(index+1,"x")*64 for category in categories};wrong=dict(values);wrong["encryption"]="f"*64
        else: values={category:f"{suffix}:{category}" for category in categories};wrong=dict(values);wrong["message-authentication"]=values["encryption"]
        pf[observed]=fvalue(typ,values);pv[expected]=pvalue(typ,values);ff[observed]=fvalue(typ,(wrong if suffix=="requirement_record_identities" else values))
    edges=[]
    for i,left in enumerate(categories):
        for right in categories[i+1:]: edges.append({"from":left,"to":right})
    pv[pairs]=pvalue("directed_graph",edges);args.append({"op":"identity_map_values_differ_for_pairs_parameter","fact":f"{pre}.requirement_record_identities","parameter":pairs})
    result=exact_profile(facts,params,{"op":"all","args":args},pf,ff,pv,"Seven exact current record identities are pairwise distinct, and every record's version, digest, purpose, algorithms, parameters, inputs, outputs, and validation rules match independently reviewed content.","All seven category keys have nonempty fields, but encryption and message authentication reuse one requirement record, so the records are not separate or non-overlapping.")
    result["fixtures"]["counterexample"]=copy.deepcopy(result["fixtures"]["fail"]);result["fixtures"]["counterexample"]["description"]="A parser can repeat one complete encryption record under all seven category names and satisfy field-presence checks; the pairwise record-identity relation rejects that reuse."
    return result


def cors_invariant_profile():
    cid="USEQ-BBED00FD";ordinal=1;pre=prefix(cid,ordinal);required=f"{pre}.required_endpoint_ids";origins=f"{pre}.effective_trusted_origin_set_digests";expected_origins=f"{pre}.approved_trusted_origin_set_digests";credentials=f"{pre}.credentialed_access_enabled";wildcard_disabled=f"{pre}.wildcard_origin_disabled";policy=f"{pre}.effective_origin_policy_digests";expected_policy=f"{pre}.reviewed_origin_policy_digests"
    facts=[fact(origins,"digest_map","Canonical digest of the exact trusted-origin set configured for every endpoint.","environment","Read-only effective cross-origin configuration."),fact(credentials,"boolean_map","Direct effective credentialed-cross-origin-access Boolean for every endpoint.","environment","Read-only effective credentials flag."),fact(wildcard_disabled,"boolean_map","Direct Boolean obtained from the effective allowed-origin field: true exactly when no wildcard origin is present.","environment","Lossless parsed effective allowed-origin list."),fact(policy,"digest_map","Direct digest of the effective origin-policy bytes applied to every endpoint.","environment","Read-only effective origin policy.")]
    params=[parameter(required,"string_set","Complete in-scope endpoint identities sealed independently."),parameter(expected_origins,"digest_map","Exact independently reviewed trusted-origin set digest for every endpoint."),parameter(expected_policy,"digest_map","Exact independently reviewed origin-policy digest for every endpoint.")]
    ids=["endpoint-A","endpoint-B"];origin_values={"endpoint-A":"a"*64,"endpoint-B":"b"*64};credential_values={"endpoint-A":True,"endpoint-B":False};disabled_values={"endpoint-A":True,"endpoint-B":False};policy_values={x:"c"*64 for x in ids}
    args=[{"op":"map_keys_eq_set_parameter","fact":x,"parameter":required} for x in [origins,credentials,wildcard_disabled,policy]]+[{"op":"digest_map_eq_parameter","fact":origins,"parameter":expected_origins},{"op":"digest_map_eq_parameter","fact":policy,"parameter":expected_policy},{"op":"boolean_map_implies_fact","fact":credentials,"other_fact":wildcard_disabled}]
    pf={origins:fvalue("digest_map",origin_values),credentials:fvalue("boolean_map",credential_values),wildcard_disabled:fvalue("boolean_map",disabled_values),policy:fvalue("digest_map",policy_values)};ff=copy.deepcopy(pf);ff[wildcard_disabled]=fvalue("boolean_map",{"endpoint-A":False,"endpoint-B":False});pv={required:pvalue("string_set",ids),expected_origins:pvalue("digest_map",origin_values),expected_policy:pvalue("digest_map",policy_values)}
    result=exact_profile(facts,params,{"op":"all","args":args},pf,ff,pv,"Reviewed trusted origins and policy digests apply to every endpoint, and the fixed implication requires wildcard absence whenever credentials are enabled.","endpoint-A has credentialed access enabled and a wildcard origin; even a policy parameter that expected both cannot make the mechanical implication pass.")
    result["fixtures"]["counterexample"]=copy.deepcopy(result["fixtures"]["fail"]);result["fixtures"]["counterexample"]["description"]="A policy-selected expected-state map could redefine credentialed-plus-wildcard as acceptable; the literal implication rejects that unsafe pair independently of policy."
    return result


def cost_protection_enabled_profile():
    cid="USEQ-A0E78A59";ordinal=1;pre=prefix(cid,ordinal);required=f"{pre}.required_billing_scope_resource_class_ids";alerts=f"{pre}.bound_cost_alert_identities";protections=f"{pre}.bound_runaway_protection_identities";alert_enabled=f"{pre}.cost_alert_enabled";protection_enabled=f"{pre}.runaway_protection_enabled"
    facts=[fact(alerts,"identity_map","Direct cost-alert configuration identity bound to each billing scope and resource class.","environment","Read-only effective billing configuration."),fact(protections,"identity_map","Direct runaway-resource protection identity bound to each exact scope and class.","environment","Read-only effective resource-control configuration."),fact(alert_enabled,"boolean_map","Direct effective enabled Boolean from every bound cost-alert record.","environment","Read-only raw enabled field."),fact(protection_enabled,"boolean_map","Direct effective enabled Boolean from every bound runaway-protection record.","environment","Read-only raw enabled field.")]
    params=[parameter(required,"string_set","Complete billing scope x resource-class identities sealed independently.")]
    ids=["billing-A|compute","billing-A|storage"];alert_values={"billing-A|compute":"alert-A","billing-A|storage":"alert-B"};protection_values={"billing-A|compute":"guard-A","billing-A|storage":"guard-B"};enabled_values={x:True for x in ids}
    args=[{"op":"map_keys_eq_set_parameter","fact":x,"parameter":required} for x in [alerts,protections,alert_enabled,protection_enabled]]+[{"op":"boolean_map_all_eq","fact":alert_enabled,"boolean":True},{"op":"boolean_map_all_eq","fact":protection_enabled,"boolean":True}]
    pf={alerts:fvalue("identity_map",alert_values),protections:fvalue("identity_map",protection_values),alert_enabled:fvalue("boolean_map",enabled_values),protection_enabled:fvalue("boolean_map",enabled_values)};ff=copy.deepcopy(pf);ff[protection_enabled]=fvalue("boolean_map",{"billing-A|compute":False,"billing-A|storage":True})
    result=exact_profile(facts,params,{"op":"all","args":args},pf,ff,{required:pvalue("string_set",ids)},"Every sealed billing-scope and resource-class key has both mechanisms bound, and both raw enabled fields are literally true.","compute has a bound runaway-protection record but its enabled field is false; record presence does not satisfy enabled protection.")
    return result


def triple_digest_map_profile(control_id, ordinal, authority, subject_suffix, left_suffix, middle_suffix, right_suffix, meanings, subject_examples):
    pre=prefix(control_id,ordinal);required=f"{pre}.required_{subject_suffix}";keys=[f"{pre}.{x}" for x in [left_suffix,middle_suffix,right_suffix]]
    facts=[fact(key,"digest_map",meaning,authority,"Exact authenticated record or locally recomputed bytes named by this side of the digest relation.") for key,meaning in zip(keys,meanings)]
    args=[{"op":"map_keys_eq_set_parameter","fact":key,"parameter":required} for key in keys]+[{"op":"digest_map_eq_fact","fact":keys[0],"other_fact":keys[1]},{"op":"digest_map_eq_fact","fact":keys[0],"other_fact":keys[2]}]
    values={subject:format(index+1,"x")*64 for index,subject in enumerate(subject_examples)};pf={key:fvalue("digest_map",values) for key in keys};ff=copy.deepcopy(pf);broken=dict(values);broken[subject_examples[0]]="f"*64;ff[keys[2]]=fvalue("digest_map",broken)
    return exact_profile(facts,[parameter(required,"string_set",f"Complete {subject_suffix.replace('_',' ')} sealed independently.")],{"op":"all","args":args},pf,ff,{required:pvalue("string_set",subject_examples)},"All three direct digest maps cover the same sealed identities and match per identity.","The first two digest sources agree, but the third binds different bytes for one sealed identity.")


def emergency_access_expiry_profile():
    cid="PRC-17-020";ordinal=2;pre=prefix(cid,ordinal);required=f"{pre}.required_emergency_access_grant_ids";issue=f"{pre}.issue_times";expiry=f"{pre}.expiry_times";post=f"{pre}.post_expiry_access_enabled";maximum=f"{pre}.approved_maximum_grant_duration_seconds"
    facts=[fact(issue,"time_map","Direct authenticated issue timestamp for every emergency-access grant.","environment","Authenticated grant records."),fact(expiry,"time_map","Direct authenticated expiry timestamp for every grant.","environment","Authenticated grant records."),fact(post,"boolean_map","Direct effective access-enabled Boolean observed after each grant expiry.","environment","Read-only authorization decision or permission state after the expiry instant.")]
    params=[parameter(required,"string_set","Complete in-scope emergency-access grant identities sealed independently."),parameter(maximum,"number","Approved maximum nonnegative grant duration in seconds." )]
    ids=["grant-A","grant-B"];issue_values={"grant-A":"2026-08-29T10:00:00Z","grant-B":"2026-08-29T11:00:00Z"};expiry_values={"grant-A":"2026-08-29T10:30:00Z","grant-B":"2026-08-29T11:30:00Z"};disabled={x:False for x in ids}
    args=[{"op":"map_keys_eq_set_parameter","fact":x,"parameter":required} for x in [issue,expiry,post]]+[{"op":"time_map_delta_lte_parameter","fact":issue,"other_fact":expiry,"parameter":maximum},{"op":"boolean_map_all_eq","fact":post,"boolean":False}]
    pf={issue:fvalue("time_map",issue_values),expiry:fvalue("time_map",expiry_values),post:fvalue("boolean_map",disabled)};ff=copy.deepcopy(pf);ff[post]=fvalue("boolean_map",{"grant-A":True,"grant-B":False})
    return exact_profile(facts,params,{"op":"all","args":args},pf,ff,{required:pvalue("string_set",ids),maximum:pvalue("number",3600)},"Every grant has issue and expiry times within the fixed maximum and literal post-expiry access-disabled evidence.","grant-A expires on time but remains usable after expiry.")


def redundancy_failure_domain_profile():
    cid="PRC-27-021";ordinal=1;pre=prefix(cid,ordinal);domains=f"{pre}.replica_failure_domain_identities";required=f"{pre}.required_replica_ids";pairs=f"{pre}.required_distinct_replica_pairs"
    facts=[fact(domains,"identity_map","Direct effective failure-domain identity for every replica in the redundancy group.","environment","Read-only effective placement topology.")]
    params=[parameter(required,"string_set","Complete replica identities sealed from the intended redundancy group."),parameter(pairs,"directed_graph","Every replica pair that the independently approved survival model requires to occupy different failure domains." )]
    ids=["replica-A","replica-B","replica-C"];values={"replica-A":"zone-A","replica-B":"zone-B","replica-C":"zone-C"};edges=[{"from":"replica-A","to":"replica-B"},{"from":"replica-A","to":"replica-C"},{"from":"replica-B","to":"replica-C"}]
    pred={"op":"all","args":[{"op":"map_keys_eq_set_parameter","fact":domains,"parameter":required},{"op":"identity_map_values_differ_for_pairs_parameter","fact":domains,"parameter":pairs}]};ff=dict(values);ff["replica-B"]="zone-A"
    return exact_profile(facts,params,pred,{domains:fvalue("identity_map",values)},{domains:fvalue("identity_map",ff)},{required:pvalue("string_set",ids),pairs:pvalue("directed_graph",edges)},"Every replica pair required by the survival model has different direct failure-domain identities.","replica-A and replica-B occupy the same zone, so losing that zone defeats the intended redundancy.")


def external_match_freshness_profile(control_id, ordinal, subject_suffix, observed_suffix, tracked_suffix, observed_meaning, tracked_meaning, subject_examples):
    pre=prefix(control_id,ordinal);required=f"{pre}.required_{subject_suffix}";observed=f"{pre}.{observed_suffix}";tracked=f"{pre}.{tracked_suffix}";verified=f"{pre}.verification_times";collected=f"{pre}.evidence_collection_times";maximum=f"{pre}.approved_maximum_evidence_age_seconds"
    facts=[fact(observed,"identity_map",observed_meaning,"external_registry","Authenticated current external record."),fact(tracked,"identity_map",tracked_meaning,"external_registry","Authenticated recorded project or route value joined to the same sealed identity."),fact(verified,"time_map","Direct timestamp when each external value or route was verified.","external_registry","Authenticated provider or directory verification event."),fact(collected,"time_map","Direct evidence collection timestamp for each sealed subject.","external_registry","Authenticated scanner collection record.")]
    params=[parameter(required,"string_set",f"Complete {subject_suffix.replace('_',' ')} sealed independently."),parameter(maximum,"number","Approved maximum external evidence age in seconds.")]
    values={subject:f"identity-{index}" for index,subject in enumerate(subject_examples)};verified_values={subject:"2026-08-29T10:00:00Z" for subject in subject_examples};collected_values={subject:"2026-08-29T10:30:00Z" for subject in subject_examples}
    args=[{"op":"map_keys_eq_set_parameter","fact":x,"parameter":required} for x in [observed,tracked,verified,collected]]+[{"op":"identity_map_eq_fact","fact":observed,"other_fact":tracked},{"op":"time_map_delta_lte_parameter","fact":verified,"other_fact":collected,"parameter":maximum}]
    pf={observed:fvalue("identity_map",values),tracked:fvalue("identity_map",values),verified:fvalue("time_map",verified_values),collected:fvalue("time_map",collected_values)};ff=copy.deepcopy(pf);broken=dict(values);broken[subject_examples[0]]="different-external-identity";ff[observed]=fvalue("identity_map",broken)
    return exact_profile(facts,params,{"op":"all","args":args},pf,ff,{required:pvalue("string_set",subject_examples),maximum:pvalue("number",3600)},"Every tracked value equals the current authenticated external value and was verified within the fixed maximum age.","The external record changed for one subject even though its verification timestamp is recent.")


def tamper_evident_record_profile(control_id, ordinal, subject_suffix, subject_examples):
    pre=prefix(control_id,ordinal);required=f"{pre}.required_{subject_suffix}";stored=f"{pre}.stored_record_digests";recomputed=f"{pre}.locally_recomputed_record_byte_digests";mutation=f"{pre}.unauthorized_mutation_or_deletion_enabled"
    facts=[fact(stored,"digest_map","Direct stored integrity digest for every record.","environment","Authenticated storage metadata."),fact(recomputed,"digest_map","Locally recomputed digest of exact retrieved record bytes.","environment","Exact retrieved record bytes."),fact(mutation,"boolean_map","Direct effective Boolean indicating whether unauthorized mutation or deletion is enabled for each record store binding.","environment","Read-only effective storage permission and retention configuration.")]
    params=[parameter(required,"string_set",f"Complete {subject_suffix.replace('_',' ')} sealed independently.")]
    values={subject:format(index+1,"x")*64 for index,subject in enumerate(subject_examples)};disabled={subject:False for subject in subject_examples}
    args=[{"op":"map_keys_eq_set_parameter","fact":x,"parameter":required} for x in [stored,recomputed,mutation]]+[{"op":"digest_map_eq_fact","fact":stored,"other_fact":recomputed},{"op":"boolean_map_all_eq","fact":mutation,"boolean":False}]
    pf={stored:fvalue("digest_map",values),recomputed:fvalue("digest_map",values),mutation:fvalue("boolean_map",disabled)};ff=copy.deepcopy(pf);ff[mutation]=fvalue("boolean_map",{subject:(subject==subject_examples[0]) for subject in subject_examples})
    return exact_profile(facts,params,{"op":"all","args":args},pf,ff,{required:pvalue("string_set",subject_examples)},"Every retrieved record matches its stored digest and direct effective permissions forbid unauthorized mutation and deletion.","The record bytes still match today, but the effective store grants unauthorized mutation for one record, so future integrity is unprotected.")


def number_map_threshold_profile(control_id, ordinal, authority, suffix, semantics, direction="lte"):
    pre=prefix(control_id,ordinal);observed=f"{pre}.{suffix}";expected=f"{pre}.approved_thresholds"
    facts=[fact(observed,"number_map",semantics,authority,"Complete direct raw numeric values in canonical units, keyed by every sealed subject and objective.")]
    params=[parameter(expected,"number_map","Complete independently approved numeric thresholds in matching canonical units and keys.")]
    op=f"number_map_{direction}_parameter";limits={"case-A":10,"case-B":20}
    if direction in ("lte","lt"):good={"case-A":9,"case-B":19};bad={"case-A":11,"case-B":19}
    else:good={"case-A":11,"case-B":21};bad={"case-A":9,"case-B":21}
    return exact_profile(facts,params,{"op":op,"fact":observed,"parameter":expected},{observed:fvalue("number_map",good)},{observed:fvalue("number_map",bad)},{expected:pvalue("number_map",limits)},"Every direct raw numeric value satisfies its independently approved threshold.","One direct raw numeric value violates its independently approved threshold.")


def case_state_profile(control_id, ordinal, authority, suffix, semantics):
    pre=prefix(control_id,ordinal);states=f"{pre}.{suffix}";revisions=f"{pre}.execution_revision_identities";inputs=f"{pre}.execution_input_digests";required=f"{pre}.required_case_ids";expected_states=f"{pre}.expected_case_terminal_states";expected_revisions=f"{pre}.expected_execution_revision_identities";expected_inputs=f"{pre}.expected_execution_input_digests"
    facts=[fact(states,"state_map",semantics,authority,"Direct raw terminal state or output identifier keyed by every independently sealed execution case."),fact(revisions,"identity_map","Direct assessed executable or release revision identity used by every sealed execution case.",authority,"Authenticated execution invocation record for every sealed case."),fact(inputs,"digest_map","Direct digest of the exact fixture, configuration, and bound inputs used by every sealed execution case.",authority,"Exact authenticated execution input manifest bytes for every sealed case.")]
    params=[parameter(required,"string_set","Complete stable case identities from the independently sealed versioned case matrix."),parameter(expected_states,"state_map","Exact expected direct terminal state keyed by every sealed case."),parameter(expected_revisions,"identity_map","Exact assessed executable or release revision identity keyed by every sealed case."),parameter(expected_inputs,"digest_map","Exact approved fixture, configuration, and input digest keyed by every sealed case.")]
    args=[]
    for fact_id in [states,revisions,inputs]:args.append({"op":"map_keys_eq_set_parameter","fact":fact_id,"parameter":required})
    args += [{"op":"state_map_eq_parameter","fact":states,"parameter":expected_states},{"op":"identity_map_eq_parameter","fact":revisions,"parameter":expected_revisions},{"op":"digest_map_eq_parameter","fact":inputs,"parameter":expected_inputs}]
    case_ids=["case-A","case-B"];good_states={"case-A":"allowed","case-B":"rejected"};good_revisions={"case-A":"release-A","case-B":"release-A"};good_inputs={"case-A":"a"*64,"case-B":"b"*64}
    pf={states:fvalue("state_map",good_states),revisions:fvalue("identity_map",good_revisions),inputs:fvalue("digest_map",good_inputs)};ff=copy.deepcopy(pf);ff[states]=fvalue("state_map",{"case-A":"unexpected","case-B":"rejected"});pv={required:pvalue("string_set",case_ids),expected_states:pvalue("state_map",good_states),expected_revisions:pvalue("identity_map",good_revisions),expected_inputs:pvalue("digest_map",good_inputs)}
    return exact_profile(facts,params,{"op":"all","args":args},pf,ff,pv,"Every sealed case has the exact expected direct terminal state and is bound to the exact assessed revision and inputs.","One case has an unexpected direct terminal state even though revision and input bindings remain valid.")


def time_order_profile(control_id, ordinal, authority, earlier_suffix, later_suffix, earlier_semantics, later_semantics, strict=True):
    pre=prefix(control_id,ordinal); earlier=f"{pre}.{earlier_suffix}"; later=f"{pre}.{later_suffix}"; required=f"{pre}.required_event_ids"
    facts=[fact(earlier,"time_map",earlier_semantics,authority,"Authenticated raw earlier-event timestamps keyed by every sealed event identity."),fact(later,"time_map",later_semantics,authority,"Authenticated raw later-event timestamps keyed by the same sealed event identities.")]
    params=[parameter(required,"string_set","Complete stable event identities sealed independently before event collection.")]
    args=[{"op":"map_keys_eq_set_parameter","fact":earlier,"parameter":required},{"op":"map_keys_eq_set_parameter","fact":later,"parameter":required},{"op":"time_map_before_fact" if strict else "time_map_lte_fact","fact":earlier,"other_fact":later}]
    pf={earlier:fvalue("time_map",{"event-A":"2026-08-29T10:00:00Z","event-B":"2026-08-29T11:00:00Z"}),later:fvalue("time_map",{"event-A":"2026-08-29T10:01:00Z","event-B":"2026-08-29T11:01:00Z"})}
    ff=copy.deepcopy(pf);ff[later]=fvalue("time_map",{"event-A":"2026-08-29T09:59:00Z","event-B":"2026-08-29T11:01:00Z"})
    return exact_profile(facts,params,{"op":"all","args":args},pf,ff,{required:pvalue("string_set",["event-A","event-B"])},"Every sealed event has raw timestamps in the required order.","One sealed event records the later event before the required earlier event.")


def time_delta_profile(control_id, ordinal, authority, start_suffix, end_suffix, maximum_suffix, start_semantics, end_semantics):
    pre=prefix(control_id,ordinal);start=f"{pre}.{start_suffix}";end=f"{pre}.{end_suffix}";required=f"{pre}.required_event_ids";maximum=f"{pre}.{maximum_suffix}"
    facts=[fact(start,"time_map",start_semantics,authority,"Complete authenticated raw start timestamps keyed by every sealed event."),fact(end,"time_map",end_semantics,authority,"Complete authenticated raw end timestamps keyed by the same sealed events.")]
    params=[parameter(required,"string_set","Complete stable event identities sealed independently before collection."),parameter(maximum,"number","Approved maximum nonnegative elapsed seconds sealed independently.")]
    args=[{"op":"map_keys_eq_set_parameter","fact":start,"parameter":required},{"op":"map_keys_eq_set_parameter","fact":end,"parameter":required},{"op":"time_map_delta_lte_parameter","fact":start,"other_fact":end,"parameter":maximum}]
    pf={start:fvalue("time_map",{"event-A":"2026-08-29T10:00:00Z","event-B":"2026-08-29T11:00:00Z"}),end:fvalue("time_map",{"event-A":"2026-08-29T10:30:00Z","event-B":"2026-08-29T11:30:00Z"})};ff=copy.deepcopy(pf);ff[end]=fvalue("time_map",{"event-A":"2026-08-29T12:00:00Z","event-B":"2026-08-29T11:30:00Z"})
    return exact_profile(facts,params,{"op":"all","args":args},pf,ff,{required:pvalue("string_set",["event-A","event-B"]),maximum:pvalue("number",3600)},"Every sealed event has a nonnegative raw elapsed time within the approved maximum.","One sealed event exceeds the approved maximum elapsed time.")


def distinct_identities_profile(control_id, ordinal, authority, left_suffix, right_suffix, left_semantics, right_semantics):
    pre=prefix(control_id,ordinal);left=f"{pre}.{left_suffix}";right=f"{pre}.{right_suffix}"
    facts=[fact(left,"identity",left_semantics,authority,"Direct first raw identity from the authoritative source."),fact(right,"identity",right_semantics,authority,"Direct second raw identity from the authoritative source.")]
    pred={"op":"not","arg":{"op":"identity_eq_fact","fact":left,"other_fact":right}};pf={left:fvalue("identity","identity-A"),right:fvalue("identity","identity-B")};ff={left:fvalue("identity","identity-A"),right:fvalue("identity","identity-A")}
    return exact_profile(facts,[],pred,pf,ff,{},"The two direct raw identities are distinct.","The two roles or channels resolve to the same identity.")


def pinned_tool_profile(control_id, ordinal, authority, tool_suffix, config_suffix, tool_semantics, config_semantics):
    pre=prefix(control_id,ordinal);tool=f"{pre}.{tool_suffix}";config=f"{pre}.{config_suffix}";expected_tool=f"{pre}.required_{tool_suffix}";expected_config=f"{pre}.required_{config_suffix}"
    facts=[fact(tool,"identity",tool_semantics,authority,"Authenticated exact tool invocation identity and version."),fact(config,"digest",config_semantics,authority,"Exact tool configuration and policy bytes used by the invocation.")]
    params=[parameter(expected_tool,"identity","Independently approved exact tool identity and version."),parameter(expected_config,"digest","Independently approved exact configuration and policy digest.")]
    pred={"op":"all","args":[{"op":"identity_eq_parameter","fact":tool,"parameter":expected_tool},{"op":"digest_eq_parameter","fact":config,"parameter":expected_config}]};pf={tool:fvalue("identity","tool-A@1"),config:fvalue("digest","a"*64)};ff={tool:fvalue("identity","tool-B@1"),config:fvalue("digest","a"*64)};pv={expected_tool:pvalue("identity","tool-A@1"),expected_config:pvalue("digest","a"*64)}
    return exact_profile(facts,params,pred,pf,ff,pv,"The exact invoked tool and configuration match independently approved identities.","The invoked tool identity differs from the independently approved tool.")


def review_bypass_audit_profile():
    cid="PRC-07-008";ordinal=1;pre=prefix(cid,ordinal);required=f"{pre}.required_review_bypass_execution_ids";actors=f"{pre}.bypass_event_actor_identities";subjects=f"{pre}.bypass_event_subject_identities";times=f"{pre}.bypass_event_timestamps";stored=f"{pre}.stored_bypass_event_digests";recomputed=f"{pre}.locally_recomputed_bypass_event_byte_digests";modes=f"{pre}.audit_storage_enforcement_modes";retention=f"{pre}.audit_retention_policy_digests";expected_actors=f"{pre}.expected_bypass_actor_identities";expected_subjects=f"{pre}.expected_bypass_subject_identities";approved_modes=f"{pre}.approved_append_only_storage_modes";expected_retention=f"{pre}.approved_audit_retention_policy_digests"
    facts=[fact(actors,"identity_map","Direct administrator actor identity on every raw review-bypass audit event.","environment","Authenticated raw event actor fields."),fact(subjects,"identity_map","Direct affected repository and branch identity on every raw bypass event.","environment","Authenticated raw event subject fields."),fact(times,"time_map","Direct timestamp on every raw bypass event.","environment","Authenticated raw event time fields."),fact(stored,"digest_map","Digest recorded by storage for every bypass event object.","environment","Authenticated audit-store object metadata."),fact(recomputed,"digest_map","Locally recomputed digest of the exact retrieved bypass event bytes.","environment","Exact retrieved audit event bytes."),fact(modes,"state_map","Direct effective audit-storage enforcement mode for every event object.","environment","Read-only effective audit-store write and mutation policy."),fact(retention,"digest_map","Direct digest of the effective retention-lock configuration bound to every event object.","environment","Exact read-only effective audit retention configuration bytes.")]
    params=[parameter(required,"string_set","Complete review-bypass execution identities sealed independently from the assessed action inventory."),parameter(expected_actors,"identity_map","Expected administrator actor identity keyed by every sealed bypass execution."),parameter(expected_subjects,"identity_map","Expected repository and branch identity keyed by every sealed bypass execution."),parameter(approved_modes,"string_set","Independently approved raw enforcement modes that prohibit mutation and deletion."),parameter(expected_retention,"digest_map","Independently approved retention-lock policy digest keyed by every sealed bypass event.")]
    args=[{"op":"map_keys_eq_set_parameter","fact":x,"parameter":required} for x in [actors,subjects,times,stored,recomputed,modes,retention]]+[{"op":"identity_map_eq_parameter","fact":actors,"parameter":expected_actors},{"op":"identity_map_eq_parameter","fact":subjects,"parameter":expected_subjects},{"op":"digest_map_eq_fact","fact":stored,"other_fact":recomputed},{"op":"state_map_values_in_parameter","fact":modes,"parameter":approved_modes},{"op":"digest_map_eq_parameter","fact":retention,"parameter":expected_retention}]
    ids=["bypass-A","bypass-B"];actor_values={"bypass-A":"admin-A","bypass-B":"admin-B"};subject_values={"bypass-A":"repo-A.branch-main","bypass-B":"repo-A.branch-release"};time_values={"bypass-A":"2026-08-29T10:00:00Z","bypass-B":"2026-08-29T11:00:00Z"};digest_values={"bypass-A":"a"*64,"bypass-B":"b"*64};mode_values={x:"append-only-retention-locked" for x in ids};retention_values={x:"c"*64 for x in ids}
    pf={actors:fvalue("identity_map",actor_values),subjects:fvalue("identity_map",subject_values),times:fvalue("time_map",time_values),stored:fvalue("digest_map",digest_values),recomputed:fvalue("digest_map",digest_values),modes:fvalue("state_map",mode_values),retention:fvalue("digest_map",retention_values)};ff=copy.deepcopy(pf);ff[recomputed]=fvalue("digest_map",{"bypass-A":"d"*64,"bypass-B":"b"*64});pv={required:pvalue("string_set",ids),expected_actors:pvalue("identity_map",actor_values),expected_subjects:pvalue("identity_map",subject_values),approved_modes:pvalue("string_set",["append-only-retention-locked"]),expected_retention:pvalue("digest_map",retention_values)}
    result=exact_profile(facts,params,{"op":"all","args":args},pf,ff,pv,"Every bypass execution has an actor- and subject-bound event whose retrieved bytes match storage metadata under the approved append-only retention policy.","An event record is present and attributed, but its retrieved bytes no longer match the storage digest.")
    result["fixtures"]["counterexample"]=copy.deepcopy(result["fixtures"]["fail"]);result["fixtures"]["counterexample"]["description"]="A provider could label storage immutable and return an event identity after modifying its bytes; local digest recomputation against the stored digest rejects that event."
    return result


def simultaneous_metrics_profile():
    cid="USEQ-C8DBF731";ordinal=1;pre=prefix(cid,ordinal);required=f"{pre}.required_metric_ids";scopes=f"{pre}.metric_observation_scope_identities";streams=f"{pre}.metric_stream_identities";active=f"{pre}.metric_stream_active_states";expected_scopes=f"{pre}.expected_common_metric_scope_identities";required_true=f"{pre}.required_active_state"
    facts=[fact(scopes,"identity_map","Direct pipeline, view, and observation-window identity bound to each required metric kind.","environment","Read-only effective metric view and observation-window bindings."),fact(streams,"identity_map","Direct active metric stream identity for every required metric kind.","environment","Read-only effective telemetry stream configuration."),fact(active,"boolean_map","Direct effective active state for every required metric stream.","environment","Read-only effective telemetry stream state.")]
    params=[parameter(required,"string_set","The six required metric-kind identities sealed independently."),parameter(expected_scopes,"identity_map","Expected common pipeline, view, and observation-window identity keyed by every required metric kind."),parameter(required_true,"boolean","Required active metric-stream state true." )]
    metric_ids=["errors","latency-distribution","queue-age","resource-use","throughput","unit-cost"];scope_values={x:"pipeline-A.view-A.window-A" for x in metric_ids};stream_values={x:f"stream-{x}" for x in metric_ids};active_values={x:True for x in metric_ids}
    args=[{"op":"map_keys_eq_set_parameter","fact":x,"parameter":required} for x in [scopes,streams,active]]+[{"op":"identity_map_eq_parameter","fact":scopes,"parameter":expected_scopes},{"op":"boolean_map_all_eq_parameter","fact":active,"parameter":required_true}]
    pf={scopes:fvalue("identity_map",scope_values),streams:fvalue("identity_map",stream_values),active:fvalue("boolean_map",active_values)};ff=copy.deepcopy(pf);bad_scope=dict(scope_values);bad_scope["unit-cost"]="pipeline-A.view-B.window-A";ff[scopes]=fvalue("identity_map",bad_scope);pv={required:pvalue("string_set",metric_ids),expected_scopes:pvalue("identity_map",scope_values),required_true:pvalue("boolean",True)}
    return exact_profile(facts,params,{"op":"all","args":args},pf,ff,pv,"All six active raw metric streams share the exact sealed pipeline, view, and observation window.","One metric is active in a different view, so the metrics are not monitored together.")


def boolean_implies_profile(control_id, ordinal, authority, antecedent_suffix, consequent_suffix, antecedent_semantics, consequent_semantics):
    pre=prefix(control_id,ordinal);left=f"{pre}.{antecedent_suffix}";right=f"{pre}.{consequent_suffix}";required=f"{pre}.required_subject_ids"
    facts=[fact(left,"boolean_map",antecedent_semantics,authority,"Complete direct raw antecedent booleans keyed by the sealed subjects."),fact(right,"boolean_map",consequent_semantics,authority,"Complete direct raw consequent booleans keyed by the same sealed subjects.")]
    params=[parameter(required,"string_set","Complete stable subject identities sealed independently before collection.")]
    args=[{"op":"map_keys_eq_set_parameter","fact":left,"parameter":required},{"op":"map_keys_eq_set_parameter","fact":right,"parameter":required},{"op":"boolean_map_implies_fact","fact":left,"other_fact":right}]
    pf={left:fvalue("boolean_map",{"subject-A":True,"subject-B":False}),right:fvalue("boolean_map",{"subject-A":True,"subject-B":False})};ff=copy.deepcopy(pf);ff[right]=fvalue("boolean_map",{"subject-A":False,"subject-B":False})
    return exact_profile(facts,params,{"op":"all","args":args},pf,ff,{required:pvalue("string_set",["subject-A","subject-B"])},"For every sealed subject, a true antecedent has a true direct raw consequent.","One sealed subject has a true antecedent and false consequent.")


def combine_profiles(*profiles):
    facts=[];params=[];args=[];pf={};ff={};pv={}
    for profile in profiles:
        facts.extend(profile["facts"]);params.extend(profile["params"]);args.append(profile["predicate"])
        pf.update(profile["fixtures"]["pass"]["facts"]);ff.update(profile["fixtures"]["fail"]["facts"]);pv.update(profile["fixtures"]["pass"]["parameters"])
    return exact_profile(facts,params,{"op":"all","args":args},pf,ff,pv,"Every clause-specific raw sub-relation covers the complete sealed inventory and satisfies its exact relation.","At least one clause-specific raw sub-relation is broken even though the remaining signals look valid.")


def required_field_maps_profile(control_id, ordinal, authority, subject_suffix, fields, subject_semantics, subject_examples=None):
    pre=prefix(control_id,ordinal); subjects=f"{pre}.required_{subject_suffix}"
    facts=[];params=[parameter(subjects,"string_set",subject_semantics+" Sealed independently before record parsing.")];args=[]
    good_subjects=subject_examples or ["subject-A","subject-B"];pf={};ff={};pv={subjects:pvalue("string_set",good_subjects)}
    for index,(suffix,map_type,semantics) in enumerate(fields):
        key=f"{pre}.{suffix}";facts.append(fact(key,map_type,semantics,authority,"Direct lossless parsed values keyed by every subject in the independently sealed complete inventory."));args.append({"op":"map_keys_eq_set_parameter","fact":key,"parameter":subjects})
        if map_type=="digest_map": values={subject:format((i%8)+1,"x")*64 for i,subject in enumerate(good_subjects)}
        elif map_type=="boolean_map": values={subject:True for subject in good_subjects}
        elif map_type=="number_map": values={subject:i+1 for i,subject in enumerate(good_subjects)}
        elif map_type=="time_map": values={subject:f"2026-08-{20+i:02d}T10:00:00Z" for i,subject in enumerate(good_subjects)}
        else: values={subject:f"value-{i}" for i,subject in enumerate(good_subjects)}
        pf[key]=fvalue(map_type,values);ff[key]=fvalue(map_type,({subject:values[subject] for subject in good_subjects[:-1]} if index==0 else values))
    return exact_profile(facts,params,{"op":"all","args":args},pf,ff,pv,"Every required direct raw field map covers every subject in the independently sealed subject inventory.","At least one required direct raw field is absent for an independently inventoried subject.")


def number_profile(control_id,ordinal,authority,suffix,dimensions,direction="lte"):
    pre=prefix(control_id,ordinal); observed=f"{pre}.{suffix}"; limits=f"{pre}.approved_{suffix}"
    keys=[f"subject-A.{x}" for x in dimensions]
    facts=[fact(observed,"number_map",f"Raw numeric values in canonical units keyed by the full sealed subject-by-dimension cross product for: {', '.join(dimensions)}.",authority,"Authenticated raw measurements or read-only effective numeric configuration over every independently inventoried subject." )]
    params=[parameter(limits,"number_map",f"Approved numeric values, canonical units, and exact complete sealed subject-by-dimension keys for: {', '.join(dimensions)}.")]
    op=f"number_map_{direction}_parameter"; limit_values={x:10 for x in keys}
    if direction == "eq":
        good={x:10 for x in keys}; bad=dict(good); bad[keys[0]]=11
    elif direction in ("lte", "lt"):
        good={x:9 for x in keys}; bad=dict(good); bad[keys[0]]=11
    else:
        good={x:11 for x in keys}; bad=dict(good); bad[keys[0]]=9
    return exact_profile(facts,params,{"op":op,"fact":observed,"parameter":limits},{observed:fvalue("number_map",good)},{observed:fvalue("number_map",bad)},{limits:pvalue("number_map",limit_values)},"Every complete raw numeric dimension satisfies its independently sealed value.","One complete raw numeric dimension violates its independently sealed value.")


def analysis_run(control_id,ordinal,statement,threshold=False):
    pre=prefix(control_id,ordinal); subjects=f"{pre}.analyzed_subject_ids"; required=f"{pre}.required_subject_ids"; state=f"{pre}.raw_run_terminal_state"; analyzer=f"{pre}.analyzer_identity"; expected_analyzer=f"{pre}.required_analyzer_identity"; config=f"{pre}.analyzer_config_digest"; expected_config=f"{pre}.required_config_digest"
    facts=[fact(subjects,"string_set","Raw subject identities emitted in the analyzer input manifest.","executed","Authenticated analyzer input manifest."),fact(state,"state","Raw terminal state from the analyzer process record.","executed","Authenticated execution record, not a provider compliance verdict."),fact(analyzer,"identity","Authenticated analyzer executable identity and version.","executed","Analyzer invocation record."),fact(config,"digest","Digest of exact analyzer configuration and rules.","executed","Exact configuration bytes used by the run.")]
    params=[parameter(required,"string_set","Complete scanner-sealed subject inventory."),parameter(expected_analyzer,"identity","Approved analyzer identity and version."),parameter(expected_config,"digest","Approved configuration and rules digest.")]
    args=[{"op":"set_eq_parameter","fact":subjects,"parameter":required},{"op":"state_in","fact":state,"strings":["completed"]},{"op":"identity_eq_parameter","fact":analyzer,"parameter":expected_analyzer},{"op":"digest_eq_parameter","fact":config,"parameter":expected_config}]
    pf={subjects:fvalue("string_set",["subject-A","subject-B"]),state:fvalue("state","completed"),analyzer:fvalue("identity","analyzer-A@1"),config:fvalue("digest","c"*64)};ff=copy.deepcopy(pf);ff[subjects]=fvalue("string_set",["subject-A"])
    pv={required:pvalue("string_set",["subject-A","subject-B"]),expected_analyzer:pvalue("identity","analyzer-A@1"),expected_config:pvalue("digest","c"*64)}
    if threshold:
        findings=f"{pre}.normalized_finding_counts"; accepted=f"{pre}.accepted_finding_count_limits"
        facts.append(fact(findings,"number_map","Raw normalized finding counts keyed by every severity or rule dimension named by the approved profile.","executed","Lossless normalized analyzer findings from the complete run."));params.append(parameter(accepted,"number_map","Accepted maximum count per exact severity or rule dimension."));args.append({"op":"number_map_lte_parameter","fact":findings,"parameter":accepted});pf[findings]=fvalue("number_map",{"critical":0,"high":1});ff[findings]=fvalue("number_map",{"critical":1,"high":1});pv[accepted]=pvalue("number_map",{"critical":0,"high":1})
    return exact_profile(facts,params,{"op":"all","args":args},pf,ff,pv,f"Complete raw run bindings satisfy the exact clause: {statement}","The complete run omits an inventoried subject or violates an independently sealed run binding or threshold.")


def empty_diagnostics(control_id,ordinal,authority,diagnostic_suffix,include_digest=False):
    pre=prefix(control_id,ordinal); diagnostics=f"{pre}.{diagnostic_suffix}";analyzer=f"{pre}.analyzer_identity";expected_analyzer=f"{pre}.required_analyzer_identity";config=f"{pre}.analyzer_config_digest";expected_config=f"{pre}.required_analyzer_config_digest";subjects=f"{pre}.analyzed_subject_ids";required_subjects=f"{pre}.required_subject_ids"
    facts=[fact(diagnostics,"string_set","Raw stable diagnostic or matching-node identities emitted by the pinned parser/analyzer.",authority,"Pinned deterministic parser/analyzer output with source spans."),fact(analyzer,"identity","Direct analyzer executable identity and version used for this exact run.",authority,"Authenticated pinned parser or analyzer invocation record."),fact(config,"digest","Direct digest of the exact parser or analyzer configuration and rules used.",authority,"Exact pinned configuration and rule bytes used by the invocation."),fact(subjects,"string_set","Direct stable identities from the exact parser/analyzer input manifest.",authority,"Authenticated lossless input manifest for the pinned parser/analyzer invocation.")]
    params=[parameter(expected_analyzer,"identity","Independently approved analyzer executable identity and version."),parameter(expected_config,"digest","Independently approved exact analyzer configuration and rule digest."),parameter(required_subjects,"string_set","Complete independently sealed source or document subject inventory for this check.")]
    pred={"op":"all","args":[{"op":"identity_eq_parameter","fact":analyzer,"parameter":expected_analyzer},{"op":"digest_eq_parameter","fact":config,"parameter":expected_config},{"op":"set_eq_parameter","fact":subjects,"parameter":required_subjects},{"op":"set_eq","fact":diagnostics,"strings":[]}]};pf={diagnostics:fvalue("string_set",[]),analyzer:fvalue("identity","analyzer-A@1"),config:fvalue("digest","c"*64),subjects:fvalue("string_set",["src-A","src-B"])};ff={diagnostics:fvalue("string_set",["src/example:1:1"]),analyzer:fvalue("identity","analyzer-A@1"),config:fvalue("digest","c"*64),subjects:fvalue("string_set",["src-A","src-B"])};pv={expected_analyzer:pvalue("identity","analyzer-A@1"),expected_config:pvalue("digest","c"*64),required_subjects:pvalue("string_set",["src-A","src-B"])}
    if include_digest:
        left=f"{pre}.source_tree_digest";right=f"{pre}.regenerated_tree_digest";facts += [fact(left,"digest","Digest of exact sealed source tree bytes.",authority,"Complete source inventory."),fact(right,"digest","Digest of pinned-tool regenerated output tree.",authority,"Pinned formatter output.")];pred={"op":"all","args":[pred,{"op":"digest_eq_fact","fact":left,"other_fact":right}]};pf[left]=fvalue("digest","a"*64);pf[right]=fvalue("digest","a"*64);ff[left]=fvalue("digest","a"*64);ff[right]=fvalue("digest","a"*64)
    return exact_profile(facts,params,pred,pf,ff,pv,"The bound pinned analyzer emits no matching diagnostics and every required digest relation holds.","A raw diagnostic identifies a prohibited construct or deterministic regeneration mismatch.")


EXACT = {}
CURRENT_CLAUSES = {
    (binding["control_id"], clause["ordinal"])
    for binding in json.loads(BINDINGS.read_text())["bindings"]
    for clause in binding["clauses"]
}

def add(key, profile):
    if key in CURRENT_CLAUSES:
        EXACT[key] = profile

# Direct digest and identity relations.
for key,left,right,algo in [
    (("PRC-03-003",1),"release_record_digest","exact_artifact_bytes_digest",True),
]: add(key,digest_pair(*key,"artifact",left,right,"",algo))
add(("PRC-08-021",2),previous_artifact_registry_profile())

# Exact parsed/configured sets and subject-property coverage.
add(("PRC-08-005",1),set_parameter_profile("PRC-08-005",1,"environment","mandatory_build_gate_check_ids",["check-A","check-B"],"Raw check identifiers marked mandatory in effective build-gate configuration."))
add(("PRC-20-005",1),set_parameter_profile("PRC-20-005",1,"repository","explicitly_listed_allowed_type_ids",["type-A","type-B"],"Raw allowed type identifiers parsed from the authoritative configuration or document.","eq"))
add(("PRC-21-024",1),set_parameter_profile("PRC-21-024",1,"environment","effective_trusted_proxy_and_header_ids",["proxy-A","x-forwarded-for"],"Raw exact proxy identities and header names trusted by effective proxy configuration.","eq"))
add(("PRC-21-032",1),set_parameter_profile("PRC-21-032",1,"environment","inventoried_dns_record_ids",["dns-A","dns-B"],"Raw stable DNS record identities discovered from authoritative DNS data.","eq"))
add(("PRC-34-007",1),set_parameter_profile("PRC-34-007",1,"environment","available_recovery_point_ids",["recovery-A","recovery-B"],"Raw retrievable backup or recovery-point identities from the authoritative store."))
for cid,ordinal,authority,props,sems in [
    ("PRC-33-007",1,"environment",["encrypted_backup_ids"],{"encrypted_backup_ids":"Raw backup identities whose stored object encryption metadata names an approved encryption state."}),
]: add((cid,ordinal),property_coverage(cid,ordinal,authority,props,sems))

add(("PRC-16-024",1),state_map_allowed_profile("PRC-16-024",1,"environment","effective_account_states",["disabled"],"Direct effective account lifecycle state for every independently inventoried test, dormant, abandoned, or terminated account."))
add(("PRC-21-031",1),boolean_fields_inventory_profile("PRC-21-031",1,"environment",[("transfer_lock_enabled_by_domain","Direct effective registrar transfer-lock boolean for every domain."),("modification_protection_enabled_by_domain","Direct effective registrar modification-protection boolean for every domain.")],True,"in-scope domain identities"))
add(("PRC-28-023",1),boolean_fields_inventory_profile("PRC-28-023",1,"environment",[("cost_alert_enabled_by_resource_scope","Direct effective cost-alert enabled boolean for every resource scope."),("runaway_resource_protection_enabled_by_resource_scope","Direct effective runaway-resource protection enabled boolean for every resource scope.")],True,"resource-scope identities"))
add(("PRC-28-031",1),required_field_maps_profile("PRC-28-031",1,"environment","workload_resource_dimension_keys",[("configured_resource_request_values","number_map","Direct configured resource-request numeric value for every independently sealed workload-resource key."),("configured_resource_limit_values","number_map","Direct configured resource-limit numeric value for every independently sealed workload-resource key.")],"Complete workload x applicable resource-dimension keys sealed from the deployment inventory and project resource policy."))
add(("PRC-28-005",1),resource_label_profile())
add(("PRC-28-032",1),runtime_policy_enforcement_profile())
add(("PRC-29-030",1),pipeline_monitor_profile())

add(("PRC-28-018",1),two_sets("PRC-28-018",1,"environment","asset_inventory_resource_ids","deployed_infrastructure_resource_ids",("Raw stable asset inventory identities.","Raw stable identities enumerated from deployed infrastructure.")))
add(("PRC-37-007",1),two_sets("PRC-37-007",1,"environment","sandbox_account_ids","production_account_ids",("Raw sandbox account identities.","Raw production account identities."),"set_disjoint_fact"))

# Numeric thresholds and bounds.
for cid,ord_,auth,suffix,dims,direction in [
    ("PRC-25-016",2,"executed","cold_start_measurements",["cold_start_p95_ms"],"lte"),
    ("PRC-25-024",1,"repository","configured_regression_thresholds",["latency_regression_pct","error_regression_pct"],"eq"),
    ("PRC-26-010",2,"executed","capacity_fixture_dimensions",["storage_gib","requests_per_second"],"gte"),
    ("PRC-26-016",2,"executed","remaining_capacity_by_scenario",["node_loss","region_loss","zone_loss"],"gte"),
    ("PRC-27-037",1,"executed","tested_recovery_seconds",["recovery_seconds"],"lte"),
    ("PRC-33-021",1,"executed","tested_recovery_seconds",["recovery_seconds"],"lte"),
    ("PRC-33-023",1,"executed","tested_dr_capacity",["requests_per_second"],"gte"),
    ("USEQ-427C1469",1,"environment","availability_by_objective",["critical_journey_availability_pct","api_availability_pct"],"gte"),
    ("USEQ-4B1B01AF",1,"environment","effective_worker_queue_limits",["attempts","age_seconds","execution_seconds","concurrency","queue_depth","payload_bytes","memory_bytes"],"lte"),
    ("USEQ-64AE632C",1,"environment","effective_storage_query_limits",["memory_bytes","disk_bytes","entry_bytes","result_bytes","query_complexity","retention_seconds"],"lte"),
    ("USEQ-685621E2",1,"environment","effective_file_processing_limits",["file_bytes","file_count","archive_depth","expanded_bytes","dimensions","page_count","processing_seconds"],"lte"),
    ("USEQ-D4E7515B",1,"executed","tested_recovery_objectives",["recovery_time_seconds","recovery_point_seconds"],"lte"),
    ("USEQ-E87054DC",1,"environment","effective_request_limits",["request_bytes","url_bytes","header_bytes","nesting_depth","field_count","query_complexity","batch_size"],"lte"),
]: add((cid,ord_),number_profile(cid,ord_,auth,suffix,dims,direction))

# Analyzer executions and raw diagnostic sets.
for cid in ["PRC-31-003","PRC-31-004","PRC-31-005","PRC-31-006","PRC-31-007","PRC-31-008","PRC-31-009","PRC-31-010","PRC-31-011","PRC-31-012","PRC-31-013","PRC-31-014","PRC-31-015"]:
    add((cid,1),analysis_run(cid,1,"analysis has run"))
add(("PRC-14-002",2),analysis_run("PRC-14-002",2,"accessibility scan scope and thresholds",True))
add(("USEQ-898B0514",2),analysis_run("USEQ-898B0514",2,"bound scan scope and thresholds",True))
add(("PRC-09-016",2),time_delta_profile("PRC-09-016",2,"executed","vulnerability_disclosure_times","included_or_triggered_scan_times","approved_maximum_scan_cadence_seconds","Direct disclosure timestamp for every supported-release and newly disclosed vulnerability relation.","Direct timestamp when the exact vulnerability was included in or triggered a scan of the exact supported release."))
add(("PRC-31-021",1),fixed_finding_retest_profile())
add(("USEQ-7E74EFCD",1),state_predicate_wait_profile())
add(("USEQ-0E8F2EB1",1),combine_profiles(
    pinned_tool_profile("USEQ-0E8F2EB1",1,"executed","formatter_identity","formatter_configuration_digest","Direct exact formatter executable identity and version used for regeneration.","Direct digest of the exact repository formatter configuration used for regeneration."),
    digest_pair("USEQ-0E8F2EB1",1,"executed","source_tree_digest","pinned_formatter_output_tree_digest","All assessed source bytes must equal regeneration by the independently approved formatter and configuration.",False),
))
add(("USEQ-27E5FAAF",1),empty_diagnostics("USEQ-27E5FAAF",1,"repository","jumpable_wall_clock_elapsed_subtraction_node_ids"))
add(("USEQ-DAF77C8F",1),empty_diagnostics("USEQ-DAF77C8F",1,"repository","prohibited_construct_diagnostic_ids"))
add(("USEQ-FE04149F",1),empty_diagnostics("USEQ-FE04149F",1,"repository","prohibited_style_construct_ids",True))

# Documentation and exact set/state checks.
for cid,vals,suffix,meaning in [
    ("PRC-36-002",["architecture","data-flow","dependencies","deployment","recovery"],"documented_topic_content","architecture topic content"),
    ("PRC-36-004",["build-command","test-command"],"documented_command_text","literal runnable command text"),
    ("PRC-36-005",["coding","review","testing","release"],"documented_convention_content","convention content"),
]: add((cid,1),nonempty_string_map_profile(cid,1,"repository",suffix,vals,f"Direct parsed nonempty {meaning} keyed by its exact required topic.",f"Exact required topic identities sealed from the control statement: {', '.join(vals)}."))
add(("USEQ-3726C504",1),empty_diagnostics("USEQ-3726C504",1,"repository","markdown_validator_error_ids"))
add(("USEQ-69B5C63D",1),versioned_repository_items_profile())
add(("USEQ-9E10A498",1),cryptographic_requirements_profile())

# Source inventory/property relations. A pinned analyzer must cover the complete
# sealed source inventory; provider-selected "resolved" sets are not sufficient.
add(("USEQ-589E6801",1),combine_profiles(
    required_field_maps_profile("USEQ-589E6801",1,"repository","injected_binding_ids",[("resolved_owner_identities","identity_map","Direct syntax- and framework-resolved owner identity for every injected binding."),("resolved_lifetime_or_scope_identities","identity_map","Direct syntax- and framework-resolved lifetime or scope identity for every injected binding.")],"Complete independently sealed injected-binding identities."),
    empty_diagnostics("USEQ-589E6801",1,"repository","unsupported_or_unresolved_dynamic_binding_diagnostic_ids"),
))
add(("USEQ-7414620B",1),combine_profiles(
    required_field_maps_profile("USEQ-7414620B",1,"repository","protocol_surface_ids",[("versioned_schema_identities","identity_map","Direct analyzer-resolved explicit versioned schema identity for every request, response, event, and error surface."),("versioned_compatibility_rule_set_identities","identity_map","Direct analyzer-resolved explicit versioned compatibility-rule set identity for every surface.")],"Complete independently sealed repository-derived request, response, event, and error surface identities."),
    empty_diagnostics("USEQ-7414620B",1,"repository","unresolved_surface_or_binding_diagnostic_ids"),
))

# Complete parsed documentation, source, package, and container field coverage.
for cid,values,suffix in [
    ("USEQ-D083AE3C",["access","alerts","capacity","continuity","cost","observability","privacy","runbooks","security","support"],"acceptance_criteria_topic_ids"),
    ("USEQ-DDCBBEB7",["format-migration","media-refresh","post-action-verification","re-encryption","re-signing","timestamp-renewal"],"preservation_action_ids"),
]: add((cid,1),set_parameter_profile(cid,1,"repository",suffix,values,f"Direct stable identifiers parsed from the current versioned document for every named {suffix.replace('_',' ')}.","contains"))
add(("USEQ-D083AE3C",1),nonempty_string_map_profile("USEQ-D083AE3C",1,"repository","acceptance_criteria_content_by_topic",["access","alerts","capacity","continuity","cost","observability","privacy","runbooks","security","support"],"Direct independently parsed acceptance-criterion content for each named topic.","The ten exact acceptance-criteria topic identities sealed from the clause."))

for cid,subjects,fields in [
    ("USEQ-3EAF96A6","incident_playbook_ids",[(x,"identity_map",f"Direct parsed nonempty {x.replace('_',' ')} field for each incident playbook.") for x in ["detection_steps","containment_steps","recovery_steps","validation_steps","communication_steps","escalation_steps","owner_identity","evidence_preservation_steps"]]),
    ("USEQ-4D94F1DC","cryptographic_incident_playbook_ids",[(x,"identity_map",f"Direct parsed nonempty {x.replace('_',' ')} field for each cryptographic incident playbook.") for x in ["detection_steps","containment_steps","authority_identity","revocation_or_migration_steps","recovery_steps","validation_steps","communication_steps","escalation_steps","evidence_preservation_steps"]]),
    ("USEQ-7DDE7DE6","key_response_procedure_ids",[(x,"identity_map",f"Direct parsed nonempty {x.replace('_',' ')} field for each key response procedure.") for x in ["detection_steps","containment_steps","revocation_or_disablement_steps","replacement_steps","dependent_update_steps","recovery_steps","verification_steps","escalation_steps","evidence_preservation_steps","accountable_role_identity"]]),
    ("USEQ-88DC1597","pipeline_rollback_procedure_ids",[(x,"identity_map",f"Direct parsed nonempty {x.replace('_',' ')} field for each pipeline rollback or roll-forward procedure.") for x in ["source_version_identity","target_version_identity","prerequisite_steps","command_identity","authority_identity","verification_steps","failure_handling_steps"]]),
    ("USEQ-DCBA84F8","pipeline_runbook_ids",[(x,"identity_map",f"Direct parsed nonempty {x.replace('_',' ')} field for each pipeline runbook.") for x in ["detection_steps","containment_steps","correction_or_replay_steps","integrity_verification_steps","escalation_steps","owner_identity","recovery_steps","pipeline_version_identity"]]),
]: add((cid,1),required_field_maps_profile(cid,1,"repository",subjects,fields,f"Complete parsed stable identities for every required {subjects.replace('_',' ')}."))

add(("USEQ-B086FE62",1),interface_binding_profile())

add(("USEQ-DDCBBEB7",1),required_field_maps_profile("USEQ-DDCBBEB7",1,"repository","preservation_action_ids",[(x,"identity_map",f"Direct parsed nonempty {x.replace('_',' ')} for every applicable preservation action.") for x in ["interval_identity","owner_identity","subject_inventory_identity","verification_step_set_identity"]],"Complete parsed stable identities for every selected preservation action."))

add(("PRC-20-009",1),empty_diagnostics("PRC-20-009",1,"repository","user_filename_used_as_trusted_path_node_ids"))
add(("USEQ-ED684D68",1),combine_profiles(
    required_field_maps_profile("USEQ-ED684D68",1,"repository","protocol_surface_ids",[(x,"identity_map",f"Direct analyzer-resolved {x.replace('_',' ')} for every protocol surface.") for x in ["contract_identity","field_set_identity","type_set_identity","requiredness_contract_identity","outcome_contract_identity","compatibility_identity"]],"Complete independently sealed syntax-derived request, command, query, event, and error producer and consumer surface identities."),
    empty_diagnostics("USEQ-ED684D68",1,"repository","unresolved_protocol_surface_or_contract_diagnostic_ids"),
))

add(("PRC-09-005",1),sbom_exact_profile())
add(("PRC-09-012",1),map_parameter_profile("PRC-09-012",1,"repository","identity_map","resolved_dependency_registry_identities","Direct registry identity for every resolved dependency from the exact lockfile or resolver record in repository authority."))
add(("PRC-09-012",2),map_parameter_profile("PRC-09-012",2,"external_registry","identity_map","resolved_dependency_publisher_identities","Direct authoritative publisher identity for every resolved dependency."))

add(("PRC-28-030",1),privileged_approval_profile())

# Effective CI policy raw fields, not provider-supplied conclusion booleans.
add(("PRC-07-005",1),branch_policy_profile("PRC-07-005",1,"minimum_reviews"))
add(("PRC-07-005",2),branch_policy_profile("PRC-07-005",2,"bypass_denied"))
add(("PRC-07-009",1),branch_policy_profile("PRC-07-009",1,"mandatory_checks"))
add(("PRC-07-013",1),branch_policy_profile("PRC-07-013",1,"rewrite_denied"))

add(("PRC-07-008",1),keyed_map_coverage_profile("PRC-07-008",1,"environment","identity_map","immutable_attributable_audit_event_ids","Direct immutable audit-event identity attributed to the administrator for every review-bypass path execution."))
add(("USEQ-0B41367E",1),protected_change_policy_profile())
add(("USEQ-1D76E3B6",1),required_field_maps_profile("USEQ-1D76E3B6",1,"environment","asynchronous_operation_type_ids",[(x,"identity_map",f"Direct effective nonempty {x.replace('_',' ')} behavior binding for every asynchronous operation type.") for x in ["status_contract_identity","cancellation_contract_identity","retry_contract_identity","timeout_contract_identity","reconciliation_contract_identity"]],"Complete effective asynchronous operation type identities."))
add(("USEQ-2210FD4E",1),map_equal_profile("USEQ-2210FD4E",1,"environment","identity_map","product_or_process_change_workflow_run_ids","documentation_change_workflow_run_ids","Direct workflow-run identity for every sealed product or process change.","Direct workflow-run identity containing the documentation change for the same sealed change."))
add(("USEQ-48587D52",1),time_order_profile("USEQ-48587D52",1,"environment","required_verification_success_times","artifact_signing_times","Direct successful required-verification timestamp for every signing event.","Direct artifact-signing timestamp for the same event."))
add(("USEQ-517E9DBD",1),separate_release_authorities_profile())
add(("USEQ-97E5A6C1",1),policy_negative_cases_profile("USEQ-97E5A6C1","critical_failure_input_states","unrelated_control_score_vector_digests",["critical-fail-high-other-scores","critical-fail-low-other-scores"],"Direct critical-failure input boolean supplied to the pinned evaluator for every case.","Direct canonical digest of unrelated control scores varied across each case."))
add(("USEQ-9F7B27E7",1),policy_negative_cases_profile("USEQ-9F7B27E7","critical_journey_inaccessible_input_states","low_value_page_result_vector_digests",["journey-inaccessible-many-pages-pass","journey-inaccessible-no-pages-pass"],"Direct critical-journey inaccessible input boolean supplied to the pinned evaluator for every case.","Direct canonical digest of varied low-value page results for every case."))
add(("USEQ-B9F93CB2",1),mixed_version_policy_profile())

add(("USEQ-A8A37E99",1),combine_profiles(
    pinned_tool_profile("USEQ-A8A37E99",1,"environment","gate_policy_evaluator_identity","gate_policy_evaluator_configuration_digest","Direct pinned deterministic gate-policy evaluator identity and version.","Direct digest of the exact effective aggregation expression and evaluator rules."),
    scalar_parameter_profile("USEQ-A8A37E99",1,"environment","digest","effective_release_gate_policy_digest","Direct digest of the exact effective release-gate policy bytes."),
    expected_field_maps_profile("USEQ-A8A37E99",1,"environment","gate_kind_ids",[("effective_gate_identities_by_kind","identity_map","Direct separate effective gate identity keyed by integrity, safety, and availability gate kind.")],"The exact required integrity, safety, and availability gate-kind identities.",["integrity","safety","availability"]),
    expected_field_maps_profile("USEQ-A8A37E99",1,"environment","aggregation_case_ids",[("aggregation_fixture_input_digests","digest_map","Direct exact aggregation input digest for every independently sealed integrity-, safety-, and availability-value fixture."),("integrity_or_safety_failure_states","boolean_map","Direct raw integrity-or-safety failure boolean for every sealed aggregation fixture."),("aggregate_failure_states","boolean_map","Direct raw aggregate failed-decision boolean for every sealed aggregation fixture.")],"Complete independently versioned aggregation cases that vary availability while integrity or safety fails.",["integrity-fails-high-availability","safety-fails-high-availability"]),
))
add(("USEQ-4321C4D8",1),combine_profiles(
    required_field_maps_profile("USEQ-4321C4D8",1,"environment","supplier_product_ids",[(x,"identity_map",f"Direct submitted nonempty {x.replace('_',' ')} field for every supplier product.") for x in ["algorithm_set_identity","library_set_identity","hardware_root_identity","certificate_set_identity","protocol_set_identity","key_custody_identity","validation_status_identity","update_mechanism_identity"]],"Complete effective supplier product identities."),
    case_state_profile("USEQ-4321C4D8",1,"environment","missing_disclosure_field_gate_states","Direct gate terminal state for every sealed fixture omitting one applicable supplier disclosure field."),
))
add(("USEQ-CF1B0C98",1),exception_gate_profile())
add(("USEQ-E8367B76",1),release_signature_profile())

# Spot-audited completeness and raw-relation overrides.
add(("PRC-07-008",1),review_bypass_audit_profile())
add(("PRC-33-007",1),state_map_allowed_profile("PRC-33-007",1,"environment","backup_encryption_states",["aes256-gcm","provider-managed-approved"],"Direct stored-object encryption algorithm or state identifier for every independently inventoried backup."))
add(("USEQ-C8DBF731",1),simultaneous_metrics_profile())
add(("USEQ-E87054DC",1),combine_profiles(
    pinned_tool_profile("USEQ-E87054DC",1,"environment","request_policy_canonicalizer_identity","request_policy_canonicalizer_configuration_digest","Direct pinned deterministic effective-request-policy canonicalizer identity and version.","Direct digest of canonicalizer rules for content-type and encoding sets."),
    number_profile("USEQ-E87054DC",1,"environment","effective_request_limits",["request_bytes","url_bytes","header_bytes","nesting_depth","field_count","query_complexity","batch_size"],"lte"),
    expected_field_maps_profile("USEQ-E87054DC",1,"environment","entrypoint_ids",[("effective_allowed_content_type_set_digests","digest_map","Direct canonical digest of the exact allowed content-type set at every entry point."),("effective_allowed_encoding_set_digests","digest_map","Direct canonical digest of the exact allowed content-encoding set at every entry point.")],"Complete independently sealed in-scope request entry-point identities.",["entrypoint-A","entrypoint-B"]),
))

# Explicitly audited bounded-execution domains. Each entry is intentionally assigned
# to a concrete raw terminal-state domain; this is not keyword or family routing.
EXECUTION_CASE_GROUPS = {
    "release_gate_terminal_states": """
PRC-08-005:2 PRC-08-019:1 PRC-25-024:2 PRC-34-015:2 PRC-37-018:1 PRC-42-015:2
USEQ-01EEF7ED:1 USEQ-06A9124D:1 USEQ-1F838C1D:1 USEQ-4321C4D8:2 USEQ-48996A09:1 USEQ-4D20534E:1
USEQ-5797F11A:1 USEQ-609408D4:1 USEQ-61A0F909:1 USEQ-8574E516:1 USEQ-9B3A9AA9:1 USEQ-A0338698:1
USEQ-A4A4897A:1 USEQ-B0B7E4B2:1 USEQ-CF1B0C98:2 USEQ-D63DD56F:1 USEQ-E0A45D6A:1 USEQ-ECC22163:1
USEQ-E8367B76:2
""",
    "security_and_authorization_case_states": """
PRC-16-010:2 PRC-16-027:2 PRC-21-012:1 PRC-28-024:1 PRC-35-007:1
USEQ-36E65D58:1 USEQ-4F9A4886:1 USEQ-4D94F1DC:2 USEQ-5FBAC758:1 USEQ-7DDE7DE6:2
USEQ-91D01B47:1 USEQ-9D7169D2:2 USEQ-B0B4CA24:1 USEQ-BBED00FD:2 USEQ-E87054DC:2
USEQ-ED684D68:2 USEQ-EFB39081:2 USEQ-F49BC55A:1
""",
    "recovery_and_resilience_case_states": """
PRC-30-011:1 PRC-34-006:1 PRC-34-018:1
USEQ-0AA61EFD:1 USEQ-23FADC9E:1 USEQ-2937AC0D:1 USEQ-42F63C13:1 USEQ-4B1B01AF:2
USEQ-4D92636C:2 USEQ-5218074C:2 USEQ-52FBF8D0:1 USEQ-64AE632C:2 USEQ-685621E2:2
USEQ-6B3537D2:2 USEQ-6E36BDB8:1 USEQ-719E4479:1 USEQ-84C1C16C:1
USEQ-87DCD770:1 USEQ-88DC1597:2 USEQ-8E08DDE8:1 USEQ-9E774722:1 USEQ-DCBA84F8:2
USEQ-E1B20661:1 USEQ-E4B21F0D:1 USEQ-F471F071:1 USEQ-FBA075B3:1 USEQ-FDD640C1:1
""",
    "interface_and_user_case_states": """
PRC-24-029:2 PRC-32-006:1
USEQ-0291C731:1 USEQ-376AA3EF:1 USEQ-53280330:1 USEQ-68334A49:1 USEQ-7414620B:2
USEQ-79FA79CC:1 USEQ-7AD07D0F:2 USEQ-8A50C66F:2 USEQ-9BC1EE8B:2 USEQ-AC43E131:1
USEQ-B9F93CB2:2 USEQ-BE07FB46:2 USEQ-C27619F9:2 USEQ-E35FCCCE:1 USEQ-ED0206BA:2
USEQ-F6E4CD10:2
""",
    "test_and_tool_case_states": """
USEQ-0CDDDBAD:1 USEQ-078B9EDE:1 USEQ-1223A7F4:1 USEQ-4488AC37:1 USEQ-49CDA8B6:1
USEQ-54517178:1 USEQ-57E6BCF0:1 USEQ-6A1FC94F:1 USEQ-857E50CE:1 USEQ-9993C108:1
USEQ-9CE7EDE5:1 USEQ-9EA30976:1 USEQ-9FC91F6A:1 USEQ-A6A392BF:1 USEQ-C07FE87E:1
USEQ-CB1343BF:1 USEQ-CBC5C76C:1 USEQ-CC592135:1 USEQ-D2D9767A:1 USEQ-DD834B76:2
USEQ-ED0E5F15:1 USEQ-F4B43A58:1
""",
    "data_integrity_and_reconciliation_case_states": """
USEQ-14B6379E:1 USEQ-1CA9F46F:2 USEQ-39C82E33:1 USEQ-8C9DCBEF:1 USEQ-CB5A99DA:1
USEQ-E12B5B11:1 USEQ-E7843EF8:1
""",
}

EXECUTION_MEASUREMENT_KEYS = """
USEQ-03172925:1 USEQ-2834E5D4:1 USEQ-34C19231:1 USEQ-789C8367:1 USEQ-82AC181E:1
USEQ-A84EEA72:1 USEQ-A936D8F1:1 USEQ-C262BE4B:1 USEQ-D2F05431:1 USEQ-D3926FEF:1 USEQ-DFD85174:1
"""

def parse_key_tokens(value):
    result=[]
    for token in value.split():
        cid,ordinal=token.split(":");result.append((cid,int(ordinal)))
    return result

_binding_statements={(b["control_id"],c["ordinal"]):c["statement"] for b in json.loads(BINDINGS.read_text())["bindings"] for c in b["clauses"]}
for suffix,tokens in EXECUTION_CASE_GROUPS.items():
    for key in parse_key_tokens(tokens):
        if key not in EXACT:
            add(key,case_state_profile(key[0],key[1],"executed",suffix,"Direct raw terminal output, state, rejection, transition, or emitted result keyed by every independently sealed case required by this clause: "+_binding_statements[key]))
for key in parse_key_tokens(EXECUTION_MEASUREMENT_KEYS):
    if key in _binding_statements and key not in EXACT:
        add(key,keyed_map_coverage_profile(key[0],key[1],"executed","number_map","raw_measurements_by_sealed_measurement_id","Direct raw numeric values in the clause-declared units and windows keyed by every independently sealed subject and measurement identity: "+_binding_statements[key]))

MEASUREMENT_DIMENSIONS = {
    ("USEQ-03172925",1): (["feedback_latency_approved_percentile","queue_time"],["feedback-A|feedback_latency_p95","feedback-A|queue_time"],"feedback-source x approved percentile-or-queue-time"),
    ("USEQ-34C19231",1): (["handshake_time","signing_time","verification_time","bandwidth","storage","memory","cpu","energy","queue_time","timeout_count","denial_of_service_impact"],[f"workload-A|participant-A|{dimension}" for dimension in ["handshake_time","signing_time","verification_time","bandwidth","storage","memory","cpu","energy","queue_time","timeout_count","denial_of_service_impact"]],"cryptographic-workload x participant-profile x measurement"),
    ("USEQ-789C8367",1): (["lock_duration","compute_use","io","replication_impact","storage_impact","total_duration"],[f"change-A|workload-A|{dimension}" for dimension in ["lock_duration","compute_use","io","replication_impact","storage_impact","total_duration"]],"change x workload-scenario x measurement"),
    ("USEQ-82AC181E",1): (["blocked_time","queue_time","rework","dependency_age","decision_latency"],[f"flow-A|{dimension}" for dimension in ["blocked_time","queue_time","rework","dependency_age","decision_latency"]],"delivery-flow subject x measurement"),
    ("USEQ-A84EEA72",1): (["change_lead_time","deployment_frequency","failed_deployment_recovery_time","change_fail_rate","deployment_rework_rate"],[f"service-A|{dimension}" for dimension in ["change_lead_time","deployment_frequency","failed_deployment_recovery_time","change_fail_rate","deployment_rework_rate"]],"team-or-service x DORA measurement"),
    ("USEQ-A936D8F1",1): (["recovery_point","recovery_time","data_loss","duplicate_effects","service_degradation"],[f"recovery-exercise-A|{dimension}" for dimension in ["recovery_point","recovery_time","data_loss","duplicate_effects","service_degradation"]],"recovery exercise x measurement"),
    ("USEQ-C262BE4B",1): (["end_to_end_waiting_time","queue_time"],["flow-A|end_to_end_waiting_time","flow-A|queue_time"],"delivery-flow subject x waiting-time measurement"),
    ("USEQ-D2F05431",1): (["duration","resource_use","locking","service_impact","recovery_time"],[f"change-A|{dimension}" for dimension in ["duration","resource_use","locking","service_impact","recovery_time"]],"change-scenario x measurement"),
    ("USEQ-DFD85174",1): (["flaky_outcome_count_and_total_run_count_by_test","flaky_outcome_count_and_total_run_count_by_suite","flaky_outcome_count_and_total_run_count_by_environment","flaky_outcome_count_and_total_run_count_by_branch","flaky_outcome_count_and_total_run_count_by_dependency","flaky_outcome_count_and_total_run_count_by_failure_signature"],[f"{group}-A|{count}" for group in ["test","suite","environment","branch","dependency","failure_signature"] for count in ["flaky_outcome_count","total_run_count"]],"sealed grouping-bucket x raw count kind"),
}
for _measurement_key,(_dimensions,_examples,_scope) in MEASUREMENT_DIMENSIONS.items():
    add(_measurement_key,measurement_coverage_profile(_measurement_key[0],_measurement_key[1],_dimensions,_examples,_scope))

add(("USEQ-F00285A1",1),measurement_coverage_profile("USEQ-F00285A1",1,["every independently approved latency percentile","tail_maximum_or_approved_tail_distribution_statistic"],["workload-A|journey-A|p50","workload-A|journey-A|p95","workload-A|journey-A|p99","workload-A|journey-A|tail_max"],"latency-workload x user-journey x approved-statistic"))
add(("USEQ-57378001",1),expected_field_maps_profile("USEQ-57378001",1,"executed","slo_consumption_scenario_ids",[("scenario_input_digests","digest_map","Direct digest of each independently sealed acute-fast-burn or chronic-slow-burn input fixture."),("emitted_alert_identities","identity_map","Direct emitted alert event identity for every sealed consumption scenario."),("detection_terminal_states","state_map","Direct raw detection state for every sealed consumption scenario.")],"Exactly the independently approved acute-fast-burn and chronic-slow-burn scenario identities sealed before execution.",["acute-fast-burn","chronic-slow-burn"]))

# Terminal states cannot substitute for a clause-specific oracle. These audited case
# assignments remain research routing only until a definition below adds every named
# raw result field and exact relation required by that clause.
for _case_group_tokens in EXECUTION_CASE_GROUPS.values():
    for _case_key in parse_key_tokens(_case_group_tokens):
        EXACT.pop(_case_key,None)
for _mixed_measurement_key in [("USEQ-2834E5D4",1),("USEQ-D3926FEF",1)]:
    EXACT.pop(_mixed_measurement_key,None)

add(("USEQ-7330C584",1),required_field_maps_profile("USEQ-7330C584",1,"executed","test_suite_ids",[("setup_start_times","time_map","Direct authenticated setup start timestamp for every suite."),("setup_finish_times","time_map","Direct authenticated setup finish timestamp for every suite."),("setup_resource_and_action_manifest_identities","identity_map","Direct exact resource and action manifest identity for every suite under the sealed complexity measure."),("computed_setup_duration_seconds","number_map","Deterministically normalized setup duration computed from the raw start and finish timestamps."),("computed_setup_complexity_values","number_map","Deterministically normalized complexity value computed from the raw resource and action manifest under the sealed measure.")],"Complete test-suite identities sealed independently before execution."))
add(("USEQ-9754BE06",1),required_field_maps_profile("USEQ-9754BE06",1,"executed","exercise_ids",[("communication_activity_identities","identity_map","Direct exercised communication activity identities for every exercise."),("escalation_activity_identities","identity_map","Direct exercised escalation activity identities for every exercise."),("evidence_preservation_activity_identities","identity_map","Direct exercised evidence-preservation activity identities for every exercise."),("customer_impact_activity_identities","identity_map","Direct exercised customer-impact activity identities for every exercise."),("exercise_terminal_states","state_map","Direct raw outcome state for every exercise.")],"Complete exercise identities sealed independently before execution."))
add(("USEQ-BD7F03F5",1),exact_test_trace_profile())
add(("USEQ-EA5C03BB",1),combine_profiles(
    required_field_maps_profile("USEQ-EA5C03BB",1,"executed","quarantined_test_ids",[("quarantine_start_times","time_map","Direct quarantine start timestamp for every quarantined test."),("quarantine_end_times","time_map","Direct removal or repair completion timestamp for every expired quarantined test.")],"Complete quarantined-test identities sealed independently from test history."),
    case_state_profile("USEQ-EA5C03BB",1,"executed","expired_quarantine_terminal_states","Direct removed or repaired terminal state for every sealed expired-quarantine case."),
))

def reforecast_profile():
    cid="USEQ-EEBE0AC0";ordinal=1;pre=prefix(cid,ordinal);old=f"{pre}.obsolete_forecast_dates";new=f"{pre}.reforecast_dates";throughput=f"{pre}.actual_throughput_measurements";work=f"{pre}.discovered_work_measurements";required=f"{pre}.required_forecast_ids"
    facts=[fact(old,"time_map","Direct prior forecast date for every sealed forecast identity.","executed","Authenticated prior forecast records."),fact(new,"time_map","Direct newly computed forecast date for every sealed forecast identity.","executed","Authenticated reforecast output records."),fact(throughput,"number_map","Direct actual throughput used by every reforecast.","executed","Authenticated actual completed-work measurements."),fact(work,"number_map","Direct discovered remaining-work value used by every reforecast.","executed","Authenticated discovered-work inventory and measurement.")]
    params=[parameter(required,"string_set","Complete forecast identities sealed independently before collection.")];args=[{"op":"map_keys_eq_set_parameter","fact":x,"parameter":required} for x in [old,new,throughput,work]];args.append({"op":"not","arg":{"op":"all","args":[{"op":"time_map_lte_fact","fact":old,"other_fact":new},{"op":"time_map_gte_fact","fact":old,"other_fact":new}]}})
    pv={required:pvalue("string_set",["forecast-A","forecast-B"])};pf={old:fvalue("time_map",{"forecast-A":"2026-09-01T00:00:00Z","forecast-B":"2026-09-02T00:00:00Z"}),new:fvalue("time_map",{"forecast-A":"2026-09-03T00:00:00Z","forecast-B":"2026-09-04T00:00:00Z"}),throughput:fvalue("number_map",{"forecast-A":8,"forecast-B":9}),work:fvalue("number_map",{"forecast-A":80,"forecast-B":90})};ff=copy.deepcopy(pf);ff[new]=copy.deepcopy(pf[old])
    return exact_profile(facts,params,{"op":"all","args":args},pf,ff,pv,"Every reforecast binds actual throughput and discovered work and changes the obsolete forecast dates.","Actual inputs are present but the obsolete forecast dates were preserved unchanged.")
add(("USEQ-EEBE0AC0",1),reforecast_profile())
add(("USEQ-F2601113",1),combine_profiles(
    required_field_maps_profile("USEQ-F2601113",1,"executed","calculation_ids",[("raw_measurement_artifact_digests","digest_map","Direct digest of the complete raw measurement artifact for every calculation."),("calculation_formula_identities","identity_map","Direct pinned formula and tool identity for every calculation."),("calculation_configuration_digests","digest_map","Direct digest of exact calculation configuration for every calculation.")],"Complete calculation identities sealed independently before execution."),
    map_equal_profile("USEQ-F2601113",1,"executed","number_map","recorded_calculation_outputs","independently_replayed_calculation_outputs","Direct recorded numeric output for every calculation.","Direct independently replayed numeric output from the retained raw measurements, formula, and configuration."),
))
for _unsupported_exact_key in [("USEQ-EA5C03BB",1),("USEQ-EEBE0AC0",1),("USEQ-F2601113",1)]:
    EXACT.pop(_unsupported_exact_key,None)

# Multi-digest artifact bindings that fit the existing DSL exactly.
def multi_digest(control_id,ordinal,pairs):
    pre=prefix(control_id,ordinal);facts=[];args=[];pf={};ff={}
    for i,(left_suffix,right_suffix,meaning) in enumerate(pairs):
        left=f"{pre}.{left_suffix}";right=f"{pre}.{right_suffix}";facts += [fact(left,"digest",f"Raw digest for {left_suffix.replace('_',' ')}.","artifact",meaning),fact(right,"digest",f"Independently obtained digest for {right_suffix.replace('_',' ')}.","artifact",meaning)];args.append({"op":"digest_eq_fact","fact":left,"other_fact":right});pf[left]=fvalue("digest","a"*64);pf[right]=fvalue("digest","a"*64);ff[left]=fvalue("digest","a"*64);ff[right]=fvalue("b" if False else "digest","b"*64 if i==0 else "a"*64)
    return exact_profile(facts,[],{"op":"all","args":args},pf,ff,{},"Every independent artifact and configuration digest binding matches.","At least one complete artifact or configuration digest binding differs.")
add(("USEQ-CA466120",1),merged_review_verification_profile())
add(("USEQ-4C2E2D56",1),deployed_reviewed_tested_profile())

# The current closed DSL cannot prove computed duration = finish-start or evaluate
# the independently versioned complexity formula over raw setup actions/resources.
EXACT.pop(("USEQ-7330C584",1),None)
add(("PRC-33-023",1),required_load_capacity_profile())

# Clause-specific effective-environment and external-registry raw maps. Each row
# names its exact bounded subject domain and direct fields; expected maps are sealed
# independently from policy, inventory, or the authoritative counterparty record.
for _cid,_ordinal,_authority,_subject_suffix,_subject_examples,_fields,_inventory in [
    ("PRC-08-020",1,"environment","artifact_registry_ids",["registry-A","registry-B"],[("effective_authentication_method_identities","identity_map","Direct effective authentication method identity for each registry."),("effective_access_policy_digests","digest_map","Direct digest of the effective registry access policy bytes.")],"Complete in-scope artifact registry identities."),
    ("PRC-08-020",2,"environment","artifact_registry_ids",["registry-A","registry-B"],[("effective_authorization_policy_digests","digest_map","Direct digest of each effective registry authorization policy."),("effective_immutability_mode_identities","identity_map","Direct effective registry storage immutability mode identity.")],"Complete in-scope artifact registry identities."),
    ("PRC-08-020",3,"environment","artifact_registry_ids",["registry-A","registry-B"],[("effective_retention_policy_digests","digest_map","Direct digest of each effective registry retention policy."),("effective_minimum_retention_seconds","number_map","Direct numeric minimum retention duration configured for each registry.")],"Complete in-scope artifact registry identities."),
    ("PRC-09-016",1,"environment","supported_release_component_ids",["release-A|component-A","release-B|component-B"],[("vulnerability_monitor_registration_identities","identity_map","Direct monitoring-service registration identity for each supported released component."),("vulnerability_monitor_registration_states","state_map","Direct effective registration state returned for each supported released component.")],"Complete supported release x component identities."),
    ("PRC-15-042",1,"environment","old_version_ids",["version-old-A","version-old-B"],[("retirement_dates","time_map","Direct configured retirement date for every old version."),("version_monitor_identities","identity_map","Direct active version-monitor identity bound to every old version."),("version_monitor_states","state_map","Direct effective state of every bound old-version monitor.")],"Complete old-version identities from the supported/deployed version inventory."),
    ("PRC-17-015",1,"environment","temporary_permission_ids",["temp-permission-A","temp-permission-B"],[("configured_expiry_times","time_map","Direct configured expiry timestamp for every temporary permission."),("automatic_expiry_controller_identities","identity_map","Direct automatic expiry controller identity bound to every temporary permission."),("post_expiry_permission_states","state_map","Direct effective permission state observed by the independently sealed after-expiry check for every temporary permission.")],"Complete temporary permission identities plus independently scheduled after-expiry observations."),
    ("PRC-20-021",1,"environment","abandoned_failed_or_orphaned_upload_ids",["abandoned-upload-A","failed-upload-A","orphaned-upload-A"],[("cleanup_operation_identities","identity_map","Direct cleanup operation identity for every abandoned, failed, or orphaned upload."),("post_cleanup_storage_states","state_map","Direct storage state after cleanup for every upload."),("cleanup_completion_times","time_map","Direct cleanup completion timestamp for every upload.")],"Complete abandoned, failed, and orphaned upload identities from the upload lifecycle event inventory."),
    ("PRC-21-003",1,"environment","certificate_ids",["certificate-A","certificate-B"],[("certificate_byte_digests","digest_map","Locally recomputed digest of each exact certificate byte sequence."),("covered_name_set_digests","digest_map","Canonical digest of exact names covered by each certificate."),("validated_trust_chain_identities","identity_map","Direct chain identity selected by the pinned trust validator for each certificate."),("trust_validation_states","state_map","Direct terminal output from the pinned certificate name and trust-chain validator.")],"Complete certificate identities sealed from the effective endpoint inventory."),
    ("PRC-21-004",1,"environment","certificate_lifecycle_event_kind_ids",["issuance","renewal","expiration","revocation"],[("monitor_identities_by_event_kind","identity_map","Direct monitor identity bound to each certificate lifecycle event kind."),("monitor_states_by_event_kind","state_map","Direct effective state of each lifecycle monitor."),("monitor_scope_digests_by_event_kind","digest_map","Direct digest of certificate inventory and event scope covered by each monitor.")],"Exactly issuance, renewal, expiration, and revocation event kinds."),
    ("PRC-28-007",1,"environment","patch_subject_ids",["runtime-A","image-A","host-A"],[("installed_patch_version_identities","identity_map","Direct installed security patch version for every subject."),("risk_policy_required_patch_identities","identity_map","Direct required patch version resolved from the independently selected risk policy for every subject."),("patch_observation_times","time_map","Direct authenticated observation time of installed patch state.")],"Complete patch-governed runtime, image, and host identities."),
    ("PRC-28-012",1,"environment","network_path_ids",["frontend-to-api","api-to-db"],[("effective_segmentation_policy_digests","digest_map","Direct digest of effective segmentation policy for every network path."),("effective_private_connectivity_identities","identity_map","Direct effective private-connectivity binding identity for every path."),("reviewed_network_design_digests","digest_map","Direct independently approved network-design binding digest for every path.")],"Complete in-scope network path identities."),
    ("PRC-29-006",1,"environment","telemetry_stream_ids",["logs-A","metrics-A","traces-A"],[("observed_release_version_identities","identity_map","Direct release-version identity observed in each telemetry stream."),("deployed_release_version_identities","identity_map","Direct independently derived deployed release-version identity for each telemetry stream scope.")],"Complete required telemetry stream identities."),
    ("PRC-30-008",1,"environment","required_service_hour_slot_ids",["monday-09","monday-10"],[("primary_coverage_identity_by_slot","identity_map","Direct primary on-call coverage identity for every required service-hour slot."),("secondary_coverage_identity_by_slot","identity_map","Direct secondary on-call coverage identity for every required service-hour slot.")],"Complete required service-hour slot identities sealed from the support schedule."),
    ("PRC-31-027",2,"environment","vulnerability_reporting_channel_ids",["security-form","security-contact-page"],[("resolved_endpoint_identities","identity_map","Direct resolved live endpoint identity for every reporting channel."),("nonmutating_probe_response_digests","digest_map","Direct digest of the bounded non-mutating availability response."),("probe_terminal_states","state_map","Direct terminal state of each bounded availability probe.")],"Complete vulnerability-reporting channel identities."),
    ("PRC-35-018",2,"environment","recipient_view_document_ids",["team-A|view-A|doc-v1","team-B|view-B|doc-v1"],[("delivery_or_access_event_identities","identity_map","Direct authenticated delivery or access event for each expected recipient and manifested view/document version."),("delivered_view_or_document_digests","digest_map","Direct digest of the exact manifested telemetry view or document delivered or accessed."),("delivery_or_access_times","time_map","Direct event timestamp for delivery or access.")],"Complete expected recipient x manifested view or document-version identities."),
    ("USEQ-02E3B3DC",2,"environment","obsolete_runbook_access_route_ids",["route-old-A","route-old-B"],[("resolved_effective_successor_procedure_identities","identity_map","Direct successor procedure identity resolved by each obsolete-runbook access route."),("effective_route_configuration_digests","digest_map","Direct digest of the effective redirect or alias configuration for each route.")],"Complete declared access-route identities for obsolete runbooks."),
    ("USEQ-0701D270",1,"environment","publication_topic_ids",["service-status","maintenance","accessibility","privacy","security-contact","support"],[("published_content_artifact_digests","digest_map","Direct digest of the current published artifact for each required information topic."),("publication_endpoint_identities","identity_map","Direct declared publication endpoint for each topic."),("publication_states","state_map","Direct effective publication state at each endpoint.")],"Exactly service status, maintenance, accessibility, privacy, security contact, and support topics."),
    ("USEQ-0A48A6ED",1,"external_registry","standard_ids",["standard-A","standard-B"],[("issuing_body_current_baseline_identities","identity_map","Direct current published baseline identity from each issuing body."),("referenced_draft_status_identities","identity_map","Direct current issuing-body draft status for every referenced standard."),("issuing_body_record_digests","digest_map","Direct digest of authenticated issuing-body metadata.")],"Complete referenced standard identities."),
    ("USEQ-16C2ACD8",1,"external_registry","referenced_document_ids",["document-A","document-B"],[("issuing_body_status_identities","identity_map","Direct current status identity published by the issuing body for each referenced document."),("tracked_status_identities","identity_map","Direct status identity recorded by the project for each document."),("issuing_body_metadata_digests","digest_map","Direct digest of current authenticated issuing-body metadata.")],"Complete referenced document identities."),
    ("USEQ-18804590",2,"external_registry","reusable_asset_ids",["asset-A","asset-B"],[("published_policy_artifact_digests","digest_map","Direct digest of the currently published compatibility, deprecation, and migration policy for each asset."),("published_asset_and_version_binding_identities","identity_map","Direct asset and supported-version binding parsed from each published policy."),("published_policy_version_identities","identity_map","Direct immutable policy version identity."),("published_policy_times","time_map","Direct authenticated policy publication time."),("public_retrieval_endpoint_identities","identity_map","Direct public retrieval endpoint identity for each policy.")],"Complete reusable asset identities."),
    ("USEQ-2916316A",1,"environment","event_workflow_ids",["workflow-A","workflow-B"],[("dead_letter_handler_identities","identity_map","Direct enabled dead-letter handler identity for each workflow."),("failure_destination_identities","identity_map","Direct enabled failure destination identity."),("replay_mechanism_identities","identity_map","Direct enabled replay mechanism identity."),("reconciliation_mechanism_identities","identity_map","Direct enabled reconciliation mechanism identity."),("mechanism_enabled_state_identities","state_map","Direct effective enabled state for the workflow-bound mechanism set.")],"Complete sealed event workflow identities."),
    ("USEQ-2B757A11",1,"environment","closed_defect_ids",["defect-A","defect-B"],[("completion_criterion_set_digests","digest_map","Direct digest of every satisfied independently versioned completion-criterion set for each closed defect."),("required_evidence_object_set_digests","digest_map","Direct digest of the complete evidence-object bindings for each defect."),("criteria_satisfied_times","time_map","Direct time every completion criterion was satisfied."),("evidence_bound_times","time_map","Direct time all required evidence was bound."),("defect_close_times","time_map","Direct authenticated close transition time.")],"Complete defect identities transitioned to closed."),
    ("USEQ-2C522C55",1,"environment","emergency_action_ids",["revocation","suspension","reissue","algorithm-disablement","customer-communication"],[("accountable_authority_identities","identity_map","Direct accountable authority identity assigned to each emergency action."),("protected_contact_route_identities","identity_map","Direct protected current contact route identity for each action."),("route_verification_times","time_map","Direct authenticated last verification time for each route."),("route_state_identities","state_map","Direct current route availability state.")],"Exactly revocation, suspension, reissue, algorithm disablement, and customer communication action identities."),
    ("USEQ-49154A5B",1,"environment","analyzer_rule_ids",["rule-A","rule-B"],[("effective_severity_threshold_identities","identity_map","Direct effective severity threshold identity explicitly bound to every analyzer rule."),("effective_configuration_source_states","state_map","Direct raw source state showing each effective value is explicit rather than an inherited unreviewed default."),("effective_analyzer_configuration_digests","digest_map","Direct digest of effective configuration bytes containing each binding.")],"Complete sealed analyzer-rule identities."),
    ("USEQ-4B7B4908",1,"environment","supported_interoperability_pair_ids",["legacy-legacy","legacy-pq-hybrid","pq-pq"],[("available_test_environment_identities","identity_map","Direct available test-environment identity for each supported interoperability pair."),("participant_identity_set_digests","digest_map","Canonical digest of explicit participant identities for each pair."),("participant_version_set_digests","digest_map","Canonical digest of participant versions."),("algorithm_set_digests","digest_map","Canonical digest of algorithms."),("trust_material_set_digests","digest_map","Canonical digest of trust material."),("network_configuration_digests","digest_map","Direct digest of network configuration.")],"Complete independently supported legacy-only, hybrid, and post-quantum participant-pair identities."),
    ("USEQ-542BC36E",2,"environment","active_compensating_control_ids",["control-A","control-B"],[("bound_control_identities","identity_map","Direct effective control identity for each active compensating control record."),("bound_monitor_identities","identity_map","Direct effective monitor identity bound to each control."),("control_enabled_states","state_map","Direct effective enabled state of each control."),("monitor_enabled_states","state_map","Direct effective enabled state of each bound monitor."),("assessed_scope_identities","identity_map","Direct scope identity shared by the bound control and monitor.")],"Complete active compensating-control identities."),
    ("USEQ-5CA932E3",2,"environment","nonretired_article_ids",["article-A","article-B"],[("bounded_search_result_article_identities","identity_map","Direct article identity returned by the bounded search for each required article."),("search_query_input_digests","digest_map","Direct digest of exact stable identifier and independently sealed lookup terms used for each query."),("search_result_states","state_map","Direct terminal state of every bounded search query.")],"Complete non-retired article identities."),
    ("USEQ-68821BC1",1,"environment","deployed_ai_service_ids",["ai-service-A","ai-service-B"],[("deployed_model_content_digests","digest_map","Direct deployed model content digest."),("deployed_prompt_content_digests","digest_map","Direct deployed prompt content digest."),("deployed_retrieval_configuration_digests","digest_map","Direct retrieval configuration digest."),("deployed_routing_configuration_digests","digest_map","Direct routing configuration digest."),("deployed_tool_configuration_digests","digest_map","Direct tool configuration digest."),("deployed_policy_configuration_digests","digest_map","Direct policy configuration digest.")],"Complete deployed AI service identities."),
    ("USEQ-692CA8B7",1,"external_registry","tracked_standard_ids",["standard-A","standard-B"],[("tracked_status_identities","identity_map","Direct project-tracked status identity for each standard."),("issuing_body_current_status_identities","identity_map","Direct current status identity from each issuing body."),("issuing_body_record_digests","digest_map","Direct digest of authenticated issuing-body status metadata.")],"Complete tracked standard identities."),
    ("USEQ-84F487EA",1,"environment","publication_surface_ids",["developer-portal","service-catalog"],[("published_service_boundary_digests","digest_map","Direct digest of current approved service-boundary content."),("published_supported_use_case_set_digests","digest_map","Direct digest of supported use cases."),("published_non_goal_set_digests","digest_map","Direct digest of non-goals."),("published_dependency_set_digests","digest_map","Direct digest of dependencies."),("published_escalation_model_digests","digest_map","Direct digest of escalation model."),("declared_publication_endpoint_identities","identity_map","Direct endpoint identity for each surface.")],"Complete platform-service audience publication-surface identities."),
    ("USEQ-A0E78A59",1,"environment","billing_scope_resource_class_ids",["billing-A|compute","billing-A|storage"],[("bound_cost_alert_identities","identity_map","Direct cost-alert identity bound to each billing scope and resource class."),("bound_runaway_protection_identities","identity_map","Direct runaway-resource protection identity bound to each exact scope and class."),("cost_alert_enabled_states","state_map","Direct effective cost-alert state."),("runaway_protection_enabled_states","state_map","Direct effective runaway-protection state.")],"Complete billing scope x resource-class identities."),
    ("USEQ-BBED00FD",1,"environment","endpoint_ids",["endpoint-A","endpoint-B"],[("effective_trusted_origin_set_digests","digest_map","Canonical digest of exact trusted origins configured for every endpoint."),("credentialed_access_states","state_map","Direct effective credentialed cross-origin access state."),("wildcard_origin_states","state_map","Direct effective wildcard-origin state."),("reviewed_origin_policy_digests","digest_map","Direct independently reviewed origin-policy digest applied to each endpoint.")],"Complete in-scope endpoint identities."),
    ("USEQ-C731071C",1,"environment","public_route_kind_ids",["correction","takedown","appeal","incident-reporting"],[("resolved_live_endpoint_identities","identity_map","Direct resolved endpoint identity for each route kind."),("nonmutating_probe_response_digests","digest_map","Direct digest of expected non-mutating availability response."),("availability_probe_states","state_map","Direct terminal state of each route probe.")],"Exactly correction, takedown, appeal, and incident-reporting route identities."),
    ("USEQ-D8F0B0C5",1,"environment","publication_surface_ids",["developer-portal","service-catalog"],[("published_maintenance_window_digests","digest_map","Direct digest of current approved maintenance-window content."),("published_deprecation_schedule_digests","digest_map","Direct digest of deprecation schedules."),("published_breaking_change_digests","digest_map","Direct digest of breaking-change content."),("published_reliability_limitation_digests","digest_map","Direct digest of reliability limitations."),("declared_publication_endpoint_identities","identity_map","Direct endpoint identity for each publication surface.")],"Complete platform-service audience publication-surface identities."),
    ("USEQ-EEB12085",2,"environment","internal_publication_surface_ids",["internal-portal-A","internal-catalog-A"],[("exposed_evaluation_result_artifact_digests","digest_map","Direct digest of exact bound evaluation-result artifact exposed by each surface."),("exposed_limitation_set_digests","digest_map","Direct digest of the exact limitation set exposed."),("exposed_approved_use_boundary_set_digests","digest_map","Direct digest of the approved-use boundary set exposed."),("sealed_internal_audience_identities","identity_map","Direct audience identity authorized for each surface.")],"Complete declared internal publication-surface identities."),
]:
    add((_cid,_ordinal),expected_field_maps_profile(_cid,_ordinal,_authority,_subject_suffix,_fields,_inventory,_subject_examples))

add(("USEQ-BBED00FD",1),cors_invariant_profile())
add(("USEQ-A0E78A59",1),cost_protection_enabled_profile())


def generic_gap(family, statement):
    if family == "execution_evidence": return "execution_case_matrix_join"
    if family == "artifact_integrity": return "artifact_graph_or_cryptographic_relation"
    if family == "ci_policy": return "ci_rule_event_relational_join"
    if family in {"environment_evidence","container_iac"}: return "keyed_resource_state_or_event_join"
    if family == "analysis_adapter": return "analysis_finding_or_run_relation"
    if family in {"source_ast","structured_document"}: return "bounded_static_structure_query"
    if family == "package_metadata": return "package_graph_and_string_digest_map"
    return "typed_relational_record_join"


def predicate_ops(node):
    result=set()
    if isinstance(node,dict):
        if isinstance(node.get("op"),str): result.add(node["op"])
        for value in node.values(): result.update(predicate_ops(value))
    elif isinstance(node,list):
        for value in node: result.update(predicate_ops(value))
    return result


def main():
    doc=json.loads(BINDINGS.read_text());rows=[(b,c) for b in doc["bindings"] for c in b["clauses"] if c["checker_family"]!="structured_record"]
    if not rows: raise SystemExit("no current non-structured clauses")
    definitions=[]; missing=Counter(); missing_controls=defaultdict(set); wrapper_counts=Counter(); reclass=[]
    for binding,clause in rows:
        cid=binding["control_id"];ordinal=clause["ordinal"];key=(cid,ordinal);statement=clause["statement"]
        base={"control_id":cid,"control_revision":binding["revision"],"control_semantic_sha256":binding["semantic_sha256"],"clause_ordinal":ordinal,"clause_id":clause["clause_id"],"clause_statement":statement,"checker_family":clause["checker_family"],"required_authority":clause["evidence_authority"],"corrected_checker_family":clause["checker_family"],"corrected_required_authority":clause["evidence_authority"]}
        wrapper=None
        for label,start in GENERIC_PREFIXES.items():
            if statement.startswith(start): wrapper=label;break
        if cid in RECLASSIFY or wrapper:
            if cid in RECLASSIFY:
                status="classification_error";reason=RECLASSIFY[cid];gaps=[];reclass.append({"control_id":cid,"clause_ordinal":ordinal,"clause_id":clause["clause_id"],"reason":reason})
            else:
                status="classification_error";reason=f"The current clause is a semantic wrapper ({wrapper}) that delegates the original control outcome to an unspecified provider or acceptance contract; no closed predicate can decide it.";gaps=[];wrapper_counts[wrapper]+=1;reclass.append({"control_id":cid,"clause_ordinal":ordinal,"clause_id":clause["clause_id"],"reason":reason})
            continue
        elif key in EXACT:
            profile=copy.deepcopy(EXACT[key])
            authoritative_source={
                "repository":"the exact assessed repository bytes plus the authenticated parser input manifest",
                "artifact":"the exact immutable artifact, manifest, or retained object bytes identified by the clause",
                "environment":"the read-only effective control-plane configuration, state, or authenticated event record identified by the clause",
                "executed":"the authenticated invocation input manifest and direct raw output or event record for the exact bounded execution",
                "external_registry":"the authenticated current registry, publisher, provider, or issuing-body record identified by the clause",
                "structured_record":"the exact authenticated structured record bytes identified by the clause",
            }[clause["evidence_authority"]]
            sealed_names=[item["parameter_id"].rsplit(".",1)[-1] for item in profile["params"]]
            for contract in profile["facts"]:
                contract["source_requirement"]=(f"For {cid} clause {ordinal}, collect from {authoritative_source}: {contract['raw_value_semantics']} Preserve exact subject keys and bind them to {', '.join(sealed_names) if sealed_names else 'the singular exact clause subject'} before evaluation.")
            fact_labels=[item["fact_id"].rsplit(".",1)[-1] for item in profile["facts"]]
            parameter_labels=[item["parameter_id"].rsplit(".",1)[-1] for item in profile["params"]]
            pass_facts=profile["fixtures"]["pass"]["facts"]
            fail_facts=profile["fixtures"]["fail"]["facts"]
            changed=[name.rsplit(".",1)[-1] for name in pass_facts if pass_facts.get(name)!=fail_facts.get(name)]
            unchanged=[name.rsplit(".",1)[-1] for name in pass_facts if pass_facts.get(name)==fail_facts.get(name)]
            changed_text=", ".join(changed) if changed else "a declared raw fact"
            operations=sorted(predicate_ops(profile["predicate"]))
            changed_semantics=[item["raw_value_semantics"] for item in profile["facts"] if item["fact_id"].rsplit(".",1)[-1] in changed]
            unchanged_semantics=[item["raw_value_semantics"] for item in profile["facts"] if item["fact_id"].rsplit(".",1)[-1] in unchanged]
            original_pass=profile["fixtures"]["pass"]["description"]
            original_fail=profile["fixtures"]["fail"]["description"]
            original_counter=profile["fixtures"]["counterexample"]["description"]
            proof_detail="; ".join(item["raw_value_semantics"] for item in profile["facts"][:4])
            sealed_detail="; ".join(item["source_requirement"] for item in profile["params"][:3]) if profile["params"] else "the two observations are compared directly without a provider-selected expected value"
            review_reason=(
                f"For {cid} clause {ordinal}, the executable decision for '{statement}' is: "
                f"{original_pass} It evaluates the raw fields {', '.join(fact_labels)} with "
                f"the exact operators {', '.join(operations)}. Scope and expected values are "
                f"compiled from {sealed_detail}; evidence collection cannot change them."
            )
            if original_counter.startswith("A tempting partial signal") or original_counter.startswith("All generic execution states"):
                counterexample=(
                    f"For {cid} clause {ordinal}, keep {unchanged_semantics[0] if unchanged_semantics else 'all other raw observations'} "
                    f"valid but break {changed_semantics[0] if changed_semantics else changed_text}. "
                    f"The exact program returns Fail because {original_fail.lower()}"
                )
            else:
                counterexample=original_counter
            first_fact=fact_labels[0]
            profile["fixtures"]["pass"]["description"]=f"{original_pass} Applied to {cid} clause {ordinal}: {statement}"
            profile["fixtures"]["fail"]["description"]=f"{original_fail} Broken observation: {changed_semantics[0] if changed_semantics else changed_text}."
            profile["fixtures"]["blocked"]["description"]=f"Evaluation stops for {cid} clause {ordinal} because the source for {profile['facts'][0]['raw_value_semantics']} is incomplete."
            profile["fixtures"]["counterexample"]["description"]=counterexample
            base.update({"classification_status":"exact_predicate","classification_error_reason":None,"raw_fact_contracts":profile["facts"],"sealed_parameter_contracts":profile["params"],"predicate":profile["predicate"],"required_runtime_ops":operations,"collector_contract":{"collector_id":f"prc.collect.{prefix(cid,ordinal).replace('_','.')}@0.1","required_sources":list(dict.fromkeys(x["source_requirement"] for x in profile["facts"])),"inventory_contract":f"For {cid} clause {ordinal}, every map or set key must bind to the independently sealed domain named by: {', '.join(parameter_labels) if parameter_labels else 'the exact scalar subject in the clause'}; provider-selected nonempty output is not completeness.","normalization_contract":f"For {cid} clause {ordinal}, emit only direct typed values for {', '.join(fact_labels)}; do not emit a compliance result, score, or semantic shortcut.","completeness_contract":f"For {cid} clause {ordinal}, any omitted key or value, parse error, permission gap, unsupported subject, or source conflict affecting {', '.join(fact_labels)} sets complete=false.","freshness_contract":f"For {cid} clause {ordinal}, bind every observation for {', '.join(fact_labels)} to the program maximum evidence age and the exact assessed revision.","provider_status":"unregistered"},"fixtures":profile["fixtures"],"review_reason":review_reason,"counterexample_analysis":counterexample})
        else:
            status="missing_reusable_ops";gap=generic_gap(clause["checker_family"],statement);gaps=[gap];missing[gap]+=1;missing_controls[gap].add(cid);reason=f"This deterministic clause requires {gap.replace('_',' ')} over raw heterogeneous records; the current closed DSL has no honest representation for the full statement. A scalar or provider verdict would allow a pass while part of the clause is broken."
            continue
        definitions.append(base)
    output={"schema_version":"prc.control-check-program-definitions/v0.1","scope":"non_structured","source_binding_catalog_sha256":hashlib.sha256(BINDINGS.read_bytes()).hexdigest(),"definition_count":len(definitions),"definitions":definitions}
    OUT.write_text(json.dumps(output,ensure_ascii=False,indent=2)+"\n")
    print(json.dumps({"catalog_non_structured_clause_count":len(rows),"exact_definition_count":len(definitions),"not_emitted_pending_strength_or_runtime_support":len(rows)-len(definitions),"semantic_wrapper_count":sum(wrapper_counts.values()),"reclassification_candidate_count":len(reclass),"missing_reusable_op_summary":dict(sorted(missing.items()))},indent=2))

if __name__=="__main__": main()
