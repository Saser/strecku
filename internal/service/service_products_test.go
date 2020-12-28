package service

import (
	"context"
	"testing"

	pb "github.com/Saser/strecku/api/v1"
	"github.com/Saser/strecku/resources/stores/products"
	"github.com/Saser/strecku/resources/testresources"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func productLess(u1, u2 *pb.Product) bool {
	return u1.Name < u2.Name
}

func TestService_GetProduct(t *testing.T) {
	ctx := context.Background()
	c := serveAndDial(ctx, t, seed(ctx, t))
	for _, test := range []struct {
		desc        string
		req         *pb.GetProductRequest
		wantProduct *pb.Product
		wantCode    codes.Code
	}{
		{
			desc:        "OK",
			req:         &pb.GetProductRequest{Name: testresources.Beer.Name},
			wantProduct: testresources.Beer,
			wantCode:    codes.OK,
		},
		{
			desc:        "EmptyName",
			req:         &pb.GetProductRequest{Name: ""},
			wantProduct: nil,
			wantCode:    codes.InvalidArgument,
		},
		{
			desc:        "InvalidName",
			req:         &pb.GetProductRequest{Name: testresources.Alice.Name}, // name of a user
			wantProduct: nil,
			wantCode:    codes.InvalidArgument,
		},
		{
			desc:        "NotFound",
			req:         &pb.GetProductRequest{Name: testresources.Pills.Name},
			wantProduct: nil,
			wantCode:    codes.NotFound,
		},
	} {
		t.Run(test.desc, func(t *testing.T) {
			product, err := c.GetProduct(ctx, test.req)
			if diff := cmp.Diff(product, test.wantProduct, protocmp.Transform()); diff != "" {
				t.Errorf("c.GetProduct(%v, %v) product = != test.wantProduct (-got +want)\n%s", ctx, test.req, diff)
			}
			if got := status.Code(err); got != test.wantCode {
				t.Errorf("status.Code(%v) = %v; want %v", err, got, test.wantCode)
			}
		})
	}
}

func TestService_ListProducts(t *testing.T) {
	ctx := context.Background()
	c := serveAndDial(ctx, t, seed(ctx, t))
	for _, test := range []struct {
		desc     string
		req      *pb.ListProductsRequest
		wantResp *pb.ListProductsResponse
		wantCode codes.Code
	}{
		{
			desc: "OK_Bar",
			req: &pb.ListProductsRequest{
				Parent:    testresources.Bar.Name,
				PageSize:  0,
				PageToken: "",
			},
			wantResp: &pb.ListProductsResponse{
				Products: []*pb.Product{
					testresources.Beer,
				},
			},
			wantCode: codes.OK,
		},
		{
			desc: "OK_Mall",
			req: &pb.ListProductsRequest{
				Parent:    testresources.Mall.Name,
				PageSize:  0,
				PageToken: "",
			},
			wantResp: &pb.ListProductsResponse{
				Products: []*pb.Product{
					testresources.Jeans,
				},
			},
			wantCode: codes.OK,
		},
		{
			desc: "InvalidParent",
			req: &pb.ListProductsRequest{
				Parent:    testresources.Alice.Name, // name of a user
				PageSize:  0,
				PageToken: "",
			},
			wantResp: nil,
			wantCode: codes.InvalidArgument,
		},
		{
			desc: "NegativePageSize",
			req: &pb.ListProductsRequest{
				Parent:    testresources.Bar.Name,
				PageSize:  -1,
				PageToken: "",
			},
			wantResp: nil,
			wantCode: codes.InvalidArgument,
		},
		{
			desc: "PaginationUnimplemented_PositivePageSize",
			req: &pb.ListProductsRequest{
				Parent:    testresources.Bar.Name,
				PageSize:  1,
				PageToken: "",
			},
			wantResp: nil,
			wantCode: codes.Unimplemented,
		},
		{
			desc: "PaginationUnimplemented_NonEmptyPageToken",
			req: &pb.ListProductsRequest{
				Parent:    testresources.Bar.Name,
				PageSize:  0,
				PageToken: "token",
			},
			wantResp: nil,
			wantCode: codes.Unimplemented,
		},
	} {
		t.Run(test.desc, func(t *testing.T) {
			resp, err := c.ListProducts(ctx, test.req)
			if diff := cmp.Diff(
				resp, test.wantResp, protocmp.Transform(),
				protocmp.FilterField(new(pb.ListProductsResponse), "products", protocmp.SortRepeated(productLess)),
			); diff != "" {
				t.Errorf("c.ListProducts(%v, %v) resp != test.wantResp (-got +want)\n%s", ctx, test.req, diff)
			}
			if got := status.Code(err); got != test.wantCode {
				t.Errorf("status.Code(%v) = %v; want %v", err, got, test.wantCode)
			}
		})
	}
}

func TestService_CreateProduct(t *testing.T) {
	ctx := context.Background()
	for _, test := range []struct {
		desc        string
		req         *pb.CreateProductRequest
		wantProduct *pb.Product
		wantCode    codes.Code
	}{
		{
			desc: "OK",
			req: &pb.CreateProductRequest{
				Parent:  testresources.Pharmacy.Name,
				Product: testresources.Pills,
			},
			wantProduct: testresources.Pills,
			wantCode:    codes.OK,
		},
	} {
		t.Run(test.desc, func(t *testing.T) {
			c := serveAndDial(ctx, t, seed(ctx, t))
			product, err := c.CreateProduct(ctx, test.req)
			if diff := cmp.Diff(
				product, test.wantProduct, protocmp.Transform(),
				protocmp.IgnoreFields(new(pb.Product), "name"),
			); diff != "" {
				t.Errorf("c.CreateProduct(%v, %v) product != test.wantProduct (-got +want)\n%s", ctx, test.req, diff)
			}
			if got := status.Code(err); got != test.wantCode {
				t.Errorf("status.Code(%v) = %v; want %v", err, got, test.wantCode)
			}
		})
	}
}

func TestService_UpdateProduct(t *testing.T) {
	ctx := context.Background()
	// Test scenario(s) where the update is successful.
	t.Run("OK", func(t *testing.T) {
		oldBeer := products.Clone(testresources.Beer)
		newBeer := products.Clone(oldBeer)
		newBeer.DisplayName = "New Beer"
		newBeer.FullPriceCents = oldBeer.FullPriceCents - 1000
		newBeer.DiscountPriceCents = oldBeer.DiscountPriceCents - 1000
		for _, test := range []struct {
			desc string
			req  *pb.UpdateProductRequest
			want *pb.Product
		}{
			{
				desc: "NoOp_NilUpdateMask",
				req: &pb.UpdateProductRequest{
					Product:    oldBeer,
					UpdateMask: nil,
				},
				want: oldBeer,
			},
			{
				desc: "NoOp_AllPaths",
				req: &pb.UpdateProductRequest{
					Product:    oldBeer,
					UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"display_name", "full_price_cents", "discount_price_cents"}},
				},
				want: oldBeer,
			},
			{
				desc: "NoOp_NoPaths",
				req: &pb.UpdateProductRequest{
					Product:    newBeer,
					UpdateMask: &fieldmaskpb.FieldMask{Paths: nil},
				},
				want: oldBeer,
			},
			{
				desc: "FullUpdate_NilUpdateMask",
				req: &pb.UpdateProductRequest{
					Product:    newBeer,
					UpdateMask: nil,
				},
				want: newBeer,
			},
			{
				desc: "FullUpdate_AllPaths",
				req: &pb.UpdateProductRequest{
					Product:    newBeer,
					UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"display_name", "full_price_cents", "discount_price_cents"}},
				},
				want: newBeer,
			},
			{
				desc: "PartialUpdate_FullProduct_DisplayName",
				req: &pb.UpdateProductRequest{
					Product:    newBeer,
					UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"display_name"}},
				},
				want: func() *pb.Product {
					product := products.Clone(oldBeer)
					product.DisplayName = newBeer.DisplayName
					return product
				}(),
			},
			{
				desc: "PartialUpdate_FullProduct_FullPriceCents",
				req: &pb.UpdateProductRequest{
					Product:    newBeer,
					UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"full_price_cents"}},
				},
				want: func() *pb.Product {
					product := products.Clone(oldBeer)
					product.FullPriceCents = newBeer.FullPriceCents
					return product
				}(),
			},
			{
				desc: "PartialUpdate_FullProduct_DiscountPriceCents",
				req: &pb.UpdateProductRequest{
					Product:    newBeer,
					UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"discount_price_cents"}},
				},
				want: func() *pb.Product {
					product := products.Clone(oldBeer)
					product.DiscountPriceCents = newBeer.DiscountPriceCents
					return product
				}(),
			},
			{
				desc: "PartialUpdate_PartialProduct_DisplayName",
				req: &pb.UpdateProductRequest{
					Product: &pb.Product{
						Name:        oldBeer.Name,
						DisplayName: "New Beer",
					},
					UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"display_name"}},
				},
				want: func() *pb.Product {
					product := products.Clone(oldBeer)
					product.DisplayName = "New Beer"
					return product
				}(),
			},
			{
				desc: "PartialUpdate_PartialProduct_FullPriceCents",
				req: &pb.UpdateProductRequest{
					Product: &pb.Product{
						Name:           oldBeer.Name,
						FullPriceCents: -4000,
					},
					UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"full_price_cents"}},
				},
				want: func() *pb.Product {
					product := products.Clone(oldBeer)
					product.FullPriceCents = -4000
					return product
				}(),
			},
			{
				desc: "PartialUpdate_PartialProduct_DiscountPriceCents",
				req: &pb.UpdateProductRequest{
					Product: &pb.Product{
						Name:               oldBeer.Name,
						DiscountPriceCents: -2000,
					},
					UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"discount_price_cents"}},
				},
				want: func() *pb.Product {
					product := products.Clone(oldBeer)
					product.DiscountPriceCents = -2000
					return product
				}(),
			},
		} {
			t.Run(test.desc, func(t *testing.T) {
				c := serveAndDial(ctx, t, seed(ctx, t))
				product, err := c.UpdateProduct(ctx, test.req)
				if diff := cmp.Diff(product, test.want, protocmp.Transform()); diff != "" {
					t.Errorf("c.UpdateProduct(%v, %v) product != test.want (-got +want)\n%s", ctx, test.req, diff)
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
			req  *pb.UpdateProductRequest
			want codes.Code
		}{
			{
				desc: "EmptyDisplayName",
				req: &pb.UpdateProductRequest{
					Product: func() *pb.Product {
						beer := products.Clone(testresources.Beer)
						beer.DisplayName = ""
						return beer
					}(),
					UpdateMask: nil,
				},
				want: codes.InvalidArgument,
			},
			{
				desc: "FullPriceCentsPositive",
				req: &pb.UpdateProductRequest{
					Product: func() *pb.Product {
						beer := products.Clone(testresources.Beer)
						beer.FullPriceCents = 1000
						return beer
					}(),
					UpdateMask: nil,
				},
				want: codes.InvalidArgument,
			},
			{
				desc: "DiscountPriceCentsPositive",
				req: &pb.UpdateProductRequest{
					Product: func() *pb.Product {
						beer := products.Clone(testresources.Beer)
						beer.DiscountPriceCents = 1000
						return beer
					}(),
					UpdateMask: nil,
				},
				want: codes.InvalidArgument,
			},
			{
				desc: "DiscountPriceCentsLargerThanFullPriceCents",
				req: &pb.UpdateProductRequest{
					Product: func() *pb.Product {
						beer := products.Clone(testresources.Beer)
						beer.DiscountPriceCents = beer.FullPriceCents - 1000
						return beer
					}(),
					UpdateMask: nil,
				},
				want: codes.InvalidArgument,
			},
			{
				desc: "NotFound",
				req: &pb.UpdateProductRequest{
					Product:    testresources.Cocktail,
					UpdateMask: nil,
				},
				want: codes.NotFound,
			},
			{
				desc: "InvalidUpdateMask",
				req: &pb.UpdateProductRequest{
					Product:    testresources.Beer,
					UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"invalid"}},
				},
				want: codes.InvalidArgument,
			},
		} {
			t.Run(test.desc, func(t *testing.T) {
				c := serveAndDial(ctx, t, seed(ctx, t))
				_, err := c.UpdateProduct(ctx, test.req)
				if got := status.Code(err); got != test.want {
					t.Errorf("status.Code(%v) = %v; want %v", err, got, test.want)
				}
			})
		}
	})
}

