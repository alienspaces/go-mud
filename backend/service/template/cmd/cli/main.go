package main

import (
	"os"

	"gitlab.com/alienspaces/go-mud/backend/core/cli"
	"gitlab.com/alienspaces/go-mud/backend/core/log"
	"gitlab.com/alienspaces/go-mud/backend/core/store"

	"gitlab.com/alienspaces/go-mud/backend/service/template/internal/cli"
	templateConfig "gitlab.com/alienspaces/go-mud/backend/service/template/internal/config"
)

func main() {
	l := log.NewDefaultLogger()

	c, err := templateConfig.NewConfig(nil, false)
	if err != nil {
		l.Error("failed new config >%v<", err)
		os.Exit(1)
	}

	// If a new logger instance variable is instantiated, the existing logger instance will be unused
	// and not be garbage collected during the run loop
	l = log.NewLogger(c)

	s, err := store.NewStore(c, l)
	if err != nil {
		l.Error("failed new store >%v<", err)
		os.Exit(1)
	}

	r, err := runner.NewRunner(c, l)
	if err != nil {
		l.Error("failed new runner >%v<", err)
		os.Exit(1)
	}

	cli, err := cli.NewCLI(c, l, s, r)
	if err != nil {
		l.Error("failed new server >%v<", err)
		os.Exit(1)
	}

	args := make(map[string]interface{})

	err = cli.Run(args)
	if err != nil {
		l.Error("failed server run >%v<", err)
		os.Exit(1)
	}

	os.Exit(0)
}
