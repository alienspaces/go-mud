package model

import (
	"fmt"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// ValidateActionMonsterObjectRec - validates creating and updating a game record
func (m *Model) ValidateActionMonsterObjectRec(rec *record.ActionMonsterObject) error {

	if rec.ActionID == "" {
		return fmt.Errorf("failed validation, ActionID is empty")
	}
	if rec.MonsterInstanceID == "" {
		return fmt.Errorf("failed validation, MonsterInstanceID is empty")
	}
	if rec.ObjectInstanceID == "" {
		return fmt.Errorf("failed validation, ObjectInstanceID is empty")
	}
	if rec.IsStashed == rec.IsEquipped {
		return fmt.Errorf("failed validation, IsStashed equals IsEquipped")
	}

	return nil
}

// ValidateDeleteActionMonsterObjectRec - validates it is okay to delete a game record
func (m *Model) ValidateDeleteActionMonsterObjectRec(recID string) error {

	return nil
}
