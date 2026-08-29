# Pitlord

Pitlord is the architecture-policy and diagnosis component of the Warlock toolchain.
It evaluates repository-owned rules against immutable Arcana graph snapshots built from
Lexicon's polyglot language facts. Pitlord does not maintain language adapters or parse
source code itself.

## Current capabilities

Pitlord currently provides:

- repository forbidden-content, required-content, and required/forbidden path rules that run without Arcana;
- modular policies composed through relative `includes`;
- named architecture areas with path, exclusion, and symbol-kind selectors;
- forbidden direct dependencies between paths or areas;
- exact ownership checks for unowned and multiply owned nodes;
- cycle detection over the declared area-dependency graph;
- deterministic source and target evidence with Arcana spans and identities;
- evidence-level baselines that suppress existing findings without hiding new ones;
- text, JSON, and SARIF 2.1.0 output;
- architecture-community inspection through Arcana;
- a policy-free `scan` command with the deterministic `pitlord.scan.v1` finding envelope, repository-relative scope, severity counts, and text/JSON output;
- a language-neutral dependency-pressure detector that aggregates normalized Arcana relationships by repository file path, separates production from common test/benchmark/tooling peers, requires extreme outgoing fan-out plus cross-region boundary spread, and defers conventional composition seams and highly reused central hubs to the detector families that own those roles;
- a calibrated language-neutral `dependency-knots` detector that finds strongly connected production-file components over static source-dependency relations, uses semantic namespace/module regions before directory fallback, excludes runtime calls from source-cycle direction, and uses density only as a narrowed single-region fallback;
- an advisory language-neutral `boundary-cohesion` detector that evaluates semantic namespace/module regions for broad independent target-region spread combined with both relatively and absolutely weak internal collaboration, while deferring incoming-heavy regions to the detector family that owns central coordination;
- a calibrated advisory language-neutral `hub-bottleneck` detector that requires extreme direct fan-in plus a many-to-many behavioral waist: `calls`/`reads`/`writes` must arrive from at least six architectural regions, leave toward at least four regions, and preserve at least 40% outgoing-to-incoming behavioral flow;
- a calibrated advisory language-neutral `impact-blast-radius` detector with two evidence lanes: repository-relative broad direct change exposure backed by non-metadata dependency evidence, or unusually broad transitive dependent surfaces that expand through at least three substantial independent first-hop branches with no single branch owning more than 80% of the closure;
- a `calibrate` evaluation harness that verifies pinned corpus revisions and scores detector output against machine-readable frozen reference labels without treating unlabelled findings as ground truth;
- direct Homunculus specimen-manifest conversion and expected-diagnostic validation;
- exact Homunculus mutation verification across baseline and mutated snapshots;
- snapshot-to-snapshot policy diffing with introduced, resolved, and persistent evidence;
- embedded JSON Schemas for policy and baseline files; and
- paginated Arcana graph loading for large repositories.

The implementation has been exercised against deterministic Homunculus specimens and a
64,069-node Space Rocks Arcana snapshot.

## Build and verify

```text
go build ./cmd/pitlord
go test ./...
go vet ./...
```

Repository-only content and path policies require only the Pitlord executable. Policies
containing dependency, ownership, or area-cycle rules additionally require an Arcana
executable and current snapshot. `--arcana` is an explicit executable override. Otherwise
Pitlord prefers the Arcana associated with the repository's prepared Grimoire/Lexicon
state, then checks an adjacent/shared Grimoire installation and `GRIMOIRE_HOME`, and uses
`PATH` only as the final fallback.

## Install

```text
go install github.com/Lokee86/pitlord/cmd/pitlord@v0.1.2
```

## Quick start

Validate a policy without loading a repository:

```text
pitlord validate --policy pitlord.json
```

Check a repository policy:

```text
pitlord check \
  --repo /path/to/repository \
  --policy pitlord.json
```

Content and path rules scan the repository directly. Graph rules also resolve
`.arcana/CURRENT`. An immutable snapshot directory may instead be provided explicitly:

