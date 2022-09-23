package runner

import (
	"github.com/urfave/cli/v2"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/harness"
)

// LoadSeedData -
func (rnr *Runner) LoadSeedData(c *cli.Context) error {

	rnr.Log.Info("** Loading Seed Data **")

	config := SeedDataConfig()

	h, err := harness.NewTesting(rnr.Config, rnr.Log, rnr.Store, config)
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
