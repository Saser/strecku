//+build integrationtest

package main

import (
	testingv1 "github.com/Saser/strecku/backend/gen/api/testing/v1"
	"github.com/Saser/strecku/backend/internal/impl/inmemory"
	"google.golang.org/grpc"
)

func registerServers(server *grpc.Server, impl *inmemory.Impl) {
	registerCommon(server, impl)
	testingv1.RegisterResetAPIServer(server, impl)
}
