# Architecture

Parent index: [Pitlord Documentation](INDEX.md)

## Purpose

This document defines Pitlord's implemented ownership, data flow, evaluation model, snapshot use, and runtime boundaries.

## Overview

Pitlord is a standalone repository-policy evaluator and deterministic architecture-diagnostics CLI. Lightweight content and path rules scan source files directly; semantic dependency, ownership, and cycle rules evaluate Arcana's immutable repository graph. The policy-free generalized scan surface owns a stable finding contract and repository-relative scope. Its first two opinionated detectors measure outgoing cross-file dependency pressure and strongly connected dependency knots using language-neutral Arcana relationships and production repository paths. Pitlord is not a parser, language server, graph store, or source mutation engine.

## Data flow

```text
source repository
  -> Pitlord content/path policy evaluation
  -> optional Lexicon adapters and normalized facts
  -> optional immutable Lexicon and Arcana graph snapshots
  -> Pitlord semantic area projection and policy evaluation
  -> text, JSON, SARIF, or baseline output

Arcana snapshot
  -> Pitlord policy-free generalized scan contract
  -> deterministic ordered findings
  -> text or JSON output

pinned corpus + calibration reference
  -> Git revision verification
  -> generalized scan
  -> detector-specific labelled evaluation
  -> mismatch and precision/recall report
```

## Component ownership

Lexicon owns:

- language parsing and semantic resolution;
- normalized symbols and relationships;
- stable node identities and source spans; and
- content-addressed analysis snapshots.

Arcana owns:

- graph ingestion and packed storage;
- immutable graph snapshots;
- node listing, adjacency, traversal, impact, and communities; and
- protocol-level pagination and query validation.

Pitlord owns:

- repository content and path checks;
- modular policy composition;
- named architecture areas;
- repository-owned policy and rule semantics;
- the generalized scan finding, advisory-versus-guard disposition, severity, scope, evidence, required-outcome, and ordering contract;
- calibration-reference validation and deterministic scoring of generalized detector output against frozen corpus ground truth;
- opinionated generalized detector semantics;
- area-graph projection;
- ownership and cycle interpretation;
- deterministic diagnostic evidence and fingerprints; and
- output formats and baselines.

Homunculus owns controlled source mutations and expected architecture deltas used to
calibrate Pitlord behavior.

## Repository and snapshot loading

Repository content and path policies scan only the static roots implied by their globs and do not start Arcana. Graph policies resolve `.arcana/CURRENT` or accept an explicit immutable snapshot directory. The policy-free `scan` command also requires a resolvable immutable snapshot so generalized findings are always tied to graph state. Both shipped generalized detectors reuse one production file projection that removes common test, benchmark, sample, generated, vendor, tooling, and devtool paths. `dependency-pressure` evaluates Arcana's broader normalized dependency-like relation set; a boundary-aware candidate must cross the raw fan-out trigger, rank at or above the 95th percentile of active production fan-out, reach at least three distinct target directories, and have at least 35% of its outgoing file dependencies cross those directory regions. Conventional composition seams are excluded from that detector, compatibility paths are capped at warning severity, and strongly reused central hubs with at least 20 incoming file dependents and fan-in at least equal to fan-out are deferred to bottleneck/blast-radius analysis. `dependency-knots` uses a narrower directional relation set—imports, calls, inheritance/implementation, trait use, overrides, includes, depends-on, and conversions—and runs an iterative deterministic strongly connected component pass. In a boundary-aware repository scope it reports only components spanning at least two directories; in a single-directory narrowed scope it requires at least 6% of possible internal directed file edges before reporting a knot. Pitlord queries Arcana through `arcana.query.v1`; it does not read Arcana's storage files.

`calibrate` adds a reference-evaluation layer around the same scan engine. The reference pins a Git source revision, detector, and labelled path/path-prefix expectations. Revision mismatch fails closed before graph evaluation. Required and absent labels contribute to labelled precision/recall, allowed labels can constrain severity without forcing the dependency-pressure detector to own unrelated maintenance concerns, and scan findings that match no label remain explicitly unlabelled rather than being assumed correct.

Node listings are paginated. Pitlord validates counts, offsets, page continuity, and
terminal state, and fails closed when an older Arcana binary reports truncation without a
continuation offset.

Pitlord separately calculates:

- source prefixes whose nodes must be inspected; and
- source prefixes whose outgoing relationships are required.

This allows ownership-only checks to scan large repositories without loading any adjacency
data. Dependency and area-cycle rules load only relationships from their source areas.

## Evaluation model

Policy normalization occurs before graph loading so invalid selectors cannot cause partial
evaluation. Area matching uses normalized repository-relative paths, exclusions, and node
kinds.

Direct dependency rules evaluate Arcana edges without rebuilding graph algorithms.
Ownership rules evaluate area membership. Area-cycle rules build a small graph whose nodes
are declared policy areas and whose edges are supported by Arcana evidence, then run a
deterministic strongly connected component pass over that policy graph.

Diagnostics are grouped by rule. Evidence is sorted deterministically and retains exact
source and target nodes for human output, SARIF, baselines, and snapshot diffs.

## Mutation verification

Homunculus real-repository mutation manifests use Lexicon-qualified `path::name` symbols
and declare expected added and removed relationships. Pitlord resolves the source symbol in
both snapshots, asks Arcana for its filtered outgoing neighbors, and compares the exact
target qualified name.

Added expectations require absent→present. Removed expectations require present→absent.
This verifies the architectural delta directly instead of inferring success from changed
source text or from a broad policy diagnostic.

The first real Space Rocks development mutation verified both sides of a call-chain bypass:

- introduced `handleDebugAddScore --calls--> AddPlayerScore`; and
- removed `handleDebugAddScore --calls--> addDebugScoreForPlayer`.

Both snapshots were built from the manifest's exact base commit and mutated worktree using
the same game-server source root.

## Snapshot diffs

Pitlord can evaluate one policy against two immutable snapshots concurrently and compare
flattened evidence by stable fingerprint. The result preserves introduced, resolved, and
persistent evidence. This is the initial changed-snapshot seam for CI and future Warlock
runtime integration; it does not yet avoid evaluating unchanged policy scopes.

## Runtime model

Pitlord is currently an on-demand CLI. It starts Arcana protocol processes for bounded
queries and has no daemon state of its own. The caller context owns cancellation and any
deadline; process, protocol, pagination, or response-completeness failure prevents a
partial graph-backed compliance result. See [Arcana process boundary](ARCANA_PROCESS_BOUNDARY.md).
A future Warlock runtime may trigger Pitlord against changed snapshots, but Pitlord remains
independently usable and owns its policy semantics.

## Related docs

- [Policy reference](POLICY.md)
- [Baselines and CI](BASELINES-AND-CI.md)
- [Documentation coverage](development/documentation-coverage.md)
- [Current limitations](limits/current-limitations.md)

## Notes

Future Warlock integration must preserve Pitlord's independent CLI use and ownership of policy semantics.
