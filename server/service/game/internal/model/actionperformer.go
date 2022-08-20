package model

import (
	"database/sql"
	"fmt"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/mapper"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

type characterActionFunc func(characterInstanceViewRec *record.CharacterInstanceView, actionRec *record.Action, locationInstanceRecordSet *record.LocationInstanceViewRecordSet) (*record.Action, error)

func (m *Model) performCharacterAction(
	characterInstanceViewRec *record.CharacterInstanceView,
	actionRec *record.Action,
	locationInstanceRecordSet *record.LocationInstanceViewRecordSet,
) (*record.Action, error) {

	l := m.Logger("performCharacterAction")

	actionFuncs := map[string]characterActionFunc{
		"move":  m.performActionMove,
		"look":  m.performActionLook,
		"stash": m.performActionStash,
		"equip": m.performActionEquip,
		"drop":  m.performActionDrop,
	}

	actionFunc, ok := actionFuncs[actionRec.ResolvedCommand]
	if !ok {
		msg := fmt.Sprintf("Action >%s< not supported", actionRec.ResolvedCommand)
		l.Warn(msg)
		return nil, fmt.Errorf(msg)
	}

	var err error
	actionRec, err = actionFunc(characterInstanceViewRec, actionRec, locationInstanceRecordSet)
	if err != nil {
		l.Warn("failed performing action >%s< >%v<", actionRec.ResolvedCommand, err)
		return nil, err
	}

	l.Debug("Have updated dungeon action record >%v<", actionRec)

	return actionRec, nil
}

func (m *Model) performActionMove(
	characterInstanceViewRec *record.CharacterInstanceView,
	actionRec *record.Action,
	locationInstanceRecordSet *record.LocationInstanceViewRecordSet,
) (*record.Action, error) {

	l := m.Logger("performActionMove")

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
		return nil, fmt.Errorf("moving monster instances is currently not supported")
	}

	return actionRec, nil
}

func (m *Model) performActionLook(
	characterInstanceViewRec *record.CharacterInstanceView,
	actionRec *record.Action,
	locationInstanceRecordSet *record.LocationInstanceViewRecordSet,
) (*record.Action, error) {

	l := m.Logger("performActionLook")

	// TODO: Register the number of times the character has looked at any other
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

func (m *Model) performActionStash(
	characterInstanceViewRec *record.CharacterInstanceView,
	actionRec *record.Action,
	locationInstanceRecordSet *record.LocationInstanceViewRecordSet,
) (*record.Action, error) {

	l := m.Logger("performActionStash")

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
		return nil, fmt.Errorf("monster instances stashing objects is currently not supported")
	}

	return actionRec, nil
}

func (m *Model) performActionEquip(
	characterInstanceViewRec *record.CharacterInstanceView,
	actionRec *record.Action,
	locationInstanceRecordSet *record.LocationInstanceViewRecordSet,
) (*record.Action, error) {

	l := m.Logger("performActionEquip")

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
		return nil, fmt.Errorf("monster instances equipping objects is currently not supported")
	}

	return actionRec, nil
}

func (m *Model) performActionDrop(
	characterInstanceViewRec *record.CharacterInstanceView,
	actionRec *record.Action,
	locationInstanceRecordSet *record.LocationInstanceViewRecordSet,
) (*record.Action, error) {

	l := m.Logger("performActionDrop")

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
		return nil, fmt.Errorf("monster instances dropping objects is currently not supported")
	}

	return actionRec, nil
}
