package main

import (
	streckuv1 "github.com/Saser/strecku/backend/gen/api/strecku/v1"
	"github.com/Saser/strecku/backend/internal/impl/inmemory"
	"google.golang.org/grpc"
)

func registerCommon(server *grpc.Server, impl *inmemory.Impl) {
	streckuv1.RegisterUserAPIServer(server, impl)
	streckuv1.RegisterStoreAPIServer(server, impl)
	streckuv1.RegisterRoleAPIServer(server, impl)
	streckuv1.RegisterProductAPIServer(server, impl)
}
