package model

import (
	"fmt"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// ValidateMonsterInstanceRec - validates creating and updating a game record
func (m *Model) ValidateMonsterInstanceRec(rec *record.MonsterInstance) error {

	// New monster
	if rec.ID == "" {
		if rec.DungeonInstanceID == "" {
			return fmt.Errorf("failed validation, DungeonID is empty")
		}
		if rec.LocationInstanceID == "" {
			return fmt.Errorf("failed validation, LocationID is empty")
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

// ValidateDeleteMonsterInstanceRec - validates it is okay to delete a game record
func (m *Model) ValidateDeleteMonsterInstanceRec(recID string) error {

	return nil
}
