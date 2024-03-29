package cli

import (
	"gitlab.com/alienspaces/go-mud/backend/core/type/configurer"
	"gitlab.com/alienspaces/go-mud/backend/core/type/logger"
	"gitlab.com/alienspaces/go-mud/backend/core/type/preparer"
	"gitlab.com/alienspaces/go-mud/backend/core/type/runnable"
	"gitlab.com/alienspaces/go-mud/backend/core/type/storer"
)

// CLI -
type CLI struct {
	Config             configurer.Configurer
	Log                logger.Logger
	Store              storer.Storer
	RepositoryPreparer preparer.Repository
	QueryPreparer      preparer.Query
	Runner             runnable.Runnable
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

	// TODO OX-71: Is this needed?
	// err := cli.Store.Init()
	// if err != nil {
	// 	return err
	// }

	// TODO: alerting, retries
	return cli.Runner.Init(cli.Store)
}

// Run -
func (cli *CLI) Run(args map[string]interface{}) error {

	// TODO:
	// - alerting on errors
	// - retries on start up
	// - reload  on config changes
	return cli.Runner.Run(args)
}
