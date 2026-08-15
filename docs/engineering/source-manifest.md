# Source consolidation manifest

This manifest records how the two supplied archives were incorporated. Archive content was treated as source material; operational instructions inside it were not executed.

## Consolidation result

- Source documents reviewed: **197**
- Source checkbox lines reviewed: **14,132**
- New non-duplicative lifecycle controls imported: **4,917**
- Existing production-readiness controls retained: **1,421**
- Final unique control set: **6,338**
- Published lifecycle manuals: **16**, followed by the existing production-readiness review

## Source archives

- `universal-software-engineering-quality-standards.zip` — SHA-256 `cde6ff7a3c89cbd4bd10046a1a4a8232a84dc1e25ced8d73e0a69d4d6d5b9bda`
- `universal-software-engineering-gap-supplement.zip` — SHA-256 `5c4bd9cae92a54013bfa6f69c3be2656dd74c974ee043ae147c0573db3e4b16d`

## Consolidation rules

- The quality archive's production master was not copied because this repository already maintains that review with stable `PRC-` identifiers.
- The gap supplement's `Consolidated controls from the prior corpus` sections were not copied a second time. Only expanded gap-closure controls, category evidence, and category no-go controls were candidates for import.
- Repeated applicability boilerplate and normalized exact duplicates were retained once at their earliest lifecycle phase.
- Exact matches to existing `PRC-` controls were not reissued with new identifiers.
- Imported controls use deterministic `USEQ-` identifiers derived from normalized control text so evidence references remain stable across regeneration.
- Similar-looking controls were retained when their wording expressed a materially different scope, condition, or evidence obligation.

## Source inventory

