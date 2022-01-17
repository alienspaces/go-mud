package model

import (
	"database/sql"
	"fmt"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

func (m *Model) performDungeonCharacterAction(
	dungeonCharacterRec *record.DungeonCharacter,
	dungeonActionRec *record.DungeonAction,
	dungeonLocationRecordSet *DungeonLocationRecordSet,
) (*record.DungeonAction, error) {

	actionFuncs := map[string]func(dungeonCharacterRec *record.DungeonCharacter, dungeonActionRec *record.DungeonAction, dungeonLocationRecordSet *DungeonLocationRecordSet) (*record.DungeonAction, error){
		"move":  m.performDungeonActionMove,
		"look":  m.performDungeonActionLook,
		"stash": m.performDungeonActionStash,
	}

	actionFunc, ok := actionFuncs[dungeonActionRec.ResolvedCommand]
	if !ok {
		msg := fmt.Sprintf("Action >%s< not supported", dungeonActionRec.ResolvedCommand)
		m.Log.Warn(msg)
		return nil, fmt.Errorf(msg)
	}

	var err error
	dungeonActionRec, err = actionFunc(dungeonCharacterRec, dungeonActionRec, dungeonLocationRecordSet)
	if err != nil {
		m.Log.Warn("Failed performing action >%s< >%v<", dungeonActionRec.ResolvedCommand, err)
		return nil, err
	}

	m.Log.Info("Have updated dungeon action record >%v<", dungeonActionRec)

	return dungeonActionRec, nil
}

func (m *Model) performDungeonActionMove(
	dungeonCharacterRec *record.DungeonCharacter,
	dungeonActionRec *record.DungeonAction,
	dungeonLocationRecordSet *DungeonLocationRecordSet,
) (*record.DungeonAction, error) {

	if dungeonActionRec.DungeonCharacterID.Valid {
		// Character move direction
		dungeonCharacterRec.DungeonLocationID = dungeonActionRec.ResolvedTargetDungeonLocationID.String

		err := m.UpdateDungeonCharacterRec(dungeonCharacterRec)
		if err != nil {
			m.Log.Warn("Failed updated dungeon character record >%v<", err)
			return nil, err
		}

		// Update dungeon action record
		dungeonActionRec.DungeonLocationID = dungeonActionRec.ResolvedTargetDungeonLocationID.String
	} else if dungeonActionRec.DungeonMonsterID.Valid {
		// Monster move direction
		return nil, fmt.Errorf("moving monsters is currently not supported")
	}

	return dungeonActionRec, nil
}

func (m *Model) performDungeonActionLook(
	dungeonCharacterRec *record.DungeonCharacter,
	dungeonActionRec *record.DungeonAction,
	dungeonLocationRecordSet *DungeonLocationRecordSet,
) (*record.DungeonAction, error) {

	// TODO: Register the number of times the character has looked at any other
	// location, object, monster of character. Looking at anything multiple times
	// should result in additional information within X turns.

	if dungeonActionRec.ResolvedTargetDungeonLocationID.Valid {
		m.Log.Info("Looking at location ID >%s<", dungeonActionRec.ResolvedTargetDungeonLocationID.String)

	} else if dungeonActionRec.ResolvedTargetDungeonObjectID.Valid {
		m.Log.Info("Looking at object ID >%s<", dungeonActionRec.ResolvedTargetDungeonObjectID.String)

	} else if dungeonActionRec.ResolvedTargetDungeonMonsterID.Valid {
		m.Log.Info("Looking at monster ID >%s<", dungeonActionRec.ResolvedTargetDungeonMonsterID.String)

	} else if dungeonActionRec.ResolvedTargetDungeonCharacterID.Valid {
		m.Log.Info("Looking at character ID >%s<", dungeonActionRec.ResolvedTargetDungeonCharacterID.String)
	}

	return dungeonActionRec, nil
}

func (m *Model) performDungeonActionStash(
	dungeonCharacterRec *record.DungeonCharacter,
	dungeonActionRec *record.DungeonAction,
	dungeonLocationRecordSet *DungeonLocationRecordSet,
) (*record.DungeonAction, error) {

	if dungeonActionRec.DungeonCharacterID.Valid {
		// Character stash object
		dungeonObjectID := dungeonActionRec.ResolvedStashedDungeonObjectID.String
		if dungeonObjectID == "" {
			msg := "resolved stashed dungeon object ID is empty, cannot stash object"
			m.Log.Warn(msg)
			return nil, fmt.Errorf(msg)
		}

		dungeonObjectRec, err := m.GetDungeonObjectRec(dungeonObjectID, true)
		if err != nil {
			m.Log.Warn("Failed getting dungeon object record >%v<", err)
			return nil, err
		}

		dungeonObjectRec.DungeonLocationID = sql.NullString{}
		dungeonObjectRec.DungeonCharacterID = dungeonActionRec.DungeonCharacterID
		dungeonObjectRec.IsStashed = true

		err = m.UpdateDungeonObjectRec(dungeonObjectRec)
		if err != nil {
			m.Log.Warn("Failed updating dungeon object record >%v<", err)
			return nil, err
		}

	} else if dungeonActionRec.DungeonMonsterID.Valid {
		// Monster stash object
		return nil, fmt.Errorf("monsters stashing objects is currently not supported")
	}

	return dungeonActionRec, nil
}
