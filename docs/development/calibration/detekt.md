# detekt Calibration Baseline

Parent index: [Diagnosis Calibration Baselines](INDEX.md)

## Purpose

This document freezes the manually adjudicated reference architecture state of the detekt corpus used for Pitlord dependency-pressure calibration.

## Overview

- Corpus: `kotlin-detekt`
- Source revision: `f9e1d5cc239ab740ce499b1edb36b872012648e2`
- Commit subject: `Remove unused version property (#9546)`
- Primary language: Kotlin
- Baseline role: external real-world control

The baseline records corpus state only. It does not prescribe detector changes.

## Repository architecture

The repository defines explicit module ownership for its major surfaces:

- `detekt-api/` — public extension API
- `detekt-core/` — core analysis engine
- `detekt-cli/` — command-line interface
- `detekt-gradle-plugin/` — Gradle integration
- `detekt-rules-*/` — rule implementations by category
- `detekt-report-*/` — report implementations
- `detekt-test/` — testing utilities

These module boundaries are the primary architectural units for this baseline.

## Current Pitlord observation

A repository-wide dependency-pressure scan reports eight advisory warnings: five production files and three test files. All are warnings; none are high or critical.

The production findings are:

- `detekt-core/src/main/kotlin/dev/detekt/core/config/validation/ConfigValidation.kt`
- `detekt-rules-style/src/main/kotlin/dev/detekt/rules/style/UnusedPrivateFunction.kt`
- `detekt-rules-style/src/main/kotlin/dev/detekt/rules/style/UnusedPrivateProperty.kt`
- `detekt-rules-style/src/main/kotlin/dev/detekt/rules/style/UnusedVariable.kt`
- `detekt-rules-style/src/main/kotlin/dev/detekt/rules/style/VarCouldBeVal.kt`

Module-scoped scans currently report:

- `detekt-core/src/main/kotlin/` — one warning on `ConfigValidation.kt`
- `detekt-rules-style/src/main/kotlin/` — no findings

The difference between repository-wide and module-scoped output is preserved as an observation, not treated here as ground truth.

## Manually adjudicated production state

### Core configuration validation

| Scope | Expected interpretation | Notes |
| --- | --- | --- |
| `detekt-core/.../config/validation/ConfigValidation.kt` | clean; focused validation orchestration | A 99-line orchestration file for configuration validation. It loads configured validators, validates against the baseline config, renders notifications, dispatches YAML/config-specific validation, and promotes warnings when configured. Its dependencies remain within the core configuration-validation/reporting/tooling path. |

The surrounding `config/validation` package is already decomposed into focused validator implementations including default-property, deprecated-property, invalid-property, and missing-rule validators. The inspected `ConfigValidation.kt` file coordinates those components rather than accumulating their implementation responsibilities.

### Style rules

| Scope | Expected interpretation | Notes |
| --- | --- | --- |
| `UnusedPrivateFunction.kt` | clean; cohesive rule implementation | Implements one rule and its visitor for collecting private function declarations/references and resolving Kotlin operator/delegate cases. |
| `UnusedPrivateProperty.kt` | clean; cohesive rule implementation | Implements one rule and visitor for private properties and constructor parameters, including analysis-API symbol resolution needed to determine usage. |
| `UnusedVariable.kt` | clean; cohesive rule implementation | Implements one rule and visitor for local-variable declaration/reference tracking. |
| `VarCouldBeVal.kt` | clean; cohesive rule implementation | Implements one rule and assignment visitor for determining whether mutable declarations are actually reassigned, including object-literal escape cases. |

The four warned style files are 131–221 lines and remain centered on one rule each. Their Kotlin PSI and Analysis API dependency breadth supports the semantics of those individual rules rather than unrelated architectural ownership.

## Non-production peer class

The three repository-wide findings under test source sets are test architecture, not production dependency-pressure findings in this baseline. Their inspected breadth is consistent with tests spanning the production surfaces they exercise.

Current observed test findings are:

- `detekt-core/src/test/kotlin/dev/detekt/core/AnalyzerSpec.kt` — analyzer behavior across configuration, rules, findings, and analysis modes
- `detekt-core/src/test/kotlin/dev/detekt/core/RuleDescriptorKtTest.kt` — rule descriptor/configuration behavior across API rule types and modes
- `detekt-rules-ktlint-wrapper/src/test/kotlin/dev/detekt/rules/ktlintwrapper/RulesWhichCantBeCorrectedSpec.kt` — focused cross-rule regression coverage for non-correctable ktlint-wrapper rules

## Ground-truth summary

No currently warned production file demonstrates unrelated architectural responsibility sprawl in the inspected revision.

`ConfigValidation.kt` remains a focused coordinator inside an already decomposed validation package. The four warned style-rule files each implement one semantically cohesive static-analysis rule. The three remaining warnings are test-source findings and belong to a separate peer class.

This baseline therefore records the eight current warnings as corpus observations while classifying the five inspected production targets as clean within their stated module responsibilities.

## Related docs

- [Diagnosis Calibration Baselines](INDEX.md)
- [jsoup](jsoup.md)
- [Dapper](dapper.md)
- [Space Rocks Clean Control v1](space-rocks-clean-v1.md)
- [Architecture diagnosis plan](../../planning/architecture-diagnosis-and-policy-authoring.md)

## Notes

This baseline is tied to the exact source revision above. Re-adjudicate only if that corpus revision changes or source evidence demonstrates that one of these frozen judgments is wrong.