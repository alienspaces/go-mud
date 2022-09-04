package schema

import (
	"time"
)

// DungeonCharacterResponse -
type DungeonCharacterResponse struct {
	Response
	Data []DungeonCharacterData `json:"data"`
}

// DungeonCharacterData -
type DungeonCharacterData struct {
	ID                  string    `json:"id,omitempty"`
	DungeonID           string    `json:"dungeon_id,omitempty"`
	DungeonName         string    `json:"dungeon_name,omitempty"`
	DungeonDescription  string    `json:"dungeon_description,omitempty"`
	LocationID          string    `json:"dungeon_location_id,omitempty"`
	LocationName        string    `json:"location_name,omitempty"`
	LocationDescription string    `json:"location_description,omitempty"`
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
