package model

import (
	"fmt"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// ValidateDungeonActionMonsterRec - validates creating and updating a game record
func (m *Model) ValidateDungeonActionMonsterRec(rec *record.DungeonActionMonster) error {

	if rec.RecordType == "" {
		return fmt.Errorf("failed validation, RecordType is empty")
	}
	if rec.DungeonActionID == "" {
		return fmt.Errorf("failed validation, DungeonActionID is empty")
	}
	if rec.DungeonLocationID == "" {
		return fmt.Errorf("failed validation, DungeonLocationID is empty")
	}
	if rec.DungeonMonsterID == "" {
		return fmt.Errorf("failed validation, DungeonMonsterID is empty")
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

// ValidateDeleteDungeonActionMonsterRec - validates it is okay to delete a game record
func (m *Model) ValidateDeleteDungeonActionMonsterRec(recID string) error {

	return nil
}
