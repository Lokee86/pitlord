# Documentation Policy

Parent index: [Pitlord Documentation](INDEX.md)

## Purpose

This document defines Pitlord's repository-local specialization of the shared documentation standard.

## Overview

Pitlord follows the shared CLI-product profile with the `stateful` capability. Current behavior is separated into policy reference, architecture, CI/operations, development, limits, and planning documents.

## Documentation ownership

| Information | Canonical owner |
| --- | --- |
| Product purpose, installation, quick start, command overview | `README.md` |
| Exact rule fields, defaults, validation, and semantics | `docs/POLICY.md` |
| Ownership, data flow, evaluation, snapshots, and process behavior | `docs/ARCHITECTURE.md` |
| Baselines, fingerprints, SARIF, and CI workflows | `docs/BASELINES-AND-CI.md` |
| Package, command, and stateful-flow coverage | `docs/development/documentation-coverage.md` |
| Critical invariants and protecting tests | `docs/development/behavioral-contract-matrix.md` |
| Human-adjudicated architecture and semantic-analyzer calibration evidence | `docs/development/calibration/` |
| Future work | `docs/planning/roadmap.md` |
| Current gaps and defects | `docs/limits/current-limitations.md` |
| Agent editing and verification rules | `AGENTS.md` |

## Rules

- Reference documentation must match the embedded schema, normalization, validation, diagnostics, and tests.
- Architecture documentation describes implemented ownership only.
- Repository-only rule behavior must identify that Arcana is not required.
- Stateful artifacts such as baselines, snapshots, fingerprints, and mutation manifests require lifecycle and compatibility documentation.
- Every production package, public command family, independent stateful flow, and machine-readable contract has a current documentation owner.
- Current behavior must not live only in planning.
- Existing limitations remain explicit until their removal condition is met.
- Every documentation folder contains `INDEX.md` except `stubs/`.
- Normal documents contain Purpose, Overview, Related docs, and Notes.

## Code maps and tests

Implementation-facing architecture and development documents include code maps and focused tests. A code map does not replace prose explaining responsibility, lifecycle, state, failure behavior, and non-ownership boundaries.

## Related docs

- [Shared documentation standard](../.standards/docs/documentation-standard.md)
- [Documentation procedure](documentation-procedure.md)
- [Documentation coverage](development/documentation-coverage.md)

## Notes

The shared standard is authoritative for cross-repository requirements. This policy may be stricter but must not silently weaken it.
