package testproducts

import (
	pb "github.com/Saser/strecku/api/v1"
	"github.com/Saser/strecku/resources/stores/teststores"
)

var (
	Bar_Beer = &pb.Product{
		Name:               teststores.Bar.Name + "/products/66a691bd-1387-444c-8623-dde2b0a13aee",
		Parent:             teststores.Bar.Name,
		DisplayName:        "Beer",
		FullPriceCents:     -5000,
		DiscountPriceCents: -2500,
	}
	Bar_Cocktail = &pb.Product{
		Name:               teststores.Bar.Name + "/products/66120556-8ebc-4175-904c-f0e7b227d844",
		Parent:             teststores.Bar.Name,
		DisplayName:        "Cocktail",
		FullPriceCents:     -10000,
		DiscountPriceCents: -7500,
	}
)
