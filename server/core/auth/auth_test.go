package auth

import (
	"os"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/go-boilerplate/server/core/config"
	"gitlab.com/alienspaces/go-boilerplate/server/core/log"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/configurer"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/logger"
)

// NewDefaultDependencies -
func NewDefaultDependencies() (configurer.Configurer, logger.Logger, error) {

	// configurer
	c, err := config.NewConfig(nil, false)
	if err != nil {
		return nil, nil, err
	}

	configVars := []string{
		// general
		"APP_HOST",
		// logger
		"APP_SERVER_LOG_LEVEL",
		// auth
		"APP_SERVER_JWT_SIGNING_KEY",
	}
	for _, key := range configVars {
		err = c.Add(key, false)
		if err != nil {
			return nil, nil, err
		}
	}

	// logger
	l, err := log.NewLogger(c)
	if err != nil {
		return nil, nil, err
	}

	return c, l, nil
}

func TestNewAuth(t *testing.T) {

	tests := []struct {
		name      string
		setup     func() error
		expectErr bool
	}{
		{
			name: "Success",
			setup: func() error {
				os.Setenv("APP_SERVER_JWT_SIGNING_KEY", gofakeit.UUID())
				return nil
			},
			expectErr: false,
		},
		{
			name: "Fail",
			setup: func() error {
				os.Unsetenv("APP_SERVER_JWT_SIGNING_KEY")
				return nil
			},
			expectErr: true,
		},
	}

	for _, tc := range tests {

		t.Logf("Running test >%s<", tc.name)

		// Setup
		if tc.setup != nil {
			err := tc.setup()
			require.NoError(t, err, "Setup returns without error")
		}

		c, l, err := NewDefaultDependencies()
		require.NoError(t, err, "NewDefaultDependencies returns without error")

		auth, err := NewAuth(c, l)
		if tc.expectErr == true {
			require.Error(t, err, "NewAuth returns expected error")
			continue
		}
		require.NoError(t, err, "NewAuth returns without error")
		require.NotNil(t, auth, "NewAuth returns a new auth")
	}
}

func TestEncode(t *testing.T) {

	tests := []struct {
		name      string
		claims    func() *Claims
		setup     func() error
		init      func(auth *Auth) error
		expectErr bool
	}{
		{
			name: "Successful encode",
			claims: func() *Claims {
				roles := []string{
					gofakeit.Word(),
					gofakeit.Word(),
					gofakeit.Word(),
				}
				identity := map[string]interface{}{
					gofakeit.Word(): gofakeit.Uint32(),
					gofakeit.Word(): gofakeit.Word(),
					gofakeit.Word(): gofakeit.Bool(),
				}
				c := Claims{
					Roles:    roles,
					Identity: identity,
				}
				return &c
			},
			setup: func() error {
				os.Setenv("APP_SERVER_JWT_SIGNING_KEY", gofakeit.UUID())
				return nil
			},
			init: func(auth *Auth) error {
				auth.JwtExpiryMinutes = -10
				return nil
			},
			expectErr: false,
		},
	}

	for _, tc := range tests {

		t.Logf("Running test >%s<", tc.name)

		// Setup
		if tc.setup != nil {
			err := tc.setup()
			require.NoError(t, err, "Setup returns without error")
		}

		c, l, err := NewDefaultDependencies()
		require.NoError(t, err, "NewDefaultDependencies returns without error")

		auth, err := NewAuth(c, l)
		require.NoError(t, err, "NewAuth returns without error")
		require.NotNil(t, auth, "NewAuth returns a new auth")

		// Init
		if tc.init != nil {
			err := tc.init(auth)
			require.NoError(t, err, "Init returns without error")
		}

		// EncodeJWT
		token, err := auth.EncodeJWT(tc.claims())
		if tc.expectErr == true {
			require.Error(t, err, "EncodeJWT returns expected error")
			continue
		}

		require.NoError(t, err, "EncodeJWT returns without error")
		require.NotEmpty(t, token, "EncodeJWT returns a token")
	}
}

func TestDecode(t *testing.T) {

	tests := []struct {
		name      string
		claims    func() *Claims
		setup     func() error
		init      func(auth *Auth) error
		expectErr bool
	}{
		{
			name: "Successful decode",
			claims: func() *Claims {
				roles := []string{
					gofakeit.Word(),
					gofakeit.Word(),
					gofakeit.Word(),
				}
				identity := map[string]interface{}{
					gofakeit.Word(): gofakeit.Uint32(),
					gofakeit.Word(): gofakeit.Word(),
					gofakeit.Word(): gofakeit.Bool(),
				}
				c := Claims{
					Roles:    roles,
					Identity: identity,
				}
				return &c
			},
			setup: func() error {
				os.Setenv("APP_SERVER_JWT_SIGNING_KEY", gofakeit.UUID())
				return nil
			},
			init: func(auth *Auth) error {
				auth.JwtExpiryMinutes = 10
				return nil
			},
			expectErr: false,
		},
		{
			name: "Failed, expired",
			claims: func() *Claims {
				roles := []string{
					gofakeit.Word(),
					gofakeit.Word(),
					gofakeit.Word(),
				}
				identity := map[string]interface{}{
					gofakeit.Word(): gofakeit.Uint32(),
					gofakeit.Word(): gofakeit.Word(),
					gofakeit.Word(): gofakeit.Bool(),
				}
				c := Claims{
					Roles:    roles,
					Identity: identity,
				}
				return &c
			},
			setup: func() error {
				os.Setenv("APP_SERVER_JWT_SIGNING_KEY", gofakeit.UUID())
				return nil
			},
			init: func(auth *Auth) error {
				auth.JwtExpiryMinutes = -10
				return nil
			},
			expectErr: true,
		},
	}

	for _, tc := range tests {

		t.Logf("Running test >%s<", tc.name)

		// Setup
		if tc.setup != nil {
			err := tc.setup()
			require.NoError(t, err, "Setup returns without error")
		}

		c, l, err := NewDefaultDependencies()
		require.NoError(t, err, "NewDefaultDependencies returns without error")

		auth, err := NewAuth(c, l)
		require.NoError(t, err, "NewAuth returns without error")
		require.NotNil(t, auth, "NewAuth returns a new auth")

		// Init
		if tc.init != nil {
			err := tc.init(auth)
			require.NoError(t, err, "Init returns without error")
		}

		// EncodeJWT
		token, err := auth.EncodeJWT(tc.claims())
		require.NoError(t, err, "EncodeJWT returns without error")
		require.NotEmpty(t, token, "EncodeJWT returns a token")

		// DecodeJWT
		claims, err := auth.DecodeJWT(token)
		if tc.expectErr {
			require.Error(t, err, "DecodeJWT returns expected error")
			continue
		}

		require.NoError(t, err, "DecodeJWT returns without error")
		require.NotEmpty(t, claims, "DecodeJWT returns claims")
	}
}
