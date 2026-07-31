# Documentation Procedure

Parent index: [Pitlord Documentation](INDEX.md)

## Purpose

This document defines the required process for changing Pitlord documentation.

## Overview

Documentation changes accompany implementation changes that affect rules, commands, schemas, diagnostics, ownership, state, CI, baselines, snapshots, mutations, or extension seams.

## Procedure

1. Identify the changed command, rule family, package owner, stateful flow, or machine-readable contract.
2. Update the exact reference owner and implemented architecture owner in the same change.
3. Update baseline or CI operations when fingerprints, result classes, exit behavior, or SARIF change.
4. Update the embedded JSON Schema and reference examples together when policy fields change.
5. Update the coverage map when commands, packages, or independent flows are added, removed, renamed, or reassigned.
6. Update the behavioral-contract matrix when a durable invariant or protecting test changes.
7. Put future work in planning and current gaps in limits.
8. Update every affected `INDEX.md` and relative link.
9. Run Pitlord tests, policy validation, and the shared documentation checker.
10. Report documentation impact and known gaps.

## Verification

Run:

```bash
go test ./...
python .standards/docs_policy/check.py --repo .
go run ./cmd/pitlord validate --policy ../engineering-standards/policies/pitlord/documentation-core.json
```

When reviewing a pull request, also run the shared checker with `--changed-from` against the merge base.

## Related docs

- [Documentation policy](documentation-policy.md)
- [Policy reference](POLICY.md)
- [Development coverage](development/documentation-coverage.md)
- [Behavioral contract matrix](development/behavioral-contract-matrix.md)

## Notes

Schema-only and documentation-only changes still require tests because embedded schemas, examples, fixtures, and generated output may diverge.
