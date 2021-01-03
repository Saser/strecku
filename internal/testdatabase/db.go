package testdatabase

import (
	"context"
	"database/sql"
	"testing"

	"github.com/Saser/strecku/internal/database"
	"github.com/cenkalti/backoff/v4"
	"golang.org/x/sync/errgroup"
)

func retryConnString(ctx context.Context, tdb *TestDatabase) (string, error) {
	var connString string
	op := func() error {
		var err error
		connString, err = tdb.ConnString()
		return err
	}
	b := backoff.WithContext(backoff.NewExponentialBackOff(), ctx)
	if err := backoff.Retry(op, b); err != nil {
		return "", err
	}
	return connString, nil
}

func DB(ctx context.Context, t *testing.T, migrationsPath string) *sql.DB {
	t.Helper()
	tdb := New(migrationsPath)
	ctx, cancel := context.WithCancel(ctx)
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return tdb.Serve(ctx)
	})
	t.Cleanup(func() {
		cancel()
		if err := g.Wait(); err != nil {
			t.Errorf("g.Wait() = %v; want nil", err)
		}
	})
	connString, err := retryConnString(ctx, tdb)
	if err != nil {
		t.Fatalf("ConnString(ctx, tdb) err = %v; want nil", err)
	}
	db, err := database.Open(ctx, connString)
	if err != nil {
		t.Fatalf("database.Open(ctx, %q) err = %v; want nil", connString, err)
	}
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Errorf("db.Close() = %v; want nil", err)
		}
	})
	return db
}
