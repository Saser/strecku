package testmemberships

import (
	pb "github.com/Saser/strecku/api/v1"
	"github.com/Saser/strecku/resources/stores/teststores"
	"github.com/Saser/strecku/resources/users/testusers"
)

var (
	Alice_Bar = &pb.Membership{
		Name:          "memberships/bd96a64b-7da2-4254-a315-b82675548a8f",
		User:          testusers.Alice.Name,
		Store:         teststores.Bar.Name,
		Administrator: false,
	}
	Alice_Mall = &pb.Membership{
		Name:          "memberships/9681bce8-14da-4a11-a812-8245ccd1c911",
		User:          testusers.Alice.Name,
		Store:         teststores.Mall.Name,
		Administrator: false,
	}
	Bob_Bar = &pb.Membership{
		Name:          "memberships/ad8a0fc4-1482-4f00-b69c-f6d26104e504",
		User:          testusers.Bob.Name,
		Store:         teststores.Bar.Name,
		Administrator: false,
	}
	Bob_Mall = &pb.Membership{
		Name:          "memberships/70d30eac-d059-4712-9d16-f5ec5926d4f0",
		User:          testusers.Bob.Name,
		Store:         teststores.Mall.Name,
		Administrator: false,
	}
)
