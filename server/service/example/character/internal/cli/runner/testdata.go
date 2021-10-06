package runner

import (
	"github.com/urfave/cli/v2"

	"gitlab.com/alienspaces/go-boilerplate/server/core/repository"
	"gitlab.com/alienspaces/go-boilerplate/server/service/character/internal/harness"
	"gitlab.com/alienspaces/go-boilerplate/server/service/character/internal/record"
)

// LoadTestData -
func (rnr *Runner) LoadTestData(c *cli.Context) error {

	rnr.Log.Info("** Load Test Data **")

	// harness
	config := harness.DataConfig{
		CharacterConfig: []harness.CharacterConfig{
			// Maize
			{
				Record: record.Character{
					Record: repository.Record{
						ID: "1d3f8d0b-b7b3-4569-a099-8f1b6e2a2c71",
					},
					Name: "Maize",
				},
			},
			// Veronica
			{
				Record: record.Character{
					Record: repository.Record{
						ID: "6992c452-dadf-47fd-99fa-64287b44e475",
					},
					Name: "Veronica",
				},
			},
			// Audrey
			{
				Record: record.Character{
					Record: repository.Record{
						ID: "cf0371e6-10e3-4594-a3fb-fd8253cccf2a",
					},
					Name: "Audrey",
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
