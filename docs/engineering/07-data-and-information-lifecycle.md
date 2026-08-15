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

## Data Contracts, Information Governance, Records, and Digital Preservation

_Consolidated from `final consolidated corpus/04-data-privacy-information-governance-records-analytics.md#Data Contracts, Information Governance, Records, and Digital Preservation`; 255 non-duplicative controls._

### Data and information governance

- [ ] **USEQ-95DD179A** — Inventory structured, semi-structured, unstructured, operational, analytical, event, document, media, model, log, and archival information assets.
- [ ] **USEQ-167FE265** — Assign an accountable business or mission owner and an operational steward to every material data domain or information collection.
- [ ] **USEQ-02ABB799** — Define decision rights for creation, access, definition, quality, sharing, correction, retention, legal hold, archival, and disposal.
- [ ] **USEQ-08751C45** — Classify data by confidentiality, integrity, availability, privacy, safety, regulatory, contractual, and evidentiary impact.
- [ ] **USEQ-9C040A43** — Document authorized purposes and prohibited uses.
- [ ] **USEQ-92923736** — Distinguish source records, authoritative master data, derived data, cached data, replicas, indexes, reports, and convenience copies.
- [ ] **USEQ-874AE34A** — Define the system of record and source of truth for each material concept and state.
- [ ] **USEQ-CBF98788** — Identify all processors, recipients, consumers, regions, storage locations, and transfer mechanisms.
- [ ] **USEQ-C12D6B1F** — Maintain a glossary of business terms with owner, definition, scope, examples, exclusions, and relationships.
- [ ] **USEQ-C06600A4** — Resolve conflicting definitions rather than allowing systems to use the same term for different concepts silently.
- [ ] **USEQ-2AE396F3** — Apply governance proportionately to user harm, decision impact, scale, sensitivity, and irreversibility.
- [ ] **USEQ-6A11EAE1** — Include data obligations in product requirements, architecture, supplier contracts, incident plans, and decommissioning.
- [ ] **USEQ-F4A19FAB** — Monitor whether governance rules are implemented in systems and workflows rather than existing only as policy.
- [ ] **USEQ-AC2A241F** — Review ownership and classification after acquisitions, reorganizations, new purposes, new regions, or material model changes.

### Data contracts and producer-consumer agreements

- [ ] **USEQ-51A0D545** — Define an explicit contract for every material data product, event stream, file exchange, analytical table, report feed, or shared schema.
- [ ] **USEQ-12270F27** — Identify producer, owner, consumers, intended purposes, prohibited uses, support channel, and lifecycle state.
- [ ] **USEQ-A05CA791** — Define schema, field meaning, type, units, encoding, nullability, identifiers, enumerations, constraints, and relationships.
- [ ] **USEQ-E850D3BC** — Define event semantics, trigger, ordering, delivery, duplication, correction, replay, and deletion behavior.
- [ ] **USEQ-B94AD2D5** — Define freshness, completeness, accuracy, consistency, uniqueness, latency, availability, retention, and volume expectations where relevant.
- [ ] **USEQ-043E2F00** — Define privacy classification, access model, geographic restrictions, aggregation rules, and approved downstream uses.
- [ ] **USEQ-9342B12A** — Define compatibility rules for adding, changing, deprecating, renaming, or removing fields and values.
- [ ] **USEQ-354DB2A1** — Define how unknown fields, unknown enumeration values, duplicates, malformed records, and late data are handled.
- [ ] **USEQ-32210956** — Define effective dates and temporal semantics for slowly changing or correction-prone data.
- [ ] **USEQ-805AD65B** — Provide representative valid, invalid, boundary, and evolution examples.
- [ ] **USEQ-6486D18C** — Publish machine-readable schemas and validation rules where practical.
- [ ] **USEQ-89DE0F18** — Version contracts independently of implementation and retain change history.
- [ ] **USEQ-691F686B** — Notify actual consumers of material changes and provide migration windows and tooling.
- [ ] **USEQ-CBC5C76C** — Use producer, consumer, and compatibility tests to verify the contract before release.
- [ ] **USEQ-888341EF** — Detect undocumented consumers before destructive changes.
- [ ] **USEQ-FB4915D7** — Keep contract ownership and support active for as long as consumers are expected to rely on the data.

### Metadata, catalog, glossary, and semantic consistency

- [ ] **USEQ-EF917148** — Catalog every material data set, stream, report, document collection, model input, and model output.
- [ ] **USEQ-36B7A068** — Use stable identifiers for cataloged assets across rename, relocation, storage migration, and platform change.
- [ ] **USEQ-4B9EBD98** — Record title, description, owner, steward, classification, source, schema, lineage, freshness, quality, access, retention, license, and lifecycle status.
- [ ] **USEQ-6B5E4D4F** — Distinguish metadata about the asset from data contained within the asset.
- [ ] **USEQ-1185F811** — Use controlled vocabularies and reference data for repeated concepts where consistency matters.
- [ ] **USEQ-F5B1C025** — Align field and entity names with the approved glossary or document justified domain-specific distinctions.
- [ ] **USEQ-95AD022D** — Record semantic versioning and effective dates for definitions and classifications.
- [ ] **USEQ-79D7B7C5** — Make metadata understandable to humans and processable by machines where interoperability requires it.
- [ ] **USEQ-81D23683** — Use a metadata registry model such as ISO/IEC 11179 principles when consistent data-element definition and registration are material.
- [ ] **USEQ-EA1BD426** — Use interoperable catalog vocabularies such as DCAT when publishing or federating data catalogs on the Web.
- [ ] **USEQ-B5964E35** — Validate catalog metadata for completeness, correctness, links, access, and freshness.
- [ ] **USEQ-E686AE0F** — Prevent catalog descriptions from exposing sensitive data, credentials, internal topology, or confidential business context.
- [ ] **USEQ-CD74884E** — Archive or retire catalog entries deliberately and preserve historical identity where required.
- [ ] **USEQ-1EAC6711** — Monitor metadata drift between catalog, schema, code, storage, and actual observed data.

