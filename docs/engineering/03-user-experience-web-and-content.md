# User experience, web, and content

_Phase 3 of 16 in the [complete engineering review](00-overview.md)._

Research, accessibility, localization, interaction design, frontend quality, content, SEO, and web-platform behavior.

Assess every control as **Pass**, **Fail**, **Blocked**, or **Not Applicable**. A Pass needs current evidence for the exact product, revision, artifact, configuration, environment, and data state under review.

## Accessibility

_Consolidated from `quality standards/03-ux-ui-accessibility/01-accessibility.md`; 21 non-duplicative controls._

### Universal controls

- [ ] **USEQ-3F288998** — The required accessibility standard, version, and conformance level are documented.
- [ ] **USEQ-898B0514** — Automated accessibility scans pass within the approved threshold.
- [ ] **USEQ-F96520F3** — Focus order is meaningful, visibly indicated, and not obscured by overlays, sticky regions, dialogs, or banners.
- [ ] **USEQ-9D05AB06** — Page titles, headings, landmarks, regions, lists, and tables use meaningful semantics.
- [ ] **USEQ-F2997A17** — Custom controls provide equivalent semantics and keyboard behavior.
- [ ] **USEQ-036E12C2** — Audio and video provide required captions, transcripts, and descriptions.
- [ ] **USEQ-C0550F45** — Information is not communicated by color, shape, position, or sound alone.
- [ ] **USEQ-4C2290A7** — Text and meaningful visual elements meet applicable contrast requirements.
- [ ] **USEQ-FAF81EE6** — Content reflows and remains usable under zoom, text enlargement, and text-spacing changes.
- [ ] **USEQ-76EE07DC** — Reduced-motion and other relevant user preferences are respected.
- [ ] **USEQ-AFACCA2E** — Forms have explicit labels, instructions, required-state indications, and accessible format guidance.
- [ ] **USEQ-AD9014F4** — Errors are identified, described, associated with fields, summarized, and announced appropriately.
- [ ] **USEQ-C31A9626** — Dynamic status changes are announced to assistive technologies.
- [ ] **USEQ-5071D1FC** — Authentication does not rely solely on inaccessible puzzles or prohibited cognitive tests.
- [ ] **USEQ-AC783CD7** — Time limits can be extended, disabled, or otherwise meet applicable requirements.
- [ ] **USEQ-55A051B4** — Repeated entry is minimized and help is located consistently.
- [ ] **USEQ-0D6B7B8D** — Page language and text direction are declared correctly.
- [ ] **USEQ-29ED8FF9** — Dialogs, menus, tabs, tooltips, notifications, and live regions follow accessible interaction patterns.
- [ ] **USEQ-2163C314** — Third-party embedded components meet the target or have an accessible alternative.
- [ ] **USEQ-2223809C** — Testing includes representative screen-reader and assistive-technology combinations.
- [ ] **USEQ-A9BB5F5C** — Accessibility defects have owners, severity, remediation dates, and regression coverage.

## Localization and Internationalization

_Consolidated from `quality standards/03-ux-ui-accessibility/02-localization-and-internationalization.md`; 12 non-duplicative controls._

### Universal controls

- [ ] **USEQ-E25EAE4D** — Every supported locale, language, script, time zone, currency, and region is documented.
- [ ] **USEQ-563F7189** — Text expansion and contraction do not break layout or truncate meaning.
- [ ] **USEQ-7AD07D0F** — Right-to-left and mixed-direction content are tested.
- [ ] **USEQ-84C5B4DF** — Unicode normalization and confusable characters are handled where security-sensitive.
- [ ] **USEQ-3860FBC6** — Names, addresses, phone numbers, identifiers, and personal details are not constrained to one cultural format without need.
- [ ] **USEQ-A44F53C5** — Date, time, number, decimal, currency, percentage, unit, and plural formatting are correct.
- [ ] **USEQ-993E9BCA** — Locale-specific sorting, collation, search, tokenization, and case behavior are correct.
- [ ] **USEQ-529B57C9** — Legal, privacy, consent, safety, pricing, and cancellation text is correct for each applicable locale.
- [ ] **USEQ-BD87BE8C** — Translated error, security, recovery, safety, and support content is understandable.
- [ ] **USEQ-1DCEC7EE** — Localization does not alter identifiers, signatures, tokens, URLs, machine data, or security-sensitive canonical forms.
- [ ] **USEQ-84BBB844** — Support coverage matches languages promised to users.
- [ ] **USEQ-08A393B8** — Translation providers and translation memories receive only approved data.

## Human-Centred Design and User Research

_Consolidated from `quality standards/03-ux-ui-accessibility/03-human-centered-design-and-user-research.md`; 15 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-9F00B439** — Understand users, goals, tasks, workflows, physical and social environments, constraints, and assistive strategies.
- [ ] **USEQ-21D170D6** — Include users throughout discovery, design, evaluation, launch, and improvement rather than only at final validation.
- [ ] **USEQ-E64D2BC2** — Represent diverse abilities, literacy, language, culture, devices, connectivity, expertise, and stressful contexts.
- [ ] **USEQ-1D3938D7** — Define usability and quality-in-use goals before selecting interface solutions.
- [ ] **USEQ-01709A04** — Create testable hypotheses and prototypes at the lowest fidelity that can answer the question.
- [ ] **USEQ-067A60A9** — Evaluate complete journeys, transitions, errors, recovery, and cross-channel handoffs.
- [ ] **USEQ-8FA0B608** — Observe real task performance rather than relying only on preference ratings.
- [ ] **USEQ-A33D2471** — Separate researcher facilitation from leading participants toward intended behavior.
- [ ] **USEQ-735FCF5F** — Record participant characteristics, tasks, methods, limitations, and evidence.
- [ ] **USEQ-A507D3B0** — Prioritize findings by user impact, frequency, severity, reach, and recoverability.
- [ ] **USEQ-44B1D1DA** — Verify fixes with users rather than assuming design review is sufficient.
- [ ] **USEQ-D299D8CB** — Ensure research artifacts do not expose participant identity or sensitive context.
- [ ] **USEQ-C5C1050F** — Share findings in forms that influence product, requirements, architecture, support, and operations.
- [ ] **USEQ-BDCFA5BA** — Maintain a research repository with provenance, access control, and expiry rules.
- [ ] **USEQ-D96A88ED** — Revisit user models as the population and product evolve.

## Information Architecture and Content Design

_Consolidated from `quality standards/03-ux-ui-accessibility/04-information-architecture-and-content-design.md`; 15 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-724603EE** — Organize information around user goals and mental models rather than internal organizational structure.
- [ ] **USEQ-AB908F8E** — Use consistent labels for the same concepts and distinct labels for different concepts.
- [ ] **USEQ-107041D9** — Define content hierarchy, navigation paths, relationships, metadata, and search behavior.
- [ ] **USEQ-845A4D53** — Ensure users can determine where they are, what they can do, and how to return.
- [ ] **USEQ-6D09E88C** — Design for recognition and progressive disclosure rather than unnecessary memorization.
- [ ] **USEQ-83265957** — Use plain, precise, respectful language appropriate to user literacy and context.
- [ ] **USEQ-95CE72AD** — Place critical instructions, constraints, costs, risks, and consequences before commitment.
- [ ] **USEQ-4B8463E6** — Write actionable error, empty, loading, success, warning, and recovery content.
- [ ] **USEQ-4921112D** — Avoid jargon, unexplained acronyms, ambiguous pronouns, and culturally narrow assumptions.
- [ ] **USEQ-F0806925** — Ensure headings, link text, labels, instructions, and calls to action are meaningful out of context.
- [ ] **USEQ-9B8E587D** — Keep legal and technical accuracy without making essential information incomprehensible.
- [ ] **USEQ-EDE89DB2** — Design content for translation, text expansion, bidirectionality, speech, and assistive technologies.
- [ ] **USEQ-48449C3C** — Define content ownership, review cadence, effective dates, localization status, and retirement.
- [ ] **USEQ-7EC2A731** — Test navigation and comprehension with representative users.
- [ ] **USEQ-EE6096D0** — Prevent duplicate, contradictory, stale, or orphaned content across product and support channels.

## Interaction and Visual Design

_Consolidated from `quality standards/03-ux-ui-accessibility/05-interaction-and-visual-design.md`; 15 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-6A1ECCAB** — Make controls discoverable and their purpose understandable before activation.
- [ ] **USEQ-F1EEB77C** — Use consistent interaction patterns for equivalent actions and states.
- [ ] **USEQ-7B99D597** — Provide immediate, accurate feedback for user actions and system state changes.
- [ ] **USEQ-3AF57456** — Distinguish available, selected, focused, disabled, pending, destructive, and completed states.
- [ ] **USEQ-54482BE7** — Prevent accidental destructive action through safe defaults, confirmation, preview, undo, or staged commitment.
- [ ] **USEQ-50035160** — Preserve user work across recoverable errors, refreshes, timeouts, and navigation.
- [ ] **USEQ-ECDA8FDD** — Use visual hierarchy to communicate importance without depending on color, position, or shape alone.
- [ ] **USEQ-42B235C1** — Maintain sufficient contrast, legibility, target size, spacing, and focus visibility.
- [ ] **USEQ-FED86988** — Avoid motion, flashing, parallax, and animation that create distraction, nausea, seizure risk, or hidden delay.
- [ ] **USEQ-A04BD2E2** — Ensure layout remains coherent under zoom, reflow, translation, long content, and user font settings.
- [ ] **USEQ-24C6FA97** — Design all normal, empty, loading, partial, stale, offline, error, and recovery states.
- [ ] **USEQ-EC571D80** — Keep system status and long-running progress truthful; do not display false precision.
- [ ] **USEQ-A5BF5409** — Make keyboard, touch, pointer, switch, voice, and assistive interaction behavior coherent.
- [ ] **USEQ-E9AC951D** — Use platform conventions unless deviation has a demonstrated user benefit.
- [ ] **USEQ-13A71C35** — Evaluate visual polish only after correctness, accessibility, comprehension, and task success are satisfied.

## Usability and Task Success

_Consolidated from `quality standards/03-ux-ui-accessibility/06-usability-and-task-success.md`; 15 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-938B15E0** — Define representative tasks, users, contexts, success criteria, and severity thresholds.
- [ ] **USEQ-EC681CC8** — Measure task completion, correctness, time, errors, assistance, confidence, and recovery as appropriate.
- [ ] **USEQ-03C4A44D** — Test first-use learnability and repeated-use efficiency separately.
- [ ] **USEQ-AD024376** — Test infrequent, stressful, high-stakes, and interruption-prone tasks.
- [ ] **USEQ-6CC1B95E** — Identify slips, mistakes, mode errors, ambiguous feedback, and hidden system state.
- [ ] **USEQ-A53E1887** — Reduce unnecessary steps, decisions, data entry, context switching, and repeated information.
- [ ] **USEQ-9A42EBDF** — Provide sensible defaults without concealing material choices or consequences.
- [ ] **USEQ-D4386484** — Support undo, cancel, save, resume, and recovery where users can reasonably need them.
- [ ] **USEQ-21D27DAC** — Prevent users from reaching irreversible invalid states.
- [ ] **USEQ-49142045** — Ensure help is contextual, searchable, consistent, and not a substitute for fixable design problems.
- [ ] **USEQ-EE6BEDAD** — Use severity based on impact, frequency, persistence, and ability to recover.
- [ ] **USEQ-8A50C66F** — Include assistive-technology and low-connectivity contexts in usability evidence.
- [ ] **USEQ-339013AF** — Test end-to-end journeys that cross devices, channels, teams, and external services.
- [ ] **USEQ-1CA9F46F** — Compare observed outcomes with defined quality-in-use objectives.
- [ ] **USEQ-C23FAC50** — Do not treat satisfaction alone as proof of effective or safe use.

## Design Systems and Interface Reuse

_Consolidated from `quality standards/03-ux-ui-accessibility/07-design-systems-and-interface-reuse.md`; 15 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-0F00C6D6** — Define the design system's scope, consumers, ownership, contribution model, and support policy.
- [ ] **USEQ-31BD9A24** — Separate design tokens, primitives, components, patterns, templates, content guidance, and product-specific composition.
- [ ] **USEQ-0F527338** — Provide accessible semantics and interaction behavior by default.
- [ ] **USEQ-F4C6894A** — Define variants intentionally and prevent arbitrary configuration combinations.
- [ ] **USEQ-7C23FD4C** — Keep APIs minimal, coherent, typed or schema-defined, and backward-compatible where promised.
- [ ] **USEQ-1AE96CDB** — Document usage, non-usage, states, examples, accessibility behavior, and migration guidance.
- [ ] **USEQ-1D542805** — Test visual, semantic, behavioral, keyboard, screen-reader, localization, and responsive behavior.
- [ ] **USEQ-A83EFAA8** — Version releases and communicate breaking changes before adoption.
- [ ] **USEQ-DDE326BC** — Provide codemods, migration tools, or clear upgrade paths for material changes.
- [ ] **USEQ-215174DB** — Measure adoption, divergence, defects, accessibility issues, maintenance cost, and user consistency.
- [ ] **USEQ-454D02E1** — Allow justified exceptions without permitting uncontrolled forks.
- [ ] **USEQ-5DBE2F75** — Feed product findings back into shared components only when the pattern is genuinely reusable.
- [ ] **USEQ-FC92CC8B** — Avoid forcing one component to serve incompatible purposes through excessive parameters.
- [ ] **USEQ-EFAC915C** — Deprecate duplicated and unsafe components with a controlled migration plan.
- [ ] **USEQ-04C474D7** — Ensure component reuse does not replace journey-level usability and accessibility testing.

## Responsive, Cross-Device, and Adaptive Design

_Consolidated from `quality standards/03-ux-ui-accessibility/08-responsive-cross-device-and-adaptive-design.md`; 15 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-8E136E1B** — Define supported device, viewport, browser, operating environment, input, and assistive-technology combinations.
- [ ] **USEQ-B05285EB** — Use content and task priority to determine adaptation rather than simply shrinking layouts.
- [ ] **USEQ-5402B181** — Preserve essential content, controls, context, and error recovery at every supported size.
- [ ] **USEQ-64655145** — Avoid horizontal scrolling for ordinary reading and operation except where content semantics require it.
- [ ] **USEQ-B87C9294** — Support zoom, reflow, text spacing, orientation, and user font preferences.
- [ ] **USEQ-614E2300** — Provide equivalent operation across touch, pointer, keyboard, voice, and assistive input.
- [ ] **USEQ-F22716E4** — Avoid hover-only, precision-only, gesture-only, or device-motion-only functionality without alternatives.
- [ ] **USEQ-5F887C8D** — Account for notches, safe areas, browser chrome, virtual keyboards, and dynamic viewport changes.
- [ ] **USEQ-4D92636C** — Test low-memory, low-power, slow-network, high-latency, and intermittent-connectivity conditions.
- [ ] **USEQ-5DD9239E** — Use device capabilities only after permission, feature detection, and graceful fallback.
- [ ] **USEQ-63A74C0A** — Prevent adaptive behavior from changing semantics or hiding user choices unpredictably.
- [ ] **USEQ-24784FD3** — Maintain session, state, and task continuity across orientation and viewport changes.
- [ ] **USEQ-CC9F5F44** — Verify media, tables, diagrams, charts, dialogs, and long forms at constrained sizes.
- [ ] **USEQ-A0040AE7** — Test representative physical devices, not only simulated viewports.
- [ ] **USEQ-B17084E9** — Document unsupported limitations and provide a safe fallback or alternate channel.

## Ethical Design and Dark-Pattern Prevention

_Consolidated from `quality standards/03-ux-ui-accessibility/09-ethical-design-and-dark-pattern-prevention.md`; 15 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-8F11375D** — Present material costs, conditions, data uses, renewal, lock-in, risk, and consequences before commitment.
- [ ] **USEQ-A9660898** — Make refusal, cancellation, opt-out, deletion, and complaint paths as understandable and practicable as acceptance.
- [ ] **USEQ-6FA11062** — Do not use preselection, obstruction, nagging, false scarcity, false urgency, confirm-shaming, or visual interference to manipulate choice.
- [ ] **USEQ-57BE8B37** — Do not disguise advertising, sponsored content, affiliate influence, or paid ranking.
- [ ] **USEQ-F033C102** — Do not create account or purchase consent through inactivity, ambiguity, unrelated bundling, or misleading button hierarchy.
- [ ] **USEQ-E5885EF0** — Use defaults that protect user interests and can be changed without penalty.
- [ ] **USEQ-600EE2F9** — Avoid infinite engagement patterns where they conflict with user wellbeing or stated goals.
- [ ] **USEQ-E2FA47BC** — Allow users to review, correct, and confirm consequential data before submission.
- [ ] **USEQ-1B0A0B04** — Distinguish reversible exploration from binding commitment.
- [ ] **USEQ-769B1307** — Do not exploit children, distress, disability, low literacy, financial vulnerability, or cognitive bias.
- [ ] **USEQ-7F004109** — Test whether representative users understand the choice and its consequences.
- [ ] **USEQ-250AC804** — Log and review complaints, reversals, accidental purchases, consent withdrawal, and cancellation failure as design-quality signals.
- [ ] **USEQ-92CBFB7D** — Require independent review for monetization or growth experiments that alter user autonomy.
- [ ] **USEQ-C5413AA8** — Prevent incentive targets from rewarding manipulative conversion at the expense of long-term trust.
- [ ] **USEQ-06E018CC** — Remove deceptive patterns even when they improve short-term metrics.

