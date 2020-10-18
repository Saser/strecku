package memberships

import (
	"testing"

	pb "github.com/Saser/strecku/api/v1"
	"github.com/Saser/strecku/resources/stores"
	"github.com/Saser/strecku/resources/testresources"
	"github.com/Saser/strecku/resources/users"
)

func TestValidate(t *testing.T) {
	for _, test := range []struct {
		membership *pb.Membership
		want       error
	}{
		{
			membership: &pb.Membership{
				Name:          testresources.Alice_Bar.Name,
				User:          "",
				Store:         testresources.Bar.Name,
				Administrator: false,
				Discount:      false,
			},
			want: users.ErrNameInvalidFormat,
		},
		{
			membership: &pb.Membership{
				Name:          testresources.Alice_Bar.Name,
				User:          testresources.Alice.Name,
				Store:         "",
				Administrator: false,
				Discount:      false,
			},
			want: stores.ErrNameInvalidFormat,
		},
	} {
		if got := Validate(test.membership); got != test.want {
			t.Errorf("Validate(%v) = %v; want %v", test.membership, got, test.want)
		}
	}
}
