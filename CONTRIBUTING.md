# Contributing

## Local setup

Install the versions documented in [README.md](README.md):

- Go
- `just`

No separately installed Go lint or security tools are required. They are
pinned in `justfile` and downloaded automatically by their recipes. Each tool
runs with an isolated module dependency graph rather than becoming part of the
application module.

Before opening or updating a pull request, run:

```sh
just check
```

Use `just --list` to discover individual recipes. Prefer an individual recipe
while iterating, then run the complete contract before pushing.

The event decoder is an untrusted-input boundary. Its fuzz seed corpus runs as
part of ordinary tests; run an active bounded fuzzing session with:

```sh
just fuzz-event
just fuzz-event 30s
```

Active fuzzing is intentionally opt-in rather than part of every CI run.

## Architecture and decisions

Read [ARCHITECTURE.md](ARCHITECTURE.md) and
[INVARIANTS.md](INVARIANTS.md) before changing application behavior or system
boundaries. `ARCHITECTURE.md` describes what exists; it must not present a
proposal as an implemented component.

Architecturally significant changes require a Markdown Architectural Decision
Record in [`docs/adr/`](docs/adr/README.md). Copy the repository template and
include the proposed or accepted decision in the same pull request as its
implementation. Routine implementation choices do not need ADRs. Supersede an
accepted ADR with a new record instead of rewriting its historical rationale.

## Coverage

Run statement coverage locally with:

```sh
just coverage
```

Run `just coverage-html` to also write `.build/coverage.html`. The bootstrap
implementation had 72.7% statement coverage when coverage tracking was added;
that number is historical context, not a target.

Coverage is an observability signal. Review uncovered changed code and trends,
but do not add tests merely to increase a percentage. Tests should exercise
meaningful behavior, boundary conditions, and failure semantics. Generated
code may be excluded when it is introduced; ordinary production code should
not be excluded to improve the report.

## Continuous integration

`.github/workflows/ci.yml` orchestrates pull-request and `main` branch CI. Each
check lives in one reusable workflow named `<domain>.<atom>.yml`. The caller
job is named for its domain and the called job is named for its atom, producing
grouped names such as `go / test` in the GitHub Actions run.

When the same atom must run in several component directories, define the
matrix in `ci.yml` and pass the directory to the reusable workflow. Do not add
one-element matrices.

`.github/workflows/lint-workflows.yml` is deliberately independent of
`ci.yml`. It runs actionlint even when an edit prevents the main orchestrator
from loading correctly.

Independent checks must remain independent so one failure does not hide other
diagnostics. The final `ci / gate` job waits for every required CI atom, always
runs, and succeeds only when all of them succeeded.

The `go / coverage` atom publishes the native Go coverage profile and HTML
report as workflow artifacts and reports coverage to Codecov. Codecov project
and patch statuses are informational. A Codecov outage must not block a merge;
failure to generate or preserve the native report must block it.

The `go / lint` atom emits both human-readable output and SARIF. The SARIF file
is uploaded to GitHub Code Scanning for durable alert history, while the local
linter exit status remains authoritative. Dismissing an alert in GitHub does
not suppress the linter. Accepted findings require a narrow source or
repository configuration suppression with an explanation.

Checkov should follow the same enforcement-plus-SARIF pattern when Terraform
is first added. Do not add an empty infrastructure scanner before then.

## Dependency updates

Dependabot checks Go modules and GitHub Actions each Monday morning in the
repository's local time zone. Minor and patch version updates are grouped per
ecosystem; major updates remain separate for focused review. Security updates
remain separate from the scheduled version-update groups.

GitHub Actions stay pinned to full commit SHAs. Keep the release annotation on
the same line so Dependabot updates both the SHA and its human-readable version
comment.

Dependabot only understands supported dependency manifests. The pinned
`module@version` development tools in `justfile` therefore remain manual
updates until they move to a supported Go tool manifest or another updater is
introduced.

## Repository safety

Do not commit:

- `.local/` content
- credentials, tokens, private keys, certificates, or `.env` files
- Terraform state or non-example variable files
- employer-owned code, prompts, configuration, or documents

Cloud resources must not be provisioned unless the task explicitly requires
it and the target identity and subscription are known to be personal.
