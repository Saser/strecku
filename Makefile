.PHONY: all
all: \
	bazel-build \
	bazel-test

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
	$(BAZEL) run //:gazelle -- update-repos -from_file=go.mod

# bazel-build: build the entire project using `bazel`.
.PHONY: bazel-build
bazel-build: $(BAZEL)
	$(BAZEL) build //...

# bazel-test: run all tests in the entire project using `bazel`.
.PHONY: bazel-test
bazel-test: $(BAZEL)
	$(BAZEL) test --test_output=errors //...

# bazel-buildifier: run the `buildifier` tool to format Bazel build files.
.PHONY: bazel-buildifier
bazel-buildifier: $(BAZEL)
	$(BAZEL) run //:buildifier

# circleci-build: run the `build` job using a local CircleCI executor.
.PHONY: circleci-build
circleci-build: $(CIRCLECI)
	$(CIRCLECI) local execute --job build
