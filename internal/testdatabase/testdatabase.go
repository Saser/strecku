package testdatabase

import (
	"context"
	"fmt"
	"net/url"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/ory/dockertest/v3"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	Version      = 2
	User         = "testdatabase"
	Password     = "password"
	DatabaseName = "teststrecku"
)

func buildConnString(res *dockertest.Resource) string {
	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(User, Password),
		Host:     res.GetHostPort("5432/tcp"),
		Path:     DatabaseName,
		RawQuery: "sslmode=disable",
	}
	return u.String()
}

func newDockerPool(t *testing.T) *dockertest.Pool {
	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatal(err)
	}
	return pool
}

func newContainer(t *testing.T, dockerPool *dockertest.Pool) *dockertest.Resource {
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
	return res
}

func newDBPool(ctx context.Context, t *testing.T, dockerPool *dockertest.Pool, connString string) *pgxpool.Pool {
	var dbPool *pgxpool.Pool
	if err := dockerPool.Retry(func() error {
		var err error
		dbPool, err = pgxpool.Connect(ctx, connString)
		return err
	}); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(dbPool.Close)
	return dbPool
}

func newMigrate(t *testing.T, migrationsPath string, connString string) *migrate.Migrate {
	m, err := migrate.New("file://"+migrationsPath, connString)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		srcErr, dbErr := m.Close()
		if srcErr != nil {
			t.Error(srcErr)
		}
		if dbErr != nil {
			t.Error(dbErr)
		}
		if t.Failed() {
			t.FailNow()
		}
	})
	return m
}

func doMigrate(ctx context.Context, t *testing.T, m *migrate.Migrate) {
	go func() {
		<-ctx.Done()
		m.GracefulStop <- true
	}()
	if err := m.Migrate(Version); err != nil {
		t.Fatal(err)
	}
}

// New starts a Docker container running Postgres and returns a *pgxpool.Pool
// connected to that container. migrationsPath should be a path (absolute or
// relative) to the root `database` directory containing all database
// migrations. The returned pool does not need to be closed by the caller, it is
// handled automatically. New should be called rarely, as starting a container
// and connecting to it is a costly operation.
func New(ctx context.Context, t *testing.T, migrationsPath string) *pgxpool.Pool {
	dockerPool := newDockerPool(t)
	res := newContainer(t, dockerPool)
	connString := buildConnString(res)
	dbPool := newDBPool(ctx, t, dockerPool, connString)
	m := newMigrate(t, migrationsPath, connString)
	doMigrate(ctx, t, m)
	return dbPool
}
