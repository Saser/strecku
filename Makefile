# all: run everything needed for a complete build that passes CI.
.PHONY: all
all: \
	fix \
	lint \
	build \
	test

# fix: run tools that generate build files, fix formatting, etc.
.PHONY: fix
fix: \
	gazelle \
	buildifier \
	gofumports \
	go-mod-fix

# lint: run linters for build files, Go files, etc.
.PHONY: lint
lint: \
	lint-gazelle \
	lint-buildifier \
	lint-gofumports \
	lint-golangci-lint \
	go-lint-mod-fix

# build: build the entire project.
.PHONY: build
build: \
	bazel-build

# test: test the entire project.
.PHONY: test
test: \
	bazel-test

# WD: the absolute path to the current working directory. It is used for referring to the root directory of this project
# instead of using e.g. `.` to refer to "this directory". This is necessary when invoking tools such as `gofumports`
# using Bazel.
WD := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))

GIT_LS_FILES := git ls-files --exclude-standard --cached --others

BUILD_FILES := $(shell $(GIT_LS_FILES) -- '*BUILD.bazel')
GO_FILES := $(shell $(GIT_LS_FILES) -- '*.go')

ifeq ($(CI),true)
	BAZELRC := build/ci/.bazelrc
else
	BAZELRC := .bazelrc
endif
BAZEL_FLAGS := \
	--bazelrc=$(BAZELRC)

.PHONY: \
	build/tools/bazel \
	build/tools/circleci \
	build/tools/gobin
build/tools/bazel \
build/tools/circleci \
build/tools/gobin:
	git submodule update --init --recursive '$@'

include build/tools/bazel/Makefile
build/tools/bazel/Makefile: build/tools/bazel
	@# included in submodule: build/tools/bazel
include build/tools/circleci/Makefile
build/tools/circleci/Makefile: build/tools/circleci
	@# included in submodule: build/tools/circleci
include build/tools/gobin/Makefile
build/tools/gobin/Makefile: build/tools/gobin
	@# included in submodule: build/tools/gobin

# bazel-info: display information about the Bazel server.
.PHONY: bazel-info
bazel-info: $(BAZEL)
	$(BAZEL) $(BAZEL_FLAGS) info

# bazel-build: build the entire project using `bazel`.
.PHONY: bazel-build
bazel-build: $(BAZEL)
	$(BAZEL) $(BAZEL_FLAGS) build //...

# bazel-test: run all tests in the entire project using `bazel`.
.PHONY: bazel-test
bazel-test: $(BAZEL)
	$(BAZEL) $(BAZEL_FLAGS) test //...

GOBIN_CACHE_DIR := build/.gobincache
$(GOBIN_CACHE_DIR):
	mkdir --parent '$@'

GAZELLE_VERSION := 0.18.2
GAZELLE_CACHE_DIR := $(GOBIN_CACHE_DIR)/gazelle/$(GAZELLE_VERSION)
GAZELLE := $(GAZELLE_CACHE_DIR)/gazelle

$(GAZELLE_CACHE_DIR):
	mkdir --parent '$@'

$(GAZELLE): $(GOBIN) | $(GAZELLE_CACHE_DIR)
	GOBIN=$(GAZELLE_CACHE_DIR) $(GOBIN) github.com/bazelbuild/bazel-gazelle/cmd/gazelle@$(GAZELLE_VERSION)

# gazelle: generate `BUILD.bazel` files using `gazelle`.
.PHONY: gazelle
gazelle: $(GAZELLE)
	$(GAZELLE) fix
	$(GAZELLE) update-repos \
		-from_file=go.mod \
		-prune

# lint-gazelle: check that re-generating build files does not create any modifications to committed files.
.PHONY: lint-gazelle
lint-gazelle: gazelle
	scripts/git-verify-no-diff.bash \
		WORKSPACE \
		$(BUILD_FILES)

