package memberships

import (
	"testing"

	pb "github.com/Saser/strecku/api/v1"
	"github.com/Saser/strecku/resourcename"
	"github.com/Saser/strecku/testresources"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
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
			want: resourcename.ErrInvalidName,
		},
	} {
		if got := Validate(test.membership); !cmp.Equal(got, test.want, cmpopts.EquateErrors()) {
			t.Errorf("Validate(%v) = %v; want %v", test.membership, got, test.want)
		}
	}
}
