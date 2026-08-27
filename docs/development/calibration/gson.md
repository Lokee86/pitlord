# Gson Calibration Baseline

Parent index: [Diagnosis Calibration Baselines](INDEX.md)

## Purpose

This document freezes the manually adjudicated reference architecture state of the Gson corpus used for Pitlord dependency-pressure calibration.

## Overview

- Corpus: `java-gson`
- Source revision: `aebc51a56ca0793c13b841c29f73433b82446695`
- Commit subject: `Bump the maven group across 1 directory with 9 updates (#3071)`
- Language: Java
- Baseline role: external real-world control

The baseline records corpus state only. It does not prescribe detector changes.

## Current Pitlord observation

A repository-wide dependency-pressure scan reports 20 advisory findings: 3 high and 17 warning. Eight are production Gson library files. Eleven are functional tests and one is a shrinker test fixture.

## Manually adjudicated production state

| Scope | Expected interpretation | Notes |
| --- | --- | --- |
| `gson/src/main/java/com/google/gson/Gson.java` | intentional broad core API; maintenance watch | Main serialization/deserialization engine and public entry point. Adapter lookup/cache, reader/writer handling and conversion orchestration are intrinsic to the abstraction; its size and centrality make it a legitimate maintenance watch. |
| `gson/src/main/java/com/google/gson/GsonBuilder.java` | intentional broad configuration API; maintenance watch | Public builder for the full Gson configuration surface. Breadth is expected, but the large option/factory assembly surface is change-sensitive. |
| `gson/src/main/java/com/google/gson/internal/bind/ReflectiveTypeAdapterFactory.java` | intentional algorithmic seam; maintenance watch | Central reflective binding factory coordinating field discovery, access policy and adapter construction. Broad but cohesive. |
| `gson/src/main/java/com/google/gson/internal/Excluder.java` | clean | Focused exclusion/version/annotation policy. |
| `gson/src/main/java/com/google/gson/internal/bind/DefaultDateTypeAdapter.java` | clean | Focused date adapter implementation/factory. |
| `gson/src/main/java/com/google/gson/internal/bind/MapTypeAdapterFactory.java` | clean | Focused map adapter factory. |
| `gson/src/main/java/com/google/gson/internal/bind/TreeTypeAdapter.java` | clean | Focused tree serializer/deserializer adapter bridge. |
| `gson/src/main/java/com/google/gson/internal/bind/TypeAdapters.java` | clean; intentional adapter catalog | Registry/factory home for built-in type adapters; broad references are intrinsic to catalog ownership. |

No inspected production warning demonstrated unrelated responsibility sprawl. The central public engine, builder, and reflective binder are maintenance watches because of their role and size, not because their fan-out is unexpected.

## Non-production peer class

The eleven current findings under `gson/src/test/java/com/google/gson/functional/` are functional tests and are not production architecture findings in this baseline. `test-shrinker/src/main/java/com/example/ClassWithJsonAdapterAnnotation.java` belongs to the shrinker test fixture module and is also non-production.

## Ground-truth summary

Gson has a deliberately central serialization engine surrounded by specialized adapter factories. The broadest current findings correspond to those intended core abstractions. The audit freezes `Gson`, `GsonBuilder`, and `ReflectiveTypeAdapterFactory` as maintenance watches while treating the remaining warned production adapters as cohesive.

## Related docs

- [Diagnosis Calibration Baselines](INDEX.md)
- [Architecture diagnosis plan](../../planning/architecture-diagnosis-and-policy-authoring.md)

## Notes

This baseline is tied to the exact source revision above. Re-adjudicate only if that corpus revision changes or source evidence demonstrates that one of these frozen judgments is wrong.
