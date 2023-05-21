package model

import (
	coresql "gitlab.com/alienspaces/go-mud/backend/core/sql"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

// GetMonsterEquippedObjectRecs -
func (m *Model) GetMonsterEquippedObjectRecs(monsterID string) ([]*record.Object, error) {

	l := m.Logger("GetMonsterEquippedObjectRecs")

	l.Debug("Getting monster ID >%s< equipped object records", monsterID)

	mor := m.MonsterObjectRepository()
	or := m.ObjectRepository()

	monsterObjectRecs, err := mor.GetMany(
		&coresql.Options{
			Params: []coresql.Param{
				{
					Col: "monster_id",
					Val: monsterID,
				},
				{
					Col: "is_equipped",
					Val: true,
				},
			},
		},
	)
	if err != nil {
		return nil, err
	}

	objectRecs := []*record.Object{}
	for _, monsterObjectRec := range monsterObjectRecs {
		objectRec, err := or.GetOne(monsterObjectRec.ObjectID, nil)
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

	monsterObjectRecs, err := mor.GetMany(
		&coresql.Options{
			Params: []coresql.Param{
				{
					Col: "monster_id",
					Val: monsterID,
				},
				{
					Col: "is_stashed",
					Val: true,
				},
			},
		},
	)
	if err != nil {
		return nil, err
	}

	objectRecs := []*record.Object{}
	for _, monsterObjectRec := range monsterObjectRecs {
		objectRec, err := or.GetOne(monsterObjectRec.ObjectID, nil)
		if err != nil {
			return nil, err
		}
		objectRecs = append(objectRecs, objectRec)
	}

	return objectRecs, nil
}
