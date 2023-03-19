package model

import (
	"fmt"

	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

// validateActionMonsterObjectRec - validates creating and updating an action monster object record
func (m *Model) validateActionMonsterObjectRec(rec *record.ActionMonsterObject) error {

	if rec.ActionMonsterID == "" {
		return fmt.Errorf("failed validation, ActionMonsterID is empty")
	}
	if rec.ObjectInstanceID == "" {
		return fmt.Errorf("failed validation, ObjectInstanceID is empty")
	}
	if rec.IsStashed == rec.IsEquipped {
		return fmt.Errorf("failed validation, IsStashed equals IsEquipped")
	}

	return nil
}

// validateDeleteActionMonsterObjectRec - validates it is okay to delete an action monster object record
func (m *Model) validateDeleteActionMonsterObjectRec(recID string) error {

	return nil
}
