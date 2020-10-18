package testresources

import (
	"testing"

	pb "github.com/Saser/strecku/api/v1"
	"github.com/Saser/strecku/resources/memberships"
)

func TestMembershipsValid(t *testing.T) {
	for _, membership := range []*pb.Membership{
		Bar_Alice,
		Mall_Alice,
		Bar_Bob,
		Mall_Bob,
	} {
		if err := memberships.Validate(membership); err != nil {
			t.Errorf("memberships.Validate(%v) = %v; want nil", membership, err)
		}
	}
}
