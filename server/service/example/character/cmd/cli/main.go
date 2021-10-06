package main

import (
	"fmt"
	"os"

	"gitlab.com/alienspaces/go-boilerplate/server/core/cli"
	"gitlab.com/alienspaces/go-boilerplate/server/core/config"
	"gitlab.com/alienspaces/go-boilerplate/server/core/log"
	"gitlab.com/alienspaces/go-boilerplate/server/core/store"

	"gitlab.com/alienspaces/go-boilerplate/server/service/character/internal/cli/runner"
)

func main() {

	c, err := config.NewConfig(nil, false)
	if err != nil {
		fmt.Printf("Failed new config >%v<", err)
		os.Exit(1)
	}

	configVars := []string{
		// general
		"APP_SERVER_ENV",
		// logger
		"APP_SERVER_LOG_LEVEL",
		// database
		"APP_SERVER_DB_HOST",
		"APP_SERVER_DB_PORT",
		"APP_SERVER_DB_NAME",
		"APP_SERVER_DB_USER",
		"APP_SERVER_DB_PASSWORD",
	}
	for _, key := range configVars {
		err := c.Add(key, true)
		if err != nil {
			fmt.Printf("Failed adding config item >%v<", err)
			os.Exit(1)
		}
	}

	l, err := log.NewLogger(c)
	if err != nil {
		fmt.Printf("Failed new logger >%v<", err)
		os.Exit(1)
	}

	s, err := store.NewStore(c, l)
	if err != nil {
		fmt.Printf("Failed new store >%v<", err)
		os.Exit(1)
	}

	r := runner.NewRunner()

	cli, err := cli.NewCLI(c, l, s, r)
	if err != nil {
		fmt.Printf("Failed new cli >%v<", err)
		os.Exit(1)
	}

	args := make(map[string]interface{})

	err = cli.Run(args)
	if err != nil {
		fmt.Printf("Failed cli run >%v<", err)
		os.Exit(1)
	}

	os.Exit(0)
}
