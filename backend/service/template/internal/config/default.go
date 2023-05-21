package config

import (
	coreconfig "gitlab.com/alienspaces/go-mud/backend/core/config"
)

func NewConfig(extra []coreconfig.Item, dotEnv bool) (*coreconfig.Config, error) {
	items := DefaultRequiredItemKeys()
	items = append(items, coreconfig.NewItems(DefaultItemKeys(), false)...)
	items = append(items, extra...)

	conf, err := coreconfig.NewConfigWithDefaults(items, dotEnv)
	if err != nil {
		return nil, err
	}

	return conf, nil
}

func DefaultItemKeys() []string {
	return []string{}
}

func DefaultRequiredItemKeys() []coreconfig.Item {
	return coreconfig.NewItems(defaultRequiredItemKeys(), true)
}

func defaultRequiredItemKeys() []string {
	return []string{
		coreconfig.AppServerEnv,
	}
}
