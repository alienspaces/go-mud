package model

import (
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// ValidateDungeonActionCharacterRec - validates creating and updating a game record
func (m *Model) ValidateDungeonActionCharacterRec(rec *record.DungeonActionCharacter) error {

	return nil
}

// ValidateDeleteDungeonActionCharacterRec - validates it is okay to delete a game record
func (m *Model) ValidateDeleteDungeonActionCharacterRec(recID string) error {

	return nil
}
