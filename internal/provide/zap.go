package provide

import (
	"log"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
	"golang.org/x/xerrors"
)

func ZapDevelopmentLogger(
	level zap.AtomicLevel,
) (*zap.Logger, error) {
	log.Print("creating zap development logger")
	cfg := zap.NewDevelopmentConfig()
	cfg.Level = level
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	options := []zap.Option{
		zap.AddCaller(),
	}
	logger, err := cfg.Build(options...)
	if err != nil {
		return nil, xerrors.Errorf("provide zap development logger: %w", err)
	}
	logger.Info("created zap development logger")
	return logger, nil
}

func ZapTestLogger(
	t *testing.T,
) *zap.Logger {
	t.Log("creating zap test logger")
	options := []zap.Option{
		zap.AddCaller(),
	}
	logger := zaptest.NewLogger(t, zaptest.WrapOptions(options...))
	logger.Info("created zap test logger")
	return logger
}
