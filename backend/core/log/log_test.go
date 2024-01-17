package log

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/go-mud/backend/core/config"
	"gitlab.com/alienspaces/go-mud/backend/core/type/logger"
)

func TestNewLogger(t *testing.T) {

	envVars := map[string]string{
		config.AppServerLogLevel:  "debug",
		config.AppServerLogPretty: "true",
	}
	for key, val := range envVars {
		require.NoError(t, os.Setenv(key, val), "Set environment value")
	}

	c, err := config.NewConfig([]config.Item{
		{
			Key:      config.AppServerLogLevel,
			Required: true,
		},
		{
			Key:      config.AppServerLogPretty,
			Required: true,
		},
	}, false)
	require.NoError(t, err, "NewConfig returns without error")

	l, err := NewLogger(c)
	require.NoError(t, err, "NewLogger returns without error")
	require.NotNil(t, l, "NewLogger is not nil")

	l.Debug("Test level >%s<", "debug")
	l.Info("Test level >%s<", "info")
	l.Warn("Test level >%s<", "warn")
	l.Error("Test level >%s<", "error")

	l.WithContext("correlation-id", "abcdefg")

	l.Debug("Test level >%s<", "debug")

	l.WithContext("correlation-id", "hijklmn")

	l.Debug("Test level >%s<", "debug")

	l.Level(logger.ErrorLevel)

	l.Debug("Test level >%s<", "debug")
	l.Info("Test level >%s<", "info")
	l.Warn("Test level >%s<", "warn")
	l.Error("Test level >%s<", "error")
}

func TestNewLoggerWithConfig(t *testing.T) {

	l, err := NewLoggerWithConfig(Config{
		Level:  "debug",
		Pretty: true,
	})
	require.NoError(t, err, "NewLogger returns without error")
	require.NotNil(t, l, "NewLogger is not nil")

	l.Debug("Test level >%s<", "debug")
	l.Info("Test level >%s<", "info")
	l.Warn("Test level >%s<", "warn")
	l.Error("Test level >%s<", "error")

	l.WithContext("correlation-id", "abcdefg")

	l.Debug("Test level >%s<", "debug")

	l.WithContext("correlation-id", "hijklmn")

	l.Debug("Test level >%s<", "debug")

	l.Level(logger.ErrorLevel)

	l.Debug("Test level >%s<", "debug")
	l.Info("Test level >%s<", "info")
	l.Warn("Test level >%s<", "warn")
	l.Error("Test level >%s<", "error")
}
