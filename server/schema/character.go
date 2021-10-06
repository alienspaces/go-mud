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
	ID               string    `json:"id,omitempty"`
	PlayerID         string    `json:"player_id,omitempty"`
	Name             string    `json:"name"`
	Avatar           string    `json:"avatar"`
	Strength         int       `json:"strength"`
	Dexterity        int       `json:"dexterity"`
	Intelligence     int       `json:"intelligence"`
	AttributePoints  int64     `json:"attribute_points,omitempty"`
	ExperiencePoints int64     `json:"experience_points,omitempty"`
	Coins            int64     `json:"coins,omitempty"`
	CreatedAt        time.Time `json:"created_at,omitempty"`
	UpdatedAt        time.Time `json:"updated_at,omitempty"`
}
