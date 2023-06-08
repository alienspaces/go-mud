package server

import (
	"fmt"

	"gitlab.com/alienspaces/go-mud/backend/core/config"
	"gitlab.com/alienspaces/go-mud/backend/core/type/configurer"
)

type Config struct {
	AppServerEnv             string
	AppServerHost            string
	AppServerHome            string
	AppServerPort            string
	AppImageTagFeatureBranch string
	AppImageTagSHA           string
}

func NewConfig(c configurer.Configurer) (*Config, error) {
	cfg := Config{
		AppServerEnv:  c.Get(config.AppServerEnv),
		AppServerHost: c.Get(config.AppServerHost),
		AppServerHome: c.Get(config.AppServerHome),
		AppServerPort: c.Get(config.AppServerPort),
	}

	// Services need to validate core configuration that may have been applied as
	// the core runner does not know what any specific service requires.

	return &cfg, nil
}

// ValidateHTTP validates the minimum required configuration has been provided
// to run the HTTP server.
func (c *Config) ValidateHTTP() error {
	if c.AppServerEnv == "" {
		return fmt.Errorf("missing required configuration AppServerEnv")
	}
	if c.AppServerHost == "" {
		return fmt.Errorf("missing required configuration AppServerHost")
	}
	if c.AppServerHome == "" {
		return fmt.Errorf("missing required configuration AppServerHome")
	}
	if c.AppServerPort == "" {
		return fmt.Errorf("missing required configuration AppServerPort")
	}
	return nil
}
