package store

import (
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/go-mud/server/core/config"
	"gitlab.com/alienspaces/go-mud/server/core/log"
)

func newDependencies(t *testing.T) (*config.Config, *log.Log, *Store) {
	c, err := config.NewConfigWithDefaults([]config.Item{}, false)
	require.NoError(t, err, "NewConfig returns without error")

	l, err := log.NewLogger(c)
	require.NoError(t, err, "NewLogger returns without error")

	s, err := NewStore(c, l)
	require.NoError(t, err, "NewStore returns without error")
	require.NotNil(t, s, "NewStore returns a store")

	return c, l, s
}

func TestInit(t *testing.T) {
	tests := map[string]struct {
		setup   func(*Store) func()
		wantErr bool
	}{
		"ok": {},
		"new db connection with unsupported db type": {
			setup: func(s *Store) func() {
				oldDb := s.Database
				oldDbConn := s.DB

				s.DB = nil
				s.Database = "abc"

				return func() {
					s.DB = oldDbConn
					s.Database = oldDb
				}
			},
			wantErr: true,
		},
	}

	_, _, s := newDependencies(t)

	for tcName, tc := range tests {

		t.Logf("Running test >%s<", tcName)

		func() {

			if tc.setup != nil {
				teardown := tc.setup(s)
				defer func() {
					teardown()
				}()
			}

			err := s.Init()
			if tc.wantErr {
				require.Error(t, err, "Init returns with error")
			} else {
				require.NoError(t, err, "Init returns without error")
			}
		}()
	}
}

func TestGetDb(t *testing.T) {
	tests := map[string]struct {
		setup   func(*Store) func()
		wantErr bool
	}{
		"existing db connection": {
			setup: func(s *Store) func() {
				_, _ = s.GetDb()

				return func() {
					s.DB = nil
				}
			},
		},
		"new db connection": {
			setup: func(s *Store) func() {
				oldDbConn := s.DB

				s.DB = nil

				return func() {
					s.DB = oldDbConn
				}
			},
			wantErr: false,
		},
		"new db connection with unsupported db type": {
			setup: func(s *Store) func() {
				oldDb := s.Database
				oldDbConn := s.DB

				s.DB = nil
				s.Database = "abc"

				return func() {
					s.DB = oldDbConn
					s.Database = oldDb
				}
			},
			wantErr: true,
		},
		"new db connection with wrong info": {
			setup: func(s *Store) func() {
				oldEnvValue := s.Config.Get("APP_SERVER_DB_HOST")
				oldDbConn := s.DB

				s.DB = nil
				s.Config.Set("APP_SERVER_DB_HOST", "a")

				return func() {
					s.DB = oldDbConn
					s.Config.Set("APP_SERVER_DB_HOST", oldEnvValue)
				}
			},
			wantErr: true,
		},
	}

	_, _, s := newDependencies(t)

	for tcName, tc := range tests {

		t.Logf("Running test >%s<", tcName)

		func() {
			if tc.setup != nil {
				teardown := tc.setup(s)
				defer func() {
					teardown()
				}()
			}

			db, err := s.GetDb()
			if tc.wantErr {
				require.Error(t, err, "GetDb returns with error")
				return
			}

			require.NoError(t, err, "GetDb returns without error")
			require.NotNil(t, db, "GetDb returns db struct")
		}()
	}
}

func TestGetTx(t *testing.T) {
	tests := map[string]struct {
		setup   func(*Store) func()
		wantErr bool
	}{
		"existing db connection": {},
		"no db connection": {
			setup: func(s *Store) func() {
				oldDbConn := s.DB

				s.DB = nil

				return func() {
					s.DB = oldDbConn
				}
			},
			wantErr: true,
		},
	}

	_, _, s := newDependencies(t)
	err := s.Init()
	require.NoError(t, err, "Init returns without error")

	for tcName, tc := range tests {

		t.Logf("Running test >%s<", tcName)

		func() {
			if tc.setup != nil {
				teardown := tc.setup(s)
				defer func() {
					teardown()
				}()
			}

			tx, err := s.GetTx()
			if tc.wantErr {
				require.Error(t, err, "GetTx returns with error")
				return
			}

			require.NoError(t, err, "GetTx returns without error")
			require.NotNil(t, tx, "GetTx returns tx struct")
		}()
	}
}
