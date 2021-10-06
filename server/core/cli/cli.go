package cli

import (
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/configurer"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/logger"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/preparer"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/runnable"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/storer"
)

// CLI -
type CLI struct {
	Config  configurer.Configurer
	Log     logger.Logger
	Store   storer.Storer
	Prepare preparer.Preparer
	Runner  runnable.Runnable
}

// NewCLI -
func NewCLI(c configurer.Configurer, l logger.Logger, s storer.Storer, r runnable.Runnable) (*CLI, error) {

	cli := CLI{
		Config: c,
		Log:    l,
		Store:  s,
		Runner: r,
	}

	err := cli.Init()
	if err != nil {
		return nil, err
	}

	return &cli, nil
}

// Init -
func (cli *CLI) Init() error {

	err := cli.Store.Init()
	if err != nil {
		return err
	}

	// TODO: alerting, retries
	return cli.Runner.Init(cli.Config, cli.Log, cli.Store)
}

// Run -
func (cli *CLI) Run(args map[string]interface{}) error {

	// TODO:
	// - alerting on errors
	// - retries on start up
	// - reload  on config changes
	return cli.Runner.Run(args)
}