## Browser, Client, and Usability Behavior

_Consolidated from `quality standards/06-frontend/01-browser-client-usability.md`; 21 non-duplicative controls._

### Universal controls

- [ ] **USEQ-448A4605** — A supported browser, device, operating-system, and viewport matrix is documented.
- [ ] **USEQ-01713736** — Critical journeys work on every supported browser and device class.
- [ ] **USEQ-7D7FF031** — Layout works at supported viewport sizes and orientations.
- [ ] **USEQ-32F92A25** — Zoom, text enlargement, and user font settings do not hide essential functionality.
- [ ] **USEQ-EA8F2704** — Touch, mouse, keyboard, speech, switch, and assistive input work appropriately where supported.
- [ ] **USEQ-852BA49B** — Loading, empty, success, warning, error, partial, timeout, and offline states are designed and tested.
- [ ] **USEQ-0ADCC52B** — Validation is understandable and associated with the affected input.
- [ ] **USEQ-1FE2D09F** — Back, forward, refresh, bookmarks, history, and deep links behave correctly.
- [ ] **USEQ-174891AC** — Client-side caching and browser storage do not expose another user's data or retain unnecessary sensitive data.
- [ ] **USEQ-39524F1B** — Failed-script and third-party-script failure behavior degrades acceptably.
- [ ] **USEQ-00520E0D** — No unhandled exceptions or material console errors remain in critical paths.
- [ ] **USEQ-41B14913** — Source maps, diagnostic information, and stack traces are restricted appropriately.
- [ ] **USEQ-C1C3DABE** — Error pages do not expose internal details.
- [ ] **USEQ-5C64BD44** — Download, print, clipboard, file selection, and upload behavior work where supported.
- [ ] **USEQ-AE01960E** — Long text, translations, unusual names, bidirectional text, and user-supplied content are handled safely.
- [ ] **USEQ-7D3A3905** — Destructive operations require appropriate confirmation, preview, or reversal.
- [ ] **USEQ-A4BBF142** — Navigation and controls are consistent; users are not trapped in dead ends.
- [ ] **USEQ-1D48139F** — Session-timeout warnings and recovery preserve work where appropriate.
- [ ] **USEQ-9F07016C** — Autosave is transparent and reliable.
- [ ] **USEQ-9E8C9625** — Analytics and experimentation code do not alter critical behavior unexpectedly.
- [ ] **USEQ-D6F52943** — Dark patterns, deceptive defaults, hidden costs, and misleading urgency are absent.

## Public Content, Search Engines, Metadata, and Indexing

_Consolidated from `quality standards/06-frontend/02-public-content-search-and-indexing.md`; 12 non-duplicative controls._

### Universal controls

- [ ] **USEQ-B5E1906E** — Public and private pages have correct indexing and crawler rules.
- [ ] **USEQ-9A31C0F5** — Authentication and authorization, not crawler directives, protect private data.
- [ ] **USEQ-A7DD4F29** — Redirects preserve intended security, privacy, and indexing behavior.
- [ ] **USEQ-0D005E7B** — Status codes are correct; soft-error behavior is avoided.
- [ ] **USEQ-A608556D** — Sitemaps are accurate and do not expose private, preview, or temporary URLs.
- [ ] **USEQ-8E2DD59D** — Metadata, structured data, social previews, feeds, and snippets contain no sensitive information.
- [ ] **USEQ-242BBE23** — Structured data matches visible content and is not misleading.
- [ ] **USEQ-D8DD5868** — Deleted, private, expired, or restricted content is removed from caches and indexing workflows.
- [ ] **USEQ-F9D17ABE** — Staging, preview, test, and administrative environments are not unintentionally indexed.
- [ ] **USEQ-EF32EE87** — URL generation avoids unbounded duplicate spaces and crawler traps.
- [ ] **USEQ-27ABC548** — User-generated content cannot manipulate site-wide metadata, canonical links, redirects, or scripts.
- [ ] **USEQ-F12635E9** — Search-result caching and public snapshots respect deletion, privacy, and access-control changes.

## Offline and Progressive Web Capabilities

_Consolidated from `quality standards/06-frontend/03-offline-and-progressive-web-capabilities.md`; 9 non-duplicative controls._

### Universal controls

- [ ] **USEQ-D802C2BE** — Service-worker and offline-cache scope is restricted.
- [ ] **USEQ-BA85F05E** — Cache versioning, update, invalidation, and rollback are reliable.
- [ ] **USEQ-2AA3EF40** — Updates do not leave clients permanently on vulnerable or incompatible versions.
- [ ] **USEQ-4C15FCB3** — Offline changes have clear conflict, ordering, merge, retry, and synchronization behavior.
- [ ] **USEQ-CB659F6B** — Logout, account deletion, tenant switch, and permission change clear or invalidate relevant offline data.
- [ ] **USEQ-2013C300** — Shared-device and lost-device risks are addressed.
- [ ] **USEQ-C2A76E22** — Background synchronization has bounded network, battery, storage, and compute use.
- [ ] **USEQ-C63F668E** — Offline status is communicated accurately so users do not believe unsynchronized actions are complete.
- [ ] **USEQ-03C6755E** — Cached application shells cannot expose data or functionality from a previous user or tenant.

## Frontend Architecture, Components, and State

_Consolidated from `quality standards/06-frontend/04-frontend-architecture-components-and-state.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-CBE25307** — Separate domain behavior, server data, local presentation state, navigation state, and transient interaction state.
- [ ] **USEQ-727A7AEB** — Keep the server or authoritative service as source of truth for security-sensitive and shared business state.
- [ ] **USEQ-F9D875BA** — Use unidirectional or otherwise explicit data flow that makes updates traceable.
- [ ] **USEQ-3448CF6A** — Keep components cohesive around a clear responsibility and accessible contract.
- [ ] **USEQ-0E07913F** — Distinguish reusable primitives from product-specific compositions.
- [ ] **USEQ-BFADCD19** — Avoid components with excessive flags, hidden modes, or unrelated variants.
- [ ] **USEQ-3F4F0E45** — Keep rendering free from unintended side effects.
- [ ] **USEQ-84AAC34F** — Centralize effects only when ownership and lifecycle are clearer, not merely by convention.
- [ ] **USEQ-030328F1** — Cancel or ignore stale requests when navigation, identity, query, or component lifetime changes.
- [ ] **USEQ-44C08441** — Prevent race conditions between optimistic updates, retries, invalidation, and server responses.
- [ ] **USEQ-B2DF5A8F** — Model loading, empty, stale, partial, error, offline, permission, and completed states explicitly.
- [ ] **USEQ-436F8255** — Preserve focus, announcement, scroll, and user input through dynamic updates.
- [ ] **USEQ-1BA69E71** — Keep routing, deep links, refresh, back/forward, and multi-tab behavior consistent.
- [ ] **USEQ-D775D7C6** — Avoid storing sensitive data or reusable credentials in client-accessible persistence without necessity and threat analysis.
- [ ] **USEQ-60F112DD** — Define hydration, cache, persistence, version migration, logout, and reset behavior.
- [ ] **USEQ-438D73DE** — Make state changes diagnosable without exposing personal data.
- [ ] **USEQ-84E85608** — Use error boundaries or equivalent containment so one component failure does not destroy unrelated work.
- [ ] **USEQ-1C46ED62** — Control bundle and runtime dependency boundaries.
- [ ] **USEQ-9F675947** — Test components in isolation and complete journeys in representative browsers.
- [ ] **USEQ-82B9D9C4** — Prevent client-side checks from becoming trusted authorization or validation boundaries.

## Forms, User Input, and Client Validation

_Consolidated from `quality standards/06-frontend/05-forms-user-input-and-client-validation.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-F4415492** — Use native semantic controls when they satisfy the required interaction.
- [ ] **USEQ-D2FDC656** — Provide visible programmatically associated labels and instructions.
- [ ] **USEQ-9AF004CC** — Identify required, optional, format, unit, range, and privacy implications before entry.
- [ ] **USEQ-DFBBBC98** — Use appropriate input modes and autocomplete semantics without exposing sensitive data.
- [ ] **USEQ-367C2E65** — Accept valid international names, addresses, identifiers, numbers, and characters unless a real constraint exists.
- [ ] **USEQ-10E68DD7** — Validate on the server as the authority and use client validation for timely assistance.
- [ ] **USEQ-9D556409** — Avoid rejecting input while the user is still composing unless the rule is certain and helpful.
- [ ] **USEQ-D4865C00** — Show errors near the field and in an accessible summary for multi-field failure.
- [ ] **USEQ-339AADF6** — Preserve valid input after recoverable failure.
- [ ] **USEQ-3E47809F** — Move or announce focus appropriately without trapping users.
- [ ] **USEQ-E598349A** — Distinguish validation, conflict, authorization, timeout, dependency, and submission errors.
- [ ] **USEQ-68077DB5** — Prevent duplicate submission through state, idempotency, and clear progress feedback.
- [ ] **USEQ-D70CBF17** — Allow review and correction before consequential submission.
- [ ] **USEQ-CD711588** — Provide undo or cancellation where consequences are reversible.
- [ ] **USEQ-47950324** — Do not disable paste, password managers, or assistive input without an overriding evidenced need.
- [ ] **USEQ-F2975590** — Protect secrets and sensitive values from display, browser history, analytics, logging, and unintended autocomplete.
- [ ] **USEQ-124ABB79** — Handle long forms through grouping, progress, save, resume, and timeout recovery where appropriate.
- [ ] **USEQ-00271F96** — Maintain keyboard order and visible focus.
- [ ] **USEQ-F6E4CD10** — Test autofill, browser navigation, refresh, multi-tab, slow network, offline, and server-side error behavior.
- [ ] **USEQ-5E3BF420** — Use anti-automation controls that do not create inaccessible or discriminatory barriers.

## Frontend Performance and Runtime Efficiency

_Consolidated from `quality standards/06-frontend/06-frontend-performance-and-runtime-efficiency.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-378CA6B2** — Define user-centric budgets for loading, interaction, visual stability, memory, payload, and energy on representative devices and networks.
- [ ] **USEQ-58AA81DE** — Measure real-user distributions segmented by relevant device, network, geography, and journey.
- [ ] **USEQ-F7DFB348** — Prioritize critical content and interaction before nonessential functionality.
- [ ] **USEQ-DE1F929C** — Minimize initial code, data, style, font, image, and third-party payloads.
- [ ] **USEQ-B0511771** — Split and load resources according to actual navigation and interaction needs.
- [ ] **USEQ-49FAB56F** — Compress, resize, encode, and cache media appropriately without harming accessibility or quality.
- [ ] **USEQ-AF52281A** — Avoid long main-thread tasks and excessive layout, paint, hydration, or rendering work.
- [ ] **USEQ-CB6F9C32** — Prevent layout shifts through known dimensions, stable placeholders, and deliberate dynamic insertion.
- [ ] **USEQ-F4263DF1** — Virtualize or paginate large collections without breaking semantics, findability, focus, or assistive technology.
- [ ] **USEQ-A23E20A3** — Debounce, throttle, batch, or defer work only when correctness and feedback remain intact.
- [ ] **USEQ-1AC4BD88** — Avoid unnecessary polling and duplicated requests.
- [ ] **USEQ-AC3AAF8A** — Cancel obsolete work and release event listeners, observers, timers, workers, and large retained objects.
- [ ] **USEQ-E8BDEA08** — Measure and constrain third-party script CPU, network, privacy, failure, and availability impact.
- [ ] **USEQ-42F3193D** — Use caching with correct identity, authorization, locale, and invalidation dimensions.
- [ ] **USEQ-5218074C** — Test cold, warm, empty-cache, returning-user, low-memory, slow-network, and degraded-service states.
- [ ] **USEQ-1B19048B** — Prevent performance optimizations from serving stale or cross-user data.
- [ ] **USEQ-4ACFF818** — Set automated regression budgets and investigate statistically meaningful deterioration.
- [ ] **USEQ-ED992502** — Profile before optimizing and verify end-to-end benefit after change.
- [ ] **USEQ-981ACAE8** — Keep telemetry overhead bounded.
- [ ] **USEQ-74C20BFC** — Treat responsiveness and accessibility as linked quality outcomes.

## Frontend Testing

_Consolidated from `quality standards/06-frontend/07-frontend-testing.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-EC885C79** — Test behavior through user-observable semantics rather than internal implementation structure.
- [ ] **USEQ-1D37E90B** — Cover component states, variants, boundaries, permissions, locales, themes, and responsive sizes.
- [ ] **USEQ-E1BA397D** — Test keyboard, focus, semantics, announcements, zoom, reflow, and representative assistive technology.
- [ ] **USEQ-ED0206BA** — Test routing, refresh, deep links, back, forward, multiple tabs, and session changes.
- [ ] **USEQ-6B3537D2** — Simulate slow, failed, duplicated, reordered, stale, and canceled network responses.
- [ ] **USEQ-C27619F9** — Test optimistic updates, conflict, rollback, retry, and cache invalidation.
- [ ] **USEQ-9D6AF433** — Verify forms with valid, invalid, partial, international, long, pasted, autofilled, and server-rejected input.
- [ ] **USEQ-B2DB35B9** — Use visual regression for meaningful layout and styling changes while reviewing semantic and behavioral quality separately.
- [ ] **USEQ-0030FC91** — Run critical journeys in every supported browser and representative physical device class.
- [ ] **USEQ-1F409E8C** — Avoid fixed sleeps; wait for observable state or controlled time.
- [ ] **USEQ-625A2E05** — Keep test data deterministic and isolated.
- [ ] **USEQ-5706BF1A** — Do not over-mock browser or service behavior that contract tests can validate more faithfully.
- [ ] **USEQ-D6D5DEBA** — Detect unhandled errors, rejected promises, console failures, network failures, and accessibility violations.
- [ ] **USEQ-607CE5B2** — Test without relying solely on pointer interaction.
- [ ] **USEQ-2892F8F1** — Verify source maps, production optimization, content security, caching, and client telemetry in a production-like build.
- [ ] **USEQ-B4573403** — Test third-party failure and blocked-script behavior.
- [ ] **USEQ-8719391E** — Measure flakiness and quarantine only with owner, risk, and expiry.
- [ ] **USEQ-E1394AF2** — Preserve screenshots, traces, videos, network logs, and accessibility trees for material failures as appropriate.
- [ ] **USEQ-A0BC5AEB** — Run production-safe synthetic checks for critical client journeys.
- [ ] **USEQ-DC069975** — Convert escaped frontend defects into durable regression tests.

## Frontend Security and Privacy

