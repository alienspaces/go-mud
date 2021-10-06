package store

import (
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/go-boilerplate/server/core/config"
	"gitlab.com/alienspaces/go-boilerplate/server/core/log"
)

func TestNewStore(t *testing.T) {

	// configurer
	c, err := config.NewConfig([]config.Item{}, false)
	require.NoError(t, err, "NewConfig returns without error")

	configVars := []string{
		// database
		"APP_SERVER_DB_HOST",
		"APP_SERVER_DB_PORT",
		"APP_SERVER_DB_NAME",
		"APP_SERVER_DB_USER",
		"APP_SERVER_DB_PASSWORD",
	}
	for _, key := range configVars {
		require.NoError(t, c.Add(key, true), "Add config item")
	}

	// logger
	l, err := log.NewLogger(c)
	require.NoError(t, err, "NewLogger returns without error")

	// storer
	s, err := NewStore(c, l)
	require.NoError(t, err, "NewStore returns without error")
	require.NotNil(t, s, "NewStore returns a store")
}
