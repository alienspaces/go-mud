package model

import (
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// ValidateLocationObjectRec - validates creating and updating a game record
func (m *Model) ValidateLocationObjectRec(rec *record.LocationObject) error {

	return nil
}

// ValidateDeleteLocationObjectRec - validates it is okay to delete a game record
func (m *Model) ValidateDeleteLocationObjectRec(recID string) error {

	return nil
}
