package model

import (
	"fmt"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// ValidateMonsterRec - validates creating and updating a game record
func (m *Model) ValidateMonsterRec(rec *record.Monster) error {

	// New monster
	if rec.ID == "" {
		if rec.DungeonID == "" {
			return fmt.Errorf("failed validation, DungeonID is empty")
		}
		if rec.LocationID == "" {
			return fmt.Errorf("failed validation, LocationID is empty")
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

// ValidateDeleteMonsterRec - validates it is okay to delete a game record
func (m *Model) ValidateDeleteMonsterRec(recID string) error {

	return nil
}
