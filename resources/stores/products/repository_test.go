package products

import (
	"context"
	"fmt"
	"testing"

	pb "github.com/Saser/strecku/api/v1"
	"github.com/Saser/strecku/resourcename"
	"github.com/Saser/strecku/testresources"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/protobuf/testing/protocmp"
)

func productLess(p1, p2 *pb.Product) bool {
	return p1.Name < p2.Name
}

func TestNotFoundError_Error(t *testing.T) {
	err := &NotFoundError{Name: testresources.Beer.Name}
	want := fmt.Sprintf("product not found: %q", testresources.Beer.Name)
	if got := err.Error(); !cmp.Equal(got, want) {
		t.Errorf("err.Error() = %q; want %q", got, want)
	}
}

func TestNotFoundError_Is(t *testing.T) {
	for _, test := range []struct {
		err    *NotFoundError
		target error
		want   bool
	}{
		{
			err:    &NotFoundError{Name: testresources.Beer.Name},
			target: &NotFoundError{Name: testresources.Beer.Name},
			want:   true,
		},
		{
			err:    &NotFoundError{Name: testresources.Beer.Name},
			target: &NotFoundError{Name: testresources.Cocktail.Name},
			want:   false,
		},
		{
			err:    &NotFoundError{Name: testresources.Beer.Name},
			target: fmt.Errorf("product not found: %q", testresources.Beer.Name),
			want:   false,
		},
	} {
		if got := test.err.Is(test.target); got != test.want {
			t.Errorf("test.err.Is(%v) = %v; want %v", test.target, got, test.want)
		}
	}
}

func TestExistsError_Error(t *testing.T) {
	err := &ExistsError{Name: testresources.Beer.Name}
	want := fmt.Sprintf("product exists: %q", testresources.Beer.Name)
	if got := err.Error(); !cmp.Equal(got, want) {
		t.Errorf("err.Error() = %q; want %q", got, want)
	}
}

func TestExistsError_Is(t *testing.T) {
	for _, test := range []struct {
		err    *ExistsError
		target error
		want   bool
	}{
		{
			err:    &ExistsError{Name: testresources.Beer.Name},
			target: &ExistsError{Name: testresources.Beer.Name},
			want:   true,
		},
		{
			err:    &ExistsError{Name: testresources.Beer.Name},
			target: &ExistsError{Name: testresources.Cocktail.Name},
			want:   false,
		},
		{
			err:    &ExistsError{Name: testresources.Beer.Name},
			target: fmt.Errorf("product exists: %q", testresources.Beer.Name),
			want:   false,
		},
	} {
		if got := test.err.Is(test.target); got != test.want {
			t.Errorf("test.err.Is(%v) = %v; want %v", test.target, got, test.want)
		}
	}
}

func TestNewRepository(t *testing.T) {
	NewRepository()
}

func TestRepository_LookupProduct(t *testing.T) {
	ctx := context.Background()
	r := SeedRepository(t, []*pb.Product{testresources.Beer})
	for _, test := range []struct {
		desc        string
		name        string
		wantProduct *pb.Product
		wantErr     error
	}{
		{
			desc:        "OK",
			name:        testresources.Beer.Name,
			wantProduct: testresources.Beer,
			wantErr:     nil,
		},
		{
			desc:        "EmptyName",
			name:        "",
			wantProduct: nil,
			wantErr:     resourcename.ErrInvalidName,
		},
		{
			desc:        "NotFound",
			name:        testresources.Cocktail.Name,
			wantProduct: nil,
			wantErr:     &NotFoundError{Name: testresources.Cocktail.Name},
		},
	} {
		t.Run(test.desc, func(t *testing.T) {
			store, err := r.LookupProduct(ctx, test.name)
			if diff := cmp.Diff(store, test.wantProduct, protocmp.Transform()); diff != "" {
				t.Errorf("r.LookupProduct(%v, %q) product != test.wantProduct (-got +want)\n%s", ctx, test.name, diff)
			}
			if got, want := err, test.wantErr; !cmp.Equal(got, want, cmpopts.EquateErrors()) {
				t.Errorf("r.LookupProduct(%v, %q) err = %v; want %v", ctx, test.name, got, want)
			}
		})
	}
}

