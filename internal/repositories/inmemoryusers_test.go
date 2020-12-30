package repositories

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestInMemoryUsers(t *testing.T) {
	r := NewInMemoryUsers()
	suite.Run(t, NewUsersTestSuite(r))
}
