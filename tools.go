//+build tools

package tools

import (
	_ "github.com/golang-migrate/migrate/v4/cmd/migrate"
	_ "github.com/googleapis/api-linter/cmd/api-linter"
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
)
