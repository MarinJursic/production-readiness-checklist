# AI, ML, and AI-assisted development

_Phase 15 of 16 in the [complete engineering review](00-overview.md)._

AI governance, data and model supply chains, evaluation, MLOps, agents, human oversight, incidents, and AI-assisted engineering.

Assess every control as **Pass**, **Fail**, **Blocked**, or **Not Applicable**. A Pass needs current evidence for the exact product, revision, artifact, configuration, environment, and data state under review.

## AI, Machine Learning, Autonomous, and LLM Systems

_Consolidated from `quality standards/15-conditional-domains/05-ai-machine-learning-and-llm-systems.md`; 25 non-duplicative controls._

### Universal controls

- [ ] **USEQ-2C5A990D** — Inventory every model, model version, provider, prompt, instruction template, tool, data source, retrieval system, embedding store, safety component, and evaluation set.
- [ ] **USEQ-08ACDAA3** — Classify the consequences of incorrect, biased, unsafe, manipulated, unavailable, or overly confident output.
- [ ] **USEQ-32B13D60** — Model instructions cannot override authorization, privacy, safety, compliance, or business rules.
- [ ] **USEQ-197BC632** — The model cannot choose credentials, tenants, resources, recipients, permissions, or financial destinations beyond the user's authorization.
- [ ] **USEQ-BC3602E0** — Model output is validated before use in code, queries, commands, templates, transactions, identity decisions, or security controls.
- [ ] **USEQ-96538752** — Sensitive data is excluded from prompts, retrieval, fine-tuning, evaluation, and training unless specifically approved.
- [ ] **USEQ-0D090FAC** — Provider retention, logging, human-review, model-training, and data-location terms are understood.
- [ ] **USEQ-A518D322** — Vector stores and embeddings preserve tenant boundaries, retention, deletion, and access changes.
- [ ] **USEQ-DC5087FC** — Retrieval poisoning, malicious documents, data contamination, model poisoning, and supply-chain compromise are addressed.
- [ ] **USEQ-A85B7AA6** — Training, fine-tuning, evaluation, prompt, and model data have provenance and governance.
- [ ] **USEQ-2E49E2AA** — Model, prompt, tool, retrieval, and safety changes are versioned, reviewed, tested, and approved.
- [ ] **USEQ-296AC402** — Evaluation sets cover accuracy, reliability, security, safety, privacy, bias, fairness, refusal, adversarial behavior, and misuse as applicable.
- [ ] **USEQ-C203537F** — Evaluations represent production languages, user groups, data distributions, and edge cases.
- [ ] **USEQ-0B71A9C7** — Hallucination, uncertainty, sources, and limitations are communicated appropriately.
- [ ] **USEQ-26AC7838** — High-impact decisions have suitable human review, explanation, appeal, and correction paths.
- [ ] **USEQ-E48B3005** — Output moderation and abuse detection match product risk.
- [ ] **USEQ-009D08C1** — Model denial of service, context exhaustion, tool loops, recursive calls, and cost amplification are bounded.
- [ ] **USEQ-D7B17CB8** — Rate, token, context, tool, concurrency, and spending limits are configured.
- [ ] **USEQ-C9282DB5** — Model and provider outages have tested fallback behavior.
- [ ] **USEQ-EF076436** — Model drift, quality regression, safety regression, and provider behavior changes are monitored.
- [ ] **USEQ-6D4A4F4E** — AI telemetry protects user privacy and does not unnecessarily record sensitive prompts or outputs.
- [ ] **USEQ-701EF00D** — Logs record model, prompt, tool, retrieval, and policy versions sufficiently for investigation.
- [ ] **USEQ-1F35FBC7** — Users know when they are interacting with or materially affected by AI where required.
- [ ] **USEQ-0C2884C7** — Applicable transparency, documentation, risk-assessment, copyright, data, human-oversight, and automated-decision rules are mapped.
- [ ] **USEQ-E5718654** — Retirement removes obsolete models, credentials, prompts, indexes, retained data, and tools safely.

## AI, Machine Learning, MLOps, and AI-Assisted Development Master Checklist

_Consolidated from `gap supplement/05-ai-ml-mlops-and-ai-assisted-development.md`; 318 non-duplicative controls._

### Expanded gap-closure controls

#### AI governance, inventory, and accountability

- [ ] **USEQ-C529B9D3** — Apply this checklist whenever a learned, probabilistic, generative, ranking, recommendation, classification, optimization, forecasting, biometric, or autonomous component materially influences product behavior or engineering work.
- [ ] **USEQ-814E1CE9** — Maintain a current inventory of models, model versions, prompts, system instructions, agents, tools, retrieval sources, embeddings, feature stores, training pipelines, evaluation sets, safety systems, providers, deployment endpoints, and downstream consumers.
- [ ] **USEQ-FD2B24B1** — Assign accountable business, product, model, data, security, privacy, safety, legal, compliance, accessibility, operations, and human-oversight owners according to impact.
- [ ] **USEQ-74CF2922** — Classify each system by affected users, decision consequence, reversibility, autonomy, scale, sensitivity, legal status, abuse potential, and dependence on probabilistic output.
- [ ] **USEQ-FAD9068A** — Document the intended purpose, excluded purposes, users, affected parties, operating conditions, expected benefit, assumptions, limitations, prohibited actions, and retirement criteria.
- [ ] **USEQ-7E87BDD7** — Define an AI risk-management lifecycle covering design, acquisition, data, development, evaluation, deployment, monitoring, incident response, change, and retirement.
- [ ] **USEQ-F7AD4215** — Require proportionately independent review for high-impact, safety-related, rights-affecting, financial, employment, education, healthcare, identity, law-enforcement, or public-sector uses.
- [ ] **USEQ-823F07DE** — Map applicable laws, sector rules, contracts, organizational policies, standards, licenses, intellectual-property rights, data-use commitments, and provider terms before development or procurement.
- [ ] **USEQ-D63FB3F4** — Maintain an AI system impact assessment that considers direct, indirect, cumulative, group, societal, environmental, accessibility, labor, and downstream effects.
- [ ] **USEQ-6A97D822** — Identify people who may be affected without being direct users and provide appropriate protection, notice, recourse, and monitoring.
- [ ] **USEQ-D8BEBA66** — Define risk appetite, unacceptable outcomes, stop conditions, escalation thresholds, approval authorities, and residual-risk owners before production use.
- [ ] **USEQ-96B9FB18** — Separate responsibility for building, evaluating, approving, operating, auditing, and accepting material AI risk where independence is required.
- [ ] **USEQ-B8096925** — Keep a traceable decision record for model selection, data selection, evaluation thresholds, safety controls, launch, exceptions, and material changes.
- [ ] **USEQ-DABACCB2** — Review the AI risk assessment whenever the purpose, population, geography, data, provider, model, prompt, tool, autonomy, or downstream use changes.
- [ ] **USEQ-8867EC56** — Prevent an experimental label from being used to bypass security, privacy, accessibility, safety, legal, documentation, or reliability requirements when real users or production data are involved.
- [ ] **USEQ-B399DE71** — Provide governance for shadow AI use, unregistered models, personal accounts, browser extensions, local models, and unsanctioned data uploads.
- [ ] **USEQ-AC261112** — Define who can pause, disable, roll back, isolate, or permanently retire an AI capability during an incident.

#### Use-case boundaries and consequential decisions

- [ ] **USEQ-041598F3** — Use AI only where its demonstrated benefit exceeds the risks and a simpler deterministic method is not clearly more suitable.
- [ ] **USEQ-0D09CDA5** — Do not delegate a decision merely because data and a model are available; establish a legitimate purpose and decision owner.
- [ ] **USEQ-FB01B22D** — Define which outputs are advisory, draft, ranked, automated, determinative, or capable of triggering tools and side effects.
- [ ] **USEQ-7D39EBC4** — Enforce purpose and authority outside the model so prompts or outputs cannot expand the system's permitted role.
- [ ] **USEQ-20555083** — Identify decisions that require a human, dual control, professional judgment, due process, or a deterministic rule.
- [ ] **USEQ-2FC26B87** — Keep irreversible, destructive, safety-critical, rights-affecting, financial, legal, and externally binding actions behind independently validated controls and appropriate confirmation.
- [ ] **USEQ-8521E67B** — Define confidence, uncertainty, abstention, escalation, and fallback behavior for each material task.
- [ ] **USEQ-48E1F506** — Prefer refusal, safe degradation, or handoff when inputs are outside the validated operating domain.
- [ ] **USEQ-99E4E2DD** — Ensure the system does not present generated content, predictions, classifications, or recommendations as established facts beyond available evidence.
- [ ] **USEQ-4549C1EB** — Prevent conversational fluency, personalization, anthropomorphism, emotional language, or interface design from overstating capability, authority, consciousness, confidentiality, or professional qualification.
- [ ] **USEQ-4820BFF1** — Disclose material AI involvement when users need it to make an informed choice or when law, policy, contract, or safety requires it.
- [ ] **USEQ-A79E9D12** — Provide a non-AI or assisted path where exclusion, accessibility, due process, or essential-service access requires one.
- [ ] **USEQ-34C36FCB** — Make material decisions contestable and correctable, with an understandable route to human review where appropriate.
- [ ] **USEQ-22AA7070** — Do not infer or use highly sensitive traits unless necessity, legality, validity, proportionality, and protection have been established.
- [ ] **USEQ-26F1B463** — Do not repurpose outputs for a materially different decision without a new impact assessment and validation.
- [ ] **USEQ-3FA42F8C** — Prevent secondary users from treating an output as valid for populations, contexts, languages, or decisions outside the documented scope.
- [ ] **USEQ-9C03D06C** — Define the acceptable frequency and impact of false positives, false negatives, ranking errors, hallucinations, refusals, and harmful recommendations by use case.
- [ ] **USEQ-89A25C8A** — Require safer defaults for children, vulnerable people, coercive contexts, and users who may mistake the system for a human or authority.

#### Data rights, provenance, preparation, and governance

- [ ] **USEQ-C65B1084** — Document the origin, ownership, license, consent, collection context, allowed purposes, geographic constraints, sensitivity, quality, representation, and retention of every material data source.
- [ ] **USEQ-4F1A4158** — Verify that training, fine-tuning, retrieval, evaluation, logging, feedback, and monitoring uses are authorized independently rather than assuming one permission covers all uses.
- [ ] **USEQ-480594DD** — Exclude data that was unlawfully obtained, contractually prohibited, inappropriately sensitive, materially inaccurate, or incapable of required deletion and governance.
- [ ] **USEQ-C72C1169** — Minimize personal, confidential, copyrighted, regulated, and security-sensitive data throughout the AI lifecycle.
- [ ] **USEQ-9E30F902** — Use privacy-preserving, synthetic, aggregated, pseudonymized, or deidentified data where it can meet the validated objective with less risk.
- [ ] **USEQ-23095819** — Assess whether deidentification remains effective after linkage, embedding, memorization, model access, retrieval, and output generation.
- [ ] **USEQ-7C20FBBF** — Preserve data lineage from raw source through filtering, labeling, transformation, splitting, augmentation, training, evaluation, and deployment.
- [ ] **USEQ-210D4834** — Record dataset, label, feature, prompt, retrieval-corpus, and preprocessing versions in every reproducible experiment and release.
- [ ] **USEQ-A53F8A85** — Define data-quality requirements for relevance, representativeness, completeness, accuracy, timeliness, consistency, duplication, contamination, and labeling reliability.
- [ ] **USEQ-839FD21F** — Measure coverage across intended languages, regions, demographics, disabilities, devices, domains, edge conditions, and rare but consequential cases.
- [ ] **USEQ-3A16FE4E** — Investigate sampling, measurement, historical, survivorship, selection, label, annotation, and feedback-loop bias.
- [ ] **USEQ-F5094424** — Use qualified annotators, documented instructions, disagreement handling, quality sampling, worker protection, and fair labor practices for human labeling.
- [ ] **USEQ-4CEB862D** — Separate training, validation, test, red-team, and audit sets sufficiently to prevent leakage and optimistic evaluation.
- [ ] **USEQ-49A7E146** — Detect benchmark contamination, near duplicates, memorized test examples, and evaluation-set tuning.
- [ ] **USEQ-7F84BAAD** — Protect data pipelines, labeling systems, feature stores, retrieval stores, and feedback channels against poisoning, unauthorized change, and provenance loss.
- [ ] **USEQ-F61A6BFD** — Record deletions, corrections, consent withdrawals, legal holds, and source restrictions, and propagate them to derived datasets and retrieval systems as required.
- [ ] **USEQ-41415588** — Define whether and how a trained artifact can honor deletion, unlearning, suppression, retraining, or future exclusion obligations.
- [ ] **USEQ-5B1BC6E2** — Prevent prompts, conversations, feedback, traces, and outputs from becoming training data by default unless the purpose, notice, permission, retention, and safeguards permit it.
- [ ] **USEQ-A764BBDF** — Review public or scraped data for access conditions, privacy expectations, robots or contractual restrictions, intellectual-property risk, harmful content, and vulnerable-population impact.
- [ ] **USEQ-E13A6AE0** — Document unavoidable data gaps and restrict deployment or confidence claims accordingly.

