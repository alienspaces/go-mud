package config

import (
	"fmt"

	"gitlab.com/alienspaces/go-mud/backend/core/config"
	"gitlab.com/alienspaces/go-mud/backend/core/type/configurer"
)

const (
	AppServerTurnDuration string = "APP_SERVER_TURN_DURATION"
)

type Config struct {
	config.Config
}

var _ configurer.Configurer = &Config{}

func NewConfig() (*Config, error) {

	// Core default required items
	items := config.NewItems(config.DefaultRequiredDBItemKeys(), true)
	items = append(items, config.NewItems(config.DefaultRequiredItemKeys(), true)...)
	items = append(items, config.NewItems(config.DefaultItemKeys(), false)...)

	// Additional service required items
	items = append(items, config.NewItems([]string{
		AppServerTurnDuration,
	}, true)...)

	cc, err := config.NewConfig(items, false)
	if err != nil {
		return nil, fmt.Errorf("NewConfig failed >%v<", err)
	}

	if cc == nil {
		return nil, fmt.Errorf("NewConfig returned a nil config")
	}

	c := Config{
		Config: *cc,
	}

	return &c, nil
}
