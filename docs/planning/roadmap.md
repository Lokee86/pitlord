# Roadmap

Parent index: [Pitlord Planning](INDEX.md)

## Purpose

This document records planned Pitlord work and unresolved sequencing.

## Overview

Pitlord remains independently usable while expanding from architecture-policy enforcement into a deterministic, opinionated architecture and complexity governor for coding agents. Its default generalized guard should enforce broadly useful clean-code and clean-architecture constraints without requiring repository-specific policy, while repository policy continues to encode local intent.

The central product principle is deterministic supervision: coding agents may be probabilistic, but Pitlord's observation, judgment, evidence, and verification must be mechanically reproducible. Pitlord should prevent complexity accumulation without introducing another autonomous model-review loop around the agent.

## Current status

Implementation is underway. The policy-free `scan` command and deterministic `pitlord.scan.v1` finding contract are in place. The language-neutral `dependency-pressure` detector has completed its first 16-corpus calibration pass. The `dependency-knots` detector has completed tuning against its frozen 16-corpus detector-specific reference set: the current projection preserves all 7 required knots with 0 labelled false positives, 0 false negatives, 0 severity mismatches, and 0 unlabelled findings. The `dependency-depth` detector is implemented and independently calibrated across the same corpora: its 17-expectation projection scores 1 TP, 16 TN, 0 FP, 0 FN, 0 severity mismatches, and 0 unlabelled findings, with the required positive in JMH. The region-level `unstable-dependency-direction` detector is implemented and precision-calibrated across the same 16 corpora: its projection scores 0 TP, 16 TN, 0 FP, 0 FN, 0 severity mismatches, and 0 unlabelled findings; JMH generator/runtime and Maven compatibility-layer near misses anchor the current false-positive controls, while real-world recall remains unvalidated. The file/region-level `boundary-bypass` detector is implemented and false-positive calibrated across the same 16 corpora: its projection scores 0 TP, 16 TN, 0 FP, 0 FN, 0 severity mismatches, and 0 unlabelled findings; real-world positive recall remains unvalidated. The `boundary-cohesion` detector is also implemented and tuned against 30 frozen real-corpus negative controls with 0 false positives and 0 unlabelled findings; its real-world recall remains unvalidated because that projection has no required positive. The `hub-bottleneck` detector is implemented and tuned against a frozen 16-corpus projection with one required real-world positive: 1 TP, 43 TN, 0 FP, 0 FN, 0 severity mismatches, and 0 unlabelled findings. The `impact-blast-radius` detector has independent real-world positive coverage: its corrected two-lane 16-corpus projection scores 19 TP, 52 TN, 0 FP, 0 FN, 0 severity mismatches, and 0 unlabelled findings across 150 frozen expectations. The next planned architectural scanner is symbol-level intermediary bypass, using the known Space Rocks debug-score mutation as the first concrete granularity target.

## Expected ownership

- Pitlord owns anomaly interpretation, diagnosis identity, accepted architectural decisions, policy authoring semantics, rule evaluation, evidence identity, baselines, and machine-readable diagnostics.
- Arcana owns graph storage, traversal, communities, impact, paths, strongly connected components, and other structural measurements.
- Lexicon owns language facts.
- Pitlord owns codemap coverage and changed-code-to-owning-document guard semantics.
- Demon Docs owns Markdown structure/schema/navigation plus codemap extraction and target resolution.
- Warlock owns scheduling, enrolled-repository orchestration, and interactive diagnosis/policy presentation.

## Planned work

