# kotlinx.coroutines Calibration Baseline

Parent index: [Diagnosis Calibration Baselines](INDEX.md)

## Purpose

This document freezes the manually adjudicated reference architecture state of the kotlinx.coroutines corpus used for Pitlord dependency-pressure calibration.

## Overview

- Corpus: `kotlin-coroutines`
- Source revision: `04ada74fae2e8914ae92ece34e06e80bb15385e9`
- Commit subject: `Relax the contributing rules on changing the documentation`
- Language: Kotlin multiplatform
- Baseline role: external real-world control

The baseline records corpus state only. It does not prescribe detector changes.

## Current Pitlord observation

A repository-wide dependency-pressure scan reports 20 advisory findings: 2 high and 18 warning. Nineteen are production source files and one is a JVM stress test. Reported production fan-out is only 4-6 files per finding.

## Manually adjudicated production state

### Core algorithm/runtime maintenance watches

| Scope | Expected interpretation | Notes |
| --- | --- | --- |
| `kotlinx-coroutines-core/common/src/flow/SharedFlow.kt` | intentional algorithmic seam; maintenance watch | SharedFlow public contract and synchronized replay/buffer/subscriber implementation form one complex concurrency abstraction. |
| `kotlinx-coroutines-core/common/src/flow/operators/Share.kt` | intentional flow-sharing orchestrator; maintenance watch | Converts cold flows into shared/state flows and manages sharing lifecycle. |
| `kotlinx-coroutines-core/common/src/CancellableContinuationImpl.kt` | intentional concurrency runtime seam; maintenance watch | Core cancellable-continuation state machine. |
| `kotlinx-coroutines-core/common/src/JobSupport.kt` | intentional concurrency runtime seam; maintenance watch | Central Job lifecycle/state machine. |
| `kotlinx-coroutines-core/common/src/flow/StateFlow.kt` | intentional algorithmic seam; maintenance watch | StateFlow contract and implementation are central concurrent-flow machinery. |

### Compatibility maintenance state

| Scope | Expected interpretation | Notes |
| --- | --- | --- |
| `kotlinx-coroutines-core/common/src/channels/BroadcastChannel.kt` | compatibility maintenance debt | Deprecated `BroadcastChannel`/`ConflatedBroadcastChannel` implementations remain as a substantial historical API surface alongside the newer SharedFlow model. This is real retained compatibility debt, not an unrelated-responsibility violation. |
| `kotlinx-coroutines-core/common/src/channels/Broadcast.kt` | compatibility maintenance watch | Deprecated broadcast builder surface layered over coroutine/channel machinery. |
| `kotlinx-coroutines-core/jvm/src/channels/Actor.kt` | compatibility maintenance watch | Cohesive channel-backed actor API surface, but a legacy-style API seam worth keeping distinct from core runtime findings. |

### Cohesive production findings

The following current findings are frozen as clean/intentional implementations within their platform or feature ownership:

- `kotlinx-coroutines-core/common/src/AbstractCoroutine.kt`
- `kotlinx-coroutines-core/common/src/Builders.common.kt`
- `kotlinx-coroutines-core/common/src/channels/Produce.kt`
- `kotlinx-coroutines-core/common/src/internal/LimitedDispatcher.kt`
- `kotlinx-coroutines-core/common/src/sync/Mutex.kt`
- `kotlinx-coroutines-core/jvm/src/CoroutineContext.kt`
- `kotlinx-coroutines-core/jvm/src/Executors.kt`
- `kotlinx-coroutines-core/jvm/src/internal/MainDispatchers.kt`
- `kotlinx-coroutines-core/jvm/src/scheduling/Dispatcher.kt`
- `kotlinx-coroutines-core/web/src/internal/JSDispatcher.kt`
- `ui/kotlinx-coroutines-android/src/HandlerDispatcher.kt`

Their dependency counts reflect core/platform adaptation and coroutine primitives rather than demonstrated unrelated responsibility sprawl.

## Non-production peer class

`kotlinx-coroutines-core/jvm/test/scheduling/WorkQueueStressTest.kt` is a stress test and is not production architecture in this baseline.

## Ground-truth summary

kotlinx.coroutines contains legitimately complex concurrent state machines and flow implementations, but the current detector is firing at very small absolute file fan-out values. The frozen architecture state treats the major state machines as maintenance watches, preserves the deprecated BroadcastChannel family as real compatibility debt, and treats the remaining warned production files as cohesive feature/platform implementations. No current finding is frozen as unrelated-responsibility ownership pathology.

## Related docs

- [Diagnosis Calibration Baselines](INDEX.md)
- [Architecture diagnosis plan](../../planning/architecture-diagnosis-and-policy-authoring.md)

## Notes

This baseline is tied to the exact source revision above. Re-adjudicate only if that corpus revision changes or source evidence demonstrates that one of these frozen judgments is wrong.