_Consolidated from `quality standards/06-frontend/08-frontend-security-and-privacy.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-532876B7** — Treat all client code, state, validation, and feature visibility as attacker-controlled.
- [ ] **USEQ-BDDB106C** — Enforce authorization, validation, pricing, limits, and business invariants on trusted servers.
- [ ] **USEQ-6236CA98** — Encode or sanitize untrusted content for the exact browser context.
- [ ] **USEQ-DA8CCC7C** — Avoid dangerous dynamic code, markup, URL, style, and template construction.
- [ ] **USEQ-7374C24C** — Use a restrictive content security policy appropriate to the application.
- [ ] **USEQ-E9140CCA** — Inventory, minimize, pin, and monitor third-party scripts and browser dependencies.
- [ ] **USEQ-180B6B5A** — Load sensitive functionality and data only after authorization, not merely hide it visually.
- [ ] **USEQ-1D83B905** — Keep reusable credentials, secrets, sensitive personal data, and security answers out of URLs and unnecessary client storage.
- [ ] **USEQ-585C301E** — Scope and secure cookies according to session needs.
- [ ] **USEQ-1EBC2C98** — Clear appropriate state, caches, storage, workers, and memory on logout, account switch, or deletion.
- [ ] **USEQ-472165AA** — Prevent cross-user data leakage through browser cache, shared devices, multi-tab state, and service workers.
- [ ] **USEQ-137DBBC0** — Restrict cross-origin communication, frames, openers, messaging, and embedded content to validated origins and schemas.
- [ ] **USEQ-8E1EED4B** — Validate redirect, callback, download, and navigation destinations.
- [ ] **USEQ-9FD78746** — Protect state-changing browser requests against cross-site request forgery.
- [ ] **USEQ-64F30E91** — Prevent sensitive data from analytics, error reporting, session replay, clipboard, referrers, previews, and screenshots where control is possible.
- [ ] **USEQ-67A5F5CE** — Request browser and device permissions contextually and explain their purpose.
- [ ] **USEQ-CE2A9BF5** — Provide functional degradation when nonessential tracking or third-party content is blocked.
- [ ] **USEQ-10A564C2** — Make privacy choices accessible and as easy to withdraw as grant.
- [ ] **USEQ-F102F8DA** — Test browser-extension, injected-content, autofill, and shared-device risks according to threat model.
- [ ] **USEQ-E3FA430B** — Monitor client security policy violations without collecting excessive personal data.

## Web Quality, SEO, Accessibility, and Content Master Checklist

_Consolidated from `gap supplement/01-web-quality-seo-accessibility-and-content.md`; 254 non-duplicative controls._

### Expanded gap-closure controls

#### Scope, ownership, and quality governance

- [ ] **USEQ-F62F0FFB** — Apply the SEO controls whenever public discovery through general, vertical, enterprise, or AI-assisted search is an intended acquisition or support channel.
- [ ] **USEQ-3E7F8909** — Apply the accessibility controls to every human-facing interface, including authenticated areas, support tools, administrative tools, embedded experiences, generated documents, and transactional communications.
- [ ] **USEQ-EF73B7C0** — Define the supported browser, device, assistive-technology, crawler, rendering, locale, and network conditions for each critical journey.
- [ ] **USEQ-E38F347D** — Maintain an inventory of public origins, hostnames, URL spaces, page types, templates, feeds, media libraries, generated documents, and embedded third-party surfaces.
- [ ] **USEQ-7BD51AB1** — Assign accountable product, engineering, content, SEO, accessibility, privacy, security, and operations owners for the surfaces in scope.
- [ ] **USEQ-557C8048** — Define measurable goals for discoverability, task success, content accuracy, accessibility, performance, privacy, and resilience without treating ranking position as the sole success measure.
- [ ] **USEQ-116267A5** — Map each page type to its intended audience, search intent, user task, business purpose, indexability, canonical policy, retention policy, and owner.
- [ ] **USEQ-7F036200** — Include SEO, accessibility, semantic markup, metadata, performance, and content checks in the definition of done for every affected change.
- [ ] **USEQ-67A6F5CB** — Require design-system components and content templates to meet accessibility and web-platform requirements before broad reuse.
- [ ] **USEQ-7E41BD98** — Prevent an accessibility overlay, toolbar, or automated remediation product from being treated as a substitute for accessible design, implementation, and testing.
- [ ] **USEQ-4012FF4F** — Require third-party widgets, consent tools, payment components, maps, chat, advertising, analytics, and embeds to meet the same release criteria or provide an equivalent accessible and privacy-preserving alternative.
- [ ] **USEQ-419F3CF7** — Document every intentional deviation from current search-engine guidance, web standards, or accessibility targets together with the expected impact and monitoring plan.
- [ ] **USEQ-F06C0099** — Maintain a deprecation and retirement process for obsolete pages, templates, components, redirects, structured-data types, and public endpoints.
- [ ] **USEQ-8F7990CC** — Include SEO and accessibility impact in architecture decisions, design reviews, procurement, supplier review, and incident postmortems.
- [ ] **USEQ-89404A36** — Test representative production content rather than only empty templates or idealized fixture data.

#### People-first content and editorial quality

- [ ] **USEQ-B1AFF4A1** — Create content primarily to satisfy an identifiable user need rather than to manipulate search visibility.
- [ ] **USEQ-A83B6A87** — Ensure every indexable page provides original utility, authoritative information, a useful transaction, or a clear navigational purpose.
- [ ] **USEQ-51900BB2** — State the intended audience, question, task, and desired outcome before producing or materially revising content.
- [ ] **USEQ-78C6933C** — Verify factual claims against reliable sources and retain source provenance for high-impact content.
- [ ] **USEQ-55EBFF21** — Identify authors, reviewers, organizations, credentials, conflicts, or editorial responsibility where trust and accountability depend on them.
- [ ] **USEQ-412E5759** — Display publication, review, and material-update dates when freshness affects user decisions.
- [ ] **USEQ-8CF64C13** — Define review cadences based on volatility, harm from inaccuracy, legal obligations, and user dependence.
- [ ] **USEQ-1C249F21** — Correct or withdraw inaccurate content promptly and preserve an appropriate correction history for consequential material.
- [ ] **USEQ-F08863AE** — Remove, consolidate, redirect, noindex, or clearly archive content that is obsolete, duplicative, unsupported, or no longer useful.
- [ ] **USEQ-B2E8D30C** — Ensure page titles, headings, summaries, labels, and calls to action accurately describe the visible content and outcome.
- [ ] **USEQ-0DFBDFD4** — Avoid keyword stuffing, hidden text, doorway pages, cloaking, scraped content, deceptive redirects, link schemes, and other manipulative practices.
- [ ] **USEQ-89F2E41C** — Do not generate large volumes of substantially similar or low-value pages merely to capture query variations.
- [ ] **USEQ-2CB91942** — Subject AI-assisted or automatically generated content to accountable human review, source verification, originality checks, legal review where applicable, and the same quality thresholds as human-authored content.
- [ ] **USEQ-0A702F02** — Disclose synthetic, sponsored, affiliate, advertorial, or materially automated content when users need that information to judge trust or when policy or law requires it.
- [ ] **USEQ-3692B38E** — Separate advertising, sponsored placements, and affiliate relationships clearly from editorial content.
- [ ] **USEQ-E09BE5FD** — Prevent advertisements, interstitials, consent prompts, or promotional elements from obscuring the primary content or creating deceptive interactions.
- [ ] **USEQ-1D35FB17** — Use plain language appropriate to the audience while preserving necessary technical, legal, or safety precision.
- [ ] **USEQ-42217D89** — Provide definitions, examples, prerequisites, limitations, and next steps where users could otherwise misunderstand or fail the task.
- [ ] **USEQ-7E4355CC** — Write link text, button labels, headings, instructions, and error messages that remain understandable out of visual context.
- [ ] **USEQ-77285B0F** — Ensure translated and localized content receives linguistic and subject-matter review rather than relying blindly on literal machine translation.
- [ ] **USEQ-3C25D0F9** — Prevent internal search results, faceted combinations, autogenerated tags, empty categories, and empty or scaffolding pages from becoming indexable low-value inventory by default.
- [ ] **USEQ-2F5B5FE8** — Maintain an editorial style, terminology, source, citation, correction, accessibility, and archival policy.
- [ ] **USEQ-9FDB2D05** — Measure content success using task completion, qualified engagement, satisfaction, support reduction, conversion quality, and harm indicators, not raw page views alone.

#### Information architecture, URLs, and internal linking

- [ ] **USEQ-998C9BFF** — Design the information architecture around user tasks, concepts, and stable domain entities rather than the current organizational chart or implementation details.
- [ ] **USEQ-EEA2F3BA** — Use stable, unique, absolute URLs for resources intended to be referenced or indexed.
- [ ] **USEQ-D4F93C36** — Keep public URLs understandable, durable, and free of unnecessary session identifiers, secrets, personal data, volatile state, or tracking parameters.
- [ ] **USEQ-615CFF9A** — Define one normalization policy for host, scheme, port, path case, trailing separators, encoded characters, query ordering, and default documents.
- [ ] **USEQ-E55B31F7** — Ensure every indexable page is reachable through ordinary crawlable links from an appropriate navigational path unless deliberate isolation is documented.
- [ ] **USEQ-FADBC6AC** — Use native link semantics for navigation and ensure destination URLs exist in the rendered markup.
- [ ] **USEQ-84DD0A05** — Use descriptive anchor text that communicates destination and purpose without manipulative repetition.
- [ ] **USEQ-0608EB2A** — Prevent orphan pages by reconciling the content inventory, sitemap inventory, internal link graph, and indexed inventory.
- [ ] **USEQ-0C5EAD50** — Keep important pages within a reasonable navigational depth and make hierarchy understandable through navigation, headings, and breadcrumbs where useful.
- [ ] **USEQ-C40E8A5B** — Ensure breadcrumbs represent the actual information hierarchy and use valid structured data when published.
- [ ] **USEQ-0AE2AC0A** — Control faceted navigation, filters, sorting, pagination, calendars, session state, and user-generated parameters so they do not create unbounded URL spaces.
- [ ] **USEQ-60E4D076** — Define which parameter combinations are indexable, canonical, blocked from crawling, or excluded from generation.
- [ ] **USEQ-F548A11F** — Preserve meaningful state in URLs when users need bookmarking, sharing, recovery, or browser navigation.
- [ ] **USEQ-FFA8759D** — Avoid fragment-only navigation for distinct resources that require independent discovery unless the rendering and indexing behavior is deliberately verified.
- [ ] **USEQ-90FEF1BF** — Update internal links to final destinations instead of relying indefinitely on redirect chains.
- [ ] **USEQ-13B4D3B9** — Check for broken links, circular redirects, mixed schemes, malformed URLs, and links to retired or unauthorized content on every release.
- [ ] **USEQ-AEFD2DB5** — Provide useful navigation and recovery from not-found, expired, unauthorized, and removed resources without returning misleading success responses.
- [ ] **USEQ-8159A214** — Ensure public taxonomy, tags, categories, and labels have owners, definitions, lifecycle rules, and minimum-content thresholds.

#### Crawling, discovery, and server behavior

- [ ] **USEQ-C67E7B13** — Use crawler-control files only to manage crawling; never rely on them to protect confidential or unauthorized resources.
- [ ] **USEQ-9E5F39EC** — Validate crawler-control syntax, host scope, path matching, case behavior, encoding, and production deployment.
- [ ] **USEQ-7F6BBCA4** — Ensure resources required to understand or render public content are not unintentionally blocked from permitted crawlers.
- [ ] **USEQ-21F93E8E** — Use page-level indexing directives consistently and verify that crawlers can access a page when they must observe a noindex directive.
- [ ] **USEQ-34BA8053** — Generate sitemaps from authoritative canonical inventory rather than from every routable URL.
- [ ] **USEQ-9BC1EE8B** — Include only permitted, canonical, successful, indexable URLs in sitemaps.
- [ ] **USEQ-AFA12A76** — Use accurate modification timestamps only when content changed materially.
- [ ] **USEQ-4683497B** — Split large sitemaps predictably, publish an index when required, and monitor submission and processing errors.
- [ ] **USEQ-89C938D9** — Provide discoverable feeds or specialized sitemaps for video, image, news, or other vertical content when those channels are in scope.
- [ ] **USEQ-28004D24** — Return protocol-correct status codes for success, redirect, client error, removal, rate limiting, temporary failure, and permanent failure.
- [ ] **USEQ-5A386FAB** — Avoid soft not-found responses, infinite redirect loops, excessive redirect chains, and redirects to irrelevant destinations.
- [ ] **USEQ-DC970915** — Choose permanent versus temporary redirects according to the intended resource transition and test method preservation where non-read operations are involved.
- [ ] **USEQ-2954B401** — Use a gone response or equivalent explicit removal signal when permanent deletion should be communicated and no replacement exists.
- [ ] **USEQ-0F6E8B12** — Ensure maintenance, overload, and dependency failures do not return apparently successful indexable pages containing error content.
- [ ] **USEQ-CE7C4C81** — Bound crawlable search, calendar, pagination, sort, filter, and autogenerated spaces to prevent crawler traps and resource exhaustion.
- [ ] **USEQ-3143FA2D** — Rate-limit abusive automation without blocking legitimate accessibility tools, customers, or approved crawlers indiscriminately.
- [ ] **USEQ-54FA437D** — Monitor crawler traffic, response codes, latency, bytes, rendering resources, and abnormal parameter growth from server logs.
- [ ] **USEQ-32244778** — Ensure canonical content remains available and sufficiently stable during crawl spikes, launches, migrations, and incident recovery.
- [ ] **USEQ-5E65D991** — Document handling for unsupported, malicious, unidentified, or spoofed crawlers rather than trusting user-agent strings alone.

#### Rendering and JavaScript-dependent content

- [ ] **USEQ-08CC7F5F** — Ensure critical content, metadata, links, and controls are present in or reliably produced into the document structure that supported user agents and permitted crawlers can process.
- [ ] **USEQ-D623C1A7** — Do not require scrolling, clicking, consent to unrelated processing, authentication, geolocation, or other interaction merely to expose public primary content to a crawler.
- [ ] **USEQ-6B22FD39** — Verify rendered output, not only server source, for every public template that depends on client execution.
- [ ] **USEQ-ED6F5E59** — Detect hydration, rendering, routing, chunk-loading, and script-execution failures and provide a usable recovery path.
- [ ] **USEQ-3AD4C9C0** — Ensure server-rendered and client-rendered states do not contradict one another in content, metadata, canonicals, language, or index directives.
- [ ] **USEQ-CE7DF0F9** — Keep unique titles, descriptions, canonical links, language annotations, and structured data correct after client navigation.
- [ ] **USEQ-BCDC2782** — Use history and navigation APIs without creating duplicate URLs, broken back/forward behavior, or inaccessible focus management.
- [ ] **USEQ-52EAA91D** — Load below-the-fold media lazily without withholding in-viewport primary content or relying on unsupported interaction triggers.
- [ ] **USEQ-DF700EBF** — Reserve layout space for asynchronous content to prevent disruptive movement.
- [ ] **USEQ-C1A02873** — Provide meaningful content and core task completion when nonessential scripts, analytics, advertising, experimentation, or third-party widgets fail.
- [ ] **USEQ-75DA91D2** — Ensure script bundles, styles, fonts, images, and API endpoints needed for rendering are accessible under the intended crawler and browser policies.
- [ ] **USEQ-4C4530E0** — Do not serve materially different content to crawlers and users except for documented, legitimate adaptation that preserves meaning and access.
- [ ] **USEQ-22E12A1A** — Test rendering under slow networks, constrained devices, blocked third parties, disabled storage, and partial script failure.
- [ ] **USEQ-34B8A477** — Prevent client routing from returning a generic success shell for resources that do not exist or are unauthorized.

#### Indexing, canonicalization, and duplicate control

- [ ] **USEQ-519BAA2E** — Define one canonical resource for every set of substantially duplicate public URLs.
- [ ] **USEQ-2461DD10** — Emit self-referential canonical links on canonical pages where the policy uses canonical markup.
- [ ] **USEQ-F0610D78** — Ensure canonical targets are permitted, successful, indexable, equivalent in content and language, and not redirected through avoidable chains.
- [ ] **USEQ-4B5311D6** — Do not use canonical markup as a substitute for redirects when a resource has actually moved.
- [ ] **USEQ-CDFD96B2** — Do not combine contradictory canonical, redirect, noindex, alternate-language, sitemap, and internal-link signals.
- [ ] **USEQ-974E39ED** — Control duplicates created by tracking parameters, print views, sorting, filtering, casing, protocols, hosts, mobile variants, previews, and syndication.
- [ ] **USEQ-F6B9FCB2** — Use cross-domain canonicalization only with verified ownership, equivalent content, and an explicit syndication policy.
- [ ] **USEQ-D58940DA** — Define pagination and infinite-scroll behavior so each accessible state remains navigable and content is not silently undiscoverable.
- [ ] **USEQ-0C59EF9E** — Prevent preview, staging, test, administrative, customer-specific, personalized, and private URLs from entering public indexes.
- [ ] **USEQ-35B07DAE** — Ensure access-control changes, deletion, legal removal, privacy requests, and embargoes propagate to public caches, search indexes, sitemaps, feeds, previews, and archives under organizational control.
- [ ] **USEQ-0903D47F** — Audit the indexed inventory against the intended canonical inventory and investigate unexplained additions or omissions.
- [ ] **USEQ-3B04ADCB** — Monitor duplicate-title, duplicate-description, near-duplicate-content, canonical-conflict, and discovered-not-indexed patterns by template and cause.

#### Metadata, snippets, social previews, and structured data

- [ ] **USEQ-1DFBD521** — Provide a unique, concise, descriptive title for every important indexable page.
- [ ] **USEQ-E868688C** — Provide useful summaries when they improve discovery, while accepting that search systems may choose different snippets.
- [ ] **USEQ-B0BDE9F7** — Keep titles and summaries aligned with the visible page and avoid exaggerated, misleading, or repetitive wording.
- [ ] **USEQ-95D3ED85** — Prevent secrets, personal data, internal identifiers, experiment details, or restricted content from entering metadata, previews, feeds, or structured data.
- [ ] **USEQ-C512162B** — Use heading structure to express document hierarchy rather than to style text or repeat keywords.
- [ ] **USEQ-B19E7322** — Define social-preview titles, descriptions, images, locale, type, and canonical URL when social sharing is in scope.
- [ ] **USEQ-37B79059** — Ensure social-preview images have appropriate rights, safe cropping, meaningful context, and no sensitive information.
- [ ] **USEQ-5E0FD1C7** — Use structured data only for content actually present and accessible on the page.
- [ ] **USEQ-F7D8AE42** — Meet all required and recommended properties for the selected structured-data type where applicable.
- [ ] **USEQ-5479368E** — Use stable identifiers and consistent entity relationships across pages when describing the same organization, person, product, event, article, or other entity.
- [ ] **USEQ-4F3773E5** — Ensure ratings, prices, availability, dates, authorship, offers, and other factual properties are current and not fabricated or selectively misleading.
- [ ] **USEQ-03E02192** — Validate structured data during development, in deployed output, and after client rendering.
- [ ] **USEQ-F191E6BD** — Monitor enhancement reports, parsing errors, policy violations, and eligibility changes after release.
- [ ] **USEQ-D1E537DD** — Remove or update markup immediately when visible content, eligibility, ownership, or policy changes.
- [ ] **USEQ-ABF5EFE2** — Do not mark up hidden, irrelevant, deceptive, prohibited, or user-inaccessible content.
- [ ] **USEQ-4A45244D** — Treat structured-data eligibility as conditional; do not promise a rich result or depend on one for product correctness.

#### International and multilingual discovery

- [ ] **USEQ-44B42594** — Use distinct, stable URLs for materially different language or regional versions intended for independent discovery.
- [ ] **USEQ-4098BB41** — Declare document language and text direction correctly for every version and for mixed-language passages where needed.
- [ ] **USEQ-6A74A758** — Use alternate-language annotations with valid language or region codes, reciprocal relationships, canonical consistency, and an intentional default where appropriate.
- [ ] **USEQ-2E93154B** — Do not canonicalize translated pages to a different language when each translation is intended to be indexed.
- [ ] **USEQ-A94D316D** — Avoid automatic locale redirects that prevent users or crawlers from reaching a chosen version.
- [ ] **USEQ-BD8807AC** — Provide an accessible language and region selector with persistent user choice and clear labels.
- [ ] **USEQ-E76EC703** — Keep content, navigation, metadata, canonicals, structured data, currency, units, dates, addresses, legal text, and support paths internally consistent for each locale.
- [ ] **USEQ-07FEB1DB** — Review translated titles, snippets, link text, headings, error messages, and structured data for natural language and local search intent.
- [ ] **USEQ-DFCD6B5C** — Prevent mixed-locale templates, untranslated fragments, broken bidirectional text, and duplicated machine translations from becoming indexable.
- [ ] **USEQ-8CBD9069** — Handle unavailable translations explicitly rather than silently returning unrelated or partially translated content.
- [ ] **USEQ-59B682D8** — Map regional variants to actual differences in product, law, currency, shipping, availability, or audience rather than creating artificial duplicate pages.
- [ ] **USEQ-9C2B42C3** — Test language negotiation, caching, redirects, alternate annotations, and shared URLs without allowing one locale to poison another locale's cache.

#### Media, documents, feeds, and downloadable content

- [ ] **USEQ-64BE37EA** — Give important images descriptive filenames, surrounding context, dimensions, and appropriate alternative text without keyword stuffing.
- [ ] **USEQ-52AF70C5** — Use responsive media sources and formats while preserving quality, accessibility, rights, and indexability.
- [ ] **USEQ-F151E393** — Provide captions for prerecorded synchronized media and live captions where required by the target standard.
- [ ] **USEQ-12BD5F1B** — Provide transcripts, audio description, sign-language interpretation, or equivalent alternatives when required by content and audience needs.
- [ ] **USEQ-C93BC042** — Publish video on a stable watch page with a unique title, description, thumbnail, duration, availability, and playable primary media when video discovery is intended.
- [ ] **USEQ-06CA04D6** — Ensure media thumbnails accurately represent the content and are not deceptive or harmful out of context.
- [ ] **USEQ-EA5B740A** — Retain rights, license, attribution, consent, age, territorial, and removal metadata for published media.
- [ ] **USEQ-B32DA0FB** — Make generated PDFs, office documents, e-books, reports, invoices, statements, and manuals accessible, searchable, tagged, ordered, and usable with assistive technology.
- [ ] **USEQ-20583932** — Provide an accessible HTML alternative when a downloadable format cannot meet the required experience.
- [ ] **USEQ-ECC37B8A** — Ensure feeds contain stable identifiers, valid dates, canonical links, safe markup, complete enclosures, and no private data.
- [ ] **USEQ-8C9F5F06** — Prevent expired, paywalled, restricted, removed, or private media from remaining available through alternate encodings, thumbnails, transcripts, feeds, caches, or direct object URLs.
- [ ] **USEQ-35DC440B** — Test media playback, captions, controls, keyboard access, focus, screen-reader output, reduced motion, autoplay, bandwidth adaptation, and error recovery.

#### Site moves, redesigns, and large-scale URL changes

- [ ] **USEQ-A31027E4** — Create a complete old-to-new URL map before changing domains, protocols, paths, hosts, rendering models, taxonomies, or content platforms.
- [ ] **USEQ-4C2CB846** — Preserve the closest equivalent destination rather than redirecting all retired pages to a home page or generic category.
- [ ] **USEQ-C6F37D51** — Test redirect coverage, chains, loops, status codes, query handling, fragments, methods, and cache behavior at representative scale.
- [ ] **USEQ-9936F4FC** — Update canonical links, alternate-language annotations, structured data, internal links, sitemaps, feeds, robots rules, analytics, advertising, email templates, documentation, and external integrations.
- [ ] **USEQ-9E2A4E1C** — Verify ownership of old and new domains and retain the old origin and redirects long enough for users, integrations, and search systems to transition.
- [ ] **USEQ-9831EAF2** — Freeze unrelated high-risk changes during a major migration unless their combined risk is explicitly managed.
- [ ] **USEQ-9EE07195** — Benchmark traffic, crawl, index coverage, rankings, conversions, errors, performance, accessibility, and revenue before migration.
- [ ] **USEQ-2A898C7A** — Annotate the migration in monitoring and compare old and new URL cohorts throughout the observation period.
- [ ] **USEQ-7D37957D** — Keep a tested rollback or forward-fix plan for routing, templates, metadata, redirects, DNS, certificates, and infrastructure.
- [ ] **USEQ-F783E50B** — Do not remove redirects, old-domain control, or monitoring merely because initial traffic appears stable.
- [ ] **USEQ-F34583CD** — Document intentional losses, consolidated content, changed user journeys, and residual external-link dependencies.

#### Page experience, performance, and resilience

- [ ] **USEQ-DE2C17B6** — Set performance budgets for document size, script execution, style work, image weight, font weight, request count, third-party work, memory, and energy where relevant.
- [ ] **USEQ-F2E2BC28** — Measure loading, interaction, and visual stability using real-user field data at meaningful percentiles, segmented by device class, network, geography, and page type.
- [ ] **USEQ-924C487E** — Use current Core Web Vitals thresholds as an external baseline where applicable; at this snapshot the good thresholds are LCP at most 2.5 seconds, INP at most 200 milliseconds, and CLS at most 0.1 at the 75th percentile.
- [ ] **USEQ-4A5DD6CA** — Measure complete user journeys and business-critical interactions in addition to generic page metrics.
- [ ] **USEQ-F4889FB2** — Identify and optimize the actual largest content element, interaction delays, long tasks, layout shifts, and dependency bottlenecks.
- [ ] **USEQ-3DA4B98E** — Prioritize visible and task-critical resources without starving accessibility, correctness, security, or user controls.
- [ ] **USEQ-B9097015** — Avoid blocking rendering on nonessential third-party code, experimentation, advertising, support, or analytics.
- [ ] **USEQ-4BF32EA4** — Use caching according to HTTP semantics and prevent shared caches from serving personalized, unauthorized, stale-dangerous, or locale-inappropriate content.
- [ ] **USEQ-EFD19F60** — Define freshness, validation, invalidation, stale-if-error, offline, and recovery behavior for every important cache layer.
- [ ] **USEQ-4AD387B2** — Ensure compression, connection reuse, prioritization, and content negotiation do not create compatibility, privacy, or cache-key errors.
- [ ] **USEQ-B7AE59E5** — Limit intrusive interstitials and ensure necessary notices remain dismissible, accessible, proportionate, and non-deceptive.
- [ ] **USEQ-35878084** — Keep mobile layouts and interactions fully functional without hiding essential content or actions.
- [ ] **USEQ-4E8D466F** — Test cold cache, warm cache, signed-in, signed-out, first visit, return visit, low-memory, reduced-data, slow-network, and dependency-failure states.
- [ ] **USEQ-AAC001AD** — Detect performance regressions automatically and block releases that exceed approved budgets on critical journeys without an accepted exception.
- [ ] **USEQ-7889543B** — Ensure error, maintenance, rate-limit, consent, and fallback pages meet accessibility, metadata, status-code, and performance requirements.

#### Accessibility governance and evaluation depth

- [ ] **USEQ-B59E6A65** — Use WCAG 2.2 Level AA as the default general-purpose web baseline unless a stricter applicable requirement is identified.
- [ ] **USEQ-CF9499FF** — Create an accessibility conformance matrix mapping every applicable success criterion to representative components, templates, journeys, evidence, defects, and owners.
- [ ] **USEQ-7366DBC3** — Use a defined evaluation methodology such as WCAG-EM for representative sampling and record why the sample covers the whole product.
- [ ] **USEQ-85D0E81C** — Combine automated checks, manual inspection, keyboard testing, screen-reader testing, zoom and reflow testing, voice-input testing, and cognitive walkthroughs.
- [ ] **USEQ-49BB6BE5** — Do not treat automated scan scores as conformance because many accessibility requirements require human judgment.
- [ ] **USEQ-08F0BCEE** — Include people with disabilities in discovery, usability testing, acceptance, and remediation prioritization proportionate to product impact.
- [ ] **USEQ-73E002E7** — Test with representative combinations of browsers, operating systems, screen readers, magnification, switch control, voice control, alternative input, and user settings.
- [ ] **USEQ-192A929F** — Apply WCAG2ICT principles to non-web software and documents delivered as part of the product experience.
- [ ] **USEQ-2B6D8E2E** — Apply ATAG 2.0 when the product enables users to create or publish content, including the accessibility of the authoring interface and support for producing accessible output.
- [ ] **USEQ-A2C56DCC** — Require procurement and supplier evaluations to include accessibility evidence, limitations, remediation commitments, and regression obligations.
- [ ] **USEQ-187C6958** — Publish an accurate accessibility statement identifying scope, target, known limitations, contact route, and response process.
- [ ] **USEQ-52510127** — Provide an accessible method to report barriers without forcing users through the inaccessible feature they are reporting.
- [ ] **USEQ-59B242E7** — Prioritize defects by blocked task, affected populations, frequency, irreversibility, safety, privacy, and availability of alternatives rather than scanner severity alone.
- [ ] **USEQ-A23AF001** — Prevent accessibility regressions through component tests, linting, automated browser checks, manual release sampling, and design-system governance.
- [ ] **USEQ-2257CD10** — Train designers, writers, engineers, testers, product managers, procurement staff, support staff, and content authors for their accessibility responsibilities.
- [ ] **USEQ-21DF4CEE** — Track time to acknowledge, workaround, remediate, verify, and prevent recurrence of reported barriers.

#### Advanced accessible interaction and content controls

- [ ] **USEQ-D701A677** — Provide a mechanism to bypass repeated blocks and reach the primary content or task efficiently.
- [ ] **USEQ-D8723FD3** — Preserve meaningful reading and focus order when layout changes responsively or content is reflowed.
- [ ] **USEQ-E6ECA02A** — Do not lock content to one display orientation unless that orientation is essential.
- [ ] **USEQ-34340503** — Ensure pointer gestures have simple alternatives and pointer cancellation prevents accidental activation.
- [ ] **USEQ-71058355** — Ensure the visible label is contained in the programmatic accessible name for controls operated by speech.
- [ ] **USEQ-6D759A8B** — Ensure target size and spacing meet the selected accessibility target, including the WCAG 2.2 minimum-target criterion where applicable.
- [ ] **USEQ-10BB16EB** — Keep focused elements fully visible and not hidden by author-created content at the selected conformance level.
- [ ] **USEQ-FC04F3AD** — Offer an alternative to dragging interactions and do not make fine motor precision the only way to complete a task.
- [ ] **USEQ-7181C27B** — Make help mechanisms appear in a consistent relative order across repeated page sets.
- [ ] **USEQ-3AC2FF32** — Avoid asking users to re-enter information already provided in the same process unless necessary for security, validity, or a stated exception.
- [ ] **USEQ-330C2B73** — Support accessible authentication without requiring users to recall, transcribe, or solve cognitive-function tests without an allowed alternative.
- [ ] **USEQ-ABA32111** — Provide accessible alternatives for CAPTCHAs and ensure anti-abuse controls do not exclude users with disabilities.
- [ ] **USEQ-8FBA2834** — Announce validation, saving, loading, completion, timeout, queue, connection, and error status without unexpectedly moving focus.
- [ ] **USEQ-8B45BCD1** — Manage focus deliberately when opening and closing dialogs, changing routes, adding content, deleting content, or completing asynchronous steps.
- [ ] **USEQ-819A5DEB** — Ensure modal and nonmodal overlays have correct semantics, inert background behavior, escape paths, focus restoration, and no keyboard trap.
- [ ] **USEQ-72BF99F3** — Keep destructive, legal, financial, health, and irreversible submissions reviewable, correctable, and confirmable where required.
- [ ] **USEQ-11D31D6B** — Provide visible and programmatic instructions before input rather than relying on transient hint text.
- [ ] **USEQ-62BEA935** — Ensure data visualizations expose equivalent names, values, relationships, trends, filters, and conclusions in accessible form.
- [ ] **USEQ-E237588C** — Avoid flashing that exceeds applicable thresholds and provide controls for motion triggered by interaction.
- [ ] **USEQ-38269A8C** — Respect reduced motion, contrast, color-scheme, text-size, and other supported user preferences without losing information or function.
- [ ] **USEQ-E9FB952D** — Test cognitive load, reading complexity, timeout pressure, interruption recovery, error prevention, and consistency for critical journeys.

#### Web-platform conformance and progressive enhancement

- [ ] **USEQ-4489E942** — Use native semantic elements and platform behavior before creating custom equivalents.
- [ ] **USEQ-15CF62AD** — Use ARIA only when native semantics cannot express the required behavior, and never use ARIA to override correct native behavior without evidence.
- [ ] **USEQ-6ABAB1C2** — Validate generated markup, document structure, element nesting, duplicate identifiers, labels, names, and relevant conformance requirements.
- [ ] **USEQ-C8435CCF** — Keep the document outline, landmark structure, form relationships, table semantics, and interactive roles valid after client rendering.
- [ ] **USEQ-01BF9E5D** — Follow HTTP method, status, redirect, conditional request, range, content negotiation, authentication, and caching semantics.
- [ ] **USEQ-E5B926F7** — Set content types and character encodings accurately and prevent content sniffing where unsafe.
- [ ] **USEQ-CB46C21C** — Use the Vary mechanism and cache keys correctly when representation depends on request headers, identity, locale, device adaptation, or encoding.
- [ ] **USEQ-1BD3AB1A** — Preserve safe operation when storage is unavailable, cookies are restricted, scripts are delayed, network requests fail, or optional capabilities are denied.
- [ ] **USEQ-5B2D29F8** — Use feature detection and graceful fallback rather than brittle user-agent assumptions.
- [ ] **USEQ-10978819** — Define a support policy for platform features and provide fallback for unsupported capabilities that affect critical tasks.
- [ ] **USEQ-8459C0AE** — Avoid relying on experimental or unstable platform behavior for critical operation without compatibility evidence and a fallback.
- [ ] **USEQ-A500A820** — Keep URLs, history, reload, deep linking, copy, print, download, and open-in-new-context behavior coherent.
- [ ] **USEQ-1F736DE8** — Do not intercept standard browser behavior without a documented user benefit and equivalent accessibility.
- [ ] **USEQ-E577B44C** — Ensure progressive enhancement does not expose sensitive actions or data in a less-protected fallback path.
- [ ] **USEQ-D6F19695** — Validate content security, cross-origin isolation, embedding, permissions, referrer, and opener policies together with functionality and accessibility.
- [ ] **USEQ-253D2699** — Test extensions, password managers, translation tools, reader modes, high-contrast modes, and assistive technology for avoidable interference.

#### Monitoring, diagnostics, and continuous improvement

- [ ] **USEQ-535A89D9** — Verify ownership of relevant search-engine webmaster and analytics properties using individual, least-privilege access and recovery controls.
- [ ] **USEQ-3FBF31A0** — Monitor crawl errors, indexing states, manual actions, security issues, sitemap processing, structured-data enhancements, page experience, and removal requests.
- [ ] **USEQ-E73D403C** — Use server logs to distinguish crawler discovery, fetch, rendering-resource, redirect, rate-limit, and failure problems.
- [ ] **USEQ-34256FC3** — Monitor intended versus actual indexed URL counts by host, locale, template, content type, and canonical state.
- [ ] **USEQ-E16A44DF** — Track search queries and landing pages for user intent gaps, misleading snippets, harmful results, support demand, and content decay.
- [ ] **USEQ-9A41E78D** — Monitor accessibility errors from automated tests, manual audits, user reports, support contacts, and production telemetry without collecting unnecessary sensitive data.
- [ ] **USEQ-D8C8E1CE** — Measure performance and task success from real users while respecting consent, minimization, sampling, retention, and privacy requirements.
- [ ] **USEQ-A28B0595** — Annotate releases, migrations, incidents, experiments, template changes, content campaigns, and crawler-policy changes on relevant dashboards.
- [ ] **USEQ-1680F8E4** — Alert on sudden changes in status codes, redirect volume, crawl load, indexability directives, canonical targets, structured-data validity, accessibility failures, and performance budgets.
- [ ] **USEQ-11218020** — Use synthetic checks for critical public pages, metadata, indexing directives, links, forms, authentication, and accessibility-sensitive interactions.
- [ ] **USEQ-600FBBC0** — Retain historical crawl, index, traffic, content, accessibility, and performance evidence long enough to diagnose slow regressions and migrations.
- [ ] **USEQ-CED74251** — Investigate unexplained improvements as well as declines to detect spam, accidental exposure, bot traffic, measurement errors, or unintended indexing.
- [ ] **USEQ-1C4522F8** — Turn repeated defects into component, template, content-system, design-system, test, lint, policy, or platform improvements.
- [ ] **USEQ-EAB2E602** — Review this checklist after major search guidance changes, accessibility-standard changes, browser changes, site migrations, incidents, or material shifts in audience and content.

#### Specific release blockers

- [ ] **USEQ-2C8552A0** — Do not launch when private, embargoed, customer-specific, administrative, personal, or secret-bearing content is publicly reachable or indexable.
- [ ] **USEQ-91DC65AF** — Do not launch a critical human journey that cannot be completed using the required keyboard and assistive-technology combinations.
- [ ] **USEQ-0DB446A8** — Do not launch when a critical page returns misleading success status for errors, unauthorized states, or nonexistent resources.
- [ ] **USEQ-B00C0C53** — Do not launch a large URL migration without tested mappings, redirects, monitoring, and a recovery plan.
- [ ] **USEQ-75E4E148** — Do not launch when templates can create an unbounded crawl space or a large volume of low-value indexable pages.
- [ ] **USEQ-2672431A** — Do not launch when canonical, noindex, robots, redirect, locale, or cache signals materially contradict one another.
- [ ] **USEQ-22F45EF1** — Do not launch when primary content or controls disappear under realistic script, network, third-party, zoom, reflow, or assistive-technology conditions.
- [ ] **USEQ-6071D598** — Do not claim WCAG conformance based only on automated testing or a third-party overlay.
- [ ] **USEQ-F0E44E18** — Do not publish structured data, authorship, ratings, prices, availability, or provenance that is false, stale, hidden, or unsupported by visible content.
- [ ] **USEQ-751E3994** — Do not accept a material accessibility or search defect without an owner, user-impact assessment, compensating path, monitoring, remediation date, and expiry.

### Required evidence

- [ ] **USEQ-7EB6BE46** — Retain the applicability record, scope, assumptions, owners, reviewers, and decision date.
- [ ] **USEQ-1B7F2728** — Retain objective verification results tied to the exact assessed revision.
- [ ] **USEQ-551CBF82** — Retain representative evidence for successful, failed, boundary, misuse, degraded, and recovery behavior.
- [ ] **USEQ-5D5338EB** — Retain unresolved defects and risks with severity, impact, compensating controls, owner, and target date.
- [ ] **USEQ-8E12CCA7** — Retain evidence that controls continue to operate after deployment rather than only before release.
- [ ] **USEQ-0E5A3C0E** — Retain changes made after the evidence snapshot and show which controls were reassessed.

### Category no-go conditions

- [ ] **USEQ-3822B169** — Do not approve on checklist completion percentage when one material risk remains uncontrolled.
- [ ] **USEQ-436CB3E7** — Do not accept a tool, framework, supplier, certification, or policy claim as evidence without verifying the deployed outcome.
- [ ] **USEQ-11ED9627** — Do not mark a control Not Applicable merely because it is difficult, expensive, unfamiliar, or owned by another team.
- [ ] **USEQ-CD505063** — Do not proceed when the organization cannot detect, contain, recover from, and communicate a foreseeable high-impact failure.

## Accessibility, SEO, Content Quality, and Digital Discoverability

_Consolidated from `final consolidated corpus/02-human-experience-accessibility-content-seo-internationalization.md#Accessibility, SEO, Content Quality, and Digital Discoverability`; 321 non-duplicative controls._

