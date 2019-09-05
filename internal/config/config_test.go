package config

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

const (
	valid                = "testdata/valid.toml"
	invalidSyntax        = "testdata/invalid_syntax.toml"
	invalidLoggerLevel   = "testdata/invalid_logger_level.toml"
	invalidDBConnString  = "testdata/invalid_db_conn_string.toml"
	invalidDBConnTimeout = "testdata/invalid_db_conn_timeout.toml"
)

func TestLoadFile(t *testing.T) {
	var cfg Config
	t.Run("valid", func(t *testing.T) {
		require.NoError(t, LoadFile(valid, &cfg))
		assert.Equal(t, zap.NewAtomicLevelAt(zap.DebugLevel), cfg.Logger.Level)
		assert.Equal(t, "someConnString", cfg.DB.ConnString)
		assert.Equal(t, 10*time.Second, cfg.DB.ConnTimeout)
	})
	t.Run("invalid", func(t *testing.T) {
		for _, file := range []string{
			invalidSyntax,
			invalidLoggerLevel,
			invalidDBConnString,
			invalidDBConnTimeout,
		} {
			file := file
			t.Run(fmt.Sprintf("file=%v", file), func(t *testing.T) {
				assert.Error(t, LoadFile(file, &cfg))
			})
		}
	})
	t.Run("env_override", func(t *testing.T) {
		for key, value := range map[string]interface{}{
			"STRECKU_LOGGER_LEVEL":   "info",
			"STRECKU_DB_CONNSTRING":  "someOtherConnString",
			"STRECKU_DB_CONNTIMEOUT": "15s",
		} {
			require.NoError(t, os.Setenv(key, fmt.Sprint(value)))
		}
		require.NoError(t, LoadFile(valid, &cfg))
		assert.Equal(t, zap.NewAtomicLevelAt(zap.InfoLevel), cfg.Logger.Level)
		assert.Equal(t, "someOtherConnString", cfg.DB.ConnString)
		assert.Equal(t, 15*time.Second, cfg.DB.ConnTimeout)
	})
}
