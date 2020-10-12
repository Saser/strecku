package service

import (
	"context"
	"net"
	"testing"
	"time"

	pb "github.com/Saser/strecku/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

func seed(t *testing.T) *Service {
	t.Helper()
	return New()
}

func serveAndDial(ctx context.Context, t *testing.T, svc *Service) pb.StreckUClient {
	t.Helper()
	srv := grpc.NewServer()
	pb.RegisterStreckUServer(srv, svc)
	lis := bufconn.Listen(bufSize)
	go func() {
		if err := srv.Serve(lis); err != nil {
			t.Errorf("srv.Serve(%v) = %v; want nil", lis, err)
		}
	}()
	t.Cleanup(srv.GracefulStop)
	dial := func(context.Context, string) (net.Conn, error) { return lis.Dial() }
	cc, err := grpc.DialContext(
		ctx,
		"bufconn",
		grpc.WithBlock(),
		grpc.WithInsecure(),
		grpc.WithContextDialer(dial),
	)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	return pb.NewStreckUClient(cc)
}

func TestSeed(t *testing.T) {
	seed(t)
}

func TestDial(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	serveAndDial(ctx, t, seed(t))
}
