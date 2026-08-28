# Calibration References

Parent index: [Diagnosis Calibration Baselines](../INDEX.md)

## Purpose

This folder contains the machine-readable `pitlord.calibration.v1` scoring projections for the frozen real-world calibration corpora.

## Overview

Each JSON reference pins one exact corpus revision and translates its manually adjudicated architecture baseline into detector-specific `required`, `absent`, or `allowed` expectations. The human-readable Markdown baseline remains the source of architectural rationale.

## Direct files

- `space-rocks-clean-v1.json`
- `dapper.json`
- `jsoup.json`
- `detekt.json`
- `lexicanter.json`
- `polly.json`
- `spectre-console.json`
- `gson.json`
- `hikaricp.json`
- `jmh.json`
- `maven.json`
- `kotlinx-coroutines.json`
- `now-in-android.json`
- `dunit-lotusscript.json`
- `jsonparser-lotusscript.json`
- `volt-mx-lotusscript-toolkit.json`

## Related docs

- [Calibration evaluation harness](../HARNESS.md)
- [Diagnosis calibration baselines](../INDEX.md)

## Notes

Reference files are versioned ground truth. Detector tuning must not rewrite them merely to improve scores.
