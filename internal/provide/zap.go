package provide

import (
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
