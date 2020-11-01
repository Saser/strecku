package service

import (
	"context"
	"testing"

	pb "github.com/Saser/strecku/api/v1"
	"github.com/Saser/strecku/resources/stores/memberships"
	"github.com/Saser/strecku/resources/testresources"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func membershipLess(u1, u2 *pb.Membership) bool {
	return u1.Name < u2.Name
}

func TestService_GetMembership(t *testing.T) {
	ctx := context.Background()
	c := serveAndDial(ctx, t, seed(t))
	for _, test := range []struct {
		desc           string
		req            *pb.GetMembershipRequest
		wantMembership *pb.Membership
		wantCode       codes.Code
	}{
		{
			desc:           "OK",
			req:            &pb.GetMembershipRequest{Name: testresources.Bar_Alice.Name},
			wantMembership: testresources.Bar_Alice,
			wantCode:       codes.OK,
		},
		{
			desc:           "EmptyName",
			req:            &pb.GetMembershipRequest{Name: ""},
			wantMembership: nil,
			wantCode:       codes.InvalidArgument,
		},
		{
			desc:           "InvalidName",
			req:            &pb.GetMembershipRequest{Name: testresources.Alice.Name}, // name of a user
			wantMembership: nil,
			wantCode:       codes.InvalidArgument,
		},
		{
			desc:           "NotFound",
			req:            &pb.GetMembershipRequest{Name: testresources.Mall_Bob.Name},
			wantMembership: nil,
			wantCode:       codes.NotFound,
		},
	} {
		t.Run(test.desc, func(t *testing.T) {
			membership, err := c.GetMembership(ctx, test.req)
			if diff := cmp.Diff(membership, test.wantMembership, protocmp.Transform()); diff != "" {
				t.Errorf("c.GetMembership(%v, %v) membership = != test.wantMembership (-got +want)\n%s", ctx, test.req, diff)
			}
			if got := status.Code(err); got != test.wantCode {
				t.Errorf("status.Code(%v) = %v; want %v", err, got, test.wantCode)
			}
		})
	}
}

func TestService_ListMemberships(t *testing.T) {
	ctx := context.Background()
	c := serveAndDial(ctx, t, seed(t))
	for _, test := range []struct {
		desc     string
		req      *pb.ListMembershipsRequest
		wantResp *pb.ListMembershipsResponse
		wantCode codes.Code
	}{
		{
			desc: "OK_Bar",
			req: &pb.ListMembershipsRequest{
				Parent:    testresources.Bar.Name,
				PageSize:  0,
				PageToken: "",
			},
			wantResp: &pb.ListMembershipsResponse{
				Memberships: []*pb.Membership{
					testresources.Bar_Alice,
					testresources.Bar_Bob,
				},
			},
			wantCode: codes.OK,
		},
		{
			desc: "OK_Mall",
			req: &pb.ListMembershipsRequest{
				Parent:    testresources.Mall.Name,
				PageSize:  0,
				PageToken: "",
			},
			wantResp: &pb.ListMembershipsResponse{
				Memberships: []*pb.Membership{
					testresources.Mall_Alice,
				},
			},
			wantCode: codes.OK,
		},
		{
			desc: "InvalidParent",
			req: &pb.ListMembershipsRequest{
				Parent:    testresources.Alice.Name, // name of a user
				PageSize:  0,
				PageToken: "",
			},
			wantResp: nil,
			wantCode: codes.InvalidArgument,
		},
		{
			desc: "NegativePageSize",
			req: &pb.ListMembershipsRequest{
				Parent:    testresources.Bar.Name,
				PageSize:  -1,
				PageToken: "",
			},
			wantResp: nil,
			wantCode: codes.InvalidArgument,
		},
		{
			desc: "PaginationUnimplemented_PositivePageSize",
			req: &pb.ListMembershipsRequest{
				Parent:    testresources.Bar.Name,
				PageSize:  1,
				PageToken: "",
			},
			wantResp: nil,
			wantCode: codes.Unimplemented,
		},
		{
			desc: "PaginationUnimplemented_NonEmptyPageToken",
			req: &pb.ListMembershipsRequest{
				Parent:    testresources.Bar.Name,
				PageSize:  0,
				PageToken: "token",
			},
			wantResp: nil,
			wantCode: codes.Unimplemented,
		},
	} {
		t.Run(test.desc, func(t *testing.T) {
			resp, err := c.ListMemberships(ctx, test.req)
			if diff := cmp.Diff(
				resp, test.wantResp, protocmp.Transform(),
				protocmp.FilterField(new(pb.ListMembershipsResponse), "memberships", protocmp.SortRepeated(membershipLess)),
			); diff != "" {
				t.Errorf("c.ListMemberships(%v, %v) resp != test.wantResp (-got +want)\n%s", ctx, test.req, diff)
			}
			if got := status.Code(err); got != test.wantCode {
				t.Errorf("status.Code(%v) = %v; want %v", err, got, test.wantCode)
			}
		})
	}
}

func TestService_CreateMembership(t *testing.T) {
	ctx := context.Background()
	for _, test := range []struct {
		desc           string
		req            *pb.CreateMembershipRequest
		wantMembership *pb.Membership
		wantCode       codes.Code
	}{
		{
			desc: "OK",
			req: &pb.CreateMembershipRequest{
				Parent:     testresources.Mall.Name,
				Membership: testresources.Mall_Bob,
			},
			wantMembership: testresources.Mall_Bob,
			wantCode:       codes.OK,
		},
		{
			desc: "AlreadyExists_Name",
			req: &pb.CreateMembershipRequest{
				Parent: testresources.Bar.Name,
				Membership: func() *pb.Membership {
					membership := memberships.Clone(testresources.Mall_Bob) // Mall_Bob does not exist ...
					membership.Name = testresources.Bar_Bob.Name            // ... but Bar_Bob does.
					return membership
				}(),
			},
			wantMembership: nil,
			wantCode:       codes.AlreadyExists,
		},
		{
			desc: "AlreadyExists_StoreAndUser",
			req: &pb.CreateMembershipRequest{
				Parent: testresources.Bar.Name,
				Membership: func() *pb.Membership {
					membership := memberships.Clone(testresources.Bar_Bob)                                         // Bar_Bob already exists ...
					membership.Name = testresources.Bar.Name + "/memberships/26ce041f-3d49-43f6-9420-5fda0b340008" // ... but not with this name.
					return membership
				}(),
			},
			wantMembership: nil,
			wantCode:       codes.AlreadyExists,
		},
	} {
		t.Run(test.desc, func(t *testing.T) {
			c := serveAndDial(ctx, t, seed(t))
			membership, err := c.CreateMembership(ctx, test.req)
			if diff := cmp.Diff(
				membership, test.wantMembership, protocmp.Transform(),
				protocmp.IgnoreFields(new(pb.Membership), "name"),
			); diff != "" {
				t.Errorf("c.CreateMembership(%v, %v) membership != test.wantMembership (-got +want)\n%s", ctx, test.req, diff)
			}
			if got := status.Code(err); got != test.wantCode {
				t.Errorf("status.Code(%v) = %v; want %v", err, got, test.wantCode)
			}
		})
	}
}

func TestService_UpdateMembership(t *testing.T) {
	ctx := context.Background()
	// Test scenario(s) where the update is successful.
	t.Run("OK", func(t *testing.T) {
		oldBarAlice := memberships.Clone(testresources.Bar_Alice)
		newBarAlice := memberships.Clone(oldBarAlice)
		newBarAlice.Administrator = true
		newBarAlice.Discount = true
		for _, test := range []struct {
			desc string
			req  *pb.UpdateMembershipRequest
			want *pb.Membership
		}{
			{
				desc: "NoOp_NilUpdateMask",
				req: &pb.UpdateMembershipRequest{
					Membership: oldBarAlice,
					UpdateMask: nil,
				},
				want: oldBarAlice,
			},
			{
				desc: "NoOp_AllPaths",
				req: &pb.UpdateMembershipRequest{
					Membership: oldBarAlice,
					UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"administrator", "discount"}},
				},
				want: oldBarAlice,
			},
			{
				desc: "NoOp_NoPaths",
				req: &pb.UpdateMembershipRequest{
					Membership: newBarAlice,
					UpdateMask: &fieldmaskpb.FieldMask{Paths: nil},
				},
				want: oldBarAlice,
			},
			{
				desc: "FullUpdate_NilUpdateMask",
				req: &pb.UpdateMembershipRequest{
					Membership: newBarAlice,
					UpdateMask: nil,
				},
				want: newBarAlice,
			},
			{
				desc: "FullUpdate_AllPaths",
				req: &pb.UpdateMembershipRequest{
					Membership: newBarAlice,
					UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"administrator", "discount"}},
				},
				want: newBarAlice,
			},
			{
				desc: "PartialUpdate_FullMembership_Administrator",
				req: &pb.UpdateMembershipRequest{
					Membership: newBarAlice,
					UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"administrator"}},
				},
				want: func() *pb.Membership {
					membership := memberships.Clone(oldBarAlice)
					membership.Administrator = true
					return membership
				}(),
			},
			{
				desc: "PartialUpdate_FullMembership_Discount",
				req: &pb.UpdateMembershipRequest{
					Membership: newBarAlice,
					UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"discount"}},
				},
				want: func() *pb.Membership {
					membership := memberships.Clone(oldBarAlice)
					membership.Discount = true
					return membership
				}(),
			},
			{
				desc: "PartialUpdate_PartialMembership_Administrator",
				req: &pb.UpdateMembershipRequest{
					Membership: &pb.Membership{
						Name:          testresources.Bar_Alice.Name,
						Administrator: true,
					},
					UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"administrator"}},
				},
				want: func() *pb.Membership {
					membership := memberships.Clone(oldBarAlice)
					membership.Administrator = true
					return membership
				}(),
			},
			{
				desc: "PartialUpdate_PartialMembership_Discount",
				req: &pb.UpdateMembershipRequest{
					Membership: &pb.Membership{
						Name:     testresources.Bar_Alice.Name,
						Discount: true,
					},
					UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"discount"}},
				},
				want: func() *pb.Membership {
					membership := memberships.Clone(oldBarAlice)
					membership.Discount = true
					return membership
				}(),
			},
		} {
			t.Run(test.desc, func(t *testing.T) {
				c := serveAndDial(ctx, t, seed(t))
				membership, err := c.UpdateMembership(ctx, test.req)
				if diff := cmp.Diff(membership, test.want, protocmp.Transform()); diff != "" {
					t.Errorf("c.UpdateMembership(%v, %v) membership != test.want (-got +want)\n%s", ctx, test.req, diff)
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
			req  *pb.UpdateMembershipRequest
			want codes.Code
		}{
			{
				desc: "UpdateUser_EmptyUser",
				req: &pb.UpdateMembershipRequest{
					Membership: func() *pb.Membership {
						bar := memberships.Clone(testresources.Bar_Alice)
						bar.User = ""
						return bar
					}(),
					UpdateMask: nil,
				},
				want: codes.InvalidArgument,
			},
			{
				desc: "UpdateUser_OtherUser",
				req: &pb.UpdateMembershipRequest{
					Membership: func() *pb.Membership {
						bar := memberships.Clone(testresources.Bar_Alice)
						bar.User = testresources.Bob.Name
						return bar
					}(),
					UpdateMask: nil,
				},
				want: codes.InvalidArgument,
			},
			{
				desc: "NotFound",
				req: &pb.UpdateMembershipRequest{
					Membership: testresources.Mall_Bob,
					UpdateMask: nil,
				},
				want: codes.NotFound,
			},
			{
				desc: "InvalidUpdateMask",
				req: &pb.UpdateMembershipRequest{
					Membership: testresources.Bar_Alice,
					UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"invalid"}},
				},
				want: codes.InvalidArgument,
			},
			{
				desc: "InvalidUpdateMask_User",
				req: &pb.UpdateMembershipRequest{
					Membership: testresources.Bar_Alice,
					UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"user"}},
				},
				want: codes.InvalidArgument,
			},
		} {
			t.Run(test.desc, func(t *testing.T) {
				c := serveAndDial(ctx, t, seed(t))
				_, err := c.UpdateMembership(ctx, test.req)
				if got := status.Code(err); got != test.want {
					t.Errorf("status.Code(%v) = %v; want %v", err, got, test.want)
				}
			})
		}
	})
}

