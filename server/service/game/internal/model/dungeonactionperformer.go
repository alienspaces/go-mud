package model

import (
	"fmt"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

func (m *Model) performDungeonCharacterAction(
	dungeonCharacterRec *record.DungeonCharacter,
	dungeonActionRec *record.DungeonAction,
	dungeonLocationRecordSet *DungeonLocationRecordSet,
) (*record.DungeonAction, error) {

	actionFuncs := map[string]func(dungeonCharacterRec *record.DungeonCharacter, dungeonActionRec *record.DungeonAction, dungeonLocationRecordSet *DungeonLocationRecordSet) (*record.DungeonAction, error){
		"move": m.performDungeonActionMove,
		"look": m.performDungeonActionLook,
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

	if dungeonActionRec.DungeonCharacterID.Valid {
		// TODO: Looking at anything multiple times should result in additional information within a short timeframe

	} else if dungeonActionRec.DungeonMonsterID.Valid {
		// TODO: Monster look at current room, a direction, another monster, a character, or an object
		return nil, fmt.Errorf("monsters looking is currently not supported")
	}

	return dungeonActionRec, nil
}
