package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"strings"

	pb "github.com/Saser/strecku/api/v1"
	"github.com/Saser/strecku/internal/repositories"
	"github.com/Saser/strecku/internal/service"
	"github.com/Saser/strecku/resources/stores/memberships"
	"github.com/Saser/strecku/resources/stores/payments"
	"github.com/Saser/strecku/resources/stores/products"
	"github.com/Saser/strecku/resources/stores/purchases"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

func main() {
	srv := grpc.NewServer()
	log.Print("created gRPC server")

	svc := service.New(
		repositories.NewInMemoryUsers(),
		repositories.NewInMemoryStores(),
		memberships.NewRepository(),
		products.NewRepository(),
		purchases.NewRepository(),
		payments.NewRepository(),
	)
	log.Print("created StreckU service")

	pb.RegisterStreckUServer(srv, svc)
	log.Print("registered StreckU service on gRPC server")

	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Print(err)
		return
	}
	log.Print("listening on address :8080")
	defer func() {
		if err := lis.Close(); err != nil {
			// This specific error is not exported, which
			// is something that has caused a lot of
			// people a lot of trouble. It will be
			// exported in Go 1.16. See
			// https://github.com/golang/go/issues/4373
			// for more details.
			//
			// Why do we ignore this error? The gRPC
			// server will use this listener to accept
			// connections. When the gRPC server is
			// stopped (see the call to srv.GracefulStop()
			// below), it will also close the listener. We
			// will close it a second time in this
			// deferred function, and doing so will cause
			// the "use of closed network connection"
			// error which is safe to ignore.
			//
			// TODO(issues/5): remove this check once
			// Go 1.16 is released.
			if strings.Contains(err.Error(), "use of closed network connection") {
				return
			}
			log.Print(err)
			return
		}
		log.Print("closed clistener")
	}()

	var g errgroup.Group
	g.Go(func() error {
		log.Print("serving gRPC server")
		return srv.Serve(lis)
	})

	sigChan := make(chan os.Signal, 1)
	defer close(sigChan)
	signal.Notify(sigChan, os.Interrupt)

	<-sigChan
	log.Print("interrupt received, shutting down")
	srv.GracefulStop()

	if err := g.Wait(); err != nil {
		log.Print(err)
		return
	}
	log.Print("goodbye")
}
