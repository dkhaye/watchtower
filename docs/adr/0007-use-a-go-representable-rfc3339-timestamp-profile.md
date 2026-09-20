---
status: accepted
date: 2026-09-20
decision-makers: [repository owner]
consulted: []
informed: []
---

# ADR-0007: Use a Go-representable RFC 3339 timestamp profile

## Context and Problem Statement

ADR-0003 required RFC 3339 event timestamps but did not say whether the event
contract accepts every representation permitted by that specification. Full
RFC 3339 includes inserted leap seconds, fractional seconds of unbounded
length, and `-00:00` with semantics distinct from a known zero offset. Go's
`time.Time` cannot faithfully represent leap-second notation or precision
beyond nanoseconds, and it models an instant rather than the provenance of its
source offset.

Relying on `time.Parse` alone is also insufficient: its accepted language does
not exactly match RFC 3339. The repository needs a precise timestamp contract
that can be tested without claiming semantics the normalized type cannot
preserve.

This decision supersedes ADR-0003 while retaining its minimal six-field event
envelope and raw-payload preservation requirements.

## Decision Drivers

- The accepted timestamp language must be explicit and exhaustively testable.
- A successfully decoded timestamp must be faithfully representable as
  `time.Time` without silent precision loss or leap-second normalization.
- Exact source spelling and offset metadata must remain available in the raw
  event.
- The initial event model should remain small until a real telemetry source
  demonstrates a need for richer time semantics.

## Considered Options

- Accept all RFC 3339 timestamps using a custom timestamp type.
- Adopt an explicit, Go-representable profile of RFC 3339.
- Continue combining partial lexical checks with `time.Parse` behavior.

## Decision Outcome

Chosen option: **an explicit, Go-representable RFC 3339 profile**.

The decoder accepts timestamps with:

- a four-digit year from `0000` through `9999` and a valid Gregorian date;
- two-digit hour, minute, and second fields, with seconds restricted to
  `00` through `59`;
- `T` or `t` between the date and time;
- no fractional seconds, or a period followed by one through nine digits; and
- `Z` or `z`, or a signed numeric offset from `00:00` through `23:59`.

The `-00:00` convention is accepted. The normalized `time.Time` represents
the event instant only; the exact offset spelling and its "unknown local
offset" meaning remain in `Record.Raw()`.

The decoder rejects leap-second spellings and fractional precision beyond
nine digits even though unrestricted RFC 3339 permits them. A future need for
those representations requires a richer timestamp type and a superseding ADR.
Future Watchtower producers should emit uppercase `T` and `Z` and no more
precision than their source actually provides.

### Consequences

- Good, because accepted timestamps have a complete contract independent of
  undocumented standard-library parser behavior.
- Good, because decoding cannot silently truncate sub-nanosecond precision.
- Good, because exact source notation remains available for audit and replay.
- Bad, because RFC 3339 producers that emit leap seconds or more than nine
  fractional digits must be rejected or adapted at a deliberate boundary.
- Bad, because offset provenance is not present in the normalized timestamp
  and callers must inspect the raw event if it matters.

### Confirmation

A table-driven decoder conformance suite covers every profile production,
minimum and maximum field values, Gregorian calendar boundaries, case
variants, numeric offsets, leap-second rejection, and fractional-precision
rejection. Fuzzing continues to assert exact raw-payload retention for every
accepted event.

## Pros and Cons of the Options

### Accept all RFC 3339 timestamps with a custom type

- Good, because every standards-conforming producer can be represented.
- Good, because leap-second and arbitrary-precision notation can remain typed.
- Bad, because comparison, ordering, serialization, and conversion semantics
  become substantially more complex before a source requires them.

### Adopt an explicit, Go-representable profile

- Good, because the contract maps directly to the current domain type.
- Good, because the restrictions are small, documented, and mechanically
  testable.
- Bad, because it is a documented subset rather than unrestricted RFC 3339.

### Continue relying on partial checks and `time.Parse`

- Good, because it requires little code.
- Bad, because the accepted language is an accidental combination of two
  implementations.
- Bad, because contract gaps are discovered one example at a time.

## More Information

- [RFC 3339: Date and Time on the Internet](https://www.rfc-editor.org/rfc/rfc3339)
- [Go issue 54580: strict RFC 3339 parsing](https://go.dev/issue/54580)
- [ADR-0001: Keep raw telemetry with decoded events](0001-keep-raw-telemetry-with-decoded-events.md)
- [ADR-0003: Start with a minimal event contract](0003-start-with-a-minimal-event-contract.md)
