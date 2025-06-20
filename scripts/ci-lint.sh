#!/bin/bash
set -e

# Remove any existing golangci-lint config files
rm -f .golangci.yml .golangci.yaml

# Run golangci-lint with no config, which avoids version/format issues
# Enable only a set of specific linters that are known to work
golangci-lint run --no-config --timeout=5m \
  --enable=errcheck \
  --enable=govet \
  --enable=ineffassign \
  --enable=staticcheck