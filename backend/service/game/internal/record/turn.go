package record

import (
	"database/sql"

	"gitlab.com/alienspaces/go-mud/backend/core/repository"
)

type Turn struct {
	DungeonInstanceID string       `db:"dungeon_instance_id"`
	TurnCount         int          `db:"turn_count"`
	IncrementedAt     sql.NullTime `db:"incremented_at"`
	repository.Record
}
