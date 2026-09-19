# Watchtower Project Context

## Purpose

Watchtower is a personal engineering and learning project focused on:

- security observability
- distributed systems
- cloud infrastructure
- Azure
- Terraform
- Go
- event streaming
- reliability engineering
- detection engineering
- production-quality CI/CD

The project is intended to resemble a small but serious security-telemetry
platform rather than a tutorial application.

It should provide hands-on experience designing, implementing, operating,
breaking, and improving the kinds of systems used to collect and analyze
security telemetry.

## Engineering goals

Watchtower should eventually:

1. ingest security-relevant events from multiple sources;
2. preserve raw telemetry for debugging, forensics, and replay;
3. normalize events into a stable internal representation;
4. enrich events where appropriate;
5. support detection/correlation;
6. support forensic queries;
7. expose its own operational health;
8. tolerate partial failure and backpressure;
9. support safe replay without overwhelming live processing;
10. provide meaningful reliability and security guarantees.

The system should be built incrementally. Early implementations should be
simple, but architectural boundaries should be chosen deliberately so that
real cloud infrastructure can replace local implementations without rewriting
the domain model.

## Cost constraint

This is a personal project.

Target:

- steady-state Azure cost: <= $5/month where practical
- temporary/lab-mode cost: <= $20/month
- expensive resources should be ephemeral
- all infrastructure should be easy to identify and destroy

Do not provision cloud infrastructure merely because it is convenient.

## Cloud strategy

Azure is the primary cloud because learning Azure is one of the project goals.

Infrastructure should be managed with Terraform.

Likely Azure services include:

- Azure Storage / Blob Storage
- Azure Monitor / Log Analytics
- Azure Event Hubs
- Azure Container Apps
- Managed Identities / Entra ID
- Azure Key Vault
- potentially Microsoft Sentinel later

Prefer identity-based authentication over long-lived credentials.

Do not introduce Kubernetes simply because it is familiar.

## AWS-to-Azure mental mapping

Useful conceptual translations:

- IAM -> Entra ID / Azure RBAC
- CloudTrail -> Azure Activity Log
- CloudWatch -> Azure Monitor / Log Analytics
- MSK / event streaming -> Azure Event Hubs
- KMS / Secrets Manager -> Key Vault
- EKS -> AKS
- Lambda -> Azure Functions
- VPC -> VNet

These mappings are conceptual rather than exact equivalences.

## Event pipeline direction

The conceptual pipeline is:

    source telemetry
          |
          v
        ingest
          |
          v
       preserve raw
          |
          v
       normalize
          |
          v
        enrich
          |
          v
    detect/correlate
          |
          v
      alert/query

Important concepts to preserve as the design evolves:

- raw event preservation
- event time vs ingest time vs processing time
- at-least-once delivery
- idempotent processing
- deduplication
- schema evolution
- backpressure
- bounded memory
- replay
- poison/malformed events
- dead-letter handling
- consumer lag
- live traffic priority over replay traffic
- safe throttling of replay
- explicit failure semantics

Replay should eventually use a controlled path rather than blindly dumping
historical traffic back into the live pipeline.

## Initial local architecture

Before introducing Azure Event Hubs, the first local implementation may use
newline-delimited JSON as a deliberately simple durable message transport.

The important abstraction boundary is serialized messages / raw bytes, not
already-parsed domain objects.

Transport implementations should be replaceable without changing the event
domain model.

A minimal event model starts from:

- id
- timestamp
- actor
- action
- target
- source
- raw/original representation

This is intentionally incomplete and should evolve based on real telemetry
requirements.

The initial ingester should focus on:

- parsing
- validation
- rejecting malformed input safely
- preserving original payloads
- explicit restart/failure semantics

Avoid prematurely building search, databases, enrichment, or detection into
the local v0.

## Language

Primary application language: Go.

Reasons:

- strong cloud/infrastructure ecosystem
- straightforward concurrency model
- good networking/HTTP/JSON support
- simple static deployment
- excellent built-in tooling
- useful ecosystem exposure

Rust is a preferred language personally and may be used experimentally for
specific components later.

Avoid Python unless an ecosystem dependency provides a compelling reason to
use it.

## Go quality strategy

Expected baseline checks:

- gofmt
- go mod tidy cleanliness
- go mod verify
- go vet / golangci-lint
- go test
- go test -race
- go build
- govulncheck

Likely golangci-lint checks include:

- govet
- staticcheck
- errcheck
- ineffassign
- unused
- errorlint
- gosec
- misspell

Targeted native Go fuzzing should be used at trust boundaries such as event
parsing.

Mutation testing may be explored later, probably using Gremlins, but should
not initially block every PR.

## CI philosophy

CI should maximize diagnostic information per iteration.

Do NOT create a serial pipeline where one failure causes unrelated checks to
be skipped.

