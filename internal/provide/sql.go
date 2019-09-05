package provide

import (
	"context"
	"database/sql"
	"time"

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
	var pingErr error
	pingCtxDone := pingCtx.Done()
	ticker := time.NewTicker(1 * time.Second)
loop:
	for {
		select {
		case <-pingCtxDone:
			pingErr = pingCtx.Err()
			break loop
		case <-ticker.C:
			if err := db.PingContext(pingCtx); err != nil {
				logger.Warn("pinging DB failed, retrying", zap.Error(err))
				continue
			}
			pingErr = nil
			break loop
		}
	}
	if pingErr != nil {
		return nil, nil, xerrors.Errorf("provide DB pool: %w", err)
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
