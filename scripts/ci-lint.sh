#!/bin/bash
set -e

# Remove any existing golangci-lint config files
rm -f .golangci.yml .golangci.yaml

# Run golangci-lint with flags instead of config file
golangci-lint run --timeout=5m