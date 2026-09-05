# Roadmap

Parent index: [Pitlord Planning](INDEX.md)

## Purpose

This document records planned Pitlord work and unresolved sequencing.

## Overview

Pitlord remains independently usable while expanding from architecture-policy enforcement into a deterministic, opinionated architecture and complexity governor for coding agents. Its default generalized guard should enforce broadly useful clean-code and clean-architecture constraints without requiring repository-specific policy, while repository policy continues to encode local intent.

The central product principle is deterministic supervision: coding agents may be probabilistic, but Pitlord's observation, judgment, evidence, and verification must be mechanically reproducible. Pitlord should prevent complexity accumulation without introducing another autonomous model-review loop around the agent.

## Current status

Implementation is underway. The policy-free `scan` command and deterministic `pitlord.scan.v1` finding contract are in place. The language-neutral `dependency-pressure` detector has completed its first 16-corpus calibration pass. The `dependency-knots` detector has completed tuning against its frozen 16-corpus detector-specific reference set: the current projection preserves all 7 required knots with 0 labelled false positives, 0 false negatives, 0 severity mismatches, and 0 unlabelled findings. The `dependency-depth` detector is implemented and independently calibrated across the same corpora: its 17-expectation projection scores 1 TP, 16 TN, 0 FP, 0 FN, 0 severity mismatches, and 0 unlabelled findings, with the required positive in JMH. The region-level `unstable-dependency-direction` detector is implemented and precision-calibrated across the same 16 corpora: its projection scores 0 TP, 16 TN, 0 FP, 0 FN, 0 severity mismatches, and 0 unlabelled findings; JMH generator/runtime and Maven compatibility-layer near misses anchor the current false-positive controls, while real-world recall remains unvalidated. The file/region-level `boundary-bypass` detector is implemented and false-positive calibrated across the same 16 corpora: its projection scores 0 TP, 16 TN, 0 FP, 0 FN, 0 severity mismatches, and 0 unlabelled findings; real-world positive recall remains unvalidated. The `boundary-cohesion` detector is also implemented and tuned against 30 frozen real-corpus negative controls with 0 false positives and 0 unlabelled findings; its real-world recall remains unvalidated because that projection has no required positive. The `hub-bottleneck` detector is implemented and tuned against a frozen 16-corpus projection with one required real-world positive: 1 TP, 43 TN, 0 FP, 0 FN, 0 severity mismatches, and 0 unlabelled findings. The `impact-blast-radius` detector has independent real-world positive coverage: its corrected two-lane 16-corpus projection scores 19 TP, 52 TN, 0 FP, 0 FN, 0 severity mismatches, and 0 unlabelled findings across 150 frozen expectations. The ninth generalized detector, `symbol-intermediary-bypass`, is implemented and calibrated: all 16 frozen clean corpora remain quiet, while the diff-pinned Space Rocks Homunculus development mutation produces the required warning with no unlabelled output. The tenth generalized detector, `cross-file-intermediary-bypass`, is also implemented and precision-calibrated: all 16 frozen clean corpora remain quiet, and the diff-pinned Space Rocks MatchDecision mutation produces exactly one required warning. Both intermediary-bypass positives were used while shaping their detectors, so independent validation/holdout recall remains required.

### Checkpoint — 2026-08-31

The generalized-detector expansion is paused here with ten detector families implemented. The latest work fixed Homunculus bypass-benchmark selection so requested development, validation, and holdout counts are chosen with deterministic constrained selection instead of a greedy pass that could starve later splits. Pitlord then added the separate `cross-file-intermediary-bypass` detector rather than weakening the same-file detector. Its frozen calibration is 1 TP / 16 TN / 0 FP / 0 FN: the MatchDecision mutation is a diff-pinned development positive, all 16 clean corpora remain quiet, and the stream-runtime mutation is deliberately excluded because the post-mutation graph does not establish a defensible ownership direction.

At that checkpoint, the planned sequence was:

