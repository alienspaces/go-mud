package schema

import (
	"time"
)

// ActionResponse -
type ActionResponse struct {
	Response
	Data []ActionResponseData `json:"data"`
}

// ActionResponseData -
type ActionResponseData struct {
	ID              string           `json:"id,omitempty"`
	Command         string           `json:"command"`
	Narrative       string           `json:"narrative"`
	Location        ActionLocation   `json:"location"`
	Character       *ActionCharacter `json:"character,omitempty"`
	Monster         *ActionMonster   `json:"monster,omitempty"`
	EquippedObject  *ActionObject    `json:"equipped_object,omitempty"`
	StashedObject   *ActionObject    `json:"stashed_object,omitempty"`
	DroppedObject   *ActionObject    `json:"dropped_object,omitempty"`
	TargetObject    *ActionObject    `json:"target_object,omitempty"`
	TargetCharacter *ActionCharacter `json:"target_character,omitempty"`
	TargetMonster   *ActionMonster   `json:"target_monster,omitempty"`
	TargetLocation  *ActionLocation  `json:"target_location,omitempty"`
	CreatedAt       time.Time        `json:"created_at,omitempty"`
	UpdatedAt       time.Time        `json:"updated_at,omitempty"`
}

// ActionRequest -
type ActionRequest struct {
	Request
	Data ActionRequestData `json:"data"`
}

// ActionRequestData -
type ActionRequestData struct {
	Sentence string `json:"sentence"`
}

type ActionLocation struct {
	Name        string                    `json:"name"`
	Description string                    `json:"description"`
	Direction   string                    `json:"direction,omitempty"`
	Directions  []string                  `json:"directions"`
	Characters  []ActionLocationCharacter `json:"characters,omitempty"`
	Monsters    []ActionLocationMonster   `json:"monsters,omitempty"`
	Objects     []ActionLocationObject    `json:"objects,omitempty"`
}

type ActionLocationCharacter struct {
	Name string `json:"name"`
	// Health and fatigue is always assigned to show how wounded or
	// tired a character at a location appears
	Health         int `json:"health"`
	Fatigue        int `json:"fatigue"`
	CurrentHealth  int `json:"current_health"`
	CurrentFatigue int `json:"current_fatigue"`
}

type ActionCharacter struct {
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
	// Coins are only assigned for the character performing
	// the action so that a characters coins are not visible
	// to other players.
	Coins int `json:"coins,omitempty"`
	// ExperiencePoints are only assigned for the character performing
	// the action so that a characters experience points are not visible
	// to other players.
	ExperiencePoints int `json:"experience_points"`
	// ExperiencePoints are only assigned for the character performing
	// the action so that a characters attribute points are not visible
	// to other players.
	AttributePoints int `json:"attribute_points"`
	// Equipped objects are always assigned for the character performing
	// the action or a target character so that equipped objects are
	// visible to all players.
	EquippedObjects []ActionObject `json:"equipped_objects,omitempty"`
	// Stashed objects are only assigned for the character performing
	// the action so that stashed objects are not visible to other players.
	StashedObjects []ActionObject `json:"stashed_objects,omitempty"`
	// TODO: (game) Add effects currently applied
}

type ActionLocationMonster struct {
	Name string `json:"name"`
	// Health and fatigue is always assigned to show how wounded or
	// tired a monster at a location appears
	Health         int `json:"health"`
	Fatigue        int `json:"fatigue"`
	CurrentHealth  int `json:"current_health"`
	CurrentFatigue int `json:"current_fatigue"`
}

type ActionMonster struct {
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
	// Equipped objects are always assigned for the monster performing
	// the action or a target monster so that equipped objects are
	// visible to all players.
	EquippedObjects []ActionObject `json:"equipped_objects,omitempty"`
	// TODO: (game) Add effects currently applied
}

type ActionLocationObject struct {
	Name string `json:"name"`
}

type ActionObject struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IsStashed   bool   `json:"is_stashed"`
	IsEquipped  bool   `json:"is_equipped"`
	// TODO: (game) Add effects that are applied
}
