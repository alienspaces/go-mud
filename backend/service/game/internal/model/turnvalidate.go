package model

import (
	"fmt"

	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

// validateTurnRec - validates creating and updating a turn record
func (m *Model) validateTurnRec(rec *record.Turn) error {
	l := m.Logger("validateTurnRec")

	// New turn
	if rec.ID == "" {
		// Can only have a single turn record per dungeon instance
		recs, err := m.GetTurnRecs(
			map[string]interface{}{
				"dungeon_instance_id": rec.DungeonInstanceID,
			}, nil, false,
		)
		if err != nil {
			l.Warn("failed to get many turn records >%v<", err)
		}
		if len(recs) != 0 {
			err := fmt.Errorf("turn record for dungeon instance ID >%s< already exists, cannot create turn record", rec.DungeonInstanceID)
			l.Warn(err.Error())
			return err
		}
	} else {
		// TODO: validate turn duration
		l.Info("turn duration >%d", m.turnDuration)
	}

	if rec.DungeonInstanceID == "" {
		return fmt.Errorf("failed validation, DungeonInstanceID is empty")
	}

	return nil
}

// validateDeleteTurnRec - validates it is okay to delete a turn record
func (m *Model) validateDeleteTurnRec(recID string) error {

	return nil
}
