.PHONY: all
all: \
	bazel-build \
	bazel-test

# WD: the absolute path to the current working directory. It is used for referring to the root directory of this project
# instead of using e.g. `.` to refer to "this directory". This is necessary when invoking tools such as `gofumports`
# using Bazel.
WD := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))

GIT_LS_FILES := git ls-files --exclude-standard --cached --others

BUILD_FILES := $(shell $(GIT_LS_FILES) -- '*BUILD.bazel')
GO_FILES := $(shell $(GIT_LS_FILES) -- '*.go')

.PHONY: \
	build/tools/bazel \
	build/tools/circleci
build/tools/bazel \
build/tools/circleci:
	git submodule update --init --recursive '$@'

include build/tools/bazel/Makefile
build/tools/bazel/Makefile: build/tools/bazel
	@# included in submodule: build/tools/bazel
include build/tools/circleci/Makefile
build/tools/circleci/Makefile: build/tools/circleci
	@# included in submodule: build/tools/circleci

# bazel-info: display information about the Bazel server.
.PHONY: bazel-info
bazel-info: $(BAZEL)
	$(BAZEL) info

# bazel-gazelle: generate `BUILD.bazel` files using `gazelle`.
.PHONY: bazel-gazelle
bazel-gazelle: $(BAZEL)
	$(BAZEL) run //:gazelle -- fix
	$(BAZEL) run //:gazelle -- update-repos -from_file=go.mod -prune

# bazel-build: build the entire project using `bazel`.
.PHONY: bazel-build
bazel-build: $(BAZEL)
	$(BAZEL) build //...

# bazel-test: run all tests in the entire project using `bazel`.
.PHONY: bazel-test
bazel-test: $(BAZEL)
	$(BAZEL) test --test_output=errors //...

# bazel-gofumports: run the `gofumports` Go code formatter.
.PHONY: bazel-gofumports
bazel-gofumports: $(BAZEL)
	$(BAZEL) run @cc_mvdan_gofumpt//gofumports -- -w $(WD)

# bazel-buildifier: run the `buildifier` tool to format Bazel build files.
.PHONY: bazel-buildifier
bazel-buildifier: $(BAZEL)
	$(BAZEL) run //:buildifier

# bazel-lint-gazelle: check that re-generating build files does not create any modifications to committed files.
.PHONY: bazel-lint-gazelle
bazel-lint-gazelle: bazel-gazelle
	scripts/git-verify-no-diff.bash \
		WORKSPACE \
		$(BUILD_FILES)

# bazel-lint-buildifier: check that formatting build files does not create any modifications to committed files.
.PHONY: bazel-lint-buildifier
bazel-lint-buildifier: bazel-buildifier $(BAZEL)
	scripts/git-verify-no-diff.bash \
		WORKSPACE \
		$(BUILD_FILES)
	$(BAZEL) $(BAZEL_FLAGS) run //:buildifier -- --lint=warn

# bazel-lint-gofumports: check that formatting Go files does not create any modifications to committed files.
.PHONY: bazel-lint-gofumports
bazel-lint-gofumports: bazel-gofumports
	scripts/git-verify-no-diff.bash \
		$(GO_FILES)

# circleci-build: run the `build` job using a local CircleCI executor.
.PHONY: circleci-build
circleci-build: $(CIRCLECI)
	$(CIRCLECI) local execute --job build