### Lineage, provenance, and reproducibility

- [ ] **USEQ-26BD4D59** — Record how material data was created, collected, transformed, joined, filtered, corrected, aggregated, labeled, and published.
- [ ] **USEQ-7BB67AA3** — Identify the agents, systems, versions, inputs, activities, timestamps, parameters, and approvals involved.
- [ ] **USEQ-8B4E6080** — Preserve lineage across batch, streaming, manual, spreadsheet, model, and third-party transformations.
- [ ] **USEQ-EB9853B3** — Distinguish observed facts, user assertions, inferred values, estimates, predictions, and generated content.
- [ ] **USEQ-C99DC1D2** — Record the original source and chain of custody for data used in high-impact decisions or evidence.
- [ ] **USEQ-B8F0AAB7** — Link derived data to source versions and transformation versions sufficiently to reproduce or explain it.
- [ ] **USEQ-C4AF5B57** — Capture code, query, model, configuration, environment, and dependency versions needed for reproducibility.
- [ ] **USEQ-6CCF8BFD** — Use stable provenance concepts and interoperable representations, such as W3C PROV, where cross-system exchange matters.
- [ ] **USEQ-04F72893** — Protect provenance records against unauthorized alteration and deletion.
- [ ] **USEQ-BF2AC04E** — Do not claim provenance proves truth; use it to support assessment of origin, process, and accountability.
- [ ] **USEQ-C6567EE2** — Expose lineage to authorized consumers so impact analysis and trust decisions can be made.
- [ ] **USEQ-DE4F354D** — Detect broken or incomplete lineage for critical reports, models, and regulatory outputs.
- [ ] **USEQ-456A5FAE** — Retain provenance for at least as long as the associated data must be trusted, explained, reproduced, or audited.
- [ ] **USEQ-147A0E45** — Include external supplier and manually supplied inputs in lineage rather than treating them as originless.

### Data quality requirements and measurement

- [ ] **USEQ-F38AF100** — Define quality in relation to intended use rather than as an abstract score.
- [ ] **USEQ-5F1BBBE9** — Specify required accuracy, completeness, consistency, uniqueness, validity, timeliness, freshness, precision, traceability, and accessibility as applicable.
- [ ] **USEQ-AA2AD283** — Define each quality rule, calculation, scope, threshold, owner, frequency, and decision consequence.
- [ ] **USEQ-C8612C88** — Distinguish source defects, transformation defects, delayed data, legitimate exceptions, and unknown quality.
- [ ] **USEQ-2FEFB2A0** — Use both field-level and cross-record business-rule validation.
- [ ] **USEQ-ECD94E0D** — Validate referential integrity, cardinality, totals, balances, ranges, temporal order, and state transitions.
- [ ] **USEQ-1643B475** — Use reconciliations against independent sources for balances, payments, inventory, entitlements, counts, and other critical aggregates.
- [ ] **USEQ-D2EFEAEA** — Profile distributions and detect unexpected shifts, missing categories, outliers, and sudden sparsity.
- [ ] **USEQ-5FD46969** — Account for selection, measurement, survivorship, reporting, sampling, label, and historical bias.
- [ ] **USEQ-F3FA2A0C** — Keep quality measurements segmented by source, cohort, region, tenant, time, and workflow where aggregation would hide harm.
- [ ] **USEQ-81D9F64E** — Do not mask serious defects through averages or by excluding failed records from denominators.
- [ ] **USEQ-12FFD88F** — Publish quality limitations and fitness-for-use guidance to consumers.
- [ ] **USEQ-E8DA7B35** — Prevent consumers from interpreting unverified or stale data as authoritative.
- [ ] **USEQ-D50F764D** — Track remediation and recurrence, and repair the upstream process where possible.

### Data observability and operational monitoring

- [ ] **USEQ-6CA16C56** — Monitor arrival, freshness, volume, schema, distribution, nulls, duplicates, quality rules, lineage, and processing outcomes for critical data flows.
- [ ] **USEQ-2D2D929C** — Monitor both technical pipeline success and semantic correctness.
- [ ] **USEQ-27E6480B** — Detect jobs that succeed technically while producing empty, stale, partial, duplicated, or implausible results.
- [ ] **USEQ-1419EF4A** — Correlate data incidents with source, transformation, deployment, configuration, schema, and supplier changes.
- [ ] **USEQ-0C14AC39** — Use end-to-end lineage to identify affected downstream assets and decisions.
- [ ] **USEQ-D84A4747** — Alert on user or business impact rather than every benign variation.
- [ ] **USEQ-DCBA84F8** — Provide runbooks for late data, schema breakage, corrupt partitions, duplication, source outages, and failed reconciliation.
- [ ] **USEQ-F735026B** — Define who can stop publication, quarantine data, replay processing, correct records, and notify consumers.
- [ ] **USEQ-FA240AC7** — Keep monitoring independent enough to detect failure of the primary pipeline.
- [ ] **USEQ-71AFC1FD** — Preserve raw evidence needed to investigate anomalies, subject to privacy and retention constraints.
- [ ] **USEQ-56FC3173** — Monitor quality-rule drift and stale thresholds.
- [ ] **USEQ-496E1CD9** — Annotate data dashboards with deployments, backfills, corrections, and methodology changes.
- [ ] **USEQ-57DE4607** — Treat silent data corruption and materially wrong high-impact outputs as incidents.
- [ ] **USEQ-49EBCD42** — Communicate incident scope, affected periods, assets, consumers, decisions, and correction status.

### Master data, reference data, identity, and identifiers

