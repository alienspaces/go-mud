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
	Record                 record.Dungeon
	LocationConfig  []LocationConfig
	DungeonCharacterConfig []DungeonCharacterConfig
	DungeonMonsterConfig   []DungeonMonsterConfig
	DungeonObjectConfig    []DungeonObjectConfig
	ActionConfig    []ActionConfig
	DungeonInstanceConfig[]
}

// DungeonInstanceConfig -
type DungeonInstanceConfig struct {

}

// LocationConfig -
type LocationConfig struct {
	Record                record.Location
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

// ActionConfig -
type ActionConfig struct {
	CharacterName string
	MonsterName   string
	Command       string
}

// DungeonCharacterConfig -
type DungeonCharacterConfig struct {
	Record       record.DungeonCharacter
	LocationName string
}

// DungeonMonsterConfig -
type DungeonMonsterConfig struct {
	Record       record.DungeonMonster
	LocationName string
}

// DungeonObjectConfig - 
type DungeonObjectConfig struct {
	Record        record.DungeonObject
	LocationName  string
	CharacterName string
	MonsterName   string
}
