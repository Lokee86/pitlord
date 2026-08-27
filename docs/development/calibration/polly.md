# Polly Calibration Baseline

Parent index: [Diagnosis Calibration Baselines](INDEX.md)

## Purpose

This document freezes the manually adjudicated reference architecture state of the Polly corpus used for Pitlord dependency-pressure calibration.

## Overview

- Corpus: `cs-polly`
- Source revision: `9fd716c1e84288a6a8a1c7ad9426fb8a8e66942d`
- Commit subject: `Bump zizmorcore/zizmor-action from 0.5.7 to 0.6.0 (#3172)`
- Language: C#/.NET
- Baseline role: external real-world control

The baseline records corpus state only. It does not prescribe detector changes.

## Current Pitlord observation

A repository-wide dependency-pressure scan reports 20 advisory findings: 10 high and 10 warning. Fifteen are production library files. The remaining five are a benchmark, a documentation snippet, and three test files.

## Manually adjudicated production state

| Scope | Expected interpretation | Notes |
| --- | --- | --- |
| `src/Polly.Core/CircuitBreaker/CircuitBreakerResiliencePipelineBuilderExtensions.cs` | clean; intentional composition seam | Focused builder/factory surface that assembles circuit-breaker options, controller, telemetry, and strategy. |
| `src/Polly.Core/CircuitBreaker/Controller/CircuitStateController.cs` | maintenance watch | Cohesive circuit-breaker state machine, but necessarily stateful and change-sensitive: locking, transitions, timing, callbacks, telemetry, outcomes, and manual control converge here. |
| `src/Polly.Core/Hedging/Controller/TaskExecution.cs` | maintenance watch | Cohesive hedging-attempt lifecycle coordinator with cancellation, contexts, timing, outcomes, and callbacks. |
| `src/Polly.Core/Timeout/TimeoutResilienceStrategy.cs` | clean; intentional strategy seam | Timeout generation, cancellation, telemetry, callbacks, and outcome handling all belong to the timeout strategy. |
| `src/Polly/AsyncPolicy.ExecuteOverloads.cs` | intentional broad compatibility API; maintenance watch | One fragment of the partial legacy `AsyncPolicy` type. Its overload matrix is deliberate public API compatibility rather than unrelated ownership, but the large compatibility surface remains costly to maintain. |
| `src/Polly.Core/Fallback/FallbackResilienceStrategy.cs` | clean | Focused fallback strategy. |
| `src/Polly.Core/Hedging/HedgingResilienceStrategy.cs` | clean | Focused hedging strategy. |
| `src/Polly.Core/Retry/RetryResilienceStrategy.cs` | clean | Focused retry strategy; delay, jitter, predicates, callbacks and telemetry are intrinsic to the behavior. |
| `src/Polly.Core/Simmy/Behavior/ChaosBehaviorStrategy.cs` | clean | Focused chaos-behavior injection strategy. |
| `src/Polly.Core/Simmy/Outcomes/ChaosOutcomeStrategy.cs` | clean | Focused chaos-outcome injection strategy. |
| `src/Polly.Core/Utils/Pipeline/PipelineComponentFactory.cs` | clean; intentional composition seam | Internal factory for composing strategies, pipelines and lifecycle/disposal behavior. |
| `src/Polly/CircuitBreaker/AsyncCircuitBreakerPolicy.cs` | intentional legacy compatibility seam; maintenance watch | Cohesive adapter in the supported legacy policy layer. |
| `src/Polly/CircuitBreaker/CircuitBreakerPolicy.cs` | intentional legacy compatibility seam; maintenance watch | Synchronous counterpart to the legacy circuit-breaker adapter. |
| `src/Polly/Fallback/FallbackPolicy.cs` | intentional legacy compatibility seam; maintenance watch | Cohesive legacy fallback policy adapter. |
| `src/Polly/Retry/RetryPolicy.cs` | intentional legacy compatibility seam; maintenance watch | Cohesive legacy retry policy adapter. |

The coexistence of the modern `Polly.Core` resilience-strategy model and the supported legacy `Polly` policy model is real baseline maintenance pressure. It is intentional compatibility architecture, not evidence that each warned adapter independently violates ownership.

## Non-production peer class

The following current findings are not production architecture in this baseline:

- `bench/Polly.Core.Benchmarks/BridgeBenchmark.cs` — benchmark
- `src/Snippets/Docs/Chaos.Fault.cs` — documentation snippet
- `test/Polly.Core.Tests/CircuitBreaker/CircuitBreakerManualControlTests.cs` — test
- `test/Polly.Extensions.Tests/DependencyInjection/PollyServiceCollectionExtensionTests.cs` — test
- `test/Polly.Specs/Bulkhead/BulkheadAsyncSpecs.cs` — test

## Ground-truth summary

Polly is not broadly unhealthy. Most current dependency-pressure findings are intentional strategy, composition, or compatibility seams. The stateful circuit-breaker and hedging controllers are legitimate maintenance watches, and the legacy API layer creates known compatibility maintenance pressure, but the audit did not establish unrelated-responsibility sprawl in the warned production files.

## Related docs

- [Diagnosis Calibration Baselines](INDEX.md)
- [Architecture diagnosis plan](../../planning/architecture-diagnosis-and-policy-authoring.md)

## Notes

This baseline is tied to the exact source revision above. Re-adjudicate only if that corpus revision changes or source evidence demonstrates that one of these frozen judgments is wrong.
