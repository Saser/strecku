package testmemberships

import (
	"testing"

	"github.com/Saser/strecku/resources/memberships"
	streckuv1 "github.com/Saser/strecku/saser/strecku/v1"
)

func TestValid(t *testing.T) {
	for _, membership := range []*streckuv1.Membership{
		Alice_Bar,
	} {
		if err := memberships.Validate(membership); err != nil {
			t.Errorf("memberships.Validate(%v) = %v; want nil", membership, err)
		}
	}
}
