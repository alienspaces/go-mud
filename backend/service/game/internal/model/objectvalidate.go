package model

import (
	"fmt"

	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

// validateObjectRec - validates creating and updating an object record
func (m *Model) validateObjectRec(rec *record.Object) error {

	if rec.Name == "" {
		return fmt.Errorf("failed validation, Name is empty")
	}

	return nil
}

// validateDeleteObjectRec - validates it is okay to delete an object record
func (m *Model) validateDeleteObjectRec(recID string) error {

	return nil
}
