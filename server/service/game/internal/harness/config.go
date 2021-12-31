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