1. **Finish recall evidence before adding more detector families.** Find or construct independent positive corpora for `boundary-cohesion`, file/region `boundary-bypass`, and `unstable-dependency-direction`; obtain validation/holdout positives for both intermediary-bypass detectors that were not used during tuning; and broaden `dependency-depth` beyond its single JMH positive.
2. **Decide detector disposition.** Classify each generalized detector as advisory-only or eligible for blocking use based on its calibration evidence, keeping insufficiently validated families advisory.
3. **Build the deterministic post-change guard loop.** Coding-agent change → Pitlord scan/guard → finite mechanical findings → repair → deterministic re-verification, with no model judgement in pass/fail.
4. **Harden reproducibility and disposition persistence.** Guarantee stable finding identity/order/severity for the same snapshot and add explicit human handling for intentional or ambiguous architecture before moving into user-friendly policy promotion.

Do not resume broad detector expansion by default from this checkpoint; calibration depth and the guard workflow are the immediate priorities.

### Checkpoint — 2026-09-05

Since the August checkpoint, Pitlord added a general analyzer registry with explicit graph requirements and semantic-capability gating, one graph-free native analyzer integration for Rust Clippy, and two graph-backed language-neutral semantic rules: `swallowed-error` and `unobserved-outcome`. Lexicon now supplies semantic capabilities, error-handler/action facts, downstream error-flow dispositions, and outcome obligations across Rust, TypeScript, JavaScript, and Python; Pitlord consumes those facts without owning language syntax.

The first real-source semantic calibration pass is complete enough to establish the current proof boundary but not broad precision or recall. Generated semantic sources are suppressed using trusted markers, Python import fallback is recognized as recovery, and downstream `fallback` / `enclosing-propagation` dispositions reduce Space Rocks Python `swallowed-error` findings from five to one while preserving Lexicanter's genuine Rust swallowed-error finding. `continuation` remains explicitly non-dispositive so arbitrary later work cannot hide a swallowed error. The frozen semantic rerun currently produces one Lexicanter Rust `swallowed-error`, one Space Rocks TypeScript `unobserved-outcome`, and one unresolved Space Rocks Python `swallowed-error` advisory.

The current planned sequence is:

1. **Broaden independent evidence before promotion.** Continue the architecture-detector recall work from the August checkpoint and add purpose-selected semantic corpora for Rust, TypeScript/JavaScript, and Python so neither detector nor semantic-rule disposition rests on development examples or a tiny language sample.
2. **Decide advisory versus blocking disposition from evidence.** Architecture detectors and semantic rules remain advisory unless independent calibration supports deterministic blocking use.
3. **Build the deterministic post-change guard loop.** Coding-agent change -> Pitlord scan/guard -> finite mechanical findings -> repair -> deterministic re-verification, with no model judgement in pass/fail.
4. **Harden reproducibility and human disposition.** Guarantee stable finding identity/order/severity for the same snapshot and add repository-owned handling for intentional or ambiguous architecture before policy promotion.
5. **Prefer differentiated semantic work.** New semantic rules should target repository-wide, cross-file, boundary, or cross-language problems. Do not expand native linter wrappers by default merely to duplicate ESLint, Ruff, or equivalent ecosystem diagnostics; the Clippy integration remains the normalization proof and useful Rust surface.

The semantic analyzer foundation does not reopen broad local-lint expansion. The immediate product priorities remain evidence quality, deterministic guard workflow, reproducibility, and repository-aware diagnosis.

## Expected ownership

- Pitlord owns anomaly interpretation, diagnosis identity, accepted architectural decisions, policy authoring semantics, rule evaluation, evidence identity, baselines, and machine-readable diagnostics.
- Arcana owns graph storage, traversal, communities, impact, paths, strongly connected components, and other structural measurements.
- Lexicon owns language facts.
- Pitlord owns codemap coverage and changed-code-to-owning-document guard semantics.
- Demon Docs owns Markdown structure/schema/navigation plus codemap extraction and target resolution.
- Warlock owns scheduling, enrolled-repository orchestration, and interactive diagnosis/policy presentation.

## Planned work

