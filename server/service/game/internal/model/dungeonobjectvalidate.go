package model

import (
	"fmt"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// ValidateDungeonObjectRec - validates creating and updating a game record
func (m *Model) ValidateDungeonObjectRec(rec *record.DungeonObject) error {

	if rec.ID == "" {
		if rec.DungeonID == "" {
			return fmt.Errorf("failed validation, DungeonID is empty")
		}
		if !rec.DungeonLocationID.Valid {
			return fmt.Errorf("failed validation, DungeonLocationID is empty")
		}
		if rec.Name == "" {
			return fmt.Errorf("failed validation, Name is empty")
		}
	}

	return nil
}

// ValidateDeleteDungeonObjectRec - validates it is okay to delete a game record
func (m *Model) ValidateDeleteDungeonObjectRec(recID string) error {

	return nil
}
