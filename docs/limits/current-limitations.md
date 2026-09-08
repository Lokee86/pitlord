# Current Limitations

Parent index: [Pitlord Limits](INDEX.md)

## Purpose

This document records current Pitlord limitations that materially affect use or extension.

## Overview

Pitlord is functional as an on-demand CLI, but several product and integration surfaces remain intentionally incomplete.

## Limitations

### Generalized architecture detectors

The canonical per-detector evidence classification is maintained in [Diagnosis Calibration Baselines](../development/calibration/INDEX.md#architecture-detector-evidence-status). Evidence maturity is separate from runtime disposition; all ten detectors remain advisory.

`pitlord scan` currently implements ten generalized advisory detectors: calibrated language-neutral outgoing cross-file `dependency-pressure`, calibrated language-neutral `dependency-knots` over strongly connected static source-dependency components, calibrated language-neutral `dependency-depth` over dominant acyclic behavioral corridors, region-level `unstable-dependency-direction`, file/region-level `boundary-bypass`, same-file symbol-level `symbol-intermediary-bypass`, cross-file symbol-level `cross-file-intermediary-bypass`, `boundary-cohesion` over semantic architectural regions, calibrated language-neutral `hub-bottleneck` analysis over many-to-many behavioral coordination waists, and calibrated language-neutral `impact-blast-radius` analysis over broad direct exposure or independently amplified transitive dependent surfaces. All ten detectors have frozen detector-specific projections across the same 16 real-world corpora; symbol intermediary bypass and cross-file intermediary bypass additionally have diff-pinned real Space Rocks development mutation positives.

Sixteen frozen real-world corpora now have machine-readable detector references and can be evaluated reproducibly with `pitlord calibrate`. The current dependency-pressure projection scores 3 labelled true positives and 105 labelled true negatives with 0 labelled false positives, 0 labelled false negatives, and 0 severity mismatches across those pinned revisions. Twenty-six detector findings remain deliberately unlabelled—11 in Space Rocks, 13 in Maven, and 2 in JMH—and are not counted as successes. The detector now separates common non-production paths, requires both extreme production-peer fan-out and cross-directory boundary spread, exempts conventional composition seams, caps compatibility-path severity, and defers strongly reused central hubs to later detector families.

Dependency pressure remains advisory. Its findings still target physical files rather than semantic aggregates such as C# partial types or Go packages; conventional seam recognition is deterministic naming/role policy rather than semantic inference; directory regions are only an approximation of architectural boundaries; and the 26 unlabelled outputs still require future adjudication or coverage from other detector families before the real-world reference set is exhaustive.

The dependency-knot corpus audit remains frozen as the tuning reference. Across the 16 pinned corpora, the initial detector produced 29 findings: 7 required real knot-pressure components, 14 allowed bounded/intentional cyclic seams, and 8 false architectural cycles. The tuned detector preserves the 7 required knots while removing the frozen false positives and severity errors: the current projection scores 7 true positives, 39 true negatives, 0 false positives, 0 false negatives, 0 severity mismatches, and 0 unlabelled findings. It now excludes runtime `calls` from source-cycle direction, uses Arcana semantic namespace or multi-file-module regions before filesystem-directory fallback, recognizes GUT as test tooling, suppresses tiny conventional implementation seams, and does not promote bounded two-region cycles solely because they are dense. Remaining knot limitations are chiefly the physical-file finding scope, reliance on available Arcana semantic containers, and deterministic naming/path heuristics for conventional implementation seams.

The dependency-depth detector is advisory and has one independently audited required real-world positive. It propagates depth only through calls, inheritance/implementation, trait use, overrides, includes, explicit dependencies, and conversions; import-only and generic reference-only edges do not establish depth. A finding requires at least four hops across four architectural regions with three region transitions, while cycles, shared convergence points, and highly reused central foundations terminate the walk. Across 17 frozen expectations the current projection scores 1 TP / 16 TN / 0 FP / 0 FN with 0 severity mismatches and 0 unlabelled findings. The required JMH corridor provides real-world recall evidence, but a single positive does not establish broad recall across architectural styles.

The unstable-dependency-direction detector is deliberately conservative and currently has no required real-world positive. It applies a region-level Stable Dependencies Principle over static source dependencies: both regions need at least three files and three total region couplings, at least two source files must support the boundary, the source must be incoming-heavy, the target outgoing-heavy, and target instability must exceed source instability by at least 0.25. Region SCCs are deferred to dependency-knot ownership. The raw 16-corpus pass emitted zero findings, so a separate near-miss audit was required. Only JMH generator-core -> runner (`0.444 -> 0.625`, gap `0.181`) and Maven compatibility model-building -> model-io (`0.425 -> 0.667`, gap `0.242`) survived the cycle, region-size, coupling, and multi-source gates; source audit showed both are intentional dependencies, with Maven's case explicitly deprecated compatibility architecture. Relaxing the cross-half rule also surfaced normal adapter and shared-infrastructure dependencies. The frozen projection therefore scores 0 TP / 16 TN / 0 FP / 0 FN with 0 severity mismatches and 0 unlabelled findings. This validates precision controls only; real-world recall remains unvalidated.

The boundary-bypass detector is deliberately conservative and currently has no required real-world positive. It detects file-level behavioral shortcuts only when at least three same-region peers establish an intermediary, direct access is an exceptional minority, the intermediary is materially concentrated on the downstream region, and the downstream target is not already broadly shared. Import/reference-only edges and nested source/target ownership trees remain quiet. The first broad triad rule emitted 18 false gateway inferences across Gson and Space Rocks; after source audit and tuning, the frozen 16-corpus projection scores 0 TP / 16 TN / 0 FP / 0 FN with 0 severity mismatches and 0 unlabelled findings. This validates false-positive controls only. The known Space Rocks `handleDebugAddScore -> addDebugScoreForPlayer -> AddPlayerScore` mutation remains intentionally outside this file/region detector because both local symbols occupy one file; the separate symbol-level detector owns that case.

The `symbol-intermediary-bypass` detector owns resolved same-file call-chain shortcuts. It requires an abandoned one-hop local wrapper, at least three active sibling wrapper paths toward the same downstream file, and a shared support symbol outside the downstream target file; constructor targets and continued resolved same-name overload routes remain quiet. Its 16 whole-repository negative references score 0 TP / 16 TN / 0 FP / 0 FN with 0 severity mismatches and 0 unlabelled findings. A diff-pinned Space Rocks Homunculus development mutation supplies one required real-repository positive and scores 1 TP / 0 FN with no severity or unlabelled mismatches. That positive was used to shape the detector, so it does not establish unbiased recall. Validation and holdout Homunculus mutations remain required before claiming broader symbol-bypass recall.

The `cross-file-intermediary-bypass` detector owns a narrower resolved-call shape across source files. Its 16 whole-repository negative references score 0 TP / 16 TN / 0 FP / 0 FN, and its diff-pinned Space Rocks MatchDecision development specimen scores 1 TP / 0 FN with no unlabelled or severity mismatches. It requires caller, intermediary, and target to remain inside one semantic region; the intermediary must be a short one-hop wrapper with at least two remaining callers; a same-named peer in the intermediary file must still route through it; and the downstream target may be reached from at most two caller files after the shortcut. Same-name intermediary/target forwarding and broadly shared utilities remain quiet. Its first raw rule surfaced legitimate parallel APIs and utilities in Dapper, Polly, Spectre.Console, JMH, jsoup, Maven, and Space Rocks; source-audited suppressors reduce the final 16-corpus projection to 0 TP / 16 TN / 0 FP / 0 FN. The only frozen Space Rocks Homunculus mutation satisfying the final detector contract is the MatchDecision chain used during detector development, while the stream-runtime bypass is statically ambiguous after mutation. No independent recall score is therefore claimed yet.

The boundary/cohesion detector is deliberately conservative. It evaluates only semantic namespace or multi-file-module regions with at least six production files, measures outward spread using static dependencies, and measures internal support with a broader relationship family that includes runtime calls and read/write edges. A finding requires at least six outgoing cross-region dependencies, at least three independent target regions, at least 0.25 target regions per member, top-decile peer reach, internal support in the weakest peer quartile, and no more than 0.5 internal support relationships per member. Incoming-heavy regions are deferred to hub/bottleneck analysis. Against the 30 frozen real-corpus negative controls it produces 30 true negatives, 0 false positives, and 0 unlabelled findings. Those references contain no required real-world positive, so this result validates the current false-positive controls but not real-world recall; synthetic topology remains the positive control and the detector remains advisory rather than guard-ready.

The hub/bottleneck detector is also advisory but has one frozen required real-world positive. The untouched direct-centrality version emitted 37 findings and scored 1 TP / 10 TN / 33 FP / 0 FN with 2 severity mismatches. The tuned detector treats direct fan-in only as candidate generation; the decisive evidence is many-to-many behavioral flow over `calls`, `reads`, and `writes`. A candidate must receive behavior from at least six architectural regions, dispatch behavior to at least four regions, and have behavioral outgoing degree at least 40% of behavioral incoming degree. This suppresses central data/contracts, result/domain models, base abstractions, and narrow logging/telemetry facades while preserving the known Space Rocks networking gravity well. The current frozen projection scores 1 TP / 43 TN / 0 FP / 0 FN, 0 severity mismatches, and 0 unlabelled findings. Recall is still constrained by only one required real-world positive, so this is not broad proof of bottleneck recall across architectures.

The impact/blast-radius detector is advisory but now has corrected real-world recall coverage. Independent corpus auditing showed that the original all-negative calibration had confused healthy shared foundations with absence of impact. Against the corrected references, the untouched transitive rule scores 1 TP / 8 TN / 44 FP / 18 FN, while the historical branch-only tuning scores 0 TP / 52 TN / 0 FP / 19 FN: it removed inherited-closure noise by suppressing every required positive. The current two-lane detector restores broad direct exposure while retaining the branch discriminator for low-direct-fan-in transitive amplification. Direct exposure uses repository-relative fan-in, cross-region breadth, and non-metadata dependency evidence; transitive exposure still requires at least three substantial independent first-hop branches with no branch above 80% of reach. Across 150 frozen expectations the current projection scores 19 TP / 52 TN / 0 FP / 0 FN, with 0 severity mismatches and 0 unlabelled findings. The detector remains advisory because a broad change surface can be intentional even when its exposure is measured reliably.

**Removal condition:** boundary/cohesion, file/region boundary-bypass, and unstable-dependency-direction gain required real-world positive coverage, both intermediary-bypass detectors gain independent validation/holdout positive coverage beyond their development mutations, dependency-depth gains broader positive diversity beyond its single JMH corridor, and each generalized detector has sufficient recall evidence to support its intended advisory or guard disposition without weakening the current false-positive controls.

### Semantic analyzers and native-tool normalization

The graph-backed `semantic` analyzer group currently contains two advisory rules over Lexicon-normalized facts: `swallowed-error` and `unobserved-outcome`. `swallowed-error` suppresses a handler only when Lexicon proves a local propagate/record/recover action or a downstream `fallback`, `enclosing-propagation`, or `intentional-suppression` disposition. A generic `continuation` fact is retained as evidence but does not suppress the finding. `unobserved-outcome` reports only adapter-proven fallible or async operations whose outcome has no normalized consumption action.

Real-source semantic calibration is intentionally narrower than the 16-corpus architecture suite because most frozen architecture corpora do not exercise the implemented Rust, TypeScript/JavaScript, or Python semantic adapters. The three final audited findings are now machine-readable strict calibration references with exact path/span matching and fully-labelled regression enforcement: one genuine Lexicanter Rust `swallowed-error`, one Space Rocks TypeScript `unobserved-outcome`, and one intentionally unresolved Space Rocks Python `swallowed-error`. This is evidence that the current proof boundaries can remove known noise without hiding the known Lexicanter positive; three development-era examples are still not sufficient for broad precision or recall claims.

Two adapter-scale ceilings also remain visible from calibration. Lexicanter's checked-in generated JavaScript bundle is still loaded into the TypeScript compiler program before semantic-fact suppression, so whole-repository TypeScript analysis remains impractical for that corpus. Large Rust repositories can also exceed the current whole-repository adapter latency budget; Lore's 861-file Rust pass exceeded the five-minute calibration bound. These are Lexicon analysis limits exposed through Pitlord's semantic workflow, not Pitlord rule failures.

Pitlord also ships one graph-free native analyzer integration: Rust Clippy. It normalizes upstream diagnostics into the Pitlord finding contract but does not make Pitlord the owner of Clippy rule semantics. Equivalent ESLint, Ruff, or other per-language wrappers are not a current expansion priority; new semantic work should justify itself through repository-wide, cross-file, boundary, or cross-language value rather than duplicate ordinary local linting.

**Removal condition:** purpose-selected semantic corpora provide broader independent coverage for each supported language and rule; Lexicon resolves the generated-bundle and large-repository analysis ceilings relevant to Pitlord semantic scans; and any semantic rule promoted beyond advisory status has sufficient independent evidence for that disposition.

### Change-aware policy scope

Snapshot diff evaluates the full selected policy against both snapshots before comparing evidence. It does not yet restrict evaluation to changed policy scopes.

**Removal condition:** diff planning can safely identify affected rules and graph regions without missing introduced evidence.

### Documentation semantics

`pitlord docs` now owns two blocking repository-health checks: every Arcana code file must resolve from at least one authored codemap entry in the selected Demon Docs root, and every changed mapped code file must have at least one owning codemap document changed in the same Git range. Pitlord consumes Demon Docs' schema-1 codemap export and does not parse or validate Markdown itself.

Pitlord still does not own document structure or schema enforcement, index completeness, link resolution, document classes, heading/front-matter rules, or navigation health. Those checks are Demon Docs CI responsibilities. The guard currently evaluates the current resolved codemap dataset and current Arcana file inventory; historical reconstruction of a deleted file's prior codemap ownership is not a separate data source.

**Removal condition:** none currently planned for the ownership boundary; richer historical deletion evidence may be added if current-state coverage plus Git change evidence proves insufficient.

### Arcana query bounds and deadline policy

The Arcana process adapter honors caller cancellation through `exec.CommandContext`, but Pitlord does not yet define one product-wide default timeout. Callers using an unbounded context can therefore wait indefinitely for a stuck Arcana process.

Pitlord requests Arcana's current maximum of 10,000 relationships for each `neighbors` query and fails closed if Arcana reports truncation. The current `arcana.query.v1` neighbor operation has no continuation offset, so a single source node with more than 10,000 matching neighbors cannot yet be loaded completely even though ordinary large repositories can exceed the protocol's 1,000-item default safely.

**Removal condition:** CLI and Warlock invocation contracts adopt an explicit bounded deadline policy, and Arcana exposes continuation semantics for neighbor queries so Pitlord can load arbitrarily large bounded pages without truncation.

### Runtime integration

Pitlord has no independent daemon. Warlock may invoke it against changed snapshots, but central scheduling and UI compliance reporting are not yet integrated.

**Removal condition:** Warlock owns and ships the integration.

## Affected systems

- Snapshot diff and CI optimization.
- Warlock runtime and compliance UI.
- Cross-tool documentation governance.

## Status

All entries are current as of September 5, 2026.

## Related docs

- [Roadmap](../planning/roadmap.md)
- [Architecture](../ARCHITECTURE.md)
- [Shared adoption model](../../.standards/docs/adoption.md)

## Notes

The documentation-semantics boundary is deliberate: Pitlord owns codemap coverage and changed-code documentation guarding; Demon Docs CI owns document structure, schema, and navigation.
