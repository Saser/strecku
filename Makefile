include tools.mk

proto_files := $(wildcard api/saser/strecku/v1/*.proto)
go_module := $(shell go list -m)

.PHONY: lint
lint: \
	$(api-linter) \
	$(proto-files)
lint:
	$(api-linter) \
		--proto-path=third_party/api-common-protos \
		--config=.api-linter.yml \
		$(proto_files)

.PHONY: generate
generate: \
	$(proto_files) \
	$(protoc) \
	$(protoc-gen-go) \
	$(protoc-gen-go-grpc)
generate:
	$(protoc) \
		--proto_path=api \
		--proto_path=third_party/api-common-protos \
		--plugin='$(protoc-gen-go)' \
		--go_out=. \
		--go_opt=module='$(go_module)' \
		--plugin='$(protoc-gen-go-grpc)' \
		--go-grpc_out=. \
		--go-grpc_opt=module='$(go_module)' \
		$(proto_files)