#### Model, component, and provider supply chain

- [ ] **USEQ-9AE49D07** — Inventory every foundation model, adapter, checkpoint, tokenizer, embedding model, reranker, safety model, library, runtime, dataset, evaluation tool, and hosted service in the deployed system.
- [ ] **USEQ-18F634F5** — Record publisher, source, version, immutable digest, license, provenance, release date, support status, known limitations, and approved use for each component.
- [ ] **USEQ-CB45E043** — Obtain models and artifacts through authenticated, integrity-checked, approved channels.
- [ ] **USEQ-ED7BBF05** — Scan model files and supporting artifacts for unsafe serialization, executable payloads, secrets, malware, unexpected code, and incompatible licenses.
- [ ] **USEQ-71ADB251** — Treat model cards, benchmark claims, certifications, safety reports, and provider assurances as evidence inputs rather than proof of local fitness.
- [ ] **USEQ-3F6CA513** — Evaluate a provider's security, privacy, data use, retention, training, isolation, availability, change-control, incident, audit, deletion, geographic, support, and exit terms.
- [ ] **USEQ-27B05EE5** — Ensure provider contracts address input and output rights, confidential data, model improvement, subprocessors, breach notification, service changes, deletion, portability, and continuity.
- [ ] **USEQ-122CF489** — Prevent a provider from changing a model, safety policy, retention behavior, endpoint, price, or material capability without detection and impact review.
- [ ] **USEQ-8B70FE23** — Pin model and API versions when stability is required and monitor announced and silent behavior changes.
- [ ] **USEQ-FD546F9D** — Maintain an alternative, fallback, degraded mode, or orderly shutdown path for critical provider dependence.
- [ ] **USEQ-6B13F795** — Evaluate concentration, geopolitical, export-control, sanctions, availability, lock-in, and proprietary-format risks.
- [ ] **USEQ-6518E104** — Track vulnerabilities and security advisories for models, runtimes, orchestration libraries, data tools, vector stores, accelerators, and inference infrastructure.
- [ ] **USEQ-CA69E6B9** — Require signed or otherwise verifiable release artifacts and maintain a software and model bill of materials proportionate to risk.
- [ ] **USEQ-DFDEDE7C** — Control access to model registries and prevent unauthorized replacement, rollback, deletion, promotion, or public publication.
- [ ] **USEQ-C493B4D5** — Separate experimental, internal, customer-specific, and production models and credentials.
- [ ] **USEQ-6A22F149** — Review community models, plugins, tools, skills, prompts, and agent packages as untrusted supply-chain components.

#### Experimentation and reproducibility

- [ ] **USEQ-7F60FE82** — State a falsifiable hypothesis, baseline, target metric, guardrails, population, and decision rule before a material experiment.
- [ ] **USEQ-73B85E38** — Record code, data, model, prompt, hyperparameter, seed, environment, hardware, dependency, compiler, runtime, and configuration versions.
- [ ] **USEQ-F2B8DD90** — Make material experiments reproducible within defined tolerances or document why nondeterminism prevents exact reproduction.
- [ ] **USEQ-37A37302** — Capture random seeds and sources of nondeterminism without falsely implying that one seed establishes robustness.
- [ ] **USEQ-9D121A7B** — Use representative baselines, including a no-model or simpler-model baseline where relevant.
- [ ] **USEQ-965D57C2** — Prevent test-set and production-outcome leakage into model selection.
- [ ] **USEQ-675A6565** — Predefine primary, secondary, safety, fairness, privacy, cost, and operational metrics.
- [ ] **USEQ-EEEFE531** — Use sample sizes, uncertainty intervals, repeated runs, statistical tests, and multiple-comparison controls appropriate to the decision.
- [ ] **USEQ-7E7D22B7** — Report negative, inconclusive, and failed experiments rather than retaining only favorable results.
- [ ] **USEQ-FFA16547** — Record experiment lineage and approval so a production artifact can be traced to its evidence.
- [ ] **USEQ-8EA244EE** — Review whether offline metrics predict real user outcomes before relying on them.
- [ ] **USEQ-49F0F65F** — Use controlled online experiments only when exposure, consent, monitoring, stopping, rollback, and harm controls are adequate.
- [ ] **USEQ-8EF29710** — Prevent an experiment from changing a consequential decision rule without the required approval and notice.
- [ ] **USEQ-8D4A14CC** — Test interactions among model, prompt, retrieval, tool, policy, UI, latency, and fallback rather than evaluating each in isolation only.
- [ ] **USEQ-2E44B61A** — Re-run affected evaluations after any material data, model, prompt, tool, provider, infrastructure, or policy change.
- [ ] **USEQ-7EB60DD2** — Retain enough artifacts and metadata to investigate a later incident without retaining unnecessary sensitive data.

#### Training, fine-tuning, and model construction

- [ ] **USEQ-CB5CDE0E** — Define the training objective and verify that optimizing it does not create unacceptable proxy behavior or incentives.
- [ ] **USEQ-9257FF1A** — Validate preprocessing, tokenization, normalization, augmentation, feature generation, and label transformations independently.
- [ ] **USEQ-ED1D1451** — Prevent train-serving skew by sharing validated transformation definitions or verifying exact equivalence.
- [ ] **USEQ-B7458F12** — Use bounded, reviewed, and observable training jobs with resource, cost, security, and privacy limits.
- [ ] **USEQ-75F53D23** — Protect training credentials, checkpoints, intermediate artifacts, logs, and distributed communication.
- [ ] **USEQ-2C07474B** — Validate checkpoint integrity before resume, fine-tuning, merge, quantization, conversion, or deployment.
- [ ] **USEQ-AF891A4E** — Review fine-tuning and preference data for poisoning, hidden triggers, sensitive content, policy conflict, and representation gaps.
- [ ] **USEQ-D244FD2A** — Measure catastrophic forgetting, capability regression, safety regression, and unexpected transfer after fine-tuning or adaptation.
- [ ] **USEQ-B93196E6** — Control model merging, distillation, compression, pruning, quantization, and conversion as material changes requiring renewed evaluation.
- [ ] **USEQ-30D5B880** — Document hyperparameter search space, stopping criteria, selection rule, compute budget, and discarded candidates.
- [ ] **USEQ-9CA4ABA8** — Prevent production selection based on one favorable run or an unrepresentative benchmark.
- [ ] **USEQ-CD2A10DD** — Assess memorization and extraction risk, especially for rare, identifying, confidential, or copyrighted examples.
- [ ] **USEQ-04E2B43F** — Apply privacy-enhancing training methods when required and validate their actual utility and privacy tradeoffs.
- [ ] **USEQ-2454E46F** — Track environmental and resource impact where training or inference is material and consider smaller or more efficient alternatives.
- [ ] **USEQ-188AEBC2** — Prevent training infrastructure from writing unrestricted artifacts to public or weakly controlled storage.
- [ ] **USEQ-1584DD7E** — Review generated training data for model collapse, inherited errors, homogenization, attribution loss, and circular evaluation.
- [ ] **USEQ-32D1AB90** — Validate safety tuning in the actual supported languages and contexts rather than assuming transfer from a dominant language.

#### Evaluation design and quality assurance

- [ ] **USEQ-B5FE544A** — Create an evaluation plan derived from intended use, affected parties, threat model, hazards, quality attributes, and failure consequences.
- [ ] **USEQ-0A64DD0D** — Evaluate correctness, usefulness, robustness, reliability, calibration, uncertainty, safety, security, privacy, fairness, accessibility, latency, throughput, cost, and environmental impact as applicable.
- [ ] **USEQ-F16359EE** — Use task-representative, current, sufficiently difficult, independently reviewed evaluation sets.
- [ ] **USEQ-07D4E2DB** — Include normal, boundary, malformed, ambiguous, adversarial, multilingual, accessibility-related, rare, and high-consequence cases.
- [ ] **USEQ-69C50645** — Measure per-group and per-context performance rather than relying only on aggregate averages.
- [ ] **USEQ-771A15A8** — Report distribution, confidence intervals, worst-case behavior, and failure examples in addition to a headline score.
- [ ] **USEQ-5FDC029A** — Evaluate both false-positive and false-negative harm and the distribution of burden created by errors.
- [ ] **USEQ-356476AA** — Calibrate scores and confidence claims against observed outcome frequencies where probabilities are presented or used.
- [ ] **USEQ-03C47BD1** — Test abstention and escalation quality, including whether the model abstains disproportionately or at the wrong times.
- [ ] **USEQ-F1649342** — Use blinded or independent human evaluation where subjective judgment is material.
- [ ] **USEQ-4F510C10** — Define evaluator qualifications, rubrics, disagreement resolution, inter-rater reliability, and conflict controls.
- [ ] **USEQ-13A7712A** — Protect evaluation workers from harmful content and provide suitable training, support, and exposure limits.
- [ ] **USEQ-C8AE2073** — Assess accessibility of AI-generated interfaces, summaries, captions, descriptions, code, documents, and interaction patterns.
- [ ] **USEQ-C6205AB5** — Validate factuality using authoritative references and source-grounding checks suitable to the domain.
- [ ] **USEQ-230C12A5** — Evaluate citation correctness, source relevance, quote accuracy, and whether generated references actually support claims.
- [ ] **USEQ-BCC4C62F** — Test long-context behavior, context ordering, truncation, conflicting sources, stale retrieval, and distractors when relevant.
- [ ] **USEQ-6E978130** — Evaluate multi-turn consistency, memory boundaries, identity continuity, refusal persistence, and recovery from earlier mistakes.
- [ ] **USEQ-DF342A34** — Test the integrated product under production-like latency, concurrency, quotas, failures, moderation, and UI conditions.
- [ ] **USEQ-999B68D7** — Use adversarial testing and red teams independent from the builders for systems with material impact.
- [ ] **USEQ-E5EDCDBA** — Track evaluation-set version, model version, configuration, prompt, tool permissions, and result provenance.
- [ ] **USEQ-0EF2BECE** — Set release thresholds and no-go failures before observing the final results where practical.
- [ ] **USEQ-0AC929A5** — Retain representative failure cases as regression tests without exposing sensitive or harmful data unnecessarily.

#### Fairness, inclusion, accessibility, and human factors

