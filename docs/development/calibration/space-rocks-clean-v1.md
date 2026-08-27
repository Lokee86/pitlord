# Space Rocks Clean Control v1

Parent index: [Diagnosis Calibration Baselines](INDEX.md)

## Purpose

This document freezes the manually adjudicated reference architecture baseline for the Space Rocks clean control used to tune Pitlord's generalized dependency-pressure diagnosis.

## Overview

Space Rocks remains a qualified clean architecture control, but `clean` does not mean that every file has low dependency breadth or that no maintenance pressure exists. The correct baseline contains one genuine maintenance-pressure hotspot, a small number of intentional broad seams worth watching, and many graph outliers that should not be interpreted as architectural violations.

Pitlord must be calibrated against these judgments rather than against the current number or severity of findings.

## Frozen identity

- Corpus: `space-rocks-clean-v1`
- Source checkout: `workspace/corpus/space-rocks-clean-v1`
- Pinned commit: `361dcdb4054725c8f1bb58727ff402fcfedf6d18`
- Commit subject: `feat(gameplay): add deathmatch spawn profiles`
- Reference adjudication frozen: 2026-08-27
- Primary detector under calibration: `dependency-pressure`

Generated Grimoire, Lexicon, Arcana, Warlock, build, and test state does not redefine this source identity.

## Adjudication vocabulary

| Classification | Meaning for tuning |
| --- | --- |
| `pressure` | Pitlord should emit an advisory because the structure creates genuine maintenance or comprehension pressure. |
| `watch` | Breadth is intentional, but another signal such as size or behavioral growth makes continued expansion risky. Warning severity from dependency count alone is too strong. |
| `intentional-breadth` | High dependency breadth is intrinsic to the component's architectural role and should not itself produce a dependency-pressure warning. |
| `clean` | The inspected production structure is cohesive and should not produce this detector's warning. |
| `peer-class-exception` | The node belongs to tests, benchmarks, generated code, tooling, or another role that must not be judged against ordinary production peers. |
| `granularity-false-positive` | File-level graph breadth exaggerates architectural breadth because the language or repository organizes one architectural unit across many files. |

## Service-level ground truth

| Region | Expected interpretation |
| --- | --- |
| Client | One genuine pressure hotspot; a small number of intentional broad seams worth watching; remaining dependency-count outliers should be suppressed or compared within a more appropriate role. |
| Game server | Clean for this detector. Current file-level warnings largely confuse Go source-file decomposition and intentional orchestration with architectural boundary breadth. |
| API server | Clean for this detector. Current warnings are small-service statistical false positives. |
| Player-data | Clean for production code. Test-only outliers are a peer-class issue. |
| Diagnostic aggregator | Clean for this detector. Current warnings describe cohesive service/storage components in a small peer population. |

## Frozen labelled examples

### Client

| Path | Classification | Expected detector behavior | Rationale |
| --- | --- | --- | --- |
| `client/scripts/networking/client_connection_service.gd` | `pressure` | Advisory warning/high is justified; remediation must preserve its public networking-facade ownership. | At 597 lines it exceeds the repository's own actively-changing split-candidate guideline and combines connection lifecycle, public outbound facade, websocket auth state, realtime-session composition, inbound/tooling relay, room-operation tracing, telemetry, recovery, observability, and transport metrics. Much of the breadth is intentional, but the file is a genuine maintenance gravity well. |
| `client/scripts/shell/app_entry.gd` | `watch` | No dependency-pressure warning from fan-out alone; at most informational/watch-level pressure when size or behavioral growth corroborates it. | Canonical documentation defines this file as the Godot composition root. Broad construction and wiring dependencies are its job. At roughly 384 lines it is large enough that new behavior should remain constrained. |
| `client/scripts/shell/gameplay_menu_flow.gd` | `watch` | Do not warn from dependency count alone; size/responsibility detectors may watch it. | The file remains cohesive around pause, game-over, spectate, and gameplay-menu presentation, but its size is above the repository's review range. |
| `client/scripts/gameplay/runtime/gameplay_flow_composer.gd` | `intentional-breadth` | Suppress dependency-pressure warning unless independent evidence shows responsibility drift. | It is explicitly a composer that constructs focused gameplay flows and therefore legitimately touches many collaborators. |
| `client/scripts/gameplay/gameplay_composition.gd` | `intentional-breadth` | Suppress dependency-pressure warning unless independent evidence shows responsibility drift. | Composition and high-level routing are the file's primary responsibility. |
| `client/scripts/protocol/realtime/realtime_router.gd` | `intentional-breadth` | Suppress dependency-pressure warning from fan-out alone. | Packet-family/lane routing necessarily references the focused lane states, appliers, trackers, assemblers, and recovery machinery it coordinates. |
| `client/scripts/session/room_session_controller.gd` | `clean` | No dependency-pressure warning expected. | Cohesive room/lobby session seam with bounded routing and state ownership. |
| `client/scripts/world/projectile_sync.gd` | `clean` | No dependency-pressure warning expected. | Cohesive projectile presentation/synchronization concern despite several supporting dependencies. |
| `client/legacy/player_render/player_sync.gd` | `clean` | No dependency-pressure warning expected. | Cohesive player presentation synchronization split into focused lifecycle/interpolation/target/presentation collaborators. |
| Client devtools/runtime-scenario orchestration | `peer-class-exception` | Compare against tooling/orchestration peers; do not use ordinary production-file medians. | These files intentionally coordinate broad test/devtools surfaces and should not teach production dependency thresholds. |
| Client tests | `peer-class-exception` | Exclude from ordinary production dependency-pressure peers. | Integration and behavioral tests legitimately touch many production files. |

### Game server