```text
pitlord check \
  --snapshot /path/to/repository/.arcana/snapshots/<digest> \
  --policy pitlord.json \
  --format json
```

Exit codes:

- `0`: no unsuppressed diagnostics and any Homunculus expectation matched;
- `1`: unsuppressed diagnostics exist or a Homunculus expectation mismatched;
- `2`: usage, configuration, policy, snapshot, or Arcana execution failure.

## Policy example

```json
{
  "version": 1,
  "areas": [
    {"id": "api", "paths": ["api"]},
    {"id": "service", "paths": ["service"]},
    {"id": "storage", "paths": ["storage"]}
  ],
  "rules": [
    {
      "id": "api-must-not-access-storage",
      "type": "forbid_dependency",
      "from_areas": ["api"],
      "to_areas": ["storage"],
      "relations": ["imports", "depends-on", "calls"]
    },
    {
      "id": "source-must-have-one-owner",
      "type": "require_ownership",
      "scope_paths": ["."],
      "source_kinds": ["file"]
    },
    {
      "id": "areas-must-be-acyclic",
      "type": "forbid_area_cycles",
      "relations": ["imports", "depends-on", "calls"]
    }
  ]
}
```

See [docs/POLICY.md](docs/POLICY.md) for the full policy contract.

## Commands

```text
pitlord check
pitlord baseline
pitlord validate
pitlord schema
pitlord verify-mutation
pitlord diff
pitlord scan
pitlord calibrate
pitlord inspect
pitlord generate
pitlord version
```

Use `pitlord schema --kind policy` or `--kind baseline` to print the embedded Draft
2020-12 JSON Schema.

`pitlord inspect` exposes Arcana's architecture-community summary for policy authoring:

```text
pitlord inspect \
  --repo /path/to/repository \
  --path-prefix services/game-server \
  --relations calls,imports,references
```

`pitlord scan` is the policy-free generalized-guard entry point. `dependency-pressure` reports unusually broad outgoing cross-file dependency pressure using Arcana's language-neutral normalized relationship taxonomy. It compares production peers, requires a top-5% outgoing outlier to cross at least three target regions with at least 35% boundary spread, keeps conventional entrypoint/composition/controller/invoker/factory seams out of that detector, and defers heavily reused central hubs to the detector family that owns them. `dependency-knots` separately projects static source-dependency relations such as imports, inheritance, implementation, trait use, overrides, includes, depends-on, and conversions into a production file graph; runtime calls are excluded so callback dispatch cannot manufacture reverse source ownership. It finds strongly connected components deterministically and reports significant cyclic components with a mechanical requirement that mutual reachability be broken. `boundary-cohesion` evaluates semantic namespace/module regions with six or more production files, requiring broad independent target-region reach plus internal collaboration that is both in the weakest peer quartile and at or below 0.5 support relationships per member; runtime calls and reads/writes may support internal cohesion even though they do not define dependency-knot direction. `hub-bottleneck` then owns central coordination: candidates must already be extreme direct fan-in outliers, but are reported only when behavioral flow (`calls`, `reads`, `writes`) enters from at least six architectural regions, exits toward at least four regions, and outgoing behavioral degree is at least 40% of incoming behavioral degree. This suppresses shared contracts, result/domain models, base abstractions, and narrow cross-cutting facades whose centrality is not a many-to-many coordination waist. `impact-blast-radius` owns change exposure through two lanes. Direct exposure requires at least `max(6, ceil(3% of production peers))` direct dependents, at least two incoming architectural regions, and at least two dependents supported by non-metadata relationships; normal-sized scopes require top-5% fan-in unless the surface already reaches at least 20% of peers or at least 10 regions, and final breadth requires at least 5% peer reach or 10 regions. Files that do not qualify directly may qualify through transitive amplification: at least 25% peer reach, top-5% active impact, at least 2x amplification, at least three impacted regions, at least three substantial independent first-hop branches, and at most 80% dominant-branch share. This keeps metadata-only centrality and inherited single-gateway closure quiet while still identifying real shared-foundation change exposure. All five detectors are advisory. Dependency knots, hub/bottleneck, and impact/blast-radius now have required real-corpus positives; impact's corrected 16-corpus projection scores 19 TP / 52 TN / 0 FP / 0 FN with no severity mismatches or unlabelled findings. Boundary/cohesion still relies on frozen real-corpus negative controls plus synthetic positive coverage, so real-world recall is not yet claimed for that detector family:

