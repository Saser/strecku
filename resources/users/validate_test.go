package users

import (
	"testing"

	pb "github.com/Saser/strecku/api/v1"
)

func TestValidate(t *testing.T) {
	for _, test := range []struct {
		user *pb.User
		want error
	}{
		{user: &pb.User{Name: "users/6f2d193c-1460-491d-8157-7dd9535526c6", EmailAddress: "user@example.com", DisplayName: "User"}, want: nil},
		{user: &pb.User{Name: "", EmailAddress: "user@example.com", DisplayName: "User"}, want: ErrNameEmpty},
		{user: &pb.User{Name: "invalidprefix/6f2d193c-1460-491d-8157-7dd9535526c6", EmailAddress: "user@example.com", DisplayName: "User"}, want: ErrNameInvalidFormat},
		{user: &pb.User{Name: "6f2d193c-1460-491d-8157-7dd9535526c6", EmailAddress: "user@example.com", DisplayName: "User"}, want: ErrNameInvalidFormat},
		{user: &pb.User{Name: "users/6f2d193c-1460-491d-8157-7dd9535526c6", EmailAddress: "", DisplayName: "User"}, want: ErrEmailAddressEmpty},
		{user: &pb.User{Name: "users/6f2d193c-1460-491d-8157-7dd9535526c6", EmailAddress: "user@example.com", DisplayName: ""}, want: ErrDisplayNameEmpty},
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
		{name: "users/not a UUID", want: ErrNameInvalidFormat},
	} {
		if got := ValidateName(test.name); got != test.want {
			t.Errorf("ValidateName(%q) = %v; want %v", test.name, got, test.want)
		}
	}
}
