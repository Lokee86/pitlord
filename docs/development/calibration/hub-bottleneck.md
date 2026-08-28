# Hub/Bottleneck Calibration Baseline

Parent index: [Diagnosis Calibration Baselines](INDEX.md)

## Purpose

This document freezes the first manually adjudicated real-corpus reference state for Pitlord's `hub-bottleneck` detector before real-corpus tuning.

## Overview

The baseline distinguishes central coordination pressure from healthy centrality. A file is not unhealthy merely because many files depend on it: public contracts, domain roots, result models, base abstractions, logging/telemetry facades, and compatibility types can all be intentionally central. The detector should identify a narrower condition where broad reuse and meaningful outward behavioral responsibility converge into a maintenance gravity well.

The 16 pinned corpora contain one required real-world positive, four allowed maintenance-watch seams, and 43 absent controls. Unlike the boundary/cohesion projection, this baseline therefore provides a real positive recall constraint as well as false-positive controls.

## Untouched detector

The initial detector commit is `e1254da` (`Add initial hub bottleneck detector`). It evaluates production files using the broad normalized dependency projection and requires:

- at least eight active production peers;
- direct fan-in of at least eight and at least three times the active-peer median;
- fan-in at or above the 95th percentile;
- at least three incoming architectural regions;
- at least four outgoing production-file dependencies; and
- no family-wide condition where more than 12.5% of production files meet the candidate rule.

The detector is advisory. Pure high-fan-in leaves and pure high-fan-out composition roots are synthetic negative controls, while the hub-heavy and entangled topology fixtures expose four central two-sided hubs.

Across the 16 frozen real corpora, the untouched detector emitted **37 findings**:

| Corpus | Findings |
| --- | ---: |
| Dapper | 0 |
| Polly | 3 |
| Spectre.Console | 0 |
| DUnit LotusScript | 0 |
| Gson | 3 |
| HikariCP | 0 |
| JMH | 4 |
| jsoup | 3 |
| Maven | 20 |
| JSONParser LotusScript | 0 |
| kotlinx.coroutines | 1 |
| detekt | 0 |
| Now in Android | 0 |
| Lexicanter | 0 |
| Space Rocks Clean v1 | 3 |
| Volt MX LotusScript Toolkit | 0 |

Maven reached the detector's 20-finding cap. That is the strongest evidence that raw fan-in plus ordinary outgoing dependency breadth is not yet a useful bottleneck judgment for large framework cores.

## Frozen adjudication

### Required positive

`client/scripts/networking/client_connection_service.gd` in Space Rocks is the required positive. Its existing corpus audit already freezes it as a genuine maintenance gravity well: a 597-line networking facade with connection lifecycle, public outbound API, authentication/session composition, inbound/tooling relay, tracing, telemetry, recovery, observability, and transport metrics. Its broad centrality is intentional, but enough responsibility has accumulated that continued growth creates real maintenance pressure.

### Allowed maintenance watches

Four central seams may remain visible but are not required findings:

- Dapper `Dapper/SqlMapper.cs` — central partial aggregate with known maintenance pressure; the untouched detector is currently quiet here.
- Gson `Gson.java` — public serialization/deserialization engine whose breadth is intrinsic but change-sensitive.
- jsoup `Element.java` — broad public DOM abstraction whose centrality is intentional but large enough to watch.
- Maven API `Session.java` — deliberately broad public session facade; high severity from centrality alone is not justified.

The last three were emitted by the untouched detector. `Gson.java` and Maven `Session.java` were reported as `high`, which exceeds their frozen maximum `warning` severity.

### Absent false-positive classes

The remaining emitted candidates are frozen absent. Representative classes are:

- Polly `ResilienceStrategyTelemetry` — focused cross-cutting telemetry helper;
- Polly partial `AsyncPolicy.ContextAndKeys` — file-level projection of one legacy partial API;
- Polly `PolicyBase` — shared base abstraction;
- Gson `TypeAdapter` and `JsonReader` — core reusable extension/streaming abstractions;
- JMH `BenchmarkParams`, `Result`, `BenchmarkResult`, and `IterationResult` — shared parameter/result models;
- jsoup `Document` and `Node` — intentional DOM root/base abstractions;
- Maven API, compatibility, descriptor, request, project, artifact, and session models — mostly data/contracts whose centrality comes from type/reference use rather than coordination ownership;
- kotlinx.coroutines `AbstractCoroutine` — cohesive internal coroutine base abstraction;
- Space Rocks `internal/logging/logger.go` — intentional cross-cutting logging facade; and
- Space Rocks `runtime/asteroid.go` — cohesive game-domain entity whose file-level centrality is not an architectural bottleneck.

Zero-output corpora retain explicit absent controls so later tuning cannot create new warnings silently.

## Behavioral relation audit

Temporary diagnostic builds measured `calls`, `reads`, and `writes` on both sides of the 37 untouched candidates. The diagnostic instrumentation was not committed and did not alter the frozen initial detector or labels.

