package model

import (
	"fmt"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// ValidateObjectRec - validates creating and updating a game record
func (m *Model) ValidateObjectRec(rec *record.Object) error {

	if rec.Name == "" {
		return fmt.Errorf("failed validation, Name is empty")
	}

	return nil
}

// ValidateDeleteObjectRec - validates it is okay to delete a game record
func (m *Model) ValidateDeleteObjectRec(recID string) error {

	return nil
}
