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
| Dependency-pressure detection is language-neutral at the Pitlord layer, deduplicates symbol edges by repository path, ignores non-dependency containment edges, and keeps balanced graphs quiet | `internal/scan/dependency_pressure_test.go` |
| Baselines suppress evidence rather than entire rules | `internal/baseline/baseline_test.go` |
| Snapshot diffs classify introduced, resolved, and persistent evidence by stable fingerprint | `internal/snapshotdiff/compare_test.go` |
| Mutation verification requires exact absent/present transitions | `internal/mutation/verify_test.go` |
| Arcana pagination fails closed on invalid or truncated responses | `internal/arcana/client_test.go`, `internal/arcana/pagination.go` coverage through graph-loading tests |
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
