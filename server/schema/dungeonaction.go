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
	DroppedObject   *ObjectDetailedData    `json:"dropped_object,omitempty"`
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
	// Health and fatigue is always assigned to show how wounded or
	// tired a character at a location appears
	Health         int `json:"health"`
	Fatigue        int `json:"fatigue"`
	CurrentHealth  int `json:"current_health"`
	CurrentFatigue int `json:"current_fatigue"`
}

type CharacterDetailedData struct {
	Name                string `json:"name"`
	Description         string `json:"description"`
	Strength            int    `json:"strength"`
	Dexterity           int    `json:"dexterity"`
	Intelligence        int    `json:"intelligence"`
	CurrentStrength     int    `json:"current_strength"`
	CurrentDexterity    int    `json:"current_dexterity"`
	CurrentIntelligence int    `json:"current_intelligence"`
	Health              int    `json:"health"`
	Fatigue             int    `json:"fatigue"`
	CurrentHealth       int    `json:"current_health"`
	CurrentFatigue      int    `json:"current_fatigue"`
	// Equipped objects are always assigned for the character
	// performing the action or a target character so that
	// equipped objects are visible to all players
	EquippedObjects []ObjectDetailedData `json:"equipped_objects,omitempty"`
	// Stashed objects are only assigned for the character
	// performing the action so that stashed objects are not
	// exposed to all players
	StashedObjects []ObjectDetailedData `json:"stashed_objects,omitempty"`
	// TODO: Add effects currently applied
}

type MonsterData struct {
	Name string `json:"name"`
	// Health and fatigue is always assigned to show how wounded or
	// tired a monster at a location appears
	Health         int `json:"health"`
	Fatigue        int `json:"fatigue"`
	CurrentHealth  int `json:"current_health"`
	CurrentFatigue int `json:"current_fatigue"`
}

type MonsterDetailedData struct {
	Name                string `json:"name"`
	Description         string `json:"description"`
	Strength            int    `json:"strength"`
	Dexterity           int    `json:"dexterity"`
	Intelligence        int    `json:"intelligence"`
	CurrentStrength     int    `json:"current_strength"`
	CurrentDexterity    int    `json:"current_dexterity"`
	CurrentIntelligence int    `json:"current_intelligence"`
	Health              int    `json:"health"`
	Fatigue             int    `json:"fatigue"`
	CurrentHealth       int    `json:"current_health"`
	CurrentFatigue      int    `json:"current_fatigue"`
	// Equipped objects are always assigned for the monster
	// performing the action or a target monster so that
	// equipped objects are visible to all players
	EquippedObjects []ObjectDetailedData `json:"equipped_objects,omitempty"`
	// Stashed objects are only assigned for the monster
	// performing the action so that stashed objects are not
	// exposed to all players
	StashedObjects []ObjectDetailedData `json:"stashed_objects,omitempty"`
	// TODO: Add effects currently applied
}

type ObjectData struct {
	Name string `json:"name"`
}

type ObjectDetailedData struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IsStashed   bool   `json:"is_stashed"`
	IsEquipped  bool   `json:"is_equipped"`
	// TODO: Add effects that are applied
}
