#!/bin/bash
set -e

# Remove any existing golangci-lint config files
rm -f .golangci.yml .golangci.yaml

echo "Running individual linters to avoid metalinter issues..."

# Initialize error status
ERROR=0

# Run errcheck
echo "Running errcheck..."
if ! golangci-lint run --no-config --timeout=5m --enable=errcheck --disable=govet,ineffassign,staticcheck; then
  ERROR=1
  echo "errcheck failed!"
fi

# Run govet
echo "Running govet..."
if ! golangci-lint run --no-config --timeout=5m --enable=govet --disable=errcheck,ineffassign,staticcheck; then
  ERROR=1
  echo "govet failed!"
fi

# Run ineffassign
echo "Running ineffassign..."
if ! golangci-lint run --no-config --timeout=5m --enable=ineffassign --disable=errcheck,govet,staticcheck; then
  ERROR=1
  echo "ineffassign failed!"
fi

# Run staticcheck
echo "Running staticcheck..."
if ! golangci-lint run --no-config --timeout=5m --enable=staticcheck --disable=errcheck,govet,ineffassign; then
  ERROR=1
  echo "staticcheck failed!"
fi

if [ $ERROR -eq 0 ]; then
  echo "✅ All linters passed!"
else
  echo "❌ One or more linters failed!"
  exit 1
fi