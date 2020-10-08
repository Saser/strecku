package memberships

import (
	"errors"
	"fmt"

	pb "github.com/Saser/strecku/api/v1"
	"github.com/Saser/strecku/resources/stores"
	"github.com/Saser/strecku/resources/users"
)

var (
	ErrNameEmpty         = errors.New("name is empty")
	ErrNameInvalidFormat = fmt.Errorf("name must have format %q", prefix+"<uuid>")
)

func Validate(membership *pb.Membership) error {
	if err := ValidateName(membership.Name); err != nil {
		return err
	}
	if err := users.ValidateName(membership.User); err != nil {
		return err
	}
	if err := stores.ValidateName(membership.Store); err != nil {
		return err
	}
	return nil
}
