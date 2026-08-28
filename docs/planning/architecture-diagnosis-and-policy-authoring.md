# Architecture Diagnosis and Policy Authoring

Parent index: [Pitlord Planning](INDEX.md)

## Purpose
This document defines the planned expansion of Pitlord from primarily policy enforcement into a repository architecture diagnosis and policy-authoring workflow.

## Overview
Pitlord should grow from primarily enforcing known architecture rules into a workflow that can discover suspicious structures, help a developer decide whether they are intentional, and turn those decisions into durable guardrails without requiring direct JSON-schema authoring.
## Status
In progress. The policy-free `scan` command, deterministic `pitlord.scan.v1` finding envelope, calibrated language-neutral `dependency-pressure` detector, and initial language-neutral `dependency-knots` detector are implemented. Knot real-corpus calibration, the remaining generalized detector families, disposition workflow, and policy-authoring flow remain planned. This document otherwise describes target behavior rather than shipped detector behavior.

## Problem
Current Pitlord is strongest after architectural intent has already been encoded. `check` deterministically enforces repository-owned rules, `analyze` summarizes declared architecture areas, and `inspect` exposes Arcana architecture communities. Those foundations leave a missing first step:

> Given a repository whose problems are not yet known, identify structurally suspicious architecture and explain where a developer should investigate.

The JSON policy model is also too low-level to be the normal product interface. Users should not need to understand Pitlord's schema, area selectors, relationship arrays, or rule representation to create useful guardrails.
## Product direction
Pitlord should support a complete architecture-governance loop while also acting as an opinionated, zero-configuration complexity guard for coding agents:

```text
agent changes code
  -> deterministic structural analysis
  -> deterministic generalized judgment
  -> pass, or a finite set of bounded failures
  -> agent repairs only those failures
  -> deterministic verification
  -> continue building
```

The product objective is not merely to identify unusual architecture. Pitlord should help prevent repositories from degrading past the point where humans or coding agents can reliably understand and modify them. Its built-in opinions should encode broadly useful clean-code and clean-architecture principles such as avoiding dependency cycles, excessive coupling, weak cohesion, unstable dependency direction, oversized hubs, broad blast radius, incoherent boundaries, and other measurable forms of structural complexity.

The concepts remain distinct:

- **Inspect** exposes structural facts with minimal opinion.
- **Diagnose** applies the generalized opinionated model without requiring repository policy and explains suspicious or unhealthy structure.
- **Generalized guard evaluation** applies the same deterministic model as an immediate pass/fail safeguard for agent and CI workflows without first requiring repository-specific rules.
- **Analyze** explains the repository relative to architecture areas already declared by the user.
- **Check** deterministically enforces explicit repository policy in addition to any selected generalized guard profile.

Repository-specific policy remains the mechanism for encoding local architectural intent. The generalized guard exists before that intent is authored and protects against common complexity failure modes by default.

## Determinism as a product principle
Pitlord must not become another autonomous agent supervising a coding agent. Open-ended review and remediation loops compound probabilistic errors, consume context, and create token churn. Pitlord's value is that the supervisory boundary is deterministic.

The required control loop is:

```text
probabilistic implementation
  -> deterministic observation
  -> deterministic judgment
  -> bounded agent action
  -> deterministic verification
```

For the same repository snapshot, Pitlord version, configuration, and generalized rule profile, Pitlord should produce the same findings, ordering, severity, evidence, fingerprints, and exit status.

A finding must define a stable target that an agent can satisfy rather than invite another round of architectural interpretation. For example, a dependency-cycle finding should identify the exact cycle and require that the cycle cease to exist. The coding agent retains freedom over how to repair the structure; Pitlord decides only whether the measurable condition remains true.

This leads to a hard ownership rule:

> LLMs may consume Pitlord findings and perform repairs. LLMs must not be required for Pitlord to determine whether architecture passes.

Any optional natural-language explanation or remediation assistance must remain downstream of the deterministic finding and must never alter compliance semantics.

## Immediate agent application
The primary near-term consumer of the generalized guard is a coding agent operating in a repository over many successive changes. The guard should minimize the amount of reasoning the agent must spend on self-review by returning a bounded set of actionable failures with exact evidence and a required structural outcome.

A useful machine-facing result should answer four questions deterministically:

1. What condition failed?
2. What exact evidence caused it to fail?
3. What repository scope is implicated?
4. What measurable condition must become false or return within bounds for the check to pass?

Pitlord should not require a second model to decide whether a finding matters, locate the affected structure, or determine whether the repair succeeded. This is the mechanism by which Pitlord reduces autonomous-loop error compounding and token churn while allowing agents to keep building.

## Diagnosis semantics
A diagnosis is a **generalized architectural judgment backed by deterministic evidence**. Some diagnoses remain advisory because Pitlord cannot know that every unusual hub, deep chain, or cross-boundary dependency is wrong. Other detector classes can support an opinionated default failure when the condition is broadly and mechanically undesirable, such as newly introduced dependency cycles.

