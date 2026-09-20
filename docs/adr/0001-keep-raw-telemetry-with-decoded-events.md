---
status: accepted
date: 2026-09-19
decision-makers: [repository owner]
consulted: []
informed: []
---

# ADR-0001: Keep raw telemetry with decoded events

## Context and Problem Statement

Parsing and normalization can contain bugs, omit source-specific fields, or
change meaning as schemas evolve. Watchtower must support debugging, forensic
analysis, and future reprocessing without treating today's parser as the only
truth about an event.

## Decision Drivers

- Original evidence must remain available for replay and reprocessing.
- Domain processing needs a typed, validated representation.
- A transport implementation must not define the domain model.
- Callers must not be able to mutate retained raw evidence accidentally.

## Considered Options

- Keep only the parsed event.
- Keep only raw bytes and parse them repeatedly.
- Keep a typed event and an immutable copy of the exact raw payload together.

## Decision Outcome

Chosen option: **keep a typed event and an immutable raw payload together**,
because it preserves original evidence while giving downstream code a stable,
validated contract.

### Consequences

- Good, because parser changes can be tested against and applied to old data.
- Good, because unknown source fields survive even when the current model does
  not understand them.
- Good, because transports exchange serialized bytes rather than domain types.
- Bad, because retaining and defensively copying payloads consumes additional
  memory and storage.
- Bad, because access to raw telemetry requires stronger privacy and logging
  discipline.

### Confirmation

Decoder tests verify byte-for-byte retention, independence from caller buffer
mutation, forward-compatible unknown fields, and errors that do not expose raw
telemetry values. Reviews apply `INV-001` through `INV-003`.

## Pros and Cons of the Options

### Keep only the parsed event

- Good, because it minimizes memory and makes downstream use simple.
- Bad, because omitted or misparsed evidence cannot be recovered.

### Keep only raw bytes

- Good, because it preserves the source representation exactly.
- Bad, because every consumer must parse and validate independently.
- Bad, because parser behavior and errors become inconsistent across stages.

### Keep both representations

- Good, because it supports typed processing and faithful preservation.
- Bad, because it creates copying, storage, and access-control costs.

## More Information

- [`INVARIANTS.md`](../../INVARIANTS.md)
- [`ARCHITECTURE.md`](../../ARCHITECTURE.md)
