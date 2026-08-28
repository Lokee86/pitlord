# Diagnosis Calibration Baselines

Parent index: [Development Documentation](../INDEX.md)

## Purpose

This index owns manually adjudicated reference architecture ground truth produced by calibration audits and used to calibrate Pitlord's generalized diagnosis detectors.

## Overview

Calibration baselines pin an exact corpus revision and record the expected architectural interpretation independently of Pitlord's current detector output. Detector changes are evaluated against these labels instead of redefining the expected result when thresholds or graph heuristics change.

A baseline may distinguish real architectural pressure, intentional broad seams, clean regions, peer-class exceptions, and graph-granularity false positives. The target is architectural judgment quality, not a predetermined warning count.

## Direct files

- [Calibration evaluation harness](HARNESS.md) — Machine-readable `pitlord.calibration.v1` references, pinned-revision verification, labelled scoring, and command usage.
- [Calibration references](references/INDEX.md) — The 64 machine-readable detector scoring projections for the 16 frozen corpora across dependency pressure, dependency knots, boundary/cohesion, and hub/bottleneck analysis.
- [Dependency-knot baseline](dependency-knots.md) — First real-corpus knot audit, including required package knots, allowed cyclic seams, and false callback/directory cycles.
- [Boundary/cohesion baseline](boundary-cohesion.md) — First real-corpus boundary/cohesion audit, including the untouched zero-output failure, 30 negative controls, and the limits of positive-corpus coverage.
- [Hub/bottleneck baseline](hub-bottleneck.md) — Centrality audit and tuning record: one required real-world maintenance gravity well, 33 labelled untouched false positives, and the tuned many-to-many behavioral-waist projection at 1 TP / 43 TN / 0 FP / 0 FN.
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