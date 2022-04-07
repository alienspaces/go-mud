package record

import (
	"gitlab.com/alienspaces/go-mud/server/core/repository"
)

type Dungeon struct {
	Name        string `db:"name"`
	Description string `db:"description"`
	repository.Record
}

type DungeonInstance struct {
	DungeonID string `db:"dungeon_id"`
	repository.Record
}
