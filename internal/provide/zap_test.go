package provide

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestZapDevelopmentLogger(t *testing.T) {
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
			_, err := ZapDevelopmentLogger(zap.NewAtomicLevelAt(level))
			assert.NoError(t, err)
		})
	}
}

func TestZapTestLogger(t *testing.T) {
	_ = ZapTestLogger(t)
}
