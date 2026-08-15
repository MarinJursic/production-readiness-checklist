# Environments, quality, and experience

> Verify configuration, functional correctness, testing, frontend behavior, and accessibility.

Sections 10–14 of the master checklist. For each applicable item, capture status, owner, evidence, release, environment, reviewer, and evidence age.

## 10. Environments, configuration, and secrets

### 10.1 Environment separation

- [ ] **PRC-10-001** — Development, testing, staging, and production are meaningfully separated.
- [ ] **PRC-10-002** — Production credentials do not work in non-production environments.
- [ ] **PRC-10-003** — Non-production users cannot reach production control planes.
- [ ] **PRC-10-004** — Staging is sufficiently production-like to validate topology, integrations, migrations, and deployment behavior.
- [ ] **PRC-10-005** — Production data is not copied to non-production without authorization and protection.
- [ ] **PRC-10-006** — Test accounts, bypasses, simulators, and debug endpoints are absent or inaccessible in production.
- [ ] **PRC-10-007** — Non-production email, payment, notifications, and webhooks cannot accidentally contact real users.
- [ ] **PRC-10-008** — Environment boundaries are documented and tested.
- [ ] **PRC-10-009** — Production access is logged and periodically reviewed.

### 10.2 Configuration

- [ ] **PRC-10-010** — Configuration is version-controlled or otherwise auditable.
- [ ] **PRC-10-011** — Configuration changes require review and approval.
- [ ] **PRC-10-012** — Configuration has schema and type validation.
- [ ] **PRC-10-013** — Missing configuration fails safely rather than activating insecure defaults.
- [ ] **PRC-10-014** — Secure defaults are used.
- [ ] **PRC-10-015** — Environment-specific values are explicit.
- [ ] **PRC-10-016** — Configuration drift is detected.
- [ ] **PRC-10-017** — Production configuration can be reconstructed.
- [ ] **PRC-10-018** — Configuration rollback is tested.
- [ ] **PRC-10-019** — Debug mode, verbose exception output, profiling, test routes, and development consoles are disabled.
- [ ] **PRC-10-020** — Time zones, locales, character encodings, and units are explicit.
- [ ] **PRC-10-021** — Time synchronization is monitored.
- [ ] **PRC-10-022** — Resource limits and timeouts are explicitly configured.
- [ ] **PRC-10-023** — Production source maps and diagnostic artifacts do not expose secrets or inappropriate internal information.
- [ ] **PRC-10-024** — Configuration ownership and expiration dates are recorded.
- [ ] **PRC-10-025** — Emergency configuration changes are captured and reviewed.

### 10.3 Feature flags

- [ ] **PRC-10-026** — Every flag has an owner.
- [ ] **PRC-10-027** — Every flag has a documented purpose.
- [ ] **PRC-10-028** — Default behavior is safe when the flag service fails.
- [ ] **PRC-10-029** — Flag state is included in release and incident records.
- [ ] **PRC-10-030** — Access to change high-impact flags is restricted and audited.
- [ ] **PRC-10-031** — Flag combinations are tested.
- [ ] **PRC-10-032** — A kill switch exists for high-risk new functionality where appropriate.
- [ ] **PRC-10-033** — Flag changes can be rolled back quickly.
- [ ] **PRC-10-034** — Temporary flags have removal dates.
- [ ] **PRC-10-035** — Stale flags and dead code are removed.

### 10.4 Secrets

- [ ] **PRC-10-036** — Secrets are stored in an approved secret-management system.
- [ ] **PRC-10-037** — No secrets exist in source, images, frontend bundles, logs, tickets, documentation, chat history, or analytics.
- [ ] **PRC-10-038** — Secrets are unique across environments.
- [ ] **PRC-10-039** — Each workload uses its own scoped identity or credential.
- [ ] **PRC-10-040** — Credentials follow least privilege.
- [ ] **PRC-10-041** — Short-lived credentials are preferred where supported.
- [ ] **PRC-10-042** — Secret retrieval and use are audited.
- [ ] **PRC-10-043** — Secret rotation has been tested.
- [ ] **PRC-10-044** — Emergency revocation has been tested.
- [ ] **PRC-10-045** — Credential expiry and rotation failures are monitored.
- [ ] **PRC-10-046** — Old credentials are invalidated after rotation.
- [ ] **PRC-10-047** — Backup and recovery procedures preserve or reconstruct necessary secrets securely.
- [ ] **PRC-10-048** — Break-glass credentials are protected, monitored, and periodically tested.
- [ ] **PRC-10-049** — No default or vendor-supplied credentials remain.
- [ ] **PRC-10-050** — Suspected historical secret exposure has been addressed by rotation, not merely deletion from source.