- [ ] **USEQ-A833E0CF** — Define authoritative sources and stewardship for shared entities, code sets, taxonomies, and reference values.
- [ ] **USEQ-2F541F2A** — Use stable identifiers that do not encode mutable personal or business meaning unless necessary.
- [ ] **USEQ-C4FF2BD1** — Distinguish internal identifiers, public identifiers, natural keys, surrogate keys, and display labels.
- [ ] **USEQ-F709F56A** — Do not use easily guessed identifiers as authorization controls.
- [ ] **USEQ-570E35A2** — Define matching, deduplication, merge, split, and survivorship rules.
- [ ] **USEQ-9C69756D** — Preserve source identifiers and merge history for audit and reversal.
- [ ] **USEQ-CCB4B27E** — Protect against incorrect entity resolution and cross-person or cross-tenant merging.
- [ ] **USEQ-C1CE528F** — Allow correction of names, identities, classifications, and relationships without corrupting historical records.
- [ ] **USEQ-B8F1BBFD** — Version reference data and record effective dates and retirement dates.
- [ ] **USEQ-2953795E** — Ensure producers and consumers agree on code-set meaning, unknown values, and deprecation behavior.
- [ ] **USEQ-07F84FE4** — Validate check digits and structured identifiers without rejecting legitimate international formats unnecessarily.
- [ ] **USEQ-70E63BE2** — Prevent identifier reuse when historical ambiguity or security risk would result.
- [ ] **USEQ-7335C7DE** — Use pseudonymous identifiers for cross-system linkage only after privacy and reidentification risk review.
- [ ] **USEQ-0B4C02A3** — Monitor duplicate, orphaned, conflicting, and unmapped identifiers.

### Data modeling, integrity, and transactional behavior

- [ ] **USEQ-FE212FC1** — Model entities, values, relationships, events, temporal state, and lifecycle according to domain meaning.
- [ ] **USEQ-092595F5** — Distinguish facts, commands, events, snapshots, projections, and audit records.
- [ ] **USEQ-D8311858** — Enforce required uniqueness, referential integrity, valid ranges, and state constraints at reliable boundaries.
- [ ] **USEQ-F4B2A123** — Use transactions for operations that require atomicity and define boundaries explicitly.
- [ ] **USEQ-331C3046** — Document isolation level, concurrency behavior, lost-update prevention, and conflict resolution.
- [ ] **USEQ-B512CB86** — Define consistency requirements for replicas, caches, indexes, data lakes, warehouses, and derived views.
- [ ] **USEQ-7AF7CA68** — Prevent schema designs that allow contradictory representations of the same business state without reconciliation.
- [ ] **USEQ-E7FFA764** — Preserve history when users, regulation, audit, or temporal queries require it.
- [ ] **USEQ-22F7572F** — Distinguish event time, processing time, ingestion time, effective time, and recorded time.
- [ ] **USEQ-BAA2572C** — Define correction behavior without erasing required history or duplicating current state.
- [ ] **USEQ-413EC892** — Use constraints and checks close to the data rather than relying entirely on application convention.
- [ ] **USEQ-28CF339C** — Keep tenant, security, privacy, and regional boundaries visible in the model.
- [ ] **USEQ-2C8FC339** — Model deletion, redaction, retention, legal hold, and archival states explicitly where relevant.
- [ ] **USEQ-9E4EF632** — Review models with domain experts and actual consumers, not only storage specialists.

### Schema evolution, migrations, backfills, and repair

- [ ] **USEQ-008BD1A2** — Design schema and contract changes for the actual compatibility window and mixed-version operation.
- [ ] **USEQ-426954FC** — Prefer additive changes before destructive changes when consumers need time to migrate.
- [ ] **USEQ-69B5C63D** — Version schemas, transformations, and migration logic.
- [ ] **USEQ-CD128F39** — Test migrations against representative size, skew, nulls, duplicates, malformed legacy data, and boundary values.
- [ ] **USEQ-789C8367** — Measure lock, compute, I/O, replication, storage, and duration impact.
- [ ] **USEQ-D8E756A2** — Make migrations resumable, idempotent, observable, and safe after partial failure.
- [ ] **USEQ-62764F79** — Define preconditions, checkpoints, postconditions, validation, rollback feasibility, and roll-forward recovery.
- [ ] **USEQ-B8372734** — Create a recovery point before destructive or high-risk changes.
- [ ] **USEQ-D86FD28C** — Verify row counts, checksums, constraints, totals, balances, relationships, samples, and business invariants after migration.
- [ ] **USEQ-0B1BEE26** — Keep old and new representations synchronized only for a bounded, monitored transition.
- [ ] **USEQ-C9FBC4DE** — Do not delete old fields or values until actual consumers have migrated and recovery needs have expired.
- [ ] **USEQ-7902EC23** — Treat backfills as production changes with authorization, rate controls, monitoring, and incident readiness.
- [ ] **USEQ-833B7C73** — Distinguish repairs from falsification; preserve who changed what, why, when, and based on which evidence.
- [ ] **USEQ-61FAD800** — Rehearse the exact automation intended for production.

### Data sharing, exchange, and interoperability

- [ ] **USEQ-6421ECBB** — Define purpose, recipient, legal basis or authorization, fields, frequency, format, security, retention, and deletion for every material exchange.
- [ ] **USEQ-49553535** — Minimize exchanged data to what the recipient needs.
- [ ] **USEQ-B086FE62** — Use explicit schemas, media types, encodings, language tags, units, identifiers, and time semantics.
- [ ] **USEQ-2D1F4C4C** — Validate inbound data as untrusted and outbound data against the contract.
- [ ] **USEQ-E18571A1** — Authenticate and authorize producers and consumers.
- [ ] **USEQ-669A8451** — Protect confidentiality and integrity in transit and at rest according to classification.
- [ ] **USEQ-1D063B03** — Use idempotency, sequence, checksums, manifests, and acknowledgments where delivery integrity requires them.
- [ ] **USEQ-66C674E8** — Define duplicate, missing, out-of-order, late, partial, and rejected-file behavior.
- [ ] **USEQ-AD5B6E3C** — Reconcile both sides of critical exchanges.
- [ ] **USEQ-43E86B36** — Ensure export formats preserve meaning, accessibility, provenance, and relationships.
- [ ] **USEQ-476CDA79** — Prevent spreadsheet formula injection, active content, unsafe macros, and ambiguous formatting in exports.
- [ ] **USEQ-3FB36668** — Document format and API deprecation and retain parsers only for supported versions.
- [ ] **USEQ-3A11E515** — Ensure termination of a relationship revokes access and triggers required return or deletion.
- [ ] **USEQ-711BA21C** — Keep an exit and portability path for supplier-held data.

