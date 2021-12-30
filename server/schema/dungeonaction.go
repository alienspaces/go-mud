package schema

import (
	"time"
)

// DungeonActionResponse -
type DungeonActionResponse struct {
	Response
	Data []DungeonActionResponseData `json:"data"`
}

// DungeonActionResponseData -
type DungeonActionResponseData struct {
	ID              string                 `json:"id,omitempty"`
	Command         string                 `json:"command"`
	CommandResult   string                 `json:"command_result"`
	Location        LocationData           `json:"location"`
	Character       *CharacterDetailedData `json:"character,omitempty"`
	Monster         *MonsterDetailedData   `json:"monster,omitempty"`
	EquippedObject  *ObjectDetailedData    `json:"equipped_object,omitempty"`
	StashedObject   *ObjectDetailedData    `json:"stashed_object,omitempty"`
	TargetObject    *ObjectDetailedData    `json:"target_object,omitempty"`
	TargetCharacter *CharacterDetailedData `json:"target_character,omitempty"`
	TargetMonster   *MonsterDetailedData   `json:"target_monster,omitempty"`
	TargetLocation  *LocationData          `json:"target_location,omitempty"`
	CreatedAt       time.Time              `json:"created_at,omitempty"`
	UpdatedAt       time.Time              `json:"updated_at,omitempty"`
}

// DungeonActionRequest -
type DungeonActionRequest struct {
	Request
	Data DungeonActionRequestData `json:"data"`
}

// DungeonActionRequestData -
type DungeonActionRequestData struct {
	Sentence string `json:"sentence"`
}

type LocationData struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Direction   string          `json:"direction,omitempty"`
	Directions  []string        `json:"directions"`
	Characters  []CharacterData `json:"characters,omitempty"`
	Monsters    []MonsterData   `json:"monsters,omitempty"`
	Objects     []ObjectData    `json:"objects,omitempty"`
}

type CharacterData struct {
	Name string `json:"name"`
}

type CharacterDetailedData struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	Strength     int    `json:"strength"`
	Dexterity    int    `json:"dexterity"`
	Intelligence int    `json:"intelligence"`
	Health       int    `json:"health"`
	Fatigue      int    `json:"fatigue"`
}

type MonsterData struct {
	Name string `json:"name"`
}

type MonsterDetailedData struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	Strength     int    `json:"strength"`
	Dexterity    int    `json:"dexterity"`
	Intelligence int    `json:"intelligence"`
	Health       int    `json:"health"`
	Fatigue      int    `json:"fatigue"`
}

type ObjectData struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ObjectDetailedData struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IsStashed   bool   `json:"is_stashed"`
	IsEquipped  bool   `json:"is_equipped"`
}
