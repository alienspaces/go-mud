package record

import (
	"gitlab.com/alienspaces/go-mud/server/core/repository"
)

// Repository record types
type Dungeon struct {
	Name        string `db:"name"`
	Description string `db:"description"`
	repository.Record
}

type DungeonInstance struct {
	DungeonID string `db:"dungeon_id"`
	repository.Record
}

type DungeonInstanceView struct {
	DungeonID   string `db:"dungeon_id"`
	Name        string `db:"name"`
	Description string `db:"description"`
	repository.Record
}

// Query record types
type DungeonInstanceCapacity struct {
	DungeonInstanceID             string `db:"dungeon_instance_id"`
	DungeonInstanceCharacterCount int    `db:"dungeon_instance_character_count"`
	DungeonID                     string `db:"dungeon_id"`
	DungeonLocationCount          int    `db:"dungeon_location_count"`
}
