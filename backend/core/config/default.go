package config

func NewConfigWithDefaults(extra []Item, dotEnv bool) (*Config, error) {
	items := NewItems(DefaultRequiredDBItemKeys(), true)
	items = append(items, NewItems(DefaultRequiredItemKeys(), true)...)
	items = append(items, NewItems(DefaultItemKeys(), false)...)
	items = append(items, extra...)

	conf, err := NewConfig(items, dotEnv)
	if err != nil {
		return nil, err
	}

	return conf, nil
}

func DefaultItemKeys() []string {
	return []string{}
}

func DefaultRequiredItemKeys() []string {
	return []string{
		// general
		AppServerEnv,
		AppServerHost,
		AppServerPort,
		// logger
		AppServerLogLevel,
		AppServerLogPretty,
	}
}

func DefaultRequiredDBItemKeys() []string {
	return []string{
		// database
		AppServerDBHost,
		AppServerDBPort,
		AppServerDBName,
		AppServerDBUser,
		AppServerDBPassword,
		AppServerDBMaxOpenConnections,
		AppServerDBMaxIdleConnections,
		AppServerDBMaxIdleTimeMins,
	}
}
