package model

import (
	"fmt"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// ValidateDungeonActionCharacterObjectRec - validates creating and updating a game record
func (m *Model) ValidateDungeonActionCharacterObjectRec(rec *record.DungeonActionCharacterObject) error {

	if rec.RecordType == "" {
		return fmt.Errorf("failed validation, RecordType is empty")
	}
	if rec.DungeonActionID == "" {
		return fmt.Errorf("failed validation, DungeonActionID is empty")
	}
	if rec.DungeonCharacterID == "" {
		return fmt.Errorf("failed validation, DungeonCharacterID is empty")
	}
	if rec.DungeonObjectID == "" {
		return fmt.Errorf("failed validation, DungeonObjectID is empty")
	}

	return nil
}

// ValidateDeleteDungeonActionCharacterObjectRec - validates it is okay to delete a game record
func (m *Model) ValidateDeleteDungeonActionCharacterObjectRec(recID string) error {

	return nil
}
