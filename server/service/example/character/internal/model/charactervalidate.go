package model

import (
	"fmt"

	"gitlab.com/alienspaces/go-boilerplate/server/service/character/internal/record"
)

// ValidateCharacterRec - validates creating and updating a mage record
func (m *Model) ValidateCharacterRec(rec *record.Character) error {

	// required fields
	if rec.PlayerID == "" {
		return fmt.Errorf("property PlayerID is required")
	}
	if rec.Name == "" {
		return fmt.Errorf("property Name is required")
	}

	return nil
}

// ValidateDeleteCharacterRec - validates it is okay to delete a mage record
func (m *Model) ValidateDeleteCharacterRec(recID string) error {

	return nil
}