Independent checks should start in parallel from the same workflow trigger.

Conceptually:

                +-- format
                +-- tidy
                +-- lint
    PR ---------+-- test
                +-- race
                +-- build
                +-- vulnerability
                        |
                        v
                       gate

Individual checks should not depend on each other unless there is a real data
dependency.

A final gate job should:

- wait for all required checks;
- run even when one or more checks fail;
- inspect all job results;
- succeed only if every required check succeeded.

Branch protection should depend on the aggregate gate rather than the
implementation details of every individual job.

Where practical, CI jobs should invoke repository-local commands so that
developers and agents can execute exactly the same checks locally.

The repository, not GitHub Actions YAML, should define what correctness means.

## Agentic-development philosophy

Codex may perform substantial implementation work.

The human remains responsible for:

- architecture
- tradeoffs
- system boundaries
- reliability semantics
- security assumptions
- understanding and defending the resulting implementation

Major design decisions should not be made silently by an agent.

The desired development loop is:

    problem
      |
    human design
      |
    design review / challenge
      |
    implementation agent
      |
    CI + automated review
      |
    rework
      |
    human understanding / approval

Separately, some coding practice will intentionally be done without AI.

## Orchestration

Native Codex workflows and OpenAI Symphony should be explored before building
new personal orchestration infrastructure.

Desired eventual behavior includes:

- task dependency DAGs
- parallel execution of unblocked work
- isolated workspaces
- PR creation
- full CI observation
- automated review
- rework after CI/review feedback
- rebase/conflict handling
- merge/human gates
- dependent-task activation after merge

Do not copy or derive code, prompts, scripts, configuration, or other material
from employer-owned/internal orchestration systems.

Any personal orchestration code must be clean-room work.

## Security principles

Watchtower itself should be treated as an attack target.

Threats to consider include:

- forged telemetry
- deleted telemetry
- altered telemetry
- compromised ingestion identities
- excessive access to sensitive logs
- replay abuse
- telemetry floods / denial of service
- poisoned normalization
- malicious or malformed event payloads
- compromised CI/CD
- supply-chain compromise

Relevant concepts include:

- least privilege
- AuthN vs AuthZ
- RBAC / ABAC
- workload identity
- OAuth/OIDC
- TLS / mTLS
- PKI
- secrets management
- encryption
- trust boundaries
- immutable/tamper-resistant storage
- provenance
- auditability

## Repository hygiene

This is a clean-room personal project.

Never commit:

- employer code
- employer prompts
- employer configuration
- employer credentials/tokens
- proprietary architectural documents
- secrets
- Terraform state
- private keys
- .env credentials

The project may eventually be presented publicly as a portfolio artifact, so
all committed material should be safe to publish.

## Current state

The repository exists locally and has project-level AGENTS.md instructions.

The first application boundary decodes one serialized JSON payload into a
minimal typed event while retaining an immutable copy of the exact original
bytes. It accepts unknown fields for forward-compatible evolution and
distinguishes malformed JSON from a structurally invalid event. Transport,
persistence, rejection policy, and replay remain unimplemented.

`ARCHITECTURE.md` documents the implemented system, `INVARIANTS.md` codifies
non-negotiable constraints, and `docs/adr/` records significant architectural
decisions using MADR 4.0.

Cloud provisioning should not begin until the Azure identity/subscription
being used is clearly personal and under the user's control.

The first implementation work should remain intentionally small enough to
understand every architectural assumption before Event Hubs is introduced.

## Repository commands and CI structure

The repository uses `just` as its local command runner. Go remains responsible
for package dependency analysis, compilation, testing, and build caching.

GitHub Actions uses a central `ci.yml` orchestrator for pull requests and
pushes to `main`. Each check is an independently callable reusable workflow
named `<domain>.<atom>.yml`. Repeated checks across component directories
should use a caller-side matrix. Workflow linting remains a separate top-level
workflow so errors in the central orchestrator do not hide actionlint results.

Native Go coverage profiles are the source of truth for coverage reporting.
Coverage is tracked as a diagnostic signal through CI artifacts and Codecov,
not enforced as a percentage threshold. The bootstrap baseline was 72.7%
statement coverage when tracking was introduced.

Static-analysis atoms should retain their local pass/fail behavior while also
uploading SARIF to GitHub Code Scanning when the tool supports it. A dismissal
in GitHub does not replace a narrow, documented suppression in repository
configuration. Checkov should be introduced with the first Terraform code and
follow this pattern.

Dependabot monitors Go modules and SHA-pinned GitHub Actions. Routine minor and
patch updates are grouped by ecosystem; major and security updates receive
focused pull requests. Tool versions embedded in `justfile` are outside
Dependabot's manifest support and remain an explicit maintenance item.
