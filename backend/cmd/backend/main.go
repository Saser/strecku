package main

import (
	"log"
	"net"

	"github.com/Saser/strecku/backend/internal/impl/inmemory"
	"google.golang.org/grpc"
)

func main() {
	server := grpc.NewServer()
	impl := inmemory.New()
	registerServers(server, impl)
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("error: %+v", err)
	}
	if err := server.Serve(listener); err != nil {
		log.Fatalf("error: %+v", err)
	}
	log.Println("goodbye!")
}
