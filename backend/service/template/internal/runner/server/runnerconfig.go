package runner

import (
	"gitlab.com/alienspaces/go-mud/backend/core/server"
	"gitlab.com/alienspaces/go-mud/backend/core/type/configurer"
)

// The following constants are used to source environment variables when
// establishing runner configuration.
const (
// Add here..
// EnvKeyAppAPIServerXxx
)

// Config includes core server Config along with additional service
// specific configuration.
type Config struct {
	server.Config
	// Add here..
	// AppAPIServerXxx
}

func NewConfig(c configurer.Configurer) (*Config, error) {
	ccfg, err := server.NewConfig(c)
	if err != nil {
		return nil, err
	}
	cfg := Config{
		Config: *ccfg,
		// Add here..
		// AppAPIServerXxx: c.Get(EnvKeyAppAPIServerXxx),
	}

	// Core configuration provides core configuration validation functions
	err = cfg.ValidateHTTP()
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
