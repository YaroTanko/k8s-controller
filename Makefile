# Makefile for k8s-controller

# Variables
APP_NAME := k8s-controller
VERSION := $(shell git describe --tags --always --dirty)
BUILD_TIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
COMMIT_HASH := $(shell git rev-parse --short HEAD)
LDFLAGS := -ldflags "-X k8s-controller/cmd.Version=$(VERSION) -X k8s-controller/cmd.BuildTime=$(BUILD_TIME) -X k8s-controller/cmd.CommitHash=$(COMMIT_HASH)"
GO := go
GOFMT := gofmt
GOIMPORTS := goimports
GOLINT := golangci-lint
DOCKER := docker

# Build binary
.PHONY: build
build:
	$(GO) build $(LDFLAGS) -o bin/$(APP_NAME) main.go

# Run application locally
.PHONY: run
run:
	$(GO) run main.go

# Run server command
.PHONY: server
server:
	$(GO) run main.go server

# Run serve command
.PHONY: serve
serve:
	$(GO) run main.go serve

# Build Docker image
.PHONY: docker-build
docker-build:
	$(DOCKER) build -t $(APP_NAME):$(VERSION) -t $(APP_NAME):latest .

# Clean up
.PHONY: clean
clean:
	rm -rf bin/
	rm -f coverage.out

# Format code
.PHONY: fmt
fmt:
	$(GOFMT) -w .

# Lint code
.PHONY: lint
lint:
	./scripts/lint.sh

# Run tests
.PHONY: test
test:
	$(GO) test ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	$(GO) test -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out

# Build for multiple platforms
.PHONY: build-all
build-all:
	GOOS=linux GOARCH=amd64 $(GO) build $(LDFLAGS) -o bin/$(APP_NAME)-linux-amd64 main.go
	GOOS=darwin GOARCH=amd64 $(GO) build $(LDFLAGS) -o bin/$(APP_NAME)-darwin-amd64 main.go
	GOOS=darwin GOARCH=arm64 $(GO) build $(LDFLAGS) -o bin/$(APP_NAME)-darwin-arm64 main.go
	GOOS=windows GOARCH=amd64 $(GO) build $(LDFLAGS) -o bin/$(APP_NAME)-windows-amd64.exe main.go

# Default target
.PHONY: default
default: build