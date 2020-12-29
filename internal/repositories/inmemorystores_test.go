package repositories

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestInMemoryStores(t *testing.T) {
	newStores := func() Stores { return NewInMemoryStores() }
	suite.Run(t, &StoresTestSuite{newStores: newStores})
}
