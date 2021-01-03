package testresources

import pb "github.com/Saser/strecku/api/v1"

var (
	Bar_Alice_Beer1 = &pb.Purchase{
		Name: Bar.Name + "/purchases/8e386dfa-1085-4d0d-99a1-33540cec25c3",
		User: Alice.Name,
		Lines: []*pb.Purchase_Line{
			{
				Description: Beer.DisplayName,
				Quantity:    1,
				PriceCents:  Beer.FullPriceCents,
				Product:     Beer.Name,
			},
		},
	}
	Bar_Alice_Cocktail1 = &pb.Purchase{
		Name: Bar.Name + "/purchases/926eeccc-0d6d-4bab-a5da-19e79995aeb1",
		User: Alice.Name,
		Lines: []*pb.Purchase_Line{
			{
				Description: Cocktail.DisplayName,
				Quantity:    1,
				PriceCents:  Cocktail.FullPriceCents,
				Product:     Cocktail.Name,
			},
		},
	}
	Bar_Alice_Beer2_Cocktail2 = &pb.Purchase{
		Name: Bar.Name + "/purchases/ecfb87bd-9ba4-40e4-a62c-75a14227b037",
		User: Alice.Name,
		Lines: []*pb.Purchase_Line{
			{
				Description: Beer.DisplayName,
				Quantity:    2,
				PriceCents:  Beer.FullPriceCents,
				Product:     Beer.Name,
			},
			{
				Description: Cocktail.DisplayName,
				Quantity:    2,
				PriceCents:  Cocktail.FullPriceCents,
				Product:     Cocktail.Name,
			},
		},
	}
	Mall_Alice_Jeans1 = &pb.Purchase{
		Name: Mall.Name + "/purchases/f4037d30-6ecc-4fe6-9b90-0caba0594335",
		User: Alice.Name,
		Lines: []*pb.Purchase_Line{
			{
				Description: Jeans.DisplayName,
				Quantity:    1,
				PriceCents:  Jeans.FullPriceCents,
				Product:     Jeans.Name,
			},
		},
	}
)
