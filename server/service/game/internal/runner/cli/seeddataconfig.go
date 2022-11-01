package runner

import (
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/harness"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// The following data is actual game data used to seed a new game server
func SeedDataConfig() harness.DataConfig {
	return seedDataConfig
}

var seedDataConfig = harness.DataConfig{
	ObjectConfig: []harness.ObjectConfig{
		{
			Record: record.Object{
				Name:                "Rusted Sword",
				Description:         "A rusted sword.",
				DescriptionDetailed: "A rusted sword with a chipped blade and a worn leather handle.",
			},
		},
		{
			Record: record.Object{
				Name:                "Silver Key",
				Description:         "A silver key.",
				DescriptionDetailed: "A silver key with fine runes in a language you do not understand engraved along the edge.",
			},
		},
		{
			Record: record.Object{
				Name:                "Bone Dagger",
				Description:         "A bone dagger.",
				DescriptionDetailed: "A bone dagger.",
			},
		},
		{
			Record: record.Object{
				Name:                "Vial Of Ogre Blood",
				Description:         "A large vial of ogre blood.",
				DescriptionDetailed: "A large vial of ogre blood.",
			},
		},
	},
	MonsterConfig: []harness.MonsterConfig{
		{
			Record: record.Monster{
				Name:        "Grumpy Dwarf",
				Description: "A particularly grumpy specimen of a dwarf",
			},
			MonsterObjectConfig: []harness.MonsterObjectConfig{
				{
					Record: record.MonsterObject{
						IsEquipped: true,
					},
					ObjectName: "Bone Dagger",
				},
				{
					Record: record.MonsterObject{
						IsStashed: true,
					},
					ObjectName: "Vial Of Ogre Blood",
				},
			},
		},
		{
			Record: record.Monster{
				Name:        "Angry Goblin",
				Description: "A particularly angry specimen of a goblin",
			},
		},
	},
	DungeonConfig: []harness.DungeonConfig{
		{
			Record: record.Dungeon{
				Name:        "Cave",
				Description: "A dark and damp stone cave.",
			},
			LocationConfig: []harness.LocationConfig{
				{
					Record: record.Location{
						Name:        "Cave Entrance",
						Description: "A large cave entrance.",
						IsDefault:   true,
					},
					NorthLocationName: "Cave Tunnel",
					LocationMonsterConfig: []harness.LocationMonsterConfig{
						{
							MonsterName: "Grumpy Dwarf",
						},
					},
					LocationObjectConfig: []harness.LocationObjectConfig{
						{
							ObjectName: "Rusted Sword",
						},
					},
				},
				{
					Record: record.Location{
						Name:        "Cave Tunnel",
						Description: "A cave tunnel descends into the mountain.",
					},
					NorthLocationName:     "Cave Room",
					SouthLocationName:     "Cave Entrance",
					NorthwestLocationName: "Narrow Tunnel",
					LocationMonsterConfig: []harness.LocationMonsterConfig{
						{
							MonsterName: "Angry Goblin",
						},
					},
					LocationObjectConfig: []harness.LocationObjectConfig{},
				},
				{
					Record: record.Location{
						Name:        "Cave Room",
						Description: "A large cave room.",
					},
					SouthLocationName: "Cave Tunnel",
					LocationObjectConfig: []harness.LocationObjectConfig{
						{
							ObjectName: "Silver Key",
						},
					},
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
		},
	},
}
