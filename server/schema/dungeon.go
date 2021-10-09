package schema

import (
	"time"
)

// DungeonResponse -
type DungeonResponse struct {
	Response
	Data []DungeonData `json:"data"`
}

// DungeonRequest -
type DungeonRequest struct {
	Request
	Data DungeonData `json:"data"`
}

// DungeonData -
type DungeonData struct {
	ID          string    `json:"id,omitempty"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}