- [ ] **USEQ-B4BBA5B4** — Identify protected, vulnerable, underrepresented, and contextually disadvantaged groups relevant to the actual use.
- [ ] **USEQ-25746967** — Assess whether the problem definition, target, label, proxy, data, interface, deployment, or feedback loop can create inequitable outcomes.
- [ ] **USEQ-F6FE00F8** — Select fairness concepts and metrics that correspond to the concrete harm rather than applying a generic parity metric blindly.
- [ ] **USEQ-313223A4** — Document incompatible fairness objectives and who authorized the chosen tradeoff.
- [ ] **USEQ-A2551695** — Evaluate intersectional and small-group outcomes where aggregate categories can hide harm.
- [ ] **USEQ-9F00F826** — Investigate missingness, lower-quality data, language coverage, assistive-technology use, and differing access conditions as sources of unequal error.
- [ ] **USEQ-CF63258B** — Test whether thresholds, abstention, escalation, identity checks, fraud controls, and moderation burden groups differently.
- [ ] **USEQ-0A97658F** — Ensure human reviewers receive enough context and training to correct rather than amplify model bias.
- [ ] **USEQ-9D4D9C88** — Avoid presenting protected-trait inference, emotion inference, personality inference, or risk scoring as objective fact when scientific validity is weak or context dependent.
- [ ] **USEQ-D8207FA6** — Do not use accessibility needs, disability signals, language, device quality, or behavioral differences as harmful proxies.
- [ ] **USEQ-10FA0F4D** — Make explanations and recourse understandable to the affected audience and accessible through supported modalities.
- [ ] **USEQ-8298717B** — Provide reasonable alternatives when voice, vision, motion, cognition, literacy, language, or device assumptions exclude users.
- [ ] **USEQ-F3A9EE15** — Test whether generated descriptions, captions, summaries, translations, and recommendations preserve essential meaning for disabled users.
- [ ] **USEQ-F38A1B69** — Include affected communities and domain experts in design and evaluation proportionate to risk.
- [ ] **USEQ-4EA0B410** — Monitor disparate impact and complaint patterns after deployment, not only before launch.
- [ ] **USEQ-84239BB3** — Prevent personalization from creating discriminatory prices, eligibility, reach, service quality, or information access without lawful and justified grounds.

#### Robustness, reliability, and uncertainty

- [ ] **USEQ-EAECD5B4** — Define the validated operating domain, input ranges, languages, modalities, task types, and environmental conditions.
- [ ] **USEQ-07AA5338** — Detect or safely handle out-of-distribution, corrupted, adversarial, incomplete, stale, and conflicting input.
- [ ] **USEQ-6EF822C9** — Measure sensitivity to paraphrase, ordering, formatting, spelling, noise, prompt framing, irrelevant context, and minor data changes.
- [ ] **USEQ-8607FB55** — Evaluate stability across repeated runs and quantify nondeterminism where it affects user expectations or decisions.
- [ ] **USEQ-25B208C7** — Use deterministic validation for formats, ranges, permissions, calculations, business rules, and side effects that must be exact.
- [ ] **USEQ-1C0FC607** — Do not use model confidence as a substitute for calibrated uncertainty or external validation.
- [ ] **USEQ-26D2EEF9** — Define fallback behavior for model timeout, quota exhaustion, provider outage, safety-system outage, retrieval failure, tool failure, malformed output, and policy conflict.
- [ ] **USEQ-5AEB2EF3** — Ensure degraded modes fail visibly and safely rather than silently substituting lower-quality or differently governed behavior.
- [ ] **USEQ-4C760B55** — Prevent retries from creating duplicate external actions or unbounded cost.
- [ ] **USEQ-553B4185** — Use idempotency, transaction boundaries, compensation, and confirmation for model-initiated side effects.
- [ ] **USEQ-11D17067** — Bound context, output, recursion, tool calls, steps, time, memory, compute, network, and monetary cost.
- [ ] **USEQ-50EEA27A** — Test long-running and agentic workflows for drift from the original goal, accumulating error, and unsafe continuation.
- [ ] **USEQ-5E9067FC** — Verify that a human can pause, inspect, correct, resume, and terminate consequential workflows.
- [ ] **USEQ-5EDE5164** — Monitor upstream data freshness and distinguish missing evidence from evidence of absence.
- [ ] **USEQ-9131BDD3** — Prevent stale models, embeddings, policies, prompts, or retrieval indexes from being treated as current without detection.
- [ ] **USEQ-270F3C01** — Validate behavior after dependency, hardware, quantization, compiler, runtime, locale, and infrastructure changes.
- [ ] **USEQ-98682F5F** — Define recovery and reconciliation for incomplete AI-driven workflows.

#### AI and model security

- [ ] **USEQ-9DF37D00** — Threat-model model theft, extraction, inversion, membership inference, training-data leakage, poisoning, evasion, adversarial examples, hidden triggers, unsafe serialization, and supply-chain compromise.
- [ ] **USEQ-15EDD801** — Restrict access to training data, model weights, checkpoints, system prompts, evaluation sets, embeddings, safety policies, and privileged endpoints according to sensitivity.
- [ ] **USEQ-B15DC6F6** — Use separate identities and least privilege for training, evaluation, deployment, retrieval, monitoring, and administration.
- [ ] **USEQ-B7EDE7E8** — Prevent public or tenant-scoped inputs from crossing user, tenant, region, or confidentiality boundaries through caches, memory, retrieval, logs, fine-tuning, or provider behavior.
- [ ] **USEQ-6BAFF5D7** — Apply authentication, authorization, rate limits, quotas, anomaly detection, and abuse controls to inference and model-management endpoints.
- [ ] **USEQ-EBFC942E** — Protect model and prompt intellectual property without relying on secrecy as the only safety control.
- [ ] **USEQ-73282252** — Test for extraction through repeated querying, probability outputs, embeddings, error messages, and side channels.
- [ ] **USEQ-7E587625** — Minimize output detail, logits, confidence data, or debug information when it materially increases extraction or inference risk.
- [ ] **USEQ-987DF294** — Protect feedback and reinforcement channels from coordinated manipulation and sybil attacks.
- [ ] **USEQ-37715BDE** — Validate uploaded models, adapters, datasets, notebooks, prompts, tools, and serialized artifacts in isolated environments.
- [ ] **USEQ-DE3A04CA** — Prevent training or inference code from executing untrusted content with host, cloud, network, or secret access.
- [ ] **USEQ-FC4C23B5** — Monitor unexpected model size, digest, registry, permission, network, cost, latency, refusal, and output-distribution changes.
- [ ] **USEQ-8A02B572** — Use tamper-evident provenance and promotion records for high-impact models.
- [ ] **USEQ-E82328B8** — Red-team cross-tenant leakage, secret elicitation, policy bypass, encoded attacks, multilingual attacks, indirect instruction, and persistence.
- [ ] **USEQ-1FA84F69** — Include AI systems in vulnerability management, incident response, forensic readiness, backup, recovery, and supplier-risk processes.
- [ ] **USEQ-C9DEAD72** — Define a coordinated disclosure path for model, prompt, evaluation, data, and agent vulnerabilities.
- [ ] **USEQ-13602AAE** — Prevent security filters from becoming the sole authorization or data-loss-prevention boundary.

#### Generative AI, retrieval, prompts, tools, and agents

- [ ] **USEQ-AF20A1A9** — Treat system instructions, developer instructions, user prompts, retrieved documents, web pages, files, tool results, memory, metadata, and model outputs as inputs with distinct trust levels.
- [ ] **USEQ-82711812** — Enforce instruction priority, authorization, data boundaries, and tool policy in trusted code rather than relying only on natural-language instructions.
- [ ] **USEQ-3CC8B60E** — Assume direct and indirect prompt injection can occur and design so successful injection cannot exceed independently enforced authority.
- [ ] **USEQ-0D2EF3BB** — Keep secrets, credentials, hidden policies, unrelated user data, internal reasoning artifacts, and privileged context out of model context unless strictly necessary and protected.
- [ ] **USEQ-9066469A** — Sanitize and label retrieved content without treating sanitization as a complete prompt-injection defense.
- [ ] **USEQ-2B691284** — Authorize retrieval before content enters the model and recheck authorization before output or side effects.
- [ ] **USEQ-341FBC64** — Preserve tenant, user, purpose, region, retention, and sensitivity metadata through chunking, embedding, indexing, retrieval, caching, and citation.
- [ ] **USEQ-8B6857A6** — Validate retrieval freshness, source identity, relevance, authority, coverage, contradiction, and permission.
- [ ] **USEQ-25F2FD1D** — Provide citations or source traces when needed, and verify that every cited source exists and supports the associated claim.
- [ ] **USEQ-B5516B69** — Distinguish model-generated statements from retrieved or deterministic facts in the product architecture and user interface where material.
- [ ] **USEQ-EA11C0C3** — Constrain structured output with schema validation and reject, repair, or escalate malformed output safely.
- [ ] **USEQ-309DD707** — Encode or sanitize generated output for its destination context before rendering, executing, storing, sending, querying, or compiling it.
- [ ] **USEQ-17CCC437** — Never execute generated code, commands, queries, templates, URLs, configurations, or workflow definitions without independent policy, validation, isolation, and approval appropriate to impact.
- [ ] **USEQ-DE757839** — Give each tool a narrow identity, explicit operation set, parameter schema, resource scope, budget, timeout, and audit trail.
- [ ] **USEQ-8BFD9B81** — Require renewed authorization for high-impact tool calls and prevent the model from self-granting permissions.
- [ ] **USEQ-191B4C8D** — Use previews and human confirmation for destructive, financial, public, legal, security, privacy, account, infrastructure, or irreversible actions.
- [ ] **USEQ-7F7AA7C2** — Limit agent recursion, delegation, parallelism, child agents, self-modification, memory, context expansion, and external communication.
- [ ] **USEQ-D296B587** — Prevent agents from creating hidden persistence, accounts, credentials, scheduled work, infrastructure, or subscriptions without explicit authority.
- [ ] **USEQ-3425618B** — Isolate browser, interpreter, file, network, code-execution, and computer-control environments and reset them between users or tasks as required.
- [ ] **USEQ-5212224D** — Restrict outbound network destinations and defend against server-side request forgery, data exfiltration, malicious downloads, and unsafe redirects.
- [ ] **USEQ-6AC50F7F** — Treat external tool output as untrusted even when the tool itself is approved.
- [ ] **USEQ-773557D7** — Verify that memory creation, retrieval, correction, expiry, deletion, and cross-session use are transparent and authorized.
- [ ] **USEQ-11C0DFF6** — Do not store sensitive conversation memory merely because personalization is convenient.
- [ ] **USEQ-C3984280** — Test jailbreaks, prompt leakage, policy conflict, role confusion, multi-turn escalation, context poisoning, tool misuse, hidden instructions, and malicious files.
- [ ] **USEQ-F2FE87FF** — Ensure refusals do not reveal prohibited material and do not prevent legitimate safety, accessibility, or recovery use without recourse.
- [ ] **USEQ-FE69F909** — Moderate both inputs and outputs according to the actual harm model while retaining the minimum necessary data.
- [ ] **USEQ-952B24B4** — Provide a kill switch for tool use and a separate capability to disable the model while preserving essential deterministic product functions.

#### MLOps, deployment, and release engineering

