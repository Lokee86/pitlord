# Now in Android Calibration Baseline

Parent index: [Diagnosis Calibration Baselines](INDEX.md)

## Purpose

This document freezes the manually adjudicated reference architecture state of the Now in Android corpus used for Pitlord dependency-pressure calibration.

## Overview

- Corpus: `kotlin-nowinandroid`
- Source revision: `7d45eae4f8720a0c77f507712ba2437ff974b6ed`
- Commit subject: `Merge pull request #2100 from FletchMcKee/graph-update`
- Language: Kotlin/Android
- Baseline role: external modular-application control

The baseline records corpus state only. It does not prescribe detector changes.

## Current Pitlord observation

A repository-wide dependency-pressure scan reports five advisory findings: one high and four warning. Four are under `core/data/src/test`. The only production finding is `core/designsystem/.../Theme.kt`, with four outgoing file dependencies.

## Manually adjudicated production state

| Scope | Expected interpretation | Notes |
| --- | --- | --- |
| `core/designsystem/src/main/kotlin/com/google/samples/apps/nowinandroid/core/designsystem/theme/Theme.kt` | clean | Cohesive Compose theme definition: color schemes, dynamic theming, background/gradient/tint composition locals and the app `MaterialTheme` wrapper. |

## Non-production peer class

The following current findings are tests, not production architecture:

- `core/data/src/test/kotlin/com/google/samples/apps/nowinandroid/core/data/repository/OfflineFirstNewsRepositoryTest.kt`
- `core/data/src/test/kotlin/com/google/samples/apps/nowinandroid/core/data/UserNewsResourceTest.kt`
- `core/data/src/test/kotlin/com/google/samples/apps/nowinandroid/core/data/repository/OfflineFirstTopicsRepositoryTest.kt`
- `core/data/src/test/kotlin/com/google/samples/apps/nowinandroid/core/database/model/PopulatedNewsResourceKtTest.kt`

## Ground-truth summary

The current dependency-pressure scan identifies no production architectural debt in this corpus. The only production warning is a cohesive design-system theme file; the other four findings are test-source breadth.

## Related docs

- [Diagnosis Calibration Baselines](INDEX.md)
- [Architecture diagnosis plan](../../planning/architecture-diagnosis-and-policy-authoring.md)

## Notes

This baseline is tied to the exact source revision above. Re-adjudicate only if that corpus revision changes or source evidence demonstrates that one of these frozen judgments is wrong.
