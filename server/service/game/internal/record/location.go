package record

import (
	"database/sql"

	"gitlab.com/alienspaces/go-mud/server/core/repository"
)

type Location struct {
	DungeonID           string         `db:"dungeon_id"`
	Name                string         `db:"name"`
	Description         string         `db:"description"`
	IsDefault           bool           `db:"is_default"`
	NorthLocationID     sql.NullString `db:"north_location_id"`
	NortheastLocationID sql.NullString `db:"northeast_location_id"`
	EastLocationID      sql.NullString `db:"east_location_id"`
	SoutheastLocationID sql.NullString `db:"southeast_location_id"`
	SouthLocationID     sql.NullString `db:"south_location_id"`
	SouthwestLocationID sql.NullString `db:"southwest_location_id"`
	WestLocationID      sql.NullString `db:"west_location_id"`
	NorthwestLocationID sql.NullString `db:"northwest_location_id"`
	UpLocationID        sql.NullString `db:"up_location_id"`
	DownLocationID      sql.NullString `db:"down_location_id"`
	repository.Record
}

type LocationInstance struct {
	DungeonInstanceID           string         `db:"dungeon_instance_id"`
	LocationID                  string         `db:"location_id"`
	NorthLocationInstanceID     sql.NullString `db:"north_location_instance_id"`
	NortheastLocationInstanceID sql.NullString `db:"northeast_location_instance_id"`
	EastLocationInstanceID      sql.NullString `db:"east_location_instance_id"`
	SoutheastLocationInstanceID sql.NullString `db:"southeast_location_instance_id"`
	SouthLocationInstanceID     sql.NullString `db:"south_location_instance_id"`
	SouthwestLocationInstanceID sql.NullString `db:"southwest_location_instance_id"`
	WestLocationInstanceID      sql.NullString `db:"west_location_instance_id"`
	NorthwestLocationInstanceID sql.NullString `db:"northwest_location_instance_id"`
	UpLocationInstanceID        sql.NullString `db:"up_location_instance_id"`
	DownLocationInstanceID      sql.NullString `db:"down_location_instance_id"`
	repository.Record
}

type LocationInstanceView struct {
	DungeonID                   string         `db:"dungeon_id"`
	LocationID                  string         `db:"location_id"`
	DungeonInstanceID           string         `db:"dungeon_instance_id"`
	Name                        string         `db:"name"`
	Description                 string         `db:"description"`
	IsDefault                   bool           `db:"is_default"`
	NorthLocationInstanceID     sql.NullString `db:"north_location_instance_id"`
	NortheastLocationInstanceID sql.NullString `db:"northeast_location_instance_id"`
	EastLocationInstanceID      sql.NullString `db:"east_location_instance_id"`
	SoutheastLocationInstanceID sql.NullString `db:"southeast_location_instance_id"`
	SouthLocationInstanceID     sql.NullString `db:"south_location_instance_id"`
	SouthwestLocationInstanceID sql.NullString `db:"southwest_location_instance_id"`
	WestLocationInstanceID      sql.NullString `db:"west_location_instance_id"`
	NorthwestLocationInstanceID sql.NullString `db:"northwest_location_instance_id"`
	UpLocationInstanceID        sql.NullString `db:"up_location_instance_id"`
	DownLocationInstanceID      sql.NullString `db:"down_location_instance_id"`
	repository.Record
}
