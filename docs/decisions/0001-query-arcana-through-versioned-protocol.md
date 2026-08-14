# ADR-0001: Query Arcana through a versioned protocol

## Purpose

This record explains why Pitlord obtains semantic graph evidence through Arcana's versioned query protocol instead of reading Arcana storage directly.

## Overview

Arcana remains the graph-storage and traversal owner. Pitlord sends bounded requests, validates complete protocol responses, and fails closed before policy publication when process or compatibility guarantees fail.

Status: Accepted
Date: 2026-07-31
Owners: Pitlord Arcana client boundary and Arcana query protocol
Supersedes: none
Superseded by: none

## Context

Pitlord needs semantic graph evidence for dependency, ownership, cycle, analysis, diff, and mutation-verification rules. Arcana owns graph storage and traversal. Reading Arcana's packed storage directly would couple Pitlord to an internal format, duplicate validation, and blur graph ownership.

## Decision

Pitlord queries immutable Arcana snapshots by starting the configured Arcana executable and using the exact `arcana.query.v1` JSONL protocol. Pitlord validates every response and fails closed on process, compatibility, pagination, or completeness failure. Repository-only rules remain independent and do not start Arcana.

## Consequences

### Ownership and dependencies

- Arcana owns graph storage, validation, and query semantics.
- Pitlord owns policy semantics and interpretation of returned evidence.
- `internal/arcana` is the only Pitlord process adapter for Arcana.
- Pitlord must not read `.arcana` snapshot storage files directly.

### State, lifecycle, and operations

- Each bounded query process receives one immutable snapshot path.
- Caller context owns cancellation and deadline.
- Child failure prevents a graph-backed compliance result.
- Stderr and protocol error payloads remain diagnostic evidence.

### Compatibility and migration

- Protocol changes require a new protocol identifier or compatible extension proven by contract tests.
- A persistent or pooled Arcana process requires a superseding lifecycle ADR.

## Alternatives considered

- **Read Arcana storage directly:** rejected because it duplicates storage knowledge and bypasses Arcana validation and query ownership.
- **Embed Arcana as a library:** deferred because Pitlord and Arcana are independently usable products with separately evolving implementations.
- **Approximate graph rules from source text after failure:** rejected because it would change rule semantics and could report false compliance.

## Verification

- `internal/arcana` pagination and graph-loading tests.
- Process-boundary tests for cancellation, incompatible protocol IDs, and incomplete response streams.
- Policy tests proving repository-only rules do not require Arcana.

## Risks and debt

- There is no standardized default deadline yet; callers must provide a bounded context where required.
- Process startup cost may justify future pooling, but ownership and cancellation must remain explicit.

## References

- [Arcana process boundary](../ARCANA_PROCESS_BOUNDARY.md)
- [Architecture](../ARCHITECTURE.md)

## Related docs

- [Policy reference](../POLICY.md)
- [Behavioral-contract matrix](../development/behavioral-contract-matrix.md)

## Notes

A future persistent or pooled Arcana integration must preserve immutable snapshot identity, exact protocol validation, cancellation, and fail-closed policy semantics.
