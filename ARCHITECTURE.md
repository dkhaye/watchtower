# Architecture

This document describes the architecture that exists today. It is a living
map for contributors, not a promise that every future component is already
implemented. Proposed changes belong in an architectural decision record
(ADR) before this document presents them as current behavior.

## System purpose

Watchtower is an incrementally built security-telemetry platform. Its eventual
job is to ingest security-relevant events, preserve their original evidence,
normalize and enrich them, run detections, and support investigation and
controlled replay.

The current implementation is deliberately smaller: a Go command, an event
decoding boundary, and the repository quality system around them. There is no
transport, persistence layer, cloud deployment, or detection engine yet.

## Current system

```text
serialized JSON payload
         |
         v
  internal/event.Decode
         |
         +----> typed Event
         |
         +----> exact raw payload copy
```

`event.Decode` is a trust boundary. It accepts one serialized payload, checks
its syntax and minimum schema, and returns a `Record` containing both the typed
event and an immutable copy of the input bytes.

The decoder performs no I/O. A future transport will frame messages and pass
serialized bytes into this boundary. A future persistence component will
store raw telemetry before downstream processing can make lossy changes.

## Repository structure

```text
.
├── cmd/watchtower/       # Process entry point and command-line behavior
├── internal/event/       # Minimal event contract and JSON decoding boundary
├── docs/adr/             # Architectural decisions and their rationale
├── .github/workflows/    # Independently callable CI atoms and orchestrators
├── AGENTS.md             # Contributor and Codex instructions/review rules
├── INVARIANTS.md         # Non-negotiable system and repository constraints
├── CONTRIBUTING.md       # Local development and CI conventions
└── justfile              # Repository-local verification contract
```

New packages should represent real architectural boundaries. Do not create
layers or interfaces solely in anticipation of hypothetical implementations.

## Event contract

The initial normalized event contains:

- `id`: source- or producer-assigned stable event identifier;
- `timestamp`: event time encoded as RFC 3339, including fractional seconds;
- `actor`: identity responsible for the activity;
- `action`: activity that occurred;
- `target`: resource or object affected by the activity; and
- `source`: telemetry system that produced the event.

All fields are currently strings except the decoded timestamp. This is a
minimal learning contract, not a universal security-event schema. Real source
requirements should drive future structure.

Unknown JSON fields are accepted for forward compatibility and remain present
in the raw payload. The decoder does not silently trim or otherwise normalize
field values. Blank required fields are rejected. Serialized events larger
than 1 MiB are rejected before JSON parsing; future transports must apply the
same or a stricter limit while reading so the initial input allocation is also
bounded.

## Failure model

Decoding distinguishes three caller-actionable classes:

- oversized payloads, which are rejected before parsing;
- malformed JSON, where the serialized representation is invalid; and
- invalid events, where the JSON is valid but does not satisfy the event
  contract.

Errors identify the violated contract without copying telemetry values into
the error text. Transport-level retry, rejection, dead-letter, and checkpoint
semantics have not been chosen yet and do not belong in the decoder.

## Architectural boundaries

The following boundaries guide the next increments:

1. **Transport** frames and delivers serialized messages. It does not pass
   already-parsed domain objects across the ingestion boundary.
2. **Serialization and validation** decode bytes into the event contract while
   keeping the original bytes available for preservation.
3. **Persistence** durably stores raw telemetry and, separately where useful,
   normalized representations and processing state.
4. **Domain processing** normalizes, enriches, and detects without depending
   on a specific transport implementation.
5. **Operational policy** owns retry, backpressure, dead-letter, checkpoint,
   and replay behavior explicitly rather than hiding it in helpers.

Only the serialization and validation boundary exists in application code
today.

## Delivery and replay direction

Watchtower is designed around at-least-once delivery and idempotent processing.
Duplicates are expected and must not silently corrupt derived state. No
component may claim exactly-once behavior without documenting the protocol and
failure assumptions that establish it.

Replay will use an explicit, bounded path that cannot overwhelm or become
indistinguishable from live traffic. Those mechanics remain future decisions.

## Security and trust boundaries

Telemetry is untrusted and may be malformed, forged, oversized, malicious, or
sensitive. Preserving raw telemetry is not permission to log or broadly expose
it. Future components must apply explicit size limits, least privilege, and
data-access controls at their boundaries.

Azure infrastructure will be declared with Terraform. Workloads should use
managed identity and Entra authentication instead of long-lived credentials.
No Azure resources currently exist.

## Build, test, and delivery

`just check` is the local correctness contract. GitHub Actions invokes the
same repository commands through independent reusable workflow atoms, then
combines their results in `ci / gate`.

Tests cover behavior and failure semantics. Coverage is recorded as a
diagnostic signal rather than enforced as a percentage target. Static-analysis
findings remain hard CI failures and are also uploaded as SARIF for durable
triage.

## Architectural decisions

Accepted and proposed decisions are indexed in [`docs/adr/README.md`](docs/adr/README.md).
The invariants derived from those decisions are collected in
[`INVARIANTS.md`](INVARIANTS.md).

## Near-term evolution

The next likely components are:

1. a bounded newline-delimited JSON transport reader;
2. explicit per-record rejection and continuation semantics;
3. a synthetic event generator;
4. a local ingestion command; and
5. durable raw-event storage.

Each remains subject to its own acceptance criteria and ADR when it changes an
architecturally significant contract.
