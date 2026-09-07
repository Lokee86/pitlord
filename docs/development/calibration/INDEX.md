# Diagnosis Calibration Baselines

Parent index: [Development Documentation](../INDEX.md)

## Purpose

This index owns manually adjudicated real-source calibration evidence for Pitlord's generalized architecture detectors and language-neutral semantic analyzers.

## Overview

Architecture calibration baselines pin an exact corpus revision and record the expected architectural interpretation independently of Pitlord's current detector output. Detector changes are evaluated against these labels instead of redefining the expected result when thresholds or graph heuristics change.

Semantic calibration records may instead adjudicate concrete rule findings and adapter evidence when the existing architecture corpus does not cover the required languages. In both cases, calibration evidence stays separate from current implementation claims and records limitations explicitly rather than treating warning count as ground truth.

## Direct files

- [Calibration evaluation harness](HARNESS.md) — Machine-readable `pitlord.calibration.v1` references, pinned-revision verification, labelled scoring, and command usage.
- [Calibration references](references/INDEX.md) — The 165 machine-readable scoring projections: ten architecture projections for each of the 16 frozen corpora, two diff-pinned Space Rocks intermediary-bypass mutation positives, and three strict semantic analyzer references.
- [Semantic lint calibration](semantic-lints.md) — Real-source calibration of `swallowed-error` and `unobserved-outcome` across Rust, TypeScript/JavaScript, and Python, including generated-source handling, downstream error-flow dispositions, and current corpus/scalability limits.
- [Dependency-knot baseline](dependency-knots.md) — First real-corpus knot audit, including required package knots, allowed cyclic seams, and false callback/directory cycles.
- [Dependency-depth baseline](dependency-depth.md) — Independent 16-corpus behavioral-depth audit with one required JMH cross-region corridor and 16 negative controls; final projection: 1 TP / 16 TN / 0 FP / 0 FN.
- [Boundary-bypass baseline](boundary-bypass.md) — File/region intermediary-bypass audit: 18 initial false gateway inferences, tuned shared-target/concentration/ownership controls, and a final 0 TP / 16 TN / 0 FP / 0 FN projection with recall explicitly unvalidated.
- [Symbol-intermediary-bypass baseline](symbol-intermediary-bypass.md) — Symbol-level repeated-pipeline bypass calibration: one diff-pinned real Space Rocks mutation positive, 16 whole-repository negative controls, and source-audited jsoup/Maven false-positive suppression.
- [Cross-file-intermediary-bypass baseline](cross-file-intermediary-bypass.md) — Conservative cross-file seam-bypass calibration: one diff-pinned Space Rocks MatchDecision development positive, 16 whole-repository negative controls, and source-audited suppression of overload, forwarding, utility, and symmetric-wrapper false positives.
- [Unstable-dependency-direction baseline](unstable-dependency-direction.md) — Region-level Stable Dependencies Principle audit: zero raw findings, source-audited JMH and Maven near misses, retained SCC/cross-half/support/gap discriminators, and a 0 TP / 16 TN / 0 FP / 0 FN precision projection with recall explicitly unvalidated.
- [Boundary/cohesion baseline](boundary-cohesion.md) — First real-corpus boundary/cohesion audit, including the untouched zero-output failure, 30 negative controls, and the limits of positive-corpus coverage.
- [Hub/bottleneck baseline](hub-bottleneck.md) — Centrality audit and tuning record: one required real-world maintenance gravity well, 33 labelled untouched false positives, and the tuned many-to-many behavioral-waist projection at 1 TP / 43 TN / 0 FP / 0 FN.
- [Impact/blast-radius baseline](impact-blast-radius.md) — Corrected independent corpus audit and two-lane tuning with 19 required positives, 79 allowed broad maintenance watches, and 52 negative controls; final projection: 19 TP / 52 TN / 0 FP / 0 FN.
- [Space Rocks Clean Control v1](space-rocks-clean-v1.md) — Manually adjudicated dependency-pressure reference baseline for the frozen multi-language Space Rocks control.
- [Dapper](dapper.md) — Manually adjudicated C#/.NET reference baseline covering project boundaries, partial types, and the central `SqlMapper` aggregate.
- [jsoup](jsoup.md) — Manually adjudicated Java reference baseline covering HTTP, DOM, parser, safety, and selector package responsibilities.
- [detekt](detekt.md) — Manually adjudicated Kotlin reference baseline covering core configuration validation, style-rule implementations, and test-source findings.
- [Lexicanter](lexicanter.md) — Manually adjudicated Svelte/TypeScript and Rust reference baseline with known application-layout maintenance debt.
- [Polly](polly.md) — C#/.NET resilience-library baseline covering strategy seams, state-machine pressure, compatibility APIs, tests, benchmarks, and snippets.
- [Spectre.Console](spectre-console.md) — C# console/UI baseline covering partial public facades, prompt/rendering controllers, catalogs, widgets, and tests.
- [Gson](gson.md) — Java serialization-library baseline covering public facades, reflection and adapter factories, registries, and test peers.
- [HikariCP](hikaricp.md) — Java connection-pool baseline covering pool lifecycle/runtime ownership and test peers.
- [JMH](jmh.md) — Java benchmark-framework baseline covering execution orchestrators and concentrated generator maintenance pressure.
- [Apache Maven](maven.md) — Large Java multi-module baseline covering compatibility debt, model/project construction pressure, CLI, plugin, lifecycle, and testing-support seams.
- [kotlinx.coroutines](kotlinx-coroutines.md) — Kotlin multiplatform concurrency baseline covering state machines, platform adapters, deprecated compatibility surfaces, and tests.
- [Now in Android](now-in-android.md) — Kotlin/Android modular-application baseline separating theme composition from test-source breadth.
- [DUnit LotusScript](dunit-lotusscript.md) — Small LotusScript testing-framework baseline with zero dependency-pressure findings and intentional runner seams.
- [JSONParser LotusScript](jsonparser-lotusscript.md) — Single-file LotusScript parser baseline preserving intra-unit parser/container maintenance pressure despite zero graph findings.
- [Volt MX LotusScript Toolkit](volt-mx-lotusscript-toolkit.md) — LotusScript integration baseline covering async workflow/state-contract maintenance pressure despite zero dependency-pressure findings.

## Baseline rules

- Pin an exact source commit or immutable specimen identity.
- Keep reference adjudication separate from current Pitlord output.
- Record why an expected finding is meaningful or why a graph outlier should be suppressed.
- Preserve legitimate pre-existing debt as baseline truth rather than modifying the corpus to make it appear clean.
- Treat tests, benchmarks, generated sources, composition roots, routers, orchestrators, and language/package structure as potentially distinct peer classes.
- Re-adjudicate only when the corpus identity changes or concrete evidence shows the frozen judgment was wrong.

## Related docs

- [Behavioral contract matrix](../behavioral-contract-matrix.md)
- [Architecture diagnosis plan](../../planning/architecture-diagnosis-and-policy-authoring.md)
- [Current limitations](../../limits/current-limitations.md)

## Notes

Synthetic Arcana topology tests remain deterministic detector tests. These baselines complement them with real-world architectural ground truth.