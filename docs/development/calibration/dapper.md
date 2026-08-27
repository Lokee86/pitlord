# Dapper Calibration Baseline

Parent index: [Diagnosis Calibration Baselines](INDEX.md)

## Purpose

This document freezes the human-adjudicated architecture state of the Dapper corpus used for Pitlord dependency-pressure calibration.

## Overview

- Corpus: `cs-dapper`
- Source revision: `72a54c475f75e18cb93cba0809d00a5e6e49efd9`
- Commit subject: `Add ReferenceTrimmer and remove unused references (#2191)`
- Language: C#/.NET
- Baseline role: external clean/real-world control

The baseline records corpus state only. It does not prescribe detector changes.

The frozen checkout is treated as the corpus identity; generated analysis state under ignored Warlock/Grimoire directories is not part of the source baseline.

## Current Pitlord observation

A repository-wide dependency-pressure scan reports 12 advisory findings: 7 production files, 3 test files, and 2 benchmark files. Scanning the major production projects independently currently returns no dependency-pressure findings for:

- `Dapper/`
- `Dapper.EntityFramework/`
- `Dapper.ProviderTools/`
- `Dapper.Rainbow/`
- `Dapper.SqlBuilder/`

The difference between repository-wide and project-scoped output is preserved as an observation, not treated here as ground truth.

## Human-adjudicated production state

| Scope | Expected interpretation | Notes |
| --- | --- | --- |
| `Dapper/DefaultTypeMap.cs` | clean | Cohesive default member/constructor mapping strategy. References to `SqlMapper` members cross physical files because `SqlMapper` is a partial type. |
| `Dapper/CommandDefinition.cs` | clean | Cohesive command value/setup abstraction. Current dependency count does not indicate responsibility sprawl. |
| `Dapper/DynamicParameters.cs` | clean with size watch | `DynamicParameters` is a single partial type spread across 3 files and about 474 lines. Broad references remain within its parameter-binding responsibility. |
| `Dapper/SqlMapper.GridReader.cs` + `Dapper/SqlMapper.GridReader.Async.cs` | clean with size watch | One `GridReader` partial type spread across 2 files and about 728 lines. The sync/async split is intentional physical decomposition of one abstraction. |
| `Dapper/SqlMapper.cs` + `Dapper/SqlMapper*.cs` partial fragments | genuine maintenance-pressure region | `SqlMapper` is one central public partial type spread across 26 source files and about 7,353 lines. Dapper's own historical notes record that `SqlMapper.cs` was split because the monolithic file had become too unmaintainable. This is real architectural/maintenance pressure at the aggregate type level, not independent pressure in each physical fragment. |
| `Dapper/SqlMapper.Async.cs` | part of `SqlMapper` aggregate | Do not adjudicate as an independent architectural unit. |
| `Dapper/SqlMapper.GridReader*.cs` | part of `GridReader` aggregate | Do not adjudicate each physical fragment independently. |

## Non-production peer classes

Repository-wide findings under `tests/` and `benchmarks/` are not production architecture findings in this baseline. Tests and benchmarks are separate peer classes and may legitimately touch broad production surfaces.

Current observed non-production findings include:

- `tests/Dapper.Tests/AsyncTests.cs`
- `tests/Dapper.Tests/MiscTests.cs`
- `tests/Dapper.Tests/ParameterTests.cs`
- `benchmarks/Dapper.Tests.Performance/Benchmarks.Belgrade.cs`
- `benchmarks/Dapper.Tests.Performance/LegacyTests.cs`

## Ground-truth summary

The Dapper production architecture should not be characterized as broadly unhealthy from the current repository-wide file fan-out findings.

The principal baseline condition to preserve is the large, central `SqlMapper` aggregate. It is legitimate maintenance pressure and should remain visible in corpus ground truth, while its many partial-class source files must not be treated as separate architectural responsibilities merely because they are separate files.

## Related docs

- [Diagnosis Calibration Baselines](INDEX.md)
- [Space Rocks Clean Control v1](space-rocks-clean-v1.md)
- [Architecture diagnosis plan](../../planning/architecture-diagnosis-and-policy-authoring.md)

## Notes

This baseline is tied to the exact source revision above. Re-adjudicate only if that corpus revision changes or source evidence demonstrates that one of these frozen judgments is wrong.