1. Add required real-world positive calibration coverage for boundary/cohesion, file/region boundary-bypass, and unstable-dependency-direction; add independent validation/holdout positives for both intermediary-bypass detectors; broaden dependency-depth beyond its single JMH positive; and add purpose-selected semantic corpora for `swallowed-error` and `unobserved-outcome` across the supported semantic languages before broader recall or guard-readiness claims.
2. Classify generalized detectors and semantic rules explicitly as advisory or blocking guard conditions and make every blocking finding a bounded, mechanically verifiable repair target for coding agents.
3. Add an immediate post-change guard workflow: agent changes code, Pitlord evaluates deterministically, the agent repairs the finite failures, and Pitlord verifies them without any model-based adjudication.
4. Guarantee reproducible findings for the same snapshot, Pitlord version, configuration, and generalized profile, including finding identity, ordering, severity, evidence, and exit status.
5. Add explicit human disposition for ambiguous advisory diagnoses: create a repository guardrail, accept the structure as intentional, or leave it unresolved.
6. Add user-friendly policy authoring that can promote a diagnosis into validated policy without requiring JSON-schema knowledge, including a current-repository match preview before publication.
7. Keep accepted architectural decisions semantically distinct from evidence baselines and define a durable repository-owned persistence contract for them.
8. Extend calibration evidence without weakening existing separation cases, and select any next semantic rule from demonstrated repository-wide, cross-file, boundary, or cross-language gaps rather than expanding per-language native-linter wrappers by default.
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
- Which repository-aware semantic fact families justify new Pitlord rules beyond the current local `swallowed-error` and `unobserved-outcome` proofs.

## Implemented references

