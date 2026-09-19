set dotenv-load := false
set shell := ["bash", "-euo", "pipefail", "-c"]

binary := ".build/watchtower-linux-amd64"
coverage_profile := ".build/coverage.out"
coverage_html := ".build/coverage.html"
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
    mkdir -p .build
    go run {{golangci_lint}} run ./...

# Run unit tests.
test:
    go test ./...

# Run unit tests with statement coverage and print the function summary.
coverage:
    mkdir -p .build
    go test -covermode=atomic -coverprofile={{coverage_profile}} ./...
    go tool cover -func={{coverage_profile}}

# Generate a browsable HTML coverage report after running coverage.
coverage-html: coverage
    go tool cover -html={{coverage_profile}} -o {{coverage_html}}

# Run unit tests with the race detector.
race:
    go test -race ./...

# Fuzz the untrusted event-decoding boundary for a bounded duration.
fuzz-event duration="10s":
    go test -run='^$' -fuzz=FuzzDecode -fuzztime={{duration}} ./internal/event

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
check: fmt-check tidy-check verify vet lint test coverage race build vuln lint-workflows
