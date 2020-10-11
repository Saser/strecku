package users

import (
	"testing"

	pb "github.com/Saser/strecku/api/v1"
	"github.com/Saser/strecku/resources/testresources"
)

func TestValidate(t *testing.T) {
	for _, test := range []struct {
		user *pb.User
		want error
	}{
		{
			user: &pb.User{
				Name:         testresources.Alice.Name,
				EmailAddress: "",
				DisplayName:  testresources.Alice.DisplayName,
			},
			want: ErrEmailAddressEmpty,
		},
		{
			user: &pb.User{
				Name:         testresources.Alice.Name,
				EmailAddress: testresources.Alice.EmailAddress,
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
