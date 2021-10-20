package model

import (
	"fmt"

	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

func (m *Model) performDungeonCharacterAction(
	dungeonActionRec *record.DungeonAction,
	dungeonLocationRecordSet *DungeonLocationRecordSet,
) (*record.DungeonAction, error) {

	actionFuncs := map[string]func(dungeonActionRecord *record.DungeonAction, dungeonLocationRecordSet *DungeonLocationRecordSet) (*record.DungeonAction, error){
		"move": m.performDungeonActionMove,
	}

	// const actionFuncs = {
	// 	move: (dungeonActionRecord: DungeonActionRepositoryRecord, records: DungeonLocationRecordSet) =>
	// 		this.performDungeonActionMove(dungeonActionRecord, records),
	// 	look: (dungeonActionRecord: DungeonActionRepositoryRecord, records: DungeonLocationRecordSet) =>
	// 		this.performDungeonActionLook(dungeonActionRecord, records),
	// 	equip: (dungeonActionRecord: DungeonActionRepositoryRecord, records: DungeonLocationRecordSet) =>
	// 		this.performDungeonActionEquip(dungeonActionRecord, records),
	// 	stash: (dungeonActionRecord: DungeonActionRepositoryRecord, records: DungeonLocationRecordSet) =>
	// 		this.performDungeonActionStash(dungeonActionRecord, records),
	// 	drop: (dungeonActionRecord: DungeonActionRepositoryRecord, records: DungeonLocationRecordSet) =>
	// 		this.performDungeonActionDrop(dungeonActionRecord, records),
	// };

	actionFunc, ok := actionFuncs[dungeonActionRec.ResolvedCommand]
	if !ok {
		msg := fmt.Sprintf("Action >%s< not supported", dungeonActionRec.ResolvedCommand)
		m.Log.Warn(msg)
		return nil, fmt.Errorf(msg)
	}

	var err error
	dungeonActionRec, err = actionFunc(dungeonActionRec, dungeonLocationRecordSet)
	if err != nil {
		m.Log.Warn("Failed performing action >%s< >%v<", dungeonActionRec.ResolvedCommand, err)
		return nil, err
	}

	m.Log.Info("Have updated dungeon action record >%v<", dungeonActionRec)

	return dungeonActionRec, nil
}

func (m *Model) performDungeonActionMove(
	dungeonActionRec *record.DungeonAction,
	dungeonLocationRecordSet *DungeonLocationRecordSet,
) (*record.DungeonAction, error) {

	if dungeonActionRec.DungeonCharacterID != "" {
		// Move character
		characterRec := dungeonLocationRecordSet.CharacterRec
		characterRec.DungeonLocationID = dungeonActionRec.ResolvedTargetDungeonLocationID

		err := m.UpdateDungeonCharacterRec(characterRec)
		if err != nil {
			m.Log.Warn("Failed updated dungeon character record >%v<", err)
			return nil, err
		}

		// Update dungeon action record
		dungeonActionRec.DungeonLocationID = dungeonActionRec.ResolvedTargetDungeonLocationID
	} else if dungeonActionRec.DungeonMonsterID != "" {
		// Move monster
		return nil, fmt.Errorf("moving monsters is currently not supported")
	}

	return dungeonActionRec, nil
}
