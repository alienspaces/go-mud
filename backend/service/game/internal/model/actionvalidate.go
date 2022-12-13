package model

import (
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

// validateActionRec - validates creating and updating a game record
func (m *Model) validateActionRec(rec *record.Action) error {

	return nil
}

// validateDeleteActionRec - validates it is okay to delete a game record
func (m *Model) validateDeleteActionRec(recID string) error {

	return nil
}
