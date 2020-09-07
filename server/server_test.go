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
	certFile = "testdata/cert.crt"
	keyFile  = "testdata/cert.key"
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

func TestServer_GetUser(t *testing.T) {
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
	otherRoot := f.backdoorCreateUser(t, &streckuv1.CreateUserRequest{
		User: &streckuv1.User{
			EmailAddress: "other-root@example.com",
			DisplayName:  "Other Root",
			Superuser:    true,
		},
		Password: "other root password",
	})
	userPassword := "user password"
	user := f.backdoorCreateUser(t, &streckuv1.CreateUserRequest{
		User: &streckuv1.User{
			EmailAddress: "user@example.com",
			DisplayName:  "User",
			Superuser:    false,
		},
		Password: userPassword,
	})
	otherUser := f.backdoorCreateUser(t, &streckuv1.CreateUserRequest{
		User: &streckuv1.User{
			EmailAddress: "other-user@example.com",
			DisplayName:  "Other User",
			Superuser:    false,
		},
		Password: "other user password",
	})

	// A superuser should be able to get any other (existing) user.
	t.Run("AsSuperuser", func(t *testing.T) {
		client := f.authClient(ctx, t, root.EmailAddress, rootPassword)
		for _, test := range []struct {
			name     string
			req      *streckuv1.GetUserRequest
			wantUser *streckuv1.User
			wantCode codes.Code
		}{
			{name: "Me", req: &streckuv1.GetUserRequest{Name: "users/me"}, wantUser: root},
			{name: "Root", req: &streckuv1.GetUserRequest{Name: root.Name}, wantUser: root},
			{name: "OtherRoot", req: &streckuv1.GetUserRequest{Name: otherRoot.Name}, wantUser: otherRoot},
			{name: "User", req: &streckuv1.GetUserRequest{Name: user.Name}, wantUser: user},
			{name: "OtherUser", req: &streckuv1.GetUserRequest{Name: otherUser.Name}, wantUser: otherUser},
			{name: "Nonexistent", req: &streckuv1.GetUserRequest{Name: "nonexistent"}, wantCode: codes.NotFound},
			{name: "EmptyName", req: &streckuv1.GetUserRequest{Name: ""}, wantCode: codes.InvalidArgument},
		} {
			t.Run(test.name, func(t *testing.T) {
				gotUser, err := client.GetUser(ctx, test.req)
				if gotCode := status.Code(err); gotCode != test.wantCode {
					t.Errorf("status.Code(%v) = %v; want %v", err, gotCode, test.wantCode)
				}
				if test.wantCode != codes.OK {
					return
				}
				if diff := cmp.Diff(gotUser, test.wantUser, protocmp.Transform()); diff != "" {
					t.Errorf("got != want (-got +want):\n%s", diff)
				}
			})
		}
	})

	// TODO: a non-superuser A should eventually be able to get another user
	// B, if B is a member of the store that A is an administrator of. For
	// now, however, a non-superuser is only able to get themselves. All
	// other (valid) requests should result in a PermissionDenied error.
	t.Run("AsNormalUser", func(t *testing.T) {
		client := f.authClient(ctx, t, user.EmailAddress, userPassword)
		for _, test := range []struct {
			name     string
			req      *streckuv1.GetUserRequest
			wantUser *streckuv1.User
			wantCode codes.Code
		}{
			{name: "Me", req: &streckuv1.GetUserRequest{Name: "users/me"}, wantUser: user},
			{name: "Root", req: &streckuv1.GetUserRequest{Name: root.Name}, wantCode: codes.PermissionDenied},
			{name: "OtherRoot", req: &streckuv1.GetUserRequest{Name: otherRoot.Name}, wantCode: codes.PermissionDenied},
			{name: "User", req: &streckuv1.GetUserRequest{Name: user.Name}, wantUser: user},
			{name: "OtherUser", req: &streckuv1.GetUserRequest{Name: otherUser.Name}, wantCode: codes.PermissionDenied},
			{name: "Nonexistent", req: &streckuv1.GetUserRequest{Name: "nonexistent"}, wantCode: codes.PermissionDenied},
			{name: "Empty", req: &streckuv1.GetUserRequest{Name: ""}, wantCode: codes.InvalidArgument},
		} {
			t.Run(test.name, func(t *testing.T) {
				gotUser, err := client.GetUser(ctx, test.req)
				if gotCode := status.Code(err); gotCode != test.wantCode {
					t.Errorf("status.Code(%v) = %v; want %v", err, gotCode, test.wantCode)
				}
				if test.wantCode != codes.OK {
					return
				}
				if diff := cmp.Diff(gotUser, test.wantUser, protocmp.Transform()); diff != "" {
					t.Errorf("got != want (-got +want):\n%s", diff)
				}
			})
		}
	})
}