### Scope, ownership, and measurable outcomes

- [ ] **USEQ-76DB1417** — Inventory every public page, authenticated page, document, media asset, embedded component, native or hybrid app surface, authoring interface, and search-visible endpoint in scope.
- [ ] **USEQ-EB22BF99** — Identify which surfaces are intended to be public, indexable, searchable, shareable, archived, private, or removed.
- [ ] **USEQ-6B251F18** — Assign accountable owners for accessibility, technical discoverability, content quality, metadata, structured data, translations, and search monitoring.
- [ ] **USEQ-DDC9A009** — Document the supported user groups, disability needs, assistive technologies, browsers, devices, locales, search engines, crawlers, and content channels.
- [ ] **USEQ-B04D0821** — Define accessibility conformance targets by product surface and jurisdiction rather than relying on a single organization-wide label.
- [ ] **USEQ-D66C6A9A** — Define discoverability objectives that prioritize user value, accurate representation, and legitimate acquisition rather than ranking manipulation.
- [ ] **USEQ-0103E2BB** — Define measurable success criteria for task completion, accessibility defects, search coverage, crawl errors, indexing, content freshness, and performance.
- [ ] **USEQ-6B827DAC** — Distinguish mandatory conformance requirements from aspirational quality improvements.
- [ ] **USEQ-B9158921** — Identify third-party components and content that can prevent accessibility or discoverability conformance.
- [ ] **USEQ-62232063** — Ensure procurement and supplier contracts include applicable accessibility, remediation, metadata, rendering, availability, and change-notification obligations.
- [ ] **USEQ-9AC255A5** — Maintain an inventory of public domains, subdomains, alternate hosts, preview environments, staging environments, microsites, and retired properties.
- [ ] **USEQ-0BB3E806** — Prevent preview, test, administrative, personalized, confidential, and pre-release content from becoming publicly discoverable.
- [ ] **USEQ-D5B7E6DD** — Document how accessibility and search behavior are affected by personalization, experiments, geolocation, consent, authentication, and feature flags.
- [ ] **USEQ-FCC69E21** — Include accessibility and discoverability acceptance criteria in product requirements and Definition of Done.
- [ ] **USEQ-51165919** — Give users a clear channel to report accessibility, content, metadata, or search-discovery problems.
- [ ] **USEQ-F6820614** — Review this scope after redesigns, domain changes, platform migrations, acquisitions, localization expansion, and introduction of generated content.

