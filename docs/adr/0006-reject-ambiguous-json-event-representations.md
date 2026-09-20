---
status: accepted
date: 2026-09-19
decision-makers: [repository owner]
consulted: []
informed: []
---

# ADR-0006: Reject ambiguous JSON event representations

## Context and Problem Statement

Watchtower preserves serialized telemetry and also derives typed fields used
for identity, correlation, and future deduplication. Some JSON representations
do not have a consistent interpretation across parsers. Invalid UTF-8 and lone
UTF-16 surrogate escapes can be replaced with the same Unicode replacement
character, while duplicate object member names can be interpreted using the
first value, the last value, or an error.

Accepting those representations would allow the preserved evidence and the
typed event to acquire different meanings during replay or inspection.

## Decision Drivers

- A serialized event must have one unambiguous interpretation.
- Typed identities must not silently collide after decoding.
- Preserved raw telemetry must remain safe to replay through another
  standards-conforming implementation.
- Validation work must remain bounded by the serialized event size limit.

## Considered Options

- Accept the behavior provided by Go's standard JSON decoder.
- Reject ambiguity only in the required top-level event fields.
- Reject ambiguous encoding and duplicate object members throughout the event.

## Decision Outcome

Chosen option: **reject ambiguous encoding and duplicate object members
throughout the event**.

The decoder requires valid UTF-8, rejects unpaired Unicode surrogate escapes
in all string values and member names, and requires object member names to be
unique at every nesting level after JSON escape decoding. Validation errors do
not include member names or values.

### Consequences

- Good, because an accepted event has a stable interpretation across parsing,
  preservation, and replay boundaries.
- Good, because distinct malformed identities cannot silently normalize to the
  same typed string.
- Good, because escaped and literal forms of the same member name are treated
  as duplicates.
- Bad, because some producers whose parsers tolerate duplicate members must
  correct their output before Watchtower accepts it.
- Bad, because validation performs an additional bounded pass over each event.

### Confirmation

Decoder tests reject invalid UTF-8, unpaired surrogate escapes, and duplicate
member names at the top level and inside unknown nested values. Tests also
accept valid surrogate pairs, literal replacement characters, escaped
backslashes, nested unknown values, and large JSON numbers.

## Pros and Cons of the Options

### Accept the standard decoder behavior

- Good, because it requires no additional validation.
- Bad, because invalid encodings can collapse into the same typed value.
- Bad, because duplicate-member interpretation depends on the parser.

### Reject ambiguity only in required top-level fields

- Good, because it protects the fields used by the current event contract.
- Bad, because preserved unknown fields could become ambiguous when a future
  schema begins interpreting them.
- Bad, because nested evidence could still change meaning across parsers.

### Reject ambiguity throughout the event

- Good, because the whole accepted representation remains safe to preserve and
  replay.
- Good, because future schema evolution cannot expose ambiguity that was
  accepted earlier inside an unknown field.
- Bad, because validation must recursively inspect the bounded payload.

## More Information

- [ADR-0001: Keep raw telemetry with decoded events](0001-keep-raw-telemetry-with-decoded-events.md)
- [ADR-0003: Start with a minimal event contract](0003-start-with-a-minimal-event-contract.md)
- [ADR-0004: Bound serialized events to 1 MiB](0004-bound-serialized-events.md)
- [`ARCHITECTURE.md`](../../ARCHITECTURE.md)
