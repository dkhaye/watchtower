# Watchtower

Watchtower is a personal learning project for security observability, Azure
infrastructure, Go, Terraform, and distributed systems. It is intentionally
being built in small, explicit steps rather than as a tutorial-scale facsimile
of a production platform.

The repository currently contains a minimal Go command, the first serialized
event-decoding boundary, and the quality gates that future implementation must
satisfy. Transport, persistence, cloud resources, and detection are not
implemented yet.

Start with [ARCHITECTURE.md](ARCHITECTURE.md) for the current system map,
[INVARIANTS.md](INVARIANTS.md) for non-negotiable constraints, and the
[architectural decision log](docs/adr/README.md) for the rationale behind
significant choices.

## Prerequisites

- [Go 1.27.1](https://go.dev/dl/)
- [just 1.58.0](https://github.com/casey/just/releases/tag/1.58.0)

Run the complete local verification contract with:

```sh
just check
```

Run `just --list` to see the independently executable checks. The first run
downloads the pinned Go-based development tools recorded in `justfile`. Each
tool executes with its own module dependency graph so tool dependencies cannot
affect the application or one another.

Measure statement coverage and generate a local HTML report with:

```sh
just coverage
just coverage-html
```

The HTML report is written to `.build/coverage.html`. Coverage is treated as a
diagnostic signal rather than a target: CI records the report and changed-line
coverage, but does not reject a change solely because of a percentage.

Build and run the bootstrap command with:

```sh
just build
./.build/watchtower-linux-amd64 --version
```

The build target is Linux/amd64, so the produced executable will not run
directly on macOS. Use `go run ./cmd/watchtower --version` for a native local
smoke test.

See [CONTRIBUTING.md](CONTRIBUTING.md) for the development and CI conventions.
