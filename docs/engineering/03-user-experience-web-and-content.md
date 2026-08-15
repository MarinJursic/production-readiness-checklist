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
