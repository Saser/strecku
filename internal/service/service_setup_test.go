package service

import (
	"context"
	"net"
	"testing"
	"time"

	pb "github.com/Saser/strecku/api/v1"
	"github.com/Saser/strecku/resources/stores"
	"github.com/Saser/strecku/resources/stores/memberships"
	"github.com/Saser/strecku/resources/stores/products"
	"github.com/Saser/strecku/resources/stores/purchases"
	"github.com/Saser/strecku/resources/testresources"
	"github.com/Saser/strecku/resources/users"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

func seed(t *testing.T) *Service {
	t.Helper()
	userRepo := users.SeedRepository(
		t,
		[]*pb.User{
			testresources.Alice,
			testresources.Bob,
		},
		[]string{
			testresources.AlicePassword,
			testresources.BobPassword,
		},
	)
	storeRepo := stores.SeedRepository(
		t,
		[]*pb.Store{
			testresources.Bar,
			testresources.Mall,
		},
	)
	membershipRepo := memberships.SeedRepository(
		t,
		[]*pb.Membership{
			testresources.Bar_Alice,
			testresources.Bar_Bob,
			testresources.Mall_Alice,
		},
	)
	productRepo := products.SeedRepository(
		t,
		[]*pb.Product{
			testresources.Beer,
			testresources.Jeans,
		},
	)
	purchaseRepo := purchases.SeedRepository(
		t,
		[]*pb.Purchase{
			testresources.Bar_Alice_Beer1,
			testresources.Mall_Alice_Jeans1,
		},
	)
	return New(userRepo, storeRepo, membershipRepo, productRepo, purchaseRepo)
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
