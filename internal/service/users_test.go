package service

import (
	"context"
	"testing"

	pb "github.com/Saser/strecku/api/v1"
	"github.com/Saser/strecku/resources/testresources"
	"github.com/Saser/strecku/resources/users"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func userLess(u1, u2 *pb.User) bool {
	return u1.Name < u2.Name
}

func TestService_GetUser(t *testing.T) {
	ctx := context.Background()
	c := serveAndDial(ctx, t, seed(t))
	for _, test := range []struct {
		desc     string
		req      *pb.GetUserRequest
		wantUser *pb.User
		wantCode codes.Code
	}{
		{
			desc:     "OK",
			req:      &pb.GetUserRequest{Name: testresources.Alice.Name},
			wantUser: testresources.Alice,
			wantCode: codes.OK,
		},
		{
			desc:     "EmptyName",
			req:      &pb.GetUserRequest{Name: ""},
			wantUser: nil,
			wantCode: codes.InvalidArgument,
		},
		{
			desc:     "InvalidName",
			req:      &pb.GetUserRequest{Name: testresources.Bar.Name}, // name of a store
			wantUser: nil,
			wantCode: codes.InvalidArgument,
		},
		{
			desc:     "NotFound",
			req:      &pb.GetUserRequest{Name: testresources.Carol.Name}, // name of a store
			wantUser: nil,
			wantCode: codes.NotFound,
		},
	} {
		t.Run(test.desc, func(t *testing.T) {
			user, err := c.GetUser(ctx, test.req)
			if diff := cmp.Diff(user, test.wantUser, protocmp.Transform()); diff != "" {
				t.Errorf("c.GetUser(%v, %v) user = != test.wantUser (-got +want)\n%s", ctx, test.req, diff)
			}
			if got := status.Code(err); got != test.wantCode {
				t.Errorf("status.Code(%v) = %v; want %v", err, got, test.wantCode)
			}
		})
	}
}

func TestService_ListUsers(t *testing.T) {
	ctx := context.Background()
	c := serveAndDial(ctx, t, seed(t))
	for _, test := range []struct {
		desc     string
		req      *pb.ListUsersRequest
		wantResp *pb.ListUsersResponse
		wantCode codes.Code
	}{
		{
			desc:     "OK",
			req:      &pb.ListUsersRequest{PageSize: 0, PageToken: ""},
			wantResp: &pb.ListUsersResponse{Users: []*pb.User{testresources.Alice, testresources.Bob}},
			wantCode: codes.OK,
		},
		{
			desc:     "NegativePageSize",
			req:      &pb.ListUsersRequest{PageSize: -1, PageToken: ""},
			wantResp: nil,
			wantCode: codes.InvalidArgument,
		},
		{
			desc:     "PaginationUnimplemented_PositivePageSize",
			req:      &pb.ListUsersRequest{PageSize: 1, PageToken: ""},
			wantResp: nil,
			wantCode: codes.Unimplemented,
		},
		{
			desc:     "PaginationUnimplemented_NonEmptyPageToken",
			req:      &pb.ListUsersRequest{PageSize: 0, PageToken: "token"},
			wantResp: nil,
			wantCode: codes.Unimplemented,
		},
	} {
		t.Run(test.desc, func(t *testing.T) {
			resp, err := c.ListUsers(ctx, test.req)
			if diff := cmp.Diff(
				resp, test.wantResp, protocmp.Transform(),
				protocmp.FilterField(new(pb.ListUsersResponse), "users", protocmp.SortRepeated(userLess)),
			); diff != "" {
				t.Errorf("c.ListUsers(%v, %v) resp != test.wantResp (-got +want)\n%s", ctx, test.req, diff)
			}
			if got := status.Code(err); got != test.wantCode {
				t.Errorf("status.Code(%v) = %v; want %v", err, got, test.wantCode)
			}
		})
	}
}

func TestService_CreateUser(t *testing.T) {
	ctx := context.Background()
	for _, test := range []struct {
		desc     string
		req      *pb.CreateUserRequest
		wantUser *pb.User
		wantCode codes.Code
	}{
		{
			desc: "OK",
			req: &pb.CreateUserRequest{
				User:     testresources.Carol,
				Password: testresources.CarolPassword,
			},
			wantUser: testresources.Carol,
			wantCode: codes.OK,
		},
		{
			desc: "DuplicateEmailAddress",
			req: &pb.CreateUserRequest{
				User: func() *pb.User {
					carol := users.Clone(testresources.Carol)
					carol.EmailAddress = testresources.Alice.EmailAddress
					return carol
				}(),
				Password: testresources.CarolPassword,
			},
			wantUser: nil,
			wantCode: codes.AlreadyExists,
		},
	} {
		t.Run(test.desc, func(t *testing.T) {
			c := serveAndDial(ctx, t, seed(t))
			user, err := c.CreateUser(ctx, test.req)
			if diff := cmp.Diff(
				user, test.wantUser, protocmp.Transform(),
				protocmp.IgnoreFields(new(pb.User), "name"),
			); diff != "" {
				t.Errorf("c.CreateUser(%v, %v) user != test.wantUser (-got +want)\n%s", ctx, test.req, diff)
			}
			if got := status.Code(err); got != test.wantCode {
				t.Errorf("status.Code(%v) = %v; want %v", err, got, test.wantCode)
			}
		})
	}
}

func TestService_UpdateUser(t *testing.T) {
	ctx := context.Background()
	// Test scenario(s) where the update is successful.
	t.Run("OK", func(t *testing.T) {
		oldAlice := users.Clone(testresources.Alice)
		newAlice := users.Clone(oldAlice)
		newAlice.EmailAddress = "new-alice@example.com"
		newAlice.DisplayName = "New Alice"
		for _, test := range []struct {
			desc string
			req  *pb.UpdateUserRequest
			want *pb.User
		}{
			{
				desc: "NoOp_NilUpdateMask",
				req: &pb.UpdateUserRequest{
					User:       oldAlice,
					UpdateMask: nil,
				},
				want: oldAlice,
			},
			{
				desc: "NoOp_AllPaths",
				req: &pb.UpdateUserRequest{
					User:       oldAlice,
					UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"email_address", "display_name"}},
				},
				want: oldAlice,
			},
			{
				desc: "NoOp_NoPaths",
				req: &pb.UpdateUserRequest{
					User:       newAlice,
					UpdateMask: &fieldmaskpb.FieldMask{Paths: nil},
				},
				want: oldAlice,
			},
			{
				desc: "FullUpdate_NilUpdateMask",
				req: &pb.UpdateUserRequest{
					User:       newAlice,
					UpdateMask: nil,
				},
				want: newAlice,
			},
			{
				desc: "FullUpdate_AllPaths",
				req: &pb.UpdateUserRequest{
					User:       newAlice,
					UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"email_address", "display_name"}},
				},
				want: newAlice,
			},
			{
				desc: "PartialUpdate_FullUser_EmailAddress",
				req: &pb.UpdateUserRequest{
					User:       newAlice,
					UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"email_address"}},
				},
				want: func() *pb.User {
					newAlicePartial := users.Clone(oldAlice)
					newAlicePartial.EmailAddress = newAlice.EmailAddress
					return newAlicePartial
				}(),
			},
			{
				desc: "PartialUpdate_FullUser_EmailAddress",
				req: &pb.UpdateUserRequest{
					User:       newAlice,
					UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"display_name"}},
				},
				want: func() *pb.User {
					newAlicePartial := users.Clone(oldAlice)
					newAlicePartial.DisplayName = newAlice.DisplayName
					return newAlicePartial
				}(),
			},
			{
				desc: "PartialUpdate_PartialUser_EmailAddress",
				req: &pb.UpdateUserRequest{
					User: &pb.User{
						Name:         newAlice.Name,
						EmailAddress: newAlice.EmailAddress,
					},
					UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"email_address"}},
				},
				want: func() *pb.User {
					newAlicePartial := users.Clone(oldAlice)
					newAlicePartial.EmailAddress = newAlice.EmailAddress
					return newAlicePartial
				}(),
			},
			{
				desc: "PartialUpdate_PartialUser_DisplayName",
				req: &pb.UpdateUserRequest{
					User: &pb.User{
						Name:        newAlice.Name,
						DisplayName: newAlice.DisplayName,
					},
					UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"display_name"}},
				},
				want: func() *pb.User {
					newAlicePartial := users.Clone(oldAlice)
					newAlicePartial.DisplayName = newAlice.DisplayName
					return newAlicePartial
				}(),
			},
		} {
			t.Run(test.desc, func(t *testing.T) {
				c := serveAndDial(ctx, t, seed(t))
				user, err := c.UpdateUser(ctx, test.req)
				if diff := cmp.Diff(user, test.want, protocmp.Transform()); diff != "" {
					t.Errorf("c.UpdateUser(%v, %v) user != test.want (-got +want)\n%s", ctx, test.req, diff)
				}
				if got, want := status.Code(err), codes.OK; got != want {
					t.Errorf("status.Code(%v) = %v; want %v", err, got, want)
				}
			})
		}
	})
	// Test scenario(s) where the update fails.
	t.Run("Errors", func(t *testing.T) {
		for _, test := range []struct {
			desc string
			req  *pb.UpdateUserRequest
			want codes.Code
		}{
			{
				desc: "DuplicateEmailAddress",
				req: &pb.UpdateUserRequest{
					User: func() *pb.User {
						alice := users.Clone(testresources.Alice)
						alice.EmailAddress = testresources.Bob.EmailAddress
						return alice
					}(),
					UpdateMask: nil,
				},
				want: codes.AlreadyExists,
			},
			{
				desc: "EmptyDisplayName",
				req: &pb.UpdateUserRequest{
					User: func() *pb.User {
						alice := users.Clone(testresources.Alice)
						alice.DisplayName = ""
						return alice
					}(),
					UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"display_name"}},
				},
				want: codes.InvalidArgument,
			},
			{
				desc: "NotFound",
				req: &pb.UpdateUserRequest{
					User:       testresources.Carol,
					UpdateMask: nil,
				},
				want: codes.NotFound,
			},
			{
				desc: "InvalidUpdateMask",
				req: &pb.UpdateUserRequest{
					User:       testresources.Alice,
					UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"invalid"}},
				},
				want: codes.InvalidArgument,
			},
		} {
			t.Run(test.desc, func(t *testing.T) {
				c := serveAndDial(ctx, t, seed(t))
				_, err := c.UpdateUser(ctx, test.req)
				if got := status.Code(err); got != test.want {
					t.Errorf("status.Code(%v) = %v; want %v", err, got, test.want)
				}
			})
		}
	})
}

