package model

import (
	"fmt"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// ValidateCharacterInstanceRec - validates creating and updating a game record
func (m *Model) ValidateCharacterInstanceRec(rec *record.CharacterInstance) error {

	// New character
	if rec.ID == "" {
		if rec.Strength+rec.Intelligence+rec.Dexterity > defaultAttributePoints {
			return fmt.Errorf("new character attributes exceeds allowed maximum of %d", defaultAttributePoints)
		}

		if rec.Strength == 0 {
			return fmt.Errorf("failed validation, Strength is empty")
		}
		if rec.Dexterity == 0 {
			return fmt.Errorf("failed validation, Dexterity is empty")
		}
		if rec.Intelligence == 0 {
			return fmt.Errorf("failed validation, Intelligence is empty")
		}
		if rec.Health == 0 {
			return fmt.Errorf("failed validation, Health is empty")
		}
		if rec.Fatigue == 0 {
			return fmt.Errorf("failed validation, Fatigue is empty")
		}
	}

	return nil
}

// ValidateDeleteCharacterInstanceRec - validates it is okay to delete a game record
func (m *Model) ValidateDeleteCharacterInstanceRec(recID string) error {

	return nil
}
