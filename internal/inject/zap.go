//+build wireinject

package inject

import (
	"github.com/Saser/strecku/internal/config"
	"github.com/Saser/strecku/internal/provide"
	"github.com/google/wire"
	"go.uber.org/zap"
)

func ZapDevelopmentLoggerFromConfig(
	config *config.Config,
) (*zap.Logger, func(), error) {
	panic( // panic to be replaced by generated code
		wire.Build(
			provide.ConfigLoggerLevel,
			provide.ZapDevelopmentLogger,
		),
	)
}
