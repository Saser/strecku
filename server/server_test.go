package server

import (
	"context"
	"net"
	"testing"

	streckuv1 "github.com/Saser/strecku/saser/strecku/v1"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/testing/protocmp"
)

const bufSize = 1024 * 1024

const (
	certFile = "testcert.crt"
	keyFile  = "testcert.key"
)

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
	creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
	if err != nil {
		t.Fatal(err)
	}
	s := grpc.NewServer(grpc.Creds(creds))
	streckuv1.RegisterStreckUServer(s, f.srv)
	go func() {
		if err := s.Serve(f.lis); err != nil {
			t.Errorf("s.Serve(f.lis) = %v", err)
		}
	}()
	t.Cleanup(s.GracefulStop)
	return f
}

func (f *fixture) client(ctx context.Context, t *testing.T, opts ...grpc.DialOption) streckuv1.StreckUClient {
	t.Helper()
	dial := func(context.Context, string) (net.Conn, error) { return f.lis.Dial() }
	opts = append(opts, grpc.WithContextDialer(dial))
	creds, err := credentials.NewClientTLSFromFile(certFile, "")
	if err != nil {
		t.Fatal(err)
	}
	opts = append(opts, grpc.WithTransportCredentials(creds))
	cc, err := grpc.DialContext(ctx, "localhost", opts...)
	if err != nil {
		t.Fatal(err)
	}
	return streckuv1.NewStreckUClient(cc)
}

func (f *fixture) backdoorCreateUser(t *testing.T, req *streckuv1.CreateUserRequest) *streckuv1.User {
	t.Helper()
	user := req.User
	user.Name = newUserName()
	f.srv.users = append(f.srv.users, user)
	f.srv.userKeys[user.Name] = len(f.srv.users) - 1
	t.Cleanup(func() {
		delete(f.srv.userKeys, user.Name)
	})
	f.srv.passwords[user.Name] = req.Password
	t.Cleanup(func() {
		delete(f.srv.passwords, user.Name)
	})
	return user
}

func TestServer_AuthenticateUser(t *testing.T) {
	ctx := context.Background()

	f := setUp(t)
	password := "password"
	user := f.backdoorCreateUser(t, &streckuv1.CreateUserRequest{
		User: &streckuv1.User{
			EmailAddress: "user@example.com",
			DisplayName:  "User",
		},
		Password: password,
	})
	client := f.client(ctx, t)

	for _, test := range []struct {
		name     string
		req      *streckuv1.AuthenticateUserRequest
		wantUser *streckuv1.User
		wantCode codes.Code
	}{
		{name: "Correct", req: &streckuv1.AuthenticateUserRequest{EmailAddress: user.EmailAddress, Password: password}, wantUser: user},
		{name: "IncorrectEmailAddress", req: &streckuv1.AuthenticateUserRequest{EmailAddress: "incorrect@example.com", Password: password}, wantCode: codes.Unauthenticated},
		{name: "IncorrectPassword", req: &streckuv1.AuthenticateUserRequest{EmailAddress: user.EmailAddress, Password: "incorrect password"}, wantCode: codes.Unauthenticated},
	} {
		t.Run(test.name, func(t *testing.T) {
			authUser, err := client.AuthenticateUser(ctx, test.req)
			if got := status.Code(err); got != test.wantCode {
				t.Errorf("status.Code(%v) = %v; want %v", err, got, test.wantCode)
			}
			if diff := cmp.Diff(authUser, test.wantUser, protocmp.Transform()); diff != "" {
				t.Errorf("authUser != test.wantUser (-got +want):\n%v", diff)
			}
		})
	}
}
