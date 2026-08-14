# Arcana Process Boundary

Parent index: [Pitlord Documentation](INDEX.md)

## Purpose

This document defines process ownership, deadlines, protocol compatibility, failure behavior, and observability when Pitlord queries Arcana.

## Overview

Pitlord starts Arcana as a bounded child process only for policy operations that require graph evidence. It sends JSONL requests through `arcana.query.v1`, reads JSONL responses, validates the complete response set, and terminates the child through the caller's context when evaluation is cancelled or its deadline expires.

Repository-only rules do not start Arcana.

## Ownership

Pitlord owns:

- deciding whether selected policy requires graph evidence;
- resolving the immutable snapshot passed to Arcana;
- constructing bounded requests and stable request IDs;
- supplying cancellation and deadlines through `context.Context`;
- validating protocol identity, response identity, pagination, completeness, and operation errors;
- converting process and protocol failure into a failed Pitlord evaluation.

Arcana owns:

- opening and validating its snapshot;
- query semantics, graph traversal, pagination results, and protocol response payloads;
- its process-local resources and clean exit after stdin completion or cancellation.

The shell owns neither boundary. Pitlord invokes the configured Arcana executable directly with explicit arguments.

## Deadline and cancellation contract

`internal/arcana` uses `exec.CommandContext`. The immediate caller owns the context and therefore the operation deadline. Cancellation terminates the Arcana child and the evaluation returns an error; Pitlord does not continue with partial graph evidence.

A caller that supplies an unbounded background context accepts an unbounded process wait. CLI composition and future Warlock integration should provide a deadline appropriate to repository size and operation type. A product-wide default timeout has not yet been standardized.

## Protocol contract

Pitlord requires the exact protocol identifier `arcana.query.v1`. It rejects:

- another protocol identifier;
- malformed JSON;
- duplicate or missing logical response IDs;
- fewer or more complete responses than requested;
- failed operations without an error payload;
- invalid pagination continuity or terminal state;
- truncated legacy responses that provide no continuation offset.

Pitlord does not infer compatibility from an Arcana executable version string. Protocol compatibility is established by successful exact request/response validation.

## Failure and degradation

Graph-backed policy evaluation fails closed when Arcana cannot start, is cancelled, exits unsuccessfully, emits incompatible or incomplete output, or rejects an operation. Stderr is included in the returned process error when present.

Pitlord does not silently downgrade a graph rule into a repository-only approximation. Independent repository-only rules can still be run as a separate evaluation that selects no graph-backed rules, but one mixed evaluation does not publish a partial compliant result after graph failure.

## Observability

Operational errors should preserve:

- configured Arcana command;
- immutable snapshot identity or path;
- Pitlord operation and policy source;
- cancellation or deadline state;
- Arcana stderr when available;
- protocol operation and Arcana error code/message;
- request and response counts for completeness failures.

Normal output remains deterministic policy evidence rather than child-process logs.

## Related docs

- [Architecture](ARCHITECTURE.md)
- [Policy reference](POLICY.md)
- [Current limitations](limits/current-limitations.md)
- [Behavioral-contract matrix](development/behavioral-contract-matrix.md)
- [ADR-0001](decisions/0001-query-arcana-through-versioned-protocol.md)

## Notes

Future daemon or pooled-process integration requires a new lifecycle owner and ADR. It must preserve immutable snapshot identity, exact protocol validation, cancellation, and fail-closed policy semantics.
