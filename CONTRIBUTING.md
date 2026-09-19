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

## Repository safety

Do not commit:

- `.local/` content
- credentials, tokens, private keys, certificates, or `.env` files
- Terraform state or non-example variable files
- employer-owned code, prompts, configuration, or documents

Cloud resources must not be provisioned unless the task explicitly requires
it and the target identity and subscription are known to be personal.
