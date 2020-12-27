package testdatabase

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/ory/dockertest/v3"
)

const (
	User         = "testdatabase"
	Password     = "password"
	DatabaseName = "teststrecku"
)

// New starts a Docker container running Postgres and returns a *pgxpool.Pool
// connected to that container. The returned pool does not need to be closed by
// the caller, it is handled automatically. New should be called rarely, as
// starting a container and connecting to it is a costly operation.
func New(ctx context.Context, t *testing.T) *pgxpool.Pool {
	dockerPool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatal(err)
	}
	res, err := dockerPool.Run("postgres", "13.1", []string{
		fmt.Sprintf("POSTGRES_USER=%s", User),
		fmt.Sprintf("POSTGRES_PASSWORD=%s", Password),
		fmt.Sprintf("POSTGRES_DB=%s", DatabaseName),
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := dockerPool.Purge(res); err != nil {
			t.Fatal(err)
		}
	})
	var dbPool *pgxpool.Pool
	if err := dockerPool.Retry(func() error {
		connStr := strings.Join([]string{
			"host=localhost",
			"sslmode=disable",
			fmt.Sprintf("port=%s", res.GetPort("5432/tcp")),
			fmt.Sprintf("dbname=%s", DatabaseName),
			fmt.Sprintf("user=%s", User),
			fmt.Sprintf("password=%s", Password),
		}, " ")
		var err error
		dbPool, err = pgxpool.Connect(ctx, connStr)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(dbPool.Close)
	return dbPool
}
