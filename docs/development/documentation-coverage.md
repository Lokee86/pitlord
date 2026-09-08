# Documentation Coverage

Parent index: [Development Documentation](INDEX.md)

## Purpose

This document maps Pitlord production ownership to canonical current documentation.

## Overview

Coverage includes packages, public command families, machine-readable contracts, and independent stateful flows. A mapping is valid only when the linked document actually explains the relevant behavior.

## Commands

| Surface | Implementation | Canonical current owner |
| --- | --- | --- |
| `check`, policy execution | `cmd/pitlord/check.go`, `internal/checker`, `internal/policy` | [Commands](../COMMANDS.md), [Policy](../POLICY.md), [Architecture](../ARCHITECTURE.md) |
| `validate`, embedded schemas | `cmd/pitlord/validate.go`, `internal/schema`, `internal/policy` | [Commands](../COMMANDS.md), [Policy](../POLICY.md) |
| `baseline` | `cmd/pitlord/baseline.go`, `internal/baseline` | [Commands](../COMMANDS.md), [Baselines and CI](../BASELINES-AND-CI.md) |
| `diff` | `cmd/pitlord/diff.go`, `internal/snapshotdiff` | [Commands](../COMMANDS.md), [Architecture](../ARCHITECTURE.md), [Baselines and CI](../BASELINES-AND-CI.md) |
| `inspect`, `analyze` | `cmd/pitlord/inspect.go`, `cmd/pitlord/analyze.go`, `internal/arcana`, `internal/report` | [Commands](../COMMANDS.md), [Architecture](../ARCHITECTURE.md) |
| `scan`, analyzer registry and policy-free architecture/semantic/native-analyzer execution | `cmd/pitlord/scan.go`, `internal/scan`, `internal/report` | [Commands](../COMMANDS.md), [Architecture](../ARCHITECTURE.md), [Current limitations](../limits/current-limitations.md), [Semantic calibration](calibration/semantic-lints.md) |
| `docs`, codemap coverage and changed-document guard | `cmd/pitlord/docs.go`, `internal/docguard`, `internal/arcana` | [Commands](../COMMANDS.md), [Architecture](../ARCHITECTURE.md), [Current limitations](../limits/current-limitations.md) |
| `calibrate`, frozen-corpus detector/analyzer evaluation and strict regression gating | `cmd/pitlord/calibrate.go`, `internal/calibration`, `internal/scan`, `internal/report` | [Commands](../COMMANDS.md), [Calibration harness](calibration/HARNESS.md), [Architecture](../ARCHITECTURE.md) |
| `generate`, Homunculus policy conversion | `cmd/pitlord/generate.go`, `internal/policy/homunculus.go` | [Commands](../COMMANDS.md), [Architecture](../ARCHITECTURE.md) |
| `verify-mutation` | `cmd/pitlord/verify_mutation.go`, `internal/mutation` | [Commands](../COMMANDS.md), [Architecture](../ARCHITECTURE.md) |
| `schema` | `cmd/pitlord/schema.go`, `internal/schema` | [Commands](../COMMANDS.md), [Policy](../POLICY.md) |

## Package ownership

| Package | Responsibility | Canonical current owner |
| --- | --- | --- |
| `internal/policy` | Rule model, normalization, validation, repository and graph evaluation | [Policy](../POLICY.md), [Architecture](../ARCHITECTURE.md) |
| `internal/arcana` | Arcana provider discovery, bounded protocol queries, and graph loading | [Architecture](../ARCHITECTURE.md), [Arcana process boundary](../ARCANA_PROCESS_BOUNDARY.md) |
| `internal/scan` | Analyzer registry and requirement gating; policy-free finding contract; graph-backed architecture detectors; language-neutral semantic capability validation plus `swallowed-error` and `unobserved-outcome`; graph-free Rust Clippy normalization; source-role projection; severity/scope normalization; deterministic ordering; and scan result assembly | [Architecture](../ARCHITECTURE.md), [Current limitations](../limits/current-limitations.md), [Semantic calibration](calibration/semantic-lints.md) |
| `internal/calibration` | Strict calibration reference loading; detector or analyzer/rule/language targeting; optional exact source-location matching; pinned revision and optional worktree-diff verification; labelled scoring, severity bounds, unlabelled-finding accounting, and opt-in fully-labelled mismatch semantics | [Calibration harness](calibration/HARNESS.md), [Architecture](../ARCHITECTURE.md) |
| `internal/docguard` | Demon Docs schema-1 codemap ingestion, Arcana code-file coverage, Git change loading, and changed-code-to-owning-document evaluation | [Architecture](../ARCHITECTURE.md), [README](../../README.md) |
| `internal/checker` | Check orchestration and report assembly | [Architecture](../ARCHITECTURE.md) |
| `internal/baseline` | Evidence fingerprints and suppression | [Baselines and CI](../BASELINES-AND-CI.md) |
| `internal/report` | Text, JSON, SARIF, and architecture reports | [README](../../README.md), [Baselines and CI](../BASELINES-AND-CI.md) |
| `internal/schema` | Embedded policy and baseline schemas | [Policy](../POLICY.md) |
| `internal/snapshot` | Snapshot resolution | [Architecture](../ARCHITECTURE.md) |
| `internal/snapshotdiff` | Introduced, resolved, and persistent evidence comparison | [Baselines and CI](../BASELINES-AND-CI.md) |
| `internal/mutation` | Exact Homunculus mutation-delta verification | [Architecture](../ARCHITECTURE.md) |

## Stateful flows

| Flow | Current owner |
| --- | --- |
| Policy include loading, normalization, and validation | [Policy](../POLICY.md) |
| Repository traversal and direct content/path evaluation | [Architecture](../ARCHITECTURE.md), [Policy](../POLICY.md) |
| Arcana executable discovery, snapshot resolution, and paginated loading | [Architecture](../ARCHITECTURE.md), [Arcana process boundary](../ARCANA_PROCESS_BOUNDARY.md) |
| Scan analyzer selection -> graph-requirement/capability validation -> architecture, semantic, and/or native analyzer execution -> normalized deterministic finding assembly | [Architecture](../ARCHITECTURE.md), [Current limitations](../limits/current-limitations.md), [Semantic calibration](calibration/semantic-lints.md) |
| Pinned-corpus revision/worktree verification -> exact built-in analyzer selection -> optional rule/language/location matching -> labelled scoring and optional fully-labelled regression gate | [Calibration harness](calibration/HARNESS.md), [Architecture](../ARCHITECTURE.md) |
| Demon Docs codemap export -> Arcana code-file inventory -> Git changed-file comparison -> documentation guard | [Architecture](../ARCHITECTURE.md), [README](../../README.md) |
| Evidence fingerprinting and baseline suppression | [Baselines and CI](../BASELINES-AND-CI.md) |
| Snapshot-to-snapshot evidence classification | [Baselines and CI](../BASELINES-AND-CI.md) |
| Mutation-manifest verification | [Architecture](../ARCHITECTURE.md) |

## Related docs

- [Behavioral contract matrix](behavioral-contract-matrix.md)
- [Documentation policy](../documentation-policy.md)

## Notes

Update this map whenever an independent stateful flow, command family, machine-readable contract, or package responsibility changes.
