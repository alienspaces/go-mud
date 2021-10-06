package schema

import "time"

// AuthRefreshResponse -
type AuthRefreshResponse struct {
	Response
	Data []AuthRefreshData `json:"data"`
}

// AuthRefreshRequest -
type AuthRefreshRequest struct {
	Request
	Data AuthRefreshData `json:"data"`
}

// AuthRefreshData -
type AuthRefreshData struct {
	PlayerID    string    `json:"account_id,omitempty"`
	PlayerName  string    `json:"account_name,omitempty"`
	PlayerEmail string    `json:"account_email,omitempty"`
	Token        string    `json:"token,omitempty"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
	UpdatedAt    time.Time `json:"updated_at,omitempty"`
}
