package harness

import (
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

// DataConfig -
type DataConfig struct {
	ObjectConfig    []ObjectConfig
	MonsterConfig   []MonsterConfig
	CharacterConfig []CharacterConfig
	DungeonConfig   []DungeonConfig
}

// DungeonConfig -
type DungeonConfig struct {
	Record                record.Dungeon
	LocationConfig        []LocationConfig
	DungeonInstanceConfig []DungeonInstanceConfig
}

// ObjectConfig -
type ObjectConfig struct {
	Record record.Object
}

// MonsterConfig -
type MonsterConfig struct {
	Record              record.Monster
	MonsterObjectConfig []MonsterObjectConfig
}

// MonsterObjectConfig -
type MonsterObjectConfig struct {
	Record record.MonsterObject
	// ObjectName is used to resolve the object identifier of the resulting record
	ObjectName string
}

// CharacterConfig -
type CharacterConfig struct {
	Record                record.Character
	CharacterObjectConfig []CharacterObjectConfig
}

// CharacterObjectConfig -
type CharacterObjectConfig struct {
	Record record.CharacterObject
	// ObjectName is used to resolve the object identifier of the resulting record
	ObjectName string
}

// LocationConfig -
type LocationConfig struct {
	Record record.Location
	// [Direction]LocationName is used to resolve the location identifiers of the resulting record
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

	// Location Monsters
	LocationMonsterConfig []LocationMonsterConfig

	// Location Objects
	LocationObjectConfig []LocationObjectConfig
}

type LocationMonsterConfig struct {
	Record record.LocationMonster
	// MonsterName is used to resolve the monster identifier of the resulting record
	MonsterName string
}

type LocationObjectConfig struct {
	Record record.LocationObject
	// ObjectName is used to resolve the object identifier of the resulting record
	ObjectName string
}

// DungeonInstanceConfig -
type DungeonInstanceConfig struct {
	CharacterInstanceConfig []CharacterInstanceConfig
	ActionConfig            []ActionConfig
}

// CharacterInstanceConfig -
type CharacterInstanceConfig struct {
	Name string
}

// ActionConfig -
type ActionConfig struct {
	CharacterName string
	MonsterName   string
	Command       string
}
