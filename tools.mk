build/tools/gobin/Makefile:
	git submodule update --init --recursive $(dir $@)
include build/tools/gobin/Makefile

# Tools shall be sorted in alphabetical order.

# GOBIN_BINDIR is the top-level directory for where executable binaries should be placed. The actual binaries are placed
# in subdirectories under GOBIN_BINDIR, based on the tool's name and version.
GOBIN_BINDIR := build/tools/bin

# `gofumports` is a tool for formatting Go files and insert/removing/sorting import statements. `gofumports` enforces a
# stricter formatting than `go fmt`.
GOFUMPORTS_VERSION := 96300e3d49fbb3b7bc9c6dc74f8a5cc0ef46f84b
GOFUMPORTS_BINDIR := $(GOBIN_BINDIR)/gofumports/$(GOFUMPORTS_VERSION)
GOFUMPORTS := $(GOFUMPORTS_BINDIR)/gofumports
$(GOFUMPORTS): | $(GOBIN) $(GOFUMPORTS_BINDIR)
	GOBIN=$(GOFUMPORTS_BINDIR) $(GOBIN) mvdan.cc/gofumpt/gofumports@$(GOFUMPORTS_VERSION)

# `golangci-lint` is a tool to run numerous linters for Go code.
GOLANGCI_LINT_VERSION := v1.17.1
GOLANGCI_LINT_BINDIR := $(GOBIN_BINDIR)/golangci-lint/$(GOLANGCI_LINT_VERSION)
GOLANGCI_LINT := $(GOLANGCI_LINT_BINDIR)/golangci-lint
$(GOLANGCI_LINT): | $(GOBIN) $(GOLANGCI_LINT_BINDIR)
	GOBIN=$(GOLANGCI_LINT_BINDIR) $(GOBIN) github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

BINDIRS := \
	$(GOFUMPORTS_BINDIR) \
	$(GOLANGCI_LINT_BINDIR)
$(BINDIRS):
	mkdir --parent '$@'

# tools-all: install all tools. Useful when setting up a new development environment or in CI. The latter advantage
# comes from the fact that all tools are installed to the same directory (see the `GOBIN_BINDIR` variable above) and
# thus can easily be cached in a CI environment.
.PHONY: tools-all
tools-all: | \
	$(GOFUMPORTS) \
	$(GOLANGCI_LINT)
