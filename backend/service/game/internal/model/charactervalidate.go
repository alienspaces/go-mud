package model

import (
	"fmt"

	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

// validateCharacterRec - validates creating and updating a character record
func (m *Model) validateCharacterRec(rec *record.Character) error {

	// New character
	if rec.ID == "" {
		if rec.Strength+rec.Intelligence+rec.Dexterity > defaultAttributePoints {
			return fmt.Errorf("new character attributes exceeds allowed maximum of %d", defaultAttributePoints)
		}

		if rec.Name == "" {
			return fmt.Errorf("failed validation, Name is empty")
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

// validateDeleteCharacterRec - validates it is okay to delete a character record
func (m *Model) validateDeleteCharacterRec(recID string) error {

	return nil
}
