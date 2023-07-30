package model

import (
	"database/sql"
	"fmt"

	"gitlab.com/alienspaces/go-mud/backend/core/null"
	coresql "gitlab.com/alienspaces/go-mud/backend/core/sql"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/calculator"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/mapper"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

type PerformActionArgs struct {
	ActionRec                 *record.Action
	CharacterInstanceViewRec  *record.CharacterInstanceView
	MonsterInstanceViewRec    *record.MonsterInstanceView
	LocationInstanceRecordSet *record.LocationInstanceViewRecordSet
}

type characterActionFunc func(args *PerformActionArgs) (*record.Action, error)

func checkPerformActionArgs(args *PerformActionArgs) error {
	if args.ActionRec == nil {
		msg := "args ActionRec missing, cannot perform action"
		return fmt.Errorf(msg)
	}
	if args.ActionRec.TurnNumber == 0 {
		err := fmt.Errorf("args ActionRec turn is zero, cannot perform action")
		return err
	}
	if args.CharacterInstanceViewRec == nil && args.MonsterInstanceViewRec == nil {
		err := fmt.Errorf("args CharacterInstanceViewRec and MonsterInstanceViewRec are missing, cannot perform action")
		return err
	}
	if args.LocationInstanceRecordSet == nil {
		err := fmt.Errorf("args LocationInstanceRecordSet missing, cannot perform action")
		return err
	}
	return nil
}

func (m *Model) performAction(args *PerformActionArgs) (*record.Action, error) {
	l := m.loggerWithFunctionContext("performCharacterAction")

	if err := checkPerformActionArgs(args); err != nil {
		l.Warn("failed checking performer args >%v<", err)
		return nil, err
	}

	actionRec := args.ActionRec

	actionFuncs := map[string]characterActionFunc{
		record.ActionCommandMove:   m.performActionMove,
		record.ActionCommandLook:   m.performActionLook,
		record.ActionCommandUse:    m.performActionUse,
		record.ActionCommandEquip:  m.performActionEquip,
		record.ActionCommandStash:  m.performActionStash,
		record.ActionCommandDrop:   m.performActionDrop,
		record.ActionCommandAttack: m.performActionAttack,
	}

	actionFunc, ok := actionFuncs[actionRec.ResolvedCommand]
	if !ok {
		msg := fmt.Sprintf("Action >%s< not supported", actionRec.ResolvedCommand)
		l.Warn(msg)
		return nil, fmt.Errorf(msg)
	}

	l.Info("Performing action resolved command >%s<", actionRec.ResolvedCommand)

	var err error
	actionRec, err = actionFunc(args)
	if err != nil {
		l.Warn("failed performing action >%v<", err)
		return nil, err
	}

	l.Debug("Have updated dungeon action record >%v<", actionRec)

	return actionRec, nil
}

func (m *Model) performActionMove(args *PerformActionArgs) (*record.Action, error) {
	l := m.loggerWithFunctionContext("performActionMove")

	if err := checkPerformActionArgs(args); err != nil {
		l.Warn("failed checking performer args >%v<", err)
		return nil, err
	}

	actionRec := args.ActionRec
	characterInstanceViewRec := args.CharacterInstanceViewRec
	monsterInstanceViewRec := args.MonsterInstanceViewRec

	if actionRec.CharacterInstanceID.Valid {
		// Character move direction
		characterInstanceRec, err := mapper.CharacterInstanceViewToCharacterInstance(l, characterInstanceViewRec)
		if err != nil {
			l.Warn("failed mapping character instance view to character instance >%v<", err)
			return nil, err
		}

		characterInstanceRec.LocationInstanceID = actionRec.ResolvedTargetLocationInstanceID.String

		err = m.UpdateCharacterInstanceRec(characterInstanceRec)
		if err != nil {
			l.Warn("failed updated dungeon character instance record >%v<", err)
			return nil, err
		}

		// Update dungeon action record
		actionRec.LocationInstanceID = actionRec.ResolvedTargetLocationInstanceID.String
	} else if actionRec.MonsterInstanceID.Valid {
		// Monster move direction
		monsterInstanceRec, err := mapper.MonsterInstanceViewToMonsterInstance(l, monsterInstanceViewRec)
		if err != nil {
			l.Warn("failed mapping monster instance view to monster instance >%v<", err)
			return nil, err
		}

		monsterInstanceRec.LocationInstanceID = actionRec.ResolvedTargetLocationInstanceID.String

		err = m.UpdateMonsterInstanceRec(monsterInstanceRec)
		if err != nil {
			l.Warn("failed updated dungeon monster instance record >%v<", err)
			return nil, err
		}

		// Update dungeon action record
		actionRec.LocationInstanceID = actionRec.ResolvedTargetLocationInstanceID.String
	}

	return actionRec, nil
}

func (m *Model) performActionLook(args *PerformActionArgs) (*record.Action, error) {
	l := m.loggerWithFunctionContext("performActionLook")

	if err := checkPerformActionArgs(args); err != nil {
		l.Warn("failed checking performer args >%v<", err)
		return nil, err
	}

	actionRec := args.ActionRec

	// TODO: (game) Register the number of times the character has looked at any other
	// location, object, monster of character. Looking at anything multiple times
	// should result in additional information within X turns.

	if actionRec.ResolvedTargetLocationInstanceID.Valid {
		l.Debug("Looking at location ID >%s<", actionRec.ResolvedTargetLocationInstanceID.String)

	} else if actionRec.ResolvedTargetObjectInstanceID.Valid {
		l.Debug("Looking at object ID >%s<", actionRec.ResolvedTargetObjectInstanceID.String)

	} else if actionRec.ResolvedTargetMonsterInstanceID.Valid {
		l.Debug("Looking at monster ID >%s<", actionRec.ResolvedTargetMonsterInstanceID.String)

	} else if actionRec.ResolvedTargetCharacterInstanceID.Valid {
		l.Debug("Looking at character ID >%s<", actionRec.ResolvedTargetCharacterInstanceID.String)
	}

	return actionRec, nil
}

func (m *Model) performActionUse(args *PerformActionArgs) (*record.Action, error) {
	l := m.loggerWithFunctionContext("performActionUse")

	if err := checkPerformActionArgs(args); err != nil {
		l.Warn("failed checking performer args >%v<", err)
		return nil, err
	}

	actionRec := args.ActionRec

	if actionRec.ResolvedTargetMonsterInstanceID.Valid {
		l.Debug("Using object ID >%s< on monster ID >%s<", null.NullStringToString(actionRec.ResolvedTargetObjectInstanceID), null.NullStringToString(actionRec.ResolvedTargetMonsterInstanceID))
	} else if actionRec.ResolvedTargetCharacterInstanceID.Valid {
		l.Debug("Using object ID >%s< on character ID >%s<", null.NullStringToString(actionRec.ResolvedTargetObjectInstanceID), null.NullStringToString(actionRec.ResolvedTargetCharacterInstanceID))
	} else {
		l.Debug("Using object ID >%s<", null.NullStringToString(actionRec.ResolvedTargetObjectInstanceID))
	}

	return actionRec, nil
}

func (m *Model) performActionStash(args *PerformActionArgs) (*record.Action, error) {
	l := m.loggerWithFunctionContext("performActionStash")

	if err := checkPerformActionArgs(args); err != nil {
		l.Warn("failed checking performer args >%v<", err)
		return nil, err
	}

	actionRec := args.ActionRec

	if actionRec.CharacterInstanceID.Valid {
		// Character stash object
		objectInstanceID := actionRec.ResolvedStashedObjectInstanceID.String
		if objectInstanceID == "" {
			msg := "resolved stashed dungeon object instance ID is empty, cannot stash object"
			l.Warn(msg)
			return nil, fmt.Errorf(msg)
		}

		objectInstanceRec, err := m.GetObjectInstanceRec(objectInstanceID, coresql.ForUpdate)
		if err != nil {
			l.Warn("failed getting dungeon object instance record >%v<", err)
			return nil, err
		}

		objectInstanceRec.LocationInstanceID = sql.NullString{}
		objectInstanceRec.CharacterInstanceID = actionRec.CharacterInstanceID
		objectInstanceRec.IsStashed = true
		objectInstanceRec.IsEquipped = false

		err = m.UpdateObjectInstanceRec(objectInstanceRec)
		if err != nil {
			l.Warn("failed updating dungeon object instance record >%v<", err)
			return nil, err
		}

	} else if actionRec.MonsterInstanceID.Valid {
		// Monster stash object
		objectInstanceID := actionRec.ResolvedStashedObjectInstanceID.String
		if objectInstanceID == "" {
			msg := "resolved stashed dungeon object instance ID is empty, cannot stash object"
			l.Warn(msg)
			return nil, fmt.Errorf(msg)
		}

		objectInstanceRec, err := m.GetObjectInstanceRec(objectInstanceID, coresql.ForUpdate)
		if err != nil {
			l.Warn("failed getting dungeon object instance record >%v<", err)
			return nil, err
		}

		objectInstanceRec.LocationInstanceID = sql.NullString{}
		objectInstanceRec.MonsterInstanceID = actionRec.MonsterInstanceID
		objectInstanceRec.IsStashed = true
		objectInstanceRec.IsEquipped = false

		err = m.UpdateObjectInstanceRec(objectInstanceRec)
		if err != nil {
			l.Warn("failed updating dungeon object instance record >%v<", err)
			return nil, err
		}
	}

	return actionRec, nil
}

func (m *Model) performActionEquip(args *PerformActionArgs) (*record.Action, error) {
	l := m.loggerWithFunctionContext("performActionEquip")

	if err := checkPerformActionArgs(args); err != nil {
		l.Warn("failed checking performer args >%v<", err)
		return nil, err
	}

	actionRec := args.ActionRec

	if actionRec.CharacterInstanceID.Valid {
		// Character equip object
		objectInstanceID := actionRec.ResolvedEquippedObjectInstanceID.String
		if objectInstanceID == "" {
			msg := "resolved equipped dungeon object instance ID is empty, cannot equipe object"
			l.Warn(msg)
			return nil, fmt.Errorf(msg)
		}

		objectInstanceRec, err := m.GetObjectInstanceRec(objectInstanceID, coresql.ForUpdate)
		if err != nil {
			l.Warn("failed getting dungeon object instance record >%v<", err)
			return nil, err
		}

		objectInstanceRec.LocationInstanceID = sql.NullString{}
		objectInstanceRec.CharacterInstanceID = actionRec.CharacterInstanceID
		objectInstanceRec.IsEquipped = true
		objectInstanceRec.IsStashed = false

		err = m.UpdateObjectInstanceRec(objectInstanceRec)
		if err != nil {
			l.Warn("failed updating dungeon object instance record >%v<", err)
			return nil, err
		}

	} else if actionRec.MonsterInstanceID.Valid {
		// Monster equip object
		objectInstanceID := actionRec.ResolvedEquippedObjectInstanceID.String
		if objectInstanceID == "" {
			msg := "resolved equipped dungeon object instance ID is empty, cannot equipe object"
			l.Warn(msg)
			return nil, fmt.Errorf(msg)
		}

		objectInstanceRec, err := m.GetObjectInstanceRec(objectInstanceID, coresql.ForUpdate)
		if err != nil {
			l.Warn("failed getting dungeon object instance record >%v<", err)
			return nil, err
		}

		objectInstanceRec.LocationInstanceID = sql.NullString{}
		objectInstanceRec.MonsterInstanceID = actionRec.MonsterInstanceID
		objectInstanceRec.IsEquipped = true
		objectInstanceRec.IsStashed = false

		err = m.UpdateObjectInstanceRec(objectInstanceRec)
		if err != nil {
			l.Warn("failed updating dungeon object instance record >%v<", err)
			return nil, err
		}
	}

	return actionRec, nil
}

func (m *Model) performActionDrop(args *PerformActionArgs) (*record.Action, error) {
	l := m.loggerWithFunctionContext("performActionDrop")

	if err := checkPerformActionArgs(args); err != nil {
		l.Warn("failed checking performer args >%v<", err)
		return nil, err
	}

	actionRec := args.ActionRec

	if actionRec.CharacterInstanceID.Valid {
		// Character drop object
		objectInstanceID := actionRec.ResolvedDroppedObjectInstanceID.String
		if objectInstanceID == "" {
			msg := "resolved dropped dungeon object instance ID is empty, cannot drop object"
			l.Warn(msg)
			return nil, fmt.Errorf(msg)
		}

		objectInstanceRec, err := m.GetObjectInstanceRec(objectInstanceID, coresql.ForUpdate)
		if err != nil {
			l.Warn("failed getting dungeon object instance record >%v<", err)
			return nil, err
		}

		objectInstanceRec.LocationInstanceID = sql.NullString{
			String: actionRec.LocationInstanceID,
			Valid:  true,
		}
		objectInstanceRec.CharacterInstanceID = sql.NullString{}
		objectInstanceRec.IsStashed = false
		objectInstanceRec.IsEquipped = false

		l.Debug("Updating dropped object instance >%#v<", objectInstanceRec)

		err = m.UpdateObjectInstanceRec(objectInstanceRec)
		if err != nil {
			l.Warn("failed updating dungeon object instance record >%v<", err)
			return nil, err
		}

	} else if actionRec.MonsterInstanceID.Valid {
		// Monster drop object
		objectInstanceID := actionRec.ResolvedDroppedObjectInstanceID.String
		if objectInstanceID == "" {
			msg := "resolved dropped dungeon object instance ID is empty, cannot drop object"
			l.Warn(msg)
			return nil, fmt.Errorf(msg)
		}

		objectInstanceRec, err := m.GetObjectInstanceRec(objectInstanceID, coresql.ForUpdate)
		if err != nil {
			l.Warn("failed getting dungeon object instance record >%v<", err)
			return nil, err
		}

		objectInstanceRec.LocationInstanceID = sql.NullString{
			String: actionRec.LocationInstanceID,
			Valid:  true,
		}
		objectInstanceRec.MonsterInstanceID = sql.NullString{}
		objectInstanceRec.IsStashed = false
		objectInstanceRec.IsEquipped = false

		l.Debug("Updating dropped object instance >%#v<", objectInstanceRec)

		err = m.UpdateObjectInstanceRec(objectInstanceRec)
		if err != nil {
			l.Warn("failed updating dungeon object instance record >%v<", err)
			return nil, err
		}
	}

	return actionRec, nil
}

// TODO: 10-implement-effects: Calculate to-hit, weapon damage, effects etc
func (m *Model) performActionAttack(args *PerformActionArgs) (*record.Action, error) {
	l := m.loggerWithFunctionContext("performActionAttack")

	if err := checkPerformActionArgs(args); err != nil {
		l.Warn("failed checking performer args >%v<", err)
		return nil, err
	}

	actionRec := args.ActionRec
	characterInstanceRec := args.CharacterInstanceViewRec
	monsterInstanceRec := args.MonsterInstanceViewRec
	locationInstanceRecordSet := args.LocationInstanceRecordSet

	if null.NullStringIsValid(actionRec.CharacterInstanceID) && characterInstanceRec != nil {

		// TODO: 10-implement-effects:
		// Get resolved equipped weapon or the primary equipped weapon being used to attack.
		if null.NullStringIsValid(actionRec.ResolvedEquippedObjectInstanceID) {
			l.Info("Attacking with weapon")
		}

		// TODO: 10-implement-effects:
		// Rename and refactor this function so it returns an array of effects. Pass the
		// character instance record and the equipped item that is being used to attack.
		dmg, err := calculator.CalculateCharacterDamage(characterInstanceRec, locationInstanceRecordSet)
		if err != nil {
			l.Warn("failed calculating character damage >%v<", err)
			return nil, err
		}

		if null.NullStringIsValid(actionRec.ResolvedTargetCharacterInstanceID) {

			l.Info("Character attacking character")

			tciRec, err := m.GetCharacterInstanceRec(null.NullStringToString(actionRec.ResolvedTargetCharacterInstanceID), coresql.ForUpdate)
			if err != nil {
				l.Warn("failed getting character instance record >%v<", err)
				return nil, err
			}

			// TODO: 10-implement-effects:
			// Apply resulting effects to the target character instance instead of directly
			// subtracting damage from the characters health.
			tciRec.Health -= dmg

			err = m.UpdateCharacterInstanceRec(tciRec)
			if err != nil {
				l.Warn("failed updating character instance record >%v<", err)
				return nil, err
			}

		} else if null.NullStringIsValid(actionRec.ResolvedTargetMonsterInstanceID) {

			l.Info("Character attacking monster")

			tmiRec, err := m.GetMonsterInstanceRec(null.NullStringToString(actionRec.ResolvedTargetMonsterInstanceID), coresql.ForUpdate)
			if err != nil {
				l.Warn("failed getting monster instance record >%v<", err)
				return nil, err
			}

			// TODO: 10-implement-effects:
			// Apply resulting effects to the target monster instance instead of directly
			// subtracting damage from the monsters health.
			tmiRec.Health -= dmg

			err = m.UpdateMonsterInstanceRec(tmiRec)
			if err != nil {
				l.Warn("failed updating monster instance record >%v<", err)
				return nil, err
			}
		}

	} else if null.NullStringIsValid(actionRec.MonsterInstanceID) && monsterInstanceRec != nil {

		// TODO: 10-implement-effects:
		// Get resolved equipped weapon or the primary equipped weapon being used to attack.
		if null.NullStringIsValid(actionRec.ResolvedEquippedObjectInstanceID) {
			l.Info("Attacking with weapon")
		}

		// TODO: 10-implement-effects:
		// Rename and refactor this function so it returns an array of effects. Pass the
		// character instance record and the equipped item that is being used to attack.
		dmg, err := calculator.CalculateMonsterDamage(monsterInstanceRec, locationInstanceRecordSet)
		if err != nil {
			l.Warn("failed calculating monster damage >%v<", err)
			return nil, err
		}

		if null.NullStringIsValid(actionRec.ResolvedTargetCharacterInstanceID) {

			l.Info("Monster attacking character")

			tciRec, err := m.GetCharacterInstanceRec(null.NullStringToString(actionRec.ResolvedTargetCharacterInstanceID), coresql.ForUpdate)
			if err != nil {
				l.Warn("failed getting character instance record >%v<", err)
				return nil, err
			}

			if tciRec == nil {
				err := fmt.Errorf("failed getting character instance record ID >%s<", null.NullStringToString(actionRec.ResolvedTargetCharacterInstanceID))
				l.Warn(err.Error())
				return nil, err
			}

			// TODO: 10-implement-effects:
			// Apply resulting effects to the target character instance instead of directly
			// subtracting damage from the characters health.
			tciRec.Health -= dmg

			err = m.UpdateCharacterInstanceRec(tciRec)
			if err != nil {
				l.Warn("failed updating character instance record >%v<", err)
				return nil, err
			}

		} else if null.NullStringIsValid(actionRec.ResolvedTargetMonsterInstanceID) {

			l.Info("Monster attacking monster")

			tmiRec, err := m.GetMonsterInstanceRec(null.NullStringToString(actionRec.ResolvedTargetMonsterInstanceID), coresql.ForUpdate)
			if err != nil {
				l.Warn("failed getting monster instance record >%v<", err)
				return nil, err
			}

			// TODO: 10-implement-effects:
			// Apply resulting effects to the target monster instance instead of directly
			// subtracting damage from the monsters health.
			tmiRec.Health -= dmg

			err = m.UpdateMonsterInstanceRec(tmiRec)
			if err != nil {
				l.Warn("failed updating monster instance record >%v<", err)
				return nil, err
			}
		}
	}

	return actionRec, nil
}
