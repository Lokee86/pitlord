# Pitlord Policy v1

A Pitlord policy is a JSON document with `version`, optional relative `includes`, optional
named `areas`, and one or more `rules`. Unknown fields are rejected. Pitlord normalizes
paths, names, kinds, globs, and relations before evaluation and rejects invalid or
ineffective rules.

## Policy composition

A policy may include other policy files relative to its own directory:

```json
{
  "version": 1,
  "includes": [
    "areas/game-server.json",
    "rules/repository.json"
  ],
  "rules": [
    {"id": "readme-required", "type": "require_path", "path": "README.md"}
  ]
}
```

Includes are loaded recursively, deduplicated by resolved path, and flattened before
validation. Duplicate area or rule IDs remain errors. Include cycles are rejected.

## Repository content and path rules

Repository rules do not require Lexicon, Arcana, or a graph snapshot.

`forbid_content` scans matching files line by line for either a literal or regular
expression:

```json
{
  "id": "no-global-random",
  "type": "forbid_content",
  "literal": "math/rand",
  "include_paths": ["server/**/*.go"],
  "exclude_paths": ["server/rng/**", "server/**/*_test.go"]
}
```

Use exactly one of `literal`, `regex`, or the compatibility pair `pattern` plus optional
`pattern_type`. Repository-relative globs support `*`, `?`, character classes, and `**`.
Pitlord skips nested worktrees, dependency caches, editor caches, and Warlock tool-state
directories.

`require_path` requires at least one matching file or directory. `forbid_path` rejects
every matching file or directory:

```json
{"id": "control-required", "type": "require_path", "path": "server/control.go"}
```

```json
{"id": "no-legacy-debug-files", "type": "forbid_path", "path": "server/debug_*.go"}
```

## Areas

An area is a reusable ownership and dependency selector:

```json
{
  "id": "storage",
  "description": "Persistence boundary",
  "paths": ["server/storage"],
  "exclude_paths": ["server/storage/testdata"],
  "kinds": ["file", "function", "type"]
}
```

`paths` are repository-relative prefixes. A prefix matches itself and descendants. `.`
matches the repository root. `exclude_paths` remove matches. An empty `kinds` list accepts
all Arcana node kinds.

Area IDs must be unique. A node may match zero, one, or several areas; ownership rules
control whether those states are permitted.

## Forbidden dependencies

`forbid_dependency` reports direct normalized Arcana relationships from a source selector
to a target selector.

```json
{
  "id": "api-must-not-access-storage",
  "type": "forbid_dependency",
  "severity": "error",
  "message": "api must use the service boundary",
  "from_areas": ["api"],
  "to_areas": ["storage"],
  "relations": ["imports", "calls"]
}
```

Source and target selectors may use named areas, direct paths, or both:

- `from_areas`, `to_areas`;
- `from_paths`, `to_paths`;
- `from_exclude_paths`, `to_exclude_paths`;
- `source_kinds`, `target_kinds`; and
- `relations`.

At least one source selector and one target selector are required. If `relations` is
omitted, Pitlord uses the supported dependency-like normalized relations.

Each rule produces one diagnostic containing deterministic evidence for every matching
source-relation-target edge.

## Required ownership

`require_ownership` checks nodes in a scope against all declared areas.

```json
{
  "id": "source-files-need-one-owner",
  "type": "require_ownership",
  "scope_paths": ["server"],
  "scope_exclude_paths": ["server/generated"],
  "source_kinds": ["file"]
}
```

By default, the rule reports:

- nodes that match no declared area; and
- nodes that match more than one declared area.

`source_kinds` defaults to `file`. `allow_unowned` or `allow_overlaps` may permit one
failure class. Enabling both is rejected because the rule would enforce nothing.

Ownership-only policies load node membership but do not request outgoing relationships
from Arcana.

## Forbidden area cycles

`forbid_area_cycles` projects symbol relationships onto the selected area graph and
reports strongly connected components containing at least two areas.

```json
{
  "id": "areas-must-be-acyclic",
  "type": "forbid_area_cycles",
  "areas": ["api", "service", "storage"],
  "relations": ["imports", "calls"]
}
```

If `areas` is omitted, all declared areas are selected. At least two selected areas are
required. `source_kinds`, `target_kinds`, and `relations` can narrow the projected graph.

Pitlord reports one deterministic representative source edge for each directed area pair
inside a cyclic component. Self-area relationships do not create area cycles.

## Severity and messages

Every rule accepts:

- `severity`: `error` or `warning`, defaulting to `error`; and
- `message`: a custom diagnostic message.

Pitlord supplies a deterministic default message when one is omitted.

## Normalized relations

Pitlord consumes Lexicon/Arcana relation names. Common dependency-like relations include:

- `imports`;
- `depends-on`;
- `calls`;
- `references`;
- `extends`;
- `implements`;
- `includes`;
- `routes-to`;
- `communicates-with`; and
- `observed-calls`.

A policy should name only relations supported by the active language adapters and Arcana
snapshot.

## Validation and schema

```text
pitlord validate --policy pitlord.json
pitlord schema --kind policy
```

The embedded schema is Draft 2020-12. Runtime validation remains authoritative because it
also resolves includes and checks cross-references, defaults, glob validity, regular
expressions, and semantic no-op conditions.
