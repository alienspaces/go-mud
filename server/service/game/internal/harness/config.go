package harness

import (
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
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
	Record         record.Dungeon
	LocationConfig []LocationConfig
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
	Record                 record.Character
	CharacterObjectConfig  []CharacterObjectConfig
	CharacterDungeonConfig *CharacterDungeonConfig
}

// CharacterObjectConfig -
type CharacterObjectConfig struct {
	Record record.CharacterObject
	// ObjectName is used to resolve the object identifier of the resulting record
	ObjectName string
}

// CharacterDungeonConfig creates a character instances inside a new or existing
// dungeon instance.
type CharacterDungeonConfig struct {
	// DungeonName is used to resolve the dungeon identifier of the dungeon and
	// dungeon instance a character instance will be created in.
	DungeonName string
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

	// TODO: A character can only have one active instance in a dungeon at a time.
	// When configuration contains multiple  CharacterInstanceConfig definitions
	// the result of applying all ActionConfig definitions should result in only
	// one active character instance.
	CharacterInstanceConfig []CharacterInstanceConfig

	// TODO: Change this to TurnConfig that contains a list of ActionConfig
	// definitions for character and monster instances to perform
	ActionConfig []ActionConfig
}

// CharacterInstanceConfig -
type CharacterInstanceConfig struct {
	Record record.CharacterInstance
	// CharacterName is used to resolve the required character
	// identifier of the character instance record
	CharacterName string
	// LocationName is used to resolve the required location
	// identifier of the character instance record
	LocationName string
}

// TODO: Change this to TurnConfig that contains a list of ActionConfig
// definitions for character and monster instances to perform

// ActionConfig -
type ActionConfig struct {
	// XxxxName is used to resolve the required character or
	// monster instance identifier of the action record
	CharacterName string
	MonsterName   string
	Command       string
}
