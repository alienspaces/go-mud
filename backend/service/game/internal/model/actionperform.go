package model

import (
	"database/sql"
	"fmt"

	"gitlab.com/alienspaces/go-mud/backend/core/nullstring"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/mapper"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

// TODO: (game) Determine whether we need to pass the character/monster instance view
// record around everywhere or whether a more specific PerformActionArgs definition
// (like ResolveActionArgs) would clean this up.

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
	l := m.Logger("performCharacterAction")

	if err := checkPerformActionArgs(args); err != nil {
		l.Warn("failed checking performer args >%v<", err)
		return nil, err
	}

	actionRec := args.ActionRec

	actionFuncs := map[string]characterActionFunc{
		"move":   m.performActionMove,
		"look":   m.performActionLook,
		"use":    m.performActionUse,
		"stash":  m.performActionStash,
		"equip":  m.performActionEquip,
		"drop":   m.performActionDrop,
		"attack": m.performActionAttack,
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
	l := m.Logger("performActionMove")

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
	l := m.Logger("performActionLook")

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
	l := m.Logger("performActionUse")

	if err := checkPerformActionArgs(args); err != nil {
		l.Warn("failed checking performer args >%v<", err)
		return nil, err
	}

	actionRec := args.ActionRec

	if actionRec.ResolvedTargetMonsterInstanceID.Valid {
		l.Debug("Using object ID >%s< on monster ID >%s<", nullstring.ToString(actionRec.ResolvedTargetObjectInstanceID), nullstring.ToString(actionRec.ResolvedTargetMonsterInstanceID))
	} else if actionRec.ResolvedTargetCharacterInstanceID.Valid {
		l.Debug("Using object ID >%s< on character ID >%s<", nullstring.ToString(actionRec.ResolvedTargetObjectInstanceID), nullstring.ToString(actionRec.ResolvedTargetCharacterInstanceID))
	} else {
		l.Debug("Using object ID >%s<", nullstring.ToString(actionRec.ResolvedTargetObjectInstanceID))
	}

	return actionRec, nil
}

func (m *Model) performActionStash(args *PerformActionArgs) (*record.Action, error) {
	l := m.Logger("performActionStash")

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

		objectInstanceRec, err := m.GetObjectInstanceRec(objectInstanceID, true)
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

		objectInstanceRec, err := m.GetObjectInstanceRec(objectInstanceID, true)
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
	l := m.Logger("performActionEquip")

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

		objectInstanceRec, err := m.GetObjectInstanceRec(objectInstanceID, true)
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

		objectInstanceRec, err := m.GetObjectInstanceRec(objectInstanceID, true)
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
	l := m.Logger("performActionDrop")

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

		objectInstanceRec, err := m.GetObjectInstanceRec(objectInstanceID, true)
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

		objectInstanceRec, err := m.GetObjectInstanceRec(objectInstanceID, true)
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
	l := m.Logger("performActionAttack")

	if err := checkPerformActionArgs(args); err != nil {
		l.Warn("failed checking performer args >%v<", err)
		return nil, err
	}

	actionRec := args.ActionRec

	if actionRec.CharacterInstanceID.Valid {

		// TODO: 10-implement-effects: Get equipped weapon for character, establish attack bonuses, damage rating etc
		if nullstring.IsValid(actionRec.ResolvedEquippedObjectInstanceID) {
			l.Info("Attacking with weapon")
		}

		if nullstring.IsValid(actionRec.ResolvedTargetCharacterInstanceID) {
			l.Info("Character attacking character")
			tciRec, err := m.GetCharacterInstanceRec(nullstring.ToString(actionRec.ResolvedTargetCharacterInstanceID), true)
			if err != nil {
				l.Warn("failed getting character instance record >%s<", err)
				return nil, err
			}

			tciRec.Health -= 1

			err = m.UpdateCharacterInstanceRec(tciRec)
			if err != nil {
				l.Warn("failed updating character instance record >%s<", err)
				return nil, err
			}

			// TODO: 12-implement-death: Remove character instance when dead

		} else if nullstring.IsValid(actionRec.ResolvedTargetMonsterInstanceID) {
			l.Info("Character attacking monster")
			tmiRec, err := m.GetMonsterInstanceRec(nullstring.ToString(actionRec.ResolvedTargetMonsterInstanceID), true)
			if err != nil {
				l.Warn("failed getting monster instance record >%s<", err)
				return nil, err
			}

			tmiRec.Health -= 1

			err = m.UpdateMonsterInstanceRec(tmiRec)
			if err != nil {
				l.Warn("failed updating monster instance record >%s<", err)
				return nil, err
			}

			// TODO: 12-implement-death: Remove monster instance when dead

		}
	} else if actionRec.MonsterInstanceID.Valid {

		// TODO: Get equipped weapon for monster, establish attack bonuses, damage rating etc
		if nullstring.IsValid(actionRec.ResolvedEquippedObjectInstanceID) {
			l.Info("Attacking with weapon")
		}

		if nullstring.IsValid(actionRec.ResolvedTargetCharacterInstanceID) {
			l.Info("Monster attacking character")
			tciRec, err := m.GetCharacterInstanceRec(nullstring.ToString(actionRec.ResolvedTargetCharacterInstanceID), true)
			if err != nil {
				l.Warn("failed getting character instance record >%s<", err)
				return nil, err
			}

			tciRec.Health -= 1

			err = m.UpdateCharacterInstanceRec(tciRec)
			if err != nil {
				l.Warn("failed updating character instance record >%s<", err)
				return nil, err
			}

			// TODO: 12-implement-death: Remove character instance when dead

		} else if nullstring.IsValid(actionRec.ResolvedTargetMonsterInstanceID) {
			l.Info("Monster attacking monster")
			tmiRec, err := m.GetMonsterInstanceRec(nullstring.ToString(actionRec.ResolvedTargetMonsterInstanceID), true)
			if err != nil {
				l.Warn("failed getting monster instance record >%s<", err)
				return nil, err
			}

			tmiRec.Health -= 1

			err = m.UpdateMonsterInstanceRec(tmiRec)
			if err != nil {
				l.Warn("failed updating monster instance record >%s<", err)
				return nil, err
			}

			// TODO: 12-implement-death: Remove monster instance when dead
		}
	}

	return actionRec, nil
}
