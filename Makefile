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
	go-mod-fix \
	prototool-format \
	prototool-generate

# lint: run linters for build files, Go files, etc.
.PHONY: lint
lint: \
	lint-gofumports \
	lint-golangci-lint \
	go-lint-mod-fix \
	lint-prototool \
	lint-prototool-generate

# build: build the entire project.
.PHONY: build
build: \
	go-build \
	prototool-compile

# test: test the entire project.
.PHONY: test
test: \
	go-test

# WD: the absolute path to the current working directory. It is used for referring to the root directory of this project
# instead of using e.g. `.` to refer to "this directory". This is not strictly necessary for most tools, but it improves
# the reliability of invoking them by using absolute paths instead of relative paths.
WD := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))

# `tools.mk` contains targets and recipes for building and installing all tools required by this Makefile.
include tools.mk

GIT_LS_FILES := git ls-files --exclude-standard --cached --others

# GO_FILES contains all Go files not ignored by Git, except for:
#   * Go files generated from Protobuf files (:!:*.pb.go)
GO_FILES := $(shell $(GIT_LS_FILES) -- '*.go' ':!:*.pb.go')
# PB_GO_FILES contains all Go files not ignored by Git which are generated from Protobuf files.
PB_GO_FILES := $(shell $(GIT_LS_FILES) -- '*.pb.go')
# PROTO_FILES contains all Protobuf files not ignored by Git.
PROTO_FILES := $(shell $(GIT_LS_FILES) -- '*.proto')

build/tools/circleci/Makefile:
	git submodule update --init --recursive '$(dir $@)'
include build/tools/circleci/Makefile

# gofumports: run the `gofumports` Go code formatter.
.PHONY: gofumports
gofumports: $(GOFUMPORTS)
	$(GOFUMPORTS) -w $(GO_FILES)

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

# prototool-compile: make sure Protobuf files compile, but do not generate code.
.PHONY: prototool-compile
prototool-compile: $(PROTOTOOL)
	$(PROTOTOOL) compile

# prototool-format: format Protobuf files according to `prototool`.
.PHONY: prototool-format
prototool-format: $(PROTOTOOL)
	$(PROTOTOOL) format --fix --overwrite

# prototool-generate: generate implementation code for Protobuf files.
# `protoc-gen-go` needs to in $PATH in order for `prototool` to be able to use it.
.PHONY: prototool-generate
prototool-generate: $(PROTOTOOL) $(PROTOC_GEN_GO)
	PATH=$(WD)/$(dir $(PROTOC_GEN_GO)):$(PATH) $(PROTOTOOL) generate

# lint-prototool-generate: verify that Go files generated from Protobuf files are up to date.
.PHONY: lint-prototool-generate
lint-prototool-generate: prototool-generate
	scripts/git-verify-no-diff.bash \
		$(PB_GO_FILES)

# lint-prototool-lint: lint Protobuf files using `prototool`.
.PHONY: lint-prototool
lint-prototool: $(PROTOTOOL) prototool-format
	scripts/git-verify-no-diff.bash \
		$(PROTO_FILES)
	$(PROTOTOOL) lint

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
