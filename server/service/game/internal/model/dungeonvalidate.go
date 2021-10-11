package model

import (
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// ValidateDungeonRec - validates creating and updating a game record
func (m *Model) ValidateDungeonRec(rec *record.Dungeon) error {

	return nil
}

// ValidateDeleteDungeonRec - validates it is okay to delete a game record
func (m *Model) ValidateDeleteDungeonRec(recID string) error {

	return nil
}
