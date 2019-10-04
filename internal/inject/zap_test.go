package inject

import (
	"fmt"
	"testing"

	"github.com/Saser/strecku/internal/config"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestZapDevelopmentLoggerFromConfig(t *testing.T) {
	for _, level := range []zapcore.Level{
		zapcore.DebugLevel,
		zapcore.InfoLevel,
		zapcore.WarnLevel,
		zapcore.ErrorLevel,
		zapcore.DPanicLevel,
		zapcore.PanicLevel,
		zapcore.FatalLevel,
	} {
		level := level
		t.Run(fmt.Sprintf("level=%v", level), func(t *testing.T) {
			cfg := &config.Config{
				Logger: struct {
					Level zap.AtomicLevel
				}{
					Level: zap.NewAtomicLevelAt(level),
				},
			}
			_, cleanup, err := ZapDevelopmentLoggerFromConfig(cfg)
			require.NoError(t, err)
			defer cleanup()
		})
	}
}
