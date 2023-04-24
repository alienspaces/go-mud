package model

import (
	"fmt"

	"gitlab.com/alienspaces/go-mud/backend/core/nullstring"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

// validateMonsterInstanceMemoryRec - validates creating and updating a monster instance record
func (m *Model) validateMonsterInstanceMemoryRec(rec *record.MonsterInstanceMemory) error {

	// New monster
	if rec.ID == "" {
		if rec.MemoryCommand == "" {
			return fmt.Errorf("failed validation, MemoryCommand is empty")
		}
		if rec.MemoryType == "" {
			return fmt.Errorf("failed validation, MemoryType is empty")
		}
		if rec.TurnNumber == 0 {
			return fmt.Errorf("failed validation, TurnNumber is empty")
		}
		if !nullstring.IsValid(rec.MemoryLocationInstanceID) &&
			!nullstring.IsValid(rec.MemoryCharacterInstanceID) &&
			!nullstring.IsValid(rec.MemoryMonsterInstanceID) &&
			!nullstring.IsValid(rec.MemoryObjectInstanceID) {
			return fmt.Errorf("failed validation, one of MemoryLocationInstanceID, MemoryCharacterInstanceID, MemoryMonsterInstanceID, MemoryObjectInstanceID must be set")
		}
	}

	return nil
}

// validateDeleteMonsterInstanceMemoryRec - validates it is okay to delete a monster instance record
func (m *Model) validateDeleteMonsterInstanceMemoryRec(recID string) error {

	return nil
}
