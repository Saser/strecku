package testresources

import pb "github.com/Saser/strecku/api/v1"

var (
	Bar_Alice_Payment = &pb.Payment{
		Name:        Bar.Name + "/payments/5cbd5d1c-08b8-4611-b5ac-936583b0fcb9",
		User:        Alice.Name,
		Description: "Alice's payment",
		AmountCents: 10000,
	}
	Bar_Bob_Payment = &pb.Payment{
		Name:        Bar.Name + "/payments/605b61bd-9906-4427-a71e-730cc4b60c58",
		User:        Bob.Name,
		Description: "Bob's payment",
		AmountCents: 20000,
	}
	Bar_Carol_Payment = &pb.Payment{
		Name:        Bar.Name + "/payments/12490c0b-c72d-47c7-bf48-17787c18e173",
		User:        Carol.Name,
		Description: "Carol's payment",
		AmountCents: 20000,
	}
)
