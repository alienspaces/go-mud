package model

import (
	"fmt"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// ValidateDungeonActionMonsterObjectRec - validates creating and updating a game record
func (m *Model) ValidateDungeonActionMonsterObjectRec(rec *record.DungeonActionMonsterObject) error {

	if rec.DungeonActionID == "" {
		return fmt.Errorf("failed validation, DungeonActionID is empty")
	}
	if rec.DungeonMonsterID == "" {
		return fmt.Errorf("failed validation, DungeonMonsterID is empty")
	}
	if rec.DungeonObjectID == "" {
		return fmt.Errorf("failed validation, DungeonObjectID is empty")
	}
	if rec.IsStashed == rec.IsEquipped {
		return fmt.Errorf("failed validation, IsStashed equals IsEquipped")
	}

	return nil
}

// ValidateDeleteDungeonActionMonsterObjectRec - validates it is okay to delete a game record
func (m *Model) ValidateDeleteDungeonActionMonsterObjectRec(recID string) error {

	return nil
}
