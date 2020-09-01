include tools.mk

proto_files := $(wildcard saser/strecku/v1/*.proto)

server/testdata/cert.key server/testdata/cert.crt:
	openssl \
		req \
		-x509 \
		-sha256 \
		-days 3650 \
		-subj '/CN=localhost' \
		-addext 'subjectAltName=DNS:localhost' \
		-newkey rsa:4096 \
		-nodes \
		-keyout 'server/testdata/cert.key' \
		-out 'server/testdata/cert.crt'

.PHONY: testcert
testcert: server/testdata/cert.key server/testdata/cert.crt

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
		--go_opt=module=github.com/Saser/strecku \
		--plugin='$(protoc-gen-go-grpc)' \
		--go-grpc_out=. \
		--go-grpc_opt=module=github.com/Saser/strecku \
		$(proto_files)