### Accessibility standards and evaluation methodology

- [ ] **USEQ-1F9E3B23** — Use WCAG 2.2 as the general web-content baseline unless a binding requirement specifies another version or stricter target.
- [ ] **USEQ-698BB07E** — Document the selected conformance level and every additional requirement arising from law, procurement, contract, safety, or user research.
- [ ] **USEQ-EFFCCD14** — Use a defined evaluation methodology such as WCAG-EM 2.0 to establish scope, representative samples, states, processes, and reporting.
- [ ] **USEQ-FE8CC474** — Include complete processes and multi-step journeys rather than evaluating isolated pages only.
- [ ] **USEQ-7758C65E** — Include common pages, templates, components, states, error paths, authenticated areas, documents, media, and dynamically generated variants in the evaluation sample.
- [ ] **USEQ-3801322B** — Evaluate content at supported zoom, text spacing, viewport, orientation, contrast, color-scheme, forced-color, reduced-motion, and platform settings.
- [ ] **USEQ-A371549C** — Combine automated testing, expert manual review, keyboard testing, assistive-technology testing, and user testing; do not treat automation as proof of conformance.
- [ ] **USEQ-6D505348** — Record the exact standards, techniques, test tools, versions, assistive technologies, browsers, operating systems, devices, samples, and limitations used.
- [ ] **USEQ-A4088B58** — Distinguish normative requirements from informative techniques and tool heuristics.
- [ ] **USEQ-88E2245C** — Classify defects by user impact, task blockage, frequency, reach, reversibility, and legal or safety significance.
- [ ] **USEQ-9F7B27E7** — Do not average inaccessible critical journeys against large numbers of passing low-value pages.
- [ ] **USEQ-38B52E84** — Treat a task-blocking defect in a critical journey as a release blocker unless an equally effective accessible alternative is available and verified.
- [ ] **USEQ-7BCD3094** — Require independent accessibility review for high-impact, public-sector, education, employment, finance, healthcare, safety, or essential-service products.
- [ ] **USEQ-1F39E16F** — Include people with relevant disabilities in research and validation when the product impact warrants it.
- [ ] **USEQ-50EFB4BA** — Compensate participants appropriately and protect their privacy and dignity.
- [ ] **USEQ-B7ADCF31** — Publish an accessibility statement that accurately describes scope, known limitations, contact routes, and response expectations.
- [ ] **USEQ-8AC646AD** — Keep conformance reports tied to a specific product version and date; do not reuse them after material change without reassessment.
- [ ] **USEQ-1F6E49BB** — Reassess accessibility continuously rather than only before launch.

### Semantic structure, content alternatives, and programmatic relationships

- [ ] **USEQ-23049AC4** — Use native semantic elements and platform controls whenever they satisfy the required behavior.
- [ ] **USEQ-0FD25F2E** — Give every page or view a meaningful, unique title that reflects its purpose or topic.
- [ ] **USEQ-A11CE7A5** — Use headings in a logical hierarchy that communicates document structure rather than visual styling alone.
- [ ] **USEQ-F4EE10A3** — Expose landmarks, regions, lists, tables, definitions, quotations, emphasis, and other structural relationships programmatically.
- [ ] **USEQ-F65C7A78** — Give every interactive element an accurate accessible name, role, value, state, and description where needed.
- [ ] **USEQ-EB41475C** — Ensure visible labels and accessible names are consistent enough for speech-input users to identify controls.
- [ ] **USEQ-43615406** — Provide equivalent text alternatives for informative images, icons, charts, diagrams, maps, and controls.
- [ ] **USEQ-1C090BD9** — Mark decorative content so that it does not create noise for assistive technologies.
- [ ] **USEQ-44FF600B** — Provide long descriptions, data tables, transcripts, or equivalent alternatives for complex visual information.
- [ ] **USEQ-CE72BA84** — Ensure images of text are not used when real text can provide the same presentation, except for legitimate exceptions.
- [ ] **USEQ-B96A7EB0** — Use table headers, captions, scopes, and relationships that make data tables understandable without visual position.
- [ ] **USEQ-4D11DB36** — Expose status messages, validation results, loading states, progress, notifications, and asynchronous updates without forcing focus changes.
- [ ] **USEQ-BA50C7B0** — Preserve semantics when content is reordered, virtualized, collapsed, expanded, filtered, or rendered lazily.
- [ ] **USEQ-438CC44E** — Ensure custom widgets follow established accessibility interaction patterns and do not reproduce native controls poorly.
- [ ] **USEQ-352ACD0C** — Use ARIA only where needed, apply supported roles and properties correctly, and never use ARIA to conceal broken semantics.
- [ ] **USEQ-33D779E1** — Ensure generated markup remains valid enough for user agents to derive consistent semantics.
- [ ] **USEQ-531C7FDB** — Declare the natural language of pages and language changes within content.
- [ ] **USEQ-ED944B62** — Expose reading order, relationships, and alternative content in exported and downloadable documents.

### Keyboard, focus, pointer, voice, and alternative input

- [ ] **USEQ-1ACA514F** — Make every actionable function operable with a keyboard without requiring pointer input.
- [ ] **USEQ-8FBB63F5** — Prevent keyboard traps and provide a documented exit from any constrained interaction.
- [ ] **USEQ-B05009C5** — Maintain a logical focus order that follows meaning and task sequence.
- [ ] **USEQ-DB2C10EF** — Provide a clearly visible focus indicator with sufficient contrast and area.
- [ ] **USEQ-018113B7** — Ensure focused elements are not hidden behind sticky headers, banners, dialogs, virtual keyboards, or overlays.
- [ ] **USEQ-E32BDB30** — Move focus intentionally after route changes, dialog opening, content replacement, deletion, and other context changes.
- [ ] **USEQ-C434159D** — Restore focus to a sensible location after dialogs, menus, popovers, and temporary surfaces close.
- [ ] **USEQ-591A7C39** — Provide skip mechanisms or equivalent navigation for repeated blocks.
- [ ] **USEQ-E3635275** — Avoid single-character shortcuts that conflict with assistive technology, or make them remappable or disableable.
- [ ] **USEQ-54E8800E** — Provide non-drag alternatives for functionality that depends on dragging.
- [ ] **USEQ-53DCF354** — Provide alternatives to path-based, multipoint, or complex gestures.
- [ ] **USEQ-6E86BBDD** — Make pointer targets large enough and sufficiently separated for the target population and context.
- [ ] **USEQ-639D4B06** — Do not trigger actions solely on pointer-down when cancellation or reversal is needed.
- [ ] **USEQ-9387A6EF** — Support pointer cancellation and undo for consequential actions.
- [ ] **USEQ-4054AA2A** — Do not require device motion, shaking, tilting, or biometric gestures when an accessible alternative can be provided.
- [ ] **USEQ-59674F93** — Ensure visible control text can be used by voice-control users to activate the corresponding control.
- [ ] **USEQ-A5F42D98** — Test switch-control, voice-control, eye-gaze, or other alternative input when relevant to the supported population.
- [ ] **USEQ-A66695B5** — Ensure timeout, focus, and interaction behavior remains correct under slower input and assistive-technology use.

### Visual presentation, audio, video, motion, and sensory safety

- [ ] **USEQ-9645FB34** — Meet the required contrast ratios for text, controls, focus indicators, icons, charts, and other meaningful visual information.
- [ ] **USEQ-047A21BC** — Do not rely on color, position, shape, sound, motion, or another single sensory characteristic to communicate meaning.
- [ ] **USEQ-83473448** — Allow text to be resized and spaced without loss of content, function, or readability.
- [ ] **USEQ-24B7B3FD** — Allow content to reflow without two-dimensional scrolling except where the content intrinsically requires it.
- [ ] **USEQ-E9678D71** — Avoid clipping, overlap, disappearance, and unusable overlays at supported zoom and viewport sizes.
- [ ] **USEQ-417FF771** — Honor user preferences for reduced motion, contrast, color scheme, text size, captions, and other accessibility settings where available.
- [ ] **USEQ-16826B5E** — Avoid flashes and visual patterns that can exceed seizure or physical-reaction thresholds.
- [ ] **USEQ-EB964FEA** — Provide controls to pause, stop, hide, or reduce moving, blinking, auto-updating, or auto-playing content where required.
- [ ] **USEQ-81BE1D6D** — Do not start audio automatically unless the user can stop or control it immediately.
- [ ] **USEQ-B06A6E6F** — Provide synchronized captions for prerecorded and live audio where required.
- [ ] **USEQ-499B8B78** — Provide transcripts for audio-only content and equivalent alternatives for video-only content.
- [ ] **USEQ-5966A22D** — Provide audio description or an equivalent description of important visual information where required.
- [ ] **USEQ-0F32C504** — Ensure media controls are keyboard operable, labeled, visible, and usable with assistive technology.
- [ ] **USEQ-00D74F4E** — Preserve caption timing, speaker identification, sound cues, and readability.
- [ ] **USEQ-7E071F27** — Do not encode essential information in low-contrast placeholders or transient toast messages.
- [ ] **USEQ-296488AE** — Ensure charts and visualizations expose labels, units, scales, relationships, uncertainty, and underlying data accessibly.
- [ ] **USEQ-E9136303** — Ensure maps have non-map alternatives for essential tasks and information.
- [ ] **USEQ-9013CE98** — Test forced-colors and high-contrast modes so that controls, focus, selection, errors, and state remain perceivable.
- [ ] **USEQ-33A15DF0** — Keep line length, spacing, typography, hierarchy, and density suitable for reading and cognitive access.
- [ ] **USEQ-DDF0FE8A** — Provide user control over distracting animations, parallax, carousels, and background media.

### Cognitive accessibility, forms, authentication, and error prevention

- [ ] **USEQ-47E15B30** — Use clear, concrete, consistent language appropriate to the audience.
- [ ] **USEQ-EA2525B5** — Explain unfamiliar terms, abbreviations, icons, and consequences at the point of need.
- [ ] **USEQ-BF3E67FB** — Keep navigation, help, labels, and repeated controls consistent across the product.
- [ ] **USEQ-1DCF1FC6** — Break complex tasks into understandable steps and show progress and remaining work.
- [ ] **USEQ-3EC1EBBA** — Allow users to review, correct, save, and resume complex or high-stakes tasks.
- [ ] **USEQ-EAAE1A20** — Minimize memory burden and avoid requiring users to re-enter information already available to the system.
- [ ] **USEQ-E1902ED3** — Do not use deceptive urgency, confusing defaults, hidden costs, forced continuity, or other dark patterns.
- [ ] **USEQ-8ABF95B0** — Label every form control programmatically and visibly unless an accessible equivalent is unambiguous.
- [ ] **USEQ-E35FCCCE** — Identify required fields, formats, constraints, and examples before submission.
- [ ] **USEQ-EAF2117D** — Use appropriate input purpose and autocomplete semantics where they improve accessibility and privacy.
- [ ] **USEQ-79FA79CC** — Associate errors with the affected fields and provide an accessible error summary for multi-field failures.
- [ ] **USEQ-C69C5821** — Explain errors in plain language and provide specific recovery instructions.
- [ ] **USEQ-FC0BCDA1** — Preserve valid input after recoverable errors and prevent duplicate submission.
- [ ] **USEQ-38360ECC** — Provide confirmation, review, undo, or reversal for legal, financial, destructive, or high-impact submissions.
- [ ] **USEQ-47A0D215** — Warn users before session expiry and permit extension or data preservation where required.
- [ ] **USEQ-BB3E050D** — Do not make authentication depend solely on solving cognitive-function tests, transcribing distorted content, or remembering arbitrary secrets.
- [ ] **USEQ-375A6696** — Support password managers, paste, accessible multifactor methods, and accessible account recovery.
- [ ] **USEQ-9DB6D254** — Provide alternatives when biometrics, CAPTCHAs, device possession, or a particular communication channel are inaccessible.
- [ ] **USEQ-4BC4B0A8** — Ensure help and human-support routes are accessible and do not require the inaccessible step that caused the problem.
- [ ] **USEQ-EE9AEBED** — Test error recovery under screen-reader, magnification, voice-control, and keyboard-only use.

### Responsive, mobile, assistive-technology, and document interoperability

