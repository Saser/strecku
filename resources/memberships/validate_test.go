package memberships

import (
	"testing"

	pb "github.com/Saser/strecku/api/v1"
	"github.com/Saser/strecku/resources/memberships/testmemberships"
	"github.com/Saser/strecku/resources/stores"
	"github.com/Saser/strecku/resources/stores/teststores"
	"github.com/Saser/strecku/resources/users"
	"github.com/Saser/strecku/resources/users/testusers"
)

func TestValidate(t *testing.T) {
	for _, test := range []struct {
		membership *pb.Membership
		want       error
	}{
		{
			membership: &pb.Membership{
				Name:          testmemberships.Alice_Bar.Name,
				User:          "",
				Store:         teststores.Bar.Name,
				Administrator: false,
				Discount:      false,
			},
			want: users.ErrNameEmpty,
		},
		{
			membership: &pb.Membership{
				Name:          testmemberships.Alice_Bar.Name,
				User:          testusers.Alice.Name,
				Store:         "",
				Administrator: false,
				Discount:      false,
			},
			want: stores.ErrNameEmpty,
		},
	} {
		if got := Validate(test.membership); got != test.want {
			t.Errorf("Validate(%v) = %v; want %v", test.membership, got, test.want)
		}
	}
}
