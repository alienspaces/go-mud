package model

import (
	"fmt"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// ValidateLocationRec - validates creating and updating a game record
func (m *Model) ValidateLocationRec(rec *record.Location) error {

	if rec.Name == "" {
		return fmt.Errorf("failed validation, Name is empty")
	}
	if rec.Description == "" {
		return fmt.Errorf("failed validation, Name is empty")
	}

	return nil
}

// ValidateDeleteLocationRec - validates it is okay to delete a game record
func (m *Model) ValidateDeleteLocationRec(recID string) error {

	return nil
}
