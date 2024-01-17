package cave

// The Cave is full of aggressive monsters and traps that are all hell bent on
// cutting and crushing the life from anyone and anything  that ventures within.

import (
	"gitlab.com/alienspaces/go-mud/backend/core/repository"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/harness"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

func CharacterConfig() []harness.CharacterConfig {
	return []harness.CharacterConfig{
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
	}
}

func MonsterConfig() []harness.MonsterConfig {
	return []harness.MonsterConfig{
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
		{
			Record: record.Monster{
				Record: repository.Record{
					ID: "fb032d39-48a4-4806-bd56-f9cba910bbf4",
				},
				Name:        "Giant Grey Rat",
				Description: "A very large grey rat.",
			},
		},
	}
}

func ObjectConfig() []harness.ObjectConfig {
	return []harness.ObjectConfig{
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
	}
}

func DungeonConfig() harness.DungeonConfig {
	return harness.DungeonConfig{
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
				LocationMonsterConfig: []harness.LocationMonsterConfig{
					{
						MonsterName: "Giant Grey Rat",
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
				LocationMonsterConfig: []harness.LocationMonsterConfig{
					{
						MonsterName: "Angry Goblin",
					},
				},
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
				LocationMonsterConfig: []harness.LocationMonsterConfig{
					{
						MonsterName: "Grumpy Dwarf",
					},
				},
			},
		},
		DungeonInstanceConfig: []harness.DungeonInstanceConfig{
			{
				CharacterInstanceConfig: []harness.CharacterInstanceConfig{
					{
						Name: "Barricade",
					},
				},
				TurnConfig: []harness.TurnConfig{
					{
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
}
