package model

import (
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// ValidateActionRec - validates creating and updating a game record
func (m *Model) ValidateActionRec(rec *record.Action) error {

	return nil
}

// ValidateDeleteActionRec - validates it is okay to delete a game record
func (m *Model) ValidateDeleteActionRec(recID string) error {

	return nil
}
