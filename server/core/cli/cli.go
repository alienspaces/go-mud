package cli

import (
	"gitlab.com/alienspaces/go-mud/server/core/type/configurer"
	"gitlab.com/alienspaces/go-mud/server/core/type/logger"
	"gitlab.com/alienspaces/go-mud/server/core/type/modeller"
	"gitlab.com/alienspaces/go-mud/server/core/type/preparer"
	"gitlab.com/alienspaces/go-mud/server/core/type/runnable"
	"gitlab.com/alienspaces/go-mud/server/core/type/storer"
)

// CLI -
type CLI struct {
	Config  configurer.Configurer
	Log     logger.Logger
	Store   storer.Storer
	Model   modeller.Modeller
	Prepare preparer.Preparer
	Runner  runnable.Runnable
}

// NewCLI -
func NewCLI(c configurer.Configurer, l logger.Logger, s storer.Storer, m modeller.Modeller, r runnable.Runnable) (*CLI, error) {

	cli := CLI{
		Config: c,
		Log:    l,
		Store:  s,
		Model:  m,
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
	return cli.Runner.Init(cli.Config, cli.Log, cli.Store, cli.Model)
}

// Run -
func (cli *CLI) Run(args map[string]interface{}) error {

	// TODO:
	// - alerting on errors
	// - retries on start up
	// - reload  on config changes
	return cli.Runner.Run(args)
}
