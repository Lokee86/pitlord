# Diagnosis Calibration Baselines

Parent index: [Development Documentation](../INDEX.md)

## Purpose

This index owns human-adjudicated architecture ground truth used to calibrate Pitlord's generalized diagnosis detectors.

## Overview

Calibration baselines pin an exact corpus revision and record the expected architectural interpretation independently of Pitlord's current detector output. Detector changes are evaluated against these labels instead of redefining the expected result when thresholds or graph heuristics change.

A baseline may distinguish real architectural pressure, intentional broad seams, clean regions, peer-class exceptions, and graph-granularity false positives. The target is architectural judgment quality, not a predetermined warning count.

## Direct files

- [Space Rocks Clean Control v1](space-rocks-clean-v1.md) — Human-adjudicated dependency-pressure baseline for the frozen multi-language Space Rocks control.
- [Dapper](dapper.md) — Human-adjudicated C#/.NET baseline covering project boundaries, partial types, and the central `SqlMapper` aggregate.

## Baseline rules

- Pin an exact source commit or immutable specimen identity.
- Keep human adjudication separate from current Pitlord output.
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