### Records management and evidentiary information

- [ ] **USEQ-71C30CC9** — Identify which transactions, decisions, approvals, communications, contracts, consents, notices, changes, and events must become records.
- [ ] **USEQ-8DFEBE9C** — Capture records at the point the business activity occurs or as soon as reliably possible.
- [ ] **USEQ-D145D470** — Preserve content, context, structure, identity, time, relationships, and necessary metadata.
- [ ] **USEQ-6820421D** — Ensure records are authentic, reliable, integral, usable, and discoverable for their required lifetime.
- [ ] **USEQ-AB5CF8A5** — Prevent unauthorized alteration, substitution, deletion, or backdating.
- [ ] **USEQ-4978DF17** — Preserve version and amendment history where a record can change.
- [ ] **USEQ-BF076338** — Distinguish records from transient working data, convenience copies, caches, and backups.
- [ ] **USEQ-EBF07F7B** — Apply records controls regardless of file format, storage system, communication channel, or supplier.
- [ ] **USEQ-F09291DC** — Record access, export, amendment, hold, disposition, and administrative actions when required.
- [ ] **USEQ-4B91F6F1** — Ensure electronic signatures and approvals retain identity, intent, scope, time, and verification evidence.
- [ ] **USEQ-D41732CE** — Preserve relationships between a record and attachments, referenced policies, source data, and approvals.
- [ ] **USEQ-C5ECB4A3** — Make records accessible to authorized users and investigators without exposing unrelated information.
- [ ] **USEQ-8A56F759** — Include chat, collaboration, ticketing, source-control, automated decisions, and machine events when they constitute business records.
- [ ] **USEQ-A243E554** — Review records requirements when workflows are automated or moved to new platforms.

### Retention, legal hold, disposition, and deletion

- [ ] **USEQ-99FDFA93** — Create a retention schedule based on legal, contractual, operational, historical, security, privacy, and user needs.
- [ ] **USEQ-89FC8D59** — Apply retention by data and record category rather than retaining everything indefinitely.
- [ ] **USEQ-CCB936EC** — Define the event that starts each retention period.
- [ ] **USEQ-34751429** — Automate retention and disposition where rules are deterministic and preserve review for exceptions.
- [ ] **USEQ-28E4DB6E** — Apply legal or investigation holds promptly and prevent covered information from normal deletion.
- [ ] **USEQ-D23C610B** — Record hold scope, authority, owner, start, review, release, and affected systems.
- [ ] **USEQ-013AF61E** — Ensure holds propagate to relevant primary stores, archives, collaboration systems, exports, and suppliers.
- [ ] **USEQ-562793D2** — Release holds formally and resume the correct retention schedule.
- [ ] **USEQ-D85B40F1** — Delete or anonymize data when retention expires and no overriding obligation remains.
- [ ] **USEQ-BAD0B1E6** — Apply deletion to replicas, caches, indexes, derived stores, test copies, and downstream processors as required.
- [ ] **USEQ-4EBFFEDF** — Document how immutable backups and archives handle expired or deleted data.
- [ ] **USEQ-9E03D082** — Use verifiable deletion or cryptographic erasure where the threat model and technology support it.
- [ ] **USEQ-234901AA** — Prevent disposal from destroying required audit, consent, transaction, or provenance evidence.
- [ ] **USEQ-1CF079B2** — Retain disposition logs without retaining the disposed sensitive content unnecessarily.

### Digital preservation and long-term usability

- [ ] **USEQ-EE767AA5** — Distinguish backup for recovery from archival preservation for long-term authenticity and usability.
- [ ] **USEQ-B12AA085** — Define the designated future users and the knowledge they will need to understand preserved information.
- [ ] **USEQ-1C2CC92A** — Select preservation scope based on legal, historical, scientific, contractual, product, and evidentiary value.
- [ ] **USEQ-61987E09** — Preserve content, metadata, provenance, rights, context, representation information, and relationships.
- [ ] **USEQ-854D1510** — Use durable, documented, non-obsolete formats or maintain a migration and emulation strategy.
- [ ] **USEQ-3307E3AE** — Keep format, software, dependency, encoding, compression, and rendering information needed for future access.
- [ ] **USEQ-4488AC37** — Use fixity checks and detect bit-level corruption.
- [ ] **USEQ-1DBDA1D2** — Replicate preserved information across appropriately independent failure domains.
- [ ] **USEQ-0E7EE682** — Protect archives from unauthorized alteration, deletion, ransomware, and privilege abuse.
- [ ] **USEQ-936A5306** — Record preservation events, agents, rights, objects, transformations, validations, and failures using a model such as PREMIS where appropriate.
- [ ] **USEQ-BF787168** — Validate readability and interpretability periodically, not only storage availability.
- [ ] **USEQ-B272D712** — Test format migration and verify semantic, visual, structural, accessibility, and evidentiary fidelity.
- [ ] **USEQ-952007B8** — Preserve original objects when policy requires them, even after normalized or migrated versions are created.
- [ ] **USEQ-5701C9BF** — Document unavoidable loss and obtain authorized acceptance.
- [ ] **USEQ-C69460B4** — Plan transfer, succession, and exit if the preserving organization or supplier ceases operation.

### Privacy, security, access, and ethical data use