func TestRepository_ListProducts(t *testing.T) {
	ctx := context.Background()
	want := []*pb.Product{
		testresources.Beer,
		testresources.Cocktail,
		testresources.Pills,
		testresources.Lotion,
	}
	r := SeedRepository(t, want)
	stores, err := r.ListProducts(ctx)
	if diff := cmp.Diff(
		stores, want, protocmp.Transform(),
		cmpopts.SortSlices(productLess),
	); diff != "" {
		t.Errorf("r.ListProducts(%v) stores != want (-got +want)\n%s", ctx, diff)
	}
	if err != nil {
		t.Errorf("r.ListProducts(%v) err = %v; want nil", ctx, err)
	}
}

func TestRepository_FilterProducts(t *testing.T) {
	ctx := context.Background()
	r := SeedRepository(t, []*pb.Product{
		testresources.Beer,
		testresources.Cocktail,
		testresources.Pills,
		testresources.Lotion,
	})
	for _, test := range []struct {
		name      string
		predicate func(*pb.Product) bool
		want      []*pb.Product
	}{
		{
			name:      "NoneMatching",
			predicate: func(*pb.Product) bool { return false },
			want:      nil,
		},
		{
			name:      "OneMatching",
			predicate: func(product *pb.Product) bool { return product.Name == testresources.Beer.Name },
			want: []*pb.Product{
				testresources.Beer,
			},
		},
		{
			name: "MultipleMatching",
			predicate: func(product *pb.Product) bool {
				return product.Name == testresources.Beer.Name || product.Name == testresources.Cocktail.Name
			},
			want: []*pb.Product{
				testresources.Beer,
				testresources.Cocktail,
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			stores, err := r.FilterProducts(ctx, test.predicate)
			if diff := cmp.Diff(
				stores, test.want, protocmp.Transform(),
				cmpopts.EquateEmpty(),
				cmpopts.SortSlices(productLess),
			); diff != "" {
				t.Errorf("r.FilterProducts(%v, test.predicate) stores != test.want (-got +want)\n%s", ctx, diff)
			}
			if got, want := err, error(nil); !cmp.Equal(got, want) {
				t.Errorf("r.FilterProducts(%v, test.predicate) err = %v; want %v", ctx, got, want)
			}
		})
	}
}

func TestRepository_CreateProduct(t *testing.T) {
	ctx := context.Background()
	duplicateName := Clone(testresources.Cocktail)
	duplicateName.Name = testresources.Beer.Name
	for _, test := range []struct {
		desc    string
		product *pb.Product
		want    error
	}{
		{
			desc:    "OneProductOK",
			product: testresources.Cocktail,
			want:    nil,
		},
		{
			desc:    "DuplicateName",
			product: duplicateName,
			want:    &ExistsError{Name: testresources.Beer.Name},
		},
	} {
		t.Run(test.desc, func(t *testing.T) {
			r := SeedRepository(t, []*pb.Product{testresources.Beer})
			if got := r.CreateProduct(ctx, test.product); !cmp.Equal(got, test.want) {
				t.Errorf("r.CreateProduct(%v, %v) = %v; want %v", ctx, test.product, got, test.want)
			}
		})
	}
}

