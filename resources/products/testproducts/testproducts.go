package testproducts

import (
	pb "github.com/Saser/strecku/api/v1"
	"github.com/Saser/strecku/resources/stores/teststores"
)

var Bar_Beer = &pb.Product{
	Name:               teststores.Bar.Name + "/products/66a691bd-1387-444c-8623-dde2b0a13aee",
	Parent:             teststores.Bar.Name,
	DisplayName:        "Beer",
	FullPriceCents:     -5000,
	DiscountPriceCents: -2500,
}