---

## 11. Functional correctness and business logic

- [ ] **PRC-11-001** — Every requirement maps to at least one verification activity.
- [ ] **PRC-11-002** — Every critical user journey passes end to end.
- [ ] **PRC-11-003** — Positive, negative, boundary, and invalid-input cases pass.
- [ ] **PRC-11-004** — Empty-state behavior is correct.
- [ ] **PRC-11-005** — Maximum-size and high-volume behavior is correct.
- [ ] **PRC-11-006** — Duplicate submissions do not create unintended duplicate effects.
- [ ] **PRC-11-007** — Refresh, retry, browser navigation, and network interruption do not corrupt state.
- [ ] **PRC-11-008** — Concurrent requests do not violate business invariants.
- [ ] **PRC-11-009** — Race conditions have been tested around balances, inventory, permissions, limits, and state transitions.
- [ ] **PRC-11-010** — State transitions reject invalid orderings.
- [ ] **PRC-11-011** — Operations that must be idempotent are demonstrably idempotent.
- [ ] **PRC-11-012** — Partial failures cannot leave an ambiguous or silently inconsistent result.
- [ ] **PRC-11-013** — Transactions are atomic where required.
- [ ] **PRC-11-014** — Eventually consistent operations expose correct pending and completion states.
- [ ] **PRC-11-015** — Reconciliation detects missing, duplicated, or mismatched records.
- [ ] **PRC-11-016** — Numeric precision and rounding are correct.
- [ ] **PRC-11-017** — Money calculations use correct decimal, currency, tax, and rounding rules.
- [ ] **PRC-11-018** — Time calculations handle time zones, daylight-saving changes, leap years, month boundaries, and clock differences.
- [ ] **PRC-11-019** — Unicode, normalization, collation, and case behavior do not create duplicates or bypasses.
- [ ] **PRC-11-020** — Limits and quotas are enforced server-side.
- [ ] **PRC-11-021** — Approval and four-eyes workflows cannot be bypassed.
- [ ] **PRC-11-022** — Cancellation, refund, reversal, deletion, and undo operations are complete.
- [ ] **PRC-11-023** — Account suspension and termination affect every relevant session, token, worker, and integration.
- [ ] **PRC-11-024** — Import and export behavior is complete and validated.
- [ ] **PRC-11-025** — Long-running operations can resume safely or fail clearly.
- [ ] **PRC-11-026** — Scheduled jobs handle missed runs, overlapping runs, retries, and duplicate execution.
- [ ] **PRC-11-027** — Background jobs preserve user, authorization, and tenant context.
- [ ] **PRC-11-028** — Dead-lettered operations can be inspected and replayed safely.
- [ ] **PRC-11-029** — Administrative actions produce the intended effect and audit trail.
- [ ] **PRC-11-030** — Support tools cannot create states impossible through the normal product.
- [ ] **PRC-11-031** — Failure defaults to the safer state.
- [ ] **PRC-11-032** — Error messages tell users what they can do next without exposing sensitive details.
- [ ] **PRC-11-033** — Known limitations are documented accurately.

---

## 12. Testing strategy and evidence

