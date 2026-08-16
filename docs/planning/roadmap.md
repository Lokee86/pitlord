# Roadmap

Parent index: [Pitlord Planning](INDEX.md)

## Purpose

This document records planned Pitlord work and unresolved sequencing.

## Overview

Pitlord remains independently usable while expanding from architecture-policy enforcement into a broader architecture diagnosis, decision, and repository-governance engine within Warlock.

## Current status

Active planning.

## Expected ownership

- Pitlord owns anomaly interpretation, diagnosis identity, accepted architectural decisions, policy authoring semantics, rule evaluation, evidence identity, baselines, and machine-readable diagnostics.
- Arcana owns graph storage, traversal, communities, impact, paths, strongly connected components, and other structural measurements.
- Lexicon owns language facts.
- Demon Docs and the shared checker own Markdown structure, navigation, and documentation coverage checks.
- Warlock owns scheduling, enrolled-repository orchestration, and interactive diagnosis/policy presentation.

## Planned work

1. Add policy-free architecture diagnosis over Arcana structural evidence, beginning with dependency knots, hub/bottleneck anomalies, boundary-coupling/cohesion anomalies, and impact hotspots.
2. Add explicit human disposition for diagnoses: create a guardrail, accept the structure as intentional, or leave it unresolved.
3. Add user-friendly policy authoring that can promote a diagnosis into validated policy without requiring JSON-schema knowledge, including a current-repository match preview before publication.
4. Keep accepted architectural decisions semantically distinct from evidence baselines and define a durable repository-owned persistence contract for them.
5. Calibrate diagnosis ranking against Arcana's deterministic modular, entangled, hub-heavy, layered, and dense-subsystem topologies before relying on real-repository tuning.
6. Add reusable policy packs with explicit versioning and compatibility.
7. Integrate shared documentation-policy results into Warlock repository health.
8. Add change-aware evaluation planning without weakening full-policy correctness.
9. Improve cross-snapshot, diagnosis-history, and baseline migration diagnostics.

See [Architecture diagnosis and policy authoring](architecture-diagnosis-and-policy-authoring.md) for the product workflow and detector direction.

## Acceptance criteria

Each planned feature must preserve deterministic evidence, explicit ownership, repository-only execution where applicable, schema validation, and independent CLI use. Diagnosis must remain distinguishable from enforcement: a suspicion is not automatically a CI violation, and accepted architecture is not an evidence baseline.

## Open decisions

- Exact persistence and versioning contract for accepted architectural decisions.
- Stable diagnosis identity across detector and snapshot changes.
- Which detector classes can promote into existing policy rules and which require new rule families.
- How reusable policy-pack versions are resolved and pinned by repositories.
- Whether Warlock invokes Pitlord directly or through a shared daemon job protocol.
- Which change-planning evidence is sufficient to safely skip unaffected graph regions.

## Implemented references

- `require_content` is implemented and documented in [Policy reference](../POLICY.md).
- Shared documentation governance is defined in [Engineering Standards](../../../engineering-standards/docs/INDEX.md).

## Related docs

- [Architecture diagnosis and policy authoring](architecture-diagnosis-and-policy-authoring.md)
- [Architecture](../ARCHITECTURE.md)
- [Current limitations](../limits/current-limitations.md)

## Notes

This roadmap is not a statement of shipped behavior.