- [ ] **USEQ-F3DA7D88** — Use version control and review for training code, inference code, prompts, policies, feature definitions, data transformations, evaluation logic, deployment manifests, and infrastructure.
- [ ] **USEQ-13FDD0E8** — Use immutable identifiers for models, datasets, prompts, retrieval indexes, feature definitions, policies, and evaluation suites.
- [ ] **USEQ-CFD67264** — Trace every production response or decision to the relevant deployed versions where impact and privacy permit.
- [ ] **USEQ-EDEEF0CC** — Separate development, experimentation, validation, shadow, canary, and production environments and data access.
- [ ] **USEQ-60454DA3** — Require reproducible, authenticated, integrity-checked build and promotion paths for models and associated artifacts.
- [ ] **USEQ-7C1585D4** — Prevent direct manual replacement of a production model or prompt without an auditable emergency process.
- [ ] **USEQ-3BC056D6** — Validate feature and preprocessing equivalence between training and serving.
- [ ] **USEQ-1EEA8D0C** — Use a model registry or equivalent controlled inventory with stages, owners, evidence, approvals, and retirement state.
- [ ] **USEQ-9CCF6F87** — Package required tokenizer, preprocessing, schema, policy, calibration, and compatibility metadata with the model release.
- [ ] **USEQ-767F858F** — Test backward and forward compatibility among model, client, feature, schema, retrieval, prompt, tool, and output consumers.
- [ ] **USEQ-004EB3A6** — Use shadow evaluation, canary cohorts, champion-challenger comparison, or staged rollout where impact warrants it.
- [ ] **USEQ-8002BB09** — Define automated and human rollout gates for quality, safety, security, fairness, latency, throughput, availability, cost, and business integrity.
- [ ] **USEQ-E63926B7** — Maintain a previous known-good model, prompt, policy, and retrieval index when rollback is feasible.
- [ ] **USEQ-67F31DA5** — Test rollback and roll-forward, including state, memory, embeddings, cache, and schema compatibility.
- [ ] **USEQ-39098812** — Handle provider or model nondeterminism when comparing canary and control populations.
- [ ] **USEQ-2EEA3BC7** — Prevent online learning, feedback updates, and self-modification from entering production without bounded, reviewed promotion.
- [ ] **USEQ-683E2947** — Validate accelerator, runtime, compiler, quantization, serving engine, batching, and hardware changes as model changes when output can differ.
- [ ] **USEQ-37EF3EDD** — Capacity-test inference, retrieval, safety filters, tool services, queues, and provider quotas under peak and abuse load.
- [ ] **USEQ-58DDDF7C** — Use admission control, prioritization, backpressure, caching, batching, and load shedding without violating isolation, freshness, or user safety.
- [ ] **USEQ-014BF5FE** — Monitor cost per task, token or compute amplification, runaway agents, retry storms, and unexpectedly expensive inputs.
- [ ] **USEQ-45CF19C0** — Include AI components in deployment records, configuration drift detection, disaster recovery, and production-readiness evidence.

#### Production monitoring, drift, and feedback loops

- [ ] **USEQ-FFCBD920** — Define production indicators for input quality, output quality, safety, security, fairness, refusal, calibration, latency, availability, saturation, cost, and user outcomes.
- [ ] **USEQ-81C152F9** — Monitor input, feature, label, concept, population, retrieval, embedding, output, policy, and provider drift as applicable.
- [ ] **USEQ-9666549A** — Use reference distributions and thresholds appropriate to each signal; do not treat every statistical difference as harmful drift.
- [ ] **USEQ-F280064F** — Detect silent upstream schema, semantic, unit, category, timestamp, and missingness changes.
- [ ] **USEQ-30EC58C5** — Measure quality using delayed ground truth, expert review, user correction, audits, reconciliation, and representative sampling where direct labels are unavailable.
- [ ] **USEQ-99CE268F** — Protect feedback metrics against gaming, selection bias, survivorship bias, harassment, coordinated manipulation, and automation.
- [ ] **USEQ-5C19854F** — Separate user preference from factual correctness, safety, legal compliance, and protected rights.
- [ ] **USEQ-5D28AC18** — Monitor by relevant language, region, device, group, task, tenant, model version, prompt version, provider, and failure mode.
- [ ] **USEQ-27A8B5D5** — Track harmful exposure and downstream impact, not only the number of generated outputs or blocked requests.
- [ ] **USEQ-627BB632** — Detect data leakage, repeated sensitive strings, anomalous citations, prohibited tool use, unexpected network destinations, and cross-user similarity where appropriate.
- [ ] **USEQ-1257C1CC** — Monitor model and provider changes using recurring golden tests, behavioral fingerprints, and change-point analysis.
- [ ] **USEQ-DA805FFF** — Ensure monitoring itself does not collect excessive prompts, outputs, personal data, secrets, or copyrighted content.
- [ ] **USEQ-A8219EEA** — Provide protected sampling and access controls for human review of production interactions.
- [ ] **USEQ-CB3AFE81** — Define when drift triggers investigation, recalibration, retraining, policy change, rollback, restricted scope, or suspension.
- [ ] **USEQ-5795AC64** — Track unresolved model limitations, recurring user corrections, escalations, appeals, incidents, and near misses.
- [ ] **USEQ-078062DD** — Close the loop by converting production failures into data fixes, evaluation cases, product changes, safety controls, documentation, and regression tests.
- [ ] **USEQ-F4627883** — Review whether optimization feedback loops amplify popularity, bias, fraud, polarization, low-quality content, or exclusion.
- [ ] **USEQ-F7E09645** — Maintain an external reporting path for affected users, researchers, customers, and security reporters.

#### Human oversight, explanations, and recourse

- [ ] **USEQ-61953451** — Define the exact role, authority, competence, workload, information, and response time required of human reviewers.
- [ ] **USEQ-27053648** — Ensure a human-in-the-loop control is operationally real rather than a nominal approval under impossible volume or time pressure.
- [ ] **USEQ-8410B1CC** — Provide reviewers the original input, relevant evidence, model output, uncertainty, applicable policy, prior history, and limitations needed to decide.
- [ ] **USEQ-027CBA6B** — Prevent automation bias by training reviewers, presenting uncertainty, and making independent assessment possible.
- [ ] **USEQ-5283D0F9** — Measure reviewer disagreement, override, fatigue, consistency, delay, and error.
- [ ] **USEQ-372E6D8F** — Escalate cases beyond a reviewer's competence or authority.
- [ ] **USEQ-09AA097A** — Record consequential human overrides and reasons without discouraging justified correction.
- [ ] **USEQ-EA629B07** — Provide explanations that are faithful to the actual decision mechanism and useful for the user's action or appeal.
- [ ] **USEQ-FF8035D0** — Do not present generated rationales as causal explanations unless their fidelity has been established.
- [ ] **USEQ-3836C247** — Explain material data categories, factors, thresholds, source limitations, and human involvement at an appropriate level.
- [ ] **USEQ-64875AE5** — Provide meaningful notice and recourse without exposing trade secrets, security controls, other people's data, or evasion-sensitive detail unnecessarily.
- [ ] **USEQ-1DC6F83B** — Allow correction of inaccurate source data and propagate corrections to future decisions where applicable.
- [ ] **USEQ-A3805124** — Ensure appeals are reviewed by a qualified person or independent mechanism and are not decided solely by the same unchanged model.
- [ ] **USEQ-82AF898D** — Track restoration and remediation after a successful appeal, including downstream records and notifications.
- [ ] **USEQ-906B5981** — Provide additional safeguards where users cannot reasonably understand, contest, or avoid the system.

#### Transparency, documentation, and assurance

- [ ] **USEQ-4E890ABE** — Maintain system documentation describing purpose, architecture, data, models, prompts, tools, operating domain, metrics, limitations, risks, controls, owners, and change history.
- [ ] **USEQ-57485A3D** — Maintain model, data, prompt, retrieval, evaluation, and deployment documentation proportionate to impact.
- [ ] **USEQ-8E763CF7** — Document known failure modes and conditions under which outputs must not be relied upon.
- [ ] **USEQ-268830F2** — State whether outputs are generated, predicted, retrieved, transformed, or deterministically calculated where the distinction is material.
- [ ] **USEQ-85E488CB** — Document source attribution, intellectual-property treatment, data retention, provider use, human review, and complaint paths.
- [ ] **USEQ-C56A451D** — Publish appropriate user-facing transparency without revealing security-sensitive implementation detail.
- [ ] **USEQ-140A1C57** — Make contractual and marketing claims consistent with measured evidence and supported versions.
- [ ] **USEQ-CD81C32A** — Do not imply certification, neutrality, fairness, accuracy, safety, explainability, or human equivalence beyond the actual assessment scope.
- [ ] **USEQ-33B02D9F** — Maintain an auditable mapping from risks to controls, tests, monitoring, incidents, and responsible owners.
- [ ] **USEQ-275D0D2D** — Use independent assessment, conformity evaluation, external audit, red teaming, or expert review where consequence and obligation warrant it.
- [ ] **USEQ-D4ABA242** — Retain evidence sufficient to reproduce or reconstruct consequential releases and decisions for the required period.
- [ ] **USEQ-A687DF4E** — Review documentation whenever the system changes and make stale documents visibly invalid.

#### AI-assisted software engineering

- [ ] **USEQ-7E3E586D** — Define approved AI coding, review, documentation, testing, design, operations, and support tools and the data classifications each may receive.
- [ ] **USEQ-D13D7904** — Prevent source code, secrets, credentials, customer data, incident evidence, vulnerabilities, proprietary designs, and regulated data from being submitted to unapproved services.
- [ ] **USEQ-DC32CD99** — Configure provider retention, training, sharing, logging, regional processing, and access settings according to organizational policy.
- [ ] **USEQ-0E18CBD8** — Treat generated code, tests, configurations, commands, infrastructure, migrations, dependencies, licenses, and documentation as untrusted contributions requiring accountable human review.
- [ ] **USEQ-9A66039A** — Make the human author responsible for understanding, validating, maintaining, and defending accepted AI-generated changes.
- [ ] **USEQ-482FDF4A** — Require the same requirement, architecture, security, privacy, accessibility, test, review, and release controls for generated and manually written work.
- [ ] **USEQ-CAE04E9B** — Do not merge, deploy, execute, publish, or apply generated output solely because it compiles, passes a narrow test, or appears plausible.
- [ ] **USEQ-700BB2C2** — Verify APIs, library names, versions, options, standards references, citations, and commands against authoritative sources.
- [ ] **USEQ-7221EBC7** — Check generated code for invented functions, deprecated practices, insecure defaults, hidden side effects, excessive privilege, weak error handling, and incorrect concurrency.
- [ ] **USEQ-91CFE8F2** — Review provenance, license, attribution, trade-secret, and substantial-similarity risks before accepting generated material.
- [ ] **USEQ-3A1BED2B** — Do not let generated tests merely reproduce the same misunderstanding as generated implementation; derive tests independently from requirements and invariants.
- [ ] **USEQ-3654A2CD** — Use adversarial, boundary, mutation, integration, and production-like validation for generated changes proportionate to impact.
- [ ] **USEQ-86CAAD7F** — Require explicit review for authentication, authorization, cryptography, payments, privacy, data deletion, migrations, infrastructure, safety, and incident code.
- [ ] **USEQ-CC9F4CF0** — Prevent agents from committing, opening privileged changes, modifying policies, rotating secrets, deploying, deleting resources, or changing production without scoped identity, policy gates, audit, and required approval.
- [ ] **USEQ-1D87F2FE** — Sandbox generated commands and code; inspect effects before execution and avoid granting broad filesystem, network, cloud, repository, or secret access.
- [ ] **USEQ-1AE8D9E6** — Keep AI tool instructions, repository context files, reusable prompts, and agent policies versioned, reviewed, and protected from untrusted content.
- [ ] **USEQ-54BE67BA** — Treat issue text, documentation, dependency files, comments, test output, web pages, and repository content as possible indirect instructions to coding agents.
- [ ] **USEQ-1765CB7B** — Measure AI-assisted development using quality, reliability, security, review effort, rework, learning, flow, and outcomes rather than generated volume or apparent speed alone.
- [ ] **USEQ-8F607BCB** — Monitor whether AI use increases change size, dependency sprawl, duplicated code, shallow understanding, review burden, incidents, or knowledge concentration.
- [ ] **USEQ-29A64D47** — Provide training on verification, security, privacy, licensing, prompt injection, overreliance, and tool limitations.
- [ ] **USEQ-6A8DD723** — Preserve meaningful authorship and review accountability in commit, approval, and incident records.
- [ ] **USEQ-7071707F** — Allow workers to report unsafe or low-quality AI pressure and do not mandate use where it reduces quality or accessibility.
- [ ] **USEQ-0B94F8FC** — Regularly reassess whether the tool remains beneficial as models, terms, behavior, data practices, and organizational needs change.

