package provide

import (
	"github.com/Saser/strecku/internal/config"
	"go.uber.org/zap"
)

func ConfigLoggerLevel(
	config *config.Config,
) zap.AtomicLevel {
	return config.Logger.Level
}