- [ ] **USEQ-AE14ECFB** — Test representative critical journeys across the supported browser, operating-system, device, and assistive-technology matrix.
- [ ] **USEQ-7F48A297** — Include multiple major screen readers and browser combinations appropriate to the user population.
- [ ] **USEQ-B09FCA61** — Include screen magnification, zoom, reflow, voice control, switch control, captions, and platform accessibility settings as applicable.
- [ ] **USEQ-7D7CBD51** — Ensure orientation is not locked unless a specific orientation is essential.
- [ ] **USEQ-34A9C539** — Ensure virtual keyboards, safe areas, display cutouts, and magnification do not cover controls or instructions.
- [ ] **USEQ-86065889** — Ensure touch exploration and mobile screen-reader reading order match the visual and task order.
- [ ] **USEQ-6E1DB57C** — Apply WCAG 2.2 principles to native and hybrid mobile applications and supplement them with platform accessibility guidance.
- [ ] **USEQ-E6287819** — Ensure offline, cached, and installed experiences preserve accessibility.
- [ ] **USEQ-4ABB7338** — Ensure notifications, deep links, share sheets, widgets, and system dialogs remain understandable and operable.
- [ ] **USEQ-45B46C2B** — Produce accessible PDF, office, email, and other downloadable content, using an applicable document standard such as PDF/UA where required.
- [ ] **USEQ-CACB9B46** — Preserve document tags, headings, reading order, bookmarks, language, alternate text, tables, links, and form semantics during generation and conversion.
- [ ] **USEQ-226021B5** — Do not publish scanned-image documents without accurate text recognition and structural remediation when accessible text is required.
- [ ] **USEQ-FC22FBD4** — Provide an accessible HTML or equivalent alternative when a document format cannot meet user needs.
- [ ] **USEQ-1C53CC23** — Ensure print and export layouts do not remove essential information or accessibility.
- [ ] **USEQ-A00433C7** — Verify that assistive-technology behavior remains correct after component-library, browser-engine, platform, or PDF-generation upgrades.
- [ ] **USEQ-2ED3E53A** — Document unsupported combinations and provide a workable alternative rather than silently failing.

### User agents, embedded viewers, media players, and assistive interoperability

- [ ] **USEQ-7A4E29FA** — Apply UAAG 2.0 as informative guidance when the product retrieves, renders, navigates, transforms, or plays web or digital content as a browser, viewer, reader, extension, media player, embedded web view, document viewer, or similar user agent.
- [ ] **USEQ-B12A0439** — Make the user agent interface itself perceivable, operable, understandable, and compatible with assistive technologies.
- [ ] **USEQ-584BDA9F** — Expose rendered content, controls, selection, focus, caret, viewport, playback, status, and document structure through supported accessibility APIs.
- [ ] **USEQ-43A0D27C** — Preserve content semantics and accessibility information rather than flattening or replacing them during rendering, transformation, reader mode, translation, annotation, or export.
- [ ] **USEQ-E5B715BD** — Let users control text size, zoom, colors, contrast, fonts, spacing, animation, autoplay, audio, captions, playback rate, and other presentation characteristics where the platform role permits it.
- [ ] **USEQ-07226877** — Provide keyboard and alternative-input access to browser chrome, media controls, document navigation, extensions, permissions, downloads, history, and settings.
- [ ] **USEQ-9EC00A36** — Provide navigation by headings, landmarks, links, controls, tables, pages, time, chapters, annotations, search results, and other meaningful structures when present.
- [ ] **USEQ-30ACE3C6** — Communicate security, permission, download, certificate, mixed-content, popup, and content-blocking messages accessibly.
- [ ] **USEQ-0527C48F** — Ensure accessibility extensions and assistive technologies do not require excessive privileges or expose browsing and content data without informed control.
- [ ] **USEQ-C22C2973** — Test embedded viewers and web views with the host application, content, accessibility tree, focus transfer, keyboard handling, zoom, and assistive technologies as one integrated experience.
- [ ] **USEQ-804BFA3B** — Provide an accessible fallback or alternate application when a required format or embedded user agent cannot meet critical user needs.
- [ ] **USEQ-8DD986C8** — Treat user-agent and assistive-technology compatibility regressions as product defects even when the underlying content has not changed.

### Authoring tools, design systems, third parties, and procurement

- [ ] **USEQ-56988474** — Apply ATAG 2.0 when the product lets users create or edit web or digital content.
- [ ] **USEQ-8CD38B0C** — Make the authoring interface itself accessible to authors with disabilities.
- [ ] **USEQ-E900C7C6** — Generate accessible output by default and preserve accessible information during editing, import, export, copying, and transformation.
- [ ] **USEQ-3C4F60FD** — Prompt authors for missing alternate text, labels, headings, captions, language, table headers, and other required accessibility information.
- [ ] **USEQ-547FB7E6** — Do not overwrite correct author-supplied accessibility data with lower-quality generated values.
- [ ] **USEQ-F96E1B38** — Provide accessible templates, examples, validators, repair suggestions, and previews.
- [ ] **USEQ-8E6F38D9** — Make accessibility checking discoverable and integrated into normal authoring workflows.
- [ ] **USEQ-4F0EF52C** — Include accessibility requirements in design-system component contracts and acceptance tests.
- [ ] **USEQ-4202F7C4** — Provide documented keyboard behavior, semantics, states, content constraints, and known limitations for every shared component.
- [ ] **USEQ-433B5762** — Prevent design tokens from allowing inaccessible contrast, focus, motion, spacing, or target sizes without an explicit reviewed exception.
- [ ] **USEQ-9282E33E** — Require accessibility review before introducing or materially changing a shared component.
- [ ] **USEQ-8C17324F** — Version accessibility behavior and communicate breaking changes to consumers.
- [ ] **USEQ-E6611C2E** — Evaluate third-party widgets, payment flows, identity screens, chat tools, media players, consent tools, and embeds before adoption.
- [ ] **USEQ-F28C400D** — Contractually require timely remediation or an accessible alternative for critical third-party barriers.
- [ ] **USEQ-F8516AC8** — Keep a replacement or isolation plan for inaccessible supplier components.
- [ ] **USEQ-8DC4E850** — Require vendors to provide current, scoped, evidence-based accessibility conformance information rather than marketing claims.
- [ ] **USEQ-C6CE00CE** — Test the integrated product; supplier documentation alone is not evidence that the deployed experience is accessible.
- [ ] **USEQ-2AE7001B** — Ensure procurement scoring considers actual task accessibility, not only presence of a conformance report.

### Accessibility regression, reporting, and operations

- [ ] **USEQ-00EB0815** — Automate checks for repeatable detectable defects in components, templates, rendered pages, documents, and critical flows.
- [ ] **USEQ-393454F3** — Use automated rules as a floor and track false positives, false negatives, rule coverage, and tool versions.
- [ ] **USEQ-0532FA61** — Keep manual test scripts for keyboard, focus, zoom, text spacing, screen readers, forms, media, and cognitive access.
- [ ] **USEQ-BAA82799** — Include accessibility assertions in component, integration, end-to-end, visual-regression, and document-generation tests where reliable.
- [ ] **USEQ-D3049D6B** — Require accessibility review when content structure, interaction patterns, navigation, design tokens, generated documents, or major dependencies change.
- [ ] **USEQ-CD0F0EFF** — Monitor production for missing names, invalid structure, contrast regressions, inaccessible dialogs, broken skip links, and document-generation failures where feasible.
- [ ] **USEQ-53280330** — Collect accessibility feedback without requiring users to disclose a diagnosis.
- [ ] **USEQ-0AC7F3A1** — Triage reports based on actual task impact and affected users, not only automated severity.
- [ ] **USEQ-3E16ED92** — Provide temporary accessible alternatives while material defects are being repaired.
- [ ] **USEQ-A6A392BF** — Retest fixes with the input mode or assistive technology that exposed the defect.
- [ ] **USEQ-264DBF6E** — Link recurring defects to component, authoring, process, training, or governance fixes rather than patching pages individually.
- [ ] **USEQ-F6CFFA0D** — Maintain an accessibility defect backlog with owner, impact, affected surfaces, workaround, target date, and evidence of closure.
- [ ] **USEQ-CF5FB21E** — Include accessibility incidents in operational and post-incident review when users lose access to essential functions.
- [ ] **USEQ-0C55957B** — Monitor third-party changes for accessibility regression.
- [ ] **USEQ-4FE30275** — Review accessibility statements and conformance reports whenever scope or known limitations change.
- [ ] **USEQ-036ED9BF** — Measure critical-task success and barrier frequency, not only raw issue counts.
- [ ] **USEQ-D5B4830A** — Prevent accessibility metrics from being gamed by excluding hard pages, states, users, or assistive technologies.
- [ ] **USEQ-6727CBCA** — Retain sufficient evaluation evidence to support internal assurance, procurement, customer, and regulatory review.

### SEO, discoverability, and ethical search governance

- [ ] **USEQ-B7F56A81** — Define SEO as making legitimate public content discoverable, understandable, and useful; do not treat it as manipulation of ranking systems.
- [ ] **USEQ-A4774D9F** — Follow applicable search-engine technical requirements, spam policies, and webmaster guidance.
- [ ] **USEQ-2FE2A2D4** — Prioritize people-first content and task value over keyword density, doorway pages, link schemes, scaled low-value pages, or hidden content.
- [ ] **USEQ-76C0B794** — Do not promise ranking, indexing, rich results, traffic, or generative-search inclusion because these are not fully controlled by the publisher.
- [ ] **USEQ-D77DC459** — Ensure search optimization does not reduce accessibility, privacy, security, performance, or content accuracy.
- [ ] **USEQ-E969C3C9** — Assign owners for technical SEO, editorial quality, structured data, international targeting, migrations, and monitoring.
- [ ] **USEQ-D33DFAE7** — Document which engines, regional services, vertical search products, and AI or content crawlers matter to the product.
- [ ] **USEQ-8959440E** — Create a policy for crawler access, content licensing, training use, archival access, and high-volume automated retrieval.
- [ ] **USEQ-625225E7** — Distinguish crawler directives from authentication and authorization; never rely on crawler compliance to protect confidential information.
- [ ] **USEQ-E906E553** — Ensure search-visible titles, descriptions, snippets, dates, prices, availability, ratings, and identity claims accurately reflect the page.
- [ ] **USEQ-0A3AA7F8** — Prevent experiments, personalization, consent states, and geolocation from serving materially misleading content to crawlers or users.
- [ ] **USEQ-9B0A4681** — Do not cloak by serving deceptive content to crawlers that differs materially from what users receive.
- [ ] **USEQ-73D5B4DD** — Disclose sponsored, affiliate, user-generated, and automated content where required and apply appropriate link attributes and editorial controls.
- [ ] **USEQ-A0D8CC51** — Maintain evidence for factual, medical, legal, financial, safety, and other high-impact claims.
- [ ] **USEQ-36B07DA8** — Use third-party SEO tools as diagnostic aids, not as an authoritative substitute for official guidance or product judgment.
- [ ] **USEQ-0B3A3885** — Review search practices after major search-engine policy changes and document material decisions.

### Crawlability, indexability, crawler controls, and exposure safety

- [ ] **USEQ-B9CD4463** — Ensure every intended public page is reachable through stable crawlable links or an intentional discovery mechanism.
- [ ] **USEQ-FB71685A** — Ensure intended indexable pages return a successful status and usable indexable content.
- [ ] **USEQ-838FA31D** — Ensure pages that must remain private require authentication or equivalent access control.
- [ ] **USEQ-B4F0B76C** — Use robots.txt in accordance with RFC 9309 and understand that it controls compliant crawling, not confidentiality or authorization.
- [ ] **USEQ-50530FEB** — Keep robots.txt syntactically valid, reachable, monitored, and consistent across relevant hosts and protocols.
- [ ] **USEQ-15FD8AF5** — Do not block scripts, styles, images, or other resources required for a crawler to understand and render intended public content.
- [ ] **USEQ-39597E04** — Use `noindex` or equivalent supported indexing controls for public-but-nonindexable content, and ensure crawlers can retrieve the directive when required.
- [ ] **USEQ-36224CFC** — Use response-header indexing controls for non-HTML resources where appropriate.
- [ ] **USEQ-4D234580** — Do not use robots.txt as a substitute for `noindex`, deletion, authentication, or canonicalization.
- [ ] **USEQ-162168BC** — Prevent staging, preview, search-result, filter, cart, account, administrative, debug, and internal pages from accidental indexing.
- [ ] **USEQ-314CF702** — Prevent personal, confidential, tokenized, one-time, or sensitive query parameters from entering public indexes, analytics, logs, sitemaps, or links.
- [ ] **USEQ-85574DF8** — Ensure deleted content returns the intended permanent or temporary status and is removed from sitemaps and internal links.
- [ ] **USEQ-CFC54B83** — Use temporary unavailability statuses during short maintenance rather than returning misleading success pages.
- [ ] **USEQ-4539572E** — Avoid soft errors in which a missing, blocked, or failed resource returns a successful status with an error message.
- [ ] **USEQ-6587FF90** — Handle DNS, TLS, redirect, timeout, and server failures so crawlers do not index unintended fallback or error content.
- [ ] **USEQ-14629F73** — Bound crawler traffic and expensive URL spaces without hiding important content.
- [ ] **USEQ-01A25EAE** — Control faceted navigation, calendars, internal search, session identifiers, sort orders, and unbounded parameters that create crawl traps.
- [ ] **USEQ-FB490C89** — Detect scraper or bot abuse separately from legitimate search crawling and verify crawler identity before granting special treatment.
- [ ] **USEQ-6908CAC0** — Keep private media and documents outside public URL spaces or protect them independently of page-level directives.
- [ ] **USEQ-B8CC27A7** — Test crawl and indexing controls from the public network against the deployed production configuration.

### URLs, canonicalization, redirects, site architecture, and links

- [ ] **USEQ-166444D1** — Use stable, human-understandable, consistently encoded URLs without unnecessary identifiers or volatile state.
- [ ] **USEQ-4D56F18B** — Choose one intended canonical URL for substantially duplicate content and keep canonical signals consistent.
- [ ] **USEQ-2A892965** — Align canonical declarations, redirects, internal links, sitemaps, alternate-language links, and structured data.
- [ ] **USEQ-FF9B9125** — Do not canonicalize materially different pages merely to suppress duplicate-content warnings.
- [ ] **USEQ-621650F8** — Use permanent redirects for permanent moves and temporary redirects only for genuinely temporary routing.
- [ ] **USEQ-4B6B04A6** — Avoid redirect chains, loops, mass redirects to irrelevant destinations, and redirects that discard essential path or query meaning.
- [ ] **USEQ-7BDA3D27** — Preserve user intent and equivalent content during redirects and site migrations.
- [ ] **USEQ-235346A8** — Normalize host, scheme, trailing-slash, case, encoding, and default-document variants consistently.
- [ ] **USEQ-A21A412A** — Prevent open redirects and untrusted redirect targets.
- [ ] **USEQ-B5FBA63D** — Make links real, crawlable links with meaningful anchor text when navigation or discovery is intended.
- [ ] **USEQ-1FD6CF1D** — Do not require pointer-only interaction, client-only events, or form submission for discovery of essential public pages.
- [ ] **USEQ-DA53A718** — Use a logical information architecture with clear category, hierarchy, breadcrumb, and contextual relationships.
- [ ] **USEQ-65B7D11F** — Keep important content within a reasonable navigation path without relying solely on sitemaps.
- [ ] **USEQ-21260A94** — Identify orphan pages and either link, merge, archive, redirect, noindex, or remove them intentionally.
- [ ] **USEQ-3056A833** — Avoid deceptive internal-link patterns, repetitive keyword stuffing, and hidden links.
- [ ] **USEQ-79E91255** — Use appropriate attributes for sponsored, untrusted user-generated, and intentionally non-endorsed links.
- [ ] **USEQ-83F4A9B9** — Check outgoing and internal links for breakage, unsafe destinations, obsolete redirects, and changed ownership.
- [ ] **USEQ-023C3010** — Prevent expired domains, abandoned subdomains, and dangling DNS or hosting references from becoming takeover or reputation risks.
- [ ] **USEQ-9256F95C** — Preserve meaningful URLs and redirect mappings long enough to satisfy users, external links, contracts, and retention needs.
- [ ] **USEQ-D5EEA737** — Document URL ownership and change approval for critical public routes.

### Rendering, JavaScript, mobile parity, and search performance

- [ ] **USEQ-292A2F47** — Ensure essential content, links, metadata, canonical signals, and structured data are available to supported crawlers after rendering.
- [ ] **USEQ-EAA91E46** — Test both initial server responses and rendered output for pages that depend on client-side execution.
- [ ] **USEQ-4D7D1F86** — Do not require user gestures, login, unsupported storage, denied consent, or fragile client state for crawlers to access intended public content.
- [ ] **USEQ-ABD33443** — Ensure rendering failures do not expose empty shells, perpetual loading states, generic errors, or unrelated fallback content.
- [ ] **USEQ-5E42DAC9** — Use progressive enhancement or an equivalent robust rendering strategy for critical public information where practical.
- [ ] **USEQ-98189E72** — Avoid dynamic rendering or crawler-specific rendering unless a documented temporary need exists and parity is continuously verified.
- [ ] **USEQ-08D607D9** — Keep crawler-facing and user-facing content materially equivalent.
- [ ] **USEQ-26628827** — Ensure mobile and desktop versions contain equivalent primary content, metadata, structured data, alternate text, and indexing directives.
- [ ] **USEQ-097BA4D0** — Use responsive behavior or explicit alternate relationships consistently when separate mobile URLs exist.
- [ ] **USEQ-AB375D1C** — Do not hide important mobile content behind interactions unavailable during crawl or rendering.
- [ ] **USEQ-D60F26B0** — Ensure lazy-loaded content and images can be discovered without unsupported scrolling or interaction.
- [ ] **USEQ-56F842F4** — Keep critical rendered links in the document structure rather than generating them only after opaque events.
- [ ] **USEQ-6B6351C1** — Meet approved user-performance budgets and monitor Core Web Vitals using field data where applicable.
- [ ] **USEQ-C2EB8D2B** — Use the 75th percentile segmented by meaningful device classes and user populations for user-facing performance decisions.
- [ ] **USEQ-D9711648** — Optimize loading, responsiveness, and visual stability without sacrificing accessibility, correctness, or content completeness.
- [ ] **USEQ-9DAE647D** — Control third-party scripts, tag managers, ads, consent tools, and experimentation code that can delay rendering or change indexable content.
- [ ] **USEQ-A58CEC40** — Ensure server-side rendering, pre-rendering, hydration, streaming, and caching produce consistent semantics and metadata.
- [ ] **USEQ-1F4F7014** — Test rendering and indexing after framework, browser, CDN, bot-detection, consent, or edge-compute changes.

