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
	ActionID              string           `json:"action_id"`
	ActionCommand         string           `json:"action_command"`
	ActionNarrative       string           `json:"action_narrative"`
	ActionTurnNumber      int              `json:"action_turn_number"`
	ActionSerialNumber    int16            `json:"action_serial_number"`
	ActionLocation        ActionLocation   `json:"action_location"`
	ActionCharacter       *ActionCharacter `json:"action_character,omitempty"`
	ActionMonster         *ActionMonster   `json:"action_monster,omitempty"`
	ActionEquippedObject  *ActionObject    `json:"action_equipped_object,omitempty"`
	ActionStashedObject   *ActionObject    `json:"action_stashed_object,omitempty"`
	ActionDroppedObject   *ActionObject    `json:"action_dropped_object,omitempty"`
	ActionTargetObject    *ActionObject    `json:"action_target_object,omitempty"`
	ActionTargetCharacter *ActionCharacter `json:"action_target_character,omitempty"`
	ActionTargetMonster   *ActionMonster   `json:"action_target_monster,omitempty"`
	ActionTargetLocation  *ActionLocation  `json:"action_target_location,omitempty"`
	ActionCreatedAt       time.Time        `json:"action_created_at,omitempty"`
	ActionUpdatedAt       time.Time        `json:"action_updated_at,omitempty"`
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
		if strings.EqualFold(alc[i].CharacterName, name) {
			return true
		}
	}
	return false
}

type ActionLocationMonsters []ActionLocationMonster

func (alm ActionLocationMonsters) HasMonsterWithName(name string) bool {
	for i := range alm {
		if strings.EqualFold(alm[i].MonsterName, name) {
			return true
		}
	}
	return false
}

type ActionLocationObjects []ActionLocationObject

func (alo ActionLocationObjects) HasObjectWithName(name string) bool {
	for i := range alo {
		if strings.EqualFold(alo[i].ObjectName, name) {
			return true
		}
	}
	return false
}

type ActionLocation struct {
	LocationName        string                   `json:"location_name"`
	LocationDescription string                   `json:"location_description"`
	LocationDirection   string                   `json:"location_direction,omitempty"`
	LocationDirections  []string                 `json:"location_directions"`
	LocationCharacters  ActionLocationCharacters `json:"location_characters,omitempty"`
	LocationMonsters    ActionLocationMonsters   `json:"location_monsters,omitempty"`
	LocationObjects     ActionLocationObjects    `json:"location_objects,omitempty"`
}

// ActionLocationCharacter describes a character that is at a location
type ActionLocationCharacter struct {
	CharacterName           string `json:"character_name"`
	CharacterHealth         int    `json:"character_health"`  // Health is always assigned to show how wounded a character at a location appears
	CharacterFatigue        int    `json:"character_fatigue"` // Fatigue is always assigned to show how tired a character at a location appears
	CharacterCurrentHealth  int    `json:"character_current_health"`
	CharacterCurrentFatigue int    `json:"character_current_fatigue"`
}

type ActionCharacter struct {
	CharacterName                string         `json:"character_name"`
	CharacterDescription         string         `json:"character_description"`
	CharacterStrength            int            `json:"character_strength"`
	CharacterDexterity           int            `json:"character_dexterity"`
	CharacterIntelligence        int            `json:"character_intelligence"`
	CharacterCurrentStrength     int            `json:"character_current_strength"`
	CharacterCurrentDexterity    int            `json:"character_current_dexterity"`
	CharacterCurrentIntelligence int            `json:"character_current_intelligence"`
	CharacterHealth              int            `json:"character_health"`
	CharacterFatigue             int            `json:"character_fatigue"`
	CharacterCurrentHealth       int            `json:"character_current_health"`
	CharacterCurrentFatigue      int            `json:"character_current_fatigue"`
	CharacterCoins               int            `json:"character_coins,omitempty"`            // Coins are only assigned for the character performing the action so that a characters coins are not visible to other players.
	CharacterExperiencePoints    int            `json:"character_experience_points"`          // ExperiencePoints are only assigned for the character performing the action so that a characters experience points are not visible to other players.
	CharacterAttributePoints     int            `json:"character_attribute_points"`           // AttributePoints are only assigned for the character performing the action so that a characters attribute points are not visible to other players.
	CharacterEquippedObjects     []ActionObject `json:"character_equipped_objects,omitempty"` // Equipped objects are always assigned for the character performing the action or a target character so that equipped objects are visible to all players.
	CharacterStashedObjects      []ActionObject `json:"character_stashed_objects,omitempty"`  // Stashed objects are only assigned for the character performing the action so that stashed objects are not visible to other players.
}

type ActionLocationMonster struct {
	MonsterName           string `json:"monster_name"`
	MonsterHealth         int    `json:"monster_health"`  // Health is always assigned to show how wounded a monster at a location appears
	MonsterFatigue        int    `json:"monster_fatigue"` // Fatigue is always assigned to show how tired a monster at a location appears
	MonsterCurrentHealth  int    `json:"monster_current_health"`
	MonsterCurrentFatigue int    `json:"monster_current_fatigue"`
}

type ActionMonster struct {
	MonsterName                string         `json:"monster_name"`
	MonsterDescription         string         `json:"monster_description"`
	MonsterStrength            int            `json:"monster_strength"`
	MonsterDexterity           int            `json:"monster_dexterity"`
	MonsterIntelligence        int            `json:"monster_intelligence"`
	MonsterCurrentStrength     int            `json:"monster_current_strength"`
	MonsterCurrentDexterity    int            `json:"monster_current_dexterity"`
	MonsterCurrentIntelligence int            `json:"monster_current_intelligence"`
	MonsterHealth              int            `json:"monster_health"`
	MonsterFatigue             int            `json:"monster_fatigue"`
	MonsterCurrentHealth       int            `json:"monster_current_health"`
	MonsterCurrentFatigue      int            `json:"monster_current_fatigue"`
	MonsterEquippedObjects     []ActionObject `json:"monster_equipped_objects,omitempty"` // Equipped objects are always assigned for the monster performing the action or a target monster so that equipped objects are visible to all players.
}

type ActionLocationObject struct {
	ObjectName string `json:"object_name"`
}

type ActionObject struct {
	ObjectName        string `json:"object_name"`
	ObjectDescription string `json:"object_description"`
	ObjectIsStashed   bool   `json:"object_is_stashed"`
	ObjectIsEquipped  bool   `json:"object_is_equipped"`
}
