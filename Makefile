.PHONY: lint
lint: saser/strecku/v1/*.proto
	api-linter \
		--config=.api-linter.yml \
		$?

.PHONY: generate
generate: saser/strecku/v1/*.proto
	protoc \
		--go_out=genproto \
		$?
