package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

const (
	valid                = "testdata/valid.toml"
	invalid_syntax       = "testdata/invalid_syntax.toml"
	invalid_logger_level = "testdata/invalid_logger_level.toml"
)

func TestLoadFile(t *testing.T) {
	var cfg Config
	t.Run("valid", func(t *testing.T) {
		require.NoError(t, LoadFile(valid, &cfg))
		assert.Equal(t, zap.NewAtomicLevelAt(zap.DebugLevel), cfg.Logger.Level)
	})
	t.Run("invalid", func(t *testing.T) {
		for _, file := range []string{
			invalid_syntax,
			invalid_logger_level,
		} {
			file := file
			t.Run(fmt.Sprintf("file=%v", file), func(t *testing.T) {
				assert.Error(t, LoadFile(file, &cfg))
			})
		}
	})
	t.Run("env_override", func(t *testing.T) {
		for key, value := range map[string]interface{}{
			"STRECKU_LOGGER_LEVEL": "info",
		} {
			require.NoError(t, os.Setenv(key, fmt.Sprint(value)))
		}
		require.NoError(t, LoadFile(valid, &cfg))
		assert.Equal(t, zap.NewAtomicLevelAt(zap.InfoLevel), cfg.Logger.Level)
	})
}
