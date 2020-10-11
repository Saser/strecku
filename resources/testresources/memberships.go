package testresources

import pb "github.com/Saser/strecku/api/v1"

var (
	Alice_Bar = &pb.Membership{
		Name:          "memberships/bd96a64b-7da2-4254-a315-b82675548a8f",
		User:          Alice.Name,
		Store:         Bar.Name,
		Administrator: false,
		Discount:      false,
	}
	Alice_Mall = &pb.Membership{
		Name:          "memberships/9681bce8-14da-4a11-a812-8245ccd1c911",
		User:          Alice.Name,
		Store:         Mall.Name,
		Administrator: false,
		Discount:      false,
	}
	Bob_Bar = &pb.Membership{
		Name:          "memberships/ad8a0fc4-1482-4f00-b69c-f6d26104e504",
		User:          Bob.Name,
		Store:         Bar.Name,
		Administrator: false,
		Discount:      false,
	}
	Bob_Mall = &pb.Membership{
		Name:          "memberships/70d30eac-d059-4712-9d16-f5ec5926d4f0",
		User:          Bob.Name,
		Store:         Mall.Name,
		Administrator: false,
		Discount:      false,
	}
)
