package products

import (
	"context"
	"fmt"
	"testing"

	pb "github.com/Saser/strecku/api/v1"
	"github.com/Saser/strecku/resources/products/testproducts"
	"github.com/Saser/strecku/resources/stores"
	"github.com/Saser/strecku/resources/stores/teststores"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/protobuf/testing/protocmp"
)

func productLess(p1, p2 *pb.Product) bool {
	return p1.Name < p2.Name
}

func TestNotFoundError_Error(t *testing.T) {
	err := &NotFoundError{Name: testproducts.Bar_Beer.Name}
	want := fmt.Sprintf("product not found: %q", testproducts.Bar_Beer.Name)
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
			err:    &NotFoundError{Name: testproducts.Bar_Beer.Name},
			target: &NotFoundError{Name: testproducts.Bar_Beer.Name},
			want:   true,
		},
		{
			err:    &NotFoundError{Name: testproducts.Bar_Beer.Name},
			target: &NotFoundError{Name: testproducts.Bar_Cocktail.Name},
			want:   false,
		},
		{
			err:    &NotFoundError{Name: testproducts.Bar_Beer.Name},
			target: fmt.Errorf("product not found: %q", testproducts.Bar_Beer.Name),
			want:   false,
		},
	} {
		if got := test.err.Is(test.target); got != test.want {
			t.Errorf("test.err.Is(%v) = %v; want %v", test.target, got, test.want)
		}
	}
}

func TestExistsError_Error(t *testing.T) {
	err := &ExistsError{Name: testproducts.Bar_Beer.Name}
	want := fmt.Sprintf("product exists: %q", testproducts.Bar_Beer.Name)
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
			err:    &ExistsError{Name: testproducts.Bar_Beer.Name},
			target: &ExistsError{Name: testproducts.Bar_Beer.Name},
			want:   true,
		},
		{
			err:    &ExistsError{Name: testproducts.Bar_Beer.Name},
			target: &ExistsError{Name: testproducts.Bar_Cocktail.Name},
			want:   false,
		},
		{
			err:    &ExistsError{Name: testproducts.Bar_Beer.Name},
			target: fmt.Errorf("product exists: %q", testproducts.Bar_Beer.Name),
			want:   false,
		},
	} {
		if got := test.err.Is(test.target); got != test.want {
			t.Errorf("test.err.Is(%v) = %v; want %v", test.target, got, test.want)
		}
	}
}

func TestRepository_LookupProduct(t *testing.T) {
	ctx := context.Background()
	r := SeedRepository(t, []*pb.Product{testproducts.Bar_Beer})
	for _, test := range []struct {
		desc        string
		name        string
		wantProduct *pb.Product
		wantErr     error
	}{
		{
			desc:        "OK",
			name:        testproducts.Bar_Beer.Name,
			wantProduct: testproducts.Bar_Beer,
			wantErr:     nil,
		},
		{
			desc:        "EmptyName",
			name:        "",
			wantProduct: nil,
			wantErr:     ErrNameEmpty,
		},
		{
			desc:        "NotFound",
			name:        testproducts.Bar_Cocktail.Name,
			wantProduct: nil,
			wantErr:     &NotFoundError{Name: testproducts.Bar_Cocktail.Name},
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
		testproducts.Bar_Beer,
		testproducts.Bar_Cocktail,
		testproducts.Pharmacy_Pills,
		testproducts.Pharmacy_Lotion,
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
		testproducts.Bar_Beer,
		testproducts.Bar_Cocktail,
		testproducts.Pharmacy_Pills,
		testproducts.Pharmacy_Lotion,
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
			predicate: func(product *pb.Product) bool { return product.Name == testproducts.Bar_Beer.Name },
			want: []*pb.Product{
				testproducts.Bar_Beer,
			},
		},
		{
			name:      "MultipleMatching",
			predicate: func(product *pb.Product) bool { return product.Parent == teststores.Bar.Name },
			want: []*pb.Product{
				testproducts.Bar_Beer,
				testproducts.Bar_Cocktail,
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
	duplicateName := Clone(testproducts.Bar_Cocktail)
	duplicateName.Name = testproducts.Bar_Beer.Name
	for _, test := range []struct {
		desc    string
		product *pb.Product
		want    error
	}{
		{
			desc:    "OneProductOK",
			product: testproducts.Bar_Cocktail,
			want:    nil,
		},
		{
			desc:    "DuplicateName",
			product: duplicateName,
			want:    &ExistsError{Name: testproducts.Bar_Beer.Name},
		},
	} {
		t.Run(test.desc, func(t *testing.T) {
			r := SeedRepository(t, []*pb.Product{testproducts.Bar_Beer})
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
		r := SeedRepository(t, []*pb.Product{testproducts.Bar_Beer})
		oldBeer := Clone(testproducts.Bar_Beer)
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
		r := SeedRepository(t, []*pb.Product{testproducts.Bar_Beer})
		for _, test := range []struct {
			desc   string
			modify func(beer *pb.Product)
			want   error
		}{
			{
				desc:   "EmptyName",
				modify: func(beer *pb.Product) { beer.Name = "" },
				want:   ErrNameEmpty,
			},
			{
				desc:   "EmptyParent",
				modify: func(beer *pb.Product) { beer.Parent = "" },
				want:   stores.ErrNameEmpty,
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
				desc:   "UpdateParent",
				modify: func(beer *pb.Product) { beer.Parent = teststores.Pharmacy.Name },
				want:   ErrUpdateParent,
			},
			{
				desc:   "NotFound",
				modify: func(beer *pb.Product) { beer.Name = testproducts.Bar_Cocktail.Name },
				want:   &NotFoundError{Name: testproducts.Bar_Cocktail.Name},
			},
		} {
			t.Run(test.desc, func(t *testing.T) {
				updated := Clone(testproducts.Bar_Beer)
				test.modify(updated)
				if got := r.UpdateProduct(ctx, updated); !cmp.Equal(got, test.want, cmpopts.EquateErrors()) {
					t.Errorf("r.UpdateProduct(%v, %v) = %v; want %v", ctx, updated, got, test.want)
				}
			})
		}
	})
}
