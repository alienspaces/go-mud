package runner

import (
	"fmt"

	"github.com/urfave/cli/v2"

	command "gitlab.com/alienspaces/go-mud/backend/core/cli"
	"gitlab.com/alienspaces/go-mud/backend/core/prepare"
	"gitlab.com/alienspaces/go-mud/backend/core/type/configurer"
	"gitlab.com/alienspaces/go-mud/backend/core/type/logger"
	"gitlab.com/alienspaces/go-mud/backend/core/type/modeller"
	"gitlab.com/alienspaces/go-mud/backend/core/type/preparer"

	"gitlab.com/alienspaces/go-mud/backend/service/template/internal/model"
)

// Runner -
type Runner struct {
	command.Runner
}

// NewRunner -
func NewRunner(c configurer.Configurer, l logger.Logger) (*Runner, error) {

	r := Runner{}
	r.DeferModelInitialisation = true

	r.Log = l
	if r.Log == nil {
		msg := "logger undefined, cannot init CLI runner"
		r.Log.Error(msg)
		return nil, fmt.Errorf(msg)
	}
	r.Log = r.Log.WithApplicationContext("cli")

	r.Config = c
	if r.Config == nil {
		msg := "configurer undefined, cannot init CLI runner"
		r.Log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	r.RepositoryPreparerFunc = r.RepositoryPreparer
	r.QueryPreparerFunc = r.QueryPreparer
	r.ModellerFunc = r.Modeller

	// https://github.com/urfave/cli/blob/master/docs/v2/manual.md
	r.App = &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "test",
				Aliases: []string{"t"},
				Usage:   "Runs the test command",
				Action:  r.TestCommand,
			},
		},
	}

	return &r, nil
}

// TestCommand -
func (rnr *Runner) TestCommand(c *cli.Context) error {

	rnr.Log.Info("** Test Command **")

	return nil
}

// RepositoryPreparer -
func (rnr *Runner) RepositoryPreparer() (preparer.Repository, error) {

	rnr.Log.Info("** Template Repository **")

	p, err := prepare.NewRepositoryPreparer(rnr.Log)
	if err != nil {
		rnr.Log.Warn("Failed new preparer repository >%v<", err)
		return nil, err
	}

	return p, nil
}

// QueryPreparer -
func (rnr *Runner) QueryPreparer() (preparer.Query, error) {

	rnr.Log.Info("** Template Query **")

	p, err := prepare.NewQueryPreparer(rnr.Log)
	if err != nil {
		rnr.Log.Warn("Failed new prepare config >%v<", err)
		return nil, err
	}

	return p, nil
}

// Modeller -
func (rnr *Runner) Modeller() (modeller.Modeller, error) {

	rnr.Log.Info("** Template Model **")

	m, err := model.NewModel(rnr.Config, rnr.Log, rnr.Store)
	if err != nil {
		rnr.Log.Warn("Failed new model >%v<", err)
		return nil, err
	}

	return m, nil
}
