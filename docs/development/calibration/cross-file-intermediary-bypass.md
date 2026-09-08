# Cross-File-Intermediary-Bypass Calibration

Parent index: [Diagnosis Calibration Baselines](INDEX.md)

## Purpose

This document records the first real-repository calibration of Pitlord's `cross-file-intermediary-bypass` detector.

## Overview

The detector is intentionally narrower than a generic call-chain-bypass recognizer. Homunculus can produce many mechanically valid `caller -> intermediary -> target` shortcuts that do not imply an architectural violation. Pitlord reports only cases where the graph contains deterministic evidence that the intermediary is an established, narrow same-region seam.

## Detector contract

A finding requires all of the following:

- resolved symbol-level `calls` evidence;
- caller and intermediary in different production files;
- caller, intermediary, and target in the same semantic namespace/module region;
- an intermediary with exactly one resolved call target outside its own file;
- a short intermediary declaration of at most five source lines;
- at least two callers remaining on the intermediary after the exceptional route is considered;
- a same-file intermediary caller whose symbol name matches the exceptional caller;
- a downstream target reached from at most two caller files after the direct shortcut; and
- different intermediary and target symbol names, suppressing same-name overload/forwarding chains.

Same-file repeated-pipeline shortcuts remain owned by `symbol-intermediary-bypass`.

## Development evidence

The motivating Space Rocks development specimen is:

```text
Control.MatchDecision
  -> Game.matchDecisionLocked
  -> Game.evaluateMatchDecisionLocked
```

Homunculus rewrites the first edge so `Control.MatchDecision` calls `evaluateMatchDecisionLocked` directly. The remaining graph still shows `Game.MatchDecision` and other callers routing through the short `matchDecisionLocked` intermediary, while the downstream evaluator remains narrowly owned.

This specimen originated from Homunculus's deterministic holdout split, but it was then used while shaping this detector. It is therefore contaminated development evidence for this detector and is not counted as independent holdout recall validation.

A second Homunculus mutation involving `streamruntime.StepContinuousBulletStreams` was explicitly rejected as positive ground truth. After mutation, two same-named one-hop wrappers in different files both call the same target and neither establishes a defensible direction from static call evidence. Pitlord must remain quiet rather than infer which wrapper "should" own the route from Homunculus's hidden mutation history.

## False-positive audit

An initial broader rule produced organic findings in Dapper, Polly, Spectre.Console, JMH, jsoup, Maven, and clean Space Rocks. Source inspection showed legitimate overloads, partial-class forwarding, inheritance/clone patterns, shared utility calls, and symmetric public wrappers rather than bypass violations.

The final suppressors were derived from those source distinctions rather than from repository names:

- same-name intermediary/target chains are excluded;
- the intermediary must be short;
- the intermediary must remain established by multiple callers;
- a same-named peer in the intermediary file must preserve the route; and
- broad downstream targets are excluded.

After those changes, raw scans of all 16 frozen clean corpora produced zero `cross-file-intermediary-bypass` findings.

## Frozen positive

The verified MatchDecision worktree is pinned to source revision `361dcdb4054725c8f1bb58727ff402fcfedf6d18` and tracked diff SHA-256 `834aeb40ad87969144277ecce73582e06f2969a437d1b69ef183a011ad9e75d3`.

A fresh full-repository Lexicon/Arcana rebuild followed by `pitlord calibrate --fail-on-mismatch` produced exactly one detector finding:

```text
cross-file-intermediary-bypass:5f4d74d53c67852b
services/game-server/internal/game/control_match.go :: MatchDecision
warning
```

The frozen positive scores 1 TP / 0 FP / 0 FN with precision 1 and recall 1 for that single labelled development specimen.

## Calibration status

Across the 16 frozen clean corpora, the negative projection is:

```text
0 TP / 16 TN / 0 FP / 0 FN
```

Together with the development specimen, the frozen labelled evidence is 1 TP / 16 TN / 0 FP / 0 FN. This establishes the current precision controls and verifies one real-repository positive shape. It does not establish independent recall because the only in-scope Space Rocks mutation was used while shaping the detector.

The regenerated pinned Homunculus corpus contains 94 mechanically eligible call-chain bypasses, but only MatchDecision satisfies this detector's declared structural contract. The stream-runtime validation mutation remains deliberately outside scope because static post-mutation evidence does not establish which symmetric wrapper owns the route.

Future recall evidence should come from a different repository or a newly frozen corpus containing cross-file intermediary seams that satisfy the detector contract before mutation selection.

## Related docs

- [Calibration evaluation harness](HARNESS.md)
- [Calibration reference set](REFERENCE-SET.md)
- [Symbol-intermediary-bypass baseline](symbol-intermediary-bypass.md)
- [Architecture](../../ARCHITECTURE.md)
- [Current limitations](../../limits/current-limitations.md)

## Notes

The MatchDecision specimen is development evidence because it shaped the detector. It must not be presented as independent holdout recall validation.
