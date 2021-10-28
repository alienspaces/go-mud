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
	Action     ActionData      `json:"action"`
	Location   LocationData    `json:"location"`
	Character  CharacterData   `json:"character,omitempty"`
	Monster    MonsterData     `json:"monster,omitempty"`
	Characters []CharacterData `json:"characters"`
	Monsters   []MonsterData   `json:"monsters"`
	Objects    []ObjectData    `json:"objects"`
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

type ActionData struct {
	ID                             string    `json:"id,omitempty"`
	Command                        string    `json:"command"`
	CommandResult                  string    `json:"command_result"`
	EquippedDungeonObjectName      string    `json:"equipped_dungeon_object_name,omitempty"`
	StashedDungeonObjectName       string    `json:"stashed_dungeon_object_name,omitempty"`
	TargetDungeonObjectName        string    `json:"target_dungeon_object_name,omitempty"`
	TargetDungeonCharacterName     string    `json:"target_dungeon_character_name,omitempty"`
	TargetDungeonMonsterName       string    `json:"target_dungeon_monster_name,omitempty"`
	TargetDungeonLocationDirection string    `json:"target_dungeon_location_direction,omitempty"`
	TargetDungeonLocationName      string    `json:"target_dungeon_location_name,omitempty"`
	CreatedAt                      time.Time `json:"created_at,omitempty"`
	UpdatedAt                      time.Time `json:"updated_at,omitempty"`
}

type LocationData struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Directions  []string `json:"directions"`
}

type CharacterData struct {
	Name string `json:"name"`
}

type MonsterData struct {
	Name string `json:"name"`
}

type ObjectData struct {
	Name string `json:"name"`
}
