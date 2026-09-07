# Calibration Evaluation Harness

Parent index: [Diagnosis Calibration Baselines](INDEX.md)

## Purpose

This document owns the current machine-readable calibration contract used to score Pitlord architecture detectors, semantic analyzers, and other built-in analyzers against frozen corpus ground truth.

## Overview

`pitlord calibrate` evaluates one pinned corpus at a time against a strict JSON reference. It reuses the normal scan engine, executes the exact built-in analyzer named by the reference, verifies the corpus revision and any optional frozen worktree diff before evaluation, and reports only what the frozen labels support. Legacy detector references remain valid and resolve their `detector` directly to the matching architecture analyzer.

## Reference contract

Calibration references live under `docs/development/calibration/references/` and use schema `pitlord.calibration.v1`.

Each reference pins:

- corpus name;
- exact Git source revision;
- optional SHA-256 of the complete tracked worktree diff;
- exactly one execution target: legacy `detector` or built-in `analyzer` ID;
- optional `rule_id` and `language` selectors when analyzer targeting is used;
- optional `require_fully_labelled` strictness for regression-gate references;
- optional scan path prefix; and
- labelled exact-path or path-prefix expectations.

An exact-path expectation may additionally include a `location` object with `start_line` and optional start/end columns and end line. Location matching is intended for source-local findings such as semantic lint diagnostics where multiple findings may exist in one file.

`worktree_diff_sha256` is used for controlled mutation specimens. Pitlord does not apply the mutation; it hashes `git diff --binary --full-index --no-ext-diff HEAD --` and fails closed if the supplied worktree does not match the frozen specimen.

Each expectation has one target-specific behavior:

- `required` — the selected detector/analyzer/rule/language target must emit at least one matching finding;
- `absent` — any matching target finding is a labelled false positive;
- `allowed` — presence or absence is not scored as precision/recall failure because the frozen architecture condition is real but may belong to another detector or a severity-only watch.

Optional `min_severity` and `max_severity` bounds score severity independently from presence.

## Evaluation

Run one pinned corpus with:

```text
pitlord calibrate \
  --repo C:/!bin/workspace/corpus/java-jsoup \
  --reference docs/development/calibration/references/jsoup.json
```

The command verifies Git `HEAD` against `source_revision` before scanning. When `worktree_diff_sha256` is present it also verifies the exact tracked diff. `--skip-revision-check` skips both state checks and exists only for controlled diagnostics; it should not be used for canonical calibration runs.

`--format json` emits the machine-readable evaluation result. `--fail-on-mismatch` exits `1` when a labelled false positive, false negative, or severity mismatch exists. For references with `require_fully_labelled: true`, unlabelled findings selected by the reference also make the result a mismatch. Without `--fail-on-mismatch` the command reports mismatches but exits `0`, allowing an intentionally failing calibration baseline to be measured.

## Scoring semantics

The harness reports:

- true positives and false negatives over `required` expectations;
- true negatives and false positives over `absent` expectations;
- severity mismatches over matched `required` or `allowed` expectations;
- labelled precision and recall; and
- unlabelled target findings separately.

Unlabelled findings are never silently treated as correct. By default they remain informational because many architecture references are intentionally incomplete projections. A reference with `require_fully_labelled: true` converts any unlabelled matching finding into a regression-gate mismatch without changing how the finding is counted in the summary.

An expectation may use `path_prefix: "."` to match the entire repository. Detector-wide negative controls use this form so any new finding is scored against the absent expectation instead of escaping as unlabelled output. Source-local semantic references instead use exact path plus location and set `require_fully_labelled: true`, so additional same-target findings remain visible and fail canonical regression runs.

## Ground-truth ownership

Markdown calibration records remain the human-readable evidence. JSON reference files are machine-readable scoring projections of those frozen judgments. Architecture references project architectural adjudication; semantic references project source-level rule adjudication from the semantic calibration record. A reference must not invent a classification that contradicts its human-readable owner.

Changing detector behavior does not justify changing a reference. References change only when:

- the pinned corpus revision changes;
- the frozen architectural adjudication is corrected by evidence; or
- the detector-specific scoring projection is shown to misrepresent that adjudication.

That final case matters when a frozen architectural problem is real but the current detector family does not own its signal. For example, the JMH generator concentration and Maven model/project-builder concentration remain frozen structural pressure, but are `allowed` rather than `required` for `dependency-pressure` because their evidence is source concentration/decomposition rather than necessarily anomalous outgoing dependency breadth.

## Related docs

- [Calibration baselines](INDEX.md)
- [Machine-readable references](references/INDEX.md)
- [Architecture](../../ARCHITECTURE.md)
- [Current limitations](../../limits/current-limitations.md)
- [Behavioral contract matrix](../behavioral-contract-matrix.md)

## Notes

The harness evaluates exact paths, path prefixes, and optional exact source locations. Semantic aggregate scopes such as a C# partial type remain represented conservatively until scan findings can target those higher-level units directly.
