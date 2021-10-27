package model

import (
	"fmt"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// ValidateDungeonCharacterRec - validates creating and updating a game record
func (m *Model) ValidateDungeonCharacterRec(rec *record.DungeonCharacter) error {

	// New character
	if rec.ID == "" {
		if rec.Strength+rec.Intelligence+rec.Dexterity > defaultAttributePoints {
			return fmt.Errorf("new character attributes exceeds allowed maximum of %d", defaultAttributePoints)
		}
	}

	return nil
}

// ValidateDeleteDungeonCharacterRec - validates it is okay to delete a game record
func (m *Model) ValidateDeleteDungeonCharacterRec(recID string) error {

	return nil
}
