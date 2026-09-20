---
status: accepted
date: 2026-09-19
decision-makers: [repository owner]
consulted: []
informed: []
---

# ADR-0004: Bound serialized events to 1 MiB

## Context and Problem Statement

Telemetry is attacker-controlled input. Parsing, retaining, and copying an
unbounded serialized event can consume excessive CPU and memory even when its
known fields are small. The decoder needs a defensible ceiling before parsing,
and future transports need a shared upper bound while reading.

## Decision Drivers

- Untrusted input must not trigger unbounded parsing or copying work.
- The limit must be simple to apply consistently at transport boundaries.
- Initial synthetic and Azure activity events do not need multi-megabyte
  payloads.
- Exceeding the limit must be distinguishable from malformed JSON and an
  invalid event contract.

## Considered Options

- Accept payloads of any size.
- Set a 1 MiB serialized-event limit.
- Make the limit configurable immediately.

## Decision Outcome

Chosen option: **set a 1 MiB serialized-event limit**. `event.Decode` rejects
larger byte slices before JSON validation or copying. Future transports must
enforce the same or a stricter limit while reading so they do not allocate an
oversized input before the decoder can reject it.

### Consequences

- Good, because parser and defensive-copy work have a clear upper bound.
- Good, because callers can handle oversized input separately from malformed
  or structurally invalid events.
- Good, because one exported constant gives future local transports a shared
  default.
- Bad, because legitimate source events above 1 MiB cannot be ingested without
  a future decision.
- Bad, because decoder enforcement alone does not bound a transport's initial
  read allocation.

### Confirmation

Tests verify that a payload one byte over the limit is classified as oversized
before its malformed contents are considered. Reviews require future
transports to apply a bound during reads.

## Pros and Cons of the Options

### Accept payloads of any size

- Good, because no legitimate event is rejected solely for size.
- Bad, because CPU, memory, and copying work remain attacker-controlled.

### Set a 1 MiB serialized-event limit

- Good, because the rule is explicit, deterministic, and ample for the
  project's initial event sources.
- Bad, because the chosen threshold may need revision after real source data
  is observed.

### Make the limit configurable immediately

- Good, because deployments could accommodate different source envelopes.
- Bad, because no configuration system or distinct deployment requirements
  exist yet.
- Bad, because configuration could weaken the invariant without review.

## More Information

- [`INVARIANTS.md`](../../INVARIANTS.md)
- [ADR-0003: Start with a minimal event contract](0003-start-with-a-minimal-event-contract.md)
