package dependencies

import (
	"fmt"

	"gitlab.com/alienspaces/go-mud/server/core/config"
	"gitlab.com/alienspaces/go-mud/server/core/log"
	"gitlab.com/alienspaces/go-mud/server/core/store"
	"gitlab.com/alienspaces/go-mud/server/service/template/internal/model"
)

func Default() (*config.Config, *log.Log, *store.Store, *model.Model, error) {

	// Configurer
	c, err := config.NewConfig(nil, false)
	if err != nil {
		fmt.Printf("failed new config >%v<", err)
		return nil, nil, nil, nil, err
	}

	configVars := []string{
		// general
		"APP_SERVER_ENV",
		"APP_SERVER_PORT",
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
		err := c.Add(key, true)
		if err != nil {
			fmt.Printf("failed adding config item >%v<", err)
			return nil, nil, nil, nil, err
		}
	}

	// Logger
	l, err := log.NewLogger(c)
	if err != nil {
		fmt.Printf("failed new logger >%v<", err)
		return nil, nil, nil, nil, err
	}

	// Storer
	s, err := store.NewStore(c, l)
	if err != nil {
		fmt.Printf("failed new store >%v<", err)
		return nil, nil, nil, nil, err
	}

	// Modeller
	m, err := model.NewModel(c, l, s)
	if err != nil {
		fmt.Printf("failed new model >%v<", err)
		return nil, nil, nil, nil, err
	}

	return c, l, s, m, nil
}
