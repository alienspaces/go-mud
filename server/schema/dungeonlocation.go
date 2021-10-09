package schema

import (
	"time"
)

// DungeonLocationResponse -
type DungeonLocationResponse struct {
	Response
	Data []DungeonLocationData `json:"data"`
}

// DungeonLocationRequest -
type DungeonLocationRequest struct {
	Request
	Data DungeonLocationData `json:"data"`
}

// DungeonLocationData -
type DungeonLocationData struct {
	ID                         string    `json:"id,omitempty"`
	DungeonID                  string    `json:"dungeon_id,omitempty"`
	Name                       string    `json:"name"`
	Description                string    `json:"description"`
	Default                    bool      `json:"default"`
	NorthDungeonLocationID     string    `json:"north_dungeon_location_id,omitempty"`
	NorthEastDungeonLocationID string    `json:"northeast_dungeon_location_id,omitempty"`
	EastDungeonLocationID      string    `json:"east_dungeon_location_id,omitempty"`
	SouthEastDungeonLocationID string    `json:"southeast_dungeon_location_id,omitempty"`
	SouthDungeonLocationID     string    `json:"south_dungeon_location_id,omitempty"`
	SouthWestDungeonLocationID string    `json:"southwest_dungeon_location_id,omitempty"`
	WestDungeonLocationID      string    `json:"west_dungeon_location_id,omitempty"`
	NorthWestDungeonLocationID string    `json:"northwest_dungeon_location_id,omitempty"`
	UpDungeonLocationID        string    `json:"up_dungeon_location_id,omitempty"`
	DownDungeonLocationID      string    `json:"down_dungeon_location_id,omitempty"`
	CreatedAt                  time.Time `json:"created_at,omitempty"`
	UpdatedAt                  time.Time `json:"updated_at,omitempty"`
}
