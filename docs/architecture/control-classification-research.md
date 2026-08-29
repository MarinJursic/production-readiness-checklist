# Control-classification research basis

The binary classification contract is intentionally stricter than “a script can
find something related to this sentence.” It follows established assessment and
policy-as-code boundaries:

- [NIST SP 800-53A Rev. 5](https://csrc.nist.gov/pubs/sp/800/53/a/r5/final)
  distinguishes examining evidence, interviewing people, and testing behavior,
  and requires assessment procedures to be tailored to the system and its risk
  tolerance. An interview or accountable judgment is therefore not converted
  into a deterministic result merely because software can ask the question.
- [NIST OSCAL assessment results](https://pages.nist.gov/OSCAL/learn/concepts/layer/assessment/assessment-results/)
  bind machine-readable findings to an assessment plan, assessed subjects,
  activities, observations, evidence, and time. The scanner mirrors those
  boundaries with explicit subject, inventory, authority, and freshness fields.
- [Open Policy Agent's policy language](https://www.openpolicyagent.org/docs/policy-language)
  evaluates declarative rules over structured input. The deterministic runtime
  similarly uses a closed typed predicate language rather than executing text
  taken from a checklist or target repository.
- [SLSA artifact verification](https://slsa.dev/spec/v1.1/verifying-artifacts)
  requires an artifact to match its provenance, signature, trusted builder, and
  expected fields. This supports the rule that a digest, attestation, or provider
  conclusion is not enough unless it is bound to the exact artifact and approved
  expectations.
- [OpenSSF Scorecard](https://github.com/ossf/scorecard) describes its automated
  checks as security-health heuristics. This scanner keeps useful partial signals
  separate from proof of a broader control so that a heuristic cannot silently
  become an authoritative Pass.

The authoritative classification rules remain in
[`control-classification.md`](control-classification.md). Keeping the research
notes separate lets the methodology artifact retain its stable review digest.
