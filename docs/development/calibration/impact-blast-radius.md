# Impact / Blast-Radius Calibration Baseline

Parent index: [Diagnosis Calibration Baselines](INDEX.md)

## Purpose

This document freezes the first manually adjudicated real-corpus reference state for Pitlord's `impact-blast-radius` detector before real-corpus tuning.

## Overview

The detector owns **transitive change amplification**, not direct centrality. `hub-bottleneck` asks whether a file is a many-to-many coordination waist; `impact-blast-radius` asks whether a change to a file can propagate through an unusually large and architecturally broad dependent set.

The initial file-level closure is intentionally being calibrated conservatively. A file does not become an architectural hotspot merely because it sits below a public catalog, annotation API, base visitor, result model, exception hierarchy, or other stable shared abstraction. File aggregation can inherit the entire downstream closure of a gateway even when the leaf does not independently own that change surface.

The 16 pinned corpora contain **no required real-world positive for this detector yet**. The frozen references therefore constrain false positives, while the deterministic synthetic transitive-amplification fixture supplies the current positive control. Real-world recall is not claimed.

## Untouched detector

The initial detector is commit `a21ba96` (`Add initial impact blast radius detector`). It computes transitive production-file dependents over Pitlord's normalized dependency relation family and requires:

- at least eight production peers;
- at least eight transitive dependents;
- transitive reach across at least 25% of production peers;
- reach at or above the 95th percentile of active transitive impact;
- at least 2x transitive amplification over direct fan-in;
- impacted dependents spanning at least three architectural regions; and
- no family-wide condition where more than 12.5% of production files meet the candidate rule.

The detector is advisory. Its synthetic positive is a low-direct-fan-in foundation whose dependents expand through several downstream branches. A pure high-fan-in shared leaf remains quiet, and family-wide reachability patterns are suppressed instead of reported as isolated hotspots.

## Untouched real-corpus sweep

Across the 16 frozen corpora, the untouched detector emitted **74 findings**:

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

## Frozen adjudication

All 74 untouched findings are frozen **absent** for this detector. This does not assert that changing those files has zero downstream effect. It asserts that the file-level closure does not establish an independently owned architectural blast-radius hotspot.

The main false-positive classes are:

- **catalog-inherited impact** — Spectre.Console concrete border implementations and related rendering primitives inherit the downstream closure of the shared border catalog;
- **shared API/contracts** — Maven annotations, API enums/contracts, compatibility interfaces, Gson exceptions/tokens, kotlinx annotations, and detekt's minimal visitor base are deliberately stable shared surfaces;
- **shared framework utilities** — JMH statistics/file/temp utilities and HikariCP timing/metrics contracts have broad consumers without independently owning a large unstable subsystem;
- **shared parser/runtime utilities** — jsoup internal helpers sit low in parser/runtime dependency chains but are not independent architecture pressure points; and
- **application utility inheritance** — Lexicanter `layouts.ts` is a shared layout utility, while the corpus's previously frozen maintenance pressure is elsewhere in the application surface.

This aligns with the existing corpus audits. The known structural maintenance points in Gson, HikariCP, JMH, Maven, detekt, Lexicanter, and Space Rocks are different files or architectural conditions from the untouched impact candidates. The reference set therefore does not promote generic downstream reach into a diagnosis merely because a file is foundational.

Seven zero-output corpora retain one exact absent control each so later tuning cannot create new findings silently.

## Diagnostic relation audit

A temporary, uncommitted diagnostic build removed generic `references` and `annotates` from transitive propagation while keeping the initial thresholds and closure algorithm unchanged. The hypothesis was rejected early:

- Spectre.Console remained at the 20-finding cap;
- HikariCP remained at 3 findings;
- JMH only dropped from 14 to 12;
- Gson dropped from 10 to 4; and
- jsoup remained at 4.

The diagnostic sweep was stopped after that evidence was sufficient. Generic references contribute noise, but they are not the primary failure mode. The deeper issue is **inherited transitive closure through shared gateways and stable foundations**.

## Frozen scoring projection

The 16 `*.impact-blast-radius.json` references contain **81 expectations**:

- 0 `required`;
- 81 `absent`; and
- 0 `allowed`.

Against the untouched `a21ba96` detector, the frozen projection is:

```text
TP: 0
TN: 7
FP: 74
FN: 0
severity mismatches: 0
unlabelled findings: 0
```

Because the real-corpus projection has no required positive, labelled recall is not a meaningful validation claim. Synthetic topology remains the positive constraint while tuning targets the observed false-positive classes.

## Tuning requirements

Tuning must preserve the synthetic transitive-amplification positive without encoding corpus-specific names or paths. In particular:

- direct fan-in remains owned by `hub-bottleneck` unless transitive propagation adds independent impact;
- a single registry, catalog, facade, or gateway must not automatically transfer its entire downstream closure to every leaf beneath it;
- stable annotations, contracts, base types, exceptions, enums, and simple shared utilities must not become hotspots solely because many higher-level files ultimately depend on them;
- transitive impact should require evidence that the candidate independently supports multiple dependent branches or otherwise contributes non-trivially to the broad change surface; and
- the detector must remain advisory until real positive recall coverage exists.

## Related docs

- [Calibration evaluation harness](HARNESS.md)
- [Calibration references](references/INDEX.md)
- [Hub/bottleneck baseline](hub-bottleneck.md)
- [Architecture diagnosis plan](../../planning/architecture-diagnosis-and-policy-authoring.md)

## Notes

These are manually adjudicated calibration references, not independently reviewed universal ground truth. Re-adjudicate only when a pinned corpus revision changes or concrete source evidence disproves a frozen judgment.