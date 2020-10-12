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
