# JMH Calibration Baseline

Parent index: [Diagnosis Calibration Baselines](INDEX.md)

## Purpose

This document freezes the manually adjudicated reference architecture state of the JMH corpus used for Pitlord dependency-pressure calibration.

## Overview

- Corpus: `java-jmh`
- Source revision: `a194eead0136bb66e5e59e4fdb2e18543e730929`
- Commit subject: `7904196: Test PerfAsmMethodParsingTest.checkJDK fails with JDK-26`
- Language: Java
- Baseline role: external real-world control

The baseline records corpus state only. It does not prescribe detector changes.

## Current Pitlord observation

A repository-wide dependency-pressure scan reports five production findings: one high and four warnings, all under `jmh-core`.

## Manually adjudicated production state

| Scope | Expected interpretation | Notes |
| --- | --- | --- |
| `jmh-core/src/main/java/org/openjdk/jmh/runner/Runner.java` | intentional execution orchestrator; maintenance watch | Public JMH runner coordinates benchmark discovery, option expansion, locking, profiler setup, action plans, execution and result output. Broad and large, but those responsibilities belong to benchmark-run orchestration. |
| `jmh-core/src/main/java/org/openjdk/jmh/runner/BaseRunner.java` | intentional runtime seam; maintenance watch | Shared lower-level benchmark execution/fork lifecycle. Central and stateful but cohesive with the runner subsystem. |
| `jmh-core/src/main/java/org/openjdk/jmh/generators/core/BenchmarkGenerator.java` | real structural maintenance pressure | A ~942-line generator coordinator spanning benchmark discovery/validation, grouping, state integration, and substantial source emission. All work belongs to code generation, but multiple change surfaces are co-located in one very large unit. |
| `jmh-core/src/main/java/org/openjdk/jmh/generators/core/BenchmarkGeneratorUtils.java` | clean; algorithm utility seam | Utilities supporting the generator's source/model analysis. |
| `jmh-core/src/main/java/org/openjdk/jmh/generators/core/StateObjectHandler.java` | real structural maintenance pressure | A ~934-line state-generation coordinator spanning state validation, dependency/cycle resolution, binding, lifecycle wiring, auxiliary counters, and source emission. Cohesive at subsystem level but materially concentrated. |

The two generator classes are frozen as genuine decomposition/maintenance pressure, not as unrelated-responsibility ownership violations.

## Ground-truth summary

JMH contains intentionally broad orchestration at both ends of its lifecycle. Runner breadth matches benchmark-execution ownership. The generator subsystem also contains two materially concentrated ~900-line units whose responsibilities remain related but create real structural maintenance pressure. Dependency breadth alone does not distinguish those cases; the source-level concentration does.

## Related docs

- [Diagnosis Calibration Baselines](INDEX.md)
- [Architecture diagnosis plan](../../planning/architecture-diagnosis-and-policy-authoring.md)

## Notes

This baseline is tied to the exact source revision above. Re-adjudicate only if that corpus revision changes or source evidence demonstrates that one of these frozen judgments is wrong.
