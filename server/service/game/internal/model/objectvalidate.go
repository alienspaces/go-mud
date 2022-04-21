package model

import (
	"fmt"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// ValidateObjectRec - validates creating and updating a game record
func (m *Model) ValidateObjectRec(rec *record.Object) error {

	if rec.ID == "" {
		if rec.DungeonID == "" {
			return fmt.Errorf("failed validation, DungeonID is empty")
		}
		if !rec.LocationID.Valid &&
			!rec.CharacterID.Valid &&
			!rec.MonsterID.Valid {
			return fmt.Errorf("failed validation, all of LocationID, CharacterID and MonsterID are empty")
		}
		if rec.Name == "" {
			return fmt.Errorf("failed validation, Name is empty")
		}
	}

	return nil
}

// ValidateDeleteObjectRec - validates it is okay to delete a game record
func (m *Model) ValidateDeleteObjectRec(recID string) error {

	return nil
}
