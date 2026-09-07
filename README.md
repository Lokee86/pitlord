# Pitlord

Pitlord is the architecture-policy and diagnosis component of the Warlock toolchain.
It evaluates repository-owned rules against immutable Arcana graph snapshots built from
Lexicon's polyglot language facts. Pitlord does not parse source code itself; language-specific
analyzers may also normalize diagnostics from native tooling into Pitlord's finding contract.

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
- a Rust `clippy` analyzer that runs Cargo Clippy in machine-readable mode, deduplicates repeated target diagnostics, and normalizes Clippy rule IDs, source spans, severity, and structured suggestions including machine-applicable edits;
- a graph-backed `semantic` analyzer group whose rules declare language-neutral capabilities supplied by Lexicon adapters; `swallowed-error` reports handlers with neither a local error action nor a proven downstream fallback, enclosing propagation, or explicit suppression disposition, while `unobserved-outcome` reports proven fallible/async outcomes that are discarded without consumption across Rust, TypeScript, JavaScript, and Python;
- a language-neutral dependency-pressure detector that aggregates normalized Arcana relationships by repository file path, separates production from common test/benchmark/tooling peers, requires extreme outgoing fan-out plus cross-region boundary spread, and defers conventional composition seams and highly reused central hubs to the detector families that own those roles;
- a calibrated language-neutral `dependency-knots` detector that finds strongly connected production-file components over static source-dependency relations, uses semantic namespace/module regions before directory fallback, excludes runtime calls from source-cycle direction, and uses density only as a narrowed single-region fallback;
- a calibrated advisory language-neutral `dependency-depth` detector that finds narrow acyclic behavioral dependency corridors spanning at least four architectural regions, excludes import/reference-only propagation, and stops depth inheritance at cycles, shared convergence points, and highly reused central foundations;
- an advisory language-neutral `unstable-dependency-direction` detector that applies a region-level Stable Dependencies Principle to static source dependencies, requiring an incoming-heavy source region to depend on an outgoing-heavy target region across a material instability gap while deferring region cycles to dependency-knot analysis;
- an advisory language-neutral `boundary-bypass` detector that finds exceptional file-level behavioral shortcuts around an established cross-region intermediary while suppressing shared targets, weak gateway concentration, import/reference-only reach, and nested ownership trees;
- an advisory `symbol-intermediary-bypass` detector that finds direct symbol calls bypassing an abandoned one-hop local wrapper only when at least three sibling wrapper paths and a shared support step establish the repeated behavioral pipeline; constructor targets and continued same-name overload routes remain quiet;
- an advisory `cross-file-intermediary-bypass` detector that finds same-region direct calls around an established short one-hop intermediary in another file only when a same-named peer still routes through that intermediary and the downstream target has narrow caller-file ownership;
- an advisory language-neutral `boundary-cohesion` detector that evaluates semantic namespace/module regions for broad independent target-region spread combined with both relatively and absolutely weak internal collaboration, while deferring incoming-heavy regions to the detector family that owns central coordination;
- a calibrated advisory language-neutral `hub-bottleneck` detector that requires extreme direct fan-in plus a many-to-many behavioral waist: `calls`/`reads`/`writes` must arrive from at least six architectural regions, leave toward at least four regions, and preserve at least 40% outgoing-to-incoming behavioral flow;
- a calibrated advisory language-neutral `impact-blast-radius` detector with two evidence lanes: repository-relative broad direct change exposure backed by non-metadata dependency evidence, or unusually broad transitive dependent surfaces that expand through at least three substantial independent first-hop branches with no single branch owning more than 80% of the closure;
- a `calibrate` evaluation harness that verifies pinned corpus revisions plus optional exact worktree-diff hashes and scores one exact built-in detector/analyzer target against machine-readable frozen reference labels, including optional rule/language/source-location selectors and strict fully-labelled regression mode;
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
pitlord docs
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

