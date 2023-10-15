#!/usr/bin/env bash

# ensure the execution will stop if any command fails (returns non-zero value)
set -e -o pipefail

echo "check if go.mod needs to be updated"
# Check if go.mod needs to be updated.
if go mod tidy -v 2>&1 | grep -q 'updates to go.mod needed'; then
  echo "please run go mod tidy and commit the changes"
  exit 1
fi

## this will retrieve all of the .go files that have been
## changed since the last commit
go_staged_files=$(git diff --cached --diff-filter=ACM --name-only -- '*.go')


if [[ $go_staged_files == "" ]]; then
  # when there are no staged go files, we can skip the rest of the checks.
  echo "no go files in staged changes. skipping gofmt, goimports, go vet and staticcheck"
else
    if ! command -v gofmt &> /dev/null ; then
        echo "gofmt not installed or available in the PATH" >&2
        exit 1
    fi


    if ! command -v goimports &> /dev/null ; then
      echo "goimports not installed or available in the PATH" >&2
      exit 1
    fi

    for file in $go_staged_files; do
        printf "[gofmt, goimports] %s\n" "$file"
        goimports -l -w "$file"
        gofmt -l -w "$file"
        git add "$file"
    done
fi

make lint
make test