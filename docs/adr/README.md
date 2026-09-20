# Architectural Decision Records

This directory records architecturally significant Watchtower decisions using
[MADR 4.0](https://adr.github.io/madr/).

## Decision log

- [ADR-0000: Use MADR for architectural decisions](0000-use-madr.md)
- [ADR-0001: Keep raw telemetry with decoded events](0001-keep-raw-telemetry-with-decoded-events.md)
- [ADR-0002: Use at-least-once delivery with idempotent processing](0002-use-at-least-once-delivery.md)
- [ADR-0003: Start with a minimal event contract](0003-start-with-a-minimal-event-contract.md) (superseded by ADR-0007)
- [ADR-0004: Bound serialized events to 1 MiB](0004-bound-serialized-events.md)
- [ADR-0005: Keep the public repository self-contained](0005-keep-repository-self-contained.md)
- [ADR-0006: Reject ambiguous JSON event representations](0006-reject-ambiguous-json-event-representations.md)
- [ADR-0007: Use a Go-representable RFC 3339 timestamp profile](0007-use-a-go-representable-rfc3339-timestamp-profile.md)

## When an ADR is required

Add an ADR in the same pull request that introduces or changes a significant:

- system or component boundary;
- data contract or persistence model;
- delivery, retry, checkpoint, or failure semantic;
- security or trust boundary;
- deployment or infrastructure model; or
- dependency that materially constrains the architecture.

Routine implementation details do not require ADRs.

## Lifecycle

1. Copy [`template.md`](template.md) to the next zero-padded sequence number.
2. Use a short kebab-case filename such as `0003-bound-event-size.md`.
3. Open the ADR as `proposed` when a decision still needs agreement.
4. Change it to `accepted` when the implementing pull request is approved.
5. Use `rejected` when an option is deliberately not adopted.
6. Do not rewrite an accepted decision to make history look cleaner. Create a
   new ADR and mark the prior record `superseded`, linking both records.
7. Use `deprecated` when a decision no longer applies and has no replacement.

An accepted ADR may receive non-semantic corrections such as fixed links or
typos. Changes to its rationale or outcome require a new ADR.
