package service

import (
	"context"
	"testing"

	pb "github.com/Saser/strecku/api/v1"
	"github.com/Saser/strecku/resources/stores"
	"github.com/Saser/strecku/testresources"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func storeLess(u1, u2 *pb.Store) bool {
	return u1.Name < u2.Name
}

func TestService_GetStore(t *testing.T) {
	ctx := context.Background()
	c := serveAndDial(ctx, t, seed(ctx, t))
	for _, test := range []struct {
		desc      string
		req       *pb.GetStoreRequest
		wantStore *pb.Store
		wantCode  codes.Code
	}{
		{
			desc:      "OK",
			req:       &pb.GetStoreRequest{Name: testresources.Bar.Name},
			wantStore: testresources.Bar,
			wantCode:  codes.OK,
		},
		{
			desc:      "EmptyName",
			req:       &pb.GetStoreRequest{Name: ""},
			wantStore: nil,
			wantCode:  codes.InvalidArgument,
		},
		{
			desc:      "InvalidName",
			req:       &pb.GetStoreRequest{Name: testresources.Alice.Name}, // name of a user
			wantStore: nil,
			wantCode:  codes.InvalidArgument,
		},
		{
			desc:      "NotFound",
			req:       &pb.GetStoreRequest{Name: testresources.Pharmacy.Name},
			wantStore: nil,
			wantCode:  codes.NotFound,
		},
	} {
		t.Run(test.desc, func(t *testing.T) {
			store, err := c.GetStore(ctx, test.req)
			if diff := cmp.Diff(store, test.wantStore, protocmp.Transform()); diff != "" {
				t.Errorf("c.GetStore(%v, %v) store = != test.wantStore (-got +want)\n%s", ctx, test.req, diff)
			}
			if got := status.Code(err); got != test.wantCode {
				t.Errorf("status.Code(%v) = %v; want %v", err, got, test.wantCode)
			}
		})
	}
}

func TestService_ListStores(t *testing.T) {
	ctx := context.Background()
	c := serveAndDial(ctx, t, seed(ctx, t))
	for _, test := range []struct {
		desc     string
		req      *pb.ListStoresRequest
		wantResp *pb.ListStoresResponse
		wantCode codes.Code
	}{
		{
			desc:     "OK",
			req:      &pb.ListStoresRequest{PageSize: 0, PageToken: ""},
			wantResp: &pb.ListStoresResponse{Stores: []*pb.Store{testresources.Bar, testresources.Mall}},
			wantCode: codes.OK,
		},
		{
			desc:     "NegativePageSize",
			req:      &pb.ListStoresRequest{PageSize: -1, PageToken: ""},
			wantResp: nil,
			wantCode: codes.InvalidArgument,
		},
		{
			desc:     "PaginationUnimplemented_PositivePageSize",
			req:      &pb.ListStoresRequest{PageSize: 1, PageToken: ""},
			wantResp: nil,
			wantCode: codes.Unimplemented,
		},
		{
			desc:     "PaginationUnimplemented_NonEmptyPageToken",
			req:      &pb.ListStoresRequest{PageSize: 0, PageToken: "token"},
			wantResp: nil,
			wantCode: codes.Unimplemented,
		},
	} {
		t.Run(test.desc, func(t *testing.T) {
			resp, err := c.ListStores(ctx, test.req)
			if diff := cmp.Diff(
				resp, test.wantResp, protocmp.Transform(),
				protocmp.FilterField(new(pb.ListStoresResponse), "stores", protocmp.SortRepeated(storeLess)),
			); diff != "" {
				t.Errorf("c.ListStores(%v, %v) resp != test.wantResp (-got +want)\n%s", ctx, test.req, diff)
			}
			if got := status.Code(err); got != test.wantCode {
				t.Errorf("status.Code(%v) = %v; want %v", err, got, test.wantCode)
			}
		})
	}
}

func TestService_CreateStore(t *testing.T) {
	ctx := context.Background()
	for _, test := range []struct {
		desc      string
		req       *pb.CreateStoreRequest
		wantStore *pb.Store
		wantCode  codes.Code
	}{
		{
			desc: "OK",
			req: &pb.CreateStoreRequest{
				Store: testresources.Mall,
			},
			wantStore: testresources.Mall,
			wantCode:  codes.OK,
		},
	} {
		t.Run(test.desc, func(t *testing.T) {
			c := serveAndDial(ctx, t, seed(ctx, t))
			store, err := c.CreateStore(ctx, test.req)
			if diff := cmp.Diff(
				store, test.wantStore, protocmp.Transform(),
				protocmp.IgnoreFields(new(pb.Store), "name"),
			); diff != "" {
				t.Errorf("c.CreateStore(%v, %v) store != test.wantStore (-got +want)\n%s", ctx, test.req, diff)
			}
			if got := status.Code(err); got != test.wantCode {
				t.Errorf("status.Code(%v) = %v; want %v", err, got, test.wantCode)
			}
		})
	}
}

func TestService_UpdateStore(t *testing.T) {
	ctx := context.Background()
	// Test scenario(s) where the update is successful.
	t.Run("OK", func(t *testing.T) {
		oldBar := stores.Clone(testresources.Bar)
		newBar := stores.Clone(oldBar)
		newBar.DisplayName = "New Bar"
		for _, test := range []struct {
			desc string
			req  *pb.UpdateStoreRequest
			want *pb.Store
		}{
			{
				desc: "NoOp_NilUpdateMask",
				req: &pb.UpdateStoreRequest{
					Store:      oldBar,
					UpdateMask: nil,
				},
				want: oldBar,
			},
			{
				desc: "NoOp_AllPaths",
				req: &pb.UpdateStoreRequest{
					Store:      oldBar,
					UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"display_name"}},
				},
				want: oldBar,
			},
			{
				desc: "NoOp_NoPaths",
				req: &pb.UpdateStoreRequest{
					Store:      newBar,
					UpdateMask: &fieldmaskpb.FieldMask{Paths: nil},
				},
				want: oldBar,
			},
			{
				desc: "FullUpdate_NilUpdateMask",
				req: &pb.UpdateStoreRequest{
					Store:      newBar,
					UpdateMask: nil,
				},
				want: newBar,
			},
			{
				desc: "FullUpdate_AllPaths",
				req: &pb.UpdateStoreRequest{
					Store:      newBar,
					UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"display_name"}},
				},
				want: newBar,
			},
		} {
			t.Run(test.desc, func(t *testing.T) {
				c := serveAndDial(ctx, t, seed(ctx, t))
				store, err := c.UpdateStore(ctx, test.req)
				if diff := cmp.Diff(store, test.want, protocmp.Transform()); diff != "" {
					t.Errorf("c.UpdateStore(%v, %v) store != test.want (-got +want)\n%s", ctx, test.req, diff)
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
			req  *pb.UpdateStoreRequest
			want codes.Code
		}{
			{
				desc: "EmptyDisplayName",
				req: &pb.UpdateStoreRequest{
					Store: func() *pb.Store {
						bar := stores.Clone(testresources.Bar)
						bar.DisplayName = ""
						return bar
					}(),
					UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"display_name"}},
				},
				want: codes.InvalidArgument,
			},
			{
				desc: "NotFound",
				req: &pb.UpdateStoreRequest{
					Store:      testresources.Pharmacy,
					UpdateMask: nil,
				},
				want: codes.NotFound,
			},
			{
				desc: "InvalidUpdateMask",
				req: &pb.UpdateStoreRequest{
					Store:      testresources.Bar,
					UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"invalid"}},
				},
				want: codes.InvalidArgument,
			},
		} {
			t.Run(test.desc, func(t *testing.T) {
				c := serveAndDial(ctx, t, seed(ctx, t))
				_, err := c.UpdateStore(ctx, test.req)
				if got := status.Code(err); got != test.want {
					t.Errorf("status.Code(%v) = %v; want %v", err, got, test.want)
				}
			})
		}
	})
}

