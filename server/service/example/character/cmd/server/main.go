package main

import (
	"fmt"
	"os"

	"gitlab.com/alienspaces/go-boilerplate/server/core/config"
	"gitlab.com/alienspaces/go-boilerplate/server/core/log"
	"gitlab.com/alienspaces/go-boilerplate/server/core/server"
	"gitlab.com/alienspaces/go-boilerplate/server/core/store"
	"gitlab.com/alienspaces/go-boilerplate/server/service/character/internal/server/runner"
)

func main() {

	c, err := config.NewConfig(nil, false)
	if err != nil {
		fmt.Printf("Failed new config >%v<", err)
		os.Exit(0)
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
			fmt.Printf("Failed adding config item >%v<", err)
			os.Exit(0)
		}
	}

	l, err := log.NewLogger(c)
	if err != nil {
		fmt.Printf("Failed new logger >%v<", err)
		os.Exit(0)
	}

	s, err := store.NewStore(c, l)
	if err != nil {
		fmt.Printf("Failed new store >%v<", err)
		os.Exit(0)
	}

	r := runner.NewRunner()

	svr, err := server.NewServer(c, l, s, r)
	if err != nil {
		fmt.Printf("Failed new server >%v<", err)
		os.Exit(0)
	}

	args := make(map[string]interface{})

	err = svr.Run(args)
	if err != nil {
		fmt.Printf("Failed server run >%v<", err)
		os.Exit(0)
	}

	os.Exit(1)
}
