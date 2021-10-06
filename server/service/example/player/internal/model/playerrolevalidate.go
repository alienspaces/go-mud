package model

import (
	"gitlab.com/alienspaces/go-boilerplate/server/service/player/internal/record"
)

// ValidatePlayerRoleRec - validates creating and updating a player record
func (m *Model) ValidatePlayerRoleRec(rec *record.PlayerRole) error {

	return nil
}

// ValidateDeletePlayerRoleRec - validates it is okay to delete a player record
func (m *Model) ValidateDeletePlayerRoleRec(recID string) error {

	return nil
}
