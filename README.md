# Pitlord

Pitlord is the architecture-policy and deterministic diagnosis component of the Warlock toolchain. It evaluates repository-owned rules against source files and immutable Arcana graph snapshots built from Lexicon's polyglot facts. Pitlord does not parse source code itself; language-specific analyzers may normalize native-tool diagnostics into the same finding contract.

## What Pitlord owns

Pitlord currently provides:

- repository content and path policy that can run without Arcana;
- modular JSON policy with named architecture areas, dependency rules, ownership rules, and area-cycle checks;
- deterministic evidence, baselines, snapshot diffs, text/JSON/SARIF output, and embedded policy/baseline schemas;
- policy-free generalized scanning through a registered analyzer engine;
- ten language-neutral graph-backed architecture detectors;
- cross-language semantic rules over Lexicon-normalized facts, currently `swallowed-error` and `unobserved-outcome`;
- graph-free Rust Clippy normalization without taking ownership of Clippy rule semantics;
- declared-area analysis and Arcana architecture-community inspection for policy authoring;
- Demon Docs codemap coverage and changed-code documentation guarding;
- pinned-corpus calibration and exact Homunculus mutation verification.

The core product distinction is deterministic supervision of repository structure and agent-generated changes, not recreating native language lint ecosystems. Architecture detector mechanics live in [Architecture](docs/ARCHITECTURE.md); their evidence maturity and current limits live in [Current limitations](docs/limits/current-limitations.md).

## Build and install

```text
go build ./cmd/pitlord
go test ./...
go vet ./...
```

```text
go install github.com/Lokee86/pitlord/cmd/pitlord@v0.1.2
```

Repository-only content and path policies require only Pitlord. Graph-backed policy, architecture scans, semantic scans, `analyze`, and `inspect` additionally require Arcana and current graph state. `--arcana` is an explicit executable override where supported. Otherwise Pitlord resolves Arcana from `PITLORD_ARCANA_COMMAND`, the repository's current Lexicon configuration, an adjacent Lexicon + Arcana installation, a nearby development checkout, or `PATH`. Retired Grimoire provider state is not part of active discovery.

See [Arcana process boundary](docs/ARCANA_PROCESS_BOUNDARY.md) for the exact discovery and failure contract.

## Quick start

Validate policy without loading repository state:

```text
pitlord validate --policy pitlord.json
```

Enforce policy:

```text
pitlord check \
  --repo /path/to/repository \
  --policy pitlord.json
```

Run Pitlord's policy-free architecture diagnosis:

```text
pitlord scan \
  --repo /path/to/repository \
  --format json
```

Inspect how declared areas map onto the repository graph:

```text
pitlord analyze \
  --repo /path/to/repository \
  --policy pitlord.json \
  --format text
```

Inspect Arcana architecture communities before writing policy:

```text
pitlord inspect \
  --repo /path/to/repository \
  --path-prefix services/game-server
```

## Command model

The four commonly confused surfaces have distinct roles:

- `inspect` discovers repository structure from Arcana communities;
- `analyze` projects declared policy areas and reports structural coverage/relationships;
- `check` enforces explicit repository-owned policy;
- `scan` runs generalized Pitlord/native analyzers without authored policy.

Public commands:

| Command | Purpose |
| --- | --- |
| `check` | Enforce repository policy. |
| `validate` | Validate policy composition and schema semantics. |
| `schema` | Print embedded policy or baseline schema. |
| `baseline` | Freeze current evidence fingerprints. |
| `diff` | Compare policy evidence across immutable snapshots. |
| `scan` | Run policy-free architecture, semantic, and/or native analyzers. |
| `analyze` | Analyze declared architecture areas without enforcing policy. |
| `inspect` | Inspect Arcana architecture-community summaries. |
| `docs` | Enforce codemap coverage and changed-code documentation ownership. |
| `calibrate` | Evaluate one exact built-in detector/analyzer against pinned reference labels. |
| `verify-mutation` | Verify exact Homunculus architecture deltas. |
| `generate` | Convert a Homunculus specimen manifest into Pitlord policy. |
| `version` | Print the Pitlord version. |

See [Command reference](docs/COMMANDS.md) for exact command roles and usage.