func TestService_DeleteProduct(t *testing.T) {
	ctx := context.Background()
	// Test scenario(s) where the delete is successful.
	t.Run("OK", func(t *testing.T) {
		c := serveAndDial(ctx, t, seed(ctx, t))
		{
			req := &pb.DeleteProductRequest{Name: testresources.Beer.Name}
			_, err := c.DeleteProduct(ctx, req)
			if got, want := status.Code(err), codes.OK; got != want {
				t.Errorf("status.Code(%v) = %v; want %v", err, got, want)
			}
		}
		{
			req := &pb.GetProductRequest{Name: testresources.Beer.Name}
			_, err := c.GetProduct(ctx, req)
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
			req  *pb.DeleteProductRequest
			want codes.Code
		}{
			{
				desc: "EmptyName",
				req:  &pb.DeleteProductRequest{Name: ""},
				want: codes.InvalidArgument,
			},
			{
				desc: "NotFound",
				req:  &pb.DeleteProductRequest{Name: testresources.Cocktail.Name},
				want: codes.NotFound,
			},
		} {
			t.Run(test.desc, func(t *testing.T) {
				_, err := c.DeleteProduct(ctx, test.req)
				if got := status.Code(err); got != test.want {
					t.Errorf("status.Code(%v) = %v; want %v", err, got, test.want)
				}
			})
		}
	})
}