### Sitemaps, metadata, structured data, media, and result presentation

- [ ] **USEQ-3C32DD0A** — Publish sitemaps only when they improve discovery or monitoring, and include only canonical, intended, accessible URLs.
- [ ] **USEQ-CD50C547** — Keep sitemap modification dates truthful and update files when content changes materially.
- [ ] **USEQ-F8B7471E** — Split and index large sitemaps within protocol limits and monitor retrieval failures.
- [ ] **USEQ-7ACEC211** — Remove deleted, redirected, duplicate, noindex, private, and erroring URLs from sitemaps.
- [ ] **USEQ-E3F49EF9** — Provide accurate unique page titles and concise descriptions that match visible content.
- [ ] **USEQ-BC78A39A** — Use metadata and social-preview data that accurately identifies the page, image, locale, and publisher.
- [ ] **USEQ-5BF1A157** — Prevent private data, tokens, internal names, unreleased information, or misleading copy from entering metadata and previews.
- [ ] **USEQ-C6547A77** — Use structured data only when it accurately describes visible content and meets the applicable feature guidelines.
- [ ] **USEQ-532EA584** — Validate structured data syntax, required properties, types, URLs, dates, prices, availability, identity, and relationships.
- [ ] **USEQ-D70319B2** — Keep structured data synchronized with visible content and authoritative backend state.
- [ ] **USEQ-8F3F7D97** — Do not create fake ratings, reviews, authorship, organizations, products, jobs, events, FAQs, or other entities.
- [ ] **USEQ-A00E833A** — Use stable identifiers and appropriate same-as or identity references without creating false equivalence.
- [ ] **USEQ-30CC73DE** — Treat rich-result eligibility as conditional and never as a guaranteed outcome.
- [ ] **USEQ-4A5CB8A6** — Provide image dimensions, descriptive alternatives, suitable quality, licensing and credit metadata, and stable URLs where relevant.
- [ ] **USEQ-7F57174C** — Provide accurate video titles, descriptions, thumbnails, durations, transcripts, captions, key moments, and availability where applicable.
- [ ] **USEQ-2CF8CC0E** — Provide publication, modification, author, publisher, and correction information when relevant to trust and freshness.
- [ ] **USEQ-10ABDDAA** — Ensure search snippets and previews do not expose content that users cannot legitimately access.
- [ ] **USEQ-B6CAF88F** — Test metadata and structured data after templates, CMS, localization, product feeds, inventory, or pricing logic changes.

### International and multilingual discoverability

- [ ] **USEQ-8244B358** — Use valid BCP 47 language tags and declare page language and regional variation accurately.
- [ ] **USEQ-F0424A30** — Use locale-aware content, currency, date, address, units, and terminology rather than mechanical word substitution.
- [ ] **USEQ-3ABDCE0A** — Provide stable distinct URLs for localized content when independent indexing is intended.
- [ ] **USEQ-3CA1A6EE** — Use reciprocal alternate-language annotations such as `hreflang` consistently when supported and applicable.
- [ ] **USEQ-70B062EB** — Include a self-reference and an appropriate fallback where alternate-language annotations require them.
- [ ] **USEQ-80D06A16** — Ensure canonical and alternate-language signals do not contradict each other.
- [ ] **USEQ-6CEB7564** — Keep localized pages equivalent in essential content, metadata, structured data, navigation, and accessibility.
- [ ] **USEQ-7B2A7EE1** — Do not automatically redirect users or crawlers solely by inferred language or location without an accessible override and stable URL.
- [ ] **USEQ-0D4E0314** — Allow users and crawlers to reach every supported locale through crawlable links.
- [ ] **USEQ-214B3B45** — Do not treat machine-translated, unreviewed, low-value pages as a scale strategy.
- [ ] **USEQ-6E51BB7D** — Use qualified human review for high-impact, legal, safety, medical, financial, and culturally sensitive content.
- [ ] **USEQ-748B12CD** — Preserve translated titles, descriptions, alternate text, captions, structured data, and error content.
- [ ] **USEQ-516A203A** — Monitor each locale separately for indexing, errors, search demand, content staleness, and broken alternate mappings.
- [ ] **USEQ-69E6DA5E** — Retire or redirect locales deliberately when support ends, and communicate the change to users.

### Content quality, editorial governance, trust, and provenance

- [ ] **USEQ-72E01ED6** — Give every material public content area an accountable owner, intended audience, purpose, review cadence, and retirement rule.
- [ ] **USEQ-3203A8EA** — Require original value, factual accuracy, clarity, completeness, and alignment with the user’s likely task.
- [ ] **USEQ-39D8A36D** — Identify authoritative sources and retain citations or evidence for material claims.
- [ ] **USEQ-4B623335** — Distinguish fact, opinion, estimate, advertisement, user content, generated content, and uncertainty.
- [ ] **USEQ-4E963679** — Disclose authorship, editorial responsibility, sponsorship, conflicts, and material automation where relevant.
- [ ] **USEQ-16010080** — Provide publication, review, update, correction, and expiry dates when freshness matters.
- [ ] **USEQ-9D49F492** — Review rapidly changing, high-impact, legal, medical, financial, safety, and security content on a risk-based schedule.
- [ ] **USEQ-B3215B14** — Correct material errors promptly and preserve correction history where user reliance or accountability requires it.
- [ ] **USEQ-35EC77D7** — Do not publish scaled duplicate, paraphrased, scraped, autogenerated, or thin content without meaningful user value and review.
- [ ] **USEQ-88073B7A** — Ensure generated content is reviewed for factuality, unsafe advice, plagiarism, licensing, bias, privacy, and fabricated sources.
- [ ] **USEQ-E4D4D700** — Protect editorial workflows from unauthorized publishing, account takeover, malicious embeds, and content supply-chain compromise.
- [ ] **USEQ-43033DAF** — Use controlled taxonomies, terminology, metadata, and content models so equivalent concepts remain consistent.
- [ ] **USEQ-13A01ED0** — Apply plain-language and content-design principles without removing necessary precision.
- [ ] **USEQ-46372623** — Keep navigation labels, page titles, headings, metadata, and visible content semantically aligned.
- [ ] **USEQ-4A14B2BB** — Manage content duplication intentionally through reuse, canonicalization, consolidation, or audience-specific differentiation.
- [ ] **USEQ-3B3E0164** — Archive or remove obsolete content rather than leaving it authoritative by accident.
- [ ] **USEQ-E3D759B4** — Preserve provenance, approvals, source material, licenses, and change history for regulated or high-trust content.
- [ ] **USEQ-A56DA435** — Use machine-readable provenance or signed content credentials when authenticity and chain of custody materially matter.
- [ ] **USEQ-461FBD9D** — Do not represent provenance metadata, watermarks, or detector output as infallible proof of truth or human authorship.
- [ ] **USEQ-D4854974** — Provide accessible contact, complaint, correction, and takedown routes.

### Migrations, monitoring, incident response, and release gates

- [ ] **USEQ-68F88426** — Create a search and accessibility impact plan before domain, URL, CMS, rendering, design-system, navigation, localization, or content migrations.
- [ ] **USEQ-635F7843** — Inventory and map old URLs, titles, metadata, canonical signals, alternate-language relationships, structured data, content, and accessibility behaviors.
- [ ] **USEQ-539B3D6F** — Preserve or intentionally redirect valuable URLs and external references.
- [ ] **USEQ-C3C380C4** — Test redirects, status codes, canonicalization, sitemaps, robots, noindex, rendering, metadata, structured data, accessibility, and performance before cutover.
- [ ] **USEQ-B2D6EBB1** — Keep rollback or rapid correction paths for broken public access, indexing, navigation, or critical accessibility.
- [ ] **USEQ-C8B3010A** — Monitor crawl failures, indexing changes, canonical conflicts, structured-data errors, manual actions, security issues, and traffic anomalies after release.
- [ ] **USEQ-AA34BD2A** — Monitor critical task accessibility, automated defect trends, support reports, and third-party regressions after release.
- [ ] **USEQ-42E0CA7F** — Annotate monitoring with releases, migrations, experiments, template changes, and crawler-policy changes.
- [ ] **USEQ-061489A9** — Investigate material loss of discoverability or accessibility as a product incident rather than assuming it is normal volatility.
- [ ] **USEQ-EF1F5B4F** — Validate whether an observed search change is caused by technical defects, content changes, policy violations, demand shifts, or external ranking systems before acting.
- [ ] **USEQ-914703A7** — Do not reverse accessibility, security, privacy, or content-quality controls solely to chase rankings.
- [ ] **USEQ-8E4B503A** — Retain pre-migration and post-migration evidence and compare representative critical pages and journeys.
- [ ] **USEQ-22D2A29C** — Define stop, rollback, and escalation criteria before launch.
- [ ] **USEQ-22F07F8C** — Block release when critical content is inaccessible, unintended private content is indexable, intended public content cannot be rendered, or primary routes produce incorrect status or canonical behavior.
- [ ] **USEQ-1363AEE7** — Block release when generated metadata, structured data, or page content materially misrepresents identity, price, availability, authorship, reviews, safety, or eligibility.
- [ ] **USEQ-08B33170** — Require sign-off from accessibility, content, product, security/privacy, and discoverability owners for high-impact migrations.
- [ ] **USEQ-B8E52DF7** — Review monitoring and issue ownership through the full stabilization period.
- [ ] **USEQ-92E6A944** — Convert incidents and recurring defects into shared-component, authoring, test, monitoring, and governance improvements.

### Completion rule

- [ ] **USEQ-FF92E18E** — Every applicable control in this document passes with current evidence, or an authorized exception meets all exception requirements.
- [ ] **USEQ-11FE3AD0** — No open issue in this scope exceeds the approved risk tolerance or violates a mandatory legal, contractual, safety, security, privacy, accessibility, or integrity requirement.
- [ ] **USEQ-A3C013FA** — Owners have defined ongoing monitoring, reassessment triggers, and review dates so that completion does not become a one-time claim.

## Final Gap Closure — Human Factors, Operator Experience, and Usable Security

_Consolidated from `final consolidated corpus/02-human-experience-accessibility-content-seo-internationalization.md#Final Gap Closure — Human Factors, Operator Experience, and Usable Security`; 182 non-duplicative controls._

### Human-factors governance and context of use

- [ ] **USEQ-427EE1A8** — Identify every human role that develops, configures, operates, administers, supports, audits, responds to incidents in, or is materially affected by the system.
- [ ] **USEQ-2A552C91** — Document each role’s goals, authority, competence, working conditions, information needs, constraints, and foreseeable errors.
- [ ] **USEQ-F8D1DFE7** — Study real work rather than relying only on written procedures or stakeholder assumptions.
- [ ] **USEQ-39A892D6** — Include routine work, peak demand, degraded operation, emergencies, maintenance, recovery, and handover conditions in the context-of-use analysis.
- [ ] **USEQ-8E8F5314** — Include remote, mobile, shared-device, low-bandwidth, noisy, low-light, high-stress, and interrupted working conditions where they can occur.
- [ ] **USEQ-F44AD905** — Identify mismatches between work as designed, work as documented, and work as actually performed.
- [ ] **USEQ-E0D7619B** — Treat recurring workarounds as evidence of a design or process problem rather than automatically as individual noncompliance.
- [ ] **USEQ-9A16A5C6** — Assign accountable ownership for human-factors risks and corrective actions.
- [ ] **USEQ-8EC94B9A** — Include representative operators, support staff, administrators, and affected users throughout design and evaluation.
- [ ] **USEQ-3C5F986B** — Use qualified human-factors or usability expertise when error could create material safety, security, privacy, financial, legal, or availability harm.
- [ ] **USEQ-77332F89** — Define measurable human-performance outcomes such as task success, error rate, completion time, workload, recovery, comprehension, and trust calibration.
- [ ] **USEQ-AB0DA53A** — Establish release-blocking thresholds for unacceptable operator error, excessive workload, inaccessible controls, or misleading system feedback.
- [ ] **USEQ-A813B589** — Reassess the context of use when roles, staffing, automation, interfaces, procedures, locations, or risk levels change.
- [ ] **USEQ-BA3769AE** — Preserve evidence of human-factors research, design rationale, evaluations, limitations, and unresolved risks.

### Mental workload, attention, and cognitive demand

- [ ] **USEQ-A664B79E** — Design task demand so that ordinary operation does not require sustained overload or unsafe underload.
- [ ] **USEQ-CADBDD14** — Measure cognitive workload for critical and high-frequency tasks using more than subjective opinion alone where practical.
- [ ] **USEQ-B620FB9B** — Avoid requiring users to remember information that the system can safely display, derive, or retain.
- [ ] **USEQ-73F4A01B** — Keep essential context visible while a person makes or confirms a consequential decision.
- [ ] **USEQ-640548BB** — Minimize unnecessary context switching, mode switching, navigation, re-entry, and reconciliation across tools.
- [ ] **USEQ-E6D98903** — Group information according to the user’s task and decision sequence rather than internal system structure alone.
- [ ] **USEQ-A420F8C0** — Present the most decision-relevant information first without hiding material uncertainty or exceptions.
- [ ] **USEQ-1564BAB3** — Avoid dense displays that force operators to scan irrelevant information during time-critical work.
- [ ] **USEQ-2EB12A3D** — Avoid fragmented workflows that require manual copying between systems when secure integration is feasible.
- [ ] **USEQ-47FFC033** — Make current mode, scope, tenant, environment, identity, and target object continuously apparent for high-impact operations.
- [ ] **USEQ-D97FCDB1** — Distinguish production, testing, simulation, and training environments unmistakably.
- [ ] **USEQ-BE29D663** — Prevent visual similarity between safe and dangerous actions from causing selection errors.
- [ ] **USEQ-8221FFF7** — Provide progressive disclosure without concealing information needed to recognize risk.
- [ ] **USEQ-A23A7EBC** — Preserve user orientation when navigating long or complex workflows.
- [ ] **USEQ-715B7543** — Support interruption and resumption without losing state, intent, provenance, or required checks.
- [ ] **USEQ-E16E8D1B** — Mark stale information and display its source and age when freshness affects decisions.
- [ ] **USEQ-C1794E3F** — Make uncertainty, confidence, data quality, and missing information visible where they affect action.
- [ ] **USEQ-95832A5E** — Avoid presenting estimates, predictions, or inferred states as confirmed facts.
- [ ] **USEQ-36BD3262** — Test critical workflows under realistic time pressure, interruptions, and competing demands.
- [ ] **USEQ-D06DE367** — Monitor whether production workload, queue size, or staffing changes invalidate the evaluated workload assumptions.

### Situation awareness and system-state visibility

- [ ] **USEQ-BE3B3181** — Show what the system is doing, what it has completed, what remains pending, and what failed.
- [ ] **USEQ-ABA0C499** — Distinguish requested, accepted, queued, executing, completed, partially completed, canceled, rolled back, and indeterminate states.
- [ ] **USEQ-DFD2097E** — Make hidden automation and background activity visible when it can affect user decisions.
- [ ] **USEQ-16CF89DB** — Show the scope and consequence of an action before commitment.
- [ ] **USEQ-5160A553** — Provide previews for bulk, destructive, irreversible, cross-tenant, or high-cost changes.
- [ ] **USEQ-B76A180E** — Show dependencies, blockers, and downstream effects that a reasonable operator needs to understand.
- [ ] **USEQ-50AC1C1C** — Indicate whether displayed data is live, cached, replicated, estimated, delayed, or reconciled.
- [ ] **USEQ-86F123A0** — Expose degraded modes and disabled safeguards prominently.
- [ ] **USEQ-1C852350** — Ensure status indicators reflect authoritative system state rather than optimistic client assumptions.
- [ ] **USEQ-AB406DFB** — Prevent success messages from appearing before the authoritative operation succeeds.
- [ ] **USEQ-12D988CB** — Provide a reliable way to inspect the final result of consequential actions.
- [ ] **USEQ-F416163A** — Preserve an understandable history of material actions and state transitions.
- [ ] **USEQ-65FBB13B** — Make concurrent edits, conflicts, locks, and pending approvals visible.
- [ ] **USEQ-2F4D8421** — Communicate failover, fallback, and recovery state without requiring specialist diagnosis.
- [ ] **USEQ-C7710710** — Verify that operators can correctly identify abnormal and unsafe states in evaluation.

### Alarms, alerts, notifications, and interruption design

