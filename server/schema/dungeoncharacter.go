package schema

import (
	"time"
)

// DungeonCharacterResponse -
type DungeonCharacterResponse struct {
	Response
	Data []DungeonCharacterData `json:"data"`
}

// DungeonCharacterRequest -
type DungeonCharacterRequest struct {
	Request
	Data DungeonCharacterData `json:"data"`
}

// DungeonCharacterData -
type DungeonCharacterData struct {
	ID                  string    `json:"id,omitempty"`
	DungeonID           string    `json:"dungeon_id,omitempty"`
	DungeonLocationID   string    `json:"dungeon_location_id,omitempty"`
	Name                string    `json:"name"`
	Strength            int       `json:"strength"`
	Dexterity           int       `json:"dexterity"`
	Intelligence        int       `json:"intelligence"`
	CurrentStrength     int       `json:"current_strength"`
	CurrentDexterity    int       `json:"current_dexterity"`
	CurrentIntelligence int       `json:"current_intelligence"`
	Health              int       `json:"health"`
	Fatigue             int       `json:"fatigue"`
	Coins               int       `json:"coins,omitempty"`
	ExperiencePoints    int       `json:"experience_points"`
	AttributePoints     int       `json:"attribute_points"`
	CreatedAt           time.Time `json:"created_at,omitempty"`
	UpdatedAt           time.Time `json:"updated_at,omitempty"`
}
