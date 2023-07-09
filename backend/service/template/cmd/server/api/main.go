package main

import (
	"os"

	coreconfig "gitlab.com/alienspaces/go-mud/backend/core/config"
	"gitlab.com/alienspaces/go-mud/backend/core/log"
	"gitlab.com/alienspaces/go-mud/backend/core/server"
	"gitlab.com/alienspaces/go-mud/backend/core/store"

	templateConfig "gitlab.com/alienspaces/go-mud/backend/service/template/internal/config"
	runner "gitlab.com/alienspaces/go-mud/backend/service/template/internal/runner/api"
)

func main() {
	l := log.NewDefaultLogger()

	// Required environment sourced variables.
	items := coreconfig.NewItems(
		[]string{}, true,
	)

	// Optional environment sourced variables.

	c, err := templateConfig.NewConfig(items, false)
	if err != nil {
		l.Error("failed new config >%v<", err)
		os.Exit(1)
	}

	// If a new logger instance variable is instantiated, the existing logger instance will be unused
	// and not be garbage collected during the run loop
	l, err = log.NewLogger(c)
	if err != nil {
		l.Error("failed new logger >%v<", err)
		os.Exit(1)
	}

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

	svc, err := server.NewServer(c, l, s, r)
	if err != nil {
		l.Error("failed new server >%v<", err)
		os.Exit(1)
	}

	args := make(map[string]interface{})

	err = svc.Run(args)
	if err != nil {
		l.Error("failed server run >%v<", err)
		os.Exit(1)
	}

	os.Exit(0)
}
