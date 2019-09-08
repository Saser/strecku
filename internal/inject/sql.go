//+build wireinject

package inject

import (
	"context"
	"database/sql"

	"github.com/Saser/strecku/internal/config"
	"github.com/Saser/strecku/internal/provide"
	"github.com/google/wire"
	"go.uber.org/zap"
)

func PostgresDBPoolFromConfig(
	ctx context.Context,
	logger *zap.Logger,
	config *config.Config,
) (*sql.DB, func(), error) {
	panic( // panic to be replaced by generated code
		wire.Build(
			provide.ConfigDBConnString,
			provide.ConfigDBConnTimeout,
			provide.PostgresDBPool,
		),
	)
}
