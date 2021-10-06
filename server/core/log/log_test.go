package log

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/go-boilerplate/server/core/config"
)

func TestLogger(t *testing.T) {

	// config
	c, err := config.NewConfig([]config.Item{}, false)
	require.Nil(t, err, "NewConfig returns without error")

	envVars := map[string]string{
		// logger
		"APP_SERVER_LOG_LEVEL": "debug",
	}
	for key, val := range envVars {
		require.NoError(t, os.Setenv(key, val), "Set environment value")
	}

	l, err := NewLogger(c)
	require.NoError(t, err, "NewLogger returns without error")
	require.NotNil(t, l, "NewLogger is not nil")

	l.Debug("Test level >%s<", "debug")
	l.Info("Test level >%s<", "info")
	l.Warn("Test level >%s<", "warn")
	l.Error("Test level >%s<", "error")

	l.Context("correlation-id", "abcdefg")

	l.Debug("Test level >%s<", "debug")

	l.Context("correlation-id", "hijklmn")

	l.Debug("Test level >%s<", "debug")

	l.Level(ErrorLevel)

	l.Debug("Test level >%s<", "debug")
	l.Info("Test level >%s<", "info")
	l.Warn("Test level >%s<", "warn")
	l.Error("Test level >%s<", "error")
}
