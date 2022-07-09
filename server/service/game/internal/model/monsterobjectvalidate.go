package model

import (
	"fmt"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// ValidateMonsterObjectRec - validates creating and updating a game record
func (m *Model) ValidateMonsterObjectRec(rec *record.MonsterObject) error {

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

// ValidateDeleteMonsterObjectRec - validates it is okay to delete a game record
func (m *Model) ValidateDeleteMonsterObjectRec(recID string) error {

	return nil
}
