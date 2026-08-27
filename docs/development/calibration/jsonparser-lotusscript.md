# JSONParser LotusScript Calibration Baseline

Parent index: [Diagnosis Calibration Baselines](INDEX.md)

## Purpose

This document freezes the manually adjudicated reference architecture state of the JSONParser LotusScript corpus used for Pitlord dependency-pressure calibration.

## Overview

- Corpus: `jsonparser-ls`
- Source revision: `2406bc4431225ef83276700bda853ae247bd9cfd`
- Commit subject: `Update json.ls`
- Language: LotusScript
- Baseline role: external small-library control

The baseline records corpus state only. It does not prescribe detector changes.

## Current Pitlord observation

The repository-wide dependency-pressure scan reports zero findings.

## Manually adjudicated production state

The repository has one substantive source file, `json.ls`, containing `JSONParser`, `JSONObject`, and `JSONArray`.

| Scope | Expected interpretation | Notes |
| --- | --- | --- |
| `json.ls` / `JSONParser` | maintenance watch | Roughly 300-line single-file parser library. Object, array, scalar, string and property parsing share mutable parser state and are tightly co-located. This is understandable at the corpus's current size but is the dominant maintenance-pressure seam. |
| `json.ls` / `JSONObject` | clean | Small map-like JSON object wrapper. |
| `json.ls` / `JSONArray` | clean | Small array wrapper. |

The audit does not freeze parsing correctness concerns as architectural findings. The relevant architectural fact is that essentially all parser behavior resides inside one source unit, so cross-file dependency pressure is expected to be silent even when intra-file maintenance pressure exists.

## Ground-truth summary

JSONParser is a compact monolithic library rather than a dependency-heavy architecture. Zero Pitlord findings is correct for cross-file dependency pressure, while `JSONParser` remains a legitimate maintenance watch because the repository's substantive parsing behavior is concentrated in one unit.

## Related docs

- [Diagnosis Calibration Baselines](INDEX.md)
- [Architecture diagnosis plan](../../planning/architecture-diagnosis-and-policy-authoring.md)

## Notes

This baseline is tied to the exact source revision above. Re-adjudicate only if that corpus revision changes or source evidence demonstrates that one of these frozen judgments is wrong.
