package server

import (
	"context"
	"net"
	"testing"

	"github.com/Saser/strecku/auth"
	streckuv1 "github.com/Saser/strecku/saser/strecku/v1"
	"github.com/google/go-cmp/cmp"
	"golang.org/x/crypto/bcrypt"
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

func (f *fixture) authClient(ctx context.Context, t *testing.T, emailAddress, password string) streckuv1.StreckUClient {
	t.Helper()
	return f.client(ctx, t, grpc.WithPerRPCCredentials(auth.Basic{
		Username: emailAddress,
		Password: password,
	}))
}

func (f *fixture) backdoorCreateUser(t *testing.T, req *streckuv1.CreateUserRequest) *streckuv1.User {
	t.Helper()
	user := req.User
	user.Name = newUserName()
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err)
	}
	f.srv.users = append(f.srv.users, &userEntry{
		user: user,
		hash: hash,
	})
	f.srv.userIndices[user.Name] = len(f.srv.users) - 1
	t.Cleanup(func() {
		delete(f.srv.userIndices, user.Name)
	})
	f.srv.userKeys[user.EmailAddress] = user.Name
	t.Cleanup(func() {
		delete(f.srv.userKeys, user.EmailAddress)
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
		{name: "MissingEmailAddress", req: &streckuv1.AuthenticateUserRequest{EmailAddress: "", Password: password}, wantCode: codes.InvalidArgument},
		{name: "MissingPassword", req: &streckuv1.AuthenticateUserRequest{EmailAddress: user.EmailAddress, Password: ""}, wantCode: codes.InvalidArgument},
		{name: "EmptyRequest", req: &streckuv1.AuthenticateUserRequest{}, wantCode: codes.InvalidArgument},
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

func TestServer_CreateUser_AsSuperuser(t *testing.T) {
	ctx := context.Background()

	f := setUp(t)
	rootPassword := "root password"
	root := f.backdoorCreateUser(t, &streckuv1.CreateUserRequest{
		User: &streckuv1.User{
			EmailAddress: "root@example.com",
			DisplayName:  "Root",
			Superuser:    true,
		},
		Password: rootPassword,
	})
	client := f.authClient(ctx, t, root.EmailAddress, rootPassword)

	// These test cases are for valid requests (but not necessarily requests
	// that successfully create a user).
	t.Run("Valid", func(t *testing.T) {
		for _, test := range []struct {
			name string
			req  *streckuv1.CreateUserRequest
		}{
			{
				name: "NormalUser",
				req: &streckuv1.CreateUserRequest{
					User: &streckuv1.User{
						EmailAddress: "user@example.com",
						DisplayName:  "User",
					},
					Password: "user password",
				},
			},
			{
				name: "Superuser",
				req: &streckuv1.CreateUserRequest{
					User: &streckuv1.User{
						EmailAddress: "another-root@example.com",
						DisplayName:  "Another root",
						Superuser:    true,
					},
					Password: "another root password",
				},
			},
		} {
			t.Run(test.name, func(t *testing.T) {
				user, err := client.CreateUser(ctx, test.req)
				if got, want := status.Code(err), codes.OK; got != want {
					t.Errorf("CreateUser: status.Code(%v) = %v; want %v", err, got, want)
				}
				if diff := cmp.Diff(user, test.req.User, protocmp.Transform(), protocmp.IgnoreFields(user, "name")); diff != "" {
					t.Errorf("user != test.req.User (-got +want):\n%s", diff)
				}

				// After a user has been created, it should be possible to authenticate them.
				authUser, err := client.AuthenticateUser(ctx, &streckuv1.AuthenticateUserRequest{
					EmailAddress: user.EmailAddress,
					Password:     test.req.Password,
				})
				if got, want := status.Code(err), codes.OK; got != want {
					t.Errorf("AuthenticateUser: status.Code(%v) = %v; want %v", err, got, want)
				}
				if diff := cmp.Diff(authUser, user, protocmp.Transform()); diff != "" {
					t.Errorf("authUser != user (-got +want):\n%s", diff)
				}
			})
		}
	})

	// These tests are for requests that are in some way invalid, such as having
	// missing arguments.
	t.Run("Invalid", func(t *testing.T) {
		for _, test := range []struct {
			name     string
			req      *streckuv1.CreateUserRequest
			wantCode codes.Code
		}{
			{
				name: "MissingEmailAddress",
				req: &streckuv1.CreateUserRequest{
					User: &streckuv1.User{
						EmailAddress: "",
						DisplayName:  "User",
					},
					Password: "user password",
				},
				wantCode: codes.InvalidArgument,
			},
			{
				name: "MissingDisplayName",
				req: &streckuv1.CreateUserRequest{
					User: &streckuv1.User{
						EmailAddress: "user@example.com",
						DisplayName:  "",
					},
					Password: "user password",
				},
				wantCode: codes.InvalidArgument,
			},
			{
				name: "MissingPassword",
				req: &streckuv1.CreateUserRequest{
					User: &streckuv1.User{
						EmailAddress: "user@example.com",
						DisplayName:  "User",
					},
					Password: "",
				},
				wantCode: codes.InvalidArgument,
			},
		} {
			t.Run(test.name, func(t *testing.T) {
				_, err := client.CreateUser(ctx, test.req)
				if got := status.Code(err); got != test.wantCode {
					t.Errorf("status.Code(%v) = %v; want %v", err, got, test.wantCode)
				}
			})
		}
	})
}
