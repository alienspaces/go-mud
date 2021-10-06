package runner

import (
	"github.com/urfave/cli/v2"

	"gitlab.com/alienspaces/go-boilerplate/server/core/repository"
	"gitlab.com/alienspaces/go-boilerplate/server/service/player/internal/harness"
	"gitlab.com/alienspaces/go-boilerplate/server/service/player/internal/record"
)

// LoadTestData -
func (rnr *Runner) LoadTestData(c *cli.Context) error {

	rnr.Log.Info("** Load Test Data **")

	// harness
	config := harness.DataConfig{
		PlayerConfig: []harness.PlayerConfig{
			{
				Record: record.Player{
					Record: repository.Record{
						ID: "a11f45c3-a0c7-4f79-a90d-52585c9e1473",
					},
				},
			},
		},
	}

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