- [ ] **USEQ-51E193D6** — Apply least privilege and purpose limitation to data access.
- [ ] **USEQ-AC0E1928** — Use individual or workload identities and prohibit shared access where accountability matters.
- [ ] **USEQ-85A89DBD** — Separate administrative, analytical, support, model-training, and production access.
- [ ] **USEQ-7C5BDF9F** — Mask, tokenize, pseudonymize, aggregate, or synthesize data where full fidelity is unnecessary.
- [ ] **USEQ-5A5A6DAF** — Assess reidentification and linkage risk before releasing deidentified data.
- [ ] **USEQ-B178EA24** — Encrypt data and backups according to classification and key-management requirements.
- [ ] **USEQ-14B33384** — Prevent secrets and personal data from entering logs, traces, analytics, test fixtures, source code, and public catalogs unnecessarily.
- [ ] **USEQ-626A7CF3** — Audit sensitive reads, exports, changes, and bulk operations.
- [ ] **USEQ-F3236573** — Use query, export, and rate controls to prevent enumeration and mass extraction.
- [ ] **USEQ-93B7958E** — Review new purposes and secondary uses before processing existing data.
- [ ] **USEQ-08199A24** — Ensure users can exercise applicable access, correction, deletion, restriction, objection, and portability rights across downstream systems.
- [ ] **USEQ-D564BA95** — Assess discrimination, exclusion, manipulation, and safety risk from data collection and use.
- [ ] **USEQ-05BD5FB0** — Do not use apparent data availability as evidence of permission or ethical acceptability.
- [ ] **USEQ-EDE596A0** — Notify affected owners and consumers when data confidentiality, integrity, availability, or provenance is compromised.

### Analytics, reporting, experimentation, and derived decisions

- [ ] **USEQ-823B644E** — Define the decision or user outcome each metric, report, dashboard, or experiment supports.
- [ ] **USEQ-D14A4E95** — Document metric formula, population, exclusions, time window, source, refresh, owner, and known limitations.
- [ ] **USEQ-821F40F1** — Keep denominator and missing-data treatment explicit.
- [ ] **USEQ-71B2CB94** — Prevent changes in definition from appearing as real-world changes without annotation and backfill policy.
- [ ] **USEQ-6BD8DCC9** — Version reports, semantic models, transformation logic, and experiment assignments.
- [ ] **USEQ-00A52128** — Reconcile critical reports with source systems and independent totals.
- [ ] **USEQ-C0F9D6FB** — Prevent dashboards from presenting stale, partial, estimated, or low-quality data without clear status.
- [ ] **USEQ-AE91E43D** — Use statistical methods appropriate to sampling, multiple comparisons, uncertainty, and causal claims.
- [ ] **USEQ-3352E54A** — Do not infer individual performance or protected characteristics from weak proxies.
- [ ] **USEQ-E2D57064** — Protect experimentation from sample-ratio mismatch, interference, assignment leakage, and premature stopping.
- [ ] **USEQ-4E23A195** — Preserve experiment hypothesis, design, assignment, exposure, analysis, results, and decision history.
- [ ] **USEQ-24B4D0EC** — Ensure downstream decisions can be traced to the data and model versions used.
- [ ] **USEQ-BC0C13E1** — Provide correction and appeal routes for high-impact data-driven decisions where appropriate.
- [ ] **USEQ-66275973** — Review whether metrics create harmful incentives or encourage manipulation.

### Machine learning and model data

- [ ] **USEQ-29D69EEF** — Inventory training, fine-tuning, evaluation, retrieval, prompt, feedback, and monitoring data separately.
- [ ] **USEQ-8128F64B** — Record source, rights, consent, collection method, transformations, labeling, quality, representativeness, and prohibited uses.
- [ ] **USEQ-EB53C881** — Prevent train-test contamination and evaluation leakage.
- [ ] **USEQ-DBCA7FD5** — Version datasets, splits, labels, feature logic, prompts, retrieval indexes, and model inputs.
- [ ] **USEQ-B17BB267** — Assess coverage, imbalance, historical bias, label uncertainty, adversarial contamination, and distribution shift.
- [ ] **USEQ-1F4A0CFC** — Remove or protect secrets, personal data, copyrighted material, and regulated information according to approved use.
- [ ] **USEQ-00D122F0** — Define deletion and correction propagation into derived datasets, indexes, embeddings, and retraining workflows.
- [ ] **USEQ-874DCDFD** — Use provenance sufficient to reproduce or explain model behavior and evaluations.
- [ ] **USEQ-2B72FFCA** — Keep human labeling instructions, qualifications, disagreement, quality review, and working conditions documented.
- [ ] **USEQ-2C7DBAB6** — Validate synthetic data for fidelity, privacy leakage, bias, and unintended memorization.
- [ ] **USEQ-3EDB888F** — Monitor production input and output drift against approved operating bounds.
- [ ] **USEQ-4CD1651B** — Do not treat model-generated labels or data as ground truth without validation.
- [ ] **USEQ-F0CEE610** — Preserve model cards, data cards, evaluation evidence, and deployment context for high-impact systems.
- [ ] **USEQ-D3C328B9** — Apply access and supply-chain controls to model weights, datasets, and embeddings.

### Backup, restore, archive, and recovery distinctions

- [ ] **USEQ-40113369** — Define recovery point and recovery time objectives for operational data.
- [ ] **USEQ-D486C282** — Define preservation and retention objectives separately from disaster recovery.
- [ ] **USEQ-E0F9E3DA** — Include schemas, configuration, keys or recovery mechanisms, metadata, lineage, permissions, and external dependencies in recovery planning.
- [ ] **USEQ-553211DF** — Encrypt and isolate backups from ordinary production compromise.
- [ ] **USEQ-496019E2** — Monitor backup completion, size, integrity, retention, and deletion.
- [ ] **USEQ-8270225E** — Test full, partial, point-in-time, clean-environment, and regional restoration as applicable.
- [ ] **USEQ-F328E54D** — Validate restored constraints, totals, lineage, permissions, deletions, revoked identities, and current secrets.
- [ ] **USEQ-0F570870** — Prevent restore from replaying charges, notifications, webhooks, jobs, or obsolete permissions unintentionally.
- [ ] **USEQ-41E9E282** — Reconcile restored data with external systems and events that occurred after the recovery point.
- [ ] **USEQ-27CFB60F** — Document which data can be reconstructed and which must be backed up.
- [ ] **USEQ-E43E0612** — Ensure archives are not used as an undocumented production dependency.
- [ ] **USEQ-9849FEA2** — Ensure backups are not treated as permanent archives without preservation metadata and usability validation.
- [ ] **USEQ-71B8E83F** — Protect restore tooling and credentials as critical privileged capabilities.
- [ ] **USEQ-3549FA02** — Retain restore evidence and remediate failed exercises.

