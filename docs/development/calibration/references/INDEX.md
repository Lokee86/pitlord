# Calibration References

Parent index: [Diagnosis Calibration Baselines](../INDEX.md)

## Purpose

This folder contains the machine-readable `pitlord.calibration.v1` scoring projections for frozen real-world architecture and semantic calibration corpora.

## Overview

Each JSON reference pins one exact corpus revision and translates its manually adjudicated calibration record into `required`, `absent`, or `allowed` expectations. Architecture references target one detector. Semantic references target one analyzer/rule/language tuple and may require exact source locations plus fully labelled output. The corresponding Markdown calibration record remains the human-readable source of rationale.

The 16 frozen architecture corpora currently have ten projections each, plus two diff-pinned mutation positives. Three additional semantic references pin the first audited Rust, TypeScript, and Python semantic findings:

- `<corpus>.json` — calibrated `dependency-pressure` projection;
- `<corpus>.dependency-knots.json` — frozen `dependency-knots` projection used to measure the initial knot detector before tuning;
- `<corpus>.dependency-depth.json` — independently audited `dependency-depth` projection separating cross-region behavioral corridors from contained, import-only, cyclic, convergent, and shared-foundation paths;
- `<corpus>.unstable-dependency-direction.json` — independently audited region-level `unstable-dependency-direction` false-positive projection covering cycle deferral, supported stability gradients, JMH generator/runtime coupling, and Maven compatibility layering while real-world positive recall remains unvalidated;
- `<corpus>.boundary-bypass.json` — independently audited file/region `boundary-bypass` false-positive projection covering nested ownership, shared targets, weak gateway concentration, and clean controls while real-world positive recall remains unvalidated;
- `<corpus>.symbol-intermediary-bypass.json` — whole-repository absent control for the symbol-level repeated-pipeline bypass detector;
- `space-rocks-symbol-bypass.symbol-intermediary-bypass.json` — diff-pinned required positive for the known Space Rocks `handleDebugAddScore` mutation;
- `<corpus>.cross-file-intermediary-bypass.json` — whole-repository absent control for the conservative cross-file seam-bypass detector;
- `space-rocks-match-decision.cross-file-intermediary-bypass.json` — diff-pinned required positive for the Space Rocks `Control.MatchDecision` bypass mutation;
- `<corpus>.boundary-cohesion.json` — frozen `boundary-cohesion` projection used to constrain real-corpus false positives while synthetic topology supplies the current positive control; and
- `<corpus>.hub-bottleneck.json` — frozen `hub-bottleneck` projection separating genuine central maintenance gravity from healthy shared abstractions and data/contracts; and
- `<corpus>.impact-blast-radius.json` — independently audited `impact-blast-radius` projection covering required broad direct/transitive change surfaces, allowed intentional maintenance watches, and absent inherited/metadata controls;
- `lexicanter.semantic-swallowed-error-rust.json` — strict Rust `swallowed-error` reference for the frozen Lexicanter parse-error handler;
- `space-rocks-clean-v1.semantic-unobserved-outcome-typescript.json` — strict TypeScript `unobserved-outcome` reference for the frozen font-readiness Promise; and
- `space-rocks-clean-v1.semantic-swallowed-error-python.json` — strict Python `swallowed-error` reference for the remaining ambiguous secondary-parse handler.

Architecture projections intentionally need not classify the same architectural condition as required. A real source-concentration problem, for example, can be outside dependency-pressure, dependency-knot, dependency-depth, unstable-dependency-direction, boundary-bypass, symbol-intermediary-bypass, cross-file-intermediary-bypass, boundary/cohesion, hub/bottleneck, and impact/blast-radius ownership.

## Direct files

### Semantic analyzers

- `lexicanter.semantic-swallowed-error-rust.json`
- `space-rocks-clean-v1.semantic-unobserved-outcome-typescript.json`
- `space-rocks-clean-v1.semantic-swallowed-error-python.json`

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

### Dependency depth

- `space-rocks-clean-v1.dependency-depth.json`
- `dapper.dependency-depth.json`
- `jsoup.dependency-depth.json`
- `detekt.dependency-depth.json`
- `lexicanter.dependency-depth.json`
- `polly.dependency-depth.json`
- `spectre-console.dependency-depth.json`
- `gson.dependency-depth.json`
- `hikaricp.dependency-depth.json`
- `jmh.dependency-depth.json`
- `maven.dependency-depth.json`
- `kotlinx-coroutines.dependency-depth.json`
- `now-in-android.dependency-depth.json`
- `dunit-lotusscript.dependency-depth.json`
- `jsonparser-lotusscript.dependency-depth.json`
- `volt-mx-lotusscript-toolkit.dependency-depth.json`

### Unstable dependency direction

- `space-rocks-clean-v1.unstable-dependency-direction.json`
- `dapper.unstable-dependency-direction.json`
- `jsoup.unstable-dependency-direction.json`
- `detekt.unstable-dependency-direction.json`
- `lexicanter.unstable-dependency-direction.json`
- `polly.unstable-dependency-direction.json`
- `spectre-console.unstable-dependency-direction.json`
- `gson.unstable-dependency-direction.json`
- `hikaricp.unstable-dependency-direction.json`
- `jmh.unstable-dependency-direction.json`
- `maven.unstable-dependency-direction.json`
- `kotlinx-coroutines.unstable-dependency-direction.json`
- `now-in-android.unstable-dependency-direction.json`
- `dunit-lotusscript.unstable-dependency-direction.json`
- `jsonparser-lotusscript.unstable-dependency-direction.json`
- `volt-mx-lotusscript-toolkit.unstable-dependency-direction.json`

### Boundary bypass

