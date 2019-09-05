//+build integrationtest

package provide

import (
	"context"
	"testing"

	"github.com/Saser/strecku/internal/config"
	"github.com/stretchr/testify/require"
)

func TestPostgresDBPool(t *testing.T) {
	ctx := context.Background()
	logger := ZapTestLogger(t)
	var cfg config.Config
	require.NoError(t, config.LoadFile("../../configs/integration_test.toml", &cfg))
	db, cleanupDB, err := PostgresDBPool(ctx, logger, cfg.DB.ConnString, cfg.DB.ConnTimeout)
	require.NoError(t, err)
	defer cleanupDB()
	require.NoError(t, db.Ping())
}
