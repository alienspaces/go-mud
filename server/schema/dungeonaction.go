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
	ID              string         `json:"id,omitempty"`
	Command         string         `json:"command"`
	CommandResult   string         `json:"command_result"`
	Location        LocationData   `json:"location"`
	Character       *CharacterData `json:"character,omitempty"`
	Monster         *MonsterData   `json:"monster,omitempty"`
	EquippedObject  *ObjectData    `json:"equipped_object,omitempty"`
	StashedObject   *ObjectData    `json:"stashed_object,omitempty"`
	TargetObject    *ObjectData    `json:"target_object,omitempty"`
	TargetCharacter *CharacterData `json:"target_character,omitempty"`
	TargetMonster   *MonsterData   `json:"target_monster,omitempty"`
	TargetLocation  *LocationData  `json:"target_location,omitempty"`
	CreatedAt       time.Time      `json:"created_at,omitempty"`
	UpdatedAt       time.Time      `json:"updated_at,omitempty"`
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
	Name        string `json:"name"`
	Description string `json:"description"`
}

type MonsterData struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ObjectData struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
