# AI-assisted review

AI agents can accelerate a production-readiness review, especially when controls can be evaluated from source, tests, infrastructure definitions, configuration schemas, and repository documentation. They should be used as evidence finders and gap detectors—not as autonomous release authorities.

## What an agent can often inspect

- build, test, lint, and deployment configuration;
- dependency manifests, lockfiles, and update policy;
- authentication, authorization, input handling, and security-sensitive code paths;
- database migrations and backward-compatibility patterns;
- automated tests and their relationship to critical journeys;
- infrastructure-as-code, observability definitions, and runbooks;
- documented architecture, recovery, incident, and support procedures;
- the current diff and the controls it could invalidate.

## What usually needs human or external evidence

- the exact production configuration and deployed artifact;
- successful backup restoration, failover, capacity, or incident drills;
- alert delivery and on-call authority in real operating conditions;
- legal, regulatory, privacy, accessibility, and contractual applicability;
- third-party contracts, organizational processes, and staff availability;
- penetration-test findings or production telemetry not present in the repository;
- the authority to accept residual risk and approve a launch.

An agent must label these controls **Blocked** or **Unknown**, not infer a pass from missing information.

## Recommended workflow

1. Put `CLAUDE.md` at the root of the application repository, or copy its instructions into the agent context.
2. Choose [a full lifecycle review](../prompts/full-readiness-review.md), [a release-diff review](../prompts/release-diff-review.md), or [an evidence challenge](../prompts/evidence-challenge.md).
3. Give the agent the release identifier, intended environment, architecture overview, critical journeys, and applicable conditional modules.
4. Allow read-only inspection and relevant test commands. Require permission before code or configuration changes.
5. Require a citation for each Pass: file and line, test output, command result, artifact, or external evidence reference.
6. Have a qualified person review high-impact conclusions and every proposed exception.

## Useful output contract

Ask the agent to return:

1. release and scope understood;
2. immediate no-go findings;
3. a control table with `USEQ-*` or `PRC-*` ID, status, evidence, confidence, owner, and next action;
4. unknowns that require human or external evidence;
5. changed controls invalidated by the release diff;
6. prioritized blockers and recommended verification commands;
7. an explicit statement that the agent is not making the final release decision.

## Anti-patterns

- Asking “is this production ready?” without release scope or evidence criteria.
- Treating a repository-only release review as a complete lifecycle assessment when governance, product, operational, or external evidence was never supplied.
- Treating the absence of a visible vulnerability as proof of security.
- Allowing the agent to mark organizational or production-only controls as passed from source code.
- Counting passed boxes or reporting a readiness percentage without blocker analysis.
- Letting an agent accept risk on behalf of a named human owner.
- Running destructive, load, failover, migration, or production tests without explicit authorization and safe boundaries.

## Review the reviewer

Challenge the agent’s findings. Check that cited evidence supports the exact control, pertains to the exact release, comes from the intended environment, is current enough, and demonstrates behavior rather than intent. The [evidence challenge prompt](../prompts/evidence-challenge.md) is designed for a second-pass review.
