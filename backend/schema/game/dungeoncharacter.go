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
	ID                  string                        `json:"id,omitempty"`
	Name                string                        `json:"name"`
	Strength            int                           `json:"strength"`
	Dexterity           int                           `json:"dexterity"`
	Intelligence        int                           `json:"intelligence"`
	CurrentStrength     int                           `json:"current_strength"`
	CurrentDexterity    int                           `json:"current_dexterity"`
	CurrentIntelligence int                           `json:"current_intelligence"`
	Health              int                           `json:"health"`
	Fatigue             int                           `json:"fatigue"`
	CurrentHealth       int                           `json:"current_health"`
	CurrentFatigue      int                           `json:"current_fatigue"`
	Coins               int                           `json:"coins,omitempty"`
	ExperiencePoints    int                           `json:"experience_points"`
	AttributePoints     int                           `json:"attribute_points"`
	Dungeon             *DungeonCharacterDungeonData  `json:"dungeon,omitempty"`
	Location            *DungeonCharacterLocationData `json:"location,omitempty"`
	CreatedAt           time.Time                     `json:"created_at,omitempty"`
	UpdatedAt           time.Time                     `json:"updated_at,omitempty"`
}

type DungeonCharacterDungeonData struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type DungeonCharacterLocationData struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
