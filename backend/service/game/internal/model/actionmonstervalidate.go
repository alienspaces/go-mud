package model

import (
	"fmt"

	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

// validateActionMonsterRec - validates creating and updating an action monster record
func (m *Model) validateActionMonsterRec(rec *record.ActionMonster) error {

	if rec.RecordType == "" {
		return fmt.Errorf("failed validation, RecordType is empty")
	}
	if rec.ActionID == "" {
		return fmt.Errorf("failed validation, ActionID is empty")
	}
	if rec.LocationInstanceID == "" {
		return fmt.Errorf("failed validation, LocationInstanceID is empty")
	}
	if rec.MonsterInstanceID == "" {
		return fmt.Errorf("failed validation, MonsterInstanceID is empty")
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

	return nil
}

// validateDeleteActionMonsterRec - validates it is okay to delete an action monster record
func (m *Model) validateDeleteActionMonsterRec(recID string) error {

	return nil
}
