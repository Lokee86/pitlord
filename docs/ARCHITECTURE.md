# Architecture

Parent index: [Pitlord Documentation](INDEX.md)

## Purpose

This document defines Pitlord's implemented ownership, data flow, evaluation model, snapshot use, and runtime boundaries.

## Overview

Pitlord is a standalone repository-policy evaluator and deterministic architecture-diagnostics CLI. Lightweight content and path rules scan source files directly; semantic dependency, ownership, and cycle rules evaluate Arcana's immutable repository graph. The documentation guard consumes Demon Docs codemap ownership, Arcana's code-file inventory, and Git change evidence without parsing Markdown itself. The policy-free generalized scan surface owns a stable finding contract and repository-relative scope. Its ten opinionated detectors measure outgoing cross-file dependency pressure, strongly connected dependency knots, narrow acyclic dependency-depth corridors, unstable region-level dependency direction, exceptional file/region boundary bypasses, same-file repeated-pipeline intermediary bypasses, conservative cross-file intermediary bypasses, semantic-region boundary/cohesion anomalies, many-to-many behavioral coordination bottlenecks, and broad direct or independently amplified transitive blast radius using language-neutral Arcana relationships and production repository paths. Pitlord is not a parser, language server, graph store, or source mutation engine.

The scan engine executes registered analyzers rather than hard-coding detector calls. Each analyzer declares stable metadata, whether it requires the Arcana graph, and any normalized semantic capabilities it requires from language adapters. The default registry contains the ten language-neutral architecture analyzers and therefore still requires an immutable Arcana snapshot. The `semantic` group is also graph-backed: Lexicon adapters emit language-neutral capability, error-handler, and handler-action facts, Arcana stores those facts unchanged as graph nodes and relationships, and Pitlord applies reusable semantic rules without parsing language syntax. Its first rule, `swallowed-error`, requires control-flow, error-handling, call, and source-span capabilities and currently works across Rust, TypeScript, and Python. Analyzer implementations that do not require graph state can run through the same engine without a graph loader; the first graph-free implementation is the Rust `clippy` analyzer, which runs `cargo clippy --message-format=json --all-targets --no-deps`, retains only `clippy::...` diagnostics, deduplicates repeated target emissions, and translates primary spans plus rustc suggestion applicability into Pitlord findings and structured edits. Scan findings retain the existing detector/disposition/evidence contract while adding rule identity, analyzer identity, optional language/category metadata, exact source spans, and structured suggested fixes with applicability and text edits.

## Data flow

