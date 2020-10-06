package users

import (
	"errors"
	"fmt"
	"strings"

	pb "github.com/Saser/strecku/api/v1"
	"github.com/google/uuid"
)

var (
	ErrNameEmpty         = errors.New("name is empty")
	ErrNameInvalidFormat = fmt.Errorf("name must have format %q", prefix+"<uuid>")
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

func ValidateName(name string) error {
	if name == "" {
		return ErrNameEmpty
	}
	if !strings.HasPrefix(name, prefix) {
		return ErrNameInvalidFormat
	}
	if _, err := uuid.Parse(strings.TrimPrefix(name, prefix)); err != nil {
		return ErrNameInvalidFormat
	}
	return nil
}
