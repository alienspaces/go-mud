package model

import (
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// ValidateDungeonActionObjectRec - validates creating and updating a game record
func (m *Model) ValidateDungeonActionObjectRec(rec *record.DungeonActionObject) error {

	return nil
}

// ValidateDeleteDungeonActionObjectRec - validates it is okay to delete a game record
func (m *Model) ValidateDeleteDungeonActionObjectRec(recID string) error {

	return nil
}
