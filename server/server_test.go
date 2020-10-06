package server

import (
	"context"
	"net"
	"testing"

	pb "github.com/Saser/strecku/api/v1"
	"github.com/Saser/strecku/auth"
	"github.com/Saser/strecku/resources/stores"
	"github.com/Saser/strecku/resources/users"
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
		srv: New(users.NewRepository(), stores.NewRepository()),
		lis: bufconn.Listen(bufSize),
	}
	creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
	if err != nil {
		t.Fatal(err)
	}
	s := grpc.NewServer(grpc.Creds(creds))
	pb.RegisterStreckUServer(s, f.srv)
	go func() {
		if err := s.Serve(f.lis); err != nil {
			t.Errorf("s.Serve(f.lis) = %v", err)
		}
	}()
	t.Cleanup(s.GracefulStop)
	return f
}

func (f *fixture) client(ctx context.Context, t *testing.T, opts ...grpc.DialOption) pb.StreckUClient {
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
	return pb.NewStreckUClient(cc)
}

func (f *fixture) authClient(ctx context.Context, t *testing.T, emailAddress, password string) pb.StreckUClient {
	t.Helper()
	return f.client(ctx, t, grpc.WithPerRPCCredentials(auth.Basic{
		Username: emailAddress,
		Password: password,
	}))
}

func (f *fixture) backdoorCreateUser(ctx context.Context, t *testing.T, req *pb.CreateUserRequest) *pb.User {
	t.Helper()
	user := req.User
	user.Name = users.GenerateName()
	if err := f.srv.userRepo.CreateUser(ctx, user, req.Password); err != nil {
		t.Fatalf("f.srv.userRepo.CreateUser(%v, %v, %q) = %v; want nil", ctx, user, req.Password, err)
	}
	t.Cleanup(func() {
		if err := f.srv.userRepo.DeleteUser(ctx, user.Name); err != nil {
			t.Fatalf("f.srv.userRepo.DeleteUser(%v, %q) = %v; want nil", ctx, user.Name, err)
		}
	})
	return user
}

func (f *fixture) backdoorCreateStore(ctx context.Context, t *testing.T, req *pb.CreateStoreRequest) *pb.Store {
	t.Helper()
	store := req.Store
	store.Name = stores.GenerateName()
	if err := f.srv.storeRepo.CreateStore(ctx, store); err != nil {
		t.Fatalf("f.srv.storeRepo.CreateStore(%v, %v) = %v; want nil", ctx, store, err)
	}
	t.Cleanup(func() {
		if err := f.srv.storeRepo.DeleteStore(ctx, store.Name); err != nil {
			t.Fatalf("f.srv.storeRepo.DeleteStore(%v, %q) = %v; want nil", ctx, store.Name, err)
		}
	})
	return store
}

