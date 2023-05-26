package store

import (
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/go-mud/backend/core/config"
	"gitlab.com/alienspaces/go-mud/backend/core/log"
)

func newDependencies(t *testing.T) (*config.Config, *log.Log, *Store) {
	c, err := config.NewConfigWithDefaults([]config.Item{}, false)
	require.NoError(t, err, "NewConfig returns without error")

	l := log.NewLogger(c)

	s, err := NewStore(c, l)
	require.NoError(t, err, "NewStore returns without error")
	require.NotNil(t, s, "NewStore returns a store")

	return c, l, s
}

func TestGetDb(t *testing.T) {
	tests := map[string]struct {
		setup   func(*Store)
		wantErr bool
	}{
		"existing db connection": {
			setup: func(s *Store) {
				_, _ = s.GetDb()
			},
		},
		"new db connection": {
			setup: func(s *Store) {
				s.DB = nil
			},
			wantErr: false,
		},
		"new db connection with unsupported db type": {
			setup: func(s *Store) {
				s.Database = "abc"
			},
			wantErr: true,
		},
		"new db connection with wrong info": {
			setup: func(s *Store) {
				s.DB = nil
				s.Config.Set(config.AppServerDBHost, "a")
			},
			wantErr: true,
		},
	}

	for tcName, tc := range tests {

		t.Logf("Running test >%s<", tcName)

		func() {
			_, _, s := newDependencies(t)

			tc.setup(s)

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
			wantErr: false,
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
