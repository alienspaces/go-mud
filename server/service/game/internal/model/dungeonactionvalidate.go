package model

import (
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// ValidateDungeonActionRec - validates creating and updating a game record
func (m *Model) ValidateDungeonActionRec(rec *record.DungeonAction) error {

	return nil
}

// ValidateDeleteDungeonActionRec - validates it is okay to delete a game record
func (m *Model) ValidateDeleteDungeonActionRec(recID string) error {

	return nil
}