### Data incident response, release gates, and assurance

- [ ] **USEQ-193E7A83** — Define data-incident severities for confidentiality loss, corruption, deletion, staleness, duplication, lineage loss, and materially wrong decisions.
- [ ] **USEQ-4F9D5FEF** — Provide playbooks for quarantine, publication stop, source isolation, replay, backfill, correction, rollback, notification, and consumer impact analysis.
- [ ] **USEQ-691A9DDE** — Identify every downstream asset, user, decision, report, model, and external recipient affected by a data incident.
- [ ] **USEQ-AAEC783F** — Preserve evidence and maintain a timeline of source, transformation, publication, detection, and correction.
- [ ] **USEQ-BD64650D** — Communicate uncertainty and affected time ranges honestly.
- [ ] **USEQ-D698D5A6** — Block release when critical data contracts are undefined or untested.
- [ ] **USEQ-741CEC56** — Block release when migration, backfill, or repair lacks verified recovery and post-change reconciliation.
- [ ] **USEQ-87BA665A** — Block release when a critical output can be stale, empty, duplicated, cross-tenant, or materially wrong without detection.
- [ ] **USEQ-5ACB080F** — Block release when required records, provenance, retention, legal hold, or deletion behavior is missing.
- [ ] **USEQ-72023963** — Block release when data access, export, or sharing exceeds approved purpose or authorization.
- [ ] **USEQ-F8A6115D** — Attach schemas, contracts, lineage, quality results, migration evidence, reconciliation, access review, retention rules, and recovery tests to release evidence.
- [ ] **USEQ-37F5668B** — Require data-owner sign-off for changes to authoritative definitions, high-impact reports, sensitive sharing, or irreversible transformation.
- [ ] **USEQ-9EA79399** — Convert incidents into contract, validation, lineage, monitoring, and governance improvements.
- [ ] **USEQ-1147F53E** — Reassess data controls after material source, model, consumer, or regulatory change.

## Final Gap Closure — Data Accountability, eDiscovery, Legal Holds, and Evidentiary Integrity

_Consolidated from `final consolidated corpus/04-data-privacy-information-governance-records-analytics.md#Final Gap Closure — Data Accountability, eDiscovery, Legal Holds, and Evidentiary Integrity`; 122 non-duplicative controls._

### Governing-body data accountability

- [ ] **USEQ-3A8E697F** — Assign governing-body or executive accountability for material data, records, analytics, and information risks.
- [ ] **USEQ-7DD055FC** — Define delegated authority for data ownership, stewardship, custodianship, privacy, security, records, and analytics decisions.
- [ ] **USEQ-D6AAEE4D** — Maintain an authoritative inventory of material data assets, systems of record, records classes, custodians, processors, and repositories.
- [ ] **USEQ-54F0B304** — Classify data according to confidentiality, integrity, availability, privacy, evidentiary, safety, financial, and societal impact.
- [ ] **USEQ-DF3A9CDD** — Define decision rights for collection, use, sharing, transformation, retention, archival, legal hold, and deletion.
- [ ] **USEQ-A0AEB0D7** — Require data owners to approve material changes to purpose, semantics, quality requirements, access, or retention.
- [ ] **USEQ-D2DA1696** — Escalate unresolved data-quality, integrity, lineage, or rights conflicts to an authority independent of delivery pressure.
- [ ] **USEQ-30D278A2** — Include data risk and value in portfolio, architecture, acquisition, and retirement decisions.
- [ ] **USEQ-68B33A8A** — Measure data outcomes such as fitness, trustworthiness, timeliness, correction, rights fulfilment, and incident recurrence.
- [ ] **USEQ-70FFD8B8** — Review whether incentives encourage excess collection, concealed defects, misleading metrics, or indefinite retention.
- [ ] **USEQ-0A577D72** — Ensure data accountability continues across suppliers, subprocessors, joint ventures, and shared platforms.
- [ ] **USEQ-59260332** — Preserve records of material data decisions, assumptions, evidence, dissent, and approvals.
- [ ] **USEQ-86FFEAC1** — Review data governance after incidents, litigation, regulatory action, mergers, major migrations, and new data uses.

### Discovery-readiness inventory and mapping

- [ ] **USEQ-00EEEE2B** — Identify systems that can contain records or information relevant to disputes, investigations, audits, or lawful requests.
- [ ] **USEQ-EC109500** — Include collaboration tools, email, messaging, tickets, source control, logs, cloud storage, endpoints, mobile devices, backups, archives, analytics, AI systems, and supplier systems as applicable.
- [ ] **USEQ-FC0EAAB7** — Record repository owners, administrators, locations, formats, retention behavior, search capability, export capability, and access controls.
- [ ] **USEQ-A831DD6E** — Identify ephemeral, encrypted, compressed, proprietary, deleted, versioned, or dynamically generated information.
- [ ] **USEQ-ADC96CBD** — Identify where data is replicated, cached, indexed, transformed, summarized, or embedded.
- [ ] **USEQ-2015A1C9** — Map custodians and organizational roles to systems and data sources.
- [ ] **USEQ-CDF23B7C** — Record time-zone, clock, identity, naming, and identifier conventions needed to correlate evidence.
- [ ] **USEQ-69EDB129** — Distinguish authoritative records from convenience copies and derived views.
- [ ] **USEQ-2C6A80AE** — Document supplier assistance, contractual rights, fees, deadlines, and technical constraints for collection.
- [ ] **USEQ-B55ED826** — Test whether required information can be found and exported before an urgent matter occurs.
- [ ] **USEQ-A82DB616** — Keep the discovery-readiness inventory current through system onboarding, change, and retirement.

