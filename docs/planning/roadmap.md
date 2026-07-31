# Roadmap

Parent index: [Pitlord Planning](INDEX.md)

## Purpose

This document records planned Pitlord work and unresolved sequencing.

## Overview

Pitlord remains independently usable while becoming the architecture-policy and repository-governance engine within Warlock.

## Current status

Active planning.

## Expected ownership

- Pitlord owns rule semantics, evaluation, evidence identity, baselines, and machine-readable diagnostics.
- Arcana owns graph storage and traversal.
- Lexicon owns language facts.
- Demon Docs and the shared checker own Markdown structure, navigation, and documentation coverage checks.
- Warlock owns scheduling, enrolled-repository orchestration, and compliance presentation.

## Planned work

1. Add reusable policy packs with explicit versioning and compatibility.
2. Integrate shared documentation-policy results into Warlock repository health.
3. Add change-aware evaluation planning without weakening full-policy correctness.
4. Expand policy analysis and authoring assistance using Arcana architecture summaries.
5. Improve cross-snapshot and baseline migration diagnostics.

## Acceptance criteria

Each planned feature must preserve deterministic evidence, explicit ownership, repository-only execution where applicable, schema validation, and independent CLI use.

## Open decisions

- How reusable policy-pack versions are resolved and pinned by repositories.
- Whether Warlock invokes Pitlord directly or through a shared daemon job protocol.
- Which change-planning evidence is sufficient to safely skip unaffected graph regions.

## Implemented references

- `require_content` is implemented and documented in [Policy reference](../POLICY.md).
- Shared documentation governance is defined in [Engineering Standards](../../../engineering-standards/docs/INDEX.md).

## Related docs

- [Architecture](../ARCHITECTURE.md)
- [Current limitations](../limits/current-limitations.md)

## Notes

This roadmap is not a statement of shipped behavior.