## Generalized scan

`pitlord scan` uses registered analyzer groups:

- `architecture` — the default graph-backed group containing `dependency-pressure`, `dependency-knots`, `dependency-depth`, `unstable-dependency-direction`, `boundary-bypass`, `symbol-intermediary-bypass`, `cross-file-intermediary-bypass`, `boundary-cohesion`, `hub-bottleneck`, and `impact-blast-radius`;
- `semantic` — graph-backed Pitlord semantic rules over Lexicon-normalized capabilities;
- `clippy` — graph-free Rust Clippy normalization.

Groups can be combined:

```text
pitlord scan \
  --repo /path/to/repository \
  --analyzers architecture,semantic,clippy \
  --format json
```

All current architecture detectors and semantic rules are advisory. Pitlord fails closed when a selected analyzer's required graph state or normalized semantic capabilities are unavailable.

For detailed detector thresholds and suppression semantics, see [Architecture](docs/ARCHITECTURE.md). For calibration strength, unsupported claims, and current evidence gaps, see [Current limitations](docs/limits/current-limitations.md).

## Policy

A minimal architecture policy can declare areas and enforce relationships between them:

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

See [Policy reference](docs/POLICY.md) for the full contract.

## Baselines and snapshot diffs

Create a baseline from current findings:

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

Compare the same policy across two immutable snapshots:

```text
pitlord diff \
  --before-snapshot /path/to/before/.arcana/snapshots/<digest> \
  --repo /path/to/after \
  --policy pitlord.json
```

See [Baselines and CI](docs/BASELINES-AND-CI.md) for fingerprints, SARIF, and CI behavior.

## Documentation guard

`pitlord docs` consumes Demon Docs resolved codemap ownership, Arcana's repository code-file inventory, and Git change evidence. Pitlord owns code-file coverage and changed-code-to-owning-document checks; Demon Docs continues to own Markdown structure, schema, indexes, links, and navigation.

```text
pitlord docs \
  --repo /path/to/repository \
  --changed-from origin/main \
  --format json
```

## Calibration and mutation evidence

`pitlord calibrate` runs exactly one built-in detector or analyzer against a pinned machine-readable reference. References can constrain rule, language, exact source location, severity, and fully-labelled regression behavior.

```text
pitlord calibrate \
  --repo /path/to/pinned/corpus \
  --reference testdata/calibration/jsoup.json \
  --fail-on-mismatch \
  --format json
```

Homunculus mutation manifests can be converted to policy with `generate` or verified directly with `verify-mutation`. See [Calibration harness](docs/development/calibration/HARNESS.md) and [Architecture](docs/ARCHITECTURE.md) for the evidence contracts.

## Exit codes

For commands that enforce or compare findings:

- `0` — success with no blocking finding or mismatch;
- `1` — unsuppressed/introduced diagnostics or expectation mismatch;
- `2` — usage, configuration, policy, snapshot, protocol, or execution failure.

Commands that only inspect or emit metadata use the applicable subset of these statuses.

## Ownership boundaries

- Lexicon owns parsing, semantic resolution, normalized facts, identities, and spans.
- Arcana owns graph ingestion, immutable snapshots, adjacency, traversal, and graph algorithms.
- Pitlord owns policy, area projection, generalized judgments, evidence identity, calibration scoring, and presentation.
- Demon Docs owns document parsing, documentation structure/schema, and codemap authoring/resolution.
- Homunculus owns deterministic source mutations and expected architecture deltas.

## Documentation

- [Documentation index](docs/INDEX.md)
- [Command reference](docs/COMMANDS.md)
- [Policy reference](docs/POLICY.md)
- [Architecture](docs/ARCHITECTURE.md)
- [Arcana process boundary](docs/ARCANA_PROCESS_BOUNDARY.md)
- [Baselines and CI](docs/BASELINES-AND-CI.md)
- [Current limitations](docs/limits/current-limitations.md)
- [Roadmap](docs/planning/roadmap.md)

## License

Pitlord is available under the [PolyForm Shield License 1.0.0](LICENSE.md). Competing products and services require a separate commercial license. See [LICENSING.md](LICENSING.md).
