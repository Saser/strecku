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
	Pharmacy_Pills = &pb.Product{
		Name:               teststores.Pharmacy.Name + "/products/156773ff-4424-446b-8d4d-ce5808359386",
		Parent:             teststores.Pharmacy.Name,
		DisplayName:        "Pills",
		FullPriceCents:     -1000,
		DiscountPriceCents: -1000,
	}
	Pharmacy_Lotion = &pb.Product{
		Name:               teststores.Pharmacy.Name + "/products/03c826e1-156d-4767-a564-3087e5deff05",
		Parent:             teststores.Pharmacy.Name,
		DisplayName:        "Lotion",
		FullPriceCents:     -2000,
		DiscountPriceCents: -1500,
	}
)
