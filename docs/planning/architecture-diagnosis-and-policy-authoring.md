# Architecture Diagnosis and Policy Authoring

Parent index: [Pitlord Planning](INDEX.md)

## Purpose
This document defines the planned expansion of Pitlord from primarily policy enforcement into a repository architecture diagnosis and policy-authoring workflow.

## Overview
Pitlord should grow from primarily enforcing known architecture rules into a workflow that can discover suspicious structures, help a developer decide whether they are intentional, and turn those decisions into durable guardrails without requiring direct JSON-schema authoring.
## Status
Planned. This document is not a statement of shipped behavior.

## Problem
Current Pitlord is strongest after architectural intent has already been encoded. `check` deterministically enforces repository-owned rules, `analyze` summarizes declared architecture areas, and `inspect` exposes Arcana architecture communities. Those foundations leave a missing first step:

> Given a repository whose problems are not yet known, identify structurally suspicious architecture and explain where a developer should investigate.

The JSON policy model is also too low-level to be the normal product interface. Users should not need to understand Pitlord's schema, area selectors, relationship arrays, or rule representation to create useful guardrails.
## Product direction
Pitlord should support a complete architecture-governance loop:

```text
inspect -> diagnose -> review -> codify or accept -> enforce -> detect regression
```

The concepts remain distinct:

- **Inspect** exposes structural facts with minimal opinion.
- **Diagnose** identifies policy-free architectural suspicions from repository structure.
- **Analyze** explains the repository relative to architecture areas already declared by the user.
- **Check** deterministically enforces explicit policy.

Diagnosis is the missing discovery layer in front of the existing policy engine, not a replacement for it.

## Diagnosis semantics
A diagnosis is a **candidate architectural concern**, not automatically a violation. Pitlord cannot know that every hub, deep chain, or cross-boundary dependency is wrong.

Diagnosis should:

- require no repository policy;
- produce deterministic evidence and stable candidate identity where practical;
- explain why a structure is unusual;
- rank candidates relative to the repository or comparable peers;
- avoid failing CI by default; and
- support explicit human disposition.

Initial detectors should make measurable structural claims rather than vague judgments such as "this code is spaghetti."

### Initial detector set
1. **Dependency knots and cyclic clusters** — strongly connected or mutually dependent regions ranked by size, density, and architectural significance.
2. **Hub and bottleneck candidates** — nodes, files, packages, or communities with anomalously high fan-in, fan-out, or both.
3. **Boundary-coupling and cohesion anomalies** — regions with unusually high cross-boundary relationships or weak internal cohesion relative to peers.
4. **Impact and blast-radius hotspots** — nodes whose transitive dependents cover an unusually large portion of the repository or subsystem.

Likely follow-on detectors include deep dependency/call chains, intermediary bypasses, dead or orphaned regions, graph-community/filesystem mismatches, architectural erosion across snapshots, unstable dependency direction, and change-aware anomalies when source-history evidence is available.

## Relative anomaly model
Pitlord should avoid universal magic thresholds such as "fan-out greater than 20 is bad." Repository size, language, framework, and subsystem role make fixed limits unreliable.

Diagnosis should primarily use relative evidence such as:

- percentile within the repository or subsystem;
- deviation from peer/community median;
- boundary-to-internal relationship ratio;
- impact radius as a fraction of the relevant graph;
- cycle size and density relative to surrounding structure; and
- change relative to a previous snapshot.

A useful diagnosis should explain itself concretely, for example:

```text
internal/api is in the top 0.4% of outgoing dependency degree
and has 6.7x the median fan-out of comparable packages.
```

The scoring model must remain deterministic and inspectable.

## Human disposition of suspicions
The developer supplies architectural intent. Every diagnosis should support at least three states:

- **Create guardrail** — this structure is undesirable; turn the decision into enforceable policy.
- **Accept structure** — this structure is intentional; persist that decision and stop presenting the same condition as suspicious in the same scope.
- **Leave unresolved** — keep the suspicion visible for later investigation.

Acceptance is not machine learning. Pitlord accumulates explicit repository decisions and uses them to refine what it presents as unresolved architecture.

### Acceptance is not a baseline
Accepted architecture and baselines have different meanings:

- **Accepted architecture:** "This structure is intentional and should not be diagnosed as a problem here."
- **Baseline:** "This policy violation exists today; suppress this existing evidence while still detecting new violations."

