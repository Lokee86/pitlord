# Calibration Evaluation Harness

Parent index: [Diagnosis Calibration Baselines](INDEX.md)

## Purpose

This document owns the current machine-readable calibration contract used to score Pitlord generalized detectors against frozen corpus ground truth.

## Overview

`pitlord calibrate` evaluates one pinned corpus at a time against a strict JSON reference. It reuses the normal generalized scan engine, verifies the corpus revision before evaluation, and reports only what the frozen labels support.

## Reference contract

Calibration references live under `docs/development/calibration/references/` and use schema `pitlord.calibration.v1`.

Each reference pins:

- corpus name;
- exact Git source revision;
- detector under evaluation;
- optional scan path prefix; and
- labelled path or path-prefix expectations.

Each expectation has one detector-specific behavior:

- `required` — the detector must emit at least one matching finding;
- `absent` — any matching finding is a labelled false positive;
- `allowed` — presence or absence is not scored as precision/recall failure because the frozen architecture condition is real but may belong to another detector or a severity-only watch.

Optional `min_severity` and `max_severity` bounds score severity independently from presence.

## Evaluation

Run one pinned corpus with:

```text
pitlord calibrate \
  --repo C:/!bin/workspace/corpus/java-jsoup \
  --reference docs/development/calibration/references/jsoup.json
```

The command verifies Git `HEAD` against `source_revision` before scanning. `--skip-revision-check` exists only for controlled diagnostics and should not be used for canonical calibration runs.

`--format json` emits the machine-readable evaluation result. `--fail-on-mismatch` exits `1` when a labelled false positive, false negative, or severity mismatch exists. Without that flag the command reports mismatches but exits `0`, allowing an intentionally failing detector baseline to be measured.

## Scoring semantics

The harness reports:

- true positives and false negatives over `required` expectations;
- true negatives and false positives over `absent` expectations;
- severity mismatches over matched `required` or `allowed` expectations;
- labelled precision and recall; and
- unlabelled detector findings separately.

Unlabelled findings are never silently treated as correct or incorrect. They indicate that the frozen reference does not currently assign detector-specific judgment to that finding.

## Ground-truth ownership

The Markdown corpus baselines remain the human-readable architectural evidence. The JSON reference files are detector-specific scoring projections of those frozen judgments. A reference must not invent a classification that contradicts its corpus baseline.

Changing detector behavior does not justify changing a reference. References change only when:

- the pinned corpus revision changes;
- the frozen architectural adjudication is corrected by evidence; or
- the scoring projection is shown to misrepresent that adjudication.

## Related docs

- [Calibration baselines](INDEX.md)
- [Machine-readable references](references/INDEX.md)
- [Architecture](../../ARCHITECTURE.md)
- [Current limitations](../../limits/current-limitations.md)
- [Behavioral contract matrix](../behavioral-contract-matrix.md)

## Notes

The harness currently evaluates exact paths and path prefixes. Semantic aggregate scopes such as a C# partial type remain represented conservatively until generalized findings can target those higher-level units directly.
