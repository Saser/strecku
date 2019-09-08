package migrations

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // blank import needed for `file` driver
	"go.uber.org/zap"
	"golang.org/x/xerrors"
)

const (
	Version = 4
)

type Migrator struct {
	logger  *zap.Logger
	migrate *migrate.Migrate
}

func NewMigrator(
	logger *zap.Logger,
	db *sql.DB,
	path string,
) (*Migrator, error) {
	logger.Info("creating new migrator", zap.String("path", path))
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, xerrors.Errorf("new migrator: %w", err)
	}
	m, err := migrate.NewWithDatabaseInstance(fmt.Sprintf("file://%v", path), "postgres", driver)
	if err != nil {
		return nil, xerrors.Errorf("new migrator: %w", err)
	}
	migrator := &Migrator{
		logger:  logger,
		migrate: m,
	}
	logger.Info("created new migrator")
	return migrator, nil
}

func (m *Migrator) CheckVersion() error {
	m.logger.Debug("checking version", zap.Int("expected", Version))
	current, dirty, err := m.migrate.Version()
	if err != nil {
		return xerrors.Errorf("migrations: check current: %w", err)
	}
	m.logger.Debug("got version", zap.Bool("dirty", dirty), zap.Uint("current", current))
	if dirty {
		return xerrors.New("migrations: check current: current is dirty")
	}
	if current != Version {
		return xerrors.Errorf("migrations: check current: expecting current %v, got %v", Version, current)
	}
	return nil
}

func (m *Migrator) Perform() error {
	if err := m.migrate.Up(); err != nil {
		return xerrors.Errorf("migrations: perform: %w", err)
	}
	return nil
}
