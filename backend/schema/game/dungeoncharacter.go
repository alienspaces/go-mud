package schema

import (
	"time"

	"gitlab.com/alienspaces/go-mud/backend/schema"
)

// DungeonCharacterResponse -
type DungeonCharacterResponse struct {
	schema.Response
	Data []DungeonCharacterData `json:"data"`
}

// DungeonCharacterData -
type DungeonCharacterData struct {
	DungeonID                    string    `json:"dungeon_id,omitempty"`
	DungeonName                  string    `json:"dungeon_name,omitempty"`
	DungeonDescription           string    `json:"dungeon_description,omitempty"`
	LocationID                   string    `json:"location_id,omitempty"`
	LocationName                 string    `json:"location_name,omitempty"`
	LocationDescription          string    `json:"location_description,omitempty"`
	CharacterID                  string    `json:"character_id,omitempty"`
	CharacterName                string    `json:"character_name"`
	CharacterStrength            int       `json:"character_strength"`
	CharacterDexterity           int       `json:"character_dexterity"`
	CharacterIntelligence        int       `json:"character_intelligence"`
	CharacterCurrentStrength     int       `json:"character_current_strength"`
	CharacterCurrentDexterity    int       `json:"character_current_dexterity"`
	CharacterCurrentIntelligence int       `json:"character_current_intelligence"`
	CharacterHealth              int       `json:"character_health"`
	CharacterFatigue             int       `json:"character_fatigue"`
	CharacterCurrentHealth       int       `json:"character_current_health"`
	CharacterCurrentFatigue      int       `json:"character_current_fatigue"`
	CharacterCoins               int       `json:"character_coins,omitempty"`
	CharacterExperiencePoints    int       `json:"character_experience_points"`
	CharacterAttributePoints     int       `json:"character_attribute_points"`
	CharacterCreatedAt           time.Time `json:"character_created_at,omitempty"`
	CharacterUpdatedAt           time.Time `json:"character_updated_at,omitempty"`
}
