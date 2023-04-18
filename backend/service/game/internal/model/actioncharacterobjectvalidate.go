package model

import (
	"fmt"

	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

// validateActionCharacterObjectRec - validates creating and updating an action character object record
func (m *Model) validateActionCharacterObjectRec(rec *record.ActionCharacterObject) error {

	if rec.ActionCharacterID == "" {
		return fmt.Errorf("failed validation, ActionCharacterID is empty")
	}
	if rec.ObjectInstanceID == "" {
		return fmt.Errorf("failed validation, ObjectInstanceID is empty")
	}
	if rec.IsStashed == rec.IsEquipped {
		return fmt.Errorf("failed validation, IsStashed equals IsEquipped")
	}

	return nil
}

// validateDeleteActionCharacterObjectRec - validates it is okay to delete an action character object record
func (m *Model) validateDeleteActionCharacterObjectRec(recID string) error {

	return nil
}
