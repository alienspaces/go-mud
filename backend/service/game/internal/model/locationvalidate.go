package model

import (
	"fmt"

	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

// validateLocationRec - validates creating and updating a game record
func (m *Model) validateLocationRec(rec *record.Location) error {

	if rec.Name == "" {
		return fmt.Errorf("failed validation, Name is empty")
	}
	if rec.Description == "" {
		return fmt.Errorf("failed validation, Description is empty")
	}

	return nil
}

// validateDeleteLocationRec - validates it is okay to delete a game record
func (m *Model) validateDeleteLocationRec(recID string) error {

	return nil
}
