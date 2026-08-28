# Calibration References

Parent index: [Diagnosis Calibration Baselines](../INDEX.md)

## Purpose

This folder contains the machine-readable `pitlord.calibration.v1` scoring projections for the frozen real-world calibration corpora.

## Overview

Each JSON reference pins one exact corpus revision and translates its manually adjudicated architecture baseline into detector-specific `required`, `absent`, or `allowed` expectations. The Markdown calibration baseline remains the source of architectural rationale.

The 16 frozen corpora currently have three projections each:

- `<corpus>.json` — calibrated `dependency-pressure` projection;
- `<corpus>.dependency-knots.json` — frozen `dependency-knots` projection used to measure the initial knot detector before tuning;
- `<corpus>.boundary-cohesion.json` — frozen `boundary-cohesion` projection used to constrain real-corpus false positives while synthetic topology supplies the current positive control.

The projections intentionally need not classify the same architectural condition as required. A real source-concentration problem, for example, can be outside dependency-pressure, dependency-knot, and boundary/cohesion ownership.

## Direct files

### Dependency pressure

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

### Dependency knots

- `space-rocks-clean-v1.dependency-knots.json`
- `dapper.dependency-knots.json`
- `jsoup.dependency-knots.json`
- `detekt.dependency-knots.json`
- `lexicanter.dependency-knots.json`
- `polly.dependency-knots.json`
- `spectre-console.dependency-knots.json`
- `gson.dependency-knots.json`
- `hikaricp.dependency-knots.json`
- `jmh.dependency-knots.json`
- `maven.dependency-knots.json`
- `kotlinx-coroutines.dependency-knots.json`
- `now-in-android.dependency-knots.json`
- `dunit-lotusscript.dependency-knots.json`
- `jsonparser-lotusscript.dependency-knots.json`
- `volt-mx-lotusscript-toolkit.dependency-knots.json`

### Boundary / cohesion

- `space-rocks-clean-v1.boundary-cohesion.json`
- `dapper.boundary-cohesion.json`
- `jsoup.boundary-cohesion.json`
- `detekt.boundary-cohesion.json`
- `lexicanter.boundary-cohesion.json`
- `polly.boundary-cohesion.json`
- `spectre-console.boundary-cohesion.json`
- `gson.boundary-cohesion.json`
- `hikaricp.boundary-cohesion.json`
- `jmh.boundary-cohesion.json`
- `maven.boundary-cohesion.json`
- `kotlinx-coroutines.boundary-cohesion.json`
- `now-in-android.boundary-cohesion.json`
- `dunit-lotusscript.boundary-cohesion.json`
- `jsonparser-lotusscript.boundary-cohesion.json`
- `volt-mx-lotusscript-toolkit.boundary-cohesion.json`

## Related docs

- [Calibration evaluation harness](../HARNESS.md)
- [Dependency-knot baseline](../dependency-knots.md)
- [Boundary/cohesion baseline](../boundary-cohesion.md)
- [Diagnosis calibration baselines](../INDEX.md)

## Notes

Reference files are versioned reference ground truth. Detector tuning must not rewrite them merely to improve scores.