1. Add required real-world positive calibration coverage for boundary/cohesion, file/region boundary-bypass, and unstable-dependency-direction before making recall or guard-readiness claims for those detector families.
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
- The third generalized detector, `boundary-cohesion`, is implemented as an advisory semantic-region judgment. Static dependency-like relationships measure cross-boundary spread while a broader relation family including calls and reads/writes measures internal collaboration. Candidates require six or more production files, broad independent target-region reach, top-decile peer reach, and both relatively and absolutely weak internal support; incoming-heavy regions are deferred to bottleneck analysis. Its 30 frozen real-corpus negative controls all remain quiet with 0 unlabelled findings, while deterministic synthetic topology supplies the current positive control. A required real-world positive corpus is still needed before recall or guard readiness can be claimed.
- The fourth generalized detector, `hub-bottleneck`, is implemented as an advisory file-level coordination-waist judgment. Extreme direct fan-in creates a candidate, but the finding requires many-to-many behavioral flow over `calls`, `reads`, and `writes`: at least six incoming behavioral regions, four outgoing behavioral regions, and an outgoing/incoming behavioral-degree ratio of at least 0.4. Reference-only hubs, shared contracts, base abstractions, narrow facades, and composition roots remain quiet. Against the frozen 16-corpus projection the tuned detector scores 1 TP, 43 TN, 0 FP, 0 FN, 0 severity mismatches, and 0 unlabelled findings; the one required real-world positive limits the strength of recall claims.
- The fifth generalized detector, `impact-blast-radius`, is implemented as an advisory two-lane file-level change-exposure judgment. Broad direct exposure uses repository-relative fan-in, cross-region breadth, and non-metadata dependency evidence; lower-direct-fan-in transitive exposure requires independently expanding first-hop branches and rejects dominant single-gateway closure. Against 150 independently audited frozen expectations the current detector scores 19 TP, 52 TN, 0 FP, 0 FN, 0 severity mismatches, and 0 unlabelled findings.
- The sixth generalized detector, `dependency-depth`, is implemented as an advisory dominant acyclic behavioral-corridor judgment. It excludes import/reference-only propagation, requires at least four hops across four architectural regions with three transitions, and terminates at cycles, shared convergence points, and highly reused central foundations. Its independently audited 16-corpus projection scores 1 TP, 16 TN, 0 FP, 0 FN, 0 severity mismatches, and 0 unlabelled findings; the required JMH corridor is the current real-world recall anchor.
- The seventh generalized detector, `boundary-bypass`, is implemented as an advisory file/region intermediary judgment. It requires at least three same-region peers to establish a gateway, minority direct behavioral access, material gateway concentration on the downstream region, and a non-shared target; import/reference-only edges and nested ownership trees remain quiet. Its frozen projection scores 0 TP, 16 TN, 0 FP, 0 FN with no severity or unlabelled mismatches, so real-world recall is not yet claimed. Symbol-level intermediary bypass remains separate follow-on work.
- The eighth generalized detector, `unstable-dependency-direction`, is implemented as an advisory region-level Stable Dependencies Principle judgment. It uses static source-dependency direction, requires incoming-heavy source and outgoing-heavy target regions separated by at least a 0.25 instability gap, requires multi-file boundary support, and defers region SCCs to dependency-knot ownership. The raw 16-corpus pass produced no findings; independent near-miss audit retained JMH generator/runtime and Maven deprecated compatibility boundaries as absent controls. Its frozen projection scores 0 TP, 16 TN, 0 FP, 0 FN with no severity or unlabelled mismatches, so this is precision calibration rather than recall validation.
- The blocking `pitlord docs` documentation guard is implemented. It consumes Demon Docs' schema-1 codemap export, Arcana's code-file inventory, and Git changes from the merge base through the current worktree; every code file must have at least one codemap owner and changed mapped code must change at least one owning document. Markdown structure and schema validation remain Demon Docs responsibilities.
- `require_content` is implemented and documented in [Policy reference](../POLICY.md).
- Shared documentation governance is defined in [Engineering Standards](../../../engineering-standards/docs/INDEX.md).

## Related docs

- [Architecture diagnosis and policy authoring](architecture-diagnosis-and-policy-authoring.md)
- [Architecture](../ARCHITECTURE.md)
- [Current limitations](../limits/current-limitations.md)

## Notes

This roadmap is not a statement of shipped behavior.
