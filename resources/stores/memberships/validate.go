package memberships

import (
	pb "github.com/Saser/strecku/api/v1"
	"github.com/Saser/strecku/resources/users"
)

func Validate(membership *pb.Membership) error {
	if err := ValidateName(membership.Name); err != nil {
		return err
	}
	if err := users.ValidateName(membership.User); err != nil {
		return err
	}
	return nil
}