- `space-rocks-clean-v1.boundary-bypass.json`
- `dapper.boundary-bypass.json`
- `jsoup.boundary-bypass.json`
- `detekt.boundary-bypass.json`
- `lexicanter.boundary-bypass.json`
- `polly.boundary-bypass.json`
- `spectre-console.boundary-bypass.json`
- `gson.boundary-bypass.json`
- `hikaricp.boundary-bypass.json`
- `jmh.boundary-bypass.json`
- `maven.boundary-bypass.json`
- `kotlinx-coroutines.boundary-bypass.json`
- `now-in-android.boundary-bypass.json`
- `dunit-lotusscript.boundary-bypass.json`
- `jsonparser-lotusscript.boundary-bypass.json`
- `volt-mx-lotusscript-toolkit.boundary-bypass.json`

### Symbol intermediary bypass

- `space-rocks-clean-v1.symbol-intermediary-bypass.json`
- `dapper.symbol-intermediary-bypass.json`
- `jsoup.symbol-intermediary-bypass.json`
- `detekt.symbol-intermediary-bypass.json`
- `lexicanter.symbol-intermediary-bypass.json`
- `polly.symbol-intermediary-bypass.json`
- `spectre-console.symbol-intermediary-bypass.json`
- `gson.symbol-intermediary-bypass.json`
- `hikaricp.symbol-intermediary-bypass.json`
- `jmh.symbol-intermediary-bypass.json`
- `maven.symbol-intermediary-bypass.json`
- `kotlinx-coroutines.symbol-intermediary-bypass.json`
- `now-in-android.symbol-intermediary-bypass.json`
- `dunit-lotusscript.symbol-intermediary-bypass.json`
- `jsonparser-lotusscript.symbol-intermediary-bypass.json`
- `volt-mx-lotusscript-toolkit.symbol-intermediary-bypass.json`
- `space-rocks-symbol-bypass.symbol-intermediary-bypass.json` — diff-pinned required-positive mutation

### Cross-file intermediary bypass

- `space-rocks-clean-v1.cross-file-intermediary-bypass.json`
- `dapper.cross-file-intermediary-bypass.json`
- `jsoup.cross-file-intermediary-bypass.json`
- `detekt.cross-file-intermediary-bypass.json`
- `lexicanter.cross-file-intermediary-bypass.json`
- `polly.cross-file-intermediary-bypass.json`
- `spectre-console.cross-file-intermediary-bypass.json`
- `gson.cross-file-intermediary-bypass.json`
- `hikaricp.cross-file-intermediary-bypass.json`
- `jmh.cross-file-intermediary-bypass.json`
- `maven.cross-file-intermediary-bypass.json`
- `kotlinx-coroutines.cross-file-intermediary-bypass.json`
- `now-in-android.cross-file-intermediary-bypass.json`
- `dunit-lotusscript.cross-file-intermediary-bypass.json`
- `jsonparser-lotusscript.cross-file-intermediary-bypass.json`
- `volt-mx-lotusscript-toolkit.cross-file-intermediary-bypass.json`
- `space-rocks-match-decision.cross-file-intermediary-bypass.json` — diff-pinned required-positive mutation

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

### Impact / blast radius

- `space-rocks-clean-v1.impact-blast-radius.json`
- `dapper.impact-blast-radius.json`
- `jsoup.impact-blast-radius.json`
- `detekt.impact-blast-radius.json`
- `lexicanter.impact-blast-radius.json`
- `polly.impact-blast-radius.json`
- `spectre-console.impact-blast-radius.json`
- `gson.impact-blast-radius.json`
- `hikaricp.impact-blast-radius.json`
- `jmh.impact-blast-radius.json`
- `maven.impact-blast-radius.json`
- `kotlinx-coroutines.impact-blast-radius.json`
- `now-in-android.impact-blast-radius.json`
- `dunit-lotusscript.impact-blast-radius.json`
- `jsonparser-lotusscript.impact-blast-radius.json`
- `volt-mx-lotusscript-toolkit.impact-blast-radius.json`

### Hub / bottleneck

- `space-rocks-clean-v1.hub-bottleneck.json`
- `dapper.hub-bottleneck.json`
- `jsoup.hub-bottleneck.json`
- `detekt.hub-bottleneck.json`
- `lexicanter.hub-bottleneck.json`
- `polly.hub-bottleneck.json`
- `spectre-console.hub-bottleneck.json`
- `gson.hub-bottleneck.json`
- `hikaricp.hub-bottleneck.json`
- `jmh.hub-bottleneck.json`
- `maven.hub-bottleneck.json`
- `kotlinx-coroutines.hub-bottleneck.json`
- `now-in-android.hub-bottleneck.json`
- `dunit-lotusscript.hub-bottleneck.json`
- `jsonparser-lotusscript.hub-bottleneck.json`
- `volt-mx-lotusscript-toolkit.hub-bottleneck.json`

## Related docs

- [Calibration evaluation harness](../HARNESS.md)
- [Semantic lint calibration](../semantic-lints.md)
- [Dependency-knot baseline](../dependency-knots.md)
- [Dependency-depth baseline](../dependency-depth.md)
- [Unstable-dependency-direction baseline](../unstable-dependency-direction.md)
- [Boundary-bypass baseline](../boundary-bypass.md)
- [Symbol-intermediary-bypass baseline](../symbol-intermediary-bypass.md)
- [Cross-file-intermediary-bypass baseline](../cross-file-intermediary-bypass.md)
- [Boundary/cohesion baseline](../boundary-cohesion.md)
- [Hub/bottleneck baseline](../hub-bottleneck.md)
- [Impact/blast-radius baseline](../impact-blast-radius.md)
- [Diagnosis calibration baselines](../INDEX.md)

## Notes

Reference files are versioned reference ground truth. Detector tuning must not rewrite them merely to improve scores.