The generalized engine should therefore distinguish detector disposition from detector evidence:

- **advisory** — structurally suspicious and worth review, but not a default blocking failure;
- **guard** — sufficiently general and mechanically defined to block continued agent work under the selected generalized profile; and
- **repository policy** — explicit local intent authored by the repository and enforced independently of generalized defaults.

Diagnosis should:

- require no repository policy;
- produce deterministic evidence and stable candidate identity where practical;
- explain why a structure is unhealthy or unusual;
- rank relative findings against the repository or comparable peers;
- make blocking versus advisory disposition explicit and deterministic;
- support immediate machine consumption by coding agents; and
- support explicit human disposition where architectural intent is genuinely ambiguous.

Initial detectors should make measurable structural claims rather than vague judgments such as "this code is spaghetti."

### Initial detector set
1. **Dependency knots and cyclic clusters** — strongly connected or mutually dependent regions ranked by size, density, and architectural significance.
2. **Hub and bottleneck candidates** — nodes, files, packages, or communities with anomalously high fan-in, fan-out, or both. The first shipped detector covers file-level cross-file dependency pressure; package/community aggregation remains follow-on work.
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

The first `dependency-pressure` detector now has this calibration seam. Its Arcana-family fixtures preserve the characteristic shapes of all five topology families at the same 64-node scale: modular and layered controls remain quiet; hub-heavy and entangled graphs surface their isolated hub files; and a dense subsystem is not mislabeled as a handful of individual hubs. To preserve that separation, the detector suppresses the isolated-hub diagnosis when more than 12.5% of repository files cross the same pressure threshold; that broader condition belongs to the boundary/cohesion detector.

Live Pitlord calibration also established two important interpretation rules. Explicit Arcana `file` nodes define the repository peer population, while symbol/type/module relationships are aggregated onto those file paths; virtual namespaces such as `@stdlib/*` and package-only paths do not enter the peer set. The first detector judges anomalous outgoing dependency pressure rather than pure fan-in. High fan-in with little or no outgoing coupling is deferred to the impact/blast-radius detector instead of being treated as evidence that a stable shared model should be split.

The completed 16-corpus calibration pass replaced that initial coarse behavior with a boundary-aware production-peer model. Common test, benchmark, sample, generated, vendor, tooling, and devtool paths no longer enter the production peer population. Boundary-aware candidates must remain raw fan-out outliers, rank in the top 5% of active production peers, cross at least three target directories, and have at least 35% cross-directory spread. Conventional composition/entrypoint/controller/invoker/factory seams are deferred, compatibility paths do not escalate above warning, and strongly reused central hubs are left to the bottleneck/blast-radius family. Against the frozen detector-specific references, the current projection scores 3 true positives and 105 true negatives with zero labelled false positives, false negatives, or severity mismatches. Twenty-six findings remain explicitly unlabelled rather than being counted as successes, so dependency pressure remains advisory while those cases and the remaining detector families are completed.

## Initial implementation direction
1. Build on the implemented policy-free `scan` seam and `pitlord.scan.v1` contract by adding the four initial detector families, stable evidence, severity, explicit advisory-versus-guard disposition, and mechanically verifiable required outcomes.
2. Expose the engine through `diagnose` for explanation and through an immediate generalized guard path suitable for post-change agent verification and CI. Do not require repository policy for either path.
3. Make every blocking finding a bounded mechanical target that an agent can repair and Pitlord can re-evaluate without model judgment.
4. Calibrate deterministic thresholds and repository-relative statistics against Arcana's synthetic topology families before real-repository tuning.
5. Add explicit accepted-diagnosis persistence only for genuinely ambiguous advisory structures; do not make acceptance a way to suppress mechanically defined guard failures casually.
6. Connect appropriate diagnosis candidates to human-readable repository-specific policy proposals with current-repository impact preview and safe publication through the existing validation path.

The current `check`, baseline, diff, report, and policy-validation machinery remains the downstream enforcement foundation. No implementation step should introduce an LLM dependency into pass/fail evaluation.

## Acceptance criteria
A developer or coding agent should be able to:

1. point Pitlord at an unfamiliar repository without first writing policy;
2. receive a bounded, evidence-backed set of generalized architectural judgments;
3. distinguish advisory findings from deterministic generalized guard failures;
4. obtain exact evidence, implicated scope, and a measurable passing condition for every blocking finding;
5. repair a finding without requiring another model to reinterpret whether the problem exists or whether it is fixed;
6. rerun Pitlord against the same state and receive reproducible findings, ordering, severity, evidence, and exit status;
7. explicitly accept genuinely intentional advisory structures without abusing baselines;
8. turn an undesirable advisory diagnosis into repository-specific policy without editing JSON;
9. preview what the generated policy means and matches before saving it; and
10. have later Pitlord checks deterministically enforce both generalized guardrails and repository-specific decisions.

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
