# Watchtower

Watchtower is a personal learning project for security observability,
Azure infrastructure, Go, Terraform, and distributed systems.

## Goals

- Learn Azure through direct use of Terraform and Azure-native services.
- Build a realistic but inexpensive security telemetry platform.
- Practice production-quality Go.
- Favor explicit architecture and operational reasoning over clever code.
- Keep steady-state Azure cost below $5/month where practical.
- Keep all infrastructure easily destroyable.

## Engineering principles

- Prefer simple, boring designs over unnecessary abstraction.
- Do not introduce Kubernetes unless explicitly requested.
- Use Terraform for Azure resources.
- Use managed identity / Entra authentication instead of long-lived secrets when possible.
- Never commit credentials, Terraform state, .env files, certificates, or private keys.
- Add tests for meaningful behavior.
- Run formatting, tests, and static checks before considering a task complete.
- Explain significant architectural decisions in docs or code comments where appropriate.

## Agent behavior

- Do not make major architectural choices silently.
- If a task leaves a major design decision unresolved, stop and explain the options before implementing.
- Small implementation choices may be made autonomously.
- Do not provision cloud resources unless the task explicitly asks for deployment.
- Do not use company code, prompts, configs, or proprietary material.
- This repository is a clean-room personal project.

## Commands

- Format: `go fmt ./...`
- Test: `go test ./...`
- Vet: `go vet ./...`
