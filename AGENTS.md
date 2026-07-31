# AGENTS.md

Pitlord is a deterministic repository-policy and architecture-diagnostics CLI.

## Read first

- `README.md`
- `docs/INDEX.md`
- `docs/ARCHITECTURE.md`
- `docs/POLICY.md`
- `docs/documentation-policy.md`
- `docs/documentation-procedure.md`
- `docs/development/documentation-coverage.md`

## Ownership rules

- Keep CLI entry points and argument wiring under `cmd/pitlord` thin.
- Keep policy models, validation, normalization, and evaluation under `internal/policy`.
- Keep graph protocol integration under `internal/arcana`.
- Keep baselines, reports, schemas, snapshots, mutations, and snapshot diffs in their existing focused packages.
- Preserve deterministic ordering, evidence identity, and machine-readable output contracts.
- Repository-only rules must not require Arcana.
- Do not weaken validation to accept ineffective or ambiguous policy rules.
- Exclude nested `.worktrees/` and `.workingtrees/` from repository traversal.

## Documentation

Documentation is part of the implementation.

Before completing a change:

1. Update affected policy, CLI, schema, architecture, baseline, diagnostics, and operational documentation in the same change.
2. Update `docs/development/documentation-coverage.md` when a command, package, rule family, stateful flow, or ownership boundary changes.
3. Update `docs/development/behavioral-contract-matrix.md` when a durable invariant or its protecting test changes.
4. Keep current behavior separate from planning and limitations.
5. Run the configured documentation and normal test gates.
6. Do not report documentation as complete or current unless the checks pass and known semantic gaps are disclosed.

## Completion report

```text
Documentation impact:
- Inspected:
- Updated:
- Not affected:
- Compliance check:
- Known documentation gaps:
```
