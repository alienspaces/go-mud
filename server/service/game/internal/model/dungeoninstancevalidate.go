package model

import (
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// ValidateDungeonInstanceRec - validates creating and updating a game record
func (m *Model) ValidateDungeonInstanceRec(rec *record.DungeonInstance) error {

	return nil
}

// ValidateDeleteDungeonInstanceRec - validates it is okay to delete a game record
func (m *Model) ValidateDeleteDungeonInstanceRec(recID string) error {

	return nil
}
