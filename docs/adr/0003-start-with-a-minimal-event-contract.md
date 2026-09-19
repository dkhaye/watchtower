---
status: accepted
date: 2026-09-19
decision-makers: [repository owner]
consulted: []
informed: []
---

# ADR-0003: Start with a minimal event contract

## Context and Problem Statement

Security telemetry sources represent identities, actions, and resources in
different ways. Watchtower needs enough shared structure to validate and
process initial synthetic events without pretending to have already designed
a universal security-event schema.

## Decision Drivers

- The first contract must be source-neutral and easy to explain.
- Stable event identity and event time are needed for future deduplication and
  timeline reasoning.
- Schema evolution must not discard source fields the model does not yet know.
- Real telemetry requirements, rather than speculation, should justify richer
  actor and target structures.

## Considered Options

- Keep events entirely source-specific with no normalized fields.
- Design a comprehensive normalized security schema before ingesting data.
- Begin with a minimal required envelope and preserve unknown fields in the raw
  representation.

## Decision Outcome

Chosen option: **a minimal required envelope**, containing `id`, `timestamp`,
`actor`, `action`, `target`, and `source`. Timestamp is RFC 3339 event time;
the other fields are initially non-blank strings. Unknown JSON fields are
accepted and retained in the raw payload.

### Consequences

- Good, because initial components can share a typed event without depending
  on a particular telemetry provider.
- Good, because additive source fields do not break decoding or disappear from
  preserved evidence.
- Good, because event time is explicit rather than conflated with future ingest
  and processing clocks.
- Bad, because string actors and targets cannot yet represent complex identity
  or resource relationships.
- Bad, because future schema changes will require compatibility decisions and
  likely explicit versioning.

### Confirmation

Decoder tests require all six fields, enforce RFC 3339 timestamps, accept an
unknown `schema_version` field, and verify exact raw-payload retention.

## Pros and Cons of the Options

### Source-specific events only

- Good, because no source detail is forced into a premature shared model.
- Bad, because every downstream component must understand every source.

### Comprehensive normalized schema

- Good, because it could provide rich common semantics immediately.
- Bad, because the design would be speculative before real source integration.
- Bad, because incorrect normalization choices would become expensive
  contracts.

### Minimal required envelope

- Good, because it establishes only the common fields needed by the next local
  increments.
- Good, because raw preservation provides an escape hatch for unknown fields.
- Bad, because it knowingly leaves important security semantics unmodeled.

## More Information

- [ADR-0001: Keep raw telemetry with decoded events](0001-keep-raw-telemetry-with-decoded-events.md)
- [`ARCHITECTURE.md`](../../ARCHITECTURE.md)
