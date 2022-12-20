package model

import (
	"fmt"

	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

// validateMonsterObjectRec - validates creating and updating a monster object record
func (m *Model) validateMonsterObjectRec(rec *record.MonsterObject) error {

	if rec.MonsterID == "" {
		return fmt.Errorf("failed validation, MonsterID is empty")
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

// validateDeleteMonsterObjectRec - validates it is okay to delete a monster object record
func (m *Model) validateDeleteMonsterObjectRec(recID string) error {

	return nil
}
