package main

import (
	"log"
	"net"

	streckuv1 "github.com/Saser/strecku/backend/gen/api/v1"
	"github.com/Saser/strecku/backend/internal/impl/inmemory"
	"google.golang.org/grpc"
)

func main() {
	impl := inmemory.New()
	server := grpc.NewServer()
	streckuv1.RegisterUserAPIServer(server, impl)
	streckuv1.RegisterStoreAPIServer(server, impl)
	streckuv1.RegisterRoleAPIServer(server, impl)
	streckuv1.RegisterProductAPIServer(server, impl)
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("error: %+v", err)
	}
	if err := server.Serve(listener); err != nil {
		log.Fatalf("error: %+v", err)
	}
	log.Println("goodbye!")
}
