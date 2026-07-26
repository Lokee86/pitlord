# Baselines and CI

Pitlord baselines operate at evidence level. A baseline does not disable a rule or hide an
entire diagnostic class; it suppresses only previously accepted source evidence while new
evidence under the same rule still fails the check.

## Create a baseline

```text
pitlord baseline \
  --repo . \
  --policy pitlord.json \
  --output pitlord.baseline.json
```

The baseline records one entry per current evidence item with:

- a SHA-256 evidence fingerprint;
- rule and issue identity;
- source and optional target paths;
- relation name;
- ownership or cycle areas; and
- source and target area identities for area-cycle evidence.

## Fingerprint stability

Fingerprints prefer Arcana's stable semantic node identity. When identity is unavailable,
Pitlord falls back to the node key or source location. This allows findings to survive
renames when Lexicon and Arcana preserve semantic identity, without making broad fuzzy
matches that could suppress unrelated evidence.

Rule ID, issue class, relation, target identity, component areas, and area-edge identity
are also part of the fingerprint. Changing the architectural meaning creates a new finding.

## Check against a baseline

```text
pitlord check \
  --repo . \
  --policy pitlord.json \
  --baseline pitlord.baseline.json
```

Exit code `0` means no unsuppressed findings remain. Text and JSON summaries report how
many baseline findings were suppressed.

Validate a baseline independently:

```text
pitlord validate --baseline pitlord.baseline.json
pitlord schema --kind baseline
```

## SARIF

```text
pitlord check \
  --repo . \
  --policy pitlord.json \
  --format sarif
```

Pitlord emits SARIF 2.1.0 with:

- one result per evidence item;
- the Pitlord rule ID and severity;
- source location and span;
- related target location for dependency evidence;
- issue, relation, and area properties; and
- a stable `pitlordEvidenceFingerprint/v1` partial fingerprint.

The command still exits `1` when results exist. CI should capture the SARIF output before
propagating or interpreting the exit status.

## Snapshot-to-snapshot gating

A repository can gate only newly introduced evidence without checking in a baseline file:

```text
pitlord diff \
  --before-snapshot /path/to/base/.arcana/snapshots/<digest> \
  --repo . \
  --policy pitlord.json
```

Evidence is classified by the same fingerprint contract as baselines:

- introduced: present only after;
- resolved: present only before; and
- persistent: present in both snapshots.

The command exits `1` only when introduced findings exist. Text and JSON include all three
classes; SARIF contains introduced findings only.

## Suggested CI sequence

1. Build or obtain the current Lexicon snapshot.
2. Synchronize the Arcana snapshot.
3. Run `pitlord validate --policy pitlord.json`.
4. Run `pitlord check` with the repository baseline, if adopted.
5. Optionally compare against the merge-base Arcana snapshot with `pitlord diff`.
6. Optionally emit a SARIF run for code-scanning upload.

Pitlord evaluates immutable Arcana snapshots. CI should therefore record the snapshot
identity alongside the policy and baseline revision used for the result.
