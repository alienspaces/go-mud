package model

import (
	"fmt"

	coresql "gitlab.com/alienspaces/go-mud/backend/core/sql"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

type ResolveActionTurnArgs struct {
	ActionRec         *record.Action
	EntityType        EntityType
	EntityInstanceID  string
	DungeonInstanceID string
}

func (m *Model) resolveActionTurn(args *ResolveActionTurnArgs) (*record.Action, error) {
	l := m.loggerWithContext("resolveActionTurn")

	if args == nil {
		err := fmt.Errorf("required args is nill, cannot resolve action turn")
		l.Warn(err.Error())
		return nil, err
	}

	actionRec := args.ActionRec
	if actionRec == nil {
		err := fmt.Errorf("required args ActionRec is nill, cannot resolve action turn")
		l.Warn(err.Error())
		return nil, err
	}

	l.Info("Resolving action turn >%#v<", args)

	// Get the dungeon entity instance turn record
	q := m.DungeonEntityInstanceTurnQuery()

	recs, err := q.GetMany(
		&coresql.Options{
			Params: []coresql.Param{
				{
					Col: "dungeon_instance_id",
					Val: args.DungeonInstanceID,
				},
				{
					Col: "entity_type",
					Val: args.EntityType,
				},
				{
					Col: "entity_instance_id",
					Val: args.EntityInstanceID,
				},
			},
		},
	)
	if err != nil {
		l.Warn("failed querying dungeon entity instance turns >%v<", err)
		return nil, err
	}

	// When no records are returned it would mean the character or monster
	// has yet to perform their first action.
	if len(recs) == 0 {
		actionRec.TurnNumber = 1
		return actionRec, err
	}

	if len(recs) > 1 {
		err := fmt.Errorf("unexepected number of dungeon entity instance turn records >%d< return for dungeon instance ID >%s< entity type >%s< entity instance ID >%s<", len(recs), args.DungeonInstanceID, args.EntityType, args.EntityInstanceID)
		l.Warn(err.Error())
		return nil, err
	}

	rec := recs[0]
	if rec.EntityInstanceTurnNumber >= rec.DungeonInstanceTurnNumber {
		msg := fmt.Sprintf("dungeon instance turn >%d< is less than or equal to entity instance turn >%d<", rec.DungeonInstanceTurnNumber, rec.EntityInstanceTurnNumber)
		l.Warn(msg)
		return nil, NewActionTooEarlyError(rec.DungeonInstanceTurnNumber, rec.EntityInstanceTurnNumber)
	}

	// A character or monster can choose to not execute an action for every turn so
	// whenever the dungeon instance turn number is greater than the entity instance
	// turn number we will just assign the current dungeon instance turn number.
	if rec.DungeonInstanceTurnNumber > rec.EntityInstanceTurnNumber {
		actionRec.TurnNumber = rec.DungeonInstanceTurnNumber
	}

	return actionRec, nil
}
