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
	gofumports \
	go-mod-fix

# lint: run linters for build files, Go files, etc.
.PHONY: lint
lint: \
	lint-gofumports \
	lint-golangci-lint \
	go-lint-mod-fix

# build: build the entire project.
.PHONY: build
build: \
	go-build

# test: test the entire project.
.PHONY: test
test: \
	go-test

# `tools.mk` contains targets and recipes for building and installing all tools required by this Makefile.
include tools.mk

# WD: the absolute path to the current working directory. It is used for referring to the root directory of this project
# instead of using e.g. `.` to refer to "this directory". This is not strictly necessary for most tools, but it improves
# the reliability of invoking them by using absolute paths instead of relative paths.
WD := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))

GIT_LS_FILES := git ls-files --exclude-standard --cached --others

GO_FILES := $(shell $(GIT_LS_FILES) -- '*.go')

build/tools/circleci/Makefile:
	git submodule update --init --recursive '$(dir $@)'
include build/tools/circleci/Makefile

# gofumports: run the `gofumports` Go code formatter.
.PHONY: gofumports
gofumports: $(GOFUMPORTS)
	$(GOFUMPORTS) -w '$(WD)'

# lint-gofumports: check that formatting Go files does not create any modifications to committed files.
.PHONY: lint-gofumports
lint-gofumports: gofumports
	scripts/git-verify-no-diff.bash \
		$(GO_FILES)

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

# go-build: builds the module using `go build`.
.PHONY: go-build
go-build:
	go build ./...

# go-test: tests the module using `go test`.
.PHONY: go-test
go-test:
	go test -race -cover ./...

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
