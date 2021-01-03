package service

import (
	"context"
	"testing"

	pb "github.com/Saser/strecku/api/v1"
	"github.com/Saser/strecku/resources/stores/purchases"
	"github.com/Saser/strecku/testresources"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func purchaseLess(u1, u2 *pb.Purchase) bool {
	return u1.Name < u2.Name
}

func TestService_GetPurchase(t *testing.T) {
	ctx := context.Background()
	c := serveAndDial(ctx, t, seed(ctx, t))
	for _, test := range []struct {
		desc         string
		req          *pb.GetPurchaseRequest
		wantPurchase *pb.Purchase
		wantCode     codes.Code
	}{
		{
			desc:         "OK",
			req:          &pb.GetPurchaseRequest{Name: testresources.Bar_Alice_Beer1.Name},
			wantPurchase: testresources.Bar_Alice_Beer1,
			wantCode:     codes.OK,
		},
		{
			desc:         "EmptyName",
			req:          &pb.GetPurchaseRequest{Name: ""},
			wantPurchase: nil,
			wantCode:     codes.InvalidArgument,
		},
		{
			desc:         "InvalidName",
			req:          &pb.GetPurchaseRequest{Name: testresources.Alice.Name}, // name of a user
			wantPurchase: nil,
			wantCode:     codes.InvalidArgument,
		},
		{
			desc:         "NotFound",
			req:          &pb.GetPurchaseRequest{Name: testresources.Bar_Alice_Cocktail1.Name},
			wantPurchase: nil,
			wantCode:     codes.NotFound,
		},
	} {
		t.Run(test.desc, func(t *testing.T) {
			purchase, err := c.GetPurchase(ctx, test.req)
			if diff := cmp.Diff(purchase, test.wantPurchase, protocmp.Transform()); diff != "" {
				t.Errorf("c.GetPurchase(%v, %v) purchase = != test.wantPurchase (-got +want)\n%s", ctx, test.req, diff)
			}
			if got := status.Code(err); got != test.wantCode {
				t.Errorf("status.Code(%v) = %v; want %v", err, got, test.wantCode)
			}
		})
	}
}

func TestService_ListPurchases(t *testing.T) {
	ctx := context.Background()
	c := serveAndDial(ctx, t, seed(ctx, t))
	for _, test := range []struct {
		desc     string
		req      *pb.ListPurchasesRequest
		wantResp *pb.ListPurchasesResponse
		wantCode codes.Code
	}{
		{
			desc: "OK_Bar",
			req: &pb.ListPurchasesRequest{
				Parent:    testresources.Bar.Name,
				PageSize:  0,
				PageToken: "",
			},
			wantResp: &pb.ListPurchasesResponse{
				Purchases: []*pb.Purchase{
					testresources.Bar_Alice_Beer1,
				},
			},
			wantCode: codes.OK,
		},
		{
			desc: "OK_Mall",
			req: &pb.ListPurchasesRequest{
				Parent:    testresources.Mall.Name,
				PageSize:  0,
				PageToken: "",
			},
			wantResp: &pb.ListPurchasesResponse{
				Purchases: []*pb.Purchase{
					testresources.Mall_Alice_Jeans1,
				},
			},
			wantCode: codes.OK,
		},
		{
			desc: "InvalidParent",
			req: &pb.ListPurchasesRequest{
				Parent:    testresources.Alice.Name, // name of a user
				PageSize:  0,
				PageToken: "",
			},
			wantResp: nil,
			wantCode: codes.InvalidArgument,
		},
		{
			desc: "NegativePageSize",
			req: &pb.ListPurchasesRequest{
				Parent:    testresources.Bar.Name,
				PageSize:  -1,
				PageToken: "",
			},
			wantResp: nil,
			wantCode: codes.InvalidArgument,
		},
		{
			desc: "PaginationUnimplemented_PositivePageSize",
			req: &pb.ListPurchasesRequest{
				Parent:    testresources.Bar.Name,
				PageSize:  1,
				PageToken: "",
			},
			wantResp: nil,
			wantCode: codes.Unimplemented,
		},
		{
			desc: "PaginationUnimplemented_NonEmptyPageToken",
			req: &pb.ListPurchasesRequest{
				Parent:    testresources.Bar.Name,
				PageSize:  0,
				PageToken: "token",
			},
			wantResp: nil,
			wantCode: codes.Unimplemented,
		},
	} {
		t.Run(test.desc, func(t *testing.T) {
			resp, err := c.ListPurchases(ctx, test.req)
			if diff := cmp.Diff(
				resp, test.wantResp, protocmp.Transform(),
				protocmp.FilterField(new(pb.ListPurchasesResponse), "purchases", protocmp.SortRepeated(purchaseLess)),
			); diff != "" {
				t.Errorf("c.ListPurchases(%v, %v) resp != test.wantResp (-got +want)\n%s", ctx, test.req, diff)
			}
			if got := status.Code(err); got != test.wantCode {
				t.Errorf("status.Code(%v) = %v; want %v", err, got, test.wantCode)
			}
		})
	}
}

func TestService_CreatePurchase(t *testing.T) {
	ctx := context.Background()
	for _, test := range []struct {
		desc         string
		req          *pb.CreatePurchaseRequest
		wantPurchase *pb.Purchase
		wantCode     codes.Code
	}{
		{
			desc: "OK",
			req: &pb.CreatePurchaseRequest{
				Parent:   testresources.Bar.Name,
				Purchase: testresources.Bar_Alice_Cocktail1,
			},
			wantPurchase: testresources.Bar_Alice_Cocktail1,
			wantCode:     codes.OK,
		},
	} {
		t.Run(test.desc, func(t *testing.T) {
			c := serveAndDial(ctx, t, seed(ctx, t))
			purchase, err := c.CreatePurchase(ctx, test.req)
			if diff := cmp.Diff(
				purchase, test.wantPurchase, protocmp.Transform(),
				protocmp.IgnoreFields(new(pb.Purchase), "name"),
			); diff != "" {
				t.Errorf("c.CreatePurchase(%v, %v) purchase != test.wantPurchase (-got +want)\n%s", ctx, test.req, diff)
			}
			if got := status.Code(err); got != test.wantCode {
				t.Errorf("status.Code(%v) = %v; want %v", err, got, test.wantCode)
			}
		})
	}
}

func TestService_UpdatePurchase(t *testing.T) {
	ctx := context.Background()
	// Test scenario(s) where the update is successful.
	t.Run("OK", func(t *testing.T) {
		oldPurchase := purchases.Clone(testresources.Bar_Alice_Beer1)
		newPurchase := purchases.Clone(oldPurchase)
		newPurchase.Lines[0] = &pb.Purchase_Line{
			Description: testresources.Cocktail.DisplayName,
			Quantity:    2,
			PriceCents:  testresources.Cocktail.FullPriceCents,
			Product:     testresources.Cocktail.Name,
		}
		newPurchase.Lines = append(newPurchase.Lines, &pb.Purchase_Line{
			Description: testresources.Beer.DisplayName,
			Quantity:    10,
			PriceCents:  testresources.Beer.FullPriceCents,
			Product:     testresources.Beer.Name,
		})
		for _, test := range []struct {
			desc string
			req  *pb.UpdatePurchaseRequest
			want *pb.Purchase
		}{
			{
				desc: "NoOp_NilUpdateMask",
				req: &pb.UpdatePurchaseRequest{
					Purchase:   oldPurchase,
					UpdateMask: nil,
				},
				want: oldPurchase,
			},
			{
				desc: "NoOp_AllPaths",
				req: &pb.UpdatePurchaseRequest{
					Purchase:   oldPurchase,
					UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"lines"}},
				},
				want: oldPurchase,
			},
			{
				desc: "NoOp_NoPaths",
				req: &pb.UpdatePurchaseRequest{
					Purchase:   newPurchase,
					UpdateMask: &fieldmaskpb.FieldMask{Paths: nil},
				},
				want: oldPurchase,
			},
			{
				desc: "FullUpdate_NilUpdateMask",
				req: &pb.UpdatePurchaseRequest{
					Purchase:   newPurchase,
					UpdateMask: nil,
				},
				want: newPurchase,
			},
			{
				desc: "FullUpdate_AllPaths",
				req: &pb.UpdatePurchaseRequest{
					Purchase:   newPurchase,
					UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"lines"}},
				},
				want: newPurchase,
			},
		} {
			t.Run(test.desc, func(t *testing.T) {
				c := serveAndDial(ctx, t, seed(ctx, t))
				purchase, err := c.UpdatePurchase(ctx, test.req)
				if diff := cmp.Diff(purchase, test.want, protocmp.Transform()); diff != "" {
					t.Errorf("c.UpdatePurchase(%v, %v) purchase != test.want (-got +want)\n%s", ctx, test.req, diff)
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
			req  *pb.UpdatePurchaseRequest
			want codes.Code
		}{
			{
				desc: "UpdateUser",
				req: &pb.UpdatePurchaseRequest{
					Purchase: func() *pb.Purchase {
						purchase := purchases.Clone(testresources.Bar_Alice_Beer1)
						purchase.User = testresources.Bob.Name
						return purchase
					}(),
					UpdateMask: nil,
				},
				want: codes.InvalidArgument,
			},
			{
				desc: "Line_PriceCentsPositive",
				req: &pb.UpdatePurchaseRequest{
					Purchase: func() *pb.Purchase {
						purchase := purchases.Clone(testresources.Bar_Alice_Beer1)
						purchase.Lines[0].PriceCents = 1000
						return purchase
					}(),
					UpdateMask: nil,
				},
				want: codes.InvalidArgument,
			},
			{
				desc: "Line_QuantityZero",
				req: &pb.UpdatePurchaseRequest{
					Purchase: func() *pb.Purchase {
						purchase := purchases.Clone(testresources.Bar_Alice_Beer1)
						purchase.Lines[0].Quantity = 0
						return purchase
					}(),
					UpdateMask: nil,
				},
				want: codes.InvalidArgument,
			},
			{
				desc: "Line_QuantityNegative",
				req: &pb.UpdatePurchaseRequest{
					Purchase: func() *pb.Purchase {
						purchase := purchases.Clone(testresources.Bar_Alice_Beer1)
						purchase.Lines[0].Quantity = -1
						return purchase
					}(),
					UpdateMask: nil,
				},
				want: codes.InvalidArgument,
			},
			{
				desc: "Line_DescriptionEmpty",
				req: &pb.UpdatePurchaseRequest{
					Purchase: func() *pb.Purchase {
						purchase := purchases.Clone(testresources.Bar_Alice_Beer1)
						purchase.Lines[0].Description = ""
						return purchase
					}(),
					UpdateMask: nil,
				},
				want: codes.InvalidArgument,
			},
			{
				desc: "Line_ProductFromOtherStore",
				req: &pb.UpdatePurchaseRequest{
					Purchase: func() *pb.Purchase {
						purchase := purchases.Clone(testresources.Bar_Alice_Beer1)
						purchase.Lines[0].Product = testresources.Jeans.Name
						return purchase
					}(),
					UpdateMask: nil,
				},
				want: codes.InvalidArgument,
			},
			{
				desc: "NotFound",
				req: &pb.UpdatePurchaseRequest{
					Purchase:   testresources.Bar_Alice_Cocktail1,
					UpdateMask: nil,
				},
				want: codes.NotFound,
			},
			{
				desc: "InvalidUpdateMask",
				req: &pb.UpdatePurchaseRequest{
					Purchase:   testresources.Bar_Alice_Beer1,
					UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"invalid"}},
				},
				want: codes.InvalidArgument,
			},
		} {
			t.Run(test.desc, func(t *testing.T) {
				c := serveAndDial(ctx, t, seed(ctx, t))
				_, err := c.UpdatePurchase(ctx, test.req)
				if got := status.Code(err); got != test.want {
					t.Errorf("status.Code(%v) = %v; want %v", err, got, test.want)
				}
			})
		}
	})
}

