package model

import (
	coresql "gitlab.com/alienspaces/go-mud/backend/core/sql"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

// DeleteMonsterInstance -
func (m *Model) DeleteMonsterInstance(monsterID string) error {
	l := m.loggerWithFunctionContext("DeleteMonsterInstance")

	monsterInstanceRec, err := m.GetMonsterInstanceViewRecByMonsterID(monsterID)
	if err != nil {
		l.Warn("failed getting monster instance view record >%v<", err)
		return err
	}

	monsterObjectInstanceRecs, err := m.GetMonsterInstanceObjectInstanceRecs(monsterInstanceRec.MonsterID)
	if err != nil {
		l.Warn("failed getting monster object instance view records >%v<", err)
		return err
	}

	for idx := range monsterObjectInstanceRecs {
		err := m.DeleteObjectInstanceRec(monsterObjectInstanceRecs[idx].ID)
		if err != nil {
			l.Warn("failed deleting monster object instance record >%v<", err)
			return err
		}
	}

	err = m.DeleteMonsterInstanceRec(monsterInstanceRec.ID)
	if err != nil {
		l.Warn("failed deleting monster instance record >%v<", err)
		return err
	}

	return nil
}

// GetMonsterInstanceObjectInstanceRecs -
func (m *Model) GetMonsterInstanceObjectInstanceRecs(monsterID string) ([]*record.ObjectInstance, error) {

	m.Log.Info("Getting monster ID >%s< object records", monsterID)

	r := m.ObjectInstanceRepository()

	return r.GetMany(
		&coresql.Options{
			Params: []coresql.Param{
				{
					Col: "monster_instance_id",
					Val: monsterID,
				},
			},
		},
	)
}

// GetMonsterInstanceEquippedObjectInstanceRecs -
func (m *Model) GetMonsterInstanceEquippedObjectInstanceRecs(monsterID string) ([]*record.ObjectInstance, error) {

	m.Log.Info("Getting monster ID >%s< equipped object records", monsterID)

	r := m.ObjectInstanceRepository()

	return r.GetMany(
		&coresql.Options{
			Params: []coresql.Param{
				{
					Col: "monster_instance_id",
					Val: monsterID,
				},
				{
					Col: "is_equipped",
					Val: true,
				},
			},
		},
	)
}

// GetMonsterInstanceStashedObjectInstanceRecs -
func (m *Model) GetMonsterInstanceStashedObjectInstanceRecs(monsterID string) ([]*record.ObjectInstance, error) {

	m.Log.Info("Getting monster ID >%s< stashed object records", monsterID)

	r := m.ObjectInstanceRepository()

	return r.GetMany(
		&coresql.Options{
			Params: []coresql.Param{
				{
					Col: "monster_instance_id",
					Val: monsterID,
				},
				{
					Col: "is_stashed",
					Val: true,
				},
			},
		},
	)
}

// GetMonsterInstanceObjectInstanceViewRecs -
func (m *Model) GetMonsterInstanceObjectInstanceViewRecs(monsterID string) ([]*record.ObjectInstanceView, error) {

	m.Log.Info("Getting monster ID >%s< object records", monsterID)

	r := m.ObjectInstanceViewRepository()

	return r.GetMany(
		&coresql.Options{
			Params: []coresql.Param{
				{
					Col: "monster_instance_id",
					Val: monsterID,
				},
			},
		},
	)
}

// GetMonsterInstanceEquippedObjectInstanceViewRecs -
func (m *Model) GetMonsterInstanceEquippedObjectInstanceViewRecs(monsterID string) ([]*record.ObjectInstanceView, error) {

	m.Log.Info("Getting monster ID >%s< equipped object records", monsterID)

	r := m.ObjectInstanceViewRepository()

	return r.GetMany(
		&coresql.Options{
			Params: []coresql.Param{
				{
					Col: "monster_instance_id",
					Val: monsterID,
				},
				{
					Col: "is_equipped",
					Val: true,
				},
			},
		},
	)
}

// GetMonsterInstanceStashedObjectInstanceViewRecs -
func (m *Model) GetMonsterInstanceStashedObjectInstanceViewRecs(monsterID string) ([]*record.ObjectInstanceView, error) {

	m.Log.Info("Getting monster ID >%s< stashed object records", monsterID)

	r := m.ObjectInstanceViewRepository()

	return r.GetMany(
		&coresql.Options{
			Params: []coresql.Param{
				{
					Col: "monster_instance_id",
					Val: monsterID,
				},
				{
					Col: "is_stashed",
					Val: true,
				},
			},
		},
	)
}
