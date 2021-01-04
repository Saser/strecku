package testdatabase

import (
	"context"
	"errors"
	"flag"
	"net/url"
	"sync"

	"github.com/Saser/strecku/internal/database"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/ory/dockertest/v3"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	version  = 2
	user     = "strecku"
	password = "password"
	dbName   = "strecku"
)

var (
	ErrNotReady = errors.New("testdatabase: not ready")

	deleteContainer = flag.Bool("testdatabase_delete_container", true, "Whether the database container should be deleted during cleanup.")
)

type TestDatabase struct {
	migrationsPath string

	mu         sync.RWMutex
	connString string
}

func New(migrationsPath string) *TestDatabase {
	return &TestDatabase{
		migrationsPath: migrationsPath,
	}
}

func (t *TestDatabase) Serve(ctx context.Context) (err error) {
	// Whenever and for whatever reason this method returns, we are not
	// ready to serve a connection string.
	defer func() {
		t.mu.Lock()
		defer t.mu.Unlock()
		t.connString = ""
	}()

	pool, err := dockertest.NewPool("")
	if err != nil {
		return err
	}

	container, err := pool.Run("postgres", "13.1", []string{
		"POSTGRES_USER=" + user,
		"POSTGRES_PASSWORD=" + password,
		"POSTGRES_DB=" + dbName,
	})
	if err != nil {
		return err
	}
	if *deleteContainer {
		defer func() {
			if cErr := container.Close(); cErr != nil && err == nil {
				err = cErr
			}
		}()
	}

	u := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(user, password),
		Host:   container.GetHostPort("5432/tcp"),
		Path:   dbName,
	}
	q := url.Values{}
	q.Add("sslmode", "disable")
	u.RawQuery = q.Encode()
	connString := u.String()

	db, err := database.Open(ctx, connString)
	if err != nil {
		return err
	}
	defer func() {
		if cErr := db.Close(); cErr != nil && err == nil {
			err = cErr
		}
	}()

	srcInstance, err := source.Open("file://" + t.migrationsPath)
	if err != nil {
		return err
	}
	defer func() {
		if cErr := srcInstance.Close(); cErr != nil && err == nil {
			err = cErr
		}
	}()

	dbInstance, err := postgres.WithInstance(db, new(postgres.Config))
	if err != nil {
		return err
	}
	defer func() {
		if cErr := dbInstance.Close(); cErr != nil && err == nil {
			err = cErr
		}
	}()

	m, err := migrate.NewWithInstance("file", srcInstance, "postgres", dbInstance)
	if err != nil {
		return err
	}
	go func() {
		<-ctx.Done()
		m.GracefulStop <- true
	}()
	if err := m.Migrate(version); err != nil {
		return err
	}

	t.mu.Lock()
	t.connString = connString
	t.mu.Unlock()

	<-ctx.Done()

	return nil
}

func (t *TestDatabase) ConnString() (string, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	if t.connString == "" {
		return "", ErrNotReady
	}
	return t.connString, nil
}
