# Watchtower

Watchtower is a personal learning project for security observability, Azure
infrastructure, Go, Terraform, and distributed systems. It is intentionally
being built in small, explicit steps rather than as a tutorial-scale facsimile
of a production platform.

The repository currently contains a minimal Go command and the quality gates
that future implementation must satisfy. Event ingestion, cloud resources,
storage, and detection are not implemented yet.

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

Build and run the bootstrap command with:

```sh
just build
./.build/watchtower-linux-amd64 --version
```

The build target is Linux/amd64, so the produced executable will not run
directly on macOS. Use `go run ./cmd/watchtower --version` for a native local
smoke test.

See [CONTRIBUTING.md](CONTRIBUTING.md) for the development and CI conventions.
