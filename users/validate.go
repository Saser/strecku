package users

import (
	"errors"
	"fmt"
	"strings"

	streckuv1 "github.com/Saser/strecku/saser/strecku/v1"
)

var (
	ErrNameEmpty         = errors.New("name is empty")
	ErrNameInvalidPrefix = fmt.Errorf("name must have prefix %q", CollectionID+"/")
	ErrEmailAddressEmpty = errors.New("email address is empty")
	ErrDisplayNameEmpty  = errors.New("display name is empty")
)

func Validate(user *streckuv1.User) error {
	if user.Name == "" {
		return ErrNameEmpty
	}
	if !strings.HasPrefix(user.Name, CollectionID+"/") {
		return ErrNameInvalidPrefix
	}
	if user.EmailAddress == "" {
		return ErrEmailAddressEmpty
	}
	if user.DisplayName == "" {
		return ErrDisplayNameEmpty
	}
	return nil
}
