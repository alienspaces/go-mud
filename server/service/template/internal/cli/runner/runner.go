package runner

import (
	"fmt"

	"github.com/urfave/cli/v2"

	command "gitlab.com/alienspaces/go-mud/server/core/cli"
	"gitlab.com/alienspaces/go-mud/server/core/type/configurer"
	"gitlab.com/alienspaces/go-mud/server/core/type/logger"
	"gitlab.com/alienspaces/go-mud/server/core/type/modeller"
	"gitlab.com/alienspaces/go-mud/server/service/template/internal/model"
)

// Runner -
type Runner struct {
	command.Runner
}

// NewRunner -
func NewRunner(c configurer.Configurer, l logger.Logger) (*Runner, error) {

	r := Runner{}

	r.Log = l
	if r.Log == nil {
		msg := "logger undefined, cannot init runner"
		r.Log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	r.Config = c
	if r.Config == nil {
		msg := "configurer undefined, cannot init runner"
		r.Log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	// https://github.com/urfave/cli/blob/master/docs/v2/manual.md
	r.App = &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "load-test-data",
				Aliases: []string{"l"},
				Usage:   "Load a set of test data",
				Action:  r.LoadTestData,
			},
			{
				Name:    "load-seed-data",
				Aliases: []string{"s"},
				Usage:   "Load production seed data",
				Action:  r.LoadSeedData,
			},
			{
				Name:    "test",
				Aliases: []string{"t"},
				Usage:   "Runs the test command",
				Action:  r.TestCommand,
			},
		},
	}

	r.ModellerFunc = r.Modeller

	return &r, nil
}

// TestCommand -
func (rnr *Runner) TestCommand(c *cli.Context) error {

	rnr.Log.Info("** Template Test Command **")

	return nil
}

// Modeller -
func (rnr *Runner) Modeller() (modeller.Modeller, error) {

	rnr.Log.Info("** Template Modeller **")

	m, err := model.NewModel(rnr.Config, rnr.Log, rnr.Store)
	if err != nil {
		rnr.Log.Warn("Failed new model >%v<", err)
		return nil, err
	}

	return m, nil
}
