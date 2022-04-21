package model

import (
	"fmt"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// ValidateActionCharacterObjectRec - validates creating and updating a game record
func (m *Model) ValidateActionCharacterObjectRec(rec *record.ActionCharacterObject) error {

	if rec.ActionID == "" {
		return fmt.Errorf("failed validation, ActionID is empty")
	}
	if rec.CharacterInstanceID == "" {
		return fmt.Errorf("failed validation, CharacterInstanceID is empty")
	}
	if rec.ObjectInstanceID == "" {
		return fmt.Errorf("failed validation, ObjectInstanceID is empty")
	}
	if rec.IsStashed == rec.IsEquipped {
		return fmt.Errorf("failed validation, IsStashed equals IsEquipped")
	}

	return nil
}

// ValidateDeleteActionCharacterObjectRec - validates it is okay to delete a game record
func (m *Model) ValidateDeleteActionCharacterObjectRec(recID string) error {

	return nil
}
