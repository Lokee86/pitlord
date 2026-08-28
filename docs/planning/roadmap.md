# Roadmap

Parent index: [Pitlord Planning](INDEX.md)

## Purpose

This document records planned Pitlord work and unresolved sequencing.

## Overview

Pitlord remains independently usable while expanding from architecture-policy enforcement into a deterministic, opinionated architecture and complexity governor for coding agents. Its default generalized guard should enforce broadly useful clean-code and clean-architecture constraints without requiring repository-specific policy, while repository policy continues to encode local intent.

The central product principle is deterministic supervision: coding agents may be probabilistic, but Pitlord's observation, judgment, evidence, and verification must be mechanically reproducible. Pitlord should prevent complexity accumulation without introducing another autonomous model-review loop around the agent.

## Current status

Implementation is underway. The policy-free `scan` command and deterministic `pitlord.scan.v1` finding contract are in place. The language-neutral `dependency-pressure` detector has completed its first 16-corpus calibration pass. The `dependency-knots` detector has also completed tuning against its frozen 16-corpus detector-specific reference set: the current projection preserves all 7 required knots with 0 labelled false positives, 0 false negatives, 0 severity mismatches, and 0 unlabelled findings. Boundary/cohesion analysis and blast-radius analysis remain next.

## Expected ownership

- Pitlord owns anomaly interpretation, diagnosis identity, accepted architectural decisions, policy authoring semantics, rule evaluation, evidence identity, baselines, and machine-readable diagnostics.
- Arcana owns graph storage, traversal, communities, impact, paths, strongly connected components, and other structural measurements.
- Lexicon owns language facts.
- Demon Docs and the shared checker own Markdown structure, navigation, and documentation coverage checks.
- Warlock owns scheduling, enrolled-repository orchestration, and interactive diagnosis/policy presentation.

## Planned work

1. Add a policy-free generalized architecture scan over Arcana structural evidence, beginning with dependency knots, hub/bottleneck anomalies, boundary-coupling/cohesion anomalies, and impact hotspots.
2. Classify generalized detectors explicitly as advisory or blocking guard conditions and make all blocking findings bounded, mechanically verifiable repair targets for coding agents.
3. Add an immediate post-change guard workflow: agent changes code, Pitlord evaluates deterministically, the agent repairs the finite failures, and Pitlord verifies them without any model-based adjudication.
4. Guarantee reproducible findings for the same snapshot, Pitlord version, configuration, and generalized profile, including finding identity, ordering, severity, evidence, and exit status.
5. Add explicit human disposition for ambiguous advisory diagnoses: create a repository guardrail, accept the structure as intentional, or leave it unresolved.
6. Add user-friendly policy authoring that can promote a diagnosis into validated policy without requiring JSON-schema knowledge, including a current-repository match preview before publication.
7. Keep accepted architectural decisions semantically distinct from evidence baselines and define a durable repository-owned persistence contract for them.
8. Extend the shipped Arcana-topology-family calibration to each generalized detector, then tune against representative real repositories without weakening the synthetic separation cases.
9. Add reusable policy packs with explicit versioning and compatibility.
10. Integrate shared documentation-policy results into Warlock repository health.
11. Add change-aware evaluation planning without weakening full-policy correctness.
12. Improve cross-snapshot, diagnosis-history, and baseline migration diagnostics.

See [Architecture diagnosis and policy authoring](architecture-diagnosis-and-policy-authoring.md) for the product workflow and detector direction.

## Acceptance criteria

Each planned feature must preserve deterministic evidence, explicit ownership, repository-only execution where applicable, schema validation, and independent CLI use. LLMs may consume Pitlord findings and perform repairs, but they must never be required to determine pass/fail state. Generalized advisory findings must remain distinguishable from blocking guard failures and repository-specific policy violations, and accepted architecture must remain distinct from an evidence baseline.

## Open decisions

- Exact persistence and versioning contract for accepted architectural decisions.
- Stable diagnosis identity across detector and snapshot changes.
- Which detector classes can promote into existing policy rules and which require new rule families.
- How reusable policy-pack versions are resolved and pinned by repositories.
- Whether Warlock invokes Pitlord directly or through a shared daemon job protocol.
- Which change-planning evidence is sufficient to safely skip unaffected graph regions.

## Implemented references

- The policy-free `scan` command and `pitlord.scan.v1` result envelope are implemented.
- The first generalized detector, `dependency-pressure`, reports anomalous outgoing cross-file dependency pressure from normalized Arcana relations without language-specific syntax assumptions. It uses explicit Arcana file nodes as the peer set, aggregates symbol/type/module relations onto those files, excludes virtual external namespaces and common non-production paths, and has completed its first frozen real-corpus calibration pass.
- The second generalized detector, `dependency-knots`, is implemented as an advisory strongly connected component judgment over static production-file dependencies. It uses iterative deterministic SCC traversal, Arcana semantic namespace or multi-file-module regions with directory fallback, a density fallback for narrowed single-region scopes, and synthetic modular/entangled/hub-heavy/layered/dense-subsystem calibration fixtures. Runtime `calls` are excluded from source-cycle direction, common test/addon tooling is separated, tiny conventional implementation seams are deferred, and bounded two-region components do not escalate solely from density. Against its frozen 16-corpus references, the tuned projection is 7 TP, 39 TN, 0 FP, 0 FN, 0 severity mismatches, and 0 unlabelled findings.
- `require_content` is implemented and documented in [Policy reference](../POLICY.md).
- Shared documentation governance is defined in [Engineering Standards](../../../engineering-standards/docs/INDEX.md).

## Related docs

- [Architecture diagnosis and policy authoring](architecture-diagnosis-and-policy-authoring.md)
- [Architecture](../ARCHITECTURE.md)
- [Current limitations](../limits/current-limitations.md)

## Notes

This roadmap is not a statement of shipped behavior.
