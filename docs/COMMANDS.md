# Command Reference

Parent index: [Pitlord Documentation](INDEX.md)

## Purpose

This document defines Pitlord's public CLI surfaces and the intended distinction between commands that can otherwise look overlapping.

## Overview

Pitlord keeps policy enforcement, generalized diagnosis, architecture-model analysis, repository inspection, CI support, and calibration evidence as separate commands. The command groups below make those roles explicit without changing their runtime ownership.

## Primary workflow

| Command | Role | Requires Arcana |
| --- | --- | --- |
| `check` | Enforce an explicit repository policy and return compliance diagnostics. | Only for graph-backed rules. |
| `scan` | Run policy-free registered analyzers and return normalized advisory findings. | Only when a selected analyzer requires graph state. |
| `analyze` | Project declared policy areas onto the repository graph and report ownership/dependency structure without enforcing policy rules. | Yes. |
| `inspect` | Inspect Arcana architecture-community summaries to discover structure before authoring policy. | Yes. |

The distinction is deliberate:

- use `inspect` to discover repository structure;
- use `analyze` to understand how a proposed area model maps onto that structure;
- use `check` to enforce explicit repository-owned rules;
- use `scan` for Pitlord's generalized architecture, semantic, or native-tool judgments without authoring a policy.

## Policy and CI commands

### `check`

```text
pitlord check --repo <root> --policy <pitlord.json> [--baseline <baseline.json>]
```

`check` evaluates repository content/path rules directly and graph-backed dependency, ownership, or area-cycle rules through Arcana. It also accepts `--homunculus-manifest` as an alternate policy source for generated specimens.

### `validate`

```text
pitlord validate --policy <pitlord.json>
```

Loads, composes, normalizes, and validates policy without evaluating a repository.

### `schema`

```text
pitlord schema --kind policy|baseline
```

Prints the embedded Draft 2020-12 JSON Schema for the selected contract.

### `baseline`

```text
pitlord baseline --repo <root> --policy <pitlord.json> [--output <baseline.json>]
```

Creates evidence-level baseline fingerprints from current policy findings.

### `diff`

```text
pitlord diff --before-snapshot <path> --repo <after-root> --policy <pitlord.json>
```

Evaluates the same policy against two immutable snapshots and classifies evidence as introduced, resolved, or persistent.

### `docs`

```text
pitlord docs --repo <root> --changed-from <revision> [--ddocs-root <path>]
```

Runs the blocking documentation guard over Demon Docs codemap ownership, Arcana code-file inventory, and Git change evidence.

## Diagnosis commands

### `scan`

```text
pitlord scan --repo <root> [--analyzers architecture,semantic,clippy] [--path-prefix src] [--format text|json]
```

The default `architecture` group runs Pitlord's ten graph-backed architecture detectors. `semantic` runs Pitlord-owned cross-language semantic rules over Lexicon-normalized capabilities stored in Arcana. `clippy` normalizes Rust Clippy diagnostics and is graph-free. Analyzer groups may be combined.

All architecture detectors and current semantic rules are advisory. Exact detector mechanics and thresholds belong in [Architecture](ARCHITECTURE.md); evidence strength and known limitations belong in [Current limitations](limits/current-limitations.md).

### `analyze`

```text
pitlord analyze --repo <root> --policy <pitlord.json> [--format text|json|dot]
```

`analyze` requires a policy that declares at least one area. It projects those areas over Arcana and reports structural coverage and area relationships. Useful options include:

- `--scope` and `--scope-exclude` for ownership coverage boundaries;
- `--ownership-kinds` for the node kinds included in coverage analysis;
- `--relations` for normalized architecture relationships;
- `--example-limit` for bounded unowned/overlap examples;
- `--format dot` for graph visualization.

Unlike `check`, `analyze` does not report policy compliance. It is an authoring and inspection surface for the declared architecture model.

### `inspect`

```text
pitlord inspect --repo <root> [--path-prefix src]
```

Exposes Arcana architecture-community summaries for repository discovery and policy authoring.

## Evidence and development commands

These commands are public and supported, but primarily exist to maintain deterministic evidence and integration contracts rather than the ordinary repository workflow.

### `calibrate`

```text
pitlord calibrate --repo <corpus> --reference <reference.json> [--format text|json]
```

Runs exactly one built-in detector or analyzer against a pinned machine-readable calibration reference. `--fail-on-mismatch` turns the result into a regression gate; strict references can additionally reject unexpected unlabelled findings.

### `verify-mutation`

```text
pitlord verify-mutation --manifest <architecture-mutation.json> --baseline-snapshot <path>
```

Verifies exact absent-to-present and present-to-absent architecture relationship transitions declared by a Homunculus mutation manifest.

### `generate`

```text
pitlord generate --homunculus-manifest <manifest.json> [--output pitlord.json]
```

Converts a Homunculus specimen manifest into reusable Pitlord policy.

### `version`

```text
pitlord version
```

Prints the Pitlord version.

## General options and behavior

Graph-backed commands accept `--arcana` as an explicit Arcana executable override where applicable. Normal discovery is defined in [Arcana process boundary](ARCANA_PROCESS_BOUNDARY.md).

Common exit semantics are:

- `0` — command succeeded and no blocking finding/mismatch was present;
- `1` — a command-specific blocking diagnostic, introduced finding, or expectation mismatch was present;
- `2` — usage, configuration, policy, snapshot, protocol, or execution failure.

Individual commands may intentionally use only the subset of these statuses relevant to their role.

## Related docs

- [Policy reference](POLICY.md)
- [Architecture](ARCHITECTURE.md)
- [Baselines and CI](BASELINES-AND-CI.md)
- [Calibration harness](development/calibration/HARNESS.md)
- [Current limitations](limits/current-limitations.md)

## Notes

The command reference owns CLI role and usage descriptions. Detailed policy semantics belong in the policy reference, detector mechanics belong in architecture, and calibration claims belong in development calibration records and current limitations.
