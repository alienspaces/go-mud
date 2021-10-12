package model

import (
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// ValidateDungeonActionMonsterRec - validates creating and updating a game record
func (m *Model) ValidateDungeonActionMonsterRec(rec *record.DungeonActionMonster) error {

	return nil
}

// ValidateDeleteDungeonActionMonsterRec - validates it is okay to delete a game record
func (m *Model) ValidateDeleteDungeonActionMonsterRec(recID string) error {

	return nil
}
