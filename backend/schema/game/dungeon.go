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
	DungeonID          string    `json:"dungeon_id"`
	DungeonName        string    `json:"dungeon_name"`
	DungeonDescription string    `json:"dungeon_description"`
	DungeonCreatedAt   time.Time `json:"dungeon_created_at,omitempty"`
	DungeonUpdatedAt   time.Time `json:"dungeon_updated_at,omitempty"`
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
