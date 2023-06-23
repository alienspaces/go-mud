package schema

import (
	"strings"
	"time"

	"gitlab.com/alienspaces/go-mud/backend/schema"
)

// ActionResponse -
type ActionResponse struct {
	schema.Response
	Data []ActionResponseData `json:"data"`
}

// ActionResponseData -
type ActionResponseData struct {
	ID              string           `json:"id"`
	Command         string           `json:"command"`
	Narrative       string           `json:"narrative"`
	TurnNumber      int              `json:"turn_number"`
	SerialNumber    int16            `json:"serial_number"`
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
	schema.Request
	Data ActionRequestData `json:"data"`
}

// ActionRequestData -
type ActionRequestData struct {
	Sentence string `json:"sentence"`
}

type ActionLocationCharacters []ActionLocationCharacter

func (alc ActionLocationCharacters) HasCharacterWithName(name string) bool {
	for i := range alc {
		if strings.EqualFold(alc[i].Name, name) {
			return true
		}
	}
	return false
}

type ActionLocationMonsters []ActionLocationMonster

func (alm ActionLocationMonsters) HasMonsterWithName(name string) bool {
	for i := range alm {
		if strings.EqualFold(alm[i].Name, name) {
			return true
		}
	}
	return false
}

type ActionLocationObjects []ActionLocationObject

func (alo ActionLocationObjects) HasObjectWithName(name string) bool {
	for i := range alo {
		if strings.EqualFold(alo[i].Name, name) {
			return true
		}
	}
	return false
}

type ActionLocation struct {
	Name        string                   `json:"name"`
	Description string                   `json:"description"`
	Direction   string                   `json:"direction,omitempty"`
	Directions  []string                 `json:"directions"`
	Characters  ActionLocationCharacters `json:"characters,omitempty"`
	Monsters    ActionLocationMonsters   `json:"monsters,omitempty"`
	Objects     ActionLocationObjects    `json:"objects,omitempty"`
}

// ActionLocationCharacter describes a character that is at a location
type ActionLocationCharacter struct {
	Name           string `json:"name"`
	Health         int    `json:"health"`  // Health is always assigned to show how wounded a character at a location appears
	Fatigue        int    `json:"fatigue"` // Fatigue is always assigned to show how tired a character at a location appears
	CurrentHealth  int    `json:"current_health"`
	CurrentFatigue int    `json:"current_fatigue"`
}

type ActionCharacter struct {
	Name                string         `json:"name"`
	Description         string         `json:"description"`
	Strength            int            `json:"strength"`
	Dexterity           int            `json:"dexterity"`
	Intelligence        int            `json:"intelligence"`
	CurrentStrength     int            `json:"current_strength"`
	CurrentDexterity    int            `json:"current_dexterity"`
	CurrentIntelligence int            `json:"current_intelligence"`
	Health              int            `json:"health"`
	Fatigue             int            `json:"fatigue"`
	CurrentHealth       int            `json:"current_health"`
	CurrentFatigue      int            `json:"current_fatigue"`
	Coins               int            `json:"coins,omitempty"`            // Coins are only assigned for the character performing the action so that a characters coins are not visible to other players.
	ExperiencePoints    int            `json:"experience_points"`          // ExperiencePoints are only assigned for the character performing the action so that a characters experience points are not visible to other players.
	AttributePoints     int            `json:"attribute_points"`           // AttributePoints are only assigned for the character performing the action so that a characters attribute points are not visible to other players.
	EquippedObjects     []ActionObject `json:"equipped_objects,omitempty"` // Equipped objects are always assigned for the character performing the action or a target character so that equipped objects are visible to all players.
	StashedObjects      []ActionObject `json:"stashed_objects,omitempty"`  // Stashed objects are only assigned for the character performing the action so that stashed objects are not visible to other players.
}

type ActionLocationMonster struct {
	Name           string `json:"name"`
	Health         int    `json:"health"`  // Health is always assigned to show how wounded a monster at a location appears
	Fatigue        int    `json:"fatigue"` // Fatigue is always assigned to show how tired a monster at a location appears
	CurrentHealth  int    `json:"current_health"`
	CurrentFatigue int    `json:"current_fatigue"`
}

type ActionMonster struct {
	Name                string         `json:"name"`
	Description         string         `json:"description"`
	Strength            int            `json:"strength"`
	Dexterity           int            `json:"dexterity"`
	Intelligence        int            `json:"intelligence"`
	CurrentStrength     int            `json:"current_strength"`
	CurrentDexterity    int            `json:"current_dexterity"`
	CurrentIntelligence int            `json:"current_intelligence"`
	Health              int            `json:"health"`
	Fatigue             int            `json:"fatigue"`
	CurrentHealth       int            `json:"current_health"`
	CurrentFatigue      int            `json:"current_fatigue"`
	EquippedObjects     []ActionObject `json:"equipped_objects,omitempty"` // Equipped objects are always assigned for the monster performing the action or a target monster so that equipped objects are visible to all players.
}

type ActionLocationObject struct {
	Name string `json:"name"`
}

type ActionObject struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IsStashed   bool   `json:"is_stashed"`
	IsEquipped  bool   `json:"is_equipped"`
}
