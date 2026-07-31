# Pitlord Maintainer Map

Parent index: [Pitlord Documentation](INDEX.md)

## Purpose

This document routes maintainers to the canonical documentation and primary implementation boundary for common Pitlord changes.

## Overview

Use this map when the owner of a policy, analysis, baseline, mutation, or reporting change is unclear. It is not a package inventory and does not replace the detailed architecture and policy references.

## Change-area routing

| Change area | Canonical documentation | Primary implementation boundary | Verification |
| --- | --- | --- | --- |
| CLI commands and check orchestration | [Architecture](ARCHITECTURE.md), [Policy reference](POLICY.md) | `cmd/pitlord/`, `internal/checker/` | CLI and checker tests |
| Policy loading, validation, matching, ownership, and dependency rules | [Policy reference](POLICY.md), [Architecture](ARCHITECTURE.md) | `internal/policy/`, `internal/schema/` | Policy and schema tests |
| Arcana graph loading and architecture analysis | [Architecture](ARCHITECTURE.md) | `internal/arcana/`, `internal/policy/analysis*.go` | Arcana and analysis tests |
| Baseline fingerprints and accepted findings | [Baselines and CI](BASELINES-AND-CI.md) | `internal/baseline/` | Baseline and fingerprint tests |
| Snapshot selection and before/after comparison | [Baselines and CI](BASELINES-AND-CI.md), [Architecture](ARCHITECTURE.md) | `internal/snapshot/`, `internal/snapshotdiff/` | Snapshot and diff tests |
| Mutation plans and post-mutation verification | [Architecture](ARCHITECTURE.md) | `internal/mutation/` | Mutation load and verification tests |
| Text, architecture, area-analysis, and SARIF reports | [Baselines and CI](BASELINES-AND-CI.md), [Policy reference](POLICY.md) | `internal/report/` | Report and SARIF tests |
| CI, documentation coverage, and release gates | [Baselines and CI](BASELINES-AND-CI.md), [Documentation coverage](development/documentation-coverage.md) | `.github/workflows/`, `docs/development/` | Full Go, policy, and documentation checks |

## Boundaries

- `internal/policy/` owns rule semantics; `internal/report/` only renders evaluated findings.
- `internal/arcana/` supplies graph evidence but does not own policy decisions.
- Baselines accept known evidence fingerprints; they do not redefine rule behavior.
- Maintainer navigation belongs here; focused package and flow details remain in `ARCHITECTURE.md` and `POLICY.md`.

## Related docs

- [Architecture](ARCHITECTURE.md)
- [Policy reference](POLICY.md)
- [Documentation coverage](development/documentation-coverage.md)

## Notes

Update this map when command families, policy owners, analysis providers, snapshot flows, or report surfaces move.
