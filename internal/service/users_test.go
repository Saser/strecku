package service

import (
	"context"
	"testing"

	pb "github.com/Saser/strecku/api/v1"
	"github.com/Saser/strecku/resources/testresources"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/testing/protocmp"
)

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
