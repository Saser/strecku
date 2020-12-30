package database

import (
	"context"
	"database/sql"

	"github.com/cenkalti/backoff/v4"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func Open(ctx context.Context, connString string) (*sql.DB, error) {
	var db *sql.DB
	open := func() error {
		var err error
		db, err = sql.Open("pgx", connString)
		if err != nil {
			return err
		}
		return db.PingContext(ctx)
	}
	b := backoff.WithContext(backoff.NewExponentialBackOff(), ctx)
	if err := backoff.Retry(open, b); err != nil {
		return nil, err
	}
	return db, nil
}
