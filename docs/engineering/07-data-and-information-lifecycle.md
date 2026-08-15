# Data and information lifecycle

_Phase 7 of 16 in the [complete engineering review](00-overview.md)._

Governance, contracts, semantics, modeling, quality, lineage, migrations, analytics, records, and preservation.

Assess every control as **Pass**, **Fail**, **Blocked**, or **Not Applicable**. A Pass needs current evidence for the exact product, revision, artifact, configuration, environment, and data state under review.

## Data Stores, Queues, Caches, Search, and Integrity

_Consolidated from `quality standards/08-data/01-data-stores-queues-caches-search-and-integrity.md`; 21 non-duplicative controls._

### Universal controls

- [ ] **USEQ-2C2F0415** — The data model enforces required fields, uniqueness, referential integrity, valid ranges, and business invariants.
- [ ] **USEQ-EFAFB10D** — Transactions cover operations requiring atomicity.
- [ ] **USEQ-36CAB126** — Locking, concurrency, and conflict-resolution strategies are documented.
- [ ] **USEQ-1B52ED82** — Consumers are idempotent where redelivery is possible.
- [ ] **USEQ-25FC8C04** — Dead-letter queues have monitoring, ownership, retention, inspection, and safe replay procedures.
- [ ] **USEQ-7A4C279B** — Reconciliation detects dropped, duplicated, delayed, corrupted, and inconsistent records.
- [ ] **USEQ-AA39F5AD** — Replication lag is measured and reflected in application behavior.
- [ ] **USEQ-DDCCA31C** — Cache keys include every required identity, tenant, locale, permission, experiment, and content-variant dimension.
- [ ] **USEQ-584B5CAF** — Cache invalidation is correct for security- and correctness-sensitive data.
- [ ] **USEQ-335E7D4B** — Cache failure cannot disclose data or violate invariants.
- [ ] **USEQ-F6712574** — Search indexes enforce access control, deletion, retention, and tenant boundaries.
- [ ] **USEQ-46D419C9** — Character encoding, collation, normalization, case behavior, sorting, and tokenization are intentional.
- [ ] **USEQ-E79CE4F7** — Dates and times use an unambiguous storage representation.
- [ ] **USEQ-4E1895B9** — Data-quality indicators cover completeness, validity, consistency, timeliness, duplication, and lineage.
- [ ] **USEQ-DCF9A997** — Data-repair procedures are reviewed, reversible where possible, and audited.
- [ ] **USEQ-689DD559** — Manual production data changes require authorization, peer review where material, and an audit trail.

### Schema changes and migrations

- [ ] **USEQ-BC1488EE** — Migration duration, locking, load, and resource consumption are measured.
- [ ] **USEQ-6223E816** — Online migrations avoid unacceptable blocking and downtime.
- [ ] **USEQ-1C5AB31F** — Application and schema changes remain backward- and forward-compatible during mixed-version deployment where required.
- [ ] **USEQ-7255C1AD** — Deprecated columns, data, indexes, and compatibility paths are removed only after all consumers have migrated.
- [ ] **USEQ-C449239F** — Migrations are rehearsed using the same automation used for production.

## Imports, Exports, Bulk Operations, and Data Portability

_Consolidated from `quality standards/08-data/02-import-export-bulk-operations-and-portability.md`; 14 non-duplicative controls._

### Universal controls

- [ ] **USEQ-507DCCE5** — Import formats, schemas, limits, encodings, and supported versions are explicit.
- [ ] **USEQ-29F90E73** — Import files and records are validated and bounded.
- [ ] **USEQ-14DD5460** — Imports are idempotent or use robust duplicate detection.
- [ ] **USEQ-C720565E** — Partial imports have clear rollback, retry, checkpoint, and reconciliation behavior.
- [ ] **USEQ-E58DE99C** — Import previews and validation reports are available for high-impact changes.
- [ ] **USEQ-4E113308** — Bulk operations enforce per-record authorization and tenant isolation.
- [ ] **USEQ-3A4139B8** — Bulk operations have rate, concurrency, memory, and downstream-resource limits.
- [ ] **USEQ-94620FDB** — Destructive bulk operations require confirmation, scope preview, authorization, and audit.
- [ ] **USEQ-348958E2** — Exports include only authorized data and respect filters, permissions, residency, retention, and legal holds.
- [ ] **USEQ-8830BDCC** — Export schemas are predictable and versioned.
- [ ] **USEQ-C40BA1BF** — Exports escape active spreadsheet or document content.
- [ ] **USEQ-7FAD3DB9** — Export generation cannot exhaust shared resources.
- [ ] **USEQ-54E0A23C** — Export links and files are protected, expire appropriately, and are deleted according to policy.
- [ ] **USEQ-53208459** — Data-portability formats are documented and machine-readable where required.

## Data Governance, Ownership, and Accountability