func TestServer_ListUsers(t *testing.T) {
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
	userPassword := "user password"
	user := f.backdoorCreateUser(t, &streckuv1.CreateUserRequest{
		User: &streckuv1.User{
			EmailAddress: "user@example.com",
			DisplayName:  "User",
			Superuser:    false,
		},
		Password: userPassword,
	})

	// A superuser should be able to list all users.
	t.Run("AsSuperuser", func(t *testing.T) {
		client := f.authClient(ctx, t, root.EmailAddress, rootPassword)
		for _, test := range []struct {
			name     string
			req      *streckuv1.ListUsersRequest
			wantResp *streckuv1.ListUsersResponse
			wantCode codes.Code
		}{
			{
				name: "OK",
				req:  &streckuv1.ListUsersRequest{},
				wantResp: &streckuv1.ListUsersResponse{
					Users: []*streckuv1.User{root, user},
				},
			},
			{name: "NegativePageSize", req: &streckuv1.ListUsersRequest{PageSize: -1}, wantCode: codes.InvalidArgument},
			{name: "NonEmptyPageToken", req: &streckuv1.ListUsersRequest{PageToken: "token"}, wantCode: codes.Unimplemented},
			{name: "PositivePageSize", req: &streckuv1.ListUsersRequest{PageSize: 1}, wantCode: codes.Unimplemented},
		} {
			t.Run(test.name, func(t *testing.T) {
				gotResp, err := client.ListUsers(ctx, test.req)
				if gotCode := status.Code(err); gotCode != test.wantCode {
					t.Errorf("status.Code(%v) = %v; want %v", err, gotCode, test.wantCode)
				}
				if test.wantCode != codes.OK {
					return
				}
				if diff := cmp.Diff(
					gotResp, test.wantResp, protocmp.Transform(),
					protocmp.SortRepeatedFields(test.wantResp, "users"),
				); diff != "" {
					t.Errorf("-got +want:\n%s", diff)
				}
			})
		}
	})

	// TODO: a non-superuser A should eventually be able to list all users that
	// are members of any store that A is an administrator for. For now,
	// however, A is only able to list themselves.
	t.Run("AsNormalUser", func(t *testing.T) {
		client := f.authClient(ctx, t, user.EmailAddress, userPassword)
		for _, test := range []struct {
			name     string
			req      *streckuv1.ListUsersRequest
			wantResp *streckuv1.ListUsersResponse
			wantCode codes.Code
		}{
			{
				name: "OK",
				req:  &streckuv1.ListUsersRequest{},
				wantResp: &streckuv1.ListUsersResponse{
					Users: []*streckuv1.User{user},
				},
			},
			{name: "NegativePageSize", req: &streckuv1.ListUsersRequest{PageSize: -1}, wantCode: codes.InvalidArgument},
			{name: "NonEmptyPageToken", req: &streckuv1.ListUsersRequest{PageToken: "token"}, wantCode: codes.Unimplemented},
			{name: "PositivePageSize", req: &streckuv1.ListUsersRequest{PageSize: 1}, wantCode: codes.Unimplemented},
		} {
			t.Run(test.name, func(t *testing.T) {
				gotResp, err := client.ListUsers(ctx, test.req)
				if gotCode := status.Code(err); gotCode != test.wantCode {
					t.Errorf("status.Code(%v) = %v; want %v", err, gotCode, test.wantCode)
				}
				if test.wantCode != codes.OK {
					return
				}
				if diff := cmp.Diff(
					gotResp, test.wantResp, protocmp.Transform(),
					protocmp.SortRepeatedFields(test.wantResp, "users"),
				); diff != "" {
					t.Errorf("-got +want:\n%s", diff)
				}
			})
		}
	})
}

