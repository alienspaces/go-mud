package store

import (
	"testing"

	"github.com/stretchr/testify/require"

	coreConfig "gitlab.com/alienspaces/go-mud/backend/core/config"
	"gitlab.com/alienspaces/go-mud/backend/core/log"
)

func Test_getPostgresDB(t *testing.T) {
	setThenRestoreEnv := func(key string, value string) func(*coreConfig.Config) func() {
		return func(conf *coreConfig.Config) func() {
			oldEnvValue := conf.Get(key)
			conf.Set(key, value)

			return func() {
				conf.Set(key, oldEnvValue)
			}
		}
	}

	tests := map[string]struct {
		setup   func(*coreConfig.Config) func()
		wantErr bool
	}{
		"without APP_SERVER_DB_HOST": {
			setup:   setThenRestoreEnv("APP_SERVER_DB_HOST", ""),
			wantErr: true,
		},
		"without APP_SERVER_DB_PORT": {
			setup:   setThenRestoreEnv("APP_SERVER_DB_PORT", ""),
			wantErr: true,
		},
		"without APP_SERVER_DB_NAME": {
			setup:   setThenRestoreEnv("APP_SERVER_DB_NAME", ""),
			wantErr: true,
		},
		"without APP_SERVER_DB_USER": {
			setup:   setThenRestoreEnv("APP_SERVER_DB_USER", ""),
			wantErr: true,
		},
		"without APP_SERVER_DB_PASSWORD": {
			setup:   setThenRestoreEnv("APP_SERVER_DB_PASSWORD", ""),
			wantErr: true,
		},
		"invalid config": {
			setup:   setThenRestoreEnv("APP_SERVER_DB_HOST", "a"),
			wantErr: true,
		},
		"with config": {
			wantErr: false,
		},
	}

	conf, err := coreConfig.NewConfigWithDefaults(nil, false)
	require.NoError(t, err)

	logger, err := log.NewLogger(conf)
	require.NoError(t, err)

	for tcName, tc := range tests {

		t.Logf("Running test >%s<", tcName)

		func() {
			if tc.setup != nil {
				teardown := tc.setup(conf)
				defer func() {
					teardown()
				}()
			}

			db, err := getPostgresDB(conf, logger)
			defer func() {
				if db != nil {
					db.Close()
				}
			}()
			if tc.wantErr {
				require.Error(t, err, "getPostgresDB returns with error")
				return
			}

			require.NoError(t, err, "getPostgresDB returns without error")
			require.NotNil(t, db, "getPostgresDB returns db connection")
		}()
	}
}
