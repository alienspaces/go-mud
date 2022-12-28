package model

import (
	"fmt"
	"time"

	"gitlab.com/alienspaces/go-mud/backend/core/nulltime"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

type IncrementDungeonInstanceTurnResult struct {
	Record           *record.Turn
	WaitMilliseconds int64
	Incremented      bool
}

// IncrementDungeonInstanceTurnRec -
func (m *Model) IncrementDungeonInstanceTurn(recordID string) (*IncrementDungeonInstanceTurnResult, error) {
	l := m.Logger("IncrementDungeonInstanceTurn")
	l.Info("Attempting to increment dungeon instance ID >%s< turn", recordID)

	recs, err := m.GetTurnRecs(
		map[string]interface{}{
			"dungeon_instance_id": recordID,
		}, nil, true,
	)
	if err != nil {
		l.Warn("failed getting turn records >%v<", err)
		return nil, err
	}

	var rec *record.Turn
	if len(recs) > 1 {
		err := fmt.Errorf("unexpected number of turn records returned >%d< for dunge instance ID >%s<", len(recs), recordID)
		l.Warn(err.Error())
		return nil, err
	}

	if len(recs) == 0 {
		rec = &record.Turn{
			DungeonInstanceID: recordID,
		}
		err := m.CreateTurnRec(rec)
		if err != nil {
			l.Warn("failed creating turn record >%v<", err)
			return nil, err
		}
	} else {
		rec = recs[0]
	}

	// Check time since last turn increment
	sinceLastIncremented := time.Since(nulltime.ToTime(rec.IncrementedAt))
	l.Info("Last incremented duration >%d<", sinceLastIncremented.Milliseconds())
	l.Info("Turn duration             >%d<", m.turnDuration.Milliseconds())

	if sinceLastIncremented < m.turnDuration {
		l.Info("Too early to increment, since last incremented %d < duration %d", sinceLastIncremented.Milliseconds(), m.turnDuration.Milliseconds())
		return &IncrementDungeonInstanceTurnResult{
			Record:           nil,
			WaitMilliseconds: m.turnDuration.Milliseconds() - sinceLastIncremented.Milliseconds(),
			Incremented:      false,
		}, nil
	}

	l.Info("Can increment, since last incremented %d > duration %d", sinceLastIncremented.Milliseconds(), m.turnDuration.Milliseconds())

	rec.TurnCount++
	rec.IncrementedAt = nulltime.FromTime(time.Now().UTC())

	err = m.UpdateTurnRec(rec)
	if err != nil {
		l.Warn("failed updating turn record >%v<", err)
		return nil, err
	}

	return &IncrementDungeonInstanceTurnResult{
		Record:      rec,
		Incremented: true,
	}, nil
}
