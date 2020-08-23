include tools.mk

proto_files := $(wildcard saser/strecku/v1/*.proto)

.PHONY: lint
lint: \
	$(api-linter) \
	$(proto-files)
lint:
	$(api-linter) \
		--config=.api-linter.yml \
		$(proto_files)

.PHONY: generate
generate: \
	$(proto_files) \
	$(protoc) \
	$(protoc-gen-go)
generate:
	$(protoc) \
		--plugin='$(protoc-gen-go)' \
		--go_out=genproto \
		$(proto_files)
