package model

import (
	"fmt"

	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

// validateMonsterRec - validates creating and updating a game record
func (m *Model) validateMonsterRec(rec *record.Monster) error {

	// New monster
	if rec.ID == "" {
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

	if rec.Name == "" {
		return fmt.Errorf("failed validation, Name is empty")
	}

	return nil
}

// validateDeleteMonsterRec - validates it is okay to delete a game record
func (m *Model) validateDeleteMonsterRec(recID string) error {

	return nil
}
