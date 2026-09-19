set dotenv-load := false
set shell := ["bash", "-euo", "pipefail", "-c"]

binary := ".build/watchtower-linux-amd64"
golangci_lint := "github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.13.2"
govulncheck := "golang.org/x/vuln/cmd/govulncheck@v1.8.0"
actionlint := "github.com/rhysd/actionlint/cmd/actionlint@v1.7.12"

# List available recipes.
default:
    @just --list

# Format Go source files in place.
fmt:
    gofmt -w .

# Fail when any Go source file is not formatted.
fmt-check:
    #!/usr/bin/env bash
    files="$(gofmt -l .)"
    if [[ -n "${files}" ]]; then
      printf 'The following files are not formatted:\n%s\n' "${files}"
      exit 1
    fi

# Update Go module metadata in place.
tidy:
    go mod tidy

# Fail when go.mod or go.sum is not tidy without modifying either file.
tidy-check:
    go mod tidy -diff

# Verify downloaded module content against go.sum.
verify:
    go mod verify

# Run the Go vet analyzer.
vet:
    go vet ./...

# Run the configured Go linters.
lint:
    go run {{golangci_lint}} run ./...

# Run unit tests.
test:
    go test ./...

# Run unit tests with the race detector.
race:
    go test -race ./...

# Build a reproducible Linux/amd64 executable.
build version="dev":
    mkdir -p "$(dirname '{{binary}}')"
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w -X main.version={{version}}" -o "{{binary}}" ./cmd/watchtower

# Scan reachable Go code for known vulnerabilities.
vuln:
    go run {{govulncheck}} ./...

# Lint every GitHub Actions workflow.
lint-workflows:
    go run {{actionlint}}

# Run the complete local verification contract.
check: fmt-check tidy-check verify vet lint test race build vuln lint-workflows