#### AI incidents, recovery, and retirement

- [ ] **USEQ-0993C353** — Define incident categories for harmful output, data leakage, unauthorized action, model compromise, poisoning, provider change, systemic bias, safety-control failure, model theft, fraud, and widespread misinformation.
- [ ] **USEQ-DE662B2B** — Make incident detection, declaration, containment, evidence preservation, communication, recovery, and post-incident review part of the ordinary response program.
- [ ] **USEQ-49974079** — Retain the minimum necessary input, output, model, prompt, tool, retrieval, policy, version, identity, and timeline evidence to investigate material events.
- [ ] **USEQ-0DD20497** — Enable rapid disablement of a model, prompt, tool, data source, retrieval collection, tenant, capability, or provider without destroying unrelated essential service.
- [ ] **USEQ-19BC3298** — Revoke credentials, isolate endpoints, freeze artifacts, preserve provenance, and stop feedback ingestion during suspected compromise.
- [ ] **USEQ-89D0100E** — Identify all users, decisions, content, records, and downstream systems affected by a faulty release or compromised source.
- [ ] **USEQ-3251BED0** — Provide correction, reprocessing, notification, restoration, compensation, appeal, or human review appropriate to the impact.
- [ ] **USEQ-972E5264** — Test rollback to a known-good model and verify that memories, caches, embeddings, queues, and prior unsafe outputs do not persist unnoticed.
- [ ] **USEQ-A7317A41** — Define retraining and reevaluation criteria after incidents; do not resume solely because the immediate symptom disappeared.
- [ ] **USEQ-9E544317** — Retire obsolete models, prompts, datasets, indexes, feature definitions, credentials, endpoints, and provider access.
- [ ] **USEQ-21F4B26E** — Preserve required records while deleting data and artifacts no longer authorized or needed.
- [ ] **USEQ-C1801C4C** — Prevent retired artifacts from being accidentally redeployed through caches, registries, old clients, automation, or disaster recovery.
- [ ] **USEQ-85F03B12** — Track successor compatibility and communicate material behavior changes to users and downstream consumers.
- [ ] **USEQ-6F0FBC39** — Review whether retirement changes legal holds, reproducibility, audit, deletion, or long-term support obligations.

#### AI release blockers

- [ ] **USEQ-6EDFEFA5** — Do not deploy an AI system whose intended purpose, prohibited uses, affected population, decision authority, data provenance, and accountable owner are unknown.
- [ ] **USEQ-28269E66** — Do not deploy a consequential system without representative evaluation, defined thresholds, known limitations, recourse, monitoring, and an effective stop mechanism.
- [ ] **USEQ-8CED23FB** — Do not treat a provider benchmark, model card, certification, demo, or automated evaluation as sufficient local assurance.
- [ ] **USEQ-B22A2372** — Do not allow model output to authorize itself, bypass deterministic access control, or directly trigger high-impact action without independent controls.
- [ ] **USEQ-926988B6** — Do not expose secrets, unrelated user data, privileged system context, or unauthorized retrieval content to a model.
- [ ] **USEQ-FDE77A89** — Do not deploy when material test leakage, benchmark contamination, unresolved cross-user leakage, poisoning, model compromise, or unsafe artifact provenance is suspected.
- [ ] **USEQ-7D42569E** — Do not deploy a system that cannot abstain, degrade, or hand off safely outside its validated operating domain.
- [ ] **USEQ-88BA6A09** — Do not use a human-review label when reviewers lack the authority, competence, context, capacity, or time to prevent harm.
- [ ] **USEQ-28F27AB9** — Do not accept fairness, safety, privacy, or security performance based only on aggregate results that conceal material subgroup or high-consequence failures.
- [ ] **USEQ-24C04A1C** — Do not allow online learning, autonomous self-modification, or unrestricted agent capability without bounded governance, evaluation, authorization, and recovery.
- [ ] **USEQ-5D5D0DDC** — Do not use AI-generated code or infrastructure in production without accountable review and the same evidence required for human-authored work.
- [ ] **USEQ-5C8EEFD1** — Do not continue operation after a provider or model change invalidates the evidence package until affected controls are reassessed.

## AI-Assisted Engineering, Agentic Development, and Content Provenance

_Consolidated from `final consolidated corpus/10-ai-ml-agents-ai-assisted-engineering-provenance.md#AI-Assisted Engineering, Agentic Development, and Content Provenance`; 289 non-duplicative controls._

### AI-assisted engineering governance and permitted use

- [ ] **USEQ-2F8E6FA3** — Maintain an inventory of AI-assisted development, review, testing, documentation, design, search, operations, support, and content tools used by the organization.
- [ ] **USEQ-0E231454** — Assign an accountable owner for each approved tool, provider, model family, deployment mode, integration, and use case.
- [ ] **USEQ-5DA9BD1E** — Define permitted, restricted, and prohibited uses based on data sensitivity, code criticality, intellectual-property risk, safety, security, privacy, regulation, and customer commitments.
- [ ] **USEQ-9F2C9B84** — Require risk review before AI is used for security controls, authorization rules, cryptography, safety logic, financial calculations, migrations, production changes, legal interpretation, or other high-consequence work.
- [ ] **USEQ-3FB8AC14** — Distinguish advisory assistance from delegated decision-making and autonomous execution.
- [ ] **USEQ-7322984E** — Ensure a named human remains accountable for every accepted artifact, decision, change, deployment, and communication produced with AI assistance.
- [ ] **USEQ-26AC7D70** — Do not treat provider availability, popularity, benchmark scores, or marketing claims as assurance of suitability.
- [ ] **USEQ-3AC1CD29** — Document model and tool limitations, known failure modes, supported contexts, retention behavior, geographic processing, and provider dependencies.
- [ ] **USEQ-3B2F2EE8** — Define the minimum review and verification required for each class of AI-generated output.
- [ ] **USEQ-EE003D19** — Prohibit bypassing review, testing, separation of duties, change control, or release gates because content was generated quickly.
- [ ] **USEQ-8EF9CA6E** — Prohibit representing AI-generated work as independently authored, reviewed, tested, or certified when it was not.
- [ ] **USEQ-AB9ADB62** — Train users in hallucination, automation bias, data leakage, prompt injection, insecure code, fabricated citations, licensing uncertainty, and overreliance risks.
- [ ] **USEQ-12008C5B** — Provide an accessible process for reporting harmful, insecure, inaccurate, biased, infringing, or confidential AI output.
- [ ] **USEQ-3F338012** — Review governance after model changes, provider terms changes, incidents, new integrations, new agent capabilities, or expansion into higher-risk work.
- [ ] **USEQ-5E567DE3** — Apply stricter controls to autonomous or background agents than to read-only interactive assistants.

### Tool, model, provider, and integration inventory

- [ ] **USEQ-5F29B160** — Record the provider, model identifier, model version or update channel, hosting location, integration version, enabled features, data flows, tools, plugins, and permissions for each AI capability.
- [ ] **USEQ-746C8D31** — Record whether prompts, files, code, outputs, telemetry, feedback, embeddings, and metadata are retained or used for provider training.
- [ ] **USEQ-AC2FD44B** — Record the contractual terms, service levels, security commitments, privacy terms, subprocessors, data locations, deletion behavior, and incident-notification obligations.
- [ ] **USEQ-29597410** — Identify whether model behavior can change without an explicit customer-controlled version update.
- [ ] **USEQ-100B3A4A** — Pin or validate model versions for workflows where reproducibility and regression assurance matter.
- [ ] **USEQ-054DF3AB** — Monitor model deprecation, behavior changes, safety-policy changes, context limits, pricing changes, and regional availability.
- [ ] **USEQ-1CAE8F75** — Inventory extensions, browser integrations, editor plugins, local agents, command-line tools, code-review bots, and hidden AI features in existing products.
- [ ] **USEQ-C849B749** — Prevent unapproved tools from accessing source repositories, tickets, production data, customer content, secrets, or internal documentation.
- [ ] **USEQ-BC94E90D** — Use organization-managed identities and configuration rather than uncontrolled personal accounts for business use.
- [ ] **USEQ-D4210FCC** — Restrict tool permissions and repository access to the minimum required scope.
- [ ] **USEQ-8BFAB4A7** — Review integration webhooks, OAuth grants, tokens, local file access, network access, and command execution.
- [ ] **USEQ-81875F3C** — Ensure provider compromise or outage has an isolation, disablement, replacement, and continuity plan.
- [ ] **USEQ-58409064** — Verify that disabling the AI feature also removes active credentials, agents, callbacks, scheduled tasks, and data flows.
- [ ] **USEQ-F936BCAB** — Reconcile the approved inventory with actual network, repository, identity, procurement, and endpoint telemetry.
- [ ] **USEQ-46F439A2** — Expire temporary pilots automatically unless they receive formal approval.

### Confidentiality, privacy, data minimization, and intellectual property

- [ ] **USEQ-19F84B67** — Classify source code, prompts, tickets, designs, logs, customer data, credentials, vulnerabilities, contracts, and documents before they are supplied to an AI system.
- [ ] **USEQ-3E8114EE** — Do not submit secrets, private keys, reusable tokens, production credentials, payment data, regulated data, or unredacted sensitive personal data unless explicitly approved and technically protected.
- [ ] **USEQ-5BC689C8** — Minimize prompt and context data to what is necessary for the task.
- [ ] **USEQ-C0D20FDF** — Redact or synthesize examples where real identities, customer content, incidents, vulnerabilities, or proprietary algorithms are not necessary.
- [ ] **USEQ-3227DC4E** — Prevent tools from automatically indexing repositories, drives, conversations, or documents beyond the approved scope.
- [ ] **USEQ-949B76FD** — Verify retention, deletion, model-training, abuse-monitoring, human-review, and backup terms before supplying confidential material.
- [ ] **USEQ-393DB5A1** — Ensure data-subject rights, contractual deletion, legal holds, and retention requirements extend to prompts, outputs, embeddings, caches, and provider copies where applicable.
- [ ] **USEQ-D3455DAB** — Assess cross-border transfer and data-residency implications of provider processing.
- [ ] **USEQ-BDBA143D** — Do not assume output is free of third-party intellectual property, confidential information, personal data, or license obligations.
- [ ] **USEQ-C6ED3A2C** — Review generated code, text, media, designs, and tests for substantial similarity, copied notices, trademarks, confidential patterns, and incompatible licenses when risk warrants it.
- [ ] **USEQ-8FDF8863** — Preserve required attribution, license, and provenance information for incorporated material.
- [ ] **USEQ-A41AE117** — Prevent prompts from requesting unauthorized access, circumvention, copying, de-anonymization, or reconstruction of restricted information.
- [ ] **USEQ-A873CC6F** — Ensure provider telemetry and organizational monitoring do not create disproportionate employee surveillance.
- [ ] **USEQ-9555C8D3** — Document legitimate purpose, access, retention, and review for AI-use telemetry.
- [ ] **USEQ-69C5D0B6** — Test deletion and account termination across local caches, provider stores, embeddings, logs, and integrations.

### Prompt, context, retrieval, and instruction integrity

