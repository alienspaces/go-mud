package record

import (
	"gitlab.com/alienspaces/go-mud/backend/core/repository"
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
