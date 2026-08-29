# Impact / Blast-Radius Calibration Baseline

Parent index: [Diagnosis Calibration Baselines](INDEX.md)

## Purpose

This document records the corrected real-corpus calibration baseline for Pitlord's `impact-blast-radius` detector.

## Overview

The first calibration pass adjudicated only files emitted by the detector. That was methodologically incomplete: it could measure false positives, but it could not discover real hotspots that the detector never emitted. The reference set was therefore re-audited independently across all 16 frozen corpora before further tuning.

## Detector ownership

`impact-blast-radius` owns **change exposure through dependency reach**.

It is distinct from `hub-bottleneck`:

- `hub-bottleneck` asks whether many callers converge on a many-to-many behavioral coordination waist;
- `impact-blast-radius` asks whether changing a file exposes an unusually broad set of production dependents.

A cohesive or intentionally stable foundation can still be a real impact hotspot. A finding is advisory exposure, not a claim that the file should be split.

Two impact shapes belong to this detector:

1. **broad direct exposure** — an unusually large and architecturally broad set of production files directly depends on a shared foundation or contract; and
2. **transitive amplification** — moderate direct fan-in expands through multiple independent downstream branches into an unusually broad dependent surface.

A low-level leaf must not inherit an entire registry, catalog, or gateway's blast radius merely because one first-hop consumer is broadly reused.

## Historical untouched detector

The initial detector is commit `a21ba96` (`Add initial impact blast radius detector`). It computed transitive production-file dependents over Pitlord's normalized dependency relation family and required:

- at least eight production peers;
- at least eight transitive dependents;
- transitive reach across at least 25% of production peers;
- reach at or above the 95th percentile of active transitive impact;
- at least 2x transitive amplification over direct fan-in;
- impacted dependents spanning at least three architectural regions; and
- no family-wide condition where more than 12.5% of production files met the candidate rule.

Across the 16 frozen corpora it emitted **74 findings**:

| Corpus | Findings |
| --- | ---: |
| Dapper | 0 |
| Polly | 0 |
| Spectre.Console | 20 |
| DUnit LotusScript | 0 |
| Gson | 10 |
| HikariCP | 3 |
| JMH | 14 |
| jsoup | 4 |
| Maven | 20 |
| JSONParser LotusScript | 0 |
| kotlinx.coroutines | 1 |
| detekt | 1 |
| Now in Android | 0 |
| Lexicanter | 1 |
| Space Rocks Clean v1 | 0 |
| Volt MX LotusScript Toolkit | 0 |

Spectre.Console and Maven both reached the 20-finding cap.

## First adjudication defect

The original reference projection classified every emitted finding as `absent` and added one quiet control for zero-output corpora. That produced 81 all-negative expectations.

This was useful for identifying inherited-closure noise, but it was not valid recall calibration. Ground truth had been derived from detector output rather than from the corpora themselves. In particular, the projection incorrectly treated intentionally broad foundations and contracts as non-hotspots simply because they were architecturally healthy.

That contradicted the existing detector ownership decision recorded during dependency-pressure calibration: **high fan-in with little or no outgoing coupling was deliberately deferred to impact/blast-radius**.

The all-negative projection is therefore retained only as historical evidence of the first tuning mistake. It is no longer the calibration ground truth.

## Independent corpus audit

The corrected audit examined the frozen corpora independently of either impact detector implementation.

For each corpus the audit compared:

- direct production-file fan-in;
- transitive dependent reach and repository fraction;
- impacted architectural-region breadth;
- first-hop branch concentration;
- incoming relation families, separating metadata/import-only use from inheritance and behavioral use;
- existing manually adjudicated corpus architecture notes; and
- source for representative candidate contracts and foundations.

The audit intentionally distinguishes **impact** from **architectural defect**. Stable shared contracts may be `required` or `allowed` impact findings even when no refactor is warranted.

Representative required positives include:

- Polly `src/Polly.Core/ResilienceContext.cs` — shared per-execution pipeline state used across resilience strategies and telemetry;
- Spectre.Console `src/Spectre.Console/Rendering/IRenderable.cs` — foundational rendering contract;
- Gson `Gson.java`, `TypeAdapter.java`, and `stream/JsonReader.java` — core runtime and extension/streaming contracts;
- HikariCP `metrics/IMetricsTracker.java` — cross-cutting metrics contract;
- JMH `results/Result.java` — foundational benchmark-result contract;
- jsoup `Node.java`, `Element.java`, `Parser.java`, and `Validate.java` — DOM/parser/validation foundations with broad direct consumers;
- Maven `api/.../Session.java` and `impl/.../MavenProject.java` — central build/session and project-model contracts;
- kotlinx.coroutines `CoroutineDispatcher.kt` — core dispatcher extension contract;
- detekt `Rule.kt` — core rule extension base with hundreds of direct rule dependents;
- Lexicanter `stores.ts` and `types.ts` — application-wide state and domain contracts;
- Space Rocks `services/game-server/internal/game/game.go` — central game-runtime aggregate; and
- Volt MX LotusScript Toolkit `NotesHttpJsonRequestHelper.lss` — shared HTTP/JSON integration contract.

