package record

import (
	"database/sql"

	"gitlab.com/alienspaces/go-mud/server/core/repository"
)

type DungeonLocation struct {
	DungeonID                  string         `db:"dungeon_id"`
	Name                       string         `db:"name"`
	Description                string         `db:"description"`
	Default                    bool           `db:"default"`
	NorthDungeonLocationID     sql.NullString `db:"north_dungeon_location_id"`
	NortheastDungeonLocationID sql.NullString `db:"northeast_dungeon_location_id"`
	EastDungeonLocationID      sql.NullString `db:"east_dungeon_location_id"`
	SoutheastDungeonLocationID sql.NullString `db:"southeast_dungeon_location_id"`
	SouthDungeonLocationID     sql.NullString `db:"south_dungeon_location_id"`
	SouthwestDungeonLocationID sql.NullString `db:"southwest_dungeon_location_id"`
	WestDungeonLocationID      sql.NullString `db:"west_dungeon_location_id"`
	NorthwestDungeonLocationID sql.NullString `db:"northwest_dungeon_location_id"`
	UpDungeonLocationID        sql.NullString `db:"up_dungeon_location_id"`
	DownDungeonLocationID      sql.NullString `db:"down_dungeon_location_id"`
	repository.Record
}

type DungeonLocationInstance struct {
	DungeonLocationID string `db:"dungeon_location_id"`
	repository.Record
}
