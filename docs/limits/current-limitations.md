# Current Limitations

Parent index: [Pitlord Limits](INDEX.md)

## Purpose

This document records current Pitlord limitations that materially affect use or extension.

## Overview

Pitlord is functional as an on-demand CLI, but several product and integration surfaces remain intentionally incomplete.

## Limitations

### Generalized architecture detectors

`pitlord scan` currently implements two generalized advisory detectors: calibrated language-neutral outgoing cross-file `dependency-pressure`, and an initial language-neutral `dependency-knots` detector over strongly connected directional dependency components. Boundary-coupling/cohesion anomalies and blast-radius hotspots remain to be implemented. Both detectors have deterministic synthetic topology coverage; only dependency pressure has completed the frozen real-corpus calibration pass.

Sixteen frozen real-world corpora now have machine-readable detector references and can be evaluated reproducibly with `pitlord calibrate`. The current dependency-pressure projection scores 3 labelled true positives and 105 labelled true negatives with 0 labelled false positives, 0 labelled false negatives, and 0 severity mismatches across those pinned revisions. Twenty-six detector findings remain deliberately unlabelled—11 in Space Rocks, 13 in Maven, and 2 in JMH—and are not counted as successes. The detector now separates common non-production paths, requires both extreme production-peer fan-out and cross-directory boundary spread, exempts conventional composition seams, caps compatibility-path severity, and defers strongly reused central hubs to later detector families.

Dependency pressure remains advisory. Its findings still target physical files rather than semantic aggregates such as C# partial types or Go packages; conventional seam recognition is deterministic naming/role policy rather than semantic inference; directory regions are only an approximation of architectural boundaries; and the 26 unlabelled outputs still require future adjudication or coverage from other detector families before the real-world reference set is exhaustive.

Dependency knots are earlier in calibration. The first real sanity pass produced one directional 2-file component in Dapper, one 59-file component in jsoup, and no knot findings in detekt or Lexicanter. Those observations are not yet frozen reference judgments. The detector currently treats directory boundaries as architectural regions, intentionally excludes generic `references`/read/write edges from the cycle graph, and reports one finding per qualifying strongly connected component; package/type aggregation and real-corpus false-positive calibration remain open.

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

All entries are current as of August 28, 2026.

## Related docs

- [Roadmap](../planning/roadmap.md)
- [Architecture](../ARCHITECTURE.md)
- [Shared adoption model](../../.standards/docs/adoption.md)

## Notes

The documentation-semantics boundary is deliberate: Pitlord enforces repository policy; Demon Docs and the shared checker own document structure and navigation.
