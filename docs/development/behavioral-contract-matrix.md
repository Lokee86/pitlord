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
| Generalized scans require no Pitlord policy and preserve deterministic finding ordering plus non-null finding/evidence arrays | `internal/scan/scan_test.go`, `cmd/pitlord/scan_test.go`, `internal/report/scan_test.go` |
| Dependency-pressure detection uses explicit Arcana file nodes as language-neutral peers, aggregates semantic edges by file, excludes virtual/package-only and common non-production paths, requires top-5% fan-out plus cross-directory boundary spread, defers conventional composition seams and strongly reused central hubs, caps compatibility-path severity, and separates isolated hubs from dense regional pressure | `internal/scan/dependency_pressure_test.go`, `internal/scan/source_role_test.go`, `internal/scan/dependency_pressure_calibration_test.go` |
| Dependency-knot detection uses production file nodes, static source-dependency relations rather than runtime calls or generic references/read/write edges, deterministic iterative SCC traversal, Arcana namespace/multi-file-module regions with directory fallback, density fallback in single-region scopes, stable component-level finding identity, and warning-level bounded two-region seams | `internal/scan/dependency_knots_test.go`, `internal/scan/dependency_regions_test.go`, `internal/scan/dependency_knots_calibration_test.go`, `internal/scan/dependency_scc.go` |
| Boundary-cohesion detection evaluates Arcana semantic namespaces or multi-file modules rather than directory-only regions, measures outward spread with static dependency-like relations, measures internal support with a broader relation family including calls and reads/writes, requires six-member semantic peers, independent target-region spread plus top-decile reach, both lower-quartile and <=0.5-per-member internal support, and defers incoming-heavy hubs | `internal/scan/boundary_cohesion_test.go`, `internal/scan/boundary_cohesion_metrics.go`, `internal/scan/dependency_regions_test.go` |
| Calibration references are strict, pinned to a source revision, score required/absent/allowed expectations deterministically, constrain severity independently, and never silently score unlabelled findings | `internal/calibration/reference_test.go`, `internal/calibration/reference_files_test.go`, `internal/calibration/evaluate_test.go`, `cmd/pitlord/calibrate_test.go`, `internal/report/calibration_test.go` |
| Baselines suppress evidence rather than entire rules | `internal/baseline/baseline_test.go` |
| Snapshot diffs classify introduced, resolved, and persistent evidence by stable fingerprint | `internal/snapshotdiff/compare_test.go` |
| Mutation verification requires exact absent/present transitions | `internal/mutation/verify_test.go` |
| Arcana provider discovery prefers explicit/configured and repository-prepared Grimoire/Lexicon providers before generic PATH fallback | `internal/arcana/command_test.go` |
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
