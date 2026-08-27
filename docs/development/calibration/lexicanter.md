# Lexicanter Calibration Baseline

Parent index: [Diagnosis Calibration Baselines](INDEX.md)

## Purpose

This document freezes the manually adjudicated reference architecture state of the Lexicanter corpus used for Pitlord dependency-pressure calibration.

## Overview

- Corpus: `lexicanter`
- Source revision: `eac754788b3cf18a930c085c1c49f8f353e18107`
- Commit subject: `many small fixes I forgot to make individual commits for`
- Languages represented in the prepared graph: Svelte/TypeScript and Rust
- Baseline role: external real-world control with known maintenance debt

The baseline records corpus state only. It does not prescribe detector changes.

## Current Pitlord observation

A repository-wide dependency-pressure scan reports three advisory production findings:

- `src/app/App.svelte` — high; 13 outgoing file dependencies
- `src/app/layouts/File.svelte` — warning; 9 outgoing file dependencies
- `src/app/layouts/Settings.svelte` — warning; 9 outgoing file dependencies

A scan scoped to `src/app/layouts/` currently reports no findings.

The difference between repository-wide and layout-scoped output is preserved as an observation, not treated here as ground truth.

## Manually adjudicated production state

| Scope | Expected interpretation | Notes |
| --- | --- | --- |
| `src/app/App.svelte` | clean; intentional application composition root | Approximately 128 lines. Owns top-level tab/layout composition, application-close handling, window controls, version/platform display, and top-level layout resizing. Its broad import surface is intrinsic to assembling the application UI. |
| `src/app/layouts/File.svelte` | real maintenance pressure / baseline debt | Approximately 490 lines. The file tab owns filesystem open/import behavior, legacy file migration/loading, wholesale language-state hydration, documentation initialization, pronunciation regeneration, layout restoration, database verification/retrieval and sync decisions, plus the related UI. These are related to file lifecycle, but the implementation combines persistence, migration, synchronization, domain-state reconstruction, and presentation orchestration in one component. |
| `src/app/layouts/Settings.svelte` | real architectural pressure / baseline debt | Approximately 909 lines. The component owns local settings persistence, experimental-mode state, database account verification/sync/deletion, theme selection and custom-theme storage, autosave behavior, lect rename/deletion and cross-domain data rewrites, relative-lexicon import, and layout reset/import/export alongside the settings UI. This is a broad set of operational and domain responsibilities concentrated in one component. |

## Baseline debt details

### `File.svelte`

The dependency-pressure finding points at a real maintainability seam, although dependency count alone does not express the full problem.

The inspected implementation combines several distinct operational concerns:

- selecting save locations and opening/importing files;
- compatibility handling for legacy `.lexc` versions;
- reconstructing most of the global `Language` state;
- rebuilding documentation and pronunciation-derived state;
- restoring window/panel layout state;
- verifying database credentials and comparing local/remote file versions;
- downloading remote state and deciding whether to overwrite/upload;
- presenting loading/error state in the same component.

This is frozen as known baseline maintenance debt. The corpus should not be modified merely to make this control appear clean.

### `Settings.svelte`

This is the strongest architectural-pressure specimen in the Lexicanter baseline.

The inspected implementation contains functions for:

- online file-version lookup;
- experimental/alchemy settings;
- settings-file migration and persistence;
- database-account verification and synchronization;
- theme and custom-theme management;
- autosave scheduling;
- lect rename/deletion with mutations across lexicon, phrasebook, pronunciations, orthographies, and etymological data;
- relative-lexicon import;
- remote database deletion;
- layout reset, import, and export.

These responsibilities are broader than one cohesive settings presentation concern and constitute known pre-existing architectural debt in the pinned corpus revision.

## Ground-truth summary

Lexicanter is a useful real-world control precisely because it is not uniformly clean.

- `App.svelte` should not be treated as architectural debt merely because it has the broadest dependency surface; it is the application composition root.
- `File.svelte` is a genuine maintenance-pressure seam and should remain labelled as baseline debt.
- `Settings.svelte` is genuine architectural pressure and the clearest current dependency-pressure true-positive candidate in this corpus.

The frozen judgments above are reference labels for the pinned revision. They are independent of whether the current detector emits, suppresses, or changes severity for those files under different scan scopes.

## Related docs

- [Diagnosis Calibration Baselines](INDEX.md)
- [detekt](detekt.md)
- [jsoup](jsoup.md)
- [Dapper](dapper.md)
- [Space Rocks Clean Control v1](space-rocks-clean-v1.md)
- [Architecture diagnosis plan](../../planning/architecture-diagnosis-and-policy-authoring.md)

## Notes

This baseline is tied to the exact source revision above. Re-adjudicate only if that corpus revision changes or source evidence demonstrates that one of these frozen judgments is wrong.
