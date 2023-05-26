package salesforce

import (
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/go-mud/backend/core/config"
	"gitlab.com/alienspaces/go-mud/backend/core/log"
	"gitlab.com/alienspaces/go-mud/backend/core/type/configurer"
)

// newClient - test helper to create a new client
func newClient() (*Client, error) {

	c, err := config.NewConfigWithDefaults(nil, false)
	if err != nil {
		return nil, err
	}

	l := log.NewLogger(c)

	source, err := NewClient(l, GetSourceConfig(c))
	if err != nil {
		return nil, err
	}

	return source, nil
}

func TestNewClient(t *testing.T) {
	scc, err := newClient()
	require.NoError(t, err, "TestNewClient does not return an error")
	require.NotNil(t, scc, "TestNewClient returns a source client")
}

func GetSourceConfig(c configurer.Configurer) ClientConfig {
	return ClientConfig{
		Version:      "v54.0",
		Username:     "notarealusername",
		Password:     "notarealpassword",
		ClientID:     "notarealclientid",
		ClientSecret: "notarealclientsecret",
		Environment:  "sandbox",
	}
}
