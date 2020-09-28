package memberships

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Saser/strecku/resources/stores"
	"github.com/Saser/strecku/resources/users"
	streckuv1 "github.com/Saser/strecku/saser/strecku/v1"
	"github.com/google/uuid"
)

var (
	ErrNameEmpty         = errors.New("name is empty")
	ErrNameInvalidFormat = fmt.Errorf("name must have format %q", prefix+"<uuid>")
)

func Validate(membership *streckuv1.Membership) error {
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
