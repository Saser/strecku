package repositories

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestInMemoryUsers(t *testing.T) {
	newUsers := func() Users { return NewInMemoryUsers() }
	suite.Run(t, &UsersTestSuite{newUsers: newUsers})
}