The audit also retains explicit negative controls for inherited or metadata-only cases, including Spectre.Console concrete border implementations, Maven marker annotations, kotlinx.coroutines annotation declarations, detekt marker annotations, and low-direct-fan-in helper leaves whose apparent reach is inherited through one downstream gateway.

Ambiguous but legitimately broad shared contracts are `allowed` rather than incorrectly asserted absent. This includes selected compatibility APIs, exception/value contracts, metrics/runtime watches, shared utilities, and central configuration surfaces.

DUnit LotusScript and JSONParser LotusScript do not provide meaningful file-level positive coverage because their production peer populations are too small. The frozen Now in Android Arcana graph exposes no cross-file dependency evidence for this detector and therefore remains a negative control rather than an invented positive.

## Corrected frozen projection

The 16 `*.impact-blast-radius.json` references now contain **150 expectations**:

- **19 `required`** real-world positives;
- **79 `allowed`** intentional or ambiguous broad maintenance watches; and
- **52 `absent`** negative controls.

The additional 36 `allowed` expectations were frozen from the independent direct-fan-in percentile/region audit before corpus scoring of the two-lane implementation. They cover clearly broad shared contracts such as console/runtime interfaces, serializer/parser contracts, benchmark parameters/results, session/project models, coroutine abstractions, and generated/shared service contracts that would be inappropriate to suppress merely to avoid unlabelled findings.

No label in this corrected projection was selected because a tuned detector emitted it. The required positives were selected from the independent graph/source audit, including many files that neither historical detector reported.

Against the exact untouched `a21ba96` implementation, the corrected projection scores:

```text
TP: 1
TN: 8
FP: 44
FN: 18
severity mismatches: 0
unlabelled findings: 0
```

The initial detector therefore had both precision and recall failures. Its only required hit was HikariCP `IMetricsTracker.java`; most required direct-impact foundations were suppressed by the `>=2x` amplification rule.

## Historical branch-gated tuning

Commit `b681dfa` (`Tune impact blast radius detector`) added an independent first-hop branch discriminator:

- at least 3 first-hop dependent branches with exclusive downstream contribution;
- each substantial branch contributes at least 5% of total transitive reach, with a two-file minimum; and
- the largest first-hop branch contributes no more than 80% of total transitive reach.

That correctly solved the inherited-closure failure mode. Concrete leaves beneath one catalog or gateway stopped receiving the gateway's whole downstream surface.

However, applying the branch rule to **all** impact candidates removed the direct-foundation ownership lane entirely. Against the corrected frozen projection, `b681dfa` scores:

```text
TP: 0
TN: 52
FP: 0
FN: 19
severity mismatches: 0
unlabelled findings: 0
```

The zero-false-positive result was therefore not a successful calibration result: it was achieved by suppressing every required real-world positive.

## Correct tuning requirement

Further tuning must preserve two independent evidence lanes:

### Direct exposure

A file may be an impact hotspot when unusually broad direct production fan-in reaches across the architecture, even when transitive/direct amplification is close to 1x. This is the detector family that owns the high-fan-in shared-foundation condition deferred by dependency pressure.

Metadata-only marker annotations should not become high-impact findings solely because many files import or annotate with them.

### Transitive amplification

For low/moderate direct fan-in, broad transitive closure still requires evidence that multiple first-hop branches independently contribute substantial downstream surface. The `b681dfa` branch-concentration discriminator remains useful for this lane.

### Non-goals

- Broad impact does not imply decomposition is required.
- A stable shared contract may remain intentionally central.
- A single catalog/gateway does not transfer its blast radius to every leaf beneath it.
- Metadata-only reuse is weaker impact evidence than behavioral, inheritance, or concrete contract dependency.

The detector remains advisory while the corrected real-world recall calibration is tuned.

## Related docs

- [Calibration evaluation harness](HARNESS.md)
- [Calibration references](references/INDEX.md)
- [Hub/bottleneck baseline](hub-bottleneck.md)
- [Architecture diagnosis plan](../../planning/architecture-diagnosis-and-policy-authoring.md)

## Notes

These are manually adjudicated calibration references, not universal architectural ground truth. Re-adjudicate when a pinned corpus revision changes or concrete source/graph evidence disproves a frozen judgment.