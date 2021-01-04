package database

import (
	"context"
	"database/sql"
)

func InTx(ctx context.Context, db *sql.DB, f func(tx *sql.Tx) error) (err error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if rErr := tx.Rollback(); rErr != nil {
				err = rErr
			}
		} else {
			if cErr := tx.Commit(); cErr != nil {
				err = cErr
			}
		}
	}()
	return f(tx)
}