func TestService_DeleteUser(t *testing.T) {
	ctx := context.Background()
	// Test scenario(s) where the delete is successful.
	t.Run("OK", func(t *testing.T) {
		c := serveAndDial(ctx, t, seed(t))
		{
			req := &pb.DeleteUserRequest{Name: testresources.Alice.Name}
			_, err := c.DeleteUser(ctx, req)
			if got, want := status.Code(err), codes.OK; got != want {
				t.Errorf("status.Code(%v) = %v; want %v", err, got, want)
			}
		}
		{
			req := &pb.GetUserRequest{Name: testresources.Alice.Name}
			_, err := c.GetUser(ctx, req)
			if got, want := status.Code(err), codes.NotFound; got != want {
				t.Errorf("status.Code(%v) = %v; want %v", err, got, want)
			}
		}
	})
	// Test scenario(s) where the delete fails.
	t.Run("Errors", func(t *testing.T) {
		c := serveAndDial(ctx, t, seed(t))
		for _, test := range []struct {
			desc string
			req  *pb.DeleteUserRequest
			want codes.Code
		}{
			{
				desc: "EmptyName",
				req:  &pb.DeleteUserRequest{Name: ""},
				want: codes.InvalidArgument,
			},
			{
				desc: "NotFound",
				req:  &pb.DeleteUserRequest{Name: testresources.Carol.Name},
				want: codes.NotFound,
			},
		} {
			t.Run(test.desc, func(t *testing.T) {
				_, err := c.DeleteUser(ctx, test.req)
				if got := status.Code(err); got != test.want {
					t.Errorf("status.Code(%v) = %v; want %v", err, got, test.want)
				}
			})
		}
	})
}
