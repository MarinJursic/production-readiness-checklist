---
title: Production Readiness Checklist
description: 1,421 evidence-driven checks for shipping secure, reliable, and supportable web applications with confidence.
hide:
  - toc
---

<!-- markdownlint-disable MD025 -->

<div class="prc-hero" markdown>

<div class="prc-hero__copy" markdown>

<span class="prc-eyebrow">OPEN-SOURCE RELEASE READINESS SYSTEM</span>

# Ship with evidence, not optimism

Turn production readiness into a reviewable decision with **1,421 technology-neutral controls**, stable evidence IDs, release templates, and honest AI-assisted audits.

<div class="prc-hero__actions">
  <a class="md-button md-button--primary" href="checklists/01-release-foundations/">Start the checklist <span aria-hidden="true">→</span></a>
  <a class="md-button prc-star-button" href="https://github.com/MarinJursic/production-readiness-checklist" target="_blank" rel="noopener noreferrer"><span aria-hidden="true">★</span> Star on GitHub</a>
  <span class="prc-share-group">
    <button class="md-button prc-share-button" type="button" data-prc-share hidden>Share this project</button>
    <span class="prc-share-status" role="status" aria-live="polite"></span>
  </span>
</div>

</div>

<div class="prc-hero__visual" aria-label="Example production-readiness evidence cards">
  <div class="prc-signal-card prc-signal-card--primary">
    <span class="prc-signal-card__id">PRC-34-017</span>
    <strong>Rollback rehearsed</strong>
    <span class="prc-signal-card__status"><span aria-hidden="true">●</span> Evidence ready</span>
  </div>
  <div class="prc-signal-card">
    <span class="prc-signal-card__id">PRC-02-006</span>
    <strong>Restore demonstrated</strong>
    <span class="prc-signal-card__status"><span aria-hidden="true">●</span> Independently reviewed</span>
  </div>
  <div class="prc-signal-card">
    <span class="prc-signal-card__id">PRC-31-001</span>
    <strong>Security validated</strong>
    <span class="prc-signal-card__status"><span aria-hidden="true">●</span> Current release</span>
  </div>
</div>

</div>

<!-- markdownlint-enable MD025 -->

<div class="prc-stats" aria-label="Checklist overview">
  <div><strong>1,421</strong><span>evidence-driven controls</span></div>
  <div><strong>43</strong><span>numbered sections</span></div>
  <div><strong>10</strong><span>focused review tracks</span></div>
  <div><strong>4</strong><span>honest evidence states</span></div>
</div>

## Choose your review path

<!-- markdownlint-disable MD030 -->

<div class="grid cards prc-paths" markdown>

-   <span class="prc-card-number">01</span> **I am preparing a release**

    ---

    Run the no-go screen, select applicable tracks, collect evidence, and record the accountable decision.

    [Open the 15-minute quick start →](guides/getting-started.md)

-   <span class="prc-card-number">02</span> **I want an AI review**

    ---

    Give Claude or another coding agent strict evidence rules so it reports unknowns instead of inventing passes.

    [Use the AI-assisted workflow →](guides/ai-assisted-review.md)

-   <span class="prc-card-number">03</span> **I am adopting this for a team**

    ---

    Assign track owners, define evidence freshness, tailor applicability, and keep risk acceptance visible.

    [Copy the release assessment →](records/release-assessment.md)

</div>

<!-- markdownlint-enable MD030 -->

## From release candidate to accountable decision

<div class="prc-steps">
  <div><span>1</span><strong>Identify</strong><p>Pin the exact artifact, configuration, migrations, flags, environment, and scope.</p></div>
  <div><span>2</span><strong>Stop early</strong><p>Screen all immediate no-go conditions before investing in the full assessment.</p></div>
  <div><span>3</span><strong>Prove</strong><p>Attach current evidence and an accountable owner to every applicable control.</p></div>
  <div><span>4</span><strong>Decide</strong><p>Sign off, deploy progressively, verify production, and retain the decision record.</p></div>
</div>

## Explore the checklist

<div class="prc-track-grid">
  <a href="checklists/01-release-foundations/"><span>01</span><strong>Release foundations</strong><small>No-go gates, identity, and scope</small></a>
  <a href="checklists/02-product-risk-architecture/"><span>02</span><strong>Product, risk, and architecture</strong><small>Intent, ownership, and threat modeling</small></a>
  <a href="checklists/03-source-build-supply-chain/"><span>03</span><strong>Source, build, and supply chain</strong><small>CI/CD, provenance, SBOMs, and licenses</small></a>
  <a href="checklists/04-environments-quality-experience/"><span>04</span><strong>Quality and experience</strong><small>Configuration, testing, frontend, and accessibility</small></a>
  <a href="checklists/05-application-security/"><span>05</span><strong>Application security</strong><small>Identity, authorization, input, transport, and crypto</small></a>
  <a href="checklists/06-data-privacy-performance/"><span>06</span><strong>Data, privacy, and performance</strong><small>Integrity, migrations, privacy, capacity, and overload</small></a>
  <a href="checklists/07-reliability-operations/"><span>07</span><strong>Reliability and operations</strong><small>Resilience, observability, recovery, and deployment</small></a>
  <a href="checklists/08-maintenance-vendors-compliance/"><span>08</span><strong>Vendors and compliance</strong><small>Operability, third parties, legal, and regulatory modules</small></a>
  <a href="checklists/09-conditional-modules/"><span>09</span><strong>Conditional modules</strong><small>Payments, SaaS, AI, real-time, UGC, and more</small></a>
  <a href="checklists/10-evidence-and-decision/"><span>10</span><strong>Evidence and decision</strong><small>Evidence package, sign-offs, and go/no-go rules</small></a>
</div>

## Evidence, not checkbox theater

<div class="prc-principle" markdown>

> Every applicable requirement has current evidence; no known risk exceeds the organization’s tolerance; critical user journeys meet defined objectives; and the organization can detect, contain, roll back, restore, support, and communicate failures.

A checked box without evidence is only a claim. Use **Pass**, **Fail**, **Blocked**, or **Not Applicable**—and never hide a release blocker inside an average score.

[Read the readiness principle](checklists/00-readiness-principle.md) · [Review references and limitations](references.md)

</div>

<div class="prc-community-panel" markdown>

### Help more teams ship responsibly

If this checklist saves your team from one avoidable incident, help another team find it.

<a class="md-button md-button--primary" href="https://github.com/MarinJursic/production-readiness-checklist" target="_blank" rel="noopener noreferrer"><span aria-hidden="true">★</span> Star the repository</a>
<span class="prc-share-group">
  <button class="md-button prc-share-button" type="button" data-prc-share hidden>Share this project</button>
  <span class="prc-share-status" role="status" aria-live="polite"></span>
</span>

</div>
