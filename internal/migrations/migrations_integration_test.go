//+build integrationtest

package migrations

import (
	"context"
	"database/sql"
	"testing"

	"github.com/Saser/strecku/internal/config"
	"github.com/Saser/strecku/internal/inject"
	"github.com/Saser/strecku/internal/provide"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func pool(t *testing.T, logger *zap.Logger) (*sql.DB, func()) {
	ctx := context.Background()
	var cfg config.Config
	require.NoError(t, config.LoadFile("../../configs/integration_test.toml", &cfg))
	db, cleanupDB, err := inject.PostgresDBPoolFromConfig(ctx, logger, &cfg)
	require.NoError(t, err)
	return db, cleanupDB
}

func migrator(t *testing.T, logger *zap.Logger, db *sql.DB) *Migrator {
	m, err := NewMigrator(logger, db, "sql")
	require.NoError(t, err)
	return m
}

func TestNewMigrator(t *testing.T) {
	logger := provide.ZapTestLogger(t)
	db, cleanupDB := pool(t, logger)
	defer cleanupDB()
	_ = migrator(t, logger, db)
}

func TestMigrator_CheckVersion(t *testing.T) {
	logger := provide.ZapTestLogger(t)
	db, cleanupDB := pool(t, logger)
	defer cleanupDB()
	m := migrator(t, logger, db)
	tx, err := db.Begin()
	require.NoError(t, err)
	require.Error(t, m.CheckVersion())
	require.NoError(t, m.Perform())
	require.NoError(t, m.CheckVersion())
	require.NoError(t, tx.Rollback())
}

func XxxTestMigrator_Perform(t *testing.T) {
	logger := provide.ZapTestLogger(t)
	db, cleanupDB := pool(t, logger)
	defer cleanupDB()
	m := migrator(t, logger, db)
	tx, err := db.Begin()
	require.NoError(t, err)
	require.NoError(t, m.Perform())
	require.NoError(t, tx.Rollback())
}
