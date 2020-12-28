package users

import "errors"

var (
	ErrPasswordEmpty = errors.New("empty password")
)

func ValidatePassword(password string) error {
	if password == "" {
		return ErrPasswordEmpty
	}
	return nil
}
