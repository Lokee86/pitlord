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

A temporary diagnostic build measured only outgoing `calls`, `reads`, and `writes` for the 37 untouched candidates. The diagnostic instrumentation was not committed and does not alter the frozen initial detector.

The audit exposed a strong structural distinction. Many Maven candidates have **zero behavioral fan-out** despite large raw fan-in and outgoing dependency counts: `RemoteRepository`, API `Session`, API `Type`, `BuilderProblem`, `MessageBuilderFactory`, compatibility `Artifact`, `ArtifactRepository`, `ModelBuildingRequest`, `MavenExecutionRequest`, and `AbstractMojo` all fall into this class. Their apparent two-sided centrality is largely signature/type/reference structure.

By contrast, the required Space Rocks networking facade has 11 behavioral outgoing dependencies against 13 direct dependents. Other notable measurements were Gson `Gson` 15/28, jsoup `Element` 24/30, Space Rocks logging 5/24, and Space Rocks asteroid 6/17 (behavioral outgoing/direct incoming).

This supports a tuning direction based on a **two-sided behavioral waist** rather than raw graph degree. The frozen labels must remain unchanged while that hypothesis is tested.

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

## Tuning requirements

Tuning must preserve the frozen positive and synthetic topology separation without encoding corpus-specific names or paths. In particular:

- high fan-in alone is not a defect;
- broad type/signature/reference use must not be mistaken for behavioral coordination;
- shared base classes, domain models, request/result models, and cross-cutting logging/telemetry facades require structural evidence beyond centrality;
- dependency-pressure continues to own isolated fan-out pressure;
- boundary/cohesion continues to own weak semantic-region boundaries;
- transitive blast radius remains a separate detector family; and
- any stronger bottleneck rule should be grounded in two-sided behavioral flow, architectural-region breadth, or another deterministic structural discriminator rather than naming conventions.

## Related docs

- [Calibration evaluation harness](HARNESS.md)
- [Calibration references](references/INDEX.md)
- [Space Rocks Clean Control v1](space-rocks-clean-v1.md)
- [Architecture diagnosis plan](../../planning/architecture-diagnosis-and-policy-authoring.md)

## Notes

These are manually adjudicated calibration references, not independently reviewed universal ground truth. Re-adjudicate only when a pinned corpus revision changes or concrete source evidence disproves a frozen judgment.
