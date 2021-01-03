package testresources

import pb "github.com/Saser/strecku/api/v1"

var (
	Bar_Alice = &pb.Membership{
		Name:          Bar.Name + "/memberships/bd96a64b-7da2-4254-a315-b82675548a8f",
		User:          Alice.Name,
		Administrator: false,
		Discount:      false,
	}
	Bar_Bob = &pb.Membership{
		Name:          Bar.Name + "/memberships/ad8a0fc4-1482-4f00-b69c-f6d26104e504",
		User:          Bob.Name,
		Administrator: false,
		Discount:      false,
	}
	Mall_Alice = &pb.Membership{
		Name:          Mall.Name + "/memberships/9681bce8-14da-4a11-a812-8245ccd1c911",
		User:          Alice.Name,
		Administrator: false,
		Discount:      false,
	}
	Mall_Bob = &pb.Membership{
		Name:          Mall.Name + "/memberships/70d30eac-d059-4712-9d16-f5ec5926d4f0",
		User:          Bob.Name,
		Administrator: false,
		Discount:      false,
	}
)
