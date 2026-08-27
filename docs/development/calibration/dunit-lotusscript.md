# DUnit LotusScript Calibration Baseline

Parent index: [Diagnosis Calibration Baselines](INDEX.md)

## Purpose

This document freezes the manually adjudicated reference architecture state of the DUnit LotusScript corpus used for Pitlord dependency-pressure calibration.

## Overview

- Corpus: `dunit-lotusscript`
- Source revision: `e013a21cc511fcadc93afa1e9095ac9da0b50984`
- Commit subject: `Merge branch 'develop'`
- Language: LotusScript
- Baseline role: external real-world control

The baseline records corpus state only. It does not prescribe detector changes.

## Current Pitlord observation

The repository-wide dependency-pressure scan reports zero findings.

## Manually adjudicated production state

| Scope | Expected interpretation | Notes |
| --- | --- | --- |
| `ODP/Code/ScriptLibraries/ru.livescripts.dunit.core.lss` | clean; intentional framework core | Co-locates assertions, test context/result, abstract test and runner machinery in one small LotusScript library. The responsibilities remain within the test-framework core. |
| `ODP/Code/ScriptLibraries/ru.livescripts.dunit.matchers.lss` | clean | Matcher implementations belong to one assertion/matching library. |
| `ODP/Code/Agents/TestRunner.lsa` | clean; intentional orchestration seam | Thin runner adapter over the framework core. |
| `ODP/Code/Agents/DemoTestRunner.lsa` | non-production/demo breadth | Demonstration/test fixture that intentionally aggregates matcher and assertion examples. |

No obvious architectural debt was established by the inspected principal units.

## Ground-truth summary

DUnit is a compact testing framework. Its coarse LotusScript script-library units naturally contain multiple related classes, but the inspected responsibilities remain cohesive. Zero dependency-pressure findings is consistent with the frozen architecture state.

## Related docs

- [Diagnosis Calibration Baselines](INDEX.md)
- [Architecture diagnosis plan](../../planning/architecture-diagnosis-and-policy-authoring.md)

## Notes

This baseline is tied to the exact source revision above. Re-adjudicate only if that corpus revision changes or source evidence demonstrates that one of these frozen judgments is wrong.
