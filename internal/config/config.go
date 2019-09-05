package config

import (
	"strings"
	"time"

	internalviper "github.com/Saser/strecku/internal/viper"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/xerrors"
)

const Prefix = "strecku"

type Config struct {
	Logger struct {
		Level zap.AtomicLevel
	}
	DB struct {
		ConnString  string
		ConnTimeout time.Duration
	}
}

func LoadFile(filePath string, cfg *Config) error {
	v := viper.New()
	v.SetConfigFile(filePath)
	v.SetEnvPrefix(Prefix)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		return xerrors.Errorf("load file: %w", err)
	}
	unmarshal := func(dc *mapstructure.DecoderConfig) {
		dc.ErrorUnused = true
	}
	hook := viper.DecodeHook(
		mapstructure.ComposeDecodeHookFunc(
			internalviper.ZapAtomicLevelDecodeHookFunc,
			mapstructure.StringToTimeDurationHookFunc(),
		),
	)
	if err := v.Unmarshal(cfg, unmarshal, hook); err != nil {
		return xerrors.Errorf("load file: %w", err)
	}
	return nil
}
