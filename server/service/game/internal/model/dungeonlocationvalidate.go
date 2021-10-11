package model

import (
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// ValidateDungeonLocationRec - validates creating and updating a game record
func (m *Model) ValidateDungeonLocationRec(rec *record.DungeonLocation) error {

	return nil
}

// ValidateDeleteDungeonLocationRec - validates it is okay to delete a game record
func (m *Model) ValidateDeleteDungeonLocationRec(recID string) error {

	return nil
}
