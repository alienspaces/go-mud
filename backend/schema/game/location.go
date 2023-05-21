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
	DungeonID                   string    `json:"dungeon_id"`
	LocationID                  string    `json:"location_id"`
	LocationName                string    `json:"location_name"`
	LocationDescription         string    `json:"location_description"`
	LocationDefault             bool      `json:"location_default"`
	LocationNorthLocationID     string    `json:"location_north_dungeon_location_id,omitempty"`
	LocationNorthEastLocationID string    `json:"location_northeast_dungeon_location_id,omitempty"`
	LocationEastLocationID      string    `json:"location_east_dungeon_location_id,omitempty"`
	LocationSouthEastLocationID string    `json:"location_southeast_dungeon_location_id,omitempty"`
	LocationSouthLocationID     string    `json:"location_south_dungeon_location_id,omitempty"`
	LocationSouthWestLocationID string    `json:"location_southwest_dungeon_location_id,omitempty"`
	LocationWestLocationID      string    `json:"location_west_dungeon_location_id,omitempty"`
	LocationNorthWestLocationID string    `json:"location_northwest_dungeon_location_id,omitempty"`
	LocationUpLocationID        string    `json:"location_up_dungeon_location_id,omitempty"`
	LocationDownLocationID      string    `json:"location_down_dungeon_location_id,omitempty"`
	LocationCreatedAt           time.Time `json:"location_created_at,omitempty"`
	LocationUpdatedAt           time.Time `json:"location_updated_at,omitempty"`
}