| Package | Source document | Source checks | Candidates | Imported | Covered or excluded | Destination or treatment |
| --- | --- | ---: | ---: | ---: | ---: | --- |
| quality standards | `00-core/01-applicability-evidence-and-exceptions.md` | 61 | 61 | 61 | 0 | Phase 1: Governance and foundations |
| quality standards | `00-core/02-universal-quality-attribute-model.md` | 70 | 70 | 62 | 8 | Phase 1: Governance and foundations |
| quality standards | `00-core/03-standards-and-frameworks-crosswalk.md` | 16 | 16 | 8 | 8 | Phase 1: Governance and foundations |
| quality standards | `00-core/04-universal-definition-of-done.md` | 43 | 43 | 35 | 8 | Phase 1: Governance and foundations |
| quality standards | `00-core/05-master-no-go-gates.md` | 33 | 33 | 25 | 8 | Phase 1: Governance and foundations |
| quality standards | `00-core/06-final-evidence-package.md` | 30 | 30 | 21 | 9 | Phase 1: Governance and foundations |
| quality standards | `00-core/07-glossary.md` | 8 | 8 | 0 | 8 | Phase 1: Governance and foundations |
| quality standards | `00-core/08-corpus-maintenance-and-versioning.md` | 28 | 28 | 20 | 8 | Phase 1: Governance and foundations |
| quality standards | `00-core/09-assessment-record-template.md` | 8 | 8 | 0 | 8 | Phase 1: Governance and foundations |
| quality standards | `00-core/10-category-index.md` | 8 | 0 | 0 | 8 | navigation-only source index |
| quality standards | `01-governance/01-governance-ownership-and-risk.md` | 34 | 34 | 19 | 15 | Phase 1: Governance and foundations |
| quality standards | `01-governance/02-third-party-and-supplier-readiness.md` | 31 | 31 | 15 | 16 | Phase 1: Governance and foundations |
| quality standards | `01-governance/03-legal-regulatory-and-contractual-applicability.md` | 42 | 42 | 27 | 15 | Phase 1: Governance and foundations |
| quality standards | `01-governance/04-quality-management-and-continual-improvement.md` | 63 | 63 | 55 | 8 | Phase 1: Governance and foundations |
| quality standards | `01-governance/05-engineering-governance-and-decision-rights.md` | 58 | 58 | 15 | 43 | Phase 1: Governance and foundations |
| quality standards | `01-governance/06-project-program-and-portfolio-management.md` | 58 | 58 | 15 | 43 | Phase 1: Governance and foundations |
| quality standards | `01-governance/07-engineering-measurement-and-metrics.md` | 58 | 58 | 15 | 43 | Phase 1: Governance and foundations |
| quality standards | `01-governance/08-people-competence-culture-and-sustainable-work.md` | 58 | 58 | 15 | 43 | Phase 1: Governance and foundations |
| quality standards | `01-governance/09-ethics-and-responsible-technology.md` | 58 | 58 | 15 | 43 | Phase 1: Governance and foundations |
| quality standards | `02-product/01-product-and-business-readiness.md` | 27 | 27 | 14 | 13 | Phase 2: Product and requirements |
| quality standards | `02-product/02-product-strategy-vision-and-outcomes.md` | 57 | 57 | 14 | 43 | Phase 2: Product and requirements |
| quality standards | `02-product/03-product-discovery-and-problem-validation.md` | 58 | 58 | 15 | 43 | Phase 2: Product and requirements |
| quality standards | `02-product/04-requirements-engineering.md` | 60 | 60 | 17 | 43 | Phase 2: Product and requirements |
| quality standards | `02-product/05-prioritization-roadmaps-and-scope-control.md` | 58 | 58 | 15 | 43 | Phase 2: Product and requirements |
| quality standards | `02-product/06-product-metrics-experimentation-and-learning.md` | 58 | 58 | 15 | 43 | Phase 2: Product and requirements |
| quality standards | `02-product/07-product-lifecycle-deprecation-and-retirement.md` | 58 | 58 | 15 | 43 | Phase 2: Product and requirements |
| quality standards | `03-ux-ui-accessibility/01-accessibility.md` | 39 | 39 | 21 | 18 | Phase 3: User experience, web, and content |
| quality standards | `03-ux-ui-accessibility/02-localization-and-internationalization.md` | 23 | 23 | 12 | 11 | Phase 3: User experience, web, and content |
| quality standards | `03-ux-ui-accessibility/03-human-centered-design-and-user-research.md` | 58 | 58 | 15 | 43 | Phase 3: User experience, web, and content |
| quality standards | `03-ux-ui-accessibility/04-information-architecture-and-content-design.md` | 58 | 58 | 15 | 43 | Phase 3: User experience, web, and content |
| quality standards | `03-ux-ui-accessibility/05-interaction-and-visual-design.md` | 58 | 58 | 15 | 43 | Phase 3: User experience, web, and content |
| quality standards | `03-ux-ui-accessibility/06-usability-and-task-success.md` | 58 | 58 | 15 | 43 | Phase 3: User experience, web, and content |
| quality standards | `03-ux-ui-accessibility/07-design-systems-and-interface-reuse.md` | 58 | 58 | 15 | 43 | Phase 3: User experience, web, and content |
| quality standards | `03-ux-ui-accessibility/08-responsive-cross-device-and-adaptive-design.md` | 58 | 58 | 15 | 43 | Phase 3: User experience, web, and content |
| quality standards | `03-ux-ui-accessibility/09-ethical-design-and-dark-pattern-prevention.md` | 58 | 58 | 15 | 43 | Phase 3: User experience, web, and content |
| quality standards | `04-architecture-design/01-system-understanding-and-threat-modeling.md` | 26 | 26 | 13 | 13 | Phase 4: Architecture and design |
| quality standards | `04-architecture-design/02-performance-and-user-experience-efficiency.md` | 31 | 31 | 17 | 14 | Phase 4: Architecture and design |
| quality standards | `04-architecture-design/03-capacity-scalability-and-overload-control.md` | 33 | 33 | 14 | 19 | Phase 4: Architecture and design |
| quality standards | `04-architecture-design/04-reliability-resilience-and-failure-engineering.md` | 41 | 41 | 17 | 24 | Phase 4: Architecture and design |
| quality standards | `04-architecture-design/05-architecture-governance-and-decision-records.md` | 58 | 58 | 15 | 43 | Phase 4: Architecture and design |
| quality standards | `04-architecture-design/06-modularity-cohesion-coupling-and-boundaries.md` | 58 | 58 | 15 | 43 | Phase 4: Architecture and design |
| quality standards | `04-architecture-design/07-solid-design-principles.md` | 58 | 58 | 15 | 43 | Phase 4: Architecture and design |
| quality standards | `04-architecture-design/08-dry-kiss-yagni-and-simplicity.md` | 58 | 58 | 15 | 43 | Phase 4: Architecture and design |
| quality standards | `04-architecture-design/09-reusability-and-software-product-lines.md` | 58 | 58 | 15 | 43 | Phase 4: Architecture and design |
| quality standards | `04-architecture-design/10-state-management.md` | 58 | 58 | 15 | 43 | Phase 4: Architecture and design |
| quality standards | `04-architecture-design/11-concurrency-consistency-and-transactions.md` | 58 | 58 | 15 | 43 | Phase 4: Architecture and design |
| quality standards | `04-architecture-design/12-distributed-and-event-driven-systems.md` | 58 | 58 | 15 | 43 | Phase 4: Architecture and design |
| quality standards | `04-architecture-design/13-api-and-integration-architecture.md` | 58 | 58 | 15 | 43 | Phase 4: Architecture and design |
| quality standards | `04-architecture-design/14-evolvability-technical-debt-and-changeability.md` | 58 | 58 | 15 | 43 | Phase 4: Architecture and design |
| quality standards | `04-architecture-design/15-interoperability-portability-and-compatibility.md` | 58 | 58 | 15 | 43 | Phase 4: Architecture and design |
| quality standards | `04-architecture-design/16-sustainable-software-and-resource-efficiency.md` | 58 | 58 | 15 | 43 | Phase 4: Architecture and design |
| quality standards | `05-code-implementation/01-functional-correctness-and-business-logic.md` | 37 | 37 | 18 | 19 | Phase 5: Code quality and implementation |
| quality standards | `05-code-implementation/02-input-validation-encoding-and-safe-processing.md` | 31 | 31 | 14 | 17 | Phase 5: Code quality and implementation |
| quality standards | `05-code-implementation/03-maintainability-and-long-term-operability.md` | 34 | 34 | 14 | 20 | Phase 5: Code quality and implementation |
| quality standards | `05-code-implementation/04-universal-code-quality.md` | 65 | 65 | 22 | 43 | Phase 5: Code quality and implementation |
| quality standards | `05-code-implementation/05-readability-naming-and-style.md` | 63 | 63 | 20 | 43 | Phase 5: Code quality and implementation |
| quality standards | `05-code-implementation/06-abstractions-interfaces-and-contracts.md` | 63 | 63 | 20 | 43 | Phase 5: Code quality and implementation |
| quality standards | `05-code-implementation/07-error-handling-and-defensive-programming.md` | 65 | 65 | 22 | 43 | Phase 5: Code quality and implementation |
| quality standards | `05-code-implementation/08-resource-lifecycle-memory-and-cleanup.md` | 63 | 63 | 20 | 43 | Phase 5: Code quality and implementation |
| quality standards | `05-code-implementation/09-concurrent-asynchronous-and-parallel-code.md` | 63 | 63 | 20 | 43 | Phase 5: Code quality and implementation |
| quality standards | `05-code-implementation/10-configuration-and-feature-flag-code-quality.md` | 63 | 63 | 20 | 43 | Phase 5: Code quality and implementation |
| quality standards | `05-code-implementation/11-dependency-selection-and-hygiene.md` | 63 | 63 | 20 | 43 | Phase 5: Code quality and implementation |
| quality standards | `05-code-implementation/12-refactoring-and-legacy-code.md` | 63 | 63 | 20 | 43 | Phase 5: Code quality and implementation |
| quality standards | `05-code-implementation/13-code-review-and-work-product-review.md` | 63 | 63 | 20 | 43 | Phase 5: Code quality and implementation |
| quality standards | `05-code-implementation/14-testability.md` | 63 | 63 | 20 | 43 | Phase 5: Code quality and implementation |
| quality standards | `05-code-implementation/15-reusability-in-code.md` | 63 | 63 | 20 | 43 | Phase 5: Code quality and implementation |
| quality standards | `06-frontend/01-browser-client-usability.md` | 36 | 36 | 21 | 15 | Phase 3: User experience, web, and content |
| quality standards | `06-frontend/02-public-content-search-and-indexing.md` | 22 | 22 | 12 | 10 | Phase 3: User experience, web, and content |
| quality standards | `06-frontend/03-offline-and-progressive-web-capabilities.md` | 22 | 22 | 9 | 13 | Phase 3: User experience, web, and content |
| quality standards | `06-frontend/04-frontend-architecture-components-and-state.md` | 63 | 63 | 20 | 43 | Phase 3: User experience, web, and content |
| quality standards | `06-frontend/05-forms-user-input-and-client-validation.md` | 63 | 63 | 20 | 43 | Phase 3: User experience, web, and content |
| quality standards | `06-frontend/06-frontend-performance-and-runtime-efficiency.md` | 63 | 63 | 20 | 43 | Phase 3: User experience, web, and content |
| quality standards | `06-frontend/07-frontend-testing.md` | 63 | 63 | 20 | 43 | Phase 3: User experience, web, and content |
| quality standards | `06-frontend/08-frontend-security-and-privacy.md` | 63 | 63 | 20 | 43 | Phase 3: User experience, web, and content |
| quality standards | `07-backend-services/01-apis-webhooks-and-integrations.md` | 47 | 47 | 31 | 16 | Phase 6: Application services and APIs |
| quality standards | `07-backend-services/02-authentication-enrollment-and-recovery.md` | 46 | 46 | 25 | 21 | Phase 6: Application services and APIs |
| quality standards | `07-backend-services/03-authorization-and-tenant-isolation.md` | 31 | 31 | 12 | 19 | Phase 6: Application services and APIs |
| quality standards | `07-backend-services/04-sessions-cookies-and-tokens.md` | 28 | 28 | 13 | 15 | Phase 6: Application services and APIs |
| quality standards | `07-backend-services/05-file-and-content-processing.md` | 31 | 31 | 13 | 18 | Phase 6: Application services and APIs |
| quality standards | `07-backend-services/06-multi-tenant-systems.md` | 22 | 22 | 8 | 14 | Phase 6: Application services and APIs |
| quality standards | `07-backend-services/07-high-risk-admin-and-support-tools.md` | 24 | 24 | 15 | 9 | Phase 6: Application services and APIs |
| quality standards | `07-backend-services/08-enterprise-identity-provisioning-and-organization-lifecycle.md` | 20 | 20 | 12 | 8 | Phase 6: Application services and APIs |
| quality standards | `07-backend-services/09-backend-service-design.md` | 63 | 63 | 20 | 43 | Phase 6: Application services and APIs |
| quality standards | `07-backend-services/10-background-jobs-scheduling-and-queues.md` | 63 | 63 | 20 | 43 | Phase 6: Application services and APIs |
| quality standards | `07-backend-services/11-caching-search-and-derived-data.md` | 63 | 63 | 20 | 43 | Phase 6: Application services and APIs |
| quality standards | `08-data/01-data-stores-queues-caches-search-and-integrity.md` | 51 | 51 | 21 | 30 | Phase 7: Data and information lifecycle |
| quality standards | `08-data/02-import-export-bulk-operations-and-portability.md` | 22 | 22 | 14 | 8 | Phase 7: Data and information lifecycle |
| quality standards | `08-data/03-data-governance-ownership-and-accountability.md` | 63 | 63 | 20 | 43 | Phase 7: Data and information lifecycle |
| quality standards | `08-data/04-data-architecture-and-modeling.md` | 63 | 63 | 20 | 43 | Phase 7: Data and information lifecycle |
| quality standards | `08-data/05-data-quality.md` | 63 | 63 | 20 | 43 | Phase 7: Data and information lifecycle |
| quality standards | `08-data/06-metadata-lineage-provenance-and-catalogs.md` | 63 | 63 | 20 | 43 | Phase 7: Data and information lifecycle |
| quality standards | `08-data/07-schema-evolution-migrations-and-data-repair.md` | 63 | 63 | 20 | 43 | Phase 7: Data and information lifecycle |
| quality standards | `08-data/08-analytics-bi-and-decision-data.md` | 63 | 63 | 20 | 43 | Phase 7: Data and information lifecycle |
| quality standards | `09-security/01-dependencies-sbom-and-licenses.md` | 28 | 28 | 15 | 13 | Phase 8: Security and cryptography |
| quality standards | `09-security/02-transport-browser-dns-and-network-security.md` | 45 | 45 | 15 | 30 | Phase 8: Security and cryptography |
| quality standards | `09-security/03-cryptography-and-key-management.md` | 31 | 31 | 13 | 18 | Phase 8: Security and cryptography |
| quality standards | `09-security/04-security-validation-and-vulnerability-management.md` | 35 | 35 | 14 | 21 | Phase 8: Security and cryptography |
| quality standards | `09-security/05-incident-response-and-crisis-readiness.md` | 32 | 32 | 14 | 18 | Phase 8: Security and cryptography |
| quality standards | `09-security/06-public-forms-anonymous-access-and-abuse-resistance.md` | 18 | 18 | 10 | 8 | Phase 8: Security and cryptography |
| quality standards | `09-security/07-security-governance-and-risk-management.md` | 63 | 63 | 20 | 43 | Phase 8: Security and cryptography |
| quality standards | `09-security/08-secure-software-development-lifecycle.md` | 63 | 63 | 20 | 43 | Phase 8: Security and cryptography |
| quality standards | `09-security/09-threat-modeling-and-abuse-case-analysis.md` | 63 | 63 | 20 | 43 | Phase 8: Security and cryptography |
| quality standards | `09-security/10-application-security-engineering.md` | 63 | 63 | 20 | 43 | Phase 8: Security and cryptography |
| quality standards | `09-security/11-identity-access-and-privileged-security.md` | 63 | 63 | 20 | 43 | Phase 8: Security and cryptography |
| quality standards | `09-security/12-secrets-management.md` | 63 | 63 | 20 | 43 | Phase 8: Security and cryptography |
| quality standards | `09-security/13-software-supply-chain-security.md` | 63 | 63 | 20 | 43 | Phase 8: Security and cryptography |
| quality standards | `09-security/14-security-testing-and-assurance.md` | 63 | 63 | 20 | 43 | Phase 8: Security and cryptography |
| quality standards | `09-security/15-security-monitoring-detection-and-response.md` | 63 | 63 | 20 | 43 | Phase 8: Security and cryptography |
| quality standards | `09-security/16-zero-trust-and-privileged-access.md` | 63 | 63 | 20 | 43 | Phase 8: Security and cryptography |
| quality standards | `09-security/17-fraud-abuse-and-automation-resistance.md` | 63 | 63 | 20 | 43 | Phase 8: Security and cryptography |
| quality standards | `10-privacy/01-privacy-and-data-protection.md` | 41 | 41 | 16 | 25 | Phase 9: Privacy and data protection |
| quality standards | `10-privacy/02-privacy-governance-and-engineering.md` | 63 | 63 | 20 | 43 | Phase 9: Privacy and data protection |
| quality standards | `10-privacy/03-notices-consent-and-preferences.md` | 63 | 63 | 20 | 43 | Phase 9: Privacy and data protection |
| quality standards | `10-privacy/04-data-rights-retention-deletion-and-portability.md` | 63 | 63 | 20 | 43 | Phase 9: Privacy and data protection |
| quality standards | `10-privacy/05-deidentification-pseudonymization-and-anonymization.md` | 63 | 63 | 20 | 43 | Phase 9: Privacy and data protection |
| quality standards | `10-privacy/06-cross-border-transfers-processors-and-sharing.md` | 63 | 63 | 20 | 43 | Phase 9: Privacy and data protection |
| quality standards | `10-privacy/07-childrens-sensitive-and-high-risk-data.md` | 63 | 63 | 20 | 43 | Phase 9: Privacy and data protection |
| quality standards | `11-testing-quality-assurance/01-test-strategy-and-evidence.md` | 35 | 35 | 13 | 22 | Phase 10: Verification and testing |
| quality standards | `11-testing-quality-assurance/02-unit-and-component-testing.md` | 63 | 63 | 20 | 43 | Phase 10: Verification and testing |
| quality standards | `11-testing-quality-assurance/03-integration-and-contract-testing.md` | 63 | 63 | 20 | 43 | Phase 10: Verification and testing |
| quality standards | `11-testing-quality-assurance/04-system-end-to-end-and-acceptance-testing.md` | 63 | 63 | 20 | 43 | Phase 10: Verification and testing |
| quality standards | `11-testing-quality-assurance/05-regression-change-impact-and-compatibility-testing.md` | 63 | 63 | 20 | 43 | Phase 10: Verification and testing |
| quality standards | `11-testing-quality-assurance/06-test-automation-quality.md` | 63 | 63 | 20 | 43 | Phase 10: Verification and testing |
| quality standards | `11-testing-quality-assurance/07-test-data-and-test-environments.md` | 63 | 63 | 20 | 43 | Phase 10: Verification and testing |
| quality standards | `11-testing-quality-assurance/08-nonfunctional-quality-attribute-testing.md` | 63 | 63 | 20 | 43 | Phase 10: Verification and testing |
| quality standards | `11-testing-quality-assurance/09-exploratory-risk-based-and-usability-testing.md` | 63 | 63 | 20 | 43 | Phase 10: Verification and testing |
| quality standards | `11-testing-quality-assurance/10-static-analysis-formal-methods-and-verification.md` | 63 | 63 | 20 | 43 | Phase 10: Verification and testing |
| quality standards | `11-testing-quality-assurance/11-defect-management-and-root-cause.md` | 63 | 63 | 20 | 43 | Phase 10: Verification and testing |
| quality standards | `12-delivery-cicd/01-source-control-and-change-management.md` | 24 | 24 | 9 | 15 | Phase 11: Developer experience, platform, and delivery |
| quality standards | `12-delivery-cicd/02-build-cicd-and-artifact-integrity.md` | 28 | 28 | 11 | 17 | Phase 11: Developer experience, platform, and delivery |
| quality standards | `12-delivery-cicd/03-environments-configuration-flags-and-secrets.md` | 47 | 47 | 24 | 23 | Phase 11: Developer experience, platform, and delivery |
| quality standards | `12-delivery-cicd/04-continuous-integration.md` | 63 | 63 | 20 | 43 | Phase 11: Developer experience, platform, and delivery |
| quality standards | `12-delivery-cicd/05-continuous-delivery-and-deployment.md` | 63 | 63 | 20 | 43 | Phase 11: Developer experience, platform, and delivery |
| quality standards | `12-delivery-cicd/06-reproducible-hermetic-and-trusted-builds.md` | 63 | 63 | 20 | 43 | Phase 11: Developer experience, platform, and delivery |
| quality standards | `12-delivery-cicd/07-infrastructure-as-code-and-environment-engineering.md` | 63 | 63 | 20 | 43 | Phase 11: Developer experience, platform, and delivery |
| quality standards | `12-delivery-cicd/08-deployment-strategies-and-progressive-delivery.md` | 63 | 63 | 20 | 43 | Phase 11: Developer experience, platform, and delivery |
| quality standards | `12-delivery-cicd/09-database-and-data-release-engineering.md` | 63 | 63 | 20 | 43 | Phase 11: Developer experience, platform, and delivery |
| quality standards | `12-delivery-cicd/10-rollback-rollforward-and-kill-switches.md` | 63 | 63 | 20 | 43 | Phase 11: Developer experience, platform, and delivery |
| quality standards | `13-operations-sre/01-infrastructure-and-platform-readiness.md` | 55 | 55 | 21 | 34 | Phase 12: Operations, SRE, and support |
| quality standards | `13-operations-sre/02-observability-logging-metrics-traces-and-audit.md` | 38 | 38 | 13 | 25 | Phase 12: Operations, SRE, and support |
| quality standards | `13-operations-sre/03-alerting-oncall-and-runbooks.md` | 32 | 32 | 15 | 17 | Phase 12: Operations, SRE, and support |
| quality standards | `13-operations-sre/04-backup-disaster-recovery-and-business-continuity.md` | 33 | 33 | 17 | 16 | Phase 12: Operations, SRE, and support |
| quality standards | `13-operations-sre/05-sli-slo-error-budgets-and-reliability-governance.md` | 63 | 63 | 20 | 43 | Phase 12: Operations, SRE, and support |
| quality standards | `13-operations-sre/06-service-management-and-customer-support.md` | 63 | 63 | 20 | 43 | Phase 12: Operations, SRE, and support |
| quality standards | `13-operations-sre/07-incident-problem-management-and-postmortems.md` | 63 | 63 | 20 | 43 | Phase 12: Operations, SRE, and support |
| quality standards | `13-operations-sre/08-operational-capacity-performance-and-efficiency.md` | 63 | 63 | 20 | 43 | Phase 12: Operations, SRE, and support |
| quality standards | `13-operations-sre/09-cost-efficiency-and-finops.md` | 63 | 63 | 20 | 43 | Phase 12: Operations, SRE, and support |
| quality standards | `13-operations-sre/10-asset-patch-and-technology-lifecycle-management.md` | 63 | 63 | 20 | 43 | Phase 12: Operations, SRE, and support |
| quality standards | `13-operations-sre/11-decommissioning-and-secure-retirement.md` | 63 | 63 | 20 | 43 | Phase 12: Operations, SRE, and support |
| quality standards | `13-operations-sre/12-operational-access-and-production-change.md` | 63 | 63 | 20 | 43 | Phase 12: Operations, SRE, and support |
| quality standards | `14-documentation/01-documentation-governance.md` | 63 | 63 | 20 | 43 | Phase 13: Documentation and knowledge |
| quality standards | `14-documentation/02-product-requirements-and-decision-documentation.md` | 63 | 63 | 20 | 43 | Phase 13: Documentation and knowledge |
| quality standards | `14-documentation/03-architecture-and-adr-documentation.md` | 63 | 63 | 20 | 43 | Phase 13: Documentation and knowledge |
| quality standards | `14-documentation/04-code-api-and-data-documentation.md` | 63 | 63 | 20 | 43 | Phase 13: Documentation and knowledge |
| quality standards | `14-documentation/05-user-help-training-and-support-content.md` | 63 | 63 | 20 | 43 | Phase 13: Documentation and knowledge |
| quality standards | `14-documentation/06-operations-runbooks-release-and-incident-documentation.md` | 63 | 63 | 20 | 43 | Phase 13: Documentation and knowledge |
| quality standards | `15-conditional-domains/01-payments-billing-and-money-movement.md` | 28 | 28 | 18 | 10 | Phase 16: Specialized domains and release assurance |
| quality standards | `15-conditional-domains/02-user-generated-content-social-and-marketplaces.md` | 26 | 26 | 17 | 9 | Phase 14: Trust, safety, and ecosystems |
| quality standards | `15-conditional-domains/03-email-sms-push-and-notifications.md` | 26 | 26 | 14 | 12 | Phase 14: Trust, safety, and ecosystems |
| quality standards | `15-conditional-domains/04-realtime-streaming-and-collaboration.md` | 23 | 23 | 12 | 11 | Phase 16: Specialized domains and release assurance |
| quality standards | `15-conditional-domains/05-ai-machine-learning-and-llm-systems.md` | 40 | 40 | 25 | 15 | Phase 15: AI, ML, and AI-assisted development |
| quality standards | `15-conditional-domains/06-safety-critical-and-physically-consequential-systems.md` | 23 | 23 | 5 | 18 | Phase 16: Specialized domains and release assurance |
| quality standards | `15-conditional-domains/07-media-voice-video-and-live-communications.md` | 21 | 21 | 13 | 8 | Phase 16: Specialized domains and release assurance |
| quality standards | `15-conditional-domains/08-geolocation-sensors-and-device-capabilities.md` | 20 | 20 | 12 | 8 | Phase 16: Specialized domains and release assurance |
| quality standards | `15-conditional-domains/09-search-ranking-recommendations-and-personalization.md` | 21 | 21 | 13 | 8 | Phase 14: Trust, safety, and ecosystems |
| quality standards | `15-conditional-domains/10-analytics-experimentation-advertising-and-attribution.md` | 23 | 23 | 15 | 8 | Phase 16: Specialized domains and release assurance |
| quality standards | `15-conditional-domains/11-ecommerce-orders-inventory-fulfillment-and-returns.md` | 20 | 20 | 12 | 8 | Phase 16: Specialized domains and release assurance |
| quality standards | `15-conditional-domains/12-blockchain-smart-contracts-and-irreversible-ledgers.md` | 20 | 20 | 12 | 8 | Phase 16: Specialized domains and release assurance |
| quality standards | `15-conditional-domains/13-iot-device-control-industrial-and-cyber-physical.md` | 20 | 20 | 12 | 8 | Phase 16: Specialized domains and release assurance |
| quality standards | `15-conditional-domains/14-open-source-projects-and-public-packages.md` | 63 | 63 | 20 | 43 | Phase 16: Specialized domains and release assurance |
| quality standards | `15-conditional-domains/15-mobile-desktop-and-installed-clients.md` | 63 | 63 | 20 | 43 | Phase 16: Specialized domains and release assurance |
| quality standards | `15-conditional-domains/16-cloud-native-and-managed-platforms.md` | 63 | 63 | 20 | 43 | Phase 16: Specialized domains and release assurance |
| quality standards | `15-conditional-domains/17-containers-and-orchestration.md` | 63 | 63 | 20 | 43 | Phase 16: Specialized domains and release assurance |
| quality standards | `15-conditional-domains/18-serverless-and-managed-execution.md` | 63 | 63 | 20 | 43 | Phase 16: Specialized domains and release assurance |
| quality standards | `15-conditional-domains/19-healthcare-and-regulated-sensitive-systems.md` | 63 | 63 | 20 | 43 | Phase 16: Specialized domains and release assurance |
| quality standards | `15-conditional-domains/20-public-sector-and-high-assurance-systems.md` | 63 | 63 | 20 | 43 | Phase 16: Specialized domains and release assurance |
| quality standards | `15-conditional-domains/21-libraries-sdks-clis-and-developer-tools.md` | 63 | 63 | 20 | 43 | Phase 16: Specialized domains and release assurance |
| quality standards | `15-conditional-domains/22-games-immersive-and-high-engagement-products.md` | 63 | 63 | 20 | 43 | Phase 16: Specialized domains and release assurance |
| quality standards | `16-release-production/01-universal-production-readiness-master-checklist.md` | 1,425 | 0 | 0 | 1,425 | covered by the repository's stable PRC production-readiness controls |
| quality standards | `16-release-production/02-release-control-record.md` | 31 | 31 | 18 | 13 | Phase 16: Specialized domains and release assurance |
| quality standards | `16-release-production/03-deployment-and-release-engineering.md` | 38 | 38 | 18 | 20 | Phase 16: Specialized domains and release assurance |
| quality standards | `16-release-production/04-post-deployment-verification.md` | 30 | 30 | 14 | 16 | Phase 16: Specialized domains and release assurance |
| quality standards | `16-release-production/05-required-signoffs.md` | 21 | 21 | 13 | 8 | Phase 16: Specialized domains and release assurance |
| quality standards | `16-release-production/06-go-conditional-go-and-no-go.md` | 32 | 32 | 12 | 20 | Phase 16: Specialized domains and release assurance |
| quality standards | `INDEX.md` | 8 | 0 | 0 | 8 | package documentation or manifest, not a product control checklist |
| quality standards | `MANIFEST.md` | 0 | 0 | 0 | 0 | package documentation or manifest, not a product control checklist |
| quality standards | `README.md` | 0 | 0 | 0 | 0 | package documentation or manifest, not a product control checklist |
| gap supplement | `00-gap-audit-and-consolidation-map.md` | 10 | 0 | 0 | 10 | package documentation, manifest, or adoption audit rather than product controls |
| gap supplement | `01-web-quality-seo-accessibility-and-content.md` | 592 | 254 | 254 | 338 | Phase 3: User experience, web, and content |
| gap supplement | `02-universal-code-quality-design-and-implementation.md` | 726 | 215 | 205 | 521 | Phase 5: Code quality and implementation |
| gap supplement | `03-developer-experience-platform-engineering-and-economics.md` | 650 | 213 | 203 | 447 | Phase 11: Developer experience, platform, and delivery |
| gap supplement | `04-trust-safety-content-integrity-and-extension-ecosystems.md` | 303 | 201 | 191 | 112 | Phase 14: Trust, safety, and ecosystems |
| gap supplement | `05-ai-ml-mlops-and-ai-assisted-development.md` | 436 | 328 | 318 | 118 | Phase 15: AI, ML, and AI-assisted development |
| gap supplement | `06-data-contracts-semantics-records-and-preservation.md` | 486 | 256 | 246 | 240 | Phase 7: Data and information lifecycle |
| gap supplement | `07-cryptographic-agility-obsolescence-and-long-term-resilience.md` | 338 | 212 | 202 | 136 | Phase 8: Security and cryptography |
| gap supplement | `MANIFEST.md` | 0 | 0 | 0 | 0 | package documentation, manifest, or adoption audit rather than product controls |
| gap supplement | `README.md` | 0 | 0 | 0 | 0 | package documentation, manifest, or adoption audit rather than product controls |

## Interpretation boundary

This project is an independent implementation-oriented synthesis. It does not reproduce or replace authoritative standards, laws, contracts, regulatory rules, or certification criteria. Applicable authoritative requirements prevail.
