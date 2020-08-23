package server

import (
	"context"
	"net"
	"testing"

	streckuv1 "github.com/Saser/strecku/saser/strecku/v1"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/testing/protocmp"
)

const bufSize = 1024 * 1024

type fixture struct {
	srv *Server
	lis *bufconn.Listener
}

func setUp(t *testing.T) *fixture {
	t.Helper()

	f := &fixture{
		srv: New(),
		lis: bufconn.Listen(bufSize),
	}
	s := grpc.NewServer()
	streckuv1.RegisterStreckUServer(s, f.srv)
	go func() {
		if err := s.Serve(f.lis); err != nil {
			t.Errorf("s.Serve(f.lis) = %v", err)
			t.FailNow()
		}
	}()
	t.Cleanup(s.GracefulStop)
	return f
}

func (f *fixture) dial(context.Context, string) (net.Conn, error) {
	return f.lis.Dial()
}

func (f *fixture) insecureClient(ctx context.Context, t *testing.T) streckuv1.StreckUClient {
	t.Helper()
	cc, err := grpc.DialContext(ctx, "bufconn", grpc.WithContextDialer(f.dial), grpc.WithInsecure())
	if err != nil {
		t.Fatal(err)
	}
	return streckuv1.NewStreckUClient(cc)
}

func (f *fixture) backdoorCreateUser(t *testing.T, user *streckuv1.User, password string) {
	t.Helper()
	f.srv.users = append(f.srv.users, user)
	f.srv.userKeys[user.Name] = len(f.srv.users) - 1
	t.Cleanup(func() {
		delete(f.srv.userKeys, user.Name)
	})
	f.srv.passwords[user.Name] = password
	t.Cleanup(func() {
		delete(f.srv.passwords, user.Name)
	})
}

func TestServer_AuthenticateUser(t *testing.T) {
	ctx := context.Background()

	f := setUp(t)
	emailAddress := "user@example.com"
	password := "password"
	user := &streckuv1.User{
		Name:         newUserName(),
		EmailAddress: "user@example.com",
		DisplayName:  "User",
	}
	f.backdoorCreateUser(t, user, password)
	client := f.insecureClient(ctx, t)

	t.Run("Correct", func(t *testing.T) {
		resp, err := client.AuthenticateUser(ctx, &streckuv1.AuthenticateUserRequest{
			EmailAddress: emailAddress,
			Password:     password,
		})
		if err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(user, resp.User, protocmp.Transform()); diff != "" {
			t.Errorf("user != resp.User (-want +got):\n%v", diff)
		}
		if got := resp.Token; got == "" {
			t.Errorf("resp.Token = %q, want non-empty string", got)
		}
	})
	t.Run("IncorrectEmailAddress", func(t *testing.T) {
		_, err := client.AuthenticateUser(ctx, &streckuv1.AuthenticateUserRequest{
			EmailAddress: "incorrect@user.com",
			Password:     password,
		})
		if got, want := status.Code(err), codes.Unauthenticated; got != want {
			t.Errorf("status.Code(err) = %v; want %v", got, want)
		}
	})
	t.Run("IncorrectPassword", func(t *testing.T) {
		_, err := client.AuthenticateUser(ctx, &streckuv1.AuthenticateUserRequest{
			EmailAddress: emailAddress,
			Password:     "incorrect password",
		})
		if got, want := status.Code(err), codes.Unauthenticated; got != want {
			t.Errorf("status.Code(err) = %v; want %v", got, want)
		}
	})
}