func TestService_DeletePurchase(t *testing.T) {
	ctx := context.Background()
	// Test scenario(s) where the delete is successful.
	t.Run("OK", func(t *testing.T) {
		c := serveAndDial(ctx, t, seed(ctx, t))
		{
			req := &pb.DeletePurchaseRequest{Name: testresources.Bar_Alice_Beer1.Name}
			_, err := c.DeletePurchase(ctx, req)
			if got, want := status.Code(err), codes.OK; got != want {
				t.Errorf("status.Code(%v) = %v; want %v", err, got, want)
			}
		}
		{
			req := &pb.GetPurchaseRequest{Name: testresources.Bar_Alice_Beer1.Name}
			_, err := c.GetPurchase(ctx, req)
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
			req  *pb.DeletePurchaseRequest
			want codes.Code
		}{
			{
				desc: "EmptyName",
				req:  &pb.DeletePurchaseRequest{Name: ""},
				want: codes.InvalidArgument,
			},
			{
				desc: "NotFound",
				req:  &pb.DeletePurchaseRequest{Name: testresources.Bar_Alice_Cocktail1.Name},
				want: codes.NotFound,
			},
		} {
			t.Run(test.desc, func(t *testing.T) {
				_, err := c.DeletePurchase(ctx, test.req)
				if got := status.Code(err); got != test.want {
					t.Errorf("status.Code(%v) = %v; want %v", err, got, test.want)
				}
			})
		}
	})
}
