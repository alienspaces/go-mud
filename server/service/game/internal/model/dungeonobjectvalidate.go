package model

import (
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// ValidateDungeonObjectRec - validates creating and updating a game record
func (m *Model) ValidateDungeonObjectRec(rec *record.DungeonObject) error {

	return nil
}

// ValidateDeleteDungeonObjectRec - validates it is okay to delete a game record
func (m *Model) ValidateDeleteDungeonObjectRec(recID string) error {

	return nil
}
