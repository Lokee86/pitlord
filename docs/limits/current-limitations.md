# Current Limitations

Parent index: [Pitlord Limits](INDEX.md)

## Purpose

This document records current Pitlord limitations that materially affect use or extension.

## Overview

Pitlord is functional as an on-demand CLI, but several product and integration surfaces remain intentionally incomplete.

## Limitations

### Generalized architecture detectors

`pitlord scan` currently implements one generalized detector: language-neutral outgoing cross-file dependency pressure. Dependency knots, boundary-coupling/cohesion anomalies, and blast-radius hotspots remain to be implemented. The detector is calibrated against Arcana-style modular, entangled, hub-heavy, layered, and dense-subsystem topology fixtures and has now been exercised across Pitlord plus representative C#, Java, Kotlin, Svelte/Rust, GDScript, and Go/polyglot repositories.

Sixteen frozen real-world corpora now have machine-readable detector references and can be evaluated reproducibly with `pitlord calibrate`. The current detector remains not ready for blocking use: whole-repository medians can overstate pressure in large fine-grained repositories, tests and benchmarks currently share the production peer population, language-level aggregates such as C# partial types and Go packages are flattened to files, and central domain/entrypoint files can receive overly generic split recommendations. The detector therefore remains advisory while peer grouping, source-role handling, semantic granularity, absolute/statistical floors, and remediation wording are refined.

**Removal condition:** the initial deterministic generalized detector set is implemented, calibrated against Arcana synthetic topology families, and validated across representative real repositories with stable source-role peer groups and acceptable false-positive rates.

### Change-aware policy scope

Snapshot diff evaluates the full selected policy against both snapshots before comparing evidence. It does not yet restrict evaluation to changed policy scopes.

**Removal condition:** diff planning can safely identify affected rules and graph regions without missing introduced evidence.

### Documentation semantics

Pitlord can enforce required paths and content, but it does not understand Markdown index completeness, link resolution, document classes, or semantic implementation coverage. Those checks are owned by Demon Docs and the shared documentation checker.

**Removal condition:** none currently planned; this is an explicit tool boundary rather than a defect.

### Arcana query bounds and deadline policy

The Arcana process adapter honors caller cancellation through `exec.CommandContext`, but Pitlord does not yet define one product-wide default timeout. Callers using an unbounded context can therefore wait indefinitely for a stuck Arcana process.

Pitlord requests Arcana's current maximum of 10,000 relationships for each `neighbors` query and fails closed if Arcana reports truncation. The current `arcana.query.v1` neighbor operation has no continuation offset, so a single source node with more than 10,000 matching neighbors cannot yet be loaded completely even though ordinary large repositories can exceed the protocol's 1,000-item default safely.

**Removal condition:** CLI and Warlock invocation contracts adopt an explicit bounded deadline policy, and Arcana exposes continuation semantics for neighbor queries so Pitlord can load arbitrarily large bounded pages without truncation.

### Runtime integration

Pitlord has no independent daemon. Warlock may invoke it against changed snapshots, but central scheduling and UI compliance reporting are not yet integrated.

**Removal condition:** Warlock owns and ships the integration.

## Affected systems

- Snapshot diff and CI optimization.
- Warlock runtime and compliance UI.
- Cross-tool documentation governance.

## Status

All entries are current as of August 27, 2026.

## Related docs

- [Roadmap](../planning/roadmap.md)
- [Architecture](../ARCHITECTURE.md)
- [Shared adoption model](../../.standards/docs/adoption.md)

## Notes

The documentation-semantics boundary is deliberate: Pitlord enforces repository policy; Demon Docs and the shared checker own document structure and navigation.
