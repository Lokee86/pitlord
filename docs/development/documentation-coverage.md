# Documentation Coverage

Parent index: [Development Documentation](INDEX.md)

## Purpose

This document maps Pitlord production ownership to canonical current documentation.

## Overview

Coverage includes packages, public command families, machine-readable contracts, and independent stateful flows. A mapping is valid only when the linked document actually explains the relevant behavior.

## Commands

| Surface | Implementation | Canonical current owner |
| --- | --- | --- |
| `check`, policy execution | `cmd/pitlord/check.go`, `internal/checker`, `internal/policy` | [README](../../README.md), [Policy](../POLICY.md), [Architecture](../ARCHITECTURE.md) |
| `validate`, embedded schemas | `cmd/pitlord/validate.go`, `internal/schema`, `internal/policy` | [Policy](../POLICY.md) |
| `baseline` | `cmd/pitlord/baseline.go`, `internal/baseline` | [Baselines and CI](../BASELINES-AND-CI.md) |
| `diff` | `cmd/pitlord/diff.go`, `internal/snapshotdiff` | [Architecture](../ARCHITECTURE.md), [Baselines and CI](../BASELINES-AND-CI.md) |
| `inspect`, `analyze` | `cmd/pitlord/inspect.go`, `cmd/pitlord/analyze.go`, `internal/arcana`, `internal/report` | [README](../../README.md), [Architecture](../ARCHITECTURE.md) |
| `scan`, policy-free generalized guard foundation | `cmd/pitlord/scan.go`, `internal/scan`, `internal/report` | [README](../../README.md), [Architecture](../ARCHITECTURE.md), [Current limitations](../limits/current-limitations.md) |
| `docs`, codemap coverage and changed-document guard | `cmd/pitlord/docs.go`, `internal/docguard`, `internal/arcana` | [README](../../README.md), [Architecture](../ARCHITECTURE.md), [Current limitations](../limits/current-limitations.md) |
| `calibrate`, frozen-corpus detector evaluation | `cmd/pitlord/calibrate.go`, `internal/calibration`, `internal/scan`, `internal/report` | [README](../../README.md), [Calibration harness](calibration/HARNESS.md), [Architecture](../ARCHITECTURE.md) |
| `generate`, Homunculus policy conversion | `cmd/pitlord/generate.go`, `internal/policy/homunculus.go` | [README](../../README.md), [Architecture](../ARCHITECTURE.md) |
| `verify-mutation` | `cmd/pitlord/verify_mutation.go`, `internal/mutation` | [README](../../README.md), [Architecture](../ARCHITECTURE.md) |
| `schema` | `cmd/pitlord/schema.go`, `internal/schema` | [Policy](../POLICY.md) |

## Package ownership

| Package | Responsibility | Canonical current owner |
| --- | --- | --- |
| `internal/policy` | Rule model, normalization, validation, repository and graph evaluation | [Policy](../POLICY.md), [Architecture](../ARCHITECTURE.md) |
| `internal/arcana` | Arcana provider discovery, bounded protocol queries, and graph loading | [Architecture](../ARCHITECTURE.md), [Arcana process boundary](../ARCANA_PROCESS_BOUNDARY.md) |
| `internal/scan` | Policy-free generalized finding contract, source-role peer classification, dependency-pressure boundary/role judgment, dependency-knot static-relation SCC projection, dependency-depth dominant acyclic behavioral-corridor judgment, unstable-dependency-direction region stability-gradient judgment, boundary-bypass file/region intermediary judgment, symbol-intermediary-bypass repeated-call-pipeline judgment, cross-file-intermediary-bypass same-region seam judgment, semantic architectural-region projection, boundary/cohesion spread and internal-support judgment, hub/bottleneck direct-centrality candidate generation plus many-to-many behavioral-flow judgment, impact/blast-radius direct-exposure plus independent-transitive-impact judgment, severity judgment, scope normalization, deterministic ordering, and scan result assembly | [Architecture](../ARCHITECTURE.md), [Current limitations](../limits/current-limitations.md) |
| `internal/calibration` | Strict calibration reference loading, pinned revision and optional worktree-diff verification, labelled detector scoring, severity bounds, whole-repository expectation matching, and unlabelled-finding accounting | [Calibration harness](calibration/HARNESS.md), [Architecture](../ARCHITECTURE.md) |
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
| Policy-free generalized scan peer classification, dependency-pressure judgment, dependency-knot static-relation/semantic-region SCC judgment, dependency-depth dominant acyclic behavioral-corridor judgment, unstable-dependency-direction region stability-gradient judgment, boundary-bypass file/region intermediary judgment, symbol-intermediary-bypass repeated-call-pipeline judgment, cross-file-intermediary-bypass same-region seam judgment, boundary/cohesion semantic-region judgment, hub/bottleneck behavioral-waist judgment, impact/blast-radius direct-exposure and independent-transitive-impact judgment, and result assembly | [Architecture](../ARCHITECTURE.md), [Current limitations](../limits/current-limitations.md) |
| Pinned-corpus calibration revision/worktree-state verification and labelled detector scoring | [Calibration harness](calibration/HARNESS.md), [Architecture](../ARCHITECTURE.md) |
| Demon Docs codemap export -> Arcana code-file inventory -> Git changed-file comparison -> documentation guard | [Architecture](../ARCHITECTURE.md), [README](../../README.md) |
| Evidence fingerprinting and baseline suppression | [Baselines and CI](../BASELINES-AND-CI.md) |
| Snapshot-to-snapshot evidence classification | [Baselines and CI](../BASELINES-AND-CI.md) |
| Mutation-manifest verification | [Architecture](../ARCHITECTURE.md) |

## Related docs

- [Behavioral contract matrix](behavioral-contract-matrix.md)
- [Documentation policy](../documentation-policy.md)

## Notes

Update this map whenever an independent stateful flow, command family, machine-readable contract, or package responsibility changes.