- [ ] **USEQ-525FCCE0** — Treat user prompts, retrieved documents, source comments, issue text, web content, dependency metadata, tool output, and model output as untrusted input.
- [ ] **USEQ-464BA3EB** — Separate system policy, developer instructions, user requests, retrieved content, and tool results so that untrusted text cannot silently become authoritative instructions.
- [ ] **USEQ-FC5D2312** — Do not rely on prompt wording alone to enforce authorization, confidentiality, financial limits, safety, or change-control policy.
- [ ] **USEQ-D1A9A000** — Validate identity, authorization, tenant, data classification, and purpose before retrieving or supplying context.
- [ ] **USEQ-22B94B7A** — Apply access control before retrieval rather than filtering unauthorized results after they have reached the model.
- [ ] **USEQ-F5A43B34** — Preserve tenant and user boundaries in indexes, embeddings, caches, conversation memory, tool calls, and logs.
- [ ] **USEQ-DCD08BD7** — Detect and handle direct and indirect prompt injection, malicious repository instructions, poisoned documentation, and adversarial tool output.
- [ ] **USEQ-5BD7D44C** — Limit the sources, domains, paths, repositories, branches, documents, and record types the system may retrieve.
- [ ] **USEQ-17A2E4C0** — Record source provenance and freshness for material facts used in generated output.
- [ ] **USEQ-AD2E1C13** — Require primary or authoritative sources for claims whose accuracy materially affects engineering or compliance decisions.
- [ ] **USEQ-0E651845** — Reject or clearly mark unsupported citations, fabricated identifiers, nonexistent APIs, invented standards, and unverifiable facts.
- [ ] **USEQ-D19D4771** — Bound context size and prioritize relevant evidence rather than truncating critical policy or requirements silently.
- [ ] **USEQ-3C9530AD** — Do not permit hidden conversation state or memory to alter high-impact behavior without visibility and control.
- [ ] **USEQ-DA4C03E8** — Make context expiration, refresh, and invalidation rules explicit after source changes or permission changes.
- [ ] **USEQ-3B2A4D90** — Test the system with conflicting instructions, malicious comments, deceptive filenames, encoded payloads, and poisoned retrieval content.

### Human accountability, review, and automation-bias controls

- [ ] **USEQ-45E865D0** — Require the accepting engineer or reviewer to understand the generated artifact well enough to explain its purpose, assumptions, risks, and failure modes.
- [ ] **USEQ-706392C8** — Reject code or configuration that no accountable maintainer can understand and support.
- [ ] **USEQ-27ADA51E** — Require independent review for high-impact AI-generated changes even when the generator also performs a review.
- [ ] **USEQ-3DF6CB04** — Ensure reviewers know which portions were generated, transformed, or suggested by AI when that knowledge affects review strategy.
- [ ] **USEQ-159E86E9** — Prevent the same AI output from serving as implementation, expected result, reviewer, and approval evidence without independent validation.
- [ ] **USEQ-1023D68F** — Use checklists and evidence prompts that direct reviewers to correctness, security, privacy, performance, accessibility, maintainability, licensing, and operational impact.
- [ ] **USEQ-C2BD1A86** — Require reviewers to inspect surrounding code and system behavior rather than only the generated diff.
- [ ] **USEQ-97B09BF8** — Verify all assumptions, external facts, API behavior, version compatibility, and cited requirements against authoritative sources.
- [ ] **USEQ-568D42D2** — Challenge plausible-looking complexity, abstractions, dependencies, and comments that lack a demonstrated need.
- [ ] **USEQ-35A8F0B6** — Require explicit approval for deletions, migrations, permission changes, security-control changes, and production commands.
- [ ] **USEQ-5B11CA85** — Use dual control for irreversible, financially consequential, safety-related, or organization-wide actions.
- [ ] **USEQ-7299FE1A** — Measure escaped defects and review failures involving AI-assisted work without using metrics to hide or discourage appropriate use.
- [ ] **USEQ-A0D8343C** — Design interfaces to show uncertainty, missing evidence, tool actions, changed files, affected systems, and unresolved warnings.
- [ ] **USEQ-00A677A3** — Prevent speed or productivity targets from encouraging rubber-stamping of generated work.
- [ ] **USEQ-EC36170F** — Retain accountability records for accepted high-impact outputs and decisions.

### AI-generated source code and implementation quality

- [ ] **USEQ-BB79B60F** — Apply every ordinary code-quality, architecture, security, privacy, performance, accessibility, and maintainability rule to AI-generated code without reduction.
- [ ] **USEQ-28217B46** — Verify that generated code satisfies the exact requirement rather than a superficially similar generic pattern.
- [ ] **USEQ-B5C04073** — Confirm names, abstractions, boundaries, data types, error behavior, state transitions, and concurrency semantics match the existing system.
- [ ] **USEQ-439A8D87** — Remove unnecessary duplication, speculative generality, dead code, excessive comments, unused dependencies, and invented configuration.
- [ ] **USEQ-0A35B7B4** — Check for fabricated APIs, obsolete methods, invalid options, deprecated protocols, unsupported versions, and nonexistent packages.
- [ ] **USEQ-80CB3947** — Check for injection, unsafe deserialization, path traversal, authorization bypass, secret exposure, insecure randomness, weak cryptography, and dangerous defaults.
- [ ] **USEQ-069C7D6E** — Check for integer, numeric, time, locale, Unicode, encoding, null, overflow, and precision errors.
- [ ] **USEQ-90FD5CD0** — Check for resource leaks, unbounded work, missing timeouts, unsafe retries, race conditions, deadlocks, cancellation defects, and duplicate effects.
- [ ] **USEQ-FDF5A263** — Check for missing observability, misleading logs, swallowed errors, broad exception handling, and loss of diagnostic context.
- [ ] **USEQ-E92C614F** — Check for accidental cross-tenant access, global mutable state, cache-key omissions, and authorization context loss.
- [ ] **USEQ-2BB3F020** — Check whether generated code duplicates an existing internal capability or violates established architecture decisions.
- [ ] **USEQ-8BBA5659** — Prefer small changes that fit the current design over broad rewrites generated without full system context.
- [ ] **USEQ-471DDE50** — Require generated public interfaces and persistent schemas to undergo normal compatibility and evolution review.
- [ ] **USEQ-FF208505** — Run formatting, type checking, static analysis, dependency analysis, tests, and relevant security checks on the exact accepted output.
- [ ] **USEQ-592A4BA4** — Retain regression tests for defects found in generated code and update organizational guidance to prevent recurrence.

### AI-generated tests, fixtures, and verification artifacts

- [ ] **USEQ-239E5251** — Review whether generated tests assert intended requirements or merely reproduce current implementation behavior.
- [ ] **USEQ-BF92D3B4** — Prevent the same model context from generating code and expected outputs without independent oracle validation.
- [ ] **USEQ-7BAE92AE** — Verify that tests fail when the intended behavior is deliberately broken.
- [ ] **USEQ-36C03659** — Use mutation or controlled fault insertion to evaluate assertion strength for critical generated tests.
- [ ] **USEQ-524C0709** — Ensure generated tests cover negative, boundary, error, recovery, authorization, concurrency, and data-integrity behavior rather than only happy paths.
- [ ] **USEQ-12534F05** — Review generated mocks and stubs for incorrect assumptions about real dependencies.
- [ ] **USEQ-F374893D** — Ensure generated fixtures do not contain real personal data, credentials, proprietary examples, or unsafe production identifiers.
- [ ] **USEQ-ED97CA2A** — Control randomness and retain seeds for generated property, fuzz, and scenario tests.
- [ ] **USEQ-44E65FA4** — Remove redundant or brittle tests that create maintenance cost without meaningful assurance.
- [ ] **USEQ-44661241** — Verify that snapshots and golden files represent approved outcomes and do not normalize accidental output.
- [ ] **USEQ-BAE2D2AE** — Review test names, diagnostics, setup, cleanup, isolation, and failure messages for maintainability.
- [ ] **USEQ-9142AF7A** — Ensure generated test data exercises representative Unicode, locale, time, size, tenant, role, and lifecycle cases.
- [ ] **USEQ-16875252** — Run generated tests against the actual artifact and production-like configuration where relevant.
- [ ] **USEQ-46A93573** — Track defects escaped despite generated tests and improve prompts, review, or techniques based on root cause.
- [ ] **USEQ-F7978D60** — Do not count generated test volume as quality evidence without coverage, assertion, and fault-detection analysis.

### AI-generated requirements, designs, architecture, and decisions

- [ ] **USEQ-6DD50F52** — Validate generated requirements through user research, stakeholder agreement, domain evidence, legal review, and operational feasibility.
- [ ] **USEQ-DD4E49F4** — Reject invented user needs, regulations, standards, constraints, metrics, competitors, incidents, or customer commitments.
- [ ] **USEQ-B4A6F52A** — Identify assumptions and uncertainty explicitly rather than converting them into authoritative prose.
- [ ] **USEQ-3E4A6B5A** — Verify that generated acceptance criteria are measurable, complete, noncontradictory, and technology-neutral where intended.
- [ ] **USEQ-3BCFBFCB** — Review generated architecture against actual scale, data, trust boundaries, failure domains, existing systems, skills, cost, and lifecycle constraints.
- [ ] **USEQ-1CD27635** — Require alternatives, tradeoffs, consequences, and reversal strategy in material generated design decisions.
- [ ] **USEQ-B9E44CC3** — Avoid adopting fashionable patterns, services, abstractions, or distributed architectures without demonstrated need.
- [ ] **USEQ-7BA2175A** — Verify diagrams, flows, schemas, sequence descriptions, and dependency maps against the deployed or intended system.
- [ ] **USEQ-4344CC44** — Threat-model generated architectures independently and check for hidden trust, control-plane, tenant, and data-flow assumptions.
- [ ] **USEQ-B77FC1DD** — Review generated data models for integrity, retention, privacy, migration, localization, and reporting requirements.
- [ ] **USEQ-F3633846** — Review generated performance and capacity estimates against measurements and defensible workload models.
- [ ] **USEQ-D7408E34** — Ensure generated decisions preserve accessibility, supportability, observability, rollback, recovery, and decommissioning.
- [ ] **USEQ-5E4618A5** — Record AI assistance and source evidence in material architecture decision records where relevant to future review.
- [ ] **USEQ-E7FE1B6C** — Do not let generated documentation substitute for accountable stakeholder decisions.
- [ ] **USEQ-AE086A58** — Revalidate generated designs after requirements, providers, standards, risks, or constraints change.

### AI-generated infrastructure, configuration, queries, and migrations

- [ ] **USEQ-24BBED4A** — Treat generated infrastructure, deployment, policy, database, query, migration, and operational code as high-impact executable artifacts.
- [ ] **USEQ-802762AD** — Run schema validation, static analysis, policy checks, security analysis, plan or dry-run review, and peer review before execution.
- [ ] **USEQ-764B2420** — Verify resource names, accounts, regions, networks, identities, permissions, quotas, retention, encryption, and deletion protection.
- [ ] **USEQ-ECEF8251** — Check that generated permissions are least-privilege and do not use broad wildcards for convenience.
- [ ] **USEQ-43965569** — Check that generated network rules do not expose administrative, database, storage, metadata, or internal services unnecessarily.
- [ ] **USEQ-A6D8DA54** — Check that generated configuration does not enable debug mode, insecure defaults, public access, unbounded resources, or disabled verification.
- [ ] **USEQ-B9138B79** — Evaluate generated database queries for correctness, authorization, injection, locking, resource cost, indexing, consistency, and data minimization.
- [ ] **USEQ-F11E331B** — Rehearse generated migrations against representative data and validate compatibility, duration, locking, restart, rollback, roll-forward, and integrity checks.
- [ ] **USEQ-F7C0B574** — Require previews and explicit confirmation before destructive or irreversible actions.
- [ ] **USEQ-D20D5B1C** — Prevent agents from executing production commands with credentials inherited from a developer session without controlled elevation.
- [ ] **USEQ-6A91D81E** — Use isolated sandboxes and nonproduction environments for generated command experimentation.
- [ ] **USEQ-53136269** — Record exact generated commands, plans, affected resources, approvals, execution output, and resulting state.
- [ ] **USEQ-21445766** — Verify idempotency and safe rerun behavior for automation that may be retried.
- [ ] **USEQ-386CA87E** — Ensure generated cleanup does not delete shared, retained, legally held, or customer-owned resources.
- [ ] **USEQ-7D1589E9** — Maintain a tested rollback or recovery path before high-impact generated changes are applied.

