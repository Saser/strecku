package testresources

import pb "github.com/Saser/strecku/api/v1"

var (
	Alice_Beer1 = &pb.Purchase{
		Name:  "purchases/8e386dfa-1085-4d0d-99a1-33540cec25c3",
		User:  Alice.Name,
		Store: Bar.Name,
		Lines: []*pb.Purchase_Line{
			{
				Description: Beer.DisplayName,
				Quantity:    1,
				PriceCents:  Beer.FullPriceCents,
				Product:     Beer.Name,
			},
		},
	}
	Alice_Cocktail1 = &pb.Purchase{
		Name:  "purchases/926eeccc-0d6d-4bab-a5da-19e79995aeb1",
		User:  Alice.Name,
		Store: Bar.Name,
		Lines: []*pb.Purchase_Line{
			{
				Description: Cocktail.DisplayName,
				Quantity:    1,
				PriceCents:  Cocktail.FullPriceCents,
				Product:     Cocktail.Name,
			},
		},
	}
	Alice_Beer2_Cocktail2 = &pb.Purchase{
		Name:  "purchases/ecfb87bd-9ba4-40e4-a62c-75a14227b037",
		User:  Alice.Name,
		Store: Bar.Name,
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
)