func TestServer_GetUser(t *testing.T) {
	ctx := context.Background()

	f := setUp(t)
	rootPassword := "root password"
	root := f.backdoorCreateUser(ctx, t, &pb.CreateUserRequest{
		User: &pb.User{
			EmailAddress: "root@example.com",
			DisplayName:  "Root",
			Superuser:    true,
		},
		Password: rootPassword,
	})
	otherRoot := f.backdoorCreateUser(ctx, t, &pb.CreateUserRequest{
		User: &pb.User{
			EmailAddress: "other-root@example.com",
			DisplayName:  "Other Root",
			Superuser:    true,
		},
		Password: "other root password",
	})
	userPassword := "user password"
	user := f.backdoorCreateUser(ctx, t, &pb.CreateUserRequest{
		User: &pb.User{
			EmailAddress: "user@example.com",
			DisplayName:  "User",
			Superuser:    false,
		},
		Password: userPassword,
	})
	otherUser := f.backdoorCreateUser(ctx, t, &pb.CreateUserRequest{
		User: &pb.User{
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
			req      *pb.GetUserRequest
			wantUser *pb.User
			wantCode codes.Code
		}{
			{name: "Me", req: &pb.GetUserRequest{Name: "users/me"}, wantUser: root},
			{name: "Root", req: &pb.GetUserRequest{Name: root.Name}, wantUser: root},
			{name: "OtherRoot", req: &pb.GetUserRequest{Name: otherRoot.Name}, wantUser: otherRoot},
			{name: "User", req: &pb.GetUserRequest{Name: user.Name}, wantUser: user},
			{name: "OtherUser", req: &pb.GetUserRequest{Name: otherUser.Name}, wantUser: otherUser},
			{name: "Nonexistent", req: &pb.GetUserRequest{Name: "nonexistent"}, wantCode: codes.NotFound},
			{name: "EmptyName", req: &pb.GetUserRequest{Name: ""}, wantCode: codes.InvalidArgument},
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
			req      *pb.GetUserRequest
			wantUser *pb.User
			wantCode codes.Code
		}{
			{name: "Me", req: &pb.GetUserRequest{Name: "users/me"}, wantUser: user},
			{name: "Root", req: &pb.GetUserRequest{Name: root.Name}, wantCode: codes.PermissionDenied},
			{name: "OtherRoot", req: &pb.GetUserRequest{Name: otherRoot.Name}, wantCode: codes.PermissionDenied},
			{name: "User", req: &pb.GetUserRequest{Name: user.Name}, wantUser: user},
			{name: "OtherUser", req: &pb.GetUserRequest{Name: otherUser.Name}, wantCode: codes.PermissionDenied},
			{name: "Nonexistent", req: &pb.GetUserRequest{Name: "nonexistent"}, wantCode: codes.PermissionDenied},
			{name: "Empty", req: &pb.GetUserRequest{Name: ""}, wantCode: codes.InvalidArgument},
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
	root := f.backdoorCreateUser(ctx, t, &pb.CreateUserRequest{
		User: &pb.User{
			EmailAddress: "root@example.com",
			DisplayName:  "Root",
			Superuser:    true,
		},
		Password: rootPassword,
	})
	userPassword := "user password"
	user := f.backdoorCreateUser(ctx, t, &pb.CreateUserRequest{
		User: &pb.User{
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
			req      *pb.ListUsersRequest
			wantResp *pb.ListUsersResponse
			wantCode codes.Code
		}{
			{
				name: "OK",
				req:  &pb.ListUsersRequest{},
				wantResp: &pb.ListUsersResponse{
					Users: []*pb.User{root, user},
				},
			},
			{name: "NegativePageSize", req: &pb.ListUsersRequest{PageSize: -1}, wantCode: codes.InvalidArgument},
			{name: "NonEmptyPageToken", req: &pb.ListUsersRequest{PageToken: "token"}, wantCode: codes.Unimplemented},
			{name: "PositivePageSize", req: &pb.ListUsersRequest{PageSize: 1}, wantCode: codes.Unimplemented},
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
			req      *pb.ListUsersRequest
			wantResp *pb.ListUsersResponse
			wantCode codes.Code
		}{
			{
				name: "OK",
				req:  &pb.ListUsersRequest{},
				wantResp: &pb.ListUsersResponse{
					Users: []*pb.User{user},
				},
			},
			{name: "NegativePageSize", req: &pb.ListUsersRequest{PageSize: -1}, wantCode: codes.InvalidArgument},
			{name: "NonEmptyPageToken", req: &pb.ListUsersRequest{PageToken: "token"}, wantCode: codes.Unimplemented},
			{name: "PositivePageSize", req: &pb.ListUsersRequest{PageSize: 1}, wantCode: codes.Unimplemented},
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
	root := f.backdoorCreateUser(ctx, t, &pb.CreateUserRequest{
		User: &pb.User{
			EmailAddress: "root@example.com",
			DisplayName:  "Root",
			Superuser:    true,
		},
		Password: rootPassword,
	})
	userPassword := "user password"
	user := f.backdoorCreateUser(ctx, t, &pb.CreateUserRequest{
		User: &pb.User{
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
			req      *pb.CreateUserRequest
			wantCode codes.Code
		}{
			{
				name: "CreateSuperuser",
				req: &pb.CreateUserRequest{
					User: &pb.User{
						EmailAddress: "other-root@example.com",
						DisplayName:  "Other Root",
						Superuser:    true,
					},
					Password: "other root password",
				},
			},
			{
				name: "CreateNormalUser",
				req: &pb.CreateUserRequest{
					User: &pb.User{
						EmailAddress: "other-user@example.com",
						DisplayName:  "Other User",
						Superuser:    false,
					},
					Password: "other user password",
				},
			},
			{
				name: "DuplicateEmailAddress",
				req: &pb.CreateUserRequest{
					User: &pb.User{
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
				req: &pb.CreateUserRequest{
					User: &pb.User{
						EmailAddress: "",
						DisplayName:  "User",
					},
					Password: userPassword,
				},
				wantCode: codes.InvalidArgument,
			},
			{
				name: "EmptyDisplayName",
				req: &pb.CreateUserRequest{
					User: &pb.User{
						EmailAddress: "other-user@example.com",
						DisplayName:  "",
					},
					Password: userPassword,
				},
				wantCode: codes.InvalidArgument,
			},
			{
				name: "EmptyPassword",
				req: &pb.CreateUserRequest{
					User: &pb.User{
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
			req  *pb.CreateUserRequest
			want codes.Code
		}{
			{
				name: "CreateSuperuser",
				req: &pb.CreateUserRequest{
					User: &pb.User{
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
				req: &pb.CreateUserRequest{
					User: &pb.User{
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
				req: &pb.CreateUserRequest{
					User: &pb.User{
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
				req: &pb.CreateUserRequest{
					User: &pb.User{
						EmailAddress: "",
						DisplayName:  "User",
					},
					Password: userPassword,
				},
				want: codes.InvalidArgument,
			},
			{
				name: "EmptyDisplayName",
				req: &pb.CreateUserRequest{
					User: &pb.User{
						EmailAddress: "other-user@example.com",
						DisplayName:  "",
					},
					Password: userPassword,
				},
				want: codes.InvalidArgument,
			},
			{
				name: "EmptyPassword",
				req: &pb.CreateUserRequest{
					User: &pb.User{
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

func TestServer_GetStore(t *testing.T) {
	ctx := context.Background()

	f := setUp(t)
	rootPassword := "root password"
	root := f.backdoorCreateUser(ctx, t, &pb.CreateUserRequest{
		User: &pb.User{
			EmailAddress: "root@example.com",
			DisplayName:  "Root",
			Superuser:    true,
		},
		Password: rootPassword,
	})
	userPassword := "user password"
	user := f.backdoorCreateUser(ctx, t, &pb.CreateUserRequest{
		User: &pb.User{
			EmailAddress: "user@example.com",
			DisplayName:  "User",
			Superuser:    false,
		},
		Password: userPassword,
	})
	store := f.backdoorCreateStore(ctx, t, &pb.CreateStoreRequest{
		Store: &pb.Store{
			DisplayName: "Store",
		},
	})

	t.Run("AsSuperuser", func(t *testing.T) {
		client := f.authClient(ctx, t, root.EmailAddress, rootPassword)
		for _, test := range []struct {
			name      string
			req       *pb.GetStoreRequest
			wantStore *pb.Store
			wantCode  codes.Code
		}{
			{
				name:      "OK",
				req:       &pb.GetStoreRequest{Name: store.Name},
				wantStore: store,
			},
			{
				name:     "EmptyName",
				req:      &pb.GetStoreRequest{Name: ""},
				wantCode: codes.InvalidArgument,
			},
			// TODO: reenable the below test after migrating to teststores.
			//{
			//	name:     "Nonexistent",
			//	req:      &pb.GetStoreRequest{Name: "nonexistent"},
			//	wantCode: codes.NotFound,
			//},
		} {
			t.Run(test.name, func(t *testing.T) {
				gotStore, err := client.GetStore(ctx, test.req)
				if gotCode := status.Code(err); gotCode != test.wantCode {
					t.Errorf("status.Code(%v) = %v; want %v", err, gotCode, test.wantCode)
				}
				if test.wantCode != codes.OK {
					return
				}
				if diff := cmp.Diff(gotStore, test.wantStore, protocmp.Transform()); diff != "" {
					t.Errorf("-got +want:\n%s", diff)
				}
			})
		}
	})

	t.Run("AsNormalUser", func(t *testing.T) {
		client := f.authClient(ctx, t, user.EmailAddress, userPassword)
		for _, test := range []struct {
			name      string
			req       *pb.GetStoreRequest
			wantStore *pb.Store
			wantCode  codes.Code
		}{
			{
				name:      "OK",
				req:       &pb.GetStoreRequest{Name: store.Name},
				wantStore: store,
			},
			{
				name:     "EmptyName",
				req:      &pb.GetStoreRequest{Name: ""},
				wantCode: codes.InvalidArgument,
			},
			// TODO: reenable the below test after migrating to teststores.
			//{
			//	name:     "Nonexistent",
			//	req:      &pb.GetStoreRequest{Name: "nonexistent"},
			//	wantCode: codes.PermissionDenied,
			//},
		} {
			t.Run(test.name, func(t *testing.T) {
				gotStore, err := client.GetStore(ctx, test.req)
				if gotCode := status.Code(err); gotCode != test.wantCode {
					t.Errorf("status.Code(%v) = %v; want %v", err, gotCode, test.wantCode)
				}
				if test.wantCode != codes.OK {
					return
				}
				if diff := cmp.Diff(gotStore, test.wantStore, protocmp.Transform()); diff != "" {
					t.Errorf("-got +want:\n%s", diff)
				}
			})
		}
	})
}

func TestServer_ListStores(t *testing.T) {
	ctx := context.Background()

	f := setUp(t)
	rootPassword := "root password"
	root := f.backdoorCreateUser(ctx, t, &pb.CreateUserRequest{
		User: &pb.User{
			EmailAddress: "root@example.com",
			DisplayName:  "Root",
			Superuser:    true,
		},
		Password: rootPassword,
	})
	userPassword := "user password"
	user := f.backdoorCreateUser(ctx, t, &pb.CreateUserRequest{
		User: &pb.User{
			EmailAddress: "user@example.com",
			DisplayName:  "User",
			Superuser:    false,
		},
		Password: userPassword,
	})
	testStores := []*pb.Store{
		{DisplayName: "Foo Store"},
		{DisplayName: "Bar Store"},
		{DisplayName: "Quux Store"},
	}
	for i, store := range testStores {
		testStores[i] = f.backdoorCreateStore(ctx, t, &pb.CreateStoreRequest{Store: store})
	}

	type testCase struct {
		name     string
		req      *pb.ListStoresRequest
		wantResp *pb.ListStoresResponse
		wantCode codes.Code
	}
	testCases := []testCase{
		{
			name: "OK",
			req:  &pb.ListStoresRequest{},
			wantResp: &pb.ListStoresResponse{
				Stores:        testStores,
				NextPageToken: "",
			},
			wantCode: codes.OK,
		},
		{
			name:     "NegativePageSize",
			req:      &pb.ListStoresRequest{PageSize: -1},
			wantResp: nil,
			wantCode: codes.InvalidArgument,
		},
		{
			name:     "PositivePageSize",
			req:      &pb.ListStoresRequest{PageSize: 1},
			wantResp: nil,
			wantCode: codes.Unimplemented,
		},
		{
			name:     "NonEmptyPageToken",
			req:      &pb.ListStoresRequest{PageToken: "invalid"},
			wantResp: nil,
			wantCode: codes.Unimplemented,
		},
	}
	testF := func(client pb.StreckUClient) func(*testing.T) {
		return func(t *testing.T) {
			for _, test := range testCases {
				t.Run(test.name, func(t *testing.T) {
					gotResp, err := client.ListStores(ctx, test.req)
					if gotCode := status.Code(err); gotCode != test.wantCode {
						t.Errorf("status.Code(%v) = %v; want %v", err, gotCode, test.wantCode)
					}
					if diff := cmp.Diff(
						gotResp, test.wantResp,
						protocmp.Transform(),
						protocmp.SortRepeatedFields(gotResp, "stores"),
					); diff != "" {
						t.Errorf("gotResp != test.wantResp (-got +want)\n%s", diff)
					}
				})
			}
		}
	}
	t.Run("AsSuperuser", testF(f.authClient(ctx, t, root.EmailAddress, rootPassword)))
	t.Run("AsNormalUser", testF(f.authClient(ctx, t, user.EmailAddress, userPassword)))
}

func TestServer_CreateStore(t *testing.T) {
	ctx := context.Background()

	f := setUp(t)
	rootPassword := "root password"
	root := f.backdoorCreateUser(ctx, t, &pb.CreateUserRequest{
		User: &pb.User{
			EmailAddress: "root@example.com",
			DisplayName:  "Root",
			Superuser:    true,
		},
		Password: rootPassword,
	})
	rootClient := f.authClient(ctx, t, root.EmailAddress, rootPassword)
	userPassword := "user password"
	user := f.backdoorCreateUser(ctx, t, &pb.CreateUserRequest{
		User: &pb.User{
			EmailAddress: "user@example.com",
			DisplayName:  "User",
			Superuser:    false,
		},
		Password: userPassword,
	})
	userClient := f.authClient(ctx, t, user.EmailAddress, userPassword)

	type testCase struct {
		name      string
		req       *pb.CreateStoreRequest
		wantStore *pb.Store
		wantCode  codes.Code
	}
	testCases := []testCase{
		{
			name: "OK",
			req: &pb.CreateStoreRequest{
				Store: &pb.Store{
					DisplayName: "Store",
				},
			},
			wantStore: &pb.Store{
				DisplayName: "Store",
			},
		},
		{
			name: "EmptyDisplayName",
			req: &pb.CreateStoreRequest{
				Store: &pb.Store{
					DisplayName: "",
				},
			},
			wantCode: codes.InvalidArgument,
		},
	}
	testF := func(client pb.StreckUClient, test testCase) func(t *testing.T) {
		return func(t *testing.T) {
			gotStore, err := client.CreateStore(ctx, test.req)
			if gotCode := status.Code(err); gotCode != test.wantCode {
				t.Errorf("status.Code(%v) = %v; want %v", err, gotCode, test.wantCode)
			}
			if test.wantCode != codes.OK {
				return
			}
			if diff := cmp.Diff(
				gotStore, test.wantStore, protocmp.Transform(),
				protocmp.IgnoreFields(gotStore, "name"),
			); diff != "" {
				t.Errorf("-got +want:\n%s", diff)
			}
		}
	}
	t.Run("AsSuperuser", func(t *testing.T) {
		for _, test := range testCases {
			t.Run(test.name, testF(rootClient, test))
		}
	})
	t.Run("AsNormalUser", func(t *testing.T) {
		for _, test := range testCases {
			t.Run(test.name, testF(userClient, test))
		}
	})
}
