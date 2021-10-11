package model

import (
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// ValidateDungeonCharacterRec - validates creating and updating a game record
func (m *Model) ValidateDungeonCharacterRec(rec *record.DungeonCharacter) error {

	return nil
}

// ValidateDeleteDungeonCharacterRec - validates it is okay to delete a game record
func (m *Model) ValidateDeleteDungeonCharacterRec(recID string) error {

	return nil
}
