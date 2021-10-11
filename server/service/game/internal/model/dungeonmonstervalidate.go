package model

import (
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// ValidateDungeonMonsterRec - validates creating and updating a game record
func (m *Model) ValidateDungeonMonsterRec(rec *record.DungeonMonster) error {

	return nil
}

// ValidateDeleteDungeonMonsterRec - validates it is okay to delete a game record
func (m *Model) ValidateDeleteDungeonMonsterRec(recID string) error {

	return nil
}