| Path | Classification | Expected detector behavior | Rationale |
| --- | --- | --- | --- |
| `services/game-server/internal/game/combat.go` | `granularity-false-positive` | No dependency-pressure warning expected from current file fan-out. | The Go `game` package is intentionally decomposed across many sibling files. Calls into combat damage requests/application, runtime entities, scoring/death helpers, constants, and collision helpers inflate file-to-file fan-out without representing equivalent architectural boundaries. The file remains cohesive around collision/combat consequences. |
| `services/game-server/cmd/game-server/main.go` | `intentional-breadth` | No dependency-pressure warning from fan-out alone. | Process composition root: constructs co-hosted services, HTTP routes, rooms, tooling, player-data, reporting, auth, diagnostics, lifecycle, and shutdown. Reusable simulation is kept out of it as required. |
| `services/game-server/internal/game/game.go` | `intentional-breadth` | No dependency-pressure warning from fan-out alone. | Central `Game` aggregate and constructor; broad dependencies are intrinsic to assembling the authoritative simulation aggregate. |
| `services/game-server/internal/game/session.go` | `clean` | No dependency-pressure warning expected. | Cohesive player-session and respawn behavior. |
| `services/game-server/internal/game/simulation.go` | `clean` | No dependency-pressure warning expected. | Small simulation orchestrator; breadth reflects sequencing of focused simulation steps rather than responsibility sprawl. |
| `services/game-server/internal/game/simulation_radial_effects.go` | `clean` | No dependency-pressure warning expected. | Cohesive radial-effect simulation and damage application. |
| `services/game-server/internal/protocol/realtime/active.go` | `clean` | No dependency-pressure warning expected from current breadth. | Cohesive active realtime candidate/result construction, scheduling, encoding, and metrics pipeline. |
| Game-server tests | `peer-class-exception` | Exclude from ordinary production peers. | Cross-cutting tests legitimately exercise multiple collaborators and files. |

### API server

| Path/region | Classification | Expected detector behavior | Rationale |
| --- | --- | --- | --- |
| `app/controllers/api/auth/discord_controller.rb` | `clean` | No dependency-pressure warning expected. | Small controller coherently orchestrates one OAuth flow. |
| `app/controllers/api/auth/discord_login_sessions_controller.rb` | `clean` | No dependency-pressure warning expected. | Small controller owns one login-session exchange flow. |
| `app/controllers/api/auth/sessions_controller.rb` | `clean` | No dependency-pressure warning expected. | Small conventional session controller. |
| `app/lib/observability/worker_runtime.rb` | `clean` | No dependency-pressure warning expected. | Cohesive observability worker lifecycle. |
| API tests | `peer-class-exception` | Exclude from ordinary production peers. | Test dependency breadth is not production architectural breadth. |

### Player-data

Production code inspected under the service scope is `clean` for dependency pressure. The observed `configured_runtime_test.go` outlier is a `peer-class-exception`, not baseline production debt.

### Diagnostic aggregator

| Path | Classification | Expected detector behavior | Rationale |
| --- | --- | --- | --- |
| `internal/diagnosticreports/service.go` | `clean` | No dependency-pressure warning expected. | Cohesive report validation, redaction, creation, storage, retrieval, and observability service. |
| `internal/storage/jsonlstore/report_store.go` | `clean` | No dependency-pressure warning expected. | Cohesive JSONL report-store lifecycle. |
| `internal/storage/jsonlstore/report_retention.go` | `clean` | No dependency-pressure warning expected. | Cohesive retention/archive rewrite implementation. |
| Diagnostic tests/e2e tests | `peer-class-exception` | Exclude from ordinary production peers. | Cross-boundary test breadth is expected. |

## Known baseline maintenance pressure

`client/scripts/networking/client_connection_service.gd` is the one confirmed pre-existing maintenance-pressure hotspot discovered during this audit. It remains an intentional networking facade and is not an architecture-policy violation, but further unrelated responsibility should not accumulate there. Future dedicated responsibility/size/boundary detectors may strengthen or refine this judgment.

This condition is part of the frozen ground truth. Do not modify the clean corpus merely to make the detector return zero findings.

## Calibration consequences

The current file-level `3x median` dependency-pressure model is insufficient as a generalized architectural judgment because it conflates several distinct conditions:

1. production code versus tests, benchmarks, generated code, and tooling;
2. source files versus package/module/community architectural units;
3. composition roots, routers, orchestrators, and facades versus ordinary implementation modules;
4. small-service statistical outliers versus meaningful absolute complexity; and
5. dependency count versus dependency diversity across architectural regions.

Future tuning should preserve the true positive at `ClientConnectionService` without producing warnings for healthy composition roots, same-package Go decomposition, small cohesive controllers, or cross-cutting tests.

## Scoring expectations

For this baseline, tuning quality should be measured by labelled judgment rather than raw finding count:

- `pressure`: should be found with an appropriate advisory severity and role-preserving remediation;
- `watch`: must not be escalated solely because of raw fan-out;
- `intentional-breadth`, `clean`, and `granularity-false-positive`: dependency-pressure warning is a false positive unless new independent evidence changes the judgment;
- `peer-class-exception`: must be excluded or compared within an appropriate peer population;
- remediation that tells an intentional composition root/router/facade simply to `split` is considered incorrect even if the graph outlier itself is real.

## Related docs

- [Calibration baselines](INDEX.md)
- [Architecture diagnosis plan](../../planning/architecture-diagnosis-and-policy-authoring.md)
- [Current limitations](../../limits/current-limitations.md)
- [Behavioral contract matrix](../behavioral-contract-matrix.md)

## Notes

This is manually adjudicated reference tuning ground truth for the pinned source revision. A later Space Rocks commit requires a new baseline identity or an explicit requalification; it must not silently replace this one.