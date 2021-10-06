package server

import (
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/go-boilerplate/server/core/config"
	"gitlab.com/alienspaces/go-boilerplate/server/core/log"
	"gitlab.com/alienspaces/go-boilerplate/server/core/store"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/configurer"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/logger"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/storer"
)

// NewDefaultDependencies -
func NewDefaultDependencies() (configurer.Configurer, logger.Logger, storer.Storer, error) {

	// configurer
	c, err := config.NewConfig(nil, false)
	if err != nil {
		return nil, nil, nil, err
	}

	configVars := []string{
		// logger
		"APP_SERVER_LOG_LEVEL",
		// database
		"APP_SERVER_DB_HOST",
		"APP_SERVER_DB_PORT",
		"APP_SERVER_DB_NAME",
		"APP_SERVER_DB_USER",
		"APP_SERVER_DB_PASSWORD",
		// jwt signing key
		"APP_SERVER_JWT_SIGNING_KEY",
	}
	for _, key := range configVars {
		err = c.Add(key, true)
		if err != nil {
			return nil, nil, nil, err
		}
	}

	// logger
	l, err := log.NewLogger(c)
	if err != nil {
		return nil, nil, nil, err
	}

	// storer
	s, err := store.NewStore(c, l)
	if err != nil {
		return nil, nil, nil, err
	}

	return c, l, s, nil
}

func TestNewServer(t *testing.T) {

	c, l, s, err := NewDefaultDependencies()
	require.NoError(t, err, "NewDefaultDependencies returns without error")

	tr := Runner{}

	ts, err := NewServer(c, l, s, &tr)
	require.NoError(t, err, "NewServer returns without error")
	require.NotNil(t, ts, "Test server is not nil")
}
