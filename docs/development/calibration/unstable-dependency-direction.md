# Unstable-Dependency-Direction Calibration

Parent index: [Diagnosis Calibration Baselines](INDEX.md)

## Purpose

This document records the first real-corpus audit and false-positive calibration of Pitlord's region-level `unstable-dependency-direction` detector.

## Overview

`unstable-dependency-direction` applies the Stable Dependencies Principle to Pitlord's language-neutral repository graph. It looks for static source dependencies that point from a comparatively stable architectural region toward a materially less stable region.

The audit used the same 16 pinned corpora as the other generalized detector families. The raw detector was run before detector-specific references were created. Its zero-output result was not treated as validation; the strongest suppressed and relaxed candidates were enumerated and source-audited before the negative controls were frozen.

## Owned shape

For an architectural region, Pitlord measures region-level instability as:

```text
I = outgoing regions / (incoming regions + outgoing regions)
```

A finding currently requires all of the following:

- source and target each contain at least three production files;
- source and target each participate in at least three independent incoming-plus-outgoing region couplings;
- at least two source files establish the source-to-target boundary;
- the source is incoming-heavy (`incoming > outgoing`, therefore `I < 0.5`);
- the target is outgoing-heavy (`outgoing > incoming`, therefore `I > 0.5`);
- target instability exceeds source instability by at least `0.25`; and
- source and target are not members of the same region-level strongly connected component.

The detector uses the same static source-dependency relation family as `dependency-knots`: imports, inheritance/implementation, trait use, overrides, includes, explicit dependencies, and conversions. Runtime calls do not define source dependency direction.

Region-level SCCs are deliberately deferred to `dependency-knots`; unstable direction is meant to describe an acyclic architectural gradient, not re-label a dependency cycle.

## Raw corpus result

The committed initial detector emitted no findings across all 16 pinned corpora:

```text
raw findings  0
corpora      16
```

That result created a recall question rather than a success claim. The audit therefore enumerated the structural gates independently.

Several repositories contained apparent stable-to-unstable region edges before SCC filtering. Polly, Gson, HikariCP, jsoup, kotlinx.coroutines, and Lexicanter all had such edges, but every one was inside a region-level strongly connected component. Those relationships already belong to cycle analysis and support retaining the SCC deferral.

After cycle deferral, region-size, coupling, and multi-source boundary support, only two cross-half near-misses survived:

```text
JMH    generator core -> runner   I 0.444 -> 0.625   gap 0.181
Maven  building       -> io       I 0.425 -> 0.667   gap 0.242
```

Both were source-audited before deciding whether to lower the `0.25` threshold.

## JMH near miss

The JMH boundary is supported by two generator files:

- `jmh-core/src/main/java/org/openjdk/jmh/generators/core/BenchmarkGenerator.java`;
- `jmh-core/src/main/java/org/openjdk/jmh/generators/core/CompilerControlPlugin.java`.

They depend on runner-side `BenchmarkList`/`BenchmarkListEntry`, `CompilerHints`, and related runtime metadata because the generator writes the benchmark-list and compiler-hint resources consumed by the runner. This is a coherent generator/runtime contract, not evidence that a stable domain boundary has accidentally taken ownership of volatile implementation details.

Lowering the threshold far enough to report the `0.181` gap would therefore create a false positive in the current corpus.

## Maven near miss

The strongest Maven boundary is in the compatibility model-builder tree:

```text
compat/maven-model-builder/.../model/building
    -> compat/maven-model-builder/.../model/io
```

Five files in `building` establish six dependencies into three `io` files. The measured instability gap is `0.242`, just below the detector threshold.

Source inspection shows that this is explicitly legacy compatibility architecture. `ModelProcessor` is deprecated since Maven 4.0.0, directs callers to `org.apache.maven.api.services.ModelBuilder`, and extends the legacy `ModelReader` contract from the `io` package. Reporting this boundary would penalize an intentional compatibility layer rather than identify newly inverted ownership.

This case is the closest observed supported near miss and is the main reason the `0.25` gap was not relaxed during this calibration.

## Relaxed positive-gap audit

A second audit removed the cross-half requirement and considered every supported acyclic dependency for which target instability was merely greater than source instability.

That broader rule exposed healthy dependency shapes rather than a stronger positive set. Representative examples include:

- JMH reflection adapters depending on generator-core abstractions (`I 0.333 -> 0.444`);
- Maven command-specific `mvnup`, `mvnenc`, and `mvnsh` implementations depending on shared invoker infrastructure;
- Maven project-building code depending on resolver infrastructure; and
- Maven implementation code depending on public XML service contracts.

No supported acyclic positive-gap boundary in any of the 16 corpora reached a `0.25` gap. Relaxing the source-stable/target-unstable half crossing would therefore broaden the detector into normal adapter and layering relationships without producing a source-validated violation.

Space Rocks contained three non-cyclic positive-gap boundaries under the relaxed enumeration, but none survived the minimum region/coupling/multi-source architectural-support requirements.

## Frozen projection

The 16 detector-specific references contain one absent expectation per corpus and no required positive:

```text
true positives        0
true negatives       16
false positives       0
false negatives       0
severity mismatches   0
unlabelled findings   0
```

JMH and Maven anchor the two strongest source-audited cross-half near misses. The remaining corpus references preserve quiet controls at the same pinned revisions.

This is **false-positive calibration only**. It does not validate real-world recall. The deterministic synthetic stable-to-unstable fixture remains the current positive control.

## Current interpretation

The first corpus pass supports the detector's major discriminators:

1. **Cycle ownership.** Region SCCs belong to `dependency-knots`, not unstable-direction analysis.
2. **Cross-half direction.** A generic positive instability gradient is too broad and captures healthy adapters and shared infrastructure.
3. **Boundary support.** One-file edges are too weak to establish an architectural direction judgment.
4. **Material instability gap.** The closest supported real-corpus near miss is an intentional deprecated compatibility boundary at `0.242`; lowering the current `0.25` threshold would immediately reduce precision.

None of those observations establish recall. A real positive is still required to determine whether the current rule is appropriately selective or overly conservative.

## Removal condition

Real-world recall claims require at least one independently adjudicated required unstable-direction positive, preferably a controlled real-repository mutation that reverses an otherwise stable dependency direction, while preserving the current 16 negative controls.

## Related docs

- [Calibration references](references/INDEX.md)
- [Dependency-knot baseline](dependency-knots.md)
- [Dependency-depth baseline](dependency-depth.md)
- [Architecture](../../ARCHITECTURE.md)
- [Current limitations](../../limits/current-limitations.md)

## Notes

An all-negative real-corpus projection is a precision constraint, not evidence that the detector is complete.
