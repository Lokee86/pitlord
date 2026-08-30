# Boundary-Bypass Calibration

Parent index: [Diagnosis Calibration Baselines](INDEX.md)

## Purpose

This document records the independent real-corpus audit and false-positive calibration of Pitlord's file/region-level `boundary-bypass` detector.

## Overview

`boundary-bypass` looks for an exceptional direct behavioral dependency that reaches past an established intermediary between architectural regions. The detector is intentionally narrower than general cross-boundary coupling: it requires evidence that peer files from one upstream region normally use the same gateway and that the gateway is materially concentrated on the downstream region.

The audit used the same 16 pinned corpora as the other generalized detector families. Ground truth was established by source and topology inspection before the final detector-specific references were frozen.

## Owned shape

The detector owns file-level triads of this form:

```text
upstream peer files -> gateway file -> downstream region

exceptional upstream file ----------> downstream file
```

A finding requires all of the following:

- at least three peer files from the same upstream region depend on the gateway;
- direct access from that upstream region is an exceptional minority, at most 34% of the gateway-peer count;
- the direct edge is behavioral or explicit dependency evidence (`calls`, `reads`, `writes`, `depends-on`, `includes`, or `converts-to`), not import/reference-only reach;
- source, gateway, and downstream target occupy distinct architectural regions;
- the downstream target is not already shared from more than two source regions;
- the gateway reaches at least two files in the downstream region; and
- at least 40% of the gateway's outgoing behavioral targets are in that downstream region.

A source/target pair inside the same nested filesystem ownership tree is not treated as a bypass merely because a sibling region also touches the target.

## Initial corpus result

The first triad rule was deliberately broad and emitted 18 findings:

```text
Gson         6
Space Rocks 12
other 14     0
```

Source inspection showed that these were not defensible gateway violations. Gson public/shared abstractions such as `Gson`, `TypeAdapter`, `JsonPrimitive`, and `JsonReader` happened to reach internal helpers; their popularity did not make them mandatory gateways to those helpers. Space Rocks similarly made shared runtime state look like the required route to pickup and weapon types because many game files depended on that state.

The audit established two additional suppressors rather than accepting those findings:

1. **Downstream concentration.** A gateway must be materially specialized toward the downstream region, not merely popular upstream and incidentally connected downstream.
2. **Nested ownership.** Parent/child source trees such as Gson `internal.bind -> internal` and Space Rocks `game -> game/entities/...` do not establish a sibling gateway requirement.

The shared-target breadth rule also suppresses foundations reached independently from many architectural regions.

## Independent near-miss audit

After tuning, all 16 corpora were quiet. That result was not accepted as proof of recall. A separate near-miss enumeration inspected the strongest suppressed triads with the detector thresholds relaxed.

Representative controls included:

- JMH profiler paths reaching `TimeValue`, `Statistics`, and `ScoreFormatter`: targets are shared from three to seven regions rather than owned behind one gateway.
- Maven `DefaultSettingsBuilder -> ProtoSession` through `ProblemCollector`: the target is a shared API/session surface reached from four regions.
- jsoup routes through `Normalizer`, `Attributes`, and `Element`: downstream node/parser primitives are shared across several regions.
- Lexicanter layouts/components through `stores.ts` to `utils/layouts.ts`: the proposed gateway has only one target in that downstream utility region.
- Space Rocks rooms through game-mode resolution to team configuration: the team configuration target is independently reached from four regions.
- Gson `ReflectiveTypeAdapterFactory` through `ConstructorConstructor` to reflection helpers: gateway concentration is too weak to establish ownership.

These cases were frozen as absent regression controls where practical. No source-audited file/region bypass positive was found in the pinned corpora.

## Frozen projection

The 16 detector-specific references contain 16 absent expectations and no required positive:

```text
true positives        0
true negatives       16
false positives       0
false negatives       0
severity mismatches   0
unlabelled findings   0
```

This validates the current false-positive controls only. It does **not** validate real-world recall. The deterministic synthetic bypass fixture remains the current positive control.

## Symbol-level boundary

This first detector operates on file-level dependency edges and architectural regions. It does not own an intermediary bypass that occurs entirely between symbols in one file.

Space Rocks provides a concrete known case from prior mutation verification:

```text
handleDebugAddScore
-> addDebugScoreForPlayer
-> AddPlayerScore
```

A mutation that changed the first call to reach `AddPlayerScore` directly was a genuine symbol-level intermediary bypass, but both local functions live in the same source file. The current file-level detector cannot observe that lost intermediary edge. Symbol-level intermediary-bypass detection therefore remains separate follow-on work rather than being inferred from file topology.

## Removal condition

Real-world recall claims require at least one independently adjudicated required file/region bypass positive, preferably including a controlled real-repository mutation, while preserving the current 16 negative controls.

## Related docs

- [Calibration references](references/INDEX.md)
- [Boundary/cohesion baseline](boundary-cohesion.md)
- [Architecture](../../ARCHITECTURE.md)
- [Current limitations](../../limits/current-limitations.md)

## Notes

An all-negative real-corpus projection is a recall gap, not evidence that the detector is complete.
