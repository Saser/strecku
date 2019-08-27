package provide

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4/log/zapadapter"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
	"golang.org/x/xerrors"
)

func PGXPool(
	ctx context.Context,
	logger *zap.Logger,
	connString string,
	connTimeout time.Duration,
) (*pgxpool.Pool, func(), error) {
	logger.Info("creating pgx connection pool", zap.String("connString", connString))
	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, nil, xerrors.Errorf("provide pgx connection pool: %w", err)
	}
	poolConfig.ConnConfig.Logger = zapadapter.NewLogger(logger)
	connCtx, connCancel := context.WithTimeout(ctx, connTimeout)
	defer connCancel()
	pool, err := pgxpool.Connect(connCtx, connString)
	if err != nil {
		return nil, nil, xerrors.Errorf("provide pgx connection pool: %w", err)
	}
	logger.Info("created pgx connection pool")
	cleanup := func() {
		logger.Info("closing pgx connection pool")
		pool.Close()
		logger.Info("closed pgx connection pool")
	}
	return pool, cleanup, nil
}
