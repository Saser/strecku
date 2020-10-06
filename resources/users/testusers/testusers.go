package testusers

import pb "github.com/Saser/strecku/api/v1"

var (
	Alice = &pb.User{
		Name:         "users/6f2d193c-1460-491d-8157-7dd9535526c6",
		EmailAddress: "alice@example.com",
		DisplayName:  "Alice",
		Superuser:    false,
	}
	AlicePassword = "Alice's password"

	Bob = &pb.User{
		Name:         "users/1c0334cf-9eb2-40b2-accc-43157fedb7ca",
		EmailAddress: "bob@example.com",
		DisplayName:  "Bob",
		Superuser:    false,
	}
	BobPassword = "Bob's password"

	Carol = &pb.User{
		Name:         "users/0bcd2540-e067-41ae-951a-bf95db0817fb",
		EmailAddress: "carol@example.com",
		DisplayName:  "Carol",
	}
	CarolPassword = "Carol's password"
)