### Dependencies, packages, licenses, and generated supply-chain choices

- [ ] **USEQ-A42E53F8** — Verify every suggested package, image, action, plugin, model, data set, API, and service exists and comes from the intended publisher.
- [ ] **USEQ-ED0BD0F8** — Check spelling, namespace, registry, repository, ownership, signatures, provenance, release history, and known compromise before adoption.
- [ ] **USEQ-D2946F6A** — Prefer established existing dependencies already approved for the system when they satisfy the requirement.
- [ ] **USEQ-E4E7B8C5** — Reject unnecessary dependencies introduced to solve trivial problems or because they appear in generic training examples.
- [ ] **USEQ-EC5B4CD4** — Verify supported versions, maintenance health, security process, license, data use, portability, and replacement feasibility.
- [ ] **USEQ-69099CAD** — Pin or constrain dependencies according to supply-chain policy and generate updated SBOM and provenance evidence.
- [ ] **USEQ-8BF5FC3E** — Review copied or generated license headers, attribution, notices, and source obligations for accuracy.
- [ ] **USEQ-6C4E1F4A** — Do not assume generated code is original or compatible with the project license.
- [ ] **USEQ-0BE33354** — Search for suspiciously distinctive fragments or notices when generated output resembles known external code.
- [ ] **USEQ-36A0CE74** — Prevent AI from choosing packages or services solely through popularity or unverifiable recommendation.
- [ ] **USEQ-0D1E4378** — Require architecture review before generated changes introduce a new database, queue, framework, cloud service, model provider, or operational platform.
- [ ] **USEQ-6C07A1CF** — Monitor AI-suggested dependencies for advisories, ownership change, deprecation, and malicious releases after adoption.
- [ ] **USEQ-31D5D8DE** — Record why each new dependency is necessary and who owns its lifecycle.
- [ ] **USEQ-0A9114C1** — Remove hallucinated, unused, duplicate, obsolete, or insecure dependencies before merge.
- [ ] **USEQ-4C7730ED** — Include models, prompts, adapters, embeddings, data sets, and agent tools in relevant supply-chain inventories.

### Agentic tools, autonomous execution, and capability containment

- [ ] **USEQ-E1582655** — Define each agent capability explicitly, including files, repositories, commands, network destinations, APIs, credentials, data, and environments it may access.
- [ ] **USEQ-EE41DE02** — Deny capabilities by default and grant the minimum scope required for the current task.
- [ ] **USEQ-DBCE39F6** — Use separate identities and sandboxes for agents rather than sharing a user or production administrator session.
- [ ] **USEQ-C5639E97** — Keep untrusted prompts, issue content, code comments, retrieved documents, and web pages from expanding agent permissions.
- [ ] **USEQ-A1439915** — Require human confirmation at clearly identified boundaries for destructive, external, privileged, financial, production, publication, and communication actions.
- [ ] **USEQ-87469F66** — Show the exact proposed action, target, scope, side effects, credentials, and rollback implications before confirmation.
- [ ] **USEQ-D22830A2** — Prevent an agent from modifying its own policy, approval rules, audit records, safety controls, or credentials.
- [ ] **USEQ-15C830F0** — Bound runtime, token use, cost, network requests, file changes, command count, concurrency, recursion, retries, and spawned agents.
- [ ] **USEQ-71EE868E** — Restrict outbound network destinations and validate downloaded content and tools.
- [ ] **USEQ-A348D57E** — Prevent agents from exfiltrating source, secrets, personal data, or private context through prompts, URLs, logs, errors, or generated artifacts.
- [ ] **USEQ-E07F70B6** — Use allowlisted commands or isolated execution for high-risk environments.
- [ ] **USEQ-EA332D5C** — Record tool calls, inputs, outputs, file changes, decisions, approvals, and errors in an auditable form that protects sensitive data.
- [ ] **USEQ-172D4874** — Provide immediate pause, cancellation, credential revocation, and kill-switch capabilities.
- [ ] **USEQ-26B3A07B** — Recover safely from partial execution and identify every side effect already performed.
- [ ] **USEQ-D31234E4** — Test prompt injection, confused-deputy behavior, privilege escalation, hidden instructions, malicious tool output, and infinite-loop scenarios.
- [ ] **USEQ-B55613A3** — Do not permit unsupervised production changes until the agent has demonstrated reliable behavior under realistic adversarial evaluation and the residual risk is authorized.

### Model, prompt, adapter, and AI-component supply-chain assurance

- [ ] **USEQ-EDAB9D12** — Inventory base models, fine-tunes, adapters, prompts, system instructions, safety policies, retrieval indexes, embeddings, data sets, evaluators, and external tools used in each AI-assisted workflow.
- [ ] **USEQ-9832AA48** — Record source, provider, version, hash or immutable identifier, license, permitted use, training or fine-tuning data claims, and support status where available.
- [ ] **USEQ-68821BC1** — Bind deployed AI configuration to a release record so behavior changes can be traced.
- [ ] **USEQ-21EAD928** — Verify downloaded model and adapter integrity and retrieve them through authenticated trusted channels.
- [ ] **USEQ-9509D8FC** — Assess model, prompt, adapter, and data poisoning risks before adoption.
- [ ] **USEQ-C41FAB76** — Review provider and community models for unsafe serialization, executable code, custom loaders, hidden tools, and excessive permissions.
- [ ] **USEQ-9F95E96A** — Do not execute untrusted model repositories, notebooks, setup scripts, or custom code with production credentials.
- [ ] **USEQ-90F32C34** — Scan model and data packages for secrets, personal data, malicious files, license obligations, and unexpected components.
- [ ] **USEQ-23C5B858** — Separate evaluation, development, and production credentials and data.
- [ ] **USEQ-C9B74797** — Control who may change prompts, retrieval sources, model routing, tool definitions, safety settings, and thresholds.
- [ ] **USEQ-886D57EB** — Review and test changes even when the provider labels them as minor or safety-related.
- [ ] **USEQ-6B2D33C4** — Maintain a rollback path to a known evaluated model and configuration.
- [ ] **USEQ-57129030** — Monitor provider outages, model withdrawals, silent updates, policy changes, and compromised distribution channels.
- [ ] **USEQ-F05A29D5** — Include AI components in SBOM-like, model-card, system-card, or provenance records appropriate to the system.
- [ ] **USEQ-B016EE0A** — Define retirement procedures for models, embeddings, prompts, credentials, retained data, and downstream derivatives.

### Evaluation, red teaming, regression, and calibrated trust

- [ ] **USEQ-F1783FEA** — Define evaluations for the actual engineering tasks, repositories, languages, domains, users, tools, and risk conditions rather than relying only on general benchmarks.
- [ ] **USEQ-EB79B3ED** — Measure functional correctness, security, privacy, licensing, maintainability, performance, accessibility, documentation quality, and operational safety as applicable.
- [ ] **USEQ-0C102A58** — Include incomplete context, conflicting requirements, obsolete documentation, malicious instructions, unusual code, large changes, and unfamiliar domains.
- [ ] **USEQ-8235AA68** — Evaluate refusal and escalation behavior when the tool lacks evidence, permission, capability, or confidence.
- [ ] **USEQ-0F3DD770** — Evaluate whether the tool fabricates APIs, packages, citations, standards, test results, vulnerabilities, or completed actions.
- [ ] **USEQ-D21E0A71** — Test direct and indirect prompt injection, data exfiltration, tool abuse, confused deputy, secret discovery, sandbox escape, and policy bypass.
- [ ] **USEQ-9994DB54** — Test output under model updates, prompt changes, retrieval changes, tool changes, provider routing changes, and degraded dependencies.
- [ ] **USEQ-F4C1D548** — Use independently curated evaluation sets and protect them from contamination where scores are used for decisions.
- [ ] **USEQ-B854BC2B** — Retain representative historical failures and incidents as regression cases.
- [ ] **USEQ-7633CD04** — Use human expert review for semantic, architectural, security, legal, and safety qualities that automated scoring cannot establish.
- [ ] **USEQ-08934850** — Measure uncertainty and confidence calibration where the system exposes confidence or autonomous thresholds.
- [ ] **USEQ-53A8E42B** — Define acceptable performance by task and harm, not one aggregate score.
- [ ] **USEQ-48A484CF** — Investigate subgroup, language, accessibility, repository, and domain differences rather than averaging them away.
- [ ] **USEQ-B0B7E4B2** — Require evaluation before expanding permissions, autonomy, data access, user population, or production scope.
- [ ] **USEQ-EEB12085** — Publish internal limitations and approved use boundaries alongside evaluation results.
- [ ] **USEQ-3C2B0FE1** — Reevaluate periodically and after incidents, model changes, new attacks, and meaningful shifts in production tasks.

### AI-generated documentation, citations, and knowledge integrity

- [ ] **USEQ-48B9998C** — Verify technical facts, commands, configuration, API behavior, version requirements, standards status, dates, and external links against authoritative sources.
- [ ] **USEQ-C1C0BF4C** — Do not include fabricated citations, quotations, people, products, vulnerabilities, decisions, metrics, or incident details.
- [ ] **USEQ-C5181864** — Distinguish facts, inference, recommendation, example, placeholder, and uncertainty clearly.
- [ ] **USEQ-A3559530** — Ensure generated documentation matches the exact released artifact and supported configuration.
- [ ] **USEQ-7B8E61BE** — Review setup, migration, rollback, recovery, security, accessibility, and troubleshooting instructions through actual execution.
- [ ] **USEQ-2D5D5C89** — Protect internal, personal, confidential, and security-sensitive information during generation and publication.
- [ ] **USEQ-75862A50** — Ensure generated user content is understandable, accessible, localized appropriately, and free of deceptive claims.
- [ ] **USEQ-F06F2315** — Review generated release notes so they do not omit breaking changes, security impact, required customer action, or known limitations.
- [ ] **USEQ-44ADFB43** — Use source attribution and content provenance where readers need to assess authority and freshness.
- [ ] **USEQ-FAC8F384** — Version generated documentation and invalidate it after material product or provider changes.
- [ ] **USEQ-327CFE20** — Prevent AI summaries from replacing authoritative records when detail or legal meaning matters.
- [ ] **USEQ-E547539A** — Ensure search and knowledge systems surface the authoritative current source rather than an unreviewed generated duplicate.
- [ ] **USEQ-EE1E3012** — Provide a correction route and update propagated copies after material errors are found.
- [ ] **USEQ-E1EEAF8B** — Label synthetic or AI-generated content when required by law, contract, policy, audience expectation, or risk.
- [ ] **USEQ-80D53A5E** — Retain the reviewed final artifact rather than depending on transient conversation history.

### Content provenance, authenticity, and C2PA-style credentials

