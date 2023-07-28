package model

import (
	"fmt"
	"time"

	"gitlab.com/alienspaces/go-mud/backend/core/null"
	coresql "gitlab.com/alienspaces/go-mud/backend/core/sql"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

type IncrementDungeonInstanceTurnArgs struct {
	DungeonInstanceID string
	TurnDuration      *time.Duration
}

type IncrementDungeonInstanceTurnResult struct {
	Record           *record.Turn
	WaitMilliseconds int64
	Incremented      bool
}

// IncrementDungeonInstanceTurnRec -
func (m *Model) IncrementDungeonInstanceTurn(args *IncrementDungeonInstanceTurnArgs) (*IncrementDungeonInstanceTurnResult, error) {
	l := m.loggerWithContext("IncrementDungeonInstanceTurn")

	l.Info("Increment dungeon instance ID >%s< turn", args.DungeonInstanceID)

	recs, err := m.GetTurnRecs(
		&coresql.Options{
			Params: []coresql.Param{
				{
					Col: "dungeon_instance_id",
					Val: args.DungeonInstanceID,
				},
			},
		},
	)
	if err != nil {
		l.Warn("failed getting turn records >%v<", err)
		return nil, err
	}

	var rec *record.Turn
	if len(recs) > 1 {
		err := fmt.Errorf("unexpected number of turn records returned >%d< for dunge instance ID >%s<", len(recs), args.DungeonInstanceID)
		l.Warn(err.Error())
		return nil, err
	}

	if len(recs) == 0 {
		rec = &record.Turn{
			DungeonInstanceID: args.DungeonInstanceID,
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
	turnDuration := m.turnDuration
	if args.TurnDuration != nil {
		turnDuration = *args.TurnDuration
	}

	sinceLastIncremented := time.Since(null.NullTimeToTime(rec.IncrementedAt))
	l.Debug("Last incremented duration >%d<", sinceLastIncremented.Milliseconds())
	l.Debug("Turn duration             >%d<", turnDuration.Milliseconds())

	if sinceLastIncremented < turnDuration {
		l.Debug("Too early to increment, since last incremented %d < duration %d", sinceLastIncremented.Milliseconds(), turnDuration.Milliseconds())
		return &IncrementDungeonInstanceTurnResult{
			Record:           nil,
			WaitMilliseconds: turnDuration.Milliseconds() - sinceLastIncremented.Milliseconds(),
			Incremented:      false,
		}, nil
	}

	l.Debug("Can increment, since last incremented %d > duration %d", sinceLastIncremented.Milliseconds(), turnDuration.Milliseconds())

	rec.TurnNumber++
	rec.IncrementedAt = null.NullTimeFromTime(time.Now().UTC())

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
