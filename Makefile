include tools.mk

proto_files := $(wildcard saser/strecku/v1/*.proto)

server/testcert.key server/testcert.crt:
	openssl \
		req \
		-x509 \
		-sha256 \
		-days 3650 \
		-subj '/CN=localhost' \
		-addext 'subjectAltName=DNS:localhost' \
		-newkey rsa:4096 \
		-nodes \
		-keyout 'server/testcert.key' \
		-out 'server/testcert.crt'

.PHONY: testcert
testcert: server/testcert.key server/testcert.crt

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
	$(protoc-gen-go) \
	$(protoc-gen-go-grpc)
generate:
	$(protoc) \
		--plugin='$(protoc-gen-go)' \
		--go_out=. \
		--plugin='$(protoc-gen-go-grpc)' \
		--go-grpc_out=. \
		$(proto_files)
