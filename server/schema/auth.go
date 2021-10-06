package schema

import "time"

// AuthResponse -
type AuthResponse struct {
	Response
	Data []AuthData `json:"data"`
}

// AuthRequest -
type AuthRequest struct {
	Request
	Data AuthData `json:"data"`
}

// AuthData -
type AuthData struct {
	Provider          string    `json:"provider,omitempty"`
	ProviderAccountID string    `json:"provider_account_id,omitempty"`
	ProviderToken     string    `json:"provider_token,omitempty"`
	PlayerID          string    `json:"player_id,omitempty"`
	PlayerName        string    `json:"player_name,omitempty"`
	PlayerEmail       string    `json:"player_email,omitempty"`
	Token             string    `json:"token,omitempty"`
	CreatedAt         time.Time `json:"created_at,omitempty"`
	UpdatedAt         time.Time `json:"updated_at,omitempty"`
}
