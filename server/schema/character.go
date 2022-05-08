package schema

import (
	"time"
)

// CharacterResponse -
type CharacterResponse struct {
	Response
	Data []CharacterData `json:"data"`
}

// CharacterRequest -
type CharacterRequest struct {
	Request
	Data CharacterData `json:"data"`
}

// CharacterData -
type CharacterData struct {
	ID                  string    `json:"id,omitempty"`
	DungeonID           string    `json:"dungeon_id,omitempty"`
	LocationID          string    `json:"dungeon_location_id,omitempty"`
	Name                string    `json:"name"`
	Strength            int       `json:"strength"`
	Dexterity           int       `json:"dexterity"`
	Intelligence        int       `json:"intelligence"`
	CurrentStrength     int       `json:"current_strength"`
	CurrentDexterity    int       `json:"current_dexterity"`
	CurrentIntelligence int       `json:"current_intelligence"`
	Health              int       `json:"health"`
	Fatigue             int       `json:"fatigue"`
	CurrentHealth       int       `json:"current_health"`
	CurrentFatigue      int       `json:"current_fatigue"`
	Coins               int       `json:"coins,omitempty"`
	ExperiencePoints    int       `json:"experience_points"`
	AttributePoints     int       `json:"attribute_points"`
	CreatedAt           time.Time `json:"created_at,omitempty"`
	UpdatedAt           time.Time `json:"updated_at,omitempty"`
}
