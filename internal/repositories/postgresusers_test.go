package repositories

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/Saser/strecku/internal/testdatabase"
	"github.com/stretchr/testify/suite"
)

func TestPostgresUsers(t *testing.T) {
	ctx := context.Background()
	pool := testdatabase.New(ctx, t, "../../database")
	newUsers := func() Users {
		tx, err := pool.Begin(ctx)
		if err != nil {
			t.Fatal(err)
		}
		defer func() {
			if err := tx.Commit(ctx); err != nil {
				t.Fatal(err)
			}
		}()
		tables := []string{
			"users",
			"stores",
		}
		sql := fmt.Sprintf("TRUNCATE %s CASCADE;", strings.Join(tables, ", "))
		if _, err := tx.Exec(ctx, sql); err != nil {
			t.Fatal(err)
		}
		return NewPostgresUsers(pool)
	}
	suite.Run(t, &UsersTestSuite{newUsers: newUsers})
}
