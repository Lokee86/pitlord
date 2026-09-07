# Behavioral Contract Matrix

Parent index: [Development Documentation](INDEX.md)

## Purpose

This document maps critical Pitlord invariants to focused tests and release gates.

## Overview

The matrix prevents durable safety and compatibility contracts from depending only on broad package coverage.

## Contracts

| Contract | Primary tests |
| --- | --- |
| Repository-only rules do not require an Arcana graph | `internal/policy/filesystem_test.go`, `internal/checker/checker_test.go` |
| Forbidden content reports exact matching source lines | `internal/policy/filesystem_test.go` |
| Required content reports each selected file without a match and the missing selection case | `internal/policy/filesystem_test.go` |
| Policy normalization rejects ambiguous, incompatible, or ineffective fields | `internal/policy/validate_test.go` |
| Embedded schema matches supported policy and baseline contracts | `internal/schema/schema_test.go`, `cmd/pitlord/schema_test.go` |
| Includes are relative, recursive, deduplicated, and cycle-safe | `internal/policy/load_test.go` |
| Evidence ordering and reporting are deterministic | `internal/policy/evaluate_test.go`, `internal/report/report_test.go` |
| Documentation guarding consumes Demon Docs resolved codemap ownership rather than parsing Markdown, requires exact per-file coverage, treats resolved glob files as coverage but not directory-only targets, and requires a changed mapped code file to change at least one owning document | `internal/docguard/evaluate_test.go`, `cmd/pitlord/docs_test.go` |
| Generalized scans require no Pitlord policy and preserve deterministic finding ordering plus non-null finding/evidence arrays | `internal/scan/scan_test.go`, `cmd/pitlord/scan_test.go`, `internal/report/scan_test.go` |
| Analyzer selection is deterministic and fail-closed: duplicate/unknown groups fail, graph loading is required only when selected analyzers declare it, and semantic analyzers reject snapshots that do not advertise their required normalized capabilities | `internal/scan/scan_test.go`, `internal/scan/semantic_swallowed_error_test.go`, `internal/scan/semantic_unobserved_outcome_test.go`, `cmd/pitlord/scan_test.go` |
| `swallowed-error` remains language-neutral: local propagate/record/recover actions and proven downstream fallback, enclosing-propagation, or explicit-suppression dispositions suppress the advisory, while generic `continuation` does not | `internal/scan/semantic_swallowed_error_test.go` |
| `unobserved-outcome` reports only adapter-proven fallible/async outcome obligations without contained consumption; Pitlord does not infer Rust, TypeScript/JavaScript, or Python syntax itself | `internal/scan/semantic_unobserved_outcome_test.go` |
| Rust Clippy execution is graph-free, preserves upstream `clippy::...` rule identity and primary spans, deduplicates repeated Cargo target diagnostics, and emits structured edits only from applicable rustc suggestions | `internal/scan/rust_clippy_test.go`, `cmd/pitlord/scan_test.go` |
| Dependency-pressure detection uses explicit Arcana file nodes as language-neutral peers, aggregates semantic edges by file, excludes virtual/package-only and common non-production paths, requires top-5% fan-out plus cross-directory boundary spread, defers conventional composition seams and strongly reused central hubs, caps compatibility-path severity, and separates isolated hubs from dense regional pressure | `internal/scan/dependency_pressure_test.go`, `internal/scan/source_role_test.go`, `internal/scan/dependency_pressure_calibration_test.go` |
| Dependency-knot detection uses production file nodes, static source-dependency relations rather than runtime calls or generic references/read/write edges, deterministic iterative SCC traversal, Arcana namespace/multi-file-module regions with directory fallback, density fallback in single-region scopes, stable component-level finding identity, and warning-level bounded two-region seams | `internal/scan/dependency_knots_test.go`, `internal/scan/dependency_regions_test.go`, `internal/scan/dependency_knots_calibration_test.go`, `internal/scan/dependency_scc.go` |
| Dependency-depth detection walks only acyclic behavioral/architectural relations, excludes import/reference-only propagation, requires >=4 hops across >=4 architectural regions with >=3 region transitions, permits branching only when one branch remains at least two hops deeper, and terminates at cycles, shared convergence points, or highly reused central foundations | `internal/scan/dependency_depth_test.go`, `internal/scan/dependency_depth_test_helpers_test.go`, `internal/scan/dependency_depth.go`, `internal/scan/dependency_depth_metrics.go` |
| Unstable-dependency-direction detection aggregates static source dependencies to architectural regions, requires >=3 files and >=3 total region couplings on both sides, >=2 source files supporting the boundary, an incoming-heavy source, an outgoing-heavy target, and >=0.25 instability increase, while deferring region SCCs to dependency-knot ownership | `internal/scan/unstable_dependency_direction_test.go`, `internal/scan/unstable_dependency_direction.go`, `internal/scan/unstable_dependency_direction_metrics.go`, `internal/scan/dependency_scc.go` |
| Boundary-bypass detection requires a minority direct behavioral shortcut around a gateway established by >=3 same-region peers, requires >=2 gateway targets and >=40% gateway concentration in the downstream region, rejects targets reached from >2 regions, excludes import/reference-only edges and nested source/target ownership, and remains explicitly file/region-level | `internal/scan/boundary_bypass_test.go`, `internal/scan/boundary_bypass_ownership_test.go`, `internal/scan/boundary_bypass.go`, `internal/scan/boundary_bypass_metrics.go` |
| Boundary-cohesion detection evaluates Arcana semantic namespaces or multi-file modules rather than directory-only regions, measures outward spread with static dependency-like relations, measures internal support with a broader relation family including calls and reads/writes, requires six-member semantic peers, independent target-region spread plus top-decile reach, both lower-quartile and <=0.5-per-member internal support, and defers incoming-heavy hubs | `internal/scan/boundary_cohesion_test.go`, `internal/scan/boundary_cohesion_metrics.go`, `internal/scan/dependency_regions_test.go` |
| Hub-bottleneck detection requires extreme production-file fan-in but reports only a many-to-many behavioral waist: `calls`/`reads`/`writes` must arrive from >=6 architectural regions, leave toward >=4 regions, and maintain >=0.4 outgoing/incoming behavioral degree; reference-only hubs, shared leaves, narrow behavioral facades, and composition roots remain quiet | `internal/scan/hub_bottleneck_test.go`, `internal/scan/hub_bottleneck.go`, `internal/scan/dependency_regions_test.go` |
| Impact/blast-radius detection has two deterministic lanes: broad direct exposure requires repository-relative extreme fan-in, cross-region breadth, and non-metadata dependency evidence; lower-direct-fan-in transitive exposure requires top-5% reach over >=25% of peers, >=2x amplification, >=3 regions, >=3 substantial independent first-hop branches, and <=80% dominant-branch share. Metadata-only direct fan-in, single-region direct fan-in, inherited single-gateway closure, dominant-branch reach, ordinary layering, and family-wide transitive reachability remain quiet | `internal/scan/impact_blast_radius_test.go`, `internal/scan/impact_blast_radius.go`, `internal/scan/impact_direct_exposure.go`, `internal/scan/impact_blast_radius_metrics.go` |
| Calibration references are strict and source-revision pinned; they target exactly one legacy detector or built-in analyzer, may further constrain rule/language and exact source location, score required/absent/allowed expectations deterministically, constrain severity independently, keep unlabelled findings explicit, and can opt into fully-labelled regression failure | `internal/calibration/reference_test.go`, `internal/calibration/reference_files_test.go`, `internal/calibration/evaluate_test.go`, `cmd/pitlord/calibrate_test.go`, `internal/report/calibration_test.go` |
| Baselines suppress evidence rather than entire rules | `internal/baseline/baseline_test.go` |
| Snapshot diffs classify introduced, resolved, and persistent evidence by stable fingerprint | `internal/snapshotdiff/compare_test.go` |
| Mutation verification requires exact absent/present transitions | `internal/mutation/verify_test.go` |
| Arcana discovery prefers explicit Pitlord configuration and repository-prepared Lexicon + Arcana installations, supports combined-release and nearby source-checkout layouts, and does not consult retired Grimoire provider state | `internal/arcana/command_test.go` |
| Arcana graph loading requests the maximum neighbor page, validates returned relationship counts, and fails closed on truncation; node-list pagination also fails closed on invalid continuation state | `internal/arcana/client_test.go`, `internal/arcana/qualified_test.go`, `internal/arcana/pagination.go` coverage through graph-loading tests |
| Arcana subprocesses terminate when the caller context expires | `internal/arcana/protocol_process_test.go` |
| Incompatible Arcana protocol identifiers fail closed | `internal/arcana/protocol_process_test.go` |
| Interrupted Arcana response streams cannot publish partial graph evidence | `internal/arcana/protocol_process_test.go` |
| SARIF preserves evidence identity and locations | `internal/report/sarif_test.go` |

## Release gate

The minimum release gate is:

```bash
go test ./...
python ../engineering-standards/tools/docs_policy/check.py --repo .
```

## Related docs

- [Documentation coverage](documentation-coverage.md)
- [Architecture](../ARCHITECTURE.md)
- [Baselines and CI](../BASELINES-AND-CI.md)

## Notes

Update this matrix when an invariant changes, a focused test moves, or a new persistent or compatibility boundary is introduced.
