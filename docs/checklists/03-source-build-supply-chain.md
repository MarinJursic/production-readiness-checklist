# Source, build, and supply chain

> Protect source changes, build integrity, dependencies, provenance, and licensing.

Sections 7–9 of the master checklist. For each applicable item, capture status, owner, evidence, release, environment, reviewer, and evidence age.

## 7. Source control and change management

- [ ] **PRC-07-001** — One authoritative source repository is identified.
- [ ] **PRC-07-002** — Repository access uses individual identities rather than shared accounts.
- [ ] **PRC-07-003** — Strong authentication is required for repository access.
- [ ] **PRC-07-004** — Access follows least privilege.
- [ ] **PRC-07-005** — Protected branches prevent unreviewed direct changes.
- [ ] **PRC-07-006** — Material changes require independent review.
- [ ] **PRC-07-007** — Critical code areas require designated reviewers.
- [ ] **PRC-07-008** — Review cannot be bypassed silently by administrators.
- [ ] **PRC-07-009** — Merge requirements include the required automated checks.
- [ ] **PRC-07-010** — Security-sensitive changes receive appropriate specialist review.
- [ ] **PRC-07-011** — Infrastructure, pipeline, authorization, cryptography, payment, and migration changes receive heightened review.
- [ ] **PRC-07-012** — Changes are traceable from requirement or issue to code, review, build, artifact, and deployment.
- [ ] **PRC-07-013** — Commit history and release tags are protected against unauthorized rewriting.
- [ ] **PRC-07-014** — Secret scanning covers current changes and repository history.
- [ ] **PRC-07-015** — Large binaries, generated files, and vendored code are controlled.
- [ ] **PRC-07-016** — Generated code has a documented trusted source and regeneration process.
- [ ] **PRC-07-017** — Emergency changes follow a documented break-glass process and receive retrospective review.
- [ ] **PRC-07-018** — Repository backup and recovery have been tested.
- [ ] **PRC-07-019** — Departed personnel and obsolete automation identities no longer have access.
- [ ] **PRC-07-020** — Repository webhooks, apps, bots, and deploy keys are inventoried and restricted.

---

## 8. Build, CI/CD, and artifact integrity

SLSA 1.2 organizes software supply-chain assurance around verifiable provenance and build integrity. The checklist below applies those outcomes without requiring a particular build platform. ([slsa.dev](https://slsa.dev/spec/v1.2/))

- [ ] **PRC-08-001** — Builds run from reviewed source.
- [ ] **PRC-08-002** — Builds are automated and repeatable.
- [ ] **PRC-08-003** — A clean build produces equivalent output from the same declared inputs.
- [ ] **PRC-08-004** — Build inputs and dependency versions are pinned sufficiently to prevent unexpected substitution.
- [ ] **PRC-08-005** — Builds fail closed when required verification cannot run.
- [ ] **PRC-08-006** — The production artifact is built once and promoted, not independently rebuilt in each environment.
- [ ] **PRC-08-007** — Deployable artifacts are immutable.
- [ ] **PRC-08-008** — Every artifact has a cryptographic digest.
- [ ] **PRC-08-009** — Build provenance records source, builder, inputs, commands or workflow, time, and output digest.
- [ ] **PRC-08-010** — Artifact signatures or attestations are verified before deployment.
- [ ] **PRC-08-011** — Build identities and signing keys use least privilege.
- [ ] **PRC-08-012** — Build and signing keys are stored outside source code and build output.
- [ ] **PRC-08-013** — Build runners are patched and isolated appropriately.
- [ ] **PRC-08-014** — Untrusted contribution workflows cannot access production secrets.
- [ ] **PRC-08-015** — Pull-request builds cannot overwrite trusted release artifacts.
- [ ] **PRC-08-016** — CI jobs have minimum required network, repository, cloud, and registry permissions.
- [ ] **PRC-08-017** — Pipeline definitions receive the same review protection as application code.
- [ ] **PRC-08-018** — Required unit, integration, security, policy, and packaging checks cannot be skipped silently.
- [ ] **PRC-08-019** — Failed or canceled checks cannot be misrepresented as passed.
- [ ] **PRC-08-020** — Artifact registries enforce authentication, authorization, immutability, and retention.
- [ ] **PRC-08-021** — The previous known-good artifact remains available.
- [ ] **PRC-08-022** — Build logs do not expose secrets or sensitive source material.
- [ ] **PRC-08-023** — Build caches cannot substitute unverified or cross-trust artifacts.
- [ ] **PRC-08-024** — Release tags, artifacts, provenance, and deployment records agree.
- [ ] **PRC-08-025** — Compromise of the build platform has a documented containment and recovery procedure.
- [ ] **PRC-08-026** — A supply-chain incident drill has tested revocation, rebuild, resigning, and redeployment.

---

## 9. Dependencies, SBOM, and licenses

A software bill of materials should cover the complete release, while vulnerability prioritization should account for known exploitation and deployment context rather than merely counting CVEs. ([cisa.gov](https://www.cisa.gov/topics/information-communications-technology-supply-chain-security/sbom))

- [ ] **PRC-09-001** — Generate an SBOM for every released application artifact.
- [ ] **PRC-09-002** — Include direct and transitive components.
- [ ] **PRC-09-003** — Include runtime, operating-system, container, client, server, plugin, model, and embedded components as applicable.
- [ ] **PRC-09-004** — Distinguish build-only, development-only, optional, and runtime dependencies.
- [ ] **PRC-09-005** — Record exact versions, package identifiers, sources, hashes, and relationships.
- [ ] **PRC-09-006** — Record license information and required notices.
- [ ] **PRC-09-007** — Link the SBOM to the released artifact digest.
- [ ] **PRC-09-008** — Verify the SBOM against the final packaged artifact.
- [ ] **PRC-09-009** — Ensure no unexpected undeclared components are present.
- [ ] **PRC-09-010** — Remove unused dependencies.
- [ ] **PRC-09-011** — Remove unsupported and end-of-life dependencies.
- [ ] **PRC-09-012** — Check whether dependencies come from intended publishers and registries.
- [ ] **PRC-09-013** — Protect against dependency confusion, namespace collision, and typographical substitution.
- [ ] **PRC-09-014** — Assess maintainer health and abandonment risk for critical dependencies.
- [ ] **PRC-09-015** — Monitor upstream security advisories.
- [ ] **PRC-09-016** — Continuously scan released components for newly disclosed vulnerabilities.
- [ ] **PRC-09-017** — Prioritize vulnerabilities using exploit activity, reachability, exposure, asset criticality, data sensitivity, and compensating controls.
- [ ] **PRC-09-018** — Check relevant known-exploited-vulnerability catalogs.
- [ ] **PRC-09-019** — Define remediation targets by actual risk.
- [ ] **PRC-09-020** — Verify fixes or mitigations after deployment.
- [ ] **PRC-09-021** — Maintain a process for emergency replacement of a compromised dependency.
- [ ] **PRC-09-022** — Document ownership of vendored, forked, or internally patched components.
- [ ] **PRC-09-023** — Test updates before promotion.
- [ ] **PRC-09-024** — Verify license compatibility among combined components.
- [ ] **PRC-09-025** — Provide required attributions, source offers, notices, and redistribution terms.
- [ ] **PRC-09-026** — Ensure proprietary or confidential code is not unintentionally incorporated into distributable artifacts.
- [ ] **PRC-09-027** — Retain SBOMs for historical releases and incident investigation.

---
