package users

import (
	"errors"
	"fmt"
	"strings"

	streckuv1 "github.com/Saser/strecku/saser/strecku/v1"
	"github.com/google/uuid"
)

const prefix = CollectionID + "/"

var (
	ErrNameEmpty         = errors.New("name is empty")
	ErrNameInvalidFormat = fmt.Errorf("name must have format %q", prefix+"<uuid>")
	ErrEmailAddressEmpty = errors.New("email address is empty")
	ErrDisplayNameEmpty  = errors.New("display name is empty")
)

func Validate(user *streckuv1.User) error {
	if user.Name == "" {
		return ErrNameEmpty
	}
	if !strings.HasPrefix(user.Name, prefix) {
		return ErrNameInvalidFormat
	}
	if _, err := uuid.Parse(strings.TrimPrefix(user.Name, prefix)); err != nil {
		return ErrNameInvalidFormat
	}
	if user.EmailAddress == "" {
		return ErrEmailAddressEmpty
	}
	if user.DisplayName == "" {
		return ErrDisplayNameEmpty
	}
	return nil
}
