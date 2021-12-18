package model

import (
	"fmt"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// ValidateDungeonActionObjectRec - validates creating and updating a game record
func (m *Model) ValidateDungeonActionObjectRec(rec *record.DungeonActionObject) error {

	if rec.RecordType == "" {
		return fmt.Errorf("failed validation, RecordType is empty")
	}
	if rec.DungeonActionID == "" {
		return fmt.Errorf("failed validation, DungeonActionID is empty")
	}
	if rec.DungeonLocationID == "" {
		return fmt.Errorf("failed validation, DungeonLocationID is empty")
	}
	if rec.DungeonObjectID == "" {
		return fmt.Errorf("failed validation, DungeonObjectID is empty")
	}

	return nil
}

// ValidateDeleteDungeonActionObjectRec - validates it is okay to delete a game record
func (m *Model) ValidateDeleteDungeonActionObjectRec(recID string) error {

	return nil
}
