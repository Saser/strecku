package users

import (
	"testing"

	streckuv1 "github.com/Saser/strecku/saser/strecku/v1"
)

func TestValidate(t *testing.T) {
	for _, test := range []struct {
		user *streckuv1.User
		want error
	}{
		{user: &streckuv1.User{Name: "users/example", EmailAddress: "user@example.com", DisplayName: "User"}, want: nil},
		{user: &streckuv1.User{Name: "", EmailAddress: "user@example.com", DisplayName: "User"}, want: ErrNameEmpty},
		{user: &streckuv1.User{Name: "invalidprefix/example", EmailAddress: "user@example.com", DisplayName: "User"}, want: ErrNameInvalidPrefix},
		{user: &streckuv1.User{Name: "users/example", EmailAddress: "", DisplayName: "User"}, want: ErrEmailAddressEmpty},
		{user: &streckuv1.User{Name: "users/example", EmailAddress: "user@example.com", DisplayName: ""}, want: ErrDisplayNameEmpty},
	} {
		if got := Validate(test.user); got != test.want {
			t.Errorf("Validate(%v) = %v; want %v", test.user, got, test.want)
		}
	}
}
