package model

import (
	"fmt"

	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

// validateActionObjectRec - validates creating and updating an action object record
func (m *Model) validateActionObjectRec(rec *record.ActionObject) error {

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

// validateDeleteActionObjectRec - validates it is okay to delete an action object record
func (m *Model) validateDeleteActionObjectRec(recID string) error {

	return nil
}
