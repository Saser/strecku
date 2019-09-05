package provide

import (
	"github.com/Saser/strecku/internal/config"
	"go.uber.org/zap"
)

func LoggerLevelFromConfig(
	config *config.Config,
) zap.AtomicLevel {
	return config.Logger.Level
}
