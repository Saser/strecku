package testmemberships

import (
	"github.com/Saser/strecku/resources/stores/teststores"
	"github.com/Saser/strecku/resources/users/testusers"
	streckuv1 "github.com/Saser/strecku/saser/strecku/v1"
)

var Alice_Bar = &streckuv1.Membership{
	Name:          "memberships/bd96a64b-7da2-4254-a315-b82675548a8f",
	User:          testusers.Alice.Name,
	Store:         teststores.Bar.Name,
	Administrator: false,
}
