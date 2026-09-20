---
status: accepted
date: 2026-09-19
decision-makers: [repository owner]
consulted: []
informed: []
---

# ADR-0002: Use at-least-once delivery with idempotent processing

## Context and Problem Statement

Failures can occur after an event has been processed but before its progress is
durably acknowledged. Avoiding all duplicate delivery would require coupling
processing, side effects, and checkpoints into a stronger protocol than the
initial local system or likely managed transports can honestly guarantee.

## Decision Drivers

- Events must not be silently lost merely to avoid duplicates.
- Failure and restart behavior must be explainable and testable.
- Transport implementations should remain replaceable.
- The design must work with managed event systems that commonly redeliver.

## Considered Options

- At-most-once delivery.
- At-least-once delivery with idempotent processing and deduplication where
  side effects require it.
- A system-wide exactly-once guarantee.

## Decision Outcome

Chosen option: **at-least-once delivery with idempotent processing**, because
it makes potential duplication explicit while prioritizing durable event
handling over an unsupported exactly-once claim.

### Consequences

- Good, because retry and crash recovery do not require silently dropping
  events.
- Good, because the model maps to likely local and Azure transports.
- Bad, because consumers and side effects must tolerate duplicate events.
- Bad, because stable identifiers, deduplication scope, and retention will
  require explicit future decisions.
- Bad, because replay must be distinguishable from accidental redelivery when
  side effects differ.

### Confirmation

Future transports will acknowledge progress only according to documented
durability semantics. Stateful consumers will test repeated delivery of the
same event identifier. Reviews reject undocumented exactly-once claims.

## Pros and Cons of the Options

### At-most-once delivery

- Good, because consumers do not handle duplicates.
- Bad, because crashes and transient failures can lose unprocessed telemetry.

### At-least-once delivery

- Good, because recovery favors eventual processing over silent loss.
- Bad, because idempotency and deduplication become application concerns.

### System-wide exactly-once delivery

- Good, because it would simplify consumer-visible outcomes if achieved.
- Bad, because the claim requires end-to-end transactional assumptions that
  do not currently exist.
- Bad, because a casual claim would hide real failure modes rather than solve
  them.

## More Information

- [`INVARIANTS.md`](../../INVARIANTS.md)
- [`ARCHITECTURE.md`](../../ARCHITECTURE.md)