OWASP ASVS 5.0.0 is the current stable application-security verification standard and provides testable requirements independent of a specific implementation technology. OWASP WSTG and SAMM complement it with testing and program-lifecycle practices. ([owasp.org](https://owasp.org/www-project-application-security-verification-standard/))

- [ ] **PRC-12-001** — A documented test strategy exists for the release risk level.
- [ ] **PRC-12-002** — Tests cover unit, component, integration, contract, system, end-to-end, smoke, and regression levels as appropriate.
- [ ] **PRC-12-003** — Tests prioritize critical and high-risk behavior rather than relying solely on aggregate coverage.
- [ ] **PRC-12-004** — Coverage gaps in critical code are explicitly reviewed.
- [ ] **PRC-12-005** — Test data is representative of real size, shape, encoding, and edge conditions.
- [ ] **PRC-12-006** — Personal or confidential test data is synthetic, masked, or properly authorized.
- [ ] **PRC-12-007** — Tests are deterministic enough to act as release gates.
- [ ] **PRC-12-008** — Flaky gating tests are fixed or formally quarantined with risk review.
- [ ] **PRC-12-009** — Test failures cannot be dismissed without evidence.
- [ ] **PRC-12-010** — Tests validate failure and recovery, not only successful operation.
- [ ] **PRC-12-011** — Contract tests cover internal and external interfaces.
- [ ] **PRC-12-012** — Compatibility tests cover supported clients and versions.
- [ ] **PRC-12-013** — Migration, upgrade, mixed-version, downgrade, rollback, and roll-forward scenarios are tested.
- [ ] **PRC-12-014** — Performance tests cover load, stress, spike, endurance, and capacity limits.
- [ ] **PRC-12-015** — Resilience tests inject relevant dependency, network, storage, and resource failures.
- [ ] **PRC-12-016** — Security testing covers source, dependencies, deployed behavior, infrastructure, configuration, and business logic.
- [ ] **PRC-12-017** — Accessibility testing combines automation with manual keyboard and assistive-technology checks.
- [ ] **PRC-12-018** — Property-based, fuzz, mutation, model-based, or combinatorial testing is used for high-risk parsers and stateful logic where useful.
- [ ] **PRC-12-019** — Critical calculations use independent test or reconciliation methods.
- [ ] **PRC-12-020** — Tests run against production-like topology and configuration.
- [ ] **PRC-12-021** — Tests verify the actual artifact intended for production.
- [ ] **PRC-12-022** — Evidence identifies test version, data, environment, build, and result.
- [ ] **PRC-12-023** — Failed tests are linked to fixes and successful retests.
- [ ] **PRC-12-024** — Material changes trigger the relevant regression suite.
- [ ] **PRC-12-025** — Independent user acceptance is complete.
- [ ] **PRC-12-026** — An independent penetration test is completed where risk, contract, or regulation warrants it.
- [ ] **PRC-12-027** — Testing limitations and untested conditions are disclosed in the readiness decision.

---

## 13. Frontend, usability, and browser behavior

- [ ] **PRC-13-001** — A supported browser and device matrix is documented.
- [ ] **PRC-13-002** — Unsupported clients receive reasonable behavior or a clear message.
- [ ] **PRC-13-003** — Critical journeys work on every supported browser.
- [ ] **PRC-13-004** — Layout works at supported viewport sizes.
- [ ] **PRC-13-005** — Portrait and landscape behavior is correct where relevant.
- [ ] **PRC-13-006** — Zoom and text enlargement do not hide essential functionality.
- [ ] **PRC-13-007** — Touch, mouse, keyboard, and assistive input work appropriately.
- [ ] **PRC-13-008** — Loading, empty, success, warning, error, partial, timeout, and offline states are designed.
- [ ] **PRC-13-009** — Forms preserve appropriate input after recoverable errors.
- [ ] **PRC-13-010** — Validation is understandable and points to the affected field.
- [ ] **PRC-13-011** — Back, forward, refresh, bookmarks, and deep links behave correctly.
- [ ] **PRC-13-012** — Multiple tabs and windows do not corrupt session or workflow state.
- [ ] **PRC-13-013** — Client-side caching does not expose another user’s data.
- [ ] **PRC-13-014** — Browser storage contains no unnecessary sensitive data.
- [ ] **PRC-13-015** — Logout and account deletion clear appropriate client state.
- [ ] **PRC-13-016** — Client-side authorization is never treated as the security boundary.
- [ ] **PRC-13-017** — JavaScript-disabled or failed-script behavior degrades acceptably where required.
- [ ] **PRC-13-018** — Third-party script failure does not break critical operation unnecessarily.
- [ ] **PRC-13-019** — No unhandled exceptions or material console errors remain.
- [ ] **PRC-13-020** — Source maps, diagnostic data, and stack traces are appropriately restricted.
- [ ] **PRC-13-021** — Error pages do not expose internals.
- [ ] **PRC-13-022** — Download, print, copy, paste, and upload behaviors work where supported.
- [ ] **PRC-13-023** — Long text, translated text, unusual names, and bidirectional text are handled.
- [ ] **PRC-13-024** — Destructive operations require suitable confirmation or reversal.
- [ ] **PRC-13-025** — Navigation and controls are consistent.
- [ ] **PRC-13-026** — Users are not trapped in dead ends.
- [ ] **PRC-13-027** — Session timeout warnings and recovery preserve work where appropriate.
- [ ] **PRC-13-028** — Autosave behavior is transparent and reliable.
- [ ] **PRC-13-029** — User feedback is not lost during asynchronous operations.
- [ ] **PRC-13-030** — Client and server validation produce consistent outcomes.
- [ ] **PRC-13-031** — Analytics and experimentation code do not alter critical behavior.
- [ ] **PRC-13-032** — Dark patterns, deceptive defaults, and misleading urgency are absent.

---

## 14. Accessibility

WCAG 2.2 is a technology-neutral W3C Recommendation. It requires both testable technical conformance and human evaluation; W3C also notes that even the highest conformance level cannot meet every possible user need. A strong general baseline is **WCAG 2.2 Level AA**, with stricter controls where law, contract, or user risk requires them. ([w3.org](https://www.w3.org/TR/WCAG22/))

- [ ] **PRC-14-001** — The required WCAG version and conformance level are documented.
- [ ] **PRC-14-002** — Automated accessibility scans pass within the accepted threshold.
- [ ] **PRC-14-003** — Critical journeys have been tested manually using only a keyboard.
- [ ] **PRC-14-004** — No keyboard trap exists.
- [ ] **PRC-14-005** — Focus order follows a meaningful sequence.
- [ ] **PRC-14-006** — Focus is visibly indicated.
- [ ] **PRC-14-007** — Focus is not obscured by sticky headers, dialogs, banners, or overlays.
- [ ] **PRC-14-008** — Page titles, headings, landmarks, regions, and lists use meaningful semantics.
- [ ] **PRC-14-009** — Every interactive control has an accessible name, role, state, and value.
- [ ] **PRC-14-010** — Custom controls expose equivalent semantics and keyboard behavior.
- [ ] **PRC-14-011** — Images have appropriate alternative text or are correctly marked decorative.
- [ ] **PRC-14-012** — Audio and video have required captions, transcripts, and descriptions.
- [ ] **PRC-14-013** — Information is not communicated by color alone.
- [ ] **PRC-14-014** — Text and meaningful visual elements meet contrast requirements.
- [ ] **PRC-14-015** — Content reflows and remains usable when zoomed.
- [ ] **PRC-14-016** — Text spacing changes do not break content.
- [ ] **PRC-14-017** — Touch and pointer targets are sufficiently usable.
- [ ] **PRC-14-018** — Dragging has a non-drag alternative where required.
- [ ] **PRC-14-019** — Motion, animation, autoplay, and flashing are controlled safely.
- [ ] **PRC-14-020** — “Reduce motion” or equivalent preferences are respected where relevant.
- [ ] **PRC-14-021** — Forms have explicit labels and understandable instructions.
- [ ] **PRC-14-022** — Required fields and formats are conveyed accessibly.
- [ ] **PRC-14-023** — Errors are identified, described, and associated with fields.
- [ ] **PRC-14-024** — Error summaries and status changes are announced to assistive technologies.
- [ ] **PRC-14-025** — Authentication does not rely solely on inaccessible puzzles or memory tests.
- [ ] **PRC-14-026** — Time limits can be extended or disabled where required.
- [ ] **PRC-14-027** — Repeated entry of the same information is minimized.
- [ ] **PRC-14-028** — Help is located consistently.
- [ ] **PRC-14-029** — Language and text direction are declared correctly.
- [ ] **PRC-14-030** — Tables expose proper headers and relationships.
- [ ] **PRC-14-031** — Dialogs, menus, tabs, tooltips, and notifications follow accessible interaction patterns.
- [ ] **PRC-14-032** — Dynamic content changes are announced appropriately.
- [ ] **PRC-14-033** — Downloaded documents are accessible.
- [ ] **PRC-14-034** — Third-party embedded components meet the required level or have an accessible alternative.
- [ ] **PRC-14-035** — Testing includes representative screen-reader combinations.
- [ ] **PRC-14-036** — Testing includes users with disabilities where the product’s impact warrants it.
- [ ] **PRC-14-037** — An accessibility statement and issue-reporting path exist.
- [ ] **PRC-14-038** — Accessibility defects have owners, severity, and remediation dates.
- [ ] **PRC-14-039** — Accessibility remains part of regression testing after launch.

---
