package repositories

import (
	"context"
	"testing"

	"github.com/Saser/strecku/internal/testdatabase"
	"github.com/stretchr/testify/suite"
)

func TestPostgresUsers(t *testing.T) {
	if testing.Short() {
		t.Skipf("skipping: -short is set")
	}
	ctx := context.Background()
	db := testdatabase.DB(ctx, t, "../../database")
	r := NewPostgresUsers(db)
	suite.Run(t, NewUsersTestSuite(r))
}
