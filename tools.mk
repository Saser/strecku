# tools.mk: rules for installing tools used by this project.

tools := tools
$(tools):
	mkdir -p '$@'

# protoc: the protobuf compiler
protoc := $(tools)/protoc
protoc_version := 3.13.0
protoc_versioned_dir := $(tools)/protoc-$(protoc_version)
protoc_versioned_archive := $(protoc_versioned_dir).zip
protoc_versioned := $(protoc_versioned_dir)/bin/protoc
# TODO: make this more platform-independent (`linux` is specified in the archive URL.)
$(protoc_versioned_archive): | $(tools)
	curl \
		--fail \
		--location \
		--show-error \
		--silent \
		--output '$@' \
		'https://github.com/protocolbuffers/protobuf/releases/download/v$(protoc_version)/protoc-$(protoc_version)-linux-x86_64.zip'

$(protoc_versioned_dir): $(protoc_versioned_archive)
	unzip \
		'$<' \
		-d '$@'

$(protoc_versioned): $(protoc_versioned_dir)

$(protoc): $(protoc_versioned)
	ln \
		--symbolic \
		--relative \
		'$<' \
		'$@'

# api-linter: Go tool to lint protobuf services
api-linter := $(tools)/api-linter
$(api-linter): go.mod go.sum
	go \
		build \
		-o='$@' \
		github.com/googleapis/api-linter/cmd/api-linter

# protoc-gen-go: protoc plugin to generate Go code for protobufs.
protoc-gen-go := $(tools)/protoc-gen-go
$(protoc-gen-go): go.mod go.sum
	go \
		build \
		-o='$@' \
		google.golang.org/protobuf/cmd/protoc-gen-go

# protoc-gen-go-grpc: protoc plugin to generate Go code for protobuf services.
protoc-gen-go-grpc := $(tools)/protoc-gen-go-grpc
$(protoc-gen-go-grpc): go.mod go.sum
	go \
		build \
		-o='$@' \
		google.golang.org/grpc/cmd/protoc-gen-go-grpc
