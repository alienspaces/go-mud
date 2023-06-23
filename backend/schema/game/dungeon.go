package schema

import (
	"time"

	"gitlab.com/alienspaces/go-mud/backend/schema"
)

// DungeonResponse -
type DungeonResponse struct {
	schema.Response
	Data []DungeonData `json:"data"`
}

// DungeonRequest -
type DungeonRequest struct {
	schema.Request
	Data DungeonData `json:"data"`
}

// DungeonData -
type DungeonData struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

// DungeonEnterRequest
type DungeonEnterRequest struct {
	schema.Request
	Data DungeonEnterData `json:"data"`
}

// DungeonEnterData
type DungeonEnterData struct {
	CharacterID string `json:"character_id,omitempty"`
}

// DungeonExitRequest
type DungeonExitRequest struct {
	schema.Request
	Data DungeonExitData `json:"data"`
}

// DungeonExitData
type DungeonExitData struct {
	CharacterID string `json:"character_id,omitempty"`
}
