package memberships

import (
	"testing"

	pb "github.com/Saser/strecku/api/v1"
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
				Name:          testresources.Bar_Alice.Name,
				User:          "",
				Administrator: false,
				Discount:      false,
			},
			want: users.ErrNameInvalidFormat,
		},
	} {
		if got := Validate(test.membership); got != test.want {
			t.Errorf("Validate(%v) = %v; want %v", test.membership, got, test.want)
		}
	}
}
