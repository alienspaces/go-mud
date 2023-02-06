package config

func NewConfigWithDefaults(extra []Item, dotEnv bool) (*Config, error) {
	items := NewItems(DefaultRequiredItemKeys(), true)
	items = append(items, extra...)

	conf, err := NewConfig(items, dotEnv)
	if err != nil {
		return nil, err
	}

	return conf, nil
}

func DefaultItemKeys() []string {
	return []string{
		// build info for documentation
		"APP_IMAGE_TAG_FEATURE_BRANCH",
		"APP_IMAGE_TAG_SHA",
	}
}

func DefaultRequiredItemKeys() []string {
	return []string{
		AppServerHome,
		AppServerEnv,
		AppServerPort,
		AppServerLogLevel,
		AppServerLogPretty,
		AppServerDbHost,
		AppServerDbPort,
		AppServerDbName,
		AppServerDbUser,
		AppServerDbPassword,
		AppServerDbMaxOpenConnections,
		AppServerDbMaxIdleConnections,
		AppServerDbMaxIdleTimeMins,
		AppServerSchemaPath,
		AppServerJwtSigningKey,
	}
}
