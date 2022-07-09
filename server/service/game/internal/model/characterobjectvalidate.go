package model

import (
	"fmt"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// ValidateCharacterObjectRec - validates creating and updating a game record
func (m *Model) ValidateCharacterObjectRec(rec *record.CharacterObject) error {

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

// ValidateDeleteCharacterObjectRec - validates it is okay to delete a game record
func (m *Model) ValidateDeleteCharacterObjectRec(recID string) error {

	return nil
}
