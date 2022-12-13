package model

import (
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

// validateDungeonRec - validates creating and updating a game record
func (m *Model) validateDungeonRec(rec *record.Dungeon) error {

	return nil
}

// validateDeleteDungeonRec - validates it is okay to delete a game record
func (m *Model) validateDeleteDungeonRec(recID string) error {

	return nil
}
