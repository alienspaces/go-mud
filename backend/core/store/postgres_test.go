package store

import (
	"testing"

	"github.com/stretchr/testify/require"

	coreconfig "gitlab.com/alienspaces/go-mud/backend/core/config"
	"gitlab.com/alienspaces/go-mud/backend/core/log"
	"gitlab.com/alienspaces/go-mud/backend/core/type/storer"
)

func Test_getPostgresDB(t *testing.T) {
	// setThenRestoreEnv := func(key string, value string) func(*coreconfig.Config) func() {
	// 	return func(conf *coreconfig.Config) func() {
	// 		oldEnvValue := conf.Get(key)
	// 		conf.Set(key, value)

	// 		return func() {
	// 			conf.Set(key, oldEnvValue)
	// 		}
	// 	}
	// }

	c, err := coreconfig.NewConfigWithDefaults(nil, false)
	require.NoError(t, err)

	l := log.NewLogger(c)

	s, err := NewStore(c, l)
	require.NoError(t, err)

	defaultConnectionConfig, err := s.GetConnectionConfig()
	require.NoError(t, err)

	tests := map[string]struct {
		connectionConfig func() *storer.ConnectionConfig
		wantErr          bool
	}{
		"without APP_SERVER_DB_HOST": {
			connectionConfig: func() *storer.ConnectionConfig {
				c := *defaultConnectionConfig
				c.Host = ""
				return &c
			},
			wantErr: true,
		},
		"without APP_SERVER_DB_PORT": {
			connectionConfig: func() *storer.ConnectionConfig {
				c := *defaultConnectionConfig
				c.Port = ""
				return &c
			},
			wantErr: true,
		},
		"without APP_SERVER_DB_NAME": {
			connectionConfig: func() *storer.ConnectionConfig {
				c := *defaultConnectionConfig
				c.Database = ""
				return &c
			},
			wantErr: true,
		},
		"without APP_SERVER_DB_USER": {
			connectionConfig: func() *storer.ConnectionConfig {
				c := *defaultConnectionConfig
				c.User = ""
				return &c
			},
			wantErr: true,
		},
		"without APP_SERVER_DB_PASSWORD": {
			connectionConfig: func() *storer.ConnectionConfig {
				c := *defaultConnectionConfig
				c.Password = ""
				return &c
			},
			wantErr: true,
		},
		"missing config": {
			connectionConfig: func() *storer.ConnectionConfig {
				return nil
			},
			wantErr: true,
		},
		"valid config": {
			connectionConfig: func() *storer.ConnectionConfig {
				c := *defaultConnectionConfig
				return &c
			},
			wantErr: false,
		},
	}

	for tcName, tc := range tests {

		t.Logf("Running test >%s<", tcName)

		func() {
			connectionConfig := tc.connectionConfig()

			db, err := connectPostgresDB(l, connectionConfig)
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