func TestService_DeleteMembership(t *testing.T) {
	ctx := context.Background()
	// Test scenario(s) where the delete is successful.
	t.Run("OK", func(t *testing.T) {
		c := serveAndDial(ctx, t, seed(t))
		{
			req := &pb.DeleteMembershipRequest{Name: testresources.Bar_Alice.Name}
			_, err := c.DeleteMembership(ctx, req)
			if got, want := status.Code(err), codes.OK; got != want {
				t.Errorf("status.Code(%v) = %v; want %v", err, got, want)
			}
		}
		{
			req := &pb.GetMembershipRequest{Name: testresources.Bar_Alice.Name}
			_, err := c.GetMembership(ctx, req)
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
			req  *pb.DeleteMembershipRequest
			want codes.Code
		}{
			{
				desc: "EmptyName",
				req:  &pb.DeleteMembershipRequest{Name: ""},
				want: codes.InvalidArgument,
			},
			{
				desc: "NotFound",
				req:  &pb.DeleteMembershipRequest{Name: testresources.Mall_Bob.Name},
				want: codes.NotFound,
			},
		} {
			t.Run(test.desc, func(t *testing.T) {
				_, err := c.DeleteMembership(ctx, test.req)
				if got := status.Code(err); got != test.want {
					t.Errorf("status.Code(%v) = %v; want %v", err, got, test.want)
				}
			})
		}
	})
}
