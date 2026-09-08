# Boundary/Cohesion Calibration Baseline

Parent index: [Diagnosis Calibration Baselines](INDEX.md)

## Purpose

This document freezes the first manually adjudicated real-corpus reference state for Pitlord's `boundary-cohesion` detector before any real-corpus tuning.

## Overview

The baseline separates the untouched detector observation, manually adjudicated negative controls, frozen scoring projection, and tuned detector result. The 16 pinned corpora constrain false positives; deterministic synthetic topology currently supplies the positive control because no independently selected required real-world boundary/cohesion specimen is frozen yet.

## Untouched detector result

The initial detector commit is `9b80a6a` (`Add initial boundary cohesion detector`). It emitted **zero findings across all 16 frozen corpora**.

That result is not interpreted as evidence that all 16 repositories have perfect boundaries. The metric audit exposed a defect in the initial statistical rule: Tukey outlier fences were applied directly to bounded `[0,1]` boundary/cohesion ratios. In multiple corpora the outgoing upper fence exceeded `1.0` while the internal-cohesion lower fence was negative, so the conjunction could not be satisfied by any real region.

The frozen real-corpus projection therefore records representative **absent** controls only. It contains no `required` real-world positive. The deterministic synthetic fixture remains the positive control for this phase.

This means real-corpus precision can be calibrated, but real-corpus recall cannot yet be claimed for this detector. A future independently selected positive corpus is required before guard-mode promotion.

## Manually adjudicated controls

The audit focused on the strongest measured weak-cohesion/outward-coupling regions and checked whether they represented actual ownership failure or legitimate architecture shapes.

| Corpus | Control | Frozen interpretation |
| --- | --- | --- |
| Dapper | `Dapper` namespace / EntityFramework adapter | Partial-type/public-API decomposition and a small provider adapter. File/namespace projection is not a cohesion defect. |
| Polly | hedging strategy / legacy wrappers | Cohesive strategy-family and compatibility-wrapper seams. Independent strategy files can legitimately depend outward on shared resilience infrastructure. |
| Spectre.Console | rendering / CI enrichment | Rendering is a deliberate subsystem; CI enrichers are an implementation catalog. Low same-region static edge count is not enough to call either incohesive. |
| DUnit | demo runner agents | Demo/test runner seam; not production boundary debt. |
| Gson | SQL adapters / annotations / reflect API | Adapter and declarative API families. `reflect` is also strongly reused; fan-in belongs to later bottleneck analysis. |
| HikariCP | Micrometer metrics | Provider adapter seam. |
| JMH | reflection generator / benchmark infra | Cohesive backend adapter plus high-fan-in runtime API. |
| jsoup | safety / public facade | Cohesive layered safety subsystem and intentional public facade. |
| Maven | core IT support / model SPI / CLI invoker / lifecycle builder | Integration-test support, service API, composition, and orchestration seams. These are important false-positive controls for region statistics. |
| JSONParser LotusScript | `json.ls` | Single-unit implementation cannot establish a multi-file boundary/cohesion anomaly. |
| kotlinx.coroutines | common internal / Rx3 | Multiplatform runtime support and platform-adapter module. |
| detekt | ktlint/style rule packages | Large catalogs of independent rule implementations. High outward dependency count with little peer coupling is intrinsic to the plugin shape. |
| Now in Android | design-system theme | Cohesive design-system composition; current prepared graph does not expose a suspicious multi-file active region here. |
| Lexicanter | `src/app/layouts` | Physical presentation bucket, not an independently evidenced semantic ownership boundary. Existing `File.svelte`/`Settings.svelte` maintenance debt is file-responsibility pressure, not proof that the whole layouts directory is incohesive. |
| Space Rocks | shell / networking / Go game package | Composition shell, focused networking seam, and same-package decomposition. The clean control must remain quiet for this detector. |
| Volt MX LotusScript toolkit | EchoDoc agents | Small integration/example family; outward calls do not establish a broken production boundary. |

## Calibration implications

The frozen audit establishes several requirements for tuning without changing the reference labels:

- bounded ratios need percentile/rank comparisons rather than unconstrained Tukey fences;
- internal cohesion and cross-boundary dependency are not the same relation family;
- runtime calls may support evidence that files inside one region actually collaborate, even though runtime call edges must not define static dependency cycles;
- leaf catalogs/adapters can have very high outward ratios and near-zero internal static dependencies without being incohesive;
- high-fan-in central APIs belong to bottleneck/blast-radius analysis;
- physical directories are weak evidence of architectural ownership when no language namespace/module exists;
- broad composition, builder, invoker, compatibility, and platform-adapter seams require structural caution; and
- absolute boundary volume is less useful than **how many independent architectural regions a region reaches across**.

## Frozen scoring projection

The 16 `*.boundary-cohesion.json` references contain 30 `absent` expectations and no required/allowed expectations.

Against the untouched detector:

```text
TP: 0
TN: 30
FP: 0
FN: 0
severity mismatches: 0
unlabelled findings: 0
```

The harness reports a denominator-default precision/recall of `1.0`; those values are vacuous here because there are no labelled positive expectations. Do not cite them as detector recall evidence.

## Tuned detector result

The detector was tuned against the unchanged references in commit `84c2168` (`Tune boundary cohesion detector`). The final 16-corpus projection remains:

```text
TP: 0
TN: 30
FP: 0
FN: 0
severity mismatches: 0
unlabelled findings: 0
```

The first tuned pass surfaced two unlabelled warnings: JMH's `runner.format` namespace and Maven's compatibility `project.artifact` namespace. Source inspection showed both were cohesive families whose internal support was merely low relative to peers, not low in absolute terms. JMH had 1.17 internal support relationships per file and Maven had 0.62 per file. A general absolute weak-cohesion ceiling of 0.5 internal support relationships per file removed both without path- or framework-specific exceptions.

The tuned detector now:

- evaluates only Arcana semantic namespaces or multi-file modules rather than treating physical directories alone as ownership boundaries;
- uses static dependency-like relationships to measure outward boundary spread;
- uses a broader support family, including runtime calls and read/write relationships, to measure whether files inside a region actually collaborate;
- requires at least six production files and six outgoing cross-region dependencies;
- requires at least three independent target regions, at least 0.25 target regions per member, and top-decile target-region reach among semantic peers;
- requires internal support to be both in the weakest peer quartile and no more than 0.5 relationships per file; and
- defers incoming-heavy regions whose fan-in is at least their fan-out to bottleneck/blast-radius analysis.

The synthetic weak/outward region remains the positive control and the catalog, runtime-collaboration, incoming-hub, modular, entangled, hub-heavy, layered, and dense-subsystem synthetic controls remain quiet. The real-corpus reference set still contains no independently selected required positive, so real-world recall remains unvalidated and this detector remains advisory rather than guard-ready.

## Related docs

- [Calibration evaluation harness](HARNESS.md)
- [Calibration reference set](REFERENCE-SET.md)
- [Architecture diagnosis plan](../../planning/architecture-diagnosis-and-policy-authoring.md)

## Notes

These are manually adjudicated reference baselines, not independently human-reviewed ground truth. Re-adjudicate only when a pinned corpus revision changes or concrete source evidence disproves a frozen judgment.