- [ ] **USEQ-798E6E1D** — Determine which images, audio, video, documents, designs, reports, releases, and public communications require provenance or authenticity information.
- [ ] **USEQ-5D158B5E** — Define what provenance claims mean and avoid presenting them as proof that the content is true, harmless, complete, or unbiased.
- [ ] **USEQ-EDD07EC1** — Use a recognized provenance standard such as the current C2PA specification when interoperable content credentials are required.
- [ ] **USEQ-D9985041** — Cryptographically bind provenance manifests to the exact asset or defined derivative relationship.
- [ ] **USEQ-22BC5813** — Protect signing identities and keys and define who is authorized to make each assertion.
- [ ] **USEQ-EE822CDB** — Record creation, editing, transformation, generation, source, ingredient, tool, organization, and timestamp information only when accurate and permitted.
- [ ] **USEQ-E59769F2** — Minimize personal, precise-location, device, and confidential metadata in provenance records.
- [ ] **USEQ-C4489851** — Provide user interfaces that explain provenance in understandable layers rather than only displaying a trust icon.
- [ ] **USEQ-841334CB** — Show missing, invalid, broken, revoked, or untrusted provenance without claiming the underlying content is necessarily false.
- [ ] **USEQ-CF367760** — Preserve provenance through supported transformations, exports, compression, transcoding, publication, and distribution where feasible.
- [ ] **USEQ-2A380FCA** — Detect and disclose when a workflow strips or cannot preserve provenance.
- [ ] **USEQ-5A50E71F** — Validate trust chains, certificates, timestamps, revocation, content bindings, and manifest structure.
- [ ] **USEQ-941EA8C0** — Protect against manifest substitution, replay, downgrade, detached-asset confusion, and misleading partial histories.
- [ ] **USEQ-0C275FAF** — Define how edited, combined, redacted, and derived content updates provenance.
- [ ] **USEQ-D078F7CC** — Test provenance handling across supported applications, platforms, media formats, and accessibility technologies.
- [ ] **USEQ-29786656** — Maintain a recovery and key-rotation process that does not invalidate legitimate historical content without explanation.

### Verifiable claims, credentials, signatures, and evidence integrity

- [ ] **USEQ-2325DA4F** — Use verifiable credentials or signed evidence only when cryptographic verification, issuer accountability, portability, or selective disclosure provides a real requirement.
- [ ] **USEQ-EE4C0408** — Define issuer, holder, verifier, subject, claims, purpose, audience, validity, revocation, status, and trust policy explicitly.
- [ ] **USEQ-F7C75D9C** — Use current interoperable standards such as the W3C Verifiable Credentials Data Model 2.0 family when applicable.
- [ ] **USEQ-5FE49B3A** — Verify issuer identity, key authorization, proof purpose, domain or challenge binding, validity period, status, schema, and audience before accepting a claim.
- [ ] **USEQ-6BF7123C** — Do not treat a valid signature as proof that the claim is true, current, authorized for the present purpose, or ethically appropriate.
- [ ] **USEQ-B791AEA6** — Request and disclose only the claims necessary for the transaction.
- [ ] **USEQ-93F664BC** — Support selective disclosure and unlinkability where the use case and chosen standard permit it.
- [ ] **USEQ-A80DCC8B** — Prevent identifiers and status checks from enabling unnecessary correlation or tracking.
- [ ] **USEQ-403D0FC0** — Protect credential issuance, storage, presentation, recovery, revocation, and deletion.
- [ ] **USEQ-AABF5B36** — Handle lost, compromised, expired, suspended, and superseded credentials.
- [ ] **USEQ-CBD654F9** — Validate canonicalization, serialization, proof suites, cryptographic parameters, and implementation interoperability.
- [ ] **USEQ-F2B1299D** — Protect against replay, presentation substitution, confused-deputy, verifier impersonation, and malicious issuer metadata.
- [ ] **USEQ-AE60A435** — Provide understandable user consent and display the verifier, requested claims, purpose, and consequences before presentation.
- [ ] **USEQ-BF615392** — Maintain audit evidence without retaining unnecessary credential contents.
- [ ] **USEQ-00A91D07** — Test offline, degraded, revoked, clock-skewed, and trust-store-update scenarios.
- [ ] **USEQ-6EBF4486** — Plan crypto-agile and post-quantum migration for long-lived credentials and evidence.

### Synthetic media, generated content, and public trust

- [ ] **USEQ-FB59A946** — Define when generated or materially altered content must be labeled, watermarked, credentialed, disclosed, or restricted.
- [ ] **USEQ-24A27761** — Make disclosures visible, accessible, durable, and understandable in the context where content is consumed.
- [ ] **USEQ-6D2ED17A** — Avoid disclosures that disappear when content is downloaded, embedded, cropped, translated, printed, or shared.
- [ ] **USEQ-87A2D323** — Prevent generated content from impersonating real people, organizations, authorities, or official communications without authorization and appropriate safeguards.
- [ ] **USEQ-DFB7F6ED** — Obtain necessary consent and rights for likeness, voice, biometric, personal, copyrighted, and confidential material.
- [ ] **USEQ-9CD66A9B** — Protect children, vulnerable users, and high-risk contexts with stricter generation and distribution controls.
- [ ] **USEQ-31086253** — Review generated public claims for factual accuracy, defamation, discrimination, safety, financial, medical, legal, and civic risks.
- [ ] **USEQ-9D635145** — Use human review before publishing high-impact generated content under an organizational identity.
- [ ] **USEQ-8C754D48** — Do not present detection scores as definitive proof that content is or is not AI-generated.
- [ ] **USEQ-5710654C** — Combine provenance, source verification, context, policy, and investigation rather than relying on one synthetic-content detector.
- [ ] **USEQ-C731071C** — Provide correction, takedown, appeal, and incident-reporting routes.
- [ ] **USEQ-02B7A344** — Record the approved final content, sources, model and tool versions, material transformations, reviewer, and publication channels.
- [ ] **USEQ-10EDBDFD** — Monitor for misuse of organizational names, marks, credentials, voices, and identities in synthetic media.
- [ ] **USEQ-D3FD11E3** — Prepare response procedures for forged, misleading, compromised, or incorrectly credentialed content.
- [ ] **USEQ-97B7BEBB** — Preserve evidence needed to investigate provenance while respecting privacy and legal constraints.

### Monitoring, incident response, rollback, and retirement

- [ ] **USEQ-0CD1F465** — Monitor AI-assisted workflows for unauthorized data access, secret exposure, insecure output, fabricated facts, malicious dependencies, excessive permissions, abnormal tool use, cost spikes, and repeated human overrides.
- [ ] **USEQ-E44BDA17** — Monitor model, provider, prompt, retrieval, permission, policy, and integration changes that can alter behavior.
- [ ] **USEQ-91F85F40** — Detect whether agents perform actions outside the approved task, scope, identity, repository, environment, or time window.
- [ ] **USEQ-D20D2DD1** — Provide rapid disablement for each model, provider, integration, tool, retrieval source, agent, and autonomous capability.
- [ ] **USEQ-E0D7CFC5** — Revoke associated credentials, webhooks, sessions, scheduled tasks, and network access when disabled.
- [ ] **USEQ-39F46163** — Maintain fallback workflows that do not depend on the unavailable or compromised AI system for essential operation.
- [ ] **USEQ-A737503A** — Define incident playbooks for confidential-data leakage, malicious action, unsafe code, poisoned context, provider compromise, model change, provenance-key compromise, and synthetic impersonation.
- [ ] **USEQ-4EE501B5** — Identify every repository, artifact, deployment, document, decision, and communication affected by a compromised or defective AI workflow.
- [ ] **USEQ-B7A152C3** — Preserve prompts, context, outputs, tool calls, approvals, model versions, and audit evidence proportionately and lawfully for investigation.
- [ ] **USEQ-FA7DAC9D** — Notify affected customers, users, employees, providers, and authorities according to impact and obligations.
- [ ] **USEQ-9CA9E828** — Search for related generated defects and artifacts after an incident rather than fixing only the first instance.
- [ ] **USEQ-0A7A2F6C** — Revalidate accepted output after a material model or tool defect is discovered.
- [ ] **USEQ-2C7F55DC** — Perform root-cause analysis covering governance, permissions, interface design, review, evaluation, training, incentives, and technical controls.
- [ ] **USEQ-3EEAD754** — Retire obsolete models, prompts, embeddings, agents, indexes, credentials, provider data, and generated artifacts according to lifecycle and retention requirements.
- [ ] **USEQ-86EF8E1F** — Test disablement, rollback, provider replacement, data deletion, and evidence recovery periodically.

### Release blockers and AI-assistance evidence package

- [ ] **USEQ-4858F7E6** — Block release when material AI-generated code, configuration, migration, policy, or architecture has no accountable human reviewer who understands it.
- [ ] **USEQ-3D7EFEF0** — Block release when generated claims, citations, APIs, packages, standards, or compatibility assumptions have not been verified.
- [ ] **USEQ-6465C673** — Block release when AI assistance exposed secrets, personal data, customer content, vulnerabilities, or proprietary material and containment is incomplete.
- [ ] **USEQ-E5E6E4D2** — Block release when an agent executed high-impact actions outside the approved capability, identity, environment, or confirmation boundary.
- [ ] **USEQ-FD59DB11** — Block release when critical generated tests use the same unverified logic as the implementation or have not demonstrated fault detection.
- [ ] **USEQ-C2B2D033** — Block release when model, prompt, retrieval, tool, or provider changes invalidate prior assurance evidence.
- [ ] **USEQ-BA8ED02F** — Block release when generated dependencies, licenses, provenance, or supplier identities remain unresolved.
- [ ] **USEQ-14FDF5E4** — Block release when disabling or rolling back the AI capability is not possible for a high-risk workflow.
- [ ] **USEQ-41AEF92F** — Retain the use case, approved tool and model versions, prompts or prompt templates where appropriate, retrieved sources, permissions, tool calls, generated artifacts, human changes, reviews, tests, and approvals.
- [ ] **USEQ-3CFE6905** — Retain evaluation scope, data, adversarial cases, results, limitations, model changes, incidents, and accepted residual risk.
- [ ] **USEQ-85D0CB75** — Record data classifications, provider processing terms, retention settings, deletion evidence, and intellectual-property review where material.
- [ ] **USEQ-A3940500** — Record provenance and authenticity evidence for public content and signed claims where used.
- [ ] **USEQ-DC177F9F** — Require independent sign-off for autonomous production action, safety-related output, security controls, identity decisions, financial movement, or high-impact public communications.
- [ ] **USEQ-54B598C4** — Ensure ongoing monitoring detects behavior drift, unauthorized use, provider changes, data leakage, and control degradation after release.
- [ ] **USEQ-AEE1EB03** — State the exact scope of AI assistance and residual uncertainty rather than claiming the work was fully verified merely because automated checks passed.

## Standards and source references

- [ISO/IEC 42001:2023 — AI management systems](https://www.iso.org/standard/81230.html)
- [ISO/IEC 23894:2023 — AI risk management](https://www.iso.org/standard/77304.html)
- [NIST AI Risk Management Framework 1.0](https://www.nist.gov/itl/ai-risk-management-framework)
- [NIST AI RMF Generative AI Profile](https://www.nist.gov/publications/artificial-intelligence-risk-management-framework-generative-artificial-intelligence)
- [NIST SP 800-218A — Secure Software Development Practices for Generative AI and Dual-Use Foundation Models](https://csrc.nist.gov/pubs/sp/800/218/a/final)
- [ISO/IEC 25059:2023 — Quality model for AI systems](https://www.iso.org/standard/80655.html)
- [ISO/IEC 42005:2025 — AI system impact assessment](https://www.iso.org/standard/44545.html)
- [ISO/IEC 5338:2023 — AI system life-cycle processes](https://www.iso.org/standard/81118.html)
- [OWASP Artificial Intelligence Security Verification Standard](https://owasp.org/www-project-artificial-intelligence-security-verification-standard/)
- [OWASP AI Testing Guide](https://owasp.org/www-project-ai-testing-guide/)
- [DORA 2025 State of AI-assisted Software Development](https://dora.dev/research/2025/dora-report/)

---

[Previous phase](14-trust-safety-and-ecosystems.md) · [Next: Phase 16: Specialized domains and release assurance](16-specialized-domains-and-release-assurance.md)
