package model

import (
	"fmt"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// ValidateObjectInstanceRec - validates creating and updating a game record
func (m *Model) ValidateObjectInstanceRec(rec *record.ObjectInstance) error {

	if rec.ID == "" {
		if rec.DungeonInstanceID == "" {
			return fmt.Errorf("failed validation, DungeonInstanceID is empty")
		}
		if !rec.LocationInstanceID.Valid &&
			!rec.CharacterInstanceID.Valid &&
			!rec.MonsterInstanceID.Valid {
			return fmt.Errorf("failed validation, all of LocationInstanceID, CharacterInstanceID and MonsterInstanceID are empty")
		}
	}

	return nil
}

// ValidateDeleteObjectInstanceRec - validates it is okay to delete a game record
func (m *Model) ValidateDeleteObjectInstanceRec(recID string) error {

	return nil
}
