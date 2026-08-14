# Current Limitations

Parent index: [Pitlord Limits](INDEX.md)

## Purpose

This document records current Pitlord limitations that materially affect use or extension.

## Overview

Pitlord is functional as an on-demand CLI, but several product and integration surfaces remain intentionally incomplete.

## Limitations

### Change-aware policy scope

Snapshot diff evaluates the full selected policy against both snapshots before comparing evidence. It does not yet restrict evaluation to changed policy scopes.

**Removal condition:** diff planning can safely identify affected rules and graph regions without missing introduced evidence.

### Documentation semantics

Pitlord can enforce required paths and content, but it does not understand Markdown index completeness, link resolution, document classes, or semantic implementation coverage. Those checks are owned by Demon Docs and the shared documentation checker.

**Removal condition:** none currently planned; this is an explicit tool boundary rather than a defect.

### Arcana deadline policy

The Arcana process adapter honors caller cancellation through `exec.CommandContext`, but Pitlord does not yet define one product-wide default timeout. Callers using an unbounded context can therefore wait indefinitely for a stuck Arcana process.

**Removal condition:** CLI and Warlock invocation contracts adopt an explicit bounded deadline policy with repository-size evidence and focused timeout tests.

### Runtime integration

Pitlord has no independent daemon. Warlock may invoke it against changed snapshots, but central scheduling and UI compliance reporting are not yet integrated.

**Removal condition:** Warlock owns and ships the integration.

## Affected systems

- Snapshot diff and CI optimization.
- Warlock runtime and compliance UI.
- Cross-tool documentation governance.

## Status

All entries are current as of July 30, 2026.

## Related docs

- [Roadmap](../planning/roadmap.md)
- [Architecture](../ARCHITECTURE.md)
- [Shared adoption model](../../.standards/docs/adoption.md)

## Notes

The documentation-semantics boundary is deliberate: Pitlord enforces repository policy; Demon Docs and the shared checker own document structure and navigation.
