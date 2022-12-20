package record

import (
	"gitlab.com/alienspaces/go-mud/backend/core/repository"
)

type Turn struct {
	DungeonInstanceID string `db:"dungeon_instance_id"`
	Turn              int    `db:"turn"`
	repository.Record
}
