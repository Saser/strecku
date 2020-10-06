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
)
