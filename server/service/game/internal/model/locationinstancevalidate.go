package model

import (
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// ValidateLocationInstanceRec - validates creating and updating a game record
func (m *Model) ValidateLocationInstanceRec(rec *record.LocationInstance) error {

	return nil
}

// ValidateDeleteLocationInstanceRec - validates it is okay to delete a game record
func (m *Model) ValidateDeleteLocationInstanceRec(recID string) error {

	return nil
}
