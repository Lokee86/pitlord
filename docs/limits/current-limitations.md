# Current Limitations

Parent index: [Pitlord Limits](INDEX.md)

## Purpose

This document records current Pitlord limitations that materially affect use or extension.

## Overview

Pitlord is functional as an on-demand CLI, but several product and integration surfaces remain intentionally incomplete.

## Limitations

### Generalized architecture detectors

`pitlord scan` currently implements five generalized advisory detectors: calibrated language-neutral outgoing cross-file `dependency-pressure`, calibrated language-neutral `dependency-knots` over strongly connected static source-dependency components, `boundary-cohesion` over semantic architectural regions, calibrated language-neutral `hub-bottleneck` analysis over many-to-many behavioral coordination waists, and calibrated language-neutral `impact-blast-radius` analysis over broad direct exposure or independently amplified transitive dependent surfaces. All five detectors have deterministic synthetic topology coverage and frozen detector-specific projections across the same 16 real-world corpora.

Sixteen frozen real-world corpora now have machine-readable detector references and can be evaluated reproducibly with `pitlord calibrate`. The current dependency-pressure projection scores 3 labelled true positives and 105 labelled true negatives with 0 labelled false positives, 0 labelled false negatives, and 0 severity mismatches across those pinned revisions. Twenty-six detector findings remain deliberately unlabelled—11 in Space Rocks, 13 in Maven, and 2 in JMH—and are not counted as successes. The detector now separates common non-production paths, requires both extreme production-peer fan-out and cross-directory boundary spread, exempts conventional composition seams, caps compatibility-path severity, and defers strongly reused central hubs to later detector families.

Dependency pressure remains advisory. Its findings still target physical files rather than semantic aggregates such as C# partial types or Go packages; conventional seam recognition is deterministic naming/role policy rather than semantic inference; directory regions are only an approximation of architectural boundaries; and the 26 unlabelled outputs still require future adjudication or coverage from other detector families before the real-world reference set is exhaustive.

The dependency-knot corpus audit remains frozen as the tuning reference. Across the 16 pinned corpora, the initial detector produced 29 findings: 7 required real knot-pressure components, 14 allowed bounded/intentional cyclic seams, and 8 false architectural cycles. The tuned detector preserves the 7 required knots while removing the frozen false positives and severity errors: the current projection scores 7 true positives, 39 true negatives, 0 false positives, 0 false negatives, 0 severity mismatches, and 0 unlabelled findings. It now excludes runtime `calls` from source-cycle direction, uses Arcana semantic namespace or multi-file-module regions before filesystem-directory fallback, recognizes GUT as test tooling, suppresses tiny conventional implementation seams, and does not promote bounded two-region cycles solely because they are dense. Remaining knot limitations are chiefly the physical-file finding scope, reliance on available Arcana semantic containers, and deterministic naming/path heuristics for conventional implementation seams.

The boundary/cohesion detector is deliberately conservative. It evaluates only semantic namespace or multi-file-module regions with at least six production files, measures outward spread using static dependencies, and measures internal support with a broader relationship family that includes runtime calls and read/write edges. A finding requires at least six outgoing cross-region dependencies, at least three independent target regions, at least 0.25 target regions per member, top-decile peer reach, internal support in the weakest peer quartile, and no more than 0.5 internal support relationships per member. Incoming-heavy regions are deferred to hub/bottleneck analysis. Against the 30 frozen real-corpus negative controls it produces 30 true negatives, 0 false positives, and 0 unlabelled findings. Those references contain no required real-world positive, so this result validates the current false-positive controls but not real-world recall; synthetic topology remains the positive control and the detector remains advisory rather than guard-ready.

The hub/bottleneck detector is also advisory but has one frozen required real-world positive. The untouched direct-centrality version emitted 37 findings and scored 1 TP / 10 TN / 33 FP / 0 FN with 2 severity mismatches. The tuned detector treats direct fan-in only as candidate generation; the decisive evidence is many-to-many behavioral flow over `calls`, `reads`, and `writes`. A candidate must receive behavior from at least six architectural regions, dispatch behavior to at least four regions, and have behavioral outgoing degree at least 40% of behavioral incoming degree. This suppresses central data/contracts, result/domain models, base abstractions, and narrow logging/telemetry facades while preserving the known Space Rocks networking gravity well. The current frozen projection scores 1 TP / 43 TN / 0 FP / 0 FN, 0 severity mismatches, and 0 unlabelled findings. Recall is still constrained by only one required real-world positive, so this is not broad proof of bottleneck recall across architectures.

The impact/blast-radius detector is advisory but now has corrected real-world recall coverage. Independent corpus auditing showed that the original all-negative calibration had confused healthy shared foundations with absence of impact. Against the corrected references, the untouched transitive rule scores 1 TP / 8 TN / 44 FP / 18 FN, while the historical branch-only tuning scores 0 TP / 52 TN / 0 FP / 19 FN: it removed inherited-closure noise by suppressing every required positive. The current two-lane detector restores broad direct exposure while retaining the branch discriminator for low-direct-fan-in transitive amplification. Direct exposure uses repository-relative fan-in, cross-region breadth, and non-metadata dependency evidence; transitive exposure still requires at least three substantial independent first-hop branches with no branch above 80% of reach. Across 150 frozen expectations the current projection scores 19 TP / 52 TN / 0 FP / 0 FN, with 0 severity mismatches and 0 unlabelled findings. The detector remains advisory because a broad change surface can be intentional even when its exposure is measured reliably.

**Removal condition:** boundary/cohesion gains required real-world positive coverage, and each generalized detector has sufficient recall evidence to support its intended advisory or guard disposition without weakening the current false-positive controls.

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
