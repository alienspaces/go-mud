package model

import (
	"fmt"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// validateObjectRec - validates creating and updating a game record
func (m *Model) validateObjectRec(rec *record.Object) error {

	if rec.Name == "" {
		return fmt.Errorf("failed validation, Name is empty")
	}

	return nil
}

// validateDeleteObjectRec - validates it is okay to delete a game record
func (m *Model) validateDeleteObjectRec(recID string) error {

	return nil
}
