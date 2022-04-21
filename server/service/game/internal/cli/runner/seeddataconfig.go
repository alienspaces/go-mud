package runner

import (
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/harness"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

func SeedDataConfig() harness.DataConfig {
	return seedDataConfig
}

var seedDataConfig = harness.DataConfig{
	DungeonConfig: []harness.DungeonConfig{
		{
			Record: record.Dungeon{
				Name: "Cave",
			},
			LocationConfig: []harness.LocationConfig{
				{
					Record: record.Location{
						Name:        "Cave Entrance",
						Description: "A large cave entrance.",
						Default:     true,
					},
					NorthLocationName: "Cave Tunnel",
				},
				{
					Record: record.Location{
						Name:        "Cave Tunnel",
						Description: "A cave tunnel descends into the mountain.",
					},
					NorthLocationName:     "Cave Room",
					SouthLocationName:     "Cave Entrance",
					NorthwestLocationName: "Narrow Tunnel",
				},
				{
					Record: record.Location{
						Name:        "Cave Room",
						Description: "A large cave room.",
					},
					SouthLocationName: "Cave Tunnel",
				},
				{
					Record: record.Location{
						Name:        "Narrow Tunnel",
						Description: "A narrow tunnel gradually descending into the darkness.",
					},
					NorthwestLocationName: "Dark Narrow Tunnel",
					SoutheastLocationName: "Cave Tunnel",
				},
				{
					Record: record.Location{
						Name:        "Dark Narrow Tunnel",
						Description: "A dark narrow tunnel.",
					},
					SoutheastLocationName: "Narrow Tunnel",
					DownLocationName:      "Dark Room",
				},
				{
					Record: record.Location{
						Name:        "Dark Room",
						Description: "A dark room.",
					},
					UpLocationName: "Dark Narrow Tunnel",
				},
			},
			DungeonMonsterConfig: []harness.DungeonMonsterConfig{
				{
					Record: record.DungeonMonster{
						Name: "Grumpy Dwarf",
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
						Name:                "Rusted Sword",
						Description:         "A rusted sword.",
						DescriptionDetailed: "A rusted sword with a chipped blade and a worn leather handle.",
					},
					LocationName: "Cave Entrance",
				},
				{
					Record: record.DungeonObject{
						Name:                "Silver Key",
						Description:         "A silver key.",
						DescriptionDetailed: "A silver key with fine runes in a language you do not understand engraved along the edge.",
					},
					LocationName: "Cave Room",
				},
				{
					Record: record.DungeonObject{
						Name:                "Bone Dagger",
						Description:         "A bone dagger.",
						DescriptionDetailed: "A bone dagger.",
						IsEquipped:          true,
					},
					MonsterName: "Grumpy Dwarf",
				},
				{
					Record: record.DungeonObject{
						Name:                "Vial Of Ogre Blood",
						Description:         "A large vial of ogre blood.",
						DescriptionDetailed: "A large vial of ogre blood.",
						IsStashed:           true,
					},
					MonsterName: "Grumpy Dwarf",
				},
			},
		},
	},
}