GOFUMPORTS_VERSION := 96300e3d49fbb3b7bc9c6dc74f8a5cc0ef46f84b
GOFUMPORTS_CACHE_DIR := $(GOBIN_CACHE_DIR)/gofumports/$(GOFUMPORTS_VERSION)
GOFUMPORTS := $(GOFUMPORTS_CACHE_DIR)/gofumports

$(GOFUMPORTS_CACHE_DIR):
	mkdir --parent '$@'

$(GOFUMPORTS): $(GOBIN) | $(GOFUMPORTS_CACHE_DIR)
	GOBIN=$(GOFUMPORTS_CACHE_DIR) $(GOBIN) mvdan.cc/gofumpt/gofumports@$(GOFUMPORTS_VERSION)

# gofumports: run the `gofumports` Go code formatter.
.PHONY: gofumports
gofumports: $(GOFUMPORTS)
	$(GOFUMPORTS) -w '$(WD)'

# lint-gofumports: check that formatting Go files does not create any modifications to committed files.
.PHONY: lint-gofumports
lint-gofumports: gofumports
	scripts/git-verify-no-diff.bash \
		$(GO_FILES)

BUILDIFIER_VERSION := 0.28.0
BUILDIFIER_CACHE_DIR := $(GOBIN_CACHE_DIR)/buildifier/$(BUILDIFIER_VERSION)
BUILDIFIER := $(BUILDIFIER_CACHE_DIR)/buildifier

$(BUILDIFIER_CACHE_DIR):
	mkdir --parent '$@'

$(BUILDIFIER): $(GOBIN) | $(BUILDIFIER_CACHE_DIR)
	GOBIN=$(BUILDIFIER_CACHE_DIR) $(GOBIN) github.com/bazelbuild/buildtools/buildifier@$(BUILDIFIER_VERSION)

# buildifier: run the `buildifier` tool to format Bazel build files.
.PHONY: buildifier
buildifier: $(BUILDIFIER)
	$(BUILDIFIER) -r '$(WD)'

# lint-buildifier: check that formatting build files does not create any modifications to committed files.
.PHONY: lint-buildifier
lint-buildifier: buildifier $(BUILDIFIER)
	scripts/git-verify-no-diff.bash \
		WORKSPACE \
		$(BUILD_FILES)
	$(BUILDIFIER) \
		--lint=warn \
		-r '$(WD)'

GOLANGCI_LINT_VERSION := v1.17.1
GOLANGCI_LINT_CACHE_DIR := $(GOBIN_CACHE_DIR)/golangci-lint/$(GOLANGCI_LINT_VERSION)
GOLANGCI_LINT := $(GOLANGCI_LINT_CACHE_DIR)/golangci-lint

$(GOLANGCI_LINT_CACHE_DIR):
	mkdir --parent '$@'

$(GOLANGCI_LINT): $(GOBIN) | $(GOLANGCI_LINT_CACHE_DIR)
	GOBIN=$(GOLANGCI_LINT_CACHE_DIR) $(GOBIN) github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

# lint-golangci-lint: lint Go files using a number of linters.
.PHONY: lint-golangci-lint
lint-golangci-lint: $(GOLANGCI_LINT)
	@# interfacer: disabled since its author has deprecated it
	$(GOLANGCI_LINT) run \
		--enable-all \
		--disable interfacer

# circleci-build: run the `build` job using a local CircleCI executor.
.PHONY: circleci-build
circleci-build: $(CIRCLECI)
	$(CIRCLECI) local execute --job build

# go-mod-fix: update and format `go.mod` and `go.sum` files.
.PHONY: go-mod-fix
go-mod-fix:
	go mod tidy -v
	go mod edit -fmt

# go-lint-mod: make sure the `go.mod` and `go.sum` files are properly formatted and up to date.
.PHONY: go-lint-mod
go-lint-mod-fix: go-mod-fix
	scripts/git-verify-no-diff.bash \
		go.mod \
		go.sum