func TestRepository_UpdateProduct(t *testing.T) {
	ctx := context.Background()
	// Test scenario where the update is successful.
	t.Run("OK", func(t *testing.T) {
		r := SeedRepository(t, []*pb.Product{testresources.Beer})
		oldBeer := Clone(testresources.Beer)
		newBeer := Clone(oldBeer)
		newBeer.DisplayName = "New Beer"
		newBeer.FullPriceCents = -1500
		newBeer.DiscountPriceCents = -1000
		if err := r.UpdateProduct(ctx, newBeer); err != nil {
			t.Errorf("r.UpdateProduct(%v, %v) = %v; want nil", ctx, newBeer, err)
		}
		store, err := r.LookupProduct(ctx, newBeer.Name)
		if diff := cmp.Diff(store, newBeer, protocmp.Transform()); diff != "" {
			t.Errorf("r.LookupProduct(%v, %q) store != newBeer (-got +want)\n%s", ctx, newBeer.Name, diff)
		}
		if err != nil {
			t.Errorf("r.LookupProduct(%v, %q) err = %v; want nil", ctx, newBeer.Name, err)
		}
	})

	// Test scenario where the update fails.
	t.Run("Errors", func(t *testing.T) {
		r := SeedRepository(t, []*pb.Product{testresources.Beer})
		for _, test := range []struct {
			desc   string
			modify func(beer *pb.Product)
			want   error
		}{
			{
				desc:   "EmptyName",
				modify: func(beer *pb.Product) { beer.Name = "" },
				want:   resourcename.ErrInvalidName,
			},
			{
				desc:   "EmptyDisplayName",
				modify: func(beer *pb.Product) { beer.DisplayName = "" },
				want:   ErrDisplayNameEmpty,
			},
			{
				desc:   "PositiveFullPrice",
				modify: func(beer *pb.Product) { beer.FullPriceCents = 1000 },
				want:   ErrFullPricePositive,
			},
			{
				desc:   "PositiveDiscountPrice",
				modify: func(beer *pb.Product) { beer.DiscountPriceCents = 1000 },
				want:   ErrDiscountPricePositive,
			},
			{
				desc:   "DiscountPriceHigherThanFullPrice",
				modify: func(beer *pb.Product) { beer.DiscountPriceCents = beer.FullPriceCents - 1000 },
				want:   ErrDiscountPriceHigherThanFullPrice,
			},
			{
				desc:   "NotFound",
				modify: func(beer *pb.Product) { beer.Name = testresources.Cocktail.Name },
				want:   &NotFoundError{Name: testresources.Cocktail.Name},
			},
		} {
			t.Run(test.desc, func(t *testing.T) {
				updated := Clone(testresources.Beer)
				test.modify(updated)
				if got := r.UpdateProduct(ctx, updated); !cmp.Equal(got, test.want, cmpopts.EquateErrors()) {
					t.Errorf("r.UpdateProduct(%v, %v) = %v; want %v", ctx, updated, got, test.want)
				}
			})
		}
	})
}

func TestRepository_DeleteProduct(t *testing.T) {
	ctx := context.Background()

	t.Run("OK", func(t *testing.T) {
		r := SeedRepository(t, []*pb.Product{
			testresources.Beer,
			testresources.Cocktail,
		})
		if err := r.DeleteProduct(ctx, testresources.Beer.Name); err != nil {
			t.Errorf("r.DeleteProduct(%v, %q) = %v; want nil", ctx, testresources.Beer.Name, err)
		}
		for _, test := range []struct {
			desc        string
			name        string
			wantProduct *pb.Product
			wantErr     error
		}{
			{
				desc:        "LookupDeleted",
				name:        testresources.Beer.Name,
				wantProduct: nil,
				wantErr:     &NotFoundError{Name: testresources.Beer.Name},
			},
			{
				desc:        "LookupExisting",
				name:        testresources.Cocktail.Name,
				wantProduct: testresources.Cocktail,
				wantErr:     nil,
			},
		} {
			t.Run(test.desc, func(t *testing.T) {
				store, err := r.LookupProduct(ctx, test.name)
				if diff := cmp.Diff(store, test.wantProduct, protocmp.Transform()); diff != "" {
					t.Errorf("r.LookupProduct(%v, %q) store != test.wantProduct (-got +want)\n%s", ctx, test.name, diff)
				}
				if !cmp.Equal(err, test.wantErr, cmpopts.EquateErrors()) {
					t.Errorf("r.LookupProduct(%v, %q) err = %v; want %v", ctx, test.name, err, test.wantErr)
				}
			})
		}
	})
	t.Run("Errors", func(t *testing.T) {
		r := SeedRepository(t, []*pb.Product{testresources.Beer})
		for _, test := range []struct {
			desc string
			name string
			want error
		}{
			{
				desc: "EmptyName",
				name: "",
				want: resourcename.ErrInvalidName,
			},
			{
				desc: "NotFound",
				name: testresources.Cocktail.Name,
				want: &NotFoundError{Name: testresources.Cocktail.Name},
			},
		} {
			t.Run(test.desc, func(t *testing.T) {
				if got := r.DeleteProduct(ctx, test.name); !cmp.Equal(got, test.want, cmpopts.EquateErrors()) {
					t.Errorf("r.DeleteProduct(%v, %q) = %v; want %v", ctx, test.name, got, test.want)
				}
			})
		}
	})
}
