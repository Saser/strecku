package testmemberships

import (
	"testing"

	pb "github.com/Saser/strecku/api/v1"
	"github.com/Saser/strecku/resources/memberships"
)

func TestValid(t *testing.T) {
	for _, membership := range []*pb.Membership{
		Alice_Bar,
	} {
		if err := memberships.Validate(membership); err != nil {
			t.Errorf("memberships.Validate(%v) = %v; want nil", membership, err)
		}
	}
}
