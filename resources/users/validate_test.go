package users

import (
	"testing"

	pb "github.com/Saser/strecku/api/v1"
	"github.com/Saser/strecku/resources/users/testusers"
)

func TestValidate(t *testing.T) {
	for _, test := range []struct {
		user *pb.User
		want error
	}{
		{
			user: &pb.User{
				Name:         testusers.Alice.Name,
				EmailAddress: "",
				DisplayName:  testusers.Alice.DisplayName,
			},
			want: ErrEmailAddressEmpty,
		},
		{
			user: &pb.User{
				Name:         testusers.Alice.Name,
				EmailAddress: testusers.Alice.EmailAddress,
				DisplayName:  "",
			},
			want: ErrDisplayNameEmpty,
		},
	} {
		if got := Validate(test.user); got != test.want {
			t.Errorf("Validate(%v) = %v; want %v", test.user, got, test.want)
		}
	}
}
