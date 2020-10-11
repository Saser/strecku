package users

import (
	"errors"

	pb "github.com/Saser/strecku/api/v1"
)

var (
	ErrEmailAddressEmpty = errors.New("email address is empty")
	ErrDisplayNameEmpty  = errors.New("display name is empty")
)

func Validate(user *pb.User) error {
	if err := ValidateName(user.Name); err != nil {
		return err
	}
	if user.EmailAddress == "" {
		return ErrEmailAddressEmpty
	}
	if user.DisplayName == "" {
		return ErrDisplayNameEmpty
	}
	return nil
}
