package server

import (
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/go-mud/server/core/config"
	"gitlab.com/alienspaces/go-mud/server/core/log"
	"gitlab.com/alienspaces/go-mud/server/core/store"
	"gitlab.com/alienspaces/go-mud/server/core/type/configurer"
	"gitlab.com/alienspaces/go-mud/server/core/type/logger"
	"gitlab.com/alienspaces/go-mud/server/core/type/storer"
)

// NewDefaultDependencies -
func NewDefaultDependencies() (configurer.Configurer, logger.Logger, storer.Storer, error) {

	// configurer
	c, err := config.NewConfigWithDefaults(nil, false)
	if err != nil {
		return nil, nil, nil, err
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

	defer func() {
		db, err := s.GetDb()
		require.NoError(t, err, "getDb should return no error")

		err = db.Close()
		require.NoError(t, err, "close db should return no error")
	}()

	tr := Runner{}
	tr.Config = c
	tr.Log = l

	ts, err := NewServer(c, l, s, &tr)
	require.NoError(t, err, "NewServer returns without error")
	require.NotNil(t, ts, "Test server is not nil")
}
