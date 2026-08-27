# Apache Maven Calibration Baseline

Parent index: [Diagnosis Calibration Baselines](INDEX.md)

## Purpose

This document freezes the manually adjudicated reference architecture state of the Apache Maven corpus used for Pitlord dependency-pressure calibration.

## Overview

- Corpus: `java-maven`
- Source revision: `b646a7a38e920f7590e7b25fb43b5756d8bf2b4d`
- Commit subject: `Add test for prefixed Maven elements (#10971)`
- Language: Java
- Baseline role: large multi-module real-world control

The baseline records corpus state only. It does not prescribe detector changes. Lexicon and Arcana were current for this revision; Grimoire source/document indexing did not complete within the preparation timeout, so adjudication used direct source inspection rather than assuming complete Grimoire coverage.

## Current Pitlord observation

A repository-wide dependency-pressure scan reports 20 high advisory findings. Eighteen belong to Maven API, compatibility, CLI, core, or implementation modules. Two belong to `impl/maven-testing` support code.

## Manually adjudicated production state

### Public API and session

| Scope | Expected interpretation | Notes |
| --- | --- | --- |
| `api/maven-api-core/src/main/java/org/apache/maven/api/Session.java` | intentional broad public facade; maintenance watch | Session exposes the build/session service surface and convenience operations across repositories, artifacts, projects and services. Breadth is deliberate API ownership. |
| `impl/maven-impl/src/main/java/org/apache/maven/impl/AbstractSession.java` | intentional broad session implementation; maintenance watch | Internal implementation of the broad session contract. Central integration seam rather than an isolated domain object. |

### Compatibility layer

| Scope | Expected interpretation | Notes |
| --- | --- | --- |
| `compat/maven-compat/src/main/java/org/apache/maven/plugin/internal/DefaultPluginManager.java` | compatibility maintenance debt | Legacy plugin-management compatibility surface. Cohesive within its role, but its existence preserves a parallel historical integration layer. |
| `compat/maven-compat/src/main/java/org/apache/maven/project/artifact/MavenMetadataSource.java` | compatibility maintenance debt | Legacy project/artifact metadata bridge. |
| `compat/maven-compat/src/main/java/org/apache/maven/repository/legacy/LegacyRepositorySystem.java` | compatibility maintenance debt | Explicit legacy repository-system adapter. |
| `compat/maven-embedder/src/main/java/org/apache/maven/cli/MavenCli.java` | intentional broad CLI/composition seam; maintenance watch | Legacy/embedder CLI entry point with naturally broad bootstrap and execution ownership. |
| `compat/maven-embedder/src/main/java/org/apache/maven/cli/internal/BootstrapCoreExtensionManager.java` | compatibility/bootstrap maintenance watch | Core-extension bootstrap in the compatibility CLI stack. |
| `compat/maven-model-builder/src/main/java/org/apache/maven/model/building/DefaultModelBuilder.java` | compatibility structural maintenance debt | ~1,099-line legacy effective-model pipeline spanning model reading, inheritance, profiles, interpolation, validation, dependency/plugin management, and parent resolution. The work is related, but materially concentrated inside the supported compatibility layer. |
| `compat/maven-model-builder/src/main/java/org/apache/maven/model/building/DefaultModelBuilderFactory.java` | clean; intentional factory/composition seam | Factory assembling model-builder collaborators. |

The `compat/` modules are frozen as real compatibility maintenance debt at the architectural layer level. Individual classes remain cohesive within that supported legacy boundary.

### Current CLI and core runtime

| Scope | Expected interpretation | Notes |
| --- | --- | --- |
| `impl/maven-cli/src/main/java/org/apache/maven/cling/extensions/BootstrapCoreExtensionManager.java` | intentional bootstrap seam; maintenance watch | Current CLI core-extension bootstrap coordinator. |
| `impl/maven-cli/src/main/java/org/apache/maven/cling/invoker/LookupInvoker.java` | clean; intentional invocation seam | Dependency/service lookup stage in CLI invocation. |
| `impl/maven-cli/src/main/java/org/apache/maven/cling/invoker/mvn/MavenInvoker.java` | intentional broad invocation orchestrator; maintenance watch | Maven-specific CLI invocation lifecycle. |
| `impl/maven-core/src/main/java/org/apache/maven/DefaultMaven.java` | intentional core build orchestrator; maintenance watch | Coordinates Maven execution at the application-core boundary. |
| `impl/maven-core/src/main/java/org/apache/maven/lifecycle/internal/DefaultLifecycleExecutionPlanCalculator.java` | intentional algorithmic seam; maintenance watch | Lifecycle-plan calculation integrates project, plugin and execution-plan concerns intrinsic to planning. |
| `impl/maven-core/src/main/java/org/apache/maven/lifecycle/internal/concurrent/BuildPlanExecutor.java` | intentional concurrency/orchestration seam; maintenance watch | Executes build plans across concurrent project scheduling/runtime concerns. |
| `impl/maven-core/src/main/java/org/apache/maven/plugin/internal/DefaultMavenPluginManager.java` | broad subsystem coordinator; maintenance watch | Plugin descriptor loading, dependency resolution, realms, component discovery and mojo configuration converge in one plugin-management subsystem. This is substantial pressure, but inspected responsibilities remain plugin lifecycle ownership. |
| `impl/maven-core/src/main/java/org/apache/maven/project/DefaultProjectBuilder.java` | real structural maintenance pressure | ~1,164-line project-building coordinator spanning model construction, repositories/session wiring, dependency/artifact setup, lifecycle/plugin setup, validation, and reactor concerns. Cohesive around project construction, but materially concentrated. |
| `impl/maven-impl/src/main/java/org/apache/maven/impl/model/DefaultModelBuilder.java` | real structural maintenance pressure | ~2,457-line current effective-model pipeline spanning source reading, parent/mixin resolution, profiles, interpolation, normalization, validation, caching, and result assembly. This is a genuine decomposition pressure point even though the responsibilities share model-building ownership. |

## Testing-support peer class

The following are production-source files of the `maven-testing` support module rather than shipped Maven core architecture and are treated as testing-support peers:

- `impl/maven-testing/src/main/java/org/apache/maven/testing/plugin/MojoExtension.java`
- `impl/maven-testing/src/main/java/org/apache/maven/testing/plugin/stubs/RepositorySystemSupplier.java`

## Ground-truth summary

Maven is a deliberately broad multi-module build system with many central orchestration seams. The audit does not treat all 20 high findings as independent architectural failures. Frozen debt includes the supported compatibility architecture under `compat/` plus material concentration in both model-builder implementations and `DefaultProjectBuilder`. Other CLI, plugin, lifecycle, session, and concurrency coordinators remain broad maintenance watches whose dependency breadth largely matches subsystem ownership.

## Related docs

- [Diagnosis Calibration Baselines](INDEX.md)
- [Architecture diagnosis plan](../../planning/architecture-diagnosis-and-policy-authoring.md)

## Notes

This baseline is tied to the exact source revision above. Re-adjudicate only if that corpus revision changes or source evidence demonstrates that one of these frozen judgments is wrong.
