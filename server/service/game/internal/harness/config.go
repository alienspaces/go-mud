package harness

import (
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// DataConfig -
type DataConfig struct {
	DungeonConfig []DungeonConfig
}

// DungeonConfig -
type DungeonConfig struct {
	Record record.Dungeon
	// DefaultDungeonLocationName string
	DungeonLocationConfig  []DungeonLocationConfig
	DungeonCharacterConfig []DungeonCharacterConfig
	DungeonMonsterConfig   []DungeonMonsterConfig
	DungeonObjectConfig    []DungeonObjectConfig
	DungeonActionConfig    []DungeonActionConfig
}

type DungeonLocationConfig struct {
	Record                record.DungeonLocation
	NorthLocationName     string
	NortheastLocationName string
	EastLocationName      string
	SoutheastLocationName string
	SouthLocationName     string
	SouthwestLocationName string
	WestLocationName      string
	NorthwestLocationName string
	UpLocationName        string
	DownLocationName      string
}

type DungeonActionConfig struct {
	CharacterName string
	MonsterName   string
	Command       string
}

type DungeonCharacterConfig struct {
	Record       record.DungeonCharacter
	LocationName string
}

type DungeonMonsterConfig struct {
	Record       record.DungeonMonster
	LocationName string
}

type DungeonObjectConfig struct {
	Record       record.DungeonObject
	LocationName string
}

var DefaultDataConfig = DataConfig{
	DungeonConfig: []DungeonConfig{
		{
			Record: record.Dungeon{
				Name: "Cave",
			},
			DungeonLocationConfig: []DungeonLocationConfig{
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
					NorthLocationName: "Cave Room",
					SouthLocationName: "Cave Entrance",
				},
				{
					Record: record.DungeonLocation{
						Name:        "Cave Room",
						Description: "A large cave room.",
					},
					SouthLocationName: "Cave Tunnel",
				},
			},
			DungeonCharacterConfig: []DungeonCharacterConfig{
				{
					Record: record.DungeonCharacter{
						Name: "Barricade",
					},
					LocationName: "Cave Entrance",
				},
			},
			DungeonMonsterConfig: []DungeonMonsterConfig{
				{
					Record: record.DungeonMonster{
						Name: "Angry Goblin",
					},
					LocationName: "Cave Tunnel",
				},
			},
			DungeonObjectConfig: []DungeonObjectConfig{
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
			DungeonActionConfig: []DungeonActionConfig{
				{
					CharacterName: "Barricade",
					Command:       "look north",
				},
			},
		},
	},
}
