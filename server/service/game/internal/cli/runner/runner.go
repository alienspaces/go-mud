package runner

import (
	"github.com/urfave/cli/v2"

	command "gitlab.com/alienspaces/go-mud/server/core/cli"
	"gitlab.com/alienspaces/go-mud/server/core/type/modeller"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/model"
)

// Runner -
type Runner struct {
	command.Runner
}

// NewRunner -
func NewRunner() *Runner {

	r := Runner{}

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

	return &r
}

// TestCommand -
func (rnr *Runner) TestCommand(c *cli.Context) error {

	rnr.Log.Info("** Game Test Command **")

	return nil
}

// Modeller -
func (rnr *Runner) Modeller() (modeller.Modeller, error) {

	rnr.Log.Info("** Game Modeller **")

	m, err := model.NewModel(rnr.Config, rnr.Log, rnr.Store)
	if err != nil {
		rnr.Log.Warn("Failed new model >%v<", err)
		return nil, err
	}

	return m, nil
}
