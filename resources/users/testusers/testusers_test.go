package testusers

import (
	"testing"

	"github.com/Saser/strecku/resources/users"
	streckuv1 "github.com/Saser/strecku/saser/strecku/v1"
)

func TestValid(t *testing.T) {
	for _, user := range []*streckuv1.User{
		Alice,
		Bob,
		Carol,
		David,
	} {
		if err := users.Validate(user); err != nil {
			t.Errorf("users.Validate(%v) = %v; want nil", user, err)
		}
	}
}
