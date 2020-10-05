package testusers

import streckuv1 "github.com/Saser/strecku/saser/strecku/v1"

var (
	Alice = &streckuv1.User{
		Name:         "users/6f2d193c-1460-491d-8157-7dd9535526c6",
		EmailAddress: "alice@example.com",
		DisplayName:  "Alice",
		Superuser:    false,
	}
	AlicePassword = "Alice's password"

	Bob = &streckuv1.User{
		Name:         "users/1c0334cf-9eb2-40b2-accc-43157fedb7ca",
		EmailAddress: "bob@example.com",
		DisplayName:  "Bob",
		Superuser:    false,
	}
	BobPassword = "Bob's password"

	Carol = &streckuv1.User{
		Name:         "users/0bcd2540-e067-41ae-951a-bf95db0817fb",
		EmailAddress: "carol@example.com",
		DisplayName:  "Carol",
	}
	CarolPassword = "Carol's password"

	David = &streckuv1.User{
		Name:         "users/9c05668c-a75b-4667-a9a9-a9faa93993fe",
		EmailAddress: "david@example.com",
		DisplayName:  "David",
	}
	DavidPassword = "David's password"
)
