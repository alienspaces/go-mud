package runner

import (
	"github.com/urfave/cli/v2"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/harness"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// LoadSeedData -
func (rnr *Runner) LoadSeedData(c *cli.Context) error {

	rnr.Log.Info("** Loading Seed Data **")

	// harness
	config := harness.DataConfig{
		DungeonConfig: []harness.DungeonConfig{
			{
				Record: record.Dungeon{
					Name: "Cave",
				},
				DungeonLocationConfig: []harness.DungeonLocationConfig{
					{
						Record: record.DungeonLocation{
							Name:        "Cave Entrance",
							Description: "A large cave entrance.",
							Default:     true,
						},
						NorthLocationName: "Cave Tunnel",
					},
					{
						Record: record.DungeonLocation{
							Name:        "Cave Tunnel",
							Description: "A cave tunnel descends into the mountain.",
						},
						NorthLocationName:     "Cave Room",
						SouthLocationName:     "Cave Entrance",
						NorthwestLocationName: "Narrow Tunnel",
					},
					{
						Record: record.DungeonLocation{
							Name:        "Cave Room",
							Description: "A large cave room.",
						},
						SouthLocationName: "Cave Tunnel",
					},
					{
						Record: record.DungeonLocation{
							Name:        "Narrow Tunnel",
							Description: "A narrow tunnel gradually descending into the darkness.",
						},
						NorthwestLocationName: "Dark Narrow Tunnel",
						SoutheastLocationName: "Cave Tunnel",
					},
					{
						Record: record.DungeonLocation{
							Name:        "Dark Narrow Tunnel",
							Description: "A dark narrow tunnel.",
						},
						SoutheastLocationName: "Narrow Tunnel",
						DownLocationName:      "Dark Room",
					},
					{
						Record: record.DungeonLocation{
							Name:        "Dark Room",
							Description: "A dark room.",
						},
						UpLocationName: "Dark Narrow Tunnel",
					},
				},
				DungeonMonsterConfig: []harness.DungeonMonsterConfig{
					{
						Record: record.DungeonMonster{
							Name: "White Cat",
						},
						LocationName: "Cave Entrance",
					},
					{
						Record: record.DungeonMonster{
							Name: "Angry Goblin",
						},
						LocationName: "Cave Tunnel",
					},
				},
				DungeonObjectConfig: []harness.DungeonObjectConfig{
					{
						Record: record.DungeonObject{
							Name: "Rusted Sword",
						},
						LocationName: "Cave Entrance",
					},
					{
						Record: record.DungeonObject{
							Name: "Silver Key",
						},
						LocationName: "Cave Room",
					},
				},
			},
		},
	}

	h, err := harness.NewTesting(rnr.Config, rnr.Log, rnr.Store, rnr.Model, config)
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
