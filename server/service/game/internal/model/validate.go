package model

import (
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// ValidateGameRec - validates creating and updating a game record
func (m *Model) ValidateGameRec(rec *record.Game) error {

	return nil
}

// ValidateDeleteGameRec - validates it is okay to delete a game record
func (m *Model) ValidateDeleteGameRec(recID string) error {

	return nil
}
