#!/usr/bin/env bash

# ensure the execution will stop if any command fails (returns non-zero value)
set -e -o pipefail

echo "execute go test"
go test -race -short -timeout 60s ./...