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
		{user: &streckuv1.User{Name: "users/6f2d193c-1460-491d-8157-7dd9535526c6", EmailAddress: "user@example.com", DisplayName: "User"}, want: nil},
		{user: &streckuv1.User{Name: "", EmailAddress: "user@example.com", DisplayName: "User"}, want: ErrNameEmpty},
		{user: &streckuv1.User{Name: "invalidprefix/6f2d193c-1460-491d-8157-7dd9535526c6", EmailAddress: "user@example.com", DisplayName: "User"}, want: ErrNameInvalidFormat},
		{user: &streckuv1.User{Name: "6f2d193c-1460-491d-8157-7dd9535526c6", EmailAddress: "user@example.com", DisplayName: "User"}, want: ErrNameInvalidFormat},
		{user: &streckuv1.User{Name: "users/6f2d193c-1460-491d-8157-7dd9535526c6", EmailAddress: "", DisplayName: "User"}, want: ErrEmailAddressEmpty},
		{user: &streckuv1.User{Name: "users/6f2d193c-1460-491d-8157-7dd9535526c6", EmailAddress: "user@example.com", DisplayName: ""}, want: ErrDisplayNameEmpty},
	} {
		if got := Validate(test.user); got != test.want {
			t.Errorf("Validate(%v) = %v; want %v", test.user, got, test.want)
		}
	}
}

func TestValidateName(t *testing.T) {
	for _, test := range []struct {
		name string
		want error
	}{
		{name: "users/6f2d193c-1460-491d-8157-7dd9535526c6", want: nil},
		{name: "", want: ErrNameEmpty},
		{name: "invalidprefix/6f2d193c-1460-491d-8157-7dd9535526c6", want: ErrNameInvalidFormat},
		{name: "6f2d193c-1460-491d-8157-7dd9535526c6", want: ErrNameInvalidFormat},
	} {
		if got := ValidateName(test.name); got != test.want {
			t.Errorf("ValidateName(%q) = %v; want %v", test.name, got, test.want)
		}
	}
}
