package schema

import "time"

// PlayerRequest -
type PlayerRequest struct {
	Request
	Data PlayerData `json:"data"`
}

// PlayerResponse -
type PlayerResponse struct {
	Response
	Data []PlayerData `json:"data"`
}

// PlayerData -
type PlayerData struct {
	ID                string    `json:"id,omitempty"`
	Name              string    `json:"name,omitempty"`
	Email             string    `json:"email,omitempty"`
	Provider          string    `json:"provider,omitempty"`
	ProviderAccountID string    `json:"provider_account_id,omitempty"`
	CreatedAt         time.Time `json:"created_at,omitempty"`
	UpdatedAt         time.Time `json:"updated_at,omitempty"`
}

// PlayerCharacterRequest -
type PlayerCharacterRequest struct {
	Request
	Data PlayerCharacterData `json:"data"`
}

// PlayerCharacterResponse -
type PlayerCharacterResponse struct {
	Request
	Data []PlayerCharacterData `json:"data"`
}

// PlayerCharacterData -
type PlayerCharacterData struct {
	ID          string    `json:"id,omitempty"`
	PlayerID    string    `json:"account_id,omitempty"`
	CharacterID string    `json:"entity_id,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}
