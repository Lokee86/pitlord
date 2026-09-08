# Dependency Depth Calibration

Parent index: [Diagnosis Calibration Baselines](INDEX.md)

## Purpose

This document records the independent real-corpus audit and calibration of Pitlord's `dependency-depth` detector.

## Overview

`dependency-depth` owns narrow, acyclic behavioral dependency corridors rather than cycles, broad fan-out, or blast radius. The calibration used the same 16 pinned corpora as the other generalized detector families and established ground truth from corpus inspection before freezing detector-specific references.

The final projection contains one required real-world positive and 16 absent controls:

- 1 true positive;
- 16 true negatives;
- 0 false positives;
- 0 false negatives;
- 0 severity mismatches; and
- 0 unlabelled findings.

## Relation semantics

Depth propagates only through relationships that can represent behavioral or architectural delegation:

- `calls`;
- `implements`;
- `extends`;
- `uses-trait`;
- `overrides`;
- `includes`;
- `depends-on`; and
- `converts-to`.

Import-only and generic reference-only edges do not establish depth. This distinction was required by the corpus audit: an apparent four-hop Maven corridor depended on an import-only middle edge and disappeared when depth was restricted to behavioral/architectural relationships.

## Structural boundaries

A candidate must form a dominant acyclic corridor rather than merely participate in a long reachable path. The walk terminates at:

- dependency cycles, which belong to `dependency-knots`;
- shared convergence points, so a common downstream implementation does not transfer its depth to every caller; and
- highly reused central foundations, so shared infrastructure does not manufacture deep-chain findings upstream.

Branching remains eligible only when one branch is materially deeper than the alternatives. The current dominant-depth gap is two hops.

## Calibrated threshold

The final real-corpus separation is:

- at least 4 dependency hops;
- at least 4 architectural regions; and
- at least 3 region transitions.

After import-only propagation was removed, the 16-corpus audit produced one corridor satisfying all three conditions. Spectre.Console also contained a four-hop behavioral path, but it crossed only two architectural regions and remains a negative control. All other audited corpora topped out below the required cross-region depth.

## Required real-world positive

JMH contains the source-validated corridor:

```text
jmh-core/src/main/java/org/openjdk/jmh/Main.java
→ jmh-core/src/main/java/org/openjdk/jmh/runner/Runner.java
→ jmh-core/src/main/java/org/openjdk/jmh/runner/format/OutputFormatFactory.java
→ jmh-core/src/main/java/org/openjdk/jmh/runner/format/TextReportFormat.java
→ jmh-core/src/main/java/org/openjdk/jmh/results/format/ResultFormatFactory.java
```

The chain is behavioral rather than import-derived: `Main` constructs and executes `Runner`; `Runner` delegates output-format creation to `OutputFormatFactory`; that factory constructs `TextReportFormat`; and `TextReportFormat` delegates final result rendering to `ResultFormatFactory`.

The finding remains advisory. A four-hop cross-region corridor can be deliberate, but it is a useful deterministic maintenance signal because routine behavior traverses successive intermediary seams.

## Negative controls

Each frozen corpus contributes at least one independently selected absent control. Important separation cases include:

- Spectre.Console's four-hop exception/rendering pipeline, which remains contained to two architectural regions;
- Maven's former apparent four-hop corridor, whose critical middle hop was import-only;
- shorter cross-region command/controller pipelines in Space Rocks, detekt, Maven, and Kotlin coroutines;
- single-region inheritance or implementation chains; and
- import/reference-only and shared-foundation paths.

## Final result

Across the 16 pinned corpus revisions:

```text
TP  1
TN 16
FP  0
FN  0
severity mismatches 0
unlabelled findings 0
```

The detector therefore has one independently audited real-world recall case while preserving all frozen negative controls. That is sufficient for the current advisory disposition, but not evidence that every harmful deep architecture shape is covered.

## Related docs

- [Calibration reference set](REFERENCE-SET.md)
- [Calibration evaluation harness](HARNESS.md)
- [Architecture](../../ARCHITECTURE.md)
- [Current limitations](../../limits/current-limitations.md)
- [Architecture diagnosis plan](../../planning/architecture-diagnosis-and-policy-authoring.md)

## Notes

Thresholds are frozen from observed corpus separation rather than chosen to produce a target warning count. Future changes require new evidence, not detector-output-driven relabelling.