func TestServer_CreateUser(t *testing.T) {
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
	userPassword := "user password"
	user := f.backdoorCreateUser(t, &streckuv1.CreateUserRequest{
		User: &streckuv1.User{
			EmailAddress: "user@example.com",
			DisplayName:  "User",
			Superuser:    false,
		},
		Password: userPassword,
	})

	// A superuser can always create a new user, given that the request is valid.
	t.Run("AsSuperuser", func(t *testing.T) {
		client := f.authClient(ctx, t, root.EmailAddress, rootPassword)
		for _, test := range []struct {
			name     string
			req      *streckuv1.CreateUserRequest
			wantCode codes.Code
		}{
			{
				name: "CreateSuperuser",
				req: &streckuv1.CreateUserRequest{
					User: &streckuv1.User{
						EmailAddress: "other-root@example.com",
						DisplayName:  "Other Root",
						Superuser:    true,
					},
					Password: "other root password",
				},
			},
			{
				name: "CreateNormalUser",
				req: &streckuv1.CreateUserRequest{
					User: &streckuv1.User{
						EmailAddress: "other-user@example.com",
						DisplayName:  "Other User",
						Superuser:    false,
					},
					Password: "other user password",
				},
			},
			{
				name: "DuplicateEmailAddress",
				req: &streckuv1.CreateUserRequest{
					User: &streckuv1.User{
						EmailAddress: user.EmailAddress,
						DisplayName:  "User Again",
						Superuser:    false,
					},
					Password: userPassword,
				},
				wantCode: codes.AlreadyExists,
			},
			{
				name: "EmptyEmailAddress",
				req: &streckuv1.CreateUserRequest{
					User: &streckuv1.User{
						EmailAddress: "",
						DisplayName:  "User",
					},
					Password: userPassword,
				},
				wantCode: codes.InvalidArgument,
			},
			{
				name: "EmptyDisplayName",
				req: &streckuv1.CreateUserRequest{
					User: &streckuv1.User{
						EmailAddress: "other-user@example.com",
						DisplayName:  "",
					},
					Password: userPassword,
				},
				wantCode: codes.InvalidArgument,
			},
			{
				name: "EmptyPassword",
				req: &streckuv1.CreateUserRequest{
					User: &streckuv1.User{
						EmailAddress: "other-user@example.com",
						DisplayName:  "Other User",
					},
					Password: "",
				},
				wantCode: codes.InvalidArgument,
			},
		} {
			t.Run(test.name, func(t *testing.T) {
				got, err := client.CreateUser(ctx, test.req)
				if gotCode := status.Code(err); gotCode != test.wantCode {
					t.Fatalf("status.Code(%v) = %v; want %v", err, gotCode, test.wantCode)
				}
				if diff := cmp.Diff(
					got, test.req.User, protocmp.Transform(),
					protocmp.IgnoreFields(got, "name"),
				); test.wantCode == codes.OK && diff != "" {
					t.Errorf("-got +want:\n%s", diff)
				}
			})
		}
	})

	// A normal user can never create a new user, and should always receive
	// PermissionDenied errors (for valid requests).
	t.Run("AsNormalUser", func(t *testing.T) {
		client := f.authClient(ctx, t, user.EmailAddress, userPassword)
		for _, test := range []struct {
			name string
			req  *streckuv1.CreateUserRequest
			want codes.Code
		}{
			{
				name: "CreateSuperuser",
				req: &streckuv1.CreateUserRequest{
					User: &streckuv1.User{
						EmailAddress: "other-root@example.com",
						DisplayName:  "Other Root",
						Superuser:    true,
					},
					Password: "other root password",
				},
				want: codes.PermissionDenied,
			},
			{
				name: "CreateNormalUser",
				req: &streckuv1.CreateUserRequest{
					User: &streckuv1.User{
						EmailAddress: "other-user@example.com",
						DisplayName:  "Other User",
						Superuser:    false,
					},
					Password: "other user password",
				},
				want: codes.PermissionDenied,
			},
			{
				name: "DuplicateEmailAddress",
				req: &streckuv1.CreateUserRequest{
					User: &streckuv1.User{
						EmailAddress: user.EmailAddress,
						DisplayName:  "User Again",
						Superuser:    false,
					},
					Password: userPassword,
				},
				want: codes.PermissionDenied,
			},
			{
				name: "EmptyEmailAddress",
				req: &streckuv1.CreateUserRequest{
					User: &streckuv1.User{
						EmailAddress: "",
						DisplayName:  "User",
					},
					Password: userPassword,
				},
				want: codes.InvalidArgument,
			},
			{
				name: "EmptyDisplayName",
				req: &streckuv1.CreateUserRequest{
					User: &streckuv1.User{
						EmailAddress: "other-user@example.com",
						DisplayName:  "",
					},
					Password: userPassword,
				},
				want: codes.InvalidArgument,
			},
			{
				name: "EmptyPassword",
				req: &streckuv1.CreateUserRequest{
					User: &streckuv1.User{
						EmailAddress: "other-user@example.com",
						DisplayName:  "Other User",
					},
					Password: "",
				},
				want: codes.InvalidArgument,
			},
		} {
			t.Run(test.name, func(t *testing.T) {
				_, err := client.CreateUser(ctx, test.req)
				if got := status.Code(err); got != test.want {
					t.Errorf("status.Code(%v) = %v; want %v", err, got, test.want)
				}
			})
		}
	})
}
