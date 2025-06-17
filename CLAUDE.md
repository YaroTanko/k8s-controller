# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Test Commands
- Build: `go build`
- Run: `go run main.go`
- Test all: `go test ./...`
- Test single file: `go test ./path/to/package/file_test.go`
- Test single function: `go test -run TestFunctionName ./path/to/package`
- Lint: `golangci-lint run`
- Format code: `gofmt -w .`

## Code Style Guidelines
- **Formatting**: Follow Go standard formatting with `gofmt`
- **Imports**: Group imports by standard library, third-party, and local packages
- **Naming**: Use CamelCase for exported names, camelCase for non-exported
- **Error Handling**: Check all errors, don't use panic or log.Fatal in library code
- **Comments**: Document all exported functions, types, and constants
- **Types**: Prefer strong typing, use interfaces for dependency injection
- **Dependencies**: Manage with go modules, minimize external dependencies
- **Testing**: Write unit tests for all functionality, use table-driven tests

This project uses Cobra for CLI commands structure.