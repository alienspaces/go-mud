package runner

import (
	"github.com/urfave/cli/v2"

	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/harness"
)

// LoadTestData -
func (rnr *Runner) LoadTestData(c *cli.Context) error {

	rnr.Log.Info("** Load Test Data **")

	config := TestDataConfig()

	h, err := harness.NewTesting(rnr.Config, rnr.Log, rnr.Store, config)
	if err != nil {
		rnr.Log.Warn("Failed new testing harness >%v<", err)
		return err
	}

	// harness commit data
	h.ShouldCommitData = true

	_, err = h.Setup()
	if err != nil {
		rnr.Log.Warn("Failed testing harness setup >%v<", err)
		return err
	}

	rnr.Log.Info("All done")

	return nil
}
