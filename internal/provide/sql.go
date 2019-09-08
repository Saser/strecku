package provide

import (
	"context"
	"database/sql"
	"time"

	"github.com/cenkalti/backoff/v3"
	_ "github.com/lib/pq" // blank import needed for PostgreSQL driver
	"go.uber.org/zap"
	"golang.org/x/xerrors"
)

func PostgresDBPool(
	ctx context.Context,
	logger *zap.Logger,
	connString string,
	connTimeout time.Duration,
) (*sql.DB, func(), error) {
	logger.Info("creating new DB pool", zap.String("connString", connString))
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, nil, xerrors.Errorf("provide DB pool: %w", err)
	}
	pingCtx, pingCancel := context.WithTimeout(ctx, connTimeout)
	defer pingCancel()
	operation := func() error { return db.PingContext(pingCtx) }
	policy := backoff.WithContext(backoff.NewExponentialBackOff(), pingCtx)
	if err := backoff.Retry(operation, policy); err != nil {
		return nil, nil, xerrors.Errorf("provide DB pool: %w")
	}
	cleanup := func() {
		logger.Info("closing DB pool")
		if err := db.Close(); err != nil {
			logger.Error("closing DB pool failed", zap.Error(err))
			return
		}
		logger.Info("closed DB pool")
	}
	return db, cleanup, nil
}
