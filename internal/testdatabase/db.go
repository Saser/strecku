package testdatabase

import (
	"context"
	"database/sql"
	"testing"

	"github.com/Saser/strecku/internal/database"
	"github.com/golang-migrate/migrate/v4"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const version = 2

func DB(ctx context.Context, t *testing.T, migrationsPath string) *sql.DB {
	mu.Lock()
	defer mu.Unlock()
	check()

	connString := defaultContainer.ConnString()

	// Open a database connection to the container. This also makes sure
	// that the database is up and running.
	db, err := database.Open(ctx, connString)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Fatal(err)
		}
	})

	// Create a new migration runner, which will create its own database
	// connection.
	m, err := migrate.New("file://"+migrationsPath, connString)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		srcErr, dbErr := m.Close()
		if srcErr != nil {
			t.Fatal(srcErr)
		}
		if dbErr != nil {
			t.Fatal(dbErr)
		}
	}()

	// Make sure that when the context is done, the migration is gracefully
	// stopped.
	go func() {
		<-ctx.Done()
		m.GracefulStop <- true
	}()

	// Finally, run the migrations.
	if err := m.Migrate(version); err != nil {
		t.Fatal(err)
	}

	return db
}
