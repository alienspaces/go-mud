package model

import (
	"fmt"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// ValidateActionObjectRec - validates creating and updating a game record
func (m *Model) ValidateActionObjectRec(rec *record.ActionObject) error {

	if rec.RecordType == "" {
		return fmt.Errorf("failed validation, RecordType is empty")
	}
	if rec.ActionID == "" {
		return fmt.Errorf("failed validation, ActionID is empty")
	}
	if rec.LocationInstanceID == "" {
		return fmt.Errorf("failed validation, LocationID is empty")
	}
	if rec.ObjectInstanceID == "" {
		return fmt.Errorf("failed validation, DungeonObjectID is empty")
	}

	return nil
}

// ValidateDeleteActionObjectRec - validates it is okay to delete a game record
func (m *Model) ValidateDeleteActionObjectRec(recID string) error {

	return nil
}
