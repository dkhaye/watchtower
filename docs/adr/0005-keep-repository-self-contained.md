---
status: accepted
date: 2026-09-19
decision-makers: [repository owner]
consulted: []
informed: []
---

# ADR-0005: Keep the public repository self-contained

## Context and Problem Statement

Contributor instructions are part of the public project interface. Naming or
requiring uncommitted machine-specific files creates an undocumented dependency
that outside contributors cannot satisfy and can disclose information that is
irrelevant to the project.

## Decision Drivers

- A fresh clone must contain every project-specific file named by contributor
  instructions.
- Required setup must be public and reproducible.
- Public documentation must discuss the project rather than contributor-specific
  workstation state.
- Ignored files may support individual workflows without becoming repository
  dependencies.

## Considered Options

- Allow contributor instructions to reference optional machine-specific files.
- Keep references generic without naming individual files.
- Prohibit public references and dependencies on uncommitted machine-specific
  files and scripts.

## Decision Outcome

Chosen option: **prohibit public references and dependencies on uncommitted
machine-specific files and scripts**. Every file or script named by committed
code, documentation, automation, or contributor instructions must be committed
or produced through documented, reproducible setup.

### Consequences

- Good, because a fresh clone is sufficient to understand and use the project.
- Good, because outside contributors receive the same instructions as the
  repository owner.
- Good, because public documentation cannot accidentally describe ignored
  workstation state.
- Bad, because useful local conventions cannot be mentioned in repository
  instructions unless they are made reproducible and generally available.

### Confirmation

Reviews search committed content for references to unavailable files and
scripts, and verify that required generated artifacts have reproducible setup
instructions. `INV-012` and the repository code-review rules make violations
explicit review findings.

## Pros and Cons of the Options

### Allow machine-specific references

- Good, because individual contributors can attach extra instructions easily.
- Bad, because a fresh clone is incomplete and public documentation can expose
  irrelevant contributor details.

### Keep references generic

- Good, because it avoids naming a particular unavailable file.
- Bad, because it still makes hidden workstation state part of the expected
  contributor workflow.

### Prohibit references and dependencies

- Good, because the repository remains portable, public, and auditable.
- Bad, because any genuinely required helper must be committed or given a
  reproducible installation path.

## More Information

- [`INVARIANTS.md`](../../INVARIANTS.md)
- [`CONTRIBUTING.md`](../../CONTRIBUTING.md)
