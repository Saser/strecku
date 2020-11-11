package testresources

import pb "github.com/Saser/strecku/api/v1"

var (
	Beer = &pb.Product{
		Name:               Bar.Name + "/products/66a691bd-1387-444c-8623-dde2b0a13aee",
		DisplayName:        "Beer",
		FullPriceCents:     -5000,
		DiscountPriceCents: -2500,
	}
	Cocktail = &pb.Product{
		Name:               Bar.Name + "/products/66120556-8ebc-4175-904c-f0e7b227d844",
		DisplayName:        "Cocktail",
		FullPriceCents:     -10000,
		DiscountPriceCents: -7500,
	}
	Jeans = &pb.Product{
		Name:               Mall.Name + "/products/032b6523-9503-4ac4-95c5-622f723d91f4",
		DisplayName:        "Jeans",
		FullPriceCents:     -50000,
		DiscountPriceCents: -40000,
	}
	Shirt = &pb.Product{
		Name:               Mall.Name + "/products/110436c0-ab22-4813-96ae-e079ab90a5e6",
		DisplayName:        "Shirt",
		FullPriceCents:     -30000,
		DiscountPriceCents: -25000,
	}
	Pills = &pb.Product{
		Name:               Pharmacy.Name + "/products/156773ff-4424-446b-8d4d-ce5808359386",
		DisplayName:        "Pills",
		FullPriceCents:     -1000,
		DiscountPriceCents: -1000,
	}
	Lotion = &pb.Product{
		Name:               Pharmacy.Name + "/products/03c826e1-156d-4767-a564-3087e5deff05",
		DisplayName:        "Lotion",
		FullPriceCents:     -2000,
		DiscountPriceCents: -1500,
	}
)
