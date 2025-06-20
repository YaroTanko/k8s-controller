#!/bin/bash
set -e

# Run golangci-lint with flags instead of config file
golangci-lint run --timeout=5m