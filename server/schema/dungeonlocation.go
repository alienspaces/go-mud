package schema

import (
	"time"
)

// LocationResponse -
type LocationResponse struct {
	Response
	Data []LocationData `json:"data"`
}

// LocationRequest -
type LocationRequest struct {
	Request
	Data LocationData `json:"data"`
}

// LocationData -
type LocationData struct {
	ID                  string    `json:"id,omitempty"`
	DungeonID           string    `json:"dungeon_id,omitempty"`
	Name                string    `json:"name"`
	Description         string    `json:"description"`
	Default             bool      `json:"default"`
	NorthLocationID     string    `json:"north_dungeon_location_id,omitempty"`
	NorthEastLocationID string    `json:"northeast_dungeon_location_id,omitempty"`
	EastLocationID      string    `json:"east_dungeon_location_id,omitempty"`
	SouthEastLocationID string    `json:"southeast_dungeon_location_id,omitempty"`
	SouthLocationID     string    `json:"south_dungeon_location_id,omitempty"`
	SouthWestLocationID string    `json:"southwest_dungeon_location_id,omitempty"`
	WestLocationID      string    `json:"west_dungeon_location_id,omitempty"`
	NorthWestLocationID string    `json:"northwest_dungeon_location_id,omitempty"`
	UpLocationID        string    `json:"up_dungeon_location_id,omitempty"`
	DownLocationID      string    `json:"down_dungeon_location_id,omitempty"`
	CreatedAt           time.Time `json:"created_at,omitempty"`
	UpdatedAt           time.Time `json:"updated_at,omitempty"`
}
