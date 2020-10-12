package users

func ValidatePassword(password string) error {
	if password == "" {
		return ErrEmptyPassword
	}
	return nil
}
