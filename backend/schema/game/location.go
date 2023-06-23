package schema

import (
	"time"

	"gitlab.com/alienspaces/go-mud/backend/schema"
)

// LocationResponse -
type LocationResponse struct {
	schema.Response
	Data []LocationData `json:"data"`
}

// LocationRequest -
type LocationRequest struct {
	schema.Request
	Data LocationData `json:"data"`
}

// LocationData -
type LocationData struct {
	ID                  string    `json:"id"`
	DungeonID           string    `json:"dungeon_id"`
	Name                string    `json:"name"`
	Description         string    `json:"description"`
	Default             bool      `json:"default"`
	NorthLocationID     string    `json:"north_location_id,omitempty"`
	NorthEastLocationID string    `json:"northeast_location_id,omitempty"`
	EastLocationID      string    `json:"east_location_id,omitempty"`
	SouthEastLocationID string    `json:"southeast_location_id,omitempty"`
	SouthLocationID     string    `json:"south_location_id,omitempty"`
	SouthWestLocationID string    `json:"southwest_location_id,omitempty"`
	WestLocationID      string    `json:"west_location_id,omitempty"`
	NorthWestLocationID string    `json:"northwest_location_id,omitempty"`
	UpLocationID        string    `json:"up_location_id,omitempty"`
	DownLocationID      string    `json:"down_location_id,omitempty"`
	CreatedAt           time.Time `json:"created_at,omitempty"`
	UpdatedAt           time.Time `json:"updated_at,omitempty"`
}
