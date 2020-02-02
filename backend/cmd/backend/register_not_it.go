//+build !integrationtest

package main

import (
	"github.com/Saser/strecku/backend/internal/impl/inmemory"
	"google.golang.org/grpc"
)

func registerServers(server *grpc.Server, impl *inmemory.Impl) {
	registerCommon(server, impl)
}
