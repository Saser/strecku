package users

func ValidatePassword(password string) error {
	if password == "" {
		return ErrPasswordEmpty
	}
	return nil
}
