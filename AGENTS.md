# Watchtower

Watchtower is a personal learning project for security observability,
Azure infrastructure, Go, Terraform, and distributed systems.

Read these files before substantial work:

- `ARCHITECTURE.md`
- `INVARIANTS.md`
- `docs/PROJECT_CONTEXT.md`
- `.local/INTERVIEW_CONTEXT.md` when it exists locally

`ARCHITECTURE.md` describes the system that exists today.
`INVARIANTS.md` contains non-negotiable system and repository constraints.
`docs/PROJECT_CONTEXT.md` contains durable project context.
`.local/INTERVIEW_CONTEXT.md` contains private interview-preparation context
and must never be committed.

## Goals

- Learn Azure through direct use of Terraform and Azure-native services.
- Build a realistic but inexpensive security telemetry platform.
- Practice production-quality Go.
- Develop security-observability and distributed-systems interview fluency.
- Favor explicit architecture and operational reasoning over clever code.
- Keep steady-state Azure cost below $5/month where practical.
- Keep all infrastructure easy to identify and destroy.

## Engineering principles

- Prefer simple, boring designs over unnecessary abstraction.
- Do not introduce Kubernetes unless explicitly requested.
- Use Terraform for Azure resources.
- Prefer managed identity / Entra authentication over long-lived secrets.
- Preserve clear boundaries between transport, serialization, domain logic,
  persistence, and detection.
- Make failure semantics explicit.
- Prefer at-least-once + idempotent processing over casually claiming
  exactly-once behavior.
- Preserve raw telemetry where needed for debugging, replay, and forensics.
- Add tests for meaningful behavior.
- Explain significant architectural decisions in docs or ADRs where useful.

## Architecture decisions

- Store architectural decision records in `docs/adr/` using its MADR template.
- Add or supersede an ADR in the same pull request that changes a significant
  system boundary, data contract, delivery or failure semantic, trust boundary,
  persistence model, deployment model, or constraining dependency.
- Do not require ADRs for routine implementation details.
- Do not rewrite an accepted ADR's rationale or outcome. Supersede it with a
  new ADR and preserve the historical record.
- Keep `ARCHITECTURE.md` synchronized with implemented architecture rather
  than presenting proposals as current behavior.

## Agent behavior

- Do not make major architectural choices silently.
- If a task leaves an important design decision unresolved, stop and explain
  the options before implementing.
- Small implementation choices may be made autonomously.
- Do not provision cloud resources unless the task explicitly asks for it.
- Do not introduce new infrastructure merely for convenience.
- Do not use employer code, prompts, configs, tokens, documents, or other
  proprietary material.
- This repository is a clean-room personal project.
- Never commit `.local/`.
- Never commit credentials, Terraform state, `.env` files, certificates,
  private keys, or secrets.

## Development workflow

For substantial changes:

1. Understand the task and acceptance criteria.
2. Inspect relevant project context and existing code.
3. Surface unresolved architectural decisions before coding.
4. Implement the smallest coherent change.
5. Run all relevant local checks.
6. Fix all discovered failures, not just the first one.
7. Summarize meaningful design decisions and remaining risks.

Do not treat passing tests as sufficient if the implementation violates the
stated architecture, security model, or operational requirements.

## CI philosophy

Independent checks should run independently and report all failures in one
iteration.

Avoid serial CI chains where one failing check prevents unrelated checks from
running.

The repository should define correctness through local commands; GitHub
Actions should primarily orchestrate those commands in parallel.

A final aggregate CI gate should succeed only if every required check passes.

## Go commands

Baseline commands:

- Format: `gofmt`
- Test: `go test ./...`
- Race: `go test -race ./...`
- Vet: `go vet ./...`
- Build: `go build ./...`
- Module verify: `go mod verify`

Additional commands such as `golangci-lint`, `govulncheck`, fuzzing, and
mutation testing may be added as the project evolves.

Prefer repository-local CI entrypoints once they exist rather than duplicating
CI logic directly in GitHub Actions YAML.

## Code Review Rules

### Raw telemetry integrity

- Flag any processing path that discards, mutates, or replaces the exact
  serialized telemetry before it is available for durable preservation. A
  parsed representation may accompany raw telemetry but may not replace it.
- Flag logs and error messages that expose raw telemetry or sensitive event
  field values. Safe behavior identifies the violated contract without copying
  payload data into diagnostic text.

### Honest and bounded processing

- Flag exactly-once claims that do not define their scope, protocol, and
  failure assumptions. Safe processing assumes at-least-once delivery and
  makes idempotency or deduplication explicit where duplicates have side
  effects.
- Flag unbounded reads, buffers, queues, retries, concurrency, or replay.
  Safe behavior defines a limit and observable failure or backpressure
  semantics.

### Architecture governance

- Flag changes to an invariant, system boundary, data contract, delivery or
  failure semantic, trust boundary, persistence model, or deployment model
  when the change neither follows an accepted ADR nor includes a proposed ADR.
  Routine implementation details do not require ADRs.
