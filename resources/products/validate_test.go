package products

import (
	"testing"

	pb "github.com/Saser/strecku/api/v1"
	"github.com/Saser/strecku/resources/stores"
	"github.com/Saser/strecku/resources/testresources"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestValidate(t *testing.T) {
	for _, test := range []struct {
		product *pb.Product
		want    error
	}{
		{
			product: &pb.Product{
				Name:               testresources.Beer.Name,
				Parent:             "",
				DisplayName:        testresources.Beer.DisplayName,
				FullPriceCents:     testresources.Beer.FullPriceCents,
				DiscountPriceCents: testresources.Beer.DiscountPriceCents,
			},
			want: stores.ErrNameInvalidFormat,
		},
		{
			product: &pb.Product{
				Name:               testresources.Beer.Name,
				Parent:             testresources.Beer.Parent,
				DisplayName:        "",
				FullPriceCents:     testresources.Beer.FullPriceCents,
				DiscountPriceCents: testresources.Beer.DiscountPriceCents,
			},
			want: ErrDisplayNameEmpty,
		},
		{
			product: &pb.Product{
				Name:               testresources.Beer.Name,
				Parent:             testresources.Beer.Parent,
				DisplayName:        testresources.Beer.DisplayName,
				FullPriceCents:     10,
				DiscountPriceCents: testresources.Beer.DiscountPriceCents,
			},
			want: ErrFullPricePositive,
		},
		{
			product: &pb.Product{
				Name:               testresources.Beer.Name,
				Parent:             testresources.Beer.Parent,
				DisplayName:        testresources.Beer.DisplayName,
				FullPriceCents:     testresources.Beer.FullPriceCents,
				DiscountPriceCents: 10,
			},
			want: ErrDiscountPricePositive,
		},
		{
			product: &pb.Product{
				Name:               testresources.Beer.Name,
				Parent:             testresources.Beer.Parent,
				DisplayName:        testresources.Beer.DisplayName,
				FullPriceCents:     testresources.Beer.FullPriceCents,
				DiscountPriceCents: testresources.Beer.FullPriceCents - 10,
			},
			want: ErrDiscountPriceHigherThanFullPrice,
		},
	} {
		if got := Validate(test.product); !cmp.Equal(got, test.want, cmpopts.EquateErrors()) {
			t.Errorf("Validate(%v) = %v; want %v", test.product, got, test.want)
		}
	}
}
