package runner

import (
	"gitlab.com/alienspaces/go-mud/server/core/repository"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/harness"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// The following data is test data used to seed a test game servers
func TestDataConfig() harness.DataConfig {
	return testDataConfig
}

var testDataConfig = harness.DataConfig{
	ObjectConfig: []harness.ObjectConfig{
		{
			Record: record.Object{
				Record: repository.Record{
					ID: "54cf320b-6485-4e86-973e-6a5016f809fd",
				},
				Name:                "Rusted Sword",
				Description:         "A rusted sword.",
				DescriptionDetailed: "A rusted sword with a chipped blade and a worn leather handle.",
			},
		},
		{
			Record: record.Object{
				Record: repository.Record{
					ID: "86fa6a84-c23a-45de-8ede-f79966e2ce07",
				},
				Name:                "Silver Key",
				Description:         "A silver key.",
				DescriptionDetailed: "A silver key with fine runes in a language you do not understand engraved along the edge.",
			},
		},
		{
			Record: record.Object{
				Record: repository.Record{
					ID: "792477df-f63d-4cba-b7a4-695a1f896c67",
				},
				Name:                "Dull Bronze Ring",
				Description:         "A bronze ring.",
				DescriptionDetailed: "A dull bronze ring.",
			},
		},
	},
	MonsterConfig: []harness.MonsterConfig{
		{
			Record: record.Monster{
				Record: repository.Record{
					ID: "1e8179aa-fc2e-4f5a-abe1-e70a237739f5",
				},
				Name:        "Grumpy Dwarf",
				Description: "A particularly grumpy specimen of a dwarf",
			},
		},
		{
			Record: record.Monster{
				Record: repository.Record{
					ID: "e25cbb71-8fac-4734-a0c1-4c00df729beb",
				},
				Name:        "Angry Goblin",
				Description: "A particularly angrey specimen of a goblin",
			},
		},
	},
	CharacterConfig: []harness.CharacterConfig{
		{
			Record: record.Character{
				Record: repository.Record{
					ID: "38efe8fc-a228-484b-b476-ff0d961942a6",
				},
				Name: "Barricade",
			},
			CharacterObjectConfig: []harness.CharacterObjectConfig{
				{
					Record: record.CharacterObject{
						IsEquipped: true,
					},
					ObjectName: "Dull Bronze Ring",
				},
			},
		},
	},
	DungeonConfig: []harness.DungeonConfig{
		{
			Record: record.Dungeon{
				Record: repository.Record{
					ID: "55087d68-dc17-41ed-bb53-12dc636ac196",
				},
				Name:        "Cave",
				Description: "A dark and damp stone cave.",
			},
			LocationConfig: []harness.LocationConfig{
				{
					Record: record.Location{
						Record: repository.Record{
							ID: "b47febf0-3c51-405e-8f41-abc18d20392a",
						},
						Name:        "Cave Entrance",
						Description: "A large cave entrance.",
						IsDefault:   true,
					},
					NorthLocationName: "Cave Tunnel",
					LocationObjectConfig: []harness.LocationObjectConfig{
						{
							ObjectName: "Rusted Sword",
						},
					},
					LocationMonsterConfig: []harness.LocationMonsterConfig{
						{
							MonsterName: "Grumpy Dwarf",
						},
					},
				},
				{
					Record: record.Location{
						Record: repository.Record{
							ID: "8bfd0e7d-7249-43f7-ab8b-176f32b962bb",
						},
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
				},
				{
					Record: record.Location{
						Record: repository.Record{
							ID: "08c75bd1-13e7-4b44-a594-c87a125885d0",
						},
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
						Record: repository.Record{
							ID: "4a6697e9-11df-4fc4-8a9b-8e26a0c64a21",
						},
						Name:        "Narrow Tunnel",
						Description: "A narrow tunnel gradually descending into the darkness.",
					},
					NorthwestLocationName: "Dark Narrow Tunnel",
					SoutheastLocationName: "Cave Tunnel",
				},
				{
					Record: record.Location{
						Record: repository.Record{
							ID: "495c7346-caa0-4993-bb9a-b9e48fee0ac1",
						},
						Name:        "Dark Narrow Tunnel",
						Description: "A dark narrow tunnel.",
					},
					SoutheastLocationName: "Narrow Tunnel",
					DownLocationName:      "Dark Room",
				},
				{
					Record: record.Location{
						Record: repository.Record{
							ID: "dc2206de-795b-4864-8d52-36c55d80d33e",
						},
						Name:        "Dark Room",
						Description: "A dark room.",
					},
					UpLocationName: "Dark Narrow Tunnel",
				},
			},
			DungeonInstanceConfig: []harness.DungeonInstanceConfig{
				{
					CharacterInstanceConfig: []harness.CharacterInstanceConfig{
						{
							Name: "Barricade",
						},
					},
					ActionConfig: []harness.ActionConfig{
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
