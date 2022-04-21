package record

import (
	"database/sql"

	"gitlab.com/alienspaces/go-mud/server/core/repository"
)

type Object struct {
	DungeonID           string         `db:"dungeon_id"`
	LocationID          sql.NullString `db:"dungeon_location_id"`
	CharacterID         sql.NullString `db:"character_id"`
	MonsterID           sql.NullString `db:"monster_id"`
	Name                string         `db:"name"`
	Description         string         `db:"description"`
	DescriptionDetailed string         `db:"description_detailed"`
	IsStashed           bool           `db:"is_stashed"`
	IsEquipped          bool           `db:"is_equipped"`
	repository.Record
}

type ObjectInstance struct {
	ObjectID            string         `db:"object_id"`
	DungeonInstanceID   string         `db:"dungeon_instance_id"`
	LocationInstanceID  sql.NullString `db:"location_instance_id"`
	CharacterInstanceID sql.NullString `db:"character_instance_id"`
	MonsterInstanceID   sql.NullString `db:"monster_instance_id"`
	IsStashed           bool           `db:"is_stashed"`
	IsEquipped          bool           `db:"is_equipped"`
	repository.Record
}

type ObjectInstanceView struct {
	ObjectID            string         `db:"object_id"`
	DungeonInstanceID   string         `db:"dungeon_instance_id"`
	LocationInstanceID  sql.NullString `db:"location_instance_id"`
	CharacterInstanceID sql.NullString `db:"character_instance_id"`
	MonsterInstanceID   sql.NullString `db:"monster_instance_id"`
	Name                string         `db:"name"`
	Description         string         `db:"description"`
	DescriptionDetailed string         `db:"description_detailed"`
	IsStashed           bool           `db:"is_stashed"`
	IsEquipped          bool           `db:"is_equipped"`
	repository.Record
}
