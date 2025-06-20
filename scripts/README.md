# Scripts

This directory contains utility scripts for development and CI.

## CI Scripts

- `ci-lint.sh` - Runs golangci-lint without using a config file, avoiding potential config format issues. Used by both the Makefile and GitHub Actions.