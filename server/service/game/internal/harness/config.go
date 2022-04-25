package harness

import (
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// DataConfig -
type DataConfig struct {
	CharacterConfig []CharacterConfig
	DungeonConfig   []DungeonConfig
}

// DungeonConfig -
type DungeonConfig struct {
	Record                record.Dungeon
	LocationConfig        []LocationConfig
	MonsterConfig         []MonsterConfig
	ObjectConfig          []ObjectConfig
	DungeonInstanceConfig []DungeonInstanceConfig
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

// LocationConfig -
type LocationConfig struct {
	Record record.Location
	// [Direction]LocationName is used to resolve the required
	// location identifiers of the location record
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

// CharacterConfig -
type CharacterConfig struct {
	Record record.Character
}

// MonsterConfig -
type MonsterConfig struct {
	Record record.Monster
	// LocationName is used to resolve the required location
	// identifier of the monster record
	LocationName string
}

// ObjectConfig -
type ObjectConfig struct {
	Record record.Object
	// XxxxName is used to resolve the required location, character
	// or monster identifier of the object record
	LocationName  string
	CharacterName string
	MonsterName   string
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
