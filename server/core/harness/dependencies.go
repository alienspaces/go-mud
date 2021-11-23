package harness

import (
	"gitlab.com/alienspaces/go-mud/server/core/config"
	"gitlab.com/alienspaces/go-mud/server/core/log"
	"gitlab.com/alienspaces/go-mud/server/core/store"
	"gitlab.com/alienspaces/go-mud/server/core/type/configurer"
	"gitlab.com/alienspaces/go-mud/server/core/type/logger"
	"gitlab.com/alienspaces/go-mud/server/core/type/storer"
)

func NewDefaultConfig() (configurer.Configurer, error) {

	c, err := config.NewConfig(nil, false)
	if err != nil {
		return nil, err
	}

	configVars := []string{
		// logger
		"APP_SERVER_LOG_LEVEL",
		// database
		"APP_SERVER_DB_HOST",
		"APP_SERVER_DB_PORT",
		"APP_SERVER_DB_NAME",
		"APP_SERVER_DB_USER",
		"APP_SERVER_DB_PASSWORD",
		// schema
		"APP_SERVER_SCHEMA_PATH",
		// jwt signing key
		"APP_SERVER_JWT_SIGNING_KEY",
	}
	for _, key := range configVars {
		err = c.Add(key, true)
		if err != nil {
			return nil, err
		}
	}
	return c, nil
}

func NewDefaultLogger(c configurer.Configurer) (logger.Logger, error) {

	l, err := log.NewLogger(c)
	if err != nil {
		return nil, err
	}

	return l, nil
}

func NewDefaultStorer(c configurer.Configurer, l logger.Logger) (storer.Storer, error) {

	s, err := store.NewStore(c, l)
	if err != nil {
		return nil, err
	}

	err = s.Init()
	if err != nil {
		return nil, err
	}

	return s, nil
}

func NewDefaultDependencies() (configurer.Configurer, logger.Logger, storer.Storer, error) {

	// configurer
	c, err := NewDefaultConfig()
	if err != nil {
		return nil, nil, nil, err
	}

	// logger
	l, err := NewDefaultLogger(c)
	if err != nil {
		return nil, nil, nil, err
	}

	// storer
	s, err := NewDefaultStorer(c, l)
	if err != nil {
		return nil, nil, nil, err
	}

	return c, l, s, nil
}
