package harness

import (
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

var DefaultDataConfig = DataConfig{
	CharacterConfig: []CharacterConfig{
		{
			Record: record.Character{
				Name: "Barricade",
			},
		},
	},
	DungeonConfig: []DungeonConfig{
		{
			Record: record.Dungeon{
				Name: "Cave",
			},
			LocationConfig: []LocationConfig{
				{
					Record: record.Location{
						Name:        "Cave Entrance",
						Description: "A large cave entrance.",
						IsDefault:   true,
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
			MonsterConfig: []MonsterConfig{
				// 0
				{
					Record: record.Monster{
						Name: "Grumpy Dwarf",
					},
					LocationName: "Cave Entrance",
				},
				// 1
				{
					Record: record.Monster{
						Name: "Angry Goblin",
					},
					LocationName: "Cave Tunnel",
				},
			},
			ObjectConfig: []ObjectConfig{
				// 0
				{
					Record: record.Object{
						Name:                "Rusted Sword",
						Description:         "A rusted sword.",
						DescriptionDetailed: "A rusted sword with a chipped blade and a worn leather handle.",
					},
					LocationName: "Cave Entrance",
				},
				// 1
				{
					Record: record.Object{
						Name:                "Silver Key",
						Description:         "A silver key.",
						DescriptionDetailed: "A silver key with fine runes in a language you do not understand engraved along the edge.",
					},
					LocationName: "Cave Room",
				},
				// 2
				{
					Record: record.Object{
						Name:                "Dull Bronze Ring",
						Description:         "A dull bronze ring.",
						DescriptionDetailed: "A dull bronze ring.",
						IsEquipped:          true,
					},
					CharacterName: "Barricade",
				},
				// 3
				{
					Record: record.Object{
						Name:                "Blood Stained Pouch",
						Description:         "A blood stained pouch.",
						DescriptionDetailed: "A blood stained pouch.",
						IsStashed:           true,
					},
					CharacterName: "Barricade",
				},
				// 4
				{
					Record: record.Object{
						Name:                "Bone Dagger",
						Description:         "A bone dagger.",
						DescriptionDetailed: "A bone dagger.",
						IsEquipped:          true,
					},
					MonsterName: "Grumpy Dwarf",
				},
				// 5
				{
					Record: record.Object{
						Name:                "Vial Of Ogre Blood",
						Description:         "A large vial of ogre blood.",
						DescriptionDetailed: "A large vial of ogre blood.",
						IsStashed:           true,
					},
					MonsterName: "Grumpy Dwarf",
				},
			},
			// MOVE:
			// ActionConfig: []ActionConfig{
			// 	{
			// 		CharacterName: "Barricade",
			// 		Command:       "look north",
			// 	},
			// 	{
			// 		CharacterName: "Barricade",
			// 		Command:       "look grumpy dwarf",
			// 	},
			// },
		},
	},
}