### Legal holds and preservation duties

- [ ] **USEQ-EF106894** — Define who can issue, approve, modify, release, and audit a legal or investigative hold.
- [ ] **USEQ-CDBF3942** — Trigger preservation promptly when litigation, investigation, audit, or another duty is reasonably anticipated or formally received.
- [ ] **USEQ-501834EB** — Define the matter scope, custodians, data sources, subjects, date ranges, preservation method, and responsible owners.
- [ ] **USEQ-B6755AD2** — Suspend conflicting deletion, overwrite, rotation, compaction, archival expiry, and account-closure processes.
- [ ] **USEQ-BA845436** — Apply holds to downstream copies, exports, derived data, archives, and supplier-held information as required.
- [ ] **USEQ-0A31DC29** — Prevent held information from being altered, concealed, or destroyed without authorization.
- [ ] **USEQ-380A0A17** — Notify affected custodians clearly and require acknowledgment where appropriate.
- [ ] **USEQ-6EE4DE8B** — Preserve confidentiality, privilege, privacy, and need-to-know restrictions during hold administration.
- [ ] **USEQ-1E1EFD59** — Track hold coverage, acknowledgment, collection, changes, exceptions, and release.
- [ ] **USEQ-1CFB49B8** — Reassess scope as facts, claims, custodians, systems, or legal duties change.
- [ ] **USEQ-666E4FF9** — Validate that technical holds actually stop disposition in each relevant repository.
- [ ] **USEQ-6A030672** — Monitor repository migrations and account offboarding so holds remain effective.
- [ ] **USEQ-D4BAC821** — Release holds through an authorized documented process when preservation duty ends.
- [ ] **USEQ-1EFE79A1** — Resume ordinary retention and defensible disposition after release without deleting information still subject to another hold.
- [ ] **USEQ-7DE623E7** — Preserve evidence of hold issuance, operation, validation, and release.

### Defensible preservation and collection

- [ ] **USEQ-AD5B82C4** — Select preservation and collection methods proportionate to the matter, source, volatility, burden, and evidentiary need.
- [ ] **USEQ-E3238C67** — Preserve original content, metadata, relationships, versions, timestamps, permissions, and context needed for interpretation.
- [ ] **USEQ-8C9A995F** — Use repeatable, documented, and validated collection procedures.
- [ ] **USEQ-14580CAC** — Record who collected what, when, where, how, under whose authority, and using which tools and versions.
- [ ] **USEQ-65B40E66** — Calculate and retain integrity hashes or equivalent fixity evidence where appropriate.
- [ ] **USEQ-E6E5A19C** — Preserve source time-zone and clock information.
- [ ] **USEQ-53CEDBC1** — Avoid changing access times, metadata, content, or application state unnecessarily during collection.
- [ ] **USEQ-2F2046BA** — Document unavoidable transformations, normalization, filtering, or loss.
- [ ] **USEQ-DF0A4C8B** — Capture encrypted material with lawful key-access and recovery planning.
- [ ] **USEQ-D00DDFDC** — Collect dynamic, collaborative, threaded, or linked content in a form that preserves material context.
- [ ] **USEQ-5DD374B2** — Preserve references between messages, attachments, comments, revisions, and source records.
- [ ] **USEQ-007BF418** — Handle deleted, hidden, archived, or inaccessible data according to documented legal and proportionality decisions.
- [ ] **USEQ-784EB923** — Separate privileged, restricted, personal, regulated, and trade-secret information with appropriate controls.
- [ ] **USEQ-7006D883** — Validate collection completeness through counts, samples, queries, reconciliation, or independent review.
- [ ] **USEQ-F5E05710** — Retain collection logs, errors, exclusions, retry history, and unresolved gaps.
- [ ] **USEQ-937ABB34** — Protect collected material from unauthorized access, alteration, malware, and accidental disclosure.

### Chain of custody and evidentiary integrity

- [ ] **USEQ-2AA81182** — Assign a unique identifier to each collected evidence set and material item.
- [ ] **USEQ-052AE940** — Maintain an unbroken record of possession, transfer, access, copying, transformation, and disposition.
- [ ] **USEQ-143D9842** — Restrict evidence access to authorized individuals with a documented purpose.
- [ ] **USEQ-8FE404EE** — Use tamper-evident storage and integrity verification proportionate to evidentiary risk.
- [ ] **USEQ-39C82E33** — Verify integrity before and after transfer, processing, restoration, and production.
- [ ] **USEQ-2C35861F** — Preserve original evidence separately from working copies.
- [ ] **USEQ-A993560C** — Record tool versions, configuration, time, environment, and operator for material processing.
- [ ] **USEQ-F6B54DAC** — Validate tools used for acquisition, parsing, conversion, deduplication, search, and export.
- [ ] **USEQ-B63156A8** — Document known tool limitations and parsing failures.
- [ ] **USEQ-1292500D** — Protect audit and custody records from unauthorized modification or deletion.
- [ ] **USEQ-1111C3BC** — Use secure, authenticated transfer channels and recipient verification.
- [ ] **USEQ-CC7E87FB** — Maintain controlled evidence storage during provider or platform migration.
- [ ] **USEQ-B4D527BC** — Ensure disaster recovery preserves evidence integrity and custody history.
- [ ] **USEQ-AA4171CE** — Define authorized destruction and retain proof of disposition when evidence is no longer required.

### Processing, search, review, and production

