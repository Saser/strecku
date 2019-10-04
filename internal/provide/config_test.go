package provide

import (
	"fmt"
	"testing"

	"github.com/Saser/strecku/internal/config"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestLoggerLevelFromConfig(t *testing.T) {
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
			expected := zap.NewAtomicLevelAt(level)
			require.Equal(t, expected, LoggerLevelFromConfig(cfg))
		})
	}
}
