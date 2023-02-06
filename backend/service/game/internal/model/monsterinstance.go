package model

import (
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

// GetMonsterInstanceObjectInstanceRecs -
func (m *Model) GetMonsterInstanceObjectInstanceRecs(monsterID string) ([]*record.ObjectInstance, error) {

	m.Log.Info("Getting monster ID >%s< object records", monsterID)

	r := m.ObjectInstanceRepository()

	return r.GetMany(map[string]interface{}{
		"monster_instance_id": monsterID,
	}, nil, false)
}

// GetMonsterInstanceEquippedObjectInstanceRecs -
func (m *Model) GetMonsterInstanceEquippedObjectInstanceRecs(monsterID string) ([]*record.ObjectInstance, error) {

	m.Log.Info("Getting monster ID >%s< equipped object records", monsterID)

	r := m.ObjectInstanceRepository()

	return r.GetMany(map[string]interface{}{
		"monster_instance_id": monsterID,
		"is_equipped":         true,
	}, nil, false)
}

// GetMonsterInstanceStashedObjectInstanceRecs -
func (m *Model) GetMonsterInstanceStashedObjectInstanceRecs(monsterID string) ([]*record.ObjectInstance, error) {

	m.Log.Info("Getting monster ID >%s< stashed object records", monsterID)

	r := m.ObjectInstanceRepository()

	return r.GetMany(map[string]interface{}{
		"monster_instance_id": monsterID,
		"is_stashed":          true,
	}, nil, false)
}

// GetMonsterInstanceObjectInstanceViewRecs -
func (m *Model) GetMonsterInstanceObjectInstanceViewRecs(monsterID string) ([]*record.ObjectInstanceView, error) {

	m.Log.Info("Getting monster ID >%s< object records", monsterID)

	r := m.ObjectInstanceViewRepository()

	return r.GetMany(map[string]interface{}{
		"monster_instance_id": monsterID,
	}, nil, false)
}

// GetMonsterInstanceEquippedObjectInstanceViewRecs -
func (m *Model) GetMonsterInstanceEquippedObjectInstanceViewRecs(monsterID string) ([]*record.ObjectInstanceView, error) {

	m.Log.Info("Getting monster ID >%s< equipped object records", monsterID)

	r := m.ObjectInstanceViewRepository()

	return r.GetMany(map[string]interface{}{
		"monster_instance_id": monsterID,
		"is_equipped":         true,
	}, nil, false)
}

// GetMonsterInstanceStashedObjectInstanceViewRecs -
func (m *Model) GetMonsterInstanceStashedObjectInstanceViewRecs(monsterID string) ([]*record.ObjectInstanceView, error) {

	m.Log.Info("Getting monster ID >%s< stashed object records", monsterID)

	r := m.ObjectInstanceViewRepository()

	return r.GetMany(map[string]interface{}{
		"monster_instance_id": monsterID,
		"is_stashed":          true,
	}, nil, false)
}
