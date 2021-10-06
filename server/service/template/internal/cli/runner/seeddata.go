package runner

import (
	"github.com/urfave/cli/v2"

	"gitlab.com/alienspaces/go-boilerplate/server/service/template/internal/harness"
)

// LoadSeedData -
func (rnr *Runner) LoadSeedData(c *cli.Context) error {

	rnr.Log.Info("** Load Seed Data **")

	// harness
	config := harness.DataConfig{}

	h, err := harness.NewTesting(config)
	if err != nil {
		rnr.Log.Warn("Failed new testing harness >%v<", err)
		return err
	}

	// harness commit data
	h.CommitData = true

	err = h.Setup()
	if err != nil {
		rnr.Log.Warn("Failed testing harness setup >%v<", err)
		return err
	}

	rnr.Log.Info("All done")

	return nil
}
