package harness

import (
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

var DefaultDataConfig = DataConfig{
	ObjectConfig: []ObjectConfig{
		{
			Record: record.Object{
				Name:                "Rusted Sword",
				Description:         "A rusted sword.",
				DescriptionDetailed: "A rusted sword with a chipped blade and a worn leather handle.",
			},
		},
		{
			Record: record.Object{
				Name:                "Rusted Helmet",
				Description:         "A rusted helmet.",
				DescriptionDetailed: "A rusted helmet pitted with dents.",
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
				Name:                "Dull Bronze Ring",
				Description:         "A dull bronze ring.",
				DescriptionDetailed: "A dull bronze ring.",
			},
		},
		{
			Record: record.Object{
				Name:                "Blood Stained Pouch",
				Description:         "A blood stained pouch.",
				DescriptionDetailed: "A blood stained pouch.",
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
		{
			Record: record.Object{
				Name:                "Stone Mace",
				Description:         "A stone mace.",
				DescriptionDetailed: "A stone mace.",
			},
		},
		{
			Record: record.Object{
				Name:                "Chipped Hammer",
				Description:         "A chipped hammer.",
				DescriptionDetailed: "A chipped hammer.",
			},
		},
		{
			Record: record.Object{
				Name:                "Chipped Breastplate",
				Description:         "A chipped breastplate.",
				DescriptionDetailed: "A chipped breastplate.",
			},
		},
	},
	MonsterConfig: []MonsterConfig{
		{
			Record: record.Monster{
				Name:        "Grumpy Dwarf",
				Description: "A particularly grumpy specimen of a common dwarf.",
			},
			MonsterObjectConfig: []MonsterObjectConfig{
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
				Description: "A particularly angry specimen of a common goblin.",
			},
			MonsterObjectConfig: []MonsterObjectConfig{},
		},
	},
	CharacterConfig: []CharacterConfig{
		{
			Record: record.Character{
				Name: "Barricade",
			},
			CharacterObjectConfig: []CharacterObjectConfig{
				{
					Record: record.CharacterObject{
						IsEquipped: true,
					},
					ObjectName: "Dull Bronze Ring",
				},
				{
					Record: record.CharacterObject{
						IsStashed: true,
					},
					ObjectName: "Blood Stained Pouch",
				},
			},
		},
		{
			Record: record.Character{
				Name: "Legislate",
			},
			CharacterObjectConfig: []CharacterObjectConfig{
				{
					Record: record.CharacterObject{
						IsEquipped: true,
					},
					ObjectName: "Stone Mace",
				},
			},
		},
		{
			Record: record.Character{
				Name: "Bolster",
			},
			CharacterObjectConfig: []CharacterObjectConfig{
				{
					Record: record.CharacterObject{
						IsEquipped: true,
					},
					ObjectName: "Chipped Hammer",
				},
				{
					Record: record.CharacterObject{
						IsEquipped: true,
					},
					ObjectName: "Chipped Breastplate",
				},
			},
		},
	},
	DungeonConfig: []DungeonConfig{
		{
			Record: record.Dungeon{
				Name:        "Cave",
				Description: "A dark and damp stone cave.",
			},
			LocationConfig: []LocationConfig{
				{
					Record: record.Location{
						Name:        "Cave Entrance",
						Description: "A large cave entrance.",
						IsDefault:   true,
					},
					NorthLocationName: "Cave Tunnel",
					LocationMonsterConfig: []LocationMonsterConfig{
						{
							MonsterName: "Grumpy Dwarf",
						},
					},
					LocationObjectConfig: []LocationObjectConfig{
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
					LocationMonsterConfig: []LocationMonsterConfig{
						{
							MonsterName: "Angry Goblin",
						},
					},
					LocationObjectConfig: []LocationObjectConfig{
						{
							ObjectName: "Rusted Helmet",
						},
					},
				},
				{
					Record: record.Location{
						Name:        "Cave Room",
						Description: "A large cave room.",
					},
					SouthLocationName:     "Cave Tunnel",
					LocationMonsterConfig: []LocationMonsterConfig{},
					LocationObjectConfig: []LocationObjectConfig{
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
					LocationMonsterConfig: []LocationMonsterConfig{},
					LocationObjectConfig:  []LocationObjectConfig{},
				},
				{
					Record: record.Location{
						Name:        "Dark Narrow Tunnel",
						Description: "A dark narrow tunnel.",
					},
					SoutheastLocationName: "Narrow Tunnel",
					DownLocationName:      "Dark Room",
					LocationMonsterConfig: []LocationMonsterConfig{},
					LocationObjectConfig:  []LocationObjectConfig{},
				},
				{
					Record: record.Location{
						Name:        "Dark Room",
						Description: "A dark room.",
					},
					UpLocationName:        "Dark Narrow Tunnel",
					LocationMonsterConfig: []LocationMonsterConfig{},
					LocationObjectConfig:  []LocationObjectConfig{},
				},
			},
			DungeonInstanceConfig: []DungeonInstanceConfig{
				{
					CharacterInstanceConfig: []CharacterInstanceConfig{
						{
							Name: "Barricade",
						},
						{
							Name: "Legislate",
						},
					},
					ActionConfig: []ActionConfig{
						{
							CharacterName: "Barricade",
							Command:       "look",
						},
						{
							MonsterName: "Grumpy Dwarf",
							Command:     "look",
						},
					},
				},
			},
		},
	},
}
