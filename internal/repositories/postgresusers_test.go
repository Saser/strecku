package repositories

import (
	"context"
	"testing"

	"github.com/Saser/strecku/internal/testdatabase"
	"github.com/stretchr/testify/suite"
)

func TestPostgresUsers(t *testing.T) {
	ctx := context.Background()
	testdatabase.Init()
	t.Cleanup(testdatabase.Cleanup)
	db := testdatabase.DB(ctx, t, "../../database")
	r := NewPostgresUsers(db)
	suite.Run(t, NewUsersTestSuite(r))
}