The audit exposed the structural discriminator that ultimately survived tuning. Many Maven candidates had **zero behavioral fan-out** despite large raw fan-in and outgoing dependency counts: `RemoteRepository`, API `Session`, API `Type`, `BuilderProblem`, `MessageBuilderFactory`, compatibility `Artifact`, `ArtifactRepository`, `ModelBuildingRequest`, `MavenExecutionRequest`, and `AbstractMojo` all fell into this class. Their apparent two-sided centrality was largely signature/type/reference structure.

Region breadth sharpened that distinction. Polly and all JMH candidates dispatched behavior into only one architectural region. Space Rocks logging received behavior from seven regions but dispatched into two; the asteroid domain entity was 3 incoming behavioral regions -> 2 outgoing. By contrast, the required Space Rocks networking facade formed a genuine many-to-many waist: **13 behavioral dependents from 8 regions -> 11 behavioral dependencies across 7 regions**.

Other high-centrality examples demonstrated why behavioral breadth alone still needed a balance constraint. Maven `RepositoryUtils` was 21 behavioral dependents from 15 regions -> 8 dependencies across 5 regions, while `MavenProject` was 65 from 24 -> 7 across 6. Their outward behavioral responsibility was small relative to how broadly they were consumed. Requiring behavioral outgoing degree to retain at least 40% of behavioral incoming degree separates these shared/core abstractions from the networking gravity well without corpus-specific names or paths.

## Frozen scoring projection

The 16 `*.hub-bottleneck.json` references contain:

- 1 `required` expectation;
- 4 `allowed` expectations; and
- 43 `absent` expectations.

Against the untouched `e1254da` detector:

```text
TP: 1
TN: 10
FP: 33
FN: 0
severity mismatches: 2
unlabelled findings: 0
```

Labelled precision is approximately **2.9%** and labelled recall is **100%**. Recall is constrained by only one required real-world positive, so the result demonstrates that the detector preserves that known maintenance gravity well; it is not broad proof of bottleneck recall across architectures.

## Tuned detector

The tuned implementation is commit `f38dea3` (`Tune hub bottleneck detector`). It preserves the initial direct-centrality gate as candidate generation, then requires a many-to-many behavioral waist.

A candidate must first satisfy the original broad-centrality conditions:

- at least eight active production peers;
- direct fan-in of at least eight and at least three times the active-peer median;
- fan-in at or above the 95th percentile;
- at least three incoming architectural regions;
- at least four outgoing production-file dependencies; and
- no family-wide condition where more than 12.5% of production files meet the candidate rule.

The decisive judgment then uses only `calls`, `reads`, and `writes`:

- behavioral dependents must span at least **6 architectural regions**;
- behavioral dependencies must span at least **4 architectural regions**; and
- behavioral outgoing degree must be at least **40%** of behavioral incoming degree.

High severity requires a stronger behavioral waist: at least 24 behavioral dependents, 16 behavioral dependencies, 8 incoming behavioral regions, and 6 outgoing behavioral regions. Findings remain advisory.

Synthetic calibration was adjusted to match the clarified claim. The positive fixture is now a true many-to-many behavioral hub. The previous hub-heavy and entangled fixtures, whose incoming centrality was primarily `references`, are negative controls for reference-driven centrality. Pure shared leaves, narrow behavioral facades, composition roots, modular, layered, and dense-subsystem topologies also remain quiet.

Against the frozen 16-corpus references, the tuned detector scores:

```text
TP: 1
TN: 43
FP: 0
FN: 0
severity mismatches: 0
unlabelled findings: 0
```

The required Space Rocks `client_connection_service.gd` finding is preserved. All 33 untouched labelled false positives are removed, including the Maven framework-contract explosion, and the four allowed maintenance-watch seams do not create scoring failures. No frozen reference was changed during tuning.

This is strong evidence that the many-to-many behavioral-waist rule controls the observed false-positive classes while retaining the known positive. It is still only one required real-world positive, so broad bottleneck recall across architectures is not established and the detector remains advisory.

## Ownership boundaries

- high fan-in alone is not a defect;
- broad type/signature/reference use is candidate evidence, not proof of coordination;
- shared base classes, domain models, request/result models, and narrow cross-cutting facades require behavioral evidence before they can be called bottlenecks;
- dependency-pressure continues to own isolated fan-out pressure;
- boundary/cohesion continues to own weak semantic-region boundaries;
- `hub-bottleneck` owns direct many-to-many behavioral coordination waists; and
- transitive impact/blast radius remains a separate detector family.

## Related docs

- [Calibration evaluation harness](HARNESS.md)
- [Calibration references](references/INDEX.md)
- [Space Rocks Clean Control v1](space-rocks-clean-v1.md)
- [Architecture diagnosis plan](../../planning/architecture-diagnosis-and-policy-authoring.md)

## Notes

These are manually adjudicated calibration references, not independently reviewed universal ground truth. Re-adjudicate only when a pinned corpus revision changes or concrete source evidence disproves a frozen judgment.
