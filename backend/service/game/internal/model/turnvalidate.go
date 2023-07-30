package model

import (
	"fmt"

	coresql "gitlab.com/alienspaces/go-mud/backend/core/sql"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

// validateTurnRec - validates creating and updating a turn record
func (m *Model) validateTurnRec(rec *record.Turn) error {
	l := m.loggerWithFunctionContext("validateTurnRec")

	// New turn
	if rec.ID == "" {
		// Can only have a single turn record per dungeon instance
		recs, err := m.GetTurnRecs(
			&coresql.Options{
				Params: []coresql.Param{
					{
						Col: "dungeon_instance_id",
						Val: rec.DungeonInstanceID,
					},
				},
			},
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
		l.Debug("turn duration >%d", m.turnDuration)
		currRec, err := m.GetTurnRec(rec.ID, nil)
		if err != nil {
			l.Warn("failed getting existing turn record >%v<", err)
			return err
		}
		if rec.TurnNumber != currRec.TurnNumber+1 {
			err := fmt.Errorf("updated turn number >%d< is not an increment of current turn number >%d<", rec.TurnNumber, currRec.TurnNumber)
			l.Warn(err.Error())
			return err
		}
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
