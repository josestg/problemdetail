#!/usr/bin/env bash

# ensure the execution will stop if any command fails (returns non-zero value)
set -e -o pipefail

echo "execute go vet"
go vet ./...

if ! command -v staticcheck > /dev/null; then
  echo "staticcheck not installed or available in the PATH" >&2
  exit 1
else
    echo "execute staticcheck"
    staticcheck ./...
fi

if ! command -v govulncheck > /dev/null; then
  echo "govulncheck not installed or available in the PATH" >&2
  exit 1
else
    echo "execute govulncheck"
    govulncheck ./...
fi

if ! command -v gosec > /dev/null; then
  echo "gosec not installed or available in the PATH" >&2
  exit 1
else
    echo "execute gosec, will take a while. please be patient!"
    gosec -exclude-generated -quiet ./...
fi