func TestService_DeleteStore(t *testing.T) {
	ctx := context.Background()
	// Test scenario(s) where the delete is successful.
	t.Run("OK", func(t *testing.T) {
		c := serveAndDial(ctx, t, seed(ctx, t))
		{
			req := &pb.DeleteStoreRequest{Name: testresources.Bar.Name}
			_, err := c.DeleteStore(ctx, req)
			if got, want := status.Code(err), codes.OK; got != want {
				t.Errorf("status.Code(%v) = %v; want %v", err, got, want)
			}
		}
		{
			req := &pb.GetStoreRequest{Name: testresources.Bar.Name}
			_, err := c.GetStore(ctx, req)
			if got, want := status.Code(err), codes.NotFound; got != want {
				t.Errorf("status.Code(%v) = %v; want %v", err, got, want)
			}
		}
	})
	// Test scenario(s) where the delete fails.
	t.Run("Errors", func(t *testing.T) {
		c := serveAndDial(ctx, t, seed(ctx, t))
		for _, test := range []struct {
			desc string
			req  *pb.DeleteStoreRequest
			want codes.Code
		}{
			{
				desc: "EmptyName",
				req:  &pb.DeleteStoreRequest{Name: ""},
				want: codes.InvalidArgument,
			},
			{
				desc: "NotFound",
				req:  &pb.DeleteStoreRequest{Name: testresources.Pharmacy.Name},
				want: codes.NotFound,
			},
		} {
			t.Run(test.desc, func(t *testing.T) {
				_, err := c.DeleteStore(ctx, test.req)
				if got := status.Code(err); got != test.want {
					t.Errorf("status.Code(%v) = %v; want %v", err, got, test.want)
				}
			})
		}
	})
}
