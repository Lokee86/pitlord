# Symbol-Intermediary-Bypass Calibration

Parent index: [Diagnosis Calibration Baselines](INDEX.md)

## Purpose

This document records the first real-repository calibration of Pitlord's symbol-level `symbol-intermediary-bypass` detector.

## Overview

The calibration freezes one real Space Rocks development mutation and a 16-corpus negative projection, while keeping development evidence distinct from independent validation or holdout recall.

## Owned shape

The detector owns a narrower condition than file/region `boundary-bypass`: one caller bypasses a local behavioral intermediary while sibling call paths still establish that intermediary layer at symbol granularity.

A finding currently requires:

- a direct caller-to-target `calls` edge;
- an uncalled one-hop local intermediary that still points to that target;
- at least three active sibling wrapper paths in the same source file toward the same downstream file;
- a support symbol shared by the bypassing caller and those sibling wrapper callers, outside the downstream target file;
- a behavioral target rather than a constructor; and
- no evidence that the caller still routes through another same-name intermediary overload.

Tests, generated/vendor/sample code, recognized test addons, and other non-production source roles remain excluded. Runtime `devtools` source is intentionally eligible because the first real positive lives there.

## Real positive

The first positive is the previously verified Space Rocks Homunculus development mutation at base revision:

```text
361dcdb4054725c8f1bb58727ff402fcfedf6d18
```

The mutation rewrites:

```text
handleDebugAddScore
    -> addDebugScoreForPlayer
    -> AddPlayerScore
```

into:

```text
handleDebugAddScore
    -----------------> AddPlayerScore
```

while leaving the sibling score/lives handlers on their wrapper routes. The mutation also preserves the shared `resolveCommandTargetPlayerIDs` support step.

The frozen worktree diff is a one-line source change in:

```text
services/game-server/internal/devtools/player_counters.go
```

and is pinned by calibration hash:

```text
a3df7f53c63db4b4ec7b1d5b9484938dabce474422cd8c431ec66df7285a5445
```

Pitlord does not apply the mutation. Homunculus retains mutation ownership; `pitlord calibrate` only verifies the base Git revision and optional worktree-diff hash before evaluating the detector.

On the frozen mutated worktree the detector emits exactly one finding:

```text
handleDebugAddScore bypasses an established local intermediary
```

The evidence identifies:

- direct call: `handleDebugAddScore -> AddPlayerScore`;
- abandoned intermediary: `addDebugScoreForPlayer -> AddPlayerScore`;
- three intact sibling wrapper paths; and
- shared support step `resolveCommandTargetPlayerIDs`.

The untouched Space Rocks corpus remains quiet.

## False-positive calibration

The first broad symbol rule treated any abandoned one-hop wrapper beside three active same-file wrappers as an intermediary convention. Real-corpus auditing showed that this was too broad.

### GUT test framework

The untouched Space Rocks corpus initially produced findings inside vendored GUT test-framework files under `client/addons/gut`. These are not production architecture and are excluded by source role. The exclusion is narrow enough that `internal/devtools` remains eligible.

### jsoup fluent API and validation helpers

The first 16-corpus pass then produced 20 capped findings in jsoup, concentrated in wrapper-heavy fluent APIs such as `HttpConnection`. Public convenience methods independently call shared validation helpers such as `notEmptyParam` and `notNullParam`; those helper calls do not establish ownership through another fluent setter.

The discriminator that removed this family is repeated pipeline evidence: active peer wrapper callers must share an additional support symbol with the bypassing caller, outside the downstream target file. This preserves the Space Rocks handler pipeline while making jsoup quiet.

### Maven overloads and construction

Maven produced five initial findings in `DefaultModelBuilder`. Four disappeared under the shared-pipeline requirement. The remaining case treated a convenience `derive` overload around `DefaultModelBuilderResult` construction as an abandoned intermediary even though `build` still routes through another `derive` overload.

Constructor targets are now outside this detector. Construction/factory ownership is a distinct architectural question, and Arcana may not resolve every overloaded call strongly enough for behavioral bypass analysis. Same-name resolved overloads are also treated as continued intermediary use rather than abandonment.

## Frozen negative projection

The final detector was rerun against all 16 pinned real-world corpora:

```text
corpora               16
symbol findings        0
```

Detector-specific references use a repository-root absent expectation, so any future `symbol-intermediary-bypass` finding in those clean controls is scored as a false positive rather than left unlabelled.

The clean projection is therefore:

```text
true positives         0
true negatives        16
false positives        0
false negatives        0
severity mismatches    0
unlabelled findings    0
```

The diff-pinned Space Rocks mutation adds one required positive and currently scores:

```text
true positives         1
false negatives        0
severity mismatches    0
unlabelled findings    0
```

## Interpretation

This calibration establishes that Pitlord can detect the known real symbol-level bypass while preserving the 16 frozen clean repositories. It does not establish broad recall: the Space Rocks mutation is the development mutation used to shape the detector, not an untouched validation or holdout positive.

The next recall test should use Homunculus validation and holdout call-chain bypass mutations without changing the detector in response to individual examples. Those experiments should preserve the current 16-corpus false-positive controls.

## Calibration-state pinning

`pitlord.calibration.v1` now supports optional `worktree_diff_sha256` in addition to `source_revision`. When present, calibration verifies the SHA-256 of:

```text
git diff --binary --full-index --no-ext-diff HEAD --
```

before scanning. This allows controlled dirty worktrees to be frozen as calibration specimens without moving source-mutation ownership into Pitlord.

A reference can also use `path_prefix: "."` to cover the entire repository. This is used by the 16 negative controls so new detector output cannot hide as an unlabelled finding outside one preselected path.

## Related docs

- [Calibration evaluation harness](HARNESS.md)
- [Calibration reference set](REFERENCE-SET.md)
- [Boundary-bypass baseline](boundary-bypass.md)
- [Architecture](../../ARCHITECTURE.md)
- [Current limitations](../../limits/current-limitations.md)

## Notes

The detector intentionally depends on resolved Arcana `calls` relationships. Unresolved dynamic dispatch cannot be reconstructed reliably by Pitlord and should remain outside deterministic pass/fail claims until the underlying semantic graph can represent it.
