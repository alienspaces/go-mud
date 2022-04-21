package model

import (
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// ValidateLocationRec - validates creating and updating a game record
func (m *Model) ValidateLocationRec(rec *record.Location) error {

	return nil
}

// ValidateDeleteLocationRec - validates it is okay to delete a game record
func (m *Model) ValidateDeleteLocationRec(recID string) error {

	return nil
}