- [ ] **USEQ-18028100** — Define processing rules, search methods, filters, deduplication, threading, date handling, and exclusions before use.
- [ ] **USEQ-355B9882** — Validate that search syntax and indexing behavior match the intended query semantics.
- [ ] **USEQ-9DD20C62** — Test representative known documents to verify recall and precision.
- [ ] **USEQ-F91F42B5** — Document search terms, versions, iterations, reviewers, and decision rationale.
- [ ] **USEQ-DCAE9F4C** — Preserve the ability to reproduce the processed population and result set.
- [ ] **USEQ-25FE2994** — Keep original and normalized representations linked.
- [ ] **USEQ-00F56848** — Avoid deduplication that removes materially different metadata, recipients, permissions, or context.
- [ ] **USEQ-395FF585** — Use defensible sampling and quality control for large review populations.
- [ ] **USEQ-93A61AF7** — Train reviewers on scope, relevance, privilege, confidentiality, and escalation.
- [ ] **USEQ-042D0B84** — Measure and address inconsistent review decisions.
- [ ] **USEQ-E0A87576** — Apply redaction without altering unredacted source evidence.
- [ ] **USEQ-EF45812D** — Validate redactions cannot be reversed through layers, metadata, hidden text, or alternate representations.
- [ ] **USEQ-C589DC56** — Verify productions include agreed fields, metadata, numbering, formats, load files, and relationships.
- [ ] **USEQ-548B21E0** — Confirm exports do not include information outside authorized scope.
- [ ] **USEQ-75CD5DEE** — Encrypt and authenticate productions in transit and at rest.
- [ ] **USEQ-094254BF** — Record what was produced, to whom, when, under what authority, and with which integrity evidence.
- [ ] **USEQ-338A4A7F** — Retain production validation and exception records.

### Privacy, proportionality, privilege, and cross-border constraints

- [ ] **USEQ-8A4B8912** — Limit preservation, collection, review, and production to information reasonably necessary for the authorized purpose.
- [ ] **USEQ-3CD2B624** — Consider less intrusive methods and staged discovery before broad collection.
- [ ] **USEQ-75FA260D** — Apply jurisdiction, employment, communications, secrecy, blocking, localization, and transfer requirements.
- [ ] **USEQ-1493E5C4** — Obtain qualified legal review for cross-border collection and disclosure.
- [ ] **USEQ-F7581CD7** — Use access restrictions, pseudonymization, filtering, confidentiality terms, or protective procedures where appropriate.
- [ ] **USEQ-CC4CBEA4** — Identify and protect legally privileged, confidential, health, child, biometric, financial, authentication, and other sensitive information.
- [ ] **USEQ-1A33A5A8** — Prevent discovery tooling from becoming an unrestricted surveillance capability.
- [ ] **USEQ-4040EE6C** — Keep investigative and litigation data separate from ordinary analytics and product use.
- [ ] **USEQ-490FD16A** — Prohibit secondary use of preserved or collected material unless separately authorized.
- [ ] **USEQ-2DCEA90A** — Delete or return matter data when obligations and authorized purposes end.
- [ ] **USEQ-2EFD0053** — Include suppliers and external counsel in privacy, security, breach, and deletion obligations.
- [ ] **USEQ-69B301CD** — Record proportionality, burden, risk, and alternative-method decisions.

### Records, backups, archives, and legal-hold distinctions

- [ ] **USEQ-FD46F2D9** — Define which information constitutes an authoritative record and which is a backup, archive, replica, cache, or transient copy.
- [ ] **USEQ-D3571578** — Do not use backup retention as a substitute for records management.
- [ ] **USEQ-3191EF19** — Do not assume ordinary backup systems provide searchable, complete, or legally defensible preservation.
- [ ] **USEQ-D1BB5019** — Document whether and how holds apply to immutable backups and rolling backup media.
- [ ] **USEQ-50CFA71B** — Define when restoration from backup is required, proportionate, technically feasible, and authorized.
- [ ] **USEQ-2F26A1EC** — Prevent restored backups from silently reintroducing expired credentials, deleted data, or obsolete system state.
- [ ] **USEQ-EF2DB422** — Preserve archive representation information needed to interpret records over time.
- [ ] **USEQ-CEAE5A5D** — Validate fixity, readability, and retrievability of held archival records periodically.
- [ ] **USEQ-0552820A** — Ensure decommissioning plans account for active holds, records obligations, and future discovery needs.
- [ ] **USEQ-FB12DE3D** — Maintain tools or migration paths for proprietary formats retained beyond product support.

### Investigation and audit readiness

- [ ] **USEQ-6456A239** — Define procedures for internal investigations, regulator requests, external audits, and lawfully authorized disclosures.
- [ ] **USEQ-014CFEB0** — Separate fact finding, legal advice, technical response, disciplinary decisions, and public communication responsibilities.
- [ ] **USEQ-2B91F303** — Preserve investigation independence and protect against conflicts of interest.
- [ ] **USEQ-CE36DD64** — Prevent subjects of an investigation from controlling relevant evidence or access decisions.
- [ ] **USEQ-FF473CF7** — Establish secure case-management, communication, and evidence repositories.
- [ ] **USEQ-DC93FE5C** — Record allegations, scope, authority, evidence, findings, limitations, decisions, and corrective actions.
- [ ] **USEQ-E888F262** — Distinguish verified facts, analysis, inference, allegation, and unresolved uncertainty.
- [ ] **USEQ-B26054AD** — Provide fair review and response opportunities where required.
- [ ] **USEQ-02B57300** — Preserve anonymity or confidentiality of reporters to the extent possible.
- [ ] **USEQ-6A813913** — Protect against retaliation and monitor for it.
- [ ] **USEQ-D1208FCF** — Coordinate incident containment with evidence preservation so neither objective is ignored.
- [ ] **USEQ-7A78E7C0** — Test discovery and investigation readiness through exercises or representative retrieval drills.
- [ ] **USEQ-3F523D13** — Track time to locate, preserve, collect, validate, and produce required information.
- [ ] **USEQ-436EB79A** — Treat inability to preserve or retrieve material records as a production and governance risk.

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
