# HikariCP Calibration Baseline

Parent index: [Diagnosis Calibration Baselines](INDEX.md)

## Purpose

This document freezes the manually adjudicated reference architecture state of the HikariCP corpus used for Pitlord dependency-pressure calibration.

## Overview

- Corpus: `java-hikaricp`
- Source revision: `a4d93f4f85517f90e632b795486d7102e933d7ff`
- Commit subject: `fixup circleci runner config`
- Language: Java
- Baseline role: external real-world control

The baseline records corpus state only. It does not prescribe detector changes.

## Current Pitlord observation

A repository-wide dependency-pressure scan reports three advisory warnings: two production pool classes and one test.

## Manually adjudicated production state

| Scope | Expected interpretation | Notes |
| --- | --- | --- |
| `src/main/java/com/zaxxer/hikari/pool/HikariPool.java` | intentional broad runtime seam; maintenance watch | Primary connection-pool runtime. Acquisition, eviction, housekeeping, pool state, connection creation, shutdown, metrics hooks and management operations are all part of pool ownership. Large and stateful, but cohesive. |
| `src/main/java/com/zaxxer/hikari/pool/PoolBase.java` | intentional broad runtime foundation; maintenance watch | Base pool infrastructure for data-source setup, connection validation/reset, network timeout and shared pool resources. Broad because it underpins the pool runtime. |

`src/test/java/com/zaxxer/hikari/pool/TestConnections.java` is test architecture and is excluded from production ground truth.

## Ground-truth summary

HikariCP deliberately concentrates connection-pool lifecycle behavior in `PoolBase` and `HikariPool`. Both are legitimate maintenance watches because they are central stateful runtime classes, but the audit did not establish unrelated responsibility sprawl from the current findings.

## Related docs

- [Diagnosis Calibration Baselines](INDEX.md)
- [Architecture diagnosis plan](../../planning/architecture-diagnosis-and-policy-authoring.md)

## Notes

This baseline is tied to the exact source revision above. Re-adjudicate only if that corpus revision changes or source evidence demonstrates that one of these frozen judgments is wrong.
