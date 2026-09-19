# Invariants

These rules are non-negotiable constraints on Watchtower. They apply whenever
the relevant subsystem exists. An architectural decision record (ADR) may
supersede an invariant, but an implementation or review comment may not ignore
one silently.

## INV-001: Preserve original telemetry

The exact serialized telemetry received at the ingestion boundary must remain
available for durable preservation before any lossy transformation. A parsed
or normalized representation does not replace the original evidence.

## INV-002: Treat telemetry as untrusted and sensitive

Telemetry must be validated at trust boundaries. Raw payloads and sensitive
field values must not appear in logs or error messages by default. Preservation
does not imply unrestricted access or disclosure.

## INV-003: Keep architectural boundaries explicit

Transport, serialization, validation, domain processing, persistence, and
detection must not be coupled merely for convenience. Transports deliver
serialized bytes rather than already-parsed domain objects.

## INV-004: Make delivery semantics honest

The system assumes at-least-once delivery and requires idempotent processing
where duplicates have side effects. Exactly-once behavior must not be claimed
without a documented protocol, scope, and failure model.

## INV-005: Bound work and memory

Reads, buffers, queues, retries, concurrency, and replay rates must have
explicit bounds. When a bound is reached, the failure or backpressure behavior
must be observable and documented.

## INV-006: Preserve distinct clocks

Event time, ingest time, and processing time have different meanings and must
not be silently collapsed into one timestamp when those clocks are introduced.

## INV-007: Keep replay controlled

Replay must be identifiable, rate limited, observable, and prevented from
overwhelming live processing. Historical processing must not accidentally
repeat live-only side effects such as notifications.

## INV-008: Use personal, disposable infrastructure

Cloud resources must be declared with Terraform, clearly attributable to this
project, inexpensive at rest, and straightforward to destroy. Provisioning
requires an explicitly identified personal Azure identity and subscription.

## INV-009: Prefer workload identity over credentials

Azure workloads must prefer managed identity and Entra authentication. Secrets,
private keys, credentials, Terraform state, and non-example environment files
must never be committed.

## INV-010: Preserve clean-room provenance

Employer-owned code, prompts, configurations, tokens, documents, and
proprietary implementation details must not enter this repository.

## INV-011: Keep correctness local and CI diagnostic

Repository-local commands define correctness. Independent CI checks run
independently so one failure does not hide unrelated failures, and the final
aggregate gate succeeds only when every required check succeeds.
