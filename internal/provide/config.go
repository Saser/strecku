package provide

import (
	"time"

	"github.com/Saser/strecku/internal/config"
	"go.uber.org/zap"
)

func ConfigLoggerLevel(
	config *config.Config,
) zap.AtomicLevel {
	return config.Logger.Level
}

func ConfigDBConnString(
	config *config.Config,
) string {
	return config.DB.ConnString
}

func ConfigDBConnTimeout(
	config *config.Config,
) time.Duration {
	return config.DB.ConnTimeout
}
