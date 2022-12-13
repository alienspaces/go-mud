package config

func NewConfigWithDefaults(extra []Item, dotEnv bool) (*Config, error) {
	items := NewItems(DefaultRequiredDBItemKeys(), true)
	items = append(items, extra...)

	return NewConfigWithStatelessDefaults(items, dotEnv)
}

func NewConfigWithStatelessDefaults(extra []Item, dotEnv bool) (*Config, error) {
	items := NewItems(DefaultRequiredStatelessItemKeys(), true)
	items = append(items, NewItems(DefaultItemKeys(), false)...)
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

func DefaultRequiredDBItemKeys() []string {
	return []string{
		// database
		"APP_SERVER_DB_HOST",
		"APP_SERVER_DB_PORT",
		"APP_SERVER_DB_NAME",
		"APP_SERVER_DB_USER",
		"APP_SERVER_DB_PASSWORD",
		"APP_SERVER_DB_MAX_OPEN_CONNECTIONS",
		"APP_SERVER_DB_MAX_IDLE_CONNECTIONS",
		"APP_SERVER_DB_MAX_IDLE_TIME_MINS",
	}
}

func DefaultRequiredStatelessItemKeys() []string {
	return []string{
		// general
		"APP_SERVER_ENV",
		"APP_SERVER_HOST",
		"APP_SERVER_PORT",
		// logger
		"APP_SERVER_LOG_LEVEL",
		"APP_SERVER_LOG_PRETTY",
	}
}
