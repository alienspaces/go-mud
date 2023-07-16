package store

import (
	"testing"

	"github.com/stretchr/testify/require"

	coreconfig "gitlab.com/alienspaces/go-mud/backend/core/config"
	"gitlab.com/alienspaces/go-mud/backend/core/log"
)

func Test_getPostgresDB(t *testing.T) {
	c, err := coreconfig.NewConfigWithDefaults(nil, false)
	require.NoError(t, err)

	l, err := log.NewLogger(c)
	require.NoError(t, err, "NewLogger returns without error")

	s, err := NewStore(c, l)
	require.NoError(t, err, "NewStore returns without error")

	defaultConnectionConfig, err := s.GetConnectionConfig()
	require.NoError(t, err)

	tests := map[string]struct {
		config  func() *Config
		wantErr bool
	}{
		"without APP_SERVER_DB_HOST": {
			config: func() *Config {
				c := *defaultConnectionConfig
				c.Host = ""
				return &c
			},
			wantErr: true,
		},
		"without APP_SERVER_DB_PORT": {
			config: func() *Config {
				c := *defaultConnectionConfig
				c.Port = ""
				return &c
			},
			wantErr: true,
		},
		"without APP_SERVER_DB_NAME": {
			config: func() *Config {
				c := *defaultConnectionConfig
				c.Database = ""
				return &c
			},
			wantErr: true,
		},
		"without APP_SERVER_DB_USER": {
			config: func() *Config {
				c := *defaultConnectionConfig
				c.User = ""
				return &c
			},
			wantErr: true,
		},
		"without APP_SERVER_DB_PASSWORD": {
			config: func() *Config {
				c := *defaultConnectionConfig
				c.Password = ""
				return &c
			},
			wantErr: true,
		},
		"missing config": {
			config: func() *Config {
				return nil
			},
			wantErr: true,
		},
		"valid config": {
			config: func() *Config {
				c := *defaultConnectionConfig
				return &c
			},
			wantErr: false,
		},
	}

	for tcName, tc := range tests {

		t.Logf("Running test >%s<", tcName)

		func() {
			config := tc.config()

			db, err := connectPostgresDB(l, config)
			defer func() {
				if db != nil {
					db.Close()
				}
			}()
			if tc.wantErr {
				require.Error(t, err, "connectPostgresDB returns with error")
				return
			}

			require.NoError(t, err, "connectPostgresDB returns without error")
			require.NotNil(t, db, "connectPostgresDB returns db connection")
		}()
	}
}
