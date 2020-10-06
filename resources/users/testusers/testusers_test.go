package testusers

import (
	"testing"

	pb "github.com/Saser/strecku/api/v1"
	"github.com/Saser/strecku/resources/users"
)

func TestValid(t *testing.T) {
	for _, user := range []*pb.User{
		Alice,
		Bob,
		Carol,
	} {
		if err := users.Validate(user); err != nil {
			t.Errorf("users.Validate(%v) = %v; want nil", user, err)
		}
	}
}
