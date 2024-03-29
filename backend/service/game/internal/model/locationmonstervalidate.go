package model

import (
	"fmt"

	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

// validateLocationMonsterRec - validates creating and updating a location monster record
func (m *Model) validateLocationMonsterRec(rec *record.LocationMonster) error {

	// New monster
	if rec.ID == "" {
		if rec.LocationID == "" {
			return fmt.Errorf("failed validation, LocationID is empty")
		}
		if rec.MonsterID == "" {
			return fmt.Errorf("failed validation, MonsterID is empty")
		}
		if rec.SpawnMinutes < 0 || rec.SpawnMinutes > 60 {
			return fmt.Errorf("failed validation, SpawnMinutes must be greater than or equal to 0 and less than or equal to 60")
		}
		if rec.SpawnPercentChance < 1 || rec.SpawnPercentChance > 100 {
			return fmt.Errorf("failed validation, SpawnPercentChance must be greater than 0 and less than or equal to 100")
		}
	}

	return nil
}

// validateDeleteLocationMonsterRec - validates it is okay to delete a location monster record
func (m *Model) validateDeleteLocationMonsterRec(recID string) error {

	return nil
}