- [ ] **USEQ-F38882AC** — Define a purpose, owner, severity, audience, response, and expiry for every operational alarm or notification.
- [ ] **USEQ-51B2809E** — Alert only when a recipient can take a meaningful action or needs awareness for a defined decision.
- [ ] **USEQ-93F20D56** — Prioritize alarms by consequence and urgency rather than by the emitting component’s internal severity.
- [ ] **USEQ-4A12FEFB** — Use distinct, consistent, and perceivable signals for materially different urgency levels.
- [ ] **USEQ-02CBA158** — Prevent low-value notifications from obscuring urgent conditions.
- [ ] **USEQ-B62D2C33** — Correlate duplicate symptoms and suppress cascades without hiding independent failures.
- [ ] **USEQ-02E00509** — State what happened, who or what is affected, when it began, current severity, and the recommended next action.
- [ ] **USEQ-B20693F4** — Provide direct access to relevant evidence and runbooks where safe.
- [ ] **USEQ-A762DBD7** — Avoid alarm text that assumes undocumented tribal knowledge.
- [ ] **USEQ-E6D37933** — Preserve critical notifications across channel or device failure.
- [ ] **USEQ-AAB9C8C0** — Ensure acknowledgment does not imply resolution.
- [ ] **USEQ-734401CD** — Escalate unacknowledged or unresolved alarms according to impact and time.
- [ ] **USEQ-39882218** — Prevent indefinite silencing, muting, snoozing, or suppression of critical alarms.
- [ ] **USEQ-AF409BCA** — Record who suppressed an alarm, why, for how long, and under what compensating monitoring.
- [ ] **USEQ-CC514E43** — Test notification delivery, comprehension, and response under realistic conditions.
- [ ] **USEQ-8F6C0C32** — Measure nuisance rate, missed-event rate, response time, and operator workload.
- [ ] **USEQ-3FE97034** — Review alarm floods and repeated dismissals as design defects.
- [ ] **USEQ-16958561** — Provide accessible visual, auditory, and haptic alternatives where appropriate.
- [ ] **USEQ-5C37698B** — Avoid relying on color, sound, motion, or a single channel alone for critical information.
- [ ] **USEQ-AB27A3DF** — Design customer notifications so urgency is accurate and not manipulative.

### Usable security and privacy

- [ ] **USEQ-DACA1500** — Make the secure and privacy-preserving path the easiest ordinary path.
- [ ] **USEQ-C8DCD39E** — Avoid security controls that routinely force unsafe workarounds.
- [ ] **USEQ-0426D010** — Explain security and privacy consequences at the moment a person can act on them.
- [ ] **USEQ-0EA95712** — Use language understandable to the intended audience rather than internal security terminology.
- [ ] **USEQ-0C50C691** — Distinguish recommendations, warnings, hard prohibitions, and legally required choices.
- [ ] **USEQ-456F31EA** — Explain why a requested permission, data item, or high-risk action is necessary.
- [ ] **USEQ-6D1E347D** — Request permissions in context and no earlier or broader than required.
- [ ] **USEQ-7348B9D9** — Show the exact account, organization, tenant, resource, audience, scope, duration, and privilege being granted.
- [ ] **USEQ-A01FE895** — Make privilege escalation, impersonation, delegation, sharing, export, and public exposure unmistakable.
- [ ] **USEQ-6E14F14E** — Provide safe defaults and avoid preselecting broad access or unnecessary data use.
- [ ] **USEQ-93B51E78** — Make revocation, withdrawal, logout, deletion, and recovery easy to find and complete.
- [ ] **USEQ-105BAE9F** — Prevent consent fatigue by removing unnecessary prompts and combining only genuinely related choices.
- [ ] **USEQ-3800D247** — Do not train users to approve unexpected authentication or authorization prompts reflexively.
- [ ] **USEQ-0FD9EE53** — Make authentication prompts resistant to misdirection and phishing where the platform permits.
- [ ] **USEQ-707AB163** — Show enough transaction context in confirmation steps to detect substitution or tampering.
- [ ] **USEQ-084225DB** — Avoid requiring secrets to be transcribed, shared, or stored insecurely.
- [ ] **USEQ-7D94E516** — Support password managers, accessible authenticators, and secure recovery methods.
- [ ] **USEQ-7BFE9065** — Design account recovery to resist social engineering without excluding legitimate users.
- [ ] **USEQ-A15BEB43** — Provide clear explanations and safe remediation when access is denied.
- [ ] **USEQ-4F91139B** — Do not reveal protected existence, membership, identity, or resource details through denial messages.
- [ ] **USEQ-2D3C741B** — Warn before data leaves an expected boundary or is shared with a new audience or provider.
- [ ] **USEQ-C567208A** — Make privacy settings reflect actual backend behavior and downstream processing.
- [ ] **USEQ-DEA38BEB** — Test whether users correctly understand sharing, consent, deletion, retention, and recovery outcomes.
- [ ] **USEQ-9124EC32** — Monitor security-control abandonment, repeated bypass, support burden, and error as usability indicators.

### Human error prevention, tolerance, and recovery

- [ ] **USEQ-438E928C** — Identify foreseeable slips, mistakes, omissions, mode errors, sequence errors, and misinterpretations for critical tasks.
- [ ] **USEQ-91B96F84** — Eliminate error opportunities through design before relying on warnings or training.
- [ ] **USEQ-631BFFE1** — Constrain inputs and actions to valid ranges and states where practical.
- [ ] **USEQ-E2E7D7E8** — Use independent confirmation for actions whose harm is severe and not readily reversible.
- [ ] **USEQ-CE0F0834** — Make confirmations specific to the action rather than generic approval prompts.
- [ ] **USEQ-2356FC80** — Avoid confirmation prompts so frequent that users approve them without reading.
- [ ] **USEQ-6DC3E16D** — Provide undo, rollback, draft, simulation, staging, or reversible commitment where feasible.
- [ ] **USEQ-A73F5710** — Preserve work after validation errors, session expiry, network loss, or recoverable failure.
- [ ] **USEQ-8C2C03F4** — Identify the exact error, its scope, what remains valid, and the safest recovery path.
- [ ] **USEQ-8FBF78F7** — Do not blame the user for system ambiguity, stale state, or conflicting constraints.
- [ ] **USEQ-2BD6C29C** — Prevent retry from duplicating charges, messages, records, or destructive effects.
- [ ] **USEQ-88382116** — Make partial completion and uncertain outcomes explicit.
- [ ] **USEQ-967EE82E** — Provide reconciliation when the system cannot know whether an external action completed.
- [ ] **USEQ-A8F84FD1** — Prevent one mistaken bulk selection from affecting an unbounded population.
- [ ] **USEQ-5B6A2F99** — Use safe limits, staged execution, previews, and sampling for mass operations.
- [ ] **USEQ-577EA46D** — Distinguish recoverable warnings from conditions that must block action.
- [ ] **USEQ-CF95BD3A** — Preserve forensic and audit evidence while still enabling safe correction.
- [ ] **USEQ-A8477083** — Analyze near misses and recovered errors, not only incidents that caused harm.
- [ ] **USEQ-7E1BE30B** — Convert recurring user errors into design improvements, tests, and controls.

### Human–automation interaction and decision support

- [ ] **USEQ-7434557A** — Define which decisions remain human, which are automated, and which require shared control.
- [ ] **USEQ-E4F11AA5** — Assign accountability for outcomes even when automation selects or executes the action.
- [ ] **USEQ-9F72B693** — Make automation boundaries, authority, inputs, limitations, and failure modes understandable to operators.
- [ ] **USEQ-C8E8E618** — Prevent operators from assuming automation is correct merely because it is authoritative or statistically sophisticated.
- [ ] **USEQ-3E6BC74A** — Prevent operators from rejecting reliable automation merely because its rationale is unfamiliar.
- [ ] **USEQ-2CC2C206** — Calibrate trust through accurate performance information, limitations, and feedback.
- [ ] **USEQ-049ECDC6** — Display the evidence, rules, assumptions, and uncertainty needed to review consequential recommendations.
- [ ] **USEQ-B9894D69** — Allow a qualified person to challenge, override, defer, or escalate automation where appropriate.
- [ ] **USEQ-A3AB4FAD** — Record overrides and outcomes without discouraging justified intervention.
- [ ] **USEQ-9C53B1F0** — Detect automation bias, complacency, skill decay, and out-of-the-loop loss of awareness.
- [ ] **USEQ-240CA3B2** — Preserve manual competence and fallback procedures where automation failure would be material.
- [ ] **USEQ-6D6B3F70** — Avoid requiring humans to monitor highly reliable automation continuously for rare events without effective attention support.
- [ ] **USEQ-65264138** — Design handoff from automation to human control with sufficient time, context, and authority.
- [ ] **USEQ-40CA4CDB** — Prevent sudden transfer of control when the person cannot reasonably regain situation awareness.
- [ ] **USEQ-D09F15F2** — Define safe behavior when automation and human instructions conflict.
- [ ] **USEQ-D77DF566** — Test automation interaction across normal, edge, degraded, and adversarial conditions.
- [ ] **USEQ-9BF88D5A** — Monitor whether real-world automation performance or operator behavior departs from evaluated assumptions.

### Staffing, fatigue, shifts, and handovers

- [ ] **USEQ-2754CC69** — Staff critical functions at a level that permits safe operation, review, incident response, leave, and training.
- [ ] **USEQ-D84F078F** — Avoid depending on sustained overtime or chronic interruption to maintain required service levels.
- [ ] **USEQ-E4E634C5** — Identify tasks whose error likelihood rises materially with fatigue, isolation, monotony, or time pressure.
- [ ] **USEQ-0FABCB71** — Set workload and shift practices appropriate to the consequences of error.
- [ ] **USEQ-0DF978DF** — Ensure critical decisions can receive a second qualified review when risk warrants it.
- [ ] **USEQ-1C4167A5** — Provide overlap and structured handover for ongoing incidents, migrations, releases, and high-risk operations.
- [ ] **USEQ-76A1AB5D** — Include current state, recent changes, unresolved alarms, risks, decisions, and next actions in handover records.
- [ ] **USEQ-EE688BD3** — Verify that the receiving person understands and accepts the handover.
- [ ] **USEQ-68B07B41** — Avoid changing too many critical roles simultaneously without continuity controls.
- [ ] **USEQ-2F7EFB21** — Provide backup personnel for every critical operational capability.
- [ ] **USEQ-4D83FB81** — Detect and address unsafe on-call load, alert burden, and repeated sleep disruption.
- [ ] **USEQ-86BFA927** — Enable responders to declare themselves unfit or unavailable without retaliation.
- [ ] **USEQ-CDAF1A2F** — Include contractors, suppliers, and outsourced operators in staffing and handover design.
- [ ] **USEQ-69D2A7A8** — Review human-capacity assumptions before launches, seasonal peaks, migrations, and organizational changes.

### Procedures, training, and competence in use

- [ ] **USEQ-2B3CE832** — Keep procedures aligned with the actual system, interfaces, permissions, and failure modes.
- [ ] **USEQ-090057F0** — Make the authoritative procedure easy to locate during time-critical work.
- [ ] **USEQ-B6DF1544** — Distinguish mandatory steps from optional guidance.
- [ ] **USEQ-10C48526** — Explain the purpose and hazard controlled by critical steps.
- [ ] **USEQ-3B4727FD** — Validate procedures through realistic walkthroughs and simulations.
- [ ] **USEQ-5E37EAFF** — Train personnel on abnormal, degraded, emergency, and recovery conditions, not only routine success paths.
- [ ] **USEQ-A4F91F23** — Assess demonstrated competence rather than attendance alone.
- [ ] **USEQ-F11266F2** — Refresh training when the system, threat, procedure, or role materially changes.
- [ ] **USEQ-0AA2E884** — Provide safe practice environments for high-impact operations.
- [ ] **USEQ-C691A030** — Ensure simulations cannot affect production data or users.
- [ ] **USEQ-9754BE06** — Include communication, escalation, evidence preservation, and customer impact in exercises.
- [ ] **USEQ-C5B812F6** — Capture exercise findings and verify corrective-action effectiveness.
- [ ] **USEQ-56926CF3** — Avoid procedures so complex or lengthy that compliance is unrealistic under actual conditions.
- [ ] **USEQ-141DC9C4** — Provide concise checklists for infrequent high-risk tasks while preserving detailed supporting guidance.
- [ ] **USEQ-BD879838** — Test whether new personnel can perform critical tasks without undocumented assistance.

### Inclusive operator and workforce experience

- [ ] **USEQ-0484AE31** — Make internal and administrative tools meet an explicit accessibility target proportionate to their importance.
- [ ] **USEQ-7E095AAC** — Include workers with disabilities in evaluation of critical tools and procedures.
- [ ] **USEQ-32C2F9AE** — Provide alternatives to audio, visual, fine-motor, memory, speech, and biometric interaction where appropriate.
- [ ] **USEQ-0CA3DF0E** — Ensure emergency and security procedures are accessible.
- [ ] **USEQ-BC3E5A49** — Avoid role requirements that exclude people unnecessarily when the barrier can be removed by design or accommodation.
- [ ] **USEQ-0AC37768** — Support language, literacy, cultural, and domain-knowledge differences among intended operators.
- [ ] **USEQ-2370BDFB** — Use terminology consistently across interfaces, documentation, alerts, and training.
- [ ] **USEQ-58670489** — Avoid humor, idiom, abbreviations, and culturally specific assumptions in safety- or security-critical instructions.
- [ ] **USEQ-0946EC5B** — Provide accessible accommodations without weakening accountability or security.
- [ ] **USEQ-E4E1E40D** — Protect health, disability, and accommodation information from unnecessary disclosure.
- [ ] **USEQ-511286BE** — Evaluate whether productivity metrics penalize accessibility needs, careful review, mentoring, or incident-prevention work.
- [ ] **USEQ-CF9A3C7A** — Provide accessible reporting, escalation, and support channels.

### Human-factors assurance and release gates

- [ ] **USEQ-BB1DA3B6** — Define representative tasks, roles, environments, data, devices, and stress conditions for evaluation.
- [ ] **USEQ-FA348CF7** — Combine observation, task-performance measures, interviews, workload assessment, error analysis, and production evidence.
- [ ] **USEQ-1FA1F50C** — Include novice, experienced, occasional, and substitute users where their behavior differs materially.
- [ ] **USEQ-8E0AE768** — Evaluate critical tasks using the exact production interface and configuration.
- [ ] **USEQ-5B8BCABC** — Record deviations between test conditions and real operation.
- [ ] **USEQ-FE495998** — Do not approve a critical workflow solely because trained designers or developers can complete it.
- [ ] **USEQ-EE93F582** — Treat repeated confusion, unsafe workaround, alarm dismissal, or recovery failure as a release blocker when impact is material.
- [ ] **USEQ-CA8D6F84** — Re-evaluate human performance after major interface, staffing, automation, or process changes.
- [ ] **USEQ-15D04E70** — Monitor task failure, support demand, near misses, overrides, abandonment, and recovery after launch.
- [ ] **USEQ-1364D5A3** — Preserve a human-factors risk register and link unresolved risks to owners, monitoring, and expiry.
- [ ] **USEQ-C01D9ECD** — Require independent review when human error could create catastrophic or systemic harm.
- [ ] **USEQ-33BC0B54** — Include human-factors evidence in production-readiness and assurance cases.

## Standards and source references

- [W3C Web Content Accessibility Guidelines 2.2](https://www.w3.org/TR/WCAG22/)
- [ISO 9241-210:2019 — Human-centred design](https://www.iso.org/standard/77520.html)
- [ISO/IEC 25010:2023 — Product quality model](https://www.iso.org/standard/78176.html)
- [ISO/IEC 25019:2023 — Quality-in-use model](https://www.iso.org/standard/78177.html)
- [ISO/IEC/IEEE 26514:2022 — Information for users](https://www.iso.org/standard/77451.html)
- [ISO/IEC 29100:2024 — Privacy framework](https://www.iso.org/standard/85938.html)
- [RFC 9110 — HTTP Semantics](https://www.rfc-editor.org/rfc/rfc9110)
- [OWASP Application Security Verification Standard 5.0.0](https://owasp.org/www-project-application-security-verification-standard/)
- [ISO/IEC/IEEE 15939:2017 — Measurement process](https://www.iso.org/standard/71197.html)
- [ISO/IEC/IEEE 29119-2:2021 — Test processes](https://www.iso.org/standard/79428.html)
- [ISO/IEC/IEEE 29119-4:2021 — Test techniques](https://www.iso.org/standard/79430.html)
- [OWASP Top 10 — 2025](https://owasp.org/www-project-top-ten/)
- [Google Search Essentials](https://developers.google.com/search/docs/essentials)
- [Google SEO Starter Guide](https://developers.google.com/search/docs/fundamentals/seo-starter-guide)
- [Google structured data general guidelines](https://developers.google.com/search/docs/appearance/structured-data/sd-policies)
- [Google JavaScript SEO basics](https://developers.google.com/search/docs/crawling-indexing/javascript/javascript-seo-basics)
- [Google multilingual and multiregional guidance](https://developers.google.com/search/docs/specialty/international/managing-multi-regional-sites)
- [Google site move guidance](https://developers.google.com/search/docs/crawling-indexing/site-move-with-url-changes)
- [W3C WCAG Evaluation Methodology](https://www.w3.org/WAI/test-evaluate/conformance/wcag-em/)
- [W3C Authoring Tool Accessibility Guidelines 2.0](https://www.w3.org/TR/ATAG20/)
- [W3C WCAG2ICT](https://www.w3.org/TR/wcag2ict-22/)
- [WHATWG HTML Living Standard](https://html.spec.whatwg.org/)
- [RFC 9111 — HTTP Caching](https://www.rfc-editor.org/rfc/rfc9111)
- [Web Vitals](https://web.dev/articles/vitals)

---

[Previous phase](02-product-and-requirements.md) · [Next: Phase 4: Architecture and design](04-architecture-and-design.md)
