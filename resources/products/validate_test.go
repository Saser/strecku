package products

import (
	"testing"

	"github.com/Saser/strecku/resources/stores"

	pb "github.com/Saser/strecku/api/v1"
	"github.com/Saser/strecku/resources/products/testproducts"
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
				Name:               testproducts.Bar_Beer.Name,
				Parent:             "",
				DisplayName:        testproducts.Bar_Beer.DisplayName,
				FullPriceCents:     testproducts.Bar_Beer.FullPriceCents,
				DiscountPriceCents: testproducts.Bar_Beer.DiscountPriceCents,
			},
			want: stores.ErrNameEmpty,
		},
		{
			product: &pb.Product{
				Name:               testproducts.Bar_Beer.Name,
				Parent:             testproducts.Bar_Beer.Parent,
				DisplayName:        "",
				FullPriceCents:     testproducts.Bar_Beer.FullPriceCents,
				DiscountPriceCents: testproducts.Bar_Beer.DiscountPriceCents,
			},
			want: ErrDisplayNameEmpty,
		},
		{
			product: &pb.Product{
				Name:               testproducts.Bar_Beer.Name,
				Parent:             testproducts.Bar_Beer.Parent,
				DisplayName:        testproducts.Bar_Beer.DisplayName,
				FullPriceCents:     10,
				DiscountPriceCents: testproducts.Bar_Beer.DiscountPriceCents,
			},
			want: ErrFullPricePositive,
		},
		{
			product: &pb.Product{
				Name:               testproducts.Bar_Beer.Name,
				Parent:             testproducts.Bar_Beer.Parent,
				DisplayName:        testproducts.Bar_Beer.DisplayName,
				FullPriceCents:     testproducts.Bar_Beer.FullPriceCents,
				DiscountPriceCents: 10,
			},
			want: ErrDiscountPricePositive,
		},
		{
			product: &pb.Product{
				Name:               testproducts.Bar_Beer.Name,
				Parent:             testproducts.Bar_Beer.Parent,
				DisplayName:        testproducts.Bar_Beer.DisplayName,
				FullPriceCents:     testproducts.Bar_Beer.FullPriceCents,
				DiscountPriceCents: testproducts.Bar_Beer.FullPriceCents - 10,
			},
			want: ErrDiscountPriceHigherThanFullPrice,
		},
		//{
		//	product: &pb.Product{
		//		Name:               testproducts.Bar_Beer.Name,
		//		Parent:             testproducts.Bar_Beer.Parent,
		//		DisplayName:        testproducts.Bar_Beer.DisplayName,
		//		FullPriceCents:     testproducts.Bar_Beer.FullPriceCents,
		//		DiscountPriceCents: testproducts.Bar_Beer.DiscountPriceCents,
		//	},
		//	want: nil,
		//},
	} {
		if got := Validate(test.product); !cmp.Equal(got, test.want, cmpopts.EquateErrors()) {
			t.Errorf("Validate(%v) = %v; want %v", test.product, got, test.want)
		}
	}
}