The persistence format for accepted decisions remains open, but it must be repository-owned, reviewable, deterministic, and semantically separate from evidence baselines.

## User-friendly policy authoring
The JSON schema should remain a stable machine-readable storage and interchange format, not the normal authoring experience.

A user should be able to create policy from a diagnosis without understanding JSON. A coupling suspicion could become an authoring form equivalent to:

```text
API should not directly depend on Storage
From: internal/api/**
To: internal/storage/**
Relationships: calls, imports
Exceptions: none
```

Pitlord or Warlock then generates the required areas, selectors, and rule representation. The same authoring experience should support manual guardrail creation when no diagnosis exists. The schema-oriented editor can remain available as an advanced representation.

### Policy proposal safety
Before a generated guardrail is saved, Pitlord should preview:

- what nodes or areas the selectors include;
- what evidence it would currently match;
- whether it would create immediate violations;
- whether it overlaps or conflicts with existing policy; and
- the final human-readable rule meaning.

One-click generation must not silently create an unexpectedly broad rule.
## Planned Warlock experience
Warlock should present architecture health rather than primarily a JSON policy editor. A repository summary may expose:

```text
7 new suspicions
3 accepted architectural decisions
12 enforced guardrails
2 current policy violations
```

A diagnosis detail could present:

```text
HIGH - Cross-boundary coupling

internal/api -> internal/storage
47 direct relationships
6.2x peer median
Introduced across 9 files

[Investigate] [Create Guardrail] [Accept Structure]
```

The raw policy editor remains useful for advanced users and exact inspection, but the normal workflow should be decision-oriented.

## Ownership boundary
**Arcana measures the graph. Pitlord judges the measurements. Warlock presents and orchestrates the workflow.**

Arcana owns graph algorithms and structural facts such as degrees, adjacency, strongly connected components, paths, call chains, communities, reachability, impact, and boundary relationship counts.

Pitlord owns anomaly interpretation, repository-relative normalization, diagnosis classes/severity/confidence, evidence aggregation, candidate identity, accepted architectural decisions, policy proposals, promotion to guardrails, policy evaluation, reports, and baselines.

Warlock owns repository selection, jobs, integration-state persistence where appropriate, and interactive review/authoring presentation.

## Calibration strategy
Arcana's deterministic synthetic topology families provide controlled detector fixtures:

- `modular` — comparatively healthy control;
- `entangled` — dependency knots and cross-coupling;
- `hub-heavy` — bottlenecks;
- `layered` — directionality and depth; and
- `dense-subsystem` — local over-coupling.

Initial diagnosis should be tuned against those controlled topologies before evaluation against real Warlock, Continuity, Space Rocks, and other repositories. Thresholds and rankings should be justified by observed separation between known structures rather than chosen arbitrarily.

## Initial implementation direction
1. Add a policy-free `diagnose` capability with the four initial detector families, deterministic machine-readable output, and enough evidence for Warlock presentation.
2. Add explicit accepted-diagnosis persistence and disposition handling.
3. Connect diagnosis candidates to human-readable policy proposals with current-repository impact preview and safe publication through the existing validation path.

The current `check`, baseline, diff, report, and policy-validation machinery remains the downstream enforcement foundation.

## Acceptance criteria
A developer should be able to:

1. point Pitlord at an unfamiliar repository without first writing policy;
2. receive a bounded, evidence-backed set of architectural suspicions;
3. investigate why each suspicion was raised;
4. explicitly accept intentional structures without abusing baselines;
5. turn an undesirable suspicion into a guardrail without editing JSON;
6. preview what the guardrail means and matches before saving it; and
7. have later Pitlord checks deterministically enforce that decision.

## Open decisions
- Exact persistence contract for accepted architectural decisions.
- Stable candidate identity across graph and detector-version changes.
- Representation of diagnosis severity/confidence separately from policy violation severity.
- Which detector parameters are fixed, repository-relative, or user-tunable.
- Which diagnoses map into existing rule types and which require new rule families.
- Whether `analyze` remains strictly declared-area analysis.
- How Warlock presents historical diagnosis changes across snapshots.

## Related docs
- [Roadmap](roadmap.md)
- [Architecture](../ARCHITECTURE.md)
- [Policy reference](../POLICY.md)
- [Baselines and CI](../BASELINES-AND-CI.md)

## Notes
The existing deterministic policy engine is not discarded. It becomes the enforcement half of a broader discovery, decision, and governance workflow.
