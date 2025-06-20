# Scripts

This directory contains utility scripts for development and CI.

## CI Scripts

- `ci-lint.sh` - Simple linting script that uses Go's built-in tools instead of golangci-lint. This avoids compatibility issues between golangci-lint and different Go versions. The script runs:
  - `go vet` - Examines Go source code and reports suspicious constructs
  - `gofmt` - Formats Go code according to standard guidelines
  - `go mod tidy` - Ensures module dependencies are up to date

This script is used by both the Makefile (`make lint`) and GitHub Actions CI workflows.