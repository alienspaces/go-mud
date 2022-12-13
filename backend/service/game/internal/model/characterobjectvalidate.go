package model

import (
	"fmt"

	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

// validateCharacterObjectRec - validates creating and updating a game record
func (m *Model) validateCharacterObjectRec(rec *record.CharacterObject) error {

	if rec.CharacterID == "" {
		return fmt.Errorf("failed validation, CharacterID is empty")
	}
	if rec.ObjectID == "" {
		return fmt.Errorf("failed validation, ObjectID is empty")
	}
	if rec.IsStashed && rec.IsEquipped {
		return fmt.Errorf("failed validation, object cannot be stashed and equipped")
	}
	if !rec.IsStashed && !rec.IsEquipped {
		return fmt.Errorf("failed validation, object must be stashed or equipped")
	}

	return nil
}

// validateDeleteCharacterObjectRec - validates it is okay to delete a game record
func (m *Model) validateDeleteCharacterObjectRec(recID string) error {

	return nil
}
