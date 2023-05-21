package server

import (
	"fmt"

	"gitlab.com/alienspaces/go-mud/backend/core/type/configurer"
)

// The following constants are used to source environment variables when
// establishing runner configuration.
const (
	EnvKeyAppVariant               = "APP_VARIANT"
	EnvKeyAppServerEnv             = "APP_SERVER_ENV"
	EnvKeyAppServerHost            = "APP_SERVER_HOST"
	EnvKeyAppServerHome            = "APP_SERVER_HOME"
	EnvKeyAppServerPort            = "APP_SERVER_PORT"
	EnvKeyAppImageTagFeatureBranch = "APP_IMAGE_TAG_FEATURE_BRANCH"
	EnvKeyAppImageTagSHA           = "APP_IMAGE_TAG_SHA"
)

type Config struct {
	AppVariant               string
	AppServerEnv             string
	AppServerHost            string
	AppServerHome            string
	AppServerPort            string
	AppImageTagFeatureBranch string
	AppImageTagSHA           string
}

func NewConfig(c configurer.Configurer) (*Config, error) {
	cfg := Config{
		AppVariant:               c.Get(EnvKeyAppVariant),
		AppServerEnv:             c.Get(EnvKeyAppServerEnv),
		AppServerHost:            c.Get(EnvKeyAppServerHost),
		AppServerHome:            c.Get(EnvKeyAppServerHome),
		AppServerPort:            c.Get(EnvKeyAppServerPort),
		AppImageTagFeatureBranch: c.Get(EnvKeyAppImageTagFeatureBranch),
		AppImageTagSHA:           c.Get(EnvKeyAppImageTagSHA),
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
