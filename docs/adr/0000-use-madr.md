---
status: accepted
date: 2026-09-19
decision-makers: [repository owner]
consulted: []
informed: []
---

# ADR-0000: Use MADR for architectural decisions

## Context and Problem Statement

Watchtower is a learning project whose value depends on understanding and
defending its architectural choices. Those choices need a durable history that
captures context and tradeoffs without turning routine implementation work
into a heavyweight design process.

## Decision Drivers

- Decisions must remain understandable after their original discussion is gone.
- Records must be easy for humans and coding agents to find and review.
- The format should be recognizable rather than repository-specific.
- Recording a decision should require no dedicated service or runtime.

## Considered Options

- Markdown Architectural Decision Records (MADR) 4.0.
- Michael Nygard's minimal ADR format.
- Unstructured design notes.
- No persistent decision records.

## Decision Outcome

Chosen option: **MADR 4.0**, because it is a recognized Markdown convention
that explicitly captures decision drivers, considered options, consequences,
and confirmation while remaining lightweight.

### Consequences

- Good, because architectural rationale lives beside the code it governs.
- Good, because contributors and reviewers share a predictable structure.
- Bad, because significant changes require maintaining another artifact.
- Bad, because a template can encourage unnecessary ceremony if applied to
  routine implementation choices.

### Confirmation

Architecturally significant pull requests include a proposed or accepted ADR,
and `ARCHITECTURE.md` links to this decision log.

## Pros and Cons of the Options

### MADR 4.0

- Good, because it captures alternatives and their tradeoffs explicitly.
- Good, because it is plain Markdown with maintained templates.
- Bad, because its full form is longer than the original minimal format.

### Michael Nygard's minimal ADR format

- Good, because status, context, decision, and consequences are extremely lean.
- Bad, because considered options and confirmation would rely on author habit.

### Unstructured design notes

- Good, because authors can write whatever a decision seems to need.
- Bad, because important context becomes inconsistent and harder to scan.

### No persistent decision records

- Good, because it creates no documentation overhead.
- Bad, because rationale would be reconstructed from code and pull-request
  history.

## More Information

- [MADR](https://adr.github.io/madr/)
- [ADR templates](https://adr.github.io/adr-templates/)