- The policy-free `scan` command and `pitlord.scan.v1` result envelope are implemented.
- The first generalized detector, `dependency-pressure`, reports anomalous outgoing cross-file dependency pressure from normalized Arcana relations without language-specific syntax assumptions. It uses explicit Arcana file nodes as the peer set, aggregates symbol/type/module relations onto those files, excludes virtual external namespaces and common non-production paths, and has completed its first frozen real-corpus calibration pass.
- The second generalized detector, `dependency-knots`, is implemented as an advisory strongly connected component judgment over static production-file dependencies. It uses iterative deterministic SCC traversal, Arcana semantic namespace or multi-file-module regions with directory fallback, a density fallback for narrowed single-region scopes, and synthetic modular/entangled/hub-heavy/layered/dense-subsystem calibration fixtures. Runtime `calls` are excluded from source-cycle direction, common test/addon tooling is separated, tiny conventional implementation seams are deferred, and bounded two-region components do not escalate solely from density. Against its frozen 16-corpus references, the tuned projection is 7 TP, 39 TN, 0 FP, 0 FN, 0 severity mismatches, and 0 unlabelled findings.
- The third generalized detector, `boundary-cohesion`, is implemented as an advisory semantic-region judgment. Static dependency-like relationships measure cross-boundary spread while a broader relation family including calls and reads/writes measures internal collaboration. Candidates require six or more production files, broad independent target-region reach, top-decile peer reach, and both relatively and absolutely weak internal support; incoming-heavy regions are deferred to bottleneck analysis. Its 30 frozen real-corpus negative controls all remain quiet with 0 unlabelled findings, while deterministic synthetic topology supplies the current positive control. A required real-world positive corpus is still needed before recall or guard readiness can be claimed.
- The fourth generalized detector, `hub-bottleneck`, is implemented as an advisory file-level coordination-waist judgment. Extreme direct fan-in creates a candidate, but the finding requires many-to-many behavioral flow over `calls`, `reads`, and `writes`: at least six incoming behavioral regions, four outgoing behavioral regions, and an outgoing/incoming behavioral-degree ratio of at least 0.4. Reference-only hubs, shared contracts, base abstractions, narrow facades, and composition roots remain quiet. Against the frozen 16-corpus projection the tuned detector scores 1 TP, 43 TN, 0 FP, 0 FN, 0 severity mismatches, and 0 unlabelled findings; the one required real-world positive limits the strength of recall claims.
- The fifth generalized detector, `impact-blast-radius`, is implemented as an advisory two-lane file-level change-exposure judgment. Broad direct exposure uses repository-relative fan-in, cross-region breadth, and non-metadata dependency evidence; lower-direct-fan-in transitive exposure requires independently expanding first-hop branches and rejects dominant single-gateway closure. Against 150 independently audited frozen expectations the current detector scores 19 TP, 52 TN, 0 FP, 0 FN, 0 severity mismatches, and 0 unlabelled findings.
- The sixth generalized detector, `dependency-depth`, is implemented as an advisory dominant acyclic behavioral-corridor judgment. It excludes import/reference-only propagation, requires at least four hops across four architectural regions with three transitions, and terminates at cycles, shared convergence points, and highly reused central foundations. Its independently audited 16-corpus projection scores 1 TP, 16 TN, 0 FP, 0 FN, 0 severity mismatches, and 0 unlabelled findings; the required JMH corridor is the current real-world recall anchor.
- The seventh generalized detector, `boundary-bypass`, is implemented as an advisory file/region intermediary judgment. It requires at least three same-region peers to establish a gateway, minority direct behavioral access, material gateway concentration on the downstream region, and a non-shared target; import/reference-only edges and nested ownership trees remain quiet. Its frozen projection scores 0 TP, 16 TN, 0 FP, 0 FN with no severity or unlabelled mismatches, so real-world recall is not yet claimed. Same-file call-chain shortcuts remain intentionally outside this detector and are owned by the separate symbol-level detector.
- The eighth generalized detector, `unstable-dependency-direction`, is implemented as an advisory region-level Stable Dependencies Principle judgment. It uses static source-dependency direction, requires incoming-heavy source and outgoing-heavy target regions separated by at least a 0.25 instability gap, requires multi-file boundary support, and defers region SCCs to dependency-knot ownership. The raw 16-corpus pass produced no findings; independent near-miss audit retained JMH generator/runtime and Maven deprecated compatibility boundaries as absent controls. Its frozen projection scores 0 TP, 16 TN, 0 FP, 0 FN with no severity or unlabelled mismatches, so this is precision calibration rather than recall validation.
- The ninth generalized detector, `symbol-intermediary-bypass`, is implemented as an advisory resolved-call repeated-pipeline judgment. A direct caller-to-target edge is reportable only when an abandoned one-hop local wrapper is backed by at least three intact sibling wrapper paths plus a shared non-target support step; constructor targets and continued resolved same-name overload routes remain quiet. The 16 frozen clean-corpus references score 16 TN / 0 FP with no unlabelled findings, and the diff-pinned Space Rocks development mutation scores 1 TP / 0 FN at warning severity. Independent recall remains unvalidated because the positive shaped the detector.
- The tenth generalized detector, `cross-file-intermediary-bypass`, is implemented as an advisory same-region cross-file seam judgment. It requires a short one-hop intermediary with at least two remaining callers, a same-named peer in the intermediary file that still preserves the route, narrow downstream caller-file ownership, and distinct intermediary/target names. The 16 frozen clean-corpus references score 16 TN / 0 FP, while the diff-pinned Space Rocks MatchDecision development specimen scores 1 TP / 0 FN at warning severity with no unlabelled output. The stream-runtime Homunculus mutation remains intentionally outside scope because its post-mutation wrapper direction is statically ambiguous.
- The scan analyzer registry is implemented. Built-in groups declare graph requirements and normalized semantic capabilities, duplicate or unknown groups fail closed, and graph loading occurs only when selected analyzers require it.
- The graph-backed `semantic` analyzer group is implemented with `swallowed-error` and `unobserved-outcome`. Lexicon supplies language-neutral semantic facts for Rust, TypeScript, JavaScript, and Python; Pitlord owns reusable rule interpretation without parsing those languages. `swallowed-error` now consumes proven downstream error-flow dispositions while keeping generic `continuation` non-dispositive.
- Rust Clippy normalization is implemented as the first graph-free native analyzer. Pitlord preserves upstream Clippy rule identity, spans, and applicable structured suggestions, but broader per-language native-linter wrapper expansion is not a default roadmap priority.
- Real-source semantic calibration is recorded in [Semantic lint calibration](../development/calibration/semantic-lints.md). The current frozen rerun preserves one genuine Lexicanter Rust swallowed error, one Space Rocks TypeScript unobserved outcome, and one unresolved Space Rocks Python swallowed-error advisory; broader semantic precision/recall remains unclaimed.
- The blocking `pitlord docs` documentation guard is implemented. It consumes Demon Docs' schema-1 codemap export, Arcana's code-file inventory, and Git changes from the merge base through the current worktree; every code file must have at least one codemap owner and changed mapped code must change at least one owning document. Markdown structure and schema validation remain Demon Docs responsibilities.
- `require_content` is implemented and documented in [Policy reference](../POLICY.md).
- Shared documentation governance is defined in [Engineering Standards](../../../engineering-standards/docs/INDEX.md).

## Related docs

- [Architecture diagnosis and policy authoring](architecture-diagnosis-and-policy-authoring.md)
- [Architecture](../ARCHITECTURE.md)
- [Current limitations](../limits/current-limitations.md)

## Notes

This roadmap is not a statement of shipped behavior.
