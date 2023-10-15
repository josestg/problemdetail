#!/usr/bin/make -f

# Choosing the shell
# - [docs](https://www.gnu.org/software/make/manual/html_node/Choosing-the-Shell.html)
SHELL =/bin/bash


# The default target is to prepare the development environment.
all: tools lint

.PHONY: lint
lint: hack/go-lint.sh
	@chmod +x hack/go-lint.sh
	@echo "Running linter."
	@hack/go-lint.sh
	@echo "Linter done."

.PHONY: test
test: hack/go-unittest.sh
	@chmod +x hack/go-unittest.sh
	@echo "Running unit tests."
	@hack/go-unittest.sh
	@echo "Unit tests done."


# Install all development tools, these tools are used by pre-commit hook.
tools: hack/install-tools.sh
	@echo "Installing tools"
	@hack/install-tools.sh
	@echo "Tools installed"


# Enable pre-commit hook.
setup-pre-commit:
	@echo "Setting up pre-commit hook"
	@cp -f hack/pre-commit.sh .git/hooks/pre-commit
	@chmod +x .git/hooks/pre-commit
