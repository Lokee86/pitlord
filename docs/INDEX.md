# Pitlord Documentation

Parent index: [Pitlord](../README.md)

## Purpose

This index is the canonical entry point for Pitlord documentation.

## Overview

Pitlord documentation separates exact policy and CLI contracts, implemented architecture, CI and baseline operation, development coverage, plans, limitations, and repository documentation rules.

## Direct files

- [Command reference](COMMANDS.md) — Public CLI roles, usage, and the distinction between inspect/analyze/check/scan.
- [Architecture](ARCHITECTURE.md) — Implemented ownership, runtime flow, evaluation, snapshots, and mutation verification.
- [Arcana process boundary](ARCANA_PROCESS_BOUNDARY.md) — Process ownership, deadlines, protocol compatibility, failure behavior, and observability.
- [Architecture decisions](decisions/INDEX.md) — Consequential ownership, process, persistence, and compatibility decisions.
- [Maintainer map](MAINTAINER_MAP.md) — Routes common changes to canonical documentation and primary implementation boundaries.
- [Policy reference](POLICY.md) — Exact policy schema and rule semantics.
- [Baselines and CI](BASELINES-AND-CI.md) — Evidence fingerprints, baselines, SARIF, and snapshot gating.
- [Documentation policy](documentation-policy.md) — Pitlord-specific documentation ownership and required shapes.
- [Documentation procedure](documentation-procedure.md) — Required documentation workflow and verification.

## Direct folders

- [Development](development/INDEX.md) — Coverage, behavioral contracts, testing, and extension ownership.
- [Limits](limits/INDEX.md) — Current product limitations and incomplete behavior.
- [Planning](planning/INDEX.md) — Future and unresolved work.

## Related docs

- [Shared documentation standard](../../engineering-standards/docs/documentation-standard.md)
- [README](../README.md)

## Notes

The root README remains the product entry point. These documents own the detailed contracts.
