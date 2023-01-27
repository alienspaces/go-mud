package cabin

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
					ID: "9d7b1765-5c76-461a-aa87-0e2c99803c0f",
				},
				Name: "Bolster",
			},
		},
	}
}

func MonsterConfig() []harness.MonsterConfig {
	return []harness.MonsterConfig{
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
					ID: "1075fada-b173-4e23-97c0-32ce013786e4",
				},
				Name:                "Yellow Chewed Bone",
				Description:         "A yellow chewed bone.",
				DescriptionDetailed: "A yellowed and chewed human arm bone.",
			},
		},
	}
}

func DungeonConfig() harness.DungeonConfig {
	return harness.DungeonConfig{
		Record: record.Dungeon{
			Record: repository.Record{
				ID: "34c5b913-3079-42a6-8228-3f1fb8f20dbe",
			},
			Name:        "Cabin",
			Description: "A wood cabin.",
		},
		LocationConfig: []harness.LocationConfig{
			{
				Record: record.Location{
					Record: repository.Record{
						ID: "87096f48-8c7a-4512-b134-a6c77662de9b",
					},
					Name:        "Cabin Verandah",
					Description: "A wooden boarded verandah.",
					IsDefault:   true,
				},
				NorthLocationName: "Cabin Room",
			},
			{
				Record: record.Location{
					Record: repository.Record{
						ID: "cf24c4f0-13bb-470a-8b3b-12e80c575c8c",
					},
					Name:        "Cabin Room",
					Description: "A mostly empty cabin room.",
				},
				SouthLocationName: "Cabin Verandah",
				LocationObjectConfig: []harness.LocationObjectConfig{
					{
						ObjectName: "Yellow Chewed Bone",
					},
				},
				LocationMonsterConfig: []harness.LocationMonsterConfig{
					{
						MonsterName: "Giant Grey Rat",
					},
				},
			},
		},
		DungeonInstanceConfig: []harness.DungeonInstanceConfig{
			{
				CharacterInstanceConfig: []harness.CharacterInstanceConfig{
					{
						Name: "Bolster",
					},
				},
				TurnConfig: []harness.TurnConfig{
					{
						ActionConfig: []harness.ActionConfig{
							{
								CharacterName: "Bolster",
								Command:       "look",
							},
						},
					},
					{
						ActionConfig: []harness.ActionConfig{
							{
								CharacterName: "Bolster",
								Command:       "look north",
							},
						},
					},
				},
			},
		},
	}
}
