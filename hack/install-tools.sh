#!/usr/bin/env bash

# ensure the execution will stop if any command fails (returns non-zero value)
set -e -o pipefail

# install staticcheck if not installed yet.
if ! command -v staticcheck > /dev/null; then
  echo "Installing staticcheck"
  go install honnef.co/go/tools/cmd/staticcheck@latest
else
  echo "staticcheck already installed"
fi

# install goimports if not installed yet.
if ! command -v goimports > /dev/null; then
  echo "Installing goimports"
  go install golang.org/x/tools/cmd/goimports@latest
else
    echo "goimports already installed"
fi

# install gosec if not installed yet.
if ! command -v gosec > /dev/null; then
  echo "Installing gosec"
  curl -sfL https://raw.githubusercontent.com/securego/gosec/master/install.sh | sh -s -- -b "$(go env GOPATH)"/bin
else
  echo "gosec already installed"
fi

# install govulncheck if not installed yet.
if ! command -v govulncheck > /dev/null; then
  echo "Installing govulncheck"
  go install golang.org/x/vuln/cmd/govulncheck@latest
else
    echo "govulncheck already installed"
fi
