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
	DungeonID          string    `json:"dungeon_id"`
	DungeonName        string    `json:"dungeon_name"`
	DungeonDescription string    `json:"dungeon_description"`
	DungeonCreatedAt   time.Time `json:"dungeon_created_at,omitempty"`
	DungeonUpdatedAt   time.Time `json:"dungeon_updated_at,omitempty"`
}

// DungeonEnterRequest
type DungeonEnterRequest struct {
	Request
	Data DungeonEnterData `json:"data"`
}

// DungeonEnterData
type DungeonEnterData struct {
	CharacterID string `json:"character_id,omitempty"`
}

// DungeonExitRequest
type DungeonExitRequest struct {
	Request
	Data DungeonExitData `json:"data"`
}

// DungeonExitData
type DungeonExitData struct {
	CharacterID string `json:"character_id,omitempty"`
}