```text
pitlord scan \
  --repo /path/to/repository \
  --path-prefix services/game-server \
  --format json
```

`pitlord calibrate` runs the selected generalized detector against a frozen corpus reference. It verifies the repository Git revision by default and distinguishes required findings, findings that must be absent, allowed maintenance-watch findings, severity mismatches, and unlabelled detector output:

```text
pitlord calibrate \
  --repo /path/to/pinned/corpus \
  --reference docs/development/calibration/references/jsoup.json \
  --format json
```

Use `--fail-on-mismatch` when the reference is expected to pass. Canonical detector-tuning runs omit that flag so current false positives and false negatives can be measured without aborting the batch.

## Baselines and CI

Create a baseline from the current findings:

```text
pitlord baseline \
  --repo /path/to/repository \
  --policy pitlord.json \
  --output pitlord.baseline.json
```

Then fail only on new evidence:

```text
pitlord check \
  --repo /path/to/repository \
  --policy pitlord.json \
  --baseline pitlord.baseline.json
```

For code-scanning systems:

```text
pitlord check --repo . --policy pitlord.json --format sarif
```

See [docs/BASELINES-AND-CI.md](docs/BASELINES-AND-CI.md) for fingerprint and CI behavior.

## Snapshot diffs

Compare the same policy across two immutable Arcana snapshots:

```text
pitlord diff \
  --before-snapshot /path/to/before/.arcana/snapshots/<digest> \
  --repo /path/to/after \
  --policy pitlord.json
```

The diff classifies evidence as introduced, resolved, or persistent using the same stable
evidence fingerprints as baselines. It exits `1` only when introduced evidence exists.
`--format sarif` emits only introduced findings, making it suitable for pull-request gates.

## Homunculus integration

Generated specimens can be checked directly:

```text
pitlord check \
  --repo /path/to/specimen \
  --homunculus-manifest /path/to/specimen/homunculus.manifest.json
```

Architecture areas and forbidden edges in the specimen manifest become Pitlord policy.
Expected diagnostic IDs are compared with the actual result. A reusable policy can also
be emitted:

```text
pitlord generate \
  --homunculus-manifest /path/to/specimen/homunculus.manifest.json \
  --output pitlord.json
```

Real-repository mutation contracts can be verified exactly against matched snapshots:

```text
pitlord verify-mutation \
  --manifest /path/to/mutated/.homunculus/architecture-mutation.json \
  --baseline-snapshot /path/to/baseline/.arcana/snapshots/<digest> \
  --repo /path/to/mutated/source-root
```

For every expected added relationship, Pitlord requires absent→present. For every expected
removed relationship, it requires present→absent. Qualified Lexicon `path::name` symbols
are resolved through Arcana rather than matched by source text.

## Ownership boundaries

- Lexicon owns parsing, semantic resolution, normalized relationships, identities, and spans.
- Arcana owns graph ingestion, immutable snapshots, adjacency, traversal, and graph algorithms.
- Pitlord owns policy, area projection, generalized scan judgments, rule evaluation, diagnostic identity, and presentation.
- Homunculus owns deterministic source mutations and their expected architecture deltas.

See [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) for the runtime and data-flow design.

## Documentation

- [Documentation index](docs/INDEX.md)
- [Policy reference](docs/POLICY.md)
- [Architecture](docs/ARCHITECTURE.md)
- [Baselines and CI](docs/BASELINES-AND-CI.md)
- [Documentation coverage](docs/development/documentation-coverage.md)
- [Current limitations](docs/limits/current-limitations.md)
- [Roadmap](docs/planning/roadmap.md)

## License

Pitlord is available under the [PolyForm Shield License 1.0.0](LICENSE.md). Competing products and services require a separate commercial license. See [LICENSING.md](LICENSING.md).
