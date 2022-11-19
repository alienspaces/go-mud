package model

import (
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

// GetMonsterEquippedObjectRecs -
func (m *Model) GetMonsterEquippedObjectRecs(monsterID string) ([]*record.Object, error) {

	l := m.Logger("GetMonsterEquippedObjectRecs")

	l.Debug("Getting monster ID >%s< equipped object records", monsterID)

	mor := m.MonsterObjectRepository()
	or := m.ObjectRepository()

	monsterObjectRecs, err := mor.GetMany(map[string]interface{}{
		"monster_id":  monsterID,
		"is_equipped": true,
	}, nil, false)
	if err != nil {
		return nil, err
	}

	objectRecs := []*record.Object{}
	for _, monsterObjectRec := range monsterObjectRecs {
		objectRec, err := or.GetOne(monsterObjectRec.ObjectID, false)
		if err != nil {
			return nil, err
		}
		objectRecs = append(objectRecs, objectRec)
	}

	return objectRecs, nil
}

// GetMonsterStashedObjectRecs -
func (m *Model) GetMonsterStashedObjectRecs(monsterID string) ([]*record.Object, error) {

	l := m.Logger("GetMonsterStashedObjectRecs")

	l.Debug("Getting monster ID >%s< stashed object records", monsterID)

	mor := m.MonsterObjectRepository()
	or := m.ObjectRepository()

	monsterObjectRecs, err := mor.GetMany(map[string]interface{}{
		"monster_id": monsterID,
		"is_stashed": true,
	}, nil, false)
	if err != nil {
		return nil, err
	}

	objectRecs := []*record.Object{}
	for _, monsterObjectRec := range monsterObjectRecs {
		objectRec, err := or.GetOne(monsterObjectRec.ObjectID, false)
		if err != nil {
			return nil, err
		}
		objectRecs = append(objectRecs, objectRec)
	}

	return objectRecs, nil
}