`pitlord scan` is the policy-free generalized-guard and analyzer entry point. The default `--analyzers architecture` selection runs the ten Arcana-backed architecture detectors. `--analyzers semantic` runs Pitlord-native semantic rules over normalized Lexicon capability facts stored in Arcana; `swallowed-error` and `unobserved-outcome` operate across Rust, TypeScript, JavaScript, and Python without containing any of those languages' syntax rules. `--analyzers clippy` runs Rust Clippy without requiring Arcana, and comma-separated groups such as `architecture,semantic,clippy` combine result streams under the same finding contract. Clippy findings preserve the upstream `clippy::...` rule ID, repository-relative source span, and structured suggested edits when rustc marks a suggestion as applicable. `dependency-pressure` reports unusually broad outgoing cross-file dependency pressure using Arcana's language-neutral normalized relationship taxonomy. It compares production peers, requires a top-5% outgoing outlier to cross at least three target regions with at least 35% boundary spread, keeps conventional entrypoint/composition/controller/invoker/factory seams out of that detector, and defers heavily reused central hubs to the detector family that owns them. `dependency-knots` separately projects static source-dependency relations such as imports, inheritance, implementation, trait use, overrides, includes, depends-on, and conversions into a production file graph; runtime calls are excluded so callback dispatch cannot manufacture reverse source ownership. It finds strongly connected components deterministically and reports significant cyclic components with a mechanical requirement that mutual reachability be broken. `dependency-depth` analyzes acyclic behavioral dependency flow over calls, inheritance/implementation, trait use, overrides, includes, explicit dependencies, and conversions. A finding requires a narrow dominant corridor of at least four hops spanning at least four architectural regions with at least three region transitions. Import-only and reference-only edges do not propagate depth, and cycles, shared convergence points, or highly reused central foundations terminate the walk rather than transferring downstream depth to callers. `unstable-dependency-direction` evaluates static source dependency direction between architectural regions using `I = outgoing regions / (incoming + outgoing regions)`. It requires at least three files and three total region couplings on both sides, at least two source files supporting the boundary, an incoming-heavy source, an outgoing-heavy target, and an instability increase of at least 0.25; region-level strongly connected components are deferred to `dependency-knots`. `boundary-bypass` detects an exceptional direct behavioral edge from an upstream region into a downstream region only when at least three upstream peers establish the same gateway, direct access remains a minority, the gateway is materially concentrated on the downstream region, the target is not already a broadly shared foundation, and source/target files do not merely occupy one nested ownership tree. Import/reference-only reach does not establish a bypass. `symbol-intermediary-bypass` separately operates on resolved symbol-level calls: a direct target call is reportable only when an abandoned one-hop local wrapper is supported by at least three intact sibling wrapper paths and a shared non-target support step; constructors and continued same-name overload routes remain quiet. `cross-file-intermediary-bypass` owns a narrower cross-file shape: caller, intermediary, and target must remain inside one semantic region; the intermediary must be a short one-hop wrapper with at least two remaining callers; a same-named peer in the intermediary file must still route through it; and the downstream target may be called from at most two files after the shortcut. Same-named forwarding chains and broad shared utilities are suppressed. `boundary-cohesion` evaluates semantic namespace/module regions with six or more production files, requiring broad independent target-region reach plus internal collaboration that is both in the weakest peer quartile and at or below 0.5 support relationships per member; runtime calls and reads/writes may support internal cohesion even though they do not define dependency-knot direction. `hub-bottleneck` then owns central coordination: candidates must already be extreme direct fan-in outliers, but are reported only when behavioral flow (`calls`, `reads`, `writes`) enters from at least six architectural regions, exits toward at least four regions, and outgoing behavioral degree is at least 40% of incoming behavioral degree. This suppresses shared contracts, result/domain models, base abstractions, and narrow cross-cutting facades whose centrality is not a many-to-many coordination waist. `impact-blast-radius` owns change exposure through two lanes. Direct exposure requires at least `max(6, ceil(3% of production peers))` direct dependents, at least two incoming architectural regions, and at least two dependents supported by non-metadata relationships; normal-sized scopes require top-5% fan-in unless the surface already reaches at least 20% of peers or at least 10 regions, and final breadth requires at least 5% peer reach or 10 regions. Files that do not qualify directly may qualify through transitive amplification: at least 25% peer reach, top-5% active impact, at least 2x amplification, at least three impacted regions, at least three substantial independent first-hop branches, and at most 80% dominant-branch share. This keeps metadata-only centrality and inherited single-gateway closure quiet while still identifying real shared-foundation change exposure. All ten detectors are advisory. Dependency knots, dependency depth, hub/bottleneck, impact/blast-radius, symbol intermediary bypass, and cross-file intermediary bypass have required real-repository positives. Both intermediary-bypass positives are diff-pinned Homunculus mutations used during detector development, so neither supplies independent validation/holdout recall evidence. Dependency depth's independently audited 16-corpus projection scores 1 TP / 16 TN / 0 FP / 0 FN with no severity mismatches or unlabelled findings; impact's corrected projection scores 19 TP / 52 TN / 0 FP / 0 FN. Unstable dependency direction and boundary bypass each score 0 TP / 16 TN / 0 FP / 0 FN against independently audited false-positive controls, while boundary/cohesion also relies on frozen real-corpus negative controls plus synthetic positive coverage; real-world recall is not yet claimed for those three families:

```text
pitlord scan \
  --repo /path/to/repository \
  --path-prefix services/game-server \
  --format json
```

Run the cross-language semantic lint group against an Arcana snapshot built from current Lexicon facts:

```text
pitlord scan \
  --repo /path/to/repository \
  --analyzers semantic \
  --format json
```

Run only Rust Clippy through the same scan contract, with no Arcana snapshot required:

```text
pitlord scan \
  --repo /path/to/rust-repository \
  --analyzers clippy \
  --format json
```

`pitlord docs` is the blocking documentation guard. Demon Docs remains authoritative for codemap parsing and resolution: Pitlord consumes the schema-1 dataset from `ddocs codemaps export`, while Arcana supplies the repository code-file inventory. Every Arcana code file must resolve from at least one authored codemap entry in the selected Demon Docs root. With `--changed-from`, every changed mapped code file must also have at least one of its owning codemap documents changed in the same Git range. Directory-only codemap entries do not satisfy per-file coverage; exact or symbol-backed file resolutions and resolved glob matches do. Pitlord does not validate Markdown structure, document schemas, indexes, or navigation; those remain Demon Docs CI responsibilities.

```text
pitlord docs \
  --repo /path/to/repository \
  --changed-from origin/main \
  --format json
```

`pitlord calibrate` runs the exact built-in detector or analyzer named by a frozen corpus reference. Legacy `detector` references remain valid; analyzer references can additionally constrain `rule_id`, `language`, and exact source location. The command verifies the repository Git revision by default and, when supplied, the exact SHA-256 of the tracked worktree diff. It distinguishes required findings, findings that must be absent, allowed findings, severity mismatches, and unlabelled target output:

```text
pitlord calibrate \
  --repo /path/to/pinned/corpus \
  --reference docs/development/calibration/references/jsoup.json \
  --format json
```

Use `--fail-on-mismatch` when the reference is expected to pass. References with `require_fully_labelled: true` also fail when the selected detector/analyzer produces an unexpected unlabelled finding. Canonical exploratory detector-tuning runs may omit that flag so current false positives and false negatives can be measured without aborting the batch.

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

- Lexicon owns parsing, semantic resolution, normalized relationships, semantic capabilities, error-handler/action facts, downstream error-flow facts, outcome-obligation facts, identities, and spans.
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
