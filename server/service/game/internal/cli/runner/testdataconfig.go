package runner

import (
	"gitlab.com/alienspaces/go-mud/server/core/repository"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/harness"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

var testDataConfig = harness.DataConfig{
	DungeonConfig: []harness.DungeonConfig{
		{
			Record: record.Dungeon{
				Record: repository.Record{
					ID: "55087d68-dc17-41ed-bb53-12dc636ac196",
				},
				Name: "Cave",
			},
			DungeonLocationConfig: []harness.DungeonLocationConfig{
				{
					Record: record.DungeonLocation{
						Record: repository.Record{
							ID: "b47febf0-3c51-405e-8f41-abc18d20392a",
						},
						Name:        "Cave Entrance",
						Description: "A large cave entrance.",
						Default:     true,
					},
					NorthLocationName: "Cave Tunnel",
				},
				{
					Record: record.DungeonLocation{
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
					Record: record.DungeonLocation{
						Record: repository.Record{
							ID: "08c75bd1-13e7-4b44-a594-c87a125885d0",
						},
						Name:        "Cave Room",
						Description: "A large cave room.",
					},
					SouthLocationName: "Cave Tunnel",
				},
				{
					Record: record.DungeonLocation{
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
					Record: record.DungeonLocation{
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
					Record: record.DungeonLocation{
						Record: repository.Record{
							ID: "dc2206de-795b-4864-8d52-36c55d80d33e",
						},
						Name:        "Dark Room",
						Description: "A dark room.",
					},
					UpLocationName: "Dark Narrow Tunnel",
				},
			},
			DungeonCharacterConfig: []harness.DungeonCharacterConfig{
				{
					Record: record.DungeonCharacter{
						Record: repository.Record{
							ID: "38efe8fc-a228-484b-b476-ff0d961942a6",
						},
						Name: "Barricade",
					},
					LocationName: "Cave Entrance",
				},
			},
			DungeonMonsterConfig: []harness.DungeonMonsterConfig{
				{
					Record: record.DungeonMonster{
						Record: repository.Record{
							ID: "1e8179aa-fc2e-4f5a-abe1-e70a237739f5",
						},
						Name: "Grumpy Dwarf",
					},
					LocationName: "Cave Entrance",
				},
				{
					Record: record.DungeonMonster{
						Record: repository.Record{
							ID: "e25cbb71-8fac-4734-a0c1-4c00df729beb",
						},
						Name: "Angry Goblin",
					},
					LocationName: "Cave Tunnel",
				},
			},
			DungeonObjectConfig: []harness.DungeonObjectConfig{
				{
					Record: record.DungeonObject{
						Record: repository.Record{
							ID: "54cf320b-6485-4e86-973e-6a5016f809fd",
						},
						Name:                "Rusted Sword",
						Description:         "A rusted sword.",
						DescriptionDetailed: "A rusted sword with a chipped blade and a worn leather handle.",
					},
					LocationName: "Cave Entrance",
				},
				{
					Record: record.DungeonObject{
						Record: repository.Record{
							ID: "86fa6a84-c23a-45de-8ede-f79966e2ce07",
						},
						Name:                "Silver Key",
						Description:         "A silver key.",
						DescriptionDetailed: "A silver key with fine runes in a language you do not understand engraved along the edge.",
					},
					LocationName: "Cave Room",
				},
				{
					Record: record.DungeonObject{
						Record: repository.Record{
							ID: "792477df-f63d-4cba-b7a4-695a1f896c67",
						},
						Name:                "Dull Bronze Ring",
						Description:         "A bronze ring.",
						DescriptionDetailed: "A dull bronze ring.",
						IsEquipped:          true,
					},
					CharacterName: "Barricade",
				},
			},
			DungeonActionConfig: []harness.DungeonActionConfig{
				{
					CharacterName: "Barricade",
					Command:       "look north",
				},
			},
		},
	},
}