_Consolidated from `quality standards/08-data/03-data-governance-ownership-and-accountability.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-DDB8BB36** — Assign a business owner and technical steward to every material data domain and data product.
- [ ] **USEQ-04E6873D** — Define authoritative sources, permitted uses, users, quality obligations, sensitivity, and lifecycle.
- [ ] **USEQ-39926AB5** — Classify data by confidentiality, integrity, availability, privacy, safety, regulatory, and business impact.
- [ ] **USEQ-6BFA579E** — Maintain an inventory of stores, flows, copies, exports, derived data, models, backups, and third-party transfers.
- [ ] **USEQ-47E3AA76** — Define access according to least privilege, purpose, role, tenant, and environment.
- [ ] **USEQ-87C9E971** — Require approval and audit for sensitive production data access and manual changes.
- [ ] **USEQ-01F51D9B** — Define quality contracts between producers and consumers.
- [ ] **USEQ-FA08D468** — Make ownership of shared identifiers, definitions, metrics, and reference data explicit.
- [ ] **USEQ-6E790781** — Prevent unauthorized secondary use, combination, enrichment, or retention.
- [ ] **USEQ-303E3614** — Control data sharing through documented agreements and technical enforcement.
- [ ] **USEQ-BADAB0C6** — Maintain lineage from source through transformations to decisions and reports.
- [ ] **USEQ-94F759EB** — Define retention, deletion, archival, legal hold, and restoration behavior.
- [ ] **USEQ-CBEB4AAF** — Review data risk when introducing new sources, joins, inferences, models, or destinations.
- [ ] **USEQ-084B2C8A** — Make data incidents, quality failures, and rights requests routable to accountable owners.
- [ ] **USEQ-C5054653** — Measure governance by risk and data outcomes rather than catalog completion alone.
- [ ] **USEQ-7B62C473** — Resolve conflicting definitions through an authoritative governance process.
- [ ] **USEQ-4546BD78** — Ensure suppliers and processors meet equivalent data obligations.
- [ ] **USEQ-0775A316** — Remove orphaned stores, copies, exports, accounts, and pipelines.
- [ ] **USEQ-33F0D22A** — Preserve evidence of data decisions and changes.
- [ ] **USEQ-A8D04A5A** — Review governance as products, jurisdictions, users, and purposes evolve.

## Data Architecture and Modeling

_Consolidated from `quality standards/08-data/04-data-architecture-and-modeling.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-5EC68102** — Model domain concepts and relationships independently of a specific storage technology before physical optimization.
- [ ] **USEQ-4907F297** — Use a shared vocabulary and define every material entity, attribute, identifier, status, unit, and relationship.
- [ ] **USEQ-847FF844** — Distinguish natural, surrogate, external, display, tenant, and correlation identifiers.
- [ ] **USEQ-1264F0CD** — Define null, absent, unknown, not-applicable, zero, empty, and deleted semantics explicitly.
- [ ] **USEQ-790E6BAC** — Encode invariants through schema and constraints where enforcement is reliable.
- [ ] **USEQ-F3F98125** — Normalize authoritative data enough to avoid contradictory facts, and denormalize only with ownership and reconciliation.
- [ ] **USEQ-0EB7A9E0** — Define cardinality, optionality, uniqueness, referential, temporal, and lifecycle constraints.
- [ ] **USEQ-BA72D75D** — Choose precision, scale, encoding, collation, normalization, time zone, and unit deliberately.
- [ ] **USEQ-F86428E5** — Represent history and effective time where past truth or auditability matters.
- [ ] **USEQ-ED589523** — Separate immutable facts from mutable projections and current-state summaries.
- [ ] **USEQ-485A3BAD** — Align partitions, tenancy, residency, and encryption boundaries with governance requirements.
- [ ] **USEQ-BFEC9150** — Model access patterns without making one query pattern corrupt the domain model.
- [ ] **USEQ-35CA3885** — Avoid polymorphic or untyped structures that hide incompatible meanings without governance.
- [ ] **USEQ-5AC48987** — Version schemas and define compatibility for producers and consumers.
- [ ] **USEQ-5E48A4D4** — Document ownership and source of truth for every duplicated field.
- [ ] **USEQ-FAA9F327** — Design deletion and legal hold into relationships rather than relying on ad hoc cleanup.
- [ ] **USEQ-2D08184E** — Make bulk import, export, correction, and migration semantics explicit.
- [ ] **USEQ-E76EC4E9** — Validate models with representative data volumes, edge cases, and users.
- [ ] **USEQ-AD2A2D72** — Prevent internal storage layout from becoming an unnecessary public contract.
- [ ] **USEQ-87C063CE** — Review the model when workarounds and repeated joins show a mismatch with the domain.

## Data Quality

_Consolidated from `quality standards/08-data/05-data-quality.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-96E785D6** — Define quality dimensions and thresholds according to each intended use and consequence of error.
- [ ] **USEQ-0EB8C058** — Identify critical data elements and the business rules they must satisfy.
- [ ] **USEQ-31AFDEB2** — Validate at collection and ingestion while preserving raw evidence needed for investigation.
- [ ] **USEQ-EB05E76F** — Distinguish invalid, missing, late, duplicated, contradictory, stale, inferred, and uncertain data.
- [ ] **USEQ-15440CF2** — Record provenance, collection time, effective time, processing time, and quality status where material.
- [ ] **USEQ-008E403E** — Use referential, domain, uniqueness, range, format, temporal, and cross-field checks.
- [ ] **USEQ-26DF5DCE** — Reconcile counts, totals, balances, checksums, and invariants across boundaries.
- [ ] **USEQ-483E3476** — Monitor drift, schema change, source change, late arrival, and silent truncation.
- [ ] **USEQ-AC86E38E** — Segment quality by source, cohort, region, tenant, device, and time when aggregate metrics hide problems.
- [ ] **USEQ-188396D5** — Prevent downstream systems from treating unknown or low-confidence values as authoritative.
- [ ] **USEQ-A4DEB1F3** — Provide quarantine, correction, replay, and backfill procedures.
- [ ] **USEQ-D2553028** — Preserve audit history for material corrections.
- [ ] **USEQ-219C5F63** — Assign owners and response objectives to failed quality rules.
- [ ] **USEQ-5D148E5D** — Test pipelines with malformed, boundary, duplicate, reordered, late, and partial data.
- [ ] **USEQ-544292A2** — Measure false positives and false negatives in automated quality checks.
- [ ] **USEQ-2B4C5C56** — Communicate freshness, completeness, and limitations to consumers.
- [ ] **USEQ-AF818D57** — Prevent remediation from destroying evidence or creating inconsistent copies.
- [ ] **USEQ-59968AE7** — Verify that derived metrics remain semantically valid after source or transformation changes.
- [ ] **USEQ-041F521D** — Use user and operator reports as quality signals.
- [ ] **USEQ-6A36C3D2** — Retire quality rules that no longer correspond to a real information need, and add rules from incidents and escaped defects.

## Metadata, Lineage, Provenance, and Catalogs

_Consolidated from `quality standards/08-data/06-metadata-lineage-provenance-and-catalogs.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-FF3D5CCB** — Catalog material data sets, streams, schemas, metrics, reports, models, features, and exports.
- [ ] **USEQ-E18F00AC** — Record owner, steward, purpose, source, consumers, sensitivity, residency, retention, quality, and support status.
- [ ] **USEQ-9443D242** — Maintain business definitions and technical schemas together with clear authority.
- [ ] **USEQ-A58054A8** — Capture lineage across ingestion, transformation, aggregation, model use, export, and deletion.
- [ ] **USEQ-96D946D3** — Version metadata and retain historical meaning for past reports and incidents.
- [ ] **USEQ-EB438196** — Make lineage granular enough to assess impact without exposing unnecessary sensitive values.
- [ ] **USEQ-0C9F0B2B** — Connect code, configuration, jobs, schemas, tests, dashboards, and incidents to data assets.
- [ ] **USEQ-C4B33E7A** — Detect undocumented stores, fields, pipelines, and transfers.
- [ ] **USEQ-D7D1C856** — Mark deprecated, experimental, untrusted, sampled, synthetic, and derived assets clearly.
- [ ] **USEQ-22AEFF86** — Expose freshness, quality, completeness, and known limitations to consumers.
- [ ] **USEQ-4A84E28E** — Control access to sensitive metadata that could reveal identities, vulnerabilities, or business secrets.
- [ ] **USEQ-820C8588** — Use automated discovery as assistance, not a substitute for accountable validation.
- [ ] **USEQ-CD20611C** — Review lineage after schema, source, transformation, supplier, or model changes.
- [ ] **USEQ-65EE7446** — Ensure deletion, consent, and legal-hold propagation can be traced.
- [ ] **USEQ-129479F9** — Support impact analysis before breaking or retiring an asset.
- [ ] **USEQ-B50D752D** — Prevent multiple conflicting definitions of critical business measures.
- [ ] **USEQ-10B8922B** — Make catalog search and terminology understandable to non-specialist users.
- [ ] **USEQ-0A167C3B** — Define metadata retention and retirement.
- [ ] **USEQ-37837369** — Measure catalog accuracy, coverage, use, stale ownership, and unresolved quality issues.
- [ ] **USEQ-C478C8F8** — Verify metadata against deployed reality periodically.