```text
source repository
  -> Pitlord content/path policy evaluation
  -> optional Lexicon adapters and normalized facts
  -> optional immutable Lexicon and Arcana graph snapshots
  -> Pitlord semantic area projection and policy evaluation
  -> text, JSON, SARIF, or baseline output

registered scan analyzers
  -> analyzer requirement and semantic-capability selection
  -> optional Lexicon-normalized semantic facts in Arcana
  -> optional Arcana snapshot/graph load
  -> normalized Pitlord scan finding contract
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
- normalized symbols, relationships, and semantic capability/error-handler facts;
- stable node identities and source spans; and
- content-addressed analysis snapshots.

Arcana owns:

- graph ingestion and packed storage;
- immutable graph snapshots;
- node listing, adjacency, traversal, impact, and communities; and
- protocol-level pagination and query validation.

Pitlord owns:

- repository content and path checks;
- per-code-file codemap coverage and changed-code-to-owning-document guard semantics;
- modular policy composition;
- named architecture areas;
- repository-owned policy and rule semantics;
- the generalized scan analyzer registry and finding contract, including rule/analyzer identity, advisory-versus-guard disposition, severity, scope, language/category metadata, source locations, evidence, required outcomes, structured fixes, and deterministic ordering;
- reusable semantic lint rules over adapter-declared capabilities, without language-specific parsing in Pitlord;
- calibration-reference validation and deterministic scoring of generalized detector output against frozen corpus ground truth;
- opinionated generalized detector semantics;
- area-graph projection;
- ownership and cycle interpretation;
- deterministic diagnostic evidence and fingerprints; and
- output formats and baselines.

Demon Docs owns Markdown parsing, documentation schema and structure enforcement, configured documentation-root selection, and authored codemap extraction/resolution. Pitlord consumes the exported codemap dataset but does not reinterpret Markdown syntax.

Homunculus owns controlled source mutations and expected architecture deltas used to
calibrate Pitlord behavior.

## Documentation guard

`pitlord docs` is a separate blocking guard rather than a generic policy rule. It asks Demon Docs for the current schema-1 codemap dataset, loads all Arcana `file` nodes for the current or explicit immutable snapshot, and compares the merge base of `--changed-from` and `HEAD` against the current worktree. A code file passes coverage only when at least one dataset entry resolves to that exact file; resolved file/symbol targets and non-directory pattern matches qualify, while a directory target alone does not cover every descendant.

The codemap ownership relation also defines change expectations. When a mapped code file changes, at least one document currently mapping that file must appear in the same Git change set. Findings and expected owner-document lists are sorted deterministically. Pitlord does not decide whether a document has valid front matter, headings, indexes, links, schema, or other structural semantics; Demon Docs CI owns those checks.

## Repository and snapshot loading

Repository content and path policies scan only the static roots implied by their globs and do not start Arcana. Graph policies resolve `.arcana/CURRENT` or accept an explicit immutable snapshot directory. The policy-free `scan` CLI defaults to the `architecture` analyzer group and therefore requires a resolvable immutable snapshot so generalized findings are tied to graph state. `--analyzers semantic` is also graph-backed and requires an Arcana snapshot containing the normalized semantic capabilities declared by the selected rules; unsupported or missing capabilities fail closed. `--analyzers clippy` selects the graph-free Rust analyzer and requires only a Cargo project plus Clippy. Comma-separated groups such as `architecture,semantic,clippy` combine all three streams. The underlying scan engine only loads Arcana when at least one selected analyzer declares a graph requirement. Eight file/region generalized detectors reuse the shared production file projection that removes common test, benchmark, sample, generated, vendor, tooling, devtool, and recognized test-addon paths. `symbol-intermediary-bypass` and `cross-file-intermediary-bypass` use a callable-symbol projection with equivalent non-production exclusions but deliberately retain runtime `devtools` source because development tooling can still contain executable architectural boundaries. `dependency-pressure` evaluates Arcana's broader normalized dependency-like relation set; a boundary-aware candidate must cross the raw fan-out trigger, rank at or above the 95th percentile of active production fan-out, reach at least three distinct target directories, and have at least 35% of its outgoing file dependencies cross those directory regions. Conventional composition seams are excluded from that detector, compatibility paths are capped at warning severity, and strongly reused central hubs with at least 20 incoming file dependents and fan-in at least equal to fan-out are deferred to `hub-bottleneck`. `dependency-knots` uses a narrower static source-dependency relation set—imports, inheritance/implementation, trait use, overrides, includes, depends-on, and conversions—and deliberately excludes runtime `calls` so callback/interface dispatch cannot manufacture reverse source ownership. Knot significance uses Arcana language namespaces or multi-file modules when those semantic containers are available, falling back to filesystem directories only for files without stronger ownership evidence. Tiny two-file cycles across a conventional parent/`impl`, parent/`internal`, or parent/`implementation` seam are deferred as bounded implementation-detail cycles rather than reported as architectural knots. In a boundary-aware repository scope it reports only components spanning at least two architectural regions; in a single-region narrowed scope it requires at least 6% of possible internal directed file edges before reporting a knot. `dependency-depth` then evaluates acyclic behavioral dependency corridors over calls, inheritance/implementation, trait use, overrides, includes, explicit dependencies, and conversions. It requires at least four hops, four architectural regions, and three region transitions, with each hop either singular or materially deeper than its alternatives. Import-only/reference-only edges do not establish depth. Cycles, shared convergence points, and highly reused central foundations terminate the walk so downstream implementation depth is not inherited by unrelated callers. `unstable-dependency-direction` projects the same static source-dependency relation family onto architectural regions and measures instability as outgoing-region count divided by total incoming-plus-outgoing region count. A candidate requires at least three files and three region couplings on each side, at least two source files supporting the boundary, an incoming-heavy source, an outgoing-heavy target, and at least a 0.25 instability increase. Region-level strongly connected components are deferred to `dependency-knots`. `boundary-bypass` evaluates file-level behavioral shortcuts across three distinct architectural regions. At least three peer files in one upstream region must establish a gateway; direct upstream-to-downstream access must remain at or below 34% of that peer count; the gateway must send at least two behavioral dependencies and at least 40% of its outgoing behavioral targets into the downstream region; and the downstream target may be reached from at most two source regions. Import/reference-only edges, nested source/target ownership trees, broadly shared targets, and weakly concentrated gateways remain quiet. `symbol-intermediary-bypass` operates on resolved Arcana `calls` at symbol granularity. A candidate requires a direct caller-to-target call, an uncalled one-hop local intermediary to that target, at least three active sibling wrapper paths toward the same downstream file, and an additional support symbol shared by those peer callers and the bypassing caller outside the downstream target file. Constructor targets are excluded, and a caller that still reaches a resolved same-name intermediary overload is treated as preserving the intermediary family. `cross-file-intermediary-bypass` separately requires caller and intermediary to occupy different production files inside one semantic region, a short one-hop intermediary with at least two remaining callers, a same-named peer in the intermediary file that still preserves the route, a narrow downstream target reached from at most two caller files, and distinct intermediary/target names; these constraints suppress overload, shared-utility, symmetric-wrapper, and broad-target triangles. `boundary-cohesion` evaluates only semantic namespace or multi-file-module regions with at least six production files. Static dependency-like relations define outward boundary edges, while a broader support family including runtime calls and read/write relations measures internal collaboration. A candidate must have at least six outgoing cross-region dependencies, reach at least three independent target regions, reach at least 0.25 target regions per member, rank in the top decile of semantic peers by target-region reach, and have internal support both in the weakest peer quartile and at or below 0.5 support relationships per member. Regions with incoming fan-in at least equal to outgoing fan-out are deferred to `hub-bottleneck`. `hub-bottleneck` first requires an extreme direct fan-in candidate: at least eight incoming production-file dependents, at least three times the active-peer median, top-5% active-peer fan-in, at least three incoming architectural regions, and at least four direct outgoing production dependencies. It then applies the discriminator that removed the real-corpus false positives: only `calls`, `reads`, and `writes` count as behavioral flow, which must enter from at least six architectural regions, leave toward at least four architectural regions, and retain an outgoing/incoming behavioral-degree ratio of at least 0.4. This makes shared type contracts and data/domain roots quiet even when their signature/reference centrality is extreme. `impact-blast-radius` has two evidence lanes. Broad direct exposure requires at least `max(6, ceil(3% of production peers))` direct dependents, at least two incoming architectural regions, and at least two dependents supported by non-metadata relationships such as calls, inheritance/implementation, reads/writes, includes, explicit dependencies, or conversions. In scopes with at least 32 peers it normally requires top-5% direct fan-in, with exceptions for surfaces reaching at least 20% of peers or at least 10 architectural regions; final breadth requires either 5% direct peer reach or 10 regions. Files that do not qualify directly may qualify through transitive amplification: at least 25% peer reach, top-5% active impact, at least 2x amplification over direct fan-in, at least three impacted regions, at least three substantial independent first-hop branches, and no single branch contributing more than 80% of the closure. Generic reference/import/annotation reuse does not establish direct-impact evidence, and the branch discriminator prevents a leaf from inheriting one gateway's downstream closure. Pitlord queries Arcana through `arcana.query.v1`; it does not read Arcana's storage files.

`calibrate` adds a reference-evaluation layer around the same scan engine. The reference pins a Git source revision, detector, and labelled path/path-prefix expectations, and may additionally pin the SHA-256 of the complete tracked worktree diff for a controlled mutation specimen. Revision or worktree-diff mismatch fails closed before graph evaluation. Required and absent labels contribute to labelled precision/recall, allowed labels can constrain severity without forcing the dependency-pressure detector to own unrelated maintenance concerns, and scan findings that match no label remain explicitly unlabelled rather than being assumed correct.

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
the same game-server source root. That mutation is a symbol-level intermediary bypass inside one source file, so it remains intentionally outside file/region `boundary-bypass`. The separate `symbol-intermediary-bypass` detector now owns that shape directly from symbol-level call evidence and the mutation is frozen as its first diff-pinned real-repository positive. Because this development-split mutation was used to shape the detector, independent validation/holdout mutations are still required for broader recall claims.

## Snapshot diffs

Pitlord can evaluate one policy against two immutable snapshots concurrently and compare
flattened evidence by stable fingerprint. The result preserves introduced, resolved, and
persistent evidence. This is the initial changed-snapshot seam for CI and future Warlock
runtime integration; it does not yet avoid evaluating unchanged policy scopes.

## Runtime model

Pitlord is currently an on-demand CLI. Registered analyzers execute in deterministic registration order; duplicate or missing analyzer IDs fail closed, analyzer failures are propagated, declared semantic capabilities are validated before rule execution, and only analyzers that declare a graph requirement cause Arcana graph loading. The CLI exposes `architecture`, `semantic`, and `clippy` analyzer groups, which may be combined. It starts Arcana protocol processes for bounded
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