## Schema Evolution, Migrations, and Data Repair

_Consolidated from `quality standards/08-data/07-schema-evolution-migrations-and-data-repair.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-509E348F** — Classify each change as additive, compatible, breaking, destructive, reparative, or behavioral.
- [ ] **USEQ-43AD1E0F** — Analyze every producer, consumer, query, export, cache, index, report, backup, and integration affected.
- [ ] **USEQ-C8E2FD24** — Use expand-migrate-contract or an equivalent staged approach when versions coexist.
- [ ] **USEQ-DB25A2EB** — Deploy readers tolerant of the new and old representation before writers require the new representation.
- [ ] **USEQ-A3949A80** — Make migrations resumable, idempotent, observable, and bounded.
- [ ] **USEQ-DD77C289** — Test against representative volume, skew, malformed legacy records, and production-like contention.
- [ ] **USEQ-5FDAB7FB** — Estimate duration, lock impact, log growth, replication lag, storage, and failure recovery.
- [ ] **USEQ-D9C07812** — Create a verified recovery point before destructive change.
- [ ] **USEQ-AB1CDFE8** — Define whether rollback is safe; when it is not, provide a tested roll-forward plan.
- [ ] **USEQ-AC3EF432** — Validate preconditions and stop rather than applying a migration to an unexpected state.
- [ ] **USEQ-F2F92B55** — Verify counts, checksums, constraints, totals, samples, and business invariants afterward.
- [ ] **USEQ-47CBE497** — Keep old fields and code until all consumers have migrated and rollback windows close.
- [ ] **USEQ-702C04BB** — Prevent dual-write divergence through transactional patterns or reconciliation.
- [ ] **USEQ-7D7F2085** — Record every repair's scope, reason, approver, code, affected records, and before/after evidence.
- [ ] **USEQ-71D3AA4E** — Use dry runs and impact previews for bulk repair.
- [ ] **USEQ-D9B1C1E4** — Protect privacy, retention, and legal-hold semantics during backfill and copy.
- [ ] **USEQ-6E24DDDE** — Throttle work to preserve service objectives and downstream capacity.
- [ ] **USEQ-EF110B0C** — Monitor errors, lag, locks, saturation, and quality during execution.
- [ ] **USEQ-13513246** — Avoid irreversible manual production edits without review and reproducible scripts.
- [ ] **USEQ-D7F7DCBE** — Retire temporary migration infrastructure, permissions, and compatibility paths after verification.

## Analytics, Business Intelligence, and Decision Data

_Consolidated from `quality standards/08-data/08-analytics-bi-and-decision-data.md`; 20 non-duplicative controls._

### Category-specific universal rules

- [ ] **USEQ-E1AC670A** — Define every decision metric with purpose, formula, population, inclusion, exclusion, time basis, unit, owner, and source.
- [ ] **USEQ-C1D6F926** — Use governed semantic definitions across dashboards and reports.
- [ ] **USEQ-37D5ABE7** — Preserve lineage from source events through transformations to reported values.
- [ ] **USEQ-661381D2** — Validate event generation, deduplication, identity stitching, bot filtering, late data, and backfill behavior.
- [ ] **USEQ-BE717543** — Distinguish operational, financial, experimental, forecast, and exploratory data.
- [ ] **USEQ-3EDB9D35** — Label provisional, sampled, modeled, incomplete, delayed, and corrected values.
- [ ] **USEQ-93C92CBA** — Use appropriate statistical methods and communicate uncertainty, confidence, and practical significance.
- [ ] **USEQ-4DB89BD6** — Prevent dashboards from encouraging decisions outside the metric's valid scope.
- [ ] **USEQ-A3F2C18E** — Segment results where aggregate values can conceal failures or harmed groups.
- [ ] **USEQ-BB25C45C** — Protect personal, sensitive, and commercially confidential information through minimization and access control.
- [ ] **USEQ-13C8956B** — Avoid unnecessary individual-level tracking when aggregate evidence is sufficient.
- [ ] **USEQ-39EDC90D** — Reconcile critical financial and operational metrics with authoritative systems.
- [ ] **USEQ-4989F7C1** — Version transformations and reports so past decisions can be reproduced.
- [ ] **USEQ-380BC207** — Test metric changes and obtain stakeholder approval before replacing established definitions.
- [ ] **USEQ-2EB60C43** — Monitor freshness, completeness, pipeline failure, schema drift, and unexpected distribution changes.
- [ ] **USEQ-4B85F9C2** — Define retention and deletion for raw, transformed, exported, and cached analytical data.
- [ ] **USEQ-1ED5CCAF** — Audit sensitive exports and broad query access.
- [ ] **USEQ-E039CC8E** — Provide a correction and communication process for materially wrong reports.
- [ ] **USEQ-FB1819A1** — Retire unused dashboards and contradictory measures.
- [ ] **USEQ-B55432F8** — Pair dashboards with context, narrative, and responsible interpretation.

## Data Contracts, Semantics, Records, and Digital Preservation Master Checklist

_Consolidated from `gap supplement/06-data-contracts-semantics-records-and-preservation.md`; 246 non-duplicative controls._

### Expanded gap-closure controls

#### Data governance and product accountability

- [ ] **USEQ-C547E075** — Treat every material dataset, event stream, metric, feature, model input, report, archive, record collection, and derived product as an owned product with defined consumers and consequences.
- [ ] **USEQ-CDE52000** — Assign accountable data owner, steward, technical custodian, security owner, privacy owner, records owner, quality owner, and operational contact where applicable.
- [ ] **USEQ-D46B4888** — Maintain an inventory of authoritative sources, replicas, transformations, interfaces, consumers, classifications, locations, jurisdictions, retention classes, and criticality.
- [ ] **USEQ-D1AAA044** — Define the business purpose, permitted uses, prohibited uses, affected parties, decision dependence, expected quality, and retirement criteria for each data product.
- [ ] **USEQ-9B0AB803** — Map legal, regulatory, contractual, privacy, intellectual-property, security, records, sovereignty, localization, archival, and sector obligations.
- [ ] **USEQ-2E6C4E1E** — Establish decision rights for schema, semantics, quality thresholds, access, sharing, correction, retention, deletion, publication, and deprecation.
- [ ] **USEQ-1A7102F9** — Separate data ownership from unrestricted access; ownership creates accountability, not personal possession.
- [ ] **USEQ-2AD17106** — Define authoritative systems of record and prevent multiple sources from claiming incompatible authority without reconciliation rules.
- [ ] **USEQ-7D1A2726** — Classify data by confidentiality, integrity, availability, privacy, safety, financial, legal, evidentiary, and business impact.
- [ ] **USEQ-8A244596** — Document critical data elements and the business rules whose violation would create material harm.
- [ ] **USEQ-73BC563D** — Require data impact analysis for new collection, new fields, changed semantics, new consumers, cross-border movement, new inference, and retirement.
- [ ] **USEQ-1FB06B55** — Establish a data issue, exception, waiver, and escalation process with accountable risk acceptance.
- [ ] **USEQ-AD3273DB** — Review data governance when organizational ownership, product purpose, jurisdiction, architecture, or supplier changes.
- [ ] **USEQ-0197FA6A** — Fund stewardship, metadata, quality, lineage, migration, retention, and preservation work as product capabilities rather than incidental cleanup.

#### Data contracts and producer-consumer agreements

- [ ] **USEQ-D8A791A2** — Use a versioned data contract for material datasets, events, files, APIs, features, reports, and shared tables whose changes can affect independent consumers.
- [ ] **USEQ-C9C31A30** — Identify producer, owner, steward, approved consumers, support contact, lifecycle state, and change authority in the contract.
- [ ] **USEQ-18AD2497** — Define field names, meanings, types, cardinality, nullability, required status, default behavior, constraints, identifiers, units, formats, precision, encoding, and examples.
- [ ] **USEQ-A553302B** — Define record, event, collection, partition, ordering, deduplication, update, deletion, correction, and late-arrival semantics.
- [ ] **USEQ-B5A21EC4** — Define event time, processing time, effective time, observation time, time zone, clock basis, validity interval, and timestamp precision where relevant.
- [ ] **USEQ-31C5E44E** — Define quality objectives for completeness, validity, uniqueness, consistency, accuracy, timeliness, freshness, reconciliation, and acceptable loss.
- [ ] **USEQ-3B112125** — Define delivery objectives for latency, availability, durability, throughput, retention, replay, backfill, recovery, and support.
- [ ] **USEQ-93060CD4** — Define confidentiality, privacy, purpose, region, access, encryption, masking, logging, retention, and deletion requirements.
- [ ] **USEQ-E729AA07** — Define provenance, lineage, source authority, transformation, aggregation, and derivation expectations.
- [ ] **USEQ-8E555230** — Define error, quarantine, dead-letter, rejection, partial acceptance, retry, duplicate, and reconciliation behavior.
- [ ] **USEQ-29044E1F** — Define compatibility promises, supported versions, deprecation period, migration assistance, and termination conditions.
- [ ] **USEQ-01525D88** — Define cost, quota, fair-use, capacity, and high-volume consumer expectations where material.
- [ ] **USEQ-4E491D18** — Represent contracts in machine-readable form where it improves validation without replacing human-readable semantics.
- [ ] **USEQ-AD01F86C** — Store contracts with version control, review, ownership, release history, and immutable identifiers.
- [ ] **USEQ-70D26B9C** — Validate produced data against the contract at the earliest trustworthy boundary.
- [ ] **USEQ-BB5C83FD** — Validate consumer assumptions explicitly rather than relying only on producer-side schema checks.
- [ ] **USEQ-0132ADA9** — Run contract tests in CI and against representative deployed interfaces.
- [ ] **USEQ-63B27806** — Reject, quarantine, or visibly flag contract violations rather than silently coercing values into misleading success.
- [ ] **USEQ-CE3470D3** — Prevent one consumer's undocumented dependency from silently becoming a permanent public contract; discover and resolve such dependencies.
- [ ] **USEQ-231C09BC** — Require a contract change proposal to identify affected consumers, migration path, evidence, rollout, rollback, and observation plan.
- [ ] **USEQ-1C751A87** — Distinguish additive structural compatibility from semantic compatibility; a field can remain structurally present while its meaning changes incompatibly.
- [ ] **USEQ-0F7847A2** — Treat changes to units, timezone, category definitions, enumerations, default values, precision, ordering, frequency, population, methodology, and source as potentially breaking.
- [ ] **USEQ-8156074E** — Support dual publication, adapters, versioned views, or coordinated migration when consumers cannot change atomically.
- [ ] **USEQ-BFA2B85F** — Monitor consumer adoption and retire old contract versions only after evidence that required consumers have migrated.
- [ ] **USEQ-A20E36B7** — Provide a clear process for emergency corrections when continuing bad data is more harmful than preserving compatibility.

#### Semantic modeling and interoperability

- [ ] **USEQ-AA8B7EF7** — Define each domain concept independently from its current storage, API, interface, or organizational representation.
- [ ] **USEQ-BFA7B27F** — Use a controlled vocabulary, glossary, ontology, taxonomy, code set, or semantic model where multiple teams or systems must share meaning.
- [ ] **USEQ-7A78C028** — Assign identifiers that are stable, unique in their namespace, non-recycled, and independent of mutable display attributes.
- [ ] **USEQ-F730C612** — Document identifier scope, issuer, format, lifecycle, merge, split, alias, supersession, and retirement behavior.
- [ ] **USEQ-8EF33366** — Distinguish entity identity from account identity, tenancy, location, version, state, and presentation.
- [ ] **USEQ-EAFEDDAC** — Define units using an accepted system and include unit metadata whenever ambiguity is possible.
- [ ] **USEQ-ED2A5112** — Define currency, exchange-rate source, effective time, precision, rounding, and minor-unit behavior for monetary data.
- [ ] **USEQ-2FBAE1B8** — Define date, time, duration, interval, recurrence, calendar, locale, and daylight-saving semantics explicitly.
- [ ] **USEQ-5D81ECC6** — Use unambiguous machine representations for dates and times and preserve original source timezone or offset when needed for meaning or evidence.
- [ ] **USEQ-336A10C7** — Define whether missing, unknown, not applicable, withheld, redacted, zero, empty, false, and not yet observed are distinct states.
- [ ] **USEQ-1EB23A6A** — Avoid overloaded sentinel values and magic codes.
- [ ] **USEQ-17062C73** — Define code-set owner, version, hierarchy, mappings, deprecations, unknown values, and localization.
- [ ] **USEQ-87876B70** — Preserve source values when normalization or mapping could lose evidence or require later reinterpretation.
- [ ] **USEQ-509B2FD6** — Record semantic mappings and transformation rules between systems and test for information loss.
- [ ] **USEQ-72D57CC1** — Use open, documented, stable formats and vocabularies when interoperability and longevity outweigh proprietary optimization.
- [ ] **USEQ-3F9D0214** — Define character encoding, normalization, collation, language, script, direction, case, transliteration, and comparison behavior.
- [ ] **USEQ-5A23AF38** — Support names, addresses, identifiers, and cultural concepts without assuming one country's format.
- [ ] **USEQ-5504E2F7** — Document aggregation grain, population, denominator, inclusion, exclusion, weighting, and suppression rules for metrics.
- [ ] **USEQ-F72D313D** — Prevent one metric name from referring to different formulas or populations across products without qualification.
- [ ] **USEQ-C349FEF7** — Distinguish measured, reported, inferred, estimated, imputed, predicted, and manually corrected values.
- [ ] **USEQ-D1D271BC** — Attach uncertainty, confidence, error bounds, method, and provenance where values are not exact.
- [ ] **USEQ-FF236B14** — Version semantic definitions independently when meaning can change without a physical schema change.
- [ ] **USEQ-553684CD** — Use a compatibility and mapping strategy when adopting a new vocabulary or standard.

#### Metadata, cataloging, discovery, and FAIR data

- [ ] **USEQ-061DB99D** — Create metadata sufficient for authorized users to find, understand, assess, access, use, cite, and retire each material data asset.
- [ ] **USEQ-62C19DBD** — Include title, description, owner, contact, purpose, source, scope, population, grain, schema, semantics, format, quality, classification, license, rights, location, freshness, retention, lineage, and lifecycle state as applicable.
- [ ] **USEQ-EDEF1598** — Use globally unique and persistent identifiers for durable published data assets where appropriate.
- [ ] **USEQ-8605BA17** — Make metadata searchable through a catalog or equivalent discovery mechanism.
- [ ] **USEQ-4DE4A813** — Keep metadata available even when the data is restricted, while preventing sensitive metadata leakage.
- [ ] **USEQ-6B30EA7B** — Use interoperable metadata vocabularies and machine-readable formats where exchange or automation is intended.
- [ ] **USEQ-E4039352** — Publish access methods, authentication, authorization, protocol, distribution format, and service-level expectations.
- [ ] **USEQ-1D593FAE** — Describe data license, permitted use, attribution, redistribution, derivative work, and ethical constraints.
- [ ] **USEQ-C50F0C7E** — Record spatial, temporal, language, jurisdictional, and subject coverage.
- [ ] **USEQ-8BC37E42** — Document quality measurements and known limitations rather than labeling data simply high quality.
- [ ] **USEQ-A970B216** — Indicate whether data is authoritative, derived, experimental, deprecated, archived, synthetic, anonymized, or sample data.
- [ ] **USEQ-D64CE299** — Maintain metadata through schema, ownership, location, purpose, quality, and lifecycle changes.
- [ ] **USEQ-460A0879** — Prevent stale catalog entries from pointing to missing, unauthorized, or semantically changed assets.
- [ ] **USEQ-A89FAE90** — Link related versions, distributions, schemas, code sets, documentation, provenance, policies, and successor assets.
- [ ] **USEQ-ACC6ADA4** — Support citation using persistent identifiers, version, publication date, publisher, and access date where appropriate.
- [ ] **USEQ-46AAFE3F** — Use discoverable, accessible, interoperable, and reusable principles in balance with privacy, security, ethics, contracts, and legitimate access controls.
- [ ] **USEQ-F4E4BDBF** — Do not make protected data publicly accessible merely to satisfy discoverability or openness goals.
- [ ] **USEQ-343A79D3** — Measure catalog coverage, metadata completeness, search success, ownership freshness, and unresolved consumer questions.

#### Provenance, lineage, and reproducibility

- [ ] **USEQ-DF1A54E5** — Record where data originated, who or what generated it, when it was collected, under what authority, and through which transformations it passed.
- [ ] **USEQ-C16864A9** — Distinguish source provenance, custody provenance, computational lineage, business lineage, and decision lineage.
- [ ] **USEQ-41917408** — Capture lineage across files, streams, databases, APIs, manual changes, spreadsheets, reports, models, caches, indexes, and exports according to risk.
- [ ] **USEQ-48D5EED4** — Record transformation code, query, configuration, parameters, environment, input versions, output version, and execution identity for material derivations.
- [ ] **USEQ-269F8617** — Make lineage granular enough to identify affected outputs when a source, field, transformation, or quality rule fails.
- [ ] **USEQ-D84F4D47** — Preserve links among raw observation, corrected value, derived value, aggregate, report, and decision where accountability requires them.
- [ ] **USEQ-3943C39A** — Track manual overrides and repairs with actor, reason, evidence, before value, after value, scope, and approval.
- [ ] **USEQ-4844D6D4** — Prevent lineage metadata from being editable by the same unrestricted path as the data it attests to when tamper resistance is required.
- [ ] **USEQ-62D8609C** — Use integrity protection and immutable or append-only records for high-assurance provenance.
- [ ] **USEQ-767577C4** — Record deletion, redaction, masking, anonymization, aggregation, and access-controlled transitions in lineage without retaining prohibited content.
- [ ] **USEQ-0BDC155E** — Represent agents, activities, entities, roles, timestamps, generation, derivation, attribution, and invalidation consistently.
- [ ] **USEQ-3E1BB875** — Validate that lineage itself is complete, current, authorized, and understandable.
- [ ] **USEQ-957B905F** — Test impact analysis by tracing a known source change to downstream assets and consumers.
- [ ] **USEQ-442BB18D** — Make reproducible data products rebuildable from declared source versions, code, configuration, and environment within documented tolerances.
- [ ] **USEQ-08E68045** — Retain source snapshots or immutable references when future reconstruction is required and legally permitted.
- [ ] **USEQ-37CC1E72** — Document when external or mutable sources prevent exact reproduction.
- [ ] **USEQ-4F5E910E** — Use lineage during incident response, rights requests, legal hold, migration, deprecation, quality correction, and model investigation.

#### Data quality and observability

- [ ] **USEQ-A360A36C** — Define quality in relation to an explicit use; no dataset is universally fit for every purpose.
- [ ] **USEQ-A419AD6A** — Set measurable expectations for accuracy, completeness, consistency, credibility, currentness, accessibility, compliance, confidentiality, efficiency, precision, traceability, understandability, availability, portability, and recoverability as applicable.
- [ ] **USEQ-8889705F** — Define critical fields, records, populations, relationships, and aggregates that receive stricter controls.
- [ ] **USEQ-95BE25A6** — Validate type, format, domain, range, length, pattern, uniqueness, referential integrity, conditional rules, and cross-field invariants.
- [ ] **USEQ-987FE684** — Validate semantic and business rules that structural schema validation cannot detect.
- [ ] **USEQ-2834E5D4** — Measure freshness, latency, lateness, missing partitions, volume, distribution, cardinality, duplication, nulls, and schema change.
- [ ] **USEQ-7C6AFFC6** — Monitor cross-system reconciliation, balances, counts, totals, checksums, and business invariants.
- [ ] **USEQ-DE28B24B** — Detect sudden and gradual drift in values, categories, units, populations, source mix, and collection methods.
- [ ] **USEQ-9730C66A** — Establish baselines that account for legitimate seasonality, growth, launches, and known events.
- [ ] **USEQ-55B3F90C** — Route bad data to quarantine or a visible degraded state rather than silently publishing misleading results.
- [ ] **USEQ-479F352A** — Define whether consumers should fail closed, use stale data, use partial data, fall back, or continue with warnings for each quality failure.
- [ ] **USEQ-E1F0C88A** — Prevent quality monitoring from relying exclusively on the same faulty source or transformation it checks.
- [ ] **USEQ-8E9CE593** — Provide end-to-end freshness and correctness indicators for critical data journeys.
- [ ] **USEQ-FCC555F6** — Record quality incidents with affected assets, consumers, decisions, time range, severity, containment, correction, and prevention.
- [ ] **USEQ-42ED1555** — Prioritize quality defects by user and business harm rather than the number of failed rows alone.
- [ ] **USEQ-5D5CFB00** — Create ownership and response objectives for data alerts and avoid unowned noisy checks.
- [ ] **USEQ-2937AC0D** — Test alerting, quarantine, backfill, correction, and reconciliation procedures.
- [ ] **USEQ-D3926FEF** — Measure issue recurrence, time to detection, time to containment, time to correction, consumer impact, and prevention completion.
- [ ] **USEQ-A089E927** — Publish known quality limitations to consumers and revoke a fitness claim when evidence no longer supports it.

#### Pipelines, transformations, and reproducible processing

- [ ] **USEQ-0BEE9F84** — Treat ingestion, transformation, orchestration, quality, publication, and deletion logic as production software subject to review, testing, versioning, security, and observability.
- [ ] **USEQ-3A5EF905** — Make pipeline operations idempotent or provide explicit duplicate-detection and reconciliation.
- [ ] **USEQ-EF1C1BE8** — Define event ordering, partitioning, watermark, late-arrival, correction, retraction, replay, and backfill semantics.
- [ ] **USEQ-5EDE4039** — Prevent backfills and replays from duplicating charges, messages, side effects, metrics, or downstream events.
- [ ] **USEQ-4CC06BF6** — Use atomic publication, versioned snapshots, transactional boundaries, or equivalent mechanisms so consumers do not observe half-written products.
- [ ] **USEQ-40982F28** — Record input and output manifests, counts, checksums, partitions, schema versions, and quality results.
- [ ] **USEQ-A1E6753E** — Bound concurrency, memory, storage, network, compute, file, and cost use.
- [ ] **USEQ-F11A318F** — Use checkpoints and resumability for long-running processes without accepting inconsistent partial state.
- [ ] **USEQ-5E9DEA31** — Prevent a failed task from marking a larger workflow successful.
- [ ] **USEQ-3ACA21DD** — Make retries bounded and safe and distinguish transient from permanent failures.
- [ ] **USEQ-AA48781B** — Use dead-letter and quarantine mechanisms with ownership, retention, privacy, inspection, correction, and safe replay.
- [ ] **USEQ-9E0A29BD** — Validate source and destination authorization on every run and prevent service identities from accumulating unnecessary access.
- [ ] **USEQ-CB390717** — Prevent sensitive data from leaking into logs, temporary files, error payloads, samples, metrics, or debugging tools.
- [ ] **USEQ-E4B21F0D** — Test pipeline behavior under malformed input, missing data, duplicate data, delayed data, schema change, dependency failure, quota, and partial outage.
- [ ] **USEQ-FDD640C1** — Test historical replay against the exact transformation and semantic versions appropriate to each period.
- [ ] **USEQ-E967BBE3** — Separate correction of source truth from compensating transformation patches and document both.
- [ ] **USEQ-6FE5FEFF** — Verify that optimization, parallelization, caching, approximate computation, and incremental processing preserve required semantics.
- [ ] **USEQ-88DC1597** — Maintain rollback or roll-forward procedures for pipeline code, schemas, and published data.

#### Master, reference, and shared data

- [ ] **USEQ-D8AD6E3C** — Identify master and reference data whose inconsistency can affect multiple systems or decisions.
- [ ] **USEQ-FE0920E1** — Assign an authoritative owner and stewardship process for each master entity and reference code set.
- [ ] **USEQ-85EC6FC9** — Define creation, matching, deduplication, merge, split, survivorship, correction, versioning, and retirement rules.
- [ ] **USEQ-114C2920** — Prevent automated entity resolution from merging records without calibrated confidence, evidence, and recovery appropriate to impact.
- [ ] **USEQ-9E40D1CF** — Preserve source identifiers and merge history so erroneous consolidation can be reversed.
- [ ] **USEQ-16F85DA9** — Define hierarchy and relationship effective dates and prevent current-state queries from rewriting historical truth.
- [ ] **USEQ-0500410F** — Version reference data and make consumers aware of compatible and retired versions.
- [ ] **USEQ-4F9A4886** — Handle unknown, pending, other, deprecated, and unmapped values explicitly.
- [ ] **USEQ-2913F1EE** — Reconcile replicated master data and monitor propagation delay.
- [ ] **USEQ-30E0E329** — Prevent one tenant, region, business unit, or supplier from modifying globally shared reference data without authorized governance.
- [ ] **USEQ-242D0174** — Provide bulk correction and migration tooling with preview, approval, audit, rate limits, rollback, and reconciliation.

#### Access, privacy, security, and ethical data use

- [ ] **USEQ-7D80E96A** — Enforce least privilege, purpose limitation, tenant isolation, and separation of duties for collection, query, export, correction, administration, and deletion.
- [ ] **USEQ-1B85AD0F** — Authorize access at dataset, row, column, field, object, purpose, region, and time boundaries as required by risk.
- [ ] **USEQ-8F6092B7** — Prevent derived data, aggregates, embeddings, logs, caches, search indexes, samples, and exports from bypassing source restrictions.
- [ ] **USEQ-31E7AFA5** — Use encryption, tokenization, masking, pseudonymization, confidential computing, secure enclaves, or equivalent controls where justified by the threat model.
- [ ] **USEQ-8BF130CE** — Do not present masking as anonymization when values remain linkable or reversible.
- [ ] **USEQ-031BFAEA** — Assess reidentification and linkage risk for aggregates, sparse datasets, location, time series, and released statistics.
- [ ] **USEQ-AB5D4632** — Apply disclosure controls such as minimum group size, suppression, perturbation, query limits, and review where public or broad analytical access creates risk.
- [ ] **USEQ-8E63E088** — Prevent differencing, repeated-query, and composition attacks against protected aggregates.
- [ ] **USEQ-5DCED033** — Log and review sensitive reads, bulk exports, privilege changes, overrides, and unusual queries.
- [ ] **USEQ-F680A8D1** — Limit analyst and support access to production data and provide safer governed workspaces or synthetic data.
- [ ] **USEQ-1C2BE353** — Ensure purpose, consent, legal basis, contractual restriction, and user preference travel with data where decisions depend on them.
- [ ] **USEQ-6F35F15A** — Prevent discriminatory, manipulative, exploitative, or scientifically unsupported secondary use.
- [ ] **USEQ-07FE40D4** — Review data sharing and combination for harms not present in either source alone.
- [ ] **USEQ-565ACD95** — Define and enforce data residency, localization, cross-border transfer, and subprocessor requirements.
- [ ] **USEQ-24385111** — Ensure rights requests, corrections, restrictions, objections, deletion, and portability reach all relevant copies and consumers.
- [ ] **USEQ-16A208D6** — Protect metadata, lineage, catalog, and quality systems because they can reveal sensitive structure and activity.

#### Records management and evidentiary integrity

- [ ] **USEQ-9F7F3524** — Identify information that constitutes an official, contractual, financial, legal, safety, security, clinical, administrative, or business record.
- [ ] **USEQ-D92A0B88** — Define which events and documents must be captured as records and at what point they become authoritative.
- [ ] **USEQ-EAC18E35** — Maintain records that are authentic, reliable, integral, usable, and connected to their business context.
- [ ] **USEQ-781B65B7** — Capture creator, authority, date, business process, classification, version, relationships, and disposition metadata.
- [ ] **USEQ-5DEED201** — Prevent ordinary editing from silently rewriting an authoritative record; preserve supersession and correction history.
- [ ] **USEQ-B718E81C** — Use append-only, versioned, signed, timestamped, witnessed, or otherwise tamper-evident controls where evidentiary risk warrants them.
- [ ] **USEQ-AA7CF9C6** — Protect records from unauthorized access, alteration, deletion, corruption, format loss, and context loss.
- [ ] **USEQ-0BA61237** — Define file plan, classification, naming, aggregation, access, retention trigger, retention period, disposition authority, and final action.
- [ ] **USEQ-808F7EA8** — Base retention on legal, regulatory, contractual, operational, historical, scientific, user, and privacy requirements rather than unlimited convenience.
- [ ] **USEQ-AD5F5267** — Suspend normal disposition promptly when a legal hold, investigation, audit, incident, or preservation order applies.
- [ ] **USEQ-2092A320** — Document hold scope, authority, custodians, systems, dates, preservation actions, release, and conflicts with deletion duties.
- [ ] **USEQ-284AE7F4** — Prevent backup rotation, lifecycle automation, user deletion, migration, or decommissioning from destroying held records.
- [ ] **USEQ-79B32491** — Apply defensible disposition with authorization, evidence, and verification that all required copies and derivatives were addressed.
- [ ] **USEQ-E7DC518F** — Distinguish correction from deletion and preserve the appropriate audit relationship.
- [ ] **USEQ-6B0C5547** — Ensure records remain searchable, retrievable, readable, understandable, and exportable for the required period.
- [ ] **USEQ-D00A49D8** — Test retrieval by representative request, date range, actor, transaction, and legal matter.
- [ ] **USEQ-9E095079** — Control exports so metadata, ordering, attachments, relationships, signatures, and context remain intact.
- [ ] **USEQ-01C4B98F** — Document custody transfer and chain of custody for evidence and high-assurance records.
- [ ] **USEQ-5101BEAA** — Review record systems after schema, format, supplier, platform, encryption, identity, and retention changes.

#### Digital preservation and long-term access

- [ ] **USEQ-E8B2C973** — Identify digital information whose required life exceeds the expected life of its current format, application, encryption, key, storage medium, vendor, or organization.
- [ ] **USEQ-B59A351D** — Define the designated community and the knowledge needed to understand preserved information independently of the original creators.
- [ ] **USEQ-17F1226D** — Preserve content together with representation information, provenance, context, fixity, rights, identifiers, and packaging information.
- [ ] **USEQ-2DAEA0AE** — Distinguish bit preservation from preservation of meaning, behavior, appearance, executability, and evidentiary validity.
- [ ] **USEQ-7EB59709** — Use checksums or stronger integrity evidence at ingest, storage, replication, transfer, audit, and access.
- [ ] **USEQ-94877209** — Store independent copies across failure domains and periodically verify every copy rather than assuming replication equals integrity.
- [ ] **USEQ-AA74E24D** — Monitor media age, error rates, format obsolescence, software support, hardware availability, vendor viability, cryptographic strength, and access dependencies.
- [ ] **USEQ-719E4479** — Plan and test media refresh, replication, format migration, normalization, emulation, encapsulation, or documentation strategies.
- [ ] **USEQ-8D51D59D** — Validate that migration preserves significant properties, metadata, relationships, timestamps, signatures, accessibility, and usability.
- [ ] **USEQ-1585B6D6** — Retain original bitstreams when lawful and useful, alongside normalized or migrated versions.
- [ ] **USEQ-F8D49065** — Record every preservation action and maintain lineage among originals, derivatives, migrations, and access copies.
- [ ] **USEQ-CF5DB4C0** — Preserve software, configuration, schemas, code sets, fonts, codecs, keys, licenses, documentation, and execution environment when required to interpret or reproduce content.
- [ ] **USEQ-9AF15653** — Plan for encrypted content whose future key, algorithm, certificate, or identity infrastructure may disappear.
- [ ] **USEQ-3C15ABE7** — Use archival information packages and transfer manifests appropriate to the repository and risk.
- [ ] **USEQ-FBA075B3** — Validate ingest completeness and reject or quarantine unsupported, corrupt, encrypted-without-key, or unidentified content.
- [ ] **USEQ-DA245FBE** — Define access controls, privacy restrictions, embargoes, redactions, and lawful disclosure for the preservation period.
- [ ] **USEQ-555290A4** — Make preserved content discoverable without exposing protected metadata.
- [ ] **USEQ-3A0E99D4** — Test recovery and intelligibility using people and systems not involved in original creation.
- [ ] **USEQ-E677E927** — Ensure succession, custody transfer, provider exit, organizational closure, and funding risk do not orphan preserved information.
- [ ] **USEQ-88157CEC** — Review preservation plans periodically and before any platform, supplier, format, or cryptographic transition.

#### Backup, archive, cache, and record distinctions

- [ ] **USEQ-808E256B** — Define separately the purposes of operational backup, disaster recovery, high availability, cache, analytical replica, legal record, business archive, and long-term preservation.
- [ ] **USEQ-CAC49AE0** — Do not treat a backup as a searchable records archive or preservation repository unless it meets the corresponding requirements.
- [ ] **USEQ-A6AC6CC0** — Do not treat replication as backup when corruption, deletion, or compromise propagates to every replica.
- [ ] **USEQ-7919C039** — Do not treat a cache or search index as authoritative source truth unless explicitly designed and governed as one.
- [ ] **USEQ-8E523AC4** — Ensure retention, deletion, legal hold, access, encryption, restoration, and verification rules match each copy's purpose.
- [ ] **USEQ-ACA98CFF** — Track where data can persist after deletion from the primary system and document the permitted expiration path.
- [ ] **USEQ-C6373244** — Prevent disaster recovery from reintroducing revoked credentials, deleted records, invalid reference data, obsolete schemas, or superseded policies without detection.
- [ ] **USEQ-87DCD770** — Test restoration at data-product and business-process level, not only storage-volume level.
- [ ] **USEQ-4D3582CB** — Reconcile restored data with external systems, later transactions, and irreversible side effects.

#### Export, portability, publication, and interchange

- [ ] **USEQ-6CDDC958** — Provide documented, stable, accessible, and secure export formats appropriate to the user's and consumer's legitimate purpose.
- [ ] **USEQ-BB807B59** — Use open or widely implementable formats for long-lived or interoperable exports where practical.
- [ ] **USEQ-A1BBEC55** — Include schema, semantics, units, timezone, encoding, identifiers, provenance, version, and quality information needed to interpret an export.
- [ ] **USEQ-0554B993** — Preserve relationships, attachments, order, history, and context rather than exporting isolated values that cannot be understood.
- [ ] **USEQ-75A76A46** — Validate export completeness, authorization, tenant isolation, privacy, formula injection, active content, size, and resumability.
- [ ] **USEQ-15993798** — Use deterministic ordering or manifests and checksums where recipients must verify completeness.
- [ ] **USEQ-D766C812** — Define incremental, full, point-in-time, and corrected export semantics.
- [ ] **USEQ-6C4053B9** — Support large exports asynchronously with progress, expiry, revocation, audit, and safe retry.
- [ ] **USEQ-A1908B5F** — Prevent export links and files from remaining accessible beyond their authorized purpose.
- [ ] **USEQ-AB993B12** — Provide machine-readable APIs or bulk interfaces where repeated manual export would create risk and operational burden.
- [ ] **USEQ-F471F071** — Test import of exported data or provide a documented independent validation path when portability is claimed.
- [ ] **USEQ-0140D0B4** — Publish datasets only after rights, privacy, security, ethics, accessibility, quality, licensing, metadata, and reidentification review.
- [ ] **USEQ-0A2C8A1D** — Version published data and preserve access to cited versions or clearly record withdrawal.
- [ ] **USEQ-B2FD73C8** — Communicate corrections and withdrawals to known consumers when bad data can cause material harm.

#### Data lifecycle, deprecation, and retirement

- [ ] **USEQ-551FD6F6** — Define creation, active use, sharing, replication, archival, retention, preservation, disposition, and retirement states.
- [ ] **USEQ-758BF94D** — Stop collecting data when the approved purpose ends unless a new authorized purpose is established.
- [ ] **USEQ-7FA64563** — Identify every producer, consumer, replica, derivative, report, model, integration, and retained copy before retirement.
- [ ] **USEQ-D6334B23** — Provide deprecation notice, migration guidance, compatibility period, and owner support proportionate to consumer dependence.
- [ ] **USEQ-DB720526** — Prevent new consumers from adopting a deprecated asset without an approved exception.
- [ ] **USEQ-CA5E222A** — Freeze or preserve a final authoritative version when legal, scientific, contractual, audit, or historical value requires it.
- [ ] **USEQ-73B2C81A** — Delete or irreversibly deidentify data that no longer has an authorized purpose, subject to holds and preservation duties.
- [ ] **USEQ-5A31A9D2** — Verify deletion across primary stores, replicas, caches, indexes, queues, temporary files, analytics, providers, and access credentials.
- [ ] **USEQ-0C7088E7** — Record residual presence in immutable backups and the controls preventing ordinary use until expiry.
- [ ] **USEQ-7B1C10E9** — Revoke service accounts, keys, shares, catalog entries, schedules, webhooks, and automation associated with retired data.
- [ ] **USEQ-2740E000** — Monitor for continued reads or writes after retirement and investigate undeclared consumers.
- [ ] **USEQ-A9B7A92F** — Preserve lineage from retired assets to successors and document semantic differences.
- [ ] **USEQ-80753EC8** — Conduct a post-retirement review to verify cost, risk, rights, contracts, records, and consumer obligations were closed.

#### Data release blockers

- [ ] **USEQ-D7F7986D** — Do not publish or depend on a material data product whose owner, purpose, source, semantics, contract, quality expectations, and consumers are unknown.
- [ ] **USEQ-E4A69C3C** — Do not make a structurally compatible change whose semantic impact on consumers has not been assessed.
- [ ] **USEQ-D7668FC7** — Do not silently coerce, drop, impute, merge, normalize, or reinterpret invalid data when doing so can conceal material error.
- [ ] **USEQ-A7E04443** — Do not proceed with a migration, backfill, replay, or repair without representative rehearsal, observability, reconciliation, and a recovery strategy.
- [ ] **USEQ-FD112BB9** — Do not use data for a new consequential purpose without verifying authority, fitness, fairness, privacy, and downstream impact.
- [ ] **USEQ-B0E18CAC** — Do not release a public or broadly shared dataset with unresolved rights, confidentiality, reidentification, provenance, quality, or licensing risk.
- [ ] **USEQ-396303A1** — Do not delete data subject to a valid hold, records duty, user promise, or preservation obligation.
- [ ] **USEQ-4AB86041** — Do not retain data indefinitely merely because future use is conceivable.
- [ ] **USEQ-319EB9A9** — Do not claim archival or preservation capability when only untested backup copies exist.
- [ ] **USEQ-DAB12EEC** — Do not retire a source until required consumers, records, lineage, exports, holds, and successor mappings have been addressed.

## Standards and source references

- [ISO/IEC 25012:2008 — Data quality model](https://www.iso.org/standard/35736.html)
- [ISO/IEC 25024:2015 — Measurement of data quality](https://www.iso.org/standard/35749.html)
- [ISO/IEC 25010:2023 — Product quality model](https://www.iso.org/standard/78176.html)
- [ISO/IEC 29100:2024 — Privacy framework](https://www.iso.org/standard/85938.html)
- [ISO/IEC 27701:2025 — Privacy information management systems](https://www.iso.org/standard/85819.html)
- [ISO/IEC 38500:2024 — Governance of IT](https://www.iso.org/standard/81684.html)
- [ISO/IEC/IEEE 42010:2022 — Architecture description](https://www.iso.org/standard/74393.html)
- [ISO/IEC/IEEE 15939:2017 — Measurement process](https://www.iso.org/standard/71197.html)
- [ISO/IEC/IEEE 15289:2019 — Life-cycle information items](https://www.iso.org/standard/74909.html)
- [ISO/IEC/IEEE 12207:2026 — Software life cycle processes](https://www.iso.org/standard/90219.html)
- [ISO/IEC/IEEE 29119-2:2021 — Test processes](https://www.iso.org/standard/79428.html)
- [W3C Data on the Web Best Practices](https://www.w3.org/TR/dwbp/)
- [W3C PROV-O: The PROV Ontology](https://www.w3.org/TR/prov-o/)
- [W3C Data Catalog Vocabulary 3](https://www.w3.org/TR/vocab-dcat-3/)
- [FAIR Guiding Principles](https://www.go-fair.org/fair-principles/)
- [Open Data Contract Standard](https://bitol-io.github.io/open-data-contract-standard/latest/)
- [ISO 15489-1:2016 — Records management](https://www.iso.org/standard/62542.html)
- [ISO 14721:2025 — Open archival information system reference model](https://www.iso.org/standard/87471.html)

---

[Previous phase](06-application-services-and-apis.md) · [Next: Phase 8: Security and cryptography](08-security-and-cryptography.md)
